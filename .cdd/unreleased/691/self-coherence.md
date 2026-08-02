# α Self-Coherence — cnos#691

## §R0 — AC-by-AC verification evidence

### AC1 — `MEMORY.md` rewritten to the ranked model; triadic framing retired

**Change:** `docs/reference/runtime/MEMORY.md` was rewritten from "Memory in cnos — lean triadic model" (v0.2.0) to "Memory in cnos — ranked model" (v0.3.0). The episodic/reflective/working-continuity ontology (three fixed named classes bound to `threads/adhoc/`, `threads/reflections/`, `state/conversation.json`) is removed entirely — not relabeled. The doc now states #690's ranked model verbatim: r0 = raw evidence (writer-owned orphan-ref boxes), rN = home-compacted rollups that `reads:` cite rank N-1, cadence (daily=r1/weekly=r2/monthly=r3) explicitly called out as operational not structural, and "Promotion ≠ rank" as its own numbered point (§ The model, point 6).

**Oracle evidence:**
```
$ grep -n "triadic\|α — episodic\|β — reflective\|γ — working continuity" docs/reference/runtime/MEMORY.md
(no output — exit 1, zero hits)
```
- `r0`/`rN` rank model: stated in "Principle" and "The model (KISS/YAGNI)" §§1-2.
- `reads:` provenance field: stated as point 5 ("Provenance... `reads:`") and in the "Entry schema" YAML block.
- "Promotion ≠ rank" stated as its own bolded point 6, with cross-reference to the essay's "Promotion is not a rank" section.
- **Status header** now reads: "Canonical doctrine, ratified by #690... The runtime implementation of the box topology... does **not** exist yet — that is #690 Sub 2... through Sub 5... future work. This document states the doctrine that implementation must converge on, not a description of what already runs." This distinguishes doctrine-canonical (now) from implementation-status (future) per the oracle's explicit instruction not to overclaim.

**Verdict: PASS.**

### AC2 — `AGENT-MEMORY-LOG-STRUCTURED.md` reconciled

**Change:** Light-touch reconciliation, not a rewrite (per friction note 2 / gamma-scaffold §3). Added:
1. A "Doctrine cross-reference" callout immediately after the title, stating `MEMORY.md` is the concrete canonical doctrine this essay's law feeds (ratified #690, 2026-08-02); clarifying that the essay's richer illustrative schema (`inputs:`/`basis:`/`coherence:`) and `MEMORY.md`'s minimal `reads:` field name the same provenance requirement, not two different structural models; and clarifying that the essay's generic "activation logs" references describe the r0-substrate *role*, not the specific `AGENT-ACTIVATION-LOG-v0.md` convention (which the essay itself cross-links, now flagged superseded-for-memory).
2. "Related analogues" section: added `MEMORY.md` as the lead canonical-doctrine entry, and annotated the `AGENT-ACTIVATION-LOG-v0.md` link as "historical convention, superseded for agent-memory purposes by `MEMORY.md`."

The essay's own rank law (r0/r1/rN, "Rank is not frequency," "Promotion is not a rank," citation-bearing compaction) was already aligned with #690 and is untouched — no contradiction existed there before this cell, per friction note 2.

**Oracle evidence:**
- Rank terminology check: both docs use `r0`, `rN`/`r1`/`r2`/`r3`, "promotion ≠ rank" / "Promotion is not a rank," and a provenance/citation requirement (`reads:` in MEMORY.md; `inputs:`/citation discipline in the essay) — now explicitly reconciled as the same requirement via the new cross-reference note, not left as an unstated coincidence.
- Cross-reference: essay now opens with a link to `docs/reference/runtime/MEMORY.md` and closes "Related analogues" with the same link — no longer reads as a free-standing, disconnected philosophy piece.
- No factual claim survives unreconciled with AC3: the essay's mentions of "activation logs" (e.g. "The current agent activation log shape is already a good `r0` substrate," "Activation logs can still own cross-body continuity") are scoped by the new top-of-file note to describe the r0-substrate *role* generically, not to assert the `AGENT-ACTIVATION-LOG-v0.md` convention's specific ref/path scheme is still authoritative — consistent with that doc's new superseded-for-memory framing (AC3).

