// Package restore implements the cn deps restore flow in Go:
// lockfile → index lookup → fetch (HTTP or local) → SHA-256 verify →
// tar extract → validate cn.package.json.
//
// Also provides lockfile generation from a package index
// (GenerateLockFromIndex) for the `cn deps lock` command.
//
// Fetch supports both HTTP URLs (remote index) and relative file paths
// (local dev workflow: cn build → dist/ → cn deps restore).
package restore

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/usurobor/cnos/src/go/internal/pkg"
)

// FindIndexPath walks up from hubPath looking for dist/packages/index.json.
func FindIndexPath(hubPath string) string {
	dir := hubPath
	for {
		candidate := filepath.Join(dir, "dist", "packages", "index.json")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return filepath.Join(hubPath, "dist", "packages", "index.json")
		}
		dir = parent
	}
}

// --- IO wrappers for pure parsers (mirrors OCaml src/ocaml/cmd/ vs src/ocaml/lib/ split) ---

// ReadLockfile reads and parses a lockfile from disk.
func ReadLockfile(path string) (*pkg.Lockfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read lockfile %s: %w", path, err)
	}
	return pkg.ParseLockfile(data)
}

// ReadManifest reads and parses a deps.json manifest from disk.
func ReadManifest(path string) (*pkg.Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read manifest %s: %w", path, err)
	}
	return pkg.ParseManifest(data)
}

// ReadPackageIndex reads and parses a package index from disk.
func ReadPackageIndex(path string) (*pkg.PackageIndex, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read package index %s: %w", path, err)
	}
	return pkg.ParsePackageIndex(data)
}

// ValidatePackageManifest reads cn.package.json from pkgDir and
// validates that it declares a name matching expectedName.
func ValidatePackageManifest(pkgDir, expectedName string) error {
	path := filepath.Join(pkgDir, "cn.package.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("missing cn.package.json: %w", err)
	}
	return pkg.ValidatePackageManifestData(data, expectedName)
}

// Result records the outcome of restoring one package.
type Result struct {
	Name    string
	Version string
	Err     error // nil = success
}

// Restore installs every package declared in the lockfile into
// .cn/vendor/packages/. It collects all errors rather than stopping
// at the first failure.
//
// indexPath is the path to the package index file. When the index
// contains relative URLs (e.g. "cnos.core-1.0.0.tar.gz" from cn build),
// they are resolved relative to the directory containing indexPath.
func Restore(ctx context.Context, hubPath, indexPath string) ([]Result, error) {
	lockPath := filepath.Join(hubPath, ".cn", "deps.lock.json")
	lf, err := ReadLockfile(lockPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil // no lockfile = nothing to restore
		}
		return nil, fmt.Errorf("read lockfile: %w", err)
	}
	if len(lf.Packages) == 0 {
		return nil, nil
	}

	idx, err := ReadPackageIndex(indexPath)
	if err != nil {
		return nil, fmt.Errorf("load package index: %w", err)
	}

	indexDir := filepath.Dir(indexPath)

	var results []Result
	for _, dep := range lf.Packages {
		r := restoreOne(ctx, hubPath, idx, dep, indexDir)
		results = append(results, r)
	}
	return results, nil
}

