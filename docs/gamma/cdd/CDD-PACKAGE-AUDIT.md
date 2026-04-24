# CDD Triad Skill-Program Audit — Converged Report

**Date:** 2026-04-24
**Scope:** `src/packages/cnos.cdd/` structural coherence review

## Scope

```
src/packages/cnos.cdd/
  cn.package.json
  commands/cdd-verify/cn-cdd-verify
  skills/cdd/CDD.md
  skills/cdd/SKILL.md
  skills/cdd/alpha/SKILL.md
  skills/cdd/beta/SKILL.md
  skills/cdd/gamma/SKILL.md
  skills/cdd/issue/SKILL.md
  skills/cdd/design/SKILL.md
  skills/cdd/plan/SKILL.md
  skills/cdd/review/SKILL.md
  skills/cdd/review/checklist.md
  skills/cdd/release/SKILL.md
  skills/cdd/post-release/SKILL.md
```

Related cross-package/root surfaces read because the package references them:

```
docs/gamma/cdd/
docs/gamma/cdd/{version}/POST-RELEASE-ASSESSMENT.md
docs/gamma/cdd/{PLAN,SELF-COHERENCE,GATE}-TEMPLATE.md
docs/gamma/cdd/RATIONALE.md
src/packages/cnos.eng/skills/eng/
src/packages/cnos.core/skills/{design,skill}/
```

## Severity key

- **D** = causes wrong agent action or wrong release/closure state
- **C** = structural contract, load, or discovery break
- **B** = ambiguity likely to drift execution
- **A** = polish / navigability

Citation convention: Section refs are authoritative where GitHub raw rendering compresses markdown lines. L refs preserve the repository/local audit line anchors where available.

---

## 1. Directory layout as program structure

### Current state

The directory tree is mostly a valid skill program:

```
cdd/SKILL.md          package-visible loader / public entrypoint
cdd/CDD.md            canonical lifecycle algorithm
cdd/alpha/            α role module
cdd/beta/             β role module
cdd/gamma/            γ role module
cdd/issue/            issue-writing lifecycle module
cdd/design/           design lifecycle module
cdd/plan/             implementation-planning lifecycle module
cdd/review/           review lifecycle module
cdd/release/          release lifecycle module
cdd/post-release/     observe / assess / close lifecycle module
```

Execution model intended by CDD.md: γ dispatches → α produces → β reviews → β releases → β assesses → γ/β close-out surfaces feed next observation cycle.

### Findings

#### 1.1 — Lifecycle/role decomposition is largely sound — A

**Evidence:**
- `skills/cdd/` contains alpha, beta, gamma, issue, design, plan, post-release, release, review, CDD.md, SKILL.md.
- CDD.md §1.4 and §4.1 define the triadic role flow and lifecycle steps.

**Finding:** The directory structure mostly matches the program's control structure.

**Patch:** No structural edit required.

#### 1.2 — `plan/` exists but is under-wired — C

**Evidence:**
- `skills/cdd/plan/SKILL.md` exists and declares `name: plan`.
- CDD.md §4.4 Tier 1 parenthetical omits `plan`.
- CDD.md §5.3 row 6b names artifact "Plan" but assigns Skill = `design`.

**Why this matters:** Agents can load and cite `design` for implementation planning even though `plan/SKILL.md` is the module that owns implementation sequencing.

**Patch:**
- In CDD.md §4.4: change `issue, design, review, release, post-release` to `issue, design, plan, review, release, post-release`
- In CDD.md §5.3 row 6b: change Skill: `design` to: `plan`

#### 1.3 — `review/checklist.md` is orphan/stale — C

**Evidence:**
- `review/` contains `SKILL.md` and `checklist.md`.
- `beta/SKILL.md` load order loads `review/SKILL.md`, `release/SKILL.md`, `post-release/SKILL.md`, but not `review/checklist.md`.
- `review/SKILL.md` does not load `checklist.md`.
- `review/checklist.md` P1.2 says "Only author deletes branch."
- `release/SKILL.md` §2.6a says β deletes merged branches after release push.
- `gamma/SKILL.md` §2.10 says γ closure requires merged remote branches cleaned up.

**Why this matters:** A sidecar checklist that is not loaded cannot constrain execution. If loaded later, it would contradict the release/closure cleanup model.

**Patch option A:** Delete `review/checklist.md`.

**Patch option B:** Add to `review/SKILL.md` §4: "Also load `review/checklist.md` as the condensed pre-submit checklist." Then change `checklist.md` P1.2 to: "Branch cleanup follows release/closure policy: β deletes merged remote branches after release push; γ verifies cleanup at closure and may clean anything still present."

**Recommended:** Option A unless operators still use `checklist.md` manually. The full checklist in `review/SKILL.md` already supersedes it.

#### 1.4 — `post-release/` name versus observe / assess / close lifecycle names — B

**Evidence:**
- CDD.md §4.1 names steps 11–13 as Observe, Assess, Close.
- `post-release/SKILL.md` title and frontmatter say "post-release" and "assessment".
- `post-release/SKILL.md` already says it executes CDD §9 Assessment and §10 Closure.

**Finding:** This is not a missing module. It is a naming/glossary issue: `post-release/` is the umbrella module for steps 11–13.

**Patch:** Add a lead line to `post-release/SKILL.md` after the title: "This module implements CDD steps 11–13: post-release observation, β-owned assessment, and close-out/closure evidence."

No split required unless the operator wants strict one-directory-per-phase mapping. If split later:
```
post-release/           umbrella loader
post-release/observe/
post-release/assessment/
post-release/closeout/
```

---

## 2. Naming coherence

### Current state

