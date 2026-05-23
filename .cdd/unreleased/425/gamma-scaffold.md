# γ scaffold — cycle/425

**Issue:** [cnos#425](https://github.com/usurobor/cnos/issues/425) — Capture remote-runner-delegation primitive as cnos doctrine + execute its first use (v3.82.0 tag retarget workflow). Single cycle bundling capture + first exercise.

**Mode:** docs + skill patch + workflow + receipt; γ+α+β-collapsed-on-δ (per `ROLES.md §4`-precedent for design-essay + receipt-discipline cycles with mechanical AC oracles; α=β permitted because the primary product is markdown + YAML with grep-checked and shape-checked invariants, paralleling the cycle-414 (DECREASING-INCOHERENCE.md) and cycle-424 (CELL-OF-CELLS.md) precedents).

**Branch:** `cycle/425` from `origin/main` (`ecf63114` — the cycle-424 post-CELL-OF-CELLS merge).

## Surface (5-deliverable bundle + close-out artifacts)

Five-deliverable bundle authored under the doctrine-before-first-use precedence rule:

- **D1:** `docs/gamma/essays/BOX-AND-THE-RUNNER.md` (new; 288 lines; DRAFT v0.1.0). Authored against the operator's seed content from the dispatch prompt: 10 H2 sections (subtitle + Point, Progress note, What I realized, What this means for cnos, What not to celebrate, Memory to keep, The 6-field receipt convention, First use, References). Seed prose preserved verbatim where applicable; new synthesis sections (6-field receipt convention with template; first-use anchor; References) added per issue body. Frontmatter has 12 `related:` entries spanning gamma essays + cnos.cdd kernel + handoff + release-effector + operator + ROLES + the workflow + the receipt.

- **D2:** `docs/gamma/essays/README.md` (edit; 2 lines added). One Document Map row appended after `CELL-OF-CELLS.md`; one Reading Order item added at position 7 (last). Reading Order choice: **position 7 (last)** — the essay is the operational discipline (remote-runner receipt convention) that complements the prior synthesis essay (CELL-OF-CELLS.md, position 6). Placing it after the synthesis preserves the synthesis-then-discipline reading flow; placing it earlier would require forward-references to the system-layer integration thesis the essay extends. No alternative considered seriously — the position-7 placement is structurally indicated.

- **D3:** `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` (edit; §8 added; 53 lines, file grew 319 → 372, ≥ 95% of pre-cycle). New §8 "Remote-runner delegation — δ-class effect surface" with 6 sub-sections (§8.1 Required output; §8.2 Authoring-order rule; §8.3 One-shot workflows self-delete; §8.4 Composition with the existing outward membrane; §8.5 What this is not; §8.6 First use). Cites `BOX-AND-THE-RUNNER.md` as canonical doctrine. Cites `release-effector/SKILL.md` and `operator/SKILL.md` as composing skills.

- **D4:** `.github/workflows/repoint-3.82.0.yml` (new; 60 lines). One-shot, push-triggered workflow. Triggers on `push: branches: [main], paths: ['.github/workflows/repoint-3.82.0.yml']`. Three substantive job steps (checkout / git identity / tag force-move) plus self-delete. Header comments (lines 1–18) cite doctrine, skill section, receipt, and issue.

- **D5:** `.cdd/unreleased/425/remote-runner-receipt-3.82.0.md` (new; 161 lines). First instantiated 6-field receipt under the new doctrine. All 6 fields populated (§5 Evidence is post-run-fillable per hard rule 2; the shape is named in advance). Plus expected-effect (5 steps), failure-modes (6 rows + mitigations), acceptance-criteria (6 numbered + partial/rejection branches), V/δ composition note.

Close-out artifacts at `.cdd/unreleased/425/`:

- `gamma-scaffold.md` (this file)
- `self-coherence.md` (α-side AC verification, all 11)
- `beta-review.md` (β-side R1 review verdict; 11 binding findings, all PASS)
- `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md` (per-role closeouts)
- `cdd-iteration.md` (**substantive** — records the `cdd-protocol-gap` finding "boundary model incomplete — remote-runner delegation needed naming" with disposition `patch-landed-in-cycle` since this cycle IS the patch)
- INDEX.md row appended at `.cdd/iterations/INDEX.md`.

## Implementation contract (pinned by δ at dispatch; verified at this scaffold)

| Axis | Pinned value | Conforms? |
|---|---|---|
| Language | Markdown + YAML | yes |
| CLI integration target | None new; uses GitHub Actions | yes |
| Package scoping | New essay + README pointer + delta/SKILL.md section + new workflow + new receipt; close-out artifacts | yes (exact 5-deliverable + standard close-out bundle) |
| Authority | `GITHUB_TOKEN` with `contents: write` for the workflow runner | yes (encoded in workflow `permissions:` + receipt §4) |
| Backward compat | Existing essays unchanged; delta/SKILL.md preserved at ≥ 95% original line count (372 ≥ 303) | yes (mechanical diff: 0 lines changed on other essays; 5 other prior gamma essays + alpha/essays/ + docs/essays/ all byte-identical to origin/main) |
| Doctrine-before-first-use precedence | All 5 deliverables in commit-set; workflow trigger path-filtered on own path so it fires only post-merge | yes (commit `334f1ca6` atomic; trigger pinned) |

## AC oracle approach (issue body verbatim)

| AC | Oracle (mechanical) | Surface |
|----|---------------------|---------|
| AC1 | file exists, line count 200–400 | `wc -l docs/gamma/essays/BOX-AND-THE-RUNNER.md` → 288 |
| AC2 | `grep -c "^## "` ≥ 9 | 10 |
| AC3 | `grep -c "remote-runner delegation\|effect surface\|6-field receipt\|agent boundary"` ≥ 5 | 17 |
| AC4 | Document Map row + Reading Order entry | `grep -c BOX-AND-THE-RUNNER docs/gamma/essays/README.md` → 2 |
| AC5 | `grep -c "remote-runner\|BOX-AND-THE-RUNNER" delta/SKILL.md` ≥ 2; file line count ≥ 95% pre-cycle (319 × 0.95 = 303) | 12; 372 |
| AC6 | workflow has `name`, `on:` with branches + paths + permissions; 3 substantive steps + self-delete | verified by read |
| AC7 | receipt file exists; contains all 6 named field rows | `grep -c "^## [1-6]\." receipt.md` → 6 |
| AC8 | CDD.md diff 0; cds/cdr/handoff diff 0; only delta/SKILL.md in cdd/skills/ | 0; 0; only delta/SKILL.md |
| AC9 | no mojibake (Greek used natively where applicable) | essay does not require Greek (subject is workflows/runners); skill section inherits Greek discipline from surrounding file |
| AC10 | only BOX-AND-THE-RUNNER + README in essays | `git diff --name-only origin/main -- docs/gamma/essays/` → exactly 2 |
| AC11 | schemas/runtime/scripts diff 0 | 0 lines across schemas/, src/go/, cnos.core/, cnos.eng/, cnos.kata/, cnos.cdd.kata/, scripts/ |

## Branch + commit shape

- **α-425** (`334f1ca6`): "α-425: BOX-AND-THE-RUNNER essay + delta SKILL §8 + repoint-3.82.0 workflow + receipt" — all 5 deliverables in a single atomic commit (doctrine-before-first-use precedence rule requires the bundle to land atomically; workflow trigger fires only post-merge by which point doctrine is on main).
- **β-425** (`50de609a`): "β-425: R1 review APPROVED (11/11 PASS) + role closeouts (α, β)" — beta-review.md (11 binding findings, all PASS) + alpha-closeout.md + beta-closeout.md.
- **γ-425** (forthcoming this commit): "γ-425: close-outs (γ + scaffold + self-coherence) + substantive cdd-iteration + INDEX row" — gamma-scaffold.md (this file) + self-coherence.md + gamma-closeout.md + cdd-iteration.md (substantive) + INDEX.md row.

Push to `origin/cycle/425`; do NOT merge to main (operator's call). Merge instruction reported in γ closeout.

## Critical refusal conditions surfaced during authoring

- **Doctrine-before-first-use precedence held.** All 5 deliverables land in commit `334f1ca6`; workflow's push trigger fires only on changes to its own path on main, so it cannot run before the cycle merges. There is no window in which the workflow exists on main without the doctrine.

- **Verbatim preservation of seed.** The dispatch brief states: "preserve where possible, light essay-flow polish only." α held to this in the six core prose sections (Point, Progress note, What I realized, What this means for cnos, What not to celebrate, Memory to keep). The seed's load-bearing claims — the primitive name "remote-runner delegation", the rule "Any artifact that can make another system execute is an effect surface", the 6-field list, the "sandbox guarding one limb" framing, the trick-vs-protocol distinction, the closing principle "the agent boundary is not the place where the model sits; it is the full path from intention to effect" — all appear verbatim. Essay-class scaffolding (paragraph transitions, the effect-surfaces code block, cross-references to existing δ/release-effector doctrine) is additive, not substance-displacing.

- **Receipt evidence is post-run-fillable, not invented.** Hard rule 2 authorizes §5 Evidence to be a post-run-fillable placeholder. α populated §5 with a 5-row evidence-shape table naming what evidence will be filled when the run completes (workflow run URL pattern, tag-move proof pattern, release re-publish URL pattern, self-delete commit SHA pattern, workflow-file-gone proof pattern), rather than fabricating evidence. The shape declaration enables downstream consumers to detect a missing receipt as a structural signal.

- **Tag target hardcoded.** The workflow's tag-move step uses `git tag -f 3.82.0 fd1d654e` — the literal SHA, not `main`, not `HEAD`. The receipt's §6 acceptance criterion is anchored on the release body matching `.cdd/releases/3.82.0/RELEASE.md` (which is the canonical body that lives at `fd1d654e`). The current main HEAD (`ecf63114` at branch creation; now `50de609a` on this cycle branch) has CELL-OF-CELLS content that landed post-tag and is NOT the correct retarget point.

- **No kernel / CDS / CDR / handoff changes.** Verified by mechanical diff. Only `delta/SKILL.md` touched in `cnos.cdd/skills/`.

- **No schemas / runtime / scripts.** Verified by mechanical diff (0 lines across all 7 excluded paths).

- **Citation closed loop.** Each of the 5 deliverables cites the others. Verified by reading each artifact:
  - Essay → cites skill, workflow, receipt, plus existing doctrine (CELL-OF-CELLS, delta/SKILL.md, release-effector/SKILL.md, operator/SKILL.md, etc.)
  - README → cites essay (Document Map row + Reading Order item)
  - Skill §8 → cites essay as canonical doctrine; cites workflow + receipt in §8.6
  - Workflow → cites essay, skill section, receipt, issue in header comments
  - Receipt → cites essay, skill section, workflow, issue in header
  No artifact stands alone.

Filed by γ@cnos (γ+α+β-collapsed-on-δ) on 2026-05-23.
