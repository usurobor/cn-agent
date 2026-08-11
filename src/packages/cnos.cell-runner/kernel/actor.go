package kernel

// Actor is the "assisted" seam: the boundary between the mechanical kernel and
// whatever performs each step. v0 ships deterministic stub actors
// (../actor/stub); agent-backed actors (real α/β/CC judgment) implement this
// SAME interface next (v1) with nothing in the kernel or Drive loop changing.
//
// The kernel is uniform across all three cells: every cell runs the same six
// steps. A cell's telos is expressed purely by WHICH Actor is bound:
//   - WC: Alpha edits the target file; Validate asserts acceptance predicates.
//   - PC: Alpha emits a wave (one WC contract); the artifact is the relation graph.
//   - CC: Alpha IS the Measure step (produces the CM); DeltaAccept emits the
//     judgment. ("Measure" in DESIGN.md is realized as the CC actor's α — the
//     kernel stays uniform; see ../DESIGN.md §Build notes.)
type Actor interface {
	// Alpha produces the cell's artifact and sets ArtifactPresent.
	Alpha(rc *RunContext) error
	// Beta reviews the artifact and sets ReviewApproved or ReviewRejected.
	Beta(rc *RunContext) error
	// Gamma closes: computes the receipt from transition evidence and sets
	// ReceiptBound.
	Gamma(rc *RunContext) error
	// Validate mechanically checks the artifact against the contract's
	// acceptance predicates and sets VPass or VFail (recording the verdict
	// into the receipt).
	Validate(rc *RunContext) error
	// DeltaAccept is δ on the pass path: records the accept boundary decision
	// and finalizes the receipt (and, for a CC, emits the judgment).
	DeltaAccept(rc *RunContext) error
	// DeltaHold is δ on the fail path: computes/finalizes a STOP receipt with a
	// reject boundary decision (fail closed).
	DeltaHold(rc *RunContext) error
}

// stepFuncs is the actor-step registry: it maps the step ids named in the
// transition table to Actor method calls. Like guardFuncs, it is a fixed
// generic vocabulary — Drive dispatches through it and never switches on a
// state name.
var stepFuncs = map[string]func(Actor, *RunContext) error{
	"alpha_produce": func(a Actor, rc *RunContext) error { return a.Alpha(rc) },
	"beta_review":   func(a Actor, rc *RunContext) error { return a.Beta(rc) },
	"gamma_close":   func(a Actor, rc *RunContext) error { return a.Gamma(rc) },
	"v_validate":    func(a Actor, rc *RunContext) error { return a.Validate(rc) },
	"delta_accept":  func(a Actor, rc *RunContext) error { return a.DeltaAccept(rc) },
	"delta_hold":    func(a Actor, rc *RunContext) error { return a.DeltaHold(rc) },
}
