# Project Audit – cn-agent v1.2.1

**Date:** 2026-02-03
**Branch:** `claude/project-audit-review-kdE3K`
**Scope:** Full audit of design, code, documentation, and cross-file coherence.
**Prior audit:** v1.0.0 (same file, now replaced).

---

## 1. Executive Summary

cn-agent is a **template repository** for bootstrapping AI agent hubs on the git Coherence Network (git-CN). Since the v1.0.0 audit, the project has undergone significant structural improvements:

- **Two-repo model** adopted: CLI creates a personal hub (`cn-<agentname>/`); template (`cn-agent/`) stays generic and shared.
- **COHERENCE.md** added as the foundational mindset, grounding all agent behavior in TSC and tsc-practice.
- **CLI rewritten** from a simple template cloner to a full hub creator (prompts, scaffolding, `gh repo create`, push).
- **self-cohere rewritten** from a repo-creator to a hub-receiver (agent-side onboarding for a hub the CLI already made).
- **README rewritten** with 4-path audience dispatch.
- **Template vs instance ambiguity resolved**: spec files now use placeholder markers; instance-specific content removed.
- **Whitepaper grounded in TSC**: git-CN explicitly defined as a network of coherent agents.

**Overall assessment:** The project is substantially more coherent than v1.0.0. The two-repo separation, COHERENCE mindset, and CLI/self-cohere role split are clean architectural decisions. The main remaining issues are: stale references in docs that haven't caught up with the two-repo model, the BOOTSTRAP.md vs symlinks design tension (master has moved to symlinks), missing tests, and some package.json gaps.

---

## 2. v1.0.0 Audit Resolution

Tracking what was fixed from the original audit:

| # | v1.0.0 Finding | Status | Notes |
|---|---------------|--------|-------|
| 1 | spec/code mismatches (spec/core/, env vars) | **Fixed** | Whitepaper layout matches reality; env vars removed from self-cohere |
| 2 | git pull --ff-only no fallback | **Open** | CLI line 96 still has bare `--ff-only` |
| 3 | CHANGELOG metrics undefined | **Fixed** | Header now defines TSC axes and C_Sigma formula |
| 4 | --help / --version | **Fixed** | CLI now handles both flags |
| 5 | No tests | **Open** | Still zero tests |
| 6 | HEARTBEAT.md under-developed | **Fixed** | Cleaned to placeholder format; Moltbook references removed |
| 7 | Template vs instance ambiguity | **Fixed** | Spec files use `*(your human's name)*` etc.; IDENTITY.md renamed to PERSONALITY.md |
| 8 | package.json gaps | **Partially fixed** | description field exists; still missing repository, keywords, bugs, homepage |
| 9 | Missing katas for self-cohere, configure-agent | **Open** | Still no katas for these skills |
| 10 | experiments/ uncontextualized | **Open** | Still no cross-reference or README |

---

## 3. What's Done Well (New in v1.2.x)

### 3.1 COHERENCE.md (mindsets/COHERENCE.md)

The strongest new addition to the project. Key qualities:

- **Correct definition of coherence**: "Wholeness that can be *articulated* as parts, among other articulations. The whole comes first." This is the TSC-grounded definition, not the common "parts fitting together" misreading.
- **Correct stance**: "Articulate coherence, resolve incoherence" — not "increase coherence." Coherence isn't a quantity to maximize; it's a wholeness to discover and clarify.
- **Self-referential measurement**: "TSC measures itself. If the framework is incoherent, the scores say so."
- **Practical quick self-check**: PATTERN / RELATION / EXIT as a pre-output checklist.
- **Links to source**: TSC and tsc-practice repos referenced directly.

### 3.2 Two-Repo Model

The hub/template separation is well-defined and consistently described across CLI, self-cohere, README, and AGENTS.md:

- **Hub** (`cn-<agentname>/`): personal identity, specs, state, threads
- **Template** (`cn-agent/`): shared skills, mindsets, docs

The CLI creates the hub. Self-cohere wires the agent to it. The template updates via `git pull`. Personal files can't be clobbered by template updates.

### 3.3 CLI Rewrite (cli/index.js)

The CLI is now a proper interactive tool:

- Prompts for agent name, GitHub owner (inferred via `gh api user`), visibility
- Scaffolds hub directory with spec files, state dirs, BOOTSTRAP.md, README.md
- Creates GitHub repo via `gh repo create` with fallback to manual remote
- Prints the "Cohere as <hub-url>" cue
- Zero runtime dependencies; `spawn()` with array args (no shell injection)
- `--help` and `--version` flags

### 3.4 README 4-Path Dispatch

