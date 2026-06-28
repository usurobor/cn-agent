# CTB Language Spec

**Version:** 0.1 (draft-normative)
**Date:** 2026-04-26
**Status:** Normative reference for the skill-module / invocation layer of CTB. Fully realized today in `cnos.cdd`, with `alpha/SKILL.md` as the first concrete `calls_dynamic` case; prescriptive for the wider package set as it migrates. Expected to evolve with the kernel and runtime.

---

## 0. Purpose and status

This document specifies the language-level model of CTB: what a skill is, what it declares, how it is named, scoped, invoked, composed, and bounded against effects.

The signature surface defined in §2 — frontmatter, scope, visibility, contracts, calls — is fully realized today in `cnos.cdd`. The cnos.core skills (`skill`, `write`, `design`, `compose`, `naturalize`, `audit`) declare only the older minimal frontmatter (`name`, `description`, `triggers`, sometimes `artifact_class` and `kata_surface`) and state their input/output contracts in prose inside the body. For those skills this spec is **prescriptive**: it states the form they should take when migrated. Migration is not a precondition for adopting the spec elsewhere.

The spec also closes the seams the Vision surfaces — loader model, pipeline semantics, contract grammar, scoped authority — that the practice has named but not yet formalized.

It is **not**:

- the strategy document (see `CTB-v4.0.0-VISION.md`)
- the rationale document (see `SEMANTICS-NOTES.md`)
- the kernel grammar of `.coh` terms (defined separately by the v1.x kernel spec)

Where this spec and the Vision disagree, this spec governs the language model. Where this spec and the kernel spec disagree on terms and equations, the kernel spec governs.

The spec uses **MUST**, **MUST NOT**, **SHOULD**, **MAY** in the IETF sense.

---

## 1. Core semantic unit

The unit of the language is the **skill**.

A skill is:

- a **module** — a self-contained file (or directory rooted at a `SKILL.md`) with a stable identity
- carrying a **callable signature** — a declared contract that other skills, runtimes, and operators can dispatch against
- transforming one **artifact** along one **axis**, deterministically

A skill is not a class. It has no private state, no inheritance, no implicit ambient behavior. It is a function with a contract, packaged as a module.

A directory containing a `SKILL.md` MAY contain sub-skill directories. The root `SKILL.md` is the module's entry; sub-skills are addressable as nested modules under the same package.

```
skills/
  cdd/
    SKILL.md            # root skill, public entry
    CDD.md              # canonical method (data, not a skill)
    alpha/SKILL.md      # sub-skill
    beta/SKILL.md       # sub-skill
    gamma/SKILL.md      # sub-skill
    issue/SKILL.md      # lifecycle sub-skill
    ...
```

---

## 2. Signature

Every skill MUST declare a signature in YAML frontmatter at the top of `SKILL.md`. The signature is the public face of the module.

### 2.1 Required fields

| Field | Type | Meaning |
|---|---|---|
| `name` | identifier | Unique within its package; used as the dispatch handle. |
| `description` | string | One-sentence statement of what the skill transforms. |
| `artifact_class` | enum | One of `skill`, `runbook`, `reference`, `deprecated`. |
| `kata_surface` | enum | One of `embedded`, `external`, `none`. `none` is permitted only for non-skill classes. |
| `governing_question` | string | The one question this skill answers. |

### 2.2 Invocation fields

| Field | Type | Meaning |
|---|---|---|
| `triggers` | list of strings | Situations under which the loader MAY select this skill. |
| `scope` | enum | One of `global`, `role-local`, `task-local`. See §3. |
| `visibility` | enum | One of `public`, `internal`. Defaults to `internal`. See §4. |

### 2.3 Contract fields

| Field | Type | Meaning |
|---|---|---|
| `inputs` | list | What the skill assumes on entry. |
| `outputs` | list | What the skill guarantees on exit. |
| `requires` | list | Preconditions outside the input set (environment, files present, prior state). |

### 2.4 Composition fields

| Field | Type | Meaning |
|---|---|---|
| `calls` | list of static call entries | Statically declared call edges. See §2.4.1. |
| `calls_dynamic` | list of dynamic call sources | Data-driven call sets resolved at runtime. See §2.4.2. |
| `runs_after` | list of skill names | Pipeline ordering: skills that MUST have completed first. |
| `runs_before` | list of skill names | Pipeline ordering: skills that MUST NOT have completed first. |
| `excludes` | list of skill names | Skills that MUST NOT load with this skill in the same activation. |

