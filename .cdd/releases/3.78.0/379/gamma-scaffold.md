---
cycle: 379
issue: "https://github.com/usurobor/cnos/issues/379"
date: "2026-05-19"
dispatch_configuration: "§5.1 canonical multi-session (one `claude -p` per role) — escalated from §5.2 hub default per operator/SKILL.md §5.3 because: 7 ACs, new contract surface (new skill + Go renderer), cross-repo proposal"
base_sha: "7a7f7152af649ea93cebe5f909fbf1397d809547"
source_proposal:
  repo: "usurobor/cn-sigma"
  path: ".cdd/iterations/cross-repo/cnos/agent-activate-skill/"
  disposition: "accepted"
  bundle_commit: "1a4e25f75"
  filing_commit: "8f153c15e"
  delta: "AC4 capability-matrix refinement folded in from cn-sigma branch HEAD bdda457f5 per threads/adhoc/20260519-git-read-and-untested-limits.md — single body-capability gate replaced with three-tier capability matrix"
mode: "design-and-build"
mode_rationale: "Design converges in issue body during γ; α writes the skill and evolves the Go renderer in the same cycle. Not MCA — no separate committed DESIGN.md path; design is in the issue body."
---

# γ Scaffold — #379

## Gap

cnos exposes activation as a behavior defined in Go code: `src/go/internal/activate/activate.go` `writePrompt` hardcodes the load order (persona, operator, kernel, deps, latest reflection) and the section emission (`## Read first`, `## Kernel`, `## Persona`, …). Every other cnos behavior is defined as a skill with frontmatter, triggers, governing question, and load semantics — `cnos.core/skills/agent/{cap,clp,mca,mci,coherent,agent-ops}/SKILL.md` is the precedent. Activation is the exception.

Two consequences follow:

1. **The "everything is a skill" invariant is violated for activation.** A non-cn body cannot fetch the activation procedure because it does not exist as a fetchable artifact.
2. **Bodies without `cn` reinvent activation per session.** Claude Code on the web, Codex sessions, and Claude.ai chats with WebFetch must depend on the operator manually pointing at the hub and reading files in an improvised order — the "I wake up incoherent by default" failure documented in cn-sigma `threads/adhoc/20260325-session2-learnings.md`.

This cycle creates `src/packages/cnos.core/skills/agent/activate/SKILL.md` as the single source of truth and evolves `cn activate` to render from it. The skill is body-agnostic and names the capability matrix (shell+git, fetch-only, no-fetch) so future bodies pick the load path their environment actually supports.

## Peer enumeration (gamma/SKILL.md §2.2a)

Before authoring §Gap, γ enumerated the directories named by the issue's impact graph and grepped for the proposed surface:

```bash
# (a) The target skill directory does not yet exist
ls src/packages/cnos.core/skills/agent/activate/
# → No such file or directory  (confirms additive, not consolidation)

# (b) The agent/ namespace exists with peer skills (cap, clp, mca, mci, coherent, agent-ops, …)
ls src/packages/cnos.core/skills/agent/
# → agent-ops/ ca-conduct/ cap/ cbp/ clp/ coherence-test/ coherent/ communicating/
#   configure-agent/ emoji-language/ hello-world/ human-interaction/ mca/ mci/
#   onboarding/ reflect/ self-cohere/
#   — no activate/ entry; the new skill slots in alongside these

# (c) "agent/activate" path is not referenced anywhere in src/
rg 'agent/activate' src/
# → No files found  (no dangling references to address)

# (d) The same-word peer at cdd/activation/ exists and is shipped
ls src/packages/cnos.cdd/skills/cdd/activation/
# → SKILL.md (50381 bytes), templates/  — distinct concern (CDD repo activation),
#   must be cited and disambiguated per AC6

# (e) The Go renderer to evolve exists at the path AC7 names
ls src/go/internal/activate/
# → activate.go (11475 bytes), activate_test.go (19145 bytes)

# (f) writePrompt currently hardcodes the ordering
grep -n 'writePrompt\|Read first' src/go/internal/activate/activate.go
# → 66: writePrompt(opts.Stdout, ...)
#   250: func writePrompt(...)
#   263: fmt.Fprintf(w, "## Read first\n")
#   — confirms the hardcoded ordering AC7 displaces
```

All six checks confirm: gap is real and additive at the canonical path. The peer skill at `cdd/activation/SKILL.md` is the named "same word, different concern" surface to cite per AC6.

## Mode

