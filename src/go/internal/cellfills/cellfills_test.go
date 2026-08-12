package cellfills_test

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellfills"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
)

var testSkills = []string{"cnos.eng:eng/code", "cnos.eng:eng/test", "cnos.eng:eng/go", "cnos.eng:eng/write-functional"}

func testRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
	run("init", "-q", "-b", "main")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("base\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	run("add", "-A")
	run("commit", "-qm", "base")
	return dir
}

type resolvedDecl struct {
	Skills []struct {
		Ref    string `json:"ref"`
		SHA256 string `json:"sha256"`
	} `json:"skills"`
}

// D2: the same canonical skill bodies and digests load from an INSTALLED hub
// tree while the process runs somewhere else entirely — skill authority is
// the hub, never the working directory.
func TestInstalledHubLoadsFromForeignCwd(t *testing.T) {
	repo := testRepo(t)
	hub := t.TempDir()
	installed := cellfills.InstalledPackages(hub)
	for _, ref := range testSkills {
		pkg, path, _ := strings.Cut(ref, ":")
		dir := filepath.Join(installed, pkg, "skills", filepath.FromSlash(path))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# body of "+ref+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	// Run from a directory unrelated to both the hub and the repository.
	foreign := t.TempDir()
	restore, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(foreign); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(restore) })

	reg := cellfills.Assemble(hub)
	decl := json.RawMessage(`{"fill":"cds.patch","cognition":{"provider":"fake","model":""},` +
		`"skills":["cnos.eng:eng/code","cnos.eng:eng/test","cnos.eng:eng/go","cnos.eng:eng/write-functional"]}`)
	got, err := reg.ConstructAlpha(context.Background(), decl)
	if err != nil {
		t.Fatalf("installed-hub construction from a foreign cwd failed: %v", err)
	}
	var rd resolvedDecl
	if err := json.Unmarshal(got.Decl, &rd); err != nil {
		t.Fatal(err)
	}
	// Same canonical identities and digests as the direct-tree construction.
	want, err := cellskill.LoadAll(cellskill.Tree{Root: installed}, testSkills)
	if err != nil {
		t.Fatal(err)
	}
	for i, s := range want {
		if rd.Skills[i].Ref != s.Ref || rd.Skills[i].SHA256 != s.SHA256 {
			t.Fatalf("skill %d: got %+v, want %s/%s", i, rd.Skills[i], s.Ref, s.SHA256)
		}
	}

	// An uninstalled skill fails closed — there is no fallback search.
	missing := strings.Replace(string(decl), `"cnos.eng:eng/go"`, `"cnos.eng:eng/nope"`, 1)
	if _, err := reg.ConstructAlpha(context.Background(), json.RawMessage(missing)); err == nil {
		t.Fatal("an uninstalled skill must fail construction, not fall back")
	}

	// The assembled registry also carries the ONE subject adapter this binary
	// ships. Without it a spec declaring a subject cannot be built at all, so a
	// registry that forgot to wire it would fail every CDS cell at construction
	// rather than silently running one unpinned.
	if reg.PinSubject == nil {
		t.Fatal("the assembled registry wires no subject adapter")
	}
	authored, err := json.Marshal(map[string]string{
		"kind": "git.snapshot/0.1", "repo": repo, "base_sha": "HEAD",
	})
	if err != nil {
		t.Fatal(err)
	}
	pinned, err := reg.PinSubject(context.Background(), authored)
	if err != nil {
		t.Fatalf("the wired subject adapter must pin a real repository: %v", err)
	}
	var snapshot struct {
		BaseSHA string `json:"base_sha"`
	}
	if err := json.Unmarshal(pinned, &snapshot); err != nil {
		t.Fatal(err)
	}
	if len(snapshot.BaseSHA) != 40 {
		t.Fatalf("the wired adapter did not pin the base: %s", pinned)
	}
}
