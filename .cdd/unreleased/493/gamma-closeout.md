# γ close-out — cnos#493

## Scope note

This file is written at the **δ §9.5 converge boundary**: β returned
`verdict: converge` at R0 (`beta-review.md`, head `1f76b126`), and δ's
§9.5 per-R[N] artifact contract's `converge` row requires all three
closeout artifacts (`alpha-closeout.md`, `beta-closeout.md`,
`gamma-closeout.md`) to land on `cycle/493` before δ writes the
`status:review` return token. This is **not** the full γ §2.7/§2.10
closure declaration — that closure gate (post-release assessment,
cycle-directory move to `.cdd/releases/{X.Y.Z}/{N}/`, `RELEASE.md`,
hub-memory update, `gh issue view 493` close-state assertion,
merged-branch cleanup, δ release-boundary preflight) presumes a merge
and a release boundary that have **not happened yet**: PR #634
(`cycle/493` → `main`) is open and draft (`gh pr view 634` →
`state: OPEN`, `isDraft: true`), `gh issue view 493 --json state` →
`OPEN`, and `main` is untouched by this cycle's work. δ's next step
per §9.6 is to open (or confirm) the cycle-PR and request the
`in-progress → review` transition; the operator merges later. What
follows is the **process-gap audit** portion of γ's role at this
boundary — did the cycle reveal friction worth converting into a
follow-up, or not — written as a pre-merge, cycle-level retrospective,
not a release-time artifact.

## Cycle summary

- **Issue:** cnos#493 — "label-doctor: install + verify cnos.core
  canonical labels in repos" (live-discovery follow-up to cnos#491).
- **Shape:** single R0 pass, converged with no `iterate` round and no
  RC — but that R0 pass spans **two separate α sessions** across a
  mechanical requeue (detailed below), not one continuous session.
- **Result:** a new Go module (`src/packages/cnos.core/commands/
  label-doctor/`) implementing a dependency-free GitHub-REST label
  audit+repair engine against `labels.json`; a `cn label doctor`
  kernel subcommand; `repoinstall.go`'s `ensureCanonicalDispatchLabels()`
  stub replaced with a real in-process call; a CI guard exercising
  both the drift-fails and clean-passes directions through the real
  `Run()` entry point; a committed live audit artifact. All 5 ACs pass,
  independently re-derived by β from live `gh label list` output, a
  rebuilt binary, and re-run test suites (32/32 new-module tests,
  15/15 `src/go` package suites, no regression). Live CI green on the
  review head (`gh api .../check-runs`, independently pulled by β).

## The cnos#630 resumed-from-mechanical-reversion shape: first production witness

This is the load-bearing process observation of this cycle, and it is
named explicitly here rather than treated as background noise.

**What happened, per the branch's own record.** `CLAIM-REQUEST.yml`
states it directly: a prior firing's claim on cnos#493 produced
`cycle/493` + draft PR #634 + **12 implementation commits**, then died
before `REVIEW-REQUEST.yml` was written. `transitions.json`'s
`propose_status_todo_with_matter` rule (the mechanism cnos#630 itself
introduced) mechanically reconciled `status:in-progress → status:todo`
on the next scan, preserving the branch and the draft PR rather than
discarding them. A later dispatch-wake firing (`wake: cds-dispatch`,
`protocol: cds`, per `CLAIM-REQUEST.yml`'s `requested_transition: "todo
-> in-progress"`) then re-claimed the cell explicitly citing "cnos#630
resumed-from-mechanical-reversion shape (delta/SKILL.md §9.11)" as the
reason, and resumed rather than restarted.

**Why this is a first, not a repeat.** cnos#630 itself — the cycle
that introduced the `propose_status_todo_with_matter` rule and
`delta/SKILL.md` §9.11's doctrine text — disclosed in its own
`gamma-closeout.md` §"Disclosed debt" that "§9.11's doctrine addition
has no live wake-invoked-δ empirical witness yet," explicitly noting
that cnos#630's own cycle ran under bootstrap-δ (single-session
δ-as-γ), not a genuine wake firing, and that the missing witness was
"a legitimate disclosed limitation... mirrors §9.10's own
unresolved-witness history." **cnos#493 is that witness.** This cycle
ran under wake-invoked-δ (a real `cds-dispatch` claim/re-claim cycle,
per `CLAIM-REQUEST.yml`'s `wake:`/`protocol:` fields — not a bootstrap
exception), and it is the first cycle in which a real dispatch-wake
firing claimed a `status:todo` cell that the recovery scanner had
mechanically reconciled from a dead run, per the exact shape §9.11
names.

