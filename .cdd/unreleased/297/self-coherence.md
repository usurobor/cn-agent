## Gap

**Issue:** #297 — docs(ctb): make TSC formal grounding explicit for tri(), witnesses, and composition  
**Version/mode:** MCA — documentation update only, no code changes  
**Branch:** cycle/297

CTB's current `SEMANTICS-NOTES.md` §15.1 references C≡ and the α/β/γ evaluators by paraphrase, without naming TSC as the formal upstream or citing specific sections. The TSC repo (usurobor/tsc) already contains:

- C≡ formal term algebra with `tri(·,·,·)`, equivalence, normal form, and evaluators
- Proof that α/β/γ evaluators are algebraically independent via distinct idempotent profiles
- TSC Core measurement framework with dimensional scores, aggregate C_Σ, confidence intervals, provenance, and a composition bound
- TSC Operational witness model, witness floors, verification state machine, fail-fast verdict logic, and provenance bundle

Without explicit grounding, CTB risks re-deriving TSC concepts informally, weakening the tri() claim to an under-specified design metaphor, and building ctb-check without using TSC's witness-independence pattern.

The fix is documentation: update `SEMANTICS-NOTES.md` to make the TSC relationship explicit, with specific section references.

---

## Skills

**Tier 1:**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` (canonical lifecycle)
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (α role)

**Tier 2:**
- `src/packages/cnos.core/skills/write/SKILL.md` (writing — docs-only change)

**Tier 3 (per dispatch):**
- `src/packages/cnos.core/skills/write/SKILL.md`

No code skill applies; this is a docs-only MCA. Bootstrap not required (no new version snapshot directory; the change updates an existing non-normative notes file).
