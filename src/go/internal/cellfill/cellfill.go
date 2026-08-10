// Package cellfill defines fill-owned seat construction
// (msg-cn-pi-cnos-cds-fill-construction-51, operator-ratified).
//
// A seat is ONE tagged immutable value: its "fill" selects a constructor, and
// every other field in the same object is that constructor's argument. The
// generic runner's complete construction algorithm is
//
//	factory := registry.lookup(seat["fill"])
//	seatFn  := factory.construct(seat)
//
// and nothing more — the runner never learns that a fill needs cognition, a
// worktree, skills, a model, or a provider. A factory is registered already
// closed over whatever subsystems it composes; construction returns immutable
// values and starts no session.
package cellfill

import (
	"encoding/json"
	"fmt"

	"github.com/usurobor/cnos/src/go/internal/cellkernel"
)

// Constructed is what a factory returns: the seat function, the CANONICAL
// resolved declaration the closure records and digests (the factory owns its
// canonical form — deterministic bytes, holes resolved), and the truthful
// mode class of the work this seat produces.
type Constructed struct {
	Decl json.RawMessage
	Mode cellkernel.ExecutionMode
}

type ConstructedAlpha struct {
	Constructed
	Seat cellkernel.Alpha
}

type ConstructedBeta struct {
	Constructed
	Seat cellkernel.Beta
}

// AlphaFactory and BetaFactory construct one side from its complete resolved
// declaration. The declaration arrives strict-decoded by the factory itself:
// each fill owns the meaning AND the exact shape of its arguments.
type AlphaFactory func(decl json.RawMessage) (ConstructedAlpha, error)
type BetaFactory func(decl json.RawMessage) (ConstructedBeta, error)

// Registry is the small statically assembled fill map. No DI container, no
// reflection, no service locator — the assembly point lists its fills.
type Registry struct {
	Alpha map[string]AlphaFactory
	Beta  map[string]BetaFactory
}

// FillID extracts the tag that selects a constructor. It is the ONLY field
// the generic path reads from a seat declaration.
func FillID(decl json.RawMessage) (string, error) {
	var tag struct {
		Fill string `json:"fill"`
	}
	if err := json.Unmarshal(decl, &tag); err != nil {
		return "", fmt.Errorf("seat declaration is not an object: %w", err)
	}
	if tag.Fill == "" {
		return "", fmt.Errorf("seat declaration has no fill")
	}
	return tag.Fill, nil
}

// ConstructAlpha dispatches an alpha declaration. Unknown fills fail here,
// before any seat or provider is touched.
func (r Registry) ConstructAlpha(decl json.RawMessage) (ConstructedAlpha, error) {
	id, err := FillID(decl)
	if err != nil {
		return ConstructedAlpha{}, err
	}
	f, ok := r.Alpha[id]
	if !ok {
		return ConstructedAlpha{}, fmt.Errorf("unknown alpha fill %q", id)
	}
	return f(decl)
}

func (r Registry) ConstructBeta(decl json.RawMessage) (ConstructedBeta, error) {
	id, err := FillID(decl)
	if err != nil {
		return ConstructedBeta{}, err
	}
	f, ok := r.Beta[id]
	if !ok {
		return ConstructedBeta{}, fmt.Errorf("unknown beta fill %q", id)
	}
	return f(decl)
}

// CombineModes states the episode's truthful mode from its two seats: any
// stub work makes the whole episode non-authoritative; otherwise any rented
// seat makes it cognitive; otherwise it is mechanical and reproducible.
func CombineModes(a, b cellkernel.ExecutionMode) cellkernel.ExecutionMode {
	switch {
	case a == cellkernel.ModeStub || b == cellkernel.ModeStub:
		return cellkernel.ModeStub
	case a == cellkernel.ModeCognitive || b == cellkernel.ModeCognitive:
		return cellkernel.ModeCognitive
	default:
		return cellkernel.ModeMechanical
	}
}

// StrictDecode is the shared decode discipline for fill arguments: exact
// case-sensitive keys via DisallowUnknownFields over canonical lowercase
// structs, no trailing data. Fills use it so Go and CUE reject the same
// shapes. (Null and duplicate keys are rejected once, at the envelope walk.)
func StrictDecode(decl json.RawMessage, into any) error {
	dec := json.NewDecoder(bytesReader(decl))
	dec.DisallowUnknownFields()
	if err := dec.Decode(into); err != nil {
		return err
	}
	var extra json.RawMessage
	if dec.Decode(&extra) == nil {
		return fmt.Errorf("trailing data after seat declaration")
	}
	return nil
}