**Did the routing work as designed?** Yes, on every point §9.11's
routing sequence specifies:

1. **δ read the existing branch state first, and did not re-dispatch
   γ for a fresh scaffold.** `gamma-scaffold.md` already existed
   (committed `a703808a` during the original session); this pass's δ
   read it and dispatched α directly against the existing scaffold —
   no second scaffold was authored, no scope was re-litigated.
2. **δ dispatched the next incomplete role, per the artifacts already
   present, not a blanket restart.** `self-coherence.md` already
   carried `## Gap` / `## Skills` / `## ACs` with no review-readiness
   signal and no `beta-review.md` — exactly the "partially-written
   self-coherence.md" case §9.11 step 2 names, which routes to α (not
   β) under α's own `alpha/SKILL.md` §4 "Resumption" protocol. δ
   dispatched α, and α's `self-coherence.md` §Self-check names the
   resumption explicitly and independently re-verified every AC rather
   than trusting the prior prose.
3. **No implementation commits were discarded.** All 12 of the
   original session's implementation commits (the label-doctor module,
   CLI wiring, `repoinstall.go` stub replacement, CI guard, live audit
   artifact) remain on the branch, untouched by this resumption. The
   resumption's own commits are strictly additive: one correction to
   AC2's stale claim, the three missing `self-coherence.md` sections,
   a rebase onto `origin/main`, and the review-readiness signal.
4. **From that point, §9.3 applied unmodified** — β reviewed once and
   converged; no re-dispatch of α, no iterate round.

**Disposition:** this closes the loop cnos#630's own closeout left
open. No new tracking issue is needed for "§9.11 needs a production
witness" — it now has one, and the routing performed exactly as the
doctrine specifies. This paragraph, plus the `CLAIM-REQUEST.yml` +
`self-coherence.md` §Self-check / §CDD Trace records this cycle
already carries, **is** the evidence; nothing further needs filing.

## Root-cause triage — the original session's mid-cycle death

**What died, precisely.** The original session wrote `## Gap`, `##
Skills`, and `## ACs` (three of `alpha/SKILL.md` §2.5's six mandated
`self-coherence.md` sections, committed incrementally as that section
requires) plus 12 implementation commits, and then stopped before
writing `## Self-check`, `## Debt`, or `## CDD Trace` — the three
sections that come after `## ACs` in §2.5's mandated order.
`self-coherence.md`'s own `## Self-check` section (written by this
resumption) names the observable consequence directly: the truncation
"left CI's 'CDD artifact ledger validation (I6)' check red
(`self-coherence.md sections — missing required sections: CDD Trace
Self-check or Debt`)" — i.e., the failure is independently confirmed
by CI's own ledger check, not merely asserted by the resuming session.
Root cause is not directly observable from the branch (no crash log or
timeout trace is part of the CDD artifact set), but the shape — steady
incremental commits, then a hard stop mid-artifact with no partial,
uncommitted section fragment left behind — is consistent with a
session or stream timeout hitting during α's incremental writing
between committed sections, not a logic error or a deliberate exit.

