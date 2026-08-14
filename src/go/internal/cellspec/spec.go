// Package cellspec loads a serialized cell spec (the JSON IR a compiled
// main.cell — or a hand-authored contract file — produces) and binds it to a
// runnable cellkernel.Spec through fill-owned seat construction
// (msg-cn-pi-cnos-cds-fill-construction-51, operator-ratified).
//
// The loader is semantically blind to seats: alpha and beta are each ONE
// tagged object whose "fill" selects a constructor from a statically
// assembled registry, and whose remaining fields belong to that constructor.
// This package validates the generic envelope strictly (Pi #32 D5, round-6
// D2: pinned version, no unknown/trailing/duplicate/case-aliased keys, no
// null anywhere) and resolves `$param` holes in place; every fill-specific
// rule lives with the fill and its CUE overlay.
//
// See docs/architecture/CDS-CELL-MIGRATION.md and schemas/cdd/spec.cue.
package cellspec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"sort"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
	"github.com/usurobor/cnos/src/go/internal/cellmethod"
)

// SchemaVersion is the pinned cell-spec version; a spec must declare it exactly.
const SchemaVersion = "cnos.cellspec.v0"

// CellSpec is the serialized, strictly-validated cell. Alpha and Beta are
// kept raw: their shape belongs to their fills.
//
// Methodology is raw for a different reason: its shape belongs to cellmethod,
// which is the one parser for it. This package would otherwise be a second
// place a methodology declaration is read, and two readers of one document are
// two chances to disagree about it (eng/go §2.17). It is OPTIONAL here — a
// cell may declare none, which is what every cell did before the bundle
// existed — and a fill that cannot act without one refuses at its own
// constructor, because only the fill knows that.
type CellSpec struct {
	Version     string               `json:"version"`
	Contract    ContractSpec         `json:"contract"`
	ProtocolID  string               `json:"protocol_id"`
	Params      map[string]ParamSpec `json:"params,omitempty"`
	Methodology json.RawMessage      `json:"methodology,omitempty"`
	Alpha       json.RawMessage      `json:"alpha"`
	Beta        json.RawMessage      `json:"beta"`
}

type ContractSpec struct {
	ID               string        `json:"id"`
	Goal             string        `json:"goal"`
	RequiredEvidence []RequiredRef `json:"required_evidence,omitempty"`
}

// RequiredRef names a required evidence ref and the producer role authorized
// to mint it (Pi D2).
type RequiredRef struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Producer string `json:"producer"`
}

// ParamSpec is a typed hole. There is deliberately no "kind": under
// fill-owned construction a hole is substituted as a string wherever it
// sits, and what the value MEANS is the fill's business — a skill ref is
// validated by loading it, not by a label in the generic envelope.
type ParamSpec struct {
	Required bool `json:"required"`
	// Pointer so a declared empty-string default ("model": "" for the fake
	// provider) is distinguishable from no default at all.
	Default *string  `json:"default,omitempty"`
	Domain  []string `json:"domain,omitempty"`
}

