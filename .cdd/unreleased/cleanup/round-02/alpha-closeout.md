# α Close-out — Docs/Structure Remediation (Round 02)

**Role:** α (implementer)
**Scope:** DOCS/STRUCTURE only (docs/ tree + root `*.md` + nav/links). Source, OCaml, `scripts/`, and package `SKILL.md`/`*.cue` behavior deferred to the code pass.
**Input:** `round-02/beta-review.md` (F8–F11).
**Files changed:** 11 (all under `docs/`). No commit — κ commits.

---

## Finding → action

| # | Sev | Status | Files | What changed |
|---|-----|--------|-------|---------------|
| F8 | HIGH | **FIXED (rewrite)** | `docs/reference/governance/DOCUMENTATION-SYSTEM.md` | Rewrote the doc around one governing question — "how the cnos docs tree is organized and where a given document lives." Removed the retired α/β/γ folder taxonomy teaching entirely (§1 triad-axis table, §1.2 `docs/alpha/{scope}/` bundle scheme, §6 alpha-root legacy migration, §7/§9 α/β/γ placement). Replaced with the real intent-directory layout (`quickstart/ concepts/ guides/ reference/ architecture/ development/ papers/ evidence/`, verified against the tree + `docs/README.md`). Kept the still-valid concepts re-homed onto reality: feature-bundle contract (README names one canonical spec), single cnos-release version lineage, supersession. Corrected the version-freeze section to describe reality — **no `X.Y.Z/` snapshot directories exist anywhere under `docs/`, and `docs/README.md` states frozen history lives in git, not folders on `main`** — so the old version-directory machinery was cut, not re-homed. Added the α/β/γ-is-a-measurement-not-folders statement (matches `README.md`/`GLOSSARY.md`). Header bumped to `Version: 3.82.0` (current release, per VERSION). **361 → 122 lines**; zero `docs/alpha/`. |
| F9 | MED | **FIXED (8-row swap)** | `architecture/INVARIANTS.md`, `architecture/HUB-PLACEMENT-MODELS.md`, `papers/COHERENCE-SYSTEM.md`, `reference/ctb/CTB-v4.0.0-VISION.md` (×2), `reference/ctb/LANGUAGE-SPEC.md` (×2), `reference/ctb/SEMANTICS-NOTES.md`, `development/design/README.md` | Applied β's F9 table verbatim; every target verified present before writing. Two self-referential canonical-path headers fixed first (`INVARIANTS.md:6` → `docs/architecture/INVARIANTS.md`; `HUB-PLACEMENT-MODELS.md:468` → `docs/architecture/HUB-PLACEMENT-MODELS.md`). `docs/alpha/doctrine/FOUNDATIONS.md` → `docs/papers/FOUNDATIONS.md`; the four `SKILL-ARCHITECTURE.md` cites → `docs/papers/SKILL-ARCHITECTURE.md`; `COHERENCE-FOR-AGENTS.md` → `docs/concepts/doctrine/coherence-for-agents/COHERENCE-FOR-AGENTS.md`; `docs/alpha/essays/` → `docs/papers/`. |
| F10 | MED | **FIXED (3 files repointed + 1 banner)** | `architecture/cognitive-substrate/CAR.md` (×2), `papers/COHERENCE-SYSTEM.md`, `reference/runtime/ORCHESTRATORS.md` (×3), `reference/packages/PACKAGE-RESTRUCTURING.md` (banner) | Replaced the phantom `cnos.pm` on the live surface with real shipping packages, verified against `src/packages/` + README's package table (no invention). CAR §3.1 role-pack list `cnos.eng, cnos.pm` → `cnos.eng, cnos.cds`; CAR §5 `cnos.pm/` tree block → `cnos.cds/` with its real skills (`cds/lifecycle/SKILL.md, cds/selection/SKILL.md`). COHERENCE-SYSTEM §5.5 `cnos.pm` → `cnos.cds`. **ORCHESTRATORS.md examples repointed to `cnos.core`, not β's suggested cnos.cds/cnos.eng** — the `daily`/`weekly` commands and the `daily-review` orchestrator in those examples are actually declared by `cnos.core` (`src/packages/cnos.core/cn.package.json` + `orchestrators/daily-review/`), so `cnos.core` is the source-verified owner. PACKAGE-RESTRUCTURING.md (Draft #186): added a status banner ("Draft proposal — not the shipped structure; see README package table"), body left intact per β's special-case. |
| F11 | LOW/policy | **APPLIED** | — (decision + policy application below) | Applied β's ratified rule: a doc is frozen-and-left iff it carries a dated/superseded/status header AND is not linked as current authority; only the live surface is corrected. No per-file "paths historical" notes added. |

---

## F11 policy application — frozen-left vs live-fixed

**Live-and-fixed** (cited as current authority or on the active reader surface):
- `DOCUMENTATION-SYSTEM.md` — cited by README + `reference/governance/README.md` as "single source of truth"; failed the second clause → rewritten (F8).
- The 8 F9 reference-surface files (INVARIANTS, HUB-PLACEMENT, COHERENCE-SYSTEM, the three CTB docs, design/README) — bare code-span nav paths on live docs → swapped.
- The F10 live specs (CAR, COHERENCE-SYSTEM §5.5, ORCHESTRATORS) → phantom package repointed.
- `PACKAGE-RESTRUCTURING.md` — Draft under `reference/packages/` that reads as authoritative; its `cnos.pm` mentions are intentional (it proposes retiring it), so per β it got a **status banner** rather than a body edit — the minimum marker that makes it qualify as clearly-proposal, not live structure.

**Frozen-and-left** (dated/superseded/status header + not live authority — untouched):
- `CHANGELOG.md` (dated release ledger), `.cn-sigma/logs/*` (dated peer logs).
- `docs/development/plans/{PLAN-v3.13.0-docs-governance, PR-docs-governance-v3.13.0}.md` (dated 2026-03-23), `docs/development/plans/CAR-implementation-plan.md` (`## Status: Implemented`, not linked as authority).
- `docs/evidence/design/WRITER-PACKAGE.md` (`Superseded` header; β's 25-hit leave-list).
- `docs/reference/packages/DESIGN-266-dist-out-of-git.md` (issue-#266 design record; its one `docs/alpha/` hit describes the *content* of a frozen INVARIANTS line, not a nav path — β-ratified frozen).
- `docs/papers/EXECUTABLE-SKILLS.md:398` (published-paper bibliography citing `WHITEPAPER.md`, which exists nowhere in the repo — β flagged leave-unless-locatable; not locatable, left as a frozen reference-list entry).
- Intentional "these were retired" statements naming the folders to forbid them: `GLOSSARY.md:18` and the two new lines I wrote in `DOCUMENTATION-SYSTEM.md:28,121`. These are correct usages, not stale citations.

---

## Verification (live surface)

### `docs/alpha/` — cleared on the live reference surface

`git grep -n 'docs/alpha/' -- docs/ '*.md' ':!.cdd/'` → every survivor is frozen-or-intentional:
- Frozen records: `CHANGELOG.md` (8), `.cn-sigma/logs/*` (3), `plans/PLAN-v3.13.0-docs-governance.md`, `plans/PR-docs-governance-v3.13.0.md`, `evidence/design/WRITER-PACKAGE.md` (25), `papers/EXECUTABLE-SKILLS.md:398` (bibliography), `reference/packages/DESIGN-266-dist-out-of-git.md:75`.
- Intentional retirement statements: `reference/governance/GLOSSARY.md:18`, `reference/governance/DOCUMENTATION-SYSTEM.md:28,121`.
- **Zero live-authority bare-path `docs/alpha/` citations remain.** All 8 F9 rows + F8's 12 in-file hits gone.

### `cnos.pm` — cleared on the live reference surface

`git grep -n 'cnos\.pm' -- docs/ '*.md' ':!.cdd/'` → survivors are frozen or the bannered Draft:
- Frozen: `CHANGELOG.md` (3, dated history), `docs/development/plans/CAR-implementation-plan.md` (3, `Status: Implemented` plan).
- Bannered Draft: `docs/reference/packages/PACKAGE-RESTRUCTURING.md` (7 — intentional; it proposes retiring `cnos.pm`; now carries the not-shipped banner, incl. the banner's own reference).
- **Zero phantom `cnos.pm` on live specs** — CAR / COHERENCE-SYSTEM / ORCHESTRATORS all repointed.

### Regression — no F1/F2 reintroduction

- `git grep -n 'src/agent/' -- docs/ '*.md' ':!.cdd/' | wc -l` → **43, unchanged** from the Round-1 baseline (every hit still a CHANGELOG/plan/banner/package record; I touched no `src/agent/` line).
- README / CLI.md / OPERATOR.md untouched this round → still agree on shipped-vs-planned.
- All 8 F9 swap targets + the F10 `cnos.core`/`cnos.cds` names verified present on disk before writing → no new broken paths.

---

## Newly-discovered drift (out of F8–F11 scope — flag for a future β round)

A sibling `docs/beta/` / `docs/gamma/` folder-path cluster exists on the live surface (β's F9 tabled only `docs/alpha/`). Examples: `development/design/README.md:23-24` (`docs/gamma/essays/`, `docs/gamma/cdd/`), `concepts/lineage/{LINEAGE,ORIGIN}.md` canonical-path headers (`docs/beta/...`), `papers/{CCNF-AND-TYPED-TRUST,RELEASE-LEVEL-CLASSIFICATION}.md` (`docs/gamma/ENGINEERING-LEVELS.md`), `reference/ctb/CTB-v4.0.0-VISION.md:341`. **Not remediated** — it is not in F8–F11, and critically `docs/gamma/cdd/{X.Y.Z}/` is a *live CDD snapshot-path convention* referenced across active `src/packages/cnos.cdd` / `cnos.cds` skill files (out of scope, code-entangled), so a correct sweep must be coordinated with the code pass, not blind-swapped in docs alone. Same class as F9; recommend β file it and split the pure-docs `docs/beta`/`docs/gamma/essays` rows from the code-entangled `docs/gamma/cdd` convention.
