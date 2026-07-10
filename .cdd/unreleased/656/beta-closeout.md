# β close-out — cycle #656

## Verdict: converge (single round, R0)

## Summary

Independent verification (adversarial-review subagent + hand re-derivation
of the CUE closedness claims + a fresh end-to-end smoke test, not a re-read
of α's self-report) found three real correctness bugs. All three are fixed
on this branch with regression tests that fail against the pre-fix code:

1. Ledger silently dropped a removed dispatch workflow (`DriftRemoved`
   added) — found by α's own pre-review self-check.
2. A stale workflow file from a prior successful install could get its
   ledger render-contract metadata silently corrupted by a later, failed,
   different-identity install attempt — found by independent adversarial
   review. Fixed by threading `runDispatchCds`'s own `rendered` bool
   through instead of re-deriving "did the render happen" from file
   existence.
3. A confirmed sha256 mismatch (real drift) could be silently reported as
   "no drift" if the classification step itself failed — found by the
   same adversarial review. Fixed with a new `DriftUnclassified` value
   that still counts toward the top-level `Drift` bool.

Full finding writeups + fix evidence: `beta-review.md`.

## Verification performed independent of α's narrative

- Full clean-state rebuild + test + vet, both affected Go modules.
- CUE self-test re-run.
- A live, from-scratch end-to-end smoke test (real `cn build`, real local
  index, real `cn repo install`, real `cn repo status`) — not scripted
  from α's own transcript.
- A dedicated adversarial-review subagent, briefed on the issue/design
  context but NOT on α's self-coherence narrative, specifically hunting
  for correctness bugs in the drift-classification and ledger-write logic.
- Hand-crafted an ad hoc CUE fixture to independently confirm the
  closedness-enforces-timestamp-freeness claim, rather than trusting the
  checked-in fixture corpus alone.

## Non-blocking observations (recorded, not iterate-worthy)

- `sha256File` duplicated between two packages (~8 lines each) — not
  worth a shared package yet.
- A `deps.lock.json`-present-but-`deps.json`-absent package drift
  direction isn't surfaced by any field — narrow, requires a hand-edited
  manifest, out of this phase's literal Acceptance text.
- `cmd_repo_status.go`'s presentation layer has no line-by-line stdout
  test — consistent with `hubstatus.Run`'s own precedent elsewhere in this
  codebase; the underlying `Status` struct is exhaustively tested.
- `cn repo doctor` is named in the design doc's Phase 1 command-surface
  table but claimed by none of the wave's six filed sub-issues (#657–
  #661) — not a gap in THIS cycle's own Acceptance text, and not filing a
  new issue for it here since it's speculative which future phase should
  own it; flagging for operator/planner awareness rather than unilaterally
  scoping a new phase.

## Follow-up issue check

Searched existing issues before considering any new filing
(`gh issue list --search "repo doctor"`, `--search "cn repo lifecycle"`).
All six wave phases already have dedicated sub-issues (#656 this one,
#657–#661). No new issue filed — every Debt item above is either already
covered by a later phase's existing sub-issue, or narrow/speculative
enough that filing now would be premature (mirrors the cnos#608 precedent
of reasoning through non-filing rather than filing reflexively).
