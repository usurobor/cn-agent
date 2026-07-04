# γ R0 scaffold — cnos#575

## Issue reference

- **Issue:** [usurobor/cnos#575](https://github.com/usurobor/cnos/issues/575) — "cds/fsm: route claim, hard-block, release-back-to-queue through the FSM (Phase 3 — Sub 2 of #583)"
- **Mode:** design-and-build (per issue header)
- **Parent:** #583 (master wave — mechanical dispatch-cell architecture)
- **Precondition:** #584 (Sub 1 — mechanism/cognition doctrine) — landed, CLOSED. Verified via `gh issue view 584 --json state` → `CLOSED`.
- **Protocol:** cds
- **Wake run id:** 28703190110 (opaque, recorded for correlation only)
- **Cycle branch:** `cycle/575`, created from `origin/main`

### Base-SHA discrepancy (recorded, not blocking)

δ handed γ a pinned main SHA of `a0756669a17825c86721e845427c85b1838f1c0d`. At branch-creation time, `origin/main` HEAD had advanced to `8f97b604689f0d5b39a2ed6b3541a34e759cab9a` (two commits ahead: `1004572` agent-admin heartbeat log entry, `8f97b60` board-map regen). `git diff --stat` between the two SHAs touches only `.cn-sigma/logs/20260704.md` and `docs/development/board/{board-data.json,index.html}` — no overlap with any surface this cycle touches (`transitions.json`, the Go FSM engine, `cds-dispatch/SKILL.md`, `dispatch-protocol/SKILL.md`, or CDD role skills). Per `gamma/SKILL.md` staleness-check discipline, this is named explicitly rather than silently absorbed. `cycle/575` was created from `origin/main` at `8f97b604689f0d5b39a2ed6b3541a34e759cab9a` (the branch-creation instruction says "from `origin/main`", not from the stale pin). No re-load of canonical skills was triggered because the delta is non-overlapping.

### Empirical anchors

- **cnos#389 / #391 / #392 / #393** — the implementation-contract-drift lineage. #389 shipped in the wrong language because γ's prompt didn't pin one; #391 shipped at the wrong package path / as a separate binary for the same reason; #392 is the first cycle δ pinned the 7-axis contract ad hoc and it worked; #393 codified the `## Implementation contract` block as mandatory. This scaffold pins all 7 axes explicitly below — no TBD rows — precisely because #575 is Go-native CDS tooling extending an existing `cn` subcommand, exactly the shape #391 got wrong.
- **cnos#568/#569/#574** — the FSM lineage this cycle extends. #568 built the read-only Phase-1 evaluator + declarative `transitions.json`; #569 gave the FSM `--apply` mutation authority for `in-progress → review`; #574 hardened the review guard from `any_true` to `all_true` on `[review_request_present, pr_exists, pr_has_commits]`. AC5 requires none of this regresses.
- **cnos#368** — the blind-requeue failure class AC3 explicitly protects against (never re-dispatch/requeue over existing matter without checking).
- **cnos#516** — the repair-context guard on `changes → todo` (a structurally similar "don't silently proceed without evidence" pattern to what AC2's hard-block guard needs).

## Surfaces α is expected to touch

| Surface | Path | What changes |
|---|---|---|
| Transition table (data) | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` | New rules under `state: "todo"` (claim → in-progress), a new hard-block rule reachable from `state: "in-progress"` (→ blocked), and a new release-back-to-queue rule under `state: "in-progress"` (→ todo) that fires even while `run_active` is true (mirrors the #569 Phase-2 review-request pattern — see Friction note 1). Possibly new entries in the `guards` doc-map (data-only annotations; the doc strings are illustrative, not enforced). |
| FSM engine (Go, guard registry) | `src/packages/cnos.issues/commands/issues-fsm/table.go` | `guardFuncs` map (currently 9 entries, lines 90–100) likely needs new predicate(s) — see Friction note 1 for exactly which are missing. |
| FSM engine (Go, fact model) | `src/packages/cnos.issues/commands/issues-fsm/snapshot.go` | `FactSnapshot` struct likely needs new fields to carry the new guards' underlying observations (mirror the `ReviewRequestPresent` / `REVIEW-REQUEST.yml` pattern at lines 55–57). |
| FSM engine (Go, live-fetch assembly) | `src/packages/cnos.issues/commands/issues-fsm/fetch.go` | The live-path fact assembly (mirrors `case "REVIEW-REQUEST.yml": snap.ReviewRequestPresent = true` at line ~119) needs equivalent wiring for any new marker file(s)/fact(s) AC1–AC3 introduce. |
| CLI entrypoint (no structural change expected) | `src/go/internal/cli/cmd_issues_fsm.go` | Thin dispatch wrapper only (`Run` delegates to `issuesfsm.Run`); no change expected unless a new flag is needed. Flag surface currently: `evaluate --issue N [--fixture path] [--table path] [--apply] [--repo owner/name]`. |
| Go tests | `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go` | New `TestAC1_*` / `TestAC2_*` / `TestAC3_*` test functions (issue's own AC numbering — do not collide with the file's existing `TestAC1`–`TestAC7` which number *this file's own* historical ACs, not #575's; suggest naming e.g. `TestAC575_1_ClaimRoutedThroughFSM`, `TestAC575_2_HardBlockRequiresEvidence`, `TestAC575_3_ReleaseBlockedOverMatter` to avoid symbol collision — see Friction note 4). |
| Go test fixtures | `src/packages/cnos.issues/commands/issues-fsm/testdata/*.json` | New fixture files, TDD fail-then-pass. Suggested names (α may rename, but must follow the existing `testdata/<state>-<condition>.json` convention): `todo-claimable.json`, `todo-competing-run.json`, `in-progress-block-with-evidence.json`, `in-progress-block-no-evidence.json`, `in-progress-release-no-matter.json`, `in-progress-release-with-matter.json`. |
| CDS dispatch wake doctrine | `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` | The "Lifecycle transitions" table (~lines 258–272) explicitly labels claim / hard-block / release-back-to-queue as "direct label write" — this prose must change to "requests via `cn issues fsm evaluate --issue {N} --apply`" for those three, mirroring the existing β-converge row's phrasing. The frontmatter description (line 3) also says "...does not apply labels except for the four named lifecycle transitions..." — check whether "applies labels" language needs a request-vs-write distinction too. |
| Dispatch protocol doctrine | `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` | The claim sequence's steps 5–6 (~lines 130–132, `cds-dispatch/SKILL.md`) and the generic dispatch-protocol claim-sequence description (~lines 60–66, `dispatch-protocol/SKILL.md`: "4. **Claim** — ... remove `status:todo`, add `status:in-progress`") currently describe a direct label write. Per the issue's explicit scope this file is in-scope for AC4 — but it is **protocol-agnostic** (shared by cds/cdr/future cdw), while the FSM (`cn issues fsm`) is a cds-only tool today. See Friction note 2 — this is a real design question for α/δ, not one γ resolves here. |
| (possibly) δ doctrine | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.6 | The return-tokens table's "status:blocked + reason" and "claim release (race)" rows currently describe direct label transitions, not FSM requests. Not named in the issue's explicit AC4 surface list (`transitions.json` + dispatch doctrine), but it is the same "direct write" prose pattern. See Friction note 3 — flagged, not decided. |

## Per-AC oracle list

### AC1 — claim routed through the FSM

**Invariant:** `todo → in-progress` is applied by the FSM on passing guards (no competing active run; dispatchable contract), not a direct wake write.

**Oracle commands:**
- `git show cycle/575:src/packages/cnos.cds/skills/cds/fsm/transitions.json | jq '.transitions[] | select(.state=="todo")'` — must show a new rule with `target_state: "in-progress"` and named guard(s), not just the existing unconditional `valid`/`none` rule.
- `cd src/packages/cnos.issues/commands/issues-fsm && go test ./... -run TestAC575_1 -v` — new fixture-driven test(s) pass.
- `cn issues fsm evaluate --issue <N> --fixture testdata/todo-claimable.json --table ../../../cnos.cds/skills/cds/fsm/transitions.json` (or the Go-test equivalent via `Evaluate()`) — decision block shows `outcome: proposed`, `enabled_transition: todo -> in-progress`.
- `rg -n "direct label write" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` on the claim row — must no longer say "direct label write" for claim.

**Negative case:** a fixture with a competing active run (`testdata/todo-competing-run.json`) must evaluate to `outcome: valid` / `action: none` (or `blocked`, per α's design — but NOT `proposed` with `target_state: in-progress`). `TestAC575_1_ClaimBlockedOverCompetingRun` (or equivalent name) asserts this.

**Surface:** `transitions.json`, `table.go`/`snapshot.go`/`fetch.go` (if a new guard is needed — see Friction note 1), `cds-dispatch/SKILL.md`, new tests + fixtures.

### AC2 — hard-block routed through the FSM

**Invariant:** transition to `status:blocked` is FSM-applied only with explicit STOP/escalation evidence present.

**Oracle commands:**
- `git show cycle/575:src/packages/cnos.cds/skills/cds/fsm/transitions.json | jq '.transitions[] | select(.state=="in-progress") | .rules[] | select(.target_state=="blocked" or .outcome=="proposed")'` — a new rule proposing `status:blocked` exists, gated on a named evidence guard (not just the existing "dead run" fallback rules, which propose `todo` or delta-recovery, never `blocked`).
- `go test ./... -run TestAC575_2 -v` in `cnos.issues/commands/issues-fsm`.
- Fixture `testdata/in-progress-block-with-evidence.json` → `outcome: proposed`, `target_state: blocked`.
- Fixture `testdata/in-progress-block-no-evidence.json` → NOT `proposed`/`blocked` (either `valid`/`none` or an explicit `blocked`-with-missing-evidence outcome that does NOT apply the transition — α picks; must not silently apply).
- `rg -n "direct label write" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` on the hard-block row — must be gone.

**Negative case:** hard-block requested without evidence is refused — `TestAC575_2_HardBlockRefusedWithoutEvidence` (or equivalent).

**Surface:** `transitions.json`, Go FSM engine (new guard + FactSnapshot field — see Friction note 1), `cds-dispatch/SKILL.md`, tests.

### AC3 — release-back-to-queue routed through the FSM

**Invariant:** `in-progress → todo` (pre-work claim release) is FSM-applied only when no matter has been produced since claim; does NOT collide with the existing dead-run reconciliation rules (the last two rules already in `transitions.json`'s `in-progress` block, which handle *dead-run* reconciliation, not *live-run pre-work release*).

**Oracle commands:**
- `git show cycle/575:src/packages/cnos.cds/skills/cds/fsm/transitions.json | jq '.transitions[] | select(.state=="in-progress") | .rules'` — a NEW rule distinct from the existing "dead run with no matter" fallback rule (that rule requires `all_false: [run_active]`; the new release rule must fire even while `run_active` is true, analogous to how the #569 review-request rule fires despite `run_active` — see the `_doc_phase2` comment in `transitions.json` and Friction note 1).
- `go test ./... -run TestAC575_3 -v`.
- Fixture `testdata/in-progress-release-no-matter.json` (run_active=true, no commits/PR, explicit release-intent evidence present) → `outcome: proposed`, `target_state: todo`.
- Fixture `testdata/in-progress-release-with-matter.json` (run_active=true, release-intent evidence present, but `branch_has_commits` or `pr_exists` true) → NOT proposed as a requeue; routes to delta-recovery instead (reuse `propose_delta_recovery` action, or an equivalent explicit refusal) — this is the #368 blind-requeue protection.
- `rg -n "direct label write" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` on the release-back-to-queue row — must be gone.

**Negative case:** a release attempt with matter present is refused / routed to delta-recovery, never blind-requeued. `TestAC575_3_ReleaseBlockedOverMatter`.

**Surface:** `transitions.json` (new rule, positioned so it does not shadow or get shadowed by the existing dead-run rules — rule order matters, first-match-wins per `table.go`'s `Evaluate`), Go FSM engine, tests.

### AC4 — workers/wake request, don't write

**Invariant:** the CDS dispatch prompt/protocol instructs the wake to request all four transitions (claim, hard-block, release-back-to-queue, plus the existing review) through the FSM; no surviving prose describes a direct status write for these.

**Oracle commands:**
- `rg -n "direct label write" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` → zero matches (currently 3 matches at the claim/hard-block/release-back-to-queue rows of the "Lifecycle transitions" table, approx. lines 262–265; confirmed present via `rg -n "direct label write" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` at scaffold time).
- `rg -n "wake writes the label directly" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` (line ~258 prose: "For the claim, hard-block, and release-back-to-queue events, the wake writes the label directly") → zero matches after the fix; replaced with "requests via `cn issues fsm evaluate --issue {N} --apply`" language mirroring the β-converge row.
- `rg -n "remove.?label.*status:todo|add.?label.*status:in-progress" src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` (the claim-sequence steps 5–6, ~lines 130–132 of `cds-dispatch/SKILL.md`, and the generic claim-sequence description ~lines 60–66 of `dispatch-protocol/SKILL.md`) — check whether these need updating too; see Friction note 2 for the scope question (protocol-agnostic file, cds-only FSM).
- Manual read-through of both files' "Lifecycle transitions" / "Claim operation" sections confirms all four named transitions (claim, hard-block, release-back-to-queue, review) use FSM-request phrasing.

**Negative case:** any direct-write prose remains for these three transitions — caught by the `rg` zero-match oracle above.

**Surface:** `cds-dispatch/SKILL.md`, `dispatch-protocol/SKILL.md` (doctrine-only diff, no runtime code change on this AC).

### AC5 — invariants preserved; gates green

**Invariant:** the #574 review-guard tightening + `run_active` non-gating hold; `cell_kind` observed-only (`TestSeam_CellKindNotEnforced` unmodified); dispatch guards (#516/#524) + I1/I2/I4/I5/I6 + install-wake-golden + Go + Package + Binary all green.

**Oracle commands:**
- `cd src/packages/cnos.issues/commands/issues-fsm && go test ./... -v` — full existing suite green, in particular:
  - `TestAC574_ReviewRequestAloneNoLongerValid`, `TestAC574_ReviewPartialEvidenceBlocked`, `TestAC574_ReviewWithPRStillValid`, `TestAC574_InProgressBranchOnlyNoLongerProposesReview`, `TestAC574_InProgressWithMatterStillProposesReview`, `TestAC574_InProgressRunActiveNonGatingPreserved` — all must pass unmodified (no fixture used by these renamed, deleted, or altered in shape).
  - `TestSeam_CellKindNotEnforced` and `TestSeam_CellKindDefaultedWhenAbsent` — must pass unmodified; if a new guard is added, it MUST NOT reference `cell_kind` (that would flip the seam from observed-only to enforced — a scope violation per the issue's explicit "cell_kind enforcement… out of scope").
  - `TestAC4_DeadInProgressNoMatterProposesRequeue`, `TestAC5_DeadInProgressWithCommitsProposesDeltaRecovery`, `TestAC5_HealthyActiveInProgressIsValid` — must pass unmodified; the new AC3 release rule must be positioned in `transitions.json`'s `in-progress` rule list so it does not shadow (fire before) these three existing rules for their existing fixtures. Concretely: if the new release rule doesn't require `review_request_present`-equivalent explicit release-intent evidence, it MUST be positioned so the existing fixtures (`in-progress-dead-no-matter.json`, `in-progress-dead-with-commits.json`, `in-progress-active.json`) still hit their original rules first (first-match-wins per `table.go` `Evaluate`).
- `cd src/go && go build ./... && go test ./...` — full Go module green (Package/Binary gates).
- `gh workflow list` / CI run status on the pushed branch for I1/I2/I4/I5/I6, install-wake-golden — all green (checked at β/converge time, not at R0).
- `rg -n "cell_kind" src/packages/cnos.cds/skills/cds/fsm/transitions.json` → zero matches (no rule references `cell_kind`; the seam stays observation-only, consistent with the existing table).

**Negative case:** any of the above regress.

**Surface:** CI; the full existing FSM test suite.

## Source-of-truth table

| Claim / surface | Canonical source | Status |
|---|---|---|
| Transition table (data, single source of truth for CDS status transitions) | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` | Shipped (Phase 1/2, hardened #574); this cycle adds rules, does not restructure |
| Generic FSM evaluator engine (guard registry, rule matching, `Evaluate()`) | `src/packages/cnos.issues/commands/issues-fsm/table.go` | Shipped; guard registry (`guardFuncs`, 9 entries) is the extension point for new guards |
| Fact model (`FactSnapshot`) | `src/packages/cnos.issues/commands/issues-fsm/snapshot.go` | Shipped; extension point for new observed facts |
| Live fact assembly (GitHub API + git) | `src/packages/cnos.issues/commands/issues-fsm/fetch.go` | Shipped; extension point for wiring new facts to live observation |
| CLI entrypoint | `src/go/internal/cli/cmd_issues_fsm.go` | Shipped; thin dispatch wrapper, `evaluate [--apply]` is the only sub-verb |
| Existing FSM test suite + fixture format | `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go` + `testdata/*.json` | Shipped; TDD fail-then-pass convention this cycle must follow |
| CDS dispatch wake doctrine (claim/hard-block/release-back-to-queue as direct writes — the prose AC4 corrects) | `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` §"Lifecycle transitions" | Current, pre-#575 (says "direct label write" for 3 of 4 transitions) |
| Generic dispatch-protocol claim sequence | `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` §"Claim operation" | Current; protocol-agnostic (cds/cdr/cdw); see Friction note 2 |
| δ wake-invoked return-token contract | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.6 | Current; not in the issue's explicit AC4 surface list, but carries analogous direct-write prose — see Friction note 3 |
| Mechanism/cognition boundary doctrine (this cycle builds against it) | landed via #584 (doctrine only, no runtime change) | Landed, CLOSED |
| Master wave context | #583 | OPEN, tracks this sub |

## The α dispatch prompt

```text
You are α. Project: cnos.

Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.

Issue: gh issue view 575 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/575

Tier 3 skills:
- src/packages/cnos.core/skills/write/SKILL.md
- src/packages/cnos.cds/skills/cds/CDS.md
- src/packages/cnos.cdd/skills/cdd/issue/SKILL.md
- src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md
- src/packages/cnos.core/skills/eng/SKILL.md
- src/packages/cnos.core/skills/eng/ship/SKILL.md
- src/packages/cnos.core/skills/eng/test/SKILL.md
- src/packages/cnos.core/skills/go/SKILL.md

Scaffold: .cdd/unreleased/575/gamma-scaffold.md (this file — read it in full before coding;
it carries the per-AC oracle list, the surfaces you are expected to touch, the source-of-truth
table, and named friction notes you must resolve as part of implementation, not defer).

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Go — cnos's CDS FSM tooling is Go-native (the engine at `src/packages/cnos.issues/commands/issues-fsm/` and the `cn` CLI kernel at `src/go/`); no new runtime introduced by this cycle. |
| CLI integration target | `cn` subcommand — extends the existing `cn issues fsm evaluate [--apply]` sub-verb family (`src/go/internal/cli/cmd_issues_fsm.go`); this cycle does NOT add a new sub-verb or a new binary. |
| Package scoping | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` (data, new rules) + `src/packages/cnos.issues/commands/issues-fsm/{table.go,snapshot.go,fetch.go,issuesfsm_test.go,testdata/*.json}` (Go engine, only if new guard predicates are required — see Friction note 1) + `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` + `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` (doctrine prose). Do NOT touch `src/go/internal/cli/cmd_issues_fsm.go` unless a new CLI flag is genuinely required (none is expected — `--apply` already exists and covers all four transitions uniformly). |
| Existing-binary disposition | Extend `cn` — no new binary. The `issues-fsm` package is Go-source co-located under `src/packages/cnos.issues/commands/issues-fsm/` per the existing cnos#568/#392 precedent (package-command exec-dispatch does NOT apply here; this is compiled-in-kernel Go source). |
| Runtime dependencies | None beyond the existing Go toolchain and the already-vendored GitHub REST client used by `fetch.go`'s live path. No new external dependency. |
| JSON/wire contract preservation | `transitions.json`'s schema (`states`, `guards`, `transitions[].{state,trigger,rules[]}`, `rules[].{all_true,any_true,all_false,outcome,action,target_state,repair_pass,reason,evidence_guards}`) MUST remain backward-compatible. Do NOT restructure existing keys; ADD new rules/guards additively. The existing `review` transition rule (state `"review"`) and its two rules (valid-with-full-evidence / blocked-with-partial-evidence) MUST NOT change in shape or outcome for any existing fixture. `Decision`'s exported field set (`decision.go`) is additive-only per its own doc comments (`ApplyAttempted`/`Applied` pattern) — follow that precedent if a new field is needed; do not remove or retype an existing field. |
| Backward-compat invariant | `run_active` stays non-gating for the review-request path (#569) — do not add a guard that makes `run_active` block the worker's own in-flight review/claim/release requests (see the `_doc_phase2` comment in `transitions.json` explaining why gating on `run_active` would make the worker's own-run request unusable). The #574 review-guard tightening (`all_true` on `[review_request_present, pr_exists, pr_has_commits]` for `in-progress → review` and for `review` state validity) is UNCHANGED — do not touch the existing `review`-adjacent rules. `cell_kind` stays observed-only: no new rule's `all_true`/`any_true`/`all_false` may reference a `cell_kind`-derived guard (`TestSeam_CellKindNotEnforced` must still pass, unmodified, covering every fixture including your new ones — extend the test's fixture list to include your new fixtures if the test iterates a fixed list; do not create a parallel un-covered test). |
```

## The β dispatch prompt

```text
You are β. Project: cnos.

Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.

Issue: gh issue view 575 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/575

Scaffold: .cdd/unreleased/575/gamma-scaffold.md (this file). Walk the per-AC oracle list
below independently — do not trust α's self-coherence report as a substitute for running
the oracle commands yourself.

Per-AC oracle list to walk independently (full detail in the scaffold's "Per-AC oracle
list" section — reproduced here as the checklist you execute):

- AC1 (claim): confirm a new `todo`-state rule in `transitions.json` proposes
  `in-progress` gated on named guards; run the new AC1 test(s); confirm a competing-run
  fixture does NOT propose the claim transition; confirm `cds-dispatch/SKILL.md`'s claim
  row no longer says "direct label write".
- AC2 (hard-block): confirm a new `in-progress`-state rule proposes `blocked` gated on
  explicit evidence; run the new AC2 test(s); confirm a no-evidence fixture is refused;
  confirm `cds-dispatch/SKILL.md`'s hard-block row no longer says "direct label write".
- AC3 (release-back-to-queue): confirm the new release rule is DISTINCT from the existing
  dead-run reconciliation rules (check rule order in `transitions.json`'s `in-progress`
  block — first-match-wins, so verify the new rule doesn't shadow or get shadowed
  incorrectly); confirm a matter-present fixture is refused / routed to delta-recovery,
  never blind-requeued (the #368 protection); confirm `cds-dispatch/SKILL.md`'s
  release-back-to-queue row no longer says "direct label write".
- AC4 (prose): run `rg -n "direct label write" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`
  and confirm zero matches; run `rg -n "wake writes the label directly"` on the same file;
  read through `dispatch-protocol/SKILL.md`'s claim-sequence section and confirm α's
  resolution of Friction note 2 is coherent (either the file is updated with matching
  FSM-request language, or a documented reason is given for why it stays protocol-agnostic
  and out of scope for the label-writing language change).
- AC5 (regression): run the FULL existing `cnos.issues/commands/issues-fsm` test suite
  (`go test ./... -v`) and confirm every pre-existing `TestAC574_*`, `TestSeam_CellKind*`,
  `TestAC4_*`, `TestAC5_*` test still passes UNMODIFIED in assertion (renaming for
  namespace-collision reasons is acceptable per Friction note 4, but the assertions and
  fixtures they test against must not have changed meaning); run `go build ./...` +
  `go test ./...` from `src/go`; confirm `rg -n "cell_kind" transitions.json` has zero
  matches; confirm CI gates (I1/I2/I4/I5/I6, install-wake-golden, Go, Package, Binary)
  are green on the pushed branch.

## Implementation-contract verification (β Rule 7)

Confirm the diff matches every pinned axis in the α prompt's `## Implementation contract`
table above: Go only, no new binary, extends `cn issues fsm evaluate --apply`, package
scoping confined to `cnos.cds/skills/cds/fsm/`, `cnos.issues/commands/issues-fsm/`,
`cds-dispatch/SKILL.md`, `dispatch-protocol/SKILL.md` (plus test/fixture files under the
Go package), no new runtime dependency, `transitions.json` schema additive-only,
`run_active` non-gating + #574 guard + `cell_kind` observed-only all preserved. Any
diff drifting from this (new binary, new language, new CLI sub-verb, restructured JSON
schema, a `cell_kind`-consuming guard) is a REQUEST CHANGES, severity D,
classification `implementation-contract` — regardless of whether the behavioral ACs pass.

## Scope guardrails (restated from the issue — do not expand)

- Do NOT start Sub 3/4 of #583 (mechanical PR-open / checkpoint / recovery runtime).
- Do NOT implement the mechanical PR-open-on-first-matter extraction (#583's AC4) — that
  is a separate sub.
- Do NOT touch `cell_kind` enforcement — it stays observed-only.
  `TestSeam_CellKindNotEnforced` must remain unmodified in intent (extend its fixture
  list if it iterates one; do not weaken or remove its assertion).
- Do NOT introduce new status labels or taxonomy. The four transitions this cycle covers
  (claim, hard-block, release-back-to-queue, review) use the EXISTING label set
  (`status:todo`, `status:in-progress`, `status:blocked`, `status:review`) — no new
  `status:*` value.
- Do NOT build a Demo 0.
- Do NOT touch #585 / #586 (not cited in the issue body directly, but named in the
  dispatch task as explicitly out of bounds — if either surfaces as tempting adjacent
  work, defer it).
- Do NOT restructure `transitions.json`'s existing `review` / `changes` state blocks —
  additive-only for `todo` and `in-progress`.

## Friction notes (for α to resolve — γ does not decide these)

**1. The three new guards are NOT already expressible in the existing guard registry.**
`table.go`'s `guardFuncs` (9 entries: `run_active`, `branch_exists`,
`branch_has_commits`, `pr_exists`, `pr_has_commits`, `review_request_present`,
`repair_contract_present`, `cdd_artifacts_present`, `checks_passing`) does not have:
  - a guard distinguishing "a run OTHER than the current requesting run is active" from
    the existing `run_active` (which only asks "is *some* run queued/in_progress for this
    issue" — and per the `_doc_phase2` comment in `transitions.json`, the calling wake's
    OWN run is by definition `in_progress` at the moment it requests a transition
    synchronously, exactly the tension #569 solved for the review-request rule by NOT
    gating on `run_active`). AC1's guard is phrased "no *competing* active run" — this
    needs either (a) a new fact/guard that can distinguish self-run from other-run (would
    require FactSnapshot to carry a run-id and compare against the invoking run's own id,
    a nontrivial new fact-assembly requirement in `fetch.go`'s live path), or (b) a
    reframing where "no competing active run" is satisfied structurally by the claim
    sequence's pre-existing serialized-claim guard (per `dispatch-protocol/SKILL.md`
    §2.3 layers 1–2) and the FSM guard for claim is simply "dispatchable contract"
    (label well-formedness) without needing a *new* run-based guard at all, since the
    claim-time race is already handled upstream by the claim sequence itself. α must
    pick one framing and justify it in `self-coherence.md`; this scaffold does not
    prescribe which.
  - a guard for "explicit STOP/escalation evidence present" (AC2's hard-block guard).
    No existing fact or marker file captures this. The closest existing pattern is
    `ReviewRequestPresent` / `REVIEW-REQUEST.yml` (a marker file under
    `.cdd/unreleased/{issue}/`, wired in `fetch.go` at the `case "REVIEW-REQUEST.yml":`
    branch). α likely needs an analogous new marker (e.g. `BLOCK-REQUEST.yml`, or a
    named comment-scan) + a new `FactSnapshot` field + a new `guardFuncs` entry. The
    `cds-dispatch/SKILL.md` prose already describes the wake posting a "`STOP`/`BLOCKED`
    comment naming the missing evidence" in its deliverable-integrity preflight section
    (§"No-deliverable rule") — α should check whether that existing comment convention
    can be repurposed as the evidence source, or whether a dedicated marker file is
    cleaner (marker files are easier to observe deterministically than comment-text
    parsing).
  - a guard for "no matter produced since claim" that fires DURING an active run (AC3's
    release-back-to-queue guard). The existing "dead run with no matter" rule
    (`all_false: [run_active]`, `outcome: proposed`, `action: propose_status_todo`) is
    explicitly a DIFFERENT case — a reconciliation sweep finding an abandoned issue,
    not a synchronous pre-work release request from the wake's own active run. AC3
    needs a NEW rule, positioned before the existing dead-run rules in the `in-progress`
    rule list, gated on an explicit release-intent marker (mirroring the
    `review_request_present` pattern — e.g. a hypothetical `release_requested` guard +
    fact + marker file) `all_true`'d with `all_false: [branch_has_commits, pr_exists,
    pr_has_commits]` (no matter), and — symmetrically to the review-request rule's
    paired "blocked" rule — a second new rule for "release requested but matter
    present" that does NOT requeue (routes to `propose_delta_recovery` instead, per the
    issue's explicit AC3 negative case).

  **Net assessment for δ/α:** the FSM *engine* (rule-matching, guard composition,
  first-match-wins evaluation) is generic and needs NO changes — it already supports
  arbitrary new guards and rules via pure JSON-table edits, PROVIDED the underlying
  Go guard predicates exist. All three new guards this issue needs DO require new Go
  code (new `FactSnapshot` fields + new `guardFuncs` entries + new `fetch.go` live-path
  wiring), following the `review_request_present`/`REVIEW-REQUEST.yml` precedent
  exactly. This is a small, well-bounded Go change (not an engine redesign) but it is
  real code, not just JSON-table authoring — the issue's phrasing ("Define fact-model
  guards... Add the transition-table rules") already anticipates this; γ names it
  explicitly here so α does not under-scope the Go-side work.

**2. `dispatch-protocol/SKILL.md` is protocol-agnostic; the FSM is cds-only today.**
The issue's Scope explicitly names `cnos.core/skills/agent/dispatch-protocol/SKILL.md`
as a surface to update for AC4. But that skill's claim-sequence prose (steps 4–5,
"remove `status:todo`, add `status:in-progress`") is shared doctrine for every future
protocol-owned dispatch wake (cds/cdr/future cdw), while `cn issues fsm` is presently
cds-specific tooling (data lives under `cnos.cds/skills/cds/fsm/`; no `cnos.cdr`
or `cnos.cdw` equivalent exists yet). If α rewrites the generic claim-sequence
description to say "request via `cn issues fsm evaluate --apply`" unconditionally,
that would be false for any non-cds protocol wake today. α should resolve this by
either: (a) scoping the dispatch-protocol edit to a CDS-specific callout/note rather
than rewriting the generic claim-sequence steps themselves (e.g. "concrete protocol
packages that have an FSM — currently cds — MUST request lifecycle transitions through
it rather than writing labels directly; see `cds-dispatch/SKILL.md` for the reference
shape"), or (b) determining the generic claim sequence's "direct write" framing was
always meant as "the generic default absent a protocol-specific FSM" and adding that
qualifier. Either resolution is acceptable; α names the choice explicitly in
`self-coherence.md` — γ does not pre-decide this because it is a real design call
about how much of dispatch-protocol's generic prose a cds-only FSM should reshape.

**3. `delta/SKILL.md` §9.6 carries the same "direct write" pattern for hard-block and
claim-release, but is NOT named in the issue's explicit AC4 surface list.** The issue's
AC4 surface is stated as "dispatch doctrine" and its oracle names `cds-dispatch/SKILL.md`
+ `dispatch-protocol/SKILL.md` specifically. `delta/SKILL.md`'s return-tokens table
(§9.6) describes "status:blocked + reason" and "claim release (race)" as label
transitions without FSM-request phrasing (only the "status:review" row explicitly says
"a *request* for the transition... δ does NOT write the label directly"). This is the
identical prose pattern AC4 targets, one level up the doctrinal stack (δ's contract,
which `cds-dispatch/SKILL.md` implements). γ flags this as adjacent but does not add it
to α's required surface list since the issue text does not name it — α should read
§9.6 while working AC1–AC3 and decide (and record in `self-coherence.md`) whether
leaving it as-is is a genuine non-goal (the issue's "Out of scope" list does not
mention it either) or a gap worth a follow-up issue. Do not silently edit
`delta/SKILL.md` as an unscoped addition without naming the decision.

**4. Test-function-name collision risk.** `issuesfsm_test.go`'s existing tests are
named `TestAC1_*` through `TestAC7_*`, `TestAC569_*`, `TestAC574_*` — these numbers are
this FILE's own historical AC numbering (cnos#568/#569/#574's AC1–AC7, not #575's).
New tests for #575's AC1–AC3 MUST NOT be named bare `TestAC1_*`/`TestAC2_*`/`TestAC3_*`
(collision with existing #568 Phase-1 tests of the same name). Use an issue-qualified
prefix, e.g. `TestAC575_1_*`, `TestAC575_2_*`, `TestAC575_3_*`, consistent with how
`TestAC569_*` and `TestAC574_*` already qualify by issue number.

**5. Fixture naming convention.** Existing fixtures follow `<state>-<condition>.json`
(e.g. `in-progress-review-request-with-matter.json`, `changes-with-repair.json`). New
fixtures should follow the same pattern; γ suggested names above are illustrative, not
mandatory — α may rename for clarity as long as the AC-oracle `rg`/`go test` commands
in this scaffold are updated to match in `self-coherence.md` if names diverge.

**6. `all_false` rule check for the new `in-progress` release rule.** Because
`table.go`'s `Evaluate()` is first-match-wins over `tr.Rules` in file order, the new
AC3 release rule(s) MUST be positioned in `transitions.json`'s `in-progress` block
BEFORE the existing "healthy active" rule (`all_true: [run_active]` → valid/none) if
the release rule is meant to fire while `run_active` is true (per Friction note 1's
framing that it must fire synchronously during the requesting wake's own run) — placing
it after that rule would make it dead code, silently shadowed. α must verify rule
ordering empirically against ALL existing `in-progress` fixtures (not just the new
ones) via the full `go test ./... -v` run named in AC5's oracle.
