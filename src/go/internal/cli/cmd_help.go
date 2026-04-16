package cli

import (
	"context"
	"fmt"
)

// HelpCmd implements the "help" command.
// It lists all available commands with their summaries.
type HelpCmd struct {
	Registry *Registry
}

func (c *HelpCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:     "help",
		Summary:  "Show available commands",
		Source:   SourceKernel,
		Tier:     TierKernel,
		NeedsHub: false, // help is always available
	}
}

func (c *HelpCmd) Run(_ context.Context, inv Invocation) error {
	// `cn help <noun>` — show the subcommand listing for a noun group.
	// Falls through to the full listing when <noun> is not a group, so
	// `cn help <unknown>` remains informative rather than erroring.
	if len(inv.Args) > 0 {
		if PrintGroup(inv.Stdout, c.Registry, inv.Args[0]) {
			return nil
		}
	}

	hasHub := inv.HubPath != ""

	// Kernel commands are always part of the binary and always listed —
	// users can see the full kernel surface without creating a hub.
	// Repo-local and package commands are only meaningful inside a hub.
	var kernelCmds, repoLocalCmds, packageCmds []Command
	for _, cmd := range c.Registry.All() {
		spec := cmd.Spec()
		switch spec.Source {
		case SourceKernel:
			kernelCmds = append(kernelCmds, cmd)
		case SourceRepoLocal:
			if hasHub {
				repoLocalCmds = append(repoLocalCmds, cmd)
			}
		case SourcePackage:
			if hasHub {
				packageCmds = append(packageCmds, cmd)
			}
		}
	}

	fmt.Fprintf(inv.Stdout, "cn — cnos kernel\n\n")
	fmt.Fprintf(inv.Stdout, "Usage: cn <command> [args...]\n\n")

	if len(kernelCmds) > 0 {
		fmt.Fprintf(inv.Stdout, "Kernel commands:\n")
		for _, cmd := range kernelCmds {
			spec := cmd.Spec()
			suffix := ""
			if !hasHub && spec.NeedsHub {
				suffix = "  (requires hub)"
			}
			fmt.Fprintf(inv.Stdout, "  %-12s %s%s\n", spec.Name, spec.Summary, suffix)
		}
	}

	if len(repoLocalCmds) > 0 {
		fmt.Fprintf(inv.Stdout, "\nRepo-local commands:\n")
		for _, cmd := range repoLocalCmds {
			spec := cmd.Spec()
			fmt.Fprintf(inv.Stdout, "  %-12s %s\n", spec.Name, spec.Summary)
		}
	}

	if len(packageCmds) > 0 {
		fmt.Fprintf(inv.Stdout, "\nPackage commands:\n")
		for _, cmd := range packageCmds {
			spec := cmd.Spec()
			fmt.Fprintf(inv.Stdout, "  %-12s %s", spec.Name, spec.Summary)
			if spec.Package != "" {
				fmt.Fprintf(inv.Stdout, "  [%s]", spec.Package)
			}
			fmt.Fprintln(inv.Stdout)
		}
	}

	if !hasHub {
		fmt.Fprintf(inv.Stdout, "\n⚠ No hub found — run 'cn init' to create one, or cd into an existing hub.\n")
	}

	return nil
}
