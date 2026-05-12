---
name: post-release
description: Assess the released cycle, decide iteration, and close with explicit evidence and triage.
artifact_class: skill
kata_surface: embedded
governing_question: How does γ assess a released version, decide follow-up, and close the cycle with explicit evidence?
visibility: internal
parent: cdd
triggers:
  - post-release
  - assess
  - close
scope: task-local
inputs:
  - released version state
  - alpha close-out
  - beta close-out
  - production verification evidence
outputs:
  - post-release assessment
  - cycle-iteration disposition
  - closure evidence
requires:
  - released version exists
  - γ role active
calls: []
---

# Post-Release Assessment

This module implements CDD steps 11–13: post-release observation, γ-owned assessment (PRA), and close-out / closure evidence. After every release, assess what shipped, what the system looks like now, and what to do next. This is CDD §9 (Assessment) and §10 (Closure) executed as a concrete procedure.

Canonical artifact locations (PRA path, close-out paths, snapshot dirs, tag policy) are defined in `CDD.md` §5.3a (Artifact Location Matrix). Tags are bare `X.Y.Z` everywhere; `v`-prefixed tags are legacy and warn-only.

## Who

**γ owns the post-release assessment.** The PRA is a cycle-level observation artifact — it measures α's implementation, β's review quality, and the cycle's economics. β assessing its own review quality is a self-grading problem that weakens the independence CDD exists to provide. γ holds the cycle-level observational authority that the assessment requires.

**β owns `git merge` into main and the β close-out.** β's close-out captures the review context and merge evidence. **δ owns tag/release/deploy** (the release boundary per `CDD.md` §1.4 β algorithm and `operator/SKILL.md` §3.4). The PRA is a separate artifact written by γ after β's merge and close-out are complete.

**Handoff:** β completes release + β close-out → γ reads both close-outs (α + β) + the shipped artifacts → γ writes the PRA. If γ's session ends before the assessment is complete, the assessment is the first task of its next session.

Exception: the operator may explicitly reassign the assessment to another agent. The reassignment must name the target agent and the reason.

## When

After every `git tag` + `gh release create`. No exceptions — patch, minor, or major.

## Output

The assessment produces one artifact at the canonical path `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` (for the CDD package: `docs/gamma/cdd/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md`). All version strings in the artifact are bare `X.Y.Z`; do not use a `v` prefix.

The artifact has the following sections:

