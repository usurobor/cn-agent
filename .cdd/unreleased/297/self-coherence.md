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

---

## ACs

### AC1: TSC formal upstream stated

**Criterion:** `SEMANTICS-NOTES.md` §15.1 explicitly states TSC is the formal upstream for CTB's triadic carrier claim, referencing C≡ §1.2, §2.1–2.2, §3.1–3.4, and TSC Core §7.2.

**Evidence:** §15.1 now opens with "TSC (usurobor/tsc) is the formal upstream for CTB's triadic carrier claim." It explicitly references:
- C≡ §1.2 (three positions hold one-as-two without collapse)
- C≡ §2.1–2.2 (term syntax, equivalence relation, normal form)
- C≡ §3.1–3.4 (α/β/γ evaluators with idempotent profiles)
- TSC Core §7.2 (algebraic independence, no Eckmann-Hilton collapse)

All four required references are present with specific section citations.

**Met:** ✅

---

### AC2: CTB keeps the right level of claim

**Criterion:** The update preserves the distinction — TSC provides formal backing; CTB does not overclaim metaphysical triadic; CTB uses `tri()` as operational carrier.

**Evidence:** §15.1 contains "What CTB claims and does not claim" subsection. It quotes TSC Core §0 directly ("Non-claims: Reality is fundamentally triadic (metaphysical)."), states CTB inherits that restraint, and distinguishes two levels explicitly: "TSC provides formal backing for CTB's `tri()` carrier and checking shape" vs. "CTB's agent-execution model is a practical generalization of the TSC carrier, not a derivation from TSC measurement."

**Met:** ✅

---

### AC3: Witness model mapped to CTB close-outs

**Criterion:** The docs map TSC Operational witness/verdict flow to CTB close-out discipline, showing the HANDSHAKE→MEASURE→WITNESS→DIAGNOSE/VERDICT→ACCEPT/REJECT to orient→intervene→witness→close-out mapping. CTB close-out forms described as agent-execution generalizations of TSC's discipline.

**Evidence:** §15.6 "TSC-Oper: witness model, close-outs, and ctb-check" provides both the state machine flow side-by-side and a table mapping each TSC-Oper outcome to the corresponding CTB close-out form (`accepted`, `repair-needed`, `structured-failure`, `blocked`, `close-with-debt`). Explicitly states CTB does not implement TSC measurement while making the structural parallel clear.

**Met:** ✅

---

### AC4: Witness theater mitigation points to TSC-Oper

**Criterion:** The docs explicitly connect the "witness theater" risk (v0.2 draft §15) to TSC-Oper's W1–W4 witness-independence model. ctb-check should check independently-grounded witnesses, not just field presence.

**Evidence:**
- `SEMANTICS-NOTES.md` §15.6 "Witness theater and TSC-Oper W1–W4" lists all four witnesses (W1 S₃ permutation, W2 role-gauge independence, W3 scale equivariance, W4 variance/Lipschitz) with CTB-specific readings, and states "ctb-check should not merely verify that witness fields are present — that is W2 at minimum, not the full model."
- `LANGUAGE-SPEC-v0.2-draft.md` §15 now ends with: "The structural basis for witness independence is TSC-Oper's W1–W4 witness model. See `SEMANTICS-NOTES.md` §15.6 for the mapping."

**Met:** ✅

---

### AC5: Composition bound informs join semantics without overclaiming

**Criterion:** The docs reference TSC Core §10's composition bound (`C_Σ(P) ≥ geometric_mean(C_Σ(Pᵢ)) - ε_comp`) and state the CTB implication carefully: composed runs must name join loss; joins should not hide degradation; any stronger rule is a CTB policy extension.

**Evidence:** §15.3 "Composition by dimension" now includes "Composition bound and join semantics" subsection with the exact theorem statement, the CTB implication ("composed runs must name coupling penalties (join loss) where applicable, and joins should not hide degradation introduced by composition"), and the explicit non-overclaim ("Any stronger rule … is a CTB policy extension, not a consequence of TSC Core §10 directly. CTB should not derive that rule from the theorem unless it is formally closed.").

**Met:** ✅

---

### AC6: ctb-check dependency noted

**Criterion:** The update records that ctb-check v0 should draw from TSC-Oper's witness-independence pattern. Minimum checker implication: field presence is not enough; future checker stages should prefer independent witness signals.

**Evidence:** §15.6 "ctb-check v0 dependency" explicitly states: "Field presence is the minimum check, not the sufficient one. The minimum checker implication from TSC-Oper is: for each required witness field, verify that the evidence is independently grounded rather than self-asserted where the checking surface permits that distinction."

**Met:** ✅
