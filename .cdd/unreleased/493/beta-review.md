<!--
section-manifest:
  planned: [Verdict, Contract Integrity, Issue Contract, Diff Context, Architecture, Findings]
  completed: [Verdict, Contract Integrity, Issue Contract, Diff Context, Architecture, Findings]
-->

# β Review — cnos#493 (R0)

## Verdict

**verdict: converge**

All 5 ACs independently re-verified against real command output (not against α's self-coherence.md prose) — live `gh label list` re-run, module tests re-run from a clean shell (32/32 pass), the built `cn` binary exercised directly (`cn label doctor --help`), CI re-pulled via `gh api`/`gh run list`. Implementation-contract conformance (Role Rule 7) confirmed across all 7 pinned axes. Scope guardrails (no `cnos.cds`/`cnos.cdr`/`cnos.cdw` touch, no dispatch-protocol doctrine rewrite, `internal/cell/cell.go` untouched) hold. No D-severity findings.

Base at this review: `origin/main` re-fetched at review time = `dd9693ca`'s parent lineage unchanged since α's rebase (α's Review-readiness signal cites `origin/main@90e9c8b2`; `git merge-base origin/main cycle/493` reconfirmed identical — no new commits landed on `main` between α's signal and this review). Reviewed head: `dd9693ca40a8b8d007652df170379584a7ac511f` ("self-coherence.md §Review-readiness — round 1, ready for β").

## Contract Integrity

**γ artifact presence** — `.cdd/unreleased/493/gamma-scaffold.md` present on `origin/cycle/493` (221 lines, committed `a703808a`). ✓ (pre-merge-gate row 4 precondition satisfied, informationally — β does not execute the merge gate in this wake-invoked-δ cycle, per this cycle's dispatch instructions; row 4's *artifact-completeness* substance still applies and is verified here.)

**Role Rule 7 — implementation-contract conformance.** Verified each of the 7 pinned axes against the diff independently:

| Axis | Pinned value | β's independent check | Result |
|---|---|---|---|
| Language | Go 1.24, no shell script | `git diff --name-only main..HEAD` shows only `.go`/`.json`/`.md`/`.yml` files; all new domain logic in `src/packages/cnos.core/commands/label-doctor/*.go`. `gofmt -l` on every changed `.go` file → clean. `go vet ./...` in both the new module and `src/go` → clean. | conforms |
| CLI integration | New kernel subcommand `cn label-doctor`, thin wrapper mirroring `cmd_issues_fsm.go`, `NeedsHub: false` | Rebuilt `/tmp/cn-b` from this worktree; `/tmp/cn-b label doctor --help` → exit 0, prints the correct summary; `/tmp/cn-b label-doctor --help` (flat form) → resolves to the `label` group listing, not the command — confirmed this is `dispatch.go`'s pre-existing, documented behavior for any hyphenated `Spec().Name` (same as `issues-fsm`→`issues fsm`), not an α improvisation. Read `cmd_label_doctor.go` directly: imports only `context` and the `labeldoctor` package — no `os`/`net/http`/`encoding/json`/`path/filepath`. `Spec()` omits `NeedsHub` (zero value `false`), matching the pinned hub-less requirement — same idiom `cmd_activate.go`/`cmd_help.go` use for their `false` case. | conforms |
| Package scoping | `src/packages/cnos.core/commands/label-doctor/`, own `go.mod`, 5th `go.work` `use` entry, new `build.yml` module-test step | `go.work` read directly: `use (./src/go ./src/packages/cnos.cdd/commands/cdd-verify ./src/packages/cnos.issues/commands/issues-map ./src/packages/cnos.issues/commands/issues-fsm ./src/packages/cnos.core/commands/label-doctor)` — 5th entry present. `build.yml` diff shows a new `"Test label-doctor module (cnos#493 AC4/AC5)"` step, `working-directory: src/packages/cnos.core/commands/label-doctor`, mirroring the `issues-fsm` step exactly. Domain logic is NOT in `src/go/internal/` (confirmed: `src/go/internal/repoinstall/repoinstall.go` only imports `labeldoctor`, contains no label-audit logic itself). | conforms |
| Existing-binary disposition | Stub replaced (not paralleled); in-process call, not subprocess | Read `repoinstall.go`'s diff directly: `ensureCanonicalDispatchLabels()`'s stub body (`return fmt.Errorf("canonical dispatch labels not ensured: cnos#493 ...")`) is gone; `grep -n "cnos#493 label-install mechanism is not yet available" src/go/internal/repoinstall/repoinstall.go` → 0 hits (independently re-ran, not trusted from self-coherence.md). New body calls `labeldoctor.Doctor(ctx, labeldoctor.Options{...})` in-process — no `exec.Command` to a `cn-label-doctor` binary anywhere in `repoinstall.go`. `internal/cell/cell.go` is byte-identical to `main` (`git diff main..HEAD -- src/go/internal/cell/cell.go` → empty) — the `gh`-shellout label-transition path was not touched, confirming the pinned "separate subsystem, do not touch" instruction was honored. | conforms |
| Runtime dependencies | Go stdlib only (`net/http`, `encoding/json`, `os`) + `$GITHUB_TOKEN`/`$GH_TOKEN`; no `gh` CLI shellout | Read `github.go` directly: `ghRequest`/`ghListLabels`/`ghCreateLabel`/`ghUpdateLabel` are all `net/http` + `encoding/json`, no `gh` CLI invocation anywhere. `resolve.go` shells out to `git remote get-url origin` via `os/exec` — this is *not* a `gh` CLI call and is exactly the mechanism the scaffold's own friction note 2 and the Implementation-contract row itself directed α to build ("α needs a small git-remote-parsing utility... parse `git remote get-url origin`"); not a contract deviation. | conforms |
| JSON/wire contract | N/A for GitHub's label API; `labels.json` schema (`cn.labels.v1`) unchanged, data-only edit permitted | `git diff main..HEAD -- src/packages/cnos.core/labels.json` shows only `dispatch:cell.description` changed (149 bytes → 92 bytes, to fit GitHub's 100-char API cap) — no schema/shape change, no other field touched. Independently confirmed the new value is byte-identical to what's live on GitHub (`gh label list` output, see AC2 below). | conforms |
| Backward compat | `--dispatch none` unaffected; existing `cnos#493`-string tests updated (not deleted) to reflect the new non-stub behavior | `ensureCanonicalDispatchLabels` is called only from `runDispatchCds` (the `--dispatch cds` branch) — confirmed by reading `repoinstall.go`'s call graph directly; `--dispatch none`'s code path contains zero references to `labeldoctor`. `git diff main..HEAD -- repoinstall_test.go cmd_repo_install_test.go \| grep -n cnos#493` shows the tests were rewritten to assert the string is now ABSENT from the error (`if strings.Contains(err.Error(), "cnos#493") { t.Errorf(...) }`), not deleted or weakened — this is the correct backward-compat-test-churn shape the contract explicitly permits. | conforms |

