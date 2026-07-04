# γ R0 scaffold — cnos#591

## Issue reference

- **Issue:** [usurobor/cnos#591](https://github.com/usurobor/cnos/issues/591) — "cds/dispatch: mechanical checkpoint + PR-open-on-first-matter finalizer (Sub B of #583)"
- **Mode:** design-and-build
- **Parent:** #583 (master wave — mechanical dispatch cell). **Sub B** of #583 (AC4) — the smallest mechanical extraction that structurally kills the "matter but no PR" strand class.
- **Precondition:** #575 (FSM Phase 3 — claim/hard-block/release-back-to-queue routed through `cn issues fsm evaluate --apply`) — landed, merged to `main`. Verified: `main`@`fe4b6d0b8f32b6af2a00cc8922b85c3e70876819` contains `claim_request_present`/`block_request_present`/`release_request_present` guards in `src/packages/cnos.cds/skills/cds/fsm/transitions.json`.
- **Protocol:** cds
- **Wake run id:** local-manual-dispatch (this cycle ran as a manually-invoked `cds-dispatch` wake firing, executed as `sigma`; no substrate Actions run id — recorded on the claim comment)
- **Cycle branch:** `cycle/591`, created from `origin/main`@`fe4b6d0b8f32b6af2a00cc8922b85c3e70876819`
- **Base-SHA discrepancy:** none — branch created at the exact SHA δ observed as `main` HEAD at claim time.

## Governing rule (restated from the issue)

> If a run ends with matter, there must be a branch/commit/PR — or an explicit blocked/STOP record. No more "branch has work, issue stranded in `status:in-progress`, no PR."

Checkpoint + PR creation are a **mechanical runtime/finalizer responsibility**, not a final cognitive-agent step. Workers produce matter; runtime checkpoints + opens/updates the PR; the FSM owns status labels; skills provide cognition only.

## Empirical anchors

- **cnos#368/#504/#568/#571/#574(×2)** — the strand-class recurrence: a run ends with matter on the branch (or only in the workspace) and no PR, stranding the issue at `status:in-progress`. AC9's fixture must prove this class is now unreachable.
- **cnos#575** — Phase 3 FSM authority (claim/hard-block/release/review all routed through `cn issues fsm evaluate --apply`); this cycle builds ON TOP of that stable authority, must not touch it.
- **cnos#500** — `cn cell return` / `cn cell resume`: the precedent for adding a new `cn cell <verb>` subcommand (package `src/go/internal/cell`, dispatched from `src/go/internal/cli/cmd_cell.go`). This cycle follows that exact precedent (`cn cell finalize`) rather than inventing a new command family.
- **cnos#496 (write-fence)** — the existing `if: always()` mechanical post-run step pattern in the rendered dispatch workflow (`Write fence — dispatch_activation_log_write_violation`, sourced from the renderer at `src/packages/cnos.core/commands/install-wake/cn-install-wake` lines ~979–1047, gated on `activation_log_writer == false`). This cycle's finalizer step is a SIBLING step to that fence, using the SAME `if: always()` + `CN_WAKE_BASE_SHA` baseline convention, but gated on `role == "dispatch"` instead (see Friction note 1 — this is a δ-pinned design decision, not left to α).
- **install-wake-golden.yml AC7/AC8 leak audit** (lines 629–689 of `.github/workflows/install-wake-golden.yml`) — the renderer source (`cn-install-wake`) is grepped for the literal strings `protocol:cds`, `cdr`, `cdw`, `dispatch:cell`, `status:todo` (case-insensitive) and MUST contain zero matches, including in comments. **This is a hard constraint on the new renderer code this cycle adds** — see Implementation contract row "Renderer leak-audit constraint" below. α must not write any comment or string literal containing those tokens anywhere in `cn-install-wake`.

## Design rule (restated, binding)

- **Workers** produce matter (file changes, commits, CDD artifacts).
- **Runtime** (the new `cn cell finalize` command + its renderer-emitted `if: always()` step) checkpoints matter (commit/push) and opens/updates the PR.
- **FSM** owns status labels — the finalizer never writes a label and never calls `cn issues fsm evaluate --apply` (see Friction note 2 for why AC8's "requests via --apply, never writes labels" is satisfied vacuously here, not by adding an unneeded FSM call).
- **Skills** (SKILL.md prose) provide cognition only — they control nothing; this cycle's mechanism must work even if the cognitive session dies immediately after producing matter.

## Surfaces α is expected to touch

| Surface | Path | What changes |
|---|---|---|
| New Go package logic | `src/go/internal/cell/cell.go` | Add a `Finalizer` type (mirrors `Returner`/`Resumer` shape: injectable `RunGH`/`RunGHJSON` funcs, package-level `runGH`/`runGHJSON` fallbacks) implementing matter-detection + checkpoint + idempotent PR create/update. Add `FinalizeArgs` + `ParseFinalizeArgs`. |
| New CLI wiring | `src/go/internal/cli/cmd_cell.go` | Add `CellFinalizeCmd` (mirrors `CellReturnCmd`/`CellResumeCmd` shape: `Spec()`, `Help()`, `Run()`). |
| CLI command registration | wherever `CellReturnCmd`/`CellResumeCmd` are registered (find via `grep -rn "CellReturnCmd{}" src/go/cmd/cn` or `src/go/internal/cli`) | Register `CellFinalizeCmd` as `cell-finalize`, dispatched under the `cn cell finalize` noun-verb pair (mirror however `cell return`/`cell resume` are wired — likely a `cell` dispatcher file; find it, do not guess its name). |
| Go tests | `src/go/internal/cell/cell_test.go` (create if absent — check first) | Table/fixture-driven tests for matter-detection decision logic (pure function, no subprocess) covering AC2 (all 4 matter signals), AC4 (idempotent PR create/update via injected `RunGHJSON`/`RunGH` fakes), AC5 (no-matter no-op), AC6 (PR-open failure → STOP evidence, no label touched), AC9 (the strand-class fixture: run ended, no active run, matter exists, no PR → finalizer creates/recovers the PR). |
| Renderer (substrate emission) | `src/packages/cnos.core/commands/install-wake/cn-install-wake` | Add a new `if: always()` step, siblings to the existing write-fence step, gated on `[ "$role" = "dispatch" ]` (NOT on `activation_log_writer`, a different concern — see Friction note 1). The step builds `cn` fresh (`cd src/go && go build -o /tmp/cn-finalize ./cmd/cn`) and runs `/tmp/cn-finalize cell finalize --base-sha "${CN_WAKE_BASE_SHA:-}"` (no `--issue` flag — the command self-detects the issue number from the current branch name / `.cdd/unreleased/*/` diff; see Implementation contract "Issue-number auto-detection"). **Constraint: zero occurrences of `protocol:cds`/`cdr`/`cdw`/`dispatch:cell`/`status:todo` anywhere in this file, including new comments** (install-wake-golden.yml AC7/AC8 audit). |
| Golden fixture (regenerate, do not hand-edit) | `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` | Regenerate via `./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch` after the renderer change lands — this file is rendered output, never hand-edited. |
| Live substrate artifact (regenerate, do not hand-edit) | `.github/workflows/cnos-cds-dispatch.yml` | Regenerate via `./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml` so the live workflow's sha256 matches the golden's (install-wake-golden.yml's "Verify live cds-dispatch workflow matches golden" step). |
| Doctrine (small, scoped edit) | `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` | Add a short subsection (near the existing "Invoke δ in wake-invoked mode" step 7, which currently says "δ opens (or updates) a pull request scoped to the cell...") naming that PR creation/checkpointing is now a MECHANICAL finalizer step that runs regardless of whether the cognitive session reaches step 7 — δ's own job at step 7 narrows to ensuring `REVIEW-REQUEST.yml` exists and requesting the `status:review` transition; it no longer needs to be the thing that opens the PR (the finalizer already will have). Do NOT restructure the rest of the file. Do NOT touch the "Lifecycle transitions" table (out of scope — no new label semantics). |
| Go tests (existing suite) | `src/go/internal/cell/` existing tests for `Returner`/`Resumer` | Must remain green, unmodified in assertion. |

## Per-AC oracle list

### AC1 — mechanical checkpoint/finalizer exists, not prompt-only prose

**Invariant:** the finalizer is an explicit `cn` subcommand + an `if: always()` workflow step, not prose asking the agent to open the PR earlier.

**Oracle commands:**
- `grep -n "CellFinalizeCmd" src/go/internal/cli/cmd_cell.go` — exists.
- `cd src/go && go build ./... ` — builds clean; `./cn help cell` (or equivalent) lists `finalize`.
- `grep -n 'if: always()' src/packages/cnos.core/commands/install-wake/cn-install-wake` — at least two occurrences (the existing write-fence + the new finalizer step).
- `git show cycle/591:src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml | grep -n "cell finalize"` — the rendered workflow actually invokes the new command.

**Negative case:** the only change is prose in `cds-dispatch/SKILL.md` with no new `cn` subcommand and no new workflow step — REQUEST CHANGES.

### AC2 — matter-detection across all four signals

**Invariant:** the finalizer detects matter as the OR of: uncommitted file changes, commits beyond base on `cycle/{N}`, CDD artifacts under `.cdd/unreleased/{N}/`, an existing `cycle/{N}` branch.

**Oracle commands:**
- `cd src/go && go test ./internal/cell/... -run TestFinalize_MatterDetection -v` — a table test with 4+ cases, one per signal isolated, each independently sufficient to flip `hasMatter` to true; a fifth case with all four false confirms `hasMatter` is false only then.
- Read the matter-detection function directly (`grep -n "func.*[Mm]atter" src/go/internal/cell/cell.go`) — confirm it is a pure function over an explicit facts struct (mirrors `FactSnapshot` in `issues-fsm/snapshot.go`), not inlined imperative logic scattered through `Finalize`. This makes AC9's strand-class fixture constructible without a real git/gh subprocess.

**Negative case:** a facts struct with all four signals false must NOT report matter.

### AC3 — matter causes branch + commit + PR creation/update (draft)

**Invariant:** when matter exists, the finalizer ensures `cycle/{N}` exists, commits any uncommitted changes, pushes, and creates (or updates) a **draft** PR.

**Oracle commands:**
- `go test ./internal/cell/... -run TestFinalize_Checkpoint -v` — with injected `RunGH`/`RunGHJSON` fakes (mirror `Returner`'s injection pattern exactly), assert: (a) no existing PR + matter → a `gh pr create --draft ...` invocation is issued with `--head cycle/{N}` and a body containing `Refs #{N}`; (b) the created PR is draft (`--draft` flag present in the captured args).
- Manual smoke (recorded in self-coherence.md, not required to be a unit test): actually run `cn cell finalize` against `cycle/591` itself once real matter exists on this branch (this cycle's own commits) and confirm a real draft PR is opened/updated — this cycle IS the dogfood case.

**Negative case:** no matter → no `gh pr create` call at all (see AC5).

### AC4 — PR creation/update is idempotent

**Invariant:** a rerun against a `cycle/{N}` that already has an open PR updates that PR (or is a clean no-op) — never creates a duplicate.

**Oracle commands:**
- `go test ./internal/cell/... -run TestFinalize_Idempotent -v` — inject `RunGHJSON` to return one existing open PR for `--head cycle/{N}` (mirror `gh pr list --head cycle/{N} --state open --json number,url`); assert the finalizer does NOT call `gh pr create` again; assert a second `Finalize` call with the same injected state produces the same result (byte-identical decision) as the first.

**Negative case:** two consecutive `Finalize` calls against unchanged facts create two PRs — FAIL.

### AC5 — no matter, no empty PR

**Invariant:** if none of the four matter signals hold, the finalizer is a clean no-op: no branch created, no commit, no `gh pr create` call, exit 0.

**Oracle commands:**
- `go test ./internal/cell/... -run TestFinalize_NoMatterNoOp -v` — facts struct with all four signals false → `Finalize` returns nil, and the injected `RunGH`/`RunGHJSON` fakes record zero invocations.

**Negative case:** an empty run creates a PR anyway — FAIL (this is the AC9 strand-class's mirror-image failure mode).

### AC6 — PR-open failure records STOP/BLOCKED evidence, never claims review

**Invariant:** if `gh pr create` (or the commit/push step) fails, the finalizer does NOT pretend success, writes a STOP-evidence artifact, and — critically — never touches `REVIEW-REQUEST.yml` or any label.

**Oracle commands:**
- `go test ./internal/cell/... -run TestFinalize_PROpenFailure -v` — inject a `RunGH` fake that returns an error on the `gh pr create` call; assert `Finalize` returns a non-nil error naming `cell_finalize_pr_open_failed` (or equivalent precise marker, following the `review_return_*` / `review_resume_*` naming convention already used in this package); assert a `.cdd/unreleased/{N}/FINALIZE-STOP.md` (or similarly named) evidence file is written with the failure detail; assert no `gh issue edit` / label-mutation call was made (grep the fake's captured calls for `--add-label`/`--remove-label` — zero).

**Negative case:** failure is swallowed silently (`Finalize` returns nil) — FAIL.

### AC7 — PR is draft unless review-readiness already proven

**Invariant:** every PR the finalizer creates is opened with `--draft`; the finalizer never flips it out of draft and never requests `status:review`.

**Oracle commands:**
- Re-use the AC3 test's captured `gh pr create` args — assert `--draft` is present unconditionally.
- `grep -n "fsm evaluate\|REVIEW-REQUEST" src/go/internal/cell/cell.go` scoped to the new `Finalizer` code — zero matches (the finalizer must not touch either surface; that remains δ's job at step 7 of the wake-invoked cycle, per the small doctrine edit in `cds-dispatch/SKILL.md`).

**Negative case:** the finalizer creates a non-draft PR, or writes `REVIEW-REQUEST.yml`, or calls `cn issues fsm evaluate` — FAIL (scope creep into δ's review-request responsibility).

### AC8 — FSM remains the only label-writer; finalizer never writes labels

**Invariant:** zero label-mutation calls anywhere in the new `Finalizer` code path.

**Oracle commands:**
- `grep -n "add-label\|remove-label\|issue edit" src/go/internal/cell/cell.go` scoped to new code — zero matches outside the pre-existing `Returner.applyLabelTransition` (which is a DIFFERENT, pre-existing command's code, untouched by this cycle).
- Every test in the `TestFinalize_*` family asserts (via the injected fake's call log) that no `gh issue edit ... --add-label ... / --remove-label ...` invocation occurred, for both the success and failure paths.

**Negative case:** any new code path calls `gh issue edit` with a label flag — FAIL, regardless of which label.

**Design note (binding, not left to α):** this AC is satisfied by the finalizer having **no label-writing code path at all** — not by adding a spurious `cn issues fsm evaluate --apply` call that has no transition to request. See Friction note 2.

### AC9 — the prior strand class is covered by fixture/proof

**Invariant:** run ended · no active run · matter exists · no PR → the finalizer creates/recovers the PR (or blocks with evidence); this exact scenario (the #574/#568 strand) no longer strands.

**Oracle commands:**
- `go test ./internal/cell/... -run TestFinalize_AC9_StrandClassRecovered -v` — a fixture reconstructing exactly: branch `cycle/{N}` exists with commits beyond base (matter), no PR exists (`RunGHJSON` fake returns empty PR list), `run_active` is irrelevant to this command (the finalizer runs post-hoc, it does not itself check `run_active` — that is the FSM's claim-guard concern, not this command's; see Friction note 3). Assert: the finalizer calls `gh pr create --draft` and the resulting decision is "PR created", not "no-op" and not "blocked".

**Negative case:** the same fixture but the finalizer reports no-op or fails silently — this IS the bug class AC9 exists to kill; FAIL.

### AC10 — all gates green

**Invariant:** I1, I2, I4, I5, I6, install-wake-golden, dispatch-repair-preflight, dispatch-closeout-integrity, Go, Package, Binary all pass on the pushed branch.

**Oracle commands (checked at β/converge time, not R0):**
- `cd src/go && go build ./... && go test ./...` — full Go module green.
- `./src/packages/cnos.core/commands/install-wake/cn-install-wake agent-admin && ./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch && git diff --exit-code -- src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` — clean (goldens match re-render after α commits the regenerated goldens).
- `sha256sum .github/workflows/cnos-cds-dispatch.yml src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` — identical.
- The AC7/AC8 leak-audit greps from `install-wake-golden.yml` lines 671–689, run manually against `cn-install-wake` — zero matches for both the admin-shape and dispatch-shape leak sets.
- `bash scripts/ci/check-dispatch-repair-preflight.sh` and the closeout-integrity equivalent (find via `ls scripts/ci/`) — green (this cycle does not touch repair/closeout-integrity logic, so these should be unaffected regressions checks).

## Source-of-truth table

| Claim / surface | Canonical source | Status |
|---|---|---|
| `cn cell` command family precedent (`return`, `resume`) | `src/go/internal/cell/cell.go` + `src/go/internal/cli/cmd_cell.go` | Shipped (cnos#500); this cycle extends it additively with `finalize` |
| CDS FSM (claim/hard-block/release/review via `--apply`) | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` + `src/packages/cnos.issues/commands/issues-fsm/` | Shipped (cnos#568/#569/#574/#575); NOT touched by this cycle — the finalizer never calls it |
| Dispatch workflow renderer | `src/packages/cnos.core/commands/install-wake/cn-install-wake` | Shipped; extension point for the new `if: always()` finalizer step, gated on `role == "dispatch"` |
| Existing write-fence step (sibling pattern) | same renderer, ~lines 979–1047 | Shipped (cnos#496); this cycle's step follows the identical `if: always()` + `CN_WAKE_BASE_SHA` convention |
| cds-dispatch wake manifest + prompt | `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` | Current; small scoped doctrine edit near "Invoke δ" step 7 |
| Golden fixture + live workflow parity gate | `.github/workflows/install-wake-golden.yml` | Shipped; this cycle's renderer change MUST regenerate both goldens and keep the AC7/AC8 leak audit clean |
| Master wave context | #583 | OPEN, tracks this sub |
| Sub A (mechanism/cognition doctrine) | #584 | Merged |
| FSM Phase 3 (precondition) | #575 | Merged |

## The α dispatch prompt

```text
You are α. Project: cnos.

Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.

Issue: gh issue view 591 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/591

Tier 3 skills:
- src/packages/cnos.core/skills/write/SKILL.md
- src/packages/cnos.cdd/skills/cdd/issue/SKILL.md
- src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md
- src/packages/cnos.core/skills/eng/SKILL.md
- src/packages/cnos.core/skills/eng/ship/SKILL.md
- src/packages/cnos.core/skills/eng/evolve/SKILL.md
- src/packages/cnos.core/skills/eng/process-economics/SKILL.md
- src/packages/cnos.core/skills/eng/test/SKILL.md
- src/packages/cnos.core/skills/go/SKILL.md

Scaffold: .cdd/unreleased/591/gamma-scaffold.md (this file — read it in full before coding;
it carries the per-AC oracle list, the surfaces you are expected to touch, the source-of-truth
table, the implementation contract below, and named friction notes you must resolve, not defer).

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Go — extends the existing `cn` CLI kernel (`src/go/`) and the `cell` package (`src/go/internal/cell/`); no new runtime. |
| CLI integration target | New sub-verb `cn cell finalize`, added to the existing `cell` command family (cnos#500 precedent: `cell return`, `cell resume`). Do NOT create a new binary, a new top-level noun, or a shell-script-only finalizer. |
| Command signature | `cn cell finalize [--base-sha SHA]` — NO required `--issue` flag. The command self-detects the target issue number: (1) if the current git branch matches `cycle/(\d+)`, use that N; (2) else, scan for uncommitted or (if `--base-sha` given) newly-committed paths under `.cdd/unreleased/{N}/` and use the N with matter. If no N can be determined by either method, this IS the AC5 no-matter case — exit 0, no-op. (An optional `--issue N` override MAY be added for testability/manual invocation, but auto-detection is the primary path the renderer's mechanical step relies on — it cannot know N in advance.) |
| Package scoping | `src/go/internal/cell/cell.go` (new `Finalizer` type + matter-detection pure function + `FinalizeArgs`/`ParseFinalizeArgs`) + `src/go/internal/cli/cmd_cell.go` (new `CellFinalizeCmd`) + wherever the `cell` commands are registered (find it — do not guess) + `src/go/internal/cell/cell_test.go` (new or extended) + `src/packages/cnos.core/commands/install-wake/cn-install-wake` (new renderer step) + regenerated goldens (`src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`, `.github/workflows/cnos-cds-dispatch.yml`) + a small scoped doctrine edit in `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`. |
| Matter-detection shape | A pure function over an explicit facts struct (mirror `FactSnapshot` in `src/packages/cnos.issues/commands/issues-fsm/snapshot.go`) — NOT inline imperative git-shelling scattered through `Finalize`. This is what makes AC2/AC9 fixture-testable without a real git/gh subprocess. |
| Testability shape | Mirror `Returner`'s injection pattern EXACTLY: `RunGH func(ctx, args, w) error` and `RunGHJSON func(ctx, args) ([]byte, error)` fields on `Finalizer`, defaulting to package-level `runGH`/`runGHJSON` when nil. Git plumbing (status/branch/commit/push) may use real `exec.Command` calls in production but MUST be reachable through an equivalently injectable seam for the matter-detection facts specifically (the facts struct itself should be constructible directly in tests without invoking git at all — see AC2/AC9). |
| Renderer gating | The new `if: always()` finalizer step in `cn-install-wake` is gated on `[ "$role" = "dispatch" ]` (the variable already set earlier in the script at the role-enum-check block), NOT on `activation_log_writer`. This makes the finalizer apply to any future dispatch-shaped wake, not just `cds-dispatch` — consistent with the renderer's protocol-agnostic architecture. Do NOT reuse the `activation_log_writer` gate for this — it is a different concern (channel-log write ownership vs. matter checkpointing). |
| Renderer leak-audit constraint (HARD) | `.github/workflows/install-wake-golden.yml`'s AC7/AC8 step greps `cn-install-wake` for the literal (case-insensitive) strings `protocol:cds`, `cdr`, `cdw`, `dispatch:cell`, `status:todo` and requires ZERO matches, including in comments. The new renderer code and ALL its comments MUST avoid these exact substrings. Refer to the dispatch role generically ("dispatch-shaped wakes", "role: dispatch") — never name the concrete protocol or its labels in the renderer file. |
| No FSM interaction | The `Finalizer` MUST NOT call `cn issues fsm evaluate`, MUST NOT write `REVIEW-REQUEST.yml`, and MUST NOT run any `gh issue edit --add-label/--remove-label`. It only touches: git (branch/commit/push) and `gh pr create`/`gh pr list`/`gh pr edit` (for idempotent update, if you choose to update an existing PR's body/state rather than treat "PR exists" as a pure no-op — either is acceptable, name the choice in self-coherence.md). |
| PR shape | Always `--draft`. Title/body should reference the issue (`Refs #{N}` or `Part of #{N}`) so `gh pr list --head cycle/{N}` type lookups and future closeout-integrity checks can find it. Base branch: `main`. |
| STOP-evidence artifact | On PR-open (or checkpoint) failure, write `.cdd/unreleased/{N}/FINALIZE-STOP.md` (create the directory if absent) naming the failure class and the raw error. Use an error marker string following this package's existing convention (`review_return_*`, `review_resume_*`) — e.g. `cell_finalize_pr_open_failed`, `cell_finalize_checkpoint_failed`. |
| Existing-binary disposition | Extend `cn` — no new binary, no new package-command exec-dispatch surface. |
| Runtime dependencies | None beyond the existing `git` and `gh` CLI dependencies this package already assumes (see `runGH`/`runGHJSON` in `cell.go`). |
| Backward-compat invariant | `CellReturnCmd`/`CellResumeCmd` and their existing tests are untouched. `transitions.json` and the FSM engine are untouched (zero diff). The write-fence step's existing YAML in `cn-install-wake` (~lines 979–1047) is untouched — the new step is ADDED, not interleaved into it. |
```

## The β dispatch prompt

```text
You are β. Project: cnos.

Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.

Issue: gh issue view 591 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/591

Scaffold: .cdd/unreleased/591/gamma-scaffold.md (this file). Walk the per-AC oracle list
independently — do not trust α's self-coherence report as a substitute for running the
oracle commands yourself.

Per-AC oracle list to walk independently (full detail in the scaffold's "Per-AC oracle
list" section):

- AC1: confirm `CellFinalizeCmd` exists and is wired into the CLI; confirm a NEW
  `if: always()` step exists in the renderer output distinct from the write-fence step.
- AC2: run the matter-detection table test; confirm each of the four signals independently
  flips the decision, and all-false does not.
- AC3: confirm the checkpoint+PR-create path via the injected-fake test; confirm the PR is
  draft, targets `main`, and the branch/head naming is `cycle/{N}`.
- AC4: confirm idempotency — an existing-open-PR fixture does NOT trigger a second create.
- AC5: confirm a no-matter fixture makes zero git/gh calls and exits 0.
- AC6: confirm a PR-open failure writes FINALIZE-STOP.md-equivalent evidence, returns a
  precise error, and makes NO label-mutation call.
- AC7: confirm every created PR is `--draft`; confirm the Finalizer code never references
  `REVIEW-REQUEST.yml` or `cn issues fsm evaluate`.
- AC8: grep the new Finalizer code for `add-label`/`remove-label`/`issue edit` — must be
  zero matches. This is the one β should be most skeptical about — re-read the diff line
  by line, don't just trust a clean grep on the final state (check the diff didn't
  temporarily/conditionally add one either).
- AC9: confirm the strand-class fixture test exists and asserts PR-creation (not no-op,
  not silent failure) when matter exists with no PR.
- AC10 (checked before converge, not per-round): run `cd src/go && go build ./... && go
  test ./...`; regenerate BOTH goldens via `cn-install-wake` and confirm zero diff; confirm
  live workflow sha256 == golden sha256; manually run the AC7/AC8 leak-audit greps from
  install-wake-golden.yml against `cn-install-wake` and confirm zero matches for BOTH the
  admin-shape and dispatch-shape leak sets — this is a HARD gate α was warned about; treat
  any leak as severity-blocking, not a nit.

## Implementation-contract verification (β Rule 7)

Confirm the diff matches every pinned axis in the α prompt's `## Implementation contract`
table: Go only, extends the EXISTING `cell` command family (no new binary, no new noun),
`cn cell finalize` with NO required `--issue` flag (self-detecting), matter-detection is a
pure function over an explicit facts struct, `Finalizer` injectable exactly like `Returner`,
renderer gated on `role == "dispatch"` (not `activation_log_writer`), zero
`protocol:cds`/`cdr`/`cdw`/`dispatch:cell`/`status:todo` substrings anywhere in
`cn-install-wake` including comments, zero FSM/label-mutation calls anywhere in the new
code, PR always draft targeting `main`. Any drift (new binary, new CLI noun, imperative
non-pure matter detection untestable without subprocess, wrong gate variable, a leaked
protocol-string constant, an FSM call, a non-draft PR) is REQUEST CHANGES, severity D,
classification `implementation-contract` — regardless of whether the behavioral ACs
happen to pass.

## Scope guardrails (restated from the issue — do not expand)

- Do NOT build mechanical recovery/finalizer for stranded states BEYOND checkpoint+PR-open
  (that is Sub C, a separate future issue).
- Do NOT add per-phase checkpointing of the full lifecycle beyond "matter exists →
  checkpoint + PR".
- Do NOT change what skills decide (cognition stays with γ/α/β; this cycle only moves
  mechanical actions out of the cognitive path).
- Do NOT add new status labels/taxonomy; do NOT touch `cell_kind` enforcement; do NOT
  change the wake source model; do NOT start Demo 0.
- Do NOT touch #564. Do NOT touch #585/#586 unless directly required (they are not).
- Do NOT modify `transitions.json` or any FSM Go file — zero diff there.
- Do NOT restructure the write-fence step; add a sibling step instead.
```

## Friction notes (for α to resolve — γ does not decide these, except where explicitly pinned above)

**1. Renderer gate variable choice — PINNED, not open.** Unlike #575's friction notes, this
one is a δ decision, stated here for transparency rather than left to α: the new finalizer
step is gated on `role == "dispatch"`, read from the same `$role` shell variable the
renderer already computes at its role-enum-check block (~line 488 of `cn-install-wake`).
This applies the finalizer to any future dispatch-shaped wake uniformly (matches the
framework's protocol-agnostic architecture — see `delta/SKILL.md` §9.1's framing that
wake-invoked-δ is "protocol-agnostic at the framework level"). Reusing
`activation_log_writer` instead would conflate two independent concerns (channel-log
write-ownership vs. matter-checkpointing) that happen to currently coincide for
`cds-dispatch` but need not for a future wake.

**2. Why AC8 does NOT need a `cn issues fsm evaluate --apply` call.** The issue's AC8 text
("the finalizer requests transitions via `cn issues fsm evaluate --apply`, never writes
labels") is boilerplate carried over from the broader dispatch-protocol lifecycle-transition
language used elsewhere in this wave (claim/hard-block/release/review). This SPECIFIC
cell's finalizer has no lifecycle transition to request: opening a draft PR is explicitly
NOT a `status:review` transition (see the issue's "Important distinction" section) and
matter-checkpointing is not `status:blocked` either (a checkpoint failure is recorded as
evidence for a HUMAN or a later δ hard-block decision, not auto-escalated to `status:blocked`
by this mechanical step — auto-escalating would be scope creep into δ's hard-block judgment,
which the issue's "Out of scope" list protects: "Changing what skills decide"). The correct
reading of AC8 is therefore satisfied **vacuously**: the finalizer contains zero
label-mutation code paths and zero FSM calls, so "never writes labels" holds trivially and
"the FSM remains the only mechanism that applies lifecycle status labels" holds because
this new code never attempts to be a second one. Do not add a no-op or speculative FSM call
just to make the AC8 oracle "look" satisfied by presence of the string `fsm evaluate` — the
oracle checks for ABSENCE of label-writing, not presence of an FSM call.

**3. `run_active` is not this command's concern.** The FSM's claim guard (`todo → in-progress`)
already checks `run_active` to prevent competing claims (cnos#575). The finalizer runs
strictly AFTER a cell has already been claimed and worked — by the time it runs, whether
some OTHER run is active is irrelevant to "does matter exist on this branch right now."
Do not add a `run_active` check to the finalizer; AC9's fixture explicitly does not
supply one.

**4. Finding the `cell` command registration site.** `cmd_cell.go` defines `CellReturnCmd`/
`CellResumeCmd` but this scaffold did not locate where they are registered into the `cn`
noun-verb dispatch table (likely `src/go/cmd/cn/main.go` or an `internal/cli` registry
file). α must find this via `grep -rn "CellReturnCmd{}" src/go` and wire `CellFinalizeCmd`
in the same place, following the exact same registration shape — do not invent a
different wiring mechanism for the new command.

**5. Idempotent PR update vs. no-op.** AC4 requires "a rerun updates the existing PR; no
duplicate PRs." α may implement "update" as either (a) an actual `gh pr edit` call
refreshing the body/labels-of-PR (not issue labels) on each run, or (b) a true no-op that
just re-reports the existing PR's URL without calling `gh pr edit` at all, IF the PR body
already contains everything needed and nothing new needs refreshing on this pass. Either
satisfies "no duplicate PRs"; α picks and documents which in self-coherence.md. Do not
treat "must update" as requiring a wasteful edit call on every single finalize invocation
if nothing changed.
