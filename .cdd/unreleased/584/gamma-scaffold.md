# γ Scaffold — cnos#584

**Issue:** [usurobor/cnos#584](https://github.com/usurobor/cnos/issues/584) — "arch(cds/cdd): codify mechanical cell runtime vs cognitive skills (mechanism/cognition boundary)"
**Parent:** #583 (master wave — mechanical dispatch-cell architecture). **Sub 1 of 5. Lands first.**
**Mode:** design-and-build, **doctrine-only — no runtime change**.
**Protocol:** `cds`.
**Base SHA:** `9309de97d7e6d90637012839163e8d0511b56ca6` (`origin/main` HEAD at scaffold time).

> **Note on the pinned base SHA drift.** The wake-invoked-δ input contract named `9277c7e3585e8ac054fa212fe96cf9aa1d9d011e` as the current main SHA. At branch pre-flight, `git rev-parse origin/main` returned `9309de97d7e6d90637012839163e8d0511b56ca6` instead. The single intervening commit (`9309de97` — "board-map: regenerate docs/development/board from live open issues") is an unrelated docs-generation commit with no overlap with this cycle's surfaces. Per γ's branch pre-flight discipline (`gamma/SKILL.md` §2.5 Step 3a — "base SHA known"), γ branched from the verified current HEAD rather than the stale pinned SHA. `cycle/584` was created from `origin/main@9309de97`.

**Wake run id:** `https://github.com/usurobor/cnos/actions/runs/28696482644`

---

## Governing rule to codify (operator-stated, per issue body)

> **Cells are mechanical and defer cognition to skills. Skills don't control anything.**

The cell is the execution engine: it owns control flow and all side effects, and **calls** a skill when it needs a judgment. A skill returns cognition; it never opens a PR, writes a label, or drives the lifecycle. Mechanism ≠ policy; policy (cognition) does not execute.

γ's read: this is the same distinction the repo has already been building toward empirically (δ's wake-invoked-mode §9 in `delta/SKILL.md` already separates "δ requests the transition" from "the FSM decides whether the guards allow it" — see `cn issues fsm evaluate --issue {N} --apply`). Sub 1's job is to make that distinction **explicit, general doctrine** that the framework role skills (γ/α/β) and the two audited dispatch-mechanics skills conform to in *prose*, before Subs 2–4 touch any runtime. This sub does **not** build the mechanical runtime that executes PR-open / transition-request as machine actions — it names the boundary the later subs build against.

---

## Per-AC oracle list

### AC1 — the governing rule is codified as doctrine

- **Invariant:** the cell-framework doctrine states "cells are mechanical; cognition is deferred to skills; skills don't control anything," with a worked example (one mechanical step vs. one `[call skill]` cognitive delegation).
- **Oracle (mechanical check):** `rg -n -i "cells are mechanical" src/packages/cnos.cdd/ docs/papers/` returns a hit at a **stable canonical path** (not the issue body); and that path is `rg`-discoverable *from* at least one already-loaded framework surface (`CDD.md` and/or `gamma/SKILL.md` and/or `CELL-OF-CELLS.md` contain a cross-reference — a relative link or explicit path string — pointing at the new doctrine section).
- **Negative case:** `rg -i "cells are mechanical|skills don't control anything" src/packages/cnos.cdd/ docs/papers/ src/packages/cnos.cds/ src/packages/cnos.core/` returns **zero** hits outside issue #583/#584 text — i.e. the rule lives only in chat/issue prose, not in a doc a future session loads by default.
- **Surface:** cell-framework doctrine — α picks the canonical home (see Source-of-truth table below) and states the choice explicitly in `self-coherence.md`.

### AC2 — the dispatch lifecycle is classified mechanical vs. cognitive

