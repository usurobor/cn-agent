# α Close-Out — 3.60.0 (#235, PR #276)

## Cycle summary

| | |
|---|---|
| Issue | #235 — *B2: cn build --check validates command entrypoints + skill paths against manifest*. Selected by γ from the encoding-lag table (issue listed as "growing"). Single-issue B2 bundle per `docs/gamma/cdd/ISSUE-CONSOLIDATION-ANALYSIS.md` §3 B2. Closed on merge of PR #276 (`Closes #235`). |
| PR | #276 — *#235 cn build --check: validate command entrypoints + skill paths*. Squash-merged to main as `d814e16`. |
| Release | Pending — δ disconnect not yet cut at close-out write time. Anticipated 3.60.0 (next minor; feature add). δ may renumber. |
| Dispatch | γ → α: project `cnos`, issue #235, Tier 3 skills `eng/go`, `eng/test`. Branch `claude/alpha-tier-3-skills-M8Vce` (harness-assigned, non-canonical per CDD §4.2 — same shape β has flagged in prior cycles). |
| Work shape | Small substantial cycle (MCA). 2 files in PR #276 (both modified); +381 / −0 net; 1 commit on-branch (`f7d27b4`). |
| Tier 2 skills applied | Coding bundle equivalents: `eng/go` (runtime work), `eng/test` (proof depth). Already covered by issue's Tier 3 list. |
| Tier 3 skills applied | `eng/go` (parse/read split honored — new validator reads `cn.package.json` once via `pkg.ParseFullManifestData`, the canonical full-manifest parser already used by `internal/discover/discover.go`; no new parallel parser introduced; explicit error wrapping; deterministic `sort.Strings` on issue lists for stable CI output); `eng/test` (negative-space coverage — entrypoint-is-directory, escapes-package-root, container-namespace-exempt, resource-subdir-exempt; surface-enumeration test planning with explicit pass + fail + exempt for each AC). |
| MCA selection | Issue named the gap (manifest claim ↔ filesystem fact mismatch shipping silently) and the canonical fix path (one validator pass in `pkgbuild/build.go::CheckOne`). α's job was to land the surgical patch + tests + verify the existing I1 CI gate already routed the new error classes. No new CI workflow needed. |
| Review rounds | **1** — R1 APPROVED at `f7d27b4` (state=COMMENTED per `review/SKILL.md` §7.1 shared-identity convention; body verdict explicit). β merged at `d814e16` after CI 7/7 green. |
| Findings against α | **None.** β's review surfaced 4 observations (N1–N4) explicitly out-of-scope for this PR (pre-existing parallel-parser debt, pre-existing `runBuild`/`runCheck` non-gating, pre-existing symlink containment, γ-side issue-quality on a stale doc path). All filed for γ's PRA / future cycle planning, not for action on this PR. |
| Mechanical ratio | 0/0 — N/A (no findings). |
| Local verification | `go build ./... && go vet ./... && go test ./...` green; `go test -race ./internal/pkgbuild/...` green; `go build -o /tmp/cn ./cmd/cn && /tmp/cn build --check` against live `src/packages/` reports all 5 packages valid; AC4 deliberate-break dry-run on a copy of `cnos.cdd` surfaces both new error classes with exit 1. CI on merged head `d814e16`: 7/7 green (go, kata-tier1, kata-tier2, I1 Package/source-drift, I2 Protocol-contract-schema-sync, 2× notify). |
| Concurrent main activity | None during review. Last main commit pre-cycle was `48d99a3` (3.59.3 δ disconnect, 2026-04-26); main did not advance during R1 → merge. |

## Findings (α-side observations)

**None.** R1 returned APPROVED with zero findings. β explicitly framed N1–N4 as observations, not findings, and named each as pre-existing or out-of-scope.

## Cycle Iteration (CDD §9.1)

### Triggers fired

