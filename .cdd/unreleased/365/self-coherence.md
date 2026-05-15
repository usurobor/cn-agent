---
manifest:
  sections:
    - Gap
    - Skills
    - ACs
    - Self-check
    - Debt
    - CDD-Trace
    - Review-readiness
  completed: [Gap, Skills]
---

# Self-Coherence — #365 (I6 CDD artifact validator era-mismatch on v3.75.0/v3.76.0)

## Gap

`src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` enforces current
artifact requirements on cycles released before those requirements were in
force. The CI job `cdd-artifact-check` (`I6 CDD artifact ledger validation`)
fails on every push to every branch with **87 errors / 426 total checks**
(reproduced locally on `cycle/365@4fa0bd05`; the issue body cited 77 — drift
is because more v3.75/v3.76 cycles have shipped since the issue was filed).
Every one of the 87 errors is on a `v3.75.0` (126) or `v3.76.0` (2) cycle:

```
src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify --all 2>&1 \
  | grep "❌" | grep -oE "v3\.[0-9]+\.[0-9]+" | sort | uniq -c
    126 v3.75.0
      2 v3.76.0
```

The validator already has an `is_legacy_version` helper that downgrades
errors to warnings for older releases. The current cutoff is `minor < 74`
(i.e. cycles pre-`v3.74.0` warn; `v3.74.0`+ strict), but the file header
comment at line 253 claims the cutoff is `v3.65.0`. Both the threshold and
the comment need to move forward. The actual era boundary is `v3.77.0`:
`v3.74.0` already has every required artifact (cycle 327 ships clean),
`v3.75.0` and `v3.76.0` predate the artifact lockdown, and `v3.77.0` is
the next un-shipped minor where strict enforcement becomes correct.

**Mode (selected):** cutoff bump rather than per-artifact era policy. γ's
peer enumeration confirms the failure shape is sharp at a single boundary —
α-closeout/γ-closeout absence collapses to one threshold (`v3.77.0`). A
per-artifact era table would add three knobs all set to the same value;
that is complexity without payoff. The single-function change keeps the
validator simple and matches γ's "lightweight path" framing.

## Skills

| Tier | Skill | Use |
|------|-------|-----|
| 1 | `cdd/CDD.md` | Lifecycle contract; artifact-location matrix §5.3a |
| 1 | `cdd/alpha/SKILL.md` | Role contract; incremental self-coherence §2.5 |
| 2 | `cnos.eng/skills/eng/*` | Tier 2 always-applicable engineering bundle |
| 3 | `cnos.core/skills/write/SKILL.md` | Doc + commit-message style (CHANGELOG row, self-coherence prose) |
| 3 | `cnos.eng/skills/eng/tool/SKILL.md` | Bash standards — `set -euo pipefail`, idempotence, exit codes |
| 3 | `cnos.eng/skills/eng/test/SKILL.md` | Test invariant first; oracle = exit code + warn/fail line for each era boundary |
| 3 | `cnos.eng/skills/eng/ship/SKILL.md` | Bug-fix flow: failing test → fix → passing test → no regressions |

No `eng/{language}` skill is loaded — the validator is a single bash file
and γ confirmed no other language surface is involved. The diff is
bash + bash test harness + Markdown (CHANGELOG row + this file).
