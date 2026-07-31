# CTB

CTB is a triadic agent-composition language. `tri()` is the common carrier. A skill is a narrow agent viewed as `tri(input, transform, witnessed-output)`. An agent is a scoped process viewed as `tri(orientation, intervention, witness)`. A protocol is a relation among agents viewed as `tri(roles/capabilities, interaction, close-outs)`. Composition preserves the triadic carrier and closure requires an inspectable witness.

This package holds two documents because the language has concerns that should not share a file. The Spec sets the rules current implementations enforce. The Notes record the conceptual moves behind those rules so that future revisions can be made deliberately rather than through drift.

## Document Map

| Document | Role |
|----------|------|
| [LANGUAGE-SPEC.md](./LANGUAGE-SPEC.md) | The normative language spec. Governs skill-module shape, signature, scope, composition, and the effect boundary. |
| [SEMANTICS-NOTES.md](./SEMANTICS-NOTES.md) | Conceptual rationale and harvest notes. Non-normative. |

## Authority

`LANGUAGE-SPEC.md` is the canonical spec — it governs the language model and conformance.

If the Spec and the kernel `.coh` grammar disagree on terms, the kernel governs. The Notes never govern; they preserve reasoning.
