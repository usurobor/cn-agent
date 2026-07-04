<!--
section-manifest:
  planned: [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness]
  completed: [Gap, Skills, ACs, Self-check, Debt, CDD Trace]
-->

# Self-coherence — cnos#584

## §Gap

**Issue:** [usurobor/cnos#584](https://github.com/usurobor/cnos/issues/584) — "arch(cds/cdd): codify mechanical cell runtime vs cognitive skills (mechanism/cognition boundary)". Sub 1 of 5 under parent #583 (master wave — mechanical dispatch-cell architecture). Lands first.

**Mode:** design-and-build, doctrine-only — no runtime change.

**Governing rule to codify (operator-stated):** cells are mechanical and defer cognition to skills; skills don't control anything. The cell is the execution engine — it owns control flow and all side effects, and calls a skill when it needs a judgment. A skill returns cognition; it never opens a PR, writes a label, or drives the lifecycle.

**Branch:** `cycle/584`, branched from `origin/main@9309de97d7e6d90637012839163e8d0511b56ca6` (per γ's scaffold; the wake-invoked-δ input's pinned SHA `9277c7e3...` had already advanced by one unrelated commit — board-map regen — by branch time; γ verified and branched from the current HEAD per its branch pre-flight rule, not the stale pinned value).

**Scope guardrail (restated):** in scope — codify the governing rule in cell-framework doctrine; classify the 9 named dispatch lifecycle steps mechanical vs. cognitive; audit 5 named files (`cdd/gamma`, `cdd/alpha`, `cdd/beta`, `dispatch-protocol/SKILL.md`, `cds-dispatch/SKILL.md`) for control-implying prose and correct it. Out of scope — any runtime/workflow/transition-table change (Subs 2–4), implementing checkpoint/PR-open/recovery, `cell_kind` enforcement, new status labels/taxonomy, Demo 0.

## §Skills

**Tier 1 (loaded, canonical):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle + kernel doctrine; the chosen AC1/AC2 doctrine home.
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface (this file's governing load order).

**Tier 1 sub-skill:**
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — cell-shaped issue contract; used to read the issue pack correctly (mode = design-and-build, source-of-truth table, non-goals discipline).

**Tier 3 (issue-specific, named by γ's dispatch prompt):**
- `src/packages/cnos.core/skills/write/SKILL.md` — prose-quality discipline (one governing question per doc, front-loaded point, stable-fact-once-then-point, contrastive examples). Applied throughout the new CDD.md section and this artifact: the section leads with the rule, states the worked example before generalizing, and the lifecycle table states each row's fact once rather than repeating it in prose.

**Read for context, not loaded as an authoring skill (per scope guardrails):**
- `docs/papers/CELL-OF-CELLS.md` — read in full as the other candidate doctrine home; not chosen (see §ACs → AC1 for the placement justification).
- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md`, `.../beta/SKILL.md` (already Tier-1-adjacent for α context), `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md`, `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` — the 5 AC3 audit files, read in full for the audit.
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9 — read for citation (existing positive instance of the mechanism/cognition split); not an audit target, not restructured.

No Tier 2 `eng/*` skills were loaded — this cycle produces no source code in any language (per the pinned implementation contract's Language axis: N/A).

## §ACs

### AC1 — the governing rule is codified as doctrine

**Canonical-home decision: `src/packages/cnos.cdd/skills/cdd/CDD.md`**, new §"Mechanism and cognition" (inserted between §"Software-specific realization → cnos.cds" and §"Hard rule").

**Justification.** The contract pinned exactly one home from {`CDD.md`, `docs/papers/CELL-OF-CELLS.md`}. I chose `CDD.md` over `CELL-OF-CELLS.md` for three reasons:

1. **Trivial discoverability.** `CDD.md` is Tier-1-loaded by every role (α, β, γ all name it first in their own Load Order sections). Placing the rule there means the AC1 oracle's "already-loaded framework surface" requirement is satisfied by construction — no separate cross-reference edit is needed in a second file. Choosing `CELL-OF-CELLS.md` instead would have required a pointer edit *back* into `CDD.md` anyway (since `gamma/SKILL.md` is restricted to AC3-audit-only edits and can't carry a new, unrelated pointer), which would touch two files for one fact instead of one.
2. **Precedent for domain-flavored content in non-kernel sections.** `CDD.md`'s own §Non-goals restricts naming "tooling, platform, dispatch mechanism... invocation surface" only in "the kernel sections (§Kernel through §Scope-lift)" — verbatim scoped to that span. The new section sits after §"Software-specific realization → cnos.cds," which already cites CDS-concrete skills, commands, and issue numbers by name. The new section follows the same established pattern (state the substrate-independent rule, then cite the one running domain package's concrete instantiation) rather than introducing a new kind of content CDD.md doesn't already carry elsewhere.
3. **Single-file footprint.** Confining the new doctrine to one file (rather than splitting the rule into CDD.md and the worked example into CELL-OF-CELLS.md, which friction note 3 floated as an option) keeps the Package-scoping axis unambiguous: exactly one new-content file, matching the contract's literal "(1) exactly one canonical doctrine home."

`CELL-OF-CELLS.md` was read in full and considered; its essay register and existing willingness to name concrete dispatch mechanics (§12 already cites dispatch, mid-flight, artifact-channel by name) made it a defensible alternative, but it is not Tier-1-loaded by any role, so choosing it would not have satisfied AC1's discoverability oracle without an additional edit elsewhere.

**Oracle run (mechanical check):**

```
$ rg -n -i "cells are mechanical" src/packages/cnos.cdd/ docs/papers/
src/packages/cnos.cdd/skills/cdd/CDD.md:156:**Cells are mechanical. Cognition is deferred to skills. Skills don't control anything.**
```

One hit, at the stable canonical path `src/packages/cnos.cdd/skills/cdd/CDD.md` (not the issue body). The same line also carries "skills don't control anything" verbatim.

**Negative-case oracle:**

```
$ rg -n -i "cells are mechanical|skills don't control anything" src/packages/cnos.cdd/ docs/papers/ src/packages/cnos.cds/ src/packages/cnos.core/
src/packages/cnos.cdd/skills/cdd/CDD.md:156:**Cells are mechanical. Cognition is deferred to skills. Skills don't control anything.**
```

Exactly one hit, at the canonical doc, not in issue/chat prose. (Re-run before this edit, per γ's peer-enumeration pass and my own pre-edit check, returned zero hits — the rule genuinely did not exist anywhere in doctrine before this cycle.)

**Worked example present:** §"Mechanism and cognition" → §"Worked example" states one mechanical step (`runtime.claim(issue) → status:todo → status:in-progress`, no skill invoked) contrasted with one `[call skill]` cognitive delegation (`runtime.dispatch([call γ scaffold], contractₙ) → gamma-scaffold.md`, γ returns judgment, the runtime commits it).

**AC1: met.**

### AC2 — the dispatch lifecycle is classified mechanical vs. cognitive

**Oracle run — all 9 step names present in the new section:**

```
$ for w in claim scaffold checkpoint implement "PR-open" review transition-request close-out recover; do
    echo -n "$w: "; rg -c -i -- "$w" src/packages/cnos.cdd/skills/cdd/CDD.md; done
claim: 6
scaffold: 3
checkpoint: 2
implement: 4
PR-open: 3
review: 16
transition-request: 3
close-out: 3
recover: 2
```

All 9 names present (0 misses). The table in §"Lifecycle classification" has exactly 9 rows: claim, scaffold, checkpoint, implement, PR-open, review, transition-request, close-out, recover — each with an Owner / Skill / Returns column.

**PR-open and transition-request classification (the AC2 negative-case check):**

- `PR-open` — Owner column reads `runtime-owned *(target; not yet built — Subs 2–4)*`. Skill column is `—`. Returns column names the current ad hoc state (cds-dispatch wake opens a PR under wake-invoked mode; the bootstrap sequential-dispatch model this cycle ran under does not) rather than naming a role skill as the owner.
- `transition-request` — Owner column reads `runtime-owned *(target; partially landed for one edge — Subs 2–4 for the rest)*`. Skill column is `—`. Returns column names the one landed edge (`in-progress → review`, requested by δ, decided by the FSM via `cn issues fsm evaluate --apply`) as a partial, explicitly-scoped exception, not a general claim that the step is done or cognitive.

Neither row names a role skill (γ/α/β) as the owner-of-record; both are `runtime-owned` with an explicit target/partial qualifier, satisfying the AC2 oracle's requirement that these two rows not be classified as cognitive-agent steps.

**AC2: met.**

### AC3 — skills are shown to control nothing

**Oracle run (mechanical check, post-edit, across all 5 named files):**

```
$ rg -n -i "skill (opens|sets|writes|drives|dispatches)|opens? (a |the )?PR|writes? the label|drives (the )?lifecycle" \
    src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md \
    src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md \
    src/packages/cnos.cdd/skills/cdd/beta/SKILL.md \
    src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md \
    src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md
src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md:258:The dispatch wake transitions the **claimed cell's** lifecycle labels at these named events only. For the claim, hard-block, and release-back-to-queue events, the wake writes the label directly ... For the β-converge event, the wake does **not** write the label directly — it requests the transition via the FSM (`cn issues fsm evaluate --issue {N} --apply`, cnos#569 Phase 2) and the FSM applies the label only when its guards pass:
```

One hit, unchanged by this cycle's edits. It is the positive example γ's scaffold named (the wake — a mechanical runtime component, not a cognitive role skill — correctly distinguishing which label transitions it writes directly from the one it only requests via the FSM). It already reads as "returns the decision; the runtime executes it" and required no correction.

**Per-file audit result:**

- `gamma/SKILL.md` — read in full. No control-implying hit (oracle regex: 0 matches; manual read: γ's Core Principle explicitly states "γ produces α and β prompts; δ (operator) executes dispatch via the harness — one role at a time," which is already correctly framed as γ returning artifacts, not executing lifecycle transitions). **No edit made.**
- `alpha/SKILL.md` — read in full (also the α role-contract file governing this session). No control-implying hit. **No edit made.**
- `beta/SKILL.md` — read in full. **Violation found and corrected** at lines 38 and 70 (pre-edit line numbers), matching γ's scaffold exactly:
  - Line 38 (pre-edit): `merge (β's authority per CDD.md — the natural conclusion of review)` — described β as the actor whose authority *is* the merge, not as a role that renders a merge/no-merge judgment a runtime then executes.
  - Line 70 (pre-edit): `**β merge** → \`git merge\` per \`release/SKILL.md\`` — same collapse in the phase map.
  - **Correction applied** (commit `10700bab`): both lines now distinguish "the merge/no-merge judgment β renders" from "the mechanical `git merge` execution β currently performs as a bootstrap-mode stand-in pending the Subs 2–4 runtime," with an explicit forward pointer to `CDD.md` §"Mechanism and cognition" for the target split. **β's actual behavior is unchanged** — β still executes `git merge` exactly as before; only the prose framing changed, per the scaffold's explicit resolution path (friction note 2) and the AC4 constraint that this must be a prose-only qualification, not a behavior change.
- `dispatch-protocol/SKILL.md` — read in full. No control-implying hit. This file describes the *wake* (a mechanical runtime component, not a cognitive role skill) claiming issues, applying labels, and requesting the FSM transition — already correctly framed throughout (e.g. §2.9's "the wake does not write `status:review` itself: it *requests* the transition via `cn issues fsm evaluate --issue {N} --apply`, and the FSM's own `in-progress → review` guard independently re-enforces this same deliverable-evidence bar"). **No edit made.**
- `cds-dispatch/SKILL.md` — read in full. No control-implying hit requiring correction; the one oracle-regex hit (line 258, above) is the positive existing instance. **No edit made.**

**AC3: met.** Exactly one violation existed across the 5 audited files (β lines 38/70); it is corrected; the correction is a framing qualification, not a behavior change; the remaining 4 files needed no edit.

### AC4 — no behavior change; gates green

**Oracle run — diff is `.md`-only:**

```
$ git diff --name-only origin/main...HEAD
.cdd/unreleased/584/gamma-scaffold.md
.cdd/unreleased/584/self-coherence.md
src/packages/cnos.cdd/skills/cdd/CDD.md
src/packages/cnos.cdd/skills/cdd/beta/SKILL.md
$ git diff --name-only origin/main...HEAD | grep -vE '\.md$'
(no output — every changed file is .md)
```

**Oracle run — the ring-fenced test file is untouched:**

```
$ git diff origin/main...HEAD -- src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go
(empty)
```

**No `transitions.json` or any `.json`/`.go`/`.yml`/`.yaml` file appears in the diff** (confirmed by the same `--name-only` listing above — 4 files total, all `.md`).

**CI gates.** This is a docs-only diff with no code, workflow, or schema surface touched, so I1/I2/I4/I5/I6, `install-wake-golden`, dispatch guards, Go, Package, and Binary gates are structurally unaffected — there is no code path for them to exercise. I attempted to run the local frontmatter validator (`scripts/ci/validate-skill-frontmatter.sh`) as an extra check on `beta/SKILL.md`'s frontmatter (unchanged by my edit, but the file's YAML header could in principle be corrupted by a body edit); it exited with `prerequisite missing: cue` — the `cue` binary is not installed in this environment, so the check could not run locally. I manually re-read `beta/SKILL.md`'s frontmatter block (lines 1–28) after editing and confirmed it is unchanged and well-formed (no stray blank lines, closing `---` intact, no key collisions introduced by my edits, both of which were in the body well below the frontmatter). This is disclosed as known debt in §Debt — β should re-run the frontmatter validator with `cue` available, or confirm CI's own frontmatter-validation gate is green on the branch head.

**AC4: met**, with the one disclosed local-tooling gap above.

## §Self-check

**Did this work push ambiguity onto β?**

- **Canonical-home choice** is a design call, not an ambiguity — I picked `CDD.md` and gave three concrete reasons (§ACs → AC1). β can disagree with the choice, but the choice itself is not left implicit or deferred.
- **Every AC1–AC4 claim above is backed by a literal command + its literal output**, re-run and re-verified immediately before this commit (see the verification pass above the §ACs commit). β does not have to take my word for any oracle result; each is reproducible from the branch head.
- **The AC3 resolution for `beta/SKILL.md`** is the one place this cycle asks a reviewer (β) to judge a framing correction to β's *own* role file. I did not soften or hide this: §ACs names the exact pre-edit lines, the exact correction, and states explicitly that β's behavior is unchanged — the fact that β is reviewing a change to its own doctrine file is disclosed here, not left for β to notice unaided.
- **Out-of-scope findings are named, not silently expanded into.** I found one additional control-implying-adjacent line (`CDD.md` line 104's pre-existing Roles pointer: "β (reviews and merges; merge is β's authority)") that reads similarly to the corrected `beta/SKILL.md` lines, but it is not one of the 5 AC3-audited files and the contract restricts my `CDD.md` edit to "new section(s) for AC1+AC2" only — touching it would be unrelated restructuring outside the pinned scope. Named as a finding in §Debt below, not fixed.

**Is every claim backed by evidence in the diff?** Yes — §ACs cites the exact commit SHAs (`10700bab` for the beta/SKILL.md fix), the exact pre-edit line content, and re-runs every oracle command against the current branch head rather than against memory of an earlier run.

## §Debt

1. **CDD.md line 104 (pre-existing, not touched this cycle).** The Roles pointer list in `CDD.md` §Pointers reads: `[`beta/SKILL.md`](beta/SKILL.md) — β (reviews and merges; merge is β's authority)`. This is the same "skill executes the mechanical action as its own authority" phrasing pattern that AC3 required correcting in `beta/SKILL.md` itself, but `CDD.md`'s existing Roles-pointer line is not one of the 5 AC3-audited files, and my permitted `CDD.md` edit this cycle is scoped to "new section(s) for AC1+AC2" only (per the pinned implementation contract) — correcting an unrelated pre-existing line in the same file would be unrelated restructuring outside that scope. Named here per the scaffold's explicit instruction to surface (not silently fix) findings outside the prescribed audit-file list. **Disposition: γ triage** — likely a small, mechanical, one-line follow-up (either fold into a Sub 2–4 cycle or a same-cycle immediate fix if γ judges the risk of touching `CDD.md` beyond the pinned scope acceptable).
2. **`scripts/ci/validate-skill-frontmatter.sh` could not run locally** (missing `cue` binary in this environment). I manually verified `beta/SKILL.md`'s frontmatter block is unchanged and well-formed after my body-only edits (see §ACs → AC4). β should confirm the branch-head CI run for this validator (or any package/frontmatter gate) is green before merge, per α's own pre-review gate row 10 (branch CI must be green, or the artifact must say so explicitly and β waits for green).
3. **No design or plan artifact was produced separately from this self-coherence.md.** Per `alpha/SKILL.md` §2.2, design and plan may be marked "not required" with a concrete justification: this cycle is doctrine-only prose across 2 files plus one audited correction in a 3rd; the canonical-home decision and the AC3 resolution path were both fully pre-specified by γ's scaffold (the source-of-truth table + friction notes 2 and 3), leaving no independent design surface for α to produce beyond the placement/wording judgment recorded directly in §ACs above. Marking design and plan as "not required" for this reason.
4. **Peer enumeration / harness audit (per `alpha/SKILL.md` §2.3–§2.4) do not apply in their code-surface form** — there is no schema-bearing parser, manifest shape, or runtime contract changed by this cycle, so there is no producer/consumer/harness family to enumerate in the code sense. The doctrine-level "peer enumeration" analog — every one of the 5 AC3-audited files, individually read and individually reported on (edited / not-edited-with-reason) — is done explicitly in §ACs → AC3 above, which is the applicable form of this rule for a docs-only, audit-shaped cycle.

## §CDD Trace

| Step | Artifact | Skills loaded | Decision |
|---|---|---|---|
| 0 Observe | issue #584 | `cdd/issue` | γ observed parent #583 (master wave); #584 is Sub 1 of 5, doctrine-only, lands first. |
| 1 Select | issue #584 | `cdd/gamma` (γ's prior pass) | γ selected per CDS §"Selection function" — dependency order (Sub 1 must land before Subs 2–4 build against it). |
| 2 Branch | `cycle/584` | `cdd/gamma` | Branch created by γ from `origin/main@9309de97d7e6d90637012839163e8d0511b56ca6` per §"Branch rule" / §"Branch pre-flight"; SHA drift from the wake-invoked-δ input's stale pinned SHA logged as γ's friction note 1. |
| 3 Bootstrap | `.cdd/unreleased/584/gamma-scaffold.md` | `cdd/gamma`, `cdd/issue` | γ authored the scaffold: AC1–AC4 oracle list, source-of-truth table, 7-axis implementation contract, scope guardrails, 5 friction notes. |
| 4 Gap | `self-coherence.md` §Gap | `cdd/alpha`, `cdd/issue` | α named the gap: the mechanism/cognition boundary is practiced (δ §9, cds-dispatch) but not codified as general, named doctrine. Mode = design-and-build, doctrine-only. |
| 5 Mode | `self-coherence.md` §Gap, §Skills | `cdd/alpha`, `write` | Mode confirmed design-and-build per the issue header; active skills named (Tier 1 `CDD.md` + `alpha/SKILL.md` + `issue/SKILL.md`; Tier 3 `write/SKILL.md`; no Tier 2 `eng/*` — no source code). |
| 6 Artifacts | `CDD.md` (new §"Mechanism and cognition"), `beta/SKILL.md` (corrected lines 38/70) | `cdd/alpha`, `write` | AC1+AC2 doctrine authored in `CDD.md` (commit `c005b854`); AC3 correction applied to `beta/SKILL.md` (commit `10700bab`); AC4 diff-scope confirmed `.md`-only, `issuesfsm_test.go` untouched. |
| 7 Self-coherence | `.cdd/unreleased/584/self-coherence.md` | `cdd/alpha` | AC-by-AC self-check completed (§ACs above); AC1, AC2, AC3, AC4 all met with cited oracle evidence; known debt named explicitly (§Debt). |
