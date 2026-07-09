# β review — cnos#640

## §R0

**Base SHA:** `origin/main` at `2eeb40325463186abf76c38e961f840deceb9a43` (re-fetched synchronously before computing this review's diff base; matches the SHA α's self-coherence.md §Gap and §Review-readiness record — no further drift on `main` since α's readiness signal).
**Cycle-branch head SHA:** `b468ba7f38ae4a2ed1073c746175b97493c0b275`.
**Diff reviewed:** `git diff origin/main...cycle/640` (20 files changed, 1555 insertions(+), 2 deletions(-)).

This review re-derives independently against `gamma-scaffold.md` §1's oracle list. α's `self-coherence.md` framing is treated as hypothesis, not evidence — every claim below was re-checked against the live diff, a re-run test suite, or a re-run grep, not restated from α's artifact.

### AC1 — dispatching a design-first issue can no longer leave body-hold and label state contradictory

**Verdict: PASS.**

- `TestDispatch_ContradictoryShape_BodyCorrectedLabelUntouched` (`src/packages/cnos.issues/commands/issues-dispatch/dispatch_test.go:99`) builds `contradictoryBody` (`reconcile_test.go:9`) — I checked this literal string against the actual #640 issue body I fetched via `gh issue view 640 --json body` and it reconstructs the live phrase verbatim: `"**Not dispatched** — \`status:ready\`. Design-first; dispatch on explicit operator authorization."` preceded by `---`. The fixture pairs this body with `status:todo` (the drifted contradiction, not the "about to dispatch" case) and asserts `BodyChanged==true`, exactly 1 `EditBody` call, zero label-mutation calls, and the corrected body text.
- `TestDispatch_AlreadyCleanAndTodo_SafeNoOp` (line 165) proves the negative case: a clean body + `status:todo` produces `NoOp==true` and zero mutation calls of any kind.
- `TestDispatch_Idempotent_SecondRunIsNoOp` (line 221) proves running `Dispatch` twice against the same stateful fixture converges: first run applies both actions, second run is `NoOp==true` against the post-first-run state.
- `TestDispatch_BodyEditFails_LabelNeverFlipped` (line 261) proves the "no half-applied contradiction" property: I read `dispatch.go:144-160` myself — the body edit runs before the label flip by construction, and an injected `EditBody` failure returns before any label-mutation call, leaving `status:ready` untouched (never flipped to `status:todo` with the hold phrase still present). This is the exact contradiction shape the primitive exists to prevent, and the test proves it cannot be self-inflicted by a partial failure.
- All 14 tests in the package pass: I ran `go test ./src/packages/cnos.issues/commands/issues-dispatch/... -v` myself; output shows 14 `--- PASS` lines, 0 `--- FAIL`. This is an automated unit-test suite against an injectable fake (`fakeIssueStore`) and a local `httptest.Server` (`TestLiveDefaults_RouteThroughSharedPrimitives`), not a prose claim or a manual live-network check — satisfies the oracle's explicit requirement.
- Candidate (3) (contradiction gate) was not built as a separate mechanism; α discloses this explicitly in §Debt with reasoning (the detection logic is reused inside `Dispatch` itself; a standalone scanning gate is a materially different, out-of-scope shape). γ's scaffold marked candidate (3) optional/cheap-if-convenient, not mandatory. I accept α's judgment call here — it is disclosed, not silently dropped, and does not gate AC1's oracle (the oracle's mandatory bullet is the fixture pair, which is satisfied; the gate bullet is conditional on "if built").

### AC2 — the wake no longer has to decide "body says held, label says go"

**Verdict: PASS.**

