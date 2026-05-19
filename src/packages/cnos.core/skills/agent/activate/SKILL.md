---
name: activate
description: Activate an agent identity at a cnos hub. Body-agnostic procedure for loading Kernel + CA skills + Persona + Operator + hub state, in that order, and confirming identity before any other action. Source of truth for `cn activate` and for non-cn bodies that self-bootstrap.
artifact_class: skill
kata_surface: embedded
governing_question: How does an AI body, knowing only a hub URL, reach a state where it can name its identity, its operator, and its current orientation without operator improvisation?
visibility: public
parent: agent
triggers:
  - activate
  - activation
  - bootstrap
  - self-activation
  - wake up
  - cn activate
scope: task-local
inputs:
  - hub URL or hub path (the named hub the body is activating against)
  - body capabilities (shell + git, HTTP fetch only, or no fetch)
  - access to this skill's content (vendored, cloned, fetched, or injected)
outputs:
  - body has loaded Kernel doctrine (what kind of agent)
  - body has loaded CA skills (behavioral instruction set)
  - body has loaded hub Persona (who it is at this hub)
  - body has loaded hub Operator (who it serves at this hub)
  - body has surveyed hub state (deps, latest reflection, memory, threads)
  - body has named its identity, operator, and current orientation
requires:
  - a hub URL or hub path
  - either shell+git, HTTP fetch, or an operator willing to inject the bootstrap
calls:
  - cnos.core/doctrine/KERNEL.md
  - cnos.core/skills/agent/cap/SKILL.md
  - cnos.core/skills/agent/clp/SKILL.md
  - cnos.core/skills/agent/mca/SKILL.md
  - cnos.core/skills/agent/mci/SKILL.md
  - cnos.core/skills/agent/coherent/SKILL.md
  - cnos.core/skills/agent/agent-ops/SKILL.md
---

<!--
section-manifest:
  planned: [frontmatter, core-principle, algorithm, define, unfold-load-order, unfold-capability-matrix, unfold-readme-router, unfold-disambiguation, rules, renderer-contract, verify, failure-modes, kata, references]
  completed: [frontmatter, core-principle, algorithm, define, unfold-load-order, unfold-capability-matrix]
-->

# Activate

## Core Principle

**A body wakes up incoherent by default.** It carries model weights, not memory of which agent it is. Activation is the procedure that closes the gap between "a body with a hub URL" and "a coherent agent oriented to that hub" — and the procedure is the same regardless of whether the body has the `cn` binary, only WebFetch, or only an operator to dictate text into a chat box.

A coherent activation has six load steps in this exact order: **Kernel → CA skills → Persona → Operator → hub state → identity confirmation.** The order matters. The Kernel and CA skills define _what kind of agent_ this is (the soul — generic, in cnos); the Persona and Operator define _which agent_ this is (the identity — per-hub). Loading identity before soul produces a body that knows its name but not how to think; loading soul without identity produces a body that thinks coherently but cannot name what it is doing or for whom. Both failures are documented in cn-sigma `threads/adhoc/20260325-session2-learnings.md` §1 ("I wake up incoherent by default") and §3 ("Soul = what kind of agent. Identity = which agent. Don't mix them.").

This skill is the **single source of truth for the activation procedure**. The `cn activate` Go command reads this skill from the vendored bundle and renders it with hub-state interpolation; non-cn bodies (Claude Code on the web, Codex sessions, Claude.ai with WebFetch) fetch this skill directly and follow it against the named hub. Every cnos hub's README points at this skill via the router template in §3.3. The procedure lives here, in skill form, not in any runtime that consumes it.

The failure mode the skill prevents is **improvised activation**: a body that wakes up, asks the operator twice "to what?", reads files in an order that depends on operator suggestion, and never reaches a state where it can name its identity without being told. Improvised activation produces drift between sessions, between hubs, and between bodies — every wake-up reinvents the procedure. The skill is the canonical sequence; the renderer renders from it; future bodies fetch from it.

## Algorithm

