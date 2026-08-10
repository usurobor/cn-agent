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
//   - `--safe-mode` disables USER AND PROJECT customization — CLAUDE.md,
//     skills, plugins, hooks, MCP servers, custom commands/agents and
//     auto-memory — while keeping authentication, model selection and the
//     built-in tools. Repository-local guidance is not legitimate implicit
//     context here: it would be a second, unreceipted component definition
//     beside the fill's ordered, digested skills. Project guidance a seat
//     should have must be declared AS a skill.
//
//     Stated exactly (Pi #56 B1): this closes the user/project layer, not
//     every layer. Vendor-managed policy can remain higher-authority context
//     supplied by the execution substrate, and nothing here detects or
//     overrides it. So the honest claim is that the digested skills are the
//     only context THIS CELL contributes — not the only context that exists.
//
//   - `--tools` RESTRICTS the available built-in set. `--allowedTools` merely
//     pre-approves tools that remain available, so using it would have left
//     Bash reachable while claiming otherwise; it must never appear here.
//
//   - `--permission-mode acceptEdits` is what AUTHORIZES the edits this seat
//     exists to make. Availability and approval are two different things: with
//     the tool set restricted but no mode declared, Write and Edit are offered
//     and then not approved, so the same resolved cell would either fall back
//     on whatever ambient permission settings the host happens to carry or
//     produce no patch at all (Pi #56 D1). Sealing the mode here is what makes
//     the episode depend on the declaration rather than the environment. It
//     approves edits only — Bash is absent from the tool set entirely, and
//     bypassPermissions is never used.
//
//   - `--no-session-persistence` keeps the adapter stateless.
func ClaudeArgv(model string) []string {
	return []string{
		"-p",
		"--model", model,
		"--safe-mode",
		"--no-session-persistence",
		"--tools", ToolSurface,
		"--permission-mode", "acceptEdits",
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
