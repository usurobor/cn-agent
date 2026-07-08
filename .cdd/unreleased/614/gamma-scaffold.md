## gamma-scaffold.md — cycle/614

**cell_kind:** `doctrine`

**Issue:** #614 — cdd: amend CELL-KINDS (#570) — add `planning` cell kind + mandatory terminal learning / ε-observations section.

**Mode:** design-and-build (doctrine). **Doctrine-only; no behavior change** (AC5). No CUE schema field changes, no FSM guard changes, no verifier enforcement changes.

## Source-of-truth table

| Surface | Path | What changes |
|---|---|---|
| Cell-kind taxonomy | `src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md` | Add `planning` as an 11th cell kind (AC1); add the mandatory terminal `learning` section as doctrine, applicable to all cell kinds (AC2); update the FSM-awareness table row for `planning`; update the intro count ("10 cell kinds" → "11 cell kinds"). |
| γ role skill | `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` | State explicitly that γ binds ε-observations/learning into the receipt at closeout (AC3) — frontmatter `description`, Core Principle, and the `gamma-closeout.md` authoring bullet in §2.7. |
| Kernel doc | `src/packages/cnos.cdd/skills/cdd/CDD.md` | Update the `CELL-KINDS.md` pointer (line ~94) to include `planning` in the enumerated cell-kind list and to name the mandatory learning-section rule, so the kernel doc and the taxonomy doc do not drift (AC4). |
| Receipt schema surface | `schemas/cdd/receipt.cue` | Add a **documentation-only** cross-reference comment naming `CELL-KINDS.md`'s mandatory `learning` section as a future schema candidate (not yet a typed field) — satisfies "link from ... the receipt schema surface" without adding a schema field (preserves AC5 no-behavior-change; verifier enforcement stays explicitly deferred per the issue's Scope). |

## Per-AC oracle (mechanical, docs-only)

