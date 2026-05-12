# Self-Coherence — Cycle #350

## Gap

**Issue:** #350 — cdd(core): define wave primitive — multi-cycle coordination surface

**Problem:** CDD defines single-cycle artifacts (`.cdd/unreleased/{N}/`) and single-cycle coordination (γ → α → β → close-out). But when δ runs a sequence of related cycles (a "wave"), coordination artifacts live in `/tmp`, inter-cycle state is carried in operator memory, and there is no manifest, status surface, or closure protocol.

**Evidence:** Wave dispatch on 2026-05-12 (#286 → #295 → #294 → #293 → #296 → #305) used `/tmp/wave-status.md` and `/tmp/gamma-prompt-{N}.md` — ephemeral, uncommitted, lost on reboot. See `.cdd/iterations/wave-2026-05-12.md` findings #1, #2, #3.

**Impact:** Wave artifacts are not auditable (no git history), δ dispatch prompts are recreated from scratch each wave, inter-cycle dependencies are tracked in operator memory, wave closure has no canonical signal.

**Version:** CDD 3.15.0  
**Mode:** docs-only, design-and-build

## Skills

**Tier 1 (CDD lifecycle):**
- `CDD.md` — canonical lifecycle and role contract
- `alpha/SKILL.md` — α role surface and algorithm
- `design/SKILL.md` — design artifact requirements (not required for this docs-only cycle)
- `plan/SKILL.md` — implementation sequencing (not required for this single-target cycle)

**Tier 2 (Engineering always-applicable):**
- `eng/document` — documentation conventions and structure
- `eng/git` — commit message and branch conventions

**Tier 3 (Issue-specific):**
- `write/SKILL.md` — writing so the reader understands on first pass (target for operator/SKILL.md updates)