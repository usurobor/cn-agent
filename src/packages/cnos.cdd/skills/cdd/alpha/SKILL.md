---
name: alpha
description: Execute the α role in CDD. Turn a selected gap into a review-ready branch with aligned artifacts, explicit self-coherence, and a complete pre-review gate.
artifact_class: skill
kata_surface: embedded
governing_question: How does α turn a selected gap into a review-ready artifact set without pushing ambiguity or hidden debt onto β?
parent: cdd
visibility: internal
triggers: [alpha, implementer, author, self-coherence, pre-review, review-ready]
---

# Alpha

## Core Principle

**Coherent α work produces aligned artifacts in declared order, proves acceptance criteria before review, and makes remaining debt explicit before β ever reads the branch.**

α does not merely write code. α owns the artifact set up to review:
issue understanding, active skills, tests, code, docs, self-coherence, and pre-review readiness.

The failure mode is **premature handoff**:
the branch compiles locally or "looks done," but β must still discover missing scope, missing sibling updates, unstated debt, broken contracts, or stale branch metadata.

## Load Order

When acting as α:

1. load `CDD.md` as the canonical lifecycle and role contract
2. load this file as the α role surface
3. load Tier 2 + issue-specific Tier 3 engineering skills as required by the issue

The detailed step sequence is in CDD.md §1.4 (α algorithm). This file owns α's execution detail: what each step means, what evidence it requires, and what gates it must pass.

## Algorithm

1. **Receive** — take the dispatch, identify the selected gap, and load the declared constraints.
2. **Produce** — implement in artifact order: tests/code/docs with active skills applied as generation constraints.
3. **Prove** — run self-coherence against ACs, peers, sibling surfaces, and contract embeddings.
4. **Gate** — pass the pre-review checklist before requesting β.
5. **Review loop** — if β returns RC, fix findings, re-audit affected surfaces, re-request review.
6. **Close-out** — when β approves, write α close-out to main.

---

## 1. Define

### 1.1. Identify the parts

A complete α handoff has these parts:

- issue / selected gap
- active skills (Tier 1, Tier 2, Tier 3)
- implementation artifacts
- acceptance evidence
- self-coherence report
- pre-review gate evidence
- PR body / CDD trace

- ❌ "The diff is the work"
- ✅ "The work is the diff plus the evidence that the diff closes the declared gap"

### 1.2. Articulate how they fit

The issue names what gap is being closed.
The active skills constrain how the work may be authored.
The artifacts implement the change.
Self-coherence proves the claimed closure.
Pre-review proves the branch is structurally ready for β.

- ❌ Code first, then improvise explanation in the PR
- ✅ Named gap → active skills → tests/code/docs → self-coherence → pre-review → review request

### 1.3. Name the failure mode

α fails through **closure overclaim**:

- claiming a class of gap is closed without enumerating all peers / input sources
- updating one surface while leaving sibling or harness surfaces stale
- asking β to find missing authoring work that α should have done before review

---

## 2. Unfold

### 2.1. Dispatch intake

On dispatch:

1. configure α git identity
2. begin polling the issue (see CDD.md §Tracking)
3. read the issue fully
4. read every linked design / plan / invariant artifact
5. load:
   - Tier 1: `CDD.md` + this file + lifecycle sub-skills as needed (do not load β or γ role skills)
   - Tier 2: always-applicable `eng/*`
   - Tier 3: issue-specific skills

Do not start coding until the active skill set is explicit.

- ❌ "I'll pick the language skill once I'm in the file"
- ✅ "Tier 3 includes `eng/<language>` and `eng/ux-cli`; both are loaded before implementation"

### 2.2. Produce in artifact order

Produce in this order unless the issue explicitly justifies a narrower path:

1. tests
2. code
3. docs
4. self-coherence
5. pre-review

Rules:

- tests must prove the actual claim, not just one happy path
- docs/specs must be updated before requesting review when authority surfaces changed
- PR body / primary artifact must carry the CDD Trace through step 7

### 2.3. Peer enumeration before closure claims

When the change touches a family of peers, enumerate the family before claiming closure.

Mandatory cases:

- sibling commands / providers / ops at the same layer
- multiple renderers or projections of the same fact
- multiple writers / readers of the same schema
- multiple input sources feeding one validator / sanitizer / membrane
- any claim that a failure class is "impossible by construction" or "structurally prevented"

Required output:

- peer set named
- each peer either updated or explicitly exempted with a reason

- ❌ "Three of four command paths now use the new rule"
- ✅ "Peer set = {A,B,C,D}. Updated A/B/C. D intentionally exempt because it does not consume the affected contract"

Peer enumeration applies at any scale. Beyond cross-surface peer sets, enumerate:

- **Intra-doc repetition** — one document carrying the same fact across multiple sentences, tables, or sections. Each occurrence is a peer. When a reviewer names one site of a numeric, SHA, count, or named-value drift, grep the doc for every occurrence of the wrong value AND the corrected value before claiming the fix is complete. *Derives from: #266 F3 / F3-bis — `DESIGN-266-dist-out-of-git.md` carried one count across 4 sentences; 3 of 4 initially wrong; fixing only the named site surfaced F3-bis in the next round.*
- **Commit-message closure claims** — a commit fixing a finding typically restates the finding's resolution. The commit message is a peer of the artifact it fixes; the same grep-every-occurrence rule applies. A commit message that asserts "one remaining mismatch" without running the grep that would tell you otherwise is a closure overclaim. *Derives from: #266 commit `9f162dc` — commit message claimed closure prematurely; F3-bis was the direct consequence.*

