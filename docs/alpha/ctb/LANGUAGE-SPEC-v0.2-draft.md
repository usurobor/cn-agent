# CTB Language Spec

**Version:** 0.2 (draft-normative)
**Date:** 2026-04-28
**Status:** Normative reference for the agent-module / invocation layer of CTB. Supersedes v0.1 by unifying "skill" and "agent" into one primitive.

---

## 0. Purpose and status

This document specifies the language-level model of CTB: what an agent is, what it declares, how it is named, scoped, invoked, composed, and bounded against effects.

The core simplification over v0.1: there is one primitive — the **agent**. What v0.1 called a "skill" is a narrow-scope agent with a single axis and short lifetime. What practice calls a "role" is an agent with role-local scope and a loop. What practice calls a "triad" is three agents composing into one agent at the next scale. The composition model is the same at every level: dispatch with scoped authority.

It is **not**:

- the strategy document (see `CTB-v4.0.0-VISION.md`)
- the rationale document (see `SEMANTICS-NOTES.md`)
- the kernel grammar of `.coh` terms (defined separately by the v1.x kernel spec)

Where this spec and the Vision disagree, this spec governs the language model. Where this spec and the kernel spec disagree on terms and equations, the kernel spec governs.

The spec uses **MUST**, **MUST NOT**, **SHOULD**, **MAY** in the IETF sense.

---

## 1. Two kinds

The language has exactly two kinds:

### 1.1 Agent-type

An **agent-type** is a type declaration that constrains agents. It defines:

- a **loop** — the ordered operations every agent of this type must follow (e.g. observe → identify → execute → verify)
- **invariants** — named laws (MUST/MUST NOT) that hold across all dispatches
- **required bindings** — what each instance must provide (identity, strategy, style, memory model)

An agent-type is not callable. It is not dispatched. It is satisfied.

