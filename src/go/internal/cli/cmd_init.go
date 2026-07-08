package cli

import (
	"context"
	"fmt"

	"github.com/usurobor/cnos/src/go/internal/hubinit"
)

// InitCmd implements the "init" command — creates a new hub.
type InitCmd struct{}

func (c *InitCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:     "init",
		Summary:  "Initialize a new hub",
		Source:   SourceKernel,
		Tier:     TierKernel,
		NeedsHub: false,
	}
}

func (c *InitCmd) Run(ctx context.Context, inv Invocation) error {
	var hubName string
	if len(inv.Args) > 0 {
		// cnos#612 AC3: init takes no flags. Before this fix, an
		// unrecognized "--flag" (e.g. "--help", reaching here whenever a
		// caller invokes InitCmd.Run directly rather than through main.go's
		// generic --help interception) was treated as the positional hub
		// name — validHubName allows '-', so "--help" scaffolded a stray
		// "cn---help/" directory. Refuse instead of scaffolding.
		if isFlag(inv.Args[0]) {
			fmt.Fprintf(inv.Stderr, "✗ cn init: unrecognized flag %q\n\n", inv.Args[0])
			fmt.Fprintf(inv.Stderr, "Usage: cn init [HUB_NAME]\n\ninit takes no flags; HUB_NAME is an optional positional argument.\n")
			return fmt.Errorf("init: unrecognized flag %q", inv.Args[0])
		}
		hubName = inv.Args[0]
	}
	// If no name given, hubinit.Run derives it from cwd.
	return hubinit.Run(ctx, hubName, inv.Stdout, inv.Stderr)
}
