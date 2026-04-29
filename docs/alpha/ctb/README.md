# CTB

CTB is a triadic agent-composition language. `tri()` is the common carrier. A skill is a narrow agent viewed as `tri(input, transform, witnessed-output)`. An agent is a scoped process viewed as `tri(orientation, intervention, witness)`. A protocol is a relation among agents viewed as `tri(roles/capabilities, interaction, close-outs)`. Composition preserves the triadic carrier and closure requires an inspectable witness.

This package holds four documents because the language has concerns that should not share a file. The Vision sets direction and explains why the language exists. The v0.1 Spec sets the rules current implementations enforce. The v0.2 draft defines the migration target (agent as primitive, triadic carrier, composition operators). The Notes record the conceptual moves behind those rules so that future revisions can be made deliberately rather than through drift.

## Document Map

| Document | Role |
|----------|------|
| [CTB-v4.0.0-VISION.md](./CTB-v4.0.0-VISION.md) | Strategy, motivation, roadmap. Non-normative. |
| [LANGUAGE-SPEC.md](./LANGUAGE-SPEC.md) | v0.1 skill-module baseline. Normative for current v0.1 conformance. |
| [LANGUAGE-SPEC-v0.2-draft.md](./LANGUAGE-SPEC-v0.2-draft.md) | Draft-normative agent-module target. Supersedes v0.1 conceptually, not yet promoted. |
| [SEMANTICS-NOTES.md](./SEMANTICS-NOTES.md) | Conceptual rationale and harvest notes. Non-normative. |

## Authority

Until v0.2 is promoted, `LANGUAGE-SPEC.md` governs current v0.1 conformance. `LANGUAGE-SPEC-v0.2-draft.md` governs the migration target.

If v0.2 and v0.1 disagree, the disagreement is expected and must be resolved by promotion or revision, not by inference.

If the Vision and the Spec disagree on language semantics, the Spec governs. If the Spec and the kernel `.coh` grammar disagree on terms, the kernel governs. The Notes never govern; they preserve reasoning.
