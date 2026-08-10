// Package cellcog rents cognition for a cell seat.
//
// The kernel is cognition-free by construction: it knows only the `Alpha` and
// `Beta` function shapes and never learns that a model exists. This package
// implements those seats over a narrow provider port, keeping the split
// eng/write-functional requires:
//
//	pure core   RenderAlphaPrompt(contract, skills) → prompt
//	            ParseAlphaResponse(bytes)           → AlphaOutput
//	effect edge Provider.Complete(ctx, prompt)      → raw text  (adapters only)
//
// A rented seat has no more authority than a mechanical one. It returns
// artifact CANDIDATES `{id, kind, text}` — there are no provenance, verdict,
// or status fields for a model to forge, its contract copy is frozen, and V
// still checks required evidence positionally. A provider that answers with
// nonsense produces a failing episode, never a false accepted one.
package cellcog

import (
	"context"
	"errors"
)

// Provider is the narrow port this package needs from a cognition backend
// (eng/go §2.3 — the interface is defined by its consumer, not its
// implementer). Implementations are adapters: they own their effects
// (subprocess, network) and nothing else in this package performs any.
type Provider interface {
	// Name identifies the provider. It is disclosed in the closure through
	// the resolved spec's parameters, so a reader can always tell what held
	// the seat.
	Name() string

	// Complete answers a rendered prompt with the provider's raw text.
	// Implementations must honor ctx and bound the response.
	Complete(ctx context.Context, prompt string) (string, error)
}

// ErrNoProvider is returned when a rented seat is built without a provider.
// It fails closed at the seat, before any prompt is rendered.
var ErrNoProvider = errors.New("cellcog: no provider")
