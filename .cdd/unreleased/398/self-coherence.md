---
manifest:
  planned: [Gap, Mode, ACs, CDD Trace, Implementation contract conformance, Review readiness]
  completed: [Gap, Mode, ACs, CDD Trace, Implementation contract conformance, Review readiness]
---

# Self-coherence — cycle #398 (Phase 4b of #366)

## Gap

Operator/SKILL.md hosts both δ-role doctrine (gate authority, override, route discipline) and harness substrate (dispatch invocation shells, observability flags, worktree management, identity discipline, polling loops, timeout recovery procedures). #366 Phase 4 names a δ split that requires three sibling carvings: 4a δ-role boundary → `delta/SKILL.md`, **4b harness substrate** (this cycle), 4c release effector. Without 4b, the harness substrate has no canonical home; cycles #371 #373 #384 surfaced the doctrine but pinned it into operator/SKILL.md prose rather than a substrate skill, conflating the WHY (δ gate authority) with the HOW (dispatch shell flags). γ-skill currently cites operator/SKILL.md for dispatch mechanics; β-skill Pre-merge gate Row 1 has no canonical doctrine to cross-reference for the preventive `--worktree` rule. Phase 5 (γ shrink) cannot land until harness mechanics are out of γ-adjacent skills.

Grep-evidence that the gap is real:

- `rg "claude -p"` in `operator/SKILL.md` (pre-cycle): 16 hits — pure mechanics literals embedded in δ-role prose.
- `rg "extensions.worktreeConfig|--worktree"` in `cdd/` (pre-cycle): hits only in `beta/SKILL.md` Row 1 (β-side enforcement) and `alpha/SKILL.md` §2.6 row 14 (α-side pre-review-gate); no canonical doctrine surface; the preventive rule prescribed by cnos#373 was nowhere stated.
- `rg "harness/SKILL.md"` in `cdd/` (pre-cycle): 0 hits — no harness skill existed.

## Mode

design-and-build; γ+α+β-collapsed-on-δ (single-session via Agent tool / parent session, per `operator/SKILL.md` §5.2). Issue body pinned the implementation contract; this cycle codifies the contract surface without reimplementing the `cn dispatch` Go code (which is preserved as-is).

## ACs

