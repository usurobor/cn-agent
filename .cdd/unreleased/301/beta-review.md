# β review — #301

## Round 1

**Verdict:** REQUEST CHANGES

**Branch:** `claude/cnos-alpha-tier3-skills-MuE2P`
**Base:** `a8e67b7` (origin/main)
**Head reviewed:** `ed5f218`
**Branch CI state:** unverified-from-β-environment — see Notes
**Merge instruction:** Do not merge yet. After α lands the round-1 fixes
on-branch and CI is green on the resulting HEAD, β will re-review and
issue the merge instruction.

## §2.0 Issue Contract

### AC Coverage

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| 1 | Stratified CUE schema (hard-gate / spec-required-exception-backed / optional-with-default / reserved / open) | yes | met | `schemas/skill.cue` L20–71. Hard-gate are unconditional CUE values; exception-backed are `?`; struct is open (`...`) so unknown package-local keys (e.g. `parent:`) pass. Schema field set matches AC1. |
| 2 | Extraction script with `set -euo pipefail`, fail-fast prereqs, deterministic ordering, NO_COLOR, machine-readable diagnostics; surface boundary script-vs-CUE-vs-filesystem | yes | met | `tools/validate-skill-frontmatter.sh`. `set -euo pipefail` (L29), `need cue/jq/awk/find` exit 2 (L40–50), `find … \| sort` (L280), NO_COLOR honored AND non-TTY (L32–36), per-line `path :: field :: rule :: reason :: fix` (L116–117). Surface boundary documented at L1–9 and held in code. β confirmed exit-2 path locally with `PATH=/usr/bin:/bin` (rc=2). |
| 3 | CI job `skill-frontmatter-check` (I5) added; pinned action tag + pinned CUE version; notify aggregation updated | yes | met | `.github/workflows/build.yml` L242–260. `cue-lang/setup-cue@v1.0.1` (L252), CUE `v0.13.2` (L254), runs `--self-test` (L257) and full validation (L260). `notify.needs` includes `skill-frontmatter-check` (L265); aggregation row `I5:${{ needs.skill-frontmatter-check.result }}` (L285). |
| 4 | `schemas/skill-exceptions.json` field-specific exceptions with `reason` + `spec_ref`; hard-gate not exceptable | yes | met | All 43 entries are field-specific (not file-wide), each carries `reason` and `spec_ref`. `SPEC_REQUIRED_EXCEPTION_BACKED` array in script (L92) restricts the exception path to `artifact_class/kata_surface/inputs/outputs`; hard-gate failures surface from `cue vet` regardless. |
| 5 | Exception list is intended to shrink; entries addressable | yes | met | `schemas/README.md` "How to shrink the exception list" gives the mechanical procedure; each entry's `spec_ref` names the LANGUAGE-SPEC section to author from. |
| 6 | Static `calls` validate against package skill root (not caller dir); invalid targets fail; `calls_dynamic` not validated for existence | yes | met | `package_skill_root_of` (L244–272) computes the common ancestor of every SKILL.md under the package's `skills/` subtree, matching the issue example (`cdd/alpha` → `cdd/design/SKILL.md`). β confirmed with the broken-calls fixture. `calls_dynamic` source/constraint shape only (schema L58–62), no existence check. |
| 7 | Enums for `scope`, `visibility`, `artifact_class`, `kata_surface` | yes | met | `schemas/skill.cue` L31, L42, L35, L36. `bad-enum-scope` fixture exercises the failure path. |
| 8 | Negative proof — fixtures or `--self-test`; the four minimum negative cases | yes | met | `schemas/fixtures/skill-frontmatter/{valid,invalid}/`: 3 positive (minimal, full, full/sub), 4 negative (missing-hard-gate, bad-enum-scope, broken-calls, missing-required-no-exception) each with `.expect` sidecar. β confirmed `--self-test` rc=0 locally with CUE `v0.13.2`. |
| 9 | `schemas/README.md` documents what is/isn't validated, exceptions, shrink procedure, local-run, CI integration with pinned versions, schema-update procedure | yes | met | `schemas/README.md` covers every required topic; cross-refs LANGUAGE-SPEC v0.1, v0.2 draft, #289 (v0.2 stabilization), #303 (`ctb-check` v0). |

### Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `schemas/skill.cue` | yes | added | new authority surface for machine-checkable shape |
| `tools/validate-skill-frontmatter.sh` | yes | added | extraction + script-side checks |
| `schemas/skill-exceptions.json` | yes | added | debt ledger (43 entries) |
| `schemas/README.md` | yes | added | operator-facing summary; covers AC9 in full |
| `schemas/fixtures/skill-frontmatter/{valid,invalid}/` | yes | added | AC8 positive/negative coverage |
| `.github/workflows/build.yml` | yes | updated | I5 job + notify aggregation |
| 46 `SKILL.md` frontmatter backfills under `src/packages/**` | yes | updated | hard-gate fields backfilled where missing; YAML-quoted descriptions/inputs in `cdd/design/SKILL.md` and `cdd/gamma/SKILL.md` for CUE strictness; `cdd/release/SKILL.md` `calls: -writing` → `calls: []` |
| `docs/alpha/ctb/LANGUAGE-SPEC.md` v0.1 | n/a | unchanged | normative source — left untouched, as scope demands |
| `docs/alpha/ctb/LANGUAGE-SPEC-v0.2-draft.md` | n/a | unchanged | non-goal, left untouched |

### CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `.cdd/unreleased/301/self-coherence.md` (gap, mode, AC mapping, CDD Trace, peer enumeration, harness audit, polyglot re-audit, known debt, review-readiness signal) | yes (substantial cycle) | yes | 197 lines on branch HEAD; review-readiness section present at L187. Three artifact-internal numeric/SHA inconsistencies named in Findings F2 and F3. |
| Tests / fixtures | yes (substantial cycle) | yes | `--self-test` integration with 4 negative + 3 positive fixtures; full-corpus validation acts as integration test on production data. |
| Branch artifacts on cycle branch (not PR) | yes | yes | per CDD §Tracking dyad-on-one-branch model. β review committed to same branch. |

## Architecture Check

| Check | Result | Notes |
|---|---|---|
| Reason to change preserved | yes | schema (shape), script (mechanism), workflow (CI gate), exception list (debt ledger), README (doc) — four-and-a-bit surfaces, each with one reason to change |
| Policy above detail preserved | yes | CUE owns shape/type/enum policy; shell owns discovery/extraction/exceptions/calls-fs mechanism; the script does not encode shape rules and the schema does not encode discovery |
| Interfaces remain truthful | yes | exit codes 0/1/2 are distinct; per-finding diagnostics are uniform (`path :: field :: rule :: reason :: fix`); `--self-test`, `--root`, `--file`, `--help` all behave as documented; β verified each mode locally |
| Registry model remains unified | yes | one CUE schema is the single source of truth for the shape; the README's stratification table is documentation derived from it; the exception list is the single inverse ledger |
| Source / artifact / installed boundary preserved | yes | source = SKILL.md frontmatter; the validator enforces source shape; no build artifact is involved |
| Runtime surfaces remain distinct | yes | I5 is a CI-time validator on source files; it does not collapse with `cn build` (I1) or with the cnos.core skill loader (which is the runtime consumer of frontmatter) |
| Degraded paths visible and testable | yes | `schemas/skill-exceptions.json` is the visible degraded path; entries are field-specific and `reason`-tagged; the README documents how to shrink it; the validator runs against it on every push (testable by running `--self-test`) |

No "no" answers. Architecture Check passes. Approval is therefore not
gated on architecture grounds.

## Findings

| # | Finding | Evidence | Severity | Type |
|---|---------|----------|----------|------|
| F1 | `release/SKILL.md` prose still references a non-existent "writing" skill in two sites in the file α touched. The actual skill is named `write` (`cnos.core/skills/write/SKILL.md` L2: `name: write`). α's own `calls: -writing` removal removed the static-call surface; the prose surface in the same file was left stale. `find src/packages -name SKILL.md -path '*writing*'` returns 0 matches. | `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` L103 ("If release notes or CHANGELOG wording are being authored, load the writing skill and record it in the CDD Trace.") and L217 ("skills loaded: release, plus writing if used") | C | judgment |
| F2 | Numeric drift in `.cdd/unreleased/301/self-coherence.md`: claims "45 entries" / "45 file-specific entries" three times for `schemas/skill-exceptions.json`. Actual count is **43** (`jq 'length' schemas/skill-exceptions.json` → 43). | `.cdd/unreleased/301/self-coherence.md` L27 ("a debt ledger of 45 file-specific entries"), L53 ("`schemas/skill-exceptions.json` (45 entries…)"), L144 ("Exception list is large — 45 entries") | A | mechanical |
| F3 | Review-readiness signal head SHA lags actual branch HEAD by one commit. The signal in `.cdd/unreleased/301/self-coherence.md` L187 names `head SHA 6e5ce21`; actual HEAD on `origin/claude/cnos-alpha-tier3-skills-MuE2P` is `ed5f218`. The `α(#301): self-coherence — refresh head SHA for the readiness commit` commit (`ed5f218`) updated the signal from `8adfd44` (implementation) to `6e5ce21`, and the act of committing that update advanced HEAD to `ed5f218` — leaving the new value once again one commit behind. | `git rev-parse origin/claude/cnos-alpha-tier3-skills-MuE2P` = `ed5f218`; `.cdd/unreleased/301/self-coherence.md` L187 names `6e5ce21` | A | mechanical |

