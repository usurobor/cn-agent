// Schema ID: cnos.cellkernel.episode-closure.v0
//
// #EpisodeClosure — the terminal artifact `cn cell run` emits under the FIDO/
// functional doctrine (msg-cn-pi-cnos-cell-runner-fido-functional-44): ONE
// immutable episode record with ONE scope-lift digest, plus verdict/decision/
// status that pure functions re-derive from the record
// (cellkernel.VerifyClosure). Artifact provenance is positional — an artifact's
// producer is the record side it sits on, never a stamp. `declared_protocol`
// is provenance; `protocol_validated` is false in v0. A `stub` run is a
// non-authoritative `simulated`, never accepted.
package cdd

#Role: "alpha" | "beta"

#Artifact: {
	id:       string & !=""
	kind:     string & !=""
	encoding: "utf8"
	text:     string
}

#StationRecord: {
	execution_id: string & !=""
	artifacts: [...#Artifact]
}

#Failure: {
	class:  "contract_unmet" | "invalid_record" | "invalid_identity"
	detail: string
}

#Verdict: {
	pass: bool
	failures?: [...#Failure]
}

#RequiredRef: {id: string, kind: string, producer: #Role}

#EpisodeRecord: {
	canon:          "cnos.cellkernel.episode-record-canon.v0"
	episode_id:     string & !=""
	// How the work was produced: "stub" fabricated it (non-authoritative
	// `simulated`), "mechanical" is deterministic and reproducible from the
	// record, "cognitive" means a provider held a seat — authoritative work
	// that re-running does not reproduce.
	execution_mode: "stub" | "mechanical" | "cognitive"
	resolved_spec: {
		version:           string & !=""
		declared_protocol: string & !=""
		// Opaque at the generic output boundary (Pi round-7 C2): the builtin
		// stub|bool whitelist is an INPUT rule (#CellSpec / cellspec); the
		// kernel and its closure accept any non-empty profile. Coherence with
		// execution_mode is re-derived by the verifier, not encoded here.
		profile: string & !=""
		params?: {[string]: string}
		alpha_skills: [...string]
		beta_skills: [...string]
	}
	contract: {
		id:   string & !=""
		goal: string
		required_evidence?: [...#RequiredRef]
	}
	alpha: #StationRecord
	matter: {data: string}
	beta: #StationRecord
	review: {pass: bool, notes: string}
}

#EpisodeClosure: {
	closure_schema:     "cnos.cellkernel.episode-closure.v0"
	protocol_validated: false
	status:             "accepted" | "degraded" | "rejected" | "needs_repair" | "simulated"
	decision:           "accept" | "release" | "override" | "reject" | "repair_dispatch"
	verdict:            #Verdict
	receipt: {
		record:            #EpisodeRecord
		scope_lift_digest: =~"^[0-9a-f]{64}$"
	}
	repair?: {
		reason: string
		failed?: [...string]
	}
}
