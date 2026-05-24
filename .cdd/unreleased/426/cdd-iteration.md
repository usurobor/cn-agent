# ε cdd-iteration — cycle/426

**Issue:** [cnos#426](https://github.com/usurobor/cnos/issues/426) — Publish v3.82.0 GH release directly (second use of remote-runner-delegation doctrine; sidesteps broken release.yml build).
**Mode:** workflow + receipt only; γ+α+β-collapsed-on-δ. **Substantive** entry: records two `cdd-tooling-gap` findings with disposition `next-MCA`, both surfaced during cnos#425's first use of the remote-runner-delegation primitive and dispositioned in this cycle.

## Context

This cycle is the second use of the `remote-runner delegation` primitive landed in cnos#425. The first use (`.github/workflows/repoint-3.82.0.yml`) was authored to force-move the `3.82.0` tag from its prior commit (`fb6149c`) to `fd1d654e` so that `release.yml` would re-publish the GH release with the correct body. The tag-move step succeeded; `git ls-remote origin refs/tags/3.82.0` resolves to `fd1d654e` post-cnos#425. But the *downstream* expected cascade (tag push → release.yml triggers on the tag → release.yml's publish job creates the GH release) did NOT complete:

1. `release.yml`'s `build` job had already failed once on the original tag push (run 26340314834 at commit `fb6149c`) with `could not read Username for 'https://github.com': terminal prompts disabled` — a git-auth failure inside the build matrix. 2 of 4 platforms produced binaries before cancel. The `publish` job (gated on `needs: build`) therefore never ran.
2. The cnos#425 tag force-push from `fb6149c` to `fd1d654e` did NOT re-trigger `release.yml` at all — GH Actions' `on.push.tags` filter does not fire reliably on tag force-update events. So no new run was produced to even attempt the build.

These two infrastructure facts together blocked v3.82.0's GH release publication via the standard pipeline. cnos#426 (this cycle) sidesteps both via direct `softprops/action-gh-release@v1` invocation in a one-shot workflow; both findings are filed here as `next-MCA` so the standard pipeline is properly repaired in subsequent cycles rather than being permanently bypassed.

## Findings

### F1 — cdd-tooling-gap: release.yml build job has git-auth failure on tag-push

**Class:** `cdd-tooling-gap`
**Severity:** binding (blocks future releases from auto-publishing binaries via the standard pipeline)
**Disposition:** `next-MCA` — file a follow-on issue with the first AC below; sequence behind cnos#426's merge so the v3.82.0 GH release is published first, then repair the build path.

**The gap.** `.github/workflows/release.yml`'s `build` job, when triggered by a tag push (`on: push: tags: ['[0-9]*.[0-9]*.[0-9]*']`), fails with `could not read Username for 'https://github.com': terminal prompts disabled` on one or more matrix positions. Some `git push` or `git clone` step inside the build matrix expects HTTP basic-auth credentials that the runner's environment does not supply, and the runner's `git` is configured to disable interactive prompts (correct behavior on CI). The exact failing command has not been pinpointed in this cycle (out of scope per Hard Rule 3); candidate locations include:

- The `actions/setup-go@v5` step's `cache-dependency-path` lookup if Go module fetching attempts to clone a private dependency over HTTPS without a configured auth helper;
- The `go test ./...` step if it fetches private modules at test time;
- The `scripts/kata/run-all.sh` step (Tier 1 verification) if any kata under `scripts/kata/` performs a `git clone` or `git push` against a private remote;
- A `git fetch` inside the build for version stamping if `-ldflags "-X main.version=${{ github.ref_name }}"` indirectly invokes a git command that needs auth.

**Empirical anchor.** GitHub Actions run 26340314834, tag `3.82.0` at commit `fb6149c`. 2 of 4 platforms (matrix positions) produced binaries before the cancel; the failure was deterministic across the failing platforms rather than a flake.

**Why this is a `cdd-tooling-gap` rather than a `cdd-protocol-gap` or `cdd-skill-gap`.** Tooling-gap means a substrate / pipeline mechanism is broken or misconfigured. Protocol-gap would mean the cnos protocol shape itself lacks a naming for something. Skill-gap would mean an existing skill's content is incomplete. This is the tooling case: `release.yml` is the substrate that converts a tag push into a built+published+smoke-tested+notified release; one job in that substrate is broken. The cnos protocol around it (δ-as-role authority for releases, release-effector mechanics, β-side release authoring, remote-runner-delegation for one-shot workaround) is all functional and correctly applied — it is the substrate underneath that needs repair.

**Why disposition `next-MCA` is correct.**
- `patch-landed-in-cycle` is wrong because this cycle's matter does not modify `release.yml` (Hard Rule 3 explicitly forbids it; the cycle's scope is the workaround, not the repair).
- `no-patch` is wrong because the gap is real and recurring (every future release would re-hit the same failure if `release.yml` is not repaired).
- `next-MCA` is correct: a follow-on cycle with the scope "diagnose and fix `release.yml`'s `build` job git-auth failure on tag-push" is required. The cycle's first AC is given below; the issue scope is contained (single workflow file; mechanical CI verification via test tag push).

