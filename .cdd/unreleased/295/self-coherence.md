# Self-Coherence Report: Issue #295

## Gap

**Issue:** #295 - `cn dispatch` — role identity rotation primitive (unblocks #286 encapsulation)

**Version:** Development (unreleased)
**Mode:** Substantial change - new command implementation with backend architecture

**What exists:** CDD has triadic role structure (α/β/γ) with defined dispatch prompt format and cycle/{N} branch convention. γ currently produces dispatch prompts, but the operator manually instantiates α and β sessions. There is no harness-provided dispatch mechanism to enable autonomous role coordination.

**What is expected:** γ should be able to invoke isolated α or β role identity sessions through a single command surface (`cn dispatch`) that preserves CDD role isolation and branch/artifact continuity. The command must provide cognitive isolation via fresh Claude CLI sessions while maintaining shared branch and .cdd/unreleased/{N}/ coordination.

**Where they diverge:** The dispatch prompt format exists, the cycle/{N} branch convention exists, and the role skills exist, but there is no command/harness boundary that turns a prompt into an isolated α or β role invocation. Without this primitive, γ cannot close #286's encapsulation gap, leaving the operator as the manual relay channel for every α/β rotation.

**Gap significance:** This preserves CDD correctness but wastes operator attention and prevents γ from becoming autonomous inside the cycle. It is the M12 friction class from the 3.61.0 cycle.

## Skills

**Tier 1 (CDD Lifecycle):**
- `CDD.md` — canonical lifecycle and role contract
- `alpha/SKILL.md` — α role surface and algorithm
- `write/SKILL.md` — coherent writing standards

**Tier 2 (Engineering - Always Applicable):**
- All `eng/*` skills as required

**Tier 3 (Issue-Specific):**
- `cnos.core/skills/design` — CLI/harness/runtime boundary and no future-as-present
- `eng/tool` — command behavior, prereq checks, deterministic diagnostics
- `eng/test` — fixtures, positive/negative proof, stub backend
- `eng/ux-cli` — help text, structured output, failure messages
- `eng/go` — Go implementation in the `cn` CLI

**Active skill set loaded and constraining generation.**