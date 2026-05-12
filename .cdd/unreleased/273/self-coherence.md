# Self-coherence — cycle #273

## Gap

**Issue:** [#273](https://github.com/usurobor/cnos/issues/273) — Rebase-collision integrity guard: prevent silent loss of upstream content at integration  
**Version:** 3.61.0+  
**Mode:** triadic, non-release

The gap is silent loss of upstream content during rebase-integration cycles. Two confirmed instances in γ #268: COHERENCE-FOR-AGENTS.md and CTB vision §8.5.2 were added on `main` while γ branches were in flight, rebased away silently, and required manual restoration post-hoc. No existing skill or CI mechanism detects this failure class.

Current state: operator-dependent manual detection via post-merge content review.  
Target state: pre-push git hook blocks upstream content loss before it reaches remote, with bypass for intentional deletions.

The gap blocks process-integrity (P1): routine CDD cycles risk losing doctrine content silently.

## Skills

**Tier 1:**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` (canonical lifecycle)
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (α role surface and algorithm)

**Tier 2:**
- `src/packages/cnos.eng/skills/eng/tool/SKILL.md` (mechanical script standards)
- `src/packages/cnos.eng/skills/eng/document/SKILL.md` (documentation verification against source)
- `src/packages/cnos.eng/skills/eng/test/SKILL.md` (prove invariants from testing)
- `src/packages/cnos.eng/skills/eng/ship/SKILL.md` (ship code to production/merge)

**Tier 3:**
- `src/packages/cnos.core/skills/skill/SKILL.md` (skill-program coherence for new eng/ship surface)