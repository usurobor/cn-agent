// schema/judgment.cue — the CC's δ-output: a judgment over a measurement.
//
// The CC consumes a #CM and emits exactly one #Judgment. The verdict drives
// the loop: request_planning routes to a PC; request_working routes straight
// to a WC; coherent terminates the loop. Exploratory spike.
package schema

// #Judgment is the CC's parent-facing verdict computed from a #CM.
#Judgment: {
	verdict:        "coherent" | "request_planning" | "request_working"
	cm_ref:         string
	target:         string
	decided_at_ref: string
}
