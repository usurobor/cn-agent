# cnos Architectural Invariants

**Version:** 1.2.0
**Status:** Active
**Doc-Class:** constitutive
**Canonical-Path:** `docs/alpha/architecture/INVARIANTS.md`
**Owns:** architectural invariants and transition constraints that CDD must validate
**Does-Not-Own:** implementation tactics, coding-style heuristics, release notes, or detailed migration plans

---

## 0. Purpose

This document defines the constraints that keep cnos one coherent system. CDD must validate these constraints in every substantial cycle.

This document distinguishes two kinds of constraints:

- **Active invariants** — already in force; a cycle must not violate them
- **Transition constraints** — target-state architecture decisions; until fully landed, a cycle must not move the system away from them

The distinction matters because some parts of the architecture are already canonical, while others are still active transition work.

---

## 1. Active invariants

### INV-001 — One package substrate

There is exactly one package system. Every distributable unit — doctrine, mindsets, skills, commands, orchestrators, templates, extensions — is distributed through `cn.package.v1` and the package restore/install flow.

**Validation:**

- no second package manifest family
- no second restore/install path
- no separate command/plugin package ecosystem

---

### INV-002 — Package metadata is language-neutral

The package boundary is the correct place for polyglot behavior. Package metadata does not assume the implementation language of:

- a command
- a provider
- a helper executable

**Validation:**

- manifests and package index do not branch on language
- compatibility is expressed through metadata and runtime contracts, not language-specific package formats

---

### INV-003 — Commands, providers, orchestrators, and skills remain distinct

- **skills** choose
- **commands** dispatch
- **orchestrators** execute
- **providers** provide capability

These runtime surfaces may live in the same package. They do not become the same runtime thing.

**Validation:**

- command registry and provider registry are separate
- provider contracts are not reused as command contracts
- orchestrators are not treated as skills
- command dispatch is not keyword activation

---

### INV-004 — Kernel owns policy

The kernel may load or invoke commands, orchestrators, and providers. The kernel alone decides:

- precedence
- activation
- enable/disable state
- permissions
- routing
- protocol acceptance

Packages may declare what they provide. They do not define policy.

**Validation:**

- package/provider code cannot widen its own authority
- routing/permission/preference rules live in the kernel, not in package payloads

---

### INV-005 — Protocol semantics stay above transport/provider implementations

Envelope rules, receipts, routing invariants, rejection policy, and dedupe semantics remain kernel-owned. Concrete transport/provider implementations may vary. Protocol meaning does not.

**Validation:**

- providers implement transport/capabilities
- providers do not redefine protocol truth

---

### INV-006 — Hub runtime state is not source

`.cn/` is active runtime state for a working hub. It is not the authored source tree.

**Validation:**

- `.cn/` is not treated as authored package source
- doctor validates `.cn/` as runtime state, not source content

---

## 2. Transition constraints

These are target-state architectural decisions already accepted by design, but not all of them may yet be fully implemented everywhere. Until they are fully landed, cycles must move the system toward them, not away from them.

### T-001 — The kernel is a package-compatible kernel

`cnos.kernel` is an embedded `cn.package.v1`-compatible manifest. The platform is an instance of its own package/command model.

**Validation:**

- kernel manifest uses the same schema/parser family as normal packages
- kernel commands normalize into the same runtime command descriptor model as other commands

---

### T-002 — The kernel remains minimal and trusted

The kernel owns bootstrap, registries, policy, protocol semantics, and runtime surfaces. It does not own every command, every provider, or all intelligence-bearing content.

**Validation:**

- convenience behavior moves to packages when it does not need to be built-in
- built-in surface trends toward the bootstrap kernel, not toward accretion

---

### T-003 — Go is the sole implementation language

The trusted kernel is monolingual Go. No CGo. No embedded interpreters. No dynamic loading. OCaml is ceased — no OCaml CI, no OCaml releases, no OCaml changes. The existing OCaml source remains in-tree as dead code until Go Phase 5 replaces the full runtime, at which point it is deleted.

**Validation:**

