<!-- sections: [issue, mode, surfaces, peer-enumeration, scope-boundary, ac-oracle, design-source-citations, diff-scope, dispatch-config, tier3-skills] -->
<!-- completed: [issue, mode, surfaces, peer-enumeration, scope-boundary, ac-oracle, design-source-citations, diff-scope, dispatch-config, tier3-skills] -->

# γ Scaffold — Cycle #394

**Date:** 2026-05-21
**Issue:** [#394](https://github.com/usurobor/cnos/issues/394) — Sub 2 (#376): bootstrap cnos.cdr package skeleton (cn.package.json + README.md + SKILL.md)
**Parent:** [#376](https://github.com/usurobor/cnos/issues/376) (CDR v0.1 master; Sub 2 wraps Sub 1's CDR.md with package metadata)
**Branch:** `cycle/394`
**Base SHA:** `e531dba0` (origin/main HEAD at session start; ε for #392)
**γ identity:** gamma / gamma@cdd.cnos
**Dispatch config:** γ+α+β-collapsed-on-δ (single Claude Code session; per breadth-2026-05-12 wave manifest precedent + cycles 375/377/378/388/390 empirical validation). β-α-collapse acknowledged: this is a docs-and-metadata mechanical mirroring cycle whose AC oracles are mechanical (`jq -e`, `test -f`, `rg -c`, schema parity via `jq 'keys'`, `cn build --check`).

---

## Issue

**Gap:** [cnos#390](https://github.com/usurobor/cnos/issues/390) (Sub 1) shipped `src/packages/cnos.cdr/skills/cdr/CDR.md` — the canonical instantiation contract. That file lives inside a directory tree that is **not yet a real cnos package**: `src/packages/cnos.cdr/` lacks `cn.package.json` (so `cn build --check` rejects it implicitly — it is not on the package-discovery list), lacks a `README.md` (so the package has no human-facing entrypoint), and lacks `skills/cdr/SKILL.md` (so the skill-loader pattern cannot dispatch CDR.md as a loaded skill). Without these three files, cnos.cdr is a directory of orphan markdown — neither loadable by `cn` nor reachable by the skill-loader.

**Goal:** Author the three package-skeleton files such that:
1. `cn.package.json` declares the package per `cn.package.v1` schema with `name: cnos.cdr`, `version: 0.1.0`, `kind: package`, `engines.cnos` declared.
2. `README.md` overviews the package (what CDR is, package structure, status, cross-references to Sub 1's CDR.md and the forthcoming Sub 3 / Sub 4 surfaces).
3. `skills/cdr/SKILL.md` declares the loader pattern with standard frontmatter (`name`, `description`, `artifact_class`, `governing_question`, `triggers`, `scope`, `inputs`, `outputs`, `requires`, `calls`) — naming CDR.md + the forthcoming Sub 3 role overlays as `calls:` entries.

**Priority:** P2 — gates cnos.cdr being a loadable package.

**Work-shape:** small docs-and-metadata cycle. Three new files in `src/packages/cnos.cdr/` plus cycle evidence under `.cdd/unreleased/394/`. No code, no schemas, no existing files modified.

**Mode:** design-and-build (γ+α+β collapsed). The design half is mechanical mirroring of cnos.cdd's package shape (cn.package.json schema, README.md structure, SKILL.md loader pattern). The build half authors the three files.

---

## Surfaces γ expects α to touch

1. `src/packages/cnos.cdr/cn.package.json` — **new** (~10 lines; cn.package.v1 schema mirror of cnos.cdd's package.json, omitting `commands` because no commands are declared by Sub 2).
2. `src/packages/cnos.cdr/README.md` — **new** (~40-60 lines; mirror of cnos.cdd/README.md adapted to CDR's research-discipline framing and pointing to CDR.md as primary doctrine).
3. `src/packages/cnos.cdr/skills/cdr/SKILL.md` — **new** (~80-120 lines; mirror of cnos.cdd/skills/cdd/SKILL.md adapted to CDR's role overlays from Sub 3; declared as `calls:` but not authored).
4. `.cdd/unreleased/394/gamma-scaffold.md` — this file.
5. `.cdd/unreleased/394/design-notes.md` — design-half record (mirror-map from cnos.cdd; CDR-specific deltas; loader pattern).
6. `.cdd/unreleased/394/self-coherence.md` — α self-coherence per AC1–AC6.
7. `.cdd/unreleased/394/beta-review.md` — β-collapsed review (mechanical AC re-check).
8. `.cdd/unreleased/394/{alpha,beta,gamma}-closeout.md` — close-out artifacts.
9. `.cdd/unreleased/394/cdd-iteration.md` — closure-gate per `post-release/SKILL.md §5.6b`.
10. `.cdd/iterations/INDEX.md` — append cycle 394 row.

**Out-of-scope surfaces (explicitly):**
- `src/packages/cnos.cdr/skills/cdr/CDR.md` — Sub 1's deliverable; cited not modified.
- `src/packages/cnos.cdr/skills/cdr/{alpha,beta,gamma,delta,epsilon}/SKILL.md` — Sub 3.
- Empirical-anchor mapping doc — Sub 4.
- New cnos package content classes (`pkg.ContentClasses` is fixed per `pkg.go`).
- Runtime dependencies in `cn.package.json` (the schema does not carry a dependencies field; package interrelation is via skill-loader cross-reference, not manifest declaration — verified by inspecting cnos.cdd's package.json which omits any deps field).

---

## Peer enumeration (§2.2a — evidence)

- "cnos.cdd's package.json shape is `{schema, name, version, kind, engines, commands}`" — confirmed: `jq 'keys' src/packages/cnos.cdd/cn.package.json` returns `["commands", "engines", "kind", "name", "schema", "version"]`.
- "cnos.cdd does not declare runtime dependencies in its package.json" — confirmed: `jq 'keys' src/packages/cnos.cdd/cn.package.json` shows no `dependencies` or `deps` field; same for cnos.core, cnos.eng, cnos.kata, cnos.cdd.kata.
- "package content classes are discovered by filesystem presence, not declared in manifest" — confirmed by `src/go/internal/pkg/pkg.go` lines 137-139: `ContentClasses = {"doctrine", "mindsets", "skills", "extensions", "commands", "orchestrators", "katas"}`; comment block at lines 132-135 reads "Filesystem presence — not manifest JSON fields — determines which classes a package contains (PACKAGE-SYSTEM.md §3: 'Content classes are discovered by directory presence')".
- "cnos.cdd/skills/cdd/SKILL.md is the exemplar loader" — confirmed: file exists; 108 lines; frontmatter declares `name`, `description`, `artifact_class`, `kata_surface`, `governing_question`, `visibility`, `triggers`, `scope`, `inputs`, `outputs`, `requires`, `calls`; body declares Load order + Rule + Role skills + Lifecycle sub-skills + Conflict rule.
- "Sub 1 already shipped `src/packages/cnos.cdr/skills/cdr/CDR.md`" — confirmed: 617 lines; sections preamble + Architecture choice + Persona/Protocol/Project + Six-field contract + Empirical anchor + Related documents + Non-goals.
- "`src/packages/cnos.cdr/cn.package.json` does not exist" — confirmed: `ls src/packages/cnos.cdr/` at scaffold time returned only `skills/` (after this scaffold finishes, ls returns README.md + cn.package.json + skills/).
- "`cn build --check` is the package-discovery command" — confirmed: built `cn` from `src/go/cmd/cn/`; `cn build --check` lists all packages including cnos.cdr after the package.json was authored; output: `✓ cnos.cdr: valid` + `✓ All packages valid.`

No claim of "X does not exist" asserted without ls/grep evidence.

---

## Scope boundary

**In scope (Sub 2):** authoring the three package-skeleton files (cn.package.json + README.md + skills/cdr/SKILL.md). Mirroring cnos.cdd's package shape mechanically. Cross-referencing Sub 1's CDR.md as primary doctrine. Naming Sub 3's forthcoming role overlays as `calls:` entries without authoring them.

**Out of scope:**
- Modifying CDR.md (Sub 1's deliverable; cite, do not modify).
- Authoring role-overlay skills (Sub 3 of #376).
- Authoring empirical-anchor mapping (Sub 4 of #376).
- Adding new package content classes (the `pkg.ContentClasses` list is fixed).
- Introducing runtime dependencies (the cn.package.v1 schema does not carry a deps field; package interrelation is by skill-loader cross-reference).
- Adding commands (the `commands` field is omitted for v0.1; no CLI surfaces introduced).
- Authoring `cph` project bindings (project layer is out of CDR-protocol-layer scope per CDR.md §"Persona, Protocol, Project").

**Sister-cycle awareness:**
- cnos#376 — parent master. Sub 2 close-out comment confirms Sub 2 shipped, AC1-AC6 met, Sub 3/4 unblocked.
- cnos#390 — Sub 1; CDR.md cited as primary doctrine; not modified.
- cnos#388 — architectural inheritance; cited in README.md cross-reference list.

---

## AC oracle approach

α verifies each AC using the oracles embedded in the issue body. β re-runs the same mechanical checks against the branch HEAD.

**AC1 — `cn.package.json` present with correct shape:**
- `test -f src/packages/cnos.cdr/cn.package.json` exits 0.
- `jq -e '.name == "cnos.cdr"' src/packages/cnos.cdr/cn.package.json` exits 0.
- `jq -e '.schema == "cn.package.v1"' src/packages/cnos.cdr/cn.package.json` exits 0.
- `jq -e '.version == "0.1.0"' src/packages/cnos.cdr/cn.package.json` exits 0.
- `jq -e '.kind == "package"' src/packages/cnos.cdr/cn.package.json` exits 0.
- `jq -e '.engines.cnos' src/packages/cnos.cdr/cn.package.json` exits 0.
- Schema parity: `jq 'keys' src/packages/cnos.cdr/cn.package.json` ⊆ `jq 'keys' src/packages/cnos.cdd/cn.package.json` (cdr omits `commands`, which is allowed — cnos.eng and cnos.cdd.kata also omit it).

**AC2 — `README.md` present and overviews the package:**
- `test -f src/packages/cnos.cdr/README.md` exits 0.
- `rg -c "CDR.md|cdr/CDR.md" src/packages/cnos.cdr/README.md` returns ≥1.
- Section presence: heading "# CDR" or "# CDR — Coherence-Driven Research"; package structure section; cross-reference to CDR.md.

**AC3 — `skills/cdr/SKILL.md` loader present:**
- `test -f src/packages/cnos.cdr/skills/cdr/SKILL.md` exits 0.
- Frontmatter parses (YAML fences `---` ... `---` at file top).
- Required frontmatter fields present: `name`, `description`, `artifact_class`, `governing_question`, `triggers`, `scope`.
- Loader pattern: file body includes "Load order" section; lists CDR.md + Sub 3 role overlays in `calls:` (with note that they are forthcoming).

**AC4 — Package is loadable / discoverable:**
- `cn build --check` (from a built kernel binary) lists cnos.cdr as a valid package.
- Output line `✓ cnos.cdr: valid` present.

**AC5 — Sub 1's CDR.md cross-references inherited:**
- `rg "skills/cdr/CDR.md|CDR\.md" src/packages/cnos.cdr/{README.md,skills/cdr/SKILL.md}` returns ≥2 hits (one per file).

**AC6 — No surface fusion:**
- `rg -i "matter type|review oracle|gate verdict|falsifiability|reproduction|claim_status" src/packages/cnos.cdr/cn.package.json` returns 0 (package.json carries no doctrinal content).
- `rg -i "matter type|review oracle|gate verdict|falsifiability" src/packages/cnos.cdr/README.md` returns 0 in normative sections (cited as "see CDR.md §Field N" is allowed; embedded doctrinal prose is not).
- `rg -i "empirical anchor|cph.*mapping|six-field declaration" src/packages/cnos.cdr/README.md` returns 0 (cnos.cdr/README.md references the empirical anchor but does not embed mapping content — Sub 4 owns mapping prose).
- SKILL.md loader names sub-skills but does not define them: `calls:` list is structural, not procedural; no role-procedure prose in SKILL.md body.

---

## Design-source citations

The cycle's design half cites:

1. `src/packages/cnos.cdd/cn.package.json` — exemplar manifest shape. Mirror schema; omit `commands` (no commands declared for cnos.cdr v0.1).
2. `src/packages/cnos.cdd/README.md` — exemplar README. Mirror structure (What X Does → Package Structure → Quick Start → License); adapt to CDR's research-discipline framing.
3. `src/packages/cnos.cdd/skills/cdd/SKILL.md` — exemplar loader. Mirror frontmatter shape + Load order + Rule + Role skills + Conflict rule; adapt `calls:` list to CDR's α/β/γ/δ/ε (note: cdd's loader lists α/β/γ + lifecycle sub-skills; cdr's loader lists α/β/γ/δ/ε per `ROLES.md §1` five-role grammar, all marked as Sub 3 forthcoming).
4. `src/packages/cnos.cdr/skills/cdr/CDR.md` — Sub 1's deliverable; the doctrinal anchor. README.md + SKILL.md both cross-reference it.
5. `src/go/internal/pkg/pkg.go` lines 107-167 (FullPackageManifest, ContentClasses, EnginesJSON, PackageCommandJSON) — the schema source-of-truth for what fields are required, what content classes exist, and how filesystem presence governs class membership.

---

## Expected diff scope

| Surface | Expected delta |
|---|---|
| `src/packages/cnos.cdr/cn.package.json` | new (~10 lines) |
| `src/packages/cnos.cdr/README.md` | new (~45-65 lines) |
| `src/packages/cnos.cdr/skills/cdr/SKILL.md` | new (~80-120 lines) |
| `.cdd/unreleased/394/gamma-scaffold.md` | new (this file; ~250 lines) |
| `.cdd/unreleased/394/design-notes.md` | new (~150-200 lines) |
| `.cdd/unreleased/394/self-coherence.md` | new (~100-150 lines) |
| `.cdd/unreleased/394/beta-review.md` | new (~80-120 lines) |
| `.cdd/unreleased/394/alpha-closeout.md` | new (~60-80 lines) |
| `.cdd/unreleased/394/beta-closeout.md` | new (~60-80 lines) |
| `.cdd/unreleased/394/gamma-closeout.md` | new (~100-130 lines) |
| `.cdd/unreleased/394/cdd-iteration.md` | new (~30-60 lines) |
| `.cdd/iterations/INDEX.md` | +1 row |

Total: ~1000-1300 line net change across ~12 files. Docs-and-metadata-only. No code, no schemas modified.

---

## Dispatch configuration

**γ+α+β collapsed on δ.** Single Claude Code session. The β-α-collapse is acknowledged: α ≠ β within a session is structurally compromised but acceptable for **docs-and-metadata mechanical-mirroring class** because β oracles are mechanical:
- File-presence (`test -f`).
- JSON schema validation (`jq -e '.field == "value"'`).
- Schema parity (`jq 'keys'` set comparison).
- Cross-reference presence (`rg -c "CDR.md"`).
- Frontmatter parse (YAML fences + required field check).
- Package discoverability (`cn build --check` exit 0 with `✓ cnos.cdr: valid`).

No subjective judgment β could provide beyond α's mechanical verification.

Pre-flight check result:
```
γ pre-dispatch check — 2026-05-21:
  cycle/394 branch: created from origin/main HEAD e531dba0
  .cdd/unreleased/394/gamma-scaffold.md exists locally: YES (this file)
  issue #394 is open: YES
  branch pre-flight: PASS
  peer enumeration: PASS — all referenced surfaces confirmed by jq/ls/rg
  scope boundary: documented (Sub 2 only; Sub 3/4 deferred; Sub 1 CDR.md cited not modified)
  issue quality gate: PASS (six ACs with mechanical oracles)
  dispatch config: γ+α+β-collapsed on δ; docs-and-metadata; mechanical β gates
  refusal conditions: none triggered at scaffold time
```

---

## Tier 3 skills

Named explicitly:
- `src/packages/cnos.core/skills/skill/SKILL.md` — skill-frontmatter discipline; loader-pattern reference.
- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` — design-half discipline (mirror-map authoring).
- `src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md` — AC oracle authoring; β re-checks the issue ACs.
- `src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` §5.6b — `cdd-iteration.md` cadence; closure-gate.

Tier 1 (read-only doctrinal):
- `src/packages/cnos.cdd/{cn.package.json, README.md, skills/cdd/SKILL.md}` — exemplar package shape.
- `src/packages/cnos.cdr/skills/cdr/CDR.md` — Sub 1's deliverable; primary doctrine anchor.
- `src/go/internal/pkg/pkg.go` — `FullPackageManifest` + `ContentClasses` schema source-of-truth.
