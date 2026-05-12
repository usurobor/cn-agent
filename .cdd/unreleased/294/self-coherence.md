# Self-Coherence — Cycle 294

## Gap

**Issue:** #294 — cn cdd status N — single-shot cycle TLDR (read-only γ tooling)

**Version/Mode:** Standard substantial change cycle per CDD.md §1.1 — introduces new command tooling, spans design/code/tests/docs, requires CDD artifact set.

**Selected Gap:** Missing γ command for cycle state introspection. γ assembles cycle TLDR manually (read issue, git rev-parse, list artifacts, evaluate closure gate, summarize). Information is derivable from git + filesystem + GitHub API but no `cn` command produces structured TLDR. Operator and γ both ask "TLDR current state" 1-2x per cycle without tooling.

**Gap Type:** Tooling gap — mechanical read-only command to project existing state into structured TLDR format.

**Governing CDD.md Clause:** Tooling that reduces manual repetitive coordinator work without changing process semantics (supports §1.4 γ algorithm efficiency).

## Skills

**Tier 1 (Always loaded):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface and algorithm

**Tier 2 (Always applicable):**
- `src/packages/cnos.eng/skills/eng/*` — engineering skill bundle per cnos.eng/README.md

**Tier 3 (Issue-specific):**
- `src/packages/cnos.eng/skills/eng/tool/SKILL.md` — shell script standards, fail-fast, deterministic diagnostics (command is executable shell script)
- `src/packages/cnos.eng/skills/eng/test/SKILL.md` — fixture coverage, golden output, negative cases (AC7 requires test fixture)
- `src/packages/cnos.eng/skills/eng/ux-cli/SKILL.md` — terminal UX, exit codes, readable diagnostics (human-read output for γ/operator)
- `src/packages/cnos.core/skills/skill/SKILL.md` — package manifest discipline for command registry (discoverable via `cn` registry per existing convention)

**Constraints applied:** Read-only constraint (no commits/pushes/API writes), bash + git + gh CLI only, NO_COLOR support, 5-section structured output, deterministic exit codes, role attribution via git log --format=%an.