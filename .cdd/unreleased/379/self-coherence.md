<!--
section-manifest:
  planned: [gap, skills, acs, self-check, debt, cdd-trace, review-readiness]
  completed: [gap, skills]
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

## Skills

**Tier 1 — always-on for α/CDD:**

- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical algorithm and lifecycle authority. §1.4 large-file authoring rule applied: SKILL.md authored section-by-section with HTML-comment manifest; self-coherence.md authored section-by-section per α/SKILL.md §2.5.
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface. §2.2 artifact order followed (skill artifact → tests → code → docs → self-coherence); §2.5 incremental self-coherence; §2.6 pre-review gate exercised in §Review-readiness; §2.6 row 14 retroactive author-amend executed (worktree config had γ identity from session start; rebase --exec --reset-author applied across α's six commits before signaling).

**Tier 2 — always-applicable engineering:**

- `src/packages/cnos.eng/skills/eng/go/SKILL.md` — Go authoring conventions. Applied to renderer evolution (`activate.go`): small functions, explicit fallback path, error-free reads on optional skill, line-oriented parser instead of regex.
- `src/packages/cnos.core/skills/write/SKILL.md` — prose discipline for the SKILL.md artifact (clarity over decoration; rule of three: principle → unfold → kata).

**Tier 3 — issue-specific (per issue body §Skills to load + γ scaffold §Tier 3 list):**

- `src/packages/cnos.core/skills/skill/SKILL.md` — load-bearing for AC2 conformance. Frontmatter keys (`name`, `description`, `artifact_class: skill`, `kata_surface`, `governing_question`, `triggers`, `scope`, `inputs`, `outputs`, `parent`, `visibility`, `requires`, `calls`) all present; body follows Define → Unfold → Rules → Verify → Final test → Kata progression; embedded kata at §7 satisfies `kata_surface: embedded`.
- `src/packages/cnos.core/skills/write/SKILL.md` — prose structure.
- `src/packages/cnos.core/skills/design/SKILL.md` — design-and-build mode design step (substrate independence between cnos and per-hub layers; single source of truth for the activation procedure; defer realization choices about which body capability is used).
- `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` — cited and disambiguated in SKILL.md §2.4 (AC6). Not modified — non-goal per γ scaffold failure-mode 10.
- `src/packages/cnos.eng/skills/eng/go/SKILL.md` — renderer evolution (AC7).
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — form authority. No in-cycle issue-pack reconciliation was required; the issue body remained authoritative throughout.
