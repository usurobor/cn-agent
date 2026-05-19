<!--
section-manifest:
  planned: [r1-verdict, contract-integrity, acs, issue-contract, surface-containment, findings, ci-status, artifact-completeness, notes]
  completed: [r1-verdict, contract-integrity, acs, issue-contract, surface-containment, findings, ci-status, artifact-completeness, notes]
-->

# β review — #379 agent/activate skill

## Round 1 — verdict

**Verdict:** APPROVED

**Round:** 1
**Fixed this round:** N/A (R1)
**Review SHA:** `e78aa9829f07e82c567c2684b45af4601f9d1387` (cycle/379 HEAD)
**Review base:** `7a7f7152af649ea93cebe5f909fbf1397d809547` (origin/main HEAD, re-fetched synchronously at review start)
**Branch CI state:** N/A — repo has no required workflows for the Go package; β ran `go test -count=1 ./internal/activate/...` locally on cycle HEAD (PASS, 29 tests) and on the merge tree (PASS).
**Merge instruction:** `git merge --no-ff cycle/379` into main with `Closes #379` in the merge commit body.

The cycle ships exactly what the issue body specifies: a new conformant skill at `src/packages/cnos.core/skills/agent/activate/SKILL.md`, an evolved Go renderer that reads the skill's machine-readable load-order block, and a source-of-truth test that demonstrates skill edits drive renderer output. Surface containment is exact (five files; all allowed). Pre-merge gate rows 1–4 pass. No unresolved findings of any severity. Search space closed — no further blocker found in the contract.

---

## §2.0.0 Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | yes | Issue body §Status truth table matches what shipped; α's self-coherence §Gap reiterates without conflation. |
| Canonical sources/paths verified | yes | All paths in SKILL.md resolve in tree (`KERNEL.md`, CA skill set, cdd/activation/SKILL.md, activate.go, activate_test.go, cmd_activate.go). |
| Scope/non-goals consistent | yes | Cycle scope sizing held; non-goals untouched (no hub README patches, no cmd_activate.go change, no cnos.xyz, no cn-doctor edits). |
| Constraint strata consistent | yes | "Everything is a skill" invariant restored for activation; layering rule preserved in skill §2.1. |
| Exceptions field-specific/reasoned | yes | Fallback path in §4.3 is field-specific to "skill not vendored" and reasoned (`cn activate` must remain useful before `cn deps restore`). |
| Path resolution base explicit | yes | Renderer reads from `<hub>/.cn/vendor/packages/cnos.core/skills/agent/activate/SKILL.md`; base is the hub root resolved at `Run`. |
| Proof shape adequate | yes | AC7 oracle = positive (edit changes output) + negative (no-op pretender would produce identical output) + fallback (manifest-only); all three exercised by tests. |
| Cross-surface projections updated | yes | Existing `TestReadFirstSection_*` tests updated to the new canonical ordering; α flagged the observable-output delta in §Debt 2 of self-coherence. |
| No witness theater / false closure | yes | `TestSkillIsSourceOfTruthForReadFirstOrder` two-phase edit + `out1 != out2` coherence assertion is real evidence, not a file-exists check. |
| PR body matches branch files | n/a | No PR opened; cycle ships via direct merge by β per CDD `release/SKILL.md`. |
| γ artifacts present (gamma-scaffold.md) | yes | `.cdd/unreleased/379/gamma-scaffold.md` present on cycle branch (rule 3.11b compliance). |

---

## §ACs — per-AC verdict

