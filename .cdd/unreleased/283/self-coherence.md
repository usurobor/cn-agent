# α Artifact — Cycle #283

This file is the primary branch artifact for cycle #283 per CDD.md §5.4. It carries the gap, mode, active skills, impact graph, CDD Trace, AC mapping, known debt, review-readiness signal, fix-round appendices, and final α close-out.

## Gap

CDD's triadic coordination model uses GitHub Pull Requests for agent-to-agent communication: α creates a PR (`gh pr create`), all three roles poll PR state (`gh pr view`, comment threads, review state, CI checks), β posts review findings as PR comments, β merges via `gh pr merge`, and close-outs reference PR numbers throughout. PRs are a GitHub UI surface, not a repo artifact — every PR action consumes context budget on polling/subscription mechanics that produce no durable repo evidence. The `.cdd/unreleased/` directory already exists for cycle-local artifacts; the natural artifact-exchange surface is `.cdd/unreleased/{N}/{role}.md`, one cycle directory per issue keyed by issue number, one role file per role inside. Multiple `Derives from` annotations (#274, #268, #266) document repeated failures of the PR-polling model: first-iteration absorption of pre-existing PRs, branch-glob templates that assume issue-number encoding the harness doesn't produce, transient PR-state rows that drift between write and read, and shared-GitHub-identity review-state degradation.

The three failure modes the PR coordination layer creates:

1. **Polling ceremony.** Every cycle pays a PR-existence/PR-status/PR-comments/PR-CI quadruple-poll that silently degrades when any of: `gh auth` is unavailable, MCP and shell tools have unequal reachability, branch-globs miss harness-encoded names, or shared identity defeats GitHub-native review state. These are observable failures because they happened, not theoretical risk.
2. **Surface duplication.** A PR body restates the gap; a PR comment restates a finding; a release CI auto-generates from PR titles; a close-out references a PR number. Each is a copy of information already (or potentially) in a repo artifact. Drift between copies is silent.
3. **Coordination over a non-repo surface.** A PR is not a git artifact. It cannot be `git diff`ed, cannot be merged via `git merge`, cannot be moved into the release directory by a `mv`. Every coordination action on a PR is a side effect outside the repo's audit trail.

The change replaces PR-based coordination with `.cdd/unreleased/{N}/` — one cycle directory per issue, with multiple role-distinguished files inside (`self-coherence.md`, `beta-review.md`, `alpha-closeout.md`, `beta-closeout.md`, `gamma-closeout.md`) — all coordination expressed as commits to those files. Issues remain (gap-naming); branches remain (isolation); PRs are removed from the triadic protocol.

## Mode

