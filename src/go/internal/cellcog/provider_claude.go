package cellcog

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ToolSurface is the built-in tool set a producing seat is offered. File tools
// only: a seat needs to change files, not command the host.
const ToolSurface = "Read,Write,Edit,Glob,Grep"

// ClaudeCLI rents workspace cognition from the Claude Code CLI. Stateless: a
// Work call is one fresh bounded subprocess; no session is started or kept.
type ClaudeCLI struct {
	Model   string        // exact model id; required
	Bin     string        // default "claude"
	Timeout time.Duration // default 10m
}

func (ClaudeCLI) Name() string { return "claude-cli" }

// ClaudeArgv is the exact invocation, built purely so it can be asserted
// without spawning anything.
//
//   - `--tools` RESTRICTS the available built-in set. `--allowedTools` only
//     pre-approves tools that remain available, so using it would have left
//     Bash reachable while claiming otherwise; it must never appear here.
//   - `--setting-sources ""` loads no user/project/local settings, and
//     `--strict-mcp-config` admits no ambient MCP servers, so local
//     customization cannot become a second, unreceipted component definition
//     beside the fill's digested skills.
//   - `--no-session-persistence` keeps the adapter stateless.
//
// What this does NOT provide is OS confinement. The honest authority is the
// offered tool surface plus the runtime-measured worktree: whatever a seat
// touches elsewhere simply never becomes evidence.
func ClaudeArgv(model string) []string {
	return []string{
		"-p",
		"--model", model,
		"--no-session-persistence",
		"--setting-sources", "",
		"--strict-mcp-config",
		"--tools", ToolSurface,
		"--output-format", "text",
	}
}

func (c ClaudeCLI) Work(ctx context.Context, dir, prompt string) error {
	if dir == "" {
		return fmt.Errorf("claude-cli: Work needs a directory")
	}
	bin := c.Bin
	if bin == "" {
		bin = "claude"
	}
	return runCLI(ctx, bin, dir, prompt, ClaudeArgv(c.Model), c.Timeout)
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
