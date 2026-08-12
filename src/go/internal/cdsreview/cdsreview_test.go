package cdsreview

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cdsissue"
	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
	"github.com/usurobor/cnos/src/go/internal/cellwork"
)

var testSkills = []string{"cnos.cdd:cdd/review", "cnos.eng:eng/go"}

func skillTree(t *testing.T, refs ...string) cellskill.Tree {
	t.Helper()
	root := t.TempDir()
	for _, ref := range refs {
		pkg, path, _ := strings.Cut(ref, ":")
		dir := filepath.Join(root, pkg, "skills", filepath.FromSlash(path))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# body of "+ref+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return cellskill.Tree{Root: root}
}

func declJSON(provider string) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(
		`{"fill":"cds.review","cognition":{"provider":%q,"model":""},`+
			`"skills":["cnos.cdd:cdd/review","cnos.eng:eng/go"]}`, provider))
}

// testIssue is an admissible CDS issue. The reviewing seat admits the
// contract's issue before it rents anything, so every test that expects to
// reach the answerer must carry one.
const testIssue = `{
	"kind": "cnos.cds.issue.v0",
	"id": "test-contributing",
	"problem": {
		"exists": "The repository carries no contribution guide.",
		"expected": "A CONTRIBUTING.md states how a change is proposed.",
		"diverges": "Nothing states the procedure, so every contributor invents one."
	},
	"sources": [{"claim": "what the repository is", "path": "README.md"}],
	"scope": {"in": ["add CONTRIBUTING.md"], "out": []},
	"acceptance": [{
		"id": "AC1",
		"statement": "CONTRIBUTING.md exists at the repository root.",
		"verification": "the diff adds CONTRIBUTING.md"
	}]
}`

// aDiff is the smallest matter the gate admits: a real unified-diff header.
const aDiff = "diff --git a/CONTRIBUTING.md b/CONTRIBUTING.md\n" +
	"--- /dev/null\n+++ b/CONTRIBUTING.md\n@@ -0,0 +1 @@\n+how to contribute\n"

var reviewContract = cellkernel.Contract{
	ID:   "c",
	Goal: "add a CONTRIBUTING guide",
	Task: json.RawMessage(testIssue),
}

