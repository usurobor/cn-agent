package cellkernel

import (
	"context"
	"errors"
	"testing"
)

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
					t.Fatalf("want ErrInvalidClosure, got err=%v status=%q", err, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("statusOf(pass=%v,%q)=%q, want %q", tc.pass, tc.dec, got, tc.want)
			}
		})
	}
}

func TestEmptyCellRunsToAccepted(t *testing.T) {
	res, err := RunEpisode(context.Background(), EmptySpec())
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Status != Accepted {
		t.Fatalf("status: want accepted, got %q", res.Status)
	}
	if res.Receipt.Contract.ID != "cell-0" || res.ContractHash == "" {
		t.Fatalf("receipt did not bind the frozen contract: %+v", res.Receipt)
	}
	t.Logf("CCNF trace: episode=%s matter=%q review.pass=%v -> verdict.pass=%v -> decision=%s -> status=%s",
		res.EpisodeID, res.Matter.Data, res.Review.Pass, res.Verdict.Pass, res.Decision, res.Status)
}

func TestBoolCellAccepts(t *testing.T) {
	res, err := RunEpisode(context.Background(), BoolSpec(true))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Status != Accepted {
		t.Fatalf("status: want accepted, got %q", res.Status)
	}
	// The α "bool" evidence is bound and stamped with the α producer.
	var found bool
	for _, e := range res.Receipt.Evidence {
		if e.ID == "bool" && e.Producer == RoleAlpha {
			found = true
		}
	}
	if !found {
		t.Fatalf("bool evidence not bound with alpha producer: %+v", res.Receipt.Evidence)
	}
}

func TestBoolCellNeedsRepair(t *testing.T) {
	res, err := RunEpisode(context.Background(), BoolSpec(false))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Status != NeedsRepair || res.Repair == nil {
		t.Fatalf("want needs_repair with RepairRequest, got %q %+v", res.Status, res.Repair)
	}
}

// --- Malfunction vs contract-unmet --------------------------------------

type brokenBeta struct{}

func (brokenBeta) Review(context.Context, BetaInput) (BetaResult, error) {
	return BetaResult{}, errors.New("backend unavailable")
}

func TestBetaMalfunctionReturnsError(t *testing.T) {
	s := EmptySpec()
	s.Beta = brokenBeta{}
	if _, err := RunEpisode(context.Background(), s); err == nil {
		t.Fatal("want error from beta malfunction, got nil")
	}
}

func TestNilSeatsFailClosed(t *testing.T) {
	t.Run("nil alpha", func(t *testing.T) {
		s := EmptySpec()
		s.Alpha = nil
		if _, err := RunEpisode(context.Background(), s); err == nil {
			t.Fatal("want error for nil alpha")
		}
	})
	t.Run("typed-nil beta", func(t *testing.T) {
		s := EmptySpec()
		var b *typedNilBeta // typed nil
		s.Beta = b
		if _, err := RunEpisode(context.Background(), s); err == nil {
			t.Fatal("want error for typed-nil beta")
		}
	})
}

type typedNilBeta struct{}

func (*typedNilBeta) Review(context.Context, BetaInput) (BetaResult, error) {
	return BetaResult{}, nil
}

func TestCancelledContextFailsClosed(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := RunEpisode(ctx, EmptySpec()); err == nil {
		t.Fatal("want error for cancelled context")
	}
}

// --- Adversarial: evidence authority (D2) -------------------------------

type forgingAlpha struct{}

func (forgingAlpha) Produce(context.Context, Contract) (AlphaResult, error) {
	return AlphaResult{Matter: Matter{Data: "junk"}}, nil
}

type rejectBeta struct{}

func (rejectBeta) Review(context.Context, BetaInput) (BetaResult, error) {
	return BetaResult{Review: Review{Pass: false, Notes: "reject"}}, nil
}

func TestNoSelfCertification(t *testing.T) {
	s := Spec{
		Contract: Contract{ID: "c", Goal: "cannot self-certify"},
		Alpha:    forgingAlpha{},
		Beta:     rejectBeta{},
	}
	res, err := RunEpisode(context.Background(), s)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Status == Accepted || res.Verdict.Pass {
		t.Fatalf("self-cert leaked: status=%q verdict.pass=%v", res.Status, res.Verdict.Pass)
	}
}

// alphaMintingBetaEvidence tries to satisfy a β-owned required ref itself.
type alphaMintingBetaEvidence struct{}

func (alphaMintingBetaEvidence) Produce(context.Context, Contract) (AlphaResult, error) {
	return AlphaResult{
		Matter: Matter{Data: "m"},
		Evidence: []EvidenceRef{
			{ID: "diff", Kind: "diff", Content: "d"},
			// α forges a β-owned ref; the runtime will stamp it producer=alpha.
			{ID: "beta_review", Kind: "review", Content: "fake"},
		},
	}, nil
}

type acceptNoEvidenceBeta struct{}

func (acceptNoEvidenceBeta) Review(context.Context, BetaInput) (BetaResult, error) {
	return BetaResult{Review: Review{Pass: true}}, nil
}