```markdown
## Post-Release Assessment — X.Y.Z

**CI status on merge SHA:** green / red / pending — [run URL]

### 1. Coherence Measurement
- **Baseline:** PREV — α _, β _, γ _
- **This release:** X.Y.Z — α _, β _, γ _
- **Delta:** which axes improved / held / regressed and why
- **Coherence contract closed?** Expected effect achieved? If not, what remains?

### 2. Encoding Lag
| Issue | Title | Type | Design | Impl | Lag |
|-------|-------|------|--------|------|-----|

**MCI/MCA balance:** balanced / freeze MCI / resume MCI
**Rationale:** ...

### 3. Process Learning
**What went wrong:** ...
**What went right:** ...
**Skill patches:** (committed Y/N, link if Y)
**Active skill re-evaluation:** (for each review finding: skill underspecified → patch / application gap → note / already covered → skip)
**CDD improvement disposition:** (patch landed: [description] / no patch needed: [justification])

### 4. Review Quality
**Cycles this release:** N
**Avg review rounds:** N.N (target: ≤1 docs, ≤2 code)
**Superseded cycles:** N (target: 0)

**Per-cycle round counts:**

| Cycle | Issue | Mode | Rounds | Binding findings (R1) | Notes |
|-------|-------|------|--------|----------------------|-------|

**Per-cycle dispatch telemetry** (optional initially; mandatory after ~10 cycles of data accumulate to validate the §1.6c heuristic):

| Cycle | `dispatch_seconds_budget` | `dispatch_seconds_actual` | `commit_count_at_termination` |
|-------|--------------------------|--------------------------|-------------------------------|

- `dispatch_seconds_budget` — the timeout budget set in the dispatch prompt (per `CDD.md §1.6c(a)`)
- `dispatch_seconds_actual` — wall-clock seconds from dispatch start to agent exit (SIGTERM or clean)
- `commit_count_at_termination` — number of commits α pushed to the cycle branch before session end

These three fields accumulate per-cycle data so the heuristic constants in `CDD.md §1.6c(a)` (`120s × ac_count` docs floor, `180s × ac_count` code floor) can be tightened from informed-guess to empirically-validated after sufficient cycles. If `commit_count_at_termination = 0` and `dispatch_seconds_actual < dispatch_seconds_budget`, the agent was SIGTERM'd without checkpointing — consult `operator/SKILL.md §timeout-recovery`.

**Finding-class breakdown** (across cycles in this release):

| Class | Definition | Count |
|---|---|---|
| **mechanical** | Caught by grep/diff/script | N |
| **wiring** | "X is wired into Y" but isn't (see review/SKILL.md 3.13c) | N |
| **honest-claim** | Doc claims something code/data doesn't back (review/SKILL.md 3.13) | N |
| **judgment** | Design/coherence assessment | N |
| **contract** | Work contract incoherent | N |

**Mechanical ratio:** N% (threshold: 20% → file process issue)
**Honest-claim ratio:** N% (target: <30% — high ratio means α docs are drifting from artifacts; patch by tightening review/SKILL.md 3.13 application or by improving α self-coherence templates)
**Action:** none / filed #NN

### 4a. CDD Self-Coherence
- **CDD α:** _/4 — (artifact integrity)
- **CDD β:** _/4 — (surface agreement)
- **CDD γ:** _/4 — (cycle economics)
- **Weakest axis:** α / β / γ
- **Action:** none / patch skill / patch doc / automate check

### 4b. Cycle Iteration
- **Triggered by** (per CDD.md §9.1 thresholds):
  - review rounds > 2
  - mechanical ratio > 20% with ≥ 10 findings
  - avoidable tooling/environmental failure
  - CI red on merge commit (post-merge)
  - loaded skill failed to prevent a finding
  - none
- **Root cause:** ...
- **Disposition:** patch landed now / next MCA #NN / no patch with reason
- **Evidence:** commit / issue / note

### 5. Production Verification

**Scenario:** [what to test]
**Before this release:** [what happened]
**After this release:** [what should happen]
**How to verify:** [concrete steps]
**Result:** [pass / fail / deferred]

### 6. CDD Closeout

| Step | Artifact | Skills loaded | Decision |
|------|----------|--------------|----------|
| 11 Observe | observation surface | post-release | runtime/design alignment result |
| 12 Assess | POST-RELEASE-ASSESSMENT.md | post-release | assessment completed |
| 13 Close | immediate fixes / next MCA issue | post-release (+ others if used) | cycle closed / deferred outputs committed |

### 6a. Invariants Check

If the project maintains an architectural invariants document, confirm which constraints were touched this cycle and their status:

| Constraint | Touched? | Status (preserved / tightened / revised / N/A) |
|---|---|---|

If no invariants document exists, omit this section.

### 7. Next Move
**Next MCA:** #NN — title
**Owner:** ...
**Branch:** ...
**First AC:** ...
**MCI frozen until shipped?** yes / no
**Rationale:** ...

**Closure evidence (CDD §10):**
- Immediate outputs executed: yes / no
  - (list each with link/commit)
- Deferred outputs committed: yes / no
  - (issue # / owner / branch / first AC / freeze state for each)

**Immediate fixes** (executed in this session):
- ...

### 8. Hub Memory
- **Daily reflection:** [path] — committed at [sha]
- **Adhoc thread(s) updated:** [path(s)] — committed at [sha]
```

## Procedure

### Step 1: Score

Score α/β/γ for the release. Compare to baseline (previous release score from CHANGELOG TSC table).

Rules:
- Score the release, not the intent
- If an axis regressed, name the cause
- If an axis held, say whether that's expected or stagnation

### Step 2: Update CHANGELOG TSC table

Add a row in the canonical bare-version format. The Level and Rounds columns are required (per CDD.md §9.1 and `release/SKILL.md` §2.4):

```
| X.Y.Z | C_Σ | α | β | γ | Level | Rounds | Coherence note |
```

The coherence note describes which incoherence was reduced, not what feature was added. The Level column records the cycle-level engineering level (L5 / L6 / L7) per CDD.md §9.1. The Rounds column records the review-round count for the cycle (e.g. `1`, `2`, `3`); for releases bundling multiple cycles, sum or list (`1+2`).