// Parse strictly decodes and validates the generic envelope of a serialized
// cell spec. Seat interiors are validated by their fills at Build time and by
// the fill's CUE overlay at vet time.
func Parse(data []byte) (CellSpec, error) {
	// Duplicate keys and nulls: the schema admits none, and a null collection
	// would silently decode to a Go nil (Pi round-6 D2). The walk lives in
	// cellfill because the run input needs the identical rule against the
	// identical CUE-side expectation.
	if err := cellfill.NoDuplicateKeysOrNull(data); err != nil {
		return CellSpec{}, fmt.Errorf("cellspec: %w", err)
	}
	if err := checkExactKeys(data); err != nil {
		return CellSpec{}, fmt.Errorf("cellspec: %w", err)
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var s CellSpec
	if err := dec.Decode(&s); err != nil {
		return CellSpec{}, fmt.Errorf("cellspec: decode: %w", err)
	}
	// Strict EOF: any non-whitespace byte after the first value — including a
	// stray delimiter like `]` or `}` — is rejected. dec.More() is an
	// iteration helper, not an EOF check, so decode a second value and
	// require io.EOF.
	var extra json.RawMessage
	if err := dec.Decode(&extra); err != io.EOF {
		return CellSpec{}, fmt.Errorf("cellspec: trailing data after JSON object")
	}

	if s.Version != SchemaVersion {
		return CellSpec{}, fmt.Errorf("cellspec: version must be %q, got %q", SchemaVersion, s.Version)
	}
	if s.Contract.ID == "" {
		return CellSpec{}, fmt.Errorf("cellspec: contract.id is required")
	}
	if s.Contract.Goal == "" { // parity with #CellSpec (Pi round-5 D2)
		return CellSpec{}, fmt.Errorf("cellspec: contract.goal is required")
	}
	if s.ProtocolID == "" { // opaque provenance (Pi PR-#718-fido β D6)
		return CellSpec{}, fmt.Errorf("cellspec: protocol_id is required")
	}
	for side, decl := range map[string]json.RawMessage{"alpha": s.Alpha, "beta": s.Beta} {
		if len(decl) == 0 {
			return CellSpec{}, fmt.Errorf("cellspec: %s is required", side)
		}
		if _, err := cellfill.FillID(decl); err != nil {
			return CellSpec{}, fmt.Errorf("cellspec: %s: %w", side, err)
		}
	}
	if err := validateEvidence(s.Contract.RequiredEvidence); err != nil {
		return CellSpec{}, fmt.Errorf("cellspec: %w", err)
	}
	return s, nil
}

func validateEvidence(refs []RequiredRef) error {
	seen := make(map[string]bool)
	for _, r := range refs {
		if r.ID == "" || r.Kind == "" {
			return fmt.Errorf("required_evidence entries need non-empty id and kind")
		}
		if r.Producer != string(cellkernel.RoleAlpha) && r.Producer != string(cellkernel.RoleBeta) {
			return fmt.Errorf("required_evidence %q has invalid producer %q (want alpha|beta)", r.ID, r.Producer)
		}
		if seen[r.ID] {
			return fmt.Errorf("duplicate required_evidence id %q", r.ID)
		}
		seen[r.ID] = true
	}
	return nil
}

// Resolved is a cell spec with its parameter holes filled in place — the seat
// declarations and the methodology are complete objects ready for their
// consumers.
type Resolved struct {
	Spec CellSpec
	// Methodology is empty when the spec declared none.
	Methodology json.RawMessage
	Alpha       json.RawMessage
	Beta        json.RawMessage
}

// Resolve fills `$param` holes from `given` (with defaults and closed
// domains) IN PLACE inside the methodology and both seat declarations: any
// string value equal to `$name` is replaced by the parameter's value, wherever
// it sits in the tree. Unresolved authored JSON carries holes in the same
// positions the resolved tree fills.
func (s CellSpec) Resolve(given map[string]string) (Resolved, error) {
	for name := range given {
		if _, ok := s.Params[name]; !ok {
			return Resolved{}, fmt.Errorf("unknown parameter %q (not declared by cell %q)", name, s.Contract.ID)
		}
	}
	vals := make(map[string]string)
	var missing []string
	for name, p := range s.Params {
		// PRESENCE, not emptiness (Pi #59 C2). `--param p=` supplies the empty
		// string deliberately — it is how a fake's meaningless model is
		// written — and testing `given[name] != ""` silently reclassified that
		// as "absent", so an explicit empty either picked up a default or was
		// reported missing. Whether empty is LEGAL is the declared domain's
		// question, and then the fill's; it is not this loop's.
		v, supplied := given[name]
		switch {
		case supplied:
			vals[name] = v
		case p.Default != nil:
			vals[name] = *p.Default
		case p.Required:
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return Resolved{}, fmt.Errorf("missing required parameter(s): %s (supply with --param <name>=<value>)", strings.Join(missing, ", "))
	}
	for name, p := range s.Params {
		if v, ok := vals[name]; ok && len(p.Domain) > 0 && !slices.Contains(p.Domain, v) {
			return Resolved{}, fmt.Errorf("parameter %q value %q not in domain [%s]", name, v, strings.Join(p.Domain, ", "))
		}
	}

	// The methodology carries holes for the same reason a seat does — a cell
	// declaring `$language` is one cell, not three — so it is spliced through
	// the identical walk. A second substitution rule for the same `$name`
	// spelling is how a hole comes to mean two things.
	var method json.RawMessage
	if len(s.Methodology) > 0 {
		m, err := splice(s.Methodology, s.Params, vals)
		if err != nil {
			return Resolved{}, fmt.Errorf("cellspec: methodology: %w", err)
		}
		method = m
	}
	alpha, err := splice(s.Alpha, s.Params, vals)
	if err != nil {
		return Resolved{}, fmt.Errorf("cellspec: alpha: %w", err)
	}
	beta, err := splice(s.Beta, s.Params, vals)
	if err != nil {
		return Resolved{}, fmt.Errorf("cellspec: beta: %w", err)
	}
	return Resolved{Spec: s, Methodology: method, Alpha: alpha, Beta: beta}, nil
}

// splice walks a seat's JSON tree replacing `$name` string values in place.
// A hole must reference a declared parameter; a declared-but-unfilled
// optional hole is an error at the point of use, not a silent empty string.
func splice(decl json.RawMessage, declared map[string]ParamSpec, vals map[string]string) (json.RawMessage, error) {
	var v any
	if err := json.Unmarshal(decl, &v); err != nil {
		return nil, err
	}
	filled, err := spliceValue(v, declared, vals)
	if err != nil {
		return nil, err
	}
	return json.Marshal(filled)
}

func spliceValue(v any, declared map[string]ParamSpec, vals map[string]string) (any, error) {
	switch t := v.(type) {
	case string:
		if !strings.HasPrefix(t, "$") {
			return t, nil
		}
		name := strings.TrimPrefix(t, "$")
		// A hole is MALFORMED or UNDECLARED, and the two are different facts
		// (Pi #59 C1). Checking only declaration made a malformed hole report
		// "undeclared" by accident — an illegal name cannot be declared, so
		// the wrong check happened to fire. It is the same predicate the
		// parameter declaration uses, because a hole spelling IS a name.
		if !validParamName(name) {
			return nil, fmt.Errorf("hole %q is malformed: %q is not a legal identifier "+
				"(letters, digits and underscore, not starting with a digit)", t, name)
		}
		if _, ok := declared[name]; !ok {
			return nil, fmt.Errorf("hole %q references undeclared parameter %q", t, name)
		}
		val, ok := vals[name]
		if !ok {
			return nil, fmt.Errorf("hole %q is unfilled (parameter %q has no value or default)", t, name)
		}
		return val, nil
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, e := range t {
			f, err := spliceValue(e, declared, vals)
			if err != nil {
				return nil, err
			}
			out[k] = f
		}
		return out, nil
	case []any:
		out := make([]any, 0, len(t))
		for _, e := range t {
			f, err := spliceValue(e, declared, vals)
			if err != nil {
				return nil, err
			}
			out = append(out, f)
		}
		return out, nil
	default:
		return v, nil
	}
}

// Binding is the per-run contract content a caller freezes into the cell: the
// admitted issue, the admitted design, and the PINNED subject.
//
// The three are `json.RawMessage` and stay that way. This package is the
// generic loader — it already refuses to learn what a fill means, and it
// refuses to learn what an issue or a repository is for the same reason. It
// carries the bytes to the one place the record digests them and reads none of
// them. A zero Binding is a cell with no run input, which is every cell that
// existed before the run input did.
//
// Admission, not this type, decides whether a given profile REQUIRES the three
// to be present: a generic loader that demanded an issue would have made the
// CDS profile's rule everyone's.
type Binding struct {
	Issue   json.RawMessage
	Design  json.RawMessage
	Subject json.RawMessage
}

// Build dispatches both seat declarations through the fill registry and binds
// the result to a runnable kernel Spec plus the RunMeta the record carries:
// the complete canonical resolved declarations and the truthful combined
// mode. The registry arrives from the assembly point (the CLI domain) — this
// package never names a fill.
//
// bind is spliced into the kernel contract and therefore into the canonical
// record bytes and the one scope-lift digest. It is a Build argument rather
// than something a caller sets on the returned Spec so that there is one
// construction point for the contract: a caller that assembled half of it here
// and half of it afterwards would have two places for the frozen value to
// come from.
func (r Resolved) Build(ctx context.Context, reg cellfill.Registry, bind Binding) (cellkernel.Spec, cellkernel.RunMeta, error) {
	// BOTH fills' DECLARED requirements are checked against the binding before
	// EITHER constructor runs. A missing subject is decided entirely by two
	// values both already in hand — the registration and `bind` — so nothing
	// about it needs a seat to exist. Constructing first made the refusal wait
	// for a provider adapter and every skill body to be built, and then arrive
	// from a station as an episode malfunction, which is not what a run with no
	// subject is. The assessing seat is checked here rather than at its own
	// constructor for the identical reason: a beta that cannot act without the
	// subject — `cds.assess` reconstructs the candidate from it — otherwise
	// discovered the fact from inside Review, one whole produced side later.
	// This package still learns nothing about fills: it reads one declared bool
	// per side and never asks what a subject is for.
	if err := r.checkDeclaredNeeds(reg, bind); err != nil {
		return cellkernel.Spec{}, cellkernel.RunMeta{}, err
	}

	constructive, adversarial, err := r.methodology(reg)
	if err != nil {
		return cellkernel.Spec{}, cellkernel.RunMeta{}, err
	}

	alpha, err := reg.ConstructAlpha(ctx, r.Alpha, constructive)
	if err != nil {
		return cellkernel.Spec{}, cellkernel.RunMeta{}, fmt.Errorf("cellspec: alpha: %w", err)
	}
	beta, err := reg.ConstructBeta(ctx, r.Beta, adversarial)
	if err != nil {
		return cellkernel.Spec{}, cellkernel.RunMeta{}, fmt.Errorf("cellspec: beta: %w", err)
	}

	meta := cellkernel.RunMeta{
		ExecutionMode: cellfill.CombineModes(alpha.Mode, beta.Mode),
		ResolvedSpec: cellkernel.ResolvedSpec{
			Version:          r.Spec.Version,
			DeclaredProtocol: r.Spec.ProtocolID,
			Alpha:            alpha.Decl,
			Beta:             beta.Decl,
		},
	}
	return cellkernel.Spec{Contract: r.contract(bind), Alpha: alpha.Seat, Beta: beta.Seat}, meta, nil
}

// checkDeclaredNeeds reads each side's DECLARED requirements and refuses a
// binding that cannot satisfy them — before either constructor runs.
//
// One rule, applied to each side in turn, alpha first: a cell whose two seats
// both need a subject refuses with the producing side named, as it did before
// the assessing side could declare anything.
func (r Resolved) checkDeclaredNeeds(reg cellfill.Registry, bind Binding) error {
	alphaID, alphaFill, err := reg.LookupAlpha(r.Alpha)
	if err != nil {
		return fmt.Errorf("cellspec: alpha: %w", err)
	}
	betaID, betaFill, err := reg.LookupBeta(r.Beta)
	if err != nil {
		return fmt.Errorf("cellspec: beta: %w", err)
	}
	for _, side := range []struct {
		role         cellkernel.Role
		id           string
		needsSubject bool
	}{
		{cellkernel.RoleAlpha, alphaID, alphaFill.NeedsSubject},
		{cellkernel.RoleBeta, betaID, betaFill.NeedsSubject},
	} {
		if side.needsSubject && len(bind.Subject) == 0 {
			return fmt.Errorf("cellspec: %s fill %q requires contract.subject, and no run input supplied one "+
				"(pass --input)", side.role, side.id)
		}
	}
	return nil
}

// methodology loads the declared bundle ONCE and returns the two projections.
//
// Once, here, before either seat exists: a declared methodology is loaded
// whatever the fills are, so a cell whose bundle names an uninstalled skill
// fails at the cell, not later and only if some seat happened to ask.
//
// BOTH projections come off the SAME load. Loading twice — once per seat —
// would put two reads of one declaration on the two sides of the episode, which
// is the drift the single bundle exists to remove: the producing seat and the
// assessing seat would then be held to two collections that nothing could tell
// apart if a skill body changed between the reads. Two projections of one
// digest cannot be two methodologies.
func (r Resolved) methodology(reg cellfill.Registry) (constructive, adversarial cellmethod.View, err error) {
	if len(r.Methodology) == 0 {
		return cellmethod.View{}, cellmethod.View{}, nil
	}
	bundle, bodies, err := cellmethod.Load(reg.Skills, r.Methodology)
	if err != nil {
		return cellmethod.View{}, cellmethod.View{}, fmt.Errorf("cellspec: %w", err)
	}
	return cellmethod.Constructive(bundle, bodies), cellmethod.Adversarial(bundle, bodies), nil
}

// contract freezes the kernel contract: the cell's own fields, plus the three
// opaque slots exactly as admission pinned them.
func (r Resolved) contract(bind Binding) cellkernel.Contract {
	req := make([]cellkernel.RequiredRef, 0, len(r.Spec.Contract.RequiredEvidence))
	for _, e := range r.Spec.Contract.RequiredEvidence {
		req = append(req, cellkernel.RequiredRef{ID: e.ID, Kind: e.Kind, Producer: cellkernel.Role(e.Producer)})
	}
	return cellkernel.Contract{
		ID:               r.Spec.Contract.ID,
		Goal:             r.Spec.Contract.Goal,
		Issue:            bind.Issue,
		Design:           bind.Design,
		Subject:          bind.Subject,
		RequiredEvidence: req,
	}
}

// checkExactKeys walks the known GENERIC object shapes and requires every key
// to be in the shape's exact (case-sensitive) set — closing encoding/json's
// case-insensitive field matching (Pi round-6 D2). Seat interiors are fill
// territory: each fill's StrictDecode enforces its own exact shape, so only
// the envelope shapes are checked here. Not a schema engine.
func checkExactKeys(data []byte) error {
	var root map[string]json.RawMessage
	if err := json.Unmarshal(data, &root); err != nil {
		return err
	}
	if err := keysIn("spec", root, "version", "contract", "protocol_id", "params", "methodology", "alpha", "beta"); err != nil {
		return err
	}
	if err := objectKeys(root["contract"], "contract", "id", "goal", "required_evidence"); err != nil {
		return err
	}
	var contract struct {
		RequiredEvidence []json.RawMessage `json:"required_evidence"`
	}
	if len(root["contract"]) > 0 {
		if err := json.Unmarshal(root["contract"], &contract); err != nil {
			return err
		}
	}
	for _, ref := range contract.RequiredEvidence {
		if err := objectKeys(ref, "required_evidence entry", "id", "kind", "producer"); err != nil {
			return err
		}
	}
	var params map[string]json.RawMessage
	if len(root["params"]) > 0 {
		if err := json.Unmarshal(root["params"], &params); err != nil {
			return err
		}
	}
	for name, p := range params {
		if !validParamName(name) {
			return fmt.Errorf("parameter %q is not a legal identifier "+
				"(letters, digits and underscore, not starting with a digit)", name)
		}
		if err := objectKeys(p, fmt.Sprintf("parameter %q", name), "required", "default", "domain"); err != nil {
			return err
		}
	}
	return nil
}

// validParamName is the ONE identifier grammar, mirrored exactly by
// `#ParamName` in schemas/cdd/spec.cue. A name is also a hole spelling —
// `$name` — so a name legal here and illegal there would resolve in Go and be
// rejected by CUE, which is the divergence this closes (Pi #55 C1). Written
// out rather than compiled as a regexp: it is four comparisons, and the
// authority it mirrors is a pattern, not a library.
func validParamName(name string) bool {
	if name == "" {
		return false
	}
	for i, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r == '_':
		case i > 0 && r >= '0' && r <= '9':
		default:
			return false
		}
	}
	return true
}

// objectKeys checks one raw JSON object's keys against an exact allowed set.
// Absent or non-object values are left for the strict decode to judge.
func objectKeys(raw json.RawMessage, where string, allowed ...string) error {
	if len(raw) == 0 || raw[0] != '{' {
		return nil
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return err
	}
	return keysIn(where, m, allowed...)
}

func keysIn(where string, m map[string]json.RawMessage, allowed ...string) error {
	for k := range m {
		if !slices.Contains(allowed, k) {
			return fmt.Errorf("%s has unknown key %q (keys are exact and case-sensitive)", where, k)
		}
	}
	return nil
}