**First AC for the MCA follow-on:** "`release.yml`'s `build` job succeeds on a tag-push event (verified by pushing a `v0.0.0-mca-test` tag and observing all 4 matrix positions complete the build + test + package + Tier-1-verification steps without a git-auth failure)."

**Recommended additional scope hints for the MCA follow-on** (not binding; just signal for the next γ):
- Add a `git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"` step early in the build job to ensure any `git` invocation has token-based HTTPS auth available, OR
- Pinpoint and remove the auth-needing git invocation (the more surgical fix);
- Add a CI smoke test that pushes a no-op test tag in a sandbox repo or under a `if: github.event.head_commit.message contains '[release-test]'` guard, so the fix is verifiable before the next real release.

### F2 — cdd-tooling-gap: GH Actions on.push.tags does not reliably trigger on tag force-update

**Class:** `cdd-tooling-gap`
**Severity:** binding (any future tag-retarget cycle will hit the same non-trigger, blocking the doctrine-named expected cascade)
**Disposition:** `next-MCA` — file a follow-on doctrine patch to `BOX-AND-THE-RUNNER.md`'s known-failure-modes section (or, if no such section exists yet, add one).

**The gap.** When a git tag is **force-updated** (the tag already exists pointing at commit A; `git push --force origin <tag>` moves it to commit B), GitHub Actions' `on: push: tags: [<pattern>]` filter does NOT reliably trigger a new workflow run. The behavior appears to be: tag *creation* triggers a workflow run; tag *deletion* does not (correctly, per docs); tag *force-update* sometimes triggers and sometimes does not, with no clear documented rule. Empirically in cnos#425's first use, the force-move of `3.82.0` from `fb6149c` to `fd1d654e` did NOT produce a new `release.yml` run — but the workflow was correctly set up to trigger on the tag pattern, and the tag pattern matched the moved tag.

This is significant because the cnos#425 receipt's "Expected effect" (§Expected effect step 4) explicitly stated "release.yml triggers on the new tag" as the load-bearing downstream cascade. The cascade did not fire. The receipt's failure-mode "Workflow does not trigger" covered this generally; this finding makes the specific failure mode (tag force-update non-trigger) into a doctrine-named known failure mode so future receipts can plan around it rather than re-discover it.

**Empirical anchor.** cnos#425's tag force-push (workflow run for `.github/workflows/repoint-3.82.0.yml` completed successfully, tag SHA verified moved to `fd1d654e`), no subsequent `release.yml` run produced on the moved tag. The Actions tab shows the cycle/425 cycle of runs (repoint workflow + its self-delete trigger), but no release.yml run keyed on `3.82.0` post-retarget.

**Why this is a `cdd-tooling-gap` rather than `cdd-protocol-gap` or `cdd-skill-gap`.** Tooling-gap means a substrate behavior is broken or undocumented in a way that affects cnos cycles. Protocol-gap would mean the cnos cell algorithm itself has an undefined edge case. Skill-gap would mean a cnos skill's content is incomplete. This is the tooling case: GH Actions' substrate has an undocumented (or under-documented) non-trigger behavior on tag force-update; cnos cycles that rely on tag-push cascade need to know about it. The `BOX-AND-THE-RUNNER.md` doctrine has a "known failure modes" gap (the essay does not currently enumerate substrate-side failure modes for remote-runner cycles); the patch is to add a section that names this and any other known substrate quirks the doctrine should warn future cycles about.

**Why disposition `next-MCA` is correct (alternative considered: documentation patch).**
- The narrowest possible patch is a documentation addition to `BOX-AND-THE-RUNNER.md` (adding a "Known substrate failure modes" section listing this non-trigger behavior + any others). That would be a "patch-landed-in-cycle" disposition if filed inside cnos#426 — BUT Hard Rule 6 of cnos#426 ("No essay or README edits") forbids it.
- Therefore the patch must be a separate cycle. That makes it `next-MCA`.
- `no-patch` is wrong because the gap is real and recurring (any future tag-retarget cycle will re-encounter it; the doctrine should warn rather than silently repeat the surprise).

**First AC for the MCA follow-on:** "`docs/gamma/essays/BOX-AND-THE-RUNNER.md` has a 'Known substrate failure modes' section (or equivalent) that explicitly names tag-force-update non-trigger as a known GH Actions behavior, with a recommendation that remote-runner cycles relying on `on.push.tags` cascade after a tag force-move should explicitly trigger the downstream workflow (e.g., via `workflow_dispatch` or direct `gh workflow run`) rather than rely on the substrate cascade."

