# External CC Ratification — cnos#662 (PC-D0)

## Corrective recursive measurement addendum — frozen R7 only

**Role:** independent external Cohering Cell, outside the Sigma activation lineage  
**Measured matter:** `docs/architecture/CELL-RUNTIME-CLASSES.md` at exact SHA `2d6b93cc4e69e5b413a80bd8e352cb0a004da460`  
**Measured receipt head:** `a0d39293a27cfe57b49dacff696345b1ee2cdb40`  
**Current-channel warning:** PR #667 has since moved to R8 matter `30a9455afaf7f9f63ca0c60b5c91477b321ca49f`, receipt head `4c7703940705c0ed88431408e31c19da445d1c10`. Its matter SHA-256 is `db171e69...`, not frozen R7's `80e0d8c...`. Per the immutable-target rule, **this comment does not measure or ratify R8**.

This addendum supersedes the **measurement and scope conclusions** in my earlier comment `5015281024`, not the two concrete blockers or their conservative `hold` effect. I withdraw its ad-hoc `CΣ=0.663874`, its telos-fit PASS, and its “no separate contract/FSM blocker” conclusion. That number was not produced from an explicitly frozen recursive target projection through `coh`, and the matter/source/evidence triad was too shallow for this architecture. The two earlier findings remain evidenced and have already been accepted into R8; no conclusion about the sufficiency of those repairs is made here.

## Instrument, target, and standing

- TSC source: `usurobor/tsc@26aab5023f03dc7d0abf82e5fdba20134fc6adad`.
- Released engine: `coh 0.12.0 (016c511)`.
- Measurement instruction: TSC `runtime/SELF-MEASURE.md` v3.2.4, SHA-256 `ca26d7a1a4dc6bd73e0afff558ed0342d42daeee193d4169de2c51b762759391`.
- Matter bytes: SHA-256 `80e0d8c68a3d8affabdd4bd14848cbb3f0bee27b078e435f61433e42a0ee89e0` at both the frozen matter commit and R7 receipt head.
- Six ordered, frozen projections were measured: L0 three-cell architecture, L1 kernel, L2 typed contracts, L3 FSMs, L4 execution/receipts, and their 30-file recursive union.
- Mechanical arm: three repetitions per projection, with exact equality of numeric scores and axis details.
- Semantic arm: one validated standard v3.2.4 witness per frozen prompt, ingested through `coh --mode hybrid --llm-response`. This is `k=1`, so it has **no semantic consistency standing**. I do not claim TSC Operational ACCEPT/REJECT, N≥30, bootstrap CI, OOD, scale, variance, or Lipschitz witnesses.

TSC computes pair coherence from the witness discrepancies as
`Coh(a,b)=exp(-φ(δ(a,b)))`, `φ(δ)=δ/(1-δ)`, and reports the symmetric aggregate
`CΣ=(sα·sβ·sγ)^(1/3)`. The semantic system result is therefore
`(0.36×0.24×0.30)^(1/3)=0.295945`. The geometric cross-level aggregate of L0–L4 is `0.425230`. These are diagnostics without standing; the directly evidenced hard-invariant failures below gate disposition independently of either number.

| Projection | Files | Mechanical α/β/γ → CΣ | Semantic α/β/γ → CΣ | Semantic δαβ/δβγ/δγα | Bottleneck |
|---|---:|---|---|---|---|
| L0 — three cells | 4 | .992/.985/.843 → **.937364** | .63/.40/.45 → **.484029** | .44/.50/.36 | β |
| L1 — one kernel | 4 | .949/.435/.843 → **.703264** | .43/.48/.58 → **.492848** | .45/.40/.42 | α |
| L2 — typed surfaces | 8 | .953/.850/.800 → **.865220** | .60/.30/.38 → **.408964** | .58/.60/.52 | β |
| L3 — cell/wave FSM | 12 | .998/.435/.950 → **.744362** | .44/.32/.28 → **.340346** | .58/.66/.60 | γ |
| L4 — #662 execution | 11 | .916/1.000/.925 → **.946381** | .38/.42/.46 → **.418726** | .55/.52/.50 | α |
| Recursive system | 30 | .925/.255/.850 → **.585211** | .36/.24/.30 → **.295945** | .68/.72/.62 | β |

