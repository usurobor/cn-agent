---
name: naturalize
description: Remove surface signatures of machine-generated prose without changing meaning or structure. Use after write, before final publication.
artifact_class: skill
kata_surface: embedded
governing_question: How do we remove machine-like surface signatures without changing what the text says?
triggers:
  - editing AI-drafted prose for publication
  - de-branding a structurally coherent draft
  - matching a draft's surface register to its audience
  - checking whether a finished piece still sounds generated
---

# Naturalize

## Core Principle

A coherent draft can still sound generated. Surface naturalness is a separate axis from structural coherence: stock word choice, sentence-shape habits, synonym cycling, and register mismatch.

This skill transforms a structurally coherent prose draft along the axis of surface naturalness.

Naturalness fails through unnoticed signature.

write governs structure, evidence, abstraction control, and authority. This skill runs after write and touches surface only. If a fix changes the claim, the evidence, or the structure, revert it and send the issue back to write.

## Algorithm

1. Define — classify the artifact, name the parts, state the contract, the governing question, and the failure mode.
2. Unfold — declare the handoff from write, then run four passes in order: vocabulary, rhythm, repetition, register.
3. Rules — apply imperative rules with stable IDs and contrastive ❌/✅ examples.
4. Verify — check that the signatures are removed or justified and that meaning and structure are unchanged.

---

## 1. Define

### 1.0. Classify the artifact first

Naturalize applies to prose that already satisfies write's contract or has already been through write. It is not a drafting skill and it is not a structural revision skill.

Valid inputs:

- essay, post, memo, README, or review whose structure is already settled
- AI-drafted prose that needs surface cleanup before publication
- human-written prose whose wording or rhythm still sounds generated

Not valid inputs:

- a draft with unresolved structural problems
- a draft whose claims need evidence, sourcing, or scope correction
- code, tables, or lookup-oriented reference material

- ❌ Run naturalize on a draft whose sections still do the wrong jobs
- ✅ Run write first; use naturalize once the draft already works

### 1.1. Identify the parts

Surface signature has four families:

- vocabulary — stock AI words and helper verbs
- rhythm — sentence-shape habits such as trailing -ing tails, compulsive triads, and default contrast formulas
- repetition — synonym cycling where deliberate reuse would be cleaner
- register — mismatch between contractions, formality, and audience

- ❌ Scan the draft once for "anything that sounds AI"
- ✅ Scan one family at a time

### 1.2. Declare the contract

- Input: structurally coherent prose plus an intended register
- Output: the same meaning and structure, with listed surface signatures removed or justified and the register matched to the context

- ❌ "Make it sound more human"
- ✅ "Keep the meaning and structure; remove or justify the listed surface signatures"

### 1.3. Name the governing question

> This skill teaches how to remove machine-like surface signatures without changing what the text says.

- ❌ One skill teaching structure, voice capture, and signature removal at once
- ✅ One skill teaching surface signature removal only

### 1.4. Name the failure mode

Naturalness fails through unnoticed signature.

Typical signs:

- stock words survive because the writer reads past them
- sentences fall into patterned shapes the content did not earn
- synonyms cycle because repetition felt "bad"
- contractions and formality drift away from the audience

- ❌ "It sounds off"
- ✅ "The prose still carries surface signatures the reader can hear"

---

## 2. Unfold

### 2.1. Declare pipeline position

- Runs after: write
- Runs before: final publication, copyfit, or final proof
- Assumes: meaning, structure, evidence, and authority are already settled
- Produces: a surface-clean draft with the same meaning and structure

- ❌ Use naturalize to repair structure or evidence
- ✅ Use naturalize once only surface cleanup remains

### 2.2. Run the four passes in order

Run vocabulary, then rhythm, then repetition, then register. Finish each pass on the whole draft before starting the next.

- ❌ Fix a stock word, then a repetition issue, then a rhythm issue in the same read-through
- ✅ Finish the vocabulary pass before touching rhythm

### 2.3. Scan against the reference, not from memory

Load references/ai-tells.md for each pass. Memory misses the patterns the writer has already normalized.

- ❌ "I know what AI language looks like"
- ✅ Open the reference and scan against the active pass

### 2.4. Keep surface fixes on the surface

If a fix changes the claim, the proof, the scope, or the document shape, revert it. That work belongs to write.

- ❌ Rewrite the paragraph because the sentence shape felt patterned
- ✅ Adjust the sentence shape only; if the claim itself is weak, flag it for write

### 2.5. Keep legitimate uses when they are load-bearing

The reference is not a blacklist. A listed word or pattern may stay when it names the thing exactly or when the content genuinely earns the shape.

