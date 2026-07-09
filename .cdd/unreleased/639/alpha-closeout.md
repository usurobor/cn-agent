# α close-out — cnos#639

## Scope note

This close-out is written at the **δ §9.5 converge boundary** — β
returned `verdict: converge` at R0 (`beta-review.md`, head SHA reviewed
`649c8649`), and this file is one of the three artifacts that boundary
requires (`alpha-closeout.md` + `beta-closeout.md` + `gamma-closeout.md`),
per `delta/SKILL.md` §9.5's per-R[N] artifact contract for the `converge`
row. As with cnos#626's and cnos#640's equivalent closeouts, this is
**not** the post-merge/post-release close-out `alpha/SKILL.md` §2.8
describes — no PR exists yet for `cycle/639`, and opening the cycle-PR
plus requesting the `status:review` transition is δ's next action per
§9.6, not this pass's. γ, α, and β roles for this closeout round are
collapsed into a single authoring pass per Sigma's engineering-persona
protocol commitment #5 (β-α-collapse-on-δ for skill/docs-class cycles),
since this cycle's own AC oracle is mechanical (grep counts, byte-diffs,
table cross-checks) and β's R0 review (zero findings, `verdict:
converge`) is already independently on record.

## Summary — what was built

cnos#639 asked for the `run_class` taxonomy that classifies each
`cds-dispatch` wake firing to be reconciled — it was inconsistently
enumerated across three doctrine surfaces (`dispatch-protocol/SKILL.md`
listing 4 values and missing `resumed_from_matter` entirely;
`cds-dispatch/SKILL.md` listing 5, the most complete but still missing an
operator-authorized "scope continuation" shape; `delta/SKILL.md`
narrating individual resume shapes with no canonical enum at all) and
incomplete — cnos#626 hit the scope-continuation shape three times across
its own R0/R1/R2 rounds and self-disclosed the gap explicitly each time
rather than silently mislabeling.

What shipped, entirely inside γ's oracle:

1. **One canonical enum, one canonical home.** `cds-dispatch/SKILL.md`
   §"Repair re-entry preflight" → Step A is now the single source of
   truth: a four-check ordered decision procedure
   (`resumed_from_matter` → `scope_continuation` → `repair_pass` →
   `first_pass` terminal), each non-terminal check carrying an explicit
   positive discriminator, plus a mutual-exclusivity table and the
   five-shape mapping table (#593, #630, #614, #626 R0, #626 R1/R2) γ's
   scaffold required be shown as a table in the doctrine itself, not just
   asserted in `self-coherence.md`.
   `dispatch-protocol/SKILL.md` §2.8 now states the full 6-value set and
   points at `cds-dispatch/SKILL.md` as canonical rather than restating a
   narrower competing list.
2. **`scope_continuation` named with defined re-entry behavior.** The
   enum gains the sixth value; `delta/SKILL.md` gains new §9.13, sibling
   to §9.10/§9.11 (§9.12 is #626's own, unrelated,
   cell/substrate-identity-boundary doctrine — untouched), stating
   explicitly that δ does NOT re-scaffold merged matter and does NOT
   treat prior `CDDArtifacts` as a conflict, citing #626 R1/R2 as the
   empirical worked example exactly as §9.10 cites cycle/497 and §9.11
   cites #630.
3. **Recovery-vs-resume resolved as one class**, explicitly, not left
   implicit: `resumed_from_matter` spans a produce/consume pair, stated
   in words rather than only implied by the existing FSM/scanner
   behavior.
4. **Reset-branch-vs-first_pass resolved**: the ordered procedure checks
   issue-level history (not branch-commit-count) before ever falling
   through to `first_pass`, citing #626 R2's concrete branch-reset-to-
   clean-main-with-prior-merged-rounds case by artifact path.

Full diff against `origin/main`: three doctrine `.md` files
(`cds-dispatch/SKILL.md`, `dispatch-protocol/SKILL.md`,
`delta/SKILL.md`), two mechanically-regenerated non-`.md` files
(`cnos-cds-dispatch.golden.yml`, `.github/workflows/cnos-cds-dispatch.yml`),
and three `.cdd/unreleased/639/` cell artifacts (γ's `CLAIM-REQUEST.yml` +
`gamma-scaffold.md`, this cycle's `self-coherence.md`) — 809 insertions /
29 deletions total, zero Go source touched, zero `transitions.json`
touched, zero label definitions added or removed.

## AC-by-AC (brief; full evidence in self-coherence.md)

- **AC1** — closed: the five concrete prior firings (not hypotheticals)
  each map to exactly one class under the new procedure, with artifact
  citations in the doctrine's own mapping table; only the terminal
  `first_pass` is exclusion-defined, the other three non-terminal classes
  each carry a named positive discriminator.
- **AC2** — closed: `scope_continuation` is in the enum with δ re-entry
  behavior at §9.13, detection-order justified against the other three
  checks (not appended without reasoning).
- **AC3** — closed: both sub-questions (recovery/resume one class;
  issue-level history gates the reset-branch case) resolved with explicit
  sentences, and the ordered procedure itself (not just prose elsewhere)
  implements the resolution.
- **AC4** — closed: `grep -n "run_class"` across all three surfaces shows
  no narrower/divergent restatement; the CI guard's exact literal
  substrings were preserved (only additive changes); golden + live
  workflow regenerated via `cn install-wake` and verified byte-identical
  to a fresh render.
- **AC5** — closed: doctrine-only diff, zero `.go` files, zero
  `transitions.json` diff, zero label changes; Go/Package/Binary gates
  confirmed running and green (no-op as expected, verified rather than
  assumed); both CI guard scripts green locally.

## Debt disclosed

- **I4 (lychee link-check) and I5 (cue frontmatter validation) not run
  locally.** Same tool-unavailability gap #626's own R1 self-coherence.md
  disclosed for the identical gates. No repo links were added or changed
  by this round's edits and no SKILL.md frontmatter block was touched
  (body prose only), so both were expected green in real CI — named as a
  gap rather than silently assumed. (β's independent close-out
  observation: this resolved cleanly — real CI ran and passed both.)
- **`cn cdd verify` (I6) and `cn build --check` (I1) not independently
  invocable in this session's environment.** The `cdd-verify` Go module
  builds and vets clean, but its binary is plugin-shaped rather than a
  standalone CLI in this sandbox, so the actual `cn cdd verify
  --unreleased` / `cn build --check` invocations were not run locally.
  (β's independent close-out observation: this also resolved cleanly —
  real CI ran and passed both.)
- **Live-firing validation is out of scope for this cell, by
  construction.** This cell documents the taxonomy; it does not wire the
  classifier to enforce it (explicit non-goal per the issue body). The
  new `scope_continuation` value and the tightened `repair_pass`
  discriminator have not yet been observed driving an actual
  dispatch-wake firing's classification decision — doctrine describing
  what a firing SHOULD do, verified self-consistent and CI-guard-
  compliant, but not yet exercised live. Mirrors #626 R1's own
  "structurally only possible after this change merges and the wake
  fires for real" disclosure pattern.
- **CELL-KINDS.md / CDD.md cross-reference — confirmed not required, not
  added.** γ's scaffold read (`run_class` is a δ/dispatch-runtime
  classification, distinct from `cell_kind`'s `FactSnapshot.CellKind`
  seam) was independently re-checked (`grep -n "run_class"` against both
  files → 0 hits in all three candidate files) rather than silently
  accepted or silently skipped.
- **Repair-contract machinery (Steps B–E) untouched.** `scope_continuation`'s
  routing explicitly bypasses Steps B–E, mirroring §9.11's pattern; this
  cell did not add or modify any Steps-B–E logic, consistent with the
  issue's "no dispatch-behavior change" non-goal.

None of the above debt items required an implementation change or a
scope adjustment — each is a disclosed, by-design limit of this cell's
scope or session environment, checked against the actual diff rather than
asserted from habit.

## No scaffold corrections needed this cycle

Unlike cnos#640 (where α self-caught two factual errors in
`gamma-scaffold.md` — a taken failure-mode-catalogue slot and a
misattributed file citation), γ's scaffold for #639 required no
correction: every line-range citation, the CI-guard's exact literal
substrings, and the five concrete prior-firing shapes named in the
per-AC oracle list held up as written when checked against the live
repository state at implementation time. The one live-state fact that did
change mid-cycle was `origin/main` advancing by four unrelated doc-only
commits (board-map regeneration) — handled by a clean rebase and a
SHA-citation re-stamp (commit `beccf098`), not a scaffold error.

## Calibrated success claim

What is verified, as of this boundary:

- AC1–AC5 are all met, per the evidence enumerated in
  `self-coherence.md` §ACs and independently re-derived by β from source
  (not from this file's or self-coherence.md's narrative) in
  `beta-review.md`.
- The CI-guard hard constraint held: no literal substring the guard
  `grep -qF`s for was renamed or removed; only additive values were
  introduced.
- The golden/live-workflow regeneration is mechanical, not hand-edited,
  independently reproduced from source by both α (at implementation time)
  and β (at review time, re-running the actual renderer rather than
  trusting a pre-existing sha comparison).

What is explicitly **not** attempted or claimed:

- No live `cds-dispatch` firing has yet exercised the new
  `scope_continuation` value or the tightened `repair_pass`
  discriminator under real operating conditions — the doctrine is inert
  prose until the next wake firing reads it post-merge, mirroring the
  same gap named for §9.10/§9.11 in the cnos#630 closeout precedent and
  for AC3/AC4 in the cnos#626 closeout precedent.
- No new Go field was added to `FactSnapshot` for "issue has prior
  merged/converged rounds" — this remains a comment/PR-history read the
  wake/δ performs today, honestly named as such rather than claimed to
  already exist as a typed field (building one would have been a Go
  change, out of scope for this doctrine-only cell).
- No PR has been opened and no merge to `main` has happened — that is
  δ's next step (§9.6 `status:review` transition), not something this
  cycle's α/β passes have done or claim to have done.

**Claim:** the taxonomy reconciliation (AC1–AC5) is correctly and
completely implemented, independently verified by β's from-scratch
re-derivation of every grep/diff/CI-guard check rather than by trusting
this file's or `self-coherence.md`'s narrative. This is a clean,
single-round R0 convergence with no debt requiring further work inside
this cell's scope. It is ready for δ to proceed to the `status:review`
transition.
