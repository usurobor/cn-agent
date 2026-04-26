# α Close-Out — 3.59.0 (#230, PR #274)

## Cycle summary

| | |
|---|---|
| Issue | #230 — *B1: distribution chain honesty — restore version-skip + release-bootstrap smoke*. Selected by γ from the 3.58.0 PRA encoding-lag table (issue listed as "growing"). Subsumes #238. Closed on merge of PR #274 (`Closes #230`). |
| PR | #274 — *α #230: distribution chain honesty — version-skip + release-bootstrap smoke*. Squash-merged to main as `9980e3f`. |
| Release | 3.59.0 — minor bump from 3.58.0. β tagged and shipped as commit `9dd30d9` (`release: 3.59.0 — distribution chain honesty (#230, PR #274)`). β close-out at `.cdd/releases/3.59.0/beta/CLOSE-OUT.md`. |
| Dispatch | γ → α: project `cnos`, issue #230, Tier 3 skills `eng/go`, `eng/test`, `eng/tool`. Branch `claude/alpha-tier-3-skills-IZOsO` (non-canonical per CDD §4.2; harness-assigned, not α-chosen — β flagged for γ). |
| Work shape | Substantial runtime + tooling fix (MCA, L6). 8 files in PR #274 (6 modified + 2 new); +784 / −22 net; 2 commits on-branch (`6cbd225` initial implementation, `93ea1d6` round-1 fix). |
| Tier 2 skills applied | Drawn from `eng/{go,test,tool}` per the issue's Tier 3 list — that set already covers the bundle (`eng/go` runtime work, `eng/test` proof-depth, `eng/tool` for the smoke script template). |
| Tier 3 skills applied | `eng/go` (parse/read split honored: pure parser in `internal/pkg`, IO wrappers in `restore` and `doctor` consume it; structured `slog.WarnContext` on every degraded path); `eng/test` (negative-space coverage — version drift, unreadable manifest, half-state, unparseable JSON; anti-overlay assertion via v1-only marker file); `eng/tool` (smoke script: `set -euo pipefail`, NO_COLOR, prereq checks, exit codes 0/1/2 with documented meaning, machine-readable trailing `RESULT:` line). |
| MCA selection | Issue named the gap as a two-source disagreement (lockfile = integrity authority, installed `cn.package.json` = installed-version authority) and listed the canonical fix path. α's job was to land the surgical patch (version-aware skip + reinstall) plus the new gate (release-bootstrap smoke + CI workflow), under DESIGN-CONSTRAINTS §1/§2.2/§3.2/§6.3. The L6 leverage move is making the runtime surface (doctor) report the same disagreement that drove the silent skip — so a stale vendor cannot stay invisible after a missed restore. |
| Review rounds | **2** — R1 REQUEST CHANGES + 4 findings (F1 C, F2 B, F3 B, F4 A); R2 APPROVED at `93ea1d6`. β merged at `9980e3f` after CI green on round 2. |
| Findings against α | 4 total — F1 (C, judgment, doctor parser duplication); F2 (B, judgment, smoke header lied about checksums authority); F3 (B, judgment, untested `(unparseable manifest)` branch); F4 (A, mechanical, dead `BIN_VERSION_OUT` capture). All resolved on-branch in `93ea1d6` per `review/SKILL.md` §7.0. No D-level findings, no AC misses. |
| Mechanical ratio | 1/4 = 25% (F4). Above the §9.1 20% threshold but with N=4 ≪ 10-finding floor; not an automatic process-debt trigger. |
| Local verification | `go build ./... && go vet ./... && go test ./...` green on every review-round head; `go test -race` on touched packages green; `bash -n` on smoke clean; offline probe exits 2 with `RESULT: skipped (offline)`; both YAML workflows parse via `yaml.safe_load`. CI on merged head `9980e3f` (round 2): 7/7 green (go, kata-tier1, kata-tier2, I1 package/source-drift, I2 protocol-contract-schema-sync, 2× notify). |
| Concurrent main activity | None during review. Last main commit pre-cycle was `eafc230` (γ #268 close-out, 2026-04-25); main did not advance during R1 → R2 → merge. |

## Findings (α-side observations)

### Finding F1 — doctor parser duplication (C, judgment)

**What happened.** Round-1 `doctor.checkPackages` parsed installed `cn.package.json` via an inline anonymous-struct unmarshal:

```go
var pm struct {
    Version string `json:"version"`
}
if err := json.Unmarshal(manifestData, &pm); err != nil { ... }
```

The same PR introduced `pkg.ParseInstalledManifestData([]byte) (PackageManifest, error)` as a pure parser, with `restore.ReadInstalledManifest` as its IO wrapper — exactly the surface doctor needed. Two parsers for the same fact ("installed package version") landed in the same diff.

**Pattern.** `eng/go` §2.17 (Parse vs Read: pure functions on bytes, IO on paths, no mixing) + `design/SKILL.md` §3.2 (one source of truth per fact) + `cdd/review` §2.2.8 (authority-surface conflict). I introduced the pure parser specifically because it was needed in two places (restore + doctor), then hand-rolled an inline unmarshal in the second place anyway. The skill was loaded; the skill named the rule; I broke the rule in the same diff that obeyed it on the adjacent surface.

**Resolution.** Round 2 commit `93ea1d6`: import `internal/pkg` into `doctor.go`, replace the inline block with `pm, err := pkg.ParseInstalledManifestData(manifestData)`, and add a comment citing eng/go §2.17 + design §3.2 so the constraint is locally visible. No new import cycle (`activation/index.go` already imports `pkg`). Net change: ~5 lines, behavior preserved (same diagnostic strings, same `StatusFail` mapping).

**Surfaces involved.** `internal/doctor/doctor.go` (the inline parse, now removed); `internal/pkg/pkg.go` (the pure parser, now consumed by both call sites); `internal/restore/restore.go` (the IO wrapper, unchanged).

### Finding F2 — smoke header lied about integrity authority (B, judgment)

**What happened.** The smoke script's file header listed the verification step as "verify tarball SHA-256 against `checksums.txt`," and the script downloaded `checksums.txt` from the release. The actual SHA verification compared each tarball against `index.json[name][version].sha256`, never reading `checksums.txt`. The downloaded file was the binary-asset checksums file (covering `cn-<platform>` builds, generated by `release.yml`'s `Generate binary checksums` step), not the tarball authority.

