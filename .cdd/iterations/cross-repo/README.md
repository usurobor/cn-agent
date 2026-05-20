# Cross-Repo Trace Bundles

Each cross-repo coordination unit lives at:

  .cdd/iterations/cross-repo/{counterpart-repo}/{slug}/

Where:
- {counterpart-repo} — the repo that is *not* cnos (e.g., tsc, cph, cn-sigma, cn-rho)
- {slug}            — a short, lowercase, hyphen-separated descriptor of the coordination
                      (e.g., cdd-supercycle, bootstrap-cdr, agent-activate-skill)

This path is **direction-agnostic** — both inbound proposals (cnos as target) and outbound iteration traces (cnos as source) live under the same path shape; the directional case is named in the bundle's `LINEAGE.md`.

## Canonical surface

Protocol doctrine for cross-repo coordination — directional cases, canonical path,
STATUS state machine, bundle file set, LINEAGE.md schema, feedback-patch format,
bundle archival rule, hat-collapse attribution — lives in
`src/packages/cnos.cdd/skills/cdd/cross-repo/SKILL.md`. This README describes the
bundle convention only; the protocol governs.

## When to create a bundle

Create a bundle when one of the four directional cases (per `cdd/cross-repo/SKILL.md §2.1`) applies:

- **(a) inbound proposal** — counterpart repo emits a proposal; cnos accepts/modifies/rejects.
- **(b) outbound iteration trace** — cnos cycle produces patches landing in another repo.
- **(c) bilateral iteration** — a single session produces patches for both repos via hat-collapse.
- **(d) operator-pending bundle** — content drafted for a counterpart the current cnos session cannot reach.

## Bundle file set

The required files per directional case are defined in `cdd/cross-repo/SKILL.md §2.5`. Each
bundle's `LINEAGE.md` names the case in its first paragraph; the file set follows from the case.

## Bundle-state phases

Bundles carry one of three phases, mapped to the canonical STATUS state machine
(per `cdd/cross-repo/SKILL.md §2.4`):

| Phase | STATUS event (cases a, c with proposal-shape) | Derivation for cases b, c, d |
|---|---|---|
| `open` | `drafted` or `submitted` | Patches in flight (b); ε deliverables pending (c); operator-pending drafts with `Disposition: drafted (operator-pending)` (d) |
| `converging` | `accepted` or `modified` (no terminal `landed` yet; may carry partial per-sub `landed` rows) | (n/a for b/c/d) |
| `closed` | terminal `landed` (1:1) or master-close `landed` (master/sub) or `rejected` or `withdrawn` | All patches confirmed landed (b); both sides land + iteration recorded (c); operator effects application (d) |

Audit events (`revised`, `corrected`) do not change phase — they record post-hoc corrections without re-opening or re-closing the bundle.

A bundle is closed when its terminal phase condition is met. Case (d) bundles
typically archive on `closed`.

## Existing bundles

- `cn-rho/` — operator-pending bootstrap drafts for the planned `usurobor/cn-rho` research persona hub (case d.1)
- `cn-sigma/agent-activate-skill/` — inbound proposal, 1:1 (case a.1) → cnos#379
- `cn-sigma/discipline-section-2026-05-19/` — operator-pending feedback patch for cn-sigma persona file (case d.2)
- `cph/bootstrap-cdr/` — inbound proposal, master+sub (case a.2) → cnos#376 + subs (post-rename canonical; pre-rename mirror at `gait-support-paths/bootstrap-cdr/`)
- `cph/coherence-drift-sweep-followup-2026-05-18/` — bilateral iteration with hat-collapse (case c)
- `cph/issue-32-tightening-2026-05-19/` — operator-pending proposal as GitHub issue comment (case d.3)
- `gait-support-paths/bootstrap-cdr/` — pre-rename historical mirror of `cph/bootstrap-cdr/` (preserved as audit artifact)
- `tsc/cdd-supercycle/` — outbound iteration trace, 6 patches (case b.1) → cnos#331/#332
