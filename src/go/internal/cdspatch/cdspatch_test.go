package cdspatch

import (
	"bytes"
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

// testRepo builds a one-commit git repository and returns its path and HEAD.
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

// skillTree writes a minimal installed package tree with the given skills.
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

// testIssue is an admissible CDS issue. Every seat test needs one, because a
// cell whose contract carries no admissible issue no longer reaches a provider
// at all — which is the point of this cycle.
const testIssue = `{
	"kind": "cnos.cds.issue.v0",
	"id": "test-notes",
	"problem": {
		"exists": "The repository records nothing about the change it is under.",
		"expected": "A NOTES file records it.",
		"diverges": "There is no NOTES file."
	},
	"sources": [{"claim": "what the repository is", "path": "README.md"}],
	"scope": {"in": ["add a NOTES file"], "out": ["changing README.md"]},
	"acceptance": [{
		"id": "AC1",
		"statement": "A NOTES file exists at the repository root.",
		"verification": "the diff adds NOTES"
	}]
}`

var patchContract = cellkernel.Contract{
	ID:   "cds-code",
	Goal: "add a NOTES file",
	Task: json.RawMessage(testIssue),
	RequiredEvidence: []cellkernel.RequiredRef{
		{ID: DiffArtifactID, Kind: DiffArtifactKind, Producer: cellkernel.RoleAlpha},
	},
}

func admittedTestIssue(t *testing.T) cdsissue.Issue {
	t.Helper()
	iss, err := cdsissue.Admit([]byte(testIssue))
	if err != nil {
		t.Fatalf("the test issue must be admissible: %v", err)
	}
	return iss
}

func declJSON(provider string) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{
		"fill": "cds.patch",
		"cognition": {"provider": %q, "model": ""},
		"skills": ["cnos.eng:eng/code", "cnos.eng:eng/test", "cnos.eng:eng/go", "cnos.eng:eng/write-functional"]
	}`, provider))
}

// pinnedSubject is what construction hands both stations: the repository and
// the exact commit, resolved once. Every seat test needs one, because a
// producing seat now materializes what the CONTRACT names.
func pinnedSubject(t *testing.T, repo string) json.RawMessage {
	t.Helper()
	authored, err := json.Marshal(cellwork.Subject{Kind: cellwork.SubjectKind, Repo: repo, BaseSHA: "HEAD"})
	if err != nil {
		t.Fatal(err)
	}
	pinned, err := cellwork.Pin(context.Background(), authored)
	if err != nil {
		t.Fatalf("pin the test subject: %v", err)
	}
	return pinned
}

// contractOn is the test contract acting on repo.
func contractOn(t *testing.T, repo string) cellkernel.Contract {
	t.Helper()
	c := patchContract
	c.Subject = pinnedSubject(t, repo)
	return c
}

var testSkills = []string{"cnos.eng:eng/code", "cnos.eng:eng/test", "cnos.eng:eng/go", "cnos.eng:eng/write-functional"}

func construct(t *testing.T, provider string) cellfill.ConstructedAlpha {
	t.Helper()
	f := Factory(skillTree(t, testSkills...))
	a, err := f(context.Background(), declJSON(provider))
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	return a
}

// The constructor resolves and LOADS skill bodies: the resolved declaration
// records ordered refs + content digests, and the prompt carries the bodies.
func TestConstructionLoadsSkillsAndCanonicalizes(t *testing.T) {
	a := construct(t, "fake")
	var rd ResolvedDecl
	if err := json.Unmarshal(a.Decl, &rd); err != nil {
		t.Fatalf("resolved decl: %v", err)
	}
	if len(rd.Skills) != 4 || rd.Skills[0].Ref != "cnos.eng:eng/code" || rd.Skills[2].Ref != "cnos.eng:eng/go" {
		t.Fatalf("resolved skills wrong/unordered: %+v", rd.Skills)
	}
	for _, s := range rd.Skills {
		if len(s.SHA256) != 64 {
			t.Fatalf("skill %q has no content digest", s.Ref)
		}
	}
	seat := a.Seat.(PatchAlpha)
	prompt := RenderPrompt(patchContract, admittedTestIssue(t), seat.skills)
	if !strings.Contains(prompt, "# body of cnos.eng:eng/go") {
		t.Fatal("skill BODY was not injected into the prompt — naming is not loading")
	}
	if a.Mode != cellkernel.ModeMechanical {
		t.Fatalf("fake provider mode = %q, want mechanical", a.Mode)
	}
	// The declaration says nothing about WHERE the work happens. That is
	// contract truth, and a seat-declared copy of it is exactly the second
	// source this change removed.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(a.Decl, &raw); err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"workspace", "repo", "base_sha"} {
		if _, present := raw[forbidden]; present {
			t.Errorf("resolved patch declaration must not carry %q: %s", forbidden, a.Decl)
		}
	}
}

// The registry canonicalizes what a fill returns, so the recorded bytes — and
// therefore the digest — do not depend on how the fill happened to serialize,
// and constructing twice yields byte-identical declarations.
func TestConstructionIsCanonicalAndStable(t *testing.T) {
	reg := cellfill.CddFills()
	reg.Alpha[Fill] = Factory(skillTree(t, testSkills...))
	first, err := reg.ConstructAlpha(context.Background(), declJSON("fake"))
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	second, err := reg.ConstructAlpha(context.Background(), declJSON("fake"))
	if err != nil {
		t.Fatalf("construct again: %v", err)
	}
	if string(first.Decl) != string(second.Decl) {
		t.Fatalf("construction is not stable:\n%s\n%s", first.Decl, second.Decl)
	}
	// Canonical form: keys sorted, no insignificant whitespace.
	if !strings.HasPrefix(string(first.Decl), `{"cognition":`) {
		t.Fatalf("declaration is not canonical: %s", first.Decl)
	}
}

func TestConstructionFailsClosed(t *testing.T) {
	f := Factory(skillTree(t, testSkills...))
	bad := map[string]string{
		"unknown key":       `{"fill":"cds.patch","cognition":{"provider":"fake","model":""},"skills":["cnos.eng:eng/go"],"Extra":1}`,
		"unknown provider":  `{"fill":"cds.patch","cognition":{"provider":"clyde","model":"m"},"skills":["cnos.eng:eng/go"]}`,
		"modelless claude":  `{"fill":"cds.patch","cognition":{"provider":"claude-cli","model":""},"skills":["cnos.eng:eng/go"]}`,
		"no skills":         `{"fill":"cds.patch","cognition":{"provider":"fake","model":""},"skills":[]}`,
		"uninstalled skill": `{"fill":"cds.patch","cognition":{"provider":"fake","model":""},"skills":["cnos.eng:eng/nope"]}`,
		// A workspace is no longer a cds.patch argument at all: where the work
		// happens is the contract's to say, and a seat that could name a second
		// repository could produce a diff against a tree beta never sees.
		"declared workspace": `{"fill":"cds.patch","cognition":{"provider":"fake","model":""},"workspace":{"kind":"git-worktree","repo":".","base_sha":"x"},"skills":["cnos.eng:eng/go"]}`,
	}
	for name, decl := range bad {
		t.Run(name, func(t *testing.T) {
			if _, err := f(context.Background(), json.RawMessage(decl)); err == nil {
				t.Fatalf("%s must fail construction", name)
			}
		})
	}
}

// recordingCoder proves a NEGATIVE: that the provider was never reached. A
// coder that merely returned an error could not distinguish "refused before
// renting" from "rented and then failed".
type recordingCoder struct{ called bool }

func (*recordingCoder) Name() string { return "recording" }
func (c *recordingCoder) Work(context.Context, string, string) error {
	c.called = true
	return nil
}

// AC3: an ill-defined issue fails the cell BEFORE any cognition is rented.
// Cognition is the expensive, non-reproducible resource; spending it on a task
// nobody managed to state is the waste this admission exists to prevent, and
// the matter it would produce is matter beta could not judge.
func TestIllDefinedIssueFailsBeforeRentingCognition(t *testing.T) {
	repo, _ := testRepo(t)
	for name, task := range map[string]string{
		"absent":                 ``,
		"not an issue at all":    `{"goal":"do the thing"}`,
		"criterion unverifiable": strings.Replace(testIssue, `"the diff adds NOTES"`, `""`, 1),
		"unknown key":            strings.Replace(testIssue, `"kind"`, `"knid"`, 1),
	} {
		t.Run(name, func(t *testing.T) {
			coder := &recordingCoder{}
			seat := construct(t, "fake").Seat.(PatchAlpha)
			seat.coder = coder

			contract := contractOn(t, repo)
			contract.Task = json.RawMessage(task)
			_, err := seat.Produce(context.Background(), cellkernel.AlphaInput{Contract: contract})
			if err == nil {
				t.Fatal("an ill-defined issue must fail the cell")
			}
			if !strings.Contains(err.Error(), "cds issue") {
				t.Fatalf("failed for the wrong reason: %v", err)
			}
			if coder.called {
				t.Fatal("cognition was rented for an issue that was never admissible")
			}
		})
	}

	// The guard is not vacuous: the SAME double IS reached when the issue is
	// admissible, so `called` can be true.
	coder := &recordingCoder{}
	seat := construct(t, "fake").Seat.(PatchAlpha)
	seat.coder = coder
	if _, err := seat.Produce(context.Background(), cellkernel.AlphaInput{Contract: contractOn(t, repo)}); err != nil {
		t.Fatalf("an admissible issue must reach the coder: %v", err)
	}
	if !coder.called {
		t.Fatal("the recording coder was never called even for an admissible issue")
	}
}

// The load-bearing property: the change is measured from the worktree. A
// coder that changed nothing produces no diff, so a contract requiring one
// closes needs_repair — false completion (#514/#516) is unrepresentable.
type idleCoder struct{}

func (idleCoder) Name() string                               { return "idle" }
func (idleCoder) Work(context.Context, string, string) error { return nil }

func TestIdleCoderCannotFalselyComplete(t *testing.T) {
	repo, _ := testRepo(t)
	a := construct(t, "fake").Seat.(PatchAlpha)
	a.coder = idleCoder{}
	out, err := a.Produce(context.Background(), cellkernel.AlphaInput{Contract: contractOn(t, repo)})
	if err != nil {
		t.Fatalf("produce: %v", err)
	}
	for _, c := range out.Artifacts {
		if c.ID == DiffArtifactID {
			t.Fatal("a diff artifact was manufactured from no change")
		}
	}
	if !strings.Contains(out.Matter.Data, "no change was made") {
		t.Fatalf("matter must state the truth: %q", out.Matter.Data)
	}
}

// Case-2 honesty end to end: even a REAL measured change closes needs_repair,
// because the mechanical-unmet beta cannot judge the goal and says so. The
// diff is preserved in the closure for the independent review Case 3 adds.
// (The fake coder writes a file unrelated to the NOTES goal — exactly the
// change a passing mechanical beta would have wrongly blessed.)
func TestMeasuredChangeAwaitsIndependentReview(t *testing.T) {
	repo, head := testRepo(t)
	a := construct(t, "fake")
	betas := cellfill.CddFills()
	b, err := betas.ConstructBeta(context.Background(), json.RawMessage(`{"fill":"cdd.mechanical-unmet"}`))
	if err != nil {
		t.Fatalf("beta: %v", err)
	}
	meta := cellkernel.RunMeta{
		ExecutionMode: cellfill.CombineModes(a.Mode, b.Mode),
		ResolvedSpec: cellkernel.ResolvedSpec{
			Version: "cnos.cellspec.v0", DeclaredProtocol: "cnos.cdd.cds.receipt.v1",
			Alpha: a.Decl, Beta: b.Decl,
		},
	}
	contract := contractOn(t, repo)
	cl, err := cellkernel.RunEpisode(context.Background(),
		cellkernel.Spec{Contract: contract, Alpha: a.Seat, Beta: b.Seat}, meta)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if cl.Status != cellkernel.NeedsRepair {
		t.Fatalf("a mechanically unreviewed change closed %q, want needs_repair", cl.Status)
	}
	var diff, base string
	for _, c := range cl.Receipt.Record.Alpha.Artifacts {
		switch c.ID {
		case DiffArtifactID:
			diff = c.Text
		case BaseArtifactID:
			base = c.Text
		}
	}
	if !strings.Contains(diff, "CELL-FAKE-CHANGE.txt") {
		t.Fatal("the measured diff was not preserved in the closure")
	}
	if base != head {
		t.Fatalf("base_sha = %q, want resolved HEAD %q", base, head)
	}
	if err := cellkernel.VerifyClosure(contract, meta, cl); err != nil {
		t.Fatalf("closure must verify: %v", err)
	}
}

// The worktree is disposable and the repository is left untouched.
func TestWorktreeIsDisposableAndRepoUntouched(t *testing.T) {
	repo, head := testRepo(t)
	a := construct(t, "fake").Seat.(PatchAlpha)
	if _, err := a.Produce(context.Background(), cellkernel.AlphaInput{Contract: contractOn(t, repo)}); err != nil {
		t.Fatalf("produce: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repo, "CELL-FAKE-CHANGE.txt")); !os.IsNotExist(err) {
		t.Fatal("the seat wrote into the caller's repository")
	}
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil || strings.TrimSpace(string(out)) != head {
		t.Fatalf("repository HEAD moved: %s (%v)", out, err)
	}
	cmd = exec.Command("git", "worktree", "list")
	cmd.Dir = repo
	list, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(list), "\n") != 1 {
		t.Fatalf("a worktree outlived the episode:\n%s", list)
	}
}

// The seat fails closed on every way its inputs can be wrong, and each cause
// is named. The subject cases are contract-level now: a seat cannot be given a
// bad repository or a moving base by its own declaration, only by a contract
// that carries one.
func TestPatchAlphaFailsClosed(t *testing.T) {
	repo, _ := testRepo(t)
	seat := construct(t, "fake").Seat.(PatchAlpha)

	noCoder := seat
	noCoder.coder = nil
	if _, err := noCoder.Produce(context.Background(), cellkernel.AlphaInput{Contract: contractOn(t, repo)}); err == nil {
		t.Fatal("nil coder must fail closed")
	}

	badSubject := func(t *testing.T, s cellwork.Subject) cellkernel.Contract {
		t.Helper()
		raw, err := json.Marshal(s)
		if err != nil {
			t.Fatal(err)
		}
		c := patchContract
		c.Subject = raw
		return c
	}
	pinned := func(t *testing.T) string {
		t.Helper()
		s, err := cellwork.AdmitSubject(pinnedSubject(t, repo))
		if err != nil {
			t.Fatal(err)
		}
		return s.BaseSHA
	}
	for name, tc := range map[string]struct {
		contract cellkernel.Contract
		want     string
	}{
		"no subject":       {patchContract, "contract carries no subject"},
		"unpinned base":    {badSubject(t, cellwork.Subject{Kind: cellwork.SubjectKind, Repo: repo, BaseSHA: "HEAD"}), "is not pinned"},
		"unknown kind":     {badSubject(t, cellwork.Subject{Kind: "svn.checkout/0.1", Repo: repo, BaseSHA: pinned(t)}), "subject kind must be"},
		"not a repository": {badSubject(t, cellwork.Subject{Kind: cellwork.SubjectKind, Repo: t.TempDir(), BaseSHA: pinned(t)}), "is not a git repository"},
		"absent commit": {badSubject(t, cellwork.Subject{Kind: cellwork.SubjectKind, Repo: repo,
			BaseSHA: "0000000000000000000000000000000000000000"}), "does not resolve"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := seat.Produce(context.Background(), cellkernel.AlphaInput{Contract: tc.contract})
			if err == nil {
				t.Fatal("must fail closed")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("failed for the wrong reason: got %v, want mention of %q", err, tc.want)
			}
		})
	}
}

// A fake ignores the model, so an authored declaration may simply omit the
// key — Go's decoder yields the empty model either way. The authored CUE
// shape admits the same omission; this is the Go half of that parity
// (Pi #57 C1), and the RESOLVED declaration must still record `model: ""`
// so a receipt says what held the seat rather than what the author typed.
func TestFakeMayOmitModelAndStillRecordsIt(t *testing.T) {
	decl := json.RawMessage(`{
		"fill": "cds.patch",
		"cognition": {"provider": "fake"},
		"skills": ["cnos.eng:eng/code", "cnos.eng:eng/test", "cnos.eng:eng/go", "cnos.eng:eng/write-functional"]
	}`)
	a, err := Factory(skillTree(t, testSkills...))(context.Background(), decl)
	if err != nil {
		t.Fatalf("a fake omitting its meaningless model must construct: %v", err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(a.Decl, &raw); err != nil {
		t.Fatal(err)
	}
	var cog map[string]any
	if err := json.Unmarshal(raw["cognition"], &cog); err != nil {
		t.Fatal(err)
	}
	if m, present := cog["model"]; !present || m != "" {
		t.Fatalf("resolved cognition must record model:\"\", got %v (present=%v)", m, present)
	}
}

// subjectWatchBeta records what the reviewing station was handed. A real
// cds.review beta cannot be used here — it lives downstream of this package —
// and what is under test is the kernel-facing half anyway: whether the two
// stations are handed the same value.
type subjectWatchBeta struct{ saw *json.RawMessage }

func (b subjectWatchBeta) Review(_ context.Context, in cellkernel.BetaInput) (cellkernel.BetaOutput, error) {
	*b.saw = in.Contract.Subject
	return cellkernel.BetaOutput{Review: cellkernel.Review{Pass: true, Notes: "watched"}}, nil
}

// AC2: the author names a moving revision, the runtime pins it once, and BOTH
// stations receive those same bytes. The producing seat's own measurement is
// checked against them too — a seat that resolved `HEAD` for itself could
// measure against a commit the reviewer never sees, which is the whole failure
// this contract slot removes.
func TestBothStationsReceiveTheSamePinnedSubject(t *testing.T) {
	repo, head := testRepo(t)
	authored, err := json.Marshal(cellwork.Subject{Kind: cellwork.SubjectKind, Repo: repo, BaseSHA: "HEAD"})
	if err != nil {
		t.Fatal(err)
	}
	pinned, err := cellwork.Pin(context.Background(), authored)
	if err != nil {
		t.Fatalf("pin: %v", err)
	}
	if bytes.Equal(pinned, authored) {
		t.Fatal("the authored subject was already pinned, so this proves nothing")
	}

	contract := patchContract
	contract.Subject = pinned
	a := construct(t, "fake")
	var atBeta json.RawMessage
	meta := cellkernel.RunMeta{
		ExecutionMode: cellfill.CombineModes(a.Mode, cellkernel.ModeMechanical),
		ResolvedSpec: cellkernel.ResolvedSpec{
			Version: "cnos.cellspec.v0", DeclaredProtocol: "cnos.cdd.cds.receipt.v1",
			Alpha: a.Decl, Beta: json.RawMessage(`{"fill":"test.subject-watch"}`),
		},
	}
	cl, err := cellkernel.RunEpisode(context.Background(),
		cellkernel.Spec{Contract: contract, Alpha: a.Seat, Beta: subjectWatchBeta{&atBeta}}, meta)
	if err != nil {
		t.Fatalf("run: %v", err)
	}

	if !bytes.Equal(atBeta, pinned) {
		t.Fatalf("the reviewing station saw %s, want the pinned %s", atBeta, pinned)
	}
	if !bytes.Equal(cl.Receipt.Record.Contract.Subject, pinned) {
		t.Fatalf("the record carries %s, want the pinned %s", cl.Receipt.Record.Contract.Subject, pinned)
	}
	// The producing station's measurement was taken at that same commit: its
	// base_sha artifact is what the runtime actually materialized at.
	var measured string
	for _, art := range cl.Receipt.Record.Alpha.Artifacts {
		if art.ID == BaseArtifactID {
			measured = art.Text
		}
	}
	if measured != head {
		t.Fatalf("alpha measured at %q, want the pinned commit %q", measured, head)
	}
	var got cellwork.Subject
	if err := json.Unmarshal(pinned, &got); err != nil {
		t.Fatal(err)
	}
	if got.BaseSHA != measured {
		t.Fatalf("the contract names %q and the measurement was taken at %q", got.BaseSHA, measured)
	}
}