- I re-read `dispatch-protocol/SKILL.md` §1.2 myself (diff hunk, not α's description): the new paragraph "Labels are the sole source of truth for dispatch readiness (cnos#640)" states, in one place, that the `status:*` label alone determines dispatch readiness, that body-hold prose is non-authoritative, and that a wake/δ finding a contradiction MUST NOT defer. Two ❌/✅ example pairs reinforce this. A pointer to `cn issues dispatch` follows.
- `label-doctrine/SKILL.md` §7 gets a one-line cross-reference only, correctly avoiding restating the doctrine in two places (would itself reintroduce a two-source-of-truth pattern for the *doctrine text*, ironically) — I confirmed this is genuinely one line, not a second copy of the substantive paragraph.
- New D13 failure-mode entry in `dispatch-protocol/SKILL.md` §5 names the symptom, cites #614→#633, points at §1.2's fix. α's correction of γ's friction note 2 (D8/D9 were not actually free slots — I grepped the live pre-edit file myself via `git show origin/main:.../dispatch-protocol/SKILL.md | grep -n "D8\|D9"` and confirmed both were already taken; D13 is the correct next-available slot) checks out.
- `delta/SKILL.md` §9.2 input #1 (line 401): I read it directly — it restricts δ to the body's named fields (Mode/Gap/Impact/Source-of-truth/Constraints/ACs/Proof/Cross-references) and states "δ does NOT improvise additions." It never told δ to treat a stray body sentence as authority-bearing in the first place, so no edit is required for consistency. α's negative-result claim (no edit needed) is correct; I confirm independently rather than take α's word for it.
- No standalone contradiction gate was built, so the "≥2 negative fixtures" sub-bullet of AC2's oracle (which is conditional on "if a contradiction gate is built") is not applicable — correctly not claimed by α.

### AC3 — the observe → capture → detect-recurrence → mechanize → verify-non-recurrence loop is documented

**Verdict: PASS.**

- New section `## Process self-improvement loop (cnos#640)` in `CELL-KINDS.md` (diff hunk read directly) names all five steps in a table, each with a "What happens" and "Owner today" column: Observe (no single owner, incidental), Capture (γ, cites #614/#630), Detect recurrence (**explicitly "No current owner,"** quotes κ's comment verbatim: "nobody/nothing actually runs it," names "ε made concrete" as next step), Mechanize (whoever files/dispatches the promoted issue; MCA bar per Kernel §1.2/§2.3), Verify non-recurrence (whoever ships the mechanizing cell; cites this cell's own `dispatch_test.go` as worked example).
- I cross-checked κ's actual comment via `gh issue view 640 --json comments` myself (not trusting α's quoting) — the comment's proposed loop, its 5-numbered breakdown ("observe / capture / detect recurrence / mechanize / verify non-recurrence"), and its "nobody owns the follow-through" framing match what CELL-KINDS.md now documents. No fabrication or over-claim found.
- Negative oracle: the closing paragraph ("Honest state as of cnos#640") explicitly states detect-recurrence has no running owner and that #640 "mechanizes one specific recurring failure... directly — it does not close the general detect-recurrence gap," naming the false-closure failure mode by name. This is an honest, non-inflated closure statement.
- Cross-linking: I checked `CDD.md`'s `CELL-KINDS.md` pointer entry (diff hunk) — it now names the new section explicitly with a link to #640. I checked `CELL-KINDS.md`'s own "Cross-references" list (diff hunk) — two new entries point at `dispatch-protocol/SKILL.md` §1.2/D13 and `ROLES.md` §4b (ε's role doctrine). Not orphaned.

### AC4 — all gates green; no weakening of the operator's genuine ability to hold work design-first

**Verdict: PASS.**

- **`transitions.json` hard check (byte-identical):** I ran `git diff origin/main -- src/packages/cnos.cds/skills/cds/fsm/transitions.json` myself. Output is empty (0 lines, confirmed via `wc -l` → 0). Byte-identical, as required.
- **CI status:** α's self-coherence.md recorded branch CI as "unknown, not asserted" (no Actions runner reachable from α's sandbox). I checked live: `gh run list --branch cycle/640` shows the HEAD commit (`b468ba7f`)'s "Build" workflow run (`28995180325`) is **green**, and `gh api repos/usurobor/cnos/commits/b468ba7f.../check-runs` confirms all 10 check runs succeeded (CDD artifact ledger, Protocol contract schema sync, Go build & test, Package/source drift, Repo link validation, both dispatch guards, SKILL.md frontmatter, Binary verification, Package verification). Note: several *intermediate* commits during α's incremental `self-coherence.md` authoring (`§Gap` through `§ACs`) show CI **failures** — I checked one (`28994947842`) and confirmed the failure is the CDD-verify gate correctly flagging an in-progress `self-coherence.md` missing required sections (an expected, self-resolving artifact of incremental multi-commit authoring, not a code/build break — `Go build & test` was not among the failing jobs on those runs). This is informational, not a finding: HEAD is fully green.
- **No wake/scan self-invocation:** I ran `grep -rn "issues-dispatch\|issuesdispatch\." .github/workflows/ src/packages/cnos.cds/orchestrators/` myself — 0 hits. The new package (`src/packages/cnos.issues/commands/issues-dispatch/`) has exactly one production caller path (`cmd_issues_dispatch.go` → `cn` CLI noun-verb routing), confirmed by reading the file directly; no dispatch-wake or reconciler code path reaches it.
- **"Not automated" sentence, true location:** I ran `grep -rn "not automated" src/packages/cnos.core/skills/agent/` myself. It lives at `dispatch-protocol/SKILL.md:132` — `"The status:ready → status:todo transition is the operator's authorization event. It is not automated."` — unchanged in substance (confirmed by reading the line directly; the new §1.2 paragraph is additive, appended after this sentence, and does not alter its wording). α's correction of both γ's scaffold citation and its own β-prompt's citation (both said `label-doctrine/SKILL.md` §1.2, which is actually "Dispatchability label" with no equivalent sentence) is correct — I independently re-derived this rather than trusting α's correction at face value.
- **`status:ready` remains a valid, honored hold state:** confirmed by reading `dispatch.go:130` (`needsLabelFlip := status == "ready"`) — the primitive only ever flips a label when the issue is currently at `status:ready`; any other status (including `status:todo`, confirmed by `TestDispatch_AlreadyCleanAndTodo_SafeNoOp`) is left untouched.
- **Build/test green, re-run independently:** `go build ./src/go/... ./src/packages/cnos.cdd/commands/cdd-verify/... ./src/packages/cnos.issues/commands/issues-map/... ./src/packages/cnos.issues/commands/issues-fsm/... ./src/packages/cnos.issues/commands/issues-dispatch/...` → clean, exit 0. `go test` on the same set → all `ok`, including pre-existing `internal/cli`, `issues-fsm`, `issues-map` suites (no regression) and the new `issues-dispatch` package (14/14). Also ran `go build ./...` / `go test ./...` from `src/go/` directly → clean.

### Standard β obligations

**Rule 7 — implementation-contract coherence.** Verified each of the 7 pinned axes against the diff directly:

| Axis | Pinned | Diff conformance |
|---|---|---|
| Language | Go + Markdown only | Confirmed — no other language in the 20-file diff (`.github/ISSUE_TEMPLATE/cdd-issue.md` is Markdown; `go.work`/`go.mod` are Go-toolchain config, not a new language). |
| CLI integration target | `cn` subcommand, sibling of `issues fsm`/`issues map` | Confirmed — `main.go` diff adds `reg.Register(&cli.IssuesDispatchCmd{})` beside the two existing siblings; `IssuesDispatchCmd.Spec()` returns `Name: "issues-dispatch"`, matching the noun-verb registry pattern I read in `dispatch.go` of `internal/cli`. Not a standalone binary. |
| Package scoping | `src/go/internal/cli/cmd_issues_dispatch.go` (thin wiring) + `src/packages/cnos.issues/commands/issues-dispatch/` (domain logic), registered in `main.go` | Confirmed exactly — `cmd_issues_dispatch.go` is 42 lines, delegates immediately to `issuesdispatch.Run`; all domain logic lives in the pinned package path with its own `go.mod`. |
| Existing-binary disposition | Additive, nothing replaced | Confirmed — no existing file's behavior is deleted or altered outside the two doctrine-file additions and the template comment. |
| Runtime dependencies | GitHub REST API only, `ghRequest`/`ghAddLabel`/`ghRemoveLabel` reused-in-shape, `ghEditIssueBody` as the one new call | Confirmed by reading `fetch.go` directly — `ghEditIssueBody` (PATCH) is modeled on `ghAddLabel`'s exact shape (marshal payload → `ghRequest` → status-code check). No other new runtime dependency. |
| JSON/wire contract preservation | N/A/additive | Confirmed — new command, no existing wire contract touched. |
| Backward-compat invariant | label semantics + `transitions.json` + claim/dispatch mechanics unchanged | Confirmed via the hard `transitions.json` byte-identical check above and the `status:ready`-gated mutation logic in `dispatch.go`. |

All 7 axes conform. No D-severity `implementation-contract` finding.

**Non-goal / out-of-scope check.** `git diff --stat origin/main...cycle/640` (re-run independently) touches exactly 20 files; none under `.github/workflows/` or `.cn-sigma/` (confirmed empty diff against those paths), none referencing #618/#626/#630 recovery-behavior code, no new status label introduced (`label-doctrine/SKILL.md` diff is a 1-line cross-reference addition, not a new label declaration), and `transitions.json` is untouched. Scope guardrails hold.

**gamma-scaffold.md presence gate.** `.cdd/unreleased/640/gamma-scaffold.md` is present on `origin/cycle/640`, predates α's dispatch (commit `600fe8e2`, before α's `766df7e8`/`cb861886`). Already satisfied per γ/δ's own claim sequence.

