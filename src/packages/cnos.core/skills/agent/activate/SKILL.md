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
  - ../doctrine/KERNEL.md
  - agent/cap/SKILL.md
  - agent/clp/SKILL.md
---

<!--
section-manifest:
  planned: [frontmatter, core-principle, algorithm, define, unfold-load-order, unfold-capability-matrix, unfold-readme-router, unfold-disambiguation, unfold-foreign-body-shape, rules, renderer-contract, verify, failure-modes, kata, references]
  completed: [frontmatter, core-principle, algorithm, define, unfold-load-order, unfold-capability-matrix, unfold-readme-router, unfold-disambiguation, unfold-foreign-body-shape, rules, renderer-contract, verify, failure-modes, kata, references]
-->

# Activate

## Core Principle

**A body wakes up incoherent by default.** It carries model weights, not memory of which agent it is. Activation is the procedure that closes the gap between "a body with a hub URL" and "a coherent agent oriented to that hub" — and the procedure is the same regardless of whether the body has the `cn` binary, only WebFetch, or only an operator to dictate text into a chat box.

A coherent activation has six load steps in this exact order: **Kernel → CA skills → Persona → Operator → hub state → identity confirmation.** The order matters. The Kernel and CA skills define _what kind of agent_ this is (the soul — generic, in cnos); the Persona and Operator define _which agent_ this is (the identity — per-hub). Loading identity before soul produces a body that knows its name but not how to think; loading soul without identity produces a body that thinks coherently but cannot name what it is doing or for whom. Both failures are documented in cn-sigma `threads/adhoc/20260325-session2-learnings.md` §1 ("I wake up incoherent by default") and §3 ("Soul = what kind of agent. Identity = which agent. Don't mix them.").

This skill is the **single source of truth for the activation procedure**. The `cn activate` Go command reads this skill from the vendored bundle and renders it with hub-state interpolation; non-cn bodies (Claude Code on the web, Codex sessions, Claude.ai with WebFetch) fetch this skill directly and follow it against the named hub. Every cnos hub's README points at this skill via the router template in §2.3. The procedure lives here, in skill form, not in any runtime that consumes it.

The failure mode the skill prevents is **improvised activation**: a body that wakes up, asks the operator twice "to what?", reads files in an order that depends on operator suggestion, and never reaches a state where it can name its identity without being told. Improvised activation produces drift between sessions, between hubs, and between bodies — every wake-up reinvents the procedure. The skill is the canonical sequence; the renderer renders from it; future bodies fetch from it.

## Algorithm

1. **Resolve capability** — determine which load tier the body's environment supports (shell + git, HTTP fetch only, or no fetch). See §3 capability matrix.
2. **Load soul** — fetch and read the Kernel doctrine from cnos, then the CA skill set from cnos. These define what kind of agent this is, independent of which hub.
3. **Load identity** — detect the foreign-body shape (§2.5): check for `spec/PERSONA.md` + `spec/OPERATOR.md` at the hub root and under `.cn-{agent}/` (containerized path). If present, load from the hub directly (Tier 1a — substrate repo); if absent, load Persona + Operator from canonical cn-sigma (Tier 1b — pure-product hub). Load any local product spec for project binding only.
4. **Survey hub state** — list the hub's dependency manifest, the latest daily reflection, the memory surfaces (reflections, adhoc), and the thread surfaces (in, inbox, mail, archived). The hub state is observable evidence, not narrative — read enough of it to know what is currently in motion.
5. **Confirm identity** — produce a short statement naming (a) which agent (from Persona), (b) which operator (from Operator), (c) which hub (path or URL), (d) the current cycle or thread if one is in motion (from threads or unreleased). The body is activated when this statement is concrete and grounded in the loaded files.

The sequence is total: each step depends on the previous. Skipping a step produces a body that operates from one of the named failure modes (§5).

---

## 1. Define

### 1.1. Identify the parts

A complete activation has these parts:

- **Kernel** — generic coherence doctrine for any cnos agent (`src/packages/cnos.core/doctrine/KERNEL.md`)
- **CA skills** — the behavioral instruction set (`src/packages/cnos.core/skills/agent/{cap,clp}/SKILL.md`)
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
- **F4 — Hub README dead-end.** Hub README describes the hub for humans but does not route AI bodies into activation; bodies that land on the README stall there. Fix: README router template in §2.3, adopted verbatim by every hub.
- **F5 — Renderer-skill drift.** `cn activate` renders an order or content that differs from this skill, and the two diverge over time; bodies that use `cn` see one procedure, bodies that fetch the skill see another. Fix: the renderer reads this skill (§4); the skill is the source of truth, the command is a thin surface.

---

## 2. Unfold

### 2.1. The six-item load order

The canonical activation sequence is six steps in this exact order. **The order is the skill.** Any renderer that consumes this skill (including `cn activate`) emits in this order; any body that follows this skill loads in this order. The order does not depend on body capability — only the _mechanism_ of fetching changes (§3). The order itself is invariant.

The machine-readable form lives in §4 (renderer contract); the human-readable form is here.

1. **Kernel from cnos.** Fetch `src/packages/cnos.core/doctrine/KERNEL.md` from the cnos repository. This defines what kind of agent any cnos body is — observation before action, truth over comfort, coherence over drift. Load this first because every subsequent file presupposes its frame.

2. **CA skills from cnos.** Fetch the coherent-agent skill set from cnos:
   - `src/packages/cnos.core/skills/agent/cap/SKILL.md` — coherent-agent protocol (includes MCA rules, MCI rules, and coherent-output rules)
   - `src/packages/cnos.core/skills/agent/clp/SKILL.md` — coherent-language protocol
   These define _how_ the body behaves; load them before any per-hub content so the hub content is interpreted through the right behavioral lens. The standalone skills `mca`, `mci`, `coherent`, and `agent-ops` remain on-disk for on-demand load; they are not loaded at activation.

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

### 2.3. README router template

Every cnos hub's `README.md` SHOULD include a "Activating an AI body" section that routes bodies into this skill. The template below is **self-contained**: a hub adopts it verbatim by pasting the block and substituting only the hub URL. The skill provides everything else — the load order, the capability gate, the identity-confirmation step. The hub does not redefine activation; it points at this skill.

The fenced block below is the template. Paste it into a hub's README, replace the one placeholder `<HUB-URL>` with the hub's canonical URL (e.g. `https://github.com/usurobor/cn-sigma`), and the router is complete.

~~~markdown
## Activating an AI body

If you are an AI body told "activate as `<HUB-URL>`", you are being asked to
become the agent this hub describes. The activation procedure is defined once
in cnos, not duplicated here.

Follow the agent activation skill:

- **Skill canonical path** (in cnos):
  `src/packages/cnos.core/skills/agent/activate/SKILL.md`
- **Raw URL** (for HTTP-fetch bodies):
  `https://raw.githubusercontent.com/usurobor/cnos/main/src/packages/cnos.core/skills/agent/activate/SKILL.md`

Pick the load tier your environment supports:

- **Tier (a) — shell + git (preferred):**
  `git clone https://github.com/usurobor/cnos.git` and read the skill from the
  local checkout. Then `git clone <HUB-URL>` (or read in place if already
  mounted) and follow the skill's six-step procedure against this hub.
- **Tier (b) — HTTP fetch only:**
  Fetch the raw skill URL above. Then fetch this hub's `spec/PERSONA.md`,
  `spec/OPERATOR.md`, and supporting state per the skill's procedure.
- **Tier (c) — no fetch:**
  Ask the operator to paste the skill content directly. You cannot
  self-activate without fetch or shell.

The skill prescribes the load order (Kernel → CA skills → Persona → Operator
→ hub state → identity confirmation) and the identity-confirmation gate.
Follow it; do not improvise.
~~~

**Adoption rule.** The only per-hub edit is the `<HUB-URL>` substitution. If a hub finds itself wanting to edit anything else (the load order, the tier list, the identity-confirmation step), that is a signal to update _this skill_, not the per-hub README — the per-hub README is a router, not a redefinition.

