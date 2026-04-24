package pkgbuild

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"

	pkgtypes "github.com/usurobor/cnos/src/go/internal/pkg"
)

func TestBuildOneProducesTarball(t *testing.T) {
	repoRoot := t.TempDir()
	setupTestRepo(t, repoRoot)

	packages, err := DiscoverPackages(repoRoot)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(packages) != 1 {
		t.Fatalf("expected 1 package, got %d", len(packages))
	}

	result := BuildOne(repoRoot, packages[0])
	if result.Err != nil {
		t.Fatalf("build: %v", result.Err)
	}
	if result.SHA256 == "" {
		t.Error("expected non-empty SHA256")
	}

	tarball := filepath.Join(repoRoot, "dist", "packages", "test.pkg-1.0.0.tar.gz")
	if _, err := os.Stat(tarball); err != nil {
		t.Errorf("expected tarball at %s", tarball)
	}
}

func TestBuildOneUpdateIndex(t *testing.T) {
	repoRoot := t.TempDir()
	setupTestRepo(t, repoRoot)

	packages, _ := DiscoverPackages(repoRoot)
	result := BuildOne(repoRoot, packages[0])

	if err := UpdateIndex(repoRoot, []BuildResult{result}); err != nil {
		t.Fatalf("update index: %v", err)
	}

	indexPath := filepath.Join(repoRoot, "dist", "packages", "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read index: %v", err)
	}

	if !strings.Contains(string(data), "test.pkg") {
		t.Error("expected test.pkg in index")
	}
	if !strings.Contains(string(data), "1.0.0") {
		t.Error("expected version 1.0.0 in index")
	}
	if !strings.Contains(string(data), result.SHA256) {
		t.Error("expected SHA256 in index")
	}
}

func TestCheckOneValid(t *testing.T) {
	repoRoot := t.TempDir()
	setupTestRepo(t, repoRoot)

	packages, _ := DiscoverPackages(repoRoot)
	result := CheckOne(packages[0])

	if len(result.Issues) > 0 {
		t.Errorf("expected 0 issues for valid package, got %v", result.Issues)
	}
}

func TestCheckOneMissingContent(t *testing.T) {
	repoRoot := t.TempDir()
	pkgDir := filepath.Join(repoRoot, "src", "packages", "empty.pkg")
	os.MkdirAll(pkgDir, 0755)
	writeManifest(t, pkgDir, "empty.pkg", "1.0.0")

	packages, _ := DiscoverPackages(repoRoot)
	result := CheckOne(packages[0])

	if len(result.Issues) == 0 {
		t.Fatal("expected issues for package with no content")
	}

	found := false
	for _, issue := range result.Issues {
		if strings.Contains(issue, "no content class") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'no content class' issue, got %v", result.Issues)
	}
}

func TestClean(t *testing.T) {
	repoRoot := t.TempDir()
	setupTestRepo(t, repoRoot)

	packages, _ := DiscoverPackages(repoRoot)
	BuildOne(repoRoot, packages[0])

	distDir := filepath.Join(repoRoot, "dist", "packages")
	if _, err := os.Stat(distDir); err != nil {
		t.Fatal("dist/packages should exist after build")
	}

	if err := Clean(repoRoot); err != nil {
		t.Fatalf("clean: %v", err)
	}

	if _, err := os.Stat(distDir); !os.IsNotExist(err) {
		t.Error("dist/packages should be removed after clean")
	}
}

func TestDiscoverPackagesEmpty(t *testing.T) {
	repoRoot := t.TempDir()
	os.MkdirAll(filepath.Join(repoRoot, "src", "packages"), 0755)

	packages, err := DiscoverPackages(repoRoot)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(packages) != 0 {
		t.Errorf("expected 0 packages, got %d", len(packages))
	}
}

