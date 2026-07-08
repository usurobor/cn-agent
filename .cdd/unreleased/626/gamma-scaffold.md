# γ R0 scaffold — cnos#626

## Issue reference

- **Issue:** [usurobor/cnos#626](https://github.com/usurobor/cnos/issues/626) — "arch(cds): isolate the dispatch cell-worker from substrate/agent identity + hub state"
- **Mode:** design-and-build (doctrine-first; likely master+subs) — **operator-narrowed to a supervised design-first cell** (see "Operator directive" + "Authority — A2/CAP" sections of the issue body, dated 2026-07-08). Bar: implement ONLY the bounded/safe mechanism; STOP at doctrine + implementation plan for anything not clearly bounded within this single cell.
- **Cell kind:** `doctrine` → `implementation` (partial — see "Recommended design" below for the bounded/unbounded split).
- **Family:** #614 (the cell during whose *prior firing* the write-fence false-positived — RCA'd and fixed under #625, not #614 itself; #614's own title is an unrelated CELL-KINDS amendment, confirmed via `gh issue view 614`) / #625 (merged fix: "wake: fix write-fence false-positive on concurrent admin-wake log commits" — the tactical patch this issue's structural fix sits behind) / #496 (cnos#496 / cycle/496 — the write-fence guard AC4 ultimately retires) / #583 (mechanical-dispatch doctrine — "guards patch symptoms; the boundary is wrong," the exact pattern this issue extends) / #584 (mechanism/cognition boundary doctrine) / #467 (two-wake admin/dispatch architecture — this issue extends the split one boundary further) / #449 (per-package bot identity follow-up, cited as a future unblocker for `{agent}`-scoped substrate identity) / #630 (immediately preceding cycle — explicitly named as "next up after this cell" in #630's own scaffold; #630's own non-goals explicitly exclude #626 work) / #633 (open; FSM rule-ordering fix — **do not touch unless directly required**, per operator directive).
- **Protocol:** cds
- **Dispatch mode:** **bootstrap-δ, single-session δ-as-γ via the Agent tool (Claude Code)** — per `delta/SKILL.md` §9.1's own table. **Not** wake-invoked-δ (§9.1: wake-invoked-δ is "not yet observed in production" as of this cycle). This scaffold is itself being authored inside a checkout that has `.cn-sigma/` present — see "Friction notes" §3 for the self-referential bootstrapping consideration this creates.
- **Base SHA:** `86042ec5be4b5fb45b213c27dfcf635958f60aac` — verified: `git rev-parse HEAD` and `git rev-parse origin/main` both resolve to this SHA at scaffold time on `cycle/626` (zero drift since branch creation).
- **Cycle branch:** `cycle/626`, already created and pushed from `main@86042ec5be4b5fb45b213c27dfcf635958f60aac` (per the δ-invocation prompt; confirmed via `git branch --show-current` = `cycle/626`).
- **run_class:** `first_pass` — per `.cdd/unreleased/626/CLAIM-REQUEST.yml`: "Clean first dispatch -- no cycle/626 branch or prior claim exists." No prior `status:changes` history, no prior `cycle/626` artifacts beyond `CLAIM-REQUEST.yml` itself, matching the issue's own "Clean first dispatch" note.

## Governing observation (operator-stated, restated)

> A worker cell isn't even supposed to know about sigma or activation logs. Its job is just to produce code, reviews, etc.

Today it does know. The **core rule** the operator restates in the dispatch-authorization comment: the worker cell produces code/reviews/matter; **the wake/substrate** owns Sigma identity, hub state, activation logs, and reporting. The cell-worker should NOT perform full Sigma activation, load Persona/Operator/hub state, see or write `.cn-sigma` logs, or own activation-log behavior.

## Key finding 1 — the fused activation instruction, exact location

`src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` §"Identity and activation" (lines 99–105; identically duplicated at the rendered golden `cnos-cds-dispatch.golden.yml:99–103` and the live `.github/workflows/cnos-cds-dispatch.yml:99–103` — the prompt body is embedded verbatim in both rendered artifacts):

```
Activate per [`src/packages/cnos.core/skills/agent/activate/SKILL.md`]
(../../../cnos.core/skills/agent/activate/SKILL.md): Kernel → CA skills →
Persona → Operator → hub state → identity confirmation. The
identity-confirmation statement names you as the dispatch wake, not the
admin wake. Do NOT execute any further action until identity confirmation
completes.
```

This is the literal instruction the issue's problem statement #2 cites. It tells the thing whose job is "claim issue #N and run δ→γ/α/β" to run the **full six-step agent-at-a-hub activation** before touching the cell.

## Key finding 2 — the six-step load order this pulls in

`src/packages/cnos.core/skills/agent/activate/SKILL.md` §2.1 (lines 106–137) defines the six-item load order: (1) Kernel, (2) CA skills [cap, clp] — together "the soul, substrate-independent" — then (3) Persona, (4) Operator, (5) **hub state** ("dependency manifest, latest reflection, memory surfaces... thread surfaces: `<hub>/threads/{in,inbox,mail,archived}/`"), (6) identity confirmation. Step 5's hub-state survey is the step that, applied to cnos-as-hub, walks `.cn-sigma/` — exactly the surface #626 wants out of cell scope. §2.5 (lines 231–284) additionally specifies containerized-path detection under `.cn-{agent}/spec/PERSONA.md` / `OPERATOR.md`, i.e. the resolver explicitly walks into `.cn-sigma/` as a candidate identity path even before reaching hub-state survey.

Critically, `activate/SKILL.md` itself frames this as a **general agent-at-a-hub** procedure ("how does an AI body, knowing only a hub URL, reach a state where it can name its identity, its operator, and its current orientation" — governing question, line 6). A package-owned, protocol-scoped dispatch wake is not that shape: its identity (`cds-dispatch`, `protocol:cds`, substrate bot binding via the `{agent}` manifest variable) is **manifest-declared at render time**, not discovered by reading hub files at runtime. Applying the general hub-activation procedure to this narrower actor is the over-fusion.

## Key finding 3 — the actual substrate-identity surface, measured

- `.cn-sigma/` is git-tracked with **43 files** total (`git ls-tree -r --name-only HEAD -- .cn-sigma/ | wc -l`), of which **41 files live under `.cn-sigma/logs/`** (`git ls-tree -r --name-only HEAD -- .cn-sigma/logs/ | wc -l`). (The issue's problem statement cites "42 tracked files... 40 channel-log files"; the count has grown by one of each since the issue was filed — same order of magnitude, same structural fact: every dispatch checkout of `cnos` contains this directory.)
- `.cn-sigma/README.md` confirms by design: "This directory is Sigma's namespace when Sigma is activated against the cnos repo... Sigma's inbound mailbox at cnos... a file under `.cn-sigma/`." It is explicitly per-context state — living inside the *product* repo, not Sigma's home hub (`cn-sigma`).
- Every `cds-dispatch` checkout — including this very cycle's checkout — has `.cn-sigma/` present in the working tree because the wake's checkout step (rendered from `cn install-wake`) is a normal full-repo `actions/checkout@v4`, with no path exclusion.

## Key finding 4 — the write-fence this issue's AC4 targets, exact location

Mechanically enforced in `.github/workflows/cnos-cds-dispatch.yml` (step `Write fence — dispatch_activation_log_write_violation`, lines 358–439; identical logic in the golden fixture) and named in `cds-dispatch/SKILL.md` §"Disallowed surfaces" (line 302). Two layers: (1) working-tree `git status --porcelain` check against `.cn-*/logs/`, (2) local commit-graph diff (`$CN_WAKE_BASE_SHA..HEAD`) against the same pathspec, both `exit 1` with `dispatch_activation_log_write_violation` on a hit. This is a **guard**, not a **removed capability**: `.cn-sigma/logs/` is still checked out, still readable, still writable by any process running inside the job — the fence only catches the write *after the fact*, via git-status/diff inspection. This is exactly AC3's target ("a code cell cannot see or write `.cn-*/logs/` — the capability is removed, not just guarded").

`gh issue view 625` confirms the tactical fix landed separately and is `MERGED`: "wake: fix write-fence false-positive on concurrent admin-wake log commits." #626 is the durable structural fix behind that patch, per the issue's own framing ("The tactical fence fix ships separately (Refs #496); this issue is the durable boundary fix").

## Recommended design — the bounded/unbounded split (the operator's design-first bar, applied)

The issue names three candidate boundary moves. I evaluated each against what is achievable **inside this single cell, without new infra/permissions, without weakening skill loading, and without breaking codebase/issue-contract access** — the operator's explicit STOP conditions:

| Candidate move | Mechanism | Bounded in this cell? |
|---|---|---|
| **(3) Drop full activation for the cell; load Kernel + CA skills only** | Edit `cds-dispatch/SKILL.md` §"Identity and activation" prose (+ mirror in `delta/SKILL.md` §9.2 input-contract doctrine) to state explicitly: cell-execution cognition (δ/γ/α/β) loads **Kernel + CA skills only** (soul); it does NOT load Persona/Operator/hub-state (identity + hub survey); the wake's self-identification is manifest-declared (`{agent}` variable + bot binding), not file-discovered. Re-render via `cn install-wake cds-dispatch` into both the golden fixture and the live workflow. | **YES.** Pure prose/doctrine + a mechanical re-render step already covered by the existing `install-wake-golden` CI gate. No new Go surface, no new infra, no checkout change. Directly discharges **AC2's first disjunct** ("no longer performs a full sigma activation") — AC2 is worded as an "and/or," so this alone satisfies AC2. |
| **(1) Isolated checkout / sparse worktree excluding `.cn-sigma`** | Add a post-checkout `git sparse-checkout` step (non-cone mode: `/*` + `!/.cn-sigma/`) to the rendered workflow, gated by a new wake-manifest field (e.g. `wake.surfaces.excluded_checkout_paths`), requiring `cn install-wake` renderer (Go) changes. | **NOT cleanly bounded here.** This changes the checkout shape of a **live production** wake (`activation_state: live` since cnos#487) that other in-flight and future cells depend on for reading the full codebase, opening PRs, and running the finalizer (`cn cell finalize`, cnos#591). The operator's own STOP conditions name exactly this risk ("isolating the cell would break its ability to read the codebase or the issue contract"; "removing full activation weakens skill loading" — sparse-checkout could silently break a doc cross-reference or a Tier-3 skill path lookup that happens to resolve through a path pattern not anticipated at design time). Validating this safely needs its own AC5-equivalent regression pass (does every existing cell type still see everything it needs?) that this single cell cannot responsibly absorb alongside AC1/AC2. |
| **(2) Move `.cn-sigma/` out of the product repo entirely** | Relocate Sigma's per-context state to a separate repo/location; cnos becomes pure-codebase. | **NOT bounded here — explicit STOP condition.** "Moving the hub out of the repo requires new infra/permissions not available → escalate." This is an infra/ownership decision (where does the new location live, who has write access, does `activate/SKILL.md`'s Tier 1b pure-product-hub path already cover the read side) that the operator, not this cell, must authorize. |

**Decision: implement move (3) now; produce doctrine + plan for moves (1)/(2) and STOP.** This satisfies the operator's design-first bar precisely: the safe, bounded mechanism (dropping full activation) ships; AC1 (doctrine naming the boundary) ships alongside it; AC3/AC4 (capability removal, write-fence retirement) are named as the next wave's scope with the two candidate mechanisms and their tradeoffs recorded, but **not implemented** — because AC4 is explicitly gated on AC3 being proven first ("write-fence removed only AFTER the capability is removed"), and AC3 itself requires a checkout/infra decision this cell's STOP conditions correctly flag as not yet safe to force.

This is not a partial failure of the cell — it is the designed outcome of a "supervised design-first cell": ship what's clearly bounded, name what isn't, hand the rest back.

## Source of truth

| Claim / surface | Canonical source | Status |
|---|---|---|
| Fused activation instruction (the exact prose AC2 changes) | `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` §"Identity and activation" (lines 99–105) | On `main`; edit in place |
| Rendered mirrors of the same prose (must be regenerated, never hand-edited) | `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` (lines 99–103) + `.github/workflows/cnos-cds-dispatch.yml` (lines 99–103) | On `main`; regenerate via `cn install-wake cds-dispatch`, do not hand-edit (the `install-wake-golden` CI gate enforces byte-identity) |
| Six-step activation procedure (general agent-at-a-hub contract; NOT edited by this cell) | `src/packages/cnos.core/skills/agent/activate/SKILL.md` §2.1 (lines 106–137), §2.5 (lines 231–284) | On `main`; read-only cross-reference — this cell narrows *who is required to run the full procedure*, it does not change the procedure itself (admin wake and any true agent-at-a-hub still need all six steps) |
| δ wake-invoked input contract (where the "five inputs, no hidden state" doctrine already lives) | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.2 (lines 395–407) | On `main`; the five named inputs (claimed-issue-number, protocol identifier, current main SHA, wake run id, package-runtime context) already do NOT include Persona/Operator/hub-state — this cell makes that omission an **explicit, named doctrine statement** rather than a silent absence, satisfying AC1 |
| Write-fence guard code (the mechanism AC4 will retire, once AC3 is proven — NOT touched by this cell) | `.github/workflows/cnos-cds-dispatch.yml` lines 358–439 (+ identical golden mirror); named in `cds-dispatch/SKILL.md` §"Disallowed surfaces" line 302 | On `main`; out of scope — do not remove or weaken; AC3 is not proven this cycle |
| `.cn-sigma/` surface (the substrate-identity state being isolated) | `.cn-sigma/README.md`; `git ls-tree -r HEAD -- .cn-sigma/` (43 files, 41 under `logs/`) | On `main`; not modified by this cell — the cell changes what the *worker* is told to load/see, not the directory's contents or presence |
| `cn install-wake` renderer (authority for regenerating both rendered artifacts) | `src/packages/cnos.core/commands/install-wake/cn-install-wake` — usage: `cn install-wake cds-dispatch [--out <path>]`; default `--out` is the per-package golden, cutover path is `.github/workflows/cnos-<wake>.yml` | On `main`; invoke twice (once per `--out` target) after editing the SKILL.md body |
| install-wake-golden CI gate (verifies renderer determinism + golden/live sync) | `.github/workflows/install-wake-golden.yml` | On `main`; must stay green — this is the mechanical proof the two rendered mirrors match the edited SKILL.md |
| Two-wake admin/dispatch architecture (the split this issue extends one boundary further) | cnos#467 (closed); `cds-dispatch/SKILL.md` frontmatter `wake:` block, esp. `activation_log_writer: false`, `surfaces.disallowed` | On `main`; this cell's doctrine addition is a refinement, not a new architecture |
| The stranded-cell / write-fence RCA context (for accurate "Family" framing) | `.cdd/unreleased/614/gamma-closeout.md` line 31 ("prior firing hit an avoidable tooling failure (write-fence false-positive)... RCA'd and fixed in #625"); `gh issue view 625` (title: "wake: fix write-fence false-positive on concurrent admin-wake log commits", state MERGED) | On `main`/closed; cited for accurate lineage, not modified |
| Prior cell in this family, non-goals confirming #626 is separate scope | `.cdd/unreleased/630/gamma-scaffold.md` §"Scope" ("no #626 work (cell⇄substrate identity isolation — separate architectural wave...)") | Landed; confirms sequencing, not modified |

## Per-AC oracle list

| AC | Oracle (mechanical pass/fail) | This cell's disposition |
|---|---|---|
| **AC1** | A doctrine statement exists, on `main` after merge, naming: (a) the cell carries no substrate identity and no hub-state scope, (b) the wake owns identity + channel logs. Locate it via `grep` in `delta/SKILL.md` and/or `cds-dispatch/SKILL.md` for a named boundary section (not scattered inference from unrelated prose). | **Implemented this cell.** New subsection in `delta/SKILL.md` §9 (sibling to §9.10/§9.11 — e.g. §9.12 "cell/substrate identity boundary") + a rewritten `cds-dispatch/SKILL.md` §"Identity and activation" that states the rule in the wake's own prompt body. |
| **AC2** | The dispatch cell-worker's prompt/doctrine no longer instructs a full six-step sigma activation for cell-execution cognition (δ/γ/α/β), **and/or** no longer has `.cn-{agent}/` in its working scope. Oracle: `grep -n "Activate per.*activate/SKILL.md" cds-dispatch/SKILL.md` returns nothing (or returns only a narrowed reference scoped to Kernel+CA skills, not the six-step chain); the golden + live rendered mirrors match byte-for-byte post re-render (`install-wake-golden` CI green). | **Implemented this cell** via move (3) — the "and/or" is satisfied by the first disjunct alone. The second disjunct (checkout-scope exclusion) is explicitly NOT implemented this cell (see "Recommended design"). |
| **AC3** | A code cell cannot see or write `.cn-*/logs/` — the **capability** is removed (e.g. not present in the checkout, or checkout is scoped), not merely guarded post-hoc. Oracle: a fixture/dry-run showing the cell's working tree structurally excludes `.cn-*/logs/` (e.g. `ls .cn-sigma/logs/` fails inside the scoped checkout). | **NOT implemented this cell — doctrine + plan only, per operator STOP condition.** The two candidate mechanisms (sparse-checkout; relocate hub out of repo) are named and evaluated in "Recommended design" above with explicit tradeoffs; the implementation decision is handed back to the operator. β must confirm this is an honest STOP, not a silently dropped AC — see α prompt "When done" for the exact disclosure required. |
| **AC4** | Once AC3 holds, the `dispatch_activation_log_write_violation` write-fence is removed as redundant. | **NOT applicable this cell** — hard-gated on AC3, which is not proven this cycle. α MUST NOT touch the write-fence code (`.github/workflows/cnos-cds-dispatch.yml` lines 358–439) in any way. |
| **AC5** | No regression: FSM ownership, finalizer, reconciler, dispatch guards, and all gates stay green; cells still produce code/reviews/PRs identically. | **In scope.** Oracle: `install-wake-golden` CI green (proves golden/live-workflow sync + renderer determinism unchanged); `#516`/`#524` dispatch guard scripts (`check-dispatch-repair-preflight.sh`, `check-dispatch-closeout-integrity.sh`) still pass; no Go source touched this cell (prose-only diff) so `go test ./...` is a pure regression check, not a new-behavior check — must be unchanged/green. |
| **AC6** | `cell_kind` stays observed-only; no new taxonomy; no Demo 0. | **Guardrail, not a build target.** Verify the diff introduces no new `cell_kind` value, no Demo-0-shaped artifact, no new `status:*` label. |

## The α prompt

```
Branch: cycle/626

You are α, implementing cnos#626: "arch(cds): isolate the dispatch
cell-worker from substrate/agent identity + hub state." Read the full issue
body FIRST, including the "Operator directive — run #626 as a supervised
design-first cell" section and the "Authority — A2/CAP (narrowed for a
design-first cell)" section near the bottom (dated 2026-07-08) — these are
authoritative and narrow your scope. Read
.cdd/unreleased/626/gamma-scaffold.md (this scaffold) in full before
starting.

## The operator's bar (binding, restated)

This is a supervised design-first cell. You are authorized under Authority
A2/CAP to proceed WITHOUT asking on: in-scope CI, fixtures, prompt/prose,
PR body, and receipt honesty. But: "if the safe mechanism is not clearly
bounded, STOP at doctrine + implementation plan and hand back to the
operator rather than force an implementation." This scaffold has already
done that analysis for you (see "Recommended design" below) — your job is
to EXECUTE the bounded part and WRITE UP (not build) the unbounded part.
Do not force AC3/AC4 into this cell. No Demo 0. Do not touch #633 unless
you find, with evidence, that it is directly required (it should not be).

## What is bounded and what you implement

γ's analysis (gamma-scaffold.md "Recommended design") found ONE of the
issue's three candidate boundary moves cleanly implementable in this
single cell, without new infra, without touching the live checkout shape
of a production wake: **drop the full six-step sigma activation for
cell-execution cognition; load Kernel + CA skills only.**

1. Edit `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`
   §"Identity and activation" (currently lines 99-105). Replace the
   "Activate per activate/SKILL.md: Kernel -> CA skills -> Persona ->
   Operator -> hub state -> identity confirmation" instruction with
   doctrine stating: cell-execution cognition (delta/gamma/alpha/beta)
   loads Kernel + CA skills ONLY (soul, per activate/SKILL.md's own
   layering rule at its Sec 2.1 -- "Kernel and CA skills are
   substrate-independent... Persona and Operator are per-hub"). It does
   NOT load Persona/Operator/hub-state. The wake's own minimal
   self-identification (which wake, which protocol -- needed for claim
   comments and audit trails) is MANIFEST-DECLARED (the wake's name +
   the `{agent}` variable + bot-account binding already supplied by the
   renderer), not file-discovered by reading Persona.md. Do NOT read
   `.cn-{agent}/` or thread surfaces at any point in cell execution.
   Preserve the existing "you do NOT attach to a channel like the admin
   wake does" sentence -- it already states part of this boundary; fold
   your edit in coherently rather than duplicating it.

2. Add a new subsection to `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`
   section 9 (sibling to Sec 9.10 "resumed-from-changes shape" and Sec 9.11
   "resumed-from-mechanical-reversion shape" -- follow their heading/prose
   style), naming this AC1's doctrine: the cell carries no substrate
   identity and no hub-state scope; the wake/substrate owns identity,
   channel logs, and reporting. Cross-reference Sec 9.2's existing five-input
   contract (claimed-issue-number, protocol identifier, current main SHA,
   wake run id, package-runtime context) and state EXPLICITLY that the
   absence of Persona/Operator/hub-state from that list is intentional,
   not an oversight -- this is the doctrine-naming AC1 requires. Cite
   cnos#626 as the issue landing this subsection, mirroring how Sec 9.11
   cites cnos#630.

3. Re-render BOTH mirrors of the edited prompt body -- they are byte-identical
   copies embedded by the renderer, never hand-edited:
   - `cn install-wake cds-dispatch --out src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`
   - `cn install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml`
   Confirm `.github/workflows/install-wake-golden.yml`'s sha256-diff check
   would pass (the two rendered files' cds-dispatch section must match
   post-render; run the golden CI job locally if you have the tooling, or
   diff by hand and state which lines changed and why in self-coherence.md).

4. Do NOT touch: the write-fence code (.github/workflows/cnos-cds-dispatch.yml
   lines ~358-439, "Write fence -- dispatch_activation_log_write_violation"),
   `.cn-sigma/` contents, the checkout step / any sparse-checkout mechanism,
   `activate/SKILL.md` itself (you are narrowing WHO must run the full
   procedure, not changing the procedure), any FSM/transitions.json file,
   #633.

## What you write up but do NOT build (AC3/AC4)

Write a clearly labeled section in self-coherence.md (e.g. "## AC3/AC4 --
doctrine + plan, deferred to operator") that:
- restates the two candidate mechanisms from gamma-scaffold.md's
  "Recommended design" table (sparse-checkout / worktree isolation; moving
  .cn-sigma out of the repo) with their tradeoffs;
- states plainly that this cell does NOT implement either, per the
  operator's own STOP conditions ("isolating .cn-sigma breaks code/issue
  access"; "new infra or permissions required"; "removing full activation
  weakens skill loading");
- names what a follow-on cell would need to prove AC3 safely (e.g. a
  regression matrix over every existing cell type's checkout-path needs
  before attempting sparse-checkout on the live production wake).
This is a REQUIRED deliverable of this cell, not optional narrative -- the
operator's design-first bar treats "doctrine + plan" as a first-class
output, equal in weight to code, when the safe mechanism isn't bounded.

## Guardrails (binding -- violating any of these is a scope violation)

- No Demo 0.
- Do not touch #633 unless you find concrete evidence it is directly
  required -- if you think it might be, STOP and name the evidence in
  self-coherence.md rather than editing it.
- Do not remove, weaken, or bypass the write-fence
  (dispatch_activation_log_write_violation) -- AC4 is gated on AC3, which
  this cell does not implement.
- Do not change `activate/SKILL.md`'s six-step procedure itself -- other
  consumers (the admin wake, true agent-at-a-hub activations) still need
  the full procedure unchanged.
- Never hand-edit `.github/workflows/cnos-cds-dispatch.yml` or the golden
  fixture directly -- always regenerate via `cn install-wake cds-dispatch`.
- No new `status:*` label, no new `cell_kind`, no FSM/transitions.json
  change.
- Zero Go source changes expected for the bounded scope; if you find you
  need one, STOP and name exactly why in self-coherence.md before
  proceeding -- that is itself evidence the mechanism was not as bounded
  as this scaffold assessed.

## When done

Write self-coherence.md documenting: (a) the exact diff to
cds-dispatch/SKILL.md + delta/SKILL.md and why it discharges AC1 and AC2's
first disjunct; (b) confirmation both rendered mirrors were regenerated and
match; (c) the required AC3/AC4 doctrine+plan section (see above); (d)
confirmation you touched no Go source, no FSM table, no #633, no write-fence
code, no `.cn-sigma/` contents. Commit, push to cycle/626, append your
review-readiness signal, and stop -- you do not dispatch β yourself.
```

## The β prompt

```
Branch: cycle/626

You are β, independently reviewing α's implementation of cnos#626 on
cycle/626. Read .cdd/unreleased/626/gamma-scaffold.md (this scaffold) and
.cdd/unreleased/626/self-coherence.md (alpha's round record) in full before
forming any verdict. This is a supervised design-first cell -- part of a
correct R0 is CONFIRMING that alpha correctly stopped at doctrine+plan for
the unbounded part rather than either (a) skipping it silently or (b)
over-forcing an implementation the operator's STOP conditions warned against.
Do not take alpha's self-coherence claims as verified until you have
independently walked each item below yourself.

## Independent AC walk (do not skip any -- re-derive each verdict yourself)

- AC1: read the new delta/SKILL.md Sec 9 subsection yourself. Confirm it
  explicitly states (a) the cell carries no substrate identity / hub-state
  scope, (b) the wake/substrate owns identity + channel logs + reporting.
  Confirm it cross-references Sec 9.2's five-input contract and explicitly
  names the Persona/Operator/hub-state omission as intentional. A vague or
  implied statement is a finding, not a pass.
- AC2: grep cds-dispatch/SKILL.md yourself for "Activate per" / "Persona"
  / "Operator" / "hub state" in the "Identity and activation" section --
  confirm the six-step full-activation instruction is gone or narrowed to
  Kernel+CA-skills-only for cell execution. Confirm BOTH rendered mirrors
  (golden fixture + live workflow) were regenerated via `cn install-wake
  cds-dispatch` (not hand-edited -- diff the two files against each other;
  they must match byte-for-byte in the cds-dispatch prompt section).
  Confirm `activate/SKILL.md` itself is UNCHANGED (grep the diff -- if
  alpha touched it, that is a scope violation requiring a finding).
- AC3/AC4: confirm alpha did NOT implement a checkout-isolation mechanism
  or touch the write-fence code. Confirm self-coherence.md carries the
  required "doctrine + plan, deferred to operator" section per the alpha
  prompt's explicit requirement, and that it names both candidate
  mechanisms with tradeoffs (not a one-line punt). If alpha instead built
  a sparse-checkout mechanism or removed/weakened the write-fence, this is
  a BLOCKING finding regardless of implementation quality -- it violates
  the operator's explicit STOP condition and the AC3-before-AC4 gate.
- AC5: run (or read evidence of) the `install-wake-golden` workflow's
  logic yourself -- confirm the golden/live sha256 match and renderer
  idempotence would hold. Confirm `#516`/`#524` dispatch guard scripts
  (check-dispatch-repair-preflight.sh, check-dispatch-closeout-integrity.sh)
  are unaffected (grep the diff -- it should not touch them). Confirm zero
  Go source files changed (`git diff --stat cycle/626` against base --
  the diff should be prose/markdown only, plus the two regenerated YAML
  mirrors).
- AC6: confirm no new `cell_kind`, no new `status:*` label, no
  transitions.json change, no Demo-0-shaped artifact.

## Guardrail verification (independent of the AC table)

- Confirm #633 (src/packages/cnos.cds/skills/cds/fsm/transitions.json rule
  ordering) is untouched, unless alpha names concrete evidence it was
  directly required -- if so, scrutinize that evidence hard; the operator
  explicitly said not to touch it.
- Confirm the write-fence step
  (.github/workflows/cnos-cds-dispatch.yml, "Write fence --
  dispatch_activation_log_write_violation") is byte-identical to main
  outside of line-number drift caused by the Identity-and-activation edit
  above it.
- Confirm `.cn-sigma/` contents are untouched (git diff --stat should show
  no .cn-sigma/ paths).

## Verdict

Write .cdd/unreleased/626/beta-review.md with a full R[N] section: outcome
per AC, any findings (with severity), and `verdict: converge` or
`verdict: iterate`. A cell that correctly ships AC1+AC2 and correctly
STOPS at doctrine+plan for AC3/AC4 (per the operator's own design-first
bar) is a legitimate `converge` -- do not penalize alpha for not building
AC3/AC4 if the STOP was well-reasoned and well-documented; DO flag it as a
blocking finding if alpha either silently dropped AC3/AC4 with no writeup,
or forced an implementation the scaffold and operator both flagged as
unbounded. Commit, push, and stop -- you do not re-dispatch alpha yourself.
```

## Implementation contract (pinned by δ; α MUST NOT improvise)

| Axis | Pinned value |
|---|---|
| Language | N/A for the bounded scope — Markdown/prose skill-file edits only (`cds-dispatch/SKILL.md`, `delta/SKILL.md`). No Go source changes expected. If α's design-call analysis surfaces a genuine need for a Go change, that itself falls outside this cell's bounded scope (see α prompt guardrails) and must be flagged, not silently absorbed. |
| CLI integration target | `cn install-wake cds-dispatch` (existing renderer, invoked twice — once per `--out` target — to regenerate the golden fixture and the live workflow after the SKILL.md prose edit). No new CLI verb. The deferred AC3 mechanism, if the operator later authorizes it, would target the same renderer (`src/packages/cnos.core/commands/install-wake/cn-install-wake`) with a new manifest field — explicitly out of scope here. |
| Package scoping | `cnos.cds` (`src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` + its golden fixture) and `cnos.cdd` (`src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`). `cnos.core`'s `activate/SKILL.md` is read-only cross-reference — not edited. The rendered `.github/workflows/cnos-cds-dispatch.yml` is a generated artifact regenerated in place, not a package source file. |
| Existing-binary disposition | Preserve — no binary/CLI behavior changes. `cn install-wake` is invoked exactly as documented; no flags or defaults change. |
| Runtime dependencies | None new. |
| JSON/wire contract | N/A — no JSON/wire-format surface touched. The FSM's `FactSnapshot`/`Decision`/`ScanReport` shapes (relevant to #630/#633, not this issue) are untouched. |
| Backward-compat invariant | The admin wake (`agent-admin`) and any true agent-at-a-hub activation continue to use the full unmodified six-step `activate/SKILL.md` procedure — this cell narrows only *who is instructed to run it* (the `cds-dispatch` cell-execution path), never the procedure's content. `install-wake-golden` CI must stay green (proves renderer determinism + golden/live sync unchanged post-edit). `#516`/`#524` dispatch guard scripts pass unchanged. |

No TBD rows. Every axis above is pinned from direct repo-convention evidence gathered while reading the current fusion (Key findings 1–4) and the source-of-truth table; the one genuinely undecidable axis — what Go/infra change AC3's checkout-isolation mechanism would eventually require — is explicitly NOT pinned here because it is explicitly NOT this cell's scope; α documents the open question in self-coherence.md per the α prompt rather than guessing a value for it.

## Scope guardrails (binding)

- No Demo 0.
- Do not touch #633 (FSM rule-ordering) unless α finds and documents concrete evidence it is directly required — expected answer: not required.
- Do not remove, weaken, or bypass the `dispatch_activation_log_write_violation` write-fence. AC4 is hard-gated on AC3; AC3 is not implemented this cycle, so AC4 cannot be either.
- Do not modify `.cn-sigma/` contents, the checkout step, or introduce any sparse-checkout/worktree-isolation mechanism — that is candidate move (1)/(2), explicitly deferred to operator per "Recommended design."
- Do not modify `activate/SKILL.md`'s six-step procedure — other consumers still need it unmodified.
- Preserve FSM ownership, finalizer (`cn cell finalize`), reconciler (`cn issues fsm scan`), and all existing dispatch guards (`#516`, `#524`) green — AC5.
- `cell_kind` stays observed-only; no new taxonomy; no new `status:*` label — AC6.
- Any rendered-workflow change goes through `cn install-wake cds-dispatch`, never a hand-edit — the `install-wake-golden` CI gate enforces this.

## Friction notes

1. **AC2's "and/or" wording is the load-bearing escape valve for a responsible R0.** The issue's own AC2 text — "the dispatch cell-worker no longer performs a full sigma activation **and/or** no longer has `.cn-{agent}/` in its working scope" — means the first, cleanly-bounded disjunct alone is a legitimate, honest AC2 pass. Without this reading, a naive R0 would feel compelled to force the checkout-isolation mechanism (candidate move 1) into this cell despite the operator's explicit STOP conditions warning against exactly that. β should verify this reading is the one α actually applied, not an after-the-fact rationalization for skipping harder work — the "Recommended design" table above shows the bounded/unbounded split was reasoned *before* dispatch, not invented by α to avoid AC3.

2. **This cell is self-referential: it analyzes and potentially changes the very mechanism it is itself running under.** γ (this scaffold) is being authored under bootstrap-δ, inside a checkout that has `.cn-sigma/` present, using the exact `activate`-adjacent doctrine surfaces this issue proposes to narrow. If a future wake-invoked-δ firing of `cds-dispatch` runs this very cell (or a repair round of it), it would be executing under the OLD fused-activation prompt until α's edit merges — there is no way for this cell to "run under its own fix" mid-cycle. This is not a blocking problem (the doctrine change is inert prose until the next wake firing reads it post-merge), but it is worth naming explicitly per the task's own instruction: a cell that changes its own execution contract cannot self-validate the new contract from inside the old one. β's independent AC walk is the closest available check (reading the new prose critically, not just diffing it), not a live re-run under the new contract.

3. **The "42 tracked files / 40 channel-log files" the issue cites has already drifted to 43/41 by scaffold time** (one more of each, one week later) — this is expected, ongoing admin-wake activity, not evidence against the issue's premise; if anything it strengthens it (the surface keeps growing precisely because nothing currently isolates it from dispatch checkouts).

4. **The precise identity of "#614" required correction during scaffolding.** The issue text says "Origin: RCA of #614 (write-fence false-positive)"; `gh issue view 614` shows #614's actual title is an unrelated CELL-KINDS/planning-cell-kind amendment (closed, merged). The write-fence RCA is documented in `.cdd/unreleased/614/gamma-closeout.md` line 31 as an incident that occurred during a *prior firing* of cycle #614 itself (not what #614 is "about"), and the actual fix issue is #625 ("wake: fix write-fence false-positive on concurrent admin-wake log commits," merged). The "Family" line above reflects this corrected lineage; α/β should not be confused if they independently `gh issue view 614` and find CELL-KINDS content — that is expected and does not indicate a wrong issue number in this scaffold.

5. **`.cn-sigma/README.md`'s own documented design already half-anticipates the fix.** The README states identity/peer-graph files (PERSONA.md, OPERATOR.md) live at Sigma's *home hub* (`cn-sigma`), not in `.cn-sigma/` at cnos — "There is no local cache; this directory is for per-context state only." This means candidate move (2) (moving the hub out of the repo) is arguably *already true for identity*, and only per-context state (mailbox, logs) remains in-repo by design. This nuance doesn't change the bounded/unbounded split above (the per-context state is still the thing AC3 wants isolated), but α/β should not mistake the existing README framing for evidence that AC2's `hub-state` step is harmless — the hub-state survey step in `activate/SKILL.md` §2.1 item 5 explicitly walks memory/thread surfaces, which is a different (and still fused) concern from the Persona/Operator identity-pair resolution the README addresses.
