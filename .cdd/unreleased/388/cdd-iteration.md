# cdd-iteration — Cycle #388 (Phase 2.5 of #366)

**Cycle:** #388 — split generic schemas/cdd/ from CDS/CDR domain evidence
**Merge:** `e8958788` on `main`
**Closed:** 2026-05-21 (auto-closed via `Closes #388` in merge commit)
**Mode:** design-and-build (γ+α+β-collapsed-on-δ per `breadth-2026-05-12` wave manifest precedent; β-α-collapse acknowledged in `beta-review.md`)
**Rounds:** R1 APPROVE (no fix-round)
**ACs:** 7/7 PASS

## §1 Findings dispositioned

### F1: γ-skill `protocol_id` automation debt

- **Source:** cycle 388 β-review §Findings, observation 6
- **Class:** `cdd-skill-gap`
- **Trigger:** ε process-gap check during cycle close-out
- **Description:** Phase 2.5 introduces `protocol_id` as a required field in `#Receipt` (e.g. `cnos.cdd.cds.receipt.v1`, `cnos.cdd.cdr.receipt.v1`). Future cycles authoring receipts must declare `protocol_id` correctly. γ's current scaffold workflow does not auto-suggest the correct `protocol_id` value based on the cycle's protocol class (CDS/CDR/etc.). γ scaffolding could either prompt for the value at cycle creation or detect it from issue labels / parent class. Without automation, individual cycles risk omitting or mis-declaring `protocol_id`, which would silently break V's dispatch logic in Phase 3.
- **Root cause:** `gamma/SKILL.md` §2.5 scaffold authoring does not yet reference the new `protocol_id` dispatch convention. The convention is fresh — landed at `e8958788` — so the gap is by construction, not negligence.
- **Disposition:** `next-MCA`. The patch shape is "add a row to `gamma/SKILL.md` §2.5 Step 3a or 3b that names the `protocol_id` declaration responsibility at scaffold time, with values pinned per protocol class." Best landed alongside Phase 3 (when V starts consuming the field) or as a small standalone γ-skill patch.
- **Issue filed:** none yet. Suggest filing a small `gamma/SKILL.md §2.5 protocol_id automation` issue when Phase 3 dispatches, or folding into Phase 3's own scope as an "auxiliary skill patch."
- **First AC for the eventual MCA:** `gamma/SKILL.md` §2.5 names `protocol_id` as a required scaffold field; cites the canonical values per protocol class; references `schemas/{cds,cdr}/receipt.cue` as the source-of-truth pins.

### F2: Essay v0.1.1 update deferred

- **Source:** cycle 388 β-review §Findings, observation 5 + gamma-closeout §Debt-Forwarded point 3
- **Class:** process debt (not strictly a `cdd-*-gap`, but documentation-coherence debt)
- **Trigger:** Phase 2.5 closes essay open question #1 (option (a) split vs (b) adapter) with (a) decided. The essay's open-questions section should reflect the resolution.
- **Description:** `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` v0.1.0 (commit `7cdcd8d4`) lists open question #1 as "Should Phase 2.5 prefer split schemas (...) or a generic schema with domain adapters?" That question is now resolved by #388's option (a) decision. The essay should get a v0.1.1 update folding the resolution into its open-questions section (or moving the resolved question into a "Decisions" subsection).
- **Disposition:** `next-MCA`. Deferred per cycle 388's scope discipline — folding the essay update into #388 would expand scope. Small follow-on doc cycle.
- **Issue filed:** none yet. Suggest a small docs cycle, or fold into any future essay update.
- **First AC:** essay §"Open questions" item #1 marked resolved with link to cnos#388 and `schemas/cdd/README.md §"Architectural choice"`.

### F3: `evidence_refs` value-type union widening (build-time discovery)

- **Source:** cycle 388 β-review §Findings, observation 3
- **Class:** none (within-cycle engineering refinement; not a protocol gap)
- **Description:** The original design proposed `evidence_refs` as `[string]: #EvidenceRef` (typed open map). During CUE authoring, the value type was widened to a union to accommodate both single-ref and list-of-refs evidence shapes (CDS needs single refs like `self_coherence: #EvidenceRef`; CDR needs lists like `claim: [...#EvidenceRef]`). The union resolved the inheritance cleanly.
- **Disposition:** `drop` (positive within-cycle refinement; no protocol gap). Recorded for posterity.

## §2 No-findings observations (informational)

- AC2 false-positive on `contract.cue` comment caught + fixed pre-commit. β-mechanical grep gate worked as expected; non-event.
- `cue.mod/module.cue` added at repo root as part of the CUE module declaration. Engineering-level decision, not a protocol concern.
- CDR list min-length-1 constraint applied during AC4 implementation. Design refinement, not a gap.

## §3 Trigger assessment (per `gamma/SKILL.md` §2.8 table)

| Trigger | Fire condition | Fired? | ε note |
|---|---|---|---|
| Review churn | review rounds > 2 | **No** | R1 APPROVE on first pass. |
| Mechanical overload | mechanical ratio > 20% AND findings ≥ 10 | **No** | 6 β findings total; mechanical class only catches AC2 false-positive (caught pre-commit). |
| Avoidable tooling / environment failure | environment blocked the cycle | **No** | `cue vet` worked cleanly; the only harness pattern (origin push --delete 403) was non-blocking and previously named in wave-2026-05-19 ε iteration F5. |
| Loaded-skill miss | a loaded skill should have prevented a finding | **Borderline** | F1 (`protocol_id` automation) is a γ-skill-miss in the sense that γ doesn't have the convention codified yet. But the convention is fresh from this cycle, so the miss is constructive (the cycle creates the gap by introducing the field). Trigger does not fire; F1 records the discovery as `next-MCA`. |

## §4 INDEX update

Add to `.cdd/iterations/INDEX.md`:

```
| 388 | #388 | 2026-05-21 | 3 | 0 | 2 | 0 | .cdd/unreleased/388/cdd-iteration.md |
```

Findings: 3 (F1, F2, F3). Patches: 0 immediate (F3 was within-cycle, dropped). MCAs: 2 (F1, F2). No-patch: 0. Path: `.cdd/unreleased/388/cdd-iteration.md` (will move to `.cdd/releases/<version>/388/cdd-iteration.md` at next release per `release/SKILL.md` §2.5a).

## §5 Skill-gap candidate disposition

F1 is the only `cdd-skill-gap` candidate from this cycle. Per ε's MCA discipline (`epsilon/SKILL.md §1`), the patch shape is clear but better landed when Phase 3 dispatches (since Phase 3 is the consumer of `protocol_id`). Carry as `next-MCA` named in the Phase 3 cycle's "auxiliary skill patches" scope, or as a small standalone γ-skill patch before Phase 3 dispatches.

F2 is documentation debt; not a `cdd-*-gap` per the strict vocabulary but worth recording so the essay update is not forgotten.

## §6 Deferred outputs

- **Cycle-dir move.** `.cdd/unreleased/388/` → `.cdd/releases/<version>/388/` at next release per `release/SKILL.md` §2.5a. Not blocking; standard release-time mechanic.
- **Essay v0.1.1 update** (F2). Deferred; folded into next docs cycle.
- **γ-skill `protocol_id` patch** (F1). Deferred; folded into Phase 3 or standalone γ-skill cycle.

## §7 Next-MCA commitment

Phase 3 dispatches next (gated on Phase 2.5; gate now satisfied). cnos#376 Sub 1 can dispatch in parallel since it inherits #388's (a) architectural decision. Both should reference `schemas/cdd/README.md §"Architectural choice"` as the canonical decision-record.

Filed by ε on 2026-05-21.
