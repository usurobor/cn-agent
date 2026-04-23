# Skill Architecture

Why skills compose instead of inherit, and what that means for a coherence-first system.

---

## 1. What a Skill Is

A skill is not a class. No private state, no encapsulated mutation, no inheritable behavior. A skill is advisory text loaded into context to bias a downstream judgment.

A skill has four properties:

- **Trigger** — when it applies
- **Input contract** — what it assumes about the artifact it sees
- **Output contract** — what it guarantees about the artifact after it runs
- **Transformation along one axis** — what it actually changes

That is the signature of a function with preconditions, a body, and postconditions. Not the signature of a class.

The choice between composition and inheritance is not an engineering trade-off for skills. It is a question about what kind of object a skill is. Inheritance models a skill as a class. Composition models a skill as a function. Only one is faithful to what a skill actually does.

---

## 2. Why Inheritance Fails for Skills

Inheritance creates three problems skills cannot absorb.

### 2.1. Implicit coupling

When a child skill inherits from a parent, the child's behavior depends on the parent's rules in a way the loading system cannot inspect. If the parent changes, every child silently changes. In a system that selects skills by trigger and runs them in order, every dependency must be visible to the loader. Inherited behavior is invisible by definition.

### 2.2. Single-axis commitment

A class hierarchy commits to one variation axis. Once `essay` inherits from `write`, adding voice or form as second axes requires multiple inheritance (which most systems forbid or punish) or flattening everything into the one hierarchy guessed at months ago. Skills cross many axes — coherence, voice, form, naturalness, domain — and no single hierarchy captures them all.

### 2.3. Brittle base prompt

When the parent skill iterates, all children inherit the change without consenting. In code, this is annoying. In skills, where "behavior" is rules biasing a generative system, children change in ways their authors did not anticipate and cannot easily test. The base class problem becomes a base prompt problem.

These are not theoretical. They are the predictable failure modes of making skills inherit, and the pattern is the same as in code — only with less tooling to catch it.

---

## 3. Why Composition Fits

Composition treats each skill as an independent function with a clean contract. Combining skills means running them in sequence or selecting by trigger, each knowing only what it needs, producing what it promises.

### 3.1. Multiple axes without commitment

A draft can be transformed along coherence, then form, then voice, then naturalness, in any order the situation requires. No axis lives in a hierarchy. No axis forces a choice about which is "more fundamental." They are peers because they are axes, not levels.

### 3.2. Explicit dependencies

A composed pipeline declares what runs before what. No inherited behavior, no silent side-effects from parent changes. When a skill iterates, only artifacts that explicitly compose it are affected. The loader can inspect the pipeline.

### 3.3. Skill replaceability

A composition lets you swap one skill for another that satisfies the same contract. Inheritance lets you subclass but not easily substitute. For a system where skills are added and refined continuously, replaceability is more valuable than reuse.

The Unix pipeline is the canonical instance: small tools, clean interfaces, no shared state, composed at the call site. GoF said "favor composition over inheritance" in 1994. Thirty years of language design agreed — by weakening inheritance (Java's `final`, Kotlin's sealed classes) or removing it (Go, Rust, Clojure). Skills are a young enough abstraction to learn the lesson without relearning it.

---

## 4. Where Coherence Raises the Bar

In a system whose foundation is coherence, the question is not just which model is more flexible. It is which model preserves coherence under change.

Coherence: all parts of a system point at the same thing. A coherent skill has one axis, one trigger, one set of contracts, one source of truth for each rule. A coherent skill system has skills that combine without overlap, drift, or hidden order.

Inheritance violates coherence in two specific ways.

### 4.1. Two homes for one rule

When a child inherits from a parent, both files contain the rule — the parent declares it, the child inherits it. As the child evolves, it may override or shadow the rule. Two slightly different versions in two places, and the system has no way to know which the loader applies.

This is exactly the failure mode write names ("say a fact once, then point to it") and compose names ("move shared rules to one owner"). Inheritance institutionalizes that failure.

### 4.2. Implicit axes

A skill that inherits from another inherits the parent's reasons-to-change on top of its own. Multiple reasons to change, hidden under one file. The boundary is no longer crisp. The skill is no longer single-axis. By the system's own rules, it is no longer coherent.

### 4.3. Composition creates neither failure

Each skill owns its own rules. Each has one axis. Each skill's reasons-to-change are exactly its own. When two skills combine, the combination happens at the use site — visible and explicit — not the definition site — implicit and inherited.

Composition is not just preferable for skills. It is the only model consistent with coherence as a foundational principle.

---

## 5. The Principle

**For systems whose foundation is coherence, composition is the only architectural primitive that preserves the foundation.**

Inheritance trades coherence for concision. Concision is not worth that trade.

Skills compose because they are functions, not classes. They preserve coherence because composition keeps every dependency explicit, every rule single-owned, every axis singular. They scale because adding a new skill never requires editing an old one. They stay honest because the loader can inspect every relationship.

---

## 6. Evidence from cnos.core

The cnos.core skill set organized itself into the shape this principle predicts.

**Foundation layer** — three peers along orthogonal axes:

| Skill | Axis |
|-------|------|
| skill | artifact form and classification |
| write | prose coherence |
| design | system decomposition |

**Specialization layer** — one application of a foundation skill to a specific artifact class:

| Skill | Relationship |
|-------|-------------|
| compose | design applied to skill artifacts |

The compose ↔ design relationship confirms the architecture. compose is structurally what design looks like when applied to artifacts of class skill. The rules map almost one-to-one:

| design | compose |
|--------|---------|
| One reason to change per boundary | One axis per skill |
| One source of truth per fact | One owner per shared rule |
| Interfaces belong to consumers | Contracts belong to consumers |
| Preserve runtime surface boundaries | Compose, do not extend |

If skills could inherit, compose would be a child of design. Because they don't, compose is a peer that applies the design formula to a specialized artifact class, declares what it adds (classification, triggers, the publish-or-compose verdict), and points back for shared rules.

The skills compose to form the skill system. The rules of composition are themselves expressed compositionally. The system is coherent with itself.

---

## 7. The Deeper Test

The cnos system did not arrive at this by argument. It arrived by building. The four cnos.core skills are peers, not a hierarchy, because every attempt to make one inherit from another would have violated the system's own coherence rule. The architecture discovered itself.

That is the test of a foundational principle. Not that you can state it, but that the system you build by following it organizes itself the way the principle predicts.

Coherence demanded composition. Composition delivered coherence. Same move at different scales.

---

## Authority

This document is doctrine. It governs the architectural model for skills in cnos.

- skill governs skill form and classification
- write governs prose coherence
- design governs system decomposition
- compose governs skill boundaries and composition judgment
- This document governs why skills compose and what that means architecturally

If this document and a skill disagree on the architectural model, this document governs. If this document and COHERENCE.md disagree on what coherence means, COHERENCE.md governs.
