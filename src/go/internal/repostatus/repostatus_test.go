package repostatus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/dispatchrender"
	"github.com/usurobor/cnos/src/go/internal/pkg"
	"github.com/usurobor/cnos/src/go/internal/repostate"
)

func sha256Hex(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// setupFixtureRepo builds a minimal but complete "installed" repo layout
// directly on disk (no network, no tarballs) — the pure shape
// internal/repoinstall's Run would have produced — so repostatus tests
// exercise reading, not the installer's own write path (already covered by
// internal/repoinstall's tests).
func setupFixtureRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	manifest := pkg.Manifest{Schema: "cn.deps.v1", Profile: "cds", Packages: []pkg.ManifestDep{{Name: "cnos.core", Version: "1.0.0"}}}
	depsJSON, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, ".cn", "deps.json"), string(depsJSON)+"\n")

	lock := pkg.Lockfile{Schema: "cn.lock.v2", Packages: []pkg.LockedDep{{Name: "cnos.core", Version: "1.0.0", SHA256: "abc123"}}}
	lockJSON, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, ".cn", "deps.lock.json"), string(lockJSON)+"\n")

	writeFile(t, filepath.Join(root, ".cn", "vendor", "packages", "cnos.core", "cn.package.json"), `{"name": "cnos.core", "version": "1.0.0"}`)

	state := repostate.RepoState{
		Schema: repostate.Schema, Profile: "cds",
		Source: repostate.Source{Channel: "stable", Release: "1.0.0", Index: "https://example/index.json"},
		ManagedFiles: []repostate.ManagedFile{
			{Path: ".cn/deps.json", Kind: repostate.KindManifest, ID: "deps-manifest", SHA256: sha256Hex(t, filepath.Join(root, ".cn", "deps.json"))},
			{Path: ".cn/deps.lock.json", Kind: repostate.KindLockfile, ID: "deps-lockfile", SHA256: sha256Hex(t, filepath.Join(root, ".cn", "deps.lock.json"))},
		},
		ManagedDirs: []repostate.ManagedDir{{Path: ".cn/vendor/packages/cnos.core", Package: "cnos.core"}},
		ExternalExpectations: repostate.ExternalExpectations{
			Labels: repostate.LabelExpectations{Mode: repostate.LabelModeEnsure, Source: "cnos.core/labels.json", DeleteOnUninstall: false},
		},
	}
	data, err := state.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".cn", "repo.state.json"), data, 0644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestRun_CleanInstall_NoDrift(t *testing.T) {
	root := setupFixtureRepo(t)
	st, err := Run(context.Background(), Options{RepoRoot: root, SkipNetwork: true})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if st.Schema != StatusSchema {
		t.Errorf("schema = %q, want %q", st.Schema, StatusSchema)
	}
	if !st.Ledger.Present || st.Ledger.Reconstructed {
		t.Errorf("Ledger = %+v, want Present=true Reconstructed=false", st.Ledger)
	}
	if len(st.Packages) != 1 || !st.Packages[0].InSync {
		t.Errorf("Packages = %+v, want one in-sync cnos.core entry", st.Packages)
	}
	if st.Dispatch.Present {
		t.Errorf("Dispatch.Present = true, want false (no workflow in fixture)")
	}
	if len(st.LocalEdits) != 0 {
		t.Errorf("LocalEdits = %+v, want none", st.LocalEdits)
	}
	if st.Drift {
		t.Error("Drift = true, want false for a clean fixture")
	}
}

