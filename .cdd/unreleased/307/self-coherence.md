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

## Self-check

**Did α push ambiguity onto β?** No.
- Boundary decisions taken: Option B over Option A (named in §Gap with reason); kata bundle includes full prompt-pair + rubric (named with reason); M-series number = 5 (named with reason).
- All five ACs map to a concrete oracle and a single piece of evidence (commit SHA + line numbers).
- Peer enumeration completed (see below); no sibling skill silently inherits the move.

**Is every claim backed by evidence in the diff?** Yes.
- AC1: grep counts in §ACs (run + reported = 0/0).
- AC2/AC3: filesystem layout shipped in `a16ad77f`.
- AC4: frontmatter lines L5 + L26 in `e61c2615`.
- AC5: `## 5. External kata` body section L703 in `e61c2615`.

**Peer enumeration — kata_surface peers (skill-class peers).**
The diff modifies frontmatter `kata_surface` on `issue/SKILL.md`. The peer set is every cdd lifecycle skill with a `kata_surface` declaration:

| Skill | kata_surface | kata_ref | Status this cycle |
|---|---|---|---|
| `cdd/alpha/SKILL.md` | embedded | (none) | exempt — embedded `## 4. Embedded Kata` (L316–L344) intentional; not a §307 target |
| `cdd/beta/SKILL.md` | embedded | (none) | exempt — out of scope (issue §Non-goals limits to issue katas) |
| `cdd/gamma/SKILL.md` | embedded | (none) | exempt — out of scope |
| `cdd/design/SKILL.md` | embedded | (none) | exempt — out of scope |
| `cdd/plan/SKILL.md` | none | (none) | exempt — `kata_surface: none`; nothing to relocate |
| `cdd/issue/SKILL.md` | external (was embedded) | M5-issue-authoring/ | **updated this cycle** |
| `cdd/review/SKILL.md` | external | M2-review/ | unchanged — moved in #304 (`5a8bb3e`); the precedent this cycle mirrors |
| `cdd/release/SKILL.md` | embedded | (none) | exempt — out of scope |
| `cdd/post-release/SKILL.md` | embedded | (none) | exempt — out of scope |

Issue §Non-goals explicitly excludes "Touching review katas (already moved)" and limits scope to the three issue katas; β/γ/design/plan/release/post-release skills are not in this cycle's scope. No skill silently inherits the convention change.

**Peer enumeration — `## External kata` body-section peers.**
`review/SKILL.md` is the precedent (§External kata at L249). `issue/SKILL.md` now matches its structural shape. No other lifecycle skill currently carries a §External kata section because no other has external kata content yet — adding the section to skills without an external kata target would be a future-as-present claim.

**Peer enumeration — kata_ref consumers.**
Searched the repo for code that reads `kata_ref` (machine-readable consumers). None present today; the field is documentation-shaped per the open-schema §11 contract. No downstream surface needs updating.

**Harness audit.** Schema (`schemas/skill.cue`) is open via trailing `...`; `kata_ref` passes through without schema change required (review/SKILL.md ships with the same key already). `tools/validate-skill-frontmatter.sh` validates via `cue vet` and is unaffected. No CI workflow, fixture, or generator emits skill frontmatter — no harness drift surface exists.

**Intra-doc repetition check.** `grep -nE "Kata A|Kata B|Kata C|^### 5\." src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` → 0 matches (verified post-commit). The §5 body-removal claim is exhaustive against the doc, not just the dominant heading.

**Closure-overclaim check.** The cycle does not claim "all CDD lifecycle skills now carry external katas" — only that `issue/SKILL.md` now does, mirroring `review/SKILL.md` per #304. Issue §Non-goals are preserved unchanged.

## Debt

**Carried into this cycle (none new from this cycle):**

- **Other cdd lifecycle skills with `kata_surface: embedded`** — `alpha/`, `beta/`, `gamma/`, `design/`, `release/`, `post-release/` retain embedded kata bodies. This cycle is explicitly out-of-scope for them per issue §Non-goals; whether they should also move externalises is a separate γ decision, not α-side debt.
- **No `rubric.json` / `baseline.prompt.md` / `cdd.prompt.md` deferred** — the issue's exception-backed strata permitted deferring these per-kata if "the content shape doesn't yet warrant baseline-vs-CDD evaluation," but the M-series shape ships them uniformly and the issue-authoring drill supports method-level evaluation. M5 carries the full quartet, mirroring M2-review.
- **Per-scenario rubric and prompt-pair files not split** — Option B's three `## Worked examples` share one rubric.json and one prompt-pair (the M2-review precedent). If a future cycle wants per-scenario evaluation, that is a separate split; the issue's "Deferred" entry covers exactly this case.

**No known α debt new this cycle.**

## CDD-Trace

CDD canonical artifact order (CDD.md §5.2) mapped to this cycle:

1. **Design artifact** — *not required.* Justification: small-change relocation; the boundary call ("skill body owns rules; kata package owns drill scenarios; one source of truth per fact") is identical to the boundary #304 already decided and shipped for `review/SKILL.md` (`5a8bb3e`). The issue's §Active design constraints make this explicit ("Mirror #304's shape"). No new boundary is being established.
2. **Coherence contract** — `.cdd/unreleased/307/self-coherence.md` §Gap (commit `8c857ce2`). Names the gap, mode, scope shape, three boundary decisions (Option B / full prompt-pair / M5 numbering), and the failure mode (convention drift) the cycle prevents.
3. **Plan** — *not required.* Justification: the implementation sequencing is "create kata bundle; flip frontmatter; remove §5 body; add pointer section" — three isolated edits with no shared dependency or sequencing risk. The issue §Implementation guidance names the file paths and the pattern; no further plan adds executable detail.
4. **Tests** — *not applicable.* Justification: the change is documentation/runbook content + frontmatter shape, not code or runtime. Validation lives in `tools/validate-skill-frontmatter.sh` (CI job I5) which is invariant under this change — `kata_ref` passes the open-schema (`schemas/skill.cue` trailing `...`), `kata_surface: external` is in the enum (LANGUAGE-SPEC §2.1).
5. **Code** — n/a (no code surface touched).
6. **Docs** — issue/SKILL.md frontmatter + §5 body replaced with pointer (commit `e61c2615`); kata.md + rubric.json + prompt-pair created in M5-issue-authoring/ (commit `a16ad77f`).
7. **Self-coherence** — this file. Sections written incrementally per α §2.5: §Gap (`8c857ce2`), §Skills (`9ccfa9f5`), §ACs (`9d430391`), §Self-check (`90a90b80` + `73cebc9c` placeholder-fix), §Debt (`4d91ff04`), §CDD-Trace (this commit), §Review-readiness (next commit).

**Branch commits in cycle order:**

| # | SHA | Author email | Subject |
|---|---|---|---|
| 1 | `8c857ce2` | alpha@cdd.cnos | self-coherence §Gap |
| 2 | `a16ad77f` | alpha@cdd.cnos | create M5-issue-authoring bundle |
| 3 | `e61c2615` | alpha@cdd.cnos | flip frontmatter + replace §5 body |
| 4 | `9ccfa9f5` | alpha@cdd.cnos | self-coherence §Skills |
| 5 | `9d430391` | alpha@cdd.cnos | self-coherence §ACs |
| 6 | `90a90b80` | alpha@cdd.cnos | self-coherence §Self-check |
| 7 | `73cebc9c` | alpha@cdd.cnos | self-coherence §Self-check placeholder-fix |
| 8 | `4d91ff04` | alpha@cdd.cnos | self-coherence §Debt |

All commits authored as `alpha@cdd.cnos` per CDD.md §1.4 α step 2 — role identity is git-observable across the cycle.

## Review-readiness | round 1

- **Base SHA (origin/main HEAD at signal time):** `396d9982af2c2318400fb5c9eb400178d488f99b`
- **Head SHA (cycle/307 HEAD at signal time):** *(this commit; appended last)*
- **Merge-base HEAD ↔ origin/main:** `396d9982` (= origin/main HEAD — cycle is rebased; no drift).
- **Signal time:** 2026-04-30 08:23 UTC.
- **Branch CI on cycle/307:** **not configured** — `.github/workflows/build.yml` triggers only on `push: branches: [main]` and `pull_request: branches: [main]`; cycle branches are not gated by GitHub-side CI. Per α §2.6 row 10, this artifact says so explicitly; β waits for green at the post-merge `main` push (workflow includes I5 `skill-frontmatter-check` which validates the new `kata_ref` frontmatter via `cue vet schemas/skill.cue` + `tools/validate-skill-frontmatter.sh`). Local CI is also unavailable in α's environment (`cue` binary not installed in the sandbox) — the schema is open via trailing `...` and `kata_ref` is the same key shape `review/SKILL.md` already ships with, so the schema-conformance assertion is structural rather than executed at signal time.

**Pre-review gate (α §2.6) at signal time:**

| # | Row | State | Evidence |
|---|---|---|---|
| 1 | cycle branch rebased onto current origin/main | **pass** (transient) | merge-base = origin/main HEAD = `396d9982` |
| 2 | self-coherence.md carries CDD Trace through step 7 | **pass** (durable) | §CDD-Trace section above |
| 3 | tests present or explicit reason none apply | **pass** (durable) | §CDD-Trace step 4 — n/a (doc/runbook + frontmatter, not code/runtime) |
| 4 | every AC has evidence | **pass** (durable) | §ACs AC1–AC5 |
| 5 | known debt explicit | **pass** (durable) | §Debt — no new α debt |
| 6 | schema/shape audit when contracts changed | **pass** (durable) | §CDD-Trace step 4 — `schemas/skill.cue` open via `...`; `kata_ref` precedent in `review/SKILL.md` |
| 7 | peer enumeration when closure claim touches a family | **pass** (durable) | §Self-check kata_surface peers table (9 lifecycle skills) |
| 8 | harness audit when schema-bearing contract changed | **pass** (durable) | §Self-check — `tools/validate-skill-frontmatter.sh` invariant; no other frontmatter writers |
| 9 | post-patch re-audit covering every language in diff | **pass** (durable) | diff languages = Markdown + YAML frontmatter only; Markdown grep oracles run (AC1); YAML frontmatter shape verified against `schemas/skill.cue` |
| 10 | branch CI green on head commit | **explicit-not-applicable** (transient) | see "Branch CI on cycle/307" note above; β waits for green at post-merge main push |
| 11 | α commit author email matches `alpha@cdd.{project}` | **pass** (durable) | `git log --format=%ae origin/main..HEAD | sort -u` → `alpha@cdd.cnos` (single value across all 8 cycle commits per §CDD-Trace branch-commit table) |

**Ready for β.**
