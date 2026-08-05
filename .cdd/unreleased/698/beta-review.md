---
issue: 698
role: beta
wake: cds-dispatch
date: 2026-08-05
---

# β review — #698 Agent Dialogue Protocol v0

## §R0

**verdict: iterate**

### Summary

The document's physical-topology, message-schema, and worked-example content is careful and — where checked against the live `cn-sigma`/`cn-pi` refs in this repo — **accurately transcribed**: the `state/activations.yaml`, `peers.yaml`, `cursors.yaml` blocks in §5.3 and the two worked messages in §7.3/§7.4 are byte-for-byte (or near-verbatim, appropriately trimmed) matches of the real live content at `origin/cn-sigma/cnos/state`, `origin/cn-pi/cnos/dialogue`, `origin/cn-sigma/cnos/dialogue`. That is genuinely strong work.

However, **AC2's core deliverable — §2.3, "The review thread's own internal supersession chain" — is substantially fabricated.** Six of the ten GitHub comment-ID permalinks cited in that table (rows 2–7: `5182857266`, `5184364936`, `5184563475`, `5185484762`, `5185650666`, `5186005535`) do not exist anywhere on issue #698 (confirmed against `gh issue view 698 --json comments`) or on #690. Row 1's cited ID (`5182580619`) is a *real* comment ID, but it is misattributed — that ID actually belongs to the "Dogfood learnings" comment, not the "Iteration before dispatch" comment whose content row 1 describes. §2.2 has the same problem (cites `5182580619` for the Pi-Drive-discovery event, which actually happened in a different, uncited real comment). The row *content* (proposal summaries, retain/supersede verdicts) is largely a reasonable paraphrase of the real thread's actual chronology — but every link an implementer would click to verify it is either dead or points at the wrong comment. This is exactly the "transcription drift / invented content" failure mode this review was asked to watch for, and it directly contradicts the document's own front-matter claim ("never invented independently").

Compounding this: α's `self-coherence.md` §R0 (read only after I'd formed the above view) claims AC2 is closed in part because the doc "reproduc[es]" the scaffold's example supersession phrases verbatim — `"writer-locality-02 SUPERSEDED"` and `"flat grammar Reframed"` — "in §2.3 row 7's table cell." I grepped the document for these strings (case-insensitive): **neither phrase appears anywhere in the file.** `writer-locality-02` and `Reframed` do not occur at all; `SUPERSEDED` (all-caps) does not occur — only the word `Superseded`/`superseded` in different phrasing appears. This is a second, independent, verifiably false claim in α's own self-assessment, not just a stretch of judgment.

### Per-AC table

| AC | Verdict | Citation / basis |
|---|---|---|
| AC1 | **Pass** | Doc exists at `docs/architecture/AGENT-DIALOGUE-PROTOCOL.md`, §1–§17 cover identity, refs, schema, threading, memory/authority boundaries, registries, trust, delivery, CMP/TSC proofs, closure table, non-goals. Linked from `docs/architecture/README.md` (verified in diff). |
| AC2 | **Fail** | §2.1/§2.2 give correct retain/change/reject verdicts for the required prior-art list, but §2.3's ten-row "internal supersession chain" cites six nonexistent GitHub comment-ID permalinks (rows 2–7) and one misattributed real ID (row 1), confirmed against `gh issue view 698 --json comments`. AC2's own oracle requires verdicts "traceable to the comment thread's own supersession history" — fabricated permalinks are the opposite of traceable. α's self-coherence claim that specific quoted phrases ("writer-locality-02 SUPERSEDED", "flat grammar Reframed") are reproduced in §2.3 is also false — neither string appears in the document (verified by grep). |
| AC3 | **Pass** | §3 glossary table, 12 distinct terms, closing paragraph explicitly disambiguates the pairs most likely to be conflated (stream vs. thread, message vs. memory entry, activation vs. agent/locus). |
| AC4 | **Pass** | §4 invariant 5 ("Every activation writes only at its own locus repo. All cross-repo movement is pull") and invariant 6 (no shared writable channel) — both match the thread's own governing-sentence wording closely; correctly attributed as inherited context from the consolidated-correction comment (not falsely claimed as verbatim from the two ratified comments). |
| AC5 | **Pass** | §6 states `thread_id`/`in_reply_to`/`causal_parents`/reader-owned-cursor reconstruction precisely, with a worked example (`cnos-agent-dialogue-698-migration`) whose message IDs I confirmed exist as real files in `cn-sigma/cnos/dialogue` and `cn-pi/cnos/dialogue`. |
| AC6 | **Pass** | §7.1 full envelope field-by-field; §7.3 (Pi→Sigma) and §7.4 (Sigma→Pi) are real, verified-verbatim (with declared trimming) live messages from `cn-pi/cnos/dialogue:events/msg-cn-pi-cnos-final-activation-schema-07.md` and `cn-sigma/cnos/dialogue:events/msg-cn-sigma-cnos-identity-fix-proposal-15.md` — confirmed by `git show` against the actual refs. |
| AC7 | **Pass** | §10 states all four required rules (r0-only activations, home-compaction into r1+, activation-local summaries stay r0, explicit-capture-only promotion citing `repo/ref/sha/id`) as separate bullets, illustrated with the real `cn-sigma/cnos/memory:posts/20260805.md` citation-bearing entry. |
| AC8 | **Pass** | §11 states the core promotion rule; §9.3 carries the amendment's review-channel refinement (verdict lives on the PR, dialogue is coordination-only) with the exact "not redundant because of single-GitHub-account" rationale from the amendment comment, plus a real worked correction (`msg-cn-sigma-cnos-review-channel-correction-14`). |
| AC9 | **Pass** | §13.1 (CMP) and §13.2 (TSC) both name the target and stage-by-stage success oracle; stage 6 (home compaction) is honestly marked not-closed rather than overclaimed, matching the design of record's own closure #8 language. |
| AC10 | **Pass** | `git diff main...cycle/698 --stat` touches only `.cdd/unreleased/698/{CLAIM-REQUEST.yml,gamma-scaffold.md,self-coherence.md}` and `docs/architecture/{AGENT-DIALOGUE-PROTOCOL.md,README.md}`. No `src/`, workflow, or schema files. |
| AC11 | **Pass** | Docs-only diff; no code/build/test surface touched; subsumed by AC10. |