// TestRun_NeverWrites is the P3 oracle: repeated Run calls (default and
// --json-equivalent structured calls) must never touch RepoRoot.
func TestRun_NeverWrites(t *testing.T) {
	root := setupFixtureRepo(t)

	before := map[string][]byte{}
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		data, _ := os.ReadFile(path)
		before[path] = data
		return nil
	})

	for i := 0; i < 3; i++ {
		if _, err := Run(context.Background(), Options{RepoRoot: root, SkipNetwork: true}); err != nil {
			t.Fatalf("Run: %v", err)
		}
	}

	after := map[string][]byte{}
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		data, _ := os.ReadFile(path)
		after[path] = data
		return nil
	})

	if len(before) != len(after) {
		t.Fatalf("file count changed: before=%d after=%d", len(before), len(after))
	}
	for path, data := range before {
		if string(after[path]) != string(data) {
			t.Errorf("%s content changed after Run", path)
		}
	}
}

func TestRun_LocalEditDetected(t *testing.T) {
	root := setupFixtureRepo(t)
	// Tamper deps.json after the ledger was computed — a local edit.
	writeFile(t, filepath.Join(root, ".cn", "deps.json"), `{"schema":"cn.deps.v1","profile":"tampered","packages":[]}`)

	st, err := Run(context.Background(), Options{RepoRoot: root, SkipNetwork: true})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(st.LocalEdits) != 1 || st.LocalEdits[0].Path != ".cn/deps.json" || st.LocalEdits[0].Classification != "user_edit" {
		t.Errorf("LocalEdits = %+v, want one user_edit entry for .cn/deps.json", st.LocalEdits)
	}
	if !st.Drift {
		t.Error("Drift = false, want true when a managed file was locally edited")
	}
}

// TestRun_RemovedManagedFile covers the LocalEdit "removed" classification:
// a ledger-recorded managed file (here, a synthetic gitignore entry added
// to the fixture) that no longer exists on disk at all.
func TestRun_RemovedManagedFile(t *testing.T) {
	root := setupFixtureRepo(t)
	gitignorePath := filepath.Join(root, ".gitignore")
	writeFile(t, gitignorePath, ".cn/vendor/\n")

	statePath := filepath.Join(root, ".cn", "repo.state.json")
	data, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	state, err := repostate.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	state.ManagedFiles = append(state.ManagedFiles, repostate.ManagedFile{
		Path: ".gitignore", Kind: repostate.KindGitignore, ID: "gitignore", SHA256: sha256Hex(t, gitignorePath),
	})
	newData, err := state.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(statePath, newData, 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove(gitignorePath); err != nil {
		t.Fatal(err)
	}

	st, err := Run(context.Background(), Options{RepoRoot: root, SkipNetwork: true})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	found := false
	for _, e := range st.LocalEdits {
		if e.Path == ".gitignore" {
			found = true
			if e.Classification != "removed" {
				t.Errorf(".gitignore classification = %q, want \"removed\"", e.Classification)
			}
		}
	}
	if !found {
		t.Errorf("LocalEdits = %+v, want an entry for the removed .gitignore", st.LocalEdits)
	}
}

func TestRun_OrphanVendoredPackage(t *testing.T) {
	root := setupFixtureRepo(t)
	writeFile(t, filepath.Join(root, ".cn", "vendor", "packages", "cnos.legacy-foo", "cn.package.json"), `{"name": "cnos.legacy-foo", "version": "0.1.0"}`)

	st, err := Run(context.Background(), Options{RepoRoot: root, SkipNetwork: true})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(st.OrphanPackages) != 1 || st.OrphanPackages[0] != "cnos.legacy-foo" {
		t.Errorf("OrphanPackages = %v, want [cnos.legacy-foo]", st.OrphanPackages)
	}
	if !st.Drift {
		t.Error("Drift = false, want true when an orphan vendored package exists")
	}
}

func TestRun_NoDepsJSON_Errors(t *testing.T) {
	root := t.TempDir()
	if _, err := Run(context.Background(), Options{RepoRoot: root, SkipNetwork: true}); err == nil {
		t.Error("Run on a repo with no .cn/deps.json: got nil error, want non-nil")
	}
}

func TestRun_NoLedger_DegradesGracefully(t *testing.T) {
	root := setupFixtureRepo(t)
	if err := os.Remove(filepath.Join(root, ".cn", "repo.state.json")); err != nil {
		t.Fatal(err)
	}
	st, err := Run(context.Background(), Options{RepoRoot: root, SkipNetwork: true})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if st.Ledger.Present || !st.Ledger.Reconstructed {
		t.Errorf("Ledger = %+v, want Present=false Reconstructed=true", st.Ledger)
	}
	// Packages are still reported from deps.json/deps.lock.json/vendor
	// directly — the ledger's absence only degrades local_edits/dispatch.
	if len(st.Packages) != 1 || !st.Packages[0].InSync {
		t.Errorf("Packages = %+v, want still-reported in-sync cnos.core", st.Packages)
	}
	if len(st.LocalEdits) != 0 {
		t.Errorf("LocalEdits = %+v, want none (no ledger to compare against)", st.LocalEdits)
	}
}

// --- dispatch drift classification (A2) ---

func repoSrcRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..")
}

