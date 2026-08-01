# γ scaffold — cycle #618

**Issue:** [usurobor/cnos#618](https://github.com/usurobor/cnos/issues/618) — release-infra: FSM-capable `cn` tooling channel
**Mode:** design-and-build (cell_kind: implementation)
**Branch:** `cycle/618` (base `main@937f2755d02b8e2cfb35845ac7b78213e5c229b1`)

## Gap

`install.sh` resolves whatever GitHub currently reports as the repo's
"latest" (non-prerelease) release, and the only way to put a *new* binary
there today is either (a) cut a full, D14-held feature version, or (b) hand-patch
a one-off release outside CI, as already happened twice (see friction note 6).
Separately, the `release.yml` publish pipeline has two independent defects:
its `workflow_dispatch` path cannot target an explicit tag (`action-gh-release`
publishes against `github.ref`, which for a manual dispatch is a branch, not a
tag), and its tag-push path has a live, previously-observed checkout-auth
failure. Together these mean there is no *repeatable* way to ship the current
`main`'s `cn` (with the FSM verbs and, now, `cn repo install`) to a tenant
running the documented installer, independent of the feature-version release
cadence. This cycle's job is to add that repeatable path — not to cut a
release itself, and not to touch the feature-version cadence or the repo
installer's own internals.

## Per-AC oracle list

### AC1 — tooling-channel publish

**Oracle:** Trigger `release.yml` via `workflow_dispatch` with
`smoke-only=false`, an explicit `tag` input (e.g. `tooling-20260709-937f275`
— see friction note 5 for the naming convention), and a `prerelease=true`
input (α adds this input — see Implementation contract). After the run
completes, `gh release view <tooling-tag> --repo usurobor/cnos --json
tag_name,isPrerelease,assets` must show: `tag_name` equal to the given tag,
`isPrerelease: true`, and `assets` containing `cn-linux-x64`,
`cn-linux-arm64`, `cn-macos-x64`, `cn-macos-arm64`, `checksums.txt`,
`index.json`, and the `cnos.*` package tarballs — all published **without**
bumping or touching the D14-held feature version tag.

### AC2 — install.sh resolves an FSM-capable cn

**Oracle:** After AC1's publish, resolve the tooling channel via the
documented selector (α implements — see Implementation contract / friction
note 1) and confirm the resolved binary lists the FSM verbs:

```sh
tmp=$(mktemp -d)
CNOS_CHANNEL=tooling curl -fsSL https://raw.githubusercontent.com/usurobor/cnos/cycle/618/install.sh | BIN_DIR="$tmp/bin" sh
"$tmp/bin/cn" --help | grep -E 'issues fsm|cell finalize'
```

Both verbs must appear. **Negative control:** run the same channel-selector
logic pointed at a known pre-FSM tag (e.g. `3.82.0`, 2026-05-24) and confirm
the grep does **not** match — this proves the check discriminates rather than
passing vacuously. β must run this against the *new* tooling-channel
mechanism this cycle produces, not against the pre-existing hand-patched
3.82.2/3.82.3 releases (see friction note 6 — those already happen to pass a
naive version of this check today, which is exactly the fragile state AC1–AC4
exist to replace).

### AC3 — tag-targeted dispatch publish

**Oracle:** `gh workflow run release.yml --repo usurobor/cnos -f
smoke-only=false -f tag=<explicit-tag>` (dispatched from `main`, i.e.
`github.ref_name` at dispatch time is `main`, not `<explicit-tag>`). After
the run completes, `gh release view <explicit-tag>` must show the release
was created **at** `<explicit-tag>` — not silently attached to a "main"
pseudo-tag, and not failing with a ref/tag mismatch. This is the direct
mechanical fix for the `action-gh-release` sees `github.ref =
refs/heads/main`, not a tag" defect named in the issue body.

### AC4 — tag-push auth fixed

**Oracle:** a real lightweight/annotated tag push matching the existing
trigger patterns (`[0-9]*.[0-9]*.[0-9]*` or `v*`) triggers `release.yml`,
and the `build` job's checkout step completes without the error observed at
[run 26340314834](https://github.com/usurobor/cnos/actions/runs/26340314834)
(tag `3.82.0`, 2026-05-23T18:30Z): `fatal: could not read Username for
'https://github.com': terminal prompts disabled` (three retries, all
failed, in the `build (ubuntu-latest, linux-x64)` matrix leg's checkout
step). The run must reach `publish` green. See friction note 4 — α's first
step is confirming this specific failure still reproduces before attempting
a fix.

## Source-of-truth table

| Surface | File / ref | Authoritative for |
|---|---|---|
| Release/publish pipeline | `.github/workflows/release.yml` | build matrix, publish job, `workflow_dispatch` inputs, tag-push trigger patterns |
| Installer resolution | `install.sh` | which release `install.sh` downloads from; checksum verification; channel selection (new, this cycle) |
| External smoke oracle | `scripts/smoke/90-release-bootstrap.sh` | what "a working release" mechanically means (index.json + tarballs + `cn` binary + bootstrap chain) — reference for shape, not required to modify |
| Customer source | [cnos#606](https://github.com/usurobor/cnos/issues/606) C1/C2 | the underlying tenant-facing gap this cycle closes |
| Installer wave (dependency; **do not touch**) | [cnos#607](https://github.com/usurobor/cnos/issues/607) | the repo installer (`cn repo install`) itself — out of scope |
| Prior one-off workaround (informative only, not authoritative for this cycle's shape) | tags `3.82.1`/`3.82.2`/`3.82.3`; commit `1ebdb86a` "fix(repo install): resolve package versions from release index, not the release tag"; commit `fbb64bf9` "one-off publish release 3.82.3" | evidence that a manual, non-repeatable publish path already exists and must be *replaced* by a supported one — not re-run as the fix, not merged as part of this cycle |

## Surfaces α is expected to touch

- **`.github/workflows/release.yml`**
  - `workflow_dispatch.inputs`: add a `prerelease` boolean input (default
    `false`); keep the existing `smoke-only` and `tag` inputs.
  - `publish` job: pass `tag_name: ${{ inputs.tag }}` (when non-empty) and
    `prerelease: ${{ inputs.prerelease }}` to both `softprops/action-gh-release@v1`
    steps, so a `workflow_dispatch` publish targets the given tag instead of
    `github.ref` (AC3).
  - `build` job: the `-ldflags "-X main.version=..."` stamp currently uses
    `github.ref_name`, which for a `workflow_dispatch` run is the dispatched
    branch (e.g. `main`), not the target tag — likely needs to prefer
    `inputs.tag` when set, so the shipped binary's stamped version matches
    the tag it's published under (affects `cn setup`'s default `deps.json`
    pin, per the existing code comment at that step).
  - `build` job's `actions/checkout@v4` step (tag-push path): fix per AC4.
    α owns root-causing whether the fix is `persist-credentials`, an
    explicit `token:`, job-level `permissions`, or something else — the
    scaffold does not pre-select the mechanism.
- **`install.sh`**: add a documented, opt-in channel selector (e.g. a
  `CNOS_CHANNEL` env var) that, when set to the tooling channel, resolves
  the newest tooling/prerelease tag via the GitHub API instead of the
  default `/releases/latest` redirect. **Default (unset) behavior must be
  unchanged** — this is additive, not a replacement of the stable path.
- Optionally: a short doc note (e.g. in `docs/guides/BUILD-RELEASE.md`, which
  is otherwise stale/OCaml-era — α is not obligated to fix the rest of that
  doc) describing the channel-selector convention. Non-binding; α's call.

## Scope guardrails (issue Non-goals — α MUST NOT touch)

- **`cnos#607` / the repo installer itself.** Any change to `cn repo
  install`'s internal package-resolution logic — including the fix already
  sitting on the side lineage at commit `1ebdb86a` — is out of scope. Do not
  merge, cherry-pick, or reimplement it as part of #618.
- **Semver-range resolution.** No fuzzy/range version matching (`^`, `~`,
  etc.) anywhere in `install.sh` or the workflow.
  - **Hosted package registry.** No new registry service. This cycle only
  touches GitHub Releases (already-existing mechanism) + `install.sh`.
- **The D14-held feature version cut.** Do not tag, bump, or publish the
  actual held feature version. The tooling-channel tag must be visibly
  distinct from both `X.Y.Z` and `vX.Y.Z` feature-version tags.

## α prompt

```text
You are α. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 618 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/618
Scaffold: .cdd/unreleased/618/gamma-scaffold.md — read it in full before coding.
It carries the per-AC oracle list (AC1-AC4), the source-of-truth table, the
scope guardrails (issue Non-goals), and friction notes recording judgment
calls already made (tooling-channel mechanism, tag-naming convention). Do
not re-litigate those calls without a strong reason; if you must deviate,
name the deviation explicitly in self-coherence.md.
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | YAML (GitHub Actions workflow) + POSIX sh (install.sh); no new Go/OCaml source expected |
| CLI integration target | N/A — CI workflow + installer script, not a `cn` subcommand |
| Package scoping | `.github/workflows/release.yml`; `install.sh` (repo root) |
| Existing-binary disposition | preserve — extend, do not replace, the existing release/publish pipeline and installer |
| Runtime dependencies | None new (existing `actions/checkout@v4`, `softprops/action-gh-release@v1`, `actions/upload-artifact@v4`, GitHub REST API already used by install.sh) |
| JSON/wire contract preservation | preserve as-is — `index.json`, `checksums.txt`, package tarball formats unchanged |
| Backward-compat invariant | default (no env var set) `install.sh` behavior resolves identically to today's stable-latest release; any channel selector is additive/opt-in only |

Implement AC1-AC4 exactly as named in the scaffold's oracle list. When done,
write .cdd/unreleased/618/self-coherence.md with a `§R0` section (AC-by-AC
mechanical evidence + Implementation-contract conformance) and signal
review-readiness per your role skill. Do NOT spawn α/β/γ/δ or any other
role — you exit after signaling; δ dispatches β next.
```

## β prompt

```text
You are β. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 618 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/618
Scaffold: .cdd/unreleased/618/gamma-scaffold.md — read the per-AC oracle
list before reviewing. Walk AC1-AC4 independently by re-running the
mechanical checks yourself (dispatch runs, gh release inspections, the
install.sh channel-selector probe including the negative control) — do not
treat α's self-coherence.md claims as substitute evidence. Confirm the
Implementation contract's 7 axes were honored against the diff. Confirm the
scope guardrails (issue Non-goals) were not violated — in particular that
cnos#607/repo-installer internals were not touched.

Write .cdd/unreleased/618/beta-review.md with a `§R0` section recording your
verdict (`converge` or `iterate`) and any findings. Do NOT spawn α/β/γ/δ or
any other role — you exit after writing your verdict; δ routes next.
```

## Friction notes (judgment calls made at scaffold time)

1. **What "tooling-channel" concretely means.** The issue leaves this open
   ("a tooling-channel or prerelease"). Decision: implement it as a **GitHub
   prerelease published via an explicit `workflow_dispatch` tag input**, not
   as a `git push` of a matching semver tag. Rationale: this cleanly
   separates AC1-3 (dispatch-driven, no full version cut, and — because
   GitHub's `/releases/latest` redirect excludes prereleases by definition —
   invisible to `install.sh`'s existing default resolution) from AC4 (fixing
   the pre-existing tag-push path used for real version cuts). The
   prerelease flag is exactly the mechanism that makes the channel additive
   rather than a change to default behavior.
2. **Consequence of (1):** because prereleases are excluded from
   `/releases/latest`, `install.sh`'s default-path resolution is *structurally*
   unchanged by this cycle — the channel selector required for AC2 must be
   an explicit opt-in (e.g. an env var), not a change to the redirect-based
   default path.
3. **Out-of-scope observation, not to be acted on.** `main` at
   `937f2755d02b8e2cfb35845ac7b78213e5c229b1` does **not** contain commit
   `1ebdb86a` ("fix(repo install): resolve package versions from release
   index, not the release tag") — that fix exists only on the side lineage
   that produced the manual `3.82.2`/`3.82.3` releases. Whether/when that
   fix lands on `main` is cnos#607-wave territory; α MUST NOT merge or
   cherry-pick it here even though it is adjacent to this cycle's release
   mechanics.
4. **AC4 reproduction caveat.** The concrete historical auth failure
   (run 26340314834, tag `3.82.0`, 2026-05-23) is real, but the most recent
   tag-push run (`3.82.2`, 2026-07-08,
   [run 28966706260](https://github.com/usurobor/cnos/actions/runs/28966706260))
   failed for an unrelated reason (a `macos-arm64` `go test` failure in the
   `build` job, not a checkout-auth error). α's first step for AC4 should be
   attempting a real (or rehearsal) tag push to confirm the checkout-auth
   failure still reproduces on current `main` before diagnosing a fix — the
   bug may be intermittent (GitHub-side) rather than a static workflow
   misconfiguration, in which case the fix should still harden the checkout
   step defensively (explicit `persist-credentials`/`token`) even if a single
   rehearsal happens not to reproduce it.
5. **Tooling-tag naming convention.** Chosen: `tooling-<UTC-date>-<short-sha>`
   (e.g. `tooling-20260709-937f275`) — visually distinct from both `X.Y.Z`
   and `vX.Y.Z` feature tags at a glance, monotonically sortable by date,
   and guaranteed not to collide with the push-trigger tag globs
   (`[0-9]*.[0-9]*.[0-9]*` / `v*`) so a tooling-channel dispatch never
   free-rides the tag-push trigger path. α may choose a different concrete
   naming scheme with a stated reason, but must preserve the
   non-collision property.
6. **`install.sh` already "passes" AC2 today, by accident.** As of this
   scaffold's authoring (2026-07-09), `install.sh`'s existing "latest release"
   redirect resolves to `3.82.3`, a hand-patched release (commit `fbb64bf9`,
   "one-off publish release 3.82.3") that does include `cn repo install` and
   the FSM verbs — but that release is not reachable through a repeatable
   pipeline and its commit lineage is not on `main`. β must verify AC2
   against the *new* mechanism this cycle produces (a fresh tooling-channel
   publish), not accept the pre-existing accidental pass as evidence.
