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
//
// # Error prefixes across the cell packages
//
// One identity, stated once, where the error LEAVES the package that decided
// it: the package's own name (`cellwork:`, `cdsissue:`), or the fill id
// (`cds.patch:`) when the error is about a seat the operator declared in cell
// JSON. Everything inside names the operation or the position and repeats no
// identity — which is why nothing in THIS package prefixes at all: cellfill's
// errors are always wrapped by the fill that called it, and that fill is the
// identity an operator can act on.
//
// There were five spellings before, and `cds issue: cds issue problem has
// unknown key "extra"` is what they cost a reader.
package cellfill

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/usurobor/cnos/src/go/internal/cellcog"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellmethod"
	"github.com/usurobor/cnos/src/go/internal/cellskill"
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
// arguments. Construction may do bounded IO — pinning a revision — so it takes
// a context; it must not start or retain a session.
//
// EACH SIDE ADDITIONALLY RECEIVES THE CELL'S METHODOLOGY PROJECTION for its
// role: alpha the constructive one, beta the adversarial one. It is a separate
// parameter rather than a field of the declaration because the methodology is
// not the seat's to declare: it is declared once on the cell, loaded once, and
// projected — so a fill can consume obligations but cannot choose them, and two
// seats cannot be held to two lists. A fill that needs no methodology ignores
// it; a fill that cannot act without one says so itself, since only it knows
// that.
//
// Beta gained the parameter when the assessing fill landed, and not before: a
// parameter every beta ignored would have been a claim that something reads it.
type AlphaFactory func(ctx context.Context, decl json.RawMessage, method cellmethod.View) (ConstructedAlpha, error)
type BetaFactory func(ctx context.Context, decl json.RawMessage, method cellmethod.View) (ConstructedBeta, error)

// Admitted is the per-run contract a door let through: the exact authored bytes
// of each payload, in the form the cell will freeze them. Raw, because what an
// issue or a repository reference MEANS belongs to the profile that admitted
// it — this package carries the bytes and reads none of them, exactly as it
// carries seat declarations without learning what a fill needs.
type Admitted struct {
	Issue   json.RawMessage
	Design  json.RawMessage
	Subject json.RawMessage
}

// ErrRefused is what every non-admitting outcome wraps. It lives here rather
// than in a profile package because it is part of the Door contract: the
// generic runner has to tell "this document was judged and refused" from "the
// door malfunctioned", and it must do so without importing the profile whose
// door it happens to be dispatching.
var ErrRefused = errors.New("run input refused")

// Door decides whether one run-input document is executable, and returns the
// typed receipt of that decision ALREADY SERIALIZED. The receipt is opaque
// bytes for the same reason a seat declaration is: its shape, its kind tag and
// its vocabulary belong to the profile, and a runner that decoded it would have
// learned a profile's language in order to print it.
//
// It takes the raw document, not a decoded envelope. The envelope's kind and
// key language are as much a profile's rule as its payload shapes, so a
// document with the wrong kind is a REFUSAL with a receipt — not a read error
// the caller reports some other way. One decision, one channel.
type Door func(raw []byte) (Admitted, json.RawMessage, error)

// AlphaFill is one registered alpha: its constructor, and the per-run inputs
// the fill cannot act without.
//
// The requirement is declared BESIDE the constructor, for the same reason a
// fill's arguments sit beside its `fill` tag: the fill is the only thing that
// knows it. Registering a bare factory forced the fact to be discovered by
// running the seat, so a decisive inadmissibility — no subject at all —
// surfaced only after the constructor had built its provider adapter and read
// every skill body, and then surfaced as a station malfunction. Declaring it
// here lets the loader refuse before construction without the generic runner
// learning what any fill needs.
type AlphaFill struct {
	Construct AlphaFactory
	// NeedsSubject: this fill cannot act without contract.subject. It says
	// nothing about issue or design, and nothing about what a subject MEANS —
	// that stays the profile's.
	NeedsSubject bool
}

