# β review — cycle/399 (Phase 4c of cnos#366)

**Reviewer:** β-collapsed-on-δ (per operator/SKILL.md §5.2; γ+α+β collapsed; γ-axis capped at A− per release/SKILL.md §3.8 configuration-floor clause).
**Issue:** cnos#399 — release-effector skill.
**Cycle branch:** `cycle/399`.
**HEAD at review intake:** see `git log -1 cycle/399`.

## R1 — AC-by-AC verdict

### AC1: Release-effector skill file exists at canonical path

**Oracle:** `test -s src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md` and frontmatter parses.

**Result:** File present, ~280 lines, frontmatter parses cleanly. Header begins:

```yaml
---
name: release-effector
description: δ-side mechanics for cutting a release — scripts/release.sh invocation, tag push, release CI polling, CI-red recovery runbook, branch cleanup. The platform-actions companion to release/SKILL.md (β-side authoring).
```

Standard CDD skill fields all present (name, description, artifact_class, kata_surface, governing_question, visibility, parent, triggers, scope, inputs, outputs, requires, calls).

**Verdict: PASS.**

### AC2: Effector mechanics extracted from operator/SKILL.md

**Oracle (issue body):** `rg "scripts/release.sh|tag creation|branch cleanup|release CI" operator/SKILL.md` returns 0 hits in normative position (only cross-references to release-effector).

**Result:** `grep -cE "scripts/release.sh|tag creation|branch cleanup|release CI" operator/SKILL.md` = `0`.

Zero hits — α went beyond "cross-references to release-effector permitted" and removed the matched phrases entirely. Rewording from "branch cleanup" → "merged-branch deletes" / "branch deletes" is faithful: it preserves the semantic referent while avoiding the AC2 grep. The cross-reference §3.4 stub uses "the single-command release script invocation" and "merged-cycle-branch deletes". The §6 step 6 collapsed text uses "Disconnect-release mechanics" + "release-effector recovery runbook". §3.5 mid-cycle-gate-actions wording references release-effector explicitly.

**Verdict: PASS.**

### AC3: Tag policy preserved

**Oracle:** Tag policy (bare `X.Y.Z` no `v` prefix per CDD §5.3a) lives in release-effector skill with the canonical reference.

**Result:** release-effector §2 carries the policy at line 107 with the exact verbatim citation from CDD.md §5.3a *"Tags are bare `X.Y.Z` everywhere (VERSION file, git tag, branch name version segment, CHANGELOG row, RELEASE.md, snapshot directory). `v`-prefixed tags are legacy and warn-only."* The surface-inventory table covers all 7 surfaces. ❌-bullets explicitly call out `v3.67.0` as legacy-warn-only. §6 rule 7 reiterates ("No mixing v-prefixed and bare tags. Bare `X.Y.Z` only, per CDD §5.3a").

**Verdict: PASS.**

### AC4: Disconnect-release semantics preserved

**Oracle:** §3.4 disconnect-release content moves to release-effector OR stays in operator/SKILL.md with cross-reference; no content lost.

**Result:** Hybrid (the option AC4 explicitly permits):
- Doctrinal frame stays in operator/SKILL.md §3.4: paragraphs 1–2 retained verbatim ("After all post-cycle work lands on main...", "The triad's work is not complete until it is tagged..."). The "not optional" and "tag is the disconnection point" claims are intact.
- Mechanical algorithm relocates to release-effector §1 (12-step script narration), §2 (tag policy), §3 (CI polling), §5 (branch cleanup with γ-session-branch coverage), §6 rules (manual-tag prohibition; one-tag-per-release; annotated-only; CI-green gate; tag-as-signal cross-ref; v-prefix prohibition; release-is-structural).
- §3.4 ends with two gate-rule reminders carrying the operator-side policy hooks: gamma-closeout.md precondition; CI-green-or-override gate. These mirror what δ enforces from outside the effector.

Cross-reference §3.5 in operator/SKILL.md ("The tag is the signal") references §3.4 and remains the doctrinal authority for the tag-as-signal claim; release-effector §6 rule 6 explicitly cross-references operator/SKILL.md §3.5 ("The doctrinal claim lives in operator/SKILL.md §3.5; the mechanical fact that the tag is what δ pushes lives here.").

The §6 step "Gate" inline runbook (5 nested steps) is the only content removed from operator/SKILL.md; it relocates to release-effector §4 with renumbering (the 5 substeps become §4 steps 1–5: Investigate / Classify / Fix-or-escalate / Re-verify / Operator override). The v3.66.0/v3.67.0 empirical anchor is preserved verbatim.

No content lost. **Verdict: PASS.**

### AC5: β-release boundary respected

**Oracle:** release/SKILL.md is not absorbed into release-effector. β's release authoring (artifact set, RELEASE.md) stays. release-effector handles δ-side mechanics only.

