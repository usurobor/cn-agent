// Package pkgbuild implements the cn build pipeline per
// BUILD-AND-DIST.md: validate + package src/packages/ into
// dist/packages/ as distributable tarballs.
//
// Three modes:
//   - Build: validate src/packages/<name>/ → produce dist/packages/<name>-<version>.tar.gz
//   - Check: validate src/packages/<name>/ structurally (fast CI path, no artifacts)
//   - Clean: remove dist/packages/ artifacts
//
// Source: src/packages/<name>/ (self-contained package with cn.package.json + content)
// Output: dist/packages/<name>-<version>.tar.gz + index.json + checksums.txt
package pkgbuild

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	pkgtypes "github.com/usurobor/cnos/src/go/internal/pkg"
)

// FindContentClasses returns the subset of pkgtypes.ContentClasses
// that are present as non-empty directories inside pkgDir. This is
// the shared filesystem predicate used by `cn build --check` (validating
// sources under src/packages/) and `cn status` (inspecting installed
// packages under .cn/vendor/packages/). Presence is the single
// authority; manifest JSON fields are not consulted.
//
// Order follows pkgtypes.ContentClasses. An error walking an entry
// is treated as absence — the caller gets a conservative "not present"
// rather than a partial failure.
func FindContentClasses(pkgDir string) []string {
	var present []string
	for _, class := range pkgtypes.ContentClasses {
		classDir := filepath.Join(pkgDir, class)
		info, err := os.Stat(classDir)
		if err != nil || !info.IsDir() {
			continue
		}
		entries, err := os.ReadDir(classDir)
		if err != nil || len(entries) == 0 {
			continue
		}
		present = append(present, class)
	}
	return present
}

// PackageManifest is the parsed cn.package.json.
type PackageManifest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ParseManifestData parses cn.package.json from raw bytes.
// Pure — no IO. Mirrors the Parse* (pure) vs Read* (IO) pattern
// from internal/pkg/.
func ParseManifestData(data []byte) (*PackageManifest, error) {
	var m PackageManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	if m.Name == "" || m.Version == "" {
		return nil, fmt.Errorf("manifest missing name or version")
	}
	return &m, nil
}

// ReadManifest reads and parses a cn.package.json file from disk.
// IO wrapper around ParseManifestData.
func ReadManifest(path string) (*PackageManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return ParseManifestData(data)
}

// FindRepoRoot walks up from cwd looking for .git.
func FindRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getwd: %w", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not inside a git repository")
		}
		dir = parent
	}
}

// DiscoveredPackage pairs a directory name with its parsed manifest.
type DiscoveredPackage struct {
	DirName  string
	SrcDir   string // absolute path to src/packages/<name>/
	Manifest *PackageManifest
}

// DiscoverPackages finds all package manifests under src/packages/.
func DiscoverPackages(repoRoot string) ([]DiscoveredPackage, error) {
	srcPkgsDir := filepath.Join(repoRoot, "src", "packages")
	entries, err := os.ReadDir(srcPkgsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read src/packages/: %w", err)
	}

	var packages []DiscoveredPackage
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		manifestPath := filepath.Join(srcPkgsDir, e.Name(), "cn.package.json")
		m, err := ReadManifest(manifestPath)
		if err != nil {
			continue
		}
		packages = append(packages, DiscoveredPackage{
			DirName:  e.Name(),
			SrcDir:   filepath.Join(srcPkgsDir, e.Name()),
			Manifest: m,
		})
	}

	sort.Slice(packages, func(i, j int) bool {
		return packages[i].DirName < packages[j].DirName
	})
	return packages, nil
}

// --- Build mode ---

// BuildResult records the outcome of building one package.
type BuildResult struct {
	Name        string
	Version     string
	TarballPath string
	SHA256      string
	Err         error
}

// BuildOne validates and packages one src/packages/<name>/ into
// dist/packages/<name>-<version>.tar.gz.
func BuildOne(repoRoot string, pkg DiscoveredPackage) BuildResult {
	r := BuildResult{Name: pkg.Manifest.Name, Version: pkg.Manifest.Version}

	distDir := filepath.Join(repoRoot, "dist", "packages")
	if err := os.MkdirAll(distDir, 0755); err != nil {
		r.Err = fmt.Errorf("create dist dir: %w", err)
		return r
	}

	tarballName := fmt.Sprintf("%s-%s.tar.gz", pkg.Manifest.Name, pkg.Manifest.Version)
	tarballPath := filepath.Join(distDir, tarballName)

	hash, err := createTarGz(tarballPath, pkg.SrcDir)
	if err != nil {
		r.Err = fmt.Errorf("create tarball: %w", err)
		return r
	}

	r.TarballPath = tarballPath
	r.SHA256 = hash
	return r
}

