// Package cellcog constructs bounded, stateless provider adapters that a
// fill's constructor composes into a seat: Pi's construction boundary
// (msg-cn-pi-cnos-cds-fill-construction-51) assigns it model selection, typed
// safe argv, timeout/output bounds and capability truth, and no fill semantics
// — what a prompt means, what a worktree is for and what counts as evidence
// stay with the fill. Ports are narrow, each added only for a consumer that
// already needs it (Pi #59 B1).
//
// A declared tool surface is NOT confinement, said once here rather than
// beside every flag: a seat offered Bash reaches as far as the operator
// running Claude Code does, and only the worktree is measured, so nothing it
// touches elsewhere becomes evidence. Containment is the execution
// substrate's, and its managed policy stays above this baseline undetected.
//
// Admitted: claude-cli and fake. codex-cli is HELD (Pi #55 D1) — its
// suppression flags leave global and project AGENTS.md and discovered skills
// loading, which would make ambient instruction a second, unreceipted
// component definition. Return conditions, preserved argv research and this
// package's measurement history: docs/architecture/CDS-CELL-MIGRATION.md.
package cellcog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Coder is a workspace-EDIT port and only that: a directory goes in, nothing
// comes back but whether the process ran. Successful stdout is discarded
// because what the seat did is measured from the worktree, not taken from its
// word; a seat whose product is an ANSWER uses Answerer (Pi #56 B1).
type Coder interface {
	Name() string
	Work(ctx context.Context, dir, prompt string) error
}

// Answerer produces a VALUE rather than an edit, because a REVIEWING seat's
// whole product is a judgement: nothing is left in a worktree to measure, so
// `Coder` cannot serve it. The caller supplies the JSON Schema; what the shape
// MEANS stays with the fill. Deliberately not a widening of `Coder` — a
// producing seat must not gain a channel for reporting on itself.
type Answerer interface {
	Name() string
	Answer(ctx context.Context, prompt string, schema json.RawMessage) (json.RawMessage, error)
}

// ErrNoProvider is returned when a seat is constructed without cognition.
var ErrNoProvider = errors.New("cellcog: no provider")

// Config is what a fill passes through from the seat: provider and REQUESTED
// MODEL SELECTOR, nothing else — no cell can smuggle argv into an adapter.
type Config struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
}

// Mode is what a constructed coder TRUTHFULLY runs under: renting real
// cognition is `cognitive`, the deterministic fake is `mechanical`.
type Mode string

const (
	ModeCognitive  Mode = "cognitive"
	ModeMechanical Mode = "mechanical"
)

// ErrNoDeterministicAnswer is what NewAnswerer returns for `fake`: a REFUSAL,
// not a gap to fill later. A fake answerer was tried and rejected —
// `{"pass":false,"notes":"…"}` puts one reviewing fill's verdict vocabulary
// inside a package that owns no fill semantics, and an honest judgement
// depends on what was asked, which is the fill's knowledge. The mode still
// comes back: nothing was rented, and that much IS this package's truth.
var ErrNoDeterministicAnswer = errors.New("cellcog: the answering port has no deterministic provider (a judgement cannot be fabricated by a provider-neutral adapter)")

// NewAnswerer is New for the answering port: same closed provider set, same
// model rule, and `fake` differs in what it can honestly supply.
func NewAnswerer(cfg Config) (Answerer, Mode, error) {
	switch cfg.Provider {
	case "claude-cli":
		if cfg.Model == "" {
			return nil, "", fmt.Errorf("cellcog: provider %q requires a model selector", cfg.Provider)
		}
		return ClaudeCLI{Model: cfg.Model}, ModeCognitive, nil
	case "fake":
		// A model id the fake would ignore must not be receipted as though it
		// selected something (identical rule in the CUE overlay), checked BEFORE
		// the refusal so a malformed declaration reads as malformed.
		if cfg.Model != "" {
			return nil, "", fmt.Errorf("cellcog: provider %q takes no model, got %q", cfg.Provider, cfg.Model)
		}
		return nil, ModeMechanical, ErrNoDeterministicAnswer
	default:
		return nil, "", fmt.Errorf("cellcog: unknown provider %q (want claude-cli or fake)", cfg.Provider)
	}
}

// New constructs the adapter for a cognition declaration. The provider set is
// closed — a typo fails construction, before any invocation.
func New(cfg Config) (Coder, Mode, error) {
	switch cfg.Provider {
	case "claude-cli":
		if cfg.Model == "" {
			return nil, "", fmt.Errorf("cellcog: provider %q requires a model selector", cfg.Provider)
		}
		return ClaudeCLI{Model: cfg.Model}, ModeCognitive, nil
	case "fake":
		// A model id the fake would ignore must not be receipted as though it
		// selected something: the rule is identical in the CUE overlay.
		if cfg.Model != "" {
			return nil, "", fmt.Errorf("cellcog: provider %q takes no model, got %q", cfg.Provider, cfg.Model)
		}
		return FakeCoder{}, ModeMechanical, nil
	default:
		return nil, "", fmt.Errorf("cellcog: unknown provider %q (want claude-cli or fake)", cfg.Provider)
	}
}

// Adapter bounds. A rented seat is the one part of an episode the runtime
// cannot predict: an unbounded provider burns the budget and dies mid-cell.
const (
	// defaultTimeout was 10 minutes until measurement moved it: rented episodes
	// cluster at 460-500s, inside that 600s bound, so it was cutting off live
	// episodes rather than catching runaways. 30 minutes is ~3.6x the observed
	// peak. Compiled-in and never cell-supplied: a timeout a cell could set is
	// a cell that can decline to be bounded.
	defaultTimeout      = 30 * time.Minute
	maxOutputBytes      = 4 << 20 // matches the kernel's aggregate artifact bound
	maxStderrBytes      = 8 << 10 // diagnostics only
	waitDelay           = 2 * time.Second
	diagnosticTailBytes = 400 // bounded stderr tail carried in a timeout error
)
