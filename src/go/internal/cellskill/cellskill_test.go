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

// One root, one lookup: there is no second place to look, so a ref that is
// not installed fails instead of resolving from somewhere else. The previous
// installed-then-source fallback is gone — a second place to look would be a
// second authority.
func TestSingleRootHasNoFallback(t *testing.T) {
	installed, elsewhere := t.TempDir(), t.TempDir()
	writeSkill(t, installed, "cnos.eng:eng/go", "installed body\n")
	writeSkill(t, elsewhere, "cnos.eng:eng/test", "somewhere else\n")

	tree := Tree{Root: installed}
	got, err := tree.Load("cnos.eng:eng/go")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got.Body != "installed body\n" {
		t.Fatalf("wrong body: %q", got.Body)
	}
	if _, err := tree.Load("cnos.eng:eng/test"); err == nil {
		t.Fatal("a skill outside the single root must not resolve")
	}
}

// The digest is the identity, and it is the digest OF THE BODY that was
// loaded — a changed body is a changed skill.
func TestDigestTracksBody(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "p:a", "one\n")
	tree := Tree{Root: root}
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
	tree := Tree{Root: root}
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
	got, err := LoadAll(Tree{Root: root}, []string{"p:c", "p:a", "p:b"})
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