### Points of disagreement / independent notes (not findings — informational)

- α's self-coherence.md states branch CI is "unknown, not asserted green" — I independently confirmed it **is** observable and **is** green on HEAD, including all 10 gate jobs. This is a favorable correction (α under-claimed rather than over-claimed), not a defect.
- The git identity on this branch's commits (`sigma@cnos.cn-sigma.cnos`, environment-level) does not match the role-specific `alpha@.../beta@cdd.cnos` canonical forms per `alpha/SKILL.md` §2.6 row 14 / `beta/SKILL.md` pre-merge gate row 1. α disclosed this explicitly as known debt (path (b), environment-level legacy) rather than silently deviating. I note the same constraint applies to my own review commit in this environment; consistent with the branch's existing pattern and α's disclosed reasoning, not re-litigated here.

## Verdict: **converge**

All four AC oracles pass on independent re-derivation (not on α's self-coherence.md framing). Every fixture claim was re-run and observed to pass (14/14). The two hard/mechanical checks (`transitions.json` byte-identical; no wake self-invocation) both hold on direct grep, not restated from α's claims. The implementation contract's 7 axes all conform. Scope guardrails hold — no diff hunk touches #618, #626, #630's recovery code, adds a status label, or changes `transitions.json`. CI is green on HEAD across all 10 gates. No unresolved findings. Proceeding to the merge step is appropriate on this basis (merge/close-out execution outside this review round's scope per this dispatch's instructions).
