# α close-out — cnos#493

## Scope note

This close-out is written at the **δ §9.5 converge boundary** — β
returned `verdict: converge` at R0 (`beta-review.md`, head `1f76b126`),
and this file is one of the three artifacts that boundary requires
(`alpha-closeout.md` + `beta-closeout.md` + `gamma-closeout.md`), per
`delta/SKILL.md` §9.5's per-R[N] artifact contract for the `converge`
row. It is **not** the post-merge/post-release close-out
`alpha/SKILL.md` §2.8 otherwise describes (that flow assumes β has
already merged to `main`; here, merge/PR/release are still ahead — δ's
next step, per §9.6, is to open (or confirm) the cycle-PR and request
the `status:review` transition). The retrospective content below is
written against the full cycle record (`gamma-scaffold.md`,
`self-coherence.md`, `beta-review.md`) as it stands at this boundary,
all on `cycle/493`, none of it yet on `main`. Voice per `alpha/SKILL.md`
§2.8: factual observations and patterns only — dispositions are γ's
job, not named here.

## Summary — what was built

A `cn label-doctor` audit + repair engine for cnos.core's canonical
GitHub label set (`src/packages/cnos.core/labels.json`, schema
`cn.labels.v1`, 8 entries), plus the CLI/CI wiring that makes it
reachable and enforced:

- New Go module `src/packages/cnos.core/commands/label-doctor/`
  (`doctor.go`, `github.go`, `manifest.go`, `resolve.go`, `cli.go` +
  tests), dependency-free GitHub REST primitives (`ghListLabels`,
  `ghCreateLabel`, `ghUpdateLabel` — no `gh` CLI shellout), mirroring
  `cnos.issues/commands/issues-fsm/fetch.go`'s idiom. Added as the 5th
  entry in root `go.work`.
- A thin kernel subcommand (`src/go/internal/cli/cmd_label_doctor.go`,
  registered in `src/go/cmd/cn/main.go`), invoked as `cn label doctor`
  (noun-verb form — `dispatch.go`'s pre-existing routing for any
  hyphenated `Spec().Name`, the same shape `issues-fsm` → `issues fsm`
  already uses; not an α improvisation).
- `src/go/internal/repoinstall/repoinstall.go`'s
  `ensureCanonicalDispatchLabels()` stub (which unconditionally
  returned a `"cnos#493 label-install mechanism is not yet available"`
  error) replaced with a real in-process call into the new package —
  not left as a parallel dead path.
- A CI guard: a new `build.yml` step running the new module's test
  suite (mirroring the existing `issues-fsm` module-test-step
  pattern), with fixture-backed tests exercising both drift-fails and
  clean-passes directions through the actual `Run()` CLI entry point,
  plus the new command name added to the CLI-ergonomics `--help`
  smoke-test loop.
- A committed audit artifact (`.cdd/unreleased/493/label-audit.md`)
  recording the live before/after label state.

Live repo state at the time of β's final review: all 8 canonical
labels present, color+description byte-equal to `labels.json`
(independently re-verified by β via `gh label list`).

## The mid-cycle resumption

This cycle's commit history (`git log --oneline main..HEAD`, 20
commits) spans two distinct α sessions, not one continuous session:

1. **Original implementation session** — produced 12 commits: the γ
   scaffold, the full `label-doctor` module + CLI wiring + CI guard
   (AC1–AC5 implementation), the live audit + repair record, and the
   first four `self-coherence.md` sections (`## Gap`, `## Skills`,
   `## ACs`, and a `labels.json` fix commit interleaved between them).
   This session's run died before writing `## Self-check`, `## Debt`,
   or `## CDD Trace` — the three sections `alpha/SKILL.md` §2.5 lists
   after `## ACs`.
2. **Recovery scanner reconciliation** — a later scheduled scan
   (`cn issues fsm scan --apply`) observed the dead run against
   checkpointed matter (the `cycle/493` branch + draft PR #634, no
   `REVIEW-REQUEST.yml`) and mechanically reconciled
   `status:in-progress → status:todo`, preserving the branch and PR,
   per `transitions.json`'s `propose_status_todo_with_matter` rule.
   `CLAIM-REQUEST.yml` records this explicitly: "a prior firing's claim
   on this same issue reached this point... then died before
   REVIEW-REQUEST.yml was written; the recovery scanner's
   propose_status_todo_with_matter rule reconciled
   status:in-progress -> status:todo, preserving the branch/PR."
3. **This resumption** — re-claimed the cell per `delta/SKILL.md`
   §9.11 "resumed-from-mechanical-reversion" and, per this α's
   dispatch instruction, did not restart the cycle. It read the
   existing 12 implementation commits and the partial
   `self-coherence.md` (through `## ACs`), independently re-ran every
   AC oracle against real command output rather than trusting the
   existing prose, found one stale claim (see below), corrected it in
   place with an explicit correction note, and completed the three
   missing sections (`## Self-check`, `## Debt`, `## CDD Trace`) plus
   the review-readiness signal.

