// Package cellkernel is a reference implementation of the Coherence-Cell
// Normal Form (CCNF) single-episode kernel: the substrate-independent
// five-step closure a coherence cell runs at one scope.
//
//  1. matter   := α.produce(contract)
//  2. review   := β.review(betaInput)
//  3. receipt  := γ.close(record)                               [kernel-owned]
//  4. verdict  := V(record, receipt)                            [kernel-owned]
//  5. decision := δ.decide(receipt, verdict)                    [kernel-owned]
//
// See docs/architecture/CDS-CELL-MIGRATION.md and cnos#711. This is the
// single-episode engine: no repair loop (a future Drive wrapper), no
// composition (α-proposes / runtime-executes), no protocol dispatch.
//
// Authority model (Pi β #32 + #33). Everything a parent or another runtime
// needs to independently re-verify is bound into the receipt, and V re-derives
// it rather than trusting copied fields:
//
//   - IDENTITY (D1): each invocation gets a distinct runtime-minted episode id
//     and α/β execution ids (injectable only for deterministic tests). The
//     receipt binds a resolved-spec hash separate from the contract hash, so
//     runs that differ only in resolved input cannot share an identity.
//   - EVIDENCE CUSTODY (D2): seats return candidate {id, kind, bytes} only. The
//     runtime stamps producer/execution/digest, creates the content-addressed
//     ref, and INLINES the bytes, so the serialized receipt re-verifies out of
//     process (VerifyReceipt). β's canonical review artifact is minted by the
//     runtime from the actual review, not seat-authored.
//   - GAMMA IS SELF-PROVING (D3): the runtime builds an authoritative
//     EpisodeRecord of canonical hashes/identities; γ binds it; V recomputes
//     every hash from the receipt's own content AND compares to the record.
//   - BETA SURFACE (D4): the exact bytes handed to β are canonically hashed and
//     bound (policy id + β-input hash); V checks them.
//   - KERNEL BOUNDARY (D6): the kernel validates the spec, bounds seat output,
//     and checks cancellation between α, β, and closure — a direct Spec cannot
//     bypass cellspec's strict parse.
//
// γ, V, δ are kernel-owned and mechanical; only α and β are open seats.
package cellkernel

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Output bounds enforced at the kernel boundary (D6).
const (
	maxRequiredEvidence = 64
	maxMatterBytes      = 1 << 20 // 1 MiB
	maxReviewNotesBytes = 64 << 10
	maxEvidenceItems    = 64
	maxEvidenceBytes    = 1 << 20
)

// Role names which seat produced an artifact. The runtime assigns it.
type Role string

const (
	RoleAlpha Role = "alpha"
	RoleBeta  Role = "beta"
)

// BetaInputPolicyID identifies the review-surface policy: β sees the frozen
// contract, the matter, and α's authenticated evidence — never α private state.
const BetaInputPolicyID = "cnos.cellkernel.beta-input.v0"

// betaInputCanonVersion versions the canonical β-input encoding that is hashed.
const betaInputCanonVersion = "cnos.cellkernel.beta-input-canon.v0"

// --- Artifacts ----------------------------------------------------------

// RequiredRef names required evidence and the producer role authorized to mint
// it. V checks presence AND producer authority.
type RequiredRef struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Producer Role   `json:"producer"`
}

// Contract is the cell's input.
type Contract struct {
	ID               string        `json:"id"`
	Goal             string        `json:"goal"`
	RequiredEvidence []RequiredRef `json:"required_evidence,omitempty"`
}

func (c Contract) clone() Contract {
	cp := c
	if c.RequiredEvidence != nil {
		cp.RequiredEvidence = append([]RequiredRef(nil), c.RequiredEvidence...)
	}
	return cp
}

// Matter is α's product.
type Matter struct {
	Data string `json:"data"`
}

// Review is β's discrimination of the matter against the contract.
type Review struct {
	Pass  bool   `json:"pass"`
	Notes string `json:"notes"`
}

