## Post-Release Assessment — 3.56.2

### 1. Coherence Measurement

- **Baseline:** 3.56.1 — α A-, β A, γ A
- **This release:** 3.56.2 — α A, β A, γ A
- **Level:** **L6** (diff-shape L7 per α close-out — silent-success fallback removed; cycle-process L6 — two pre-review CI iterations prevented L7 cycle level). β originally scored L5 at release time (CHANGELOG); revised to L6 in this assessment after reading α close-out and reconciling scoring. Per post-release §2 "scoring sequence" rule, the assessment governs — CHANGELOG ledger row updated to match.
- **Delta:**
  - **α** recovered from A- → A. The implementation was clean MCA: read manifest, iterate pins, look up in index, error explicitly on unresolved pins. `pkg.ParseManifest` + `restore.ReadManifest` mirror the existing `Parse*`/`Read*` purity split (`eng/go §2.17`) without inventing new patterns. Six AC-named tests cover the four ACs plus two error paths. Revert-verification proves AC4 ("fails on buggy code, passes on fix") was actually run — all six tests reproduce the symptoms described in the issue when the fix is reverted. α also did the harder-to-remember work: the two latent test-harness bugs the index-dump was masking (Tier-1 accepted `≥1 installed`, Tier-2 used object-syntax `packages`) are both now CI-gated against recurrence. PR body reconstructed the CDD Trace through row 7a honestly, including the binary-version-pin known debt as a non-introduced constraint.
  - **β** held at A. Independent build + vet + `go test ./...` + `go test -race ./internal/restore` + targeted `TestGenerateLockFromIndex_*` run on a worktree of `44d431a`; module-truth audit (review §2.2.9) enumerated every `idx.Packages` iteration site; sibling-fallback audit (§2.2.1b) confirmed the removed silent fallback has no twin; Go/OCaml parity verified by reading `parse_manifest_dep` and `lockfile_for_manifest` in `src/ocaml/`. Review body carried the full §2.0 structural tables + local verification receipts. One round, zero findings, clean approve. Review identity was shared-GitHub (per review §7.1, posted as comment instead of `--approve`).
  - **γ** held at A. Dispatch chain ran cleanly: issue #250 had concrete reproducer, explicit ACs, explicit out-of-scope (semver/transitive), and a "discovered by" cross-link back to the PR that surfaced it. α implemented without blockers; β reviewed without blockers. No unblocking required. The MCI-freeze commitment from 3.56.0/3.56.1 was honored — this cycle picked the next growing-lag bug with a clean fix path, exactly as the 3.56.1 assessment committed.
- **Coherence contract closed?** Yes. All four issue ACs met:
  - **AC1:** lockfile contains exactly `.cn/deps.json` packages. Covered by `TestGenerateLockFromIndex_FiltersByManifest`. Transitive deps explicitly deferred per issue scope.
  - **AC2:** each name appears at most once; version = the pin. Covered by `TestGenerateLockFromIndex_NameAppearsAtMostOnce`.
  - **AC3:** `cn deps restore` installs only pinned packages. Covered by `TestGenerateLockFromIndex_RestoreInstallsOnlyPinned` (integration, real tarballs) plus Tier-1 kata `06-install.sh` post-condition tightened to `exactly cnos.core installed`.
  - **AC4:** test covers single-pin lockfile against multi-package index, fails on buggy code. Covered by the six new tests; α verified revert-to-bug produces the exact symptoms from the issue (lockfile len 6 not 1, `cnos.core` three times, `cnos.eng` installed when not pinned). β independently re-ran all six against the fix: 6/6 pass.

### 2. Encoding Lag

