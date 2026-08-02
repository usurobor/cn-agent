# γ Scaffold — cnos#691

**Issue:** [#691](https://github.com/usurobor/cnos/issues/691) — memory wave Sub 1: cnos doctrine alignment — one canonical model (MEMORY.md / AGENT-MEMORY-LOG-STRUCTURED / activation-log-v0).
**Mode:** doctrine-only. No code, no `.cn-sigma/` changes (explicit in the issue body).
**cell_kind:** `doctrine`.
**Branch:** `cycle/691` (created from `origin/main` at `737e573388a9a771ce2432fefaefe1327bf55d32`).
**Wake:** `cds-dispatch`, protocol `cds`, `run_class: first_pass`.
**Parent:** Sub 1 of #690 (memory reset wave). Per κ's operator clarification on #690: "cnos memory doctrine alignment ← running first" — this cell is explicitly authorized to run first, ahead of Subs 2–5.
**Supersedes (closed, named in #690 body):** #581, #528, #458, #100, #35.
**Family:** #684 (re-scoped substrate exploration, held PR #688 — do not merge, named as prior exploration only), #690 (parent wave tracker).

---

## 1. Per-AC oracle list

### AC1 — `MEMORY.md` rewritten to the ranked model; triadic/adhoc-vs-reflection framing retired

**Invariant:** `docs/reference/runtime/MEMORY.md` no longer presents the "lean triadic model" (α episodic / β reflective / γ working-continuity, three fixed named classes bound to three fixed paths) as the current or proposed model. It instead states #690's ranked model: `r0` = raw evidence (writer-owned orphan-ref boxes per the box topology), `rN` (N≥1) = compaction that reads rank N-1 and cites it via `reads:` provenance, cadence (daily/weekly/monthly) is operational not structural, and promotion into spec/state/protocol is a distinct event, not a rank increment.

**Oracle:**
- `grep -n "triadic\|α — episodic\|β — reflective\|γ — working continuity" docs/reference/runtime/MEMORY.md` returns zero hits after the rewrite (the triadic framing is retired, not merely relabeled).
- The rewritten doc states the `r0`/`rN` rank model, the `reads:` provenance field, and "promotion ≠ rank" explicitly (each as its own statement, not implied).
- The doc's `Status` header reflects that this is now the canonical doctrine (not "Draft (proposed — not implemented)" carried over unchanged) — word it accurately: the *doctrine* is canonical per #690's ratification; the *runtime implementation* (box topology, cursors, rollups) is still Sub 2–5 future work. Do not overclaim implementation status.

### AC2 — `AGENT-MEMORY-LOG-STRUCTURED.md` reconciled (aligned with the spec, not a competing essay)

**Invariant:** `docs/papers/AGENT-MEMORY-LOG-STRUCTURED.md` (the rN essay) is brought into alignment with #690's ranked model rather than left standing as an independent, potentially-divergent essay. The essay already uses `r0`/`r1`/`rN` rank language close to #690's — the reconciliation is confirming/tightening consistency (terminology, the `reads:`/provenance framing, "promotion ≠ rank"), not a full rewrite, and explicitly cross-linking it to `MEMORY.md` as the concrete doctrine it feeds.

**Oracle:**
- The essay and `MEMORY.md` use the same rank terminology (`r0`, `rN`, "promotion ≠ rank", provenance/citation requirement) without contradiction — no residual place where the essay implies a *different* structural model than `MEMORY.md` states.
- The essay explicitly cross-references `MEMORY.md` (or the box topology it now describes) rather than reading as a free-standing, disconnected philosophy piece.
- No factual claim in the essay (e.g. describing the activation-log convention as authoritative for agent memory) survives unreconciled with AC3's changes to `AGENT-ACTIVATION-LOG-v0.md`.

### AC3 — `AGENT-ACTIVATION-LOG-v0.md` narrowed/replaced: the box model supersedes main-tree activation logs for agent memory; v0 kept only as a historical convention with a superseded-by pointer

**Invariant:** `docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md` no longer presents itself as the current mechanism for agent memory continuity. It is explicitly marked historical/superseded, with a pointer to #690 (and this doc's rewritten `MEMORY.md`) as the successor doctrine for agent memory. The `§0`/`§0.1` writer-locality and wake-class-ownership material may stay factually descriptive of what currently runs in production (it is not being deleted — it documents real, still-operating mechanics), but the frontmatter/status and framing must not claim this convention is *the* memory model going forward.

**Oracle:**
- The doc's frontmatter (`status:`) and/or a new prose header states explicitly that this convention is superseded for agent-memory purposes by #690's box model, with a direct link/reference.
- No sentence in the doc still reads as "this is the current/only way agent memory persists across activations" without the superseded caveat attached.
- The doc is not deleted (historical value + still-descriptive-of-live-mechanics per the issue's own "kept only as a historical convention" instruction) — confirm the file still exists at the same path after the edit.

### AC4 — box topology stated authoritative; #684/#688 named as prior exploration

**Invariant:** Wherever the rewritten docs describe the go-forward memory mechanism, they state the box topology (`refs/heads/sigma/<activation-id>` per-locus write-local r0 boxes; `refs/heads/sigma/home` as home's own r0 box; home as the sole rN compactor) as the authoritative topology, and explicitly name #684 (and its held PR #688) as prior substrate *exploration* — not the controlling topology, per κ's operator clarification on #690 ("#684/#688 are prior substrate exploration, not the controlling topology; their salvage... carries into Sub 2. #684/#688 closed as subsumed.").

**Oracle:**
- `grep -n "684\|688" docs/reference/runtime/MEMORY.md docs/papers/AGENT-MEMORY-LOG-STRUCTURED.md docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md` shows #684/#688 named explicitly as prior exploration, superseded/subsumed — not silently absent, and not described as still-controlling.
- The box topology (ref paths, single-writer/orphan/fast-forward-only invariants) is stated in at least one of the three rewritten docs (most naturally `MEMORY.md`) as the authoritative go-forward shape, sourced from #690's issue body.
- Scope discipline: name the topology as doctrine (what will be true), not as already-implemented runtime state — Sub 2 (cn-sigma dry-run migration map) and later subs are the implementation; this cell is doctrine-only.

### AC5 — no residual contradiction: no doc still presents triadic-memory or main-tree activation logs as the *current* model

**Invariant:** After AC1–AC4 land, a reader following only these three docs (plus their cross-references) does not encounter a contradiction about what the current/authoritative cnos agent-memory model is. This is the cross-cutting consistency check across all three files together, not a fourth independent edit.

**Oracle:**
- Re-read all three files together after the edits; confirm every mention of "memory model," "canonical memory," or "activation log" in any of the three is consistent with: ranked (`r0`/`rN`) model, box topology, promotion-is-a-separate-event, `AGENT-ACTIVATION-LOG-v0` superseded-for-memory-purposes.
- `grep -rn "triadic" docs/` (repo-wide, not just the three files) to catch any other doc that might cross-reference the retired triadic framing and now dangle — if found, note as a follow-up (out of scope to fix elsewhere unless trivial, per the issue's explicit **Out** scope: "cn-sigma tree changes (Sub 3); the compactor mechanism (Sub 4)").
- `git diff --stat origin/main..HEAD` touches only the three named doc files (plus this cell's own `.cdd/unreleased/691/` artifacts) — no code, no `.cn-sigma/` changes, consistent with the issue's explicit **Out** scope and **Mode: Doctrine-only**.

---

## 2. Source-of-truth table

| Claim / surface | Canonical source | Status / role for α |
|---|---|---|
| Ranked model definition (`r0`/`rN`, provenance, promotion≠rank) | Issue #690 body ("The model (KISS / YAGNI)" section) | Authoritative — α rewrites `MEMORY.md` to state this model, not invent a variant |
| Box topology (ref paths, invariants) | Issue #690 body ("Ref topology" section) + operator clarification comment | Authoritative — α cites verbatim ref-path shapes |
| #684/#688 disposition | Issue #690 body ("Substrate reconciliation (#684)" comment) + operator clarification comment | Authoritative — "prior substrate exploration... not the controlling topology... closed as subsumed" |
| Current `MEMORY.md` content | `docs/reference/runtime/MEMORY.md` (pre-edit, this repo) | Superseded by this cycle — α rewrites, does not append |
| Current `AGENT-MEMORY-LOG-STRUCTURED.md` content | `docs/papers/AGENT-MEMORY-LOG-STRUCTURED.md` (pre-edit) | Reconciled, not replaced — α tightens/aligns, preserves the essay's own voice and structure where not contradictory |
| Current `AGENT-ACTIVATION-LOG-v0.md` content | `docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md` (pre-edit) | Narrowed, not deleted — α adds supersession framing, keeps §0/§0.1 descriptive material intact |
| Scope boundary (in/out) | Issue #691 body, `## Scope` section | Binding — Sub 3 (cn-sigma tree) and Sub 4 (compactor mechanism) are explicitly out of scope for this cell |

---

## 3. α prompt

You are α (implementer) for cnos#691, a doctrine-only documentation cell. Read this scaffold's §1 oracle list and §2 source-of-truth table in full before editing anything.

**Task:** rewrite/reconcile three docs so cnos has one canonical memory doctrine (#690's ranked model + box topology) instead of three competing stories:

1. `docs/reference/runtime/MEMORY.md` — rewrite from the "lean triadic model" to #690's ranked (`r0`/`rN`) model. Retire the triadic α/β/γ framing entirely (not just relabel it). State the box topology as the go-forward mechanism. Name #684/#688 as prior exploration, subsumed.
2. `docs/papers/AGENT-MEMORY-LOG-STRUCTURED.md` — reconcile terminology and claims with the rewritten `MEMORY.md`; this essay is already close (it already uses `r0`/`rN` language) — tighten, don't rewrite wholesale. Cross-link to `MEMORY.md`.
3. `docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md` — add explicit supersession framing: this convention is superseded for *agent-memory* purposes by #690's box model; keep the file (§0/§0.1 writer-locality/wake-ownership material documents live production mechanics and stays factually intact), but it must no longer read as the current memory-continuity model going forward.

**Constraints:**
- Doctrine-only. Do not touch any `.go` file, any code, any `.cn-sigma/` path, any CI config, any label/FSM definition.
- Do not invent topology details not present in #690's issue body — cite it as source of truth (§2 table).
- Preserve each doc's own voice/structure where not contradictory; this is reconciliation, not a from-scratch rewrite for #2 and #3. #1 (`MEMORY.md`) needs a fuller rewrite since its current core model is being retired.
- After editing, walk every AC's oracle in §1 and record the result in `.cdd/unreleased/691/self-coherence.md` under `§R0`.
- Do not create `threads/memory/INDEX.md`, vector/graph stores, or any new canonical memory surface — #690's KISS/YAGNI constraint applies here too (implicitly inherited; not itself an AC, but do not contradict it in the doctrine you write).

**When done:** commit + push to `cycle/691`, write `self-coherence.md §R0` with AC-by-AC verification evidence (grep output, diff stat), and signal review-ready.

---

## 4. β prompt

You are β (reviewer) for cnos#691. Independently walk each AC in §1 against the diff on `cycle/691` and α's `self-coherence.md §R0`. Do not trust α's self-report — re-run the greps/checks yourself.

Specifically verify:
- AC1: triadic framing (`α — episodic`, `β — reflective`, `γ — working continuity` headers, "lean triadic" title) is actually gone from `MEMORY.md`, not just renamed.
- AC2: no contradiction between the essay and the rewritten `MEMORY.md`.
- AC3: `AGENT-ACTIVATION-LOG-v0.md` still exists, still describes real §0/§0.1 mechanics accurately, but now carries explicit supersession-for-memory framing.
- AC4: #684/#688 named as prior exploration/subsumed, box topology stated authoritative, ref paths match #690's issue body verbatim.
- AC5: cross-doc consistency — read all three together, confirm no dangling contradiction; confirm `git diff --stat origin/main..HEAD` touches only the three docs + this cell's `.cdd/unreleased/691/` artifacts.

Write `.cdd/unreleased/691/beta-review.md §R0` with a verdict (`converge` or `iterate`) and, if `iterate`, concrete findings keyed to AC numbers.

---

## 5. Scope guardrails

- **In scope:** the three named docs, this cell's own `.cdd/unreleased/691/` artifacts.
- **Out of scope (explicit in issue body):** `.cn-sigma/` tree changes (Sub 3 of #690), the compactor mechanism (Sub 4 of #690), any code change, any FSM/label change.
- If α finds the doctrine reconciliation genuinely requires a code or `.cn-sigma/` change to be coherent, that is a scope-break — STOP and escalate via `status:blocked` rather than expanding scope silently.

## 6. Friction notes

1. This is the first cell dispatched under #690's wave; #690 itself is a large design doc (not a cell) — α should read the operator-clarification comment on #690 (the most recent comment, from κ, 2026-08-02) as the binding execution-order authority: Sub 1 (this cell) runs first, ahead of Subs 2–5.
2. `AGENT-MEMORY-LOG-STRUCTURED.md` already uses `r0`/`rN` rank language close to #690's model — this is not a from-scratch essay rewrite, it is confirmation + tightening. Do not over-edit an already-mostly-aligned essay.