`runs_after`, `runs_before`, and `excludes` are reciprocal: if A declares `excludes: [B]`, B SHOULD declare `excludes: [A]`. A loader MUST treat a one-sided exclusion as an exclusion; it MAY warn about the missing reciprocal.

#### 2.4.1 `calls` — static call entries

Each entry in `calls` is one of:

- a bare path string: `design/SKILL.md` (sugar for `path: design/SKILL.md`)
- a mapping: `path: <skill-path>`

Every static target is a concrete skill the loader can resolve at activation time. Static targets form the inspectable call graph. A skill MUST NOT dispatch to a target that is not declared in either `calls` or `calls_dynamic`.

#### 2.4.2 `calls_dynamic` — data-driven call sets

A skill that selects sub-skills based on runtime data (e.g. issue contents, role context, configuration) declares the **source** of those targets, not the targets themselves. Each entry takes the form:

```yaml
calls_dynamic:
  - source: <data-path>
    constraint: <optional shape requirement>
```

`source` names where the target list comes from (a field on the input, a manifest entry, a runtime-computed set). `constraint` MAY narrow the permissible targets (e.g. "must be skills under `eng/`").

Example, modelled on the current α role:

```yaml
calls:
  - design/SKILL.md
  - plan/SKILL.md
calls_dynamic:
  - source: issue.tier3_skills
  - source: issue.tier2_bundles
```

The loader cannot enumerate dynamic targets ahead of time, but it can verify at dispatch that every actually-resolved target is reachable by *some* declared source. A dispatch to a skill that matches no static entry and no dynamic source is a violation.

### 2.5 Worked example

From `src/packages/cnos.cdd/skills/cdd/SKILL.md`:

```yaml
name: cdd
description: Coherence-Driven Development. Use for substantial changes ...
artifact_class: skill
kata_surface: embedded
governing_question: How do we evolve a system through a substantial change ...
visibility: public
triggers: [review, PR, release, issue, design, plan, assess, post-release, ship, tag]
scope: global
inputs:
  - substantial-change context
  - active role
  - issue or PR context
outputs:
  - canonical method loaded
  - active role skill loaded
  - required lifecycle sub-skills selected
requires:
  - CDD applies
  - CDD.md exists in this directory
calls:
  - alpha/SKILL.md
  - beta/SKILL.md
  - gamma/SKILL.md
  - issue/SKILL.md
  - design/SKILL.md
  - plan/SKILL.md
  - review/SKILL.md
  - release/SKILL.md
  - post-release/SKILL.md
```

This is the canonical shape of a public, scope-global root skill that calls a fixed set of sub-skills.

---

## 3. Scope

Scope is the lexical region in which a skill's authority and mutations are valid.

### 3.1 Scope values

| Scope | Meaning |
|---|---|
| `global` | The skill is active for the whole activation. Authority spans every nested invocation. |
| `role-local` | The skill is active only while a named role is bound (e.g. α, β, γ in CDD). |
| `task-local` | The skill is active only inside a single dispatched task and its descendants. |

### 3.2 Scope is first-class

Scope is part of the signature. A conformant skill MUST declare `scope`. The loader MUST refuse to dispatch a skill into a scope it did not declare to support.

For migration only, a loader MAY treat an omitted `scope` on a legacy skill as `task-local` and SHOULD warn. This is a compatibility behavior, not core semantics: a skill that omits `scope` is not conformant, and a stable release MUST NOT ship such a skill.

### 3.3 Scope and mutation

A skill at scope S MUST NOT mutate artifacts owned by an enclosing scope. Mutation rules are normative (see §10). Concretely, in CDD:

- α (role-local) owns implementation artifacts (branch, code, tests)
- β (role-local) owns review artifacts (CR, verdict)
- γ (role-local) owns cycle artifacts (issue triage, PRA)
- α MUST NOT write to β's or γ's artifacts; β MUST NOT write to α's or γ's; γ owns the cycle but MUST NOT replace α's or β's local artifacts after they are produced.

This is the same principle as functions not sharing mutable state. The loader and runtime SHOULD enforce it.

---

## 4. Visibility and dispatch

### 4.1 Visibility

| Value | Meaning |
|---|---|
| `public` | The skill is a valid external dispatch entrypoint for its package. Its triggers are loader-visible. |
| `internal` | The skill is callable only from a skill that declares it under `calls` or admits it via `calls_dynamic`. Its triggers are advisory and MUST NOT be used as public dispatch. |

