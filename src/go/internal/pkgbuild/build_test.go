package pkgbuild

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

// --- helpers ---

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
