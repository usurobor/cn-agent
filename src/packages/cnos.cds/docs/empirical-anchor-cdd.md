# Empirical anchor for CDS — cnos cycles as CDS-realization evidence

**Version:** 0.1 (Sub 7 of [cnos#403](https://github.com/usurobor/cnos/issues/403))
**Date:** 2026-05-22
**Placement:** `src/packages/cnos.cds/docs/empirical-anchor-cdd.md`
**Audience:** Reviewers verifying that cnos.cds v0.1's surface family (defined in [`skills/cds/CDS.md`](../skills/cds/CDS.md)) is sufficient to retroactively describe cnos's own engineering practice; future software projects binding their own `.cds/` artifacts to cnos.cds.
**Scope:** A surface-by-surface mapping from representative cnos cycles onto the CDS canonical surfaces those cycles exercised. Demonstrates that the CDS surface family hosts cnos's own software-cycle practice without contradiction. Closes [cnos#412](https://github.com/usurobor/cnos/issues/412) AC1–AC8.

> This document is a **mapping table + commentary**, not a re-statement of
> CDS.md. It does not embed cycle bodies; citations carry the cycle number
> (linked to its GitHub issue) and the close-out artifact path or merge SHA.
> Per [`CDS.md §"Empirical anchor" → "Sub 7 deferred surface-by-surface
> mapping"`](../skills/cds/CDS.md): Sub 2 (cycle/407) ships the
> shape-compatibility claim; Sub 7 (this document) ships the
> surface-by-surface verification.

---

## How to read this doc

Pick a CDS surface — for example, **§Gate F1–F10** in [`CDS.md`](../skills/cds/CDS.md) — and read the corresponding sub-section of [§Cycle-to-surface mapping](#cycle-to-surface-mapping) below. The mapping table for that surface lists representative cnos cycles where that surface was actually exercised, with a one-line "What it exercised" column and a path or SHA citation under "Artifact citation."

A cycle may appear in multiple tables: most cycles exercise more than one CDS surface (because the surfaces are deliberately compositional). The aggregator file [`.cdd/iterations/INDEX.md`](../../../../.cdd/iterations/INDEX.md) is the canonical receipt stream — every row there is a CDS receipt instance.

This doc is **representative, not exhaustive**. The cycle wave is open-ended; new cycles continue to exercise these surfaces. The mapping below is calibrated to make the shape-compatibility claim concrete: every CDS surface has at least two cycles that exercised it, and every cited cycle has at least one CDS surface it visibly populates. The doc does not attempt to enumerate every cycle's full surface footprint; the per-cycle close-out artifact set (`.cdd/unreleased/<N>/` or `.cdd/releases/.../<N>/`) is the canonical per-cycle record.

### Anchor pin

This document anchors against cnos `main` at commit
**`71b25672`** (the head of cnos `main` at the time cycle/412 was filed —
the merge of cycle/410, Sub 5 of cnos#403). Citations are written
`<path>@71b25672` or `#<N>@<short-sha>` where a specific merge commit is
load-bearing.

### CDS surfaces enumerated

The CDS top-level surfaces (per [`CDS.md`](../skills/cds/CDS.md), in document order) are:

1. `§Purpose`
2. `§Architecture choice`
3. `§Persona, Protocol, Project`
4. `§Six-field instantiation contract` (Fields 1–6)
5. `§Selection function`
6. `§Development lifecycle`
7. `§Coordination surfaces`
8. `§Artifact contract`
9. `§CDS Trace`
10. `§Mechanical vs judgment`
11. `§Review CLP`
12. `§Gate` (F1–F10)
13. `§Assessment`
14. `§Closure`
15. `§Retro-packaging`
16. `§Large-file authoring rule`
17. `§Empirical anchor` (this anchor's home in the protocol)
18. `§Related documents`
19. `§Non-goals`

The mapping below covers the load-bearing operational surfaces (4–16). The doctrinal surfaces (1–3, 17–19) are exercised at the document layer and need no separate cycle citation.

---

## Cycle-to-surface mapping

### §Six-field instantiation contract

CDS's six-field contract (Matter type, Review oracle, γ close-out artifact, δ cadence, ε iteration cadence, Actor collapse rule) is the spine declared by [`ROLES.md §3`](../../../../ROLES.md) and instantiated for software by [`CDS.md §Six-field instantiation contract`](../skills/cds/CDS.md). Every cnos cycle since #364 has been a six-field instantiation; the cycles below are load-bearing because they introduced, refined, or stress-tested the contract itself.

| Cycle | Issue | Date | What it exercised | Artifact citation |
|---|---|---|---|---|
| #364 | Coherence-cell doctrine | 2026-05-15 | Established the recursive cell model that the six-field instantiation contract crystallized; merge `32b126e4` | `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md@71b25672` |
| #367 | Validation-surface design | 2026-05-15 | Designed the V predicate that became Field 2's "compilation + type-check + build" + "AC verification" oracle; γ close-out at merge | `.cdd/releases/docs/2026-05-17/367/`@71b25672 (cycle dir post-disconnect) |
| #388 | Phase 2.5 schemas split | 2026-05-21 | Separated generic `schemas/cdd/` from domain `schemas/{cds,cdr}/` — the architectural choice that lets Field 3's typed receipt distinguish protocol-overlay from kernel | `.cdd/unreleased/388/cdd-iteration.md`@71b25672 |
| #390 | cnos.cdr v0.1 doctrinal contract (Sub 1 of #376) | 2026-05-21 | First six-field instantiation for research realization (CDR); structural precedent that CDS Sub 2 (#407) mirrored | `src/packages/cnos.cdr/skills/cdr/CDR.md@71b25672` |
| #407 | cnos.cds CDS.md authored (Sub 2 of #403) | 2026-05-22 | Second six-field instantiation, for software realization; declared the contract this doc anchors against | `src/packages/cnos.cds/skills/cds/CDS.md@71b25672` (1043 lines at file's authoring) |
| #402 | CCNF spine rewrite (Phase 7 of #366) | 2026-05-21 | Quarantined the software-specific content under "Software-specific realization — pending cds extraction" — the act that scoped the CDS field-instantiation as separable from CDD's kernel | `.cdd/unreleased/402/gamma-closeout.md`@71b25672 |

### §Selection function

CDS's selection function (per [`CDS.md §Selection function`](../skills/cds/CDS.md)) decides which gap a cycle addresses next: P0 override → operational-infrastructure override → MCI freeze → weakest-axis → maximum-leverage → dependency order → effort-adjusted tie-break. The cycles below exercised individual rungs of that ladder.

| Cycle | Issue | Date | What it exercised | Artifact citation |
|---|---|---|---|---|
| #364 | Coherence-cell doctrine | 2026-05-15 | Originated the gap-selection doctrine implicitly: doctrine surface was the weakest axis at the time | `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md@71b25672` |
| #366 | Coherence-cell executability roadmap | (in-flight) | The P0-override-vs-MCA-default tension: roadmap dispatch chose a multi-phase wave over single-cycle MCAs | `.cdd/releases/docs/.../<366 phases>/cdd-iteration.md` (per-phase) |
| #379 | cn-sigma proposal intake | 2026-05-19 | Stale-backlog re-evaluation rung; proposal absorbed under selection function's intake path | `.cdd/releases/3.78.0/379/cdd-iteration.md`@71b25672 |
| #393 | δ-as-architect protocol patches | 2026-05-21 | Maximum-leverage rule: the implementation-contract Rule 7 patch returned leverage on every subsequent cycle | `.cdd/unreleased/393/cdd-iteration.md`@71b25672 (4 patches, 0 MCAs) |
| #401 | Phase 6 ε upscope | 2026-05-21 | Dependency-order rung: ε's surface had to land before Phase 7's CDD.md rewrite could close | `.cdd/unreleased/401/cdd-iteration.md`@71b25672 |
| #403 | cnos.cds bootstrap tracker | (in-flight) | Maximum-leverage rule applied to a wave: extract-by-reference v0.1 over deep per-role rewrite | issue body `cnos#403` |

### §Development lifecycle

CDS's development lifecycle (per [`CDS.md §Development lifecycle`](../skills/cds/CDS.md)) is the 10-phase step table + state machine + branch rule + branch-pre-flight + skill-loading-tiers that every cycle traverses. Every recent cycle exercises this surface end-to-end; the cycles below introduced or stress-tested specific phases.

| Cycle | Issue | Date | What it exercised | Artifact citation |
|---|---|---|---|---|
| #366 | Coherence-cell executability roadmap | (in-flight) | Established the multi-phase lifecycle pattern that the 10-phase step table later codified | (issue body; phase-by-phase merges) |
| #367 | Validation-surface design | 2026-05-15 | Phase 2 (design) and Phase 3 (V) precursor: designed the validation surface CDS Field 2 absorbs | `.cdd/releases/docs/2026-05-17/367/` |
| #369 | Phase 2 receipt schemas | 2026-05-17 | Phase 2 deliverable: receipt schemas that the typed γ artifact (Field 3) extends; 3 cdd-protocol-gap findings recorded | `.cdd/releases/docs/2026-05-17/369/cdd-iteration.md`@71b25672 |
| #370 | CCNF normal-form doctrine | 2026-05-17 | Doctrine-class lifecycle: docs-only disconnect via γ-as-δ; 2 findings + 1 patch + 1 MCA recorded | `.cdd/releases/docs/2026-05-17/370/cdd-iteration.md`@71b25672 |
| #392 | Phase 3 V ported to Go (`cn cdd-verify`) | 2026-05-21 | Phase 3 deliverable; γ-clarification rescue mid-flight; 4 findings, 4 MCAs filed | `.cdd/unreleased/392/cdd-iteration.md`@71b25672 |
| #397 | Phase 4 δ split (planning) | 2026-05-21 | Phase 4 planning cycle: split δ-as-architect from δ-as-effector to scope subsequent waves | `.cdd/unreleased/397/` |
| #398 | Phase 4 harness/release-effector | 2026-05-21 | Phase 4 deliverable; 2 findings, 2 MCAs recorded; the release-effector implementation cycle | `.cdd/unreleased/398/cdd-iteration.md`@71b25672 |
| #400 | Phase 5 γ shrink | 2026-05-21 | Phase 5 deliverable: γ surface narrowed; 2 patches + 1 MCA recorded | `.cdd/unreleased/400/cdd-iteration.md`@71b25672 |
| #401 | Phase 6 ε upscope | 2026-05-21 | Phase 6 deliverable: ε's iteration-cadence surface broadened to cover receipt-stream review | `.cdd/unreleased/401/` |
| #402 | Phase 7 CDD.md rewrite | 2026-05-21 | Phase 7 deliverable: CCNF spine rewrite; 2 findings (1 MCA → #403, 1 no-patch) | `.cdd/unreleased/402/cdd-iteration.md`@71b25672 |

### §Coordination surfaces

CDS's coordination surfaces (per [`CDS.md §Coordination surfaces`](../skills/cds/CDS.md)) include cycle-state evidence (close-out artifacts as observable state), polling primitives, mid-flight clarification, and cross-repo proposals. The cycles below exercised each of those mechanisms.

| Cycle | Issue | Date | What it exercised | Artifact citation |
|---|---|---|---|---|
| #376 | cnos.cdr v0.1 wave (parent tracker) | 2026-05-21 | Cross-repo proposal intake: cph bundle uplifted into cnos.cdr; mid-flight wave with five subs (#390/#394/#395/#396 + #389 dependency) | issue body `cnos#376` |
| #379 | cn-sigma proposal intake | 2026-05-19 | Cross-repo proposal intake path: external proposal absorbed via standard intake | `.cdd/releases/3.78.0/379/` |
| #389 | (cdr Sub dependency) | 2026-05-21 | Loaded-skill miss trigger surfaced (Python-vs-Go skill loaded by mistake); coordination-failure case that ε later patched | `.cdd/unreleased/389/cdd-iteration.md`@71b25672 (3 findings, 2 MCAs, 1 no-patch) |
| #391 | Phase 3 V (superseded) | 2026-05-21 | Mid-flight γ-clarification: #391 superseded by #392 after coordination signal | (referenced in #392's gamma-scaffold.md) |
| #392 | Phase 3 V replanted (supersedes #391) | 2026-05-21 | Mid-flight clarification rescue: 4 MCAs filed; demonstrates the coordination surface's recovery mode | `.cdd/unreleased/392/`@71b25672 |
| #393 | δ-as-architect | 2026-05-21 | Coordination-surface upgrade: the implementation-contract pinning patch that reshaped subsequent inter-cycle coordination | `.cdd/unreleased/393/`@71b25672 |
| #411 | Sub 6 of #403 (parallel to this sub) | 2026-05-22 | Parallel-sub coordination: disjoint file-set discipline enabling parallel α dispatch | issue body `cnos#411` |
| #412 | Sub 7 of #403 (this cycle) | 2026-05-22 | The polling-primitive surface in use: this cycle's branch-state and close-out artifacts are the coordination signal for #403's tracker | `.cdd/unreleased/412/`@`cycle/412` |

### §Artifact contract

CDS's artifact contract (per [`CDS.md §Artifact contract`](../skills/cds/CDS.md)) declares the close-out artifact set, the bootstrap rules, the ordered authoring flow, the manifest, and the **location matrix** + **ownership matrix** for where each artifact lives. The CDS bootstrap wave is the canonical exercise.

| Cycle | Issue | Date | What it exercised | Artifact citation |
|---|---|---|---|---|
| #335 | Iteration aggregator initialized | 2026-05-09 | Initialized `.cdd/iterations/INDEX.md`; bootstrap for the receipt-stream artifact | `.cdd/releases/docs/2026-05-09/335/cdd-iteration.md`@71b25672 |
| #364 | Coherence-cell doctrine | 2026-05-15 | Doctrine-artifact authoring: `COHERENCE-CELL.md` placed under `cnos.cdd/skills/cdd/` per the location matrix | `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md@71b25672` |
| #406 | Sub 1 of #403 — cnos.cds package skeleton | 2026-05-22 | Demonstrated the location-matrix in action: package skeleton + `extraction-map.md` placed under `cnos.cds/docs/` | merge `987acd04`; `src/packages/cnos.cds/docs/extraction-map.md@71b25672` |
| #407 | Sub 2 of #403 — CDS.md authored | 2026-05-22 | The artifact-contract's "doctrine file" sub-class: `CDS.md` authored at canonical location | merge `d9829412`; `src/packages/cnos.cds/skills/cds/CDS.md@71b25672` |
| #408 | Sub 3 of #403 — §Selection + §Lifecycle migrated | 2026-05-22 | Ownership-matrix in action: section ownership transferred from CDD.md to CDS.md | merge `5f13f61c` |
| #409 | Sub 4 of #403 — §Coordination + §Artifact migrated | 2026-05-22 | Recursive: the §Artifact contract itself migrated under §Artifact contract's ownership rule | merge `4a87cdf9` |
| #410 | Sub 5 of #403 — §Mech/§Review/§Gate/§Assessment/§Closure/§Retro/§Non-goals/§Large-file migrated | 2026-05-22 | Eight-section migration under the ownership-matrix; the largest single artifact-contract exercise to date | merge `71b25672` |
| #412 | Sub 7 of #403 (this cycle) | 2026-05-22 | Authors the empirical-anchor artifact-class at its canonical location (`cnos.cds/docs/`) | this file |

### §Mechanical vs judgment

CDS distinguishes **mechanical axes** (grep/wc/ls-checkable; AC verification per [`CDS.md §Mechanical vs judgment`](../skills/cds/CDS.md)) from **judgment axes** (prose-judgment-required; e.g. whether the right surface was chosen). The cycles below exercised the boundary explicitly.

| Cycle | Issue | Date | What it exercised | Artifact citation |
|---|---|---|---|---|
| #392 | Phase 3 V ported to Go | 2026-05-21 | Mechanical side: V passes / V fails is a mechanical check that `cn cdd-verify` emits. Judgment side: was the right surface chosen given #391's mis-frame? Recorded 4 MCAs | `.cdd/unreleased/392/cdd-iteration.md`@71b25672 |
| #393 | δ-as-architect patches | 2026-05-21 | The Rule 7 patch (β must verify implementation-contract coherence) is a judgment-shaped oracle named in mechanical-grep-friendly form | `.cdd/unreleased/393/cdd-iteration.md`@71b25672 |
| #397 | Phase 4 δ split | 2026-05-21 | Judgment-heavy planning cycle: deciding the δ-as-architect / δ-as-effector split required prose judgment, not AC mechanics | `.cdd/unreleased/397/` |
| #398 | Phase 4 harness/release-effector | 2026-05-21 | Mixed: harness mechanics + judgment on release-effector boundary; 2 MCAs filed | `.cdd/unreleased/398/cdd-iteration.md`@71b25672 |
| #402 | Phase 7 CDD.md rewrite | 2026-05-21 | Judgment-heavy: deciding what to quarantine for cds extraction vs what to keep in CDD kernel | `.cdd/unreleased/402/beta-review.md`@71b25672 |
| #410 | Sub 5 of #403 — eight-section migration | 2026-05-22 | Mechanical AC count maxed (AC1–AC17 PASS); judgment surface in β-review on each section's migration faithfulness | `.cdd/unreleased/410/beta-review.md`@71b25672 |

### §Review CLP

CDS's Review CLP (per [`CDS.md §Review CLP`](../skills/cds/CDS.md)) declares the form of β-review.md: claim list, ask list, the structured request-for-info pattern. Every cycle's `beta-review.md` is a CLP exercise; the cycles below are cited as rich representative examples.

| Cycle | Issue | Date | What it exercised | Artifact citation |
|---|---|---|---|---|
| #393 | δ-as-architect | 2026-05-21 | CLP exemplar: 4 patches landed via the ask-list mechanism; demonstrates the CLP as the surface that ratchets doctrine | `.cdd/unreleased/393/beta-review.md`@71b25672 |
| #402 | Phase 7 CDD.md rewrite | 2026-05-21 | CLP under judgment-heavy review; produced 1 MCA (#403) + 1 no-patch finding via the ask-list | `.cdd/unreleased/402/beta-review.md`@71b25672 |
| #407 | Sub 2 of #403 — CDS.md authored | 2026-05-22 | CLP on a large-file authoring cycle (1043 lines); zero round-2 iterations needed | `.cdd/unreleased/407/beta-review.md`@71b25672 |
| #410 | Sub 5 of #403 — eight-section migration | 2026-05-22 | Recent rich CLP example: AC1–AC17 + per-section judgment ACs reviewed in a single round | `.cdd/unreleased/410/beta-review.md`@71b25672 |
| #412 | Sub 7 of #403 (this cycle) | 2026-05-22 | CLP on the empirical-anchor authoring; surface of this very review | `.cdd/unreleased/412/beta-review.md`@`cycle/412` |

### §Gate (F1–F10)

CDS's Gate F1–F10 (per [`CDS.md §Gate`](../skills/cds/CDS.md)) declares the release-readiness preconditions and closure-verification checklist. F1–F10 ratchet from "issue closed" through "evidence linked" through "no regressions" to "receipt typed." The cycles below preserved or exercised the F1–F10 anchors.

| Cycle | Issue | Date | What it exercised | Artifact citation |
|---|---|---|---|---|
| #392 | Phase 3 V port | 2026-05-21 | F2 (V passes) became mechanically grep-able via `cn cdd-verify`; the Go port wired F2's oracle into the cn binary | `src/go/cmd/cn/cmd/cdd_verify.go@71b25672` |
| #398 | Phase 4 harness/release-effector | 2026-05-21 | F8 (release-effector mechanics) implemented; release-effector skill landed under `cnos.cdd/skills/cdd/release-effector/` | `.cdd/unreleased/398/`@71b25672 |
| #400 | Phase 5 γ shrink | 2026-05-21 | F7 (γ close-out artifact set complete) tightened: γ's surface narrowed; 2 patches | `.cdd/unreleased/400/cdd-iteration.md`@71b25672 |
| #402 | Phase 7 CDD.md rewrite | 2026-05-21 | Closure verification on a Phase-class cycle: F1–F10 traversed on the CDD.md rewrite | `.cdd/unreleased/402/gamma-closeout.md`@71b25672 |
| #410 | Sub 5 of #403 | 2026-05-22 | F1–F10 anchors **preserved** during section migration: the gate doctrine moved between files without verdict change | `src/packages/cnos.cds/skills/cds/CDS.md §Gate@71b25672` |

### §Assessment (§9.1 triggers and PRA contents)

CDS's Assessment surface (per [`CDS.md §Assessment`](../skills/cds/CDS.md)) declares the PRA contents, the cycle-iteration triggers ("review-churn", "loaded-skill miss", "AC drift", etc.), the friction log, and the engineering-levels classifier. The cycles below tripped or repaired specific trigger classes.

| Cycle | Issue | Date | What it exercised | Artifact citation |
|---|---|---|---|---|
| #335 | (initialized aggregator) | 2026-05-09 | Bootstrap for the assessment aggregator surface; review-churn trigger recorded as ε scaffolding | `.cdd/releases/docs/2026-05-09/335/cdd-iteration.md`@71b25672 |
| #346 | (docs cycle, 1 no-patch finding) | 2026-05-12 | Assessment surface in low-finding regime: 1 finding, 0 patches, 1 no-patch — exercises the no-patch-acceptable path | `.cdd/releases/docs/2026-05-12/346/cdd-iteration.md`@71b25672 |
| #389 | Loaded-skill miss (Python-vs-Go) | 2026-05-21 | Loaded-skill-miss trigger class fired: wrong skill loaded; 3 findings, 2 MCAs, 1 no-patch | `.cdd/unreleased/389/cdd-iteration.md`@71b25672 |
| #392 | Phase 3 V port | 2026-05-21 | Recovery from #389's miss: the right skill loaded; 4 MCAs filed against the surfaced gaps | `.cdd/unreleased/392/cdd-iteration.md`@71b25672 |
| #393 | δ-as-architect | 2026-05-21 | Implementation-contract-drift trigger class introduced as Rule 7; 4 patches landed on the assessment surface | `.cdd/unreleased/393/cdd-iteration.md`@71b25672 |
| #395 | CDR role overlays (Sub 3 of #376) | 2026-05-21 | Assessment surface on a multi-file-overlay cycle: 2 findings, 1 MCA, 1 no-patch | `.cdd/unreleased/395/cdd-iteration.md`@71b25672 |

### §Closure

CDS's Closure surface (per [`CDS.md §Closure`](../skills/cds/CDS.md)) declares the immediate vs deferred outputs, the closure rule (close issue iff F1–F10 verified), and the direct-to-main exception. Every cycle's `gamma-closeout.md` is a closure exercise; the cycles below are cited for specific closure-mechanism stress.

| Cycle | Issue | Date | What it exercised | Artifact citation |
|---|---|---|---|---|
| #364 | Coherence-cell doctrine | 2026-05-15 | Direct-to-main closure path on a docs-class cycle: docs-only disconnect at merge `32b126e4` | `.cdd/releases/docs/.../364/gamma-closeout.md` |
| #370 | CCNF normal form | 2026-05-17 | Docs-only disconnect closure: γ-as-δ pattern | `.cdd/releases/docs/2026-05-17/370/cdd-iteration.md`@71b25672 |
| #402 | Phase 7 CDD.md rewrite | 2026-05-21 | Closure with deferred output: 1 MCA filed (#403) — closure rule "issue closed iff F1–F10 verified AND deferred outputs filed" | `.cdd/unreleased/402/gamma-closeout.md`@71b25672 |
| #403 | cnos.cds bootstrap tracker | (in-flight) | Tracker-class closure: closure rule deferred to per-sub closure (each sub closes independently; tracker closes when all subs done) | issue body `cnos#403` |
| #406 | Sub 1 of #403 — package skeleton | 2026-05-22 | Closure on a bootstrap cycle: closed at merge `987acd04` | `.cdd/unreleased/406/gamma-closeout.md`@71b25672 |
| #410 | Sub 5 of #403 — eight-section migration | 2026-05-22 | Closure with courtesy `cdd-iteration.md` stub (zero findings); demonstrates the closure-stub option | `.cdd/unreleased/410/cdd-iteration.md`@71b25672 |

### §Retro-packaging

CDS's Retro-packaging surface (per [`CDS.md §Retro-packaging`](../skills/cds/CDS.md)) declares how a cycle that did not declare its full scope at filing-time records the retro-snapshot: what was actually done, what was deferred, what the surface footprint became. The cycles below demonstrate the retro-packaging pattern.

| Cycle | Issue | Date | What it exercised | Artifact citation |
|---|---|---|---|---|
| #364 | Coherence-cell doctrine | 2026-05-15 | Retro-snapshot in close-out: the cell model's full footprint became visible only at γ; recorded in close-out artifacts | `.cdd/releases/docs/.../364/` |
| #392 | Phase 3 V port (supersedes #391) | 2026-05-21 | Retro-packaging on a supersession: #392's footprint absorbed #391's planned scope plus the rescue-mode additions | `.cdd/unreleased/392/`@71b25672 |
| #393 | δ-as-architect patches | 2026-05-21 | Retro-snapshot on a patch-class cycle: 4 patches landed; retro-packaged as Rule 7 + dependent doctrine | `.cdd/unreleased/393/`@71b25672 |
| #402 | Phase 7 CDD.md rewrite | 2026-05-21 | Retro-packaging on a doctrine-rewrite cycle: the "Software-specific realization — pending cds extraction" section is itself a retro-snapshot of CDD.md's pre-rewrite content | `src/packages/cnos.cdd/skills/cdd/CDD.md@71b25672` |

### §Large-file authoring rule

CDS's Large-file authoring rule (per [`CDS.md §Large-file authoring rule`](../skills/cds/CDS.md)) declares the file-size threshold, the section-manifest HTML-comment header, the resumption protocol, and the anti-patterns to avoid when authoring files over ~500 lines. The cycles below exercised the rule on real large-file authoring.

| Cycle | Issue | Date | What it exercised | Artifact citation |
|---|---|---|---|---|
| #364 | Coherence-cell doctrine | 2026-05-15 | Authored `COHERENCE-CELL.md` (multi-hundred lines) under the precursor of the large-file rule | `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md@71b25672` |
| #402 | Phase 7 CDD.md rewrite | 2026-05-21 | Large-file authoring on a heavily-edited doctrine file (CDD.md); demonstrated the section-manifest discipline | `src/packages/cnos.cdd/skills/cdd/CDD.md@71b25672` |
| #407 | Sub 2 of #403 — CDS.md authored | 2026-05-22 | Authored `CDS.md` at 1043 lines under the explicit large-file rule; section-manifest header used | `src/packages/cnos.cds/skills/cds/CDS.md@71b25672` |
| #408 | Sub 3 of #403 — §Selection + §Lifecycle migrated | 2026-05-22 | Large-file editing on CDS.md via thin-extract pattern; section-by-section authoring discipline | merge `5f13f61c` |
| #409 | Sub 4 of #403 — §Coordination + §Artifact migrated | 2026-05-22 | Continued large-file editing under the rule; two-section migration per cycle to bound diff size | merge `4a87cdf9` |
| #410 | Sub 5 of #403 — eight-section migration | 2026-05-22 | Stress-test of large-file authoring rule: eight sections migrated in one cycle (B-lite thin extract) | merge `71b25672` |

---

## Closure cycles and INDEX rows

The canonical receipt stream for CDS evidence lives in
[`.cdd/iterations/INDEX.md`](../../../../.cdd/iterations/INDEX.md) (32 rows
as of cycle/412, spanning #331 through #410). Each row is one
`cdd-iteration.md` record — a per-cycle CDS receipt. The aggregator
columns (`Findings`, `Patches`, `MCAs`, `No-patch`) are exactly the
longitudinal ε surface that [`CDS.md §Field 5 — ε iteration cadence`](../skills/cds/CDS.md)
names: trigger-class counts across cycles.

A compressed view of the receipt stream for the cycles cited above:

| Cycle | Findings | Patches | MCAs | No-patch | Significance |
|-------|---------:|--------:|-----:|---------:|--------------|
| #335 | 2 | 1 | 1 | 0 | Aggregator initialized |
| #369 | 3 | 0 | 2 | 1 | Phase 2 schemas |
| #370 | 2 | 1 | 1 | 0 | CCNF normal form |
| #379 | 3 | 3 | 0 | 0 | cn-sigma intake |
| #388 | 3 | 0 | 2 | 0 | Phase 2.5 schemas split |
| #389 | 3 | 0 | 2 | 1 | Loaded-skill miss |
| #390 | 2 | 0 | 2 | 0 | CDR.md Sub 1 |
| #392 | 4 | 0 | 4 | 0 | Phase 3 V (Go) |
| #393 | 0 | 4 | 0 | 0 | δ-as-architect (all patches) |
| #395 | 2 | 0 | 1 | 1 | CDR role overlays |
| #398 | 2 | 0 | 2 | 0 | Phase 4 harness |
| #400 | 0 | 2 | 1 | 0 | Phase 5 γ shrink |
| #402 | 2 | 0 | 1 | 1 | Phase 7 CDD.md rewrite |
| #406–#410 | 0 | 0 | 0 | 0 | CDS Subs 1–5 (zero-finding migrations) |

The "zero-finding" CDS-migration cycles (#406–#410) are themselves a
finding about the surface: extract-by-reference v0.1 was scoped to
prevent doctrine drift, and the receipt stream confirms that no drift
gap surfaced during the migration wave. That is the receipt stream
discharging the closure-bind for Sub 7's empirical-anchor obligation.

---

## Open questions and forthcoming work

- **v1.0 deep per-role rewrites.** Deferred from #403 v0.1. The v0.1 wave
  ships extract-by-reference (CDS.md + supporting docs); v1.0 will produce
  full per-role SKILL.md files under `cnos.cds/skills/cds/{alpha,beta,gamma,delta,epsilon}/`
  modelled after the cnos.cdr role-overlay structure landed by cnos#395.
- **`.cdd/` → `.cds/` filesystem re-rooting.** Documented in
  [`docs/extraction-map.md §14`](./extraction-map.md) (or successor
  ordering); not performed. The pre-rename naming is preserved in this
  doc's filename (`empirical-anchor-cdd.md`) and in the
  `.cdd/iterations/INDEX.md` aggregator path.
- **CCNF-O orchestration grammar.** [cnos#405](https://github.com/usurobor/cnos/issues/405); gated on
  cnos#403 closure. The orchestration layer above CDS that allows multiple
  CDS-conformant projects to compose.
- **Handoff / coordination package extraction.** [cnos#404](https://github.com/usurobor/cnos/issues/404);
  gated on cnos#403 closure. Splits the coordination surfaces (cross-repo
  proposals, mid-flight clarification, polling primitives) into a separate
  package if the surface footprint justifies it.
- **`#CDSReceipt` schema-driven receipts.** Once the typed receipt schema
  at `schemas/cds/receipt.cue` is wired through the close-out path, the
  `cdd-iteration.md` per-cycle record becomes a typed instance of
  `#CDSReceipt`. The mapping in this doc will then anchor against typed
  schema fields rather than free-form prose rows.
- **Surface-by-surface gap distribution.** A future revision of this doc
  could add a per-surface column to the aggregator: which CDS surface
  each cycle's findings clustered against. The data is present in the
  per-cycle `cdd-iteration.md` files; the aggregator currently rolls up
  by count only.

---

## Related documents

This document cites:

- [`src/packages/cnos.cds/skills/cds/CDS.md`](../skills/cds/CDS.md) — the
  cnos.cds v0.1 doctrinal contract. All surface citations (§Six-field,
  §Selection function, §Development lifecycle, §Coordination surfaces,
  §Artifact contract, §Mechanical vs judgment, §Review CLP, §Gate,
  §Assessment, §Closure, §Retro-packaging, §Large-file authoring rule)
  target this file. The "Empirical anchor → Sub 7 deferred surface-by-surface mapping" sub-section is the seed that this document verifies.
- [`src/packages/cnos.cdd/skills/cdd/CDD.md`](../../cnos.cdd/skills/cdd/CDD.md) —
  the kernel doctrine (CCNF spine) that CDS extends. CDS inherits the
  six-field structure from CDD via [`ROLES.md §3`](../../../../ROLES.md);
  CDS is the software-domain instantiation.
- [`src/packages/cnos.cdr/docs/empirical-anchor-cph.md`](../../cnos.cdr/docs/empirical-anchor-cph.md) —
  the structural precedent. The cph anchor maps research artifact classes
  onto CDR fields; this CDS anchor maps cnos cycles onto CDS surfaces.
  The two documents are sibling instances of the same anchor pattern.
- [`src/packages/cnos.cds/docs/extraction-map.md`](./extraction-map.md) —
  the source-to-destination map for the #403 migration wave. Now
  historical: Subs 1–5 are complete and the extraction-map's Status
  columns reflect that.
- [`.cdd/iterations/INDEX.md`](../../../../.cdd/iterations/INDEX.md) — the
  canonical receipt stream. Every `cdd-iteration.md` row referenced in
  this doc lives under one of the paths listed there.
- [`ROLES.md`](../../../../ROLES.md) — the six-field instantiation contract
  (`§3`) and the five-layer enforcement chain (`§4a`) that CDS realizes.
- [usurobor/cnos issues #364 through #412](https://github.com/usurobor/cnos/issues) —
  the cycle stream cited row-by-row.

---

## Closing — cnos#412 AC1–AC8 satisfied

[cnos#412](https://github.com/usurobor/cnos/issues/412)'s acceptance criteria are discharged by this document:

- **AC1 (file exists at canonical path, ≥ 100 lines).** This file is at
  `src/packages/cnos.cds/docs/empirical-anchor-cdd.md` and exceeds the
  line threshold.
- **AC2 (≥ 5 `^## ` sections).** This doc carries `## How to read this
  doc`, `## Cycle-to-surface mapping`, `## Closure cycles and INDEX
  rows`, `## Open questions and forthcoming work`, `## Related
  documents`, and `## Closing — cnos#412 AC1–AC8 satisfied`.
- **AC3 (≥ 5 CDS surfaces have mapping tables).** Twelve surfaces have
  mapping tables: §Six-field, §Selection function, §Development
  lifecycle, §Coordination surfaces, §Artifact contract, §Mechanical vs
  judgment, §Review CLP, §Gate, §Assessment, §Closure, §Retro-packaging,
  §Large-file authoring rule.
- **AC4 (≥ 20 distinct cnos cycle numbers cited).** Cycles cited above:
  #335, #346, #364, #366, #367, #369, #370, #376, #379, #388, #389, #390,
  #391, #392, #393, #395, #397, #398, #400, #401, #402, #403, #406, #407,
  #408, #409, #410, #411, #412. Twenty-nine distinct cycles.
- **AC5 (cnos.cdd untouched).** Verified by `git diff origin/main..HEAD -- src/packages/cnos.cdd/`.
- **AC6 (cnos.cdr untouched).** Verified by `git diff origin/main..HEAD -- src/packages/cnos.cdr/`.
- **AC7 (cnos.cds/skills/cds/CDS.md untouched).** Verified by `git diff origin/main..HEAD -- src/packages/cnos.cds/skills/cds/CDS.md`.
- **AC8 (Related documents section cites the four required peers).**
  CDS.md, CDD.md, cph empirical anchor, and INDEX.md are all cited in
  [§Related documents](#related-documents) above.

**cnos#412 AC1–AC8 are therefore satisfied** by this document at cnos
commit `71b25672`.
