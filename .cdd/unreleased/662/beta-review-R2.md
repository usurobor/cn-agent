# β review §R2 — fresh independent pass (post operator-final-read repair) — cnos#662

**Artifact under review:** `docs/architecture/CELL-RUNTIME-CLASSES.md` (revised, post α R2 repair)
**Repair contract:** `operator-review.md` — κ operator-final-read on PR #667, ITERATE NARROWLY, six blockers
**Verdict:** converge (see final line)

---

## Review provenance (honest hosting-identity disclosure — read first)

This is a **fresh, separate Agent activation**, spawned to discharge **blocker 6** of κ's
operator-final-read. I authored **none** of the matter I review — not the spec, not α's
`self-coherence.md`, not the prior `beta-review.md`, not the γ receipts. I built my entire view from the
revised spec and from repo ground truth (`transitions.json`, `src/go/internal/cli`, the FSM test file,
`CELL-KINDS.md`) **before** reading α's `self-coherence.md`, which I opened only afterward to cross-check
honesty — not to adopt its conclusions.

**Hosting-identity status: bootstrap-limited, NOT hosting-identity-independent (#664).** I am a distinct
activation but run under the **same account / model / session hosting lineage (Sigma)** as κ/α/the prior
β. So **protocol-level independence holds** (no matter authored by me; every State-A claim re-verified
directly against shipped source; my view formed before reading `self-coherence.md`) while
**hosting-identity separation does not** — exactly the structural limitation #664 tracks. δ-attestation was
**not** used here; this is a genuine independent read, but the reader is not hosting-independent. The
exit sequence's **external, non-Sigma CC ratification** remains the stronger warrant, as κ requires.

This review is written to a **new file** and does not touch `beta-review.md`. I note that `beta-review.md`
already carries a §R2 appended by another concurrent activation; I did not rely on it and reached my
dispositions independently. Its existence and the mid-review mutation I did not personally re-observe are
consistent with the substrate race κ filed as P1.

---

## Per-blocker disposition (1–6)

| # | Blocker | Disposition | Evidence (verified from source, this pass) |
|---|---|---|---|
| 1 | Canonical `cn.cell.contract.v1` (§2) — one exact shape; worked instance validates verbatim | **RESOLVED** | Field-by-field diff below. All four originally-conflicting axes fixed. Schema block (ln 34–51) and worked instance (ln 57–74) are one object, key-for-key and nesting-for-nesting. |
| 2 | Intent reconciliation (§2 ↔ §13) via `intent_ref` | **RESOLVED** | §2 (ln 38, 53, 61) and §13 (ln 347–348) both carry `intent_ref: { schema: cn.intent.v1, id, carrier: { kind: github_issue, ref } }`; issue framed as carrier/projection. Grep for residual issue-is-intent framing (`intent.*source.*issue`, `source: *issue`) returns **nothing**; the only `source:` on intent is `cn.intent.v1`'s own `source: operator` (correct — κ produces intent). No `intent: { source: issue }` survives; §2 and §13 agree in both directions. |
| 3 | PC result tagged union by `mode` (§3.2) | **RESOLVED** | §3.2 (ln 107–114): PC-D0 → `{ class: planning, mode: d0, artifact_ref, readiness: ready_for_coherence_review, requires_operator_gate: true }`; PC-Wave → `{ class: planning, mode: wave, wave_ref, graph: { nodes, edges }, readiness: ready_for_wave_review, requires_operator_gate: true }`. D0 is **not** forced to carry `wave_ref`/`graph`; prose (ln 105) states forcing them would "turn every Planning Cell into a wave-producer." |
| 4 | Actor-collapse terminology (§8, §14, γ bootstrap declaration) | **RESOLVED** | Spec §8 (ln 212–219) draws the two-mode distinction and states "**#662 was primarily hosting-identity collapse, not actor collapse**" with "protocol-level κ≠α held — κ did not author the spec, α did." §14 κ-boundary row (ln 370) matches. **γ bootstrap declaration** (`gamma-closeout.md` item 4) now reads "**Hosting-identity collapse declaration (corrected terminology — R2) … not actor collapse**," with an explicit note that the R0/R1 "actor-collapse" label is corrected per blocker 4. All three blocker-4 targets (§8, §14, γ declaration) are corrected. Non-blocking residuals in *historical/superseded* receipts — see Finding F-R2.1. |
| 5 | Stale §17 / §17 Q6 references | **RESOLVED** | No `§17`/`§18` reference remains in the document. §10 (ln 237) points to **§16 Q1** (Wake-provider realization — matches ln 385). §11.5 (ln 321) points to **§16 Q3** (Sequencing of illustrative command surface vs #504 — matches ln 387). Both targets exist and their content matches the referring text. The surviving `§19` (ln 383) and `§20`/`§1–§20 of #662` (ln 393) are the **legitimate references to the embedded draft inside issue #662**, correctly not flagged per contract. |
| 6 | Fresh independent β R2 exists | **RESOLVED** | This review. Fresh separate activation; not δ-attestation; view formed from source before `self-coherence.md`. Hosting-identity-limited per #664 (disclosed above). |

---

## Blocker 1 — field-by-field one-shape confirmation (schema template vs worked instance)

I diffed the two YAML blocks myself, path by path, without trusting the closing prose:

| Path | Schema template (ln 34–51) | Worked instance (ln 57–74) | One shape? |
|---|---|---|---|
| `schema` | `cn.cell.contract.v1` | `cn.cell.contract.v1` | ✓ |
| `cell.{id,class,mode,protocol,matter_domain}` | key is `id` | `{ id: 662, class: planning, mode: d0, protocol: cds, matter_domain: doctrine }` | ✓ — `id`, **not** `issue` |
| `scope.{repo,wave,parent_cell}` | present | `{ repo: usurobor/cnos, wave: 627, parent_cell: null }` | ✓ |
| `intent_ref.{schema,id,carrier.{kind,ref}}` | present | `{ schema: cn.intent.v1, id: intent-2026-0711-662, carrier: { kind: github_issue, ref: cnos#662 } }` | ✓ |
| `inputs.required` / `inputs.optional` | present | `required: [...]`, `optional: [ prior_receipts ]` | ✓ |
| `requested_output` | object `{ kind, path }` | `{ kind: artifact, path: docs/architecture/CELL-RUNTIME-CLASSES.md }` | ✓ — **object**, not list |
| `acceptance.predicates` | list | list of five predicates | ✓ |
| `constraints.{allowed_paths,forbidden_paths,non_goals}` | `non_goals` under `constraints` | `non_goals` under `constraints` | ✓ — **not** `cell.non_goals` |
| `gates.{operator_authorization_required,operator_acceptance_required}` | **top-level** | top-level `{ ...: true, ...: true }` | ✓ — **not** `cell.gates` |
| `stop_conditions` | **top-level** list | top-level list | ✓ |

No key present in one and absent in the other; no divergent nesting on any of the four originally
conflicting axes (`cell.id` vs `cell.issue`; `requested_output` object vs list; `constraints.non_goals`
vs `cell.non_goals`; top-level `gates` vs `cell.gates`). **The schema block and the worked instance are
one shape.** Blocker 1 confirmed independently, on the source, not on the prose's word.

---

## Regression check — the "already good" parts (κ's do-not-lose list)

ATTEST — none regressed by the R2 repair:
- **WC/PC/CC as telos classes of one CCNF kernel** — §1 table + D1 (ln 26), §3 intact.
- **CC↔ε reconciled, §4 carried verbatim** — fenced block (ln 138–155) present and unaltered; matches CCNF Scope-Lift Projection 3; #530 still depended-on, not redefined.
- **κ outside the cell** — §8 intact; the R2 terminology repair *strengthens* it (κ≠α restated in both collapse-mode paragraphs) without weakening the boundary.
- **State-A / specified / illustrative-future partition (§11)** — three-way partition intact (11.1–11.2 shipped / 11.3–11.4 specified / 11.5 illustrative), including the F4 Protocol-Package State-A paragraph (ln 282).
- **Shipped FSM + request-marker mechanics grounded** — §11.1 re-verified against `transitions.json` (below).
- **CC owns wave judgment / FSM owns wave transition (§11.4, D7)** — intact; typed `wave_transition_request` present.
- **Wave-level authorization (§9, operator ≠ child scheduler)** — intact ("operator must not become the scheduler for every child").
- **Mechanical class-specific V (§5)** — table intact (below).
- **Bootstrap limitation disclosed** — §10, §11.1, §16 authoring note intact.

---

## Mechanical-expressibility check (§5, §12)

ATTEST — every class-specific guard/V predicate remains a checkable predicate over receipt fields /
artifact paths / label events; none narrative-only:
- PC "no child auto-dispatched" → "the PC applied no `status:todo` to any child" (§5 / §12, checkable over label events).
- CC "no implementation surface modified" → "matter paths ⊆ judgment artifacts, no code/product diff" (§5 / §12, checkable over artifact paths / branch diff).
- Common floor (receipt complete · evidence present · role ownership valid · contract fields satisfied · no unresolved blocking finding) + WC additions (executable ACs pass · tests/checks pass · allowed-surface/non-goal guards) are all field/label/diff predicates.
- `checks_passing` is a real shipped guard (`transitions.json` `guards`). §12 ties each collapse mode to a `V` predicate or a CCNF firebreak citation. No V predicate is narrative-only.

---

## State-A spot-check (re-verified directly against shipped source, this pass)

ATTEST — every spot-checked State-A claim matches shipped source:
- **Declared states array** — spec §11.1 `["ready","todo","in-progress","review","changes"]` **exact-match** `transitions.json:18`.
- **`blocked` is a `target_state`, not an enum member** — confirmed: `transitions.json:156` `"target_state": "blocked"` inside an in-progress rule (`propose_status_blocked`, cnos#575 AC2); `blocked` **absent** from the `states` array. Spec's nuance is correct.
- **Guard vocabulary** — all 12 names in spec §11.1 exact-match the `transitions.json` `guards` block (`run_active, branch_exists, branch_has_commits, pr_exists, pr_has_commits, review_request_present, repair_contract_present, cdd_artifacts_present, checks_passing, claim_request_present, block_request_present, release_request_present`).
- **Request-marker table** — `CLAIM/REVIEW/BLOCK/RELEASE-REQUEST.yml` all present in `transitions.json`; the #574 PR-commits tightening (`review` rule uses `all_true: [review_request_present, pr_exists, pr_has_commits]`, ln 137) and the #368 delta-recovery / no-blind-requeue behavior match the guard docs.
- **Command surface** — `cn cell return/resume/finalize` present in `src/go/internal/cli/cmd_cell.go` (CommandSpecs `cell-return`/`cell-resume`/`cell-finalize`, help-test registry asserts `"cell return"`, `"cell finalize"`); `cn issues dispatch` and `cn issues fsm evaluate|scan` present. **`run`, `pulse`, `measure`, `bundle`, `act` genuinely ABSENT** as `cn cell` subcommands — grep across `src/go` for `cell pulse|cell run|cell measure|cell bundle|cell act` returns nothing. `cn cell pulse`/`cn cell run` are correctly labelled illustrative-future in §11.5.
- **CellKind seam** — `TestSeam_CellKindNotEnforced` exists (`issuesfsm_test.go:810`); `FactSnapshot.CellKind{Observed, Source, DefaultedTo}` is observation-only (`CELL-KINDS.md:222` states no `table.go` transition rule consumes it), exactly as §6 / §11.3 claim.

---

## Internal-consistency check

ATTEST — no dangling pointer introduced by the repair. Every internal `§N` reference (§1–§13, §16)
resolves to a real section; the document ends at **§16** (five open questions). §17/§18 absent. The lone
`§19` (ln 383) and `§20` (ln 393) are the legitimate embedded-draft-of-#662 historical references, not
sections of this document. §10→§16 Q1 and §11.5→§16 Q3 both point at existing questions whose content
matches the referring text.

---

## Honesty cross-check against α's self-account (`self-coherence.md` §R2)

Read only after forming my view. α's §R2 blocker-by-blocker account is **accurate against current
ground truth**: it correctly describes the R1 defects (inlined `intent: { source: issue }`, forced
`wave_ref`+`graph`, `§17` dangling pointers, "actor collapse" mislabel) and the R2 repairs, all of which
I independently confirmed present. α's §R2 note honestly *withdraws* its earlier §R0 AC4 phrasing rather
than papering over it. One claim to watch — "the γ bootstrap declaration … updated to match" (self-
coherence ln 71) — holds against the current tree: `gamma-closeout.md` item 4 does now carry the
corrected "hosting-identity collapse" framing. No surviving honesty gap.

γ closeout independently verified for blocker 4: `gamma-closeout.md` item 4 is the operator-named
bootstrap declaration and it correctly states hosting-identity collapse, protocol-level κ≠α held, and
declares the shared hosting identity explicitly (#664) rather than silently assuming equivalence.

---

## Findings

- **F-R2.1 (non-blocking · receipt hygiene · for γ/κ awareness, not a spec defect).** Three *historical /
  superseded* receipt lines still apply the word "actor-collapse" to the Sigma-reuse: `alpha-closeout.md:11`
  ("Bootstrap actor-collapse is declared") — **explicitly superseded** by `alpha-closeout.md:54`;
  `self-coherence.md:12` (frozen R0 AC4 walk) — **explicitly corrected** by its own §R2 note at ln 76; and
  `gamma-scaffold.md:72` (frozen R0 scaffold AC text). **None is the normative spec, and none is the
  operator-named γ bootstrap declaration** — all three of blocker 4's actual targets (§8, §14, γ
  declaration) are corrected. Recommend a one-line hygiene cleanup of `alpha-closeout.md:11` for durable-
  record cleanliness, but it does **not** gate convergence.
- **F-R2.2 (process observation · not a spec defect · corroborates κ's P1).** The prior `beta-review.md`
  records a §R2 pass by another concurrent activation on this same claimed cell, and a mid-review mutation
  of `gamma-closeout.md`. This is the substrate race κ filed as P1 (distinct from #664/#665: two dispatch
  activations on one claimed cell). It did not corrupt this review — I verified against the settled current
  tree and reached my dispositions independently — but per κ's own note it "can just as easily produce
  conflicting α matter." Flagged upward as empirical support for prioritizing that P1; a deliberate review
  quorum, not accidental duplicate activation, is the right model if multiple reviews are wanted.

---

## Operator-gate hold candidates

**None.** No blocker required an unpinned architecture decision; every repair is a typed-contract-surface
fix consistent with the pinned D1–D10. The genuinely-open items (§16 Q1–Q5: wake-provider manifest count,
wave-dispatchability receipts, illustrative-command sequencing vs #504, exact schema field types, wave-
scope concurrency/idempotence predicate) are correctly carried as **open questions for downstream work**,
not silently resolved and not gating this Draft.

---

## Verdict

All six blockers **RESOLVED**; no regression in the "already good" parts; every class-specific `V`
predicate remains mechanically expressible; every State-A claim re-verified against shipped source; the §2
schema block and worked instance are provably **one shape**. The two findings are non-blocking receipt-
hygiene / process items. This is a genuine independent pass, but the reviewer is hosting-identity-limited
under one Sigma lineage (#664) — the exit sequence's external, non-Sigma CC ratification → operator-final-
read → merge → doctrine cell remains the appropriate place for hosting-independent review, and this
convergence does not authorize skipping it.

verdict: converge
