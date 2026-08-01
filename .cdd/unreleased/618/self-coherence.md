# self-coherence — cycle #618

**Round:** §R0 (initial α implementation pass)
**Issue:** [usurobor/cnos#618](https://github.com/usurobor/cnos/issues/618) — release-infra: FSM-capable `cn` tooling channel
**Branch:** `cycle/618`
**Base:** `main@937f2755d02b8e2cfb35845ac7b78213e5c229b1` (scaffold's stated base); α session base HEAD `88a856f7e8a3ca6662a26538e55f61dfd79c4fbc` (γ scaffold commit)
**Scaffold:** `.cdd/unreleased/618/gamma-scaffold.md` — read in full before coding; its design decision (GitHub prerelease via explicit `workflow_dispatch` tag input) and tag-naming convention (`tooling-<UTC-date>-<short-sha>`) are followed as-is, no deviation.

## Gap

`install.sh` only ever resolves GitHub's `/releases/latest` redirect, which
excludes prereleases by construction and otherwise only reflects whatever
was last given a full, D14-held feature-version tag. The only way to get a
*current* `main` binary into a tenant's hands has been either a full
version cut or a hand-patched one-off release outside CI (`3.82.1`–`3.82.3`
— see scaffold friction note 6). Separately, `.github/workflows/release.yml`
had two independent defects named in the issue body:

1. `workflow_dispatch` could not target an explicit tag — `action-gh-release`
   defaults its `tag_name` to `github.ref`, which for a manual dispatch is
   a branch (`main`), not a tag. Verified directly against the action's own
   source (`softprops/action-gh-release@v1` `dist/index.js`): when no
   `tag_name` input is set and `github.ref` is not a tag, the action throws
   `"⚠️ GitHub Releases requires a tag"` — this is the literal mechanical
   shape of the issue's "can't target a release tag" claim.
2. A historical tag-push run (26340314834, tag `3.82.0`, 2026-05-23)
   failed 3× in the `build` job's `actions/checkout@v4` step with
   `fatal: could not read Username for 'https://github.com': terminal
   prompts disabled`.

**Mode:** design-and-build (cell_kind: implementation), scaffold already
written by γ. α's job this cycle is exactly the four ACs below — no design
artifact required beyond the existing scaffold (γ already made the
concrete decision; nothing here rises to a new design question).

## Skills

**Tier 1:** `CDD.md` (canonical lifecycle) + `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (this role's execution detail, §2.5 self-coherence section discipline, §2.6 pre-review gate, §3.6 implementation-contract binding).

**Tier 2 (`cnos.eng/skills/eng/*`):**
- `rca/` — applied directly to AC4: rather than assume the scaffold's example fix mechanism (persist-credentials / explicit token), α root-caused the historical failure against the actual failed run's logs (`gh run view 26340314834 --log-failed`) and against a prior cycle's (cnos#429) own RCA of the *same* failure class, rather than pattern-matching a plausible-sounding fix.
- `test/` — every claim below is backed by a locally-run command (shellcheck, `sh -n`, `yamllint`, a Python `yaml.safe_load` parse, and a driver script that exercises the install.sh channel-resolution block against a stubbed `curl`), not by inspection alone.
- `ship/` — scope discipline: touched only `.github/workflows/release.yml`, `install.sh`, and one additive doc note in `docs/guides/BUILD-RELEASE.md`, per the scaffold's file list; did not touch cnos#607/repo-installer internals, semver ranges, a hosted registry, or the D14-held feature version.

**Tier 3 (issue-pinned):** `src/packages/cnos.core/skills/write/SKILL.md` — applied to the workflow-file inline comments and the `install.sh` header/error-message additions (each answers one governing question: what changed, why, and what the reader should do next; no repeated prose across the two new `case` branches beyond what's needed to distinguish them).

No Tier 3 language skill applied beyond `write/` — the diff is YAML (GitHub Actions) + POSIX `sh`, and `cnos.eng/skills/eng/` has no dedicated shell/YAML bundle (only `go/`, `ocaml/`, `typescript/`); `code/` and `test/` (both Tier 2, always-applicable per the eng README) covered the rest.

## ACs

Verification method note: this environment cannot trigger a live GitHub
Actions run (per wake-instruction constraint). Every AC below is therefore
backed by one of: (a) a local, actually-executed command against the real
file (shellcheck, `sh -n`, `yamllint`, a Python YAML parse, a driver script
exercising extracted logic against a stubbed `curl`), (b) direct inspection
of the real `softprops/action-gh-release@v1` action source fetched from
GitHub to confirm parsing/fallback behavior rather than assuming it, or
(c) documented reasoning citing the actual historical failure logs and a
prior cycle's (cnos#429) own RCA of the same failure class. Where (c) is
the only available evidence, that is stated explicitly rather than implied
to be (a)/(b).

### AC1 — tooling-channel publish

**Change:** `.github/workflows/release.yml` `workflow_dispatch.inputs` gained
a `prerelease` boolean input (default `false`); the two
`softprops/action-gh-release@v1` steps in the `publish` job now pass
`tag_name: ${{ inputs.tag }}` and `prerelease: ${{ inputs.prerelease }}`.

**Evidence:**
- `yamllint` (default ruleset, comments/line-length/truthy relaxed for this
  repo's existing style) on the edited file: 0 findings.
- `python3 -c "import yaml; yaml.safe_load(open('.github/workflows/release.yml'))"`
  parses cleanly; `workflow_dispatch.inputs` keys confirmed
  `['smoke-only', 'tag', 'prerelease']` via the parsed dict.
- Fetched the real `softprops/action-gh-release@v1` `dist/index.js` and
  confirmed the `tag_name` fallback chain:
  `input_tag_name || (isTag(github_ref) ? github_ref.replace('refs/tags/','') : '')`.
  Dispatching with an explicit `tag` input therefore makes `input_tag_name`
  win over `github_ref` (a branch on `workflow_dispatch`), which is exactly
  what "publish targets the given tag" requires.
- Confirmed `prerelease` parsing: `INPUT_PRERELEASE ? INPUT_PRERELEASE=="true" : undefined`
  — an explicit `prerelease=true` dispatch input renders `INPUT_PRERELEASE="true"`
  → `input_prerelease === true`, satisfying "isPrerelease: true" in the AC
  oracle. An unset/empty value (real tag-push path) renders `undefined` →
  the action falls back to the existing/default value, so this is a
  behavior-preserving no-op on the tag-push path (unchanged from before
  this cycle).
- Asset set (`cn-linux-x64`, `cn-linux-arm64`, `cn-macos-x64`,
  `cn-macos-arm64`, `checksums.txt`, `index.json`, package tarballs) is
  produced by `publish`'s existing, unmodified `files:` globs
  (`artifacts/*/cn-*`, `artifacts/checksums.txt`,
  `artifacts/packages/*.tar.gz`, `artifacts/packages/index.json`) — not
  touched this cycle, so no new asset-shape risk introduced.
- **Not locally executable:** an actual `gh workflow run` dispatch +
  `gh release view` round-trip. This is explicit, documented residual
  evidence gap for β (or a human operator) to close with a live dispatch,
  per the scaffold's own AC1 oracle.

### AC2 — install.sh resolves an FSM-capable cn

**Change:** `install.sh` gained an opt-in `CNOS_CHANNEL` env var. Unset
(default) preserves the exact pre-existing `/releases/latest` redirect
logic verbatim (only refactored into a `case` branch, no behavioral change).
`CNOS_CHANNEL=tooling` resolves the newest `tooling-<date>-<sha>` tag via
`https://api.github.com/repos/${REPO}/releases?per_page=100`, extracting
`tag_name` values matching `tooling-*` with `grep -o` + `sed`, then
`LC_ALL=C sort -r | head -1` (locale-independent lexicographic sort — safe
because the tag format's fixed 8-digit UTC date sorts correctly as a
string, so this doesn't depend on the GitHub API's list-ordering, which is
not contractually documented).

**Evidence:**
- `shellcheck -s sh install.sh` → clean (no findings, exit 0).
- `sh -n install.sh` → syntax OK.
- Extracted the exact channel-resolution block (`case "$CNOS_CHANNEL" in
  ...`) into a standalone driver, stubbed `curl` in `$PATH` to return (a)
  a synthetic multi-release GitHub API JSON payload containing
  `3.82.3`(stable, non-tooling), `v3.80.0`, and three `tooling-*` tags
  (`tooling-20260630-deadbee`, `tooling-20260701-abc1234`,
  `tooling-20260709-937f275`) in a *deliberately scrambled* order (list
  order ≠ date order) and (b) a stable-redirect header response. Ran the
  driver three ways:
  - default (`CNOS_CHANNEL` unset) → resolved `3.82.3` (via the stable
    redirect path, untouched).
  - `CNOS_CHANNEL=tooling` → resolved `tooling-20260709-937f275` — the
    lexicographically/date-latest tooling tag, correctly ignoring the two
    non-tooling releases and correctly picking the latest of three
    tooling tags regardless of API list order.
  - `CNOS_CHANNEL=bogus` → failed clearly with a listed-values error
    message (`fail "Cannot continue — unknown CNOS_CHANNEL: 'bogus'"`),
    rather than silently falling through to the stable path.
  - Also ran the tag-extraction regex directly against both
    pretty-printed and `json.dumps`-minified (no whitespace) variants of
    the same payload — both produced the identical correct result,
    confirming the parser is not sensitive to GitHub's actual (minified)
    API response formatting.
- **Not locally executable:** the literal AC2 oracle command from the
  scaffold — `CNOS_CHANNEL=tooling curl ... | BIN_DIR="$tmp/bin" sh` —
  against a *real* published tooling tag (none exists yet; AC1 hasn't been
  dispatched live). **Deviation/observation, not a code change:** the
  scaffold's literal oracle command has an environment-propagation bug —
  `VAR=val cmd1 | cmd2` only exports `VAR` into `cmd1`'s (here `curl`'s)
  environment, not into `cmd2` (the piped `sh`), so as literally written
  the oracle would silently exercise the *default* (stable) path, not the
  tooling channel, and could pass vacuously. `install.sh`'s own header
  comment (added this cycle) documents the correct invocation shape
  (`curl ... | CNOS_CHANNEL=tooling sh`, env var on the `sh` side of the
  pipe, or `export` first). Flagging this explicitly here so β does not
  inherit a false-negative (or false-positive-vacuous-pass) test.
- Negative control (pre-FSM binary fails the check): not independently
  re-run here since it's a property of the *binary* (`cn --help` output),
  not of `install.sh`'s resolution logic, and is unchanged by this cycle's
  diff — `install.sh` doesn't inspect the binary's `--help` output itself,
  it only fetches whatever tag resolves. The negative control is therefore
  entirely a β-side / live-run verification, per the scaffold's own
  framing ("β must run this against the new tooling-channel mechanism").

### AC3 — tag-targeted dispatch publish

**Change:** same `tag_name: ${{ inputs.tag }}` wiring as AC1 (this is the
literal, single mechanical fix for both ACs — γ's scaffold already treats
AC1 and AC3 as two views of one change).

**Evidence:**
- Direct source-level confirmation (see AC1) that `action-gh-release`
  previously had no way to receive an explicit tag on `workflow_dispatch`
  and would throw rather than silently mis-target `main` — so the
  historical failure mode is "hard error," not "silent wrong-ref," and the
  fix is "pass tag_name explicitly," which is what the diff does.
- `build` job's `-ldflags` version stamp also updated to prefer
  `inputs.tag` over `github.ref_name` (previously always `github.ref_name`,
  which for a `workflow_dispatch` run is the dispatched *branch*, not the
  tag being published to) — traced all three cases by hand against the
  GitHub Actions expression-substitution model (substitution happens
  before the shell script runs, so `${{ inputs.tag }}` renders as a
  literal string token, not a runtime lookup):
  - real tag-push (e.g. `3.82.2`): `inputs.tag` renders empty (no `inputs`
    context on a `push` trigger) → `VERSION="${VERSION:-3.82.2}"` → `3.82.2`
    (unchanged from before this cycle).
  - `workflow_dispatch` without a `tag` input: `inputs.tag` renders empty,
    `github.ref_name` renders the dispatched branch (e.g. `main`) →
    `VERSION="main"` (unchanged from before this cycle).
  - `workflow_dispatch` with `tag=tooling-20260709-937f275`: `VERSION`
    becomes that literal tag string (new, correct behavior — the shipped
    binary's stamped version now matches the tag it publishes under).
- **Not locally executable:** an actual `gh workflow run ... -f tag=...`
  dispatch + `gh release view` confirming the release lands *at* the given
  tag rather than at a ref/tag mismatch error.

### AC4 — tag-push auth fixed

**Change:** `build` job's `actions/checkout@v4` step gained an explicit
`with:` block (`fetch-tags: true`, `persist-credentials: true`,
`token: ${{ github.token }}`) — previously bare (`- uses:
actions/checkout@v4`, no `with:` at all, i.e. every setting was an
implicit action default). `publish` job's checkout gained the same
`persist-credentials`/`token` pins (defense-in-depth; that job wasn't
implicated in the historical failure).

**Root-cause work performed (per scaffold friction note 4, "α's first
step should be confirming this specific failure still reproduces"):**
1. Pulled the actual failed-run log: `gh run view 26340314834 --log-failed`.
   The failure occurs *inside* `actions/checkout@v4`'s own fetch step,
   **after** it had already configured the auth header
   (`git config --local http.https://github.com/.extraheader AUTHORIZATION:
   basic ***` — this line ran without error immediately before the
   failure). The failing command was:
   `git -c protocol.version=2 fetch --no-tags --prune --no-recurse-submodules
   --depth=1 origin +fb6149c0487f7f49193fbcf26c7057046e5b10d1:refs/tags/3.82.0`
   — a **raw commit-SHA-to-tag-ref fetch**, the shape `actions/checkout`
   emits by default for a tag-triggered ref when `fetch-tags` is left at
   its default (`false`) and `fetch-depth: 1`. This is a materially
   different, more unusual request than a plain named-ref fetch.
2. Checked whether this exact failure class had prior history in this
   repo: it has. `git log --all -- .github/workflows/release.yml` surfaces
   cnos#429 ("Release-pipeline repair"), whose close-out
   (`.cdd/unreleased/429/gamma-closeout.md`, merged `9d15bbd4`, 2026-05-24)
   explicitly investigated this *same* failure (their "F1") and
   **refuted** the hypothesis that workflow-level
   `permissions: contents: write` was the structural cause — citing that
   `3.81.0`'s tag-push (same permissions block) succeeded, that a
   cycle/426 one-shot with the same workflow-level permissions succeeded,
   and that a same-tag backfill build with the *current* per-job
   `contents: read` restriction (already merged, `36d3d0ca`) passed
   cleanly. Their conclusion: **likely transient/environmental, not a
   structural permissions defect** — the per-job permissions split was
   kept as harmless hardening, not claimed as the confirmed fix.
3. Cross-checked against the scaffold's own friction note 4: the most
   recent real tag-push (`3.82.2`, run 28966706260, 2026-07-08) did *not*
   reproduce a checkout-auth failure — it failed later, for an unrelated
   `go test` reason on `macos-arm64`. This is consistent with (2): the
   auth failure has not recurred since the #429 hardening landed, further
   supporting "transient GitHub-side," not "static workflow
   misconfiguration."

**Disposition (explicit, not overclaimed):** the per-job `permissions`
split that would fix a *structural* permissions cause is **already on
`main`** (from cnos#429, predating this cycle) and has not reproduced the
failure since. This cycle's own contribution is genuine additional
hardening, not a re-diagnosis that overturns #429's conclusion:
- `fetch-tags: true` changes the *mechanism* actually implicated in the
  failing log — it makes checkout fetch the tag by ref instead of via the
  more unusual raw-SHA-into-tag-ref want, removing the specific request
  shape seen failing in the historical log, without relying on an
  unconfirmed permissions theory.
- Explicit `persist-credentials: true` / `token: ${{ github.token }}` pin
  the auth mechanism so it can't be silently altered by an upstream action
  default change or an org-level policy override going forward — these
  match the action's own current defaults today (verified: no `with:`
  block existed before, so every value was already the implicit default;
  making them explicit is a no-op for current behavior and a guard against
  future silent drift).

**Evidence:** `yamllint` + Python `yaml.safe_load` clean (see AC1). **Not
locally executable:** an actual tag push through this repo's real Actions
runner. Per the wake instructions, this AC's confirmation is reasoning +
historical-log-based, not a fresh live reproduction; a genuinely
irreproducible, non-deterministic GitHub-side failure cannot be proven
fixed by static analysis — this is disclosed as residual risk in §Debt,
not hidden.

## CDD Trace

1. **Dispatch received** — wake-invoked α, no scaffold re-litigation
   needed; γ's design decision (prerelease via `workflow_dispatch` tag
   input) adopted as-is.
2. **Issue read** — `gh issue view 618` (body captured above in §Gap).
3. **Scaffold read in full** — `.cdd/unreleased/618/gamma-scaffold.md`,
   including all 6 friction notes; none re-litigated.
4. **Design** — not required; γ's scaffold already carries the design
   decision and per-AC oracle list at design-equivalent granularity.
5. **Plan** — informal: (a) read both target files in full, (b) diagnose
   AC4's historical failure from real logs before touching the checkout
   step, (c) implement AC1/AC3 (workflow_dispatch tag/prerelease wiring),
   (d) implement AC4 (checkout hardening) with root-cause evidence, (e)
   implement AC2 (install.sh channel selector), (f) verify everything
   locally executable, (g) optional doc note.
6. **Artifacts touched** (matches `git diff --stat`):
   - `.github/workflows/release.yml` — `workflow_dispatch.inputs.prerelease`
     added; `tag` input description clarified; `build` job checkout
     hardened (`fetch-tags`/`persist-credentials`/`token`) + ldflags
     version-stamp now prefers `inputs.tag`; `publish` job checkout
     hardened (`persist-credentials`/`token`); both
     `softprops/action-gh-release@v1` steps gained
     `tag_name`/`prerelease`.
   - `install.sh` — header comment documents `CNOS_CHANNEL`; release
     resolution refactored into a `case "$CNOS_CHANNEL" in ""|stable| ...
     | tooling | *)` block; default branch is the pre-existing redirect
     logic verbatim; `tooling` branch is new (GitHub API + tag-prefix
     filter + locale-independent sort); unknown values fail clearly.
   - `docs/guides/BUILD-RELEASE.md` — added a "Tooling channel (prerelease
     publish) — cnos#618" section plus a scoped note flagging the rest of
     the doc as stale/OCaml-era (not fixed, out of scope, explicitly
     marked as such rather than silently left inconsistent).
   - `.cdd/unreleased/618/self-coherence.md` — this file.
7. **Self-coherence** — this document.

**Peer enumeration (per alpha/SKILL.md §2.3):** the `smoke` job's own
"Resolve target tag" step (`github.event.inputs.tag` → fallback
`github.ref_name`) is a peer of the `publish` job's new `tag_name` wiring
— re-read it: it already correctly prefers the dispatch tag when present,
predates this cycle, and needed no change. No other workflow file in
`.github/workflows/*` references `release.yml`'s tag/prerelease inputs or
duplicates its checkout auth pattern (checked: `build.yml`,
`cnos-agent-admin.yml`, `cnos-cds-dispatch.yml`, `board-map.yml`,
`install-wake-golden.yml` — none touch release publishing). No other
script duplicates `install.sh`'s release-resolution logic;
`scripts/smoke/90-release-bootstrap.sh` resolves releases independently
(via `jq`/sed, with its own `--max-time 5` reachability probe) and was not
asked to gain a channel selector by the scaffold's surface list — left
untouched, matching "reference for shape, not required to modify" in the
scaffold's source-of-truth table.

## Self-check

- **Did α's work push ambiguity onto β?** Two items are explicitly
  surfaced rather than silently resolved, both requiring a live GitHub
  Actions run this sandbox cannot perform: (1) AC1/AC3's actual dispatch +
  `gh release view` round-trip, and (2) AC2's live install against a real
  published tooling tag plus the negative control against `3.82.0`. Both
  are named as residual evidence gaps in §ACs and §Debt, not left for β to
  discover unprompted.
- **Is every claim backed by evidence in the diff / a run command?** Yes —
  every AC section above cites either an actually-run local command (with
  its result), a directly-fetched external source (the action's own JS),
  or a specific historical log line, rather than an assumption about how
  GitHub Actions or `action-gh-release` behaves.
- **Scope guardrails honored?** Confirmed no changes to: any file under
  cnos#607's repo-installer path (`src/go/internal/repoinstall/`,
  `src/go/internal/cli/cmd_repo_install*.go` — not present in `git diff
  --stat`), any semver-range logic (none added — `install.sh`'s tooling
  branch does exact tag-prefix matching + lexicographic sort, no `^`/`~`
  ranges), any hosted registry (none added — GitHub Releases + the GitHub
  REST API are the only two backends used, both pre-existing), and the
  D14-held feature-version tag (not touched; the tooling-channel tag
  format `tooling-<date>-<sha>` is visibly distinct from both `X.Y.Z` and
  `vX.Y.Z`, and doesn't match either push-trigger glob
  `[0-9]*.[0-9]*.[0-9]*` / `v*`, so a tooling dispatch can never free-ride
  the tag-push trigger — verified by inspection of the two glob patterns
  against the naming convention).
- **Implementation contract (7 axes, `gamma-scaffold.md` table)
  conformance:**
  | Axis | Pinned | Honored? |
  |---|---|---|
  | Language | YAML + POSIX sh, no new Go/OCaml | Yes — only `.yml` and `.sh` touched |
  | CLI integration target | N/A | N/A — no `cn` subcommand touched |
  | Package scoping | `.github/workflows/release.yml`; `install.sh` | Yes — exactly these two, plus one additive doc note (explicitly optional per scaffold) |
  | Existing-binary disposition | preserve, extend don't replace | Yes — no rewrite of the pipeline/installer, only additive input/branch logic |
  | Runtime dependencies | none new | Yes — no `jq`/new binary added to `install.sh`; only `curl`+`grep`+`sed`+`sort`, all already implicitly required (curl explicitly checked for already; grep/sed/sort are POSIX baseline, already used elsewhere in the same file) |
  | JSON/wire contract preservation | `index.json`/`checksums.txt`/tarball formats unchanged | Yes — not touched at all |
  | Backward-compat invariant | default install.sh behavior unchanged; channel selector additive/opt-in | Yes — verified by the driver-script test showing the default (unset `CNOS_CHANNEL`) path is byte-identical logic to before, and by `git diff` showing that branch is a verbatim carry-over into a `case` arm |

## Debt

1. **No live GitHub Actions run performed or observable in this
   environment.** AC1, AC2 (full oracle command), AC3, and AC4 all have a
   local-reasoning/local-test ceiling — none can be *fully* closed without
   an actual `workflow_dispatch` run against `usurobor/cnos` (for AC1/AC3),
   an actual tag push (for AC4), and an actual `install.sh` run against a
   real published tooling tag plus the `3.82.0` negative control (for
   AC2). This is the single largest residual item and is squarely a β-side
   (or human-operator) verification task per the scaffold's own β prompt
   ("re-running the mechanical checks yourself ... do not treat α's
   self-coherence.md claims as substitute evidence").
2. **AC4 cannot be proven fixed with certainty** — the best available
   evidence (this cycle's log analysis + cnos#429's own prior RCA)
   concludes the original failure was *likely* transient/GitHub-side
   rather than a static misconfiguration. The `fetch-tags: true` change
   removes the specific request shape implicated in the one failing log
   this repo has, and the explicit auth pins are genuine hardening — but
   if the failure truly is a rare GitHub-backend flake, no client-side
   YAML change can guarantee it never recurs. This is disclosed here
   rather than claimed as a definitive fix.
3. **AC2 oracle command as literally written in the scaffold has an
   env-var-placement bug** (`VAR=val curl ... | sh` doesn't propagate
   `VAR` past the pipe) — documented in §ACs and in `install.sh`'s own
   header comment; not a defect in this cycle's diff, but worth β/operator
   awareness so a literal copy-paste of the scaffold's oracle command
   isn't mistaken for a failing (or vacuously passing) test of this
   cycle's actual mechanism.
4. **`install.sh`'s tooling-channel API call is unauthenticated** (60
   req/hour/IP GitHub rate limit) and uses `per_page=100` without
   pagination — acceptable for the expected release volume and consistent
   with the existing script's already-unauthenticated GitHub interactions,
   but would need revisiting if the release count or call volume grows
   substantially. Not raised by the scaffold as in-scope; noted as minor
   forward-looking debt, not a defect against AC1–AC4.
5. **`docs/guides/BUILD-RELEASE.md`'s pre-existing OCaml/dune-era content
   was not rewritten** — flagged inline instead (see the new note directly
   under the Overview table) per the scaffold's explicit "α is not
   obligated to fix the rest of that doc."
6. **Commit author identity is `sigma@cnos.cn-sigma.cnos`, not the
   `alpha@cdd.cnos` pattern named in `alpha/SKILL.md` §2.6 row 14.** This
   matches the identity already used for this cycle's prior commits on
   `origin/cycle/618` (`b9ae2222` dispatch-claim, `88a856f7` γ-scaffold —
   both under the same identity, evidently the session/environment
   identity for this wake-invoked run rather than a per-role one) and is
   the harness-provided git identity for this session; this session's
   tooling does not permit local git-config changes. Per row 14 path (b),
   disclosing explicitly rather than silently diverging from the named
   convention.
7. **Review-readiness signal not issued.** Per the wake instructions for
   this session, α's job ends at "implementation + self-coherence.md
   committed and pushed" — no `## Review-readiness` section is appended,
   and α does not poll for β nor spawn any other role. δ dispatches β
   next, per the operating instructions given for this session (this is a
   deviation from the full `alpha/SKILL.md` §2.7 "Request review" step,
   made explicit here because that skill file expects α to also signal
   review-readiness and begin polling; this session's wake instructions
   override that with "Do NOT spawn or invoke β yourself... Your job ends
   when your implementation + self-coherence.md are committed and
   pushed").
