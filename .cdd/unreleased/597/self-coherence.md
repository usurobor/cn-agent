# Self-coherence — #597

## §R0

**Claim:** AC1–AC4 (per #597's body) are satisfied by pre-existing comments; the only residual work at dispatch-comment time (patch #596's dispatch-order tree) had already been independently completed by κ's 10:16:14Z refresh before this cell claimed the issue; but that same refresh left a newly-stale runtime-blocker reference (PR #602 merged at 10:59:35Z, after the refresh).

**Evidence:**
- `gh issue view 597` comments: AC1–AC4 evidence posted 2026-07-05T03:44:13Z and 08:25:08Z (both before this cell's claim at ~12:13Z).
- `gh issue view 596` — `updatedAt: 2026-07-05T10:16:14Z`, body already contains `#597 → #598 → Sub 0.5 → #452 → #453 → operator gate` dispatch order.
- `gh pr view 602` — `mergedAt: 2026-07-05T10:59:35Z`, state `MERGED`, title confirms it reverts the broken `#594` model pin.
- `grep -n model .github/workflows/cnos-cds-dispatch.yml` on local `main` (fast-forwarded to `origin/main`) — no rejected model pin present, consistent with #602 having merged.

**Falsification check:** ran one command designed to disprove "the tree is still broken" (`gh issue view 596 --json body`, read in full) before proposing any patch — confirmed the tree fix was already live, so I did not re-do already-done work. Ran a second command (`gh pr view 602`) to disprove "the runtime is still down" before writing any blocker-status prose.

**Governing gap identified:** not the originally-named tree defect (moot), but staleness introduced by the timing gap between κ's refresh (10:16:14Z) and #602's merge (10:59:35Z) plus this cell's own claim (12:13Z) — #596 asserted a blocked state that no longer held.

**Mode:** MCA — patched #596's body directly (4 targeted line replacements, diffed and verified minimal before applying).

## Debt / follow-ups

None named. This is a bounded reconciliation; #598/Sub 0.5/#452/#453 dispatch order is unaffected by this cell's changes.