// createTarGz creates a .tar.gz from srcDir.
// Returns the hex SHA-256 of the resulting file.
func createTarGz(dest, srcDir string) (string, error) {
	f, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("create %s: %w", dest, err)
	}
	defer f.Close()

	h := sha256.New()
	mw := io.MultiWriter(f, h)

	gw := gzip.NewWriter(mw)
	tw := tar.NewWriter(gw)

	err = filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = rel
		if d.IsDir() {
			hdr.Name += "/"
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if !d.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if _, err := tw.Write(data); err != nil {
				return err
			}
		}
		return nil
	})

	// Close in order: tar → gzip → file. Check errors — gzip
	// finalization writes the footer; a flush failure means a corrupt tarball.
	if err := tw.Close(); err != nil {
		return "", fmt.Errorf("finalize tar: %w", err)
	}
	if err := gw.Close(); err != nil {
		return "", fmt.Errorf("finalize gzip: %w", err)
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("close file: %w", err)
	}

	if err != nil {
		return "", fmt.Errorf("walk %s: %w", srcDir, err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// UpdateIndex writes/updates dist/packages/index.json with build results.
func UpdateIndex(repoRoot string, results []BuildResult) error {
	distDir := filepath.Join(repoRoot, "dist", "packages")
	if err := os.MkdirAll(distDir, 0755); err != nil {
		return err
	}
	indexPath := filepath.Join(distDir, "index.json")

	type indexEntry struct {
		URL    string `json:"url"`
		SHA256 string `json:"sha256"`
	}
	type packageIndex struct {
		Schema   string                           `json:"schema"`
		Packages map[string]map[string]indexEntry `json:"packages"`
	}

	idx := packageIndex{
		Schema:   "cn.package-index.v1",
		Packages: make(map[string]map[string]indexEntry),
	}

	// Read existing index if present.
	if data, err := os.ReadFile(indexPath); err == nil {
		json.Unmarshal(data, &idx)
		if idx.Packages == nil {
			idx.Packages = make(map[string]map[string]indexEntry)
		}
	}

	for _, r := range results {
		if r.Err != nil {
			continue
		}
		if idx.Packages[r.Name] == nil {
			idx.Packages[r.Name] = make(map[string]indexEntry)
		}
		// URL is relative — release workflow makes it absolute.
		idx.Packages[r.Name][r.Version] = indexEntry{
			URL:    fmt.Sprintf("%s-%s.tar.gz", r.Name, r.Version),
			SHA256: r.SHA256,
		}
	}

	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal index: %w", err)
	}
	return os.WriteFile(indexPath, append(data, '\n'), 0644)
}

// UpdateChecksums writes dist/packages/checksums.txt.
func UpdateChecksums(repoRoot string, results []BuildResult) error {
	distDir := filepath.Join(repoRoot, "dist", "packages")
	if err := os.MkdirAll(distDir, 0755); err != nil {
		return err
	}
	checksumsPath := filepath.Join(distDir, "checksums.txt")

	var lines []string
	for _, r := range results {
		if r.Err != nil {
			continue
		}
		tarball := fmt.Sprintf("%s-%s.tar.gz", r.Name, r.Version)
		lines = append(lines, fmt.Sprintf("%s  %s", r.SHA256, tarball))
	}

	sort.Strings(lines)
	content := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(checksumsPath, []byte(content), 0644)
}

// --- Check mode ---

// CheckResult records validation status for one package.
type CheckResult struct {
	Name   string
	Issues []string
}

// CheckOne validates that src/packages/<name>/ is structurally valid.
func CheckOne(pkg DiscoveredPackage) CheckResult {
	r := CheckResult{Name: pkg.Manifest.Name}

	if pkg.Manifest.Name == "" {
		r.Issues = append(r.Issues, "missing name in cn.package.json")
	}
	if pkg.Manifest.Version == "" {
		r.Issues = append(r.Issues, "missing version in cn.package.json")
	}

	hasContent := false
	for _, class := range pkgtypes.ContentClasses {
		classDir := filepath.Join(pkg.SrcDir, class)
		if info, err := os.Stat(classDir); err == nil && info.IsDir() {
			hasContent = true
			entries, err := os.ReadDir(classDir)
			if err == nil && len(entries) == 0 {
				r.Issues = append(r.Issues, fmt.Sprintf("empty content class: %s/", class))
			}
		}
	}
	if !hasContent {
		r.Issues = append(r.Issues, "no content class directories found")
	}

	return r
}

// --- Lockfile generation ---

// GenerateLockfileData produces a cn.lock.v2 lockfile from build results.
// Uses canonical pkgtypes.Lockfile/pkgtypes.LockedDep types — single
// schema definition. Pure — no IO.
func GenerateLockfileData(results []BuildResult) ([]byte, error) {
	lf := pkgtypes.Lockfile{Schema: "cn.lock.v2"}
	for _, r := range results {
		if r.Err != nil {
			continue
		}
		lf.Packages = append(lf.Packages, pkgtypes.LockedDep{
			Name:    r.Name,
			Version: r.Version,
			SHA256:  r.SHA256,
		})
	}

	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal lockfile: %w", err)
	}
	return append(data, '\n'), nil
}

// GenerateLockfile writes a lockfile to dest from build results.
func GenerateLockfile(dest string, results []BuildResult) error {
	data, err := GenerateLockfileData(results)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("create lockfile dir: %w", err)
	}
	return os.WriteFile(dest, data, 0644)
}

// --- Clean mode ---

// Clean removes dist/packages/ artifacts.
func Clean(repoRoot string) error {
	distDir := filepath.Join(repoRoot, "dist", "packages")
	return os.RemoveAll(distDir)
}
