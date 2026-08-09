package cellkernel

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

// seqIDs is a deterministic id source for tests.
type seqIDs struct{ ep string }

func (s seqIDs) EpisodeID() string         { return s.ep }
func (s seqIDs) ExecutionID(r Role) string { return string(r) + "-" + s.ep }

func TestStatusOf(t *testing.T) {
	tests := []struct {
		name    string
		pass    bool
		dec     Decision
		want    Status
		wantErr bool
	}{
		{"pass+accept", true, Accept, Accepted, false},
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
			if err != nil || got != tc.want {
				t.Fatalf("statusOf=%q,%v want %q", got, err, tc.want)
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
	if err := VerifyReceipt(res.Receipt); err != nil {
		t.Fatalf("empty receipt must self-verify: %v", err)
	}
	t.Logf("CCNF trace: episode=%s -> verdict.pass=%v -> decision=%s -> status=%s",
		res.EpisodeID, res.Verdict.Pass, res.Decision, res.Status)
}

func TestBoolCellAcceptedAndRepair(t *testing.T) {
	acc, err := RunEpisode(context.Background(), BoolSpec(true))
	if err != nil || acc.Status != Accepted {
		t.Fatalf("bool true: status=%q err=%v", acc.Status, err)
	}
	rep, err := RunEpisode(context.Background(), BoolSpec(false))
	if err != nil || rep.Status != NeedsRepair || rep.Repair == nil {
		t.Fatalf("bool false: status=%q repair=%v err=%v", rep.Status, rep.Repair, err)
	}
}

// --- D1: identity is per-invocation, not a content alias -----------------

func TestEpisodeIdentityIsPerInvocation(t *testing.T) {
	a, _ := RunEpisode(context.Background(), BoolSpec(true))
	b, _ := RunEpisode(context.Background(), BoolSpec(true))
	if a.EpisodeID == b.EpisodeID {
		t.Fatal("same contract reused the same episode id across invocations")
	}
	if a.Receipt.AlphaExecutionID == b.Receipt.AlphaExecutionID {
		t.Fatal("execution ids are not per-invocation")
	}
}

func TestResolvedInputBoundInReceipt(t *testing.T) {
	meta := func(v string) RunMeta { return RunMeta{ResolvedSpecHash: hashJSON("spec:" + v)} }
	tr, _ := RunEpisode(context.Background(), BoolSpec(true), WithMeta(meta("true")))
	fa, _ := RunEpisode(context.Background(), BoolSpec(false), WithMeta(meta("false")))
	if tr.Receipt.SpecHash == fa.Receipt.SpecHash {
		t.Fatal("runs differing in resolved input share a resolved_spec_hash")
	}
}

// --- D2: serialized receipt re-verifies; tampering fails -----------------

func acceptedReceipt(t *testing.T) Receipt {
	t.Helper()
	res, err := RunEpisode(context.Background(), BoolSpec(true), WithIDSource(seqIDs{ep: "ep-test"}))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	return res.Receipt
}

func roundTrip(t *testing.T, rc Receipt) Receipt {
	t.Helper()
	b, err := json.Marshal(rc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out Receipt
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

func TestReceiptReVerifiesAfterSerialization(t *testing.T) {
	rc := roundTrip(t, acceptedReceipt(t))
	if err := VerifyReceipt(rc); err != nil {
		t.Fatalf("round-tripped receipt must verify: %v", err)
	}
	// The inlined evidence content is what was hashed.
	for _, e := range rc.Evidence {
		if e.Content == "" || e.Ref != "sha256:"+e.SHA256 {
			t.Fatalf("evidence not self-describing after serialization: %+v", e)
		}
	}
}

func TestTamperedReceiptFailsVerification(t *testing.T) {
	tamper := map[string]func(*Receipt){
		"flip review pass": func(r *Receipt) { r.Review.Pass = !r.Review.Pass },
		"rewrite matter":   func(r *Receipt) { r.Matter.Data = "changed" },
		"rewrite contract": func(r *Receipt) { r.Contract.Goal = "changed" },
		"forge evidence":   func(r *Receipt) { r.Evidence[0].Content = "tampered" },
		"substitute ref":   func(r *Receipt) { r.Evidence[0].Ref = "sha256:deadbeef" },
		"drop exec id":     func(r *Receipt) { r.AlphaExecutionID = "" },
	}
	for name, mut := range tamper {
		t.Run(name, func(t *testing.T) {
			rc := roundTrip(t, acceptedReceipt(t))
			mut(&rc)
			if err := VerifyReceipt(rc); err == nil {
				t.Fatalf("tamper %q passed verification", name)
			}
		})
	}
}

// --- D3: V catches a gamma binding that diverges from the record ---------

func TestValidateCatchesRecordDivergence(t *testing.T) {
	frozen := Contract{ID: "c", Goal: "g"}
	rec := EpisodeRecord{ContractHash: hashJSON(frozen), MatterHash: hashJSON(Matter{}), ReviewHash: hashJSON(Review{Pass: true}), EvidenceHash: hashJSON([]EvidenceRef(nil))}
	rc := closeReceipt(rec)
	rc.Contract = frozen // fill so internal recompute lines up
	rc.Review = Review{Pass: true}
	if v := validate(rec, rc); !v.Pass {
		t.Fatalf("faithful receipt should pass: %+v", v)
	}
	// A gamma that rewrote the matter but kept the record's hash must fail.
	bad := rc
	bad.Matter = Matter{Data: "rewritten"}
	if v := validate(rec, bad); v.Pass {
		t.Fatal("V accepted a receipt whose matter diverged from its bound hash")
	}
}

// --- D4: beta-input hash is stable and sensitive -------------------------

func TestBetaInputHashStableAndSensitive(t *testing.T) {
	base := BetaInput{Contract: Contract{ID: "c"}, ContractHash: "h", Matter: Matter{Data: "m"}, PolicyID: BetaInputPolicyID}
	if hashBetaInput(base) != hashBetaInput(base) {
		t.Fatal("beta-input hash is not stable")
	}
	other := base
	other.PolicyID = "different"
	if hashBetaInput(base) == hashBetaInput(other) {
		t.Fatal("beta-input hash ignored the policy id")
	}
}

// --- D6: kernel-boundary validation, cancellation, bounds ----------------

func TestKernelRejectsInvalidSpec(t *testing.T) {
	cases := map[string]Spec{
		"empty contract id": {Contract: Contract{ID: ""}, Alpha: NoopAlpha{}, Beta: AcceptBeta{}},
		"nil alpha":         {Contract: Contract{ID: "c"}, Alpha: nil, Beta: AcceptBeta{}},
		"bad producer":      {Contract: Contract{ID: "c", RequiredEvidence: []RequiredRef{{ID: "x", Kind: "k", Producer: "gamma"}}}, Alpha: NoopAlpha{}, Beta: AcceptBeta{}},
		"dup required":      {Contract: Contract{ID: "c", RequiredEvidence: []RequiredRef{{ID: "x", Kind: "k", Producer: RoleAlpha}, {ID: "x", Kind: "k", Producer: RoleBeta}}}, Alpha: NoopAlpha{}, Beta: AcceptBeta{}},
	}
	for name, s := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := RunEpisode(context.Background(), s); err == nil {
				t.Fatalf("invalid spec %q ran", name)
			}
		})
	}
}

// countingAlpha records whether it ran (to prove cancellation ordering).
type recordingBeta struct{ ran *bool }

func (b recordingBeta) Review(context.Context, BetaInput) (BetaResult, error) {
	*b.ran = true
	return BetaResult{Review: Review{Pass: true}}, nil
}

type cancelAfterAlpha struct{ cancel context.CancelFunc }

func (a cancelAfterAlpha) Produce(context.Context, Contract) (AlphaResult, error) {
	a.cancel()
	return AlphaResult{}, nil
}

func TestCancellationBetweenSeats(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ran := false
	s := Spec{Contract: Contract{ID: "c"}, Alpha: cancelAfterAlpha{cancel: cancel}, Beta: recordingBeta{ran: &ran}}
	if _, err := RunEpisode(ctx, s); err == nil {
		t.Fatal("want cancellation error after alpha")
	}
	if ran {
		t.Fatal("beta ran after context was cancelled")
	}
}

type oversizeAlpha struct{}

func (oversizeAlpha) Produce(context.Context, Contract) (AlphaResult, error) {
	return AlphaResult{Matter: Matter{Data: strings.Repeat("x", maxMatterBytes+1)}}, nil
}

func TestBoundedOutput(t *testing.T) {
	s := Spec{Contract: Contract{ID: "c"}, Alpha: oversizeAlpha{}, Beta: AcceptBeta{}}
	if _, err := RunEpisode(context.Background(), s); err == nil {
		t.Fatal("want error for oversized matter")
	}
}

// --- Malfunction vs contract-unmet + self-cert ---------------------------

type brokenBeta struct{}

func (brokenBeta) Review(context.Context, BetaInput) (BetaResult, error) {
	return BetaResult{}, errors.New("backend unavailable")
}

func TestBetaMalfunctionReturnsError(t *testing.T) {
	s := EmptySpec()
	s.Beta = brokenBeta{}
	if _, err := RunEpisode(context.Background(), s); err == nil {
		t.Fatal("want error from beta malfunction")
	}
}

type forgingAlpha struct{}

func (forgingAlpha) Produce(context.Context, Contract) (AlphaResult, error) {
	return AlphaResult{Matter: Matter{Data: "junk"}}, nil
}

type rejectBeta struct{}

func (rejectBeta) Review(context.Context, BetaInput) (BetaResult, error) {
	return BetaResult{Review: Review{Pass: false, Notes: "reject"}}, nil
}

func TestNoSelfCertification(t *testing.T) {
	s := Spec{Contract: Contract{ID: "c", Goal: "cannot self-certify"}, Alpha: forgingAlpha{}, Beta: rejectBeta{}}
	res, err := RunEpisode(context.Background(), s)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Status == Accepted || res.Verdict.Pass {
		t.Fatalf("self-cert leaked: status=%q", res.Status)
	}
}

// alphaMintingBetaEvidence forges a β-owned ref; the runtime stamps it
// producer=alpha and also mints the real beta_review → duplicate/authority fail.
type alphaMintingBetaEvidence struct{}

func (alphaMintingBetaEvidence) Produce(context.Context, Contract) (AlphaResult, error) {
	return AlphaResult{
		Matter: Matter{Data: "m"},
		Evidence: []EvidenceCandidate{
			{ID: "diff", Kind: "diff", Bytes: []byte("d")},
			{ID: "beta_review", Kind: "review", Bytes: []byte("fake")},
		},
	}, nil
}

func TestAlphaCannotMintBetaEvidence(t *testing.T) {
	s := Spec{
		Contract: Contract{ID: "c", Goal: "g", RequiredEvidence: []RequiredRef{
			{ID: "diff", Kind: "diff", Producer: RoleAlpha},
			{ID: "beta_review", Kind: "review", Producer: RoleBeta},
		}},
		Alpha: alphaMintingBetaEvidence{},
		Beta:  AcceptBeta{},
	}
	res, err := RunEpisode(context.Background(), s)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.Verdict.Pass {
		t.Fatal("V passed although beta_review was forged by alpha")
	}
}
