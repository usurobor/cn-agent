# self-coherence.md — cnos#639 (α, R0)

## Gap

**Issue:** [cnos#639](https://github.com/usurobor/cnos/issues/639) — "cds/doctrine: clarify the dispatch `run_class` taxonomy (first-pass / repair / recovery / continuation / scope-continuation)."

**Mode:** design-and-build (doctrine/design only — the issue's own header states this explicitly). **Cell kind:** `doctrine` (pinned by γ; no reclassification escape hatch — see γ's scaffold Friction note 1).

**Dispatch authority:** κ's 2026-07-09 operator-directive comment grants explicit A2/CAP authority — "Fix stale prose, PR metadata, and link/doc drift inside scope without asking. Escalate only if the doctrine change requires behavior/FSM/code changes." This cell's whole deliverable is doctrine prose, so the implementation below sits inside that grant; no escalation was triggered (the taxonomy reconciliation did not require any FSM/code/label change to ship as scoped).

**The gap.** The `run_class` taxonomy that classifies each dispatch-wake firing (first pass / repair re-entry / resumed-from-matter / etc.) was enumerated **inconsistently** across three doctrine surfaces:
- `dispatch-protocol/SKILL.md` §2.8 listed 4 values (`first_pass, repair_pass, manual_delta_repair, blocked`) — missing `resumed_from_matter` entirely.
- `cds-dispatch/SKILL.md` §"Repair re-entry preflight" listed 5 values (the most complete, but still missing the operator-authorized continuation shape).
- `delta/SKILL.md` narrated individual resume shapes (§9.10, §9.11) without a canonical enum.

On top of the divergence, the taxonomy was **incomplete**: cnos#626 hit an operator-authorized "scope continuation" shape three times across R0/R1/R2 (a new bounded AC scope dispatched on an issue whose prior round had already converged and merged) that matched none of the existing classes, and self-disclosed the gap explicitly each time (`.cdd/unreleased/626/self-coherence.md` §"run_class note", twice; `gamma-closeout.md`, twice) rather than silently mislabeling.

**γ's scaffold** (`.cdd/unreleased/639/gamma-scaffold.md`) is the primary spec for this round: it names the per-AC oracle list (§1), the source-of-truth table (§2), the α/β prompts (§3/§4), the scope guardrails (§5), and pre-identified friction (§6). This file follows that scaffold precisely; nothing below re-derives what γ already pinned.

**Scope guardrails honored (restated from γ's scaffold §5, verified against the actual diff in §ACs below):** no `.go` file changes; no `transitions.json` diff (byte-identical); no dispatch-behavior change; no new status labels; no `cell_kind` enforcement; no Demo 0; no changes to #626/sparse-checkout/write-fence content; `#642` was not dispatched; `scripts/ci/*.sh` was not edited; `cnos-cds-dispatch.golden.yml` and `.github/workflows/cnos-cds-dispatch.yml` were regenerated mechanically via `cn-install-wake`, never hand-edited.
