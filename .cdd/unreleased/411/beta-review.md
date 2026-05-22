# β review | cnos#411 — Sub 6 of #403

**Cycle:** cnos#411
**Branch:** `cycle/411`
**Base SHA:** `71b25672` (Merge cycle/410)
**Head SHA:** at α review-ready signal
**Round:** R1
**Reviewer:** β (collapsed wave-mode dispatch per #411 cycle frame)

## TERMS

The contract β reviews against:

- **#411 ACs (AC1–AC9).** Mechanical-leaning ACs covering CDD.md
  marker removal, CDS pointer insertion, kernel preservation, cross-ref
  re-pointing, CDS path-count threshold, cdr/cds non-touch invariants,
  anchor resolution spot-check, markdown validity.
- **#411 hard rules.** cnos.cdr untouched; CCNF kernel preserved; no
  content moves (Sub 6 is cleanup + cross-reference repair only);
  cnos.cds substantively untouched; no `.cdd/` → `.cds/` filesystem
  migration; CDS § anchor granularity preserved.
- **Implementation contract.** Markdown-only edits; no CLI; package
  scoping per #411; no runtime dependencies; no wire/JSON contract.

## POINTER

α's deliverables:

1. **CDD.md §"Software-specific realization" replacement (D1).** Heading
   renamed `→ cnos.cds` (no longer "pending"); body is a one-paragraph
   pointer to CDS.md plus a 15-bullet per-section pointer list; closing
   "Anchor convention" paragraph rewritten to declare the pre-#402
   anchor forms retired and to name the v0.1 overlay path for surfaces
   where operational expansion still lives in cdd skill files.
2. **CDD.md kernel-section minor edits (D1.1).** §"Domain packages",
   §"Pointers", §"Hard rule" final paragraph, §"Non-goals" final bullet,
   and the preamble all updated to reflect cnos.cds v0.1 having shipped
   (rather than "pending bootstrap"). AC3 permits "minor edits for
   cross-reference consistency"; these qualify.
3. **Cross-reference re-pointing (D2).** 15 files edited; α's table at
   `self-coherence.md` lines 87–116 enumerates per-pattern mappings.
4. **Extraction-map status note (D3).** Single top-of-file Status
   blockquote in `extraction-map.md`. No row-level Status edits — AC7
   compliance preserved.

## β verification (mechanical)

| AC | Check | β observation |
|---|---|---|
| AC1 | grep "pending cds extraction" CDD.md = 0 | confirmed (0) |
| AC2 | CDD.md has CDS pointer paragraph + per-section list | confirmed — body has 15 sub-anchored CDS bullets covering every former-pending marker family |
| AC3 | kernel sections preserved | confirmed — kernel doctrine in §Kernel through §Scope-lift unchanged; §Domain packages / §Pointers / §Hard rule / §Non-goals edits are factual updates (cnos.cds shipped vs pending) and explicit reference repairs, no doctrinal substance changed; preamble similarly |
| AC4 | grep "CDD.md §X" minus kernel exclusions = 0 | confirmed (0 lines) |
| AC5 | grep "cnos.cds/skills/cds/CDS.md" cdd skill files ≥ 5 | confirmed (124 hits) |
| AC6 | git diff origin/main..HEAD -- cnos.cdr/ = 0 | confirmed (0 lines) |
| AC7 | git diff origin/main..HEAD -- cnos.cds/skills/cds/CDS.md = 0; extraction-map.md only top-of-file Status note | confirmed (CDS.md untouched; extraction-map.md diff is the 5-line top blockquote only) |
| AC8 | ≥ 5 re-pointed citations resolve at cited CDS path | confirmed — 11 anchors spot-checked, all 11 resolve (§Selection function, §Branch rule, §Location matrix, §Frozen snapshot rule, §Cycle iteration triggers, §Coordination surfaces, §Polling primitives, §Field 6: Actor collapse rule, §Step table, §Ownership matrix, §Resumption protocol) |
| AC9 | CDD.md valid markdown; section ordering preserved | confirmed — heading hierarchy unchanged; no broken fenced code blocks; cross-link syntax well-formed |

## β substantive read

The diff is **mechanically clean and protocol-coherent**:

1. **Granularity preservation per gamma-scaffold pinning.** α did not
   collapse `§5.3a` and `§5.3b` into a generic CDS.md pointer; they map
   to distinct sub-anchors (§"Location matrix" and §"Frozen snapshot rule"
   / §"Ownership matrix"). Same for §Tracking variants → §"Coordination
   surfaces" / §"Polling primitives" / §"Cycle-state evidence".

2. **Local fallbacks correctly identified.** Two categories of citations
   could not be cleanly re-pointed at CDS canonical sub-anchors and α
   handled them well:
   - **role-identity-is-git-observable** (cited from
     `harness/SKILL.md`, `activation/SKILL.md`, `alpha/SKILL.md` row 14):
     does not exist in CDS.md (it's software-cycle-specific local
     property). α correctly re-points at `operator/SKILL.md` §"Git
     identity for role actors" (verified to exist at operator/SKILL.md
     line 57), preserving auditability.
   - **dispatch prompt formats** (`§1.6a` re-dispatch prompts): the
     prompt text lives in `operator/SKILL.md` §5.2 as v0.1 overlay; CDS
     names only the mechanism (§"Coordination surfaces"). α correctly
     cites both: "(`cnos.cds/skills/cds/CDS.md` §"Coordination surfaces";
     prompt format in `operator/SKILL.md` §5.2 v0.1 overlay)" — this
     preserves granularity without inventing a CDS anchor that doesn't
     exist.