- [ ] review rounds > 2 — actual: 1.
- [ ] mechanical ratio > 20% with ≥ 10 findings — N/A (0 findings).
- [x] **avoidable tooling / environmental failure** — α-side polling miss. β posted the round-1 review at 22:35:46Z. The MCP webhook subscription did not deliver a `<github-webhook-activity>` event into α's session. α did not establish a polling loop after PR-open and did not re-poll `get_reviews`. Operator intervention was the wake (`Troubleshoot why you missed review`). β's review independently documents the symmetric β-side miss in §β review-skill miss surfaced this round — same root-cause class, opposite direction. N=2 in this single cycle.
- [ ] loaded skill failed to prevent a finding — no findings to attribute.

### Friction log (α-side)

1. **Polling discipline failure (shared-identity webhook treated as primary).** CDD.md §Tracking is explicit: *"GitHub event notifications are unreliable when agents share a GitHub identity (self-authored PR comments, CI status changes, and review events may not trigger notifications). All roles must track issue and PR activity via periodic polling, not GitHub event subscriptions."* α subscribed to `mcp__github__subscribe_pr_activity` immediately after opening PR #276, made one synchronous poll, and stopped. α inverted the priority: webhook treated as primary, polling treated as the redundant cross-check. The webhook never woke α for β's R1 review; the operator had to surface it manually.

2. **`get_status` is structurally misleading for this repo.** Even at close-out write time, `mcp__github__pull_request_read method=get_status` returns `state: pending, total_count: 0` for PR #276 despite all 7 check-runs being complete and successful. cnos uses GitHub Actions check-runs (`get_check_runs`), not the legacy commit-status API. A polling loop watching `get_status` would have shown "pending" forever — silent failure on a surface α would have trusted.

3. **`get_reviews` is the only authoritative review surface and α never re-polled it.** α treated `get_reviews` as redundant with the webhook. After the initial empty poll at PR-open, α did not re-poll it. When the operator pinged α to "check the current CI status and if there are any unresolved review comments," α re-polled `get_review_comments` and `get_comments` (both still empty) but did not re-poll `get_reviews` until the operator forced a second prompt explicitly naming the missed review.

### Root cause (α-side reading)

**Polling/webhook priority inversion.** CDD.md §Tracking names polling as primary and webhook as complement. α treated webhook as primary and polling as one-shot. The single sync poll at PR-open was not a polling loop; subscribing to webhook activity is not equivalent to polling. The §Tracking spec says this directly; α did not follow it.

A secondary contributing cause: surface-selection error. α polled `get_status` (legacy commit-status, structurally wrong for GitHub Actions repos) and `get_review_comments` + `get_comments` (which carry inline-on-diff comments and issue-style comments respectively, not top-level submitted reviews). The top-level review surface is `get_reviews`, and α did not re-poll it after PR-open until the operator forced a second prompt. The MCP toolset offers four overlapping read methods on the PR object; α did not have a clear mental model of which surface carries which event class.

### Engineering level (α-side reading, for γ to adjudicate)

- **L5 (local correctness):** met. `go build ./... && go vet ./... && go test ./... && go test -race ./internal/pkgbuild/...` green. CI 7/7 green on merged head `d814e16`. No semantic findings from β.

- **L6 (system-safe execution):** met. No cross-surface drift. The new validator and runtime activation/dispatch consume the same authority surface (`pkg.ParseFullManifestData` for commands, filesystem-presence for skills). Pre-implementation peer enumeration (build-time validator vs runtime command dispatch in `internal/discover/discover.go` vs runtime skill activation in `internal/activation/index.go`) caught all three surfaces on the first pass; β's review verified peer agreement. Negative-space coverage explicit at test-write time (4 exempt-case tests guard against false-positives on resource subdirs and namespace containers).

- **L7 (system-shaping leverage):** not attempted. The cycle was MCA, not system-shaping. The new validator converts what was a silent runtime-time degradation (entrypoint missing → dispatch failure at exec; skill dir without SKILL.md → silently invisible to activation) into a noisy build-time failure — that is a one-step coherence improvement, not a system-boundary move.

