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
  - γ-authored dispatch prompts containing the `## Implementation contract` section (gamma/SKILL.md §2.5 Step 3b)
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
| Pre-merge gate validation | δ authorises; γ requests via `scripts/validate-release-gate.sh --mode pre-merge` (`CDD.md` §5.3b + `gamma/SKILL.md` §2.10) | `operator/SKILL.md` §3.1 |
| Push β-approved merge to main | δ executes on β's behalf when env/auth blocks β; this is execution of β's integration authority, not δ approval | `operator/SKILL.md` §3.1 |
| Release-boundary preflight | δ verifies merge commit, release artifacts, tag/deploy preconditions, platform readiness; decision is proceed / request changes / override (CDD §1.4 Phase 5a) | `operator/SKILL.md` §3.1 |
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

The inward face of δ runs at dispatch time, before α is routed. γ writes the protocol-level contract (gap, ACs, oracle, evidence) per `gamma/SKILL.md` §2.5 Step 3b. γ also writes the `## Implementation contract` section enumerating the 7 architectural axes:

1. Language
2. CLI integration target
3. Package scoping
4. Existing-binary disposition
5. Runtime dependencies
6. JSON/wire contract preservation
7. Backward-compat invariant

γ populates the rows from repo conventions and the issue body. **At dispatch time, δ reviews γ's dispatch prompt before routing it to α and ensures every row is populated.** If γ left a row unpopulated or marked "TBD," δ has two paths:

- **(a) Fill the row** per repo conventions (e.g. "this repo is Go-native, row 1 = Go"; "this repo uses the `cn` subcommand pattern for protocol-level commands, row 2 = `cn` subcommand"). Log the enrichment in the cycle's artifact channel (`.cdd/unreleased/{N}/gamma-clarification.md`, or a δ-specific channel if a later phase of cnos#366 has carved one) so the contract trail is auditable.
- **(b) Block dispatch and escalate to operator-as-human** if the row is genuinely undecidable — typically because the choice is part of the cycle's design question, not its execution shape, or because the row's value would commit the repo to a convention that hasn't been settled.

**Why this is δ's surface, not γ's alone.** γ writes what the cycle is *for* (gap, ACs, oracle, evidence). δ writes what the cycle's output is *shaped like* (language, package path, integration target, dependency footprint). The two contracts are distinct: γ's is protocol-level; δ's is implementation-level. γ knows the protocol; δ knows the repo's standing conventions. Mixing them produced cnos#389 (α improvised language because γ's prompt didn't name one and δ didn't catch the omission) and cnos#391 (α improvised package scoping and binary disposition for the same reason). cnos#392 was the first cycle where δ pinned the implementation contract at dispatch; the cycle succeeded specifically because of it.

### 2.1. The four-surface mesh

This section is the δ side of a four-surface mesh that cnos#393 shipped and Phase 4a (this cycle) extracts into its dedicated role-skill home:

- **γ template:** [`gamma/SKILL.md`](../gamma/SKILL.md) §2.5 Step 3b — the 7-axis `## Implementation contract (required for α prompt)` block; γ MUST NOT dispatch with empty rows.
- **δ enrichment:** this section — inward-membrane function; δ reviews γ's prompt; fills or escalates.
- **α constraint:** [`alpha/SKILL.md`](../alpha/SKILL.md) §3.6 — "Implementation contract is δ's, not α's"; α MUST NOT improvise; α surfaces unpinned rows before coding.
- **β verification:** [`beta/SKILL.md`](../beta/SKILL.md) §Role Rules Rule 7 — "Implementation-contract coherence"; β verifies the diff conforms to every pinned axis before APPROVE; non-conformance → REQUEST CHANGES, severity D, classification `implementation-contract`.

Each surface is locally self-justifying via the empirical anchors below; the mesh is for **discoverability** (a future role session loading any one finds the others), not for circular justification.

### 2.2. Phase 4a landing note

This section has landed here as Phase 4a of cnos#366. The precursor (cnos#393) anchored the δ-inward-membrane content at `operator/SKILL.md` §3a as a design-prerequisite for Phase 4; this cycle relocates the substance to `delta/SKILL.md` (its membrane-policy home) and reduces the operator-side anchor to a cross-reference. The two-sided framing (outward §1 + inward §2 of this file) is the membrane's core doctrine. **Phase 4b — harness substrate** landed in [`harness/SKILL.md`](../harness/SKILL.md) via cnos#398 (cycle/398). **Phase 4c — release-effector mechanics** landed in [`release-effector/SKILL.md`](../release-effector/SKILL.md) via cnos#399 (cycle/399). **Phase 5 — γ shrink** landed at cnos#400 (cycle/400), reducing `gamma/SKILL.md` to coordination + closure + triage.

### 2.3. Empirical anchor

cnos#389 (Python-not-Go) and cnos#391 (wrong package scoping + separate binary) are the failure-mode evidence that motivates this section. In both cycles γ's dispatch prompt under-specified the implementation contract; δ did not catch the omission at routing time; α improvised; β's behavior-only AC oracles APPROVE-d without catching the drift. cnos#392 was the first cycle where δ pinned the 7-axis contract at dispatch as an ad-hoc operator action; the cycle succeeded specifically because of it (the `cdd-iteration.md` F1–F4 forecast the four patches cnos#393 ships). cnos#393 made the inward function doctrine; cnos#397 (this cycle) implements δ-inward in `delta/SKILL.md`.

- ❌ δ routes γ's α prompt with rows blank — "α can figure it out"
- ❌ δ fills rows by guessing — no consultation with γ on intent, no anchor in repo conventions
- ❌ δ enriches but does not log the change, leaving the contract trail invisible
- ✅ δ reviews γ's `## Implementation contract` section row-by-row; enriches per repo conventions; logs in `gamma-clarification.md`; escalates the row to operator-as-human if undecidable
- ✅ δ blocks dispatch (does not route the α prompt) until every row is populated or explicitly escalated

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

Per `CDD.md` §1.4: the reassignment must name the target agent and the reason. No implicit drift.

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
- [`gamma/SKILL.md`](../gamma/SKILL.md) §2.5 Step 3b — γ-template side of the implementation-contract mesh.
- [`alpha/SKILL.md`](../alpha/SKILL.md) §3.6 — α-constraint side of the implementation-contract mesh.
- [`beta/SKILL.md`](../beta/SKILL.md) §Role Rules Rule 7 — β-verification side of the implementation-contract mesh.
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
