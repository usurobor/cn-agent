# β review — cycle #656

## Verdict: **converge**

---

## Independent verification performed (not a re-read of α's claims)

Per RECEIPT-VALIDATION.md's discipline (independent verification, not
trusting α's self-report), this review:

1. Rebuilt and re-ran the full test suite from a clean state on both
   affected Go modules (`src/go`, `src/packages/cnos.core/commands/
   label-doctor`) — `go build ./... && go vet ./... && go test ./...`,
   all green.
2. Re-ran the CUE self-test (`./scripts/ci/validate-repo-state.sh
   --self-test`) — all 6 fixtures pass/fail exactly as their filenames
   claim.
3. Re-ran a live end-to-end smoke test independent of any unit test: real
   `cn build` → real local package index → real `cn repo install` into a
   fresh `git init`-only temp repo → real `cn repo status` (clean,
   drifted-via-hand-edit, `--check` exit codes, `git status --porcelain`
   before/after write-check).
4. Dispatched an independent adversarial-review subagent (not primed with
   α's own self-coherence narrative — given only the diff and the issue
   context) specifically instructed to find correctness bugs in the P2/P3
   drift-classification and ledger-write logic, the CUE schema/fixture
   correspondence, the P3 read-only invariant, and the `internal/
   dispatchrender` extraction's behavioral parity with the pre-existing
   `runDispatchCds` code it replaced.
5. Independently re-read the CUE schema's closedness reasoning by hand-
   crafting an ad hoc fixture and confirming `cue vet` rejects it, rather
   than trusting the checked-in fixture corpus alone.

## AC-by-AC walk (code-first)

Re-verified each AC in `self-coherence.md`'s ACs section against the
actual code and actual test output, not the prose:

