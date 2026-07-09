# α close-out — cnos#640

[provisional — pending β outcome at time of writing; converged at R0, no fix-round required]

## Summary

Built the mechanization the issue named: (1) a `cn issues dispatch --issue N` primitive
(`src/packages/cnos.issues/commands/issues-dispatch/`, wired via
`src/go/internal/cli/cmd_issues_dispatch.go` as a sibling of `cn issues fsm` / `cn issues map`)
that performs the body-hold-phrase cleanup and the `status:ready → status:todo` label flip as
one caller-visible operation; and (2) doctrine changes naming labels as the sole source of
truth for dispatch readiness (`dispatch-protocol/SKILL.md` §1.2 new paragraph + D13
failure-mode catalogue entry) plus a documented observe → capture → detect-recurrence →
mechanize → verify-non-recurrence loop (`CELL-KINDS.md`, new section). No standalone
contradiction gate (candidate 3) was built — the detection logic lives inside `Dispatch`
itself, and a separate scanning gate was judged out of scope for this cell's sizing.

Diff: 21 files, 1634 insertions / 2 deletions. Go + Markdown only, matching the pinned
implementation contract's language axis.

## AC-by-AC

- **AC1** (body/label contradiction structurally corrected) — closed by the primitive itself.
  14 tests in the new package prove: the exact #614/#633 shape corrected atomically
  (`TestDispatch_ContradictoryShape_BodyCorrectedLabelUntouched`), a clean issue is a safe
  no-op (`TestDispatch_AlreadyCleanAndTodo_SafeNoOp`), the operation is idempotent
  (`TestDispatch_Idempotent_SecondRunIsNoOp`), and a body-edit failure can never produce a
  half-applied contradiction (`TestDispatch_BodyEditFails_LabelNeverFlipped`, body-edit-first
  ordering by construction). Verified live end-to-end wiring via
  `go run ./src/go/cmd/cn issues dispatch --issue 999999 --repo doesnotexist/doesnotexist`
  (exit 1, fetch error surfaced — full call chain executes).
- **AC2** (wake no longer decides "body says held, label says go") — closed via
  `dispatch-protocol/SKILL.md` §1.2's new paragraph stating labels alone gate dispatch
  readiness; `label-doctrine/SKILL.md` §7 gets a one-line cross-reference rather than a
  restatement (one canonical home for the doctrine text); `delta/SKILL.md` §9.2 input #1
  checked and found already consistent (negative result, recorded rather than skipped).