**Verdict: PASS.**

### AC3 — `AGENT-ACTIVATION-LOG-v0.md` narrowed with supersession framing

**Change:** Frontmatter `status:` field rewritten from `field convention` to `historical convention — superseded for agent-memory purposes by cnos#690 (docs/reference/runtime/MEMORY.md); §0/§0.1 writer-locality and wake-class-ownership mechanics remain current`; added a new `superseded_for_memory_by:` frontmatter list pointing at `MEMORY.md` and `cnos#690`. Added a three-paragraph prose notice immediately after the title, before §0, stating: (a) agent-memory continuity is now governed by the box-topology model in `MEMORY.md`, not this doc's `.cn-{agent}/logs/` + `.cn-{agent}/threads/activations/` mirror; (b) §0/§0.1 are explicitly *not* superseded — they document live, currently-operating write-discipline mechanics orthogonal to where memory lives; (c) §1 onward is kept as accurate historical record, not current/future source of truth.

**Oracle evidence:**
```
$ ls -la docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md
-rw-r--r-- 1 runner runner 22359 Aug  2 01:59 docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md
```
File still exists at the same path — not deleted.
- Frontmatter `status:` and a new prose header both state the supersession explicitly with a direct link (`../runtime/MEMORY.md`) and issue reference (`cnos#690`).
- §0 ("Writer Locality invariant") and §0.1 ("Wake-class writer ownership") bodies are untouched (byte-identical to pre-edit) — confirmed by `git diff` showing only additions before §0, no changes within §0/§0.1 text.
- No sentence in the document reads as "this is the current/only way agent memory persists across activations" without the superseded caveat now attached at the top — the caveat precedes all of §0 onward and explicitly scopes which sections remain current (§0/§0.1) vs. historical (§1+).

**Verdict: PASS.**

### AC4 — box topology stated authoritative; #684/#688 named as prior exploration

**Change:** `MEMORY.md`'s "Ref topology (box model; authoritative go-forward shape, per #690)" section states the ref paths verbatim from #690's issue body (`refs/heads/sigma/<activation-id>` per cnos/bumpt/tsc, `refs/heads/sigma/box`+`.../cloud` for cmp, `refs/heads/sigma/home` for cn-sigma), plus the orphan/single-writer/fast-forward-only/no-force-push/no-delete-while-registered invariants, plus an explicit "this topology is doctrine... not yet-implemented runtime state" scope sentence. The top-of-doc "Prior exploration, subsumed (not controlling)" paragraph names #684 and PR #688 explicitly, quotes κ's operator-clarification language ("prior substrate exploration, not the controlling topology... closed as subsumed"), and states the salvaged artifacts carry into Sub 2 while the direction-based ref scheme itself does not carry forward. Also listed under "Related" as "prior substrate exploration... subsumed, not controlling; closed."

