package cdsassess

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/usurobor/cnos/src/go/internal/cdsissue"
	"github.com/usurobor/cnos/src/go/internal/cellcheck"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
)

// issueWith builds an ADMITTED issue with the given acceptance ids. Admitted,
// not hand-constructed: the catalogue's claim is that it derives from the
// contract that came through the door, so these tests derive it the same way.
func issueWith(t *testing.T, ids ...string) cdsissue.Issue {
	t.Helper()
	crit := make([]map[string]string, 0, len(ids))
	for _, id := range ids {
		crit = append(crit, map[string]string{
			"id":           id,
			"statement":    "statement of " + id,
			"verification": "route for " + id,
		})
	}
	raw, err := json.Marshal(map[string]any{
		"kind": cdsissue.Kind,
		"id":   "issue-under-test",
		"problem": map[string]string{
			"exists": "what exists", "expected": "what is expected", "diverges": "where they diverge",
		},
		"sources": []map[string]string{{"claim": "a claim", "path": "a/path"}},
		"scope":   map[string]any{"in": []string{"in scope"}, "out": []string{}},

		"acceptance": crit,
	})
	if err != nil {
		t.Fatal(err)
	}
	iss, err := cdsissue.Admit(raw)
	if err != nil {
		t.Fatalf("the fixture issue is not admissible, so nothing below is testing the admitted path: %v", err)
	}
	return iss
}

func ids(c Catalogue) []string {
	out := make([]string, 0, len(c.Units))
	for _, u := range c.Units {
		out = append(out, u.ID)
	}
	return out
}

// AC1. The catalogue is exactly the acceptance ids, in the issue's order, plus
// the two check units — derived from the admitted contract and from nothing
// else.
func TestCatalogueIsAcceptanceIDsPlusTheTwoCheckUnits(t *testing.T) {
	cat := Build(issueWith(t, "AC1", "AC2", "AC3"))
	want := []string{"AC1", "AC2", "AC3", UnitMatterNonEmpty, UnitProjectVerify}
	if got := ids(cat); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("catalogue = %v, want %v", got, want)
	}
	// The count is asserted independently of the ids, because it is the rule a
	// skill-ref unit would break: a catalogue that also carried the cell's
	// methodology would be longer than the issue plus two, whatever the first
	// few ids happened to be.
	if len(cat.Units) != 3+2 {
		t.Fatalf("catalogue carries %d units, want one per acceptance criterion plus the two checks", len(cat.Units))
	}
	for i, u := range cat.Units[:3] {
		if u.Kind != KindAcceptance {
			t.Errorf("unit %d kind = %q, want acceptance", i, u.Kind)
		}
		// The criterion's own words travel on the unit, so the prompt renders one
		// list rather than walking the issue again beside it.
		if u.Statement != "statement of "+u.ID || u.Verification != "route for "+u.ID {
			t.Errorf("unit %q lost its criterion: %+v", u.ID, u)
		}
	}
	for _, u := range cat.Units[3:] {
		if u.Kind != KindCheck {
			t.Errorf("unit %q kind = %q, want check", u.ID, u.Kind)
		}
	}
	// Order is the issue's, not sorted: two criteria whose ids sort the other
	// way must still come back in the order the author declared them.
	if got := ids(Build(issueWith(t, "ZZ", "AA"))); got[0] != "ZZ" || got[1] != "AA" {
		t.Fatalf("catalogue reordered the issue's criteria: %v", got)
	}
}

// Build is pure over the issue: nothing is forced until Decide runs. Without
// this, a later change could quietly decide a unit at build time and the
// forcing rules below would be testing a value that was already made up.
func TestBuildForcesNothing(t *testing.T) {
	for _, u := range Build(issueWith(t, "AC1")).Units {
		if u.Forced != nil {
			t.Fatalf("unit %q arrives already decided: %+v", u.ID, *u.Forced)
		}
	}
}

