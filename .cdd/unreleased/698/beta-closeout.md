---
issue: 698
role: beta
wake: cds-dispatch
date: 2026-08-05
---

# β closeout — #698 Agent Dialogue Protocol v0

## Cycle summary

Two review rounds. R0: walked the γ scaffold's AC1–AC11 oracle list against α's first draft; found the doc's grounded content (live-ref examples, boundary-rule wording) accurate, but found a blocking defect in the citation layer — §2.3's supersession-chain table and two related citations in §2.2/§4 pointed at GitHub comment-ID permalinks that do not exist on the issue (fabricated) or existed but were misattributed to the wrong comment's content. Verdict: iterate, 5 numbered findings with the correct ID mapping supplied so α could fix without re-deriving. R1: independently re-verified (not just re-read α's self-coherence) every citation against freshly-fetched comment data, confirmed all fixes including 2 additional fabrications α's own sweep caught beyond my original 5, confirmed the repair touched only citation lines (9 lines total) with zero row-content rewrites, and reconfirmed AC10 scope. Verdict: converge.

## What went well

- The independent-verification discipline held: at both R0 and R1 I fetched primary data (`gh issue view`/`gh api .../issues/comments/<id>`) myself rather than trusting the doc's own citations or α's self-coherence claims at face value. This is what caught the defect — a review that only checked "does this AC's content sound right" without checking "does this specific citation resolve to a real, matching comment" would have converged at R0 and shipped fabricated citations.
- Grounding checks against live git refs (message envelopes, registry YAML, cursor SHAs) were straightforward to verify mechanically (diff against `origin/cn-sigma/cnos/...` etc.) and gave high-confidence pass/fail rather than a subjective read.

## What went wrong / process gap

- My R0 review found 5 numbered findings but did not run a *complete* mechanical sweep of every ID-shaped token in the document against the real comment-ID set — I found the fabrications by spot-checking the rows AC2's oracle specifically calls attention to (§2.3's table), not by exhaustively scanning. α's own R1 repair caught 2 more fabrications (line 132, lines 145–147) that were the same defect class but outside the specific lines my R0 findings named. Process fix for next time: when a finding class is "citation ID doesn't resolve," the review should immediately escalate to a full-document mechanical sweep for that class rather than fixing only the specific instances first noticed — I got the right *diagnosis* but not the most efficient *scope* of the fix request, costing what could have been a 3rd round if α hadn't independently swept.
- Recommend this become a standing β-review habit note in `beta/SKILL.md` for any cell whose AC requires citing external IDs (issue/comment/commit references): grep-extract the full ID pattern set once, cross-check the whole set, not just the AC's named example rows.

## Follow-up worth filing (not blocking this cell)

- Same CI/lint idea α's closeout names (mechanical "does this permalink resolve" check for docs citing GitHub comment permalinks) would also reduce β's review burden for this AC class. Not filed as a separate issue by this cycle; noting for operator triage.
