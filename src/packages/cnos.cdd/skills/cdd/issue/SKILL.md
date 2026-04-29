---
name: issue
description: Write an executable issue pack that names the incoherence, source-of-truth, constraints, proof obligations, required skills, and implementation guidance without creating ambiguity or false closure.
artifact_class: runbook
kata_surface: embedded
governing_question: How does γ turn a selected gap into an issue pack that α and β can execute, verify, and close without hidden assumptions?
visibility: internal
parent: cdd
triggers:
  - issue
scope: task-local
inputs:
  - selected gap
  - mode
  - active design constraints
  - Tier 3 skills
  - affected surfaces
  - related artifacts
  - current implementation status
outputs:
  - executable issue pack
requires:
  - γ completed observe/select
  - canonical CDD.md loaded
calls: []
---

# Issue

## Core Principle

A coherent issue is an executable work contract. It names:

1. the incoherence;
2. the impact;
3. the status of the affected system;
4. the source of truth;
5. the implementation boundary;
6. the acceptance criteria;
7. the proof or rejection mechanism;
8. the non-goals.

A coherent issue lets an engineer who was not in the conversation act without asking clarifying questions.

Failure modes:

- **ambiguity** — engineer must ask what the issue means;
- **overclaiming** — issue states target/draft/spec behavior as if it already ships;
- **false closure** — issue has good prose but no proof or rejection mechanism;
- **source drift** — issue cites the wrong canonical doc, path, branch, or runtime surface;
- **scope leak** — non-goals reappear inside acceptance criteria or work steps;
- **exception theater** — exceptions hide the rule instead of documenting temporary debt.

Operational shape:

```text
Issue = tri(gap, execution boundary, verification surface)
```

The problem names the gap.
The scope and constraints name the execution boundary.
The acceptance criteria and proof plan name the verification surface.

---

## 1. Define

### 1.1 Identify the parts

Every substantial issue MUST include:

- title;
- priority or severity;
- problem statement;
- impact;
- current status;
- source-of-truth / evidence links;
- scope;
- non-goals;
- acceptance criteria;
- proof or validation plan;
- Tier 3 skills and active constraints;
- related artifacts;
- implementation notes, if needed;
- success / closure condition.

For small issues, some sections may be short, but the issue must still be executable.

- ❌ `Fix CTB validation.`
- ✅ `CTB v0.1 frontmatter fields are specified but not mechanically checked. Add a CUE schema and CI job that reject missing hard-gate fields while allowing field-specific migration exceptions for legacy skills.`

### 1.2 State the incoherence

Problem statement format:

- **What exists:**
- **What is expected:**
- **Where they diverge:**

Keep the problem statement to 3–5 lines. Link to deeper evidence instead of reproducing the whole design discussion.

- ❌ `The package system is inconsistent and needs cleanup.`
- ✅ `SKILL.md frontmatter has a CTB v0.1 schema in LANGUAGE-SPEC.md, but no CI gate validates the 54 package skill files. Missing fields can drift indefinitely. Add structural validation before CTB v0.2 promotion.`

### 1.3 Name the impact

Impact answers:

- Who or what is affected?
- What breaks or drifts?
- Why now?
- What is blocked if this remains unresolved?

- ❌ `This would be nice to have.`
- ✅ `Without this gate, CTB frontmatter conformance depends on reviewer memory. v0.2 promotion and future ctb-check work inherit unchecked v0.1 drift.`

### 1.4 Name the status truth

Every issue touching specs, runtime, CI, language, protocol, or docs MUST distinguish:

- **shipped:** implemented and enforced today
- **current spec:** normative document currently governing conformance
- **draft target:** proposed or draft spec not yet promoted
- **planned:** accepted direction, not implemented
- **not in scope:** explicitly excluded
- **unknown:** not verified yet

Rules:

- Do not describe draft behavior as runtime behavior.
- Do not describe a planned checker as shipped enforcement.
- Do not describe a spec as normative if it is explicitly a draft.
- If the issue depends on an unverified fact, say so.

- ❌ `CTB enforces witnessed close-outs.`
- ✅ `CTB v0.2 draft defines witnessed close-outs. The shipped runtime does not yet enforce them.`

### 1.5 Name the failure mode

Use concrete failure modes, not generic concern. Examples:

- source drift
- witness theater
- unbounded repair
- broken link / stale path
- scope leak
- authority drift
- false runtime claim
- exception blanket
- CI blind spot
- operator-visible projection missing

- ❌ `This could be confusing.`
- ✅ `Failure mode: witness theater. A close-out can look complete because all fields are present while no checker verifies that failed evidence was preserved.`

