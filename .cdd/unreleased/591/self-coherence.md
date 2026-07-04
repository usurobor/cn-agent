# self-coherence — cnos#591

cycle: 591
role: alpha
round: R0

---

## §R0

### Summary of what was built

Implemented `cn cell finalize` — the mechanical checkpoint + PR-open-on-first-matter
finalizer specified in `.cdd/unreleased/591/gamma-scaffold.md` (γ R0 scaffold, cnos#591,
Sub B of #583/AC4). Surfaces touched, exactly per the scaffold's "Surfaces α is expected
to touch" table:

- `src/go/internal/cell/cell.go` — new `FinalizeArgs`/`ParseFinalizeArgs`, new
  `FinalizeFacts` (pure matter-detection facts struct, mirrors `FactSnapshot`'s shape
  idea) + `HasMatter` (pure function, AC2), new `Finalizer` type (mirrors `Returner`'s
  `RunGH`/`RunGHJSON` injection exactly, plus two additional injectable seams —
  `AssembleFacts` and `Checkpoint` — so the checkpoint+PR decision logic is testable
  with zero real git/gh subprocess calls), plus the real production git-plumbing helpers
  (`assembleFacts`, `realCheckpoint`, and the low-level `gitXxx`/`runGitCmd*` functions)
  used only when the injectable fields are left nil.
- `src/go/internal/cli/cmd_cell.go` — new `CellFinalizeCmd` (`Spec()`/`Help()`/`Run()`),
  mirroring `CellReturnCmd`/`CellResumeCmd` exactly.
- `src/go/cmd/cn/main.go` — `reg.Register(&cli.CellFinalizeCmd{})` added next to the
  existing `CellReturnCmd{}`/`CellResumeCmd{}` registrations (found via
  `grep -rn "CellReturnCmd{}" src/go` → `src/go/cmd/cn/main.go:45`, per Friction note 4).
  No other wiring was needed: `ResolveCommand` (`src/go/internal/cli/dispatch.go`)
  resolves `cn cell finalize` → registry key `cell-finalize` purely from
  `CommandSpec.Name`; there is no separate per-noun dispatch table to hand-wire.
- `src/go/internal/cell/cell_test.go` — new tests (see AC section below).
- `src/packages/cnos.core/commands/install-wake/cn-install-wake` — new `if: always()`
  step "Mechanical checkpoint + PR finalizer", a sibling to the existing write-fence
  step, gated on `[ "$role" = "dispatch" ]` (Friction note 1, pinned).
- `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` +
  `.github/workflows/cnos-cds-dispatch.yml` — regenerated (see §Golden regeneration).
- `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` — step 7 of "Invoke δ in
  wake-invoked mode" edited to narrow δ's own responsibility to the `REVIEW-REQUEST.yml`
  preflight + requesting the transition, naming PR-open as the new mechanical finalizer's
  job. No other section touched; the "Lifecycle transitions" table is untouched.

Zero diff in `transitions.json` or any FSM Go file (`src/packages/cnos.issues/commands/issues-fsm/`).
`CellReturnCmd`/`CellResumeCmd` behavior and their existing tests are unmodified.

### AC-by-AC verification

**AC1 — mechanical checkpoint/finalizer exists, not prompt-only prose.**
- `grep -n "CellFinalizeCmd" src/go/internal/cli/cmd_cell.go` → matches (`Spec`, `Help`, `Run`).
- `cd src/go && go build ./...` → clean. `/tmp/cn cell` lists `cell finalize` in the group listing.
- `grep -n 'if: always()' src/packages/cnos.core/commands/install-wake/cn-install-wake` →
  2 occurrences (write-fence + new finalizer step).
- `git diff` of the regenerated golden shows the new step and `/tmp/cn-finalize cell finalize`
  invocation landing in `cnos-cds-dispatch.golden.yml`. PASS.

**AC2 — matter-detection across all four signals.**
- `go test ./internal/cell/... -run TestFinalize_MatterDetection -v` — 6 subtests: one per
  signal isolated (uncommitted / commits-beyond-base / CDD-artifacts / branch-exists), one
  all-four-true, one all-four-false. All PASS.
- `HasMatter` (cell.go) is a pure function over `FinalizeFacts` — no git/gh subprocess,
  confirmed by reading the function body (`grep -n "func HasMatter" src/go/internal/cell/cell.go`).
  PASS.

**AC3 — matter causes branch + commit + PR creation/update (draft).**
- `go test ./internal/cell/... -run TestFinalize_Checkpoint -v` — injects `AssembleFacts`
  (matter present), a recording `Checkpoint` fake, `RunGHJSON` returning `[]` (no existing
  PR), and a recording `RunGH`. Asserts exactly one gh call, containing `pr create`,
  `--draft`, `--head cycle/591`, `--base main`, and `Refs #591` in the body. PASS.
- Manual dogfood smoke (not a unit test; see §Dogfood smoke below) confirms the real
  `assembleFacts`/`realCheckpoint` production path end-to-end against a scratch repo.

**AC4 — PR creation/update is idempotent.**
- `go test ./internal/cell/... -run TestFinalize_Idempotent -v` — `RunGHJSON` returns one
  existing open PR; asserts zero `gh pr create` calls on a first `Finalize` call and zero
  again on an independent second call against unchanged facts. PASS.

**AC5 — no matter, no empty PR.**
- `go test ./internal/cell/... -run TestFinalize_NoMatterNoOp -v` — all-false facts; asserts
  `Finalize` returns nil, `Checkpoint` never invoked, and both `RunGH`/`RunGHJSON` fakes
  recorded zero calls (the `RunGHJSON` fake is wired to fail the test if called at all).
  PASS. Confirmed live: a scratch-repo run on a non-`cycle/*` branch with no CDD matter
  printed "no matter detected — no-op" and exited 0 (§Dogfood smoke, TEST 1).

**AC6 — PR-open failure records STOP/BLOCKED evidence, never claims review.**
- `go test ./internal/cell/... -run TestFinalize_PROpenFailure -v` — injected `RunGH`
  returns an error on `pr create`; asserts the returned error contains
  `cell_finalize_pr_open_failed`, that `.cdd/unreleased/591/FINALIZE-STOP.md` was written
  naming the marker and the raw error, and that no `--add-label`/`--remove-label` call
  occurred. PASS.
- Confirmed live against a scratch repo (§Dogfood smoke, TEST 2 and TEST 3): a real push
  failure (no remote configured) produced `cell_finalize_checkpoint_failed` +
  `FINALIZE-STOP.md`; a real `gh pr list` failure (scratch repo has no GitHub host)
  produced `cell_finalize_pr_list_failed` + an overwritten `FINALIZE-STOP.md` — both exits
  non-zero, no label mutation attempted.

**AC7 — PR is draft unless review-readiness already proven.**
- The AC3 test's captured `gh pr create` args unconditionally include `--draft`.
- `grep -n "fsm evaluate\|REVIEW-REQUEST" src/go/internal/cell/cell.go` — the only matches
  are in new doc comments stating the Finalizer never touches either surface; zero
  matches in executable code. `TestFinalize_NeverTouchesLabelsOrFSM` re-asserts this at
  runtime (stdout/stderr never mentions `REVIEW-REQUEST.yml`; no gh call contains
  `fsm evaluate`). PASS.

**AC8 — FSM remains the only label-writer; finalizer never writes labels.**
- `grep -n "add-label\|remove-label\|issue edit" src/go/internal/cell/cell.go` — every
  match is inside the pre-existing `Returner`/`applyLabelTransition` code (lines
  documented in the file's original F3 comments), none inside the new `Finalizer`,
  `ensurePR`, `stopWithEvidence`, or `assembleFacts`/`realCheckpoint` code. Verified by
  reading the diff line-by-line, not just the final-state grep.
- Every `TestFinalize_*` test's injected fakes assert (implicitly, by construction of the
  fakes and explicitly in `TestFinalize_NeverTouchesLabelsOrFSM`) that no
  `--add-label`/`--remove-label` call occurred, in both success and failure paths. PASS.
- No `cn issues fsm evaluate --apply` call was added — per Friction note 2, AC8 is
  satisfied vacuously (zero label-mutation code paths at all), not by adding an
  unneeded/no-op FSM call.

**AC9 — the prior strand class is covered by fixture/proof.**
- `go test ./internal/cell/... -run TestFinalize_AC9_StrandClassRecovered -v` — fixture:
  `FinalizeFacts{Issue: 591, BranchExists: true, CommitsBeyondBase: 5}` (no uncommitted
  changes, no CDD-artifacts signal needed), no existing PR (`RunGHJSON` → `[]`). Asserts
  exactly one `pr create` call and that stdout never contains "no-op". `FinalizeFacts` has
  no `run_active` field at all — the fixture cannot supply one, by construction, per
  Friction note 3. PASS.

**AC10 — all gates green** (checked below, not a per-AC unit test).

### Golden regeneration

Order followed per the task: renderer change landed first, then the `SKILL.md` doctrine
edit, then **both** goldens were regenerated (so the golden reflects the updated prompt
body), then the live workflow was regenerated from the golden.

```
$ ./src/packages/cnos.core/commands/install-wake/cn-install-wake agent-admin
cn-install-wake: agent-admin → .../cnos-agent-admin.golden.yml (unchanged)
$ ./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch
cn-install-wake: cds-dispatch → .../cnos-cds-dispatch.golden.yml (rendered)
$ ./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml
cn-install-wake: cds-dispatch → .github/workflows/cnos-cds-dispatch.yml (rendered)
```

`agent-admin`'s golden is byte-identical (unchanged) because the new finalizer step is
gated on `role == "dispatch"`; `agent-admin` is `role: admin` and never emits it. This is
the expected, intentional shape — not a missed regeneration.

`sha256sum .github/workflows/cnos-cds-dispatch.yml src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`
→ identical hash (`6c0da28f...`) confirmed before commit.

After committing the regenerated goldens (this commit), re-running both renderer
invocations produces "(unchanged)" for both files — golden parity confirmed.

### Leak-audit (AC7/AC8 of install-wake-golden.yml; HARD gate)

Ran the exact greps from `.github/workflows/install-wake-golden.yml` lines 671-689
against the full renderer file (not just the diff), after the renderer edit landed:

```
$ grep -nE 'admin_only|disallowed_surfaces|defer_path|cell_execution' src/packages/cnos.core/commands/install-wake/cn-install-wake \
    | grep -vE '^[0-9]+:[[:space:]]*#' | grep -vE '(die|manifest_error|log_writer_mis_declaration)[[:space:]]'
(no output — 0 admin-shape leaks)

$ grep -ciE 'protocol:cds|cdr|cdw|dispatch:cell|status:todo' src/packages/cnos.core/commands/install-wake/cn-install-wake
0
```

Both leak sets are zero across the **entire file**, not just the new lines — no
pre-existing leaks were found either (CI was not already broken before this cycle).
The new renderer prose deliberately refers to "role == \"dispatch\"", "dispatch-shaped
wakes", and "a lifecycle status label" generically, never naming the concrete protocol
or a label literal.

### Design choices where the scaffold left flexibility

- **Idempotent PR update vs. no-op (Friction note 5):** chose (b) — a true no-op. When
  `gh pr list --head cycle/{N} --state open` returns an existing PR, `Finalize` reports its
  URL and returns nil without ever calling `gh pr edit`. Rationale: the finalizer's body
  content (`Refs #{N}` + a static description) never goes stale across reruns, so a
  refresh-on-every-run `gh pr edit` call would be a wasted API call with no behavior
  change — exactly the case Friction note 5 says is acceptable to skip.
- **Auto-detection fallback (`scanForIssueWithMatter`):** implemented as `git status
  --porcelain -- .cdd/unreleased/` (uncommitted) unioned with `git diff --name-only
  {base-sha}..HEAD -- .cdd/unreleased/` (newly-committed, only when `--base-sha` is
  given), taking the first issue number extracted from any matching path. This is the
  literal reading of the implementation contract's "Command signature" row.
- **Two additional injectable seams beyond `RunGH`/`RunGHJSON`** (`Finalizer.AssembleFacts`,
  `Finalizer.Checkpoint`): the scaffold pins `RunGH`/`RunGHJSON` shape exactly (honored,
  same two field names/signatures as `Returner`) but also requires the facts struct
  itself be "constructible directly in tests without invoking git at all." Since
  `Finalize` needs to assemble facts and perform git checkpoint mutation before it ever
  reaches a `RunGH`/`RunGHJSON` call, two more function-valued fields (mirroring
  `Resumer.CurrentBranch`'s existing injection precedent in this same package) were added
  so unit tests can supply a facts literal and a no-op checkpoint directly — no new
  binary, no new command family, purely internal seams on the existing `Finalizer` type.
- **STOP-evidence marker names:** `cell_finalize_checkpoint_failed`,
  `cell_finalize_pr_list_failed`, `cell_finalize_pr_open_failed`,
  `cell_finalize_pr_list_unparseable`, `cell_finalize_facts_unavailable` — following the
  `review_return_*`/`review_resume_*` convention already used in this package.

### Dogfood smoke (AC3/AC5/AC6, real end-to-end run)

Per the task's guidance, this was run against a **scratch git repo**, not the real cnos
checkout, to avoid creating a real duplicate PR against cycle/591 before δ's converge-time
PR-open (this cycle's own commits already exist on `cycle/591` with an as-yet-unknown PR
state at α time; running the real finalizer here would race that concern unnecessarily).

Built `/tmp/cn` via `cd src/go && go build -o /tmp/cn ./cmd/cn`, then:

- **TEST 1 (no matter, main branch):** `cn cell finalize` from a hub-marked (`.cn/`)
  scratch repo on a bare `main` branch with no CDD artifacts → printed
  `no matter detected — no-op (nothing to checkpoint, no PR touched).`, exit 0. Zero git
  mutation, zero gh calls (none possible: `gh` was never invoked).
- **TEST 2 (matter, no remote):** checked out `cycle/9001`, added
  `.cdd/unreleased/9001/gamma-scaffold.md`, committed → ran `cn cell finalize` → printed
  the matter-detected message, attempted `git push -u origin cycle/9001` against a repo
  with no `origin` remote configured, which failed → error returned:
  `cell_finalize_checkpoint_failed: pushing cycle/9001: git push -u origin cycle/9001: exit status 128`,
  exit 1, and `.cdd/unreleased/9001/FINALIZE-STOP.md` was written naming the marker and
  the raw error. No gh call was ever attempted (checkpoint failed first).
- **TEST 3 (matter, real local remote, no real GitHub host):** added a local bare
  repository as `origin`, reran `cn cell finalize` → checkpoint succeeded this time
  (`✓ checkpoint complete: cycle/9001 committed and pushed`), then `gh pr list` failed
  because the scratch repo's remote is not a known GitHub host → error:
  `cell_finalize_pr_list_failed: gh pr list --head cycle/9001 --state open --json number,url: exit status 1 (stderr: none of the git remotes configured for this repository point to a known GitHub host...)`,
  exit 1, `FINALIZE-STOP.md` overwritten with the new marker. No label mutation attempted
  in either failure path.

This confirms the real (non-injected) `assembleFacts`/`realCheckpoint`/`ensurePR` code
paths behave sanely under both success and two distinct realistic failure modes, without
ever touching this cycle's real GitHub PR state. The real draft-PR-create success path
(`gh pr create --draft ...` actually succeeding against a real, authenticated GitHub repo)
is exercised by the injected-fake unit tests (AC3/AC9) and will be exercised again for
real by δ's converge-time PR-open against `cycle/591` itself — the scaffold's own AC3
manual-smoke note names this cycle as "the dogfood case" for that specific step, which is
intentionally deferred to δ rather than run here to avoid a race on this branch's PR
state.

### Friction encountered

- None of the five named friction notes required a genuinely open decision beyond what
  is documented above — notes 1-4 were pinned/locatable exactly as described; note 5's
  flexibility is resolved and documented above.
- Minor: `Invocation.HubPath` (used as `Finalizer.RepoRoot`) is the discovered `.cn/`
  hub root, which for this repo is the git repository root — confirmed by reading
  `cmd_cell.go`'s existing `CellResumeCmd.Run` (`RepoRoot: inv.HubPath`) before reusing
  the same field for `CellFinalizeCmd`, rather than assuming.
- No FSM/transitions.json file needed touching; confirmed via `git diff --stat` against
  `src/packages/cnos.cds/skills/cds/fsm/transitions.json` and
  `src/packages/cnos.issues/commands/issues-fsm/` — empty, as required.

### Full verification run (AC10, checked here at R0 as well as at β/converge)

```
$ cd src/go && go build ./... && go test ./... && go vet ./...
ok (all packages, including internal/cell with the new Finalize* tests)
```

```
$ bash scripts/ci/check-dispatch-repair-preflight.sh
cnos#516 repair-preflight guard: ... rc=0
$ bash scripts/ci/check-dispatch-closeout-integrity.sh
cnos#524 closeout-integrity guard: ... rc=0
```

### Review-ready signal

**REVIEW-READY.** All ten ACs have oracle-backed evidence above; `go build`/`go test`/
`go vet` are green; both goldens regenerated and byte-identical to their live-workflow
counterpart (sha256 match); the renderer's AC7/AC8 leak-audit greps return zero matches
against the full file; `transitions.json` and all FSM Go files carry zero diff;
`CellReturnCmd`/`CellResumeCmd` are unmodified. Ready for β review.
