# PLAN-v3.22.0-eng-lane-clarity

## Implementation Plan for Engineering Lane Clarity

Status: Complete
Implements: the Engineering Lane Clarity design amendment
Purpose: Make the engineering skill lane legible at the top level and make L5 / L6 / L7 stable, local terms inside cnos.

---

## 0. Coherence Contract

### Gap

The engineering skill suite has grown into a real lane, but its top-level description is still too small and too stale. At the same time, the repo uses level language such as:

- L5
- L6
- L7

without one canonical local definition.

This creates three problems:

1. Discoverability drift
   The top-level engineering README no longer reflects the actual suite.

2. Level-language ambiguity
   "L5 / L6 / L7" are used in discussion, but not grounded in one stable repo-local reference.

3. Bundle-selection ambiguity
   There is no clear top-level guidance for which engineering skills to load for which class of work.

### Mode

MCA — revise top-level engineering documentation and add a normative levels doc.

### α / β / γ target

- α PATTERN: one explicit engineering lane map, one explicit levels definition
- β RELATION: align skill tree, engineering assessment essay, and design/review language
- γ EXIT: future skill work and engineering assessment can reference a stable shared vocabulary

### Smallest coherent intervention

Do not redesign the engineering tree. Do not rewrite the assessment essay.

Add and revise only:

- src/agent/skills/eng/README.md
- docs/gamma/ENGINEERING-LEVELS.md

Optionally add small cross-links from nearby docs.

---

## 1. Deliverables

### New

- docs/gamma/ENGINEERING-LEVELS.md

### Revised

- src/agent/skills/eng/README.md

### Optional cross-links

- docs/papers/ENGINEERING-LEVEL-ASSESSMENT.md
- other lane/index docs if they reference engineering levels or bundles

---

## 2. Core Decisions to Preserve

These are the design decisions this implementation must not drift from.

### 2.1 README is operational

eng/README.md is the lane map:

- what eng/ is for
- what skills exist
- which bundles are recommended
- which work contexts map to which bundles

### 2.2 Levels doc is normative

ENGINEERING-LEVELS.md is the shared vocabulary:

- what L5 means
- what L6 means
- what L7 means
- how those levels map to cnos work

### 2.3 Essay stays descriptive

ENGINEERING-LEVEL-ASSESSMENT.md remains a historical/descriptive artifact, not the canonical definition.

### 2.4 No tree redesign in this issue

This line clarifies the lane and the language. It does not rearrange directories or repackage skills.

---

## 3. Step Order

## Step 1 — Inventory the current engineering lane

### Goal

Capture the actual engineering skill set and ensure the README reflects the real tree, not a stale subset.

### Work

Inspect the current src/agent/skills/eng/ contents and classify each skill into one of these functional groups:

- local implementation
- design / system-shaping
- review / closeout
- tooling / operator surfaces
- authoring engineering knowledge

### Acceptance

The README will not omit currently live engineering skills.

---

## Step 2 — Rewrite eng/README.md

### Goal

Make the lane legible and useful at load time.

### Required contents

1. What `eng/` is for
2. Skill index
3. Recommended load bundles
4. How to choose a bundle
5. Short level mapping
6. Relationship to the assessment essay
7. Relationship to the new levels doc

### Bundle set to include

At minimum:

- coding bundle
- review bundle
- design bundle
- runtime / platform bundle
- tooling bundle
- writing bundle

### Acceptance

A reader can answer:

- what skills exist in the engineering lane
- what bundle to load for a given work type
- what L5/L6/L7 roughly mean before reading the full levels doc

---

## Step 3 — Add ENGINEERING-LEVELS.md

### Goal

Create a stable, local, normative definition of engineering levels.

### Required contents

1. Why the doc exists
2. Core principle
3. L5 definition
4. L6 definition
5. L7 definition
6. What the levels are not
7. Practical interpretation inside cnos
8. Mapping from engineering skills to levels
9. How to tell what level a diff is operating at
10. Relationship to the assessment essay
11. Summary

### Required shortest definitions

- L5 — local correctness
- L6 — system-safe execution
- L7 — system-shaping leverage

