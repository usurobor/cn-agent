package cellcog

import (
	"encoding/json"
	"strings"
	"testing"
)

func joined(args []string) string { return strings.Join(args, " ") }

// The Claude invocation must RESTRICT the tool surface, not merely
// pre-approve it; must disable USER AND PROJECT customization — CLAUDE.md
// included — so the receipted skills are the only component definition this
// cell contributes (vendor-managed policy is a substrate concern this makes
// no claim about); and must AUTHORIZE the edits the seat exists to make
// rather than inheriting that authority from the host.
func TestClaudeArgvIsExact(t *testing.T) {
	got := ClaudeArgv("claude-opus-4-1")
	want := []string{
		"-p",
		"--model", "claude-opus-4-1",
		"--safe-mode",
		"--no-session-persistence",
		"--tools", "Read,Write,Edit,Glob,Grep",
		"--permission-mode", "acceptEdits",
		"--output-format", "stream-json",
		"--verbose",
	}
	if joined(got) != joined(want) {
		t.Fatalf("argv drifted:\n got %q\nwant %q", got, want)
	}
}

// Availability is not approval: a seat offered Write and Edit but given no
// permission mode depends on ambient host settings for whether it may use
// them (Pi #56 D1). The mode must be present AND must be the edit-scoped one.
func TestClaudeArgvAuthorizesEditsExplicitly(t *testing.T) {
	got := joined(ClaudeArgv("m"))
	if !strings.Contains(got, "--permission-mode acceptEdits") {
		t.Errorf("argv must authorize edits with --permission-mode acceptEdits: %s", got)
	}
	for _, tool := range []string{"Write", "Edit"} {
		if !strings.Contains(got, tool) {
			t.Errorf("argv must offer %s for a patch-producing seat: %s", tool, got)
		}
	}
}

// `--allowedTools` only pre-approves tools that remain available; using it
// while claiming a restricted surface is exactly the defect this pins.
// Bypass modes and a Bash grant must never appear — acceptEdits approves
// edits, it does not widen the surface.
func TestClaudeArgvForbidsPreApprovalAndShell(t *testing.T) {
	got := joined(ClaudeArgv("m"))
	for _, forbidden := range []string{
		"--allowedTools", "--allowed-tools",
		"--dangerously-skip-permissions", "--allow-dangerously-skip-permissions",
		"bypassPermissions", "dontAsk",
		"Bash",
	} {
		if strings.Contains(got, forbidden) {
			t.Errorf("argv must not contain %q: %s", forbidden, got)
		}
	}
	if !strings.Contains(got, "--tools ") {
		t.Errorf("argv must restrict the tool surface with --tools: %s", got)
	}
	if !strings.Contains(got, "--safe-mode") {
		t.Errorf("argv must disable ambient customization with --safe-mode: %s", got)
	}
	// Exactly one permission mode may be declared: a second would let a later
	// edit quietly widen authority by appending.
	if n := strings.Count(got, "--permission-mode"); n != 1 {
		t.Errorf("argv must declare exactly one permission mode, got %d: %s", n, got)
	}
}

// codex-cli is held, and the hold has to be a fact about the code rather than
// a note in a document: its suppression flags do not cover AGENTS.md or
// discovered skills, so admitting it would let ambient instructions stand
// beside the fill's digested skills (Pi #55 D1). A model id must not make it
// constructible either.
func TestCodexIsHeld(t *testing.T) {
	for _, model := range []string{"", "gpt-5-codex"} {
		if _, _, err := New(Config{Provider: "codex-cli", Model: model}); err == nil {
			t.Fatalf("codex-cli must not construct while held (model %q)", model)
		}
	}
}

// The fake rents nothing, so a model id it would ignore must not be
// receipted as though it selected something. Identical rule in CUE.
func TestFakeRejectsModel(t *testing.T) {
	if _, _, err := New(Config{Provider: "fake", Model: "claude-opus-4-1"}); err == nil {
		t.Fatal("fake + nonempty model must fail construction")
	}
	if _, _, err := New(Config{Provider: "fake"}); err != nil {
		t.Fatalf("fake without a model must construct: %v", err)
	}
}

// The ANSWERING recipe must be strictly less authoritative than the producing
// one, and must stream for the same reason: a stalled seat should say where.
func TestClaudeAnswerArgvIsExact(t *testing.T) {
	schema := json.RawMessage(`{"type":"object"}`)
	got := ClaudeAnswerArgv("claude-opus-5", schema)
	want := []string{
		"-p",
		"--model", "claude-opus-5",
		"--safe-mode",
		"--no-session-persistence",
		"--tools", "",
		"--output-format", "stream-json",
		"--verbose",
		"--json-schema", `{"type":"object"}`,
	}
	if joined(got) != joined(want) {
		t.Fatalf("answer argv drifted:\n got %q\nwant %q", got, want)
	}
}

// A reviewer is offered NO tools and therefore declares no permission mode:
// its canonical input is (contract, matter), and a tool would let it read the
// workspace it is meant to judge from outside. Asking for edit authority it
// cannot use would be requesting power for nothing.
func TestClaudeAnswerArgvIsLessAuthoritativeThanProducing(t *testing.T) {
	answer := joined(ClaudeAnswerArgv("m", json.RawMessage(`{}`)))
	for _, forbidden := range []string{
		"--permission-mode", "acceptEdits", "bypassPermissions",
		"Read", "Write", "Edit", "Glob", "Grep", "Bash",
		"--allowedTools", "--dangerously-skip-permissions",
	} {
		if strings.Contains(answer, forbidden) {
			t.Errorf("answering argv must not contain %q: %s", forbidden, answer)
		}
	}
	if !strings.Contains(answer, "--safe-mode") {
		t.Errorf("answering argv must still suppress user/project context: %s", answer)
	}
}

// Only the TERMINAL result event is an answer. A progress event that happens
// to carry a structured payload must not be mistaken for the verdict, and a
// stream that never reaches a result has not answered at all.
func TestTerminalStructuredOutput(t *testing.T) {
	cases := []struct {
		name   string
		stream string
		want   string // want=="" means success
		value  string
	}{
		{
			name:   "terminal result wins over earlier events",
			stream: `{"type":"assistant","structured_output":{"pass":true}}` + "\n" + `{"type":"result","is_error":false,"structured_output":{"pass":false}}`,
			value:  `{"pass":false}`,
		},
		{name: "no result event", stream: `{"type":"assistant"}`, want: "no result event"},
		{name: "error result", stream: `{"type":"result","is_error":true}`, want: "error result"},
		{name: "result without structured output", stream: `{"type":"result","is_error":false}`, want: "no structured_output"},
		{name: "not ndjson", stream: `{"type":"result"} <garbage>`, want: "not NDJSON"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := terminalStructuredOutput(tc.stream)
			if tc.want == "" {
				if err != nil {
					t.Fatalf("want success, got %v", err)
				}
				if string(got) != tc.value {
					t.Fatalf("got %s, want %s", got, tc.value)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want error mentioning %q, got %v", tc.want, err)
			}
		})
	}
}