**Disposition: no patch needed — this is the incremental-commit
discipline working exactly as designed, not a gap it exposed.**
`alpha/SKILL.md` §2.5 mandates one-section-per-commit specifically
because "stream timeouts will discard partial work" if the whole file
is written in one generation. That is precisely the failure this
cycle hit, and precisely why nothing was lost: the three sections
already committed (`## Gap`, `## Skills`, `## ACs`) and all 12
implementation commits survived the death intact, on the branch, where
the recovery scanner could find them and the resumption could read and
build on them rather than re-deriving from scratch. A one-big-commit
discipline would have lost the entire session's implementation and
partial self-coherence work to the same timeout; the incremental
discipline is the reason this cycle's cost was "resume and write three
sections" rather than "restart from zero." No skill patch is filed
here — the mechanism that made this recoverable is the existing rule,
verified working, not a newly discovered gap.

## Root-cause triage — the `labels.json` `dispatch:cell` description defect

**What happened.** `labels.json`'s `dispatch:cell.description` was 149
bytes. GitHub's label API rejects any description over 100 characters
(`HTTP 422: "description is too long (maximum is 100 characters)"`).
This blocked AC2 from reaching 8/8 until the manifest's own data was
edited (commit `0135b30b`, shortening the description to 92 bytes) —
a genuine `labels.json` data defect this cycle had to fix as a
blocking dependency of an acceptance criterion, not a nice-to-have
cleanup. The fix landed on the branch before the mechanical requeue;
the live repair (re-applying the manifest fix against the actual
`usurobor/cnos` repo) had to be re-run separately during the
resumption, since the manifest edit and the live label update are two
distinct actions the tool itself makes visible (`--dry-run` vs.
apply) — `alpha-closeout.md` §"Friction observed" item 2 has the full
account.

**Disposition: Project MCI — candidate follow-up, not filed as a
separate tracking issue this pass.** Reasoning: this is a single
occurrence (one label, one manifest, discovered live rather than
recurring across cycles) — per this same closeout discipline's own
two-occurrence bar for calling something a recurring pattern (per
`gamma/SKILL.md` §2.9), one instance does not yet clear that bar, so
this is not escalated to an immediate skill/spec patch landed in this
pass. But it is not a pure one-off either: `labels.json` has no
schema-level or CI-level check that any entry's `description` stays
under GitHub's 100-character API limit, so the same defect class can
recur silently on any future manifest edit, surfacing again only when
someone next runs the repair tool live against a real repo and hits a
422 — exactly what happened here. The concrete next action, named
here rather than left implicit: a description-length assertion,
either as a unit test in `label-doctor`'s own manifest-parsing tests
(`manifest_test.go`) or as a lightweight `cn cdd verify`-style lint
step, that fails fast on any `labels.json` entry whose description
exceeds 100 bytes. Not landed in this pass (this closeout is a
retrospective document, not an implementation dispatch) — named as the
first AC a future MCA on this thread would carry.

## Cycle-iteration triggers (`gamma/SKILL.md` §2.8) — none fired

- **Review churn** (>2 rounds): did not fire — 1 round (R0 only).
- **Mechanical overload** (>20% mechanical findings, ≥10 total): did
  not fire — `beta-review.md` §Findings names 0 blocking findings and
  2 non-blocking sanity-checked observations.
