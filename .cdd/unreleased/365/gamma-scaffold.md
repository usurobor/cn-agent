---
cycle: 365
issue: "https://github.com/usurobor/cnos/issues/365"
date: "2026-05-15"
dispatch_configuration: "§5.1 multi-session via codex"
base_sha: "16aaef89"
scope: "Make cdd-verify era-appropriate for v3.75.0/v3.76.0 cycles missing alpha-closeout.md, gamma-closeout.md, and CDD Trace sections."
mode: "design-and-build"
tier_3_skills:
  - src/packages/cnos.core/skills/write/SKILL.md
  - src/packages/cnos.eng/skills/eng/tool/SKILL.md
  - src/packages/cnos.eng/skills/eng/test/SKILL.md
  - src/packages/cnos.eng/skills/eng/ship/SKILL.md
---

# γ Scaffold — #365

## Gap

`src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` enforces current
artifact requirements on cycles released before those requirements were in
force. CI job `cdd-artifact-check` (`I6 CDD artifact ledger validation`) fails
with 77 errors across the v3.75.0 and v3.76.0 release directories on every
push, blocking branch CI.

Two failure patterns:

1. **Missing `CDD Trace` section** in `self-coherence.md` for some v3.75.0
   cycles.
2. **Missing `alpha-closeout.md` / `gamma-closeout.md`** for every v3.75.0
   cycle and one v3.76.0 cycle (cycle 327 of v3.74.0 already has all
   artifacts).

The validator already has an `is_legacy_version` function (lines 251–274 of
`cn-cdd-verify`) that downgrades errors to warnings for "legacy" releases.
The current cutoff is `minor < 74` — i.e. cycles released before v3.74.0.
v3.75.0 and v3.76.0 are *not* legacy under that cutoff and therefore fail
strict enforcement on artifacts that were not yet required when they
shipped.

Source comment vs. code mismatch: the function header comment at line 253
says "Returns 0 (true) if version is before v3.65.0 (epoch threshold)"; the
code at lines 265–272 actually compares against `3.74`. This drift is a
secondary finding for the cycle to resolve.

## Peer enumeration

Before scaffolding, γ ran:

- `ls .cdd/releases/3.75.0/` → 28 cycle directories.
- `ls .cdd/releases/3.75.0/*/alpha-closeout.md 2>/dev/null | wc -l` → `0`
  (confirms every v3.75.0 cycle is missing the artifact).
- `grep -L "CDD Trace" .cdd/releases/3.75.0/*/self-coherence.md` →
  `323/`, `340/`, `341/`, `348/`, `361/` are missing the section.
- `ls .cdd/releases/3.76.0/362/` → has all current artifacts including
  `gamma-scaffold.md`. The 3.76.0 run-failure must therefore be on a
  *different* v3.76.0 cycle directory — re-verified with `gh run view
  25905265925` evidence in the issue body (lines reference 348 and 292,
  both v3.75.0). The single failing v3.76.0 entry is consistent with a
  prior CI run that included an incomplete unreleased directory.
- `ls .cdd/releases/3.74.0/` → cycle 327 only, has every artifact.
  v3.74.0 is the current cutoff *boundary* and it satisfies strict checks
  on its own — confirms the cutoff bump only needs to extend forward, not
  retroactively whitelist v3.74.0.
- `grep -rn "is_legacy\|3\.65\|3\.74\|3\.77" cn-cdd-verify` → cutoff is
  in exactly one place (`is_legacy_version`) and the comment-vs-code drift
  is in exactly one place (line 253 comment vs. lines 265–272 code).
- `grep -rn "cdd-verify" .github/workflows/` → invoked once from
  `.github/workflows/build.yml` line 269 as `cn-cdd-verify --all` plus
  `test-fixtures.sh` line 272.

No `eng/{language}` skill is needed — the validator is bash. The fix
surface is one file (the validator script) plus the existing test
fixtures.

## Mode

`design-and-build` — the design call is small but real: extend the legacy
cutoff to cover v3.75.0–v3.76.0 vs. introduce per-artifact era policies.
The lightweight path (cutoff bump + comment repair + a fixture that pins
the new boundary) closes the gap and keeps the validator simple. α makes
the final call between options in §self-coherence Mode.

## Cycle scope sizing

