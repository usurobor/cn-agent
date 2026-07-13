verdict: converge

# β review — cnos#662 (`docs/architecture/CELL-RUNTIME-CLASSES.md`)

## §R0

Independent, adversarial β review of the PC-D0 planning-cell spec. I re-walked the AC oracle
(cell-level AC1–AC7, note-level AC1–AC8) and the ten pinned decisions D1–D10 against the four
ground-truth sources in the repo, formed my own view, and only then read
`.cdd/unreleased/662/self-coherence.md`. Verdict: **converge**. No blocking finding survived
verification; the non-blocking items I considered and rejected are recorded under §R0.7 for the
record so a later reader can see they were checked, not missed.

### §R0.1 — Dimension A: Relation-consistency (CCNF, CELL-RUNTIME.md, cited issues, κ, shipped FSM)

ATTEST — coherent.

- **CCNF ε / Projection 3.** Spec §4's carried block and §4's closing paragraph frame CC as a full
  scope-`n+1` cell consuming εₙ's receipt-stream projection, "not ε wearing a different hat." This
  matches CCNF §Scope-Lift **Projection 3** exactly: `εₙ receipt-stream observation → γₙ₊₁-like
  coordination / evolution`, with the explicit projection-not-flat-renaming clause
  (`COHERENCE-CELL-NORMAL-FORM.md:276`, `:289`, `:297`). The spec's "one lineage, not two competing
  coherence surfaces" is faithful to CCNF's "the projection carries the observation, not the role."
