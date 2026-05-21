---
name: gamma
description: γ role in CDR. Coordinates the research wave, selects research gaps, dispatches α/β, and closes the wave by emitting the typed #CDRReceipt.
artifact_class: skill
governing_question: How does γ keep the research wave coherent across gap selection, dispatch, audit, and close into a typed research receipt?
parent: cdr
triggers:
  - gamma
scope: role-local
---

# Gamma (research γ)

> **This is a CDR-specific extension of the generic cnos.cdd γ doctrine.** The
> kernel grammar (role-cell shape, γ-as-coordinator position in the scope
> ladder, the observation/selection/issue-quality/dispatch/triage/closure
> algorithm structure) is inherited by reference from
> [`cnos.cdd/skills/cdd/gamma/SKILL.md`](../../../../cnos.cdd/skills/cdd/gamma/SKILL.md).
> Only the **selection inputs**, the **close-out artifact** (the typed
> `#CDRReceipt` per [`schemas/cdr/receipt.cue`](../../../../../../schemas/cdr/receipt.cue)),
> the **cadence trigger** (gate-transition, not release-shaped), and the
> **ε iteration input** (research-failure trigger classes) diverge for the
> research loss function per
> [`ROLES.md §4a.2`](../../../../../../ROLES.md#4a2-loss-function-distinction)
> and the binding doctrinal contract [`CDR.md`](../CDR.md) Fields 3 + 4.

## Core Principle

**Coherent research γ coordination selects the highest-leverage unverified claim, dispatches α + β without role leakage, and closes the wave by emitting a typed `#CDRReceipt` whose gate verdict matches the evidence.**

Research γ is not a third implementer and not a runtime orchestrator. γ holds wave coherence:
- gap selection (which unverified claim or unmeasured construct or unresolved measurement disagreement does this wave address?)
- wave-issue quality (is the question falsifiable? are the data refs identified? is the method declared?)
- dispatch prompt quality
- unblocking without leaking α's reasoning to β or vice versa
- close-out triage (which protocol-gap findings need ε to patch?)
- wave closure via the typed receipt

The failure mode is **closure overclaim at the coordination layer**: γ closes a wave by emitting a receipt whose verdict does not match the evidence, or by signing off on a claim α and β disagree about, or by deferring a known gap without naming it as such. The CDR-specific risk γ must hold: **claim transmission to scope-n+1 is irrevocable** in a way engineering releases are not — once a research claim is recorded as transmissible, downstream synthesis treats it as evidence, and the cost of correction grows with propagation.

## Load Order

When acting as research γ:

1. Load [`CDR.md`](../CDR.md) — Field 3 (γ close-out artifact; verdict vocabulary); Field 4 (δ cadence; wave-transitions); Field 5 (ε trigger classes — γ's triage routes findings here); Field 6 (actor-collapse rule).
2. Load this file as the research-γ role surface.
3. Load the generic [`cnos.cdd/skills/cdd/gamma/SKILL.md`](../../../../cnos.cdd/skills/cdd/gamma/SKILL.md) for the kernel-grammar reference — observation-then-selection ordering, the issue-quality gate, the pre-dispatch scaffold check, the "transfer artifact facts not hidden role state" unblocking rule, the closure-gate checklist. Engineering-specific surfaces (CHANGELOG TSC table; encoding lag table; `RELEASE.md` authoring; `scripts/release.sh`) do not transfer; research-specific analogues are named below.
4. Load the typed receipt schema [`schemas/cdr/receipt.cue`](../../../../../../schemas/cdr/receipt.cue) — `#CDRReceipt` is γ's close-out artifact.
5. Load Tier 2 and Tier 3 skills as the wave's research domain requires (statistical-method skills, dataset-stewardship skills, citation-management skills as the project binding declares them).

## Selection inputs (research-specific)

Engineering γ selects from CHANGELOG TSC tables, encoding-lag tables, doctor/status surfaces. Research γ selects from research-specific observation surfaces:

1. **Open-claim ledger** — unverified claims awaiting evidence; claims with INDETERMINATE verdicts from prior waves; claims with BOUNDED-GO verdicts whose bounds could be closed.
2. **Construct-stability surface** — terms or measures whose definitions have drifted across receipts (`CDR.md §"Field 5"` construct-drift trigger).
3. **Measurement-disagreement surface** — recurring discrepancies between observations from different methods or different datasets that have not been adjudicated.
4. **Citation-debt surface** — claims previously transmitted that cited weak or now-stale references; claims that should have cited external work but did not.
5. **Project-binding waves table** — the project's `.cdr/waves/` directory or equivalent (per `CDR.md §"Empirical anchor"`); per-project wave-priority ordering, dependency chains, and binding rules.

The selection rule order is set by the project binding (a research project's `<project>/.cdr/POLICY.md` or equivalent). γ does not impose a universal research-wave selection ordering at the protocol layer; the protocol layer requires only that selection name (a) the gap being closed, (b) the project-binding selection rule that nominated it, (c) the dependency chain (if any).

## Wave-issue quality gate (research-specific)

Per the engineering γ's "issue-quality gate" + the research-discipline overlay, a wave is dispatchable when:

- The research gap is **falsifiable** at the question level (some observation would settle it).
- The candidate `claim_status` for each prospective claim is named (so α and β share calibration expectations).
- The candidate `data_refs` are identified (which mount + manifest + checksum α will use).
- The candidate `method_refs` are identified (which script + commit-SHA pinning α will use).
- The expected reproduction record is named (β's re-run command, if applicable).
- The verdict-vocabulary expectations are stated (which of GO / REVISE / NO-GO / INDETERMINATE / BOUNDED-GO is likely; under what conditions each fires).
- Project-binding policies that apply to the wave's data are named (data-use policy, citation policy, construct-stability constraints).
- Non-goals exist (research-specific: which claims this wave does *not* address; which adjacent constructs are out of scope).

If the wave-issue cannot be written to that level, the work is not ready for α dispatch. γ does not compensate for a weak issue by making the dispatch prompt longer.

## Dispatch (research-specific)

γ produces α + β prompts. δ (the operator) routes them — γ does not execute dispatch directly.

The research dispatch prompts include the same structural fields as engineering dispatch:
- the wave/issue link
- the load-order pointer (`cnos.cdr/skills/cdr/<role>/SKILL.md` + `CDR.md` + role-specific Tier 3 skills)
- the project binding's policy file
- the branch / artifact-channel coordinate (the project's wave-tracking surface)

γ does **not** include in the prompt:
- α's hidden rationale, when prompting β
- β's hidden rationale, when prompting α
- a paraphrase of the wave-issue (the prompt names the issue by reference; α and β read the issue directly)
- runtime-mechanic instructions (no `claude -p` invocation prose, no polling loops — those belong to δ and the harness, not to γ's prompt content)

## Wave-transitions (per `CDR.md §"Field 4"`)

γ records the wave-transition that opens each wave and the gate verdict that closes it. The cadence is **gate-transition-shaped, not release-shaped**:

- **Wave open** — γ records the trigger that opened the wave (a NO-GO verdict from a prior wave; a GO verdict triggering downstream synthesis; a BOUNDED-GO triggering scope expansion; an INDETERMINATE triggering measurement design; or an initial selection from the open-claim ledger).
- **Wave in progress** — γ coordinates α and β through the wave's claims-evidence loop; γ unblocks (artifact-fact transfer only); γ does not author α's matter or β's verdict.
- **Wave close** — γ emits the typed `#CDRReceipt` carrying β's verdict (GO / REVISE / NO-GO / INDETERMINATE / BOUNDED-GO per `boundary_decision.action` + `transmissibility`). γ records the boundary decision; the wave is closed.

There is no "tag" or "release" or "deploy" step. The receipt is the closure artifact. Untyped wave summaries (narrative reports) cite the receipt; they do not replace it.

## Close-out artifact: the typed `#CDRReceipt`

Per `CDR.md §"Field 3"`, γ's close-out artifact is the typed receipt at the surface declared by the project binding (typically `<project>/.cdr/waves/{wave-id}/receipt.md` or `.cue`). The receipt carries:

- `claim_refs` — which claims this receipt asserts (≥1).
- `data_refs` — dataset / mount / manifest / checksum (≥1).
- `method_refs` — script paths + commit SHA (≥1).
- `result_refs` — output file paths (≥1).
- `claim_status` per claim — `observed | computed | inferred | hypothesized | indeterminate`.
- `limitations` (optional) — explicit caveats.
- `reproduction` (required when `claim_status ∈ {observed, computed}`) — β re-ran; output matched.
- generic-kernel fields inherited from `schemas/cdd/#Receipt` — `protocol_id` (pinned to `cnos.cdd.cdr.receipt.v1`), `boundary_decision`, `verdict`, `protocol_gap_count`, `protocol_gap_refs`.

The receipt is what scope-n+1 reads when the wave crosses the trust boundary. A wave that has not produced a typed `#CDRReceipt` has not closed; it has only stopped.

## Triage of close-out findings → ε

γ triages each cycle-iteration finding into one of four dispositions (per the engineering γ's CAP discipline):

1. **Immediate MCA available** — ship the protocol patch now (γ lands it or names the delegate).
2. **Protocol MCI** — file `.cdd/iterations/` (or the research-protocol-iteration analogue per project binding) + issue.
3. **Agent MCI** — update the persona hub.
4. **One-off** — drop explicitly.

The trigger classes γ routes findings into ε's surface for are CDR.md §"Field 5":
- missing data gates (data-policy oracle drift)
- overclaiming (claim/evidence-alignment oracle drift)
- unreproducible numbers (reproduction-from-clean drift)
- weak citation discipline
- recurring oracle ambiguity
- construct drift

Silence is not triage. Every finding gets a disposition.

## Independence rules (research-specific)

### 1. Preserve α/β epistemic separation

γ sees both sides because coordination requires it. γ transfers **artifact facts** (the typed receipt, the data-policy file, the citation list) — not hidden role state (α's rationale, β's reasoning prose).

### 2. Do not author the matter or the verdict

If α is stuck on calibration, γ may clarify the wave-issue or point α to the relevant policy file. γ does not draft α's claim. If β disagrees with α's matter, γ records the disagreement and routes it for resolution; γ does not author β's audit.

### 3. Close the wave on the receipt, not on a tag

Engineering's "the tag is the disconnection point" (per `cnos.cdd/skills/cdd/operator/SKILL.md` §3.4) does not transfer. Research waves close on the typed receipt; the verdict is what makes the claim transmissible. There is no analogue of a release bundle.

### 4. Triage protocol-gap findings to ε

A repeated finding across waves (citation drift, construct drift, recurring oracle ambiguity) is ε's input, not γ's to silently absorb. γ routes such findings into ε's surface (the project-binding analogue of `cdd-iteration.md`).

## Persona / protocol / project boundary

This overlay declares **what research γ does at the protocol layer**. It does not declare:

- **Who is doing the coordination** (persona — layer 1). The persona hub's `## Discipline` section names the persona's loss function; this overlay names the role function.
- **What concrete project's wave-tracking surface is used** (project — layer 4). The wave-manifest path, the receipt file format (Markdown vs CUE vs JSON), the project's policy file — all project-binding. This overlay names the typed `#CDRReceipt` shape per the schema; the project binds where it lives and how it is rendered.

## Resumption

The generic γ doctrine's resumption protocol applies. Research γ's resumption cases:

- **Wave γ-closeout / γ-scaffold analogues** — if γ resumes mid-coordination, read the wave's existing artifacts and continue from the next uncompleted section.
- **Multi-wave dispatch coordination** — if γ resumes a wave-set with partial state, read the wave manifest's current row state and continue from the next pending wave.

**Never restart completed sections.** Committed sections represent settled γ coordination; resumption preserves the decisions and continues forward.