| Factor | Reading | Splitting signal? |
|---|---|---|
| (a) New code surface | Zero new modules; one bash function changes | No |
| (b) Cross-module breadth | One file (`cn-cdd-verify`) plus optional test fixture additions | No |
| (c) Lifecycle span | code + test fixture + brief release note (CHANGELOG row stays for δ at release) | No |
| (d) MCA preconditions | Gap clear; design space narrow (cutoff vs. era table); no prior design doc cited | One factor (no stable design path) → `design-and-build`, not MCA |
| (e) Independent shippability | Monolithic — fixing CI requires the cutoff change *and* the fixtures to land together | No |

Decision: **keep whole**. 4–6 ACs; small diff.

## Acceptance criteria (γ-drafted; α refines)

These are the ACs the issue should test. α takes them as the starting
contract and refines wording in `self-coherence.md`.

1. **AC1 — CI green.** `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify --all`
   exits 0 on `origin/main` after the change, and the
   `cdd-artifact-check` job in `.github/workflows/build.yml` passes on
   `cycle/365`.
2. **AC2 — Era-appropriate handling for v3.75.0 / v3.76.0.** For every
   cycle directory under `.cdd/releases/3.75.0/` and `.cdd/releases/3.76.0/`
   that is missing `alpha-closeout.md`, `gamma-closeout.md`, or the
   `CDD Trace` section of `self-coherence.md`, the validator emits a
   **warning** rather than an **error**.
3. **AC3 — Strict enforcement preserved for current cycles.** New
   in-flight cycles under `.cdd/unreleased/{N}/` and cycles in releases
   `≥ v3.77.0` still get the strict-error path when required artifacts or
   sections are missing. Cycle 365 itself must satisfy strict enforcement
   when it ships.
4. **AC4 — Comment / code coherence.** The header comment on
   `is_legacy_version` accurately describes the actual cutoff used by the
   code. Whichever direction is chosen (raise the cutoff, or rewrite the
   helper as per-artifact era policy), the comment matches the
   implementation and names the rationale in one sentence.
5. **AC5 — Test fixture pinning the boundary.** A fixture or test case
   under `src/packages/cnos.cdd/commands/cdd-verify/test/` exercises the
   v3.75.0/v3.76.0 era-appropriate path so a future regression that
   pushes the cutoff back surfaces in `test-fixtures.sh`.
6. **AC6 — RELEASE.md / CHANGELOG note prepared.** A line in the
   `Unreleased` section of `CHANGELOG.md` (or equivalent surface used by
   the project's release prep) names the validator change so δ does not
   have to retro-author release notes at tag time. Doc-only — δ still
   owns the tag.

## Non-goals

- Do not retroactively backfill `alpha-closeout.md` or `gamma-closeout.md`
  for v3.75.0 / v3.76.0 cycles. They shipped under looser rules; the
  cycle's job is to make the validator era-appropriate, not to invent
  retroactive close-outs.
- Do not change which artifacts are required for the *current* era. AC3
  pins this explicitly.
- Do not rewrite the exception-file (`exceptions.example.yml`) flow.
  Exception backing is a separate surface; this cycle's fix is on the
  classification function, not on per-path exceptions.

## Dispatch state

- Branch: `cycle/365` from `origin/main@16aaef89`.
- Configuration: §5.1 multi-session via codex; δ dispatches α and β
  sequentially with the prompts in this cycle directory.
- Polling: γ exits after this scaffold lands; δ owns dispatch sequencing
  per `operator/SKILL.md` §1.2.

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | issue body, run 25905265925 evidence, .cdd/releases/3.75.0/, validator script | cdd | I6 fails 77/418 on every push; cycles 273…361 + 362 missing era-current artifacts |
| 1 Select | issue #365 | cdd | Selected under CDD.md §3.2 (operational infrastructure override — CI block on all branches) |
| 2 Branch | `cycle/365` | cdd | Created from `origin/main@16aaef89`, pre-flight passed |
| 3 Bootstrap | `.cdd/unreleased/365/` | cdd | Cycle directory + this scaffold (no version-snapshot stubs — validator-only cycle) |
| 4 Gap | this file §Gap | cdd | Era-mismatched validator; cutoff is `< 3.74` but current shipped cycles are 3.75/3.76 |
| 5 Mode | this file §Mode | cdd | `design-and-build`; α picks between cutoff bump and per-artifact era policy |
