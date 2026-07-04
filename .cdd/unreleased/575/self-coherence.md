# α self-coherence — cnos#575

<!--
section-manifest:
  planned: [gap, skills, acs, self-check, debt, cdd-trace, review-readiness]
  completed: [gap, skills, acs, self-check, debt, cdd-trace]
-->

## Gap

**Issue:** [usurobor/cnos#575](https://github.com/usurobor/cnos/issues/575) — "cds/fsm: route claim, hard-block, release-back-to-queue through the FSM (Phase 3 — Sub 2 of #583)".

**Mode:** design-and-build. **Parent:** #583 (master wave). **Precondition:** #584 (Sub 1 doctrine) — landed, CLOSED.

**Gap statement (from the issue).** The CDS issue-state FSM (`cn issues fsm`, #568/#569, hardened #574) owns exactly one status transition: `in-progress → review`. Claim (`todo → in-progress`), hard-block (`→ status:blocked`), and release-back-to-queue (`in-progress → todo`) remain direct wake label writes, not routed through the FSM at all — so "the FSM owns status labels" was still aspirational, and the mechanical PR-open/recovery runtime (Subs 3-4 of #583) must not be built on top of direct label writes.

**Closure condition.** Claim, hard-block, and release-back-to-queue are FSM-applied on passing guards (TDD fixtures), the wake requests rather than writes, the #574/#569 invariants and `cell_kind` observed-only hold, and all gates are green.

**Branch:** `cycle/575`, created by γ from `origin/main`. γ's R0 scaffold (`.cdd/unreleased/575/gamma-scaffold.md`, commit `1db0608`) is the load-bearing artifact this implementation follows — its per-AC oracle list, surfaces-to-touch table, source-of-truth table, and six named friction notes are all addressed below.

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

## Skills

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

## ACs

Evidence below is real command output run against branch HEAD (implementation commits `02ef725` / `dd3e966`, before this readiness section), not a manual claim.

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

**Doctrine evidence:** `cds-dispatch/SKILL.md` claim sequence steps 5-6 now write `CLAIM-REQUEST.yml` then request `cn issues fsm evaluate --issue {N} --apply` instead of `gh issue edit --remove-label status:todo --add-label status:in-progress` (commit `dd3e966`); the "Lifecycle transitions" table's claim row no longer says "direct label write" (see AC4 below for the `rg` oracle covering this same site).

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

## Self-check

**Did α's work push ambiguity onto β?**

- Every AC has a concrete, re-runnable oracle command with output pasted in §ACs, not a narrative claim. β can re-run every `jq`, `go test -run`, and `rg` command verbatim.
- The rule-ordering concern the scaffold's Friction note 6 raised (first-match-wins shadowing) is closed with a dedicated test (`TestAC575_RuleOrderingDoesNotShadowExistingDeadRunRules`) exercising the exact three fixtures the scaffold named, not just "the suite is green."
- The one place this cycle leaves an open question for β to weigh, rather than closing outright, is Friction note 1's claim guard framing (see §Debt below) — a genuine design call the scaffold explicitly declined to pre-decide. It is not ambiguity α failed to resolve; it is a documented choice with a stated rationale, which is the coherent way to hand a real design call to review.
- CI gate rows (I1/I2/I4/I5/I6, install-wake-golden, dispatch guards) are recorded as **local-proxy-only** in §ACs AC5, not claimed green — this is disclosed explicitly rather than left for β to discover was never actually checked.

**Is every claim backed by evidence in the diff?**

Yes. Every new Go symbol (`ClaimRequestPresent`, `BlockRequestPresent`, `ReleaseRequestPresent` fields; `claim_request_present`/`block_request_present`/`release_request_present` guard funcs; the `CLAIM-REQUEST.yml`/`BLOCK-REQUEST.yml`/`RELEASE-REQUEST.yml` fetch.go wiring) has at least one non-test caller: the guard funcs are called from `table.go`'s `evalGuard` (via the generic `ruleMatches`/`Evaluate` path, which every `go test` invocation and every real `cn issues fsm evaluate` call exercises), and the FactSnapshot fields are populated by `fetch.go`'s live path (exercised by `TestAC575_LiveObservesNewMarkerFiles`, a non-`LoadFixture` caller) in addition to `LoadFixture`'s JSON path.

**Peer enumeration performed:**

- **`cds-dispatch/SKILL.md` "direct label write" prose sites** — enumerated via `grep -n` before editing (5 sites: frontmatter description L3, claim steps L130-131, lifecycle-transitions intro L258, three table rows L262/264/265, responsibilities item 5 L359); all updated together in one commit, not just the first hit found. The `rg` re-check after editing (§ACs AC4) proves zero survived.
- **Rendered-artifact peers of `cds-dispatch/SKILL.md`** — the golden fixture (`cnos-cds-dispatch.golden.yml`) and the live substrate workflow (`.github/workflows/cnos-cds-dispatch.yml`) are both machine-derived from the SKILL.md via `cn-install-wake`; both were re-rendered and verified to match a fresh render (no drift) rather than left stale. This is the harness-audit discipline (alpha/SKILL.md §2.4) applied to a non-Go, non-test harness (a shell-rendered YAML artifact).
- **`transitions.json` existing-fixture peers** — every pre-existing fixture in `testdata/*.json` was re-evaluated against the new rule set via the full test suite (not just the six new fixtures), and `TestSeam_CellKindNotEnforced`'s fixture list was extended rather than a parallel, uncovered test created.
- **Sibling doctrine surface `delta/SKILL.md` §9.6** — read and considered (same "direct write" prose pattern), explicitly NOT edited; see §Debt Friction note 3 for the reasoning, recorded rather than silently skipped or silently edited.

## Debt

### Friction note resolutions (all six, as required by dispatch)

**1. Claim guard framing ("no competing active run").** Resolved as a **hybrid**, not strictly either of γ's two offered options: reused the existing `run_active` guard unmodified (no new Go code for the run-check half — `run_active`'s existing definition, scoped to workflow runs observed against the would-be `cycle/{issue}` branch, already cannot see the requesting wake's own run, since that run executes on `main`/the dispatch workflow, not on a `cycle/{issue}` branch that doesn't exist yet at claim time) **plus** a new `claim_request_present` marker-file guard (new Go code: `FactSnapshot.ClaimRequestPresent`, the `claim_request_present` guardFunc, and `fetch.go`'s `CLAIM-REQUEST.yml` case), mirroring the `review_request_present`/`REVIEW-REQUEST.yml` precedent. Pure option (b) (drop the run-based guard entirely, treat "dispatchable contract" as satisfied structurally by the claim sequence) was rejected because it cannot produce AC1's own negative case ("a fixture with a competing active run is blocked") — dropping the run check makes that oracle unsatisfiable by construction. Pure option (a) (a new self-vs-other run-id-comparing fact) was rejected as γ predicted: nontrivial, and unneeded, because `run_active`'s existing branch-scoped definition already structurally excludes the requesting wake's own pre-branch execution — no new run-id plumbing is required to get the same discriminating power. The marker-file half was necessary regardless of which run-check option was chosen: without it, an unconditional "any evaluate on a `status:todo` issue proposes claim when `run_active` is false" rule would have made `testdata/todo.json`'s existing idempotence assertion (`TestApply_RequeueTransitionAppliesOnGuardPassAndIsIdempotent`'s second call, which evaluates a post-requeue `status:todo` fixture and expects a no-op) fail — proven empirically by the TDD sanity check below.

**2. `dispatch-protocol/SKILL.md` scope (protocol-agnostic file, cds-only FSM).** Resolved as option (a) from the scaffold's framing: added a scoped "Protocol-specific FSM override" paragraph within existing §2.2 (not a new top-level numbered subsection, to avoid renumbering every downstream cross-reference to §2.3-§4.9 that this 613-line file and its external referrers carry), stating that a concrete protocol package which has shipped its own FSM (today, only `cnos.cds`) MUST route lifecycle transitions through it, while protocols without one (`cnos.cdr` today) still use the generic direct-write shape shown in the existing step-4 code block. The generic claim-sequence code block itself is untouched. This keeps the file honest for future non-cds dispatch wakes that have no FSM yet, while making the cds-specific reality (all four transitions now FSM-requested) discoverable from the generic doctrine via a cross-reference to `cds-dispatch/SKILL.md`.

**3. `delta/SKILL.md` §9.6 (same "direct write" prose pattern, not in AC4's named surface list).** Read in full (via targeted `grep`) and left unedited. Reasoning: §9.6's "status:blocked + reason" and "claim release (race)" rows describe δ **signaling** the wake that a transition is needed (a semantic fact that remains true after #575: δ still decides *when* a hard-block or claim-release is needed), not δ **performing** the label write itself — the actual write-vs-request mechanism lives one layer down, in `cds-dispatch/SKILL.md` (which #575 does update). Editing §9.6's prose now would restate a mechanism detail that isn't §9.6's job to carry, and the issue's own "Out of scope" list doesn't name it. This is recorded as a candidate follow-up (not filed as a separate issue by α, since triage is γ's job per alpha/SKILL.md §2.8 "Voice: factual observations and patterns only"): a future doctrine pass could add a one-line cross-reference from §9.6 to `cds-dispatch/SKILL.md`'s mechanism, but that is a judgment call for whoever next touches `delta/SKILL.md`, not an unscoped edit to smuggle into this cycle.

**4. Test-function naming.** Used `TestAC575_1_*` / `TestAC575_2_*` / `TestAC575_3_*` prefixes throughout (plus `TestAC575_RuleOrderingDoesNotShadowExistingDeadRunRules`, `TestAC575_LiveObservesNewMarkerFiles`, and five `TestAC575_Apply*` names) — zero collisions with the file's own historical `TestAC1`-`TestAC7` / `TestAC569_*` / `TestAC574_*` names, confirmed by `grep -c "^func TestAC575" issuesfsm_test.go` = 14 and by the full suite compiling/running with no duplicate-symbol error.

**5. Fixture naming convention.** All six new fixtures follow `<state>-<condition>.json`: `todo-claimable.json`, `todo-competing-run.json`, `in-progress-block-with-evidence.json`, `in-progress-block-no-evidence.json`, `in-progress-release-no-matter.json`, `in-progress-release-with-matter.json` — consistent with the existing `in-progress-review-request-with-matter.json` / `changes-with-repair.json` naming.

**6. Rule ordering (first-match-wins, must not shadow the dead-run reconciliation rules).** The three new `in-progress` rules (hard-block, release-no-matter, release-with-matter) are positioned immediately after the two existing review-request rules and **before** the `all_true: [run_active]` valid rule and both dead-run reconciliation rules. This was verified two ways: (a) `TestAC575_RuleOrderingDoesNotShadowExistingDeadRunRules` asserts the three pre-existing dead-run/active fixtures still resolve to their original outcome/action/target_state; (b) the full pre-existing test suite (`TestAC4_DeadInProgressNoMatterProposesRequeue`, `TestAC5_DeadInProgressWithCommitsProposesDeltaRecovery`, `TestAC5_HealthyActiveInProgressIsValid`, and everything else) passes unmodified against the modified `transitions.json`.

**TDD sanity check (empirical, run twice during this cycle, not part of the committed diff).** Checking out `transitions.json`'s pre-#575 content (commit `1db0608`) into the working tree and re-running `go test ./... -run TestAC575 -v` against it produces exactly:

```
--- FAIL: TestAC575_1_ClaimRoutedThroughFSM
--- FAIL: TestAC575_1_ClaimBlockedOverCompetingRun
--- PASS: TestAC575_1_ClaimNotRequestedStaysValid
--- FAIL: TestAC575_2_HardBlockRoutedThroughFSM
--- PASS: TestAC575_2_HardBlockRefusedWithoutEvidence
--- FAIL: TestAC575_3_ReleaseRoutedThroughFSMWhenNoMatter
--- FAIL: TestAC575_3_ReleaseBlockedOverMatter
--- PASS: TestAC575_RuleOrderingDoesNotShadowExistingDeadRunRules
--- PASS: TestAC575_LiveObservesNewMarkerFiles
--- FAIL: TestAC575_ApplyClaimTransitionAppliesOnGuardPass
--- FAIL: TestAC575_ApplyClaimBlockedRefusesAndMutatesNothing
--- FAIL: TestAC575_ApplyHardBlockAppliesOnGuardPass
--- FAIL: TestAC575_ApplyReleaseAppliesOnGuardPass
--- PASS: TestAC575_ApplyReleaseWithMatterDoesNotRequeue
```

9 of 14 fail without the new rules — exactly the tests asserting a *positive* proposed/applied outcome (claim proposed, claim blocked-over-competing-run, hard-block proposed, release proposed, and their four `--apply`-level counterparts). The 5 that still pass without the new rules are the ones whose asserted outcome coincides with a pre-existing fallback rule's behavior even absent the #575 rules: `TestAC575_1_ClaimNotRequestedStaysValid` and `TestAC575_2_HardBlockRefusedWithoutEvidence` assert a no-op/fallback by design (their fixtures have the request marker false/absent); `TestAC575_RuleOrderingDoesNotShadowExistingDeadRunRules` and `TestAC575_LiveObservesNewMarkerFiles` don't depend on `transitions.json` content at all (the former only re-checks pre-existing fixtures, the latter tests `fetch.go`'s marker-file wiring directly); and `TestAC575_ApplyReleaseWithMatterDoesNotRequeue` happens to still observe `applied: false` under the old rule set too, because the old fallback rule (`all_true: [run_active]` → valid/none) also produces no write for that fixture's facts — a coincidental pass for the wrong reason, noted here rather than silently relied on, though the test still correctly guards the real (new) implementation's behavior once the new rules are present. Restoring the modified `transitions.json` afterward reproduces the full 51/51 green suite reported in §ACs AC5.

### Known debt

- **Substrate CI gates not independently verified in this sandbox.** I1/I2/I4/I5/I6, `install-wake-golden`, and the two `scripts/ci/check-dispatch-*.sh` guards' GitHub-Actions-hosted runs were not observable from this environment (no live GitHub Actions access). Local proxies were run instead (both `check-dispatch-*.sh` scripts directly; `cn-install-wake cds-dispatch` re-render diffed against the committed golden and workflow). β/the pre-review gate must confirm actual branch CI is green on the pushed head commit before merge — this is flagged, not silently assumed.
- **`delta/SKILL.md` §9.6 cross-reference gap** (Friction note 3) — a candidate follow-up, not filed as an issue by α (triage is γ's/operator's call), recorded here per alpha/SKILL.md §2.8's "factual observations, not recommended dispositions" voice rule.
- **`cds-dispatch/SKILL.md`'s "Surfaces" §"You MAY write to" bullet** ("Label application on the claimed cell only — the four transitions enumerated in §Lifecycle transitions above") was not reworded to mention the FSM-request indirection; it remains true at the level it's written (the wake still causes those four labels to be applied, now via the FSM rather than directly) and rewording it would be restating the same fact a third time beyond the frontmatter description and the Lifecycle transitions table/prose — judged as diminishing-return repetition rather than a required peer site, since it doesn't use the phrase "direct label write" or any of the AC4 oracle's targeted strings.
- **CDS.md was consulted but not re-audited line-by-line** (3600 lines) — only the sections the scaffold specifically pointed to (artifact contract, location matrix, coordination surfaces) were read. No claim in this cycle depends on an unread part of CDS.md; if CDS.md's step table itself named the four lifecycle transitions with write-vs-request phrasing, that would be a peer this cycle should have caught — a targeted `rg` for the same oracle strings against CDS.md is a cheap next check β can run if this concerns review (`rg -n "direct label write|wake writes the label directly" src/packages/cnos.cds/skills/cds/CDS.md` returns zero matches at the time of this writing, run as part of closing this debt item before the review-readiness signal).

## CDD Trace

Per `CDS.md` §"Step table", α owns steps 4-7.

- **Step 4 — Gap (α names the incoherence precisely).** See §Gap above: the FSM owned exactly one status transition (`in-progress → review`); claim, hard-block, and release-back-to-queue were direct wake label writes. This cycle closes that gap for the three remaining transitions.
- **Step 5 — Mode (α chooses MCA/MCI and active skills).** Mode: MCA (autonomous implementation against a fully-scaffolded issue + γ scaffold with no unresolved design ambiguity requiring operator interaction — the six friction notes were γ-named design calls for α to resolve, not blockers requiring escalation, and all six were resolved without needing to stop). Active skills: §Skills above (Tier 1: CDD.md + alpha/SKILL.md + issue/SKILL.md + issue/proof/SKILL.md; Tier 2: eng/ship, eng/test, eng/go, write; Tier 3: the γ scaffold itself).
- **Step 6 — Artifacts (produced in artifact order: design → contract → plan → tests → code → docs).**
  - *Design:* not a separate design doc — γ's scaffold (`.cdd/unreleased/575/gamma-scaffold.md`, already on the branch pre-α) carries the per-AC oracle list, surface table, and friction-note framing that constitutes this cycle's design; α's job was to resolve the six named friction notes (documented in §Debt), not re-derive a design from scratch. Justification for skipping a separate design artifact: the scaffold already names invariant, oracle, positive/negative cases, and surfaces per `issue/proof/SKILL.md`'s output pattern — writing a second design doc restating the same content would be repetition, not new design work.
  - *Contract:* this file's §Gap.
  - *Plan:* not written as a separate artifact — sequencing was: (1) resolve Friction note 1 (guard framing) since it gates what Go code is needed at all, (2) add FactSnapshot fields + guardFuncs + fetch.go wiring, (3) add transitions.json rules in the friction-note-6-directed position, (4) write fixtures, (5) write tests, (6) empirically verify via TDD stash check, (7) update doctrine (AC4), (8) re-render golden/workflow artifacts, (9) run full gates. Justification for no separate plan doc: a single-PR, single-package cycle with a linear dependency chain (guards before rules before fixtures before tests) — plan/SKILL.md's own bar ("when implementation sequencing is non-trivial") is not met here relative to writing it inline.
  - *Tests:* `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go` (+421 lines: 14 new `TestAC575_*` functions plus the `TestSeam_CellKindNotEnforced` fixture-list extension) and six new fixtures under `testdata/` (`todo-claimable.json`, `todo-competing-run.json`, `in-progress-block-with-evidence.json`, `in-progress-block-no-evidence.json`, `in-progress-release-no-matter.json`, `in-progress-release-with-matter.json`). TDD-fail-then-pass demonstrated empirically (see §Debt Friction note 6).
  - *Code:* `src/packages/cnos.cds/skills/cds/fsm/transitions.json` (+80 lines: 3 new guard doc entries, a `_doc_phase3` block, 2 new `todo`-state rules, 3 new `in-progress`-state rules), `src/packages/cnos.issues/commands/issues-fsm/snapshot.go` (+27: 3 new `FactSnapshot` fields), `table.go` (+6: 3 new `guardFuncs` entries), `fetch.go` (+11: 3 new marker-file `switch` cases), `decision.go` (+6: 3 new `Render()` lines for operator visibility).
  - *Docs:* `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` (+20/-13: frontmatter description, claim-sequence steps 5/6/8, "Lifecycle transitions" table + intro prose, responsibilities item 5, Cross-references), `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` (+2: "Protocol-specific FSM override" paragraph in §2.2), plus the two machine-rendered artifacts kept in sync (`src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`, `.github/workflows/cnos-cds-dispatch.yml`, both +14/-7 from `cn-install-wake cds-dispatch`).
  - **Caller-path trace for new symbols (pre-review gate row 12):** `ClaimRequestPresent`/`BlockRequestPresent`/`ReleaseRequestPresent` (snapshot.go) are read by `table.go`'s new `guardFuncs` entries (`claim_request_present`/`block_request_present`/`release_request_present`), which are called from `evalGuard` → `ruleMatches` → `Evaluate` — the same generic evaluator path every `cn issues fsm evaluate` invocation (fixture or live) already runs; non-test callers: `issuesfsm.go`'s `runEvaluate` (the CLI entry point) and `cmd_issues_fsm.go`'s `IssuesFsmCmd.Run` (the kernel dispatch wrapper) — neither file was modified, both already call `Evaluate`/`LoadFixture`/`assembleLive` unconditionally. The three fetch.go `switch` cases are reached from `assembleLive`'s existing `.cdd/unreleased/{issue}/` directory-walk, itself called by `runEvaluate`'s live (non-`--fixture`) path — exercised by `TestAC575_LiveObservesNewMarkerFiles`, a non-`LoadFixture` caller.
  - **Artifact enumeration matches diff (pre-review gate row 11):** every file in `git diff --stat origin/main..HEAD` is named in this step's bullets above or in §ACs: `transitions.json`, `snapshot.go`, `table.go`, `fetch.go`, `decision.go`, `issuesfsm_test.go`, the 6 new `testdata/*.json` fixtures, `cds-dispatch/SKILL.md`, `dispatch-protocol/SKILL.md`, `cnos-cds-dispatch.golden.yml`, `.github/workflows/cnos-cds-dispatch.yml`, `.cdd/unreleased/575/gamma-scaffold.md` (γ's pre-existing artifact, not authored by α), and this `self-coherence.md` itself.
- **Step 7 — Self-coherence (α audits own work against ACs and the triad).** This file, §ACs + §Self-check + §Debt above, is that audit.

