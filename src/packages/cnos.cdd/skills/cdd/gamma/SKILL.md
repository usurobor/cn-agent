---
name: gamma
description: γ role in CDD. Coordination + closure + triage. γ produces dispatch prompts, unblocks via artifact facts, runs post-release assessment, triages findings, and declares cycle closure — binding the cell's ε-observations/learning section into the receipt at closeout. δ executes dispatch; harness owns the mechanics.
artifact_class: skill
kata_surface: embedded
governing_question: How does γ keep the full cycle coherent across issue creation, dispatch coordination, unblocking, assessment, and close-out — without leaking into runtime supervision mechanics?
visibility: internal
parent: cdd
triggers:
  - gamma
scope: role-local
inputs:
  - repo state
  - lag and signals
  - issue state
  - branch state
  - .cdd/unreleased/{N}/self-coherence.md (α's gap, mode, ACs, CDD Trace, review-readiness, fix-rounds)
  - .cdd/unreleased/{N}/beta-review.md (β's round-by-round verdicts + findings)
  - .cdd/unreleased/{N}/alpha-closeout.md (α close-out, post-merge)
  - .cdd/unreleased/{N}/beta-closeout.md (β close-out + release evidence, post-merge)
  - release state
  - "delta gate results (observable via git: tags, branch state)"
outputs:
  - issue pack
  - dispatch prompts
  - unblock decisions (clarification artifacts)
  - post-release assessment
  - close-out triage table
  - cycle closure declaration (gamma-closeout.md)
requires:
  - active role is γ
  - canonical CDD.md loaded
calls:
  - issue/SKILL.md
  - post-release/SKILL.md
  - operator/SKILL.md
  - harness/SKILL.md
  - release-effector/SKILL.md
---

# Gamma

## Core Principle

**Coherent γ coordination selects the highest-leverage real gap, turns it into an executable issue pack, preserves role separation during handoffs, and closes the cycle with an explicit next move.**

γ is not a third implementer and not a runtime orchestrator. γ holds cycle coherence: selection, issue quality, dispatch-prompt quality, unblocking, close-out triage, and process iteration. γ produces α and β prompts; δ (operator) executes dispatch via the harness — one role at a time. This avoids nested subprocess chains (γ spawning α/β), keeps memory pressure low, and gives δ direct visibility into each session. γ's closeout responsibility explicitly includes binding the cell's learning/ε-observations into the receipt (`CELL-KINDS.md` §"Mandatory terminal learning section") — this is a clarification of the existing closeout obligation, not new scope.

The failure mode is **orchestration by vibes**: arbitrary selection, vague issues, prompts that compensate for underspecified skills, leaked α/β reasoning across the boundary, cycle closure without learning. The related failure is **managerial residue** — verbs like *monitor* / *track* / *supervise* that produce no artifact, receipt, or boundary decision. The sweep rule (`COHERENCE-CELL.md §γ and δ Managerial-Residue Sweep`) is binding: any γ verb that produces no concrete output is suspect; rename to the actual biological function or relocate to the substrate.

## Load Order

When acting as γ:

1. Load `CDD.md` — canonical lifecycle, selection rules, role contract.
2. Load this file — the γ role surface (coordination + closure + triage).
3. Load `issue/SKILL.md` — issue-pack contract.
4. Load `post-release/SKILL.md` — γ owns the PRA (cycle-level assessment of α, β, cycle economics) and step 13a skill/spec patches.
5. Load `operator/SKILL.md` — δ-role policy frame. `harness/SKILL.md` owns dispatch mechanics (invocation shell, observability, worktree, identity, polling, timeout recovery); `release-effector/SKILL.md` owns release mechanics (`scripts/release.sh`, post-push CI polling, CI-red recovery, branch deletes). γ references these when authoring prompts that name observable-output flags or when requesting gate actions.
6. Load other lifecycle sub-skills only when the selected gap requires them.

**Canonical-skill staleness check before each γ phase change.** Before transitioning from one CDD phase to another (intake → dispatch, dispatch → close-out triage, close-out triage → PRA, PRA → closure), γ runs `git fetch --verbose origin main && git rev-parse origin/main`. If `origin/main` HEAD has advanced beyond the SHA at which canonical CDD/role skills were loaded, **re-load** `CDD.md`, this file, and the lifecycle sub-skill governing the next phase, then re-evaluate the next phase's plan against the updated canonical surfaces. This is the parallel of `beta/SKILL.md`'s pre-merge gate row 2 on the γ-axis; it catches canonical-skill snapshot drift across release boundaries (cycle #301 §9.1 trigger 1: γ proposed an out-of-spec merge option because σ's `4a0f678` "merge is β authority" had landed mid-cycle and was not in γ's session-loaded `gamma/SKILL.md`).

Canonical artifact locations (PRA, close-out paths, snapshot dirs, tag policy) are defined in `cnos.cds/skills/cds/CDS.md` §"Artifact contract" → §"Location matrix".

`cnos.cds/skills/cds/CDS.md` is the canonical source for: the ordered γ lifecycle (§"Development lifecycle" → §"Step table" — γ owns Steps 0–3 and Steps 9–13; Step 8 is β's; Steps 9–10 are δ's release-gate authority; Step 10's tag mechanics are δ's via `release-effector/SKILL.md`) and selection rule order (§"Selection function"). This file does **not** redefine that algorithm — it expands those steps into executable checks, evidence, and gates.

## Step map

- Step 1 → git identity (`harness/SKILL.md` §3), Step 2a → observe + candidates (`§2.1`), Step 2b → select + size (`§2.2`)
- Step 3 → issue pack + quality gate (`§2.3–§2.4`)
- Step 3a → create cycle branch (`§2.5` — `cycle/{N}` from `origin/main`, γ-owned pre-flight per `cnos.cds/skills/cds/CDS.md` §"Development lifecycle" → §"Branch pre-flight")
- Steps 3b–6 → dispatch + unblocking (`§2.5`)
- Steps 7–8 → release-prep artifacts (`§2.6`)
- Steps 9–10 → close-out triage (`§2.7`)
- Steps 11–13 → cycle iteration + process patching (`§2.8–§2.9`)
- Steps 14–16 → hub memory, branch cleanup, closure declaration (`§2.10`)
- Step 17 → δ disconnect release (mechanics: `release-effector/SKILL.md`; doctrinal frame: `delta/SKILL.md` §1) — not γ's

---

## 1. Define

**Parts of a coherent γ cycle:** observation inputs / candidate table / selected gap / issue pack / dispatch prompts / unblock actions (clarification artifacts) / close-out findings / next-move commitment. Observation produces candidates; selection chooses one by rule, not taste; the issue pack makes the work executable; dispatch preserves separation; unblocking restores flow without collapsing roles; close-out converts findings into immediate fixes or committed next work.

**Failure modes.** γ fails through:
- **selection drift** — choosing by excitement instead of CDD rule order
- **issue ambiguity** — α must ask for missing constraints, ACs, or skills
- **prompt compensation** — dispatch prompt grows because the issue or skill is weak
- **role leakage** — α sees β reasoning state or β sees α rationale state
- **closure amnesia** — findings noted but not triaged into patch / issue / drop
- **managerial residue** — verbs without concrete output (track, monitor, supervise) leaking from substrate into γ doctrine

- ❌ "γ just creates issues and nudges people"
- ✅ "γ owns the decision and artifact quality that make α/β work coherent"

---

## 2. Unfold

### 2.1. Step 1a — Observe and build candidates

Read the observation surfaces required by `CDD.md`:

1. **Last post-release assessment** (read first) — prior cycle's next-MCA commitment, deferred outputs, cycle iteration findings, MCI freeze state. Binding, not optional.
2. CHANGELOG TSC table.
3. Encoding lag table.
4. Doctor / status / operational-health surface.
5. **Cross-repo proposal intake** — scan `.cdd/iterations/cross-repo/cnos/*/STATUS` for the last non-comment event = `submitted` (or `drafted` with source-delegated filing-authority per [`cnos.handoff/skills/handoff/cross-repo/SKILL.md §2.3.3`](../../../../cnos.handoff/skills/handoff/cross-repo/SKILL.md)). For each candidate, read the adjacent `ISSUE.md` + optional `PATCH.diff`, check cnos target state for duplicate / already-landed work, decide `accepted` / `modified` / `rejected`, file a target issue with the `## Source Proposal` block from `cdd/issue/SKILL.md`, and emit the disposition event to source STATUS. Protocol details (8-event vocabulary, transition graph, bundle file set, LINEAGE schema, feedback-patch format) live in [`cnos.handoff/skills/handoff/cross-repo/SKILL.md`](../../../../cnos.handoff/skills/handoff/cross-repo/SKILL.md).

Build a candidate table:

```md
| Candidate | Source | Rule clause that nominates it             | Dependency | Size | Decision |
|-----------|--------|-------------------------------------------|------------|------|----------|
| #NN ...   | lag    | CDS §Selection function → §X              | ...        | ...  | ...      |
```

- ❌ "I remember issue #143 felt important"
- ✅ "Candidate table shows #143 is selected under `cnos.cds/skills/cds/CDS.md` §"Selection function" §X.x for a named reason"

### 2.2. Step 1b — Select by CDS §Selection function, then size the intervention

Apply `cnos.cds/skills/cds/CDS.md` §"Selection function" in order (P0 override → operational-infrastructure override → assessment-commitment default → stale-backlog re-evaluation → MCI freeze check → weakest-axis rule → maximum-leverage rule → dependency order → effort-adjusted tie-break → no-gap case). Required output: selected gap / decisive `cnos.cds/skills/cds/CDS.md` §"Selection function" rule-clause / rejected alternatives when non-obvious / intervention size (immediate output / small change / substantial cycle).

Sizing rule: immediate-output → execute now + continue observation; small-change → route to small-change path; blocked → name the blocking dependency before dispatch.

- ❌ "Selected because it felt highest leverage"
- ✅ "Selected #143 under `cnos.cds/skills/cds/CDS.md` §"Selection function" → §"Maximum-leverage rule"; #151 rejected because dependency unresolved"

### 2.2a. Peer enumeration at scaffold time

Before authoring §Gap in `self-coherence.md`, γ MUST: (1) list every file in directories named by the issue's impact graph (`ls -la` or `find` per directory); (2) grep for the term/symbol/surface the cycle proposes to add (`rg <term> <directories>`); (3) if any match is found, name it explicitly in §Gap — additive framing ("partially closed by X; this cycle completes it") or consolidation framing ("X overlaps the proposed surface; reconcile in scope").

A §Gap that asserts "X does not exist" without grep-evidence is a γ-side honest-claim violation analogous to α's rule 3.13(a).

**Empirical anchor — tsc cycle #36:** §Gap asserted "CI does not invoke `coh --kata` against shipped kata content"; `.github/workflows/ci.yml`'s `kata-check` job (added in 344-c) invoked `bash scripts/run-katas.sh` which auto-discovers every kata. The negation was empirically false; β R1 caught it as binding finding B-1 (RC). Cost: 1 wasted RC round.

- ❌ "X does not exist" without grep evidence
- ✅ "X does not exist — confirmed by `rg X directories/` returning no matches"

### 2.3. Step 2 — Build the issue pack

Use `issue/SKILL.md` as the base contract. A dispatchable γ issue is:
- a complete `issue/SKILL.md` issue
- **plus** a work-shape note (`substantial` / `small-change` / `immediate-output`)
- **plus** dependency notes when sequencing or blockers are real

γ does not restate Tier 1 or Tier 2 skills in the issue. γ does name Tier 3 skills, active design constraints, related artifacts, non-goals, and priority exactly as `issue/SKILL.md` requires.

If the issue cannot be written to that level, the work is not ready for α dispatch.

### 2.4. Step 2 — Pass the issue-quality gate

Before dispatch, verify:
- issue satisfies `issue/SKILL.md`
- ACs are numbered and independently testable
- every noun in ACs and work items is in scope
- non-goals exist when the issue is substantial
- Tier 3 skills are named explicitly
- active design constraints are linked and stated plainly
- related artifacts are linked or explicitly absent
- priority is stated
- work-shape is stated
- dependency notes exist when blockers or sequencing are real
- no prompt-only constraints are hiding outside the issue

Do not compensate for a weak issue by making the prompt longer. Fix the issue instead.

### 2.5. Steps 3a–5 — Create the cycle branch, dispatch, unblock without leakage

#### Step 3a — Create the cycle branch

After the issue passes the quality gate (§2.4) and **before** dispatching α/β, γ creates `cycle/{N}` from `origin/main`. The branch is the canonical coordination surface (`cnos.cds/skills/cds/CDS.md` §"Coordination surfaces" + §"Development lifecycle" → §"Branch rule"); a single named target replaces the pre-#287 model where α opened the branch.

**γ-owned branch pre-flight (`cnos.cds/skills/cds/CDS.md` §"Development lifecycle" → §"Branch pre-flight"):** `origin/cycle/{N}` does not yet exist (fail loud — one cycle = one branch); no stalled `.cdd/unreleased/{N}/` on `origin/main` (would indicate a previous cycle for `{N}` did not complete its release-time move per `release/SKILL.md` §2.5a); the issue's scope is declared in the body; base SHA known (`git rev-parse origin/main`); issue is open. Mechanics (the `git fetch` + `rev-parse --verify` + `ls-tree` + `switch -c` + `push -u` sequence) follow `cnos.cds/skills/cds/CDS.md` §"Development lifecycle" → §"Branch pre-flight" directly; γ does not restate them here.

The branch must exist on `origin` before dispatch. α and β never create branches; their prompts include a `Branch: cycle/<N>` line and they `git switch` to it.

#### Step 3b — Subscribe and dispatch α and β

##### Pre-dispatch γ scaffold check (binding gate)

**Binding rule.** γ MUST NOT proceed to α dispatch until `.cdd/unreleased/{N}/gamma-scaffold.md` exists on `origin/cycle/{N}`. If absent when γ is about to produce the α prompt, γ authors it first, commits and pushes it to `cycle/{N}`, then continues.

The scaffold names: issue / mode (substantial / small-change / immediate-output, and §5.2 wave-mode if applicable) / surfaces γ expects α to touch / AC oracle approach / empirical anchor when framed (§2.2a) / expected diff scope. It is the executable rendering of what γ already decided at selection time, not new analysis.

**Dual of `review/SKILL.md` rule 3.11b.** When `gamma-scaffold.md` is missing on the cycle branch at review time and no `## Protocol exemption` exists in the sub-issue body, β returns REQUEST CHANGES (D-severity, `protocol-compliance`). β-side enforcement guarantees the artifact exists by merge time but pays the round-trip *per RC round*; γ-side pre-dispatch enforcement pays once *per cycle* at scaffold time. The two gates are symmetric.

**Empirical anchor — cycle #369.** β R1 fired rule 3.11b as D1 against α's review-ready signal (missing scaffold); γ recovered via path (a) at `227d2373`; β R2 APPROVED at `4e179db6`; merged at `ff54f2a0`. All 10 ACs were met at R1 — the one "no" row was the missing scaffold. The γ-side gate makes rule 3.11b's RC the dual-redundant safety net rather than the routine first detector.

Mechanical check (run before producing the α prompt): `git ls-tree -r --name-only origin/cycle/<N> .cdd/unreleased/<N>/gamma-scaffold.md` must be non-empty; otherwise γ authors the scaffold, commits, pushes, and re-runs.

**Protocol exemption.** Per `review/SKILL.md` §3.11b "Exemption discoverability": γ may skip the scaffold *only if* the sub-issue body carries an explicit `## Protocol exemption` section naming the reason (emergency patch, infra-only change with operator override, wave-manifest-as-γ-artifact per `operator/SKILL.md` §10).

**Canonical section-header reminder (binding content of the scaffold).** `gamma-scaffold.md` is the first cycle-directory artifact α reads, and its own prose refers to α's forthcoming sections using this skill file's `§`-cross-reference shorthand (e.g. "§Gap"). That shorthand names a *section*, not markdown header text — the scaffold MUST NOT let that ambiguity carry into α's authoring. When a scaffold quotes or previews `self-coherence.md`'s section shape, γ writes the bare canonical form (`## Gap`, `## Skills`, `## ACs`, `## CDD Trace`, `## Self-check`, `## Debt`) verbatim, because `cn cdd verify`'s `sectionPresent()` (`src/packages/cnos.cdd/commands/cdd-verify/ledger.go`) does a literal/prefix match and a decorated `## §X` form hard-FAILs the non-lenient (pre-`beta-review.md` and release-time) validation path. *Derives from: cnos#608 → cnos#610 — #608's `self-coherence.md` had already discovered and fixed this exact header-form defect, but the lesson lived only in #608's own cycle-scoped artifact and was never promoted into a Tier 1/2 skill surface either γ or α load at dispatch time; #610's γ-authored scaffold (this same cycle family, three cycles later) did not carry the warning forward, and #610's α re-hit the identical `## §Gap` mistake, costing a full R0→R1 round. The durable fix lives at `alpha/SKILL.md` §2.5 (self-coherence authoring, the artifact actually gated); this entry is the γ-side reminder so the scaffold reinforces it rather than silently relying on α's own load order.*

##### Polling cross-reference

γ polls the issue and `origin/cycle/{N}` to react to α/β commits, β verdicts, and CI status. **Mechanics: `harness/SKILL.md` §5.4** (single-named-branch transition loop, `Monitor`-wrapped, reachability re-probe per `cnos.cds/skills/cds/CDS.md` §"Coordination surfaces" → §"Polling primitives"). γ's invariant: **polling requires a query, a wake-up mechanism, and a reachability probe.** If the environment provides neither a `Monitor`-equivalent nor a shell-wake harness, γ surfaces the gap to operator before dispatch. Cycle-dir artifacts live on `origin/cycle/{N}` (not on `main`) until release-time move; polling `main` for in-flight cycle dirs is silent.

##### Dispatch prompts + implementation contract → cnos.handoff

γ produces the dispatch prompts (γ / α / β) and δ routes them via `cn dispatch`. γ does not execute dispatch; γ produces prompts. The α prompt MUST carry the 7-axis `## Implementation contract` block γ injects per the wire-format template — γ MUST NOT dispatch with empty / "TBD" rows.

**Canonical wire format lives at [`cnos.handoff/skills/handoff/dispatch/SKILL.md`](../../../../cnos.handoff/skills/handoff/dispatch/SKILL.md)** (Sub 3 of [cnos#404](https://github.com/usurobor/cnos/issues/404), shipped under [cnos#417](https://github.com/usurobor/cnos/issues/417)). That skill owns the γ / α / β prompt templates, the 7-axis implementation-contract schema (Language; CLI integration target; Package scoping; Existing-binary disposition; Runtime dependencies; JSON/wire contract preservation; Backward-compat invariant), the prompt rules (issue reference shape, explicit `Branch:` line, `--json` convention, do-not-restate-the-algorithm, do-not-smuggle-constraints), the δ-as-inward-membrane enrichment doctrine, and the four-surface mesh declaration (γ template ↔ δ enrichment ↔ α constraint ↔ β verification).

γ's role-local obligation here: **author** the dispatch prompts at the wire-format shape; **inject** the `## Implementation contract` section into the α prompt with every row populated; **escalate** to δ when a row's value is unclear; **log** any mid-cycle re-pin in `.cdd/unreleased/{N}/gamma-clarification.md` before signaling the change. The doctrine lives at handoff; the role-local procedure (γ authors the prompts during scaffold; δ reviews before routing) lives here as cycle-lifecycle coordination.

Tooling and observability flags (`--allowedTools`, `--output-format stream-json --verbose`, `--permission-mode acceptEdits`) live in `harness/SKILL.md` §1–§2; γ does not restate them here.

#### Spec-staleness propagation

When γ (or δ) pushes spec changes to `origin/main` while α/β sessions are in-flight, loaded skills become stale. γ's responsibility:
- **Identity-rotation mode (`cn dispatch` / `claude -p`):** Not applicable — each role invocation loads skills fresh from the filesystem; the next dispatch picks them up.
- **Long-lived polling sessions (legacy):** γ writes a coordination note to `.cdd/unreleased/{N}/gamma-coordination.md` on `cycle/{N}` naming the spec change and affected skill path; or surfaces via the dispatch channel.
- **When to propagate:** Changes to `CDD.md`, role skill files (`alpha/`, `beta/`, `gamma/`, `operator/`, `delta/`, `harness/`, `release-effector/`), `release/SKILL.md`, or `review/SKILL.md` landing on `main` while in-flight.
- **When not to propagate:** Changes outside the CDD package, doc-only, or issues/threads.

**Canonical wire-format home at [`cnos.handoff/skills/handoff/mid-flight/SKILL.md §2.6`](../../../../cnos.handoff/skills/handoff/mid-flight/SKILL.md)** (Sub 4 of [cnos#404](https://github.com/usurobor/cnos/issues/404) / cnos#418). The wire-format invariant ("γ writes a coordination note when spec changes land on origin/main mid-flight; cycle branch is the substrate") lives at cnos.handoff; the consumer-specific file list above (which cdd skill files trigger propagation) is the CDD-side realization.

*Derives from #301 §9.1: γ proposed an out-of-spec merge because σ's `4a0f678` ("merge is β authority") landed mid-cycle outside γ's loaded snapshot. The reactive fix is the staleness check in Load Order; this is the proactive side.*

#### Step 5 — Unblock

When α or β is blocked, γ may: clarify requirement wording / add missing artifact links / edit the issue to state an omitted invariant / resolve scope ambiguity / provide mechanical environment help / point the role back to the governing skill or artifact.

γ may **not**: forward β's reasoning transcript to α / forward α's rationale transcript to β / author the implementation fix inside the review loop / silently change the target gap without updating the issue.

**Allowed transfer unit: artifact facts**, not hidden role state.

- ❌ "β said this design is shaky; just rewrite the parser like this"
- ✅ "The issue omitted invariant X; γ adds it to the issue and points α back to the updated artifact"

**Issue-edit cache-bust procedure.** When γ edits an issue body mid-cycle, γ commits a `gamma-clarification.md` entry on the cycle branch *before* signaling the edit. The cycle-branch transition is the signal; γ does not need to chat-relay. The clarification names: date, edit summary, and which ACs / non-goals / constraints / artifacts changed. **Canonical wire-format home at [`cnos.handoff/skills/handoff/mid-flight/SKILL.md`](../../../../cnos.handoff/skills/handoff/mid-flight/SKILL.md)** (Sub 4 of [cnos#404](https://github.com/usurobor/cnos/issues/404) / cnos#418); this section retains the role-local procedure citation.

### 2.6. Steps 6–7 — Prepare release artifacts before δ tags

In the sequential dispatch model, β exits after merge. δ runs `scripts/release.sh` per `release-effector/SKILL.md` (stamp + tag) but does not author artifacts. γ owns two release-ordering obligations before requesting the tag:

1. **Write `RELEASE.md`** — per `release/SKILL.md` §2.5. The GitHub release body, at repo root, committed to main. Without it, release CI auto-generates sparse notes.
2. **Keep cycle directories at the release-gate path** — per `release/SKILL.md` §2.5a. `.cdd/unreleased/{N}/` remains in place through δ's exact release validation and tagged disconnect. γ archives it only after δ reports release completion and green CI.

**Before any push that follows a rebase, run the eng/ship rebase-integrity gate** (see `eng/ship` § Rebase-Collision Integrity).

`RELEASE.md` must be committed before γ requests the disconnect release from δ (§2.10 → `release-effector/SKILL.md`). The directory move is deliberately later.

- ❌ Leave RELEASE.md for δ to write (δ does not author)
- ❌ Move unreleased directories before δ validates them
- ❌ Assume β handled release prep (β exits at merge in sequential model)
- ✅ γ writes RELEASE.md → δ validates/tags/verifies → γ archives cycle dirs

### 2.7. Steps 8–9 — Triage close-outs explicitly

**Collecting close-outs in sequential bounded dispatch (`cnos.cds/skills/cds/CDS.md` §"Field 6: Actor collapse rule").** β exits after writing `beta-closeout.md`; α already exited after signaling review-readiness. γ obtains both close-outs before triaging:

1. **β close-out**: verify `.cdd/unreleased/{N}/beta-closeout.md` on main; if missing, request δ to re-dispatch β.
2. **α close-out**: verify `.cdd/unreleased/{N}/alpha-closeout.md` on main; if missing, request δ to re-dispatch α via the close-out prompt (`cnos.cds/skills/cds/CDS.md` §"Coordination surfaces"; prompt format in `operator/SKILL.md` §5.2 v0.1 overlay). **γ must explicitly request this re-dispatch.**

**Cross-repo proposal close-out.** If this cycle accepted or modified a source proposal, γ emits the `landed` event to source STATUS at merge. Per [`cnos.handoff/skills/handoff/cross-repo/SKILL.md §"STATUS state machine"`](../../../../cnos.handoff/skills/handoff/cross-repo/SKILL.md) master/sub rule: per-sub `landed` per merge, terminal master-close when the master closes. A proposal touched by the cycle may not remain at `accepted` / `modified` after target work lands.

**Post-merge CI verification (mandatory).** Before authoring `gamma-closeout.md`, γ verifies CI ran green on the merge commit (`gh run list --branch main --json status,conclusion,head_sha` filtered to merge SHA). Pending → delay close-out. Red → log as §9.1 trigger (avoidable tooling failure); γ-axis grade reflects; consider rollback or follow-on fix-cycle. Green → record run URL in `gamma-closeout.md` §Post-merge verification.

`self-coherence.md` + `beta-review.md` carry the in-cycle record; the two `*-closeout.md` files are γ's primary triage inputs. The cycle directory remains at `.cdd/unreleased/{N}/` through closeout and release validation; γ archives it only after the tagged disconnect (`cnos.cds/skills/cds/CDS.md` §"Artifact contract" → §"Location matrix").

Then write the post-release assessment per `post-release/SKILL.md` at `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` (CDD package: `docs/gamma/cdd/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md`). The PRA is γ's — it measures α's implementation, β's review quality, and cycle economics. β assessing its own review is a self-grading problem.

γ triages every finding from both close-outs + the PRA using CAP:
1. **Immediate MCA available** → ship now (γ lands the skill/spec patch per `cnos.cds/skills/cds/CDS.md` §"Closure" → §"Immediate outputs"; if delegated, name delegate + deadline)
2. **Project MCI** → file / update project issue or `.cdd/` artifact
3. **Agent MCI** → update hub / adhoc thread
4. **One-off** → drop explicitly

Minimum triage record:

```md
| Finding | Source | Type | Disposition | Artifact / commit |
|---------|--------|------|-------------|-------------------|
| ...     | α/β/assessment | process/skill/... | immediate MCA / project MCI / agent MCI / drop | ... |
```

Silence is not triage. Every finding gets a disposition. Step 13a skill/spec patches are γ's to land or explicitly delegate.

### 2.8. Steps 10–11 — Enforce cycle-iteration outputs when triggers fire

Apply the cycle-iteration checks named in `CDD.md` step 10. For each fired trigger, γ must do something explicit:

| Trigger | Fire condition | Required γ action | Closure rule |
|---|---|---|---|
| Review churn | review rounds > 2 | Verify the assessment contains a `Cycle Iteration` entry naming the root cause and the chosen disposition (patch landed now / next MCA / no patch with reason). | If missing, cycle cannot close. |
| Mechanical overload | mechanical ratio > 20% **and** total findings ≥ 10 | Verify the assessment names the recurring mechanical class and whether the mechanization patch was landed now or filed as a concrete MCA. | If missing, cycle cannot close. |
| Avoidable tooling / environment failure | environment or tooling blocked the cycle in a way a guardrail could likely prevent | Verify the assessment names the failure, workaround, and disposition (guard landed now / issue filed / no spec-level fix with reason). | If missing, cycle cannot close. |
| Loaded-skill miss | a loaded skill should have prevented a finding but did not | Patch the skill now when the correction is clear; otherwise verify the assessment names the exact skill gap and the concrete next MCA. | If neither patch nor committed MCA exists, cycle cannot close. |

Each fired trigger must end in one of three states:
- patch landed now
- concrete next MCA committed
- explicit no-patch decision with reason

### 2.9. Step 13 — Run the independent γ process-gap check

Even if no `cnos.cds/skills/cds/CDS.md` §"Assessment" → §"Cycle iteration triggers" trigger fired, γ must still ask:
- Did this cycle reveal a recurring friction?
- Was any gate too weak or too vague?
- Did a role skill fail to prevent a predictable error?
- Did coordination burden show a better mechanical path?

If yes:
- patch the skill / spec now when clear, **or**
- commit the concrete next MCA (issue / owner / first AC)

If no:
- state why not in one sentence inside the assessment or γ close-out

- ❌ "No trigger fired, so nothing to do"
- ✅ "No formal trigger fired, but dispatch kept compensating for issue ambiguity; γ patches issue-quality gate now"

### 2.10. Steps 13–15 — Declare post-merge closeout, then disconnect and archive

Do not declare post-merge closeout complete until all of the following are true:

1. `.cdd/unreleased/{N}/alpha-closeout.md` exists on main
2. `.cdd/unreleased/{N}/beta-closeout.md` exists on main
3. γ has written the post-release assessment per `post-release/SKILL.md`
4. every fired cycle-iteration trigger has a `Cycle Iteration` entry with root cause and disposition
5. recurring findings were assessed for skill / spec patching
6. immediate outputs were either landed or explicitly ruled out
7. deferred outputs have issue / owner / first AC
8. next MCA is named
9. hub memory is updated
10. merged remote branch cleanup is assigned to δ after release
11. the release batch/version assignment is named; `RELEASE.md` may be authored after this post-merge closeout declaration (§2.6)
12. cycle directory remains at `.cdd/unreleased/{N}/` for exact-cycle release validation (§2.6)
13. no release-boundary action has run out of order; passing this closeout gate authorizes γ to request preflight
14. if the cycle's receipt has `protocol_gap_count > 0` (≥1 finding tagged `cdd-skill-gap` / `cdd-protocol-gap` / `cdd-tooling-gap` / `cdd-metric-gap`), `.cdd/unreleased/{N}/cdd-iteration.md` exists with each finding structured per `post-release/SKILL.md` Step 5.6b, **and** `.cdd/iterations/INDEX.md` has a row for cycle N. If any finding shipped to a different repo, `.cdd/iterations/cross-repo/{target}/{slug}/` exists with bundle + `LINEAGE.md`. If `protocol_gap_count == 0`, no iteration file is required (per [`ROLES.md §4b.4`](../../../../../../ROLES.md), [`epsilon/SKILL.md §1`](../epsilon/SKILL.md), and [`activation/SKILL.md §22`](../activation/SKILL.md)); the INDEX row is also not required for empty-findings cycles.
15. **γ MUST assert the parent issue's close state before declaring closure — this is a hard gate, not a conditional fallback.** Run `gh issue view {N} --json state --jq .state`. If the result is `CLOSED`, record the asserted state in `gamma-closeout.md` (field below) and proceed. If the result is anything other than `CLOSED` (e.g. the merge subject lacked a close-keyword per `beta/SKILL.md §"Pre-merge gate"` row 5), γ MUST run `gh issue close {N}` immediately and record the discrepancy — the pre-assertion state, the corrective action taken, and the post-correction state — in `gamma-closeout.md`. Closure MUST NOT be declared while this row is unresolved. *Derives from: cnos#368 — cycle #367 merged with a bare `(#367)` reference (no close-keyword); no γ-side assertion existed to catch the resulting OPEN state, and the issue stayed OPEN ~24h until a manual close. This row is the structural closer of that gap.*

The mechanical lifecycle oracle is phase-specific. Before merge, run
`scripts/validate-release-gate.sh --mode pre-merge --cycle N`; it must not ask
for future close-outs. After merge and after the role close-outs have landed,
γ appends the post-merge closeout declaration and the exact standalone line
`CDD-Post-Merge-Closeout: complete` to `gamma-closeout.md`, then runs
`scripts/validate-release-gate.sh --mode post-merge --cycle N`. A pre-operator
assurance receipt in that filename does not satisfy post-merge closeout without
the marker. This marker is explicitly nonterminal: `RELEASE.md`, δ's tagged
disconnect, release CI, and γ's directory archive remain pending.

Then:
- write `.cdd/unreleased/{N}/gamma-closeout.md`. Contains: cycle summary, close-out triage table, §"Cycle iteration triggers" assessment, cycle iteration, skill gap candidate dispositions, deferred outputs, hub memory evidence, next MCA, **and the asserted issue-close state from row 15 above** (the `gh issue view {N} --json state --jq .state` result at assertion time, plus, if a discrepancy was found and corrected, the discrepancy note: pre-assertion state, corrective action, post-correction state). `gamma-closeout.md` MUST also include the mandatory `learning`/`epsilon_observations` section per `CELL-KINDS.md` §"Mandatory terminal learning section" (`observations`, `process_deltas`, `reusable_patterns`, `followups`, `operator_burden`) — γ binds this section into the receipt; it is not optional narrative. **δ must not tag/release before this artifact carries `CDD-Post-Merge-Closeout: complete` on main, but that marker alone does not close the cycle.**
- update hub memory
- request δ's release-boundary preflight and tagged disconnect. After δ reports
  release CI green, archive `.cdd/unreleased/{N}/` under the release version and
  commit that move on main. Only the tag plus this archival move establish
  terminal cycle closure; the pre-release marker does not.

### 2.11. γ as autonomous coordinator

On each polling transition, γ matches the event against named decision-points and acts: act autonomously / pause / ignore. γ does not *track* the cycle as a managerial verb — γ polls (mechanics: `harness/SKILL.md` §5.4) and on each transition produces a decision (the autonomous-coordinator artifact).

**Decision-point matching:**

| Decision point | γ action |
|---|---|
| Selection commit / Scope expansion / P0 override / β-approved merge of high-blast-radius / Design-call γ cannot resolve / Conflict-of-interest escalation / Process-debt commitment | pause (structured report to operator) |
| Normal unblocking / artifact coordination | act autonomously |
| γ's own commits / no-op transitions | ignore |

**Operator-facing report shapes.** γ surfaces consolidated state, not per-transition noise:
- *TLDR* (3–5 bullets at consolidated state): issue status / cycle branch SHA + CI / last commits / open findings / blocked-on
- *Decision-request*: 1-sentence problem / options with trade-offs / γ recommendation / impact-if-delayed
- *Deferred-question batch*: collect at natural coordination pauses, surface as batch

**Kata — 3-round autonomous cycle.** γ drives a substantial cycle (≥5 ACs) with one mid-cycle clarification: (1) selects, creates issue, requests α/β dispatch; (2) α signals review-readiness, β returns RC; (3) γ writes clarification, requests α fix-round; (4) α fixes, β approves, merge completes; (5) γ writes PRA, triages, closes; (6) δ cuts release — **no operator interaction except at spawn time**. Success criterion: operator receives only the final TLDR.

---

## 3. Rules

### 3.1. Select by rule order, not taste

The candidate you like is irrelevant if a stronger `CDD.md` rule applies.

### 3.2. Name the decisive clause

Every selected gap must record the `CDD.md` clause that made it win.

### 3.3. Make the issue executable before dispatch

Prompt cleverness is not a substitute for issue quality.

### 3.4. Name only Tier 3 skills in the issue

Tier 1 and Tier 2 are already mandatory. Repeating them hides the real issue-specific constraints.

### 3.5. Preserve epistemic separation

γ sees both sides because coordination requires it. γ transfers artifact facts only.

### 3.6. Land immediate process fixes in the same cycle when possible

A missing gate discovered this cycle should not automatically become future work when the patch is already clear.

### 3.7. Do not close the cycle with unresolved triage

"Noted" is not a disposition.

### 3.8. Silence rule

γ does not surface every transition to the operator; γ surfaces only decision-points and consolidated state. A noisy γ is an unencapsulated γ.

- ❌ "α committed at SHA abc123; β is reviewing; CI is running"
- ❌ "Found 3 findings in β round 1; α is fixing finding #2"
- ❌ "Polling detected branch update; reading new self-coherence.md"
- ✅ "Selection decision needed: P0 vs committed next-MCA" (decision-point hit)
- ✅ "TLDR requested: Issue #286, cycle/286 at def456, β approved, ready for merge"
- ✅ "Cycle #286 closed. Next: #287." (consolidated state at closure)

### 3.9. Managerial-residue sweep (binding before any γ-skill patch)

Per `COHERENCE-CELL.md §γ and δ Managerial-Residue Sweep`: any γ verb that does not produce an artifact, receipt, or boundary decision is suspect. *Monitor*, *track*, *supervise*, *oversee*, *manage* are sweep-suspect; rename to the biological function (observe, discriminate, route, validate, close, transport, release, repair-dispatch) or relocate the mechanics to the substrate (`harness/SKILL.md`, `release-effector/SKILL.md`). Coordination is legitimate; runtime supervision mechanics belong below γ.

---

## 4. Embedded Kata

### Kata A — Selection

**Scenario.** Three candidates: a newly noticed feature idea / a stale process issue from two cycles ago / a small infra script fix taking five minutes.

**Task.** Select the next move; justify using `CDD.md` rule order.

**Expected.** Immediate-output executed now if truly immediate; stale issue chosen for the next substantial MCA if the governing clause forces it; explicit decisive clause named.

### Kata B — Issue quality

**Scenario.** An issue says only: "Fix package restore; it's incoherent."

**Task.** Rewrite into a dispatchable γ issue pack.

**Required fields.** Concise gap / evidence / numbered ACs / non-goals / Tier 3 skills / active design constraints / related artifacts / priority / work-shape note / dependency note when applicable.

**Common failures.** Repeats Tier 1/2 skills; vague ACs; omits non-goals; leaves α to infer the invariant.

### Kata C — Cycle iteration

**Scenario.** A cycle finished with 3 review rounds, 12 findings (5 mechanical), one loaded-skill miss, no process patch yet committed.

**Task.** State what γ must require before closure.

**Expected.** `Cycle Iteration` section in the assessment; root cause named for review churn / mechanical overload / skill miss; disposition for each fired trigger; closure blocked until landed patch, concrete next MCA, or explicit no-patch justification exists.

## Resumption

When dispatched to an artifact path that already contains a section manifest, γ follows `cnos.cds/skills/cds/CDS.md` §"Large-file authoring rule" → §"Resumption protocol": (1) read existing manifest, (2) verify completed sections are coherent and complete, (3) continue from the first uncompleted section, (4) update manifest on each completion. γ-specific resumption cases: `gamma-closeout.md` (sections `[Close-out Triage, Trigger Assessment, Cycle Iteration, Deferred Outputs, Next-MCA Commitment]`); PRA (large multi-section); dispatch coordination (mid-prompt-generation resume). **Never restart completed sections** — they encode settled γ coordination decisions.
