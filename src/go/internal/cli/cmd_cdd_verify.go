package cli

import (
	"context"
	"fmt"

	cddverify "github.com/usurobor/cnos/packages/cnos.cdd/commands/cdd-verify"
)

// CddVerifyCmd implements the "cdd-verify" command — V (Contract × Receipt
// → ValidationVerdict) plus the legacy CDD artifact-presence ledger.
//
// Per eng/go §2.18 (dispatch boundary), this file owns only the dispatch
// wiring; the domain logic lives in the cddverify package at
// src/packages/cnos.cdd/commands/cdd-verify/ (per cnos#392's package-scoping
// pin).
type CddVerifyCmd struct{}

func (c *CddVerifyCmd) Spec() CommandSpec {
	return CommandSpec{
		Name:    "cdd-verify",
		Summary: "Verify CDD cycle artifact completeness and validate V receipts",
		Source:  SourceKernel,
		Tier:    TierKernel,
		// NeedsHub is false: --receipt mode runs against arbitrary receipt
		// paths; ledger modes that need a hub fail with a clearer error from
		// the cddverify package than the kernel hub-required gate.
	}
}

func (c *CddVerifyCmd) Run(ctx context.Context, inv Invocation) error {
	code := cddverify.Run(ctx, inv.Args, inv.Stdout, inv.Stderr)
	if code == cddverify.ExitPASS {
		return nil
	}
	// Encode the exit code in the error so cmd/cn/main.go can propagate it.
	return &CddVerifyExit{Code: int(code)}
}

// CddVerifyExit carries V's exit code (0/1/2) back to the cn process so
// it can exit with the same code. main.go inspects this concretely.
type CddVerifyExit struct{ Code int }

func (e *CddVerifyExit) Error() string {
	return fmt.Sprintf("cdd-verify exited with status %d", e.Code)
}
