// cnos#671 R8 — PLAN-LOCAL, TRANSITIONAL structural schema for intent.cn-intent-v1.yaml.
//
// Design input (D9 NAMES cn.intent.v1; #627 S2/#644 formalize the durable object). Closed #Intent so
// `cue vet` structurally validates the transitional bootstrap intent projection: unknown top-level
// key rejected, typed fields, required sections. Identity/carrier resolution is a DEFERRED-GO check.
//
// Invocation:
//   cue vet ./schema/ ./intent.cn-intent-v1.yaml -d '#Intent'

package crd

#Intent: {
	schema!:               "cn.intent.v1"
	id!:                   string
	source!:               string
	captured_by!:          string
	projection_status!:    string
	authored_during_cell!: string
	first_pinned!:         string
	supersedes!:           null | string
	superseded_by!:        string
	carriers!: [...#Carrier]
	statement!:            string
	scope!:                #IntentScope
	constraints!: [...string]
	desired_outcome!:      string
}

#Carrier: {
	kind!:  string
	role!:  string
	ref!:   string
	note?:  string
	text?:  string
}

#IntentScope: {
	repo!:        string
	wave!:        string
	parent_wave!: string
	in_scope!: [...string]
	out_of_scope!: [...string]
}
