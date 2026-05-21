package cddverify

// DispatchEntry pins a protocol_id to its CUE schema package.
//
// See #388's three-package schema split: cdd-generic, cds, cdr.
// Extra schema files for cdd-generic include the contract +
// boundary_decision definitions (cue needs all definitions resolvable
// in the same package context).
type DispatchEntry struct {
	Definition  string   // e.g. "#Receipt", "#CDSReceipt", "#CDRReceipt"
	SchemaFiles []string // repo-relative paths
}

// DispatchTable maps protocol_id → schema package for `cue vet`.
//
// This mirrors the Python predecessor's DISPATCH literal exactly so
// AC2 + AC3 + AC4 + AC5 oracle behavior is preserved.
var DispatchTable = map[string]DispatchEntry{
	"cnos.cdd.generic.receipt.v1": {
		Definition: "#Receipt",
		SchemaFiles: []string{
			"schemas/cdd/contract.cue",
			"schemas/cdd/boundary_decision.cue",
			"schemas/cdd/receipt.cue",
		},
	},
	"cnos.cdd.cds.receipt.v1": {
		Definition: "#CDSReceipt",
		SchemaFiles: []string{
			"schemas/cds/receipt.cue",
		},
	},
	"cnos.cdd.cdr.receipt.v1": {
		Definition: "#CDRReceipt",
		SchemaFiles: []string{
			"schemas/cdr/receipt.cue",
		},
	},
}

// CDSRequiredEvidenceRefs lists keys a CDS receipt must declare in
// evidence_refs. Used by rule C5 for a friendlier counterfeit diagnostic
// (CUE already enforces these structurally).
var CDSRequiredEvidenceRefs = []string{
	"evidence_root",
	"self_coherence",
	"beta_review",
	"alpha_closeout",
	"beta_closeout",
	"gamma_closeout",
	"diff",
}

// CDRRequiredEvidenceRefs lists keys a CDR receipt must declare in
// evidence_refs.
var CDRRequiredEvidenceRefs = []string{
	"claim",
	"data",
	"method",
	"result",
}

// KnownProtocols returns the sorted list of protocol_ids the dispatch
// table understands (used by the unknown-protocol diagnostic so the
// diagnostic itself is deterministic regardless of map iteration).
func KnownProtocols() []string {
	out := make([]string, 0, len(DispatchTable))
	for k := range DispatchTable {
		out = append(out, k)
	}
	// stable sort without importing sort just for one slice
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j-1] > out[j]; j-- {
			out[j-1], out[j] = out[j], out[j-1]
		}
	}
	return out
}
