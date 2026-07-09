package repoinstall

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/pkg"
)

// --- fixture helpers ---

// makeTarGz creates a .tar.gz in memory containing the given files
// (map of relative-path → content), writes it to a temp file, and
// returns the raw bytes + hex SHA-256. Mirrors restore_test.go's helper
// of the same name/shape (test-local duplication of a fixture builder,
// not a production parser — eng/go §2.17 governs parser dedup, not test
// fixture helpers).
func makeTarGz(t *testing.T, files map[string]string) ([]byte, string) {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	for name, content := range files {
		hdr := &tar.Header{Name: name, Mode: 0644, Size: int64(len(content))}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	tw.Close()
	gw.Close()
	data := buf.Bytes()
	h := sha256.Sum256(data)
	return data, hex.EncodeToString(h[:])
}

func writeJSON(t *testing.T, path string, v any) {
	t.Helper()
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
}

// writeLocalIndex writes a single-package, single-version index.json plus
// its tarball to a temp "dist/packages"-shaped directory (relative URL,
// like `cn build` produces), and returns the index path.
func writeLocalIndex(t *testing.T, name, version, manifestJSON string) string {
	t.Helper()
	dir := t.TempDir()
	tarData, sha := makeTarGz(t, map[string]string{"cn.package.json": manifestJSON})
	tarName := name + "-" + version + ".tar.gz"
	if err := os.WriteFile(filepath.Join(dir, tarName), tarData, 0644); err != nil {
		t.Fatal(err)
	}
	idx := pkg.PackageIndex{
		Schema: "cn.package-index.v1",
		Packages: map[string]map[string]pkg.IndexEntry{
			name: {version: {URL: tarName, SHA256: sha}},
		},
	}
	indexPath := filepath.Join(dir, "index.json")
	writeJSON(t, indexPath, idx)
	return indexPath
}

func noopStdio() (*bytes.Buffer, *bytes.Buffer) {
	return &bytes.Buffer{}, &bytes.Buffer{}
}

// --- validateDispatch ---

// cnos#610: "cds" is no longer unconditionally refused — validateDispatch
// only rejects unrecognized --dispatch values now. The old
// TestValidateDispatch asserted "cds" always errors mentioning #609; that
// assertion is exactly the behavior this cell removes (see
// TestRun_DispatchCds_* below for the new per-precondition contract).
func TestValidateDispatch(t *testing.T) {
	cases := []struct {
		in      string
		wantErr bool
	}{
		{"", false},
		{"none", false},
		{"cds", false},
		{"bogus", true},
	}
	for _, c := range cases {
		err := validateDispatch(c.in)
		if (err != nil) != c.wantErr {
			t.Errorf("validateDispatch(%q) err = %v, wantErr %v", c.in, err, c.wantErr)
		}
	}
}

// --- ParseArgs ---

func TestParseArgs(t *testing.T) {
	a, err := ParseArgs([]string{
		"--release", "3.82.0",
		"--index", "./dist/packages/index.json",
		"--packages", "cnos.core, cnos.cdd ,cnos.cds",
		"--dispatch", "none",
		"--dry-run",
	})
	if err != nil {
		t.Fatalf("ParseArgs: %v", err)
	}
	if a.Release != "3.82.0" {
		t.Errorf("Release = %q", a.Release)
	}
	if a.IndexPath != "./dist/packages/index.json" {
		t.Errorf("IndexPath = %q", a.IndexPath)
	}
	want := []string{"cnos.core", "cnos.cdd", "cnos.cds"}
	if len(a.Packages) != len(want) {
		t.Fatalf("Packages = %v, want %v", a.Packages, want)
	}
	for i, p := range want {
		if a.Packages[i] != p {
			t.Errorf("Packages[%d] = %q, want %q", i, a.Packages[i], p)
		}
	}
	if a.Dispatch != "none" {
		t.Errorf("Dispatch = %q", a.Dispatch)
	}
	if !a.DryRun {
		t.Error("DryRun = false, want true")
	}
}

func TestParseArgs_UnknownFlag(t *testing.T) {
	_, err := ParseArgs([]string{"--bogus"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestParseArgs_MissingValue(t *testing.T) {
	_, err := ParseArgs([]string{"--release"})
	if err == nil {
		t.Fatal("expected error for missing --release value")
	}
}

// --- resolvePins ---

func TestResolvePins_PinnedVersionFound(t *testing.T) {
	idx := &pkg.PackageIndex{Packages: map[string]map[string]pkg.IndexEntry{
		"cnos.core": {"3.82.0": {URL: "u", SHA256: "s"}},
	}}
	deps, err := resolvePins(idx, []string{"cnos.core"}, "3.82.0")
	if err != nil {
		t.Fatalf("resolvePins: %v", err)
	}
	if len(deps) != 1 || deps[0].Name != "cnos.core" || deps[0].Version != "3.82.0" {
		t.Errorf("deps = %+v", deps)
	}
}

func TestResolvePins_PinnedVersionMissing(t *testing.T) {
	idx := &pkg.PackageIndex{Packages: map[string]map[string]pkg.IndexEntry{
		"cnos.core": {"1.0.0": {URL: "u", SHA256: "s"}},
	}}
	_, err := resolvePins(idx, []string{"cnos.core"}, "3.82.0")
	if err == nil {
		t.Fatal("expected error for missing pin")
	}
	if !strings.Contains(err.Error(), "cnos.core@3.82.0") {
		t.Errorf("error should name the missing pin, got: %v", err)
	}
}

func TestResolvePins_SingleVersionNoPreference(t *testing.T) {
	idx := &pkg.PackageIndex{Packages: map[string]map[string]pkg.IndexEntry{
		"cnos.core": {"9.9.9": {URL: "u", SHA256: "s"}},
	}}
	deps, err := resolvePins(idx, []string{"cnos.core"}, "")
	if err != nil {
		t.Fatalf("resolvePins: %v", err)
	}
	if len(deps) != 1 || deps[0].Version != "9.9.9" {
		t.Errorf("deps = %+v", deps)
	}
}

func TestResolvePins_MissingPackage(t *testing.T) {
	idx := &pkg.PackageIndex{Packages: map[string]map[string]pkg.IndexEntry{}}
	_, err := resolvePins(idx, []string{"cnos.core"}, "")
	if err == nil {
		t.Fatal("expected error for missing package")
	}
	if !strings.Contains(err.Error(), "cnos.core") {
		t.Errorf("error should name the missing package, got: %v", err)
	}
}

func TestResolvePins_AmbiguousVersionsWithoutRelease(t *testing.T) {
	idx := &pkg.PackageIndex{Packages: map[string]map[string]pkg.IndexEntry{
		"cnos.core": {"1.0.0": {}, "2.0.0": {}},
	}}
	_, err := resolvePins(idx, []string{"cnos.core"}, "")
	if err == nil {
		t.Fatal("expected error for ambiguous versions")
	}
	if !strings.Contains(err.Error(), "--release") {
		t.Errorf("error should suggest --release, got: %v", err)
	}
}

// --- rewriteRelativeEntriesFromBase / rewriteRelativeEntries ---

func TestRewriteRelativeEntriesFromBase(t *testing.T) {
	idx := &pkg.PackageIndex{Schema: "cn.package-index.v1", Packages: map[string]map[string]pkg.IndexEntry{
		"cnos.core": {
			"1.0.0": {URL: "cnos.core-1.0.0.tar.gz", SHA256: "s1"},
		},
		"cnos.cdd": {
			"1.0.0": {URL: "https://example.com/cnos.cdd-1.0.0.tar.gz", SHA256: "s2"},
		},
	}}
	out := rewriteRelativeEntriesFromBase(idx, "https://dl.example.com/owner/repo/releases/download/1.0.0/")

	got := out.Packages["cnos.core"]["1.0.0"].URL
	want := "https://dl.example.com/owner/repo/releases/download/1.0.0/cnos.core-1.0.0.tar.gz"
	if got != want {
		t.Errorf("relative URL rewritten = %q, want %q", got, want)
	}

	// HTTP-tarball-URL preservation: an already-absolute entry must not
	// be touched (double-prefixing would break the fetch).
	gotAbs := out.Packages["cnos.cdd"]["1.0.0"].URL
	if gotAbs != "https://example.com/cnos.cdd-1.0.0.tar.gz" {
		t.Errorf("absolute URL was mutated: %q", gotAbs)
	}

	// Original must be untouched (pure function).
	if idx.Packages["cnos.core"]["1.0.0"].URL != "cnos.core-1.0.0.tar.gz" {
		t.Error("rewriteRelativeEntriesFromBase mutated its input")
	}
}

func TestRewriteRelativeEntries_DerivesBaseFromIndexURL(t *testing.T) {
	idx := &pkg.PackageIndex{Packages: map[string]map[string]pkg.IndexEntry{
		"cnos.core": {"1.0.0": {URL: "cnos.core-1.0.0.tar.gz"}},
	}}
	out := rewriteRelativeEntries(idx, "https://dl.example.com/owner/repo/releases/download/9.9.9/index.json")
	got := out.Packages["cnos.core"]["1.0.0"].URL
	want := "https://dl.example.com/owner/repo/releases/download/9.9.9/cnos.core-1.0.0.tar.gz"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// --- End-to-end Run: local index path ---

func TestRun_LocalIndex_EndToEnd(t *testing.T) {
	indexPath := writeLocalIndex(t, "cnos.core", "9.9.9", `{"name": "cnos.core", "version": "9.9.9"}`)
	repoRoot := t.TempDir()

	stdout, stderr := noopStdio()
	res, err := Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: indexPath,
		Packages:  []string{"cnos.core"},
		Stdout:    stdout,
		Stderr:    stderr,
	})
	if err != nil {
		t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
	}
	if res.DryRun {
		t.Error("Result.DryRun = true, want false")
	}
	if res.ReleaseTag != "" {
		t.Errorf("ReleaseTag = %q, want empty (index-only flow)", res.ReleaseTag)
	}

	// AC2/AC3: .cn/deps.json + .cn/deps.lock.json exist and parse.
	depsPath := filepath.Join(repoRoot, ".cn", "deps.json")
	data, err := os.ReadFile(depsPath)
	if err != nil {
		t.Fatalf("read deps.json: %v", err)
	}
	m, err := pkg.ParseManifest(data)
	if err != nil {
		t.Fatalf("parse deps.json: %v", err)
	}
	if m.Schema != "cn.deps.v1" {
		t.Errorf("deps.json schema = %q, want cn.deps.v1", m.Schema)
	}
	if len(m.Packages) != 1 || m.Packages[0].Name != "cnos.core" || m.Packages[0].Version != "9.9.9" {
		t.Errorf("deps.json packages = %+v", m.Packages)
	}

	lockPath := filepath.Join(repoRoot, ".cn", "deps.lock.json")
	lockData, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatalf("read deps.lock.json: %v", err)
	}
	lf, err := pkg.ParseLockfile(lockData)
	if err != nil {
		t.Fatalf("parse deps.lock.json: %v", err)
	}
	if lf.Schema != "cn.lock.v2" {
		t.Errorf("lockfile schema = %q, want cn.lock.v2", lf.Schema)
	}
	if len(lf.Packages) != 1 || lf.Packages[0].SHA256 == "" {
		t.Errorf("lockfile packages = %+v, want a non-empty sha256", lf.Packages)
	}

	// AC4: package restored under the name-based vendor path.
	pkgManifest := filepath.Join(pkg.VendorPath(repoRoot, "cnos.core"), "cn.package.json")
	if _, err := os.Stat(pkgManifest); err != nil {
		t.Errorf("cn.package.json not restored: %v", err)
	}

	// AC10 regression guard: no agent-hub scaffold.
	for _, p := range []string{"spec/SOUL.md", "agent", "threads", "state"} {
		if _, err := os.Stat(filepath.Join(repoRoot, p)); !os.IsNotExist(err) {
			t.Errorf("agent-hub scaffold path %q must not exist (err=%v)", p, err)
		}
	}

	// AC8: no .github/workflows/ in base mode.
	if _, err := os.Stat(filepath.Join(repoRoot, ".github")); !os.IsNotExist(err) {
		t.Errorf(".github must not exist after base install (err=%v)", err)
	}

	// .gitignore contains the vendor entry.
	gi, err := os.ReadFile(filepath.Join(repoRoot, ".gitignore"))
	if err != nil {
		t.Fatalf("read .gitignore: %v", err)
	}
	if !strings.Contains(string(gi), ".cn/vendor/") {
		t.Errorf(".gitignore = %q, want .cn/vendor/ entry", gi)
	}
}

// AC4 (literal default set): absent --packages, Run installs exactly the
// default triple (cnos.core, cnos.cdd, cnos.cds) under
// .cn/vendor/packages/<name>/, each with a validated cn.package.json.
func TestRun_DefaultPackageSet_AllThreeRestored(t *testing.T) {
	dir := t.TempDir()
	idx := pkg.PackageIndex{Schema: "cn.package-index.v1", Packages: map[string]map[string]pkg.IndexEntry{}}
	for _, name := range DefaultPackages {
		tarData, sha := makeTarGz(t, map[string]string{
			"cn.package.json": `{"name": "` + name + `", "version": "3.82.0"}`,
		})
		tarName := name + "-3.82.0.tar.gz"
		os.WriteFile(filepath.Join(dir, tarName), tarData, 0644)
		idx.Packages[name] = map[string]pkg.IndexEntry{"3.82.0": {URL: tarName, SHA256: sha}}
	}
	indexPath := filepath.Join(dir, "index.json")
	writeJSON(t, indexPath, idx)

	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()
	res, err := Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		Release:   "3.82.0", // pins the shared version across all three, mirroring Mock A
		IndexPath: indexPath,
		Stdout:    stdout,
		Stderr:    stderr,
	})
	if err != nil {
		t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
	}
	if len(res.Manifest.Packages) != len(DefaultPackages) {
		t.Fatalf("manifest packages = %+v, want the default triple", res.Manifest.Packages)
	}
	for i, name := range DefaultPackages {
		if res.Manifest.Packages[i].Name != name {
			t.Errorf("manifest.Packages[%d].Name = %q, want %q (default order preserved)", i, res.Manifest.Packages[i].Name, name)
		}
		manifestPath := filepath.Join(pkg.VendorPath(repoRoot, name), "cn.package.json")
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			t.Errorf("%s: cn.package.json not restored: %v", name, err)
			continue
		}
		if err := pkg.ValidatePackageManifestData(data, name); err != nil {
			t.Errorf("%s: manifest validation failed: %v", name, err)
		}
	}
}