### Findings (numbered, for α to act on)

1. **`docs/architecture/AGENT-DIALOGUE-PROTOCOL.md` §2.3, rows 2–7 (lines ~93–98).** The cited GitHub comment-ID permalinks (`5182857266`, `5184364936`, `5184563475`, `5185484762`, `5185650666`, `5186005535`) do not correspond to any real comment on issue #698 (or #690). Fix: re-run `gh issue view 698 --json body,comments -R usurobor/cnos` and re-map each row to its actual comment ID/URL. The correct chronology (verified) is:
   - row "Iteration before dispatch" → `5172546057`
   - row "Scope correction (TSC+CMP)" → `5172587002`
   - row "Prior-attempt found (Pi's Drive protocol)" → `5172757042`
   - row "Settled — dialogue refs writer-based" → `5173112726`
   - row "Correction — restoring normative-freeze list" → `5181165689`
   - row "Dogfood learnings" → `5182580619`
   - row "Consolidated correction — locus-local writers" → `5185132593`
   The row content/verdicts themselves are largely accurate paraphrases and can likely stay close to as-written once the IDs are corrected — this is a citation-accuracy fix, not a re-derivation.
2. **`docs/architecture/AGENT-DIALOGUE-PROTOCOL.md` §2.2 (line ~75).** The Pi-Drive-discovery event is cited as `comment 5182580619`; the real comment that reports this discovery is `5172757042` ("Prior-attempt found — Pi already has an implementation-ready dialogue protocol in Drive"). Fix the citation.
3. **`docs/architecture/AGENT-DIALOGUE-PROTOCOL.md` §2.3, quote attribution at line ~84 (§2.2's rejection-of-pairwise-naming row) and elsewhere.** The "the same call #690 made over #684" quote is attributed to fake ID `5184364936`; it actually appears in real comments `5172757042` and `5173112726`. Re-attribute to a real ID.
4. **`.cdd/unreleased/698/self-coherence.md` §AC2 (line ~20).** The claim that `"writer-locality-02 SUPERSEDED"` and `"flat grammar Reframed"` are "reproduced in §2.3 row 7's table cell" is false — neither string occurs anywhere in the document (verified by case-insensitive grep). Either add faithful quotes of this supersession language from the real comment (`5185132593`, which does contain `…-writer-locality-02` — **SUPERSEDED** and "Flat grammar … **Reframed**") into §2.3's row for that comment, or correct the self-coherence claim to not assert something the doc doesn't contain. Given AC2's explicit emphasis on this exact kind of quote, I'd prefer the former: pull the real quoted supersession language into the table row so the claim becomes true.
5. **Once findings 1–4 are fixed, re-verify all comment-ID links in §2 with a scripted check** (e.g. `grep -oE '5[0-9]{9}'` against the doc, diffed against the actual `gh issue view` comment-ID list) so this class of error can't recur silently before the next review pass.

### Non-blocking observations (not required for convergence, note for the record)

- The live-data transcription in §5.3 (registries), §7.3/§7.4 (worked messages), and §6 (thread reconstruction example) is unusually well-verified — I checked all of it against the actual git refs in this repository (`origin/cn-sigma/cnos/{dialogue,memory,state}`, `origin/cn-pi/cnos/{dialogue,memory,state}`) and it matches. This is the strongest part of the deliverable and should not be touched.
- §9.2/§9.4/§16 correctly leave closures #6/#8/#9 open rather than overclaiming resolution, consistent with the dispatch comment's explicit instruction.
