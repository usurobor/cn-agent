# Skills Language — Design Evidence from Practice

Evidence collected from operating cnos skills that informs the design of the skills language.

---

## E1: Global aspects (cross-cutting constraints)

**Date:** 2026-04-24
**Source:** CDD skill-program maintenance, session with Sigma

### Observation

The CDD skill program had the same authoring constraint — "write large files section by section to disk" — duplicated in four places:

1. `CDD.md` α step 11 (α close-out)
2. `CDD.md` β step 10 (β close-out)
3. `gamma/SKILL.md` §2.10 (γ close-out)
4. `post-release/SKILL.md` Writing Process (PRA)

Each was worded slightly differently. When γ's version was missing, γ hit API timeouts writing a close-out because no instruction constrained the behavior.

### Fix in current skill-program

Hoisted to a single rule in `CDD.md §1.4` and removed the four per-step duplicates. The rule now says: any file >50 lines must be written section by section. Applies to all roles and all artifact types.

### What this tells us about the language

This is a **cross-cutting constraint** — it doesn't belong to any role, phase, or artifact type. It applies to any skill that produces files above a size threshold. In the current markdown-based skill format, the only way to express this is prose in a shared authority document, with manual deletion of duplicates.

A skills language should support **global aspects**: constraints declared once that apply across all skills matching a condition.

```
aspect large_file_authoring {
  applies_to: any skill where output.lines > 50
  constraint: write section by section, report after each section
  rationale: prevents context loss from API timeouts or session interruption
}
```

### Properties of a good aspect mechanism

1. **Declared once, enforced everywhere.** No duplication across skills.
2. **Condition-based application.** The aspect fires based on observable properties (file size, artifact type, role), not manual annotation on every skill.
3. **Composable.** Multiple aspects can apply to the same skill execution without conflict resolution headaches.
4. **Visible.** When a skill is loaded, the agent can see which aspects apply — aspects are not invisible magic.
5. **Overridable with justification.** A specific skill can opt out of an aspect if it states why. Silent opt-out is not allowed.

### Analogies in other systems

- **AOP (AspectJ, Spring):** cross-cutting concerns (logging, transactions) applied by pointcut matching
- **CSS:** global styles applied by selector matching, overridable by specificity
- **Middleware (Express, Rack):** pipeline stages that apply to all requests matching a pattern
- **Cargo/npm lint configs:** rules declared once, applied to all packages

The skills language version is closest to AOP but with the transparency requirement (property 4) that AOP typically lacks.

---

## E2: Authority hierarchy (which file wins)

**Date:** 2026-04-24
**Source:** CDD package audit finding 4.1

### Observation

CDD.md §1.4 says "When they diverge on role execution, the skill governs." But `cdd/SKILL.md` has a more precise conflict rule: CDD.md governs lifecycle order and artifact contracts; role skills govern local execution detail.

When two files disagree, agents need a mechanical answer for which one wins.

### What this tells us about the language

A skills language needs an explicit **authority hierarchy**:

```
authority_order:
  1. CDD.md          # lifecycle, selection, artifact contract
  2. role skill       # role-local execution detail
  3. phase skill      # phase-specific procedure
```

With a rule: a lower-authority skill cannot override lifecycle order, role ownership, or artifact obligations from a higher-authority file.

---

## E3: Visibility as a first-class property

**Date:** 2026-04-24
**Source:** CDD package audit finding 3.1, 3.2

### Observation

Root `cdd/SKILL.md` lacked `visibility: public`. Internal sub-skills had `visibility: internal` but no mechanism prevented the runtime from exposing them as public dispatch targets. The trigger overlap between root and sub-skills meant a bare keyword like "review" could resolve to an internal sub-skill.

### What this tells us about the language

- `visibility` must be a first-class, enforced property — not advisory metadata
- Internal skills are not dispatch targets from outside the package
- Trigger resolution must respect visibility: external dispatch enters through public skills only; internal triggers are advisory after the public entrypoint is loaded
- The language should enforce this at the dispatch/resolution layer, not rely on each skill author to remember the convention

---

## E4: Load order as dependency graph

**Date:** 2026-04-24
**Source:** CDD package audit §5

### Observation

Each skill declares what to load and in what order, but there's no formal graph — just prose. Finding 5.1 showed CDD.md §4.4 claiming "all CDD lifecycle skills are always loaded" while role skills only load what they need. Finding 5.2 showed CDD.md referencing gamma sections that don't exist.

### What this tells us about the language

- Load order should be a **declared dependency graph**, not prose instructions
- The language should detect: missing dependencies, cycles, references to non-existent sections
- Tier classification (always/phase-required/issue-specific) should be expressed in the dependency declaration, not in a separate prose rule that drifts from actual load behavior