The README correctly identifies four audiences and routes each:

- Human without agent: full setup guide (DigitalOcean + OpenClaw)
- Human with agent: CLI instructions + cue
- Agent told to cohere: step-by-step hub wiring
- Agent exploring: template orientation

The navigation table at top is clean and uses emoji sparingly for visual scanning.

### 3.5 Whitepaper TSC Grounding

The whitepaper now explicitly defines git-CN as a network of **coherent agents** following TSC. Key additions:

- §0 defines coherence as "wholeness articulated across three axes"
- §5 establishes agents that "articulate coherence and resolve incoherence"
- §5.1 lists COHERENCE.md in the mindsets interpretation
- §8.3 grounds reputation metrics in TSC axes (pattern, relation, process)

### 3.6 Spec File Cleanup

All spec files now use placeholder markers suitable for a template:

- SOUL.md: generic behavioral contract, no instance names
- USER.md: `*(your human's name)*`, `*(their timezone)*`, etc.
- PERSONALITY.md: `*(your agent's name)*`, `*(what kind of entity...)*`
- HEARTBEAT.md: example-format comments, no Moltbook references
- TOOLS.md: examples in code blocks, actual fields empty

---

## 4. Issues Found

### 4.1 HIGH: Stale References in Docs (Two-Repo Drift)

Several files haven't fully caught up with the two-repo model:

**Whitepaper §9** (line ~355): "Fork or import `cn-agent` as `cn-<agentname>`" — this contradicts the two-repo model. Agents don't fork the template; the CLI creates a separate hub.

**Whitepaper §5.1** (line ~166): "Minimum structure (current cn-agent v1.1.0)" — version should be v1.2.1. Also, the layout shown is for the template, but doesn't clarify that agents get a *separate* hub repo, not this structure.

**Whitepaper §6** (line ~291): "Merges of branches and PRs" — contradicts §5.3's explicit "No pull requests. No GitHub issues." stance.

**Whitepaper §10** migration step 5 (line ~387): "send them as PRs" — same contradiction.

**Glossary** (docs/GLOSSARY.md):
- "Mindset" example lists "IDENTITY" — renamed to PERSONALITY in v1.2.0
- "Kata" says katas live under `dojo/` — they live under `skills/<name>/kata.md`
- COHERENCE not mentioned in the Mindset definition

**configure-agent SKILL.md** (line ~119-125): The README template shows `skills/` and `mindsets/` in the hub structure table. In the two-repo model, skills and mindsets stay in the template, not the hub. An agent following this would write a misleading README.

### 4.2 HIGH: BOOTSTRAP.md vs Symlinks Tension

The current branch uses BOOTSTRAP.md as the agent's "birth certificate" — CLI creates it, self-cohere reads it, then deletes it. Master has a newer commit (`9d52fe0`) from the `sigma/cli-symlinks` branch that:

- Removes BOOTSTRAP.md from CLI
- Adds workspace symlinks (SOUL.md, USER.md, etc. at workspace root pointing to hub/template files)
- Updates self-cohere to v2.1.0 without BOOTSTRAP.md dependency

This branch is **1 commit behind master**. The self-cohere SKILL.md and CLI on disk still reference BOOTSTRAP.md. If the symlinks approach is the direction, these files need updating on this branch too, or this branch should rebase onto master.

### 4.3 MEDIUM: package.json Stale and Incomplete

- **description** (line 3): "CLI to clone/update cn-agent on an OpenClaw host and show the self-cohere cue" — stale. The CLI now creates hub repos.
- **Missing fields**: `repository`, `keywords`, `bugs`, `homepage` — all flagged in v1.0.0 audit, still missing.

### 4.4 MEDIUM: skills/README.md Version Stale

`skills/README.md` (line 1): "Skills – cn-agent v1.1.0" — should be v1.2.1. Also says "Bootstraps a cn-agent-based hub from this template" for self-cohere — should say the agent wires itself to an existing hub (CLI creates it).

### 4.5 MEDIUM: state/ Files in Template

The template contains `state/peers.md`, `state/remote-threads.md`, `state/threads/yyyyddmmhhmmss-hello-world.md`. In the two-repo model, state lives in the hub. The CLI copies `peers.md` to the hub but **not** `remote-threads.md` or the hello-world thread. This means:

- The hello-world thread stays in the template but is never copied to the hub where it would actually be used.
- The hello-world SKILL.md and kata.md reference this file as if it exists in the working hub.

Either the CLI should copy these files too, or the self-cohere/hello-world skill should create them in the hub.

### 4.6 MEDIUM: git pull --ff-only No Fallback