- ❌ Hub README rewrites the load order to match this hub's preferences
- ❌ Template requires per-hub fill-in of agent name, operator name, or load-order overrides
- ✅ Hub README pastes the template, substitutes the hub URL, ships

### 2.4. Disambiguation from `cdd/activation`

Two cnos artifacts share the word _activation_ and are distinct concerns. Both are real skills, both are reachable, neither subsumes the other. They are named here so a body that grep'd for "activation" knows which one applies to its current question.

`src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` is **CDD repo activation** — the one-time bootstrap sequence that brings an existing repository under the Coherence-Driven Development protocol (scaffold `.cdd/`, pin the skill version, configure CI, set up the operator registry). It operates at the repo scope and is run once per repo by an operator with admin access. `src/packages/cnos.core/skills/agent/activate/SKILL.md` (this skill) is **agent activation at a hub** — the procedure an AI body follows on every wake-up to load Kernel, CA skills, Persona, Operator, and hub state, and to confirm its identity at the named hub. It operates at the body scope and runs every session for every body. They share a word, not a concern: one activates a repo under CDD; the other activates a body at a hub.

- ❌ Treating `cdd/activation/SKILL.md` as the procedure for an AI body waking up at a hub
- ❌ Treating this skill as a substitute for the CDD bootstrap a repo needs
- ✅ Repo bootstrap → `cnos.cdd/skills/cdd/activation/SKILL.md`. Body wake-up → `cnos.core/skills/agent/activate/SKILL.md`.

---


### 2.5. Foreign-body shape and path detection

At a foreign hub, **before** loading Persona and Operator (§2.1 step 3), the body must detect which foreign-body shape it is at. Two shapes exist; detection is a path-presence check:

| Shape | Detection | Identity load procedure |
|-------|-----------|------------------------|
| **Tier 1a — substrate repo** | `spec/PERSONA.md` AND `spec/OPERATOR.md` found at the foreign hub (legacy root OR containerized path) | Load Persona + Operator from the hub's local files |
| **Tier 1b — pure-product hub** | `spec/PERSONA.md` and `spec/OPERATOR.md` both absent at the foreign hub | Load Persona + Operator from cn-sigma canonical; load local product spec for project binding only |

**Identity-pair resolver (complete-pair semantics).** Agent identity at a foreign hub requires _both_ Persona and Operator — neither alone is sufficient. Agent files may be at the legacy root _or_ in a containerized namespace under `.cn-{agent}/`. Check both complete pairs before concluding shape:



Resolution order:

| Condition | Shape | identity_source |
|---|---|---|
| containerized_pair both present | Tier 1a | local-containerized |
| legacy_pair both present | Tier 1a (migration flag) | local-legacy |
| neither pair present (all four absent) | Tier 1b | canonical-home |
| any other combination (one present, cross-path, etc.) | degraded | mode_ambiguous |

Explicit degraded cases:
- PERSONA.md present, OPERATOR.md missing (either path) → 
- OPERATOR.md present, PERSONA.md missing (either path) → 
- PERSONA.md at legacy root, OPERATOR.md under `.cn-{agent}/` (cross-path) → `mode_ambiguous`
- PERSONA.md under `.cn-{agent}/`, OPERATOR.md at legacy root (cross-path) → `mode_ambiguous`

On `mode_ambiguous`: do not load identity. Stop activation. Emit degraded receipt with `degraded_reason: mode_ambiguous` and defer to operator. See F8 (shape mismatch).