Role naming is coherent:
```
alpha ↔ α
beta  ↔ β
gamma ↔ γ
```

Lifecycle naming is mostly coherent, but later-cycle terms overlap:
```
post-release    assessment    close-out    closeout    close    closure
.cdd            protocol surface
docs/gamma/cdd  version surface
hub memory      Engineering skill naming also drifts
```

CDD.md names some `eng/*` skills that do not exist or are deprecated. `review/SKILL.md` points at `core/design`, while CDD.md frames Tier 2 as `eng/*`.

### Findings

#### 2.1 — Glossary gap: post-release / assessment / close-out — B

**Evidence:**
- CDD.md §9: "Assessment" is the post-release judgment artifact.
- CDD.md §10: "Closure" executes immediate/deferred outputs and hub memory.
- `post-release/SKILL.md` mixes "Post-Release Assessment", "CDD Closeout", "Hub Memory", and "Next Move".
- `gamma/SKILL.md` §2.7–§2.10 uses close-out triage and closure gate language.

**Why this matters:** Agents need one answer for whether "close-out" is the assessment, the role-local findings record, the γ triage record, or the final closure declaration.

**Patch:** Add to CDD.md §5, before or after Artifact Contract:

```markdown
### Terminology: post-release, assessment, close-out

- **Post-release** is the umbrella phase after release. It covers lifecycle steps 11–13: observe, assess, close.
- **Assessment** is the β-owned repo artifact: `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` unless the Artifact Location Matrix declares a more specific canonical path.
- **Close-out** is a role-local findings record written by α, β, and γ. Close-outs feed γ triage and CDD iteration. They are not a substitute for the assessment.
- **Closure** is the final state after immediate outputs, deferred outputs, hub memory, and branch cleanup are complete.
```

#### 2.2 — Release tag naming contradicts itself — D

**Evidence:**
- `release/SKILL.md` §2.6 requires bare tags like `3.15.1`, not `v3.15.1`.
- `release/SKILL.md` §3.6 repair recipe uses `vX.Y.Z`.
- `cn-cdd-verify` accepts either `v${VERSION}` or `${VERSION}`.

**Why this matters:** Release CI, version directories, changelog rows, and verifier behavior can diverge. This can make a release appear valid while the canonical tag shape is violated.

**Patch:**

In `release/SKILL.md` §3.6, replace:
```
git commit --amend --no-edit && git tag -d vX.Y.Z && git tag -a vX.Y.Z && git push --force --tags
```
with:
```
git commit --amend --no-edit && git tag -d X.Y.Z && git tag X.Y.Z && git push --force origin main X.Y.Z
```

In `cn-cdd-verify`, replace:
```
git tag -l "v${VERSION}" "$VERSION"
```
with:
```
git tag -l "$VERSION"
```
If a v-prefixed tag exists: warn "v-prefixed tag is legacy/noncanonical"; do not pass the canonical tag check solely because it exists.

#### 2.3 — Architecture skill path drift: `core/design` vs `eng/design` vs `cdd/design` — C

**Evidence:**
- `review/SKILL.md` §2.2.14 says: "Load `core/design` when this check is active."
- CDD.md §4.4 describes Tier 2 as "all general skills under `eng/`" and includes `design`.
- `skills/cdd/design/SKILL.md` is a CDD lifecycle design artifact skill.
- `src/packages/cnos.core/skills/design/SKILL.md` exists.
- `src/packages/cnos.eng/skills/eng/README.md` references `design/`, but the actual `eng` directory does not expose an `eng/design` directory.

**Why this matters:** There are three different possible meanings of "design":
1. CDD lifecycle design artifact (`cdd/design`)
2. Core architectural design skill (`cnos.core/skills/design`)
3. Engineering design skill (`eng/design`) referenced by text but absent

**Patch:** Make a design-skill authority decision. Recommended:
- `review/SKILL.md` §2.2.14 should say: "Load `src/packages/cnos.core/skills/design/SKILL.md` when this check is active."
- CDD.md §4.4 should not list `design` under `eng/` unless an actual `eng/design` skill exists.
- `eng/README.md` should either:
  a) remove `design/` from its skill index and bundles, replacing it with `cnos.core/skills/design`, or
  b) add a real `eng/design/SKILL.md` that delegates to `core/design`.

**No edit — design question for operator:** Should "design" be a core skill, an eng skill, or a CDD lifecycle skill only? Current files imply all three.

#### 2.4 — Trigger overlap: `design` ↔ `plan` — B

**Evidence:**
- `design/SKILL.md` frontmatter triggers include `plan`.
- `plan/SKILL.md` frontmatter triggers include `plan`.

**Why this matters:** If internal triggers are indexed equally, a bare "plan" dispatch can resolve to the design skill rather than the plan skill.

**Patch:** In `design/SKILL.md` frontmatter: remove `plan`. Optional replacement: add `architecture plan` or `design plan` only if phrase triggers are supported.

#### 2.5 — Skill-name drift in CDD.md §5.3 and §4.4 — C

**Evidence:**
- CDD.md §5.3 row 5 uses `eng/README`.
- CDD.md §5.3 row 6c uses `testing`; actual skill path is `eng/test`.
- CDD.md §5.3 row 6d uses `active generation skills`; this is conceptually right but not path-like.
- CDD.md §5.3 row 6e uses `writing`; actual durable-doc skill is `eng/document`. Functional-writing/dataflow skill is `eng/write-functional`.
- CDD.md §4.4 includes `write`; actual eng path is `write-functional`, while docs writing is `document`.
- CDD.md §4.4 includes `skill`; actual `eng/skill` is deprecated and points to `src/packages/cnos.core/skills/skill/SKILL.md`.