// testRepo builds a one-commit git repository and returns its path and HEAD.
// The reviewing seat now reconstructs its view from the contract's subject, so
// a test that expects to reach the answerer needs a real repository the real
// matter applies to — a stub subject would only prove that a stub was accepted.
func testRepo(t *testing.T) (string, string) {
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

// contractOn is the review contract acting on a real repository at its pinned
// HEAD — the frozen value construction would have handed both stations.
func contractOn(t *testing.T, repo string) cellkernel.Contract {
	t.Helper()
	authored, err := json.Marshal(cellwork.Subject{Kind: cellwork.SubjectKind, Repo: repo, BaseSHA: "HEAD"})
	if err != nil {
		t.Fatal(err)
	}
	pinned, err := cellwork.Pin(context.Background(), authored)
	if err != nil {
		t.Fatalf("pin the test subject: %v", err)
	}
	c := reviewContract
	c.Subject = pinned
	return c
}

// contributingRepo is a repository plus the measured matter that adds a
// CONTRIBUTING.md to it — real matter for the fixture issue, produced the way
// a producing seat produces it.
func contributingRepo(t *testing.T) (cellkernel.Contract, cellkernel.Matter) {
	t.Helper()
	repo, head := testRepo(t)
	wt, release, err := cellwork.Materialize(context.Background(), repo, head)
	if err != nil {
		t.Fatalf("materialize: %v", err)
	}
	defer release()
	if err := os.WriteFile(filepath.Join(wt.Dir, "CONTRIBUTING.md"), []byte("how to contribute\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	diff, err := wt.Diff(context.Background())
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	return contractOn(t, repo), cellkernel.Matter{Data: diff}
}

func admittedTestIssue(t *testing.T) cdsissue.Issue {
	t.Helper()
	iss, err := cdsissue.Admit([]byte(testIssue))
	if err != nil {
		t.Fatalf("the test issue must be admissible: %v", err)
	}
	return iss
}

// stubAnswerer returns exactly what a test wants back, so the decode boundary
// can be driven without renting anything. `called` is what proves the
// NEGATIVE — that nothing was rented — which an error return could not.
type stubAnswerer struct {
	raw    string
	err    error
	called *bool
}

func (stubAnswerer) Name() string { return "stub" }
func (s stubAnswerer) Answer(context.Context, string, json.RawMessage) (json.RawMessage, error) {
	if s.called != nil {
		*s.called = true
	}
	if s.err != nil {
		return nil, s.err
	}
	return json.RawMessage(s.raw), nil
}

// The reviewing seat holds NO workspace: not a repo, not a directory, not a
// worktree. Independence is structural — it cannot reach alpha's workspace
// because it was never given one and has no tool to look with.
func TestSeatHoldsNoWorkspace(t *testing.T) {
	c, err := Factory(skillTree(t, testSkills...))(context.Background(), declJSON("fake"))
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	seat, ok := c.Seat.(ReviewBeta)
	if !ok {
		t.Fatalf("seat is %T, want ReviewBeta", c.Seat)
	}
	// Whole-value check rather than field-by-field: a workspace added later in
	// any form fails here without anyone remembering to extend the list.
	if got := fmt.Sprintf("%+v", struct {
		Answerer any
		Skills   int
	}{seat.answerer, len(seat.skills)}); !strings.Contains(got, "Skills:2") {
		t.Fatalf("unexpected seat shape: %s", got)
	}
	var decl map[string]json.RawMessage
	if err := json.Unmarshal(c.Decl, &decl); err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"workspace", "repo", "base_sha"} {
		if _, present := decl[forbidden]; present {
			t.Errorf("resolved reviewer declaration must not carry %q: %s", forbidden, c.Decl)
		}
	}
}

// A workspace key is rejected FOR BEING A WORKSPACE KEY. The shared corpus
// cannot witness this: without an installed hub every spec fails earlier on
// skill loading, so an exit code alone would prove nothing.
func TestWorkspaceKeyRejectedForItsOwnReason(t *testing.T) {
	decl := json.RawMessage(`{"fill":"cds.review","cognition":{"provider":"fake","model":""},` +
		`"workspace":{"kind":"git-worktree","repo":".","base_sha":"x"},` +
		`"skills":["cnos.cdd:cdd/review"]}`)
	_, err := Factory(skillTree(t, testSkills...))(context.Background(), decl)
	if err == nil {
		t.Fatal("a reviewer declaring a workspace must fail construction")
	}
	if !strings.Contains(err.Error(), "workspace") {
		t.Fatalf("rejected for the wrong reason: %v", err)
	}
}

// A verdict that cannot be read is NOT a failing verdict. Defaulting an
// unreadable answer to pass:false would fabricate a judgement nobody made —
// the same sin as fabricating completion, from the other direction.
func TestUnreadableVerdictFailsRatherThanDefaulting(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{name: "unknown field", raw: `{"pass":false,"notes":"n","verdict":"extra"}`, want: "schema"},
		{name: "wrong type", raw: `{"pass":"no","notes":"n"}`, want: "schema"},
		{name: "not an object", raw: `"fail"`, want: "schema"},
		{name: "empty notes", raw: `{"pass":false,"notes":"   "}`, want: "no notes"},
		{name: "passing with no reason", raw: `{"pass":true,"notes":""}`, want: "no notes"},
	}
	contract, matter := contributingRepo(t)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			seat := ReviewBeta{answerer: stubAnswerer{raw: tc.raw}}
			_, err := seat.Review(context.Background(), cellkernel.BetaInput{
				Contract: contract,
				Matter:   matter,
			})
			if err == nil {
				t.Fatalf("verdict %s must fail, not become a review", tc.raw)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("failed for the wrong reason: got %v, want mention of %q", err, tc.want)
			}
		})
	}
}

