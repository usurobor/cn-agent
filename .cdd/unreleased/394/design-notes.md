<!-- sections: [intent, mirror-map, cdr-specific-deltas, loader-pattern, schema-discipline, open-questions] -->
<!-- completed: [intent, mirror-map, cdr-specific-deltas, loader-pattern, schema-discipline, open-questions] -->

# α Design Notes — Cycle #394

**Date:** 2026-05-21
**Issue:** [#394](https://github.com/usurobor/cnos/issues/394)
**Mode:** design-and-build (mechanical mirroring of cnos.cdd's package shape; γ+α+β-collapsed-on-δ)

---

## §1 Intent

Sub 1 ([cnos#390](https://github.com/usurobor/cnos/issues/390)) shipped the doctrinal contract for CDR (the file `src/packages/cnos.cdr/skills/cdr/CDR.md`) but the file lives in a directory that is **not yet a real cnos package**. `cn build --check` (the discovery oracle) looks for a `cn.package.json` at the package root — without it, the directory is invisible to the kernel. The skill-loader cannot dispatch CDR.md as a skill without a `SKILL.md` at the package-visible entrypoint per `cnos.core/skills/skill/SKILL.md` convention. The package has no human-facing entrypoint without a `README.md`.

Sub 2's job is to author the three package-skeleton files such that:
1. `cn` recognises cnos.cdr as a valid package.
2. The skill-loader can dispatch through `skills/cdr/SKILL.md` to load CDR.md and (when Sub 3 ships) the per-role overlays.
3. A human reading `src/packages/cnos.cdr/README.md` learns what CDR is, where the doctrine lives, and what is forthcoming.

The shape is **mechanical mirroring** of cnos.cdd's package: cnos.cdd is the reference package for the c-d-X family (engineering-discipline twin); cnos.cdr is its research-discipline twin in the option (a) architectural split.

---

## §2 Mirror map

The three new files map to three cnos.cdd files. The mapping is structural; the content adapts per discipline.

| Source (cnos.cdd) | Target (cnos.cdr) | Mirror discipline |
|---|---|---|
| `cn.package.json` | `cn.package.json` | Schema-identical (`cn.package.v1`); fields `{schema, name, version, kind, engines}` mirrored; `commands` **omitted** (no commands declared for cnos.cdr v0.1); `name` changes from `cnos.cdd` to `cnos.cdr`; `version` set to `0.1.0` (fresh package start, not aligned with cnos's `3.81.0` global because cnos.cdr is a new package at v0.1 — same convention as `cnos.cdd.kata: 0.3.0` and `cnos.kata: 0.2.0`); `engines.cnos` set to `>=3.81.0` (compatible-with-current; same convention as cnos.kata's `>=3.54.0`). |
| `README.md` | `README.md` | Structure mirrored: top heading + What X Does + Package Structure + (Quick Start / Forthcoming / Status) + License. Discipline content adapts: CDD's loss function is engineering (artifact improvement under repairable feedback); CDR's is research (truth-preserving claim transmission under uncertainty). README cross-references CDR.md as primary doctrine; lists Sub 3 (role overlays) and Sub 4 (empirical-anchor) as forthcoming surfaces. |
| `skills/cdd/SKILL.md` | `skills/cdr/SKILL.md` | Frontmatter shape mirrored (`name`, `description`, `artifact_class`, `kata_surface`, `governing_question`, `visibility`, `triggers`, `scope`, `inputs`, `outputs`, `requires`, `calls`). Body structure mirrored: Load order + Rule + Role overlays + Cross-protocol relationship + Conflict rule. Content adapts: cdd's loader names CDD.md + cdd-side α/β/γ/operator + lifecycle sub-skills; cdr's loader names CDR.md + cdr-side α/β/γ/δ/ε (note: cdr uses the full five-role grammar per `ROLES.md §1`; cdd's older loader pre-dates δ/ε formalization but cdr does not inherit that gap — cdr's loader names all five). Sub 3 forthcoming-status is explicit in the loader body. |

### Why omit `commands` from cn.package.json

The `commands` field in `cn.package.v1` is optional (per `pkg.go FullPackageManifest`: `Commands map[string]PackageCommandJSON` is JSON-omitempty by default). cnos.eng and cnos.cdd.kata omit it; cnos.cdd, cnos.core, cnos.kata declare it. cnos.cdr v0.1 declares no CLI commands (the issue's "CLI integration target: None" line); omitting the field is correct per schema and consistent with the `Non-goals` of the issue body ("Do NOT add new cnos package content classes" / no new commands).

### Why `version: 0.1.0` (not `3.81.0`)

The cnos repo carries a global version (`3.81.0` at base SHA) declared in `cn.json`. Packages can either align with the global (cnos.core, cnos.cdd, cnos.eng all read `3.81.0`) or carry their own (`cnos.cdd.kata: 0.3.0`, `cnos.kata: 0.2.0`). For a **new** package at v0.1 skeleton, the convention from cnos.cdd.kata/cnos.kata is `0.X.0`. cnos.cdr's first version is `0.1.0` — declaring it as `3.81.0` would falsely imply parity with the mature engineering-side overlay. The package's release-train independence is signalled by the standalone version. The future release-cycle decision (whether cnos.cdr eventually aligns with the global cnos version or stays standalone) is a Sub 3/Sub 4 question, not Sub 2.

### Why `engines.cnos: ">=3.81.0"` (range, not pin)

cnos.cdd pins to `3.81.0` exactly because cnos.cdd is the canonical engineering overlay and its in-sync version is part of the global release contract. cnos.cdr is a new package at v0.1 — pinning to `3.81.0` exactly would mean the package becomes unloadable the moment cnos releases `3.82.0`. The `>=3.81.0` range convention (used by cnos.kata `>=3.54.0` and cnos.cdd.kata `>=3.54.0`) matches the "new package, accepts current-or-newer engine" intent. This is a deliberate choice, not a copy-mistake from the exemplar.

---

## §3 CDR-specific deltas (vs cnos.cdd's package shape)

The mirroring is structural, not content-level. Four CDR-specific deltas:

### Delta 1 — Five-role grammar, not three

cnos.cdd's `skills/cdd/SKILL.md` `calls:` lists `alpha/SKILL.md`, `beta/SKILL.md`, `gamma/SKILL.md`, plus lifecycle sub-skills (`issue/`, `design/`, `plan/`, `review/`, `release/`, `post-release/`). cnos.cdr's `skills/cdr/SKILL.md` `calls:` lists `alpha/SKILL.md`, `beta/SKILL.md`, `gamma/SKILL.md`, `delta/SKILL.md`, `epsilon/SKILL.md` — the full five-role grammar per `ROLES.md §1` and CDR.md Field 5 (which explicitly names ε's role).

The cdd loader's pre-five-role shape is a legacy of cdd predating the ε formalization (`ROLES.md §1` introduced ε after cdd's loader was already authored). cnos.cdr does not inherit that gap; it ships with the full ladder from day one.

cdd's lifecycle sub-skills (`issue/`, `design/`, `plan/`, `review/`, `release/`, `post-release/`) are **not** mirrored in cnos.cdr's loader because:
1. They are engineering-cycle-specific (issue → design → plan → review → release → post-release is the CDS-cycle shape).
2. CDR.md Field 4 (δ cadence) declares the CDR wave-shape is **gate-transition-shaped, not release-shaped** — there is no `release/SKILL.md` analog for CDR.
3. Lifecycle skills for CDR are Sub 3's authoring decision; if CDR needs wave-shape lifecycle skills (e.g. `wave-open/`, `wave-close/`), Sub 3 will declare them.

The Sub 2 loader names only the role overlays; the lifecycle layer is Sub 3's question.

### Delta 2 — Cross-protocol relationship section

cnos.cdd's loader does not need a "this protocol is one of a family" section because cdd is the original protocol; there is no peer. cnos.cdr's loader **does** need such a section because:
1. The architectural-choice inheritance ([cnos#388](https://github.com/usurobor/cnos/issues/388)) makes the family structure load-bearing.
2. A reader hitting `skills/cdr/SKILL.md` needs to know that loading cdr is a class-routing act — the matter is research-class, not engineering-class.

The Cross-protocol relationship section in cnos.cdr's SKILL.md names cnos.cdd (the engineering twin), declares the kernel/per-protocol split, cross-references CDR.md §"Architecture choice" + cnos#388, and warns against misreading CDR as CDS-with-different-words.

### Delta 3 — Sub 3 forthcoming-status is explicit

cnos.cdd's loader names role skills as if they exist (`alpha/SKILL.md` etc. are real files). cnos.cdr's loader names them in `calls:` but the loader body has an explicit "Role overlays (forthcoming — Sub 3 of cnos#376)" section explaining they are not yet authored. Until Sub 3 ships, the loader runs against CDR.md alone with role discipline supplied by the persona hub (`cn-rho`) and the project binding (`cph`).

This is honest about the v0.1 skeleton status — the package is loadable and discoverable, but the role-overlay layer is in flight.

### Delta 4 — Triggers are research-vocabulary

cnos.cdd's loader triggers: `review`, `release`, `issue`, `design`, `plan`, `assess`, `post-release`, `ship`, `tag`. These are CDS-cycle verbs.

cnos.cdr's loader triggers: `research`, `claim`, `hypothesis`, `dataset`, `analysis`, `field-report`, `wave`, `reproduce`, `cite`, `overclaim`. These are CDR-discipline verbs. The trigger set is the class-routing signal — a dispatch matching `overclaim` should not load `cdd`.

---

## §4 Loader pattern

The loader pattern is a documented convention from `cnos.core/skills/skill/SKILL.md` (the meta-skill describing skills). The relevant invariants:

- **Single package-visible entrypoint.** External dispatch enters via the `cdr` skill only; sub-skills (`cdr/alpha`, etc.) are advisory, not public entrypoints.
- **Load order is declared.** The loader names the canonical doctrine (`CDR.md`) and the role overlay for the active role; everything else is Tier 2/3.
- **Conflict rule is named.** If the loader and the doctrine disagree, doctrine governs. If a role overlay and the doctrine disagree on the contract shape, doctrine governs.
- **`requires` is mechanical.** "CDR applies" + "CDR.md exists in this directory" — the loader does not run for non-research matter.

The cnos.cdr loader instantiates this pattern with one CDR-specific addition: the **class-routing** check at load time. The Cross-protocol relationship section explicitly warns that loading cdr is a routing act — if the matter is engineering-class, the loader should route to `cdd` instead.

---

## §5 Schema discipline

`cn.package.v1` schema fields and their treatment in cnos.cdr's manifest:

| Field | Required? | Value | Source-of-truth |
|---|---|---|---|
| `schema` | yes | `"cn.package.v1"` | `pkg.go FullPackageManifest.Schema` |
| `name` | yes | `"cnos.cdr"` | `pkg.go FullPackageManifest.Name`; `ValidatePackageManifestData` rejects empty |
| `version` | yes | `"0.1.0"` | `pkg.go FullPackageManifest.Version`; semver convention for new packages |
| `kind` | yes | `"package"` | `pkg.go FullPackageManifest.Kind`; all installed packages are `kind: package` |
| `engines` | yes | `{"cnos": ">=3.81.0"}` | `pkg.go EnginesJSON.Cnos`; range vs pin per §2 above |
| `commands` | optional | omitted | `pkg.go FullPackageManifest.Commands`; omitempty default; cnos.eng/cnos.cdd.kata also omit |

The schema is `cn.package.v1`. There is no `dependencies` field — package interrelation is by skill-loader cross-reference (e.g. `skills/cdr/SKILL.md` `calls:` may reference paths in other packages), not by manifest declaration. Verified by:
- `jq 'keys' src/packages/cnos.cdd/cn.package.json` returns `["commands", "engines", "kind", "name", "schema", "version"]` — no `dependencies`.
- `pkg.go FullPackageManifest` has no `Dependencies` field.
- `pkg.go ManifestDep` exists but is for the **hub's** `deps.json` (the operator's desired package list), not for inter-package declarations within a single package's manifest.

This means cnos.cdr does **not** declare a dependency on cnos.core or cnos.cdd in its package.json. The relationship is doctrinal (CDR.md §"Architecture choice" inherits cnos.cdd's kernel by reference) and skill-level (the loader cross-references CDR.md and Sub 3's overlays), not manifest-level.

---

## §6 Open questions

### OQ1 — When does cnos.cdr's version graduate from 0.X to 1.0?

The v0.1 skeleton ships only the doctrinal contract (Sub 1) + package metadata (Sub 2). Sub 3 (role overlays) and Sub 4 (empirical-anchor doc) are forthcoming. A natural version-bump cadence:
- `0.1.0` — Sub 2 ships (this cycle).
- `0.2.0` — Sub 3 ships (role overlays land).
- `0.3.0` — Sub 4 ships (empirical-anchor doc lands).
- `1.0.0` — cnos.cdr declared mature; verdict-enum schema (cycle 390 F1) and project-policy template (cycle 390 F2) pinned; first CDR cycle ships against the package and produces a valid `#CDRReceipt`.

Not in Sub 2 scope; recorded for Sub 3/4 dispatch context.

### OQ2 — Should cn.package.json declare `description`?

cnos.cdd's manifest has no `description` field; `pkg.go FullPackageManifest` does not declare one. The schema does not carry a description; human-facing description lives in `README.md`. cnos.cdr follows the same convention — no `description` field.

### OQ3 — Should the Sub 3 role overlays be staged with empty `SKILL.md` stubs in this cycle?

Considered and rejected. Stubs would:
1. Make Sub 3 a "fill stubs" cycle rather than an authoring cycle, which trivializes the role-discipline content.
2. Risk drift between the stub frontmatter and Sub 3's eventual frontmatter (each role overlay needs role-specific frontmatter that Sub 3 will author).
3. Violate the issue's `Non-goals`: "Do NOT author role overlays (Sub 3)."

The loader's `calls:` list names them by path (`alpha/SKILL.md` etc.) and the loader body declares them as forthcoming. Sub 3 will create the files at the listed paths.

### OQ4 — Should README.md cite cph by full path or just by name?

The repo cph (`usurobor/cph`) is not present in this worktree; it is an external repo (per cycle 390's design-notes and CDR.md §"Empirical anchor"). cnos.cdr/README.md cites cph by GitHub URL (`https://github.com/usurobor/cph`) in the Forthcoming surfaces paragraph; it does not embed any cph-local content. This matches CDR.md's discipline: "These paths are cited; cph-local prose is **not** embedded in this file. Embedding would violate the protocol/project separation."
