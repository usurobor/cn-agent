# γ Scaffold — cnos#639

**Issue:** [#639](https://github.com/usurobor/cnos/issues/639) — cds/doctrine: clarify the dispatch `run_class` taxonomy (first-pass / repair / recovery / continuation / scope-continuation).
**Mode:** design-and-build (doctrine/design only — the issue's own header states this explicitly).
**cell_kind:** `doctrine` (pinned — see Friction note 1; this is a pure prose/doctrine reconciliation, no mechanism to build, unlike #640's sibling shape).
**Branch:** `cycle/639` (already created from `origin/main` at `7313caea3b6b9901cc02e797be9c929393b971ee`; base verified — `git merge-base cycle/639 origin/main` equals that SHA exactly, no drift).
**Wake run id:** `28998099166` (cds-dispatch, protocol `cds`).
**Family:** #583/#584 (mechanism/cognition doctrine), #593 (reconciler/recovery — a first_pass shape with a stale claim marker, not a new run_class), #630 (`resumed_from_matter`), #626 (where the gap surfaced three times — see `.cdd/unreleased/626/self-coherence.md` §"run_class note" and `gamma-closeout.md`'s two "`run_class` taxonomy gap" entries), #614 (mandatory-learning doctrine + the concrete resumed-cell-disclosure precedent), #640 (adjacent sibling doctrine cell, same three surfaces, different gap shape — do NOT re-touch its content), #642 (ε detect-recurrence — **do NOT dispatch**).
**Authority already granted (κ, 2026-07-09 dispatch comment):** "Authority — A2/CAP: Fix stale prose, PR metadata, and link/doc drift inside scope without asking. Escalate only if the doctrine change requires behavior/FSM/code changes." α does not need to re-ask for routine doctrine-reconciliation authority within that envelope; this cell's whole deliverable *is* doctrine prose, so almost everything α does sits inside the A2/CAP grant. Escalate only if α finds the taxonomy fix is impossible without an FSM/`transitions.json`/label change (see Scope guardrails — this would mean the cell cannot ship as scoped).

---

## 1. Per-AC oracle list

### AC1 — one canonical, exhaustive, mutually-exclusive `run_class` set with observable discriminators and an ordered decision procedure; every shape observed in #626/#630/#614 maps to exactly one class

**Invariant:** after this cycle, `cds-dispatch/SKILL.md` §"Repair re-entry preflight" Step A's ordered decision procedure has enough named classes and discriminators that walking each of the five concretely-observed shapes below through the procedure produces exactly one result, and every class other than the terminal catch-all has a *positive* discriminator (not merely "none of the others").

**The five concrete shapes to walk (do not invent hypotheticals — these are real, already-recorded firings):**
1. `.cdd/unreleased/593/gamma-scaffold.md` — a prior `cds-dispatch` firing wrote `CLAIM-REQUEST.yml` and created `cycle/593`, but the FSM claim transition never applied; issue stayed `status:todo`, no CDD implementation artifact existed. Classified `first_pass`.
2. `.cdd/unreleased/630/gamma-scaffold.md` — clean first claim, `run_class: first_pass` (630 is the cycle that *invented* `resumed_from_matter` as future-facing doctrine; 630 itself was not an instance of it).
3. `.cdd/unreleased/614/alpha-closeout.md` — branch/artifacts pre-existed from an infra-interrupted prior firing, mechanically pattern-matches repair re-entry, but no operator rejection ever occurred. Classified `first_pass` with an explicit disclosure paragraph (not `resumed_from_matter`, because #630's doctrine — which names the scanner's "MECHANICAL reversion" audit-note comment as the trigger — did not exist yet at #614's time, and no such comment was posted for #614 specifically).
4. `.cdd/unreleased/626/self-coherence.md` §R0 — clean `first_pass` claim, normal cycle.
5. `.cdd/unreleased/626/self-coherence.md` §R1 and §R2 ("run_class note" sections, both) — operator posted an authorization comment reopening `status:todo` for a *new bounded AC scope* on an issue whose prior round(s) already converged and merged (R0→PR #635, R1→PR #637). No rejection evidence (rules out `repair_pass`), no scanner "MECHANICAL reversion" comment (rules out `resumed_from_matter`), branch was reset to a clean `main` base for R2 (zero commits beyond base at claim time — looks like `first_pass` at the branch level) but the *issue* carries prior converged/merged rounds. γ's own #626 artifacts named this "a fourth shape" three times and tentatively proposed the name `scope_continuation` — treat this as strong prior art (see Friction note 2), not a fresh naming decision, unless α has a concrete reason to diverge.

**Oracle:**
- For each of the five shapes above, cite the exact artifact path/section as evidence and show it maps to exactly one class under the new ordered procedure. Do this as a table in the doctrine addition itself (or in `self-coherence.md`, cross-linked) — not just asserted in prose.
- Confirm mutual exclusivity: for every pair of non-catch-all classes, name the single discriminator that cannot hold for both simultaneously (e.g., `repair_pass` requires operator-rejection evidence; `resumed_from_matter` requires the scanner's audit-note comment AND absence of rejection evidence; the new continuation class requires issue-level prior-converged-round evidence AND absence of both rejection evidence and the scanner comment).
- Negative check: only the terminal class (`first_pass`, the historical "otherwise" case) may be defined by exclusion. Every other class — `repair_pass`, `resumed_from_matter`, `manual_delta_repair`, `blocked`, and the new continuation class — must have a positive, independently-checkable discriminator named in the ordered procedure, not "none of the above."

### AC2 — the operator-authorized scope-continuation shape is a named class with defined re-entry behavior (does NOT re-scaffold merged matter, does NOT treat prior converged artifacts as a conflict)

**Invariant:** `cds-dispatch/SKILL.md`'s `run_class` enum gains a new named value for the #626-shape, and `delta/SKILL.md` gains a sibling subsection to §9.10 ("resumed-from-changes") and §9.11 ("resumed-from-mechanical-reversion") — §9.12 is already taken by #626's cell/substrate-identity-boundary doctrine, so the new subsection is **§9.13** — defining δ's re-entry behavior for this shape: does not unconditionally dispatch γ for a fresh R0 scaffold when `.cdd/unreleased/{N}/` already carries a converged-and-merged prior round's artifacts (mirroring §9.11 step 1's "if `gamma-scaffold.md` already exists, δ does NOT overwrite it and does NOT dispatch γ for a new R0 scaffold" pattern, adapted for the continuation case), and does not treat `CDDArtifacts` being non-empty + no rejection evidence as an anomaly requiring escalation.

**Oracle:**
- `grep` the new §9.13 for explicit "does NOT re-scaffold" / "does NOT treat ... as a conflict" language, structurally parallel to §9.11's existing phrasing.
- Confirm #626 R1/R2's actually-observed δ behavior (per `.cdd/unreleased/626/gamma-closeout.md` and `self-coherence.md`: δ resumed the bounded new-AC scope without re-scaffolding R0's already-merged matter) is cited as the concrete worked/empirical example, the same way §9.10 cites cycle/497 and §9.11 cites #630 as their respective first empirical witnesses.
- Confirm the new class's detection rule in `cds-dispatch/SKILL.md` Step A is checked in the correct order relative to the existing two checks (§9.11's own doctrine is explicit that the #630 check must run *before* the repair-re-entry check, because both can superficially resemble each other at a shallow glance — the new continuation check needs the same explicit ordering justification against the other three, not just being appended at the end without reasoning).

### AC3 — recovery-vs-resume and reset-branch-vs-first-pass are resolved explicitly (not left to the reader)

Two sub-questions, both from the issue body's "Observed shapes" §2:

**(a) Recovery vs. resume — one class or two?** Current doctrine (`delta/SKILL.md` §9.11 + `transitions.json`'s `propose_status_todo_with_matter` rule + `scan.go`) already treats "the reconciler mechanically requeues a dead-run-with-checkpointed-matter cell" (recovery) and "the next claim continues from that matter" (resume) as two phases of **one** class, `resumed_from_matter` — the reconciler produces the trigger condition, δ's §9.11 consumes it. α's job is to make this explicit rather than implicit: state plainly, in the doctrine, that recovery-and-resume is one `run_class` value spanning a produce/consume pair, not two candidate classes needing separate names — unless α finds concrete evidence (a real observed shape, not a hypothetical) that they need to diverge, in which case α names the split explicitly with its own discriminator.

**(b) Reset-branch-vs-first-pass — how is a clean-but-reset branch classified?** #626 R2 is the concrete empirical case: `cycle/626` was reset to a clean `main` base (zero commits beyond base at claim time — indistinguishable from a brand-new branch by branch-state alone) but the *issue* carried two already-merged prior rounds (PR #635, PR #637). γ's own R2 artifacts classified this as the continuation shape, NOT `first_pass`, precisely because branch-level state (commits beyond base, `BranchExists`) is the wrong discriminator here — it is issue-level convergence history (prior merged PRs referencing this issue number; `CDDArtifacts` non-empty from a prior round even after a branch reset; an operator continuation-authorization comment) that must gate the classification. α resolves this explicitly: the ordered decision procedure must check issue-level history (not just branch/commit state) before falling through to `first_pass`, and must say so in words, citing #626 R2 as the worked example of why branch-state alone is insufficient.

**Oracle:**
- `grep` the doctrine addition for an explicit sentence resolving (a) as one class (or an explicit justified split).
- `grep` the doctrine addition for an explicit statement that issue-level history (not branch-commit-count alone) gates the continuation-vs-first_pass distinction, with #626 R2 cited by artifact path.
- Confirm the ordered decision procedure in `cds-dispatch/SKILL.md` Step A actually implements this ordering (issue-level check happens before the branch-state-only fallthrough to `first_pass`), not just asserted in prose elsewhere while the procedure itself still checks branch state first.

### AC4 — the divergent enumerations across `dispatch-protocol/SKILL.md` / `delta/SKILL.md` / `cds-dispatch/SKILL.md` are reconciled; no active doctrine surface carries a divergent `run_class` list

**Invariant:** after this cycle, every doctrine surface that states the `run_class` enum states the *same* set of values (mutual reference is fine — e.g. "see `cds-dispatch/SKILL.md` §X for the canonical enum" — but no surface may state a *narrower or different* list as if it were complete).

**Currently divergent (verified by γ at scaffold time, re-verify at implementation time in case main moved):**
- `dispatch-protocol/SKILL.md` §2.8, line ~445: `run_class ∈ {first_pass, repair_pass, manual_delta_repair, blocked}` — **missing `resumed_from_matter` entirely**, and missing whatever new continuation class AC2 adds.
- `cds-dispatch/SKILL.md` §"Repair re-entry preflight" → `` `run_class` `` (lines ~207–216): `first_pass`, `resumed_from_matter`, `repair_pass`, `manual_delta_repair`, `blocked` — already has 5 values (this is the most complete list today) but will need the AC2 addition.
- `delta/SKILL.md` §9.10/§9.11/§9.12 do not restate the full enum as a list (they narrate individual shapes) — confirm no restated sub-list anywhere in §9 diverges from the reconciled canonical set once AC2's addition lands.

**Oracle:**
- `grep -n "run_class" src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` and confirm every enum-shaped list across all three matches, or the shorter ones defer explicitly to the canonical source rather than silently truncating.
- **Hard constraint (do not break CI):** `scripts/ci/check-dispatch-repair-preflight.sh` does a literal `grep -qF` for the substrings `run_class`, `first_pass`, `repair_pass`, `manual_delta_repair`, `blocked` in `dispatch-protocol/SKILL.md`, and for `run_class`, `first_pass`, `repair_pass` in `cds-dispatch/SKILL.md` + its golden + the live rendered workflow. **None of these literal substrings may be removed or renamed** — reconciling the enum means *adding* the missing value(s) to the shorter lists, never deleting or renaming the existing ones the CI guard greps for verbatim. Confirm this by re-running `./scripts/ci/check-dispatch-repair-preflight.sh` locally after the edit.
- If `cds-dispatch/SKILL.md`'s body text changes, its rendered golden (`cnos-cds-dispatch.golden.yml`) and the live workflow (`.github/workflows/cnos-cds-dispatch.yml`) MUST be regenerated via `cn install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml` (mechanical regeneration — never hand-edited), then verified byte-identical per the `install-wake-golden` gate. #626's own R1 commit message ("regenerated; mechanical; not hand-edited") is the precedent to follow.

### AC5 — doctrine-only; no code/FSM/label/dispatch-behavior change; gates green (I1, I2, I4, I5, I6, install-wake-golden, dispatch-repair-preflight, dispatch-closeout-integrity, Go, Package, Binary)

**Oracle:**
- `git diff --stat origin/main..HEAD` contains **zero** `.go` files and **zero** diff to `src/packages/cnos.cds/skills/cds/fsm/transitions.json`. The only non-`.md` files that may appear are the mechanically-regenerated `cnos-cds-dispatch.golden.yml` + `.github/workflows/cnos-cds-dispatch.yml` (if `cds-dispatch/SKILL.md`'s body changed), and even those must be byte-identical to a fresh `cn install-wake` render (no hand-edits).
- No label additions/removals anywhere in the diff (no new `status:*`, no new `dispatch:*`/`protocol:*` label definitions).
- CI green on every named gate for the merge commit: I1, I2, I4, I5, I6, `install-wake-golden`, `dispatch-repair-preflight`, `dispatch-closeout-integrity`, Go, Package, Binary. Since no `.go` files are touched, Go/Package/Binary gates are expected to be no-op-green; confirm they still run and pass rather than assuming.
- `scripts/ci/check-dispatch-closeout-integrity.sh` does not reference `run_class` at all (verified by γ — zero grep hits) — this cell's changes should not affect it; confirm it still passes as an unrelated-but-required gate.

---

## 2. Source-of-truth table

| Claim / surface | Canonical source | Status / role for α |
|---|---|---|
| `run_class` enum, most complete today | `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` §"Repair re-entry preflight" → §"`run_class` (recorded in the cycle receipt + this wake's return token)" (lines ~207–216) and §"Step A — classify the run" (lines ~153–166) | This is the wake's live prompt body (the actual dispatch wake that claimed #639 renders from this file) — the primary landing spot for AC1/AC2/AC3's ordered-decision-procedure rewrite and the new class addition. Edits here propagate to the golden + live workflow via `cn install-wake` (mechanical regeneration required, see AC4). |
| `run_class` enum, narrowest/stalest today | `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` §2.8 "Repair re-entry detection and preflight (cnos#516)" (lines ~435–447, enum at line ~445) | Framework-level doctrine (cnos.core, not cnos.cds) — this is where the divergence is worst (`{first_pass, repair_pass, manual_delta_repair, blocked}`, missing `resumed_from_matter` and the new class). Must gain the missing values. Also has §4.8 "Repair re-entry preflight check" (verify text) at line ~564–566 — confirm its cross-reference language stays accurate after the enum changes. |
| δ's re-entry-behavior doctrine per shape | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.10 (resumed-from-changes, cnos#500, lines ~518–568), §9.11 (resumed-from-mechanical-reversion, cnos#630, lines ~569–606), §9.12 (cell/substrate identity boundary, cnos#626 — **not a resume shape, do not confuse with 9.10/9.11's pattern**, lines ~607+) | AC2's new subsection is **§9.13**, sibling to §9.10/§9.11 in shape (trigger / evidence-on-the-branch / what-resume-means / repair-contract-applicability table, matching §9.11's own comparison-table format against §9.10). Do not renumber or edit §9.10/§9.11/§9.12's existing content — additive only. |
| CI guard — literal-string presence gate (NOT behavioral) | `scripts/ci/check-dispatch-repair-preflight.sh` | **Read-only reference; do NOT edit this file** (a CI script edit reads as code, out of scope per the issue's "No code" line). It does `grep -qF` for exact substrings (`run_class`, `first_pass`, `repair_pass`, `manual_delta_repair`, `blocked` in dispatch-protocol; `run_class`, `first_pass`, `repair_pass` in cds-dispatch/SKILL.md + golden + live workflow). Preserve every one of these literal substrings verbatim while reconciling the enum — see AC4's hard constraint. |
| CI guard — closeout integrity (unrelated) | `scripts/ci/check-dispatch-closeout-integrity.sh` | Read-only reference. Confirmed (γ grep) it does not mention `run_class` at all — unaffected by this cell, but still a required-green gate per AC5. |
| Go fact model — observable discriminators the doctrine's decision procedure may cite | `src/packages/cnos.issues/commands/issues-fsm/snapshot.go`'s `FactSnapshot` struct (fields: `BranchExists`, `CommitsBeyondBase`, `PRExists`, `PRCommitCount`, `CDDArtifacts []string`, `ReviewRequestPresent`, `ClaimRequestPresent`, `RepairContractPresent`, `CellKind`) | **Read-only reference — no Go change.** These are the only structurally-observable facts the FSM's claim-time snapshot exposes today. Note the gap: there is **no** typed field for "issue has prior merged/converged rounds" (the AC3(b) discriminator) — that fact is observable only via GitHub PR history / issue comments (not via `FactSnapshot`), which δ/the wake already read separately (issue body + comments, per `delta/SKILL.md` §9.2's five-input contract and the wake's own claim-time `gh issue view --json comments` read). Name this honestly in the doctrine — do not claim `FactSnapshot` already carries a field it doesn't; the decision procedure's issue-level check is a comment/PR-history read, not a new Go field, and building a new Go field for it is explicitly out of scope ("No code"). |
| The reconciler rule that produces the `resumed_from_matter` trigger | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` (`propose_status_todo_with_matter` rule, `in-progress` state, ~line 188) + `src/packages/cnos.issues/commands/issues-fsm/scan.go` (~lines 145, 244, 278–285) | **Read-only reference — no FSM/Go change** (explicit issue non-goal). Cite for AC1's discriminator description of `resumed_from_matter` (the scanner's "MECHANICAL reversion" audit-note phrase is this rule's `reason` text, posted verbatim by `scan.go`). |
| The five concrete worked examples for AC1's mapping table | `.cdd/unreleased/593/gamma-scaffold.md` + `gamma-closeout.md`; `.cdd/unreleased/630/gamma-scaffold.md` + `gamma-closeout.md`; `.cdd/unreleased/614/alpha-closeout.md` + `beta-closeout.md` + `gamma-closeout.md`; `.cdd/unreleased/626/self-coherence.md` (§R0/§R1/§R2, incl. both "run_class note" sections) + `gamma-closeout.md` (both "`run_class` taxonomy gap" entries) | Primary evidentiary source for AC1's five-shape mapping table. Quote/cite directly — these are the actual recorded firings, not hypotheticals γ invented. #626's artifacts already propose the tentative name `scope_continuation` twice; treat as strong prior art per Friction note 2. |
| `cn install-wake` regeneration command (if `cds-dispatch/SKILL.md` body changes) | `cn install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml` (named in `cds-dispatch/SKILL.md`'s own live-state banner, line ~95) | Mechanical, not hand-edited. Precedent: `.cdd/unreleased/626/gamma-closeout.md` "Files changed this round" — `cnos-cds-dispatch.golden.yml` + `.github/workflows/cnos-cds-dispatch.yml` "regenerated (mechanical; not hand-edited)". |
| Whether CELL-KINDS.md/CDD.md need a cross-reference (issue's own optional suggestion) | `src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md`, `docs/development/cdd/CDD.md` (pointer only; canonical is `src/packages/cnos.cdd/skills/cdd/CDD.md`) | γ's read: **not required.** `run_class` is a δ/dispatch-runtime classification (distinct from `cell_kind`, which is FactSnapshot's separate, already-defined observation seam per `CellKind` struct — see `snapshot.go` above). The issue's "if the taxonomy belongs there" is conditional and γ's assessment is it does not belong in CELL-KINDS.md/CDD.md; α may add a one-line cross-reference from CDD.md if genuinely cheap, but should not treat it as a required deliverable — confirm this read in `self-coherence.md` rather than silently skipping. |

---

## 3. α prompt

```text
You are α. Project: cnos.

Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 639 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/639
Tier 3 skills:
  - src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md
  - src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md
  - src/packages/cnos.cdd/skills/cdd/delta/SKILL.md (§9.10-§9.12 read in full; §9.13 is new)
  - src/packages/cnos.core/doctrine/KERNEL.md (§1.2, §2.2 — one source of truth per fact; cite, do not restate)

Read .cdd/unreleased/639/gamma-scaffold.md in full FIRST — it names the
per-AC oracle list (including the five CONCRETE prior firings you must
walk through the new decision procedure, not hypotheticals), the
source-of-truth table (including a hard CI-guard constraint you must not
violate), and pre-identified friction.

## What you are fixing

The `run_class` taxonomy that classifies each dispatch-wake firing
(first pass / repair re-entry / resumed-from-matter / etc.) is currently
enumerated INCONSISTENTLY across three doctrine surfaces:
dispatch-protocol/SKILL.md §2.8 lists 4 values (missing
resumed_from_matter entirely); cds-dispatch/SKILL.md lists 5 (the most
complete, but about to need a 6th); delta/SKILL.md narrates individual
shapes (§9.10, §9.11) without a canonical enum. On top of the
divergence, the taxonomy is INCOMPLETE: #626 hit an operator-authorized
"scope continuation" shape three times (a new bounded AC scope
dispatched on an issue whose prior round already converged and merged)
that matches none of the existing classes, and self-disclosed the gap
each time rather than silently mislabeling.

## Your job

1. Reconcile the three enumerations into ONE canonical list, stated once
   (cds-dispatch/SKILL.md, the wake's live prompt body, is the natural
   canonical home given AC4's oracle) with the others cross-referencing
   it rather than restating a competing list.
2. Add the missing operator-authorized scope-continuation class.
   #626's own artifacts (.cdd/unreleased/626/self-coherence.md +
   gamma-closeout.md) already proposed the tentative name
   `scope_continuation` TWICE across two rounds — treat this as strong
   prior art, not a fresh naming decision. Only diverge from it with an
   explicit, stated reason in self-coherence.md.
3. Resolve recovery-vs-resume (gamma-scaffold.md §1 AC3(a) — current
   doctrine already treats this as ONE class across a produce/consume
   pair; make that explicit) and reset-branch-vs-first-pass
   (gamma-scaffold.md §1 AC3(b) — #626 R2's concrete case: branch reset
   to clean main, but issue has prior merged rounds; the discriminator
   must be issue-level history, not branch-commit-count).
4. Write the ordered decision procedure so it actually implements (3)'s
   resolution — check issue-level continuation evidence BEFORE falling
   through to first_pass, not just asserted in prose while the checklist
   still checks branch state first.
5. Add delta/SKILL.md §9.13 (sibling to §9.10/§9.11 — §9.12 is already
   taken by #626's unrelated cell/substrate-identity-boundary doctrine)
   defining δ's re-entry behavior for the new class: does not
   re-scaffold merged matter, does not treat prior CDDArtifacts as a
   conflict. Cite #626 R1/R2 as the empirical worked example, exactly as
   §9.10 cites cycle/497 and §9.11 cites #630.

## Hard constraint — do not break the CI guard

scripts/ci/check-dispatch-repair-preflight.sh does a literal grep -qF
for the substrings `run_class`, `first_pass`, `repair_pass`,
`manual_delta_repair`, `blocked` in dispatch-protocol/SKILL.md, and
`run_class`, `first_pass`, `repair_pass` in cds-dispatch/SKILL.md + its
golden + the live rendered workflow. Reconciling the enum means ADDING
the missing values to the shorter lists — never removing or renaming
any of these literal substrings. Do NOT edit the .sh script itself (that
reads as code; out of scope). Re-run
./scripts/ci/check-dispatch-repair-preflight.sh locally to confirm
before requesting review.

If you edit cds-dispatch/SKILL.md's body text, its golden
(cnos-cds-dispatch.golden.yml) and the live workflow
(.github/workflows/cnos-cds-dispatch.yml) MUST be regenerated via
`cn install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml`
(mechanical — never hand-edited), verified byte-identical per
install-wake-golden. #626's own precedent commit
("regenerated; mechanical; not hand-edited") is the pattern to follow.

## Implementation contract (pinned by γ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Markdown/doctrine prose only. No Go, no shell, no YAML hand-edits (rendered YAML/golden files are machine-regenerated only, never hand-authored). |
| CLI integration target | N/A — no CLI surface changes. If cds-dispatch/SKILL.md's body changes, the ONLY command you run against it is `cn install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml` to regenerate the rendered artifacts; you do not add or modify any `cn` subcommand. |
| Package scoping | N/A — no new package, no new Go file, no new directory under src/go or src/packages/*/commands. Edits are confined to the three named SKILL.md doctrine files (plus their mechanically-regenerated rendered artifacts if cds-dispatch/SKILL.md changes). |
| Existing-binary disposition | N/A — no binary changes. |
| Runtime dependencies | N/A — no new runtime dependency. |
| JSON/wire contract preservation | N/A — no wire contract touched. FactSnapshot's Go struct (snapshot.go) is read-only reference material for this cell, not something you modify. |
| Backward-compat invariant | transitions.json byte-identical (no diff at all — this is a hard check, not a judgment call). No status/dispatch/protocol label added, removed, or redefined. No existing run_class value (first_pass, repair_pass, resumed_from_matter, manual_delta_repair, blocked) is renamed or removed — only additive: the missing values are added where absent, and one new value names the scope-continuation shape. |

## Scope guardrails (see gamma-scaffold.md §5 for the full list — restated here for emphasis)

Do NOT: write or modify any .go file · change transitions.json (byte-identical,
no diff) · add/remove/redefine any status/dispatch/protocol label ·
enforce cell_kind anywhere · start Demo 0 · touch #626, sparse-checkout,
or write-fence content · dispatch #642 · edit scripts/ci/*.sh · hand-edit
cnos-cds-dispatch.golden.yml or .github/workflows/cnos-cds-dispatch.yml
(regenerate only, via cn install-wake).

## Before requesting review

Write .cdd/unreleased/639/self-coherence.md per alpha/SKILL.md, mapping
each AC to concrete evidence (doctrine section names/paths, the AC1
five-shape mapping table, the AC4 grep confirming reconciliation, the
CI-guard re-run output). State your final class name for the
scope-continuation shape explicitly, with reasoning if you diverged from
`scope_continuation`. State your CELL-KINDS.md/CDD.md cross-reference
decision explicitly (gamma-scaffold.md §2's read: not required, but your
call to confirm or override).
```

---

## 4. β prompt

```text
You are β. Project: cnos.

Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 639 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/639

Read .cdd/unreleased/639/gamma-scaffold.md §1 (per-AC oracle list,
INCLUDING the five concrete prior-firing shapes) and §2 (source-of-truth
table, including the CI-guard hard constraint) FIRST — walk α's diff
against that oracle list independently. Do not trust α's
self-coherence.md framing; re-derive.

## AC-by-AC verification

- AC1: for each of the five concrete shapes named in gamma-scaffold.md
  §1 AC1 (#593, #630, #614, #626 R0, #626 R1/R2), confirm α's diff
  produces exactly one classification per shape under the NEW ordered
  procedure, with the artifact citation present (not just asserted).
  Confirm every class except the terminal fallback has a positive
  discriminator, not a "none of the above" definition.
- AC2: confirm cds-dispatch/SKILL.md's run_class enum has the new class
  name, and delta/SKILL.md has a new §9.13 (not overwriting §9.10/§9.11/
  §9.12) with explicit "does NOT re-scaffold" / "does NOT treat ... as a
  conflict" language, citing #626 R1/R2 as the empirical example. Check
  the new class's detection-order justification in Step A (why it's
  checked where it's checked relative to the other three checks).
- AC3: confirm an explicit sentence resolves recovery-vs-resume as one
  class (or an explicit justified split — check the justification is
  evidence-based, not hypothetical). Confirm an explicit sentence names
  issue-level history (not branch-commit-count alone) as the
  reset-branch discriminator, citing #626 R2. Confirm the ordered
  decision procedure ACTUALLY implements this ordering — re-read Step A
  line by line and check the issue-level check precedes the
  first_pass fallthrough, not just prose elsewhere claiming it does.
- AC4: `grep -n "run_class" src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md
  src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md
  src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` and confirm no surface
  states a narrower/divergent list as if complete. Run
  ./scripts/ci/check-dispatch-repair-preflight.sh locally and confirm it
  passes (proves none of the CI-guarded literal substrings were removed).
  If cds-dispatch/SKILL.md's body changed, confirm the golden + live
  workflow were regenerated via `cn install-wake` (not hand-edited) and
  are byte-identical to a fresh render.
- AC5: `git diff --stat origin/main..HEAD` has zero .go files, zero diff
  to transitions.json, and the only non-.md files (if any) are the
  mechanically-regenerated golden + live workflow. Confirm CI green on
  I1, I2, I4, I5, I6, install-wake-golden, dispatch-repair-preflight,
  dispatch-closeout-integrity, Go, Package, Binary on the merge commit.

## Standard β obligations

- Implementation-contract coherence (Rule 7): verify the diff conforms
  to every pinned axis in gamma-scaffold.md §3's Implementation contract
  table (doctrine/Markdown only; no Go/shell/YAML hand-edits).
- Non-goal / out-of-scope check: confirm no diff touches #626 content,
  sparse-checkout/write-fence material, scripts/ci/*.sh, any .go file,
  transitions.json, or any label definition; confirm #642 was not
  dispatched.
- gamma-scaffold.md presence gate (rule 3.11b): already satisfied — this
  file exists on cycle/639 before α's dispatch.
```

---

## 5. Scope guardrails

Restated verbatim (or near-verbatim) from the issue body's "Out of scope / Do NOT" list, plus κ's dispatch-comment "Do NOT":

- **No code.** No `.go` file changes.
- **No FSM/`transitions.json` change** — byte-identical, no diff at all.
- **No dispatch-behavior change** — this issue documents the classes; wiring the classifier to enforce them is a possible follow-on, not this cell.
- **No new status labels.**
- **No `cell_kind` enforcement** — `cell_kind` is a separate, already-defined observation seam (`FactSnapshot.CellKind`); this cell does not touch it.
- **No Demo 0.**
- **No changes to #626 / sparse-checkout / write-fence work.**
- **No lifecycle-semantics change.**
- **Do NOT dispatch #642.**
- **(γ-added, not literal issue text)** Do NOT edit `scripts/ci/*.sh` — the repair-preflight guard's literal-string checks must be satisfied by *doctrine* edits alone; editing the guard itself would read as a code change and is out of scope.
- **(γ-added)** Do NOT hand-edit `cnos-cds-dispatch.golden.yml` or `.github/workflows/cnos-cds-dispatch.yml` — regenerate mechanically via `cn install-wake` only, and only if `cds-dispatch/SKILL.md`'s body actually changed.
- **Escalation boundary** (κ's dispatch comment, A2/CAP): fix stale prose, PR metadata, and link/doc drift inside scope without asking; escalate only if the doctrine change would require behavior/FSM/code changes. If α discovers the taxonomy genuinely cannot be reconciled without an FSM/code change, that is an escalation trigger, not a license to make the change anyway.

---

## 6. Friction notes

1. **`cell_kind: doctrine` — explicit, no hedge.** Unlike #640 (which hedged between `doctrine` and `implementation` because a CLI primitive was a live candidate), #639's scope is unambiguously doctrine-only: the issue's own header states "Mode: design-and-build (doctrine/design only)" and its Do-NOT list forbids code, FSM changes, and label changes outright. γ pins `doctrine` with no reclassification escape hatch offered to α — if α's diff turns out code-dominant, that is itself a scope violation to flag and stop on, not a relabeling opportunity.
2. **Prior-art naming: `scope_continuation`.** γ found this is not a naming decision α starts fresh on — #626's own `self-coherence.md` (§R1 "run_class note") and `gamma-closeout.md` (both R1 and R2 "`run_class` taxonomy gap" entries) already proposed the exact tentative name `scope_continuation` twice, independently, across two different rounds of the same cell. This is strong convergent prior art from the actual agent that lived through the ambiguity three times. γ instructs α to treat it as the default and only diverge with a stated reason, to avoid re-litigating a naming question that direct empirical evidence (not this scaffold's opinion) already answered twice.
3. **The CI-guard literal-string trap.** `scripts/ci/check-dispatch-repair-preflight.sh` (cnos#516's regression guard) does exact `grep -qF` matches for `first_pass`, `repair_pass`, `manual_delta_repair`, `blocked` in `dispatch-protocol/SKILL.md`, and `first_pass`/`repair_pass` in `cds-dispatch/SKILL.md` + its golden + the live workflow. Because this guard predates AC4's reconciliation work, it is easy for a well-meaning prose cleanup pass to accidentally rename or restructure one of these exact substrings away (e.g. rewording "first_pass" into "first-pass run" in prose) and silently break the gate. γ flags this explicitly in both the AC4 oracle and the α/β prompts rather than trusting it to be caught only by CI red after the fact.
4. **No Go field observes "issue has prior merged/converged rounds."** γ checked `FactSnapshot` (`snapshot.go`) directly: `BranchExists`, `CommitsBeyondBase`, `PRExists`, `PRCommitCount`, `CDDArtifacts`, `ReviewRequestPresent`, `ClaimRequestPresent`, `RepairContractPresent`, `CellKind` — none of these is "this issue number has a prior merged PR / release history." AC3(b)'s resolution (issue-level history, not branch-commit-count, gates the continuation-vs-first_pass call) is therefore a claim about a fact the wake/δ must read via GitHub comments/PR-search today, not a fact already sitting in the typed snapshot. γ names this honestly in the source-of-truth table so α doesn't accidentally claim `FactSnapshot` already carries a field it doesn't, and doesn't try to add one (that would be a Go change, out of scope).
5. **`delta/SKILL.md` §9 numbering gap.** §9.9 is "Cross-references" (pre-existing, appears to close the original §9 block); §9.10/§9.11/§9.12 were appended later as siblings, not renumbered into §9's original body. AC2's new subsection continues that same append-only pattern as **§9.13** — γ confirmed by reading all four headers (`grep -n "^### 9\."`) that no §9.13 currently exists. Do not renumber 9.10–9.12.
6. **Branch base-SHA: no drift.** Unlike #640's scaffold (which had to note a one-commit drift), γ confirmed for #639 that `git merge-base cycle/639 origin/main` equals `7313caea3b6b9901cc02e797be9c929393b971ee` exactly — the SHA cited in the wake-invoked inputs. No rebase note needed.
7. **`.cdd/unreleased/639/CLAIM-REQUEST.yml` already exists on `cycle/639`** (δ's claim-sequence artifact, `run_class: first_pass` per that file's own contents — this cell's own dispatch is, unsurprisingly, an ordinary first pass, not an instance of the taxonomy gap it's fixing). γ leaves it untouched; not γ's artifact to modify.
