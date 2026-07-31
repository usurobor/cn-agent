# β Review — Code Pass, Round 01 (scope + plan the 6 deferred code-coupled items)

**Role:** β (adversarial reviewer / planner). **Branch:** `claude/repo-cleanup-code-pass` (stacked on docs PR #679). **HEAD:** `a66f4cf0`.
**Method:** every coupling claim below is re-derived from a path I opened or `git grep`'d on the live tree — not from the handoff agenda. Where the handoff and the tree disagree, the tree wins and I say so. Edited only this file; no commit; no PR.
**Mandate (unchanged):** the working tree reflects only current state; git history is the archive; dotdirs (`.cdd/`, `.cn-sigma/`, `.github/`) are the intentional history store. Code is now in scope. Standard: cnos write skill + L7.

**Headline:** All of item 1's inbound references are **citations/comments, zero are runtime loads** — no skill/command/schema/CI job breaks on deletion; every repair is a line edit. But three handoff premises do not survive contact with the tree and must go back to the operator: **(a)** two of the six "delete" docs (`LANGUAGE-SPEC-v0.2-draft.md`, `CTB-v4.0.0-VISION.md`) are **current-state governing docs, not historical records** — deleting them removes live spec, not history; **(b)** `RELEASE.md` is a **hard-required, already-current release-gate artifact** — it must not be moved or de-historicized; **(c)** `generate-release-tag-message.sh` is **not coupled** to CHANGELOG/RELEASE (reads `.cdd/` cycle artifacts). The single highest-leverage move is **item 6** (strip the design-cycle apparatus at *promotion*), which stops every future promoted doc from re-injecting exactly what rounds 8–9 removed. The whole pass is gated on **one doctrine ratification** (item 5), which is the license to delete frozen records at all.

---

## Item 1 — Delete 6 code-coupled docs + repair references

**Verified inbound refs (all `git grep`'d; classified load-bearing vs citation):** none of the 6 is loaded at runtime by any skill/command/schema. Every hit is prose, a code comment, a reading-list bullet, or a relative doc link. Deleting the doc therefore never breaks execution — it only orphans a citation, and each repair is a one-line edit.

| Doc | Is it a historical record? | Inbound refs (exact) | Load-bearing? | Repair |
|-----|---------------------------|----------------------|---------------|--------|
| `docs/reference/packages/DESIGN-266-dist-out-of-git.md` | **Yes** — `Issue:#266 / Mode:MCA / Active Skills / Engineering Level` header + Problem/Impact Graph/File Changes/Acceptance Criteria/CDD Trace. Design record. | `.gitignore:37` (comment); `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md:176` (prose "*Derives from: #266 F3 — `DESIGN-266…md` carried one count across 4 sentences…*") | No (comment + skill prose) | `.gitignore:37`: drop the filename, keep "never committed" rationale. `alpha/SKILL.md:176`: keep the #266 F3 lesson, drop the filename. |
| `docs/papers/ENGINEERING-LEVEL-ASSESSMENT.md` | **Yes** — dated point-in-time assessment (`Period: Feb 2–Mar 25 2026`, `Commits: 1,434`). Both referrers explicitly frame it as *"a historical/descriptive assessment of observed work in a period."* | `docs/development/ENGINEERING-LEVELS.md:11` (Related list), `:285` ("## 10. Relationship to the assessment essay"); `docs/papers/README.md:48` (index link); `src/packages/cnos.eng/skills/eng/README.md:9` (Related list), `:277` ("## Relationship to the level assessment essay") | No (reading lists + prose sections) | Drop the essay bullet from both `Related:` lists + `papers/README.md:48`; delete/rewrite the two "Relationship to the … essay" sections (they contrast the rubric with the essay — after deletion, keep the rubric-is-normative statement, drop the essay reference). |
| `docs/reference/ctb/CTB-v4.0.0-VISION.md` | **QUESTIONABLE — likely KEEP.** `Document type: Vision + Design Definition`, `Status: Draft for alignment (guide next 2–4 iterations)`. Cited as the **governing strategy authority** by 4 live CTB docs. Forward-looking roadmap, not a dated record. | `docs/reference/ctb/README.md:11` (table: "governs strategy"); `LANGUAGE-SPEC.md:19,508` ("Vision governs strategy/motivation/roadmap"); `SEMANTICS-NOTES.md:7,610` ("Companion to"; "the Vision governs"); `LANGUAGE-SPEC-v0.2-draft.md:21`; `docs/reference/runtime/ORCHESTRATORS.md:325,339,371` (cites §3.1, "matches CTB-v4.0.0-VISION exactly"); `src/packages/cnos.core/skills/agent/emoji-language/SKILL.md:70` ("See: [CTB v4.0.0 Vision]") | No (doc cross-refs + one skill "See:" link) | **Do not delete without operator ruling.** If kept: no action. If deleted: cascade repair — remove the README table row, rewrite the "Vision governs strategy" pointers in LANGUAGE-SPEC/SEMANTICS-NOTES (they would lose their strategy authority), rewrite ORCHESTRATORS §325/339/371, drop the emoji-language "See:" link. This is a lattice, not a line. |
| `docs/reference/ctb/LANGUAGE-SPEC-v0.2-draft.md` | **NO — current-state.** `Version: 0.2 (draft-normative)`, `Status: Draft-normative migration target`. README:18 says it **"governs the migration target."** This is an *active* draft spec, not history. | `docs/papers/ACTIVATION-NOT-DEPLOYMENT.md:178` (inside a fenced **illustrative JSON receipt** — not a manifest); `CTB-v4.0.0-VISION.md:271,272,464` (§1/§6/§1.4 pointers); `docs/reference/ctb/README.md:13,18` ("governs the migration target"); `schemas/README.md:184` ("v0.2 draft — out of scope until promoted") | No | **Recommend KEEP — deleting it violates the mandate** (removes current draft-normative spec). If operator insists on delete, it also cascades into the Vision doc's §-refs and the README governance rows. |
| `docs/development/cdd/DISPATCH-FAILURE-EVIDENCE.md` | **Yes** — "Evidence log for #295. Each triad run records failures…" Audit/evidence log. | `src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md:79` (prose "…(see `DISPATCH-FAILURE-EVIDENCE.md`, cycle #84 failure 3)") | No (skill prose) | `release-effector/SKILL.md:79`: keep the "manual tagging forbidden" rule, drop the parenthetical file/cycle citation (or restate the lesson inline). |
| `docs/development/design/cn-repo-install-MOCKS.md` | **Yes** — `Status: Design surface (pre-implementation)`; `docs/guides/INSTALL-CDS.md:253` confirms the installer **"was built against"** it (shipped). Design record for landed work. | `.github/workflows/build.yml:72` (comment); `.github/workflows/install-wake-golden.yml:692` (comment); `src/packages/cnos.core/commands/install-wake/cn-install-wake:40,1070,1082` (comments); `docs/guides/INSTALL-CDS.md:253` (citation); self-ref `:260` | **No — verified the workflows do not read the file, only name it in comments** (CI cannot go red). | Rewrite the 5 comments (`build.yml:72`, `install-wake-golden.yml:692`, `cn-install-wake:40/1070/1082`) to keep the `cnos#61x` / Mock-letter provenance, drop the file path. Remove/rewrite `INSTALL-CDS.md:253`. Then handle the dir (below). |

**Directory lifecycle:** `docs/development/design/` = `README.md` + `cn-repo-install-MOCKS.md` only. After the delete, decide: remove the dir + README, or keep the README as a survey stub. Mechanical either way.

**Classification:**
- **SAFE-MECHANICAL** — delete + repair the **4 genuine records** (DESIGN-266, ENGINEERING-LEVEL-ASSESSMENT, DISPATCH-FAILURE-EVIDENCE, cn-repo-install-MOCKS) + the `docs/development/design/` dir. No runtime coupling; all repairs are line edits; comments in CI files cannot turn CI red.
- **NEEDS-OPERATOR-DECISION** — the **CTB pair** (CTB-v4.0.0-VISION, LANGUAGE-SPEC-v0.2-draft). β recommendation: **KEEP both** — v0.2-draft is an active draft-normative migration target and the Vision is the cited strategy authority; neither is the dated/evidence-record class. If deleted anyway, it is a **multi-doc cascade repair**, not a line edit.

---

## Item 2 — `RELEASE.md` / `CHANGELOG.md` release-gate rework

**Full coupling map (each read/require re-derived):**

| Consumer | What it reads/requires | Hard or soft |
|----------|------------------------|--------------|
| `scripts/validate-release-gate.sh:55` | `RELEASE.md` **must exist at repo root** (release mode) before a tag may be cut | **HARD fail** (exit 1) |
| `.github/workflows/release.yml:165,180` | `RELEASE.md` used as GitHub-release `body_path` | Hard (release body) |
| `src/packages/cnos.cdd/commands/cdd-verify/ledger.go:96–101` | `CHANGELOG.md` must contain a `## $VERSION` entry | **HARD fail** in ledger mode (`checkFail`) |
| `src/packages/cnos.cdd/commands/cdd-verify/ledger.go:138–141` | `RELEASE.md` presence | Soft (`checkWarn`, "may have been overwritten") |
| `scripts/release.sh:88–95` | `CHANGELOG.md` has `## $VERSION` | Soft (interactive `y/N` warning, abortable) |
| `src/packages/cnos.cdd/commands/cdd-status/cn-cdd-status:159,291` | `RELEASE.md` written+committed (status step 11) | Status display |
| `scripts/test-validate-release-gate.sh`, `scripts/test-release-tag-integration.sh`, `.../test-cn-cdd-status.sh` | assert RELEASE.md present/absent behavior | **Test gate** |
| `src/packages/cnos.cds/skills/cds/CDS.md` (many) | γ authors `RELEASE.md` per release; CHANGELOG "Release Coherence Ledger" row per release | Doctrine |
| `scripts/generate-release-tag-message.sh` | **NOT coupled** — reads `.cdd/` cycle artifacts (`manifest.md`, `beta-review.md`, `beta-closeout.md`), zero CHANGELOG/RELEASE refs | — (handoff over-listed it) |

**What the tree shows about "de-historicizing" these files:**
- **`RELEASE.md` is already current-state and load-bearing.** It holds **only the current release** (v3.82.0) — it is overwritten each release, not accumulated. It carries no historical residue to strip. **Moving or deleting it fails the hard gate + 3 test scripts.** Correct action: **leave it exactly as is.** The handoff's inclusion of RELEASE.md in a "de-historicize" bucket is a category error.
- **`CHANGELOG.md` is an accumulating ledger (184 KB).** That accumulation *is* its conventional current-state function (Keep-a-Changelog + the CDS "Release Coherence Ledger"), and it is hard-coupled to `ledger.go`. "De-historicizing" it — trimming past entries, or dropping it for git-history-only — is a **release-record contract change**, not a mechanical edit, and would break `cn cdd verify` ledger mode unless the version-entry check is also changed. The dead cite at `CHANGELOG.md:1336` (`docs/beta/evidence/rca/`, dir no longer exists) sits **inside a historical changelog entry**; rewriting a point-in-time ledger entry is itself an anti-pattern — so even this "one dead link" is part of the contract question, not a free fix.

**Classification: NEEDS-OPERATOR-DECISION (release-record contract).** No SAFE-MECHANICAL sub-part exists: RELEASE.md is leave-alone, and every CHANGELOG edit is either contract-level or a questionable history-rewrite. β recommendation: **keep both under the current contract** (RELEASE.md unchanged; CHANGELOG remains the canonical human-facing ledger, gate-coupled), and treat "history is the archive" as targeting *process/design-cycle narrative*, not a conventional changelog. **RISKY-CI** if the operator elects any CHANGELOG/gate change — see risk map.

---

## Item 3 — `docs/development/board/`

**Verified coupling:** dir = `README.md` + `board-data.json` (62 KB) + `index.html` (382 KB), both data files **regenerated today (Jul 31)**. README stamps it *"Generated file — do not edit by hand."*
- `.github/workflows/board-map.yml` — the `board-map` Action runs `./cn issues map --repo … --out docs/development/board` (`:57`) and **auto-commits** the regenerated dir on issue changes (`:73–74`, bot `board-map-bot`).
- `src/packages/cnos.issues/commands/issues-map/issuesmap.go:53` — `const defaultOut = "docs/development/board"`.
- `cn-install-wake:1257,1302` — comments noting `board-map-bot` heartbeat commits the worker rebases onto.

**Assessment:** this is **not stale history** — it is a **live materialized view** of the open-issue board, auto-regenerated by CI. Under the mandate it is current-state (a cache of GitHub live data), so **no de-historicization is strictly required.** The real question is a **strategy choice**: keep the committed-and-regenerated snapshot (status quo) vs drop the committed artifact and link the live GitHub board / githack preview (removes the CI commit loop).

**Classification: NEEDS-OPERATOR-DECISION (board strategy)**, β default **KEEP / no action** (it is current, not history). If "drop": touches `issuesmap.go:53` + delete/rewrite `board-map.yml` + update `cn-install-wake` comments — mechanical, but exercises `go test` for `issues-map` (board-map.yml itself fires on issue events, not PR CI, so it will not gate the code-pass PR).

---

## Item 4 — Root `ROLES.md` / `OPERATOR.md` / `SUSTAINABILITY.md` → `docs/`

**Critical disambiguation (the handoff conflates two different files):** almost every `OPERATOR.md` hit in code/golden/kata is **`<hub>/spec/OPERATOR.md`** — a per-hub identity spec, **not** root `OPERATOR.md`. `src/go/internal/activate/activate.go:291` (`scanOperator`), `activate_test.go`, the two `.golden.yml` fixtures, and `cnos.kata/R5-activate/run.sh` all reference `spec/OPERATOR.md`. **Root `OPERATOR.md` has ZERO code/schema/CI/golden/kata referrers** (only CHANGELOG history mentions it). So relocating root OPERATOR.md is **coupling-free** except its own header.

**Full referrer list, per root file:**
- **`OPERATOR.md` (root):** no live referrers. `activate.go` + golden fixtures + kata all bind `spec/OPERATOR.md` (different file). Move = update its own header only.
- **`SUSTAINABILITY.md` (root):** `README.md:195` (`[Sustainability](SUSTAINABILITY.md)` — relative link) + its own `**Canonical-Path:** SUSTAINABILITY.md` stamp (line 6). Move = update the README link + the Canonical-Path stamp.
- **`ROLES.md` (root):** wide but **entirely text**, none CI-resolved —
  - `schemas/cdd/README.md:141`, `schemas/cdr/README.md:6,50,72,118,136`, `schemas/cds/README.md:62,113` — relative links `../../ROLES.md[#anchor]` (would rot; depth changes on move).
  - `schemas/cdr/receipt.cue:7,33,52,60`, `schemas/cds/receipt.cue:14`, `schemas/cdr/fixtures/valid-cdr-receipt.yaml:3` — **CUE/YAML comments** ("anchored in `ROLES.md §4a.3`"). `cue vet` ignores comments → **no golden/schema break.**
  - skill prose: `CDD.md:95,107`, `CELL-KINDS.md:204,222`, `COHERENCE-CELL-NORMAL-FORM.md:391`.
  - papers frontmatter `related: ROLES.md` + prose: `BOX-AND-THE-RUNNER.md`, `CCNF-AND-TYPED-TRUST.md`, `CELL-OF-CELLS.md`, `DECREASING-INCOHERENCE.md`.
  - `ROLES.md:186` itself references `<hub>/spec/OPERATOR.md`.

**Is a schema/golden path change in-scope-safe?** **Yes — verified CI-safe.** The only CI validator that resolves paths is `scripts/ci/validate-skill-frontmatter.sh` (build.yml:334,337), and it validates **only `SKILL.md` under `src/packages/` and their `calls:` frontmatter** — it never resolves `related:` in papers, doc links, or `.cue` comments. There is **no markdown link-checker** in any of the 6 workflows. The golden fixtures reference `spec/OPERATOR.md` (hub), untouched by root moves. So **relocation cannot turn CI red** — the risk is reader-facing link rot only.

**Classification: NEEDS-OPERATOR-DECISION (whether to relocate at all)** — relocation is an **organizational/newcomer-clarity choice, not a de-historicization necessity**; `ROLES.md` is current-state doctrine (single source of truth for the role ladder, cited as the normative anchor for the CDR/CDS receipt schemas). If the operator says relocate: **SAFE-MECHANICAL (CI-safe, no golden/schema break)** but a **wide edit surface** for ROLES.md (schema READMEs' relative-link depth, 6 CUE/YAML comments, 3 skill bodies, 4 papers' frontmatter+prose, plus `Canonical-Path` stamps in the moved files). OPERATOR.md and SUSTAINABILITY.md moves are trivial.

---

## Item 5 — Doctrine reconciliation: `DOCUMENTATION-SYSTEM.md §5 "Frozen history"`

**Verified conflict.** `docs/reference/governance/DOCUMENTATION-SYSTEM.md` §5, subsection **"Frozen history"**, paragraph 2:

> *"A frozen record (a dated design decision, a cycle log, a completed plan) is left in place under its dated or `Superseded` header. Its stale internal paths are not corrected — the header already marks it as history. Only the live reader surface is kept current."*

This **directly encodes the rule the whole mandate reversed** — "leave frozen records in place, don't correct stale paths." It is the standing doctrine that *forbids exactly what item 1 does* (deleting dated design/evidence records because git history is their archive). Until it is rewritten, every future round re-litigates the frozen-bucket ruling, and item 1's deletions technically violate live doctrine.

**Keep unchanged:** paragraph 1 ("Released snapshots are not kept as folders on `main` — git history is the archive… A document on `main` is always the current one") — this is already mandate-consistent. And the **"Supersession"** subsection's `Supersedes:` / `Superseded by:` lineage-stamp convention (γ ruled it a governed convention → KEEP).

**Drafted replacement for paragraph 2 (operator to confirm):**

> *A completed, dated, or superseded record — a design decision, a cycle log, a finished plan — is **not kept on the reader surface**. Git history is its archive: to read it, check out the release or commit under which it was written. The working tree on `main` reflects only current state. Records that must persist as intentional history live in the dotdirs (`.cdd/`, `.cn-sigma/`, `.github/`), which are exempt from the current-state rule. A document that remains on the reader surface is, by that fact, current — its paths and contents are kept correct, not frozen.*

**Classification: NEEDS-OPERATOR-DECISION (doctrine/design).** This is the **license clause** for the entire code pass — it should be ratified **first**, because it authorizes item 1's deletions and constrains item 5/6 wording. α must not silently edit doctrine.

---

## Item 6 — ROOT-CAUSE fix: strip the design-cycle apparatus at *promotion*

**Verified source of recurrence.** The apparatus rounds 8–9 stripped is **prescribed as the authoring template** in three live sources:
- **`src/packages/cnos.cdd/skills/cdd/design/SKILL.md §3.1 "Output Format"` (lines 279–336)** — the canonical design-doc template: `**Issue:** / **Version:** / **Mode:** MCA|MCI / **Active Skills:** / **Engineering Level:**` header, then `## Problem / ## Impact Graph / ## File Changes / ## Acceptance Criteria / ## Known Debt / ## CDD Trace`. §3.2 (`:339`) *requires* the CDD Trace in the primary branch artifact; §3.3 is a 24-item pre-submission checklist enforcing every apparatus section.
- **`.github/PULL_REQUEST_TEMPLATE.md`** — `## CDD Trace` (`:1`), `## Gap (step 4)`, `## Mode + Active Skills (step 5)`, `## Acceptance Criteria` (`:40`), `## Known Debt` (`:48`).
- **`docs/development/cdd/SELF-COHERENCE-TEMPLATE.md`** — `Mode: MCA/MCI` (`:5`), `## Acceptance Criteria Check` (`:20`), `## Known Debt` (`:33`).
- (`docs/development/cdd/PLAN-TEMPLATE.md` is **clean** — Goal/Sequencing/Steps/Checkpoints/Dependencies/Risks/Exit, no apparatus stamps. Handoff over-listed it; **exclude**.)

**The missing rule = the recurrence engine.** `design/SKILL.md` never distinguishes the **working-cell artifact** (which legitimately carries Mode/AC/CDD-Trace during the design cycle — and belongs in a `.cdd/` cycle dir, the PR, or the issue) from the **promoted reader-surface doc** (a canonical `docs/reference/**` | `docs/architecture/**` spec, which must be current-state only). There is **no promotion step that strips the apparatus** — so every design doc authored per §3.1 and landed under `docs/reference/` re-injects precisely the `## 0. Coherence Contract` / `## N. Acceptance Criteria` / `## N. File Changes` / `## N. CDD Trace` class that R8/R9 spent two rounds removing (THREAD-EVENT-MODEL, EXTENSION-REGISTRY, RUNTIME-EXTENSIONS, COGNITIVE-SUBSTRATE were all born with it). This is the **MCA that stops the loop.**

**Proposed fix (the promotion contract):**
1. **Keep the apparatus in the working cell.** The `## Problem / Mode / Active Skills / Engineering Level / Impact Graph / File Changes / Acceptance Criteria / CDD Trace` set is legitimate *cycle scaffolding* — it stays in the `.cdd/` cycle artifact / PR body / issue where the design work happens (dotdir + `.github`, both exempt).
2. **Strip at promotion.** Add a **"Promotion" rule** to `design/SKILL.md` (new subsection after §3.1): *when a design artifact becomes a canonical reader-surface document (`docs/reference/**`, `docs/architecture/**`), strip the cycle apparatus to current-state — retain a governed header (`Version:` / `Status:` / `Doc-Class:` where applicable) + `## Purpose` + the spec body + `## Alternatives` / `## Migration` / `## Non-goals`; fold genuine present-day limits into `## Limitations`; drop `Mode` / `Active Skills` / `Engineering Level` / `Impact Graph` / `File Changes` / `Acceptance Criteria` / `CDD Trace` (they live in the cycle record, which git history preserves).* Annotate §3.1's "Output Format" as **"cycle-artifact form (working cell) — not the promoted reader-surface form."**
3. **Mirror the split in the templates.** Add one line to the top of the PR template and SELF-COHERENCE-TEMPLATE noting they are **cycle scaffolding**, and that promoted docs are stripped per the design skill's Promotion rule. (The PR/self-coherence templates themselves stay full — they *are* the working cell.)

**Classification: SPLIT.**
- **SAFE-MECHANICAL (template edits):** annotate §3.1 as working-cell form + add the one-line notes to the PR/SELF-COHERENCE templates. Low risk. Note `design/SKILL.md` is a `SKILL.md` under the frontmatter validator — editing the **body** (not frontmatter) is safe; re-run the validator to confirm.
- **NEEDS-OPERATOR-DECISION (process/doctrine):** the **promotion contract itself** — what strips vs stays, and that the stripped trace's archive is git history — is a change to CDD authoring doctrine. It must be ratified **together with item 5** (they are the same "current-state-only" doctrine, one for records, one for promoted design docs).
- Minor: `design/SKILL.md:257` carries a stale `PLAN-package-system.md` example string — refresh while in-doc (harmless).

---

## Recommended execution order

**Phase 0 — Operator ratifies doctrine (blocks everything; this is the license).**
- **Item 5** — confirm the "Frozen history" rewrite (authorizes deleting frozen records at all).
- **Item 6 promotion contract** — confirm the strip-at-promotion doctrine (same current-state-only ruling, applied to design docs).

**Phase 1 — Root-cause first (stops recurrence before adding more current-state docs).**
- **Item 6 mechanical** — annotate `design/SKILL.md §3.1` + add the promotion rule; one-line notes on PR template + SELF-COHERENCE-TEMPLATE; refresh the stale example string. Re-run `validate-skill-frontmatter.sh`.

**Phase 2 — Safe-mechanical deletions/repairs (all CI-safe, all line edits).**
- **Item 1 (4 records)** — delete DESIGN-266, ENGINEERING-LEVEL-ASSESSMENT, DISPATCH-FAILURE-EVIDENCE, cn-repo-install-MOCKS; repair the citation lines; resolve `docs/development/design/` dir.
- **Item 5 mechanical** — apply the ratified frozen-history text.

**Phase 3 — Execute the surfaced decisions (only what the operator greenlights).**
- **Item 1 CTB pair** — per operator ruling (β: keep both).
- **Item 4** — relocate ROLES/OPERATOR/SUSTAINABILITY if operator opts in (CI-safe, wide edit surface for ROLES).
- **Item 3 board** — per operator ruling (β: keep/no-action).
- **Item 2 release-record contract** — **last, and isolated in its own PR** — the only item that can hard-fail CI.

Rationale for order: item 6 is highest-leverage (prevents the whole class from re-appearing), but its *doctrine half* and item 5 are the license for everything downstream, so they lead. Safe-mechanical work lands next. The one gate-coupled item (2) goes last and alone so a red gate never blocks the rest.

---

## CI / gate risk map (the branch stacks on #679; the code-pass PR must go green)

| Change | Can it turn CI red? | How to validate before merge |
|--------|--------------------|------------------------------|
| **Item 2 — RELEASE.md moved/deleted** | **YES — hard.** `validate-release-gate.sh:55` exits 1; `scripts/test-validate-release-gate.sh` + `test-release-tag-integration.sh` + `test-cn-cdd-status.sh` assert its presence. | **Do not move/delete it.** If any gate change is attempted, run all three test scripts + `validate-release-gate.sh --mode release` locally first. |
| **Item 2 — CHANGELOG contract change** | **YES — hard.** `ledger.go:96–101` `checkFail` if `## $VERSION` missing. | Run `cn cdd verify` ledger mode against the target version; keep the `## $VERSION` entry invariant or change the check in the same PR. |
| **Item 6 — editing `design/SKILL.md` (a live SKILL.md)** | Low — only if frontmatter breaks. | Body-only edits; run `validate-skill-frontmatter.sh --self-test` **and** the full run (build.yml:334,337). |
| **Item 1 — CI-file comment rewrites** (build.yml:72, install-wake-golden.yml:692, cn-install-wake) | **No** — verified comment-only; workflows do not read the MOCKS file. | Confirm the edited lines remain inside `#` comments; run the affected workflow if paranoid. |
| **Item 4 — root relocation** | **No** — no CI validator resolves the affected links (`related:`/doc links/`.cue` comments are not checked; golden fixtures use `spec/OPERATOR.md`, untouched). | Run `validate-skill-frontmatter.sh` to confirm no `calls:` target moved (ROLES/OPERATOR/SUSTAINABILITY are not `calls:` targets). Grep post-move for stale `../../ROLES.md` depths (reader-quality, not CI). |
| **Item 3 — board drop** (if chosen) | Low for PR CI — `board-map.yml` fires on issue events, not PR. `go test` for `issues-map` exercises `issuesmap.go`. | Run `go test ./...` for the `issues-map` package after changing `defaultOut`. |
| **Item 1/6 doc deletions & template edits** | No (no link-checker in CI). | Grep for dangling inbound refs to each deleted doc after the edit. |

**Top 3 CI risks:** (1) touching `RELEASE.md` or the CHANGELOG version-entry contract (item 2) — the only hard-fail coupling, quarantine in its own PR; (2) breaking `design/SKILL.md` frontmatter while editing it (item 6); (3) `issuesmap.go` `go test` if the board is dropped (item 3).

---

## Operator decisions needed before α proceeds (shortlist)

1. **Doctrine license (item 5)** — ratify the "Frozen history" rewrite. *This is the precondition for the entire pass; without it, item 1's deletions violate standing doctrine.*
2. **Promotion contract (item 6)** — ratify strip-the-apparatus-at-promotion (what strips, what stays, archive = git history). Decide **with #1** — same doctrine.
3. **CTB pair (item 1)** — are `CTB-v4.0.0-VISION.md` and `LANGUAGE-SPEC-v0.2-draft.md` historical (delete) or current governing docs (keep)? *β strongly recommends KEEP both — v0.2-draft is an active draft-normative migration target; the Vision is the cited strategy authority; neither is the dated/evidence class the mandate targets.*
4. **Release-record contract (item 2)** — keep CHANGELOG as the canonical gate-coupled ledger (β rec) or move to git-history-only? RELEASE.md stays as-is regardless. *This is the one contract change that risks the release gate.*
5. **Board strategy (item 3)** — keep the auto-regenerated committed snapshot (β rec: it is a live view, not history) or drop + link the live board?
6. **Root relocation (item 4)** — relocate ROLES/OPERATOR/SUSTAINABILITY to `docs/`, or leave at root? *Not a de-historicization necessity; CI-safe either way; ROLES has a wide (but mechanical) repair surface.*

---

## CLP

- **TERMS:** "load-bearing" = removing the target breaks a skill/command/schema/CI job at runtime; "citation" = a prose/comment/link reference repaired by editing the referring line. "Frozen record" per §5 = dated design decision / cycle log / completed plan.
- **POINTER:** every coupling claim is anchored to a path:line I opened or grep'd on HEAD `a66f4cf0`; where the handoff and tree diverge (RELEASE.md category, CTB pair currency, root-vs-hub OPERATOR.md, uncoupled tag-message script, clean PLAN-TEMPLATE) the tree ruling is stated inline.
- **EXIT:** scope + plan + classification + order + CI map + operator-decision shortlist delivered for all 6 items. The next work is the operator's Phase-0 doctrine ratification (items 5 + 6-contract); α cannot safely start deletions until the "Frozen history" clause is reversed, since it currently forbids them.