// EvidenceCandidate is what a seat returns: semantic identity + bytes only. The
// runtime assigns producer, execution id, ref, and digest (D2).
type EvidenceCandidate struct {
	ID    string
	Kind  string
	Bytes []byte
}

// EvidenceRef is an authenticated, self-verifying evidence record in the
// receipt. Content is inlined so the serialized receipt re-verifies out of
// process; Ref is the runtime-created content address.
type EvidenceRef struct {
	ID                  string `json:"id"`
	Kind                string `json:"kind"`
	Producer            Role   `json:"producer"`
	ProducerExecutionID string `json:"producer_execution_id"`
	Ref                 string `json:"ref"`     // "sha256:<hex>", runtime-created
	SHA256              string `json:"sha256"`  // digest of Content
	Content             string `json:"content"` // inline canonical bytes
}

// AlphaResult / BetaResult carry candidate evidence.
type AlphaResult struct {
	Matter   Matter
	Evidence []EvidenceCandidate
}

type BetaResult struct {
	Review   Review
	Evidence []EvidenceCandidate
}

// BetaInput is the runtime-owned review surface handed to β (D4).
type BetaInput struct {
	Contract      Contract
	ContractHash  string
	Matter        Matter
	AlphaEvidence []EvidenceRef
	PolicyID      string
	BundleHash    string
}

// Receipt is the parent-facing, self-verifying artifact γ emits.
type Receipt struct {
	EpisodeID        string            `json:"episode_id"`
	SpecHash         string            `json:"resolved_spec_hash"`
	DeclaredProtocol string            `json:"declared_protocol"`
	Profile          string            `json:"profile"`
	Params           map[string]string `json:"params,omitempty"`
	Contract         Contract          `json:"contract"`
	ContractHash     string            `json:"contract_hash"`
	Matter           Matter            `json:"matter"`
	MatterHash       string            `json:"matter_hash"`
	Review           Review            `json:"review"`
	ReviewHash       string            `json:"review_hash"`
	Evidence         []EvidenceRef     `json:"evidence_refs"`
	EvidenceHash     string            `json:"evidence_hash"`
	AlphaExecutionID string            `json:"alpha_execution_id"`
	BetaExecutionID  string            `json:"beta_execution_id"`
	PolicyID         string            `json:"beta_input_policy_id"`
	BetaInputHash    string            `json:"beta_input_hash"`
}

