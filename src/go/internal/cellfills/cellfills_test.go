package cellfills_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellfills"
	"github.com/usurobor/cnos/src/go/internal/cellmethod"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
)

var testSkills = []string{"cnos.eng:eng/code", "cnos.eng:eng/test", "cnos.eng:eng/go", "cnos.eng:eng/write-functional"}

// There is no repository here, and its absence is the assertion: constructing
// a cds.patch alpha touches no git at all now. The repository the episode acts
// on comes from the run's pinned contract subject, at Produce.

type resolvedDecl struct {
	Methodology struct {
		Role   string `json:"role"`
		SHA256 string `json:"sha256"`
	} `json:"methodology"`
}

func methodologyDecl(refs ...string) []byte {
	d, err := json.Marshal(map[string]any{"kind": cellmethod.Kind, "skills": refs})
	if err != nil {
		panic(err)
	}
	return d
}

// D2: the same canonical skill bodies and digests load from an INSTALLED hub
// tree while the process runs somewhere else entirely — skill authority is
// the hub, never the working directory.
//
// The resolver now sits on the REGISTRY and feeds the cell's one methodology
// bundle, so this is what it proves: what `Assemble` wires up resolves against
// `<hub>/.cn/vendor/packages`, and the digest the seat records is that
// bundle's.
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
	if reg.Skills == nil {
		t.Fatal("the assembled registry carries no skill authority")
	}
	bundle, bodies, err := cellmethod.Load(reg.Skills, methodologyDecl(testSkills...))
	if err != nil {
		t.Fatalf("installed-hub methodology load from a foreign cwd failed: %v", err)
	}

	// Same canonical identities and digests as a direct load against the tree.
	want, err := cellskill.LoadAll(cellskill.Tree{Root: installed}, testSkills)
	if err != nil {
		t.Fatal(err)
	}
	for i, s := range want {
		if bundle.Skills[i].Ref != s.Ref || bundle.Skills[i].SHA256 != s.SHA256 {
			t.Fatalf("skill %d: got %+v, want %s/%s", i, bundle.Skills[i], s.Ref, s.SHA256)
		}
	}

	// ...and the seat this binary registers records THAT bundle. Without this
	// the registry could resolve the hub correctly and still hand the fill
	// something else.
	decl := json.RawMessage(`{"fill":"cds.patch","cognition":{"provider":"fake","model":""}}`)
	got, err := reg.ConstructAlpha(context.Background(), decl, cellmethod.Constructive(bundle, bodies))
	if err != nil {
		t.Fatalf("construction from the installed hub failed: %v", err)
	}
	var rd resolvedDecl
	if err := json.Unmarshal(got.Decl, &rd); err != nil {
		t.Fatal(err)
	}
	if rd.Methodology.SHA256 != bundle.SHA256 || rd.Methodology.Role != string(cellmethod.RoleConstructive) {
		t.Fatalf("the seat recorded %+v, want the hub bundle %s", rd.Methodology, bundle.SHA256)
	}

	// An uninstalled skill fails closed — there is no fallback search.
	if _, _, err := cellmethod.Load(reg.Skills, methodologyDecl("cnos.eng:eng/nope")); err == nil {
		t.Fatal("an uninstalled skill must fail the methodology load, not fall back")
	}
}
