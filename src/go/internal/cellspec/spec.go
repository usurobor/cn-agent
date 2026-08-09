// Package cellspec loads a serialized cell spec (the JSON IR a compiled
// main.cell — or a hand-authored contract file — produces) and binds it to a
// runnable cellkernel.Spec.
//
// The serialized spec is data: contract, protocol_id, typed parameter holes,
// and the α/β skill lines. Parameters are Unix-shaped typed holes — required or
// optional-with-default — filled by the invoker (the CLI now, a parent cell
// later) and spliced into the seats via `$name` references. Resolving a value
// to a concrete skill and loading it is a later phase; here the holes are
// filled and spliced, and α/β are stubs, which is enough to run one real
// episode end-to-end through the corrected kernel.
//
// See docs/architecture/CDS-CELL-MIGRATION.md (Phases 0–1) and
// schemas/cdd/spec.cue (#CellSpec), the CUE contract this JSON vets against.
package cellspec

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cellkernel"
)

// CellSpec is the serialized, CUE-vettable cell (schemas/cdd/spec.cue
// #CellSpec). γ/V/δ are kernel-owned and never appear here.
type CellSpec struct {
	Contract   ContractSpec         `json:"contract"`
	ProtocolID string               `json:"protocol_id"`
	Params     map[string]ParamSpec `json:"params,omitempty"`
	Alpha      SeatSpec             `json:"alpha"`
	Beta       SeatSpec             `json:"beta"`
	Budget     *BudgetSpec          `json:"budget,omitempty"`
}

// ContractSpec is the cell's input contract.
type ContractSpec struct {
	ID               string        `json:"id"`
	Goal             string        `json:"goal"`
	RequiredEvidence []RequiredRef `json:"required_evidence,omitempty"`
}

// RequiredRef names an evidence ref the closed receipt must carry.
type RequiredRef struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

// ParamSpec is a typed hole (Unix arg/flag). Kind is the resolution family
// (e.g. "skill"); Required with no Default must be filled by the invoker.
type ParamSpec struct {
	Kind     string   `json:"kind"`
	Required bool     `json:"required"`
	Default  string   `json:"default,omitempty"`
	Domain   []string `json:"domain,omitempty"` // closed value set; empty = open
}

// SeatSpec is a seat's skill line; entries may be literal skill names or
// `$param` references spliced from resolved parameters.
type SeatSpec struct {
	Skills []string `json:"skills"`
}

// BudgetSpec bounds a cell's cognition (unused by the stub runner).
type BudgetSpec struct {
	Tokens int `json:"tokens,omitempty"`
}

// Parse decodes and structurally checks a serialized cell spec. It rejects
// unknown top-level keys so a typo fails fast (mirrors CUE closedness).
func Parse(data []byte) (CellSpec, error) {
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.DisallowUnknownFields()
	var s CellSpec
	if err := dec.Decode(&s); err != nil {
		return CellSpec{}, fmt.Errorf("decode cell spec: %w", err)
	}
	if s.Contract.ID == "" {
		return CellSpec{}, fmt.Errorf("cell spec: contract.id is required")
	}
	if s.ProtocolID == "" {
		return CellSpec{}, fmt.Errorf("cell spec: protocol_id is required")
	}
	return s, nil
}

// Resolved is a cell spec with its parameter holes filled and seat skills
// spliced — ready to bind to a kernel Spec.
type Resolved struct {
	Spec        CellSpec
	Params      map[string]string // param name → resolved value
	AlphaSkills []string          // spliced
	BetaSkills  []string
}

// Resolve fills the parameter holes from `given` (CLI `--param k=v`), applies
// defaults, and errors — with a Unix-style listing — on any unfilled required
// hole or unknown parameter. It then splices `$param` refs in the seat skills.
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
		case p.Default != "":
			vals[name] = p.Default
		case p.Required:
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return Resolved{}, fmt.Errorf("missing required parameter(s): %s (supply with --param <name>=<value>)", strings.Join(missing, ", "))
	}

	// Domain check: a filled value outside a param's closed domain fails now
	// (the "$PATH resolves, but a typo fails" property).
	for name, p := range s.Params {
		v, ok := vals[name]
		if !ok || len(p.Domain) == 0 {
			continue
		}
		if !contains(p.Domain, v) {
			return Resolved{}, fmt.Errorf("parameter %q value %q not in domain [%s]", name, v, strings.Join(p.Domain, ", "))
		}
	}

	alpha, err := splice(s.Alpha.Skills, s.Params, vals)
	if err != nil {
		return Resolved{}, fmt.Errorf("alpha skills: %w", err)
	}
	beta, err := splice(s.Beta.Skills, s.Params, vals)
	if err != nil {
		return Resolved{}, fmt.Errorf("beta skills: %w", err)
	}

	return Resolved{Spec: s, Params: vals, AlphaSkills: alpha, BetaSkills: beta}, nil
}

func contains(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

// splice replaces `$name` skill entries with their resolved parameter value. A
// `$name` referencing an undeclared parameter is a spec error; a declared but
// unfilled (optional, no default) parameter is dropped from the seat line.
func splice(skills []string, declared map[string]ParamSpec, vals map[string]string) ([]string, error) {
	out := make([]string, 0, len(skills))
	for _, sk := range skills {
		if !strings.HasPrefix(sk, "$") {
			out = append(out, sk)
			continue
		}
		name := strings.TrimPrefix(sk, "$")
		if _, ok := declared[name]; !ok {
			return nil, fmt.Errorf("skill %q references undeclared parameter %q", sk, name)
		}
		if v, ok := vals[name]; ok {
			out = append(out, v)
		}
		// declared-but-unfilled optional param: drop it from the seat line.
	}
	return out, nil
}

// KernelSpec binds the resolved cell spec to a runnable kernel Spec using stub
// α/β (Phase 1: no cognition, no real skill loading yet).
func (r Resolved) KernelSpec() cellkernel.Spec {
	req := make([]cellkernel.RequiredRef, 0, len(r.Spec.Contract.RequiredEvidence))
	for _, e := range r.Spec.Contract.RequiredEvidence {
		req = append(req, cellkernel.RequiredRef{ID: e.ID, Kind: e.Kind})
	}
	return cellkernel.Spec{
		Contract: cellkernel.Contract{
			ID:               r.Spec.Contract.ID,
			Goal:             r.Spec.Contract.Goal,
			RequiredEvidence: req,
		},
		Alpha: stubAlpha{skills: r.AlphaSkills},
		Beta:  stubBeta{skills: r.BetaSkills},
	}
}