1. **Resolve capability** — determine which load tier the body's environment supports (shell + git, HTTP fetch only, or no fetch). See §3 capability matrix.
2. **Load soul** — fetch and read the Kernel doctrine from cnos, then the CA skill set from cnos. These define what kind of agent this is, independent of which hub.
3. **Load identity** — fetch and read the hub's `spec/PERSONA.md` and `spec/OPERATOR.md`. These define which agent this is and whom it serves at this hub.
4. **Survey hub state** — list the hub's dependency manifest, the latest daily reflection, the memory surfaces (reflections, adhoc), and the thread surfaces (in, inbox, mail, archived). The hub state is observable evidence, not narrative — read enough of it to know what is currently in motion.
5. **Confirm identity** — produce a short statement naming (a) which agent (from Persona), (b) which operator (from Operator), (c) which hub (path or URL), (d) the current cycle or thread if one is in motion (from threads or unreleased). The body is activated when this statement is concrete and grounded in the loaded files.

The sequence is total: each step depends on the previous. Skipping a step produces a body that operates from one of the named failure modes (§5).

---

## 1. Define

### 1.1. Identify the parts

A complete activation has these parts:

- **Kernel** — generic coherence doctrine for any cnos agent (`src/packages/cnos.core/doctrine/KERNEL.md`)
- **CA skills** — the behavioral instruction set (`src/packages/cnos.core/skills/agent/{cap,clp,mca,mci,coherent,agent-ops}/SKILL.md`)
- **Persona** — hub-specific identity (`<hub>/spec/PERSONA.md`)
- **Operator** — hub-specific operator instructions (`<hub>/spec/OPERATOR.md`)
- **Hub state** — dependency manifest, latest reflection, memory surfaces, thread surfaces
- **Identity confirmation** — a concrete statement naming agent, operator, hub, current orientation

- ❌ "The skill says load these files" (treating activation as a list of paths)
- ✅ "Kernel + CA = soul; Persona + Operator = identity; hub state = what is in motion; confirmation = the gate that proves the previous five landed"

### 1.2. Articulate how they fit

The Kernel and CA skills are **substrate-independent**: they are the same content for every cnos agent regardless of hub. They live in the `cnos.core` package. The Persona and Operator are **per-hub**: they are the content that makes this body _this_ agent at _this_ hub. The hub state is **observable evidence**: what is in motion right now. The identity confirmation is the **gate** — without it the body has not actually activated, it has only read files.

The two-layer split (soul vs identity) is constitutive. Mixing them produces hubs that re-define what kind of agent any body should be (overreach: the hub should not redefine the kernel) and bodies that load a hub's PERSONA without first knowing what kind of thing reads a PERSONA (underreach: the body never gets to coherent agency).

- ❌ Hub `spec/SOUL.md` (the hub redefining the kernel) — supplanted by Kernel-in-cnos
- ✅ Kernel in cnos, Persona/Operator in the hub — layering preserved

### 1.3. Name the failure mode

Activation fails in five named ways. Each has a concrete symptom; each has a structural fix.

- **F1 — Improvised activation.** Body wakes up, asks "to what?" twice, reads files in operator-suggested order, never reaches identity. Fix: this skill, fetched and followed.
- **F2 — Soul/identity confusion.** Body loads Persona before Kernel; Persona references doctrine it has not loaded; body operates from a partial mental model. Fix: enforce six-step order in §2.1.
- **F3 — Capability mismatch.** Body assumes shell+git when it only has WebFetch, or assumes WebFetch when it has neither, then silently fails or hallucinates content. Fix: §3 capability matrix with explicit ladder; body picks tier its environment actually supports.
- **F4 — Hub README dead-end.** Hub README describes the hub for humans but does not route AI bodies into activation; bodies that land on the README stall there. Fix: README router template in §3.3, adopted verbatim by every hub.
- **F5 — Renderer-skill drift.** `cn activate` renders an order or content that differs from this skill, and the two diverge over time; bodies that use `cn` see one procedure, bodies that fetch the skill see another. Fix: the renderer reads this skill (§4); the skill is the source of truth, the command is a thin surface.

