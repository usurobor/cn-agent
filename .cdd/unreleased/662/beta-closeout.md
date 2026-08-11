# β closeout — cnos#662 (PC-D0: Cell Classes and Mechanical FSM)

**Verdict:** `converge` at **R1** (R0 was `iterate`) · **R0 findings:** 2 (F1 blocking = AC5/F4; F2 non-blocking = §9 citation)
**Full review:** `.cdd/unreleased/662/beta-review.md`

## Review summary

Independent, adversarial β review of `docs/architecture/CELL-RUNTIME-CLASSES.md` as planning matter, reconciled over **two independent β passes** under bootstrap concurrency (see `beta-review.md` "Review provenance"):

- **β-pass-A** — the separate Agent activation the δ driver spawned this session; returned `converge` but under-weighted the AC5/F4 protocol-package gap (filed it as a non-blocking "under-claim").
- **β-pass-B** — a concurrent bootstrap firing (preserved in git history at `d66d761b`); returned `iterate`, correctly flagging **AC5 = FAIL** (the spec carried zero CDR/CDW/#376/#403 protocol-package State-A truth that AC5/F4 requires) plus a minor §9 citation nit.

The stricter reading is correct, so the **reconciled R0 verdict is `iterate`**. α repaired both findings (F1: added the §11.2 Protocol-Package State-A paragraph; F2: fixed the §9 `doctrine_affecting` citation); **β converged at R1**. β re-walked the AC oracle (cell-level AC1–AC7, note-level AC1–AC8) and the ten pinned decisions D1–D10 across six dimensions, forming its view against repo ground truth **before** reading α's `self-coherence.md`. R1 convergence on the two enumerated findings is δ-attested (verbatim/mechanical fixes), recorded as a bootstrap-β limitation rather than a re-spawned third pass.

| Dimension | Verdict | Key evidence |
|---|---|---|
| A — relation-consistency (CCNF / CELL-RUNTIME / #530 / #500 / #504 / #644 / #654 / #583 / #584 / κ / shipped FSM) | ATTEST | CCNF Projection 3 match; #530 depended-on not redefined; #500 shipped / #504 open stated correctly |
| B — State-A/State-B honesty (D10) | ATTEST | `states` array, `blocked`-as-target, 12 guards, 4 request markers, shipped `cn` commands all exact-match source; `run`/`pulse` confirmed absent |
| C — no policy invented silently | ATTEST | κ≠α stated and never violated; CC↔ε verbatim; `cell_class` post-claim, readiness label-gated only |
| D — mechanical expressibility of guards/V (D8) | ATTEST | PC "no child auto-dispatched" and CC "no impl-surface modified" both reduced to checkable predicates |
| E — hard non-goals respected | ATTEST | `git diff --stat` = one spec file + own `.cdd` artifacts; no code, no child dispatch |
| F — cell FSM + wave FSM + 3 contracts + nomenclature + State A→B | ATTEST | all five Coherence-Loop terms present; three-way State partition consistent |

## Non-blocking items

Two items were considered and **rejected as non-findings** (`beta-review.md §R0.7`): a stale cli-wrapper doc-comment in `cmd_issues_fsm.go` (repo source outside this cell's forbidden surface; spec correct against the real package); and a cosmetic YAML-nesting quote of the γ-scaffold (D9 makes exact field placement non-load-bearing). β-pass-A's third "non-finding" (a benign CDR under-claim) was **too lenient and is superseded** by blocking finding F1 (AC5/F4) — recorded honestly rather than left as a rejected nit.

## Contradiction-rule outcome

No TRUE contradiction in the pinned architecture. No candidate operator-gate hold raised by β.

## β independence disclosure (#664)

β ran as a **separate Agent activation** spawned by the δ/dispatch driver, under the **same account/session lineage** as the authoring cell. Therefore: **α≠β holds at the protocol layer** (β did not author the matter, re-verified all State-A claims directly against source, and formed its view before reading `self-coherence.md`) but is **bootstrap-limited at the hosting layer** (same session lineage — the exact structural gap #664 tracks). This is recorded honestly rather than presented as full independence. The later CC ratification in the exit sequence provides the additional external warrant.

**β → γ handoff:** converge at R1; the cell may close and request `status:review`.

## R2 — fresh independent β review (post operator-final-read repair)

**Verdict:** `converge` at **R2** · **fresh independent pass** (NOT δ-attestation), corroborated by a second independent pass · **findings:** 2 (both non-blocking) · **operator-gate holds:** none.

After κ's operator-final-read of PR #667 returned **ITERATE NARROWLY** with six typed-contract-surface blockers, α repaired blockers 1–5 and a **fresh, separate β activation** reviewed the *entire revised contract surface* (blocker 6). This was a real independent read — δ-attestation was explicitly ruled out for substantive contract/schema corrections. Full record: `beta-review.md §R2`; a second corroborating independent pass: `beta-review-R2.md`.

| Blocker | Disposition | Evidence |
|---|---|---|
| 1 — canonical `cn.cell.contract.v1` (§2) | RESOLVED | Field-by-field diff: schema block and worked instance are one shape; `cell.id`, `requested_output` object, `non_goals` under `constraints`, top-level `gates`/`stop_conditions`. |
| 2 — intent reconciliation (§2↔§13) | RESOLVED | `intent_ref → cn.intent.v1` in both sections; issue = carrier; no residual `intent:{source:issue}`. |
| 3 — PC result tagged union (§3.2) | RESOLVED | D0 → `{mode:d0, artifact_ref, ready_for_coherence_review}`; Wave → `{mode:wave, wave_ref, graph, ready_for_wave_review}`; D0 not forced to carry `wave_ref`/`graph`. |
| 4 — actor-collapse terminology (§8/§14/γ decl.) | RESOLVED | Two-mode distinction; #662 = hosting-identity collapse (#664), not actor collapse; κ≠α preserved. |
| 5 — stale §17/§17 Q6 refs | RESOLVED | §10→§16 Q1, §11.5→§16 Q3; both targets exist and match; no §17/§18 remains. |
| 6 — fresh independent β R2 | RESOLVED | This pass (+ a second corroborating pass). |

**No regression** in the operator-named "already good" parts; **every class-specific V predicate remains mechanically expressible** (§5, §12); **State-A re-verified against source** (`transitions.json` states array + `blocked`-as-target + guard vocabulary; `cn cell pulse`/`run` confirmed absent; `TestSeam_CellKindNotEnforced` observation-only).

**Non-blocking findings:** F-R2.1 (receipt hygiene — a couple of frozen/superseded historical receipt lines still say "actor-collapse"; the normative spec §8 and the operator-named γ declaration are correct; `alpha-closeout.md` line cleaned up in R2); F-R2.2 (process — corroborates κ's separately-filed P1 substrate race; a deliberate review quorum, not accidental duplicate activation, is the right model).

**β independence disclosure (#664):** both R2 passes ran as **separate Agent activations** but under the **same Sigma account/model lineage** — **α≠β holds at the protocol layer** (neither authored the matter; both re-verified State-A against source; both formed their view before reading `self-coherence.md`), **bootstrap-limited at the hosting layer** (hosting-identity collapse, #664). The exit sequence's separate, non-Sigma CC ratification supplies the hosting-independent warrant.

**β → γ handoff (R2):** converge; the cell may re-close and re-request `status:review` for a second operator-final-read. **Do NOT merge; do NOT mark PR #667 ready; do NOT dispatch CC ratification** — those are operator-gated per the exit sequence.