### Acceptance

A repo-internal reader can interpret "this is L6" or "this is L7" without relying on external folklore.

---

## Step 4 — Align the README and levels doc

### Goal

Make sure the lane map and the level vocabulary do not conflict.

### Checks

- Skills described in eng/README.md match actual tree contents
- L7 lane in the README is consistent with the levels doc
- Bundles in the README do not contradict the levels definitions
- The README does not quietly make the essay normative
- The levels doc does not pretend to be an assessment artifact

### Acceptance

The two docs can be read together without contradiction.

---

## Step 5 — Add cross-links and nearby doc updates

### Goal

Reduce ambiguity and make the new docs discoverable.

### Minimum links

- eng/README.md links to docs/gamma/ENGINEERING-LEVELS.md
- ENGINEERING-LEVELS.md links to the assessment essay
- assessment essay notes that it is descriptive while the levels doc is normative

### Optional nearby updates

If other docs currently use L5/L6/L7 language loosely, add links or brief notes so they can point to the new canonical definition.

### Acceptance

A reader can navigate from:

- engineering lane → levels doc
- levels doc → assessment essay

without ambiguity.

---

## Step 6 — Review against real usage

### Goal

Make sure the new wording matches how the repo actually uses engineering terms and bundles.

### Review prompts

- Does the lane map reflect real activation contexts?
- Are the bundles too broad or too narrow?
- Are L5/L6/L7 defined in a way that matches recent assessment language?
- Does the README understate or overstate any skill's role?
- Is the L7 lane framed as "system-shaping leverage" rather than vague seniority?

### Acceptance

The docs improve clarity without requiring the skill tree itself to change.

---

## 4. File Changes

### Create

- docs/gamma/ENGINEERING-LEVELS.md

### Edit

- src/agent/skills/eng/README.md
- optionally docs/papers/ENGINEERING-LEVEL-ASSESSMENT.md for a small descriptive-vs-normative cross-link

---

## 5. Test / Verification Plan

This is a doc/governance change, so the verification is mostly structural and semantic.

### Structural checks

- all referenced skill directories/files exist
- all links resolve
- no stale skill names are listed
- no broken placeholders render incorrectly

### Semantic checks

- README skill groupings reflect actual activation contexts
- levels doc defines L5/L6/L7 clearly and non-circularly
- essay and levels doc are clearly differentiated

### Review focus

- Does this make the engineering lane easier to load correctly?
- Does this make future references to "L6/L7" unambiguous?
- Does it avoid overclaiming or redesigning the tree?

---

## 6. Non-goals

This plan does not:

- rearrange the engineering skill tree
- rename skill directories
- collapse or split packages
- rewrite the assessment essay
- automate skill selection
- change the content of existing engineering skills

It only makes the lane and the level language explicit and coherent.

---

## 7. Risks

### Risk 1 — Overfitting the README to today's tree

Mitigation:
- keep grouping broad
- avoid over-detailed sub-taxonomies

### Risk 2 — Treating the levels doc like a career ladder

Mitigation:
- define levels as shape of engineering impact, not title or tenure

### Risk 3 — Contradiction with future skill reorgs

Mitigation:
- keep the lane map descriptive of current contents
- avoid binding too much policy into the README

---

## 8. Acceptance Criteria

This issue/plan is complete when:

1. src/agent/skills/eng/README.md accurately reflects the current engineering skill suite.
2. docs/gamma/ENGINEERING-LEVELS.md exists and defines L5/L6/L7 in stable cnos-local terms.
3. The README defines recommended engineering bundles by work context.
4. The levels doc clearly distinguishes:
   - L5 = local correctness
   - L6 = system-safe execution
   - L7 = system-shaping leverage
5. The README and the levels doc clearly distinguish:
   - operational lane map
   - normative vocabulary
   - descriptive assessment essay

---

## 9. Summary

This plan does not change engineering behavior directly. It changes the top-level legibility of the engineering lane and the shared meaning of engineering level language.

That is the smallest coherent step that lets future engineering work, assessments, and skill selection refer to the same map and the same vocabulary.
