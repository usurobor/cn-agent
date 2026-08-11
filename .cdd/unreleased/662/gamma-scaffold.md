---
schema: cn.cdd.gamma-scaffold.v1
issue: 662
cell_class: planning
mode: PC-D0
protocol: cds
matter_domain: doctrine
main_sha: 5ca785cd
cycle_branch: cycle/662
wake: cds-dispatch
run_class: first_pass
---

# γ scaffold — cnos#662 (PC-D0: Cell Classes and Mechanical FSM architecture note)

## 1. Cell contract (read from the issue, not improvised)

```yaml
cell:
  class: planning
  protocol: cds
  matter_domain: doctrine
  issue: 662
  requested_output: [ "docs/architecture/CELL-RUNTIME-CLASSES.md" ]
  non_goals: [ implement_cell_runner, alter_ccnf, add_new_role_letters,
               replace_cdd_or_cds, file_or_dispatch_child_issues ]
  gates: { operator_authorization_required: true, operator_acceptance_required: true }
  stop_conditions: [ source_doctrine_conflicts, required_decision_unresolved ]
```

**Required output — exactly one normative planning artifact:** `docs/architecture/CELL-RUNTIME-CLASSES.md`, a sibling of `docs/architecture/CELL-RUNTIME.md` (which defines the WC/PC/CC *classes*; this note operationalizes them into the *FSM/routing*). Status header: `Proposed architecture note (realization layer). Not ratified.` — matching `CELL-RUNTIME.md`'s own header convention; ratification is a separate CDD `doctrine` cell, not this cell.

## 2. Source of truth (α reads these; α does not invent architecture)

