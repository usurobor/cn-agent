# self-coherence — cycle #656

**manifest:** sections = [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness]
**completed:** [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness]

---

## Gap

**Issue:** [cnos#656](https://github.com/usurobor/cnos/issues/656) — "repo lifecycle Phase 1 — state ledger (.cn/repo.state.json) + cn repo status"

**Mode:** implied `design-and-build` (schema design + CLI implementation; no explicit `**Mode:**` field in the issue body).
**cell_kind:** `implementation` (per `gamma-scaffold.md`'s frontmatter).

**Parent:** cnos#655 (the `cn repo` lifecycle wave master tracker). This cycle
closes Phase 1 only — the foundation phase the wave's later removal commands
(`update`/`repair`/`prune`/`uninstall`) depend on. Those later phases,
`cn repo doctor`, and the A7 install-semantics two-step flip are explicitly
out of scope (see `gamma-scaffold.md`'s "Non-goals" section).

**Gap being closed:** `cn repo install` (cnos#607/#608/#609/#610, shipped in
3.82.x) writes packages/lock/vendor/dispatch-workflow state but keeps no
durable **ownership ledger** — there is no machine-readable record of which
files/dirs/workflows CNOS itself manages versus what a repo owner added.
Without this, a future `cn repo uninstall` cannot safely distinguish
CNOS-managed surfaces from user matter, and there is no `cn repo status` to
answer "is this repo's install in sync, and with what?"

**What closes it:**
- `schemas/repo_state.cue` — a checked-in, closed CUE schema for
  `cn.repo.state.v1`, validated by `scripts/ci/validate-repo-state.sh` (a new
  CI job, `repo-state-schema-check`, mirroring the existing
  `skill-frontmatter-check` job's `cue vet` pattern). Closedness makes two
  invariants structural rather than script-enforced: timestamp-freeness (no
  timestamp field is declared) and no lock-duplication (no top-level
  `packages` key — design doc A1).
- `internal/repostate` — pure Go types (`RepoState`, `Source`, `ManagedFile`,
  `ManagedDir`, `ExternalExpectations`) matching the CUE schema's field names
  1:1, with deterministic `Marshal()` (path-sorted, trailing newline, 2-space
  indent — same convention as `.cn/deps.json`).
- `internal/repoinstall`: `cn repo install` now writes/updates
  `.cn/repo.state.json` at the end of every non-dry-run install, computing
  every `managed_files` sha256 from the file it just wrote (read back from
  disk, never recomputed from in-memory state). Every install — including a
  rerun on a pre-cnos#656 repo with no ledger yet — writes/refreshes the
  ledger (A3's "backfill on next install", scoped to the one lifecycle-write
  command that exists today).
- `internal/dispatchrender` (new, extracted out of `repoinstall`'s
  `runDispatchCds`) — the single place that knows how to invoke the vendored
  `cn-install-wake` renderer, used by both the original install-time render
  and `cn repo status`'s fresh-render drift comparison (P2/A2), so the two
  call sites cannot drift apart on argument-building logic.
- `internal/repostatus` + `cli.RepoStatusCmd` (`cn repo status`, `--json`,
  `--check`) — read-only (P3: never writes `.cn/`, verified by a dedicated
  `TestRun_NeverWrites` test and a manual `git status --porcelain`
  before/after smoke test). Reports per-package desired/locked/vendored
  in-sync status, dispatch-workflow drift classified three ways
  (`matches_ledger` / `user_edit` / `renderer_moved` — A2, via a real
  fresh-render-and-compare using `internal/dispatchrender`), canonical-label
  status (`ok`/`drifted`/`unknown`, via a small additive `label-doctor`
  change exposing `LiveLabels` so "unknown/non-canonical" labels are
  computable), orphan vendored packages, locally-edited managed files, and a
  best-effort update-available check. Degrades gracefully (never fails
  outright) when the ledger is absent (A3) or labels/network can't be
  resolved.
- `docs/development/design/cn-repo-status-MOCKS.md` (new) — mocked
  human-visible `cn repo status` output authored before implementation,
  mirroring the install wave's own `cn-repo-install-MOCKS.md` mock-first +
  receipt-parity discipline (A6). See the `mock_parity` block in
  `gamma-closeout.md`.

**Design surface:** `docs/development/design/cn-repo-lifecycle.md` — §Phase-
contract precision P1–P3 is the authoritative AC surface (the "core decision"
section's illustrative ledger JSON predates that tightening and is superseded
where the two disagree — see `gamma-scaffold.md` Decision 1).

**Base SHA (branch creation):** `59cd36e6` (main tip at claim time).

---

## Skills

**Tier 1 (lifecycle):** `CDD.md`, `cnos.cdd/skills/cdd/delta/SKILL.md` §9
(dispatch-wake-invoked mode — this cycle was run by an interactive session
acting as the `cds-dispatch` wake per `cnos.cds/orchestrators/cds-dispatch/
SKILL.md`, routing γ/α/β itself rather than three separately-invoked
sub-agents).

**Tier 2 (protocol):** `cnos.cds/skills/cds/SKILL.md` — CDS lifecycle/label
discipline; `dispatch-protocol/SKILL.md` — claim mechanics (§Claim mechanism
followed literally: scan → verify gates → CLAIM-REQUEST.yml → FSM `--apply`
→ claim comment → post-claim re-read).

**Tier 3 (task-local):** `docs/development/design/cn-repo-lifecycle.md`,
`docs/development/design/cn-repo-install-MOCKS.md` (precedent mirrored),
`schemas/cdd/README.md` (§"Architectural choice" — the "keep the proof in
CUE" doctrine this cycle applies to closedness instead of required-key
unification), `scripts/ci/validate-skill-frontmatter.sh` (the script-owns-
discovery / schema-owns-shape split mirrored for `validate-repo-state.sh`).

---

## ACs

Mapped against the issue's own Acceptance section, verbatim IDs.

**P1 — durable schema + validation gate.**
- Schema file checked in: `schemas/repo_state.cue`. ✅
- CI gate rejects a ledger with timestamps: `invalid/timestamp-field.json`
  fixture, rejected by `cue vet` via CUE closedness (no timestamp field
  declared). Verified locally: `cue vet -c -d '#RepoState' schemas/repo_state.cue
  schemas/fixtures/repo-state/invalid/timestamp-field.json` exits nonzero
  naming `installed_at: field not allowed`. ✅
- CI gate rejects duplicated package version/sha under a `packages` key:
  `invalid/duplicated-package-lock.json` fixture, rejected the same way
  (`packages: field not allowed`). ✅
- `managed_files`/workflow record missing required fields rejected:
  `invalid/managed-file-missing-fields.json` (missing id/sha256) and
  `invalid/workflow-missing-render-contract.json` (workflow kind missing
  tier/renderer/etc.) both rejected with field-specific `incomplete value`
  diagnostics. ✅
- Valid ledger passes: `valid/base-install.json` and
  `valid/with-dispatch-workflow.json` both vet clean. ✅ Additionally
  verified against a REAL `cn repo install` output (not just hand-written
  fixtures) — see Self-check below.
- Interpretation call: `managed_files` requires `id`/`sha256` on every entry
  (not just workflow kind) — see `gamma-scaffold.md` Decision 1 for why this
  reads the issue's own P1 text over the design doc's earlier illustration.

**P2 — workflow records carry the render contract.**
- `path`/`kind`/`id`/`tier`/`sha256` plus `renderer`/`renderer_package`/
  `renderer_version_source`/`agent`/`workflow_pat_secret`(name only)/
  `bot_name`/`bot_id` — all present on workflow-kind `managed_files` entries,
  enforced by CUE's `if kind == "workflow" { ... }` conditional-required
  block. `workflow_pat_secret` never carries a secret *value* — only the
  GitHub Actions secret *name* (same discipline as the existing
  `--workflow-pat-secret` flag). ✅
- Round-trips enough to re-render: proven literally, not just by shape —
  `cn repo status`'s A2 drift classification re-invokes the SAME renderer
  (`internal/dispatchrender.Render`) with the ledger's recorded
  tier/agent/renderer_package/etc. and compares the fresh output's sha256
  against the live file. `TestDispatchStatus_RendererMoved` and
  `TestDispatchStatus_UserEdit` both exercise the real vendored
  `cn-install-wake` script end-to-end (not a stub). ✅
- Drift classified matches-ledger / user-edit / renderer-moved:
  `internal/repostatus.dispatchStatus` implements exactly this three-way
  split; `differs-from-ledger-but-equals-fresh-render` is NOT flagged as a
  local edit (`TestDispatchStatus_RendererMoved`). ✅

**P3 — `cn repo status` read-only by default.**
- Never writes `.cn/repo.state.json` (or anything else), with or without any
  flag: `internal/repostatus.Run` performs zero `os.WriteFile`/`os.Mkdir`
  calls anywhere in its call graph (only `os.ReadFile`/`os.ReadDir`/
  `os.MkdirTemp` for the throwaway fresh-render comparison, which is cleaned
  up via `defer os.RemoveAll`). `TestRun_NeverWrites` diffs every file under
  RepoRoot before/after three consecutive `Run` calls. A manual smoke test
  additionally confirmed `git status --porcelain` is byte-identical before
  and after `cn repo status`, `--json`, and `--check`. ✅
- May emit a reconstructed ledger candidate under `--json` when the ledger
  is absent: `Status.Ledger.{Present,Reconstructed}` + degraded (but still
  populated from deps.json/deps.lock.json/vendor directly) `Packages`/
  `Dispatch` reporting — `TestRun_NoLedger_DegradesGracefully`. ✅

**Status reports packages/lock/vendor/workflows/labels/orphans/drift;
`--json` machine-readable; `--check` nonzero on drift.**
- All seven surfaces present in `repostatus.Status`; `--json` uses
  `json.MarshalIndent(v, "", "  ")` (the repo's one consistent convention —
  see α's research findings in the claiming session, no prior `--json`
  precedent existed elsewhere in this codebase, so this cycle establishes
  it). `--check` returns a non-nil error (nonzero exit) iff `Drift == true`;
  confirmed by a live smoke test (`.cn/deps.json` hand-edited → `--check`
  exits 1; clean repo → exits 0). ✅

**Install writes a deterministic ledger; same-input rerun → no diff.**
- `TestRun_WritesRepoState_Deterministic` runs `Run` twice against the same
  fixture server and diffs the two `.cn/repo.state.json` outputs
  byte-for-byte. `RepoState.Marshal()`'s path-sort makes managed_files/
  managed_dirs ORDER deterministic independent of Go map/slice build order
  (`TestMarshal_DeterministicAcrossBuildOrder` in `internal/repostate`). ✅

**A6 — mock_parity, `missed: 0`.**
- `docs/development/design/cn-repo-status-MOCKS.md` authored (Mocks A–D,
  invariants A1–A7/B1–B2/C1/D1–D2/E1) before `cn repo status` was
  implemented. `gamma-closeout.md` carries the parity block; see that file
  for the row-by-row observed-vs-mocked comparison.

---

## Self-check

**Real (not fixture-only) proof the schema matches what `cn repo install`
actually writes:** a throwaway test dumped a real `Run`-produced
`.cn/repo.state.json` to `/tmp/real-ledger.json` and ran
`cue vet -c -d '#RepoState' schemas/repo_state.cue /tmp/real-ledger.json` —
PASS. This is stronger evidence than hand-written fixtures alone: the Go
struct tags and the CUE schema were checked against the SAME artifact the
production code path emits, not two independently-authored shapes that
happen to agree on paper.

**Full test suite, both Go modules:**
```
$ go build ./... && go vet ./... && go test ./...
ok  	.../internal/cli
ok  	.../internal/repoinstall
ok  	.../internal/repostate
ok  	.../internal/repostatus
ok  	.../internal/... (all others unaffected)
```
```
$ cd src/packages/cnos.core/commands/label-doctor && go build ./... && go vet ./... && go test ./...
ok
```

**CUE self-test:**
```
$ ./scripts/ci/validate-repo-state.sh --self-test
✓ all checks passed.
```

**Manual end-to-end smoke** (real `cn build` → real local index → real
`cn repo install` → real `cn repo status`, no mocks): confirmed clean-install
no-drift report, confirmed local-edit detection + `--check` exit codes
(0 clean / 1 drifted) after hand-editing `.cn/deps.json`, confirmed zero
filesystem writes from `cn repo status` via `git status --porcelain`
before/after.

**Existing-behavior regression check:** `internal/repoinstall`'s and
`internal/cli`'s full pre-existing test suites (dispatch-cds identity flags,
dry-run, idempotence, label-doctor gap surfacing, sparse-checkout exclusion,
SIGMA leak-audit, etc.) all still pass unmodified except
`TestRepoInstall_Idempotent_NoDiffOnSecondRun`, which now also stages
`.cn/repo.state.json` (a legitimate update: this cycle adds a new file
`cn repo install` writes, so the test's git-add list needed the same
addition, not a weakening of the assertion).

**A genuine bug found and fixed during self-review (not by β — see Debt for
what IS still open):** the first draft of `writeRepoState`'s call site lost
the ledger entirely whenever `--dispatch cds`'s render succeeded but the
downstream cnos#493 label-doctor step failed (a common real-world shape —
label-doctor needs a resolvable git remote + GitHub token, which most test
fixtures and many first-time installs lack). Fixed by capturing
`runDispatchCds`'s error instead of returning immediately, writing the
ledger whenever the workflow file actually exists on disk regardless of the
label error, and only skipping the ledger write when the render itself never
produced a file (the AC2 "no partial render" precondition-failure path).
Two new tests pin both branches:
`TestRun_WritesRepoState_DispatchWorkflowRenderContract` (ledger written
despite the label failure) and
`TestRun_DispatchCds_MissingIdentity_FailsEarlyNoPartialWrite`'s new
assertion (no ledger when there was never a render).

---

## Debt

- **`cn repo doctor`** is named in the design doc's Phase 1 command-surface
  table but not required by the issue's own Acceptance section; deferred to
  a later phase (recorded as an explicit scope call in `gamma-scaffold.md`,
  not an oversight).
- **A3 backfill is install-only.** The design doc says "backfill on next
  install/update/repair"; `update`/`repair` do not exist in this repo yet,
  so backfill is wired into the one command that exists. The next phase
  that adds `update`/`repair` should call the same ledger-write path this
  cycle added to `internal/repoinstall`.
- **`sha256File` is duplicated** (near-identically) between
  `internal/repoinstall` and `internal/repostatus` (~8 lines each). Not
  worth a shared package for a two-line SHA-256-of-a-file helper at this
  scale; flagged for β/future β to weigh in on whether hoisting to
  `internal/repostate` (which already has file-format knowledge) is worth
  it once a third caller appears.
- **Update-available check is best-effort only.** No retry, no rate-limit
  awareness beyond `binupdate.FetchLatestRelease`'s existing behavior;
  degrades to `Checked: false` on any error. Acceptable for an informational
  field per the design doc, but a future phase might want this to also
  respect `--offline`.
- **Label "unknown" (non-canonical) detection required a small upstream
  change** to `label-doctor` (`Result.LiveLabels`, additive, non-breaking).
  This is the one change this cycle made outside its own new packages;
  flagged explicitly here for β's attention even though it's low-risk
  (new field, no existing caller reads/needs it removed, both label-doctor
  test suites still pass).
- **CLI-level rendering (`renderStatus`) is not exhaustively tested** — only
  exercised via the manual smoke test, not a Go test asserting exact stdout
  strings. The underlying `repostatus.Status` struct IS exhaustively tested;
  the presentation layer is thin (~50 lines of `fmt.Fprintf`) and mirrors
  `hubstatus.Run`'s own precedent of not unit-testing print formatting
  line-by-line.

---

## CDD Trace

- γ scaffold: `.cdd/unreleased/656/gamma-scaffold.md` (5 numbered design
  decisions + scope/non-goals + AC mapping).
- α (this file): schema + CI gate (P1) → ledger write in install (P2/A3) →
  dispatchrender extraction → cn repo status (P3) → self-review fix
  (ledger-survives-label-failure) → this self-coherence.md.
- No β/α iteration rounds this cycle (`run_class: first_pass`, single pass
  R0) — the "self-review" fix described above was caught by α re-reading its
  own work before requesting β review, not by a separate β round.

---

## Review-readiness

- Build/vet/test clean across both affected Go modules.
- CUE self-test clean; real-ledger-against-schema proof captured above.
- Manual smoke test performed (install → status → hand-edit → status →
  --check exit codes → no-write verification).
- `docs/development/design/cn-repo-status-MOCKS.md` authored pre-
  implementation per A6; `gamma-closeout.md` carries the parity block.
- Known gaps recorded in Debt above, none blocking for Phase 1's own scope.
- Ready for β review.
