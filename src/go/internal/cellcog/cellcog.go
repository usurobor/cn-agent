// Package cellcog constructs bounded, stateless provider adapters that a
// fill's constructor composes into a seat.
//
// Scope of the word "reusable" (Pi #59 B1): the PROCESS AND PROVIDER SEAM
// here is reusable — argv recipes, bounded execution, timeouts, output
// limits. The PORT is not general cognition. `Coder` supports workspace
// edits only, so a planning or research fill cannot rent it as it stands; it
// would need a returned-value port that does not exist yet. Do not describe
// this package as a general cognition subsystem.
//
// The package owns exactly what Pi's construction boundary assigns it
// (msg-cn-pi-cnos-cds-fill-construction-51): explicit model selection,
// executable invocation with typed safe arguments, timeout/cancellation/
// output bounds, stateless operation, and provider capability truth. It owns
// no fill semantics — what a prompt means, what a worktree is for, and what
// counts as evidence belong to the fill that rents the cognition.
//
// Scope, stated honestly: an adapter points a provider at one directory with
// file tools only. That is a working-directory boundary plus the provider's
// own workspace rules, NOT an OS sandbox; a provider writing an absolute path
// elsewhere is not stopped here — it merely gains nothing, because only the
// worktree is ever measured. Real containment belongs under the adapter and
// is not implemented.
//
// The admitted provider set is claude-cli and fake. codex-cli is HELD, not
// forgotten (Pi #55 D1): a seat may carry only the fill's ordered, digested
// skills, and Codex's available suppression flags do not reach far enough —
// `--ignore-user-config` suppresses only $CODEX_HOME/config.toml and
// `--ignore-rules` only execpolicy .rules, while global and project AGENTS.md
// and discovered skills still load. Admitting the provider on those flags
// would let ambient instructions become a second, unreceipted component
// definition. It returns when the execution substrate can supply a clean
// CODEX_HOME and a real run proves poisoned AGENTS.md and ambient skills do
// not reach the invocation; the argv research is preserved in
// docs/architecture/CDS-CELL-MIGRATION.md rather than in unreachable code.
package cellcog

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"
)

// Coder is a workspace-EDIT cognition port, and only that: a directory goes
// in, nothing comes back but whether the process ran. Successful stdout is
// deliberately discarded, because what the seat actually did is measured from
// the worktree rather than taken from its word.
//
// That shape is honest about its limits (Pi #56 B1): a seat whose product is
// an ANSWER rather than an edit — planning, research, review — cannot use
// this port, because there is no value to return. Case 3 may add the smallest
// returned-value port when beta actually needs one. Nothing here should be
// widened in advance to look like a general cognition framework.
type Coder interface {
	Name() string
	Work(ctx context.Context, dir, prompt string) error
}

// ErrNoProvider is returned when a seat is constructed without cognition.
var ErrNoProvider = errors.New("cellcog: no provider")

// Config is the inline cognition declaration a fill passes through from the
// seat: which provider and which REQUESTED MODEL SELECTOR. Nothing else — a
// cell cannot smuggle argv or flags into an adapter.
type Config struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
}

// Modes a constructed coder truthfully runs under: a provider that rents
// real cognition is `cognitive`; the deterministic fake is `mechanical`.
type Mode string

const (
	ModeCognitive  Mode = "cognitive"
	ModeMechanical Mode = "mechanical"
)

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

// Adapter bounds. They exist because a rented seat is the one part of an
// episode the runtime cannot predict: an unbounded provider burns the whole
// budget and dies mid-cell.
const (
	defaultTimeout = 10 * time.Minute
	maxOutputBytes = 4 << 20 // matches the kernel's aggregate artifact bound
	maxStderrBytes = 8 << 10 // diagnostics only
	waitDelay      = 2 * time.Second
)

// boundedBuffer captures at most max bytes and remembers that it had to stop.
// It never reports a short write, so a bound fails the run here rather than
// killing the child mid-stream with a broken pipe.
type boundedBuffer struct {
	max       int
	buf       bytes.Buffer
	truncated bool
}

func (b *boundedBuffer) Write(p []byte) (int, error) {
	switch room := b.max - b.buf.Len(); {
	case room >= len(p):
		b.buf.Write(p)
	case room > 0:
		b.buf.Write(p[:room])
		b.truncated = true
	default:
		b.truncated = true
	}
	return len(p), nil
}

func (b *boundedBuffer) String() string { return b.buf.String() }
