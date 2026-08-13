package cellcog

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// The one process seam every provider adapter shares. Adapters contribute a
// pure argv recipe; this file owns the bounded, stateless execution around it,
// so neither provider owns mechanics the other also needs.

// CodingToolSurface is the built-in tool set a PRODUCING seat is offered. It
// matches the live `cnos-cds-dispatch` workflow's allow-list exactly, because
// a cell is meant to mechanize what the operator already does by hand — not to
// be a weaker Claude.
//
// Bash is here deliberately, and its earlier absence was a mistake worth
// naming: a software-development seat that cannot run `go test`, `cue vet` or
// `gofmt` cannot verify its own work, so it produces plausible code it has no
// way to check. The only real α episode before this change wrote a Markdown
// file, because prose was all an unverifying seat could honestly finish.
//
// This list is a CAPABILITY DECLARATION, not a containment mechanism, and the
// distinction was previously blurred. Withholding Bash never provided
// confinement — this package claims none — it only removed the seat's ability
// to check itself. What actually bounds an episode is unchanged: a disposable
// worktree, a runtime-measured diff, credentials that never enter cell JSON,
// and `--safe-mode`. A seat with Bash has the same reach as the operator
// running Claude Code themselves; real containment belongs to the execution
// substrate, which is why CI and local use the identical surface.
const CodingToolSurface = "Read,Write,Edit,MultiEdit,Glob,Grep,Bash"

// runCLI runs one provider invocation: prompt on stdin, output bounded as it
// streams, timeout, and WaitDelay — killing the child does not by itself
// unblock Wait, because anything it spawned inherits the output pipe and holds
// it open; WaitDelay bounds that second wait.
//
// It returns captured stdout and whether the bound cut it short. Truncation
// is reported rather than decided: for a producing seat stdout is PROGRESS
// (the product is the worktree diff, so a clipped stream costs nothing),
// while for a seat whose stdout IS the product a clipped stream means an
// answer that cannot be trusted. Only the caller knows which. This increment
// has one caller and it is the producing one; the split is kept as the donor
// proved it rather than collapsed into a policy runCLI is not entitled to
// make.
//
// Nothing here is an OS sandbox. The honest authority is the offered tool
// surface plus the runtime-measured worktree: whatever a seat touches
// elsewhere simply never becomes evidence.
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
	stdout := &boundedBuffer{max: maxOutputBytes}
	stderr := &boundedBuffer{max: maxStderrBytes}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	cmd.WaitDelay = waitDelay

	err := cmd.Run()
	if ctxErr := ctx.Err(); ctxErr != nil {
		// A hang is the one failure that leaves no other trace: the child is
		// killed, no diff is measured, and no answer comes back. What it
		// managed to emit before stalling is the only evidence of where it
		// stalled, and this path previously discarded both streams.
		//
		// It goes into the ERROR, not a side file. A diagnostic written
		// somewhere durable becomes a second, unreceipted account of the
		// episode; an error is already the explicit outcome channel and
		// cannot be mistaken for evidence.
		return "", stdout.truncated, fmt.Errorf("%s did not finish within %s: %w (captured %d stdout bytes before the stall; stderr tail: %q)",
			bin, timeout, ctxErr, len(stdout.String()), tail(stderr.String(), diagnosticTailBytes))
	}
	if err != nil {
		return "", stdout.truncated, fmt.Errorf("%s failed: %w (stderr: %s)", bin, err, strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), stdout.truncated, nil
}

// tail returns at most n trailing bytes of s, kept valid UTF-8 so a truncated
// diagnostic cannot corrupt the error it is embedded in. Trailing, not
// leading: when a provider stalls, its last output is what says where.
func tail(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return "…" + strings.ToValidUTF8(s[len(s)-n:], "")
}
