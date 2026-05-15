# Coherence Cell

**Status:** Draft refactor doctrine
**Owns:** Recursive coherence-cell model for CDD redesign
**Does-Not-Own:** Current executable CDD algorithm. `src/packages/cnos.cdd/skills/cdd/CDD.md` remains the canonical CDD algorithm. This document does not replace, does not supersede, and does not silently override `CDD.md`. When this doctrine and `CDD.md` disagree about *what CDD does today*, `CDD.md` wins. When they disagree about *what CDD should become*, this document is the proposal surface, not the ratified state.

---

## Thesis

CDD is a recursive coherence-cell protocol.

```text
A child cell receives a contract.
α produces matter.
β discriminates matter.
γ closes the cell and emits a receipt.
V validates the receipt against the contract.
δ transports, accepts/rejects, and effects boundary actions.
Accepted receipts become matter at the next scope.
ε observes receipt streams to evolve the protocol.
```

The deep invariant:

```text
A closed cell is not trusted because a higher role approved it.
A closed cell is trusted because its receipt validates against its contract.
```

That is the load-bearing rephrase. CDD today is described accurately but flatly as *"CDD has α, β, γ, δ, ε roles"*. The accurate but operational description is *"CDD is a recursive coherence-cell protocol where validity, not seniority, is what makes a closed cell transmissible to the next scope."*

---

## Structural Prediction: Roles, Runtime, Validation, Release

The doctrine this document carries is structural, not metaphorical. The biological / cellular / organismic framing — roles as biological functions, runtime as substrate, validation as cell wall / receptor assay, release as boundary effection — is permitted as expository scaffolding and does not appear in the AC checks as load-bearing wording. What is load-bearing is the four-way separation predicted by the model:

| Surface              | Function                                                       | Owns                                                                  |
|----------------------|----------------------------------------------------------------|-----------------------------------------------------------------------|
| **Roles**            | α/β/γ/δ/ε as biological functions of the cell                  | Production, discrimination, closure, boundary, evolution              |
| **Runtime substrate**| Harness and platform driver                                    | Dispatch mechanics, polling, git/CI/shell mechanics, timeout recovery |
| **Validation**       | Trust boundary                                                 | `V(contract, receipt, evidence)` predicate / capability               |
| **Release / boundary effection** | Effector for irreversible platform actions         | Tag, deploy, external commit, spend, money                            |

These four are separate surfaces. They must not be fused into one skill. A skill that says "δ holds the boundary, dispatches roles via `claude -p`, polls CI, retries on timeout, runs the release driver, and validates the receipt" is a fused skill — it collapses all four onto one role doctrine. It is not one skill. The structural prediction is that fusing role function with runtime substrate, validation predicate, and release effector is a design smell — surface smearing per `src/packages/cnos.core/skills/design/SKILL.md §3.7`.

The biological metaphor is one expository route to that prediction. The doctrine survives without the metaphor; the structural prediction survives. If the metaphor is removed, the four separate surfaces still hold; if the metaphor is kept, it must translate into structural design predictions, not decorate them.

This is the AC2 invariant. The check is not "does the doc say 'biological functions'?" — that wording may appear or not. The check is "does the doc make the structural prediction that role, runtime substrate, validation, and release/boundary effection are separate surfaces that must not be fused?"

---

## The Cell

A coherence cell is the unit of work CDD recognises. A cell:

1. Receives a **contract** — the issue body plus its embedded acceptance criteria, scope, non-goals, related artifacts, and active design constraints. The contract is what α/β/γ are accountable to. Nothing else.
2. Produces **matter** — what α builds against that contract (code, prose, doctrine, schemas, depending on the protocol instantiation; for CDD: source/tests/docs).
3. Carries a **review record** — β's verdict, finding disposition, round-by-round transcript on the cycle branch.
4. Closes — γ records triage, dispositions, and a **receipt** that describes the cell to the next scope.
5. Survives validation — `V(contract, receipt, evidence)` returns a verdict. If PASS, the closed cell is transmissible. If FAIL, the closed cell is not transmissible without an explicit degraded-boundary override.
6. Is accepted at the boundary — δ accepts the receipt and effects the boundary action (tag, release, deploy), or rejects, or repair-dispatches.
7. Becomes matter at the next scope — once accepted, the receipt itself is α-level matter for the parent scope. The closed cell is no longer "a project that happened"; it is "a typed object that the parent scope reasons over."

