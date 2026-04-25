---
name: audit
description: Assess whether a package is coherent as a whole — does it have one job, do its parts serve that job, and do its parts agree with each other and with reality? Use when evaluating package health or preparing convergence work.
artifact_class: skill
kata_surface: embedded
governing_question: Is this package coherent — does it have one clear purpose, and do all its parts serve and agree on that purpose?
triggers:
  - auditing a package for coherence
  - assessing whether a package's parts agree
  - preparing convergence work for a package
  - evaluating package health before a release or restructure
---

# Coherence Audit

## Core Principle

**A coherent package has one job. Every artifact in the package serves that job. Every artifact agrees with every other artifact about what that job is, how it's done, and who does what.**

An audit measures how close the package is to that state. The output is not a list of nits — it's a judgment about the package's overall coherence, supported by evidence where the parts disagree.

## Algorithm

1. **Purpose** — what is this package's one job? State it. If you can't, that's the first finding.
2. **Structure** — does the package's shape serve the purpose? Are there parts that don't belong, or missing parts that should?
3. **Agreement** — do the parts agree with each other? Where two artifacts touch the same surface, do they say the same thing?
4. **Reality** — does what the package claims match what actually exists?
5. **Judgment** — is this package coherent? What's the smallest convergence path to make it so?

---

## 1. Purpose

### 1.1. State the package's job

In one sentence: what does this package do? What coherence does it provide that wouldn't exist without it?

- ❌ "cnos.cdd contains CDD skills and tools" (inventory, not purpose)
- ✅ "cnos.cdd makes triadic development cycles executable and verifiable"
- ❌ "cnos.core contains foundational skills" (category, not purpose)
- ✅ "cnos.core teaches the judgment patterns that every other skill assumes"

If the sentence needs "and" between unrelated clauses, the package may have two jobs.

### 1.2. Check whether every part serves the purpose

For each top-level artifact or skill in the package, ask: does this serve the stated purpose? Three possible answers:

- **Yes** — it belongs
- **No** — it should move or be removed
- **Unclear** — the purpose statement is too vague, or the artifact's role is ambiguous

If more than a few artifacts answer "unclear," the purpose statement needs sharpening before continuing.

### 1.3. Check whether anything is missing

Given the purpose, are there obvious gaps? Things the package should provide but doesn't?

- ❌ Invent speculative additions
- ✅ Name gaps only when the purpose clearly requires them and their absence causes real incoherence

---

## 2. Structure

### 2.1. Does the shape serve the purpose?

Look at how the package is organized. Does the structure make the purpose executable?

- Are skills grouped by the right axis?
- Are roles, phases, or layers clear from the structure?
- Does the structure create confusion about where something belongs?

### 2.2. Are there structural anti-patterns?

Common signs of structural incoherence:

- **Two things with one name** — different artifacts use the same term for different concepts
- **One thing with two names** — the same concept is called different things in different places
- **Orphaned artifacts** — files that nothing loads, references, or uses
- **Circular authority** — A says B governs, B says A governs
- **Layer smearing** — artifacts that mix concerns that should be separate

---

## 3. Agreement

This is where detail matters — but in service of the whole. The question is not "does line 35 match line 72" but "do these artifacts agree about what this package is and how it works?"

### 3.1. Identify shared surfaces

A shared surface is any fact, policy, location, naming convention, role boundary, or load order that appears in more than one artifact.

### 3.2. Check agreement on each shared surface

For each shared surface, ask: do the artifacts that touch it say the same thing?

When they disagree, record:
- What they disagree about
- Which artifacts are involved
- How severe the disagreement is:
  - **D (Critical)** — agents will act wrongly. Behavior-changing contract drift.
  - **C (Significant)** — agents can work around it but shouldn't have to. Authority ambiguity, stale references to important surfaces.
  - **B (Minor)** — low impact. Naming drift that doesn't cause confusion.
  - **A (Informational)** — not wrong, could be tighter.

### 3.3. Quote, don't summarize

Every disagreement must cite specific text from both sides. A finding you can't verify from the citations is not a finding.

- ❌ "The naming is inconsistent across files"
- ✅ "beta/SKILL.md:3 says 'post-release assessment'; CDD.md:89 says 'post-release assessment, close the cycle' — these assign different scope to the same term"