// AC2/AC5: .cn/deps.json and .cn/deps.lock.json are byte-identical across
// two successive runs with the same inputs (determinism + idempotence).
func TestRun_Idempotent_ByteIdenticalArtifacts(t *testing.T) {
	indexPath := writeLocalIndex(t, "cnos.core", "9.9.9", `{"name": "cnos.core", "version": "9.9.9"}`)
	repoRoot := t.TempDir()

	run := func() ([]byte, []byte) {
		stdout, stderr := noopStdio()
		_, err := Run(context.Background(), Options{
			RepoRoot:  repoRoot,
			IndexPath: indexPath,
			Packages:  []string{"cnos.core"},
			Stdout:    stdout,
			Stderr:    stderr,
		})
		if err != nil {
			t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
		}
		deps, err := os.ReadFile(filepath.Join(repoRoot, ".cn", "deps.json"))
		if err != nil {
			t.Fatal(err)
		}
		lock, err := os.ReadFile(filepath.Join(repoRoot, ".cn", "deps.lock.json"))
		if err != nil {
			t.Fatal(err)
		}
		return deps, lock
	}

	deps1, lock1 := run()
	deps2, lock2 := run()

	if !bytes.Equal(deps1, deps2) {
		t.Errorf("deps.json not byte-identical across runs:\n--- 1 ---\n%s\n--- 2 ---\n%s", deps1, deps2)
	}
	if !bytes.Equal(lock1, lock2) {
		t.Errorf("deps.lock.json not byte-identical across runs:\n--- 1 ---\n%s\n--- 2 ---\n%s", lock1, lock2)
	}
}

