package cellcog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
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
//     a boundary. `--allowedTools` is a different flag that only PRE-APPROVES
//     tools already available; it is absent because `--permission-mode`
//     already grants approval, and because using it *instead of* `--tools`
//     was a real defect once: it left the surface unrestricted while the
//     comment claimed otherwise.
//
//   - `--permission-mode acceptEdits` is what AUTHORIZES the edits this seat
//     exists to make. Availability and approval are two different things: with
//     the tool set restricted but no mode declared, Write and Edit are offered
//     and then not approved, so the same resolved cell would either fall back
//     on whatever ambient permission settings the host happens to carry or
//     produce no patch at all (Pi #56 D1). Sealing the mode makes the
//     BASELINE explicit rather than inherited from user or project defaults.
//     It covers Bash as well as file edits — verified against the CLI, which
//     ran a Bash command under this mode with `permission_denials: []` — so a
//     seat can run its own tests without a second pre-approval flag.
//     bypassPermissions is never used.
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

// ClaudeAnswerArgv is the ANSWERING recipe. It differs from the producing one
// in exactly the ways the capability differs, and is strictly less
// authoritative:
//
//   - `--tools ""` — a reviewer's canonical input is (contract, matter). File
//     tools would let it read the workspace it is meant to judge from the
//     outside, which is the independence the seat exists to provide.
//   - no `--permission-mode` — with no tools there is nothing to approve, so
//     declaring edit authority would be requesting power the seat cannot use.
//   - `--json-schema` — the provider constrains the answer to the caller's
//     shape, so the verdict is decoded rather than parsed hopefully out of
//     prose. It survives streaming: the terminal `result` event carries
//     `structured_output`.
//   - `--output-format stream-json --verbose` for the same reason as the
//     producing recipe: a reviewer that stalls should say where.
func ClaudeAnswerArgv(model string, schema json.RawMessage) []string {
	return []string{
		"-p",
		"--model", model,
		"--safe-mode",
		"--no-session-persistence",
		"--tools", NoTools,
		"--output-format", "stream-json",
		"--verbose",
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
	out, truncated, err := runCLI(ctx, bin, "", prompt, ClaudeAnswerArgv(c.Model, schema), c.Timeout)
	if err != nil {
		return nil, err
	}
	// Here the stream IS the product, so a clipped one is fatal: the terminal
	// result event may be exactly what was lost.
	if truncated {
		return nil, fmt.Errorf("claude-cli: answer stream exceeded %d bytes, so the verdict may be incomplete", maxOutputBytes)
	}
	return terminalStructuredOutput(out)
}

// terminalStructuredOutput reads the NDJSON event stream and returns the
// schema-constrained value from the terminal `result` event.
//
// Decoded as a sequence of JSON values rather than scanned by line: an event
// carrying a long verdict can exceed any line-buffer size, and a diagnostic
// that fails on long input is not a diagnostic. Progress events are skipped —
// only the terminal result is an answer, and an intermediate assistant
// message that happens to look like one is not.
func terminalStructuredOutput(stream string) (json.RawMessage, error) {
	dec := json.NewDecoder(strings.NewReader(stream))
	var (
		found      bool
		isError    bool
		structured json.RawMessage
	)
	for {
		var ev struct {
			Type       string          `json:"type"`
			IsError    bool            `json:"is_error"`
			Structured json.RawMessage `json:"structured_output"`
		}
		if err := dec.Decode(&ev); err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("claude-cli: event stream is not NDJSON: %w", err)
		}
		if ev.Type == "result" {
			found, isError, structured = true, ev.IsError, ev.Structured
		}
	}
	switch {
	case !found:
		return nil, fmt.Errorf("claude-cli: event stream ended with no result event")
	case isError:
		return nil, fmt.Errorf("claude-cli: provider reported an error result")
	case len(structured) == 0:
		return nil, fmt.Errorf("claude-cli: result carried no structured_output for the requested schema")
	}
	return structured, nil
}