// setupDispatchFixture extends setupFixtureRepo with a REAL vendored
// cn-install-wake renderer script (so dispatchrender.Render actually
// works) and an initial workflow render + matching ledger record, using
// the engine tier (no PAT/identity fixture wiring needed).
func setupDispatchFixture(t *testing.T) (root string, wf repostate.ManagedFile) {
	t.Helper()
	root = setupFixtureRepo(t)

	rendererScript, err := os.ReadFile(filepath.Join(repoSrcRoot(t), "src", "packages", "cnos.core", "commands", "install-wake", "cn-install-wake"))
	if err != nil {
		t.Fatalf("read real cn-install-wake: %v", err)
	}
	rendererPath := filepath.Join(root, ".cn", "vendor", "packages", "cnos.core", "commands", "install-wake", "cn-install-wake")
	if err := os.MkdirAll(filepath.Dir(rendererPath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(rendererPath, rendererScript, 0755); err != nil {
		t.Fatal(err)
	}

	// The renderer resolves CN_PACKAGE_ROOT from its own script_dir
	// (cnos.core) when unset, then re-pins to the sibling package that
	// actually owns orchestrators/cds-dispatch/SKILL.md — mirroring
	// internal/repoinstall's own writeDispatchFixtureIndex, which vendors
	// cnos.cds alongside cnos.core for exactly this reason.
	skillMD, err := os.ReadFile(filepath.Join(repoSrcRoot(t), "src", "packages", "cnos.cds", "orchestrators", "cds-dispatch", "SKILL.md"))
	if err != nil {
		t.Fatalf("read real cds-dispatch SKILL.md: %v", err)
	}
	skillPath := filepath.Join(root, ".cn", "vendor", "packages", "cnos.cds", "orchestrators", "cds-dispatch", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(skillPath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(skillPath, skillMD, 0644); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, ".cn", "vendor", "packages", "cnos.cds", "cn.package.json"), `{"name": "cnos.cds", "version": "1.0.0"}`)

	// cnos.cds is vendored only so the renderer can find its sibling
	// SKILL.md (see comment above) — it isn't a "real" desired package in
	// this fixture. Record it in the lockfile too so it isn't mistaken
	// for an orphan vendored package by the orphan-detection test below.
	lockPath := filepath.Join(root, ".cn", "deps.lock.json")
	lockData, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatal(err)
	}
	lock, err := pkg.ParseLockfile(lockData)
	if err != nil {
		t.Fatal(err)
	}
	lock.Packages = append(lock.Packages, pkg.LockedDep{Name: "cnos.cds", Version: "1.0.0", SHA256: "def456"})
	newLockData, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lockPath, append(newLockData, '\n'), 0644); err != nil {
		t.Fatal(err)
	}

	workflowPath := filepath.Join(root, ".github", "workflows", "cnos-cds-dispatch.yml")
	if err := dispatchrender.Render(context.Background(), dispatchrender.Options{
		RepoRoot: root, RendererPackage: "cnos.core", Tier: repostate.TierEngine, Agent: "sigma", OutPath: workflowPath,
	}); err != nil {
		t.Fatalf("initial render: %v", err)
	}

	wf = repostate.ManagedFile{
		Path: ".github/workflows/cnos-cds-dispatch.yml", Kind: repostate.KindWorkflow, ID: "cnos-cds-dispatch",
		SHA256: sha256Hex(t, workflowPath), Tier: repostate.TierEngine,
		Renderer: "cnos.core/install-wake", RendererPackage: "cnos.core", RendererVersionSource: "lock", Agent: "sigma",
	}

	// Add the workflow record to the ledger already on disk.
	data, err := os.ReadFile(filepath.Join(root, ".cn", "repo.state.json"))
	if err != nil {
		t.Fatal(err)
	}
	state, err := repostate.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	// deps.lock.json's content changed (cnos.cds appended above) after
	// setupFixtureRepo computed the ledger's original sha256 for it —
	// refresh that entry so this fixture's deps.lock.json isn't itself
	// mistaken for a local edit.
	if lockRecord := state.FindManagedFile(".cn/deps.lock.json"); lockRecord != nil {
		lockRecord.SHA256 = sha256Hex(t, lockPath)
	}
	state.ManagedFiles = append(state.ManagedFiles, wf)
	newData, err := state.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".cn", "repo.state.json"), newData, 0644); err != nil {
		t.Fatal(err)
	}
	return root, wf
}