// --dry-run must write nothing at all, not even a .cn/ directory.
func TestRun_DryRun_WritesNothing(t *testing.T) {
	indexPath := writeLocalIndex(t, "cnos.core", "9.9.9", `{"name": "cnos.core", "version": "9.9.9"}`)
	repoRoot := t.TempDir()

	stdout, stderr := noopStdio()
	res, err := Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: indexPath,
		Packages:  []string{"cnos.core"},
		DryRun:    true,
		Stdout:    stdout,
		Stderr:    stderr,
	})
	if err != nil {
		t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
	}
	if !res.DryRun {
		t.Error("Result.DryRun = false, want true")
	}

	entries, err := os.ReadDir(repoRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("dry-run must write nothing under RepoRoot; found: %v", entries)
	}

	out := stdout.String()
	for _, want := range []string{".cn/deps.json", ".cn/deps.lock.json", ".gitignore"} {
		if !strings.Contains(out, want) {
			t.Errorf("dry-run stdout missing planned-diff entry %q:\n%s", want, out)
		}
	}
	if !strings.Contains(out, "Dispatch: none") {
		t.Errorf("dry-run stdout must state dispatch mode:\n%s", out)
	}
}

func TestRun_MissingIndexFile(t *testing.T) {
	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()
	_, err := Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: filepath.Join(t.TempDir(), "nonexistent-index.json"),
		Stdout:    stdout,
		Stderr:    stderr,
	})
	if err == nil {
		t.Fatal("expected error for missing index file")
	}
	if !strings.Contains(stderr.String(), "✗") {
		t.Errorf("stderr should carry a diagnostic, got: %q", stderr.String())
	}
}

func TestRun_MissingPackageInIndex(t *testing.T) {
	indexPath := writeLocalIndex(t, "cnos.core", "9.9.9", `{"name": "cnos.core", "version": "9.9.9"}`)
	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()
	_, err := Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: indexPath,
		Packages:  []string{"cnos.core", "cnos.cdd"}, // cnos.cdd absent from the fixture index
		Stdout:    stdout,
		Stderr:    stderr,
	})
	if err == nil {
		t.Fatal("expected error for package not present in index")
	}
	if !strings.Contains(err.Error(), "cnos.cdd") {
		t.Errorf("error should name the missing package, got: %v", err)
	}
	// No partial writes on failure.
	if _, statErr := os.Stat(filepath.Join(repoRoot, ".cn")); !os.IsNotExist(statErr) {
		t.Error(".cn must not exist after a failed resolve (pin resolution happens before any write)")
	}
}

// SHA-mismatch propagation: restore.Restore's existing verify step must
// not be silently bypassed by this installer.
func TestRun_SHAMismatchPropagates(t *testing.T) {
	dir := t.TempDir()
	tarData, _ := makeTarGz(t, map[string]string{"cn.package.json": `{"name": "cnos.core", "version": "9.9.9"}`})
	tarName := "cnos.core-9.9.9.tar.gz"
	os.WriteFile(filepath.Join(dir, tarName), tarData, 0644)

	idx := pkg.PackageIndex{
		Schema: "cn.package-index.v1",
		Packages: map[string]map[string]pkg.IndexEntry{
			"cnos.core": {"9.9.9": {URL: tarName, SHA256: "deliberately-wrong-sha256"}},
		},
	}
	indexPath := filepath.Join(dir, "index.json")
	writeJSON(t, indexPath, idx)

	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()
	_, err := Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: indexPath,
		Packages:  []string{"cnos.core"},
		Stdout:    stdout,
		Stderr:    stderr,
	})
	if err == nil {
		t.Fatal("expected error for sha256 mismatch")
	}
	if !strings.Contains(stderr.String(), "sha256 mismatch") {
		t.Errorf("stderr should surface the sha256 mismatch, got: %q", stderr.String())
	}
	if _, statErr := os.Stat(pkg.VendorPath(repoRoot, "cnos.core")); !os.IsNotExist(statErr) {
		t.Error("vendor dir must not exist after a sha256 mismatch")
	}
}

// cnos#610: this test previously asserted --dispatch cds always failed
// unconditionally, naming #609 (the pre-generalization guard). That
// guard is exactly what this cell removes — TestRun_DispatchCds_
// RendersWorkflow_ThenSurfacesLabelGap above proves cds now renders
// successfully once the renderer + identity preconditions are met. The
// genuine failure path this fixture still exercises: --packages here
// carries only "cnos.core" (a minimal fixture manifest with no actual
// command content, mirroring writeLocalIndex's other callers) — the
// dispatch renderer script itself is therefore never vendored, so
// runDispatchCds's own defensive existence check fires before any
// subprocess is even spawned, and — matching Mock C2 "no partial
// render" — no .github/workflows/ file or directory may exist.
func TestRun_DispatchCds_RendererNotVendored_FailsWithNoPartialWrite(t *testing.T) {
	indexPath := writeLocalIndex(t, "cnos.core", "9.9.9", `{"name": "cnos.core", "version": "9.9.9"}`)
	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()
	_, err := Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: indexPath,
		Packages:  []string{"cnos.core"},
		Dispatch:  "cds",
		Stdout:    stdout,
		Stderr:    stderr,
	})
	if err == nil {
		t.Fatal("expected --dispatch cds to fail when the renderer script is not vendored")
	}
	if !strings.Contains(err.Error(), "dispatch renderer not found") {
		t.Errorf("error should name the missing renderer, got: %v", err)
	}
	if !strings.Contains(stderr.String(), "dispatch renderer not found") {
		t.Errorf("stderr should carry the same diagnostic, got: %q", stderr.String())
	}
	if _, statErr := os.Stat(filepath.Join(repoRoot, ".github")); !os.IsNotExist(statErr) {
		t.Error("no .github/ directory may exist after a missing-renderer failure (no partial render)")
	}
	// Base install artifacts, by contrast, are unaffected — base install
	// (C1) always precedes the dispatch render attempt.
	if _, statErr := os.Stat(filepath.Join(repoRoot, ".cn", "deps.json")); statErr != nil {
		t.Errorf("expected base install artifact .cn/deps.json to still exist: %v", statErr)
	}
}

// --- Release-resolution flows (httptest servers stand in for GitHub) ---

