# β review — cnos#615

## §R0

**Reviewer:** β, independent review round 0. This section is derived from
an independent re-walk of the code/tests/diff — α's `self-coherence.md`
claims were read but not trusted; every AC below was re-derived from the
actual diff, the actual test bodies, and live re-runs of the gates, not
from α's narrative.

**Base / head:**
- `origin/main` at review time: re-fetched; `cycle/615` head `1c18c68`
  ("alpha: append review-readiness signal for cnos#615 R0") is confirmed
  to contain current `origin/main` (rebase claim in self-coherence.md
  holds — no drift).
- Diff reviewed: `git diff origin/main...cycle/615 --stat` → 6 files,
  `.cdd/unreleased/615/{gamma-scaffold.md,self-coherence.md}` (new),
  `src/packages/cnos.issues/commands/issues-fsm/{fetch.go,issuesfsm.go,terminal.go,terminal_test.go}`.
  `terminal.go` (322 lines), `terminal_test.go` (487 lines), the `fetch.go`
  addition (+36), and the `issuesfsm.go` wiring (+82/-8) were all read in
  full.

### AC1 — recognizes closed dispatch:cell+protocol:{P} with stale labels

**PASS.** `liveListClosedCandidates` (`terminal.go:270`) issues
`GET /repos/{repo}/issues?labels=dispatch:cell,protocol:{P}&state=closed`.
This query is structurally the recognizer: GitHub's own `labels=`/`state=`
filtering means a closed issue missing either `dispatch:cell` or
`protocol:{P}`, or an open issue regardless of labels, can never appear in
the response. Verified in `TestLiveListClosedCandidates_QueryShapeAndFiltering`
(`terminal_test.go:274`): captures the request URL, asserts
`state=closed` and `labels=dispatch:cell,protocol:cds` are present, and a
fake server honoring that filter returns exactly the one true candidate
(#950). At the orchestration layer, `TestTerminal_ClosedCompletedStale_RecognizedAndCleaned`'s
dry-run half confirms the candidate is reported (non-empty "would remove
... --apply" note) without `--apply`. `TestTerminal_OpenIssueNeverTouched`
confirms an open issue never appears in the results, run alongside a real
closed candidate so absence isn't merely "nothing to process."

### AC2 — correct label removal + resolution/* addition keyed off state_reason

**PASS.** `resolutionForStateReason` (`terminal.go:108`) maps exactly:
`"completed"` → `"resolution/completed"`, `"not_planned"` →
`"resolution/not-planned"`, anything else → `""` (skip, never guess).
This is the exact pinned mapping — re-verified independently against the
operator's clarifying comment, not just read from α's claim.
`TestTerminal_ClosedCompletedStale_RecognizedAndCleaned`'s apply half
asserts the exact removal set `{900:dispatch:cell, 900:protocol:cds,
900:status:review}` and exactly `900:resolution/completed` added — not
merely "a label was added." `TestTerminal_ClosedNotPlannedStale_ResolvesToNotPlanned`
explicitly asserts `resolution/not-planned` (with an inline comment
rejecting `resolution/completed` or `resolution/wontfix`) for a
`not_planned` fixture — this is the exact ambiguity the operator's
2026-07-07T01:51:47Z comment resolved, and the test pins it correctly.
`TestTerminal_LiveDefaults_RouteThroughSharedPrimitives` proves the same
shape through the live-default wiring (fake GitHub server, not injected
fakes): exactly 1 GET, 3 DELETEs, 2 POSTs.

**Live label verification (independently re-run, not trusted from
scaffold/self-coherence time):**
```
$ gh api repos/usurobor/cnos/labels/resolution%2Fcompleted
{"name":"resolution/completed","color":"ededed","description":null}
$ gh api repos/usurobor/cnos/labels/resolution%2Fnot-planned
404 Not Found
$ gh api repos/usurobor/cnos/labels/resolution%2Fwontfix
404 Not Found
```
`terminal.go`'s `resolutionLabelColor = "ededed"` / `resolutionLabelDescription
= ""` match `resolution/completed`'s live shape exactly (empty string +
`omitempty` on the JSON tag serializes to an absent field, which GitHub
treats as null on creation — matching `resolution/completed`'s
`description: null`). `resolution/not-planned` is confirmed not to exist
yet on the live repo, consistent with α's account; the creation path
(`ghEnsureLabelExists`) is unit-tested (`TestGhEnsureLabelExists_CreatesOnFirstUse`,
`TestGhEnsureLabelExists_TreatsAlreadyExists422AsNoOp`) but was not
exercised against the live repo by α or by β (β does not touch issue
labels, per role rules) — this is a disclosed, correct debt item, not a
gap.

### AC3 — idempotent; open/active untouched

**PASS.** `TestTerminal_Idempotent_SecondPassIsNoOp` (`terminal_test.go:191`):
first pass reconciles #902 (3 removes + 1 add), second pass's fixture
simulates the post-cleanup list (candidate excluded) → zero new mutation
calls, and `ghEnsureLabelExists` called exactly once total across both
passes. `TestTerminal_OpenIssueNeverTouched` independently confirms an
open issue is never a candidate and never mutated, run alongside a real
closed candidate in the same sweep — checked explicitly, not assumed from
AC1. Idempotence is structural (falls out of the list-query filter, no
bolted-on marker), matching the scaffold's own argument.

### AC4 — transitions.json untouched; removals go through the guarded path

**PASS**, independently re-run:
```
$ git diff origin/main...cycle/615 -- src/packages/cnos.cds/skills/cds/fsm/transitions.json | wc -l
0
$ git diff origin/main...cycle/615 -- src/packages/cnos.core/labels.json | wc -l
0
$ git diff origin/main...cycle/615 -- src/packages/cnos.issues/commands/issues-fsm/scan.go
(empty)
```
Grepped `terminal.go`/`issuesfsm.go`/`fetch.go` for
`gh issue close|state.*closed.*PATCH|"state":"closed"|http.MethodPatch` →
zero hits — no close-authority path introduced anywhere in the new code.
`TerminalOptions` carries no `Table` field (unlike `ScanOptions`);
`TestTerminalOptions_HasNoFSMTableDependency` pins that shape. All label
removal in `terminal.go` calls `opts.RemoveLabel`, whose live default
(`setDefaults()`) is `ghRemoveLabel` (`fetch.go:372`) reused verbatim — no
direct/duplicate `http.MethodDelete` call exists in `terminal.go` itself
(only the one inside `ghRemoveLabel`, which is unchanged pre-existing
code). The CLI wrapper `src/go/internal/cli/cmd_issues_fsm.go` has zero
diff (`git diff --stat` confirms) — it already delegates generically, no
change was needed, consistent with α's claim.

### AC5 — the 4 required fixture scenarios are present and test what they claim

**PASS.** Read every test body in `terminal_test.go` in full (not just
function names):
1. `TestTerminal_ClosedCompletedStale_RecognizedAndCleaned` — genuinely
   asserts dry-run non-mutation + apply's exact removal/add sets (see AC2).
2. `TestTerminal_ClosedNotPlannedStale_ResolvesToNotPlanned` — genuinely
   asserts `resolution/not-planned` specifically, with an explicit
   rejection comment for the two wrong mappings.
3. `TestTerminal_AlreadyClean_NoOp` — zero candidates in, zero mutation
   calls out; genuinely tests the "list-query-excludes" no-op case.
4. `TestTerminal_OpenIssueNeverTouched` — genuinely asserts an open issue
   never appears and is never mutated, alongside an in-flight closed
   candidate.

These are Go-native injected-fake tests (`TerminalOptions` seam), not
`testdata/*.json` files — a documented, reasonable divergence from
`scan_test.go`'s fixture-file style since this pass's testing seam
(no `Table`/`FactSnapshot`) doesn't need one. This is a legitimate
choice, not a shortfall — the four scenarios the scaffold names are all
present and all pass. Two additional defense-in-depth tests
(`TestTerminal_UnrecognizedStateReason_SkippedNeverMutated`,
`TestGhEnsureLabelExists_*`) exceed the minimum bar.

### AC6 — gates green

Independently re-run (not read from α's claims):
- `go build ./...`, `go vet ./...`, `go test ./... -v` from
  `src/packages/cnos.issues/commands/issues-fsm`: clean; full suite
  (pre-existing + 13 new top-level tests) passes, 0 failures.
- `go build ./...`, `go test ./...` from `src/go`: clean, `ok` for every
  package.
- `gofmt -l .` on the touched package: clean (no output).
- Built `/tmp/cn` from `src/go` and ran
  `cn cdd verify --unreleased --exceptions .cdd/exceptions.yml` from repo
  root: **182 passed, 0 failed, 119 warnings** — matches α's claimed
  result exactly (independently reproduced, not trusted).
- **Live GitHub Actions CI on the actual head commit** (`1c18c68`),
  checked via `gh run list --branch cycle/615` and
  `gh run view <id> --json jobs`: the `Build` workflow run for `1c18c68`
  is `success`, and every job is green: `Go build & test`, `Binary
  verification`, `Package verification`, `Package/source drift (I1)`,
  `Protocol contract schema sync (I2)`, `Repo link validation (I4)`,
  `SKILL.md frontmatter validation (I5)`, `CDD artifact ledger validation
  (I6)`, `Dispatch repair-preflight guard`, `Dispatch closeout-integrity
  guard`. This is a stronger verification than α could self-certify (α
  had no Actions runner available; β does, via `gh run view`) and it
  closes α's disclosed debt item ("branch CI is not independently
  confirmed green"). `install-wake-golden` did not trigger at all
  (`gh run list --workflow install-wake-golden.yml --branch cycle/615` →
  no runs), consistent with the path-filter reasoning in Friction note 3
  — this cell touches none of its filtered paths.
- Note: the one earlier commit on this branch
  (`d262684`, "alpha: self-coherence R0...") had a **failed** Build run
  (self-coherence.md header format issue, per α's own account) — this is
  expected in-flight history, not a live problem, since the immediately
  following commit fixed it and every commit since (including the final
  head) is green.

### Guardrail verification (independent of the AC table)

- `src/packages/cnos.core/labels.json`: confirmed zero-diff (see AC4).
- No `gh issue close`/close-equivalent introduced (see AC4 grep).
- `resolution/not-planned`'s pinned color/description independently
  re-verified live to match `resolution/completed`'s current values (see
  AC2) — not merely re-read from scaffold-time notes.
- Open design call (CLI shape: new `terminal` sub-verb vs `--closed` flag
  on `scan`) is documented in `self-coherence.md` §Self-check with a
  stated rationale (issue's own wording names `ghRemoveLabel` not
  `Evaluate`/`Table` as the reuse target; folding into `scanOne` would
  either fake an FSM state through `Evaluate()` or require an awkward
  early-return). α chose the new-sub-verb shape; confirmed in the diff
  that no `Evaluate`/`Table`/`transitions.json` coupling exists anywhere
  in the closed-issue path.

### Scope-creep / guardrail check

- No new `status:*` label introduced anywhere (`labels.json` untouched;
  `terminal.go`'s only `"status:` reference is the `strings.HasPrefix`
  scan used to *find and strip* existing `status:*` labels, not to add
  one).
- No new FSM terminal state added — `transitions.json` byte-identical.
- `scan.go` diff is empty — zero behavior change to the existing
  active-lifecycle scanner, and its full pre-existing test suite passes
  unchanged in the same `go test` run as the new tests.
- No close-authority code path added (see AC4).
- `resolution/duplicate` is correctly *not* special-cased — GitHub's
  `state_reason` enum has no `duplicate` value; Friction note 2's
  reachability argument holds and the code matches it (`default: return
  ""` in `resolutionForStateReason`, which for the actually-unreachable
  case skips rather than guesses — this is over-cautious, not a
  duplicate-handling gap).

## Findings

None. No correctness, scope, or guardrail issues survived independent
re-derivation. The one item worth naming explicitly (not a finding, a
note for γ's PRA): AC5's fixtures are Go-native injected-fake tests
rather than `testdata/*.json` files, a legitimate and disclosed
divergence from `scan_test.go`'s style, not a gap — flagging only so a
future reader auditing "does this cell have `testdata/`" doesn't
misread the absence as AC5 being unmet.

## Verdict

**converge**

All six ACs are genuinely met, independently re-derived from the code,
tests, and live gates rather than taken from α's self-coherence.md
narrative. The pinned `not_planned` → `resolution/not-planned` mapping
(the one correctness-critical detail the operator's clarifying comment
resolved) is correct in the code and covered by a test that explicitly
rejects the two wrong alternatives. `transitions.json` and
`cnos.core/labels.json` are byte-identical to `origin/main`. No
close-authority code path, no new FSM terminal state, no new `status:*`
label, and zero behavior change to `scan.go`/`evaluate`. Live GitHub
Actions CI on the actual head commit (`1c18c68`) is green across every
job the scaffold's AC6 oracle names, closing α's one disclosed debt item
(CI not independently confirmed). Ready for β merge / closeout.