**design-and-build.** Design converges in the issue body during γ (Source of Truth table, AC4 capability matrix, AC3 load order, renderer interpolation contract). α executes both the design step (write the skill conforming to `cnos.core/skills/skill/SKILL.md`) and the build step (evolve `src/go/internal/activate/activate.go` to render from the skill) inside this cycle.

Mode is **not MCA**: no separate committed `DESIGN.md` / `PLAN.md` at a stable path. The issue body is the design surface; the source-of-truth rows cite the artifact paths the cycle will produce/evolve, not pre-existing design docs.

## ACs reference

All 7 ACs are defined in the issue body at https://github.com/usurobor/cnos/issues/379 — they are the binding contract for α/β. α maps each AC to evidence in `.cdd/unreleased/379/self-coherence.md` §ACs.

## Acceptance posture summary (γ pre-flagged to α)

The issue body is authoritative. The notes below are γ's compressed reading of each AC for handoff; they do not override the issue.

- **AC1 — file exists at canonical path.** `src/packages/cnos.core/skills/agent/activate/SKILL.md`, non-empty, ≥200 lines. Sibling to `cap/`, `clp/`, `coherent/`, `agent-ops/` under `agent/`. No alternative path is acceptable (e.g., `doctrine/`, `activate/` at the package root, `activation/` under core).
- **AC2 — conforms to skill skill format.** YAML frontmatter parses with required keys: `name`, `description`, `artifact_class: skill`, `governing_question`, `triggers`, `scope`, `inputs`, `outputs`. Body matches `cnos.core/skills/skill/SKILL.md`'s prescribed section structure (core principle, Define/Unfold/Rules/Verify or equivalent). Free-form prose without skill structure is a binding AC2 finding.
- **AC3 — body-agnostic six-item load order.** Numbered list naming, in this exact order: (1) Kernel from cnos, (2) CA skills from cnos, (3) Persona from hub, (4) Operator from hub, (5) hub state (deps, latest reflection, memory, threads), (6) identity confirmation. Layering rule (soul vs identity, per cn-sigma `threads/adhoc/20260325-session2-learnings.md` §3) cited. Order must be consistent with what `cn activate` emits after AC7 — i.e., the skill is the canonical ordering and the renderer follows.
- **AC4 — body-capability matrix with three tiers.** Explicit section enumerating at least three tiers, each with its load path:
  - **(a) shell + git** — Claude Code on web/phone/CLI, Codex. `git clone https://github.com/usurobor/cnos.git` + local file reads for kernel & CA skills; clone or local-read hub for Persona/Operator/state. Faster and atomic; preferred when available.
  - **(b) HTTP fetch only** — constrained Claude.ai with WebFetch, no shell. Fetches raw GitHub URLs per file.
  - **(c) no fetch** — operator must inject the bootstrap (degraded path, named honestly).
  Section must state tier (a) is preferred when available. A single uniform load path or treating WebFetch as the only non-cn surface is a binding AC4 finding.
- **AC5 — README router template.** A labeled "README router template" section with a fenced markdown block underneath. Self-contained: a hub adopts it verbatim by pasting and substituting only the hub URL. Routes an AI body told "Activate as `<hub URL>`" into fetching the activate skill. Template requiring per-hub edits beyond the hub URL is a binding AC5 finding.
- **AC6 — disambiguation from cdd/activation/SKILL.md.** Paragraph citing both paths exactly:
  - `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` (CDD repo activation — activate a repo under the CDD protocol)
  - `src/packages/cnos.core/skills/agent/activate/SKILL.md` (agent activation — activate an agent identity at a hub)
  Distinction stated in one sentence: they share a word, not a concern. Both paths must be spelled correctly.
- **AC7 — `cn activate` reads the skill, not in-Go constants.** `src/go/internal/activate/activate.go` `writePrompt` no longer sources the section ordering / "Read first" list from in-Go literals; instead it reads from the activate skill file (canonically `.cn/vendor/packages/cnos.core/skills/agent/activate/SKILL.md` when the bundle is vendored; with a sensible fallback or diagnostic when absent — current behavior on `manifest-only` deps state must not regress). `activate_test.go` carries a skill-as-source-of-truth test: editing the skill's section ordering changes the renderer's output. The interpolation surface (template positions ↔ hub-state fields, e.g. `{hub_persona_path}`, `{hub_latest_reflection}`) is either declared in the skill or in the renderer's documented contract — not implicit substitution.

## Skills to load

**Tier 1 (always-on for α/β CDD):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical algorithm
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (α) / `beta/SKILL.md` (β)

