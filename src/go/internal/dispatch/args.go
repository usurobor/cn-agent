package dispatch

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseArgs parses command-line arguments for the dispatch command
//
// Expected format: cn dispatch --role α|β --branch cycle/N [--issue N] [--project NAME] [--backend claude|stub|print] [prompt-file]
func ParseArgs(args []string) (*Args, error) {
	result := &Args{}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "--role":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--role requires a value")
			}
			i++
			role := args[i]
			// Normalize α/β to full names
			switch role {
			case "α", "alpha":
				result.Role = "alpha"
			case "β", "beta":
				result.Role = "beta"
			default:
				return nil, fmt.Errorf("invalid role %q (use α|β or alpha|beta)", role)
			}

		case arg == "--branch":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--branch requires a value")
			}
			i++
			result.Branch = args[i]

		case arg == "--issue":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--issue requires a value")
			}
			i++
			issueStr := args[i]
			issue, err := strconv.Atoi(issueStr)
			if err != nil {
				return nil, fmt.Errorf("invalid issue number %q", issueStr)
			}
			result.Issue = issue

		case arg == "--project":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--project requires a value")
			}
			i++
			result.Project = args[i]

		case arg == "--backend":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--backend requires a value")
			}
			i++
			result.Backend = args[i]

		case strings.HasPrefix(arg, "--"):
			return nil, fmt.Errorf("unknown flag %q", arg)

		default:
			// Assume it's the prompt file
			if result.PromptFile != "" {
				return nil, fmt.Errorf("multiple prompt files specified")
			}
			result.PromptFile = arg
		}
	}

	// Validate required arguments
	if result.Role == "" {
		return nil, fmt.Errorf("--role is required")
	}
	if result.Branch == "" {
		return nil, fmt.Errorf("--branch is required")
	}

	// Validate branch format (should be cycle/N)
	if !strings.HasPrefix(result.Branch, "cycle/") {
		return nil, fmt.Errorf("branch %q should be cycle/N format", result.Branch)
	}

	// Extract issue number from branch if not provided via --issue
	if result.Issue == 0 {
		branchParts := strings.Split(result.Branch, "/")
		if len(branchParts) == 2 {
			if issue, err := strconv.Atoi(branchParts[1]); err == nil {
				result.Issue = issue
			}
		}
	}

	return result, nil
}