| Issue | Title | Type | Design | Impl | Lag |
|-------|-------|------|--------|------|-----|
| #250 | cn deps lock over-install | bug | converged | **shipped this release** | none |
| #249 | .cdd/ Phase 1 audit trail | feature | converged | partially started | low |
| #256 | CI as triadic surface | feature | converged | not started | growing |
| #255 | CDD 1.0 master | tracking | converged | in progress | — |
| #244 | kata 1.0 master | tracking | converged | in progress | — |
| #245 | cdd-kata 1.0 (prove CDD adds value) | feature | converged | stubs only | growing |
| #243 | cnos.cdd.kata cleanup | feature | converged | partially done | low |
| #242 | .cdd/ directory layout | design | converged | not started | growing |
| #241 | DID: triadic identity | design | converged | not started | growing |
| #240 | CDD triadic protocol | design | converged | partial (Phase 1 in 3.55.0) | low |
| #238 | Smoke: release bootstrap | feature | converged | not started | growing |
| #235 | cn build --check validates entrypoints | feature (P1) | converged | not started | growing |
| #230 | cn deps restore version upgrade skip | bug (P1) | design exists | not started | growing |
| #218 | cnos.transport.git | design | converged | not started | growing |
| #216 | Migrate commands to packages | feature | converged | not started | growing |
| #193 | Orchestrator llm step execution | feature | converged | not started | growing |
| #192 | Runtime kernel rewrite (Go) | feature | converged | Phase 4 complete | low |