---

## 2. Unfold

### 2.1 Source-of-truth table

For substantial issues, include a source map.

Format:

```markdown
## Source of truth

| Claim / surface | Canonical source | Status | Notes |
|---|---|---|---|
| CTB v0.1 frontmatter fields | docs/alpha/ctb/LANGUAGE-SPEC.md §2, §11 | current spec | governs current skill frontmatter |
| CTB v0.2 witnessed close-outs | docs/alpha/ctb/LANGUAGE-SPEC-v0.2-draft.md §5.4–§5.6 | draft target | not runtime-enforced |
| CUE install in CI | .github/workflows/coherence.yml | implementation surface | add job I5 |
```

Rules:

- Every load-bearing claim must link to a path, issue, PR, or upstream document.
- If a claim depends on another repo, include the repo and path.
- Prefer canonical docs over secondary summaries.
- If a link points to a branch or draft, mark that status.

- ❌ `The CTB docs say this somewhere.`
- ✅ `docs/alpha/ctb/LANGUAGE-SPEC.md §2.1 lists required skill frontmatter fields.`

### 2.2 Scope and non-goals

Substantial issues MUST include:

```markdown
## Scope

In scope:
- ...

Out of scope:
- ...

Deferred:
- ...

Blocked by:
- ...
```

Rules:

- Non-goals are mandatory for issues with 3+ ACs.
- Every noun in ACs and implementation notes must be in scope.
- If a non-goal noun appears in an AC, fix the issue before filing.
- If a deferred item is necessary for success, it is not deferred; it is a dependency.

- ❌ `Non-goal: implement runtime enforcement. AC4: runtime enforces the schema.`
- ✅ `Non-goal: runtime enforcement. AC4: CI rejects malformed frontmatter using scripts/check-skill-frontmatter.sh.`

### 2.3 Constraint strata

When an issue defines or enforces rules, stratify them. Use:

- **hard gate:** failure if missing or violated; no exception
- **exception-backed:** failure unless explicitly listed with reason
- **optional/defaulted:** may be omitted; default or absence semantics defined
- **validated-if-present:** not required, but checked when declared
- **ignored/deferred:** explicitly not checked in this issue

Example:

```markdown
## Constraint strata

Hard gate:
- name
- description
- governing_question
- triggers
- scope

Exception-backed:
- artifact_class
- kata_surface
- inputs
- outputs

Optional/defaulted:
- visibility, default internal

Validated if present:
- calls
- calls_dynamic
- requires
- parent

Ignored in v0:
- CTB v0.2 witnessed close-out semantics
```

Rules:

- Hard-gate fields cannot appear in `allowed_missing`.
- Exception-backed fields require field-specific exceptions.
- Optional/defaulted fields must name the default.
- Validated-if-present fields must name what is checked.
- Deferred fields must be explicitly marked as not enforced.

### 2.4 Exceptions ledger

Exceptions are debt, not permission to ignore the rule. Exception entries MUST be specific.

Preferred shape:

```json
[
  {
    "path": "src/packages/cnos.core/skills/example/SKILL.md",
    "allowed_missing": ["artifact_class", "kata_surface"],
    "reason": "Legacy skill; v0.1 migration pending",
    "spec_ref": "docs/alpha/ctb/LANGUAGE-SPEC.md §2.1",
    "owner": "cnos.core",
    "remove_when": "frontmatter migration completes"
  }
]
```

Rules:

- JSON exceptions use fields, not comments.
- Exceptions are field-specific, not whole-file blanket passes.
- Exception examples must not contradict hard-gate rules.
- Every exception needs a reason.
- Prefer `remove_when` or follow-up issue for cleanup.

- ❌ `{ "path": "src/packages/cnos.core/skills/example/SKILL.md", "ignore": true }`
- ✅ `{ "path": "src/packages/cnos.core/skills/example/SKILL.md", "allowed_missing": ["artifact_class"], "reason": "Legacy skill; artifact class not migrated yet", "spec_ref": "LANGUAGE-SPEC.md §2.1" }`

### 2.5 Path and resolution semantics

If an issue asks implementers to validate paths, load files, resolve calls, update links, or traverse directories, it MUST define the resolution base.

Questions:

- Is the path repo-root-relative? Package-root-relative? Caller-file-relative?
- Current-working-directory-relative? Branch-relative? Generated-artifact-relative?