**Scoring sequence:** The CHANGELOG TSC entry written at release time is **provisional** — it is β's release-time score (marked as `provisional, pending γ PRA` in the level cell per `release/SKILL.md` §2.4). The post-release assessment is γ's independent score and MUST revise the CHANGELOG entry. γ updates the provisional TSC row to the final scoring values in the same commit as the PRA, replacing β's provisional markers with the final assessment. The assessment governs.

### Step 3: Encoding lag table

All converged-but-unimplemented design commitments MUST appear in the lag table. This is not optional and has no wiggle room — if a design is converged and not shipped, it appears here.

For every open design issue (issues with design docs, architecture docs, or converged plans that are not yet fully implemented) and every open process issue (repeatable failure modes from §9–§10):

| Issue | Title | Type | Design | Impl | Lag |
|-------|-------|------|--------|------|-----|

Type: `feature` (design/code gap) or `process` (development method gap).

Lag levels:
- **none** — shipped in this or prior release
- **low** — implementation in progress (branch exists, `.cdd/unreleased/{N}/self-coherence.md` present)
- **growing** — design converged, no implementation started
- **stale** — design aging without implementation plan

### Step 4: MCI/MCA balance decision

Based on the lag table:
- **Balanced:** roughly equal design and implementation activity. Continue normally.
- **Freeze MCI:** ≥3 issues at "growing" lag, OR designs outnumber implementations 3:1, OR any "SHALL" in substrate docs without runtime enforcement. Stop designing, start shipping.
- **Resume MCI:** implementation caught up to design frontier. Can advance designs again.

This decision is **mandatory**. Every release states the balance.

**What "freeze MCI" means operationally:** No new substantial design docs, plans, or architecture proposals until the committed MCA backlog is reduced below the freeze threshold. Small clarifications to existing designs are allowed; new design commitments are not.

### Step 5: Process learning

Answer three questions:
1. **What went wrong?** What broke, was caught late, or slipped through?
2. **What went right?** What process improvement paid off?
3. **Skill patches needed?** If a repeatable failure mode is identified, patch the skill NOW — not next session. This is CDD step 13a: the patch must land in the same commit as the assessment, synced across all affected surfaces: the canonical source under `src/packages/`, any package-visible loader entrypoint if affected, and any human-facing pointer surface that exposes the changed rule. A "noted for next cycle" commitment is insufficient when the failure mode is recurring — two cycles of the same mechanical failure with no spec-level fix is the trigger.
4. **Active skill re-evaluation.** For each review finding: would the declared active skills, as written, have prevented it? If not, is the skill underspecified for this pattern, or was it just not applied deeply enough? Underspecified → patch the skill (step 13a). Application gap → note it but don't patch (the skill was right, the agent didn't follow it).
5. **If no CDD improvement is possible**, state so explicitly with reasoning. "No skill patch needed" requires a justification: either (a) all findings were application gaps with adequate existing spec, (b) zero review findings this cycle, or (c) the failure mode is environmental (e.g. tooling constraint) with no spec-level fix available. Silence is not an acceptable disposition — every cycle must close the self-learning loop with either a patch or an explicit reason why not.

### Step 5.5: Review quality

For every triadic cycle in this release, record:
1. **Review rounds** — how many iterations before merge (fix-round appendices in `.cdd/unreleased/{N}/self-coherence.md` + RC verdicts in `.cdd/unreleased/{N}/beta-review.md`). Target: ≤1 for docs cycles, ≤2 for code cycles.
2. **Superseded cycles** — count of branches abandoned-not-merged and replaced. Target: 0.
3. **Finding taxonomy** — tag each review finding as `mechanical` (automatable: stale cross-refs, missing scope items, wrong branch name) or `judgment` (design coherence, architecture trade-offs).
4. **Mechanical ratio** — mechanical findings / total findings. If >20% AND total findings ≥ 10, file an issue to add the missing pre-flight check. Below 10 findings the ratio is noise — note it but don't file.

This step closes the loop from CDD §9 (Assessment). The review quality section in the output template must be filled.

### Step 5.6: CDD self-coherence

Did this cycle itself follow CDD coherently? Score each axis 1–4. The branch-level SELF-COHERENCE.md uses the template at `docs/gamma/cdd/SELF-COHERENCE-TEMPLATE.md`; this assessment-level check uses the same triadic axes:

- **CDD α (artifact integrity):** required artifacts present? Bootstrap/frozen snapshot complete? Self-coherence present?
- **CDD β (surface agreement):** canonical doc, executable skill, `.cdd/unreleased/{N}/` cycle artifacts (per the canonical filename set in `CDD.md` §Tracking), changelog, and assessment agree? Authority conflicts or stale references?
- **CDD γ (cycle economics):** review rounds within target? Superseded PRs low? Mechanical ratio under threshold? Immediate outputs executed, deferred outputs committed?

