package repoinstall

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/pkg"
	"github.com/usurobor/cnos/src/go/internal/repostate"
)

// TestRun_WritesRepoState_Deterministic covers cnos#656 P1/P2/A3: install
// writes .cn/repo.state.json, it validates against schemas/repo_state.cue's
// closed #RepoState shape (checked structurally here via the Go type — the
// CUE gate itself is scripts/ci/validate-repo-state.sh's job), no
// version/sha256 duplication under a "packages" key, and a same-input
// rerun produces a byte-identical ledger (the idempotence contract).
func TestRun_WritesRepoState_Deterministic(t *testing.T) {
	tarData, tarSHA := makeTarGz(t, map[string]string{"cn.package.json": `{"name": "cnos.core", "version": "1.0.0"}`})
	idx := pkg.PackageIndex{
		Schema: "cn.package-index.v1",
		Packages: map[string]map[string]pkg.IndexEntry{
			"cnos.core": {"1.0.0": {URL: "cnos.core-1.0.0.tar.gz", SHA256: tarSHA}},
		},
	}
	indexBody, err := json.Marshal(idx)
	if err != nil {
		t.Fatal(err)
	}

	repo := "acme/widgets"
	api, dl := newFakeGitHub(t, repo, "1.0.0", indexBody, map[string][]byte{"cnos.core-1.0.0.tar.gz": tarData})

	repoRoot := t.TempDir()
	runOnce := func() []byte {
		stdout, stderr := noopStdio()
		if _, err := Run(context.Background(), Options{
			RepoRoot:     repoRoot,
			Release:      "latest",
			Packages:     []string{"cnos.core"},
			Repo:         repo,
			APIBaseURL:   api.URL,
			DownloadBase: dl.URL,
			Stdout:       stdout,
			Stderr:       stderr,
		}); err != nil {
			t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
		}
		data, err := os.ReadFile(filepath.Join(repoRoot, ".cn", "repo.state.json"))
		if err != nil {
			t.Fatalf("read repo.state.json: %v", err)
		}
		return data
	}

	first := runOnce()

	var state repostate.RepoState
	if err := json.Unmarshal(first, &state); err != nil {
		t.Fatalf("parse repo.state.json: %v", err)
	}
	if state.Schema != repostate.Schema {
		t.Errorf("schema = %q, want %q", state.Schema, repostate.Schema)
	}
	if state.Source.Release != "1.0.0" {
		t.Errorf("source.release = %q, want 1.0.0", state.Source.Release)
	}
	if f := state.FindManagedFile(".cn/deps.json"); f == nil || f.SHA256 == "" {
		t.Errorf("managed_files missing a hashed .cn/deps.json entry: %+v", f)
	}
	if d := state.FindManagedDir("cnos.core"); d == nil || d.Path != ".cn/vendor/packages/cnos.core" {
		t.Errorf("managed_dirs missing cnos.core: %+v", d)
	}

	// No top-level "packages" key (design doc A1 — no lock duplication).
	var raw map[string]any
	if err := json.Unmarshal(first, &raw); err != nil {
		t.Fatal(err)
	}
	if _, present := raw["packages"]; present {
		t.Error("repo.state.json carries a top-level \"packages\" key — violates A1 (no lock duplication)")
	}
	for k := range raw {
		if k == "timestamp" || k == "created_at" || k == "installed_at" || k == "updated_at" {
			t.Errorf("repo.state.json carries a timestamp-shaped field %q — violates P1 (timestamp-free)", k)
		}
	}

	second := runOnce()
	if string(first) != string(second) {
		t.Errorf("repo.state.json not idempotent across reruns:\nfirst=%s\nsecond=%s", first, second)
	}
}

// TestRun_WritesRepoState_DispatchWorkflowRenderContract covers P2: a
// dispatch-cds install's workflow managed_files entry carries the full
// render contract (tier/renderer/renderer_package/renderer_version_source/
// agent), not just path/kind/id/sha256/tier.
func TestRun_WritesRepoState_DispatchWorkflowRenderContract(t *testing.T) {
	indexPath := writeDispatchFixtureIndex(t)
	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()
	_, err := Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: indexPath,
		Packages:  []string{"cnos.core", "cnos.cds"},
		Dispatch:  "cds",
		Engine:    true, // engine tier needs no PAT/identity fixture wiring
		Stdout:    stdout,
		Stderr:    stderr,
	})
	// ensureCanonicalDispatchLabels will fail in this fixture (no git
	// remote configured) — that's expected and orthogonal to what this
	// test checks (the workflow file + its would-be ledger record exist
	// before that later failure). Only fail the test on an error that
	// isn't the known label-doctor gap.
	if err == nil || !strings.Contains(err.Error(), "canonical dispatch labels not ensured") {
		t.Fatalf("Run: unexpected error %v\nstderr: %s", err, stderr.String())
	}

	workflowPath := dispatchWorkflowPath(repoRoot)
	if _, statErr := os.Stat(workflowPath); statErr != nil {
		t.Fatalf("dispatch workflow was not rendered: %v", statErr)
	}

	// The render succeeded even though label-doctor (checked above) then
	// failed — Run still writes the ledger in this case (the fix this
	// test exists to pin): losing the ledger because an orthogonal
	// obligation failed would defeat A3's backfill contract for an
	// otherwise-successful install.
	data, err := os.ReadFile(filepath.Join(repoRoot, ".cn", "repo.state.json"))
	if err != nil {
		t.Fatalf("read repo.state.json: %v", err)
	}
	var state repostate.RepoState
	if err := json.Unmarshal(data, &state); err != nil {
		t.Fatalf("parse repo.state.json: %v", err)
	}
	wf := state.FindManagedFile(".github/workflows/cnos-cds-dispatch.yml")
	if wf == nil {
		t.Fatal("no workflow managed_files entry")
	}
	if wf.Kind != repostate.KindWorkflow {
		t.Errorf("kind = %q, want %q", wf.Kind, repostate.KindWorkflow)
	}
	if wf.Tier != repostate.TierEngine {
		t.Errorf("tier = %q, want %q", wf.Tier, repostate.TierEngine)
	}
	if wf.Renderer == "" || wf.RendererPackage == "" || wf.RendererVersionSource == "" || wf.Agent == "" {
		t.Errorf("workflow render contract incomplete: %+v", wf)
	}
	if wf.SHA256 == "" {
		t.Error("workflow managed_files entry has empty sha256")
	}
}

