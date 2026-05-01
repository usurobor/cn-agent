## Gap

**Issue:** #325 — cdd: systematic lifecycle audit and refactor for role, artifact, and release coherence
**Version/mode:** MCA — lifecycle coherence patch (spec/skill refactor; no runtime implementation)
**Priority:** P1

CDD defines the development lifecycle for substantial cnos work. Recent cycles exposed gaps where CDD itself was ambiguous or internally inconsistent:

- The coordination model (sequential bounded dispatch vs. persistent polling) is not declared, so each cycle must infer it
- α close-out timing is broken in bounded dispatch: α is told to write close-out "after β merge" but α has already exited
- No role/artifact ownership matrix exists: a mechanical lookup by artifact yields no single source for owner + verifier + preconditions + failure consequences
- γ closure gate does not make `gamma-closeout.md` itself a precondition for δ tag/release
- Small-change path does not state which close-outs and release gates collapse
- No failure-mode regression surface exists to catch the known patterns (stale unreleased dirs, missing close-outs, missing RELEASE.md, δ tag before γ closure)

The work is entirely spec/skill: edits to CDD.md and all role skills. No runtime implementation.

Design artifact: `.cdd/unreleased/325/alpha-design.md`

---

## Skills

**Tier 1a (CDD authority):**
- `CDD.md` — canonical lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface

**Tier 1b (lifecycle sub-skills):**
- `src/packages/cnos.cdd/skills/cdd/design/` — not loaded separately; design principles applied via Tier 3 `cnos.core/skills/design`

**Tier 2 (general engineering):**
- Writing bundle applies: all outputs are written artifacts

**Tier 3 (issue-specific):**
- `src/packages/cnos.core/skills/write/SKILL.md` — all CDD edits are written artifacts; write skill governs prose quality, one-governing-question-per-doc discipline, stable-facts-once rule
- `src/packages/cnos.core/skills/design/SKILL.md` — lifecycle refactor is an authority-boundary and state-machine design problem; design skill governs: one reason to change per boundary, policy above detail, interfaces truthful, no premature canonicalization
- `src/packages/cnos.eng/skills/eng/test/SKILL.md` — regression fixtures and lifecycle checklists must prove invariants, not just examples; negative space is mandatory; test the failure-mode regression surface
