package cli

import (
	"context"
	"fmt"

	"github.com/usurobor/cnos/src/go/internal/activate"
)

const activateHelp = `Usage: cn activate [--claude|--codex] [HUB_DIR]

Generate a bootstrap prompt from a local hub.

Arguments:
  HUB_DIR     path to the hub (optional; discovers hub from cwd when omitted)

Flags:
  --claude    after rendering, replace the cn process with an interactive
              ` + "`claude \"$prompt\"`" + ` session (TTY preserved; requires claude on $PATH)
  --codex     after rendering, replace the cn process with an interactive
              ` + "`codex \"$prompt\"`" + ` session (TTY preserved; requires codex on $PATH)

The default (no flag) writes the rendered prompt to stdout. Diagnostics go to
stderr in every mode. This command generates a prompt only — no model is
invoked unless --claude or --codex is set.

Examples:
  cn activate
  cn activate HUB_DIR
  cn activate HUB_DIR > activation.prompt.md
  cn activate --claude HUB_DIR
  cn activate --codex HUB_DIR
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
	// Parse args: --help/-h, --claude, --codex, and a single positional HUB_DIR.
	var claudeFlag, codexFlag bool
	positionalHub := ""
	for _, a := range inv.Args {
		switch a {
		case "--help", "-h":
			fmt.Fprint(inv.Stdout, activateHelp)
			return nil
		case "--claude":
			claudeFlag = true
		case "--codex":
			codexFlag = true
		default:
			if isFlag(a) {
				fmt.Fprintf(inv.Stderr, "✗ Unknown flag: %s\n", a)
				return fmt.Errorf("activate: unknown flag %s", a)
			}
			if positionalHub == "" {
				positionalHub = a
			}
		}
	}

	// Hub resolution: explicit path wins; fall back to cwd discovery.
	hubPath := inv.HubPath
	if positionalHub != "" {
		hubPath = positionalHub
	}

	_ = claudeFlag
	_ = codexFlag

	return activate.Run(ctx, activate.Options{
		HubPath: hubPath,
		Stdout:  inv.Stdout,
		Stderr:  inv.Stderr,
	})
}

func isFlag(s string) bool {
	return len(s) > 0 && s[0] == '-'
}
