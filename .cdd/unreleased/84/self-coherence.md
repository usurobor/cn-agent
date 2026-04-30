## §Gap

**Issue:** #84 — agent(conduct): make reflection a core CA continuity requirement  
**Version/mode:** CDD 3.15.0, MCA  
**Branch:** cycle/84

**Gap:** `ca-conduct` covered PLUR, accountability, shipping, ownership, truth, learning, and preparation but did not explicitly require reflection as core coherent-agent behavior. When reflection is only an available skill or operations detail, agents complete external work while losing the learning that should change future behavior — operational patterns stay in chat, mistakes repeat across sessions, and MCIs do not migrate into skills, mindsets, or threads.

**Target:** Add a Reflect section to both `ca-conduct/SKILL.md` and its doctrine mirror `CA-CONDUCT.md`, referencing `reflect` and `mci` skills, naming concrete triggers, and explaining why reflection is required (not optional) continuity behavior.

---

## §Skills

**Tier 1:**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` (v3.15.0)
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`

**Tier 2 (Writing bundle):**
- `src/packages/cnos.eng/skills/eng/document/SKILL.md`
- `src/packages/cnos.core/skills/skill/SKILL.md`

**Tier 3 (issue-specific):**
- `src/packages/cnos.core/skills/agent/ca-conduct/SKILL.md` (primary artifact — modified)
- `src/packages/cnos.core/doctrine/CA-CONDUCT.md` (primary artifact — modified)
- `src/packages/cnos.core/skills/agent/reflect/SKILL.md` (referenced; not modified)
- `src/packages/cnos.core/skills/agent/mci/SKILL.md` (referenced; not modified)

---

## §ACs

**AC1: CA conduct includes reflection as core behavior**

Evidence: `ca-conduct/SKILL.md` line 101–116 and `CA-CONDUCT.md` line 93–108 both contain a `## Reflect` section. The section opens with "Agents wake fresh. Context does not persist across sessions by default." and introduces reflection via "A coherent agent captures operational residue before it disappears..." — framing it as required behavior, not an available option.

**PASS.**

---

**AC2: Reflection continuity is explained**

Evidence: The opening of the Reflect section explicitly covers all four required points:
- "Agents wake fresh." ✓
- "Context does not persist across sessions by default." ✓
- "Threads and reflections are continuity" ✓
- "future agents can only learn from what was recorded" ✓

**PASS.**

---

**AC3: Concrete triggers are named**

Evidence: The section's bullet list covers all five named triggers:
- "Reflect before ending substantial work" ✓ (trigger 1)
- "Capture MCIs immediately" ✓ (trigger 2 — when an MCI appears)
- "Reflect after a mistake or repeated pattern." ✓ (trigger 3)
- "Reflect after a judgment, triage, or classification future sessions may inherit." ✓ (trigger 4)
- "Reflect at daily/heartbeat cadences when appropriate." ✓ (trigger 5)

**PASS.**

---

**AC4: Reflect and MCI skills are referenced**

Evidence: The closing line of the Reflect section reads:
> Use `reflect` (`src/packages/cnos.core/skills/agent/reflect/SKILL.md`) for structured reflection and `mci` (`src/packages/cnos.core/skills/agent/mci/SKILL.md`) when the result is a behavior-changing insight.

Both full paths are present in both files.

**PASS.**

---

**AC5: No stale path remains**

Evidence: `grep -r "mindsets/CA\.md"` in the diff surfaces zero hits. Neither modified file references `cnos.core/mindsets/CA.md`. The implementation also does not add any new stale path.

Verification: `grep -rn "mindsets/CA" src/packages/cnos.core/skills/agent/ca-conduct/ src/packages/cnos.core/doctrine/` → 0 matches.

**PASS.**

---

**AC6: Conduct mirrors stay aligned**

Evidence: Both `ca-conduct/SKILL.md` and `CA-CONDUCT.md` received identical Reflect section content. Both are live authoritative surfaces — the SKILL.md is the skill (with frontmatter, kata surface) and CA-CONDUCT.md is the doctrine mirror (plain text). Both are updated in the same commit (`ce8b8108`).

**PASS.**