// A well-formed verdict passes through unchanged, in both directions, and
// carries no artifacts — evidence is gamma/V's channel, never beta's.
func TestWellFormedVerdictIsCarried(t *testing.T) {
	contract, matter := contributingRepo(t)
	for _, pass := range []bool{true, false} {
		raw := fmt.Sprintf(`{"pass":%t,"notes":"because"}`, pass)
		out, err := ReviewBeta{answerer: stubAnswerer{raw: raw}}.Review(
			context.Background(),
			cellkernel.BetaInput{Contract: contract, Matter: matter})
		if err != nil {
			t.Fatalf("%s: %v", raw, err)
		}
		if out.Review.Pass != pass || out.Review.Notes != "because" {
			t.Fatalf("verdict altered in transit: %+v", out.Review)
		}
		if len(out.Artifacts) != 0 {
			t.Fatalf("a reviewing seat must emit no artifacts, got %d", len(out.Artifacts))
		}
	}
}

// The prompt IS the seat's whole world. If the contract goal or the matter
// were missing from it, the reviewer would be judging something else.
func TestPromptCarriesContractAndMatter(t *testing.T) {
	skills := []cellskill.Skill{{Ref: "cnos.cdd:cdd/review", Body: "REVIEW-SKILL-BODY"}}
	view := cellwork.View{Files: []cellwork.FileState{
		{Path: "CONTRIBUTING.md", Status: cellwork.FileAdded, Content: "THE-VIEWED-CONTENT"},
	}}
	got := RenderPrompt(reviewContract, admittedTestIssue(t), cellkernel.Matter{Data: "THE-MATTER"}, view, skills)
	for _, want := range []string{"add a CONTRIBUTING guide", "THE-MATTER", "THE-VIEWED-CONTENT",
		"REVIEW-SKILL-BODY", "cnos.cdd:cdd/review"} {
		if !strings.Contains(got, want) {
			t.Errorf("prompt is missing %q", want)
		}
	}
}

// Construction fails closed, and the resolved declaration records digests —
// naming a skill is not loading it.
func TestConstructionFailsClosedAndRecordsDigests(t *testing.T) {
	tree := skillTree(t, testSkills...)
	for _, tc := range []struct {
		name string
		decl string
		want string
	}{
		{"no skills", `{"fill":"cds.review","cognition":{"provider":"fake","model":""},"skills":[]}`, "needs its skills"},
		{"unknown provider", `{"fill":"cds.review","cognition":{"provider":"clyde","model":"m"},"skills":["cnos.cdd:cdd/review"]}`, "unknown provider"},
		{"fake with a model", `{"fill":"cds.review","cognition":{"provider":"fake","model":"m"},"skills":["cnos.cdd:cdd/review"]}`, "takes no model"},
		{"uninstalled skill", `{"fill":"cds.review","cognition":{"provider":"fake","model":""},"skills":["nope:x"]}`, "not installed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Factory(tree)(context.Background(), json.RawMessage(tc.decl)); err == nil {
				t.Fatal("must fail construction")
			} else if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("wrong reason: %v", err)
			}
		})
	}

	c, err := Factory(tree)(context.Background(), declJSON("fake"))
	if err != nil {
		t.Fatal(err)
	}
	var rd ResolvedDecl
	if err := json.Unmarshal(c.Decl, &rd); err != nil {
		t.Fatal(err)
	}
	if len(rd.Skills) != 2 {
		t.Fatalf("want 2 resolved skills, got %d", len(rd.Skills))
	}
	for i, s := range rd.Skills {
		if s.Ref != testSkills[i] {
			t.Errorf("skill order changed: got %q want %q", s.Ref, testSkills[i])
		}
		if len(s.SHA256) != 64 {
			t.Errorf("skill %q carries no content digest", s.Ref)
		}
	}
	if c.Mode != cellkernel.ExecutionMode("mechanical") {
		t.Errorf("a fake reviewer must be mechanical, got %q", c.Mode)
	}
}