func TestUpdateChecksums(t *testing.T) {
	repoRoot := t.TempDir()
	results := []BuildResult{
		{Name: "a.pkg", Version: "1.0.0", SHA256: "aaa111"},
		{Name: "b.pkg", Version: "2.0.0", SHA256: "bbb222"},
	}

	if err := UpdateChecksums(repoRoot, results); err != nil {
		t.Fatalf("update checksums: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(repoRoot, "dist", "packages", "checksums.txt"))
	if err != nil {
		t.Fatalf("read checksums: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "aaa111  a.pkg-1.0.0.tar.gz") {
		t.Error("expected a.pkg checksum entry")
	}
	if !strings.Contains(content, "bbb222  b.pkg-2.0.0.tar.gz") {
		t.Error("expected b.pkg checksum entry")
	}
}

func TestGenerateLockfileData(t *testing.T) {
	results := []BuildResult{
		{Name: "cnos.core", Version: "1.0.0", SHA256: "abc123"},
		{Name: "cnos.eng", Version: "1.0.0", SHA256: "def456"},
		{Name: "cnos.fail", Version: "0.1.0", Err: fmt.Errorf("build failed")},
	}

	data, err := GenerateLockfileData(results)
	if err != nil {
		t.Fatalf("GenerateLockfileData: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, `"cn.lock.v2"`) {
		t.Error("expected cn.lock.v2 schema")
	}
	if !strings.Contains(content, `"cnos.core"`) {
		t.Error("expected cnos.core in lockfile")
	}
	if !strings.Contains(content, `"cnos.eng"`) {
		t.Error("expected cnos.eng in lockfile")
	}
	// Failed packages should be excluded.
	if strings.Contains(content, "cnos.fail") {
		t.Error("failed package should not appear in lockfile")
	}
}

func TestGenerateLockfile(t *testing.T) {
	dir := t.TempDir()
	lockPath := filepath.Join(dir, ".cn", "deps.lock.json")
	results := []BuildResult{
		{Name: "test.pkg", Version: "2.0.0", SHA256: "deadbeef"},
	}

	if err := GenerateLockfile(lockPath, results); err != nil {
		t.Fatalf("GenerateLockfile: %v", err)
	}

	data, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatalf("read lockfile: %v", err)
	}

	var lf struct {
		Schema   string `json:"schema"`
		Packages []struct {
			Name    string `json:"name"`
			Version string `json:"version"`
			SHA256  string `json:"sha256"`
		} `json:"packages"`
	}
	if err := json.Unmarshal(data, &lf); err != nil {
		t.Fatalf("parse lockfile: %v", err)
	}
	if lf.Schema != "cn.lock.v2" {
		t.Errorf("schema = %q, want %q", lf.Schema, "cn.lock.v2")
	}
	if len(lf.Packages) != 1 {
		t.Fatalf("expected 1 package, got %d", len(lf.Packages))
	}
	if lf.Packages[0].Name != "test.pkg" {
		t.Errorf("package name = %q, want %q", lf.Packages[0].Name, "test.pkg")
	}
}

// --- FindContentClasses ---

// TestFindContentClassesEmpty: absent and empty directories count as
// not-present. This is the conservative behaviour `cn build --check`
// relied on (an empty class dir is not a content class).
func TestFindContentClassesEmpty(t *testing.T) {
	pkgDir := t.TempDir()
	if got := FindContentClasses(pkgDir); len(got) != 0 {
		t.Errorf("expected no classes, got %v", got)
	}

	// Empty class dir is present on disk but should still be absent.
	if err := os.MkdirAll(filepath.Join(pkgDir, "skills"), 0755); err != nil {
		t.Fatal(err)
	}
	if got := FindContentClasses(pkgDir); len(got) != 0 {
		t.Errorf("empty skills/ should not count, got %v", got)
	}
}

// TestFindContentClassesAll: every class enumerated in
// pkgtypes.ContentClasses is surfaced when its directory is non-empty,
// and the returned order follows the canonical list. This is the
// single-source-of-truth check for #253.
func TestFindContentClassesAll(t *testing.T) {
	pkgDir := t.TempDir()
	// Create every canonical class directory with a marker file.
	for _, class := range pkgtypes.ContentClasses {
		classDir := filepath.Join(pkgDir, class)
		if err := os.MkdirAll(classDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(classDir, ".keep"), []byte{}, 0644); err != nil {
			t.Fatal(err)
		}
	}

	got := FindContentClasses(pkgDir)
	if len(got) != len(pkgtypes.ContentClasses) {
		t.Fatalf("got %d classes, want %d (%v)", len(got), len(pkgtypes.ContentClasses), got)
	}
	for i, class := range pkgtypes.ContentClasses {
		if got[i] != class {
			t.Errorf("position %d: got %q, want %q", i, got[i], class)
		}
	}
}

// TestFindContentClassesSubset: only non-empty canonical directories
// are returned; unknown directories alongside them are ignored.
func TestFindContentClassesSubset(t *testing.T) {
	pkgDir := t.TempDir()
	// Two canonical classes, one ignored unknown directory.
	for _, class := range []string{"skills", "katas", "not-a-class"} {
		classDir := filepath.Join(pkgDir, class)
		if err := os.MkdirAll(classDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(classDir, ".keep"), []byte{}, 0644); err != nil {
			t.Fatal(err)
		}
	}

	got := FindContentClasses(pkgDir)
	want := []string{"skills", "katas"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i, class := range want {
		if got[i] != class {
			t.Errorf("position %d: got %q, want %q", i, got[i], class)
		}
	}
}

// --- determinism (#264) ---

func TestCreateTarGzDeterministic(t *testing.T) {
	repoRoot := t.TempDir()
	setupTestRepo(t, repoRoot)

	packages, err := DiscoverPackages(repoRoot)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}

	r1 := BuildOne(repoRoot, packages[0])
	if r1.Err != nil {
		t.Fatalf("build 1: %v", r1.Err)
	}

	tarball := filepath.Join(repoRoot, "dist", "packages", "test.pkg-1.0.0.tar.gz")
	data1, err := os.ReadFile(tarball)
	if err != nil {
		t.Fatalf("read tarball 1: %v", err)
	}

	// Rebuild — same source, different wall-clock time.
	r2 := BuildOne(repoRoot, packages[0])
	if r2.Err != nil {
		t.Fatalf("build 2: %v", r2.Err)
	}

	data2, err := os.ReadFile(tarball)
	if err != nil {
		t.Fatalf("read tarball 2: %v", err)
	}

	if r1.SHA256 != r2.SHA256 {
		t.Errorf("SHA256 mismatch: %s vs %s", r1.SHA256, r2.SHA256)
	}

	if len(data1) != len(data2) {
		t.Fatalf("tarball size mismatch: %d vs %d", len(data1), len(data2))
	}
	for i := range data1 {
		if data1[i] != data2[i] {
			t.Fatalf("tarball byte mismatch at offset %d", i)
			break
		}
	}
}

// --- DerivePacklist ---
//
// Invariant: the tarball content is a derivation of the package model
// (manifest + content-class dirs + root README*/LICENSE*), not an
// arbitrary walk of the package root. Stray files and unknown subdirs
// at the root must never appear in the packlist.

// TestDerivePacklistMissingManifest: without cn.package.json there is
// no package; surface the error rather than silently producing an
// empty tarball.
func TestDerivePacklistMissingManifest(t *testing.T) {
	pkgDir := t.TempDir()
	if _, err := DerivePacklist(pkgDir); err == nil {
		t.Fatal("expected error for missing cn.package.json")
	}
}

// TestDerivePacklistManifestOnly: a package with only the manifest
// produces a one-entry packlist. Content classes are optional at the
// derivation layer (CheckOne enforces the "must have content" rule
// separately).
func TestDerivePacklistManifestOnly(t *testing.T) {
	pkgDir := t.TempDir()
	writeManifest(t, pkgDir, "bare.pkg", "1.0.0")

	got, err := DerivePacklist(pkgDir)
	if err != nil {
		t.Fatalf("derive: %v", err)
	}
	want := []string{"cn.package.json"}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestDerivePacklistIncludesContentClasses: every file under a
// recognized, non-empty content-class directory is in the packlist,
// plus the class directory itself. Multi-level nesting is preserved.
func TestDerivePacklistIncludesContentClasses(t *testing.T) {
	pkgDir := t.TempDir()
	writeManifest(t, pkgDir, "test.pkg", "1.0.0")

	mustWrite(t, filepath.Join(pkgDir, "skills", "alpha", "SKILL.md"), "# alpha\n")
	mustWrite(t, filepath.Join(pkgDir, "commands", "daily", "cn-daily"), "#!/bin/sh\n")
	mustWrite(t, filepath.Join(pkgDir, "templates", "t.md"), "template\n")

	got, err := DerivePacklist(pkgDir)
	if err != nil {
		t.Fatalf("derive: %v", err)
	}

	for _, want := range []string{
		"cn.package.json",
		"commands",
		"commands/daily",
		"commands/daily/cn-daily",
		"skills",
		"skills/alpha",
		"skills/alpha/SKILL.md",
		"templates",
		"templates/t.md",
	} {
		if !slices.Contains(got, want) {
			t.Errorf("missing %q in packlist: %v", want, got)
		}
	}
}

// TestDerivePacklistIncludesRootFiles: non-dotfile entries at the
// package root ship — README/LICENSE for projects that follow
// convention, AGENTS.md/HEARTBEAT.md/TOOLS.md for cnos.core, lib.sh
// for cnos.kata. All are exercised together so the test breaks if
// the rule ever tightens without a migration for the load-bearing
// root files that exist today.
func TestDerivePacklistIncludesRootFiles(t *testing.T) {
	pkgDir := t.TempDir()
	writeManifest(t, pkgDir, "test.pkg", "1.0.0")
	for _, name := range []string{"README.md", "LICENSE", "AGENTS.md", "lib.sh"} {
		mustWrite(t, filepath.Join(pkgDir, name), name+"\n")
	}

	got, err := DerivePacklist(pkgDir)
	if err != nil {
		t.Fatalf("derive: %v", err)
	}

	for _, want := range []string{"README.md", "LICENSE", "AGENTS.md", "lib.sh"} {
		if !slices.Contains(got, want) {
			t.Errorf("missing %q in packlist: %v", want, got)
		}
	}
}

// TestDerivePacklistExcludesStrayEntries: the negative space. Dotfile
// entries at the root and any unrecognized *subdirectory* at the root
// never ship. This is the regression the issue calls out: the walk-
// based implementation included every directory it encountered, so
// `scratch/wip.txt` leaked into tarballs. The strict README*/LICENSE*-
// only root allowlist is deferred (see DerivePacklist comment).
func TestDerivePacklistExcludesStrayEntries(t *testing.T) {
	pkgDir := t.TempDir()
	writeManifest(t, pkgDir, "test.pkg", "1.0.0")
	mustWrite(t, filepath.Join(pkgDir, "skills", "alpha", "SKILL.md"), "# alpha\n")

	// Dotfiles at root.
	mustWrite(t, filepath.Join(pkgDir, ".DS_Store"), "junk")
	mustWrite(t, filepath.Join(pkgDir, ".gitignore"), "*.swp\n")
	// Unrecognized subdirectory at root.
	mustWrite(t, filepath.Join(pkgDir, "scratch", "wip.txt"), "wip")
	mustWrite(t, filepath.Join(pkgDir, "tmp", "cache.bin"), "cache")

	got, err := DerivePacklist(pkgDir)
	if err != nil {
		t.Fatalf("derive: %v", err)
	}

	for _, forbidden := range []string{
		".DS_Store",
		".gitignore",
		"scratch",
		"scratch/wip.txt",
		"tmp",
		"tmp/cache.bin",
	} {
		if slices.Contains(got, forbidden) {
			t.Errorf("packlist must not contain %q, got %v", forbidden, got)
		}
	}
}

// TestDerivePacklistSortedLexically: determinism of the tarball's
// byte output in createTarGz depends on a stably-sorted input.
func TestDerivePacklistSortedLexically(t *testing.T) {
	pkgDir := t.TempDir()
	writeManifest(t, pkgDir, "test.pkg", "1.0.0")
	mustWrite(t, filepath.Join(pkgDir, "skills", "z", "SKILL.md"), "z")
	mustWrite(t, filepath.Join(pkgDir, "skills", "a", "SKILL.md"), "a")
	mustWrite(t, filepath.Join(pkgDir, "templates", "t.md"), "t")
	mustWrite(t, filepath.Join(pkgDir, "commands", "c", "cn-c"), "c")

	got, err := DerivePacklist(pkgDir)
	if err != nil {
		t.Fatalf("derive: %v", err)
	}
	if !sort.StringsAreSorted(got) {
		t.Errorf("packlist not sorted: %v", got)
	}
}

// TestDerivePacklistEmptyClassDirSkipped: an empty content-class
// directory does not contribute to the packlist. Matches the
// FindContentClasses "empty == absent" rule so the derivation and the
// check surface agree on what counts as a present class.
func TestDerivePacklistEmptyClassDirSkipped(t *testing.T) {
	pkgDir := t.TempDir()
	writeManifest(t, pkgDir, "test.pkg", "1.0.0")
	if err := os.MkdirAll(filepath.Join(pkgDir, "skills"), 0755); err != nil {
		t.Fatal(err)
	}

	got, err := DerivePacklist(pkgDir)
	if err != nil {
		t.Fatalf("derive: %v", err)
	}
	if slices.Contains(got, "skills") {
		t.Errorf("empty skills/ must not appear in packlist: %v", got)
	}
}

// TestDerivePacklistReusesContentClassConstant: fail loudly if a
// future class is added to pkgtypes.ContentClasses but the derivation
// silently drops it. This is the single-source-of-truth guard for
// the shared constant (issue #262 AC2).
func TestDerivePacklistReusesContentClassConstant(t *testing.T) {
	pkgDir := t.TempDir()
	writeManifest(t, pkgDir, "test.pkg", "1.0.0")
	for _, class := range pkgtypes.ContentClasses {
		mustWrite(t, filepath.Join(pkgDir, class, "x", "marker"), class)
	}

	got, err := DerivePacklist(pkgDir)
	if err != nil {
		t.Fatalf("derive: %v", err)
	}
	for _, class := range pkgtypes.ContentClasses {
		want := class + "/x/marker"
		if !slices.Contains(got, want) {
			t.Errorf("missing %q: every canonical class must flow through DerivePacklist", want)
		}
	}
}

// --- Tarball content assertions ---
//
// These tests reach past createTarGz's SHA-256 and inspect what
// actually ships. They prove the negative space: stray files are
// excluded end-to-end, and visibility metadata does not gate
// packaging.

// TestBuildOneExcludesStrayFilesFromTarball: integration proof that
// the stray files which DerivePacklist filters out never reach the
// produced .tar.gz.
func TestBuildOneExcludesStrayFilesFromTarball(t *testing.T) {
	repoRoot := t.TempDir()
	setupTestRepo(t, repoRoot)

	pkgDir := filepath.Join(repoRoot, "src", "packages", "test.pkg")
	mustWrite(t, filepath.Join(pkgDir, ".DS_Store"), "junk")
	mustWrite(t, filepath.Join(pkgDir, ".gitignore"), "*.swp\n")
	mustWrite(t, filepath.Join(pkgDir, "scratch", "wip.txt"), "wip")

	packages, err := DiscoverPackages(repoRoot)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	result := BuildOne(repoRoot, packages[0])
	if result.Err != nil {
		t.Fatalf("build: %v", result.Err)
	}

	entries := readTarEntries(t, result.TarballPath)
	for _, forbidden := range []string{
		".DS_Store", ".gitignore",
		"scratch", "scratch/", "scratch/wip.txt",
	} {
		if _, ok := entries[forbidden]; ok {
			t.Errorf("tarball must not contain %q, got entries: %v", forbidden, sortedKeys(entries))
		}
	}
}

// TestBuildOneIncludesInternalVisibilitySkill: a skill with
// visibility: internal is still shipped. Visibility affects
// activation/discovery (per #261), not packaging (issue #262 AC6).
func TestBuildOneIncludesInternalVisibilitySkill(t *testing.T) {
	repoRoot := t.TempDir()
	pkgDir := filepath.Join(repoRoot, "src", "packages", "test.pkg")
	mustWrite(t, filepath.Join(pkgDir, "skills", "public", "SKILL.md"),
		"---\nname: public\n---\n# public\n")
	mustWrite(t, filepath.Join(pkgDir, "skills", "secret", "SKILL.md"),
		"---\nname: secret\nvisibility: internal\n---\n# secret\n")
	writeManifest(t, pkgDir, "test.pkg", "1.0.0")

	packages, err := DiscoverPackages(repoRoot)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	result := BuildOne(repoRoot, packages[0])
	if result.Err != nil {
		t.Fatalf("build: %v", result.Err)
	}

	entries := readTarEntries(t, result.TarballPath)
	for _, want := range []string{
		"skills/public/SKILL.md",
		"skills/secret/SKILL.md",
	} {
		if _, ok := entries[want]; !ok {
			t.Errorf("tarball missing %q — visibility must not affect packaging. entries: %v",
				want, sortedKeys(entries))
		}
	}
}

// --- helpers ---

// mustWrite creates parent dirs for path and writes content. Fails
// the test on any IO error.
func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// readTarEntries returns the set of entry names in tarballPath. Keys
// keep whatever separator the producer wrote (regular files have no
// trailing "/", directories do). Bodies are not returned.
func readTarEntries(t *testing.T, tarballPath string) map[string]bool {
	t.Helper()
	raw, err := os.ReadFile(tarballPath)
	if err != nil {
		t.Fatalf("read tarball: %v", err)
	}
	gr, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("gzip reader: %v", err)
	}
	defer gr.Close()

	entries := make(map[string]bool)
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("tar next: %v", err)
		}
		entries[hdr.Name] = true
		// We do not care about bodies here; skip.
	}
	return entries
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func setupTestRepo(t *testing.T, repoRoot string) {
	t.Helper()
	pkgDir := filepath.Join(repoRoot, "src", "packages", "test.pkg")
	skillDir := filepath.Join(pkgDir, "skills", "test-skill")
	os.MkdirAll(skillDir, 0755)
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# Test Skill\n"), 0644)
	writeManifest(t, pkgDir, "test.pkg", "1.0.0")
}

func writeManifest(t *testing.T, pkgDir, name, version string) {
	t.Helper()
	manifest := PackageManifest{Name: name, Version: version}
	data, _ := json.MarshalIndent(manifest, "", "  ")
	os.WriteFile(filepath.Join(pkgDir, "cn.package.json"), data, 0644)
}
