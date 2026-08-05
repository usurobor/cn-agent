---
issue: 698
cycle_branch: cycle/698
main_sha: 7f249ddbb50f230d5d41287b6554ab17b5a1d1d5
protocol: cds
cell_kind: docs
---

# γ scaffold — #698 Agent Dialogue Protocol v0

## 1. Mode / gap / deliverable (restated from the issue)

Mode: design-doc / docs-only. Deliverable: `docs/architecture/AGENT-DIALOGUE-PROTOCOL.md`, structured to the issue body's §1–13 + its acceptance criteria (AC1–AC11), filled from the two ratified design-of-record comments:

1. **Design of record** — [comment 5193497595](https://github.com/usurobor/cnos/issues/698#issuecomment-5193497595), "FINALIZED DESIGN — ratified 2026-08-05 (activation = {agent, locus})".
2. **Amendments** — [comment 5195256310](https://github.com/usurobor/cnos/issues/698#issuecomment-5195256310), "Design-of-record amendments — 2026-08-05 (operator-ratified)".

Per the dispatch comment: **transcribe, do not re-derive.** The doc is a synthesis/transcription of already-ratified content, not new design work. Open items (closure #6, #8, #9) are noted as pending/deferred, not blockers.

## 2. Source-of-truth table

| Content area | Source |
|---|---|
| Identity model (`activation = {agent, locus}`), refs (dialogue/memory/state), message envelope, registries, trust modes, memory/authority boundary, closure table | Comment 5193497595 (design of record) |
| Nomenclature (ADR/CDD spelled out), review-channel boundary (PR = authority, dialogue = coordination only), signed-activation trust mechanism, account-topology closure #9 | Comment 5195256310 (amendments) |
| Prior-attempt review (OCaml-era, activation-log v0, #684/#688, #690, TSC-Pi proposal) | Issue body AC2 + comment thread (2026-08-03 through 2026-08-04 iteration comments) — retain/change/reject decisions are visible across the thread's supersession table in the "Consolidated correction" comment |
| Grounding / live refs | `cn-sigma/cnos/{dialogue,memory,state}`, `cn-pi/cnos/{dialogue,memory,state}` per the dispatch comment |
| Non-goals, success/closure condition | Issue body |
| Existing conventions this doc must not silently contradict | `docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md`, `docs/reference/runtime/MEMORY.md` |

α MUST read the full issue body and **all 11 issue comments** (not just the two cited as authoritative) before writing — the two ratified comments are the *controlling* content, but the intermediate comments carry the retain/change/reject history AC2 requires the doc to review, including explicit supersession chains (e.g. "writer-locality-02 SUPERSEDED", "flat grammar Reframed").

## 3. Per-AC oracle list (β walks this independently)

- **AC1 — Design doc exists at the correct path and covers the doc's own outline.** Oracle: `docs/architecture/AGENT-DIALOGUE-PROTOCOL.md` exists; it has numbered sections §1–13 (or a documented equivalent structure covering: identity, refs/topology, message schema, thread reconstruction, memory boundary, project-authority boundary, registries, trust modes, delivery semantics, TSC/CMP round-trip proofs, review-channel boundary, closure table, non-goals).
- **AC2 — Prior attempts reviewed.** Oracle: doc names OCaml-era design, activation-log v0, #684/#688, #690, TSC-Pi proposal, each with an explicit retain/change/reject verdict traceable to the comment thread's supersession history.
- **AC3 — Object model unambiguous.** Oracle: doc has a glossary or equivalent section distinctly defining agent, activation, home, stream, thread, cursor, message, memory, promotion, project authority — no two terms conflated.
- **AC4 — Writer-local stream topology specified.** Oracle: doc states every activation writes only at its own locus repo (no cross-repo writes; all cross-repo movement is pull), and explicitly rejects a shared mutable channel.
- **AC5 — Thread reconstruction specified.** Oracle: doc explains how a logical `thread_id` spans multiple writer streams via `in_reply_to`/`causal_parents` and reader-owned cursors.
- **AC6 — Message schema concrete.** Oracle: doc gives the full `cnos.agent-message.v1` envelope field-by-field, plus at least two worked examples: Sigma→Pi and Pi→Sigma (CMP box/cloud and/or TSC).
- **AC7 — Memory boundary explicit.** Oracle: doc states activations write r0 only, home writes r1+ via compaction, activation-local summaries are not canonical r1, and dialogue transcripts are never copied wholesale into memory (only explicit-capture promotion citing exact `repo/ref/sha/entry-id`).
- **AC8 — Project authority boundary explicit.** Oracle: doc states channel messages govern nothing until promoted into a project-native artifact (issue, ADR, CDD receipt, spec, reviewed PR, commit); review verdicts specifically live on the PR, not on dialogue (per the amendment).
- **AC9 — TSC round-trip test specified.** Oracle: doc names the first implementation target(s) (CMP box↔cloud and/or TSC Sigma↔Pi) and each one's success oracle.
- **AC10 — No implementation changes.** Oracle: `git diff main...cycle/698 --stat` touches only `docs/` (and `.cdd/unreleased/698/` cycle artifacts) — no `src/`, no workflow, no schema, no repo-migration files.
- **AC11 — Current gates remain green.** Oracle: this is a docs-only cell; no code/build/test gates are affected. β confirms no non-docs files are touched (subsumed by AC10) and that no CI-relevant surface changed.

## 4. Scope guardrails (from the issue's Non-goals + Active design constraints)

- Do NOT implement the protocol (no code).
- Do NOT migrate cn-sigma/cn-pi/tsc/cmp.
- Do NOT build pi-host.
- Do NOT create a router.
- Do NOT add cryptographic signatures as a *requirement* (the amendment's `signed-activation` mode is a documented future trust mode, not something to build here).
- Do NOT alter #690.
- Do NOT start Demo 0.
- One primitive: single-writer append-only streams. No generic router in this cell.

## 5. α prompt

You are α (implementer) for cell #698 on branch `cycle/698`. Read the full issue body of https://github.com/usurobor/cnos/issues/698 and **all 11 comments** on it (`gh issue view 698 --json body,comments` or equivalent). The two comments titled "FINALIZED DESIGN — ratified 2026-08-05" and "Design-of-record amendments — 2026-08-05" are the controlling content — transcribe them faithfully, do not re-derive or invent new design. The earlier comments in the thread give you the prior-attempt review history AC2 requires (what was proposed, corrected, and superseded, and why) plus the CMP/TSC grounding examples.

Write `docs/architecture/AGENT-DIALOGUE-PROTOCOL.md`: an implementation-ready design doc structured to the issue's §1–13 outline, satisfying AC1–AC11 exactly as scaffolded above in §3. Do not touch any file outside `docs/` and `.cdd/unreleased/698/`. When done, write `.cdd/unreleased/698/self-coherence.md` §R0: walk each AC yourself, cite the doc section/line that satisfies it, and flag any AC you could not fully close (with why). Commit and push to `cycle/698`.

## 6. β prompt

You are β (reviewer) for cell #698 on branch `cycle/698`. You did not write the doc. Independently read the issue (body + all comments) and `docs/architecture/AGENT-DIALOGUE-PROTOCOL.md` as α left it. Walk the AC oracle list in §3 above one by one — for each AC, find the specific section/passage that satisfies it (or does not). Pay particular attention to AC2 (does the doc actually cite retain/change/reject verdicts, not just list prior art?), AC4/AC7/AC8 (the three boundary rules that are easy to state loosely and hard to state precisely), and AC10 (diff scope). Write `.cdd/unreleased/698/beta-review.md` §R0 with a `verdict: converge` or `verdict: iterate` and, on iterate, a numbered findings list α can act on. Commit and push to `cycle/698`.

## 7. Friction notes

- The issue's own body opens with "**Not dispatched** — dispatch on explicit operator authorization"; the final issue comment ("Dispatching (operator-authorized, 2026-08-05)") is the authorization. Recorded in the claim comment; not itself an AC, but worth α/β knowing the dispatch is deliberate and recent, not stale.
- Closure #6 (identity/trust under one GitHub account) and #9 (account topology) are explicitly deferred to #701/#702 per the dispatch comment — the doc should note them as pending, not attempt to resolve them.