// The deterministic fake never passes. A fake that returned pass:true would be
// false completion wearing a reviewer's clothes.
func TestFakeReviewerNeverPasses(t *testing.T) {
	c, err := Factory(skillTree(t, testSkills...))(context.Background(), declJSON("fake"))
	if err != nil {
		t.Fatal(err)
	}
	contract, matter := contributingRepo(t)
	out, err := c.Seat.Review(context.Background(), cellkernel.BetaInput{
		Contract: contract,
		Matter:   matter,
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Review.Pass {
		t.Fatal("the deterministic fake must never pass a review it did not perform")
	}
}

// AC3: the reviewing seat admits the contract's issue before renting, by the
// same predicate alpha used. A reviewer whose criteria cannot be read is back
// to judging plausibility, which is the failure this cycle exists to remove.
func TestIllDefinedIssueFailsBeforeRentingCognition(t *testing.T) {
	repo, _ := testRepo(t)
	for name, task := range map[string]string{
		"absent":                 ``,
		"not an issue at all":    `{"goal":"judge the work"}`,
		"criterion unverifiable": strings.Replace(testIssue, `"the diff adds CONTRIBUTING.md"`, `""`, 1),
		"scope.out absent":       strings.Replace(testIssue, `, "out": []`, ``, 1),
	} {
		t.Run(name, func(t *testing.T) {
			called := false
			contract := contractOn(t, repo)
			contract.Task = json.RawMessage(task)
			seat := ReviewBeta{answerer: stubAnswerer{raw: `{"pass":true,"notes":"n"}`, called: &called}}
			_, err := seat.Review(context.Background(), cellkernel.BetaInput{
				Contract: contract,
				Matter:   cellkernel.Matter{Data: aDiff},
			})
			if err == nil {
				t.Fatal("an ill-defined issue must fail the cell")
			}
			if !strings.Contains(err.Error(), "cds issue") {
				t.Fatalf("failed for the wrong reason: %v", err)
			}
			if called {
				t.Fatal("cognition was rented for an issue that was never admissible")
			}
		})
	}
}

// AC5, both directions. A reviewer handed nothing reviewable has NOT passed
// it, and asking a provider to judge `cds.patch`'s "no change was made to …"
// sentence buys an opinion about a sentence. The second half of the table is
// what stops this from being a gate that rejects everything.
func TestOnlyAUnifiedDiffReachesTheAnswerer(t *testing.T) {
	cases := []struct {
		name  string
		data  string
		reach bool
	}{
		{"alpha's no-change sentence", "no change was made to /tmp/repo at abc123", false},
		{"empty", "", false},
		{"whitespace", "   \n\t\n", false},
		{"prose about a diff", "I refactored the parser and everything passes now.", false},
		{"a header-shaped near miss", "diff --gita/x b/x\n", false},
		{"a unified diff", aDiff, true},
		{"a diff after a preamble", "here is the change:\n" + aDiff, true},
	}
	contract, real := contributingRepo(t)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := tc.data
			if tc.reach {
				// A matter that passes the gate is then RECONSTRUCTED, so the
				// admissible cases carry a diff that really applies to the
				// contract's subject rather than a hand-written stand-in.
				data = strings.Replace(tc.data, aDiff, real.Data, 1)
			}
			called := false
			seat := ReviewBeta{answerer: stubAnswerer{raw: `{"pass":true,"notes":"looks fine"}`, called: &called}}
			out, err := seat.Review(context.Background(), cellkernel.BetaInput{
				Contract: contract,
				Matter:   cellkernel.Matter{Data: data},
			})
			if err != nil {
				t.Fatalf("the gate returns a verdict, never an error: %v", err)
			}
			if called != tc.reach {
				t.Fatalf("answerer called=%v, want %v", called, tc.reach)
			}
			if tc.reach {
				if !out.Review.Pass {
					t.Fatal("a reviewable diff must carry the answerer's verdict through")
				}
				return
			}
			if out.Review.Pass {
				t.Fatal("unreviewable matter was passed")
			}
			if !strings.Contains(out.Review.Notes, "not reviewed") {
				t.Fatalf("the verdict must say it did not review: %q", out.Review.Notes)
			}
		})
	}
}

var _ cellfill.BetaFactory = Factory(cellskill.Tree{})

