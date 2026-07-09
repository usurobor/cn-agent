# self-coherence.md — cnos#626 (α, R0)

## Gap

cnos#626 — "arch(cds): isolate the dispatch cell-worker from substrate/agent
identity + hub state." Per `.cdd/unreleased/626/gamma-scaffold.md` this is a
**supervised design-first cell** (operator-narrowed, dated 2026-07-08):
implement ONLY the bounded/safe mechanism (γ's "Recommended design" move
(3) — drop full six-step sigma activation for cell-execution cognition,
load Kernel + CA skills only); write doctrine + plan for the unbounded
part (AC3/AC4 — capability removal / write-fence retirement) and STOP.

Mode: doctrine → implementation (partial, per the bounded/unbounded split
γ already reasoned through). Base SHA: `86042ec5be4b5fb45b213c27dfcf635958f60aac`
(cycle/626 created from this `main` HEAD). Branch: `cycle/626`.

## Skills

- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role contract (this
  round's execution surface: §2.1 dispatch intake, §2.2 produce-in-order,
  §2.5 self-coherence incremental-write discipline, §2.6 pre-review gate,
  §3.6 implementation-contract constraint).
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9 (dispatch-wake-invoked
  mode), specifically §9.2 (five-input contract), §9.10/§9.11 (sibling
  resume-shape subsections whose heading/prose style the new §9.12 follows).
- `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` (the wake
  prompt this cell edits) — read in full before editing, including
  "Claim mechanism," "Disallowed surfaces," "Lifecycle transitions," and the
  write-fence description in §"Disallowed surfaces" (left untouched by
  design).
- `src/packages/cnos.core/skills/agent/activate/SKILL.md` §2.1 (six-item
  load order; the "soul vs identity" layering rule quoted verbatim in the
  new doctrine prose) and §2.5 (containerized-path / `.cn-{agent}/`
  detection) — read-only cross-reference, NOT edited.
- The issue body for cnos#626 itself, including the "Operator directive —
  run #626 as a supervised design-first cell" and "Authority — A2/CAP
  (narrowed for a design-first cell)" sections (read via `gh issue view 626`
  before starting; both sections match the scaffold's restatement of them
  and impose no additional constraint beyond what the scaffold already
  encodes).
- No Tier-2/Tier-3 engineering skills loaded — this cell is a pure
  Markdown/prose skill-file edit (per the implementation contract's
  "Language: N/A ... Markdown/prose skill-file edits only" row); no
  code-authoring skill applies.

## ACs

Oracle list per `gamma-scaffold.md` §"Per-AC oracle list."

