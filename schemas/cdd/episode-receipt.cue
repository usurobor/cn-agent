// Schema ID: cnos.cellkernel.episode-receipt.v0
//
// #EpisodeReceipt — the generic, self-verifying receipt the GitHub-free runner
// (`cn cell run`) emits. It is protocol-agnostic: `declared_protocol` is
// provenance and `protocol_validated` is false in v0 (the runner does not run
// protocol-specific validation). Every hash is recomputable from the receipt's
// own content, so a parent can re-verify it out of process
// (cellkernel.VerifyReceipt). See docs/architecture/CDS-CELL-MIGRATION.md and
// Pi β msg-cn-pi-cnos-cell-prototype-rereview-33.
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

#EpisodeVerdict: {
	pass: bool
	failed?: [...string]
	warnings?: [...string]
}

#EpisodeReceipt: {
	receipt_schema:     "cnos.cellkernel.episode-receipt.v0"
	protocol_validated: bool
	execution_mode:     "stub" | "mechanical"
	status:             "accepted" | "degraded" | "rejected" | "needs_repair"
	decision:           "accept" | "release" | "override" | "reject" | "repair_dispatch"
	verdict:            #EpisodeVerdict

	episode_id:         string & !=""
	resolved_spec_hash: string
	declared_protocol:  string & !=""
	profile:            string & !=""
	params?: {[string]: string}

	contract: {
		id:   string & !=""
		goal: string
		required_evidence?: [...{id: string, kind: string, producer: #Role}]
	}
	contract_hash: =~"^[0-9a-f]{64}$"
	matter: {data: string}
	matter_hash: =~"^[0-9a-f]{64}$"
	review: {pass: bool, notes: string}
	review_hash: =~"^[0-9a-f]{64}$"
	evidence_refs: [...#EvidenceRef]
	evidence_hash:         =~"^[0-9a-f]{64}$"
	alpha_execution_id:    string & !=""
	beta_execution_id:     string & !=""
	beta_input_policy_id:  string & !=""
	beta_input_hash:       =~"^[0-9a-f]{64}$"

	repair?: {
		reason: string
		failed?: [...string]
	}
}