- **Invariant:** each of the 9 named lifecycle steps — **claim, scaffold, checkpoint, implement, PR-open, review, transition-request, close-out, recover** — is labeled `runtime-owned` or `skill-call`, naming the skill (γ/α/β/δ/none) and what it returns.
- **Oracle:** the new/edited section contains a table (or equivalent enumerable structure) with exactly these 9 rows (or a documented superset/rename with 1:1 justification); `rg` for each of the 9 step names against the new section confirms presence. **PR-open** and **transition-request** rows are explicitly marked as mechanical targets — i.e. their "owner" column reads `runtime-owned` (or `runtime-owned (target; not yet built — Subs 2–4)`), not `γ/α/β skill-call`.
- **Negative case:** the table exists but PR-open or transition-request is classified as a cognitive-agent step (owner column names a role skill without qualifying it as "returns the decision; runtime executes"), OR the table is missing/incomplete for any of the 9 steps.
- **Surface:** dispatch doctrine — the same canonical home as AC1, or a clearly cross-referenced companion section.

### AC3 — skills are shown to control nothing

- **Invariant:** no role skill (`cdd/gamma`, `cdd/alpha`, `cdd/beta`) or dispatch-mechanics prose (`dispatch-protocol/SKILL.md`, `cds-dispatch/SKILL.md`) describes a skill *opening* a PR, *writing* a label, or *driving* control flow as its own action.
- **Oracle:** `rg -n -i "skill (opens|sets|writes|drives|dispatches)|opens? (a |the )?PR|writes? the label|drives (the )?lifecycle" src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md src/packages/cnos.cdd/skills/cdd/beta/SKILL.md src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` — every remaining hit post-edit reads as "returns the decision; the runtime executes it" (or is corrected to that framing), not as the skill performing the side effect.
- **Concrete pre-existing hit α must resolve (found during γ peer-enumeration, §2.2a of `gamma/SKILL.md`):** `beta/SKILL.md` line 38 — `merge (β's authority per CDD.md — the natural conclusion of review)` — and line 70 — `**β merge** → \`git merge\` per \`release/SKILL.md\`` — describe β (a cognitive role/skill) as the actor that *executes* the mechanical merge, not as a role that *returns a merge decision* which a runtime then executes. This is exactly the entanglement pattern named in the issue's Problem statement. **α must resolve this without violating AC4** (no runtime change): the resolution is a prose framing correction — naming that under the current bootstrap-mode architecture the β session temporarily executes the mechanical action as a stand-in for the not-yet-built runtime (Subs 2–4), while the doctrine states the target/general rule that skills return decisions and runtimes execute them. α should not silently rewrite β's actual current merge behavior (that would risk a behavior change) — only correct/qualify the *prose framing*.
- **Negative case:** any of the 5 audited files still describes a skill performing a lifecycle side effect as its own unqualified action.
- **Surface:** role-skill + dispatch docs (the 5 named files above).

### AC4 — no behavior change; gates green

- **Invariant:** doctrine/prose only — no runtime, workflow, transition-table, or test behavior change.
- **Oracle:** `git diff origin/main...cycle/584 --stat` shows **only** `.md` files touched (no `.go`, `.yml`/`.yaml`, `.json` schema, or `transitions.json`-class files); CI gates green — I1/I2/I4/I5/I6, `install-wake-golden`, dispatch guards, Go, Package, Binary; `git diff origin/main...cycle/584 -- src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go` (containing `TestSeam_CellKindNotEnforced`, confirmed present at that path) is **empty**.
- **Negative case:** any non-`.md` file changed, any gate red, or `TestSeam_CellKindNotEnforced` (or its containing file) modified.
- **Surface:** diff + CI.

---

## Source-of-truth table

Verified present in the working tree at scaffold time (`ls`/`test -f`, not assumed):

