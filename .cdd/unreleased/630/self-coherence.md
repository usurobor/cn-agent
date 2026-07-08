# self-coherence — cnos#630

## Gap

**Issue:** [usurobor/cnos#630](https://github.com/usurobor/cnos/issues/630) — "partial-matter in-progress wedge: give the reconciler a mechanical path out."

**Mode:** design-and-build. **Cell kind:** `implementation` (reconciler gap fix).

**The gap.** A `status:in-progress` dispatch cell whose run dies *after* producing matter (a `cycle/{N}` branch + a checkpointed draft PR, per the cnos#591 finalizer) but *before* `REVIEW-REQUEST.yml` is written had **no mechanical path out**: it could not requeue (the cnos#575 release-back-to-queue guard blocks over existing matter — the cnos#368 blind-requeue protection), could not advance (no `REVIEW-REQUEST.yml`, so `in-progress -> review` is blocked), and was never re-claimed (the dispatch selector only claims `status:todo`). Concrete evidence: cnos#614 sat wedged in `status:in-progress` for over a day across many scheduled sweeps, recoverable only by manual κ (operator) intervention. The literal wedge was `scan.go`'s `else` branch of the `propose_delta_recovery` case (the pre-fix `scan.go:252-256`), which re-evaluated to the identical no-op decision on every scan tick once matter was checkpointed.

**Round:** R0 (first implementation pass on `cycle/630`, `run_class: first_pass` per the γ scaffold — clean first dispatch, no prior claim history).

**Base:** `cycle/630`, created from `main@6143b53c9098e415c4e7e83791843bcad2313314` per `.cdd/unreleased/630/gamma-scaffold.md`.

## Skills

**Tier 1 (always loaded):** `CDD.md`; `cnos.cdd/skills/cdd/alpha/SKILL.md` (this role's contract — §2.1 dispatch intake, §2.5 self-coherence authoring, §2.6 pre-review gate, §2.7 request-review signal, all followed).

**Tier 2 (always-applicable engineering skills, per `cnos.eng/skills/eng/README.md`):** Go authoring conventions (table-driven evaluator discipline — `table.go`'s own doc comment: "never switched on a CDS-specific state name... entirely table-driven" — honored by adding the new rule as declarative JSON data, not a Go `switch` branch); test-first discipline for behavior changes (every fixture/action-name change is backed by a table-driven or `RunScan`-level test).

**Tier 3 (issue-specific):**
- `cnos.core/skills/skill` — not invoked as a distinct skill-authoring pass; the diff amends existing `SKILL.md` prose (`cds-dispatch/SKILL.md`, `delta/SKILL.md`) rather than authoring a new skill file, so the existing files' own section/heading conventions were followed by direct read-and-match rather than a separate skill-authoring skill load.
- The γ scaffold's "Recommended design" and "Open design call" sections (`.cdd/unreleased/630/gamma-scaffold.md`) functioned as the binding Tier-3 design constraint for this cycle — read in full before any code change, per α's dispatch-intake discipline (`alpha/SKILL.md` §2.1 step 4).
- `cnos.cds/skills/cds/fsm/transitions.json`'s own `_doc`/`_doc_phase2`/`_doc_phase3` blocks were read as the binding authoring convention for how new rules must be documented (each `_doc*` block names the issue that added the rule and the rationale) — the new rule in this diff follows that convention (a `reason` string naming cnos#630, the mechanism, and the guardrail it preserves).

**Implementation contract (pinned by δ; not relaxed):** Go; `cn issues fsm evaluate`/`scan` subcommand family, no new binary; package scoping `src/packages/cnos.issues/commands/issues-fsm/` + `src/packages/cnos.cds/skills/cds/fsm/transitions.json` + `cds-dispatch/SKILL.md` + `delta/SKILL.md` §9; existing-binary disposition preserved; no new runtime dependencies; JSON/wire contract additive-only; backward-compat preserved except the one AC2-redirected test (named in §ACs below). Every diff hunk in this cycle maps to one or more of these pinned rows — no improvisation on language, package location, or binary shape.
