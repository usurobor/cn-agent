# α self-coherence — cnos#575

<!--
section-manifest:
  planned: [gap, skills, acs, self-check, debt, cdd-trace, review-readiness]
  completed: [gap, skills, acs]
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

## §ACs

Evidence below is real command output run against branch HEAD (implementation commits `11b82d2` / `489aeab`, before this readiness section), not a manual claim.

### AC1 — claim routed through the FSM

**Rule (new, `transitions.json`'s `todo` block):**

```
$ jq '.transitions[] | select(.state=="todo")' src/packages/cnos.cds/skills/cds/fsm/transitions.json
{
  "state": "todo",
  "trigger": "evaluate",
  "rules": [
    {
      "all_true": ["claim_request_present"],
      "all_false": ["run_active"],
      "outcome": "proposed",
      "action": "propose_status_in_progress",
      "target_state": "in-progress",
      ...
    },
    {
      "all_true": ["claim_request_present"],
      "outcome": "blocked",
      "action": "block",
      "target_state": "",
      ...
    },
    { "outcome": "valid", "action": "none", ... }
  ]
}
```

**Test evidence:**

```
$ go test ./... -run TestAC575_1 -v
=== RUN   TestAC575_1_ClaimRoutedThroughFSM
--- PASS: TestAC575_1_ClaimRoutedThroughFSM (0.00s)
=== RUN   TestAC575_1_ClaimBlockedOverCompetingRun
--- PASS: TestAC575_1_ClaimBlockedOverCompetingRun (0.00s)
=== RUN   TestAC575_1_ClaimNotRequestedStaysValid
--- PASS: TestAC575_1_ClaimNotRequestedStaysValid (0.00s)
PASS
```

**Live CLI evidence (real `cn` binary, `--fixture` offline path):**

Positive (`testdata/todo-claimable.json`, `claim_request_present=true`, `run_active` absent):
```
Current state: todo
Decision:
  outcome: proposed
  enabled_transition: todo -> in-progress
  proposed_action: propose_status_in_progress
```

Negative (`testdata/todo-competing-run.json`, `claim_request_present=true`, `run_state=in_progress`):
```
Current state: todo
Decision:
  outcome: blocked
  enabled_transition: (none)
  proposed_action: block
```

**Doctrine evidence:** `cds-dispatch/SKILL.md` claim sequence steps 5-6 now write `CLAIM-REQUEST.yml` then request `cn issues fsm evaluate --issue {N} --apply` instead of `gh issue edit --remove-label status:todo --add-label status:in-progress` (commit `489aeab`); the "Lifecycle transitions" table's claim row no longer says "direct label write" (see AC4 below for the `rg` oracle covering this same site).

**Status: MET.** Positive and negative cases both proven at the `Evaluate()` level, the CLI `--apply` level (see `TestAC575_ApplyClaimTransitionAppliesOnGuardPass` / `TestAC575_ApplyClaimBlockedRefusesAndMutatesNothing`), and the doctrine level.

### AC2 — hard-block routed through the FSM

**Rule (new, `transitions.json`'s `in-progress` block, positioned before the `run_active` valid rule):**

```
$ jq '.transitions[] | select(.state=="in-progress") | .rules[] | select(.action=="propose_status_blocked")' src/packages/cnos.cds/skills/cds/fsm/transitions.json
{
  "all_true": ["block_request_present"],
  "outcome": "proposed",
  "action": "propose_status_blocked",
  "target_state": "blocked",
  ...
}
```

**Test evidence:**
```
$ go test ./... -run TestAC575_2 -v
=== RUN   TestAC575_2_HardBlockRoutedThroughFSM
--- PASS: TestAC575_2_HardBlockRoutedThroughFSM (0.00s)
=== RUN   TestAC575_2_HardBlockRefusedWithoutEvidence
--- PASS: TestAC575_2_HardBlockRefusedWithoutEvidence (0.00s)
PASS
```

**Live CLI evidence:**

Positive (`testdata/in-progress-block-with-evidence.json`, `block_request_present=true`, `run_state=in_progress` — proves non-gating on `run_active`, mirroring the cnos#569 review rule):
```
Current state: in-progress
Decision:
  outcome: proposed
  enabled_transition: in-progress -> blocked
  proposed_action: propose_status_blocked
```

Negative (`testdata/in-progress-block-no-evidence.json`, `block_request_present=false`):
```
Current state: in-progress
Decision:
  outcome: valid
  enabled_transition: (none)
  proposed_action: none
  reason: workflow run is active (queued/in_progress); awaiting completion, no transition proposed.
```
(Falls through to the pre-existing `run_active` rule, exactly as designed — no hard-block is silently applied without evidence.)

**Status: MET.**

### AC3 — release-back-to-queue routed through the FSM

**Rules (new, `transitions.json`'s `in-progress` block, positioned before the existing dead-run reconciliation rules):**

```
$ jq '.transitions[] | select(.state=="in-progress") | .rules[] | select(.all_true[0]=="release_request_present")' src/packages/cnos.cds/skills/cds/fsm/transitions.json
{
  "all_true": ["release_request_present"],
  "all_false": ["branch_has_commits", "pr_exists", "pr_has_commits"],
  "outcome": "proposed", "action": "propose_status_todo", "target_state": "todo", ...
}
{
  "all_true": ["release_request_present"],
  "any_true": ["branch_has_commits", "pr_exists", "pr_has_commits"],
  "outcome": "proposed", "action": "propose_delta_recovery", "target_state": "", ...
}
```

**Test evidence:**
```
$ go test ./... -run TestAC575_3 -v
=== RUN   TestAC575_3_ReleaseRoutedThroughFSMWhenNoMatter
--- PASS: TestAC575_3_ReleaseRoutedThroughFSMWhenNoMatter (0.00s)
=== RUN   TestAC575_3_ReleaseBlockedOverMatter
--- PASS: TestAC575_3_ReleaseBlockedOverMatter (0.00s)
PASS
```

**Live CLI evidence:**

Positive, no matter (`testdata/in-progress-release-no-matter.json`, `release_request_present=true`, `run_state=in_progress`, no commits/PR):
```
Decision:
  outcome: proposed
  enabled_transition: in-progress -> todo
  proposed_action: propose_status_todo
```

Negative, matter present — the #368 protection (`testdata/in-progress-release-with-matter.json`, `release_request_present=true`, `commits_beyond_base=3`, `pr_exists=true`):
```
Decision:
  outcome: proposed
  enabled_transition: (none)
  proposed_action: propose_delta_recovery
```
(Never `target_state: todo` — never a blind requeue over existing matter.)

**Distinctness from the dead-run reconciliation rules:** `TestAC575_RuleOrderingDoesNotShadowExistingDeadRunRules` proves the pre-existing `in-progress-dead-no-matter.json` / `in-progress-dead-with-commits.json` / `in-progress-active.json` fixtures (none of which set `release_request_present`) still resolve to their original rules/actions/target-states unchanged, despite the new rules being positioned earlier in the same rule list.

**Status: MET.**

### AC4 — workers/wake request, don't write

```
$ rg -n "direct label write" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md
(no output — exit 1, zero matches)

$ rg -n "wake writes the label directly" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md
(no output — exit 1, zero matches)
```

Both oracle greps that the scaffold named — previously 3 matches on the claim/hard-block/release-back-to-queue rows plus the "wake writes the label directly" lead-in prose — now return zero. All four named transitions (claim, hard-block, release-back-to-queue, β-converge) are documented as FSM requests in `cds-dispatch/SKILL.md`'s frontmatter description, claim-sequence steps 5-6-8, "Lifecycle transitions" table, and "Responsibilities" item 5 (peer enumeration: all four prose sites carrying the old "direct label write" language were located via `grep -n` before editing and updated together, not just the first hit).

`dispatch-protocol/SKILL.md` — per Friction note 2 (see §Debt below for the framing decision) — gets a scoped "Protocol-specific FSM override" paragraph in §2.2 rather than a rewrite of the generic claim-sequence code block; the generic `gh issue edit` example stays correct as the default for protocols without a package-owned FSM (`cnos.cdr` today).

**Rendered-artifact consistency:** `cds-dispatch/SKILL.md` is compiled into a golden fixture and a live GitHub Actions workflow via `cn-install-wake`. Both were re-rendered after the prose edits and are byte-identical to a fresh render (`cn-install-wake cds-dispatch` reports `(unchanged)` after the edits landed — i.e. no drift between the committed golden/workflow and the source SKILL.md). `scripts/ci/check-dispatch-repair-preflight.sh` and `scripts/ci/check-dispatch-closeout-integrity.sh` both still exit 0 against the updated surface.

**Status: MET.**

### AC5 — invariants preserved; gates green

**Full existing + new `issuesfsm` test suite** (`cd src/packages/cnos.issues/commands/issues-fsm && go test ./... -count=1 -v`):

```
ok  	github.com/usurobor/cnos/packages/cnos.issues/commands/issues-fsm	0.033s
```
51 top-level `--- PASS` lines, 0 `--- FAIL` lines (exact count from `grep -c "^--- PASS"` / `grep -c "^--- FAIL"` against the `-v` run output). 37 are pre-existing (`TestAC1`-`TestAC7`, `TestAC569_*`, `TestAC574_*`, `TestSeam_CellKind*`, `TestApply_*`, `TestObserveRemoteBranch_*`, `TestAssembleLive_*`, etc.) and pass **unmodified** in assertion/meaning — the only pre-existing test touched at all is `TestSeam_CellKindNotEnforced`, whose fixture list was extended (not whose assertion was weakened) per the pinned contract's explicit instruction. 14 are new `TestAC575_*` tests.

Specifically named in the scaffold's AC5 oracle, all present and passing unmodified: `TestAC574_ReviewRequestAloneNoLongerValid`, `TestAC574_ReviewPartialEvidenceBlocked`, `TestAC574_ReviewWithPRStillValid`, `TestAC574_InProgressBranchOnlyNoLongerProposesReview`, `TestAC574_InProgressWithMatterStillProposesReview`, `TestAC574_InProgressRunActiveNonGatingPreserved`, `TestSeam_CellKindNotEnforced`, `TestSeam_CellKindDefaultedWhenAbsent`, `TestAC4_DeadInProgressNoMatterProposesRequeue`, `TestAC5_DeadInProgressWithCommitsProposesDeltaRecovery`, `TestAC5_HealthyActiveInProgressIsValid`.

**Package/Binary gate** (`cd src/go && go build ./... && go test ./...`):
```
BUILD OK
ok  	github.com/usurobor/cnos/src/go/internal/activate	0.038s
ok  	github.com/usurobor/cnos/src/go/internal/activation	0.057s
ok  	github.com/usurobor/cnos/src/go/internal/binupdate	0.017s
ok  	github.com/usurobor/cnos/src/go/internal/cell	0.017s
ok  	github.com/usurobor/cnos/src/go/internal/cli	0.011s
ok  	github.com/usurobor/cnos/src/go/internal/discover	0.014s
ok  	github.com/usurobor/cnos/src/go/internal/dispatch	0.006s
ok  	github.com/usurobor/cnos/src/go/internal/doctor	0.246s
ok  	github.com/usurobor/cnos/src/go/internal/hubinit	0.034s
ok  	github.com/usurobor/cnos/src/go/internal/hubsetup	0.009s
ok  	github.com/usurobor/cnos/src/go/internal/hubstatus	0.010s
ok  	github.com/usurobor/cnos/src/go/internal/pkg	0.003s
ok  	github.com/usurobor/cnos/src/go/internal/pkgbuild	0.033s
ok  	github.com/usurobor/cnos/src/go/internal/restore	0.027s
```
(the `internal/cli` package's own test suite, which touches `cmd_issues_fsm.go`'s dispatch wrapper, passes — no wrapper changes were made, per the pinned contract.)

**`cell_kind` non-enforcement:**
```
$ rg -n "cell_kind" src/packages/cnos.cds/skills/cds/fsm/transitions.json
(no output — exit 1, zero matches)
```
No new rule references `cell_kind`; the seam stays observed-only. `TestSeam_CellKindNotEnforced` (unmodified assertion, extended fixture list covering all 6 new fixtures) confirms the decision is byte-identical across `cell_kind` values for every fixture including the new ones.

**CI gates (I1/I2/I4/I5/I6, install-wake-golden, dispatch guards #516/#524):** not independently re-run in this sandbox (no GitHub Actions runner available locally); the golden/workflow re-render and the two `scripts/ci/*.sh` guard scripts (both invoked directly above, both exit 0) are the local proxies for `install-wake-golden` and the #516/#524 dispatch guards respectively. Branch CI on the pushed head commit is the authoritative signal β should check per the pre-review gate (§Review-readiness below) — this is recorded as a transient row, not claimed green from local evidence alone.

**Status: MET** (behavioral invariants), **local proxy evidence only for the substrate CI gates** — see §Review-readiness for the explicit CI-state disclosure.
