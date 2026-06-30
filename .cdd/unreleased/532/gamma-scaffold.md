# γ scaffold — cycle/532

Issue: [#532](https://github.com/usurobor/cnos/issues/532) — "cdd/cds: require review-request proof before status:review"

## Mode

`design-and-build` (per issue header). MCA preconditions do not hold — no stable `DESIGN.md`/`PLAN.md` exists; the issue's own "Cycle scope sizing" table records "No stable design doc yet → Not MCA." α both designs the `REVIEW-REQUEST.yml` shape (already drafted in the issue as a strong proposal, not yet a converged schema doc) and builds the guard in this cycle.

## Peer enumeration (gamma/SKILL.md §2.2a)

Performed before writing §Gap below:

- `find scripts/ci` → two existing guards: `check-dispatch-repair-preflight.sh` (cnos#516, presence-of-contract pattern: greps doc/prompt/golden/workflow for required substrings) and `validate-skill-frontmatter.sh` (CUE-schema validation against skill frontmatter).
- `rg -i "review-request|REVIEW-REQUEST"` repo-wide → 4 hits, all unrelated: `ops/inbox/SKILL.md` and `OPERATIONS.md` use `*-review-request` as an example inbox-thread-id string; `SEMANTICS-NOTES.md` uses "review-request" as a session-type label in a formal-semantics aside; `HANDSHAKE.md` uses `alice/review-request` as a naming-convention example. **None of these define or imply an existing `REVIEW-REQUEST.yml` artifact or a deliverable-proof guard.** The issue's claim that this artifact is "Missing" is confirmed, not asserted.
- `find schemas/cds schemas/cdd` → both packages already have a `fixtures/` convention (`schemas/cds/fixtures/valid-receipt.yaml`, `counterfeit-*.yaml`; `schemas/cdd/fixtures/valid-generic-receipt.yaml`, `invalid-*.yaml`). This is the repo's established pattern for positive/negative fixture pairs and is the natural home for AC4/AC5's positive/negative review-request fixtures — **not** a new ad hoc fixture location.
- `grep "status:review"` across `dispatch-protocol/SKILL.md`, `cds-dispatch/SKILL.md`, `cds-dispatch/prompt.md`, `.github/workflows/cnos-cds-dispatch.yml` → confirmed: the actual `status:in-progress → status:review` transition is δ's step 7, currently gated **only** on "β's converge verdict" (cds-dispatch/SKILL.md:222, prompt.md:140, rendered workflow:208). There is no mechanical artifact-existence check anywhere in this chain today — confirming the issue's stated gap directly at the literal call site that will need to change.

Disposition: additive. This cycle adds a new artifact, a new guard, and new doctrine language at the existing surfaces; it does not need to invent a new fixture directory pattern or contradict an existing mechanism.

## Surfaces α is expected to touch

1. `src/packages/cnos.cdd/skills/cdd/CDD.md` — AC1: state the review-request invariant ("no matter, no review; no proof, no `status:review`") inside the existing cell-closure-sequence doctrine. Additive paragraph(s) near the existing `matterₙ`/`reviewₙ`/`receiptₙ` kernel description (lines ~15–35); do not restructure the kernel.
2. `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` — AC3/AC7: the canonical owner of the `status:*` label semantics table (lines ~77–92 list each `status:*` meaning). `status:review`'s row ("Cell complete; awaiting external human/planner review") needs the deliverable-proof precondition named here, since this file is "Source of truth" for dispatch lifecycle rules per the issue's own table.
3. `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` **and** its sibling `prompt.md` — AC3: both files carry near-duplicate content (confirmed: step 7 / line ~222 vs ~140, and the lifecycle-transitions table at ~237/~155). Per the repo's established `install-wake-golden` convention (seen in the cnos#516 guard), `prompt.md` is the source the golden/live workflow is rendered from, while `SKILL.md` is the doctrine copy — **both must be updated identically**, and the rendered `.github/workflows/cnos-cds-dispatch.yml` + `cnos-cds-dispatch.golden.yml` must be re-rendered/kept in sync (see `cn install-wake cds-dispatch` golden-render convention). α must not edit only one of the three and leave the others stale — that would itself reproduce the AC4 failure class (label/contract drift between doc and rendered substrate).
4. `.github/workflows/cnos-cds-dispatch.yml` — re-rendered output of `prompt.md`; α regenerates via the existing `cn install-wake` mechanism rather than hand-editing, consistent with the file's own banner ("rendered via `cn install-wake cds-dispatch --out ...`").
5. A new `scripts/ci/` guard script — AC4/AC5 (see dedicated section below for the pinned design).
6. A new `.cdd/unreleased/<N>/REVIEW-REQUEST.yml` schema/shape — AC2. The issue's own body already contains a complete draft shape (`schema: cdd.review-request.v1`, with `cell`, `issue`, `run_class`, `matter{branch,pr,base_sha,head_sha,changed_files}`, `artifacts{...}`, `checks[]`, `known_gaps[]`, `request{next_state,requested_by,requested_at}`). α documents this shape (in CDD.md or a dedicated doc — α's call) and the guard script (surface 5) parses/validates it. **UNPINNED below**: whether this needs a CUE schema (`schemas/cdd/review_request.cue`) alongside the doc-level shape, matching the `contract.cue`/`receipt.cue`/`boundary_decision.cue` pattern already in `schemas/cdd/`, or whether a doc-level shape + shell-script field checks is sufficient for this cycle. See UNPINNED axis list below.
7. Fixtures — AC5: positive + negative fixture pairs under `schemas/cds/fixtures/` (or a new `scripts/ci/fixtures/` if schema fixtures are judged out of scope for a shell guard — α's call, but must be consistent with existing convention; do not invent a third fixture location).
8. Closeout/doctrine guidance for AC6 — likely lands in the same `dispatch-protocol/SKILL.md` or `cds-dispatch/SKILL.md` surfaces already listed in (2)/(3), not a new file.

No Go/CLI binary surface is implicated — confirmed by reading the issue's own Scope section, which lists only doctrine, prompt, guard, and artifact-schema work.

## Per-AC oracle approach

| AC | Oracle approach |
|---|---|
| AC1 | Grep `CDD.md` for the literal invariant sentence (or a doctrinally-equivalent phrasing β can verify by reading) plus a sentence distinguishing self-check / review-request / beta-review. Mechanical: `grep` in the new repo-pattern guard script (§AC4/AC5 below) can also assert this string exists, folding AC1 into the same guard's "doctrine surface" checks — following the `check-dispatch-repair-preflight.sh` precedent of asserting doctrine strings exist alongside prompt/golden/workflow strings. |
| AC2 | Documentation-presence check: the `REVIEW-REQUEST.yml` shape (the field list above) is written down at a stable path, referenced from CDD.md and/or dispatch-protocol/SKILL.md. Oracle: grep for the field names (`schema:`, `cell:`, `matter:`, `branch:`, `pr:`, `base_sha:`, `head_sha:`, `changed_files:`, `artifacts:`, `checks:`, `known_gaps:`, `request:`, `next_state:`) in the doc surface. |
| AC3 | Grep `cds-dispatch/SKILL.md` + `prompt.md` (both copies) for language requiring `REVIEW-REQUEST.yml` emission and matter-proof before the `status:in-progress → status:review` transition at step 7. Mechanical, same guard script. |
| AC4 / AC5 | **See dedicated design below — this is the load-bearing decision γ is pinning so α does not improvise.** |
| AC6 | Grep closeout-guidance surface(s) for: "closeout must reference matter," "review request must reference matter," "no-matter run must STOP/BLOCKED." Same presence-guard pattern. |
| AC7 | Grep dispatch-protocol/SKILL.md + cds-dispatch surfaces for explicit text preserving the review-request ≠ beta-review ≠ validator-verdict ≠ delta-decision distinction (already implicit in CDD.md's five-step kernel; AC7 asks that `status:review`'s prose not be rewritten to imply acceptance). Negative case: confirm no edited sentence states or implies `status:review` ⇒ PASS/accepted. |
| AC8 | Diff review: confirm no `gh label create` / label-taxonomy-table edits beyond what's needed to add the proof-precondition prose to the existing `status:review` row. β checks the diff for new label names — none should appear. |
| AC9 | Run existing CI gates locally/in CI: I1, I2, I4, I5, I6, install-wake-golden, dispatch-repair-preflight, Go, Package, Binary. All must stay green — standard β verification, no new oracle needed beyond "CI is green on the PR." |
| AC10 | This cycle's own `.cdd/unreleased/532/` receipt and closeouts must honestly name what changed, how the guard works, the negative-proof evidence, and known gaps — verified by γ at close-out triage per `gamma/SKILL.md §2.7`, not by α's own guard script (avoids the self-grading problem AC1's "no self-certify" principle would otherwise contradict if applied reflexively). |

## AC4/AC5 guard design — pinned, not left to α's improvisation

**This is the load-bearing design decision from the operator's clarification comment** (posted 2026-06-30T19:29:09Z, immediately before claim). The operator explicitly distinguishes this guard's *shape* from the cnos#516 `check-dispatch-repair-preflight.sh` precedent:

> "For the current GitHub/CDS dispatch surface, deliverable proof means: a PR exists … the PR has commits beyond base … the branch exists … the diff is non-empty unless explicitly declared as no-op … `.cdd/unreleased/<N>/REVIEW-REQUEST.yml` exists … required closeout artifacts exist … the review request names the PR/branch/head SHA/checks."
>
> "Prompt text alone is not sufficient."

This means: **AC4's guard is structurally different from the #516 precedent and must not be implemented as a copy of it.**

`check-dispatch-repair-preflight.sh` is a **presence-of-contract** guard: it greps doctrine/prompt/golden/workflow files for required substrings, proving the *prompt still tells the cell what to do*. It runs unconditionally in CI on every PR, regardless of which issue's cycle is open, because it is checking that static doc/prompt text wasn't regressed.

AC4 needs a **deliverable-existence** guard: it must check that a *specific cell's actual matter* — PR, branch, diff, artifacts — exists for the issue being transitioned to `status:review`. This is checking dynamic, per-cycle state (a PR number, a branch name, a diff), not static doc text. A presence-of-contract grep over `prompt.md` cannot prove a PR exists; it can only prove the prompt *asks for* one.

**Pinned design:**

1. **Two-guard split**, not one guard wearing two hats:
   - **Guard A — doctrine/prompt presence** (`scripts/ci/check-review-request-preflight.sh`, modeled directly on the #516 precedent's `need()` helper pattern): asserts the required language (the AC1/AC3/AC6/AC7 strings) is present in `CDD.md`, `dispatch-protocol/SKILL.md`, `cds-dispatch/SKILL.md`, `cds-dispatch/prompt.md`, the golden, and the rendered live workflow. This runs unconditionally in CI exactly like the #516 guard — no network/gh calls, pure `grep -qF` over checked-in files. **This part is deterministic and offline by construction**, same as the precedent.
   - **Guard B — deliverable-existence check**: a function (can live in the same script or a sibling, α's call, but must be named and tested separately from Guard A) that takes a `REVIEW-REQUEST.yml` path (or a directory) as input and validates: `matter.pr` is a non-empty integer-shaped value, `matter.branch` matches `^cycle/\d+$`, `matter.base_sha`/`matter.head_sha` are present and unequal (proves commits-beyond-base without calling `git log` against a remote — see offline-testability below), `matter.changed_files` is a non-empty list (proves non-empty diff) unless a sibling `no_op: true` + `no_op_approval:` field is present, `artifacts.*` paths exist on disk relative to repo root, `request.next_state == "status:review"`, `request.requested_by` and `request.requested_at` are present.

2. **Offline/deterministic testability for Guard B is solved by treating `REVIEW-REQUEST.yml` as the oracle's *input*, not by having the guard re-derive proof via live `gh`/network calls.** The guard does not call `gh pr view` or `git fetch` to verify a PR exists on GitHub — that would make CI non-deterministic and dependent on live network/API state, which the issue's own AC4 oracle line ("CI must be deterministic and offline-testable") forbids. Instead:
   - `matter.head_sha` / `matter.base_sha` are checked **structurally** (present, well-formed git-SHA-shaped strings, not equal to each other) and, where the guard runs inside the actual PR's CI job, **cross-checked against the job's own `GITHUB_SHA`/`git rev-parse HEAD`** — that part *is* live but uses data CI already has for free (no extra API call), exactly as the existing `install-wake-golden` and `dispatch-repair-preflight` jobs already do.
   - `matter.pr` existing is proven by the **workflow trigger context** (`github.event.pull_request.number` / `gh pr view --json number` is acceptable **inside an actual PR-triggered CI run**, since that's a real PR, not a guard-script unit test) — but the **unit-testable** part of Guard B (what α writes tests for, what AC5's fixture proves) operates purely on a static `REVIEW-REQUEST.yml` fixture file and the filesystem, with no network/gh dependency at all. This is the key design split: **CI-context integration (live, runs once per real PR) vs. guard-logic unit test (offline, runs on fixtures, is what AC5 actually proves).**
   - AC5's "synthetic or fixture test" is satisfied by Guard B's **fixture-driven unit test**, not by re-running the real #524 W4 scenario against live GitHub state. α creates fixture `REVIEW-REQUEST.yml` files (valid one + several invalid ones: missing PR field, equal base/head SHA, empty `changed_files` with no `no_op` declaration, missing artifact path) under the established `schemas/cds/fixtures/` (or `scripts/ci/fixtures/review-request/`, α's call per the "Surfaces" section item 7 above) and the guard script's test entry point runs Guard B's validation function against each fixture, asserting pass/fail matches expectation. This is exactly the shape of `schemas/cds/fixtures/valid-receipt.yaml` + `counterfeit-*.yaml` and `schemas/cdd/fixtures/valid-generic-receipt.yaml` + `invalid-*.yaml` — α should follow that naming convention (`valid-review-request.yaml`, `invalid-review-request-missing-pr.yaml`, `invalid-review-request-no-diff.yaml`, etc.) rather than inventing new naming.
   - The **#524 W4 negative case specifically** ("status:review with no PR/commits/diff") is reproduced as one of these fixtures: a `REVIEW-REQUEST.yml`-shaped (or *absent* — the missing-file case is itself a fixture scenario) input representing that exact historical state, asserted to fail Guard B.

3. **Where Guard B's check actually gates the real transition**: the cds-dispatch wake's step 7 (the `status:in-progress → status:review` label transition) is the integration point. α's prompt/doctrine changes (surface 3 above) must make step 7 read: "before applying `status:review`, run Guard B against this cycle's `REVIEW-REQUEST.yml`; if it fails, do not transition — post STOP/BLOCKED instead" (per the operator's "must post STOP/BLOCKED or equivalent evidence" requirement). This is a **prompt-level instruction calling a script**, not a separate CI job gating the merge — the operator's clarification frames this as gating the *label transition event*, which happens inside the dispatch wake's own run, not as a GitHub branch-protection check on the PR. α should confirm this reading is correct against the issue body's "Surface: CI / dispatch guard" phrasing (both readings are textually supported; γ flags this as the one remaining interpretive choice below).

**UNPINNED: Guard-B invocation point — "CI job on the PR" vs "step inside the dispatch wake's own run before it applies the label."** The issue's AC4 says "Surface: CI / dispatch guard" (both, ambiguously, separated by a slash). The operator clarification says "The guard may be implemented as a workflow/script/post-step" (also ambiguous — could be a workflow *job* or a wake-internal *post-step*). γ's reading above (favoring wake-internal post-step, since that's the actual point where the label transition happens and where STOP/BLOCKED commentary is posted) is the more literal match to "the issue must not be left in `status:review` without proof" — a PR-level CI job can be green/red but cannot itself *prevent* the label edit that already happened in a separate wake run. **δ should confirm this reading before α starts**, since it changes whether Guard B's primary integration is `.github/workflows/cnos-cds-dispatch.yml` (a step inside the existing dispatch workflow, right before the label-transition step) or a net-new standalone workflow file triggered on `pull_request` / `issues: labeled (status:review)`. γ's recommendation: wake-internal post-step, consistent with how `dispatch-repair-preflight`'s logic is *also* invoked from inside `cds-dispatch/prompt.md`'s own steps (not a separate workflow) — Guard A already follows this pattern for #516; Guard B should follow the same integration shape for consistency, just gating a different transition.

## Expected diff scope

- `src/packages/cnos.cdd/skills/cdd/CDD.md` — additive doctrine paragraph(s), no restructuring.
- `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` — additive precondition language on the `status:review` row + nearby prose.
- `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` + `prompt.md` — additive step-7 language requiring `REVIEW-REQUEST.yml` + Guard B pass before the label transition; both files edited in lockstep.
- `.github/workflows/cnos-cds-dispatch.yml` (+ `cnos-cds-dispatch.golden.yml`) — regenerated via `cn install-wake`, not hand-edited, to stay byte-identical to `prompt.md` per the existing `install-wake-golden` CI convention.
- New file: `scripts/ci/check-review-request-preflight.sh` (Guard A, presence-of-contract, #516-style) — and either a second script or a clearly-separated function for Guard B (deliverable-existence), per α's structuring choice, but both must exist and both must be independently testable.
- New fixtures under the existing `schemas/cds/fixtures/` (or equivalent existing fixture convention) — valid + multiple invalid `REVIEW-REQUEST.yml` cases, including the #524 W4 reproduction.
- New `.github/workflows/` CI wiring to invoke the new guard script(s) — minimal, following the existing job pattern in `ci.yml` / `cnos-cds-dispatch.yml` (α reads the existing job definitions for the #516 guard to match shape).
- No Go source changes. No CLI (`cn`) subcommand changes beyond the existing `cn install-wake` re-render, which is mechanical regeneration, not new command surface.

## Explicit non-goals (copied verbatim from the issue body)

Out of scope:
- Demo 0.
- Wake-as-skill W4.
- Changing dispatch labels broadly.
- Changing GitHub issue status taxonomy.
- Rewriting old cycles.
- Editing frozen releases.
- Implementing full projection artifacts.
- Making beta review automatic.
- Changing who owns final delta decisions.

Deferred:
- Moving review-state transition to a separate deterministic controller.
- Full projection schema integration.
- TSC coherence requirement for review requests.
- Domain-specific CDS review-request extensions beyond the minimum proof fields.

Active design constraints (also binding non-goals in effect):
- Do not redesign the label taxonomy.
- Do not start Demo 0.
- Do not implement projection in this issue.
- Do not change wake source model.
- Do not change wake rendering behavior (beyond the mechanical re-render of step-7 content).
- Do not delete wake-provider JSON/prompt files.
- Do not make TSC mandatory.

## Findings for δ (per gamma/SKILL.md §3.5 — artifact facts, not invented scope)

1. **UNPINNED axis surfaced above** — Guard B invocation point (CI job vs wake-internal post-step). γ's recommendation is wake-internal post-step for consistency with the #516 precedent's own integration shape, but this is a genuine design choice the issue text leaves ambiguous (AC4's "Surface: CI / dispatch guard" reads as a deliberate either/or). δ should pin this before α starts, or explicitly delegate the choice to α with the reasoning above as the default.
2. The issue body is otherwise complete and well-formed per `issue/SKILL.md`'s handoff checklist — ACs are numbered and independently testable, non-goals exist and don't leak into ACs, source-of-truth paths all resolved on inspection, proof plan has oracle/positive/negative cases, related artifacts link to real existing paths (all four confirmed to exist: `CDD.md`, `dispatch-protocol/SKILL.md`, `cds-dispatch/SKILL.md`, `check-dispatch-repair-preflight.sh`). No issue-quality gap found requiring a γ issue-edit before dispatch.
3. `cds-dispatch/SKILL.md` and `prompt.md` are near-duplicates (confirmed by grep — same step numbering, same line content at offset ~82 lines). α must edit both; β's review should explicitly check both were touched identically, not just one.

---

## α dispatch prompt

```text
You are α. Project: cnos.

Load src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md and follow its load order.

Issue: gh issue view 532 --repo usurobor/cnos --json title,body,state,comments

Branch: cycle/532

Read the issue's operator-clarification comment (posted 2026-06-30T19:29:09Z) carefully —
it is the authoritative interpretation of AC4/AC5's "deliverable proof" requirement and
sharpens what the guard must mechanically check. Do not implement against the issue body
alone; the comment overrides/specifies where the body is ambiguous.

Read .cdd/unreleased/532/gamma-scaffold.md on this branch before starting. It names the
expected surfaces, the per-AC oracle approach, and — load-bearing — the pinned design for
the AC4/AC5 guard (the two-guard split: Guard A = presence-of-contract, modeled on
scripts/ci/check-dispatch-repair-preflight.sh; Guard B = deliverable-existence, fixture-driven,
offline-testable). Do not improvise a different guard shape; if the scaffold's design does not
fit what you find once you start, surface the conflict via .cdd/unreleased/532/ artifacts rather
than silently diverging.

Tier 3 skills: src/packages/cnos.core/skills/write/SKILL.md, src/packages/cnos.cdd/skills/cdd/issue/SKILL.md, src/packages/cnos.cdd/skills/cdd/issue/contract/SKILL.md, src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md, src/packages/cnos.cdd/skills/cdd/issue/constraints/SKILL.md, src/packages/cnos.core/skills/eng/SKILL.md

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | Markdown (CDD/dispatch-protocol/CDS doctrine edits) + Bash (new guard script(s) at `scripts/ci/`, matching the existing `scripts/ci/*.sh` convention, e.g. `check-dispatch-repair-preflight.sh`) + YAML (the `REVIEW-REQUEST.yml` shape + fixture files under the existing `schemas/cds/fixtures/`-style convention) |
| CLI integration target | N/A — no new `cn` subcommand. The only CLI interaction is the existing `cn install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml` regeneration mechanism, used to re-render `prompt.md` changes into the live workflow + golden, not a new command. |
| Package scoping | Doctrine edits land at their existing canonical paths: `src/packages/cnos.cdd/skills/cdd/CDD.md`, `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md`, `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` + `prompt.md` (edit both, lockstep). New guard script(s) at `scripts/ci/`. New fixtures under the existing fixture convention (`schemas/cds/fixtures/` or equivalent — match existing naming, do not invent a new fixture root). |
| Existing-binary disposition | N/A — no existing binary is replaced, deprecated, or modified. This cycle adds doctrine + a guard script + fixtures only. |
| Runtime dependencies | None new. Guard script(s) use bash + the tooling already available to `scripts/ci/*.sh` (e.g. `grep`, standard POSIX utilities, optionally `yq`/`python3` if needed for YAML field extraction — check what `scripts/ci/validate-skill-frontmatter.sh` already relies on before introducing a new YAML-parsing dependency). No network or live `gh`/GitHub API calls inside the unit-testable guard logic (Guard B must be offline-testable against fixtures per the scaffold's design). |
| JSON/wire contract preservation | The new `.cdd/unreleased/<N>/REVIEW-REQUEST.yml` schema is additive — a new artifact, not a modification of any existing wire contract (`schemas/cdd/receipt.cue`, `schemas/cds/receipt.cue`, `boundary_decision.cue`, etc. are untouched). No existing JSON/YAML contract changes shape. |
| Backward-compat invariant | Existing `dispatch-protocol/SKILL.md` and `cds-dispatch/SKILL.md` + `prompt.md` content is preserved — this is additive (new required artifact + new gate before the `status:in-progress → status:review` transition), not a rewrite. No existing label taxonomy changes (AC8 — verify no `gh label create`/label-rename touches the diff). No existing lifecycle transition table rows are removed or renamed; only the `status:review` row/step-7 prose gains a precondition. |

**UNPINNED: Guard-B invocation point — CI job on the PR vs a step inside the dispatch wake's own run (before it applies the `status:review` label).** See `.cdd/unreleased/532/gamma-scaffold.md` §"AC4/AC5 guard design" for the full reasoning. γ's recommendation is the wake-internal post-step (consistent with how the #516 `dispatch-repair-preflight` guard is also invoked from inside `cds-dispatch/prompt.md`'s own steps, not as a separate workflow). If δ has not resolved this before you start, treat γ's recommendation as the working default but name the choice explicitly in your `self-coherence.md` so β can verify the reasoning, not just the outcome.

If you encounter any other axis that is genuinely ambiguous after reading the issue and the scaffold, do not pick a value — surface it via the artifact channel before proceeding.
```

## β dispatch prompt

```text
You are β. Project: cnos.

Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.

Issue: gh issue view 532 --repo usurobor/cnos --json title,body,state,comments

Branch: cycle/532

Read .cdd/unreleased/532/gamma-scaffold.md on this branch — it names the pinned AC4/AC5
guard design (the two-guard split and the offline-testability requirement) and the
operator-clarification comment's binding interpretation of deliverable proof. Verify α's
diff against both the issue's 10 ACs and the scaffold's pinned design — in particular,
confirm Guard B is fixture-driven and does not depend on live network/gh calls in its
unit-testable path, and confirm cds-dispatch/SKILL.md and prompt.md were both edited
identically (the scaffold names them as near-duplicate surfaces that must move in lockstep).
Verify the diff against the 7-axis Implementation contract in the α prompt per
beta/SKILL.md Role Rules Rule 7.
```

