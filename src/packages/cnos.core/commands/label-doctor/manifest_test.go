package labeldoctor

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

const fixtureManifestJSON = `{
  "schema": "cn.labels.v1",
  "owner": "cnos.core",
  "doctrine": "src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md",
  "labels": [
    {"name": "status:backlog", "description": "Well-formed scope but not yet refined.", "color": "ededed", "owner": "cnos.core", "group": "lifecycle"},
    {"name": "status:review", "description": "Cell complete; awaiting review.", "color": "5319e7", "owner": "cnos.core", "group": "lifecycle"},
    {"name": "status:blocked", "description": "Gated on external input.", "color": "b60205", "owner": "cnos.core", "group": "lifecycle"}
  ]
}`

func TestLoadManifest_HappyPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labels.json")
	writeFile(t, path, fixtureManifestJSON)

	m, err := LoadManifest(path)
	if err != nil {
		t.Fatalf("LoadManifest: %v", err)
	}
	if m.Schema != "cn.labels.v1" {
		t.Errorf("Schema = %q, want cn.labels.v1", m.Schema)
	}
	if len(m.Labels) != 3 {
		t.Fatalf("len(Labels) = %d, want 3", len(m.Labels))
	}
	if m.Labels[1].Name != "status:review" || m.Labels[1].Color != "5319e7" {
		t.Errorf("Labels[1] = %+v, want status:review/5319e7", m.Labels[1])
	}
}

func TestLoadManifest_MissingFile(t *testing.T) {
	if _, err := LoadManifest(filepath.Join(t.TempDir(), "nope.json")); err == nil {
		t.Fatal("expected an error for a missing manifest file")
	}
}

func TestLoadManifest_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labels.json")
	writeFile(t, path, "{not json")
	if _, err := LoadManifest(path); err == nil {
		t.Fatal("expected an error for invalid JSON")
	}
}

func TestLoadManifest_EmptyLabels(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labels.json")
	writeFile(t, path, `{"schema":"cn.labels.v1","owner":"cnos.core","labels":[]}`)
	if _, err := LoadManifest(path); err == nil {
		t.Fatal("expected an error for a manifest with zero labels")
	}
}

// TestResolveDefaultManifestPath_SourceTree covers the primary use case
// named by cnos#493's Gap section: running inside a cnos checkout that
// carries the source-tree labels.json at
// src/packages/cnos.core/labels.json.
func TestResolveDefaultManifestPath_SourceTree(t *testing.T) {
	root := t.TempDir()
	manifestPath := filepath.Join(root, "src", "packages", "cnos.core", "labels.json")
	writeFile(t, manifestPath, fixtureManifestJSON)

	sub := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := resolveDefaultManifestPath(sub)
	if err != nil {
		t.Fatalf("resolveDefaultManifestPath: %v", err)
	}
	if got != manifestPath {
		t.Errorf("got %q, want %q", got, manifestPath)
	}
}

// TestResolveDefaultManifestPath_Vendored covers the installed-package
// use case: an arbitrary tenant repo that ran `cn repo install` has
// cnos.core vendored under .cn/vendor/packages/cnos.core/labels.json,
// with no src/packages/ source tree present at all.
func TestResolveDefaultManifestPath_Vendored(t *testing.T) {
	root := t.TempDir()
	manifestPath := filepath.Join(root, ".cn", "vendor", "packages", "cnos.core", "labels.json")
	writeFile(t, manifestPath, fixtureManifestJSON)

	got, err := resolveDefaultManifestPath(root)
	if err != nil {
		t.Fatalf("resolveDefaultManifestPath: %v", err)
	}
	if got != manifestPath {
		t.Errorf("got %q, want %q", got, manifestPath)
	}
}

// TestResolveDefaultManifestPath_PrefersSourceTree: when both candidates
// exist (a cnos checkout that has also vendored itself, an edge case but
// not impossible), the source-tree path wins — it is listed first in
// defaultManifestRelPaths and is the more "primary" source per the γ
// scaffold's framing.
func TestResolveDefaultManifestPath_PrefersSourceTree(t *testing.T) {
	root := t.TempDir()
	srcPath := filepath.Join(root, "src", "packages", "cnos.core", "labels.json")
	vendoredPath := filepath.Join(root, ".cn", "vendor", "packages", "cnos.core", "labels.json")
	writeFile(t, srcPath, fixtureManifestJSON)
	writeFile(t, vendoredPath, fixtureManifestJSON)

	got, err := resolveDefaultManifestPath(root)
	if err != nil {
		t.Fatalf("resolveDefaultManifestPath: %v", err)
	}
	if got != srcPath {
		t.Errorf("got %q, want the source-tree path %q", got, srcPath)
	}
}

func TestResolveDefaultManifestPath_NotFound(t *testing.T) {
	root := t.TempDir()
	if _, err := resolveDefaultManifestPath(root); err == nil {
		t.Fatal("expected an error when neither candidate path exists anywhere above startDir")
	}
}
