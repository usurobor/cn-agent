// schema/cm.cue — the Coherence Measurement object (the cnos↔tsc seam).
//
// #CM is the runner/provider-produced MEASUREMENT of a target (NOT the CC's
// judgment; the CC consumes a #CM and emits a #Judgment). v0's provider is a
// mechanical checker over the toy target; a tsc/agent measurement provider
// drops in behind this same schema (WC-2 of #672). Exploratory spike.
package schema

// #Defect is one measured incoherence against the target.
#Defect: {
	id:     string
	kind:   string
	detail: string
}

// #CM is the measurement object. alpha_score is a coherence score in [0,1]
// (1.0 = fully coherent); incoherence is the defect count. v0 scores
// mechanically as 1 - min(1, len(defects)); a richer provider refines the
// score behind this schema without changing the seam.
#CM: {
	target:      string
	alpha_score: number & >=0 & <=1
	defects: [...#Defect]
	provenance: {
		measured_by: string
		method:      string
		cm_version:  string
	}
	measured_at_ref: string
}
