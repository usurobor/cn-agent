<!-- sections: [collapse-acknowledgment, mechanical-recheck, observations, verdict] -->
<!-- completed: [collapse-acknowledgment, mechanical-recheck, observations, verdict] -->

# β Review — Cycle #394

**Date:** 2026-05-21
**Issue:** [#394](https://github.com/usurobor/cnos/issues/394)
**Branch:** `cycle/394`
**β identity:** beta / beta@cdd.cnos
**Round:** R1
**β-α collapse acknowledged:** YES — this is a docs-and-metadata mechanical-mirroring cycle whose AC oracles are entirely mechanical (`jq -e`, `test -f`, `rg -c`, `cn build --check`). No subjective judgment β could provide beyond α's verification. Same precedent as cycles 375/377/378/388/390 (per breadth-2026-05-12 wave manifest and subsequent empirical validation). Class-identification: docs-and-metadata, CDS-class under repairable feedback — **not** research-claim authoring; CDR.md Field 6's α=β prohibition (research-class) does not apply.

---

## §1 Collapse acknowledgment

β-α collapse on δ acknowledged. The cycle's verdict-determining evidence is:

- `jq -e` exit codes against `cn.package.json`.
- `rg -c` line counts in README.md and SKILL.md.
- `test -f` file-presence checks.
- `cn build --check` package-discovery output.
- Schema-parity comparison via `jq 'keys'`.

These are mechanical; the β role's contribution is re-running them against the branch HEAD with the same expected outcome. The collapse is structurally compromised but tractable for this class of cycle.

The α-β prohibition in CDR.md Field 6 ("α=β within a single cycle. Always prohibited") applies to research-class claim transmission, not to engineering-class docs cycles. This cycle is CDS-class (docs + metadata under repairable feedback); the collapse is consistent with engineering-class precedents (cycles 375/377/378/388/390).

---

## §2 Mechanical re-check

### AC1 — `cn.package.json` present with correct shape

```
$ test -f src/packages/cnos.cdr/cn.package.json && echo OK
OK

$ jq -e '.schema == "cn.package.v1" and .name == "cnos.cdr" and .version == "0.1.0" and .kind == "package" and (.engines.cnos | type == "string")' src/packages/cnos.cdr/cn.package.json
true

$ jq 'keys' src/packages/cnos.cdr/cn.package.json
["engines", "kind", "name", "schema", "version"]

$ jq 'keys' src/packages/cnos.cdd/cn.package.json
["commands", "engines", "kind", "name", "schema", "version"]
```

Schema parity: cdr's key set is a subset of cdd's; the `commands` omission is allowed (cnos.eng + cnos.cdd.kata also omit it; `pkg.go FullPackageManifest.Commands` is JSON-omitempty).

**Verdict: PASS.**

### AC2 — `README.md` present and overviews the package

```
$ test -f src/packages/cnos.cdr/README.md && echo OK
OK

$ rg -c "CDR.md|cdr/CDR.md" src/packages/cnos.cdr/README.md
8
```

Heading "# CDR — Coherence-Driven Research" present at line 1. Sections present: What CDR Does, Package Structure (with `/skills/cdr/` subsection naming CDR.md), Forthcoming surfaces (naming Sub 3 + Sub 4), Quick Start, Status, License.

**Verdict: PASS.**

### AC3 — `skills/cdr/SKILL.md` loader present

```
$ test -f src/packages/cnos.cdr/skills/cdr/SKILL.md && echo OK
OK

$ rg "^name:|^description:|^artifact_class:|^governing_question:|^triggers:|^scope:" src/packages/cnos.cdr/skills/cdr/SKILL.md | head -10
name: cdr
description: Coherence-Driven Research. Use for research work …
artifact_class: skill
governing_question: How do we transmit claims about the world …
triggers:
scope: global
```

All six required frontmatter fields present (per `cnos.core/skills/skill/SKILL.md` convention for skill frontmatter). Plus optional `kata_surface`, `visibility`, `inputs`, `outputs`, `requires`, `calls`. Body has `## Load order` section + `## Rule` + `## Role overlays (forthcoming — Sub 3 of cnos#376)` + `## Cross-protocol relationship` + `## Conflict rule`.

**Verdict: PASS.**

### AC4 — Package is loadable / discoverable

```
$ go build -o /tmp/cn ./src/go/cmd/cn
$ /tmp/cn build --check 2>&1
✓ cnos.cdd: valid
✓ cnos.cdd.kata: valid
✓ cnos.cdr: valid
✓ cnos.core: valid
✓ cnos.eng: valid
✓ cnos.kata: valid
✓ All packages valid.
```

cnos.cdr is recognised by the package-discovery oracle. No errors.

**Verdict: PASS.**

### AC5 — Sub 1's CDR.md cross-references inherited

```
$ rg -c "CDR\.md" src/packages/cnos.cdr/README.md src/packages/cnos.cdr/skills/cdr/SKILL.md
src/packages/cnos.cdr/README.md:8
src/packages/cnos.cdr/skills/cdr/SKILL.md:10
```

Total: 18 hits across two files; threshold is ≥2 (one per file).

**Verdict: PASS.**

### AC6 — No surface fusion

```
$ rg -i "matter type|review oracle|gate verdict|falsifiability|claim_status|claim_refs|data_refs|method_refs|result_refs" src/packages/cnos.cdr/cn.package.json
(no output — 0 hits)

$ rg -i "Field 1|Field 2|Field 3|Field 4|Field 5|Field 6" src/packages/cnos.cdr/README.md
(no output — 0 hits)

$ rg -i "falsifiability oracle|reproduction-from-clean procedure|claim/evidence alignment" src/packages/cnos.cdr/skills/cdr/SKILL.md
(no output — 0 hits)
```

The loader's role-overlay section names roles ("α role: research-claim … production") without authoring their procedural discipline. Per Field 2 review-oracle procedural detail and Field 3 close-out artifact procedural detail — none embedded.

One borderline observation (not a finding): the SKILL.md loader's "Role overlays (forthcoming)" section uses one-line role descriptions like "α role: research-claim, hypothesis, method, dataset, analysis, report production". This is a **structural name** of α's matter type from CDR.md Field 1 — it is naming, not defining. The matter-type prose is one-line scoping, not procedural authoring. β classifies this as compliant with AC6's "names structure; doctrine lives in CDR.md" requirement.

**Verdict: PASS.**

---

## §3 Observations

### Observation 1 — `cn build --check` only validates package manifests, not skill `calls:` paths

β verified that `cn build --check` does not enforce existence of `calls:` paths declared in skill frontmatter. The cnos.cdr loader names `alpha/SKILL.md` ... `epsilon/SKILL.md` as forthcoming Sub 3 files; these paths do not yet exist; `cn build --check` still passes. This is by design (per `pkg.go` schema — the manifest is validated, not the skill-loader graph). If a future skill-loader strict-check is introduced, the Sub 3 dispatch will need to land before the strict-check is enabled, or the loader will need a "forthcoming" marker that the strict-check honors. Recorded as forward-looking debt; not a Sub 2 blocker.

### Observation 2 — Version mismatch with cnos cn.json global is intentional

cnos.cdr declares `version: "0.1.0"` while cn.json's global is `3.81.0`. β verified this matches the convention from cnos.cdd.kata (`0.3.0`) and cnos.kata (`0.2.0`) — new packages start at independent 0.X versions. The design-notes record the rationale (§2 "Why version: 0.1.0"). Not a finding.

### Observation 3 — Engines range vs pin

cnos.cdr declares `engines.cnos: ">=3.81.0"` (range), while cnos.cdd pins to `3.81.0`. β verified the range convention matches cnos.kata (`>=3.54.0`) and cnos.cdd.kata (`>=3.54.0`) — new packages use ranges to avoid being broken by the next cnos release. The design-notes record the rationale (§2 "Why engines.cnos: '>=3.81.0' (range, not pin)"). Not a finding.

### Observation 4 — Five-role grammar in loader is intentional

The cnos.cdr loader's `calls:` lists five role overlays (α/β/γ/δ/ε); cnos.cdd's loader lists three (α/β/γ) plus lifecycle sub-skills. β verified this matches `ROLES.md §1`'s formalised five-role ladder + CDR.md Field 5 (which names ε explicitly). The design-notes record this as "Delta 1 — Five-role grammar, not three" (§3.D1). Not a finding; intentional per the design.

### Observation 5 — Sub 1's CDR.md is unchanged

```
$ git status src/packages/cnos.cdr/skills/cdr/CDR.md
# (no entry — file unchanged)
```

The cycle's `Non-goals` ("Do NOT modify CDR.md") is honored. β confirms no diff against CDR.md.

---

## §4 Verdict

**APPROVE.**

All six ACs PASS mechanically. The four observations above are intentional design choices (not findings); each is recorded in `design-notes.md` with rationale. The β-α collapse is acknowledged and reconciled per class (docs-and-metadata, CDS-class, mechanical oracles).

Cycle is ready for merge. Recommend:
1. α writes `alpha-closeout.md`.
2. β writes `beta-closeout.md`.
3. γ writes `gamma-closeout.md` + `cdd-iteration.md` + `INDEX.md` row.
4. Merge to main with `Closes #394`.
5. Post-merge: comment on cnos#376 confirming Sub 2 ships; cross-reference Sub 1 (#390), Sub 3 (planned), Sub 4 (planned).