func TestDispatchStatus_MatchesLedger(t *testing.T) {
	root, _ := setupDispatchFixture(t)
	st, err := Run(context.Background(), Options{RepoRoot: root, SkipNetwork: true})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !st.Dispatch.Present || st.Dispatch.Drift != DriftMatchesLedger {
		t.Errorf("Dispatch = %+v, want Present=true Drift=%s", st.Dispatch, DriftMatchesLedger)
	}
	if st.Drift {
		t.Error("Drift = true, want false when the workflow matches the ledger")
	}
}

func TestDispatchStatus_RendererMoved(t *testing.T) {
	root, _ := setupDispatchFixture(t)
	// Simulate an old ledger sha (as if the renderer output changed since
	// install) by corrupting the RECORDED sha while leaving the live file
	// as the actual fresh-render output — a fresh re-render will match the
	// live file exactly, classifying this as renderer_moved.
	statePath := filepath.Join(root, ".cn", "repo.state.json")
	data, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatal(err)
	}
	state, err := repostate.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	wfRecord := state.FindManagedFile(".github/workflows/cnos-cds-dispatch.yml")
	wfRecord.SHA256 = "0000000000000000000000000000000000000000000000000000000000000000"[:64]
	newData, err := state.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(statePath, newData, 0644); err != nil {
		t.Fatal(err)
	}

	st, err := Run(context.Background(), Options{RepoRoot: root, SkipNetwork: true})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if st.Dispatch.Drift != DriftRendererMoved {
		t.Errorf("Dispatch.Drift = %q, want %q", st.Dispatch.Drift, DriftRendererMoved)
	}
}

func TestDispatchStatus_UserEdit(t *testing.T) {
	root, _ := setupDispatchFixture(t)
	workflowPath := filepath.Join(root, ".github", "workflows", "cnos-cds-dispatch.yml")
	data, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(workflowPath, append(data, []byte("\n# hand-edited\n")...), 0644); err != nil {
		t.Fatal(err)
	}

	st, err := Run(context.Background(), Options{RepoRoot: root, SkipNetwork: true})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if st.Dispatch.Drift != DriftUserEdit {
		t.Errorf("Dispatch.Drift = %q, want %q", st.Dispatch.Drift, DriftUserEdit)
	}
	if !st.Drift {
		t.Error("Drift = false, want true for a hand-edited dispatch workflow")
	}
}
