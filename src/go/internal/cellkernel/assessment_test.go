package cellkernel

import (
	"strings"
	"testing"
)

// The kernel's rules about an assessment are about whether the value can be
// READ as a verdict, never about whether the verdict is right. Each case below
// is a value no reader could act on, and each must be refused at the seal AND
// re-refused by validateRecord — because a record that self-verifies must be
// one the honest path would have produced.
func TestAssessmentIntegrityRefusesUnreadableVerdicts(t *testing.T) {
	ok := []UnitResult{{Unit: "AC1", Disposition: DispositionPass}}
	cases := []struct {
		name  string
		units []UnitResult
		want  string
	}{
		{"blank unit id", []UnitResult{{Unit: "  ", Disposition: DispositionPass}}, "blank id"},
		{
			name:  "the same unit twice",
			units: []UnitResult{{Unit: "AC1", Disposition: DispositionPass}, {Unit: "AC1", Disposition: DispositionPass}},
			want:  "more than once",
		},
		{"unknown disposition", []UnitResult{{Unit: "AC1", Disposition: "probably"}}, "unknown disposition"},
		{
			name:  "a finding with no reason",
			units: []UnitResult{{Unit: "AC1", Disposition: DispositionFinding}},
			want:  "a judgement without a reason is not review",
		},
		{
			name:  "an unverified with a whitespace reason",
			units: []UnitResult{{Unit: "AC1", Disposition: DispositionUnverified, Reason: "  \t"}},
			want:  "a judgement without a reason is not review",
		},
		{
			name:  "a citation pointing at nothing",
			units: []UnitResult{{Unit: "AC1", Disposition: DispositionPass, Citations: []string{""}}},
			want:  "blank citation",
		},
		{"more units than the bound", tooManyUnits(), "units (>"},
		{"more text than the bound", oversizeUnits(), "exceeds"},
	}
	// The well-formed value must pass, or every negative below would be
	// satisfied by a rule that refuses everything.
	if err := assessmentIntegrity(ok); err != nil {
		t.Fatalf("a well-formed assessment must be admissible: %v", err)
	}
	if err := assessmentIntegrity(nil); err != nil {
		t.Fatalf("a review with no assessment is the shape every cell had before assessment existed: %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := sealBeta(BetaOutput{Review: Review{Assessment: tc.units}}, "beta-1"); err == nil ||
				!strings.Contains(err.Error(), tc.want) {
				t.Fatalf("seal admitted it: want %q, got %v", tc.want, err)
			}
			// ...and a verifier reaches the same conclusion from the serialized
			// record alone, with no seat in sight.
			var found bool
			for _, f := range validateRecord(EpisodeRecord{Review: Review{Assessment: tc.units}}) {
				if f.Class == InvalidRecord && strings.Contains(f.Detail, tc.want) {
					found = true
				}
			}
			if !found {
				t.Fatalf("validateRecord did not name it: %+v", validateRecord(EpisodeRecord{Review: Review{Assessment: tc.units}}))
			}
		})
	}
}

func tooManyUnits() []UnitResult {
	out := make([]UnitResult, 0, maxAssessmentUnits+1)
	for i := 0; i <= maxAssessmentUnits; i++ {
		out = append(out, UnitResult{Unit: string(rune('a'+i%26)) + strings.Repeat("x", i), Disposition: DispositionPass})
	}
	return out
}

func oversizeUnits() []UnitResult {
	return []UnitResult{{
		Unit:        "AC1",
		Disposition: DispositionFinding,
		Reason:      strings.Repeat("r", maxAssessmentBytes+1),
	}}
}

// The sealed record owns its own copy. Without this a seat that kept the slice
// it returned could write through the record the digest is taken over — the
// same rule Contract.clone applies to the opaque slots, and it is checked here
// because nothing else would notice a shared backing array.
func TestSealCopiesTheAssessment(t *testing.T) {
	units := []UnitResult{{Unit: "AC1", Disposition: DispositionFinding, Reason: "no", Citations: []string{"x:1"}}}
	sealed, err := sealBeta(BetaOutput{Review: Review{Assessment: units}}, "beta-1")
	if err != nil {
		t.Fatal(err)
	}
	units[0].Disposition = DispositionPass
	units[0].Reason = "actually yes"
	units[0].Citations[0] = "y:2"
	got := sealed.review.Assessment[0]
	if got.Disposition != DispositionFinding || got.Reason != "no" || got.Citations[0] != "x:1" {
		t.Fatalf("a seat wrote through the sealed record: %+v", got)
	}
}
