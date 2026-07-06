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
	"path/filepath"
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

func TestValidateDispatch(t *testing.T) {
	cases := []struct {
		in      string
		wantErr bool
		want609 bool
	}{
		{"", false, false},
		{"none", false, false},
		{"cds", true, true},
		{"bogus", true, false},
	}
	for _, c := range cases {
		err := validateDispatch(c.in)
		if (err != nil) != c.wantErr {
			t.Errorf("validateDispatch(%q) err = %v, wantErr %v", c.in, err, c.wantErr)
			continue
		}
		if c.want609 && (err == nil || !strings.Contains(err.Error(), "#609")) {
			t.Errorf("validateDispatch(%q) = %v, want mention of #609", c.in, err)
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

// AC9: --dispatch cds fails explicitly, before any write, with no partial
// .github/workflows/ file.
func TestRun_DispatchCds_FailsWithNoPartialWrite(t *testing.T) {
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
		t.Fatal("expected --dispatch cds to fail")
	}
	if !strings.Contains(err.Error(), "#609") {
		t.Errorf("error should name #609, got: %v", err)
	}
	if !strings.Contains(stderr.String(), "✗ --dispatch cds requires generalized wake renderer support (#609)") {
		t.Errorf("stderr should carry the exact dispatch-guard message, got: %q", stderr.String())
	}
	entries, rerr := os.ReadDir(repoRoot)
	if rerr != nil {
		t.Fatal(rerr)
	}
	if len(entries) != 0 {
		t.Errorf("--dispatch cds must write nothing; found: %v", entries)
	}
	if _, statErr := os.Stat(filepath.Join(repoRoot, ".github", "workflows", "cnos-cds-dispatch.yml")); !os.IsNotExist(statErr) {
		t.Error("no partial .github/workflows/cnos-cds-dispatch.yml may exist")
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