The mechanical/semantic gaps are material evidence: textual regularity and cross-reference density make L0, L2, and L4 look strong, while semantic measurement exposes contradictory types and authority edges. Mechanical scoring is a structural proxy here, not a substitute for the semantic arm.

<details><summary>Frozen prompt bindings</summary>

- L0: `7726a0cddccead6b5f9379518da7d4aaa8a16c0c13a94cfc065a2c38b5e5b279`
- L1: `371b4a55979d45e712fbb22619b10a38f73bfbd815e5ef1016f68624a6fbb084`
- L2: `5c48f414e50b9714ddfb8b7d24a2bd0e92282a1178aa6efc2822703ebd295af9`
- L3: `0d94a81a0d98bc78b3e0259e4fc8d87c1fac3b62789a885736c4b9d29f2669af`
- L4: `00064fe387b5d4297fc58d2c8b5aa781dc796b0af61804180262ef4072576f4a`
- Recursive system: `f26c03de3e52bd5f7a901a3904a8351a5cfe345939d7fca097aef99764afde37`

</details>

## Recursive articulation and findings

### L0 — the three-cell system

**α / properties.** WC, PC, and CC are deployment classes of one kernel, distinguished by canonical output telos: artifact, relation graph, and process judgment. Matter domain is orthogonal.

**β / relationships.** The intended system edge is `intent → relational plan → artifact work → process judgment → V → δ/FSM → next cell`, with closed receipts observed through ε at scope n+1.

**γ / system process.** Class, authority, measurement instrument, and state must remain stable or migrate explicitly as that loop evolves.

**Finding L0-B1 — the #662 class does not match its canonical output.** Frozen §3.2 defines PC output as a relation graph, yet the actual #662 contract requests `kind: artifact`, the PC result is `artifact_ref`, and the normative deliverable is a merged Markdown architecture specification. Parent `CELL-RUNTIME.md` gives the controlling example: “a doctrine change is one kernel cell of class WC (pattern artifact—a doc), matter domain doctrine.” If a typed relation graph were canonical and Markdown only its projection, PC would fit. On the frozen record, the Markdown spec itself is canonical; it is WC telos. Planning a later wave and writing its doctrine contract are different outputs and should not be collapsed by calling both “planning.”

This is not merely a naming preference. It changes V predicates, result type, downstream consumer, receipt meaning, and the claimed TSC-β telos. The frozen PC-D0 does increase relational structure, but it does so **inside an artifact**; that does not make the artifact's canonical output type a relation graph.

### L1 — one cell and its kernel

**α / properties.** Every class must instantiate `α produce → β review → γ close → V validate → δ decide`; ε is external observation, not a sixth role.

**β / relationships.** Each signature must compose without requiring a value before it is produced, and α/β/κ firebreaks must remain role propositions independent of hosting identity.

**γ / process.** Repair must re-enter one declared state, preserve history, and end in an immutable receipt whose downstream effects remain outside CC authority.

**Finding L1-A1 — receipt production is circular across doctrine and shipped schema.** CCNF places γ receipt production before V and δ. The shipped `schemas/cdd/receipt.cue` requires that receipt already contain `validation` and `boundary_decision`, which are outputs of V and δ. Parent runtime also places CM after matter/receipt collection and says V consumes the CM result, while the frozen V signature is only `Contract × Receipt` and neither the frozen contract nor receipt carries a typed measurement result. No provisional `receipt_core → measurement → V → δ → final_receipt` relation is declared. The system therefore cannot type-check its own kernel ordering.

