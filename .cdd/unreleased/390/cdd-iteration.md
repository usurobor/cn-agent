# cdd-iteration — Cycle #390 (Sub 1 of #376)

**Cycle:** #390 — Sub 1 of #376: CDR six-field instantiation contract + architecture-choice declaration
**Merge:** [filled by γ post-merge in merge SHA]
**Closed:** 2026-05-21 (auto-closed via `Closes #390` in merge commit)
**Mode:** design-and-build (γ+α+β-collapsed-on-δ per breadth-2026-05-12 wave manifest precedent + cycles 375/377/378/388 empirical validation; β-α-collapse acknowledged in `beta-review.md`)
**Rounds:** R1 APPROVE (no fix-round)
**ACs:** 6/6 PASS

## §1 Findings dispositioned

### F1: CDR gate-verdict vocabulary not enum-pinned in `schemas/cdr/receipt.cue`

- **Source:** cycle 390 α-closeout F2 + self-coherence §Debt item 4 + design-notes §OpenQuestions item 1.
- **Class:** `cdd-protocol-gap` (borderline — the schema is intentionally generic-typed; the gap is "should CDR pin a domain-specific verdict enum, or stay with the generic boundary_decision.action + transmissibility typed surface?").
- **Trigger:** ε process-gap check during cycle close-out; surfaced by the doctrinal-vocabulary-vs-typed-surface mismatch.
- **Description:** CDR.md Field 3 declares the gate-verdict vocabulary as GO / REVISE / NO-GO / INDETERMINATE / BOUNDED-GO. The typed surface in `schemas/cdr/receipt.cue` uses the generic `boundary_decision.action` enum (`accept | release | reject | repair_dispatch | override`) + `transmissibility` enum (`accepted | degraded | not_transmissible`) inherited from `schemas/cdd/`. CDR.md maps the verdict vocabulary onto the generic typed surface abstractly (e.g. BOUNDED-GO ↔ `transmissibility: degraded` + `limitations` enumerated). The mapping works but is not pinned in the schema. If cph or a downstream research project's verdict verbs drift (e.g. introducing a new verdict like `PROVISIONAL-GO`), the schema does not signal the drift structurally; V would accept any string mapping abstractly.
- **Root cause:** Sub 1 deliberately keeps the schema generic (per cnos#388's discipline — generic kernel does not name domain enums by value). Pinning a CDR-specific verdict enum would belong in a `schemas/cdr/verdict.cue` extension; that file does not exist yet.
- **Disposition:** `next-MCA`. Patch shape: "if cph or downstream research projects develop a stable verdict-name set that differs from the generic enum, add `schemas/cdr/verdict.cue` declaring a `#CDRVerdict` enum and reference it from `schemas/cdr/receipt.cue` and from `CDR.md §"Field 3"`." Best landed when Sub 4 (cph empirical-anchor doc) confirms cph's actual verdict vocabulary stabilises.
- **Issue filed:** none yet. Suggest filing alongside Sub 4 dispatch if cph's verdict set is concrete; otherwise carry as known debt.
- **First AC for the eventual MCA:** `schemas/cdr/verdict.cue` declares `#CDRVerdict: "GO" | "REVISE" | "NO-GO" | "INDETERMINATE" | "BOUNDED-GO"` (or whatever cph's stabilised set is); `CDR.md §"Field 3"` cites it; `schemas/cdr/receipt.cue` includes the verdict field structurally.

### F2: Project-specific actor-collapse stricter-floor policy template not provided

- **Source:** cycle 390 α-closeout debt item 5 + self-coherence §Debt item 5 + design-notes §OpenQuestions item 3.
- **Class:** `cdd-protocol-gap` (borderline — the doctrinal floor is set; the project-binding template is not).
- **Trigger:** ε process-gap check; surfaced by Field 6's "project-specific stricter floors are permitted" clause without a binding template.
- **Description:** CDR.md Field 6 names a doctrinal floor (α=β never for research claims; γ=δ and ε=δ permitted with conditions) and notes that project-specific stricter floors are permitted (e.g. "for externally-published claims, β must be a distinct human reviewer"). No template for `<project>/.cdr/POLICY.md` or equivalent is provided. A future research project wanting a stricter floor would have to derive its own structure.
- **Root cause:** Sub 1's scope is doctrinal contract authoring; project-binding templates are arguably out of CDR.md's scope (they belong to per-project policy authoring or to a separate "CDR project-binding starter" skill).
- **Disposition:** `next-MCA`. Patch shape: either (a) author a starter template `.cdr/POLICY.md.template` under `cnos.cdr/skills/cdr/` (Sub 3 or later), or (b) defer to cph as the canonical example — cph's own `.cdr/POLICY.md` (if/when authored) becomes the de-facto template via Sub 4's mapping.
- **Issue filed:** none yet. Suggest folding into Sub 3 (if a template is wanted alongside role overlays) or into Sub 4 (if cph's actual policy file is sufficient as the worked example).
- **First AC for the eventual MCA:** either `cnos.cdr/skills/cdr/templates/POLICY.md.template` exists with project-policy stub, **or** Sub 4's empirical-anchor doc cites `cph/.cdr/POLICY.md` as the worked example with structural commentary.

## §2 No-findings observations (informational)

- **AC1 release/deploy/tag token classification:** all 6 hits classified as disavowing-context or label-tag (non-release-tag) usage. The cycle's `tag`-in-Field-5 usage ("tagged with one of the classes above") is a label-tag distinct from release-tag; a future ε pattern-match scanning for release-vocabulary in CDR-class files should account for the label-tag semantic to avoid false-positives. Recorded; not raised to a finding because the AC oracle is satisfied as-is.
- **β-α-collapse vs Field 6 borderline:** the borderline observation in `beta-review.md` (this cycle's β-α-collapse vs Field 6's α=β prohibition) is reconciled by class-identification — docs-only contract authoring is CDS-class under repairable feedback, not research-claim authoring. The reconciliation prevents the same-session collapse pattern from being misread as a research-class waiver template. The discipline is held; no finding.
- **Architectural-choice inheritance pattern validates as reusable design principle.** cnos#388's rationale-5 ("decision-once-applied-twice") is now applied twice (once for schemas, once for skills). The pattern is reinforced as a c-d-X design principle. No finding; positive validation.

## §3 Trigger assessment (per `gamma/SKILL.md §2.8` table)

| Trigger | Fire condition | Fired? | ε note |
|---|---|---|---|
| Review churn | review rounds > 2 | **No** | R1 APPROVE on first pass. |
| Mechanical overload | mechanical ratio > 20% AND findings ≥ 10 | **No** | 7 β observations + 1 borderline; all confirming or reconciled. Mechanical ratio is high by design (contract-authoring class) but findings ≥ 10 is not met. |
| Avoidable tooling / environment failure | environment blocked the cycle | **No** | No tooling friction; the only known harness pattern (origin push --delete 403 for branch cleanup) is non-blocking and named in wave-2026-05-19 ε iteration F5. |
| Loaded-skill miss | a loaded skill should have prevented a finding | **No** | No findings raised. F1/F2 above are forward-looking next-MCA candidates, not retroactive skill misses; they surface from the doctrinal layer being completed for the first time. |

No trigger fires. Cycle ran cleanly.

## §4 INDEX update

Add to `.cdd/iterations/INDEX.md`:

```
| 390 | #390 | 2026-05-21 | 2 | 0 | 2 | 0 | .cdd/unreleased/390/cdd-iteration.md |
```

Findings: 2 (F1, F2). Patches: 0 immediate (both are forward-looking next-MCAs). MCAs: 2 (F1, F2). No-patch: 0. Path: `.cdd/unreleased/390/cdd-iteration.md` (will move to `.cdd/releases/<version>/390/cdd-iteration.md` at next release per `release/SKILL.md §2.5a`).

## §5 Skill-gap candidate disposition

Both F1 and F2 are `cdd-protocol-gap` candidates from the new CDR protocol. Per ε's MCA discipline:
- F1 patches when cph's verdict vocabulary stabilises (Sub 4 dispatch is the natural confluence).
- F2 patches when either (a) Sub 3 adds a starter template, or (b) Sub 4's cph mapping promotes cph's policy file as the worked example.

Both can be left as `next-MCA` until the natural trigger arrives; neither requires immediate patch.

## §6 Deferred outputs

- **Cycle-dir move.** `.cdd/unreleased/390/` → `.cdd/releases/<version>/390/` at next release per `release/SKILL.md §2.5a`. Not blocking; standard release-time mechanic.
- **F1 schema verdict-enum pin** — deferred to Sub 4 or later; recorded as design-notes §OpenQuestions item 1 + this iteration F1.
- **F2 project-policy template** — deferred to Sub 3 or Sub 4; recorded as design-notes §OpenQuestions item 3 + this iteration F2.

## §7 Next-MCA commitment

Sub 2 (cnos.cdr package skeleton), Sub 3 (role overlays), Sub 4 (empirical-anchor doc) can all dispatch in parallel against CDR.md as their stable doctrinal target. cnos#376 close-out comment names this commitment.

Filed by ε on 2026-05-21.