**Patch:**

In CDD.md §5.3:
- row 5 Skill: `cdd + applicable eng bundle from src/packages/cnos.eng/skills/eng/README.md`
- row 6c Skill: `eng/test`
- row 6d Skill: `active loaded Tier 2/Tier 3 generation skills, e.g. eng/code, eng/ocaml, eng/go, eng/typescript, eng/tool, eng/ux-cli`
- row 6e Skill: `eng/document` for docs; `eng/write-functional` only when the artifact is functional-programming/design prose or code-shape guidance
- row 10 Skill: `release + eng/document` when release notes/changelog prose is authored

In CDD.md §4.4:
- replace `write` with `document` and/or `write-functional` according to intent
- replace `skill` with `src/packages/cnos.core/skills/skill/SKILL.md` when authoring or modifying skills
- resolve `design` per finding 2.3

---

## 3. File-level contracts

### Current state

Most sub-skills declare:
```yaml
name
description
artifact_class
kata_surface
governing_question
parent: cdd
visibility: internal
triggers
```

The root `cdd/SKILL.md` has frontmatter and acts as the package-visible loader, but it lacks explicit `visibility: public`. `cn.package.json` declares the package and `cdd-verify` command, but no skills stanza.

### Findings

#### 3.1 — Root `cdd/SKILL.md` lacks explicit public visibility — B

**Evidence:**
- `cdd/SKILL.md` has no `visibility:` field.
- CDD.md §11.1 calls `cdd/SKILL.md` the package-visible loader entrypoint.
- Sub-skills are explicitly `visibility: internal`.

**Why this matters:** The visibility contract is implicit for the root but explicit for children.

**Patch:** Add to `src/packages/cnos.cdd/skills/cdd/SKILL.md` frontmatter: `visibility: public`. No parent needed for the root.

#### 3.2 — Trigger overlap between root and sub-skills — B

**Evidence:**
- Root `cdd/SKILL.md` triggers include `review`, `PR`, `release`, `issue`, `design`, `plan`, `assess`, `post-release`, `ship`, `tag`.
- Internal sub-skills repeat many of these triggers.
- Sub-skills are internal, but only if the runtime honors visibility strictly.

**Patch:** Add to `cdd/SKILL.md` after the load-order paragraph: "External dispatch enters through `cdd` only. Internal sub-skill triggers are advisory and are used only after `cdd` and the active role skill have been loaded. The runtime must not expose internal sub-skills as public dispatch entrypoints."

#### 3.3 — Package manifest has no skills stanza — B / design question

**Evidence:**
- `cn.package.json` declares schema/name/version/kind/engines and the `cdd-verify` command.
- It does not declare exposed or internal skills.

**Finding:** This is fine only if skill discovery is directory+frontmatter driven.

**Patch if directory+frontmatter discovery is canonical:** No edit.

**Patch if manifest inventory is still consumed:** Add a skills stanza.

**No edit — design question for operator:** Is manifest-based skill inventory retired? If yes, add a comment/doc note elsewhere so future package authors do not reintroduce stale manifest inventory.

---

## 4. Authority hierarchy

### Current state

CDD.md is the canonical lifecycle and artifact contract. Role skills expand execution detail. Lifecycle sub-skills expand phase procedure. This hierarchy is mostly stated correctly, but several files still duplicate or soften canonical rules.

### Findings

#### 4.1 — CDD/role-skill conflict rule should be sharper — B

**Evidence:**
- CDD.md §1.4 role-algorithm blurbs say "When they diverge on role execution, the skill governs."
- `cdd/SKILL.md` conflict rule is more precise: CDD.md governs ordered steps, selection rules, and artifact contract; role skills govern local execution.

**Why this matters:** A role skill could be read as authorized to change lifecycle order or artifact obligations.

**Patch:** In each CDD.md §1.4 role-algorithm sentence that says: "When they diverge on role execution, the skill governs." replace with: "When they diverge on role-local execution detail, the skill governs, provided it does not alter role ownership, lifecycle order, selection rules, dispatch contract, artifact contract, or closure obligations."

#### 4.2 — α artifact order omits design / contract / plan production — C

**Evidence:**
- CDD.md §5.2 canonical artifact order: design → coherence contract → plan → tests → code → docs → self-coherence.
- `alpha/SKILL.md` §2.2 says α produces: tests → code → docs → self-coherence → pre-review.
- `alpha/SKILL.md` §2.1 says α reads linked design/plan artifacts, but does not require producing them or explicitly marking them "not required".

**Why this matters:** α can appear compliant while skipping the first three canonical artifact-order decisions.

**Patch:** In `alpha/SKILL.md` §2.2, replace the order with:

```
Produce in CDD artifact order unless the issue explicitly justifies a narrower path:
1. design artifact if required, or explicit "not required"
2. coherence contract
3. plan if required, or explicit "not required"
4. tests
5. code
6. docs
7. self-coherence
8. pre-review
```

In `alpha/SKILL.md` Load Order, add:
- load `issue/SKILL.md` when issue quality or AC boundaries are being interpreted
- load `design/SKILL.md` when producing or judging design-required work
- load `plan/SKILL.md` when implementation sequencing is non-trivial

#### 4.3 — Review severity table contradicts review §7.0 and CDD §5.5 — C

**Evidence:**
- `review/SKILL.md` §5 says C = "APPROVED with note".
- `review/SKILL.md` §7.0 says A/B/C/D findings must be fixed before merge.
- CDD.md §5.5 says all review findings must be resolved before merge and forbids "approved with follow-up".

