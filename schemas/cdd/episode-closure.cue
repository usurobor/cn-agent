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
		// The complete RESOLVED seat declarations (fill-owned construction):
		// each is the tagged object whose fill selected its constructor, with
		// holes resolved and — for fills that load skills — ordered canonical
		// skill refs + content digests. Opaque at the generic output boundary
		// (Pi round-7 C2): fill whitelists are an input/registry rule; the
		// closure only requires the tag. Mode truth is bound to the parent-
		// trusted metadata by VerifyClosure, not encoded here.
		alpha: {fill!: string & !="", ...}
		beta: {fill!: string & !="", ...}
	}
	contract: {
		id:   string & !=""
		goal: string
		// The frozen issue and design — the admitted run contract — opaque at
		// the generic output boundary for the same reason the resolved seat
		// declarations are: whose contract language it is belongs to the
		// protocol, not to the closure schema. They appear here because the
		// record carries the whole frozen contract, which is what binds them
		// into the one scope-lift digest.
		issue?: {...}
		design?: {...}
		// The frozen subject — what the episode acted on — opaque here for the
		// same reason. It appears because the record carries the whole frozen
		// contract, which is what binds it into the one scope-lift digest and
		// what makes "both stations saw the same subject" a property of the
		// record rather than a claim about the runtime.
		subject?: {...}
		required_evidence?: [...#RequiredRef]
	}
	alpha: #StationRecord
	matter: {data: string}
	beta: #StationRecord
	review: {
		pass:  bool
		notes: string
		// The per-obligation coverage an assessing seat produced. OPTIONAL,
		// because the cells that predate assessment state a verdict and no
		// coverage at all — an empty list forced into their records would claim
		// that a review with zero units happened. The unit ids are opaque here
		// for the same reason the contract's issue is: the catalogue's
		// vocabulary belongs to the profile that built it.
		assessment?: [...#UnitResult]
	}
}

// #UnitResult is one obligation's disposition. `unverified` is its own value
// and not a shade of failure: "we checked and it is wrong" and "we could not
// check" call for different next work. A reason is required whenever the
// disposition is not `pass` — a judgement without a reason is not review — and
// that rule is expressed as a disjunction so an absent reason cannot unify with
// a finding.
#UnitResult: {
	unit: string & !=""
	citations?: [...string & !=""]
} & ({
	disposition: "pass"
	reason?:     string
} | {
	disposition: "finding" | "unverified"
	reason!:     string & !=""
})

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
