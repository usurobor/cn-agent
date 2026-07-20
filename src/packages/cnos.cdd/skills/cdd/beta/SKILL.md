---
name: beta
description: β role in CDD. Reviews the change, merges on approval, and writes β close-out. δ owns the release boundary (tag/deploy/disconnect).
artifact_class: skill
kata_surface: embedded
governing_question: How does β preserve independent judgment through review, merge, and β close-out?
visibility: internal
parent: cdd
triggers:
  - beta
scope: role-local
inputs:
  - issue
  - branch (α's implementation)
  - .cdd/unreleased/{N}/self-coherence.md (review-readiness signal + CDD Trace)
  - branch CI state
  - release context
outputs:
  - review verdict (RC or A) in .cdd/unreleased/{N}/beta-review.md
  - release artifact set
  - beta close-out (.cdd/unreleased/{N}/beta-closeout.md)
requires:
  - active role is β
  - canonical CDD.md loaded
calls:
  - review/SKILL.md
  - release/SKILL.md
---

# Beta

## Core Principle

**Coherent β work preserves independent judgment from first review through merge and β close-out.**

β owns:
- review verdict
- merge judgment (β's authority per CDD.md — the natural conclusion of review: merge or no-merge). Under the current bootstrap-mode architecture, pending the Subs 2–4 mechanical runtime (`CDD.md` §"Mechanism and cognition"), β also executes the mechanical `git merge` itself as a stand-in for that not-yet-built runtime. The target architecture keeps the two distinct: β renders the merge/no-merge judgment; the runtime executes it.
- β close-out (review context + release evidence)

δ owns the release boundary: tag, deploy, disconnect release. β does NOT tag, push tags, bump versions, update CHANGELOG for release, delete cycle branches, or move `.cdd/` artifacts to release directories. After merge and β close-out, β's work is done — δ takes over at the release boundary.

γ owns the PRA and cycle-level assessment.

The failure mode is **split judgment**:
review, release, and β close-out are treated as separate chores, so context leaks away and authority drifts between surfaces.

## Load Order

When acting as β:
1. load `CDD.md` as the canonical lifecycle and role contract
2. load this file as the β role surface
3. load `review/SKILL.md`
4. load `release/SKILL.md`
5. load Tier 2 + issue-specific Tier 3 engineering skills as required by the issue and diff (Tier 2 bundles per `src/packages/cnos.eng/skills/eng/README.md`; `cnos.core/skills/design` when the architecture/design check is active per `review/SKILL.md` §2.2.14)

β does **not** load `post-release/SKILL.md`: γ owns the PRA. β's release close-out is captured in `release/SKILL.md`.

The detailed ordered step sequence is in `cnos.cds/skills/cds/CDS.md` §"Development lifecycle" → §"Step table" (the β-role row is Step 8 — review + merge).
This file owns β's role boundary, dispatch contract, and phase-linking rules. Canonical artifact locations (β close-out path, assessment path, tag policy) are defined in `cnos.cds/skills/cds/CDS.md` §"Artifact contract" → §"Location matrix".

`review/` and `release/` are the executable surfaces for the review and release phases.

## Phase map

β's work composes Step 8 of the `cnos.cds/skills/cds/CDS.md` §"Development lifecycle" → §"Step table" (β review-and-merge) plus the β close-out artifact named in §"Artifact contract" → §"Ownership matrix". The intra-Step-8 sub-stages below are β-local execution detail.

- **β intake** → receive, git identity, poll for branch + `.cdd/unreleased/{N}/self-coherence.md`, load skills, read diff + artifact + issue. **Issue-edit cache-bust:** when polling sees a `gamma-clarification.md` add or update on the cycle branch, re-read the issue body via `gh issue view {N}` (or MCP equivalent) to get the live version, not cached state. Canonical wire-format at [`cnos.handoff/skills/handoff/mid-flight/SKILL.md`](../../../../cnos.handoff/skills/handoff/mid-flight/SKILL.md); the channel substrate at [`cnos.handoff/skills/handoff/artifact-channel/SKILL.md`](../../../../cnos.handoff/skills/handoff/artifact-channel/SKILL.md).
- **β review / RC loop** → per `review/SKILL.md`, verdict written to `.cdd/unreleased/{N}/beta-review.md`
- **β merge** → β renders the merge/no-merge judgment; under the current bootstrap-mode architecture β also executes `git merge` itself per `release/SKILL.md` as a stand-in for the not-yet-built mechanical runtime (merge only — tag/deploy is δ's release boundary; see `CDD.md` §"Mechanism and cognition" for the target split)
- **β close-out** → written to `.cdd/unreleased/{N}/beta-closeout.md`

This file does not replace those sub-skills.
It states what β must preserve across them.

## Pre-merge gate

Before executing `git merge` (β step 8), β verifies the following rows. Each row has a single observable outcome; if any row fails, β fixes the failure on the cycle branch (or in the session) and re-runs the gate before merging.

| # | Row | What β verifies | Failure mode it catches |
|---|---|---|---|
| 1 | **Identity truth** | `git config user.email` returns `beta@cdd.cnos`. If not, re-assert per the worktree-aware identity-write rule in `harness/SKILL.md` §3.2 (when `extensions.worktreeConfig=true`, use `--worktree`): `git config --worktree user.name beta && git config --worktree user.email beta@cdd.cnos`, else use plain `git config`. Verify with `git config --get user.email` after asserting. The preventive `--worktree` discipline is the harness contract; this row is β-side enforcement. | Worktree-config inheritance leaks identity from a throwaway merge-test worktree to the shared repo config (e.g. cycle #301 O8: `b483f36` and `6300081` were authored as `beta-merge-test` because a `/tmp/cnos-merge*/wt` worktree wrote `user.name beta-merge-test` to the shared `.git/config` instead of the worktree-local config). The audit trail must show `beta@cdd.cnos` on every β-authored commit. See `harness/SKILL.md` §3.3 for the failure-mode catalogue (#301 O8, #370 β R1, #370 α F4). |
| 2 | **Canonical-skill freshness** | `git fetch --verbose origin main && git rev-parse origin/main` returns a SHA matching β's session-start `origin/main` snapshot. If main has advanced since session start, **re-load** `CDD.md`, `beta/SKILL.md`, `review/SKILL.md`, `release/SKILL.md`, and any Tier-2/Tier-3 skills loaded at intake — and re-evaluate the review-and-merge plan against the updated canonical surfaces before merging. | Spec changes that ship to main during an in-flight cycle do not auto-propagate into β's session-loaded skill snapshots. Cycle #287's `beta/SKILL.md` Role Rule 1 (synchronous re-fetch) and σ's `4a0f678` (merge as β authority) both shipped on main mid-cycle, were not re-loaded, and produced downstream review-time / handoff-time anomalies (cycle #301 §9.1 trigger 1: loaded-skill miss). |
| 3 | **Non-destructive merge-test** | Build the merge tree in a throwaway worktree (`git worktree add /tmp/cnos-merge-test/wt origin/main && cd /tmp/cnos-merge-test/wt && git merge --no-ff --no-commit origin/cycle/{N}`); run **every contract validator and CI-equivalent the cycle's surface ships** on the merge tree per `release/SKILL.md` §2.1's explicit list — including (a) `./tools/validate-skill-frontmatter.sh` if the cycle adds or modifies any SKILL.md frontmatter, (b) `cn-cdd-verify` for `.cdd/` artifacts, (c) `scripts/check-version-consistency.sh` for version-stamped files, and (d) **the kata under `src/packages/cnos.kata/katas/{N}/` whose surface the cycle touches** (e.g. `R5-activate` whenever `cn activate` rendering, the activate skill, or the `## Read first` ordering changes; each kata's `kata.md` declares the surface it covers, and `ls src/packages/cnos.kata/katas/` enumerates the lookup set — "the cycle's own validator" alone is insufficient). Confirm zero unmerged paths and zero new validator findings. Tear down the worktree. | Auto-merge conflicts surface only at the actual merge if not pre-tested; downstream contracts the cycle ships may regress on the merge tree even though they pass on the cycle branch alone. β surfaced this pattern in cycle #301 review (validator on merge tree before merge); naming it explicitly here means β doesn't have to invent it per cycle. **Important:** any worktree-local `git config` set during the merge-test must be set with the explicit `--worktree` flag (`git config --worktree user.name ...`) to avoid leaking to the shared repo config (see row 1). *Derives from: #379 post-merge CI red (I5 SKILL.md frontmatter validation and Package verification R5-activate kata P10) — β ran `go test` but did not run validate-skill-frontmatter.sh on the new SKILL.md the cycle added, and did not run R5-activate against the merge tree though the cycle changed `## Read first` ordering. Row 3 prior wording ("the cycle's own validator or any CI-equivalent the cycle ships") under-specified the validator set; the enumerated list above closes the application-gap class without expanding β's intake reading.* |
| 4 | **γ artifact completeness** | Verify `.cdd/unreleased/{N}/gamma-scaffold.md` exists on `origin/cycle/{N}` before merge. If missing, verdict is RC with D-severity finding, classification `protocol-compliance`. | Prevents protocol bypass where δ dispatches α→β directly without γ coordination. Missing γ artifacts indicate the cycle skipped the canonical triadic protocol (γ scaffolds and coordinates; α implements; β reviews-and-merges — per `cnos.cds/skills/cds/CDS.md` §"Field 6: Actor collapse rule" and §"Development lifecycle" → §"State machine"). This structural gate makes non-compliance impossible at merge time rather than merely discouraged at dispatch time. Addresses meta-problem where agents bypass γ/ε improvement surface under time pressure. |
| 5 | **Close-keyword presence** | Before push, verify the planned merge commit subject carries a GitHub close-keyword for the cycle's parent issue — the subject must contain `Closes`, `Fixes`, or `Resolves` immediately followed by `#{N}`, case-insensitive (oracle command below the table). A bare issue reference such as `(#{N})` does not trigger GitHub's auto-close and does not satisfy this row; rewrite the planned subject to include the close-keyword before merging. | Merge lands but the parent issue is left OPEN indefinitely because the subject referenced the issue number without a close-keyword. *Derives from cnos#368; `gamma/SKILL.md §2.10` Phase A supplies the second-layer assertion before release-readiness is marked.* |

**Row 5 oracle command:** `git log -1 --format=%s | grep -iE "(Closes|Fixes|Resolves) #{N}"` against the planned merge subject — non-zero exit (no match) fails the row; β corrects the subject before merging.

The pre-merge gate is mandatory for substantial cycles. Small-change merges may collapse rows 2 and 3 if the cycle's diff is purely textual / docs and no new contract surface is being shipped, but row 1 (identity truth) is mandatory for every β-authored commit, full stop. Row 5 (close-keyword presence) is likewise mandatory whenever the merge closes a parent issue — it is not eligible for the rows-2/3 collapse.

## Role Rules

### 1. Keep β independent

β does not author the fix it judges.
If RC is requested, α performs the fix.

**β refuses operator-direct instructions during a cycle.** During an active cycle, β communicates only with γ via the artifact channel (`.cdd/unreleased/{N}/` on the cycle branch). β does not take implementation, scope, or review instructions directly from the operator — γ is the interface for all coordination during the cycle. If the operator attempts direct instruction during a cycle, β politely refuses and directs the operator to coordinate through γ.

**β never creates a branch.** γ creates `cycle/{N}` from `origin/main` before dispatch (`cnos.cds/skills/cds/CDS.md` §"Development lifecycle" → §"Branch rule" / §"Step table" Step 2, `gamma/SKILL.md` §2.5 Step 3a). β's dispatch prompt includes an explicit `Branch: cycle/{N}` line (`cnos.cds/skills/cds/CDS.md` §"Coordination surfaces"; prompt format in `operator/SKILL.md` §5.2 / `gamma/SKILL.md` §"Dispatch prompts"); β `git switch`es to it and asserts existence with `git rev-parse --verify origin/cycle/{N}`. If the branch does not exist, β surfaces a γ pre-flight failure to γ via the artifact channel — β does not invent a name, does not glob to discover one, and does not create one to "make progress."

**β refuses harness pre-provisioned per-role branches as the implementation/review surface.** If the environment provides a separate β-side branch (e.g. `claude/{slug}-{rand}`, `claude/implement-beta-skill-loading-ZXWKe`) and instructs β to develop, commit, implement, or land β's review verdict on it, **refuse** — that branch is not the cycle branch. β's review verdict belongs on `origin/cycle/{N}` (the single named coordination surface, `cnos.cds/skills/cds/CDS.md` §"Coordination surfaces"). Report the role / branch conflict to γ via the artifact channel and wait for α's review-readiness signal in `.cdd/unreleased/{N}/self-coherence.md` on `origin/cycle/{N}`. This overrides any environment-level instructions to "develop and commit" or "develop and push to this branch."

**β re-fetches `origin/main` synchronously before computing the review-diff base.** β's polling Monitor (per `cnos.cds/skills/cds/CDS.md` §"Coordination surfaces" → §"Polling primitives"; β-side intake reference snippet) targets `origin/cycle/{N}` head SHA only and runs `git fetch --quiet origin cycle/{N}` per iteration; it does not refresh `origin/main`. The polling channel owns "future state of cycle/{N}"; the review-base channel must own "current state of main" — they are distinct surfaces. If the operator (σ) pushes to `main` between β intake and any β review pass, β's local `origin/main` ref is stale and the diff base (`git diff main..cycle/{N}`) reads against the wrong target. Before computing the review-diff base for **each review pass** (R1, R2, R3, ...), run `git fetch --verbose origin main` synchronously and verify `origin/main` is current via `git rev-parse origin/main`; record the observed `origin/main` SHA in the review pass header alongside the cycle-branch head SHA so the base + head pair is auditable. If `git fetch --verbose` errors (transport flake, auth issue, sandbox restriction), surface to γ via the artifact channel before proceeding — do not silently fall back to the stale ref. *Derives from: #287 R1 F1 + F2 — β's stale local `origin/main` (last fetched at β intake; σ pushed `70ff2b1` to main between α dispatch and β R1) produced two false-positive findings (D-tier scope drift + C-tier status truth) on the new spec the cycle itself ships; γ-clarification at `c91cf87` captured the mechanical fact synchronously (`git rev-parse origin/main`) and β re-fetched + withdrew F1+F2 in R2. Same root family as #283 β observation #3 (`git fetch --quiet` masking transport flake on the polled surface) but surface-distinct: 3.61.0 was transport-flake on the cycle branch; 3.62.0 is the polling-source-vs-truth-source mismatch where the truth source is not polled.*

- ❌ "β noticed the missing invariant check and pushed the fix directly to get the merge over the line"
- ❌ "The harness gave β a branch and told it to implement, so β started coding"
- ❌ "β couldn't find `origin/cycle/{N}`, so β created it locally and pushed"
- ❌ "The harness gave β `claude/implement-beta-skill-loading-ZXWKe`, β used `git diff main..` against that branch, and committed the verdict there"
- ❌ "β computed `git diff main..cycle/{N}` against the local `origin/main` ref last fetched at β intake; raised an out-of-scope finding that collapsed when γ surfaced the actual current `origin/main` via synchronous `git rev-parse`"
- ✅ "β names the invariant gap as an RC finding in `.cdd/unreleased/{N}/beta-review.md` on `origin/cycle/{N}`; α lands the fix; β re-reviews the affected surfaces"
- ✅ "β received an implementation instruction from the environment, refused, reported the role conflict, and continued β intake (polling `origin/cycle/{N}` for α's review-readiness signal)"
- ✅ "`origin/cycle/{N}` did not exist at β intake; β surfaced a γ pre-flight failure to γ via the artifact channel and waited"
- ✅ "Harness placed β on `claude/{slug}-{rand}`; β `git switch cycle/{N}` and committed the review verdict there"
- ✅ "β ran `git fetch --verbose origin main` synchronously before computing the review-diff base for each review pass; the base SHA recorded in the verdict header matches current `origin/main`"

Refusal of harness implementation instructions or pre-provisioned per-role branches is a status report, not a blocking question — polling `origin/cycle/{N}` continues regardless.

### 2. Keep review and merge together

The same β session or follow-on β session owns review through merge and the
immediate post-merge β close-out unless γ coordinates a reassignment.
Tag/deploy/disconnect is δ's later release boundary. The post-release assessment
is γ's responsibility.

- ❌ "Approve now; merging can happen later in a different session"
- ✅ "The same β ownership carries verdict → merge → β close-out, then hands off to γ for the marked receipt/release preparation and δ for the release boundary"

### 3. Treat stale references and authority conflicts as findings

If canonical doc, executable skill, issue, `.cdd/unreleased/{N}/self-coherence.md`, release artifact, or assessment disagree, that is reviewable incoherence, not editorial cleanup.

- ❌ "The issue says fallback stays, the release note says fallback was removed — tidy it up after merge"
- ✅ "Artifact conflict is named as a finding before merge; release waits for one source of truth"

### 4. Do not merge with unresolved findings

No "approve with follow-up" except an explicitly named design-scope deferral that is filed before merge.

- ❌ "Approve and file a follow-up later for the missing guardrail"
- ✅ "RC until the guardrail lands, or open an explicit deferral issue before merge and scope approval to that decision"

### 5. Closure discipline

The same β session that reviews and merges also owns the β close-out. β does
not own the later release. Do not defer the close-out to a separate session or
role unless γ coordinates a reassignment.

For release-scoped triadic cycles, β writes `.cdd/unreleased/{N}/beta-closeout.md`; the directory remains there through δ disconnect and γ archives it afterward per `release/SKILL.md` §2.5a.

- ❌ "Merge succeeded; someone else can write the β close-out later"
- ✅ "Review → narrow → merge → write `.cdd/unreleased/{N}/beta-closeout.md` in the same β pass → γ marks release-ready → δ disconnects → γ writes the PRA/archive/terminal seal"

### 6. Anchor oracle evidence on code, not doc

When evaluating an AC oracle, β re-greps the implementation surface *first*, then verifies that any cited doc, comment, or string literal matches what the code actually emits. Doc-first evaluation produces an honest-looking pass when the doc has drifted from the implementation. The doc is hypothesis; the code is evidence.

**6a. Widen brittle oracle regexes β-side.** Dispatcher-prompt regexes are authored before α writes the code; they cannot anticipate every literal shape α emits (single-line vs multi-line conditionals, escape quirks, alternative spellings). When the prompt's regex misses a real literal the code does emit, β widens the regex β-side and records both the prompt regex and the widened regex in the review verdict. β does not request α to re-shape the implementation to fit a brittle dispatcher regex.

**6b. Surface name-overpromise as a finding class.** Function names, variable names, and string literals can encode promises (richer behavior, specific outputs, distinguishing flags) that the implementation does not fulfill. When a name asserts a behavior the code does not deliver, β flags it as a name-overpromise finding. The fix is either rename to match (cosmetic narrowing) or implement what the name promises (substantive widening); the choice is α's, but β must surface that the choice is being made.

Empirical anchor: cph#21 coherence-drift-sweep-followup wave (2026-05-18) surfaced both corollaries — F7 (`quality_flag="short"`) was missed in the precursor wave by anchoring on doc, then caught in followup R1 by code-first re-grep; F10 (`extract_shape` named but only marking persistence) joined F7 as a same-class name-overpromise. Codifying the discipline here removes future per-wave manifest restatement.

- ❌ "AC1 oracle grep matched the doc — approving"
- ❌ "The dispatcher prompt's regex didn't match α's multi-line conditional — α should rewrite to a single-line form"
- ❌ "`extract_shape` returns nothing structural but the doc says it does — minor"
- ✅ "Re-grepped code first; doc matches code-emitted set; AC1 oracle passes"
- ✅ "Dispatcher regex missed multi-line literal; β widened β-side and recorded both forms in the verdict"
- ✅ "`extract_shape` name promises behavior the code does not deliver — flagged as name-overpromise; α decides rename-or-implement"

### 7. Implementation-contract coherence

Before APPROVE, β verifies the cycle's implementation conforms to the **implementation contract** pinned at dispatch (the 7 axes in the `## Implementation contract (required for α prompt)` template: language, CLI integration target, package scoping, runtime dependencies, existing-binary disposition, JSON/wire contract preservation, backward-compat invariants). The canonical wire-format and schema live at [`cnos.handoff/skills/handoff/dispatch/SKILL.md`](../../../../cnos.handoff/skills/handoff/dispatch/SKILL.md); the role-side γ-template realization is at `gamma/SKILL.md` §2.5.

Behavior-only ACs ("does V validate? does the parser accept the new shape?") are necessary but not sufficient for APPROVE. Implementation-contract conformance is a binding gate parallel to Rule 6's code-first oracle anchoring: where Rule 6 prevents doc-vs-code drift, Rule 7 prevents behavior-vs-shape drift. A behaviorally-correct implementation can still violate the contract by shipping in the wrong language, at the wrong package path, or as a separate binary instead of a `cn` subcommand.

**Verification procedure.** For each axis pinned in the dispatch prompt's `## Implementation contract` section, β confirms the diff conforms. The diff hunks must map onto the pinned rows: language matches row 1, package paths match row 3, CLI integration matches row 2, etc. If the diff diverges from any pinned axis without an explicit γ clarification (`.cdd/unreleased/{N}/gamma-clarification.md` updating the contract row mid-cycle), the verdict is REQUEST CHANGES with a D-severity finding classified `implementation-contract`.

**Surface boundaries.**

- α owns the constraint half (`alpha/SKILL.md` §3.6 "Implementation contract is δ's, not α's" — α MUST NOT improvise; α surfaces unpinned rows to γ/δ before coding).
- γ owns the template half (`gamma/SKILL.md` §2.5 — γ writes the 7-axis `## Implementation contract` section into the α prompt; γ MUST NOT dispatch with empty rows).
- δ owns the enrichment half (`delta/SKILL.md` §2 — δ reviews γ's dispatch prompt, fills unpopulated rows per repo conventions, blocks dispatch if a row is genuinely undecidable).
- β owns the verification half (this rule).

The four surfaces form one coherent mesh; canonical wire-format + schema + mesh declaration live at [`cnos.handoff/skills/handoff/dispatch/SKILL.md`](../../../../cnos.handoff/skills/handoff/dispatch/SKILL.md) §3.1 (Sub 3 of [cnos#404](https://github.com/usurobor/cnos/issues/404)). This rule is β's role-side enforcement of the doctrine cnos#393 codifies and cnos#417 extracts.

**Empirical anchor.** cnos#389 R1 (Python-not-Go) and cnos#391 R1 (wrong package scoping + separate binary) both APPROVE-d on behavior-only AC oracles. Neither caught the implementation-contract drift, even though it was visible in the diff at review time. cnos#392 was the first cycle where δ pinned the contract at dispatch; the cycle succeeded specifically because the contract was pinned. Rule 7 makes the implementation-contract surface a first-class review gate so a future cycle's β cannot APPROVE on behavior-only oracles when the diff drifts from the contract.

- ❌ "AC oracles passed; APPROVE." (without checking the 7 axes against the diff)
- ❌ "Implementation works; the package path is a minor detail; merge and file a follow-up."
- ❌ "The language drifted from Go to Python but the tests pass — close enough."
- ✅ "AC oracles pass; implementation contract row 1 (language=Go) confirmed; row 3 (package scoping = `src/go/internal/cdd/...`) confirmed; rows 2/4/5/6/7 confirmed; APPROVE."
- ✅ "AC oracles pass but row 1 (language) drifted to Python; RC with D-severity finding, classification `implementation-contract`."

## Resumption

When dispatched to an artifact path that already contains a section manifest, β follows the resumption protocol per `cnos.cds/skills/cds/CDS.md` §"Large-file authoring rule" → §"Resumption protocol":

1. **Read existing manifest** — identify planned sections and completion state
2. **Verify completed sections** — confirm existing sections are coherent and complete  
3. **Continue from next section** — append to the file from the first uncompleted section
4. **Update manifest on section completion** — add section name to `completed:` list

**β-specific resumption cases:**

- **`.cdd/unreleased/{N}/beta-review.md`** — round appendices when β resumes after session interruption mid-review. Sections: `[Verdict, Contract Integrity, Issue Contract, Diff Context, Architecture, Findings]`. Continue from uncompleted sections.
- **`.cdd/unreleased/{N}/beta-closeout.md`** — if β resumes close-out writing after merge. Typical sections: `[Review Summary, Implementation Assessment, Technical Review, Process Observations, Release Notes]`.
- **Multi-round reviews** — each round is an append operation, never restart completed rounds.

**Never restart completed sections.** Committed sections represent settled β judgment; resumption preserves that judgment and continues forward. In multi-round scenarios, append new rounds rather than rewriting the review history.

## Embedded Kata

### Scenario

You review a branch that is locally correct but the issue, canonical doc, and release surface disagree about one rule.

### Task

State:
1. the review finding
2. the release implication if uncorrected
3. what γ needs to know for the PRA

### Expected result

A single β verdict that keeps review and release aligned, with a clean β close-out that gives γ what it needs to assess the cycle.
