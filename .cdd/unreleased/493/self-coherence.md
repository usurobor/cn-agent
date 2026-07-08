§R0

## Gap

**Issue:** cnos#493 — "label-doctor: install + verify cnos.core canonical labels in repos" (live-discovery follow-up to cnos#491).

**Mode:** `design-and-build` (per issue header and γ scaffold).

**Branch:** `cycle/493`, based on `main@31d7ddfa53efb67f2d996b085f4d92986b585ef9` at γ scaffold time; rebased onto `origin/main@43b81295ff8172caceaf9c0795d03b5261a606ff` before this readiness signal (see §Review-readiness for the rebase record — `origin/main` had advanced via three automated `board-map` regeneration commits since scaffold time; the rebase was clean, no conflicts, diff unaffected).

**Incoherence being closed:** `src/packages/cnos.core/labels.json` (schema `cn.labels.v1`, 8 entries: 7 `status:*` lifecycle labels + `dispatch:cell`) declares the canonical cnos.core label set. The cnos repo SHOULD carry this full set after install, but no mechanism materializes or repairs it — labels get auto-created by GitHub with default colors/empty descriptions the first time an agent applies one (as happened with `status:review` during cnos#491), and nothing detects or corrects the drift. Two pieces of already-shipped runtime code assume this mechanism exists and name cnos#493 explicitly when it doesn't:

1. `src/go/internal/repoinstall/repoinstall.go`'s `ensureCanonicalDispatchLabels()` (the `cn repo install --dispatch cds` path) was a stub that unconditionally returned an error naming cnos#493.
2. `src/go/internal/cell/cell.go`'s `Returner.preflightTargetLabel()` (the `status:review → status:changes` review-return transition) has a runtime error whose fix-it text literally says "Run label-doctor before retrying" — a command that did not exist anywhere in the repo.

**What this cycle builds** (per the γ scaffold's `## α prompt` and 7-axis Implementation contract): a new Go module `src/packages/cnos.core/commands/label-doctor/` implementing a canonical-label audit+repair engine against `labels.json`, using dependency-free GitHub REST primitives (no `gh` CLI shellout, mirroring `cnos.issues/commands/issues-fsm/fetch.go`'s style); a thin `cn label-doctor` kernel subcommand wired per the dispatch-boundary convention; replacement of `ensureCanonicalDispatchLabels()`'s stub body with a real in-process call into the new package; `go.work` + CI wiring (module test step, CLI ergonomics smoke-test entry, a CI guard that actually exercises both "fails on drift" and "passes on clean" directions); and a committed live audit + repair artifact.

Full detail, the pinned 7-axis implementation contract, the per-AC oracle, and the source-of-truth table are in `.cdd/unreleased/493/gamma-scaffold.md` (γ's dispatch scaffold, read in full before any code was written).

## Skills

**Tier 1 (loaded):** `CDD.md` (canonical lifecycle/role contract); `cnos.cdd/skills/cdd/alpha/SKILL.md` (this role's own skill — read in full before coding, including §2.5 incremental self-coherence discipline, §2.6 pre-review gate, §3.6 implementation-contract constraint discipline).

**Tier 2 (always-applicable `eng/*`, per `src/packages/cnos.eng/skills/eng/README.md`):** `eng/go` (dispatch boundary — INVARIANTS.md T-002 — thin `cli/cmd_*.go` wrappers, domain logic in an imported package; go.work module co-location precedent).

**Tier 3 (issue-specific, resolved from the γ scaffold's Source-of-truth table rather than re-derived from scratch):**
- `src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` §3 ("Ownership split") + §5 ("Manifest") — labels.json is data-only; a consumer reads it and creates/repairs labels.
- Structural precedent: `src/packages/cnos.issues/commands/issues-fsm/` (Go-module co-location shape, `fetch.go`'s net/http idiom — `ghGetJSON`/`ghRequest`/`ghAddLabel`/`ghRemoveLabel`/`ghEnsureLabelExists`, `githubAPIBase` test-seam var, `resolveDefaultTablePath`'s directory-walk idiom, `withFakeGitHub` test helper) and `src/go/internal/cli/cmd_issues_fsm.go` (thin-wrapper shape: `Run` delegates entirely to the domain package, no `Help()` method, generic `--help` fallback via `PrintCommandHelp`).
- `src/go/internal/cli/dispatch.go`'s `ResolveCommand`/`GroupMembers`/`InvocationName` — read in full after an initial `cn label-doctor --help` smoke check surfaced that a hyphenated `Spec().Name` is *always* routed through the noun-verb space form (`cn label doctor`), never a flat single hyphenated token, by design ("Hyphenated flat forms ... are NOT supported" — dispatch.go's own doc comment). See §ACs AC4 for the full write-up; this is not an improvisation on the pinned CLI-integration axis, it is dispatch.go's pre-existing, documented behavior for the Name value the contract pinned.
- `cnos.cds/skills/cds/CDS.md` §"Coordination surfaces" / §"Artifact contract" (cycle-branch-as-coordination-surface, `.cdd/unreleased/{N}/` conventions).

## ACs

All AC-by-AC verification below was run against implementation SHA `a73e1f38` (the last implementation commit before this self-coherence.md's own commits — `git log --oneline origin/main..HEAD` at write time), against the real `usurobor/cnos` repo unless noted otherwise.

### AC1 — Audit complete

**Oracle (γ scaffold):** an artifact exists listing all 8 `labels.json` entries and their live state (present-and-matching / present-but-drifted / missing); re-running the audit gets the same classification, modulo repairs this PR itself made.

**Evidence:** `.cdd/unreleased/493/label-audit.md` (committed `a73e1f38`) — the "Before" section reproduces the γ scaffold's live-audit table exactly (7 of 8 drifted, `status:blocked` missing), captured via `cn label doctor --dry-run` run against the real repo before any repair. The "After" section shows the post-repair state (7 of 8 fully canonical; `dispatch:cell`'s color corrected, description blocked by a real GitHub API constraint — see below). The audit is generic by construction: `Audit()` (`src/packages/cnos.core/commands/label-doctor/doctor.go`) diffs every `manifest.Labels` entry against the live set by name/color/description with no per-label special-casing — `status:review` (the issue's originally-cited finding) is not treated differently from any other entry (`TestAudit_ClassifiesMatchDriftedMissing`, `TestAudit_ColorCaseInsensitive_DescriptionExact`).

**Status: MET.**

### AC2 — All canonical labels exist with canonical color+description

**Oracle:** `gh label list --repo usurobor/cnos --json name,color,description` shows all 8 entries present, color+description byte-equal (case-insensitive on hex) to `labels.json`.

**Evidence (live re-verification at write time):**

```
$ gh label list --repo usurobor/cnos --limit 200 --json name,color,description | jq -c '.[] | select(.name|test("^(status:|dispatch:cell)"))'
{"color":"ededed","description":"Well-formed scope but not yet refined to dispatch readiness.","name":"status:backlog"}
{"color":"c2e0c6","description":"Spec'd; ACs converged; awaiting operator authorization.","name":"status:ready"}
{"color":"0e8a16","description":"Operator-authorized; dispatch wake may claim.","name":"status:todo"}
{"color":"fbca04","description":"Claimed by a dispatch wake; cycle running.","name":"status:in-progress"}
{"color":"5319e7","description":"Cell complete; awaiting external/operator review of the receipt/result.","name":"status:review"}
{"color":"d93f0b","description":"External review requested changes; issue must be revised before re-dispatch.","name":"status:changes"}
{"color":"b60205","description":"Gated on external input (operator, infra, external authority).","name":"status:blocked"}
{"color":"1d76db","description":null,"name":"dispatch:cell"}
```

7 of 8 are byte-equal to `labels.json`. `dispatch:cell`'s color is now canonical (`1d76db`); its description is `null` (empty), not the canonical 149-byte string. This is not a bug in the shipped tool — GitHub's label API rejects any description over 100 characters (`HTTP 422: "description is too long (maximum is 100 characters)"`), and `labels.json`'s own `dispatch:cell.description` is 149 bytes. See `.cdd/unreleased/493/label-audit.md` §"Residual gap" for the full account and §Debt below.

**Status: MET for 7/8; PARTIALLY MET for `dispatch:cell`** (color canonical, description structurally inapplicable — genuine `labels.json` data defect, out of this issue's pinned scope to silently edit; see §Debt).

**Correction (resumption pass, 2026-07-08):** the above was accurate at original write time (implementation SHA `a73e1f38`) but is now stale. A later commit on this same branch, `0135b30b` ("fix over-100 dispatch:cell label description (GitHub API limit)"), shortened `labels.json`'s `dispatch:cell.description` from 149 bytes to 92 bytes (under GitHub's 100-character API limit), removing the "structurally inapplicable" blocker this AC's original text names. During this resumption pass, `cn label doctor --repo usurobor/cnos --dry-run` showed the live label still carried the old (empty) description — the manifest fix had landed but the live repair had not yet been re-applied — so `cn label doctor --repo usurobor/cnos` (no `--dry-run`) was run against the real repo, and a follow-up `--dry-run` confirmed all 8/8 canonical labels, including `dispatch:cell`, now report `match` (color `1d76db`, description byte-equal to `labels.json`'s new 92-byte string). Full transcript in `.cdd/unreleased/493/label-audit.md` §"Update (resumption pass, 2026-07-08): residual gap closed". **AC2 status is corrected to fully MET (8/8), not 7/8 partial.**

### AC3 — `status:review` corrected

**Oracle:** `status:review.color == "5319e7"`, checked independently of AC2.

**Evidence:** confirmed live above — `status:review` color is `5319e7`, description byte-equal to canonical. Also directly exercised by `TestGhUpdateLabel_PatchesColorAndDescription` (unit-level PATCH-body assertion) and `TestDoctor_Apply_RepairsDriftAndIsIdempotent` (end-to-end repair against a fixture carrying `status:review` drift).

**Status: MET.**

### AC4 — Command exists, idempotent

**Oracle:** (a) source directory exists under `src/packages/cnos.core/commands/`; (b) running it against drifted/missing labels produces the canonical set (fixture-based, no live network required for CI); (c) running it a second time performs zero mutating API calls.

**Evidence:**

- (a) `src/packages/cnos.core/commands/label-doctor/` exists, own `go.mod` (`module github.com/usurobor/cnos/packages/cnos.core/commands/label-doctor`, `go 1.24`), added to root `go.work`'s `use (...)` block (5th entry).
- (b) `TestDoctor_Apply_RepairsDriftAndIsIdempotent` and `TestRun_Apply_RepairsLiveState` (`doctor_test.go`, `cli_test.go`) drive a fixture with drifted + missing labels through an `httptest.Server` and assert the live-store mock ends up canonical. No live network call in any test (`withFakeGitHub` redirects `githubAPIBase`, mirroring `issues-fsm`'s idiom).
- (c) `TestDoctor_Apply_RepairsDriftAndIsIdempotent`'s second `Doctor()` call against the now-repaired fixture asserts `store.posts == 1 && store.patches == 1` (unchanged from the first call) and `len(res2.Applied) == 0` — idempotence is asserted on the fake server's actual request counts, not just "exit 0 twice" (per the γ scaffold's explicit oracle wording).
- **Command registration + CLI integration:** `src/go/internal/cli/cmd_label_doctor.go` (`LabelDoctorCmd`), registered in `src/go/cmd/cn/main.go` (`reg.Register(&cli.LabelDoctorCmd{})`, same block as `IssuesFsmCmd`/`RepoInstallCmd`). `Run` delegates entirely to `labeldoctor.Run(ctx, inv.Args, inv.Stdin, inv.Stdout, inv.Stderr)` — no `os`/`net/http`/`encoding/json`/`path/filepath` import in `cmd_label_doctor.go` (verified: the CI "Dispatch boundary check" grep pattern run locally against `internal/cli/cmd_*.go` reports zero violations, see §Self-check).
- **Invocation-shape finding (real, not improvised):** because `Spec().Name = "label-doctor"` contains a hyphen, `cli/dispatch.go`'s `ResolveCommand` routes it through the noun-verb path — `cn label doctor` (two argv tokens), not a flat single hyphenated token `cn label-doctor`. This was empirically confirmed (`./cn label-doctor --help` resolves to a *group listing*, not the command; `./cn label doctor --help` resolves correctly, exit 0) and traced to `dispatch.go`'s own documented invariant: "Hyphenated flat forms (e.g., 'kata-run') are NOT supported. The canonical form is noun-verb." This governs **every** existing hyphenated-Name kernel command identically (`issues-fsm` → `issues fsm`, `cdd-verify` → `cdd verify`, `repo-install` → `repo install`, `cell-finalize`/`cell-return`/`cell-resume` → space forms) — `label-doctor` is not a special case, and this is not α relaxing the pinned CLI-integration axis: the `Spec().Name` string is registered exactly as pinned (`"label-doctor"`), and the resulting invocation shape is a deterministic, pre-existing consequence of that value passing through unmodified shared dispatch code. The contract's "single-token, matching `cn doctor`/`cn build`/`cn status` naming" framing was an inexact analogy (those three are genuinely flat, non-hyphenated Names) alongside its own, more specific and load-bearing instruction to "mirror `cmd_issues_fsm.go` exactly" — and `issues-fsm` is itself only reachable via its noun-verb space form, not a flat token. The two exemplars conflict; this cycle resolved the conflict in favor of the more specific, structurally-detailed instruction rather than pausing for a γ/δ round-trip over a downstream, fully-explained, zero-latitude mechanical consequence (see `cmd_label_doctor.go`'s doc comment for the same write-up at the code site). `cell.go`'s existing runtime error text ("Run label-doctor before retrying") uses the hyphenated form as operator-facing prose/mnemonic, not a literal invocable token — flagged in §Debt.
- CLI ergonomics smoke test (`.github/workflows/build.yml`) now includes `"label doctor"` in its `--help` exercise loop and `"label-doctor"` in both `cn help`/`cn status` leaked-hyphenated-name negative checks (peer enumeration of the existing family: `issues-fsm|cell-finalize|cell-return|cell-resume|repo-install|cdd-verify|issues-map`).
- **`ensureCanonicalDispatchLabels()` stub replaced, not paralleled:** `src/go/internal/repoinstall/repoinstall.go` now calls `labeldoctor.Doctor(ctx, labeldoctor.Options{RepoRoot: opts.RepoRoot, Stdout: opts.Stdout, Stderr: opts.Stderr})` in-process (both packages are `go.work`-linked; import mirrors `cmd_issues_fsm.go`'s cross-module `issuesfsm` import). No parallel/dead path remains — `grep -n "cnos#493 label-install mechanism is not yet available" src/go/internal/repoinstall/repoinstall.go` returns zero hits (verified).
- **Existing "cnos#493" stub-string tests updated, not deleted/weakened:** `repoinstall_test.go`'s `TestRun_DispatchCds_RendersWorkflow_ThenSurfacesLabelGap` and `TestRun_DispatchCds_SigmaDefault_NoIdentityFlagsRequired`, and `cmd_repo_install_test.go`'s `TestRepoInstall_DispatchCds_IdentityFlagsWireThrough`, now assert on `labeldoctor`'s real target-resolution failure (`"canonical dispatch labels not ensured"` + `"could not resolve target repo"`) instead of the literal `"cnos#493"` stub string — and each test explicitly asserts `"cnos#493"` no longer appears in that error, proving the stub is gone, not just relabeled. These three fixtures' `RepoRoot`/`repoDir` genuinely have no configured git remote (a bare `t.TempDir()`, or `initGitRepo` which never adds "origin"), so this exercises `labeldoctor`'s real resolution-failure path with **zero live network calls** — deterministic and CI-safe.

**Status: MET.**

### AC5 — CI guard

**Oracle:** a CI job/step runs the tool in dry-run/diff mode; fails on injected/fixture drift, passes on clean state; both directions actually exercised via test fixtures, not just described.

**Evidence:**

- `.github/workflows/build.yml`'s new `"Test label-doctor module (cnos#493 AC4/AC5)"` step (`working-directory: src/packages/cnos.core/commands/label-doctor`, `go test ./... -v`) mirrors the existing `"Test issues-fsm module (cnos#568 AC9)"` step exactly, for the same reason (a `go.work` sibling module is not covered by `src/go`'s own `go test ./...`).
- **Both directions are exercised through the actual `Run()` CLI entry point** (the same code path `cn label-doctor --dry-run` runs in production), not just the lower-level `Doctor`/`Audit` functions:
  - `TestRun_DryRun_FailsOnDrift` (`cli_test.go`): fixture with injected color drift + a missing label → `Run(..., []string{"--token","tok","--dry-run"}, ...)` returns a non-nil error, stderr names "drift detected", zero mutating calls on the fake server.
  - `TestRun_DryRun_PassesOnClean` (`cli_test.go`): fixture whose live labels already match `labels.json` exactly → `Run(...)` returns nil, zero mutating calls.
  - (Lower-level equivalents `TestDoctor_DryRun_FailsOnDrift`/`TestDoctor_DryRun_PassesOnClean` in `doctor_test.go` cover the same two directions one layer down, against `Doctor` directly.)
- Every CI run of this step exercises both directions unconditionally (they're both in the default `go test ./...` set) — a regression in either direction turns CI red on the next push, satisfying "a guard that only ever passes is not a guard."

**Status: MET.**

## Self-check

This section is written as part of a **resumption** (per `alpha/SKILL.md`
§4 "Resumption"): a prior α run implemented all 5 ACs and wrote `## Gap`
/ `## Skills` / `## ACs` but died mid-write before `## Self-check` /
`## Debt` / `## CDD Trace`, leaving CI's "CDD artifact ledger validation
(I6)" check red (`self-coherence.md sections — missing required
sections: CDD Trace Self-check or Debt`). This resumption re-verified
every AC claim against real command runs rather than trusting the
existing prose, per the dispatch's explicit instruction not to
uncritically re-read the prior claims.

**Did this cycle's work push ambiguity onto β?**

- **No new ambiguity was introduced.** All 5 ACs were independently
  re-run against real state during this resumption pass (not merely
  re-read): `cd src/packages/cnos.core/commands/label-doctor && go test
  ./... -v` → 29/29 tests pass (`PASS`, `ok
  github.com/usurobor/cnos/packages/cnos.core/commands/label-doctor
  0.078s`); `cd src/go && go build ./... && go test ./...` → builds
  clean, 15/15 packages `ok`; `gofmt -l` against every `.go` file in
  `git diff --name-only origin/main...HEAD` → zero output (clean);
  `go.work` confirmed to list the new module as its 5th `use` entry;
  `python3 -c "import yaml; yaml.safe_load(open('.github/workflows/build.yml'))"`
  → `YAML OK`; `./cn label doctor --help` (noun-verb form) exits 0,
  `./cn label-doctor --help` (flat form) resolves to the `label` group
  listing, exactly matching the existing self-coherence.md's AC4 claim
  about `dispatch.go`'s hyphenated-Name routing — independently
  reproduced, not assumed.
- **One claim was found stale and was corrected, not silently
  overwritten.** AC2's original text (written before commit `0135b30b`
  landed on this same branch) said `dispatch:cell`'s canonical
  description was "structurally inapplicable" due to GitHub's 100-char
  API limit. Live re-verification during this resumption found the
  underlying blocker had already been fixed on-branch (description
  shortened to 92 bytes) but never re-applied to the live
  `usurobor/cnos` repo. This resumption ran `cn label doctor --repo
  usurobor/cnos` (apply, not dry-run) against the real repo and
  confirmed via a follow-up `--dry-run` that all 8/8 canonical labels
  now report `match`. Both `self-coherence.md` (AC2's own subsection)
  and `label-audit.md` (the AC1 evidence artifact) carry an explicit
  correction note (§2.3 intra-doc-repetition rule — grepped both sites
  carrying the stale "149-byte"/"structurally inapplicable" claim) rather
  than a silent rewrite, so β can see exactly what changed and why. This
  means the branch now hands β a **fully-met** 8/8 AC2/AC3 state rather
  than a 7/8-partial one the prior claim would have left for β to either
  accept as "good enough" or independently discover was stale — either
  outcome would have been ambiguity pushed onto β.
- **Rebase disclosure.** `origin/main` had advanced by 22 commits since
  the γ scaffold's base SHA (`31d7ddfa5...`) — considerably more than
  "board-map regeneration only." `git rebase origin/main` completed
  cleanly with zero conflicts (`Successfully rebased and updated
  refs/heads/cycle/493`); the diff content is unchanged, only every
  commit's SHA and the base moved. Force-pushed with `--force-with-lease`.
  No β-facing ambiguity here since the rebase was mechanical and the
  post-rebase diff was re-diffed against `origin/main` to confirm
  identical file list/content (see §CDD Trace step 3).
- **Caller-path / dead-code check.** `labeldoctor.Doctor` is called from
  two live call sites, not zero: (1) `src/go/internal/cli/cmd_label_doctor.go`
  (`LabelDoctorCmd.Run` → `labeldoctor.Run`, the `cn label doctor` CLI
  entry point) and (2) `src/go/internal/repoinstall/repoinstall.go`'s
  `ensureCanonicalDispatchLabels()` (the `cn repo install --dispatch cds`
  path) — confirmed by `grep -rn "labeldoctor\.\(Run\|Doctor\)" src/go/internal`.
  Neither is a new module left uncalled.
- **Is every claim backed by evidence I just re-ran?** Yes for every AC
  row above and for this section's own claims — each cites either a
  command actually run during this resumption pass (with its actual
  output quoted or paraphrased faithfully) or a specific file/line/grep
  result. No claim in this section rests on "the prior self-coherence.md
  said so."

## Debt
