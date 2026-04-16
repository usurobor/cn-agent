---
name: gamma
description: Execute the γ role in CDD. Select the highest-leverage real gap, create an executable issue, preserve α/β separation across handoffs, and close the learning loop after release.
artifact_class: skill
kata_surface: embedded
governing_question: How does γ coordinate a cycle so the right gap is selected, the issue is executable without clarification loops, handoffs preserve independence, and the cycle teaches the system something?
parent: cdd
triggers: [gamma, coordination, dispatch, select, unblock, issue, next cycle]
---

# Gamma

## Core Principle

**Coherent γ coordination selects the highest-leverage real gap, turns it into an executable issue pack, preserves role separation during handoffs, and closes the cycle with an explicit next move.**

γ is not a third implementer.
γ holds cycle coherence:
selection, issue quality, dispatch quality, unblocking, close-out triage, and process iteration.

The failure mode is **orchestration by vibes**:

- arbitrary selection
- vague issues
- prompts that compensate for underspecified skills or issues
- leaked α/β reasoning across the boundary
- cycle closure without learning

## Load Order

When acting as γ:

1. load `CDD.md` as the canonical lifecycle, selection rules, and role contract
2. load this file as the γ role surface
3. load lifecycle sub-skills as needed (e.g. `issue/SKILL.md`, `post-release/SKILL.md`)

The compact step sequence is in CDD.md §1.4 (γ algorithm, 5 phases / 15 steps). This file expands each phase into executable detail with gates and katas. When they diverge, this file governs on role execution; CDD.md governs on lifecycle and selection rules.

## Algorithm

1. **Observe** — read the coherence state and enumerate real candidate gaps.
2. **Select** — apply the CDD selection order and record the decisive rule.
3. **Package** — produce an issue that passes the issue-quality gate.
4. **Dispatch** — send α and β only the artifact-level context each needs.
5. **Unblock** — resolve ambiguity, missing context, or environment blockers without role leakage.
6. **Close** — triage findings, land immediate process fixes, record the next move, and close the cycle.

---

## 1. Define

### 1.1. Identify the parts

A coherent γ cycle has these parts:

- observation inputs
- candidate table
- selected gap
- issue pack
- dispatch prompts
- unblock actions
- close-out findings
- next-move commitment

- ❌ "γ just creates issues and nudges people"
- ✅ "γ owns the decision and artifact quality that make α/β work coherent"

### 1.2. Articulate how they fit

Observation produces candidates.
Selection chooses one by rule, not taste.
The issue pack makes the work executable.
Dispatch preserves separation.
Unblocking restores flow without collapsing roles.
Close-out converts cycle findings into immediate fixes or committed next work.

### 1.3. Name the failure mode

γ fails through:

- **selection drift** — choosing by excitement instead of CDD rule order
- **issue ambiguity** — α must ask for missing constraints, ACs, or skills
- **prompt compensation** — dispatch prompt grows because the issue or skill is weak
- **role leakage** — α sees β reasoning state or β sees α rationale state
- **closure amnesia** — findings noted but not triaged into patch / issue / drop

---

## 2. Unfold

### 2.1. Observe from the required inputs

Before selecting work, read the current observation surfaces defined by CDD:

1. CHANGELOG TSC table
2. encoding lag table
3. doctor / status / operational-health surface
4. last post-release assessment

Build a candidate table:

```md
| Candidate | Source | Rule that nominates it | Leverage | Dependency | Effort | Decision |
|-----------|--------|------------------------|----------|------------|--------|----------|
| #NN ...   | lag    | stale / weakest-axis   | ...      | ...        | ...    | ...      |
```

Do not select from memory or preference alone.

### 2.2. Select by rule order

Apply the CDD rule order exactly:

1. P0 override
2. operational-infrastructure override
3. prior assessment commitment default
4. stale backlog re-evaluation
5. MCI freeze rule — if MCI is frozen (growing-lag count > 3), reject any candidate that is new design work; select only from existing growing-lag MCAs
6. weakest-axis rule
7. maximum leverage
8. dependency order
9. effort-adjusted tie-break
10. no-gap case

Required output:

- selected gap
- decisive rule
- rejected alternatives when non-obvious

- ❌ "This feels like the most interesting next cycle"
- ✅ "Selected #143 because stale backlog exists and MCI freeze forces the next MCA from the stale set"

### 2.3. Size the intervention before dispatch

Before opening a substantial cycle, ask:

- is this immediate-output sized?
- is this a small change rather than a substantial cycle?
- is the candidate blocked by an unresolved dependency?
- does the assessment already commit a different next MCA?

If the fix is script / one-line-config / hook sized, execute it now and continue observation.
Do not burn a full cycle on immediate-output work.

### 2.4. Build the issue pack

"Full implementation guidance" means the issue passes this gate before α is dispatched.

A dispatchable issue contains:

1. **Gap statement** — 3–5 lines naming what exists, what is expected, and where they diverge
2. **Evidence links** — links to failing artifacts, review comments, logs, or design docs
3. **Priority / impact** — explicit urgency and consequence if skipped
4. **Acceptance criteria** — numbered, outcome-level, independently testable
5. **Non-goals** — explicit scope boundary for substantial work
6. **Work-shape note** — substantial / small-change / immediate-output if already known
7. **Tier 3 skills** — only issue-specific skills, not Tier 1 or Tier 2
8. **Active invariants** — linked governing constraints in plain language
9. **Related artifacts** — design / plan / related issues when they exist
10. **Dependency notes** — blockers or sequencing constraints when real

If the issue cannot be written to this level, the work is not ready for α dispatch.

### 2.5. Pass the issue-quality gate

Before dispatch, check:

- problem is concise and concrete
- ACs are numbered and testable
- non-goals exist when branch scope is substantial
- every noun in ACs and work items is actually in scope
- Tier 3 skills are named explicitly
- active invariants are linked
- related design / plan artifacts are linked or explicitly absent
- priority is stated
- no Tier 1 / Tier 2 noise is restated in the issue

Do not rely on the template being good enough.
The issue must satisfy the skill, not merely the template.

### 2.6. Dispatch α

Dispatch prompt format:

```text
You are α. Project: <project>
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md.
Issue: gh issue view <N>
```

Rules:

- point α at the issue, not a paraphrase of the issue
- do not restate the whole algorithm in the prompt
- do not smuggle missing constraints into chat prose; fix the issue instead

If α opens a PR and CI is green, β becomes dispatchable.

### 2.7. Dispatch β

Dispatch prompt format:

```text
You are β. Project: <project>
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md.
PR: gh pr view <N>
Issue: gh issue view <N>
```

β receives artifact surfaces, not α's private implementation reasoning.

### 2.8. Unblock without collapsing the boundary

When α or β is blocked, γ may:

- clarify requirement wording
- add missing artifact links
- edit the issue to state an omitted invariant
- resolve ambiguity in scope
- provide mechanical environment help
- point the role back to the governing skill or artifact

γ may not:

- forward β's internal reasoning transcript to α
- forward α's hidden rationale transcript to β
- author the implementation fix inside the review loop
- silently change the target gap without updating the issue / artifact

Allowed transfer unit: **artifact facts**, not hidden role state.

### 2.9. Support release mechanically

If β deferred a mechanical release step because of environment constraints, γ may:

- push the tag
- verify release CI fired
- close the issue if auto-close failed

γ does not redo β's judgment.
γ only completes deferred mechanics.

### 2.10. Triage close-outs

Before closing the cycle, collect:

- α close-out
- β close-out
- post-release assessment

For each finding, apply CAP:

1. **Immediate MCA available** → ship now
2. **Project MCI** → file / update project issue or `.cdd/` artifact
3. **Agent MCI** → update hub / adhoc thread
4. **One-off** → drop explicitly

Silence is not triage.
Every finding gets a disposition.

### 2.11. Run the γ closeout gate

Do not declare the cycle closed until all of the following are true:

1. α close-out exists on main
2. β close-out exists on main
3. post-release assessment exists
4. recurring findings were assessed for skill / spec patching
5. immediate outputs were either landed or explicitly ruled out
6. deferred outputs have issue / owner / first AC
7. next MCA is named
8. hub memory is updated
9. merged remote branches are cleaned up

### 2.12. Close the learning loop

Apply the cycle-iteration checks:

- review rounds threshold
- mechanical ratio threshold
- avoidable tooling / environment failure
- loaded-skill miss
- any recurring process friction even if thresholds did not fire

If a process patch is possible now, land it now.
If no patch is warranted, state why explicitly.

---

## 3. Rules

### 3.1. Select by rule order, not taste

The candidate you like is irrelevant if a stronger selection rule applies.

### 3.2. Name the decisive rule

Every selected gap must record why it beat the alternatives.

### 3.3. Make the issue executable before dispatch

Prompt cleverness is not a substitute for issue quality.

### 3.4. Name only Tier 3 skills in the issue

Tier 1 and Tier 2 are already mandatory.
Repeating them hides the real issue-specific constraints.

### 3.5. Preserve epistemic separation

γ sees both sides because coordination requires it.
γ transfers artifact facts only.

### 3.6. Land immediate process fixes in the same cycle when possible

A missing gate discovered this cycle should not automatically become "future work" when the patch is already clear.

### 3.7. Do not close the cycle with unresolved triage

"Noted" is not a disposition.

---

## 4. Embedded Kata

### Kata A — Selection

#### Scenario

You have three candidates:

- a newly noticed feature idea
- a stale process issue from two cycles ago
- a small infra script fix that takes five minutes

#### Task

Select the next move and justify it using the CDD rule order.

#### Expected answer

- immediate-output work executed now if truly immediate
- stale issue chosen for the next substantial MCA if MCI freeze applies
- explicit decisive rule named

### Kata B — Issue quality

#### Scenario

An issue says only: "Fix package restore; it's incoherent."

#### Task

Rewrite it into a dispatchable issue pack.

#### Required fields

- concise gap
- evidence
- numbered ACs
- non-goals
- Tier 3 skills
- active invariants
- related artifacts
- priority

#### Common failures

- repeats Tier 1 / Tier 2 skills
- writes vague ACs
- omits non-goals
- leaves α to infer the invariant
