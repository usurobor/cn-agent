package cellkernel

import (
	"context"
	"errors"
	"testing"
)

// TestOutcomeOf is the table-driven contract for the four-outcome mapping
// (CCNF §Cell Outcomes): outcome is a pure function of (verdict, decision).
func TestOutcomeOf(t *testing.T) {
	tests := []struct {
		name string
		pass bool
		dec  Decision
		want Outcome
	}{
		{"pass+accept", true, Accept, Accepted},
		{"pass+release", true, Release, Accepted},
		{"fail+override", false, Override, Degraded},
		{"any+reject", true, Reject, Blocked},
		{"any+repair", false, RepairDispatch, Blocked},
		{"pass+override_inconsistent", true, Override, Invalid},
		{"fail+accept_inconsistent", false, Accept, Invalid},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := outcomeOf(Verdict{Pass: tc.pass}, tc.dec)
			if got != tc.want {
				t.Errorf("outcomeOf(pass=%v, %q) = %q, want %q", tc.pass, tc.dec, got, tc.want)
			}
		})
	}
}

// TestEmptyCellRunsToAccepted runs the empty cell through the full five-step
// closure and asserts it terminates in `accepted` with every seat having fired.
func TestEmptyCellRunsToAccepted(t *testing.T) {
	cc, err := Run(context.Background(), EmptySpec())
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if cc.Outcome != Accepted {
		t.Fatalf("outcome: want %q, got %q", Accepted, cc.Outcome)
	}
	if cc.Decision != Accept {
		t.Fatalf("decision: want %q, got %q", Accept, cc.Decision)
	}
	if !cc.Verdict.Pass {
		t.Fatalf("verdict: want PASS, got %+v", cc.Verdict)
	}
	if cc.Receipt.Contract.ID != "cell-0" {
		t.Fatalf("receipt did not bind the contract: %+v", cc.Receipt)
	}
	t.Logf("CCNF trace: matter=%q review.pass=%v -> receipt(contract=%s) -> verdict.pass=%v -> decision=%s -> outcome=%s",
		cc.Matter.Data, cc.Review.Pass, cc.Receipt.Contract.ID, cc.Verdict.Pass, cc.Decision, cc.Outcome)
}

// rejectBeta discriminates every matter as failing: it ran, it said no.
type rejectBeta struct{}

func (rejectBeta) Review(context.Context, Contract, Matter) (Review, error) {
	return Review{Pass: false, Notes: "reject for test"}, nil
}

// TestRejectingBetaBlocks: contract-unmet flows FAIL -> repair_dispatch ->
// blocked through the same kernel, and Run returns no error (the cell closed).
func TestRejectingBetaBlocks(t *testing.T) {
	s := EmptySpec()
	s.Beta = rejectBeta{}

	cc, err := Run(context.Background(), s)
	if err != nil {
		t.Fatalf("run: unexpected error: %v", err)
	}
	if cc.Verdict.Pass {
		t.Fatalf("verdict: want FAIL, got PASS")
	}
	if cc.Decision != RepairDispatch {
		t.Fatalf("decision: want %q, got %q", RepairDispatch, cc.Decision)
	}
	if cc.Outcome != Blocked {
		t.Fatalf("outcome: want %q, got %q", Blocked, cc.Outcome)
	}
}

// brokenBeta fails to run: a malfunction, distinct from a reject.
type brokenBeta struct{}

func (brokenBeta) Review(context.Context, Contract, Matter) (Review, error) {
	return Review{}, errors.New("backend unavailable")
}

// TestBetaMalfunctionReturnsError: a seat that cannot run makes Run return an
// error (the cell does not close) — distinct from a contract-unmet blocked cell.
func TestBetaMalfunctionReturnsError(t *testing.T) {
	s := EmptySpec()
	s.Beta = brokenBeta{}

	_, err := Run(context.Background(), s)
	if err == nil {
		t.Fatal("run: want error from beta malfunction, got nil")
	}
}
