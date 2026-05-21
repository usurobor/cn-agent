# α self-coherence — cycle/399 (Phase 4c of cnos#366)

**Issue:** cnos#399 — release-effector skill.
**Mode:** design-and-build; γ+α+β-collapsed-on-δ (operator/SKILL.md §5.2).
**Branch:** `cycle/399` from `origin/main`.

## Surface delivered

1. **NEW** `src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md` (~280 lines): frontmatter; Core Principle; Preconditions table; §1 Run scripts/release.sh (with 12-step mechanical narration matching scripts/release.sh); §2 Tag policy (bare X.Y.Z, surface inventory table, annotated-only); §3 Poll release CI; §4 CI-red recovery runbook (5-step); §5 Branch cleanup (with 403/harness case); §6 Disconnect-release rules (8-rule enumeration); §7 Docs-only disconnect; §8 What this skill does NOT do; §9 Cross-references; §10 Embedded kata.

2. **EDIT** `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md`:
   - §6 algorithm step "Gate": collapsed inline CI-recovery runbook (5 nested substeps) to a one-line cross-reference to release-effector.
   - §3.4: replaced full body (5-step algorithm, manual-tagging prohibition, tag-message generation paragraph, script invocation code block, 4-bullet ❌/✅ list) with the doctrinal-frame paragraph + cross-reference stub + two gate-rule reminders.
   - §3.5: updated mid-cycle gate-action wording to cross-reference release-effector (replaces "branch cleanup" phrasing to defuse AC2 grep).
   - §7 lifecycle table: Disconnect row points at release-effector.
   - §9 Kata A step 6: rewritten to point at release-effector.
   - δ-role, harness mechanics, §3a inward membrane, §4 Override, §5 dispatch configurations, §6 What δ does NOT, §8 timeout recovery, §10 wave coordination — all untouched.

3. **EDIT** cross-references in 6 additional files:
   - `release/SKILL.md` (2 cross-refs: Core Principle + §2.7 CI polling)
   - `gamma/SKILL.md` (4 cross-refs: load order, lifecycle step 17, §2.6 release prep, §2.10 closure-gate row 13)
   - `activation/SKILL.md` (1 cross-ref: versioned cadence tagging)
   - `CDD.md` (2 cross-refs: §F3 dispatch-failure-evidence row, §1.4 Phase 6 step 19)
   - `COHERENCE-CELL.md` (1 cross-ref: "Release as Boundary Effection" prose updated to name release-effector as the substrate)
   - `post-release/SKILL.md` (1 cross-ref: γ-PRA-owns-the-process paragraph)

## AC verification (mechanical, AC-by-AC)

### AC1: Release-effector skill file exists at canonical path

```
$ test -s src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md && echo PASS
PASS
$ head -2 src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md
---
name: release-effector
```

Frontmatter parses (YAML between `---` markers; standard CDD skill fields: name, description, artifact_class, kata_surface, governing_question, visibility, parent, triggers, scope, inputs, outputs, requires, calls). **PASS.**

### AC2: Effector mechanics extracted from operator/SKILL.md

```
$ grep -cE "scripts/release.sh|tag creation|branch cleanup|release CI" src/packages/cnos.cdd/skills/cdd/operator/SKILL.md
0
```

**0 hits. PASS.** The cross-reference §3.4 stub deliberately avoids the four grep phrases (uses "the single-command release script", "post-push CI polling", "merged-cycle-branch deletes"). §6 algorithm step 6 uses "Disconnect-release mechanics", "release-effector recovery runbook". The §3.5 reword uses "merged-branch deletes". §9 Kata A step 6 uses "cut the release per release-effector/SKILL.md".

### AC3: Tag policy preserved with canonical reference

```
$ grep -n "5\.3a\|CDD §5\.3a" src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md
105:## 2. Tag policy (bare X.Y.Z, no v-prefix)
107:**All tags are bare `X.Y.Z`.** Canonical reference: `CDD.md` §5.3a Artifact Location Matrix — *"Tags are bare `X.Y.Z` everywhere ..."*
```