**Finding L1-A2 — κ is both outside and equal to α in frozen R7.** Frozen §8 and its γ closeout assert κ≠α/κ-outside and literal `κ=α`. Operator authorization and disclosed common hosting (#664) can authorize a separate α activation; they cannot make unequal functional roles equal. This was the earlier CC-1 and is retained only as an R7 finding.

### L2 — contracts, CM, V, evidence, and receipts

**α / properties.** A runnable contract needs immutable matter, target projection, CM reference, thresholds/standing policy, result type, and receipt form.

**β / relationships.** `contract → matter/review → CM measurement → γ binding → V → δ` must be an immutable, typed path. CM measures; V gates; δ effects.

**γ / process.** Shipped and future schemas require explicit adapters and lossless migration; identity-only control-plane references cannot become content revisions by assertion.

**Finding L2-B1 — #662 omits the measurement contract it inherits from #628.** Parent runtime explicitly requires `matter+receipt → CM → V` and says the schemas must gain `cell_class`, `matter_domain`, **and CM fields**. Frozen #662's canonical contract/four-schema chain contains no `cm_ref`, target-bundle digest, mode, thresholds/standing policy, measurement result, pair discrepancies, consistency, provenance, or immutable receipt edge to that result. “Class-specific V” is specified, but the measurement that V is meant to gate is not representable. This is the direct answer to whether CC merely applies a CM: normally yes—the CM is an input selected before judgment; the CC output is the process judgment. A CC may separately cause a new CM artifact to be produced, but that artifact cannot retroactively warrant the same judgment.

**Finding L2-B2 — shipped and specified contracts have no adapter.** `cnos.cdd.contract.v1` and specified `cn.cell.contract.v1` are disjoint shapes without a declared versioned projection, preservation rule, or negative fixtures. The State-A/B labels are honest, but honesty alone is not a migration.

**Finding L2-G1 — γ content-binds matter but not external review bytes.** The matter SHA edge is sound. The external-β edge in R7 is a mutable comment ID/URL, although frozen §§2/11.6/13 require content hash plus retrievable snapshot for control-plane evidence. That leaves `matter → review → γ` temporally unstable. This was earlier CC-2.

### L3 — FSM coherence at cell and wave scopes

The FSM first appears as a dedicated system at L3, but it is measured triadically and projected back into L1 kernel order and L2 types.

**α / properties:** one state/event/guard/action vocabulary, typed terminal/exceptional states, state closure, totality, determinism or normative priority, and reachability.  
**β / relationships:** exactly one authority per transition; command/table parity; guards aligned with contract, V, δ, labels, class, and domain.  
**γ / process:** owned exits for every nonterminal state, repair/recovery/resume, preserved history, migrations, wave aggregation, and terminal semantics.

**Finding L3-A1/B1 — shipped cell lifecycle authority is split.** `transitions.json` is locally deterministic, but the full lifecycle is not one table-governed machine. Frozen §11.1 says review return is `changes → todo`; §11.2's live sequence says `changes → in-progress`. `ready → todo` exists as an unguarded table rule while `cn issues dispatch` also owns that boundary, and `cn cell resume` owns alternatives. `blocked` is reachable/label-doctrine state but is absent from the table's declared state set. These surfaces do not define one closed transition system or one authority per edge.

**Finding L3-G1 — the wave “FSM” is a state sketch, not an executable FSM.** It lists a sequence but lacks a total event/guard/action relation, invalid-transition semantics, complete-with-residuals predicates, child-receipt aggregation, dependency/gate dispatch rules, holding/replanning exits, terminal semantics, and a total mapping from all CC dispositions to transition requests. Narrative judgment/FSM separation is preserved, but there is no executable edge by which a valid CC result reaches the single transition authority.

### L4 — the concrete #662 cell and εₙ

**α / properties.** Exact matter, review artifacts, receipt claims, identities, and state history must match the class contract.

**β / relationships.** Every receipt edge must bind exact content, and actual role occupancy must satisfy the declared firebreaks.

**γ / process.** Rounds should repair monotonically, supersede earlier records explicitly, and close without later mutable evidence changing the result.

R3→R7 is a substantive monotonic repair history. External β is independent of Sigma, names exact `2d6b93cc`, and authored no frozen matter; α≠β therefore passes. #664 honestly discloses the hosting-identity limitation. The instance nonetheless fails because its PC class/output disagree at the canonical boundary, its γ record asserts incompatible κ propositions, no CM was selected or bound, and its load-bearing β content can change under the recorded identity.

## Hard-invariant check — frozen R7

| Invariant | Result |
|---|---|
| One CCNF kernel for WC/PC/CC | PASS as doctrine |
| Class determined by canonical output telos | **FAIL** — canonical spec artifact is labeled PC relation output |
| Class and matter domain orthogonal | PASS in the target design |
| α≠β | PASS — exact-SHA external β authored no matter |
| κ outside and κ≠α/β/γ/δ | **FAIL** — frozen matter/receipt also assert `κ=α` |
| ε is cross-cell observation, not a role | PASS |
| Immutable matter/review/CM/measurement/receipt chain | **FAIL** — review bytes and CM/measurement edge absent |
| CM/V/δ separation with type-correct data path | **FAIL** — CM fields absent and receipt ordering circular |
| Four-schema boundary | PASS as a vocabulary boundary |
| Single dependency authority | PASS inside the proposed wave contract |
| Single transition authority per FSM edge | **FAIL** — table/dispatch/resume/spec conflict |
| Complete wave-FSM semantics before executable claim | **FAIL** — state sequence lacks total transition semantics |
| Honest shipped/specified/illustrative partition | PASS as labeling; it does not cure missing adapters |
| CC judgment distinct from FSM transition/operator acceptance | PASS narratively; executable mapping absent |

## Derived CM artifact

This measurement produced a separate candidate artifact, **`cnos Recursive Cell Coherence Methodology v0.1.0`**, rather than pretending the CC judgment itself is a CM. It defines the L0–L4 recursion, class classifier, typed CM edge, cell/wave FSM tests, thirteen gating invariants, disposition rules, semantic instruction, and executable #662 calibration manifests.

- Authority file SHA-256: `1e1534674007ad8f55044e59357e80ce53260fa05c4528d324ee2be1e0155074`.
- Final 11-file CM self-measure prompt SHA-256: `b3447ecec9516827abfc5a99ef85aef64bc7503eefcceeef3c7596e5681a696c`.
- TSC schema validation: PASS.
- Mechanical self-measure, N=3 exact: α=1.000, β=1.000, γ=.7925, `CΣ=.925408`, γ bottleneck.
- Semantic self-measure, k=1: α=.93, β=.82, γ=.58, `CΣ=.761918`, γ bottleneck.
- Standing: **none**. It was derived from #662, has no held-out anchor or semantic ensemble, and current TSC cannot consume arbitrary CM declarations first-class. It is a candidate outcome for a later WC/CM lifecycle, not independent warrant for this review.

## Disposition

**DISPOSITION for frozen R7 `2d6b93cc…`: `request_planning`.** Ratification is withheld. The minimum next MCA is planning because the failures are not only two local prose/receipt repairs: canonical class ownership, the CM→receipt→V→δ type path, shipped→specified contract migration, transition authority, and complete wave-FSM semantics require a coherent relation/FSM redesign before a Working Cell should implement them. This is a more precise routing of the prior conservative `hold`.

Because PR #667 has moved, this disposition is historical and bound only to frozen R7. It is **not** a verdict on R8 `30a9455a…`; R8 requires its own exact-SHA β→γ→CC chain.

I do **not** merge, mark ready, change labels, perform δ/FSM state transition, author a repair to the frozen matter, or replace operator acceptance. This judgment feeds the operator's final-read gate.

**TL;DR: `request_planning` for frozen R7 only — not ratified at exact SHA `2d6b93cc4e69e5b413a80bd8e352cb0a004da460`; current R8 is unjudged.**

