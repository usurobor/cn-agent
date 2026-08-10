package cdspatch

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
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

var patchContract = cellkernel.Contract{
	ID:   "cds-code",
	Goal: "add a NOTES file",
	RequiredEvidence: []cellkernel.RequiredRef{
		{ID: DiffArtifactID, Kind: DiffArtifactKind, Producer: cellkernel.RoleAlpha},
	},
}

func declJSON(repo, provider string) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{
		"fill": "cds.patch",
		"cognition": {"provider": %q, "model": ""},
		"workspace": {"kind": "git-worktree", "repo": %q, "base_sha": "HEAD"},
		"skills": ["cnos.eng:eng/code", "cnos.eng:eng/test", "cnos.eng:eng/go", "cnos.eng:eng/write-functional"]
	}`, provider, repo))
}

var testSkills = []string{"cnos.eng:eng/code", "cnos.eng:eng/test", "cnos.eng:eng/go", "cnos.eng:eng/write-functional"}

func construct(t *testing.T, repo, provider string) cellfill.ConstructedAlpha {
	t.Helper()
	f := Factory(skillTree(t, testSkills...))
	a, err := f(context.Background(), declJSON(repo, provider))
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	return a
}

// The constructor resolves and LOADS skill bodies: the resolved declaration
// records ordered refs + content digests, and the prompt carries the bodies.
func TestConstructionLoadsSkillsAndCanonicalizes(t *testing.T) {
	repo, _ := testRepo(t)
	a := construct(t, repo, "fake")
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
	prompt := RenderPrompt(patchContract, seat.skills)
	if !strings.Contains(prompt, "# body of cnos.eng:eng/go") {
		t.Fatal("skill BODY was not injected into the prompt — naming is not loading")
	}
	if a.Mode != cellkernel.ModeMechanical {
		t.Fatalf("fake provider mode = %q, want mechanical", a.Mode)
	}
	// "resolved" must mean resolved: the recorded declaration names the exact
	// commit, never the moving ref the caller passed.
	if rd.Workspace.BaseSHA == "HEAD" || len(rd.Workspace.BaseSHA) != 40 {
		t.Fatalf("resolved declaration did not pin the base commit: %q", rd.Workspace.BaseSHA)
	}
}

// The registry canonicalizes what a fill returns, so the recorded bytes — and
// therefore the digest — do not depend on how the fill happened to serialize,
// and constructing twice yields byte-identical declarations.
func TestConstructionIsCanonicalAndStable(t *testing.T) {
	repo, _ := testRepo(t)
	reg := cellfill.CddFills()
	reg.Alpha[Fill] = Factory(skillTree(t, testSkills...))
	first, err := reg.ConstructAlpha(context.Background(), declJSON(repo, "fake"))
	if err != nil {
		t.Fatalf("construct: %v", err)
	}
	second, err := reg.ConstructAlpha(context.Background(), declJSON(repo, "fake"))
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
	repo, _ := testRepo(t)
	f := Factory(skillTree(t, testSkills...))
	bad := map[string]string{
		"unknown key":       `{"fill":"cds.patch","cognition":{"provider":"fake","model":""},"workspace":{"kind":"git-worktree","repo":".","base_sha":"x"},"skills":["cnos.eng:eng/go"],"Extra":1}`,
		"unknown provider":  `{"fill":"cds.patch","cognition":{"provider":"clyde","model":"m"},"workspace":{"kind":"git-worktree","repo":".","base_sha":"x"},"skills":["cnos.eng:eng/go"]}`,
		"modelless claude":  `{"fill":"cds.patch","cognition":{"provider":"claude-cli","model":""},"workspace":{"kind":"git-worktree","repo":".","base_sha":"x"},"skills":["cnos.eng:eng/go"]}`,
		"bad workspace":     `{"fill":"cds.patch","cognition":{"provider":"fake","model":""},"workspace":{"kind":"zip","repo":".","base_sha":"x"},"skills":["cnos.eng:eng/go"]}`,
		"no skills":         `{"fill":"cds.patch","cognition":{"provider":"fake","model":""},"workspace":{"kind":"git-worktree","repo":".","base_sha":"x"},"skills":[]}`,
		"uninstalled skill": `{"fill":"cds.patch","cognition":{"provider":"fake","model":""},"workspace":{"kind":"git-worktree","repo":"` + repo + `","base_sha":"HEAD"},"skills":["cnos.eng:eng/nope"]}`,
	}
	for name, decl := range bad {
		t.Run(name, func(t *testing.T) {
			if _, err := f(context.Background(), json.RawMessage(decl)); err == nil {
				t.Fatalf("%s must fail construction", name)
			}
		})
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
	a := construct(t, repo, "fake").Seat.(PatchAlpha)
	a.coder = idleCoder{}
	out, err := a.Produce(context.Background(), cellkernel.AlphaInput{Contract: patchContract})
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
	a := construct(t, repo, "fake")
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
	cl, err := cellkernel.RunEpisode(context.Background(),
		cellkernel.Spec{Contract: patchContract, Alpha: a.Seat, Beta: b.Seat}, meta)
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
	if err := cellkernel.VerifyClosure(patchContract, meta, cl); err != nil {
		t.Fatalf("closure must verify: %v", err)
	}
}

// The worktree is disposable and the repository is left untouched.
func TestWorktreeIsDisposableAndRepoUntouched(t *testing.T) {
	repo, head := testRepo(t)
	a := construct(t, repo, "fake").Seat.(PatchAlpha)
	if _, err := a.Produce(context.Background(), cellkernel.AlphaInput{Contract: patchContract}); err != nil {
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

func TestPatchAlphaFailsClosed(t *testing.T) {
	repo, _ := testRepo(t)
	in := cellkernel.AlphaInput{Contract: patchContract}
	base := construct(t, repo, "fake").Seat.(PatchAlpha)

	noCoder := base
	noCoder.coder = nil
	if _, err := noCoder.Produce(context.Background(), in); err == nil {
		t.Fatal("nil coder must fail closed")
	}
	badBase := base
	badBase.base = "no-such-rev"
	if _, err := badBase.Produce(context.Background(), in); err == nil {
		t.Fatal("an unresolvable base must fail before any work")
	}
	notRepo := base
	notRepo.repo = t.TempDir()
	if _, err := notRepo.Produce(context.Background(), in); err == nil {
		t.Fatal("a non-repository must fail closed")
	}
}
