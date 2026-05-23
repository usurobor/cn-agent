# Design Notes — Cycle #395 (Sub 3 of #376)

**Date:** 2026-05-21
**Issue:** [#395](https://github.com/usurobor/cnos/issues/395)
**Author:** γ-as-α (γ+α+β-collapsed-on-δ; design half)
**Status:** complete

This document records the design rationale for the five CDR role-overlay SKILL.md files authored by Sub 3, plus the loader stub. It is the design-half record for the design-and-build cycle.

---

## 1. Design problem

The role-overlay design problem reduces to one differential: given the engineering exemplar `cnos.cdd/skills/cdd/<role>/SKILL.md`, what changes when the loss function shifts from *artifact improvement under repairable feedback* to *truth-preserving claim transmission under uncertainty* (per `ROLES.md §4a.2`)?

The kernel grammar is identical across protocols (per `ROLES.md §1`): α produces; β reviews; γ coordinates; δ holds external gates and routes; ε observes whether the protocol itself is coherent. The role-cell shape is preserved. What changes is the matter, the oracle, the close-out artifact, the cadence trigger, and the actor-collapse constraint — i.e. the six fields CDR.md instantiates per `ROLES.md §3`.

For each of the five role overlays, the design step is: identify what each engineering-skill subsection *means* under the research loss function, replace the discipline-specific language, keep the structural shape, and cite the parent doctrine by reference.

## 2. Per-role research-function-vs-software-function differential

### 2.1 α — Research α (`alpha/SKILL.md`)

**Engineering α** (cnos.cdd) produces source code, tests, and docs. Its primary failure is **closure overclaim** (claiming a class of gap is closed without enumerating peers; updating one surface while leaving siblings stale; asking β to find missing authoring work).

**Research α** produces claims, hypotheses, methods, datasets, analyses, and reports per CDR.md Field 1. Its primary failure is **overclaim** in the research sense: stating a claim stronger than its evidence supports; not declaring a `claim_status` calibration; not providing `data_refs` for `observed` claims, not providing `method_refs` for `computed` claims, not separating `inferred` from `observed`. The research α's discipline is **evidence-bound authoring**: every claim is calibrated, every method is reproducible, every data ref carries a manifest.

The engineering α's "peer enumeration before closure claims" and "harness audit for schema-bearing changes" subsections do not transfer directly. The research α's analogues are:
- **Claim-status calibration**: every claim labelled `observed | computed | inferred | hypothesized | indeterminate` before review.
- **Evidence-ref completeness**: `data_refs`, `method_refs`, `result_refs` populated per `schemas/cdr/receipt.cue` `#CDRReceipt`; no `observed` claim without a `data_refs` entry; no `computed` claim without `method_refs` + `result_refs`.
- **Method/data provenance**: methods cite script paths + commit SHA; datasets cite mount points + manifests + checksums + data-use compliance records.
- **Limitations declaration**: explicit caveats in `limitations`. The research α writes limitations *before* β asks for them.

### 2.2 β — Research β (`beta/SKILL.md`)

**Engineering β** (cnos.cdd) reviews against acceptance criteria, runs CI-equivalent validators, holds the pre-merge gate, and merges. Its core principle is *independent judgment from first review through release and β close-out*.

**Research β** reviews against CDR.md Field 2's review oracles: falsifiability, diagnostic oracles, reproduction-from-clean, citation integrity, data-policy compliance, claim/evidence alignment. The research β's primary work is **falsifiability + reproduction enforcement**: rejecting claims too vague to falsify; re-running the producing command in a clean environment; verifying citations support the claims they back. The engineering β's "compiles + tests pass" oracle does not apply (CDR.md §Field 2 declares this explicitly). The research β's verdict is recorded in the typed `#CDRReceipt`, not via `git merge`.

Engineering's pre-merge gate (`beta/SKILL.md` "Pre-merge gate" table) has a structural analogue in research: a **pre-receipt-emission gate** verifying:
- claim_status calibration consistent with evidence type
- reproduction record present when `claim_status ∈ {observed, computed}`
- data-policy compliance attested
- citation integrity verified
- diagnostic-oracle result recorded

These are research-specific β gates. The "Identity truth" row (γ.cdd canonical email) is structurally the same — research-β's commit-time identity is `beta@cdr.cnos` (or its project-elision form per `operator/SKILL.md §"Git identity for role actors"`).

### 2.3 γ — Research γ (`gamma/SKILL.md`)

**Engineering γ** (cnos.cdd) selects gaps per `CDD.md §3` rule order, builds issue packs, dispatches α/β through δ, unblocks without leaking role state, triages close-outs, runs cycle iteration, and writes `gamma-closeout.md` per the canonical lifecycle.

**Research γ** does the same coordination work, but the cycle is a **research wave**, not an engineering cycle. The selection rules differ (research selects by *open claim* or *unverified construct* or *recurring measurement disagreement*, not by lag tables or encoding lag). The dispatch surfaces are the same (α + β prompts; δ routes). The close-out artifact is the typed **research receipt** `#CDRReceipt` (CDR.md Field 3), not a `gamma-closeout.md` text artifact at the same path. The close-out is mechanically validated by V dereferencing `protocol_id: cnos.cdd.cdr.receipt.v1`.

The engineering γ's "Pre-dispatch γ scaffold check (binding gate)" and "issue-quality gate" subsections transfer structurally: research γ also pre-flights a wave-scaffold artifact (declaring what α will produce, what β will audit, what the gate verdict criteria are) before dispatching. The research γ's gate-verdict vocabulary (GO / REVISE / NO-GO / INDETERMINATE / BOUNDED-GO per CDR.md Field 3) is what γ closes the wave with via the receipt's `boundary_decision.action` + `transmissibility` fields.

### 2.4 δ — Research δ (`operator/SKILL.md`)

**Engineering δ** (cnos.cdd) routes prompts, holds external gates (push, tag, release, branch cleanup, force-push, auth refresh), executes the disconnect release (§3.4), and may override with explicit declaration.

**Research δ** routes prompts and holds gates, but the gates are **research-protocol gates**, not release gates. The disconnect-release model (§3.4 — "the tag is the disconnection point") does not apply: research waves close on **gate-transition** (CDR.md Field 4 — "the receipt is the artifact; there is no release-bundle artifact in the engineering sense"). The research δ's gate vocabulary is:
- **Data-mounted gate**: before accepting an `observed`-status receipt, verify the data was actually mounted at the time of the claim (mount point reachable, checksum matches, manifest valid).
- **Reproduction gate**: before accepting a `computed`-status receipt, verify the reproduction record (`reproduction.output_match: true`).
- **Citation gate**: every claim citing external work has a resolvable reference that supports the claim.
- **Data-policy gate**: data-use policy compliance (consent, anonymisation, retention).

The research δ records the gate verdict (GO / REVISE / NO-GO / INDETERMINATE / BOUNDED-GO) in `#CDRReceipt.boundary_decision.action` + `transmissibility`. Wave-transitions (CDR.md Field 4) replace the engineering δ's tag-push / release-script flow. No `scripts/release.sh` analogue.

The "Git identity for role actors" subsection transfers structurally (research δ's email is `operator@cdr.cnos` or its project form), but the engineering δ's `scripts/release.sh` workflow does not — research has no analogue.

### 2.5 ε — Research ε (`epsilon/SKILL.md`)

**Engineering ε** (cnos.cdd) is short by design: it observes whether CDD is coherent, writes `cdd-iteration.md` on every cycle, applies MCA discipline, collapses onto δ in small-protocol regimes.

**Research ε** is symmetrically short. It writes `cdr-iteration.md` on every wave (analogous artifact), observing CDR.md Field 5's trigger classes: missing data gates, overclaiming, unreproducible numbers, weak citation discipline, recurring oracle ambiguity, construct drift. The MCA discipline is identical (immediate MCA / next-MCA / drop). The ε=δ collapse for small-protocol regimes carries over directly — most early CDR work will have ε collapsed onto δ.

The differential is small: ε's structural role is the same (observe whether the protocol is itself coherent); only the trigger classes change (research failure classes per CDR.md Field 5, not engineering failure classes).

## 3. Persona/protocol/project boundary discipline (all five files)

Per CDR.md §"Persona, Protocol, Project" and `ROLES.md §4a`, each role overlay sits at **layer 3 (protocol overlay)**. The boundaries it must respect:

- **No layer-1 (persona) content.** No "I am Rho" prose, no "my voice", no "my temperament", no discipline-profile assertions on behalf of a named persona. The role overlay declares the role function; whatever persona enacts the role brings its own discipline profile from its hub (`cn-rho/spec/PERSONA.md` or equivalent).
- **No layer-4 (project) content.** No project-specific dataset paths (no `/opt/gait-data/`, no cph-local mount points), no project-specific gate thresholds, no project-specific report templates. Projects bind concrete artifacts via `<project>/.cdr/`; the role overlay declares the role-level shape only.
- **No surface fusion with runtime mechanics.** No `polling`, `dispatch`, `claude -p`, `gh issue`, `cn dispatch` invocations. Cross-references to cnos.cdd (which legitimately discusses these for engineering work) are permitted only as path strings — the CDR role overlay never re-authors that machinery.
- **No release-driver effection.** No `git tag`, no `scripts/release.sh`, no merge instructions, no branch-cleanup procedures. Research waves close on receipts, not tags.

## 4. Frontmatter shape

Each role overlay uses the cnos.cdd frontmatter shape with two modifications:

- `parent: cdr` (not `cdd`).
- `description` names the research role function (claims/methods/data for α; oracles for β; etc.).
- `triggers` includes the role name only (e.g. `alpha`). The role overlay is loaded when the active role for a CDR cycle is α/β/γ/δ/ε. No CDR-specific phase triggers (issue / design / plan / review / release / post-release) are declared because lifecycle sub-skills are deferred to cds emergence per the issue's Non-goals.
- `scope: role-local`.
- `governing_question` is research-specific (e.g. "How does α produce research matter without overclaim?").

## 5. Loader stub design

`src/packages/cnos.cdr/skills/cdr/SKILL.md` is Sub 2's deliverable. Sub 2 has not shipped. This Sub authors a **minimal stub** that:

- Has standard frontmatter (name: cdr, description, artifact_class: skill, parent: <none/package-root>, governing_question, triggers, scope: global).
- Declares: "This is the loader entrypoint for the CDR protocol skill bundle."
- Names the five role-overlay sub-skills (AC5 oracle: `rg "alpha/SKILL.md|beta/SKILL.md|gamma/SKILL.md|operator/SKILL.md|epsilon/SKILL.md"` ≥5).
- Cross-references `CDR.md` as the binding doctrinal contract.
- Cross-references the engineering exemplar `cnos.cdd/skills/cdd/SKILL.md` as a structural model.
- Notes that Sub 2 will extend with package-level concerns (`cn.package.json` pointer, package `README.md` cross-reference).

If Sub 2 ships before this Sub merges, the merge integrates Sub 2's loader by appending the role-skill references to whatever Sub 2 authored. If Sub 3 ships first, Sub 2's later cycle extends this stub with package-level concerns. Either ordering is integration-safe because the role-skill references and the package metadata pointers occupy distinct sections of the loader file.

## 6. Cross-reference plan

Each role overlay carries a stable set of cross-references:

- `../CDR.md` — Sub 1 contract; binding. Cited from each overlay's intro + relevant body sections (matter type for α; review oracle for β; close-out for γ; gate vocabulary for δ; ε trigger classes for ε).
- `../../../../../ROLES.md §4a.2` (loss-function distinction) — cited in each overlay's preamble.
- `../../../cnos.cdd/skills/cdd/<role>/SKILL.md` — generic doctrine. Cited as "this overlay is a CDR-specific extension of the cnos.cdd <role> doctrine; the kernel grammar is inherited by reference."
- `../../../../../schemas/cdr/receipt.cue` — γ overlay only (and δ overlay where the gate-verdict typed surface is named).

## 7. Open questions (for ε iteration, not blocking)

1. **Project-specific stricter floors templates** — already raised by Sub 1's #390 as a known next-MCA. CDR Field 6 allows projects to impose stricter floors (e.g. "for externally-published claims, β must be a distinct human reviewer"). No template ships with the role overlays. Sub 4 may surface this via cph's `.cdr/POLICY.md`.

2. **Wave-coordination primitive** — engineering γ's `cnos.cdd/skills/cdd/gamma/SKILL.md §10` declares wave coordination for multi-cycle dispatch. Whether research γ needs an analogue surface (research-wave coordination distinct from a single research cycle) is deferred. cph's wave-shape practice may inform this in Sub 4.

3. **Construct-stability artifact** — CDR.md Field 5 names "construct drift" as an ε trigger class but does not name the project-binding template (a glossary, a deprecation registry) that would address it. Deferred to cph's practice or to a later cnos.cdr cycle.

These are recorded for ε's cdd-iteration.md and not blocking on this Sub.

## 8. Refusal-condition check

Per the dispatch contract, refusal conditions for Sub 3:

- **CDR.md missing or unreadable** — not triggered. CDR.md exists (31,739 bytes) and is fully readable.
- **A role overlay needs persona-identity content to be coherent** — not triggered. The five role overlays declare role-function only; persona is layer-1 and lives in `cn-rho`.
- **A role overlay needs project-specific content** — not triggered. Project-binding is layer-4 and lives in `<project>/.cdr/`.
- **cnos.cdd exemplar uses patterns incompatible with cnos.cdr's research focus** — not triggered. The engineering exemplars use kernel-grammar patterns (Core Principle / Load Order / Algorithm / Rules / Resumption / Embedded Kata) that transfer to research with discipline-profile replacement. Specific engineering-discipline patterns (pre-merge gate; release/SKILL.md; CI/test rows) do not transfer mechanically but are *named as not-applicable* in each overlay with the research-discipline analogue named instead.

No refusal triggered. Proceed to build phase.
