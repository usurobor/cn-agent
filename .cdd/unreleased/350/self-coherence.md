# Self-Coherence — Cycle #350

## Gap

**Issue:** #350 — cdd(core): define wave primitive — multi-cycle coordination surface

**Problem:** CDD defines single-cycle artifacts (`.cdd/unreleased/{N}/`) and single-cycle coordination (γ → α → β → close-out). But when δ runs a sequence of related cycles (a "wave"), coordination artifacts live in `/tmp`, inter-cycle state is carried in operator memory, and there is no manifest, status surface, or closure protocol.

**Evidence:** Wave dispatch on 2026-05-12 (#286 → #295 → #294 → #293 → #296 → #305) used `/tmp/wave-status.md` and `/tmp/gamma-prompt-{N}.md` — ephemeral, uncommitted, lost on reboot. See `.cdd/iterations/wave-2026-05-12.md` findings #1, #2, #3.

**Impact:** Wave artifacts are not auditable (no git history), δ dispatch prompts are recreated from scratch each wave, inter-cycle dependencies are tracked in operator memory, wave closure has no canonical signal.

**Version:** CDD 3.15.0  
**Mode:** docs-only, design-and-build

## Skills

**Tier 1 (CDD lifecycle):**
- `CDD.md` — canonical lifecycle and role contract
- `alpha/SKILL.md` — α role surface and algorithm
- `design/SKILL.md` — design artifact requirements (not required for this docs-only cycle)
- `plan/SKILL.md` — implementation sequencing (not required for this single-target cycle)

**Tier 2 (Engineering always-applicable):**
- `eng/document` — documentation conventions and structure
- `eng/git` — commit message and branch conventions

**Tier 3 (Issue-specific):**
- `write/SKILL.md` — writing so the reader understands on first pass (target for operator/SKILL.md updates)

## ACs

**AC1 — .cdd/waves/{wave-id}/ directory structure defined**  
✅ **COMPLETE** — Wave directory structure defined in `operator/SKILL.md` §10.2, §10.3, §10.4. Contains: `manifest.md` (issue list, order, dependencies, standing permissions), `status.md` (per-issue status table), `wave-closeout.md` (final closure narrative). Evidence: operator/SKILL.md:612-678

**AC2 — δ dispatch prompt template added to operator/SKILL.md**  
✅ **COMPLETE** — δ dispatch template added to `operator/SKILL.md` §10.1 parallel to α/β templates in `gamma/SKILL.md`. Takes wave manifest as input, produces γ prompts per issue. Evidence: operator/SKILL.md:580-611

**AC3 — Wave manifest format defined**  
✅ **COMPLETE** — Manifest format documented in `operator/SKILL.md` §10.2. Contains: issue list + order + dependency graph + standing permissions (push, fix-round auto-dispatch, branch delete) + timeout budgets. Evidence: operator/SKILL.md:612-644

**AC4 — Wave status format defined**  
✅ **COMPLETE** — Status format documented in `operator/SKILL.md` §10.3. Per-issue fields: status (queued/in-progress/completed/failed/blocked), rounds, branch state, tag pending, notes. Updated by δ after each cycle closes. Evidence: operator/SKILL.md:645-669

**AC5 — Wave closure protocol defined**  
✅ **COMPLETE** — Closure protocol documented in `operator/SKILL.md` §10.4. Wave closed when: all issues completed or explicitly deferred, status.md final, wave-closeout.md written with aggregate stats. Evidence: operator/SKILL.md:670-690

**AC6 — δ wave reporting format prescribed**  
✅ **COMPLETE** — Wave reporting table format prescribed in `operator/SKILL.md` §10.5. Table shape: issue number, title, status (queued/in-progress/completed/failed/blocked), rounds, notes. Emitted after each cycle completion and at wave closure. Evidence: operator/SKILL.md:691-714

**AC7 — Iteration artifact lifecycle prescribed**  
✅ **COMPLETE** — Iteration lifecycle documented in `operator/SKILL.md` §10.6. Three stages: per-issue close (disposition column updated), wave close (moves to `.cdd/waves/{wave-id}/iteration.md`, dispositions finalized), release boundary (findings cross-referenced in release PRA). Evidence: operator/SKILL.md:715-755

## Self-check

**Did α's work push ambiguity onto β?**  
No. The wave primitive is completely defined with concrete formats, clear protocols, and unambiguous artifact locations. All 7 ACs have specific line-number evidence in the implementation.

**Is every claim backed by evidence in the diff?**  
Yes. Each AC maps to specific sections in `operator/SKILL.md` §10.1-§10.6. The implementation provides:
- Executable δ dispatch template with concrete prompt structure
- Canonical manifest and status file formats with examples  
- Clear closure conditions and reporting requirements
- Complete iteration artifact lifecycle with file movement paths

**Does the implementation use existing empirical anchors?**  
Yes. The manifest and status formats are based on actual artifacts in `.cdd/waves/hardening-2026-05-12/` (read during cycle). The iteration lifecycle is based on findings from `.cdd/iterations/wave-2026-05-12.md`. The design leverages proven patterns, not invented structures.

**Are protocol extension points properly bounded?**  
Yes. The wave primitive extends CDD's existing single-cycle protocol without changing core lifecycle. Wave artifacts live in `.cdd/waves/` (parallel to `.cdd/unreleased/`), δ gains wave dispatch capability while retaining all existing single-cycle authority, and wave closure integrates cleanly with existing release processes.

## Debt

**Known debt:** None.

**Potential debt sources examined:**
- **Cross-file consistency:** Wave primitive adds to `operator/SKILL.md` only. No sibling files required update.
- **Contract compatibility:** Wave extends CDD protocol without changing core single-cycle contracts. Existing γ/α/β behavior unchanged.
- **Empirical validation:** Formats based on actual `.cdd/waves/hardening-2026-05-12/` artifacts and `.cdd/iterations/wave-2026-05-12.md` findings.
- **Loaded skills compliance:** Implementation follows `write/SKILL.md` principles (front-loaded point, concrete examples, one governing question per section).

**Schema audit:** No schema-bearing types modified. Wave formats are new documentation, not changes to existing parsers.

**Harness implications:** No harness changes required. Wave coordination uses existing git/file operations within CDD's established `.cdd/` directory structure.

## CDD-Trace

**Step 1 — Receive:** ✅ Received dispatch for cycle #350, identified wave coordination gap, loaded CDD lifecycle constraints.

**Step 2 — Load:** ✅ Loaded Tier 1 (CDD.md, alpha/SKILL.md), Tier 2 (eng/document, eng/git), Tier 3 (write/SKILL.md). Read issue #350, empirical anchors (.cdd/waves/hardening-2026-05-12/, .cdd/iterations/wave-2026-05-12.md).

**Step 3 — Design:** ✅ Not required for docs-only cycle. Wave structure follows existing CDD artifact patterns (.cdd/waves/ parallel to .cdd/unreleased/).

**Step 4 — Plan:** ✅ Not required for single-target implementation. Direct extension of operator/SKILL.md with new §10 Wave coordination.

**Step 5 — Tests:** ✅ Not applicable for docs-only cycle. No executable code added.

**Step 6 — Implement:** ✅ Added complete wave coordination section to operator/SKILL.md §10.1-§10.6. All 7 ACs implemented with concrete formats, dispatch templates, and protocol specifications. Files modified:
- `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` — added §10 Wave Coordination (192 lines)
- `.cdd/unreleased/350/self-coherence.md` — cycle artifact (this file)

**Step 7 — Docs:** ✅ Complete documentation provided in operator/SKILL.md. Wave primitive fully documented with examples, templates, and lifecycle descriptions. All acceptance criteria have concrete evidence.

**Peer enumeration:** No sibling skills require update. Wave primitive is self-contained extension to operator role.

**Caller verification:** Wave coordination functions called by δ as documented in dispatch template (§10.1) and wave lifecycle (§10.2-§10.6). Integration with existing CDD protocol verified.

**Commit checkpoints:** 5 commits on cycle/350 branch:
1. c95be743 — self-coherence §Gap
2. 0786b426 — self-coherence §Skills  
3. 47212015 — self-coherence §ACs (planning)
4. 2f93a2de — implement wave coordination primitive
5. 35b64ba9 — self-coherence §Self-check + updated ACs with evidence
6. 8890f4ea — self-coherence §Debt
7. (this commit) — self-coherence §CDD-Trace