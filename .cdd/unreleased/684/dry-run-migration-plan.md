# Dry-run migration plan — `.cn-sigma/logs/**` history (cnos#684, AC7)

**Status: PLAN ONLY. No command in this document has been executed by this cell.** α (this cycle) performs **no** destructive push, **no** history rewrite, **no** force-push, **no** ref deletion, and **no** `git rm` of live paths on `main`. Every step below is described for operator review and, eventually, operator execution — see §"Operator-executed, separately-gated step" at the bottom, which is binding, not a formality.

This plan implements the operator's **relocate-then-strip** decision (cnos#684 comment @2026-08-01T07:50:17Z): the already-committed `.cn-sigma/logs/**` content is imported to the durable orphan ref defined in `docs/reference/conventions/AGENT-ACTIVATION-CHANNEL-v1.md` §2 **before** any removal from `main`. This honors the issue's "durable ref, not deletion" framing and Kernel §2.1 (no silent drops) — see the convention doc's §6 Promotion boundary for the same no-silent-drops principle applied to channel content generally.

## 1. Target-path scope

**In scope for the strip:** `.cn-sigma/logs/**` only — the dated activation-stream files (`YYYYMMDD.md`) under that path, across all of `main`'s history. Per the operator's dispatch-authorization comment, this is "~1024 commits / ~10 MB … the dominant source of `main` HEAD churn."

**Explicitly excluded — `docs/development/board/**`.** The operator's scoping note is explicit: board-map regeneration commits are *not* part of this strip's target scope. The board is kept current-state per cnos#681, and its churn is already handled going-forward by the debounce landed in `328451d4` (tactical cadence cut on `claude/cds-dispatch-manual-trigger-w8l7ub`: agent-admin heartbeat 4×/hr → hourly, board-map per-issue-event trigger removed). If any future revision of this plan is tempted to widen the strip to include board history, that requires a *separate*, explicitly-authorized scoping decision — this document does not authorize it and this cycle does not decide it.

No other path under `main` is in scope. In particular: `.cn-sigma/state/`, `.cn-sigma/spec/`, and any convention `README` under `.cn-sigma/` are **not** touched by the strip — per the target architecture (`AGENT-ACTIVATION-CHANNEL-v1.md` §8 "Rename status"), `.cn-sigma` HEAD is meant to retain exactly this residual current-state content after the strip lands.

## 2. Sequencing (binding order)

