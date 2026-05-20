# α close-out — cycle #375

Cycle: usurobor/cnos#375 — γ-side pre-dispatch gate for `gamma-scaffold.md` (rule 3.11b symmetry).

Branch: `origin/cycle/375` at `c4d29344` (pre-merge SHA).

## Summary

Single section-level addition to `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` §2.5 Step 3b: new sub-section "Pre-dispatch γ scaffold check (binding gate)" naming `.cdd/unreleased/{N}/gamma-scaffold.md` as a binding pre-dispatch existence check, framing it as the dual of `review/SKILL.md` rule 3.11b, and citing cycle #369 as the empirical anchor. 26 insertions, 0 deletions, no surrounding restructure.

## What I touched

- `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` — added §2.5 Step 3b sub-section (binding rule + scaffold contents + dual framing + empirical anchor + bash procedure + exemption deference). Insertion point chosen at the top of Step 3b because the gate binds to α dispatch, not to branch creation.

No other file touched at α implementation time (γ scaffold pre-existed by handoff; β-review and close-outs are downstream artifacts).

## What I did not touch (per non-goals)

- `review/SKILL.md` rule 3.11b — unchanged. The γ-side gate cross-references the β-side rule but does not edit it.
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — unchanged. The issue body permitted either `gamma/SKILL.md` §2.5 or `CDD.md` §1.4 step 3; γ chose the former as the more specific surface.
- Any validator (`cn-cdd-verify`, scripts) — unchanged. Per non-goal "Do not introduce a new validator surface; the check is a SKILL-level gate."
- Any CI workflow — unchanged.
- Any other role skill (`alpha/`, `beta/`, `operator/`, `release/`, `post-release/`) — unchanged.

## Self-coherence

See `.cdd/unreleased/375/self-coherence.md` for AC1–AC4 mapping and CDD Trace. All four ACs PASS at R1; no fix rounds needed.

## Notes for next α

- The bash procedure in the new sub-section uses `git ls-tree -r --name-only origin/cycle/<N> .cdd/unreleased/<N>/gamma-scaffold.md | grep -q .`. This pattern parallels the existing pre-flight check in §2.5 Step 3a ("no stalled `.cdd/unreleased/{N}/` on main") so γ has a consistent idiom for filesystem-on-remote-branch checks.
- The exemption-channel deference (the sub-section's last paragraph) intentionally does *not* duplicate rule 3.11b's exemption-discoverability prose; it points to that prose. If rule 3.11b's exemption surface is ever moved or restated, this γ-side gate inherits the change automatically by reference.