func observation(status cellcheck.Status, steps ...cellcheck.Step) cellcheck.Observation {
	return cellcheck.Observation{Recipe: cellcheck.RecipeID, Status: status, Steps: steps}
}

func unit(t *testing.T, c Catalogue, id string) Unit {
	t.Helper()
	for _, u := range c.Units {
		if u.ID == id {
			return u
		}
	}
	t.Fatalf("the catalogue carries no unit %q", id)
	return Unit{}
}

// AC2, first direction: the checker observation decides `check:project-verify`,
// and each of its three outcomes maps to its own disposition.
func TestCheckerObservationDecidesItsUnit(t *testing.T) {
	diff := cellkernel.Matter{Data: "diff --git a/x b/x\n"}
	cases := []struct {
		status cellcheck.Status
		want   cellkernel.Disposition
	}{
		{cellcheck.Pass, cellkernel.DispositionPass},
		{cellcheck.Fail, cellkernel.DispositionFinding},
		{cellcheck.Unavailable, cellkernel.DispositionUnverified},
	}
	for _, tc := range cases {
		t.Run(string(tc.status), func(t *testing.T) {
			steps := []cellcheck.Step{{Name: "build", Status: tc.status, Exit: 2, Tail: "it went wrong"}}
			cat := Build(issueWith(t, "AC1")).Decide(diff, observation(tc.status, steps...))
			u := unit(t, cat, UnitProjectVerify)
			if u.Forced == nil || u.Forced.Disposition != tc.want {
				t.Fatalf("checker %q forced %+v, want %q", tc.status, u.Forced, tc.want)
			}
			if !strings.Contains(u.Forced.Reason, cellcheck.RecipeID) {
				t.Fatalf("the forced reason does not name the recipe: %q", u.Forced.Reason)
			}
			if tc.status != cellcheck.Pass && !strings.Contains(u.Forced.Reason, "build") {
				t.Fatalf("the forced reason does not name the step that decided it: %q", u.Forced.Reason)
			}
		})
	}
}

// AC2, second direction: a cognitive answer claiming `pass` for a unit the
// checker failed is a FAULT, not a pass and not a finding — the seat claimed
// authority over a fact it was told.
func TestAnswerContradictingTheCheckerIsAFault(t *testing.T) {
	diff := cellkernel.Matter{Data: "diff --git a/x b/x\n"}
	cat := Build(issueWith(t, "AC1")).Decide(diff, observation(cellcheck.Fail,
		cellcheck.Step{Name: "test", Status: cellcheck.Fail, Exit: 1, Tail: "FAIL"}))

	claims := func(projectVerify cellkernel.Disposition) Assessment {
		return Assessment{Units: []cellkernel.UnitResult{
			{Unit: "AC1", Disposition: cellkernel.DispositionPass},
			{Unit: UnitMatterNonEmpty, Disposition: cellkernel.DispositionPass},
			{Unit: UnitProjectVerify, Disposition: projectVerify, Reason: "looks fine to me"},
		}}
	}
	_, err := Reconcile(cat, claims(cellkernel.DispositionPass), true)
	if err == nil || !strings.Contains(err.Error(), "not the seat's to decide") {
		t.Fatalf("a seat overriding the checker must be a fault, got %v", err)
	}
	// ...and the honest direction still works: repeating what the runtime
	// measured is accepted, and the RUNTIME'S reason is what the record carries.
	units, err := Reconcile(cat, claims(cellkernel.DispositionFinding), true)
	if err != nil {
		t.Fatalf("an answer agreeing with the checker must reconcile: %v", err)
	}
	got := units[len(units)-1]
	if got.Disposition != cellkernel.DispositionFinding {
		t.Fatalf("checker fail did not force a finding: %+v", got)
	}
	if strings.Contains(got.Reason, "looks fine to me") || !strings.Contains(got.Reason, cellcheck.RecipeID) {
		t.Fatalf("a measured unit must carry the runtime's reason, not the seat's: %q", got.Reason)
	}
}

