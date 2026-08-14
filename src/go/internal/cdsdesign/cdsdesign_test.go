package cdsdesign

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cdsissue"
)

// corpusDir is the SAME directory scripts/cell-schema-check.sh vets with
// `cue vet … -d '#CDSDesign'`. Two authorities, one corpus: a document the two
// disagree about fails a gate here rather than surfacing as a run input that
// vets and will not run (or runs and will not vet).
const corpusDir = "../../../../schemas/cds/fixtures/design"

// Every negative is invalid for exactly ONE reason, and the expected substring
// pins WHICH rule fired. A fixture that broke two rules could not show that,
// and a test asserting only "some error" would pass on the wrong one.
var negatives = map[string]string{
	"design-bad-kind.json":           "kind must be",
	"design-blank-approach.json":     "approach is required",
	"design-no-invariants.json":      "invariants is required",
	"design-blank-invariant.json":    "invariants[1] is blank",
	"design-no-impact.json":          "impact is required",
	"design-blank-surface.json":      "impact[1] names no surface",
	"design-impact-without-why.json": "states no reason",
	"design-unknown-key.json":        `unknown key "alternatives"`,
	"design-mixed-case-key.json":     `unknown key "Kind"`,
}

func corpusFiles(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		t.Fatalf("read design corpus: %v", err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".json") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names
}

// The corpus must actually contain both classes. A missing or renamed
// directory would otherwise make every assertion below vacuous — the same
// failure class files_exist guards in the shell corpus.
func TestCorpusIsPopulatedAndFullyClassified(t *testing.T) {
	names := corpusFiles(t)
	valid := 0
	for _, n := range names {
		switch {
		case strings.HasPrefix(n, "valid-"):
			valid++
		case negatives[n] != "":
		default:
			t.Errorf("corpus file %q is neither a valid- fixture nor a named negative; "+
				"an unclassified fixture is vetted by CUE and by nothing here", n)
		}
	}
	if valid == 0 {
		t.Fatal("the design corpus has no positive fixture")
	}
	if len(negatives) != len(names)-valid {
		t.Fatalf("negatives table lists %d files, corpus has %d non-positive files",
			len(negatives), len(names)-valid)
	}
}

func TestCorpusAdmission(t *testing.T) {
	for _, name := range corpusFiles(t) {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(corpusDir, name))
			if err != nil {
				t.Fatal(err)
			}
			_, admitErr := Admit(data)
			if strings.HasPrefix(name, "valid-") {
				if admitErr != nil {
					t.Fatalf("valid design rejected: %v", admitErr)
				}
				return
			}
			want := negatives[name]
			if want == "" {
				t.Fatalf("no expected reason declared for %q", name)
			}
			if admitErr == nil {
				t.Fatal("must be rejected")
			}
			if !strings.Contains(admitErr.Error(), want) {
				t.Fatalf("rejected for the wrong reason:\n got: %v\nwant mention of: %q", admitErr, want)
			}
		})
	}
}

func TestAbsentDesignIsRejected(t *testing.T) {
	for _, raw := range [][]byte{nil, {}, []byte("null")} {
		if _, err := Admit(raw); err == nil {
			t.Fatalf("Admit(%q) must fail: a run input with no design has no design", raw)
		}
	}
}

// The blankness predicate is cdsissue's, and this is the statement that it is
// the SAME one rather than an equivalent-looking copy. A design carrying only
// a NO-BREAK SPACE is blank here exactly because Blank says so there; if this
// package ever grew its own predicate, the assertion below would still pass
// while the two could drift, so it compares the two functions directly rather
// than comparing two outcomes.
func TestBlanknessIsTheIssuePredicate(t *testing.T) {
	// One authored-looking but blank field, refused.
	blank := []byte(`{"kind":"` + Kind + `","approach":" ",` +
		`"invariants":["i"],"impact":[{"surface":"s","why":"w"}]}`)
	if _, err := Admit(blank); err == nil {
		t.Fatal("a NO-BREAK-SPACE-only approach must be refused")
	} else if !strings.Contains(err.Error(), "approach is required") {
		t.Fatalf("refused for the wrong reason: %v", err)
	}
	if !cdsissue.Blank(" ") {
		t.Fatal("cdsissue.Blank no longer covers NO-BREAK SPACE, so the case above " +
			"is refused by something else in this package")
	}
}