| Claim / surface | Canonical source | Status | Notes |
|---|---|---|---|
| Kernel doctrine (CCNF) | `src/packages/cnos.cdd/skills/cdd/CDD.md` | Shipped | Headers: `Kernel (CCNF)`, `Outcomes`, `Recursion modes`, `Scope-lift`, `Domain packages`, `Pointers`, `Software-specific realization → cnos.cds`, `Hard rule`, `Non-goals`. Explicitly scoped to be substrate-independent (§Non-goals: "does not name any particular tooling, platform, dispatch mechanism..."). Candidate home for AC1/AC2 if the rule is framed as kernel-level, but its own Non-goals section resists new prose — α must justify placement. |
| Cell-of-cells essay | `docs/papers/CELL-OF-CELLS.md` | Shipped | 19 numbered sections incl. `§6 Boundary: V and δ`, `§12 Relation to handoff`, `§16 Consequences for design`. Longer-form / essay register — a plausible home for the worked example, less plausible as the *stable rg-target* doctrine path (essays here read as exploratory, not binding-rule surfaces elsewhere in the repo's citation style). |
| γ role contract | `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` | Shipped | AC3 audit target. No current "skill opens/sets/writes" hits found by γ's grep pass. |
| α role contract | `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` | Shipped | AC3 audit target. No current "skill opens/sets/writes" hits found by γ's grep pass. |
| β role contract | `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` | Shipped | AC3 audit target. **Concrete hit found** — lines 38, 70 (see AC3 above): β described as executing `git merge` as its own authority/action. |
| δ role contract | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` | Shipped | **Not in the issue's AC3 audit list** — out of scope for this sub's audit per the issue body's explicit file list. γ read it in full for context (its §9 wake-invoked-mode already models "δ requests; FSM decides" — useful as an existing worked example α may cite, but α must not expand the audit scope to edit this file's prose beyond what AC1/AC2's doctrine section needs to cross-reference). |
| Dispatch-protocol mechanics | `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` | Shipped | AC3 audit target. Already contains mechanical-enforcement precedent (e.g. §"D10" activation-log write-fence, "prompt-only prohibition has been falsified; mechanical enforcement is the rule") — same spirit as this issue's governing rule, useful citation for the worked example. |
| Reference dispatch wake | `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` | Shipped | AC3 audit target. Line ~258 already distinguishes "wake writes the label directly" (claim/hard-block/release) vs. "wake requests the transition via the FSM" (β-converge) — this is a **positive existing example** of the mechanism/cognition split the issue wants named as doctrine; α should cite it, not just audit it for violations. |
| Dispatch wire-format (7-axis contract origin) | `src/packages/cnos.handoff/skills/handoff/dispatch/SKILL.md` | Shipped | Not an audit target per the issue, but γ loaded it for this scaffold's own dispatch-prompt authoring obligation. |
| AC4 test anchor | `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go` | Shipped | Confirmed present; contains `TestSeam_CellKindNotEnforced` (also referenced in `src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md` and `src/packages/cnos.cds/skills/cds/fsm/transitions.json`). Must remain byte-identical in this cycle's diff. |

**Peer enumeration (per `gamma/SKILL.md` §2.2a):** γ ran `rg -i "cells are mechanical|skills don't control anything|skill.*(opens|controls)" ` across `src/packages/cnos.cdd/`, `docs/papers/`, `src/packages/cnos.core/skills/agent/dispatch-protocol/`, `src/packages/cnos.cds/orchestrators/cds-dispatch/` before writing this scaffold. No existing doctrine section states the governing rule verbatim; the closest existing structural precedent is `delta/SKILL.md` §9.6's δ-requests/FSM-decides split and `cds-dispatch/SKILL.md`'s claim/hard-block-direct vs. β-converge-via-FSM split — both are **positive partial instances** of the rule, not yet named as general doctrine. This is additive framing, not "X does not exist" — the mechanism/cognition split is practiced in two places but not yet codified as a named, general, audited rule.

---

## α dispatch prompt

```text
You are α. Project: usurobor/cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 584 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/584
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md, src/packages/cnos.cdd/skills/cdd/issue/SKILL.md
```

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | N/A — this cycle produces Markdown/prose doctrine only; no source code in any language is authored or modified. |
| CLI integration target | N/A — no CLI surface, command, or flag is added, removed, or changed. |
| Package scoping | The diff MUST be confined to: (1) exactly one canonical doctrine home α selects and justifies from `src/packages/cnos.cdd/skills/cdd/CDD.md` or `docs/papers/CELL-OF-CELLS.md` (new section(s) for AC1+AC2); (2) targeted prose edits inside `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md`, `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md`, `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md`, `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md`, `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` (AC3 audit corrections only — no unrelated restructuring); (3) `.cdd/unreleased/584/*` artifact files. No other package paths may appear in the diff. |
| Existing-binary disposition | N/A — no binary, compiled artifact, or `cn` command is added, replaced, or deprecated in this cycle. |
| Runtime dependencies | None. |
| JSON/wire contract preservation | N/A — no JSON schema, receipt shape, `ValidationVerdict`/`BoundaryDecision` field, or `transitions.json` content is touched. Preserve as-is; if α's diff touches any `.json` file, that is itself an AC4 violation — stop and re-scope. |
| Backward-compat invariant | Doctrine/prose only — no runtime, workflow, transition-table, or test behavior change (this is AC4 verbatim, restated as the binding backward-compat constraint for this cycle). `TestSeam_CellKindNotEnforced` (`src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go`) MUST remain unmodified; CI gates I1/I2/I4/I5/I6, `install-wake-golden`, dispatch guards, Go, Package, Binary MUST stay green. |

**Design note (mode = design-and-build, doctrine-only):** α owns the design decision of *where* AC1/AC2's doctrine section lands (CDD.md vs. CELL-OF-CELLS.md vs. a new cross-referenced surface) and *how* to word the AC3 prose corrections — γ has named the oracle, the concrete pre-existing hits (β/SKILL.md lines 38/70), and the positive existing partial instances to cite (delta/SKILL.md §9.6, cds-dispatch/SKILL.md's claim/hard-block vs. β-converge split) as source material, not as a prescribed rewrite. α must record the placement decision and its justification in `self-coherence.md`.

**Known friction α inherits (see Friction notes below):** the AC3 correction to `beta/SKILL.md` risks conflating "correct the doctrine's framing" with "change what β currently does" — α must resolve this as a prose-only qualification, not a behavior change, or AC4 fails.

---

## β dispatch prompt

```text
You are β. Project: usurobor/cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 584 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/584
```

β verifies the diff against the pinned 7-axis implementation contract above (Rule 7) in addition to walking the AC1–AC4 oracle list independently. Given the "Package scoping" axis is unusually restrictive for a doctrine-only cycle (no code paths at all), β's Rule 7 check reduces to: **confirm every changed file is one of the paths named in the Package-scoping row, or `.cdd/unreleased/584/*`.** Any file outside that set is a D-severity `implementation-contract` finding regardless of whether the prose content is otherwise correct.

---

## Scope guardrails (restated verbatim from the issue's Scope section)

**In scope (doctrine + audit only):**
- Codify the governing rule in the cell-framework doctrine (`CDD.md` / `docs/papers/CELL-OF-CELLS.md` / a design-constraints surface — α picks the canonical home).
- Classify every dispatch lifecycle step as **mechanical** (runtime-owned) or **cognitive** (a `[call γ/α/β]` delegation), with the call boundary + what each skill returns.
- Audit `cdd/gamma`, `cdd/alpha`, `cdd/beta` (+ `dispatch-protocol/SKILL.md`, `cds-dispatch/SKILL.md`) prose for any language implying a skill *controls* the lifecycle (opens a PR, sets a label, drives flow); correct to "returns the decision; the runtime executes it."

**Out of scope:**
- Any runtime/workflow/transition-table change (Subs 2–4).
- Implementing checkpoint/PR-open/recovery.
- `cell_kind` enforcement; new status labels/taxonomy; Demo 0.

γ adds: no file outside the Source-of-truth table above should appear in the diff without α recording an explicit justification in `self-coherence.md` naming why the issue's file list required it (e.g. a genuinely necessary cross-reference target). `delta/SKILL.md` is explicitly **not** an audit target per the issue body — α may cite it as a positive example but should not restructure it; if α believes `delta/SKILL.md` needs a correction, that is a friction to surface via `self-coherence.md`, not a silent scope expansion.

---

## Friction notes

1. **SHA drift at branch pre-flight.** The wake-invoked input's pinned main SHA (`9277c7e3...`) had already advanced by one unrelated commit (`9309de97` — board-map regen) by the time γ ran the pre-flight check. γ branched from the verified current `origin/main` HEAD per the pre-flight rule ("base SHA known" — verified, not assumed) rather than the stale pinned value. Logged here per γ's staleness-check discipline; no `gamma-clarification.md` entry needed since this is the initial scaffold, not a mid-cycle re-pin.

2. **AC3 vs. AC4 tension on `beta/SKILL.md`.** The issue's AC3 wants control-implying language in `beta/SKILL.md` corrected ("skill... merges" → "returns the decision; runtime executes"). But `beta/SKILL.md` lines 38/70 describe β's *actual current, shipped* behavior (β does execute `git merge` today — there is no separate mechanical runtime yet; that's Subs 2–4). A literal rewrite to "the runtime executes the merge" would misdescribe what happens today and could itself be read as a behavior claim, risking AC4 ("no behavior change" — read narrowly as "no *description* of new behavior that isn't real yet"). α's resolution path (flagged in the α prompt's Design note above): frame the correction as *the general/target doctrine* stated once at the canonical doctrine home (AC1/AC2), with `beta/SKILL.md`'s existing "β merge... β's authority" prose annotated/qualified (not deleted) to distinguish "the judgment β renders (merge / no-merge)" from "the mechanical action of executing `git merge`, which β currently performs as a stand-in pending the Subs 2–4 runtime." This preserves AC4 (no behavior change — β still does exactly what it does today) while satisfying AC3 (the prose no longer reads as if β's *cognition* is what performs the mechanical act — the two are named as distinct, with the current collapse named honestly as a transitional state, not the target state).

3. **Canonical-home choice is genuinely open.** `CDD.md`'s own §Non-goals resists new prose ("does not name any particular tooling, platform, dispatch mechanism... those names belong in domain packages and runtime substrate") — but the governing rule as stated ("cells are mechanical... skills don't control anything") is arguably kernel-level (substrate-independent), not cds-specific, since it applies to the generic γ/α/β/δ contract, not just cds. `CELL-OF-CELLS.md` is essay-register and may be a better home for the worked example, while a shorter binding statement could go in `CDD.md`'s kernel sections if α argues it's substrate-independent framing rather than tooling-naming. γ deliberately leaves this design call to α (mode = design-and-build) rather than pre-deciding it — over-specifying the home would compensate for weak issue scoping per `gamma/SKILL.md` §2.4 ("do not compensate for a weak issue by making the prompt longer").

4. **No `cell_kind` / label-taxonomy scope creep.** The issue explicitly excludes `cell_kind` enforcement and new status labels from this sub. `cds-dispatch/SKILL.md`'s existing claim/hard-block-direct vs. β-converge-via-FSM split (cited above as a positive example) is adjacent to but does not require any `cell_kind` change to cite — α should quote/reference it, not extend it.

5. **γ did not expand the AC3 audit-file list.** The issue names exactly 5 files for AC3 audit (`cdd/gamma`, `cdd/alpha`, `cdd/beta`, `dispatch-protocol/SKILL.md`, `cds-dispatch/SKILL.md`). γ's own grep pass found no violating hits in `gamma/SKILL.md` or `alpha/SKILL.md` (only `beta/SKILL.md` has the concrete hit). γ did not additionally audit `delta/SKILL.md`, `operator/SKILL.md`, or `harness/SKILL.md` even though they discuss adjacent mechanics — that would be scope expansion beyond the issue's explicit list. If β's independent walk of the AC3 oracle turns up violations in files outside the 5, that is a finding to surface via `beta-review.md`, not a silent fix.