// newFakeGitHub returns (apiServer, dlServer) where apiServer answers
// /repos/{repo}/releases/latest like the real GitHub API, and dlServer
// serves index.json + tarball(s) at the release-download URL shape:
// /{repo}/releases/download/{tag}/<file>.
func newFakeGitHub(t *testing.T, repo, tag string, indexBody []byte, tarballs map[string][]byte) (*httptest.Server, *httptest.Server) {
	t.Helper()

	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"tag_name": %q, "target_commitish": "main"}`, tag)
	}))
	t.Cleanup(api.Close)

	prefix := "/" + repo + "/releases/download/" + tag + "/"
	dl := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == prefix+"index.json" {
			w.Write(indexBody)
			return
		}
		for name, data := range tarballs {
			if r.URL.Path == prefix+name {
				w.Write(data)
				return
			}
		}
		http.NotFound(w, r)
	}))
	t.Cleanup(dl.Close)

	return api, dl
}

// AC7 (--release latest): resolves a concrete tag via the GitHub releases
// API, fetches that release's index.json, rewrites its relative tarball
// URL to the release-download URL, and installs successfully — proving
// the rewrite (not just the API call) actually happened, since a fetch
// against the un-rewritten relative URL would 404 against dlServer.
func TestRun_ReleaseLatest_ResolvesAndInstalls(t *testing.T) {
	tarData, tarSHA := makeTarGz(t, map[string]string{"cn.package.json": `{"name": "cnos.core", "version": "9.9.9"}`})
	idx := pkg.PackageIndex{
		Schema: "cn.package-index.v1",
		Packages: map[string]map[string]pkg.IndexEntry{
			// Relative URL, exactly as `cn build`'s local index shape would
			// declare it — proves rewriteRelativeEntriesFromBase actually
			// runs on the release-fetched index, not only on --index inputs.
			"cnos.core": {"9.9.9": {URL: "cnos.core-9.9.9.tar.gz", SHA256: tarSHA}},
		},
	}
	indexBody, err := json.Marshal(idx)
	if err != nil {
		t.Fatal(err)
	}

	repo := "acme/widgets"
	api, dl := newFakeGitHub(t, repo, "9.9.9", indexBody, map[string][]byte{"cnos.core-9.9.9.tar.gz": tarData})

	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()
	res, err := Run(context.Background(), Options{
		RepoRoot:     repoRoot,
		Release:      "latest",
		Packages:     []string{"cnos.core"},
		Repo:         repo,
		APIBaseURL:   api.URL,
		DownloadBase: dl.URL,
		Stdout:       stdout,
		Stderr:       stderr,
	})
	if err != nil {
		t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
	}
	if res.ReleaseTag != "9.9.9" {
		t.Errorf("ReleaseTag = %q, want 9.9.9", res.ReleaseTag)
	}
	if res.Manifest.Packages[0].Version != "9.9.9" {
		t.Errorf("manifest pinned version = %q, want 9.9.9", res.Manifest.Packages[0].Version)
	}
	pkgManifest := filepath.Join(pkg.VendorPath(repoRoot, "cnos.core"), "cn.package.json")
	if _, err := os.Stat(pkgManifest); err != nil {
		t.Errorf("cnos.core not restored via release flow: %v", err)
	}
	if !strings.Contains(stdout.String(), "9.9.9") {
		t.Errorf("stdout should print the resolved release tag:\n%s", stdout.String())
	}
}

// AC7 (--release <tag>): a pinned tag skips "latest" resolution entirely
// and fetches that tag's release assets directly.
func TestRun_ReleaseTag_Pinned(t *testing.T) {
	tarData, tarSHA := makeTarGz(t, map[string]string{"cn.package.json": `{"name": "cnos.core", "version": "1.2.3"}`})
	idx := pkg.PackageIndex{
		Schema: "cn.package-index.v1",
		Packages: map[string]map[string]pkg.IndexEntry{
			"cnos.core": {"1.2.3": {URL: "cnos.core-1.2.3.tar.gz", SHA256: tarSHA}},
		},
	}
	indexBody, err := json.Marshal(idx)
	if err != nil {
		t.Fatal(err)
	}

	repo := "acme/widgets"
	// No "latest" endpoint registered on a distinct path — if Run called
	// FetchLatestRelease despite an explicit tag, the download server
	// would 404 on a tag it never served, since we only wire index.json
	// under the "1.2.3" download prefix below.
	_, dl := newFakeGitHub(t, repo, "1.2.3", indexBody, map[string][]byte{"cnos.core-1.2.3.tar.gz": tarData})

	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()
	res, err := Run(context.Background(), Options{
		RepoRoot:     repoRoot,
		Release:      "1.2.3",
		Packages:     []string{"cnos.core"},
		Repo:         repo,
		DownloadBase: dl.URL,
		Stdout:       stdout,
		Stderr:       stderr,
	})
	if err != nil {
		t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
	}
	if res.ReleaseTag != "1.2.3" {
		t.Errorf("ReleaseTag = %q, want 1.2.3", res.ReleaseTag)
	}
}

// AC7 (--index <URL>): an explicit http(s) index URL is fetched, its
// relative tarball URL is rewritten against the URL's own directory, and
// install succeeds.
func TestRun_IndexHTTPURL_RelativeRewrite(t *testing.T) {
	tarData, tarSHA := makeTarGz(t, map[string]string{"cn.package.json": `{"name": "cnos.core", "version": "5.0.0"}`})
	idx := pkg.PackageIndex{
		Schema: "cn.package-index.v1",
		Packages: map[string]map[string]pkg.IndexEntry{
			"cnos.core": {"5.0.0": {URL: "cnos.core-5.0.0.tar.gz", SHA256: tarSHA}},
		},
	}
	indexBody, err := json.Marshal(idx)
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/dist/packages/index.json":
			w.Write(indexBody)
		case "/dist/packages/cnos.core-5.0.0.tar.gz":
			w.Write(tarData)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()
	_, err = Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: srv.URL + "/dist/packages/index.json",
		Packages:  []string{"cnos.core"},
		Stdout:    stdout,
		Stderr:    stderr,
	})
	if err != nil {
		t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
	}
	pkgManifest := filepath.Join(pkg.VendorPath(repoRoot, "cnos.core"), "cn.package.json")
	if _, err := os.Stat(pkgManifest); err != nil {
		t.Errorf("cnos.core not restored via HTTP index flow: %v", err)
	}
}

// AC7 combination: --index <URL> together with an explicit --release
// <tag> (not "latest"). TestRun_IndexHTTPURL_RelativeRewrite and
// TestRun_ReleaseTag_Pinned each exercise one of these independently;
// this test exercises resolveIndex's isRemoteURL(indexArg) branch with a
// non-empty pinFromRelease at the same time — the untested flag
// combination β's R0 review flagged as narrow test-coverage debt.
func TestRun_IndexHTTPURL_WithExplicitReleaseTag(t *testing.T) {
	tarData, tarSHA := makeTarGz(t, map[string]string{"cn.package.json": `{"name": "cnos.core", "version": "7.7.7"}`})
	idx := pkg.PackageIndex{
		Schema: "cn.package-index.v1",
		Packages: map[string]map[string]pkg.IndexEntry{
			"cnos.core": {"7.7.7": {URL: "cnos.core-7.7.7.tar.gz", SHA256: tarSHA}},
		},
	}
	indexBody, err := json.Marshal(idx)
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/dist/packages/index.json":
			w.Write(indexBody)
		case "/dist/packages/cnos.core-7.7.7.tar.gz":
			w.Write(tarData)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()
	res, err := Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: srv.URL + "/dist/packages/index.json",
		Release:   "7.7.7",
		Packages:  []string{"cnos.core"},
		Stdout:    stdout,
		Stderr:    stderr,
	})
	if err != nil {
		t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
	}
	if res.ReleaseTag != "7.7.7" {
		t.Errorf("ReleaseTag = %q, want 7.7.7 (explicit --release must still pin the version even when --index is also given)", res.ReleaseTag)
	}
	pkgManifest := filepath.Join(pkg.VendorPath(repoRoot, "cnos.core"), "cn.package.json")
	if _, err := os.Stat(pkgManifest); err != nil {
		t.Errorf("cnos.core not restored via --index URL + --release combo: %v", err)
	}
}

// --- Dispatch install (cnos#610: cn repo install --dispatch cds) ---

// repoSrcRoot resolves the actual repository root (the directory
// containing src/packages/) from this test file's own path, independent
// of process cwd. The dispatch tests below vendor the REAL cn-install-wake
// renderer script and the REAL (post-cnos#610) cds-dispatch/SKILL.md
// prose into fixture package tarballs, so they exercise the actual
// renderer + actual prose this cycle shipped, not a synthetic stand-in.
func repoSrcRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	// thisFile: <repoRoot>/src/go/internal/repoinstall/repoinstall_test.go
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..")
}

func readRealFile(t *testing.T, repoRoot string, relPath ...string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(append([]string{repoRoot}, relPath...)...))
	if err != nil {
		t.Fatalf("read %s: %v", filepath.Join(relPath...), err)
	}
	return data
}

// tarEntry is a single file to write into a fixture tarball, with an
// explicit mode — unlike makeTarGz above (always 0644), the renderer
// script fixture below needs its real executable bit preserved so
// os/exec can run it once restore.Restore extracts it.
type tarEntry struct {
	name string
	mode int64
	data []byte
}

func makeTarGzEntries(t *testing.T, entries []tarEntry) ([]byte, string) {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	for _, e := range entries {
		hdr := &tar.Header{Name: e.name, Mode: e.mode, Size: int64(len(e.data))}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(e.data); err != nil {
			t.Fatal(err)
		}
	}
	tw.Close()
	gw.Close()
	data := buf.Bytes()
	h := sha256.Sum256(data)
	return data, hex.EncodeToString(h[:])
}

// writeDispatchFixtureIndex builds a package index + two tarballs:
// cnos.core (carrying the REAL commands/install-wake/cn-install-wake
// renderer script, executable bit preserved) and cnos.cds (carrying the
// REAL orchestrators/cds-dispatch/SKILL.md this cycle edited). Returns
// the index path.
func writeDispatchFixtureIndex(t *testing.T) string {
	t.Helper()
	root := repoSrcRoot(t)
	rendererScript := readRealFile(t, root, "src", "packages", "cnos.core", "commands", "install-wake", "cn-install-wake")
	skillMD := readRealFile(t, root, "src", "packages", "cnos.cds", "orchestrators", "cds-dispatch", "SKILL.md")

	dir := t.TempDir()

	coreTar, coreSHA := makeTarGzEntries(t, []tarEntry{
		{name: "cn.package.json", mode: 0644, data: []byte(`{"name": "cnos.core", "version": "9.9.9"}`)},
		{name: "commands/install-wake/cn-install-wake", mode: 0755, data: rendererScript},
	})
	coreTarName := "cnos.core-9.9.9.tar.gz"
	if err := os.WriteFile(filepath.Join(dir, coreTarName), coreTar, 0644); err != nil {
		t.Fatal(err)
	}

	cdsTar, cdsSHA := makeTarGzEntries(t, []tarEntry{
		{name: "cn.package.json", mode: 0644, data: []byte(`{"name": "cnos.cds", "version": "9.9.9"}`)},
		{name: "orchestrators/cds-dispatch/SKILL.md", mode: 0644, data: skillMD},
	})
	cdsTarName := "cnos.cds-9.9.9.tar.gz"
	if err := os.WriteFile(filepath.Join(dir, cdsTarName), cdsTar, 0644); err != nil {
		t.Fatal(err)
	}

	idx := pkg.PackageIndex{
		Schema: "cn.package-index.v1",
		Packages: map[string]map[string]pkg.IndexEntry{
			"cnos.core": {"9.9.9": {URL: coreTarName, SHA256: coreSHA}},
			"cnos.cds":  {"9.9.9": {URL: cdsTarName, SHA256: cdsSHA}},
		},
	}
	indexPath := filepath.Join(dir, "index.json")
	writeJSON(t, indexPath, idx)
	return indexPath
}

// AC1/C1/C3 + AC4/C4/C6 + AC5-positive: a non-sigma agent with full
// identity flags drives the real renderer end-to-end through Run. AC1's
// oracle ("writes exactly .github/workflows/cnos-cds-dispatch.yml after
// ... base artifacts present") and the cnos#493 label-doctor
// integration's oracle ("the dispatch-install path returns a named,
// non-silent error when the label mechanism cannot resolve its target")
// describe the SAME observed behavior here: the base install and the
// render both complete (this is what AC1 checks), and only AFTER that
// Run surfaces label-doctor's target-repo-resolution failure as its
// returned error (repoRoot here is a plain t.TempDir() with no git
// remote configured, so label-doctor cannot resolve "owner/repo" to
// call the GitHub API — this is expected, not a live network call) —
// not a silent skip, and not a partial/rolled-back render. See
// .cdd/unreleased/610/self-coherence.md §ACs for why these two ACs were
// originally verified by one test rather than two independent
// scenarios (cnos#493 replaced the label-install stub this test
// originally pinned; see .cdd/unreleased/493/self-coherence.md).
func TestRun_DispatchCds_RendersWorkflow_ThenSurfacesLabelGap(t *testing.T) {
	indexPath := writeDispatchFixtureIndex(t)
	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()

	res, err := Run(context.Background(), Options{
		RepoRoot:          repoRoot,
		IndexPath:         indexPath,
		Packages:          []string{"cnos.core", "cnos.cds"},
		Dispatch:          "cds",
		Agent:             "acme",
		WorkflowPatSecret: "ACME_WORKFLOW_PAT",
		BotName:           "acme-bot",
		BotID:             "12345678",
		Stdout:            stdout,
		Stderr:            stderr,
	})

	// AC1/C1: base install artifacts exist regardless of the AC3 gap.
	for _, f := range []string{filepath.Join(".cn", "deps.json"), filepath.Join(".cn", "deps.lock.json")} {
		if _, statErr := os.Stat(filepath.Join(repoRoot, f)); statErr != nil {
			t.Errorf("expected base install artifact %s: %v", f, statErr)
		}
	}

	// AC1/C3: the workflow renders to exactly this path.
	workflowPath := filepath.Join(repoRoot, ".github", "workflows", "cnos-cds-dispatch.yml")
	data, statErr := os.ReadFile(workflowPath)
	if statErr != nil {
		t.Fatalf("expected rendered workflow at %s: %v\nstdout: %s\nstderr: %s", workflowPath, statErr, stdout, stderr)
	}

	// cnos#493: repoRoot has no git remote configured, so label-doctor
	// cannot resolve the installing repo's "owner/repo" target — Run
	// surfaces this as a named, actionable error (not a silent skip),
	// even though the render itself (checked above) succeeded.
	if err == nil {
		t.Fatal("expected Run to surface the label-doctor target-resolution error")
	}
	if !strings.Contains(err.Error(), "canonical dispatch labels not ensured") {
		t.Errorf("error should name the canonical-dispatch-labels failure, got: %v", err)
	}
	if !strings.Contains(err.Error(), "could not resolve target repo") {
		t.Errorf("error should explain WHY (no git remote), got: %v", err)
	}
	if strings.Contains(err.Error(), "cnos#493") {
		t.Errorf("error should no longer name cnos#493 as an unshipped-mechanism stub (the mechanism now exists), got: %v", err)
	}
	if res != nil {
		t.Errorf("Run should return a nil Result on this error path (matching every other error path in this file), got %+v", res)
	}
	if !strings.Contains(stderr.String(), "canonical dispatch labels not ensured") {
		t.Errorf("stderr should also carry the diagnostic, got: %q", stderr.String())
	}

	// AC4/C4: zero sigma substrate-binding leak in the acme render.
	content := string(data)
	for _, leak := range []string{"SIGMA_WORKFLOW_PAT", "41898282", "cds-dispatch-sigma", "sigma@cnos.cn-sigma.cnos"} {
		if strings.Contains(content, leak) {
			t.Errorf("acme render leaks sigma-bound substrate token %q", leak)
		}
	}

	// AC5 positive: no tenant-visible/confusing sigma prose leak.
	for _, leak := range []string{"today: `sigma`", "agent-admin-sigma", "cn-sigma:"} {
		if strings.Contains(content, leak) {
			t.Errorf("acme render leaks tenant-visible sigma prose: %q", leak)
		}
	}
	if !strings.Contains(content, "cds-dispatch-acme") {
		t.Error("acme render should bind its own concurrency group (cds-dispatch-acme)")
	}

	// AC4/C6: workflow-scope PAT requirement + never-pushes-main fact
	// stated in stdout.
	out := stdout.String()
	if !strings.Contains(out, "needs `workflow` scope") {
		t.Errorf("stdout should state the workflow-scope PAT requirement:\n%s", out)
	}
	if !strings.Contains(out, "never pushes to main") {
		t.Errorf("stdout should state the never-pushes-main fact:\n%s", out)
	}
}

// AC5 positive (Go-level identity resolution): a bare --dispatch cds
// with no --agent/--workflow-pat-secret must NOT trip the identity gate
// (agent resolves to "sigma", which has a default substrate PAT
// binding) — preserving the implicit sigma-bound behavior AC5 requires
// as backward-compat. The render still ends in the same label-doctor
// target-resolution error as the acme case above (cnos#493 — repoRoot
// has no git remote); what this test isolates is that the identity gate
// itself does not fire for the sigma default.
func TestRun_DispatchCds_SigmaDefault_NoIdentityFlagsRequired(t *testing.T) {
	indexPath := writeDispatchFixtureIndex(t)
	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()

	_, err := Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: indexPath,
		Packages:  []string{"cnos.core", "cnos.cds"},
		Dispatch:  "cds",
		Stdout:    stdout,
		Stderr:    stderr,
	})

	workflowPath := filepath.Join(repoRoot, ".github", "workflows", "cnos-cds-dispatch.yml")
	data, statErr := os.ReadFile(workflowPath)
	if statErr != nil {
		t.Fatalf("expected rendered workflow at %s (sigma default should render, not fail identity): %v\nstdout: %s\nstderr: %s", workflowPath, statErr, stdout, stderr)
	}
	if !strings.Contains(string(data), "cds-dispatch-sigma") {
		t.Error("default (no --agent) render should bind the sigma concurrency group")
	}
	// The only error expected here is label-doctor's target-resolution
	// failure (cnos#493 — no git remote on repoRoot), not an identity
	// error.
	if err == nil || !strings.Contains(err.Error(), "canonical dispatch labels not ensured") {
		t.Errorf("expected the canonical-dispatch-labels error (not an identity error), got: %v", err)
	}
	if strings.Contains(stderr.String(), "--workflow-pat-secret is required") {
		t.Errorf("sigma default must not trip the identity gate, stderr: %q", stderr.String())
	}
}

// AC2/C2 (Mock C2 "no partial render"): a non-sigma --agent with no
// --workflow-pat-secret must fail early, before the renderer ever runs
// — nonzero exit, no .github/workflows/ directory created at all.
func TestRun_DispatchCds_MissingIdentity_FailsEarlyNoPartialWrite(t *testing.T) {
	indexPath := writeDispatchFixtureIndex(t)
	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()

	_, err := Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: indexPath,
		Packages:  []string{"cnos.core", "cnos.cds"},
		Dispatch:  "cds",
		Agent:     "acme", // non-sigma, no --workflow-pat-secret
		Stdout:    stdout,
		Stderr:    stderr,
	})

	if err == nil {
		t.Fatal("expected --dispatch cds with a non-sigma agent and no --workflow-pat-secret to fail")
	}
	if !strings.Contains(err.Error(), "--workflow-pat-secret") {
		t.Errorf("error should name the missing flag, got: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(repoRoot, ".github")); !os.IsNotExist(statErr) {
		t.Error("no .github/ directory may exist after a missing-identity failure (no partial render)")
	}

	// Base install artifacts, by contrast, are unaffected by this
	// dispatch-specific gate (C1 layering: base install already ran).
	for _, f := range []string{filepath.Join(".cn", "deps.json"), filepath.Join(".cn", "deps.lock.json")} {
		if _, statErr := os.Stat(filepath.Join(repoRoot, f)); statErr != nil {
			t.Errorf("expected base install artifact %s to still exist: %v", f, statErr)
		}
	}
}

// AC5 negative case (non-vacuous oracle): the leak grep must actually
// be capable of catching a violation, not just pass on the (now-fixed)
// real prose by construction. This renders a synthetic dispatch-shaped
// SKILL.md fixture carrying the EXACT pre-cnos#610 hardcoded-sigma
// phrasing (the "(today: `sigma`; ...)" parenthetical + the
// "agent-admin-sigma" / "cn-sigma:" citation this cycle removed from
// the real cds-dispatch/SKILL.md) for a non-sigma agent, and asserts
// the leak grep DOES find hits — proving the grep oracle used by the
// positive-case test above is a real detector, not one that would pass
// on anything.
func TestDispatchRenderer_ProseLeakGrep_CatchesPreFixSigmaPhrasing(t *testing.T) {
	root := repoSrcRoot(t)
	rendererPath := filepath.Join(root, "src", "packages", "cnos.core", "commands", "install-wake", "cn-install-wake")
	if _, err := os.Stat(rendererPath); err != nil {
		t.Fatalf("renderer script not found at %s: %v", rendererPath, err)
	}

	fixtureDir := t.TempDir()
	skillPath := filepath.Join(fixtureDir, "SKILL.md")
	preFixBody := "---\n" +
		"name: test-pre-fix-leak\n" +
		"description: \"AC5 negative-case fixture — pre-cnos#610 hardcoded-sigma prose\"\n" +
		"governing_question: n/a\n" +
		"artifact_class: wake\n" +
		"scope: global\n" +
		"kata_surface: none\n" +
		"triggers:\n  - dispatch-wake\n" +
		"inputs:\n  - \"n/a\"\n" +
		"outputs:\n  - \"n/a\"\n" +
		"wake:\n" +
		"  role: dispatch\n" +
		"  package: cnos.cds\n" +
		"  admin_only: false\n" +
		"  activation_log_writer: false\n" +
		"  input:\n    triggers:\n      - issues_labeled_selector_match\n" +
		"  output:\n" +
		"    cycle_artifact_root: \".cdd/unreleased\"\n" +
		"    artifact_class_taxonomy:\n      - gamma-scaffold\n" +
		"    cell_runtime: shell\n" +
		"  permission_intent:\n    - contents.write\n" +
		"  concurrency:\n    serialize: true\n    group: \"test-pre-fix-leak-{agent}\"\n" +
		"  agent_variable:\n    name: agent\n    default: null\n" +
		"  surfaces:\n    allowed:\n      - \"n/a\"\n    disallowed:\n      - \"n/a\"\n" +
		"  protocol: \"test\"\n" +
		"  selector:\n    include:\n      - \"protocol:test\"\n    exclude: []\n" +
		"---\n\n" +
		"# AC5 negative-case fixture body (pre-cnos#610 prose)\n\n" +
		"You substrate-execute as `{agent}` (today: `sigma`; future: per-package bot accounts per cnos#449 follow-up).\n\n" +
		"The admin wake's `agent-admin-sigma` concurrency group runs in parallel. Empirical motivator: the 2026-06-24 mixed log entries (`cn-sigma:.cn-sigma/logs/20260624.md`, four selector-scan no-ops).\n"
	if err := os.WriteFile(skillPath, []byte(preFixBody), 0644); err != nil {
		t.Fatal(err)
	}

	outPath := filepath.Join(t.TempDir(), "pre-fix-render.yml")
	cmd := exec.Command(rendererPath, "test-pre-fix-leak",
		"--manifest", skillPath,
		"--agent", "acme", "--workflow-pat-secret", "ACME_WORKFLOW_PAT",
		"--bot-name", "acme-bot", "--bot-id", "12345678",
		"--out", outPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("render pre-fix fixture: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read pre-fix render: %v", err)
	}
	content := string(data)

	var found []string
	for _, leak := range []string{"today: `sigma`", "agent-admin-sigma", "cn-sigma:"} {
		if strings.Contains(content, leak) {
			found = append(found, leak)
		}
	}
	if len(found) != 3 {
		t.Fatalf("expected the pre-fix fixture to trip ALL three leak strings for a non-sigma agent (proving the grep oracle is real, not vacuous); found only %v in:\n%s", found, content)
	}
}

// extractSparseCheckoutPatterns pulls the literal-block-scalar lines under
// a "sparse-checkout: |" key out of rendered workflow YAML, stripping the
// fixed 12-space indent the renderer emits. Returns nil if the key is
// absent (the admin-shape negative case).
func extractSparseCheckoutPatterns(t *testing.T, yaml string) []string {
	t.Helper()
	const marker = "sparse-checkout: |\n"
	idx := strings.Index(yaml, marker)
	if idx < 0 {
		return nil
	}
	rest := yaml[idx+len(marker):]
	var patterns []string
	for _, line := range strings.Split(rest, "\n") {
		if !strings.HasPrefix(line, "            ") {
			break
		}
		patterns = append(patterns, strings.TrimPrefix(line, "            "))
	}
	return patterns
}

// TestDispatchRenderer_SparseCheckoutExcludesAgentHub is the cnos#626 AC3
// proof: the dispatch-shape checkout step must structurally exclude
// .cn-{agent}/ (the capability must be GONE, not merely instructed
// against — the cell's own governing doctrine at
// cnos.cdd/skills/cdd/delta/SKILL.md SS9.12). Two parts: (a) the rendered
// YAML carries the sparse-checkout block for role:dispatch and omits it
// for role:admin (the admin wake keeps its full checkout — it owns
// .cn-{agent}/logs/ per AGENT-ACTIVATION-LOG-v0 SS0); (b) the extracted
// patterns are fed to a REAL `git sparse-checkout set --no-cone` against a
// throwaway fixture repo containing a .cn-sigma/ tree, proving the
// resulting working tree omits it while unrelated paths survive. Part (b)
// is the part that would fail if the YAML said the right words but the
// underlying git mechanism didn't actually behave as claimed.
func TestDispatchRenderer_SparseCheckoutExcludesAgentHub(t *testing.T) {
	root := repoSrcRoot(t)
	rendererPath := filepath.Join(root, "src", "packages", "cnos.core", "commands", "install-wake", "cn-install-wake")

	renderManifest := func(manifestRelPath, wakeName string) string {
		t.Helper()
		manifestPath := filepath.Join(root, filepath.FromSlash(manifestRelPath))
		outPath := filepath.Join(t.TempDir(), wakeName+".yml")
		cmd := exec.Command(rendererPath, wakeName, "--manifest", manifestPath, "--out", outPath)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			t.Fatalf("render %s: %v\nstdout: %s\nstderr: %s", wakeName, err, stdout.String(), stderr.String())
		}
		data, err := os.ReadFile(outPath)
		if err != nil {
			t.Fatalf("read %s render: %v", wakeName, err)
		}
		return string(data)
	}

	dispatchYAML := renderManifest("src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md", "cds-dispatch")
	adminYAML := renderManifest("src/packages/cnos.core/orchestrators/agent-admin/SKILL.md", "agent-admin")

	// (a) presence/absence per role.
	dispatchPatterns := extractSparseCheckoutPatterns(t, dispatchYAML)
	wantPatterns := []string{"/*", "!/.cn-sigma"}
	if len(dispatchPatterns) != len(wantPatterns) {
		t.Fatalf("dispatch render: expected sparse-checkout patterns %v, got %v\nfull YAML:\n%s", wantPatterns, dispatchPatterns, dispatchYAML)
	}
	for i, want := range wantPatterns {
		if dispatchPatterns[i] != want {
			t.Errorf("dispatch render: pattern[%d] = %q, want %q", i, dispatchPatterns[i], want)
		}
	}
	if !strings.Contains(dispatchYAML, "sparse-checkout-cone-mode: false") {
		t.Errorf("dispatch render: expected sparse-checkout-cone-mode: false alongside the pattern block")
	}
	if adminPatterns := extractSparseCheckoutPatterns(t, adminYAML); adminPatterns != nil {
		t.Errorf("admin render: expected NO sparse-checkout block (admin keeps its full checkout), got patterns %v", adminPatterns)
	}

	// (b) the extracted patterns actually work against a real git checkout.
	fixtureRepo := t.TempDir()
	runGit := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = fixtureRepo
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out.String())
		}
	}
	runGit("init", "--quiet", "-b", "main")
	runGit("config", "user.email", "test@example.com")
	runGit("config", "user.name", "test")
	mustWrite := func(rel, content string) {
		t.Helper()
		full := filepath.Join(fixtureRepo, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	mustWrite(".cn-sigma/logs/20260101.md", "channel log entry\n")
	mustWrite(".cn-sigma/spec/PERSONA.md", "persona\n")
	mustWrite("src/go/go.mod", "module example\n")
	mustWrite("README.md", "hello\n")
	runGit("add", "-A")
	runGit("commit", "--quiet", "-m", "fixture")

	setArgs := append([]string{"-C", fixtureRepo, "sparse-checkout", "set", "--no-cone"}, dispatchPatterns...)
	runGit(setArgs...)

	if _, err := os.Stat(filepath.Join(fixtureRepo, ".cn-sigma")); !os.IsNotExist(err) {
		t.Errorf(".cn-sigma: expected absent after sparse-checkout, stat err = %v", err)
	}
	for _, keep := range []string{"src/go/go.mod", "README.md"} {
		if _, err := os.Stat(filepath.Join(fixtureRepo, keep)); err != nil {
			t.Errorf("%s: expected present after sparse-checkout, got %v", keep, err)
		}
	}
}

// --- cnos#613: PAT-free mechanical FSM-engine wake tier (Mock G) ---

// TestParseArgs_Engine covers the --engine boolean flag: it parses to
// Args.Engine and defaults to false when absent.
func TestParseArgs_Engine(t *testing.T) {
	a, err := ParseArgs([]string{"--dispatch", "cds", "--engine"})
	if err != nil {
		t.Fatalf("ParseArgs: %v", err)
	}
	if !a.Engine {
		t.Error("Engine = false, want true after --engine")
	}
	if a.Dispatch != "cds" {
		t.Errorf("Dispatch = %q, want cds", a.Dispatch)
	}

	b, err := ParseArgs([]string{"--dispatch", "cds"})
	if err != nil {
		t.Fatalf("ParseArgs: %v", err)
	}
	if b.Engine {
		t.Error("Engine = true, want false when --engine absent")
	}
}

// TestValidateEngine pins the cross-field guard: --engine is only valid
// with --dispatch cds. It is meaningless on a base (or non-cds) install
// and must be refused with a clear, named error rather than silently
// ignored.
func TestValidateEngine(t *testing.T) {
	cases := []struct {
		dispatch string
		engine   bool
		wantErr  bool
	}{
		{"cds", true, false},  // the one valid combination
		{"cds", false, false}, // engine off is always fine
		{"", false, false},
		{"none", false, false},
		{"", true, true},     // --engine with base install
		{"none", true, true}, // --engine with --dispatch none
	}
	for _, c := range cases {
		err := validateEngine(c.dispatch, c.engine)
		if (err != nil) != c.wantErr {
			t.Errorf("validateEngine(%q, %v) err = %v, wantErr %v", c.dispatch, c.engine, err, c.wantErr)
		}
	}
}

// TestRun_Engine_WithoutDispatchCds_Rejected proves the guard fires
// through the real Run entry point (not just the unit above): --engine on
// a base install (--dispatch none) fails early with a named error and
// writes no .github/workflows/ at all.
func TestRun_Engine_WithoutDispatchCds_Rejected(t *testing.T) {
	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()

	_, err := Run(context.Background(), Options{
		RepoRoot: repoRoot,
		Dispatch: "none",
		Engine:   true,
		Stdout:   stdout,
		Stderr:   stderr,
	})
	if err == nil {
		t.Fatal("expected --engine with --dispatch none to be rejected")
	}
	if !strings.Contains(err.Error(), "--engine is only valid with --dispatch cds") {
		t.Errorf("error should name the guard, got: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(repoRoot, ".github")); !os.IsNotExist(statErr) {
		t.Error("no .github/ directory may exist after an --engine misuse rejection")
	}
}

// TestRun_DispatchCds_EngineTier_RendersPatFreeMechanicalWake is the G1/G2
// render proof: --dispatch cds --engine drives the REAL vendored renderer +
// REAL cds-dispatch SKILL.md end-to-end through Run (mirroring the agent-tier
// TestRun_DispatchCds_RendersWorkflow_ThenSurfacesLabelGap fixture), and the
// emitted .github/workflows/cnos-cds-dispatch.yml (a) runs the mechanical FSM
// engine on the default GITHUB_TOKEN, (b) carries NONE of the agent-tier
// PAT/OAuth/bot substrate bindings. As with the agent-tier test, Run still
// ends by surfacing label-doctor's target-resolution error (repoRoot is a
// plain t.TempDir() with no git remote) — the render itself (asserted here)
// completed first.
func TestRun_DispatchCds_EngineTier_RendersPatFreeMechanicalWake(t *testing.T) {
	indexPath := writeDispatchFixtureIndex(t)
	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()

	_, err := Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: indexPath,
		Packages:  []string{"cnos.core", "cnos.cds"},
		Dispatch:  "cds",
		Engine:    true,
		Stdout:    stdout,
		Stderr:    stderr,
	})

	workflowPath := filepath.Join(repoRoot, ".github", "workflows", "cnos-cds-dispatch.yml")
	data, statErr := os.ReadFile(workflowPath)
	if statErr != nil {
		t.Fatalf("expected rendered engine-tier workflow at %s: %v\nstdout: %s\nstderr: %s", workflowPath, statErr, stdout, stderr)
	}
	content := string(data)

	// G1: the mechanical FSM engine runs on the default GITHUB_TOKEN.
	for _, want := range []string{
		"cn issues fsm scan --protocol cds --apply",
		"cn issues fsm evaluate",
		"--apply",
		"${{ secrets.GITHUB_TOKEN }}",
		"Mechanical FSM engine (PAT-free)",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("engine render should contain %q\nfull YAML:\n%s", want, content)
		}
	}

	// G2: the agent tier's PAT / OAuth / bot substrate bindings are the
	// SOLE domain of the agent tier — none of them may appear in an engine
	// render. (claude_code_oauth_token is the lower-cased action-input form;
	// asserting both it and the secret name closes the case fully.)
	for _, leak := range []string{
		"claude-code-action",
		"CLAUDE_CODE_OAUTH_TOKEN",
		"claude_code_oauth_token",
		"SIGMA_WORKFLOW_PAT",
		"41898282",
		"bot_name",
		"bot_id",
	} {
		if strings.Contains(content, leak) {
			t.Errorf("engine render must NOT contain agent-tier binding %q\nfull YAML:\n%s", leak, content)
		}
	}

	// cnos#613 (operator review): the engine tier must name an EXPLICIT
	// least-privilege permissions block, not inherit the agent tier's broad
	// permission_intent. Enough to mutate issue-label state (issues:write),
	// check out (contents:read), and read PRs (pull-requests:read) — and NO
	// contents:write / pull-requests:write / id-token:write.
	for _, want := range []string{"contents: read", "issues: write", "pull-requests: read"} {
		if !strings.Contains(content, want) {
			t.Errorf("engine render permissions must include %q\nfull YAML:\n%s", want, content)
		}
	}
	for _, over := range []string{"contents: write", "pull-requests: write", "id-token: write"} {
		if strings.Contains(content, over) {
			t.Errorf("engine render must NOT grant %q (least privilege)\nfull YAML:\n%s", over, content)
		}
	}

	// The engine render must not fabricate a git remote to reach the API;
	// the only error expected is label-doctor's target-resolution failure
	// (no origin on repoRoot), exactly as the agent-tier fixture sees.
	if err == nil || !strings.Contains(err.Error(), "canonical dispatch labels not ensured") {
		t.Errorf("expected the canonical-dispatch-labels error (render succeeded, label step reached), got: %v", err)
	}

	// stdout states the PAT-free runtime fact for this tier.
	out := stdout.String()
	if !strings.Contains(out, "tier: engine") {
		t.Errorf("stdout should name the engine tier:\n%s", out)
	}
	if !strings.Contains(out, "GITHUB_TOKEN") {
		t.Errorf("stdout should state the GITHUB_TOKEN-only runtime:\n%s", out)
	}
}

// TestRun_DispatchCds_AgentTier_StillRendersClaudeCodeAction is the G2
// separability proof from the other side: with NO --engine flag, the
// default cds dispatch render remains the agent (claude-code-action) tier —
// the two tiers render independently, and adding the engine tier did not
// perturb the agent path.
func TestRun_DispatchCds_AgentTier_StillRendersClaudeCodeAction(t *testing.T) {
	indexPath := writeDispatchFixtureIndex(t)
	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()

	_, _ = Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: indexPath,
		Packages:  []string{"cnos.core", "cnos.cds"},
		Dispatch:  "cds",
		// Engine deliberately false — this is the agent tier.
		Stdout: stdout,
		Stderr: stderr,
	})

	workflowPath := filepath.Join(repoRoot, ".github", "workflows", "cnos-cds-dispatch.yml")
	data, statErr := os.ReadFile(workflowPath)
	if statErr != nil {
		t.Fatalf("expected rendered agent-tier workflow at %s: %v", workflowPath, statErr)
	}
	content := string(data)
	for _, want := range []string{
		"anthropics/claude-code-action@v1",
		"claude_code_oauth_token: ${{ secrets.CLAUDE_CODE_OAUTH_TOKEN }}",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("agent-tier (no --engine) render should still contain %q", want)
		}
	}
	if strings.Contains(content, "Mechanical FSM engine (PAT-free)") {
		t.Error("agent-tier render must NOT contain the engine-tier mechanical step")
	}
}

// TestRun_DispatchCds_Engine_HubLessLabelReconcile is the G3 proof: the
// dispatch install ensures canonical FSM labels HUB-LESSLY. The install
// target here is a plain t.TempDir() with a base install but NO agent hub
// (no .cn/spec, .cn/agent, .cn/threads — cn repo install never scaffolds
// one, unlike cn init) and no git "origin" remote. The label step is still
// REACHED (label-doctor is hub-less — it needs only a git remote +
// GITHUB_TOKEN, not a hub), and surfaces the actionable target-resolution
// error, NOT a "hub missing" error. That distinction is the whole point:
// the labels obligation runs off git+API, so it is never silently skipped
// for want of a hub.
func TestRun_DispatchCds_Engine_HubLessLabelReconcile(t *testing.T) {
	indexPath := writeDispatchFixtureIndex(t)
	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()

	_, err := Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: indexPath,
		Packages:  []string{"cnos.core", "cnos.cds"},
		Dispatch:  "cds",
		Engine:    true,
		Stdout:    stdout,
		Stderr:    stderr,
	})

	// No agent hub was scaffolded by the install (hub-less by construction).
	for _, hubPath := range []string{
		filepath.Join(".cn", "spec"),
		filepath.Join(".cn", "agent"),
		filepath.Join(".cn", "threads"),
	} {
		if _, statErr := os.Stat(filepath.Join(repoRoot, hubPath)); !os.IsNotExist(statErr) {
			t.Errorf("cn repo install must not scaffold an agent hub; found %s", hubPath)
		}
	}

	// The label step was REACHED (label-doctor is hub-less): the error names
	// the canonical-labels obligation and the git-remote target-resolution
	// cause — not a hub-missing failure, and not a silent skip.
	if err == nil {
		t.Fatal("expected the hub-less label-doctor reconcile to surface its target-resolution error")
	}
	if !strings.Contains(err.Error(), "canonical dispatch labels not ensured") {
		t.Errorf("error should name the labels obligation (label-doctor reached), got: %v", err)
	}
	if !strings.Contains(err.Error(), "could not resolve target repo") {
		t.Errorf("error should name the git-remote (hub-less) target cause, got: %v", err)
	}
	if strings.Contains(err.Error(), "hub") {
		t.Errorf("error must be a labels/target-resolution error, not a hub-missing error, got: %v", err)
	}
}