- **AC1** — `grep -n "^### 11\. \`planning\`" CELL-KINDS.md` finds the new section; it contains Purpose / Matter / Lifecycle (α/β/γ) / Role-vs-model note / Relationship-to-neighbors / FSM implication sub-bullets (mirror the existing 10 entries' shape).
- **AC2** — `grep -n "learning:" CELL-KINDS.md` finds a new "Mandatory terminal learning section" doctrine block with the exact YAML field schema from the issue body (`observations`, `process_deltas`, `reusable_patterns`, `followups`, `operator_burden`), framed as ε-captured-by-γ, stated as applying to **all** cell kinds (not just `planning`).
- **AC3** — `grep -n "binds" gamma/SKILL.md` finds explicit language (frontmatter description + Core Principle + §2.7 closeout bullet) stating γ binds the ε-observations/learning section into the receipt at closeout.
- **AC4** — `CDD.md`'s `CELL-KINDS.md` pointer line lists `planning` alongside the other 10 kinds; no new orphan doctrine (every new CELL-KINDS.md section is reachable from CDD.md's pointer and CELL-KINDS.md's own cross-references section).
- **AC5** — doctrine-only: `git diff --stat cycle/614 main` touches only `.md`/`.cue`-comment files (no `.go`, no `transitions.json`, no CI workflow, no label changes); CI gates green (I1, I2, I4, I5, I6, install-wake-golden, dispatch-repair-preflight, dispatch-closeout-integrity, Go, Package, Binary — whichever of these run in this repo's CI on a docs-only PR).

## Scope guardrails

- **In scope** (per issue body): `CELL-KINDS.md` planning section + learning-section doctrine + γ role-text clarification; `CDD.md` linkage; a `planning` worked example (a short illustrative example inside the new `### 11. \`planning\`` section, following the pattern of existing entries — no separate example file needed).
- **Out of scope** (per issue body, STOP conditions): no mechanical/verifier enforcement of the learning section; no `cell_kind` FSM guard changes (`table.go` / `transitions.json` untouched); no dispatch-label changes; no Demo 0 work. If any implementation step drifts toward these, STOP and escalate rather than improvise past the issue's named boundary.
- **Cell-kind self-reference note:** this cycle (#614) is itself a `doctrine` cell per the existing (pre-amendment) taxonomy — not a `planning` cell — since it produces doctrine text directly rather than an issue/wave contract. Do not conflate: `planning` cells produce plan/issue/wave *contracts* as matter; this cycle produces the doctrine describing that cell kind.

## α prompt

Implement the four source-of-truth changes above in `.cdd/unreleased/614/` cycle branch `cycle/614`, against `main@ff7bd2fe467580f725d05731ff4a132ae46ebcdb`:

1. **`CELL-KINDS.md`**:
   - Bump the "10 cell kinds" heading/intro references to "11 cell kinds" wherever they enumerate the count.
   - Insert a new `### 11. \`planning\`` section after `### 10. \`experiment\`` (before the "FSM awareness" table), following the issue body's Amendment 1 content: Purpose, Matter, Lifecycle (α produces the plan/issue/wave; β reviews for coherence/scope/proof/design boundary; γ closes: posts/links issues, verifies labels, records sequence, captures decisions, binds learning/follow-ups into the receipt), the Role-vs-model note (κ commonly acts as α of a planning cell; the Greek letter names the role, not the model/agent — do not mint a new Greek identity per model), Relationship to neighbors (`issue_authoring` produces one issue; `wave` composes child-cell receipts; `planning` produces the plan/wave contract + sequencing and may spawn `issue_authoring` + `implementation` children — β resolves exact boundaries), FSM implication (may create/update issues; does not imply immediate dispatch unless the plan explicitly says so — mirrors `issue_authoring`). Include one short worked-example line/paragraph (e.g., citing this very issue #614's own planning-shaped origin — "operator review surfaced two gaps in #570; a planning-shaped pass would have filed this as its own issue with explicit AC/scope, then handed to a `doctrine` implementation cell — illustrating the planning → issue_authoring/implementation handoff" — or an equally concrete illustrative sentence; do not invent a fictitious unrelated example when a real one is available).
   - Add the new `planning` row to the "FSM awareness" table (mirrors the other 10 rows' shape: cell kind | valid matter | allowed transition request).
   - Add a new doctrine subsection (top-level, e.g. after "Recording point (AC8)" or as its own numbered section) titled something like "Mandatory terminal learning section" stating: every terminal cell closeout MUST include the `learning` section with the exact YAML schema from the issue body (`observations`, `process_deltas`, `reusable_patterns`, `followups`, `operator_burden`), framed explicitly as "ε captured by γ" (ε observes patterns across runs; γ binds the observation into the receipt at closeout), applying to **all** cell kinds, with an explicit "required by doctrine; mechanical verifier enforcement deferred to a follow-up" caveat (mirroring #570's own "document first, enforce later" stance already cited in this file's "Observation, not enforcement" section).
   - Update the "Cross-references" section if needed so the new sections are reachable (no orphan doctrine — AC4).

2. **`gamma/SKILL.md`**:
   - Frontmatter `description:` — add a clause naming that γ binds ε-observations/learning into the receipt at closeout (keep the existing description's content; append, don't replace).
   - `## Core Principle` — add one sentence stating γ's closeout responsibility explicitly includes binding the cell's learning/ε-observations into the receipt (this is a clarification of existing responsibility, not new scope — the "cycle closure without learning" failure mode is already named in the paragraph below Core Principle; make the positive obligation explicit).
   - §2.7 (the `gamma-closeout.md` authoring bullet, ~line 337) — add that `gamma-closeout.md` must include the `learning`/`epsilon_observations` section per `CELL-KINDS.md`'s mandatory-learning-section doctrine, with the field schema.

3. **`CDD.md`**:
   - Update the `CELL-KINDS.md` pointer line (~line 94) to add `planning` to the parenthetical enumeration of cell kinds, and add a short clause noting the mandatory terminal learning-section rule (cite `CELL-KINDS.md` for the field schema — don't restate the schema in CDD.md; CDD.md stays a pointer, CELL-KINDS.md stays the source of truth).

4. **`schemas/cdd/receipt.cue`**:
   - Add a documentation-only comment (no field/type change) near the `#ProtocolGapRef` / ε-signal comment block, cross-referencing `CELL-KINDS.md`'s mandatory `learning` section as a doctrine-level requirement not yet reflected as a typed schema field — explicitly noting verifier enforcement is deferred (consistent with the issue's Scope/STOP conditions). This satisfies "link from ... the receipt schema surface" (AC4/AC1 context) without adding behavior (AC5).

Write `.cdd/unreleased/614/self-coherence.md` with a `§R0` section per `alpha/SKILL.md` §2.5 conventions: per-AC verification notes (AC1–AC5), a diff summary, and the review-ready signal. Do not touch any `.go`, `transitions.json`, CI workflow, or label-related file — this is docs-only per AC5.

## β prompt

Review cycle/614 against the 5 ACs above independently, per `beta/SKILL.md`'s review-oracle discipline (`doctrine`-kind review surface per `CELL-KINDS.md` entry 9: doctrinal coherence, cross-reference correctness, no orphaned new doctrine, no silent behavior change smuggled in under a docs label). Specifically check:

- AC1: the `planning` section mirrors the shape (all required sub-fields) of the other 10 entries; no missing sub-bullet.
- AC2: the learning-section schema in `CELL-KINDS.md` matches the issue body's YAML **exactly** (field names: `observations`, `process_deltas`, `reusable_patterns`, `followups`, `operator_burden`); framed as ε-captured-by-γ; states it applies to all cell kinds, not just `planning`.
- AC3: `gamma/SKILL.md`'s new language is unambiguous that γ (not α, not β) binds learning into the receipt.
- AC4: `CDD.md`'s pointer line is updated; nothing new in `CELL-KINDS.md` is unreachable from either `CDD.md` or `CELL-KINDS.md`'s own cross-references.
- AC5: confirm via `git diff --stat` that no `.go`/`transitions.json`/workflow/label files changed; confirm the `receipt.cue` change is comment-only (no field/type added); confirm relevant CI gates pass.

Write `.cdd/unreleased/614/beta-review.md` with a `§R0` verdict (`converge` or `iterate` + findings).

## Friction notes

- Issue explicitly names this as amending a **closed** issue (#570) rather than reopening it — γ scaffold preserves that framing; no attempt to touch #570's own receipt/closeout artifacts.
- The "receipt schema surface" AC4 clause is satisfied via a comment-only doc note in `receipt.cue` rather than a new typed field, since the issue's own Scope/STOP conditions explicitly defer "mechanical verifier enforcement" and forbid the learning section becoming "a mandatory mechanical gate in this cell." A typed CUE field would be a schema-behavior change (fails `cue vet` fixtures or requires new fixture updates) and risks tripping AC5's no-behavior-change bar; a comment is the doctrine-only-safe choice.
- This cycle produces doctrine that (once merged) makes *future* terminal closeouts require a `learning` section — but this cycle's own `gamma-closeout.md` is authored under the same doctrine being landed, so γ's closeout for #614 should itself carry a `learning` section as a dogfooding demonstration (not mechanically required by any gate yet, but consistent with the spirit of the change).
