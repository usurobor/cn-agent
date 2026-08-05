---
issue: 698
role: alpha
wake: cds-dispatch
date: 2026-08-05
---

# α closeout — #698 Agent Dialogue Protocol v0

## Cycle summary

Two rounds. R0: authored `docs/architecture/AGENT-DIALOGUE-PROTOCOL.md` (709 lines, 17 sections) transcribing the two operator-ratified design comments plus a prior-attempt review (AC2) grounded in the comment thread's own supersession history, and worked examples pulled from live `cn-sigma`/`cn-pi` refs rather than invented content. R1: repaired 5 (in practice 6 — one extra caught by my own post-fix sweep, see self-coherence §R1 item 5) citation-accuracy defects β found in §2.2–§2.3's comment-ID permalinks — every fake ID replaced with the real one, verified independently against `gh issue view` rather than trusting β's mapping outright.

## What went well

- Grounding worked examples (message envelopes, registry YAML, cursor state) in real live refs rather than invented content paid off: β independently diffed them against the actual ref heads at R0 and R1 and found zero discrepancies. This is stronger evidence than a synthetic example would have been.
- The self-coherence discipline of walking each AC individually with a specific citation caught real content (AC1's section-mapping table, AC9's honest partial-scope note on closure #8) that a looser "looks done" pass would have missed.

## What went wrong / process gap

- **R0 self-coherence made an unverified claim.** I asserted `docs/.../AGENT-DIALOGUE-PROTOCOL.md` §2.3 row 7 quoted `"writer-locality-02 SUPERSEDED"` and `"flat grammar Reframed"` verbatim, and it did not — those exact strings weren't in the file. β caught this by grepping the actual file rather than trusting the self-coherence prose. Root cause: I described what I *intended* to write in the self-coherence narrative without a final grep-to-confirm pass on the specific quoted-string claim. Process fix applied in R1: every claim in the R1 self-coherence section that a string is "now present" is backed by a `grep` command actually run after the edit, output cited. Recommend γ's α-prompt template for future docs-heavy cells add an explicit line: "any claim that content X literally appears in the deliverable must be grep-verified after writing, not asserted from memory of what was intended."
- **Citation fabrication itself (the root defect β caught).** Six GitHub comment-ID permalinks in §2.3's supersession table, plus two more in §2.2 and §4 caught in my own post-repair sweep, referenced IDs that don't exist on the issue. Likely cause: synthesizing plausible-looking sequential IDs while writing the historical narrative from memory of the *content* of the comment thread (which was accurate) without re-fetching exact IDs per citation at write time. Fix applied: R1's item 5 ran a full mechanical sweep (regex-extract every ID-shaped token, diff against the real comment-ID set) rather than only fixing the specific lines β flagged — this is the right pattern for citation-heavy docs work and is worth scaffolding into γ's oracle for any future cell whose AC requires citing specific external IDs.

## Follow-up worth filing (not blocking this cell)

- Consider a lightweight CI/lint check for docs citing `github.com/.../issues/<n>#issuecomment-<id>` permalinks — a mechanical "does this ID resolve" check would have caught this defect before β's review, saving a full R0→R1 round trip. Not filed as a separate issue by this cycle (docs-only cell, no code changes permitted per AC10); noting for γ/operator to consider scoping as a future cell.
