package cellcog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

func joined(args []string) string { return strings.Join(args, " ") }

// flagValue returns the value following the last occurrence of name. Last, not
// first: an argv that declares a flag twice is decided by what the CLI would
// honour, and a test that read the first one would miss a widening appended
// later.
func flagValue(args []string, name string) (string, bool) {
	for i := len(args) - 2; i >= 0; i-- {
		if args[i] == name {
			return args[i+1], true
		}
	}
	return "", false
}

func mustFlagValue(t *testing.T, args []string, name string) string {
	t.Helper()
	value, ok := flagValue(args, name)
	if !ok {
		t.Fatalf("argv carries no %s with a value: %q", name, args)
	}
	return value
}

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
		"--tools", "Read,Write,Edit,MultiEdit,Glob,Grep,Bash",
		"--allowedTools", "Bash",
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

// `--allowedTools` only pre-approves tools that remain available. Used
// INSTEAD OF `--tools` that is exactly the defect this pins — the surface
// stays open while the recipe claims to be restricted — so what must never
// appear is a pre-approval flag without the restricting flag beside it, not
// the pre-approval flag itself. Bypass modes must never appear at all. Bash is
// not forbidden either: it is the capability a software-development seat needs
// to verify its own work.
func TestClaudeArgvForbidsPreApprovalWithoutRestriction(t *testing.T) {
	got := joined(ClaudeArgv("m"))
	for _, forbidden := range []string{
		"--dangerously-skip-permissions", "--allow-dangerously-skip-permissions",
		"bypassPermissions", "dontAsk",
	} {
		if strings.Contains(got, forbidden) {
			t.Errorf("argv must not contain %q: %s", forbidden, got)
		}
	}
	if !strings.Contains(got, "--tools ") {
		t.Errorf("argv must restrict the tool surface with --tools: %s", got)
	}
	// Pre-approval may only name tools the surface actually offers. Approving
	// something `--tools` withheld would widen the recipe by the back door,
	// and this is the containment of the substitution defect that survives:
	// checking the RELATIONSHIP between the two values, member by member.
	// (An earlier version asked whether a pre-approval flag appeared without
	// `--tools` — which the presence check above already decides, so no argv
	// could ever reach it.)
	surface := strings.Split(mustFlagValue(t, ClaudeArgv("m"), "--tools"), ",")
	for _, name := range []string{"--allowedTools", "--allowed-tools"} {
		value, ok := flagValue(ClaudeArgv("m"), name)
		if !ok {
			continue
		}
		for _, approved := range strings.Split(value, ",") {
			if !slices.Contains(surface, approved) {
				t.Errorf("%s pre-approves %q, which the offered surface %q does not include", name, approved, surface)
			}
		}
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

// The producing surface must match the live cnos-cds-dispatch allow-list. A
// cell mechanizes what the operator does by hand, so a seat that cannot run
// its own tests is a weaker Claude than the workflow it replaces — which is
// how Bash came to be missing in the first place.
//
// The workflow is READ, not restated. An earlier version of this test named
// the same seven tools in a literal beside CodingToolSurface and called that
// parity; two hand-copied lists prove only that someone typed the same thing
// twice, and both could drift together while the test stayed green. Only one
// of the two sides may be written here, and it is not the source of truth.
func TestProducingSurfaceMatchesLiveDispatch(t *testing.T) {
	live := dispatchAllowList(t)
	surface := strings.Split(CodingToolSurface, ",")
	slices.Sort(live)
	slices.Sort(surface)
	if !slices.Equal(live, surface) {
		t.Errorf("producing surface has drifted from the live dispatch allow-list:\n surface %q\n    live %q", surface, live)
	}
	// The answering surface stays empty: that one IS load-bearing.
	if NoTools != "" {
		t.Errorf("a reviewing seat must be offered no tools, got %q", NoTools)
	}
}

// dispatchAllowList returns settings.permissions.allow from the live
// cnos-cds-dispatch workflow. Every way of not finding it is a FAILURE: a
// parity guard that skips when it cannot locate its source of truth passes
// vacuously, which is the exact defect the caller exists to close.
//
// The repo root comes from this file's own position rather than the working
// directory, so the test does not depend on where `go test` was invoked from
// (the same idiom as repoinstall's closure test and cli's install test).
func dispatchAllowList(t *testing.T) []string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	// thisFile: <root>/src/go/internal/cellcog/argv_test.go
	root := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..")
	path := filepath.Join(root, ".github", "workflows", "cnos-cds-dispatch.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read the live dispatch workflow: %v", err)
	}
	block, err := yamlLiteralBlock(string(data), "settings:")
	if err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	// The action's `settings` input is JSON carried inside the YAML block, so
	// the block is decoded as what it is rather than pattern-matched.
	var settings struct {
		Permissions struct {
			Allow []string `json:"allow"`
		} `json:"permissions"`
	}
	if err := json.Unmarshal([]byte(block), &settings); err != nil {
		t.Fatalf("%s: settings block is not the JSON the action consumes: %v", path, err)
	}
	if len(settings.Permissions.Allow) == 0 {
		t.Fatalf("%s: settings.permissions.allow is absent or empty, so there is nothing to compare against", path)
	}
	return settings.Permissions.Allow
}

// yamlLiteralBlock returns the body lines of the `key: |` literal block, still
// carrying their indentation — the caller decodes JSON, which ignores it. The
// block ends where a line dedents back to the key's own level.
//
// Hand-rolled rather than pulling in a YAML dependency: this reads one known
// key out of one rendered file, and it reports every failure to find that key
// instead of returning an empty result the caller could mistake for absence.
func yamlLiteralBlock(doc, key string) (string, error) {
	lines := strings.Split(doc, "\n")
	start, keyIndent := -1, 0
	for i, line := range lines {
		trimmed := strings.TrimLeft(line, " ")
		if trimmed == key+" |" || trimmed == key+" |-" {
			start, keyIndent = i+1, len(line)-len(trimmed)
			break
		}
	}
	if start < 0 {
		return "", fmt.Errorf("no %q literal block found", key)
	}
	var body []string
	for _, line := range lines[start:] {
		if strings.TrimSpace(line) == "" {
			body = append(body, "")
			continue
		}
		if len(line)-len(strings.TrimLeft(line, " ")) <= keyIndent {
			break
		}
		body = append(body, line)
	}
	if strings.TrimSpace(strings.Join(body, "\n")) == "" {
		return "", fmt.Errorf("%q literal block is empty", key)
	}
	return strings.Join(body, "\n"), nil
}
