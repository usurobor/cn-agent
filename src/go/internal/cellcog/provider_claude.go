package cellcog

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// ClaudeCLI rents workspace cognition from the Claude Code CLI. Stateless: a
// Work call is one fresh bounded subprocess; no session is started or kept.
type ClaudeCLI struct {
	Model   string        // requested model selector; required
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
//     produce no patch at all (Pi #56 D1). Sealing the mode makes the
//     BASELINE explicit rather than inherited from user or project defaults.
//     It approves edits only — Bash is absent from the tool set entirely, and
//     bypassPermissions is never used.
//
//     Not more than that (Pi #59 B1): this does not make the episode
//     independent of the environment. Managed substrate policy remains above
//     the declared baseline, and nothing here detects or overrides it.
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
	// stdout is discarded on purpose: a producing seat is judged by the
	// worktree diff, never by its own account of what it did.
	_, err := runCLI(ctx, bin, dir, prompt, ClaudeArgv(c.Model), c.Timeout)
	return err
}

// ClaudeAnswerArgv is the ANSWERING recipe. It differs from the producing one
// in exactly the ways the capability differs, and is strictly less
// authoritative:
//
//   - `--tools ""` — a reviewer's canonical input is (contract, matter). File
//     tools would let it read the workspace it is meant to judge from the
//     outside, which is the independence the seat exists to provide.
//   - no `--permission-mode` — with no tools there is nothing to approve, so
//     declaring edit authority would be requesting power the seat cannot use.
//   - `--output-format json` + `--json-schema` — the provider constrains the
//     answer to the caller's shape, so the verdict is decoded rather than
//     parsed hopefully out of prose.
func ClaudeAnswerArgv(model string, schema json.RawMessage) []string {
	return []string{
		"-p",
		"--model", model,
		"--safe-mode",
		"--no-session-persistence",
		"--tools", NoTools,
		"--output-format", "json",
		"--json-schema", string(schema),
	}
}

// Answer runs one bounded invocation and returns the provider's structured
// result. The envelope carries execution metadata; `structured_output` is the
// schema-constrained value, and its absence is a failure rather than an empty
// answer — a reviewer that returned nothing has not reviewed.
func (c ClaudeCLI) Answer(ctx context.Context, prompt string, schema json.RawMessage) (json.RawMessage, error) {
	if len(schema) == 0 {
		return nil, fmt.Errorf("claude-cli: Answer needs an answer schema")
	}
	bin := c.Bin
	if bin == "" {
		bin = "claude"
	}
	out, err := runCLI(ctx, bin, "", prompt, ClaudeAnswerArgv(c.Model, schema), c.Timeout)
	if err != nil {
		return nil, err
	}
	var env struct {
		IsError    bool            `json:"is_error"`
		Structured json.RawMessage `json:"structured_output"`
	}
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		return nil, fmt.Errorf("claude-cli: result envelope is not JSON: %w", err)
	}
	if env.IsError {
		return nil, fmt.Errorf("claude-cli: provider reported an error result")
	}
	if len(env.Structured) == 0 {
		return nil, fmt.Errorf("claude-cli: provider returned no structured_output for the requested schema")
	}
	return env.Structured, nil
}