| AC | Surface | Oracle | Result |
|---|---|---|---|
| AC1 | `src/packages/cnos.cdd/skills/cdd/harness/SKILL.md` exists with standard frontmatter | `test -f` + `head -1` (`---`) + `bash tools/validate-skill-frontmatter.sh` (cue schema accepts the file) | **PASS** — file exists; frontmatter validates; 0 new findings on this cycle's surfaces |
| AC2 | Observability contract codified (#371) | `rg "output-format stream-json\|stream-json" harness/SKILL.md` ≥ 1 | **PASS** — 10 hits; §2 codifies the `--output-format stream-json --verbose` requirement, JSONL line shape, log path convention, and #369/#370 falsification anchor |
| AC3 | Worktree identity write codified (#373); cross-references β Row 1 | `rg "extensions.worktreeConfig\|--worktree" harness/SKILL.md` ≥ 1 + grep for β-row-1 cross-ref | **PASS** — 32 hits; §3.2 names the `extensions.worktreeConfig=true` check + preventive `--worktree` rule; §3.4 cross-refs `beta/SKILL.md` Pre-merge gate Row 1; β-skill Row 1 now cross-refs harness §3.2 + §3.3 |
| AC4 | Parallel α worktree pre-creation codified (#384) | `rg "parallel\|pre-create\|worktree add" harness/SKILL.md` ≥ 1 | **PASS** — 23 hits; §4 ("Parallel dispatch precondition") names the binding rule, pre-creation contract, cleanup contract, cph#27/#28 empirical anchor, and worked example |
| AC5 | operator/SKILL.md mechanics extracted | `rg "claude -p\|cn dispatch\|polling loop\|git fetch.*polling\|worktree\|stream-json" operator/SKILL.md` → 0 hits OR only cross-refs to harness/SKILL.md | **PASS** — 3 remaining hits at lines 68, 116, 389; all are explicit cross-references to harness/SKILL.md (§3 identity, §5 polling, §6 timeout recovery). The §5.2.1 Agent-tool isolation reference was rephrased to avoid the literal "worktree" token (which referred to Agent-tool filesystem-isolation mode, semantically distinct from git worktrees but caught by the over-broad oracle regex) |
| AC6 | γ-skill + β-skill cross-references | `rg "harness/SKILL.md" gamma/SKILL.md beta/SKILL.md` ≥ 1 each | **PASS** — γ: 3 hits (load order, §2.5 Identity-rotation primitive, frontmatter `calls:`); β: 1 hit (Row 1 worktree-aware identity-write rule + failure-mode catalogue cross-ref) |
| AC7 | #371/#373/#384 absorbed via close-out comments | post-merge `gh issue close` with comment naming #398 | **deferred to post-merge** — close-out comments will be filed against #371, #373, #384 after cycle/398 merges to main; comments name #398 as the absorbing cycle |

## CDD Trace

- **Selection** — Issue cnos#398 (Phase 4b of #366); decisive clause: #366 Phase 4 deliverable list explicitly absorbs #371/#373/#384 as inputs. Next-MCA committed by #393 close-out.
- **Issue quality** — γ-side issue carries full 7-axis implementation contract (pinned by δ); ACs are numbered, mechanical, independently testable; non-goals named (no reimplementation of `cn dispatch`, no movement of 4a δ-boundary content, no movement of 4c release-effector content).
- **γ scaffold** — `.cdd/unreleased/398/gamma-scaffold.md` committed at `df0e96bf` before α dispatch; satisfies γ-side rule 3.11b dual.
- **α design notes** — `.cdd/unreleased/398/design-notes.md` records partition decisions (which sections move to harness, which stay in operator), the §5.2.1 Agent-tool isolation rephrasing decision, and the AC5 oracle pattern's treatment of cross-references vs literal mechanics.
- **α build** — committed at `2966d133`: 5 files changed (1 new harness skill, 3 cross-ref edits, 1 design notes); 789 insertions / 134 deletions; operator/SKILL.md net reduction of ~127 lines (extraction confirms operator is now WHY-only for the touched surfaces).
- **β-collapsed review** — see `beta-review.md`.

## Implementation contract conformance (β Rule 7 surface)

| Axis | Pinned value | Diff conforms? |
|---|---|---|
| Language | Markdown (skill files). No Go touches. | **Yes** — diff is 100% Markdown across `harness/SKILL.md` (new), `operator/SKILL.md` (edits), `gamma/SKILL.md` (edits), `beta/SKILL.md` (edits), `design-notes.md` (new). Zero `.go` files touched. |
| CLI integration target | `cn dispatch` Go subcommand (existing under `src/go/internal/dispatch/`) — codified, not reimplemented. | **Yes** — `src/go/internal/dispatch/` is untouched (`git diff --stat` confirms); harness/SKILL.md §1.3 names `cn dispatch` and explicitly says "codified, not reimplemented." |
| Package scoping | New skill at `src/packages/cnos.cdd/skills/cdd/harness/SKILL.md`. | **Yes** — exactly the path pinned. Sibling to `operator/`, `gamma/`, `beta/`, `alpha/`, etc. |
| Existing-binary disposition | `cn dispatch` Go code preserved. | **Yes** — no Go file in the diff. |
| Runtime dependencies | None new. | **Yes** — no `go.mod` / `package.json` / dependency-manifest changes. |
| JSON/wire contract | `--output-format stream-json` JSONL contract codified as δ-membrane primitive; wire format documented, not changed. | **Yes** — harness §2.1 names the flag; §2.2 names the JSONL line schema (referring to Claude Code SDK docs for the authoritative reference); harness does not invent a new schema. |
| Backward compat | `cn dispatch` invocations work as before; harness contract is documented, not changed. | **Yes** — Go code untouched; γ-skill `Identity-rotation primitive` line preserves the existing `cn dispatch --role α|β --branch cycle/N` invocation form verbatim and adds the harness cross-reference; operator/SKILL.md §Algorithm step list preserves the role-dispatch sequence verbatim, replacing only the literal `claude -p` mechanics with `via the harness` references. |

## #371 / #373 / #384 absorption trace

| Issue | Surface in harness | Lifecycle |
|---|---|---|
| #371 (observability) | §2 "Dispatch observability contract" — flag combination, JSONL log path convention, line shape, tail-and-gate forms, #369/#370 falsification anchor; §2.5 mirror for γ-side dispatch | Codified; close on cycle merge |
| #373 (worktree identity) | §3 "Git identity for role actors" — canonical form (relocated from operator), §3.2 worktree-aware identity-write binding rule (preventive `--worktree`), §3.3 failure-mode catalogue (#301 O8, #370 β R1, #370 α F4), §3.4 cross-references to β Row 1 / α §2.6 row 14 / release §2.1, §3.5 reactive recovery | Codified; close on cycle merge |
| #384 (parallel α worktrees) | §4 "Parallel dispatch precondition" — same-repo race rule, cph#27/#28 anchor, pre-creation contract, cleanup contract, β Row 3 family cross-ref, §4.6 worked example | Codified; close on cycle merge |

## Empirical anchors carried forward

The harness skill names every empirical anchor the absorbed issues cite, so future role sessions reading the substrate skill trace each rule to its falsification:

- #371: cycle #369 + #370 parallel dispatch in default text mode (2026-05-17) — black-box dispatch produced zero in-process visibility (harness §2.4).
- #373: cycle #301 O8 (first surfacing); #370 β R1 §2.1 row 1 (second surfacing — merge-test teardown reverted shared layer); #370 α F4 (third surfacing — α session-start identity overwritten by γ identity write) (harness §3.3).
- #384: cph cdr-refactor wave 2026-05-18 to 2026-05-20 F1 (cph#27 + cph#28 — α-28 worktree-add race under shared `/root/cph`) (harness §4.2).

## Review readiness

R1 review-readiness signal: the harness skill is authored, operator/gamma/beta are updated, all mechanical ACs pass per the table above (AC7 is the post-merge close-out comment work). β-collapsed review begins.

## Cross-reference to siblings (Phase 4a / 4c coordination)

- **Phase 4a (cycle/397 — δ-role split into `delta/SKILL.md`):** this cycle does not touch §3a "δ as inward membrane" (left untouched per scope) or other Phase-4a-owned δ-role doctrine surfaces. When 4a lands, the `operator/SKILL.md` → `delta/SKILL.md` rename inherits the harness cross-references intact; no further harness-side edits required.
- **Phase 4c (cycle/399 — release effector):** this cycle does not touch §3.4 "Cut the release — disconnect the triad's final state" or §3.5 "The tag is the signal" (left untouched per scope). When 4c lands, the release-effector relocation can pull from `operator/SKILL.md` §3.4 + §3.5 + `scripts/release.sh` references; the harness substrate is unaffected.

**Conflict resolution discipline if merge-time conflicts arise with cycle/397:** integrate — each cycle touches different sections (4b: dispatch/observability/worktree/identity/polling/timeout; 4a: δ-role boundary into a new file). The `calls: harness/SKILL.md` line in operator/SKILL.md frontmatter is the only structural edit; 4a's wholesale-move of operator → delta will carry that line forward.