3. **Wave-mode sub-discipline.** β has not edited cnos.cds/skills/cds/
   substantively (Subs 1–5 had that exclusive right per the wave plan);
   the only cnos.cds touch is the extraction-map.md top-of-file Status
   note, which #411 explicitly permits as optional.

4. **Wave-coherence forward.** The retired pre-#402 anchor forms are now
   structurally unresolvable in CDD.md. Subs 1–5 already migrated the
   content; Sub 6 closes the broken-anchor risk; Sub 7 (parallel) is
   independent of this work. #403 closes when Sub 6 + Sub 7 both land.

## Findings

**None.**

β scrutinised three potential failure modes:

- **(a) Cross-repo `STATUS state machine` references.** The
  cross-repo/SKILL.md formerly cited
  "`CDD.md §"Cross-repo proposal lifecycle"` (lines 234–243)". The
  pre-#402 §Cross-repo proposal lifecycle section is gone; α correctly
  re-points at CDS.md §"Coordination surfaces" → §"Cross-repo
  proposals". β verified CDS.md line 1428 hosts the canonical
  STATUS-state-machine doctrine and line 1443 has the `submitted` event
  definition cited by cross-repo/SKILL.md §2.3.3. No broken anchor.

- **(b) "CDD.md is canonical" governance language in design/SKILL.md
  example.** The example at design/SKILL.md §4.4 originally said
  "CDD.md is the normative source for the lifecycle algorithm." Post-#403
  wave, the normative source for the software-cycle lifecycle algorithm
  is CDS.md, not CDD.md (CDD.md owns only the generic CCNF kernel). α
  correctly updated the example to: "`cnos.cds/skills/cds/CDS.md` is the
  normative source for the software-cycle lifecycle algorithm (CDD.md
  owns the generic CCNF kernel)." This is a doctrinal correctness
  improvement, not a regression.

- **(c) Example PRA in post-release/SKILL.md.** Lines 456 and 468 cited
  "CDD §7.6 output format" in an example PRA. The §7.6 anchor never
  existed in current CDD.md; the example was historical. α updated to
  "CDS §Assessment output format" — coherent with the post-wave
  architecture, no semantic drift.

## EXIT

**Verdict: APPROVED**

All AC checks PASS mechanically. No findings. No fix round needed.

Merge action: as collapsed-wave Sub 6, operator δ merges `cycle/411` to
main once γ closeout and CDS-iteration courtesy stub land (γ phase
follows below). β does not push tags or release artifacts for Sub 6 —
this is a docs cleanup cycle, no version bump.
