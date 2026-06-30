# self-coherence — cycle/532

manifest:
  completed: [Gap, Skills, ACs, Self-check, Debt, CDD Trace]

## Gap

Issue: [#532](https://github.com/usurobor/cnos/issues/532) — "cdd/cds: require review-request proof before status:review"

A dispatch cell was able to set an issue to `status:review` without producing reviewable matter (observed in #524 W4: `status:review` + no PR + no commits beyond base + no diff + no receipt + no closeout + no STOP comment). CDD's closure sequence (`contract → matter → review → receipt → verdict → decision`) does not currently gate the `status:review` label transition on any mechanical proof that matter exists — the label is writable by the cell on a self-assertion alone.

The operator's pre-claim clarification comment (2026-06-30T19:29:09Z) is the authoritative, binding interpretation of AC4/AC5 and overrides any ambiguity in the issue body: the guard must be a mechanical state-integrity check, not prompt guidance alone; deliverable proof for the current GitHub/CDS surface means a PR (or approved no-op), commits beyond base (or no-op approval), a non-empty diff (or declared no-op), a `REVIEW-REQUEST.yml`, required closeout artifacts, and a review request that names PR/branch/head SHA/checks; on missing proof the guard must fail and the run must post STOP/BLOCKED, not leave `status:review` standing unproven.

**Version / mode:** `design-and-build` per the issue header (no stable `DESIGN.md`/`PLAN.md` existed at scoping time; γ confirmed MCA preconditions did not hold). This cycle both designs the `REVIEW-REQUEST.yml` shape (drafted as a strong proposal in the issue body, not yet a converged schema doc) and builds the guard in the same cycle, per the issue's own "Cycle scope sizing" decision (keep whole as one design-and-build issue).

**δ's binding correction to γ's scaffold (followed, not γ's mis-citation):** γ's scaffold claimed the cnos#516 guard is "invoked from inside `cds-dispatch/prompt.md`'s own steps" — δ verified this is factually wrong: `check-dispatch-repair-preflight.sh` runs as a standalone CI job (`dispatch-repair-preflight` in `.github/workflows/build.yml`), triggered on every PR, unrelated to the dispatch wake's own run. This correction does not change γ's underlying two-guard-split design (which δ pinned as final), only where Guard A is wired (a new CI job mirroring the #516 job's exact shape) and confirms Guard B's integration point (the dispatch wake's own run, immediately before the `status:in-progress → status:review` label transition) on its own independent merits — not because of the #516 mis-citation. See §Self-check below for how this resolution was followed exactly as δ pinned it.

## Skills

Tier 1 (loaded per dispatch §2.1, always):

- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical lifecycle and role contract.
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface; governed this cycle's execution (§2.2 artifact order, §2.3 peer enumeration, §2.6 pre-review gate, §2.7 review-readiness signal).

Tier 3 (named in the dispatch prompt; all loaded):

- `src/packages/cnos.core/skills/write/SKILL.md` — governed prose density in the doctrine additions (CDD.md, dispatch-protocol/SKILL.md §"Review-request proof gate", cds-dispatch step-7 edits): front-loaded points, one governing question per section, no throat-clearing.
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — read to interpret AC boundaries and verify the issue's own handoff-checklist completeness (confirmed complete by γ's scaffold; re-verified independently during implementation — ACs numbered/testable, non-goals present and not leaking into ACs, source-of-truth paths all resolve).
- `src/packages/cnos.cdd/skills/cdd/issue/contract/SKILL.md` — governed reading the issue's Problem/Status-truth/Scope/Non-goals sections precisely (e.g. confirming "Status truth" table's "Missing" claims before trusting them — re-verified via `rg -i review-request` repo-wide, 0 real hits beyond γ's own scaffold).
- `src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md` — governed AC-to-oracle mapping; every AC below names a concrete, runnable oracle (script invocation or grep), not "it works."
- `src/packages/cnos.cdd/skills/cdd/issue/constraints/SKILL.md` — governed the constraint-strata reading: AC8's "no label taxonomy change" is a hard gate (verified via diff grep for `gh label create`/rename, 0 hits beyond gamma-scaffold's own discussion text); cross-surface projections enumerated in §Self-check below.

Not found (noted, not blocking per dispatch instruction): `src/packages/cnos.core/skills/eng/SKILL.md` does not exist at that path — confirmed via `ls` before starting; no Go/CLI surface was touched this cycle (δ's pinned contract: "No Go source changes"), so this absence had no material effect on the work.

Role-skill-only constraint honored: did not load β/γ/δ role skills as primary governance (per `alpha/SKILL.md` §2.1 "do not load β or γ role skills"); `delta/SKILL.md` and `gamma-scaffold.md` were read as **input artifacts** (the dispatch prompt and the scaffold), not as α's own governing skill — consistent with the role boundary.

## ACs

### AC1: CDD doctrine states the invariant

**Met.** `src/packages/cnos.cdd/skills/cdd/CDD.md` gained a new §"Review-request invariant (no self-certified review state)" immediately after the kernel statement (commit `a9f4957`). It states verbatim:

> No matter, no review.
> No proof, no `status:review`.

and distinguishes self-check / review-request / β-review as three separate acts a cell must not collapse.

Oracle evidence: `grep -F "No matter, no review." src/packages/cnos.cdd/skills/cdd/CDD.md` → 1 hit. Mechanically asserted by Guard A (`scripts/ci/check-review-request-preflight.sh --guard-a`), which I ran and confirmed passes:

```
$ ./scripts/ci/check-review-request-preflight.sh --guard-a
cnos#532 Guard A (presence-of-contract): review-request doctrine present (CDD + dispatch-protocol + cds-dispatch SKILL.md + prompt.md).
```

Negative case: the doc does not describe a cell as β-reviewing its own matter — the new section explicitly says α "may not grant itself review state."

### AC2: Review-request artifact is defined

**Met.** `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` gained §"Review-request proof gate (cnos#532)" (commit `aba6628`) documenting the full `REVIEW-REQUEST.yml` shape: the example YAML block (taken from the issue's own draft, unmodified) plus a field-by-field table naming required-vs-no_op-conditional fields (`schema`, `cell`/`issue`, `run_class`, `matter.branch`/`pr`/`base_sha`/`head_sha`/`changed_files`, `artifacts.*`, `checks`, `known_gaps`, `request.next_state`/`requested_by`/`requested_at`, `no_op`/`no_op_approval`).

Oracle evidence: `grep -c "REVIEW-REQUEST.yml\|cdd.review-request.v1" src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` → multiple hits across the schema block and field table. Asserted mechanically by Guard A.

Negative case: the field table states explicitly — "A review request without a PR or changed matter is invalid unless explicitly declared `no_op: true` with `no_op_approval:` present — an undeclared empty diff fails the gate." Guard B mechanically enforces this (see AC4/AC5 below); fixture `invalid-empty-diff-no-no_op.yaml` proves it.

### AC3: CDS dispatch prompt requires review-request before status:review

**Met.** Both `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` and `prompt.md` (near-duplicate surfaces; confirmed byte-identical at the edited region via `diff <(sed -n '218,239p' SKILL.md) <(sed -n '136,157p' prompt.md)` → 0 output) edited identically at step 7 (commit `f4d9861`): step 7 now requires writing `REVIEW-REQUEST.yml` and running `scripts/ci/check-review-request-preflight.sh --guard-b --fixture ... ` (live mode) before the `status:in-progress → status:review` label transition. The lifecycle-transitions table's "β converge verdict" row was renamed to "β converge verdict **and** Guard B deliverable-proof pass," naming both as jointly necessary.

Oracle evidence: Guard A asserts `REVIEW-REQUEST.yml`, `check-review-request-preflight.sh`, `Guard B`, `STOP/BLOCKED`, `deliverable-proof pass` are present in both files — confirmed passing (see AC1 invocation above).

Negative case: prompt no longer reads "only on β's converge verdict transitions the cell's label" (the old text implying β-converge alone is sufficient) — it now requires both predicates jointly, and explicitly states "the wake does NOT apply `status:review`" on guard failure.

### AC4: status:review transition has mechanical guard

**Met**, per δ's pinned two-guard design. `scripts/ci/check-review-request-preflight.sh` implements:

- **Guard A** (`--guard-a`) — presence-of-contract, wired into `.github/workflows/build.yml` as job `review-request-preflight` (commit `c2dae48`), runs on every PR unconditionally, no network/gh calls.
- **Guard B** (`--guard-b --fixture <path>` / `--guard-b --live <path>`) — deliverable-existence check against a `REVIEW-REQUEST.yml`: `matter.pr` present + integer-shaped (unless `no_op`), `matter.branch` matches `^cycle/[0-9]+$`, `matter.base_sha`/`matter.head_sha` present and unequal, `matter.changed_files` non-empty unless `no_op:true`+`no_op_approval:` present, every `artifacts.*` path exists on disk, `request.next_state == status:review`, `request.requested_by`/`requested_at` present. Structural mode (`--fixture`) is offline/deterministic (pure awk/grep + filesystem checks, no network or git-remote/gh calls); live mode (`--live`) additionally cross-checks `head_sha` against `git rev-parse HEAD` and may use `gh pr view` — used by the dispatch wake's real invocation per δ's pinned design.

Invocation point: per δ's explicit resolution (not γ's mis-cited #516 recommendation — see §Gap), Guard B is invoked by the dispatch wake itself (`cds-dispatch/SKILL.md` + `prompt.md` step 7, `delta/SKILL.md` §9.5/§9.6) immediately before the `status:in-progress → status:review` label transition — not as a separate `pull_request`-triggered workflow job, since the event being gated (a label edit inside the wake's own run) is not observable to a generic PR-triggered CI job. No new `pull_request`-triggered workflow file was added for Guard B, per the dispatch contract's explicit instruction.

Oracle evidence — positive case (valid review request passes):

```
$ ./scripts/ci/check-review-request-preflight.sh --guard-b --fixture scripts/ci/fixtures/review-request/valid/valid-review-request.yaml
cnos#532 Guard B: scripts/ci/fixtures/review-request/valid/valid-review-request.yaml — deliverable proof valid.
```

Oracle evidence — negative case (empty-review state fails): see AC5 below for the full fixture matrix; every invalid fixture fails with a named reason.

### AC5: Previous failure mode is covered

**Met.** `scripts/ci/check-review-request-preflight.sh --self-test` runs Guard B's offline structural mode over 9 fixtures under `scripts/ci/fixtures/review-request/{valid,invalid}/` plus a missing-file case, asserting pass/fail matches the directory. Full run output (10 assertion lines, all passing — i.e. all 10 fixture checks behaved as expected):

```
$ ./scripts/ci/check-review-request-preflight.sh --self-test
  ✓ valid:   valid-no_op-with-approval.yaml — passed (expected)
  ✓ valid:   valid-review-request.yaml — passed (expected)
  ✓ invalid: invalid-524-w4-empty-review.yaml — failed (expected): ...matter.pr missing...; matter.base_sha == matter.head_sha...; matter.changed_files empty...; artifacts block missing or empty
  ✓ invalid: invalid-empty-diff-no-no_op.yaml — failed (expected): matter.changed_files empty or missing and no_op is not declared
  ✓ invalid: invalid-missing-artifact.yaml — failed (expected): artifacts path(s) do not exist on disk
  ✓ invalid: invalid-missing-pr.yaml — failed (expected): matter.pr missing or not integer-shaped
  ✓ invalid: invalid-missing-request-fields.yaml — failed (expected): request.requested_by missing; request.requested_at missing
  ✓ invalid: invalid-no-commits-beyond-base.yaml — failed (expected): matter.base_sha == matter.head_sha
  ✓ invalid: invalid-no_op-without-approval.yaml — failed (expected): no_op: true but no_op_approval is absent
  ✓ missing-file (#524 W4 repro) — failed (expected): no REVIEW-REQUEST.yml, no matter, no review
cnos#532 Guard B self-test: all fixtures matched expectation (valid pass, invalid fail, #524 W4 repro fails).
```

**#524 W4 reproduction specifically** (issue's required negative case — "status:review + no PR/diff/matter must fail mechanically"): covered two ways — (1) `invalid-524-w4-empty-review.yaml`, a fixture with no PR, `base_sha == head_sha` (no commits beyond base), empty `changed_files`, and an empty `artifacts` block, all failing simultaneously; (2) the missing-file case, representing a cell that never wrote `REVIEW-REQUEST.yml` at all. Both fail Guard B.

This fixture suite is wired into CI (`.github/workflows/build.yml` job `review-request-preflight`, step "Check Guard B fixture suite") so AC5 is mechanically re-verified on every PR, not just locally.

### AC6: Closeout guidance includes repair/evidence distinction

**Met.** `dispatch-protocol/SKILL.md` §"Review-request proof gate" §"On guard failure" names the required behavior: "If Guard B fails, the wake MUST NOT apply `status:review`. It posts a STOP/BLOCKED comment naming the missing proof and leaves the cell at `status:in-progress` (or transitions to `status:blocked`...)." `delta/SKILL.md` §9.5/§9.6 (commit `7a44d84`) mirrors this from δ's own canonical token-write perspective: "On Guard B failure, δ writes `status:blocked` instead, naming the missing proof." `cds-dispatch/SKILL.md` + `prompt.md` step 7 state the same STOP/BLOCKED requirement.

Oracle evidence: Guard A asserts the literal string "STOP/BLOCKED" is present in `dispatch-protocol/SKILL.md`, `cds-dispatch/SKILL.md`, and `prompt.md` — confirmed passing.

Negative case: none of the edited surfaces permit an empty run to claim CONVERGE or `status:review` — the new text is unconditional ("MUST NOT apply `status:review`").

### AC7: status:review remains review, not verdict

**Met.** The `status:review` label-meaning row in `dispatch-protocol/SKILL.md`'s lifecycle-labels table now reads: "`status:review` means 'ready for β review,' never 'accepted'; it carries no verdict." in addition to naming the new precondition. No edited sentence anywhere in this diff states or implies `status:review` ⇒ PASS/accepted — verified by re-reading every edited paragraph in `dispatch-protocol/SKILL.md`, `cds-dispatch/SKILL.md`, `prompt.md`, and `delta/SKILL.md` for this cycle's diff (`git diff origin/main..HEAD -- <those files>`), none of which collapses review-request, β-review, validator verdict (V), or δ's decision into a single act — the new prose explicitly separates them ("this gate only proves the *matter* a review request claims to point at is real" — `dispatch-protocol/SKILL.md`).

Negative case: confirmed no edited sentence states `status:review` implies acceptance — `grep -n "status:review" <edited files> | grep -i "accept\|pass\b"` → 0 hits (the only "PASS"-adjacent text is Guard A/B's own pass/fail vocabulary for the *guard's* verdict on proof-presence, not a claim about the cell's work being accepted).

### AC8: No broad label redesign

**Met.** `git diff origin/main..HEAD | grep -i "gh label create\|gh label delete"` → 0 hits. No new `status:*` label was created, renamed, or removed; the only change to the `status:review` table row is additive prose (a precondition clause appended to the existing "Cell complete; awaiting external human/planner review" meaning). The lifecycle-transitions table in `cds-dispatch/SKILL.md`/`prompt.md` keeps the same four named events (claim, β-converge, hard-block, release-back-to-queue) — only the β-converge row's *label* changed from "β converge verdict" to "β converge verdict **and** Guard B deliverable-proof pass," and its *to*-state and *from*-state are unchanged (`status:in-progress` → `status:review`).

### AC9: Existing gates remain green

**Partially verified locally; full confirmation is CI's job.** See §Self-check below for the explicit local-vs-CI-only breakdown per gate.

### AC10: Receipt is honest

**In progress as of this section.** This `self-coherence.md` document is the receipt for this cycle; §Debt below names every known gap honestly (no false CONVERGE claimed for anything not actually verified). α does not write `alpha-closeout.md` in this dispatch (per `alpha/SKILL.md` §2.8, close-out happens at re-dispatch after β's merge in bounded-dispatch mode) — γ's eventual close-out triage and this self-coherence document together satisfy AC10's receipt-honesty requirement for the α-side portion of the cycle.

## Self-check

### Guard-B-invocation-point decision — followed δ's resolution, not γ's mis-cited recommendation

γ's scaffold flagged the Guard-B invocation point as UNPINNED and recommended a "wake-internal post-step" — but justified that recommendation by claiming consistency with how "the #516 `dispatch-repair-preflight` guard is also invoked from inside `cds-dispatch/prompt.md`'s own steps (not a separate workflow)." δ's dispatch explicitly named this citation as factually wrong (`check-dispatch-repair-preflight.sh` is a standalone CI job, `dispatch-repair-preflight`, in `.github/workflows/build.yml`, unrelated to the wake's own run) and re-pinned Guard B's invocation point on independent grounds: "the event being gated (a label transition on a specific issue) happens inside the wake's own run, and a generic PR-triggered CI job cannot itself prevent a label edit that occurs in a separate workflow run."

I implemented per δ's resolution, explicitly: Guard B is invoked by the dispatch wake's own steps (`cds-dispatch/SKILL.md` + `prompt.md` step 7, `delta/SKILL.md` §9.5/§9.6's token-write specification) immediately before the `status:in-progress → status:review` label transition — not as a `pull_request`-triggered workflow job. I did not add a new `.github/workflows/` job for Guard B's live mode; the only new CI job (`review-request-preflight` in `build.yml`) runs Guard A's doctrine-presence checks plus Guard B's **offline fixture self-test** (structural mode only, no live cross-checks), which is a distinct thing from "Guard B live-gating a PR" — exactly as δ's contract specified ("Guard B itself does NOT need a separate `.github/workflows/` job").

### Peer enumeration

**Role-skill peers** (per `alpha/SKILL.md` §2.3): the `status:in-progress → status:review` label transition is described independently in three role-adjacent surfaces beyond the four δ explicitly named (CDD.md, dispatch-protocol/SKILL.md, cds-dispatch/SKILL.md, prompt.md). I found and updated a fourth: `delta/SKILL.md` §9.5/§9.6 — δ's own canonical specification of when it writes the `status:review` return token. This was not named in δ's dispatch contract's "Surfaces α is expected to touch" list, but leaving it stale would have been exactly the lifecycle-skill-peer-drift failure mode `alpha/SKILL.md` §2.3 warns against (its own citation: "three of four R1 findings landed in lifecycle skills... the audit's checklist did not distinguish the two classes"). Updated in commit `7a44d84`.

I additionally found `wake-provider.json`'s `responsibilities` array carries the same prose independently again (a fifth copy). I confirmed via re-render (`cn-install-wake cds-dispatch` — golden unchanged) that this field is pure manifest documentation, not substituted into the rendered workflow, so updating it carries no golden-drift risk. Updated for consistency in the same commit.

Checked and found **not** peers (no `status:review` text present, confirmed via `grep -ln "status:review"` over the candidate set): `gamma/SKILL.md`, `beta/SKILL.md`, `review/SKILL.md`. `operator-review/SKILL.md` does reference `status:review` but only describes the **downstream** operator-initiated `status:review → status:changes` transition (after the cell already reached `status:review`), which this cycle's gate does not change — left untouched, correctly, since editing it would be out of scope (the issue's non-goal: "do not change who owns final delta decisions").

### Harness audit (schema-bearing change: REVIEW-REQUEST.yml)

`REVIEW-REQUEST.yml` is a new schema-bearing artifact. Enumerated producers and consumers:

- **Producer**: the dispatch wake / δ, per the doctrine edits in `dispatch-protocol/SKILL.md`, `cds-dispatch/SKILL.md`+`prompt.md`, `delta/SKILL.md` — all updated to require emitting it before requesting `status:review`.
- **Consumer**: `scripts/ci/check-review-request-preflight.sh`'s Guard B (both structural and live modes) — the only consumer in this cycle's scope. No other script, Go binary, or CI job reads `REVIEW-REQUEST.yml` (`grep -rl "REVIEW-REQUEST" --include="*.sh" --include="*.go"` outside the new guard script and its own doc references → only the guard script and the doctrine prose).
- **Non-primary-language writers**: none exist yet to drift — this is a net-new artifact with no prior shell/CI/template/fixture writer to audit for staleness. The fixture YAML files I authored under `scripts/ci/fixtures/review-request/` are themselves the only other writers of the shape, and they are hand-authored test fixtures (not a second production writer), consistent with the existing `schemas/cds/fixtures/` / `schemas/cdd/fixtures/` convention of hand-authored positive/negative YAML.
- **Schema-validation runtime dependency check** (per δ's pinned contract): confirmed `scripts/ci/validate-skill-frontmatter.sh` uses `cue` + `jq` for its schema, but that schema is CUE-based and JSON-targeted (skill frontmatter, extracted to JSON first); it is not a YAML-parsing precedent. No script in `scripts/ci/` parses YAML structurally before this cycle. Per δ's instruction ("keep Guard B's YAML field extraction simple — grep/sed/awk is acceptable for a flat-ish schema"), I used pure `awk` (no `yq`/`python3`/new dependency) for `REVIEW-REQUEST.yml`'s flat-ish two-level-nesting shape.

### Caller-path trace for new modules

The only new "module" is `scripts/ci/check-review-request-preflight.sh` itself (a script, not a library function, so the relevant trace is: who invokes it). Callers, all confirmed by direct grep/read:

- `.github/workflows/build.yml` job `review-request-preflight` → `./scripts/ci/check-review-request-preflight.sh --guard-a` and `--self-test` (CI invocation, runs on every PR).
- `cds-dispatch/SKILL.md` + `prompt.md` step 7 (prose instruction naming the exact invocation `scripts/ci/check-review-request-preflight.sh --guard-b --fixture .cdd/unreleased/{N}/REVIEW-REQUEST.yml`) → the dispatch wake's own future real invocation (not exercised by this cycle's CI, since this cycle does not itself transition any cell's `status:review` — it ships the gate, it is not gated by it; the immediately-following cycle that requests `status:review` for *this* issue, #532, will be the first live exercise of Guard B's live mode).
- `delta/SKILL.md` §9.5 (prose instruction) → δ's own future real invocation, same caller-path as above from δ's role perspective.

No dead helper functions remain in the script — I found and removed one (`yaml_block_nonempty`) during authoring after confirming via `grep -n "yaml_block_nonempty"` that it had zero callers once I switched the `changed_files` check to a direct nested-block scan.

### Test-assertion count from runner output

Pasted directly from runner output, not manually enumerated:

- Go: `go test ./... -v 2>&1 | grep -c "^--- PASS"` → **262** individual test cases pass, 0 fail, across 14 packages.
- Guard B fixture suite: `./scripts/ci/check-review-request-preflight.sh --self-test 2>&1 | grep -c "✓\|✗"` → **10** fixture assertions (2 valid pass + 7 invalid fail + 1 missing-file-repro fail), all ✓ (matched expectation), 0 ✗.
- Guard A: 1 assertion (presence-of-contract check across 4 files × ~5-8 string patterns each, all-or-nothing pass/fail per invocation) — passes.
- `cn cdd verify --unreleased`: **106 passed, 0 failed, 78 warnings** (184 total) per its own summary line, against all unreleased cycles repo-wide (not specific to #532, since #532's own self-coherence.md was incomplete at the time the gate was last run — re-run after this section completes, see §Review-readiness).
- `cn-cdd-verify` test-fixtures.sh: 3/3 fixture tests pass ("🎉 All fixture tests passed!").

### Artifact enumeration matches diff

Every file in `git diff --stat origin/main..HEAD` is named explicitly in this document's §ACs or §CDD Trace (CDD Trace step 6, next section) — cross-checked: 23 files changed total (1 pre-existing γ-scaffold commit's file + 22 files this cycle's α commits touch, now growing with this self-coherence.md itself). No file appears in the diff without a corresponding mention.

## Debt

1. **`cn cdd verify` (I6, Go) does not know about `REVIEW-REQUEST.yml`.** Per δ's pinned contract ("No Go source changes"), I6's existing artifact-presence checker was not extended to require or validate `REVIEW-REQUEST.yml` as part of its own canonical artifact set. This means I6 alone cannot detect a cell that reached `status:review` without a `REVIEW-REQUEST.yml` — that detection is Guard B's job, invoked separately by the dispatch wake. This is consistent with the issue's own deferred items ("full projection schema integration," "domain-specific CDS review-request extensions beyond the minimum proof fields") and δ's explicit Go-exclusion, so it is not a gap in *this* cycle's closure, but it is worth naming: a future cycle could fold `REVIEW-REQUEST.yml` presence into I6 for defense-in-depth (two independent checkers instead of one).

2. **Guard B's live mode is unexercised end-to-end in this cycle.** This cycle ships the gate; it does not itself request `status:review` for issue #532 through the gate (that happens after β converges, per δ's normal cycle flow — δ will run Guard B live against this very cycle's `REVIEW-REQUEST.yml` before applying `status:review` to #532 itself, per the now-updated `delta/SKILL.md` §9.5). The fixture suite (AC5) proves Guard B's structural logic deterministically; the live cross-check paths (`git rev-parse HEAD` comparison, `gh pr view`) were authored and reasoned about but not exercised against a real open PR during this implementation session, since no PR exists yet for `cycle/532` (α does not open PRs — that is δ's job per the dispatch instruction). This is expected, not a defect: the live mode literally cannot be exercised before a PR exists.

3. **`yaml_scalar`'s field-extraction is line-oriented, not a real YAML parser.** It correctly handles the flat two-level-nesting shape `REVIEW-REQUEST.yml` actually uses (confirmed by the full fixture suite, including edge cases like empty lists `changed_files: []` and empty maps `artifacts: {}` in the #524 W4 repro fixture), but it would not correctly parse deeply nested, multi-line-string, or anchor/alias YAML. This is an intentional, named tradeoff per δ's explicit instruction to keep Guard B's extraction simple (grep/sed/awk) rather than introduce a new YAML-parsing runtime dependency. If `REVIEW-REQUEST.yml`'s shape grows more complex in a future cycle, this tradeoff should be revisited.

4. **No CUE schema for `REVIEW-REQUEST.yml`.** γ's scaffold flagged this as a genuinely open question ("whether this needs a CUE schema... or whether a doc-level shape + shell-script field checks is sufficient for this cycle"). I resolved it in favor of doc-level shape + shell-script checks, consistent with δ's pinned contract's "JSON/wire contract preservation" row ("does not modify any existing `schemas/cdd/*.cue` or `schemas/cds/*.cue` contract" — read as: do not touch existing CUE files; silent on whether to add a *new* one) and the "Runtime dependencies" row's preference for simplicity. A future cycle could add `schemas/cdd/review_request.cue` if stronger typed validation becomes warranted; named here as deferred, not silently dropped.

5. **Fixture artifact stub files are minimal placeholders.** `scripts/ci/fixtures/review-request/_artifacts-524/{alpha-closeout,self-coherence,receipt}.md` exist only so Guard B's on-disk-existence check has real paths to resolve in the offline fixture suite — they are not meant to be read as documentation and contain only a one-line provenance comment each. This is intentional and named, not an oversight.

## CDD Trace

1. **contract** — issue [#532](https://github.com/usurobor/cnos/issues/532), `design-and-build` mode, scoped by γ's `gamma-scaffold.md` (commit `4404611`, pre-existing on the branch at α dispatch time), sharpened by the operator's 2026-06-30T19:29:09Z clarification comment, and pinned by δ's dispatch (Guard A/Guard B two-script split; Guard B invocation point resolved to wake-internal; implementation-contract table: Markdown+Bash+YAML, no new `cn` subcommand, doctrine at canonical paths + `scripts/ci/` + new fixture convention, no new runtime deps, additive/backward-compatible, no label taxonomy change).

2. **design** — not a separate artifact; the `REVIEW-REQUEST.yml` shape was already drafted in the issue body as a concrete YAML example (γ confirmed: "already drafted in the issue as a strong proposal, not yet a converged schema doc"). I treated documenting that shape in `dispatch-protocol/SKILL.md` (AC2) as the design artifact for this `design-and-build` cycle, per `alpha/SKILL.md` §2.2's allowance ("design... may be marked 'not required' only with a concrete justification"). Justification: the shape was already concretely specified by the issue; the work was to document and mechanically enforce it, not invent it.

3. **plan** — not a separate artifact; γ's scaffold's "Surfaces α is expected to touch" + "Per-AC oracle approach" + "Expected diff scope" sections functioned as the plan, and I followed that sequencing (doctrine surfaces first, then the guard script, then fixtures, then CI wiring, then golden re-render, then peer-surface sweep, then self-coherence) without needing a separate `PLAN.md`.

4. **tests** — the Guard B fixture suite (`scripts/ci/fixtures/review-request/{valid,invalid}/`, 9 fixtures) IS the test suite for this cycle's only new executable logic (the guard script). No Go tests were added (no Go code changed, per δ's pinned contract). The fixture suite is itself wired into CI (`build.yml` job `review-request-preflight`), so it is not merely a local convenience — it runs on every future PR.

5. **code** — `scripts/ci/check-review-request-preflight.sh` (Guard A + Guard B + self-test harness, 420 lines). No other executable code changed.

6. **docs** (full file enumeration — every file in `git diff --stat origin/main..HEAD` is named below, exhaustively, per `alpha/SKILL.md` §2.6 row 11):
   - `.cdd/unreleased/532/gamma-scaffold.md` — γ's pre-existing scaffold (not authored by α; read as input).
   - `.github/workflows/build.yml` — new job `review-request-preflight` (Guard A + Guard B self-test).
   - `.github/workflows/cnos-cds-dispatch.yml` — regenerated (not hand-edited) via `cn-install-wake cds-dispatch --out ...`, reflecting `prompt.md`'s step-7 edit.
   - `scripts/ci/check-review-request-preflight.sh` — new guard script (Guard A + Guard B + self-test).
   - `scripts/ci/fixtures/review-request/_artifacts-524/{alpha-closeout,self-coherence,receipt}.md` — fixture stub artifacts (on-disk targets for the valid fixture's `artifacts.*` paths).
   - `scripts/ci/fixtures/review-request/invalid/invalid-524-w4-empty-review.yaml` — #524 W4 reproduction fixture (AC5).
   - `scripts/ci/fixtures/review-request/invalid/invalid-empty-diff-no-no_op.yaml` — undeclared empty diff fixture.
   - `scripts/ci/fixtures/review-request/invalid/invalid-missing-artifact.yaml` — missing on-disk artifact fixture.
   - `scripts/ci/fixtures/review-request/invalid/invalid-missing-pr.yaml` — missing PR fixture.
   - `scripts/ci/fixtures/review-request/invalid/invalid-missing-request-fields.yaml` — missing `request.requested_by`/`requested_at` fixture.
   - `scripts/ci/fixtures/review-request/invalid/invalid-no-commits-beyond-base.yaml` — `base_sha == head_sha` fixture.
   - `scripts/ci/fixtures/review-request/invalid/invalid-no_op-without-approval.yaml` — undeclared-no_op-without-approval fixture.
   - `scripts/ci/fixtures/review-request/valid/valid-no_op-with-approval.yaml` — properly-declared no-op fixture (proves the exemption path actually works, not just that its absence fails).
   - `scripts/ci/fixtures/review-request/valid/valid-review-request.yaml` — full deliverable-proof valid fixture.
   - `src/packages/cnos.cdd/skills/cdd/CDD.md` — AC1, kernel doctrine addition.
   - `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` — role-skill peer update (§9.5/§9.6 Guard B precondition on the `status:review` token write).
   - `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` — AC3/AC6, step 7 + lifecycle-transitions table.
   - `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` — regenerated via `cn-install-wake cds-dispatch` (not hand-edited).
   - `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` — AC3/AC6, identical edit to SKILL.md (confirmed byte-identical at the edited region).
   - `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` — peer-surface consistency update (`responsibilities` array prose; confirmed not substituted into rendered output, so no golden-drift risk).
   - `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` — AC2/AC3/AC7, new §"Review-request proof gate" + `status:review` row precondition.
   - `.cdd/unreleased/532/self-coherence.md` — this document.

7. **self-coherence** — this document, written incrementally per §2.5's discipline (one section per commit: §Gap `4c2c283`, §Skills `337b438`, §ACs `edbff6e`, §Self-check `06fce26`, §Debt `8ee349c`, §CDD Trace this commit).
