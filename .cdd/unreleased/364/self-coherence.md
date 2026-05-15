---
cycle: 364
issue: "https://github.com/usurobor/cnos/issues/364"
date: "2026-05-15"
mode: docs-only
manifest:
  planned:
    - Gap
    - Skills
    - ACs
    - Self-check
    - Debt
    - CDD-Trace
    - Review-readiness
  completed: []
---

# α Self-Coherence — #364

α writes this file section-by-section per `alpha/SKILL.md §2.5` (incremental write discipline). γ scaffold-time content below seeds Gap/Skills/ACs framing; α fills the evidence and trace as work lands.

## §Gap

**Issue:** #364 — Articulate CDD coherence-cell refactor doctrine. Mode: docs-only. Work shape: substantial CDD cycle.

**What exists:** CDD has a canonical lifecycle spec in `src/packages/cnos.cdd/skills/cdd/CDD.md`, a generic α/β/γ/δ/ε ladder in `ROLES.md`, role-local skills (γ, operator, alpha, beta, epsilon, activation), and an artifact-presence checker in `cn-cdd-verify`.

**What is expected:** A draft refactor doctrine document `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` that captures the recursive coherence-cell model: `contract → α/β/γ cell → receipt → validator → δ boundary → accepted receipt as next-scope matter`.

**Where they diverge:** The coherence-cell model is implicit in design discussion and scattered across role / operator / verifier / foundational-coherence surfaces. Today's `operator/SKILL.md` fuses δ boundary policy with dispatch, polling, git identity, release execution, CI recovery, branch cleanup, timeout recovery — membrane/substrate fusion. Today's `gamma/SKILL.md` carries cell-closure coordination plus runtime supervision idioms. The doctrine doc names and falsifies those fusions before implementation refactors begin.

**Peer enumeration (per `gamma/SKILL.md §2.2a`, applied as α as well):**

```bash
test ! -f src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md   # passes — file absent
rg -l "COHERENCE-CELL|coherence cell" src/packages/cnos.cdd/skills/  # no matches
```

The gap is real and additive; no existing surface partially closes it.

## §Skills

**Tier 1 (always):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical CDD lifecycle and role contract (remains authoritative; this doc does not supersede it)
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface

**Tier 2 (always-applicable engineering):**
- markdown authoring (eng default)
- `src/packages/cnos.core/skills/write/SKILL.md` — written artifact authoring

**Tier 3 (issue-specific):**
- `src/packages/cnos.core/skills/design/SKILL.md` — keep policy above detail; avoid surface smearing; distinguish role doctrine from runtime substrate; avoid premature canonicalization; one source of truth per fact; degraded paths must be visible (`§3.9` informs override-as-degraded-boundary-action language).

**β role skill is NOT loaded by α** (per `alpha/SKILL.md` load order).

## §ACs

α will populate AC-by-AC evidence after the COHERENCE-CELL.md draft lands. Each AC maps to the rg/test commands defined in #364, run against the head of the cycle branch. Empty until α drafts the artifact.

| AC | Status | Evidence ref |
|----|--------|--------------|
| AC1 | pending | — |
| AC2 | pending | — |
| AC3 | pending | — |
| AC4 | pending | — |
| AC5 | pending | — |
| AC6 | pending | — |
| AC7 | pending | — |
| AC8 | pending | — |
| AC9 | pending | — |
| AC10 | pending | — |
| AC11 | pending | — |
| AC12 | pending | — |
| AC13 | pending | — |
| AC14 | pending | — |
| AC15 | pending | — |
| AC16 | pending | — |
| AC17 | pending | — |

## §Self-check

α self-check fields populated post-implementation:

- Did α push any ambiguity onto β? (pending)
- Every claim in COHERENCE-CELL.md backed by structural prediction or explicit "illustrative / not normative" label? (pending)
- Did the doc resolve any of the five Open Questions? (must be NO)
- Did the doc pin `cnos.cdd.receipt.v1`? (must be NO)
- Did the diff touch any forbidden surface? (must be NO)

## §Debt

(populated by α after implementation)

## §CDD-Trace

CDD canonical artifact order (`CDD.md §5.2`) for this docs-only cycle:

1. **Design artifact** — marked "not required separately" (the doctrine doc IS the design artifact; this cycle's matter is the design).
2. **Coherence contract** — this `self-coherence.md` §Gap.
3. **Plan** — marked "not required separately" (single new doc + README pointer; sequencing is trivial: write doc section-by-section, then update README, then run AC oracle).
4. **Tests** — N/A for docs-only; AC oracles are the structural tests (each AC has rg/grep oracle defined in #364).
5. **Code** — N/A.
6. **Docs** — `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` (primary), `src/packages/cnos.cdd/README.md` (pointer).
7. **Self-coherence** — this file, populated incrementally.
8. **Pre-review gate** — verified before review-readiness signal.

(α fills the artifact enumeration list against `git diff --stat origin/main..HEAD` at pre-review-gate time per `alpha/SKILL.md §2.6 row 11`.)

## §Review-readiness

(α appends this section per round after pre-review-gate passes — see `alpha/SKILL.md §2.7`.)