**Why this matters:** Reviewers can approve a C finding as "non-blocking" while merge-readiness says all findings must be fixed on-branch.

**Patch:** Replace `review/SKILL.md` §5 severity table Action column with two columns:

| Severity | Meaning | Review verdict | Merge readiness |
|----------|---------|----------------|-----------------|
| D | demonstrable incoherence | REQUEST CHANGES | not merge-ready |
| C | real incoherence | APPROVED only after on-branch fix + re-check, otherwise REQUEST CHANGES | not merge-ready until fixed |
| B | improvement opportunity | APPROVED with required on-branch fix | not merge-ready until fixed |
| A | polish | APPROVED with required on-branch fix | not merge-ready until fixed |

Or shorter: "C/B/A may be lower-severity, but they are still required on-branch fixes unless explicitly deferred by design scope per §7.0."

#### 4.4 — Branch cleanup ownership duplicated — B

**Evidence:**
- CDD.md γ algorithm step 14 says γ deletes merged remote branches.
- `gamma/SKILL.md` §2.10 says merged remote branches are cleaned up, then γ deletes them.
- `release/SKILL.md` §2.6a says β deletes merged branches immediately after release push.

**Why this matters:** The action is duplicated. Agents may either both attempt cleanup or defer it because the other role "owns" it.

**Patch:** Declare:
- β owns branch cleanup immediately after release push.
- γ owns verification at closure and may clean leftovers if β could not.

In CDD.md γ step 14 and `gamma/SKILL.md` §2.10, change: "Delete merged remote branches" to: "Verify merged remote branches were cleaned by β; delete any still present only as closure cleanup."

#### 4.5 — Step 12a skill/spec patch ownership is ambiguous — C

**Evidence:**
- CDD.md §5.3 row 12a owner/producer says "γ (or releasing agent)".
- CDD.md β algorithm step 9 gives β the post-release assessment.
- `post-release/SKILL.md` Step 5 says recurring failure skill patches must land in the same commit as the assessment.
- `gamma/SKILL.md` §2.8–§2.9 says γ independently checks and patches process gaps when clear.

**Why this matters:** A required skill patch can be skipped because β thinks γ owns it, while γ thinks β must land it with the assessment.

**Patch:** In CDD.md §5.3 row 12a:
- Owner: β primary; γ fallback / independent process-gap owner
- Producer: releasing β when identified during assessment; γ when identified during close-out triage or missed by β

Add rule: "If the patch is identified in the assessment, β lands it with the assessment. If γ independently identifies a process gap after reading close-outs, γ lands a follow-on patch or commits a concrete next MCA."

---

## 5. Load order as dependency graph

### Current graph

```
cdd/SKILL.md
  -> CDD.md
  -> active role skill
  -> lifecycle sub-skills as directed
  -> Tier 2 / Tier 3 skills

alpha/SKILL.md
  -> CDD.md -> alpha/SKILL.md
  -> Tier 2 + Tier 3
  -> lifecycle sub-skills as needed

beta/SKILL.md
  -> CDD.md -> beta/SKILL.md
  -> review/SKILL.md -> release/SKILL.md -> post-release/SKILL.md
  -> Tier 2 + Tier 3 as required by issue/diff

gamma/SKILL.md
  -> CDD.md -> gamma/SKILL.md
  -> issue/SKILL.md
  -> post-release/SKILL.md when close-out
  -> other lifecycle sub-skills only when selected gap requires them

review/SKILL.md
  -> issue ACs -> PR diff -> CDD artifacts -> active skills
  -> project design constraints
  -> core/design when architecture check active

release/SKILL.md
  -> VERSION -> stamp/check scripts -> CHANGELOG.md -> RELEASE.md
  -> tag/release CI -> deployment validation
  -> post-release assessment

post-release/SKILL.md
  -> CHANGELOG TSC table -> release artifacts -> review metrics -> lag table
  -> assessment artifact -> skill patches -> hub memory
```

No hard SKILL-load cycle is required. The back-references from sub-skills to CDD.md are authority references, not execution cycles.

### Findings

#### 5.1 — Tier 1 is internally inconsistent — C

**Evidence:**
- CDD.md §4.4 says Tier 1 = all skills under `cdd/`, always loaded by every role, not optional.
- `alpha/SKILL.md` §2.1 says lifecycle sub-skills are loaded as needed.
- `gamma/SKILL.md` Load Order says other lifecycle sub-skills are loaded only when the selected gap requires them.
- `beta/SKILL.md` always loads review/release/post-release, not all lifecycle sub-skills.

**Why this matters:** The canonical load rule is stricter than the role files and likely impractical. It also makes every role load irrelevant peer role/lifecycle surfaces.

**Patch:** Rewrite CDD.md §4.4:

```
Skills are loaded in three tiers. For substantial changes, all applicable tiers are mandatory.

Tier 1a — CDD authority, always loaded
  - CDD.md
  - cdd/SKILL.md
  - the current role skill: alpha/, beta/, or gamma/

Tier 1b — CDD lifecycle phase skills, phase-required
  Load the lifecycle sub-skills required by the current phase/work shape:
  - issue
  - design
  - plan
  - review
  - release
  - post-release

Tier 1c — β closure bundle, always loaded by β
  β always loads:
  - review/SKILL.md
  - release/SKILL.md
  - post-release/SKILL.md

Tier 2 — General engineering, role/work-shape required
  Load the applicable engineering bundle from src/packages/cnos.eng/skills/eng/README.md.

Tier 3 — Issue-specific
  Named by the issue. γ names only Tier 3 in issues.
```

#### 5.2 — CDD.md references non-existent gamma sections — B

