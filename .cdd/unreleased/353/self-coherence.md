# Self-Coherence — Issue #353

## Gap

**What exists:** The claude-code-dispatch proposal names §5.2 single-session δ-as-γ via Agent tool as a valid dispatch configuration. It documents γ/δ collapse, harness push restrictions (branch sprawl), and grading-floor implications. It does **not** yet name the **shared-working-tree invariant**: parent session and Agent-tool sub-agents share the same filesystem; concurrent writes to the working tree corrupt state.

**What is expected:** A new subsection (§5.2.x — Parent-session quiescence) naming the invariant explicitly. When a sub-agent is running in the background, the parent session must enter quiescent mode: no file edits, no commits, no branch switches. Reads (status, log, file inspection) are fine. The parent's role during sub-agent runs is dispatch coordination and waiting — not concurrent work.

**Where they diverge:** Empirically observed during tsc cycle #36 dispatch. While α R1 was running in the background, the parent session (operator-as-γ) made a manual edit to the working tree. When α attempted to commit, the working-tree state was unexpected; the commit picked up parent's uncommitted changes, leading to a corrupted commit that had to be amended.

## Mode

**docs-only** — No code change, no version bump. Adding one prose subsection to `cdd/operator/SKILL.md` §5.2.

## Active Skills

- Tier 1a: CDD.md, cdd/SKILL.md, gamma/SKILL.md
- Tier 1b: issue/SKILL.md (loaded for AC verification)
- Tier 2: eng/write-functional (for functional/procedural prose)
- Tier 3: (none — pure documentation)

## Acceptance Criteria

### AC1 — §5.2.x subsection authored

**Invariant:** subsection enumerates ≥5 prohibited actions and ≥3 permitted actions during sub-agent runs.

**Oracle:** `rg 'Parent.session quiescence|quiescent' cnos:cdd/operator/SKILL.md` returns ≥1 hit.

**Positive:** a §5.2 operator reading the subsection knows exactly what they can / cannot do mid-dispatch.

**Negative:** soft "operator should avoid concurrent edits" prose — must be a concrete list.

**Surface:** `cnos:cdd/operator/SKILL.md`.

### AC2 — Failure-mode worked example present

**Invariant:** example present; names ≥2 specific git operations (`git add`, `git commit`) in the corruption sequence.

**Oracle:** `rg 'worked example|when this is violated' cnos:cdd/operator/SKILL.md` returns ≥1 hit.

**Positive:** future operators reading the example understand the failure mode without having to experience it.

**Negative:** abstract-only prescription.

**Surface:** same as AC1.

### AC3 — Sub-agent parallelism note

**Invariant:** one paragraph naming the parallel-sub-agent shared-WT case.

**Oracle:** `rg 'parallel sub.agent|concurrent.*sub.agent' cnos:cdd/operator/SKILL.md` returns ≥1 hit.

**Positive:** operators don't dispatch concurrent-write sub-agents.

**Negative:** the parallelism note is absent or omits the WT-sharing detail.

**Surface:** same as AC1.

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | — | — | Issue #353 selected from operator documentation gap |
| 1 Select | — | — | Selected gap: Parent-session quiescence invariant missing from §5.2 |
| 2 Branch | `cycle/353` | cdd | Created from origin/main at dc65c7e5 |
| 3 Bootstrap | N/A | — | docs-only cycle, no version directory needed |
| 4 Gap | self-coherence.md §Gap | cdd | Named incoherence: shared-WT safety invariant not documented |
| 5 Mode | self-coherence.md §Mode | cdd + eng/write-functional | docs-only mode declared |
| 6 Artifacts | cdd/operator/SKILL.md edit | eng/write-functional | Add §5.2.x subsection per ACs |
| 7 Self-coherence | self-coherence.md | cdd | Document implementation state |

## Known Debt

None expected for this documentation addition.

## Implementation Evidence

### AC1 — §5.2.1 subsection authored

- **Status:** ✅ COMPLETE
- **Evidence:** Added §5.2.1 "Parent-session quiescence during sub-agent runs" to `cdd/operator/SKILL.md` lines 301-323
- **Prohibited actions:** 5 enumerated (file edits, commits, branch switches, git add/restore, git pull/fetch HEAD updates)  
- **Permitted actions:** 3 categories enumerated (reads, coordination, external filesystem ops)
- **Oracle verified:** `rg 'Parent.session quiescence|quiescent'` returns 2 hits

### AC2 — Failure-mode worked example present

- **Status:** ✅ COMPLETE
- **Evidence:** "What goes wrong when this is violated" section on line 319
- **Git operations named:** `git add`, `git commit` in corruption sequence
- **Oracle verified:** `rg 'worked example|when this is violated'` returns 1 hit

### AC3 — Sub-agent parallelism note

- **Status:** ✅ COMPLETE  
- **Evidence:** "Sub-agent parallelism note" paragraph on line 321
- **Details:** Names parallel sub-agents share working tree with parent and each other
- **Oracle verified:** `rg 'parallel sub.agent|concurrent.*sub.agent'` returns 1 hit

## Review-Readiness

**Base SHA:** dc65c7e552414991c3b201513015412c5cf1ac42 (origin/main at cycle creation)
**Head SHA:** d11ca84b0dddfb0efc0b17d4a12c0d8af96d9ef3 (current cycle/353)
**Branch:** cycle/353
**CI Status:** Not applicable (docs-only change, no CI dependencies)

All acceptance criteria have been implemented and verified. The subsection has been added to the correct location in `cdd/operator/SKILL.md` within the §5.2 dispatch configuration section.

**Ready for β review.**