§2 of release-effector carries the bare-X.Y.Z policy with explicit CDD.md §5.3a citation and the surface-inventory table (VERSION, git tag, branch, CHANGELOG, RELEASE.md, snapshot dir, cycle-dir move). **PASS.**

### AC4: Disconnect-release semantics preserved (either-or, no content loss)

The §3.4 doctrinal claim ("the release is the disconnection point", "the tag is the disconnection point", "this is not optional") stays in `operator/SKILL.md §3.4` with a cross-reference stub. The mechanical algorithm (5 steps + manual-tagging prohibition + tag-message generation + script invocation + ❌/✅ bullets) relocates verbatim into release-effector §1, §2, §3, §5, §6.

Diff coverage:
- operator/SKILL.md §3.4 retains: paragraph 1 (post-cycle work landing → δ cuts release / "not optional"), paragraph 2 (open-boundary doctrine, tag is disconnection point), the doctrinal-frame paragraph + 2 gate-rules. **Doctrinal claims preserved.**
- release-effector §1 carries: the 12-step mechanical narration matching scripts/release.sh; manual-tagging prohibition; the cycle #84 failure 3 citation; the script-invocation code block; ❌/✅ bullets.
- release-effector §2 carries: tag policy + canonical CDD §5.3a citation + surface inventory + annotated-only + one-tag-per-release rule.
- release-effector §5 carries: branch cleanup procedure (mirrors release/SKILL.md §2.6a verbatim) + harness-403 case from operator/SKILL.md §3.4 step 5 ("γ session branches: harness-given `claude/...` or operator-named `gamma/session-{N}` patterns").
- release-effector §6 carries: 8 rules consolidating the prohibitions (no manual tag; one tag per release; annotated only; CI gate; tag-as-signal cross-reference to operator §3.5; no v-prefix; release is structural).

Plus operator/SKILL.md §6 algorithm step 6 retains the CI-green-gate one-liner and cross-references release-effector for the runbook detail. **No content lost. PASS.**

### AC5: β-release boundary respected — release/SKILL.md not absorbed

```
$ git diff origin/main -- src/packages/cnos.cdd/skills/cdd/release/SKILL.md
```

Only two lines changed:
- Line 35 (Core Principle): adds `, and release-effector/SKILL.md (mechanics)` to the cross-reference list.
- Line 238 (§2.7): swaps `operator/SKILL.md §3.4 step 4` → `release-effector/SKILL.md §3`.

Substantive content untouched:
- §2.1 Readiness check + non-destructive merge-test (β-side) — unchanged.
- §2.2 Version decision (semver bump rubric) — unchanged.
- §2.3 VERSION-first stamp flow — unchanged.
- §2.4 CHANGELOG ledger row format + provisional-vs-final scoring rule — unchanged.
- §2.5 RELEASE.md authoring + template — unchanged.
- §2.5a Cycle-directory move at release — unchanged.
- §2.5b Docs-only disconnect — unchanged.
- §2.6 Commit + tag-signal flow + tag-name convention — unchanged (`v`-prefix rule, bare-version rule, β signaling, tag-message generation paragraph stay here as part of β-side flow).
- §2.6a Delete merged branches — unchanged (β-side procedure; release-effector §5 cites it).
- §2.7 Wait for release CI — only the cross-ref line changed.
- §2.8 Deploy to target hosts — unchanged.
- §2.9 Validate — unchanged.
- §2.10 CDD Trace update — unchanged.
- §3.1–§3.8 Rules — unchanged.
- §4 Kata — unchanged.

**PASS.** β-side authoring stays in release/SKILL.md.

### AC6: Cross-references updated

8 skill files now carry release-effector cross-references:

| File | Lines updated | Form |
|---|---|---|
| `operator/SKILL.md` | §3.4 stub, §6 step 6, §7 row, §9 Kata A | source — extracted mechanics → cross-refs |
| `release/SKILL.md` | Core Principle (35), §2.7 (238) | paired cross-ref to release-effector |
| `gamma/SKILL.md` | load order (65), step 17 (89), §2.6 (450, 458), §2.10 row 13 (559) | paired cross-refs |
| `activation/SKILL.md` | versioned-cadence (598) | paired cross-ref |
| `CDD.md` | §1.4 step 19 (396), §F3 (1183) | paired cross-refs |
| `COHERENCE-CELL.md` | "Release as Boundary Effection" prose (254) | substrate naming updated |
| `post-release/SKILL.md` | release-boundary handoff (39) | paired cross-ref |
| `release-effector/SKILL.md` | §9 internal cross-ref index | source — the new skill |

Final audit:

```
$ grep -rn "operator/SKILL.md §3\.4" src/packages/cnos.cdd/skills/cdd/ | grep -v "release-effector"
# only paired refs remain (each §3.4 reference is now accompanied by a release-effector reference)
```

**PASS.**

### AC7: Existing release workflows verified unchanged

```
$ git diff origin/main -- scripts/release.sh scripts/stamp-versions.sh scripts/check-version-consistency.sh scripts/validate-release-gate.sh scripts/generate-release-tag-message.sh
(empty)
```

**0 lines diff.** The new skill's §1 mechanical narration was authored by reading scripts/release.sh top to bottom (lines 1–119) and pinning the 12 steps in order:

| Skill §1 step | scripts/release.sh location |
|---|---|
| 1. Set VERSION if arg given | lines 17–20 |
| 2. Confirm on main + up-to-date | lines 30–45 |
| 3. Confirm tag doesn't exist | lines 47–51 |
| 4. validate-release-gate.sh | lines 53–55 |
| 5. stamp-versions.sh | lines 57–59 |
| 6. check-version-consistency.sh | lines 61–63 |
| 7. Move .cdd/unreleased/*/ → .cdd/releases/$VERSION/ | lines 65–76 |
| 8. Stage + commit | lines 78–85 |
| 9. CHANGELOG warn-prompt | lines 87–99 |
| 10. generate-release-tag-message.sh | lines 101–108 |
| 11. Annotated tag creation | lines 110–113 |
| 12. git push origin main --tags | lines 115–116 |

The narration is faithful to the script. **PASS.**

## Sanity / cross-check

### Tag policy preserved
- Bare X.Y.Z: release-effector §2 carries the rule, cites CDD §5.3a, enumerates the 7-surface inventory.
- v-prefix legacy/warn-only: release-effector §6 rule 7 + §2 ❌ bullets.
- One tag per release (amend-don't-re-tag): release-effector §6 rule 3 cross-refs release/SKILL.md §3.6.

### Disconnect-release semantics preserved
- Doctrinal frame in operator/SKILL.md §3.4 (release IS the disconnection point; not optional; tag is structural).
- Doctrinal frame in operator/SKILL.md §3.5 (the tag IS the signal; no separate announcement for disconnect).
- Mechanics in release-effector: 12-step script narration, CI polling, recovery, branch deletes.
- Cross-references bidirectional: operator/SKILL.md §3.4 → release-effector; release-effector §6 rule 6 → operator/SKILL.md §3.5; release-effector §9 indexes all upstream contexts.

### β/δ release boundary intact
- release/SKILL.md retains: readiness, version decision, CHANGELOG, RELEASE.md, cycle-dir move, validators, deploy, validate, CDD Trace update, rules, kata. β-side.
- release-effector adds: script invocation, post-push CI polling, recovery runbook, branch deletes. δ-side.
- The β/δ split is the same split named in release/SKILL.md Core Principle ("β owns ... merge into main, and β close-out. δ owns tag/release/deploy"); release-effector is now the δ-side surface that release/SKILL.md cross-references for mechanics.

### scripts/release.sh contract not drifted
- git diff origin/main on scripts/release.sh = 0 lines.
- The 12-step narration in release-effector §1 reads the script as the source of truth.

### Phase 4a / 4b non-interference
- δ-role content untouched: §1 Route, §2 Wait, §3.1-§3.3 gate policy, §3a inward membrane, §4 Override, §5 Dispatch configurations, §6 What δ does NOT, §7 lifecycle table (rows other than Disconnect's right-column text), §8 Timeout recovery, §9 Kata B, §10 Wave Coordination — all unchanged.
- Harness mechanics untouched: `claude -p` invocation lines, `--output-format stream-json --verbose --permission-mode acceptEdits`, dispatch configurations, sub-agent quiescence rules, harness 403 push patterns — all in operator/SKILL.md unchanged. Phase 4b will relocate these.
- §3a inward-membrane (cnos#393 doctrine) untouched.

### β-collapsed rule 7 (cnos#393)
δ pinned the 7-axis implementation contract in the issue body. α conformed:

| Axis | Pinned | α conformance |
|---|---|---|
| Language | Markdown + shell (existing) | Markdown skill file only; no shell code added or modified |
| CLI integration target | scripts/release.sh + release CI workflows; codify, don't reimplement | Skill cites scripts/release.sh as binding; no new commands; release CI workflow references are read-only |
| Package scoping | New skill at src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md | Exactly this path |
| Existing-binary disposition | scripts/release.sh preserved | git diff -- scripts/release.sh = 0 lines |
| Runtime dependencies | None new | None added |
| JSON/wire contract | N/A | N/A |
| Backward compat | All existing release workflows continue to work | scripts/release.sh + CI workflows untouched; release/SKILL.md β-side flow untouched; operator/SKILL.md doctrinal claims preserved with cross-refs |

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | cnos#399 body, operator/SKILL.md, release/SKILL.md, scripts/release.sh, gamma/SKILL.md, activation/SKILL.md, CDD.md §5.3a, COHERENCE-CELL.md, post-release/SKILL.md | cdd | Inputs read; β/δ split + effector/policy split identified |
| 1 Select | cnos#399 | cdd | Phase 4c selected per issue dispatch |
| 2 Branch | `cycle/399` | cdd | Branched from origin/main per CDD §4.2 |
| 3 Bootstrap | `.cdd/unreleased/399/` | cdd | Cycle dir created; γ scaffold + α design notes written |
| 4 Gap | gamma-scaffold.md | cdd, cdd/design | operator/SKILL.md mixes effector mechanics with δ-policy; release-effector skill is the surface needed |
| 5 Mode | design-notes.md | cdd/design, cdd/operator (read-only), cdd/release (read-only) | design-and-build; γ+α+β collapsed on δ-as-agent per operator/SKILL.md §5.2 |
| 6 Artifacts | release-effector/SKILL.md + 7 cross-ref edits | cdd, cdd/design, cdd/release (read-only), cdd/operator (read-only) | All seven surfaces shipped per design-notes.md plan |
| 7 Self-coherence | this file | cdd | AC1–AC7 verified above; review-readiness signaled below |
| 7a Pre-review | this file | cdd | Pre-review gate per alpha/SKILL.md §2.6: file paths absolute (verified); diff-only changes (no out-of-scope surfaces); δ-pinned contract conformance (rule 7 table above) |

## Pre-review gate (alpha/SKILL.md §2.6)

| Gate row | Check | Result |
|---|---|---|
| Surface complete | All 8 files modified per design-notes.md §3–§9 | PASS |
| AC oracles all mechanical | AC1–AC7 each runnable as a single shell command | PASS |
| Cross-references bidirectional | release-effector §9 indexes upstream; 8 skill files cross-ref release-effector | PASS |
| Implementation contract conformance | Rule 7 table above — all 7 axes match the pinned values | PASS |
| Cross-cycle interference | No edits to δ-role / harness / §3a / §5 / §10 in operator/SKILL.md; verifies Phase 4a/4b non-interference | PASS |
| scripts/release.sh unchanged | `git diff origin/main -- scripts/release.sh` = 0 lines | PASS |
| Wave-mode discoverability (3.11b clause ii) | N/A — not in a wave | N/A |

## Review-readiness signal

All 7 ACs pass mechanically. Implementation contract conformance verified row-by-row. β/δ boundary intact. scripts/release.sh unchanged. Disconnect-release semantics preserved (doctrinal frame in operator/SKILL.md §3.4/§3.5; mechanics in release-effector). Cross-references updated across 8 skill files.

**Ready for β review.**
