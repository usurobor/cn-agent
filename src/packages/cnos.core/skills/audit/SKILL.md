---
name: audit
description: Assess whether a package is coherent as a whole — does it have one job, do its parts serve that job, and do its parts agree with each other and with reality?
artifact_class: skill
kata_surface: embedded
governing_question: How do we judge whether a package is coherent — and produce the smallest convergence path to make it so?
triggers:
  - auditing a package for coherence
  - assessing whether a package's parts agree
  - preparing convergence work for a package
  - evaluating package health before a release or restructure
---

# Coherence Audit

## Core Principle

**Coherent audit judges a package top-down: purpose first, then structure, then agreement, then reality — and produces a convergence path that closes the governing gaps.**

A package is coherent when it has one job, every part serves that job, every part agrees with every other part about how the job is done, and what the parts claim matches what actually exists.

## Algorithm

1. Define — name the parts of a coherence audit, state the formula, and name the failure mode.
2. Unfold — execute the audit top-down: purpose → structure → agreement → reality → judgment.
3. Rules — imperative constraints on how audits are conducted.
4. Verify — check that the audit artifact is itself coherent.

---

## 1. Define

### 1.1. Identify the parts

A coherence audit has five parts:

- **Purpose** — the package's one job
- **Structure** — whether the shape serves the purpose
- **Agreement** — whether the parts say the same thing about shared surfaces
- **Reality** — whether claims match the filesystem and actual usage
- **Judgment** — the verdict and convergence path

Each part depends on the one before it. Purpose must be named before structure can be assessed. Structure must be understood before agreement findings have context. Agreement findings must exist before reality checks are meaningful. All four inform the judgment.

- ❌ Start with detail cross-checks and hope a picture emerges
- ✅ Start with purpose; let every subsequent finding serve the coherence question

### 1.2. Articulate how they fit

Purpose sets the standard. Structure reveals whether the package is organized to meet it. Agreement reveals whether the parts believe they're meeting it in the same way. Reality reveals whether what they claim is true. Judgment synthesizes all four into a verdict and a path.

The coherence formula: **a package is coherent when purpose, structure, agreement, and reality all align.**

### 1.3. Name the governing question

> This skill teaches how to judge whether a package is coherent and produce the smallest convergence path to make it so.

### 1.4. Name the failure mode

Audit fails through **bottom-up drift and ungrounded findings.**

Typical signs:

- detail findings without structural context — a list of nits with no governing judgment
- findings that can't be verified from their citations
- severity inflation or deflation — everything is D, or everything is A
- convergence themes too vague to execute ("fix naming")
- purpose never stated — the audit checks agreement without knowing what the parts should agree about
- the audit duplicates an existing audit or filed issue (CAP §1.2a)

---

## 2. Unfold

### 2.1. Purpose

State the package's one job in one sentence. What coherence does it provide that wouldn't exist without it?

- ❌ "cnos.cdd contains CDD skills and tools" (inventory, not purpose)
- ✅ "cnos.cdd makes triadic development cycles executable and verifiable"

If the sentence needs "and" between unrelated clauses, the package may have two jobs. That's finding zero.

Then check: does every top-level artifact serve this purpose?

- **Yes** — it belongs
- **No** — it should move or be removed
- **Unclear** — either the purpose statement is vague or the artifact's role is ambiguous

If more than a few answer "unclear," sharpen the purpose before continuing.

### 2.2. Structure

Does the package's shape make the purpose executable?

- Are skills grouped by the right axis?
- Are roles, phases, or layers clear from the structure?
- Does the structure create confusion about where something belongs?

Structural anti-patterns:

- **Two things with one name** — different artifacts use the same term for different concepts
- **One thing with two names** — the same concept has different names in different places
- **Orphaned artifacts** — files that nothing loads, references, or uses
- **Circular authority** — A says B governs, B says A governs
- **Layer smearing** — artifacts that mix concerns that should be separate
- **Missing parts** — things the purpose clearly requires but the package doesn't provide

### 2.3. Agreement

For each surface that appears in more than one artifact — a fact, policy, location, naming convention, role boundary, or load order — check: do the artifacts agree?

When they disagree, record a finding:

| ID | Sev | Theme | Finding | Evidence |
|----|-----|-------|---------|----------|

Severity:

- **D (Critical)** — agents will act wrongly. Behavior-changing contract drift.
- **C (Significant)** — agents can work around it but shouldn't have to.
- **B (Minor)** — low impact. Naming drift that doesn't cause confusion.
- **A (Informational)** — not wrong, could be tighter.

Every finding must cite specific text from both sides. A finding you can't verify from the citations is not a finding.

- ❌ "The naming is inconsistent across files"
- ✅ "F3 | C | Role ownership | beta/SKILL.md:22 lists 'post-release assessment' as β-owned; gamma/SKILL.md:35 assigns PRA to γ"

### 2.4. Reality

Verify artifact claims against the actual system:

- Do referenced paths exist?
- Do referenced commands work as described?
- Do artifact location claims match where artifacts actually live?
- Are there artifacts that exist but are never loaded or referenced?

- ❌ Trust what a skill says about the filesystem
- ✅ Check the filesystem

### 2.5. Judgment

State the verdict. Three possible outcomes:

- **Coherent** — purpose is clear, structure serves it, parts agree, claims match reality. Minor improvements possible but no convergence work needed.
- **Locally coherent, globally drifted** — individual artifacts are well-written but they disagree with each other or with the stated purpose. Convergence work needed.
- **Incoherent** — purpose is unclear or contested, structure doesn't serve it, significant disagreements. Restructuring needed.

Then consolidate findings into convergence themes. A good theme:

- Groups findings that share a root cause
- Has a clear "done" state
- Maps to a CDD cycle or a set of ACs in one issue

Order themes by: impact on agent behavior (D first), dependency (what unblocks what), then size (small wins first).

Decide the output path:

- All A/B → note in a close-out or adhoc thread, no issue needed
- C/D findings → produce a convergence issue
- Extensive (>15 C/D) → consider splitting by theme into multiple issues

---

## 3. Rules

### 3.1. State purpose before checking details

Do not cross-check artifacts before you can state the package's one job. If you can't state it, that's finding zero.

- ❌ "F1: beta/SKILL.md and CDD.md disagree on step 9" (without knowing what the package is for)
- ✅ "The package's job is X. Given that, here's where the parts disagree"

### 3.2. Work top-down

Purpose → structure → agreement → reality → judgment. Each level depends on the one above. Do not invert the order.

- ❌ Start with a line-by-line cross-check and derive structure from the findings
- ✅ Start with purpose, let structure contextualize the detail findings

### 3.3. Quote, don't summarize

Every finding must cite the specific text or path that constitutes the disagreement. Summarized findings are uncheckable.

- ❌ "The load order seems inconsistent"
- ✅ "beta/SKILL.md:35 loads post-release; gamma/SKILL.md:4 says γ owns PRA — load order contradicts ownership"

### 3.4. Check the filesystem

Do not trust artifact claims about paths, files, or commands. Verify them.

- ❌ "The skill references cn-cdd-verify, which presumably exists"
- ✅ `ls src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` → exists / does not exist

### 3.5. Grade honestly

D means agents will do the wrong thing. A means improvement opportunity. Don't inflate for urgency theater or deflate to avoid seeming harsh.

- ❌ Everything is C because it "feels significant"
- ✅ D only when the disagreement changes agent behavior; A when it's genuinely informational

### 3.6. Make themes actionable

Each convergence theme must have a clear done state. If you can't say what "fixed" looks like, the theme is too vague.

- ❌ "Fix naming drift" (which names? fixed how? done when?)
- ✅ "Artifact location alignment (F1, F4, F9): all five skills agree on POST-RELEASE-ASSESSMENT.md path. Done when: one canonical path, all references updated."

### 3.7. One finding per disagreement

Do not bundle unrelated disagreements. Each finding should be independently fixable and independently verifiable.

### 3.8. Respect the inverse gate

Before producing an audit, check: does an existing audit already cover this scope? Is there already a filed issue with the same findings? Don't duplicate. (CAP §1.2a)

---

## 4. Verify

### 4.1. Purpose check

Was the purpose stated before detail work began? Does the verdict answer the purpose question?

### 4.2. Coverage check

Did the audit cover all artifacts in scope? Were all five levels (purpose, structure, agreement, reality, judgment) addressed?

### 4.3. Evidence check

Does every finding cite specific text or paths? Were filesystem claims verified?

### 4.4. Severity check

Are D findings genuinely behavior-changing? Are A findings genuinely informational?

### 4.5. Theme check

Does every C/D finding belong to a convergence theme? Does every theme have a clear done state?

### 4.6. Self-coherence check

Does the audit itself follow the top-down principle it demands of the artifacts it audits?

---

## 5. Output

The audit produces one document:

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
For each theme: the gap, which findings, what "fixed" looks like.

## Convergence Path
Ordered themes with rationale for the order.
```

---

## 6. Kata

### Scenario

You are auditing cnos.core — the foundational skills package containing agent/, compose/, design/, naturalize/, skill/, write/, audit/, and ops/.

### Task

Judge whether the package is coherent as a whole.

### Governing skills

audit, skill, compose, design

### Inputs

All files under src/packages/cnos.core/skills/ and cn.package.json.

### Expected artifacts

A coherence audit document per §5.

### Verification

- Was purpose stated before detail work?
- Does the verdict match the evidence?
- Does every finding cite specific text?
- Were paths verified against the filesystem?
- Are convergence themes actionable?

### Common failures

- Starting with detail cross-checks before stating purpose
- Producing a finding list without a structural judgment
- Citing from memory instead of quoting
- Themes too vague to execute

### Reflection

- Does the audit tell you whether this package is healthy?
- Could someone execute the convergence path without re-deriving the analysis?
- Is the audit itself coherent top-down?
