---
artifact: alpha-closeout
cycle: 556
issue: https://github.com/usurobor/cnos/issues/556
branch: cycle/556
author: alpha
---

# α close-out — cnos#556

## Summary

cnos#556 asked for a package home (`cnos.issues`) for the Issue Pivot board
domain, while keeping `cn issues map` as the public command. The cycle ran
two rounds:

- **R0** implemented the `cnos.issues` package (manifest, `SKILL.md`, three
  sub-skills) *and* physically relocated the Go domain implementation
  (`issuesmap.go`, `fetch.go`, `issuesmap_test.go`, `templates/`,
  `testdata/`) from `src/go/internal/issuesmap/` into
  `src/packages/cnos.issues/commands/issues-map/` as its own Go module,
  mirroring the cnos#392 `cdd-verify` precedent. β independently re-verified
  all 10 ACs and returned `verdict: converge`, flagging one non-blocking
  process note (Finding 1) about the relocation's justification.
- **δ overrode the R0 converge verdict.** The operator's clarifying comment
  on the issue contains an explicit, on-topic "Go implementation rule": *"Do
  not force Go implementation code into `src/packages/`. The active Go
  implementation may remain under `src/go/internal/issuesmap/` during the
  shim phase..."* — the opposite of what R0 implemented. δ's independent
  re-read of that comment found this and ordered an R1 repair.
- **R1** reverted the physical relocation via clean `git revert` of the two
  R0 commits that performed and referenced it, removed the now-unnecessary
  Go module wiring (`go.mod`, `go.work` entry, `src/go/go.mod`
  require+replace), reverted `cmd_issues_map.go`'s import path, and updated
  `cnos.issues/SKILL.md` + `skills/map/SKILL.md` to state — quoting the
  operator's rule verbatim — that the Go implementation lives at
  `src/go/internal/issuesmap/`, not under `src/packages/cnos.issues/`, and
  to narrate the R0→R1 history so a future reader isn't misled. β
  independently re-verified the repair and returned `verdict: converge`
  again.

## What was implemented (final state)

- `src/packages/cnos.issues/` exists: `cn.package.json` (schema
  `cn.package.v1`, `"version": "0.1.0"`, `"engines": {"cnos": ">=3.82.0"}`,
  no `commands` object), `SKILL.md`, `skills/map/SKILL.md`,
  `skills/taxonomy/SKILL.md`, `skills/triage/SKILL.md`.
