# Design: Move `dist/packages/` out of git

**Issue:** #266
**Mode:** MCA
**Active Skills:** eng/go, design, cdd/design, eng/process-economics
**Engineering Level:** L7 — changes the system boundary so the friction class disappears

## Problem

`dist/packages/` is committed into the repo. `cn build` regenerates its contents from `src/packages/`. Because both the source and the build output live in git, every PR that touches packaged content races every concurrent main-side edit to `src/packages/`:

1. α rebases, runs `cn build`, commits refreshed dist, pushes.
2. CI goes green against the PR's merge-base.
3. A main-side edit lands that touches a packaged file (even a one-line SKILL.md tweak).
4. The PR's committed `dist/packages/checksums.txt` and `index.json` now diverge from what `cn build` would produce on the post-merge tree — I3 fails on a fresh rerun.
5. α rebases again. Goto 2.

The I3 check (`.github/workflows/coherence.yml` `dist-source-sync`) exists specifically because committed dist/ can drift from source. It does not cause the race — it only makes the race visible.

Evidence (from #266): PR #265 hit D1 twice in one cycle. Rebase onto `07f021f` went green; a one-line `alpha/SKILL.md` edit at `22c0be0` re-triggered D1; rebase onto `22c0be0` went green only because main stayed still for ~30 minutes. Merged at `9a1c7df`.

The race is structural, not discipline. Post-PR #264 (deterministic tarballs) and `da92e44` (Go 1.24 pin), the rebuild is a pure function of source — meaning the committed dist/ is redundant information, and the only way it can disagree with source is if main moves.

## Constraints

- **DESIGN-CONSTRAINTS.md §2.2 (Source / artifact / installed clarity):** "What is authored (`src/`), what is built (`dist/`), and what is installed (`.cn/vendor/`) are always distinct and explicit." Committing `dist/` muddies this — `dist/packages/*.tar.gz`, `index.json`, and `checksums.txt` are authored-in-git while declaring themselves build outputs.
- **BUILD-AND-DIST.md §Migration steps** (already committed, draft): explicitly names `Add dist/ to .gitignore` as step 2 of the planned target state. The current committed-dist/ state is a transitional deviation, not the intended end state.
- **Determinism is a precondition, not a consequence.** PR #264 (`8265fe0`) + Go 1.24 pin make `cn build` a pure function of source. Without determinism, removing committed dist/ would lose release-time verifiability; with it, committed dist/ adds no verification the source itself doesn't already provide.
- **Release pipeline already attaches tarballs as GitHub release assets** (`release.yml` lines 108–147). The tarballs are already fetchable post-release independent of what is or isn't committed.
- **`cn deps restore` supports both HTTP and relative-path URLs** (`restore.go` lines 217–231). The local-dev flow (`cn build` → local `dist/packages/` → `cn deps restore`) works off an on-disk index, which `cn build` regenerates. Nothing in the restore path requires dist/ to be committed — only to exist on disk at restore time.
- **One source of truth (§1).** Source is in `src/packages/`. The committed dist/ is a second derived surface asserting the same facts. One must go.
- **No new tool or runtime dependency** — the fix must be implementable with the existing `cn build` / `cn deps restore` code paths.
- **Kata and CI workflows already run `cn build` before using dist/** (`ci.yml` `kata-tier2` line 115; `release.yml` line 111; `scripts/kata/05-build.sh`; package katas under `src/packages/cnos.kata/`). None of them depend on dist/ being pre-populated from git.

## Challenged Assumption

**The assumption being replaced:** *"The repository must carry a committed snapshot of every package artifact so that CI can verify dist/ matches source."*

This assumption was reasonable before PR #264: when tarballs were not byte-reproducible across machines/toolchains, the only way to answer *"what exactly is cnos.core-3.57.0?"* was to store the bytes themselves in git and gate them behind an I3 diff check. Committed dist/ functioned as an archived witness of a non-deterministic process.

Post-determinism, the witness is redundant: the authoritative answer to *"what is cnos.core-3.57.0?"* is already *"the bytes `cn build` produces from `src/packages/cnos.core/` at the commit tagged 3.57.0"*, and any agent can reproduce those bytes. The committed dist/ no longer adds verification; it only adds a second surface that can disagree with the first.

The replacement assumption: **the repository carries the source of truth (`src/packages/`) and the deterministic builder (`cn build`); the artifact is derived on demand, fetched from the release when published, or regenerated locally during development.**

This is the same assumption the code already holds about binaries (`cn` is not committed; `go build` produces it from source). Extending it to packages brings the artifact boundaries into alignment with DESIGN-CONSTRAINTS.md §2.2.

## Impact Graph

### Downstream consumers (read or walk toward `dist/packages/`)

| Consumer | Expects | Post-change |
|----------|---------|-------------|
| `.github/workflows/coherence.yml` `dist-source-sync` (I3) | Committed `checksums.txt` + `index.json` to diff against rebuild | **Delete** — drift is impossible once dist/ isn't committed; I1 (`cn build --check`) already validates source. |
| `.github/workflows/release.yml` line 111 | Runs `cn build` in-job → copies to `artifacts/packages/` | **Unchanged** — already produces dist/ from source at tag time. |
| `.github/workflows/ci.yml` `kata-tier2` line 115 | Runs `cn build` in-job before `cn deps restore` | **Unchanged** — does not depend on committed dist/. |
| `scripts/kata/05-build.sh` | Runs `cn build`, then checks `dist/packages/index.json` exists | **Unchanged** — builds before checking. |
| `scripts/kata/06-install.sh` | Expects `dist/packages/index.json` to exist; comment says "run 05-build.sh first" | **Unchanged** — already sequenced after build. |
| `src/packages/cnos.kata/katas/R{2,3,4}/run.sh` + `lib.sh` | Walk up from `$PWD` to find `dist/packages/index.json` | **Unchanged** — katas run under `scripts/kata/run-all.sh` which builds first. Katas already tolerate a freshly-built dist/. |
| `src/go/internal/restore/restore.go` `FindIndexPath` | Walks up from hubPath to `dist/packages/index.json` | **Unchanged** — the walk works identically whether dist/ is committed or freshly built; helpful error for missing index is a separate improvement. |
| `src/go/internal/pkg/pkg_test.go` line 76 | `t.Skip("dist/packages/index.json not present (run cn build first)")` | **Unchanged** — already defensive. |

### Upstream producers (write into `dist/packages/`)

| Producer | Writes | Post-change |
|----------|--------|-------------|
| `src/go/internal/pkgbuild/build.go` `Build`, `UpdateIndex`, `UpdateChecksums`, `Clean` | `*.tar.gz`, `index.json`, `checksums.txt` | **Unchanged** — still writes to local `dist/packages/`, which is now gitignored. |
| `.github/workflows/release.yml` line 111 `cn build` | Same files, then copies to `artifacts/packages/` for release assets | **Unchanged**. |

### Copies / embeddings of the same fact

- `dist/packages/checksums.txt` (committed) vs `cn build` output (derived) — **the duplication the change resolves**.
- `dist/packages/index.json` (committed) vs `cn build` output (derived) — **same**.
- `docs/alpha/package-system/BUILD-AND-DIST.md` §Migration steps 2 — already states "Add `dist/` to `.gitignore`" as a planned step. Needs update: strike "planned," mark the migration's dist-gitignore step as accomplished by this cycle, and remove any "currently committed" language.
- `README.md` line 168 — describes `cn build` → `dist/packages/` → `cn setup` flow. **Unchanged** — the flow still describes what `cn build` does locally.
- `docs/alpha/architecture/INVARIANTS.md` line 166 — lists `dist/packages/<name>-<version>.tar.gz` as an artifact class. **Unchanged** — the artifact still exists at that path; it's just not committed.
- `docs/alpha/package-system/PACKAGE-AUTHORING.md`, `PACKAGE-SYSTEM.md`, `DESIGN-227-distribution-pipeline.md`, `SELF-COHERENCE-227.md` — describe `cn build` producing dist/. **Unchanged** — production step still occurs; the commit step is what disappears.

### Authority relationships

- `src/packages/<name>/cn.package.json` is the authoritative source for package identity (name, version).
- `dist/packages/index.json` is derived from source + build logic. After this change, it remains derived but is no longer re-authored by being committed.
- Released tarballs on GitHub Releases become the authoritative artifact witness at release time. Between releases, the local `cn build` output is the witness for dev/CI flows.

### Rule changes rippling through all embeddings

The change touches one rule: *"dist/ is build output, never committed."* That rule must:
- appear in `.gitignore` (currently absent),
- be enforced by removing tracked files (currently 8 tarballs + 2 metadata files),
- stop being checked by the I3 drift gate (becomes vacuous / misleading),
- be reflected as accomplished-migration in `BUILD-AND-DIST.md`,
- not contradict anything else that already describes `cn build` as producing `dist/packages/` (it doesn't — producing ≠ committing).

## Proposal

### Types (what changes in the model)

Nothing in the type model changes. `pkg.PackageIndex`, `pkg.Lockfile`, `pkgbuild.Result` and the `cn.package-index.v1` / `cn.lock.v2` schemas remain byte-identical. The change is at the **persistence boundary**, not the **type boundary**.

### Persistence model — before / after

```
                      Before                                    After
src/packages/   ────► cn build ────► dist/packages/      src/packages/   ────► cn build ────► dist/packages/
 (git-tracked)                       (git-tracked)        (git-tracked)                       (gitignored)
                                            │                                                       │
                                            ├── committed                                           ├── used locally
                                            ├── diffed by I3                                        └── uploaded to Release at tag
                                            └── uploaded to Release at tag
```

### Concrete changes

1. **Remove `dist/packages/` from git index.**
   ```sh
   git rm dist/packages/*.tar.gz dist/packages/index.json dist/packages/checksums.txt
   ```
   8 tarballs + `index.json` + `checksums.txt` = 10 tracked files to remove.

2. **Add `dist/` to `.gitignore`.**
   New line (scoped broadly, per BUILD-AND-DIST.md §Migration step 2: "`dist/` is build output, never committed"):
   ```
   # Build output (cn build produces this; never committed — see DESIGN-266)
   dist/
   ```

3. **Delete the I3 `dist-source-sync` job from `.github/workflows/coherence.yml`.**
   The job compares committed `dist/packages/{checksums.txt,index.json}` against the rebuilt output. Once dist/ isn't committed, there is no "before" to diff against. The failure class the job detected (authored dist/ vs current source) cannot arise.
   Also remove `dist-source-sync` from the `notify.needs` list and the status aggregation logic.
   `coherence-build-check` (I1, `cn build --check`) remains and continues to gate structural validity of `src/packages/`.

4. **`release.yml` stays unchanged.** Line 111 already runs `cn build` in-job; lines 113–115 copy the produced artifacts into `artifacts/packages/`. The `2>/dev/null || true` tolerance handles the pre-build state. Nothing to modify.

5. **`ci.yml` `kata-tier2` stays unchanged.** Line 115 already runs `cn build` before `cn deps restore`.

6. **`restore.go` `FindIndexPath` stays unchanged.** It walks up to find `dist/packages/index.json`; that file exists once `cn build` has run locally. No code change needed for the MCA itself. (A minor UX improvement — a helpful error message when the index is missing, pointing to `cn build` — is appealing but is **out of scope for this cycle** per §Non-goals.)

7. **Update `docs/alpha/package-system/BUILD-AND-DIST.md`.** The migration section currently writes `Add dist/ to .gitignore` as a future step. Mark it as accomplished under the specific cycle (reference #266), and remove transitional language that implies dist/ is still committed. Keep the layout diagrams (they describe structure, not tracking).

### Source of truth (§1) after this change

| Fact | Canonical source | Derived |
|------|------------------|---------|
| Package identity (name, version, contents) | `src/packages/<name>/` | everything |
| Package tarball bytes | `cn build` at HEAD | local `dist/packages/*.tar.gz`; GitHub Release assets |
| Package SHA-256 | `cn build` output | `dist/packages/checksums.txt` (local); GitHub Release asset |
| Package index | `cn build` output | `dist/packages/index.json` (local); GitHub Release asset |

Every row collapses to a single canonical source. No row has two authoritative entries.

### Mechanical vs judgment boundary

- **Mechanical (enforced by the change itself):** removing committed dist/ + gitignoring `dist/` + removing the I3 job makes authored-dist drift *impossible* in this repo. No reviewer vigilance, no periodic audit, no discipline gate required.
- **Judgment (remains):** whether a package's source itself is coherent (I1, review); whether a release's tag points at the right commit (β release step); whether a hub's installed packages match its lockfile (runtime, doctor).

### Leverage

- Eliminates the rebase race class for the life of the repo. Every PR that touched packaged content under the old model paid the rebase tax; every PR that would have paid it stops paying.
- Removes an entire CI job (I3) from every push and every PR, cutting `cn build` runtime + diff-check overhead from the coherence workflow.
- Aligns dist/ with the already-accomplished treatment of `cn` binaries and node_modules: build output is never in git.
- Simplifies reviewer cognitive load — PRs stop carrying binary tarball diffs that no reviewer reads.
- Unblocks future changes to tarball internals (compression level, metadata, schema evolution) without triggering a flood of committed-dist-diff noise.

### Negative leverage

- **Fresh-clone bootstrap requires `cn build` before the first `cn deps restore`.** Today a contributor cloning the repo has dist/ ready. After the change, any command that walks up to find `dist/packages/index.json` and finds nothing will produce the existing "package index missing" error. This is strictly worse UX for the first local run. Mitigation: kata scripts and CI already build before install; the README already documents `cn build`; the `FindIndexPath` fallback already returns a predictable path. A more helpful error is deferred to a separate cycle to keep this MCA surgical.
- **Between-release artifact forensics becomes harder.** Today you can `git checkout <sha> && sha256sum dist/packages/cnos.core-3.57.0.tar.gz` to answer "what did cnos.core-3.57.0 look like at that commit?" After the change, the same question requires `git checkout <sha> && cn build && sha256sum dist/packages/cnos.core-3.57.0.tar.gz` — same answer, one extra step. This cost is only paid by forensic investigation, not normal dev/CI/release, and the determinism invariant guarantees the answer is identical.
- **The post-release record of "what bytes shipped" now lives only on GitHub Releases, not in git.** Mitigated: Release assets are immutable and include checksums; `release.yml` already publishes them.

### Simplest thing that closes the failure class

The race exists because two artifacts claim to describe the same fact (`src/packages/*` and committed `dist/packages/*`). The minimum fix is *remove one of them*. Removing source is impossible. Removing committed dist/ is the minimum. Everything else in this proposal is either a mechanical follow-on (I3 job, .gitignore) or a doc consequence (BUILD-AND-DIST.md).

## Alternatives Considered

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| **A. Move `dist/` out of git** (this proposal) | Eliminates the race by construction — nothing in-repo can drift. No new infrastructure. Returns the layout to the documented target state in BUILD-AND-DIST.md. Removes an entire CI job. Aligns dist/ with how binaries are already handled. | Slight bootstrap cost: fresh clones need `cn build` before first `cn deps restore`. Forensic rebuild adds one step. | **Chosen.** |
| **B. GitHub merge queue** | Re-runs CI against actual post-merge tree; I3 would fail pre-merge instead of post-merge. Doesn't require source changes. | Requires branch-protection config outside the repo diff (can't be shipped inside the PR). Doesn't eliminate the rebuild cost — just hides the race behind queue latency. Introduces throughput ceiling for a single-maintainer project that doesn't need one. Leaves the underlying duplication (committed dist/ vs source) in place — I3 would still be re-running the diff on every merge. A L6 fix (faster recovery) rather than an L7 fix (boundary change). | Rejected. Symptomatic, not structural. Preserves the incoherence the issue was filed to eliminate. |
| **C. Freeze main-side `src/packages/` edits while a PR touching packaged files is open** | Simplest to describe. No code or CI changes. | Social coordination, not structural. Scales only to 1–2 concurrent PRs. Imposes a "don't touch main" rule that CDD explicitly rejects as ritual. Failure class is still present, only suppressed by discipline. An L5 fix at best. | Rejected. Trades a race for a coordination gate without eliminating either. |
| **D. Keep committed dist/ but make I3 advisory** | Zero disruption. | Doesn't fix the race; it just stops surfacing it. The stated failure class (PR rebased on A, main edits land, dist stale) would still happen; we'd just ship with stale committed dist/. Defeats the point of having an I3 check at all. | Rejected. Would silence the signal the issue specifically says not to silence ("#Not in scope: Removing the I3 check"). |
| **E. Dist-manifest-only in git** (keep `index.json` + `checksums.txt`, drop `*.tar.gz`) | Partially addresses binary-artifact-diff noise. | Doesn't eliminate the race. Any main-side edit to `src/packages/<name>/*` still changes the tarball's SHA, which still changes checksums.txt and index.json. The race is identical in failure behavior to the current state; only the diff size shrinks. | Rejected. Half-measure that preserves the failure class. |

Choice rationale: A is the only candidate that makes the failure class *structurally impossible*, which is the definition of an L7 MCA per CDD §9.1.

## Process Cost / Automation Boundary

- **One-time adoption cost:** 10 files removed from git, `.gitignore` rule broadened from `dist/bin/` to `dist/` (+2 lines net, with explanatory comment), 1 CI job removed, ~10 lines of doc updated. No ongoing cost — the change removes process rather than adds.
- **Ongoing cost:** zero. The gitignore rule is passive. `cn build` already runs in CI and locally in the existing flows that need dist/.
- **Consumer of the removed I3 job:** nobody post-change. I3's function (catching authored-dist/ drift) dissolves because authored-dist/ stops existing.
- **Automation boundary:** `.gitignore` + commit removal is 100% mechanical. No discipline, checklist, or reviewer gate is required to enforce "don't commit dist/" — git itself silently drops the path.
- **Sunset criterion:** none needed — this is a removal, not an addition. The failure class either re-emerges (nothing in-repo to drift) or it doesn't.

## Non-goals

- **Changing the tarball format or hashing algorithm.** PR #264's determinism work is the precondition for this change, not an overlap with it.
- **Removing or weakening I1 (`cn build --check`).** I1 gates source validity; it is not affected by where dist/ lives.
- **Adding a helpful "index missing, run `cn build` first" error to `cn deps restore`.** Desirable but orthogonal — the current behavior (missing-index error) is correct, just terse. File separately if warranted.
- **Migrating other build outputs** (e.g. `cn` binary handling, `node_modules/`). Those are already handled correctly; no change needed.
- **Publishing packages to a remote registry** (future #227-class work). This change is a precondition, not a prerequisite.
- **Changing `FindIndexPath` walk semantics, `cn deps restore` index URL resolution, or lockfile schema.** All contracts stay byte-identical.
- **Enforcing the gitignore rule with a CI check.** Git's own tracking state is sufficient; adding a "no committed dist/" CI gate would be process-economics debt for a failure that gitignore already prevents.

## File Changes

### Delete (from git tracking; local copies regenerable by `cn build`)

- `dist/packages/cnos.cdd-3.56.2.tar.gz`
- `dist/packages/cnos.cdd-3.57.0.tar.gz`
- `dist/packages/cnos.cdd.kata-0.3.0.tar.gz`
- `dist/packages/cnos.core-3.56.2.tar.gz`
- `dist/packages/cnos.core-3.57.0.tar.gz`
- `dist/packages/cnos.eng-3.56.2.tar.gz`
- `dist/packages/cnos.eng-3.57.0.tar.gz`
- `dist/packages/cnos.kata-0.2.0.tar.gz`
- `dist/packages/index.json`
- `dist/packages/checksums.txt`

(Total: 10 files, per `git ls-files dist/packages/` at HEAD `3e850f3`. Implementation uses `git rm` against that list.)

### Edit

- `.gitignore` — add `dist/` block with DESIGN-266 reference comment.
- `.github/workflows/coherence.yml` — remove the `dist-source-sync` job; remove from `notify.needs`; update the status-aggregation shell to drop the `R2` variable and the dist reference in the status icon logic.
- `docs/alpha/package-system/BUILD-AND-DIST.md` — §Migration steps: mark step 2 (`Add dist/ to .gitignore`) accomplished, with a DESIGN-266 link; remove any remaining transitional "currently committed" language; leave the §Directory Structure and §Flow diagrams unchanged (they describe structure, which is preserved).

### Create

- `docs/alpha/package-system/DESIGN-266-dist-out-of-git.md` — this file (primary branch artifact).

### Unchanged (peers audited; no edit required)

- All Go code under `src/go/` — no call site assumes dist/ is committed.
- `release.yml` — already rebuilds in-job.
- `ci.yml` — already rebuilds in-job.
- `scripts/kata/*.sh` — already sequence `cn build` before consumption.
- `src/packages/cnos.kata/**` — already builds before walking up to find index.
- `README.md`, `INVARIANTS.md`, `PACKAGE-AUTHORING.md`, `PACKAGE-SYSTEM.md`, `DESIGN-227-distribution-pipeline.md`, `SELF-COHERENCE-227.md` — describe `cn build` producing dist/; still true post-change (production is local, not git-committed). No edit needed.

## Acceptance Criteria

- [ ] **AC1 (issue AC1):** This design doc exists at `docs/alpha/package-system/DESIGN-266-dist-out-of-git.md`, names the chosen MCA (A), lists rejected alternatives (B, C, D, E) with per-option rationale. Verification: `test -f docs/alpha/package-system/DESIGN-266-dist-out-of-git.md && grep -c '^## Alternatives Considered' …` returns 1, and the Alternatives table has rows for A/B/C/D/E.
- [ ] **AC2 (issue AC2):** The specific failure mode in #262 ("PR rebased on main-at-open + main-side edit to any `src/packages/` file → red I3 on merge commit") is structurally impossible. Verification:
  - (a) `git ls-files dist/packages/` on HEAD returns empty.
  - (b) `.gitignore` contains `dist/` (or `dist/packages/`).
  - (c) `.github/workflows/coherence.yml` contains no job named `dist-source-sync` and no `dist/packages/` reference.
  - (d) After the PR merges, any subsequent PR can edit any `src/packages/**` file on main without producing a new file diff in `dist/packages/` on any other open PR (because there is no `dist/packages/` to diff).
- [ ] **AC3 (issue AC3 — partial, α-side):** The design doc names the expected engineering level (L7) and states the post-release question β must answer: *"After N subsequent cycles touching `src/packages/`, did any rebase-race recurrence occur?"* Verification: `grep -E 'L7|post-release' docs/alpha/package-system/DESIGN-266-dist-out-of-git.md` returns hits in the Challenged Assumption, Engineering Level header, and CDD Trace. (β's post-release assessment closes this AC per CDD §9.1; α cannot verify it from the branch.)
- [ ] **AC4 (impact-graph coverage, per cdd/design §2.5):** Every consumer/producer enumerated in §Impact Graph is either updated or explicitly marked "unchanged" with a reason. Verification: the §Impact Graph tables cover `coherence.yml`, `release.yml`, `ci.yml`, `scripts/kata/*`, `src/packages/cnos.kata/**`, `restore.go`, `pkgbuild/build.go`, `pkg_test.go` — confirmed by cross-reference against `grep -rn 'dist/packages' --include='*.go' --include='*.sh' --include='*.yml' --include='*.md'`.
- [ ] **AC5 (build-and-dist.md alignment):** `docs/alpha/package-system/BUILD-AND-DIST.md` no longer describes the dist-in-git state as current. Verification: `grep -n 'dist/ to .gitignore' docs/alpha/package-system/BUILD-AND-DIST.md` shows the step marked accomplished with a #266 reference; no lines imply dist/ is currently tracked.
- [ ] **AC6 (kata still passes post-change):** `scripts/kata/run-all.sh` and the kata-tier2 sequence (`cn build` → init hub → `cn deps lock` → `cn deps restore`) succeed with an empty dist/ at start. Verification: local `scripts/kata/run-all.sh` exits 0.

## Known Debt

- **Bootstrap UX for fresh clones.** A contributor who clones and immediately runs `cn deps restore` in a hub will get an index-missing error. The error message is accurate but not self-explanatory. A UX polish pass (helpful error that says "run `cn build` first") is deferred.
- **`FindIndexPath` fallback path.** Currently returns `filepath.Join(hubPath, "dist", "packages", "index.json")` when no index is found by walking up. That's a file that doesn't exist — the subsequent `ReadPackageIndex` call will fail with a clear error. Tolerable, but not elegant. Deferred.
- **Documentation churn in the package-system doc cluster.** Several docs describe the pipeline with `dist/packages/` in prose. None of them assert dist/ is committed, and they all remain accurate about `cn build`'s output location. No text changes are required, but a later consolidation pass may want to make the "build-output only" property explicit in one or two places.
- **Post-release validation (AC3).** β must observe at least one subsequent cycle with concurrent src/packages/ edits to confirm the race doesn't recur. The assessment structure (CDD §9.1) already requires this.

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | — | — | Read issue #266 fully; read linked evidence (PR #265 commits `567e195`, `eb9a6fe`, `da92e44`); read `dist/packages/` layout, `coherence.yml` I3, `release.yml`, `ci.yml`, `restore.go`, `BUILD-AND-DIST.md`, kata scripts. Observed that BUILD-AND-DIST.md §Migration already commits to `dist/` being gitignored. |
| 1 Select | — | — | Selected gap: committed dist/ + main-side src/packages/ edits = rebase race class; candidate MCAs A/B/C filed; α to pick one and implement. |
| 4 Gap | this doc §Problem | — | Named incoherence: dist/ is build output but is also committed, creating a second authoritative surface that can disagree with source whenever main moves. |
| 5 Mode | this doc §header | eng/go, design, cdd/design, eng/process-economics | MCA, L7 (boundary change eliminates friction class). Active skills loaded before implementation. |
| 6 Artifacts | this doc + follow-on diff | eng/go (schema/audit §2.12), design (source-of-truth §3.2, register §3.6), eng/process-economics (§2.6 lightest option, §2.13 challenge) | Design complete; implementation = `git rm` + `.gitignore` + coherence.yml edit + BUILD-AND-DIST.md edit. No test changes needed (no code paths change). |
| 8 Review | — | review (pending) | β review pending. |
| 9 Gate | — | — | Pending. |
| 10 Release | — | — | Pending. |
