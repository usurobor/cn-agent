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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

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

// AlphaFactory and BetaFactory construct one side from its complete
// declaration. Each fill owns the meaning AND the exact shape of its
// arguments. Construction may do bounded IO — loading skill bodies, pinning a
// revision — so it takes a context; it must not start or retain a session.
type AlphaFactory func(ctx context.Context, decl json.RawMessage) (ConstructedAlpha, error)
type BetaFactory func(ctx context.Context, decl json.RawMessage) (ConstructedBeta, error)

// Registry is the small statically assembled fill map. No DI container, no
// reflection, no service locator — the assembly point lists its fills.
type Registry struct {
	Alpha map[string]AlphaFactory
	Beta  map[string]BetaFactory
}

// FillID extracts the tag that selects a constructor. It is the ONLY field
// the generic path reads from a seat declaration.
//
// The lookup is by exact key: decoding into a struct would let `Fill` or
// `FILL` select a constructor that the closed CUE overlay rejects.
func FillID(decl json.RawMessage) (string, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(decl, &obj); err != nil {
		return "", fmt.Errorf("seat declaration is not an object: %w", err)
	}
	raw, ok := obj["fill"]
	if !ok {
		return "", fmt.Errorf("seat declaration has no fill")
	}
	var id string
	if err := json.Unmarshal(raw, &id); err != nil || id == "" {
		return "", fmt.Errorf("seat declaration has no fill")
	}
	return id, nil
}

// ConstructAlpha dispatches an alpha declaration. Unknown fills fail here,
// before any seat or provider is touched, and the returned declaration is
// canonicalized centrally so no fill can make the record's digest depend on
// how it happened to serialize.
func (r Registry) ConstructAlpha(ctx context.Context, decl json.RawMessage) (ConstructedAlpha, error) {
	id, err := FillID(decl)
	if err != nil {
		return ConstructedAlpha{}, err
	}
	f, ok := r.Alpha[id]
	if !ok {
		return ConstructedAlpha{}, fmt.Errorf("unknown alpha fill %q", id)
	}
	c, err := f(ctx, decl)
	if err != nil {
		return ConstructedAlpha{}, err
	}
	c.Decl, err = canonical(id, c.Decl)
	return c, err
}

func (r Registry) ConstructBeta(ctx context.Context, decl json.RawMessage) (ConstructedBeta, error) {
	id, err := FillID(decl)
	if err != nil {
		return ConstructedBeta{}, err
	}
	f, ok := r.Beta[id]
	if !ok {
		return ConstructedBeta{}, fmt.Errorf("unknown beta fill %q", id)
	}
	c, err := f(ctx, decl)
	if err != nil {
		return ConstructedBeta{}, err
	}
	c.Decl, err = canonical(id, c.Decl)
	return c, err
}

// canonical re-serializes a fill's resolved declaration into one canonical
// form: object keys sorted, no insignificant whitespace, numeric literals
// preserved exactly. The scope-lift digest covers these bytes, so making the
// runtime own the form means a fill cannot destabilize a digest by changing
// how it marshals — and reordering a Go struct's fields never moves a digest.
func canonical(fill string, raw json.RawMessage) (json.RawMessage, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, fmt.Errorf("fill %q returned a non-JSON declaration: %w", fill, err)
	}
	if _, ok := v.(map[string]any); !ok {
		return nil, fmt.Errorf("fill %q returned a declaration that is not an object", fill)
	}
	out, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("fill %q: canonicalize: %w", fill, err)
	}
	return out, nil
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

// StrictDecode is the shared decode discipline for fill arguments: every key
// must match a declared json tag EXACTLY, and there must be no trailing data.
//
// DisallowUnknownFields alone is not enough. encoding/json matches field names
// case-insensitively, so `Fill`, `COGNITION`, or a nested `Provider` would
// decode happily in Go while the closed CUE overlay rejects them — the two
// authorities would accept different languages. exactKeys closes that by
// walking the declaration against the target's declared tags, so a fill's Go
// shape and its CUE shape admit the same keys by construction rather than by
// a hand-maintained list per fill.
func StrictDecode(decl json.RawMessage, into any) error {
	if err := exactKeys(decl, reflect.TypeOf(into), "seat declaration"); err != nil {
		return err
	}
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

// exactKeys checks one JSON value's object keys against the json tags of the
// Go type it will decode into, recursing through nested structs and slices.
// It validates KEYS only — types and values remain the decoder's and the
// fill's business, so this stays a boundary check rather than a schema engine.
func exactKeys(raw json.RawMessage, t reflect.Type, where string) error {
	for t != nil && (t.Kind() == reflect.Pointer || t.Kind() == reflect.Interface) {
		t = t.Elem()
	}
	if t == nil {
		return nil
	}
	switch t.Kind() {
	case reflect.Struct:
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(raw, &obj); err != nil {
			return nil // not an object here; the decoder reports the type error
		}
		allowed := make(map[string]reflect.Type, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
			if name == "" {
				name = f.Name
			}
			if name == "-" {
				continue
			}
			allowed[name] = f.Type
		}
		for k, v := range obj {
			ft, ok := allowed[k]
			if !ok {
				return fmt.Errorf("%s has unknown key %q (keys are exact and case-sensitive)", where, k)
			}
			if err := exactKeys(v, ft, where+"."+k); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		var items []json.RawMessage
		if err := json.Unmarshal(raw, &items); err != nil {
			return nil
		}
		for i, item := range items {
			if err := exactKeys(item, t.Elem(), fmt.Sprintf("%s[%d]", where, i)); err != nil {
				return err
			}
		}
	}
	return nil
}