The implementation was correct — `index.json` sha256 is what `cn deps restore` itself trusts at runtime, so the smoke verifying against it mirrors production behavior. The header was wrong, and the irrelevant download was a moving part with no contract.

**Pattern.** Header-as-contract drift, plus YAGNI on the second authority. I wrote the header at script-design time intending to use `checksums.txt` as a cross-confirmation, then implemented against `index.json` (correctly) and forgot to revise the header. β's grep across script and header found the disagreement immediately.

**Resolution.** Round 2 commit `93ea1d6`: header rewritten to "verify each tarball's SHA-256 against the index.json entry — the same authority `cn deps restore` itself trusts at runtime." `checksums.txt` download removed; replaced by a 4-line comment explaining the binary-asset role of that file and why it is not the tarball authority. Single integrity authority, one moving part, header tells the truth.

**Surfaces involved.** `scripts/smoke/90-release-bootstrap.sh` (header + download block + SHA-verification loop). No CI change; the workflow only invokes the script.

### Finding F3 — untested `(unparseable manifest)` branch (B, judgment)

**What happened.** `doctor.checkPackages` round-1 emitted three stale-state diagnostics: `(installed X, locked Y)` (version drift), `(no manifest)` (vendor dir present but `cn.package.json` missing), and `(unparseable manifest)` (manifest present but JSON malformed). The PR added `TestRunAllPackageVersionDrift` and `TestRunAllPackageManifestMissing`, covering two of three branches; the third had no test.

**Pattern.** `eng/test` §2.7 (negative space is mandatory) + §3.6 (every meaningful boundary gets a "must surface" test) + §3.13 (cover new surfaces, not just preserved ones — each new diagnostic string is a new surface element). Three diagnostics ⇒ three tests. I covered the two named in conversation with the AC and missed the third because I was thinking of "drift" as the AC's headline case and "no manifest" as the obvious half-state, without enumerating the diagnostic-string family as the test boundary.

