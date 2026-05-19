---
cycle: 379
date: "2026-05-19"
issue: "https://github.com/usurobor/cnos/issues/379"
merge_sha: "a3bf7892"
findings_count: 3
patches_count: 3
mcas_count: 0
no_patch_count: 0
---

# CDD Iteration — #379

Three `cdd-*-gap` findings surfaced during γ's post-merge CI verification: one skill-gap (F1, β/SKILL.md row 3 under-specification) and two wiring-class tooling-gaps (F2 / F3, peer surfaces not updated for the cycle's intended change). All three are patched in the same PRA commit; no findings deferred to a next-MCA.

Format per `cdd/post-release/SKILL.md` Step 5.6b.

## F1: β/SKILL.md §pre-merge gate row 3 under-specifies the contract-validator set

- **Source:** γ-triage / PRA §3 "What went wrong" item 2 / PRA §4b §9.1 trigger (loaded skill failed to prevent a finding) — `docs/gamma/cdd/3.78.0/POST-RELEASE-ASSESSMENT.md`
- **Class:** `cdd-skill-gap`
- **Trigger:** §9.1 — loaded skill failed to prevent a finding (β/SKILL.md row 3 was loaded; the gate fired; the application stopped at `go test` because row 3 said "the cycle's own validator" without enumerating which validators apply when)
- **Description:** β's pre-merge gate row 3 ("Non-destructive merge-test") prescribed running "the cycle's own validator (or any CI-equivalent the cycle ships)" on the merge tree. β ran `go test ./internal/activate/...` and `go vet` on the merge tree (documented in `beta-review.md` §"Branch CI state") and stopped. The cycle ships a new SKILL.md (validate-skill-frontmatter.sh applies) and displaces `cn activate`'s `## Read first` ordering (R5-activate kata applies). Neither validator was run. Both regressed on the merge commit and produced post-merge CI red.
- **Root cause:** Under-specified row-3 wording — "the cycle's own validator" is the natural read for "the one validator most obviously associated with the cycle's primary surface," not "every validator whose contract the cycle's surface touches." The exhaustive list lives in `release/SKILL.md §2.1` but β/SKILL.md row 3 did not cross-reference it. Application gap was made cheap by the wording.
- **Disposition:** **patch-landed**
- **Patch:** this PRA commit (β/SKILL.md §pre-merge gate row 3 enumerated list with the four validator classes named explicitly: `validate-skill-frontmatter.sh`, `cn-cdd-verify`, `check-version-consistency.sh`, and per-cycle katas under `src/packages/cnos.kata/katas/{N}/`; R5-activate named as the activate-surface kata exemplar; #379 cited inline as empirical anchor)
- **Affects:** `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` §Pre-merge gate row 3

## F2: agent/activate skill `calls:` paths violate validator's package-skill-root resolution

- **Source:** γ-triage / PRA §3 "What went wrong" item 1a — surfaced by `gh run view 26093805709` showing I5 (SKILL.md frontmatter validation) failure with 7 findings all in the new `src/packages/cnos.core/skills/agent/activate/SKILL.md`
- **Class:** `cdd-tooling-gap` (the cycle's primary deliverable, a new SKILL.md, was authored without running the validator that gates SKILL.md frontmatter; the failure is mechanical and bounded)
- **Trigger:** avoidable tooling/environmental failure (per §9.1) — CI red on merge commit
- **Description:** The new skill's `calls:` field listed seven entries with `cnos.core/...` prefixes (e.g. `cnos.core/doctrine/KERNEL.md`, `cnos.core/skills/agent/cap/SKILL.md`). `tools/validate-skill-frontmatter.sh` resolves `calls:` entries against the package-skill-root (`src/packages/cnos.core/skills/` for cnos.core), so `cnos.core/doctrine/KERNEL.md` resolved to `src/packages/cnos.core/skills/cnos.core/doctrine/KERNEL.md` which does not exist. All 7 entries failed the `calls-target-exists` check. Pre-cycle CI green → post-cycle CI red.
- **Root cause:** α authored `calls:` with the path form that read naturally as a package-qualified locator, not the form the I5 validator expects (package-skill-root-relative). The validator's resolution rule is documented in `tools/validate-skill-frontmatter.sh::package_skill_root_of` but the cycle's authoring loop did not exercise it. Combined with F1, this regression escaped both α's authoring-time check and β's pre-merge gate.
- **Disposition:** **patch-landed**
- **Patch:** this PRA commit — `calls:` entries rewritten to package-skill-root-relative paths (`../doctrine/KERNEL.md` for the doctrine reference; `agent/cap/SKILL.md`, `agent/clp/SKILL.md`, etc. for sibling skills). Validator green after patch: `./tools/validate-skill-frontmatter.sh` → "66 SKILL.md validated; no findings."
- **Affects:** `src/packages/cnos.core/skills/agent/activate/SKILL.md` (calls: frontmatter); rooted in `tools/validate-skill-frontmatter.sh` resolution semantics; backstopped by F1's β/SKILL.md row 3 patch (future cycles will run the validator pre-merge as part of the enumerated row-3 list)

## F3: R5-activate kata P10 assertion was not updated when AC3 displaced the `## Read first` ordering

- **Source:** γ-triage / PRA §3 "What went wrong" item 1b — surfaced by `gh run view 26093805709` showing Package verification failure at `R5-activate / P10: read-first order incorrect (persona=3 operator=4 kernel=1 deps=5 refl=5)`
- **Class:** `cdd-tooling-gap` (the cycle's intended change displaced a peer test contract; the peer contract was not updated to match)
- **Trigger:** avoidable tooling/environmental failure (per §9.1) — CI red on merge commit
- **Description:** The R5-activate kata's P10 test asserts the `## Read first` ordering produced by `cn activate`. The pre-#379 canonical ordering was `persona < operator < kernel < deps < reflection`; AC3 of #379 displaced this to `kernel < ca-skills < persona < operator < hub-state < identity`. α's self-coherence §Debt 2 surfaced the observable-output delta honestly but did not update the kata's assertion. β's pre-merge gate row 3 did not run R5-activate against a fresh build, so the regression reached the merge commit.
- **Root cause:** α applied harness-audit discipline (§2.4 in `alpha/SKILL.md`) to the schema-level surface (the marker-bounded ordering block consumed by `parseReadFirstOrderBlock`) but did not extend the audit to the higher-level peer (the kata that asserts the rendered output). The cycle's "harness audit" treatment was bounded to producers/consumers of the parsing schema, not consumers of the rendered prompt. Combined with F1, the kata regression escaped both α's authoring-time check and β's pre-merge gate.
- **Disposition:** **patch-landed**
- **Patch:** this PRA commit — `src/packages/cnos.kata/katas/R5-activate/run.sh` P10 assertion updated from `persona < operator < kernel < deps < refl` to `kernel < persona < operator < {deps, refl}` (deps and refl share one hub-state line in the new ordering, so their relative order is no longer asserted). `src/packages/cnos.kata/katas/R5-activate/kata.md` P10 documentation updated to name the new six-item canonical order, the source-of-truth (`src/packages/cnos.core/skills/agent/activate/SKILL.md §4.1`), and the deprecated pre-#379 form. CI Package-verification expected green on the next push.
- **Affects:** `src/packages/cnos.kata/katas/R5-activate/run.sh` (P10 ordering assertion); `src/packages/cnos.kata/katas/R5-activate/kata.md` (P10 documentation); rooted in cycle #379 AC3 displacement; backstopped by F1's β/SKILL.md row 3 patch (future activate-surface cycles will run R5-activate pre-merge per the enumerated row-3 kata-coverage clause).

## Summary

| Finding | Class | Disposition | Patch / MCA path |
|---------|-------|-------------|------------------|
| F1 — β/SKILL.md row 3 under-specifies validator set | cdd-skill-gap | patch-landed | this PRA commit (β/SKILL.md §pre-merge gate row 3 enumeration) |
| F2 — activate skill `calls:` paths violate I5 resolution | cdd-tooling-gap | patch-landed | this PRA commit (activate SKILL.md `calls:` frontmatter rewrite) |
| F3 — R5-activate kata P10 not updated for new ordering | cdd-tooling-gap | patch-landed | this PRA commit (R5-activate run.sh + kata.md P10 update) |

All three findings have backstop relationships: F1 (skill patch) is the structural fix that prevents the F2 / F3 class from recurring (β's pre-merge gate now enumerates the validators that would have caught F2 and F3 at pre-merge). F2 and F3 are mechanical patches that restore CI green on the next push.

## Aggregator update

`.cdd/iterations/INDEX.md` row added at cycle #379 close (this commit); after the cycle dir move per `gamma/SKILL.md §2.6`, this artifact's path is `.cdd/releases/3.78.0/379/cdd-iteration.md`.