| AC | Verdict | Evidence (path:line or command) | Finding |
|----|---------|---------------------------------|---------|
| AC1 — skill at canonical path, ≥200 lines | pass | `wc -l src/packages/cnos.core/skills/agent/activate/SKILL.md` → 485 lines; `test -s …` true | — |
| AC2 — conforms to skill/SKILL.md | pass | Frontmatter SKILL.md:1–39 carries name/description/artifact_class/kata_surface/governing_question/visibility/parent/triggers/scope/inputs/outputs/requires/calls. Body Define→Unfold→Rules→Verify present (SKILL.md §1, §2, §3, §5). Embedded kata at §7 SKILL.md:395–445 has scenario/task/inputs/expected-artifacts/procedure/verification/common-failures/reflection. | — (see Notes 1 on kata field-name choice; not blocking) |
| AC3 — six-item load order; layering rule cited | pass | §2.1 SKILL.md:110–145 lists Kernel → CA skills → Persona → Operator → hub state → identity confirmation; §4.1 SKILL.md:301–308 has the machine-readable block in the same order. Layering rule SKILL.md:145 cites cn-sigma `threads/adhoc/20260325-session2-learnings.md` §3. Renderer-side order verified by `TestReadFirstSection_OrderedSigma` and the skill-as-source-of-truth test. | — |
| AC4 — three-tier capability matrix | pass | §2.2 SKILL.md:155–180 names tiers (a) shell+git / (b) HTTP fetch only / (c) no fetch; each has a load-path column. Tier (a) labeled "preferred when available" in both table header and §161 paragraph. Tier (c) named honestly as degraded. Tier-selection rule SKILL.md:172–176 makes the choice observable. | — |
| AC5 — README router template | pass | Heading "### 2.3. README router template" SKILL.md:182. Fenced `~~~markdown` block SKILL.md:188–218. Paste test: `awk '/^~~~markdown$/,/^~~~$/' SKILL.md \| grep -E '<[A-Z_-]+>' \| sort -u` → only `<HUB-URL>`. Self-contained; the §2.3 adoption rule SKILL.md:220 makes the single-substitution invariant explicit. | — |
| AC6 — naming disambiguation from cdd/activation | pass | §2.4 SKILL.md:230 paragraph cites both `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` (leaf: activation) and `src/packages/cnos.core/skills/agent/activate/SKILL.md` (leaf: activate). One-sentence distinction: _"They share a word, not a concern: one activates a repo under CDD; the other activates a body at a hub."_ Paths verified exact (no typos in `activate` vs `activation`). | — |
| AC7 — `cn activate` reads the skill | pass | `writePrompt` (activate.go:368–388) iterates a `readFirst []readFirstItem` slice; no in-Go literal ordering. `loadActivateSkillOrdering` (activate.go:95–107) reads `<cnDir>/vendor/packages/cnos.core/skills/agent/activate/SKILL.md`. `parseReadFirstOrderBlock` (activate.go:127–163) scans the §4.1 marker-bounded block. `TestSkillIsSourceOfTruthForReadFirstOrder` (activate_test.go:418–478) is a two-phase edit test: phase 1 writes a fixture skill with kernel-before-persona and asserts the same ordering in renderer output; phase 2 swaps the two and asserts the swap propagates; `out1 != out2` coherence check at line 475. `TestSkillFallback_NotVendored` (activate_test.go:484–508) asserts manifest-only fallback path emits canonical ordering + stderr diagnostic + preserves `run cn deps restore` guidance. Interpolation surface explicit: SKILL.md §4.2 table + `renderReadFirstLine` switch (activate.go:482–527) is one-to-one. | — |

All seven ACs PASS. The AC7 test is the strong content-driven form (not a file-exists oracle), so the γ-flagged failure mode 6 (file-read oracle masquerading as source-of-truth) does not apply.

---

## §2.0 Issue Contract

### AC Coverage

| # | AC | In diff? | Status | Notes |
|---|----|----------|--------|-------|
| AC1 | activate skill exists at canonical path | yes | met | 485 lines, exact path |
| AC2 | conforms to skill skill format | yes | met | frontmatter + Define/Unfold/Rules/Verify + embedded kata |
| AC3 | six-item load order | yes | met | §2.1 + §4.1 + layering citation; renderer emits in same order |
| AC4 | three-tier capability matrix | yes | met | (a)/(b)/(c) with paths; (a) preferred; (c) honest |
| AC5 | README router template | yes | met | heading + fenced block + single `<HUB-URL>` placeholder |
| AC6 | naming disambiguation from cdd/activation | yes | met | both paths exact; one-sentence distinction |
| AC7 | `cn activate` reads the skill | yes | met | data-driven `writePrompt`; source-of-truth test green |

### Named Doc Updates

| Doc / File | In diff? | Status | Notes |
|------------|----------|--------|-------|
| `src/packages/cnos.core/skills/agent/activate/SKILL.md` | yes (A) | created | primary deliverable |
| `src/go/internal/activate/activate.go` | yes (M) | evolved | reads skill instead of hardcoding |
| `src/go/internal/activate/activate_test.go` | yes (M) | extended | source-of-truth + fallback + parser tests |
| `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` | no | not modified | non-goal per issue body §Scope |
| `src/go/internal/cli/cmd_activate.go` | no | not modified | CLI surface preservation (non-goal) |

