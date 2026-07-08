package cli

import (
	"bytes"
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

// Help implements HelpProvider (cnos#612 AC2) so main.go's generic
// --help/-h interception, which runs before Run, surfaces this text
// instead of the generic invocation+summary fallback.
func (c *ActivateCmd) Help() string {
	return activateHelp
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

	// Mutual exclusion fires before any rendering: an operator who set both
	// flags gets a clean error without diagnostic noise on stderr from the
	// renderer's "→ Generating activation prompt …" line.
	if claudeFlag && codexFlag {
		fmt.Fprintf(inv.Stderr, "✗ --claude and --codex are mutually exclusive — pass only one\n")
		return fmt.Errorf("activate: --claude and --codex are mutually exclusive")
	}

	// Hub resolution: explicit path wins; fall back to cwd discovery.
	hubPath := inv.HubPath
	if positionalHub != "" {
		hubPath = positionalHub
	}

	// Determine the spawn target. When set, this is the binary that will
	// replace the cn process image after the prompt is rendered.
	var spawnBinary, spawnFlag string
	switch {
	case claudeFlag:
		spawnBinary, spawnFlag = "claude", "--claude"
	case codexFlag:
		spawnBinary, spawnFlag = "codex", "--codex"
	}

	// Missing-binary detection runs before render so the operator who asked
	// for an interactive spawn but does not have the CLI installed sees the
	// pre-render error (no wasted render, no confusing diagnostic order).
	if spawnBinary != "" {
		if err := activate.CheckSpawnBinary(spawnBinary, spawnFlag); err != nil {
			fmt.Fprintf(inv.Stderr, "✗ %s\n", err)
			return err
		}
	}

	// Render. The default (no-flag) path writes straight to inv.Stdout —
	// bytes-equal to the 3.78.0 behavior pipe consumers depend on. The
	// spawn path captures into a buffer so the rendered prompt can be
	// handed to claude / codex as the bare-positional argv element.
	if spawnBinary == "" {
		return activate.Run(ctx, activate.Options{
			HubPath: hubPath,
			Stdout:  inv.Stdout,
			Stderr:  inv.Stderr,
		})
	}

	var buf bytes.Buffer
	if err := activate.Run(ctx, activate.Options{
		HubPath: hubPath,
		Stdout:  &buf,
		Stderr:  inv.Stderr,
	}); err != nil {
		return err
	}
	return activate.Spawn(spawnBinary, buf.String())
}

func isFlag(s string) bool {
	return len(s) > 0 && s[0] == '-'
}
