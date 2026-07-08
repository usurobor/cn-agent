---
name: delta
description: δ role-skill. Owns the two-sided coherence-cell membrane (outward boundary decision on receipt + V verdict; inward implementation-contract enrichment at dispatch) and the override authority that compose with V's verdict per RECEIPT-VALIDATION.md.
artifact_class: skill
kata_surface: embedded
governing_question: What does δ-as-role decide at the cell's boundary — both the receipt-facing outward gate and the dispatch-facing inward enrichment — and what authority does δ hold over override?
visibility: internal
parent: cdd
triggers:
  - boundary
  - membrane
  - override
  - receipt
  - verdict
  - gate
  - dispatch-enrichment
  - implementation-contract
scope: role-local
inputs:
  - typed receipt emitted by γ at close-out (RECEIPT-VALIDATION.md §Validation Interface)
  - ValidationVerdict emitted by V (RECEIPT-VALIDATION.md §Output contract)
  - γ-authored dispatch prompts containing the `## Implementation contract` section (wire-format at cnos.handoff/skills/handoff/dispatch/SKILL.md; γ-side authoring at gamma/SKILL.md §2.5)
  - repo conventions (language, package scoping, integration-target, etc. — for inward enrichment)
outputs:
  - BoundaryDecision recorded against (receipt, verdict) — one of {accept, release, override, reject, repair_dispatch}
  - dispatch-prompt enrichment when γ leaves implementation-contract rows unpopulated (or escalation to operator-as-human when undecidable)
  - override declarations with structured override block when proceeding against non-PASS verdict
requires:
  - active CDD cycle exists with a closed cell at scope n (receipt + verdict materialised)
  - "or: active γ dispatch prompt awaiting δ's pre-routing enrichment"
calls:
  - RECEIPT-VALIDATION.md
  - COHERENCE-CELL-NORMAL-FORM.md
  - operator/SKILL.md
---

# δ (delta) — the cell's boundary actor

## Core Principle

**δ is the coherence cell's boundary — the actor that decides what crosses from scope `n` to scope `n+1`, and the actor that enriches what crosses from γ's protocol-level contract into α's executable dispatch.**

δ is a **two-sided membrane**:

```text
δ as two-sided membrane:

  outward:  receipt + V verdict → parent-scope boundary decision   (§1)
            (V dereferences evidence from the receipt; δ records
             a BoundaryDecision against the verdict)

  inward:   γ contract → α-ready dispatch                           (§2)
            (implementation-contract enrichment happens here:
             γ writes protocol-level; δ pins implementation-level)
```

`COHERENCE-CELL-NORMAL-FORM.md` (CCNF) names δ at step 5 of the kernel — δ reads the receipt and V's verdict and records a `BoundaryDecision`. `RECEIPT-VALIDATION.md` freezes the typed interface: V emits a `ValidationVerdict`; δ emits a `BoundaryDecision`; the two surfaces are independent and compose by an explicit rule (§4 below). The δ-inward function is the complementary face — the dispatch-time enrichment that pins the implementation-contract axes γ's protocol-level contract leaves unpinned.

