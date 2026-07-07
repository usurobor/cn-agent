# α close-out — cnos#615

## Scope note

This close-out is written at the **δ §9.5 converge boundary** — β returned
`verdict: converge` at R0 (`beta-review.md`, head `a4e7128`), and this file
is one of the three artifacts that boundary requires
(`alpha-closeout.md` + `beta-closeout.md` + `gamma-closeout.md`), per
`delta/SKILL.md` §9.5's per-R[N] artifact contract for the `converge` row.
It is **not** the post-merge/post-release close-out `alpha/SKILL.md` §2.8
otherwise describes (that flow assumes β has already merged to `main`;
here, merge/PR/release are still ahead of us — δ's next step, per §9.6, is
to open the cycle-PR and request the `status:review` transition). The
retrospective content below is written against the full R0 cycle record
(`gamma-scaffold.md`, `self-coherence.md`, `beta-review.md`) as it stands
at this boundary, all on `cycle/615`, none of it yet on `main`.

## Summary — what was built

A terminal-hygiene reconciliation pass for the `cn issues fsm` command
family: new sub-verb `cn issues fsm terminal --protocol P [--apply]`
(new file `terminal.go`, +1 new primitive `ghEnsureLabelExists` in
`fetch.go`, wiring in `issuesfsm.go`). For a **closed** issue that still
carries `dispatch:cell` + `protocol:{P}` (and, incidentally, any stale
`status:*` label), the pass strips those three label classes and adds the
`resolution/*` label keyed off GitHub's own `state_reason`
(`completed` → `resolution/completed`; `not_planned` →
`resolution/not-planned`, a label this cell creates on first live use).
It is a sibling reconciler to the #593 recovery scanner (`scan.go`), not a
modification of it and not a new FSM terminal state — `transitions.json`
is untouched, `Evaluate`/`Table` are never called from the new code path,
and idempotence falls out structurally from the list query's own
`state=closed&labels=...` filter (no bolted-on "already processed"
marker was needed).

The one design call the scaffold left open — new `terminal` sub-verb vs.
a `--closed` flag folded into `scan` — was resolved in favor of the new
sub-verb, documented with rationale in `self-coherence.md` §Self-check,
and independently confirmed correct by β (`beta-review.md` §Guardrail
verification).

## AC risk areas

The two ACs carrying the most correctness risk going in were:

- **AC2's `not_planned` → `resolution/not-planned` mapping.** This is the
  exact ambiguity the operator's clarifying comment (2026-07-07T01:51:47Z)
  resolved, superseding the issue body's original "-or-equivalent"
  phrasing. Getting this backwards (mapping to `resolution/completed` or
  a `resolution/wontfix` label that does not exist) would have been a
  silent correctness bug, not a style nit — it was tested with an
  explicit assertion rejecting both wrong alternatives
  (`TestTerminal_ClosedNotPlannedStale_ResolvesToNotPlanned`), and β
  independently re-derived the same mapping from the pinned comment
  rather than trusting the test's own framing.
- **AC4's transitions.json/Evaluate isolation.** Because this pass lives
  in the same package as `scan.go` and reuses its label-mutation
  primitives, the risk was that the new code path would (even
  accidentally) route a closed issue through `Evaluate()` against
  `transitions.json`, or that a duplicate DELETE-labels HTTP call would
  be written outside `ghRemoveLabel`. Both risks were closed by structural
  proof, not just a passing test: `TerminalOptions` carries no `Table`
  field at all (pinned by
  `TestTerminalOptions_HasNoFSMTableDependency`), and the only
  `http.MethodDelete` call anywhere in the new code is the pre-existing
  one inside `ghRemoveLabel`.

Both risk areas were independently re-walked by β at R0 and confirmed
correct, not merely re-read from α's narrative.

## Debt carried forward — status at this boundary

`self-coherence.md` §Debt named three items. Status now:

1. **"branch CI is not independently confirmed green on GitHub Actions for
   this head commit"** — **closed.** β independently ran `gh run view`
   against the live Actions run for the actual head commit at review time
   (`1c18c68`, prior to β's own commit) and confirmed every named job
   green (`Go build & test`, `Binary verification`, `Package
   verification`, I1/I2/I4/I5/I6, both dispatch-integrity guards),
   with `install-wake-golden` correctly not triggering at all (path
   filter). This closes the one debt item α could not self-certify from
   a sandbox with no Actions runner — β's live-CI confirmation is strictly
   stronger evidence than α's local `go build`/`go vet`/`go test`
   equivalents, and it is now the record of truth for AC6.
2. **"no dedicated JSON fixture files under `testdata/`"** — **not a
   gap; confirmed a legitimate, disclosed divergence.** β read every test
   body in `terminal_test.go` in full and confirmed all four scaffold-
   required scenarios are genuinely exercised via the injected-fake
   (`TerminalOptions`) seam. Remains noted for anyone auditing this
   package by directory-listing convention alone.
3. **"no live mutation was performed against the real repo"** — **not
   resolved and not expected to be at this boundary**, by design: neither
   α nor β has label-mutation authority (out of role scope for both), so
   `ghEnsureLabelExists`'s live creation of `resolution/not-planned`
   remains untested against the real GitHub API until the first operator
   `--apply` run against an actual `not_planned`-closed issue. This is
   expected production behavior, not cycle debt — the creation path's two
   branches (create-on-first-use, tolerate-422-already-exists) are unit-
   tested against a fake server.

No new debt surfaced during β's review that wasn't already named in
`self-coherence.md` — β's `## Findings` section reports none.

## Calibrated success claim

What is verified, as of this boundary:

- All six scaffold ACs (AC1–AC6) pass, independently re-derived by β from
  the code/tests/diff, not merely asserted by α.
- The full pre-existing `issues-fsm` test suite (evaluate/scan) passes
  unchanged alongside the 13 new test functions — zero behavior change to
  the active-lifecycle reconciler.
- `transitions.json` and `src/packages/cnos.core/labels.json` are both
  byte-identical to `origin/main`.
- Live GitHub Actions CI is green on the actual cycle-branch head commit
  at β's review time (`1c18c68`), independently confirmed via `gh run
  view`, across every job the scaffold's AC6 oracle names.

What is **not yet verified**, because it is not this boundary's job:

- No PR has been opened and no merge to `main` has happened — that is
  δ's next step (§9.6 `status:review` transition), not something either
  α or β has done or claims to have done.
- CI has not been (and cannot yet be) observed on a post-merge `main`
  commit, since no merge exists yet.
- The `resolution/not-planned` label does not yet exist on the live repo;
  its creation is deferred to the first real `--apply` invocation.

**Claim:** the implementation is correct and complete against the
scaffold's six ACs and all stated guardrails, as verified by two
independent passes (α's own self-coherence walk and β's from-scratch
re-derivation including live CI). It is ready for δ to proceed to the
`status:review` transition. It is not yet released, merged, or tagged,
and no claim is made about post-merge CI or live-repo label state beyond
what is stated above.