- ❌ `Validate each target file exists relative to the skill's parent directory.`
- ✅ `Validate static calls targets relative to the package skill root. Example: src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md with calls: design/SKILL.md resolves to: src/packages/cnos.cdd/skills/cdd/design/SKILL.md not: src/packages/cnos.cdd/skills/cdd/alpha/design/SKILL.md`

Rules:

- Include at least one concrete resolution example.
- If current files already violate the expected shape, name that as migration debt.
- If link checking exists, mention whether the issue should satisfy it.

### 2.6 Acceptance criteria

Acceptance criteria say what done looks like. They do not merely restate work steps.

Every AC MUST be:

- numbered;
- independently testable;
- scoped;
- tied to a source or invariant;
- clear about the proof or oracle.

Preferred AC shape:

```markdown
### AC1: Frontmatter hard gate is enforced

Invariant: SKILL.md frontmatter must include the hard-gate fields.
Oracle: The checker exits non-zero and names the missing field for a fixture missing `scope`.
Positive: A fixture with all hard-gate fields passes.
Negative: A fixture missing `scope` fails, even if listed in exceptions.
Surface: scripts/check-skill-frontmatter.sh, CI job I5.
```

Rules:

- Use this expanded shape for complex/tooling/runtime issues.
- For simple issues, a shorter AC is acceptable if the oracle is obvious.
- At least one AC should cover the negative space for validation/checker/tooling issues.
- ACs must not duplicate the full implementation plan.
- ACs must not include out-of-scope nouns.

- ❌ `AC1: Schema works.`
- ✅ `AC1: Missing scope causes scripts/check-skill-frontmatter.sh to exit non-zero with a diagnostic naming the file and missing field.`

### 2.7 Proof plan

For issues that add validators, checkers, parsers, CI jobs, runtime behavior, or schemas, include a proof plan.

Required fields:

- **invariant:**
- **surface:**
- **oracle:**
- **positive case:**
- **negative case:**
- **operator-visible projection:**
- **known gap:**

Example:

```markdown
## Proof plan

Invariant: Every checked SKILL.md has required frontmatter fields or a field-specific exception.
Surface: schemas/skill.cue, CI job I5.
Oracle: Script exit code and diagnostics.
Positive: Valid fixture passes.
Negative: Missing hard-gate field fails. Invalid enum fails. Broken static calls target fails.
Operator-visible projection: CI job I5 status and notify aggregation include the new check.
Known gap: Does not validate CTB v0.2 witness/close-out semantics.
```

Rules:

- Existing CI passing is the floor, not the proof.
- Every new surface requires at least one CI-gated or documented assertion.
- If proof is partial, say what remains unproven.

### 2.8 Skills and constraints

Implementation issues MUST name Tier 3 skills. Tier 1 CDD lifecycle skills and general Tier 2 defaults do not need to be repeated unless the issue intentionally overrides or clarifies the load set.

Format:

```markdown
## Skills to load

Tier 3:
- cnos.core/skills/design
- eng/tool
- eng/test
- eng/ux-cli

Why:
- design: schema/script/CI boundary and source-of-truth decisions.
- tool: shell script and automation standards.
- test: positive/negative proof and invariants.
- ux-cli: diagnostics and contributor-facing failure output.
```

Rules:

- If a skill changes the issue shape, explain why.
- Do not list skills as decoration.
- If an issue spans docs + tooling + CI, include all three surfaces.

### 2.9 Active design constraints

Name constraints in plain language, with links.

Questions:

- What is the stable policy?
- What is volatile?
- Who owns the source of truth?
- Which boundary must not smear?
- What must remain inspectable?
- What is current state vs target state?

- ❌ `Follow project conventions.`
- ✅ `Design constraint: schema owns field shape/type/enum validation; shell script owns discovery, extraction, exceptions, filesystem checks, and diagnostics. Do not duplicate schema rules in shell.`

### 2.10 Related artifacts

The issue is the entry point. Link deeper artifacts.

Include:

- design doc, if any;
- plan, if any;
- related issues;
- related PRs;
- canonical specs;
- upstream docs;
- affected runtime/docs/CI surfaces;
- source map entries that matter.

Rules:

- Link exact paths, not prose references.
- If a companion artifact is being drafted in parallel, say how conflicts are resolved.
- If the issue updates a README/source map, include link-checking or CI status if relevant.

- ❌ `Related: CTB docs.`
- ✅ `Related: docs/alpha/ctb/LANGUAGE-SPEC-v0.2-draft.md, docs/alpha/ctb/SEMANTICS-NOTES.md, https://github.com/usurobor/tsc/blob/main/spec/c-equiv.md`

### 2.11 Cross-surface updates

When an issue changes a surface that appears in more than one place, name all projections.

