---
name: operator
description: δ role in CDR. Enforces the data-mounted gate, accepts or rejects research receipts, and records gate verdicts on the wave-transition cadence.
artifact_class: skill
governing_question: What does the research operator gate on — and what does it refuse to do — during an active CDR wave?
parent: cdr
triggers:
  - operator
  - gate
  - unblock
scope: role-local
---

# Operator (research δ)

> **This is a CDR-specific extension of the generic cnos.cdd δ doctrine.** The
> kernel grammar (role-cell shape, δ-as-external-gate-holder + δ-as-router
> position, override authority structure) is inherited by reference from
> [`cnos.cdd/skills/cdd/operator/SKILL.md`](../../../../cnos.cdd/skills/cdd/operator/SKILL.md).
> Only the **gate vocabulary** (data-mounted + reproduction + citation +
> data-policy gates), the **close-out cadence** (gate-transition, not the
> engineering disconnect cadence), and the **verdict record surface** (the
> typed `#CDRReceipt`'s `boundary_decision`) diverge for the research loss
> function per
> [`ROLES.md §4a.2`](../../../../../../ROLES.md#4a2-loss-function-distinction)
> and the binding doctrinal contract [`CDR.md`](../CDR.md) Fields 3 + 4.

## Core Principle

**Coherent research δ operation enforces the data-mounted gate, records the gate verdict on the typed receipt, and stays out of the α/β/γ reasoning.**

Research δ is not a fourth research role; δ is not scored on coherence axes. δ owns what α/β/γ cannot: cross-wave routing, gate authority, override declaration, and the data-mounted gate that prevents an `observed`-status claim from being certified without the data actually present and verifiable. The failure mode is **gate bypass** — δ accepts a research receipt whose evidence refs do not actually resolve in the operating environment, or whose data is no longer mounted at the cited mount point, or whose reproduction record was not re-verified.

## Load Order

When acting as research δ:

1. Load [`CDR.md`](../CDR.md) — Field 3 (verdict vocabulary δ records); Field 4 (wave-transition cadence δ runs); Field 6 (actor-collapse rule δ enforces, especially the α=β prohibition).
2. Load this file as the research-δ role surface.
3. Load the generic [`cnos.cdd/skills/cdd/operator/SKILL.md`](../../../../cnos.cdd/skills/cdd/operator/SKILL.md) for the kernel-grammar reference — δ's role-boundary discipline (do not implement / do not review / do not triage findings), routing pattern, override-declaration protocol, "What the operator does NOT do" (§6), git-identity-for-role-actors convention. The engineering-specific surfaces (`scripts/release.sh`; the disconnect-the-triad-via-tag at §3.4; CI-green-or-recovery-runbook for release CI; branch-cleanup after release; the wave-coordination subsection at §10's release framing) do **not** transfer; research analogues are named below.
4. Load the typed receipt schema [`schemas/cdr/receipt.cue`](../../../../../../schemas/cdr/receipt.cue) — `#CDRReceipt.boundary_decision` is the surface δ records gate verdicts on.
5. Load Tier 2 and Tier 3 skills as the wave's project binding requires (dataset-mount verification, data-policy enforcement, citation-resolution as the project's `.cdr/POLICY.md` declares).

## Gate vocabulary (research-specific)

Engineering δ holds external gates: push, tag, release CI, branch cleanup. Research δ holds **research-protocol gates**:

### 1. Data-mounted gate

Before δ certifies an `observed`-status receipt, δ verifies the cited data was actually mounted at the time of the claim:
- mount point reachable in the operating environment
- checksum on `data_refs` entries matches the recorded checksum
- manifest path resolves and matches the recorded manifest
- data-use compliance record present and current (consent windows not expired; retention policy honored)

A receipt whose `data_refs` mount point is unreachable, whose checksum mismatches, or whose manifest cannot be re-validated **does not pass the data-mounted gate** regardless of α's pre-emission self-check or β's audit. The gate closes the surface where data could disappear between α's authoring and γ's emission.

### 2. Reproduction gate

