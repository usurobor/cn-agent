package cdsissue

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

// corpusDir is the SAME directory scripts/cell-schema-check.sh vets with
// `cue vet … -d '#CDSIssue'`. Two authorities, one corpus: a document the two
// disagree about fails a gate here rather than surfacing as a cell that vets
// and will not run (or runs and will not vet).
const corpusDir = "../../../../schemas/cds/fixtures/issue"

// Every negative is invalid for exactly ONE reason, and the expected substring
// pins WHICH rule fired. A fixture that broke two rules could not show that,
// and a test asserting only "some error" would pass on the wrong one.
var negatives = map[string]string{
	"issue-bad-kind.json":                       "kind must be",
	"issue-empty-id.json":                       "id is required",
	"issue-blank-problem-line.json":             "problem.diverges is required",
	"issue-blank-unicode-whitespace.json":       "problem.exists is required",
	"issue-reserved-acceptance-id.json":         "uses the reserved",
	"issue-no-sources.json":                     "sources is required",
	"issue-source-without-path.json":            "sources[1] needs a claim and a path",
	"issue-empty-scope-in.json":                 "scope.in is required",
	"issue-missing-scope-out.json":              "scope.out must be present",
	"issue-no-acceptance.json":                  "acceptance is required",
	"issue-criterion-without-verification.json": "states no verification route",
	"issue-duplicate-acceptance-id.json":        "duplicate acceptance id",
	"issue-unknown-key.json":                    `unknown key "owner"`,
	"issue-mixed-case-key.json":                 `unknown key "Kind"`,
}

func corpusFiles(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		t.Fatalf("read issue corpus: %v", err)
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
		t.Fatal("the issue corpus has no positive fixture")
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
					t.Fatalf("valid issue rejected: %v", admitErr)
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

// scope.out PRESENT-but-empty and scope.out ABSENT decode to the same Go value,
// so presence has to be read off the raw document. This is the pair that proves
// it: the only difference between the two fixtures is the key.
func TestEmptyScopeOutIsAdmittedAndAbsentIsNot(t *testing.T) {
	present, err := os.ReadFile(filepath.Join(corpusDir, "valid-empty-scope-out.json"))
	if err != nil {
		t.Fatal(err)
	}
	iss, err := Admit(present)
	if err != nil {
		t.Fatalf("an empty out list declares that non-goals were considered: %v", err)
	}
	if len(iss.Scope.Out) != 0 {
		t.Fatalf("fixture is not the empty-out case: %v", iss.Scope.Out)
	}
	absent, err := os.ReadFile(filepath.Join(corpusDir, "issue-missing-scope-out.json"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Admit(absent); err == nil {
		t.Fatal("an absent out key must be rejected; it declares nothing")
	}
}

func TestAbsentTaskIsRejected(t *testing.T) {
	for _, raw := range [][]byte{nil, {}, []byte("null")} {
		if _, err := Admit(raw); err == nil {
			t.Fatalf("Admit(%q) must fail: a cell with no issue has no task", raw)
		}
	}
}

// Render must carry every criterion's verification route. That string is what
// makes beta's job decidable; a renderer that dropped it would leave beta with
// statements and no way to decide them.
func TestRenderCarriesEveryFieldBetaNeeds(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(corpusDir, "valid-issue.json"))
	if err != nil {
		t.Fatal(err)
	}
	iss, err := Admit(data)
	if err != nil {
		t.Fatal(err)
	}
	out := Render(iss)
	want := []string{iss.ID, iss.Problem.Exists, iss.Problem.Expected, iss.Problem.Diverges}
	for _, s := range iss.Sources {
		want = append(want, s.Claim, s.Path)
	}
	want = append(want, iss.Scope.In...)
	want = append(want, iss.Scope.Out...)
	for _, c := range iss.Acceptance {
		want = append(want, c.ID, c.Statement, c.Verification)
	}
	for _, w := range want {
		if !strings.Contains(out, w) {
			t.Errorf("rendered issue is missing %q", w)
		}
	}
	if Render(iss) != out {
		t.Fatal("Render is not deterministic")
	}
}

// The blankness predicate is enumerated by hand in TWO languages, so the
// enumeration itself has to be checked rather than trusted. This is the Go
// half: over the whole rune space, Blank must agree with unicode.IsSpace
// exactly — a rune wrongly included makes legitimate issue text illegal, one
// wrongly omitted lets a blank-looking field through. The CUE half is
// fixtures/issue/issue-blank-unicode-whitespace.json, which carries the entire
// set in one field and must be rejected by `cue vet` too.
func TestBlankEnumeratesExactlyUnicodeSpace(t *testing.T) {
	for r := rune(0); r <= unicode.MaxRune; r++ {
		if !utf8.ValidRune(r) {
			continue
		}
		if got, want := Blank(string(r)), unicode.IsSpace(r); got != want {
			t.Fatalf("Blank(%U) = %v, but unicode.IsSpace = %v", r, got, want)
		}
	}
	// Not vacuous in the other direction either: a blank run followed by real
	// text is not blank, and the empty string is.
	if !Blank("") {
		t.Error("the empty string must be blank")
	}
	if Blank(" x ") {
		t.Error("a field carrying real text between whitespace is not blank")
	}
}