### CDD Artifact Contract

| Artifact | Required? | Present? | Notes |
|----------|-----------|----------|-------|
| `.cdd/unreleased/379/gamma-scaffold.md` | yes (rule 3.11b) | yes | γ-authored; carries gap, peer enumeration, mode, ten γ-flagged failure modes |
| `.cdd/unreleased/379/self-coherence.md` | yes | yes | seven sections per α/SKILL.md §2.5; review-readiness ready-for-β signal present |
| `.cdd/unreleased/379/beta-review.md` | yes (this file) | yes | R1 verdict APPROVED |
| `.cdd/unreleased/379/beta-closeout.md` | yes (post-merge) | pending | written after merge per `release/SKILL.md` |

### Active Skill Consistency

| Skill | Required by | Loaded? | Applied? | Notes |
|-------|-------------|---------|----------|-------|
| `cnos.cdd/skills/cdd/beta/SKILL.md` | β role | yes | yes | pre-merge gate rows 1–4 exercised |
| `cnos.cdd/skills/cdd/review/SKILL.md` | β review phase | yes | yes | contract-integrity + per-AC verdict format |
| `cnos.core/skills/skill/SKILL.md` | AC2 conformance check | yes | yes | frontmatter + body structure + kata checked literally |
| `cnos.cdd/skills/cdd/release/SKILL.md` | β merge phase | yes | yes | merge-test (row 3) executed in throwaway worktree |
| `cnos.eng/skills/eng/go/SKILL.md` | AC7 Go review | yes | yes | `go test` + `go vet` clean on cycle HEAD and merge tree |

---

## Surface containment

```
$ git diff --name-status origin/main..HEAD
A    .cdd/unreleased/379/gamma-scaffold.md
A    .cdd/unreleased/379/self-coherence.md
M    src/go/internal/activate/activate.go
M    src/go/internal/activate/activate_test.go
A    src/packages/cnos.core/skills/agent/activate/SKILL.md
```

Five files; all in the allowed set. No edits to `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` (γ failure-mode 10 held); no hub README patches; no `cnos.xyz/`; no `cn-doctor` surface; no `src/go/internal/cli/cmd_activate.go` change (CLI surface preserved).

The vendored bundle path the renderer reads (`<hub>/.cn/vendor/packages/cnos.core/skills/agent/activate/SKILL.md`) is constructed in tests via the `writeVendoredActivateSkill` helper (activate_test.go:401); no in-repo vendored bundle file is shipped, which is correct — vendoring is a per-hub state, not source state.

---

## §Findings

| # | Finding | Evidence | Severity | Type |
|---|---------|----------|----------|------|

No findings at any severity. Verdict-shape lint per `review/SKILL.md` §3.4a: single terminal verdict (`APPROVED`), no conditional qualifier, no split scope.

---

## §CI status

