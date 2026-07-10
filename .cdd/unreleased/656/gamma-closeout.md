# γ close-out — cycle #656

**Scope of this artifact.** Converge-boundary close-out, written after β's
`converge` verdict, before the cycle-PR merges. Issue #656 remains `OPEN`
at this point; full release-time closure (§2.10) happens post-merge.

---

## Cycle summary

**Issue:** [cnos#656](https://github.com/usurobor/cnos/issues/656) — Phase 1
of the [cnos#655](https://github.com/usurobor/cnos/issues/655) `cn repo`
lifecycle wave.

**What shipped:** see `alpha-closeout.md` for the shipped-surface list;
`self-coherence.md` for the full design rationale (5 numbered decisions);
`beta-review.md` for the three bugs found and fixed during review.

**Review arc:** single round (R0 → converge). No iterate-round back to α —
all findings (1 self-check + 2 independent-adversarial-review) were fixed
and re-verified within the same review pass, each with a regression test
proven to fail against the pre-fix code.

---

## Process-gap audit

**Was this "review working as intended" or a scaffold gap?** Working as
intended, and unusually cleanly so: this cycle ran as a single interactive
session acting as δ (routing γ/α/β itself, per `delta/SKILL.md` §9), which
means the α/β role separation had to be enforced deliberately rather than
structurally (no separate subagent invocation forces independence). Two
mitigations were used to keep the review genuinely independent rather than
a rubber stamp: (1) α's own pre-review critical re-read, which caught the
`DriftRemoved` gap before requesting review at all; (2) a dedicated
adversarial-review subagent, briefed only on the issue/diff — not on α's
self-coherence narrative — which caught two further real bugs. Both
mechanisms are recorded explicitly in `beta-review.md` rather than
presented as if a separate human β reviewed blind, because that would
misrepresent how this cycle's role separation actually worked.

**No `gamma/SKILL.md` §2.8 cycle-iteration trigger fired:** one review
round, three findings total (well under a mechanical-overload floor), all
substantive (not tooling/environment flakes), all fixed within the same
pass.

**No new process-gap issue filed this cycle.** Unlike cnos#608's γ
close-out (which found and filed a genuine `gamma/SKILL.md` scaffold gap),
this cycle's findings were implementation bugs in α's own code, not gaps
in the CDD process/skill documents themselves — nothing here generalizes
to a skill-doctrine fix the way #608's did.

---

## Follow-up issues

Checked existing trackers first (`gh issue list --search "repo doctor"`,
`--search "cn repo lifecycle"`) — all six wave phases already have
dedicated sub-issues (#656 this one, #657 Phase 2, #658 Phase 3, #659
Phase 4, #660 Phase 5, #661 Phase 6 docs). No new issue filed:

| Debt item | Disposition |
|---|---|
| `cn repo doctor` named in the design doc's command table but not claimed by any of the six filed sub-issues | **Not filed** — speculative which future phase should own it; flagged for operator/planner awareness in `beta-closeout.md` rather than unilaterally scoping a new phase. |
| `sha256File` duplicated between `internal/repoinstall` and `internal/repostatus` | **Not filed** — ~8 lines each, two call sites; not worth a shared package yet per this repo's own DRY discipline (3+ callers threshold). |
| `deps.lock.json`-present-but-`deps.json`-absent package drift not surfaced | **Not filed** — narrow, requires a hand-edited manifest to trigger, out of this phase's literal Acceptance text. |
| `cmd_repo_status.go` presentation layer untested line-by-line | **Not filed** — consistent with `hubstatus.Run`'s own precedent elsewhere in this codebase. |
| A3 backfill scoped to `install` only (`update`/`repair` don't exist yet) | **Not filed** — `update`/`repair` are #658/#659's own territory; they inherit this cycle's `internal/repostate`/ledger-write pattern when built. |

**Total new issues this cycle: 0.** Every debt item is either already
covered by an existing filed sub-issue's future scope, or narrow/
speculative enough that filing now would be premature.

---

## Receipt parity contract (A6)

```yaml
mock_parity:
  source: docs/development/design/cn-repo-status-MOCKS.md
  source_commit: "e0dbf509 (the commit that authored this file, pre-implementation)"
  rows:
    - id: A1
      expectation: "reports source channel/release/index"
      observed: "Status.Source (Source{Channel,Release,Index}); printed as '✓ source: <channel> @ <release> (<index>)' or '✓ source: <channel> (<index>)' when release is empty (index-only installs)"
      evidence: "internal/repostatus/repostatus.go Run(); TestRun_CleanInstall_NoDrift; manual smoke test transcript in self-coherence.md"
      verdict: match
      how: "matches the mocked A1 shape exactly"
    - id: A2
      expectation: "reports each package's desired vs locked vs vendored version and an in_sync verdict"
      observed: "PackageStatus{Name,Desired,Locked,Vendored,InSync}"
      evidence: "packageStatuses(); TestRun_CleanInstall_NoDrift, TestRun_LocalEditDetected"
      verdict: match
      how: "1:1 field match with Mock C's packages[] shape"
    - id: A3
      expectation: "reports dispatch presence, tier, id, path, and drift classification"
      observed: "DispatchStatus{Present,Tier,ID,Path,Drift}; Drift takes 5 values (matches_ledger/user_edit/renderer_moved/removed/unclassified), exceeding the mocked 4-way (matches_ledger/user_edit/renderer_moved/absent)"
      evidence: "dispatchStatus(); TestDispatchStatus_MatchesLedger/RendererMoved/UserEdit/Removed/Unclassified"
      verdict: exceed
      how: "the two extra values (removed, unclassified) were added during β review to close real gaps a mocked-only 4-way split would have missed (a ledger-recorded-but-deleted workflow, and a confirmed-but-unclassifiable drift) — both are strictly more informative, never masking the mocked cases, and both count toward the top-level Drift bool so they can't silently degrade to 'looks clean'"
    - id: A4
      expectation: "reports canonical label status: ok / missing / unknown"
      observed: "LabelStatus{Status,Missing,Drifted,Unknown} — adds Drifted (canonical labels present with wrong color/description) beyond the mocked 3-field shape"
      evidence: "labelStatus(); requires a small additive label-doctor change (Result.LiveLabels) to compute Unknown"
      verdict: exceed
      how: "Drifted is a real, distinct category the mock's illustrative JSON didn't separately name but the underlying label-doctor Audit() already computes (StatusDrifted) — surfacing it is strictly more informative and doesn't conflict with the mocked fields"
    - id: A5
      expectation: "reports orphan vendored packages (materialized but absent from deps.lock.json)"
      observed: "Status.OrphanPackages []string"
      evidence: "packageStatuses() orphan detection; TestRun_OrphanVendoredPackage"
      verdict: match
      how: "matches the mocked shape and the design doc's own definition (absent from LOCK, not from manifest)"
    - id: A6
      expectation: "reports locally-edited managed files, distinguishing user_edit from ledger-vs-fresh-render drift"
      observed: "Status.LocalEdits []LocalEdit{Path,Classification}; workflow-kind entries excluded (reported via Dispatch instead, matching the mock's own example which never double-lists the workflow file here)"
      evidence: "localEditStatus(); TestRun_LocalEditDetected, TestRun_RemovedManagedFile"
      verdict: match
      how: "matches Mock B's local_edits shape and the no-duplication behavior the mock's own example implies"
    - id: A7
      expectation: "reports whether a newer release is available, read-only, degrading offline"
      observed: "UpdateAvailable{Checked,Available,Release}"
      evidence: "updateAvailableStatus(); manual smoke test showed 'unknown (network check skipped or failed)' in the sandbox's actual offline conditions"
      verdict: match
      how: "matches the mocked shape; Checked:false degrade path exercised live, not just in a unit test"
    - id: B1
      expectation: "--json emits schema cn.repo.status.v1, machine-parseable"
      observed: "StatusSchema = \"cn.repo.status.v1\"; json.MarshalIndent(st, \"\", \"  \")"
      evidence: "cmd_repo_status.go Run(); manual --json smoke test output"
      verdict: match
      how: "exact schema string match; 2-space indent matches the repo's one consistent --json/file-write convention"
    - id: B2
      expectation: "--json sets top-level drift: true iff any surface reports non-clean"
      observed: "Status.Drift bool, computed in Run() from packages/dispatch/labels/orphans/local_edits"
      evidence: "Run()'s drift aggregation; TestRun_LocalEditDetected, TestRun_OrphanVendoredPackage, TestDispatchStatus_* all assert Drift"
      verdict: match
      how: "every non-clean surface (including the two new dispatch drift values) is wired into the aggregation — verified by dedicated tests per surface, not just inspection"
    - id: C1
      expectation: "--check exits nonzero iff drift == true; exits 0 otherwise, always 0 without --check"
      observed: "RepoStatusCmd.Run() returns a non-nil error iff checkMode && st.Drift"
      evidence: "cmd_repo_status.go; manual smoke test: clean repo --check exit=0, hand-edited repo --check exit=1, hand-edited repo WITHOUT --check exit=0"
      verdict: match
      how: "all three cases verified live against the built binary, not just unit-tested"
    - id: D1
      expectation: "cn repo status never writes .cn/repo.state.json or any other file, with or without any flag"
      observed: "zero write calls reachable from repostatus.Run outside a throwaway os.MkdirTemp cleaned up via defer os.RemoveAll"
      evidence: "TestRun_NeverWrites; manual git status --porcelain before/after smoke test across default/--json/--check"
      verdict: match
      how: "verified both by a dedicated Go test walking every file under RepoRoot before/after 3 consecutive Run calls, AND by a live git-status diff on a real installed repo"
    - id: D2
      expectation: "when the ledger is absent, status still runs and reports ledger.present:false, degraded but non-empty reporting"
      observed: "Ledger{Present,Reconstructed}; Packages/Dispatch still populated directly from deps.json/deps.lock.json/vendor/the live workflow file when the ledger is missing"
      evidence: "TestRun_NoLedger_DegradesGracefully"
      verdict: match
      how: "matches the mocked degrade behavior; local_edits correctly reports empty (nothing to compare against without the ledger) rather than erroring"
    - id: E1
      expectation: "help text explicitly distinguishes cn repo status from the existing global cn status"
      observed: "repoStatusHelp const's DESCRIPTION section states this explicitly"
      evidence: "cmd_repo_status.go repoStatusHelp"
      verdict: match
      how: "verbatim distinguishing sentence present in --help output"
  summary:
    matched: 11
    exceeded: 2
    missed: 0
    exceed_justified: true
```

---

## δ closeout-integrity preflight — deliverable_evidence

```yaml
deliverable_evidence:
  pr: "#663 (cycle/656 -> main)"
  head_sha: "c06f93618b9a8ba11719b06b31ae56c08d36ed8a"
  base_sha: "e7bf83ca033c40078abd41f0eb0af3817aceb4cd"
  commits_beyond_base: 9
  closeout_artifacts: [gamma-scaffold.md, self-coherence.md, beta-review.md, alpha-closeout.md, beta-closeout.md, gamma-closeout.md]
  ci: "all 12 required checks pass at this head SHA (gh pr checks 663), including the new repo.state.json schema validation (cnos#656) job and the pre-existing Go build & test job (which caught and required the T-002 dispatch-boundary fix folded into this head)"
```

All five deliverable-evidence conditions satisfied: PR #663 exists and
references #656; PR has 9 commits beyond base; `cycle/656` exists and
differs from base; all six required closeout artifacts present (this
file included); this block names the PR number and head/base SHAs.

---

**Cycle #656's converge-boundary close-out is complete.** Full release-time
closure (§2.10 — `RELEASE.md`, the `.cdd/unreleased/656/` →
`.cdd/releases/{X.Y.Z}/656/` move, issue-close assertion) is deferred to
the post-merge γ invocation.
