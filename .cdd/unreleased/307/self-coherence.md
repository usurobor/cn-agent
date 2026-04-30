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
