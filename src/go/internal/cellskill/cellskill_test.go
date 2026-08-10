package cellskill

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

// Roots resolve like $PATH: first match wins, so an installed package shadows
// the source tree rather than the other way round.
func TestRootsResolveInOrder(t *testing.T) {
	installed, source := t.TempDir(), t.TempDir()
	writeSkill(t, installed, "cnos.eng:eng/go", "installed body\n")
	writeSkill(t, source, "cnos.eng:eng/go", "source body\n")
	writeSkill(t, source, "cnos.eng:eng/test", "only in source\n")

	tree := Tree{Roots: []string{installed, source}}
	got, err := tree.Load("cnos.eng:eng/go")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got.Body != "installed body\n" {
		t.Fatalf("first root did not win: %q", got.Body)
	}
	// A ref only the later root has still resolves.
	if _, err := tree.Load("cnos.eng:eng/test"); err != nil {
		t.Fatalf("later root not searched: %v", err)
	}
}

// The digest is the identity, and it is the digest OF THE BODY that was
// loaded — a changed body is a changed skill.
func TestDigestTracksBody(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "p:a", "one\n")
	tree := Tree{Roots: []string{root}}
	first, err := tree.Load("p:a")
	if err != nil {
		t.Fatal(err)
	}
	writeSkill(t, root, "p:a", "two\n")
	second, err := tree.Load("p:a")
	if err != nil {
		t.Fatal(err)
	}
	if first.SHA256 == second.SHA256 {
		t.Fatal("digest did not follow the body")
	}
	if len(first.SHA256) != 64 {
		t.Fatalf("digest is not a sha256: %q", first.SHA256)
	}
}

func TestLoadFailsClosed(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "p:a", "body\n")
	tree := Tree{Roots: []string{root}}
	bad := []string{
		"no-colon",
		":path",
		"pkg:",
		"../escape:a",        // traversal in the package
		"p:../../etc/passwd", // traversal in the path
		"p:/etc",             // absolute path
		"p/sub:a",            // separator in the package
		"p:missing",          // simply not installed
	}
	for _, ref := range bad {
		t.Run(ref, func(t *testing.T) {
			if _, err := tree.Load(ref); err == nil {
				t.Fatalf("ref %q must fail", ref)
			}
		})
	}
	if _, err := (Tree{}).Load("p:a"); err == nil {
		t.Fatal("no configured roots must fail closed, not silently resolve")
	}
}

// LoadAll preserves order — later skills refine earlier ones and the record
// keeps that sequence.
func TestLoadAllPreservesOrder(t *testing.T) {
	root := t.TempDir()
	for _, ref := range []string{"p:a", "p:b", "p:c"} {
		writeSkill(t, root, ref, "body of "+ref+"\n")
	}
	got, err := LoadAll(Tree{Roots: []string{root}}, []string{"p:c", "p:a", "p:b"})
	if err != nil {
		t.Fatalf("load all: %v", err)
	}
	var refs []string
	for _, s := range got {
		refs = append(refs, s.Ref)
	}
	if strings.Join(refs, ",") != "p:c,p:a,p:b" {
		t.Fatalf("order not preserved: %v", refs)
	}
}