---

## 4. Reality

### 4.1. Verify claims against the system

Where artifacts make claims about the filesystem, commands, paths, or behavior — check them.

- Do referenced paths exist?
- Do referenced skills/commands work as described?
- Do artifact location claims match where artifacts actually live?

- ❌ Trust what the skill says about the filesystem
- ✅ Check the filesystem

### 4.2. Verify against actual usage

If possible, check how the package is actually used:
- Which skills are actually loaded in practice?
- Which artifacts are actually referenced?
- Are there artifacts that exist but are never used?

---

## 5. Judgment

### 5.1. State the verdict

Is this package coherent? Three possible verdicts:

- **Coherent** — purpose is clear, structure serves it, parts agree, claims match reality. Minor improvements possible but no convergence work needed.
- **Locally coherent, globally drifted** — individual artifacts are well-written but they disagree with each other or with the stated purpose. Convergence work needed.
- **Incoherent** — purpose is unclear or contested, structure doesn't serve it, significant disagreements. Restructuring needed.

### 5.2. Name the convergence themes

Group the disagreements into themes — each theme represents a coherence gap that can be closed in one pass. A good theme:

- Has a clear "done" state
- Groups findings that share a root cause
- Maps to a CDD cycle or a set of ACs in one issue

### 5.3. Propose the convergence path

What's the smallest set of changes to make this package coherent? Order the themes by:

1. Impact on agent behavior (D findings first)
2. How many other fixes depend on this one
3. Size (prefer small wins that unblock larger work)

### 5.4. Decide the output path

- All A/B → note in a close-out or adhoc thread
- C/D findings → produce a convergence issue
- Extensive findings (>15 C/D) → consider splitting by theme into multiple issues

---

## 6. Output

The audit produces one markdown document:

```markdown
# Coherence Audit — <package>

## Purpose
One sentence: what this package does.

## Verdict
Coherent / Locally coherent, globally drifted / Incoherent

## Structure
Does the shape serve the purpose? What doesn't belong? What's missing?

## Findings
| ID | Sev | Theme | Finding | Evidence |
|----|-----|-------|---------|----------|

## Convergence Themes
For each theme: what's the gap, which findings belong to it, what does "fixed" look like.

## Convergence Path
Ordered list of themes to address, with rationale for the order.
```

---

## 7. Rules

### 7.1. Purpose first

Do not start cross-checking details before you can state the package's one job. If you can't state it, that's finding zero.

### 7.2. Top-down, then bottom-up

Start with purpose and structure. Only then check agreement at the detail level. Detail findings without structural context are noise.

### 7.3. Quote, don't summarize

Every finding cites specific text. Uncheckable findings are worthless.

### 7.4. Check the filesystem

Don't trust artifact claims about paths. Verify them.

### 7.5. Grade honestly

D means agents will act wrong. Don't inflate or deflate.

### 7.6. Themes must be actionable

Each convergence theme must have a clear done state. If you can't say what "fixed" looks like, the theme is too vague.

### 7.7. Respect the inverse gate (CAP §1.2a)

Before producing the audit, check: does an existing audit already cover this scope? Don't duplicate.

---

## 8. Kata

### Scenario

You are auditing cnos.core — the foundational skills package. It contains: agent/ (16 sub-skills), compose/, design/, naturalize/, skill/, write/, audit/, ops/, and a README.

### Task

Assess whether the package is coherent as a whole.

### Governing skills

- audit (this skill)
- skill (for classification checks)
- compose (for cross-skill coherence)
- design (for structural coherence)

### Inputs

- All files under src/packages/cnos.core/skills/
- cn.package.json for the package

### Expected artifacts

- A coherence audit document per §6

### Verification

- Was the purpose stated before detail work began?
- Does the verdict match the evidence?
- Does every finding cite specific text?
- Are convergence themes actionable with clear done states?

### Common failures

- Starting with detail cross-checks before understanding the purpose
- Producing a finding list without a structural judgment
- Grading everything the same severity
- Themes too vague to act on ("fix naming")

### Reflection

- Does the audit tell you whether this package is healthy?
- Could someone execute the convergence path without re-deriving the analysis?
