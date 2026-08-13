package cellcog

import (
	"context"
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
//   - `--tools` sets the available built-in set — see CodingToolSurface for
//     why it includes Bash and why it is a capability declaration rather than
//     a boundary.
//
//   - `--allowedTools Bash` sits BESIDE `--tools`, never instead of it. The
//     two flags answer different questions: `--tools` bounds what exists,
//     `--allowedTools` grants approval for what already exists to run without
//     asking. Using `--allowedTools` as a REPLACEMENT for `--tools` was a real
//     defect once — it pre-approves without restricting, so the surface stayed
//     open while the comment claimed otherwise — and that substitution is what
//     the tests still forbid.
//
//     It is required because `--permission-mode acceptEdits` does not reach
//     ordinary commands. Measured against the installed CLI on the donor
//     branch: `echo X > file` ran under the mode alone, while `go version`
//     came back denied (`permission_denials: [{"tool_name":"Bash",…,
//     "command":"go version"}]`); adding this flag brought that to zero
//     denials with the command actually executing. Without it a seat is
//     offered the shell it needs to verify its own work and then refused
//     permission to use it — which is the same unverifying seat the surface
//     was widened to end.
//
//   - `--permission-mode acceptEdits` is what AUTHORIZES the edits this seat
//     exists to make. Availability and approval are two different things: with
//     the tool set restricted but no mode declared, Write and Edit are offered
//     and then not approved, so the same resolved cell would either fall back
//     on whatever ambient permission settings the host happens to carry or
//     produce no patch at all (Pi #56 D1). Sealing the mode makes the
//     BASELINE explicit rather than inherited from user or project defaults.
//     It covers file edits and the small filesystem command set the mode
//     treats as part of editing; it does NOT cover ordinary commands, which is
//     why `--allowedTools Bash` is declared above. An earlier version of this
//     comment claimed the mode covered Bash outright, on evidence that only
//     ever exercised a redirect into a file — measurement against the CLI
//     falsified it. bypassPermissions is never used.
//
//     Not more than that (Pi #59 B1): this does not make the episode
//     independent of the environment. Managed substrate policy remains above
//     the declared baseline, and nothing here detects or overrides it.
//
//   - `--output-format stream-json --verbose` emits events AS THEY HAPPEN
//     instead of one dump after the run ends. The bytes are still discarded
//     on success — the product is the measured diff — but a stalled provider
//     now leaves a trail: "captured 0 bytes" and "captured 40KB then stopped"
//     are different diagnoses, and under `text` both looked identical because
//     print mode emits nothing until it finishes. `--verbose` is not optional
//     here; the CLI refuses `--print` with `stream-json` without it.
//
//   - `--no-session-persistence` keeps the adapter stateless.
func ClaudeArgv(model string) []string {
	return []string{
		"-p",
		"--model", model,
		"--safe-mode",
		"--no-session-persistence",
		"--tools", CodingToolSurface,
		"--allowedTools", "Bash",
		"--permission-mode", "acceptEdits",
		"--output-format", "stream-json",
		"--verbose",
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
	// worktree diff, never by its own account of what it did. Truncation is
	// tolerated for the same reason — the stream is progress, so losing its
	// tail costs nothing the episode depends on.
	_, _, err := runCLI(ctx, bin, dir, prompt, ClaudeArgv(c.Model), c.Timeout)
	return err
}