func TestAlphaCannotMintBetaEvidence(t *testing.T) {
	s := Spec{
		Contract: Contract{
			ID:   "c",
			Goal: "beta must produce its own review evidence",
			RequiredEvidence: []RequiredRef{
				{ID: "diff", Kind: "diff", Producer: RoleAlpha},
				{ID: "beta_review", Kind: "review", Producer: RoleBeta},
			},
		},
		Alpha: alphaMintingBetaEvidence{},
		Beta:  acceptNoEvidenceBeta{}, // β produces no evidence
	}
	res, err := RunEpisode(context.Background(), s)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Verdict.Pass {
		t.Fatal("V passed although beta_review was minted by alpha, not beta")
	}
}

// duplicateEvidenceAlpha returns the same required id twice.
type duplicateEvidenceAlpha struct{}

func (duplicateEvidenceAlpha) Produce(context.Context, Contract) (AlphaResult, error) {
	return AlphaResult{
		Matter: Matter{Data: "m"},
		Evidence: []EvidenceRef{
			{ID: "diff", Kind: "diff", Content: "a"},
			{ID: "diff", Kind: "diff", Content: "b"},
		},
	}, nil
}

func TestDuplicateEvidenceFailsV(t *testing.T) {
	s := Spec{
		Contract: Contract{ID: "c", Goal: "g",
			RequiredEvidence: []RequiredRef{{ID: "diff", Kind: "diff", Producer: RoleAlpha}}},
		Alpha: duplicateEvidenceAlpha{},
		Beta:  acceptNoEvidenceBeta{},
	}
	res, err := RunEpisode(context.Background(), s)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Verdict.Pass {
		t.Fatal("V passed with duplicate evidence ids")
	}
}

// TestForgedHashFailsV exercises the integrity check directly, since the runtime
// otherwise recomputes the hash and a seat cannot forge it.
func TestForgedHashFailsV(t *testing.T) {
	frozen := Contract{ID: "c", Goal: "g"}
	rc := Receipt{
		ContractHash: frozen.canonicalHash(),
		Review:       Review{Pass: true},
		Evidence:     []EvidenceRef{{ID: "x", Kind: "k", Producer: RoleAlpha, Content: "real", SHA256: "deadbeef"}},
	}
	v := validate(frozen, frozen.canonicalHash(), rc)
	if v.Pass {
		t.Fatal("V passed a receipt with a forged evidence hash")
	}
}

// --- Adversarial: contract freeze (D3) ----------------------------------

// mutatingAlpha tries to relax the contract it is judged against by clearing the
// required-evidence slice it received.
type mutatingAlpha struct{}

func (mutatingAlpha) Produce(_ context.Context, c Contract) (AlphaResult, error) {
	for i := range c.RequiredEvidence {
		c.RequiredEvidence[i].ID = "neutralized"
	}
	c.RequiredEvidence = nil
	return AlphaResult{Matter: Matter{Data: "m"}}, nil // produces none of the required evidence
}

func TestHostileAlphaCannotMutateFrozenContract(t *testing.T) {
	s := Spec{
		Contract: Contract{ID: "c", Goal: "g",
			RequiredEvidence: []RequiredRef{{ID: "diff", Kind: "diff", Producer: RoleAlpha}}},
		Alpha: mutatingAlpha{},
		Beta:  acceptNoEvidenceBeta{},
	}
	res, err := RunEpisode(context.Background(), s)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	// The frozen contract is unchanged, so V still requires the diff evidence
	// that the mutating alpha never produced → not accepted.
	if res.Verdict.Pass {
		t.Fatal("frozen contract was relaxed by a hostile alpha")
	}
	if len(res.Contract.RequiredEvidence) != 1 || res.Contract.RequiredEvidence[0].ID != "diff" {
		t.Fatalf("frozen contract mutated: %+v", res.Contract.RequiredEvidence)
	}
}

// --- Beta input surface (D4) --------------------------------------------

// surfaceCheckingBeta asserts it received the runtime-owned review surface.
type surfaceCheckingBeta struct{ t *testing.T }

func (b surfaceCheckingBeta) Review(_ context.Context, in BetaInput) (BetaResult, error) {
	if in.ContractHash == "" || in.PolicyID != BetaInputPolicyID || in.BundleHash == "" {
		b.t.Errorf("beta input missing surface fields: %+v", in)
	}
	if len(in.AlphaEvidence) == 0 {
		b.t.Error("beta did not receive alpha evidence as review surface")
	}
	return BetaResult{Review: Review{Pass: true}}, nil
}

func TestBetaReceivesRuntimeOwnedSurface(t *testing.T) {
	s := Spec{
		Contract: Contract{ID: "c", Goal: "g",
			RequiredEvidence: []RequiredRef{{ID: "diff", Kind: "diff", Producer: RoleAlpha}}},
		Alpha: duplicateEvidenceAlpha{}, // produces diff evidence α sees as surface
		Beta:  surfaceCheckingBeta{t: t},
	}
	if _, err := RunEpisode(context.Background(), s); err != nil {
		t.Fatalf("run: %v", err)
	}
}