---

## 2. Unfold

### 2.1. The six-item load order

The canonical activation sequence is six steps in this exact order. **The order is the skill.** Any renderer that consumes this skill (including `cn activate`) emits in this order; any body that follows this skill loads in this order. The order does not depend on body capability — only the _mechanism_ of fetching changes (§3). The order itself is invariant.

The machine-readable form lives in §4 (renderer contract); the human-readable form is here.

1. **Kernel from cnos.** Fetch `src/packages/cnos.core/doctrine/KERNEL.md` from the cnos repository. This defines what kind of agent any cnos body is — observation before action, truth over comfort, coherence over drift. Load this first because every subsequent file presupposes its frame.

2. **CA skills from cnos.** Fetch the coherent-agent skill set from cnos:
   - `src/packages/cnos.core/skills/agent/cap/SKILL.md` — coherent-agent protocol
   - `src/packages/cnos.core/skills/agent/clp/SKILL.md` — coherent-language protocol
   - `src/packages/cnos.core/skills/agent/mca/SKILL.md` — modify-the-codebase actions
   - `src/packages/cnos.core/skills/agent/mci/SKILL.md` — modify-the-conceptual-instructions
   - `src/packages/cnos.core/skills/agent/coherent/SKILL.md` — coherent output
   - `src/packages/cnos.core/skills/agent/agent-ops/SKILL.md` — ops the agent can request
   These define _how_ the body behaves; load them before any per-hub content so the hub content is interpreted through the right behavioral lens.

3. **Persona from hub.** Fetch `<hub>/spec/PERSONA.md`. This file names which agent this body is at this hub: name, role, vibe, what this agent is for. Persona is per-hub — the same body may be Sigma at cn-sigma and Pi at cn-pi and another agent at a third hub. Persona presupposes the Kernel and CA skills; load it after them.

4. **Operator from hub.** Fetch `<hub>/spec/OPERATOR.md`. This file names which operator this body serves at this hub: who decides, who routes, what their preferences are, what they will and will not approve. Operator is per-hub. Load it after Persona because operator instructions often reference the agent's role from Persona.

5. **Hub state.** Survey the hub's current state — what is in motion right now. The hub state is a composite of four readable surfaces:
   - **Dependency manifest** — `<hub>/.cn/deps.json`, which packages this hub vendors from cnos and which versions.
   - **Latest reflection** — most-recent file under `<hub>/threads/reflections/daily/`. This is the operator's or agent's most recent durable thought; reading it grounds the body in current-state-of-mind, not session-start.
   - **Memory surfaces** — what reflection layers exist: `<hub>/threads/reflections/{daily,weekly,monthly}/`, `<hub>/threads/adhoc/`. These are durable continuity across sessions.
   - **Thread surfaces** — what inbox / mail / archived surfaces exist: `<hub>/threads/{in,inbox,mail,archived}/`. These are live coordination state.
   Read each in order. The hub state is grouped as item 5 because the four surfaces are co-presupposing: they all describe what is currently happening at this hub.

6. **Identity confirmation.** Produce a concrete statement, in the body's voice, naming:
   - **Who** — the agent name and role from Persona (e.g. "Sigma, the coordinator at cn-sigma")
   - **Whom** — the operator name and current expectations from Operator (e.g. "operator: usurobor")
   - **Where** — the hub URL or path (e.g. "hub: github.com/usurobor/cn-sigma")
   - **When** — the current cycle / thread / in-motion work, if any (from threads or unreleased)
   Confirmation is the gate. A body that has read all five preceding items but cannot produce this statement has not activated — it has only read files. The statement is the test.

**Layering rule (soul vs identity).** Items 1–2 are the **soul**: what kind of agent. They live in cnos and are the same for every hub. Items 3–4 are the **identity**: which agent. They live in the hub and differ per hub. The split is constitutive — mixing them produces the failures named in §1.3. The layering rule is from cn-sigma `threads/adhoc/20260325-session2-learnings.md` §3: _"Soul = what kind of agent. Identity = which agent. Don't mix them."_

