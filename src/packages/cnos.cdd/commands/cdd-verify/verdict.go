// Package cddverify implements V (Contract × Receipt → ValidationVerdict).
//
// V dispatches on the receipt's protocol_id to a CUE schema package, runs
// `cue vet` for structural validation, then applies the counterfeit-receipt
// rules (C1–C5) from the design source frozen in #389 / #367's
// RECEIPT-VALIDATION.md / #370's COHERENCE-CELL-NORMAL-FORM.md.
//
// JSON output conforms to schemas/cdd/validation_verdict.schema.json. The
// schema is the public contract; structs in this file map 1:1.
//
// Exit codes (encoded by callers via Run):
//
//	0  result: PASS
//	1  result: FAIL
//	2  V itself errored (parse fail, missing schema package, cue missing, ...)
//
// The bash and Python predecessors at this same path are removed by cycle/392
// (cnos#392, superseding cnos#391, porting the Python design from cnos#389).
package cddverify

// ValidatorIdentity is the canonical capability name pinned by
// schemas/cdd/validation_verdict.schema.json and
// RECEIPT-VALIDATION.md.
const ValidatorIdentity = "cnos.cdd.validate_receipt"

// ValidatorVersion is the deployed validator-implementation version.
// Kept at phase3.0 to match the JSON-schema example and the Python
// predecessor (the receipt-validation contract is stable; this is a
// language re-implementation, not a semantic revision).
const ValidatorVersion = "phase3.0"

// ValidationVerdict is V's emitted verdict. Matches
// schemas/cdd/validation_verdict.schema.json $id
// cnos.cdd.validation_verdict.v1.
type ValidationVerdict struct {
	Result           string            `json:"result"` // PASS | FAIL
	FailedPredicates []FailedPredicate `json:"failed_predicates"`
	Warnings         []string          `json:"warnings"`
	Provenance       Provenance        `json:"provenance"`
}

// FailedPredicate is one row in ValidationVerdict.failed_predicates.
type FailedPredicate struct {
	Predicate   string `json:"predicate"`
	EvidenceRef string `json:"evidence_ref,omitempty"`
	Diagnostic  string `json:"diagnostic"`
}

// Provenance records who emitted the verdict and against what inputs.
type Provenance struct {
	ValidatorIdentity string    `json:"validator_identity"`
	ValidatorVersion  string    `json:"validator_version"`
	CheckedAt         string    `json:"checked_at"`
	InputRefs         InputRefs `json:"input_refs"`
}

// InputRefs records pointers to the verdict's inputs.
// Only receipt_ref is required by the schema; contract_ref and
// evidence_root_ref are emitted when known.
type InputRefs struct {
	ContractRef     string `json:"contract_ref,omitempty"`
	ReceiptRef      string `json:"receipt_ref"`
	EvidenceRootRef string `json:"evidence_root_ref,omitempty"`
}

// VerdictResult is the headline value of ValidationVerdict.Result.
const (
	ResultPASS = "PASS"
	ResultFAIL = "FAIL"
)
