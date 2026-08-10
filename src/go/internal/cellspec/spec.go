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
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"sort"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cellfill"
	"github.com/usurobor/cnos/src/go/internal/cellkernel"
)

// SchemaVersion is the pinned cell-spec version; a spec must declare it exactly.
const SchemaVersion = "cnos.cellspec.v0"

// CellSpec is the serialized, strictly-validated cell. Alpha and Beta are
// kept raw: their shape belongs to their fills.
type CellSpec struct {
	Version    string               `json:"version"`
	Contract   ContractSpec         `json:"contract"`
	ProtocolID string               `json:"protocol_id"`
	Params     map[string]ParamSpec `json:"params,omitempty"`
	Alpha      json.RawMessage      `json:"alpha"`
	Beta       json.RawMessage      `json:"beta"`
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

type ParamSpec struct {
	Kind     string `json:"kind"`
	Required bool   `json:"required"`
	// Pointer so a declared empty-string default ("model": "" for the fake
	// provider) is distinguishable from no default at all.
	Default *string  `json:"default,omitempty"`
	Domain  []string `json:"domain,omitempty"`
}

// Parse strictly decodes and validates the generic envelope of a serialized
// cell spec. Seat interiors are validated by their fills at Build time and by
// the fill's CUE overlay at vet time.
func Parse(data []byte) (CellSpec, error) {
	if err := checkNoDuplicateKeysOrNull(data); err != nil {
		return CellSpec{}, fmt.Errorf("cell spec: %w", err)
	}
	if err := checkExactKeys(data); err != nil {
		return CellSpec{}, fmt.Errorf("cell spec: %w", err)
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var s CellSpec
	if err := dec.Decode(&s); err != nil {
		return CellSpec{}, fmt.Errorf("decode cell spec: %w", err)
	}
	// Strict EOF: any non-whitespace byte after the first value — including a
	// stray delimiter like `]` or `}` — is rejected. dec.More() is an
	// iteration helper, not an EOF check, so decode a second value and
	// require io.EOF.
	var extra json.RawMessage
	if err := dec.Decode(&extra); err != io.EOF {
		return CellSpec{}, fmt.Errorf("cell spec: trailing data after JSON object")
	}

	if s.Version != SchemaVersion {
		return CellSpec{}, fmt.Errorf("cell spec: version must be %q, got %q", SchemaVersion, s.Version)
	}
	if s.Contract.ID == "" {
		return CellSpec{}, fmt.Errorf("cell spec: contract.id is required")
	}
	if s.Contract.Goal == "" { // parity with #CellSpec (Pi round-5 D2)
		return CellSpec{}, fmt.Errorf("cell spec: contract.goal is required")
	}
	if s.ProtocolID == "" { // opaque provenance (Pi PR-#718-fido β D6)
		return CellSpec{}, fmt.Errorf("cell spec: protocol_id is required")
	}
	for side, decl := range map[string]json.RawMessage{"alpha": s.Alpha, "beta": s.Beta} {
		if len(decl) == 0 {
			return CellSpec{}, fmt.Errorf("cell spec: %s is required", side)
		}
		if _, err := cellfill.FillID(decl); err != nil {
			return CellSpec{}, fmt.Errorf("cell spec: %s: %w", side, err)
		}
	}
	if err := validateEvidence(s.Contract.RequiredEvidence); err != nil {
		return CellSpec{}, fmt.Errorf("cell spec: %w", err)
	}
	for name, p := range s.Params {
		if p.Kind != "skill" && p.Kind != "value" {
			return CellSpec{}, fmt.Errorf("cell spec: parameter %q has unsupported kind %q (want \"skill\" or \"value\")", name, p.Kind)
		}
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
// declarations are complete tagged objects ready for their fills.
type Resolved struct {
	Spec  CellSpec
	Alpha json.RawMessage
	Beta  json.RawMessage
}

// Resolve fills `$param` holes from `given` (with defaults and closed
// domains) IN PLACE inside both seat declarations: any string value equal to
// `$name` is replaced by the parameter's value, wherever it sits in the seat
// tree. Unresolved authored JSON carries holes in the same positions the
// resolved tree fills.
func (s CellSpec) Resolve(given map[string]string) (Resolved, error) {
	for name := range given {
		if _, ok := s.Params[name]; !ok {
			return Resolved{}, fmt.Errorf("unknown parameter %q (not declared by cell %q)", name, s.Contract.ID)
		}
	}
	vals := make(map[string]string)
	var missing []string
	for name, p := range s.Params {
		switch {
		case given[name] != "":
			vals[name] = given[name]
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

	alpha, err := splice(s.Alpha, s.Params, vals)
	if err != nil {
		return Resolved{}, fmt.Errorf("alpha: %w", err)
	}
	beta, err := splice(s.Beta, s.Params, vals)
	if err != nil {
		return Resolved{}, fmt.Errorf("beta: %w", err)
	}
	return Resolved{Spec: s, Alpha: alpha, Beta: beta}, nil
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

// Build dispatches both seat declarations through the fill registry and binds
// the result to a runnable kernel Spec plus the RunMeta the record carries:
// the complete canonical resolved declarations and the truthful combined
// mode. The registry arrives from the assembly point (the CLI domain) — this
// package never names a fill.
func (r Resolved) Build(reg cellfill.Registry) (cellkernel.Spec, cellkernel.RunMeta, error) {
	alpha, err := reg.ConstructAlpha(r.Alpha)
	if err != nil {
		return cellkernel.Spec{}, cellkernel.RunMeta{}, fmt.Errorf("alpha: %w", err)
	}
	beta, err := reg.ConstructBeta(r.Beta)
	if err != nil {
		return cellkernel.Spec{}, cellkernel.RunMeta{}, fmt.Errorf("beta: %w", err)
	}

	req := make([]cellkernel.RequiredRef, 0, len(r.Spec.Contract.RequiredEvidence))
	for _, e := range r.Spec.Contract.RequiredEvidence {
		req = append(req, cellkernel.RequiredRef{ID: e.ID, Kind: e.Kind, Producer: cellkernel.Role(e.Producer)})
	}
	contract := cellkernel.Contract{
		ID:               r.Spec.Contract.ID,
		Goal:             r.Spec.Contract.Goal,
		RequiredEvidence: req,
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
	return cellkernel.Spec{Contract: contract, Alpha: alpha.Seat, Beta: beta.Seat}, meta, nil
}

// checkNoDuplicateKeysOrNull rejects duplicate object keys anywhere in the
// JSON, which encoding/json otherwise silently accepts (last-wins), and JSON
// null anywhere — the schema admits none, and a null collection would
// silently decode to a Go nil (Pi round-6 D2).
func checkNoDuplicateKeysOrNull(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	return walkNoDup(dec)
}

func walkNoDup(dec *json.Decoder) error {
	t, err := dec.Token()
	if err != nil {
		return err
	}
	if t == nil {
		return fmt.Errorf("null is not allowed (the cell-spec schema admits no null)")
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
	if err := keysIn("spec", root, "version", "contract", "protocol_id", "params", "alpha", "beta"); err != nil {
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
		if err := objectKeys(p, fmt.Sprintf("parameter %q", name), "kind", "required", "default", "domain"); err != nil {
			return err
		}
	}
	return nil
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
