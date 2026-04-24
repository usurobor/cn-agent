# RELEASE.md

## Outcome

Coherence delta: C_Σ **A-** (`α A-`, `β A-`, `γ A-`) · **Level:** **L6**

Skill activation moves from OCaml-only to Go, with the single source of truth for skill existence becoming filesystem presence (`<pkg>/skills/<id>/SKILL.md`) rather than a declared manifest field. `visibility: internal` frontmatter — added on 9 CDD sub-skills in 3.56.x but unenforced — is now effective at the runtime activation boundary. `cn doctor` gains a skill-activation check wired over the new Go discovery path, with a severity model where unreadable SKILL.md files fail the hub (structural break) while overlapping trigger keywords across public skills surface as informational warnings (many-to-many hints, not exclusive dispatches).

## Why it matters

Three latent problems close in the same cycle:

1. **Authority drift.** `cn.package.json` stopped declaring `sources.skills` in 3.56.x (ef53b939), but the OCaml activation validator iterated only manifest-declared IDs, so it silently turned into a no-op for cnos.core — any authoring problems in 19 core skills would never surface at `cn doctor`. The new Go path walks the filesystem (the authority `DESIGN-CONSTRAINTS.md §1` names) and evaluates every SKILL.md that actually exists, so authoring problems surface immediately.

2. **Visibility was a docstring, not a contract.** `visibility: internal` was present on 9 CDD sub-skills since 3.55 but had no runtime effect — the OCaml index didn't consume the field because it never saw the sub-skills at all. `BuildIndex` now consumes it, excludes internal skills from the public activation table, and also excludes them from trigger-conflict detection (they aren't addressable, so they can't ambiguate activation).

3. **Doctor severity matches the activation model.** The first end-to-end run against the real corpus surfaced that cnos skill keywords are intentionally overlapping across public skills — `coherence`, `verify`, `sync`, `reflect`, etc. are many-to-many activation hints rather than 1:1 dispatches. Rather than fail the hub on every overlap, doctor reports them as ○ warnings (operator-visible without escalation). Only structural breakage — a SKILL.md that cannot be read or parsed — fails.

## Added

- **`src/go/internal/activation/`** — new Go package. `frontmatter.go` (pure, no IO) parses YAML-subset frontmatter with both block-list and inline flow-sequence triggers, CRLF/BOM handling, malformed-line tolerance. `index.go` (IO) exposes `Discover` (filesystem walk), `BuildIndex` (public filter), `Validate` / `ValidateSkills` (doctor-facing issue generation).
- **Skill-activation doctor check** in `internal/doctor/doctor.go` — reports `[missing]` / `[empty]` / `[conflict]` issues, mapping only `[missing]` to `StatusFail`.
- **Real-corpus integration tests** — `TestDiscover_RealCoreSkills_HaveTriggers` + `TestValidate_RealCorpus_NoEmptyTriggers` install unmodified `src/packages/*/skills/**/SKILL.md` into a temp hub and assert the full installed corpus validates clean.

## Changed

- **`docs/alpha/package-system/PACKAGE-AUTHORING.md §8`** — rewrites the manifest-contract paragraph. Manifest owns `commands` and `engines`; filesystem owns content classes. Drops "if it's not in the manifest, the runtime doesn't know about it" — inverted now: if it's on disk, the runtime knows.
- **`docs/alpha/package-system/PACKAGE-ARTIFACTS.md`** — replaces `sources.skills` examples with filesystem-walk + `visibility: internal` narrative; the `cnos.cdd` worked example is rewritten end-to-end.
- **Skill-activation authoring checklist** — "every declared skill directory contains a SKILL.md" becomes "every skill directory under `skills/` contains a valid SKILL.md (non-empty `triggers:`; `visibility: internal` when the skill is internal to its package's orchestrator)".
- **Trigger-conflict severity** — public-public overlap is now an informational warning rather than a hub failure. Divergence from OCaml `cn_doctor.ml` severity mapping is intentional and ratified by the operator; downstream tooling that grep-ed for `✗ skill activation` on conflicts will see `○ skill activation N warning(s): [conflict] ...` instead.

## Fixed

- **Inline-list trigger parsing** — 42 of 52 production SKILL.md files used `triggers: [a, b, c]` inline-flow format, which the initial Go parser dropped. `parseInlineList` handles the dominant production form with whitespace tolerance, bare-identifier items, and empty-list (`[]` → nil) semantics.
- **Internal skills no longer ambiguate activation detection** — an internal sub-skill sharing a trigger keyword with its public parent orchestrator is not a conflict (the runtime cannot route to internal skills).

## Validation

- `go test ./...` green across all 11 Go packages (33 activation tests, 10 doctor tests).
- Local reproduction of `kata-tier2`: `cn init ci-hub && cn setup && cn deps lock && cn deps restore && cn doctor` exits rc=0 with `○ skill activation 9 warning(s)` for the overlap hints — pre-break baseline required by `R3-doctor-broken` now holds.
- `cn kata run --class runtime` passes all 4 runtime katas locally on 7e76798.
- CI on PR #263: `go`, `kata-tier1`, `kata-tier2`, `Package/source drift (I1)`, `Protocol contract schema sync (I2)` all green; only `Package dist/source sync (I3)` red (deferred per #264).
- The targeted coherence delta — `visibility: internal` being enforced at the runtime activation boundary rather than being a docstring — is proven by `BuildIndex` excluding internal skills from the public index in unit tests, and by the real-corpus integration tests walking the full `src/packages/` tree into a temp hub.

## Known Issues

- **#264** — `cn build` produces non-deterministic tarball SHA-256s because `createTarGz` captures live file ModTime and gzip writer time. `Package dist/source sync (I3)` CI check is red on this release as a result; it was red on main pre-#261 and is not a regression of this cycle. Fix scoped in #264 with concrete patch. Deferred per review §7.0 design-scope exception.
- **Cnos.core trigger overlap corpus** — 9 overlapping trigger keywords across public skills (`alignment`, `boundary`, `coherence`, `drift`, `onboard`, `reflect`, `self-check`, `sync`, `verify`) are intentional per the many-to-many model; whether any future cycle wants to disambiguate them (or introduce activation precedence) is a doctrine question not scoped to #261.
