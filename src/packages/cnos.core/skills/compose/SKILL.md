---
name: compose
description: Decide when to publish a new skill, when to compose existing skills, and how to declare the contracts that make that composition coherent.
artifact_class: skill
kata_surface: embedded
governing_question: How do we compose skills without overlap, drift, or hidden order?
triggers:
  - designing a new skill
  - deciding whether behavior belongs in a neighbor skill
  - resolving overlap between skills
  - declaring skill order for a task
  - reviewing whether a proposal should be a composition instead
---

# Composition

## Core Principle

Coherent composition keeps axis, contract, trigger, order, and ownership aligned so each skill changes one thing, loads at the right time, hands off cleanly, and shares no silent source of truth.

This skill teaches how to decide whether behavior belongs in a new skill or in a composition of existing skills, then declare the contracts and order that make that decision executable.

This skill transforms a skill proposal or skill set along the axis of compositional fitness.

Composition fails through overlap and hidden order.

The skill skill governs whether an artifact is a skill and how skill artifacts are structured. The write skill governs prose coherence and shared house style. This skill adds only what compose adds.

## Algorithm

1. Define — classify the artifact, name the parts, state the formula, the governing question, and the failure mode.
2. Unfold — declare axis, contracts, trigger, order, ownership, exclusions, and the publish-or-compose verdict.
3. Rules — write imperative rules with stable IDs and contrastive ❌/✅ examples.
4. Verify — check that the skill composes cleanly and can be practiced.

---

## 1. Define

### 1.0. Classify the artifact first

Composition is itself a skill because it has a domain-specific formula, teaches judgment about skill boundaries and order, benefits from Define → Unfold → Rules, and exposes a kata.

Apply the same test to the artifact under review. Valid classes:

- skill — teaches repeatable judgment in a domain
- runbook — teaches ordered operational steps
- reference — provides lookup-oriented material
- deprecated — is superseded and should not remain live

Use this skill only for artifacts of class skill. Reclassify everything else before you continue.

- ❌ Treat a reference as if it were a skill
- ✅ Reclassify first, then compose only true skills

### 1.1. Identify the parts

Composition has six parts:

- skill — the unit being composed
- axis — what it changes
- contract — what it assumes and guarantees
- trigger — when it should load
- order — what runs before or after it
- owner — where each shared rule or stable fact lives

If one part is missing, the composition decision is vague. If two parts blur together, skills overlap.

- ❌ Talk about composition in generalities
- ✅ Name the axis, contract, trigger, order, and owner before deciding to publish

### 1.2. Articulate how they fit

Skills compose when each one owns one axis, declares entry and exit contracts, loads on a distinct trigger, runs in explicit order, and points shared rules to one owner.

If axis, contract, trigger, order, or ownership is missing, the composition overlaps, drifts, or misfires.

- ❌ "Coherent composition is combining related skills well"
- ✅ "Coherent composition keeps axis, contract, trigger, order, and ownership aligned so each skill changes one thing and hands off cleanly"

### 1.3. Name the governing question

> This skill teaches how to decide whether work belongs in one skill, several composed skills, or no skill at all, and how to declare the composition once chosen.

- ❌ One artifact teaching composition, release policy, and versioning at once
- ✅ One artifact teaching composition judgment and declaration

### 1.4. Name the failure mode

Composition fails through overlap and hidden order.

Typical signs:

- two skills own the same rule
- one skill assumes another ran first but does not say so
- two skills trigger on the same situation without a clear distinction
- a proposed skill duplicates a viable composition
- one skill changes two axes
- a shared fact has two homes

- ❌ Generic: "doing it wrong"
- ✅ Specific: "composition fails through overlap and hidden order"

---

## 2. Unfold

### 2.1. Keep one skill on one axis

A skill that changes two axes is two skills sharing a file. Split it.

- ❌ A single skill that changes structure, voice, and format
- ✅ Separate skills for structure, voice, and format, composed in order

### 2.2. Declare the axis in one sentence

State the axis with this sentence:

> This skill transforms <artifact> along the axis of <axis>.

If the sentence needs an "and," the skill probably has two axes.

- ❌ "This skill transforms a draft along the axis of coherence and tone"
- ✅ "This skill transforms a draft along the axis of structural coherence"

### 2.3. Declare entry and exit contracts

A contract says what must be true on entry and what will be true on exit. Use this shape:

- Input:
- Output:

This skill's own contract is:

- Input: a skill proposal, skill pair, or load plan with neighboring authorities available
- Output: a composition verdict with explicit axis, trigger, contract, order, and ownership decisions

- ❌ "Takes a document and makes it better"
- ✅ "Input: prose with a governing question. Output: sections align to sub-questions and stable facts have one home"

### 2.4. Declare the trigger

A trigger distinguishes this skill from its neighbors. If the trigger is vague, the skill will be double-loaded or loaded too late.

Use this shape:

> Use when...

- ❌ "Use when writing matters"
- ✅ "Use when deciding whether new behavior belongs in a new skill or in a composition of existing skills"

### 2.5. Declare pipeline order

State what runs before, what runs after, what the skill assumes, and what it leaves true. Use this shape:

- Runs after:
- Runs before:
- Assumes:
- Produces:

This skill itself composes like this:

- Runs after: skill
- Runs before: publication, bundling, or co-load decisions
- Assumes: candidate artifacts are already classified or are being classified now
- Produces: a composition verdict, declared contracts, and ownership decisions

- ❌ Assume a neighbor skill ran first without saying so
- ✅ Declare the handoff explicitly

### 2.6. Assign one owner to every shared rule

When two skills appear to share a rule, choose one owner and replace the other copy with a pointer.

Use write for prose rules. Use skill for skill-form rules. Add a new owner only when neither existing skill owns the rule.

