# α closeout — cnos#662 (PC-D0: Cell Classes and Mechanical FSM)

**Cell class:** planning (PC), PC-D0 mode · **matter_domain:** doctrine · **doctrine_affecting:** true
**Run class:** `first_pass` (R0 → R1; δ-resumed from checkpointed matter, see below)
**Matter produced:** exactly one normative planning artifact — `docs/architecture/CELL-RUNTIME-CLASSES.md` (Status: **Draft**).

## What α produced

`docs/architecture/CELL-RUNTIME-CLASSES.md` — the realization-layer spec that operationalizes the WC/PC/CC output-telos classes named by `docs/architecture/CELL-RUNTIME.md` (#628) into a cell contract envelope, class-specific `V` predicates, the CC↔ε lineage, the cell FSM (grounded in the shipped `transitions.json`) and the specified-but-unshipped wave FSM, the human-gate and wake-topology policy, and the schema-first destination — carrying the ten operator-pinned decisions D1–D10 as settled input, not re-deriving them.

α owns and wrote this specification matter. **κ did not author it** (κ≠α, D2/F3): κ authored the #662 issue and its operator comments (the control-plane input); α authored the spec file (the cell matter). The bootstrap **hosting-identity collapse** (one Sigma lineage across separate κ/α/β activations, tracked by #664 — *not* actor collapse, per R2 blocker 4) is declared explicitly, not silently assumed — see the γ receipt's bootstrap-calibration section and spec §8.

## AC + decision coverage

Full per-AC walk (cell-level AC1–AC7, note-level AC1–AC8) and per-decision walk (D1–D10) is in `self-coherence.md §R0` (with §R1 repair note), with file/section citations. β re-walked the same oracle across two independent passes; the reconciled **R0 verdict was `iterate`** (finding F1 = AC5/F4 protocol-package State-A truth absent; F2 = a §9 citation nit), and **β converged at R1** after α's repair (`beta-review.md`).

## R1 repair (against β R0 findings)

- **F1 (AC5/F4):** added the "Protocol Package state truth (F4)" paragraph to §11.2 — `cnos.cds` shipped (#403), `cnos.cdr` shipped v0.1 (#376), `cnos.cdw` illustrative-only; every worked example uses `protocol: cds`. The R0 spec had omitted the protocol-package half of AC5 (only the `cn cell`/#500/#504 half was present); the omission was caught by a concurrent independent β pass (preserved at `d66d761b`) and is now closed.
- **F2 (§9):** re-cited `doctrine_affecting: true` to the operator-authorization comment's contract block (§2's snippet quotes only `matter_domain: doctrine`).

Both repairs are verbatim/mechanical fixes of the enumerated findings; no other AC changed.

## δ-resume adjustments (R0.1)

The R0 spec + self-coherence were committed by an earlier firing of this bootstrap session (commit `1754cf4d`) before β ran. On δ-resume (`delta/SKILL.md` §9.11), three binding-directive/honesty alignments were applied before independent β and recorded in `self-coherence.md §R0.1`:

1. **Status → `Draft`** (κ dispatch contract fixes the artifact Status at Draft for this bootstrap pre-ratification cycle; R0 had mirrored the sibling's "Proposed" header). The "Not ratified — separate CC review → operator-final-read → doctrine cell" exit framing is preserved.
2. **State-A command-surface basis → honest** (§11.2/§16): the ground truth is the shipped command *definitions* in `src/go/internal/cli` + `src/go/internal/cell`, not a `cn --help` binary run that did not occur in this environment. The shipped/unshipped partition is unchanged and was independently re-verified by β against source.
3. **Dangling `§11.6` cross-references repointed** (§11.1, §11.2×2) to the correct subsections.

No AC-level claim changed; these were directive-alignment and internal-consistency fixes on the same matter.

## Hard-STOP check

No true contradiction was found among D1–D10, the issue draft, and shipped state. The single nuance requiring exposition (shipped `states` array omits `blocked` while `blocked` remains a reachable `target_state`) is resolved by exposition in §11.1 against shipped ground truth, not by invention. **No operator-gate hold raised by α.**

## Non-goals honored

No Go/runtime/schema/FSM code; no wake changes; no child issues filed or dispatched; no PC-Wave; no child `status:todo` labels; no separate PC/CC provider; no CCNF role-semantics change. The only filesystem writes are the spec file and this cell's own `.cdd/unreleased/662/` artifacts (confirmed by β via `git diff --stat`).

## R2 repair (against κ's operator-final-read verdict on PR #667 — ITERATE NARROWLY)

κ's operator-final-read (issue comment 2026-07-16, materialized verbatim as `operator-review.md`, `schema: cn.operator-review.v1`) returned **iterate narrowly**: the architecture converges but the typed-contract surface carried six load-bearing blockers. α repaired blockers 1–5 (blocker 6 = a fresh independent β R2, spawned separately). Full per-blocker walk in `self-coherence.md §R2`; summary:

1. **Canonicalized `cn.cell.contract.v1` (§2)** — one governing envelope; the worked instance (#662's own contract) rewritten to validate against it verbatim (`cell.id`; `requested_output` object; `non_goals` under `constraints`; `gates`/`stop_conditions` top-level).
2. **Reconciled intent (§2 ↔ §13)** — envelope now carries `intent_ref → cn.intent.v1`; the GitHub issue is a **carrier/projection**, not intent identity; the inline `intent: { source: issue }` is removed.
3. **Split PC result by mode (§3.2)** — tagged union: PC-D0 `{ mode: d0, artifact_ref, readiness: ready_for_coherence_review, … }` vs PC-Wave `{ mode: wave, wave_ref, graph, readiness: ready_for_wave_review, … }`. D0 no longer forced to carry `wave_ref`/`graph`.
4. **Corrected actor-collapse terminology (§8, §14)** — distinguished **actor collapse** (one activation, two roles in one cell boundary) from **hosting-identity collapse** (separate activations, one Sigma lineage; #664). #662 was **hosting-identity collapse, not actor collapse**; **protocol-level κ≠α held — κ did not author the spec.** The γ bootstrap declaration is updated to match.
5. **Repaired stale §17 / §17 Q6 references** — §10 → §16 Q1 (wake-provider realization); §11.5 → §16 Q3 (command-surface sequencing vs #504). Both targets exist and match; pure renumber, no new policy.

**R2 Hard-STOP check:** no blocker required an unpinned architecture decision; **no operator-gate hold raised by α in R2.** Every operator-named "already good" part is preserved. The only spec write is `CELL-RUNTIME-CLASSES.md`; the only other writes are this cell's own `.cdd/unreleased/662/` receipts (incl. the new `operator-review.md`). No code, schema, FSM, wake, or child issue.

**Terminology note (R2):** the "What α produced" line above and the γ bootstrap declaration now use **hosting-identity collapse** (#664), not "actor collapse," per blocker 4; κ≠α held at the protocol layer. The frozen R0 scaffold text (`gamma-scaffold.md:72`) is left as historical R0 record.

**α → β handoff (R2):** review-ready for a **fresh independent β R2** over the entire revised contract surface. (R0 iterate → R1 converge → R2 repair per κ's six-blocker contract.) β R2 returned **converge** — two fresh separate activations (both protocol-independent, bootstrap-limited under one Sigma lineage, #664) independently converged, all six blockers RESOLVED, no operator-gate holds.

## R3 repair (§8 authorized-transitional reframe — operator-settled)

Per the operator's **direct control-plane instruction** settling the κ/α doctrine (*"You're K in sigma world executing the role of alpha of cnos planning cell since we haven't yet made planning cell operational … We're bootstrapping that three-cells agent."*), §8 of `CELL-RUNTIME-CLASSES.md` — and its mirror statements in the Thesis, §1 table, §14 reconciliation rows, §15 non-goal, and §16 authoring note — were rewritten from **"disclosed nonconformance / not a valid topology"** to the **operator-authorized transitional bootstrap posture (State A)** framing: κ=α during bootstrap is a legitimate, explicitly transitional, operator-authorized stage retired at State B (not a nonconformance, not a licensed permanent topology), still carrying a reduced *independent* warrant (hosting-identity, #664) discharged by external β; and `CELL-KINDS.md:137` is reframed from "stale/wrong" to **"State-A-accurate, scoped, superseded at State B"** (the downstream migration repoints/scopes it, not "corrects an error"). All invariant claims held unchanged — unconditional State-B κ≠α, actor-collapse-forbidden vs hosting-identity-collapse (#664) distinction, external-β-bound-to-matter-SHA provenance (§11.6), no `CELL-KINDS.md` edit (§14/§15), §12 guard. Full walk in `self-coherence.md §R3`. This is a mechanical reframe of the characterization surface only — **no unpinned decision, no operator-gate hold raised by α.**

**Matter produced (R3 reframe):** spec-only amend commit **SHA_M = `014c75ff05570be877026b36e7f74c0abfa68d0e`** (`docs/architecture/CELL-RUNTIME-CLASSES.md` only; trailer-less cnos-cell commit; force-pushed to `cycle/662`). This is the immutable matter SHA offered for exact-SHA external-β re-review (§11.6). γ closeout is **not** written here — γ binds SHA_M + the external review later, after external β converges. `CELL-KINDS.md` not edited.
