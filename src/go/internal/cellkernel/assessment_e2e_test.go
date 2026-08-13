// This file is the ONE external test package in this directory. It has to be:
// the episode it drives is assembled from the CDS fills, and those import the
// kernel, so an internal test could not name them.
//
// It belongs here rather than beside the fills because what it proves is the
// KERNEL'S half of the claim — that a whole run's assessment lands inside the
// one scope-lift digest, and that a one-byte change to any unit result breaks
// the closure's self-verification.
package cellkernel_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellfills"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
	"github.com/usurobor/cnos/src/go/internal/cellspec"
	"github.com/usurobor/cnos/src/go/internal/cellwork"
)

// AC5. A whole cell runs end to end behind the deterministic providers: alpha
// produces a real change in a real worktree, the runtime MEASURES it,
// reconstructs the candidate from (subject, matter), runs the closed checker
// against that reconstruction, assesses every catalogue unit, and closes — with
// the assessment inside the record and therefore inside the one digest.
//
// Nothing is rented. Both seats are the deterministic providers, so the episode
// is honestly `mechanical` and its outcome is reproducible.
func TestACellProducesMeasuresReconstructsChecksAssessesAndCloses(t *testing.T) {
	repo, head := gitRepo(t)
	reg := cellfills.With(skillTree(t, "cnos.eng:eng/go"))

	spec, err := cellspec.Parse(cellJSON())
	if err != nil {
		t.Fatalf("parse the cell: %v", err)
	}
	resolved, err := spec.Resolve(nil)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	subject, err := json.Marshal(cellwork.Subject{Kind: cellwork.SubjectKind, Repo: repo, BaseSHA: head})
	if err != nil {
		t.Fatal(err)
	}
	kspec, meta, err := resolved.Build(context.Background(), reg, cellspec.Binding{
		Issue:   issueJSON(t),
		Subject: subject,
	})
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if meta.ExecutionMode != cellkernel.ModeMechanical {
		t.Fatalf("nothing was rented, so the episode is mechanical, got %q", meta.ExecutionMode)
	}

	cl, err := cellkernel.RunEpisode(context.Background(), kspec, meta)
	if err != nil {
		t.Fatalf("episode: %v", err)
	}
	if err := cellkernel.VerifyClosure(kspec.Contract, meta, cl); err != nil {
		t.Fatalf("the closure does not re-verify: %v", err)
	}

	rec := cl.Receipt.Record
	// alpha really produced and the runtime really measured: the matter is the
	// diff of the change the deterministic coder made, not its account of it.
	if !strings.Contains(rec.Matter.Data, "diff --git") || !strings.Contains(rec.Matter.Data, "CELL-FAKE-CHANGE.txt") {
		t.Fatalf("the record carries no measured change: %q", rec.Matter.Data)
	}
	// ...and the assessment covers the issue's criteria plus the two check
	// units, in order.
	want := []string{"AC1", "AC2", "check:matter-nonempty", "check:project-verify"}
	got := make([]string, 0, len(rec.Review.Assessment))
	for _, u := range rec.Review.Assessment {
		got = append(got, u.Unit)
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("assessment covers %v, want %v", got, want)
	}
	// The matter unit was decided by the runtime from a diff that really exists;
	// the acceptance units were not decided by anyone, because nothing was
	// rented — and `unverified` is what that honestly is.
	byUnit := map[string]cellkernel.UnitResult{}
	for _, u := range rec.Review.Assessment {
		byUnit[u.Unit] = u
	}
	if d := byUnit["check:matter-nonempty"].Disposition; d != cellkernel.DispositionPass {
		t.Errorf("the matter unit = %q, want pass", d)
	}
	if d := byUnit["AC1"].Disposition; d != cellkernel.DispositionUnverified {
		t.Errorf("AC1 = %q; a seat that rented nothing decided nothing", d)
	}
	for _, u := range rec.Review.Assessment {
		if u.Disposition != cellkernel.DispositionPass && strings.TrimSpace(u.Reason) == "" {
			t.Errorf("unit %q is %q with no reason", u.Unit, u.Disposition)
		}
	}
	if cl.Status != cellkernel.NeedsRepair {
		t.Fatalf("status = %q, want needs_repair: not every unit passed", cl.Status)
	}
	// V named each non-pass unit, so the verdict is re-derivable from the
	// receipt rather than taken from the seat's boolean.
	if len(cl.Verdict.Failures) < len(rec.Review.Assessment) {
		t.Fatalf("V reported %d failures for %d non-passing units: %+v",
			len(cl.Verdict.Failures), len(rec.Review.Assessment), cl.Verdict.Failures)
	}

	// THE DIGEST COVERS THE ASSESSMENT. One byte changed in one unit result and
	// the closure no longer re-verifies — which is what makes the assessment a
	// verdict a reader can check rather than a decoration beside one.
	for _, mutate := range []struct {
		what string
		fn   func(*cellkernel.UnitResult)
	}{
		{"the disposition", func(u *cellkernel.UnitResult) { u.Disposition = cellkernel.DispositionPass }},
		{"the reason", func(u *cellkernel.UnitResult) { u.Reason += "." }},
		{"the unit id", func(u *cellkernel.UnitResult) { u.Unit += "x" }},
	} {
		t.Run("digest covers "+mutate.what, func(t *testing.T) {
			tampered := cl
			units := append([]cellkernel.UnitResult(nil), rec.Review.Assessment...)
			mutate.fn(&units[0])
			tampered.Receipt.Record.Review.Assessment = units
			err := cellkernel.VerifyClosure(kspec.Contract, meta, tampered)
			if err == nil || !strings.Contains(err.Error(), "scope-lift digest does not recompute") {
				t.Fatalf("changing %s left the closure verifying: %v", mutate.what, err)
			}
		})
	}
}