All 7 axes conform. No `implementation-contract` finding.

## Issue Contract

Cross-checked cnos#493's issue body + 2 comments (`gh issue view 493 --json body,comments`) against the γ scaffold and the diff: the 5 ACs, the "Constraints/non-goals" (in-scope / out-of-scope) sections, and the operator follow-up comment's required-behavior bullets (creates missing labels, updates drifted color/description idempotently, hub-less, callable by the installer path, dry-run/preview, tests against fixture labels, never silently falls back) all match verbatim what the scaffold restates and what α executed against. No drift between issue, scaffold, and diff.

Mode = `design-and-build`, correctly honored: γ pinned the command name (`cn label-doctor`) and the open design point (git-remote resolution utility) rather than leaving them for α to improvise from scratch, and α resolved the remaining design surface (drift-repair PATCH primitive, per-label-failure continuation) and disclosed both bugs found live (§Debt / `label-audit.md`) rather than silently patching over them.

## Diff Context

```
$ git diff main..HEAD --stat
 .cdd/unreleased/493/CLAIM-REQUEST.yml                              |  18 +
 .cdd/unreleased/493/gamma-scaffold.md                              | 221 +++++++++
 .cdd/unreleased/493/label-audit.md                                 | 182 +++++++
 .cdd/unreleased/493/self-coherence.md                              | 434 +++++++++++++++++
 .github/workflows/build.yml                                        |  24 +-
 docs/development/board/board-data.json                              |   2 +-
 docs/development/board/index.html                                  |   2 +-
 docs/guides/INSTALL-CDS.md                                          |  19 +-
 go.work                                                             |   1 +
 src/go/cmd/cn/main.go                                               |   1 +
 src/go/internal/cli/cmd_label_doctor.go                            |  47 ++
 src/go/internal/cli/cmd_repo_install.go                            |  17 +-
 src/go/internal/cli/cmd_repo_install_test.go                       |  24 +-
 src/go/internal/repoinstall/repoinstall.go                         |  72 ++--
 src/go/internal/repoinstall/repoinstall_test.go                    |  66 ++--
 src/packages/cnos.core/commands/label-doctor/cli.go                |  68 +++
 src/packages/cnos.core/commands/label-doctor/cli_test.go           | 137 +++++
 src/packages/cnos.core/commands/label-doctor/doctor.go             | 269 +++++++++
 src/packages/cnos.core/commands/label-doctor/doctor_test.go        | 344 ++++++++++++
 src/packages/cnos.core/commands/label-doctor/github.go             | 172 ++++++
 src/packages/cnos.core/commands/label-doctor/github_test.go        | 148 ++++++
 src/packages/cnos.core/commands/label-doctor/go.mod                |   3 +
 src/packages/cnos.core/commands/label-doctor/manifest.go           |  91 ++++
 src/packages/cnos.core/commands/label-doctor/manifest_test.go      | 141 ++++++
 src/packages/cnos.core/commands/label-doctor/resolve.go            |  45 ++
 src/packages/cnos.core/commands/label-doctor/resolve_test.go       |  86 ++++
 src/packages/cnos.core/labels.json                                 |   2 +-
 .../cnos.core/skills/agent/label-doctrine/SKILL.md                 |   2 +-
 28 files changed, 2561 insertions(+), 77 deletions(-)
```

