package cellcog

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ClaudeCLI rents workspace cognition from the Claude Code CLI. Stateless: a
// Work call is one fresh bounded subprocess; no session is started or kept.
type ClaudeCLI struct {
	Model   string        // exact model id; required
	Bin     string        // default "claude"
	Timeout time.Duration // default 10m
}

func (ClaudeCLI) Name() string { return "claude-cli" }

// Work runs the CLI confined to dir with FILE tools only — read, write, edit,
// search. No shell: a producing seat needs to change files, not command the
// host. Output is discarded up to a bound; the fill measures the worktree.
func (c ClaudeCLI) Work(ctx context.Context, dir, prompt string) error {
	if dir == "" {
		return fmt.Errorf("claude-cli: Work needs a directory")
	}
	bin := c.Bin
	if bin == "" {
		bin = "claude"
	}
	args := []string{
		"-p", "--output-format", "text",
		"--model", c.Model,
		"--permission-mode", "acceptEdits",
		"--allowedTools", "Read,Write,Edit,Glob,Grep",
	}
	return runCLI(ctx, bin, dir, prompt, args, c.Timeout)
}

// runCLI is the one subprocess shape every adapter shares: prompt on stdin,
// bounded output, timeout, and WaitDelay — killing the child does not by
// itself unblock Wait, because anything it spawned inherits the output pipe
// and holds it open; WaitDelay bounds that second wait.
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