| | |
|---|---|
| Mode | MCA |
| Cycle level (target) | L7 — system-shaping leverage. Replaces a coordination protocol with a leaner one that eliminates a recurring friction class (PR polling). |
| Tier 1 skills | `cdd/SKILL.md`, `cdd/CDD.md`, `cdd/alpha/SKILL.md` (this is CDD self-modification — no Tier 2 or Tier 3 beyond Tier 1 per the issue) |
| Tier 2 skills | not applicable per the issue's "Tier 3 skills: none beyond CDD Tier 1" declaration; Tier 2 `eng/document` would normally apply for prose discipline but is not loaded since CDD's lifecycle sub-skills (design, plan, issue) are themselves Tier 1 here |
| Tier 3 skills | none beyond CDD Tier 1 (per issue declaration) |
| Work shape | Substantial — 10 files declared in issue + the migration of existing in-version cycle artifacts |
| Affected files (modified) | `src/packages/cnos.cdd/skills/cdd/{CDD.md, SKILL.md, alpha/SKILL.md, beta/SKILL.md, gamma/SKILL.md, review/SKILL.md, release/SKILL.md, post-release/SKILL.md, operator/SKILL.md}` |
| Affected files (file moves) | `.cdd/unreleased/{alpha,beta,gamma}/{268,278}.md` → `.cdd/unreleased/{268,278}/{alpha,beta,gamma}.md` (6 renames; existing in-version cycles migrated to the new layout) |
| Affected files (declared but unchanged) | `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` — issue lists it but a full PR-ref scan returned zero hits. Leaving unchanged is correct. |
| Branch | `claude/cdd-tier-1-skills-pptds` (carries the dispatch's pre-set branch name; the cycle scope claims this branch) |
| Design artifact | not required — the issue body carries Problem, Impact, ACs, Non-goals, related artifacts, and dependencies; this α artifact carries the impact graph and CDD Trace. A separate design doc would duplicate without adding leverage. |
| Plan | not required — implementation order is `CDD.md` (canonical authority) → role skills (α/β/γ) → lifecycle sub-skills (review/release/post-release/operator) → top-level SKILL.md → in-version artifact migration → α artifact write. The order is dictated by authority hierarchy + the natural flow of the new coordination model; no novel sequencing is needed. |

## Impact graph

Per `design/SKILL.md` §2.3 (impact graph is not optional). The change is schema-bearing for the CDD coordination protocol — every consumer of "where do roles coordinate" is bound to the new model.

### Downstream consumers (modified)

| File | What changed | Why |
|---|---|---|
| `cdd/CDD.md` | §Tracking rewritten; §1.4 α/β/γ algorithms; §5.3 step 4/5/7/7a/8/9 rows; §5.3a Artifact Location Matrix; §5.4 CDD Trace template; small-change path; default flow diagram; role responsibility table; branch pre-flight; bootstrap stub note; supporting refs | Canonical authority; every other surface points here for the lifecycle and artifact contract |
| `cdd/SKILL.md` | frontmatter `triggers` (drop `PR`); `inputs` (drop `PR context`) | Loader entrypoint; trigger list visible to runtime |
| `cdd/alpha/SKILL.md` | frontmatter outputs; §2.1 dispatch intake (poll `.cdd/unreleased/{N}/`); §2.2 produce-in-order (coherence contract path); §2.5 self-coherence path; §2.6 pre-review gate (rows 2 + 10); §2.7 review-readiness signal (replaces `gh pr create` and override clause); §2.8 close-out (write `alpha-closeout.md` separately from `self-coherence.md`); §3.4 post-patch re-audit; §4 embedded kata expected artifacts | α's executable detail for steps 2–7a |
| `cdd/beta/SKILL.md` | frontmatter inputs/outputs; §1 phase map; §2.2/2.3 polling and refusal; §3 stale-references list; §5 closure (write `beta-closeout.md`) | β's role boundary; phase-link rules |
| `cdd/gamma/SKILL.md` | frontmatter inputs (drop PR state, add cycle dir); §2.5 polling shell snippets (drop PR-list, add `.cdd/unreleased/{N}/` ls-tree); §2.7 close-out triage (read `*-closeout.md` files); §2.10 closure rule (write `gamma-closeout.md`) | γ's executable detail for cycle coordination + closure |
| `cdd/review/SKILL.md` | frontmatter inputs; §2.0 PRE-GATE (drop `gh pr view`, use `git merge-base --is-ancestor`); §2.0.5/2.0.8 CDD artifact references; §2.2.8 authority-surface conflict sites; §3.7 CI gate language; §3.8 merge instruction (drop PR number); §6 output format header; §7 after-review (`git merge` action); §7.1 review identity (git author, not GitHub identity); various 2.x prose hits | β's review surface; verdict format and gate semantics |
| `cdd/release/SKILL.md` | frontmatter inputs (drop PR); §2.5 release-notes auto-gen language (commit titles, not PR titles); §2.5a `mv` command (per-cycle dir move + canonical-filename note); §2.10 CDD Trace + β close-out path | β's release surface; coordination-cycle artifact handoff into release |
| `cdd/post-release/SKILL.md` | §3 encoding-lag table key; §5.5 review quality (cycles instead of PRs; review-rounds reads `beta-review.md` rounds); §5.6 self-coherence axes; §pre-publish gate row | γ's PRA surface; cycle metrics now read role files instead of PR comments |
| `cdd/operator/SKILL.md` | §2.1 wake-up signals; §2.2 polling shell snippets (drop PR-list, add cycle-dir watch); §3.1 gate table (PR merge → branch merge to main); §3.2 anti-pattern; §6 cycle lifecycle table | δ's role surface; gate-execution against the new merge model |

### Downstream consumers (unmodified — declared but no PR refs found)

- `cdd/design/SKILL.md` — full PR-ref scan returned zero hits. Design-skill prose centers on impact-graph/AC discipline, not coordination mechanics; no edits needed.
- `cdd/issue/SKILL.md` — not in the issue's declared list, scanned defensively, zero PR refs.
- `cdd/plan/SKILL.md` — not in the issue's declared list, scanned defensively, zero PR refs.

### File migrations (in-version cycle artifacts)

| Old path | New path | Reason |
|---|---|---|
| `.cdd/unreleased/{alpha,beta,gamma}/268.md` | `.cdd/unreleased/268/{alpha,beta,gamma}-closeout.md` | Existing #268 close-outs; rename to descriptive names per the new convention. (One step: pivot `{role}/{N}.md` → `{N}/{role}.md`, then rename `{role}.md` → `{role}-closeout.md` since the existing files are close-out narratives.) |
| `.cdd/unreleased/{alpha,beta,gamma}/278.md` | `.cdd/unreleased/278/{alpha,beta,gamma}-closeout.md` | Same as #268 — existing close-outs, same migration. |
| (new) | `.cdd/unreleased/283/self-coherence.md` | This file — the cycle's primary branch artifact carrying gap, mode, ACs, CDD Trace, review-readiness, fix-rounds. |

### Out of scope (verifier kept as-is)

- `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` — has a `--pr` flag for verifying merged PR-scoped cycles. The issue scope explicitly says "All affected files are within `src/packages/cnos.cdd/skills/cdd/`," and the verifier verifies *merged* cycles (not active coordination). Leaving as-is; the legacy `CLOSE-OUT.md` paths it asserts are now warn-only per the updated §5.3a, but the verifier's own canonical-vs-warn logic doesn't need to change for that — it already accepts both paths.
- `src/packages/cnos.cdd/commands/cdd-verify/test-cn-cdd-verify.sh` — same scope.

### Authority relationships

`CDD.md` is the canonical source for the lifecycle algorithm (§1.4) and the artifact contract (§5.3, §5.3a, §5.4). Role skills add executable detail; lifecycle sub-skills add per-phase detail. When the new filename convention is named in `CDD.md` §Tracking, every other file refers to that authoritative table — no surface restates the canonical filename set.

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | issue #283 body + dispatch prompt | — | γ surfaced PR-coordination ceremony from #274/#268 derivation history; selected as next MCA |
| 1 Select | issue #283 | — | Selected gap: triadic CDD coordination uses GH PRs; replace with `.cdd/unreleased/{N}/` artifacts |
| 2 Branch | `claude/cdd-tier-1-skills-pptds` | cdd | Branch pre-set by harness/dispatch; matches §4.2 format (agent prefix + scope + suffix) |
| 3 Bootstrap | — | cdd | Not required — process/governance change with no version dir; no bundle gets a frozen snapshot in this cycle (CDD self-modification under the existing `cnos.cdd` tree) |
| 4 Gap | this file §Gap | — | Named incoherence: PR layer adds ceremony without value; `.cdd/unreleased/{N}/` exists and is the natural artifact-exchange surface |
| 5 Mode | this file §Mode | cdd, alpha (Tier 1 only per issue) | MCA, L7-shaped, Tier 3 = none beyond CDD Tier 1 |
| 6 Artifacts | CDD.md, SKILL.md, alpha/SKILL.md, beta/SKILL.md, gamma/SKILL.md, review/SKILL.md, release/SKILL.md, post-release/SKILL.md, operator/SKILL.md + 7 file migrations | cdd | All 9 affected files modified per impact graph; design/SKILL.md left unchanged (zero PR refs); existing in-version cycles migrated to new layout |
| 7 Self-coherence | this file | cdd | AC-by-AC mapping below |
| 7a Pre-review | this file §Review-readiness | cdd | Pre-review gate per `alpha/SKILL.md` §2.6 — see §Review-readiness section |

## AC Coverage

Each AC mapped to concrete evidence in the diff or migration. Evidence cites file:line where useful.

| # | AC | Met / partial / missing | Evidence |
|---|----|-------------------------|----------|
| 1 | No GH PR creation in any CDD role skill or `CDD.md` algorithm — α outputs a branch + artifacts in `.cdd/unreleased/{N}/` | **Met** | `alpha/SKILL.md` §2.7 rewritten: removed "create the PR" + override clause; replaced with "commit `.cdd/unreleased/{N}/self-coherence.md` carrying the explicit review-readiness signal." `CDD.md` §1.4 α algorithm step 7 rewritten same way. Final residual scan for `gh pr create` / `Open PR` / `request review from β` shows zero remaining instructions. |
| 2 | No GH PR polling in any role — inter-role coordination uses `.cdd/unreleased/{N}/` artifacts and `git` branch operations | **Met** | `CDD.md` §Tracking rewritten: query table dropped PR existence/status/comments rows; replaced with branch + cycle-dir state. `gamma/SKILL.md` §2.5 polling snippet uses `git ls-tree -r origin/main .cdd/unreleased/<N>/`, not `gh pr list`. `beta/SKILL.md` §1.4 step 3 polling reference shape uses cycle-dir blob SHA. `operator/SKILL.md` §2.2 supplementary polling uses `git ls-tree`, not `gh pr list`. |
| 3 | α signals review-readiness by writing a review-request artifact to `.cdd/unreleased/{N}/` | **Met** | `alpha/SKILL.md` §2.7: "commit `.cdd/unreleased/{N}/self-coherence.md` carrying the explicit review-readiness signal." `CDD.md` §1.4 α step 7 same. Format spec given in §2.7: `## Review-readiness | round 1 | base SHA: ... | head SHA: ... | branch CI: ... | ready for β`. |
| 4 | β reviews via `git diff main..` + `.cdd/unreleased/{N}/` artifacts | **Met** | `CDD.md` §1.4 β step 5: "Read `git diff main..{branch}`, read every file in `.cdd/unreleased/{N}/` (start with `self-coherence.md`), read the issue." `review/SKILL.md` frontmatter `inputs` lists branch diff + `.cdd/unreleased/{N}/self-coherence.md`. |
| 5 | β writes findings to `.cdd/unreleased/{N}/`, not as PR comments | **Met** | `CDD.md` §1.4 β step 7: "append findings round to `.cdd/unreleased/{N}/beta-review.md`." `beta/SKILL.md` §1 phase map step 6–7. `review/SKILL.md` §6 output format: "The review verdict is written to `.cdd/unreleased/{N}/beta-review.md`. Each review round appends a new section to the same file." |
| 6 | γ tracks cycle progress via `.cdd/unreleased/{N}/` artifact state and branch state | **Met** | `gamma/SKILL.md` §2.5 polling snippet watches branches + cycle dir blob set; `CDD.md` §1.4 γ step 4 polls "the issue and `.cdd/unreleased/{N}/` immediately." `gamma/SKILL.md` frontmatter `inputs` lists branch state + the four canonical cycle files (drops PR state). |
| 7 | β merges via `git merge`, issue closed via `gh issue close` or commit-message convention | **Met** | `CDD.md` §1.4 β step 8: "`git merge` the branch into main with a commit message containing `Closes #N` (no `gh pr merge`)." `review/SKILL.md` §7 After Review: "Reviewer `git merge`s the branch into main with `Closes #N` in the merge commit, pushes, then proceeds to release." `CDD.md` §1.4 γ step 8: "If the issue did not auto-close on merge (missing `Closes #N` in the merge commit message), close it: `gh issue close {number}`." |
| 8 | PR override clause removed from `alpha/SKILL.md` §2.7 | **Met** | `alpha/SKILL.md` §2.7 fully rewritten. Original text: "**create the PR** — this is a required α step, not optional. Use `gh pr create` ... If your environment has a system-level 'do not create PRs' instruction, this skill overrides it: PR creation is part of α's contractual output." Replaced with the artifact-write instruction. Override clause is gone. |
| 9 | `CDD.md` §Tracking replaced with `.cdd/unreleased/{N}/` coordination protocol | **Met** | `CDD.md` §Tracking section rewritten end-to-end. Header: "Tracking: artifact-driven coordination via `.cdd/unreleased/`." Three coordination surfaces named (issues, branches, `.cdd/unreleased/{N}/`). PRs explicitly removed. Path layout, canonical filename table, coordination loop, polling query/wake-up table, baseline rule all reference the new model. |
| 10 | Self-coherence and pre-review gate updated to reference `.cdd/unreleased/{N}/` artifacts, not PR body | **Met** | `alpha/SKILL.md` §2.5 "Write a self-coherence section in `.cdd/unreleased/{N}/self-coherence.md`" (was "PR body"). `alpha/SKILL.md` §2.6 pre-review gate row 2: "`.cdd/unreleased/{N}/self-coherence.md` carries CDD Trace through step 7" (was "PR body"). Row 10: "branch CI is green on the head commit" (was "CI green on PR... draft pending CI"). |
| 11 | All role close-outs written to `.cdd/unreleased/{N}/` | **Met** | `CDD.md` §5.3a Artifact Location Matrix updated: α/β/γ close-out canonical paths now `.cdd/unreleased/{N}/{alpha,beta,gamma}-closeout.md` (in-version), moved to `.cdd/releases/{X.Y.Z}/{N}/` at release. Legacy `.cdd/releases/{X.Y.Z}/{role}/CLOSE-OUT.md` paths marked warn-only. `alpha/SKILL.md` §2.8 writes `alpha-closeout.md`; `beta/SKILL.md` §5 writes `beta-closeout.md`; `gamma/SKILL.md` §2.10 writes `gamma-closeout.md`. |
| 12 | Pre-review gate rows for "CI green on PR" and "PR body carries trace" replaced with branch/artifact equivalents | **Met** | `alpha/SKILL.md` §2.6 row 2: "`.cdd/unreleased/{N}/self-coherence.md` carries CDD Trace through step 7" (was "PR body carries CDD Trace"). Row 10: "branch CI is green on the head commit (or, if local CI is unavailable, the artifact's review-readiness section says so explicitly and β waits for green before merge)" (was "CI is green on the head commit, or PR remains draft pending CI"). Mirror updates in `CDD.md` §5.3 step 7a row. |

## Role self-check

Did α push ambiguity onto β? Two sites worth flagging:

1. **In-version `{N}.md` filename vs post-release `CLOSE-OUT.md` filename mismatch** — pre-existing in §5.3a; not opened by this cycle, not closed by this cycle either. The release skill's `mv` preserves filenames (`.cdd/unreleased/283/alpha-closeout.md` → `.cdd/releases/X.Y.Z/283/alpha-closeout.md`), but §5.3a's "noncanonical / legacy" column names `.cdd/releases/{X.Y.Z}/alpha/CLOSE-OUT.md` as the legacy aggregate form. The verifier (`cn-cdd-verify`) currently asserts `CLOSE-OUT.md` as canonical. There's no contradiction in policy after this change — the new canonical IS the per-cycle file under `{N}/`, with `CLOSE-OUT.md` warn-only — but the verifier hasn't been updated to match. This is a consequence of the issue's explicit scope ("All affected files are within `src/packages/cnos.cdd/skills/cdd/`"), which excludes `commands/cdd-verify/`. Documented as known debt below.
2. **Existing `.cdd/releases/{3.59.0,3.60.0}/{role}/CLOSE-OUT.md` legacy artifacts not migrated** — the release directories under `.cdd/releases/` from prior releases use the old layout. These are post-release artifacts, frozen by repo policy (§5.6). Not migrated.

Every claim above is backed by file edits in the diff or by file mv operations recorded in `git status`. No claim depends on an unverified PR-side surface.

## Known debt

1. **Verifier `cn-cdd-verify` not updated.** Out of scope per the issue. With this cycle's §5.3a changes, the verifier's canonical-vs-warn logic for close-out paths is now slightly behind: it asserts `CLOSE-OUT.md` as canonical when the new policy makes per-cycle `{role}-closeout.md` canonical and `CLOSE-OUT.md` warn-only. A follow-up cycle (or a small-change cycle) can update the verifier to read the new convention. Not opened in this cycle.

2. **No verifier check yet for `.cdd/unreleased/{N}/` cycle directory hygiene.** No current mechanical check ensures (a) every `.cdd/unreleased/{N}/` directory has the canonical filenames the protocol expects, (b) that `release/SKILL.md` §2.5a actually moved every cycle directory before tagging, or (c) that filenames inside cycle directories follow the `{role}-{purpose}.md` / `{purpose}.md` convention. This is mechanizable; γ can pick it up as a follow-up MCA if drift becomes a pattern.

3. **Existing `.cdd/releases/{X.Y.Z}/` directory layouts are mixed.** Some prior releases (3.58.0) use per-cycle `{N}.md` files under `{role}/` directories; others (3.59.0, 3.60.0) use a single `CLOSE-OUT.md` per role. Neither matches the new `{X.Y.Z}/{N}/{role-prefix}.md` shape this cycle introduces. Not migrated — frozen-snapshot policy (§5.6) forbids semantic content changes after release. Future releases will use the new shape; old releases are warn-only legacy.

4. **Path-notation discrepancy with original issue body resolved by issue edit.** The first dispatch (and my early commit `aac0607`) implemented `.cdd/unreleased/{N}/{role}.md` based on the original issue body. The issue was edited mid-cycle to relax `{role}.md` into "role distinguished by filename" with examples like `alpha-closeout.md`, `beta-review.md`, `self-coherence.md`. This cycle's final implementation matches the edited issue. The mid-cycle reshape required two passes through the role-skill files (the second pass is what landed in this commit). The pattern (issue edits during α implementation invalidating prior α work) is a γ-side pattern observation more than an α-side one — α correctly stopped, re-read the issue when prompted, and reworked. Not raising as α-side process debt.

## Review-readiness | round 1

| Pre-review gate row (`alpha/SKILL.md` §2.6) | State | Evidence (observed at write time) |
|---|---|---|
| 1. branch rebased onto current `main` | **Yes (transient)** | `origin/main` HEAD = `9dbb107a50c3aa5dc0c120edf9386cc195992647` at observation time. Branch `claude/cdd-tier-1-skills-pptds` was checked out at that SHA on cycle start. The first commit on this branch (`aac0607`, the WIP checkpoint) was made on top of `9dbb107a`. The second commit (this one) is on top of `aac0607`. Re-validate immediately before β review per §2.6 transient-row rule. |
| 2. `.cdd/unreleased/{N}/self-coherence.md` carries CDD Trace through step 7 | **Yes** | This file §CDD Trace, rows 0–7a all present and filled. |
| 3. tests are present, or explicit reason none apply | **Explicit reason — none apply** | This is a pure documentation/coordination-protocol change. No code, no test surfaces. Verifier scripts under `commands/cdd-verify/` are out of scope per the issue. |
| 4. every AC has evidence | **Yes** | This file §AC Coverage, all 12 ACs marked Met with concrete file:section evidence. |
| 5. known debt is explicit | **Yes** | This file §Known debt, four items enumerated. |
| 6. schema / shape audit completed when contracts changed | **Yes** | The new `.cdd/unreleased/{N}/` filename schema is named in `CDD.md` §Tracking + §5.3a; verified all consumers (alpha, beta, gamma, review, release, post-release, operator skills) reference the canonical filenames; `grep -E "/(alpha\|beta\|gamma)\.md\b"` across the cdd skill tree returns zero non-`self-coherence.md` matches. |
| 7. peer enumeration completed when closure claim touches a family of surfaces | **Yes** | Peer set for "skills that mention coordination": {alpha, beta, gamma, review, release, post-release, operator, design, top-level SKILL.md, issue, plan}. Modified: 9 files. Exempted: design (zero PR refs), issue (zero PR refs, not in declared list), plan (zero PR refs, not in declared list). Each exemption justified by `grep` returning zero PR coordination references in the file. |
| 8. harness audit completed when a schema-bearing contract changed | **Yes (with explicit out-of-scope verifier note)** | Non-skill harnesses that consume the close-out path schema: `cn-cdd-verify` (asserts `CLOSE-OUT.md` legacy form). Out of scope per issue. Documented as known debt #1. No other harnesses found (no shell test fixtures, no CI emitters touch this schema). |
| 9. post-patch re-audit completed after any mid-cycle patch | **Yes** | Two patches in this cycle: (a) initial implementation with rigid `{N}/{role}.md` per the original issue body; (b) rework after operator surfaced the issue edit, retargeting to descriptive filenames. After patch (b), every file in the diff was re-grepped for `alpha.md` / `beta.md` / `gamma.md` / `{role}.md` references; final state has zero such references except in the historical-evidence prose ("Derives from #274 PR-list Monitor"). |
| 10. branch CI green on the head commit | **Provisional — local validation only; β waits for green** | No local CI runner available to α in this environment. The diff is markdown-only (no code, no tests). β should observe branch CI on the head commit before `git merge` per `alpha/SKILL.md` §2.6 row 10. The artifact's intent is signaled here per §2.6 transient-vs-durable rule. |

**Round-1 review-request:** ready for β.

- Branch: `claude/cdd-tier-1-skills-pptds`
- Base SHA: `9dbb107a50c3aa5dc0c120edf9386cc195992647` (origin/main at observation time)
- Head SHA: pending this commit (will be set by the commit that lands this artifact)
- Branch CI: deferred to β (local CI unavailable; markdown-only diff)
- Diff scope: 9 modified skill files + 6 in-version cycle migrations + 1 new self-coherence artifact

β: please read this file, then `.cdd/unreleased/268/{alpha,beta,gamma}-closeout.md` (verify migrations preserved content), then `git diff main..claude/cdd-tier-1-skills-pptds` for the skill-file changes. Verdict at `.cdd/unreleased/283/beta-review.md`.
