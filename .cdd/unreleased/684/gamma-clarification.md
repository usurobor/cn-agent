# δ inward-membrane enrichment — cnos#684

γ's R0 scaffold flagged two scope ambiguities (gamma-scaffold.md §"Scope ambiguities flagged for δ"). δ reviews and confirms both rather than escalating to operator — neither is genuinely undecidable; both are already resolved by the operator's own binding comment text.

## 1. AC5 split (design-landed vs. physical-strip)

**Confirmed.** The operator's binding comment ("κ dispatch authorization + bounded scope refinement") states explicitly: "the mechanical worker produces PRs into `main`; it does not and must not rewrite history, force-push `main`, or delete refs" and separately: "Deliverables are design + mechanism + dry-run plan only. No history rewrite, no `main` force-push, no ref deletion, no `git rm` of live paths on `main`." AC5's physical component ("`.cn-sigma` HEAD retains only current cursors/state ... not the dated stream") requires exactly the `git rm`-of-live-paths action the operator forbids the worker from performing. There is no reading of the operator's text under which the worker could complete AC5 physically without violating the guardrail in the same comment. γ's split (design-landed this cycle; physical-strip deferred to the operator-executed AC7 step) is the only scope consistent with both AC5 and the guardrail as pinned. Confirmed, not overridden.

## 2. AC6 verification scope-down

**Confirmed.** `delta/SKILL.md` §9.12 (cnos#626) is standing doctrine, not a per-cycle discretionary call — cell-execution cognition (δ/γ/α/β) structurally does not read or write `.cn-{agent}/logs/` regardless of checkout visibility. Overriding it would require an explicit δ override under §3 (degraded boundary action, populated override block) against a doctrine constraint, not an implementation-contract judgment call — out of scope for routine inward-membrane enrichment. γ's scope-down (specified verification procedure, not executed evidence) is the correct application of existing doctrine, not a new decision. Confirmed.

## Disposition

No scaffold amendment needed. α proceeds under gamma-scaffold.md as authored, with this clarification as the citable record that δ reviewed and confirmed both interpretive calls before routing to α.
