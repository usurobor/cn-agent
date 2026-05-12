package dispatch

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *Args
		wantErr  bool
	}{
		{
			name: "valid alpha dispatch",
			args: []string{"--role", "α", "--branch", "cycle/295"},
			expected: &Args{
				Role:   "alpha",
				Branch: "cycle/295",
				Issue:  295, // extracted from branch
			},
			wantErr: false,
		},
		{
			name: "valid beta dispatch with issue",
			args: []string{"--role", "beta", "--branch", "cycle/300", "--issue", "300"},
			expected: &Args{
				Role:   "beta",
				Branch: "cycle/300",
				Issue:  300,
			},
			wantErr: false,
		},
		{
			name: "with backend and project",
			args: []string{"--role", "α", "--branch", "cycle/295", "--backend", "stub", "--project", "cnos"},
			expected: &Args{
				Role:    "alpha",
				Branch:  "cycle/295",
				Issue:   295,
				Backend: "stub",
				Project: "cnos",
			},
			wantErr: false,
		},
		{
			name:    "missing role",
			args:    []string{"--branch", "cycle/295"},
			wantErr: true,
		},
		{
			name:    "missing branch",
			args:    []string{"--role", "alpha"},
			wantErr: true,
		},
		{
			name:    "invalid role",
			args:    []string{"--role", "gamma", "--branch", "cycle/295"},
			wantErr: true,
		},
		{
			name:    "invalid branch format",
			args:    []string{"--role", "alpha", "--branch", "main"},
			wantErr: true,
		},
		{
			name: "with prompt file",
			args: []string{"--role", "α", "--branch", "cycle/295", "custom-prompt.md"},
			expected: &Args{
				Role:       "alpha",
				Branch:     "cycle/295",
				Issue:      295,
				PromptFile: "custom-prompt.md",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseArgs(tt.args)

			if tt.wantErr && err == nil {
				t.Errorf("expected error but got none")
				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.wantErr {
				return // we expected an error and got one
			}

			// Compare results
			if result.Role != tt.expected.Role {
				t.Errorf("Role: got %q, want %q", result.Role, tt.expected.Role)
			}
			if result.Branch != tt.expected.Branch {
				t.Errorf("Branch: got %q, want %q", result.Branch, tt.expected.Branch)
			}
			if result.Issue != tt.expected.Issue {
				t.Errorf("Issue: got %d, want %d", result.Issue, tt.expected.Issue)
			}
			if result.Backend != tt.expected.Backend {
				t.Errorf("Backend: got %q, want %q", result.Backend, tt.expected.Backend)
			}
			if result.Project != tt.expected.Project {
				t.Errorf("Project: got %q, want %q", result.Project, tt.expected.Project)
			}
			if result.PromptFile != tt.expected.PromptFile {
				t.Errorf("PromptFile: got %q, want %q", result.PromptFile, tt.expected.PromptFile)
			}
		})
	}
}