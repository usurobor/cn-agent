# Write a Skill

How to add a new skill to cnos.

---

## Reading order

1. **[WRITING.md](../../../src/packages/cnos.core/mindsets/ENGINEERING.md)** — the voice. Coherent writing: nothing can be removed without loss. Read this first so every sentence you write passes the test.
2. **[skill/SKILL.md](../../../src/packages/cnos.core/skills/skill/SKILL.md)** — the structure. Four-step algorithm: Define (coherence formula), Unfold (sections in necessary order), Rules (imperative + examples), Verify.
3. **[compose/SKILL.md](../../../src/packages/cnos.core/skills/compose/SKILL.md)** — the composition test. Before publishing, name the best alternative composition and the gap it leaves.
4. **[skills/README.md](../../../src/packages/cnos.core/skills/README.md)** — the mechanics. Where to put files, what to name them, artifact classes.

---

## Steps

1. Read all docs above.
2. Create `skills/<name>/`.
3. Write `SKILL.md` following the algorithm in `skill/SKILL.md`.
4. Include a kata surface (embedded or external `kata.md`).
5. Run the compose check: does this skill earn publication, or should it be a composition of existing skills?
6. Commit and push.

---

## Examples

| Skill | Why it's a good model |
|-------|----------------------|
| [skill/SKILL.md](../../../src/packages/cnos.core/skills/skill/SKILL.md) | Teaches the format using the format |
| [code/SKILL.md](../../../src/packages/cnos.eng/skills/eng/code/SKILL.md) | Full-size skill with numbered rules, examples, process section |
| [compose/SKILL.md](../../../src/packages/cnos.core/skills/compose/SKILL.md) | Composition judgment with embedded kata |
