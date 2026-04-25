---
name: audit
description: Systematically verify that a package's skills, contracts, and artifacts are internally coherent and coherent with each other. Use when assessing package health, finding contract drift, or preparing a convergence issue.
artifact_class: skill
kata_surface: embedded
governing_question: How do we find every incoherence within and across a package's skills, contracts, and artifacts — and produce an executable convergence plan?
triggers:
  - auditing a package for coherence
  - assessing skill-program health
  - finding contract drift across skills
  - preparing a convergence issue from audit findings
  - reviewing whether a package's skills agree with each other
---

# Audit

## Core Principle

**Coherent audit finds every gap between what a package's artifacts say and what they mean — within each artifact, across artifacts, and against the system they govern.**

An audit is not a code review. A code review judges a proposed change against a contract. An audit judges the contracts themselves against each other and against reality.

An audit is not a refactor plan. A refactor plan assumes the target shape is known. An audit discovers the shape by measuring what disagrees.

## Algorithm

1. **Define** — name the package, its artifact inventory, its authority surfaces, and the failure modes you're scanning for.
2. **Inventory** — enumerate every artifact, classify it, and map its declared authorities and dependencies.
3. **Cross-check** — for each pair of artifacts that share a surface, verify agreement. For each artifact, verify internal coherence.
4. **Grade** — assign severity to each finding.
5. **Consolidate** — group findings into convergence themes and produce an executable output.

---

## 1. Define

### 1.1. Name the audit scope

State:
- the package being audited
- which artifact classes are in scope (skills, commands, docs, manifests, tests)
- which authority surfaces are in scope (role contracts, artifact locations, load orders, naming, frontmatter)

- ❌ "Audit the CDD package" (scope undefined)
- ✅ "Audit cnos.cdd: all skills under skills/cdd/, the cn-cdd-verify command, cn.package.json, and cross-references to cnos.core and cnos.eng"

### 1.2. Name the failure modes to scan for

Common audit failure modes:

- **Contract drift** — two artifacts disagree on the same fact (location, ownership, naming, policy)
- **Authority ambiguity** — a fact has no clear owner, or two artifacts both claim ownership
- **Stale reference** — an artifact points to something that no longer exists or has moved
- **Role confusion** — two roles claim the same step, or a step has no owner
- **Load-order gap** — a skill assumes another is loaded but doesn't declare it
- **Naming drift** — the same concept uses different names across artifacts
- **Frontmatter gap** — missing or incorrect frontmatter fields
- **Trigger overlap** — two skills fire on the same situation without distinction
- **Dead artifact** — an artifact exists but is never loaded or referenced

- ❌ Scan for "anything that looks wrong"
- ✅ Name which failure modes you're scanning for before reading any artifact

---

## 2. Inventory

### 2.1. Enumerate every artifact

List every file in scope. For each, record:

| Path | Artifact class | Declared authorities | Dependencies | Frontmatter status |
|------|---------------|---------------------|--------------|-------------------|

- ❌ Read artifacts in random order and note issues ad hoc
- ✅ Build the inventory first, then cross-check systematically

### 2.2. Classify each artifact

Apply the skill skill's classification test (skill §1.0):
- Is it a skill, runbook, reference, or deprecated?
- Does the declared class match the actual content?
- If misclassified, record it as a finding

### 2.3. Map authority surfaces

For each artifact, extract:
- what facts it claims to own (artifact locations, naming conventions, load orders, role boundaries)
- what facts it references from other artifacts
- what it assumes without citing

If two artifacts claim the same fact, that's a candidate finding.

---

## 3. Cross-check

### 3.1. Pairwise agreement

For every pair of artifacts that share a surface (reference the same location, role, step, naming convention, or policy):

- Do they agree on the fact?
- If they disagree, which one governs?
- If neither governs, that's an authority ambiguity finding

Work systematically. Do not rely on memory of what you read earlier.

- ❌ "I think these two agree" (from recall)
- ✅ Quote the specific text from each artifact and compare

### 3.2. Internal coherence

For each artifact:
- Does the frontmatter match the content?
- Does the artifact follow its own rules (skill §2.3)?
- Are all internal cross-references valid?
- Is the governing question singular (write §1.3)?
- If a skill, does it pass the skill skill's Final Test (skill §5)?

### 3.3. Reality check

Where possible, verify artifacts against the actual system:
- Do referenced paths exist?
- Do referenced commands work?
- Do declared load orders match actual loader behavior?
- Do artifact location claims match where artifacts actually live?

- ❌ Trust what the skill says about the filesystem
- ✅ Check the filesystem

---

## 4. Grade

### 4.1. Severity levels

Assign each finding a severity:

- **D (Critical)** — agents will produce wrong output or take wrong action. Contract drift that changes behavior.
- **C (Significant)** — agents can work around it but may not. Authority ambiguity, stale references to important surfaces.
- **B (Minor)** — cosmetic or low-impact. Naming drift that doesn't cause confusion, minor frontmatter gaps.
- **A (Informational)** — improvement opportunities. Not wrong, but could be clearer or tighter.

### 4.2. Finding format

