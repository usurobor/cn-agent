# β Review — Round 05 (full adversarial re-review; frozen-bucket ruling reversed)

**Role:** β — independent adversarial newcomer reviewer.
**Branch:** `claude/repo-cleanup-newcomer`.
**Method:** Walked the whole tree (root, `docs/`, `src/packages/`, `.cdd/`, `.cn-sigma/`, `schemas/`, `scripts/`, `tests/`). Every finding below names a path I opened or grep'd. Edited only this file.
**Mandate delta:** Rounds 1–4 converged under a rule that PRESERVED a "frozen record = leave" bucket (`round-03/beta-review.md` F16, `round-04/closeout.md` §Convergence). That ruling is reversed. Everything the prior loop called "frozen, leave it" is re-judged here as legacy narrative unless it is current-state doctrine.

---

## 1. Verdict

**REVISE.** The newcomer *reading path* is genuinely strong, but the repo contradicts its own stated law: `docs/README.md:10` and `docs/papers/README.md` declare "frozen history is not kept on `main` — git history is the archive," while the tree ships **~950 tracked `.cdd/` files (13 MB), 41 dated `.cn-sigma/` activation logs (2.2 MB), 16 dated PLAN docs, a 180 KB CHANGELOG, and a shelf of dated audits/assessments/RCAs.** The prior loop froze all of it. It is the archive, in the tree.

---

## 2. Newcomer first-impression

Cold read, top-down. **What's clear:** `README.md` lands the point in one line ("cnos gives AI agents a Git-native home"), separates *what ships* (`cn` hub/package commands) from *planned* runtime without contradiction, and draws an accurate tree map. `docs/README.md` is a clean intent-nav (`quickstart/ concepts/ guides/ reference/ development/ architecture/ papers/ evidence/`). `src/packages/cnos.*` is a legible, best-practice monorepo package layout. The primary spine — README → THESIS → concepts → papers/reference — reads coherently. A newcomer knows what cnos IS within a minute. Credit where due.

**What's noise the moment you look past the spine:** the root `README.md` tree map itself points a newcomer straight into two archives — `.cdd/` ("CDD work cells, receipts, and release evidence") and `.cn-sigma/` ("foreign-activation footprint"). Open either and you find hundreds of dated per-cycle work records and daily heartbeat logs — the working scratchpad of the *process that built the repo*, not the repo. `docs/development/plans/` is 16 implementation plans for shipped work, most stamped `Status: Complete`/`Draft`. The root carries 11 markdown files including a 510-line protocol treatise (`ROLES.md`) and a single-release changelog snapshot (`RELEASE.md`). The signal spine is clean; it floats on a large sediment of history the repo elsewhere promises it deleted.

---

## 3. Benchmark (named comparators)

| Comparator | What it does that cnos doesn't | cnos delta |
|---|---|---|
| **TigerBeetle** | `docs/` is one current-state narrative spine; design docs describe the system as it is; zero in-tree dated per-cycle work-record directories — process history lives in Git/PRs. | cnos ships `.cdd/releases/**` = **596 files** of α/β/γ closeout/review/self-coherence/scaffold records for already-shipped cycles. |
| **SQLite** | Famously lean tree; canonical docs are current-state; history lives in the Fossil timeline + one changelog page, never as hundreds of plan files. | cnos ships **16** `docs/development/plans/PLAN-*.md`, dated and version-keyed (`PLAN-v3.22.0…` `Status: Complete`). |
| **Redis** | Lean root; release notes in one generated file; no per-release directories; no hand-scored "quality ledger." | cnos root has `RELEASE.md` (one release, 3.82.0) **plus** a 180 KB / 1352-line `CHANGELOG.md` carrying a hand-authored "Release Coherence Ledger" of intuition-grade TSC letters. |
| **Zig** | README is the on-ramp; roadmap/futures live in issues and a milestone, not `*-VISION.md` files in the doc tree. | cnos ships `CTB-v4.0.0-VISION.md` (`Status: Draft for alignment`) and `MODULAR-ARCHITECTURE-REFACTOR.md` (`Status: Proposed`) as tree docs. |
| **Deno / Bun / Biome** | Monorepo packages under one directory, each self-describing; no sibling archive of build-process records. | cnos `src/packages/` **matches this well** — the package layout is a genuine strength; the archive lives *outside* `src/`, in `.cdd/`. |
| **Astro** | README newcomer on-ramp: what it is, install, first run, where to go. | cnos README **matches this well** — the on-ramp is strong; the delta is everything below the spine. |