**Tier 2 (always-applicable):**
- `src/packages/cnos.eng/skills/eng/*` — markdown + Go authoring
- `src/packages/cnos.core/skills/write/SKILL.md` — every α output is a written artifact

**Tier 3 (issue-specific, per issue body §Skills to load):**
- `src/packages/cnos.core/skills/skill/SKILL.md` — defines the format the new skill must match (AC2 conformance)
- `src/packages/cnos.core/skills/write/SKILL.md` — clarity matters; this skill instructs LLM bodies and humans
- `src/packages/cnos.core/skills/design/SKILL.md` — design-and-build mode's design step (substrate-independence, single source of truth, defer realization choices)
- `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` — same word, different concern; cited and disambiguated (AC6)
- `src/packages/cnos.eng/skills/eng/go/SKILL.md` — Go renderer evolution (AC7)
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — form authority for any in-cycle issue-pack reconciliation

`skill/SKILL.md` is the load-bearing Tier 3 skill for AC1–AC6 (the SKILL.md must conform to it); `eng/go/SKILL.md` is load-bearing for AC7 (the Go renderer must conform to it). Both must be loaded — neither displaces the other.

## Dispatch configuration

**§5.1 canonical multi-session.** γ produces α and β prompts at `/tmp/dispatch-379/{alpha,beta}-prompt.md`; δ (operator) routes each role to a fresh `claude -p` session sequentially. No Agent-tool sub-agent dispatch from a parent session. Each role loads its skills fresh from disk on dispatch.

**§5.3 escalation rationale.** Hub default is §5.2 (γ-in-hub authors prompts and dispatches via cn dispatch). Escalating to §5.1 because:
- AC count = 7 (at the §5.3 multi-session signal threshold)
- New contract surface: net-new skill artifact AND evolving the Go renderer (two coupled surfaces)
- Cross-repo proposal: cn-sigma bundle accepted with explicit delta (AC4 refinement); lineage tracking and review benefit from fresh-session role isolation

Dispatch order:
1. δ dispatches α with `/tmp/dispatch-379/alpha-prompt.md` → α writes SKILL.md + evolves activate.go + activate_test.go, commits to `cycle/379`, writes `self-coherence.md` with review-readiness signal, exits
2. δ dispatches β with `/tmp/dispatch-379/beta-prompt.md` → β reviews on `cycle/379`; on APPROVE, merges into `main` per `beta/SKILL.md` §1, writes `beta-review.md` + `beta-closeout.md`, exits
3. δ runs release-boundary actions (cycle-dir move, RELEASE.md, tag) after γ writes RELEASE.md and `gamma-closeout.md`

Per the operator/SKILL.md dispatch convention, α gets `--allowedTools "Read,Write,Bash"` and `--permission-mode acceptEdits`; β gets the same allowedTools (β needs Bash to run the AC7 Go test and the AC1/AC2/AC4–AC6 oracle commands during review).

## Failure modes to guard against (γ-side)

1. **Skill format drift (AC2).** `cnos.core/skills/skill/SKILL.md` prescribes a specific frontmatter set and a Define/Unfold/Rules/Verify body shape. Watch for: missing `artifact_class: skill`, missing `governing_question`, free-form prose with `##` headings that don't match the skill skill's section convention. β must verify against `skill/SKILL.md` literally, not by template eyeballing.

2. **Load-order drift (AC3).** The six items must be in this exact order: Kernel, CA skills, Persona, Operator, hub state, identity confirmation. The current Go renderer emits in a different order (Persona, Operator, kernel, deps, latest reflection). AC7 says the skill is canonical; the renderer must align to the skill, not the reverse. Watch for: skill ordering matched to the existing Go emission instead of the canonical six-item order.

3. **Capability matrix collapsed to one tier (AC4).** This is the cn-sigma post-filing refinement that this cycle's issue body now carries. Watch for: skill says "fetch the skill URL" without distinguishing the shell+git tier from the WebFetch-only tier; omits the no-fetch degraded path; treats WebFetch as the universal non-cn surface. Tier (a) shell+git is the preferred path and must be named first when listed by preference.

4. **README router template requires per-hub customization beyond URL (AC5).** A template that needs the hub to fill in identity name, list of operator-instructions, or load order is not "self-contained." The hub edits one thing: the hub URL. The activate skill provides the rest. Watch for: template referencing `<PERSONA_NAME>` or `<OPERATOR_INSTRUCTIONS>` as required substitutions.

