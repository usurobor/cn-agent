package labeldoctor

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// chdirFixture creates a git-initialized fixture repo at t.TempDir(),
// seeds src/packages/cnos.core/labels.json with the manifest fixture,
// configures an "origin" remote pointing at usurobor/cnos, chdir's the
// test process into it, and restores the original cwd on cleanup — so
// Run's os.Getwd()-based resolution (matching issues-fsm's
// resolveDefaultTablePath idiom) exercises the exact same path a real
// `cn label-doctor` invocation would.
func chdirFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runGit(t, root, "init", "-q")
	runGit(t, root, "remote", "add", "origin", "https://github.com/usurobor/cnos.git")
	writeFile(t, filepath.Join(root, "src", "packages", "cnos.core", "labels.json"), fixtureManifestJSON)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(cwd); err != nil {
			t.Fatal(err)
		}
	})
	return root
}

// TestRun_DryRun_FailsOnDrift drives the exact `cn label-doctor
// --dry-run` code path (flag parsing -> repo/manifest resolution ->
// Doctor) against injected drift, and asserts it exits nonzero. This is
// the AC5 CI-guard "must fail on injected drift" direction, exercised
// through Run itself (not just the lower-level Doctor/Audit functions),
// since that is the code path `cn label-doctor --dry-run` in CI actually
// runs.
func TestRun_DryRun_FailsOnDrift(t *testing.T) {
	chdirFixture(t)
	store := newFakeLabelStore([]ghLabel{
		{Name: "status:backlog", Color: "ededed", Description: "Well-formed scope but not yet refined."},
		{Name: "status:review", Color: "ededed", Description: "Cell complete; awaiting review."}, // drifted color
		// status:blocked missing
	})
	withFakeGitHub(t, store.handler())

	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), []string{"--token", "tok", "--dry-run"}, nil, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected Run --dry-run to fail on injected drift")
	}
	if !strings.Contains(stderr.String(), "drift detected") {
		t.Errorf("stderr should name the drift, got: %q", stderr.String())
	}
	if store.mutatingCalls() != 0 {
		t.Errorf("--dry-run must make zero mutating calls, got %d", store.mutatingCalls())
	}
}

// TestRun_DryRun_PassesOnClean is AC5's positive direction through the
// same Run() code path: a repo whose live labels already match the
// manifest exits 0.
func TestRun_DryRun_PassesOnClean(t *testing.T) {
	chdirFixture(t)
	store := newFakeLabelStore([]ghLabel{
		{Name: "status:backlog", Color: "ededed", Description: "Well-formed scope but not yet refined."},
		{Name: "status:review", Color: "5319e7", Description: "Cell complete; awaiting review."},
		{Name: "status:blocked", Color: "b60205", Description: "Gated on external input."},
	})
	withFakeGitHub(t, store.handler())

	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), []string{"--token", "tok", "--dry-run"}, nil, &stdout, &stderr)
	if err != nil {
		t.Fatalf("expected Run --dry-run to pass on a clean repo, got: %v\nstderr: %s", err, stderr.String())
	}
	if store.mutatingCalls() != 0 {
		t.Errorf("--dry-run must make zero mutating calls, got %d", store.mutatingCalls())
	}
}

// TestRun_Apply_RepairsLiveState exercises the non-dry-run path through
// Run itself, confirming the CLI wiring (not just Doctor) performs the
// repair.
func TestRun_Apply_RepairsLiveState(t *testing.T) {
	chdirFixture(t)
	store := newFakeLabelStore([]ghLabel{
		{Name: "status:backlog", Color: "ededed", Description: "Well-formed scope but not yet refined."},
	})
	withFakeGitHub(t, store.handler())

	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), []string{"--token", "tok"}, nil, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
	}
	if store.labels["status:review"].Color != "5319e7" {
		t.Errorf("status:review should be repaired to 5319e7, got %q", store.labels["status:review"].Color)
	}
	if _, ok := store.labels["status:blocked"]; !ok {
		t.Error("status:blocked should have been created")
	}
}

func TestRun_Help_ExitsZero(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), []string{"--help"}, nil, &stdout, &stderr)
	// flag.ContinueOnError's --help returns flag.ErrHelp; the CLI wrapper
	// (cmd_label_doctor.go) is responsible for treating that as exit 0,
	// mirroring every other kernel command's --help handling
	// (cmd_repo_install.go checks for --help/-h before parsing at all).
	// This test only proves Run's flag set actually registers a Usage
	// that names every documented flag.
	if err == nil {
		t.Fatal("flag.ContinueOnError should return an error for --help (flag.ErrHelp)")
	}
	if !strings.Contains(stderr.String(), "--repo") || !strings.Contains(stderr.String(), "--dry-run") {
		t.Errorf("usage should document --repo and --dry-run, got: %q", stderr.String())
	}
}

func TestRun_UnknownFlag_Errors(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := Run(context.Background(), []string{"--bogus"}, nil, &stdout, &stderr); err == nil {
		t.Fatal("expected an error for an unrecognized flag")
	}
}