// TestRun_DispatchCds_StaleWorkflowFileFromPriorRun_LedgerNotCorrupted pins
// a bug an adversarial review caught in the first draft of the
// dispatchErr/canWriteLedger restructure: `canWriteLedger` was derived from
// os.Stat(dispatchWorkflowPath(...)) — "does a file exist at the fixed
// output path" — which is ALSO true when a PRIOR successful run rendered
// the file and THIS run's identity gate fails before ever calling the
// renderer again. Writing the ledger in that case described an untouched
// file using THIS run's (different) opts.Agent/opts.WorkflowPatSecret,
// silently corrupting the render-contract metadata for a file nothing this
// run touched. The fix threads runDispatchCds's own `rendered` bool through
// instead of re-deriving it from file presence.
func TestRun_DispatchCds_StaleWorkflowFileFromPriorRun_LedgerNotCorrupted(t *testing.T) {
	indexPath := writeDispatchFixtureIndex(t)
	repoRoot := t.TempDir()
	stdout, stderr := noopStdio()

	// First run: sigma/engine tier, fully succeeds (engine tier skips the
	// label-doctor gap entirely — no, it doesn't; engine still calls
	// ensureCanonicalDispatchLabels. Use a real successful render by
	// tolerating the expected label-doctor failure, same as
	// TestRun_WritesRepoState_DispatchWorkflowRenderContract, then assert
	// the ledger this first run produced.
	if _, err := Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: indexPath,
		Packages:  []string{"cnos.core", "cnos.cds"},
		Dispatch:  "cds",
		Engine:    true,
		Stdout:    stdout,
		Stderr:    stderr,
	}); err == nil || !strings.Contains(err.Error(), "canonical dispatch labels not ensured") {
		t.Fatalf("first Run: unexpected error %v\nstderr: %s", err, stderr.String())
	}

	firstData, err := os.ReadFile(filepath.Join(repoRoot, ".cn", "repo.state.json"))
	if err != nil {
		t.Fatalf("read repo.state.json after first run: %v", err)
	}
	var firstState repostate.RepoState
	if err := json.Unmarshal(firstData, &firstState); err != nil {
		t.Fatal(err)
	}
	firstWF := firstState.FindManagedFile(".github/workflows/cnos-cds-dispatch.yml")
	if firstWF == nil || firstWF.Agent != "sigma" || firstWF.Tier != repostate.TierEngine {
		t.Fatalf("first run's ledger workflow entry unexpected: %+v", firstWF)
	}
	firstWorkflowBytes, err := os.ReadFile(dispatchWorkflowPath(repoRoot))
	if err != nil {
		t.Fatal(err)
	}

	// Second run: a DIFFERENT, non-sigma agent with no --workflow-pat-secret
	// and Engine: false — this fails at the identity gate INSIDE
	// runDispatchCds, before dispatchrender.Render is ever called again.
	// The workflow file on disk is the first run's sigma/engine render,
	// untouched.
	stdout2, stderr2 := noopStdio()
	_, err = Run(context.Background(), Options{
		RepoRoot:  repoRoot,
		IndexPath: indexPath,
		Packages:  []string{"cnos.core", "cnos.cds"},
		Dispatch:  "cds",
		Agent:     "acme",
		Stdout:    stdout2,
		Stderr:    stderr2,
	})
	if err == nil || !strings.Contains(err.Error(), "--workflow-pat-secret") {
		t.Fatalf("second Run: expected the identity gate to fail, got %v\nstderr: %s", err, stderr2.String())
	}

	// The live workflow file must be byte-identical to the first run's
	// render — the second run's failed identity gate must not have
	// touched it.
	secondWorkflowBytes, err := os.ReadFile(dispatchWorkflowPath(repoRoot))
	if err != nil {
		t.Fatal(err)
	}
	if string(firstWorkflowBytes) != string(secondWorkflowBytes) {
		t.Error("workflow file changed after a failed-identity-gate rerun — it should be untouched")
	}

	// The ledger must still describe the FIRST run's render (sigma,
	// engine) — not be rewritten with the second run's acme identity for
	// a file the second run never actually produced.
	secondData, err := os.ReadFile(filepath.Join(repoRoot, ".cn", "repo.state.json"))
	if err != nil {
		t.Fatalf("read repo.state.json after second run: %v", err)
	}
	var secondState repostate.RepoState
	if err := json.Unmarshal(secondData, &secondState); err != nil {
		t.Fatal(err)
	}
	secondWF := secondState.FindManagedFile(".github/workflows/cnos-cds-dispatch.yml")
	if secondWF == nil {
		t.Fatal("ledger lost its workflow entry entirely after the second run")
	}
	if secondWF.Agent != "sigma" || secondWF.Tier != repostate.TierEngine {
		t.Errorf("ledger workflow entry corrupted by the second (failed, non-rendering) run: got agent=%q tier=%q, want agent=sigma tier=engine (the first run's actual render)", secondWF.Agent, secondWF.Tier)
	}
}