## Regressions Required (D-level only)

None. No D-level findings.

## Required fixes for round 2

α must fix all three findings on this branch before merge (per
`review/SKILL.md` §7.0 — no "approved with follow-up"):

1. **F1** — change `release/SKILL.md` L103 and L217 references from
   "writing" to `write` (the canonical skill name). Optional alternative
   if α decides the load relationship is no longer real: drop the prose
   reference entirely. Wider rename debt across `CDD.md` L611 / L785
   and `eng/skills/eng/README.md` L165 ("writing bundle") is **out of
   scope** for #301 — it is pre-existing rename debt from a prior cycle
   and may be filed as a follow-up issue. F1's required scope is the
   touched file only.
2. **F2** — replace the three "45" sites in self-coherence.md with the
   actual count (43), or rephrase to a non-numeric form. Per
   `review/SKILL.md` §2.1.3 "numeric / value-repetition" rule, all sites
   must be fixed in one fix-round, not first-occurrence-only.
3. **F3** — either drop the `head SHA` value from the readiness-signal
   line (CI verifies actual HEAD anyway, so the value carries no load
   the polling cycle does not already carry), or change the convention
   so that the documented "head SHA" names the implementation commit
   (`8adfd44`) rather than the readiness commit, breaking the recursive
   refresh-then-stale loop α hit.

After α fixes the three findings on this branch and CI is green on the
new HEAD, β will re-review and issue the merge instruction (§β step 8).

## Notes

### Local validation reproduces α's claims

β fetched CUE `v0.13.2` (matching the pinned CI version) and ran the
validator inside a worktree at branch HEAD `ed5f218`:

- `./tools/validate-skill-frontmatter.sh --self-test` → **rc=0**;
  3 positive fixtures pass, 4 negative fixtures fail with the expected
  `.expect` substrings (`scope: incomplete value`, `conflicting values`,
  `calls-target-exists`, `required-or-excepted`).
- `./tools/validate-skill-frontmatter.sh` → **rc=0**;
  `✓ 56 SKILL.md validated; no findings.`
- `./tools/validate-skill-frontmatter.sh --root schemas/fixtures/skill-frontmatter/valid` → **rc=0**;
  `✓ 3 SKILL.md validated; no findings.`
- `./tools/validate-skill-frontmatter.sh --file <single>` → reports
  `✓ <path>` and exits 0.
- `PATH=/usr/bin:/bin ./tools/validate-skill-frontmatter.sh --self-test`
  → **rc=2** (`✗ prerequisite missing: cue`), confirming the fail-fast
  prereq path emits exit 2 distinct from validation-failure exit 1.

### CI green precondition (`review/SKILL.md` §3.7)

β's environment cannot reach `api.github.com` to query the
`skill-frontmatter-check` (I5) job result on `ed5f218` (HTTP 403 on
both the check-runs endpoint and the workflow-runs endpoint). Per CDD
§β step 8 deferral path (env-blocked release mechanics), CI verification
is delegated to δ (operator) at merge time. The merge instruction is
explicit: "Do not merge until CI is green on the round-2 HEAD." Local
proof above is sufficient to refute "the validator does not validate";
it is **not** sufficient to refute "the I5 GitHub-Actions job has not
yet run on this SHA." Both are required before the merge gate opens.

### Scope notes

- F1's wider sibling-fallback (`CDD.md` L611 / L785, `eng/README.md`
  L165) is **pre-existing debt from a prior `writing/` → `write/`
  rename**, not a defect this PR introduced. Per CDD §1.4 small-change
  path and `review/SKILL.md` §7.0 "design-scope deferral", α may file
  a follow-up issue rather than expand #301's scope. The required fix
  on this branch is limited to the file α touched
  (`release/SKILL.md`).
- The issue's non-goal "Rewriting non-conformant skills to add missing
  fields (separate work)" is honored: 46 backfills add hard-gate fields
  only (the minimum to keep CI green), not exception-backed fields.
  The 43-entry exception ledger is the explicit deferral.

### Process observation (informational, not a finding)

α's commit `ed5f218` ("self-coherence — refresh head SHA for the
readiness commit") is a structural pattern that's worth naming for γ's
PRA: a "head SHA = SHA of the commit that lands this readiness file"
convention is recursively self-stale, because committing the refresh
advances HEAD by one commit. Two stable conventions exist —
"head SHA = SHA of the implementation commit" or "drop the value and
let the polling protocol carry HEAD" — either of which makes the
review-readiness signal self-consistent on first write. This is recorded
here for γ; β does not block on it.

