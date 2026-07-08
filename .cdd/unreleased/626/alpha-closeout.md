# α close-out — cnos#626

## Scope note

This close-out is written at the **δ §9.5 converge boundary** — β
returned `verdict: converge` at R0 (`beta-review.md`, reviewing commits
`51c2c1c` + `87c9ed1`), and this file is one of the three artifacts that
boundary requires (`alpha-closeout.md` + `beta-closeout.md` +
`gamma-closeout.md`), per `delta/SKILL.md` §9.5's per-R[N] artifact
contract for the `converge` row. As with cnos#630's equivalent closeout,
this is **not** the post-merge/post-release close-out `alpha/SKILL.md`
§2.8 describes — no PR exists yet for `cycle/626`, and opening the
cycle-PR plus requesting the `status:review` transition is δ's next
action per §9.6, not this pass's. γ, α, and β roles for this closeout
round are collapsed into a single authoring pass per Sigma's engineering-
persona protocol commitment #5 (β-α-collapse-on-δ for skill/docs-class
cycles), since this cycle's own AC oracle is mechanical (file existence,
grep counts, byte-diffs) and β's R0 review (zero findings, `verdict:
converge`) is already independently on record.

## Summary — what was built

cnos#626 asked for the dispatch cell-worker (δ/γ/α/β cognition, as run
inside `cds-dispatch`) to stop being fused with substrate/agent identity
and hub state — the operator's own framing: "A worker cell isn't even
supposed to know about sigma or activation logs. Its job is just to
produce code, reviews, etc." The issue was explicitly narrowed by the
operator to a **supervised design-first cell**, dated 2026-07-08: ship
only the bounded/safe mechanism; STOP at doctrine + implementation plan
for anything not clearly bounded within a single cell.

γ's scaffold (`gamma-scaffold.md`) evaluated the issue's three candidate
boundary moves and found exactly one cleanly implementable without new
infra or a checkout-shape change to a live production wake: **drop the
full six-step sigma activation for cell-execution cognition; load Kernel
+ CA skills only.** That move shipped as two doctrine/prose edits plus
their required regenerated mirrors:

1. `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`
   §"Identity and activation" — the "Activate per activate/SKILL.md:
   Kernel → CA skills → Persona → Operator → hub state → identity
   confirmation" instruction is replaced with doctrine stating
   cell-execution cognition loads Kernel + CA skills only, does not load
   Persona/Operator/hub-state, and does not read `.cn-{agent}/` or
   thread surfaces. The wake's own self-identification is
   manifest-declared (the `{agent}` variable + bot-account binding), not
   file-discovered. The pre-existing "you do NOT attach to a channel
   like the admin wake does" sentence is preserved verbatim.
2. `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.12 "cell/
   substrate identity boundary (cnos#626)" — new doctrine, sibling to
   §9.10/§9.11, stating explicitly: the cell carries no substrate
   identity and no hub-state scope; the wake/substrate owns identity,
   channel logs, and reporting. It cross-references §9.2's existing
   five-input contract and names the Persona/Operator/hub-state omission
   from that list as intentional, not an oversight.
3. Both rendered mirrors of the edited `cds-dispatch/SKILL.md` prompt
   body were regenerated via the actual on-branch renderer
   (`src/packages/cnos.core/commands/install-wake/cn-install-wake
   cds-dispatch`, invoked once per `--out` target) and independently
   confirmed byte-identical to each other by both α and β (matching
   sha256 `8f68971...76ff1`), reproducing the `install-wake-golden` CI
   gate's own check.

Full diff against `origin/main`: 4 substantive files (two SKILL.md prose
files, two regenerated YAML mirrors) plus the three `.cdd/unreleased/626/`
doctrine artifacts (`gamma-scaffold.md`, `self-coherence.md`,
`beta-review.md`) — 22 insertions / 3 deletions in the substantive diff,
zero Go source touched, zero FSM/`transitions.json` touched.

## What was deliberately deferred, and why

**AC3** ("a code cell cannot see or write `.cn-*/logs/` — the capability
is removed, not merely guarded") and **AC4** (hard-gated on AC3: "once
AC3 holds, the `dispatch_activation_log_write_violation` write-fence is
removed as redundant") were **not attempted this cycle** — a deliberate
STOP, not an oversight or a silently dropped scope item.

The reason is the operator's own explicit STOP condition, restated in
the scaffold and honored exactly: *"if the safe mechanism is not clearly
bounded, STOP at doctrine + implementation plan and hand back to the
operator rather than force an implementation."* Both candidate mechanisms
for AC3 — (1) sparse-checkout/worktree isolation excluding `.cn-sigma/`,
and (2) relocating `.cn-sigma/` out of the product repo entirely — touch
surfaces the operator named explicitly as unsafe to force inside this
cell: (1) changes the checkout shape of a *live production* wake
(`activation_state: live` since cnos#487) that other in-flight and future
cells depend on for reading the full codebase and running the finalizer,
with a real risk of silently breaking a doc cross-reference or Tier-3
skill path lookup not anticipated at design time; (2) requires new
infra/permissions ("escalate," per the operator directive) that this cell
has no authority to provision.

`self-coherence.md`'s "AC3/AC4 — doctrine + plan, deferred to operator"
section is the required deliverable in place of an implementation: it
restates both candidate mechanisms with their tradeoffs and names four
concrete requirements a follow-on cell would need before attempting AC3
safely (a regression matrix over every existing cell type's checkout-path
needs; a decision on the manifest surface for checkout scoping; a
live-fire validation window on the production wake; and the explicit
AC3-before-AC4 ordering constraint). This is treated, per the operator's
own design-first bar, as a first-class output equal in weight to code —
not filler narrative.

AC2's issue text is worded as an "and/or" ("no longer performs a full
sigma activation **and/or** no longer has `.cn-{agent}/` in its working
scope"), which makes the first disjunct alone — the doctrine/prose move
that shipped — a legitimate, honest AC2 pass without needing the second
disjunct (checkout-scope exclusion, i.e. AC3's mechanism). β's review
independently confirmed this reading was applied by α before dispatch,
not invented afterward to justify skipping harder work.

## Debt disclosed

- **`cn install-wake` naming discrepancy.** The scaffold's and α prompt's
  literal invocation form (`cn install-wake cds-dispatch --out <path>`)
  does not match the actual on-branch tooling: the renderer is a
  standalone POSIX shell script at
  `src/packages/cnos.core/commands/install-wake/cn-install-wake`, not a
  subcommand of the `cn` Go binary (`cn help` lists no `install-wake`
  verb). This is a naming/packaging mismatch between prose and tooling,
  not a functional gap — the correct script was invoked with the correct
  arguments and produced output independently verified byte-identical
  against both the golden fixture and the live workflow mirror by both α
  and β. Flagged in `self-coherence.md` for visibility rather than
  silently reconciled; worth correcting the scaffold-prose convention in
  a future cycle so it stops citing a `cn` subcommand form that does not
  exist on this branch.
- **Git author identity.** Commits on this cycle carry the pre-existing
  session git identity (`sigma@cnos.cn-sigma.cnos`), not the canonical
  `alpha@cdd.cnos` pattern `alpha/SKILL.md` §2.6 row 14 names. Disclosed
  as known debt under that row's permitted path (b) — "configure
  correctly going forward and accept existing commits as legacy" — rather
  than rewritten, since amending authorship was not a requirement named
  anywhere in the α prompt and an unprompted history rewrite would be a
  destructive git operation the operator did not ask for.
- **Local-only renderer verification.** The `install-wake-golden` sha256
  + idempotence checks were reproduced locally (by both α and β,
  independently) rather than observed via a live GitHub Actions run on
  the pushed branch at the time `self-coherence.md` was written. No
  behavior difference is expected (same script, same manifest, same
  prompt-body text runs in both places), but this is named as a residual
  gap between "locally reproduced" and "CI-observed" until a live run is
  polled.

## Calibrated success claim

What is verified, as of this boundary:

- AC1 is met: `delta/SKILL.md` §9.12 explicitly states both required
  halves (no substrate identity/hub-state scope on the cell; wake/
  substrate owns identity + channel logs + reporting) and explicitly
  names the §9.2 five-input contract's Persona/Operator/hub-state
  omission as intentional — independently re-confirmed by β reading the
  actual section text, not α's paraphrase.
- AC2's first disjunct is met: the six-step full-activation instruction
  is gone from `cds-dispatch/SKILL.md`'s "Identity and activation"
  section (`grep -n "Activate per.*activate/SKILL.md"` → 0 hits,
  independently re-run by β), replaced with Kernel+CA-skills-only
  doctrine; both rendered mirrors were regenerated via the renderer (not
  hand-edited) and are byte-identical to each other, independently
  reproduced by β from source rather than trusted from a pre-existing
  sha256 comparison.
- AC5/AC6 hold as pure-regression guardrails: zero Go files changed, zero
  FSM/`transitions.json` changes, zero new `cell_kind`/`status:*` label,
  write-fence code and `.cn-sigma/` contents byte-identical to `main`,
  `activate/SKILL.md` itself unmodified — all independently re-run by β.

What is explicitly **not** attempted or claimed:

- AC3 (capability removal) and AC4 (write-fence retirement) were not
  implemented this cycle. No claim is made that the write-fence is
  redundant or that `.cn-sigma/logs/` visibility has been structurally
  removed from any cell's working tree — it has not. The doctrine+plan
  section documenting why and what a follow-on cell needs is the
  complete discharge of this cell's obligation toward AC3/AC4, not a
  partial implementation of either.
- No PR has been opened and no merge to `main` has happened — that is
  δ's next step (§9.6 `status:review` transition), not something this
  cycle's α/β passes have done or claim to have done.
- No live wake-invoked-δ firing has exercised the new §9.12 doctrine
  yet (this cycle itself ran under bootstrap-δ, per the scaffold's own
  dispatch-mode note) — the doctrine is inert prose until the next wake
  firing reads it post-merge, mirroring the same gap named for §9.10/
  §9.11 in the cnos#630 closeout precedent.

**Claim:** the bounded mechanism (AC1, AC2's first disjunct) is correctly
and completely implemented, independently verified by β's from-scratch
re-derivation of every grep/diff/sha256 check rather than by trusting
this file's or `self-coherence.md`'s narrative. AC3/AC4 are honestly and
substantively deferred to the operator per the operator's own
design-first STOP condition — this is the designed outcome of a
supervised design-first cell, not a partial failure. It is ready for δ
to proceed to the `status:review` transition.

---

## §R1 amendment — AC3/AC4 continuation (cnos#626)

R1 implemented AC3 (sparse-checkout excluding `.cn-{agent}/` from the
dispatch wake's own checkout, gated to `role == "dispatch"`) with
empirical proof at both the git-mechanism level and a new automated Go
test (`TestDispatchRenderer_SparseCheckoutExcludesAgentHub`). AC4 (write-
fence retirement) was deliberately NOT attempted, per the operator's
explicit gate and R0's own follow-on-cell prerequisite list (a live-fire
validation window this cycle cannot self-supply). Zero regressions: full
test suite green, admin wake's golden fixture re-renders byte-identical.
`run_class` for this continuation does not cleanly match
`cds-dispatch/SKILL.md`'s current taxonomy (`first_pass` /
`resumed_from_matter` / `repair_pass`) — flagged explicitly in
self-coherence.md rather than silently mislabeled.
