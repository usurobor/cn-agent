// Package model holds the domain objects the runner emits: the contract, the
// coherence-measurement object (CM), the wave, the judgment, and the receipt.
// Every struct's JSON field names match the CUE definitions in
// ../schema/*.cue, so each emitted object round-trips through `cue vet`.
//
// This is spike code (see ../DESIGN.md) — a walking skeleton for the CC→PC→WC
// loop, not canonical R12 runtime.
package model

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Contract is the pinned unit of work handed to one kernel run (#Contract).
type Contract struct {
	ID          string   `json:"id"`
	Telos       string   `json:"telos"`
	Target      string   `json:"target"`
	Instruction string   `json:"instruction"`
	Acceptance  []string `json:"acceptance"`
}

// Defect is one measured incoherence (#Defect).
type Defect struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Detail string `json:"detail"`
}

// Provenance records who/how produced a CM.
type Provenance struct {
	MeasuredBy string `json:"measured_by"`
	Method     string `json:"method"`
	CMVersion  string `json:"cm_version"`
}

// CM is the coherence-measurement object (#CM) — the cnos↔tsc seam.
type CM struct {
	Target        string     `json:"target"`
	AlphaScore    float64    `json:"alpha_score"`
	Defects       []Defect   `json:"defects"`
	Provenance    Provenance `json:"provenance"`
	MeasuredAtRef string     `json:"measured_at_ref"`
}

// Incoherence is the defect count — the loop's decreasing quantity.
func (c CM) Incoherence() int { return len(c.Defects) }

// Wave is the PC's relation-graph output (#Wave).
type Wave struct {
	ID        string     `json:"id"`
	FromCMRef string     `json:"from_cm_ref"`
	Contracts []Contract `json:"contracts"`
}

// Judgment is the CC's δ-output over a CM (#Judgment).
type Judgment struct {
	Verdict      string `json:"verdict"`
	CMRef        string `json:"cm_ref"`
	Target       string `json:"target"`
	DecidedAtRef string `json:"decided_at_ref"`
}

// TransitionRecord is one kernel edge actually traversed (#TransitionRecord).
type TransitionRecord struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Step  string `json:"step"`
	Guard string `json:"guard"`
}

// ValidationVerdict is V's mechanical output (#ValidationVerdict).
type ValidationVerdict struct {
	Verdict          string   `json:"verdict"`
	FailedPredicates []string `json:"failed_predicates"`
	Warnings         []string `json:"warnings"`
}

// BoundaryDecision is δ's action at the cell boundary (#BoundaryDecision).
type BoundaryDecision struct {
	Actor  string `json:"actor"`
	Action string `json:"action"`
}

// Receipt is the parent-facing trust surface of a closed/held cell (#Receipt).
type Receipt struct {
	ID               string             `json:"id"`
	Cell             string             `json:"cell"`
	Telos            string             `json:"telos"`
	Transitions      []TransitionRecord `json:"transitions"`
	EvidenceRefs     map[string]string  `json:"evidence_refs"`
	Validation       ValidationVerdict  `json:"validation"`
	BoundaryDecision BoundaryDecision   `json:"boundary_decision"`
	Transmissibility string             `json:"transmissibility"`
	Verdict          string             `json:"verdict"`
	ProducedAtRef    string             `json:"produced_at_ref"`
}

// DeriveTransmissibility computes the trust property from the
// (validation verdict × boundary action) pair, mirroring the structural table
// in ../schema/receipt.cue. Returns "" for an inconsistent pair so the caller
// fails closed (and `cue vet` would reject the emitted receipt anyway).
func DeriveTransmissibility(verdict, action string) string {
	switch verdict {
	case "PASS":
		switch action {
		case "accept", "release":
			return "accepted"
		case "reject", "repair_dispatch":
			return "not_transmissible"
		}
	case "FAIL":
		switch action {
		case "reject", "repair_dispatch":
			return "not_transmissible"
		case "override":
			return "degraded"
		}
	}
	return ""
}

// WriteJSON marshals v as indented JSON (a valid CUE/YAML input) to path,
// creating parent directories as needed. JSON is used rather than YAML to stay
// stdlib-only (no third-party YAML dependency); JSON is a strict subset of
// YAML and `cue vet` reads it directly. See ../DESIGN.md §Build notes.
func WriteJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}
