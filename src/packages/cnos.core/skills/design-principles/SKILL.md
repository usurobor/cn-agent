---
name: design-principles
description: Design systems so policy stays stable, volatile decisions stay hidden, interfaces stay truthful, and runtime surfaces do not smear together. Use for architectural decomposition, boundary design, package splits, registry design, and refactors.
artifact_class: skill
kata_surface: embedded
governing_question: How do we decompose a system so each part has one real reason to change, dependencies point the right way, and replacements do not break the contract?
triggers:
  - splitting a package or module
  - defining a new runtime surface
  - moving logic out of a CLI layer
  - deciding what belongs in core vs package/provider
  - introducing or revising an interface
  - reviewing architectural refactors
---

# Design Principles

## Core Principle

**Coherent design hides volatile decisions behind stable contracts and keeps policy above detail.**

A system is well designed when:

- each boundary has one real reason to change,
- dependencies point toward policy, not away from it,
- interfaces say no more and no less than the runtime can actually support,
- and the same concept is not reintroduced under a second name somewhere else.

## Algorithm

1. **Define** — name the changing parts, the stable parts, the boundaries, and the failure mode.
2. **Unfold** — design by reasons to change, information hiding, dependency direction, truthful interfaces, and explicit runtime surfaces.
3. **Apply** — keep contracts narrow, layers honest, and duplicative systems out.
4. **Verify** — check that the decomposition still matches reality after the code and docs land.

---

## 1. Define

### 1.0. Classify the design problem first

Use this skill when the question is architectural:

- what belongs in core vs package/provider
- what should be a command vs provider vs orchestrator vs skill
- how to split a large module
- how to design a package boundary
- how to define or narrow an interface
- how to stop convenience code from becoming a second system

Do not use this skill for:

- local syntax/style questions
- framework taste disputes
- one-off micro-optimizations with no boundary consequences

### 1.1. Identify the parts

A coherent design problem usually has:

- **invariants** — what must remain true
- **volatile parts** — what may be replaced
- **boundaries** — where they meet
- **contracts** — what the boundary promises
- **registries** — what is discoverable at runtime
- **sources of truth** — where the fact actually lives
- **reasons to change** — what would force this part to move

- ❌ Start by drawing folders or classes
- ✅ Start by naming what is stable, what is volatile, and why each part changes

### 1.2. Articulate how they fit

Use three questions:

- **Information hiding:** what volatile decision should this module hide?
- **Dependency direction:** which side owns policy, and do dependencies point toward it?
- **Substitutability:** if this thing claims to be replaceable, can it really be substituted without breaking callers?

A decomposition is coherent when all three answers line up.

- ❌ Split by processing steps or convenience only
- ✅ Split by hidden volatility, dependency direction, and truthful substitutability

### 1.3. Name the failure mode

Design fails through:

- **leakage** — one layer does work that belongs elsewhere
- **duplication** — the same fact lives in more than one authoritative place
- **false substitutability** — an interface claims interchangeability the implementations do not support
- **surface smearing** — commands, providers, orchestrators, and skills collapse into one vague mechanism
- **premature canonicalization** — target-state architecture is written as if already fully implemented
- **authority drift** — docs, runtime contract, help, doctor, and code disagree about what exists

- ❌ "It works if you know the convention"
- ✅ "The contract, registry, and runtime all reveal the same system"

---

## 2. Unfold

### 2.1. Decompose by reason to change

Do not decompose by processing order alone. Good decomposition groups together code that changes for the same reason and separates code that changes for different reasons.

Examples:

- CLI dispatch changes when the operator surface changes
- doctor logic changes when health rules change
- transport implementation changes when the backend changes
- protocol semantics change when the message contract changes

- ❌ Put dispatch, policy, and implementation in one file because the request passes through all three
- ✅ Keep dispatch thin and move domain logic to the module that owns the rule

### 2.2. Hide volatile decisions

A module should hide the decision most likely to vary.

Examples:

- package artifact resolution hides index/URL layout changes
- provider adapters hide transport/process details
- contract builders hide rendering/layout details
- activation loaders hide discovery rules

- ❌ Expose every internal choice at the boundary
- ✅ Hide what may change and expose only what the rest of the system actually needs

### 2.3. Keep policy above detail

Policy should not depend on the details it governs.

The kernel or core should own:

- precedence
- permissions
- routing
- protocol semantics
- acceptance/rejection rules

Packages/providers may supply implementations. They should not decide policy.

- ❌ Let a package/provider widen its own authority or rewrite precedence
- ✅ Let the kernel decide whether and how the implementation is used

### 2.4. Make interfaces truthful

Define an interface only when a real consumer needs substitution.

A truthful interface:

- is small
- maps to actual behavior
- does not promise more than all implementations can provide
- does not collapse unrelated roles into one type

- ❌ One interface for commands and providers because both "run code"
- ✅ Separate command execution from provider protocol because they have different contracts

If a thing is not really substitutable, do not hide it behind an interface.

### 2.5. Separate runtime surfaces

Keep these runtime surfaces distinct:

- **skills** choose
- **commands** dispatch
- **orchestrators** execute
- **providers** provide capability

They may live in the same package. They do not become the same runtime thing.

- ❌ Trigger commands by keyword like skills
- ❌ Treat providers as just hidden commands
- ✅ Give each surface its own registry, contract, and policy boundary

