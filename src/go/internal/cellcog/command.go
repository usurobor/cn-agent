package cellcog

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/usurobor/cnos/src/go/internal/cellbound"
)

// The one process seam every provider adapter shares: adapters contribute a
// pure argv recipe, this file the bounded execution, so neither owns mechanics
// the other also needs.

// CodingToolSurface is the built-in tool set a PRODUCING seat is offered. It
// matches the live `cnos-cds-dispatch` workflow's allow-list exactly, because
// a cell mechanizes what the operator already does by hand rather than being a
// weaker Claude. Bash is in it because a seat that cannot run `go test`, `cue
// vet` or `gofmt` cannot verify its own work; withholding it bought no
// confinement, only a seat that cannot check itself, so CI and local use the
// identical surface. What bounds an episode is the disposable worktree, the
// measured diff, credentials that never enter cell JSON, and `--safe-mode`.
const CodingToolSurface = "Read,Write,Edit,MultiEdit,Glob,Grep,Bash"

// NoTools is the surface an ANSWERING seat is offered, and unlike the coding
// surface this one IS load-bearing: a reviewer's canonical input is
// `(contract, matter)`, so any file tool would cost the independence — not the
// containment — that the seat exists to provide.
const NoTools = ""

// runCLI runs one bounded invocation, prompt on stdin. WaitDelay is set
// because killing the child does not by itself unblock Wait: anything it
// spawned inherits the output pipe and holds it open. Truncation is reported
// rather than decided — stdout is PROGRESS for a producing seat but IS the
// product for an answering one, so only the caller can judge a clipped stream.
func runCLI(ctx context.Context, bin, dir, prompt string, args []string, timeout time.Duration) (string, bool, error) {
	if _, err := exec.LookPath(bin); err != nil {
		return "", false, fmt.Errorf("%q not found in PATH: %w", bin, err)
	}
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(prompt)
	// Head-keeping: stdout is consumed from the start — Answer parses it and
	// refuses it when clipped — so the first bytes are the ones worth keeping.
	stdout := cellbound.New(cellbound.KeepHead, maxOutputBytes)
	stderr := cellbound.New(cellbound.KeepHead, maxStderrBytes)
	cmd.Stdout, cmd.Stderr = stdout, stderr
	cmd.WaitDelay = waitDelay

	err := cmd.Run()
	if ctxErr := ctx.Err(); ctxErr != nil {
		// A hang leaves no other trace — child killed, no diff, no answer — so
		// what it emitted is the only evidence of where it stalled. That goes in
		// the ERROR, not a side file: a durable diagnostic would be a second,
		// unreceipted account of the episode.
		return "", stdout.Truncated(), fmt.Errorf("%s did not finish within %s: %w (captured %d stdout bytes before the stall; stderr tail: %q)",
			bin, timeout, ctxErr, len(stdout.String()), cellbound.Tail(stderr.String(), diagnosticTailBytes))
	}
	if err != nil {
		return "", stdout.Truncated(), fmt.Errorf("%s failed: %w (stderr: %s)", bin, err, strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), stdout.Truncated(), nil
}