**Evidence:**
- CDD.md §5.3 row 12a cites `gamma/SKILL.md §2.10–§2.12`.
- `gamma/SKILL.md` ends at §2.10.

**Patch:** Change CDD.md row 12a reference to: `gamma/SKILL.md §2.7–§2.10`. Or add explicit gamma headings: §2.11 Hub memory, §2.12 Branch cleanup and closure declaration.

**Recommended:** Change the reference. `gamma/SKILL.md` already maps steps 8–15 to §2.7–§2.10.

#### 5.3 — `RATIONALE.md` reference is misplaced from package perspective — B

**Evidence:**
- CDD.md §11.2 says "`RATIONALE.md` explains why CDD takes this shape."
- No `RATIONALE.md` exists beside `src/packages/cnos.cdd/skills/cdd/CDD.md`.
- `docs/gamma/cdd/RATIONALE.md` exists.
- `docs/gamma/cdd/3.15.2/RATIONALE.md` exists as a frozen snapshot.

**Patch option A:** Change CDD.md §11.2 to: "`docs/gamma/cdd/RATIONALE.md` explains why CDD takes this shape. Frozen release snapshots may also contain version-local rationale files."

**Patch option B:** Add `src/packages/cnos.cdd/skills/cdd/RATIONALE.md` as a pointer file: "Canonical rationale: `docs/gamma/cdd/RATIONALE.md`"

**Recommended:** Option A if the docs bundle owns rationale; Option B if each skill package is expected to be self-contained.

#### 5.4 — `review/checklist.md` has no load edge — C

Same as finding 1.3. **Patch:** Delete or load explicitly.

---

## 6. Dispatch surface

### Current state

CDD.md and `gamma/SKILL.md` dispatch prompts correctly point agents at:
```
src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md
src/packages/cnos.cdd/skills/cdd/beta/SKILL.md
src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md
```

The prompt structure correctly points roles at the issue, not a paraphrase.

### Findings

#### 6.1 — Dispatch paths are correct — A

**Evidence:**
- CDD.md §1.4 prompt format names `alpha/SKILL.md`, `beta/SKILL.md`, `gamma/SKILL.md`.
- `gamma/SKILL.md` §2.5 uses the same paths.
- Those files exist.

**Patch:** No edit.

#### 6.2 — `<number>` parameter ambiguity — B

**Evidence:**
- CDD.md §1.4 says `<number>` is the issue or PR number.
- All prompt templates use `gh issue view <N>`.
- `gamma/SKILL.md` says both roles should be pointed at the issue.

**Why this matters:** β may be dispatched to a PR number while the prompt asks `gh issue view`, or α may miss the canonical issue contract.

**Patch:** Change the CDD.md parameter note to: "`<number>` is the issue number for dispatch prompts. β discovers the PR from the issue, α's review request, or GitHub notifications. A PR number may be used only for post-hoc review/recovery workflows, and the prompt must say explicitly that it is PR-scoped."

#### 6.3 — Tier 1/2/3 classification is not executable as written — C

**Evidence:**
- CDD.md §4.4 Tier 1 conflicts with role load orders.
- CDD.md §4.4 Tier 2 names skills that do not exactly match the actual eng/core inventory.
- `issue/SKILL.md` says issues should name only Tier 3 skills.
- `gamma/SKILL.md` repeats that γ names only Tier 3.

**Patch:** Apply findings 2.3, 2.5, and 5.1:
- rewrite Tier 1 into 1a/1b/1c
- map Tier 2 to actual cnos.eng bundles and core skills
- keep issue/gamma rule: issue names only Tier 3

---

## 7. Close-out and assessment surfaces

### Current state

This is the least coherent surface. Current files mention all of these as assessment/close-out locations:

```
docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md
docs/gamma/cdd/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md
.cdd/<version>/POST-RELEASE-ASSESSMENT.md
.cdd/releases/{X.Y.Z}/beta/POST-RELEASE-ASSESSMENT.md
.cdd/releases/{X.Y.Z}/beta/ASSESSMENT.md
.cdd/releases/{X.Y.Z}/{alpha,beta,gamma}/CLOSE-OUT.md
PR comments
hub memory
```

There is no single canonical answer for "where does artifact X go?"

### Findings

#### 7.1 — Post-release assessment path conflicts across spec, skill, verifier — D

**Evidence:**
- CDD.md §9.1 says the assessment is located at: `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md`
- For this package, actual observed release assessment path: `docs/gamma/cdd/3.57.0/POST-RELEASE-ASSESSMENT.md`
- `post-release/SKILL.md` Writing process says expected path example: `.cdd//POST-RELEASE-ASSESSMENT.md` which is malformed and not aligned with either CDD.md or `cdd-verify`.
- `cn-cdd-verify` accepts:
  - `.cdd/releases/$VERSION/beta/POST-RELEASE-ASSESSMENT.md`
  - `.cdd/releases/$VERSION/beta/ASSESSMENT.md`
  - `docs/gamma/cdd/$VERSION/POST-RELEASE-ASSESSMENT.md`

**Why this matters:** β and γ can each write/check a different artifact and still believe the cycle is complete. The verifier currently validates compatibility, not canonicality.

**Patch:** Add to CDD.md §5 after §5.3:

