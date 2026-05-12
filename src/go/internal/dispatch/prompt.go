package dispatch

import (
	"fmt"
	"os"
	"path/filepath"
)

// constructPrompt builds the role-appropriate prompt with dispatch-mode clause
func (d *Dispatcher) constructPrompt(args *Args) (string, error) {
	// If prompt file is provided, read it
	if args.PromptFile != "" {
		content, err := os.ReadFile(args.PromptFile)
		if err != nil {
			return "", fmt.Errorf("read prompt file %q: %w", args.PromptFile, err)
		}
		return string(content), nil
	}

	// Construct role-specific prompt
	switch args.Role {
	case "alpha":
		return d.constructAlphaPrompt(args), nil
	case "beta":
		return d.constructBetaPrompt(args), nil
	default:
		return "", fmt.Errorf("unknown role %q", args.Role)
	}
}

// constructAlphaPrompt builds the α role prompt with dispatch-mode clause
func (d *Dispatcher) constructAlphaPrompt(args *Args) string {
	prompt := fmt.Sprintf(`You are α.
Project: %s.
Load %s and follow its load order.
Issue: gh issue view %d --json title,body,state,comments
Branch: %s

Dispatch mode: identity-rotation.
Perform the bounded role step requested here, write required artifacts to the cycle branch, and return.
Do not start polling loops or wait for future role updates.
γ will re-dispatch you if another step is required.`,
		d.Project,
		d.getAlphaSkillPath(),
		args.Issue,
		args.Branch)

	return prompt
}

// constructBetaPrompt builds the β role prompt with dispatch-mode clause
func (d *Dispatcher) constructBetaPrompt(args *Args) string {
	prompt := fmt.Sprintf(`You are β.
Project: %s.
Load %s and follow its load order.
Issue: gh issue view %d --json title,body,state,comments
Branch: %s

Dispatch mode: identity-rotation.
Perform the bounded role step requested here, write required artifacts to the cycle branch, and return.
Do not start polling loops or wait for future role updates.
γ will re-dispatch you if another step is required.`,
		d.Project,
		d.getBetaSkillPath(),
		args.Issue,
		args.Branch)

	return prompt
}

// getAlphaSkillPath returns the path to the alpha skill file
func (d *Dispatcher) getAlphaSkillPath() string {
	// Look for alpha skill in the standard location
	skillPath := "src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md"
	if d.HubPath != "" {
		fullPath := filepath.Join(d.HubPath, skillPath)
		if _, err := os.Stat(fullPath); err == nil {
			return skillPath
		}
	}
	// Fallback to generic path
	return "src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md"
}

// getBetaSkillPath returns the path to the beta skill file
func (d *Dispatcher) getBetaSkillPath() string {
	// Look for beta skill in the standard location
	skillPath := "src/packages/cnos.cdd/skills/cdd/beta/SKILL.md"
	if d.HubPath != "" {
		fullPath := filepath.Join(d.HubPath, skillPath)
		if _, err := os.Stat(fullPath); err == nil {
			return skillPath
		}
	}
	// Fallback to generic path
	return "src/packages/cnos.cdd/skills/cdd/beta/SKILL.md"
}