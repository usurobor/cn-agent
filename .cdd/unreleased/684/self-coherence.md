## Gap

**Issue:** cnos#684 — "Sigma activation channel — rename `.cn-sigma/logs`, move to symmetric append-only orphan refs (independent stream, off main HEAD)."

**Version / mode:** Mode `explore` (issue frontmatter: `Mode: explore`, `Status: design-first (explore)`). This is a design-first cell, not an implementation cell for orphan refs or history rewriting. The deliverable is a design document + mechanism spec + dry-run migration plan, produced entirely inside the `cnos` repo on `cycle/684`.

**What is not-current about the pre-684 state:** `.cn-sigma/logs/` — dated `YYYYMMDD.md` files under `main`'s tree — is the Sigma foreign→home activation channel, but (a) the name "logs" is too generic for what is specifically an activation-stream, and (b) the channel is an independent communication timeline, not product state, yet it materializes in `main`'s HEAD tree (~1024 commits / ~10 MB, the dominant source of `main` HEAD churn per the operator's dispatch-authorization comment). The v0 convention (`AGENT-ACTIVATION-LOG-v0.md`) also assumes a bare-SHA cursor and a commit-to-`main` attach contract, both of which stop working once the stream moves off `main`.

**What this cycle closes:** the go-forward mechanism design (rename, orphan-ref names + writer/reader roles, amended registration schema, attach-contract sequence, orphan-ref invariants + enforcement, promotion boundary) plus a dry-run-only migration plan and a specified (not executed) content-preservation verification procedure for the already-committed `.cn-sigma/logs/**` history on `main`. It explicitly does **not** close: any physical write to `.cn-sigma/**`, any history rewrite, any ref creation/deletion, any force-push. Those remain operator-executed, separately-gated actions per the operator's binding dispatch-authorization comment (cnos#684, @2026-08-01T07:50:17Z).

**Dispatch chain:** operator authorization (2026-08-01T07:50:17Z, "κ dispatch authorization + bounded scope refinement") → γ R0 scaffold (`.cdd/unreleased/684/gamma-scaffold.md`) → δ inward-membrane enrichment confirming γ's two scope-resolution calls (`.cdd/unreleased/684/gamma-clarification.md`) → this α R0 cycle.

## Skills

**Tier 1 (loaded):**
- `CDD.md` — canonical lifecycle and role contract.
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — this role's execution detail; §2.1 dispatch intake, §2.5 self-coherence authoring (incremental, one section per commit — followed in this file), §2.6 pre-review gate, §2.7 request review.
- `issue/SKILL.md` — not separately loaded as a file; AC-boundary interpretation was already resolved by γ (scaffold) and δ (clarification) before dispatch, so α consumed their resolution rather than re-deriving AC boundaries from the raw issue.
- `design/SKILL.md`, `plan/SKILL.md` — **not required, explicit reason:** this cycle's design artifact *is* the deliverable (the convention doc itself), and γ's scaffold already carries the plan-equivalent (AC oracle table, source-of-truth table, α prompt with 4 concrete files, scope guardrails). A separate design/plan artifact preceding the convention doc would duplicate the scaffold rather than add sequencing value for a 4-file, non-code, no-implementation-sequencing cycle.

**Tier 2 (eng/* always-applicable):** not loaded as a distinct bundle — this cycle produces no code, no tests, no CI-affecting artifacts (Markdown + one unexecuted shell script). The shell script (`verify-channel-reconstruction.sh`) was authored with ordinary shell-correctness discipline (`set -euo pipefail`, quoted variables, `mktemp -d` for scratch state) but no `eng/shell` Tier 3 bundle exists in this repo's Tier 2/3 taxonomy to load explicitly; general shell hygiene was applied inline.

**Tier 3 (issue-specific, per γ's scaffold):**
- `delta/SKILL.md` §9.12 (cell/substrate identity boundary, cnos#626) — the binding doctrine constraint that shaped AC6's scope-down (verification procedure specified, not executed) and that forbids any read/write of `.cn-sigma/**` by this cell regardless of checkout visibility. Read in full per the dispatch prompt's explicit instruction.
- `docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md` — read in full as the doc being superseded; source-of-truth for what v1 must remain consistent with (not contradict/duplicate silently) per the dispatch prompt.
- `gamma-scaffold.md` and `gamma-clarification.md` — both read in full before drafting, per the dispatch prompt's explicit reading list; treated as settled (not re-litigated) per the dispatch instruction.

No skill named in the dispatch prompt or discoverable via `## Related artifacts`-style enumeration was left unread. The issue itself named no `## Related artifacts` section in the strict alpha/SKILL.md §2.1.5 sense (its `## Related` section lists other issues, not artifact paths); all named artifact-path items in `## Related` (`AGENT-ACTIVATION-LOG-v0.md`) were read.