```markdown
### 5.x Artifact Location Matrix

| Artifact | Canonical repo location | CDD package default | Noncanonical / scratch / legacy |
|---|---|---|---|
| Version snapshot directory | `docs/{tier}/{bundle}/{X.Y.Z}/` | `docs/gamma/cdd/{X.Y.Z}/` | `.cdd/` is not the frozen snapshot |
| POST-RELEASE-ASSESSMENT.md | `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` | `docs/gamma/cdd/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` | `.cdd/releases/.../ASSESSMENT.md` is triadic scratch/legacy unless explicitly declared canonical |
| α close-out | `.cdd/releases/{X.Y.Z}/alpha/CLOSE-OUT.md` | same | PR comment only for PR-scoped unreleased cycles |
| β close-out | `.cdd/releases/{X.Y.Z}/beta/CLOSE-OUT.md` | same | PR comment only for PR-scoped unreleased cycles |
| γ close-out | `.cdd/releases/{X.Y.Z}/gamma/CLOSE-OUT.md` | same | none |
| γ kata verdict | `.cdd/releases/{X.Y.Z}/gamma/KATA-VERDICT.md` when kata is available | same | optional/warn-only if kata unavailable |
| CHANGELOG ledger row | `CHANGELOG.md` | same | none |
| RELEASE.md | repo root `RELEASE.md` in release commit | same | CI auto-generated release body is not acceptable |
| Hub memory | external agent hub | external agent hub | not a repo artifact; assessment records path/commit/sha or unavailable reason |
```

Then update `post-release/SKILL.md` Writing process step 1 to: "Create the assessment file at `docs/gamma/cdd/X.Y.Z/POST-RELEASE-ASSESSMENT.md` for the CDD package, or more generally at `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md`."

Then update `cn-cdd-verify`:
- pass only on canonical assessment path
- warn for legacy `.cdd/releases/.../ASSESSMENT.md`
- fail if canonical assessment is missing for release-scoped cycles

#### 7.2 — α/β/γ close-out paths are underspecified — C

**Evidence:**
- CDD.md α and β algorithms require close-outs written to a file and committed to main, but do not specify paths.
- `gamma/SKILL.md` §2.10 requires α and β close-outs to exist on main.
- `cn-cdd-verify` accepts close-outs in PR comments or `.cdd/`.

**Why this matters:** "Committed to main" is not enough. Agents need exact paths, and the verifier must know whether PR comments are canonical or fallback evidence.

**Patch:** Add to the Artifact Location Matrix:
- α close-out: `.cdd/releases/{X.Y.Z}/alpha/CLOSE-OUT.md`
- β close-out: `.cdd/releases/{X.Y.Z}/beta/CLOSE-OUT.md`
- γ close-out: `.cdd/releases/{X.Y.Z}/gamma/CLOSE-OUT.md`

Add rule: "PR comments may satisfy close-out evidence only for PR-scoped, unreleased, non-triadic cycles. For release-scoped triadic cycles, close-outs are repo artifacts under `.cdd/releases/{X.Y.Z}/`."

Update `alpha/SKILL.md` §2.8 and `beta/SKILL.md` closure discipline to name these paths.

#### 7.3 — Hub memory boundary is mandatory but out-of-repo — B

**Evidence:**
- CDD.md §10.1 requires hub memory writes.
- `post-release/SKILL.md` template includes Hub Memory paths and SHAs.
- Hub memory is not a repo artifact.

**Why this matters:** A repo-only verifier cannot fully prove closure unless the assessment records the external memory evidence.

**Patch:** In the Artifact Location Matrix, add: "Hub memory is external agent-hub state, not a repo artifact. The assessment must record the path, commit/sha if git-backed, or an explicit unavailable reason."

In `cn-cdd-verify`:
- do not attempt to inspect hub memory directly
- check that `POST-RELEASE-ASSESSMENT.md` contains a Hub Memory section
- warn if paths/sha are absent

---

## 8. Completeness

### Current state

No literal TODO markers or unfinished stub sections were found in the inspected CDD skill tree, package manifest, or `cdd-verify` command. The main completeness issues are stale references, malformed paths, deprecated skill names, and template/path references outside the package root.

### Findings

#### 8.1 — External template references live outside the package root — B

**Evidence:**
- CDD.md §5.3 references:
  - `docs/gamma/cdd/PLAN-TEMPLATE.md`
  - `docs/gamma/cdd/SELF-COHERENCE-TEMPLATE.md`
  - `docs/gamma/cdd/GATE-TEMPLATE.md`
- `post-release/SKILL.md` references the same `docs/gamma/cdd` template area.
- These files exist in `docs/gamma/cdd/`, not under `src/packages/cnos.cdd`.

**Finding:** This is coherent if the package intentionally depends on repo-level human documentation templates. It is incoherent if packages are expected to be self-contained.

**Patch option A:** Leave paths unchanged and add to `cn.package.json` or package README: "This package has repo-level documentation template dependencies under `docs/gamma/cdd/`."

**Patch option B:** Mirror templates into `src/packages/cnos.cdd/skills/cdd/templates/` and update CDD.md §5.3 references.

**No edit — design question for operator:** Are skill packages supposed to be self-contained, or may they depend on repo-level `docs/gamma/cdd` templates?

#### 8.2 — `cn-cdd-verify` encodes compatibility conventions rather than canonical CDD — C

**Evidence:**
- accepts both bare and v-prefixed tags
- accepts assessment in multiple paths
- accepts close-outs in PR comments or `.cdd`
- warns on many missing triadic artifacts instead of enforcing canonical release-scoped closure

**Why this matters:** The verifier can pass cycles that are structurally noncanonical.

**Patch after Artifact Location Matrix lands:**

`cn-cdd-verify` should:

Required for `--version`:
- CHANGELOG contains bare version
- bare git tag exists
- canonical `POST-RELEASE-ASSESSMENT.md` exists
- release note / `RELEASE.md` evidence exists if release was just tagged