This recursion is the second load-bearing claim. CDD does not stop at "we closed the cycle." It stops at "we closed the cycle and emitted a receipt that some larger frame can now treat as one of its inputs."

---

## Roles as Cell Functions

α, β, γ, δ, ε are role functions, not job titles or positions. Each names a biological function the cell performs to remain coherent:

- **α** — produces matter against the contract.
- **β** — discriminates matter against the contract. β is the cell's immune discrimination.
- **γ** — closes the cell. γ collects matter, review record, and triage into a receipt.
- **δ** — boundary complex. δ transports contracts outward, receipts inward, and effects irreversible platform actions (tag, release, deploy). δ is not internal to the cell's reasoning.
- **ε** — protocol evolution. ε observes accepted receipt streams across cells and patches CDD doctrine, schemas, validators, and harness contracts when repeated process incoherence appears.

Roles being functions, not positions, has consequence: any agent can wear any hat across cells. Two hats may collapse onto one actor *within* a cell when the independence guarantee survives. The one collapse that is never safe is α=β within a single cell — for the reason given in *β Independence as Contagion Firebreak* below. This restates `ROLES.md §4` (hats-vs-actors) for the cell model.

---

## Receipt Validity

A receipt is the parent-facing artifact of a closed cell. It encodes:

- the contract the cell was accountable to
- the production record (α actor, artifacts, commits)
- the review record (β actor, verdict, findings, finding dispositions, merge evidence)
- the closure record (γ actor, triage, unresolved debt, next commitments, protocol_gap_count, protocol_gap_refs)
- the validation record (validator identity, validator version, result, failures, warnings)
- the boundary record (δ actor, acceptance, release tag, override block if any)

A receipt is **valid** when:

```text
V(contract, receipt, evidence) = PASS
```

A receipt is **degraded** when V fails and δ overrides anyway. A degraded receipt is still emitted — but the override block must be populated, the failed predicates named, and the boundary action recorded as degraded. A missing override block plus a non-PASS V is not a receipt — it is an incomplete closure, and the cell is not closed.

The parent scope reasons over receipts, not over the underlying matter. The receipt is the cell's surface. The cell's internal state — α's thinking, β's transcript, γ's triage notes — is below the receipt surface. The receipt is what crosses scopes.

---

## Validator Position: V, Not ζ

Validation is `V`, a typed predicate / capability invoked by δ. It is **not** a sixth Greek role by default.

```text
V : Contract × Receipt × EvidenceGraph → ValidationVerdict
```

`V` is a pure function in the doctrine sense — given (contract, receipt, evidence), it returns a verdict. It does not hold longitudinal context. It does not maintain its own observation order. It does not introduce ζ as a settled sixth role. δ invokes `V` and gates the boundary action on the verdict.

Trust is **not** issued by γ, δ, or ε. Trust is established by the relation:

```text
contract + evidence + valid receipt
```

If `V(contract, receipt, evidence) = PASS`, the cell is trusted by the parent scope. If γ approves but `V` fails, the cell is not trusted — γ's approval is necessary but not sufficient. If δ overrides a non-PASS verdict, the override is a degraded boundary action and must be receipted as such (see *Trust Boundary* and *δ Boundary Complex* below).

This positioning is deliberate. Introducing ζ as a settled sixth role would treat the validator as an agent with discretionary authority and new observation order. That is a different design — one that may eventually become necessary if `V` needs longitudinal state or judgment beyond predicate evaluation. For now, the validator is a predicate/capability, not a role. If a later cycle proves a sixth role is necessary, the doctrine can name it. Until then, ζ is not a settled doctrine and is not introduced here.

---

## Recursion Equation

The cell recursion equation in symbol form:

