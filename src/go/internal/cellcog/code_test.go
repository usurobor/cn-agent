package cellcog

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellkernel"
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

var codeContract = cellkernel.Contract{
	ID:   "cell-code",
	Goal: "add a NOTES file",
	RequiredEvidence: []cellkernel.RequiredRef{
		{ID: DiffArtifactID, Kind: DiffArtifactKind, Producer: cellkernel.RoleAlpha},
	},
}

// writeCoder creates a file; idleCoder honors the prompt and does nothing.
type writeCoder struct{ name, content string }

func (c writeCoder) Name() string { return "write" }
func (c writeCoder) Work(_ context.Context, dir, _ string) error {
	return os.WriteFile(filepath.Join(dir, c.name), []byte(c.content), 0o600)
}

type idleCoder struct{}

func (idleCoder) Name() string                               { return "idle" }
func (idleCoder) Work(context.Context, string, string) error { return nil }

type failCoder struct{}

func (failCoder) Name() string                               { return "fail" }
func (failCoder) Work(context.Context, string, string) error { return errors.New("backend died") }

// The change is measured from the worktree: the diff carries what was written.
func TestCodeAlphaMeasuresTheChange(t *testing.T) {
	repo, head := testRepo(t)
	a := CodeAlpha{Coder: writeCoder{name: "NOTES.md", content: "hello\n"}, Repo: repo, BaseRef: "HEAD"}
	out, err := a.Produce(context.Background(), cellkernel.AlphaInput{Contract: codeContract})
	if err != nil {
		t.Fatalf("produce: %v", err)
	}
	if !strings.Contains(out.Matter.Data, "NOTES.md") || !strings.Contains(out.Matter.Data, "+hello") {
		t.Fatalf("matter is not the measured diff: %q", out.Matter.Data)
	}
	var diff, base string
	for _, c := range out.Artifacts {
		switch c.ID {
		case DiffArtifactID:
			diff = c.Text
		case BaseArtifactID:
			base = c.Text
		}
	}
	if diff == "" {
		t.Fatal("no diff artifact")
	}
	if base != head {
		t.Fatalf("base_sha artifact = %q, want the resolved HEAD %q", base, head)
	}
}

// The load-bearing property: a seat that changed nothing produces no diff, so
// a contract requiring one closes needs_repair. False completion — the running
// system's #514/#516 scar, where wakes "repaired 0 of 41 items" and reported
// success — is unrepresentable rather than caught late.
func TestIdleCoderCannotFalselyComplete(t *testing.T) {
	repo, _ := testRepo(t)
	s := cellkernel.Spec{
		Contract: codeContract,
		Alpha:    CodeAlpha{Coder: idleCoder{}, Repo: repo, BaseRef: "HEAD"},
		Beta:     MatterBeta{},
	}
	meta := testMeta()
	cl, err := cellkernel.RunEpisode(context.Background(), s, meta)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if cl.Status != cellkernel.NeedsRepair {
		t.Fatalf("an unchanged worktree closed %q, want needs_repair", cl.Status)
	}
	for _, c := range cl.Receipt.Record.Alpha.Artifacts {
		if c.ID == DiffArtifactID {
			t.Fatal("a diff artifact was manufactured from no change")
		}
	}
}

// The worktree is disposable and the repository is left untouched: the seat
// never edits the caller's checkout, and no worktree survives the episode.
func TestWorktreeIsDisposableAndRepoUntouched(t *testing.T) {
	repo, head := testRepo(t)
	a := CodeAlpha{Coder: writeCoder{name: "NOTES.md", content: "hello\n"}, Repo: repo, BaseRef: "HEAD"}
	if _, err := a.Produce(context.Background(), cellkernel.AlphaInput{Contract: codeContract}); err != nil {
		t.Fatalf("produce: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repo, "NOTES.md")); !os.IsNotExist(err) {
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

func TestCodeAlphaFailsClosed(t *testing.T) {
	repo, _ := testRepo(t)
	in := cellkernel.AlphaInput{Contract: codeContract}
	if _, err := (CodeAlpha{Repo: repo, BaseRef: "HEAD"}).Produce(context.Background(), in); !errors.Is(err, ErrNoProvider) {
		t.Fatalf("nil coder: want ErrNoProvider, got %v", err)
	}
	if _, err := (CodeAlpha{Coder: failCoder{}, Repo: repo, BaseRef: "HEAD"}).Produce(context.Background(), in); err == nil {
		t.Fatal("a failing coder must reach the kernel as a seat error")
	}
	if _, err := (CodeAlpha{Coder: idleCoder{}, Repo: repo, BaseRef: "no-such-rev"}).Produce(context.Background(), in); err == nil {
		t.Fatal("an unresolvable base must fail before any work")
	}
	if _, err := (CodeAlpha{Coder: idleCoder{}, Repo: t.TempDir(), BaseRef: "HEAD"}).Produce(context.Background(), in); err == nil {
		t.Fatal("a non-repository must fail closed")
	}
}

// The whole code episode closes and verifies against the caller's own
// contract and metadata.
func TestCodeEpisodeClosesAndVerifies(t *testing.T) {
	repo, _ := testRepo(t)
	s := cellkernel.Spec{
		Contract: codeContract,
		Alpha:    CodeAlpha{Coder: FakeCoder{}, Repo: repo, BaseRef: "HEAD"},
		Beta:     MatterBeta{},
	}
	meta := testMeta()
	cl, err := cellkernel.RunEpisode(context.Background(), s, meta)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if cl.Status != cellkernel.Accepted {
		t.Fatalf("status = %q, want accepted (%+v)", cl.Status, cl.Verdict)
	}
	if err := cellkernel.VerifyClosure(codeContract, meta, cl); err != nil {
		t.Fatalf("closure must verify: %v", err)
	}
}
