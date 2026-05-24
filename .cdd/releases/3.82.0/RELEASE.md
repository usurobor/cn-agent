# 3.82.0

## Outcome

Coherence delta: CDS, CDR, and cnos.handoff reach v0.1 alongside the compact CCNF kernel. v3.82.0 pins this four-package shape as the baseline.

`VERSION` bumps `3.81.0` → `3.82.0`. CDS and CDR READMEs replace "v0.1 skeleton" / "in flight" framing with current v0.1-complete status. CCNF kernel (`CDD.md`) is byte-identical to 3.81.0.

## Why it matters

CDS, CDR, and cnos.handoff shipped across the 3.81 → 3.82 window. Their READMEs kept saying "skeleton" / "Pending Sub N" / "in flight" while the underlying packages were already v0.1-complete and load-bearing. A reader evaluating whether to depend on `cnos.cds` would conclude "no, skeleton" when the reality is "yes, ready." This release fixes the discoverability gap and gives the architecture a citable version.

## Added

- **CDS v0.1** ([#403](https://github.com/usurobor/cnos/issues/403); subs #406–#412): `cnos.cds` package with the six-field software-realization contract (CDS.md, 1043 lines), operational sub-area overlays (`selection/`, `lifecycle/`, `tracking/`, `artifacts/`, `review/`, `gate/`, `assessment/`), extraction map, and empirical-anchor doc. CDD.md "pending cds extraction" markers swept.
- **cnos.handoff v0.1** ([#404](https://github.com/usurobor/cnos/issues/404); subs #415–#420): wire-format doctrine for inter-agent / inter-activation / inter-repo handoff. Five sub-skills under `skills/handoff/`: `cross-repo/`, `dispatch/`, `mid-flight/`, `artifact-channel/`, `receipt-stream/`.
- **CDR v0.1** ([#376](https://github.com/usurobor/cnos/issues/376); subs #390 / #394 / #395 / #396): research-realization peer. CDR.md + loader + five role overlays + empirical-anchor-cph.md.
- **Three design essays** under `docs/gamma/essays/`: `DECREASING-INCOHERENCE.md` ([#414](https://github.com/usurobor/cnos/issues/414); receipt → measure → bottleneck → next contract); `CELL-OF-CELLS.md` ([#424](https://github.com/usurobor/cnos/issues/424); recursive coherence-cell system model); `BOX-AND-THE-RUNNER.md` ([#425](https://github.com/usurobor/cnos/issues/425); remote-runner-delegation primitive + 6-field receipt).
- **`delta/SKILL.md §8`** ([#425](https://github.com/usurobor/cnos/issues/425)): codifies remote-runner delegation as a δ-class effect surface with a required 6-field receipt.
- **`docs/gamma/design/`** new directory with `ccnf-o-track-a1-survey.md` ([#421](https://github.com/usurobor/cnos/issues/421)): 382-line CCNF-O Track A1 survey pinning five decisions.
- **Sigma activation bundle** ([#413](https://github.com/usurobor/cnos/issues/413)): staged at `.cdd/iterations/cross-repo/cn-sigma/sigma-activation-2026-05-22/`; operator-applied to `cn-sigma` main.

## Changed

- **`cnos.cds/README.md`** ([#422](https://github.com/usurobor/cnos/issues/422)): "v0.1 skeleton" + "Pending Sub 2" → v0.1-complete; sub-issue references for #406–#412.
- **`cnos.cdr/README.md`** ([#422](https://github.com/usurobor/cnos/issues/422)): "Role overlays in flight" → v0.1-complete; sub-issue references for #390 / #394 / #395 / #396.
- **`VERSION`** ([#422](https://github.com/usurobor/cnos/issues/422)): `3.81.0` → `3.82.0`.
- **`cnos.cdd/skills/cdd/cross-repo/SKILL.md`** ([#423](https://github.com/usurobor/cnos/issues/423)): `triggers: [cross-repo-moved, handoff-extracted]` added to the compatibility-pointer frontmatter. Fixes `TestValidate_RealCorpus_NoEmptyTriggers` build failure that landed with #416's stub.

## Unchanged

- **CCNF kernel** (`CDD.md`, `COHERENCE-CELL.md`, `COHERENCE-CELL-NORMAL-FORM.md`): byte-identical to 3.81.0.
- **`cn cdd verify`**: no Go source changes; behavior identical.
- **`scripts/release.sh`** and the release-effector skill: unchanged.

## Known Issues

- **`cnos.cds/skills/cds/SKILL.md`** carries stale "Pending Subs 3–5" wording in the loader frontmatter, anticipating per-role overlays (`alpha/beta/gamma/delta/epsilon`) that did not ship in v0.1; the wave landed operational sub-area overlays instead. `CDS.md` governs in any conflict per the loader's Conflict rule. Cleanup deferred. Recorded in cycle/422 `cdd-iteration.md`.
- **`release.yml`** build job has a git-auth failure on tag-push events ("could not read Username for 'https://github.com'"). 2 of 4 platforms produced binaries before cancel. v3.82.0 binaries were not auto-published; this release ships release-notes only. Recorded as `cdd-tooling-gap` finding in cycle/426 `cdd-iteration.md`; next-MCA fix-cycle pending.
- **`on.push.tags` does not reliably trigger on tag force-update.** Recorded as `cdd-tooling-gap` finding in cycle/426; doctrine update pending.
- **CCNF-O / TSC steering roadmap** ([#405](https://github.com/usurobor/cnos/issues/405)) remains open. Track A1 shipped the survey; Tracks A2–A6 and B1–B6 are deliberately not executed pending field evidence.

## γ note (release-notes hygiene)

This release boundary was published via two remote-runner-delegation workflows (cnos#425, cnos#426), the first explicit uses of the doctrine landed in cnos#425. The initial v3.82.0 release notes did not follow the `cnos.core/skills/write/SKILL.md` discipline (5× longer than 3.80.0/3.81.0, mixed jobs, repeated stop-condition language). Rewritten under the write skill and published via cnos#427; the discoverability gap that motivated v3.82.0 is now closed at the release-notes layer too.
