---
name: gamma
description: γ role in CDD. Observes, selects, dispatches, unblocks, runs post-release assessment, and closes the cycle.
artifact_class: skill
kata_surface: embedded
governing_question: How does γ keep the full cycle coherent across issue creation, dispatch, unblocking, assessment, and close-out?
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
  - unblock decisions
  - post-release assessment
  - gamma close-out triage
  - cycle closure
requires:
  - active role is γ
  - canonical CDD.md loaded
calls:
  - issue/SKILL.md
  - post-release/SKILL.md
  - operator/SKILL.md
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
3. load `issue/SKILL.md`
4. load `post-release/SKILL.md` — γ owns the PRA (cycle-level assessment of α, β, and cycle economics) and step 13a skill/spec patches
5. load `operator/SKILL.md` — δ owns release-phase gate execution (tag push, branch cleanup, release CI) and the disconnect release (§3.4). δ's actions are git-observable (tags, branch state). If δ is unavailable, γ may execute gates directly.
6. load other lifecycle sub-skills only when the selected gap requires them

Canonical artifact locations (PRA, close-out paths, snapshot dirs, tag policy) are defined in `CDD.md` §5.3a (Artifact Location Matrix).

`CDD.md` is the only canonical source for:
- the ordered γ lifecycle (`CDD.md` §1.4, steps 1–16; step 17 is δ's)
- selection rule order (`CDD.md` §3)

This file does **not** redefine that algorithm.
It expands those steps into executable checks, evidence, and gates.

## Step map

- Step 1 → git identity (`§2.0`), Step 2a → observe and build candidates (`§2.1`), Step 2b → select and size (`§2.2`)
- Step 3 → issue pack + issue-quality gate (`§2.3–§2.4`)
- Step 3a → create the cycle branch (`§2.5` Step 3a — `cycle/{N}` from `origin/main`, with γ-owned branch pre-flight per `CDD.md` §4.3, before any dispatch)
- Steps 3b–6 → subscribe, dispatch + unblocking (`§2.5` Step 3b)
- Steps 7–8 → deferred release mechanics (`§2.6`)
- Steps 9–10 → close-out triage (`§2.7`)
- Steps 11–13 → cycle iteration and process patching (`§2.8–§2.9`)
- Steps 14–16 → hub memory, branch cleanup, closure declaration (`§2.10`)
- Step 17 → δ disconnect release (`operator/SKILL.md` §3.4) — δ's step, not γ's

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

### 2.1. Step 1a — Observe and build candidates

Before selecting work, read the observation surfaces required by `CDD.md`:
1. **Last post-release assessment** — read this first. It contains the prior cycle's next-MCA commitment, deferred outputs, cycle iteration findings, and MCI freeze state. These are binding inputs to selection, not optional context.
2. CHANGELOG TSC table
3. encoding lag table
4. doctor / status / operational-health surface

Build a candidate table:

```md
| Candidate | Source | Rule clause that nominates it | Dependency | Size | Decision |
|-----------|--------|-------------------------------|------------|------|----------|
| #NN ...   | lag    | CDD §3.x                      | ...        | ...  | ...      |
```

Do not select from memory or preference alone.

- ❌ "I remember issue #143 felt important"
- ✅ "Candidate table shows #143 is selected under CDD §3.x for a named reason"

### 2.2. Step 1b — Select by `CDD.md`, then size the intervention

Apply `CDD.md` §3 in order.
Do **not** restate or reorder the rule list here.

Required output of selection:
- selected gap
- decisive `CDD.md` rule / clause
- rejected alternatives when non-obvious
- intervention size: immediate output / small change / substantial cycle

Sizing rule:
- if the fix is immediate-output sized, execute it now and continue observation
- if the work qualifies for small-change, route it to the small-change path instead of opening a substantial cycle
- if the selected candidate is blocked, name the blocking dependency in the record before dispatch

- ❌ "Selected because it felt highest leverage"
- ✅ "Selected #143 under CDD §3.x; #151 rejected because dependency unresolved; script fix executed immediately instead of becoming the cycle"

### 2.3. Step 2 — Build the issue pack

Use `issue/SKILL.md` as the base contract.
A dispatchable γ issue is:
- a complete `issue/SKILL.md` issue
- **plus** a work-shape note (`substantial` / `small-change` / `immediate-output`)
- **plus** dependency notes when sequencing or blockers are real

γ does not restate Tier 1 or Tier 2 skills in the issue.
γ does name Tier 3 skills, active design constraints, related artifacts, non-goals, and priority exactly as `issue/SKILL.md` requires.

If the issue cannot be written to that level, the work is not ready for α dispatch.

- ❌ "Fix package restore; it's incoherent"
- ✅ "Gap, evidence, numbered ACs, non-goals, Tier 3 skills, invariants, related artifacts, priority, work-shape note, dependency note"

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

Do not compensate for a weak issue by making the prompt longer.
Fix the issue instead.

### 2.5. Steps 3a–5 — Create the cycle branch, dispatch, unblock without leakage

#### Step 3a — Create the cycle branch

After the issue passes the issue-quality gate (§2.4) and **before** dispatching α and β, γ creates the cycle branch. The branch is the canonical coordination surface during the cycle (`CDD.md` §Tracking + §4.2); a single named target (`cycle/{N}`) replaces the pre-#287 model where α opened the branch and γ + β had to glob-discover it.

**γ-owned branch pre-flight (`CDD.md` §4.3):**

- `origin/cycle/{N}` does not yet exist (`git rev-parse --verify origin/cycle/{N}` returns non-zero — fail loud if it already exists, since one cycle = one branch)
- no stalled `.cdd/unreleased/{N}/` directory exists on `origin/main` (would indicate a previous cycle for `{N}` did not complete its release-time move per `release/SKILL.md` §2.5a)
- the issue's intended scope is declared in the issue body
- base SHA is known (`git rev-parse origin/main`)
- the issue is open

**Creation snippet:**

```bash
git fetch --quiet origin main
# pre-flight: branch does not yet exist
if git rev-parse --verify origin/cycle/<N> >/dev/null 2>&1; then
  echo "pre-flight fail: origin/cycle/<N> already exists" >&2
  exit 1
fi
# pre-flight: no stalled .cdd/unreleased/{N}/ on main
if git ls-tree -r --name-only origin/main .cdd/unreleased/<N>/ 2>/dev/null | grep -q .; then
  echo "pre-flight fail: stalled .cdd/unreleased/<N>/ on main" >&2
  exit 1
fi
# create + push
git switch -c cycle/<N> origin/main
git push -u origin cycle/<N>
```

The branch must exist on `origin` before dispatch. α and β never create branches; their prompts include a `Branch: cycle/<N>` line and they `git switch` to it (`alpha/SKILL.md`, `beta/SKILL.md` §1, `CDD.md` §1.4).

#### Step 3b — Subscribe and dispatch α and β

Immediately begin polling the issue and `origin/cycle/{N}` (see `CDD.md` §Tracking for query forms and the wake-up mechanism) — do not ask, just do it. γ must track the full cycle: issue activity, `origin/cycle/{N}` head SHA, α's `.cdd/unreleased/{N}/self-coherence.md` updates, β's `.cdd/unreleased/{N}/beta-review.md` verdicts, and branch CI status. Start polling before dispatching.

**Polling requires a query, a wake-up mechanism, and a reachability probe.** Picking a query form (`gh issue`, `git fetch`, `git ls-tree`) without confirming a wake-up form (`Monitor` stdout-as-notification or shell-wake-on-loop-exit) is silent — the loop runs but γ never reacts. Verify both before dispatch. If the environment provides neither a `Monitor`-equivalent nor a shell-wake harness, surface the gap to the operator before dispatching α and β; γ cannot autonomously coordinate without a wake-up contract.

For MCP-only γ environments (e.g. Claude Code on the web), the canonical wake-up shape is `Monitor` wrapping a transition-only stdout filter against the **single named cycle branch** (no glob — γ created it in Step 3a and knows its name):

```bash
# Baseline sync — run BEFORE the transition loop.
git fetch --quiet origin cycle/<N>
git rev-parse --verify origin/cycle/<N>
echo "baseline-cycle-dir: $(git ls-tree -r --name-only origin/cycle/<N> .cdd/unreleased/<N>/ 2>/dev/null | tr '\n' ' ')"

# Transition loop — single named branch.
prev_head=""; empty_iters=0
while true; do
  git fetch --quiet origin cycle/<N> || echo "fetch-error: cycle/<N>"
  cur_head="$(git rev-parse origin/cycle/<N> 2>/dev/null)"
  if [ "$cur_head" != "$prev_head" ] && [ -n "$cur_head" ]; then
    echo "branch-update: cycle/<N> → $cur_head"
    prev_head="$cur_head"
    empty_iters=0
  else
    empty_iters=$((empty_iters + 1))
    # CDD.md §Tracking git fetch reliability rule: every 10 empty iterations, do a synchronous re-probe.
    if [ "$empty_iters" -ge 10 ]; then
      git fetch --verbose origin cycle/<N> 2>&1 | tee /tmp/cycle-<N>-fetch.log >&2 || \
        echo "reachability-fail: cycle/<N> — surface to operator"
      empty_iters=0
    fi
  fi
  sleep 60
done
```

Each transition line becomes a `task-notification` that wakes the session. Run under `Monitor`; emit only on transition. **All cycle-dir artifacts live on `origin/cycle/{N}` — not on `main`** (per `CDD.md` §Tracking). Polling `origin/main` for `.cdd/unreleased/{N}/` is silent for in-flight cycles; γ dereferences the cycle directory on the named branch via `git ls-tree -r origin/cycle/{N} .cdd/unreleased/{N}/` when a transition fires.

Then produce both prompts at dispatch time. β begins polling the issue and `.cdd/unreleased/{N}/self-coherence.md` and starts intake while α implements — β does not need to wait for review-readiness to begin polling.

**α prompt:**
```text
You are α. Project: <project>.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md.
Issue: gh issue view <N>
Branch: cycle/<N>
Tier 3 skills: <list issue-specific skills>
```

**β prompt:**
```text
You are β. Project: <project>.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view <N>
Branch: cycle/<N>
```

Rules:
- point both roles at the issue, not a paraphrase of the issue
- include the explicit `Branch: cycle/<N>` line so α and β never have to invent or glob-discover the branch
- do not restate the algorithm in the prompt
- do not smuggle missing constraints into chat prose; fix the issue instead
- β begins polling the issue and `.cdd/unreleased/{N}/self-coherence.md` on `origin/cycle/{N}` immediately; the β skill handles waiting for α's review-readiness signal
- β receives artifact surfaces, not α's hidden implementation rationale

#### Step 5 — Unblock

When α or β is blocked, γ may:
- clarify requirement wording
- add missing artifact links
- edit the issue to state an omitted invariant
- resolve ambiguity in scope
- provide mechanical environment help
- point the role back to the governing skill or artifact

γ may **not**:
- forward β's internal reasoning transcript to α
- forward α's hidden rationale transcript to β
- author the implementation fix inside the review loop
- silently change the target gap without updating the issue / artifact

Allowed transfer unit: **artifact facts**, not hidden role state.

- ❌ "β said this design is shaky; just rewrite the parser like this"
- ✅ "The issue omitted invariant X; γ adds it to the issue and points α back to the updated artifact"

### 2.6. Steps 6–7 — Support deferred release mechanics only

If β deferred a mechanical release step because of environment constraints, δ (operator) executes the gate action. γ observes the result via git (tag appears, branch deleted). If δ is unavailable, γ may execute directly:
- push the tag
- verify release CI fired
- close the issue if auto-close failed

γ does not redo β's judgment.
γ only completes deferred mechanics. When δ handles the gate, γ observes completion via git before proceeding to close-out triage.

### 2.7. Steps 8–9 — Triage close-outs explicitly

Before close-out, collect (in-version, before release):
- `.cdd/unreleased/{N}/alpha-closeout.md` (α close-out narrative)
- `.cdd/unreleased/{N}/beta-closeout.md` (β close-out narrative + release evidence)

`self-coherence.md` and `beta-review.md` carry the in-cycle record (gap/ACs/trace and round-by-round verdicts respectively); the two `*-closeout.md` files are γ's primary triage inputs.

After release, the cycle directory moves to `.cdd/releases/{X.Y.Z}/{N}/` per `release/SKILL.md` §2.5a. The legacy aggregate paths `.cdd/releases/{X.Y.Z}/{alpha,beta,gamma}/CLOSE-OUT.md` are warn-only (pre-#283 form). See `CDD.md` §5.3a Artifact Location Matrix.

Then write the post-release assessment per `post-release/SKILL.md` at the canonical path `docs/{tier}/{bundle}/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md` (for the CDD package itself: `docs/gamma/cdd/{X.Y.Z}/POST-RELEASE-ASSESSMENT.md`). The PRA is γ's artifact — it measures α's implementation, β's review quality, and cycle economics. β assessing its own review is a self-grading problem.

γ triages all findings from both close-outs and the PRA. For each finding, record one disposition using CAP:
1. **Immediate MCA available** → ship now (γ lands the skill/spec patch per CDD §5.3 row 12a; if γ explicitly delegates the patch, name the delegate and the deadline in the triage record)
2. **Project MCI** → file / update project issue or `.cdd/` artifact
3. **Agent MCI** → update hub / adhoc thread
4. **One-off** → drop explicitly

Minimum triage record:

```md
| Finding | Source | Type | Disposition | Artifact / commit |
|---------|--------|------|-------------|-------------------|
| ...     | α/β/assessment | process/skill/... | immediate MCA / project MCI / agent MCI / drop | ... |
```

Silence is not triage.
Every finding gets a disposition.
Step 13a skill/spec patches are γ's to land or explicitly delegate — they do not become silent next-cycle work.

### 2.8. Steps 10–11 — Enforce cycle-iteration outputs when triggers fire

Apply the cycle-iteration checks named in `CDD.md` step 10.
For each fired trigger, γ must do something explicit:

| Trigger | Fire condition | Required γ action | Closure rule |
|---|---|---|---|
| Review churn | review rounds > 2 | Verify the assessment contains a `Cycle Iteration` entry naming the root cause and the chosen disposition (patch landed now / next MCA / no patch with reason). | If missing, cycle cannot close. |
| Mechanical overload | mechanical ratio > 20% **and** total findings ≥ 10 | Verify the assessment names the recurring mechanical class and whether the mechanization patch was landed now or filed as a concrete MCA. | If missing, cycle cannot close. |
| Avoidable tooling / environment failure | environment or tooling blocked the cycle in a way a guardrail could likely prevent | Verify the assessment names the failure, workaround, and disposition (guard landed now / issue filed / no spec-level fix with reason). | If missing, cycle cannot close. |
| Loaded-skill miss | a loaded skill should have prevented a finding but did not | Patch the skill now when the correction is clear; otherwise verify the assessment names the exact skill gap and the concrete next MCA. | If neither patch nor committed MCA exists, cycle cannot close. |

This is not just a checklist.
Each fired trigger must end in one of three states:
- patch landed now
- concrete next MCA committed
- explicit no-patch decision with reason

### 2.9. Step 13 — Run the independent γ process-gap check

Even if no `CDD.md` §9.1 trigger fired, γ must still ask:
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

### 2.10. Steps 13–15 — Close only after the closure gate passes

Do not declare the cycle closed until all of the following are true:
1. `.cdd/unreleased/{N}/alpha-closeout.md` exists on main
2. `.cdd/unreleased/{N}/beta-closeout.md` exists on main
3. γ has written the post-release assessment per `post-release/SKILL.md`
4. every fired cycle-iteration trigger has a `Cycle Iteration` entry with root cause and disposition
5. recurring findings were assessed for skill / spec patching
6. immediate outputs were either landed or explicitly ruled out
7. deferred outputs have issue / owner / first AC
8. next MCA is named
9. hub memory is updated
10. merged remote branches are cleaned up

Then:
- write `.cdd/unreleased/{N}/gamma-closeout.md`. The γ close-out contains: cycle summary, close-out triage table, §9.1 trigger assessment, cycle iteration, skill gap candidate dispositions, deferred outputs, hub memory evidence, and next MCA.
- update hub memory
- delete merged remote branches
- state closure explicitly: *"Cycle #N closed. Next: #M."* This is γ's last commit. δ will cut the disconnect release (step 17) — the tag appearing on main is the observable proof the cycle is fully closed.

---

## 3. Rules

### 3.1. Select by rule order, not taste

The candidate you like is irrelevant if a stronger `CDD.md` rule applies.

### 3.2. Name the decisive clause

Every selected gap must record the `CDD.md` clause that made it win.

### 3.3. Make the issue executable before dispatch

Prompt cleverness is not a substitute for issue quality.

### 3.4. Name only Tier 3 skills in the issue

Tier 1 and Tier 2 are already mandatory.
Repeating them hides the real issue-specific constraints.

### 3.5. Preserve epistemic separation

γ sees both sides because coordination requires it.
γ transfers artifact facts only.

### 3.6. Land immediate process fixes in the same cycle when possible

A missing gate discovered this cycle should not automatically become future work when the patch is already clear.

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

Select the next move and justify it using `CDD.md` rule order.

#### Expected answer

- immediate-output work executed now if truly immediate
- stale issue chosen for the next substantial MCA if the governing `CDD.md` clause forces it
- explicit decisive clause named

### Kata B — Issue quality

#### Scenario

An issue says only: "Fix package restore; it's incoherent."

#### Task

Rewrite it into a dispatchable γ issue pack.

#### Required fields

- concise gap
- evidence
- numbered ACs
- non-goals
- Tier 3 skills
- active design constraints
- related artifacts
- priority
- work-shape note
- dependency note when applicable

#### Common failures

- repeats Tier 1 / Tier 2 skills
- writes vague ACs
- omits non-goals
- leaves α to infer the invariant

### Kata C — Cycle iteration

#### Scenario

A cycle finished with:
- 3 review rounds
- 12 total findings, 5 mechanical
- one loaded-skill miss
- no process patch committed yet

#### Task

State what γ must require before closure.

#### Expected answer

- `Cycle Iteration` section added to the assessment
- root cause named for review churn / mechanical overload / skill miss
- disposition given for each fired trigger
- closure blocked until there is either a landed patch, a concrete next MCA, or an explicit no-patch justification