Required for `--version --triadic`:
- `.cdd/releases/{X.Y.Z}/alpha/CLOSE-OUT.md`
- `.cdd/releases/{X.Y.Z}/beta/CLOSE-OUT.md`
- `.cdd/releases/{X.Y.Z}/gamma/CLOSE-OUT.md`
- `.cdd/releases/{X.Y.Z}/gamma/KATA-VERDICT.md` if kata available, otherwise warning only

Warnings only:
- v-prefixed tag exists
- legacy assessment path exists
- PR comment close-out exists but canonical close-out missing

#### 8.3 — post-release CHANGELOG row omits Level — C

**Evidence:**
- `release/SKILL.md` §2.4 says the ledger row includes: Version, C_Σ, α, β, γ, Level, coherence note.
- CDD.md §9.1 says cycle level is recorded in the CHANGELOG TSC table.
- `post-release/SKILL.md` Step 2 shows: `| vX.Y.Z | C_Σ | α | β | γ | Coherence note |`

**Patch:** In `post-release/SKILL.md` Step 2, replace with: `| X.Y.Z | C_Σ | α | β | γ | Level | Coherence note |`

Also align prefix policy: use `X.Y.Z` if CHANGELOG canonical row uses bare version; use `vX.Y.Z` only if CHANGELOG explicitly requires v-prefix display.

**Recommended:** Use bare `X.Y.Z` everywhere unless CHANGELOG's human-facing display format explicitly keeps `v`.

#### 8.4 — CDD Trace template omits rows 2 / 3 / 7 / 7a — B

**Evidence:**
- CDD.md §5.3 requires evidence for Branch, Bootstrap, Self-coherence, and Pre-review.
- `alpha/SKILL.md` §2.6 requires PR body carries CDD Trace through step 7.
- CDD.md §5.4 trace template includes rows 0, 1, 4, 5, 6, 8, 9, 10.

**Why this matters:** Required lifecycle steps are not present in the example trace agents copy.

**Patch:** In CDD.md §5.4 template, add:
```
| 2 Branch       | branch              | cdd  | branch created / verified                        |
| 3 Bootstrap    | version dir          | cdd  | stubs created or explicit small-change exemption |
| 7 Self-coherence | SELF-COHERENCE.md / PR body | cdd | self-check completed                    |
| 7a Pre-review  | PR body              | cdd  | pre-review gate passed                           |
```

#### 8.5 — `design/SKILL.md` subsection numbering is inconsistent — A

**Evidence:**
- `design/SKILL.md` §2.2 "Constraints" uses local subsection numbers 3.1–3.4.
- §2.3 "Impact Graph" uses 4.x.
- §2.4 "Proposal" uses 5.x.

**Why this matters:** Navigational friction only. It can confuse citation and "load section X" prompts.

**Patch:** Renumber local subsections under each markdown heading: 2.2.1, 2.2.2, ... 2.3.1, 2.3.2, ... 2.4.1, 2.4.2, ... Or remove nested numeric prefixes and rely on markdown headings.

#### 8.6 — eng README references absent/deprecated skills — C, outside package but affects CDD dispatch

**Evidence:**
- `src/packages/cnos.eng/skills/eng/README.md` references `design/`, `review/`, and `write` in bundles/defaults.
- actual `eng` directory lacks `design/` and `review/` directories.
- actual `eng` directory has `write-functional/`, not `write/`.
- `eng/skill/SKILL.md` is deprecated and points to `cnos.core/skills/skill/SKILL.md`.

**Why this matters:** CDD.md §4.4 depends on the eng taxonomy. If eng README is stale, CDD Tier 2 cannot be made executable without fixing or overriding it.

**Patch:** In `eng/README.md`:
- replace `design/` with `src/packages/cnos.core/skills/design/SKILL.md`, or add a real `eng/design` skill
- remove `review/` from eng bundles unless a real `eng/review` skill is added
- replace `write` with `document` or `write-functional` depending on intent
- replace meta-skill path with `src/packages/cnos.core/skills/skill/SKILL.md`

Then update CDD.md §4.4 to cite the corrected eng README.

---

## Coherence map

### Files that agree strongly

**1. Role ownership:**
- CDD.md, cdd/SKILL.md, alpha/SKILL.md, beta/SKILL.md, gamma/SKILL.md
- Agreement: α produces. β reviews, merges, releases, assesses. γ selects, dispatches, unblocks, triages close-outs, closes the loop.

**2. Artifact boundaries:**
- issue/SKILL.md, design/SKILL.md, plan/SKILL.md
- Agreement: Issue owns concise problem, outcome ACs, priority, scope. Design owns gap analysis, constraints, proposal, impact graph. Plan owns sequence, per-step ACs, file changes, test strategy.

**3. Merge discipline:**
- CDD.md §5.5, beta/SKILL.md role rules, review/SKILL.md §7.0
- Agreement: No merge with unresolved findings, except explicit design-scope deferral filed before merge.

**4. β continuity:**
- beta/SKILL.md, release/SKILL.md, post-release/SKILL.md
- Agreement: β/reviewer context should carry review → merge → release → assessment.

### Files that drift

**1. Tier loading:**
- CDD.md §4.4 says all CDD lifecycle skills are always loaded by every role.
- alpha/gamma load lifecycle sub-skills as needed.
- beta loads review/release/post-release.

**2. Plan wiring:**
- `plan/SKILL.md` exists.
- CDD.md §4.4 omits plan from Tier 1 parenthetical.
- CDD.md §5.3 row 6b assigns Plan to design.

**3. Tag naming:**
- `release/SKILL.md` §2.6 requires bare tags.
- `release/SKILL.md` §3.6 uses v-prefixed repair tags.
- `cn-cdd-verify` accepts both.

