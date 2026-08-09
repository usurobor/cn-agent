// Package cellspec loads a serialized cell spec (the JSON IR a compiled
// main.cell — or a hand-authored contract file — produces) and binds it to a
// runnable cellkernel.Spec.
//
// Parsing is strict (Pi β msg-cn-pi-cnos-cell-prototype-beta-32, D5): a pinned
// schema version, no unknown/trailing/duplicate JSON, required seats present,
// required-evidence uniqueness + producer authority, and a known builtin
// profile. The spec is data; γ/V/δ are kernel-owned and never appear here.
//
// v0 has no cognition: a `profile` selects a builtin mechanical seat pair
// ("stub" for smoke, "bool" for a real independently-checked episode). Rented
// cognition (Phase 3) swaps the profile for a provider-backed seat and leaves
// this contract unchanged.
//
// See docs/architecture/CDS-CELL-MIGRATION.md and schemas/cdd/spec.cue.
package cellspec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cellkernel"
)

// SchemaVersion is the pinned cell-spec version; a spec must declare it exactly.
const SchemaVersion = "cnos.cellspec.v0"

// EpisodeReceiptSchema is the generic receipt schema the v0 runner actually
// emits. A spec's protocol_id is carried as *declared* provenance; the runner
// never claims to have validated a protocol it did not execute (Pi D1).
const EpisodeReceiptSchema = "cnos.cellkernel.episode-receipt.v0"

// knownProtocols are the protocol identifiers the runner recognizes. An unknown
// (typo'd) declared protocol fails closed; a known-but-unexecuted protocol
// (e.g. CDS in v0) is carried as provenance with execution_mode set honestly.
var knownProtocols = map[string]bool{
	EpisodeReceiptSchema:      true,
	"cnos.cdd.receipt.v1":     true,
	"cnos.cdd.cds.receipt.v1": true,
	"cnos.cdd.cdr.receipt.v1": true,
	"cnos.cdd.cdw.receipt.v1": true,
}

// IsKnownProtocol reports whether id is a recognized protocol identifier.
func IsKnownProtocol(id string) bool { return knownProtocols[id] }

// CellSpec is the serialized, strictly-validated cell.
type CellSpec struct {
	Version    string               `json:"version"`
	Contract   ContractSpec         `json:"contract"`
	ProtocolID string               `json:"protocol_id"`
	Profile    string               `json:"profile,omitempty"` // builtin seat profile; default "stub"
	Params     map[string]ParamSpec `json:"params,omitempty"`
	Alpha      *SeatSpec            `json:"alpha"`
	Beta       *SeatSpec            `json:"beta"`
}

type ContractSpec struct {
	ID               string        `json:"id"`
	Goal             string        `json:"goal"`
	RequiredEvidence []RequiredRef `json:"required_evidence,omitempty"`
}

// RequiredRef names a required evidence ref and the producer role authorized to
// mint it (Pi D2).
type RequiredRef struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Producer string `json:"producer"`
}

type ParamSpec struct {
	Kind     string   `json:"kind"`
	Required bool     `json:"required"`
	Default  string   `json:"default,omitempty"`
	Domain   []string `json:"domain,omitempty"`
}

type SeatSpec struct {
	Skills []string `json:"skills"`
}

// Parse strictly decodes and validates a serialized cell spec.
func Parse(data []byte) (CellSpec, error) {
	if err := checkNoDuplicateKeys(data); err != nil {
		return CellSpec{}, fmt.Errorf("cell spec: %w", err)
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var s CellSpec
	if err := dec.Decode(&s); err != nil {
		return CellSpec{}, fmt.Errorf("decode cell spec: %w", err)
	}
	if dec.More() { // trailing data after the first JSON value.
		return CellSpec{}, fmt.Errorf("cell spec: trailing data after JSON object")
	}

	if s.Version != SchemaVersion {
		return CellSpec{}, fmt.Errorf("cell spec: version must be %q, got %q", SchemaVersion, s.Version)
	}
	if s.Contract.ID == "" {
		return CellSpec{}, fmt.Errorf("cell spec: contract.id is required")
	}
	if s.ProtocolID == "" {
		return CellSpec{}, fmt.Errorf("cell spec: protocol_id is required")
	}
	if !IsKnownProtocol(s.ProtocolID) {
		return CellSpec{}, fmt.Errorf("cell spec: unknown protocol_id %q", s.ProtocolID)
	}
	if s.Alpha == nil || s.Beta == nil {
		return CellSpec{}, fmt.Errorf("cell spec: alpha and beta must both be present")
	}
	if s.Profile == "" {
		s.Profile = ProfileStub
	}
	if !isKnownProfile(s.Profile) {
		return CellSpec{}, fmt.Errorf("cell spec: unknown profile %q", s.Profile)
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

// Resolved is a cell spec with its parameter holes filled and seat skills
// spliced — ready to bind to a kernel Spec.
type Resolved struct {
	Spec        CellSpec
	Params      map[string]string
	AlphaSkills []string
	BetaSkills  []string
}

// Resolve fills parameter holes from `given`, applies defaults, checks domains,
// and splices `$param` refs into the seat skills.
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
	for name, p := range s.Params {
		if v, ok := vals[name]; ok && len(p.Domain) > 0 && !contains(p.Domain, v) {
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
	}
	return out, nil
}

// KernelSpec binds the resolved cell spec to a runnable kernel Spec by
// constructing the builtin profile's seat pair. ExecutionMode reports whether
// the seats are a stub (smoke) or a real mechanical profile.
func (r Resolved) KernelSpec() (cellkernel.Spec, string, error) {
	req := make([]cellkernel.RequiredRef, 0, len(r.Spec.Contract.RequiredEvidence))
	for _, e := range r.Spec.Contract.RequiredEvidence {
		req = append(req, cellkernel.RequiredRef{ID: e.ID, Kind: e.Kind, Producer: cellkernel.Role(e.Producer)})
	}
	contract := cellkernel.Contract{
		ID:               r.Spec.Contract.ID,
		Goal:             r.Spec.Contract.Goal,
		RequiredEvidence: req,
	}
	alpha, beta, mode, err := buildProfile(r)
	if err != nil {
		return cellkernel.Spec{}, "", err
	}
	return cellkernel.Spec{Contract: contract, Alpha: alpha, Beta: beta}, mode, nil
}

func contains(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

// checkNoDuplicateKeys rejects duplicate object keys anywhere in the JSON, which
// encoding/json otherwise silently accepts (last-wins).
func checkNoDuplicateKeys(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	return walkNoDup(dec)
}

func walkNoDup(dec *json.Decoder) error {
	t, err := dec.Token()
	if err != nil {
		return err
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