// Verdict is V's typed output.
type Verdict struct {
	Pass     bool     `json:"pass"`
	Failed   []string `json:"failed,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

// Decision is δ's boundary decision.
type Decision string

const (
	Accept         Decision = "accept"
	Release        Decision = "release"
	Override       Decision = "override"
	Reject         Decision = "reject"
	RepairDispatch Decision = "repair_dispatch"
)

// --- The two open seats -------------------------------------------------

type Alpha interface {
	Produce(ctx context.Context, c Contract) (AlphaResult, error)
}

type Beta interface {
	Review(ctx context.Context, in BetaInput) (BetaResult, error)
}

// Spec is a single-episode cell: a contract plus the two open seats.
type Spec struct {
	Contract Contract
	Alpha    Alpha
	Beta     Beta
}

// --- Identity source (D1) -----------------------------------------------

// IDSource mints episode and execution identities. Production uses a random
// source; tests may inject a deterministic one.
type IDSource interface {
	EpisodeID() string
	ExecutionID(role Role) string
}

type randomIDs struct{}

func (randomIDs) EpisodeID() string         { return "ep-" + randHex(16) }
func (randomIDs) ExecutionID(r Role) string { return string(r) + "-" + randHex(16) }

func randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand failure is fatal for identity; surface a distinct marker.
		return "randfail-" + hex.EncodeToString(b)
	}
	return hex.EncodeToString(b)
}

// RunMeta is the resolved-spec identity the runner binds (D1). It is hashed into
// resolved_spec_hash so runs differing only in resolved input differ in receipt.
type RunMeta struct {
	ResolvedSpecHash string
	DeclaredProtocol string
	Profile          string
	Params           map[string]string
}

type runConfig struct {
	ids  IDSource
	meta RunMeta
}

// RunOption customizes a run.
type RunOption func(*runConfig)

// WithIDSource injects a deterministic id source (tests only).
func WithIDSource(s IDSource) RunOption { return func(c *runConfig) { c.ids = s } }

// WithMeta binds the resolved-spec identity from the runner.
func WithMeta(m RunMeta) RunOption { return func(c *runConfig) { c.meta = m } }

// --- Status / result ----------------------------------------------------

type Status string

const (
	Accepted    Status = "accepted"
	Degraded    Status = "degraded"
	Rejected    Status = "rejected"
	NeedsRepair Status = "needs_repair"
)

type RepairRequest struct {
	Reason string   `json:"reason"`
	Failed []string `json:"failed,omitempty"`
}

type EpisodeResult struct {
	EpisodeID    string
	SpecHash     string
	ContractHash string
	Contract     Contract
	Matter       Matter
	Review       Review
	Receipt      Receipt
	Verdict      Verdict
	Decision     Decision
	Status       Status
	Repair       *RepairRequest
}

// EpisodeRecord is the runtime's authoritative account of what actually
// happened. γ binds it; V recomputes the receipt against it (D3).
type EpisodeRecord struct {
	EpisodeID        string
	SpecHash         string
	DeclaredProtocol string
	Profile          string
	Params           map[string]string
	Contract         Contract
	ContractHash     string
	Matter           Matter
	MatterHash       string
	Review           Review
	ReviewHash       string
	Evidence         []EvidenceRef
	EvidenceHash     string
	AlphaExec        string
	BetaExec         string
	PolicyID         string
	BetaInputHash    string
}

var ErrInvalidClosure = errors.New("cellkernel: inconsistent (verdict, decision) pair")

// RunEpisode executes the CCNF five-step closure for one episode.
func RunEpisode(ctx context.Context, s Spec, opts ...RunOption) (EpisodeResult, error) {
	cfg := runConfig{ids: randomIDs{}}
	for _, o := range opts {
		o(&cfg)
	}

	if err := validateSpec(s); err != nil { // D6: kernel-owned boundary validation.
		return EpisodeResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return EpisodeResult{}, fmt.Errorf("cellkernel: context: %w", err)
	}

	frozen := s.Contract.clone()
	contractHash := hashJSON(frozen)
	episodeID := cfg.ids.EpisodeID()
	alphaExec := cfg.ids.ExecutionID(RoleAlpha)
	betaExec := cfg.ids.ExecutionID(RoleBeta)

	aRes, err := s.Alpha.Produce(ctx, frozen.clone()) // 1
	if err != nil {
		return EpisodeResult{}, fmt.Errorf("alpha produce: %w", err)
	}
	if err := boundMatter(aRes.Matter); err != nil {
		return EpisodeResult{}, err
	}
	alphaEv, err := authenticate(aRes.Evidence, RoleAlpha, alphaExec)
	if err != nil {
		return EpisodeResult{}, fmt.Errorf("alpha evidence: %w", err)
	}
	if err := ctx.Err(); err != nil { // D6: cancellation between α and β.
		return EpisodeResult{}, fmt.Errorf("cellkernel: context after alpha: %w", err)
	}

	bin := BetaInput{
		Contract:      frozen.clone(),
		ContractHash:  contractHash,
		Matter:        aRes.Matter,
		AlphaEvidence: cloneEvidence(alphaEv),
		PolicyID:      BetaInputPolicyID,
	}
	betaInputHash := hashBetaInput(bin)
	bin.BundleHash = betaInputHash

	bRes, err := s.Beta.Review(ctx, bin) // 2
	if err != nil {
		return EpisodeResult{}, fmt.Errorf("beta review: %w", err)
	}
	if err := boundReview(bRes.Review); err != nil {
		return EpisodeResult{}, err
	}
	betaEv, err := mintBetaEvidence(bRes, betaExec) // runtime owns beta_review
	if err != nil {
		return EpisodeResult{}, fmt.Errorf("beta evidence: %w", err)
	}
	if err := ctx.Err(); err != nil { // D6: cancellation before closure.
		return EpisodeResult{}, fmt.Errorf("cellkernel: context after beta: %w", err)
	}

	evidence := append(cloneEvidence(alphaEv), betaEv...)

	record := EpisodeRecord{
		EpisodeID:        episodeID,
		SpecHash:         cfg.meta.ResolvedSpecHash,
		DeclaredProtocol: cfg.meta.DeclaredProtocol,
		Profile:          cfg.meta.Profile,
		Params:           cfg.meta.Params,
		Contract:         frozen,
		ContractHash:     contractHash,
		Matter:           aRes.Matter,
		MatterHash:       hashJSON(aRes.Matter),
		Review:           bRes.Review,
		ReviewHash:       hashJSON(bRes.Review),
		Evidence:         evidence,
		EvidenceHash:     hashJSON(evidence),
		AlphaExec:        alphaExec,
		BetaExec:         betaExec,
		PolicyID:         BetaInputPolicyID,
		BetaInputHash:    betaInputHash,
	}

	receipt := closeReceipt(record)      // 3 (kernel γ)
	verdict := validate(record, receipt) // 4 (kernel V)
	decision := decide(verdict)          // 5 (kernel δ)

	status, err := statusOf(verdict, decision)
	if err != nil {
		return EpisodeResult{}, err
	}

	res := EpisodeResult{
		EpisodeID:    episodeID,
		SpecHash:     record.SpecHash,
		ContractHash: contractHash,
		Contract:     frozen,
		Matter:       aRes.Matter,
		Review:       bRes.Review,
		Receipt:      receipt,
		Verdict:      verdict,
		Decision:     decision,
		Status:       status,
	}
	if status == NeedsRepair {
		res.Repair = &RepairRequest{Reason: "contract unmet", Failed: verdict.Failed}
	}
	return res, nil
}

func statusOf(v Verdict, d Decision) (Status, error) {
	switch {
	case v.Pass && (d == Accept || d == Release):
		return Accepted, nil
	case !v.Pass && d == Override:
		return Degraded, nil
	case d == Reject:
		return Rejected, nil
	case d == RepairDispatch:
		return NeedsRepair, nil
	default:
		return "", fmt.Errorf("%w: verdict.pass=%v decision=%q", ErrInvalidClosure, v.Pass, d)
	}
}

// --- D6: kernel-boundary spec validation + output bounds ----------------

func validateSpec(s Spec) error {
	if seatIsNil(s.Alpha) {
		return errors.New("cellkernel: spec has nil alpha")
	}
	if seatIsNil(s.Beta) {
		return errors.New("cellkernel: spec has nil beta")
	}
	if s.Contract.ID == "" {
		return errors.New("cellkernel: contract.id is empty")
	}
	if len(s.Contract.RequiredEvidence) > maxRequiredEvidence {
		return fmt.Errorf("cellkernel: too many required evidence refs (%d > %d)", len(s.Contract.RequiredEvidence), maxRequiredEvidence)
	}
	seen := make(map[string]bool)
	for _, r := range s.Contract.RequiredEvidence {
		if r.ID == "" || r.Kind == "" {
			return errors.New("cellkernel: required evidence needs non-empty id and kind")
		}
		if r.Producer != RoleAlpha && r.Producer != RoleBeta {
			return fmt.Errorf("cellkernel: required evidence %q has invalid producer %q", r.ID, r.Producer)
		}
		if seen[r.ID] {
			return fmt.Errorf("cellkernel: duplicate required evidence id %q", r.ID)
		}
		seen[r.ID] = true
	}
	return nil
}

func boundMatter(m Matter) error {
	if len(m.Data) > maxMatterBytes {
		return fmt.Errorf("cellkernel: matter exceeds %d bytes", maxMatterBytes)
	}
	return nil
}

func boundReview(r Review) error {
	if len(r.Notes) > maxReviewNotesBytes {
		return fmt.Errorf("cellkernel: review notes exceed %d bytes", maxReviewNotesBytes)
	}
	return nil
}

// --- D2: runtime evidence authentication + minting ----------------------

func authenticate(candidates []EvidenceCandidate, role Role, execID string) ([]EvidenceRef, error) {
	if len(candidates) > maxEvidenceItems {
		return nil, fmt.Errorf("too many evidence items (%d > %d)", len(candidates), maxEvidenceItems)
	}
	out := make([]EvidenceRef, 0, len(candidates))
	for _, c := range candidates {
		if c.ID == "" || c.Kind == "" {
			return nil, errors.New("evidence candidate needs non-empty id and kind")
		}
		if len(c.Bytes) > maxEvidenceBytes {
			return nil, fmt.Errorf("evidence %q exceeds %d bytes", c.ID, maxEvidenceBytes)
		}
		sum := sha256hex(c.Bytes)
		out = append(out, EvidenceRef{
			ID:                  c.ID,
			Kind:                c.Kind,
			Producer:            role,
			ProducerExecutionID: execID,
			Ref:                 "sha256:" + sum,
			SHA256:              sum,
			Content:             string(c.Bytes),
		})
	}
	return out, nil
}

// mintBetaEvidence authenticates β's candidates, then replaces/mints the
// canonical `beta_review` evidence from the ACTUAL review (D2): a seat cannot
// substitute unrelated bytes under that label.
func mintBetaEvidence(b BetaResult, execID string) ([]EvidenceRef, error) {
	filtered := make([]EvidenceCandidate, 0, len(b.Evidence))
	for _, c := range b.Evidence {
		if c.ID == "beta_review" {
			continue // runtime owns this id
		}
		filtered = append(filtered, c)
	}
	reviewBytes, _ := json.Marshal(b.Review)
	filtered = append([]EvidenceCandidate{{ID: "beta_review", Kind: "review", Bytes: reviewBytes}}, filtered...)
	return authenticate(filtered, RoleBeta, execID)
}

func cloneEvidence(in []EvidenceRef) []EvidenceRef {
	return append([]EvidenceRef(nil), in...)
}

// --- D4: canonical β-input hash -----------------------------------------

func hashBetaInput(in BetaInput) string {
	type evCanon struct{ ID, Kind, Producer, Exec, SHA string }
	type canon struct {
		Version      string
		PolicyID     string
		ContractHash string
		Matter       string
		Alpha        []evCanon
	}
	c := canon{Version: betaInputCanonVersion, PolicyID: in.PolicyID, ContractHash: in.ContractHash, Matter: in.Matter.Data}
	for _, e := range in.AlphaEvidence {
		c.Alpha = append(c.Alpha, evCanon{e.ID, e.Kind, string(e.Producer), e.ProducerExecutionID, e.SHA256})
	}
	return hashJSON(c)
}

// --- Kernel-owned γ / V / δ ----------------------------------------------

func closeReceipt(r EpisodeRecord) Receipt {
	return Receipt{
		EpisodeID:        r.EpisodeID,
		SpecHash:         r.SpecHash,
		DeclaredProtocol: r.DeclaredProtocol,
		Profile:          r.Profile,
		Params:           r.Params,
		Contract:         r.Contract,
		ContractHash:     r.ContractHash,
		Matter:           r.Matter,
		MatterHash:       r.MatterHash,
		Review:           r.Review,
		ReviewHash:       r.ReviewHash,
		Evidence:         r.Evidence,
		EvidenceHash:     r.EvidenceHash,
		AlphaExecutionID: r.AlphaExec,
		BetaExecutionID:  r.BetaExec,
		PolicyID:         r.PolicyID,
		BetaInputHash:    r.BetaInputHash,
	}
}

// validate (V) recomputes every binding from the receipt's own content and
// compares it to the authoritative record (D3). It does not re-judge β's prose.
func validate(rec EpisodeRecord, rc Receipt) Verdict {
	var failed []string
	add := func(cond bool, msg string) {
		if cond {
			failed = append(failed, msg)
		}
	}

	// Internal consistency: content must hash to its bound hash.
	add(hashJSON(rc.Contract) != rc.ContractHash, "contract hash mismatch")
	add(hashJSON(rc.Matter) != rc.MatterHash, "matter hash mismatch")
	add(hashJSON(rc.Review) != rc.ReviewHash, "review hash mismatch")
	add(hashJSON(rc.Evidence) != rc.EvidenceHash, "evidence hash mismatch")

	// Agreement with the runtime's authoritative account.
	add(rc.ContractHash != rec.ContractHash, "contract not the frozen contract")
	add(rc.MatterHash != rec.MatterHash, "matter not the produced matter")
	add(rc.ReviewHash != rec.ReviewHash, "review not the actual review")
	add(rc.EvidenceHash != rec.EvidenceHash, "evidence not the authenticated evidence")
	add(rc.AlphaExecutionID != rec.AlphaExec, "alpha execution id mismatch")
	add(rc.BetaExecutionID != rec.BetaExec, "beta execution id mismatch")
	add(rc.SpecHash != rec.SpecHash, "resolved spec hash mismatch")
	add(rc.PolicyID != rec.PolicyID, "beta-input policy mismatch")
	add(rc.BetaInputHash != rec.BetaInputHash, "beta-input hash mismatch")

	// Evidence integrity + required presence/authority + uniqueness.
	failed = append(failed, checkEvidence(rc.Contract, rc.Evidence)...)

	if !rc.Review.Pass {
		failed = append(failed, "review.pass=false")
	}

	if len(failed) > 0 {
		return Verdict{Pass: false, Failed: failed}
	}
	return Verdict{Pass: true}
}

// checkEvidence verifies each ref's content hash + ref form, id uniqueness, and
// that every required ref is present with an authorized producer. Shared by V
// and VerifyReceipt so the in-process and out-of-process checks agree.
func checkEvidence(contract Contract, refs []EvidenceRef) []string {
	var failed []string
	seen := make(map[string]int)
	for _, e := range refs {
		seen[e.ID]++
		sum := sha256hex([]byte(e.Content))
		if sum != e.SHA256 {
			failed = append(failed, "evidence hash mismatch: "+e.ID)
		}
		if e.Ref != "sha256:"+e.SHA256 {
			failed = append(failed, "evidence ref not content-addressed: "+e.ID)
		}
		if e.Producer != RoleAlpha && e.Producer != RoleBeta {
			failed = append(failed, "evidence has no producer: "+e.ID)
		}
	}
	for id, n := range seen {
		if n > 1 {
			failed = append(failed, "duplicate evidence id: "+id)
		}
	}
	for _, req := range contract.RequiredEvidence {
		ok := false
		for _, e := range refs {
			if e.ID == req.ID && e.Kind == req.Kind && e.Producer == req.Producer {
				ok = true
				break
			}
		}
		if !ok {
			failed = append(failed, fmt.Sprintf("missing/unauthorized required evidence: %s (want producer %s)", req.ID, req.Producer))
		}
	}
	return failed
}

func decide(v Verdict) Decision {
	if v.Pass {
		return Accept
	}
	return RepairDispatch
}

// --- helpers ------------------------------------------------------------

func hashJSON(v any) string {
	b, _ := json.Marshal(v)
	return sha256hex(b)
}

func sha256hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func seatIsNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Func, reflect.Map, reflect.Slice, reflect.Chan:
		return rv.IsNil()
	default:
		return false
	}
}