// restoreOne installs a single package from its lockfile entry.
// indexDir is the directory containing the package index; relative
// URLs in the index are resolved against it.
func restoreOne(ctx context.Context, hubPath string, idx *pkg.PackageIndex, dep pkg.LockedDep, indexDir string) Result {
	r := Result{Name: dep.Name, Version: dep.Version}

	pkgDir := pkg.VendorPath(hubPath, dep.Name)

	// Already installed — skip.
	// NOTE: With version-less VendorPath, this check does not detect
	// version upgrades (v1 installed, lockfile updated to v2 → directory
	// exists from v1, skip fires). Version upgrade requires comparing
	// the installed cn.package.json version against the lockfile version,
	// or comparing SHA-256 of installed state. Tracked in #230.
	if _, err := os.Stat(pkgDir); err == nil {
		slog.DebugContext(ctx, "package already installed, skipping",
			slog.String("package", dep.Name),
			slog.String("version", dep.Version))
		return r
	}

	// Lookup in index.
	entry := idx.Lookup(dep.Name, dep.Version)
	if entry == nil {
		r.Err = fmt.Errorf("package %s@%s not in index", dep.Name, dep.Version)
		return r
	}

	// Download tarball to tmp.
	tmpTar := pkg.TmpTarballPath(hubPath, dep.Name, dep.Version)
	if err := os.MkdirAll(filepath.Dir(tmpTar), 0755); err != nil {
		r.Err = fmt.Errorf("create tmp dir: %w", err)
		return r
	}
	defer func() {
		// Clean up tmp tarball regardless of outcome.
		os.Remove(tmpTar)
	}()

	if err := fetchTarball(ctx, entry.URL, tmpTar, indexDir); err != nil {
		r.Err = fmt.Errorf("fetch %s@%s from %s: %w",
			dep.Name, dep.Version, entry.URL, err)
		return r
	}

	// SHA-256 verification against the lockfile entry.
	actual, err := fileSHA256(tmpTar)
	if err != nil {
		r.Err = fmt.Errorf("compute sha256 for %s@%s: %w", dep.Name, dep.Version, err)
		return r
	}
	if actual != dep.SHA256 {
		r.Err = fmt.Errorf("sha256 mismatch for %s@%s: expected %s, got %s",
			dep.Name, dep.Version, dep.SHA256, actual)
		return r
	}

	// Extract tarball.
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		r.Err = fmt.Errorf("create vendor dir: %w", err)
		return r
	}
	if err := extractTarGz(tmpTar, pkgDir); err != nil {
		os.RemoveAll(pkgDir) // clean up partial extraction
		r.Err = fmt.Errorf("extract %s@%s: %w", dep.Name, dep.Version, err)
		return r
	}

	// Validate cn.package.json.
	if err := ValidatePackageManifest(pkgDir, dep.Name); err != nil {
		os.RemoveAll(pkgDir) // invalid package, remove
		r.Err = fmt.Errorf("validate %s@%s: %w", dep.Name, dep.Version, err)
		return r
	}

	slog.InfoContext(ctx, "restored package",
		slog.String("package", dep.Name),
		slog.String("version", dep.Version))
	return r
}

// isRemoteURL returns true if the URL has an http:// or https:// scheme.
func isRemoteURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// fetchTarball retrieves a tarball to dest. If the URL is remote (http/https),
// it downloads via HTTP. If the URL is a relative path, it resolves against
// indexDir and copies locally. This enables the local dev workflow:
// cn build → dist/packages/ → cn deps restore without a server.
func fetchTarball(ctx context.Context, url, dest, indexDir string) error {
	if isRemoteURL(url) {
		return downloadFile(ctx, url, dest)
	}
	// Local file: resolve relative to the index directory.
	src := url
	if !filepath.IsAbs(url) {
		src = filepath.Join(indexDir, url)
	}
	return copyFile(src, dest)
}

// copyFile copies src to dest. Used for local package restore.
func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open %s: %w", src, err)
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create %s: %w", dest, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	return out.Close()
}

// httpClient is the shared HTTP client for package downloads.
// Timeout mirrors OCaml's curl flags: --connect-timeout 10 --max-time 300.
var httpClient = &http.Client{
	Timeout: 300 * time.Second,
}

// downloadFile fetches url and writes it to dest via HTTP.
func downloadFile(ctx context.Context, url, dest string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status %d", resp.StatusCode)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return f.Close()
}

