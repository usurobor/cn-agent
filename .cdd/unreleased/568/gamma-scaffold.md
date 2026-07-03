# γ scaffold — cnos#568 (cds/issues: FSM Phase 1 — read-only issue-state reconciler)

**Cycle:** cycle/568 · base `main@cd61f3d26951253d539bd7111f37cf6af4de5a35` (verified against working tree at scaffold time) · protocol `cds` · run_class `first_pass`.

## Source of truth

| Fact | Value |
|---|---|
| Issue body | #568 (see issue for Mode/Gap/Impact/Constraints/ACs verbatim — this scaffold does not restate ACs, it operationalizes them) |
| Parent | #567 (wave master — "Workers produce matter. The FSM owns status labels.") |
| CLI precedent (co-location + kernel registration pattern) | `src/packages/cnos.issues/commands/issues-map/` (`issuesmap.go`, `fetch.go`, own `go.mod`), `src/go/internal/cli/cmd_issues_map.go` (thin dispatch shim), `src/go/cmd/cn/main.go:47` (`reg.Register(&cli.IssuesMapCmd{})`), `src/packages/cnos.issues/SKILL.md` §"Command-dispatch disposition" (temporary built-in shim until #216 — mirror exactly) |
| CLI resolution mechanics | `src/go/internal/cli/dispatch.go` `ResolveCommand` — noun-verb: `args[0]+"-"+args[1]` registry key. For a 3-word surface (`cn issues fsm evaluate`), register `Name: "issues-fsm"` (noun=`issues`, verb=`fsm`); `Remaining` after resolution is `["evaluate", "--issue", "568", ...]`; the command's own `Run` parses `Remaining[0]` as the sub-verb (`evaluate` is the only sub-verb in Phase 1 — do not stub `apply`, per issue's explicit out-of-scope). `src/go/internal/cli/registry.go` for `Register`/`Lookup`. |
| Dispatch-boundary rule | `src/go/internal/cli/cmd_*.go` files MUST NOT import `os`, `net/http`, `encoding/json`, `path/filepath`, etc. (enforced by build.yml "Dispatch boundary check (INVARIANTS.md T-002)", line 38-54 of `.github/workflows/build.yml`) — `cmd_issues_fsm.go` must be import+one-line-delegation only, exactly like `cmd_issues_map.go`. |
| No-external-dependency convention | `commands/issues-map/fetch.go` header comment: "REST is used over GraphQL for a dependency-free stdlib implementation" — `commands/issues-map/go.mod` has zero `require` lines. Mirror this: stdlib only (`encoding/json`, `net/http`, `os/exec` for `git`), no new third-party deps. |
| CI wiring precedent | `.github/workflows/build.yml` — the `go` job's `Test` step runs `go test ./... -v` with `working-directory: src/go`; **verified empirically this does NOT cover `commands/issues-map`** (`go list ./...` from `src/go` returns 15 packages, none under `issues-map` — Go workspace `./...` is CWD-relative, not workspace-wide). AC9's CI guard must be an **explicit new step** that `cd`s into the new module and runs `go test ./...` — do not assume the existing `go` job step covers it. |
| Package boundary (AC1) | `cnos.issues` owns the **generic** evaluator/fact-model/CLI (`src/packages/cnos.issues/commands/issues-fsm/`, new Go module, own `go.mod`, wired into `src/go/go.mod` via `require`+`replace` and `go.work` `use`, mirroring `commands/issues-map` exactly). `cnos.cds` owns the **CDS transition table** as package-owned **declarative data**, not Go literals (AC1 negative case: "transition logic exists only as prompt prose or inline Go literals" is a FAIL) — place at `src/packages/cnos.cds/skills/cds/fsm/transitions.json` (JSON: stdlib `encoding/json`, no YAML dependency, consistent with the dependency-free convention above). The generic evaluator loads a transition table from a `--table` path (defaulting to a repo-root-relative resolution of the cnos.cds path) and is unaware of "cds" as a concept beyond reading that file — this is what keeps the engine package-generic per AC1's ownership split. |
| CI gates that must stay green (AC10) | I1 (package-source-drift), I2 (protocol-contract-check), I4 (link-check), I5 (skill-frontmatter-check), I6 (cdd-artifact-check), install-wake-golden, dispatch-repair-preflight (cnos#516 guard), dispatch-closeout-integrity (cnos#524 guard), Go (`go` job), Package (`package-verify`), Binary (`binary-verify`) — all named jobs in `.github/workflows/build.yml` + `.github/workflows/install-wake-golden.yml` (confirm exact file name before citing in PR body). |

## Fact-snapshot model (per issue body's "Fact model" section)

Define a `FactSnapshot` struct (generic, `cnos.issues`-owned) carrying exactly the fields the issue enumerates: issue number + labels; workflow run state (queued/in_progress/completed+conclusion); `cycle/{N}` branch existence; branch-diff-vs-base emptiness; PR existence; PR commits-beyond-base count; `.cdd/unreleased/{N}/` artifact presence (by filename); `REVIEW-REQUEST.yml` presence; repair-contract / prior β/δ findings presence; CI/checks state where available. Two assembly paths: (a) **live** — GitHub REST (`net/http`, stdlib, GITHUB_TOKEN from env, mirroring `fetch.go`) + local `git` (`os/exec`) run from within a cnos checkout; (b) **fixture** — a JSON file fully specifying the snapshot, for hermetic fixture tests (AC3–AC7) and for the `--fixture` CLI flag (mirrors `issues-map`'s `--fixture`/`--stdin` offline path).

## AC → oracle → surface (operationalized; do not renumber, do not drop any)

| AC | Oracle (what α proves) | Surface |
|---|---|---|
| AC1 | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` exists, package-owned, declarative (states/triggers/guards/actions as JSON, not Go); the evaluator's transition logic reads from it (grep-provable: no switch/if-chain hardcoding CDS-specific state names inside the generic engine package) | package data + grep |
| AC2 | `cn issues fsm evaluate --issue {N}` (and `--fixture <path>` for offline) prints: current state, observed facts, enabled transition, blocked reason, proposed action — structured stdout | CLI + fixture test |
| AC3 | Fixture: `status:review` + no PR + no commits + no `REVIEW-REQUEST.yml` → evaluator reports blocked/invalid + missing-evidence list | test |
| AC4 | Fixture: `status:in-progress` + no active run + no/empty branch + no PR → proposes `status:todo` + requeue reason | test |
| AC5 | Fixture: `status:in-progress` + no active run + branch has commits beyond base + no PR → proposes δ-recovery (open/request recovery PR) or `blocked` if a PR cannot be created — **never** proposes `status:todo` (the #368 failure) | test |
| AC6 | Fixture: `status:changes` + missing repair contract/findings → blocks `todo` proposal. Fixture: repair contract present → enables repair-dispatch proposal marked `repair_pass` | test (positive + negative) |
| AC7 | Run evaluator twice on one fixture → identical decision; second run proposes no additional action | test |
| AC8 | No `--apply` path exists anywhere in the command; CDS dispatch prompt/protocol files (`.github/workflows/cnos-cds-dispatch.yml`, `cnos.cds/orchestrators/cds-dispatch/SKILL.md`) are byte-identical to base; grep confirms zero label-write call (`gh issue edit`, `gh api ... -X PATCH .../labels`, github REST label endpoints) in the new command's source tree | diff + grep, asserted in a test or CI step |
| AC9 | New CI step (new job or a step appended to an existing job in `build.yml`) runs `go test ./...` inside the new `issues-fsm` module; deleting a fixture-covered decision path turns that step red (demonstrate by construction, not by actually breaking it) | CI |
| AC10 | All named gates (table above) still pass after the diff — verify locally: `cd src/go && go build ./... && go vet ./... && go test ./... -v`; `cn help`/smoke; do not touch unrelated files | CI (verified pre-PR) |

## Implementation contract (δ-pinned; α MUST NOT improvise these axes)

1. **Language:** Go 1.24, matching `go.work`.
2. **CLI integration target:** compiled-in kernel command (`SourceKernel`/`TierKernel`), registered in `src/go/cmd/cn/main.go` as `reg.Register(&cli.IssuesFsmCmd{})`, dispatched by a **thin** `src/go/internal/cli/cmd_issues_fsm.go` (import + delegation only, no domain logic — the T-002 dispatch-boundary check will fail CI otherwise).
3. **Package scoping:** generic engine + CLI + fact model + fixtures under `src/packages/cnos.issues/commands/issues-fsm/` (new Go module `github.com/usurobor/cnos/packages/cnos.issues/commands/issues-fsm`, wired into `src/go/go.mod` via `require`+`replace` and into root `go.work` `use (...)`, exactly mirroring the three `commands/issues-map` wiring points). CDS transition table as JSON at `src/packages/cnos.cds/skills/cds/fsm/transitions.json`.
4. **Existing-binary disposition:** extends the existing `cn` binary; no new binary, no new package-command exec-dispatch entry in `cn.package.json` (matches the `issues-map` "no commands entry — kernel built-in shadows it" disposition; note the #216 shim caveat in a code comment, mirroring `cmd_issues_map.go`'s comment).
5. **Runtime dependencies:** stdlib only (`encoding/json`, `net/http`, `os/exec`, `context`, `flag`/manual arg parsing matching repo convention). No new `go.sum` entries.
6. **JSON/wire contract:** new surface, no backward-compat constraint from a prior version. Human-readable stdout is the AC2-required contract; a `--json` flag emitting the same decision as structured JSON is optional (nice-to-have, not an AC) — do not let it expand scope if time-constrained; the human-readable block is the binding oracle.
7. **Backward-compat invariant:** zero changes to `commands/issues-map/**`, zero changes to `.github/workflows/cnos-cds-dispatch.yml`, zero changes to any CDS dispatch prompt/protocol skill file. Zero label-mutation code path (AC8).

## Scope guardrails

- **Zero label mutation.** No `--apply`. No gh label-write call anywhere in the new tree. This is the Phase-1/Phase-2 boundary (#569 is Phase 2, out of scope).
- **Zero CDS dispatch prompt/protocol change.** Do not touch `cnos.cds/orchestrators/cds-dispatch/SKILL.md`, `dispatch-protocol/SKILL.md`, or the rendered workflow.
- **No new status labels / taxonomy change.**
- **No Demo 0 / wave-autonomy / external-service work.**

## α prompt (dispatched separately with this scaffold attached)

Implement cnos#568 Phase 1 on `cycle/568` per the AC table + implementation contract above. Write `.cdd/unreleased/568/self-coherence.md §R0` verifying each AC by number with concrete evidence (file:line, test name, command output). Commit and push to `cycle/568`. Do not open a PR yet — δ opens the PR after β converges.

## β prompt (dispatched separately after α's R0)

Independently walk the AC1–AC10 oracle table above against α's diff on `cycle/568`. Do not trust α's self-coherence claims — re-derive each: read the transition-table file directly, run the CLI yourself against the fixtures, grep for label-write calls and for prompt/protocol-file diffs, run `go build`/`go vet`/`go test` yourself. Write `.cdd/unreleased/568/beta-review.md §R0` with a per-AC verdict and an overall `verdict: converge` or `verdict: iterate` (with findings) at the top.

## Friction notes

- Existing CI does not test `commands/issues-map` today (empirically confirmed via `go list ./...`) — this is a pre-existing gap in a different package's CI wiring, not in scope for #568 to fix; α should not "fix" it as a drive-by, only ensure the *new* module's own CI step is correctly wired per AC9.
- The exact filename of `install-wake-golden.yml` should be confirmed by α/β from `.github/workflows/` rather than assumed from this scaffold's paraphrase of the issue's gate list.
