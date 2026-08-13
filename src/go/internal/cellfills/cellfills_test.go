package cellfills_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellfills"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
)

var testSkills = []string{"cnos.eng:eng/code", "cnos.eng:eng/test", "cnos.eng:eng/go", "cnos.eng:eng/write-functional"}

// There is no repository here, and its absence is the assertion: constructing
// a cds.patch alpha touches no git at all now. The repository the episode acts
// on comes from the run's pinned contract subject, at Produce.

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
}
