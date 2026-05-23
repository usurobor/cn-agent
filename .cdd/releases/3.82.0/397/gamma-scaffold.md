# γ scaffold — cycle/397

**Issue:** cnos#397 — Phase 4a of #366 — extract `delta/SKILL.md` from `operator/SKILL.md` (δ-as-role; two-sided membrane)

**Mode:** design-and-build; γ+α+β-collapsed-on-δ per breadth-2026-05-12 precedent. δ is the lone agent (γ=δ collapsed; α/β reviews collapsed). β-α-collapse acknowledged.

**Identity:** `delta@cdd.cnos` (cnos elision form per `operator/SKILL.md` §"Git identity for role actors")

## Surfaces

1. **NEW** `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` — the δ-role skill carrying the two-sided membrane (outward §3 + inward §3a) plus override semantics.
2. **EDIT** `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` — extract §3a entirely, refactor §3 to keep mechanics (release script, CI polling) but move authority/policy to delta/, extract §4 (override) to delta/. Replace extracted blocks with one-line cross-references.
3. **EDIT** `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` — update two refs: §2.5 §3b template ("δ enrichment: `operator/SKILL.md` §3a") → `delta/SKILL.md`.
4. **EDIT** `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — update ref at §3.6 ("δ ratifies and enriches" pointer to `operator/SKILL.md` §3a) → `delta/SKILL.md`.
5. **EDIT** `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` — update ref at Rule 7 ("δ owns the enrichment half") → `delta/SKILL.md`.

**Refs that STAY pointing at `operator/SKILL.md`** (operator-as-coordinator / harness / release-mechanics surfaces — Phase 4b/4c claim later):
- `post-release/SKILL.md` §3.4 reference (release mechanics — Phase 4c)
- `release/SKILL.md` two refs to §3.4 (release mechanics — Phase 4c)
- `activation/SKILL.md` several refs to §5.2/§5.3/§3.4 (dispatch configs + release mechanics)
- `review/SKILL.md` §3.11b wave-mode ref (§5.2 wave-mode = dispatch config)
- `CDD.md` refs to git identity, timeout-recovery, §3.3, §3.4 (operator/coordinator surfaces)
- `COHERENCE-CELL.md` doctrinal refs (predicts the split; references stay)
- `RECEIPT-VALIDATION.md` Phase 4 target reference (historical)
- `COHERENCE-CELL-NORMAL-FORM.md` Phase 4 target reference (historical)

## AC oracle approach (mechanical)

| AC | Oracle |
|---|---|
| AC1 | `test -f src/packages/cnos.cdd/skills/cdd/delta/SKILL.md && head -20 of file contains frontmatter with name: delta, artifact_class: skill, governing_question, parent: cdd, triggers, scope: role-local` |
| AC2 | `rg "outward|inward" delta/SKILL.md` returns explicit two-sided membrane sections; each side cross-references the other and the interacting surfaces |
| AC3 | `rg "δ-inward-membrane\|inward membrane" src/packages/cnos.cdd/skills/cdd/` returns substance hits ONLY in `delta/SKILL.md`; `operator/SKILL.md` §3a is replaced by a one-line redirect |
| AC4 | `rg "override.*degraded\|never rewrites" delta/SKILL.md` returns override-as-degraded-boundary-action with the verdict-vs-decision distinction (cnos#367 freeze) |
| AC5 | `rg -l "operator/SKILL.md.*§3a\|operator/SKILL.md.*inward membrane" cdd/` returns zero hits; gamma/alpha/beta point at `delta/SKILL.md` for δ-role content |
| AC6 | `rg "operator-as-coordinator\|dispatch coordination\|γ=δ collapse" operator/SKILL.md` still returns content; algorithm §1 + §2 Wait + §5 Dispatch + §6/§7/§8/§9/§10 still present |
| AC7 | `rg "scripts/release.sh\|gh run list" delta/SKILL.md` returns ZERO; `rg "tag\|release\|deploy" delta/SKILL.md` returns hits only in authority-naming context (not commands/mechanics) |

## Implementation-contract conformance (per cnos#393 Rule 7)

All 7 axes pinned by issue:
- Language: Markdown — confirmed (skill files are .md)
- CLI integration target: N/A — confirmed (no CLI added)
- Package scoping: `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` + edits to `operator/SKILL.md` — confirmed
- Existing-binary disposition: `operator/SKILL.md` keeps the operator-as-coordinator + dispatch role; δ-role surface moves to `delta/SKILL.md`; harness → Phase 4b; release effector → Phase 4c — confirmed
- Runtime dependencies: None — confirmed
- JSON/wire contract: N/A — confirmed
- Backward compat: All operator/SKILL.md references that addressed δ-role content updated to point to delta/; operator/SKILL.md itself preserved with operator-coordinator content intact — plan confirmed

## Refusal-condition check (negative — none triggered)

- No ambiguous-surface sections forced to improvise: the membrane-policy (outward §3 authority-policy + inward §3a) is unambiguously δ-role; release mechanics (`scripts/release.sh`, branch cleanup runbook) are unambiguously Phase 4c (per AC7); dispatch configurations §5 + timeout recovery §8 + kata §9 + wave §10 are unambiguously operator-as-coordinator/harness (Phase 4b later).
- No parallel cycle on operator/SKILL.md (last edit landed in #393; #398/#399 not yet started).
- No cross-reference update requires content change — all are pure path swaps (operator/SKILL.md §3a → delta/SKILL.md §inward; operator/SKILL.md §3.4 stays where it is).