- ❌ Restate "front-load the point" in three skills with slight wording drift
- ✅ Let write own the rule; point to it from neighboring skills

### 2.7. Declare exclusions on both sides

If two skills should not run together, both skills should say so.

- ❌ Discover the conflict only when outputs disagree
- ✅ Declare the exclusion in both artifacts

### 2.8. Prefer composition before publication

Before publishing a new skill, name the best alternative composition and the specific gap it leaves. Publish only when that gap is real and owned by no existing skill.

- ❌ Publish blog because it feels different from essay
- ✅ Compare essay + voice(casual) + write to the proposed skill, then publish only if a specific transformation is still missing

### 2.9. Expose a kata surface

Every true skill must be practice-ready. The kata appears below.

---

## 3. Rules

### 3.1. Classify before composing

Decide whether the artifact is a skill, runbook, reference, or deprecated artifact before you write composition logic.

- ❌ Add composition language to a lookup page
- ✅ Reclassify the artifact, then compose only true skills

### 3.2. State one axis

Write the axis in one sentence and reject any skill that needs two axes to explain itself.

- ❌ "Improves drafts in many ways"
- ✅ "Transforms a draft along the axis of structural coherence"

### 3.3. Declare both contracts

State what the skill needs on entry and what it guarantees on exit.

- ❌ "Takes a draft and improves it"
- ✅ "Input: prose with a governing question. Output: each section answers one sub-question"

### 3.4. State the trigger sharply

Write a trigger that distinguishes the skill from its neighbors.

- ❌ "Use when writing"
- ✅ "Use when revising expository prose for structural coherence"

### 3.5. Declare order explicitly

Name the handoff. Do not make neighboring skills infer it.

- ❌ Order inferred from file names or habit
- ✅ "Runs after write. Runs before naturalness."

### 3.6. Move shared rules to one owner

Give each shared rule or stable fact one home, then point to that home elsewhere.

- ❌ The same authority rule restated in multiple skills with drift
- ✅ One owner, named directly, with pointers from the rest

### 3.7. Declare exclusions in both skills

If A conflicts with B, A names B and B names A.

- ❌ Only one side declares the conflict
- ✅ Both skills name the exclusion

### 3.8. Name the alternative composition before publishing

A new skill must say which existing composition it replaces and what that composition cannot do.

- ❌ Publish a new skill because the topic feels distinct
- ✅ Name the missing transformation that no existing composition owns

### 3.9. Compose, do not extend

Never express skill relationships as inheritance.

- ❌ "`newsletter` extends `write`"
- ✅ "`newsletter` composes with write and adds newsletter-specific form rules"

### 3.10. Test the skill in isolation

A skill should show its effect on a sample artifact without needing the whole pipeline to make sense.

- ❌ "Its effect only appears after three other skills run"
- ✅ A before/after example shows the axis it changes

### 3.11. Keep authority explicit

Say which surface governs on disagreement.

- ❌ Let the reader infer whether skill or write governs
- ✅ "design governs system decomposition. skill governs skill form. write governs prose rules. compose applies design to artifacts of class skill, adding triggers, classification, and the publish-or-compose verdict."

---

## 4. Verify

### 4.1. Classification check

Is the artifact really a skill? If not, reclassify it before composing.

### 4.2. Parts check

Are axis, contract, trigger, order, and owner all named?

### 4.3. Overlap check

Does any rule, stable fact, or decision have two homes?

### 4.4. Trigger check

Does the trigger distinguish this skill from every neighbor, or would two skills fire on the same situation without a declared difference?

### 4.5. Pipeline check

Are runs-before, runs-after, assumptions, and products explicit?

### 4.6. Publish-or-compose check

Does the proposal name the best existing composition and the specific gap it leaves? If not, default to composition, not publication.

### 4.7. Isolation check

Can the skill's effect be shown on a sample artifact without the full pipeline?

### 4.8. Kata check

Does the kata exercise composition judgment rather than only within-skill editing?

---

## 5. Final Test

A skill composes well when:

- its class is explicit
- its axis is singular
- its entry and exit contracts are declared
- its trigger distinguishes it from neighboring skills
- its order is explicit
- each shared rule has one owner
- its existence is justified against a named alternative composition
- it can be tested in isolation
- its kata exercises the composition decision

If a proposed skill cannot name a gap beyond the best existing composition, compose instead of publishing.

---

## 6. Kata

### Scenario

You are reviewing a proposed new skill called blog. The author argues that blog posts need a casual register, looser pacing, and a stronger hook than essays.

### Task

Decide whether blog should be published as a new skill, composed from existing skills, or refused.

### Governing skills

- skill
- write
- compose

### Inputs

- the proposed blog skill draft
- the existing write skill
- a hypothetical essay skill with form rules
- a hypothetical voice skill that applies a voice profile

### Expected artifacts

- a short composition note naming what essay + voice(casual) + write already produces
- a short gap note naming what blog would add that the composition does not produce
- a verdict: publish, compose, or refuse
- if publish, a frontmatter draft plus axis, contract, trigger, and order
- if compose, a documented composition recipe
- if refuse, a note added to the proposal explaining why

### Verification

- the alternative composition is named explicitly
- the gap is specific, not atmospheric
- if the gap is empty or vague, the verdict is compose or refuse
- if the verdict is publish, the new skill passes the Final Test

### Common failures

- calling blog distinct because the topic feels distinct
- publishing blog while it silently overlaps with essay on openings or pacing
- refusing blog without documenting the composition recipe
- modeling blog as a child of essay instead of a peer that composes

### Reflection

The judgment under test is: when is a new skill warranted? The default should be composition. Publication must earn itself with a specific missing transformation.