`docs/development/board/*` are auto-regenerated board files carried by the branch's rebase onto `origin/main` (unrelated to this cycle's own work — confirmed these two files are pure regeneration artifacts, same as α's self-coherence.md CDD Trace disclosed). Every other file maps to an AC row or the implementation contract's package-scoping axis; no file is unaccounted for.

Author identity: all cycle commits carry `sigma <sigma@cnos.cn-sigma.cnos>` — this is the explicit, disclosed bootstrap/single-operator exemption (§Debt item 6 of self-coherence.md), consistent with this review session's own instructed identity handling; not a finding.

## Architecture

### AC1 — Audit complete

**MET, independently re-verified.** Re-ran `gh label list --repo usurobor/cnos --limit 200 --json name,color,description` myself and cross-referenced against `src/packages/cnos.core/labels.json`'s 8 entries: all 8 present, all byte-equal (color case-insensitive) to `labels.json` — this **is** the "after" state `label-audit.md`'s committed artifact claims. Read `label-audit.md` directly (not paraphrased from self-coherence.md): its "Before" table reproduces the γ scaffold's live-audit table exactly (7/8 drifted, `status:blocked` missing), captured via a real `cn label doctor --dry-run` run, and its "After" table shows the post-repair state. Read `doctor.go`'s `Audit()` function directly: it diffs every `manifest.Labels` entry generically by name/color/description with no per-label special-casing (`status:review` is not singled out in the code) — confirmed by reading the loop body, not inferred from a test name.

### AC2/AC3 — All canonical labels exist with canonical color+description; status:review specifically at 5319e7

