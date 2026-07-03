# γ scaffold — cycle/569

## Issue reference / mode

- Issue: [#569](https://github.com/usurobor/cnos/issues/569) — "cds/issues: FSM Phase 2 — authority flip (FSM applies labels; workers request transitions)"
- Parent/wave: #567 (master)
- Mode: **design-and-build** (per issue header). MCA preconditions are NOT met — there is no committed `DESIGN.md`/`PLAN.md` pinning the exact mechanism for "in-progress → review" proposal or the label-write shape; α does real design work inside this cycle for those two points (see §Friction notes).
- Precondition verified: #568 (Phase 1 read-only reconciler) is **CLOSED/merged** — landed at `src/packages/cnos.issues/commands/issues-fsm/`. #570 (cell-kind taxonomy) is **CLOSED/merged** — landed at `src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md`.
- **Binding scope constraint (operator comment, 2026-07-03T21:34:35Z on #569):** Phase 2 may enforce FSM-controlled transitions **only for the implementation-cell lifecycle** (the lifecycle `transitions.json` already encodes: `ready/todo/in-progress/review/changes`). Phase 2 MUST NOT introduce enforcement behavior for `issue_authoring`, `wave`, `cleanup`, `audit`, `release`, or any other `cell_kind` from `CELL-KINDS.md` until #570's taxonomy is explicitly consumed for those kinds in a future cycle. This is a hard constraint, not a suggestion — see §Scope guardrails.
- Base SHA discrepancy (recorded per γ pre-flight): the wake-invoked-δ input named **`208c07f45d130603ec953cbaea40b2cf2cd08e75`** as current main SHA. At scaffold time, `git rev-parse origin/main` returned **`0520235e1285c078eb3bc9d7eeba191b0413c53b`** — one commit ahead (`0520235`, a `board-map` auto-regeneration of `docs/development/board/{board-data.json,index.html}`, doc-data-only, does not touch any CDD/CDS/handoff/dispatch-protocol skill file). Per the wake-invoked input contract (`delta/SKILL.md` §9.2 input #3: "δ verifies this SHA against the working tree before dispatching γ"), γ branched from the **actual current `origin/main` HEAD** (`0520235e`), not the stale pinned SHA. No staleness-propagation action is required (the intervening commit is board-map data, not a role/lifecycle skill per `gamma/SKILL.md` §2.5 "Spec-staleness propagation" → "When not to propagate").

## Cycle branch

`cycle/569`, created from `origin/main@0520235e1285c078eb3bc9d7eeba191b0413c53b`, pushed to `origin/cycle/569`.

## Surfaces α is expected to touch

Go (evaluator + CLI):

- `src/packages/cnos.issues/commands/issues-fsm/issuesfsm.go` — add an `--apply` flag to `evaluate` (per AC1's oracle string `cn issues fsm evaluate --issue {N} --apply`; this is a flag on the existing `evaluate` sub-verb, not a new `apply` sub-verb — the issue's own oracle line pins this). Wire the guard-gated mutation path; update the package doc comment (currently states "Phase 1 is read-only... Label-write authority is Phase 2 (cnos#569)" — that forward-reference is now consumed here) and the `Usage:`/`fs.Usage` text (currently claims "There is no --apply flag").
- `src/packages/cnos.issues/commands/issues-fsm/decision.go` — extend `Decision`/`Render` (or add a sibling) to report whether a mutation was applied, still read-only-safe when `--apply` is absent or `outcome != "proposed"`.
- `src/packages/cnos.issues/commands/issues-fsm/fetch.go` — add the label-write primitive. **No reusable GitHub label-write helper exists anywhere in the repo** (verified: `rg "issues/%d/labels"` and `rg "RemoveLabel|AddLabel"` under `src/go` and `src/packages` return no hits) — α authors this fresh, following `fetch.go`'s existing dependency-free `net/http` + `ghGetJSON` idiom (no third-party GitHub client; same auth-header pattern: `Authorization: Bearer <token>`, `Accept: application/vnd.github+json`, `X-GitHub-Api-Version: 2022-11-28`).
- `src/packages/cnos.issues/commands/issues-fsm/issuesfsm_test.go` — extend for AC1 (guard-gated apply + idempotence) and AC3 (blocked transition mutates nothing, exits nonzero).
- `src/packages/cnos.issues/commands/issues-fsm/testdata/*.json` — new fixtures as needed (existing fixtures: `ready.json`, `review-empty.json`, `review-with-pr.json`, `changes-with-repair.json`, `changes-no-repair.json`, `in-progress-active.json`, `in-progress-dead-with-commits.json`, `in-progress-dead-no-matter.json` — `review-empty.json` is already the AC3 negative-case fixture from Phase 1's read-only reconciliation rule; reuse/extend it, do not duplicate).
- `src/go/internal/cli/cmd_issues_fsm.go` — update the doc comment (currently says "there is no 'apply' sub-verb — label-mutation authority is Phase 2 (cnos#569), out of scope here"); no structural change expected since args pass through unchanged, but verify.

CDS-owned declarative data (per `table.go`'s own doc: "adding a new CDS state or rule requires editing only this file, not Go source" — this is the table's designed extension point, not a re-derivation of Phase-1 evaluator logic):

- `src/packages/cnos.cds/skills/cds/fsm/transitions.json` — Phase 1 shipped **no rule proposing `in-progress → review`** (the `in-progress` state's rules only propose `todo` or `delta_recovery`; `review` only appears as a *reconciliation* target for an already-labeled state). AC2 requires workers to *request* the review transition instead of setting the label directly — that request has to resolve to something the evaluator can act on. α adds a rule to the `in-progress` state proposing `target_state: "review"` gated on the guards the table already declares (`review_request_present`, `pr_exists`/`pr_has_commits`, `branch_has_commits` — no new guard functions needed; reuse `table.go`'s existing `guardFuncs` registry). This is the one real design decision in this cycle — see §Friction notes.

CDS wake/protocol prose (AC2's named surface):

- `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` — step 7 of "Invoke δ in wake-invoked mode" currently reads: *"only on β's converge verdict transitions the cell's label `status:in-progress → status:review`"* — this is δ setting the label directly. Update so δ's action is: ensure `REVIEW-REQUEST.yml` + matter exist, then **request** the transition via `cn issues fsm evaluate --issue {N} --apply` rather than writing the label directly. The "Lifecycle transitions (this wake's authority)" table and the "Closeout integrity preflight" section reference the same event and need consistent wording.
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.6 "Return tokens" — the `status:review` row currently says δ "opens (or updates) a cycle-PR ... then writes the label transition + the PR-URL comment." Update the mechanism description to route through `cn issues fsm evaluate --apply` per AC2, without changing δ's authority (δ still decides *when* to request the transition; the FSM decides *whether* the guards allow it).
- `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` §2.4 "Lifecycle transitions" and §2.9 "Closeout integrity" — check for the same direct-label-write phrasing; align if present.

CI (AC3/AC5):

- `.github/workflows/build.yml` — `dispatch-closeout-integrity` job (script `scripts/ci/check-dispatch-closeout-integrity.sh`) already encodes the empty-review detector as a shell predicate independent of the Go evaluator. α decides (issue explicitly allows either surface: *"a CI/guard (or the apply path itself)"*) whether AC3 is satisfied by (a) the new `--apply` guard alone (structural: the evaluator refuses to propose/apply `in-progress → review` without evidence), (b) reusing/extending `check-dispatch-closeout-integrity.sh`, or (c) both. Do not weaken or duplicate the existing shell detector without reason.

## Per-AC oracle list

### AC1 — `--apply` is guard-gated

- **Invariant:** `--apply` mutates labels only when the evaluated transition's guards pass.
- **α proves via:** `cn issues fsm evaluate --issue {N} --apply` (or `--fixture` for hermetic tests) — positive case: a fixture with a passing-guard `proposed` outcome (e.g. `changes-with-repair.json`, or a new `in-progress`-with-review-evidence fixture) mutates the label and prints confirmation; running `--apply` again against the **post-mutation state** (a second fixture representing the new state, or — in live mode — a re-fetched snapshot) is a no-op (`outcome: valid`/`action: none`, no write attempted). Negative case: `review-empty.json`-class fixture (or the new blocked in-progress case) exits nonzero and the test asserts zero calls to the label-write path (inject a fake/counting label-writer in the test, or assert via `httptest` mock server that no write request was received).
- **β verifies via:** run both fixture classes through `--apply --fixture <path>` (or read α's unit tests directly, since live GitHub mutation isn't independently re-executable in review) and confirm: (1) the positive case's test asserts a write occurred with the correct label diff; (2) the negative case's test asserts *no* write call happened + nonzero exit; (3) the idempotence case is a real second call, not just an assertion about intent.

### AC2 — workers no longer own status labels

- **Invariant:** workers (concretely: δ in wake-invoked mode) produce matter and request transitions; the FSM applies status labels.
- **α proves via:** diff of `cnos.cds/orchestrators/cds-dispatch/SKILL.md` + `cnos.cdd/skills/cdd/delta/SKILL.md` §9.6 showing the `status:in-progress → status:review` mechanism now reads as "ensure REVIEW-REQUEST.yml + matter exist, then request via `cn issues fsm evaluate --issue {N} --apply`" instead of a direct label write. No prose anywhere in the touched files should say δ/the wake "writes the label transition" as its own act without going through the FSM apply path.
- **β verifies via:** `git diff origin/main...cycle/569 -- src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md src/packages/cnos.cdd/skills/cdd/delta/SKILL.md src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md`; grep the touched files post-diff for `status:review` transition prose and confirm no surviving "wake/δ sets/writes status:review directly" phrasing remains uncorrected.

### AC3 — empty-review is structurally impossible

- **Invariant:** `status:review` cannot be reached without deliverable proof (PR, commits-beyond-base, or `REVIEW-REQUEST.yml`).
- **α proves via:** a fixture reproducing the empty-review attempt (issue already names `review-empty.json` as the Phase-1 analogue; α adds/reuses the equivalent for the new `in-progress → review` apply path — an `in-progress` fixture with no PR, no commits, no `REVIEW-REQUEST.yml`) is blocked by `--apply` (nonzero exit, no label write). If α also extends `check-dispatch-closeout-integrity.sh`, its `--self-test` continues to pass.
- **β verifies via:** run the new negative fixture through `--apply` directly; confirm nonzero exit and no label mutation. If the CI script was touched, run `./scripts/ci/check-dispatch-closeout-integrity.sh --self-test` locally.

### AC4 — no broad label redesign

- **Invariant:** existing labels/taxonomy remain stable — no new `status:*`/`protocol:*`/`kind:*` labels, no renamed labels.
- **α proves via:** `git diff origin/main...cycle/569` contains no edits to a labels-definition surface (`label-doctrine/SKILL.md` §"Label set administration", any label-seeding script) unless explicitly justified in the diff's own commit message with a taxonomy-update rationale.
- **β verifies via:** grep the diff for `gh label create`, `label-doctrine`, or new `status:`/`protocol:` string literals outside test fixtures; confirm none are additions to the canonical label set.

### AC5 — current gates remain green

- **Invariant:** the authority flip does not regress repo health.
- **α proves via:** `go build ./... && go test ./...` (or the repo's standard `cn` build/test invocation) green locally before requesting review; the named CI jobs pass on the cycle branch: `I1` = `package-source-drift`, `I2` = `protocol-contract-check`, `I4` = `link-check`, `I5` = `skill-frontmatter-check`, `I6` = `cdd-artifact-check`, plus `install-wake-golden`, `dispatch-repair-preflight`, `dispatch-closeout-integrity`, `go`, `package-verify`, `binary-verify` (the issue's "Go, Package, Binary" job names) — all in `.github/workflows/build.yml` / `.github/workflows/install-wake-golden.yml`.
- **β verifies via:** `gh pr checks <PR>` (or `gh run list --branch cycle/569`) confirms all named jobs green on the cycle branch's latest commit before APPROVE.

## Source-of-truth table

| Claim / surface | Canonical source | Status |
|---|---|---|
| Phase 1 evaluator engine + CLI | `src/packages/cnos.issues/commands/issues-fsm/{issuesfsm,table,decision,fetch,snapshot}.go` | Shipped (read-only) |
| CDS transition table (data) | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` | Shipped; missing `in-progress → review` proposal rule (this cycle's gap) |
| CLI dispatch wrapper | `src/go/internal/cli/cmd_issues_fsm.go` | Shipped |
| Cell-kind taxonomy (deferral boundary) | `src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md` | Shipped (#570); NOT to be consumed for enforcement this cycle per operator scope note |
| CDS dispatch wake (worker whose behavior AC2 changes) | `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` | Shipped; step 7 + "Lifecycle transitions" + "Closeout integrity preflight" sections are AC2's target |
| δ wake-invoked mode / return tokens | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.5–§9.6 | Shipped; §9.6 `status:review` row is AC2's target |
| Dispatch protocol (label lifecycle authority) | `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` §2.4, §2.9 | Shipped; canonical transition-map source, cross-check for AC2 |
| Empty-review CI detector | `.github/workflows/build.yml` job `dispatch-closeout-integrity` + `scripts/ci/check-dispatch-closeout-integrity.sh` | Shipped; AC3 may extend or may rely on the FSM guard alone |
| Issue-pack contract | `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` | Shipped |
| Dispatch wire-format / 7-axis contract | `src/packages/cnos.handoff/skills/handoff/dispatch/SKILL.md` §2.3 | Shipped; canonical home for the implementation-contract block below |

## Scope guardrails (hard constraints)

1. **cell-kind deferral (binding, operator comment 2026-07-03T21:34:35Z).** This cycle enforces FSM-controlled transitions **only** for the implementation-cell lifecycle already encoded in `transitions.json` (`ready/todo/in-progress/review/changes`). Do **NOT** add enforcement logic keyed on `cell_kind` (e.g. branching behavior for `issue_authoring`, `wave`, `cleanup`, `audit`, `release`). `FactSnapshot.CellKind` stays exactly as Phase 1 left it — **observed, never consumed by a transition rule** (Phase 1's own `TestSeam_CellKindNotEnforced` locks this; do not break it). If the new `in-progress → review` rule is tempted to branch on `cell_kind`, stop — that is out of scope; gate only on the existing evidence guards (`review_request_present`, `pr_exists`, `pr_has_commits`, `branch_has_commits`).
2. **No broad label redesign** (issue AC4 + non-goals) — no new status labels, no taxonomy changes.
3. **No Demo 0, no wave autonomy, no external controller service** (issue non-goals) — this stays a repo-local reconciler invoked by command/workflow.
4. **No wake-source-model change, no package-command-discovery change** (#216 is untouched).
5. **Do not re-derive Phase-1 transition *logic*** (the `Evaluate`/`ruleMatches`/`guardFuncs` engine in `table.go`) — extend it only via the declarative `transitions.json` data, exactly as the engine was designed to be extended (per `table.go`'s own doc comment).
6. **STOP-and-ask trigger (issue "Rules"):** if safe label-apply requires broader GitHub API authority than a `contents.write`/`issues.write` token already used elsewhere in this repo's wakes, STOP and surface to operator before broadening scope — do not silently request new permissions.

## Friction notes

- The issue's AC1 oracle line (`cn issues fsm evaluate --issue {N} --apply`) pins `--apply` as a **flag on `evaluate`**, not a new `apply` sub-verb — this resolves what would otherwise be an ambiguity against `issuesfsm.go`'s current doc comment ("Do not add 'apply' here... mutation subcommands like 'apply' are Phase 2"). α should read that comment as the forward-reference now being fulfilled, not as a constraint against adding a flag.
- The one real design gap Phase 1 left for Phase 2: no `in-progress → review` proposal rule exists in `transitions.json`. This is necessary for AC2 to be mechanically true (a worker "requesting" a transition has to resolve to *some* evaluator rule) but the issue does not spell out the exact guard combination. γ's recommendation (non-binding, α may refine): gate on `any_true: [review_request_present, pr_has_commits]` mirroring the existing `review`-state reconciliation rule's evidence set, `all_false: [run_active]` mirroring the existing `in-progress` rules' pattern of only proposing once the run is dead. α owns the final guard combination; document the choice in `self-coherence.md`.
- No existing GitHub label-write helper exists anywhere in the repo (verified by grep, see §Surfaces) — this is new code, not a refactor of existing code. Idempotency across two live `--apply` calls falls out naturally from re-fetching live facts each call (post-mutation state no longer matches the `proposed` rule), but the remove-old-label / add-new-label pair should tolerate a 404 on removal (already-removed) without treating it as a hard error, to keep a manually-recovered or partially-applied state idempotent too. This is an implementation-contract-adjacent decision α should note explicitly in `self-coherence.md` rather than leave implicit.

## Implementation contract (pinned by δ; α MUST NOT improvise)

*Note: `cnos.handoff/skills/handoff/dispatch/SKILL.md` (the canonical wire-format home cited by `gamma/SKILL.md` §2.5) IS present and vendored in this repo at `src/packages/cnos.handoff/skills/handoff/dispatch/SKILL.md` — loaded and consulted for this block's shape. No absence to report.*

| Axis | Pinned value |
|---|---|
| Language | Go (repo-native; matches Phase 1's `issues-fsm` package and every sibling command under `src/packages/cnos.issues/commands/`) |
| CLI integration target | `cn` subcommand — extends the existing compiled-in kernel command `cn issues fsm evaluate` with a new `--apply` flag; no new subcommand, no standalone binary |
| Package scoping | `src/packages/cnos.issues/commands/issues-fsm/` for Go source (mirrors Phase 1's placement exactly); `src/packages/cnos.cds/skills/cds/fsm/transitions.json` for the CDS-owned declarative rule addition; `src/go/internal/cli/cmd_issues_fsm.go` for the dispatch-wrapper doc-comment update only (per `INVARIANTS.md` T-002 / `eng/go` §2.18 dispatch-boundary rule — no domain logic added here) |
| Existing-binary disposition | Preserve — `cn issues fsm` remains the single compiled-in kernel command; Phase 1's `evaluate` verb is extended in place, not replaced or forked |
| Runtime dependencies | None new — reuse `net/http` (stdlib) for the GitHub REST label-write calls, mirroring `fetch.go`'s existing `ghGetJSON` dependency-free pattern. No third-party GitHub client library. |
| JSON/wire contract preservation | Preserve as-is for `FactSnapshot`/fixture JSON shape (Phase 1's fixture files must remain loadable unchanged); `Decision`'s printed/rendered output may gain new fields (e.g. an "applied: true/false" line) additively, but existing field names/order for read-only `evaluate` (no `--apply`) output must not change, since Phase-1-era tooling/operators may already parse it |
| Backward-compat invariant | `cn issues fsm evaluate --issue N` (no `--apply`) remains exactly as read-only and side-effect-free as Phase 1 shipped it — zero label writes, zero behavior change, byte-identical decision logic for any fixture Phase 1's test suite already covers. `--apply` is strictly additive: it is off by default and only mutates on a passing-guard `proposed` outcome. |

## α prompt

```text
You are α. Project: cnos (usurobor/cnos).
Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.
Issue: gh issue view 569 --json title,body,state,comments
Branch: cycle/569
Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md, src/packages/cnos.eng/skills/eng/SKILL.md (coding bundle), src/packages/cnos.eng/skills/eng/go/SKILL.md (or repo-conventional Go bundle path), src/packages/cnos.cdd/skills/cdd/issue/SKILL.md, src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md, src/packages/cnos.cdd/skills/cdd/issue/constraints/SKILL.md

Read .cdd/unreleased/569/gamma-scaffold.md on this branch FIRST — it names the
exact surfaces to touch, the per-AC oracle list, the source-of-truth table,
and the binding scope guardrails (in particular: the cell-kind deferral
constraint from the operator's 2026-07-03T21:34:35Z comment on #569 — do NOT
add cell_kind-based enforcement branching this cycle).

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Go |
| CLI integration target | `cn` subcommand (extend `cn issues fsm evaluate` with `--apply`; no new subcommand) |
| Package scoping | `src/packages/cnos.issues/commands/issues-fsm/` (Go); `src/packages/cnos.cds/skills/cds/fsm/transitions.json` (CDS data); `src/go/internal/cli/cmd_issues_fsm.go` (doc-comment only) |
| Existing-binary disposition | Preserve |
| Runtime dependencies | None new (stdlib net/http only, mirroring fetch.go's ghGetJSON pattern) |
| JSON/wire contract preservation | Preserve as-is for FactSnapshot/fixture shape and for evaluate-without-apply output; additive-only for new Decision fields |
| Backward-compat invariant | `evaluate` without `--apply` stays exactly read-only/side-effect-free as Phase 1 shipped it |
```

## β prompt

```text
You are β. Project: cnos (usurobor/cnos).
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 569 --json title,body,state,comments
Branch: cycle/569

Read .cdd/unreleased/569/gamma-scaffold.md on this branch FIRST — walk the
per-AC oracle list (AC1-AC5) independently against α's diff and
self-coherence.md. In particular verify:
  - AC1's idempotence claim is exercised by a REAL second --apply call (or
    fixture pair), not merely asserted;
  - AC2's diff removes direct-label-write prose from cds-dispatch/SKILL.md
    and delta/SKILL.md §9.6 and replaces it with a request-via-FSM-apply
    mechanism;
  - AC3's negative fixture is blocked with nonzero exit and zero label writes;
  - AC4's diff touches no label-taxonomy-definition surface without an
    explicit justification in the commit message;
  - AC5's named CI jobs (I1/I2/I4/I5/I6, install-wake-golden,
    dispatch-repair-preflight, dispatch-closeout-integrity, go,
    package-verify, binary-verify) are green on the cycle branch's latest
    commit;
  - the scope guardrail is honored: no cell_kind-based enforcement branching
    was added (FactSnapshot.CellKind stays observed-only; Phase 1's
    TestSeam_CellKindNotEnforced-class invariant must still hold);
  - the 7-axis implementation contract in gamma-scaffold.md is satisfied by
    the diff (Rule 7).
```