// BetaFill is one registered beta, and carries the SAME declaration for the
// same reason: an assessing seat can be as unable to act without the pinned
// subject as a producing one — `cds.assess` reconstructs the candidate from
// it — and a requirement the registry cannot state is a requirement discovered
// by running the seat. Undeclared, it surfaced from inside Review as an
// episode malfunction; declared, the loader refuses the run before the
// constructor.
//
// The two sides are separate types rather than one because their constructors
// are separate types: a shared struct would need a factory field wide enough
// for both, which is exactly the `any` this package refuses to hold.
type BetaFill struct {
	Construct BetaFactory
	// NeedsSubject: as AlphaFill.NeedsSubject, on the assessing side.
	NeedsSubject bool
}

// Registry is the small statically assembled fill map. No DI container, no
// service locator — the assembly point lists its fills.
//
// Both sides map a fill id to a fill: a constructor and what that fill cannot
// act without. They are symmetric because the loader's question is symmetric —
// "can this seat act on this run input?" is asked of the assessing seat exactly
// as it is asked of the producing one.
//
// Door sits beside them because it is the same kind of thing: a profile-owned
// function the composition root wires in and the generic runner only
// dispatches. Without it the runner had to name a domain package to admit a
// document, which is the coupling the fill registry exists to prevent. A nil
// Door is a registry that admits no run input — legitimate, and the shape every
// cell had before run inputs existed.
// Skills is the ONE skill authority this binary was assembled with. It sits
// here for the same reason Door does — a profile-owned dependency the
// composition root wires in and the generic path only hands onward — and it is
// on the registry rather than inside a fill because the methodology is the
// CELL's: a resolver held privately by one fill would be a second place skills
// could come from, which is exactly the drift the single bundle removes. A nil
// resolver is a registry that can load no methodology, which is legitimate for
// every cell that declares none.
type Registry struct {
	Alpha  map[string]AlphaFill
	Beta   map[string]BetaFill
	Door   Door
	Skills cellskill.Resolver
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

// LookupAlpha resolves a seat declaration to the fill it selects, WITHOUT
// constructing it. A caller that must judge a fill's declared requirements
// before its constructor runs looks it up here; ConstructAlpha resolves through
// the same call, so there is one place a fill id becomes a fill and no second
// map to keep in step.
func (r Registry) LookupAlpha(decl json.RawMessage) (string, AlphaFill, error) {
	id, err := FillID(decl)
	if err != nil {
		return "", AlphaFill{}, err
	}
	f, ok := r.Alpha[id]
	if !ok {
		return "", AlphaFill{}, fmt.Errorf("unknown alpha fill %q", id)
	}
	return id, f, nil
}

// ConstructAlpha dispatches an alpha declaration. Unknown fills fail here,
// before any seat or provider is touched, and the returned declaration is
// canonicalized centrally so no fill can make the record's digest depend on
// how it happened to serialize.
func (r Registry) ConstructAlpha(ctx context.Context, decl json.RawMessage, method cellmethod.View) (ConstructedAlpha, error) {
	id, f, err := r.LookupAlpha(decl)
	if err != nil {
		return ConstructedAlpha{}, err
	}
	c, err := f.Construct(ctx, decl, method)
	if err != nil {
		return ConstructedAlpha{}, err
	}
	c.Decl, err = canonical(id, c.Decl)
	return c, err
}

// LookupBeta is LookupAlpha on the assessing side, and exists for the same
// reason: a caller that must judge a fill's declared requirements before its
// constructor runs needs the fill without the seat. ConstructBeta resolves
// through it, so a fill id becomes a fill in ONE place on this side too.
func (r Registry) LookupBeta(decl json.RawMessage) (string, BetaFill, error) {
	id, err := FillID(decl)
	if err != nil {
		return "", BetaFill{}, err
	}
	f, ok := r.Beta[id]
	if !ok {
		return "", BetaFill{}, fmt.Errorf("unknown beta fill %q", id)
	}
	return id, f, nil
}

// ConstructBeta dispatches a beta declaration, mirroring ConstructAlpha:
// unknown fills fail here, before any seat is touched, and the returned
// declaration is canonicalized centrally.
func (r Registry) ConstructBeta(ctx context.Context, decl json.RawMessage, method cellmethod.View) (ConstructedBeta, error) {
	id, f, err := r.LookupBeta(decl)
	if err != nil {
		return ConstructedBeta{}, err
	}
	c, err := f.Construct(ctx, decl, method)
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
	// The same EOF requirement StrictDecode carries, for the same reason: one
	// value decoded is not one value SUPPLIED. Without this, a fill returning
	// `{...} garbage` is silently canonicalized down to the leading object and
	// the trailing bytes vanish from the record (Pi #57 B1). The built-ins
	// marshal clean JSON today, so this is local robustness at a boundary that
	// must not quietly discard what it was handed.
	var extra json.RawMessage
	if err := dec.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("fill %q returned trailing data after its declaration", fill)
		}
		return nil, fmt.Errorf("fill %q returned malformed data after its declaration: %w", fill, err)
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

// StrictDecode decodes a fill's arguments with no unknown fields and no
// trailing data. It does NOT derive a key language: deriving every fill's
// accepted shape here would make the generic package learn fill semantics
// through a small schema engine — the leak the one-place fill design removed.
//
// encoding/json matches field names case-insensitively even with
// DisallowUnknownFields, so each fill pairs this with OnlyKeys over the exact
// keys ITS closed shape declares.
func StrictDecode(decl json.RawMessage, into any) error {
	dec := json.NewDecoder(bytesReader(decl))
	dec.DisallowUnknownFields()
	if err := dec.Decode(into); err != nil {
		return err
	}
	var extra json.RawMessage
	if err := dec.Decode(&extra); err != io.EOF {
		// Only a real EOF proves the declaration ended. Accepting ANY error
		// here would report "clean" for a stream that failed for some other
		// reason, which is the claim this check exists to make (Pi #55 C2).
		if err == nil {
			return fmt.Errorf("trailing data after seat declaration")
		}
		return fmt.Errorf("malformed data after seat declaration: %w", err)
	}
	return nil
}

// OnlyKeys requires one JSON object to carry no key outside `allowed`, exactly
// and case-sensitively. It is a five-line helper a fill applies to each of its
// own closed shapes — not a validator generator: types and values stay the
// decoder's and the fill's business.
func OnlyKeys(raw json.RawMessage, where string, allowed ...string) error {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return fmt.Errorf("%s is not an object: %w", where, err)
	}
	for k := range obj {
		if !slices.Contains(allowed, k) {
			return fmt.Errorf("%s has unknown key %q (keys are exact and case-sensitive)", where, k)
		}
	}
	return nil
}

