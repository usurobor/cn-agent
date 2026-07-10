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

	// The ledger write happens after runDispatchCds returns nil — with no
	// git remote, label-doctor errors, so repo.state.json is never
	// reached this run. Directly exercise writeRepoState the same way Run
	// would, to assert the workflow record shape in isolation.
	manifest := pkg.Manifest{Schema: manifestSchema, Profile: manifestProfile, Packages: []pkg.ManifestDep{{Name: "cnos.core", Version: "9.9.9"}, {Name: "cnos.cds", Version: "9.9.9"}}}
	opts := Options{RepoRoot: repoRoot, Dispatch: "cds", Engine: true, Stdout: stdout, Stderr: stderr}
	if err := writeRepoState(opts, manifest, DefaultRepo, defaultDownloadURL, "9.9.9"); err != nil {
		t.Fatalf("writeRepoState: %v", err)
	}

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