| # | Source | What it supplies |
|---|---|---|
| 1 | Issue #662 body (`κ cohering pass` + embedded draft §1–§20) | The full architecture draft (κ-corrected), findings F1–F5, AC1–AC7 (cell-level) and AC1–AC8 (note-level, §20) |
| 2 | Issue #662 operator-authorization comment (usurobor, 2026-07-13T19:14:52Z) | 10 normative decisions (D1–D10 below) that supersede unresolved wording in the embedded draft |
| 3 | `docs/architecture/CELL-RUNTIME.md` (#628, landed) | The **class** doctrine — WC/PC/CC as output-telos deployment shapes of one CCNF kernel, the four orthogonal axes, CM/V/δ separation. This note **operationalizes** it into FSM/routing; it does not restate or contradict it. |
| 4 | `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` (CCNF) | The five-step kernel (`α→β→γ→V→δ`), ε as cross-cell receipt-stream observer (not a kernel step), Scope-Lift Projection 3 |
| 5 | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` (shipped FSM) | The actual shipped state set (`ready, todo, in-progress, review, changes` + `blocked` reachable via BLOCK-REQUEST.yml) and guard vocabulary (`run_active`, `branch_exists`, `pr_exists`, `review_request_present`, `claim_request_present`, `block_request_present`, `release_request_present`, `repair_contract_present`, …) — AC6 requires the note stay compatible with this, not merely with the note's own illustrative §11 table |
| 6 | `cn --help` shipped command surface (`cell return`, `cell resume`, `cell finalize`, `issues fsm evaluate`, `issues fsm scan`, `issues dispatch`) | Ground-truths State A (§17) — `run`/`pulse`/`measure`/`bundle`/`act` are State B / illustrative, not shipped |
| 7 | #530 (open) | Owns the typed scope-lift / ε-projection artifact — this note depends on it, does not redefine it |
| 8 | #500 (closed/shipped) | review-return (`status:review → status:changes → status:in-progress → status:review`) — cite as shipped |
| 9 | #504 (open) | stale-claim recovery — cite as genuinely open, "Sub C of #583" |
| 10 | #644 (open, `status:backlog`) | PC-Wave mechanization (α produces wave graph, β reviews, γ closes) — cite as the PC-Wave mode's implementation target, not yet built |
| 11 | #654 (open, `status:backlog`) | PC D0-contract hardening — cite as the D0 contract this note's own PC-D0 mode is an instance of |
| 12 | #627 (open, parent wave) | cite as parent; AC-INT (one generic runner routes WC/PC/CC by `cell_class`, no separate PC provider) resolves note Open-Q3 |

## 3. The 10 pinned operator decisions (D1–D10) — supersede unresolved draft wording

Recorded verbatim-in-substance from the operator-authorization comment; α MUST carry these into the formalized note and MUST NOT re-litigate them:

1. **D1 — One CCNF kernel.** WC/PC/CC are output-telos classes of the same kernel (already the draft's own §1–§2; confirm, don't re-derive).
2. **D2 — κ remains outside the cell.** κ supplies operator/control-plane input; the Planning Cell's α owns and writes the specification. κ must not author the specification as α. (This governs *this cell's own execution*, and the note must also state it as general doctrine.)
3. **D3 — CC ↔ ε.** ε is the cross-cell receipt-stream projection; CC is the executable cell-class realization at scope n+1 that consumes εₙ observations and produces coherence judgment / next MCA. CC is not ε and adds no new role to the closure. (Matches issue §CC↔ε reconciliation and AC2 — carry verbatim per the issue's own instruction.)
4. **D4 — Wake topology.** v0 uses one generic package runner per protocol; the CDS runner routes Working/Planning/Cohering; no separate PC or CC wake provider in v0; CC runs as a claimed cell in v0; scheduled coherence pulses are a future extension.
5. **D5 — Human gates.** Human authority at irreversible/scope-expanding boundaries: intent acceptance; wave authorization; production/release; final acceptance of doctrine-affecting matter; explicit hold/block/escalation. A wave may pre-authorize its dependency-respecting internal nodes. The operator must not become the scheduler for every child.
6. **D6 — Cell class as an FSM dimension.** Lifecycle label-based; behavior class-based; the runtime eventually evaluates `(status, protocol, cell_class)`. Dispatch readiness remains solely `dispatch:cell + protocol:{P} + status:todo`; `cell_class` is contract/routing metadata, not a second readiness gate. (Matches draft F2/AC3 — confirm, don't re-derive.)
7. **D7 — CC judgment vs FSM transition.** CC owns the wave-state judgment; the mechanical FSM owns the wave-state transition.
8. **D8 — Class-specific V.** Common V: receipt complete; evidence present; role ownership valid; contract satisfied; no unresolved blocking finding. WC V adds implementation ACs, tests, allowed surfaces. PC V adds planning artifacts, graph acyclicity, complete child contracts, gates, STOP conditions, and no child auto-dispatched. CC V adds an evidence-backed judgment, exactly one disposition, a typed next-MCA or human gate, and no implementation-surface modification.
9. **D9 — Schema-first destination.** The document names, but does not implement: `cn.intent.v1`, `cn.cell.contract.v1`, `cn.next-mca.v1`, `cn.wave.v1`.
10. **D10 — State truth.** Clearly distinguish what is shipped now / what this architecture specifies / what remains illustrative or future. Do not present `cn cell run`, `cn cell pulse`, mechanical class routing, or scheduled CC pulses as shipped commands.

## 4. AC oracle list (α self-verifies against this; β independently re-walks it)

Cell-level ACs (issue body, "Acceptance criteria" section):

- **AC1** — Note formalized at `docs/architecture/CELL-RUNTIME-CLASSES.md`, sibling of `CELL-RUNTIME.md`; note's own AC1–AC8 (§20 of the draft) hold.
- **AC2** — F1/D3 CC↔ε wording carried verbatim (substance); cites #530 as owner of the typed ε-projection artifact.
- **AC3** — F2/D6 `cell_class` pinned as post-claim routing metadata; dispatch readiness stays label-gated per #643; contract-incompleteness routes to `status:blocked` + `degraded_reason`. No contradiction with #640/#643.
- **AC4** — F3/D2 preserved: note nowhere implies κ=α; κ≠α + actor-collapse rule stated.
- **AC5** — F4/D10 applied: State A reflects shipped reality (`cn cell return/resume/finalize`; CDR shipped #376; CDW illustrative; #504 open); future command names marked illustrative.
- **AC6** — Recorded cohering pass shows the note restores coherence across #627, #644, #654, #500, #504, #530, #583/#584, the κ boundary, and the Go mechanical-orchestration target.
- **AC7** — No child waves/implementation issues filed or dispatched by this cell without explicit operator authorization.

Note-level ACs (draft §20, must hold in the formalized doc):

- **AC1** — WC/PC/CC defined as output-telos classes of one CCNF kernel.
- **AC2** — Artifact names are not treated as cell-class identifiers.
- **AC3** — CC, PC, WC contracts distinct and example-backed.
- **AC4** — The CC → PC → CC → WC loop is specified.
- **AC5** — The FSM extension point for `cell_class` is named.
- **AC6** — Existing status lifecycle and dispatch selector remain compatible (verify against shipped `transitions.json`, not just the draft's illustrative §11 table).
- **AC7** — Mechanical guards prevent role/class collapse.
- **AC8** — State A and State B are explicit.

Plus D1–D10 (§3 above) folded in throughout, and non-goals honored: no CCNF redefinition, no CDD/CDS replacement, no new role letters, no α/β/γ/δ semantics change, no runtime implementation, no child issue filing/dispatch.

## 5. α prompt (implementation contract)

α's job: **write exactly one file**, `docs/architecture/CELL-RUNTIME-CLASSES.md`, formalizing the embedded draft in issue #662 (§1–§20 of the issue body) with:
- F1–F5 findings and D1–D10 operator decisions folded in as settled text (not left as open questions);
- style/structure matching `docs/architecture/CELL-RUNTIME.md` (Status header, Thesis-style framing where natural, a Reconciliation-map-style table, Non-goals, Open questions pared down to what's genuinely still open after D1–D10 resolve most of the draft's §19 list, and an authoring note);
- §11 "FSM Events" table cross-checked against the shipped `transitions.json` state/guard vocabulary (cite the mechanical `blocked` state reachable via `BLOCK-REQUEST.yml`, and the request-marker-file pattern — `CLAIM-REQUEST.yml`/`REVIEW-REQUEST.yml`/`BLOCK-REQUEST.yml`/`RELEASE-REQUEST.yml` — as the actual shipped mechanism behind "guard"), not merely restating the draft's simplified illustrative table as if it were shipped;
- explicit "what is shipped / what this specifies / what is illustrative-future" framing per D10, distinct sections or inline markers;
- citing #530, #500, #504, #644, #654, #627, #628, #583, #584 by number wherever the draft or D1–D10 reference them;
- an authoring note (κ-style, matching `CELL-RUNTIME.md`'s closing "Authoring note" convention) naming this cycle (cnos#662, PC-D0, cds-dispatch wake, δ wake-invoked mode) and summarizing what α reconciled.

α also writes `.cdd/unreleased/662/self-coherence.md §R0` — walking every AC in §4 above and stating how the written doc satisfies it (file:line or section citation), plus a review-ready signal.

**Hard constraints (non-goals, repeated for α):** no Go/runtime implementation; no schema implementation; no FSM code; no wake changes; no child implementation issues; no PC-Wave; no child `status:todo` labels; no dispatch of work derived from the note; no separate PC/CC provider; no modification of CCNF role semantics. If α discovers a genuine contradiction in the pinned decisions (D1–D10) or between the issue draft and shipped state that cannot be resolved by exposition alone, α STOPS and names the contradiction rather than inventing a resolution.

## 6. β prompt

β's job: independently re-walk every AC in §4 (cell-level AC1–AC7, note-level AC1–AC8) and every D1–D10 decision against the written `docs/architecture/CELL-RUNTIME-CLASSES.md`, **without** re-reading α's self-coherence.md first (review the doc against the sources, then cross-check against self-coherence.md). Specifically verify:
- relation consistency across CCNF, `CELL-RUNTIME.md`, #530, #500, #504, #644, #654, the κ boundary, and the shipped FSM (`transitions.json`);
- State A / State B honesty (D10) — no shipped-command overclaim;
- no policy was invented silently (anything not traceable to the issue draft, F1–F5, or D1–D10 is a finding);
- all class-specific guards and V predicates (D8) are mechanically expressible, not just narratively plausible;
- non-goals honored (no runtime code, no child dispatch, no CCNF/role-letter changes).

β writes `.cdd/unreleased/662/beta-review.md §R0` with a verdict (`converge` or `iterate`) and, on `iterate`, itemized findings α must address in `§R1`.

## 7. Scope guardrails

- **Allowed surface:** `docs/architecture/CELL-RUNTIME-CLASSES.md` (new file) + `.cdd/unreleased/662/*` (cell artifacts).
- **Forbidden surfaces:** `.github/workflows/**`; any Go source (`src/go/**`); any `src/packages/**/skills/**` (role/CCNF doctrine); any other `docs/architecture/*.md`; any issue creation/labeling beyond this cell's own lifecycle transitions.
- **No child dispatch.** This cell does not file or label child issues, per AC7 / the issue's non-goals / D5.

## 8. Friction notes

None yet — this is R0.
