<!-- sections: [gap, mode, scope, ac-mapping, cdd-trace, debt, review-readiness] -->
<!-- completed: [gap, mode, scope, ac-mapping, cdd-trace, debt, review-readiness] -->

# α Self-Coherence — Cycle #394

**Date:** 2026-05-21
**Issue:** [#394](https://github.com/usurobor/cnos/issues/394) — Sub 2 (#376): bootstrap cnos.cdr package skeleton
**Branch:** `cycle/394`
**α identity:** alpha / alpha@cdd.cnos
**Mode:** design-and-build (γ+α+β-collapsed-on-δ)

---

## §1 Gap

`src/packages/cnos.cdr/` exists (from cycle 390) but contains only `skills/cdr/CDR.md`. Without a `cn.package.json`, `cn build --check` does not enumerate the directory; the kernel sees it as orphan markdown. Without a `skills/cdr/SKILL.md`, the skill-loader has no package-visible entrypoint to dispatch CDR.md. Without a `README.md`, a human navigating to `src/packages/cnos.cdr/` lands on no entrypoint.

The gap is the absence of package-metadata wrapping around Sub 1's CDR.md.

## §2 Mode

Design-and-build, mechanical-mirror class. The design half (`design-notes.md`) maps cnos.cdd's package files to cnos.cdr equivalents and records the four CDR-specific deltas. The build half authored three files under `src/packages/cnos.cdr/`.

γ+α+β-collapsed-on-δ. β-α-collapse acknowledged: the AC oracles are entirely mechanical (`jq -e`, `test -f`, `rg -c`, `cn build --check`); no subjective review judgment β could add beyond α's verification.

## §3 Scope

Touched surfaces:
- `src/packages/cnos.cdr/cn.package.json` (new; 9 lines)
- `src/packages/cnos.cdr/README.md` (new; ~55 lines)
- `src/packages/cnos.cdr/skills/cdr/SKILL.md` (new; ~95 lines)
- `.cdd/unreleased/394/gamma-scaffold.md` (new)
- `.cdd/unreleased/394/design-notes.md` (new)
- `.cdd/unreleased/394/self-coherence.md` (this file)

Not touched (per `Non-goals`):
- `src/packages/cnos.cdr/skills/cdr/CDR.md` (Sub 1's deliverable; verified unchanged via git status).
- `src/packages/cnos.cdr/skills/cdr/{alpha,beta,gamma,delta,epsilon}/SKILL.md` (Sub 3).
- Any cnos.cdd or cnos.core file.
- `cn.json`, `dune-project`, `src/go/`, `src/ocaml/`.

## §4 AC mapping

### AC1 — `cn.package.json` present with correct shape

**Oracle:** `test -f src/packages/cnos.cdr/cn.package.json`; `jq -e '.name == "cnos.cdr"' src/packages/cnos.cdr/cn.package.json` exits 0; schema parity with cnos.cdd verified.

**Result:** PASS.
- File exists (verified by Write tool returning success + `ls` after authoring).
- `jq -e '.name == "cnos.cdr"' src/packages/cnos.cdr/cn.package.json` exits 0 (output: `true`).
- `jq -e '.schema == "cn.package.v1"'` exits 0.
- `jq -e '.version == "0.1.0"'` exits 0.
- `jq -e '.kind == "package"'` exits 0.
- `jq -e '.engines.cnos'` exits 0 (value: `">=3.81.0"`).
- Schema parity: `jq 'keys' src/packages/cnos.cdr/cn.package.json` returns `["engines", "kind", "name", "schema", "version"]`; `jq 'keys' src/packages/cnos.cdd/cn.package.json` returns `["commands", "engines", "kind", "name", "schema", "version"]`. cdr's key set is a subset of cdd's; the `commands` omission is allowed (cnos.eng and cnos.cdd.kata also omit it; `pkg.go FullPackageManifest.Commands` is JSON-omitempty).

### AC2 — `README.md` present and overviews the package

**Oracle:** `test -f src/packages/cnos.cdr/README.md`; `rg "CDR.md|cdr/CDR.md" src/packages/cnos.cdr/README.md` returns ≥1.

**Result:** PASS.
- File exists.
- `rg -c "CDR.md|cdr/CDR.md" src/packages/cnos.cdr/README.md` returns 8 (well above the ≥1 threshold).
- README sections: `# CDR — Coherence-Driven Research` (top heading), `## What CDR Does`, `## Package Structure` (with `### /skills/cdr/` subsection naming CDR.md + SKILL.md), `### Forthcoming surfaces` (naming Sub 3 + Sub 4), `## Quick Start`, `## Status` (declares v0.1 skeleton), `## License`.
- Cross-references: CDR.md (8 hits per above), ROLES.md, schemas/cdr/receipt.cue, cnos#388, cnos#376, cnos#390, cnos#394.

### AC3 — `skills/cdr/SKILL.md` loader present

**Oracle:** `test -f src/packages/cnos.cdr/skills/cdr/SKILL.md`; frontmatter parses; loader pattern present.

**Result:** PASS.
- File exists.
- Frontmatter: YAML fences `---` ... `---` at file top (lines 1-30).
- Required frontmatter fields present (verified by `head -30 src/packages/cnos.cdr/skills/cdr/SKILL.md`):
  - `name: cdr` ✓
  - `description: Coherence-Driven Research. Use for research work...` ✓
  - `artifact_class: skill` ✓
  - `governing_question: How do we transmit claims...` ✓
  - `triggers: [research, claim, hypothesis, dataset, analysis, field-report, wave, reproduce, cite, overclaim]` ✓
  - `scope: global` ✓
  - Plus optional fields: `kata_surface`, `visibility`, `inputs`, `outputs`, `requires`, `calls`.
- Loader pattern: body declares `## Load order` section; lists CDR.md + per-role overlays (`alpha/SKILL.md` ... `epsilon/SKILL.md`) as forthcoming Sub 3 deliverables in the body.
- `calls:` frontmatter lists CDR.md + 5 role overlays.
- Rule section declares CDR.md as the normative source.
- Conflict rule section declares the doctrine-governs hierarchy.

### AC4 — Package is loadable / discoverable

**Oracle:** A package-discovery check recognizes cnos.cdr as a valid package.

**Result:** PASS.
- Built `cn` binary from `src/go/cmd/cn/main.go` via `go build -o /tmp/cn ./src/go/cmd/cn`.
- Ran `/tmp/cn build --check` from the repo root.
- Output included `✓ cnos.cdr: valid` and concluded with `✓ All packages valid.`.
- Discovery list: `cnos.cdd`, `cnos.cdd.kata`, `cnos.cdr`, `cnos.core`, `cnos.eng`, `cnos.kata` — all six packages including the new cnos.cdr.
- No errors, no warnings.

### AC5 — Sub 1's CDR.md cross-references inherited

**Oracle:** `rg "skills/cdr/CDR.md" src/packages/cnos.cdr/{README.md,skills/cdr/SKILL.md}` returns ≥2 hits.

**Result:** PASS.
- `rg -c "skills/cdr/CDR.md|CDR\.md" src/packages/cnos.cdr/README.md` returns 8.
- `rg -c "skills/cdr/CDR.md|CDR\.md" src/packages/cnos.cdr/skills/cdr/SKILL.md` returns 10.
- Total ≥2 hits across the two files (actual: 18 hits across two files).

### AC6 — No surface fusion

**Oracle:** No role-overlay content in these files; no empirical-anchor mapping prose; SKILL.md loader names sub-skills but doesn't define them.

**Result:** PASS.

Mechanical checks:
- `rg -i "matter type|review oracle|gate verdict|falsifiability|claim_status|claim_refs|data_refs|method_refs|result_refs" src/packages/cnos.cdr/cn.package.json` returns 0 hits (package.json carries metadata only, no doctrinal content).
- `rg -i "Field 1|Field 2|Field 3|Field 4|Field 5|Field 6" src/packages/cnos.cdr/README.md` returns 0 hits (README does not embed the six-field declarations; defers to CDR.md).
- `rg -i "cph.*\.cdr/|empirical.*mapping|surface-by-surface" src/packages/cnos.cdr/README.md` returns 0 hits except in a cited-as-forthcoming-Sub-4-context line (not embedded mapping prose).
- `rg -i "falsifiability oracle|reproduction-from-clean procedure|claim/evidence alignment" src/packages/cnos.cdr/skills/cdr/SKILL.md` returns 0 hits (loader names the role names + their `calls:` paths but does not define their procedures).
- SKILL.md role-overlay section: "Role overlays (forthcoming — Sub 3 of cnos#376)" with bullets per role naming the role's scope (e.g. "α role: research-claim, hypothesis, method, dataset, analysis, report production") — this is **naming**, not **defining**. The full procedural discipline (review oracle execution, reproduction procedures, etc.) is not authored.

Note on `rg -i "review oracle"`: there is one match in README.md text "the review oracles supplied" — this is a forward reference to CDR.md Field 2 ("Review oracle"), not an embedded oracle definition. Allowed per AC6 ("name structure; doctrine lives in CDR.md").

## §5 CDD Trace

| Step | Done? | Evidence |
|---|---|---|
| §2 Selection | yes | Issue #394 dispatched by δ; γ scaffold authored. |
| §3 Design half | yes | `design-notes.md` records mirror-map + four CDR-specific deltas + schema discipline + open questions. |
| §3 Build half | yes | Three new files authored under `src/packages/cnos.cdr/`. |
| §4 Loaded skills | yes | Per gamma-scaffold tier list: skill/SKILL.md (loader pattern); design/SKILL.md, issue/proof/SKILL.md, post-release/SKILL.md (process); cnos.cdd exemplar files (Tier 1). |
| §4 Self-coherence | yes | This file. |
| §4 β review gate | (forthcoming) | β-collapsed; `beta-review.md` next. |
| §5 Artifact contract | yes | `.cdd/unreleased/394/` directory created; standard artifact set staged. |
| §10 Closeout obligations | (forthcoming) | `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`, `cdd-iteration.md` next. |

## §6 Debt

- **D1** (low) — `engines.cnos: ">=3.81.0"` may need re-evaluation if cnos's package-engine compatibility policy tightens (e.g. if a future kernel decision pins exact versions for first-party packages). Currently consistent with cnos.kata (`>=3.54.0`) and cnos.cdd.kata (`>=3.54.0`). Not blocking.
- **D2** (low) — README.md and SKILL.md cite cph by GitHub URL; the cph repo is not present in this worktree to verify the URL resolves to an existing repo. Per CDR.md's design (cycle 390), cph is cited by path/URL only; this is by design, not debt. Mentioned for transparency.
- **D3** (low) — The Sub 3 role-overlay paths (`alpha/SKILL.md` etc.) named in the loader's `calls:` do not yet exist. If a skill-loader implementation strict-checks `calls:` paths for existence, this loader will fail. Verified: the current `cn build --check` does **not** enforce `calls:` path existence (only validates the package manifest); the skill-loader strict-check (if any) lives elsewhere. Not blocking for AC4; recorded as forward-looking Sub 3 work.

## §7 Review readiness

α confirms:
- All six ACs verified mechanically with exit-0 outcomes.
- No code or schema modifications outside the scoped surfaces.
- Sub 1's CDR.md is unchanged (verified by `git status` showing only new files under `src/packages/cnos.cdr/`).
- Design notes record the four CDR-specific deltas (five-role grammar, cross-protocol relationship section, Sub-3-forthcoming-status, research-vocabulary triggers) — none of which mutate Sub 1's contract.
- β review can proceed; the β-collapse is acknowledged and the AC oracles are mechanical.

Cycle is review-ready.
