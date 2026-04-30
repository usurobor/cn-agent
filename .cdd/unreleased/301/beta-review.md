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


## Round 2

**Verdict:** APPROVED (provisional on CI green — see merge instruction)

**Branch:** `claude/cnos-alpha-tier3-skills-MuE2P`
**Base:** `a8e67b7` (origin/main)
**Round-1 verdict commit:** `40af6c0` (β round-1 RC, this file)
**Fix-round commit:** `171188e` (α round-2 F1+F2+F3 fixes)
**Appendix commit:** `55642db` (α round-2 appendix to self-coherence.md)
**Head reviewed:** `55642db`
**Branch CI state:** **unverified-from-β-environment** — required for merge; see merge instruction
**Merge instruction:** **Do not merge until CI is green on `55642db`.** Once δ confirms `skill-frontmatter-check` (I5) and all other required jobs (`go`, `binary-verify`, `package-verify`, `package-source-drift`, `protocol-contract-check`, `link-check`) report `success` on `55642db`, β (or δ acting as β per operator override) executes:

```
git merge --no-ff origin/claude/cnos-alpha-tier3-skills-MuE2P \
  -m "Closes #301: infra(ci) — CUE-based SKILL.md frontmatter validation (I5)"
```

then proceeds to release per `release/SKILL.md` (version bump, CHANGELOG, RELEASE.md, tag, push), and β writes the close-out per CDD §β step 9.

### Fixed this round

| Finding (round-1) | Severity | Fix verified at HEAD `55642db` |
|---|---|---|
| F1 — `release/SKILL.md` L103 + L217 prose still references non-existent "writing" skill | C | Both prose sites renamed `writing` → `write`. `grep -n "writing" src/packages/cnos.cdd/skills/cdd/release/SKILL.md` → 0 hits. Fix is in scope (the touched file only); wider rename debt across `CDD.md` L611 / L785 and `eng/skills/eng/README.md` L165 ("writing bundle") remains pre-existing debt for a follow-up issue, as scoped in round-1 verdict. |
| F2 — self-coherence.md "45" claim drift (3 sites) | A | All three sites updated to `43`. `grep -nE "(45 entries\|45 file-specific)" .cdd/unreleased/301/self-coherence.md` → 0 hits. Matches actual `jq 'length' schemas/skill-exceptions.json` → 43. |
| F3 — readiness-signal head SHA recursively self-stales | A | Convention changed: signal names the **implementation commit SHA** (`8adfd44`), not the readiness commit. New round-2 signal at L239 reads `implementation SHA \`8adfd44\` \| round-2 fix SHA \`171188e\``; both are stable references that do not lag HEAD. "SHA convention" paragraph at L189–195 documents the choice. Branch HEAD as visible to β is still carried by the polling protocol. |

No new findings introduced by the fix-round. Re-audit:

- `release/SKILL.md`: prose-only edit; CUE schema does not parse body prose; no schema impact.
- `.cdd/unreleased/301/self-coherence.md`: not under `src/packages/`, not in I5's scope.
- All 56 `src/packages/**/SKILL.md` unchanged from round 1; validator output unchanged.

### Local re-verification at round-2 HEAD `55642db`

β fetched the round-2 HEAD into the worktree and re-ran the validator (CUE `v0.13.2`, `jq 1.7`):

- `./tools/validate-skill-frontmatter.sh --self-test` → **rc=0**; all 3 positive + 4 negative fixtures behave as expected.
- `./tools/validate-skill-frontmatter.sh` → **rc=0**; `✓ 56 SKILL.md validated; no findings.`

Local proof confirms the validator behaves correctly on the round-2 tree. It does not substitute for the `skill-frontmatter-check` (I5) GitHub-Actions job's own `success` on `55642db` — that remains δ's gate, and is the one outstanding precondition for merge.

### Round-1 → round-2 narrowing pattern

α addressed every named finding in scope on the same branch (no new branch, no follow-up issue used as a substitute for the fix). The fix-round is two commits: implementation (`171188e`) + readiness-signal appendix (`55642db`). Pattern is clean: one round of narrowing, no scope drift, no new findings opened by the fix. β does not need a round 3.

### Outstanding precondition

The single outstanding precondition is **δ confirmation of CI green on `55642db`**, per `review/SKILL.md` §3.7. β's environment cannot reach `api.github.com` (HTTP 403 on check-runs and workflow-runs endpoints); δ holds the verification surface. Once δ confirms, the merge command above is the merge instruction; if δ is unavailable and the operator authorizes γ or β to execute the merge, the same command applies.

After merge, β writes `.cdd/unreleased/301/beta-closeout.md` (review context across rounds + release evidence) per CDD §β step 9 and the release skill's §2.10 close-out hook. γ owns the PRA.


## Round 2 — Pass 2 (re-review against fresh `origin/main`)