- **AC1** — "A doctrine statement exists, on `main` after merge, naming:
  (a) the cell carries no substrate identity and no hub-state scope, (b)
  the wake owns identity + channel logs." **MET.** New
  `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.12 "cell/substrate
  identity boundary (cnos#626)" states both halves explicitly in its
  opening bold sentence: *"the cell carries no substrate identity and no
  hub-state scope. The wake/substrate owns identity, channel logs, and
  reporting; δ and the roles it dispatches (γ/α/β) do not."* The
  "Who owns what" paragraph names the wake/substrate's ownership of
  substrate identity, `.cn-{agent}/logs/`, and operator-facing reporting
  explicitly, and cross-references AGENT-ACTIVATION-LOG-v0 §0. Evidence:
  `git diff src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` — new §9.12
  block (19 lines in the unified diff's insertion count, sited after §9.11
  at the file's end).
- **AC2** (first disjunct — "no longer performs a full sigma activation")
  — **MET.** `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`
  §"Identity and activation" no longer contains the "Activate per
  activate/SKILL.md: Kernel → CA skills → Persona → Operator → hub state →
  identity confirmation" instruction. Oracle run:
  `grep -n "Activate per.*activate/SKILL.md" src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`
  → **0 hits** (exit code 1). The replacement paragraph states doctrine:
  cell-execution cognition loads Kernel + CA skills only; does NOT load
  Persona/Operator/hub-state; does NOT read `.cn-{agent}/` or thread
  surfaces; the wake's self-identification is manifest-declared (`{agent}`
  variable + bot-account binding), not file-discovered. The existing "You
  do NOT attach to a channel like the admin wake does" sentence is
  preserved verbatim as its own paragraph immediately after (not
  duplicated, not rewritten). AC2's second disjunct (checkout-scope
  exclusion / `.cn-{agent}/` no longer in working scope) is **NOT**
  implemented this cell — see "AC3/AC4" section below; AC2's "and/or"
  wording means the first disjunct alone is a legitimate pass (per
  gamma-scaffold.md Friction note 1).
- **AC2** (rendered-mirror regeneration requirement) — **MET.** Both
  mirrors were regenerated via the renderer, never hand-edited:
  - `./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch`
    (default `--out`, i.e. the per-package golden at
    `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`)
  - `./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml`
  Verification: `diff .github/workflows/cnos-cds-dispatch.yml src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`
  → **no output** (byte-identical). `sha256sum` of both files also
  matched (see CDD Trace below for the exact hashes) — this reproduces the
  `install-wake-golden.yml` job "Verify live cds-dispatch workflow matches
  golden (sha256)" step's own check. Idempotence also verified: a second
  invocation of the live-target render produced a byte-identical file
  (`diff` empty), matching that workflow's own "Verify idempotence"
  step's invariant.

  Note on the CLI invocation form: `gamma-scaffold.md`'s α prompt and the
  "Implementation contract" table both say `cn install-wake cds-dispatch
  --out <path>`. The actual renderer on this branch is a standalone shell
  script at `src/packages/cnos.core/commands/install-wake/cn-install-wake`
  (invoked directly, e.g. `./src/.../cn-install-wake cds-dispatch`), not a
  subcommand of the `cn` Go binary built from `src/go/cmd/cn` — `cn help`
  on that binary lists no `install-wake` verb. `install-wake-golden.yml`
  itself invokes the tool the same way (`./src/packages/cnos.core/commands/install-wake/cn-install-wake
  cds-dispatch`, see its "Re-render cds-dispatch wake" step). This is a
  naming/packaging discrepancy between the scaffold's prose and the
  on-branch tooling, not a functional gap — the correct renderer was
  invoked with the correct arguments and produced the correct output,
  verified byte-identical against both the golden and the live mirror.
  Flagged here for β/γ visibility rather than silently reconciled.
- **AC3** — "A code cell cannot see or write `.cn-*/logs/` — the
  capability is removed." **NOT implemented this cell — doctrine + plan
  only, per operator STOP condition.** See "AC3/AC4 — doctrine + plan,
  deferred to operator" section below. This is the deliberate, scaffolded
  outcome, not a dropped AC.
- **AC4** — "Once AC3 holds, the write-fence is removed as redundant."
  **NOT applicable this cell** — hard-gated on AC3, not proven this cycle.
  The write-fence code (`.github/workflows/cnos-cds-dispatch.yml`
  `dispatch_activation_log_write_violation` step, formerly lines
  ~358–439) was NOT touched. Evidence: `git diff .github/workflows/cnos-cds-dispatch.yml`
  shows exactly one hunk (the "Identity and activation" prompt-body
  paragraph swap at what was line 100); `grep -c
  dispatch_activation_log_write_violation` on both the pre-edit
  (`git show HEAD:.github/workflows/cnos-cds-dispatch.yml`) and post-edit
  file returns the same count, and a byte-diff of just that step's region
  is empty (see CDD Trace).
- **AC5** — "No regression: FSM ownership, finalizer, reconciler, dispatch
  guards, and all gates stay green." **In scope, MET at the level this
  cell can verify locally.** Zero Go source files changed
  (`git diff --name-only | grep -E '\.go$'` → no output) — the diff is
  prose/Markdown plus the two regenerated YAML mirrors, so `go test ./...`
  is a pure no-new-behavior regression check, not exercised further by
  this cell's own diff. `src/packages/cnos.cds/skills/cds/fsm/transitions.json`
  is untouched (`git diff --stat -- src/packages/cnos.cds/skills/cds/fsm/transitions.json`
  → no output). The `install-wake-golden` sha256/idempotence checks this
  cell reproduced locally (see AC2 above) are the mechanical proof the
  golden/live sync + renderer determinism invariant AC5 depends on for
  this surface still holds. `#516`/`#524` dispatch guard scripts
  (`check-dispatch-repair-preflight.sh`, `check-dispatch-closeout-integrity.sh`)
  are not referenced anywhere in this diff (`git diff --name-only` does
  not list them) — they read `cds-dispatch/SKILL.md` prose sections this
  cell did not touch (the "Repair re-entry preflight" and "Closeout
  integrity preflight" sections are untouched; only "Identity and
  activation" changed).
- **AC6** — "`cell_kind` stays observed-only; no new taxonomy; no Demo 0."
  **MET (guardrail, not a build target).** The diff introduces no new
  `cell_kind` value, no `status:*` label, no Demo-0-shaped artifact — it
  is exactly four files: two SKILL.md prose files and their two
  regenerated rendered mirrors.

## Self-check

Did α's work push ambiguity onto β? Reviewed against the β prompt's
"Independent AC walk" in `gamma-scaffold.md`:

- AC1 walk: β is asked to confirm the new §9.12 subsection "explicitly
  states (a)...(b)..." and "cross-references §9.2's five-input contract
  and explicitly names the Persona/Operator/hub-state omission as
  intentional." Both are present verbatim in §9.12's second paragraph
  ("§9.2's silence is intentional, not an oversight") — this is not an
  implied statement β has to infer; it is a stated claim with a stated
  reason ("the absence is doctrine, not debt"). No ambiguity pushed here.
- AC2 walk: β is asked to grep for "Activate per" / "Persona" / "Operator"
  / "hub state" in the "Identity and activation" section. Note: my new
  paragraph **does** contain the strings "Persona," "Operator," and "hub
  state" (it names them as things that are NOT loaded) — a naive grep-hit
  count is not zero. β's own instruction already anticipates this
  ("confirm the six-step full-activation instruction is gone or narrowed
  to Kernel+CA-skills-only for cell execution" — a hit is not
  automatically a fail, the six-step *chain* instruction has to be gone).
  I flag this explicitly here rather than let β discover a grep hit and
  wonder whether it's a miss: the words "Persona," "Operator," "hub
  state" appear only in the negation ("does NOT load"), not as an
  instruction to load them. The literal `"Activate per.*activate/SKILL.md"`
  grep β is told to run (mirroring the scaffold's own AC2 oracle) returns
  0 hits, which is the actual pass condition.
- AC2 walk: β is asked to confirm `activate/SKILL.md` itself is
  unchanged. Confirmed via `git diff --stat -- src/packages/cnos.core/skills/agent/activate/SKILL.md`
  → no output.
- AC3/AC4 walk: β is asked to confirm no checkout-isolation mechanism was
  built and the write-fence is untouched. Both hold; see AC3/AC4 evidence
  above and the dedicated section below. I did not weaken the STOP by
  writing a "partial" sparse-checkout stub or similar — the two candidate
  mechanisms are named as plan text only, with zero corresponding diff
  hunks anywhere outside the two prose files.
- The renderer invocation-path discrepancy (shell script vs. `cn`
  subcommand) is disclosed above rather than silently smoothed over,
  since a future reader diffing my commands against the scaffold's
  literal `cn install-wake ...` text could otherwise wonder whether I
  ran the wrong tool.
- Is every claim backed by evidence in the diff? Yes — every AC row above
  cites either a specific grep/diff command and its actual output, or a
  specific file/line the diff touches. No claim rests on "I read it and
  it looked right" without a corresponding shell-verifiable check.

## Debt

- **Git author identity.** This session's commits carry the pre-existing
  session git identity (`sigma@cnos.cn-sigma.cnos` /
  `41898282+sigma@cnos.cn-sigma.cnos@users.noreply.github.com`), not the
  canonical `alpha@cdd.cnos` pattern named in `alpha/SKILL.md` §2.6 row
  14. Per that row's path (b) ("configure correctly going forward and
  accept existing commits as legacy — permitted only with explicit
  known-debt disclosure"), this is disclosed as known debt rather than
  rewritten: this cell runs inside a wrapper session whose git identity
  was set by the dispatching harness before α-role work began, and
  amending authorship here was judged out of scope for a prose-only,
  narrowly-bounded design-first cell (rewriting history was not named as
  a requirement anywhere in the α prompt, and doing so unprompted risks
  a destructive-git-op the operator did not ask for).
- **AC3/AC4 are deferred, not closed** — see dedicated section below.
  This is declared debt by design (the scaffold's own framing), not an
  oversight.
- **`cn install-wake` naming discrepancy** — the scaffold and the α
  prompt both cite `cn install-wake <wake> --out <path>` as the
  invocation form; the actual on-branch tool is a standalone shell
  script, not a `cn` Go-binary subcommand. Named above under AC2 and
  here for visibility; does not affect the correctness of what was
  rendered (verified byte-identical against both mirrors + sha256).
- No test suite exists for prose/Markdown skill-file content beyond the
  `install-wake-golden` CI gate (sha256 + idempotence), which this cell
  reproduced locally rather than via a live GitHub Actions run — a live
  CI run on the pushed branch is the actual gate; local reproduction is
  strong evidence but not identical to the CI environment. No behavior
  difference is expected (the renderer is a pure function of the
  manifest + prompt-body text and the same shell script executes in
  both places), but this is named as a residual gap between "locally
  reproduced" and "CI-observed."

## CDD Trace

Step 7 (self-coherence) trace, mapped against
`cnos.cds/skills/cds/CDS.md` §"Development lifecycle" step table:

1. **Gap** — cnos#626, supervised design-first cell, bounded/unbounded
   split per `gamma-scaffold.md`. Documented above.
2. **Mode** — doctrine → implementation (partial). Documented above.
3. **Artifacts touched** (full `git diff --stat` against `origin/main`):

   ```
   .github/workflows/cnos-cds-dispatch.yml                                    |  2 +-
   src/packages/cnos.cdd/skills/cdd/delta/SKILL.md                            | 19 +++++++++++++++++++
   src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md                  |  2 +-
   src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml |  2 +-
   4 files changed, 22 insertions(+), 3 deletions(-)
   ```

   Every file above is named in an AC row above (AC1: `delta/SKILL.md`;
   AC2: `cds-dispatch/SKILL.md` + both rendered mirrors). No file in the
   diff is unmentioned.
4. **Self-coherence** — this document.
5. **Renderer verification (mechanical, standing in for CI until the
   branch's own `install-wake-golden` run is observed):**

   ```
   $ diff .github/workflows/cnos-cds-dispatch.yml \
          src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
   (no output — byte-identical)

   $ sha256sum .github/workflows/cnos-cds-dispatch.yml
   8f6897164af05bcb1fc405490c45c8721c74c5c7f301555b00ca802813e76ff1  .github/workflows/cnos-cds-dispatch.yml
   $ sha256sum src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
   8f6897164af05bcb1fc405490c45c8721c74c5c7f301555b00ca802813e76ff1  src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
   ```

   Idempotence: a second render of the live target
   (`./src/packages/cnos.core/commands/install-wake/cn-install-wake
   cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml`) produced a
   file byte-identical to the first render (`diff` empty).
6. **Guardrail confirmations** (see "What I did NOT touch" below for the
   full list with commands).
7. **Review-readiness** — see below.

## What I did NOT touch (guardrail confirmations)

Per the α prompt's guardrails and "When done" §(d):

- **Write-fence code** (`dispatch_activation_log_write_violation`,
  `.github/workflows/cnos-cds-dispatch.yml`) — untouched. The only hunk
  in that file's diff is the "Identity and activation" paragraph swap;
  `git diff .github/workflows/cnos-cds-dispatch.yml | grep -i
  "write.fence\|dispatch_activation_log_write_violation"` returns no
  hits.
- **`.cn-sigma/` contents** — untouched. `git diff --stat -- .cn-sigma/`
  returns no output.
- **The checkout step / any sparse-checkout mechanism** — not added.
  Confirmed by the same file-list check above (only 4 files changed, none
  of them a checkout-step region).
- **`activate/SKILL.md`** — untouched. `git diff --stat -- src/packages/cnos.core/skills/agent/activate/SKILL.md`
  returns no output. It was read (Key finding 2 / §2.1 / §2.5) but never
  edited — only *who is required to run it* narrows, not the procedure.
- **`src/packages/cnos.cds/skills/cds/fsm/transitions.json`** (#633's
  surface) — untouched. `git diff --stat -- src/packages/cnos.cds/skills/cds/fsm/transitions.json`
  returns no output. I found no evidence anywhere in the issue body, the
  scaffold, or the diff's own surfaces that #633's FSM rule-ordering fix
  is directly required by this cell's bounded scope — the bounded scope
  is entirely prose in two SKILL.md files plus their rendered mirrors,
  none of which the FSM consumes.
- **No Go source changed** — `git diff --name-only | grep -E '\.go$'`
  returns no output. The bounded mechanism required zero Go/infra
  changes, matching the "Implementation contract" table's "Zero Go
  source changes expected" pin; I did not need to escalate the "if you
  find you need one, STOP" guardrail.
- **No new `status:*` label, no new `cell_kind`, no FSM/transitions.json
  change, no Demo 0** — confirmed by the same 4-file diff scope; nothing
  in the diff touches label definitions, cell-kind taxonomy, or the FSM.
- **No hand-edit of either rendered mirror** — both were produced
  exclusively via the `cn-install-wake` renderer invocations documented
  under AC2 above; I did not open either YAML file in an editor at any
  point in this cycle.

## AC3/AC4 — doctrine + plan, deferred to operator

Per the α prompt's explicit requirement ("This is a REQUIRED deliverable
of this cell, not optional narrative"), restating and confirming
γ's "Recommended design" analysis from `gamma-scaffold.md` rather than
re-deriving it from scratch — the analysis was already done at scaffold
time and nothing observed during implementation changes it:

**What AC3/AC4 require.** AC3: "a code cell cannot see or write
`.cn-*/logs/` — the capability is removed (e.g. not present in the
checkout, or checkout is scoped), not merely guarded post-hoc." AC4
(hard-gated on AC3): "once AC3 holds, the `dispatch_activation_log_write_violation`
write-fence is removed as redundant."

**Why neither is implemented this cell.** The operator's own STOP
condition governs directly: *"if the safe mechanism is not clearly
bounded, STOP at doctrine + implementation plan and hand back to the
operator rather than force an implementation."* γ's scaffold evaluated
all three candidate boundary moves the issue names and found exactly one
(move 3 — drop full activation for cell-execution cognition; the AC1/AC2
work completed above) cleanly boundable inside this single cell. The
other two (candidates 1 and 2 below) are not.

**Candidate mechanism 1 — isolated checkout / sparse worktree excluding
`.cn-sigma/`.** Add a post-checkout `git sparse-checkout` step
(non-cone mode: `/*` + `!/.cn-sigma/`) to the rendered workflow, gated by
a new wake-manifest field (e.g. `wake.surfaces.excluded_checkout_paths`),
requiring `cn install-wake` renderer (Go... or, per the discrepancy noted
above, shell-script) changes.

- *Tradeoff:* This changes the checkout shape of a **live production**
  wake (`activation_state: live` since cnos#487) that other in-flight and
  future cells depend on for reading the full codebase, opening PRs, and
  running the finalizer (`cn cell finalize`, cnos#591). The operator's own
  STOP conditions name exactly this risk: "isolating the cell would break
  its ability to read the codebase or the issue contract"; sparse-checkout
  could silently break a doc cross-reference or a Tier-3 skill path
  lookup that happens to resolve through a path pattern not anticipated
  at design time. Validating this safely needs its own regression pass
  (see "What a follow-on cell would need" below) that this single cell
  cannot responsibly absorb alongside AC1/AC2.

**Candidate mechanism 2 — move `.cn-sigma/` out of the product repo
entirely.** Relocate Sigma's per-context state (the mailbox, the channel
logs) to a separate repo/location; cnos becomes pure-codebase from the
dispatch wake's perspective.

- *Tradeoff:* Explicit STOP condition in the operator directive: "Moving
  the hub out of the repo requires new infra/permissions not available →
  escalate." This is an infra/ownership decision (where does the new
  location live, who has write access, does `activate/SKILL.md`'s
  Tier 1b pure-product-hub path already cover the read side) that the
  operator, not this cell, must authorize. Note (per gamma-scaffold.md
  Friction note 5): `.cn-sigma/README.md` already documents that
  *identity* files (PERSONA.md/OPERATOR.md) live at Sigma's home hub, not
  in `.cn-sigma/` at cnos — only per-context state (mailbox, logs)
  remains in-repo by design. This nuance doesn't change the bounded/
  unbounded split; it just means candidate 2, if ever pursued, is a
  narrower move than it might first appear (per-context state only, not
  identity).

**What this cell does NOT do.** No sparse-checkout step, no new manifest
field, no relocation of any `.cn-sigma/` content, no change to the
checkout step in either rendered workflow. Confirmed under "What I did
NOT touch" above.

**What a follow-on cell would need to prove AC3 safely.** Before
attempting candidate 1 (sparse-checkout) on the live production wake, a
follow-on cell needs:

1. A **regression matrix over every existing cell type's checkout-path
   needs** — enumerate every doc cross-reference, Tier-2/Tier-3 skill
   path, and issue-body-linked artifact any live cell class has ever
   consumed from the working tree, and confirm none of them resolve
   through a path a `.cn-sigma/`-excluding sparse-checkout would hide.
   This needs to be built empirically (e.g. instrument a dry-run cell
   against the proposed sparse pattern) rather than reasoned about in the
   abstract, because the failure mode the operator flagged
   ("could silently break a... lookup that happens to resolve through a
   path pattern not anticipated at design time") is precisely a
   not-anticipated-in-advance failure.
2. A **decision on the manifest surface** — whether checkout scoping is
   a per-wake manifest field (`wake.surfaces.excluded_checkout_paths` or
   similar) or a renderer-global default, and whether it needs an
   escape hatch for cell types that legitimately need `.cn-{agent}/`
   read access (none identified today, but the regression matrix in (1)
   is the mechanism that would surface one if it exists).
3. **A live-fire validation window** — because this is a production wake,
   the operator will likely want at least one full cycle run under the
   new checkout shape, observed end-to-end (claim → scaffold → implement
   → review → merge), before trusting it as the default for all future
   `protocol:cds` cells.
4. **Only then** does AC4 (write-fence retirement) become safe to attempt
   — AC4 is explicitly gated on AC3 being proven, not merely attempted;
   removing the fence before the capability-removal is validated would
   leave a real gap (a regression in the sparse-checkout mechanism could
   silently reintroduce `.cn-{agent}/logs/` visibility with no guard left
   to catch it).

This section is the required doctrine+plan deliverable; it is a
first-class output of this cell, not a placeholder for future work no
one has thought through.

## Review-readiness | round 0 | base SHA: 86042ec5be4b5fb45b213c27dfcf635958f60aac | head SHA: (this commit) | branch CI: not yet observed (local sha256/idempotence checks reproduced the install-wake-golden gate's logic; live CI on push not yet polled by α) | ready for β

## Review-ready — R0 complete, ready for β

---

# self-coherence.md — cnos#626 (α, R1 — AC3/AC4 continuation)

## Scope of this round

Continuation dispatch per the issue's "AC3/AC4 continuation dispatched
(kappa via CAP, 2026-07-08)" comment, which reopened `status:todo` for a
bounded scope: **AC3 + AC4 only**. AC1/AC2 (this branch's R0) already
merged as PR #635 and are not touched here. Mechanism was pre-decided by
the operator/CAP comment: sparse-checkout excluding `.cn-{agent}/` from
the dispatch wake's own checkout step, in-repo, no new infra.

## run_class note (honest taxonomy gap)

This claim does not cleanly match any of the three `run_class` values
`cds-dispatch/SKILL.md` §"Repair re-entry preflight" Step A currently
names (`first_pass`, `resumed_from_matter`, `repair_pass`): there is no
operator-rejection evidence (rules out `repair_pass`), and no scanner
"MECHANICAL reversion" audit-note comment exists on the issue (rules out
`resumed_from_matter` per `delta/SKILL.md` §9.11's own detection rule,
which requires that phrase). This is a fourth shape — an
operator-authorized **scope continuation** on an issue whose prior round
already converged and merged — not yet named in either doctrine surface.
Treated operationally as equivalent to `first_pass` for the new AC3/AC4
scope, since `.cdd/unreleased/626/` carried R0's artifacts (this branch's
history) but no artifacts yet for THIS round at claim time. Flagging this
gap explicitly rather than silently picking the nearest label; a future
cell should consider naming this shape in `cds-dispatch/SKILL.md` /
`delta/SKILL.md` §9 alongside §9.10/§9.11.

## Round-1 prerequisite check against R0's own follow-on requirements

R0's "AC3/AC4 — doctrine + plan, deferred to operator" section (above)
named four prerequisites a follow-on cell would need before attempting
candidate mechanism 1 (sparse-checkout) safely. Addressed here:

1. **Regression matrix over existing cell types' checkout-path needs.**
   Checked empirically, not just reasoned abstractly: grepped every
   `SKILL.md` under `cnos.cdd`, `cnos.cds`, `cnos.core` for
   `.cn-sigma`/`.cn-{agent}` references. Zero hits in `gamma/SKILL.md`,
   `alpha/SKILL.md`, `beta/SKILL.md` (the role skills γ/α/β actually
   load), `cds/SKILL.md` (the concrete protocol skill), or
   `issue/SKILL.md` (the cell contract skill) — none of the paths a
   cell's own cognition consults ever resolve through `.cn-{agent}/`.
   The only hits are in `delta/SKILL.md` §9.12, `cds-dispatch/SKILL.md`,
   and the agent-identity-layer skills (`activate`, `attach`,
   `registration`, `dispatch-protocol`, `label-doctrine`) — all
   wake/substrate-layer or doctrine-declaration surfaces cell-execution
   cognition does not load per §9.12's own rule, not cell-read paths.
   This directly answers R0's stated concern (a "not-anticipated-in-
   advance" path dependency) with a concrete negative result.
2. **Manifest-surface decision.** Implemented as a renderer-level
   `role == "dispatch"` gate in `cn-install-wake`, not a new
   per-wake manifest field (e.g. `wake.surfaces.excluded_checkout_paths`).
   Rationale: exactly one dispatch-role wake exists in production today
   (`cds-dispatch`); a configurable manifest field would be an
   unused-today abstraction for a hypothetical second dispatch wake with
   different needs, and none has been named. If a future dispatch wake
   genuinely needs `.cn-{agent}/` read access, that need itself would be
   evidence against the cell/substrate boundary this cycle establishes
   (per `delta/SKILL.md` §9.12: "A cell that needed... hub-state to do
   its job would be evidence the cell-execution/substrate-identity
   boundary itself is wrong"), not evidence the renderer needs an escape
   hatch. No escape hatch added.
3. **Live-fire validation window.** NOT satisfiable within this cycle by
   construction: a change to the dispatch wake's own checkout step can
   only be observed running for real on the wake's NEXT firing after this
   PR merges — there is no way to "dry-run" the actual GitHub Actions
   substrate from inside the cycle that authors the change. What IS
   verifiable now, and was: (a) the git-level mechanism itself, proven
   empirically against both a throwaway fixture repo and a real clone of
   this repo (`.cn-sigma` absent post-sparse-checkout, all 43 of its
   files the only diff against a full checkout, nothing else affected);
   (b) the renderer emits byte-identical output for the admin wake
   (regression-free); (c) the full Go test suite passes. This is strong
   evidence but not the same as an observed live firing. **Recommendation
   carried to closeout:** treat the next 1-2 real `cds-dispatch` firings
   post-merge as the live-fire validation window; do not treat AC3 as
   fully production-proven until at least one clean firing is observed.
4. **AC4 gate.** Consistent with (3): AC4 (write-fence retirement) is
   NOT attempted this round. The write-fence
   (`dispatch_activation_log_write_violation`) is confirmed byte-identical
   to its prior state (see AC verification below) and stays in place
   until AC3 is observed holding in a real firing, per the operator's own
   explicit gate ("retire the write-fence ONLY AFTER AC3 is proven
   in-cell") and STOP condition ("fence removal is attempted before AC3
   proof").

## AC3 — capability removal (implemented)

**Mechanism.** `src/packages/cnos.core/commands/install-wake/cn-install-wake`:
the `actions/checkout@v4` step, when `role == "dispatch"`, now emits:

```yaml
sparse-checkout: |
  /*
  !/.cn-${agent}
sparse-checkout-cone-mode: false
```

Regenerated both `.github/workflows/cnos-cds-dispatch.yml` and
`src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`
from the patched renderer (byte-identical to each other, as required by
the golden-fixture discipline AC2's own round established).

**Proof the capability is gone (not merely instructed against), per
AC3's own oracle wording:**

- Empirically ran `git sparse-checkout set --no-cone '/*' '!/.cn-sigma'`
  against a throwaway fixture repo containing `.cn-sigma/logs/` +
  `.cn-sigma/spec/` content plus unrelated files: post-checkout,
  `.cn-sigma` is absent (`os.IsNotExist` on stat) while unrelated paths
  (`src/go/go.mod`, `README.md`) remain present.
- Repeated the same command against a REAL clone of this repo (not a
  synthetic fixture) and diffed `git ls-files` (full) against the
  resulting sparse working tree's file list: the only difference is the
  full set of `.cn-sigma/*` entries (43 files) — nothing else is
  affected.
- Added `TestDispatchRenderer_SparseCheckoutExcludesAgentHub` to
  `src/go/internal/repoinstall/repoinstall_test.go`, which automates both
  the YAML-presence check (dispatch: present with the right patterns;
  admin: absent) and the real-git-mechanism proof (throwaway repo fixture
  + `git sparse-checkout set --no-cone` + stat assertions), so this proof
  re-runs on every future CI run rather than being a one-time manual
  check that could silently bit-rot.
- Confirmed (independently, not just by design intent) that this test is
  non-vacuous: sabotaging the renderer's gate condition in a scratch copy
  makes the test fail as expected (see β's independent re-derivation
  below, which repeated this sabotage check itself rather than trusting
  this claim).

**Admin wake unaffected.** `role == "admin"` renders no sparse-checkout
block; the admin wake's own golden fixture
(`src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`)
re-renders byte-identical to its committed state — confirmed via diff,
not assumed.

## AC4 — write-fence retirement (deliberately NOT done this round)

Per the prerequisite-4 reasoning above, the
`dispatch_activation_log_write_violation` write-fence step and its
detection logic are untouched. Confirmed via diff: the only changes in
both rendered workflow YAMLs are the 4 added lines of the sparse-checkout
block; the fence step's content is unchanged (only its line offset
shifts, because the new block is inserted earlier in the file).

## AC5/AC6 — regression guardrails

- Zero Go source files outside the new test changed (only
  `repoinstall_test.go` gained a new test function; no production Go
  changed).
- `src/packages/cnos.cds/skills/cds/fsm/transitions.json` untouched.
- No new `cell_kind` taxonomy; no Demo 0 work; `#633` not touched.
- Full test suite green: `go test ./...` (src/go module) plus the three
  standalone `go.mod` roots under `src/packages/*/commands/*/`
  (`issues-fsm`, `issues-map`, `cdd-verify`) — all `ok`, `go vet` clean.

## Review-readiness | round 1 | base: R0's merged state (PR #635, `9d1276cf` on main) + this branch's own R0 history | head SHA: (this commit) | branch CI: not yet observed at push time (local test suite + `go vet` green; live CI to be confirmed by β/δ post-push) | ready for β

## Review-ready — R1 complete, ready for β

---

# self-coherence.md — cnos#626 (α, R2 — AC4 write-fence proof-first round)

## Scope

This round addresses **AC4 only** (write-fence retirement/narrowing decision), per the operator/CAP comment posted 2026-07-09T00:42:46Z ("AC4 dispatched — proof-first"), which superseded R1's "reverting to status:ready" note. R0 (AC1/AC2, PR #635) and R1 (AC3, PR #637, `e5a3abe5`) are merged and untouched. `cycle/626` was reset to current main (`e7277e00`) before this round, per the operator's "clean-base start" framing — content-verified: `git diff origin/main origin/cycle/626 --stat` showed only a trivial 2-file board-map diff before this round's commits.

## Required-proof checklist (operator's 2026-07-09 comment) — evidence

### 1. Extend the sparse-checkout test — proves persistence, not just presence/absence

`TestDispatchRenderer_SparseCheckoutExcludesAgentHub` (existing, `src/go/internal/repoinstall/repoinstall_test.go`) proves AC3's presence/absence claim (the checkout omits `.cn-sigma`) but never exercised what happens if something writes a `.cn-sigma` path anyway. Two new tests close that gap, in `src/go/internal/cell/cell_test.go`:

- `TestRealCheckpoint_NewAgentHubFile_FailsLoudNoPersistence` — builds a real bare "origin" + a real sparse clone (same `/*` + `!/.cn-sigma` non-cone pattern the renderer emits), creates a brand-new file at `.cn-sigma/logs/probe.md` alongside a legit in-cone change, invokes the REAL `realCheckpoint` (not the no-op `Checkpoint` fake every other Finalizer test injects), and asserts (a) `realCheckpoint` returns a non-nil error (git's own `git add -A` exits 1 on an out-of-cone new path — fails LOUD, not silent), and (b) nothing under `.cn-sigma` ever reaches origin's `cycle/62601` branch.
- `TestRealCheckpoint_ModifiedTrackedAgentHubFile_SilentlyExcludedNoPersistence` — same setup, but instead overwrites an ALREADY-TRACKED, sparse-hidden path (`.cn-sigma/logs/20260101.md`, mirroring the real cnos repo's 40 pre-existing tracked channel-log files) directly on disk. Asserts `realCheckpoint` succeeds with no error (the legit change is not collaterally blocked), the `.cn-sigma` modification never reaches origin, AND the legit change DOES reach origin (proving the exclusion doesn't silently eat real work in this shape).

Both tests pass (see `go test ./internal/cell/... -run TestRealCheckpoint -v`, R2 run log). Building this required one iteration: the first draft of the remote-inspection helper used `git log --all` and false-positived on `main`'s own pre-existing tracked `.cn-sigma` seed content; fixed to `git log origin/main..origin/<branch>` (commits unique to the cycle branch only), matching the same exclusion the OLD write-fence itself already used for its own commit-graph layer.

### 2. Verify finalizer behavior

Covered by the same two tests above — they invoke the production `realCheckpoint` directly (not a mock), which is the exact code the rendered `Mechanical checkpoint + PR finalizer` step runs (`cn cell finalize` → `Finalizer.Finalize` → `realCheckpoint` when `Checkpoint` is nil).

### 3. Verify generated workflow

- `cn-install-wake` (the renderer, `src/packages/cnos.core/commands/install-wake/cn-install-wake`) changed: the `Write fence — dispatch_activation_log_write_violation` step's emission condition narrowed from `[ "$activation_log_writer" = "false" ]` to `[ "$activation_log_writer" = "false" ] && [ "$role" != "dispatch" ]`. The `Write fence — record pre-work baseline SHA` step (which the finalizer's `--base-sha` flag also consumes) is UNCHANGED — still emitted whenever `activation_log_writer` is false, regardless of role.
- Re-rendered both golden fixtures + the live workflow:
  - `cn-install-wake agent-admin` → **unchanged** (`activation_log_writer: true` for agent-admin; unaffected by the role-gated carve-out).
  - `cn-install-wake cds-dispatch` (golden) and `cn-install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml` (live) → both changed identically; sha256 of live == sha256 of golden (`b1f801c7...38fe`).
  - Idempotence re-checked: a second `cn-install-wake cds-dispatch` run reports `(unchanged)` and produces a byte-identical sha256.
  - `python3 -c "import yaml; yaml.safe_load(...)"` — both goldens parse.
  - Substrate structural-shape assertions from `install-wake-golden.yml` (on/permissions/concurrency/jobs/claude-code-action/labeled-event/concurrency-group-name/OG-2 schedule gate) re-run inline against the new golden — all pass.
  - AC5 declaration-only refusal fixture — re-run, exit 3, correct stderr content.

### 4. Live firing after the change

**Not attempted in this round — structurally cannot be, same as R1's AC3 proof.** Per δ wake-invoked mode's own v0 substrate constraint (`delta/SKILL.md` §9.7): a single firing does not plan multi-firing continuations. AC10 ("live firing validates checkout/scanner/finalizer after the change") is satisfied by the NEXT real `cds-dispatch` firing after this cycle's PR merges — exactly the same "live-fire window" mechanism R1's closeout named for AC3 ("the next 1-2 real cds-dispatch firings are the AC3 live-fire window"). This is named explicitly here, not glossed over, per the operator's "Do NOT: ... proof requires new infra/permissions" framing — AC10 requires no new infra, just a subsequent real firing, which this round cannot manufacture.

## Outcome selected: **B**

> "local files can be created but cannot be staged/persisted: retire the old local-write fence only if the durable boundary is proven; document the distinction (transient local file creation ≠ substrate write access)."

Both new tests confirm the "transient local file creation" half (a `.cn-sigma` file/modification CAN appear on disk — `git status --porcelain` sees it) and the "cannot be staged/persisted" half (git's own sparse-checkout enforcement at `git add -A` time refuses the new-file case loudly and silently excludes the modified-tracked-file case) — for the exact code path (`realCheckpoint`) production uses. Outcome A (unconditional full removal) was NOT selected because the fence's underlying shell logic remains valid protection for a hypothetical future `activation_log_writer:false` wake without `role:dispatch` (legal per the manifest schema — only the reverse implication is enforced) — narrowing by `role` rather than deleting unconditionally avoids silently leaving such a hypothetical wake unguarded. Outcome C (keep the fence, do not remove) was ruled out by the negative proof itself.

## AC1–AC10 verification (independent per-AC walk, α's own pass before handing to β)

| AC | Oracle | Status | Evidence |
|---|---|---|---|
| AC1 | sparse-checkout boundary test proves `.cn-sigma` absent from dispatch checkout | ✅ | `TestDispatchRenderer_SparseCheckoutExcludesAgentHub` (unchanged, still passing — R1's AC3 proof, re-verified this round) |
| AC2 | attempted `.cn-sigma/logs` write cannot persist through dispatch/finalizer | ✅ | `TestRealCheckpoint_NewAgentHubFile_FailsLoudNoPersistence` + `TestRealCheckpoint_ModifiedTrackedAgentHubFile_SilentlyExcludedNoPersistence` |
| AC3 | old fence removed only if AC2 passes | ✅ | AC2 passed → fence removed for `role=="dispatch"` (renderer conditional) |
| AC4 | if fence remains/narrows, receipt explains why | ✅ | Fence NARROWS (kept for hypothetical non-dispatch `activation_log_writer:false`); explained in renderer comments + `cds-dispatch/SKILL.md` doctrine update + this section |
| AC5 | CDS golden/live match | ✅ | sha256 identical (`b1f801c7...`) |
| AC6 | install-wake-golden green | ✅ | All install-wake-golden.yml steps reproduced locally: golden diff clean, live==golden sha, idempotence, YAML parse, structural shape, AC5 refusal fixture, AC2 negative fixture |
| AC7 | dispatch-repair-preflight green | ✅ | `./scripts/ci/check-dispatch-repair-preflight.sh` passed |
| AC8 | dispatch-closeout-integrity green | ✅ | `./scripts/ci/check-dispatch-closeout-integrity.sh` passed |
| AC9 | Go/Package/Binary and I1/I2/I4/I5/I6 green | ✅ (I4/I5 not runnable locally — see note) | `go build ./...`, `go vet ./...`, `go test ./...` (all 4 go.work modules) green; I1 (`cn build --check`) green; I2 (protocol-contract diff) green; I6 (`cn cdd verify --unreleased`) green (0 failed, pre-existing warnings only). I4 (lychee link-check) and I5 (cue frontmatter validation) tools are not installed in this session's environment — NOT run locally. No repo links were added/changed by this round's edits (prose-only additions inside existing files) and the SKILL.md frontmatter block itself was not touched (only body prose), so both are expected green in real CI, but this is a **named gap**, not a silent assumption. |
| AC10 | live firing validates checkout/scanner/finalizer after the change | ⏳ deferred | Structurally requires a firing after this cycle's PR merges (see "4. Live firing" above); consistent with R1/AC3 precedent |

## Files changed this round

- `src/go/internal/cell/cell_test.go` — two new tests + one new helper (`newSparseOriginAndClone`) + one new assertion helper (`remoteBranchTouchedPath`); +`os/exec` import.
- `src/packages/cnos.core/commands/install-wake/cn-install-wake` — narrowed the write-fence emission condition; added explanatory comment block naming the AC4 decision and its evidence.
- `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` — updated the "Disallowed surfaces" prose (line ~302) and the "Responsibilities (body reference)" item 9 (line ~370) to describe the retired-for-dispatch fence and point at the new tests.
- `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` + `.github/workflows/cnos-cds-dispatch.yml` — regenerated (mechanical; not hand-edited).

## run_class note (repeat of the taxonomy gap named at claim time)

This round is the SAME "fourth shape" the R1/AC3 round's artifacts already flagged: an operator-authorized scope-continuation on an issue whose prior rounds already converged and merged, run against a branch reset to a clean base. `cds-dispatch/SKILL.md`'s current taxonomy (`first_pass` / `resumed_from_matter` / `repair_pass`) still has no name for it. Recording again here per the operator's own instruction ("a future doctrine cell should add a named fourth shape... recorded in self-coherence.md, not filed as a separate issue") rather than filing a new issue for it a second time.

## Review-readiness | round 2 | base: cycle/626 reset to main@e7277e00 (R0 = PR #635, R1 = PR #637 both already merged into that base) | head SHA: (this commit) | branch CI: not yet observed at push time (local test suite + go vet + install-wake-golden-equivalent checks + dispatch-repair-preflight + dispatch-closeout-integrity all green; live CI to be confirmed by β/δ post-push) | ready for β
