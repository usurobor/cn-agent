# schemas/

Machine-checkable schemas for cnos artifacts. Currently scoped to
`SKILL.md` frontmatter (I5). The directory will grow as new schema-bearing
artifacts come under coherence-CI gating.

## What `skill.cue` validates

`schemas/skill.cue` defines a [CUE](https://cuelang.org/) schema for the
YAML frontmatter at the top of every `src/packages/**/SKILL.md`, per
[CTB LANGUAGE-SPEC.md](../docs/alpha/ctb/LANGUAGE-SPEC.md) §2 and §11.

The schema enforces:

- **Field shape** — type of each declared field (`name`, `description`,
  list shape of `triggers`, mapping shape of `calls_dynamic`).
- **Enum domains** — `scope ∈ {global, role-local, task-local}`,
  `visibility ∈ {internal, public}`, `artifact_class ∈ {skill, runbook,
  reference, deprecated}`, `kata_surface ∈ {embedded, external, none}`.
- **Hard-gate presence** — `name`, `description`, `governing_question`,
  `triggers`, `scope` must be present and non-empty in every skill.
- **Open extension** — unknown package-local keys (e.g. `parent`) pass
  through; the loader must ignore unknown keys per LANGUAGE-SPEC §11.

## What `skill.cue` does NOT validate

This is the v0 of skill-frontmatter validation. It deliberately does NOT
attempt:

- **Witness or close-out semantics.** That is the job of `ctb-check` v0
  (#303), which validates CTB v0.2 close-out artifacts. This schema only
  validates frontmatter shape.
- **Markdown body structure.** Headings, sections, code blocks, etc., are
  not checked. The skill-skill (`cnos.core/skills/skill/SKILL.md`)
  governs the body.
- **`calls_dynamic` target existence.** Dynamic call sources resolve at
  runtime; only the static `calls` list is checked for filesystem
  existence.
- **Cross-package call graphs.** `calls` resolution is intra-package only,
  per LANGUAGE-SPEC §2.4.1.
- **CTB v0.2 fields.** v0.2 (agent-type/agent-module distinction) is
  drafted but not promoted; this schema validates v0.1 only.

## Surface boundary: schema vs. script

The validation has two surfaces with one rule each:

- **`schemas/skill.cue`** — owns shape, type, and enum constraints. CUE
  fails any frontmatter whose fields drift from this declaration. Pure;
  no filesystem I/O.
- **`tools/validate-skill-frontmatter.sh`** — owns file discovery,
  frontmatter extraction, exception handling, and `calls`
  filesystem-existence checks. This is shell logic that wraps
  `cue vet` per skill.

Do not move shape/type/enum rules into the shell, and do not move
discovery/exceptions into CUE. Each side stays small.

## Field stratification

Every frontmatter field falls into one of four strata. Each is enforced
differently:

| Stratum | Fields | Schema | Script | Exceptionable? |
|---|---|---|---|---|
| Hard gate | `name`, `description`, `governing_question`, `triggers`, `scope` | required | (pass-through) | **no** |
| Spec-required, exception-backed | `artifact_class`, `kata_surface`, `inputs`, `outputs` | optional | enforces present-or-excepted | yes |
| Optional with default | `visibility` (default `internal`) | optional | (pass-through) | n/a |
| Reserved — validated only if present | `requires`, `calls`, `calls_dynamic`, `runs_after`, `runs_before`, `excludes` | optional | filesystem check on `calls` only | n/a |
| Unknown package-local | anything else (e.g. `parent`) | accepted (open struct) | accepted | n/a |

Hard-gate fields cannot be excepted because the loader needs them to
dispatch the skill at all (name, scope, triggers) or to communicate what
the skill is for (description, governing_question). The exception list
exists for fields the loader does not need before dispatch.

## How exceptions work

`schemas/skill-exceptions.json` is a JSON array of per-file entries:

```json
[
  {
    "path": "src/packages/cnos.core/skills/agent/coherent/SKILL.md",
    "allowed_missing": ["artifact_class", "kata_surface", "inputs", "outputs"],
    "reason": "Pre-#301 skill — v0.1 spec-required-but-exception-backed fields not yet authored.",
    "spec_ref": "LANGUAGE-SPEC.md §2.1, §2.3"
  }
]
```

Rules:

- Exceptions are **field-specific**, not file-wide. `allowed_missing` lists
  the specific fields the script will not enforce for this path. Every
  other field is enforced normally.
- Hard-gate fields (`name`, `description`, `governing_question`,
  `triggers`, `scope`) **cannot** appear in `allowed_missing`. The script
  flags them with the schema, not via the exception list.
- Each entry needs a `reason` (JSON has no comments) and ideally a
  `spec_ref` pointing at the spec section the exception defers.
- The exception list is a **debt ledger**, not a feature flag. Each entry
  represents work to do — adding the missing field — not a permanent
  carve-out. The list is intended to shrink.

## How to shrink the exception list

For each entry, the addressable change is mechanical:

1. Open the SKILL.md named in `path`.
2. For each field in `allowed_missing`, add the field to the frontmatter.
   - `artifact_class`: pick the closest of `skill`, `runbook`,
     `reference`, `deprecated`.
   - `kata_surface`: pick the closest of `embedded`, `external`, `none`.
   - `inputs` / `outputs`: lift from the body's contract section if the
     skill has one; otherwise author from the skill's stated purpose.
3. Remove the entry from `schemas/skill-exceptions.json`.
4. Run the validator locally (see below). It must remain green.

## How to run locally

You need [CUE](https://cuelang.org/) and `jq` on your `PATH`. Then:

```bash
# Validate every src/packages/**/SKILL.md against schema + exceptions.
./tools/validate-skill-frontmatter.sh

# Run the positive/negative fixture suite (#301 AC8).
./tools/validate-skill-frontmatter.sh --self-test

# Validate a single skill file (useful while editing).
./tools/validate-skill-frontmatter.sh --file src/packages/.../SKILL.md
```

`NO_COLOR=1` suppresses color. Findings are emitted one per line in the
form:

```
✗ <path> :: <field> :: <rule> :: <reason> :: <suggested fix>
```

## CI integration

The job `skill-frontmatter-check` in `.github/workflows/build.yml` (job ID
**I5**) runs both `--self-test` and the full validation on every push and
pull request. The `notify` job aggregates I5 alongside I1, I2, and I4 for
the Telegram summary.

The job pins:

- `cue-lang/setup-cue@v1.0.1` (action tag)
- `version: v0.13.2` (CUE version)

Both pins are intentional and should be bumped together when the schema
is re-tested against a new CUE release.

## How to update the schema when the CTB spec changes

When [LANGUAGE-SPEC.md](../docs/alpha/ctb/LANGUAGE-SPEC.md) changes:

1. Decide whether the change adds, removes, or re-strata a field.
2. Update `schemas/skill.cue`:
   - Adding a hard-gate field: add a non-`?` field and update the
     stratification table above. Expect mass exception-list churn.
   - Adding a spec-required-but-exception-backed field: add a `?`
     field in CUE and add it to `SPEC_REQUIRED_EXCEPTION_BACKED` in
     `tools/validate-skill-frontmatter.sh`.
   - Adding a reserved field validated only if present: add a `?` field
     in CUE; no script change.
   - Tightening an enum: update the disjunction in CUE; expect to amend
     skills that fall outside the new domain.
3. Update fixtures in `schemas/fixtures/skill-frontmatter/`:
   - At least one positive fixture must use the new field.
   - At least one negative fixture must violate the new constraint.
4. Re-run `./tools/validate-skill-frontmatter.sh --self-test` and the
   full validation; resolve any new findings before commit.

## Related issues and specs

- [#301](https://github.com/usurobor/cnos/issues/301) — this work, CUE-based skill frontmatter validation in coherence CI.
- [#289](https://github.com/usurobor/cnos/issues/289) — CTB v0.2 stabilization (this is the machine-checkable layer of v0.2 AC5).
- [#303](https://github.com/usurobor/cnos/issues/303) — `ctb-check` v0 (witness / close-out validation; out of scope here).
- `docs/alpha/ctb/LANGUAGE-SPEC.md` — v0.1 spec governing field shape, stratification, and reserved vocabulary.
- `docs/alpha/ctb/LANGUAGE-SPEC-v0.2-draft.md` — v0.2 draft (out of scope until promoted).