**δ-the-role lives here.** δ-as-actor in a running CDD cycle (γ=δ collapse cases, dispatch coordination, gate routing) lives in [`operator/SKILL.md`](../operator/SKILL.md); the harness substrate (dispatch invocation, observability, identity, polling, timeout recovery) lives in [`harness/SKILL.md`](../harness/SKILL.md) (Phase 4b of cnos#366, landed at cnos#398); release-effector mechanics live in [`release-effector/SKILL.md`](../release-effector/SKILL.md) (Phase 4c of cnos#366, landed at cnos#399). The four surfaces are deliberately separate: the role-skill is portable (any role-scope ladder instantiation per `ROLES.md` can carry it); the substrate skills are the substrate-specific realizations on the current substrate (per CCNF §Two-Layer Separation: kernel is the *what*, realization is the *how-on-this-substrate*).

> **First-time operator?** The dispatch-coordination layer (γ=δ collapse, single-session vs multi-session, harness quiescence, timeout recovery) is in [`operator/SKILL.md`](../operator/SKILL.md) (coordination) + [`harness/SKILL.md`](../harness/SKILL.md) (mechanics). δ-the-role's authority and policy live here.

---

## 1. Outward membrane — receipt + verdict → boundary decision

The outward face of δ runs after γ closes the cell and emits the receipt. V validates the receipt against the contract (per CCNF step 4 + `RECEIPT-VALIDATION.md` §Q1) and emits a `ValidationVerdict`. δ then records a `BoundaryDecision` — the outward act that determines whether and how the closed cell crosses to scope `n+1`.

```text
γ closes cell                    →   closed_cell exists; receipt drafted
γ emits receipt                  →   receipt exists; cell is "closed", not "accepted"
δ invokes V on receipt           →   V emits ValidationVerdict
δ records BoundaryDecision       →   {accept | release | override | reject | repair_dispatch}
ACCEPTED                         →   receipt transmissible to parent scope
```

(Per `RECEIPT-VALIDATION.md` §Q1 ordering rule. `V` is a typed predicate exposed as a capability per §Q2; δ invokes the capability, not the operator-facing `cn-cdd-verify` command.)

### 1.1. External actions δ-the-role authorises

δ-as-role holds **authority** for the platform actions that cross the boundary. The mechanics of executing these actions live in three companion surfaces: [`operator/SKILL.md`](../operator/SKILL.md) §3 (gate-action confirmation protocol, request-vs-observation discipline); [`harness/SKILL.md`](../harness/SKILL.md) (dispatch invocation, observability, identity, polling, timeout recovery — landed via Phase 4b of cnos#366, cycle/398); [`release-effector/SKILL.md`](../release-effector/SKILL.md) (`scripts/release.sh` invocation, post-push CI polling, recovery runbook, branch deletes — landed via Phase 4c of cnos#366, cycle/399). The authority-naming surface lives here:

| Action | Authority | Mechanics location |
|--------|-----------|--------|
| Pre-merge gate validation | δ authorises; γ requests via `scripts/validate-release-gate.sh --mode pre-merge` (`cnos.cds/skills/cds/CDS.md` §"Artifact contract" → §"Ownership matrix" + `gamma/SKILL.md` §2.10) | `operator/SKILL.md` §3.1 |
| Push β-approved merge to main | δ executes on β's behalf when env/auth blocks β; this is execution of β's integration authority, not δ approval | `operator/SKILL.md` §3.1 |
| Release-boundary preflight | δ verifies merge commit, release artifacts, tag/deploy preconditions, platform readiness; decision is proceed / request changes / override (`cnos.cds/skills/cds/CDS.md` §"Development lifecycle" → §"Step table" Step 9 — δ gate) | `operator/SKILL.md` §3.1 |
| Tag push + release | **δ is sole tag-author** — β does not tag; only δ creates tags per cycle | `release-effector/SKILL.md` |
| Branch delete | After cycle closure and merge | `release-effector/SKILL.md` §5 |
| Issue filing on external repos | Cross-project dependency | `operator/SKILL.md` §3.1 |
| Force push | Rebase required with env constraints | `operator/SKILL.md` §3.1 |
| Auth refresh | Token/permission expiry | `operator/SKILL.md` §3.1 |

The role-policy that governs the release boundary: **the triad's work is not complete until tagged.** Untagged post-cycle patches on main are an open boundary — the triad's output is still entangled with whatever comes next. The tag is the disconnection point. δ-as-role commits to this property; the `scripts/release.sh` runbook lives in [`release-effector/SKILL.md`](../release-effector/SKILL.md) (Phase 4c of cnos#366 — landed via cycle/399; the doctrinal frame stays at `operator/SKILL.md` §3.4).

**δ blocks release completion until CI is green and owns recovery on red.** This is δ-role policy:

- **CI Green** → δ declares release complete.
- **CI Red** → δ owns the failure. Investigate release logs; classify as release-specific vs pre-existing infrastructure; fix or escalate; re-verify; operator override (§3 below) is the explicit escape hatch for known pre-existing failures (cf. v3.66.0/v3.67.0 smoke failures).

The gate does not close until CI is green or δ explicitly records an override (§3).

### 1.2. Execute on request, not on observation

Gate actions fire when a role **requests** them, not when δ notices they're needed. Observing that a tag isn't pushed or a branch exists is not a gate trigger — γ's explicit request is.

- ❌ Heartbeat shows tag not pushed → δ pushes it (role leak: δ decided the gate, not γ)
- ❌ "β asked me to push the merge but I think we should wait for one more review"
- ✅ γ requests tag push → δ pushes it and confirms
- ✅ If you disagree with a gate request, declare an override (§3)

### 1.3. Report completion

After executing a gate action, confirm to the requesting role that the action completed.

- ❌ Execute silently and assume the triad will notice
- ✅ "Tag pushed: `git push origin 3.59.0` — confirmed on remote"

### 1.4. The tag is the signal

The disconnect tag (cf. `operator/SKILL.md` §3.4 for the runbook) is git-observable. γ and all future agents can see it. No separate completion signal is needed — the tag appearing on main IS the proof that all gate actions completed and the cycle is disconnected.

For mid-cycle gate actions (tag push before the disconnect, branch cleanup), confirm completion to the requesting role per §1.3. But the disconnect tag itself needs no announcement — it speaks for itself.

### 1.5. BoundaryDecision: the five outcomes δ records

Per CCNF §Cell Outcomes, the closed cell at scope `n` terminates in one of four outcomes, determined by `(verdictₙ, decisionₙ)`. δ records `decisionₙ` as one of five enumerants; the outcome is the `(verdict, decision)` pair:

| `decisionₙ` (δ records) | When | CCNF outcome |
|---|---|---|
| `accept` | `verdictₙ = PASS` and the cell is transmissible to scope `n+1` without release-boundary semantics | `accepted` |
| `release` | `verdictₙ = PASS` and the cell ships as a release boundary (tag + downstream-distributable disconnection) | `accepted` |
| `override` | `verdictₙ ≠ PASS` and δ proceeds anyway by recording a structured override block (§3 below) | `degraded` |
| `reject` | δ refuses transmission for boundary-policy reasons; the cell terminates at scope `n` with no projection | `blocked` |
| `repair_dispatch` | The cell stays open at scope `n`; a child cell at scope `n` runs under a repair contract; on child accept γₙ re-emits a fresh `receiptₙ` and steps 4–5 re-fire (CCNF §Recursion Modes) | `blocked` (within-scope) |

The `(PASS, override)` and `(non-PASS, accept/release)` pairs are **invalid** per CCNF §Cell Outcomes — these are non-terminal states; δ re-decides until the outcome is one of `{accepted, degraded, blocked}`.

---

## 2. Inward membrane — γ contract → α-ready dispatch (implementation-contract enrichment)

The inward face of δ runs at dispatch time, before α is routed. γ writes the protocol-level contract (gap, ACs, oracle, evidence). γ also writes the 7-axis `## Implementation contract` section into the α prompt. **δ reviews that section before routing**, fills any unpopulated row per repo conventions (logging the enrichment in `.cdd/unreleased/{N}/gamma-clarification.md`), or blocks dispatch and escalates to operator-as-human when a row is genuinely undecidable. **Implementation-contract decisions belong to δ; α MUST NOT improvise.**

**Canonical doctrine lives at [`cnos.handoff/skills/handoff/dispatch/SKILL.md`](../../../../cnos.handoff/skills/handoff/dispatch/SKILL.md)** (Sub 3 of [cnos#404](https://github.com/usurobor/cnos/issues/404), shipped under [cnos#417](https://github.com/usurobor/cnos/issues/417)). That skill owns:

- the 7 architectural axes (Language; CLI integration target; Package scoping; Existing-binary disposition; Runtime dependencies; JSON/wire contract preservation; Backward-compat invariant);
- the dispatch-prompt template γ injects the contract block into;
- the δ review-before-routing duty and the fill-or-escalate paths;
- the four-surface mesh (γ template / δ enrichment / α constraint / β verification);
- the cnos#389/#391/#392/#393 empirical anchors.

δ's role-local realization here: δ is the actor that performs the inward review at routing time. The doctrine — what δ is responsible for, why δ holds this authority rather than γ alone, how the contract distributes across γ/α/β — lives at the canonical handoff path. The two-sided framing of δ (outward §1 + inward §2) remains the membrane's core; this section is the inward half, with doctrine sourced from cnos.handoff and execution executed here.

### 2.1. Phase 4a / Sub 3 landing note

The inward-membrane substance landed as Phase 4a of cnos#366 (cycle/397) in this file; **Sub 3 of cnos#404 (cycle/417) extracts the doctrine into [`cnos.handoff/skills/handoff/dispatch/SKILL.md`](../../../../cnos.handoff/skills/handoff/dispatch/SKILL.md)** alongside the dispatch-prompt template and the four-surface mesh declaration. This section is now a pointer; the role-local realization (δ-as-actor at routing time) stays here. **Phase 4b — harness substrate** landed in [`harness/SKILL.md`](../harness/SKILL.md) via cnos#398 (cycle/398). **Phase 4c — release-effector mechanics** landed in [`release-effector/SKILL.md`](../release-effector/SKILL.md) via cnos#399 (cycle/399). **Phase 5 — γ shrink** landed at cnos#400 (cycle/400). The complete two-sided framing (outward §1 + inward §2 of this file) is preserved; the doctrine for §2 is sourced from handoff.

---

## 3. Override — degraded boundary action

**Override is a degraded boundary action; it is never a form of validity, and it never rewrites V's verdict.** This is the load-bearing freeze from [`RECEIPT-VALIDATION.md`](../RECEIPT-VALIDATION.md) §Q4 (cnos#367) — restated here as δ-role doctrine because δ is the actor that records overrides.

### 3.1. Override does not rewrite the ValidationVerdict

An override is a structured `override:` block inside the receipt's `boundary` block, populated by δ when δ records a `BoundaryDecision` of `override` against a non-PASS `ValidationVerdict`. The `validation` block of the receipt continues to carry V's original verdict **unchanged** — override does not rewrite, mask, or replace the `ValidationVerdict`. The presence of a non-null `boundary.override` is the structural signal to every downstream consumer that the receipt closed in degraded state.

Three explicit failure modes if the design or implementation drifts:

1. **Override never rewrites the `ValidationVerdict`.** The `validation` block carries V's emitted verdict. `validation.result` stays at FAIL; `validation.failures` stays populated. δ does not edit, redact, or replace V's verdict. The verdict is what V saw; the override is what δ did with it. The two are separate fields with separate authors.
2. **Override never substitutes for PASS in downstream consumers.** A receipt where `validation.result == FAIL` and `boundary.override` is populated is **not equivalent** to a receipt where `validation.result == PASS` and `boundary.override == null`. Downstream V at scope n+1 — reading the n-receipt as evidence — must distinguish these two cases. Treating a degraded n-receipt as PASS-equivalent would let degradation propagate silently into the parent scope.
3. **Override never emits `OVERRIDE-PASS` as a `ValidationVerdict`.** V emits PASS, FAIL, and possibly WARN. There is no `OVERRIDE-PASS` value. Override lives in the `BoundaryDecision` shape (an enumerant of what δ decided), not in the `ValidationVerdict` shape (an enumerant of what V observed). Fusing them would let δ-side decisions reach back and rewrite V-side observations.

### 3.2. Downstream-consumer detection rule (binding biconditional, per cnos#367)

```text
A receipt is degraded ⇔ boundary.override != null.
A receipt is PASS-equivalent ⇔ validation.result == PASS AND boundary.override == null.

Every downstream consumer of a receipt — including:
  - the parent scope's V (V at scope n+1)
  - the parent scope's δ
  - operators auditing the cycle
  - cn-cdd-verify when reading a closed cycle
— MUST check both validation.result and boundary.override before
treating the receipt as a clean closure. Reading only one of the
two fields is incorrect.
```

The biconditional form is load-bearing. A consumer that checks only `validation.result == PASS` and assumes clean closure is incorrect (it misses degraded cells V did not approve). A consumer that checks only `boundary.override == null` and assumes clean closure is also incorrect (it misses cells where V emitted FAIL but no override exists — which is itself an incomplete-closure state per `COHERENCE-CELL.md` §Receipt Validity: "A missing override block plus a non-PASS V is not a receipt — it is an incomplete closure"). Both fields must be checked together.

### 3.3. When δ overrides

The δ-as-role may override a triad decision (or proceed past a non-PASS verdict) when:

- The triad's direction conflicts with project priorities the triad cannot see
- A role is stuck in a loop and γ's unblocking hasn't resolved it
- External constraints force a scope change (e.g. deadline, dependency shift)
- Safety or security requires immediate action
- A known pre-existing infrastructure failure (release CI red on infra unrelated to the cycle) — the explicit escape hatch named in §1.1

Override is **for information asymmetry or hard constraints, not preference**.

- ❌ "I think the implementation should use a different approach" → that's α's job
- ❌ "The review was too harsh" → that's β's judgment
- ✅ "The API we're building against is being deprecated next week" → information the triad needs

### 3.4. Override protocol

Per `cnos.cds/skills/cds/CDS.md` §"Field 6: Actor collapse rule" and `operator/SKILL.md` §3.4 (override protocol): the reassignment must name the target agent and the reason. No implicit drift.

State:

1. What you are overriding (role assignment, scope, priority, decision, or the non-PASS verdict)
2. Why (the information or constraint the triad didn't have)
3. What the new state is

When the override is against a non-PASS `ValidationVerdict`, populate the structured `override:` block in the receipt's `boundary` block with at minimum:

- `actor` — δ identity (e.g. `delta@cdd.cnos`)
- `justification` — narrative reason for the override
- `original_validation_verdict` — V's emitted verdict or a typed reference to it
- `failed_predicates_overridden` — list of the specific failed predicate refs the override covers (blanket "override all failures" is forbidden)
- `degraded_state` — `true`

Examples:

- ❌ Edit the issue quietly and let α discover the change
- ❌ Override = "PASS" silently propagated downstream
- ✅ "Override: descoping AC4–AC5 from this cycle. Reason: timeline constraint for release by Friday. New scope: AC1–AC3 only. γ to update the issue."
- ✅ "Override: release CI red on smoke test unrelated to cycle scope (known infra issue, cnos#NNN). Boundary block populated with `failed_predicates_overridden: [release_ci_smoke]`; `validation.result` remains as V emitted; downstream V at scope n+1 reads degraded matter."

### 3.5. ValidationVerdict vs BoundaryDecision — the structural distinction

The two surfaces are independent and must not be fused. Per `RECEIPT-VALIDATION.md` §"ValidationVerdict vs BoundaryDecision":

| Surface | Emitted by | Describes | Lives in receipt block |
|---|---|---|---|
| `ValidationVerdict` | V (the predicate) | Whether the receipt satisfies the contract+evidence predicates | `validation` |
| `BoundaryDecision` | δ (this role) | What δ decided to do with the receipt at the boundary | `boundary` |

The composition rule binds them but does not collapse them. Three explicit constraints follow:

1. **V does not record boundary decisions.** V emits a verdict — PASS or FAIL — and a structured failure list. V does not record "accepted" or "rejected" or "released." Those are δ's. The `validation` block of the receipt is V-owned; the `boundary` block is δ-owned.
2. **δ does not rewrite ValidationVerdict.** δ may decide to override a FAIL verdict, but δ does not rewrite the verdict to PASS. The verdict stays FAIL; the override block carries the degraded boundary action.
3. **PASS-equivalence is the conjunction.** A receipt is PASS-equivalent for downstream consumers iff `validation.verdict == PASS` AND `boundary.override == null` (§3.2). Either condition failing makes the receipt non-PASS-equivalent.

For the embedded kata exercising override (Kata B — "α stuck on AC3 due to undocumented API change"), see [`operator/SKILL.md`](../operator/SKILL.md) §9.

---

## 4. Composition with V (RECEIPT-VALIDATION.md)

Per `RECEIPT-VALIDATION.md` §Validation Interface, the composition rule:

```text
ValidationVerdict.verdict == PASS
    → δ may record BoundaryDecision in {accept, release}
    → δ may also record BoundaryDecision in {reject, repair_dispatch} if boundary-policy
      reasons orthogonal to V demand it (release freeze, downstream readiness check)

ValidationVerdict.verdict == FAIL
    → δ may record BoundaryDecision in {reject, repair_dispatch}
    → δ may record BoundaryDecision == override only as a degraded boundary action,
      with the override block populated per §3. boundary.override != null is the
      structural signal that downstream consumers MUST detect.
```

V makes acceptance **available** to δ on PASS; V does not make acceptance **automatic**. δ holds final gate authority over what crosses the boundary, and the `BoundaryDecision` is what δ records after consulting the verdict.

For the four closed-cell outcomes (`accepted`, `degraded`, `blocked`, `invalid`) and the within-scope vs cross-scope recursion modes, see [`COHERENCE-CELL-NORMAL-FORM.md`](../COHERENCE-CELL-NORMAL-FORM.md) §Cell Outcomes and §Recursion Modes.

---

## 5. What δ-as-role does NOT do

These are role boundaries at the role-skill level (CCNF kernel surface). For the dispatch-coordinator role boundaries (δ-as-actor's coordinator hat in γ=δ collapse cases), see [`operator/SKILL.md`](../operator/SKILL.md) §6.

- **δ does not produce matter.** That's α's signature (CCNF step 1). If δ wants to fix something, δ declares an override and takes the α role explicitly under the override discipline (§3).
- **δ does not discriminate matter.** That's β's signature (CCNF step 2). If δ disagrees with β's review, δ may override the boundary decision; δ does not rewrite β's review.
- **δ does not close cells.** That's γ's signature (CCNF step 3). δ does not author or edit the receipt; γ binds evidence and emits. δ reads the receipt and acts at the boundary.
- **δ does not re-read the evidence graph.** δ's signature in CCNF (step 5) explicitly excludes evidence as input. δ trusts V's verdict. If δ doubts the verdict, δ has two paths: invoke V again (reproducible from the same inputs), or `repair_dispatch` the cell to recompute its matter.
- **δ does not pronounce non-PASS receipts valid.** Override is degraded; override is not validity (§3.1).

---

## 6. Cross-references and relationships

- [`operator/SKILL.md`](../operator/SKILL.md) — δ-as-actor coordinator surface. Dispatch loop, γ=δ collapse rules, wait discipline, dispatch configurations doctrine, embedded kata, wave coordination. Cross-references `harness/SKILL.md` for dispatch / observability / polling / identity / timeout-recovery mechanics; cross-references `release-effector/SKILL.md` for release mechanics.
- [`harness/SKILL.md`](../harness/SKILL.md) — Phase 4b substrate: dispatch invocation (`claude -p` / `cn dispatch`), JSONL observability contract, worktree management, per-actor git identity (worktree-local under `extensions.worktreeConfig`), polling and wake-up forms (Monitor-wrapped transition loops, reachability re-probes), branch-retry mechanics, timeout recovery. Landed via cnos#398.
- [`release-effector/SKILL.md`](../release-effector/SKILL.md) — δ-side release mechanics: `scripts/release.sh` invocation, bare-X.Y.Z tag policy, post-push CI polling, CI-red recovery runbook, branch deletes. Landed via Phase 4c of cnos#366 (cycle/399). δ-as-role authority for the release boundary is named in §1.1 above; mechanics live in release-effector.
- [`RECEIPT-VALIDATION.md`](../RECEIPT-VALIDATION.md) — frozen typed-interface design for V, the verdict-vs-decision distinction, and the override field shape (§Q4). Binding once Phase 3 of cnos#366 lands V as a working predicate.
- [`COHERENCE-CELL-NORMAL-FORM.md`](../COHERENCE-CELL-NORMAL-FORM.md) — kernel doctrine. δ's signature at step 5 (no evidence); the four cell outcomes; the two recursion modes; the three scope-lift projections.
- [`COHERENCE-CELL.md`](../COHERENCE-CELL.md) — predecessor doctrine. δ Boundary Complex; trust grammar; the prediction of the δ split this cycle implements.
- [`cnos.handoff/skills/handoff/dispatch/SKILL.md`](../../../../cnos.handoff/skills/handoff/dispatch/SKILL.md) — canonical wire-format + 7-axis implementation-contract schema + δ-as-inward-membrane doctrine + four-surface mesh declaration. Sub 3 of [cnos#404](https://github.com/usurobor/cnos/issues/404), shipped under [cnos#417](https://github.com/usurobor/cnos/issues/417); the role-side surfaces below are its consumer-side realizations.
- [`gamma/SKILL.md`](../gamma/SKILL.md) §2.5 — γ-template side of the implementation-contract mesh (consumer realization).
- [`alpha/SKILL.md`](../alpha/SKILL.md) §3.6 — α-constraint side of the implementation-contract mesh (consumer realization).
- [`beta/SKILL.md`](../beta/SKILL.md) §Role Rules Rule 7 — β-verification side of the implementation-contract mesh (consumer realization).
- [`release/SKILL.md`](../release/SKILL.md) and [`post-release/SKILL.md`](../post-release/SKILL.md) — β-side release authoring (RELEASE.md, CHANGELOG, version decision) and γ-side post-release assessment; release-effector cross-references both as upstream context for the δ-side mechanics it executes.

---

## 7. Phase 4 of cnos#366 — what this cycle ships and what remains

This file (`delta/SKILL.md`) is Phase 4a of cnos#366. It carries:

- The two-sided membrane doctrine (outward §1, inward §2)
- Override semantics (§3) with the verdict-vs-decision freeze from cnos#367
- Composition with V (§4)
- δ-as-role boundary constraints (§5)

**Phase 4b — harness substrate (landed).** Dispatch invocation shells and observability contract, per-actor git identity with worktree-aware writes, parallel-α worktree pre-creation, polling and wake-up forms, branch-retry under push restrictions, and timeout recovery — previously authored as inline shell prose in the δ-coordinator surface — have landed in [`harness/SKILL.md`](../harness/SKILL.md) as Phase 4b of cnos#366 (cycle/398).

**Phase 4c — release-effector mechanics (landed).** The release-effector mechanics formerly living in `operator/SKILL.md` §3.4 (the single-command release script invocation, post-push CI polling, the recovery runbook, branch deletes, manual-tag prohibition) have landed in [`release-effector/SKILL.md`](../release-effector/SKILL.md) as Phase 4c of cnos#366 (cycle/399). The β-side authoring stays in `release/SKILL.md` (§2.4 CHANGELOG, §2.5 RELEASE.md, §2.5a cycle-dir move, pre-release validators); release-effector handles the δ-side mechanics only.

**Phase 5 — γ shrink (landed).** `gamma/SKILL.md` reduced to coordination + closure + triage doctrine; runtime supervision mechanics (polling shell, dispatch invocation shells, branch pre-flight bash) now cross-reference `harness/SKILL.md` and `release-effector/SKILL.md` rather than restating. Landed via cnos#400 (cycle/400).

All three substrate surfaces are now landed. δ-as-role's authority over the platform actions is named in §1.1 above; the harness substrate lives in `harness/SKILL.md`; release-effector mechanics live in `release-effector/SKILL.md`; operator-as-coordinator routing discipline (gate-action confirmation, request-vs-observation) lives in `operator/SKILL.md` §3.

---

## 8. Remote-runner delegation — δ-class effect surface

Canonical doctrine: [`docs/papers/BOX-AND-THE-RUNNER.md`](../../../../../../docs/papers/BOX-AND-THE-RUNNER.md) (cnos#425, this section's authoring cycle). Landed alongside this section under the doctrine-before-first-use rule: the essay, this skill section, the first workflow artifact, and the first instantiated receipt all land in the same commit-set; the workflow's push trigger only fires post-merge, by which point the doctrine governs.

**δ-class rule:** Any artifact an agent commits that, when present on its branch, causes another body to execute is a δ-class effect surface. δ holds authority over its authoring, not just over the receipt that closes the cell containing it.

The pattern covers, at minimum: `.github/workflows/*.yml` files; deploy hooks; extension manifests with execution permissions; scheduled-job definitions; webhook handlers; MCP server registrations; release-effector invocation scripts the agent authors rather than invokes. The class is **artifacts that cause another body to execute** — the local-shell case is one body; remote runners, deploy hooks, third-party services, and host-process extensions are other bodies. δ's authority extends to every body the agent's writes can reach, not just the body the harness can see.

The reason the rule sits in δ-as-role rather than in α-as-producer is that the effect crosses a boundary the producer cannot validate by itself. α writes the YAML; α cannot inspect the runner that will execute it. β reviews the YAML; β cannot replay the future run. Only δ-as-role, whose signature is the boundary decision, can name what authority the runner has, where it runs, and who may accept the result. The 6-field receipt is δ's surface for recording that naming.

### 8.1. Required output: the 6-field remote-runner receipt

When δ-as-role approves the authoring of any remote-runner-triggering artifact, δ requires a receipt with six fields (per `BOX-AND-THE-RUNNER.md §"The 6-field receipt convention"`):

| Field | Content |
|---|---|
| 1. Who asked | Operator session + directive, or upstream-cell receipt ref. Agent-initiated remote-runner moves without an upstream request are forbidden. |
| 2. What runs | The actual commands the remote body will execute, at command granularity (not "publishes release" — list `git tag -f X Y`, `git push --force origin X`, etc.). |
| 3. Where it runs | The execution environment (`ubuntu-latest` on GitHub Actions; a self-hosted runner; a third-party service). |
| 4. What authority | The token, scopes, secrets, and credentials the runner has access to (`GITHUB_TOKEN` with `contents: write`; deploy keys; cloud OIDC roles). |
| 5. Evidence | The artifact that proves the run happened (workflow run URL; release URL; deployment ID). May be a post-run-fillable placeholder authored before the run, with the *shape* of the evidence named in advance. |
| 6. Who may accept | The actor (operator or named cell) authorized to declare the remote run a success, with the acceptance criterion. |

Receipt path: `.cdd/unreleased/{N}/remote-runner-receipt-{handle}.md` for the issue cycle that authors the artifact. The receipt moves with the cycle into the release directory at cycle close.

### 8.2. Authoring-order rule

The receipt and the doctrine citation must land **with or before** the artifact that will trigger execution. There must be no window in which a remote-runner-triggering artifact exists on a branch that can reach an executing runner without the accompanying receipt also existing on that branch. For push-triggered workflows, the artifact and the receipt commit together; for branch-protected workflows whose dispatch is gated, the receipt may pre-land but the artifact must not pre-land alone.

Concretely, when γ dispatches a cycle that will author a remote-runner artifact, δ at the inward membrane (§2) pins this discipline into the implementation contract: the dispatch-prompt names the receipt as a required artifact, and α-side closure is incomplete without it. β verifies presence + 6-field completeness at R1 review.

### 8.3. One-shot workflows self-delete

When the remote-runner artifact is one-shot (a release-fix workflow; a single migration job; an infrastructure repair), the artifact MUST contain a final step that deletes itself from the repository after the work completes (`git rm` + `git commit` + `git push`). This is boundary discipline applied to the effect surface itself: the artifact does not persist beyond its declared single use, so the latent-execution authority closes after the run rather than remaining open. Persistent workflows (regularly-scheduled jobs, ongoing CI) do not self-delete; they are declared persistent at authoring time and the receipt's "Who may accept" field names the actor that owns the ongoing schedule.

### 8.4. Composition with the existing outward membrane (§1)

The remote-runner receipt is one piece of evidence inside the cell's full receipt at close-out. V (per §4) validates the cycle receipt; the remote-runner sub-receipt is a typed evidence ref V dereferences when the cycle's matter includes a remote-runner-triggering artifact. δ records the BoundaryDecision against V's verdict per §1.5 as normal; the override block (§3) is populated if δ accepts the cycle despite a non-PASS verdict on the remote-runner receipt (e.g. the post-run evidence is missing and the acceptor has not yet declared acceptance).

The remote-runner receipt does **not** replace V or the BoundaryDecision — it is the surface that makes the *authoring* of the effect artifact legible. The cycle-level V/δ wall still gates whether the cycle as a whole closes accepted, degraded, blocked, or invalid.

### 8.5. What this is not

- **Not a security policy.** The receipt is a legibility discipline; it does not prevent a malicious agent from forging fields. Security boundaries (token scoping, branch protection, required reviews) are orthogonal mechanisms enforced by the substrate, not by this skill.
- **Not retroactive.** A receipt authored after the remote run happened is a degraded closure; the gap is recorded as a finding and the override block (§3) is populated. Authoring the artifact first and the receipt second is a process gap ε should aggregate.
- **Not a substitute for the release-effector runbook.** The release-effector (`release-effector/SKILL.md`) governs locally-invoked release moves with release-boundary authority; the remote-runner rule governs runner-invoked moves on any boundary. The two surfaces overlap for any release-fix that uses a remote runner (cnos#425 is the first such case); the receipt cites both skills as governing.

### 8.6. First use

cnos#425 is the first cycle to author a remote-runner artifact under this doctrine. The artifact is `.github/workflows/repoint-3.82.0.yml` (one-shot; self-deleting); the receipt is `.cdd/unreleased/425/remote-runner-receipt-3.82.0.md` (6 fields filled, evidence as post-run-fillable placeholder); the effect is force-moving the `3.82.0` tag to `fd1d654e` so `release.yml` publishes the GH release with the correct root `RELEASE.md` body. The doctrine + skill section + workflow + receipt land in the same commit-set per §8.2; the post-merge push triggers the workflow under doctrine.

---

## 9. Dispatch-wake-invoked mode

This section is the **production-mode** δ contract: what δ does when a package-owned dispatch wake (per [`cnos.core/skills/agent/dispatch-protocol/SKILL.md`](../../../../cnos.core/skills/agent/dispatch-protocol/SKILL.md) §Algorithm step 6 "Launch") hands δ a claimed cell. The wake claims; δ runs the cell. Sections §1–§8 above describe δ in receipt/dispatch/override/release contexts; this section names how δ behaves when the invocation carrier is a substrate firing of a dispatch wake rather than a human operator at the loop.

Lands the forward-reference in the body of [`cnos.cds/orchestrators/cds-dispatch/SKILL.md`](../../../../cnos.cds/orchestrators/cds-dispatch/SKILL.md) (the wake's prompt body, which forward-references this section: *"δ's wake-invoked mode is the contract that lands in cnos#467 Sub 5"*). The reference dispatch wake is `cds-dispatch` (cnos#483); δ wake-invoked mode is **protocol-agnostic** at the framework level — any future dispatch wake (`cdr-dispatch`, `cdw-dispatch`) invokes the same δ contract.

### 9.1. Bootstrap-δ vs wake-invoked-δ (load-bearing distinction)

Two empirical patterns exist in the wave; the amendment names the **destination** pattern, citing the **empirical** pattern as input observation. They are not interchangeable.

| Mode | Driver | Invocation context | Empirical witness |
|---|---|---|---|
| **bootstrap-δ** | a γ-interface Claude Code session (operator at the loop); γ-the-driver also wears the δ-the-orchestrator hat in this session | δ has a parent session that holds chat-state across role spawns; γ/α/β are spawned as Agent sub-sessions of that parent | cycle/470, cycle/476, cycle/485, cycle/486 (this cycle is one) — observed at `.cdd/releases/docs/2026-06-21/{470,476}/` and `.cdd/unreleased/{485,486}/` |
| **wake-invoked-δ** | a single substrate firing of a dispatch wake; δ is the agent the wake-firing starts; **no parent session** | δ's entire input is the claimed cell's issue body + the cycle branch + `.cdd/unreleased/{N}/` tree; chat-state does NOT carry across role spawns | not yet observed in production; lands as a destination contract via Sub 5C (cnos#487) |

The two patterns share the cell-procedure surface (γ scaffolds, α implements, β reviews, δ routes), and bootstrap-δ supplies empirical evidence for what context each role needs. They diverge structurally on one axis: bootstrap-δ has a coordinating session that survives across role spawns; wake-invoked-δ does not. **Every contract clause in §9.2–§9.7 is written for wake-invoked-δ.** A clause that only works under bootstrap-δ (i.e. requires a parent session's chat-state) is a contract defect — name it as such, do not embed it.

In bootstrap-δ, γ-the-driver implicitly violates the "δ dispatches every role" invariant (§9.3 below) because γ-the-driver is also δ-the-orchestrator. The amendment surfaces this honestly; the production routing sequence separates the two roles.

### 9.2. Input contract (what the wake hands δ at invocation)

When a dispatch wake's claim sequence completes (per `dispatch-protocol/SKILL.md` §2.2 steps 1–5) and the launch step (step 6) hands the claimed cell to the matching package runtime, the package runtime invokes δ wake-invoked mode with **five named inputs**. δ MUST be able to act on each input by reading the cycle branch state and the named source files; δ MUST NOT depend on any chat-state, environment variable, or session-local memory beyond these five.

| # | Input | Semantics |
|---|---|---|
| 1 | **claimed-issue-number** | The repository issue number the wake claimed at `dispatch-protocol/SKILL.md` §2.2 step 4. δ reads this issue's body as the cell-shaped specification per [`cdd/issue/SKILL.md`](../issue/SKILL.md) §"Minimal output pattern". The body's Mode + Gap + Impact + Source-of-truth + Constraints + ACs + Proof/rejection + Cross-references fields are the cell's contract; δ does NOT improvise additions. |
| 2 | **protocol identifier** | The qualifier identifying which concrete protocol owns this cell (e.g. `cds`, `cdr`, future `cdw`). δ uses this to locate the concrete protocol's selection/lifecycle skills inside the runtime package (input #5 below). The claimed cell's `protocol:{P}` label is authoritative; the wake guarantees `P` matches its owning protocol per `dispatch-protocol/SKILL.md` §3.7. |
| 3 | **current main SHA** | The repository's `main` HEAD commit at wake-firing time. This is the base δ uses to compute the cycle branch (`cycle/{N}` from `main@<sha>`) and the SHA δ pins in `gamma-scaffold.md` for "verify-cited-sha against filesystem state" discipline. δ verifies this SHA against the working tree before dispatching γ (see §9.6 R0 step 1). |
| 4 | **wake run id** | The substrate-emitted identifier for this wake firing. δ records it in the claim trail and surfaces it in return-token comments so operators can correlate cycle progress with the substrate firing that produced it. δ does NOT decode or parse the run id — δ treats it as an opaque identifier the wake supplies. |
| 5 | **package-runtime context** | The path to the concrete protocol's skills (e.g. for `protocol:cds` → `src/packages/cnos.cds/skills/cds/` covering selection, lifecycle, etc.) plus the framework's own role-skill paths (`src/packages/cnos.cdd/skills/cdd/{gamma,alpha,beta,delta,issue,SKILL.md}`). δ exposes these as Tier-3-loadable paths to γ/α/β at spawn time; δ does NOT pre-load them itself. |

The five inputs together fully specify the cell: who (1), under what protocol (2), against what repo state (3), in what wake firing (4), with what runtime skills available (5). Anything not on this list — chat history, prior Agent sub-session memory, the operator's local environment — δ MUST NOT depend on.

### 9.3. Routing sequence (how δ dispatches γ/α/β under wake invocation)

The production invariant is: **δ dispatches every role; γ does not spawn α/β.** This is the binding routing rule. γ scaffolds (or closes out) and exits; δ dispatches α; α implements and exits; δ dispatches β; β reviews and exits; δ iterates by re-dispatching α (or γ, if a scaffold-side ambiguity surfaces).

bootstrap-δ collapses γ-the-driver with δ-the-orchestrator and implicitly spawns α/β from the γ-driver session; this is an empirical case, not the destination architecture. Per the operator directive that authored cnos#486: the production wake-invoked routing keeps γ and δ structurally distinct even when both are agent sub-processes — γ produces prompts as scaffold artifacts; δ dispatches each role from the wake-firing's δ-agent.

The discrete sequence:

1. dispatch γ for R0 scaffold — δ spawns γ with: the claimed-issue-number; the protocol identifier; the current main SHA; the cycle branch name (`cycle/{N}`); the package-runtime context (Tier-3 paths). γ's job per [`gamma/SKILL.md`](../gamma/SKILL.md) §2.5 Steps 3a–3b is to verify the SHA against the working tree, create the cycle branch from `main@<sha>`, author `.cdd/unreleased/{N}/gamma-scaffold.md` (with the per-AC oracle list, the source-of-truth table, the α prompt, the β prompt, the scope guardrails, and the friction notes), commit, push, and exit. **γ does NOT spawn α**; γ writes the α-prompt as a scaffold artifact and returns control to δ.
2. dispatch α for R[N] implementation — δ spawns α with: the cycle branch HEAD; the scaffold path; the source-of-truth list; the issue body; and, on R[N≥1], β's prior review findings appended at `.cdd/unreleased/{N}/beta-review.md §R[N-1]`. α's job per [`alpha/SKILL.md`](../alpha/SKILL.md) §2.1 (dispatch intake), §2.5 (self-coherence), §2.6 (pre-review gate), §2.7 (request review) is to implement, write `.cdd/unreleased/{N}/self-coherence.md §R[N]`, commit, push, append the review-readiness signal, and exit. **α does NOT spawn β.**
3. dispatch β for R[N] review — δ spawns β with: the cycle branch HEAD; the scaffold path; α's self-coherence §R[N]; the source-of-truth list. β's job per [`beta/SKILL.md`](../beta/SKILL.md) §"Phase map" + §"Pre-merge gate" + §"Role Rules" is to walk the AC oracle list independently, write `.cdd/unreleased/{N}/beta-review.md §R[N]` with verdict + findings, commit, push, and exit. **β does NOT spawn α.**
4. route on β's verdict — δ reads β's verdict from the branch state. If `verdict: iterate`, δ re-dispatches α at step 2 with β's findings appended; the cell stays `status:in-progress` (the wake observes via §9.5 token discipline). If `verdict: converge`, δ dispatches γ once more for closeouts (γ-closeout per `gamma/SKILL.md` §2.7) and α/β for their closeouts; then δ writes the `status:review` return token (§9.6).
5. dispatch γ for closeout — δ spawns γ with the converge signal; γ writes `.cdd/unreleased/{N}/gamma-closeout.md` (process-gap audit per `gamma/SKILL.md` §2.7 triage), and may file follow-up issues for any unaddressed friction. α and β similarly produce their closeouts per their role skills. **δ does NOT author closeout content**; δ dispatches the role that owns each closeout.

**Per-role context discipline** (no hidden state across spawns):

| Role | Spawn-context — what the role receives |
|---|---|
| γ receives | issue number, protocol id, main SHA, cycle branch name, package-runtime context (Tier-3 paths). On closeout: the converge signal + cycle branch HEAD. |
| α receives | cycle branch HEAD, scaffold path, source-of-truth list, issue body, on R[N≥1] β's findings under `.cdd/unreleased/{N}/beta-review.md §R[N-1]`. |
| β receives | cycle branch HEAD, scaffold path, α's `.cdd/unreleased/{N}/self-coherence.md §R[N]`, source-of-truth list, issue body. |

The cycle branch + `.cdd/unreleased/{N}/` tree carries everything each role needs; no role consults hidden state. **Empirical anchor:** cycle/470's R1 finding F1 (broken relative-link path in α's prompt template — `.cdd/releases/docs/2026-06-21/470/beta-review.md §R1` records the 6-vs-5 dot-segment substantive ambiguity) demonstrated that γ-scaffold-side substantive precision is necessary: without sha-pinned + path-verified scaffolding, α's R0 inherits ambiguity that β must surface as a finding. cycle/476's R1–R3 iteration (`.cdd/releases/docs/2026-06-21/476/beta-review.md`) demonstrated the same principle in the CI mechanism class: γ scaffold's CI-step preconditions must be filesystem-verified, not assumed. The wake-invoked routing applies this discipline structurally — γ verifies; α implements against verified state; β re-verifies independently.

### 9.4. R[N] iteration token discipline

The wake must be able to observe round progress without reading chat-state or asking the substrate to inspect agent memory. The discipline is: **R[N] is wake-observable via the cycle branch state**, with issue-comment signals as the human-observable secondary mechanism mapped to dispatch-protocol lifecycle transitions.

**v0 pinned mechanism — hybrid (branch-state primary; issue-comment secondary):**

| Layer | Mechanism | What it records |
|---|---|---|
| Primary | cycle branch + `.cdd/unreleased/{N}/` tree | The branch's HEAD commit reflects the latest role artifact. `.cdd/unreleased/{N}/self-coherence.md §R[N]` and `.cdd/unreleased/{N}/beta-review.md §R[N]` are the canonical per-round records. The branch state IS the iteration state — no separate ledger. |
| Secondary | issue comments on the claimed cell | Human-observable convergence signals at the four named lifecycle transitions per `dispatch-protocol/SKILL.md` §2.4. δ writes a comment at each transition (claim record; β converge; hard block; race release) so operators see cycle progress without walking the branch. |
| Not pinned for v0 | `.cdd/unreleased/{N}/*.R{N}.json` machine-readable round artifacts | Adds machinery without observable benefit at v0; the branch state + the markdown section headers already give the same observability. A future cycle may pin this if/when a downstream consumer needs machine-readable per-round metadata. |

The per-round artifact set names the wake-observable contract: at R[N], the branch carries the artifacts §9.5 enumerates; the wake confirms round-complete by reading those paths. R[N] is **not** a token δ writes to a separate ledger — R[N] is the branch state at the moment the round's artifacts are committed and pushed.

**Why this is wake-observable without chat-state:** the substrate can (and does, on every wake firing) `git fetch origin cycle/{N}` and `git ls-tree -r origin/cycle/{N} .cdd/unreleased/{N}/` to read the round-state. The substrate can read issue comments via the substrate's issue API to surface the lifecycle signals. Neither requires reading agent memory; both work from a fresh wake firing with no prior session context.

**Empirical anchors:** cycle/470 ran R0 → R1 → R2 (`.cdd/releases/docs/2026-06-21/470/self-coherence.md` + `beta-review.md` both carry R0/R1/R2 section headers — branch-state record of each round). cycle/476 ran R0 → R1 → R2 → R3 (`.cdd/releases/docs/2026-06-21/476/`); the 3-round class-trap recurrence per cnos#478 surfaced the requirement that R[N] must be wake-observable. Without R[N] visibility on the branch, the wake cannot distinguish a CI mechanism flake (mechanical) from a substantive review iteration (doctrinal). cycle/485 ran R0 only (`.cdd/unreleased/485/` — first R0-converge under the cnos#478 mechanical-injection discipline; empirical evidence that the per-CI-step audit format absorbs the bash-e class-trap class).

### 9.5. Per-R[N] artifact contract (under `.cdd/unreleased/{N}/`)

At each R[N] boundary, the cycle branch MUST carry the artifacts named below. The wake confirms a round complete by reading these paths; the matching package's release-time tooling (per [`cdd/SKILL.md`](../SKILL.md) artifact contract) lifts them into `.cdd/releases/docs/<date>/{N}/` on merge.

Seven canonical artifact classes per [`cnos.cds/orchestrators/cds-dispatch/SKILL.md`](../../../../cnos.cds/orchestrators/cds-dispatch/SKILL.md) frontmatter `wake.output.artifact_class_taxonomy`:

| Boundary | Required artifacts at `.cdd/unreleased/{N}/` |
|---|---|
| **R0** | `gamma-scaffold.md` (authored by γ; the per-AC oracle list, source-of-truth table, α prompt, β prompt, scope guardrails, friction notes) |
| **R[N≥1]** | `gamma-scaffold.md` (unchanged from R0 unless a scaffold-side ambiguity surfaces) + `self-coherence.md` with appended `§R[N]` section (authored by α; AC verification + review-ready signal) + `beta-review.md` with appended `§R[N]` verdict (authored by β; verdict + findings if any) |
| **converge** | `gamma-scaffold.md` + `self-coherence.md` (all §R[N] sections; converge-ready signal at the last) + `beta-review.md` (all §R[N] sections; final `verdict: converge`) + `alpha-closeout.md` (authored by α; cycle-level retrospective per `alpha/SKILL.md` §2.8) + `beta-closeout.md` (authored by β; review-side retrospective) + `gamma-closeout.md` (authored by γ; process-gap audit per `gamma/SKILL.md` §2.7) + **optionally** `post-release-assessment.md` (PRA — scope-dependent per the cds-dispatch manifest's `artifact_class_notes`: "only cycles with explicit retrospective value carry it") |

R[N≥1] iteration is an **internal** cell loop — the cell stays `status:in-progress` while β-α loops continue (per `cds-dispatch/SKILL.md` §"Lifecycle transitions" and `dispatch-protocol/SKILL.md` §3.6). External lifecycle states (`status:changes`) are the operator/planner's authority post-`status:review`; β iteration is NOT a `status:changes` transition.

**On converge**, after β writes the final `verdict: converge` and the three closeouts land on the branch, δ writes the `status:review` return token (§9.6).

### 9.6. Return tokens (what δ writes back to the wake)

δ writes back to the wake at four named lifecycle transitions; these map exactly onto `dispatch-protocol/SKILL.md` (cnos#454) §2.4 and [`cnos.cds/orchestrators/cds-dispatch/SKILL.md`](../../../../cnos.cds/orchestrators/cds-dispatch/SKILL.md) frontmatter `wake.responsibilities` #5. δ MUST NOT invent new transitions; δ consumes the dispatch-protocol's transition set.

| Return token | Wake-observable mechanism | When δ writes it |
|---|---|---|
| **status:review** | a *request* for the `status:in-progress → status:review` transition on the claimed cell (`cn issues fsm evaluate --issue {N} --apply`, cnos#569 Phase 2 — δ does NOT write the label directly) + an issue comment naming the cycle-PR URL | After β's `verdict: converge` lands on the branch and all closeouts are present (§9.5 converge boundary). δ opens (or updates) a cycle-PR scoped to the claimed cell, ensures `REVIEW-REQUEST.yml` and the closeout matter exist, then runs `cn issues fsm evaluate --issue {N} --apply` — the FSM applies `status:review` only when its guards pass (PR and/or commits beyond base, plus `REVIEW-REQUEST.yml`) — and writes the PR-URL comment. δ still decides *when* to request the transition; the FSM decides *whether* the guards allow it. |
| **status:blocked + reason** | label transition `status:in-progress → status:blocked` + an issue comment naming the block class (external dependency, missing precondition, infra failure, irrecoverable substrate timeout) | When δ encounters a hard block mid-cycle: a precondition the wake-invoked mode cannot satisfy (e.g. γ's SHA-pin verification fails against the working tree; a referenced source file is broken on main; the AC oracle is genuinely ambiguous in a way the scaffold did not resolve). The block-reason comment is verbose enough that the operator can decide whether to repair-and-re-dispatch or to escalate. |
| **claim release** (race) | label transition `status:in-progress → status:todo` + an issue comment naming the race detection per `dispatch-protocol/SKILL.md` §2.2 step 5 | When δ's first action — re-verify the claim is still valid — surfaces drift. This is the wake's own race-guard; δ honors it by releasing the claim without invoking γ. (In practice this fires inside the wake's claim sequence, before δ is fully spawned; the contract names it for completeness because a delayed re-verify by δ at spawn time is a defensive option.) |
| **cycle-PR URL** | included in the `status:review` transition comment; not a separate token | δ writes the PR URL into the same comment that announces the converge transition so operators / observers have a single discoverable artifact. |

**Carve-out — `status:changes` is EXTERNAL (with mechanical translator).** The dispatch wake does NOT transition out of `status:review` and δ does NOT write `status:changes` from `status:in-progress`. `status:changes` is the operator/planner's authority on review-rejection AFTER the cell shipped to `status:review` (per `cds-dispatch/SKILL.md` §"Lifecycle transitions"; `dispatch-protocol/SKILL.md` §2.4). β iteration is **internal** — the cell stays `status:in-progress` during β-α loops; only β converge advances the cell out of that state. Future readers must not conflate β iteration with external rejection; the four return tokens above are the complete set δ writes.

**Reconciliation with `cn cell return` (cnos#500).** `cn cell return --verdict iterate` (or `--verdict reject`) applies the `status:review → status:changes` transition on behalf of the operator's stated verdict. The operator is still the authority; `cn cell return` is the mechanical translator of that authority into a label transition. This does NOT change the carve-out: the dispatch wake still never writes `status:changes`; δ's four return tokens are still the complete set δ writes during a wake-invoked cycle. `cn cell return` is an operator-facing CLI command, not a wake-internal mechanism — the operator invokes it (via HI) after reading the PR; the wake is not involved. See §9.10 for the resumed-from-changes routing shape that follows.

### 9.7. v0 substrate constraints (honestly named; substrate-agnostic)

δ wake-invoked mode operates within constraints supplied by the substrate firing. δ HONORS these constraints; δ does NOT emit substrate-specific encoding for them. The renderer (per [`wake-provider/SKILL.md`](../../../../cnos.core/skills/agent/wake-provider/SKILL.md) §3.3 substrate-leakage rule + the package-vs-renderer authority split) is the surface that maps these logical constraints onto the concrete substrate (`cn install-wake` per cnos#476 + cnos#485).

| Constraint | Logical statement δ honors | Where the substrate encodes it |
|---|---|---|
| **Substrate timeout horizon** | The wake firing has a per-firing time horizon; δ assumes the cell completes (or reaches a state-stable boundary at a R[N] artifact-committed boundary) within one substrate firing. δ does NOT plan multi-firing continuations in v0; cells too large for one firing surface as `status:blocked` with a block-class comment naming the timeout class. | Renderer authority (renders the substrate's timeout field per substrate convention). |
| **Per-protocol concurrency group** | At most one wake firing per (repository × wake protocol) at a time; concurrent firings of the same dispatch wake are serialized by the substrate's concurrency mechanism (per `dispatch-protocol/SKILL.md` (cnos#454) §2.3 layer 1). δ assumes serialization; δ does NOT serialize internally. Cross-protocol firings (e.g. `cds-dispatch` and a future `cdr-dispatch`) parallelize naturally because each owns its own concurrency group. | Renderer authority (maps to the substrate's concurrency primitive). |
| **One-claim-per-firing** | One substrate firing claims at most one cell (per `dispatch-protocol/SKILL.md` (cnos#454) §3.2). δ runs the one claimed cell to a state-stable boundary; δ does NOT pick up a second cell within the same firing. | Wake claim-sequence authority (per `dispatch-protocol/SKILL.md` §2.2). |

**v0 simplifying assumption:** *the wake fires once and the cell completes (or reaches a state-stable boundary) within one substrate firing*. The state-stable boundaries are exactly the R[N] artifact-committed boundaries enumerated in §9.5 (R0 / R[N≥1] / converge). Cells that cannot reach one of these within the substrate's per-firing horizon transition to `status:blocked` with a block-class comment per §9.6; future enhancements (cron-staggered continuations, long-running cell modes) are deferred to a post-v0 wave.

### 9.8. Relationship to substrate (descriptive only — no substrate emission)

This skill is **substrate-agnostic** per the wake-provider/SKILL.md §3.3 substrate-leakage rule applied to cnos.cdd's δ skill. The §9.2–§9.7 contracts name what δ does logically; they do NOT emit substrate-specific syntax (no workflow YAML, no `runs-on:` fields, no `GITHUB_TOKEN` / secret bindings, no `${{ }}` interpolations, no `claude-code-action` references). Substrate decisions are renderer authority (cnos#476 + cnos#485 + cnos#487); δ skill MUST NOT pin them.

For v0, the substrate is GitHub Actions and the renderer is `cn install-wake`; when the dispatch wake's firing instantiates δ, the wake's prompt body (per `cds-dispatch/SKILL.md`) cites this section as the contract δ honors. A future substrate replacement (different scheduler; different runner; different agent-execution carrier) would change the renderer and the wake-template; this skill section would be unchanged.

The single descriptive carve-out in this §9 is this paragraph: it names "GitHub Actions" + "workflow" + "claude-code-action" descriptively to anchor the reader, not to pin the contract. The contract surfaces (§9.2 input contract; §9.3 routing sequence; §9.4 iteration discipline; §9.5 artifact contract; §9.6 return tokens; §9.7 constraints) are substrate-agnostic; any substrate that supplies (a) a fire-once firing horizon, (b) a serialize-per-protocol concurrency primitive, (c) a one-claim-per-firing claim mechanism, and (d) an issue-comment + label-transition API surface, can carry the contract. The framework owns the contract; the renderer owns the substrate emission.

### 9.9. Cross-references

- [`cnos.core/skills/agent/dispatch-protocol/SKILL.md`](../../../../cnos.core/skills/agent/dispatch-protocol/SKILL.md) (cnos#454) — claim sequence (§2.2), concurrency discipline (§2.3), lifecycle transitions (§2.4), drift handling (§2.6). Wake-invoked-δ consumes the claim sequence's output and honors its lifecycle transitions; the contract above maps each return token (§9.6) to a transition this skill defines.
- [`cnos.core/skills/agent/wake-provider/SKILL.md`](../../../../cnos.core/skills/agent/wake-provider/SKILL.md) (cnos#470) — substrate-agnosticism doctrine (§3.3 substrate-leakage rule); cited as the source for §9.8 above.
- [`cnos.cds/orchestrators/cds-dispatch/SKILL.md`](../../../../cnos.cds/orchestrators/cds-dispatch/SKILL.md) (cnos#483) — the reference dispatch-shape wake; its frontmatter `wake:` block is the dispatch-shape manifest. `wake.output.artifact_class_taxonomy` is the canonical artifact set §9.5 matches; `wake.responsibilities` enumerates the lifecycle transitions §9.6 consumes.
- The body of the same [`cnos.cds/orchestrators/cds-dispatch/SKILL.md`](../../../../cnos.cds/orchestrators/cds-dispatch/SKILL.md) (cnos#483) is the wake's prompt body; it forward-references this section; the reference resolved to a landed citation once cnos#486 merged.
- [`cdd/issue/SKILL.md`](../issue/SKILL.md) — cell-shaped issue contract; δ reads the claimed cell's issue body against this skill's §"Minimal output pattern" (§9.2 input #1).
- [`cdd/gamma/SKILL.md`](../gamma/SKILL.md) — γ role contract; wake-invoked routing dispatches γ per its §2.5 (scaffold) and §2.7 (closeout). δ does not restate γ-side mechanics here.
- [`cdd/alpha/SKILL.md`](../alpha/SKILL.md) — α role contract; wake-invoked routing dispatches α per its §2.1 (intake), §2.5 (self-coherence), §2.6 (pre-review gate), §2.7 (request review), §2.8 (close-out).
- [`cdd/beta/SKILL.md`](../beta/SKILL.md) — β role contract; wake-invoked routing dispatches β per its Phase map + Pre-merge gate + Role Rules.
- **Empirical anchors:**
  - `.cdd/releases/docs/2026-06-21/470/` — cycle/470 (Sub 2 wake-provider). R1 F1: broken relative-link path in α's prompt template — γ-scaffold-side substantive ambiguity. Cited in §9.3 as evidence that γ's SHA-pin + path-verify discipline is necessary for wake-invoked routing.
  - `.cdd/releases/docs/2026-06-21/476/` — cycle/476 (Sub 3 renderer v0). R1–R3 (3 rounds): missing `set -o pipefail` + `grep -c` exit-1-on-zero-matches class-trap (per cnos#478). Cited in §9.4 as evidence that R[N] iteration must be wake-observable on the branch.
  - `.cdd/unreleased/485/` — cycle/485 (Sub 5A renderer extension). R0-only converge under cnos#478 mechanical-injection discipline. Cited in §9.4 as evidence that the per-CI-step audit format absorbs the bash-e class-trap class.
- **Successor:** cnos#487 (Sub 5C) — flips `cds-dispatch` `activation_state` to `live`, renders + commits the substrate artifact for the cds-dispatch wake, and runs a real `protocol:cds` smoke cell. cnos#487 consumes this section as the contract its smoke cell verifies.

### 9.10. resumed-from-changes shape (cnos#500)

This subsection names the `resumed-from-changes` wake-invoked mode shape: what δ does when a previously-converged cell returns from `status:changes` and re-enters the cycle.

**Trigger.** The operator has read the cell's PR at `status:review`, returned an `iterate` or `reject` verdict via `cn cell return`, and re-armed the cycle via `cn cell resume`. The issue now carries `status:changes`. On the next wake firing that claims the cell (under `dispatch:cell + protocol:{P} + status:todo`, after the operator re-labels from `status:changes → status:todo` to re-authorize dispatch), δ receives a cell with a prior-round artifact history — this is the `resumed-from-changes` shape.

**How this differs from a first-claim cell.** A first-claim cell has only `gamma-scaffold.md` on its cycle branch. A resumed-from-changes cell has the full prior-round artifact set: `gamma-scaffold.md` (R0, unchanged) + `self-coherence.md` (all §R[0..N] sections, with §R[N+1] header appended by `cn cell resume`) + `beta-review.md` (all §R[0..N] verdicts) + `operator-review.md` (the HI-authored typed verdict artifact) + optionally prior closeout artifacts. δ MUST detect this shape and route accordingly.

**Detection.** δ detects the resumed-from-changes shape by reading the cycle branch state at wake-firing time. Indicators (any one is sufficient):
- `self-coherence.md` contains two or more `§R[N]` section headers (R0 plus at least one prior implementation round)
- `operator-review.md` exists at `.cdd/unreleased/{N}/operator-review.md`
- The `status:changes → status:todo` transition occurred (branch history evidence; not always available to δ)

**Input contract (resumed-from-changes).** In addition to the five standard inputs (§9.2), δ receives:
- `operator-review.md` — the HI-authored `cn.operator-review.v1` artifact carrying the operator's findings
- prior `self-coherence.md` §R[0..N] sections — α's prior-round AC verification and review-ready signal
- prior `beta-review.md` §R[0..N] verdicts — β's prior-round review findings
- `R[N+1]` increment — the next round number, determined from the `§R[N+1]` header appended by `cn cell resume` (or computed by δ if `cn cell resume` was not invoked)

**Routing sequence (resumed-from-changes).** Follows §9.3 with modifications for the resumed shape:

1. **δ dispatches α for R[N+1]** — α receives: cycle branch HEAD; `gamma-scaffold.md`; `operator-review.md`; prior `self-coherence.md` (with §R[N+1] header); prior `beta-review.md` §R[0..N]. α's job: address `operator-review.md`'s `findings[]`; append `self-coherence.md §R[N+1]` (AC verification + review-ready signal); commit, push, append readiness signal; exit. α does NOT spawn β.
2. **δ dispatches β for R[N+1]** — β receives: cycle branch HEAD; `gamma-scaffold.md`; α's `self-coherence.md §R[N+1]`; `operator-review.md`; issue body. β's job: walk the AC oracle list independently; append `beta-review.md §R[N+1]` with verdict + findings; commit, push, exit. β does NOT spawn α.
3. **δ routes on β's R[N+1] verdict** — same as §9.3 step 4: `iterate` → re-dispatch α at R[N+2]; `converge` → dispatch γ for closeout amendment.
4. **δ dispatches γ for closeout amendment** — γ appends to `gamma-closeout.md` (per `gamma/SKILL.md §2.7`) naming the R[N+1] round, the operator findings addressed, and the recovery path. γ DOES NOT rewrite prior closeout sections; γ appends an amendment section.

**γ closeout amendment shape.** On converge at R[N+1], γ's closeout amendment appends a `§R[N+1] amendment` section to `gamma-closeout.md` carrying:
- R[N+1] round summary
- Which `operator-review.md` findings were addressed (by finding id)
- Calibrated success claim for R[N+1]
- Updated triage carryforward (any new findings; any findings closed)

**Preserved invariants.** The resumed-from-changes shape MUST preserve:
- `cycle/{N}` branch — not recreated, not rebased onto main during resume (only δ may rebase per §9.3 rules; `cn cell resume` does not rebase)
- `.cdd/unreleased/{N}/` artifact directory — not replaced, not emptied
- `gamma-scaffold.md` — unchanged (R0 γ artifact; reflects the original cell contract)
- All prior `self-coherence.md §R[0..N]` sections — preserved; only §R[N+1] is appended
- All prior `beta-review.md §R[0..N]` verdicts — preserved; only §R[N+1] is appended
- `operator-review.md` — HI-owned; not modified by α/β/γ

**Reconciliation with §9.6 carve-out.** `status:changes` is still the operator's authority; the resumed-from-changes shape does not change that. What changes is the return path: the operator uses `cn cell return` (the mechanical translator of the operator's verdict) to initiate the transition, then re-authorizes dispatch by applying `status:todo`. δ's wake then claims the cell in the resumed-from-changes shape. The dispatch wake still never writes `status:changes`; the four return tokens in §9.6 are still the complete set δ writes during a wake-invoked cycle.

**Empirical anchor.** cycle/497 R1 is the first empirical witness for the resumed-from-changes shape — executed manually (bootstrap-exception path: HI applied corrections inline as `dd819f00` rather than via `cn cell return`/`cn cell resume`). The recovery sequence (HI files `operator-review.md`; α R1 takes ownership; β R1 takes ownership; γ R1 takes ownership) is the manual stand-in for the wake-invoked routing defined here. The `degraded_recovery: human_interface_applied_operator_patch` declaration in `.cdd/unreleased/497/gamma-closeout.md §5` records the bootstrap exception. cnos#500 lands this subsection; future cycles with `iterate` verdicts use the mechanical path.

**Cross-references.**
- [`cnos.cdd/skills/cdd/operator-review/SKILL.md`](../operator-review/SKILL.md) — `cn.operator-review.v1` schema; `degraded_recovery` declaration schema; HI authoring rules.
- [`cnos.core/orchestrators/agent-admin/hi-contract.md`](../../../../cnos.core/orchestrators/agent-admin/hi-contract.md) — HI behavioral contract; prohibited surfaces.
- `src/go/internal/cell/` — Go implementation of `cn cell return` (label transition) and `cn cell resume` (cycle re-arm).
- `.cdd/unreleased/497/operator-review.md` — canonical first-use `cn.operator-review.v1` witness.
- `.cdd/unreleased/497/gamma-closeout.md §5` — canonical first `degraded_recovery` declaration.

### 9.11. resumed-from-mechanical-reversion shape (cnos#630)

This subsection names a second, structurally distinct resume shape, sibling to §9.10: what δ does when a claimed `status:todo` cell carries pre-existing `cycle/{N}` matter (a branch + an open draft PR) checkpointed by the **recovery scanner's own mechanical reversion**, not by an operator rejection. This is the AC3/AC1 "resume, not defer" behavior cnos#630 requires: the partial-matter in-progress wedge fix (`transitions.json`'s `propose_status_todo_with_matter` rule; `src/packages/cnos.issues/commands/issues-fsm/scan.go`) mechanically moves a dead-run-with-checkpointed-matter cell back to `status:todo` specifically so the NEXT claim resumes from it — but only if δ actually recognizes the shape instead of treating the pre-existing branch/PR as an unexplained anomaly.

**Trigger.** A prior wake firing claimed this cell, ran to some partial point (at minimum a `cycle/{N}` branch + a checkpointed draft PR, per the cnos#591 finalizer), and its run died before `REVIEW-REQUEST.yml` was written. A later scheduled scan (`cn issues fsm scan --apply`) observed the dead run + checkpointed matter + no `REVIEW-REQUEST.yml` and mechanically reconciled `status:in-progress -> status:todo`, preserving the branch/PR and posting an audit-note comment naming this as a MECHANICAL reversion (`transitions.json`'s `propose_status_todo_with_matter` rule `reason` text, posted verbatim by `scan.go`'s existing reconciliation-comment path). The cell is now `status:todo` again, and the next dispatch-wake firing that scans the open-issue queue claims it — this is the `resumed-from-mechanical-reversion` shape.

**How this differs from §9.10's resumed-from-changes shape.** Both shapes hand δ a cell whose `cycle/{N}` branch already carries prior-round matter, but the *cause* and the *correct response* differ:

| | §9.10 resumed-from-changes | §9.11 resumed-from-mechanical-reversion (this section) |
|---|---|---|
| Cause | Operator read the PR at `status:review` and returned `iterate`/`reject` via `cn cell return` | The prior run died mid-cycle (infra failure, timeout) before reaching `status:review`; the recovery scanner's own reconciliation requeued it |
| Evidence on the branch | `operator-review.md` present; `self-coherence.md` carries `§R[N]` for N≥1; explicit rejection findings | No `operator-review.md`; the branch may carry only a partial R0 (e.g. `gamma-scaffold.md` alone, or `gamma-scaffold.md` + a partial `self-coherence.md`); the issue carries the scanner's "MECHANICAL reversion" audit-note comment, not an operator verdict |
| What "resume" means | Address `operator-review.md`'s findings; append a NEW `§R[N+1]` round | Continue the SAME round from wherever it stopped — do not re-scaffold, do not treat the partial state as a rejection to address |
| Repair-contract machinery (`cds-dispatch/SKILL.md` §Repair re-entry preflight Steps B–E) | Does not apply (no `status:changes` involved) | Does not apply — Step A of that preflight explicitly checks for and excludes this shape (`run_class: resumed_from_matter`) before ever reaching the repair-re-entry classification, precisely so this shape's absence of a rejection is not mistaken for one |

**Detection.** δ detects the resumed-from-mechanical-reversion shape from the same five-input contract §9.2 already supplies — no new Go surface or additional wake-to-δ input is required. `cn issues fsm evaluate --issue {N} --json` (already run at claim time per `cds-dispatch/SKILL.md` §"Claim mechanism") exposes `branch_exists`, `pr_exists`, and `cdd_artifacts` (`FactSnapshot`, `src/packages/cnos.issues/commands/issues-fsm/snapshot.go`) — δ reads the cycle branch state directly (`git ls-tree -r origin/cycle/{N} .cdd/unreleased/{N}/`) at spawn time, exactly as §9.4's "wake-observable without chat-state" mechanism already does for R[N] tracking. Indicators (any one is sufficient, and `cds-dispatch/SKILL.md` §Repair re-entry preflight Step A already runs this check before δ is even invoked):
- `cycle/{N}` branch and/or an open PR already exist for this issue at claim time (`branch_exists`/`pr_exists` true in the fresh `FactSnapshot`)
- an issue comment contains the scanner's "MECHANICAL reversion" audit-note phrase (the `propose_status_todo_with_matter` rule's `reason` text)
- NO `operator-review.md` is present and NO prior `status:changes` label-history exists (the negative check that rules out §9.10 instead)

**Routing sequence (resumed-from-mechanical-reversion).** Follows §9.3 with one structural difference from a first-claim cell: δ does not unconditionally dispatch γ for a fresh R0 scaffold.

1. **δ reads the existing cycle-branch state first** — before dispatching any role, δ reads `.cdd/unreleased/{N}/` on the resumed `cycle/{N}` branch. If `gamma-scaffold.md` already exists, δ does NOT overwrite it and does NOT dispatch γ for a new R0 scaffold — the existing scaffold is still the cell's contract.
2. **δ dispatches the next incomplete role** per whatever artifacts are already present, following the same completed-sections-preserved discipline named in `alpha/SKILL.md` §4 "Resumption": if only `gamma-scaffold.md` exists, dispatch α exactly as a first-claim cell would (α's own resumption protocol, §4, handles a partially-written `self-coherence.md` if one exists); if `self-coherence.md` already carries a review-readiness signal with no `beta-review.md`, dispatch β directly. δ never re-dispatches γ for a fresh scaffold when one is already present and no operator rejection occurred.
3. **From that point, §9.3 applies unmodified** — α/β iterate, β's verdict routes convergence or another α round, and closeouts land exactly as a first-claim cell's would.

**Preserved invariants.** Identical in spirit to §9.10's list, restated for this shape's cause:
- `cycle/{N}` branch — not recreated, not discarded, not rebased onto main during resume (the cnos#368 protection this cell's own fix restates for the reconciler's own action)
- the checkpointed draft PR — not closed, not replaced; `cn cell finalize`'s own idempotence keeps re-invocation safe if a later step needs to re-checkpoint
- `.cdd/unreleased/{N}/` artifact directory — not emptied; any partial R0 artifacts already present are read, not discarded
- `gamma-scaffold.md`, if present — unchanged; δ does not re-author it

**Cross-references.**
- [`cnos.cds/orchestrators/cds-dispatch/SKILL.md`](../../../../cnos.cds/orchestrators/cds-dispatch/SKILL.md) §"Claim mechanism" (cnos#630 paragraph after step 9) — the wake-side claim behavior this section's routing continues from; §"Repair re-entry preflight" Step A — the `run_class: resumed_from_matter` classification that keeps this shape from being misrouted into the repair-contract machinery.
- `src/packages/cnos.cds/skills/cds/fsm/transitions.json` — the `propose_status_todo_with_matter` rule (`in-progress` state) that produces this shape's trigger condition.
- `src/packages/cnos.issues/commands/issues-fsm/scan.go` — the recovery scanner that applies the rule and posts the audit-note comment this section's detection reads.
- §9.10 above — the sibling resume shape (operator-driven, not mechanical); read together, the two sections cover the full "cell shows up at claim time with pre-existing matter" space this framework currently names.

### 9.12. cell/substrate identity boundary (cnos#626)

This subsection is not a resume shape like §9.10/§9.11 — it names a standing doctrine constraint that holds on every wake-invoked-δ cycle, first-claim or resumed: **the cell carries no substrate identity and no hub-state scope. The wake/substrate owns identity, channel logs, and reporting; δ and the roles it dispatches (γ/α/β) do not.**

**The rule.** Cell-execution cognition — δ itself, and every role δ dispatches (γ/α/β) — operates entirely on the five inputs §9.2 already names (claimed-issue-number, protocol identifier, current main SHA, wake run id, package-runtime context). It does NOT perform a full [`activate/SKILL.md`](../../../../cnos.core/skills/agent/activate/SKILL.md) six-step activation (Kernel → CA skills → Persona → Operator → hub state → identity confirmation), does NOT load `<hub>/spec/PERSONA.md` or `<hub>/spec/OPERATOR.md`, does NOT survey hub state (dependency manifest, reflections, thread surfaces), and does NOT read or write `.cn-{agent}/` — including `.cn-{agent}/logs/`, the admin wake's channel-log writer-locality per AGENT-ACTIVATION-LOG-v0 §0. The corresponding wake-side statement of this same rule lands in [`cnos.cds/orchestrators/cds-dispatch/SKILL.md`](../../../../cnos.cds/orchestrators/cds-dispatch/SKILL.md) §"Identity and activation" (cnos#626): the dispatch wake's own minimal self-identification (which wake, which protocol) is manifest-declared — supplied by the renderer via the `{agent}` variable and bot-account binding — not file-discovered by reading Persona/Operator at runtime.

**§9.2's silence is intentional, not an oversight.** §9.2's five-input contract already omits Persona, Operator, and hub-state from what a dispatch wake hands δ at invocation — it names claimed-issue-number, protocol identifier, current main SHA, wake run id, and package-runtime context, full stop. Before cnos#626, that omission could be read either way: an intentional narrowing, or an incomplete first pass nobody had gotten around to widening. cnos#626 makes the reading explicit: the absence is doctrine, not debt. A cell that needed Persona/Operator/hub-state to do its job would be evidence the cell-execution/substrate-identity boundary itself is wrong, not evidence that §9.2's contract is missing a row. No future amendment should add a sixth "hub state" or seventh "Persona" input to §9.2's table without first overturning this subsection.

**Who owns what.** The wake/substrate (the dispatch wake's claim sequence, the FSM transitions it requests, the admin wake's channel-log and cycle-complete reporting) owns: substrate identity (which bot account, which `{agent}` binding), channel logs (`.cn-{agent}/logs/`, admin-wake writer-locality per AGENT-ACTIVATION-LOG-v0 §0), and operator-facing reporting (claim comments, return-token comments per §9.6, the admin wake's `class: cycle-complete` entry). Cell-execution cognition (δ/γ/α/β) owns: producing code, tests, docs, and reviews, and landing the `.cdd/unreleased/{N}/` artifact set — matter, not identity or reporting. A cell that finds itself needing to know "who is Sigma" or "what did the last channel entry say" has stepped outside its own boundary; that need routes to the wake/substrate layer, not into the cell's own reading list.

**Scope note (AC3/AC4 not covered here).** This subsection names the doctrine boundary (AC1) and the input-contract omission (part of AC2). It does NOT retire the mechanical enforcement that currently backstops a narrower slice of the same boundary — the `dispatch_activation_log_write_violation` write-fence (`cds-dispatch/SKILL.md` §"Disallowed surfaces") stays in place. Removing the checkout-level capability to see `.cn-{agent}/logs/` at all (rather than merely being instructed not to read/write it) is a separate, not-yet-bounded mechanism; see cnos#626's issue body and `.cdd/unreleased/626/self-coherence.md` §"AC3/AC4" for the candidate mechanisms and why they are deferred to the operator rather than implemented in this cell.

**Cross-references.**
- [`cnos.cds/orchestrators/cds-dispatch/SKILL.md`](../../../../cnos.cds/orchestrators/cds-dispatch/SKILL.md) §"Identity and activation" (cnos#626) — the wake-side prompt-body statement of this same rule; this subsection is its cnos.cdd-side doctrine counterpart.
- §9.2 above — the five-input contract whose Persona/Operator/hub-state omission this subsection names as intentional.
- [`cnos.core/skills/agent/activate/SKILL.md`](../../../../cnos.core/skills/agent/activate/SKILL.md) §2.1 — the six-step general agent-at-a-hub procedure this subsection narrows the *audience* of, without changing the procedure itself; the admin wake and any true agent-at-a-hub activation still run all six steps.
- `docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md` §0 — channel-log writer-locality; the admin wake owns `.cn-{agent}/logs/`, cell-execution cognition does not read or write it.
- cnos#467 — the two-wake admin/dispatch architecture this subsection extends one boundary further (per cnos#626's own framing).