**Net:** cnos's *entry surface* is at parity with Astro/Zig/TigerBeetle. Its *history discipline* is the opposite of all six — it keeps in the working tree precisely the process artifacts those repos push to Git history, issues, and generated files.

---

## 4. Findings

Ordered most-severe first, content then structure. "Legacy class" tags the mandate's enumerated buckets.

### CONTENT

#### C1 — HIGH — The repo violates its own stated archival law
*Legacy class: governing.* `docs/README.md:10` — "Frozen history is not kept on `main` — git history is the archive." `docs/papers/README.md` — "Superseded drafts and frozen snapshots are not kept on `main` — git history is the archive." These are current-state doctrine and correct. The tree then ships the exact opposite: `.cdd/releases/**` (596 files), `.cn-sigma/logs/*.md` (41 dated), `docs/development/plans/PLAN-*.md` (16), `docs/development/cdd/POST-RELEASE-EPOCH-v3.12.md` + `-v3.14.md`, `docs/papers/ENGINEERING-LEVEL-ASSESSMENT.md`, `docs/evidence/AUDIT.md`, `docs/evidence/rca/**` (8), `docs/development/kata/runs/**` (15). **Fix:** this finding is the frame for C2–C8 and S1–S2 — either enforce the stated law (remove the history) or the doctrine is false. Enforce the law.

#### C2 — HIGH — Dated implementation plans for shipped work
*Legacy class: dated plan documents.* `docs/development/plans/` (16 files). Evidence of shipped status: `PLAN-v3.22.0-eng-lane-clarity.md` `Status: Complete`; `PLAN.md` `Status: Complete`; `issue-41-pass-b-wiring.md` `Status: Complete`; `CAR-implementation-plan.md` `## Status: Implemented`; `TRACEABILITY-implementation-plan.md` `Status: Implementing (Steps 1-9 complete)`. A newcomer needs the *current* design of these subsystems (which lives in `docs/reference/` + `docs/architecture/`), not the plan that produced them. **Fix:** delete the `Complete`/`Implemented` plans (git history holds them). For any still-live `Draft` plan whose target isn't yet built, relocate the design content into the owning `docs/reference/*` doc and drop the dated plan wrapper. Net target: the `plans/` directory ceases to exist.

#### C3 — HIGH — Historical assessments, audits, RCAs, and kata runs kept as files
*Legacy class: historical PRAs / audit logs / "how we got here."* `docs/papers/ENGINEERING-LEVEL-ASSESSMENT.md` (opens `**Period**: Feb 2 – Mar 25, 2026 … **Commits**: 1,434`) — a point-in-time performance audit, not a paper. `docs/development/cdd/POST-RELEASE-EPOCH-v3.12.md` + `-v3.14.md` ("retroactive epoch assessment … Individual per-release assessments were not performed at release time"). `docs/evidence/AUDIT.md` (`**Date:** 2026-02-11`, "17 design docs audited"). `docs/evidence/rca/**` (8 dated incident write-ups). `docs/development/kata/runs/**` (15 dated run/score records). Every one answers "how we got here." **Fix:** delete from tree (git history is the archive per C1). If the *rubric* they exercised is doctrine, that already lives in `ENGINEERING-LEVELS.md` / the kata templates — keep the rubric, drop the run logs.