// fileSHA256 computes the hex-encoded SHA-256 of a file.
func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// extractTarGz extracts a .tar.gz archive into destDir.
// It validates paths to prevent directory traversal.
func extractTarGz(tarball, destDir string) error {
	f, err := os.Open(tarball)
	if err != nil {
		return fmt.Errorf("open tarball: %w", err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("gzip reader: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar next: %w", err)
		}

		// Security: prevent directory traversal. The separator suffix
		// ensures /tmp/foo doesn't prefix-match /tmp/foobar.
		target := filepath.Join(destDir, hdr.Name)
		cleanDest := filepath.Clean(destDir) + string(filepath.Separator)
		if target != filepath.Clean(destDir) && !strings.HasPrefix(filepath.Clean(target), cleanDest) {
			return fmt.Errorf("tar entry %q escapes dest dir", hdr.Name)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("mkdir %s: %w", target, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("mkdir parent %s: %w", target, err)
			}
			out, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("create %s: %w", target, err)
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return fmt.Errorf("write %s: %w", target, err)
			}
			out.Close()
			if err := os.Chmod(target, hdr.FileInfo().Mode()); err != nil {
				return fmt.Errorf("chmod %s: %w", target, err)
			}
		}
	}
	return nil
}

// HasErrors returns true if any result has a non-nil error.
func HasErrors(results []Result) bool {
	for _, r := range results {
		if r.Err != nil {
			return true
		}
	}
	return false
}

// Errors returns only the failed results.
func Errors(results []Result) []Result {
	var errs []Result
	for _, r := range results {
		if r.Err != nil {
			errs = append(errs, r)
		}
	}
	return errs
}

// LockResult records the outcome of lockfile generation.
type LockResult struct {
	LockPath string
	Count    int
}

// GenerateLockFromIndex reads the hub's .cn/deps.json and produces a
// lockfile pinning exactly the packages the manifest declares, resolving
// each (name, version) pin against the package index. Writes the result
// to <hubPath>/.cn/deps.lock.json.
//
// This is the domain logic for `cn deps lock` — CLI wiring calls this,
// keeping cmd_deps.go thin. Uses canonical pkg.Lockfile/pkg.LockedDep
// types. Output is sorted by name then version for deterministic content.
//
// Errors:
//   - .cn/deps.json missing or unparseable → cannot lock without a manifest
//   - package index missing or unparseable → cannot resolve pins
//   - any pinned (name, version) not present in the index → reported in
//     a single error listing every missing pin
//
// Today's resolution semantics are exact-match: the version string in
// deps.json must equal a key under idx.Packages[name]. Semver range
// resolution is a separate concern (issue #250 out-of-scope).
func GenerateLockFromIndex(hubPath, indexPath string) (*LockResult, error) {
	manifestPath := filepath.Join(hubPath, ".cn", "deps.json")
	m, err := ReadManifest(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("read deps.json: %w", err)
	}

	idx, err := ReadPackageIndex(indexPath)
	if err != nil {
		return nil, fmt.Errorf("read index: %w", err)
	}

	lockPath := filepath.Join(hubPath, ".cn", "deps.lock.json")
	if err := os.MkdirAll(filepath.Dir(lockPath), 0755); err != nil {
		return nil, fmt.Errorf("create lockfile dir: %w", err)
	}

	lf := pkg.Lockfile{Schema: "cn.lock.v2"}
	var missing []string
	for _, dep := range m.Packages {
		entry := idx.Lookup(dep.Name, dep.Version)
		if entry == nil {
			missing = append(missing, fmt.Sprintf("%s@%s", dep.Name, dep.Version))
			continue
		}
		lf.Packages = append(lf.Packages, pkg.LockedDep{
			Name:    dep.Name,
			Version: dep.Version,
			SHA256:  entry.SHA256,
		})
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("packages pinned in deps.json not in index: %s",
			strings.Join(missing, ", "))
	}

	// Sort for deterministic output — manifest order is operator-controlled
	// and lockfile content must be reproducible regardless.
	sort.Slice(lf.Packages, func(i, j int) bool {
		if lf.Packages[i].Name != lf.Packages[j].Name {
			return lf.Packages[i].Name < lf.Packages[j].Name
		}
		return lf.Packages[i].Version < lf.Packages[j].Version
	})

	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal lockfile: %w", err)
	}
	if err := os.WriteFile(lockPath, append(data, '\n'), 0644); err != nil {
		return nil, fmt.Errorf("write lockfile: %w", err)
	}

	return &LockResult{LockPath: lockPath, Count: len(lf.Packages)}, nil
}