The SOUL template is the canonical agent-type. Concrete SOULs (e.g. cn-sigma's) are instances.

```yaml
kind: agent-type
name: coherent-agent
loop: [observe, identify, execute, verify]
invariants:
  honesty: [met means fully met, partial is partial]
  engineering: [one source of truth, MCA before MCI]
  scope: [no upward mutation, external action respects gates]
  continuity: [changes only through explicit proposal and confirmation]
requires_binding: [identity, observation_strategy, action_strategy, communication_style, memory_model]
```

### 1.2 Agent-module

An **agent-module** is a callable module. It is:

- a **self-contained directory** (rooted at a `SKILL.md`) with a stable identity
- carrying a **callable signature** — inputs, outputs, scope, visibility, dispatch edges
- carrying a **composition body** — how it combines subagent dispatches using operators (§6)
- declaring which **agent-type** it satisfies

Agent-modules vary in scope and lifetime, not in kind:

| Shape | Scope | Lifetime | Example |
|---|---|---|---|
| Narrow | `task-local` | Single invocation: inputs → outputs | `write`, `design`, `review` |
| Role | `role-local` | Bound while role is held; may loop | α, β, γ in CDD |
| Composite | `global` | Full activation; dispatches subagents | `cdd` root agent |

An agent-module is not a class. It has no inheritance. It satisfies a type and composes subagents.

```
agent cdd : CoherentAgent =
  γ.observe >> γ.select >> γ.issue >> (α ||| β) >>= γ.triage >> γ.close

agent review : CoherentAgent =
  read-diff >> check-invariants >> produce-verdict
```

### 1.3 Relationship

Types constrain. Modules compose. An agent-module declares `type: <agent-type>` in its signature. The type's invariants apply across all subagent dispatches within the module's body. A well-typed agent-module is one whose body respects its type's invariants.

### 1.4 Triadic carrier

Every agent-module MUST expose at least one triadic operational lens. The default agent lens is:

```
tri(orientation, intervention, witness)
```

where:

- **orientation** = the agent's observed boundary, identified gap, governing question, and active obligations
- **intervention** = the work, dispatches, transformations, or effect-plan synthesis performed
- **witness** = the close-out evidence, verdict, residual debt, or structured failure term

A narrow agent MAY expose the **skill lens**: `tri(input contract, transformation, witnessed output)`.

A composite agent MAY expose the **protocol lens**: `tri(roles/capabilities, interaction, close-outs)`.

The lens names are semantic annotations over triadic structure. They do not change the kernel grammar. Position labels are gauge choices (per TSC foundation); the carrier is invariant.

Composition preserves the triadic carrier: the output of one agent's witness slot feeds the next agent's orientation slot (`>>=`), or parallel agents produce independent witnesses that a join synchronizes (`|||` + `wait`). Closure requires an inspectable witness — an agent that terminates without filling the witness slot is a violation.

---

## 2. Signature

Every agent-module MUST declare a signature in YAML frontmatter at the top of its `SKILL.md`. Every agent-type MUST declare its `kind`, `name`, `loop`, `invariants`, and `requires_binding` in a package-recognized type declaration surface. The signature is the public face of the module; the type declaration is the constraint surface.

### 2.1 Required fields

| Field | Type | Meaning |
|---|---|---|
| `name` | identifier | Unique within its package; used as the dispatch handle. |
| `description` | string | One-sentence statement of what the agent does. |
| `governing_question` | string | The one question this agent answers. |
| `scope` | enum | One of `global`, `role-local`, `task-local`. See §3. |
| `type` | identifier | Which agent-type this module satisfies. See §1.1. |

### 2.2 Optional fields

| Field | Type | Meaning |
|---|---|---|
| `visibility` | enum | `public` or `internal`. Defaults to `internal`. See §4. |
| `triggers` | list of strings | Situations under which the loader MAY select this agent. |
| `inputs` | list | What the agent assumes on entry. |
| `outputs` | list | What the agent guarantees on exit. |
| `requires` | list | Preconditions outside the input set. |

### 2.3 Composition fields

| Field | Type | Meaning |
|---|---|---|
| `calls` | list of static call entries | Statically declared dispatch edges. See §2.3.1. |
| `calls_dynamic` | list of dynamic call sources | Data-driven dispatch sets resolved at runtime. See §2.3.2. |
| `runs_after` | list of agent names | Pipeline ordering: agents that MUST have completed first. |
| `runs_before` | list of agent names | Pipeline ordering: agents that MUST NOT have completed first. |
| `excludes` | list of agent names | Agents that MUST NOT load in the same activation. |

#### 2.3.1 `calls` — static dispatch entries

Each entry is a path to a subagent's `SKILL.md`. Every static target is a concrete agent the loader can resolve at activation time. An agent MUST NOT dispatch to a target not declared in `calls` or `calls_dynamic`.

#### 2.3.2 `calls_dynamic` — data-driven dispatch

A `calls_dynamic` entry declares the **source** of dispatch targets, not the targets themselves:

```yaml
calls_dynamic:
  - source: <data-path>
```

The loader cannot enumerate dynamic targets ahead of time but can verify at dispatch that every resolved target is admitted by a declared source.

---

## 3. Scope

Scope is the region in which an agent's authority and mutations are valid.

| Scope | Meaning |
|---|---|
| `global` | Active for the whole activation. Authority spans every nested dispatch. |
| `role-local` | Active only while a named role is bound. |
| `task-local` | Active only inside a single dispatched task and its descendants. |

Scope is part of the signature. A conformant agent MUST declare `scope`. The loader MUST refuse to dispatch an agent into a scope it did not declare.

**Scope and state:** A `task-local` agent is typically stateless (inputs → outputs). A `role-local` or `global` agent may maintain state across loop iterations within its scope. State MUST NOT outlive scope — when the scope exits, the agent's state exits.

**Scope and mutation:** An agent at scope S MUST NOT mutate artifacts owned by an enclosing scope. Returns flow up; mutations do not.

---

## 4. Visibility and dispatch

| Value | Meaning |
|---|---|
| `public` | Valid external dispatch entrypoint. Triggers are loader-visible. |
| `internal` | Callable only from an agent that declares it in `calls` or `calls_dynamic`. |

The loader resolves dispatch in three steps:

1. **Root dispatch.** An external trigger is matched against `public` agents' triggers. Exactly one MUST match.
2. **Role binding.** If the matched agent declares `scope: role-local`, the dispatcher MUST supply a role.
3. **Local expansion.** The active agent MAY dispatch to agents it declares in `calls` or `calls_dynamic`.

---

## 5. Invocation

Invocation is the runtime act of executing an agent's signature on concrete inputs.

```
operator
  └─ dispatch(agent = cdd, trigger = "release")      # root dispatch
       └─ bind(role = γ)                              # role binding
            └─ dispatch(agent = post-release)          # local expansion
            └─ dispatch(agent = issue)                 # local expansion
```

### 5.1 Static and dynamic dispatch

Both kinds are first-class. Static targets are declared in `calls`; dynamic targets are resolved from `calls_dynamic` sources. Dispatch to a target matching neither is forbidden.

### 5.2 Return values

Every invocation returns to its caller. Returns are data — close-out artifacts, effect plans, or witness terms. The caller owns the result.

### 5.3 Error handling

When a subagent cannot complete, it returns a witness (structured failure term). The caller decides recovery. A runtime MUST NOT silently swallow a failure.

---

## 6. Composition

Composition is the only mechanism by which agents combine. Inheritance MUST NOT exist. The composition body of an agent-module is an expression over subagent dispatches using the operators below.

### 6.1 Operators

| Operator | Name | Meaning |
|---|---|---|
| `>>` | sequence | Run A, then run B. |
| `>>=` | bind | Run A, feed its result to B. |
| `\|\|\|` | parallel | Run A and B concurrently with scoped isolation. |
| `case` | branch | Match on A's return value, dispatch accordingly. |
| `fix` | iterate | Repeat A until condition is met. |
| `wait` | join | Block until a parallel agent signals. |
| `try` | recover | Run A; if it fails, dispatch recovery. |

These are the complete core operators for v0.2. Every agent composition — from a narrow pipeline to a triadic cycle — is expressible as a combination of these operators. Packages MAY define derived operators as sugar if they lower to this core.

### 6.2 Examples

```
-- Narrow agent: pure pipeline, no subagent dispatch
agent review : CoherentAgent =
  read-diff >> check-invariants >> produce-verdict

-- Role agent: loop with branching
agent β : CoherentAgent =
  load >> wait(α.ready) >> review
  >> case verdict of
       RC → notify(α) >> wait(α.fixed) >> β    -- recursive
       A  → merge >> release >> closeout

-- Composite agent: parallel dispatch with scoped isolation
agent cdd : CoherentAgent =
  γ.observe >> γ.select >> γ.issue
  >> (α ||| β)
  >>= γ.triage >> γ.close
```

### 6.3 Rules

An agent MUST:

- own one axis (stated via `governing_question`)
- declare entry and exit contracts (`inputs`, `outputs`)
- declare its dispatch edges (`calls`, `calls_dynamic`)
- assign each shared rule one owner and point to it from elsewhere

An agent MUST NOT:

- inherit from another agent
- silently shadow a rule owned by another agent
- expose two axes in one signature
- dispatch to a target not declared in `calls` or `calls_dynamic`

### 6.4 Operator obligations

Parallel composition (`|||`) MUST declare either:
- a deterministic join, or
- a caller-visible set of independent close-outs.

A `fix` composition MUST declare either:
- a bounded stopping condition, or
- a failure/debt close-out when convergence is not reached.

A `try`/`recover` composition MUST preserve the failed witness. Recovery MUST NOT erase the failure surface.

These obligations connect to judgment doctrine: the agent must name the boundary it protects, the boundary it breaches, and the residual debt it does not call closure.

### 6.5 Composition dimensions

The operators distribute across three composition dimensions:

| Dimension | Operators | What composes |
|-----------|-----------|---------------|
| Horizontal (sequence, handoff) | `>>`, `>>=` | Witness chains: A's close-out feeds B's orientation |
| Vertical (boundary, parallel isolation) | `\|\|\|`, `case`, `wait`/`join` | Scoped operands: each runs in its own authority region |
| Deep (recurrence, repair, debt) | `fix`, `try`/`recover` | History: failure becomes witness, not silence; iteration preserves debt |

---

## 7. Global aspects

A **global aspect** is a rule that applies across every agent in a package and has exactly one owner.

A package MAY declare aspects in a package-level governance document. Each aspect MUST:

- have one owning document
- name the agents it applies to (or declare "all agents in this package")
- be referenced — not restated — by the agents it governs

The agent's **SOUL** is a global aspect: it declares invariants that apply across every agent in the activation. Skills (narrow agents) reference the SOUL's constraints; they do not restate them.

---

## 8. Authority and mutation rules

1. **One owner per fact.** A stable fact MUST have exactly one owner. Other surfaces point to the owner.
2. **Authority is explicit.** Every agent depending on another's rule MUST name the owner.
3. **Conflict resolution is machine-checkable.** When two agents disagree, the resolution rule MUST be published.
4. **No upward mutation.** An agent at scope S MUST NOT write to a scope enclosing S.

---

## 9. Evaluation and the effect-plan boundary

### 9.1 Pure layer

A narrow agent's body evaluates as a pure transformation from inputs to outputs. The output is data — typically an **effect plan**. Evaluation MUST be deterministic.

### 9.2 Looping agents

An agent with `role-local` or `global` scope MAY run a loop (e.g. UIE-V: Understand → Identify → Execute → Verify). Each iteration produces an effect plan. The loop itself is bounded by the agent's scope lifetime.

The effect-plan boundary still holds: each iteration returns a plan; a separate runtime executes it. The agent's loop is a sequence of pure evaluations, not a sequence of side effects.

### 9.3 Plan composition

When agent A dispatches agent B, A's plan MAY embed B's plan as a subterm. The composition is data.

### 9.4 Reproducibility

Two evaluations of the same agent on the same inputs MUST produce equal plans.

---

## 10. Constraints on composition

A composition is well-formed if and only if:

1. Every agent declares `name`, `description`, `governing_question`, and `scope`.
2. Public agents declare `visibility: public`.
3. Every dispatch target is reachable through the caller's `calls` or `calls_dynamic`.
4. No agent inherits from another.
5. No two agents own the same rule.
6. No agent at scope S writes to an enclosing scope.
7. No internal agent is dispatched as a public entry.
8. Conflict-resolution authority is declared and acyclic.

---

## 11. Reserved vocabulary

The following frontmatter keys are reserved:

```
name description governing_question scope visibility triggers
inputs outputs requires
calls calls_dynamic runs_after runs_before excludes
kind type loop invariants requires_binding
body closeout witness
```

The first three lines are module-signature fields. The fourth line (`kind` through `requires_binding`) is for agent-type declarations. The fifth line (`body`, `closeout`, `witness`) is reserved for future use by the triadic carrier model.

Packages MAY define additional keys. Loaders MUST ignore unknown keys but SHOULD warn.

---

## 12. Implementation status

This spec is prescriptive — it defines what the language SHOULD express. The kernel does not express all of it yet.

| Spec concept | Kernel status | Gap |
|---|---|---|
| Agent-type (`kind: agent-type`) | Not expressible | v1.x has terms, not kinds. Types can be encoded as terms but not enforced. |
| Agent-module signatures (§2) | Partially realized | YAML frontmatter exists in practice. No loader validates it. |
| Scope enforcement (§3) | Not enforced | Prose rules only. Runtime does not check scope boundaries. |
| Composition operators (§6) | `[Seq \| e1 \| e2]` only | `>>` encodable as `Seq`. `>>=`, `\|\|\|`, `case`, `fix`, `wait` have no kernel representation. |
| Type satisfaction (§1.3) | Not checkable | No mechanism to verify a module's body respects its type's invariants. |
| Effect-plan boundary (§9) | Correctly modeled | The kernel's term algebra IS the pure layer. Effect plans are terms. |
| Visibility/dispatch (§4) | Not enforced | Internal triggers can leak as public dispatch. No loader enforcement. |

**What exists today:**
- Kernel: `[L\|C\|R]`, atoms, `_`, equations over triadic terms
- Effect algebra: `[Seq \| e1 \| e2]`, `[WriteFile \| ... ]`, `[GitCommit \| ... ]`
- Signatures: YAML frontmatter on `SKILL.md` files (validated by convention, not tooling)
- Everything else: prose in markdown bodies

**Migration path** (per SEMANTICS-NOTES §9 and Vision §8):
- v1.x → prove the kernel runs, canonical formatter, one migrated skill
- v2.0 → module/import system, package format, linter (where agent-type and agent-module become first-class)
- v3.0 → capability annotations, scope enforcement, plan hashing
- v4.0 → full interop, signed artifacts, conformance suite

This spec targets v2.0. It is written now so that the migration has a defined destination.

---

## 13. Worked example: generic triad, specialized agents

This example demonstrates the core insight: the orchestration is shared; the specialization is in the bound subagents and witness rules. Both compose because both expose the same triadic carrier.

```
-- Type declaration: the coherent agent loop
agent-type coherent-agent =
  loop: [orient, intervene, witness]

-- Generic triad: any three-agent coordination
agent generic-triad : coherent-agent =
  orient >>= execute >>= verify >> close

-- Specialized: essay writing as a triadic cycle
agent write-essay : coherent-agent =
  understand-prompt
  >> outline
  >> draft
  >> (argument-review ||| style-review)
  >>= join-essay-evidence
  >> case verdict of
       accepted -> close
       repair   -> fix (revise >> re-review) until accepted else close-with-debt
       blocked  -> close-with-debt

-- Specialized: code production as a triadic cycle
agent write-code : coherent-agent =
  understand-request
  >> design
  >> implement
  >> (run-tests ||| typecheck ||| lint ||| static-review)
  >>= join-code-evidence
  >> case verdict of
       accepted -> close
       repair   -> fix (patch >> re-test) until accepted else close-with-debt
       blocked  -> close-with-debt
```

Both `write-essay` and `write-code` satisfy `coherent-agent`. Both use `|||` for parallel review/verification with scoped isolation. Both use `fix` for repair loops. Both require a witness (close-out) regardless of outcome — even `close-with-debt` fills the witness slot. The difference is the bound subagents and what counts as evidence. The triadic carrier is the same.

---

## 14. Migration from v0.1

| v0.1 concept | v0.2 equivalent |
|---|---|
| `skill` | `agent-module` (any scope) |
| SOUL (prose document) | `agent-type` declaration + instance |
| `artifact_class` | Removed from required fields. Packages MAY use it as a pragmatic label. |
| `kata_surface` | Removed from required fields. Packages MAY use it. |
| "A skill has no private state" | Narrow-scope agents have no persistent state; wider-scope agents may maintain state within scope lifetime. |
| "A skill is a function" | An agent-module is a type instance with a composition body. Narrow agents are single expressions. |
| "Composition is explicit invocation" | Composition is an expression using `>>`, `>>=`, `\|\|\|`, `case`, `fix`, `wait`, `try`. |
| Global aspect (§7) | The agent-type's invariants propagated through the dispatch tree. |

---

## Authority

This document governs the language model for agent modules, signatures, scope, invocation, composition, and the effect-plan boundary.

If this spec and the Vision disagree on language semantics, this spec governs. If this spec and the kernel spec disagree on terms or evaluation, the kernel spec governs.