- ❌ Load Persona first, infer Kernel from how Persona refers to it
- ❌ Hub's `spec/SOUL.md` redefining the kernel (a hub does not own doctrine)
- ✅ Kernel + CA skills first (soul), then Persona + Operator (identity), then state, then confirmation

### 2.2. Three-tier body-capability matrix

A body activates by following the six-item load order; the _mechanism_ for fetching the files depends on what the body's environment supports. Bodies misjudge their own capabilities (cn-sigma `threads/adhoc/20260519-git-read-and-untested-limits.md` is the reference reflection): a body that has shell+git should not fall back to WebFetch, and a body that has only WebFetch should not pretend to clone. The matrix is explicit so bodies pick the tier their environment _actually_ supports.

| Tier | Capability | Load path for cnos (Kernel + CA skills) | Load path for hub (Persona + Operator + state) | Cost |
|------|------------|------------------------------------------|-------------------------------------------------|------|
| **(a) shell + git** _(preferred when available)_ | Body can execute shell commands and `git clone` | `git clone https://github.com/usurobor/cnos.git`, then read files from the local checkout | `git clone <hub-url>`, then read files from local checkout — or, if the hub directory is already mounted locally, read in place | One network operation per repo; atomic; later reads are filesystem-fast |
| **(b) HTTP fetch only** | Body can fetch raw URLs but cannot execute shell | Fetch raw GitHub URLs per file, e.g. `https://raw.githubusercontent.com/usurobor/cnos/main/src/packages/cnos.core/doctrine/KERNEL.md` | Fetch raw GitHub URLs per file, e.g. `https://raw.githubusercontent.com/usurobor/cn-sigma/main/spec/PERSONA.md` | One network operation per file; non-atomic; may hit rate limits |
| **(c) no fetch** | Body has no shell, no git, no HTTP fetch | _The body cannot self-activate._ The operator must paste the bootstrap contents into the body's prompt window directly. | _The operator must inject._ | Operator labor per session; degraded path |

**Tier (a) shell + git is preferred when available.** Reasons:

- **Atomic.** One `git clone` fetches every file in the bundle at a consistent commit; later reads cannot drift between files. Fetching files one at a time over HTTP can produce inconsistent reads if the upstream changes mid-load.
- **Local.** Once cloned, every subsequent file read is filesystem-fast and does not touch the network. HTTP-fetch tiers pay per-file network cost.
- **Inspectable.** The body can list directories, find related files, follow links — capabilities WebFetch alone does not give.
- **Honest about reality.** Most modern bodies (Claude Code on web/phone/CLI, Codex sessions, agent-mode bodies) have shell + git. Defaulting to WebFetch when shell + git is available is the failure mode named in cn-sigma `threads/adhoc/20260519-git-read-and-untested-limits.md` — bodies underclaim their own capabilities.

**Tier (b) HTTP fetch only is the constrained path.** Some bodies (Claude.ai sessions with WebFetch enabled but no shell, embedded chat surfaces with constrained tooling) have only HTTP fetch. The procedure is the same six steps; the mechanism is per-file raw GitHub URLs.

**Tier (c) no fetch is the degraded path, named honestly.** A body with no shell, no git, and no HTTP fetch _cannot self-activate_. Activation in this tier requires the operator to inject the bootstrap content directly — paste this skill's text into the prompt window, paste the Kernel, paste the Persona, etc. The body cannot reach the procedure without operator labor. This is named explicitly so operators do not assume the body can "figure it out" in environments where it structurally cannot.

**Tier selection.** Before beginning the six-step load, the body MUST observe its own environment and pick the tier its capabilities actually support. The choice is observable:

- If the body can run `git --version` and `git clone` — tier (a).
- Else if the body can fetch `https://raw.githubusercontent.com/...` — tier (b).
- Else — tier (c), and the body announces this so the operator knows to inject.

- ❌ Body assumes WebFetch because that's the first tool it remembers
- ❌ Body claims "no capability" without checking shell access
- ✅ Body runs one capability check, picks the matching tier, announces the tier