**4. Assessment paths:**
- CDD.md §9.1 uses `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md`.
- `post-release/SKILL.md` example uses malformed `.cdd//POST-RELEASE-ASSESSMENT.md`.
- `cn-cdd-verify` accepts multiple `.cdd` and `docs` paths.
- actual observed CDD package assessments live in `docs/gamma/cdd/{version}/`.

**5. Close-out paths:**
- CDD.md says write close-outs to files on main.
- `gamma/SKILL.md` requires close-outs on main.
- `cn-cdd-verify` accepts PR comments or `.cdd` paths.
- No canonical path is declared in the skills.

**6. Branch cleanup:**
- `release/SKILL.md` makes β delete merged branches after release push.
- `gamma/SKILL.md` makes γ delete merged branches at closure.
- `review/checklist.md` says only author deletes branch.

**7. Architecture/design skill:**
- `review/SKILL.md` says `core/design`.
- CDD.md Tier 2 says `eng/design`.
- `cdd/design` is lifecycle design.
- `eng` README references design but actual eng directory lacks design.

**8. Rationale path:**
- CDD.md §11.2 says `RATIONALE.md` without path.
- no package-local `RATIONALE.md` exists.
- `docs/gamma/cdd/RATIONALE.md` exists.

**9. Gamma section refs:**
- CDD.md row 12a cites `gamma §2.10–§2.12`.
- `gamma/SKILL.md` ends at §2.10.

**10. Trace template:**
- CDD.md §5.3 requires steps 2, 3, 7, 7a evidence.
- CDD.md §5.4 template omits those rows.

---

## Highest-leverage patch set

### P1 — D/C fixes to land first

1. Add Artifact Location Matrix to CDD.md and sync: `post-release/SKILL.md`, `release/SKILL.md`, `gamma/SKILL.md`, `alpha/SKILL.md`, `beta/SKILL.md`, `cn-cdd-verify`
2. Normalize tags: bare versions everywhere; `release/SKILL.md` §3.6 repair recipe uses bare tag; `cn-cdd-verify` requires bare tag; v-prefixed tag is warning-only legacy
3. Fix plan wiring: add plan to CDD.md §4.4 Tier 1 lifecycle list; CDD.md §5.3 row 6b Skill = plan
4. Rewrite CDD.md §4.4: Tier 1a always (CDD.md + cdd/SKILL.md + current role skill), Tier 1b phase-required (lifecycle sub-skills by current phase/work shape), Tier 1c β closure bundle (review + release + post-release), Tier 2 from corrected eng/core skill inventory, Tier 3 issue-specific only
5. Patch `alpha/SKILL.md` §2.2: include design / coherence contract / plan before tests/code/docs; allow explicit "not required"; add load edges to issue/design/plan where relevant
6. Resolve review severity contradiction: split verdict from merge readiness; no C/B/A unresolved findings at merge unless design-scope deferral filed
7. Resolve step 12a ownership: β primary for patches identified during assessment; γ fallback/independent process-gap owner after close-out triage
8. Delete or wire `review/checklist.md`: recommended delete; if kept, align branch cleanup with release/closure policy
9. Fix architecture/design skill naming: choose `core/design` vs `eng/design`; update `review/SKILL.md`, CDD.md §4.4, and eng README consistently
10. Fix Skill-column names: `testing` → `eng/test`; `writing` → `eng/document` or `eng/write-functional` by intent; `eng/README` → applicable eng bundle; `skill` → `cnos.core/skills/skill` when authoring skills; `active generation skills` → active loaded Tier 2/Tier 3 generation skills
11. Add Level column to `post-release/SKILL.md` CHANGELOG row
12. Update `cn-cdd-verify` to enforce canonical paths/tags and warn on legacy

### P2 — B/A cleanup

13. Add `visibility: public` to root `cdd/SKILL.md`
14. Add glossary to CDD.md: post-release / assessment / close-out / closure
15. Fix CDD.md row 12a reference: `gamma/SKILL.md §2.7–§2.10`
16. Repoint CDD.md §11.2: `docs/gamma/cdd/RATIONALE.md` or add package-local pointer `RATIONALE.md`
17. Add CDD Trace rows: 2 Branch, 3 Bootstrap, 7 Self-coherence, 7a Pre-review
18. Remove `plan` from `design/SKILL.md` triggers
19. Clarify dispatch `<N>`: issue number for dispatch; PR discovered by role
20. Renumber `design/SKILL.md` local subsections
21. Fix `eng/README` stale references: design/review/write/skill path drift

---

## No-edit design questions for operator

1. **Is skill discovery manifest-based or directory/frontmatter-based?**
   - If directory/frontmatter-based, `cn.package.json` stays minimal.
   - If manifest-based, add a skills stanza.

2. **Should `post-release` split into separate modules?**
   - Current `post-release/` can coherently own steps 11–13 with glossary and step mapping.
   - Split only if the operator wants strict one-directory-per-lifecycle-phase.

3. **Should CDD templates live inside the package?**
   - Current package references `docs/gamma/cdd` templates.
   - Mirror templates into package only if package self-containment is required.

4. **What is the canonical "design" skill?**
   - `cdd/design` = lifecycle artifact design
   - `cnos.core/skills/design` = architecture/design reasoning
   - `eng/design` = currently referenced by text but absent
   - Choose one authority relationship and patch all references.

5. **What is `.cdd/` supposed to be?**
   - Recommended: triadic protocol/close-out scratch and role-local records, not the canonical frozen release assessment.
   - If `.cdd/` is intended as canonical, CDD.md §9.1 and `docs/gamma/cdd` assessment practice must change.