// The matter unit is decided by the runtime from the matter, on the same terms.
func TestMatterUnitIsDecidedFromTheMatter(t *testing.T) {
	cases := []struct {
		name   string
		matter string
		want   cellkernel.Disposition
		reason string
	}{
		{"a diff", "diff --git a/x b/x\n@@\n", cellkernel.DispositionPass, "unified diff"},
		{"empty", "   \n", cellkernel.DispositionFinding, "empty"},
		{"prose", "no change was made to /repo at deadbeef", cellkernel.DispositionFinding, "no unified diff"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cat := Build(issueWith(t, "AC1")).Decide(cellkernel.Matter{Data: tc.matter}, observation(cellcheck.Unavailable))
			u := unit(t, cat, UnitMatterNonEmpty)
			if u.Forced == nil || u.Forced.Disposition != tc.want {
				t.Fatalf("matter %q forced %+v, want %q", tc.matter, u.Forced, tc.want)
			}
			if !strings.Contains(u.Forced.Reason, tc.reason) {
				t.Fatalf("reason %q does not mention %q", u.Forced.Reason, tc.reason)
			}
		})
	}
}

// AC3. Coverage must be EXACT, and the four ways of failing it are four
// distinct diagnostics: a reviewer reading a fault must be able to tell an
// invented unit from a dropped one.
func TestCoverageIsExact(t *testing.T) {
	cat := Build(issueWith(t, "AC1", "AC2")).Decide(
		cellkernel.Matter{Data: "diff --git a/x b/x\n"}, observation(cellcheck.Pass))

	full := []cellkernel.UnitResult{
		{Unit: "AC1", Disposition: cellkernel.DispositionPass},
		{Unit: "AC2", Disposition: cellkernel.DispositionPass},
		{Unit: UnitMatterNonEmpty, Disposition: cellkernel.DispositionPass},
		{Unit: UnitProjectVerify, Disposition: cellkernel.DispositionPass},
	}
	// The complete answer must reconcile, or every negative below would "pass"
	// for the wrong reason.
	if _, err := Reconcile(cat, Assessment{Units: full}, true); err != nil {
		t.Fatalf("the exactly-covering answer must reconcile: %v", err)
	}

	cases := []struct {
		name  string
		units []cellkernel.UnitResult
		want  string
	}{
		{"missing", full[1:], "does not dispose of catalogue unit"},
		{"duplicate", append(append([]cellkernel.UnitResult{}, full...), full[0]), "more than once"},
		{
			name:  "reordered",
			units: []cellkernel.UnitResult{full[1], full[0], full[2], full[3]},
			want:  "out of catalogue order",
		},
		{
			name: "unknown",
			units: append(append([]cellkernel.UnitResult{}, full[:3]...),
				cellkernel.UnitResult{Unit: "AC9", Disposition: cellkernel.DispositionPass}),
			want: "which is not in the catalogue",
		},
	}
	seen := map[string]bool{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Reconcile(cat, Assessment{Units: tc.units}, true)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want a fault mentioning %q, got %v", tc.want, err)
			}
			// DISTINCT, not merely non-nil: four failures reported by one message
			// would tell a reader nothing about which rule fired.
			if seen[err.Error()] {
				t.Fatalf("this fault is not distinguishable from an earlier one: %v", err)
			}
			seen[err.Error()] = true
		})
	}
}