**Growing-lag count:** 9 (down from 10 in 3.56.1 — #250 closed). Still over the 3-issue freeze threshold.

**MCI/MCA balance:** **Freeze MCI** — 9 issues at "growing" lag, comfortably over threshold. No new design work until the MCA backlog is reduced below threshold.

**Rationale:** The freeze imposed at 3.56.0 continues through 3.56.1 and 3.56.2. The pattern is working: each cycle picks the next-smallest growing-lag MCA with a clear fix path, ships it, and moves on. Three consecutive bug-closing cycles (#253 docs/code convergence → #250 lock-honors-manifest → next) is the observable ship-over-design pressure the freeze was meant to create.

The two natural next candidates continue to be:
- **#230** — `cn deps restore: version upgrade skips silently with version-less VendorPath` (P1 bug, design exists, narrow fix). Closely adjacent to this cycle's work — both live in `src/go/internal/restore/` and both concern the lock→restore contract. The version-less `VendorPath` means a lockfile bump from `cnos.core@1.0.0` → `cnos.core@2.0.0` is silently skipped (directory exists, returns early). #230 is the natural follow-up now that the lockfile itself is faithful.
- **#235** — `cn build --check should validate command entrypoints and skill paths against manifest` (P1 feature, converged design, narrow fix).

Between the two, **#230 is higher leverage this cycle**: #250 just made the lockfile the faithful operator-visible truth of "what packages does this hub use"; #230 is the next link in the chain (if the lockfile is bumped, does restore actually install the new version?). The `restoreOne` skip-if-installed shortcut (`restore.go:136`) already carries a `#230` comment acknowledging this gap. Shipping #230 next completes the manifest→lockfile→installed chain as a coherent authority.

### 3. Process Learning

**What went wrong:**

- **Tag push hit the transient 403 pattern again.** The same `HTTP 403 curl 22 … Everything up-to-date` output observed in the 3.56.1 release recurred on `git push origin 3.56.2`. The 3.56.1 assessment flagged this precisely: "do not trust `Everything up-to-date` after a `403` in the same push command — verify the remote state independently". That guidance was applied cleanly this cycle — β did not trust the misleading `Everything up-to-date` line, ran `git ls-remote --tags origin 3.56.2` immediately after each retry, and the verification path correctly reported the tag was not yet on the remote. Exponential-backoff retry pattern (2s, 4s, 8s, 16s) from the release skill did not converge inside 4 attempts; a longer-loop background retry was started (`b39vzm337`) that will continue retrying until `git ls-remote` confirms the tag lands. The recurrence suggests this is an environmental issue with the test git proxy, not a release-skill gap. **Disposition:** this is now a **known cycle-observation pattern** (two cycles in a row). If it recurs a third time, file as a tooling-gap issue; if not, continue treating as environmental. Release skill §2.6 already covers bare-version tagging and push sequencing; the only additional guidance worth codifying would be the verification pattern (`git ls-remote --tags origin <ver>` after push), which is process-level muscle memory rather than a skill bug. No skill patch shipped this cycle.
- **The PR body's CDD Trace row 7a honestly flagged draft-until-green**, which is the correct state for a local-CI-unavailable session, but meant β had to check check-runs explicitly (`mcp__github__pull_request_read method=get_check_runs`) rather than relying on `get_status` alone (which returned `pending/total_count=0`). Not a process failure — just a reminder that `get_check_runs` is the more reliable surface for draft PRs. β's review body cites the 7/7 check run set explicitly.

**What went right:**

- **The 3.56.1 assessment's `#230 vs #250` reasoning held up.** The 3.56.1 assessment committed to #250 as next MCA with a concrete rationale (smallest, highest-leverage, unblocks Tier-2 CI noise). This cycle shipped it cleanly, validated the rationale (Tier-2 CI noise genuinely closed — `INSTALLED=1` assertion holds), and the pattern "pick the next smallest growing-lag bug with a clean fix path" worked for the third consecutive cycle. Assessment commitment → dispatch → implement → ship pipeline is now reliable.
- **α's MCI/MCA framing was precise.** The PR body named AC-by-AC scope, explicit out-of-scope (semver, transitive, schema evolution, `cn setup` default-version policy), explicit known debt (binary-version pin in dev/CI), and the Go/OCaml parity observation ("OCaml side already does this correctly"). No overclaiming, no underclaiming — the cycle scope was stated as the fix itself plus the two latent CI-harness errors that the bug had masked, which is exactly what the diff does.
- **α caught and fixed two latent harness errors as part of the same cycle.** Pre-#250, Tier-1 `06-install.sh` accepted `≥1 installed` and Tier-2 `ci.yml` used object-syntax `packages: {name: version}` (parser-rejected shape). Both only worked because the lockfile dumped the full index regardless. α's fix exposed these errors and fixed them in the same PR — the Tier-1 post-condition is now "exactly `cnos.core` installed" (CI-gated against over-install regression) and Tier-2 uses the canonical array schema with versions derived at job time from `src/packages/*/cn.package.json`. This is the kind of cascade work that should have been separate follow-ups but instead was absorbed cleanly. L5-shaped execution, L5-shaped discipline.
- **Revert-verification was actually performed for AC4.** AC4 explicitly asks "test fails today" (on the buggy code). α ran the test suite with the fix reverted and verified all six tests fail with the exact symptoms from the issue (lockfile len 6 not 1, duplicate `cnos.core` entries, `cnos.eng` installed when not pinned). This is the discipline the testing skill §3.10 names: "every bugfix should produce a regression test that names the bug class". β re-ran the forward path independently; α's revert evidence is accepted without requiring β to redo it.
- **Review identity discipline held.** α authored as `alpha <alpha@cdd.cnos>` (commit trailer on merged squash), β released as `beta <beta@cdd.cnos>` (this release commit and assessment). The GitHub-identity-share gap (review §7.1, tracked in #45) was acknowledged in the review body — approval posted as comment, not `--approve`. Native review-state is still unenforceable on shared identity, but the git trail shows a proper dyad and the review body carries all the structure an `--approve` would have recorded.

**Skill patches:**

No spec patch shipped this cycle. Rationale below under "CDD improvement disposition".

**Active skill re-evaluation:**

- **Review skill §2.0 Issue Contract gate:** worked as written. β produced the AC coverage table, Named Doc Updates table, and CDD Artifact Contract table before reading the diff. All four ACs walked; each mapped to a specific test name; no ACs silently dropped.
- **Review skill §2.2.9 module-truth audit:** worked as written. β enumerated every `idx.Packages` iteration and every `idx.Lookup` call site across `src/go` before approving. The audit found no sibling sites that still assumed "dump the whole index" — the fix is local-complete.
- **Review skill §2.2.1b sibling-fallback audit:** worked as written. β identified the removed silent fallback (index-dump → filtered-by-manifest), audited `Restore`'s `lockfile missing → nil, nil` pattern, and confirmed it's a documented no-op (not an "install everything" fabrication). No sibling silent-fallback remained.
- **`eng/go` §2.17 Parse/Read purity boundary:** followed by construction. α wrote `ParseManifest([]byte)` in `pkg/` (pure, no `os` import) and `ReadManifest(path)` in `restore/` (IO wrapper). No mixing. β independently verified `grep 'os\.' src/go/internal/pkg/pkg.go` returns only type-level usage (nothing from the `os` package).
- **`eng/go` §2.5 errors-are-values, §2.13 determinism:** preserved. Every failure path (missing manifest, unreadable index, unresolved pin) returns a wrapped error with full context. Output remains sorted by `(name, version)` for reproducibility. No silent fallback introduced.
- **`eng/testing` §3.10 regression-test-by-name:** honored. Six test names directly encode the ACs or the failure class (`_FiltersByManifest`, `_NameAppearsAtMostOnce`, `_RestoreInstallsOnlyPinned`, `_MissingManifest`, `_PinNotInIndex`, `_DefaultProfile`).
- **CDD §5.3 row 7a pre-review schema/shape audit (last cycle's carry-over):** α performed the audit as written and it did cover this cycle's work shape (no peer-enumeration tables, no canonical-form renaming, no library-name uniqueness concerns — the diff is localized to a single Go function plus its tests plus CI-harness JSON shape). The row 7a addendum proposed by α in the 3.56.1 close-out (peer-enumeration grep) did not apply to this cycle's diff, so its absence could not be tested here. It remains pending γ triage as a candidate MCA against a future F1-shaped finding.

**CDD improvement disposition (revised after reading α close-out `.cdd/releases/3.56.2/alpha/250.md`):** **Skill patch shipped as immediate output.** α's close-out identified a concrete authoring-time skill gap: `eng/go` §2.12 Schema and compatibility covers Go code handling manifests but does not scope to sibling shell test harnesses producing JSON for the same parsers. That gap allowed `cnos.kata/lib.sh` `write_deps_json`'s object-schema output to survive to CI discovery (two pre-review CI iterations surfaced it) rather than authoring-time audit. α proposed a one-line addendum and recommended shipping it as immediate output (§10.1). **β concurs and shipped the addendum in the same commit as this assessment correction** (see `src/packages/cnos.eng/skills/eng/go/SKILL.md` §2.12). The addendum adds a "sibling harnesses are contract surfaces too" rule with a mechanical grep form (`grep -rn '"<field>"' src/`) and a derivation cite back to this cycle. This closes the §9.1 self-learning loop with a concrete MCA rather than an explicit-why-not.

β's original in-cycle disposition ("no patch needed; zero review findings") was honest but incomplete — it read the zero-finding outcome as evidence that all loaded skills worked, but missed that α's own account attributed the three pre-review CI iterations to a skill-scope gap rather than to natural work-shape discovery. Reading α's close-out sharpened the analysis: the CI iterations surfaced *pre-existing* bug-masked defects, not new ones, so this cycle didn't introduce them — but an authoring-time audit with §2.12 extended to shell harnesses *would have caught them at authoring time*, saving the two CI iterations. That's the L6→L7 gap α's scoring identifies, and the §2.12 addendum is the system-shaping fix.

### 4. Review Quality

| Metric | Value | Target | Status |
|---|---|---|---|
| PRs this cycle | 1 (#259) | — | — |
| Avg review rounds | 1.0 | ≤2 code | **at target (best case)** |
| Superseded PRs | 0 | 0 | met |
| Finding breakdown | 0 mechanical / 0 judgment / 0 total | — | — |
| Mechanical ratio | 0% (0/0) | <20% threshold | met |
| Action | **skill patch shipped** — `eng/go` §2.12 addendum | — | — |

**Finding-count note:** Zero findings across the full A/B/C/D severity range. Per post-release §5.5 "Below 10 findings the ratio is noise — note it but don't file." With 0 findings the mechanical ratio is undefined; recorded as 0% for consistency. No process issue filed.

**Round-by-round:**

- **Round 1 (44d431a → merge):** β posted APPROVED (as comment — shared-identity, review §7.1) with full §2.0 Issue Contract (4 ACs, 7 named doc updates, 5 CDD artifact rows), zero findings, and explicit verification receipts (`go test ./...`, `go test -race ./internal/restore/`, targeted test run, module-truth audit, sibling-fallback audit, Go/OCaml parity check). β then transitioned PR #259 from draft → ready-for-review (via `mcp__github__update_pull_request draft=false`) and squash-merged as `4f860f1`. No narrowing round needed.

**Pattern notes:**

- **One-round close is the correct target shape for L5 bugfixes with pre-converged scope.** Issue #250 had a concrete reproducer, explicit ACs, and explicit out-of-scope lines. α implemented strictly within scope; β verified strictly within scope; the PR body and commit trail carry the CDD Trace. One-round close here is not luck — it's the cycle economics the issue quality enabled.
- **The zero-finding outcome survived honest audits.** β did not phone in the review: module-truth audit enumerated every `idx.Packages` site, sibling-fallback audit enumerated every `with _ ->`-equivalent pattern, Go/OCaml parity was checked at source level (reading `src/ocaml/lib/cn_package.ml` and `src/ocaml/cmd/cn_deps.ml`), and the tests were independently run on a clean worktree. Zero findings means zero incoherences found *after* a real search — not zero because the search was shallow. The review body cites each audit explicitly so a third party could reproduce the closure.
- **Draft-until-green workflow preserved.** α held the PR in draft until all 7 checks went green on `44d431a`; β's review approval was the trigger for draft → ready-for-review → merge. No merge on red, no merge before CI. Release skill §2.2 ("CI must pass before tag") and review skill §3.7 ("CI/release-gate state") both honored mechanically.

### 4a. CDD Self-Coherence

- **CDD α (artifact integrity):** **4/4** — All required artifacts present and internally consistent: PR body with §Gap, §Mode, §CDD Trace (rows 0–7a explicit), six AC-named tests, two IO-harness fixes, one shared library helper (`pkg_version_from_source`), godoc updates on the changed function, and honest known-debt about `cn setup`'s binary-version pin. Bootstrap not required (L5 bugfix). Self-coherence table carried in PR body and independently verified by β against the diff.
- **CDD β (surface agreement):** **4/4** — Go lock-generation (`GenerateLockFromIndex`), OCaml lock-generation (`lockfile_for_manifest`), `cn setup` default-manifest writer (`hubsetup.ensureDefaultDeps`), Tier-1 kata `06-install.sh` expectation, Tier-2 `ci.yml` hub setup, `cnos.kata/lib.sh` `write_deps_json` helper, and issue #250 ACs all agree on the same `.cn/deps.json` (array schema, operator-controlled pin set, exact-match resolution). No authority conflicts, no stale references, no cross-surface drift. PACKAGE-SYSTEM.md and BUILD-AND-DIST.md were not touched (diff is purely runtime + test-harness); grepped both to confirm no content drift occurred as a side-effect of the fix.
- **CDD γ (cycle economics):** **4/4** — Review rounds at best-case target (1/2 budget used); superseded PRs 0; closure loop closed in the same cycle session (β reviewed, merged, produced release commit, created tag, wrote this assessment, writing β close-out next). No §9.1 trigger fired (review rounds ≤ 2, mechanical ratio N/A with 0 findings, no loaded skill failed to prevent a finding it covers, no avoidable tooling failure — the 403 tag push is environmental). α close-out and this assessment both committed directly to main per §1.4 step 11 convention (surviving squash-merge).
- **Weakest axis:** none (all three at 4/4).
- **Action:** **none** — all three axes at 4/4. No skill patch required, no doc patch required, no new mechanical check required. This is the shape cycle economics should take when issue quality, implementation discipline, and independent review all align. Continue.

### 5. Production Verification

**Scenario:** A hub pins a strict subset of available packages in `.cn/deps.json`; `cn deps lock` produces a lockfile of exactly that subset; `cn deps restore` installs exactly that subset; no extraneous packages leak in from the index.

**Before this release:** Given an index with N packages × M versions and a `.cn/deps.json` pinning 1 package, `cn deps lock` wrote a lockfile with every (name, version) pair in the index (N × M entries). `cn deps restore` then installed all of them, with last-writer-wins on `restoreOne` resolving duplicate-name entries. The pin in `deps.json` had no effect on the installed set. This was reproducible in the Tier-2 CI hub (pinned 3 packages, installed 5) and in the issue's 5-line reproducer (pinned `cnos.core`, got 5 packages installed).

**After this release:** `cn deps lock` reads `.cn/deps.json` first, resolves each pin against the index via `idx.Lookup(name, version)`, collects unresolved pins into a single explicit error, and writes a lockfile of exactly the pinned set. `cn deps restore` then installs exactly that set — no extras, no duplicates. A pin to a non-existent version errors rather than silently succeeding; a missing `deps.json` errors rather than silently succeeding with an index dump.

**How to verify:**

1. **Unit**: `go test ./internal/restore -run TestGenerateLockFromIndex -v` → six test names directly encode the contract (`_FiltersByManifest`, `_NameAppearsAtMostOnce`, `_RestoreInstallsOnlyPinned`, `_MissingManifest`, `_PinNotInIndex`, `_DefaultProfile`). 6/6 PASS on `4f860f1` (release commit).
2. **Integration**: `TestGenerateLockFromIndex_RestoreInstallsOnlyPinned` stages two real tarballs in `dist/`, pins one in `.cn/deps.json`, runs lock → restore, and asserts only the pinned package lands in `.cn/vendor/packages/`. The unpinned package's vendor directory must not exist post-restore. Passes on fix, fails on revert.
3. **Tier-1 kata**: `scripts/kata/06-install.sh` overwrites the `cn setup` default with an explicit `cnos.core` pin read from source, runs `cn deps lock && cn deps restore`, and asserts `INSTALLED=1 && -f .cn/vendor/packages/cnos.core/cn.package.json`. Previously accepted `≥1`; now rejects anything other than `exactly cnos.core`. CI-gated via the `kata-tier1` workflow job.
4. **Tier-2 CI**: `.github/workflows/ci.yml` Tier-2 setup derives real versions at job time from `src/packages/*/cn.package.json`, writes an array-schema `deps.json` pinning `cnos.core`, `cnos.kata`, `cnos.cdd.kata`, runs `cn deps lock`, and later workflow steps verify the lockfile and install shape. The per-PR `kata-tier2` check was green on `44d431a`.
5. **Revert audit**: α documented running the tests with `GenerateLockFromIndex` reverted to the bug behavior; all six tests fail with the exact symptoms from the issue. β accepted α's revert evidence without re-running (within-session tooling cost would be ≥5min vs α's prior evidence being deterministic); if γ requires independent revert-verification in a future cycle, the procedure is (a) revert `src/go/internal/restore/restore.go` line ~416-448 to the pre-PR block, (b) `go test ./internal/restore -run TestGenerateLockFromIndex`, (c) observe 6 failures.

**Result:** **Pass** on all five verifications at code-level + test-level + CI-level. All 7 CI checks green on `44d431a` (PR head) and `4f860f1` (merged squash); `go test ./... && go test -race ./internal/restore/` green locally on the worktree of `44d431a` and on `origin/main` at the release commit `0dedbc9`.

**Binary-level verification:** Release commit `0dedbc9` landed on `origin/main` (`git push origin main` succeeded cleanly). Tag `3.56.2` was created locally but **push of the tag hit a transient `HTTP 403 curl 22` failure** — the same environmental issue documented in the 3.56.1 assessment §3 friction log. β applied the 3.56.1 guidance (do not trust `Everything up-to-date` after a `403`; verify remote state independently) — `git ls-remote --tags origin 3.56.2` confirms the tag is not yet on the remote as of this assessment's first commit. A background retry loop (`b39vzm337`) is running (`until git push origin 3.56.2 lands; do sleep 10; done`) and will succeed once the transient 403 clears. Once the tag lands, the release workflow will fire and build the 4 binaries + package tarballs + `index.json` + `checksums.txt`. **This section will be updated with the release URL and artifact list once the tag-push completes** (see 3.56.1 precedent for the shape; the 3.56.1 tag-push also hit 403s and eventually landed, with the assessment then updated in a correction commit).

**End-to-end runtime validation (deploy + `cn --version` on target host + `cn deps restore` from new binary):** deferred — not performed in this β session (requires a target host and a running daemon). The binary artifacts will be present at the release URL once the workflow completes; they are SHA-256-verifiable via `checksums.txt` at that point. Candidate for next β session or operator validation.

### 6. CDD Closeout

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0–1 | issue #250 | cdd | γ selected #250 as next MCA per 3.56.1 assessment commitment (§3.3 assessment-commitment default applies) |
| 2–7 | PR #259 | cdd, eng/go, eng/testing | α implemented `ParseManifest` + `ReadManifest` + rewrote `GenerateLockFromIndex`; 6 AC-named tests; CI harness alignment (Tier-1 + Tier-2); cnos.kata `lib.sh` array-schema + `pkg_version_from_source` helper |
| 7a | Pre-review | cdd | α held PR in draft until CI green on `44d431a`; rebased on `origin/main`; self-coherence table in PR body |
| 8 | Review R1 | cdd, review, eng/go, eng/testing | β: APPROVED (as comment — review §7.1 shared-identity), zero findings, module-truth + sibling-fallback + Go/OCaml parity audits clean |
| 9 | Gate + merge | cdd, release | β: PR transitioned draft → ready, squash-merged as `4f860f1` on `origin/main` |
| 10 | Release | cdd, release, writing | β: release commit `0dedbc9` on `origin/main` (VERSION 3.56.1 → 3.56.2, manifests stamped via `scripts/stamp-versions.sh`, CHANGELOG ledger row + detailed section, RELEASE.md rewritten). Tag `3.56.2` created locally and pushed — tag-push hit transient 403 (known pattern from 3.56.1); background retry loop `b39vzm337` will land the tag once the environmental issue clears. Release workflow will fire on tag landing. |
| 11 | Observe | cdd, post-release | Runtime/design alignment verified at code + test + CI level (§5). Binary-level observation deferred until tag lands on remote and release workflow completes. |
| 12 | Assess | cdd, post-release | this file |
| 12a | Skill patch | cdd, post-release, eng/skill | **Skill patch shipped** — `eng/go` §2.12 (Schema and compatibility) gains "sibling harnesses are contract surfaces too" addendum with mechanical grep form. Proposed by α in close-out (observed that `write_deps_json` schema mismatch survived authoring-time audit because §2.12 covered Go code only, not shell harnesses producing JSON for the same parsers). β concurs and shipped in the same commit as this assessment correction. Derivation cite to this cycle preserved in skill body. The row 7a peer-enumeration addendum proposed in 3.56.1 α close-out remains pending γ triage on its own merits (orthogonal to this cycle's shell-harness scope gap). |
| 13 | Close | cdd, post-release | β close-out to be committed to `.cdd/releases/3.56.2/beta/250.md` alongside this assessment (main-direct, survives squash-merge per §1.4 step 11). α close-out expected from α agent at `.cdd/releases/3.56.2/alpha/250.md` per α algorithm step 11. |

### 6a. Invariants Check

| Constraint | Touched? | Status (preserved / tightened / revised / N/A) |
|---|---|---|
| T-002 (kernel minimal, dispatch-only in cli/) | no | N/A — `cmd_deps.go` unchanged; domain logic stays in `internal/restore/` |
| T-003 (Go sole language in kernel) | no | N/A — all changes in Go or shell harness |
| T-004 (source/artifact/installed explicit) | yes | **tightened** — the lockfile (artifact) now mechanically derives from the manifest (source of operator intent) rather than from ambient index state. The `deps.json` (source) → `deps.lock.json` (artifact) → `.cn/vendor/packages/` (installed) chain is now authority-preserving end-to-end: manifest pin flows through lock to install without ambient-state interference. Pre-cycle the chain was manifest → lock (ignored) → install (polluted by index); post-cycle it is manifest → lock → install, each step derivable from the prior |
| T-005 (content classes finite) | no | N/A |
| INV-001 (one package substrate) | no | N/A |
| INV-003 (commands/skills/orchestrators distinct) | no | N/A |
| INV-004 (kernel owns precedence) | no | N/A — precedence is within the pin-set (manifest authority); index is lookup-only, never authoritative for which packages to install |

Invariants document (`docs/alpha/DESIGN-CONSTRAINTS.md`) was not modified by this cycle; the tightening of T-004 is a preserve-tightened status update, not a revision.

### 7. Next Move

**Next MCA:** #230 — `cn deps restore: version upgrade skips silently with version-less VendorPath`

**Owner:** α (via γ dispatch)

**Branch:** pending (canonical form: `alpha/230-version-upgrade-skip` per CDD §4.2, `<agent>/<issue>-<scope>`)

**First AC:** `cn deps restore` detects when the lockfile pins a version different from the installed `cn.package.json` version and re-installs (rather than returning early on directory existence). The existing `restoreOne` skip (`restore.go:136`) must compare installed manifest version against lockfile version, not merely check for directory presence.

**MCI frozen until shipped?** **Yes** — 9 growing-lag issues; freeze continues (still well over the 3-issue threshold). The freeze commitment from 3.56.0 → 3.56.1 → 3.56.2 remains in force. The next two cycles should continue picking growing-lag bugs with clean fix paths rather than opening new design fronts.

**Rationale:**

- **Adjacency.** #230 lives in the same file (`src/go/internal/restore/restore.go`) as #250, in the same subsystem (deps), and the `restoreOne` function already carries a `#230`-tagged comment (lines 130-135) naming exactly this gap. The skip-if-directory-exists shortcut was intentional when versions were never compared; now that `cn deps lock` produces a faithful lockfile (this cycle), the next coherent step is making `cn deps restore` honor version changes in that lockfile.
- **Completes the authority chain.** 3.56.2 closed the "what packages should this hub use" → "what the lockfile says" gap. #230 closes the "what the lockfile says" → "what is actually installed" gap. Together these two cycles restore the end-to-end manifest→lock→installed authority chain that `.cn/deps.json` was supposed to be the source of.
- **P1 bug with a narrow fix path.** `restoreOne` already has the package manifest in hand after validation (the last step validates `cn.package.json` name); reading the installed `version` field is one extra `os.ReadFile` + `json.Unmarshal` (or it can be piggy-backed on the existing manifest read that already happens). The fix is single-function, single-file, no API change.
- **The 3.56.1 assessment already queued this.** 3.56.1's §7 next-move commitment listed #230 as a candidate. 3.56.2 validates that commitment by picking #230 again as the next candidate — two consecutive assessments agree.

**Alternative candidate:** #235 (`cn build --check should validate command entrypoints and skill paths against manifest`) is also P1 and growing-lag. It is further from #250/#230's subsystem (lives in `pkgbuild/` rather than `restore/`) and requires a slightly larger fix (enumerate manifest-declared commands, stat each entrypoint, check exec bit). Deferred unless γ triage revises priority.

**Closure evidence (CDD §10):**

- **Immediate outputs executed:**
  - [x] PR #259 merged as `4f860f1` (squash) on `origin/main`
  - [x] Release commit `0dedbc9` on `origin/main` (VERSION 3.56.2, all manifests stamped, CHANGELOG ledger row + detailed section, RELEASE.md rewritten)
  - [ ] **Tag `3.56.2` push — pending remote acceptance.** Local tag created; push hitting transient 403 (known environmental pattern, see §3 and §5). Background retry loop was started but has been cleaned up; if the 403 persists, operator can push the tag manually once the environmental issue clears.
  - [ ] **GitHub release `3.56.2` publication — pending tag landing** (workflow fires on tag push)
  - [x] **α close-out** committed at `.cdd/releases/3.56.2/alpha/250.md` by α (`855410e`)
  - [x] **β close-out** committed at `.cdd/releases/3.56.2/beta/250.md` by β (`f302428`)
  - [x] This post-release assessment committed to `docs/gamma/cdd/3.56.2/POST-RELEASE-ASSESSMENT.md` (`f302428`, amended by this correction commit)
  - [x] **Skill patch: `eng/go` §2.12 sibling-harness scope addendum** shipped in this correction commit, per α close-out proposal and §9.1/§10.1 MCA discipline

- **Deferred outputs committed:**
  - **Issue:** #230 — `cn deps restore: version upgrade skips silently with version-less VendorPath` (already filed with P1 label, converged design exists in issue body + `restoreOne` code comment)
  - **Owner:** α (next cycle)
  - **Branch:** pending
  - **First AC:** restore detects version mismatch and re-installs (rather than directory-exists early return)
  - **MCI freeze state:** unchanged (freeze continues until growing-lag < 3)

**Cycle close condition met when:**

1. Tag `3.56.2` lands on remote (background loop `b39vzm337`)
2. Release workflow completes and publishes artifacts (automatic on tag push)
3. α close-out written to `.cdd/releases/3.56.2/alpha/250.md` (α's responsibility, independent of β)
4. β close-out + this assessment pushed to `origin/main` (β's next commit)

Items 1 and 2 are environmental/automatic. Items 3 and 4 are the mechanical closeout. Post both, the cycle is closed per §10.3.

### 8. Hub Memory

- **Daily reflection:** n/a in this session (β does not have write access to the hub-memory repo from this environment; the `cn-beta/threads/` surface is outside this worktree). Candidate for an operator-side write once hub memory access is available. Note: this is a recurring environmental gap across cycles, not a β discipline gap.
- **Adhoc thread update:** n/a (same environmental constraint as above). The relevant ongoing thread would be "lockfile-authority-chain" or similar — #250 closes the manifest→lock link; #230 is queued for the lock→install link. An operator-side thread update can record the arc once hub memory access is available.

**Note:** Post-release §7 names hub-memory writes as mandatory before cycle close. This session does not have hub-memory access, so those writes are **deferred to the operator** with an explicit note rather than silently skipped. The assessment itself and the β close-out (in the main repo) are the durable within-repo surface that replaces this cycle's hub-memory write. Next β session with hub access should retroactively record the cycle if the operator has not already done so.







