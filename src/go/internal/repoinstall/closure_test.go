package repoinstall

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// The shipped cell specs and the default installed package set must agree:
// a fill LOADS skill bodies from the installed root, so any skill a committed
// cell names has to be installable by a normal `cn repo install`. Without
// this the corpus proves a hand-assembled fixture the product cannot produce.
func TestDefaultPackagesCoverShippedCells(t *testing.T) {
	root := repoRoot(t)
	specs, err := filepath.Glob(filepath.Join(root, "schemas", "*", "fixtures", "*-cell-spec.json"))
	if err != nil || len(specs) == 0 {
		t.Fatalf("no shipped cell specs found: %v", err)
	}
	checked := 0
	for _, path := range specs {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var spec struct {
			Params map[string]struct {
				Domain []string `json:"domain"`
			} `json:"params"`
			Alpha map[string]json.RawMessage `json:"alpha"`
		}
		if err := json.Unmarshal(data, &spec); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		raw, ok := spec.Alpha["skills"]
		if !ok {
			continue // this cell's fill takes no skills
		}
		var refs []string
		if err := json.Unmarshal(raw, &refs); err != nil {
			t.Fatalf("%s: skills: %v", path, err)
		}
		// A hole stands for every value its declared domain allows, so the
		// closure must cover all of them.
		var concrete []string
		for _, ref := range refs {
			if name, isHole := strings.CutPrefix(ref, "$"); isHole {
				concrete = append(concrete, spec.Params[name].Domain...)
				continue
			}
			concrete = append(concrete, ref)
		}
		for _, ref := range concrete {
			pkg, sub, ok := strings.Cut(ref, ":")
			if !ok {
				t.Errorf("%s: skill ref %q is not <package>:<path>", filepath.Base(path), ref)
				continue
			}
			if !slices.Contains(DefaultPackages, pkg) {
				t.Errorf("%s names skill %q, but %q is not in DefaultPackages %v — "+
					"a normally installed hub could not construct this cell",
					filepath.Base(path), ref, pkg, DefaultPackages)
			}
			// The skill must also actually exist to be installed.
			if _, err := os.Stat(filepath.Join(root, "src", "packages", pkg, "skills", filepath.FromSlash(sub), "SKILL.md")); err != nil {
				t.Errorf("%s names skill %q, which is not present in src/packages: %v",
					filepath.Base(path), ref, err)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("no skill refs were checked — the assertion would pass vacuously")
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(wd, "..", "..", "..", "..")
}
