# self-coherence — cycle/532

manifest:
  completed: [Gap, Skills, ACs]

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