#### C4 — MED — `RELEASE.md` at root is a one-release narrative snapshot
*Legacy class: RELEASE as long-form narrative.* `RELEASE.md` opens `# 3.82.0` and narrates one release (Outcome / Why it matters / Added / Changed) — duplicating the top `CHANGELOG.md` entry and the `3.82.0` git tag. It also cites retired dirs (`grep -c docs/gamma RELEASE.md` → 2), so it is stale narrative *and* has dead links. The prior loop ruled it "FROZEN — leave" (`round-03/beta-review.md` F16). Reversed. **Fix:** delete from root. One release's notes belong in `CHANGELOG.md` (or a git-tag body), not a standing root file that always describes whatever release last touched it.

#### C5 — MED — CHANGELOG is 180 KB of hand-authored narrative + invented scores
*Legacy class: CHANGELOG as long-form narrative.* `CHANGELOG.md` is 1352 lines / 180 KB. Beyond version entries it carries a "Release Coherence Ledger" scoring every release on TSC letter grades — and the file itself admits "These are intuition-level ratings, not outputs from a running TSC engine." That is opinion presented as measurement (write skill 3.9). **Fix:** trim to a lean Keep-a-Changelog list (version → user-visible changes), or generate it from tags. Drop the intuition-grade ledger entirely — a hand-authored score is an invented score, which `docs/README.md:50` itself forbids for frontmatter.

#### C6 — MED — Future-state VISION / Proposed / Draft docs in the doc tree
*Legacy class: `*-VISION.md` / futures / dead states.* `docs/reference/ctb/CTB-v4.0.0-VISION.md` (`Status: Draft for alignment (intended to guide the next 2–4 major iterations)`) describes a language that does not exist. `docs/reference/ctb/LANGUAGE-SPEC-v0.2-draft.md` is a superseded draft sitting beside the real `LANGUAGE-SPEC.md`. `docs/papers/MODULAR-ARCHITECTURE-REFACTOR.md` (`Status: Proposed`, `Date: 2026-04-08`, "code layout is still monolithic") is a decision record for a refactor. `reference/` is where a newcomer looks for canonical current spec; futures and superseded drafts mislead there. **Fix:** delete the superseded `-v0.2-draft`; move genuine roadmap (CTB v4 vision, modular refactor intent) to a GitHub issue/milestone, not a spec-adjacent doc. Keep only the shipped `LANGUAGE-SPEC.md`.

#### C7 — MED — "This migration has landed / X is gone" retirement narrative
*Legacy class: retirement notes.* `docs/reference/packages/BUILD-AND-DIST.md:202` — "> This migration has landed. `src/agent/` is gone; packages are authored under …". Also the forbid-statements in `docs/reference/governance/DOCUMENTATION-SYSTEM.md` and `GLOSSARY.md` narrating retired `docs/beta/`,`docs/gamma/`. The prior loop left these as "intentional" (`round-03/beta-review.md` F16). Reversed: a current-state doc states what the layout **is**; it does not narrate the migration that produced it (write skill 3.1 — state what a thing is, not what it is not). **Fix:** rewrite each to the positive current fact ("packages are authored under `src/packages/`") and drop the "used to be / is gone / was retired" clause. Git history explains the transition.

#### C8 — LOW — Inert authoring-intent prose pointing at dead dirs
*Legacy class: authoring-intent prose.* `docs/papers/CCNF-AND-TYPED-TRUST.md:383` ("Land this document under `docs/gamma/essays/`…") and `docs/papers/DECREASING-INCOHERENCE.md:528` ("Add this file under `docs/gamma/essays/`…"). Both instructions are already satisfied (the docs live in `docs/papers/`) and both name a retired directory. Prior loop left them as "frozen." **Fix:** delete both lines — an executed instruction to file a file is pure process residue.

### STRUCTURE

