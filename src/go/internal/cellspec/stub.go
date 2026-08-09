package cellspec

import (
	"context"
	"fmt"
	"strings"

	"github.com/usurobor/cnos/src/go/internal/cellkernel"
)

// stubAlpha and stubBeta are Phase-1 placeholders: they exercise the full
// kernel closure with no rented cognition. stubAlpha "produces" a summary of
// the contract and mints one evidence ref per required-evidence entry (so V's
// binding check passes); stubBeta accepts. Phase 3 swaps these for seats that
// rent cognition through internal/dispatch.Backend.

type stubAlpha struct {
	skills []string
}

func (a stubAlpha) Produce(_ context.Context, c cellkernel.Contract) (cellkernel.AlphaResult, error) {
	matter := cellkernel.Matter{
		Data: fmt.Sprintf("stub-alpha produced for %q with skills [%s]", c.Goal, strings.Join(a.skills, ", ")),
	}
	refs := make([]cellkernel.EvidenceRef, 0, len(c.RequiredEvidence))
	for _, req := range c.RequiredEvidence {
		refs = append(refs, cellkernel.EvidenceRef{
			ID:                  req.ID,
			Kind:                req.Kind,
			Ref:                 "stub://" + req.ID,
			ProducerExecutionID: "stub-alpha",
		})
	}
	return cellkernel.AlphaResult{Matter: matter, EvidenceRefs: refs}, nil
}

type stubBeta struct {
	skills []string
}

func (b stubBeta) Review(_ context.Context, _ cellkernel.Contract, _ cellkernel.Matter) (cellkernel.BetaResult, error) {
	return cellkernel.BetaResult{
		Review: cellkernel.Review{
			Pass:  true,
			Notes: fmt.Sprintf("stub-beta accepted with skills [%s]", strings.Join(b.skills, ", ")),
		},
	}, nil
}
