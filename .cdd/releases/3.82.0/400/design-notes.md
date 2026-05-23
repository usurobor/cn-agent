# Design notes — cycle #400 (Phase 5 γ shrink)

**Author:** α (collapsed on δ)
**Date:** 2026-05-21

## 1. Section-by-section inventory of current gamma/SKILL.md (748 lines)

Each section classified KEEP (γ coordination/closure/triage; produces artifact / receipt / decision) | MOVE (mechanics; destination named) | DROP (no concrete output; managerial residue) per `COHERENCE-CELL.md §γ and δ Managerial-Residue Sweep`.

| Section (lines) | Topic | Classification | Destination / rationale |
|---|---|---|---|
| frontmatter (1-38) | YAML | KEEP | role-skill identity |
| §Core Principle (40-57) | γ holds cycle coherence | KEEP, tighten | doctrine kernel |
| §Load Order (59-78) | load CDD.md, gamma, issue, post-release, operator; staleness check | KEEP (load order), MOVE part | operator load-row mentions "Dispatch mechanics... live in `harness/SKILL.md`" — already cross-references; trim wordy operator-load row |
| §Step map (80-90) | map of CDD.md §1.4 steps to file sections | KEEP, compress | structural index |
| §1.1–§1.3 Define (94-128) | parts; how they fit; failure mode | KEEP, compress | doctrine |
| §2.1 Step 1a Observe (133-153) | observation inputs + candidate table + cross-repo intake | KEEP | produces candidate table (artifact) |
| §2.2 Step 1b Select (155-172) | apply CDD.md §3; sizing | KEEP | produces selection record (decision) |
| §2.2a Peer enumeration (174-190) | grep before scaffolding §Gap | KEEP | produces honest-claim evidence (artifact) |
| §2.3 Step 2 Build issue pack (191-205) | issue pack requirements | KEEP | produces issue pack (artifact) |
| §2.4 Step 2 Issue quality gate (207-223) | pre-dispatch gate | KEEP | produces gate decision |
| §2.5 Step 3a Create branch (227-258) | branch creation snippet + pre-flight | KEEP behavior, MOVE shell snippet | the bash snippet (lines 241-256) is mechanics; CDD.md §4.3 owns pre-flight rule; harness already documents branch mechanics. Keep the *decision* (γ creates `cycle/{N}` before dispatch; γ-owned pre-flight); trim the 16-line bash snippet to a 4-line summary + cross-ref to CDD.md §4.3 |
| §2.5 Step 3b Pre-dispatch scaffold check (262-285) | scaffold-presence gate + dual of review/SKILL.md 3.11b | KEEP (the gate; the empirical anchor); MOVE the bash snippet | the binding rule is γ doctrine (gate). The 10-line bash snippet (lines 274-282) is mechanics — but the snippet *is* the gate check; cannot be moved without losing the gate semantics. Compress to ~5 lines |
| §2.5 Step 3b Polling intro (287-323) | "Immediately begin polling" + Monitor wake-up + transition loop bash | MOVE | harness/SKILL.md §5.4 already has the identical "Single-named-branch polling (γ-side, under Monitor)" transition loop (lines 401-431). DELETE from γ; replace with 1-line "γ polls per `harness/SKILL.md` §5.4" + the requirement "polling requires query + wake-up + reachability probe" stays as γ-side rule (1 sentence) |
| §2.5 Step 3b Dispatch prompts γ/α/β (324-413) | prompt templates + Implementation contract (cnos#393) + mesh + empirical anchor + dispatch flow + identity-rotation | KEEP (cnos#393 contract is doctrinally pinned; non-goal forbids changing); MOVE the dispatch shell narration | KEEP the prompt templates, KEEP the `## Implementation contract` section verbatim (cnos#393 non-goal), KEEP the four-surface mesh narration, KEEP the empirical anchor. The "Dispatch flow" pseudo-shell (lines 396-401) is mechanics, but it's 6 lines and serves as orientation; condense to a single inline "δ dispatches γ → α → β sequentially per `operator/SKILL.md` Algorithm + `harness/SKILL.md` §1" |
| §2.5 Step 3b Spec-staleness propagation (415-424) | when to write `gamma-coordination.md` after spec change | KEEP | produces coordination note artifact |
| §2.5 Step 5 Unblock (426-447) | allowed/forbidden unblock actions + issue-edit cache-bust | KEEP | produces clarification artifact (gamma-clarification.md) |
| §2.6 Steps 6–7 Release prep (449-464) | γ writes RELEASE.md + moves cycle dirs | KEEP | produces RELEASE.md + cycle-dir move (artifacts) |
| §2.7 Steps 8–9 Triage close-outs (466-507) | collect close-outs, cross-repo close-out, post-merge CI verification, triage table, CAP disposition | KEEP | produces triage record (artifact) |
| §2.8 Steps 10–11 Cycle iteration triggers (509-525) | trigger table; closure rules | KEEP | produces Cycle Iteration entry (artifact) |
| §2.9 Step 13 Process-gap check (527-543) | independent γ check beyond §9.1 | KEEP | produces patch / next-MCA decision |
| §2.10 Steps 13–15 Closure gate (545-567) | 14-row closure gate list + gamma-closeout.md | KEEP | produces gamma-closeout.md (artifact + decision) |
| §2.11 γ as autonomous coordinator (569-612) | decision tree, decision-point matching, operator-report formats (TLDR, decision-request, deferred-question batch), kata 3-round | KEEP doctrine; MOVE operator-report-format details | The decision tree and decision-point matching are γ coordination doctrine — produces "act / pause / ignore" boundary decisions. KEEP. The "Operator-facing report formats" (TLDR template, decision-request template, deferred-question batch) are mechanics for *how* to surface to operator. MOVE to `delta/SKILL.md` (outward boundary — operator-facing formats are δ's outward surface) OR KEEP compressed as 1-line examples |
| §3.1–§3.8 Rules (616-657) | 8 short rules including 3.8 Silence rule | KEEP, compress | doctrine — each rule is a γ-decision constraint |
| §4 Kata A Selection (663-680) | embedded kata | KEEP (compress) | embedded-kata surface per skill schema |
| §4 Kata B Issue quality (682-710) | embedded kata | KEEP (compress) | embedded-kata surface |
| §Resumption (712-727) | resumption protocol for γ artifacts | KEEP (compress) | produces resumed artifact |
| §4 Kata C Cycle iteration (729-748) | embedded kata | KEEP (compress) | embedded-kata surface |

### Verb-level sweep (per AC5 — `self-coherence.md §Managerial-residue sweep` will carry the full table)

Verbs that appear with γ as the subject:
- **selects** → KEEP (produces selected gap; decision)
- **observes** → KEEP (produces candidate table; artifact)
- **builds** (candidate table, issue pack) → KEEP (artifact)
- **passes** (issue-quality gate) → KEEP (decision)
- **creates** (cycle branch) → KEEP (artifact: branch on origin)
- **writes** (scaffold, clarification, RELEASE.md, gamma-closeout.md, gamma-coordination.md) → KEEP (artifact)
- **polls** → **MOVE** (mechanics → harness/SKILL.md §5.4)
- **dispatches** → MOVE (mechanics → harness/SKILL.md §1; doctrine "γ produces prompts, δ executes dispatch" stays)
- **subscribes** → MOVE (mechanics → harness/SKILL.md §5)
- **produces** (prompts) → KEEP (artifact: prompt text)
- **runs pre-flight** → MOVE (shell snippet → CDD.md §4.3 + harness)
- **unblocks** → KEEP (produces clarification artifact)
- **edits issue** → KEEP (produces issue edit + cache-bust clarification entry)
- **moves cycle dirs** → KEEP (produces directory move on main — artifact)
- **collects close-outs** → KEEP (produces close-out triage input — artifact)
- **verifies CI green post-merge** → KEEP (produces verification record in gamma-closeout.md — receipt)
- **triages** → KEEP (produces triage table — artifact + decision)
- **enforces cycle-iteration outputs** → KEEP (produces Cycle Iteration entry — artifact)
- **patches skill/spec** → KEEP (produces commit on main — artifact)
- **requests release** → KEEP (produces gate-request boundary decision)
- **declares closure** → KEEP (produces closure declaration — receipt: gamma-closeout.md)
- **updates hub memory** → KEEP (produces hub entry — artifact)
- **deletes merged remote branches** → MOVE (mechanics → release-effector/SKILL.md §5; γ-side "requests deletion" stays as a gate-request decision)
- **reports** (operator-facing TLDR / decision-request) → KEEP, but the *report-format templates* (TLDR mech, deferred-question batch) → MOVE-or-COMPRESS to delta/SKILL.md outward boundary OR keep as 1-line examples since they produce concrete operator-facing artifacts (compressed reports). Decision: COMPRESS in γ (1 example each) + note delta/SKILL.md owns the outward boundary
- **decides** (act / pause / ignore on polling transition) → KEEP (produces autonomous-decision — boundary decision at γ level)
- **acts autonomously** → KEEP (γ coordination function)

**DROP candidates (verbs without concrete output):**
- **"manages"** — does not appear in γ as verb (good).
- **"supervises"** — does not appear (good).
- **"monitors"** — does not appear directly; "γ must track the full cycle" (line 288) is residue: *tracking* is not an artifact-producing verb. The act γ produces is *poll → match transition → decide act/pause/ignore*. The "track" framing is managerial residue. **DROP the "γ must track" phrasing**; replace with "γ polls (mechanics in harness §5.4); on transition, γ matches the decision-tree (§2.11)."
- **"oversees"** — does not appear (good).
- **"holds the full cycle"** (line 288 "γ must track the full cycle: issue activity, ... β's verdicts, and branch CI status") — KEEP narrowed: γ polls these surfaces; the *decision* γ produces on each transition is the artifact. The "must track" sentence is residue.
- **"observes"** of release CI by γ — line 480 "γ must verify CI ran green" → KEEP (γ records the run URL in gamma-closeout.md — produces verification artifact).
- **"orchestrates dispatch"** (implied in §2.5) — γ does NOT orchestrate; γ produces prompts. The "γ does not execute dispatch" line (324) is correct doctrine. KEEP that line.

**Net DROP** entries (≥1 required by AC5): the **"γ must track the full cycle"** phrasing at line 288 and the verb **"track"** as a coordination verb. Replace with explicit poll-then-decide. (β-rigor: ≥1 DROP entry honored.)

## 2. Cross-reference impact analysis

After the γ shrink, the following cross-references in destination skills need verification:

### 2.1. harness/SKILL.md
- §5.4 "Single-named-branch polling (γ-side, under `Monitor`)" — γ's transition loop lives here already. **No edit needed.** γ's new prose will cross-ref this section.
- §2.5 "Mirror in γ" — references `gamma/SKILL.md` §2.5 dispatch prompt format. **Verify the §2.5 anchor remains valid post-shrink.** Plan: keep `## 2.5. Steps 3a–5 ...` header so the anchor resolves.
- §3 Git identity / §4 Parallel dispatch / §6 Timeout recovery — no γ-content overlap; no edit.

### 2.2. delta/SKILL.md
- §1.1 cross-refs `gamma/SKILL.md` §2.10 (closure gate) and `release-effector/SKILL.md` (release mechanics). **No γ-content drift; no edit needed from γ side.**
- §2.1 cross-refs `gamma/SKILL.md` §2.5 Step 3b (implementation-contract template). **Keep §2.5 anchor; non-goal preserves the template.**
- F1: lines 59, 61, 81, 173, 316, 336, 340 have stale "Phase 4b pending / harness mechanics currently in operator" references. **REFRESH to past-tense ("Phase 4b shipped at cnos#398").**
- F2: frontmatter `requires.1` is a YAML mapping; needs to be a flat string per `schemas/skill.cue`. **FIX line 30.**

### 2.3. release-effector/SKILL.md
- Multiple cross-refs to `gamma/SKILL.md` §2.10 closure gate and §2.6 release-prep. **Keep both anchors.**
- §9 cross-references row "gamma/SKILL.md §2.10" and "gamma/SKILL.md §2 release-prep". **Verify anchors hold post-shrink.**

### 2.4. operator/SKILL.md
- §1.1 routes through γ. No γ-content overlap requiring edit.
- §3 cross-refs `gamma/SKILL.md` §2.10. **Keep §2.10 anchor.**
- §7 Cycle lifecycle table references `gamma-closeout.md` (γ artifact). No γ-content drift.

### 2.5. CDD.md
- Out of scope (Phase 7). No edit.

## 3. F1/F2 absorption plan

### F1 (delta/SKILL.md stale Phase 4b forward-refs)

Lines requiring edit per cnos#398 cdd-iteration:
- Line 59: "release-effector/SKILL.md..." — current text says "Phase 4b will further separate harness mechanics from operator-as-coordinator." → past-tense ("Phase 4b shipped at cnos#398").
- Line 61: "The dispatch-coordination layer ... is in `operator/SKILL.md`." Now harness/SKILL.md holds dispatch mechanics. → refresh to point at harness.
- Line 81: "Phase 4b will further separate harness mechanics from operator-as-coordinator." → past-tense.
- Line 173: "Phase 4b will… harness mechanics (Phase 4b) remain in `operator/SKILL.md` until that cycle relocates them." → past-tense, point at harness.
- Line 316: "Phase 4b will further separate harness mechanics into its own substrate surface; until then, harness mechanics stay in `operator/SKILL.md`." → past-tense.
- Line 336: "Phase 4b — harness substrate (pending). The harness mechanics currently living in `operator/SKILL.md` ... are slated to relocate into a harness substrate surface in a subsequent cycle." → past-tense (landed via cnos#398; harness/SKILL.md).
- Line 340: "Until Phase 4b ships, harness mechanics remain in `operator/SKILL.md`." → delete or past-tense.

### F2 (delta/SKILL.md frontmatter requires.1 type mismatch)

Line 30 currently reads:
```yaml
  - or: "active γ dispatch prompt awaiting δ's pre-routing enrichment"
```
This is a YAML mapping `{or: "..."}`. `schemas/skill.cue` requires `[...string]`. Fix:
```yaml
  - "or: active γ dispatch prompt awaiting δ's pre-routing enrichment"
```
The `or:` prefix is preserved as plain text inside the string, communicating the intended disjunctive meaning to a human reader without breaking the CUE type. This is the minimum-invasive fix; a broader frontmatter convention for disjunctive `requires` is out of scope (would be a CUE schema change).

## 4. Plan for restructured gamma/SKILL.md

Target structure (estimated line ranges):

1. Frontmatter (38 lines)
2. # Gamma + Core Principle (~15 lines) — kept tight
3. ## Load Order (~15 lines) — compressed; operator load-row trimmed
4. ## Step map (~10 lines) — kept; structural index
5. ## 1. Define — Parts / Fit / Failure (~25 lines compressed)
6. ## 2. Unfold
   - §2.1 Observe + candidate table (~20 lines, compressed cross-repo paragraph)
   - §2.2 Select + sizing (~15 lines)
   - §2.2a Peer enumeration (~15 lines, compressed empirical anchor)
   - §2.3 Issue pack (~10 lines)
   - §2.4 Issue-quality gate (~12 lines)
   - §2.5 Steps 3a–5
     - Step 3a Create branch — decision only, cross-ref to CDD.md §4.3 (~10 lines)
     - Step 3b Pre-dispatch scaffold check (binding gate) (~15 lines, compressed empirical anchor)
     - Step 3b Polling cross-ref to harness §5.4 (~5 lines)
     - Step 3b Dispatch prompts (γ/α/β) + Implementation contract block — KEPT IN FULL per cnos#393 non-goal (~60 lines)
     - Step 3b Spec-staleness propagation (~10 lines)
     - Step 5 Unblock + issue-edit cache-bust (~15 lines)
   - §2.6 Release prep (~10 lines)
   - §2.7 Close-out triage (~25 lines compressed)
   - §2.8 Cycle-iteration triggers (~15 lines)
   - §2.9 Process-gap check (~10 lines)
   - §2.10 Closure gate (~25 lines — keep all 14 rows)
   - §2.11 γ as autonomous coordinator — decision tree + matching + 1-line report-format examples (~25 lines, compressed from 44)
7. ## 3. Rules (~25 lines — 8 rules compressed)
8. ## 4. Embedded Kata A/B/C (~30 lines compressed) + Resumption (~10 lines)

Estimated total: ~450 lines (well under the 523-line cap from AC1).

## 5. Mid-cycle escalations / open items

None at design time. The β-rigor check (≥1 DROP entry) is satisfied by the "γ must track the full cycle" verb-DROP entry.