- all CI, release, and coherence workflows use Go exclusively
- no OCaml build or test steps in any workflow
- kernel build remains pure-Go
- non-Go behavior lives at package/provider boundaries, not inside kernel internals

---

### T-004 — Source → artifact → installed is explicit

The package lifecycle must have one explicit authored source, one explicit distributable artifact, and one explicit installed state.

**Current layout** (on main today):

- `src/agent/<class>/` → authored source
- `packages/<name>/<class>/` → built package (mixed: also contains manifests)
- `.cn/vendor/packages/<name>@<version>/` → installed active state

**Target layout** (per BUILD-AND-DIST.md):

- `src/packages/<name>/` → authored source (single source of truth)
- `dist/packages/<name>-<version>.tar.gz` → distributable artifact
- `.cn/vendor/packages/<name>/` → installed active state

The transition is in progress. Cycles must move toward the target layout, not away from it. Neither layout may be written as if the other does not exist.

**Validation:**

- source, artifact, and installed state are not conflated
- derived layers are not edited as if they were source
- docs that describe the layout must say which layout they describe (current or target)

---

### T-005 — Content classes are explicit and finite

The package content-class set is explicit and finite at any given version. New content classes require explicit design, not ad hoc exceptions.

**Validation:**

- no hidden special-case content pipelines
- command/orchestrator/extension behavior uses the declared content classes
- adding a new content class requires an explicit design update

---

## 3. Process constraints

### P-001 — Two-agent minimum applies to substantial cycles

For substantial CDD cycles, one agent authors and a different agent reviews. Small-change cycles may use the documented single-agent exception.

**Validation:**

- substantial cycles show distinct author/reviewer identities
- small-change path is declared explicitly when used

---

### P-002 — Findings before merge

Review findings must be resolved before merge, except for explicit, reviewer-declared design-scope deferrals that are filed before merge.

**Validation:**

- PR/review trail shows closure or explicit issue-backed deferral

---

### P-003 — Hub memory closes substantial release cycles

For release/post-release closure, the cycle must update:

- daily reflection
- relevant adhoc thread(s)

This requirement applies to the post-release closeout of substantial release cycles, not to every tiny code edit.

**Validation:**

- post-release assessment §Hub Memory references committed hub artifacts

---

## 4. Validation model

Constraints are validated by owner surface and cadence.

| Constraint | Owner surface | Cadence |
|---|---|---|
| INV-001 One package substrate | design + review | every substantial cycle touching packaging/runtime boundaries |
| INV-002 Language-neutral package metadata | design + review | every substantial cycle touching package/provider model |
| INV-003 Commands/providers/orchestrators/skills distinct | design + review | every substantial cycle touching runtime registries or package content classes |
| INV-004 Kernel owns policy | design + review | every substantial cycle touching routing, permissions, or provider activation |
| INV-005 Protocol above transport | design + review + post-release | every substantial cycle touching transport/providers; confirmed in release/assessment when applicable |
| INV-006 Hub runtime state is not source | doctor + review | continuous + every cycle touching source/build/runtime boundaries |
| T-001 Kernel is package-compatible | build/check + review | every kernel/package model cycle |
| T-002 Kernel remains minimal | review + post-release | every kernel/built-in command cycle |
| T-003 Go sole language | CI/build/release | every cycle (OCaml ceased) |
| T-004 Source/artifact/installed explicit | build/check + review | every package/build/restore cycle |
| T-005 Content classes explicit and finite | design + review | every cycle touching package content classes |
| P-001 Two-agent minimum | PR/review trail | every substantial cycle |
| P-002 Findings before merge | PR/review trail | every substantial cycle |
| P-003 Hub memory closes release cycle | post-release assessment | every substantial release cycle |

### Rule

A substantial cycle touching one of these areas must include a constraints section in the design artifact and the reviewer must confirm the affected constraints are preserved, tightened, or explicitly revised.

---

## 5. Amendment rule

This document may change only through explicit proposal and confirmation with evidence. Convenience is not enough to weaken a constraint. Persistent friction may justify restructuring one, but the restructuring must be explicit and justified.
