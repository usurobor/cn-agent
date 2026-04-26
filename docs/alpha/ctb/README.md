# CTB

CTB is the language Coherent Agents use to articulate skills as deterministic, plan-producing modules. A skill's body evaluates to data — an effect plan — that another runtime executes under explicit capability grants. Skills are scoped modules with callable signatures; they compose by explicit invocation, not by inheritance.

This package holds three documents because the language has three concerns that should not share a file. The Vision sets direction and explains why the language exists at all. The Spec sets the rules an implementation must enforce. The Notes record the conceptual moves behind those rules so that future revisions can be made deliberately rather than through drift.

## Document Map

| Document | Role |
|----------|------|
| [CTB-v4.0.0-VISION.md](./CTB-v4.0.0-VISION.md) | Strategy, motivation, roadmap, and convergence from practice. Non-normative. |
| [LANGUAGE-SPEC.md](./LANGUAGE-SPEC.md) | Normative reference for the skill-module / invocation layer: signatures, scope, dispatch, composition, the effect-plan boundary. |
| [SEMANTICS-NOTES.md](./SEMANTICS-NOTES.md) | Conceptual rationale and design discussion behind the spec. Non-normative. |

## Authority

If the Vision and the Spec disagree on language semantics, the Spec governs. If the Spec and the kernel `.coh` grammar disagree on terms, the kernel governs. The Notes never govern; they preserve reasoning.
