// Schema ID: cnos.cellkernel.episode-envelope.v0
//
// #EpisodeEnvelope — the terminal, whole-object-verifiable artifact that the
// GitHub-free runner (`cn cell run`) emits. Every field re-derives from content
// via cellkernel.VerifyEnvelope: verdict←V(receipt), decision←δ(verdict),
// status←(decision, execution_mode); resolved_spec_hash recomputes from
// resolved_spec; beta_input_hash recomputes from the receipt. It is
// protocol-agnostic — `declared_protocol` is provenance and `protocol_validated`
// is false in v0. A `stub` run is non-authoritative `simulated`, never accepted.
//
// See docs/architecture/CDS-CELL-MIGRATION.md and Pi PR-#718 β.
package cdd

#Role: "alpha" | "beta"

#EvidenceRef: {
	id:                    string & !=""
	kind:                  string & !=""
	producer:              #Role
	producer_execution_id: string & !=""
	ref:                   =~"^sha256:[0-9a-f]{64}$"
	sha256:                =~"^[0-9a-f]{64}$"
	content:               string
}

#Failure: {
	class:  "contract_unmet" | "invalid_receipt" | "invalid_evidence" | "invalid_identity" | "invalid_independence"
	detail: string
}

#Verdict: {
	pass: bool
	failures?: [...#Failure]
}

#RequiredRef: {id: string, kind: string, producer: #Role}

#Receipt: {
	episode_id: string & !=""
	contract: {
		id:   string & !=""
		goal: string
		required_evidence?: [...#RequiredRef]
	}
	contract_hash: =~"^[0-9a-f]{64}$"
	matter: {data: string}
	matter_hash: =~"^[0-9a-f]{64}$"
	review: {pass: bool, notes: string}
	review_hash: =~"^[0-9a-f]{64}$"
	evidence_refs: [...#EvidenceRef]
	evidence_hash:        =~"^[0-9a-f]{64}$"
	alpha_execution_id:   string & !=""
	beta_execution_id:    string & !=""
	beta_input_policy_id: string & !=""
	beta_input_hash:      =~"^[0-9a-f]{64}$"
}

#ResolvedSpec: {
	canon:             "cnos.cellkernel.resolved-spec-canon.v0"
	version:           string & !=""
	declared_protocol: string & !=""
	profile:           "stub" | "bool"
	params?: {[string]: string}
	alpha_skills: [...string]
	beta_skills: [...string]
	contract: #Receipt.contract
}

#EpisodeEnvelope: {
	envelope_schema:    "cnos.cellkernel.episode-envelope.v0"
	protocol_validated: false
	execution_mode:     "stub" | "mechanical"
	status:             "accepted" | "degraded" | "rejected" | "needs_repair" | "simulated"
	decision:           "accept" | "release" | "override" | "reject" | "repair_dispatch"
	verdict:            #Verdict
	resolved_spec:      #ResolvedSpec
	resolved_spec_hash: =~"^[0-9a-f]{64}$"
	receipt:            #Receipt
	repair?: {
		reason: string
		failed?: [...string]
	}
}
