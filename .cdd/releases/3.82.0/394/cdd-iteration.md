# cdd-iteration — Cycle #394 (Sub 2 of #376)

**Cycle:** #394 — Sub 2 of #376: bootstrap cnos.cdr package skeleton (cn.package.json + README.md + SKILL.md)
**Merge:** [filled by γ post-merge with merge SHA]
**Closed:** 2026-05-21 (auto-closed via `Closes #394` in merge commit)
**Mode:** design-and-build (γ+α+β-collapsed-on-δ; β-α-collapse acknowledged in `beta-review.md`)
**Rounds:** R1 APPROVE (no fix-round)
**ACs:** 6/6 PASS

## §1 Findings dispositioned

### F1: `cn build --check` does not enforce skill `calls:` path existence

- **Source:** cycle 394 β-review F1 + β-closeout F1.
- **Class:** observational (forward-looking; not a current-cycle gap).
- **Trigger:** β verification during AC4 re-check.
- **Description:** The package-validation oracle (`cn build --check`) validates manifest schema only; it does not walk the skill-loader `calls:` graph to verify referenced files exist. cnos.cdr's loader names `alpha/SKILL.md` ... `epsilon/SKILL.md` in `calls:` as forthcoming Sub 3 files; these paths do not yet exist on disk; `cn build --check` still passes. This is the current design — manifest validation is independent of skill-graph validation.
- **Root cause:** Discipline boundary by design — the package layer (manifest) and the skill layer (loader graph) are validated by separate mechanisms. There is no stricter skill-loader walk currently implemented.
- **Disposition:** `no-patch` (forward-looking observation). If a stricter skill-loader walk is later introduced (e.g. a `cn check-skills` command that validates `calls:` paths), the cnos.cdr loader will need either (a) the Sub 3 files to land before strict-check is enabled, or (b) a "forthcoming" marker in the loader frontmatter that the strict-check honors. Decision deferred until the strict-check is on someone's roadmap.
- **Issue filed:** none. Carry as informational debt.
- **First AC for the eventual MCA:** if a stricter skill-loader walk lands, this cycle's loader (or Sub 3's per-role files) is the first surface that demonstrates compliance.

### F2: cnos.cdr v0.1 skeleton is loadable but not yet dispatchable

- **Source:** cycle 394 α-closeout F3.
- **Class:** observational (by-design; not a gap).
- **Trigger:** α self-coherence during AC4 verification.
- **Description:** The package is discoverable (`✓ cnos.cdr: valid` from `cn build --check`), but the role-overlay layer (Sub 3) is not yet authored. A dispatch matching `cdr` triggers will load CDR.md as the canonical contract but will not find a per-role overlay until Sub 3 ships. The SKILL.md loader body explicitly names this as "Role overlays (forthcoming — Sub 3 of cnos#376)".
- **Root cause:** Sub 2's scope is the package skeleton; role overlays are Sub 3's scope. The "loadable but not dispatchable" state is the intended v0.1 skeleton state.
- **Disposition:** `no-patch` (by-design observation). Sub 3 dispatch will close this. Future dispatch-error messages (when the role-overlay layer is engaged) should distinguish "role overlay missing because Sub 3 unshipped" from "role overlay missing due to misconfiguration" — captured as Sub 3 dispatch context.
- **Issue filed:** none. Sub 3 of #376 (already filed under that issue's planning).
- **First AC for the eventual MCA:** Sub 3's first dispatch — when `alpha/SKILL.md` lands, the dispatch loop can engage cnos.cdr end-to-end.

## §2 No-findings observations (informational)

- **AC1 schema-parity outcome:** cnos.cdr's manifest is a proper subset of cnos.cdd's (omitting the optional `commands` field). This matches cnos.eng and cnos.cdd.kata which also omit `commands`. The schema's `commands` field is omitempty per `pkg.go FullPackageManifest.Commands` JSON tag default. No finding.
- **Engines range vs pin:** cnos.cdr uses `>=3.81.0` (range) per cnos.kata + cnos.cdd.kata convention for new packages; cnos.cdd/core/eng pin exactly to `3.81.0` because they are aligned with the kernel release-train. No central policy document records this convention; observable from the package set only. Not raised as a finding because it is consistent with prior packages and the design-notes record the rationale. Recorded for awareness in cycle 394 α-closeout debt item D1.
- **Five-role grammar vs cnos.cdd's three-role loader:** cnos.cdr ships with α/β/γ/δ/ε from day one (per `ROLES.md §1` formalised five-role ladder + CDR.md Field 5); cnos.cdd's loader pre-dates ε formalization and lists only α/β/γ + lifecycle sub-skills. cnos.cdr does not inherit cdd's pre-five-role gap. This is intentional; recorded in design-notes §3.D1. No finding; positive design signal.
- **β-α collapse reconciled per class:** the cycle ran with α+β collapsed. CDR.md Field 6 prohibits α=β for research-class. This cycle is engineering-class docs-and-metadata (CDS-class under repairable feedback), not research-class claim transmission. Same reconciliation pattern as cycle 390's beta-review borderline. The discipline boundary is held; no finding.

## §3 Trigger assessment (per `gamma/SKILL.md §2.8` table)

| Trigger | Fire condition | Fired? | ε note |
|---|---|---|---|
| Review churn | review rounds > 2 | **No** | R1 APPROVE on first pass. |
| Mechanical overload | mechanical ratio > 20% AND findings ≥ 10 | **No** | Mechanical ratio is 100% by design (docs-and-metadata cycle, mechanical oracles), but findings count is 2; ≥ 10 not met. |
| Avoidable tooling / environment failure | environment blocked the cycle | **No** | No tooling friction. `go build` and `jq` worked first try. |
| Loaded-skill miss | a loaded skill should have prevented a finding | **No** | No findings raised (F1, F2 are observational; F1 is forward-looking, F2 is by-design). |

No trigger fires. Cycle ran cleanly.

## §4 INDEX update

Add to `.cdd/iterations/INDEX.md`:

```
| 394 | #394 | 2026-05-21 | 2 | 0 | 0 | 2 | .cdd/unreleased/394/cdd-iteration.md |
```

Findings: 2 (F1, F2). Patches: 0 immediate. MCAs: 0 (both findings are observational, not protocol-gap-shaped). No-patch: 2.

## §5 Skill-gap candidate disposition

Neither F1 nor F2 rises to `cdd-protocol-gap` class. F1 is a forward-looking observation about a feature (skill-loader strict-check) that does not currently exist. F2 is a by-design observation about Sub 2's v0.1 skeleton scope.

No skill-gap candidates raised.

## §6 Deferred outputs

- **Cycle-dir move.** `.cdd/unreleased/394/` → `.cdd/releases/<version>/394/` at next release per `release/SKILL.md §2.5a`. Not blocking; standard release-time mechanic.
- **cnos#376 close-out comment** — posted post-merge by γ. Comment shape recorded in `gamma-closeout.md §3`.

## §7 Next-MCA commitment

Sub 3 (cnos.cdr role overlays) and Sub 4 (cnos.cdr empirical-anchor doc) can dispatch in parallel against cnos.cdr as a stable package target. The v0.1 skeleton is the dispatch base; neither sub needs to re-author the manifest or the loader.

Filed by ε on 2026-05-21.
