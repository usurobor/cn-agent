# KATA-C2

## Browser Capability vs App Ecosystem

**Kata ID:** C2-browser-capability-and-ecosystem-boundary
**Family:** Architecture & leverage
**Difficulty:** High
**Failure Family:** Two related but distinct architectural layers are at risk of being collapsed into one solution.
**Expected Governing Skills (hidden):** design, architecture-evolution, process-economics

---

## Public Task Statement

A new request has arrived with two demands at once:

1. Agents need a **browser capability family:**
   - page fetch
   - DOM query
   - screenshot
   - maybe later form fill / click

2. cnos should eventually support something like an **app store / marketplace:**
   - install browser tooling as one bundle
   - later install other bundles the same way
   - keep the core minimal

The question is:

> Should browser capability and marketplace/app support be designed in one architecture doc / subsystem, or should they be separated?

Your task:

1. choose the most coherent architecture direction
2. explain the challenged assumption
3. compare at least three boundary moves
4. state leverage, negative leverage, migration
5. price any added process/complexity

Do not just pick what is easiest today.

---

## Public Context

### Repeated pressure

Past requests have already exposed a pattern:

- more capability families want to enter the runtime
- there is growing pressure for installation/distribution layers above runtime
- "just add it to core" no longer scales

### Constraints

- wake-up must remain local
- runtime remains the policy decision point
- browser capability may be risky and should be isolatable
- future operator-facing install experience should feel like:
  - apps / bundles / marketplace
- but runtime execution must stay coherent and minimal

### Candidate directions

#### Option A — One combined architecture

Create one big "browser + app store + extensions" subsystem.

#### Option B — Browser as runtime extension, marketplace later

Design browser capability as a runtime extension family, and define marketplace/registry as a separate higher layer later.

#### Option C — Core browser subsystem + later app store

Put browser directly into trusted core, but keep app-store concerns separate.

---

## Public Deliverable

Return:

1. Challenged assumption
2. Boundary moves considered (at least 3)
3. Alternatives comparison
4. Selected move
5. Invariants
6. Leverage
7. Negative leverage
8. Migration
9. Process cost / complexity pricing

---

## Hidden Evaluator Rubric

### α

Strong answer should:

- identify that this is really two related layers:
  - runtime capability growth
  - ecosystem/distribution growth
- avoid collapsing them into one mega-subsystem
- explicitly challenge an assumption like:
  - "product distribution and runtime execution should be designed as one layer"
  - or "all useful capability families belong in core"

Weak answer:

- chooses the combined mega-subsystem without naming the layering failure
- or chooses core accretion with no future leverage

### β

Strong answer should align with the already-converged direction:

- browser capability belongs in the runtime extension layer
- app-store / marketplace belongs in a companion layer above runtime
- bundles/apps are not runtime primitives
- runtime contract should surface installed truth, not registry possibilities
- wake-time remains local

Weak answer:

- treats marketplace as part of the runtime extension loading mechanism
- treats app/bundle as a new runtime primitive
- blurs runtime truth and registry truth

### γ

Strong answer should:

- give a phased migration, for example:
  - phase 1: runtime extension architecture proves browser family
  - phase 2: runtime contract / doctor / traceability
  - phase 3: registry / bundles / app ecosystem
- price the added process/complexity
- state non-goals

Weak answer:

- tries to solve everything in one release
- gives no migration or process-cost story

### Hidden Checks

The strongest answer will likely choose Option B and include:

- browser family as runtime extensions
- subprocess isolation by default
- marketplace/registry as a separate companion architecture
- bundle/app as a user-facing install composition, not runtime primitive
- local wake-time truth preserved
- leverage and negative leverage both named

### Transfer Target

The evaluator should look for transfer from C1:

- explicit challenged assumption
- real boundary comparison
- layer separation
- leverage / negative leverage
- migration and process-cost reasoning