#### S1 — HIGH — `.cdd/` ships a 13 MB per-cycle process archive
*Legacy class: `.cdd/iterations/**` + `.cdd/unreleased/**` + release records.* `git ls-files .cdd | wc -l` → **950**; `du -sh .cdd` → **13 MB**. Breakdown: `.cdd/releases/**` = **596** files (per shipped cycle: `alpha-closeout.md`, `beta-review.md`, `gamma-scaffold.md`, `self-coherence.md`, `cdd-iteration.md` …); `.cdd/unreleased/**` holds cells for already-shipped issue numbers; `.cdd/iterations/**` and `.cdd/waves/**` are dated dispatch/wave records; `.cdd/proposals/gamma-iteration-after-3.61.0.md` is a spent proposal; `.cdd/legacy-exceptions.yml` exists *only* to grandfather "Historical gaps in released cycles" — self-evidence the directory is a history store. **Distinguish the live substrate:** `.cdd/CADENCE`, `CDD-VERSION`, `DISPATCH`, `OPERATORS`, `exceptions.yml` are live control config (analogous to `.github/`) and stay. The **596-file `releases/` archive, spent `unreleased/` cells, `iterations/`, `waves/`, `proposals/`** are the frozen history the repo says lives in Git. **Fix:** remove `.cdd/releases/`, `.cdd/iterations/`, `.cdd/waves/`, `.cdd/proposals/` from the tree; gitignore in-flight `.cdd/unreleased/<issue>/` working cells (keep them local, not shipped). Keep only the live control files + `unreleased/cleanup/` for this active loop. This is the single highest-impact removal in the repo.

#### S2 — HIGH — `.cn-sigma/logs/` ships 41 dated heartbeat logs (2.2 MB)
*Legacy class: agent activation logs / "how we got here."* `ls .cn-sigma/logs/2*.md` → **41** dated files (`20260530.md` … `20260709.md`); a sampled entry (`20260709.md`) is a per-wake activation heartbeat ("pure heartbeat; no change since prior wake"). `du -sh .cn-sigma` → **2.2 MB**. This is one foreign agent's operational run history, surfaced to newcomers via the root README tree map (line 125). The live doctrine (`.cn-sigma/README.md`) explaining the namespace is legitimate and stays; the dated logs are archive. **Fix:** remove `.cn-sigma/logs/2*.md` from the tree (they belong in the agent's home hub / Git history, per the namespace's own "per-context state only" framing); keep `.cn-sigma/README.md`.

#### S3 — MED — Root markdown proliferation
*Legacy class: root noise.* Root carries 11 `.md` + `VERSION`. Load-bearing for a newcomer and standard-located: `README.md`, `CONTRIBUTING.md`, `LICENSE`, `SECURITY.md`, `CODE_OF_CONDUCT.md`. Noise or misplaced: `ROLES.md` (510 lines — a full protocol treatise on the α/β/γ/δ/ε ladder; belongs in `docs/development/` or `docs/concepts/`, not root); `OPERATOR.md` (240 lines day-2 ops manual → `docs/guides/` or `docs/reference/cli/`); `RELEASE.md` (delete, C4); `SUSTAINABILITY.md` (→ fold into `.github/FUNDING` + a `docs/` page); `CHANGELOG.md` (trim, C5). `THESIS.md` is correctly under `docs/` with no root copy — good. **Fix:** relocate `ROLES.md`, `OPERATOR.md`, `SUSTAINABILITY.md` under `docs/`; delete `RELEASE.md`. Target root: `README CONTRIBUTING LICENSE SECURITY CODE_OF_CONDUCT CHANGELOG(lean) VERSION` + config.

#### S4 — LOW — `docs/reference/legacy/` is a dir that exists only to hold history
*Legacy class: history-only directory.* `docs/reference/legacy/` contains one file, `OCAML-THREAD-REFERENCE.md`, a pointer to an archived branch/tag. The root README already carries the identical "Legacy: OCaml impl archived off main" note (README.md:134–138). The directory duplicates that pointer and its name announces a history-only purpose. **Fix:** delete the directory; the README's Legacy paragraph is the one home for that fact (write skill 3.3).