// AC6. A view that was cut short cannot support a pass: the seat is told so in
// its prompt, and the runtime enforces it, because an instruction a seat can
// ignore is not a rule. The downgrade keeps the seat's own reason, so a reader
// can still see what it thought it had checked.
func TestATruncatedViewDowngradesAPassToUnverified(t *testing.T) {
	cat := Build(issueWith(t, "AC1")).Decide(
		cellkernel.Matter{Data: "diff --git a/x b/x\n"}, observation(cellcheck.Pass))
	answer := Assessment{Units: []cellkernel.UnitResult{
		{Unit: "AC1", Disposition: cellkernel.DispositionPass, Reason: "the file carries the symbol"},
		{Unit: UnitMatterNonEmpty, Disposition: cellkernel.DispositionPass},
		{Unit: UnitProjectVerify, Disposition: cellkernel.DispositionPass},
	}}

	complete, err := Reconcile(cat, answer, true)
	if err != nil {
		t.Fatalf("reconcile on a complete view: %v", err)
	}
	if complete[0].Disposition != cellkernel.DispositionPass {
		t.Fatalf("a complete view must leave a pass alone: %+v", complete[0])
	}

	partial, err := Reconcile(cat, answer, false)
	if err != nil {
		t.Fatalf("reconcile on a partial view: %v", err)
	}
	if partial[0].Disposition != cellkernel.DispositionUnverified {
		t.Fatalf("a pass on an incomplete view must become unverified, got %+v", partial[0])
	}
	if !strings.Contains(partial[0].Reason, "incomplete") ||
		!strings.Contains(partial[0].Reason, "the file carries the symbol") {
		t.Fatalf("the downgrade must state why and keep the seat's reason: %q", partial[0].Reason)
	}
	// The two MEASURED units are untouched: their disposition is the runtime's
	// own, so an incomplete view of the candidate says nothing about them.
	for _, u := range partial[1:] {
		if u.Disposition != cellkernel.DispositionPass {
			t.Fatalf("a measured unit must survive an incomplete view: %+v", u)
		}
	}
}

// A finding or an unverified from the seat is carried through unchanged: the
// downgrade above is one rule about one case, not a general rewrite of what the
// seat said.
func TestReconcileCarriesTheSeatsOwnJudgements(t *testing.T) {
	cat := Build(issueWith(t, "AC1", "AC2")).Decide(
		cellkernel.Matter{Data: "diff --git a/x b/x\n"}, observation(cellcheck.Pass))
	answer := Assessment{Units: []cellkernel.UnitResult{
		{Unit: "AC1", Disposition: cellkernel.DispositionFinding, Reason: "the goal is not met", Citations: []string{"x:1"}},
		{Unit: "AC2", Disposition: cellkernel.DispositionUnverified, Reason: "nothing shows this"},
		{Unit: UnitMatterNonEmpty, Disposition: cellkernel.DispositionPass},
		{Unit: UnitProjectVerify, Disposition: cellkernel.DispositionPass},
	}}
	units, err := Reconcile(cat, answer, false)
	if err != nil {
		t.Fatalf("reconcile: %v", err)
	}
	if units[0].Disposition != cellkernel.DispositionFinding || units[0].Reason != "the goal is not met" ||
		len(units[0].Citations) != 1 {
		t.Fatalf("a finding was rewritten: %+v", units[0])
	}
	if units[1].Disposition != cellkernel.DispositionUnverified {
		t.Fatalf("an unverified was rewritten: %+v", units[1])
	}
	if Pass(units) {
		t.Fatal("an assessment carrying a finding must not derive a passing review")
	}
	if s := Summary(units); !strings.Contains(s, "4 catalogue units") || !strings.Contains(s, "1 finding") {
		t.Fatalf("summary does not count what happened: %q", s)
	}
}