```markdown
### 4a. CDD Self-Coherence

- **CDD α:** _/4 — (rationale)
- **CDD β:** _/4 — (rationale)
- **CDD γ:** _/4 — (rationale)
- **Weakest axis:** α / β / γ
- **Action:** none / patch skill / patch doc / automate check
```

This drives action, not score-keeping:
- low α → patch artifact contract / templates / bootstrap
- low β → patch canonical doc / executable skill / authority hierarchy
- low γ → automate mechanical checks, reduce ceremony, tighten selection

If no axis scores below 3, write "no action" and move on.

### Step 5.6a: Cycle iteration

If any `CDD.md` §9.1 trigger fired, this section is mandatory.

For each fired trigger:
1. name the trigger
2. name the root cause
3. name the disposition:
   - patch landed now
   - next MCA committed
   - no patch with explicit reason
4. link the evidence (commit / issue / note)

If no `CDD.md` §9.1 trigger fired, write either:
- `No §9.1 trigger fired`, or
- the independent γ process-gap result if step 13 still found something worth patching.

### Step 5.6b: Author `cdd-iteration.md` (when applicable)

> Activation findings from `cdd/activation/SKILL.md §22` flow into this step. See [`cdd/activation/SKILL.md §22`](../activation/SKILL.md#22-cdd-iteration-cadence) for the per-repo cadence declaration, severity scale (D/C/B/A + `info`), and auto-spawn MCA trigger rule.

When the close-out triage in `gamma-closeout.md` produces ≥1 finding tagged `cdd-skill-gap`, `cdd-protocol-gap`, `cdd-tooling-gap`, or `cdd-metric-gap`, γ authors `.cdd/unreleased/{N}/cdd-iteration.md` — ε's work product (`ROLES.md §1`; `cdd/epsilon/SKILL.md §1`) — extracting only those findings into a structured form. This is the canonical home for cdd-self-improvement findings; PRA §3 prose remains as narrative, but the structured record lives here.

See [`cdd/activation/SKILL.md §22`](../activation/SKILL.md#22-cdd-iteration-cadence) for the cadence declaration: every cycle produces `cdd-iteration.md`; an empty findings list is itself signal.

**Per-finding shape:**

```markdown
### F1: <title>

- **Source:** α-closeout / β-review / γ-triage / PRA §3 / PRA §4b — name the artifact and section
- **Class:** `cdd-skill-gap` | `cdd-protocol-gap` | `cdd-tooling-gap` | `cdd-metric-gap`
- **Trigger:** §9.1 trigger N | "γ process-gap check" | "review pattern across cycles"
- **Description:** one paragraph
- **Root cause:** one paragraph
- **Disposition:** `patch-landed` | `next-MCA` | `no-patch`

If `patch-landed`:
- **Patch:** `<commit-sha>` or `<file-path>`
- **Affects:** `<cdd skill file or section>` (e.g. `cdd/review/SKILL.md` §3.13)
- **Cross-repo:** `<target-repo>` PR `#NN` if patch landed in a different repo than this cycle (otherwise omit)

If `next-MCA`:
- **Issue filed:** `<repo>` `#NN`
- **First AC:** ...

If `no-patch`:
- **Reason:** required justification per CDD §13a
```

**Aggregator update.** After writing `cdd-iteration.md`, γ updates `.cdd/iterations/INDEX.md` with one row:

```markdown
| Cycle | Issue | Date | Findings | Patches | MCAs | No-patch | Path |
|-------|-------|------|----------|---------|------|----------|------|
| {N}   | #N    | YYYY-MM-DD | M | P | A | Z | .cdd/releases/{X.Y.Z}/{N}/cdd-iteration.md |
```

For docs-only releases, the path uses the §2.5b form: `.cdd/releases/docs/{ISO-date}/{N}/cdd-iteration.md`.

**Cross-repo trace.** If any finding's disposition is `patch-landed` with a `Cross-repo` target, γ also creates `.cdd/iterations/cross-repo/{target-repo}/{slug}/` containing:

- the bundle deliverables (patch series, issue body, PR body, apply notes)
- a `LINEAGE.md` mapping each deliverable to its source cycle(s) and the target repo's PR

The cross-repo directory persists until the target PR merges; γ may delete it thereafter (the lineage is preserved in the target repo's own `cdd-iteration.md`).

**Why this artifact exists.** Without `cdd-iteration.md`, cdd-self-improvement findings are buried in PRA §3 prose. They cannot be aggregated, the protocol cannot measure its own learning rate, and cross-repo work loses traceability. Making the artifact first-class closes cdd's self-iteration loop with the same discipline cdd applies to system-level cycles. Every cycle produces the file — an empty findings list is itself signal that the cycle ran cleanly.

### Step 5.7: Production verification

After release, verify the change works in production — not just that CI passes. Design a concrete test that demonstrates the new capability (or blocked failure mode) in a real environment.

**The question:** What can the system do now that it couldn't before? Or: what failure is now impossible that was previously possible?

Write a verification scenario:

```markdown
### Production Verification

**Scenario:** [what to test]
**Before this release:** [what happened]
**After this release:** [what should happen]
**How to verify:** [concrete steps]
**Result:** [pass / fail / deferred — with evidence]
```

Rules:
- The scenario must be executable, not hypothetical
- "CI passes" is not production verification — it's build verification
- If the change is structural (L7), the test should demonstrate the boundary move, not just a single path
- If verification requires a running agent or daemon, say so — defer if the environment doesn't support it, but commit to when/how it will be verified
- If deferred, add to the next session's checklist

This is the kata: each release earns a concrete demonstration that the system changed.

### Step 6: Decide next move

Based on measurement + lag + learning, state what happens next as a **concrete commitment**:

- **Next MCA issue:** the specific issue number to implement next
- **Owner:** who is responsible for the next MCA
- **Branch:** name of branch (or "pending branch creation")
- **First AC to close:** which acceptance criterion ships first
- **MCI frozen?** Whether MCI is frozen until this MCA ships

This turns the assessment into an executable handoff, not just reflection.

**Execution rule:**
- Small process/skill fixes discovered in the assessment → execute immediately (same session)
- Everything larger → becomes the next cycle's work via the commitment above
- Do not attempt to execute a large MCA inside the post-release assessment itself

### Step 7: Hub memory

Every release changes something worth remembering across sessions. Write both:

1. **Daily reflection** — what happened this cycle, scoring, MCI freeze status, what's next. This is the operational state that the next session loads to orient.
2. **Adhoc thread update** — which ongoing thread(s) this release advances. Every release connects to at least one active thread (refactor progress, skill convergence, process evolution, etc.). Update the relevant thread with what shipped and what it means for the thread's arc. If the release doesn't connect to any existing thread, that's a signal worth noting — start one or name why.

Both writes happen **before** the cycle is considered closed. The assessment artifact lives in the repo; the hub memory is what makes it findable and contextual across sessions.

**Why mandatory:** Sessions are stateless. Without hub memory, the next session must re-derive cycle context from git log and assessment files. The daily reflection is the index; the adhoc thread is the narrative. Skipping either creates a compaction gap — proven by v3.41.0 where the assessment was written but hub memory was not, and the next session lost context.

Anti-pattern: ❌ "I'll write the reflection later" (compaction or session end erases the intent).

## Pre-publish gate

Before committing the assessment, verify mechanically:

- [ ] §1 has Baseline, This release, Delta, and Coherence contract closed
- [ ] §2 has encoding lag table with every open design/process issue
- [ ] §2 has MCI/MCA balance decision with rationale
- [ ] §3 has What went wrong, What went right, Skill patches, Active skill re-evaluation, and explicit disposition on CDD improvement (patch landed OR "no patch needed" with justification)
- [ ] §4 has all fields: cycles, review rounds, superseded cycles, finding breakdown, mechanical ratio, action
- [ ] §4 mechanical ratio: if >20% AND total findings ≥ 10, a process issue is **filed and referenced** (not just noted). Below 10 findings, note the ratio but no filing required.
- [ ] §4a CDD self-coherence: α/β/γ scored, weakest axis named, action stated (or "none" if all ≥3)
- [ ] If any `CDD.md` §9.1 trigger fired, §4b Cycle Iteration exists with trigger, root cause, disposition, and evidence.
- [ ] If close-out triage produced ≥1 finding tagged `cdd-skill-gap` / `cdd-protocol-gap` / `cdd-tooling-gap` / `cdd-metric-gap`, `.cdd/unreleased/{N}/cdd-iteration.md` exists (Step 5.6b) **and** `.cdd/iterations/INDEX.md` has a new row for cycle N. If patch landed cross-repo, `.cdd/iterations/cross-repo/{target}/{slug}/` exists with `LINEAGE.md`.
- [ ] §3/§4 skill patches: if §3 identifies a recurring failure mode or skill gap, the patch is **in this commit** (not deferred) and synced across all affected surfaces: canonical source under `src/packages/`, package-visible loader entrypoint if affected, human-facing pointer/readme surfaces if they expose the changed rule.
- [ ] §5.7 has production verification scenario (or explicit deferral with commitment)
- [ ] §6 CDD Closeout trace present with rows for observe/assess/close steps (incl. step 13a if skill patches landed)
- [ ] §7 has all 6 fields: Next MCA, Owner, Branch, **First AC**, MCI frozen?, Rationale
- [ ] §7 Closure evidence: immediate outputs listed with links, deferred outputs with issue/owner
- [ ] CHANGELOG TSC row added or updated to match assessment scores
- [ ] §8 Hub memory: daily reflection written and pushed to hub repo
- [ ] §8 Hub memory: at least one adhoc thread updated (or explicit note why none applies)
- [ ] If the triadic protocol is active, `.cdd/unreleased/{N}/gamma-closeout.md` (in-version) or `.cdd/releases/{X.Y.Z}/{N}/gamma-closeout.md` (post-release) exists and reflects α close-out, β close-out, γ triage, and final cycle status. (Legacy aggregate `.cdd/releases/{X.Y.Z}/gamma/CLOSE-OUT.md` is warn-only per `CDD.md` §5.3a; do not require it.)

This gate is mechanical. Two agents checking the same template must find the same missing fields.

## Anti-patterns

- ❌ Skip assessment for patch releases ("it's just a small fix")
- ❌ Score without the lag table ("coherence is fine" while 5 designs age)
- ❌ Note process failures without patching skills ("we should fix this")
- ❌ Ship and immediately start the next feature without assessing balance

## Examples

### Good assessment (MCI freeze triggered)

```
## Post-Release Assessment — v3.12.2

### Coherence Measurement
- Baseline: v3.12.1 — α A, β A, γ A
- This release: v3.12.2 — α A, β A+, γ A
- Delta: β improved (contract authority enforced). α/γ held.

### Encoding Lag
| #62 | RT Contract v2 | converged | branch exists | low |
| #65 | Communications | converged | not started | growing |
| #73 | Extensions | converged | not started | growing |
| #67 | Network | subsumed by #73 | not started | growing |

MCI/MCA balance: **Freeze MCI** — 3 issues at growing lag.
No new design docs until backlog reduced below threshold.

### Process Learning
Wrong: Review missed CAA.md (§2.0 gate not followed). Fixed with structural table format.
Right: Three-agent review loop caught complementary gaps.
Skill patches: review §2.0 gate (committed), CDD §7.6 output format (committed).

### Next Move
Next MCA: #62 — Runtime Contract v2
Owner: sigma
Branch: claude/runtime-contract-v2-VWKUT
First AC: CAA.md updated with wake-time architecture
MCI frozen until shipped? Yes
Rationale: Vertical self-model is foundation for #73, #65, #59.

Immediate fixes (executed this session):
- review §2.0 structural gate (eeca922)
- CDD §7.6 output format (64634fb)
```

---

## Identity migration

Cycle #343 (merged 2026-05-11) is the cutover from the deprecated `{role}@cdd.{project}` form to `{role}@{project}.cdd.cnos` (cnos elision: `{role}@cdd.cnos`). Git history is immutable; commits made before the cutover retain their original trailers. From cycle #343 onward all new commits must use the new form. Cycles that spanned the cutover may show mixed trailer forms in their history — this is acceptable during the transition window. Close-out artifacts may reference both forms transparently.

## Kata

**Scenario:** v3.29.1 just shipped. Write the post-release assessment.

1. Restate what was released and what coherence delta was targeted
2. Grade TSC honestly — where did the release fall short?
3. Identify immediate outputs (skill patches, doc fixes) and execute them now
4. Identify deferred outputs and commit them as next-cycle work
5. Check if §9.1 triggers fired (review rounds > 2, mechanical ratio > 20%, etc.)

**Verify:** Are all immediate outputs executed, not just listed? Are deferred outputs recorded concretely? Is the cycle closed?