```text
matterₙ      := αₙ.produce(contractₙ)
reviewₙ      := βₙ.review(contractₙ, matterₙ)
closed_cellₙ := γₙ.close(contractₙ, matterₙ, reviewₙ)
receiptₙ     := closed_cellₙ.receipt
verdictₙ     := V(contractₙ, receiptₙ, evidenceₙ)
ACCEPTEDₙ    := δₙ.accept(receiptₙ, verdictₙ)
αₙ₊₁         := ACCEPTEDₙ.receipt_as_matter
```

Reading this row-by-row:

- α at cell n produces matter against the cell-n contract.
- β at cell n reviews the contract-matter pair.
- γ at cell n closes the cell, fusing contract, matter, and review into a closed-cell record.
- The receipt is the parent-facing surface of the closed cell.
- `V` is invoked on the triple (contract, receipt, evidence) and emits a verdict.
- δ at cell n accepts the receipt (or rejects, or repair-dispatches) based on the verdict.
- The accepted receipt at cell n becomes α-level matter at cell n+1.

The last line is what makes the model recursive. The next scope does not re-read the cell's internal artifacts to know what happened — it reads the receipt. Cells stack.

---

## Trust Boundary

The trust-boundary judgment in inference-rule form:

```text
V(contractₙ, receiptₙ, evidenceₙ) = PASS
────────────────────────────────────────
contractₙ ⊢ receiptₙ : accepted matter
```

Read: when `V` returns PASS, the contract entails the receipt as accepted matter for the next scope.

The judgment is a typed boundary, not a discretionary hand-wave. δ may override the boundary, but the override is not validity — it is a degraded boundary action and must be receipted with the override block populated. The parent scope reading a degraded receipt knows it is reading degraded matter.

Trust is established by the relation `contract + evidence + valid receipt`. That triple is what holds. Trust is not issued by γ closing the cell, by δ tagging the release, or by ε approving the protocol. Those are role actions. Validity is what carries trust across the boundary.

---

## β Independence as Contagion Firebreak

α≠β is not bureaucracy. α≠β is the contagion firebreak.

β is the cell's immune discrimination. β observes the matter α produced from outside α's frame — the essential property that makes review add information rather than repeat α's internal validation. When β and α are the same actor in the same cell, β cannot observe α's matter from outside α's frame; the matter passes through "review" without ever being discriminated.

A receipt without independent β review is transmissible only as degraded matter. The validator can structurally check that the β actor identity differs from the α actor identity, that β-authored review findings exist, that β's verdict precedes merge, and that the merged matter is the matter β verdicted on. The validator cannot prove that β's reasoning was truly independent of α's — but it can reject obvious counterfeit receipts where the structural evidence of independent review is absent.

If a cell ships a receipt where α and β are the same actor, the receipt is immunologically compromised. The cell may still close — but the receipt must carry that fact as degraded matter, and the parent scope must know it is reading degraded matter. The contagion firebreak is the mechanism by which a defect in cell n's matter does not silently propagate as cell n+1's α-input.

This is operative immune doctrine, not loose role-style preference. The cell is the unit of containment; β is the discrimination mechanism; without β, the cell does not contain.

---

## Review-Independence Evidence

The validator cannot prove β's semantic independence completely — that would require reading β's reasoning and judging its frame. The validator can reject obvious counterfeit receipts using **structural** evidence. The receipt's review block carries fields the validator can check mechanically:

| Field                          | Structural check                                                                  |
|--------------------------------|------------------------------------------------------------------------------------|
| `alpha_actor` vs `beta_actor`  | `alpha_actor != beta_actor` (canonical role identities; same email = fail)         |
| `alpha_commit_authors`         | All α-attributed commits in the review window were authored by `alpha_actor`       |
| `beta_review_authors`          | All β-attributed review artifacts were authored by `beta_actor`                    |
| `beta-review.md` verdict       | Verdict is APPROVED with finding-disposition rows for every binding finding        |
| Reviewed artifact refs         | β's verdict names the artifact paths and SHAs that were reviewed                   |
| Verdict precedes merge         | β's APPROVED verdict commit precedes the merge commit (timestamp / ancestry)       |