Per §9.1 "lowest miss" rule: **L6** (with one operational miss in α's polling discipline that did not affect the diff itself).

## What worked

1. **Pre-implementation peer enumeration caught the surface set on the first pass.** Before writing the validator, α enumerated the three surfaces that consume the same authority: build-time validator (this PR), runtime command dispatch (`internal/discover/discover.go`), runtime skill activation (`internal/activation/index.go`). β's review confirmed all three were considered and used the same parser/authority. No round-2 needed for cross-surface alignment.

2. **AC2 rule was tightened mid-implementation when the first version false-positived.** The first draft of `checkSkillDirectories` flagged every leaf directory under `skills/` missing SKILL.md. Running `cn build --check` against live `src/packages/` immediately surfaced two false positives (`cnos.core/skills/naturalize/references/` and `cnos.eng/skills/eng/code/references/` — legitimate resource subdirectories of parent skills). α refined the rule to "top-level subtree contains SKILL.md" before any commit. The friction was caught by the live-package check, not by tests.

3. **AC4 deliberate-break dry-run gave end-to-end verification before push.** α copied the repo to `mktemp -d`, injected a fake entrypoint in `cnos.cdd/cn.package.json` and an `skills/orphan-skill/notes.md` directory, then ran the freshly-built `cn build --check` and confirmed both new errors fired with exit 1. This proved the I1 CI gate would route the new error classes through the existing exit-code mechanism without any workflow change.

4. **Surface-enumeration test planning.** Tests were planned by enumerating the new `Issues` string families (entrypoint-missing, entrypoint-is-not-regular-file, entrypoint-escapes-root, skill-dir-missing-SKILL.md) plus the exempt cases (container-namespace, resource-subdir). Each new diagnostic-string family got at least one test. This is the explicit application of `eng/test` §3.13 ("cover new surfaces, not just preserved ones") that prior cycles' close-outs (3.59.0 §F3) named as a recurring miss.

5. **Diff-local generalization of `eng/go` §2.17.** α specifically did not introduce a parallel parser. The new code uses `pkg.ParseFullManifestData` — the same canonical full-manifest parser already in use at `internal/discover/discover.go:60`. Build-time and runtime now consume one parser for the same fact. The pre-existing `pkgbuild.PackageManifest` parallel parser remains as documented out-of-scope debt (PR §Known debt; β's N1).

## What didn't

1. **Polling discipline failure.** Discussed in §Cycle Iteration friction log items 1–3 and root cause. Same shape as β's documented miss this cycle. The operator's intervention was the necessary wake — without it, α would have remained idle indefinitely under the false belief that the webhook would deliver review events.

## Observations and patterns (α-side, factual only)

### 1. Shared-identity webhook delivery fails in both directions in the same cycle

**Pattern.** Cycle #235 saw two independent webhook-routing failures across the dyad:
- **β-side**: `mcp__github__list_pull_requests head=usurobor:claude/alpha-tier-3-skills-M8Vce` and `gh pr list --search 'closes:#N'` returned empty silently when PR #276 was already open. β's close-out section *§β review-skill miss surfaced this round* documents this.
- **α-side**: `mcp__github__subscribe_pr_activity` did not deliver β's R1 review event; α never re-polled `get_reviews` until the operator prompted explicitly.

Both failures share a root cause: shared GitHub identity suppresses self-authored event delivery. CDD.md §Tracking names the rule for both roles ("All roles must track issue and PR activity via periodic polling, not GitHub event subscriptions"). Neither role followed it on this cycle.

**Surfaces this pattern touches.** CDD.md §Tracking (already names the rule); α's `alpha/SKILL.md` §2.7 (request review — names polling as a step but does not specify which read surface to poll for which event class).

**Evidence.** N=2 in this single cycle (β-side miss + α-side miss, independently observed and named by both roles). Combined with prior cycle #274 where γ + β + α all hit first-iteration absorption (CDD.md §Tracking explicit reference: "the #274 cycle, where α + β + γ in three independent role sessions all hit first-iteration absorption"), this is N=5 across two cycles in roles spanning the full triad.

### 2. `mcp__github__pull_request_read` has four overlapping read surfaces and no clear surface-to-event mapping

**Pattern.** The MCP tool exposes `get_status`, `get_check_runs`, `get_comments`, `get_review_comments`, and `get_reviews` as separate methods. They are not interchangeable:
- `get_status` returns the legacy commit-status API; for repos using GitHub Actions check-runs (cnos), it returns `pending` permanently regardless of actual CI state.
- `get_check_runs` returns the GitHub Actions check-run state (the right surface for cnos CI).
- `get_comments` returns issue-style PR comments (the comments tab).
- `get_review_comments` returns inline-on-diff review comments grouped into review threads.
- `get_reviews` returns top-level submitted reviews (APPROVE / REQUEST_CHANGES / COMMENT).

A reviewer's "verdict" lives in `get_reviews`. Inline-on-diff findings live in `get_review_comments`. They are not the same. α polled three of the five surfaces on the operator's first prompt and missed the right one (`get_reviews`).

**Surfaces this pattern touches.** CDD.md §Tracking already enumerates `gh` and MCP query forms in the table, but the table conflates "PR comments" (one row) without distinguishing `get_comments` from `get_review_comments` from `get_reviews`. Tool documentation gap, not a CDD.md gap per se.

**Evidence.** N=1 (this cycle), but the surface-confusion shape is reproducible: any role that thinks "I checked review comments" but checked the wrong one of the three review-related methods will silently miss top-level review verdicts.

### 3. Diff-local generalization of `eng/go` §2.17 was applied successfully on this cycle

**Pattern.** Prior cycle #274's close-out (3.59.0 §F1) named the failure mode: "if I have just introduced a canonical X, every other consumer of X-shaped data in this diff must consume the canonical X." On cycle #235, α explicitly applied this generalization at design time: the new `CheckOne` code reads `cn.package.json` via `pkg.ParseFullManifestData`, which is the same canonical full-manifest parser already used in `internal/discover/discover.go`. No parallel parser introduced. β's review confirmed.

The pre-existing `pkgbuild.PackageManifest` minimal-shape parser (a pre-existing parallel of `pkgtypes.PackageManifest`) remains as out-of-scope debt — explicitly named in PR §Known debt and in β's N1. The cycle did not deepen the parallelism even while consuming the full manifest for new validation.

**Surfaces this pattern touches.** `eng/go` §2.17 (already updated in 3.59.0 to name the diff-local generalization); `internal/pkgbuild/build.go` (now consumes the canonical parser for new code).

**Evidence.** N=1 (this cycle, positive case — the rule was applied successfully, the L6 cross-surface miss from #274 did not recur in the analogous shape).

### 4. AC4 dry-run pattern: copy + inject + run + diff is cheap and gives end-to-end confidence

**Pattern.** α did not modify the working tree to verify AC4. Instead: `cp -a /home/user/cnos/. $(mktemp -d)/`, inject the breaks in the copy, run the freshly-built binary against the copy, observe outputs, delete the copy. ~30 seconds of operator-visible work, end-to-end proof that the new error classes flow through I1's exit-code-driven CI failure mechanism. The pattern generalizes to any "verify CI gate fires when input is broken" requirement.

**Surfaces this pattern touches.** `eng/test` §2.9 (test artifacts and projections, not just return values — the projection here is exit code + stderr surface routed through CI); `alpha/SKILL.md` §2.6 row 8 (harness audit — the dry-run is a harness audit at the CI level).

**Evidence.** N=1 (this cycle, applied successfully). Pattern is reusable for any cycle that adds a new validator hooked into an existing CI gate.

---

Signed: α (`alpha@cdd.cnos`) · 2026-04-26 · merge commit `d814e16` · release commit pending δ disconnect · cycle level per §9.1: L6 (one operational miss in α's polling discipline; no diff findings).
