package cli

import (
	"context"
	"fmt"

	"github.com/usurobor/cnos/src/go/internal/cellfills"
	"github.com/usurobor/cnos/src/go/internal/cellrun"
)

// CellRunCmd implements "cn cell run" — the GitHub-free local cell runner. This
// file is dispatch only (eng/go §2.18); the IO, parsing, and rendering live in
// internal/cellrun.
type CellRunCmd struct{}

func (c *CellRunCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:    "cell-run",
		Summary: "Run one local cell episode from a serialized spec (no GitHub/network)",
		Source:  SourceKernel,
		Tier:    TierKernel,
	}
}

func (c *CellRunCmd) Help() string { return cellrun.Help }

func (c *CellRunCmd) Run(ctx context.Context, inv Invocation) error {
	code := cellrun.Run(ctx, cellfills.Assemble(inv.HubPath), inv.Args, inv.Stdin, inv.Stdout, inv.Stderr)
	if code == 0 {
		return nil
	}
	return &CellRunExit{Code: code}
}

// CellRunExit carries the runner's exit code back to cmd/cn/main.go.
type CellRunExit struct{ Code int }

func (e *CellRunExit) Error() string {
	return fmt.Sprintf("cell run exited with status %d", e.Code)
}