- **CELL-RUNTIME.md (#628).** Spec is a sibling that operationalizes, not restates: §1 table and §3
  header ("restated here only enough to anchor the contract-level detail `CELL-RUNTIME.md` leaves
  implicit") keep the four-orthogonal-axes frame and the CM/V/δ separation intact. No contradiction
  of the deployment-shape framing (`CELL-RUNTIME.md:36-42`).
- **#530** (typed ε-projection artifact, open): §4 line 137 and §14 depend on it, explicitly "files
  no work against #530." Correct — not redefined.
- **#500** (review-return, shipped): §11.2 line 252 cites `review → changes → in-progress → review`
  as live; consistent with the shipped `changes` state transitions in `transitions.json:231-250`.
- **#504** (stale-claim recovery, open): §11.2 line 253 marks it "genuinely open, Sub C of #583";
  the shipped `scan` sweep is correctly presented as covering only the mechanical dead-run sweep, not
  #504's fuller resume-or-escalate design.
- **#644 / #654** (open targets): §3.2, §11.4, §13, §14 name them as citable implementation targets
  the spec "specifies, does not build." No overclaim of their state.
- **#583 / #584** (mechanism/cognition boundary, landed): §8, §10 cite them as the doctrine the
  κ-outside-the-cell rule depends on. Consistent.
- **κ boundary:** §8/§14 carry it forward; §16 narrates κ (issue author) vs α (file author) as
  distinct instances. See §R0.3.
- **Shipped FSM:** verified against `transitions.json` — see §R0.2.

### §R0.2 — Dimension B: State-A / State-B honesty (D10)

ATTEST — no not-yet-shipped mechanism is presented as shipped. Every State-A claim independently
re-verified against source:

- **Declared `states` array.** Spec §11.1 line 222: `["ready","todo","in-progress","review",
  "changes"]`. Exact match to `transitions.json:18`.
- **`blocked` is a transition target, not a declared state.** Spec §11.1 line 222 states this
  precisely. Confirmed: `blocked` appears only as `target_state` in the `in-progress` block
  (`transitions.json:155-158`), never in the `states` array. This is the one nuance the spec is
  careful about, and it is correct.
- **Guard vocabulary.** Spec §11.1 line 224 lists all twelve guards; every one matches the `guards`
  map in `transitions.json:19-31` (`run_active, branch_exists, branch_has_commits, pr_exists,
  pr_has_commits, review_request_present, repair_contract_present, cdd_artifacts_present,
  checks_passing, claim_request_present, block_request_present, release_request_present`).
- **Request-marker-file pattern.** Spec §11.1 marker table matches the shipped rules:
  `CLAIM-REQUEST.yml` → todo→in-progress gated `run_active` false (`transitions.json:110-116`);
  `REVIEW-REQUEST.yml` → in-progress→review needs `pr_exists`+`pr_has_commits`, #574 tightening
  (`:137-142`); `BLOCK-REQUEST.yml` → in-progress→blocked, no further guard (`:153-158`);
  `RELEASE-REQUEST.yml` → in-progress→todo (no matter) or delta-recovery (matter exists), #368
  protection (`:161-176`). `changes → todo` repair re-entry gated `repair_contract_present`, #516
  (`:235-240`). All accurate.
- **Command surface ships as claimed.** `cn cell return/resume/finalize` present
  (`cmd_cell.go:38,120,212` → `cell-return/cell-resume/cell-finalize`). `cn issues dispatch` present
  (`cmd_issues_dispatch.go:31`). `cn issues fsm evaluate` AND `cn issues fsm scan` both ship: the
  issuesfsm package parses `evaluate`, `scan`, `terminal` sub-verbs
  (`src/packages/cnos.issues/commands/issues-fsm/issuesfsm.go:69-79`; scan is cnos#593). NOTE: the
  cli-wrapper doc-comment `cmd_issues_fsm.go:11` says "exactly one sub-verb, evaluate" — that comment
  is stale relative to the package, but the spec's claim that `scan` ships is correct against the
  actual implementation, so no finding.
- **`run`/`pulse`/`measure`/`bundle`/`act` are illustrative-future.** Confirmed absent: no
  `cell-run`/`cell-pulse`/`cell-measure`/`cell-bundle`/`cell-act` command registered anywhere in
  `src/go/internal/cli/`. Spec §11.5 correctly fences `cn cell pulse`/`cn cell run` as "not shipped."
- **CellKind seam is observation-only.** `TestSeam_CellKindNotEnforced` exists
  (`issuesfsm_test.go:810`); `CellKind{Observed, Source, DefaultedTo}` exists
  (`issuesfsm_test.go:799-806, 843`); `TestSeam_CellKindDefaultedWhenAbsent` locks
  `Source:"absent"`, `DefaultedTo:"implementation"`. Spec §6/§11.3 correctly presents `cell_class`
  as the *specified promotion* of that seam from observation to evaluation, not as already enforced.

### §R0.3 — Dimension C: No policy invented silently

ATTEST.

- κ≠α is stated (§8 line 190 "κ must not author the specification as α"; §12 line 303 "κ ≠ α/β/γ/δ")
  and κ nowhere authors as α; the actor-collapse case is required to be *declared*, not assumed
  equivalent (§8 line 190) — traceable to D2/F3.
- CC↔ε block is carried as a fenced verbatim block (§4 lines 118-135), matching D3 and consistent
  with CCNF Projection 3. (Bootstrap caveat in §R0.6: the raw issue-body §CC↔ε text is not on disk in
  my binding inputs, so I verified faithfulness-to-CCNF + internal coherence + match-to-D3 rather
  than a byte-diff against the issue body.)
- `cell_class` is post-claim routing metadata and dispatch readiness stays `dispatch:cell +
  protocol:{P} + status:todo` only (§6 line 164, §2 line 67) — traceable to D6/AC3; contract-
  incompleteness → `status:blocked` + `degraded_reason: cell_contract_incomplete` matches scaffold
  AC3. No contradiction with #643/#640 (#640 is the real shipped dispatch primitive,
  `cmd_issues_dispatch.go:9-13`).

### §R0.4 — Dimension D: Mechanical expressibility of class-specific guards/V (D8)

ATTEST — the two structural invariants the task flags are mechanically expressible, not merely
narrative:

- **PC "no child auto-dispatched"** (§5, §12 line 300): reduced to "the PC applied no `status:todo`
  to any child" — a checkable predicate over label-application events. Mechanical.
- **CC "no implementation surface modified"** (§5, §12 line 301): reduced to "matter paths ⊆ judgment
  artifacts, no code/product diff" — a checkable predicate over artifact-paths / branch diff.
  Mechanical.
- Common floor (receipt complete, evidence present, role ownership valid, contract fields satisfied,
  no unresolved blocking finding), WC additions (executable ACs pass, `checks_passing`, allowed-
  surface/non-goal guards), PC additions (graph acyclic, complete child contracts) are all predicates
  over receipt fields / artifact paths / labels. `checks_passing` is a real shipped guard
  (`transitions.json:28`). No V predicate is narrative-only.

### §R0.5 — Dimension E: Hard non-goals respected

ATTEST. `git diff --stat origin/main...cycle/662` shows exactly four files:
`.cdd/unreleased/662/{CLAIM-REQUEST.yml, gamma-scaffold.md, self-coherence.md}` and
`docs/architecture/CELL-RUNTIME-CLASSES.md`. Exactly one spec file plus the cell's own
`.cdd/unreleased/662/` artifacts. No Go/schema/FSM code; no other `docs/architecture/*.md` touched;
no `.github/workflows/**`; no child issue filed or dispatched. §15 states the non-goals reflexively
("§3.2's own guard applies reflexively to the cell that produced this note"). Clean.

### §R0.6 — Dimension F: FSM + wave FSM + three-cell contracts + Coherence-Loop nomenclature + State A→B

ATTEST — all present and internally consistent.

- Cell FSM: §11.1 (grounded in shipped `transitions.json`).
- Wave FSM: §11.4 (specified, explicitly unshipped; typed `wave_transition_request` shape).
- Three-cell contracts: §3.1 (WC), §3.2 (PC, both D0 and Wave modes), §3.3 (CC), each with input
  admissibility, matter, output-pass criteria, and a result shape.
- Coherence-Loop nomenclature — all five terms present and used consistently: **Coherence Loop**
  (§7 heading), **Cohering Cell** (§3.3), **Cell Runtime** (title / §1), **Cell Runner** (§10 line
  210 "one generic Cell Runner"), **Wake Provider** (§10 "no separate PC or CC wake provider",
  "wake-provider profiles").
- State A→State B: §11's three-way partition (11.1–11.2 shipped / 11.3–11.4 specified / 11.5
  illustrative-future) plus §10's State-B framing. The added "specified-but-unshipped" middle band is
  a legitimate refinement of the draft's binary framing, not a policy invention — it separates this
  note's own specification from both shipped reality and illustrative futures.

### §R0.7 — Non-blocking items considered and rejected (not findings)

1. Stale cli-wrapper comment `cmd_issues_fsm.go:11` ("exactly one sub-verb") contradicts the shipped
   package (evaluate/scan/terminal). This is a defect in *repo source outside this cell's surface*,
   not in the spec; the spec's claim that `scan` ships is correct against the actual implementation.
   Out of scope for this matter and correctly not touched (forbidden surface).
2. §2's worked-instance YAML nests `requested_output`/`non_goals`/`gates`/`stop_conditions` under
   `cell:` while the envelope template (§2 lines 34-51) places them at other keys. This is a faithful
   quote of the γ-scaffold §1 format, explicitly labeled as such, and D9 makes exact field placement
   non-load-bearing (schema named, not implemented). Cosmetic, not a contradiction.
3. §6 "cds, and eventually cdr/others" reads as under-claiming CDR relative to scaffold AC5's "CDR
   shipped #376." Under-claim is the safe direction for State-A honesty, the note's protocol focus is
   cds, and this concerns the FSM's protocol *dimension* rather than whether a CDR package exists.
   Not an overclaim; not actionable.

### §R0.8 — Contradiction-rule check

No TRUE contradiction found in the pinned architecture (D1–D10 / settled findings / shipped state).
The one nuance requiring exposition — the shipped `states` array omitting `blocked` while `blocked`
remains a reachable `target_state` — is a draft-shorthand imprecision resolved by exposition in
§11.1 against shipped ground truth, not an irreconcilable conflict. No candidate operator-gate hold
to raise.

### §R0.9 — β independence disclosure (#664 structural limitation)

I am a **separate Agent activation** spawned by the δ/dispatch driver, running under the **same
account/session lineage** as the cell that authored the matter. Per the #664 structural limitation,
α≠β does NOT hold at the hosting layer (same session lineage). **Protocol-level independence DOES
hold:** I did not author `docs/architecture/CELL-RUNTIME-CLASSES.md` or any of the cell's `.cdd`
artifacts; I re-verified every State-A claim directly against repo ground truth
(`transitions.json`, the CLI command definitions, the issuesfsm package, the seam test); and I
formed my full per-dimension view BEFORE reading `self-coherence.md` (read last, as instructed).
Independence is therefore **bootstrap-limited at the hosting layer, full at the protocol layer.**

---

**Verdict: converge.** The spec operationalizes CELL-RUNTIME.md without restating or contradicting
it, honestly partitions shipped vs specified vs illustrative-future state (D10) with every State-A
claim independently verified against source, carries D1–D10 as settled input, keeps every class-
specific V predicate mechanically expressible, invents no untraceable policy, and wrote exactly one
spec file plus its own cell artifacts. Ready for the exit sequence (separate CC review →
operator-final-read → merge → doctrine cell) the note's §16 names.
