# Cross-Repo Iteration Lineage: cph → cnos (coherence-drift-sweep-followup, ε)

## Source

- **Repo:** `usurobor/cph`
- **Branch:** `claude/review-repo-coherence-PNbjQ`
- **Wave:** `coherence-drift-sweep-followup-2026-05-18`
- **Wave artifacts (source side):**
  - `cph:.cdd/waves/coherence-drift-sweep-followup-2026-05-18/wave-closeout.md`
  - `cph:.cdd/waves/coherence-drift-sweep-followup-2026-05-18/receipt.md`
  - `cph:.cdd/unreleased/{22,23,24,25}/` (sub artifacts)
- **Master issue:** cph#21 (with precursor cph#16)
- **Sub issues:** cph#22, cph#23, cph#24, cph#25
- **Wave verdict:** all four subs APPROVE R1; 0 fix-rounds; 5 findings closed (F7–F11)

## ε actor

- **Played by:** cnos γ in session `claude/file-cnos-cdr-issue-fi9Ld` (this branch)
- **Authority basis:** operator delegation ("Run epsilon on cph cdd"); `epsilon/SKILL.md §2` ("ε and δ are often the same actor"); `ROLES.md §4` (role-collapse rules)
- **ε output date:** 2026-05-19

## Target patches

| Patch | Destination repo | Path | Status | Finding |
|---|---|---|---|---|
| Rule 6 — anchor oracle evidence on code | `usurobor/cnos` | `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` | landed on cnos branch `claude/file-cnos-cdr-issue-fi9Ld`; awaiting merge to cnos main | F1 + F2 (consolidated) |
| CDR §"Changelog rule" tightening | `usurobor/cph` | `CDR.md` §"Changelog rule" | feedback patch emitted (`cdr-changelog-rule.patch` in this bundle); awaiting cph γ application | F3 |
| Close cph#21 + cph#16 | `usurobor/cph` | GitHub issue state | recommendation; cph γ effects close | F4 (next-MCA) |
| Merge `claude/review-repo-coherence-PNbjQ` → main | `usurobor/cph` | branch state | operator-gated; ε surfaces only | F5 (no-patch) |

## Bundle contents

- `cdd-iteration.md` — ε's canonical output for the cph wave; named `cph:.cdd/waves/<wave>/cdd-iteration.md` once cph γ commits it; emitted here as cross-repo trace
- `cdr-changelog-rule.patch` — feedback patch for cph γ to apply (F3 / cph#21 AC6 resolution)
- `LINEAGE.md` — this file (bilateral trace)

## Bilateral trace direction

This is a **cph→cnos outbound iteration trace** authored from the cnos side because:

1. The ε actor for cph's wave was the cnos γ in this session (hat collapse).
2. The cnos session cannot write to cph; the cph-side outbound trace at `cph:.cdd/iterations/cross-repo/cnos/cph-21-followup-beta-anchor/` is cph γ's to author after applying the feedback patch.
3. The cnos-side mirror at this path serves both as the cnos record of an inbound iteration finding (F1+F2 landed in cnos) and as the staging surface for the cph-side deliverables (cdd-iteration.md content + CDR patch).

Per `cdd/post-release/SKILL.md` §5.6b, cph γ may delete the cph-side outbound bundle once cnos main carries the F1+F2 patch and cph#21 closes — the lineage is preserved in cph's own `cdd-iteration.md`.

## Protocol observations forwarded to cnos#377

This iteration surfaced two additional cross-repo protocol gaps for cnos#377 codification:

1. **ε output direction.** ε's `cdd-iteration.md` for a cph cycle whose patches land in cnos has no canonical bundle home in the protocol. This bundle's dual-purpose shape (iteration content + patch trace) is ad hoc.
2. **Hat-collapse attribution form.** An actor playing γ in repo A and ε for repo B in the same session has no protocol-defined attribution. The hat-collapse acknowledgment block at the top of `cdd-iteration.md` is an ad hoc convention.

Both observations are recorded in `cdd-iteration.md` §7 and forwarded to cnos#377's design phase.

Filed by ε on 2026-05-19.