// cellJSON is the cell under test: a CDS patch alpha and a CDS assessing beta,
// both behind the deterministic providers, over one methodology bundle.
func cellJSON() []byte {
	return []byte(`{
		"version": "cnos.cellspec.v0",
		"contract": {
			"id": "cds-code",
			"goal": "make the change the issue describes",
			"required_evidence": [{"id": "diff", "kind": "diff", "producer": "alpha"}]
		},
		"protocol_id": "cnos.cdd.cds.receipt.v1",
		"methodology": {"kind": "skills.methodology.v0", "skills": ["cnos.eng:eng/go"]},
		"alpha": {"fill": "cds.patch", "cognition": {"provider": "fake", "model": ""}},
		"beta": {"fill": "cds.assess", "cognition": {"provider": "fake", "model": ""}}
	}`)
}

func issueJSON(t *testing.T) json.RawMessage {
	t.Helper()
	raw, err := json.Marshal(map[string]any{
		"kind": "cnos.cds.issue.v0",
		"id":   "e2e-issue",
		"problem": map[string]string{
			"exists": "nothing reviews the work", "expected": "every obligation is disposed of",
			"diverges": "no episode can honestly accept",
		},
		"sources": []map[string]string{{"claim": "the assessing fill", "path": "src/go/internal/cdsassess"}},
		"scope":   map[string]any{"in": []string{"the reviewing seat"}, "out": []string{}},
		"acceptance": []map[string]string{
			{"id": "AC1", "statement": "the catalogue is exact", "verification": "go test ./internal/cdsassess/..."},
			{"id": "AC2", "statement": "the checker forces its unit", "verification": "go test ./internal/cdsassess/..."},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func gitRepo(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	run := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
		return strings.TrimSpace(string(out))
	}
	run("init", "-q", "-b", "main")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("base\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	run("add", "-A")
	run("commit", "-qm", "base")
	return dir, run("rev-parse", "HEAD")
}

func skillTree(t *testing.T, refs ...string) cellskill.Tree {
	t.Helper()
	root := t.TempDir()
	for _, ref := range refs {
		pkg, path, _ := strings.Cut(ref, ":")
		dir := filepath.Join(root, pkg, "skills", filepath.FromSlash(path))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"),
			[]byte(fmt.Sprintf("# body of %s\n", ref)), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return cellskill.Tree{Root: root}
}
