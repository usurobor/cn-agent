cell_kind: implementation

# γ scaffold — cnos#493

label-doctor: install + verify cnos.core canonical labels in repos (cnos#491 live-discovery follow-up)

Mode: `design-and-build` (per issue header). Branch: `cycle/493`. Base SHA: `31d7ddfa53efb67f2d996b085f4d92986b585ef9` (verified against `origin/main` HEAD at scaffold time — see "SHA verification" below).

Wave context (per operator's follow-up comment on the issue): sequenced after #612 (main.go's `--version`/`--help` handling — already landed, visible in `src/go/cmd/cn/main.go`); #613 (PAT-free FSM tier) depends on this landing; runs in parallel with #618 (different files — not touched by this scaffold).

## SHA verification

`git rev-parse origin/main` (and local `main`) = `31d7ddfa53efb67f2d996b085f4d92986b585ef9`, exact match to the pinned SHA in δ's contract. `cycle/493` contains exactly one commit beyond that base (`d3105cfa`, the CLAIM-REQUEST.yml). No mismatch.

## Peer enumeration (gamma/SKILL.md §2.2a)

- `grep -rln "label-doctor" src/` → hits only in `src/go/internal/cell/cell.go` and `src/go/internal/cell/cell_test.go` (a *consumer* of the not-yet-built label-doctor, citing cnos#493 as the reason a `status:changes` preflight can fail — see "Existing coupling" below) and `src/go/internal/repoinstall/repoinstall.go` (the `ensureCanonicalDispatchLabels()` stub — see below). No implementation exists anywhere.
- `grep -rln "label-doctor\|labels sync\|LabelDoctor\|labels-sync\|cn labels"` across `src/` and `.github/` → no other hits. The gap is real; this is additive, not a duplicate of in-flight work.
- `find src/packages/cnos.core/commands -maxdepth 1 -type d` → `save`, `daily`, `weekly`, `install-wake`. No `label-doctor` directory exists yet.
- `go.work` → lists exactly 4 modules: `./src/go`, `cdd-verify`, `issues-map`, `issues-fsm`. No cnos.core Go module exists yet — α is adding the first one.

## Live audit (real findings, not the issue's narrower claim)

Ran `gh label list --repo usurobor/cnos --limit 200 --json name,color,description` against `src/packages/cnos.core/labels.json` (schema `cn.labels.v1`, 8 entries, `owner: cnos.core`). The issue's Gap section names only `status:review`'s color drift as the confirmed finding ("other canonical labels MAY ALSO be missing or have drifted"). The live diff below resolves that "MAY" — **7 of 8 canonical labels are drifted, 1 is missing outright**:

| Canonical label | Live state | Drift |
|---|---|---|
| `status:backlog` | present, color `ededed` (matches canonical) | description empty vs canonical `"Well-formed scope but not yet refined to dispatch readiness."` |
| `status:ready` | present, color `ededed` | color drift (canonical `c2e0c6`) + description empty vs canonical |
| `status:todo` | present, color `ededed` | color drift (canonical `0e8a16`) + description empty vs canonical |
| `status:in-progress` | present, color `ededed` | color drift (canonical `fbca04`) + description empty vs canonical |
| `status:review` | present, color `ededed` | color drift (canonical `5319e7` — this is the issue's originally-cited finding) + description empty vs canonical |
| `status:changes` | present, color `ededed` | color drift (canonical `d93f0b`) + description empty vs canonical |
| `status:blocked` | **absent** | missing entirely |
| `dispatch:cell` | present, color `ededed` | color drift (canonical `1d76db`) + description empty vs canonical |

This is the real AC1/AC2 audit surface α's tool must detect and AC3/AC2's install command must repair. α's implementation must not special-case `status:review` — it must generically diff every entry in `labels.json` against the live label set (name/color/description), so this scaffold's table is a snapshot α's own audit artifact reproduces at implementation time, not a hand-maintained list α copies in.

## Existing coupling — why this gap is not merely cosmetic

Two pieces of shipped runtime code already assume the label-doctor mechanism exists and name cnos#493 explicitly when it doesn't:

1. `src/go/internal/repoinstall/repoinstall.go` — `runDispatchCds()` (the `cn repo install --dispatch cds` path, cnos#610) calls `ensureCanonicalDispatchLabels()` (line ~444), a **stub that unconditionally returns an error** naming cnos#493: `"canonical dispatch labels not ensured: cnos#493 label-install mechanism is not yet available; labels must be applied manually until it ships"`. This is the primary integration point α's new mechanism must wire into — replacing the stub's body, not adding a parallel path.
2. `src/go/internal/cell/cell.go` — `Returner.preflightTargetLabel()` (the `status:review → status:changes` review-return transition, cnos#500) already has a hard preflight guard that refuses the transition if `status:changes` is missing from the repo, with the error message literally instructing **"Run label-doctor before retrying"**. This is independent runtime evidence the doctrine expects a command literally named `label-doctor` to exist.

Existing tests that pin the current (stub) behavior and that α's change will legitimately need to update (not a backward-compat violation — replacing the stub *is* the point of this issue):
- `src/go/internal/repoinstall/repoinstall_test.go` — asserts `Run(...)` returns an error containing `"cnos#493"` when `--dispatch cds` is used (search `cnos#493` in that file).
- `src/go/internal/cli/cmd_repo_install_test.go` — same assertion at the CLI layer (search `cnos#493`).

## Source-of-truth table

| Claim / surface | Canonical source | Status |
|---|---|---|
| Canonical label manifest | `src/packages/cnos.core/labels.json` (schema `cn.labels.v1`) | Shipped on main; 8 entries (7 `status:*` lifecycle + 1 `dispatch:cell`) |
| Label doctrine (ownership + install responsibility) | `src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` §3 ("Ownership split" — `cn install cnos.core` creates the generic set) + §5 ("Manifest" — manifest is data-only, a consumer reads it and creates labels) | Shipped on main |
| `cn repo install --dispatch cds` (existing consumer with a stub gap) | `src/go/internal/repoinstall/repoinstall.go` `ensureCanonicalDispatchLabels()` (~line 444) + `runDispatchCds()` (~line 375) | Shipped stub; α replaces the stub body |
| CLI wiring precedent (thin-wrapper pattern) | `src/go/internal/cli/cmd_repo_install.go` (flags → `repoinstall.Options` → `repoinstall.Run`), `src/go/internal/cli/cmd_issues_fsm.go` (thinnest possible: `Run` delegates entirely to the domain package) | Shipped on main |
| Kernel-command registration | `src/go/cmd/cn/main.go` `reg.Register(&cli.IssuesFsmCmd{})` / `&cli.RepoInstallCmd{}` block | Shipped on main |
| Go-module co-location precedent (domain logic in `src/packages/{owner}/commands/{name}/`, registered as a kernel command, own `go.mod`, added to `go.work`) | `src/packages/cnos.issues/commands/issues-fsm/` (`go.mod`: `module github.com/usurobor/cnos/packages/cnos.issues/commands/issues-fsm`), `src/packages/cnos.issues/commands/issues-map/` | Shipped on main |
| GitHub REST label primitives precedent (net/http, no `gh` shellout, dependency-free) | `src/packages/cnos.issues/commands/issues-fsm/fetch.go`: `ghGetJSON` (GET), `ghRequest` (generic authed request), `ghAddLabel`/`ghRemoveLabel` (issue-label mutation), **`ghEnsureLabelExists`** (POST `/repos/{repo}/labels`, tolerates 422 "already exists" — create-only, no color/description update-on-drift) | Shipped on main (cnos#615). α's new package needs an analogous create primitive **plus a new update/PATCH primitive** (`PATCH /repos/{owner}/{repo}/labels/{name}`) this precedent does not yet have — drift repair (color/description) is not the same operation as create. |
| Competing runtime dependency style (`gh` CLI shellout, NOT the REST convention above) | `src/go/internal/cell/cell.go` `Returner.preflightTargetLabel` / `applyLabelTransition` (uses `gh label list` / `gh issue edit` shellouts) | Shipped on main — a *different* subsystem, not the pattern to imitate for this issue (see Implementation contract, Runtime dependencies row) |
| Go module workspace registration | `go.work` (repo root) — `use (./src/go ./src/packages/cnos.cdd/commands/cdd-verify ./src/packages/cnos.issues/commands/issues-map ./src/packages/cnos.issues/commands/issues-fsm)` | Shipped on main; α adds a 5th `use` entry |
| CI module-test-step precedent (separate go.work modules are NOT covered by `go test ./...` from `src/go`) | `.github/workflows/build.yml` step `"Test issues-fsm module (cnos#568 AC9)"` (`working-directory: src/packages/cnos.issues/commands/issues-fsm`, `run: go test ./... -v`) | Shipped on main; α adds an analogous step for the new module |
| cli/-boundary compliance gate (domain logic must not live in `cli/cmd_*.go`) | `.github/workflows/build.yml` step `"Dispatch boundary check (INVARIANTS.md T-002)"` — greps `internal/cli/cmd_*.go` for banned imports (`os`, `net/http`, `encoding/json`, `archive/`, `compress/`, `crypto/`, `path/filepath`) and fails CI if any thin wrapper imports them directly | Shipped on main; the new `cmd_label_doctor.go` MUST pass this unchanged — no new exemption to add |
| CLI ergonomics smoke test (must list every subcommand) | `.github/workflows/build.yml` step `"CLI ergonomics smoke test (cnos#612 F1-F4)"` — hardcoded `for c in "doctor" "status" ... "repo install"` loop asserting `./cn $c --help` exits 0 | Shipped on main; α must add the new command's name to this loop or the smoke test will not catch a broken `--help` on it (it will simply not test it — silent gap, not a red CI, so must be added deliberately) |
| Live repo label state (this scaffold's audit) | `gh label list --repo usurobor/cnos --limit 200 --json name,color,description` (run at scaffold time, 2026-07-08) | See "Live audit" table above |
| Live-discovery origin | cnos#491 smoke comment [`4784509561`](https://github.com/usurobor/cnos/issues/491#issuecomment-4784509561) + `docs/gamma/smoke/cds-dispatch-smoke-20260623.md §7` | Shipped on cycle/491 / PR #492 |

## Per-AC oracle (mechanical, checkable)

1. **AC1 — Audit complete.** Oracle: an artifact (issue comment, PR doc, or test fixture) exists listing all 8 `labels.json` entries and their live state (present-and-matching / present-but-drifted / missing) — the table above is the reference shape; α's own tool output (or a committed audit doc) is the actual artifact this AC gates. Cross-reference against the table above must match (mechanical: re-run the same `gh label list` + diff and get the same classification, modulo whatever this PR itself repairs).
2. **AC2 — All canonical labels exist with canonical color+description.** Oracle: `gh label list --repo usurobor/cnos --json name,color,description` (or the GitHub API directly) shows all 8 `status:*`/`dispatch:cell` entries present, each with `color` and `description` byte-equal (case-insensitive on hex) to `labels.json`.
3. **AC3 — `status:review` corrected.** Oracle: same API call, `status:review.color == "5319e7"` (case-insensitive). Subsumed by AC2 but named separately in the issue because it's the originally-observed instance — do not treat AC2 passing as implicitly closing AC3 without checking this row explicitly.
4. **AC4 — Command exists, idempotent.** Oracle: (a) the command's source directory exists under `src/packages/cnos.core/commands/` (mechanical: `test -d`); (b) running it against a repo with drifted/missing labels produces the canonical set (functional test, e.g. via httptest fixture mirroring `issues-fsm`'s test style — no live network call required for CI-safe testing); (c) running it a second time in a row with no external state change performs zero mutating API calls / reports "no drift" (idempotence — assert on request count against the fake server, not just "exit 0 both times").
5. **AC5 — CI guard.** Oracle: a CI job (in `.github/workflows/`, likely `build.yml` or a new dedicated workflow mirroring `install-wake-golden.yml`'s shape) runs the command in dry-run/diff mode; (a) fails when a fixture/injected drift is present; (b) passes on the clean/repaired state. Both directions must be exercised, not just the pass case — a guard that can only pass is not a guard.

## α prompt (implementer)

```
Role: α (implementer), cnos#493, cycle/493.

Branch: cycle/493 (already exists on origin, based on main@31d7ddfa53efb67f2d996b085f4d92986b585ef9).
Switch to it — do not create a new branch.

Read cnos#493's full issue body (`gh issue view 493 --repo usurobor/cnos
--json body,comments --jq '.body, (.comments[].body)'`) and this scaffold
(.cdd/unreleased/493/gamma-scaffold.md) in full before writing any code.
The scaffold's Source-of-truth table, Live audit table, and Existing
coupling section are grounded in real repo state read at scaffold time —
treat them as facts to build from, not claims to re-verify from scratch
(though you should re-run the `gh label list` audit yourself once you
begin, since live label state may have moved).

Task: (1) build a `cn label-doctor` command that audits the cnos repo's
GitHub labels against src/packages/cnos.core/labels.json and, applied,
creates missing labels and repairs drifted color/description on existing
ones, idempotently; (2) wire it into cn repo install --dispatch cds's
existing ensureCanonicalDispatchLabels() stub in
src/go/internal/repoinstall/repoinstall.go so that path stops
unconditionally failing; (3) add a CI guard that fails on injected label
drift and passes on clean state.

Follow the 7-axis Implementation contract below exactly — it is not
optional guidance, it is the pinned technical contract for this cell.
Where the contract names an open design point (target-repo resolution
from git remote), resolve it yourself and document the resolution in
self-coherence.md; do not leave it unresolved.

## Implementation contract

| Axis | Value |
|---|---|
| **Language** | Go 1.24, matching every other kernel/package command in this repo (`src/go` module + the `go.work`-listed package-command modules). No shell script (rule out the `install-wake`-style shell-script pattern — that renderer is templating-only; this is REST+diff logic, which the repo's own precedent (`issues-fsm`) already does in Go, not shell). |
| **CLI integration target** | New top-level kernel subcommand `cn label-doctor` (single-token, matching `cn doctor`/`cn build`/`cn status` naming, and matching the issue's own primary suggested name plus `cell.go`'s existing runtime error text "Run label-doctor before retrying"). Registered in `src/go/cmd/cn/main.go` (`reg.Register(&cli.LabelDoctorCmd{})`, in the same block as `IssuesFsmCmd`/`RepoInstallCmd`). Thin wrapper at `src/go/internal/cli/cmd_label_doctor.go` — mirror `cmd_issues_fsm.go` exactly: `Run` delegates entirely to the domain package, `Spec()` returns `{Name: "label-doctor", Source: SourceKernel, Tier: TierKernel, NeedsHub: false}` (hub-less: this must work in a freshly-cloned repo with no `.cn/` yet, same reasoning as `RepoInstallCmd`). The file MUST NOT import `os`/`net/http`/`encoding/json`/`path/filepath` etc. directly — the CI "Dispatch boundary check" step greps for exactly this and will fail red if violated. Flags: `--repo OWNER/NAME` (optional; live path defaults from git remote, see Runtime dependencies row), `--token` (optional; falls back to `$GITHUB_TOKEN` then `$GH_TOKEN`, mirroring `issues-fsm`'s existing flag), `--dry-run` (report-only, no mutation; required for AC5's CI guard). |
| **Package scoping** | New Go module + directory `src/packages/cnos.core/commands/label-doctor/` (domain logic, tests, own `go.mod`: `module github.com/usurobor/cnos/packages/cnos.core/commands/label-doctor`, `go 1.24` — copy `issues-fsm/go.mod`'s shape exactly). This is the **first cnos.core Go module** — α must add a 5th entry to the root `go.work` `use (...)` block or the new module is invisible to `go build`/`go test` from the workspace root and to any importer. α must also add a new CI step to `.github/workflows/build.yml` mirroring the existing `"Test issues-fsm module (cnos#568 AC9)"` step (own `working-directory`, own `go test ./... -v` — NOT covered by the top-level `src/go` `go test ./...` step, exactly as that step's own comment explains for issues-fsm/issues-map). Do not put domain logic in `src/go/internal/` — the co-location-with-cnos.issues precedent (`issues-fsm`, `issues-map`) is the one to follow for a *package-owned*, kernel-registered command; `src/go/internal/repoinstall` is the wrong precedent to copy here (that's core-kernel-owned infrastructure, not a cnos.core-package-owned label doctrine implementation). |
| **Existing-binary disposition** | `src/go/internal/repoinstall/repoinstall.go`'s `ensureCanonicalDispatchLabels()` (~line 444) stub is **replaced**, not left alongside a new parallel path. Because both `internal/repoinstall` (part of the `./src/go` module) and the new `src/packages/cnos.core/commands/label-doctor` module are `go.work`-linked, wire this as an **in-process Go function call** (`import labeldoctor "github.com/usurobor/cnos/packages/cnos.core/commands/label-doctor"`), exactly the way `cmd_issues_fsm.go` already imports `issuesfsm` cross-module — do NOT exec a subprocess/vendored-binary the way `runDispatchCds` execs `cn-install-wake` (that pattern exists because the renderer does file templating and must run as a separate artifact from a vendored package; label-doctor is pure API+diff logic with no such constraint, and an in-process call has zero extra runtime dependency). `internal/cell/cell.go`'s `preflightTargetLabel`/`applyLabelTransition` (the `gh` CLI shellout-based label-transition path) is a **separate subsystem — do not touch it and do not refactor it onto the new REST primitives** as part of this issue (out of scope per the issue body; its error message already correctly names "label-doctor" as the fix, and once this issue ships, an operator can literally run `cn label-doctor` to satisfy that guard — no code change needed there for that to become true). |
| **Runtime dependencies** | None beyond Go stdlib (`net/http`, `encoding/json`, `os`) + the existing `$GITHUB_TOKEN`/`$GH_TOKEN` env-var convention `issues-fsm` already uses. No `gh` CLI shellout (that is `cell.go`'s pattern, not this issue's — the issue's own runtime consumer, `repoinstall.go`, and its sibling `issues-fsm`, are both net/http-only; follow those, not `cell.go`). New primitives needed beyond what `issues-fsm/fetch.go` already has: a **list-repo-labels** GET (`GET /repos/{owner}/{repo}/labels`, paginated) and an **update-drifted-label** PATCH (`PATCH /repos/{owner}/{repo}/labels/{name}` with new `color`/`description`) — `ghEnsureLabelExists`'s existing POST-with-422-tolerance only covers *create*, not *repair-drift-on-existing*. Open design point α must resolve and document in `self-coherence.md`: `repoinstall.Options` has no target-repo-slug field usable for GitHub label API calls today (`Options.Repo` defaults to `usurobor/cnos` and is used for *cnos release resolution*, i.e. where to download packages FROM — not the *installing target repo*'s "owner/repo" for label API calls, which for an arbitrary tenant repo is a different value entirely). No existing helper in `src/go/internal` resolves "owner/repo" from a git remote (confirmed absent by `grep -rn "git remote\|ParseGitHubRemote\|RepoFromRemote"` across `src/go/internal` — no hits). α must add one (parse `git remote get-url origin`, handle both `https://github.com/{owner}/{repo}.git` and `git@github.com:{owner}/{repo}.git` forms) — this is a small, scoped utility, not a design blocker. |
| **JSON/wire contract** | N/A for GitHub's label API (external, already-fixed schema — standard REST label object: `name`/`color`/`description`). For `labels.json` itself: schema `cn.labels.v1` is already shipped and stable (§5 of the doctrine skill says the manifest is forward-compatible / implementers MAY add fields) — the new tool reads it as-is; no schema change is in scope. |
| **Backward compat** | `cn repo install` with `--dispatch none` (the default) is fully unaffected — `ensureCanonicalDispatchLabels()` is only reachable from the `--dispatch cds` branch. `cn repo install --dispatch cds` **does change observable behavior** (this is the point of the issue, not a violation): today it always fails with the cnos#493 stub error; after this issue, it must actually attempt to ensure the canonical labels and only fail if that genuinely cannot succeed (e.g. no token / API error) — update `repoinstall_test.go` and `cmd_repo_install_test.go`'s existing `cnos#493`-string assertions accordingly (this is expected, in-scope test churn, not scope creep). The `.cn/deps.json`/`.cn/deps.lock.json` determinism invariants and the "no agent-hub scaffold written" invariant (both documented in `repoInstallHelp` in `cmd_repo_install.go`) are unrelated to this change and must remain byte-identical. |

Non-goals (do not implement, per the issue's "Out of scope" and this
scaffold's Scope guardrails section below): cnos.cds/cnos.cdr/cnos.cdw
protocol-qualifier label installs; the dispatch-protocol
dispatch_label_missing doctrine update; bundling into the agent-admin
wake; touching internal/cell/cell.go's gh-shellout label-transition path.

Write self-coherence.md per alpha/SKILL.md, using the canonical bare
section-header form (## Gap, ## Skills, ## ACs, ## CDD Trace, ## Self-
check, ## Debt — not a decorated "## §Gap" form; cn cdd verify's
sectionPresent() does a literal/prefix match). Signal review-readiness
when all 5 ACs are met and self-coherence.md is committed on cycle/493.
```

## β prompt (reviewer)

```
Role: β (reviewer), cnos#493, cycle/493.

Branch: cycle/493. Review α's implementation of the cnos.core canonical-label
audit + `cn label-doctor` install/repair command + CI guard, against the 5
ACs in cnos#493's issue body plus the γ scaffold at
.cdd/unreleased/493/gamma-scaffold.md.

Walk the AC oracle independently — do not trust α's self-coherence.md's
claims without re-deriving each check yourself:

AC1 (audit): does a committed artifact (doc, PR description, or test
fixture) enumerate all 8 labels.json entries with their live-repo state
(matching/drifted/missing)? Does it match reality — re-run
`gh label list --repo usurobor/cnos --limit 200 --json name,color,description`
yourself and diff against labels.json; does α's artifact's classification
agree, accounting for whatever this PR itself repairs?

AC2/AC3 (repair): after α's PR (if it applies live, or against α's test
fixtures if it does not touch the live repo directly), do all 8 canonical
labels have canonical color+description? Is status:review specifically at
color 5319e7? Do not accept "AC2 passed" as closing AC3 without checking
that row explicitly — the issue names it as a separate AC on purpose.

AC4 (command): does `src/packages/cnos.core/commands/label-doctor/` exist
as a real Go module (own go.mod, added to root go.work's `use` block)? Is
it registered as a kernel command in src/go/cmd/cn/main.go and wired via a
THIN wrapper at src/go/internal/cli/cmd_label_doctor.go (Run delegates to
the domain package — check it does NOT import os/net/http/encoding/json
directly, since the "Dispatch boundary check" CI step greps for exactly
that and a false-negative there would mean the step's own regex missed a
real violation)? Run the tool's tests directly — are they fixture-based
(httptest, no live network dependency, mirroring issues-fsm's test style)?
Is idempotence actually asserted (a second dry-state run makes ZERO
mutating requests — check the test asserts request counts, not just
"exit 0 twice")? Does src/go/internal/repoinstall/repoinstall.go's
ensureCanonicalDispatchLabels() stub get REPLACED (not left as a parallel
dead path) with a real in-process call into the new package? Were
repoinstall_test.go / cmd_repo_install_test.go's existing "cnos#493"
string-assertion tests updated to match the new (non-stub) behavior rather
than just deleted or weakened?

AC5 (CI guard): is there a CI job/step that runs the tool in --dry-run (or
equivalent diff) mode? Does it actually FAIL on injected drift (find the
test/fixture that proves this — a guard that only ever passes is not
verified) and PASS on clean state? Was build.yml's "Test issues-fsm module"
step pattern followed for the new module (a step whose working-directory
is the new module, since go test ./... from src/go does NOT cover it)?
Was the new "label-doctor" command name added to build.yml's CLI
ergonomics smoke-test loop (the hardcoded `for c in ...` list) so
`cn label-doctor --help` is actually exercised?

Scope check: confirm α did NOT touch cnos.cds/cnos.cdr/cnos.cdw protocol-
qualifier label installs, did NOT rewrite dispatch-protocol's fallback
behavior (that's explicitly a separate future doctrine cycle per the
issue's "Out of scope but worth filing as separate follow-ups"), and did
NOT touch internal/cell/cell.go's gh-shellout-based label-transition path
(separate subsystem, explicitly out of scope per the γ scaffold's
Existing-binary disposition row).

Report verdict per the standard β review contract (approve / request
changes with named findings, each tagged to the AC or scope-guardrail it
violates).
```

## Scope guardrails (per issue "Constraints / non-goals")

**In scope (from the issue body verbatim, do not narrow further without cause):**
- Audit cnos repo's label set vs `labels.json`; identify missing/drifted labels.
- Materialize missing labels; correct color/description of drifted ones (matching `labels.json` exactly).
- Design + ship the install/repair command (`cn install cnos.core` or `cn label-doctor` per the issue's own "or" — this scaffold pins `cn label-doctor`, see Implementation contract).
- CI guard that fails on injected drift and passes on clean state.

**Out of scope (explicit — do not implement, even if it seems adjacent):**
- Updating `cnos.cds`/`cnos.cdr`/`cnos.cdw` protocol-qualifier label installs (`protocol:cds` etc.) — per-package concern, file separately if a similar gap exists there.
- Rewriting the dispatch-protocol skill (cnos#454) to handle missing labels more gracefully at wake-runtime (the `dispatch_label_missing` failure-class doctrine update) — named explicitly by the issue as a separate future doctrine cycle.
- Bundling label-doctor into the agent-admin wake — stays an operator/CI-invoked command for now.
- Refactoring `internal/cell/cell.go`'s `gh`-shellout label-transition path onto the new REST primitives — separate subsystem, not named in any AC.

## Friction notes (for δ / α awareness)

1. **Command name is a γ-pinned design call, not settled by the issue.** The issue text offers `cn install cnos.core` OR `cn label-doctor` as equally acceptable; the operator's follow-up comment adds `cn labels sync` / `cn install cnos.core --labels-only` as further candidates and says "operator-approved equivalent" — i.e. explicitly delegates the final name choice. This scaffold pins `cn label-doctor` (matching `cell.go`'s existing runtime error text and the issue's primary name) as the concrete AC4 target so α does not have to re-litigate naming; if δ or the operator has a strong preference for a different name, that is a legitimate pre-dispatch correction to make now rather than after α builds against this scaffold.
2. **Target-repo resolution for the live GitHub label API is genuinely undetermined today.** `repoinstall.Options.Repo` is the *source* of cnos releases, not the *installing target*'s owner/repo — see Implementation contract → Runtime dependencies row. α needs a small git-remote-parsing utility that does not exist yet anywhere in `src/go/internal` (confirmed absent by grep). This is scoped and mechanical, not a blocker, but is exactly the kind of gap that inflates α's context if not named up front — naming it here per gamma/SKILL.md §3.3 ("prompt cleverness is not a substitute for issue quality... fix the issue instead" — here, the issue itself doesn't need editing since this is normal implementation-detail discovery, but the scaffold should not leave α to rediscover the absence from scratch).
3. **`ghEnsureLabelExists` (cnos.issues, cnos#615) is create-only, not drift-repair.** It is the closest existing precedent but cannot be imported (separate Go module owned by a different package, unexported besides) and does not cover the PATCH-on-drift case this issue actually needs (AC2 requires fixing existing labels' color/description, not just creating missing ones). α should write a structurally similar (not imported) pair of primitives in the new module — mirroring the pattern, not duplicating cross-module in violation of package boundaries.
4. **The live audit found materially worse drift than the issue's Gap section describes.** The issue's Gap section names only `status:review`'s color as a *confirmed* finding, hedging the rest with "MAY ALSO be missing or have drifted." This scaffold's live audit (table above) resolves that hedge: 7 of 8 canonical labels are drifted and 1 (`status:blocked`) is missing outright. α's audit output should reflect the real severity, not just re-state the issue's narrower framing — this is not scope creep (all 8 labels are already named as in-scope by "Materialize missing labels in the cnos repo... Correct color of the existing status:review label" plus the general "identify all missing or drifted... labels" framing), just a correction of the issue's own hedged uncertainty using data this scaffold already gathered.
5. **`cell.go`'s existing preflight guard is a live empirical dependency, not just documentation.** `Returner.preflightTargetLabel` in `internal/cell/cell.go` will refuse a `status:review → status:changes` transition today if `status:changes` happens to be absent (it is currently present but drifted, not absent, so this specific guard is not firing right now — but the missing `status:blocked` label means any code path that similarly preflights on `status:blocked`'s existence, if one exists or is added later, is live-broken today). This scaffold does not require α to audit every such call site (out of scope), but β should be aware this is not a hypothetical risk being fixed defensively — cnos#500's cycle/500 R1 already hit exactly this failure mode once.