#### S5 — LOW — `docs/development/board/` ships a kanban snapshot
*Legacy class: work-record snapshot.* `docs/development/board/` = `board-data.json` + `index.html` — a point-in-time project-board render. A newcomer doc tree should not carry a frozen board dump (it is stale the moment it lands, and the live board is GitHub). **Fix:** remove from tree; link the live GitHub project board from `docs/development/README.md` if wanted.

---

## 5. Proposed ideal minimal tree

```text
cnos/
├── README.md              identity, on-ramp, map, pointers
├── CONTRIBUTING.md
├── CHANGELOG.md           lean (Keep-a-Changelog); no TSC ledger
├── CODE_OF_CONDUCT.md
├── SECURITY.md
├── LICENSE
├── VERSION
├── cn.json
├── install.sh
├── docs/
│   ├── README.md          intent-nav
│   ├── THESIS.md
│   ├── quickstart/
│   ├── concepts/          + roles.md (from root ROLES.md), doctrine/, lineage/
│   ├── guides/            + operator.md (from root OPERATOR.md), sustainability.md
│   ├── reference/         canonical CURRENT specs only (cli/ protocol/ runtime/
│   │                      packages/ ctb/{LANGUAGE-SPEC.md} governance/ schemas/)
│   ├── architecture/
│   ├── development/       CDD.md, ENGINEERING-LEVELS.md, rules/, checklists/,
│   │                      issues/, kata/{templates only}   (no plans/, no runs/, no board/)
│   └── papers/            essays only  (no ASSESSMENT, no MODULAR-REFACTOR record)
├── src/
│   ├── go/
│   └── packages/          cnos.core cnos.cdd cnos.cds cnos.cdr cnos.handoff cnos.eng cnos.kata cnos.cdd.kata
├── schemas/               (+ fixtures)
├── scripts/
├── tests/
├── cue.mod/
├── .cdd/                  LIVE control only: CADENCE CDD-VERSION DISPATCH OPERATORS
│                          exceptions.yml + gitignored in-flight cells
│                          (no releases/ iterations/ waves/ proposals/)
├── .cn-sigma/             README.md only  (no logs/)
└── .github/
```

Removed vs today: `docs/development/plans/`, `docs/development/kata/runs/`, `docs/development/board/`, `docs/evidence/rca/` + `AUDIT.md`, `docs/reference/legacy/`, `docs/reference/ctb/*-draft.md` + `*-VISION.md`, `docs/papers/ENGINEERING-LEVEL-ASSESSMENT.md` + `MODULAR-ARCHITECTURE-REFACTOR.md`, `docs/development/cdd/POST-RELEASE-EPOCH-*`, root `RELEASE.md`, `.cdd/releases|iterations|waves|proposals`, `.cn-sigma/logs`. Relocated: root `ROLES.md`/`OPERATOR.md`/`SUSTAINABILITY.md` → `docs/`.

---

## 6. Scope note

**Deferred to the code pass (not filed here):** behavior of `src/go/**`, `src/packages/**` skill logic, `schemas/**` as executable CUE/JSON contracts, and the `docs/gamma/cdd/{X.Y.Z}/` snapshot-path *convention* where it is a live behavioral contract read by package skills (per `round-03` subset-C). I did not evaluate whether that code is correct.

**Filed here regardless of code proximity:** legacy *narrative* anywhere — retirement stories (C7), authoring-intent prose (C8), dated plans/assessments (C2/C3), and the `.cdd/`/`.cn-sigma/` process archives (S1/S2). A record's being version-keyed or "frozen" is not a reason to keep it in the tree; it is the reason git history exists. Where I judged something genuinely current-state doctrine I kept it and said so: `docs/README.md` nav + archival law, `docs/papers/README.md`, `.cn-sigma/README.md` namespace doc, and the live `.cdd/` control files all stay.

**One dependency flag for whoever remediates S1:** confirm no live tool (CI gate, `cn` command, aggregator like `.cdd/iterations/INDEX.md`) reads the `.cdd/releases/**` files at runtime before deleting — if a gate does, that coupling is itself an L7 finding (a verifier that depends on an in-tree history store), and the fix is to break the coupling, not to keep the archive.
