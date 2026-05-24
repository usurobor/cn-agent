# α closeout — cycle/425

**Issue:** [cnos#425](https://github.com/usurobor/cnos/issues/425) — Capture remote-runner-delegation primitive + first use (v3.82.0 tag retarget via GH Actions).

**Mode:** docs + skill patch + workflow + receipt; γ+α+β-collapsed-on-δ.

## Matter produced

Five-deliverable bundle, all in commit `334f1ca6`:

1. **`docs/gamma/essays/BOX-AND-THE-RUNNER.md`** — new file, 288 lines, DRAFT v0.1.0. Foundational design essay capturing the remote-runner-delegation primitive as cnos doctrine. 10 H2 sections (subtitle + Point, Progress note, What I realized, What this means for cnos, What not to celebrate, Memory to keep, The 6-field receipt convention, First use, References). The six prose sections from the operator's seed are preserved verbatim with light essay-class scaffolding; the three synthesis sections (6-field receipt convention template, first-use anchor, References with 14 entries) are new but ground their claims in the seed and already-landed cnos doctrine.

2. **`docs/gamma/essays/README.md`** — 2-line edit. Document Map row appended after `CELL-OF-CELLS.md`; Reading Order item 7 added (last position; the essay closes the gamma essays sequence as the operational discipline that complements the prior synthesis essay).

3. **`src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`** — appended §8 "Remote-runner delegation — δ-class effect surface" (53 lines added; file grew from 319 to 372 lines, ≥ 95% of pre-cycle). Six sub-sections: §8.1 Required output (the 6-field table); §8.2 Authoring-order rule (doctrine-before-first-use); §8.3 One-shot workflows self-delete; §8.4 Composition with the existing outward membrane; §8.5 What this is not; §8.6 First use. Cites `BOX-AND-THE-RUNNER.md` as canonical doctrine in the opening line.

4. **`.github/workflows/repoint-3.82.0.yml`** — new file, 60 lines. One-shot, push-triggered workflow that force-moves the `3.82.0` tag to `fd1d654e` and self-deletes via final-step `git rm` + commit + push. Header comments (lines 1–18) cite the essay, skill section, receipt, and issue. Triggers on `push: branches: [main], paths: ['.github/workflows/repoint-3.82.0.yml']` so it fires only when the cycle/425 merge lands (by which point the doctrine is on main per the authoring-order rule).

5. **`.cdd/unreleased/425/remote-runner-receipt-3.82.0.md`** — new file, 161 lines. First instantiated 6-field receipt under the remote-runner doctrine. All six fields filled (§1 Who asked: operator + dispatch directives; §2 What runs: 4 steps with explicit git commands; §3 Where: `ubuntu-latest`; §4 Authority: `GITHUB_TOKEN contents: write`; §5 Evidence: post-run-fillable with shape declared; §6 Who may accept: operator with release-body match criterion). Plus expected-effect (5 steps), failure-modes (6 rows with mitigations), acceptance-criteria (6 numbered + partial/rejection branches), and V/δ composition note.

## Authoring discipline

The operator's seed content from the dispatch brief was the substantive source. α held to the "preserve where possible, light essay-flow polish only" directive:

- **§Point** — verbatim seed; one scaffolding paragraph added that restates the boundary-error claim in terms cnos already names.
- **§Progress note** — verbatim seed throughout the first-person narrative.
- **§What I realized** — verbatim for the load-bearing claim; one code block added listing six effect surfaces concretely.
- **§What this means for cnos** — verbatim for the primitive name, the rule, the 6-field list, and the "sandbox guarding one limb" framing; two paragraphs added connecting to existing δ-as-role doctrine and the release-effector precedent.
- **§What not to celebrate** — verbatim for the trick-vs-protocol framing; two paragraphs added on trust delegation.
- **§Memory to keep** — verbatim for the three-line closing and the closing principle; one restatement paragraph added.

Essay-class polish (paragraph transitions, the effect-surfaces code block, cross-references to `delta/SKILL.md` and `release-effector/SKILL.md`) was additive — no substance was displaced or rewritten.

The three new synthesis sections (6-field receipt convention, first use, References) are required additions per the issue body and ground their claims in the seed material plus already-landed cnos files.

## Mechanical AC verification

All 11 ACs PASS per `self-coherence.md`:

| AC | Verdict | Evidence |
|----|---------|----------|
| AC1 | PASS | 288 lines (200–400) |
| AC2 | PASS | 10 H2 sections (≥ 9) |
| AC3 | PASS | 17 critical-content hits (≥ 5) |
| AC4 | PASS | README Document Map row + Reading Order item 7 |
| AC5 | PASS | 12 skill-patch grep hits + 372 lines (≥ 303 = 95% × 319) |
| AC6 | PASS | workflow has named triggers + permissions + 3 substantive job steps + self-delete |
| AC7 | PASS | receipt has 6 numbered field sections, all populated |
| AC8 | PASS | CDD.md 0 lines diff; cds/cdr/handoff packages 0 lines diff; only delta/SKILL.md in cdd/skills/ |
| AC9 | PASS | no mojibake (essay uses no Greek letters; the doctrine is about workflows/runners; vacuously satisfied + manual mojibake check returned 0) |
| AC10 | PASS | exactly 2 essay files: BOX-AND-THE-RUNNER.md (added) + README.md (modified) |
| AC11 | PASS | schemas/, src/go/, cnos.core/, cnos.eng/, cnos.kata/, cnos.cdd.kata/, scripts/ all 0 lines diff |

## Refusal conditions held

All hard rules from the dispatch brief held:

- **Doctrine before workflow runs.** All 5 deliverables in one commit (`334f1ca6`); workflow trigger is path-filtered on its own path so it fires only post-merge to main, by which point the doctrine is on main.
- **Receipt explicit.** All 6 fields filled; evidence is post-run-fillable with shape declared.
- **Workflow one-shot.** Final step is `git rm` + `git commit` + `git push origin main`.
- **No CCNF kernel / CDS / CDR / handoff changes.** Verified by mechanical diff.
- **Tag target is fd1d654e.** Workflow hardcodes the SHA; receipt §2 and §6 anchor on it; not current main HEAD.
- **No schemas / runtime / scripts.** Verified by mechanical diff.
- **Greek letters used natively where applicable.** The essay does not require Greek letters (its subject is workflows/runners, not the cell algorithm); the skill section uses Greek (δ, α, β, γ) natively as inherited from the surrounding file.

## Commit

- **α-425** (`334f1ca6`): "α-425: BOX-AND-THE-RUNNER essay + delta SKILL §8 + repoint-3.82.0 workflow + receipt" — all 5 deliverables in a single commit because the doctrine-before-first-use precedence rule requires the bundle to land atomically (the workflow file cannot exist on main without the receipt and doctrine being on main too).

## Handoff

Self-coherence complete; β-review (R1) APPROVED. γ to file gamma-scaffold, gamma-closeout, cdd-iteration (substantive — see cdd-protocol-gap finding), and INDEX.md row, then push branch.

Filed by α@cnos (γ+α+β-collapsed-on-δ) on 2026-05-23.
