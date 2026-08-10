package cellcog

import (
	"strings"
	"testing"
)

func joined(args []string) string { return strings.Join(args, " ") }

// The Claude invocation must RESTRICT the tool surface, not merely
// pre-approve it, and must disable ambient customization — CLAUDE.md
// included — so the receipted skills are the only component definition.
func TestClaudeArgvIsExact(t *testing.T) {
	got := ClaudeArgv("claude-opus-4-1")
	want := []string{
		"-p",
		"--model", "claude-opus-4-1",
		"--safe-mode",
		"--no-session-persistence",
		"--tools", "Read,Write,Edit,Glob,Grep",
		"--output-format", "text",
	}
	if joined(got) != joined(want) {
		t.Fatalf("argv drifted:\n got %q\nwant %q", got, want)
	}
}

// `--allowedTools` only pre-approves tools that remain available; using it
// while claiming a restricted surface is exactly the defect this pins.
// `--dangerously-skip-permissions` and a Bash grant must never appear.
func TestClaudeArgvForbidsPreApprovalAndShell(t *testing.T) {
	got := joined(ClaudeArgv("m"))
	for _, forbidden := range []string{
		"--allowedTools", "--allowed-tools",
		"--dangerously-skip-permissions",
		"Bash", "--permission-mode",
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