// SeatDecl is the closed {fill, cognition} declaration: a seat that states its
// tag, rents cognition, and declares nothing else. Both CDS seats speak exactly
// this shape — no workspace, no skills, no model list of their own — so it is
// declared ONCE here rather than copied per fill. A second copy would be a
// second key language, and a rule added to one copy is a rule missing from the
// other, which is how the fills diverged before.
//
// It does not make this package a fill framework and does not put cognition on
// the dispatch path: nothing in Registry, FillID or ConstructAlpha reads it.
// It is a decoder a fill CALLS, in the same way the cdd fills in this package
// call OnlyKeys — the runner still learns nothing about what any fill needs.
type SeatDecl struct {
	Fill      string         `json:"fill"`
	Cognition cellcog.Config `json:"cognition"`
}

// SeatRefusal is what a fill says about ITSELF when it cannot be constructed.
// The helper decides WHEN each refusal applies, because that rule is the same
// for every seat of this shape; the fill owns the words, because "a patch
// alpha" and "an assessing beta" are each fill's own account of what it is and
// no generic decoder can write them.
type SeatRefusal struct {
	// NoMethodology completes `fill "x": ` when the cell declares no
	// methodology at all — a whole sentence naming what is missing.
	NoMethodology string
	// WrongRole completes `fill "x": ` when the projection carries the other
	// role. The role actually received is appended as `, got "..."`, so the
	// sentence stops before its comma.
	WrongRole string
}

