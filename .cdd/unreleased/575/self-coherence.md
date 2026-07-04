# α self-coherence — cnos#575

<!--
section-manifest:
  planned: [gap, skills, acs, self-check, debt, cdd-trace, review-readiness]
  completed: [gap, skills]
-->

## §Skills

**Tier 1 (canonical lifecycle + role contract):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` (implicit via `alpha/SKILL.md`'s load order)
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — the α role surface; load order followed: CDD.md → alpha/SKILL.md → lifecycle sub-skills → Tier 2/3.

**Lifecycle sub-skills loaded (per alpha/SKILL.md §2.1 step 6):**
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — AC-boundary interpretation.
- `src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md` — proof-plan discipline (positive/negative case shape).
- `src/packages/cnos.cds/skills/cds/CDS.md` — canonical step table / artifact contract (referenced for location matrix; not read line-by-line given its size, consulted for §"Artifact contract" → §"Location matrix" and §"Coordination surfaces" only).

**Tier 2 (always-applicable engineering skills), resolved to their actual repo location** — the dispatch prompt named `cnos.core/skills/eng/*`, but those skills live under `cnos.eng/skills/eng/*` in this repo (`cnos.core` only carries `mindsets/ENGINEERING.md`, not the `eng/*` skill tree); loaded from the correct path, not the prompt's literal path, since the governing content (not the path string) is what the prompt intends:
- `src/packages/cnos.eng/skills/eng/ship/SKILL.md` — TDD flow (spec → failing test → code → passing test → full suite → ship); §"Bug Fix Flow (TDD)" is the operative shape here (this cycle wires new guards into an existing, working evaluator, not a from-scratch feature).
- `src/packages/cnos.eng/skills/eng/test/SKILL.md` — invariant-first testing, positive/negative case discipline, "cover new surfaces" (§3.13).
- `src/packages/cnos.eng/skills/eng/go/SKILL.md` — Go package/type/dispatch-boundary conventions (skimmed via table-of-contents; the existing `issuesfsm` package already conforms, and this cycle's additions follow its established patterns exactly — new guard funcs, new FactSnapshot fields, new fetch.go switch cases, mirroring the `review_request_present`/`REVIEW-REQUEST.yml` precedent named throughout the scaffold).
- `src/packages/cnos.core/skills/write/SKILL.md` — prose/doc-authoring discipline for the `cds-dispatch/SKILL.md` and `dispatch-protocol/SKILL.md` edits.

**Tier 3 (issue-specific):** the scaffold's per-AC oracle list and friction notes (`.cdd/unreleased/575/gamma-scaffold.md`) functioned as the issue-specific Tier 3 surface for this cycle — read in full before any code was written, per α's dispatch-intake rule 5.

**Source files read before coding (peer enumeration of the existing implementation, per alpha/SKILL.md §2.1 step 5):** `transitions.json`, `table.go`, `snapshot.go`, `fetch.go`, `decision.go`, `issuesfsm.go`, `issuesfsm_test.go` (all existing tests + all existing `testdata/*.json` fixtures), `cmd_issues_fsm.go`, `cds-dispatch/SKILL.md` (full body), `dispatch-protocol/SKILL.md` (full body), and `delta/SKILL.md` §9.6 (grep-scoped, per Friction note 3).

## §Gap

**Issue:** [usurobor/cnos#575](https://github.com/usurobor/cnos/issues/575) — "cds/fsm: route claim, hard-block, release-back-to-queue through the FSM (Phase 3 — Sub 2 of #583)".

**Mode:** design-and-build. **Parent:** #583 (master wave). **Precondition:** #584 (Sub 1 doctrine) — landed, CLOSED.

**Gap statement (from the issue).** The CDS issue-state FSM (`cn issues fsm`, #568/#569, hardened #574) owns exactly one status transition: `in-progress → review`. Claim (`todo → in-progress`), hard-block (`→ status:blocked`), and release-back-to-queue (`in-progress → todo`) remain direct wake label writes, not routed through the FSM at all — so "the FSM owns status labels" was still aspirational, and the mechanical PR-open/recovery runtime (Subs 3-4 of #583) must not be built on top of direct label writes.

**Closure condition.** Claim, hard-block, and release-back-to-queue are FSM-applied on passing guards (TDD fixtures), the wake requests rather than writes, the #574/#569 invariants and `cell_kind` observed-only hold, and all gates are green.

**Branch:** `cycle/575`, created by γ from `origin/main`. γ's R0 scaffold (`.cdd/unreleased/575/gamma-scaffold.md`, commit `573e6bf`) is the load-bearing artifact this implementation follows — its per-AC oracle list, surfaces-to-touch table, source-of-truth table, and six named friction notes are all addressed below.

**Implementation contract (pinned by δ, restated for the record — unchanged from the scaffold, no deviation):**

| Axis | Pinned value | Honored? |
|---|---|---|
| Language | Go, no new runtime | Yes — all engine code is Go under the existing `issuesfsm` package |
| CLI integration target | Extends `cn issues fsm evaluate [--apply]`; no new sub-verb, no new binary | Yes — `cmd_issues_fsm.go` untouched, zero new flags |
| Package scoping | `cnos.cds/skills/cds/fsm/transitions.json` (data) + `cnos.issues/commands/issues-fsm/{table.go,snapshot.go,fetch.go,issuesfsm_test.go,testdata/*.json}` (Go) + `cds-dispatch/SKILL.md` + `dispatch-protocol/SKILL.md` (doctrine) | Yes — see §CDD Trace step 6 file list; zero files outside this set except the two rendered artifacts the doctrine change requires (golden fixture + live workflow, both machine-derived from `cds-dispatch/SKILL.md`) |
| Existing-binary disposition | Extend `cn`, no new binary | Yes |
| Runtime dependencies | None beyond existing Go toolchain + vendored GitHub REST client in fetch.go | Yes — zero new imports outside stdlib |
| JSON/wire contract preservation | `transitions.json` schema additive-only; `review`/`changes` state blocks unchanged in shape/outcome | Yes — see §ACs AC5 |
| Backward-compat invariants | `run_active` non-gating for review-request path; #574 review-guard tightening unchanged; `cell_kind` observed-only | Yes — see §ACs AC5 |

No axis was improvised or relaxed; no unpinned row was encountered.
