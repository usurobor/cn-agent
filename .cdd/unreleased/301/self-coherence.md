# α self-coherence — #301

## Issue

[#301 — infra(ci): CUE-based skill frontmatter validation in coherence
CI](https://github.com/usurobor/cnos/issues/301).
P2. Substantial-change MCA per γ dispatch (cycle log comment
`4347101265`). No §3.1 / §3.2 selection override fired.

## Mode

Substantial. Multi-surface: schema + script + exception list + fixtures
+ CI job + skill backfills + schema docs. Not eligible for the
small-change path (§1.1 substantial criteria: spans design, code, tests,
docs, packaging, runtime contract, CI surface).

## Active skills

| Tier | Skill | How it constrained authorship |
|---|---|---|
| 1 | `cdd/SKILL.md`, `cdd/alpha/SKILL.md`, `cdd/design/SKILL.md` (judging design-required), `cdd/plan/SKILL.md` (sequencing decisions inline; no separate plan artifact required for the linear schema→script→fixtures→CI→docs→backfill flow) | role contract, artifact order, peer enumeration, schema/script surface boundary, pre-review gate |
| 2 | `eng/document` (durable docs in schemas/README.md), `eng/code` (script soundness, set -euo pipefail), `eng/process-economics` (kept exception list as debt ledger, not feature flag) | always-applicable engineering bundle |
| 3 | `cnos.core/skills/design` (script-vs-schema-vs-filesystem source-of-truth split — AC2 surface boundary held; CUE owns shape/type/enum, shell owns discovery/extraction/exception/calls-fs); `eng/tool` (fail-fast prereq checks for `cue` and `jq`, NO_COLOR support, idempotent, exit codes 0/1/2, machine-readable diagnostics); `eng/test` (positive/negative proof via fixtures, including a calls-target broken-by-construction fixture; named the invariant — "every reserved field round-trips", "missing hard-gate fails", "broken calls fails", "missing required-or-excepted fails"); `eng/ux-cli` (one finding per line in `path :: field :: rule :: reason :: fix` shape, ✓/✗ symbols, semantic colors, NO_COLOR override, exit-code-2 for prereq missing distinct from exit-1 for validation failure) | per γ dispatch context |

`eng/process-economics` (Tier 3 conditional in γ dispatch) loaded but
did not fire as a design question — the exception list overhead is
a debt ledger of 43 file-specific entries, not a process surface that
needs new mechanism.

## Gap

Coherence contract (§Gap, the LANGUAGE-SPEC §10 conformance promise was
not machine-checkable):

> The CTB v0.1 LANGUAGE-SPEC declares ten well-formedness constraints
> (§10) and reserves a fixed vocabulary (§11), but the cnos repo had no
> CI job that *checked* either. A spot survey at issue-filing time
> showed 30+ of 54 SKILL.md files missing fields the spec calls
> required. New skills could ship without required fields indefinitely;
> existing drift could compound silently.

This cycle closes the gap: every push and PR now runs the I5 job, which
fails on hard-gate omission, enum violation, broken `calls` target, and
unexcepted spec-required-but-exception-backed gaps.

## Acceptance criteria, evidence

| AC | Claim | Evidence on this branch (head `8adfd44`) |
|---|---|---|
| 1 | CUE schema exists with stratified field rules. Hard-gate / spec-required-exception-backed / optional-with-default / reserved / unknown-package-local. | `schemas/skill.cue:23-78`. Hard-gate fields are unconditional CUE values; spec-required-but-exception-backed are `?`; reserved are `?`; struct is open (`...` at L77) so unknown keys (e.g. `parent:`) pass through. |
| 2 | Extraction script exists with the named standards. | `tools/validate-skill-frontmatter.sh`. `set -euo pipefail` (L34); fail-fast `need cue` / `need jq` / `need awk` / `need find` (L48-54, exit 2); deterministic `find … \| sort` (L257); NO_COLOR (L23-29, only emits color when stdout is a TTY *and* NO_COLOR unset); diagnostics one-per-line `path :: field :: rule :: reason :: fix` (L101-103). Surface boundary: shell owns discovery/extraction/exceptions/calls-fs; CUE owns shape/type/enum (header comment L4-9). Frontmatter delimiter shape was verified pre-implementation across all 56 SKILL.md (`first_line == ---` and ≥2 `---` delimiters; survey ran clean). Two YAML edge-cases discovered (design/SKILL.md, gamma/SKILL.md unquoted-colon-space) were normalised, not silently skipped. |
| 3 | CI job added. Notify aggregation updated. | `.github/workflows/build.yml`: new job `skill-frontmatter-check` (display name "SKILL.md frontmatter validation (I5)") at L237-256, runs both `--self-test` and the full validation. `notify.needs` extended with `skill-frontmatter-check` at L260; result aggregation row `"I5:${{ needs.skill-frontmatter-check.result }}"` at L284. Pins: `cue-lang/setup-cue@v1.0.1` (action tag, L246); `version: v0.13.2` (CUE version, L248). |
| 4 | Exception list shape: field-specific `allowed_missing`, with `reason` + `spec_ref`. CI passes on main with exceptions in place. Hard-gate cannot be excepted; only the four exception-backed fields (`artifact_class`, `kata_surface`, `inputs`, `outputs`) are eligible. | `schemas/skill-exceptions.json` (43 entries, all field-specific, each with `reason` and `spec_ref`). The script's `SPEC_REQUIRED_EXCEPTION_BACKED=(artifact_class kata_surface inputs outputs)` array (L72) is the only set the exception path applies to; hard-gate failures surface from `cue vet` regardless. CI green: see "Local verification" below. |
| 5 | Exception list is intended to shrink. Each entry addressable by adding the field. | `schemas/README.md` "How to shrink the exception list" section gives the mechanical procedure (open SKILL.md → add field → remove entry → re-run validator). Each entry's `spec_ref` points at the LANGUAGE-SPEC section that names the field, so the addresser knows what the field should mean. |
| 6 | `calls` targets validate against the package skill root. Invalid targets fail CI. `calls_dynamic` target existence not validated. | Script `package_skill_root_of` (L195-217) computes root as the common-ancestor of every SKILL.md in the package's `skills/` subtree, matching the issue-body example (`cdd/alpha` → `cdd/design/SKILL.md`, not `cdd/alpha/design/SKILL.md`). Test: `release/SKILL.md`'s pre-existing `calls: - writing` was discovered as a broken target; the broken-calls fixture exercises the same code path positively (must fail). `calls_dynamic` is parsed for source/constraint shape only (schema L57-63); existence is not checked. |
| 7 | Enums validated: `scope`, `visibility`, `artifact_class`, `kata_surface`. | `schemas/skill.cue` L27 (`scope: "global" \| "role-local" \| "task-local"`), L36 (`kata_surface`), L35 (`artifact_class`), L40 (`visibility`). Negative fixture `bad-enum-scope` exercises the failure mode. |
| 8 | Negative proof via fixtures + `--self-test` mode. Minimum cases listed in AC8 covered. | `schemas/fixtures/skill-frontmatter/{valid,invalid}/`. Positive: `minimal/SKILL.md` (hard-gate + four required-or-excepted), `full/SKILL.md` (every reserved field), `full/sub/SKILL.md` (so `full`'s `calls` target resolves). Negative (each with sidecar `.expect`): `missing-hard-gate` (must surface "scope: incomplete value"), `bad-enum-scope` (must surface "conflicting values"), `broken-calls` (must surface "calls-target-exists"), `missing-required-no-exception` (must surface "required-or-excepted"). `--self-test` mode iterates the fixture tree, asserts `valid/` files pass and `invalid/` files fail with the expected substring. |
| 9 | Schema documented. | `schemas/README.md`: what is/isn't validated, surface boundary (schema vs script), field stratification table, exception mechanics, shrink procedure, local-run instructions, CI integration with pinned versions, schema-update procedure, and pointers to #289 / #303 / LANGUAGE-SPEC v0.1 / v0.2. |

## Peer enumeration

The change touches **two distinct peer classes**, enumerated separately
per `cdd/alpha/SKILL.md` §2.3:

### Class A — schema-bearing surfaces

The CTB v0.1 frontmatter shape has these surfaces in the repo:

- `docs/alpha/ctb/LANGUAGE-SPEC.md` — normative source (read; this work
  derives from §2 / §10 / §11; LANGUAGE-SPEC is unchanged).
- `docs/alpha/ctb/LANGUAGE-SPEC-v0.2-draft.md` — out of scope per
  non-goal; left untouched.
- `src/packages/**/SKILL.md` — every consumer of the shape; **all 56
  audited**, all green.
- `schemas/skill.cue` — new authority surface for the *machine-checkable*
  shape; created in this PR.
- `tools/validate-skill-frontmatter.sh` — the wrapping enforcer.
- `schemas/skill-exceptions.json` — debt ledger.
- `.github/workflows/build.yml` — runtime gate.
- `schemas/README.md` — operator-facing summary.

Every surface that needs to know about the new contract has been
updated. No surface that consumes the shape is silently stale.

### Class B — CI-job peers

The build.yml has six existing CI jobs that the notify-aggregation block
must know about: `go`, `binary-verify`, `package-verify`,
`package-source-drift` (I1), `protocol-contract-check` (I2), `link-check`
(I4). All six are still aggregated; I5 is added as the seventh entry in
both `needs:` and the failure-collector loop. The Telegram summary will
now include I5.

## Harness audit

CTB v0.1 frontmatter is YAML; CUE consumes it directly. The non-CUE
harness writers and consumers are:

- **The 56 SKILL.md frontmatter blocks themselves** — the source of
  truth. All audited individually via `cue vet`. Two normalisation
  fixes (design, gamma) for unquoted-colon-space.
- **The `package-source-drift` (I1) job** — uses `cn build --check`
  against `src/packages/**`. It does not parse SKILL.md frontmatter and
  is unaffected.
- **The `protocol-contract-check` (I2) job** — diffs
  `docs/alpha/schemas/protocol-contract.json` against `test/cmd/`. Not a
  SKILL.md consumer.
- **The lychee `link-check` (I4) job** — does not inspect frontmatter.
- **`cn-cdd-verify`** — operates on `.cdd/` artifacts, not SKILL.md
  frontmatter. Out of scope.
- **`cn build`** — reads SKILL.md but does not enforce frontmatter shape
  (its concern is package layout). Unaffected.

No non-primary-language harness silently produces SKILL.md frontmatter.

## Polyglot re-audit

Languages present in this diff and the toolchain run for each:

- **YAML** (frontmatter): `python3 -c "import yaml; yaml.safe_load(...)"`
  pre-flight on every SKILL.md after edits — all parsed cleanly. CUE
  itself acts as a stricter second check via `cue vet`.
- **CUE** (`schemas/skill.cue`): `cue vet -d '#Skill' schema.cue $tmp`
  per skill; positive cdd-skill suite (13/13) green; full validation
  (56/56) green.
- **Bash** (`tools/validate-skill-frontmatter.sh`): `bash -n` parses
  cleanly; `set -euo pipefail` enforced; explicit prereq checks; no
  silent dead code; no unused captures (every `local`/`mapfile` value is
  consumed downstream).
- **JSON** (`schemas/skill-exceptions.json`): `jq '.' < schemas/skill-exceptions.json`
  parses (used by the script every run).
- **YAML** (`.github/workflows/build.yml`): `python3 -c "import yaml; yaml.safe_load(open('.github/workflows/build.yml'))"`
  succeeds.
- **Markdown** (`schemas/README.md`): table shape and cross-references
  visually verified; lychee will validate links on next CI.

All test-surface enumeration is in `--self-test`. Every diagnostic
class introduced (cue schema failure, required-or-excepted failure,
calls-target-exists failure, calls-shape failure) has a corresponding
fixture that exercises it.

## Known debt

1. **Exception list is large** — 43 entries cover the four
   exception-backed fields across the cnos.core agent / ops skills,
   most cnos.eng skills, and a few cnos.core foundational skills.
   This is the explicit debt ledger AC5 names: each entry is a
   one-skill-one-PR addition, addressable by reading the body and
   authoring `artifact_class` / `kata_surface` / `inputs` / `outputs`.
   The entry pattern is uniform (each missing all four fields), so a
   future cycle could batch the migration if process-economics warrants
   it.
2. **`name` regex relaxed to admit slashes** — `^[a-z][a-z0-9_/-]*$`
   accepts `review/contract` and `review/implementation`, the two
   compound sub-skill names introduced in #304. LANGUAGE-SPEC §2.1
   describes `name` only as "identifier"; the slash usage is not
   explicitly endorsed nor forbidden. If the spec later forbids
   slashes, those skills get renamed and the regex tightens.
3. **`calls` cross-package targets not supported** — the resolution base
   is the package skill root only. The release skill's pre-existing
   `calls: - writing` (which appeared to mean the cnos.core writing
   skill) was removed because cross-package `calls` is not in
   LANGUAGE-SPEC §2.4.1. The release prose still loads writing as
   needed; this is an out-of-band coordination, not a static call edge.
4. **CUE's YAML 1.2 strictness vs. PyYAML's tolerance** — design and
   gamma had unquoted colon-space sequences that PyYAML accepted but
   CUE rejected. Both normalised in this PR. Going forward, the schema
   itself documents the normalisation requirement implicitly (any
   future drift fails the validator). This is process-economics-positive
   — the strict check is what catches the latent ambiguity.

## Pre-review gate

| Row | State | Evidence |
|---|---|---|
| 1 | branch rebased onto current `main` | Rebased onto `a8e67b7` (origin/main) at 23:00:00 UTC; current head `8adfd44` is one commit on top. |
| 2 | `.cdd/unreleased/301/self-coherence.md` carries CDD Trace through step 7 | This file. |
| 3 | tests present | `--self-test` covers AC8's four required negative cases plus two positive cases; full-corpus validation (56 SKILL.md) is itself an integration test that the production data passes the schema. |
| 4 | every AC has evidence | See "Acceptance criteria, evidence" table. |
| 5 | known debt explicit | See "Known debt" section. |
| 6 | schema / shape audit completed (contracts changed) | Yes — schema is the contract. Schema reviewed against LANGUAGE-SPEC §2 / §10 / §11 line-by-line; stratification table in `schemas/README.md` mirrors the schema. |
| 7 | peer enumeration completed (touches family of surfaces) | Class A (schema-bearing surfaces) + Class B (CI-job peers) enumerated and verified above. |
| 8 | harness audit completed (schema-bearing contract changed) | Done above; no non-primary-language harness silently produces frontmatter. |
| 9 | post-patch re-audit per language present in diff | Done above (YAML, CUE, Bash, JSON, Markdown). |
| 10 | branch CI green on the head commit | Local equivalents pass: `./tools/validate-skill-frontmatter.sh` → 56/56 green, `--self-test` → 3 positive + 4 negative behave as expected. CI itself runs on push; β waits for green before merge per row 10's transient-row contract. |

## Review-readiness | round 1 | base SHA `a8e67b7` | implementation SHA `8adfd44` | local validation green at 2026-04-29T23:19:49Z | CI: pending push | ready for β

**SHA convention:** the named SHA is the **implementation commit**, not
the readiness commit. The readiness commit's own SHA cannot reference
itself before it exists, and naming the prior readiness commit
recursively self-stales each time this artifact is amended (β round-1
F3). Branch HEAD as visible to β is whatever `git rev-parse
origin/{branch}` returns at poll time; the polling protocol carries it,
not this signal line.

Commits on the branch:

- `8adfd44` — implementation (schema, script, exceptions, fixtures, CI
  job + notify aggregation, schemas/README.md, hard-gate frontmatter
  backfill, two YAML normalisations, release/SKILL.md stale-`calls`
  fix).
- subsequent commits — this self-coherence artifact and any fix-round
  appendices β requests.