Each finding must include:

| ID | Severity | Category | Finding | Surface 1 | Surface 2 | Evidence |
|----|----------|----------|---------|-----------|-----------|----------|

- **ID** — sequential within the audit (F1, F2, ...)
- **Category** — which failure mode from §1.2
- **Surface 1 / Surface 2** — the two artifacts that disagree (or one artifact + reality for internal/reality findings)
- **Evidence** — quoted text or concrete observation

- ❌ "The load order seems off"
- ✅ "F7 | C | Load-order gap | beta/SKILL.md step 5 loads post-release but gamma/SKILL.md §2.7 says γ owns PRA | beta/SKILL.md:35 | gamma/SKILL.md:99 | β loads a skill it no longer owns"

---

## 5. Consolidate

### 5.1. Group findings into convergence themes

Cluster related findings into themes — each theme becomes a candidate work item or AC in a convergence issue.

A good theme:
- groups findings that share a root cause or a shared surface
- can be fixed in one coherent pass
- has a clear "done" state

- ❌ One giant "fix everything" theme
- ❌ One finding per theme (too granular)
- ✅ "Artifact location drift (F1, F4, F9, F12) — four artifacts disagree on where POST-RELEASE-ASSESSMENT.md lives"

### 5.2. Produce the audit artifact

The audit output is a markdown document with:

1. **Scope** — what was audited and which failure modes were scanned
2. **Inventory** — the artifact table from §2.1
3. **Findings** — the finding table from §4.2, grouped by severity
4. **Convergence themes** — grouped findings with proposed ACs
5. **Recommended issue shape** — if the findings warrant a convergence issue, sketch its ACs, scope, and non-goals

### 5.3. Decide the output path

- If findings are all A/B → note in a close-out or adhoc thread, no issue needed
- If findings include C/D → produce a convergence issue per the issue skill
- If findings are extensive (>15 C/D findings) → consider splitting into multiple issues by subsystem or theme

---

## 6. Rules

### 6.1. Inventory before judgment

Build the full artifact table before recording any finding. Incomplete inventory produces incomplete audits.

### 6.2. Quote, don't summarize

Every finding must cite the specific text or path that constitutes the disagreement. Summarized findings are uncheckable.

- ❌ "The naming is inconsistent"
- ✅ "beta/SKILL.md:3 says 'post-release assessment'; gamma/SKILL.md:35 says 'PRA'; CDD.md:89 says 'post-release assessment, close the cycle'"

### 6.3. Check the filesystem

Do not trust artifact claims about paths, files, or commands. Verify them.

- ❌ "The skill says cn-cdd-verify exists, so it does"
- ✅ `ls src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify`

### 6.4. One finding per disagreement

Do not bundle unrelated disagreements into one finding. Each finding should be independently fixable.

### 6.5. Grade honestly

D means agents will do the wrong thing. Don't grade everything C to avoid seeming harsh, and don't grade everything D for urgency theater.

### 6.6. Themes must be actionable

Each convergence theme must have a clear done state. If you can't say what "fixed" looks like, the theme is too vague.

### 6.7. Respect the inverse gate (CAP §1.2a)

Before producing the audit artifact, check: does an existing audit already cover this scope? Is there already a filed issue with the same findings? Don't duplicate.

---

## 7. Verify

### 7.1. Coverage check

- Did the inventory include every artifact in scope?
- Were all declared failure modes scanned?

### 7.2. Evidence check

- Does every finding cite specific text or paths?
- Were filesystem claims verified?

### 7.3. Severity check

- Are D findings genuinely behavior-changing?
- Are A findings genuinely informational?

### 7.4. Theme check

- Does every C/D finding belong to a convergence theme?
- Does every theme have a clear done state?

### 7.5. Output check

- Does the audit artifact follow the template in §5.2?
- If a convergence issue is recommended, does the sketch satisfy the issue skill?

---

## 8. Kata

### Scenario

You are auditing cnos.core — the foundational skills package. It contains: agent/ (16 sub-skills), compose/, design/, naturalize/, skill/, write/, ops/, and a README.

### Task

Audit the package for internal coherence.

### Governing skills

- audit (this skill)
- skill (classification and structure)
- compose (cross-skill coherence)
- design (contract and boundary coherence)

### Inputs

- All files under src/packages/cnos.core/skills/
- cn.package.json for the package
- Any cross-references to other packages

### Expected artifacts

- Artifact inventory table
- Finding table with severity, category, and evidence
- Convergence themes with proposed ACs
- Recommended issue shape (if C/D findings exist)

### Verification

- Was the inventory complete before findings were recorded?
- Does every finding cite specific text?
- Were paths verified against the filesystem?
- Are convergence themes actionable with clear done states?

### Common failures

- Starting to write findings before completing the inventory
- Citing from memory instead of quoting the artifact
- Grading everything the same severity
- Producing themes that are too vague to act on ("fix naming")
- Duplicating an existing audit or filed issue

### Reflection

- Did the audit find incoherences that a casual read would miss?
- Are the convergence themes the right granularity for CDD cycles?
- Would an agent following these themes produce a more coherent package?
