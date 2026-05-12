# Self-Coherence — Cycle #323

## Gap

**Issue:** #323 — fix(activate): scanner misses threads/inbox/

**Problem:** `cn activate` scans `threads/in`, `threads/mail`, and `threads/archived` but does not scan `threads/inbox/`. The canonical reference hub (cn-sigma) uses `threads/inbox/` as its primary inbound message surface, making activation output blind to the actual inbox.

**Expected:** `cn activate` should scan `threads/inbox/` in addition to the current surfaces.

**Version/Mode:** MCA (Minimal Correct Action) — direct fix to existing scanner function.

## Skills

**Tier 1a — CDD authority (loaded):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md`
- `src/packages/cnos.cdd/skills/cdd/SKILL.md`
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`

**Tier 2 — General engineering (coding bundle, loaded):**
- `src/packages/cnos.eng/skills/eng/go/SKILL.md` (primary language skill)
- `src/packages/cnos.eng/skills/eng/ux-cli/SKILL.md` (CLI interface standards)
- `src/packages/cnos.eng/skills/eng/test/SKILL.md` (testing discipline)

**Tier 3 — Issue-specific:** None explicitly specified in issue #323.

**Work Shape:** Small local bugfix (L5 — local correctness)
**Active Skills:** go + ux-cli + test