Examples:

- **CI job added:** update workflow job; update notification aggregation; document local command; add expected diagnostics; ensure link/check job sees new docs if relevant.
- **README source map changed:** update root README; update docs/README; run link checker; preserve status labels.
- **Schema added:** add schema; add fixtures; add exceptions; add docs; add CI invocation.

Rule: If a new status, check, command, or runtime surface appears anywhere, the issue must ask where else that status is projected.

- ❌ `Add CI job I5.`
- ✅ `Add CI job I5 and update notify job status aggregation so failure summaries include I5.`

### 2.12 Implementation guidance

Implementation guidance is allowed when it prevents ambiguity, but it must not become a hidden plan.

Use implementation guidance for:

- exact branch;
- identity;
- command name;
- file paths;
- resolution base;
- required CI action;
- required local commands;
- known migration debt;
- source-of-truth split.

Do not use it to replace a plan when the work requires multi-step execution.

Format:

```markdown
## Implementation guidance

Branch: sigma/cue-skill-frontmatter
Git identity: sigma <sigma@cdd.cnos>
Local command: scripts/check-skill-frontmatter.sh
CI: Add job I5 in .github/workflows/coherence.yml.
CUE: Use cue-lang/setup-cue with pinned action and CUE version.
```

### 2.13 Handoff checklist

Before filing or dispatching, check:

- [ ] Problem states exists / expected / divergence.
- [ ] Impact says who cares and what is blocked.
- [ ] Status truth is explicit.
- [ ] Source-of-truth paths resolve.
- [ ] Scope and non-goals do not contradict ACs.
- [ ] Hard gates do not appear in exception examples.
- [ ] Optional/defaulted fields name defaults.
- [ ] Path resolution base is explicit.
- [ ] ACs are numbered and independently testable.
- [ ] Proof plan has oracle, positive case, negative case.
- [ ] New surfaces include operator-visible projections.
- [ ] CI additions include notification/status implications.
- [ ] Related artifacts include upstream/canonical docs.
- [ ] Examples obey the rules they are demonstrating.
- [ ] Known gaps are named honestly.
- [ ] No draft/target behavior is described as shipped.

---

## 3. Rules

3.1. **Engineer can act without asking**
  - If the author disappeared, an engineer should still be able to ship or reject the work.

3.2. **Problem is concise**
  - Problem is 3–5 lines. Deeper analysis belongs in design docs, plans, or linked comments.

3.3. **Status is truthful**
  - Every issue must distinguish current implementation, current spec, draft target, planned work, and non-goals when those differ.

3.4. **Source of truth is explicit**
  - Every load-bearing claim points to the canonical source. No "the docs say" without a path.

3.5. **ACs are testable**
  - Each AC has a clear pass/fail. Complex ACs name invariant, oracle, positive case, and negative case.

3.6. **Negative space is mandatory for gates**
  - Validation, checker, schema, parser, runtime, and CI issues must include at least one negative proof.

3.7. **Non-goals are mandatory for substantial issues**
  - Any issue with 3+ ACs needs explicit non-goals.

3.8. **Exceptions are debt**
  - Exceptions must be specific, reasoned, and removable. Blanket exceptions are incoherent.

3.9. **Examples must obey the issue's own rules**
  - If the issue defines hard gates, examples may not show hard gates as exception-backed.

3.10. **Path bases must be named**
  - Every path validation rule must name its resolution base and give one concrete example.

3.11. **Cross-surface projections must be checked**
  - If the issue changes CI, docs, runtime status, commands, schemas, or source maps, name the projected surfaces that must also update.

3.12. **No future-as-present**
  - Do not write target architecture, draft specs, or planned enforcement as if already implemented.

3.13. **Issue owns entry point, not full depth**
  - Issue owns problem, impact, scope, ACs, priority, proof boundary, and links. Design owns gap analysis. Plan owns step sequencing.

---

## 4. Output Pattern

Use this pattern for substantial implementation issues:

```markdown
# Title

Priority:
Status:

## Problem

What exists:
What is expected:
Where they diverge:

## Impact

Who/what is affected:
What is blocked:
Why now:

## Status truth

Shipped:
Current spec:
Draft target:
Not in scope:
Unknown / unverified:

## Source of truth

| Claim / surface | Canonical source | Status | Notes |
|---|---|---|---|

## Scope

In scope:
- ...

Out of scope:
- ...

Deferred:
- ...

Blocked by:
- ...

## Constraint strata

Hard gate:
- ...

Exception-backed:
- ...

Optional/defaulted:
- ...

Validated if present:
- ...

Ignored/deferred:
- ...

## Acceptance criteria

### AC1: ...

Invariant:
Oracle:
Positive:
Negative:
Surface:

### AC2: ...

## Proof plan

Invariant:
Surface:
Oracle:
Positive case:
Negative case:
Operator-visible projection:
Known gap:

## Skills to load

Tier 3:
- ...

Why:
- ...

## Active design constraints

- ...

## Implementation guidance

Branch:
Identity:
Files:
Commands:
CI:
Notes:

## Related artifacts

- ...

## Non-goals

- ...

## Success / closure condition

This issue is closeable when:
- all ACs are met;
- the proof plan passes;
- non-goals remain unviolated;
- known gaps are either resolved or explicitly carried as named debt.

## Handoff checklist

[ ] ...
```

For small issues, compress the shape but preserve:

- problem
- impact
- ACs
- non-goals if substantial
- source-of-truth
- proof or validation surface

---

## 5. Katas

### 5.1 Kata A — schema validation gate

**Given:** A spec defines fields. Current files may violate it. CI does not enforce it. Some legacy files need temporary exceptions.

**Produce an issue that includes:**
current spec path; hard gate fields; exception-backed fields; optional/defaulted fields; enum list; exception shape; path resolution base; positive and negative fixtures; CI job; notification/status aggregation; non-goals.

**Anti-patterns:**
exception example exempts a hard-gate field; CI job added but notify aggregation omitted; schema and shell both own the same validation rule; path resolution base omitted; "verify CI passes" without negative proof.

**Check:** Can an engineer implement without asking which fields are hard failures? Can reviewer reject a missing negative test? Can CI failure tell contributor what to fix?

### 5.2 Kata B — README/source-map alignment

**Given:** The project framing changed. Root README is stale. Docs README has broken links. Upstream formal foundation must be linked. Runtime enforcement does not yet exist.

**Produce an issue that includes:**
current framing; desired framing; affected README files; source map layers; exact doc paths; status labels; link-checking expectation; runtime truth caveat; non-goals.

**Anti-patterns:**
making draft CTB the top-level project frame; linking to stale paths; omitting upstream formal docs; claiming runtime enforcement exists; updating root README but not docs README.

**Check:** Do readers see the current system frame? Do links resolve? Do status labels prevent overclaiming?

### 5.3 Kata C — checker against witness theater

**Given:** A spec defines witness fields. No checker enforces them. A well-written close-out can still be fabricated.

**Produce an issue that includes:**
the risk name; required witness fields; what v0 checks structurally; what v0 does not prove semantically; valid and invalid fixtures; future independence checks; non-goals.

**Anti-patterns:**
field presence treated as semantic truth; no invalid fixture; checker claims runtime enforcement; TSC or upstream verifier treated as already integrated.

**Check:** Does at least one class of false closure become mechanically rejectable? Does the issue avoid claiming more than v0 can prove?

---

## 6. Review Pass

Before filing, run this review:

### 6.1 Contradiction check

- Do any ACs include out-of-scope nouns?
- Do examples violate rules stated elsewhere?
- Do exceptions exempt hard gates?
- Does proof plan promise behavior non-goals exclude?

### 6.2 Source check

- Do paths resolve?
- Are canonical docs preferred?
- Are upstream dependencies linked?
- Are branch/draft/current statuses marked?

### 6.3 Proof check

- Is there an oracle?
- Is there a positive case?
- Is there a negative case?
- Is operator-visible output included when users/operators rely on it?

### 6.4 Design check

- Is there one source of truth per fact?
- Do boundaries have one reason to change?
- Does the issue avoid smearing runtime surfaces?
- Does it avoid future-as-present?

### 6.5 Handoff check

- Can α implement from this issue alone?
- Can β review from the ACs and proof plan?
- Can γ close without hidden conversation context?

---

## 7. Summary

An issue is not just a request. It is the first executable boundary of the work.

A good issue:

- names the gap;
- names the source of truth;
- names status honestly;
- scopes the work;
- defines done;
- defines rejection;
- loads the right skills;
- names non-goals;
- preserves operator-visible consequences;
- prevents contradiction before dispatch.

The issue is coherent when the next agent can act, test, review, and close without recovering hidden context.

---

The main change is philosophical but practical:

```text
Old issue skill:  make the issue clear enough to execute.
New issue skill:  make the issue clear enough to execute, reject, and close honestly.
```

That is what practice exposed. We were not mainly missing better prose. We were missing mandatory checks for status truth, source-of-truth alignment, exception discipline, negative proof, and cross-surface consequences.