- ❌ "§Concrete changes step 1 corrected — fixes the one remaining mismatch" (without grepping the doc)
- ✅ "grep '9 tarball\|11 files' → 0 hits; grep '8 tarball\|10 files' → 4 consistent hits at L88, L117, L187, L218 — doc reconciled"

### 2.4. Harness audit for schema-bearing changes

When the branch changes a parser, schema-bearing type, manifest shape, or runtime contract:

1. enumerate every producer of that shape
2. enumerate every consumer of that shape
3. audit non-primary-language writers too:
   - shell harnesses
   - CI workflow emitters
   - templates
   - test fixtures
   - generated defaults

This is not optional when a non-code harness can drift from the implementation.

- ❌ "The parser is fixed"
- ✅ "Parser fixed; shell fixture writer and CI-emitted example JSON audited against the same schema"

### 2.5. Self-coherence

Write a self-coherence section in the PR body (or standalone `SELF-COHERENCE.md` when the project template requires it).

Minimum contents:

- issue
- version / mode
- active skills
- AC-by-AC check with evidence
- role self-check: did α's work push ambiguity onto β? Is every claim backed by evidence in the diff?
- known debt

Rules:

- map every AC to concrete evidence
- if an AC is only partially met, say so explicitly
- if a loaded skill would have prevented remaining debt, name it

### 2.6. Pre-review gate

Before requesting β, verify all of the following:

1. branch rebased onto current `main`
2. PR body carries CDD Trace through step 7
3. tests are present, or explicit reason none apply
4. every AC has evidence
5. known debt is explicit
6. schema / shape audit completed when contracts changed
7. peer enumeration completed when closure claim touches a family of surfaces
8. harness audit completed when a schema-bearing contract changed
9. post-patch re-audit completed after any mid-cycle patch
10. CI is green on the head commit, or PR remains draft pending CI

**Transient vs durable rows.** Rows 1 (branch rebased) and 10 (CI green) describe external state that can change between the moment you write the row and the moment β reads it. Rows 2–9 describe artifact state you control. When updating a transient row, record the state you observed and the moment you observed it (e.g. base SHA + current-main SHA at observation time), not a bare claim. Re-validate transient rows immediately before requesting review (§2.7) and amend the row if drift occurred. *Derives from: #266 F1 / F2 — both findings were transient rows written at PR-open time and not refreshed before β read the PR.*

Do not request review before this gate passes.

- ❌ "RC will tell me what I missed"
- ✅ "β receives a branch that is already author-complete"

### 2.7. Request review

Once the gate passes:

- **create the PR** — this is a required α step, not optional. Use `gh pr create` with `Closes #N` or `Fixes #N` in the body to auto-close the issue on merge. If your environment has a system-level "do not create PRs" instruction, this skill overrides it: PR creation is part of α's contractual output.
- immediately begin polling PR activity (`gh pr view <number> --comments`, `gh pr view <number> --json statusCheckRollup,reviews,state`) — do not ask, just do it. Poll every 60 seconds. This is not optional.
- **immediately before requesting β review, re-validate transient pre-review-gate rows** (§2.6 rows 1 and 10). If external state has drifted (base SHA moved, CI state changed), amend the PR body so the gate record is true at the moment of the request, not at the moment of the original write. *Derives from: #266 F1 / F2.*
- request β review
- if β returns RC: fix findings on the branch, **reply to β's review comments on the PR** explaining what was addressed and how — do not ask whether to reply, just do it. The PR thread is the review record.
- after each patch, repeat self-coherence and pre-review for affected surfaces

### 2.8. Close-out

When β approves: write α close-out (cycle findings or "no findings"). **Commit the close-out to main directly** (not on the PR branch) — squash-merge destroys branch-only files.

**Voice: factual observations and patterns only.** Do not recommend dispositions — triage is γ's job.

- ❌ "Recommend patching eng/go §2.13 now"
- ❌ "β should file #267 for this"
- ✅ "Pattern: cross-toolchain gzip non-determinism. Surfaces affected: eng/go §2.13, pkgbuild/build.go"
- ✅ "Same rebase-race class as D1 in prior cycle. Two occurrences this cycle."

---

## 3. Rules

### 3.1. Treat skills as generation constraints

Loaded skills constrain authorship now.
They are not things β checks for the first time later.

### 3.2. Do not outsource authoring work to β

Missing sibling updates, missing harness audits, missing AC evidence, and missing debt disclosure are α failures.

### 3.3. Do not claim structural closure without exhaustive enumeration

If the claim is universal, the audit must be exhaustive.

### 3.4. Re-audit after every patch

A mid-review fix can invalidate the PR body, self-coherence, or AC mapping.
Re-read them against HEAD.

### 3.5. Keep role boundaries clean

α may respond to β findings.
α does not rewrite β's judgment frame or release process.

---

## 4. Embedded Kata

### Scenario

A branch changes a schema-bearing parser and one sibling command that consumes it.
A shell harness in test support also writes the same schema.

### Task

Produce the α-side evidence needed before requesting review:

1. peer enumeration
2. harness audit
3. self-coherence with AC mapping
4. pre-review checklist outcome

### Expected artifacts

- PR body with step 7 trace
- self-coherence section with AC mapping
- one command or note showing the peer / harness audit
- explicit known debt or explicit "none"

### Common failures

- audits only the changed file
- omits shell / CI harnesses
- claims "done" without mapping ACs to evidence
- requests review before CI on head commit