- ❌ Delete every listed word on sight
- ✅ Keep the pattern when its use is specific, necessary, and audience-fit

### 2.6. Expose a kata surface

Embedded below.

---

## 3. Rules

### 3.1. Run the passes in order

Finish one pass on the whole draft before starting the next.

- ❌ Mix vocabulary and rhythm edits in the same sweep
- ✅ Vocabulary pass complete, then rhythm pass

### 3.2. Scan against the reference

Use references/ai-tells.md for each pass. Do not rely on memory.

- ❌ "I already know the usual tells"
- ✅ Load the reference and check the draft against it

### 3.3. Prefer ordinary words and plain helper verbs

Replace stock signal words and inflated helper verbs when they are only carrying tone.

- ❌ "The report serves as a guide."
- ✅ "The report is a guide."

### 3.4. Break patterned sentence habits the content did not earn

Cut decorative -ing tails, forced triads, and default contrast formulas when a direct sentence says more.

- ❌ "The company expanded into new markets, changing the industry."
- ✅ "The company expanded into new markets and changed the industry."
- ✅ "The company expanded into new markets." (if the second clause was filler)

### 3.5. Reuse the right word instead of cycling synonyms

Repeat the accurate word when the referent is the same. Vary only when the meaning changes.

- ❌ "The system handles errors. The tool manages failures. The platform processes exceptions."
- ✅ "The system handles errors. It also handles failures and exceptions."

### 3.6. Match contractions and formality to register

Casual prose usually wants contractions. Formal prose often does not. Keep the piece on one register unless the shift is deliberate.

- ❌ "It is not a big deal." (in a casual personal blog)
- ✅ "It's not a big deal." (in a casual personal blog)
- ✅ "It is not a significant issue." (in a formal report)

### 3.7. Preserve meaning and structure

If an edit changes the claim or reorganizes the prose, revert it. Naturalize is not allowed to rewrite the argument.

- ❌ Delete a sentence because it sounds patterned and lose the reason it was there
- ✅ Keep the reason, change only the surface shape

### 3.8. Refer content problems back to write

When the real problem is structure, evidence, scope, or authority, stop editing and hand the draft back to write.

- ❌ Remove a vague source phrase and leave the claim unsupported
- ✅ Mark the issue for write, then continue the surface pass

---

## 4. Verify

### 4.1. Pattern check

For each family, any listed pattern that remains must have a reason to remain.

### 4.2. Meaning check

Compare input and output sentence by sentence. If meaning changed, revert the edit.

### 4.3. Structure check

Section and paragraph boundaries should be unchanged. If they moved, the skill exceeded its axis.

### 4.4. Register check

State the target register plainly — casual, neutral, or formal — then confirm contractions, vocabulary, and sentence pressure match it.

### 4.5. Read-aloud check

Read the output aloud. If it sounds recited rather than said, run the relevant pass again.

---

## 5. Final Test

A draft has been naturalized when:

- stock vocabulary and helper-verb tells are removed or justified
- rhythm patterns are removed or justified
- synonym cycling is gone
- the register fits the audience
- meaning is unchanged
- structure is unchanged

If meaning changed, revert and send the problem to write. If structure changed, revert and send the problem to write. If a listed pattern remains without a reason, run the relevant pass again.

---

## 6. Kata

### Scenario

A 600-word personal blog post was drafted from the writer's notes with an LLM. It has already passed write: the post has one job, the sections answer that job, and the argument is settled. It still carries surface signatures: two stock AI words, one inflated helper verb, a forced three-item list whose content is really two things, one decorative -ing tail, synonym cycling around the same referent, and a fully formal register in a casual venue.

### Task

Run the four passes in order. Produce a draft that keeps the same meaning and structure while removing or justifying the surface signatures.

### Governing skills

- skill
- write
- naturalize

### Inputs

- the structurally coherent draft
- references/ai-tells.md
- the target register: casual personal blog

### Expected artifacts

- the edited draft
- a short change log by pass
- any notes handed back to write if a surface fix would change meaning or structure

### Verification

- every change maps to one of the four passes
- no sentence means something different after the edit
- section and paragraph boundaries are unchanged
- the register fits the casual venue
- the read-aloud no longer sounds generated

### Common failures

- running naturalize before write
- fixing a content problem as if it were a surface problem
- replacing one stock signal word with a different stock signal word
- cycling synonyms twice instead of reusing the right word
- changing structure while pretending only the surface changed

### Reflection

The judgment under test is: when is a pattern a tell, and when is it simply the right word or sentence shape? The reference lists candidates. The skill decides which instances stay.