**Oracle evidence:**
```
$ grep -n "684\|688" docs/reference/runtime/MEMORY.md docs/papers/AGENT-MEMORY-LOG-STRUCTURED.md docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md
docs/reference/runtime/MEMORY.md:7:**Prior exploration, subsumed (not controlling):** #684 and its held PR #688 explored a direction-based ref scheme...
docs/reference/runtime/MEMORY.md:82:- Issue #684 / PR #688 — prior substrate exploration (direction-based ref scheme); subsumed, not controlling; closed.
```
(0 hits in the other two docs — expected; MEMORY.md is "most naturally" the location per the oracle's own wording, and is where the box topology is stated.)
- Box topology ref paths in `MEMORY.md` match #690's issue body's "Ref topology" section verbatim (cross-checked against the issue body fetched via `gh issue view 690`).
- Scope discipline maintained: the topology is introduced as "the ratified go-forward shape," with an explicit sentence that migration onto it is Sub 2/3/5 future work, not already-implemented runtime state.

**Verdict: PASS.**

### AC5 — no residual contradiction across the three docs

**Cross-doc read:** Re-read all three files together post-edit.
- `MEMORY.md`: ranked model, box topology, promotion-is-a-separate-event, explicitly marks `AGENT-ACTIVATION-LOG-v0.md` as superseded-for-memory in its own "Supersedes" header and "Related" section.
- `AGENT-MEMORY-LOG-STRUCTURED.md`: same rank vocabulary, now explicitly cross-linked to `MEMORY.md` as the concrete doctrine, and its activation-log references scoped to not contradict AC3.
- `AGENT-ACTIVATION-LOG-v0.md`: explicitly marks itself superseded-for-memory-purposes with a direct pointer back to `MEMORY.md`/#690, while its §0/§0.1 mechanics (which the other two docs do not touch or contradict) remain live.
No mention of "memory model," "canonical memory," or "activation log" in any of the three docs is left asserting the retired triadic model or the un-caveated activation-log-as-current-memory-mechanism claim.

**Repo-wide triadic grep (informational, out of scope per issue's explicit Out scope):**
```
$ grep -rl "triadic" docs/ | sort
docs/architecture/CELL-RUNTIME.md
docs/concepts/doctrine/coherence-for-agents/COHERENCE-FOR-AGENTS.md
docs/development/board/board-data.json
docs/development/board/index.html
docs/papers/CELL-OF-CELLS.md
docs/papers/COHERENCE-SYSTEM.md
docs/papers/CTB-LANGUAGE-SPEC-v0.2-draft.md
docs/papers/CTB-v4.0.0-VISION.md
docs/papers/DECREASING-INCOHERENCE.md
docs/papers/EXECUTABLE-SKILLS.md
docs/papers/FOUNDATIONS.md
docs/papers/SHIPPING-SOFTWARE-AFTER-AI.md
docs/reference/README.md
docs/reference/ctb/README.md
docs/reference/ctb/SEMANTICS-NOTES.md
docs/reference/governance/GLOSSARY.md
docs/reference/runtime/ORCHESTRATORS.md
```
Spot-checked each hit's context: all 17 refer to the unrelated TSC/CTB "triadic carrier" concept (`tri()`, α/β/γ *coherence-measurement* axes — pattern/relation/process) — a different, still-current doctrine (TSC's triadic coherence check), not the retired MEMORY.md episodic/reflective/working-continuity memory split. None of these dangle or reference the retired memory model. No follow-up needed; confirmed out of scope per the issue's explicit **Out** scope (cn-sigma tree changes / compactor mechanism) and not touched.

**Diff-stat scope check:**
```
$ git diff --stat
 docs/papers/AGENT-MEMORY-LOG-STRUCTURED.md         |   5 +-
 .../conventions/AGENT-ACTIVATION-LOG-v0.md         |  11 +-
 docs/reference/runtime/MEMORY.md                   | 356 ++++-----------------
 3 files changed, 69 insertions(+), 303 deletions(-)
```
Only the three named doctrine docs are modified. No `.go` file, no `.cn-sigma/` path, no CI config, no FSM/label file touched. `.cdd/unreleased/691/gamma-scaffold.md` (122 lines) is already on `cycle/691` from δ's prior commit (`fefea38`), preceding this α pass; this cell's own `.cdd/unreleased/691/self-coherence.md` (this file) is the only additional artifact this α pass adds.

**Verdict: PASS.**

## Summary

All five ACs pass. The three docs now present one coherent doctrine: ranked (`r0`/`rN`) memory with box topology, `reads:` provenance, promotion-as-a-separate-event, and `AGENT-ACTIVATION-LOG-v0.md` narrowed to historical/superseded-for-memory-purposes while its live writer-locality mechanics (§0/§0.1) stay intact. No scope break encountered — no code, `.cn-sigma/`, CI, or FSM/label change was needed to make the doctrine coherent.