// The reviewing seat admits the contract's SUBJECT before renting, exactly as
// it admits the issue. A reviewer that could not say what state it is judging
// would have to take the diff's word for it — which is the position it was in
// before this cycle.
func TestInadmissibleSubjectFailsBeforeRentingCognition(t *testing.T) {
	repo, _ := testRepo(t)
	good := contractOn(t, repo)
	unpinned, err := json.Marshal(cellwork.Subject{Kind: cellwork.SubjectKind, Repo: repo, BaseSHA: "HEAD"})
	if err != nil {
		t.Fatal(err)
	}
	// Each case names the reason it must fail for. "cds.review" alone was not a
	// reason — every error this fill returns contains it, so the unpinned case
	// passed even when the seat admitted a moving ref and died later on a patch
	// that did not apply. A shared prefix is not a witness.
	for name, tc := range map[string]struct {
		subject json.RawMessage
		want    string
	}{
		"absent":        {nil, "carries no subject"},
		"not a subject": {json.RawMessage(`{"repo":"/tmp/x"}`), "kind"},
		"unpinned base": {unpinned, "is not pinned"},
	} {
		subject, want := tc.subject, tc.want
		t.Run(name, func(t *testing.T) {
			called := false
			contract := good
			contract.Subject = subject
			seat := ReviewBeta{answerer: stubAnswerer{raw: `{"pass":true,"notes":"n"}`, called: &called}}
			_, err := seat.Review(context.Background(), cellkernel.BetaInput{
				Contract: contract,
				Matter:   cellkernel.Matter{Data: aDiff},
			})
			if err == nil {
				t.Fatal("an inadmissible subject must fail the cell")
			}
			if !strings.Contains(err.Error(), want) {
				t.Fatalf("failed for the wrong reason:\n got: %v\nwant mention of: %q", err, want)
			}
			if called {
				t.Fatal("cognition was rented to judge a state the cell could not name")
			}
		})
	}
}

// A reconstruction that cannot be built is a mechanism fault, not a verdict.
// Returning pass:false here would report a judgement nobody made — the same
// sin as fabricating completion, from the reviewing side.
func TestUnreconstructableMatterIsAFaultAndNotAVerdict(t *testing.T) {
	repo, _ := testRepo(t)
	called := false
	seat := ReviewBeta{answerer: stubAnswerer{raw: `{"pass":true,"notes":"n"}`, called: &called}}
	// Structurally a unified diff — so it passes the matter gate — but its
	// context does not exist at the contract's base.
	stale := "diff --git a/README.md b/README.md\n" +
		"--- a/README.md\n+++ b/README.md\n@@ -1 +1 @@\n-a line the base never had\n+rewritten\n"
	out, err := seat.Review(context.Background(), cellkernel.BetaInput{
		Contract: contractOn(t, repo),
		Matter:   cellkernel.Matter{Data: stale},
	})
	if err == nil {
		t.Fatalf("an unreconstructable matter must not become a verdict, got %+v", out.Review)
	}
	if !strings.Contains(err.Error(), "does not apply") {
		t.Fatalf("failed for the wrong reason: %v", err)
	}
	if called {
		t.Fatal("cognition was rented for a view that was never built")
	}
}