**Resolution.** Round 2 commit `93ea1d6`: `TestRunAllPackageManifestUnparseable` writes `{not valid json` as the manifest body, asserts `StatusFail` with `"unparseable manifest"` in the diagnostic. Mirrors the shape of `TestRunAllPackageManifestMissing`; `~15` lines net. All four package-related doctor tests now pass: `Missing`, `VersionDrift`, `ManifestUnparseable`, `ManifestMissing`.

**Surfaces involved.** `internal/doctor/doctor_test.go` (test added); the underlying production code (`doctor.go`'s third branch) was already correct, just unproven.

### Finding F4 — dead `BIN_VERSION_OUT` capture (A, mechanical)

**What happened.** The smoke script captured the binary's status output:

```bash
if ! BIN_VERSION_OUT="$("$BIN" status 2>&1 || true)"; then
  : # status fails outside a hub — that is fine; we just want it to execute.
fi
```

`BIN_VERSION_OUT` was never read again. The `|| true` masked the exit code, so the `if !` branch could not fire — the body was unreachable. The intent was "verify the binary executes," which a no-capture form expresses without dead code.

**Pattern.** Mechanical leftover from an earlier draft where the captured output was inspected; the inspection was dropped during a smoke-script simplification pass and the capture wasn't pruned with it. β's grep over the script found the unused variable and the unreachable branch.

**Resolution.** Round 2 commit `93ea1d6`: replaced lines 127–129 with `"$BIN" status >/dev/null 2>&1 || true` and a 4-line comment naming what the line guards against (missing loader, wrong arch, exec bit lost in transit). Same intent, no captured-but-unused variable, no unreachable code.

**Surfaces involved.** `scripts/smoke/90-release-bootstrap.sh` only.

## Cycle Iteration (CDD §9.1)

### Triggers fired

- [ ] review rounds > 2 — actual: 2. At threshold, not above.
- [ ] mechanical ratio > 20% with ≥ 10 findings — actual: 25% (1/4) but N=4 ≪ 10-finding floor; below the §9.1 automatic-trigger threshold.
- [ ] avoidable tooling / environmental failure — none. Local sandbox lacks outbound network to `api.github.com`; the smoke's offline-graceful path was an explicit AC, so the constraint produced an on-design exit code (2) at every local run, not a friction event.
- [x] **loaded skill failed to prevent a finding** — `eng/go` was a declared active skill. §2.17 (Parse vs Read) names the rule directly: pure parsers on bytes, IO wrappers on paths, and (per §3.9 smell list) "producer-owned interfaces with one implementation" / parallel parsers for the same fact. F1 was specifically a parallel parser landed in the same diff that introduced the canonical pure parser.

β's close-out independently records this same trigger.

### Friction log (α-side)

1. **Asymmetric application of the parse/read rule across two consumers in one diff.** I introduced `pkg.ParseInstalledManifestData` because two surfaces (`restore.restoreOne`, `doctor.checkPackages`) needed to read the same `cn.package.json` `version` field. I wired the first consumer canonically (`restore.ReadInstalledManifest` IO wrapper → `pkg.ParseInstalledManifestData` pure parser). I wrote the second consumer (`doctor.checkPackages`) with an inline anonymous-struct unmarshal, never importing `pkg`. The two edits happened in the same authoring session, ~30 lines apart in conceptual ordering. The skill (`eng/go` §2.17) was loaded; the canonical exemplar was already on the keyboard; the second consumer still hand-rolled the parse.

2. **Negative-space coverage planned by AC, not by surface enumeration.** I wrote two doctor tests (`TestRunAllPackageVersionDrift`, `TestRunAllPackageManifestMissing`) for the two diagnostics I had been thinking about while writing the AC text — "drift" (the headline case) and "no manifest" (the obvious half-state). I did not enumerate the diagnostic-string family as a test boundary, so the third diagnostic (`unparseable manifest`) shipped without coverage. F3 is the direct consequence. The skill (`eng/test` §3.6 + §3.13) covers this, but my test-planning surface was AC-shaped, not diagnostic-string-shaped.

3. **Header-as-contract drift on a draft-then-revise edit.** F2: I drafted the smoke header listing `checksums.txt` as the verification authority, then implemented the verification against `index.json` (correctly), then forgot to revise the header. The header → implementation drift is one of the same shapes I caught and called out earlier in this PR for `pkg.PackageManifest` (additive, non-breaking) — except in that case I named the constraint up-front ("schema audit"), and in the smoke-header case I had no schema-audit-equivalent for "the header is a contract for the script body."

4. **Mechanical leftover from a draft simplification.** F4: dead `BIN_VERSION_OUT` capture was a leftover from an earlier draft where the captured output was inspected (e.g. log it, parse a version line). The inspection got dropped during simplification; the capture survived. No skill names the "delete inspection, also delete capture" pairing as a rule; the smell is dead code, which any review pass catches but which my pre-review-gate row 9 (post-patch re-audit) missed because the smoke is shell, not Go, and my Go-test-driven post-patch loop does not exercise `set -u` against the shell unused-variable case.

### Root cause (α-side reading)

**For F1 (the C-severity miss): proximity-to-canonical-exemplar did not transfer.** I had the canonical exemplar (`restore.ReadInstalledManifest`) on the keyboard one file away from where I wrote the inline parse. The skill (`eng/go` §2.17) was loaded. The PR body §Mode listed `eng/go` as an active skill and the "parse/read split" as the named application. None of those was sufficient at authoring time. The pattern I exhibited is: a skill whose abstract rule is well-internalized produces correct individual call-sites but does not automatically generalize to "if I have just introduced a canonical X, every other consumer of X-shaped data in this diff must consume the canonical X." The generalization lives one step removed from §2.17's literal prose; β named this in their close-out independently.

**For F2/F3/F4 (the B/B/A misses): pre-review-gate row 9 (post-patch re-audit) was Go-shaped.** My re-audit loop after every change was `go vet + go test ./... + go test -race`. That covers Go. It does not cover header-contract drift (F2), test-surface enumeration completeness (F3), or shell unused-variable / unreachable-branch detection (F4). The pre-review-gate text in `alpha/SKILL.md` §2.6 names "post-patch re-audit" as a row but does not name the audit's coverage requirement against the diff's actual language mix. My re-audit was language-monoculture; the diff was polyglot (Go + shell + YAML + Markdown).

### Engineering level (α-side reading, for γ to adjudicate)

- **L5 (local correctness):** met on the merged diff. `go build ./... && go vet ./... && go test ./... && go test -race ./internal/{pkg,restore,doctor}/...` green on every review-round head; CI 7/7 green on round-2 head `93ea1d6`. Smoke `bash -n` clean, offline branch exits 2 cleanly, both YAML files parse via `yaml.safe_load`. No code semantic finding from β.

- **L6 (system-safe execution):** partial miss. Four findings reached β review — one C (cross-surface authority duplication, F1), two B (cross-surface header drift F2 + negative-space test gap F3), one A mechanical (F4). The C-severity miss is a real cross-surface coherence failure, not a hostile-only-with-careful-review nit; the same diff that introduced the canonical parser landed a parallel one. Caught at R1, all four resolved in one round-2 commit, no semantic ripple, no further findings — but L6's criterion is "cross-surface drift did not reach review," and one cross-surface drift did reach review.

- **L7 (system-shaping leverage):** achieved on the substantive axis. Pre-cycle, the install authority chain could lie silently in two places (lockfile bump silently ignored at restore time; runtime surface couldn't see the disagreement). Post-cycle, both lies are mechanically prevented: restore reads the installed manifest and reinstalls on drift with a structured `slog.WarnContext`; doctor reads the same manifest and reports `StatusFail` with named diagnostics. The release-bootstrap smoke + `release-smoke.yml` + `release.yml` ldflags-stamping form a new gate that catches the "released binary cannot self-bootstrap" failure class on every future tagged release. Three failure classes turned into mechanical fail-fasts.

Per §9.1 "lowest miss" rule: **L6 (with an L7 MCA shipped, capped by the L6 cross-surface duplication miss)**.

## What worked

1. **AC-by-AC self-coherence map in the PR body.** Every AC mapped to one named evidence anchor (test name, file path, or workflow trigger). β's R1 review used the same row structure to verify each AC met. Zero interpretation gap between issue, PR body, and review for the AC contract — all four R1 findings were on coherence dimensions (parser duplication, header truth, test coverage, dead code), not AC misses.

2. **Anti-overlay assertion in `TestRestoreVersionDriftReinstalls`.** The AC text was "install v1 → bump v2 → restore == v2." A weaker reading would extract on top of the v1 tree (v2 manifest wins; v1 stragglers stay). I added a `v1-only-marker` file pre-restore and asserted it was *gone* post-restore — proving the reinstall path actually wipes before extracting. This is one step stronger than the AC requires; β's review explicitly noted it as "anti-overlay test, stronger than AC requires."

3. **Verifying the offline branch locally.** I executed `scripts/smoke/90-release-bootstrap.sh` from the sandbox after writing it, observed the offline probe exit 2 with `RESULT: skipped (offline)`, and treated that as the operational proof for AC5 (graceful offline skip). Without this exercise I would have shipped an AC5 implementation that compiled but might have silently failed open. The local exercise confirmed the contract round-trip (probe → exit code → workflow `::warning::` mapping in `release-smoke.yml`).

4. **Catching the "binary version is `dev` in CI" gap.** While drafting the smoke, I traced the chain `cn setup` → `deps.json` → `cn deps lock` → `cn deps restore` and discovered that `release.yml` did not stamp `-ldflags` for the binary version. Without that stamping, `cn setup` would write `cnos.core@dev` and `cn deps lock` would fail against the released index. Adding the ldflags fix in the same PR meant the smoke I shipped could pass on the very next release; without it, the smoke would have correctly failed every release as a known limitation.

5. **Small round-2 diff scope.** Round 2 was strictly four named anchors — 3 files, 50 / 17 lines. No drive-bys, no scope creep, no new findings introduced. β's R2 review used "round-2 narrowing" as the explicit framing; the narrow patch let β re-read the cycle in one synchronous diff inspection rather than re-walking the full 8-file diff.

## What didn't

1. **F1 cross-surface duplication while the canonical exemplar was visible.** Discussed in §Cycle Iteration root cause. The skill was loaded, the canonical exemplar was a visible file away, and I still hand-rolled the parallel parser in the second consumer. This is the L6 miss this cycle exhibited.

2. **F2 header lied about the integrity authority.** Header was a draft-time intention; implementation diverged correctly; header was not revised. There is no schema-audit-equivalent rule in the pre-review gate for "header strings are a contract against the script body." (Pattern observation only — not a recommendation.)

3. **F3 negative-space coverage planned by AC, not by surface enumeration.** I tested the diagnostics I had been thinking about while writing; I did not enumerate the diagnostic-string family. `eng/test` §3.6 + §3.13 cover this; my application was AC-narrative-shaped.

4. **F4 dead-code shell leftover surviving simplification.** One A-mechanical leftover. Caught by β's grep over the script. My pre-review-gate row 9 ran `go vet` + `go test`; neither examines shell scripts.

5. **Polling discipline gap (β review absorbed by webhook, not surfaced by polling).** Per CDD §Tracking, polling is mandatory for shared-identity agents because webhook delivery is unreliable. I subscribed to PR webhook activity and set up a Monitor watching git-observable transitions (branch SHAs, main SHA), but I did not at first set up a baseline-pull on the PR comments / reviews surface. Result: β posted the round-1 review ~1.5h after my PR opened; the webhook did not deliver to my session; my git-only Monitor could not see comment activity; and I discovered the review only when the operator asked. The operator's intervention was the necessary signal. The fix during the cycle was to add the synchronous initial-state pull and to keep the Monitor narrowly scoped to git-observable transitions, with the understanding that comments / CI / reviews require periodic in-conversation MCP polling on top of the Monitor channel.

## Observations and patterns (α-side, factual only)

### 1. "Loaded canonical exemplar in the same diff" is a stronger constraint than "loaded skill"

**Pattern.** When a diff introduces a new canonical surface (here: `pkg.ParseInstalledManifestData`), every other consumer of the same fact in the same diff is structurally bound to the new canonical. The bind is one step removed from `eng/go` §2.17's literal "Parse vs Read" prose — §2.17 names the abstract rule (pure on bytes, IO on paths, no mixing) but does not name the diff-local generalization ("if X just got a canonical implementation in this PR, every other consumer of X in this PR consumes that implementation"). My F1 fit exactly that pattern.

**Surfaces this pattern touches.** `eng/go` §2.17 (parse/read split) and §3.9 (smell list — "producer-owned interfaces with one implementation"); `cdd/review` §2.2.8 (authority-surface conflict) — the cross-surface review check that caught it.

**Evidence.** N=1 (this cycle). β named the same generalization in their close-out independently.

### 2. Test-planning by AC narrative under-covers diagnostic-string families

**Pattern.** F3 was a negative-space gap on the third diagnostic-string branch (`unparseable manifest`). I tested the two diagnostics I had been thinking about while writing the AC (drift, no manifest); I did not enumerate the diagnostic family as the test boundary. The skill (`eng/test` §3.13 — "cover new surfaces, not just preserved ones") covers this if the surface enumeration is explicit; my surface enumeration was implicit-from-AC-narrative.

**Surfaces this pattern touches.** `eng/test` §2.7 (negative space mandatory) + §3.6 (every meaningful boundary gets a test) + §3.13 (cover new surfaces).

**Evidence.** N=1 (this cycle). The same shape ("each new diagnostic string is a new surface element to test") was visible to me at authoring time and I did not apply it.

### 3. Pre-review-gate row 9 ("post-patch re-audit") is language-monoculture in practice

**Pattern.** My re-audit loop was Go-shaped: `go vet`, `go test ./...`, `go test -race` on touched packages. That covers Go. It does not cover header-contract drift in shell (F2), header-vs-implementation truth in shell (F2), unused-variable detection in shell (F4), or test-surface enumeration completeness across diagnostic strings (F3). Three of four R1 findings were in surfaces my re-audit loop did not exercise.

**Surfaces this pattern touches.** `alpha/SKILL.md` §2.6 row 9 (post-patch re-audit) — names the row, does not name the audit's coverage requirement against the diff's actual language mix.

**Evidence.** N=1 (this cycle). Polyglot diffs (Go + shell + YAML + Markdown) need polyglot re-audits; my re-audit was monoglot.

### 4. Structured `slog.WarnContext` on degraded paths makes the runtime surface inspectable

**Pattern.** Every reinstall reason in `restore.restoreOne` emits a structured log with named attributes (`installed_version`, `lockfile_version`, `expected_version`, `error`). The same fields are exactly what `doctor.checkPackages` later surfaces in human-readable diagnostics (`stale: cnos.core (installed 1.0.0, locked 2.0.0)`). The two surfaces use the same vocabulary because they answer the same question for two audiences (machine logs vs operator diagnostics). DESIGN-CONSTRAINTS §6.3 (degraded-path visibility) is met by the same vocabulary appearing in both places, not by separate explanations.

**Surfaces this pattern touches.** `eng/go` §2.8 (structured logging via `slog`); DESIGN-CONSTRAINTS §6.3.

**Evidence.** N=1 (this cycle), but the pattern is reusable: any degraded-path policy decision should emit one structured log with the same attributes that the read-side surface (doctor / status) later renders.

### 5. Polling discipline must include a synchronous baseline pull at session start

**Pattern.** Transition-only Monitor polling has a structural blind spot: state that exists *before* the Monitor's first iteration is silently absorbed into the baseline and never surfaces as an event. β's round-1 review posted in that window; my Monitor (started after PR-open) had `prev=""` on iteration 1 and never emitted the round-1 transition because there was no transition — only a comment that already existed when the loop started. The fix is structural: pair every transition-only Monitor with an initial synchronous pull of the same surface, so the past is observed via one channel and the future via the other. β observed and named the same pattern from the dyad's other side in their close-out.

**Surfaces this pattern touches.** CDD §Tracking (polling reference scripts) — names the wake-up mechanism but does not name the initial-state pull as a precondition.

**Evidence.** N=1 (this cycle), but this is the second cycle in three (β's close-out also reports the absorption pattern from the β side) where the polling channel didn't catch a transition that existed before the loop armed.

---

Signed: α (`alpha@cnos.local`) · 2026-04-26 · merge commit `9980e3f` · release commit `9dd30d9` · cycle level per §9.1: L6 (L7 MCA shipped, capped by L6 cross-surface duplication miss).
