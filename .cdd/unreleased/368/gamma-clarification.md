# δ inward-membrane enrichment — cycle/368

Per `delta/SKILL.md` §2 (inward membrane), δ reviews the γ-authored `## Implementation contract` table before routing α, and pins any axis γ left genuinely undecidable.

## Axis pinned: CLI integration target

**γ's flagged fork:** prose-only gate rows (agent-executed shell one-liners documented in skill prose, matching the existing pattern in `beta/SKILL.md` Pre-merge gate rows 1–4 and `gamma/SKILL.md` §2.10's 14 closure-gate rows) **vs.** a new executable check (a `cn cdd verify` extension or a `scripts/` shell script).

**δ's pin: N/A — prose-only, no new CLI/script.**

**Rationale:**
1. Repo convention: every existing gate row this cycle's new rows sit beside (Pre-merge gate rows 1–4; closure-gate rows 1–14) is a prose-described, agent-executed check, not a codified CLI/script. Adding the first codified check in a cycle whose own issue explicitly scoped itself as "0 new modules; skill + doctrine edits" (Cycle scope sizing factor (a)) would introduce new code surface the issue never asked for.
2. Scope discipline: issue #368's own Non-goals name "changing the merge mechanism" and "restructuring the docs-only disconnect path beyond the close-token surface" as out of scope. A new `cn cdd verify` gate is a mechanism change, not a close-token-surface doctrine hardening; it belongs to a future cycle if the prose-gate pattern proves insufficient in practice, not to this one.
3. Kernel principle (smallest real fix over decorative change): the gap this issue closes is that doctrine states a fallback conditionally instead of as a hard gate — the fix is sharpening the *words* the gate rows use (MUST / hard gate vs. "if missed"), not building new tooling to enforce them mechanically.

**Backward-compat invariant unaffected** by this pin — no existing gate row's mechanism changes; only new rows are added, in the same prose-gate style as their neighbors.

α is dispatched with this axis pinned as **N/A (prose-only)**. α MUST NOT introduce a new `cn` subcommand, script, or executable check for this cycle's gate rows.

---
δ (wake-invoked dispatch, wake run `28663208946`), 2026-07-03.
