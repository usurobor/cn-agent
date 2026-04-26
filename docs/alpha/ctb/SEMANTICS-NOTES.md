# CTB Semantics Notes

**Status:** Conceptual rationale and design discussion. Non-normative.
**Date:** 2026-04-26
**Companion to:** `LANGUAGE-SPEC.md` (normative), `CTB-v4.0.0-VISION.md` (strategic).

---

## 0. Why this document

The Vision sets direction. The Spec sets rules. Some material belongs in neither: the deeper observations about *why* the rules take the shape they do, and *what practice taught us* before we could name it.

This document collects that material. It is the place where conceptual moves are recorded so that future revisions of the spec can be made deliberately rather than through drift.

If a claim here ever becomes load-bearing for an implementer, it should migrate into the spec.

---

## 1. Skills are modules plus callable signatures

The earliest framing was that a skill is a function. That is correct but incomplete. A function is a callable; a skill is a callable *plus* a self-contained packaging of the rules, contracts, and kata that make the call meaningful.

Three properties hold simultaneously:

- A skill has a **module identity** — a stable path, a set of files, an authority surface.
- A skill has a **callable signature** — name, triggers, scope, inputs, outputs, calls.
- A skill has **executable content** — today prose, tomorrow `.coh` equations, eventually a hybrid.

Treating a skill as only a function loses the module identity, which is what makes authority and ownership tractable. Treating a skill as only a module loses the signature, which is what makes invocation tractable. Both halves matter; the spec models both.

This is why the language unit is a *module with a signature*, not a *function*.

---

## 2. Why functions, not classes

The architectural argument is in `docs/alpha/doctrine/SKILL-ARCHITECTURE.md`. The short form:

- Inheritance creates implicit coupling. Skills cannot afford implicit coupling because the loader must inspect every dependency.
- Inheritance commits to one variation axis. Skills cross many axes — coherence, voice, form, naturalness, domain — and no single hierarchy captures them.
- Inheritance creates two homes for one rule. The skill system's own coherence rule (one owner per fact) forbids that.

Composition fits because it keeps every dependency explicit, every rule single-owned, every axis singular. It scales because adding a new skill never requires editing an old one. It stays honest because the loader can inspect every relationship.

The deeper observation: composition is not a style preference here. It is the only model consistent with coherence as a foundational principle. A skill system whose foundation is coherence can have only one architectural primitive, and composition is that primitive.

---

## 3. Why scope is first-class

In most languages scope is implicit — lexical, structural, derived from where a definition appears. In CTB scope is part of the signature.

The reason is the agent runtime. A skill at task-local scope and a skill at global scope mean different things at execution time:

- A global skill participates in the whole activation. Its rules apply across every nested call.
- A role-local skill is bound only while a role is held. When the role exits, the skill's authority exits.
- A task-local skill is bound only inside one dispatched task. Sibling tasks cannot see it.

This is the only structure that makes authority enforceable. Without explicit scope, "this rule applies here" becomes "this rule applies wherever the runtime happens to load it" — which is exactly the failure mode that produced the four-fold duplication of the large-file rule before it was hoisted to a global aspect.

Scope is first-class because authority is real, and authority needs a region.

---

## 4. Why composition is explicit invocation

A skill composes by being *called*, not by being mixed in.

The implication is that the call site is the composition site. There is no place where two skills "merge" before runtime. A skill stays whole; another skill invokes it; the invocation is where their behaviors meet.

This is why the spec prescribes:

- `calls` lists every dispatchable target
- `runs_after` / `runs_before` declare order at the use site
- inheritance MUST NOT exist
- shared rules have one owner; other surfaces point to the owner

Each of these constraints is the same principle expressed at a different surface: combination happens at the call site, visibly, and never at the definition site, implicitly.

The Unix pipeline is the canonical instance. CTB inherits the same shape: small modules, explicit contracts, no shared mutable state, composed at the call site.

---

## 5. Why global aspects exist

The strict reading of "skills compose" would forbid any rule that is not owned by exactly one skill. Practice broke that reading.

The CDD package surfaced four cross-cutting requirements that no individual skill could own:

1. A rule about how large files are written, applicable to every artifact-producing skill.
2. A rule about authority precedence between the canonical method and role skills.
3. A rule about which triggers are public dispatch entries.
4. A rule about the load-order graph between skills.