**Recommended additional scope hints for the MCA follow-on** (not binding):
- The section could also enumerate other known substrate quirks: (a) the second-firing of `on.push.paths` triggers when a self-deleting workflow's deletion commit fires the same path-filter (cosmetic noise, not a correctness issue, observed in cnos#425's run); (b) any GH Actions quirks around `actions/checkout@v4` token scoping that cycle authors should know about.
- The MCA cycle is small (single essay edit + a paragraph or two in §"What this means for cnos" or a new §"Known substrate failure modes"); should fit in a single γ+α+β-collapsed-on-δ cycle.

## Protocol-gap signals (across receipt-stream)

This cycle files 2 `cdd-tooling-gap` findings; neither is a `cdd-protocol-gap`. Looking across the receipt stream (per `cnos.handoff/skills/handoff/receipt-stream/SKILL.md`):

- The cnos#425 cdd-iteration filed 1 `cdd-protocol-gap` finding (boundary-model incomplete — remote-runner delegation needed naming). That finding's disposition was `patch-landed-in-cycle` (the doctrine + skill section + first artifact + first receipt all landed atomically).
- This cycle's 2 `cdd-tooling-gap` findings are the *first downstream-substrate findings* the cnos#425 doctrine surfaced. The pattern is: a new primitive's first use exposes substrate quirks that were previously latent; the second use's cdd-iteration is where they get filed as findings. This is healthy ε behavior — the new doctrine is doing its job by making substrate-side cracks legible.
- ε should track whether subsequent remote-runner cycles continue to surface substrate-side findings, and whether the doctrine's "Known substrate failure modes" section (F2's patch) is updated as new quirks are discovered.

## Non-findings (worth recording)

- **Doctrine inheritance pattern established.** This cycle ships only artifact + receipt; doctrine (essay + skill §8) is inherited from main from cnos#425. The pattern (first use: doctrine + artifact + receipt atomically; second-and-later use: artifact + receipt atomically, doctrine inherited) is a *structural pattern* this cycle establishes. ε does not file this as a separate finding because the pattern is implicit in `delta/SKILL.md §8.2`'s authoring-order rule (it says "with or before"; the second-use case is the "before" branch of the disjunction). But ε notes here that the pattern is a candidate for explicit naming in `delta/SKILL.md §8.2` or `BOX-AND-THE-RUNNER.md`'s "Known patterns" section if the second-use case recurs.

- **Idempotency of remote-runner moves is a healthy invariant.** cnos#425's `git tag -f` is idempotent (running it twice produces the same tag, modulo the no-op second-firing path-trigger). cnos#426's `softprops/action-gh-release@v1` is also idempotent (running it twice updates the release in place with the same body). Both remote-runner cycles to date have idempotent target actions. This is a positive structural property: idempotent remote-runner moves are easier to reason about than non-idempotent ones (a failure-and-retry produces the same correct end state). ε observation, not a finding; future doctrine guidance could prefer idempotent target actions where feasible.

- **Worktree-switch pattern for at-tag checkouts.** cnos#426's self-delete step needs `git fetch origin main` + `git checkout main` before `git rm` because step 1's checkout uses `ref: '3.82.0'` (the working tree is detached on the tag, not on main). cnos#425's workflow did not need this switch (its checkout was implicitly at main). The pattern difference is: checkout-at-non-main-ref implies worktree-switch-before-self-delete. ε observation, not a finding; could be added to F2's doctrine patch as part of the "Known patterns" section.

- **Checkout-at-tag as a release-body-source pin.** This cycle's most load-bearing correctness mechanism is `ref: '3.82.0'` in step 1. The motivation: main HEAD has post-tag content (CELL-OF-CELLS from cycle/424, BOX-AND-THE-RUNNER from cycle/425) that would publish the wrong body if the workflow used `release.yml`'s `ref: main` pattern. This makes the release boundary's "what body does the release actually contain" question answerable structurally rather than by hoping main hasn't drifted. ε observation, not a finding; could be added to `release-effector/SKILL.md` as a pattern for any release-fix one-shot that publishes notes (rather than re-running the full build pipeline).

- **β-independence collapse named, not papered over.** This cycle ran as `γ+α+β-collapsed-on-δ`, inheriting cnos#425's first-use precedent (which itself inherited cycle-414 / cycle-424). The receipt is closed-as-degraded at the structural-independence axis. `beta-closeout.md` names the collapse explicitly with the cnos#425 / cycle-414 / cycle-424 precedent inheritance and the dispatch authorization. The discipline of naming the collapse rather than hiding it is good ε hygiene and inherited from prior cycles.

## Verdict

`tooling_gap_count: 2`. `protocol_gap_count: 0`.

Two `cdd-tooling-gap` findings (F1: release.yml build-job git-auth failure on tag-push; F2: GH Actions `on.push.tags` non-trigger on tag force-update), both with disposition `next-MCA`. Each has a first AC specified. No `patch-landed-in-cycle` or `no-patch` dispositions.

INDEX.md row: `findings=2, patches=0, MCAs=2, no-patch=0`.

Filed by ε@cnos (γ+α+β-collapsed-on-δ) on 2026-05-24.