### 2.6. Separate source, artifact, and installed state

Keep three layers explicit:

- **source** — what humans edit
- **artifact** — what build produces
- **installed** — what the runtime consumes

A system becomes hard to reason about when these collapse.

- ❌ Edit generated package output as if it were source
- ✅ Make the build pipeline explicit and keep generated/install layers derived

### 2.7. Keep registries explicit

Anything the runtime discovers should have:

- one descriptor model
- one registry builder
- one precedence rule
- one status/doctor/help story

That applies to:

- commands
- providers
- orchestrators
- activation entries

- ❌ Discovery by scattered special cases
- ✅ One normalized runtime descriptor regardless of source form

### 2.8. Preserve determinism and receipts

Architectural boundaries are stronger when they leave visible evidence.

For important runtime operations, define:

- what gets registered
- what gets rendered
- what gets logged
- what gets receipted
- what doctor/status can inspect

- ❌ Degraded path only visible in transient logs
- ✅ Degraded path visible in receipts or structured runtime state

---

## 3. Rules

### 3.1. One reason to change per boundary

If a module, package, or registry changes for several unrelated reasons, the boundary is probably wrong.

### 3.2. One source of truth per fact

If two places claim authority over the same fact, the design is already drifting.

- ❌ Package metadata in one doc, runtime contract in another, neither canonical
- ✅ One canonical source, many derived views

### 3.3. Policy belongs in the kernel/core

Precedence, permissions, routing, and protocol acceptance belong in the trusted core.

### 3.4. Details belong in packages/providers

Implementations that can vary should sit below the policy layer and behind explicit contracts.

### 3.5. Interfaces belong to consumers

Define interfaces where a consumer actually needs substitution. Do not create broad "future-proof" interfaces with one implementation and no real variation pressure.

### 3.6. Registries must be normalized

Different source forms may exist. They should normalize into one runtime descriptor model.

- ❌ Three command kinds, three unrelated registration paths
- ✅ Kernel, repo-local, and package commands all normalize into one `CommandSpec`

### 3.7. Preserve explicit runtime surface boundaries

Do not smear:

- skills into commands
- commands into providers
- providers into orchestrators
- orchestrators into skills

### 3.8. Keep transition constraints honest

If something is target state, say so. Do not write future architecture as if already fully implemented.

### 3.9. Degraded paths must be visible

A fallback or degraded behavior must be:

- intentional
- inspectable
- testable
- reviewable

### 3.10. Prefer package/install cohesion over topic labeling

Package boundaries should follow common install/use units, not just conceptual similarity.

### 3.11. Keep the bootstrap kernel small

Built-ins should justify themselves by bootstrap necessity, not convenience.

### 3.12. Build and review the actual boundary

Before push/merge, verify:

- the docs
- the manifest/registry
- the runtime contract
- the implementation

all describe the same system.

---

## 4. Verify

### 4.1. Boundary check

- What is the actual reason this part changes?
- Does the current boundary match that reason?

### 4.2. Policy check

- Does policy still live above detail?
- Can a package/provider widen its own authority?

### 4.3. Contract check

- Is the interface truthful?
- Can implementations really substitute for one another?

### 4.4. Registry check

- Is there one normalized runtime descriptor?
- Do help/status/doctor/runtime contract agree on what exists?

### 4.5. Source/artifact/install check

- Is it clear what is authored, what is built, and what is installed?
- Are any derived layers being treated as source?

### 4.6. Review/build check

- Are docs and code aligned?
- Are degraded paths visible and test-covered?

---

## 5. Kata

### 5.1. Kata A — dispatch-only CLI

A CLI package currently mixes dispatch and domain logic. You need to make dispatch-only.

**Refactor so:**

- CLI owns argument mapping and invocation only
- domain packages own doctor/status/restore/build logic
- help/registry behavior stays consistent

**Skills:** design-principles, review, writing

**Given:**

- one mixed CLI file
- one target runtime registry model
- one expected help/status surface

**Produce:**

- thin CLI wrapper
- extracted domain module(s)
- preserved registry/help behavior
- explicit source of truth for command metadata

**Check:**

- dispatch layer no longer owns domain logic
- docs and runtime model still align
- no hidden special-case registration drift

**Anti-patterns:**

- moving logic but leaving side-effect policy in the wrapper
- breaking help/doctor/status consistency
- inventing a second registry

**Reflection:**

- Did the boundary now match the true reason to change?
- Did the kernel keep policy while details moved outward?

### 5.2. Kata B — command vs provider split

A tool can be used both by operators and by the runtime. You need to decide whether it is a command, a provider, or both.

**Design the package and runtime surfaces so:**

- operator invocation stays exact-dispatch
- runtime capability use stays under provider policy
- the package can contain both without merging the contracts

**Skills:** design-principles, cdd/design, review

**Given:**

- one concrete tool candidate
- current package-system and runtime-extension rules
- current command registry model

**Produce:**

- one package-level design
- explicit command/provider contracts
- explicit kernel-owned policy boundary

**Check:**

- no contract smearing
- one package substrate, not two ecosystems
- runtime and operator stories both remain clear

**Anti-patterns:**

- calling providers "just commands"
- adding a second plugin/package system
- hiding policy in the package/tool implementation

**Reflection:**

- Did the design preserve one model while allowing different source forms?
- Did the package boundary simplify or complicate the system?
