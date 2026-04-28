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

## 13. There are no skills — only agents

The v0.1 spec drew a line between "skill" (pure function, no state) and the agents that invoke skills (stateful, looping, coordinating). Practice erased that line. A "skill" is a narrow-scope agent. A "role" is a wider-scope agent. A "triad" is three agents composing into one agent at the next scale. There is one primitive — the agent — varying in scope and lifetime, not in kind.

This simplification is not cosmetic. It resolves the tension between "skills are pure with no state" and the immediate practical need for state (role loops, memory, coordination). The resolution is the same one functional programming found: state within scope, effects as data, composition by dispatch.

### 13.1 The FP correspondence

The agent composition problem is the same problem FP solved: how to compose pure transformations when the real world has state, effects, and sequencing.

| FP concept | Agent equivalent |
|---|---|
| Pure function | Narrow-scope agent: inputs → effect plan |
| Monad / effect system | Agent loop (UIE-V): sequences effects within scoped state |
| IO monad | Runtime bridge: executes effect plans under capability grants |
| Algebraic effects | `calls` / `calls_dynamic`: declared effect vocabulary |
| Closure / captured environment | SOUL: activation context threaded into every nested dispatch |
| Higher-order function | Agent dispatching subagents (γ dispatching α/β) |
| Referential transparency | Same inputs → same plan (reproducibility) |
| Type class / trait | Global aspect: cross-cutting constraint, one owner, no inheritance |
| Module / functor | Agent-as-module with callable signature |
| Scope / lifetime / region | `task-local` (stack frame) / `role-local` (region) / `global` (process) |
| do-notation | UIE-V loop: each step yields a plan; the loop binds them |

This is not analogy. It is the same structure, discovered independently through agent coordination practice. The correspondence matters for three reasons:

1. **Known design space.** Implementers can draw from decades of work on effect systems, region types, and module functors instead of reinventing them.
2. **Constraint inevitability.** The spec's rules (no upward mutation, scoped state, explicit dispatch, composition without inheritance) are not arbitrary choices — they are the same constraints FP arrived at for the same reasons.
3. **Migration path.** The move from prose skills to `.coh` programs is the move from informal equational reasoning to a language with a type system. The FP correspondence tells us what that type system needs: effect types, region-bounded state, and module signatures.

### 13.2 Key structural identifications

**Effect plans are monadic.** An agent returns `Plan a`, not `a`. The runtime interprets the plan. This is the same separation as Haskell's `IO a` — the description of an effect is pure; the execution is not. The plan-as-data discipline (spec §9) is the monadic discipline.

**Agent loops are do-notation.** The UIE-V loop (Understand → Identify → Execute → Verify) is sequencing in the effect monad. Each step produces a plan; the loop binds them into a composite plan. A narrow agent has no loop — it is a single expression. A role agent runs a loop — it is a do-block. Same monad, different syntax.

**Scope is region-based.** `task-local` is a stack frame: state lives for one call and dies. `role-local` is a region: state lives while the role is bound and dies when the role exits. `global` is the process. The rule "state MUST NOT outlive scope" is the region discipline. The rule "no upward mutation" is the no-escape constraint.

**The SOUL is a closure.** It is the captured environment threaded into every nested dispatch. Every agent in the activation sees the SOUL's invariants — the same way every function in a closure sees the captured bindings. The SOUL is not "configuration" — it is the lexical environment of the agent's entire execution.

**Global aspects are type classes.** A type class defines a constraint that any type in scope must satisfy. A global aspect defines a constraint that any agent in the package must satisfy. Both are resolved at the package/module level, not per-instance. Both forbid duplication (orphan instances / restated aspects). Both provide a single canonical implementation per scope.

### 13.3 What the v0.1 spec got wrong

v0.1 said: "A skill has no private state, no inheritance, no implicit ambient behavior. It is a function with a contract." The "no state" claim was immediately contradicted by practice — role agents maintain state across loop iterations, the SOUL persists across the activation, memory carries across sessions.

The correct claim is: **an agent has no state that escapes its scope**. Narrow agents have no persistent state because their scope is a single invocation. Role agents maintain state within their role lifetime. The SOUL persists because it is global-scoped. None of these escape their declared region. The purity guarantee is not "no state" — it is "no state leakage."

This is exactly the move from "pure functions with no side effects" to "effects tracked and bounded by a type system." The language grew up.

---

## 14. Composition operators: how agents combine

§13 establishes that agents are the single primitive. §6 of the spec says "composition is the only mechanism by which agents combine." But the spec only declares what an agent CAN dispatch (`calls`), not HOW dispatches combine. This is like having a type signature with no function body.

Practice already uses a small set of composition operators — they are just written in prose. Naming them:

### 14.1 The operator set