- **Avoidable tooling/environment failure**: the mid-cycle death
  (above) is the candidate trigger to weigh explicitly. Assessed as
  **not an avoidable-failure trigger in the "guardrail was missing"
  sense** — the guardrail that would prevent data loss from exactly
  this failure class (incremental commit-per-section) was already in
  place and is why nothing was lost; there is no additional guardrail
  this closeout is recommending. The mechanical-reversion recovery
  path itself (cnos#630's own contribution) is the tooling that turned
  a dead run into a resumable one rather than a stuck cell — also
  already landed, also confirmed working (see witness section above).
- **Loaded-skill miss**: none identified. `alpha/SKILL.md` §2.5's
  incremental-commit discipline and §4's resumption protocol both
  fired correctly and were both exercised for real by this cycle, not
  just cited abstractly.

Per `gamma/SKILL.md` §2.9's independent process-gap check: the
recurring friction worth naming is the `labels.json` description-length
class named above (disposed as a candidate follow-up, not yet filed,
per the one-occurrence bar); no gate was found too weak or vague
beyond that; no role skill failed to prevent a predictable error (the
opposite — the incremental-commit skill rule is why this cycle's death
was recoverable); no coordination burden suggested a better mechanical
path beyond what §9.11 already provides and this cycle just confirmed
works.

## Triage table

| Finding | Source | Type | Disposition | Artifact / commit |
|---|---|---|---|---|
| First production (wake-invoked-δ, non-bootstrap) empirical witness of `delta/SKILL.md` §9.11 "resumed-from-mechanical-reversion" | `CLAIM-REQUEST.yml`; `self-coherence.md` §Self-check, §CDD Trace step 2; this file §"cnos#630 resumed-from-mechanical-reversion shape" | doctrine/process | confirmed working as designed — no patch needed; resolves cnos#630's own disclosed "no live witness yet" debt item; recorded here as the closing evidence | `CLAIM-REQUEST.yml`; `.cdd/unreleased/493/self-coherence.md` §Self-check / §CDD Trace; commits `51ab7c58` (re-claim), `d9bbf45d` (§Self-check) |
| Original session's mid-cycle death (self-coherence.md truncated after `## ACs`, before `## Self-check`/`## Debt`/`## CDD Trace`) | `self-coherence.md` §Self-check (CI's I6 ledger check went red); this file §"Root-cause triage — the original session's mid-cycle death" | process/tooling | one-off, no patch needed — the `alpha/SKILL.md` §2.5 incremental-commit-per-section discipline is why zero work was lost; confirms the existing rule works, does not expose a new gap | `.cdd/unreleased/493/self-coherence.md` §Self-check; commits `bddf631e`/`0d2da6b7`/`bc3bb8c8` (surviving pre-death sections), `d9bbf45d`/`35514b6e`/`d504ef3a` (resumption's completion) |
| `labels.json`'s `dispatch:cell` description exceeded GitHub's 100-char API cap (149 bytes), blocking AC2 until fixed | `self-coherence.md` §Debt item 1, §ACs AC2 correction; `alpha-closeout.md` §"Friction observed" item 1; `beta-review.md` §Findings item 2 | data-defect / test-coverage gap | Project MCI — candidate follow-up, not filed as a separate tracking issue this pass (single occurrence); concrete next action named: a description-length assertion in `label-doctor`'s manifest tests or a `cn cdd verify`-style lint step | commit `0135b30b` (the fix); `.cdd/unreleased/493/label-audit.md` (residual-gap record) |
| `resolve.go` shells out to `git` via `os/exec` — could be misread against the contract's "no shellout" headline sentence | `beta-review.md` §Findings item 1 | documentation / contract-reading | drop — the contract's own body text explicitly directs this construction; not a violation, no action needed | `beta-review.md` §Findings item 1; `self-coherence.md` §Debt items 2–3 |
| `cell.go`'s runtime error text ("Run label-doctor before retrying") uses the hyphenated flat form as prose, not the real `cn label doctor` invocation | `self-coherence.md` §Debt item 5; `alpha-closeout.md` §"Friction observed" item 3 | documentation | drop — cosmetic, self-correcting on first use (operator typing the hyphenated form is routed to the `label` group listing showing the real subcommand); `cell.go` is explicitly out of scope per the γ scaffold's guardrails | `self-coherence.md` §Debt item 5 |
| Git-remote-parsing utility (`resolve.go`) is `github.com`-only, `origin`-only; `githubAPIBase` hardcoded to `api.github.com` | `self-coherence.md` §Debt items 2–3 | scope (deliberate) | drop — explicitly in-scope-as-narrow per the γ scaffold's own framing ("a small, scoped utility, not a design blocker"); GHE support was never in scope | `self-coherence.md` §Debt items 2–3 |
| Vendored-install manifest-resolution path (`.cn/vendor/...labels.json`) fixture-tested only, not live-verified against a real `cn repo install`'d tenant repo | `self-coherence.md` §Debt item 4 | test-coverage | drop — `usurobor/cnos` was the only repo this cycle had live GitHub API access to; a real cross-repo live verification is a different cycle's scope, not a gap this cycle should have closed | `self-coherence.md` §Debt item 4 |
| Commit author identity (`sigma@cnos.cn-sigma.cnos`) rather than canonical `alpha@cdd.cnos` role-identity pattern | `self-coherence.md` §Debt item 6; `beta-review.md` §Diff Context | process (disclosed exemption) | drop — explicit, named bootstrap/single-operator exemption for this cycle's dispatching context, consistent across α and β; not a defect | `self-coherence.md` §Debt item 6 |

## What γ is not asserting at this boundary

- No claim about the parent issue's close state (`gh issue view 493
  --json state` → `OPEN`, confirmed above) — the issue is not yet
  closed, and closure only follows a merge that has not happened. The
  `gamma/SKILL.md` §2.10 row-15 issue-close assertion is a
  full-closure-gate item, not a converge-boundary item; it belongs to
  a later γ pass after the operator merges PR #634.
- No post-release assessment (PRA) is written here — `gamma/SKILL.md`
  §2.7 makes the PRA a post-merge artifact measuring released work;
  nothing has released yet. This is intentional scope-limiting, named
  explicitly, not an oversight.
- No hub-memory update, cycle-directory move (`.cdd/unreleased/493/` →
  `.cdd/releases/{X.Y.Z}/493/`), or `RELEASE.md` — all are
  full-closure-gate items (`gamma/SKILL.md` §2.10), out of scope for
  the §9.5 converge boundary. These happen after the operator reviews
  and merges PR #634, not now.
- No claim is made that PR #634 has been reviewed or merged — `gh pr
  view 634` shows `state: OPEN`, `isDraft: true` as of this writing.

## Next step

Per `delta/SKILL.md` §9.5, all three converge-boundary closeout
artifacts now exist on `cycle/493`. δ's next action is to ensure the
cycle-PR (#634) and any required `REVIEW-REQUEST.yml` are in place and
request the `status:review` transition via `cn issues fsm evaluate
--issue 493 --apply` per §9.6 — γ is not requesting that transition
itself and is not touching issue #493's labels or comments in this
pass. The external human operator reviews and merges PR #634 later;
the post-merge steps this closeout explicitly declines to take (issue
close assertion, PRA, cycle-directory move, hub-memory update,
branch cleanup) follow at that point, not now.

## Deliverable evidence (δ, cnos#524 closeout-integrity preflight)

Per `cds-dispatch/SKILL.md` §"Closeout integrity preflight" — recorded here, on the cycle's closeout, immediately before δ requests the `status:in-progress → status:review` transition:

```yaml
deliverable_evidence:
  pr: "#634 (cycle/493 -> main)"
  head_sha: "682f815dd372511ef82b89c3b98e9e89ab6ec610"
  base_sha: "90e9c8b2272052381848d1787441e7d12d994627"
  commits_beyond_base: 23
  closeout_artifacts: [gamma-scaffold.md, self-coherence.md, beta-review.md, alpha-closeout.md, beta-closeout.md, gamma-closeout.md]
```

All five checks named in the closeout-integrity preflight hold: (1) PR #634 exists and its body opens with `Closes #493`; (2) `cycle/493` HEAD (`682f815d`) differs from `base_sha` by 23 commits; (3) the `cycle/493` branch exists and differs from base (confirmed by the same commit count); (4) all six required `.cdd/unreleased/493/` closeout artifacts are present on the branch (listed above; independently confirmed via `git ls-tree -r origin/cycle/493 .cdd/unreleased/493/`); (5) this block itself names the PR number as evidence. PR #634 was marked ready-for-review (undrafted) and its title/body updated to reflect the converged cycle (previously an automated draft-checkpoint placeholder) as part of this same δ pass, per `delta/SKILL.md` §9.6's "δ opens (or updates) a cycle-PR scoped to the claimed cell."
