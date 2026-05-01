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