**MET (8/8), independently re-verified live.** Re-ran the audit myself (not copy-pasting α's transcript):

```
$ gh label list --repo usurobor/cnos --limit 200 --json name,color,description | jq -c 'sort_by(.name)[] | select(.name|test("^(status:|dispatch:cell)"))'
```

produced all 8 entries with color+description byte-equal to `src/packages/cnos.core/labels.json` (case-insensitive on hex):

- `status:backlog` ededed / "Well-formed scope but not yet refined to dispatch readiness." — match
- `status:ready` c2e0c6 / "Spec'd; ACs converged; awaiting operator authorization." — match
- `status:todo` 0e8a16 / "Operator-authorized; dispatch wake may claim." — match
- `status:in-progress` fbca04 / "Claimed by a dispatch wake; cycle running." — match
- `status:review` **5319e7** / "Cell complete; awaiting external/operator review of the receipt/result." — match (AC3 checked explicitly, not folded into AC2's pass)
- `status:changes` d93f0b / "External review requested changes; issue must be revised before re-dispatch." — match
- `status:blocked` b60205 / "Gated on external input (operator, infra, external authority)." — match
- `dispatch:cell` 1d76db / "Executable cell; claimable by a dispatch wake with status:todo and a matching protocol:{id}." — match (this required the `labels.json` 149→92-byte description shortening; independently verified new text is 92 bytes with `python3 -c "print(len(...))"`, under GitHub's 100-char label-description cap, and is byte-identical to the live label)

All 8/8 fully canonical. This directly resolves the residual gap self-coherence.md's own §Self-check/§Debit disclosed as fixed during the resumption pass — independently confirmed correct, not accepted on trust. `TestGhUpdateLabel_PatchesColorAndDescription` (unit PATCH-body assertion) and `TestDoctor_Apply_RepairsDriftAndIsIdempotent` (end-to-end fixture repair, asserting the drifted `status:review` fixture reaches `5319e7`) independently re-run: both pass.

### AC4 — Command exists, idempotent

**MET, independently re-verified against real code and a rebuilt binary.**

- **(a) Module exists.** `src/packages/cnos.core/commands/label-doctor/` present, own `go.mod` (`module github.com/usurobor/cnos/packages/cnos.core/commands/label-doctor`, `go 1.24`), 5th `go.work` `use` entry confirmed by direct read.
- **(b) Fixture-based tests, no live network.** Ran `cd src/packages/cnos.core/commands/label-doctor && go test ./... -v` myself from a clean shell: **32/32 PASS, 0 FAIL** (exact count matches self-coherence.md's claimed 32 — independently reproduced, not trusted). Read `cli_test.go`/`doctor_test.go` directly: `withFakeGitHub` redirects `githubAPIBase` to an `httptest.Server`, mirroring `issues-fsm`'s idiom — zero live network calls in any test.
- **(c) Idempotence asserted on request counts, not just exit codes.** Read `TestDoctor_Apply_RepairsDriftAndIsIdempotent` directly: first `Doctor()` call asserts `store.posts == 1 && store.patches == 1` (1 create + 1 repair); second call against the now-repaired fixture asserts `len(res2.Applied) == 0` **and** `store.posts != 1 || store.patches != 1` fails the test if either counter moved — i.e., the second call is asserted to make **zero additional mutating HTTP requests** on the fake server, exactly the oracle wording the γ scaffold and β prompt both demanded ("assert on request count against the fake server, not just 'exit 0 both times'"). This is not a superficial exit-code check.
- **Rebuilt binary, live-exercised.** `cd src/go && go build -o /tmp/cn-b ./cmd/cn` — clean build. `/tmp/cn-b label doctor --help` → exit 0, correct summary text. `/tmp/cn-b label-doctor --help` (flat, hyphenated form) → resolves to the `label` group listing, not the command — independently reproduces the exact invocation-shape finding self-coherence.md names, confirming it is real `dispatch.go` behavior, not a misdescribed claim.
- **Stub replaced, not paralleled.** Confirmed above under Contract Integrity (Existing-binary-disposition row).
- **`cnos#493`-string tests updated, not deleted.** Confirmed above under Contract Integrity (Backward-compat row) — `repoinstall_test.go`/`cmd_repo_install_test.go` now assert the string's **absence**, a stronger check than mere deletion.
- **`src/go` regression check.** `cd src/go && go build ./... && go test ./...` — clean build, all 15 package suites (`internal/activate` through `internal/restore`, including `internal/cell` and `internal/repoinstall`) report `ok`. No regression.

### AC5 — CI guard

**MET, independently re-verified both directions through the real `Run()` entry point.**

Read `cli_test.go` directly:
- `TestRun_DryRun_FailsOnDrift`: fixture with injected color drift (`status:review` at `ededed`) + a missing label (`status:blocked`) → `Run(ctx, []string{"--token","tok","--dry-run"}, ...)` asserted to return non-nil error, stderr contains `"drift detected"`, and `store.mutatingCalls() == 0`.
- `TestRun_DryRun_PassesOnClean`: fixture whose live labels already match canonical → `Run(...)` asserted to return nil, `store.mutatingCalls() == 0`.

Both drive the actual `Run()` CLI entry point (the same path `cn label-doctor --dry-run` executes in CI), not just the lower-level `Doctor`/`Audit` functions — confirmed by reading the test bodies, not accepting the claim. `.github/workflows/build.yml`'s diff confirmed directly: new `"Test label-doctor module (cnos#493 AC4/AC5)"` step (`working-directory: src/packages/cnos.core/commands/label-doctor`, `go test ./... -v`), mirroring the `issues-fsm` module-test-step pattern for the same not-covered-by-root-`go test ./...` reason. `"label doctor"` added to the CLI-ergonomics `--help` exercise loop; `"label-doctor"` added to both hyphenated-name-leak negative-check regexes (`cn help` and `cn status`). Both directions are unconditionally in the default `go test ./...` set, so a regression in either turns CI red on the next push — a guard that only ever passes is not a guard, and this one demonstrably isn't.

**Live CI, independently re-pulled** (not trusted from self-coherence.md's transcript): `gh run list --repo usurobor/cnos --branch cycle/493` and `gh api repos/usurobor/cnos/commits/dd9693ca.../check-runs` both confirm the review head (`dd9693ca`) is green on every job — `Go build & test`, `Package/source drift (I1)`, `Protocol contract schema sync (I2)`, `Repo link validation (I4)`, `SKILL.md frontmatter validation (I5)`, `CDD artifact ledger validation (I6)`, `Dispatch repair-preflight guard (cnos#516)`, `Dispatch closeout-integrity guard (cnos#524)`, `Binary verification`, `Package verification` — all `success`, both `push` and `pull_request` triggers. PR #634 (`usurobor/cnos`, still `OPEN`) shows the identical head SHA and the same all-`SUCCESS` `statusCheckRollup`.

### Scope guardrails

Confirmed independently, not accepted from self-coherence.md's own scope-check claim:
- `git diff --name-only main..HEAD | grep -iE "cnos\.cds|cnos\.cdr|cnos\.cdw|dispatch-protocol"` → **0 hits**. No protocol-qualifier label install touched, no dispatch-protocol doctrine rewrite.
- `git diff main..HEAD -- src/go/internal/cell/cell.go` → **empty**. The `gh`-shellout-based `preflightTargetLabel`/`applyLabelTransition` path is byte-identical to `main`; not refactored onto the new REST primitives, and its own runtime error text ("Run label-doctor before retrying") was left as-is — correctly flagged in α's §Debt item 5 as a cosmetic prose gap rather than silently patched (patching it would itself have been an out-of-scope edit to a file the guardrails name as untouchable).
- No new agent-admin-wake bundling anywhere in the diff (confirmed by the file list above — no `cds-dispatch`/wake-orchestrator files touched).

## Findings

None blocking.

**Sanity-checked, not blocking:**

1. **`resolve.go` shells out to the `git` binary via `os/exec`** (`git remote get-url origin`) rather than parsing `.git/config` directly. This is a process-exec dependency the Implementation-contract's Runtime-dependencies row's headline sentence ("None beyond Go stdlib... no `gh` CLI shellout") could be read to exclude at first glance. Checked against the row's own body text, though: the contract explicitly instructs α to "parse `git remote get-url origin`" for this exact utility — the exclusion is scoped to `gh` CLI shellouts specifically (the competing pattern in `cell.go`), not to `git` itself, which every checkout in this repo already depends on. Not a contract violation; noting only so a future reader of the row's headline sentence in isolation doesn't misread it as broader than the row's own detail intends.
2. **`labels.json`'s `dispatch:cell` description was edited** (149→92 bytes) to fit GitHub's 100-character API cap — a data change, not named in the issue's original "Constraints/non-goals" as an explicit in-scope action, though it was a necessary, disclosed, minimal fix without which AC2 could never reach 8/8 (GitHub's API hard-rejects the original value). Fully disclosed in `label-audit.md`'s "Residual gap"/"Update" sections and self-coherence.md's §Debt item 1, with an accurate before/after byte count. Correctly scoped and transparent; not a finding requiring action.

Neither item is a D-severity implementation-contract or scope violation; both are pre-existing disclosures independently confirmed accurate rather than new discoveries. Ready for γ closeout.