1. **AC6 precondition gate — MUST pass first.** Run `.cdd/unreleased/684/verify-channel-reconstruction.sh` (or an updated version of it) against a full-history checkout of `main` and the post-relocate orphan-ref import. A `FAIL` exit (missing blobs) blocks every step below. See that script's header for why this cell cannot run it itself.
2. **Blast-radius sign-off (§4 below).** Every checklist item in §4 must have a stated, reviewed remap/re-anchor strategy before the strip step (step 5) runs — not merely "considered," but actually executed or scheduled where the strategy requires prior action (e.g. #682 coordination, PR freeze).
3. **Rollback ref created (§5 below).** The pre-rewrite backup ref must exist and be verified pushed *before* any rewrite touches `main`.
4. **Relocate step** (import; see §3.1) — populates the orphan ref(s) with the full `.cn-sigma/logs/**` history. This step does not touch `main` and is comparatively low-risk; it can run ahead of the strip step with no rollback-ref dependency of its own, but AC6 verification (step 1) is only meaningful *after* this step completes, so in wall-clock terms this precedes step 1.
5. **Strip step** (§3.2) — the actual `main` history rewrite. Gated on steps 1–3 above all being green. This is the operator-executed, separately-gated action; see the closing section.

(Restating the wall-clock order: relocate → AC6 verify → blast-radius sign-off → rollback ref → strip. Steps 2 and 3 may proceed in parallel with each other but both must complete before step 5.)

## 3. `git filter-repo`-style command plan (NOT executed)

All commands below are illustrative — they describe what *would* run, in what order, against what inputs. None have been invoked by this cell. `git filter-repo` is assumed as the tool (per the operator's comment naming "`git filter-repo` (or equivalent)"); an equivalent tool (e.g. `git-filter-branch` — discouraged for performance, or BFG for the strip-only half) could substitute without changing this plan's structure.

### 3.1 Analyze first (obtain exact commit counts — AC7 oracle item (c))

Before writing any rewrite command with real numbers, the operator runs `git filter-repo --analyze` against a disposable clone of `main`:

```bash
# Disposable clone — never the working repo directly.
git clone --no-local <cnos-repo> /tmp/cnos-684-analyze
cd /tmp/cnos-684-analyze
git filter-repo --analyze
# Inspect .git/filter-repo/analysis/path-all-sizes.txt and
# .git/filter-repo/analysis/README.md for:
#   - exact commit count touching .cn-sigma/logs/**
#   - total blob count / bytes under that path
#   - which commits would become empty once that path is removed
#     (a commit that ONLY touched .cn-sigma/logs/** in that snapshot)
```

This plan intentionally does **not** hardcode a commit count — the operator's dispatch comment gives an approximate figure ("~1024 commits / ~10 MB") that this plan treats as directional, not authoritative. `--analyze` output is the authoritative source at execution time, since `main` will have advanced between this plan's authoring and its execution.

### 3.2 Relocate step (import into the orphan ref)

Populates `refs/heads/channels/sigma/cnos-to-home` (the foreign→home ref defined in `AGENT-ACTIVATION-CHANNEL-v1.md` §2) with the full historical content of `.cn-sigma/logs/**`, preserving history (not just final-state content) so the import is itself an auditable append-only record, consistent with the channel's own append-only design.

```bash
# Fresh clone, disposable — filter-repo rewrites in place.
git clone --no-local <cnos-repo> /tmp/cnos-684-relocate
cd /tmp/cnos-684-relocate

# Keep ONLY the target path's history; everything else is discarded from
# this disposable clone's perspective (the source repo is untouched — this
# clone becomes the orphan-ref payload, not a replacement for main).
git filter-repo --path .cn-sigma/logs/ --path-rename .cn-sigma/logs/:

# Result: a rewritten history containing only .cn-sigma/logs/** commits,
# with paths rebased to the ref's root (path-rename strips the
# .cn-sigma/logs/ prefix so the orphan ref's tree is the log files
# directly, not nested three directories deep).

# Graft this as an orphan ref onto the real cnos-to-home target:
cd /tmp/cnos-684-relocate
git remote add target <cnos-repo>
git push target HEAD:refs/heads/channels/sigma/cnos-to-home
# (First push to a not-yet-existing ref is NOT a force-push and is NOT a
#  rewrite of any existing ref — it is a normal ref creation. This is the
#  one step in this plan that is closest to "safe to run early," since it
#  creates new state rather than mutating existing state.)
```

**Verification checkpoint:** immediately after this step, run `verify-channel-reconstruction.sh` (§"Sequencing" step 1) against this import before proceeding.

### 3.3 Strip step (the actual `main` rewrite — NOT executed)

Removes `.cn-sigma/logs/**` from `main`'s history entirely (not just HEAD — every ancestor commit that touched the path is rewritten), which is what actually shrinks `main`'s object count and produces the "`.cn-sigma` HEAD retains only current cursors/state" end state named in AC5.

```bash
# Disposable clone of main, NOT the operator's working checkout.
git clone --no-local <cnos-repo> /tmp/cnos-684-strip
cd /tmp/cnos-684-strip

# Invert-match: keep everything EXCEPT the target path.
# docs/development/board/** is deliberately absent from this invocation —
# its presence here would be the scoping violation flagged in §1.
git filter-repo --invert-paths --path .cn-sigma/logs/

# --analyze (§3.1) should be re-run against this output to confirm the
# resulting commit graph matches expectations before push.
```

This produces a **rewritten `main`** with every commit SHA after the first commit that ever touched `.cn-sigma/logs/**` changed (parent-hash cascade — `git filter-repo`, like any history rewrite, changes the SHA of every descendant of a rewritten commit, not only commits that directly touched the target path). This is the step with by far the largest blast radius (§4) and is never run against the live repo directly — only against a disposable clone, whose result becomes the force-pushed replacement for `main` in the operator-executed step.

## 4. Blast-radius checklist

Every item below requires a stated remap/re-anchor strategy before the strip step (§3.3) is authorized to run for real.

### 4.1 Tags `3.82.x`

**Impact.** `git filter-repo` rewrites refs it can see in the same repo, including tags, by default. A `3.82.x` tag pointing at (or descended from, via any ancestor) a commit that touched `.cn-sigma/logs/**` gets a new SHA if filter-repo is allowed to rewrite it in place.

**Strategy.** Treat release tags as immutable identifiers — do not let the strip silently move them:
1. Before the strip clone runs, enumerate all `3.82.x` tags and their current SHAs: `git tag -l '3.82.*' | xargs -I{} git rev-parse {}`. Record this list as the pre-rewrite baseline.
2. Run the strip (§3.3) with filter-repo's tag handling left at default (it will remap tags automatically using its internal commit-map) inside the disposable clone — do not push yet.
3. Inspect `.git/filter-repo/commit-map` (old-SHA → new-SHA) for every baseline tag's old SHA; confirm each maps to a new SHA that is still reachable from the rewritten `main`.
4. **Decision point for the operator:** either (a) accept the remapped tag SHAs — the tag names stay the same, the SHA they point to changes, requires re-pushing tags (itself a mutation operators must explicitly authorize, since GitHub treats a tag SHA change as effectively a forced tag update), or (b) explicitly preserve the *old* SHA for `3.82.x` tags by excluding them from the rewrite's ref scope and instead recording their pre-rewrite SHA as immutable historical pointers (resolvable only via the rollback ref, §5, since the old SHA becomes unreachable from the new `main`). This decision is **not made by this plan** — it is named here as a required operator decision before step 5 runs.

### 4.2 Open PRs

**Impact.** Every open PR's base is (transitively) some commit on `main`. After the strip, that base commit's SHA no longer exists on `main` — the PR's diff base is gone, and most Git hosts will show the PR as needing a rebase at best, or as unable to compute a merge base at worst.

**Strategy.**
1. Snapshot every open PR immediately before the strip: `gh pr list --state open --json number,headRefName,baseRefName,headRefOid,baseRefOid`.
2. Use the strip's `commit-map` (old-SHA → new-SHA, same artifact as §4.1 step 3) to translate each PR's recorded `baseRefOid` to its post-rewrite equivalent.
3. For each open PR, the operator (or the PR's author) re-bases the PR branch onto the new `main` tip, using the commit-map to resolve any commit references the PR branch itself carries (e.g. if the PR branch's own history includes commits authored before the strip whose parent chain touches the rewritten range).
4. **Recommended operational step (not part of the rewrite itself):** schedule the strip during a PR-freeze window — no new PRs opened, existing PRs' authors notified in advance — to bound the number of PRs needing this treatment to whatever is open at freeze time.

### 4.3 cnos#682's CDD-by-SHA custody

**Impact.** cnos#682 ("architecture: dematerialize closed CDD cells from HEAD…") tracks CDD cell artifacts by exact commit SHA reachable through `main`'s ancestry — the whole point of that design is that a closed cell's artifacts remain findable via `main`'s history at a specific SHA, even after the cell's working files are dematerialized from HEAD's tree. A full-history rewrite of `main` (§3.3) changes essentially every commit SHA from the first `.cn-sigma/logs/**`-touching commit forward — which, per the operator's ~1024-commit estimate, is likely to be the large majority of `main`'s commit graph. Any SHA citation `#682`'s custody mechanism has recorded (or will record) against a pre-strip commit becomes stale the moment the strip lands.

**Strategy.** This cell does **not** modify #682's scope (per the operator's guardrail and γ's scaffold — #682 is referenced here only as a blast-radius item). The remap/re-anchor strategy named here is a **coordination requirement**, not an action this cell takes:
1. The strip step (§3.3) MUST NOT be authorized to run until #682's custody mechanism either (a) has landed and provides a documented remap procedure that consumes a filter-repo-style `commit-map` (old-SHA → new-SHA) to re-anchor its recorded citations, or (b) #682's owner has explicitly signed off that no such remap is needed at strip time (e.g. because #682 has not yet landed any SHA citations that would be affected).
2. Whoever runs the strip preserves the `commit-map` artifact (§4.1 step 3) as a durable output of the rewrite — not a disposable intermediate — specifically so #682's remap procedure (whenever it lands) has the translation table it needs, even if the remap itself happens after the strip.
3. This is named as a **hard sequencing dependency**: the strip (§3.3) is gated on #682-side readiness in addition to the AC6/rollback gates above. This plan does not resolve #682's readiness — it only names the dependency so the operator does not discover it at strip time.

## 5. Rollback ref

Before the strip step (§3.3) is authorized to run for real (i.e. before any force-push replaces the live `main`), the operator creates a full mirror backup of pre-rewrite `main`:

```bash
# Full mirror clone of the live repo, pushed as a new, clearly-named ref —
# NOT a tag on the live repo (tags are themselves subject to the §4.1
# remap question; a backup ref sidesteps that ambiguity entirely by living
# outside the tag namespace this plan is reasoning about).
git clone --mirror <cnos-repo> /tmp/cnos-684-backup.git
cd /tmp/cnos-684-backup.git
git push origin refs/heads/main:refs/heads/backup/main-pre-cnos684-rewrite
```

**Naming:** `refs/heads/backup/main-pre-cnos684-rewrite`. Rationale: `backup/` prefix keeps it out of normal branch listings and out of any pattern-based branch-protection rule scoped to `channels/*` or release branches; `main-pre-cnos684-rewrite` names both what it's a backup of and which cycle's rewrite it precedes, so a future reader does not have to cross-reference a date to understand why the ref exists. This ref should be protected against deletion (same branch-protection mechanism as `AGENT-ACTIVATION-CHANNEL-v1.md` §5's invariant #5) for as long as any downstream consumer (§4.1's tag decision, §4.3's #682 remap) might still need to resolve a pre-rewrite SHA.

## 6. Operator-executed, separately-gated step

**Everything above this line is a plan, not an action.** Per the operator's binding decision (cnos#684 comment @2026-08-01T07:50:17Z): "the mechanical worker produces PRs into `main`; it does not and must not rewrite history, force-push `main`, or delete refs." Concretely, the following actions are **explicitly out of this cell's scope, to be performed only by the operator (or an operator-authorized actor), as a separate, later, explicitly-gated action**:

- Running `git filter-repo` (§3.2 relocate, §3.3 strip) against anything other than a disposable clone.
- Pushing the relocate step's output to the real `refs/heads/channels/sigma/cnos-to-home` (this specific push — first-creation of a not-yet-existing ref — is the lowest-risk item in this plan, but it still is not performed by this cell, since it requires the `.cn-sigma/**`-adjacent context this cell structurally cannot access).
- Force-pushing the strip step's rewritten history to the live `main`.
- Creating, moving, or deleting any tag.
- Deleting any ref.
- Running the AC6 verification script (`verify-channel-reconstruction.sh`) against real `.cn-sigma/**` content.

**Self-verifiable claim:** this plan document contains no evidence of having been executed. `git log main` shows no rewritten SHAs beyond this cycle's own additive doc commits; `git tag` shows no new or moved tags; no force-push has occurred. β (or any later reviewer) can confirm this directly: `git diff main...cycle/684 --stat` touches only doc/plan files under `docs/reference/conventions/` and `.cdd/unreleased/684/`, never `.cn-sigma/**`, never `.github/workflows/*`, never a ref or tag.
