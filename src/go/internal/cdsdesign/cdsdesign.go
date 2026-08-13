// Package cdsdesign owns the CDS typed design: its shape and its admission
// predicate. Nothing else.
//
// The design is the issue's counterpart and is deliberately a SEPARATE
// document: an issue that could also state architecture would let the problem
// be redefined by the same act that proposes the change, which is exactly the
// collapse CELL-SYSTEM-DESIGN §10.1 forbids. The two travel together in one
// run input, are admitted separately, and are frozen separately into the
// contract.
//
// The shape is the minimum decidable one — approach, invariants, impact — and
// nothing beyond it. Alternatives, trade-offs, authority changes and known
// debt are real parts of a design as §10.1 describes it, and every one of them
// is prose a structural gate cannot decide; adding fields a gate cannot judge
// would grow the authoring burden without growing what admission proves.
//
// Pure: no IO, no paths. Admit takes bytes.
package cdsdesign

import (
	"encoding/json"
	"fmt"

	"github.com/usurobor/cnos/src/go/internal/cdsissue"
	"github.com/usurobor/cnos/src/go/internal/cellfill"
)

// Kind is the pinned design schema tag. A design must declare it exactly; the
// CUE definition #CDSDesign pins the same string. It follows the repo's
// `cnos.<domain>.<name>.vN` convention, as its sibling cnos.cds.issue.v0 does.
const Kind = "cnos.cds.design.v0"

// Design is the proposed system change in the three parts a structural gate
// can actually decide.
type Design struct {
	Kind string `json:"kind"`
	// Approach is how the change is made, in prose. One field rather than a
	// decomposition, because any finer structure here would be a shape the
	// gate cannot check.
	Approach string `json:"approach"`
	// Invariants is what must stay true after the change. At least one: a
	// design that names nothing it must not break has not been thought about,
	// and there is no such change in a system that already has invariants.
	Invariants []string `json:"invariants"`
	// Impact is the surfaces the change reaches. At least one, because a
	// design that touches nothing is not a design.
	Impact []Surface `json:"impact"`
}

// Surface binds one reached surface — a path or a component name — to why the
// change reaches it. The `why` is required for the same reason a criterion's
// verification route is: a named surface with no stated reason leaves a later
// reader guessing whether it belongs in the change at all.
type Surface struct {
	Surface string `json:"surface"`
	Why     string `json:"why"`
}

// Admit decodes and admits a design document. Every rule is mechanical.
//
// Blankness is cdsissue.Blank — IMPORTED, not re-transcribed. The predicate is
// an enumeration of unicode.IsSpace written out by hand and mirrored in CUE's
// #NonBlank; a second copy here would be a third hand-written enumeration to
// keep in step, and the divergence that produced the first one (RE2's `\s` vs
// strings.TrimSpace) is documented at its definition.
//
// The closed key language is stated with cellfill.OnlyKeys at every level, not
// left to the decoder: encoding/json matches field names case-insensitively
// even with DisallowUnknownFields, so `Kind` or a nested `Why` would decode in
// Go while the closed #CDSDesign rejects it — the two authorities would admit
// different documents.
func Admit(raw []byte) (Design, error) {
	if len(raw) == 0 {
		return Design{}, fmt.Errorf("cds design: run input carries no design")
	}
	if err := exactShape(raw); err != nil {
		return Design{}, fmt.Errorf("cds design: %w", err)
	}

	var d Design
	if err := cellfill.StrictDecode(raw, &d); err != nil {
		return Design{}, fmt.Errorf("cds design: %w", err)
	}

	if d.Kind != Kind {
		return Design{}, fmt.Errorf("cds design: kind must be %q, got %q", Kind, d.Kind)
	}
	if cdsissue.Blank(d.Approach) {
		return Design{}, fmt.Errorf("cds design: approach is required")
	}
	if len(d.Invariants) == 0 {
		return Design{}, fmt.Errorf("cds design: invariants is required (at least one thing the change must not break)")
	}
	for i, inv := range d.Invariants {
		if cdsissue.Blank(inv) {
			return Design{}, fmt.Errorf("cds design: invariants[%d] is blank", i)
		}
	}
	if len(d.Impact) == 0 {
		return Design{}, fmt.Errorf("cds design: impact is required (at least one surface the change reaches)")
	}
	for i, s := range d.Impact {
		if cdsissue.Blank(s.Surface) {
			return Design{}, fmt.Errorf("cds design: impact[%d] names no surface", i)
		}
		if cdsissue.Blank(s.Why) {
			return Design{}, fmt.Errorf("cds design: impact %q states no reason", s.Surface)
		}
	}
	return d, nil
}

// exactShape states the closed key language at every object level. Unlike the
// issue's, it has no present-vs-absent case to read off the raw document: every
// field here is required and non-empty, so an absent key and an empty one are
// refused by the same rule.
func exactShape(raw json.RawMessage) error {
	if err := cellfill.OnlyKeys(raw, "cds design", "kind", "approach", "invariants", "impact"); err != nil {
		return err
	}
	impact, ok := cellfill.Field(raw, "impact")
	if !ok {
		return nil
	}
	var items []json.RawMessage
	if err := json.Unmarshal(impact, &items); err != nil {
		// Not an array: the strict decode below reports the type error. This
		// function reports on the key language only.
		return nil
	}
	for _, item := range items {
		if len(item) == 0 || item[0] != '{' {
			continue
		}
		if err := cellfill.OnlyKeys(item, "cds design impact", "surface", "why"); err != nil {
			return err
		}
	}
	return nil
}