These checks are about whether the *structural* evidence of independent review is present. They cannot detect collusion or hidden coordination. A receipt that passes all six checks may still describe a cycle in which α and β colluded in a back-channel. The validator can flag the cell as immunologically uncompromised at the structural level, but it cannot certify it as semantically independent. Validator design acknowledges this and frames the checks as rejection criteria for counterfeit receipts, not as proof of semantic independence.

This is the AC9 invariant. The validator can reject obvious counterfeit, not prove independence completely; both halves matter.

---

## δ Boundary Complex

δ is **not** merely a passive membrane. δ is a boundary complex with at least three internal functions:

1. **Membrane / boundary policy** — what the cell exports outward and what the parent accepts inward. δ holds gate authority over what crosses the boundary.
2. **Transport** — δ transports contracts outward to the cell (dispatching α/β/γ to operate on the contract) and receipts inward from the cell (carrying closure records to the parent scope).
3. **Effector for irreversible platform actions** — tag, release, deploy, external commit, spend, money. The cell's matter is reversible internally; the effector's output is not.

δ invokes validation; δ does not embody validation. An override is not validity. It is a degraded boundary action and must be receipted as such with the override block populated.

The exclusion polarity (AC10 invariant): δ **must not contain runtime substrate mechanics**. Specifically, the following do not belong in δ role doctrine — they belong below δ in the harness/platform-driver substrate:

- `claude -p` invocation
- `cn dispatch` mechanics
- `git config` setup for role identities
- `git fetch` polling and the wake-up loop
- branch retries on push 403
- timeout recovery procedures
- CI polling between releases
- shell mechanics (worktree setup, temp directories, signal handling)

These belong harness — they belong platform driver — they belong below δ — they belong below delta — they belong runtime substrate. They are outside δ role doctrine. They are outside delta role doctrine. The doctrine document does not contain runtime substrate; runtime substrate must not be present as role doctrine. Naming these mechanics as substrate that belongs below δ is the structural prediction that AC10 checks; treating them as δ role doctrine is the failure mode AC10 catches.

Release execution belongs to a release driver / boundary effector surface, not to the whole δ role doctrine. δ holds *the policy* for what crosses the boundary; the *mechanics* of doing so live in a separate runtime substrate layer that δ invokes. Today, `operator/SKILL.md` fuses both. The refactor split is deferred; this document names the fusion and predicts the split.

---

## Runtime Substrate and Harness (excluded from δ role doctrine)

The runtime substrate is the layer below δ that carries out the mechanical actions δ's boundary policy decides on. It contains:

- session dispatch: `claude -p`, `cn dispatch`, sub-agent invocation, identity-rotation primitive
- branch substrate: branch create / fetch / push / retry, `git config user.email` setup, worktree setup, `git fetch` polling, transition-only stdout filter, Monitor wake-up
- platform substrate: shell mechanics, signal handling, timeout recovery, CI polling, release script invocation
- harness substrate: per-environment quirks (push 403 fresh-branch chains, MCP-only γ environments, sandbox transport flake)

These are real, necessary, and ship in the repo today (mostly inside `operator/SKILL.md` and `gamma/SKILL.md`). The doctrine question is *where they live*. They live below δ, not inside δ role doctrine. The refactor target is a thin δ skill that names boundary policy and an explicit runtime-substrate surface (harness / platform driver) that names the mechanics. The exclusion polarity is explicit: runtime substrate mechanics belong below δ in the harness/platform-driver substrate.

Once the substrate is named explicitly, δ's role doctrine shrinks to its actual function: invoke validation, accept/reject/repair-dispatch receipts, transport contracts/receipts, effect irreversible platform actions. Everything else is substrate.

---

## Release as Boundary Effection

Release is a boundary action, not a role's process. Tagging the version, pushing the tag, publishing the GitHub release, deploying — these are irreversible effector actions on a platform δ has authority over but does not own. The release driver is the substrate that executes those actions on δ's authorisation; δ holds the boundary policy that decides when to authorise.

