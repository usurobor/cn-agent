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
  completed: [frontmatter, core-principle, algorithm, define]
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
