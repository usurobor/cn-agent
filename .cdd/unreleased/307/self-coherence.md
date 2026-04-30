# α self-coherence — #307

skill(cdd/issue): move issue katas to cnos.cdd.kata package

## Gap

**Issue:** #307 — `issue/SKILL.md` §5 carries three embedded katas (5.1 schema validation gate, 5.2 README/source-map alignment, 5.3 checker against witness theater); the skill declares `kata_surface: embedded`. Per repo convention CDD-method drills live in `src/packages/cnos.cdd.kata/katas/M<N>-<name>/`; no other cdd lifecycle skill (`alpha/`, `beta/`, `gamma/`, `design/`, `plan/`, etc.) carries sibling kata content. PR #304 round-5 finding 4 named this cleanup explicitly and deferred it to this follow-up.

**Mode:** small-change (P3, deferred from #304, not blocking).

**Scope shape:**
- in scope: relocate the three kata bodies; flip `issue/SKILL.md` frontmatter to `external` + `kata_ref`; add `## External kata` pointer body section; remove §5 body.
- out of scope: rewriting kata content; touching review katas (already moved in #304); promoting CTB v0.2.

**Boundary decisions taken at issue-author time:**
1. **Option B over Option A** — one bundle (`M5-issue-authoring`) carrying three `## Worked examples`, mirroring #304's `5a8bb3e` shape on `M2-review/`. The active design constraint "Mirror #304's shape" governs absent explicit justification to diverge; none exists.
2. **Bundle includes `rubric.json` + `baseline.prompt.md` + `cdd.prompt.md`** — issue-authoring quality supports method-level baseline-vs-CDD evaluation (the same way M0/M1/M2/M3/M4 do). Defer-with-reason was permitted by the issue's exception-backed strata but unwarranted here; the M-series shape is uniform.
3. **M-series number = 5** — next free integer after M4-full-cycle, matching the strict one-method-per-directory pattern shipped today.

**Failure mode this cycle prevents:**
- convention drift: leaving `issue/` as the only cdd lifecycle skill carrying embedded katas after #304 fixed the symmetric case for `review/` invites the next split to copy the wrong pattern (issue body §Impact).

## Skills

**Tier 1 (CDD lifecycle):**
- `CDD.md` — canonical lifecycle and role contract
- `cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface (this dispatch)
- `cnos.cdd/skills/cdd/issue/SKILL.md` — under modification this cycle; loaded as the affected surface, not as authoring discipline (γ already authored the issue pack)

**Tier 2 (always-applicable engineering):**
- `cnos.eng/skills/eng/markdown` — issue/SKILL.md and kata.md edits
- `cnos.eng/skills/eng/yaml` — frontmatter shape (key addition + enum flip)

**Tier 3 (issue-specific):**
- `cnos.core/skills/design` — boundary decision: skill body owns rules, kata package owns drill scenarios; one source of truth per fact; one reason to change. The active design constraint "Source-of-truth split" governs the §5 → kata-package move; "Mirror #304's shape" picks Option B (one bundle with three worked examples) over Option A (three separate M-series bundles).

The issue's §Skills to load names exactly `cnos.core/skills/design`; Tier 1 + Tier 2 are loaded per α §2.1 step 6 without restatement. No skill is loaded as decoration — design directly governs the §5 boundary call and the option choice.

## ACs

### AC1 — §5 Katas body removed from issue/SKILL.md

- **Invariant:** issue skill carries no embedded kata bodies.
- **Oracle 1:** `grep -nE "^### 5\." src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` → 0 matches.
- **Oracle 2:** `grep -nE "Kata A|Kata B|Kata C" src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` → 0 matches.
- **Evidence:** commit `e61c2615` removed `## 5. Katas` + `### 5.1/5.2/5.3 Kata A|B|C` blocks; §5 now reads `## 5. External kata` with a pointer body.
- **Status:** met.

### AC2 — New kata location(s) exist under cnos.cdd.kata

- **Invariant:** each issue-kata scenario is reachable from `src/packages/cnos.cdd.kata/katas/`.
- **Oracle:** filesystem ls.
- **Evidence:** `src/packages/cnos.cdd.kata/katas/M5-issue-authoring/` exists with `kata.md`, `rubric.json`, `baseline.prompt.md`, `cdd.prompt.md` (commit `a16ad77f`). The bundle carries three `## Worked examples` (Example 1 schema validation gate, Example 2 README/source-map alignment, Example 3 checker against witness theater) — Option B per the issue's structural choice.
- **Status:** met.

### AC3 — Each new kata directory carries kata.md matching M-series shape

- **Invariant:** kata.md follows M0–M4 form (class / level / purpose / scenario / required artifacts).
- **Oracle:** required headers present.
- **Evidence:** `M5-issue-authoring/kata.md` opens with `**Class:** method`, `**Default level target:** L6`, `**Purpose:**`, `## Scenario`, `## Required artifacts`, `## Scoring`, `## Worked examples` — identical structural shape to `M2-review/kata.md` (the precedent set in #304's `5a8bb3e`).
- **Status:** met.

### AC4 — issue/SKILL.md frontmatter declares external kata

- **Invariant:** `kata_surface: external`; `kata_ref` points at the new kata location.
- **Oracle:** frontmatter parse.
- **Evidence:** L5 reads `kata_surface: external` (was `embedded`); L26 reads `kata_ref: src/packages/cnos.cdd.kata/katas/M5-issue-authoring/`. The path resolves to a directory containing `kata.md`. Schema (`schemas/skill.cue`) accepts `kata_ref` via the open-schema trailing `...` (LANGUAGE-SPEC §11 — loaders MUST ignore unknown keys); `review/SKILL.md` uses the same key shape.
- **Status:** met.

### AC5 — issue/SKILL.md body has External kata pointer

- **Invariant:** a human-readable section names the kata location and what it exercises.
- **Oracle:** `grep -nE "^## 5\. External kata" src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` → 1 match (L703).
- **Evidence:** `## 5. External kata` body section names the path `src/packages/cnos.cdd.kata/katas/M5-issue-authoring/`, lists the three drill scenarios it exercises, and references the frontmatter `kata_ref` for machine-readable linkage — mirroring the `## External kata` section that `review/SKILL.md` carries (L249–L255).
- **Status:** met.