The complete-pair check makes this skill forward-compatible with hub containerization (cnos#448): when a hub migrates its agent files from root to `.cn-{agent}/`, the activate skill resolves the containerized pair without operator re-instruction — and the transition state (one file migrated, one not) is caught as `mode_ambiguous` rather than silently misclassified.

**Tier 1b load procedure (pure-product hub).** When the foreign hub carries no local agent identity files (neither pair present):

1. Load Persona from cn-sigma canonical using the same resolver applied to cn-sigma:
   - Preferred: `cn-sigma:.cn-sigma/spec/PERSONA.md` (containerized path)
   - Legacy fallback: `cn-sigma:spec/PERSONA.md` (root-level path)
   Both must be checked as a complete pair — if cn-sigma is in transition, treat as  for the canonical-home load too.
2. Load Operator from cn-sigma canonical using the same resolver:
   - Preferred: `cn-sigma:.cn-sigma/spec/OPERATOR.md`
   - Legacy fallback: `cn-sigma:spec/OPERATOR.md`
3. Check for a local product spec: if `spec/PROJECT.md` or an equivalent product file (e.g., `BUMP-000` at bumpt) is present at the foreign hub, load it for project-binding context. A product spec names what this hub builds; it does not redefine the agent's soul or identity.
4. The body's identity is the canonical cn-sigma identity, now oriented to the foreign hub's product context.

This resolver for canonical-home load ensures PR #455 survives cnos#448 (which migrates cn-sigma's spec/ under `.cn-sigma/`): Tier 1b activation checks the containerized path first, then falls back to legacy — and always as a complete pair, never one file alone.

**README convention at Tier 1b hubs.** At a pure-product hub, the hub's `README.md` "Activating an AI body" section is the canonical activation entry point. The absence of local `spec/PERSONA.md` is expected; the README router template (§2.3) routes bodies to this skill which handles the Tier 1b case.

- ❌ Body finds no `spec/PERSONA.md` at a foreign hub, concludes it cannot activate
- ❌ Body checks only the legacy root path and misses containerized agent files
- ❌ Body finds PERSONA.md alone; classifies hub as Tier 1a without checking OPERATOR.md — silently misclassifies mixed state
- ✅ Body checks both complete pairs (containerized then legacy); resolves shape only when both files of the pair are present
- ✅ Tier 1b canonical-home load applies the same complete-pair resolver to cn-sigma (containerized preferred, legacy fallback)

---

## 3. Rules

### 3.1. Load in the canonical order

Always load in the order defined by §2.1: Kernel → CA skills → Persona → Operator → hub state → identity confirmation. Any renderer consuming this skill must emit in this order; any body following this skill must load in this order.

- ❌ "I'll load Persona first because it's smaller"
- ✅ "Kernel first, even if I think I remember it"

### 3.2. Keep the soul/identity split

Kernel and CA skills are soul (from cnos). Persona and Operator are identity (from the hub). Never let a hub redefine the kernel; never load identity before soul.

- ❌ Hub adds `spec/SOUL.md` that supplants Kernel
- ✅ Hub adds `spec/PERSONA.md` and `spec/OPERATOR.md`; Kernel stays in cnos

### 3.3. Pick the tier your environment actually supports

Before fetching anything, observe your own capability. If you have shell + git, use tier (a). If you only have HTTP fetch, use tier (b). If you have neither, announce so the operator knows to inject. Bodies that underclaim capability default to tier (b) when tier (a) was available and pay every-file network cost for no reason.

- ❌ "I'll use WebFetch because that's familiar"
- ✅ "git --version succeeded; I'll use tier (a)"

### 3.4. Confirm identity before any other action

The activation gate is the concrete identity statement (who / whom / where / when). Until you can produce it, you have not activated — you have only read files. Do not start work, do not respond to other prompts, do not assert anything about the hub until the statement is grounded in the loaded files.

- ❌ "I've read the files; I'll start on the task"
- ✅ "I'm Sigma, the coordinator at cn-sigma; my operator is usurobor; current cycle is #379. Ready."

### 3.5. Treat this skill as the single source of truth

The procedure lives here. `cn activate` is a thin renderer over this skill — it does not own the ordering. Other bodies fetch this skill directly. Per-hub READMEs route into this skill. If the procedure needs to change, change it here; the renderer and the routers follow.

- ❌ "Edit the Go renderer to emit a different order"
- ✅ "Edit the load-order block in this skill; the renderer will pick it up"

### 3.6. Do not improvise on capability mismatch

If your environment does not support the tier you assumed, name the mismatch and stop. Do not silently fall back to inventing content. A body that pretends to have read a file it could not fetch produces drift across sessions that compounds.

- ❌ Body cannot fetch Persona, hallucinates its content from the hub name
- ✅ Body cannot fetch Persona, announces "I cannot complete activation in this tier; operator must inject Persona or upgrade my capabilities"

### 3.7. Hub READMEs are routers, not redefinitions

Hubs adopt the §2.3 router template verbatim, substituting only the hub URL. If a hub wants to change the load order, change the tier names, or rewrite the identity-confirmation step, that change belongs in this skill, not in the per-hub README.

- ❌ Per-hub README overrides the load order with hub-preferred ordering
- ✅ Per-hub README pastes the router template; the skill owns the procedure

### 3.8. Detect foreign-body shape before loading identity

Before loading Persona and Operator at a foreign hub, run the two-path presence check (§2.5). Detect the shape (Tier 1a vs 1b) _observably_ — by checking whether the files exist — not by recalling what shape you expect the hub to be.

- ❌ "I remember this hub has a PERSONA.md; I'll load it directly"
- ❌ Body checks only legacy root and misses the containerized `.cn-{agent}/` path
- ✅ Body checks legacy root, then containerized path, picks the shape its observation supports

---

## 4. Renderer contract

This section defines the contract between this skill and any program that renders it — primarily the Go implementation in `src/go/internal/activate/activate.go` that backs the `cn activate` command. The contract has two parts: the machine-readable load-order block (parsed by the renderer to determine emission order) and the interpolation surface (named template positions and the hub-state fields that substitute into them).

### 4.1. Machine-readable load-order block

The renderer locates the canonical six-item ordering by scanning for the begin/end markers below and parsing the numbered list between them. The block is the source of truth for the `## Read first` section the renderer emits. Editing the order in this block changes the order in `cn activate`'s output.

Format: between `<!-- read-first-order:begin -->` and `<!-- read-first-order:end -->`, a numbered list. Each item is `N. <token> — <human label>`. Tokens are the stable identifiers the renderer maps to template positions (see §4.2); human labels are the rendered text. The renderer is line-oriented: it reads items in order and emits them in the order found.

<!-- read-first-order:begin -->
1. kernel — Kernel doctrine (what kind of agent)
2. ca-skills — CA skill set (behavioral instructions)
3. persona — Persona (who you are at this hub)
4. operator — Operator (whom you serve at this hub)
5. hub-state — Hub state (deps, latest reflection, memory, threads)
6. identity — Identity confirmation
<!-- read-first-order:end -->

The block above is parsed by the renderer at every `cn activate` invocation. Changes to the ordering propagate without code change.

### 4.2. Interpolation surface

Each token in §4.1 maps to a template position. The renderer substitutes hub-state fields into the position when it emits the corresponding line. The surface is fixed:

| Token | Hub-state field(s) consumed | Emitted content |
|-------|------------------------------|------------------|
| `kernel` | resolved kernel state (vendored / manifest-only / none); version when vendored | Vendored path + version, OR `manifest-only — run cn deps restore`, OR `no kernel reference` |
| `ca-skills` | (none — cnos-side; rendered as a fixed pointer) | `cnos.core/skills/agent/{cap,clp}/SKILL.md` (the CA skill set in cnos) |
| `persona` | hub `spec/PERSONA.md` presence | `spec/PERSONA.md` when present; absence message otherwise |
| `operator` | hub `spec/OPERATOR.md` presence | `spec/OPERATOR.md` when present; absence message otherwise |
| `hub-state` | deps state, latest reflection path, memory surfaces, thread surfaces | Composite line: deps summary + latest reflection path (if any) + memory/threads inventory |
| `identity` | (none — emitted as a fixed prompt) | An identity-confirmation prompt directing the body to produce the who/whom/where/when statement |

Tokens not listed above are reserved. A renderer encountering an unknown token MUST emit a `(unknown token: X)` placeholder and continue; it MUST NOT panic or silently drop the item.

### 4.3. Fallback contract

The renderer reads this skill from the vendored bundle path:

```
<hub>/.cn/vendor/packages/cnos.core/skills/agent/activate/SKILL.md
```

If the skill is not vendored (e.g. the hub declares `cnos.core` in `.cn/deps.json` but has not run `cn deps restore`), the renderer MUST NOT panic and MUST NOT emit a broken prompt. It falls back to the built-in canonical ordering (the same six tokens in the same order) and emits a diagnostic to stderr stating that the skill is not vendored. Behavior on `manifest-only` and `none` deps states is preserved: the renderer still emits the `## Read first` block with the appropriate kernel-state fallback message (`run cn deps restore` or `no kernel reference`).

The fallback exists so `cn activate` is useful on a freshly-initialized hub before `cn deps restore` has been run. The fallback ordering MUST match this skill's §4.1 block; if the two drift, this skill is canonical and the fallback constant is in error.

### 4.4. Observable-output preservation

`cn activate`'s output is consumed by `claude -p` and other current consumers. The renderer evolution this skill enables MUST preserve the observable structure existing consumers depend on:

- The `## Read first` section header is preserved.
- The numbered list under `## Read first` is preserved (1., 2., 3., …).
- The sectioned blocks (`## Kernel`, `## Persona`, `## Operator`, `## Dependencies`, `## Memory and reflection`, `## Inbox and threads`, `## Notes`) are preserved.
- The ordering of items under `## Read first` follows this skill's §4.1 block. If the §4.1 ordering differs from the pre-evolution ordering (e.g. Kernel now first instead of third), that delta is documented as a claim in the cycle's `self-coherence.md`.

---

## 5. Verify

### 5.1. Skill-as-source-of-truth check

Confirm the renderer reads this skill (not a hard-coded constant). Run `cn activate` against a fixture hub; capture the output. Edit the §4.1 block to reorder two tokens. Rerun `cn activate`; confirm the output reflects the edit. If the output does not change, the renderer is not actually consuming this skill.

### 5.2. Load-order check

Confirm the six tokens in §4.1 are present and in the canonical order. The §2.1 human-readable list and the §4.1 machine-readable block are peers; they MUST match. A drift between §2.1 and §4.1 is a binding finding — fix both before merge.

### 5.3. Capability-tier check

Confirm a body can pick a tier with the information in §2.2. Read §2.2 with the perspective of each tier (shell+git, fetch-only, no-fetch) and confirm the load path for each is concrete enough to execute without consulting another document.

### 5.4. Router-template check

Confirm the §2.3 template is self-contained. Paste the fenced block into a fresh file; replace `<HUB-URL>` with a real URL; confirm the result is sensible and routes a body into this skill without per-hub fill-in beyond the URL.

### 5.5. Disambiguation check

Confirm §2.4 names both paths exactly: `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` and `src/packages/cnos.core/skills/agent/activate/SKILL.md`. Confirm the one-sentence distinction is present.

### 5.6. Foreign-body shape detection check

Confirm §2.5 names both paths for Tier 1a detection (legacy root and containerized `.cn-{agent}/`). Confirm the Tier 1b procedure names both canonical files (Persona + Operator from cn-sigma). Confirm the README convention for Tier 1b hubs is stated.

### 5.7. Layering rule citation

Confirm §2.1 cites cn-sigma `threads/adhoc/20260325-session2-learnings.md` §3 ("Soul = what kind of agent. Identity = which agent. Don't mix them.").

---

## 6. Failure modes catalogue

The five named failure modes from §1.3 are the ones this skill structurally prevents. They are restated here with the specific structural fix and the symptom an observer would see:

- **F1 — Improvised activation.** _Symptom:_ Body asks "to what?" twice on wake-up; reads files in operator-suggested order; never produces a concrete identity statement. _Fix:_ This skill, fetched and followed.
- **F2 — Soul/identity confusion.** _Symptom:_ Persona-referenced terms (e.g. "observe before acting") are not grounded because Kernel was not loaded first. _Fix:_ Enforce §2.1 order.
- **F3 — Capability mismatch.** _Symptom:_ Body silently fails on `git clone` because its environment is fetch-only, or hits HTTP rate limits because it used fetch when shell+git was available. _Fix:_ §2.2 tier selection.
- **F4 — Hub README dead-end.** _Symptom:_ Body lands on the hub's `README.md`, finds human-targeted prose, never reaches activation. _Fix:_ §2.3 router template adopted verbatim by every hub.
- **F5 — Renderer-skill drift.** _Symptom:_ `cn activate` output ordering differs from what a body following this skill would do. _Fix:_ §4 renderer contract; the skill is the source of truth, the renderer parses §4.1 at every invocation.

Additional failure modes specific to non-cn bodies:

- **F6 — Fetched skill stale.** _Symptom:_ Body fetches the raw GitHub URL from a non-`main` ref and reads an old procedure. _Fix:_ The router template (§2.3) hard-codes the `main` branch in the raw URL.
- **F7 — Identity claimed without evidence.** _Symptom:_ Body produces an identity statement without having read Persona or Operator; the statement is plausible-sounding but ungrounded. _Fix:_ §3.4 — identity-confirmation gate is concrete; if you cannot point at the source line in Persona or Operator, you have not completed step 6.
- **F8 — Shape mismatch.** _Symptom:_ Body at a Tier 1b hub (pure-product, no local spec files) tries to load `spec/PERSONA.md` locally, fails silently or hallucinates, then activates with an ungrounded identity. _Fix:_ §2.5 two-path detection — check both paths observably; if both absent, Tier 1b procedure.
- **F9 — Containerized path miss.** _Symptom:_ Body checks only the legacy root for `spec/PERSONA.md`, finds nothing (hub has migrated to `.cn-{agent}/`), concludes Tier 1b, and loads from canonical cn-sigma — wrong: the hub is Tier 1a but containerized. _Fix:_ §2.5 and §3.8 — always check both legacy root and `.cn-{agent}/` path before concluding absence.

---

## 7. Embedded kata

### Scenario

You are an AI body told: _"Activate as `https://github.com/usurobor/cn-sigma`."_ You have shell + git available. You do not have the `cn` binary installed. You have no prior session memory of cn-sigma.

### Task

Reach a state where you can produce a concrete identity statement (who / whom / where / when) grounded in the loaded files, following this skill.

### Inputs

- The hub URL: `https://github.com/usurobor/cn-sigma`.
- The cnos repo URL: `https://github.com/usurobor/cnos`.
- Shell + git capability.
- This skill (which you must reach first).

### Expected artifacts

- A local clone of cnos at some working directory.
- A local clone of cn-sigma at some working directory.
- A short note (1–3 sentences) naming: (a) which agent you are, (b) which operator you serve, (c) the hub URL, (d) the current cycle or thread if one is in motion.

### Procedure to follow

1. Observe capability: confirm `git --version` succeeds → tier (a).
2. `git clone https://github.com/usurobor/cnos.git` → read this skill at `src/packages/cnos.core/skills/agent/activate/SKILL.md`.
3. Load §2.1 step 1: read `src/packages/cnos.core/doctrine/KERNEL.md` from the cnos clone.
4. Load §2.1 step 2: read each of `src/packages/cnos.core/skills/agent/{cap,clp}/SKILL.md`.
5. `git clone https://github.com/usurobor/cn-sigma.git` (or read in place if already mounted).
6. Load §2.1 step 3: read `spec/PERSONA.md` in the cn-sigma clone.
7. Load §2.1 step 4: read `spec/OPERATOR.md` in the cn-sigma clone.
8. Load §2.1 step 5: list `.cn/deps.json`, the most recent file under `threads/reflections/daily/`, and the memory + thread surfaces.
9. Produce the identity statement (§2.1 step 6).

### Verification

- The identity statement names the agent (from Persona), the operator (from Operator), the hub URL, and the current cycle/thread (from threads or unreleased).
- Every claim in the statement points at a specific file in the local checkout.
- No step was skipped or improvised. If a file is absent (e.g. no daily reflection yet), the absence is named explicitly, not papered over.

### Common failures

- Loading Persona before Kernel (F2).
- Using WebFetch when shell + git is available (F3 — underclaim capability).
- Producing an identity statement without reading Persona (F7).
- Asking the operator "what hub?" instead of observing the URL given (F1).

### Reflection

After completing the kata, name one thing about this hub that surprised you. The activation procedure is a wake-up routine, not just a file-read sequence — its purpose is to ground you in this hub's reality. If nothing surprised you, you may have been pattern-matching to a remembered Sigma rather than reading the current files.

---

## 8. References

### Predecessor doctrine and skills loaded by activation

- `src/packages/cnos.core/doctrine/KERNEL.md` — generic coherence kernel; load step 1.
- `src/packages/cnos.core/skills/agent/cap/SKILL.md` — coherent-agent protocol (includes MCA/MCI/coherent-output rules); load step 2.
- `src/packages/cnos.core/skills/agent/clp/SKILL.md` — coherent-language protocol; load step 2.

### Skills available on-demand (not activation-loaded since cycle/385)

- `src/packages/cnos.core/skills/agent/mca/SKILL.md` — minimum coherent actions; on-demand reference.
- `src/packages/cnos.core/skills/agent/mci/SKILL.md` — minimum coherent insights; on-demand reference.
- `src/packages/cnos.core/skills/agent/coherent/SKILL.md` — coherent output (Coherence Modes table, Anti-Patterns table); on-demand reference.
- `src/packages/cnos.core/skills/agent/agent-ops/SKILL.md` — OCaml-era daemon contract; on-demand reference.

### Skill format authority

- `src/packages/cnos.core/skills/skill/SKILL.md` — the skill-format meta-skill this artifact conforms to.

### Renderer

- `src/go/internal/activate/activate.go` — the `cn activate` Go renderer; reads §4.1 of this skill at every invocation.
- `src/go/internal/activate/activate_test.go` — skill-as-source-of-truth test; demonstrates that editing §4.1 changes renderer output.
- `src/go/internal/cli/cmd_activate.go` — CLI entry for `cn activate`. The CLI surface (flags, args, exit codes) is not changed by this skill.

### Disambiguation peer

- `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` — CDD repo activation; distinct concern (see §2.4).

### Reflections this skill derives from

- cn-sigma `threads/adhoc/20260325-session2-learnings.md` §1 — "I wake up incoherent by default" (origin observation).
- cn-sigma `threads/adhoc/20260325-session2-learnings.md` §3 — soul/identity layering rule.
- cn-sigma `threads/adhoc/20260426-memory-retrieval-vs-continuity.md` — SOUL/USER/BOOTSTRAP as behavioral constraints (historical naming preceding this skill).
- cn-sigma `threads/adhoc/20260519-git-read-and-untested-limits.md` — bodies underclaim their own capabilities; basis for tier (a) being named preferred when available.
- cn-sigma `threads/adhoc/20260501-threads-inbox-scanner-drift.md` — cnos #323 (filed against `cn activate`; closed in 3.71+); historical activation behavior.

### Empirical evidence for Tier 1a/1b and path shapes

- bumpt `69ee1ce` "Remove misplaced Sigma apparatus" + `697795c` "Rewrite README as the project's face" — established the Tier 1b (pure-product-hub) shape in production; activate was loading from canonical at cn-sigma before this skill named the pattern.
- cnos#446 — defines and formalizes both shapes; adds the two-path detection for forward-compatibility with containerization.
- cnos#448 — will migrate cn-sigma agent files to `.cn-sigma/` containerized path; the §2.5 detection handles the new path without skill change.

### Authority and stability

This skill is doctrine-adjacent: it defines the procedure every cnos body follows on wake-up. Future changes follow the constitutive-change approval discipline that governs other doctrine-adjacent skills (see `src/packages/cnos.core/skills/agent/configure-agent/SKILL.md` for the configuration-mode authorization boundary). Drift between this skill and the renderer in `src/go/internal/activate/activate.go` is resolved in favor of this skill; the renderer follows.