No implementation commit was discarded or redone. The resumption's own
work is additive: one correction to `## AC2`'s text (marked as a
correction, not a silent rewrite), the three missing sections, and a
rebase onto `origin/main` (which had advanced 22 commits since the γ
scaffold's base SHA across both sessions combined).

## Friction observed

1. **A real `labels.json` data defect surfaced as a blocking
   dependency of AC2, not a nice-to-have.** `labels.json`'s
   `dispatch:cell.description` was 149 bytes. GitHub's label API
   rejects any description over 100 characters (`HTTP 422:
   "description is too long (maximum is 100 characters)"`). This meant
   AC2 ("all canonical labels exist with canonical color+description")
   could not reach 8/8 until the manifest's own data was edited —
   `dispatch:cell.description` was shortened to 92 bytes (commit
   `0135b30b`, landed on this branch during the original implementation
   session, ahead of the AC1 audit-artifact commit). This was not
   optional scope-creep-avoidance; without it, AC2 was structurally
   unreachable for one of the 8 labels regardless of how correct the
   repair engine's code was.
2. **The manifest fix and the live repair were two separate actions,
   and they landed in two different sessions.** The `labels.json` edit
   (commit `0135b30b`) shortened the description in the source-of-truth
   file, but the live `usurobor/cnos` label on GitHub still carried the
   old (empty, GitHub-defaulted) description at the start of this
   resumption — the manifest fix had not yet been re-applied against
   the live repo. This resumption ran `cn label doctor --repo
   usurobor/cnos` (apply, not dry-run) and confirmed via a follow-up
   `--dry-run` that all 8/8 labels, including `dispatch:cell`, then
   reported `match`. The two-step nature of "fix the manifest" vs.
   "apply the manifest to the live repo" is a distinction this cycle's
   own tool makes visible (dry-run vs. apply) but is easy to lose track
   of across a session boundary.
3. **`Spec().Name`'s hyphen routes through the noun-verb dispatch path,
   not a flat token, and this was only confirmed empirically, not
   assumed from the contract's naming precedent.** The implementation
   contract's CLI-integration axis described the target as "single-token,
   matching `cn doctor`/`cn build`/`cn status` naming" while also
   instructing "mirror `cmd_issues_fsm.go` exactly." Those two
   instructions conflict for any `Spec().Name` containing a hyphen: the
   three named exemplars (`doctor`/`build`/`status`) are genuinely flat,
   non-hyphenated names, while `issues-fsm` (the structural precedent to
   mirror) is itself only reachable via its noun-verb space form
   (`cn issues fsm`), per `dispatch.go`'s own documented invariant that
   hyphenated flat forms are not supported. `cn label-doctor --help`
   resolves to the `label` group listing, not the command; `cn label
   doctor --help` is the real invocation. `cell.go`'s existing runtime
   error text ("Run label-doctor before retrying") uses the hyphenated
   form as operator-facing prose/mnemonic rather than a literal
   invocable token, which is a small, pre-existing prose/documentation
   gap this cycle did not correct (see `self-coherence.md` §Debt item 5
   — `cell.go` is out of scope per the γ scaffold's guardrails).

## Debt not otherwise covered in self-coherence.md §Debt

`self-coherence.md` §Debt already lists 6 items (the `dispatch:cell`
gap resolved-not-carried-forward; the git-remote-parsing utility's
narrow `github.com`-only / `origin`-only scope; `githubAPIBase`'s
hardcoded default; the vendored-install manifest-resolution path's
fixture-only verification; `cell.go`'s stale hyphenated-form prose; and
the bootstrap git-identity exemption). No additional debt surfaced
during this close-out pass beyond what that section already names —
the resumption's own re-verification pass (re-running all test suites,
`gofmt`, `go vet`, the live label audit, and the CLI smoke checks from
a clean shell) did not surface anything `self-coherence.md` §Debt does
not already disclose.

## Calibrated success claim

What is verified, as of this boundary:

- All 5 scaffold ACs (AC1–AC5) pass, independently re-derived by β from
  the code/tests/live GitHub state — every test cited in
  `beta-review.md` was re-run by β from a clean shell and the rebuilt
  `cn` binary was exercised directly, not read off this α's pasted
  output.
- 32/32 tests pass in the new `label-doctor` module; 15/15 package
  suites pass in `src/go` (no regression).
- Live GitHub Actions CI is green on the cycle-branch head commit at
  β's review time (`1f76b126`'s parent, `dd9693ca`), independently
  re-pulled by β via `gh api`.
- Live `usurobor/cnos` label state is 8/8 canonical, independently
  re-verified by β against `labels.json`.

What is **not yet verified**, because it is not this boundary's job:

- No merge to `main` has happened yet — δ's next step per §9.6 is the
  `status:review` transition request, not something this α does.
- CI has not been (and cannot yet be) observed on a post-merge `main`
  commit, since no merge exists yet.
- No claim is made here about the parent issue's close state or any
  post-release assessment — those are γ's later, post-merge actions.
