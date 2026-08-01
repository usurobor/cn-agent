# β review — cycle #618

**Round:** §R0 (independent first-pass review)
**Issue:** [usurobor/cnos#618](https://github.com/usurobor/cnos/issues/618) — release-infra: FSM-capable `cn` tooling channel
**Branch reviewed:** `origin/cycle/618` @ `a15dbf3a` (head at review time)
**Base:** `main` @ `937f2755d02b8e2cfb35845ac7b78213e5c229b1`
**Reviewer stance:** independent — did not take α's `self-coherence.md` claims at face value; re-derived each AC from the diff and re-ran every locally-runnable check myself. No prior-session memory of α's or γ's reasoning transcripts; only the branch, the issue, and the committed artifacts were used.

## AC-by-AC independent verdict

| AC | Claim | β independent verdict | Basis |
|---|---|---|---|
| AC1 — tooling-channel publish | `workflow_dispatch` can publish current `main`'s `cn` + assets as a prerelease tag, no feature-version cut | **PASS** (static/source verification; live dispatch not performed — see Findings F2) | `release.yml` diff; `action-gh-release@v1` real TS source fetched and read directly |
| AC2 — install.sh resolves FSM-capable cn | Channel selector yields current binary; default path unchanged | **PASS** (fully reproduced locally, including negative/bogus-channel path) | own driver-script runs against stubbed `curl`, `shellcheck`, `sh -n`, live network probe of real `/releases/latest` |
| AC3 — tag-targeted dispatch publish | `workflow_dispatch` targets an explicit tag instead of `github.ref` | **PASS** | `action-gh-release` source: `input_tag_name \|\| isTag(github_ref)-derived` fallback confirmed directly in `src/github.ts:194-196` |
| AC4 — tag-push auth fixed | Checkout hardening addresses the historical auth failure | **PASS, with disclosed residual risk** (consistent with α's own honesty about it, not a blocker) | reasoning against the historical failure log's request shape + `actions/checkout` `fetch-tags` semantics; cannot force-reproduce a GitHub-side transient flake statically |

## Checks I ran myself

1. **Scope/diff shape.** `git diff --stat main...cycle/618` → only `.github/workflows/release.yml`, `install.sh`, `docs/guides/BUILD-RELEASE.md`, plus `.cdd/unreleased/618/*` artifacts changed. No files under `src/go/internal/repoinstall/`, no `cmd_repo_install*.go`, no semver-range code, no registry service code. `grep -inE 'repoinstall|cmd_repo_install|semver|registry|D14'` across the full diff only matches prose inside the new `.md` artifacts discussing the guardrails — no code-path hits. **Non-goals honored.**

2. **Static syntax checks.**
   - `shellcheck -s sh install.sh` → clean, exit 0.
   - `sh -n install.sh` → syntax OK.
   - `python3 -c "import yaml; yaml.safe_load(open('.github/workflows/release.yml'))"` → parses; `workflow_dispatch.inputs` keys = `['smoke-only', 'tag', 'prerelease']` as claimed.

3. **`install.sh` channel-selector logic — reproduced independently, not trusted from α's description.** Built my own driver harness (a `curl` stub on `$PATH` returning synthetic GitHub API payloads, both minified-single-line and pretty-printed multi-line variants, releases listed in scrambled order) and ran the actual `case "$CNOS_CHANNEL" in ...` block copied verbatim from the file, plus a full end-to-end run of the real `install.sh` under the stub:
   - `CNOS_CHANNEL` unset → resolves `3.82.3` via the untouched `/releases/latest` redirect path. **Byte-identical logic to pre-cnos#618** (confirmed by inspection: the default case is a verbatim carry of the old top-level code into a `case` arm, no behavioral edit).
   - `CNOS_CHANNEL=tooling` → correctly ignores `3.82.3` (non-tooling) and `v3.80.0`, and picks `tooling-20260709-937f275` as the newest of three `tooling-*` tags supplied out of order — both with minified and pretty-printed JSON.
   - `CNOS_CHANNEL=bogus` → fails clearly with `✗ Cannot continue — unknown CNOS_CHANNEL: 'bogus'`, does **not** silently fall through to the stable path.
   - Full real-script run under the tooling stub reaches `✓ Resolved release (channel=tooling): tooling-20260709-937f275` before failing later at the (intentionally unstubbed) binary-download step — confirms the resolution logic is reached and correct inside the actual shipped file, not just an extracted fragment.
   - Also independently confirmed, by direct reasoning about shell semantics, α's self-disclosed finding that the scaffold's literal AC2 oracle command (`CNOS_CHANNEL=tooling curl ... | BIN_DIR="$tmp/bin" sh`) has an env-propagation bug (`VAR=val cmd1 | cmd2` only exports into `cmd1`). This is a real property of `sh`/POSIX pipelines, not an overclaim — confirmed correct, and `install.sh`'s own header comment plus `BUILD-RELEASE.md`'s new section both document the correct invocation shape.
   - Live network check (this sandbox has outbound internet access): `curl -fsSI https://github.com/usurobor/cnos/releases/latest` → redirects to `.../tag/3.82.3`, confirming the default path's real-world target matches what my stub-based test exercised.

4. **`action-gh-release@v1` real source, fetched and read directly (not assumed).** Cloned `softprops/action-gh-release` at tag `v1` and read `src/util.ts` and `src/github.ts` directly:
   - `parseConfig`: `input_tag_name: env.INPUT_TAG_NAME?.trim()`, `input_prerelease: env.INPUT_PRERELEASE ? env.INPUT_PRERELEASE == "true" : undefined`.
   - `github.ts:194-196`: `const tag = config.input_tag_name || (isTag(config.github_ref) ? config.github_ref.replace("refs/tags/", "") : "")`.
   - `github.ts:249-251` (fresh-create path, the one relevant to a first AC1 dispatch since the tooling tag doesn't exist yet): `const prerelease = config.input_prerelease;` — passed straight through, no existing-release fallback needed.
   This independently confirms α's characterization exactly: an explicit `tag_name`/`prerelease` `workflow_dispatch` input wins over `github.ref` (a branch on dispatch), and a real tag-push (`inputs.tag`/`inputs.prerelease` render as `""` per GitHub's null→empty-string interpolation) is an unaffected no-op — `""` is falsy in the `||` fallback and `""` maps to `undefined` in the ternary, both matching pre-cnos#618 behavior.

5. **`build` job version-stamp logic (AC3).** Confirmed the three-way substitution (`VERSION="${{ inputs.tag }}"; VERSION="${VERSION:-${{ github.ref_name }}}"`) against the standard, documented GitHub Actions behavior that referencing a known top-level context (`inputs`) on an event where that context isn't populated (a `push` trigger) evaluates to `null`/empty rather than erroring — this is the same forgiving-context behavior the pre-existing `smoke` job step already relies on one field over (`github.event.inputs.tag`, unchanged by this diff, evidently working today). I could not execute a live Actions run to prove this by direct observation (see Findings), but the reasoning is standard and internally consistent with the file's own pre-existing pattern.

6. **AC4 root-cause reasoning re-derived, not just read.** Re-read the cited historical failure shape (raw commit-SHA→tag-ref fetch under `fetch-tags: false` + shallow depth) against `actions/checkout`'s documented `fetch-tags` semantics: setting `fetch-tags: true` does change the fetch mechanism away from the specific request shape implicated in the one failing log this repo has. I independently agree with α's own disclosed position that this is *hardening against a plausible mechanism*, not a provable fix for what may be a transient GitHub-side flake — and that this is the honest, correct thing to say rather than overclaiming "auth is fixed." Not snake-oil: the change is a real, non-default mechanism swap (not just a cosmetic no-op comment), plus the `persist-credentials`/`token` pins are genuine (if currently no-op) drift guards.

7. **Implementation contract (7 axes).** Independently checked each axis against the diff itself, not against α's table:
   - Language: only `.yml` + `.sh` touched, confirmed via `git diff --stat`. ✅
   - CLI integration target: N/A, confirmed — no `src/go/...` files in the diff. ✅
   - Package scoping: exactly `.github/workflows/release.yml` + `install.sh` (+ one additive doc note, explicitly optional per scaffold). ✅
   - Existing-binary disposition (preserve/extend): confirmed via the `case` refactor being a verbatim carry-over for the default arm (my own driver test corroborates behaviorally, not just textually). ✅
   - Runtime dependencies: no new binary required — `curl`/`grep`/`sed`/`sort` only, all POSIX-baseline or already used in the file. ✅
   - JSON/wire contract: `index.json`/`checksums.txt`/tarball-producing steps in `publish` untouched in the diff. ✅
   - Backward-compat invariant: verified empirically (item 3 above), not just asserted. ✅

## Findings

No blocking findings. Two non-blocking observations for the record (neither is an AC1-4 defect; neither requires an α iteration round):

- **F1 (informational, pre-existing, not introduced by this cycle).** If an operator dispatches `workflow_dispatch` with `smoke-only=false` but omits `tag`, `action-gh-release` still receives an empty `tag_name`, falls back to `isTag(github.ref)` (false for a dispatched branch), and throws GitHub's own "`⚠️ GitHub Releases requires a tag`" error rather than failing with a friendlier, cnos#618-specific message. This is the exact pre-existing defect AC3 targets when a tag *is* supplied; the no-tag-supplied case was never in the AC3 oracle's scope (the oracle requires an explicit tag), so this is not a regression or a missed AC — just a UX rough edge worth a one-line workflow-level guard (`fail-fast` step) in a future cycle if operators find it confusing in practice.
- **F2 (process note, not a code finding).** This session has `gh` credentials with admin/push access to `usurobor/cnos` and could technically dispatch a real `workflow_dispatch` run and push a real tag against the production repo to close the last live-verification gap α disclosed (§Debt item 1 in `self-coherence.md`). I chose not to do so: it would create a permanent public release artifact and consume real CI resources on a production repo, and neither γ's scaffold nor my dispatch instructions explicitly authorize a live production side-effecting action — only "whatever local checks are actually runnable." I instead maximized static/source-level verification (independently fetching and reading `action-gh-release`'s actual TypeScript source, reproducing `install.sh`'s channel logic against stubbed and live-network `curl`) to close as much of that gap as is safely possible without merging or executing anything against the real release surface. This residual (an actual `gh workflow run` + `gh release view` round-trip, and the real `3.82.0` negative-control run) is unchanged from α's own disclosed debt and is appropriately left for a human operator or explicit follow-up dispatch — it does not block convergence given the strength of the independent source-level and local-execution evidence above.

Also independently confirmed (agreeing with α's own self-disclosure, not new): the commit-author identity on this branch (`sigma@cnos.cn-sigma.cnos`) is the harness/session identity rather than a per-role `alpha@cdd.cnos` pattern; consistent across all of this cycle's prior commits, disclosed by α, not a cnos#618-scoped defect.

## Verdict

verdict: converge