- **P1** (schema + CI gate): confirmed `schemas/repo_state.cue` is a
  closed `#RepoState` definition (no trailing `...`); confirmed by hand
  that adding an arbitrary field to a valid fixture makes `cue vet` fail
  with `field not allowed`. Confirmed all 4 invalid fixtures fail for
  their claimed reason (read each `cue vet` diagnostic, not just the
  script's pass/fail summary).
- **P2** (render contract + A2 classification): confirmed
  `internal/dispatchrender.Render` is the ONE code path both
  `repoinstall.runDispatchCds` and `repostatus.dispatchStatus` call —
  read both call sites to confirm neither independently reconstructs
  renderer args (which would risk the exact drift this extraction exists
  to prevent).
- **P3** (read-only): read `internal/repostatus/repostatus.go` end to end
  for any `os.WriteFile`/`os.Mkdir`/`os.Create`/`os.Rename` call reachable
  from `Run` — found none outside the throwaway `os.MkdirTemp` (cleaned up
  via `defer os.RemoveAll`) used for the fresh-render comparison, which
  writes only inside the temp dir, never `RepoRoot`. Confirmed
  `labeldoctor.Doctor` with `DryRun: true` returns before its apply loop
  (read `doctor.go`'s `Doctor` function directly).
- **Determinism**: confirmed `RepoState.Marshal`'s `Sort()` call is
  unconditional (not opt-in), and that every current producer
  (`repoinstall.writeRepoState`) builds `ManagedFiles`/`ManagedDirs` via
  fixed sequential appends with no map iteration in between — the sort
  key (`Path`) is unique for every entry this code path can currently
  produce.

## Findings

Three real issues were found and fixed during this review (two by an
independent adversarial subagent, one by α's own pre-review self-check).
All three are fixed on this branch with regression tests; none required
an iterate round back to α — the fixes were applied and re-verified
within this same review pass.

### Finding 1 (β, self-check before requesting review) — ledger silently dropped a removed dispatch workflow

`repostatus.dispatchStatus` originally reported `Present: true` for a
ledger-recorded workflow whose file no longer existed on disk (the
`liveErr != nil` branch fell through to the same `status := DispatchStatus
{Present: true, ...}` used for the matching-content case, before setting
`Drift: DriftUnknown`). A deleted-but-still-ledger-recorded dispatch
workflow was invisible to both the human report and the top-level `Drift`
bool.

**Fixed**: a distinct `DriftRemoved` value, `Present: false`, counted
toward the top-level `Drift` bool. Test: `TestDispatchStatus_Removed`.

### Finding 2 (independent adversarial review) — stale workflow file from a prior run corrupts the ledger's render-contract metadata

`canWriteLedger` (the guard deciding whether `cn repo install` writes/
updates `.cn/repo.state.json` after a `--dispatch cds` failure) was
originally derived from `os.Stat(dispatchWorkflowPath(...))` — "does a
file exist at the fixed output path." That check can't distinguish "this
run's render just produced this file" from "a PRIOR successful run
rendered this file, and THIS run's identity gate failed before ever
calling the renderer again." In the latter case the ledger would be
rewritten describing an untouched file using THIS run's (different,
possibly invalid) `opts.Agent`/`opts.WorkflowPatSecret`/`opts.BotName`/
`opts.BotID` — silently corrupting the render-contract metadata a later
`cn repo status` fresh-render comparison depends on for correctness.

**Fixed**: `runDispatchCds`'s signature now returns `(rendered bool, err
error)`; `Run`'s `canWriteLedger` gates on `rendered`, not file existence.
Test: `TestRun_DispatchCds_StaleWorkflowFileFromPriorRun_LedgerNotCorrupted`
(runs a real successful sigma/engine install, then a real failing
acme-no-PAT rerun, and asserts both the live file and the ledger's
workflow record are byte-identical/unchanged from the first run).

### Finding 3 (independent adversarial review) — a confirmed drift could be silently reported as "no drift"

`repostatus.dispatchStatus`'s fresh-render classification step (which
runs only after `liveSHA != wf.SHA256` — i.e., drift is already an
established fact) fell back to `DriftUnknown` on every failure path
(temp-dir creation, renderer invocation, reading the fresh render back).
`DriftUnknown` is excluded from the top-level `Drift` bool (correctly, for
its OTHER use — "no ledger record to compare against at all") — so a
confirmed sha256 mismatch whose classification step happened to fail
(e.g. the vendored renderer package was removed, or `/tmp` is
unwritable) would report as clean.

**Fixed**: introduced `DriftUnclassified` — distinct from `DriftUnknown`,
counted toward the top-level `Drift` bool — for "confirmed differs from
ledger, but couldn't determine which kind of drift." Test:
`TestDispatchStatus_Unclassified` (hand-edits the live file AND corrupts
the ledger's `renderer_package` to a nonexistent package, forcing the
classification step to fail on a confirmed mismatch).

## Non-blocking observations (not iterate-worthy, recorded for future phases)

- `sha256File` is duplicated near-identically between `internal/
  repoinstall` and `internal/repostatus`. Already flagged in α's own Debt
  section; agreed this is not worth a shared package at two call sites of
  ~8 lines each.
- A package present in `.cn/deps.lock.json` but no longer in
  `.cn/deps.json` (a manifest/lock drift in the OTHER direction from what
  `orphan_packages` checks) is not surfaced by any current field. Narrow,
  requires a hand-edited `deps.json` to trigger, and out of this phase's
  literal Acceptance text — recorded as Debt, not a blocker.
- `cmd_repo_status.go`'s `renderStatus` (human-readable presentation) has
  no dedicated Go test asserting exact stdout strings — verified only via
  manual smoke test. Consistent with `hubstatus.Run`'s own precedent
  (thin `fmt.Fprintf` presentation layers aren't unit-tested line-by-line
  elsewhere in this codebase either); the underlying `repostatus.Status`
  struct IS exhaustively tested.

## One more finding, caught by CI rather than review — disclosed for completeness

After this review's findings 1–3 were fixed and the PR (#663) was opened,
the `Go build & test` CI job failed on a check this review hadn't run
locally: the dispatch-boundary guard (INVARIANTS.md T-002 / eng/go §2.18 —
`cmd_*.go` CLI wrapper files must not import `os`/`net/http`/
`encoding/json`/etc. directly; that logic belongs in the domain package).
`cli/cmd_repo_status.go` imported `encoding/json` directly for the
`--json` output path. Fixed by adding `(*repostatus.Status).JSON()` and
having the CLI wrapper call it instead. This is disclosed here rather than
silently folded into the commit history because it is exactly the kind of
gap a code-only review (however adversarial) can miss when it doesn't also
run every CI-gate script locally — recorded honestly as "caught by CI,
not by this review."

## Verdict rationale

All P1/P2/P3 acceptance criteria are met with real (not just fixture-only)
evidence, including a genuine end-to-end proof that a real `cn repo
install` output validates against the checked-in CUE schema. Three real
correctness bugs were found by adversarial review and self-check, plus one
dispatch-boundary convention violation caught by CI; all four are fixed
with targeted regression tests/re-verification (the three correctness bugs
each have a test proven to fail on pre-fix code; the boundary-convention
fix is verified by the CI job itself plus a local re-scan of every
`cmd_*.go` file for the same forbidden-import pattern). No known-but-
unfixed correctness gap remains in this cycle's scope. Converge.
