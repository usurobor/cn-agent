<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness] -->
<!-- completed: [Gap, Skills] -->

# α self-coherence — #370

## Gap

`COHERENCE-CELL.md` (#364, merge `32b126e4`) names the coherence-cell *organism* — receipt rule, four-way structural separation, role-as-cell-function. `RECEIPT-VALIDATION.md` (#367, merge `37ac1c75`) names the parent-facing *receptor* — `V`'s typed interface, verdict/decision distinction, δ-authoritative validation. `#366`'s roadmap §"Recursion the implementation must instantiate" carries the algorithm in informal mathematical notation but is a roadmap, not a doctrine surface.

The *algorithm* — the recursion that takes a closed cell at scope `n` and projects it as α-matter, β-discrimination, and γ-coordination at scope `n+1` — is implicit across these four surfaces but stated nowhere. Without a kernel-layer doc that names the algorithm, Phase 3 (validator) would have to derive V's contract from a recursion no doctrine names; Phase 4 (δ split) would intuit δ's signature per cycle; Phase 6 (ε relocation) would move ε without a kernel statement of its signature; Phase 7 (`CDD.md` rewrite) would have to derive the kernel mid-rewrite — exactly the failure mode two-layer separation exists to prevent.

This cycle (`#370`, Phase 1.5 of `#366`) closes that gap by producing `COHERENCE-CELL-NORMAL-FORM.md` — the substrate-independent kernel-layer companion at the doctrine layer.

**Mode:** docs-only (single new doctrine file under `src/packages/cnos.cdd/skills/cdd/` + cycle evidence under `.cdd/unreleased/370/`).

**Cycle:** 370. **Branch:** `cycle/370`. **Base SHA:** `704365d2`.

## Skills

**Tier 1 (always-on for α CDD):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical algorithm
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface (esp. §1.4 large-file authoring discipline, §2.5 incremental self-coherence, §2.6 pre-review gate, §2.7 review-request)

**Tier 2 (always-applicable):**
- `src/packages/cnos.eng/skills/eng/*` — markdown / prose authoring (applied implicitly via the write skill)
- `src/packages/cnos.core/skills/write/SKILL.md` — every α output is a written artifact

**Tier 3 (issue-specific):**
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — form authority for any issue-pack reconciliation during cycle (not exercised — the issue and γ-scaffold were stable across the cycle)
- `src/packages/cnos.core/skills/design/SKILL.md` — kernel-articulation discipline: substrate-independence, two-layer separation, single source of truth, defer realization choices. **Load-bearing** for this cycle's authoring — without it the kernel would drift into operational details and fail AC7
- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` — predecessor doctrine; read as source-of-truth for receipt rule, four-way separation, role-as-cell-function; the organism this kernel describes algorithmically
- `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` — typed-interface design peer; read as source-of-truth for V's interface, verdict/decision distinction, δ-authoritative validation; the receptor this kernel positions at step 4–5