**Trigger:** the new `beta/SKILL.md` Role Rule 1 on `origin/main` (shipped in cycle #287, version 3.62.0) requires β to **re-fetch `origin/main` synchronously before computing the review-diff base for each review pass**. β's round-1 and round-2 review passes used `origin/main = a8e67b7` (last fetched at β intake). β re-fetched and re-evaluated.

**Updated state:**

| Surface | Round-1/Round-2 base | Pass-2 (current) |
|---|---|---|
| `origin/main` | `a8e67b7` (β intake) | `9d6a0fa` (sigma `docs: rewrite README for v4 convergence`) |
| Cycle branch HEAD | `55642db` → `acfa0cf` (β round-2 verdict) | `acfa0cf` (unchanged) |
| Cycle branch base (merge-base) | `a8e67b7` | `a8e67b7` (unchanged — branch never rebased) |

**What advanced on `main` during the review window:** cycle #287 shipped (v3.62.0). Touched these CDD-skill files: `CDD.md`, `alpha/SKILL.md`, `beta/SKILL.md`, `gamma/SKILL.md`, `operator/SKILL.md`, `review/SKILL.md`. Also touched: `README.md`, `CHANGELOG.md`, `RELEASE.md`, `VERSION`, `cn.json`, all three `cn.package.json` files, plus the `.cdd/releases/3.62.0/287/` and `docs/gamma/cdd/3.62.0/POST-RELEASE-ASSESSMENT.md` artifacts.

### File overlap audit

| File touched on both sides | Conflict on auto-merge? |
|---|---|
| `src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md` | no — auto-merged cleanly. Cycle's edit was a one-line YAML quoting (`- "delta gate results (observable via git: tags, branch state)"`) at L22; main's edit restructured the `## γ algorithm` body. Disjoint regions; auto-merge succeeds without human intervention. |

No other file overlap. `cdd/release/SKILL.md` and `cdd/design/SKILL.md` (cycle-only) and `cdd/CDD.md`, `cdd/alpha/SKILL.md`, `cdd/beta/SKILL.md`, `cdd/operator/SKILL.md`, `cdd/review/SKILL.md` (main-only) do not collide.

### Merge tree validation

β built the merge tree (`git merge --no-ff --no-commit origin/claude/cnos-alpha-tier3-skills-MuE2P` against `origin/main = 9d6a0fa`) in a throwaway worktree:

- Auto-merge succeeded; zero unmerged paths.
- `./tools/validate-skill-frontmatter.sh --self-test` on the merge tree → **rc=0** (3 positive + 4 negative as expected).
- `./tools/validate-skill-frontmatter.sh` on the merge tree → **rc=0**; `✓ 56 SKILL.md validated; no findings.`
- All 6 `cdd/*/SKILL.md` files updated by main pass the I5 schema after merge — main did not introduce frontmatter that the cycle's new schema would reject.

### Findings stability under fresh base

| Round-1 finding | Pass-2 status |
|---|---|
| F1 — `release/SKILL.md` prose references "writing" | **stable** — content-level finding, independent of `origin/main` advance. α's fix at `171188e` still applies cleanly post-merge (cycle's edit is the only one to L103/L217). |
| F2 — self-coherence "45 entries" drift | **stable** — internal artifact, independent of main. α's fix unchanged. |
| F3 — readiness signal head SHA recursion | **stable** — convention-level finding, independent of main. α's fix unchanged. |

**No new findings** emerge against the fresh `origin/main`. None of the round-1 findings collapse on fresh fetch (contrast #287 R1 F1/F2 which collapsed on fresh main — those were scope-drift findings against a stale spec; #301's F1–F3 are content-level findings against the issue's ACs, immune to base advance).

### Approval status

**APPROVED stands.** All three round-1 findings are fixed at `171188e`; round-2 fix appendix at `55642db`; β round-2 verdict at `acfa0cf`. The merge into the current `origin/main = 9d6a0fa` is conflict-free and the I5 validator passes on the merge result.

### Updated merge instruction

The merge target advanced; the merge command is unchanged in shape (the branch name resolves dynamically):

```
git fetch --verbose origin main
git checkout main && git pull --ff-only
git merge --no-ff origin/claude/cnos-alpha-tier3-skills-MuE2P \
  -m "Closes #301: infra(ci) — CUE-based SKILL.md frontmatter validation (I5)"
git push origin main
```

After merge, proceed with release per `release/SKILL.md` (§2.2 version decision, §2.3 VERSION-first stamp, §2.4 CHANGELOG, §2.5 RELEASE.md, §2.5a move `.cdd/unreleased/301/` → `.cdd/releases/{X.Y.Z}/301/`, §2.6 commit + tag + push, §2.7 wait for release CI, §2.8 deploy, §2.9 validate). β writes `.cdd/unreleased/301/beta-closeout.md` per CDD §β step 9.

### Outstanding precondition

Same as round 2: **δ confirmation of CI green on the cycle branch HEAD `acfa0cf`** (β env 403 on `api.github.com`). The merge tree's local validator pass refutes "the validator does not validate"; it does not refute "GitHub-Actions has not yet completed `success` on this SHA." The latter remains δ's gate.

If the operator authorizes β (or δ) to execute the merge in this session — given (a) auto-merge clean, (b) merge-tree validator green, (c) all round-1 findings addressed, (d) approval recorded — β can proceed; otherwise, the verdict + merge instruction above is the durable hand-off and δ executes when CI is verified.