| Operator | FP equivalent | Meaning | Example |
|---|---|---|---|
| `>>` | sequence / then | Do A, then B | `implement >> self-cohere >> request-review` |
| `>>=` | bind / chain | Do A, feed result to B | `review >>= \findings → notify(α, findings)` |
| `\|\|\|` | parallel / fan-out | Dispatch both concurrently, scoped isolation | `dispatch(α \|\|\| β)` — β polls while α implements |
| `case` | pattern match | Branch on return value | `case verdict of RC → loop; A → merge` |
| `fix` / `loop` | fixed point / iterate | Repeat until condition met | RC loop: `fix (α.fix >> β.re-review) until A` |
| `wait` / `join` | await / barrier | Block until parallel agent signals | `β.wait(α.ready)` |
| `try` / `recover` | catch / handle | Handle failure, decide recovery | γ unblocking stuck α or β |

These are exactly the operators in any concurrent effect system. The agent-specific property is that each operand in `|||` runs in its own scope — α cannot see β's state.

### 14.2 CDD expressed in operators

The CDD triadic cycle, currently ~300 lines of prose algorithms across CDD.md and three role skills, is one composition expression:

```
cdd-cycle =
  γ.observe
  >> γ.select
  >> γ.issue
  >> (α ||| β)                    -- parallel dispatch, scoped isolation
  >>= γ.triage                    -- bind: both close-outs flow to γ
  >> γ.close

α =
  load
  >> implement
  >> self-cohere
  >> request-review
  >> fix-loop                     -- fix: repeat until β approves
  >> closeout

β =
  load
  >> wait(α.ready)                -- join: block until α signals
  >> review
  >> case verdict of
       RC → notify(α) >> wait(α.fixed) >> β   -- recursive: re-review
       A  → merge >> release >> closeout

fix-loop =
  fix (wait(β.findings) >> fix-findings >> re-cohere >> re-request)
  until (β.verdict == A)
```

This is not a proposal for concrete syntax. It is a demonstration that the prose algorithms are already this expression — they just take 300 lines to say it because natural language cannot express `|||`, `>>=`, or `fix` concisely.

### 14.3 What this means for the spec

LANGUAGE-SPEC v0.2 §6 (Composition) needs to name these operators. The spec currently says "composition is explicit invocation" — correct, but incomplete. Explicit invocation has structure: it is sequential, parallel, conditional, iterative, or binding. The operators are not optional extras; they are the composition model.

The minimum addition to the spec:

1. **Sequential** (`>>`) — run A then B. Already implicit in every algorithm section.
2. **Bind** (`>>=`) — run A, feed its return to B. Already implicit in "close-outs flow to γ."
3. **Parallel** (`|||`) — run A and B concurrently with scoped isolation. Already implicit in "dispatch α and β at the same time."
4. **Case** — branch on return value. Already implicit in "if RC... if A..."
5. **Fix** — iterate until condition. Already implicit in the RC loop.
6. **Wait/join** — synchronize with a parallel agent's signal. Already implicit in "β waits for α's PR."

Each of these is already in the practice. The spec just needs to name them so that:
- Implementations can parse them
- Linters can verify composition well-formedness
- New triadic or multi-agent protocols can be expressed without 300 lines of prose

### 14.4 Prior art: four traditions solving the same problem

The composition operators in §14.1 were not invented here. Four traditions independently converged on the same structure. Each contributes something CTB needs.

**Erlang/OTP: supervision trees as composition.**

The CDD triad IS an Erlang supervision tree. γ is the root supervisor; α and β are child workers. The RC loop is a `one_for_one` restart: γ restarts α while β continues watching. "Let it crash" is the same principle as "return a witness to the caller" — don't handle errors locally, let the supervisor decide.

OTP behaviours (`gen_server`, `gen_statem`) are agent archetypes: contracts a process implements, not classes it inherits from. Our agent shapes (narrow/role/composite) map directly: `gen_server` = narrow agent (request-response), `gen_statem` = role agent (state machine with loop), `supervisor` = composite agent.

**Process calculi: session types and choreography.**

CDD.md is a choreography — it describes the global interaction pattern from an omniscient view. Each role skill is the local endpoint projection. The rule "role skill must not contradict CDD.md" is the choreographic faithfulness condition: if projections are faithful, the system is deadlock-free.

Session types give agent signatures a formal reading. α's session type: `!implement . !self-cohere . !request-review . μX.(?RC . !fix . X + ?A . !closeout)` — send implementation, send self-coherence, send review-request, then loop (receive RC, send fix, repeat) or (receive approve, send closeout). This IS the α algorithm. The LANGUAGE-SPEC signature (`inputs`, `outputs`, `calls`) is already a session type — it just lacks the temporal structure (what comes first, what repeats).

Channel restriction (`(ν c) P`) maps to scoped artifacts. `.cdd/unreleased/{N}/` is a restricted channel — only agents in this cycle can read/write it. Scope IS the ν-binder.

**Rx/Reactive Streams: observable composition.**