`visibility` defaults to `internal`. A skill that omits `visibility` MUST NOT be dispatched from outside its package.

### 4.2 Dispatch model

The loader resolves a dispatch in three steps:

1. **Root dispatch.** An external trigger (operator, harness, peer) is matched against the union of `triggers` across all `public` skills in scope. Exactly one skill MUST match, or dispatch fails.
2. **Role binding.** If the matched skill declares `scope: role-local`, the dispatcher MUST also supply a role. If `scope: global`, any nested role binding is established by the matched skill itself (e.g. cdd binds α, β, or γ).
3. **Local expansion.** Once bound, the active skill MAY dispatch to skills it lists under `calls` or admits via `calls_dynamic`. Internal sub-skills are reachable only along these declared edges.

Triggers on `internal` skills are advisory: they document under what condition the parent skill should call the sub-skill. The runtime MUST NOT promote them to public entrypoints.

---

## 5. Invocation

Invocation is the runtime act of executing a skill's signature on concrete inputs.

### 5.1 Invocation primitives

```
operator
  └─ dispatch(skill = cdd, trigger = "release")     # root dispatch
       └─ bind(role = γ)                             # role binding
            └─ call(skill = post-release/SKILL.md)   # local expansion (static call)
            └─ call(skill = issue/SKILL.md)          # local expansion (static call)
```

### 5.2 Static and dynamic calls

There are two kinds of declared call edges. Both are first-class.

A **static call** is a target listed under `calls` (§2.4.1). The target is a concrete skill path. The loader can resolve and verify it at activation. Static calls form the inspectable dispatch graph.

A **dynamic call** is a target *set* declared under `calls_dynamic` (§2.4.2). The set is sourced from runtime data (issue contents, role context, configuration). The loader cannot enumerate the targets ahead of time; it can only verify, at dispatch, that each resolved target is admitted by some declared source.

Both kinds are reachable from the skill's body. Selection among static targets is permitted at runtime; resolution of a dynamic source happens at runtime by definition.

What is forbidden is dispatch to a target that matches no static entry and no dynamic source. The call graph stays inspectable: every edge the runtime takes was declared, even if its concrete endpoint was data.

> Restated: declared static edges may be selected dynamically; declared dynamic sources may resolve to any target the source legitimately produces; targets reachable by neither path are unreachable.

This split exists because practice needs both. The CDD lifecycle skills are declared statically (`design/SKILL.md`, `plan/SKILL.md`); the Tier 2 / Tier 3 skills the issue names are declared dynamically (`source: issue.tier3_skills`). Forcing either pattern into the other distorts the model.

### 5.3 Return values

Every invocation returns to its caller. Returns are data — close-out artifacts, effect plans, or witness terms. The caller owns the result and decides what to do with it.

Lower scope MUST NOT mutate higher scope. A sub-skill returns; it does not push.

### 5.4 Error handling

When a sub-skill cannot complete, it returns a witness (a structured failure term). The caller decides recovery: fix inline, re-dispatch, or escalate to its own caller.

A runtime MUST NOT silently swallow a sub-skill failure. It MUST either route the witness to the caller or surface it to the operator.

---

## 6. Composition

Composition is the only mechanism by which skills combine. Inheritance MUST NOT exist.

### 6.1 Composition rules

A skill MUST:

- own one axis (state it as a single sentence)
- declare entry and exit contracts (`inputs`, `outputs`)
- declare its trigger sharply enough to distinguish it from neighbors
- declare order against neighbors when order matters (`runs_after`, `runs_before`)
- assign each shared rule one owner and point to it from elsewhere
- declare exclusions reciprocally with the skill it cannot run with

A skill MUST NOT:

- inherit from another skill
- silently shadow a rule owned by another skill
- assume an order that is not declared
- expose two axes in one signature

### 6.2 Publish-or-compose gate

A new skill MUST justify its existence against the best alternative composition of existing skills. If the gap is not specific, the verdict is compose, not publish. (See `cnos.core/skills/compose/SKILL.md` §3.8.)

### 6.3 Example: composition without inheritance

`compose` is `design` applied to artifacts of class `skill`. It is not a child of `design`. It is a peer that:

- adopts `design`'s formula (one reason to change, one source of truth, truthful contracts)
- adds skill-specific concerns (classification, triggers, the publish-or-compose verdict)
- points to `design` as the owner of the shared rules it does not redefine

(See `docs/alpha/doctrine/SKILL-ARCHITECTURE.md` §6.)

---

## 7. Global aspects

