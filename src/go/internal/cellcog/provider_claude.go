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
// spawning anything. Each flag, and the obvious alternative it rejects:
//
//   - `--safe-mode` closes the USER AND PROJECT layer (CLAUDE.md, skills,
//     plugins, hooks, MCP, commands, auto-memory), keeping auth and model:
//     repo-local guidance would be a second, unreceipted component definition
//     beside the fill's digested skills, so such guidance is declared AS one.
//   - `--tools` bounds what EXISTS; CodingToolSurface says why Bash is in it.
//   - `--allowedTools Bash` sits BESIDE `--tools`, never instead of it — the
//     substitution pre-approves without restricting, and the tests forbid it.
//     Required because `acceptEdits` does not reach ordinary commands:
//     measured, `go version` was denied under the mode alone, zero denials with
//     the flag. Without it the seat is offered a shell it may not use.
//   - `--permission-mode acceptEdits` AUTHORIZES those edits, since being
//     offered a tool is not approval to use it: undeclared, the same resolved
//     cell inherits ambient host permissions or produces no patch at all (Pi
//     #56 D1). Covers editing, not ordinary commands; bypassPermissions never.
//   - `--output-format stream-json --verbose` emits events AS THEY HAPPEN, so a
//     stall leaves a trail; `text` emits nothing until the run ends, which made
//     "captured 0 bytes" and "captured 40KB then stopped" identical. The CLI
//     refuses `--print` with `stream-json` unless `--verbose` is given.
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
	// stdout and truncation both ignored, for the reason on Coder.
	_, _, err := runCLI(ctx, bin, dir, prompt, ClaudeArgv(c.Model), c.Timeout)
	return err
}

// ClaudeAnswerArgv is the ANSWERING recipe, strictly less authoritative:
//
//   - `--tools ""` — a reviewer's canonical input is (contract, matter); file
//     tools would let it read from inside the workspace it must judge from the
//     outside, losing the independence the seat exists to provide.
//   - no `--permission-mode` — with no tools there is nothing to approve, so
//     declaring edit authority would request power the seat cannot use.
//   - `--json-schema` constrains the answer to the caller's shape, so the
//     verdict is decoded rather than parsed hopefully out of prose; the
//     terminal `result` event carries it, so it survives streaming.
//   - stream-json for the producing recipe's reason: a stall should say where.
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

// Answer returns `structured_output`; its absence is a failure, not an empty
// answer, since a reviewer that returned nothing has not reviewed. The
// directory is empty because an answering seat has no tools to scope.
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

// terminalStructuredOutput takes the schema-constrained value from the
// terminal `result` event. Decoded as a sequence of JSON values rather than
// scanned by line: an event carrying a long verdict can exceed any line-buffer
// size, and a diagnostic that fails on long input is not a diagnostic.
// Progress events are skipped — one that looks like an answer is not one.
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
