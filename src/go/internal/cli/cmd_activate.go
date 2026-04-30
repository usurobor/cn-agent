package cli

import (
	"context"
	"fmt"

	"github.com/usurobor/cnos/src/go/internal/activate"
)

const activateHelp = `Usage: cn activate [HUB_DIR]

Generate a bootstrap prompt from a local hub.

Arguments:
  HUB_DIR   path to the hub (optional; discovers hub from cwd when omitted)

The generated prompt goes to stdout. Diagnostics go to stderr.
This command generates a prompt only. No model is invoked.

Examples:
  cn activate
  cn activate HUB_DIR
  cn activate HUB_DIR > activation.prompt.md
  cn activate HUB_DIR | claude -p "Activate this cnos hub using the bootstrap prompt on stdin."
`

// ActivateCmd implements the "activate" command — generates a bootstrap prompt from hub state.
type ActivateCmd struct{}

func (c *ActivateCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:     "activate",
		Summary:  "Generate a bootstrap prompt from hub state",
		Source:   SourceKernel,
		Tier:     TierKernel,
		NeedsHub: false,
	}
}

func (c *ActivateCmd) Run(ctx context.Context, inv Invocation) error {
	// Handle --help / -h before any hub resolution.
	for _, a := range inv.Args {
		if a == "--help" || a == "-h" {
			fmt.Fprint(inv.Stdout, activateHelp)
			return nil
		}
	}

	// Hub resolution: explicit path wins; fall back to cwd discovery.
	hubPath := inv.HubPath
	if len(inv.Args) > 0 && !isFlag(inv.Args[0]) {
		hubPath = inv.Args[0]
	}

	return activate.Run(ctx, activate.Options{
		HubPath: hubPath,
		Stdout:  inv.Stdout,
		Stderr:  inv.Stderr,
	})
}

func isFlag(s string) bool {
	return len(s) > 0 && s[0] == '-'
}