- **AC3** (loop documented) — closed via `CELL-KINDS.md`'s new "Process self-improvement
  loop (cnos#640)" section: all five steps named, today's owner named per step honestly
  (detect-recurrence has no current owner — quotes κ's comment verbatim rather than papering
  over it), cites #614→#633 as the worked example, links κ's meta-comment as the seed,
  cross-linked from `CDD.md` and `CELL-KINDS.md`'s own cross-reference list.
- **AC4** (gates green, no weakening of hold authority) — `transitions.json` byte-identical to
  `origin/main` before and after rebase; the "not automated" operator-authorization sentence
  unchanged in substance (see corrections below for its true file location); zero wake/scan
  call sites reach the new primitive (`grep -rn "issues-dispatch\|issuesdispatch\."
  .github/workflows/ src/packages/cnos.cds/orchestrators/` → 0 hits); `status:ready` remains
  gated correctly (`needsLabelFlip := status == "ready"`); local build/test green across all
  touched packages, no sibling regression. Branch CI itself was not observable from α's
  sandboxed environment — recorded as unknown, not asserted, in self-coherence.md.

## Debt disclosed

- **No standalone contradiction gate (AC1 candidate 3).** γ's scaffold offered this as
  optional. α judged it net-negative complexity for this cycle: a same-repo scanning gate over
  *all* open issues is a materially different shape (a selector loop, not single-issue
  orchestration) than what the pinned implementation contract scoped, and the
  operator/human-invoked-only framing (AC4) does not cleanly extend to a background scanner
  without a fresh design call. Left as a candidate for a future, separately-scoped cell.
- **"ε made concrete" left as a named commitment, not a filed issue.** The AC3 loop doctrine's
  detect-recurrence row names the gap and its shape but explicitly defers filing the issue that
  builds a periodic ε pass to "a follow-up action for γ/κ." α chose the commitment form over a
  fabricated placeholder issue number, per the α prompt's explicit "issue number, or an explicit
  'next: ε made concrete' commitment" latitude — filing the issue itself is a different
  `cell_kind` (issue-authoring) than this cell's scope.
- **Git identity does not match the canonical `alpha@cdd.cnos` form** (session identity is
  `sigma@cnos.cn-sigma.cnos`, matching γ's and δ's own commits on this branch). Disclosed as
  path (b) — environment-level legacy, accepted rather than retroactively rebased — because a
  retroactive rebase-and-reauthor would split the cycle's history into a mixed-identity form
  without achieving the underlying goal across the whole branch.
- **Branch CI (GitHub Actions) not directly observable from α's sandbox.** α ran the closest
  local equivalent (`cn cdd verify --unreleased`, real build, real tests) and recorded live CI
  as "unknown, not asserted" rather than claiming green — β was the first surface able to
  observe it directly (and did; see beta-review.md).
- **No automated end-to-end test against the real GitHub API.** All 14 tests run against
  injected fakes or a local `httptest.Server`, per the scaffold's explicit oracle requirement
  ("NOT a live-network manual check") — `ghEditIssueBody`'s status-code handling follows
  `ghAddLabel`'s pattern by construction, not by live observation.

## Two self-caught corrections to γ's scaffold

α re-verified two premises in `gamma-scaffold.md` against the live repository state before
acting on them, rather than trusting the scaffold's framing, and found both wrong:

1. **D8/D9 failure-mode slot.** The scaffold's friction note 2 flagged D8/D9 in
   `dispatch-protocol/SKILL.md` §5's failure-mode catalogue as "unused/unassigned" — a
   "natural free slot" for the new body/label-drift entry. α ran `grep -n "D8\|D9"` against the
   live file before writing anything and found both already defined and taken (D8 "Hard-coded
   runtime in launch step," D9 "Claim-comment failure not released"). The new entry landed as
   D13 (after the existing D12) instead. Named explicitly in self-coherence.md §ACs AC2 rather
   than silently substituted.
2. **AC4's "not automated" sentence, wrong file cited.** Both `gamma-scaffold.md` §1 AC4 and its
   own β prompt (§4) cited `label-doctrine/SKILL.md` §1.2 as the home of the operator-
   authorization "not automated" sentence. α verified via
   `grep -rn "not automated" src/packages/cnos.core/skills/agent/` that the sentence actually
   lives in `dispatch-protocol/SKILL.md` §1.2 — `label-doctrine/SKILL.md`'s own §1.2 is a
   different section ("Dispatchability label") with no equivalent sentence. Named explicitly so
   β could re-verify the correction rather than re-trust the scaffold's citation (β did
   independently re-derive this and confirmed it correct).

Both corrections were caught by re-running the falsification-gate discipline (Kernel §1.1 —
"run one command designed to disprove it") against the scaffold's claims before acting on them,
not by assuming the scaffold was authoritative because γ wrote it. Neither correction changed
scope or ACs; both are citation/slot-selection fixes.

## Cell-kind

Confirmed `doctrine` (γ's tentative pin), not reclassified to `implementation` despite the diff
being Go-dominant by line count — the dominant matter produced is the doctrine fix (labels as
sole source of truth; the loop documented honestly); the Go primitive serves that doctrine
rather than the reverse. Reasoning recorded in self-coherence.md §Review-readiness as a
low-stakes, explicitly-authorized judgment call per γ's own scaffold framing.

## Pattern note

Same class as the D8/D9 and AC4-citation corrections above: a scaffold authored before the
implementer reads the live file can carry premises that were true at scaffold-authoring time (or
never verified against the live file at all) and drift stale, or were simply wrong, by
implementation time. Two occurrences this cycle, both self-caught, both cheap to fix once
caught. No recommendation on disposition — that is γ's triage call.