// The answer schema pins the catalogue's own unit ids and length, so a provider
// that honours it cannot invent or drop a unit in the first place. Reconcile
// still re-checks — the schema is a convenience, not the authority — but a
// schema that did not name the units would make the constraint decorative.
func TestAnswerSchemaPinsTheCatalogue(t *testing.T) {
	cat := Build(issueWith(t, "AC1", "AC2"))
	var schema struct {
		Properties struct {
			Units struct {
				MinItems int `json:"minItems"`
				MaxItems int `json:"maxItems"`
				Items    struct {
					Required   []string `json:"required"`
					Properties struct {
						Unit        struct{ Enum []string } `json:"unit"`
						Disposition struct{ Enum []string } `json:"disposition"`
					} `json:"properties"`
				} `json:"items"`
			} `json:"units"`
		} `json:"properties"`
	}
	raw := AnswerSchema(cat)
	if err := json.Unmarshal(raw, &schema); err != nil {
		t.Fatalf("the answer schema is not JSON a provider could be given: %v\n%s", err, raw)
	}
	u := schema.Properties.Units
	if u.MinItems != len(cat.Units) || u.MaxItems != len(cat.Units) {
		t.Fatalf("schema admits %d..%d units for a catalogue of %d", u.MinItems, u.MaxItems, len(cat.Units))
	}
	if strings.Join(u.Items.Properties.Unit.Enum, ",") != strings.Join(ids(cat), ",") {
		t.Fatalf("schema unit enum = %v, want the catalogue %v", u.Items.Properties.Unit.Enum, ids(cat))
	}
	if got := strings.Join(u.Items.Properties.Disposition.Enum, ","); got != "pass,finding,unverified" {
		t.Fatalf("schema disposition enum = %q", got)
	}
	if got := strings.Join(u.Items.Required, ","); got != "unit,disposition,reason" {
		t.Fatalf("schema does not require a reason with every disposition: %q", got)
	}
}

// The deterministic judge never passes an acceptance unit, and repeats every
// measured one exactly. Both halves are load-bearing: the first is what stops a
// fake standing in for review, the second is what stops the fake tripping the
// contradiction fault the rented seat is held to.
func TestRefusingJudgePassesNothingAndContradictsNothing(t *testing.T) {
	cat := Build(issueWith(t, "AC1", "AC2")).Decide(
		cellkernel.Matter{Data: "diff --git a/x b/x\n"}, observation(cellcheck.Pass))
	a, err := refusingJudge{}.Judge(t.Context(), "", cat)
	if err != nil {
		t.Fatalf("judge: %v", err)
	}
	units, err := Reconcile(cat, a, true)
	if err != nil {
		t.Fatalf("the deterministic answer must reconcile against the catalogue it was built from: %v", err)
	}
	for i, u := range units {
		want := cellkernel.DispositionUnverified
		if f := cat.Units[i].Forced; f != nil {
			want = f.Disposition
		}
		if u.Disposition != want {
			t.Fatalf("unit %q = %q, want %q", u.Unit, u.Disposition, want)
		}
		if strings.TrimSpace(u.Reason) == "" {
			t.Fatalf("unit %q carries no reason", u.Unit)
		}
	}
	if Pass(units) {
		t.Fatal("a seat that rented no cognition must not derive a passing review")
	}
}