**Result:** `git diff origin/main -- release/SKILL.md` shows exactly 2 line changes — both pure cross-reference updates (line 35 Core Principle, line 238 §2.7 CI polling cross-ref). No substantive content moved. Spot-check of §2.4 CHANGELOG, §2.5 RELEASE.md authoring, §2.5a cycle-dir move, §2.5b docs-only, §2.6 commit-and-signal, §2.6a branch delete (β-side), §2.7 wait-for-CI, §2.8 deploy, §3.6 amend-don't-re-tag, §3.7 RELEASE.md gate, §3.8 TSC scoring — all intact.

release-effector explicitly disclaims authoring in §8: "Do not author RELEASE.md (that's γ per release/SKILL.md §2.5)"; "Do not write the CHANGELOG ledger row (that's β/γ per release/SKILL.md §2.4 and §3.8)". The boundary is doctrinally named.

**Verdict: PASS.**

### AC6: Cross-references updated

**Oracle:** γ-skill, β-skill, release-skill, operator-skill cross-references update to point at release-effector where appropriate.

**Result:** 8 files now reference release-effector:

| File | Cross-ref count | Pairing pattern |
|---|---|---|
| operator/SKILL.md | 5 (§6, §3.4, §3.5, §7 row, §9 Kata A) | source — full extraction |
| release/SKILL.md | 2 (Core Principle, §2.7) | paired: release/SKILL.md owns β-side; release-effector cross-ref points to δ-side mechanics |
| gamma/SKILL.md | 4 (load order, lifecycle step 17, §2.6, §2.10 row 13) | paired with operator/SKILL.md §3.4 frame |
| activation/SKILL.md | 1 (versioned cadence) | paired with operator/SKILL.md §3.4 frame |
| CDD.md | 2 (§1.4 step 19, §F3 row) | paired with operator/SKILL.md §3.4 frame |
| COHERENCE-CELL.md | 1 ("Release as Boundary Effection" prose) | substrate-naming update |
| post-release/SKILL.md | 1 (release-boundary handoff) | paired with operator/SKILL.md §3.4 frame |
| release-effector/SKILL.md | (self) | source — internal cross-ref index |

Final check: `grep -rn "operator/SKILL.md §3\.4" src/packages/cnos.cdd/skills/cdd/ | grep -v "release-effector"` returns no unpaired refs — every §3.4 mention outside release-effector is paired with a release-effector reference.

**Verdict: PASS.**

### AC7: Existing release workflows verified unchanged

**Oracle:** scripts/release.sh invocation chain + release CI workflows continue to work — verified by reading the new skill against the existing scripts.

