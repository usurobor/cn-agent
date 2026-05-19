<!--
section-manifest:
  planned: [gap, skills, acs, self-check, debt, cdd-trace, review-readiness]
  completed: [gap]
-->

# Self-coherence — α #379

## Gap

**Issue:** [#379 — agent/activate skill: single source of truth for agent self-activation](https://github.com/usurobor/cnos/issues/379).
**Mode:** design-and-build. The design surface is the issue body (Source of Truth table, AC3 load order, AC4 capability matrix, AC7 renderer-interpolation contract). α executes design and build in the same cycle — no separate committed `DESIGN.md`.
**Version:** unreleased (cycle directory `.cdd/unreleased/379/`); release boundary follows β's merge and γ's cycle close.
**Base SHA:** `7a7f7152af649ea93cebe5f909fbf1397d809547` (γ scaffold base = origin/main HEAD at dispatch time; cycle branch rebased onto current `origin/main` — confirmed in §Review-readiness).

**Named gap.** cnos exposes activation as procedure-in-Go: `src/go/internal/activate/activate.go` `writePrompt` hardcodes the section ordering and the "Read first" list in in-Go literals. Every other cnos behavior is a skill (`cnos.core/skills/agent/{cap,clp,mca,mci,coherent,agent-ops}/SKILL.md` is the precedent). Two consequences follow: (1) the "everything is a skill" invariant is violated for activation; (2) non-cn bodies (Claude Code on the web, Codex sessions, Claude.ai with WebFetch) cannot reach the procedure because it does not exist as a fetchable artifact — they fall into the "I wake up incoherent by default" failure named in cn-sigma `threads/adhoc/20260325-session2-learnings.md` §1.

**Closure shape.** This cycle creates `src/packages/cnos.core/skills/agent/activate/SKILL.md` as the single source of truth for agent activation, and evolves `src/go/internal/activate/activate.go` to render from it. The skill prescribes a six-item canonical load order (Kernel → CA skills → Persona → Operator → hub state → identity confirmation), a three-tier body-capability matrix (shell+git preferred, HTTP fetch only, no-fetch operator-injected), a README router template hubs adopt verbatim, and a disambiguation paragraph distinguishing this skill from `cnos.cdd/skills/cdd/activation/SKILL.md` (CDD repo activation, a distinct concern). The renderer parses the skill's §4.1 machine-readable load-order block at every invocation and emits in that order; an in-Go fallback preserves manifest-only and no-deps behavior when the skill is not yet vendored.

**Non-goals (held).** No hub README patches (cn-sigma, cn-pi) — deferred to per-hub cycles. No `cnos.xyz` rendering service. No `cn doctor` validation of router presence. No edits to `cnos.cdd/skills/cdd/activation/SKILL.md` beyond citation in §2.4 disambiguation. No changes to the `cn activate` CLI surface (flags, arguments, exit codes — `src/go/internal/cli/cmd_activate.go` unchanged).
