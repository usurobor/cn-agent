<!-- section-manifest
completed: [Gap, Skills, ACs]
remaining: [Self-check, Debt, CDD Trace, guard-inventory, Review-readiness]
-->

# self-coherence — cycle #600

## §Gap

**Issue:** [cnos#600](https://github.com/usurobor/cnos/issues/600) — "cds/runtime: consolidate strand-era guards and scaffold now that the mechanical lifecycle is closed (Sub D of #583)".

**Mode:** `design-and-build` (per the issue's own `**Mode:**` field).

**cell_kind:** `audit` — recorded per γ's scaffold; observed-only, does not gate any FSM transition (see AC4 below).

**Parent:** #583 (Sub D). **Preconditions (all merged, per issue body):** #584, #575, #591, #593.

**Base SHA at branch creation:** `1c5dd993fcd25b5bdca3843555d240895e3212ee`. `origin/main` advanced twice while this cycle ran (to `65a55ae5`, then to `eb94445b`, both trivial `activate+attach` heartbeat-log commits touching only `.cn-sigma/logs/*` — no overlap with this cycle's surfaces). The cycle branch was rebased onto `eb94445b` and force-pushed; see §CDD Trace / §Review-readiness for the exact SHAs.

**The gap this cycle closes:** after Sub B (`cn cell finalize`, #591) and Sub C (`cn issues fsm scan`, #593) made checkpoint/PR-open and dead-run recovery mechanical, a layer of CI guards and scaffold accreted *while the strand problem was still unsolved* (cnos#516 repair re-entry preflight, cnos#524 W4 closeout integrity, the cnos#574 review-guard tightening, the cnos#568/#570 `cell_kind` seam). This cycle is an **audit-first consolidation**: classify every strand-era guard KEEP/REMOVE/FOLD/NARROW, cite a live replacement for anything removed/folded, and leave no lifecycle invariant weakened. It is explicitly **not** a feature cycle — no new FSM behavior, no new labels, no `cell_kind` enforcement, no Demo 0.

**Dispatch intake:** `git fetch origin cycle/600 && git switch cycle/600` — branch pre-existed (γ created it), no branch was created by α, per the constraint. γ's `gamma-scaffold.md` (already on the branch at dispatch time) was read in full, along with the live issue body via `gh issue view 600`, before any classification work began. γ's guard-surface inventory table was treated as a starting hypothesis to verify, not settled fact, per the dispatch's explicit instruction.

## §Skills

**Tier 1 (loaded in full):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical kernel/role-contract doctrine.
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — the α role surface governing this cycle: dispatch intake (§2.1), produce-in-artifact-order (§2.2), peer enumeration (§2.3), harness audit for schema-bearing changes (§2.4), incremental self-coherence discipline (§2.5), pre-review gate (§2.6, including the 15-row checklist), request-review (§2.7), the implementation-contract rule (§3.6).

**Tier 2/3, per the scaffold's "Skills to load" list and the actual diff surface:**
- `eng/process-economics` — the "guard-accretion vs boundary-fix" lens the issue names as its own governing frame: is a guard still doing boundary-fix work, or has the boundary moved (to the FSM/finalizer/reconciler) leaving the guard as pure accretion? Applied directly to the classification of both CI scripts.
- `go` — Tier 3, applied when reading/verifying `issuesfsm.go`, `scan.go`, `table.go`, `issuesfsm_test.go`, `scan_test.go`, and running `go test`/`go vet`.
- `eng/test` — applied to the AC2/AC3 oracle discipline (cite the specific test function, run it, prove green — not "the tests should cover this").
- Bash conventions (no dedicated Tier-3 `eng/bash` skill was loaded as a distinct file; the two guard scripts were edited following the existing script's own house style — `set -euo pipefail`, a `need()` helper, `::error::` GitHub Actions annotations — rather than introducing a new style).

**Not loaded:** `design/SKILL.md`, `plan/SKILL.md` as full artifacts — see §ACs / design-and-plan-not-required justification below. β/γ role skills were not loaded per α's Tier-1 restriction (`alpha/SKILL.md` §2.1 step 6: "do not load β or γ role skills").

**Active-skill-as-generation-constraint note (rule 3.1):** `eng/process-economics`'s guard-accretion lens directly shaped the classification method used throughout §ACs below — every KEEP/NARROW/FOLD/REMOVE decision below is argued in terms of "has the boundary that used to require this guard moved to a mechanical enforcer," not "does this guard look old."

**Design/plan artifact status:** marked **not required**, per `alpha/SKILL.md` §2.2's "silent omission is incomplete" bar — justification: this is a scoped consolidation cycle whose design was already done in γ's scaffold (guard-surface inventory, AC oracle table, implementation contract, STOP conditions all pre-declared); α's job is classification + narrow editing against that pre-existing design, not new design. No genuine design question surfaced during the audit that the scaffold didn't anticipate (see §Debt for the one adjacent-but-not-design-blocking observation about `install-wake-golden` redundancy, which was resolved by NOT acting on it, not by a design decision).

## §ACs

Base for all "current state" reads below: branch HEAD at time of audit, rebased onto `origin/main` at `eb94445b` (see §CDD Trace for exact SHAs). All test runs below were re-executed by α directly (not inherited from γ's scaffold or sub-agent paraphrase) via `go test ./src/packages/cnos.issues/commands/issues-fsm/... ./src/go/internal/cell/...`.

### AC1 — Guard inventory exists; every strand-era guard classified

Final classification table lives in `.cdd/unreleased/600/guard-inventory.md` (written as its own file per the scaffold's "your choice, name it clearly" option). Every row carries exactly one of KEEP / REMOVE / FOLD / NARROW — no blank/TBD cells. Summary (full detail + evidence in the dedicated file):

| Surface | Final classification |
|---|---|
| `scripts/ci/check-dispatch-repair-preflight.sh` (#516) | **KEEP** (was γ's NARROW hypothesis — corrected: already narrow by construction, no functional change; header strengthened per AC5) |
| `scripts/ci/check-dispatch-closeout-integrity.sh` (#524) — presence-of-contract half | **KEEP** (header strengthened per AC5) |
| `scripts/ci/check-dispatch-closeout-integrity.sh` (#524) — `--self-test` / `closeout_violation()` half | **FOLD** (removed; replacement cited and proven green) |
| cnos#574 review-guard tightenings in `transitions.json` | **KEEP** (this is itself the replacement other rows cite; no dead/superseded variant found alongside it) |
| `cell_kind` seam (`snapshot.go`, `fetch.go`, `decision.go`, `CELL-KINDS.md`) | **KEEP** (observed-only, non-enforced, confirmed via fresh full read — no misleading language found; no edit made) |
| Prompt-presence doctrine (`dispatch-protocol/SKILL.md` §2.8/§2.9/§4.8/§4.9) | **KEEP**, with §2.9/§4.9/D12 edited to remove now-stale claims about the folded bash self-test |
| TODO/scaffold comments (#568/#570/#575/#591) across `issuesfsm.go`, `scan.go`, `table.go` | **No dead scaffold found** (negative result after full-file reads — see below) |

### AC2 — No guard removed without a cited live replacement

The only REMOVE/FOLD row is the `--self-test` half of `check-dispatch-closeout-integrity.sh`. Replacement citation:

- `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go:100` `TestAC3_EmptyReviewBlocked` — status:review, no PR, no commits, no REVIEW-REQUEST.yml → `blocked`. Matches self-test case `(review, no, no) → violation`.
- `issuesfsm_test.go:500` `TestApply_EmptyReviewStateBlocked` — same case, exercised through the `--apply` CLI path (not just `Evaluate` directly), confirming the CLI-level enforcement too.
- `issuesfsm_test.go:587` `TestAC574_ReviewPartialEvidenceBlocked` — `pr_exists=true`, `pr_has_commits=false` → `blocked`. Matches self-test case `(review, yes, no) → violation`.
- `issuesfsm_test.go:602` `TestAC574_ReviewWithPRStillValid` — PR + commits → `valid`. Matches self-test case `(review, yes, yes) → ok`.
- The fourth self-test case (`in-progress, no, no → ok`) tested only that the bash predicate doesn't fire outside `status:review` — this is true by construction of `transitions.json`'s rule scoping (review-evidence rules live under the `review`/`in-progress→review` blocks only) and needs no dedicated Go test; it was never an independent behavioral claim, only a guard against the shell predicate's own scope.

Proof these are green, run directly by α (not inherited): `go test ./src/packages/cnos.issues/commands/issues-fsm/... ./src/go/internal/cell/...` →
```
ok  	github.com/usurobor/cnos/packages/cnos.issues/commands/issues-fsm	0.029s
ok  	github.com/usurobor/cnos/src/go/internal/cell	0.010s
```
(full verbose per-test output re-verified in §CDD Trace step 6). No other REMOVE/FOLD rows exist this cycle — AC2 is otherwise vacuously satisfied.

### AC3 — Strand fixtures still pass after consolidation

Ran (from repo root, using `go.work` to unify both modules): `go test ./src/packages/cnos.issues/commands/issues-fsm/... ./src/go/internal/cell/... -v 2>&1`. Result: **all packages `ok`, 0 failures.** Every named fixture/test from the scaffold's AC3 oracle table individually confirmed `--- PASS`:

| Invariant | Test | Result |
|---|---|---|
| empty-review | `TestAC3_EmptyReviewBlocked`, `TestApply_EmptyReviewStateBlocked` | PASS |
| repair re-entry | `TestAC6_ChangesWithoutRepairContextBlocked`, `TestAC6_ChangesWithRepairContextEnablesRepairPass` | PASS |
| matter-without-PR | `TestScan_DeadWithMatterNoPR_FinalizesAndComments` | PASS |
| dead-run-no-matter | `TestScan_DeadNoMatter_RequeuesToTodo`, `TestAC4_DeadInProgressNoMatterProposesRequeue` | PASS |
| dead-run-with-matter | `TestAC5_DeadInProgressWithCommitsProposesDeltaRecovery`, `TestScan_DeadWithPRNoReviewRequest_FinalizeNoOpNoComment` | PASS |
| review-request guarded transition | `TestAC569_InProgressReviewRequestWithMatterProposesReview`, `TestAC569_InProgressReviewRequestNoMatterBlocked`, `TestAC574_ReviewRequestAloneNoLongerValid` | PASS |
| cell_kind non-enforcement | `TestSeam_CellKindNotEnforced` | PASS |

This run was performed **after** the `--self-test` fold and the two script/doc edits landed — i.e. it is the post-consolidation state, not a pre-change baseline. No Go source or test file was modified this cycle (confirmed: `git diff origin/main..HEAD --stat` shows only `scripts/ci/*.sh` + `dispatch-protocol/SKILL.md` + `.cdd/unreleased/600/*` touched — no `.go` file in the diff), so AC3's fixtures were never at risk of being touched; this run proves they remain green with the guard-script consolidation applied.

### AC4 — `cell_kind` resolved

- **Retained observed-only:** confirmed via full read of `snapshot.go:96-125` (the `CellKind` field + struct + `normalizeCellKind` doc comments), `fetch.go:149-218` (`assembleLive`'s `cell_kind:` line parse, `cellKindLinePattern`, `parseCellKind`), `decision.go:92-96` (printed in the decision block). All three consistently state: observed, not consumed by any transition rule, default `"implementation"` when absent.
- **Not enforced:** `table.go` (read in full, 217 lines) has zero references to `CellKind` in any rule — confirmed by full read, not just grep.
- **`TestSeam_CellKindNotEnforced` intact and unweakened:** read in full (`issuesfsm_test.go:810-854`). It iterates **10 fixtures × 5 `cell_kind` values** (`implementation, issue_authoring, wave, cleanup, recovery`) and asserts `Outcome`/`TargetState`/`Action` are byte-identical to the baseline for every combination. Confirmed `--- PASS` in the current test run. **No edit was made to this test or any fixture it loads** — the "unweakened" requirement is satisfied by non-modification, verified.
- **No misleading TODO language:** re-read `CELL-KINDS.md`'s "Observation, not enforcement" section (lines 156-158) and the "Future typed field (deferred, Phase B)" paragraph (line 154) in full, plus `snapshot.go`/`fetch.go`/`decision.go`'s cell_kind-specific comments. **Judgment (resolving γ's friction note 3):** none of this qualifies as misleading/forgotten scaffold. Every passage names three things together — (1) observation-only, (2) not consumed by any rule, (3) enforcement deferred to a specific named future issue/phase (cnos#570's taxonomy plug-in, "a future, explicitly scoped FSM Phase 2 issue," "Phase B" CUE schema work). The two bare "not yet" phrasings found (`CELL-KINDS.md:123`, and the quoted issue text at `:152`) are self-qualifying in context (line 123 points forward in the same doc; line 152's "not yet" is inside a verbatim quotation of the original issue's own recommendation, not free-standing prose) and were judged NOT misleading. **No edit was made to `CELL-KINDS.md` or any cell_kind source file** — this is a negative-result resolution: AC4 is satisfied by confirming the existing state is already correct, not by producing a diff.
- **No new transition rule reads `CellKind`:** confirmed by the same full `table.go` read above.

### AC5 — Prompt-presence guards narrowed if kept, with explicit protects-statement

Both kept CI scripts now state explicitly, in their own header comments, what they protect and what they do not:

- `check-dispatch-repair-preflight.sh` header (new paragraph): "this script proves the PROMPT/protocol/golden/live-workflow text still states the repair-re-entry contract in words. It does NOT prove the FSM actually blocks a repair re-entry without repair context — that behavioral half is proven live by `transitions.json`'s `changes -> todo` rule ... and by `issuesfsm_test.go`'s `TestAC6_ChangesWithoutRepairContextBlocked` / `TestAC6_ChangesWithRepairContextEnablesRepairPass`."
- `check-dispatch-closeout-integrity.sh` header (rewritten): explicitly names the four Go tests that now carry the empty-review invariant, states the script "now only proves the PROMPT still says the words, not that the FSM behaves correctly," and cites the specific `transitions.json` guard (`review_request_present && pr_exists && pr_has_commits` all_true).

Both statements are in the script's own header comment, not left implicit, satisfying the scaffold's AC5 oracle ("the script's own header comment ... says this in so many words").

### AC6 — No lifecycle invariant weakened

- No `transitions.json` rule changed shape or content this cycle (confirmed: `transitions.json` does not appear in `git diff origin/main..HEAD --stat`).
- No test was deleted without a replacement citation (the only removed logic — the bash `closeout_violation`/`self_test` functions — has its replacement cited in AC2 above and proven green).
- No CI gate was removed from `build.yml`; `build.yml` is untouched this cycle (confirmed: not in the diff). The `--self-test` CLI branch inside `check-dispatch-closeout-integrity.sh` is removed, but the script itself, and its invocation in `build.yml:332` (`run: ./scripts/ci/check-dispatch-closeout-integrity.sh`, no arguments), are unchanged — the CI gate still runs, still exits non-zero on drift, and the invariant the self-test proved is independently proven by the Go suite (AC2).
- No ambiguity was found on this axis that would require a STOP (see §Debt for the one adjacent finding that was considered and explicitly NOT acted on because it fell outside this cycle's charter, rather than being a live ambiguity about *this* cycle's own changes).

### AC7 — All gates green

Gates run and confirmed locally by α (CI-hosted run status is β's/the merge gate's responsibility per the pre-review gate row 10, but the underlying scripts/tests these gates wrap were all re-run locally and are green):
- `./scripts/ci/check-dispatch-repair-preflight.sh` → exit 0.
- `./scripts/ci/check-dispatch-closeout-integrity.sh` → exit 0 (post-fold).
- `go test ./src/packages/cnos.issues/commands/issues-fsm/... ./src/go/internal/cell/...` → both packages `ok`.
- `bash -n` syntax check on both edited scripts → clean.
- Named CI gates (I1, I2, I4, I5, I6, install-wake-golden, Go, Package, Binary) were not independently re-run by α outside of GitHub Actions (no local harness for the full workflow suite); these are the branch-CI-green rows of the pre-review gate (§CDD Trace / §Review-readiness names the head SHA whose CI status β should confirm green before merge). `install-wake-golden` is unaffected by this cycle's diff (no orchestrator/golden/live-workflow file touched — confirmed via `git diff origin/main..HEAD --stat`).
