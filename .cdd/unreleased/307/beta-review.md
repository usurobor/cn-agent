# β Review — #307

skill(cdd/issue): move issue katas to cnos.cdd.kata package

---

## Round 1

**Verdict:** REQUEST CHANGES

**Round:** 1
**Base SHA (origin/main):** `396d9982af2c2318400fb5c9eb400178d488f99b`
**Head SHA (cycle/307):** `7889727ac1f558dfcbcef04f7f281ede956717b5`
**Branch CI state:** not configured on cycle branches (per α self-coherence §Review-readiness; CI gates on main push; β proceeds)
**Merge instruction (on approval):** `git merge cycle/307 -m "Closes #307: skill(cdd/issue): move issue katas to cnos.cdd.kata package"` into main

---

## §2.0.0 Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes | shipped / current spec / planned labels used correctly; `kata_surface: embedded` (was) vs `external` (now) is accurate state transition |
| Canonical sources/paths verified | yes | `docs/alpha/ctb/LANGUAGE-SPEC.md` §2.1 confirms `external` is in enum; `src/packages/cnos.cdd.kata/katas/M5-issue-authoring/` confirmed to exist at head |
| Scope/non-goals consistent | yes | diff touches exactly the three in-scope files (issue/SKILL.md + M5 bundle); non-goals unviolated (no content rewriting, review katas untouched, CTB v0.2 not promoted) |
| Constraint strata consistent | yes | hard gate (kata.md present) ✓; exception-backed (rubric.json, prompt-pair) — α shipped them; optional `kata_ref` key is present and resolves |
| Exceptions field-specific/reasoned | yes | no exceptions exercised; `schemas/skill.cue` open-schema trailing `...` is pre-existing, not a new exception |
| Path resolution base explicit | yes | `kata_ref` value is repo-root-relative; resolves to a directory containing `kata.md` |
| Proof shape adequate | yes | AC1–AC5 each carry invariant / oracle / evidence; §Self-check runs the AC1 grep oracles and reports 0 matches |
| Cross-surface projections updated | **no** | `docs/gamma/cdd/KATAS.md` Tier 3 table lists M0–M4 only; M5-issue-authoring is not added. KATAS.md is self-described as "a secondary index" and the mirror precedent (#304 / `5a8bb3e`) did not update it, but the §Debt section says "No known α debt new this cycle" without naming this gap — the omission is silent, not explicit |
| No witness theater / false closure | yes | §Self-check runs AC1 oracles and reports counts; §Peer enumeration lists all nine lifecycle skills; evidence is grep + filesystem, not prose assertion |
| PR body matches branch files | n/a | triadic protocol; no PR body; issue scope vs diff verified directly |

One "no" row: cross-surface projection (KATAS.md). Proceeding to implementation review; verdict accounts for the contract finding.

---

## §2.0 Issue Contract

### AC Coverage

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| AC1 | §5 Katas body removed from issue/SKILL.md | yes | **met** | `grep -nE "^### 5\." ...` → 0 matches; `grep -nE "Kata A\|Kata B\|Kata C" ...` → 0 matches; verified against cycle/307 head |
| AC2 | New kata location(s) exist under cnos.cdd.kata | yes | **met** | `src/packages/cnos.cdd.kata/katas/M5-issue-authoring/` exists with all four required files |
| AC3 | Each new kata directory carries kata.md matching M-series shape | yes | **met** | kata.md opens with `**Class:**`, `**Default level target:**`, `**Purpose:**`, `## Scenario`, `## Required artifacts`, `## Scoring`, `## Worked examples` — mirrors M2-review shape exactly |
| AC4 | issue/SKILL.md frontmatter declares external kata | yes | **met** | L5 `kata_surface: external`; L26 `kata_ref: src/packages/cnos.cdd.kata/katas/M5-issue-authoring/` — path resolves |
| AC5 | issue/SKILL.md body has External kata pointer | yes | **met** | `grep -nE "^## 5\. External kata" ...` → L703; body names the kata path and the three drill scenarios |

### Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` | yes | updated | frontmatter + §5 body replaced with pointer |
| `src/packages/cnos.cdd.kata/katas/M5-issue-authoring/kata.md` | yes | created | M-series shape; three worked examples |
| `src/packages/cnos.cdd.kata/katas/M5-issue-authoring/rubric.json` | yes | created | mirrors M2-review structure |
| `src/packages/cnos.cdd.kata/katas/M5-issue-authoring/baseline.prompt.md` | yes | created | matches M2-review text |
| `src/packages/cnos.cdd.kata/katas/M5-issue-authoring/cdd.prompt.md` | yes | created | CDD mode instructions |
| `docs/gamma/cdd/KATAS.md` | **no** | **gap** | Tier 3 table lists M0–M4; M5-issue-authoring not added |

### CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `.cdd/unreleased/307/self-coherence.md` | yes | yes | gap, mode, ACs, self-check, debt, CDD-Trace, review-readiness all present |
| `.cdd/unreleased/307/beta-review.md` | yes | yes (this file) | written incrementally per review/SKILL.md output format |
| Design artifact | no | n/a | small-change relocation; boundary identical to #304; α named "not required" with reason |
| Plan artifact | no | n/a | three-file sequencing with no shared dependency; α named "not required" with reason |
| Tests | no | n/a | doc/frontmatter change; α named "not applicable" with reason |

### Active Skill Consistency

| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|
| `CDD.md` | Tier 1 | yes | yes | β algorithm, artifact contract, CDD Trace followed |
| `beta/SKILL.md` | Tier 1 | yes | yes | load order, role boundary, incremental write discipline |
| `review/SKILL.md` | Tier 1c | yes | yes | this review |
| `release/SKILL.md` | Tier 1c | yes | loaded | not yet executed; loaded for β closure phase |
| `cnos.core/skills/design` | Tier 3 (issue §Skills) | yes | yes | boundary decision (skill body vs kata package); one source of truth per fact; one reason to change |

---

## Architecture Check

Change touches the boundary between `issue/SKILL.md` (skill body — authoring rules) and `cnos.cdd.kata/katas/M5-issue-authoring/` (kata package — drill scenarios). `cnos.core/skills/design` is active.

| Check | Result | Notes |
|---|---|---|
| Reason to change preserved | yes | `issue/SKILL.md` changes when issue-authoring rules change; kata bundle changes when drill scenario changes — two separate reasons, now separated |
| Policy above detail preserved | yes | policy (authoring rules) stays in skill body; details (drill exercises) move to kata package |
| Interfaces remain truthful | yes | `kata_ref` promises a path; path resolves to a directory containing `kata.md` |
| Registry model remains unified | n/a | no registry in play |
| Source/artifact/installed boundary preserved | n/a | all content is authored source; no build artifact or install layer involved |
| Runtime surfaces remain distinct | yes | skill body / kata package remain distinct; no surface smearing |
| Degraded paths visible and testable | n/a | no degraded paths |

Architecture check: all rows pass or n/a.

---

## Findings

| # | Finding | Evidence | Severity | Type |
|---|---------|----------|----------|------|
| F1 | `docs/gamma/cdd/KATAS.md` Tier 3 table lists M0–M4 only; M5-issue-authoring not listed. The §Debt section in `self-coherence.md` says "No known α debt new this cycle" without naming this gap. The omission is silent rather than explicit. | `KATAS.md` lines 94–100: table rows end at M4 `\| M4 \| Full cycle \| end-to-end CDD loop vs ad hoc on the same change \|`; no M5 row. §Debt: "No known α debt new this cycle." | B | judgment / contract |

**Note on precedent:** #304 (the mirror precedent) updated M2-review content without adding a new M-series row to KATAS.md (M2 pre-existed). This cycle creates a new M5 bundle, so a new row is warranted. KATAS.md is self-described as a "secondary index" — the package is the source of truth — but the catalog is now stale and the stale state is unnamed in debt. Resolution options (either is acceptable): (a) add M5 row to KATAS.md on-branch (2-line fix, 1 commit); (b) append an explicit debt entry in §Debt naming the gap as deferred with reason. Option (a) is preferred; option (b) is acceptable given KATAS.md's secondary-index framing.

## Regressions Required (D-level only)

None (no D findings).

## Notes

- All five AC oracles verified directly against cycle/307 head: AC1 grep → 0 matches; AC2 filesystem → 4 files present; AC3 kata.md headers → match M2-review shape; AC4 frontmatter → `kata_surface: external` + `kata_ref` resolves; AC5 body → `## 5. External kata` at L703.
- M5 kata bundle shape verified against M2-review: same structural headers, same rubric.json format, same baseline.prompt.md text, same cdd.prompt.md structure.
- Design skill applied: source-of-truth split (skill body ↔ kata package) is the correct boundary; one reason to change per surface; no duplication.
- Peer enumeration in self-coherence covers all nine cdd lifecycle skills — none silently inherit the change.
- CI gate: not configured on cycle branches; CI runs on main push (job I5 `skill-frontmatter-check` via `tools/validate-skill-frontmatter.sh`). β accepts α's explicit-not-applicable on row 10 of the pre-review gate.
- The single finding (F1) is resolution-straightforward: a 2-line KATAS.md table addition or an explicit debt entry. No design decision required.

---

## Round 2

**Verdict:** APPROVED

**Round:** 2
**Base SHA (origin/main):** `396d9982af2c2318400fb5c9eb400178d488f99b`
**Head SHA (cycle/307):** `e71831dc3112920bd2bd7f5e9554cf571d919521`
**Fixed this round:** `e71831dc` — closes F1 (KATAS.md stale catalog)
**Branch CI state:** not configured on cycle branches (unchanged from R1; CI runs on main push)
**Merge instruction:** `git merge cycle/307 -m "Closes #307: skill(cdd/issue): move issue katas to cnos.cdd.kata package"` into main

---

## §2.0.0 Contract Integrity — R2 (narrowed to F1 fix scope)

| Check | Result | Notes |
|---|---|---|
| Cross-surface projections updated | **yes** | F1 resolved: `docs/gamma/cdd/KATAS.md` M5 row added; M0–M5 reference updated in Architecture section |

All other R1 contract integrity rows carry forward unchanged — fix-round commit does not touch those surfaces.

## F1 Resolution Verification

**F1:** `docs/gamma/cdd/KATAS.md` Tier 3 table stale (M0–M4 only; M5-issue-authoring not listed; omission unnamed in §Debt).

**Fix commit:** `e71831dc` authored as `alpha@cdd.cnos` — "α(#307): fix-round 1 — add M5 row to KATAS.md (β R1 F1)"

**Verification:**
- Diff `origin/main..origin/cycle/307 -- docs/gamma/cdd/KATAS.md` shows two changes: (1) Architecture bullet: "M0–M4" → "M0–M5"; (2) Tier 3 table: `| M5 | Issue authoring | schema validation + README alignment + witness-theater check |` row added after M4.
- Oracle: `grep -n "M5" docs/gamma/cdd/KATAS.md` → 2 matches (L15: reference update, L89: table row). Both expected matches present.
- Fix scope: only `docs/gamma/cdd/KATAS.md` and `.cdd/unreleased/307/self-coherence.md` (fix-round appendix) changed in `e71831dc`. No frontmatter, schema, skill body, or kata content touched.

**F1 status:** RESOLVED.

## Findings — R2

| # | Finding | Evidence | Severity | Type |
|---|---------|----------|----------|------|
| — | None | F1 resolved; no new incoherence introduced by fix commit | — | — |

## Notes

- R2 scope: verification of F1 fix only. R1 AC coverage, architecture check, and skill consistency remain as assessed in R1 — fix-round commit does not touch those surfaces.
- Fix commit identity: `alpha@cdd.cnos` — role-observable separation maintained in identity-rotation.
- Approval explicitly closes the search space: no remaining blocker found in the contract, implementation, or fix-round scope.
- β proceeds to merge and close-out per `release/SKILL.md`.
