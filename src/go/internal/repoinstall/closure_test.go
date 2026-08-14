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
// the runtime LOADS skill bodies from the installed root, so any skill a
// committed cell names has to be installable by a normal `cn repo install`.
// Without this the corpus proves a hand-assembled fixture the product cannot
// produce.
//
// The refs are read from the CELL's methodology bundle, which is where a cell
// declares them. While this read `alpha.skills` it went on passing after that
// key was deleted — every spec took the "this cell names no skills" branch, and
// only the `checked == 0` guard below stood between that and a green vacuous
// assertion.
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
			Methodology *struct {
				Skills []string `json:"skills"`
			} `json:"methodology"`
		}
		if err := json.Unmarshal(data, &spec); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		if spec.Methodology == nil {
			continue // this cell declares no methodology
		}
		refs := spec.Methodology.Skills
		// A hole stands for every value its declared domain allows, so the
		// closure must cover all of them. A hole whose parameter is missing or
		// whose domain is empty expands to NOTHING — the skill would silently
		// go unchecked while the surrounding fixed refs keep `checked > 0` and
		// the assertion green. That is the vacuity this guard closes
		// (Pi #55 C2), so an unbounded hole is a failure, not a skip.
		var concrete []string
		for _, ref := range refs {
			name, isHole := strings.CutPrefix(ref, "$")
			if !isHole {
				concrete = append(concrete, ref)
				continue
			}
			p, declared := spec.Params[name]
			if !declared {
				t.Errorf("%s: skill hole %q references undeclared parameter %q",
					filepath.Base(path), ref, name)
				continue
			}
			if len(p.Domain) == 0 {
				t.Errorf("%s: skill hole %q has no domain, so it expands to no skill "+
					"and this closure check would not cover it", filepath.Base(path), ref)
				continue
			}
			concrete = append(concrete, p.Domain...)
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
