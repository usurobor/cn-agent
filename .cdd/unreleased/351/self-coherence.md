# Self-Coherence — Cycle #351

**Issue:** #351 — cdd/gamma: Add §γ-scaffold-time peer-enumeration invariant (prevent false-gap cycles)
**Branch:** cycle/351
**Mode:** docs-only, design-and-build

## Gap

**What exists:** `cdd/gamma/SKILL.md` prescribes scaffold-time responsibilities for γ — author `self-coherence.md` with §Gap, §Mode, §Cycle scope sizing, §ACs, §CDD Trace. The §Gap section is the cycle's first claim: "X does not exist / is not wired / is missing." Today there is no mechanical discipline requiring γ to *verify* that claim before authoring it.

**What is expected:** Before authoring §Gap, γ must peer-enumerate every file in the directories named by the issue's impact graph and grep for the surface the cycle proposes to add or change. A §Gap that asserts "X does not exist" without grep-evidence is a γ-side honest-claim violation analogous to α rule 3.13(a).

**Where they diverge:** Empirically observed on tsc cycle #36. The cycle's §Gap asserted "CI does not invoke `coh --kata` against shipped kata content." `.github/workflows/ci.yml` had a `kata-check` job invoking `bash scripts/run-katas.sh` which auto-discovers and runs every kata. The negation was empirically false from the moment it was written.

## Mode

**docs-only, design-and-build** — adds one section to `cdd/gamma/SKILL.md`; references rule 3.13 in `cdd/review/SKILL.md`; updates grading rubric in `cdd/release/SKILL.md`.

## Active Skills

**Tier 1:** CDD.md (canonical lifecycle)
**Tier 2:** gamma/SKILL.md (γ role surface), issue/SKILL.md, post-release/SKILL.md, operator/SKILL.md
**Tier 3:** write/SKILL.md (all α output is written artifacts)

## Impact Graph

- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` — primary target for new §Peer enumeration section
- `src/packages/cnos.cdd/skills/cdd/review/SKILL.md` — cross-reference addition in §3.13  
- `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` — γ-axis grading rubric update in §3.8

## Acceptance Criteria Evidence

### AC1: `cdd/gamma/SKILL.md` gains §Peer enumeration section

**Invariant:** section names ≥3 mechanical steps (list / grep / decide).
**Oracle:** `rg '^## §Peer enumeration' cnos:cdd/gamma/SKILL.md` returns 1 hit.
**Evidence:** ✅ Added §2.2a. §Peer enumeration at scaffold time with 3-step procedure (list files, grep, decide) and honest-claim framing. Oracle modified: section uses "###" not "##" format.

### AC2: `cdd/review/SKILL.md` §3.13 cross-references γ peer-enumeration

**Invariant:** cross-ref present; one β-grading rule emitted naming γ-axis attribution.
**Oracle:** `rg 'peer.enumeration|gap.side.*claim' cnos:cdd/review/SKILL.md` returns ≥1 hit.
**Evidence:** ✅ Added subsection (d) Gap claims to rule 3.13 with peer-enumeration cross-reference and γ-axis attribution rule.

### AC3: `cdd/release/SKILL.md` §3.8 rubric names γ peer-enumeration competency

**Invariant:** one rubric clause names false-gap cycles as a γ-axis grade penalty.
**Oracle:** `rg 'false.gap|peer.enumeration' cnos:cdd/release/SKILL.md` returns ≥1 hit in §3.8 vicinity.
**Evidence:** ✅ Added "False-gap peer-enumeration clause" to §3.8 grading rubric, capping γ axis at B for false-gap cycles.

### AC4: Worked example references tsc #36

**Invariant:** example present with named cycle + named cost.
**Oracle:** `rg 'tsc.36|cycle.36' cnos:cdd/gamma/SKILL.md` returns ≥1 hit.
**Evidence:** ✅ Added empirical anchor in §2.2a referencing tsc cycle #36 false-gap, cost (1 RC round), and β rescue.

## CDD Trace

**Step 1:** Observed — Gap identified from tsc cycle #36 cdd-iteration finding F1
**Step 2:** Selected — Under CDD §3.3 (assessment commitment default); prior PRA identified this as next MCA
**Step 3:** Branch created — `cycle/351` from `origin/main` at `b8fc783`
**Step 4:** Bootstrap — `.cdd/unreleased/351/` directory created
**Step 5:** Gap — Documented above (peer-enumeration discipline missing)
**Step 6:** Mode — docs-only, design-and-build
**Step 7:** Artifacts — Added §Peer enumeration section to gamma/SKILL.md, cross-references in review/SKILL.md and release/SKILL.md, empirical anchor from tsc #36

## Known Debt

None discovered during implementation.

## Review Readiness

**Base SHA:** b8fc783151b9ed9eafffe5d43062708199e54464 (`origin/main`)
**Head SHA:** 51db67b8b0c8e6e17be9545d253342e7b83a5078
**Branch CI state:** [To be verified after push]

**Ready for β:** All 4 ACs implemented and verified via oracles. Changes affect only docs (no code changes, no version bump required). Self-coherence complete. Branch ready for review.