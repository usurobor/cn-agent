// see schemas/cdd/README.md §Scope-Lift Invariant.
//
// Schema ID: cnos.cdd.receipt.v1
//
// Types the parent-facing artifact of the closed coherence cell — the surface
// every downstream consumer reads instead of crawling the in-cell narrative
// files. References #ValidationVerdict, #BoundaryDecision, and #Transmissibility
// from boundary_decision.cue (same package; no import statement needed). The
// receipt is *computed* at γ close-out from the cycle's evidence graph — it is
// not free-text-authored (RECEIPT-VALIDATION.md §Q5 derivation rule).
//
// The recursive scope-lift projection this schema preserves:
//
//   closed α/β/γ cell at scope n          → α-matter at scope n+1
//   δₙ #BoundaryDecision (this receipt)   → βₙ₊₁-like discrimination at n+1
//   εₙ #ProtocolGapRef stream             → γₙ₊₁-like coordination at n+1
//
// (See schemas/cdd/README.md §Scope-Lift Invariant for the projection and
// COHERENCE-CELL.md §Recursion Equation for the underlying doctrine.)
//
// v1 is the *schema-artifact* version, not a CDD protocol version.

package cdd

// #ProtocolGapRef is one entry in the ε signal — a structured reference to
// a protocol gap surfaced during the cycle. ε's matter is the protocol
// itself; gap refs are the typed observations Phase 6's ε reads across
// receipt streams (RECEIPT-VALIDATION.md §Q3 / §Q4 + COHERENCE-CELL.md
// §ε Artifact Rule).
//
// `source` enumerates where the gap surfaced; the four-class taxonomy from
// #366 Phase 6 (cdd-skill-gap | cdd-protocol-gap | cdd-tooling-gap |
// cdd-metric-gap) is *not* typed here because that taxonomy classifies the
// gap, not its source. Source is "where in the cycle artifacts did this gap
// observation come from?"
#ProtocolGapRef: {
	id:     string
	source: "receipt" | "artifact" | "issue" | "review" | "closeout"
	ref:    string
}

// #Receipt is the parent-facing trust surface. Every field below is required
// (no `?`) unless explicitly marked optional. Missing any required field
// fails `cue vet -c -d '#Receipt' {schemas} {fixture}.yaml`.
//
// The transmissibility-derivation if-block at the bottom enforces the AC4
// verdict × action → transmissibility table structurally. Authors cannot
// relitigate the table: a fixture asserting transmissibility: accepted on
// a non-PASS verdict (or PASS + override) fails vet.
#Receipt: {
	// V's emitted verdict — δ records this block from V's return; never
	// rewrites it (RECEIPT-VALIDATION.md §Q4).
	validation: #ValidationVerdict

	// δ's parent-facing decision — non-null, required. Missing this field
	// is the structural signal that no δ-authoritative validation has been
	// invoked yet; such a receipt is in the closure-emitted-but-not-accepted
	// intermediate state (RECEIPT-VALIDATION.md §Q1 ordering rule) and is
	// not a valid receipt at parent scope. AC5 enforces.
	boundary_decision: #BoundaryDecision

	// Derived: the parent-scope-observable trust property. The structural
	// derivation below pins this from (validation.verdict × boundary_decision.action).
	// An author-asserted value inconsistent with the table fails vet.
	transmissibility: #Transmissibility

	// ε protocol-iteration signal. Required in v1 so Phase 6 reads it
	// without forcing a v2 bump. A receipt with no observed protocol gaps
	// sets protocol_gap_count: 0 and protocol_gap_refs: []. The consistency
	// constraint below ties count and refs length together (AC5).
	protocol_gap_count: int & >=0
	protocol_gap_refs: [...#ProtocolGapRef]
	// consistency: count == len(refs). Authors asserting drifted values fail vet.
	protocol_gap_count: len(protocol_gap_refs)

	// Evidence references — typed strings. Required so Phase 3's V has a
	// pinned input set. The graph schema itself is deferred to Phase 3
	// (#366 roadmap); this cycle types refs only.
	evidence_root_ref:  string
	self_coherence_ref: string
	beta_review_ref:    string
	alpha_closeout_ref: string
	beta_closeout_ref:  string
	gamma_closeout_ref: string
	diff_ref:           string
	ci_refs?: [...string]

	// Transmissibility derivation. The six valid rows from the AC4 table are
	// encoded as explicit if-branches; the two invalid rows (PASS+override,
	// non-PASS+accept/release) leave transmissibility unifying to _|_ so
	// vet rejects.
	//
	// | verdict | action                     | transmissibility   |
	// |---------|----------------------------|---------------------|
	// | PASS    | accept / release           | accepted            |
	// | PASS    | reject / repair_dispatch   | not_transmissible   |
	// | PASS    | override                   | INVALID (no PASS to override) |
	// | FAIL    | reject / repair_dispatch   | not_transmissible   |
	// | FAIL    | override                   | degraded            |
	// | FAIL    | accept / release           | INVALID (loses trust trail)   |
	if validation.verdict == "PASS" {
		if boundary_decision.action == "accept" || boundary_decision.action == "release" {
			transmissibility: "accepted"
		}
		if boundary_decision.action == "reject" || boundary_decision.action == "repair_dispatch" {
			transmissibility: "not_transmissible"
		}
		if boundary_decision.action == "override" {
			// invalid: no PASS verdict to override. See §Scope-Lift Invariant.
			transmissibility: _|_
		}
	}
	if validation.verdict == "FAIL" {
		if boundary_decision.action == "accept" || boundary_decision.action == "release" {
			// invalid: accept on FAIL loses the trust trail. Use override + degraded.
			transmissibility: _|_
		}
		if boundary_decision.action == "reject" || boundary_decision.action == "repair_dispatch" {
			transmissibility: "not_transmissible"
		}
		if boundary_decision.action == "override" {
			transmissibility: "degraded"
		}
	}

	// Open: unknown extension keys pass through. The hard-gate fields above
	// are still required; this `...` only widens the schema to accommodate
	// future v1-compatible additions (mirrors schemas/skill.cue precedent).
	...
}
