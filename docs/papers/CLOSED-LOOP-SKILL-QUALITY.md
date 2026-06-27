# Closed-Loop Skill Quality

**Date:** 2026-03-25
**Context:** v3.17.0 release, PR #111 post-mortem, kata framework development

---

## The Problem

PR #111 (Runtime Extensions Phase 1) took 6 review rounds and had a 53% mechanical finding ratio — the second consecutive cycle at that level. The OCaml skill was loaded. The agent wrote `List.hd`, bare `with _ ->`, and `Hashtbl` mutation inside `List.map` anyway. Skills existed but didn't change behavior.

## The Diagnosis

Claude, the PR author, named it directly: "Stop reviewing against skills, start generating against them." The pattern was clear — skills were used as post-hoc checklists (write code → check against skill → find violations → fix), not as generation constraints that shape code as it's being written.

## The Method Patch

CDD §4.4: at step 5 (mode selection), name 2–3 governing skills as hard generation constraints. All others are reference only. Review checks consistency with the declared active skills. Findings that a loaded active skill would have prevented are process debt.

Design and review skills gained L7 checks: challenged assumptions, leverage/negative leverage, process economics, architecture leverage checks. Package copies synced across all surfaces.

## The Measurement Gap

But patching the method doesn't prove the method works. "We added §4.4" is not evidence that §4.4 reduces mechanical findings. We needed measurement.

## The Framework

The kata evaluation framework measures whether a skill changes the **shape of failure**, not just the score:

1. **Repair** — can the skill fix the immediate failure?
2. **Transfer** — does the skill help on a different but same-family task?
3. **Cost** — is the gain worth the extra context, time, and complexity?

Three experimental arms run in fresh sessions:

- **Cold** — no skills, baseline
- **All plausible skills** — skill soup control
- **Selected governing skills** — the 2–3 skills chosen before the attempt

If selected skills beat skill soup without blowing up cost, skill selection discipline is proven. If they only fix the same task but not a transfer task, the skill is prompt-local, not reusable.

## The Kata Packets

Three families matching the three kinds of engineering work:

- **Family A (Local implementation):** Registry conflict determinism, engine compatibility rejection. Tests coding + testing skills.
- **Family B (Cross-surface coherence):** Runtime contract doc/runtime parity review. Tests review + documenting skills.
- **Family C (Architecture leverage):** Capability growth boundary, browser + app ecosystem separation. Tests design + architecture-evolution + process-economics skills.

Each packet has a public task statement, public context, hidden evaluator rubric with α/β/γ/efficiency scoring, hidden checks, and a transfer variant.

The packets are drawn from real cnos problems — not academic exercises. They test exactly the engineering the agents actually do.

## The Worked Examples

Synthetic results from Claude across all three families show the same pattern:

| Family | Cold β | Selected β | Δ |
|--------|--------|-----------|---|
| A1 (implementation) | 2 | 3 | +1 |
| B1 (review) | 1 | 4 | **+3** |
| C1 (architecture) | 2 | 4 | **+2** |

B1 shows the largest β swing — the review skill transforms failure shape from "local code bias" to "cross-surface contract reasoning." All three classified as skill-effective. But these are hypothetical runs, not empirical data. The framework is ready; the evidence doesn't exist yet.

## The Failure Classification

Each kata outcome is classified as one of:

- **Knowledge-limited** — the model doesn't know something the skill can't teach
- **Process-limited** — the skill addresses the right gap
- **Selection-limited** — the wrong skill was chosen
- **Context-limited** — the skill helps but costs too much context
- **Architecture-limited** — the task needs a system change, not a skill

This classification is the exit condition for action: keep the skill if it transfers cheaply, tighten it if cost explodes, swap the selected skills if selection was wrong, and stop treating the problem as a skill problem when the blockage is architectural.

## The Closed Loop

The result is a closed loop for skill quality:

```
Write skill → test with kata → measure repair + transfer + cost
→ classify failure → keep / tighten / swap / stop → iterate
```

This is the missing piece between "we have skills" and "our skills work."

## The Meta-Move

The project's own failure data (PR #111, 53% mechanical ratio, type disambiguation bug across 3 rounds) was used to design the measurement framework for the project's own skills. The kata packets test the exact engineering problems the agents encounter. The framework measures the exact skills meant to prevent those problems.

cnos is measuring its own coherence tools with its own coherence methods. That's the loop closing.

## What Remains

1. **Run the katas with real agents.** The worked examples are synthetic. Real data is needed.
2. **Transfer testing.** Task B variants exist for all families. Repair alone doesn't prove skill quality.
3. **All-skills arm.** The "skill soup" control hasn't been run. The hypothesis that selected deep skills beat broad loading is untested.
4. **Longitudinal measurement.** Does the 53% mechanical ratio actually drop after §4.4 is in practice? That's the ultimate test.