A **global aspect** is a rule that applies across every skill in a package and has exactly one owner.

### 7.1 Why aspects exist

Practice surfaced this requirement: the rule "write large files section by section to disk" appeared in four CDD role/phase skills with slight wording drift. When one copy was missing, a real failure occurred. The rule was hoisted to `CDD.md §1.4` and the duplicates removed (Vision §8.5.1 E1).

This is a cross-cutting concern. It does not belong in any one skill, but it must be enforced for all of them.

### 7.2 Aspect declaration

A package MAY declare aspects in a package-level governance document. Each aspect MUST:

- have one owning document
- name the skills it applies to (or declare "all skills in this package")
- be referenced — not restated — by the skills it governs

A skill MUST NOT restate an aspect it inherits from a package-level owner. It MAY point to the aspect.

### 7.3 Authority

When two surfaces (an aspect and a local skill rule) disagree, the aspect's owning document governs unless the local rule is explicitly carved out. Authority MUST be explicit; readers and runtimes MUST NOT need to infer it.

---

## 8. Authority and mutation rules

### 8.1 One owner per fact

A stable fact (a rule, a precedence statement, a contract surface) MUST have exactly one owner. Other surfaces point to the owner. Restatement with drift is an authority violation.

### 8.2 Authority declarations are explicit

Every skill that depends on another skill's rule MUST name the owner directly. A skill that overrides a default behavior MUST declare the override and the owner it overrides.

### 8.3 Conflict resolution

When two skills disagree, the resolution rule MUST be machine-checkable. Examples from CDD:

- `CDD.md` is the only normative source for selection, lifecycle, dispatch contract, closeout obligations, and cycle-iteration rules.
- Role skills MAY add execution detail, gates, and evidence requirements for steps `CDD.md` already names — but MAY NOT redefine those steps.
- If `cdd/SKILL.md` and `CDD.md` disagree, `CDD.md` governs.

A package MUST publish its conflict-resolution rules. The runtime MUST refuse to load a package whose conflict rules are inconsistent.

### 8.4 No upward mutation

A skill at scope S MUST NOT write to a scope strictly enclosing S. Returns flow up; mutations do not.

---

## 9. Evaluation and the effect-plan boundary

### 9.1 Pure layer

A skill body is evaluated as a pure transformation from inputs to outputs. The output of a skill is data — typically an **effect plan** (a term that describes intended side effects). Evaluation MUST be deterministic: same inputs, same outputs, no hidden time, randomness, or ambient I/O.

(This is the kernel CTB property; the language layer inherits it.)

### 9.2 Effect-plan boundary

The boundary between the language and the world is the effect plan. The language returns a plan; a separate runtime executes it under explicit capability grants.

This boundary MUST NOT be crossed inside a skill body. A skill that performs I/O directly is not a CTB skill.

### 9.3 Plan composition

When skill A calls skill B, A's plan MAY embed B's plan as a subterm, sequence it, or branch on it. The composition is itself data.

### 9.4 Reproducibility

Two evaluations of the same skill on the same inputs MUST produce equal plans. A peer MUST be able to re-evaluate and verify equality. Plans SHOULD be content-addressable (plan hash).

---

## 10. Constraints on composition (summary of MUST/MUST NOT)

A composition is well-formed if and only if:

1. Every skill declares `name`, `description`, `artifact_class`, `kata_surface`, and `governing_question`.
2. Every skill declares `scope`. Public skills declare `visibility: public`.
3. Every skill declares `inputs` and `outputs`.
4. Every dispatch target is reachable through the caller's `calls` or `calls_dynamic`.
5. No skill inherits from another.
6. No two skills own the same rule.
7. No declared exclusion is one-sided in a checked package (warning) or is one-sided in a stable release (error).
8. No skill at scope S writes to an enclosing scope.
9. No internal skill is dispatched as a public entry.
10. Conflict-resolution authority is declared and acyclic.

A loader SHOULD verify (1)–(10) before activation. A linter MUST report violations of (1), (2), (4), (5), (6), (8), (9), (10).

---

## 11. Reserved vocabulary

The following frontmatter keys are reserved by this spec and MUST NOT be repurposed by packages:

```
name description artifact_class kata_surface governing_question
scope visibility triggers
inputs outputs requires
calls calls_dynamic runs_after runs_before excludes
```

Packages MAY define additional keys. Loaders MUST ignore unknown keys but SHOULD warn.

---

## 12. Examples from the current package set

### 12.1 A scoped, public, dispatching skill (fully realized)