Before δ certifies a `computed`-status receipt, δ verifies the reproduction record (`reproduction.output_match: true`). If β ran reproduction-from-clean and recorded a match, δ accepts; if β skipped reproduction or recorded a mismatch, δ does not certify. The gate is structural — δ does not re-run reproduction itself (that is β's role); δ verifies the record is present and well-formed.

### 3. Citation gate

For every claim citing external work, the cited reference is resolvable (DOI / URL / archive ref) and the cited content supports the claim. δ verifies citation resolvability mechanically (the reference dereferences); β's prior audit covered citation-support-substance. The gate prevents transmission of claims whose external grounding has decayed (link-rot, archived-paper retraction).

### 4. Data-policy gate

For every `data_refs` entry, the project binding's data-use policy is satisfied. δ verifies compliance by reference to the project's `<project>/.cdr/POLICY.md` (or equivalent project-binding file). The gate's specifics are project-bound; the gate's *presence* is protocol-bound.

### Verdict-recording (per `CDR.md §"Field 3"`)

After gates pass (or fail), δ records the verdict on the typed receipt:

- **GO** → `boundary_decision.action: accept` + `transmissibility: accepted`. All gates pass; the claim is transmissible at full strength.
- **REVISE** → `boundary_decision.action: repair_dispatch`. A named gap (calibration, evidence, citation, policy) requires α to fix; δ routes the receipt back via γ.
- **NO-GO** → `boundary_decision.action: reject`. The claim is not transmissible; further work is required.
- **INDETERMINATE** → the receipt carries `claim_status: indeterminate` for the affected claim(s); δ records the determination shape (what would change it).
- **BOUNDED-GO** → `transmissibility: degraded` + `limitations` enumerating the bound. The claim is transmissible within the named caveat.

δ does not invent new verdicts. δ does not map a known-failing gate to GO. δ does not skip the data-mounted gate because "the data was here yesterday."

## Wave-transition cadence (per `CDR.md §"Field 4"`)

Research δ runs a **wave-transition cadence**, not a release cadence. The cadence trigger that opens the next wave is the prior wave's gate verdict:
- NO-GO → may trigger a follow-up wave designed to close the gap that produced the NO-GO.
- GO → may trigger a downstream synthesis wave (combining the newly-transmissible claim with prior claims).
- BOUNDED-GO → may trigger a scope-expansion wave aimed at closing the bound.
- INDETERMINATE → may trigger a measurement-design wave aimed at producing the missing evidence.

δ records the trigger that opened each wave on the wave's coordination artifact. An unrecorded wave-open trigger breaks the receipt-stream signal ε reads (per `CDR.md §"Field 5"`).

## What research δ does NOT do

The engineering δ's "What the operator does NOT do" list transposes directly:

- **Do not author research matter.** That is α. If δ wants to fix a claim, declare an override and take the α role explicitly — and the role-separation discipline (α≠β within a wave per `CDR.md §"Field 6"`) still applies.
- **Do not audit claims.** That is β. If δ disagrees with β's verdict, declare an override per `cnos.cdd/skills/cdd/operator/SKILL.md` §4.
- **Do not triage protocol-gap findings.** That is γ → ε. δ may surface a finding δ observed during gate-execution, but δ does not assign the disposition.
- **Do not rewrite γ's dispatch prompts.** γ owns dispatch prompt quality.
- **Do not certify a receipt without running the gates.** Even if α and β both signed off, the data-mounted gate may reveal that the cited data is no longer reachable; δ does not assume gate satisfaction.
- **Do not communicate directly with α or β during in-wave work.** γ is the bridge. δ-to-γ; γ-to-α/β. δ remains external to the triad.

## What research δ does NOT do (engineering-specific surfaces explicitly out of scope)

The following engineering δ surfaces do **not** apply to research δ:

- **No `git tag` / `scripts/release.sh` / VERSION-stamping.** Research waves close on the typed receipt; there is no analogue of the disconnect-via-tag at `cnos.cdd/skills/cdd/operator/SKILL.md` §3.4. The receipt is the disconnection point.
- **No release CI / release-workflow polling.** A research wave's "CI" analogue is β's reproduction-from-clean (which β runs); δ verifies the record, not the execution.
- **No "cut the release" disconnection.** Per `CDR.md §"Field 4"`: "the receipt is the artifact; there is no release-bundle artifact in the engineering sense." δ records the verdict; the wave is closed.
- **No branch-delete post-cycle.** Branch lifecycle for research waves is a project-binding concern (the project's `.cdr/waves/` may use branches, may use directories, may use external tracking); the protocol does not bind δ to a branch-cleanup step.

## Git identity (per the generic operator doctrine)

Research δ configures git identity in the form `{role}@cdr.{project}.cdd.cnos` (or its elision form when the project is cnos itself: `{role}@cdr.cnos`). The role local-part for research δ is `operator` (mirroring the engineering δ's convention). The DNS namespace adds `cdr` to disambiguate from the engineering `{role}@cdd.{project}.cdd.cnos` form; the canonical right-to-left reading is "operator at the cdr protocol overlay running inside the cdd protocol family at the cnos repo, for the project named {project}."

```text
# general form (non-cnos projects):
git config user.name "operator"
git config user.email "operator@cdr.{project}.cdd.cnos"

# cnos itself (elision form):
git config user.name "operator"
git config user.email "operator@cdr.cnos"
```

The identity-truth row applies (per `cnos.cdd/skills/cdd/beta/SKILL.md` pre-merge gate row 1, transposed): every research-δ-authored artifact has δ's canonical email on it; identity drift is a finding.

## Override protocol

Per the generic δ doctrine §4, research δ may override α/β/γ decisions when:

- The wave's direction conflicts with project priorities the triad cannot see.
- A role is stuck and γ's unblocking has not resolved it.
- External constraints force a scope change (a dataset is withdrawn; a citation is retracted; a data-use policy changes mid-wave).
- Safety, ethics, or regulatory constraints require immediate action.

The override declares: (i) what is being overridden, (ii) why (the information the triad did not have), (iii) what the new state is. The override is recorded on the wave artifact channel; silence is structural drift.

Per the generic doctrine, override is for **information asymmetry or hard constraints, not preference**. δ does not override β's verdict because δ disagrees with the calibration call — that is β's judgment.

## Independence rule

δ enforces the α=β prohibition per `CDR.md §"Field 6"`. A research-claim wave whose receipt cannot distinguish α and β authorship (single signature, single session, no review record) is structurally invalid; δ rejects such a receipt regardless of its content. Engineering-class collapse precedents (γ+α+β-collapsed-on-δ for mechanical refactor cycles) do not transfer; δ does not permit them for research-claim transmission.

## Persona / protocol / project boundary

This overlay declares **what research δ does at the protocol layer**. It does not declare:

- **Who is enacting δ** (persona — layer 1). Persona hubs declare the persona's discipline.
- **What concrete data mounts, policy files, or wave-tracking surfaces δ operates on** (project — layer 4). The data-mounted gate verifies *what the project binding says the data should be*; the project binding declares the mount points, the manifests, the policy file. This overlay declares the gate's existence and shape; projects declare its concrete bindings.