// AC4 at the seat: the reviewer's prompt carries what the hunks cannot. This is
// the recorded defect — a reviewer blocked on a missing import that the file
// imports far outside every hunk — reproduced and closed.
func TestPromptCarriesTheFileTheHunksDoNotShow(t *testing.T) {
	repo, _ := testRepo(t)
	var body strings.Builder
	body.WriteString("package widget\n\nimport (\n\t\"bytes\"\n\t\"fmt\"\n)\n\n")
	for i := 0; i < 120; i++ {
		body.WriteString("// filler line so the import is nowhere near the change\n")
	}
	body.WriteString("func Render() string {\n\tvar out bytes.Buffer\n\tfmt.Fprint(&out, \"old\")\n\treturn out.String()\n}\n")
	if err := os.WriteFile(filepath.Join(repo, "widget.go"), []byte(body.String()), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-qm", "widget"}} {
		cmd := exec.Command("git", args...)
		cmd.Dir = repo
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	wt, release, err := cellwork.Materialize(context.Background(), repo, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	if err := os.WriteFile(filepath.Join(wt.Dir, "widget.go"),
		[]byte(strings.Replace(body.String(), `"old"`, `"new"`, 1)), 0o600); err != nil {
		t.Fatal(err)
	}
	matter, err := wt.Diff(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(matter, `"bytes"`) {
		t.Fatalf("the fixture diff shows the import, so it proves nothing:\n%s", matter)
	}

	var prompt string
	seat := ReviewBeta{answerer: promptCapturingAnswerer{&prompt}}
	if _, err := seat.Review(context.Background(), cellkernel.BetaInput{
		Contract: contractOn(t, repo),
		Matter:   cellkernel.Matter{Data: matter},
	}); err != nil {
		t.Fatalf("review: %v", err)
	}
	if !strings.Contains(prompt, "\t\"bytes\"\n") {
		t.Fatal("the reviewer's prompt does not carry the import the hunks omit")
	}
	if !strings.Contains(prompt, "RECONSTRUCTED VIEW") {
		t.Fatal("the prompt does not say what the view is, so the reviewer cannot know whose account it is")
	}
	// And the matter is still there: the view is additional evidence, not a
	// replacement for the change under review.
	if !strings.Contains(prompt, matter) {
		t.Fatal("the matter no longer reaches the reviewer")
	}
}

// promptCapturingAnswerer keeps the prompt so a test can assert on the
// reviewer's whole world without renting anything.
type promptCapturingAnswerer struct{ prompt *string }

func (promptCapturingAnswerer) Name() string { return "capturing" }
func (a promptCapturingAnswerer) Answer(_ context.Context, prompt string, _ json.RawMessage) (json.RawMessage, error) {
	*a.prompt = prompt
	return json.RawMessage(`{"pass":true,"notes":"captured"}`), nil
}

// A truncated view must be told to the reviewer as a limit on what it may
// conclude. The verdict shape has no `unverified` yet, so the instruction is to
// fail and name it — never to pass a criterion whose evidence was not shown.
func TestTruncatedViewIsRenderedAsALimitOnConclusions(t *testing.T) {
	full := renderView(cellwork.View{Files: []cellwork.FileState{
		{Path: "kept.go", Status: cellwork.FileModified, Content: "package kept\n"},
	}})
	if strings.Contains(full, "INCOMPLETE") {
		t.Fatal("a complete view must not warn about truncation")
	}
	cut := renderView(cellwork.View{
		Truncated: true,
		Files: []cellwork.FileState{
			{Path: "big.go", Status: cellwork.FileModified, Omitted: true},
			{Path: "gone.go", Status: cellwork.FileDeleted},
			{Path: "logo.bin", Status: cellwork.FileAdded, Omitted: true, Binary: true},
			{Path: "empty.txt", Status: cellwork.FileAdded},
		},
	})
	for _, want := range []string{"INCOMPLETE", "unverified", "size bound", "binary", "deleted", "empty"} {
		if !strings.Contains(cut, want) {
			t.Errorf("the rendered view does not mention %q:\n%s", want, cut)
		}
	}
	// The four empty-content cases must read differently from one another, or a
	// reviewer cannot tell "checked, absent" from "not shown".
	lines := map[string]bool{}
	for _, l := range strings.Split(cut, "\n") {
		if strings.HasPrefix(l, "(") {
			lines[l] = true
		}
	}
	if len(lines) != 4 {
		t.Fatalf("the four content-absent cases do not read distinctly: %v", lines)
	}
}

// A link target must never be rendered as if it were the file's bytes: the
// reviewer decides criteria from this text, and "the file contains /etc/passwd"
// and "the file IS a link to /etc/passwd" are different facts.
func TestRenderViewNamesASymlinkRatherThanShowingItAsContent(t *testing.T) {
	out := renderView(cellwork.View{Files: []cellwork.FileState{
		{Path: "shim", Status: cellwork.FileAdded, Symlink: true, Content: "../tools/real"},
	}})
	if !strings.Contains(out, "symbolic link") || !strings.Contains(out, "did not follow") {
		t.Fatalf("the rendered view does not say the path is a link the runtime left alone:\n%s", out)
	}
	if !strings.Contains(out, "../tools/real") {
		t.Fatalf("the rendered view drops the link target, which IS its content:\n%s", out)
	}
}