- The Go domain implementation remains at `src/go/internal/issuesmap/`
  (final location — R0's relocation was reverted in R1), unchanged in
  behavior. `cmd_issues_map.go` is unchanged from its pre-cycle shape
  (already thin; import path never actually needed to change in the final
  state).
- `docs/development/board/README.md:42`'s path reference: unchanged from
  pre-cycle (points at `src/go/internal/issuesmap/`), since R1 reverted
  R0's edit.
- No change to `.github/workflows/board-map.yml`, no Node generator, no
  package-command-discovery redesign, no #216 work.

## Actual diff footprint (final, origin/main..HEAD)

Net of the R0-relocate/R1-revert round-trip, the merge-candidate diff is
much smaller than R0's own diff suggested mid-cycle:

- **New:** `src/packages/cnos.issues/cn.package.json`, `SKILL.md`,
  `skills/{map,taxonomy,triage}/SKILL.md` (doctrine/skill surfaces only —
  this is the actual net-new content this cycle contributes).
- **Unchanged, net of revert:** `src/go/internal/issuesmap/*`,
  `src/go/internal/cli/cmd_issues_map.go`, `src/go/cmd/cn/main.go`,
  `go.work`, `src/go/go.mod`, `docs/development/board/README.md`,
  `.github/workflows/board-map.yml`.
- **CDD artifacts:** `gamma-scaffold.md`, `self-coherence.md` (R0 + §R1),
  `beta-review.md` (R0 + §R1), plus this closeout round's three files.

The round-trip (relocate in R0, revert in R1) means the shipped code diff
is nearly a no-op against `origin/main` for the Go domain files — the only
durable code-level change from this cycle is the new `cnos.issues` package
metadata/doctrine tree. This is a smaller footprint than γ's R0 scaffold
anticipated, and smaller than what R0's own `self-coherence.md` reported at
review-readiness (which included the now-reverted Go-module wiring as a
substantial part of the "done" claim).

## What worked

- The cnos#392 (`cdd-verify`) precedent was a genuinely useful architectural
  reference for shaping the package/skill surfaces, independent of the
  location question that R1 had to correct — the `cnos.issues` package
  metadata shape, the "Go-source co-location ≠ declared package command"
  distinction, and the doctrine-file structure all held up through both
  rounds unmodified.
- The fixture-driven before/after board-output comparison (AC7) gave a
  clean, reproducible, non-live-API-dependent proof of behavior
  preservation in both R0 and R1 — this made the R1 repair fast to
  re-verify (β's R1 review reused the same fixture-driven method).
- `git revert` on the two R0 commits cleanly restored byte-identical
  pre-R0 file contents with correct rename detection — no manual conflict
  resolution needed, which kept the R1 repair mechanically low-risk despite
  being a full architectural reversal.
- Every AC had a locally-reproduced command/diff/log as evidence in both
  rounds, not a bare claim — this is what let β independently re-verify
  quickly rather than needing to re-derive evidence from scratch.

## Friction encountered

**The R0-relocation itself was the dominant friction of this cycle** —
γ's R0 scaffold pinned the relocation as a *binding* guardrail, framed as
"restated from the operator's clarifying comment," when it was actually γ's
own architectural judgment call (mirroring cnos#392) that ran contrary to
the operator's specific, on-topic instruction on this exact axis. α
implemented the scaffold's guardrail faithfully and correctly per what was
pinned — the defect was upstream of α's execution, in the scaffold's
characterization of the operator's comment, not in α's reading of the
scaffold. This cost a full R0→R1 repair round: two revert commits, a
doctrine-update commit, and a second full β re-review, where a correct R0
scaffold would have produced a single-round cycle. See `gamma-closeout.md`
for the triage of this finding — it is primarily a γ-scaffold-authoring
finding, not an α-execution finding, but it is the single largest source of
wasted cycle effort this cycle produced, so it belongs in this retrospective
too.

**The cdd-artifact-check (I6) section-header naming trap (R0).** While
writing `self-coherence.md` incrementally (one section per commit, per
`alpha/SKILL.md` §2.5's convention of prose-citing sections as "§Gap",
"§Skills", etc.), I used literal Markdown headers of the form `## §Gap`,
`## §Skills`, `## §ACs`, `## §CDD Trace` — matching the "§"-prefixed prose
convention used throughout `alpha/SKILL.md` and γ's own scaffold when
*citing* sections. `cdd-verify`'s `ledger.go` (`sectionPresent()`) does an
exact literal-string match against `"## Gap"`, `"## Skills"`, `"## ACs"`,
`"## CDD Trace"` (no "§") — so all four sections read as missing, and
remote CI hard-failed on four consecutive pushes (`16a0fcc6`, `4fda6fc9`,
`dc1e392c`, `e889e5ba`). Compounding factor: this cycle was classified
`"small-change"` (not `"triadic"`) by `classifyCycleType()` because
`beta-review.md` did not yet exist at that point in R0 — indistinguishable,
by file presence alone, from a genuinely small change, which routes missing
sections through `checkSmallChangeArtifacts()`'s `forUnreleased=false` path
(hard `checkFail`) rather than the `forUnreleased=true` `checkWarn` path a
recognized-as-in-progress triadic cycle gets. This was not caught by local
`cn cdd verify` runs before the first push (a working-directory mistake in
one early invocation masked the real signal) and was only caught via
**remote CI**, costing four failed pushes before the root cause was
diagnosed and fixed (renaming the six section headers to their literal
non-"§" form). See `gamma-closeout.md`'s triage table for disposition.

## Calibrated self-assessment

**R0 implementation quality (of what was actually asked): high, but the
scaffold it was faithfully implementing was wrong on one axis.** Every AC
oracle γ specified was met with real, reproducible evidence in both rounds;
no scope creep occurred in either round (no Node generator, no #216 work,
no workflow edit, no taxonomy/effort semantic change); the CI section-header
bug was self-diagnosed and fixed within the same round, with the fix
verified by 4 consecutive clean local runs plus a green remote run before
signaling review-readiness. I do not consider the R0-relocation itself an
α execution defect — the scaffold's guardrail was binding and unambiguous,
and second-guessing a scaffold's explicit architectural pin against my own
reading of the issue's more permissive body text is not α's role per the
CDD role separation (γ decides scope/shape, α implements and can escalate
via `gamma-clarification.md` if it finds a concrete blocker — I found none,
because the relocation was technically clean and internally consistent,
just wrong against the operator's specific instruction that γ's scaffold
had mischaracterized).

**R1 repair quality: high.** The revert was mechanically clean
(byte-identical restoration, confirmed via direct diff against the pre-R0
tree), the doctrine update was honest about the R0→R1 history rather than
silently erasing it, and CI was re-verified green (all 10 jobs) before
signaling review-readiness a second time. β's independent R1 re-verification
found zero new issues.

**Where I would calibrate differently in hindsight:** nothing in α's own
execution — the one thing I would flag is that α's `self-coherence.md` (R0)
treated γ's scaffold's guardrail as settled fact worth restating rather than
re-deriving from the operator's own comment text independently. α does not
own scaffold-correctness verification as a primary duty, but a cheap
"does the scaffold's characterization of the operator's comment match what
I can see by reading the comment myself" cross-check — which β actually did
perform in R0 (Finding 1) — would have been a low-cost second read α could
also have done, and did not.