No required CI workflows configured for the Go package at branch-push time (verified by absence of `.github/workflows/` matching `go test` on cycle/379). Per α/SKILL.md §2.6 row 10 escape hatch (cited explicitly in α's self-coherence Review-readiness section), β substitutes the rule 3.10 binding gate with a local-equivalent run on cycle HEAD:

```
$ cd src/go && go test -count=1 ./internal/activate/...
ok  	github.com/usurobor/cnos/src/go/internal/activate	0.042s

$ cd src/go && go vet ./internal/activate/...
(no output — clean)
```

Source-of-truth test specifically:

```
$ go test -v -count=1 -run TestSkillIsSourceOfTruth ./internal/activate/...
--- PASS: TestSkillIsSourceOfTruthForReadFirstOrder (0.00s)
PASS
```

Merge-tree verification (β/SKILL.md row 3 non-destructive merge-test):

```
$ git worktree add /tmp/cnos-merge379/wt origin/main
$ cd /tmp/cnos-merge379/wt
$ git config --worktree user.name beta && git config --worktree user.email beta@cdd.cnos
$ git merge --no-ff --no-commit origin/cycle/379
Automatic merge went well; stopped before committing as requested
$ git diff --name-only --diff-filter=U   # unmerged paths
(empty — no conflicts)
$ cd src/go && go test -count=1 ./internal/activate/...
ok  	github.com/usurobor/cnos/src/go/internal/activate	0.047s
```

Worktree torn down. Identity (`beta@cdd.cnos`) preserved in the parent repo throughout.

---

## §Artifact completeness

`.cdd/unreleased/379/gamma-scaffold.md` is present on `origin/cycle/379` (verified by `ls .cdd/unreleased/379/` returning `gamma-scaffold.md self-coherence.md`). Rule 3.11b gate satisfied; no exemption needed.

---

## §Notes

1. **Embedded-kata field choice (informational; not a finding).** `skill/SKILL.md` §2.4 lists eight elements a kata must define: scenario, task, governing skills, inputs, expected artifacts, verification, common failures, reflection. α's kata (SKILL.md §7) carries seven of these literally (scenario, task, inputs, expected artifacts, verification, common failures, reflection) and adds an eighth element labeled "Procedure to follow" instead of "Governing skills." For an embedded kata in the activate skill itself, the governing skill is trivially this skill plus the load-time skills it calls (Kernel, CA set, Persona, Operator) — which the "Procedure to follow" section makes concrete by walking through their exercise. The kata is fully exercisable; this is naming, not a structural gap. Logged as an observation for the cycle's PRA, not as a finding (per `review/SKILL.md` rule 3.5: no phantom blockers).

2. **Renderer fallback duplicates canonical ordering.** α's §Debt 3 names this honestly: `canonicalReadFirstOrdering` in `activate.go` mirrors SKILL.md §4.1, so drift is structurally possible. α's mitigations are sound — the vendored skill is preferred whenever present, `TestSkillFallback_NotVendored` asserts the fallback path's behavior, and SKILL.md §4.3 declares the skill canonical on drift. The future `//go:embed` direction α names is the right durable fix and belongs in a follow-on cycle, not this one. No finding here.

3. **Observable-output delta is honest debt.** α's §Debt 2 names the `## Read first` ordering change explicitly (persona/operator/kernel → kernel/ca-skills/persona/operator/hub-state/identity). Section headers (`## Read first`, `## Kernel`, …) are preserved, so consumers that parse by section header are unaffected; only consumers doing string-match within `## Read first` see different content. The change is the intended consequence of AC7 and AC3 — making it surface honestly in self-coherence is the right behavior.

4. **γ-flagged failure modes audit (10 items in `gamma-scaffold.md`).** β cross-checked the AC7 source-of-truth test against γ's failure mode 6 (file-read oracle masquerading as source-of-truth). The test is two-phase with a coherence assertion (`out1 != out2`), not a file-read; the pattern γ pre-flagged is not present. Failure modes 1–5, 7–10 are also held: AC3 ordering matches §4.1 and the renderer output (1); AC4 has three tiers + load paths + tier (a) preferred + tier (c) honest (2,3); AC5 paste-tested (4); AC6 paths exact, no typos (5); fallback present and tested (7); interpolation surface explicit (8); cdd/activation/SKILL.md not modified (9); no CLI surface change in cmd_activate.go (10).

5. **Identity-truth row 1.** At session start, the worktree-local git config inherited `alpha@cdd.cnos`/`alpha` (worktree shared with α's session). β re-asserted at the worktree scope (`git config --worktree user.name beta && git config --worktree user.email beta@cdd.cnos`) before any β-authored commit. Confirmed by `git config --show-scope --get-all user.email` showing `worktree beta@cdd.cnos`. The merge-test worktree at `/tmp/cnos-merge379/wt` set its own worktree-local config explicitly with `--worktree` flag to prevent leak into the parent repo (β/SKILL.md row 3 caveat). Post-teardown, parent-repo identity verified still `beta@cdd.cnos`.

6. **No `gamma-clarification.md` on cycle branch.** Issue-edit cache-bust per β/SKILL.md "Issue-edit cache-bust" rule not triggered: only `gamma-scaffold.md` and `self-coherence.md` present in `.cdd/unreleased/379/`. Issue body read from `gh issue view 379` at session start is binding.

---

## After-review action

Per `review/SKILL.md` §After Review: APPROVED → β merges `cycle/379` into `main` with `Closes #379`, writes `beta-closeout.md`, and hands off to γ for PRA / δ for release boundary. β does NOT tag, push tags, bump version, update CHANGELOG, delete the cycle branch, or move `.cdd/unreleased/379/` to a release directory — those are δ's release boundary per `beta/SKILL.md` Role Rule 5 and `release/SKILL.md` §2.5a.
