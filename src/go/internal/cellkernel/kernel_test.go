package cellkernel

import (
	"context"
	"errors"
	"testing"
)

// TestStatusOf is the table-driven contract for the (verdict, decision) → status
// mapping (Pi D1): consistent pairs map to a Status; inconsistent pairs are a
// typed error, never a returned closed cell.
func TestStatusOf(t *testing.T) {
	tests := []struct {
		name    string
		pass    bool
		dec     Decision
		want    Status
		wantErr bool
	}{
		{"pass+accept", true, Accept, Accepted, false},
		{"pass+release", true, Release, Accepted, false},
		{"fail+override", false, Override, Degraded, false},
		{"any+reject", true, Reject, Rejected, false},
		{"fail+repair", false, RepairDispatch, NeedsRepair, false},
		{"pass+override_inconsistent", true, Override, "", true},
		{"fail+accept_inconsistent", false, Accept, "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := statusOf(Verdict{Pass: tc.pass}, tc.dec)
			if tc.wantErr {
				if !errors.Is(err, ErrInvalidClosure) {
					t.Fatalf("statusOf(pass=%v, %q): want ErrInvalidClosure, got err=%v status=%q", tc.pass, tc.dec, err, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("statusOf(pass=%v, %q): unexpected error: %v", tc.pass, tc.dec, err)
			}
			if got != tc.want {
				t.Errorf("statusOf(pass=%v, %q) = %q, want %q", tc.pass, tc.dec, got, tc.want)
			}
		})
	}
}

// TestEmptyCellRunsToAccepted runs Case 0 through the full closure and asserts
// it terminates `accepted` with the contract bound.
func TestEmptyCellRunsToAccepted(t *testing.T) {
	res, err := RunEpisode(context.Background(), EmptySpec())
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Status != Accepted {
		t.Fatalf("status: want %q, got %q", Accepted, res.Status)
	}
	if res.Decision != Accept {
		t.Fatalf("decision: want %q, got %q", Accept, res.Decision)
	}
	if !res.Verdict.Pass {
		t.Fatalf("verdict: want PASS, got %+v", res.Verdict)
	}
	if res.Receipt.Contract.ID != "cell-0" {
		t.Fatalf("receipt did not bind the contract: %+v", res.Receipt)
	}
	t.Logf("CCNF trace: matter=%q review.pass=%v -> receipt(contract=%s) -> verdict.pass=%v -> decision=%s -> status=%s",
		res.Matter.Data, res.Review.Pass, res.Receipt.Contract.ID, res.Verdict.Pass, res.Decision, res.Status)
}

// TestBoolCellAccepts: Case 1, bool true → accepted, and the required α evidence
// is bound into the receipt.
func TestBoolCellAccepts(t *testing.T) {
	res, err := RunEpisode(context.Background(), BoolSpec(true))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Status != Accepted {
		t.Fatalf("status: want %q, got %q", Accepted, res.Status)
	}
	if !hasEvidence(res.Receipt.EvidenceRefs, "bool", "value") {
		t.Fatalf("required evidence not bound: %+v", res.Receipt.EvidenceRefs)
	}
}

// TestBoolCellNeedsRepair: Case 1, bool false → contract-unmet → needs_repair,
// parent stays open, no error, Repair request set.
func TestBoolCellNeedsRepair(t *testing.T) {
	res, err := RunEpisode(context.Background(), BoolSpec(false))
	if err != nil {
		t.Fatalf("run: unexpected error: %v", err)
	}
	if res.Status != NeedsRepair {
		t.Fatalf("status: want %q, got %q", NeedsRepair, res.Status)
	}
	if res.Decision != RepairDispatch {
		t.Fatalf("decision: want %q, got %q", RepairDispatch, res.Decision)
	}
	if res.Repair == nil {
		t.Fatal("needs_repair result must carry a RepairRequest")
	}
}

// --- Pi-mandated negative tests -----------------------------------------

// brokenBeta fails to run: a malfunction, distinct from a reject.
type brokenBeta struct{}

func (brokenBeta) Review(context.Context, Contract, Matter) (BetaResult, error) {
	return BetaResult{}, errors.New("backend unavailable")
}

// TestBetaMalfunctionReturnsError: a seat that cannot run makes RunEpisode
// return an error (the episode does not close) — distinct from a needs_repair
// contract-unmet episode.
func TestBetaMalfunctionReturnsError(t *testing.T) {
	s := EmptySpec()
	s.Beta = brokenBeta{}

	if _, err := RunEpisode(context.Background(), s); err == nil {
		t.Fatal("run: want error from beta malfunction, got nil")
	}
}

// TestNilSeatsFailClosed (D4): a nil α or β is a wrapped error before any seat
// runs — never a panic.
func TestNilSeatsFailClosed(t *testing.T) {
	t.Run("nil alpha", func(t *testing.T) {
		s := EmptySpec()
		s.Alpha = nil
		if _, err := RunEpisode(context.Background(), s); err == nil {
			t.Fatal("want error for nil alpha, got nil")
		}
	})
	t.Run("nil beta", func(t *testing.T) {
		s := EmptySpec()
		s.Beta = nil
		if _, err := RunEpisode(context.Background(), s); err == nil {
			t.Fatal("want error for nil beta, got nil")
		}
	})
}

// forgingAlpha produces matter but claims, via evidence, that a rejecting review
// "passed". It models a hostile producer trying to self-certify.
type forgingAlpha struct{}

func (forgingAlpha) Produce(context.Context, Contract) (AlphaResult, error) {
	return AlphaResult{Matter: Matter{Data: "junk"}}, nil
}

// rejectBeta discriminates every matter as failing: it ran, it said no.
type rejectBeta struct{}

func (rejectBeta) Review(context.Context, Contract, Matter) (BetaResult, error) {
	return BetaResult{Review: Review{Pass: false, Notes: "reject for test"}}, nil
}

// TestNoSelfCertification (D2): a rejecting β cannot be turned into an
// acceptance. Because γ/V/δ are kernel-owned (not injectable via Spec), there is
// no seam through which α or a caller can rewrite β's review. A rejecting β
// therefore always closes needs_repair, never accepted.
func TestNoSelfCertification(t *testing.T) {
	s := Spec{
		Contract: Contract{ID: "cell-self-cert", Goal: "cannot self-certify"},
		Alpha:    forgingAlpha{},
		Beta:     rejectBeta{},
	}
	res, err := RunEpisode(context.Background(), s)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Status == Accepted {
		t.Fatalf("self-certification leaked: rejecting beta produced status %q", res.Status)
	}
	if res.Verdict.Pass {
		t.Fatal("V passed a rejecting review — binding validation failed")
	}
}

// TestMissingRequiredEvidenceFailsV: if the contract requires an evidence ref
// that no seat produced, V FAILs even when β passed (Pi D2/Q4 binding check).
func TestMissingRequiredEvidenceFailsV(t *testing.T) {
	s := Spec{
		Contract: Contract{
			ID:               "cell-needs-evidence",
			Goal:             "requires a diff that alpha never produced",
			RequiredEvidence: []RequiredRef{{ID: "diff", Kind: "diff"}},
		},
		Alpha: NoopAlpha{}, // produces no evidence
		Beta:  AcceptBeta{},
	}
	res, err := RunEpisode(context.Background(), s)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Verdict.Pass {
		t.Fatal("V passed despite missing required evidence")
	}
	if res.Status != NeedsRepair {
		t.Fatalf("status: want %q, got %q", NeedsRepair, res.Status)
	}
}