5. **AC6 path misspelling.** Both paths must be exactly:
   - `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md`
   - `src/packages/cnos.core/skills/agent/activate/SKILL.md`
   Common typo: `cnos.core/skills/cdd/activation/...` (wrong package) or `agent/activation/SKILL.md` (wrong leaf — it's `activate`, not `activation`, in core).

6. **AC7 partial: renderer reads the skill but ignores it.** β must verify the test in `activate_test.go` actually demonstrates skill content drives output — i.e., editing the skill's ordering line in a fixture changes the `cn activate` output. A test that asserts the skill file exists but does not exercise content-driven output passes the file-read oracle while failing the source-of-truth invariant.

7. **AC7 regression on `manifest-only` deps state.** Current behavior: when `.cn/deps.json` declares cnos.core but it's not yet restored, `cn activate` emits an instructive "run cn deps restore" message in the `## Read first` and `## Kernel` blocks. The renderer evolution must preserve this — if the skill is not vendored at `.cn/vendor/packages/cnos.core/skills/agent/activate/SKILL.md`, the renderer needs a sensible fallback or error, not a panic or silent empty output. Watch for: file-read with no error check; nil skill content emitting a broken prompt.

8. **Interpolation surface drift.** The issue body's Implementation guidance section names this explicitly: the renderer's interpolation surface (what hub-state fields substitute into what template positions, e.g., `{hub_persona_path}`, `{hub_latest_reflection}`) must be declared — either in the skill itself or in the renderer's documented contract. Implicit substitution that drifts is a binding finding. Watch for: skill is plain prose with no template positions, and the renderer composes output by mixing skill content with hard-coded Go strings — that re-introduces the hardcoded-ordering problem AC7 displaces.

9. **Observable-output drift for current consumers.** `cn activate` is invoked by `claude -p` and current consumers parse its output (per the issue's active design constraints). Output deltas must be either zero or documented as claims in `self-coherence.md`. Watch for: section header renames (`## Read first` → `## Bootstrap`), block reorderings that break upstream parsers, or removal of the `## Kernel` / `## Persona` blocks.

10. **Non-goal violation.** The issue's Non-goals are firm: do NOT patch any hub README (cn-sigma, cn-pi), do NOT stand up `cnos.xyz`, do NOT migrate other procedures from code to skill form, do NOT touch `cdd/activation/SKILL.md`, do NOT change the `cn activate` CLI surface (flags, args, exit codes), do NOT define a `cn doctor` check. The diff must touch only:
    - `src/packages/cnos.core/skills/agent/activate/SKILL.md` (NEW)
    - `src/go/internal/activate/activate.go` (MODIFY)
    - `src/go/internal/activate/activate_test.go` (MODIFY)
    - `.cdd/unreleased/379/*.md` (cycle evidence — required)
    Plus any vendored-bundle path the renderer reads from, if α needs to update vendored copies for the test to pass. Any other path triggers a binding finding.

## Cross-repo lineage

This cycle accepts a proposal from `usurobor/cn-sigma:.cdd/iterations/cross-repo/cnos/agent-activate-skill/`. The bundle was merged to cn-sigma main at `1a4e25f75` carrying the AC4 capability-matrix refinement; cnos main has a target-side mirror branch `sigma/cross-repo-mirror-agent-activate-skill` @ `212d5239` pending operator merge.

Close-out responsibilities (γ at cycle close):
- Verify the cn-sigma source proposal `STATUS` receives a `landed` event when this cycle merges, or that the branch emits a feedback patch with that event
- Record the landed cnos commit / artifact path in the `landed` event
- Do not let the proposal remain at `accepted` after target work has landed (per gamma/SKILL.md §2.7)

## Branch and sequence

- **Branch:** `cycle/379` (created from `origin/main` at `7a7f7152`, pushed to origin)
- **α target:** write `src/packages/cnos.core/skills/agent/activate/SKILL.md` + evolve `src/go/internal/activate/activate.go` and `activate_test.go` + `self-coherence.md` on `cycle/379`; signal review-readiness
- **β target:** review on `cycle/379`; on APPROVE, merge into `main` per `beta/SKILL.md` §1; write `beta-review.md` + `beta-closeout.md`
- **γ close-out target:** `gamma-closeout.md` after both α and β close-outs land; PRA per `post-release/SKILL.md`

## External gates needed from δ at scaffold time

None. γ's scaffolding (branch creation, push, scaffold artifact, prompt authoring) is complete in this dispatch. δ's next action is to dispatch α with `/tmp/dispatch-379/alpha-prompt.md`.
