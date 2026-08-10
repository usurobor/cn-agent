package cellcog

import (
	"strings"
	"testing"
)

func joined(args []string) string { return strings.Join(args, " ") }

// The Claude invocation must RESTRICT the tool surface, not merely
// pre-approve it, and must not inherit ambient configuration.
func TestClaudeArgvIsExact(t *testing.T) {
	got := ClaudeArgv("claude-opus-4-1")
	want := []string{
		"-p",
		"--model", "claude-opus-4-1",
		"--no-session-persistence",
		"--setting-sources", "",
		"--strict-mcp-config",
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
}

// Codex must run stateless and ignore ambient configuration, while
// authentication stays ambient.
func TestCodexArgvIsExact(t *testing.T) {
	got := CodexArgv("gpt-5-codex", "/w")
	want := []string{
		"exec",
		"--model", "gpt-5-codex",
		"--ephemeral",
		"--ignore-user-config",
		"--ignore-rules",
		"--sandbox", "workspace-write",
		"--skip-git-repo-check",
		"--cd", "/w",
		"-",
	}
	if joined(got) != joined(want) {
		t.Fatalf("argv drifted:\n got %q\nwant %q", got, want)
	}
}

func TestCodexArgvForbidsDangerousModes(t *testing.T) {
	got := joined(CodexArgv("m", "/w"))
	for _, forbidden := range []string{
		"danger-full-access", "--yolo", "--dangerously-bypass-approvals-and-sandbox",
		"--full-auto",
	} {
		if strings.Contains(got, forbidden) {
			t.Errorf("argv must not contain %q: %s", forbidden, got)
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