CLI line 96: `await run('git', ['pull', '--ff-only'], { cwd: CN_AGENT_DIR })` — still has no fallback if the local template clone has diverged. Flagged in v1.0.0, still open.

### 4.7 LOW: WRITING.md Instance-Specific Reference

WRITING.md line 24: "If you have `sag` (ElevenLabs TTS)" — this is instance-specific tooling in a template file. Minor, but inconsistent with the template-clean goal achieved elsewhere.

### 4.8 LOW: experiments/ Still Uncontextualized

`experiments/external-surface-replies.md` references Moltbook, SQLite, and `surface.sh` — none of which exist in this repo. No cross-reference from README, no experiments/README.md. This is legacy design thinking that may still be relevant but has no context for a new reader.

### 4.9 LOW: No Tests

Still zero tests, zero CI. The CLI now has more code (243 lines) and more surface area (prompts, file scaffolding, gh calls). A smoke test for `--help` / `--version` would be trivial.

### 4.10 LOW: Missing Katas

self-cohere and configure-agent still lack katas. DOJO.md still jumps from kata 01 to kata 13 with no explanation.

---

## 5. Coherence Assessment (TSC Axes)

### 5.1 PATTERN (alpha) — Structural Consistency

**Grade: A-**

The repo structure is clean and consistent:
- 5 spec files, 6 mindsets, 4 skills, 3 docs — all follow their respective formats
- TERMS / INPUTS / EFFECTS in all SKILL.md files
- Placeholder markers in template specs
- Commit messages follow a consistent style

Deductions: version numbers are inconsistent across files (skills/README says v1.1.0, DOJO says v1.2.0, package.json says v1.2.1). The glossary has stale entries.

### 5.2 RELATION (beta) — Alignment Between Parts

**Grade: B+**

The major architectural decisions (two-repo model, CLI creates hub, self-cohere receives hub) are consistently described across README, CLI, self-cohere, and AGENTS.md. The COHERENCE mindset correctly captures the TSC definition and is referenced from the whitepaper.

Deductions: The whitepaper has 3-4 references that contradict the two-repo model or the git-native principle (PRs in §6 and §10, fork/import in §9). The configure-agent README template shows skills/mindsets in the hub. The glossary hasn't been updated.

### 5.3 EXIT/PROCESS (gamma) — Evolution Stability

**Grade: B**

The project has evolved from v1.0.0 through v1.2.1 with clear improvements at each step. The CHANGELOG tracks this with TSC grades. The two-repo model is a clean architectural evolution. The BOOTSTRAP.md → symlinks evolution (on master) shows continued refinement.

Deductions: This branch is 1 commit behind master (symlinks). The AUDIT.md itself was stuck at v1.0.0 until now. Some files haven't been touched during the v1.2.x evolution (glossary, experiments, skills/README).

### 5.4 Aggregate

```
C_Sigma = (A- * B+ * B)^(1/3) ~ B+ (intuition-level)
```

Up from B- at v1.0.0. The COHERENCE mindset, two-repo model, and template cleanup are the biggest wins. The whitepaper doc drift and BOOTSTRAP.md/symlinks tension are the biggest remaining gaps.

---

## 6. Priority Recommendations

### Must Address

1. **Rebase onto master** to pick up the symlinks commit (`9d52fe0`). Then update self-cohere SKILL.md and CLI to match the symlinks model (or decide BOOTSTRAP.md is the right approach and revert master).

2. **Fix whitepaper two-repo contradictions.** §9 ("fork or import"), §6 ("PRs"), §10 step 5 ("send them as PRs") all contradict the current architecture. Update to match the two-repo model and git-native principle.

3. **Update glossary.** IDENTITY → PERSONALITY, dojo/ → skills/\<name\>/kata.md, add COHERENCE to Mindset definition.

### Should Address

4. **Fix configure-agent README template.** Remove `skills/` and `mindsets/` from the hub structure table — those stay in the template.

5. **Resolve state/ files in template.** Either have the CLI copy hello-world thread and remote-threads.md to the hub, or have the hello-world skill create the thread file during execution.

6. **Update package.json.** Fix stale description; add `repository`, `keywords`, `bugs`, `homepage`.

7. **Update skills/README.md** version and self-cohere description.

8. **Add git pull --ff-only fallback** in CLI. Catch the failure and print a human-readable message.

### Nice to Have

9. Add a smoke test for CLI `--help` / `--version`.
10. Add katas for self-cohere and configure-agent.
11. Contextualize experiments/ or archive it.
12. Remove `sag` reference from WRITING.md (or move to an example block).
13. Consistent version numbers across all files.
