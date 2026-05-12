# α Dispatch — Cycle #345

```
You are α. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 345 --json title,body,state,comments
Branch: cycle/345
Tier 3 skills:
  - src/packages/cnos.core/skills/write/SKILL.md
  - src/packages/cnos.core/skills/skill/SKILL.md
```

## Context for δ

- Allowed tools: `--allowedTools "Read,Write,Bash"`
- Git identity for α commits on this cycle: `Alpha <alpha@cdd.cnos>`
- Mode: docs-only (design-and-build; design converged in issue body — no separate design doc).
- Branch: `cycle/345` — γ created it; α checks out, never creates.
- On completion α sets `status: review-ready` in `.cdd/unreleased/345/self-coherence.md` and commits + pushes to `cycle/345`.

## What α must implement (per issue #345)

**AC1 — `ROLES.md` at repo root, §§1–8 present.**

Placement decision: issue open question #1 recommends `ROLES.md` at repo root (matches `CHANGELOG.md` / `README.md` convention; no dir scaffolding needed). α follows the recommendation unless there is a concrete reason to deviate; note the decision in self-coherence.md.

Required sections (§§ keyed to AC1 oracle `rg '^## §' ROLES.md` → 8 hits):
- `## §1 The role ladder` — five-role table (exactly 5 rows); verbs: produces / reviews / coordinates / operates / iterates. Include the "self-application footnote" per issue open question #9: this doc itself was produced by α, reviewed by β, coordinated by γ, dispatched by δ, and ε's work emerged from the very findings it surfaces.
- `## §2 Scope escalation as nested observation` — ≤3 paragraphs; orders 0–4, one-sentence example each. Cite Bateson / von Foerster as ambient prior art, briefly.
- `## §3 Instantiation contract` — exactly 6 fields: matter type, review oracle, γ close-out artifact, δ cadence, ε iteration cadence, and one additional field α judges as necessary.
- `## §4 cdd as the reference instantiation` — link to `src/packages/cnos.cdd/skills/cdd/CDD.md`.
- `## §5 cdw as the planned sibling` — ≤200 words; four one-liners (α-prose, β-clarity/coherence/source-of-truth, γ-cycle, δ-pipeline); names cdw as "separate issue / future cycle" explicitly. No `cdw/` directory.
- `## §6 Naming convention for new c-d-X protocols` — letter scheme; cdd and cdw as the first two.
- `## §7 Role-name stability` — α/β/γ/δ/ε fixed across instantiations; verbs fixed; matter and oracle vary.
- `## §8 Glossary` — ≥6 terms: role, matter, frame, instantiation, scope-escalation, order-of-observation.

Oracles: `wc -l ROLES.md` → 250–500. `rg '^## §' ROLES.md` → 8 hits.

**AC2 — `cdd/CDD.md` top-of-file pointer.**

Add a single pointer at the top of `src/packages/cnos.cdd/skills/cdd/CDD.md` (within the first 20 lines) naming cdd as the reference instantiation of the role pattern and linking to `ROLES.md`. One pointer, no claim that the role structure is cdd-original.

Oracle: `head -20 src/packages/cnos.cdd/skills/cdd/CDD.md` contains `ROLES.md`.

**AC3 — `cdd/epsilon/SKILL.md` stub.**

Create `src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md` (30–100 lines). Required:
- Frontmatter: name, description, artifact_class, parent, scope, governing_question.
- §1: ε's cdd-side scope — protocol-iteration via `cdd-iteration.md` + MCA discipline. Cross-reference `ROLES.md §1` row 5 and `cdd/post-release/SKILL.md` Step 5.6b.
- §2 (or §1 paragraph): ε's relationship to δ — often the same actor; separation becomes warranted when protocol-iteration volume justifies dedicated attention.
- Negative invariant: no claim ε is required as a separate human/agent.

Oracle: `wc -l epsilon/SKILL.md` → 30–100.

**AC4 — `post-release/SKILL.md` Step 5.6b re-attribution.**

At `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` Step 5.6b, add a sentence (or amend the opening sentence) naming cdd-iteration.md as ε's work product, not γ's by-product. Single-sentence change; no other prose changes.

Oracle: `rg 'ε' src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` → hit at Step 5.6b.

**AC5 — Pattern self-applies (close-out attribution).**

This AC is γ's close-out obligation, not an α implementation task. α should note it as a constraint in self-coherence.md §Known Debt so γ sees it:

> "AC5 — cdd-iteration.md for this cycle must open with `(ε)` attribution when γ writes it at close-out. No α action required during implementation."

**AC6 — `ROLES.md §5` names cdw as planned-not-shipped.**

Already captured in AC1 §5 requirements above. Double-check: `rg 'separate issue\|future cycle' ROLES.md` → hit in §5.

## Non-goals (do not touch)

- Any `cnos:cdw/` directory — cdw is planned-not-shipped in this cycle.
- Renaming `cdd/operator/SKILL.md` (operator = δ, but the file keeps its name).
- Higher-order roles (ζ, η, …) — reserved but undefined.
- Any code changes.
- `scripts/`, `VERSION`, `CHANGELOG.md`, `cn.json` — fully docs-only cycle.

## Dispatch note

δ: dispatch α with `claude -p` against the cnos repo on branch `cycle/345`. α loads its SKILL.md first and reads the full issue via `gh issue view 345 --json title,body,state,comments`. The branch already exists on origin — α does not create it.