These are aspects of the package, not of any one skill. The pragmatic move was to give them a single owning document (typically the package's canonical method file) and to forbid restatement in individual skills.

The conceptual move is more important: aspects exist because the system has more than one structuring axis. Skills compose along the call axis; aspects compose along the cross-cutting axis. Both axes are real. The language must model both, or it forces every cross-cutting concern to be duplicated and drift.

The spec's §7 is the minimum mechanism for this.

---

## 6. Why static and dynamic calls must both exist

A skill needs to be statically inspectable — the loader must know which targets it can dispatch — and runtime-flexible — the skill must choose among targets based on runtime state.

The reconciling rule is:

- The set of *possible* call targets is declared statically (`calls`).
- The choice *among* declared targets is made dynamically (by trigger match against runtime state).

This gives the loader a complete static call graph while leaving selection policy in the caller's hands.

The alternative — fully static dispatch — makes branching unworkable. The alternative — fully dynamic dispatch — destroys inspectability. The split is principled.

The CDD γ role illustrates it: γ has a fixed `calls` list (α, β, lifecycle sub-skills), but γ chooses among them based on cycle phase. The loader can verify the graph; γ controls the policy.

---

## 7. The triadic package taught the invocation model

The CDD triad surfaced an execution model the language did not yet name:

| Language concept | Triadic equivalent |
|---|---|
| Call stack | `operator → γ → α/β` |
| Function scope | Each role owns its local artifacts |
| Arguments | Dispatch prompt + issue context |
| Return value | Close-out artifact |
| No upward mutation | Lower role MUST NOT write higher role's artifacts |
| Tail-call | γ resolves trivial findings inline rather than dispatching a new cycle |
| Error handling | γ unblocks when α or β is stuck |

These were not designed as language features. They were the only way the package could operate without losing coherence. When we read them back, they were already a call model with scoped authority — the structure the spec now names explicitly.

The package was the prototype runtime. The spec is what the runtime asked for.

---

## 8. The plan-as-data discipline

A skill returns a plan. A plan is data. A separate runtime executes the plan.

The discipline this enforces:

- Skills are pure. Two evaluations on the same input produce equal plans.
- Plans are inspectable. A peer can read, hash, lint, or refuse a plan before execution.
- Effects are gated. The runtime executes only what its capability set permits.
- Determinism is bounded above the runtime, not below. Below is the world.

The plan is the boundary between language and world. It is also the primary artifact of the language: what one skill produces is what another skill consumes.

When two skills compose, their plans compose. When a triad releases a cycle, the cycle's external artifact is a plan that another agent can re-evaluate. The whole-to-whole composition principle (one-as-three appears at its boundary as one articulated whole) is supported by this: the boundary artifact is data.

---

## 9. The relation between practice-side markdown and the future formal syntax

Today a skill body is prose. Tomorrow it will be a `.coh` term. The signature is already structured (frontmatter); the body is not.

The migration path:

1. Today: signature in YAML, body in markdown. The signature is mechanically checked; the body is read.
2. Near term: signature in YAML, body has a structured contract section that the linter parses. Rules are still prose but bound to declared contracts.
3. Mid term: signature in YAML, body is a `.coh` term that evaluates to an effect plan. Prose becomes annotation, not authority.
4. Long term: a hybrid surface where structured terms and prose annotation share a canonical syntax.

The spec is positioned to be stable across this migration: it constrains the signature and the call model, both of which are already in YAML. The body's evolution is a kernel-spec question, not a language-spec question.

The current markdown body is therefore not a temporary scaffolding to be discarded. It is a record of the rules that the future formal syntax must be able to express. Each rule the language eventually formalizes was already running in prose.

---

## 10. On the name "CTB"

CTB is the right internal name. It is anchored in the theory stack (Coherence Calculus → TSC → CTB → coherent agents / git-CN) and in the repository layout. It is precise enough for research, architecture docs, and versioned evolution.

Whether CTB is the right *public* name is a different question. "Triple Bar" is meaningful inside the stack and opaque outside it. A future public name might foreground coherence, plans, skills, or articulation.

The decision rule: the public name can be revisited when the spec, the runtime, and onboarding stabilize. Until then, naming is a distraction. CTB stays.

---

## 11. The strongest claims

If only three sentences from this document survive into spec or doctrine, they are:

1. **Skills are scoped modules with callable signatures, not functions and not classes.**
2. **Composition is explicit invocation. Inheritance is forbidden because it institutionalizes the failure mode the system's own coherence rule names.**
3. **Invocation is a real runtime model — root dispatch, role binding, local expansion — and scope is the region in which authority and mutation are valid.**

These are the load-bearing observations. The spec's normative content is the consequence.

---

## 12. Open conceptual questions

- **What is the right contract grammar for `inputs` / `outputs`?** Today they are prose. A typed shape would let the runtime check satisfaction. But typed shapes risk forcing premature precision; the practice may not yet know what types should exist.
- **How do aspects compose across packages?** When two packages declare aspects that overlap (e.g. both forbid duplication of file-writing rules), is the union the effective aspect, or does the importing package's aspect win?
- **What is the calling convention for return values?** Positional terms map cleanly to the kernel; named close-out fields map cleanly to the practice. Both are needed; the relation between them is unsettled.
- **What is a witness?** The kernel uses witnesses as counterexamples to totality. The skill language uses close-outs as structured failure terms. These are probably the same object at different scales, but the equivalence has not been worked out.
- **Is `requires` a precondition or documentation?** If the loader enforces it, packages must be careful about what they put there. If it is documentation, the field is weaker than it looks.

These are conceptual debts. They do not block the spec, but they will eventually be paid.

---

## Authority

This document is non-normative. It exists to preserve conceptual progress and to make spec revisions deliberate.

If this document and `LANGUAGE-SPEC.md` disagree on a rule, the spec governs.

If this document and `CTB-v4.0.0-VISION.md` disagree on strategy, the Vision governs.

If a claim here becomes load-bearing for implementers, it should be moved into the spec rather than relied on here.
