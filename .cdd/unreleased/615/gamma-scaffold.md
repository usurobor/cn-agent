# γ R0 scaffold — cnos#615

## Issue reference

- **Issue:** [usurobor/cnos#615](https://github.com/usurobor/cnos/issues/615) — "cds/reconciler: clear terminal lifecycle labels on closed dispatch cells"
- **Mode:** design-and-build
- **Cell kind:** `implementation` (reconciler hygiene tail)
- **Family:** extends the mechanical reconciler #593 (`cn issues fsm scan`); sibling hygiene to the #583 mechanical-dispatch wave. Related: #570 (cell-kind taxonomy), #600 / #532 (manual precedent this mechanizes).
- **Protocol:** cds
- **Wake run id:** `cds-dispatch/local-manual-20260707` (opaque; recorded, not decoded)
- **Base SHA:** `9c780a1266aec0ad0896a94823c40ac0cbdcd0cd` — verified against the working tree (`git rev-parse HEAD` matched) **and** against `.cdd/unreleased/615/CLAIM-REQUEST.yml`'s `head_commit` (`9c780a1...`, committed to main at this same SHA) and the issue's claim-record comment (`head: 9c780a1266aec0ad0896a94823c40ac0cbdcd0cd`). All three agree.
- **Cycle branch:** `cycle/615`, created from `main@9c780a12` (see Friction note 4 re: one docs-only commit that has since landed on `origin/main` ahead of this SHA).
- **run_class:** `first_pass` (per the claim comment: no prior `status:changes` history, no prior `cycle/615` branch, no prior `.cdd/unreleased/615/` artifacts beyond `CLAIM-REQUEST.yml`).

## Governing rule (restated from the issue)

**Reconciler pass, NOT a new FSM terminal state.** For a **closed** `dispatch:cell + protocol:{P}` issue still carrying `status:*` / `dispatch:cell` / `protocol:{P}`: strip those three label classes and add the `resolution/*` label keyed off GitHub's `state_reason` (`completed` → `resolution/completed`; `not_planned` → `resolution/not-planned`, a label this cell creates). `transitions.json` and the FSM engine (`table.go`/`decision.go`/`Evaluate`) are **out of scope** — this pass never evaluates a state transition; it observes GitHub's own close metadata and reconciles label taxonomy, nothing more.

## Key finding: this is structurally independent of the FSM evaluator, unlike #593

#593's recovery scanner (`scan.go`) iterates **open** issues and calls the *same* `Evaluate(table, snap)` engine `evaluate` uses — it is glue around the existing FSM. This cell's terminal-hygiene pass has no FSM state to evaluate at all: a closed issue's "next state" isn't a transitions-table lookup, it's a direct function of GitHub's own `state_reason` on the **Issue** resource. The issue's own scope line says "reuse the existing `ghRemoveLabel` machinery in `fetch.go`" — it does **not** say reuse `Evaluate`/`Table` (`scan.go`/`table.go`). That is a deliberate signal: this is a sibling reconciler in the same package, not a mode of `scan`.

**Recommendation (not a mandate — see α prompt § "Open design call"):** a new file `terminal.go` in the same package, with its own `RunTerminalSweep` engine, decoupled entirely from `activeScanStates`/`Evaluate`/`Table`. It needs only two GitHub-side facts per candidate issue: the issue's current label set, and `state_reason`. It needs zero local-git/CI observation (no branch, no PR, no run state) — this makes it simpler and faster to test than `scan.go`, not a variant of it.

## Design (recommended shape; α may diverge with documented rationale)

New sub-verb: **`cn issues fsm terminal --protocol P [--apply]`**, in `src/packages/cnos.issues/commands/issues-fsm/` (new file `terminal.go`, wired into `issuesfsm.go`'s `Run` switch alongside `evaluate`/`scan` — same `--apply`-flag-not-a-separate-verb pattern the package doc comment already states).

1. **List** — GitHub REST `GET /repos/{repo}/issues?labels=dispatch:cell,protocol:{P}&state=closed&per_page=100` (paginate if needed). This single query is the AC1 recognizer: it can only return issues that are **closed** AND still carry **both** `dispatch:cell` and `protocol:{P}` — which is exactly "still carrying the live selectors" from the issue's own wording. The response already includes each issue's full `labels[]` and `state_reason` — no second API call needed per issue (mirrors `liveListActiveIssues`'s single-query-then-client-filter shape in `scan.go`).
2. **Per issue** — from the label set already in hand: collect every `status:*` label present (there should be at most one per label-doctrine, but the code should not assume exactly one — strip whichever `status:*` labels are found), plus `dispatch:cell`, plus `protocol:{P}`.
3. **Resolve resolution label from `state_reason`:**
   - `"completed"` → `resolution/completed` (label already exists in the repo — verified live: color `ededed`, description `null`).
   - `"not_planned"` → `resolution/not-planned` (label does **not** exist yet — this cell creates it, color `ededed` / no description, to match `resolution/completed` exactly, per the operator's pinned clarifying comment).
   - anything else (in practice: this branch should be unreachable for a `state=closed` query — GitHub's `state_reason` enum is only `completed` / `not_planned` / `reopened`/null, and `reopened` cannot appear on a currently-closed issue — but code defensively for an unmapped/unknown value): **do not mutate any label on that issue**; report it as skipped/unresolved so an operator can see it, exactly like `scan.go`'s `blocked` outcome never mutates. Never guess a resolution label (mirrors the `cnos#368` "never blind-act on missing evidence" doctrine already embedded in `scan.go`).
4. **Apply (only under `--apply`, guard-gated on step 3 having resolved a concrete resolution label):**
   - Remove each `status:*` label found, via `ghRemoveLabel` (reused verbatim from `fetch.go`; its existing 404-tolerance makes each removal independently idempotent).
   - Remove `dispatch:cell` via `ghRemoveLabel`.
   - Remove `protocol:{P}` via `ghRemoveLabel`.
   - Ensure the resolved resolution label exists in the repo (new small helper, `ghEnsureLabelExists` or equivalent — `POST /repos/{repo}/labels {name, color, description}`, tolerating a 422 "already_exists" response exactly the way `ghRemoveLabel` tolerates a 404 — same idiom, new endpoint). Only `resolution/not-planned` will actually need creating; the helper should be unconditional/idempotent so it is also a harmless no-op for `resolution/completed`.
   - Add the resolution label via `ghAddLabel` (reused verbatim; GitHub's add-labels endpoint is itself idempotent).
5. **Idempotence is structural, not bolted on:** after step 4 runs once, the issue no longer carries `dispatch:cell` **and** `protocol:{P}` (both were just removed) — so step 1's own list query excludes it on every subsequent run. AC3 falls out of the list query's own filter, the same way #593's scanner's idempotence falls out of `todo` not being in `activeScanStates`. No separate "already done" check should be needed; if α finds one is needed anyway, name why in `self-coherence.md`.

## Source of truth

| Claim / surface | Canonical source | Status |
|---|---|---|
| Reconciler precedent this extends | `src/packages/cnos.issues/commands/issues-fsm/scan.go` + `.cdd/unreleased/593/gamma-scaffold.md` (or its released copy) | On `main` |
| Reusable label-mutation primitives | `src/packages/cnos.issues/commands/issues-fsm/fetch.go` — `ghRequest`, `ghAddLabel`, `ghRemoveLabel` (404-tolerant), `ghGetJSON` | On `main`; reuse verbatim, do not fork |
| CLI dispatch shape / package doc | `src/packages/cnos.issues/commands/issues-fsm/issuesfsm.go` (`Run`, `runEvaluate`, `runScan` as the pattern to mirror for `runTerminal`) | On `main` |
| FSM transition table — **must remain byte-unchanged** | `src/packages/cnos.cds/skills/cds/fsm/transitions.json` (states: `ready, todo, in-progress, review, changes` — no `done`/`closed` state exists or should be added) | On `main` |
| Generic label manifest (does **not** cover `resolution/*` today — pre-existing, out of this cell's scope to fix) | `src/packages/cnos.core/labels.json` | On `main`; only `status:*` + `dispatch:cell` declared |
| Label-doctrine (ownership split; confirms `resolution/*` is outside the doctrine's protected generic-lifecycle set) | `src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` §1 | On `main` |
| Live label facts (verified via `gh api` at scaffold time) | `resolution/completed` exists (`color: ededed`, `description: null`); `resolution/duplicate` and `resolution/superseded` also exist; `resolution/not-planned` and `resolution/wontfix` do **not** exist | GitHub live state, 2026-07-07 |
| Close-reason → resolution mapping (pinned, binding) | Issue #615's operator clarifying comment (posted 2026-07-07T01:51:47Z) — supersedes the issue body's original "`not_planned` → `resolution/wontfix`-or-equivalent" phrasing | Issue comment; binding |
| CI gates named in AC6 | `.github/workflows/build.yml` (jobs: `go`, `binary-verify` = "Binary", `package-verify` = "Package", `package-source-drift` = I1, `protocol-contract-check` = I2, `link-check` = I4, `skill-frontmatter-check` = I5, `cdd-artifact-check` = I6, `dispatch-repair-preflight`, `dispatch-closeout-integrity`); `.github/workflows/install-wake-golden.yml` | On `main`; see Friction note 3 for exact job-name mapping and the install-wake-golden path-filter caveat |
| Manual precedent this mechanizes | `gh issue view 600` / `gh issue view 532` — both closed, `stateReason: COMPLETED`, both already carry `resolution/completed` and **no** `status:*`/`dispatch:cell`/`protocol:cds` (manually cleaned) — useful as a "target post-state" reference, not a fixture source (their raw pre-cleanup label state was not preserved) | GitHub live state |

## Per-AC oracle list

| AC | Oracle (mechanical pass/fail) |
|---|---|
| **AC1** | `cn issues fsm terminal --protocol cds` (no `--apply`) against a fixture/fake-server set containing a closed `dispatch:cell + protocol:cds` issue still carrying a `status:*` label → the tool's report names that issue as a recognized candidate (non-empty "would reconcile" style note), without requiring `--apply`. Negative: the same call against a closed issue **not** carrying `dispatch:cell`/`protocol:cds` in its label set never appears in the candidate list (list-query filter itself excludes it). |
| **AC2** | With `--apply`: for the AC1 candidate, after the call, the fake-server request log shows exactly one `DELETE .../labels/status%3A{state}` (or one per `status:*` label present), one `DELETE .../labels/dispatch%3Acell`, one `DELETE .../labels/protocol%3Acds`, and one `POST .../labels` adding the resolution label whose name matches the fixture's `state_reason` per the pinned mapping (`completed`→`resolution/completed`, `not_planned`→`resolution/not-planned`). Assert the exact label name added, not just "a label was added." |
| **AC3** | Idempotence, mirroring `scan_test.go`'s `TestScan_Idempotent` shape: run the sweep twice. First pass reconciles the AC1/AC2 candidate. Second pass's `ListClosed...` call (against a fake server / injected fixture reflecting the post-cleanup label state — i.e., the issue no longer carries `dispatch:cell`+`protocol:cds`) returns zero candidates for that issue → zero new label-mutation requests on the second pass. Separately: an **open** dispatch:cell+protocol:cds issue (any status) run through the same sweep is never touched — assert zero label requests referencing that issue's number across both passes. |
| **AC4** | (a) `git diff` (or a repo-wide `grep` in CI) shows zero changes to `src/packages/cnos.cds/skills/cds/fsm/transitions.json` across the whole cycle branch. (b) Every label-removal call in the new code traces to `ghRemoveLabel` (reused from `fetch.go`) via a `grep -n "ghRemoveLabel\|http.MethodDelete" <new file>` audit — no direct/duplicate DELETE-labels HTTP call written outside that shared primitive. |
| **AC5** | Four fixtures/fake-server scenarios exist and are exercised by tests: (1) closed + `state_reason: completed` + stale `status:*`/`dispatch:cell`/`protocol:cds` → cleaned per AC2. (2) closed + `state_reason: not_planned` + stale labels → cleaned, with `resolution/not-planned` (not `resolution/completed`) added. (3) closed + already clean (no `status:*`/`dispatch:cell`/`protocol:cds` present — i.e., excluded by the list query itself, or present-but-already-resolution-tagged) → no-op, zero mutation calls. (4) an **open** dispatch:cell+protocol:cds issue → never appears as a candidate, zero mutation calls, regardless of its `status:*` value. |
| **AC6** | All of: `go test ./...` and `go vet ./...` from `src/go` (job **Go**); `go test ./...` from `src/packages/cnos.issues/commands/issues-fsm` (same module-boundary requirement `scan.go`/AC9-of-#568 already established — this new file lives in the same module); **Binary** (`binary-verify` job, Tier 1 kata); **Package** (`package-verify` job, Tier 2 kata); **I1** (`./cn build --check` — package/source drift; only relevant if this cell touches any `cn.package.json`-tracked source, which it should not); **I2** (protocol-contract-check diff — unaffected, no schema files touched); **I4** (`lychee` link check — relevant because this cell adds `.cdd/unreleased/615/*.md`; any relative links inside those files must resolve); **I5** (SKILL.md frontmatter validation — unaffected unless a SKILL.md is edited, which this cell should not need); **I6** (`cn cdd verify --unreleased --exceptions .cdd/exceptions.yml` — directly relevant, this cell's own CDD artifacts must satisfy the ledger); **install-wake-golden** (path-filtered on `install-wake`/orchestrator files + the two dispatch workflow files — this cell should not touch any of those paths, so the workflow will most likely not even trigger; see Friction note 3); **dispatch-repair-preflight** and **dispatch-closeout-integrity** (presence-of-contract guards over the dispatch prompt/protocol surfaces — unaffected, should stay green automatically since this cell never touches those files). |

## Scope (restated from issue)

**In scope:** a new terminal-hygiene reconciliation pass (new Go file(s) in `src/packages/cnos.issues/commands/issues-fsm/`, new CLI sub-verb, new tests + fixtures); creation of the `resolution/not-planned` GitHub label (color/description matching `resolution/completed`); wiring to reuse `ghRemoveLabel`/`ghAddLabel` from `fetch.go`.

**Out of scope (non-goals, restated):** no FSM terminal state, no new `status:*` label, no change to `transitions.json`, no change to active-cell (`scan.go`/`evaluate`) behavior, no close authority (this pass never calls `gh issue close` or any close-equivalent — it only acts on issues GitHub already reports as closed), no `resolution/duplicate` handling (GitHub's `state_reason` has no `duplicate` value to key off — see Friction note 2), no Demo 0, no edits to `src/packages/cnos.core/labels.json` (that manifest already excludes `resolution/*` entirely — pre-existing drift, not this cell's to fix).

## The α prompt

```
Branch: cycle/615

You are α, implementing cnos#615: "cds/reconciler: clear terminal lifecycle
labels on closed dispatch cells." Read the full issue body and the pinned
operator clarifying comment (close-reason -> resolution-label mapping table)
before starting. Read .cdd/unreleased/615/gamma-scaffold.md in full — it is
your scaffold contract for this round.

## What you are building

A terminal-hygiene reconciliation pass, sibling to the existing #593
recovery scanner (`cn issues fsm scan`), NOT a modification of it and NOT a
new FSM state. For every CLOSED `dispatch:cell + protocol:{P}` issue that
still carries `status:*` / `dispatch:cell` / `protocol:{P}`:
  - remove those labels (reuse `ghRemoveLabel` from
    src/packages/cnos.issues/commands/issues-fsm/fetch.go verbatim — do not
    fork or duplicate its DELETE-labels logic)
  - add the correct `resolution/*` label keyed off GitHub's `state_reason`
    on the Issue resource:
      state_reason "completed"    -> resolution/completed  (already exists)
      state_reason "not_planned"  -> resolution/not-planned (YOU create this
                                       label as part of this cell — color
                                       "ededed", no description, to match
                                       resolution/completed exactly; verify
                                       resolution/completed's live color via
                                       `gh api repos/{owner}/{repo}/labels/resolution%2Fcompleted`
                                       before hardcoding, in case it has
                                       drifted since scaffold time)
      anything else (should be unreachable for a closed issue given
        GitHub's state_reason enum is only completed/not_planned/reopened,
        and reopened cannot occur on a currently-closed issue) -> do NOT
        mutate any label on that issue; report it as skipped. Never guess a
        resolution label.

## Implementation location (recommended; you may diverge — see "Open design
call" below)

New file `src/packages/cnos.issues/commands/issues-fsm/terminal.go`, mirroring
`scan.go`'s shape (an injectable-dependency Options struct + a Run function +
live defaults), but with NO dependency on `Evaluate`/`Table`/transitions.json
at all — this pass never evaluates an FSM state; it only reads GitHub's
`state_reason` and the issue's current label set (both come back in a single
`GET /repos/{repo}/issues?labels=dispatch:cell,protocol:{P}&state=closed`
call — no per-issue local git/CI observation is needed, unlike scan.go's
assembleLive). Wire a new `terminal` sub-verb into `issuesfsm.go`'s `Run`
dispatch switch, mirroring `runEvaluate`/`runScan`'s flag shape (--protocol
required, --apply, --repo, --token). Add a small idempotent
"ensure this label exists with this color/description" helper in fetch.go
(POST /repos/{repo}/labels, tolerating a 422 "already_exists" response the
same way ghRemoveLabel tolerates 404) — this is the one genuinely new HTTP
primitive this cell needs; everything else (remove, add) reuses fetch.go
verbatim.

## Idempotence (AC3) falls out of the list-query filter

Once a closed issue's `dispatch:cell` and `protocol:{P}` labels are both
removed, your own list query (`labels=dispatch:cell,protocol:{P}&state=closed`)
naturally excludes it from every future run. Do not add a separate "already
processed" marker/check unless you find the list-query-filter argument
insufficient — if so, name exactly why in self-coherence.md.

## Open design call you must resolve and document

The issue phrases the CLI shape ambiguously ("a terminal/--closed sweep").
Two shapes are plausible: (a) a new `terminal` sub-verb (recommended above,
because the issue explicitly names fetch.go's ghRemoveLabel as the reuse
target, not scan.go's Evaluate/Table engine — this cell has nothing to do
with FSM evaluation); (b) a `--closed` flag bolted onto the existing `scan`
sub-verb. If you choose (b), you must NOT let it touch scanOne's Evaluate()
call path for the closed-issue branch — a closed issue has no meaningful FSM
"current state" to evaluate against transitions.json (AC4 forbids any
transitions.json coupling). Pick one, implement it, and write one paragraph
in self-coherence.md naming which you chose and why — mirroring the
documented-design-call precedent in cnos#574 AC4 (see fetch.go's
observeRemoteBranch doc comment for that precedent's shape).

## Fixtures (AC5)

Provide four test scenarios under testdata/ (JSON fixtures or a fake-HTTP-
server-driven test, whichever matches your chosen implementation's testing
seam — mirror scan_test.go's fakeAssembleFacts/withFakeGitHub patterns):
  1. closed, state_reason=completed, stale status:*/dispatch:cell/protocol:cds
     present -> cleaned, resolution/completed added.
  2. closed, state_reason=not_planned, stale labels present -> cleaned,
     resolution/not-planned added (NOT resolution/completed, NOT
     resolution/wontfix).
  3. closed, already clean (list query naturally excludes it, or it is
     already resolution-tagged) -> zero mutation calls.
  4. open dispatch:cell+protocol:cds issue (any status) -> never a
     candidate, zero mutation calls.

## Guardrails (binding — violating any of these is a scope violation, not a
judgment call)

- Zero-diff on src/packages/cnos.cds/skills/cds/fsm/transitions.json.
- Zero changes to scan.go's Evaluate-call path, activeScanStates, or any
  existing evaluate/scan test's behavior — run the full existing
  issues-fsm test suite and confirm no existing test's expected output
  changed.
- Do not add gh issue close (or any close-equivalent) anywhere in this
  code path — this pass only acts on issues GitHub already reports closed.
- Do not edit src/packages/cnos.core/labels.json — that manifest does not
  cover resolution/* today (not even resolution/completed is listed there);
  that is pre-existing drift, out of this cell's scope to fix.
- No resolution/duplicate handling — GitHub's state_reason has no
  "duplicate" value; do not invent a special case for it.
- Write self-coherence.md documenting: your CLI-shape decision (open design
  call above), the exact resolution/not-planned label color/description you
  used and where you verified resolution/completed's live values from, and
  confirmation that every AC1-AC6 oracle in the scaffold passes.

## When done

Commit, push to cycle/615, append your review-readiness signal to
self-coherence.md, and stop — you do not dispatch β yourself.
```

## The β prompt

```
Branch: cycle/615

You are β, independently reviewing α's implementation of cnos#615 on
cycle/615. Read .cdd/unreleased/615/gamma-scaffold.md (this scaffold) and
.cdd/unreleased/615/self-coherence.md (α's round record) in full before
forming any verdict. Do not take α's self-coherence claims as verified until
you have independently walked each oracle below yourself.

## Independent AC walk (do not skip any — re-derive each verdict yourself)

- AC1: find the new code's "list closed candidates" path. Confirm it queries
  (or fixture-equivalently filters on) closed + dispatch:cell + protocol:{P}
  + at least one of status:*/dispatch:cell/protocol:{P} present. Confirm a
  closed issue lacking dispatch:cell or protocol:{P} is excluded.
- AC2: confirm the removal set is exactly {every status:* label found,
  dispatch:cell, protocol:{P}} and the added label is exactly the one the
  pinned mapping requires for the fixture's state_reason (completed ->
  resolution/completed; not_planned -> resolution/not-planned). Reject if
  the not_planned case maps to resolution/completed or resolution/wontfix
  instead of resolution/not-planned — that is the exact ambiguity the
  operator's clarifying comment resolved, and getting it backwards is a
  correctness bug, not a style nit.
- AC3: run (or read) the idempotence test yourself. Confirm a second pass
  against post-cleanup state produces zero new label-mutation calls, AND
  confirm an open dispatch:cell+protocol:{P} issue is never touched across
  either pass — check this explicitly, do not assume it from AC1 alone.
- AC4: `git diff main...cycle/615 -- src/packages/cnos.cds/skills/cds/fsm/transitions.json`
  must be empty. Grep the new code for any direct HTTP DELETE against the
  labels endpoint that bypasses ghRemoveLabel — there should be none.
- AC5: confirm all four fixture scenarios exist and are exercised by a
  passing test each (completed-stale-cleaned, not_planned-stale-cleaned,
  already-clean-noop, open-untouched).
- AC6: confirm CI is green on Go, Package, Binary, I1, I2, I4, I5, I6,
  install-wake-golden (or confirm it did not need to trigger — path-filtered
  on install-wake/orchestrator files this cell should not touch; if it DID
  trigger and is red, that is a real finding), dispatch-repair-preflight,
  dispatch-closeout-integrity. Also independently run the full existing
  issues-fsm test suite (`go test ./... -v` from
  src/packages/cnos.issues/commands/issues-fsm) and confirm no pre-existing
  test's expected behavior changed (a regression here would itself violate
  the "no active-lifecycle behavior change" non-goal even if every new test
  passes).

## Guardrail verification (independent of the AC table)

- Confirm src/packages/cnos.core/labels.json was not edited.
- Confirm no gh issue close / close-equivalent call was introduced.
- Confirm resolution/not-planned's color/description actually matches
  resolution/completed's live values (not just what the scaffold assumed at
  R0 scaffold time — re-verify live via `gh api
  repos/{owner}/{repo}/labels/resolution%2Fcompleted` yourself in case it
  drifted, and check the created label matches).
- Confirm the "open design call" (CLI shape: new sub-verb vs --closed flag
  on scan) is documented in self-coherence.md with a stated rationale, and
  that if the --closed-flag path was chosen, it demonstrably does not route
  a closed issue through Evaluate()/transitions.json.

## Verdict

Write .cdd/unreleased/615/beta-review.md with a full R[N] section: outcome
per AC, any findings (with severity), and `verdict: converge` or
`verdict: iterate`. Commit, push, and stop — you do not re-dispatch α
yourself.
```

## Scope guardrails (binding on α)

- Zero-diff on `src/packages/cnos.cds/skills/cds/fsm/transitions.json` — this is the single most load-bearing non-goal in the issue; the FSM owns the pre-close lifecycle only.
- Zero behavior change to `scan.go`/`evaluate`'s existing active-lifecycle reconciliation — the full pre-existing `issues-fsm` test suite must pass unchanged.
- No close authority anywhere in the new code path (no `gh issue close`, no state-transition equivalent) — this pass only reacts to issues GitHub already reports closed.
- No edits to `src/packages/cnos.core/labels.json` (pre-existing gap: it does not track `resolution/*` at all today; out of scope to fix here).
- No `resolution/duplicate` special-casing — not reachable via `state_reason`.
- All label mutation must route through `ghRemoveLabel`/`ghAddLabel` (reused) plus exactly one new idempotent "ensure label exists" primitive for label creation — no parallel/duplicate HTTP mutation path.

## Friction notes

1. **CLI shape is genuinely ambiguous in the issue text** ("a terminal/`--closed` sweep"). I recommend a new `terminal` sub-verb (parallel to `evaluate`/`scan`) over a `--closed` flag on `scan`, because the issue explicitly names `fetch.go`'s `ghRemoveLabel` as the reuse target and never mentions `scan.go`'s `Evaluate`/`Table` engine — this pass has no FSM state to evaluate, and folding it into `scanOne`'s switch would either (a) force a fake "state" through `Evaluate()` against `transitions.json` (violating AC4) or (b) require an early-return branch inside `scanOne` that never touches `Evaluate` for the closed case, which is more awkward than a clean sibling function. I left this as an explicit open call for α with a documented-rationale requirement rather than mandating it, per γ's "do not smuggle constraints outside the issue" discipline — the issue itself left it open with "e.g."

2. **GitHub's `state_reason` has no native `duplicate` value.** The operator's clarifying-comment table lists a third row (`duplicate` → optional `resolution/duplicate`, "not present — out of scope unless specified"), but GitHub's REST API `state_reason` field on an Issue is only ever `completed`, `not_planned`, `reopened`, or `null` — there is no `duplicate` enum value to key off mechanically. An issue closed "as a duplicate" via GitHub's UI reports `state_reason: not_planned` (GitHub's own UI folds "duplicate" under its "not planned" close reason). Consequence: the third table row is not independently reachable by this pass's `state_reason`-keyed logic at all — any duplicate-closed issue will receive `resolution/not-planned` under this cell's rule, which is already correct per the `not_planned` mapping and requires no extra branch. AC5 does not require a dedicated "duplicate" fixture for this reason. Flagging this so α doesn't spend time building unreachable dead code and so β doesn't file a false-negative finding expecting a `duplicate` branch that cannot exist.

3. **AC6 gate names, resolved to concrete jobs/scripts** (verified in `.github/workflows/build.yml` and `.github/workflows/install-wake-golden.yml` at scaffold time): I1 = `package-source-drift` job (`./cn build --check`); I2 = `protocol-contract-check` job (diff `docs/reference/schemas/protocol-contract.json` vs `tests/fixtures/protocol-contract.json`); I4 = `link-check` job (`lychee`); I5 = `skill-frontmatter-check` job; I6 = `cdd-artifact-check` job (`cn cdd verify --unreleased`); Go = the `go` job; Package = `package-verify`; Binary = `binary-verify`; `dispatch-repair-preflight`/`dispatch-closeout-integrity` = their like-named jobs (presence-of-contract guards, unaffected by this cell). **`install-wake-golden` is path-filtered** (`.github/workflows/install-wake-golden.yml`'s `on.push.paths`/`on.pull_request.paths` only cover `install-wake`/orchestrator files + the two dispatch workflow YAMLs) — since this cell should touch none of those paths, the workflow will most likely not trigger at all for this PR. "Green" for that gate in practice means "did not need to run"; if it unexpectedly *does* run and fails, that is a real signal something touched a path it shouldn't have. The issue's AC6 does not name "I3" — that gap is pre-existing (possibly folded per the `cnos#600` consolidation referenced in the two `check-dispatch-*.sh` scripts) and is not this cell's concern.

4. **`origin/main` has advanced one commit beyond the pinned base SHA.** At scaffold time, `git fetch origin main && git rev-parse origin/main` returned `58dbca7608c08eebcaf9fb63443d71f5e16726c7` — one commit (`58dbca7`, "board-map: regenerate docs/development/board from live open issues") ahead of the pinned `9c780a1266aec0ad0896a94823c40ac0cbdcd0cd`. Confirmed `9c780a12` **is** an ancestor of `origin/main` (linear, no divergence) and the one intervening commit is a docs-only automated board regeneration unrelated to `issues-fsm`/labels — no rebase was performed before branching, since the pinned SHA matches the claim record exactly and the drift is non-conflicting. If α's diff touches `docs/development/board/**`, treat that as an unrelated pre-existing file, not part of this cell's scope.

5. **`resolution/completed`'s exact live color/description was verified once at scaffold time** (`color: ededed`, `description: null`, via `gh api repos/usurobor/cnos/labels/resolution%2Fcompleted`) but is not re-verified by this scaffold at implementation time — the α prompt above instructs α to re-verify live immediately before creating `resolution/not-planned`, in case it drifts between scaffold and implementation.
