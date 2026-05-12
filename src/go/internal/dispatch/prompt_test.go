package dispatch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConstructPrompt(t *testing.T) {
	// Setup test directory
	testDir := t.TempDir()

	dispatcher := &Dispatcher{
		Project: "cnos",
		HubPath: testDir,
	}

	tests := []struct {
		name     string
		args     *Args
		expected []string // strings that should be present in the prompt
	}{
		{
			name: "alpha prompt",
			args: &Args{
				Role:   "alpha",
				Branch: "cycle/295",
				Issue:  295,
			},
			expected: []string{
				"You are α",
				"Project: cnos",
				"Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md",
				"Issue: gh issue view 295 --json title,body,state,comments",
				"Branch: cycle/295",
				"Dispatch mode: identity-rotation",
				"Do not start polling loops",
				"γ will re-dispatch you",
			},
		},
		{
			name: "beta prompt",
			args: &Args{
				Role:   "beta",
				Branch: "cycle/300",
				Issue:  300,
			},
			expected: []string{
				"You are β",
				"Project: cnos",
				"Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md",
				"Issue: gh issue view 300 --json title,body,state,comments",
				"Branch: cycle/300",
				"Dispatch mode: identity-rotation",
				"Do not start polling loops",
				"γ will re-dispatch you",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt, err := dispatcher.constructPrompt(tt.args)
			if err != nil {
				t.Errorf("constructPrompt failed: %v", err)
				return
			}

			for _, expected := range tt.expected {
				if !strings.Contains(prompt, expected) {
					t.Errorf("Expected %q in prompt, got:\n%s", expected, prompt)
				}
			}
		})
	}
}

func TestConstructPromptWithFile(t *testing.T) {
	// Setup test directory and prompt file
	testDir := t.TempDir()
	promptFile := filepath.Join(testDir, "custom-prompt.md")
	customPrompt := "This is a custom prompt for testing"

	if err := os.WriteFile(promptFile, []byte(customPrompt), 0644); err != nil {
		t.Fatalf("Failed to write test prompt file: %v", err)
	}

	dispatcher := &Dispatcher{
		Project: "cnos",
		HubPath: testDir,
	}

	args := &Args{
		Role:       "alpha",
		Branch:     "cycle/295",
		Issue:      295,
		PromptFile: promptFile,
	}

	prompt, err := dispatcher.constructPrompt(args)
	if err != nil {
		t.Errorf("constructPrompt failed: %v", err)
		return
	}

	if prompt != customPrompt {
		t.Errorf("Expected custom prompt content, got: %s", prompt)
	}
}

func TestConstructPromptInvalidRole(t *testing.T) {
	dispatcher := &Dispatcher{
		Project: "cnos",
	}

	args := &Args{
		Role:   "gamma", // invalid role
		Branch: "cycle/295",
		Issue:  295,
	}

	_, err := dispatcher.constructPrompt(args)
	if err == nil {
		t.Errorf("Expected error for invalid role")
	}

	if !strings.Contains(err.Error(), "unknown role") {
		t.Errorf("Expected 'unknown role' error, got: %v", err)
	}
}

func TestConstructPromptMissingFile(t *testing.T) {
	dispatcher := &Dispatcher{
		Project: "cnos",
	}

	args := &Args{
		Role:       "alpha",
		Branch:     "cycle/295",
		Issue:      295,
		PromptFile: "nonexistent-file.md",
	}

	_, err := dispatcher.constructPrompt(args)
	if err == nil {
		t.Errorf("Expected error for missing prompt file")
	}
}