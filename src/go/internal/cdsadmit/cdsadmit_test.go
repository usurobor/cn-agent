package cdsadmit

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cdsdesign"
	"github.com/usurobor/cnos/src/go/internal/cdsissue"
	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellinput"
)

// corpusDir is the SAME directory scripts/cell-schema-check.sh vets with
// `cue vet … -d '#CDSRunInput'`. Two authorities, one corpus.
const corpusDir = "../../../../schemas/cds/fixtures/runinput"

// Every negative is invalid for exactly ONE reason. The expected outcome pins
// WHICH ARM refused — a document that is merely "not admitted" would not show
// that an absent payload and a malformed one are different repairs — and the
// expected substring pins which rule inside that arm fired.
var negatives = map[string]struct {
	outcome Outcome
	reason  string
}{
	"runinput-bad-kind.json":           {OutcomeRejected, `kind must be "cnos.cds.run-input.v0"`},
	"runinput-no-issue.json":           {OutcomeIncomplete, "run input carries no issue"},
	"runinput-no-design.json":          {OutcomeIncomplete, "run input carries no design"},
	"runinput-no-subject.json":         {OutcomeIncomplete, "run input carries no subject"},
	"runinput-malformed-issue.json":    {OutcomeRejected, "cds issue: problem.diverges is required"},
	"runinput-malformed-design.json":   {OutcomeRejected, "cds design: approach is required"},
	"runinput-wrong-kind-subject.json": {OutcomeRejected, "cellwork: subject kind must be"},
	"runinput-unknown-key.json":        {OutcomeRejected, `unknown key "owner"`},
	"runinput-mixed-case-key.json":     {OutcomeRejected, `unknown key "Kind"`},
}

// The table above is exercised through Decide — THE function the runner
// dispatches, not a local reassembly of it. That matters: the earlier version
// of this file reconstructed the envelope-then-payload sequence itself and
// claimed to be "the whole door as the runner walks it", while the runner
// actually treated an envelope failure as a read error and exited differently
// with no receipt at all. The claim was true of the helper and false of the
// product. Calling the shipped function is what makes it checkable.

func corpusFiles(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		t.Fatalf("read run-input corpus: %v", err)
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

// A missing or renamed corpus directory would make every assertion below
// vacuous — the same failure class files_exist guards in the shell corpus.
func TestCorpusIsPopulatedAndFullyClassified(t *testing.T) {
	names := corpusFiles(t)
	valid := 0
	for _, n := range names {
		switch {
		case strings.HasPrefix(n, "valid-"):
			valid++
		case negatives[n].reason != "":
		default:
			t.Errorf("corpus file %q is neither a valid- fixture nor a named negative; "+
				"an unclassified fixture is vetted by CUE and by nothing here", n)
		}
	}
	if valid == 0 {
		t.Fatal("the run-input corpus has no positive fixture")
	}
	if len(negatives) != len(names)-valid {
		t.Fatalf("negatives table lists %d files, corpus has %d non-positive files",
			len(negatives), len(names)-valid)
	}
}

// AC1: a malformed issue, a malformed design and a wrong-kind subject each
// refuse with their OWN reason. The table above is what makes "their own"
// checkable — three documents refused by one generic message would satisfy a
// weaker test and would leave an author with nothing to repair.
func TestCorpusAdmission(t *testing.T) {
	for _, name := range corpusFiles(t) {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(corpusDir, name))
			if err != nil {
				t.Fatal(err)
			}
			admitted, receipt, admitErr := Decide(data)
			if strings.HasPrefix(name, "valid-") {
				if admitErr != nil {
					t.Fatalf("valid run input refused: %v", admitErr)
				}
				if receipt.Outcome != OutcomeAdmitted {
					t.Fatalf("outcome = %q, want %q", receipt.Outcome, OutcomeAdmitted)
				}
				if len(admitted.Issue) == 0 || len(admitted.Design) == 0 || len(admitted.Subject) == 0 {
					t.Fatal("an admitted run input must carry all three payloads")
				}
				return
			}
			want := negatives[name]
			if want.reason == "" {
				t.Fatalf("no expected reason declared for %q", name)
			}
			if admitErr == nil {
				t.Fatal("must be refused")
			}
			if !errors.Is(admitErr, cellfill.ErrRefused) {
				t.Fatalf("a refusal must wrap ErrRefused: %v", admitErr)
			}
			if receipt.Outcome != want.outcome {
				t.Fatalf("outcome = %q, want %q (reason: %s)", receipt.Outcome, want.outcome, receipt.Reason)
			}
			if !strings.Contains(receipt.Reason, want.reason) {
				t.Fatalf("refused for the wrong reason:\n got: %s\nwant mention of: %q", receipt.Reason, want.reason)
			}
			if admitted.Issue != nil || admitted.Design != nil || admitted.Subject != nil {
				t.Fatalf("a refusal must freeze nothing: %+v", admitted)
			}
		})
	}
}

