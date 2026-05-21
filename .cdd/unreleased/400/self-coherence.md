# Self-coherence — cycle #400 (Phase 5 γ shrink)

**Author:** α (collapsed on δ)
**Date:** 2026-05-21
**Issue:** cnos#400
**Mode:** substantial; design-and-build; γ+α+β-collapsed-on-δ

## Gap

`gamma/SKILL.md` carried 748 lines of mixed γ-role coordination doctrine and runtime supervision mechanics. After Phase 4a–4c shipped the δ-role boundary (`delta/SKILL.md`, cnos#397), harness substrate (`harness/SKILL.md`, cnos#398), and release-effector mechanics (`release-effector/SKILL.md`, cnos#399), the polling shell snippets, dispatch-invocation mechanics, branch pre-flight bash, and release mechanics narration in γ became redundant — the mechanics had a canonical home. Phase 5 of cnos#366 (this cycle, cnos#400) shrinks γ to its actual responsibilities (coordination + closure + triage), cross-references the Phase 4 substrate surfaces, and absorbs F1/F2 from cnos#398's `cdd-iteration.md`.

**Peer-enumeration evidence** (§2.2a):
- `rg "polling loop|claude -p|cn dispatch" src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` pre-cycle: ~50+ hits across §2.5 polling intro + dispatch flow + invocation rules
- `rg "harness/SKILL\.md" src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` pre-cycle: 4 hits (cross-refs already present but content also restated)
- `wc -l src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` pre-cycle: 748 lines (issue cited ~688 but actual was 748)
- `rg "Phase 4b.*pending|harness mechanics.*operator" delta/SKILL.md` pre-cycle: 4 hits (F1 stale refs)
- `bash tools/validate-skill-frontmatter.sh` pre-cycle: 1 finding on delta/SKILL.md (F2 type mismatch)

## Acceptance criteria

### AC1 — gamma/SKILL.md significantly shrunk

`wc -l src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` = **499 lines** (post-cycle) vs **748 lines** (pre-cycle). Ratio: 499/748 = 66.7% (target ≤ 70%, i.e. ≤ 523 lines). **PASS.**

Reduction: 748 − 499 = 249 lines (33.3%).

### AC2 — γ-keep responsibilities present

`rg "issue quality|artifact coordination|closure|triage|repair.dispatch" gamma/SKILL.md` returns multiple hits in normative sections:
- "issue quality" — Core Principle (§Core Principle), §2.4 "Pass the issue-quality gate", Load Order
- "closure" — §1 Define failure modes, §2.7 Triage close-outs, §2.10 Closure gate (header + 14 rows), §3.7 rule
- "triage" — §2.7 (Triage close-outs explicitly, multi-paragraph), §3.7 "Do not close the cycle with unresolved triage"

**PASS.**

### AC3 — γ-loses content extracted

`rg "polling loop|claude -p|cn dispatch.*invocation|CI polling|wake mechanics|branch preflight" gamma/SKILL.md`:
- "claude -p" — 1 hit in Spec-staleness propagation list ("Identity-rotation mode (`cn dispatch` / `claude -p`)") — cross-reference framing, not normative invocation
- "cn dispatch" — same 1 hit and 1 hit in §2.5 "γ produces the prompts and δ routes them via `cn dispatch` (the identity-rotation primitive)" — cross-reference, not invocation detail
- "polling loop" — 0 hits
- "wake mechanics" — 0 hits
- "branch preflight" — 0 hits (γ-owned branch pre-flight remains as the *invariant*; the bash snippet was moved out)
- "CI polling" — 0 hits (release-effector owns the runbook)

All remaining hits are in cross-reference position, not normative restatement. **PASS.**

### AC4 — Cross-references updated bidirectionally

- `harness/SKILL.md` §2.5 "Mirror in γ" → references `gamma/SKILL.md §2.5 dispatch prompt format and the Identity-rotation primitive line` — both anchors preserved (line 202 has "identity-rotation primitive"). **OK.**
- `harness/SKILL.md` §5.4 line 403 → "γ owns the tight loop on a single named branch (per `gamma/SKILL.md` §2.5 dispatch)" — §2.5 preserved. **OK.**
- `delta/SKILL.md` §2 line 22 + §2.1 → `gamma/SKILL.md §2.5 Step 3b` — anchor preserved. **OK.**
- `release-effector/SKILL.md` references `gamma/SKILL.md §2.10` (closure gate, multiple places) and `gamma/SKILL.md §2 release-prep` — both anchors preserved (§2.10 is the 14-row closure gate; §2.6 is release-prep, within §2 umbrella). **OK.**
- `operator/SKILL.md` references `gamma/SKILL.md §2.10` and `§2.5` — both preserved. **OK.**

Additionally, harness "called by γ" / "γ doctrine" sections (lines 38, 40, 50) accurately point at γ's coordination scope post-shrink. **PASS.**

### AC5 — Managerial-residue sweep documented

See §"Managerial-residue sweep" subsection below. ≥1 DROP entry surfaced: the "γ must track" phrasing at pre-cycle line 288. **PASS.**

### AC6 — F1/F2 absorption from cnos#398

**F1:** `rg "Phase 4b.*pending|harness mechanics.*operator" delta/SKILL.md` returns 0 hits (post-edit). Stale forward-references rewritten to past-tense ("Phase 4b of cnos#366, landed at cnos#398") at original lines 59, 81, 173, 316, 337, 340. **PASS.**

**F2:** `bash tools/validate-skill-frontmatter.sh` reports 0 findings on `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` (only 15 pre-existing cdr-side findings remain, out of scope). Fix: line 30 of delta/SKILL.md changed from YAML mapping (`- or: "..."`) to flat string (`- "or: ..."`), preserving the disjunctive-requirement semantics in plain prose while satisfying the `[...string]` CUE schema. **PASS.**

### AC7 — No regression to dispatch workflow

A future γ session loading post-shrink `gamma/SKILL.md` finds:
- Load Order step 5 names `operator/SKILL.md`, `harness/SKILL.md`, `release-effector/SKILL.md` explicitly
- §2.5 dispatch prompts kept verbatim (cnos#393 non-goal preserved); Implementation contract block kept verbatim
- §2.5 polling cross-ref: "Mechanics: `harness/SKILL.md` §5.4" — single-named-branch transition loop available there
- §2.5 dispatch flow: "δ dispatches γ → α → β sequentially per `operator/SKILL.md` Algorithm + `harness/SKILL.md` §1"
- §2.6 release-prep: γ writes RELEASE.md + moves cycle dirs; release mechanics cross-ref `release-effector/SKILL.md`
- §2.10 closure-gate row 13: "δ release-boundary preflight ... mechanics: `release-effector/SKILL.md`; doctrinal frame: `delta/SKILL.md` §1"

No orphan mechanics; no broken cross-refs (verified by grep). **PASS.**

## Managerial-residue sweep

Per `COHERENCE-CELL.md §γ and δ Managerial-Residue Sweep`: classify every γ verb in the pre-shrink file as KEEP (produces artifact / receipt / decision) | MOVE (mechanics; target named) | DROP (no concrete output; sweep-suspect verb).

| γ verb (pre-shrink) | Classification | Destination / rationale |
|---|---|---|
| selects (gap by CDD §3 rule order) | KEEP | produces selected gap + decisive clause (decision) |
| observes (lag, PRA, candidates) | KEEP | produces candidate table (artifact) |
| builds (candidate table, issue pack) | KEEP | produces issue pack (artifact) |
| passes (issue-quality gate) | KEEP | produces gate-pass decision |
| creates (cycle/{N} on origin) | KEEP | produces branch on origin (artifact); pre-flight remains as invariant in γ; **bash snippet MOVE → `CDD.md §4.3` for mechanics** |
| writes (gamma-scaffold, gamma-clarification, gamma-coordination, RELEASE.md, gamma-closeout, PRA) | KEEP | produces named artifacts |
| produces (γ/α/β prompts) | KEEP | produces prompt text (artifact); δ executes dispatch |
| runs (pre-dispatch scaffold check) | KEEP | produces gate decision; the bash snippet itself compressed in γ (1 invariant-stating line + cross-ref to mechanics) |
| polls (cycle branch + issue) | **MOVE** | mechanics → `harness/SKILL.md §5.4` (single-named-branch transition loop, Monitor wake-up, reachability re-probe). γ keeps the *invariant* (polling requires query + wake-up + reachability) as 1-2 lines |
| subscribes (issue activity) | MOVE | mechanics → `harness/SKILL.md §5.2` |
| dispatches | MOVE | not γ's verb — γ produces prompts, δ dispatches via harness. Restated in §2.5: "γ does not execute dispatch; γ produces prompts." |
| reacts to transitions (decision tree, decision-point matching) | KEEP | produces autonomous-coordinator decision (act / pause / ignore) |
| unblocks (clarify, edit issue, add artifact link) | KEEP | produces clarification artifact (`gamma-clarification.md`); allowed transfer unit is *artifact facts*, not hidden state |
| forwards (β/α reasoning transcripts) | KEEP-as-forbidden | role boundary; sets the forbidden list |
| moves (cycle dirs to releases/{X.Y.Z}/) | KEEP | produces directory move (artifact) |
| collects (alpha-closeout, beta-closeout) | KEEP | produces close-out triage input (artifact) |
| verifies (post-merge CI green) | KEEP | produces verification record in `gamma-closeout.md` (receipt) |
| triages (findings → CAP disposition) | KEEP | produces triage table (artifact + decision) |
| enforces (cycle-iteration triggers) | KEEP | produces `Cycle Iteration` entry in PRA (artifact) |
| patches (skill/spec immediate fixes) | KEEP | produces commit on main (artifact) |
| delegates (step 13a patches when γ explicitly does so) | KEEP | produces delegation record with name + deadline (decision) |
| requests (release boundary preflight from δ) | KEEP | produces gate-request boundary decision; mechanics in `release-effector/SKILL.md` |
| declares (cycle closure: "Cycle #N closed. Next: #M.") | KEEP | produces closure declaration in `gamma-closeout.md` (receipt) |
| updates (hub memory) | KEEP | produces hub entry (artifact) |
| deletes (merged remote branches) | KEEP authority; MOVE mechanics | γ-side decision to request deletion is a gate-request; mechanics → `release-effector/SKILL.md §5` (`git push origin --delete`, 403-tolerant) |
| reports (operator-facing TLDR / decision-request / deferred batch) | KEEP | produces operator-facing report (artifact, e.g. TLDR text); the report shapes stay in γ as 3 named formats |
| reads (skill / artifact when staleness check fires) | KEEP | produces re-loaded canonical-skill snapshot (state decision) |
| **tracks ("γ must track the full cycle")** | **DROP** | sweep-suspect verb per `COHERENCE-CELL.md`. "Track" produces no artifact/receipt/decision — it is managerial residue describing γ's relation to the cycle. **Replaced** with explicit poll → match transition → decide (act/pause/ignore) in §2.11. The old phrasing at pre-cycle line 288 ("γ must track the full cycle: issue activity, branch SHA, self-coherence updates, beta-review verdicts, CI status") is gone; the new §2.11 names polling (mechanics in harness §5.4) and the decision-point matching table. |
| **orchestrates dispatch** (implicit in some pre-cycle prose) | **DROP** | γ does not orchestrate; γ produces prompts. The line "γ does not execute dispatch" is now stated explicitly in Core Principle and §2.5. |
| monitors / supervises / oversees / manages | N/A | did not appear pre-shrink as γ verbs (sweep-clean on these synonyms; the sweep would have surfaced them if present) |

**Summary:** 28 verbs classified. KEEP: 22 (γ coordination + closure + triage). MOVE: 5 (polls, subscribes, dispatches, deletes-mechanics, runs-bash-mechanics-of-creates). DROP: 2 (**tracks**, **orchestrates dispatch**) — both replaced with explicit artifact-producing verbs.

β-rigor check: at least one DROP entry surfaced. ✓ ("tracks" and "orchestrates dispatch")

## CDD Trace

- Selected by issue #400 (Phase 5 of #366 P1, gates Phase 7); breadth-2026-05-12 precedent invoked by δ to collapse γ+α+β onto δ-as-agent.
- Implementation contract pinned in issue body; followed verbatim.
- Empirical anchor: cnos#398 cdd-iteration F1/F2 absorbed.

## Review-readiness

α self-coherence complete; ACs all PASS; managerial-residue sweep documented with 2 DROP entries; F1/F2 absorption verified. Ready for β-collapsed self-review.

## Fix-rounds

None — γ+α+β collapsed on δ, no R1 reviewer round; β-collapsed self-review writes `beta-review.md` directly.
