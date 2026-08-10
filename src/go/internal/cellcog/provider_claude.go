package cellcog

import (
	"context"
	"fmt"
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

// ClaudeArgv is the exact recipe, built purely so it can be asserted without
// spawning anything.
//
//   - `--safe-mode` disables CLAUDE.md, skills, plugins, hooks, MCP servers,
//     custom commands/agents and auto-memory while keeping authentication,
//     model selection and the built-in tools. Repository-local guidance is not
//     legitimate implicit context here: it would be a second, unreceipted
//     component definition beside the fill's ordered, digested skills. Project
//     guidance that a seat should have must be declared AS a skill.
//   - `--tools` RESTRICTS the available built-in set. `--allowedTools` merely
//     pre-approves tools that remain available, so using it would have left
//     Bash reachable while claiming otherwise; it must never appear here.
//   - `--no-session-persistence` keeps the adapter stateless.
func ClaudeArgv(model string) []string {
	return []string{
		"-p",
		"--model", model,
		"--safe-mode",
		"--no-session-persistence",
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