`cnos.cdd/skills/cdd/SKILL.md` declares the spec's signature directly in frontmatter:

- `visibility: public`, `scope: global`
- `triggers` that route external dispatch to CDD
- `inputs`, `outputs`, `requires` as structured fields
- `calls` listing the static call graph (role skills + lifecycle sub-skills)
- normative authority deferred to a sibling document (`CDD.md`)

It is the only public entry to its package. Its sub-skills are internal even though they have triggers.

### 12.2 A role skill with static and dynamic calls (fully realized — first concrete `calls_dynamic` case)

`cnos.cdd/skills/cdd/alpha/SKILL.md` is the first skill to instantiate the `calls_dynamic` field. Its frontmatter declares:

- `visibility: internal`, `scope: role-local`
- `inputs`, `outputs`, `requires` as structured fields
- a static `calls` set (`design/SKILL.md`, `plan/SKILL.md`)
- a `calls_dynamic` set sourced from the issue:

```yaml
calls:
  - design/SKILL.md
  - plan/SKILL.md
calls_dynamic:
  - source: issue.tier2_bundles
  - source: issue.tier3_skills
```

This is the canonical case for `calls_dynamic`: the static set is the lifecycle skills α always reaches for; the dynamic set is the Tier 2 / Tier 3 skills the issue itself names. Any target the runtime resolves through `issue.tier2_bundles` or `issue.tier3_skills` is admitted; anything outside both static and dynamic edges is unreachable.

### 12.3 A pure transformation skill (prescriptive)

`cnos.core/skills/write/SKILL.md` is a single-axis transformation along prose coherence. Its body declares input and output contracts in §1 and §2 prose; its current frontmatter declares only `name`, `description`, and `triggers`.

When migrated, its signature should declare:

- `artifact_class: skill`, `kata_surface: embedded`
- `governing_question`
- `scope: task-local`
- `visibility: public` (it is composed by name from other packages)
- `inputs` / `outputs` lifted from the body's contract section
- no `calls` (it does not dispatch sub-skills)

It is composed by neighboring skills via `runs_after: [write]` (e.g. `naturalize`).

### 12.4 A peer that applies a foundation skill to a class (prescriptive)

`cnos.core/skills/compose/SKILL.md` is `design` applied to artifacts of class `skill`. Today its body states the input/output contract in §2.3 prose. Its current frontmatter declares only `name`, `description`, `artifact_class`, `kata_surface`, `governing_question`, and `triggers`.

When migrated, its signature should declare structured contracts and explicit authority pointers. It already demonstrates the no-inheritance rule in its body: it does not extend `design`; it composes with `design` and `write` and adds what it owns (classification, triggers, publish-or-compose verdict).

---

## 13. Open language questions (deferred)

These are language-level questions that this version of the spec does not close:

- Concrete grammar for `inputs` / `outputs` beyond a free-form list. Today they are prose; the spec needs a typed shape.
- The precise calling convention for return values: positional terms, named close-out fields, or both.
- Whether `requires` is a precondition the loader checks or only documentation.
- How aspects compose across packages when two packages declare overlapping aspects.
- Plan-hash canonicalization rules.

These are tracked in `SEMANTICS-NOTES.md` and revisited in each minor version of this spec.

---

## 14. Conformance

A loader, linter, or runtime conforms to CTB Language Spec v0.1 when:

- It accepts every signature shape defined in §2.
- It enforces §3 scope rules and §4 visibility rules.
- It verifies §10 constraints and reports violations.
- It refuses dispatch to undeclared targets (§5.2).
- It treats effect plans as data (§9) and never lets a skill body cross the effect boundary.

Conformance is checked against the conformance suite shipped with the spec (forthcoming).

---

## Authority

This document governs the language model for skill modules, signatures, scope, invocation, composition, and the effect-plan boundary.

- `CTB-v4.0.0-VISION.md` governs strategy, motivation, and roadmap.
- `SEMANTICS-NOTES.md` governs deeper conceptual rationale and design discussion.
- The CTB v1.x kernel spec governs term grammar, equation form, and evaluation of `.coh` programs.
- `docs/alpha/doctrine/SKILL-ARCHITECTURE.md` governs the architectural argument for composition over inheritance.

If this spec and the Vision disagree on language semantics, this spec governs. If this spec and the kernel spec disagree on terms or evaluation, the kernel spec governs. If this spec and `SKILL-ARCHITECTURE.md` disagree on the composition model, this spec governs (the doctrine is the argument; the spec is the rule).
