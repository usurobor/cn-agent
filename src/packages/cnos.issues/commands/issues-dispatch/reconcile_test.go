package issuesdispatch

import "testing"

// contradictoryBody reconstructs the exact #614/#633 shape named in
// gamma-scaffold.md §1 AC1: a design-first issue body ending in the
// bold-markdown hold phrase, preceded by a "---" separator (the live
// fixtures' common shape).
const contradictoryBody = "## Design-first change\n\nSome description of the change.\n\n**Acceptance criteria**\n- AC1: does the thing\n\n---\n\n**Not dispatched** — `status:ready`. Design-first; dispatch on explicit operator authorization."

const contradictoryBodyWant = "## Design-first change\n\nSome description of the change.\n\n**Acceptance criteria**\n- AC1: does the thing"

// plainHoldBody is the non-bold, no-separator variant -- proves the
// pattern is not over-fit to the exact markdown decoration.
const plainHoldBody = "Some issue text.\n\nNot dispatched — status:ready. Design-first; dispatch on explicit operator authorization."

const plainHoldBodyWant = "Some issue text."

func TestReconcileBody_StripsHoldPhrase_BoldWithSeparator(t *testing.T) {
	got, changed := ReconcileBody(contradictoryBody)
	if !changed {
		t.Fatalf("expected changed=true for a body carrying the hold phrase")
	}
	if got != contradictoryBodyWant {
		t.Fatalf("cleaned body = %q, want %q", got, contradictoryBodyWant)
	}
}

func TestReconcileBody_StripsHoldPhrase_PlainNoSeparator(t *testing.T) {
	got, changed := ReconcileBody(plainHoldBody)
	if !changed {
		t.Fatalf("expected changed=true for a body carrying the hold phrase")
	}
	if got != plainHoldBodyWant {
		t.Fatalf("cleaned body = %q, want %q", got, plainHoldBodyWant)
	}
}

// TestReconcileBody_CleanBodyIsSafeNoOp is AC1's negative oracle: a body
// with no hold phrase must come back byte-identical, changed=false --
// ReconcileBody must never rewrite a body that has nothing to clean.
func TestReconcileBody_CleanBodyIsSafeNoOp(t *testing.T) {
	clean := "## Design-first change\n\nSome description.\n\n**Acceptance criteria**\n- AC1: does the thing"
	got, changed := ReconcileBody(clean)
	if changed {
		t.Fatalf("expected changed=false for an already-clean body, got changed=true, body=%q", got)
	}
	if got != clean {
		t.Fatalf("clean body must be returned byte-identical, got %q, want %q", got, clean)
	}
}

// TestReconcileBody_DoesNotMatchDifferentPhrasingFamily proves the
// pattern does not over-match #618's different hold convention ("Filing-
// only — operator holds dispatch."), which gamma-scaffold.md §5 pins as
// out of scope for this cell -- a future occurrence of a different
// phrasing family is a new pattern to add explicitly, not something this
// primitive should guess at.
func TestReconcileBody_DoesNotMatchDifferentPhrasingFamily(t *testing.T) {
	body := "Some issue text.\n\nFiling-only — operator holds dispatch."
	got, changed := ReconcileBody(body)
	if changed {
		t.Fatalf("expected changed=false for a different hold-phrase family (#618's phrasing), got changed=true, body=%q", got)
	}
	if got != body {
		t.Fatalf("body must be returned byte-identical, got %q, want %q", got, body)
	}
}
