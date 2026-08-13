package cellmethod

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cellskill"
)

var refs = []string{"cnos.eng:eng/code", "cnos.eng:eng/test", "cnos.eng:eng/go", "cnos.eng:eng/write-functional"}

// tree writes an installed package tree whose SKILL.md bodies are distinct, so
// a projection that dropped or reordered one is visible rather than absorbed
// into identical text.
func tree(t *testing.T, refs ...string) string {
	t.Helper()
	root := t.TempDir()
	for i, ref := range refs {
		writeSkill(t, root, ref, fmt.Sprintf("# body of %s\nobligation %d\n", ref, i))
	}
	return root
}

func writeSkill(t *testing.T, root, ref, body string) {
	t.Helper()
	pkg, path, _ := strings.Cut(ref, ":")
	dir := filepath.Join(root, pkg, "skills", filepath.FromSlash(path))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func declFor(refs ...string) []byte {
	d, err := json.Marshal(map[string]any{"kind": Kind, "skills": refs})
	if err != nil {
		panic(err)
	}
	return d
}

func load(t *testing.T, root string, refs ...string) (Bundle, []cellskill.Skill) {
	t.Helper()
	b, bodies, err := Load(cellskill.Tree{Root: root}, declFor(refs...))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	return b, bodies
}

// AC1, this package's half: the digest identifies the ordered (ref,
// body-digest) list and nothing else. Every clause is a separate way the
// digest could be wrong, and each is asserted by moving exactly one thing.
func TestTheDigestCoversTheOrderedRefAndBodyDigestList(t *testing.T) {
	root := tree(t, refs...)
	base, bodies := load(t, root, refs...)

	if len(base.Skills) != len(refs) {
		t.Fatalf("bundle carries %d skills, want %d", len(base.Skills), len(refs))
	}
	for i, r := range base.Skills {
		if r.Ref != refs[i] {
			t.Fatalf("skill %d is %q, want %q — order is meaning", i, r.Ref, refs[i])
		}
		if r.SHA256 != bodies[i].SHA256 || len(r.SHA256) != 64 {
			t.Fatalf("skill %q carries no body digest: %q", r.Ref, r.SHA256)
		}
	}

	// Stable: the same tree loaded twice is the same methodology.
	if again, _ := load(t, root, refs...); again.SHA256 != base.SHA256 {
		t.Fatal("the digest is not stable across two loads of the same tree")
	}

	// Content, not location: an identical tree at a different path is the same
	// methodology. Without this the digest could be covering filenames.
	if elsewhere, _ := load(t, tree(t, refs...), refs...); elsewhere.SHA256 != base.SHA256 {
		t.Fatal("the digest changed with the installed root: it is not covering content")
	}

	// ONE BYTE in one body.
	moved := tree(t, refs...)
	writeSkill(t, moved, refs[2], "# body of "+refs[2]+"\nobligation 2\n"+"!")
	if changed, _ := load(t, moved, refs...); changed.SHA256 == base.SHA256 {
		t.Fatal("a one-byte change to a skill body did not change the methodology digest")
	}

	// Reordering. Later skills refine earlier ones, so a different order is a
	// different methodology even though the set is identical.
	swapped := []string{refs[1], refs[0], refs[2], refs[3]}
	if reordered, _ := load(t, root, swapped...); reordered.SHA256 == base.SHA256 {
		t.Fatal("reordering the bundle did not change its digest")
	}

	// Dropping one.
	if shorter, _ := load(t, root, refs[:3]...); shorter.SHA256 == base.SHA256 {
		t.Fatal("dropping a skill did not change the methodology digest")
	}
}

func TestLoadFailsClosed(t *testing.T) {
	root := tree(t, refs...)
	bad := map[string]string{
		"not an object":    `["cnos.eng:eng/go"]`,
		"unknown key":      `{"kind":"` + Kind + `","skills":["cnos.eng:eng/go"],"extra":1}`,
		"mixed-case key":   `{"Kind":"` + Kind + `","skills":["cnos.eng:eng/go"]}`,
		"wrong kind":       `{"kind":"skills.methodology.v1","skills":["cnos.eng:eng/go"]}`,
		"no skills":        `{"kind":"` + Kind + `","skills":[]}`,
		"uninstalled":      `{"kind":"` + Kind + `","skills":["cnos.eng:eng/nope"]}`,
		"trailing data":    `{"kind":"` + Kind + `","skills":["cnos.eng:eng/go"]} {}`,
		"malformed suffix": `{"kind":"` + Kind + `","skills":["cnos.eng:eng/go"]} @`,
	}
	for name, decl := range bad {
		t.Run(name, func(t *testing.T) {
			if _, _, err := Load(cellskill.Tree{Root: root}, []byte(decl)); err == nil {
				t.Fatalf("%s must not load", name)
			}
		})
	}
	// A nil resolver is a composition fault, not an empty methodology: it must
	// say so rather than panic or return a bundle covering nothing.
	if _, _, err := Load(nil, declFor(refs...)); err == nil ||
		!strings.Contains(err.Error(), "no skill resolver") {
		t.Fatalf("a nil resolver must be named: %v", err)
	}
}

// AC2. Both projections carry every loaded body, in order, byte-identical, and
// differ ONLY in the leading role wrapper.
//
// What this proves is narrow and worth stating: the two texts agree because
// neither selects anything, so "every obligation appears in both" is trivially
// true. It becomes a real property when a projection can drop or transform an
// obligation. Until then this test guards the identity claim the package
// header makes, not a preservation property it does not yet have.
func TestBothProjectionsCarryEveryBodyInOrder(t *testing.T) {
	root := tree(t, refs...)
	b, bodies := load(t, root, refs...)
	c, a := Constructive(b, bodies), Adversarial(b, bodies)

	if c.Role != RoleConstructive || a.Role != RoleAdversarial {
		t.Fatalf("roles are wrong: %q / %q", c.Role, a.Role)
	}
	if c.SHA256 != b.SHA256 || a.SHA256 != b.SHA256 {
		t.Fatal("a view must carry the digest of the bundle it came from")
	}
	if c.Empty() || a.Empty() {
		t.Fatal("a projection of a loaded bundle is not empty")
	}

	if !strings.HasPrefix(c.Text, constructiveWrapper) || !strings.HasPrefix(a.Text, adversarialWrapper) {
		t.Fatal("a projection must lead with its role wrapper")
	}
	if constructiveWrapper == adversarialWrapper {
		t.Fatal("the two wrappers are identical: the projections are indistinguishable")
	}
	cRest := strings.TrimPrefix(c.Text, constructiveWrapper)
	aRest := strings.TrimPrefix(a.Text, adversarialWrapper)
	if cRest != aRest {
		t.Fatalf("the projections differ below the role wrapper:\n--- constructive ---\n%s\n--- adversarial ---\n%s", cRest, aRest)
	}

	// ...and what they share is every body, verbatim, in declaration order. The
	// equality above would hold for two empty texts.
	at := -1
	for _, s := range bodies {
		i := strings.Index(cRest, s.Body)
		if i < 0 {
			t.Fatalf("body of %q is not carried verbatim", s.Ref)
		}
		if i <= at {
			t.Fatalf("body of %q is out of declaration order", s.Ref)
		}
		at = i
	}
}

// AC5, this package's half: the adversarial projection has no production
// caller, and this test fails when it gains one — at which point the sentence
// in Adversarial's doc comment is what must be rewritten.
//
// The scan is proven able to find a caller before it is trusted to report
// none: Constructive HAS one (cellspec builds the producing seat's view), and
// a scan that could not see it would report both as uncalled.
func TestAdversarialHasNoProductionCaller(t *testing.T) {
	calledFrom := productionReferences(t, "cellmethod.Adversarial(")
	if len(calledFrom) != 0 {
		t.Fatalf("the adversarial projection now has production callers %v — "+
			"it is no longer consumer-less, and Adversarial's doc comment says it is", calledFrom)
	}
	if len(productionReferences(t, "cellmethod.Constructive(")) == 0 {
		t.Fatal("the scan found no caller of Constructive either: it is not looking at production sources")
	}
}

// productionReferences returns the non-test Go files under the module that
// mention `needle`, excluding this package's own directory (a package refers to
// its own identifiers unqualified, so a self-reference would be a false
// positive of a different kind).
func productionReferences(t *testing.T, needle string) []string {
	t.Helper()
	root, err := filepath.Abs("../..") // src/go, the module root
	if err != nil {
		t.Fatal(err)
	}
	self, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	var found []string
	seen := 0
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if filepath.Dir(path) == self {
			return nil
		}
		seen++
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(data), needle) {
			rel, _ := filepath.Rel(root, path)
			found = append(found, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	// A walk that read nothing would report "no callers" for every needle.
	if seen < 50 {
		t.Fatalf("the walk saw only %d production files under %s", seen, root)
	}
	return found
}
