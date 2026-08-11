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

// ToolSurface is the built-in tool set a PRODUCING seat is offered. File tools
// only: a seat needs to change files, not command the host.
const ToolSurface = "Read,Write,Edit,Glob,Grep"

// NoTools is the surface an ANSWERING seat is offered. A reviewer's canonical
// input is `(contract, matter)` and nothing else, so giving it file tools
// would let it reach outside that input and read the very workspace it is
// meant to judge from the outside. Its authority is therefore strictly less
// than a producer's: no tools, and consequently no permission mode to declare.
const NoTools = ""

// runCLI runs one provider invocation: prompt on stdin, output bounded as it
// streams, timeout, and WaitDelay — killing the child does not by itself
// unblock Wait, because anything it spawned inherits the output pipe and holds
// it open; WaitDelay bounds that second wait.
//
// It returns captured stdout. A producing seat discards it — what it did is
// measured from the worktree — while an answering seat's whole product IS that
// output.
//
// Nothing here is an OS sandbox. The honest authority is the offered tool
// surface plus the runtime-measured worktree: whatever a seat touches
// elsewhere simply never becomes evidence.
func runCLI(ctx context.Context, bin, dir, prompt string, args []string, timeout time.Duration) (string, error) {
	if _, err := exec.LookPath(bin); err != nil {
		return "", fmt.Errorf("%q not found in PATH: %w", bin, err)
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
		return "", fmt.Errorf("%s did not finish within %s: %w (captured %d stdout bytes before the stall; stderr tail: %q)",
			bin, timeout, ctxErr, len(stdout.String()), tail(stderr.String(), diagnosticTailBytes))
	}
	if err != nil {
		return "", fmt.Errorf("%s failed: %w (stderr: %s)", bin, err, strings.TrimSpace(stderr.String()))
	}
	if stdout.truncated {
		return "", fmt.Errorf("%s produced more than %d bytes of output", bin, maxOutputBytes)
	}
	return stdout.String(), nil
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