// Every value this package hands the kernel must survive the kernel's own
// record boundary. Checked here rather than assumed, because a reason this
// package left blank would surface as an episode malfunction at seal time with
// nothing pointing back to the fill that produced it.
func TestReconciledUnitsSatisfyTheKernelBoundary(t *testing.T) {
	cat := Build(issueWith(t, "AC1")).Decide(cellkernel.Matter{Data: ""}, observation(cellcheck.Unavailable))
	a, err := refusingJudge{}.Judge(t.Context(), "", cat)
	if err != nil {
		t.Fatal(err)
	}
	units, err := Reconcile(cat, a, false)
	if err != nil {
		t.Fatal(err)
	}
	review := cellkernel.Review{Pass: Pass(units), Notes: Summary(units), Assessment: units}
	cl, err := cellkernel.RunEpisode(t.Context(),
		cellkernel.Spec{
			Contract: cellkernel.Contract{ID: "c", Goal: "g"},
			Alpha:    cellkernel.NoopAlpha{},
			Beta:     fixedBeta{review},
		},
		cellkernel.RunMeta{
			ExecutionMode: cellkernel.ModeMechanical,
			ResolvedSpec: cellkernel.ResolvedSpec{
				Version: "cnos.cellspec.v0", DeclaredProtocol: "p",
				Alpha: json.RawMessage(`{"fill":"a"}`), Beta: json.RawMessage(`{"fill":"b"}`),
			},
		})
	if err != nil {
		t.Fatalf("the kernel refused an assessment this package produced: %v", err)
	}
	if cl.Status != cellkernel.NeedsRepair {
		t.Fatalf("status = %q, want needs_repair", cl.Status)
	}
	// Every non-pass unit is its own named contract-unmet failure, so the
	// verdict is re-derivable from the record rather than taken from the seat.
	var named int
	for _, f := range cl.Verdict.Failures {
		if strings.Contains(f.Detail, "assessment unit ") {
			named++
		}
	}
	if named != len(units) {
		t.Fatalf("V named %d of %d non-pass units: %+v", named, len(units), cl.Verdict.Failures)
	}
}

// fixedBeta hands the kernel one prepared review, so this test is about the
// kernel's record boundary and not about the seat that would have produced it.
type fixedBeta struct{ review cellkernel.Review }

func (b fixedBeta) Review(context.Context, cellkernel.BetaInput) (cellkernel.BetaOutput, error) {
	return cellkernel.BetaOutput{Review: b.review}, nil
}

// The vocabulary and the reason rule are Reconcile's, not the kernel's seal.
// Both are ways the SEAT answered wrongly, and a rule enforced only downstream
// surfaces as an episode malfunction naming neither this fill nor the provider
// that produced the answer — a true failure attributed to the wrong component.
func TestReconcileRejectsABadDispositionAndAWithheldReason(t *testing.T) {
	cat := Build(issueWith(t, "AC1")).Decide(
		cellkernel.Matter{Data: "diff --git a/x b/x\n"}, observation(cellcheck.Pass))
	rest := []cellkernel.UnitResult{
		{Unit: UnitMatterNonEmpty, Disposition: cellkernel.DispositionPass},
		{Unit: UnitProjectVerify, Disposition: cellkernel.DispositionPass},
	}

	for name, tc := range map[string]struct {
		unit cellkernel.UnitResult
		want string
	}{
		"unknown disposition": {
			cellkernel.UnitResult{Unit: "AC1", Disposition: "probably-fine", Reason: "looks ok"},
			"the vocabulary is pass, finding or unverified",
		},
		"finding with no reason": {
			cellkernel.UnitResult{Unit: "AC1", Disposition: cellkernel.DispositionFinding},
			"with no reason",
		},
		"unverified with a blank reason": {
			cellkernel.UnitResult{Unit: "AC1", Disposition: cellkernel.DispositionUnverified, Reason: "   "},
			"with no reason",
		},
	} {
		unit, want := tc.unit, tc.want
		t.Run(name, func(t *testing.T) {
			answer := Assessment{Units: append([]cellkernel.UnitResult{unit}, rest...)}
			_, err := Reconcile(cat, answer, true)
			if err == nil {
				t.Fatal("Reconcile must refuse; the seat answered outside the contract")
			}
			if !strings.Contains(err.Error(), want) {
				t.Fatalf("refused for the wrong reason:\n got: %v\nwant mention of: %q", err, want)
			}
		})
	}

	// Not vacuous: the same shape with a real disposition and a real reason
	// reconciles, so the rules above reject the defect and not the fixture.
	ok := Assessment{Units: append([]cellkernel.UnitResult{
		{Unit: "AC1", Disposition: cellkernel.DispositionFinding, Reason: "the diff does not add the file"},
	}, rest...)}
	if _, err := Reconcile(cat, ok, true); err != nil {
		t.Fatalf("a well-formed answer must reconcile: %v", err)
	}
}