// The three AC1 documents must not be refused by the SAME sentence. Distinct
// outcomes alone would not show it, and neither would three passing subtests
// above: this compares the reasons to each other.
func TestTheThreeRefusalReasonsAreDistinct(t *testing.T) {
	seen := make(map[string]string)
	for _, name := range []string{
		"runinput-malformed-issue.json",
		"runinput-malformed-design.json",
		"runinput-wrong-kind-subject.json",
	} {
		data, err := os.ReadFile(filepath.Join(corpusDir, name))
		if err != nil {
			t.Fatal(err)
		}
		_, receipt, err := Decide(data)
		if err == nil {
			t.Fatalf("%s must be refused", name)
		}
		if prev, dup := seen[receipt.Reason]; dup {
			t.Fatalf("%s and %s were refused by the same sentence %q", prev, name, receipt.Reason)
		}
		seen[receipt.Reason] = name
	}
}

// The authored subject is admitted while it still names a moving revision.
// This is the whole reason admission calls ParseSubject and not AdmitSubject:
// requiring 40 hex at the door would make `HEAD` — the thing a human writes —
// inadmissible and would push pinning onto the author.
func TestAnAuthoredSubjectMayStillNameAMovingRevision(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(corpusDir, "valid-run-input.json"))
	if err != nil {
		t.Fatal(err)
	}
	admitted, _, err := Decide(data)
	if err != nil {
		t.Fatalf("valid run input refused: %v", err)
	}
	if !strings.Contains(string(admitted.Subject), `"base_sha":"HEAD"`) &&
		!strings.Contains(string(admitted.Subject), `"base_sha": "HEAD"`) {
		t.Fatalf("the fixture is not the moving-revision case: %s", admitted.Subject)
	}
}

// The admitted payloads are the AUTHORED BYTES, not a re-serialization. What
// is frozen into the contract, digested, and later shown to a reader has to be
// what the author wrote — a normalizing pass would make the record a record of
// what this package understood.
func TestAdmittedPayloadsAreTheAuthoredBytes(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(corpusDir, "valid-run-input.json"))
	if err != nil {
		t.Fatal(err)
	}
	in, err := cellinput.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	admitted, _, err := Decide(data)
	if err != nil {
		t.Fatal(err)
	}
	for name, pair := range map[string][2]string{
		"issue":   {string(admitted.Issue), string(in.Issue)},
		"design":  {string(admitted.Design), string(in.Design)},
		"subject": {string(admitted.Subject), string(in.Subject)},
	} {
		if pair[0] != pair[1] {
			t.Errorf("admitted %s is not the authored bytes:\n got: %s\nwant: %s", name, pair[0], pair[1])
		}
	}
}

// relate is unit-tested DIRECTLY, because it cannot be reached in a failing
// state through Admit: both of its conjuncts are also enforced by the facet
// admissions upstream. Calling it with values the facets would never have
// produced is the only way to make it fire, and a guard nothing can exercise
// is a guard that proves nothing.
func TestRelateRejectsWhatTheFacetsWouldHaveCaught(t *testing.T) {
	goodDesign := cdsdesign.Design{Impact: []cdsdesign.Surface{{Surface: "s", Why: "w"}}}
	goodIssue := cdsissue.Issue{Acceptance: []cdsissue.Criterion{{ID: "AC1"}, {ID: "AC2"}}}

	if err := relate(goodIssue, goodDesign); err != nil {
		t.Fatalf("a well-formed pair must relate: %v", err)
	}
	dup := cdsissue.Issue{Acceptance: []cdsissue.Criterion{{ID: "AC1"}, {ID: "AC1"}}}
	if err := relate(dup, goodDesign); err == nil {
		t.Error("duplicate acceptance ids must not relate")
	} else if !strings.Contains(err.Error(), "not unique") {
		t.Errorf("wrong reason: %v", err)
	}
	blank := cdsdesign.Design{Impact: []cdsdesign.Surface{{Surface: " ", Why: "w"}}}
	if err := relate(goodIssue, blank); err == nil {
		t.Error("a blank impact surface must not relate")
	} else if !strings.Contains(err.Error(), "names no surface") {
		t.Errorf("wrong reason: %v", err)
	}
}

// Every receipt this door emits states who decided the issue was EXECUTABLE,
// not merely well-formed — because this door decides only the second (Pi #81
// C2). A receipt that carried an outcome and stayed silent on that question
// would let a reader take a structural pass for a semantic one, which is the
// overclaim the whole programme exists to stop.
func TestEveryReceiptNamesWhoAttestedSemanticAdequacy(t *testing.T) {
	for name := range negatives {
		raw, err := os.ReadFile(filepath.Join(corpusDir, name))
		if err != nil {
			t.Fatal(err)
		}
		_, r, _ := Decide(raw)
		if r.SemanticAdequacy != SemanticAdequacyOperatorAttested {
			t.Errorf("%s: receipt says %q about semantic adequacy, want %q",
				name, r.SemanticAdequacy, SemanticAdequacyOperatorAttested)
		}
	}
	valid, err := os.ReadFile(filepath.Join(corpusDir, "valid-run-input.json"))
	if err != nil {
		t.Fatal(err)
	}
	_, r, _ := Decide(valid)
	if r.Outcome != OutcomeAdmitted {
		t.Fatalf("the corpus positive must be admitted, got %q: %s", r.Outcome, r.Reason)
	}
	if r.SemanticAdequacy != SemanticAdequacyOperatorAttested {
		t.Fatalf("an ADMITTED receipt is exactly where the distinction matters, and it says %q", r.SemanticAdequacy)
	}
}
