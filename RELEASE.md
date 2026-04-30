# 3.63.0 — CUE-based SKILL.md frontmatter validation in coherence CI

## Outcome

Coherence delta: C_Σ A- (`α A-`, `β A-`, `γ A-`) · **Level:** `L6 (cycle cap: L6)`

CTB v0.1 SKILL.md frontmatter now has a machine-checkable contract that runs on every push. Pre-#301 the LANGUAGE-SPEC declared ten well-formedness constraints (§10) and reserved a fixed vocabulary (§11), but no CI job checked either; a spot survey showed 30+ of 54 SKILL.md files missing fields the spec calls required, and new skills could ship without those fields indefinitely. Post-#301 a CUE schema (`schemas/skill.cue`), a wrapping shell validator (`tools/validate-skill-frontmatter.sh`), a field-specific exception ledger (`schemas/skill-exceptions.json`, 43 entries), positive/negative fixture coverage (`schemas/fixtures/skill-frontmatter/{valid,invalid}/`), and operator-facing schema docs (`schemas/README.md`) plumb a new `skill-frontmatter-check` (I5) coherence-CI job that fails on hard-gate omission, enum violation, broken `calls` target, or unexcepted spec-required-but-exception-backed gap.

## Why it matters

Frontmatter drift was silent. Required fields could be missing for years and no surface caught it; the LANGUAGE-SPEC was normative but unenforced. This release converts that into a build-time gate. The schema is the single source of truth for shape; the script owns mechanism (discovery / extraction / exception handling / filesystem checks); the exception list is a visible debt ledger that is intended to shrink. The 43 entries are the explicit deferred work — each addressable by adding the missing field to the named SKILL.md and removing the entry; each carries a `reason` and a `spec_ref` pointing at the LANGUAGE-SPEC section that names the field.

The cycle also exemplifies σ's `4a0f678` clarification (`cdd: merge is β authority, not δ — δ owns release gates (tag/deploy) only`): β reviewed, β merged, δ owns the tag/deploy gate that follows. The merge action itself required no δ authorization.

## Added

- **CUE schema for SKILL.md frontmatter** (#301): `schemas/skill.cue` enforces field shape, type, and enum constraints across the four strata (hard-gate, spec-required-but-exception-backed, optional-with-default, reserved). Schema is open (`...`) so unknown package-local keys (e.g. `parent:`) pass through per LANGUAGE-SPEC §11.
- **Validation script** (#301): `tools/validate-skill-frontmatter.sh` discovers SKILL.md files under `src/packages/`, extracts frontmatter, runs `cue vet` per file, enforces present-or-excepted on `artifact_class`/`kata_surface`/`inputs`/`outputs`, and resolves static `calls` targets against the package skill root. Modes: default (full repo), `--self-test` (run fixtures), `--root <dir>` (alt root), `--file <path>` (single file), `-h/--help`. Exit codes 0 (pass), 1 (validation failure), 2 (prerequisite missing) are distinct.
- **Field-specific exception ledger** (#301): `schemas/skill-exceptions.json` lists 43 pre-#301 SKILL.md paths with `allowed_missing` (only `artifact_class`/`kata_surface`/`inputs`/`outputs`), `reason`, and `spec_ref`. Hard-gate fields cannot be excepted. Each entry is intended to shrink as the named skill is updated.
- **Positive + negative fixture coverage** (#301): `schemas/fixtures/skill-frontmatter/{valid,invalid}/` with 3 positive (minimal, full with every reserved field, sub-skill for calls-target resolution) and 4 negative cases (missing-hard-gate, bad-enum-scope, broken-calls-target, missing-required-no-exception), each with `.expect` sidecar naming the diagnostic substring. `--self-test` asserts that every positive fixture passes and every negative fixture fails with the expected diagnostic.
- **CI job `skill-frontmatter-check` (I5)** (#301): added to `.github/workflows/build.yml` with `cue-lang/setup-cue@v1.0.1` (action tag pinned) and CUE `v0.13.2` (binary version pinned). Runs `--self-test` then full validation. The `notify` job's `needs` and result aggregation include I5 alongside I1/I2/I4.
- **Schema documentation** (#301): `schemas/README.md` covers what is and isn't validated, the script-vs-schema-vs-filesystem surface boundary, the field stratification table, exception mechanics, the shrink procedure, local-run instructions, CI integration with the pinned versions, and the schema-update procedure when LANGUAGE-SPEC changes.

## Changed

- **`release/SKILL.md`** (#301): pre-existing `calls: -writing` removed from the frontmatter (intra-package-only `calls` resolution per LANGUAGE-SPEC §2.4.1; the prior `writing` skill was renamed to `write` in an earlier cycle and the static-call entry pointed at a non-existent target). Two prose references on L103 + L217 also renamed `writing` → `write` per β R1 F1 (in-touched-file scope; wider `writing`→`write` rename debt across `CDD.md` L611+L785 and `eng/skills/eng/README.md` L165 is pre-existing debt for a follow-up issue).
- **46 SKILL.md frontmatter backfills** under `src/packages/**` (#301): hard-gate fields added where missing (the minimum to keep CI green); 2 YAML strict-quoting normalizations in `cdd/design/SKILL.md` and `cdd/gamma/SKILL.md` for unquoted-colon-space sequences PyYAML accepts but CUE rejects.
- **`CDD.md` Role table + δ section** (`4a0f678`, σ): merge is β authority — no δ authorization required. δ's authorities are `(tag/deploy)`, not `(merge/tag/branch)`. Step 9 in §5.3 is `Merge | β | git merge of cycle/{N} into main (β authority — no δ required)`. Bundled with #301 because the cycle exemplifies the clarified contract by execution.

## Validation

- Pre-merge non-destructive merge test: `git merge --no-ff --no-commit` of `origin/claude/cnos-alpha-tier3-skills-MuE2P` against `origin/main = 4a0f678` produced one auto-merge on `cdd/gamma/SKILL.md` (disjoint regions, no human intervention needed) and zero unmerged paths.
- I5 validator on the merge tree: `--self-test` rc=0 (all 3 positive + 4 negative fixtures behave as expected with the pinned `.expect` substrings); full validation rc=0 (`✓ 56 SKILL.md validated; no findings.`).
- Exit-code-2 prerequisite path verified locally with `PATH=/usr/bin:/bin ./tools/validate-skill-frontmatter.sh --self-test` → `✗ prerequisite missing: cue` rc=2, distinct from validation-failure exit 1.
- Local CUE binary used for verification was the pinned CI version `v0.13.2` (matches `.github/workflows/build.yml`).
- Post-merge full validation on `main`: `--self-test` rc=0; full validation rc=0 (56/56 SKILL.md). The 6 `cdd/*/SKILL.md` files updated by cycle #287 (3.62.0) satisfy the new schema without any exception entry.

## Known Issues

- Wider `writing` → `write` rename debt remains in `CDD.md` L611 + L785 and `src/packages/cnos.eng/skills/eng/README.md` L165 ("writing bundle") — pre-existing debt from an earlier `writing/` → `write/` rename, out of scope for #301 (issue scope is frontmatter validation, not body prose). The I5 validator does not catch these because they are body content. Fileable as a follow-up issue.
- The 43 exception entries are by design and intended to shrink. Each is one-skill-one-PR addressable.