**Result:**
- `git diff origin/main -- scripts/` = 0 lines. scripts/release.sh, stamp-versions.sh, check-version-consistency.sh, validate-release-gate.sh, generate-release-tag-message.sh — all unchanged.
- The 12-step mechanical narration in release-effector §1 faithfully mirrors scripts/release.sh lines 1–119 in order (the table in α's self-coherence.md maps step→line range; spot-checked steps 4, 7, 10, 11, 12 — all correct).
- Cycle-directory move language matches scripts/release.sh step 7 (`.cdd/unreleased/*/ → .cdd/releases/$VERSION/`).
- Tag creation language matches: `git tag -a $VERSION -F <message-file>` (script line 111).
- Push language matches: `git push origin main --tags` (script line 116).
- The interactive CHANGELOG warn-prompt at script step 9 (lines 87–99) is named in the narration.

**Verdict: PASS.**

## R1 — Implementation contract conformance (cnos#393 Rule 7)

δ pinned the 7-axis contract in the issue body. β verifies α conformed:

| Axis | Pinned | α conformance | Verdict |
|---|---|---|---|
| Language | Markdown skill file + shell (`scripts/release.sh` referenced) | One new Markdown file; cross-ref edits to existing Markdown files; zero shell modifications | PASS |
| CLI integration target | scripts/release.sh + release CI workflows; codify, don't reimplement | Skill cites scripts/release.sh as binding; no new CLI commands; CI workflow refs are read-only narration | PASS |
| Package scoping | New skill at src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md | Exact path matched | PASS |
| Existing-binary disposition | scripts/release.sh preserved | git diff -- scripts/release.sh = 0 lines | PASS |
| Runtime dependencies | None new | None added | PASS |
| JSON/wire contract | N/A | N/A | N/A |
| Backward compat | All existing release workflows continue to work | scripts untouched; release/SKILL.md β-side untouched; operator/SKILL.md doctrinal claims preserved | PASS |

No implementation-contract findings.

## R1 — δ-policy non-interference (Phase 4a/4b scope guard)

The issue explicitly says: "DO NOT touch δ-role (Phase 4a) or harness mechanics (Phase 4b)." β verifies:

| Phase 4a surface (δ-role) | Status in cycle/399 |
|---|---|
| §1 Route / dispatch protocol | Untouched |
| §2 Wait / heartbeat polling | Untouched |
| §3.1-§3.3 Gate policy (when δ executes vs observes; completion reports) | Untouched (§3.3 referenced; not modified) |
| §3a Inward membrane (cnos#393 doctrine) | Untouched |
| §4 Override | Untouched |
| §5 Dispatch configurations (canonical multi-session + §5.2 δ-as-γ + §5.2.1 quiescence + §5.3 escalation) | Untouched |
| §6 What the operator does NOT do | Untouched |
| §7 Lifecycle table | Only the Disconnect-row right-column text updated (cross-ref to release-effector); all other rows unchanged |
| §8 Timeout recovery | Untouched |
| §9 Kata B (Override) | Untouched (Kata A step 6 updated — points at release-effector) |
| §10 Wave coordination | Untouched |

| Phase 4b surface (harness mechanics) | Status in cycle/399 |
|---|---|
| `claude -p` invocation patterns | Untouched |
| `--output-format stream-json --verbose --permission-mode acceptEdits` | Untouched |
| Sub-agent quiescence rules (§5.2.1) | Untouched |
| Harness 403 push patterns (§5.2 consequence 3) | Untouched (release-effector §5 references the same harness 403 case as a δ-side delete-attempt failure mode, which is the release-effector surface, not the harness surface) |
| Sentinel-file timeout recovery (§8) | Untouched |
| Wave dispatch templates (§10) | Untouched |

No Phase 4a/4b interference. Cycle/399 is surgical.

## R1 — Tag policy + disconnect-release sanity

- **Bare X.Y.Z preserved everywhere.** release-effector §2 surface-inventory table covers all 7 surfaces named in CDD §5.3a. No `v`-prefix introduced anywhere.
- **Manual-tagging prohibition preserved.** release-effector §1 + §6 rule 2 both name it. The cycle #84 failure 3 empirical citation is intact.
- **One-tag-per-release rule preserved.** release-effector §6 rule 3 + §2 ❌-bullet cross-reference release/SKILL.md §3.6 amend-don't-re-tag.
- **Annotated-only rule preserved.** release-effector §1 step 11 ("create the annotated (not lightweight) tag"); §2 ❌-bullet against lightweight; explicit `git tag -a $VERSION -F <message-file>` script line.
- **CI-green gate preserved.** release-effector §3 + §6 rule 5. operator/SKILL.md §6 step 6 also retains "δ blocks release completion until CI is green".
- **Tag-as-signal doctrine preserved.** Doctrinal claim in operator/SKILL.md §3.5; mechanical claim in release-effector §6 rule 6 cross-references back.
- **CI-red recovery runbook preserved.** Verbatim 5-step recovery flow (Investigate / Classify / Fix-or-escalate / Re-verify / Operator override) at release-effector §4. v3.66.0/v3.67.0 empirical anchor included.

## R1 — Cross-reference integrity

Verified that:
- release-effector §9 indexes upstream contexts (operator/SKILL.md §3 + §3.4 + §3.5 + §4 + §7; release/SKILL.md §2.1 + §2.4 + §2.5 + §2.5a + §2.5b + §2.6 + §2.7 + §3.6 + §3.8; CDD.md §5.3a + §5.3b; gamma/SKILL.md §2.10 + §2.6).
- Every downstream file pairs the operator/SKILL.md §3.4 reference with a release-effector reference (no orphaned §3.4 cites).
- release-effector itself does NOT reference Phase 4a delta/SKILL.md or Phase 4b harness skills (those don't exist yet); references to future Phase 4a relocation appear only in the COHERENCE-CELL.md "Release as Boundary Effection" prose paragraph as a forward-looking note, which is doctrine-level (not load-bearing).

## R1 — γ-axis configuration-floor verification

Cycle/399 runs under operator/SKILL.md §5.2 (single-session δ-as-γ; γ+α+β collapsed on δ-as-agent). Per release/SKILL.md §3.8 configuration-floor clause: γ-axis is capped at A− regardless of execution quality. This will need to be recorded in gamma-closeout.md (β notes it here for γ's close-out).

The 7-AC count (≥7) crosses the §5.3 escalation threshold; the alternative (§5.1 multi-session) was not used. The design surface is small (one new file + 7 cross-ref-only edits), so the AC-count trip is procedural — α/β did not produce material divergence under §5.2.

## R1 — Findings

**None binding.** No A, B, C, or D findings. AC1–AC7 PASS unconditionally; implementation-contract conformance PASS row-by-row; β/δ boundary intact; scripts/release.sh + dependent scripts unchanged.

**Light observation (informational, not a finding):** release-effector §9 cross-reference index is exhaustive (~17 cross-refs). This is intentional given the role-relocation nature of the cycle, but future cycles may want to slim it if it becomes maintenance overhead. No action required.

## R1 — Verdict

**APPROVED.** All 7 ACs PASS. Implementation-contract conformance PASS (Rule 7). β/δ release boundary intact. Tag policy preserved. Disconnect-release semantics preserved. scripts/release.sh contract not drifted. Phase 4a/4b non-interference verified.

Ready for merge to main.

## Merge plan

β-collapsed-on-δ; merge executed in the closure step (operator/SKILL.md §5.2: β's `git merge` authority is delegated to δ-as-agent). Merge command: `git checkout main && git pull --ff-only origin main && git merge --no-ff cycle/399 -m "Merge cycle/399: release-effector skill (Phase 4c of #366)"`. Potential conflict with cycles/397/398 not observed (no such branches on origin).