Agents emit streams of artifacts, not single returns. α doesn't return once — it emits self-coherence, then review-request, then fixes, then closeout. Rx models this: agents are Observables, the artifact directory is a ReplaySubject (late subscribers see history), and composition operators (`merge`, `concat`, `zip`, `takeUntil`, `retry`) map to our `|||`, `>>`, `wait`, loop-condition, and `fix`.

Rx adds push/pull unification and cancellation — both gaps in the current model. An agent should be disposable: γ cancelling α because the issue was superseded is `dispose()`.

**Schlapbach (2026): process calculus for agentic protocols.**

Recent work formalizing MCP and SGD as π-calculus processes proves that protocol safety is equivalent to schema completeness. For CTB: signature completeness IS verifiable safety. The paper also formalizes capability confinement as a process invariant — the same structure as our "no upward mutation" rule.

**What CTB should borrow:**

1. From Erlang: supervision trees as the composition model for composite agents. Restart strategies as error-recovery vocabulary.
2. From process calculi: session types to give signatures temporal structure. Choreography/projection duality to formalize CDD.md vs role skills.
3. From Rx: observable streams as the agent emission model. Cancellation as a first-class operator.
4. From Schlapbach: signature completeness as provable safety. Capability confinement as a process invariant.

### 14.5 The unified model: agent = type instance + composition body

§13 said the agent is the primitive. §14.1–§14.3 named the composition operators. §14.4 surveyed the prior art. This section unifies them.

An agent is exactly two things:

1. **Type instance** — which SOUL (agent-type) it satisfies. The type defines: the loop shape (UIE-V), the invariants (laws), and the required bindings (what each instance must provide).
2. **Composition body** — how the agent combines subagent dispatches using the operators from §14.1.

```
-- Type declaration (SOUL template)
type CoherentAgent = {
  loop: [observe, identify, execute, verify],
  invariants: {honesty, engineering, scope, continuity},
  requires: {identity, observation_strategy, action_strategy, communication_style, memory_model}
}

-- Type instance + body (concrete agent)
agent sigma : CoherentAgent = {
  identity: Sigma,
  observation_strategy: evidence over urgency,
  action_strategy: smallest change that closes the named gap,
  communication_style: direct, terse, no sycophancy,
  memory_model: hub-based
}

-- Agents at different scales — same structure
agent cdd : CoherentAgent =
  γ.observe >> γ.select >> γ.issue >> (α ||| β) >>= γ.triage >> γ.close

agent α : CoherentAgent =
  load >> implement >> self-cohere >> request-review >> fix-loop >> closeout

agent review : CoherentAgent =
  read-diff >> check-invariants >> produce-verdict
```

The difference between a narrow agent and a composite agent is not kind — it is whether the body contains subagent dispatches (`>>`, `|||`, `>>=`) or is a single pipeline of pure steps. Both satisfy the same type. Both follow the same loop. Both are constrained by the same invariants.

This resolves every structural question at once:

| Question | Answer |
|---|---|
| What is a SOUL? | A type declaration: loop shape + invariants + required bindings |
| What is a concrete agent? | A type instance (bindings) + a composition body (operators over subagents) |
| How do agents compose? | Using `>>`, `>>=`, `\|\|\|`, `case`, `fix`, `wait` in the body |
| How do agents relate to types? | They satisfy a type. They do not inherit from one another. |
| Why no inheritance? | Types constrain; bodies compose. There is nothing to inherit. |
| What is the SOUL's continuity rule? | Types don't mutate at runtime. New version = new type = re-instantiate. |
| What makes an agent well-typed? | Its body respects the type's invariants across all subagent dispatches. |
| What is a narrow vs composite agent? | Body with no subagent dispatch vs body with subagent dispatch. Same type. |

### 14.6 The two kinds in the language

The language has exactly two kinds:

1. **`agent-type`** — a type declaration. Defines loop, invariants, required bindings. The SOUL template is the canonical instance. Not callable. Not dispatchable. Applied by being satisfied.

2. **`agent-module`** — a callable module. Has a signature (name, scope, visibility, inputs, outputs, calls). Has a body (composition of subagent dispatches). Declares which agent-type it satisfies.

Every concrete thing in the system is one of these two. Types constrain. Modules compose. Together they are the complete model.

The prior art maps cleanly:

| Prior art | agent-type | agent-module |
|---|---|---|
| Haskell | typeclass definition | typeclass instance + function definition |
| Erlang/OTP | behaviour specification (`-callback`) | module implementing the behaviour |
| Process calculi | typing context Γ | typed process |
| Rx | scheduler + error strategy | observable pipeline |

---

## Authority

This document is non-normative. It exists to preserve conceptual progress and to make spec revisions deliberate.

If this document and `LANGUAGE-SPEC.md` disagree on a rule, the spec governs.

If this document and `CTB-v4.0.0-VISION.md` disagree on strategy, the Vision governs.

If a claim here becomes load-bearing for implementers, it should be moved into the spec rather than relied on here.