Today, `release/SKILL.md` and `operator/SKILL.md` both carry pieces of this. The refactor target is to name release as boundary effection and to relocate the mechanics (`scripts/release.sh`, CI polling, tag-message generation, branch cleanup) to the substrate where they belong. δ's role doctrine then states the *policy*: when is the release authorised, what gates must pass, what override is allowed. Not the *mechanics*.

This is the same structural prediction as *δ Boundary Complex*, applied to one specific class of boundary action.

---

## cn-cdd-verify Target Position

`cn-cdd-verify` today is an artifact-presence checker. It verifies that the canonical CDD cycle artifacts exist: `self-coherence.md`, `beta-review.md`, `gamma-closeout.md`, `alpha-closeout.md`, `beta-closeout.md`, the cycle directory structure, the release directory move. It does **not** validate semantic truth, review independence, receipt consistency, or contract satisfaction. The current behavior is artifact completeness; the future behavior is contract/receipt validation.

The target position for `cn-cdd-verify` is the cell-wall validator — the operational implementation of `V`. In the target state:

```text
V : Contract × Receipt × EvidenceGraph → ValidationVerdict
```

is implemented by `cn-cdd-verify`. It reads the receipt, reads the contract (issue body + AC oracles), reads the evidence graph (the cell's `.cdd/` artifacts plus the diff), and emits a verdict. The verdict has structured failure / warning lists. δ invokes it before tag and gates the release on PASS or on a populated override.

The five existing CDD files (`self-coherence.md`, `beta-review.md`, `gamma-closeout.md`, `alpha-closeout.md`, `beta-closeout.md`) may remain as evidence sources — they carry the in-cell record. The parent-facing trust surface, however, becomes the typed receipt: a single artifact that summarises the cell at the validator's interface. The existing per-cycle markdown artifacts continue to exist as evidence-graph inputs the validator reads when scoring; the typed receipt is what crosses scopes.

This document does not implement that transition. It names the target so the implementation issue (deferred) has a stable doctrine to refer back to.

---

## Illustrative Receipt Shape (non-normative)

The following is an illustrative receipt frontmatter sketch. It is **illustrative**, not normative. It is an **example**, not a canonical schema. CUE schemas (`receipt.cue`, `contract.cue`) are **deferred** to a follow-up issue; this sketch does not define a schema. The schema name below is an illustrative placeholder; this document does not pin a canonical version.

```yaml
# Illustrative; not normative. This sketch is an example only.
# CUE schema is deferred to a follow-up issue; this document does not define a schema.
schema: cnos.cdd.receipt.example
contract:
  issue: 364
  branch: cycle/364
  base_sha: d412a1e9
  cdd_version: 3.76.0
  acceptance_oracle_ref: https://github.com/usurobor/cnos/issues/364
production:
  alpha_actor: alpha@cdd.cnos
  artifacts:
    - path: src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md
      sha: <implementation-SHA>
      purpose: draft refactor doctrine
    - path: src/packages/cnos.cdd/README.md
      sha: <implementation-SHA>
      purpose: package source-map pointer
  commits:
    - sha: <implementation-SHA>
      author: alpha@cdd.cnos
review:
  beta_actor: beta@cdd.cnos
  verdict: approved
  rounds: 1
  findings:
    - id: B1
      severity: binding
      disposition: fixed
      evidence_ref: .cdd/unreleased/364/beta-review.md#round-1
  merge_sha: <merge-SHA>
closure:
  gamma_actor: gamma@cdd.cnos
  triage_ref: .cdd/unreleased/364/gamma-closeout.md
  unresolved_debt:
    - "deferred: CUE schema design (receipt.cue, contract.cue)"
    - "deferred: cn-cdd-verify receipt-validation implementation"
    - "deferred: δ/operator skill split"
    - "deferred: γ skill shrink"
    - "deferred: ε relocation"
  next_commitments:
    - "follow-up issue: receipt.cue / contract.cue schema design"
  protocol_gap_count: 0
  protocol_gap_refs: []
validation:
  validator: cn-cdd-verify
  validator_version: <future>
  result: pass
  checked_at: <future-timestamp>
  failures: []
  warnings: []
boundary:
  delta_actor: delta@cdd.cnos
  accepted: true
  release_tag: <future>
  override: null
```

The fields above are illustrative of the shape, not of the canonical schema. A future cycle will design `receipt.cue` and `contract.cue`, and at that point the schema name will be pinned. Until then, treat this sketch as a sketch.

---

## γ and δ Managerial-Residue Sweep

γ and δ are the main surfaces where managerial idioms can leak back into role doctrine. *Manager* idioms — *monitor*, *supervise*, *oversee*, *manage* — sound legitimate because coordination is legitimate. But coordination is not management: coordination produces an artifact (a triage record, a dispatch prompt, a clarification entry), a receipt (the typed cell receipt), or a boundary decision (proceed, gate, override). Management without an artifact, receipt, or boundary decision is managerial residue: it leaks process supervision into role doctrine where role doctrine should describe what the role *produces*.

The sweep rule:

```text
If a skill says "monitor", "supervise", "oversee", or "manage",
ask what artifact, receipt, or boundary decision that verb produces.
If none, it is probably managerial residue.
If yes, rename the verb to the actual biological function:
  observe
  discriminate
  route
  validate
  close
  transport
  release
  repair-dispatch
```

Coordination is legitimate — γ coordinates the α↔β loop and that coordination is one of the cell's biological functions. Runtime supervision mechanics (polling intervals, wake-up loops, branch retries) are also legitimate — they belong to the substrate. The sweep rule predicts that fusing the two surfaces into one role skill — γ-as-manager — produces L4 (or L12 in larger surfaces) of managerial residue inside role doctrine. The fix is to keep γ's role doctrine on the coordination *function* (what gets produced) and to move the supervision *mechanics* (how it's polled, when it retries) into substrate.

The same sweep applies to δ. δ's biological function is boundary policy, transport, and effection. δ's role doctrine that says "δ polls CI for green, retries on red, manages timeout recovery, supervises release workflow completion" carries managerial residue: those are substrate mechanics describing how δ's boundary decision is enacted, not what δ's boundary decision *is*. Apply the sweep rule and the residue separates from the doctrine.

---

## ε as Protocol Evolution

ε is not metabolism. ε is evolution over receipt history.

Ordinary in-cell metabolism is α producing, β reviewing, γ closing, δ effecting. That is what every cell does. ε observes the cell from outside, across many cells — ε observes accepted receipt streams and patches the protocol when repeated process incoherence appears.

ε is not a required in-cycle actor for every cell. A cell may close cleanly with no ε work. ε is triggered when receipt-stream patterns reveal a protocol gap: a recurring class of finding the protocol does not prevent, a skill that consistently fails to prevent a predictable error, a coordination burden that signals a better mechanical path. When such patterns appear in the accepted receipt stream, ε patches CDD doctrine, schemas, validators, or harness contracts.

ε observes *outside the cell*. Inside the cell, the roles are α/β/γ/δ. ε's matter is the protocol itself — the receipt schema, the role skills, the validator design, the harness contracts. ε's contract is the receipt history. ε's review is harder: ε is not in a single cell where β can discriminate, so ε's findings often pass through γ-triage and become MCAs that the next cycle implements.

The doctrine consequence is that ε must not be conscripted as ordinary metabolism. A cell that closes cleanly does not need ε. A cell whose receipt stream reveals a protocol gap is the trigger for ε work, but ε's output lands as protocol-level matter, not as cell-level matter. Whether ε relocates to `ROLES.md`, `cnos.core`, or a new protocol-iteration package is an Open Question (below); the structural prediction here is that ε is outside ordinary cell metabolism and should not be required to write a per-cycle artifact for cells with no protocol gaps.

---

## ε Artifact Rule (target doctrine)

Every receipt records:

```yaml
closure:
  protocol_gap_count: <int>
  protocol_gap_refs: [<refs>]
```

This makes "no protocol gaps" a positive signal carried by every receipt — not an absence requiring no artifact. When `protocol_gap_count > 0`, a separate `cdd-iteration.md` artifact is required and `protocol_gap_refs` points to it. When `protocol_gap_count == 0`, no separate iteration artifact is required.

This collapses two failure modes:

- The pre-receipt design forced ε to write `cdd-iteration.md` every cycle as ordinary metabolism, even when no protocol gap existed. That over-uses ε and makes the iteration artifact noise.
- A design that removed protocol-gap evidence entirely loses the positive signal that the cycle observed no gaps.

The doctrine: every receipt carries `protocol_gap_count` and `protocol_gap_refs`. Separate `cdd-iteration.md` is required only when `protocol_gap_count > 0`. The iteration artifact is a real artifact when it exists; the receipt is the always-present record of whether one is required.

---

## Practical Landing Order

This is the deferred roadmap. Each item is a future work item, named here so that follow-up issues have a doctrine to refer back to. **None of these are required by this issue.** They are future work, scoped explicitly as deferred roadmap.

1. **Add `COHERENCE-CELL.md`** — this document. (This cycle's matter.)
2. **Add `receipt.cue` and `contract.cue`** — typed CUE schemas for receipt and contract. Deferred to a follow-up issue.
3. **Refactor `cn-cdd-verify`** — from artifact-presence checker to contract/receipt validator. Implements `V(contract, receipt, evidence) → verdict`. Receipt-validation implementation deferred.
4. **Split δ into boundary skill plus harness/platform and release surfaces** — relocate runtime substrate (dispatch, polling, git, CI, release mechanics) out of `operator/SKILL.md` into a harness/platform-driver surface; rename `operator/SKILL.md` to `delta/SKILL.md` and shrink to boundary doctrine. Future work. δ split deferred.
5. **Shrink `gamma/SKILL.md`** — to issue quality, artifact coordination, closure, and triage. Relocate runtime supervision mechanics (polling, wake-up loops, branch preflights) to the harness surface. Future work. γ shrink deferred.
6. **Move ε up or declare it as protocol-iteration doctrine above ordinary CDD metabolism** — relocate ε into `ROLES.md`, `cnos.core`, or a new protocol-iteration package; remove ordinary-metabolism framing from cycle-level requirements. Future work. ε relocation deferred.

This landing order does not commit to a schedule. It commits to a structural sequence: doctrine before schema, schema before validator, validator before role splits, role splits before ε relocation. Each step depends on the previous step's stable surface. None of this implementation work is required by this issue.

---

## Open Questions

The following five Open Questions are live design tensions surfaced by this doctrine. They are **not resolved** here. They are seeded for next-cycle inheritance. A follow-up issue or cycle may resolve any one of them; this doctrine deliberately leaves them open.

1. **When does V fire?** Pre-merge (δ accept gates the merge) or post-merge (δ accept gates release/tag)? Possibly twice (γ preflight + δ authoritative). The validator predicate could be invoked at multiple boundaries; the question is which boundaries are gated.

2. **Is V a capability or a command?** A doctrine-level predicate / capability invoked by δ from within its existing flow, or a separate command/provider with its own registry that δ dispatches to? The capability framing keeps `V` as a function call; the command framing makes `V` a first-class runtime surface with its own contract.

3. **Where does ε relocate?** `ROLES.md` (alongside the generic ladder), `cnos.core` (as a core surface), or a new protocol-iteration package above `cnos.cdd`? Each location implies a different scope for ε's authority.

4. **How is an override receipted?** Structured `override` block in the receipt schema (actor, reason, failed_predicates, risk, expiry, follow_up), a sibling artifact (e.g. `override.md`), or a degraded-state flag on the validation block? The structured block is leaner; the sibling artifact preserves narrative; the flag is the minimum signal.

5. **Do existing per-role closeouts become evidence-graph inputs to V, or remain parent-facing in their own right?** The lean reading is: closeouts become evidence-graph inputs and the typed receipt is the parent-facing artifact. The alternative is closeouts remain parent-facing and the receipt is a thinner summary. The choice affects what the validator reads and what the parent scope sees.

These are inheritance for the next cycle. Resolving any of them here would close design space that the protocol may want to keep open until a real implementation issue forces the choice.
