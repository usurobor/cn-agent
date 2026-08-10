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

// ToolSurface is the built-in tool set a producing seat is offered. File tools
// only: a seat needs to change files, not command the host.
const ToolSurface = "Read,Write,Edit,Glob,Grep"

// runCLI runs one provider invocation: prompt on stdin, output bounded as it
// streams, timeout, and WaitDelay — killing the child does not by itself
// unblock Wait, because anything it spawned inherits the output pipe and holds
// it open; WaitDelay bounds that second wait.
//
// Nothing here is an OS sandbox. The honest authority is the offered tool
// surface plus the runtime-measured worktree: whatever a seat touches
// elsewhere simply never becomes evidence.
func runCLI(ctx context.Context, bin, dir, prompt string, args []string, timeout time.Duration) error {
	if _, err := exec.LookPath(bin); err != nil {
		return fmt.Errorf("%q not found in PATH: %w", bin, err)
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
		return fmt.Errorf("%s did not finish within %s: %w", bin, timeout, ctxErr)
	}
	if err != nil {
		return fmt.Errorf("%s failed: %w (stderr: %s)", bin, err, strings.TrimSpace(stderr.String()))
	}
	if stdout.truncated {
		return fmt.Errorf("%s produced more than %d bytes of output", bin, maxOutputBytes)
	}
	return nil
}