// AdmitSeatDecl decodes one {fill, cognition} declaration under its closed key
// language and checks that the projection handed to the seat is the one that
// seat takes. Everything it refuses is wrapped `fill %q: `, so a caller returns
// its zero Constructed value and this error unchanged.
//
// The key language is stated here explicitly, at both of its object shapes,
// because encoding/json matches field names case-insensitively even with
// DisallowUnknownFields: `Fill`, `Cognition` or a nested `Provider` would
// otherwise decode in Go while the closed CUE overlays reject them. `workspace`
// and `skills` are absent from the allowed sets, and those absences are
// load-bearing — a declaration naming a repository or a methodology of its own
// is refused by name here, so a second source of either cannot come back as a
// tolerated extra key.
//
// An empty projection and the wrong projection are refused HERE but worded by
// the fill: whether a seat can act with no methodology is the fill's own
// knowledge — a cell declaring none is legitimate for fills that need none —
// and this helper is called only by fills that have already decided they need
// one.
func AdmitSeatDecl(decl json.RawMessage, fill string, want cellmethod.Role,
	method cellmethod.View, refusal SeatRefusal) (SeatDecl, error) {
	if err := OnlyKeys(decl, fill, "fill", "cognition"); err != nil {
		return SeatDecl{}, fmt.Errorf("fill %q: %w", fill, err)
	}
	if cog, ok := Field(decl, "cognition"); ok {
		if err := OnlyKeys(cog, fill+".cognition", "provider", "model"); err != nil {
			return SeatDecl{}, fmt.Errorf("fill %q: %w", fill, err)
		}
	}
	var d SeatDecl
	if err := StrictDecode(decl, &d); err != nil {
		return SeatDecl{}, fmt.Errorf("fill %q: %w", fill, err)
	}
	if method.Empty() {
		return SeatDecl{}, fmt.Errorf("fill %q: %s", fill, refusal.NoMethodology)
	}
	if method.Role != want {
		// A seat handed the other role's projection would be held to the
		// opposite obligation — a producer told to falsify its own work, or a
		// reviewer told to follow the obligations it is meant to test. Nothing
		// constructs either that way today; this refuses the wiring mistake
		// rather than trusting that.
		return SeatDecl{}, fmt.Errorf("fill %q: %s, got %q", fill, refusal.WrongRole, method.Role)
	}
	return d, nil
}

// NoDuplicateKeysOrNull rejects duplicate object keys anywhere in a JSON
// document, which encoding/json otherwise silently accepts (last-wins), and
// JSON null anywhere.
//
// It lives here, beside OnlyKeys and StrictDecode, because it belongs to the
// same fact: what encoding/json quietly tolerates and a closed CUE definition
// does not. Two documents that both go to CUE — the cell spec and the run
// input — must be read the same way by both authorities, and two copies of
// this walk would be two chances to drift (eng/go §2.17, one parser per fact).
func NoDuplicateKeysOrNull(data []byte) error {
	return walkNoDup(json.NewDecoder(bytes.NewReader(data)))
}

func walkNoDup(dec *json.Decoder) error {
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if t == nil {
		return fmt.Errorf("null is not allowed (the schema admits no null)")
	}
	delim, ok := t.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		keys := make(map[string]bool)
		for dec.More() {
			kt, err := dec.Token()
			if err != nil {
				return err
			}
			key, _ := kt.(string)
			if keys[key] {
				return fmt.Errorf("duplicate key %q", key)
			}
			keys[key] = true
			if err := walkNoDup(dec); err != nil {
				return err
			}
		}
		if _, err := dec.Token(); err != nil { // consume '}'
			return err
		}
	case '[':
		for dec.More() {
			if err := walkNoDup(dec); err != nil {
				return err
			}
		}
		if _, err := dec.Token(); err != nil { // consume ']'
			return err
		}
	}
	return nil
}

// Field returns one raw member of a JSON object, if present.
func Field(raw json.RawMessage, name string) (json.RawMessage, bool) {
	var obj map[string]json.RawMessage
	if json.Unmarshal(raw, &obj) != nil {
		return nil, false
	}
	v, ok := obj[name]
	return v, ok
}
