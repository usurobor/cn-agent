---
cycle: 496
parent_issue: cnos#496
umbrella: cnos#495 (Sub 1 — first concrete enforcement of the long-arc partition)
sibling_held: cnos#497 (artifact-root design-first; do NOT pre-empt)
master_tracker: cnos#467 (closed completed; this is a post-wave follow-up)
cycle_branch: cycle/496
base_main_sha: 3f57210d95f765ce1884e0f2d6a0868e25b7e1b0
head_sha_at_scaffold: (filled by scaffold-landing commit)
self_application_paradox: true (mirrors cnos#487 Sub 5C — bootstrap-δ-claimed so live cds-dispatch won't pick it up)
role: γ
authored_by: γ@cdd.cnos (bootstrap-δ via δ-interface session)
date: 2026-06-24 (UTC)
output_contract: γ-scaffold + α prompt + β prompt
ac_set: cnos#496 body ACs 1–7, with AC7 superseded by [comment #4792858087](https://github.com/usurobor/cnos/issues/496#issuecomment-4792858087) (local-write fence guardrail)
constraints_inherited:
  - T-486-1 (variable consistency table)
  - T-486-7 (R[N] off-by-one prevention; δ-direct R1 pattern available)
  - T-486-12 (operator-final-read defense-in-depth, P1; promoted in cycle/487)
  - T-486-15 (predecessor-closeout-reading)
  - T-487-1 (variable consistency table extended to ALL surfaces, not just structural fields)
artifact_tree_doctrine: Model B (.cdd/unreleased/{N}/) preserved; cnos#497's Model A vs Model B decision is HELD and explicitly NOT pre-empted by this cycle
---

# γ-scaffold — cnos#496 Sub 1 (cnos#495 umbrella): prevent package dispatch wakes from writing activation logs

## §1. Parent issue framing

**The empirical incoherence (2026-06-24).** The `cnos-cds-dispatch` wake's provider prompt at `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` explicitly names `.cn-sigma/` under "Disallowed surfaces" and explicitly states *"the dispatch wake does NOT write channel entries; channel logs are the admin wake's writer-locality surface."* The provider manifest at `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` likewise names admin-only writer-locality semantics. **The model wrote `.cn-sigma/logs/20260624.md` four times anyway** — selector-scan no-op entries at ~05:55Z, 08:45Z, 12:38Z, and 13:58Z. The activation log is a mixed surface today: both `cnos-agent-admin` (legitimate writer) AND `cnos-cds-dispatch` (illegitimate writer) commit to it.

**The cycle's goal.** Install a **mechanical guard** that prevents this. Prompt-only prohibition has been empirically falsified. The guard must scope itself to *this wake's local writes* (staged + committed paths under `.cn-{agent}/logs/`) — not to remote-state delta on `main`, which would false-positive on legitimate concurrent admin-wake writes.

**The cycle's framing.** This is the **first concrete enforcement** of the long-arc partition declared in cnos#495's umbrella amendment ([comment #4792969173](https://github.com/usurobor/cnos/issues/495#issuecomment-4792969173)):

> Use intelligence where meaning is unresolved.
> Use mechanics where the rule is known.

> Identity can be Sigma.
> Ownership is package-local.
> Memory is summarized upward, not dumped sideways.

The cycle is narrow on purpose. Per operator's explicit guardrail: *"Keep cnos#496 narrow. It should not try to move the whole orchestration stack into Go. It should prove one mechanical boundary: package dispatch wakes cannot write activation logs."* This cycle installs THAT boundary and nothing more. The broader migration (claim sequence to Go, label transitions to Go, δ routing to Go, return-token parsing to Go) is the umbrella's later-cycle work. Sub 2 (admin dispatch-summary) and cnos#497 (artifact-root decision) are HELD; this cycle does not move them.

**Self-application paradox handling.** The issue whose work changes dispatch-wake behavior cannot be claimed by the live dispatch wake (the wake would write `.cn-sigma/logs/` during the very cycle that's supposed to repair it). Mirror cycle/487 pattern: cnos#496 is bootstrap-δ-claimed (`status:in-progress`); the live wake's selector requires `status:todo`; the issue is invisible to the live wake. The bootstrap-δ claim is recorded in [comment #4792872514](https://github.com/usurobor/cnos/issues/496#issuecomment-4792872514).

## §2. OG (Objective Gap) framing

**Today (incoherent):** every dispatch-wake firing — real claim OR no-op — runs `activate + attach` as Tier-1 skills. `attach`'s channel-append behavior appends an entry to `.cn-sigma/logs/YYYYMMDD.md` and pushes. The provider prompt forbids this; the model does it anyway. Result: package-runtime telemetry pollutes agent activation memory at every firing.

**Tomorrow (coherent):** dispatch-wake firings produce zero commits to `.cn-{agent}/logs/`. The renderer refuses to materialize a package-dispatch-class wake provider whose declaration would result in channel-append. A mechanical post-run write fence inspects the wake's own staged + committed paths and fails the workflow with `dispatch_activation_log_write_violation` if anything under `.cn-{agent}/logs/` was touched. No-op firings surface in the GitHub Actions job summary (UI-visible, structured, ephemeral) — NOT in the activation log (durable agent memory).

**The load-bearing invariant.** **Package-dispatch wakes are non-writers of `.cn-{agent}/logs/`.** The convention amendment names this; the provider field declares it per-wake; the renderer enforces it at install time; the workflow fence enforces it at run time. Belt and suspenders. The empirical proof-point demands both: the model already proved it ignores prompt-only declarations; render-time + run-time mechanical checks are the actual enforcement.

## §3. Surface inventory

Every surface this cycle touches, with the one-line change:

| # | Surface | Change |
|---|---|---|
| 1 | `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` | Amend with explicit "Writer ownership" section: admin-wake = sole writer; package-dispatch wakes = non-writers; carve-out to §0 Writer Locality invariant naming this exclusion. |
| 2 | `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` | Add `"activation_log_writer": false` field. Update `responsibilities[]` to reference the mechanical guard. Note: this manifest field is the source-of-truth for the renderer's guard decisions. |
| 3 | `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` | Reference the mechanical guard in the disallowed-surfaces section. State explicitly that the prohibition is mechanically enforced post-run (not just prompt-asserted). |
| 4 | `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` | Re-rendered to (a) NOT load `attach` for this dispatch wake, (b) include the post-run write fence step, (c) emit `activation_log_writer: false` provenance. New sha256. |
| 5 | `.github/workflows/cnos-cds-dispatch.yml` | Re-rendered to match the new golden. New sha256. |
| 6 | `src/packages/cnos.core/commands/install-wake/cn-install-wake` | Add: (a) parse `activation_log_writer` field from manifest (default `true` for backward compat; package-dispatch wakes must declare `false`); (b) refuse render when `wake_class: package-dispatch` (or `admin_only: false` && `role: dispatch`) AND `activation_log_writer: true` (mis-declaration); (c) when `activation_log_writer: false`, skip emission of `attach` in the rendered skills load list; (d) append the post-run write-fence step to the rendered workflow. New exit code `4` for "activation-log-writer mis-declaration." |
| 7 | `src/packages/cnos.core/skills/agent/attach/SKILL.md` | NO CHANGE in this cycle (operator KISS preference D1 (c): dispatch wakes don't load `attach`; `attach` itself stays unitary). If implementer discovers a forced split (e.g., `cn` substrate needs partial channel binding for dispatch claim sequence in a future cycle), THAT is a separate doctrine cycle. |
| 8 | `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md` | Add `dispatch_activation_log_write_violation` to the failure-class taxonomy. Reference this cycle's mechanical guard. (NOT a doctrine change; the failure class is the formal error name the guard emits.) |
| 9 | `.github/workflows/install-wake-golden.yml` | Add two new steps: (a) **positive write-fence test** — render cds-dispatch golden + simulate a no-op firing in a fixture environment; verify zero local writes under `.cn-{agent}/logs/`; (b) **negative write-fence test** — render a mis-declared fixture; verify rendering exits with code 4. Optional: (c) **false-positive resistance test** — run the positive test alongside a parallel fixture that writes `.cn-sigma/logs/` legitimately (simulating admin wake firing); verify the dispatch fence does NOT false-fail. |
| 10 | `src/packages/cnos.core/commands/install-wake/test-fixtures/log-writer-misdeclaration/` | New fixture: a wake-provider manifest declaring `role: dispatch + activation_log_writer: true` (mis-declaration). Used by the negative test in #9. Mirrors the cycle/485 declaration-only fixture pattern. |

## §4. AC interpretation table

The body of cnos#496 declares ACs 1–7. AC7 is superseded by the post-filing amendment in [comment #4792858087](https://github.com/usurobor/cnos/issues/496#issuecomment-4792858087). β should treat the amendment as the binding form.

### AC1 — Doctrine amendment

**Invariant:** `AGENT-ACTIVATION-LOG-v0.md` names writer ownership explicitly.

**Mechanical oracle:** the convention doc post-amendment contains a "Writer ownership" section that:
- Names admin-wake as the sole writer of `.cn-{agent}/logs/`.
- Names package-dispatch wakes as non-writers.
- Carves out the §0 Writer Locality invariant with the package-dispatch exclusion.

**Surfaces touched:** #1.

**Interpretation note:** the existing §0 says *"Every body writes only to its own repo."* That invariant is about cross-repo direction (foreign-vs-home), NOT about same-repo cross-wake ownership. The amendment must distinguish: same-repo writer-locality has historically been "the body that owns this repo writes here" — but now we have multiple wakes-as-bodies in the same repo with different ownership. The cleanest framing: add a new §0.1 (or §6.1) named "Wake-class writer ownership (same-repo)" that names the partition.

### AC2 — Provider field shipped

**Invariant:** wake-provider manifests support an `activation_log_writer: bool` field; cds-dispatch declares `false`; cnos.core agent-admin declares `true`.

**Mechanical oracle:**
- `jq -r '.activation_log_writer' src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` returns `false`.
- `jq -r '.activation_log_writer' src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` returns `true`.

**Surfaces touched:** #2, plus agent-admin's manifest (read-only check; agent-admin must declare `true` explicitly so it's not relying on the default).

**Interpretation note:** the field's default-when-absent is `true` (backward-compat with existing agent-admin behavior). Package-dispatch wakes must declare `false` explicitly. The renderer's mis-declaration check (AC4) catches package-dispatch wakes that forget to set it OR set it to `true`.

### AC3 — Dispatch wake does not run attach/channel-append

**Invariant:** when a wake's provider declares `activation_log_writer: false`, the rendered workflow does not invoke `attach`'s channel-append behavior. Operator KISS form (D1 (c)): the rendered workflow does not load `attach` at all.

**Mechanical oracle:** rendered `cnos-cds-dispatch.yml` does NOT reference `attach` (no `cn attach` invocation; no `attach` in the loaded skills section; no channel-append prose in the prompt's loaded-skills enumeration). Verified by `grep -c 'attach' .github/workflows/cnos-cds-dispatch.yml` returning zero (or only references that are documentation/comment-form, not invocation-form).

**Surfaces touched:** #4, #5, #6.

**Interpretation note:** the rendered workflow's loaded-skills list is shaped by the renderer reading the provider's declared `loaded_skills` (or equivalent field) AND the `activation_log_writer` field. When the renderer sees `activation_log_writer: false`, it omits `attach` from the loaded skills regardless of what the provider's `loaded_skills` says. (Defensive: the provider should already NOT list `attach`, but the renderer enforces it anyway.)

### AC4 — Renderer guard

**Invariant:** `cn install-wake` refuses to materialize a package-dispatch-class wake provider that mis-declares `activation_log_writer: true`. Refuses to materialize a wake whose `loaded_skills` includes `attach` when `activation_log_writer: false`.

**Mechanical oracle:**
- Run `cn install-wake test-fixtures/log-writer-misdeclaration --out /tmp/render.yml --activation-state-override live`; exit code 4; stderr names the mis-declaration field.
- Run `cn install-wake cds-dispatch --out /tmp/render.yml`; exit code 0 (legitimate render).

**Surfaces touched:** #6, #10.

**Interpretation note:** D3 fork — render-time refusal (a) is preferred over runtime-only enforcement (b) because mis-declared providers shouldn't reach the substrate at all. Belt-and-suspenders applies: render-time refuses the mis-declaration; run-time fence catches drift if a hand-edited workflow bypasses the renderer.

### AC5 — cds-dispatch scrub

**Invariant:** live cds-dispatch firings produce zero commits to `.cn-{agent}/logs/`.

**Mechanical oracle:**
- After cycle/496 merges: trigger a no-op cds-dispatch firing (scheduled sweep with empty selector).
- Observe `git log -- .cn-sigma/logs/` shows no new commits attributable to the cds-dispatch firing.
- Repeat for a synthetic real-claim firing (e.g., a smoke cell similar to cnos#491).

**Surfaces touched:** all of #2 through #6.

**Interpretation note:** this AC is empirically verifiable post-merge by watching the live wake's behavior. β should mark it provisionally green at PR review (CI guard satisfies the surrogate); the operator's post-merge check provides ground truth.

### AC6 — Historical evidence preserved

**Invariant:** the 2026-06-24 mixed log entries are NOT rewritten.

**Mechanical oracle:** `git log -- .cn-sigma/logs/20260624.md` post-cycle shows the original commits intact (4 cds-dispatch no-op entries + the admin entries). No `git filter-branch`, no `git rebase -i`, no rewrite of any commit touching `.cn-sigma/logs/20260624.md`.

**Surfaces touched:** none (the AC is a *constraint*; no positive action).

**Interpretation note:** the historical entries are the proof-point for *why* the mechanical guard is necessary. Rewriting them would erase the evidence motivating the cycle. Cycle/496's PR body cites these entries by SHA as the empirical witness.

### AC7 — Mechanical write fence for package dispatch wakes (THE CORE DELIVERABLE — AMENDED FORM)

**Invariant:** a package dispatch wake must have a mechanical guard proving it did NOT write `.cn-{agent}/logs/` during the run.

**Amended v0 acceptable proof (from [comment #4792858087](https://github.com/usurobor/cnos/issues/496#issuecomment-4792858087)):**

The write fence must prove **the dispatch wake itself did not create, modify, commit, or push paths under `.cn-{agent}/logs/`**.

Preferred v0 proof: **inspect the dispatch run's local diff / staged paths / commit paths before push.** Acceptable mechanisms:
- `git status --porcelain -- .cn-*/logs/` after the wake's work phase; expect empty.
- `git diff --name-only HEAD@{1} HEAD -- .cn-*/logs/` for any commit attributed to this wake's run; expect empty.
- Pre-push hook or workflow step that scans the local commit graph this wake created for paths under `.cn-*/logs/`; fail if any present.

**Explicit guard against false-positives:** the fence MUST NOT fail merely because another wake (e.g., the agent-admin wake firing concurrently on its own concurrency group) wrote `.cn-{agent}/logs/` to remote `main` during this run. Concurrency groups separate admin (`agent-admin-sigma`) and dispatch (`cds-dispatch-sigma`); they can run in parallel; remote-state delta is not attributable.

**Mechanical oracle:** `install-wake-golden.yml` has a step that:
1. Runs a synthetic no-op fixture firing of cds-dispatch.
2. The fixture explicitly: runs a parallel synthetic admin-wake fixture that DOES write to `.cn-sigma/logs/YYYYMMDD.md` on a fixture-only branch during the dispatch run.
3. The dispatch wake's fence inspects only its OWN staged + committed paths.
4. Verifies zero local file changes under `.cn-{agent}/logs/` attributable to the dispatch run.
5. Verifies the parallel admin-wake's writes do NOT cause the dispatch fence to fail (the false-positive guard).

A second fixture (deliberately-misbehaving) injects an `attach`-channel-append into a fixture dispatch provider:
1. The dispatch wake attempts to write `.cn-sigma/logs/YYYYMMDD.md`.
2. The fence detects the local write.
3. Fence fires with `dispatch_activation_log_write_violation`.

Both fixtures pass green in CI.

**Surfaces touched:** #4, #5, #6 (renderer appends the fence step), #8 (dispatch-protocol skill declares the failure class), #9 (CI fixtures), #10 (mis-declaration fixture).

**Interpretation note (LOAD-BEARING):** the local-vs-remote distinction is small in code, large in correctness. A naive `git diff main@before main@after -- .cn-*/logs/` check would false-positive on every legitimate concurrent admin-wake write. The fence must be scoped to *this wake's local writes only* — `git status --porcelain` on the runner's working tree, OR `git log $LAST_FETCH_SHA..HEAD --author=<this-wake's-author>` against the local commit graph this wake created. NOT a comparison of remote-state delta.

## §5. Design forks

Three forks the cycle must resolve. Operator preferences noted where stated; rationale follows the decision.

### D1 — Attach skill gating mechanism

How does the rendered workflow avoid `attach`'s channel-append behavior for dispatch wakes?

- **(a) Split attach** into `attach-identity` (identity-load; safe for any wake) + `attach-channel` (channel-append; admin-only). Dispatch loads identity only.
- **(b) Keep attach unitary**; add an `attach.channel_append: bool` flag in the skill body; provider declares wake-class default; skill body branches on flag.
- **(c) KISS** — dispatch wakes don't load `attach` at all. Only `activate` (identity load). The activation-channel binding work is admin-wake-only.

**Recommendation: (c).** Operator-stated KISS preference; smallest blast radius; no `attach` source modification; ownership clean. The activate skill's identity-load is sufficient for dispatch attribution (commit author, identity surface). Channel binding is admin-only.

**When to fall back to (b):** if a future package-dispatch wake genuinely needs partial channel-binding behavior (e.g., needs to record a structured per-firing event without it being a "channel entry" in the activation-log sense). Treat that as a separate doctrine cycle, NOT this cycle.

**Reject (a):** over-engineering for v0; introduces a skill-name expansion that touches every consumer of `attach`; the empirical drift today is solved by (c).

### D2 — Write-fence implementation

How does the post-run guard mechanically detect a write to `.cn-{agent}/logs/`?

- **(a) Workflow-level step.** Final step in the rendered workflow runs `git status --porcelain -- '.cn-*/logs/'` after the wake's work phase; non-empty → fail the job with `dispatch_activation_log_write_violation`.
- **(b) Git hook.** Renderer installs a `pre-push` hook on the runner that scans the local commit graph for `.cn-{agent}/logs/` paths; vetos the push if found.
- **(c) `cn` substrate command.** New `cn dispatch verify-fence` subcommand (bash addition to `cn-install-wake`'s sibling or a new `cn-dispatch-verify-fence` script); invoked by the rendered workflow's final step.

**Recommendation: (a).** Per operator's "keep cnos#496 narrow" guardrail: this cycle should not move the orchestration stack into Go (or grow new `cn` subcommands as a side effect). Workflow-level step is the smallest viable mechanical fence. The check is ~5 lines of bash in the rendered workflow.

**When to fall back to (c):** if multiple subsequent cycles add similar mechanical guards (artifact-path enforcement, label-set verification, etc.), consolidate them into a `cn dispatch verify` umbrella command. That's a future-cycle decision based on empirical pattern emergence. Today, ONE guard → workflow step.

**Reject (b):** git hooks on ephemeral CI runners are fragile (hook scripts persist only for the run; no audit trail outside the workflow log). Workflow-level steps are first-class CI surface — visible in the Actions UI, persist in the run log, easy to debug.

### D3 — Renderer guard scope

When does the renderer reject a mis-declared package-dispatch provider?

- **(a) Render-time refusal.** `cn install-wake` exits nonzero (new exit code 4) if `wake_class == package-dispatch` AND `activation_log_writer == true` (or `loaded_skills` contains `attach` when `activation_log_writer == false`).
- **(b) Render-with-warning + runtime fence.** Renderer warns to stderr but still produces output; runtime fence catches actual violations.

**Recommendation: (a).** Fail fast: a mis-declared provider shouldn't reach the substrate at all. Run-time fence (AC7) is belt-and-suspenders, NOT primary defense. Two-layer enforcement: render-time validation + run-time fence.

**Reject (b):** prompts the same prompt-prohibition failure mode (warning ignored; substrate carries broken declarations). Hard refusal is the correct primary defense.

### D4 — No-op evidence destination (lower-stakes; operator already directional)

Where do no-op dispatch firings surface, post-cycle?

- **(a) GitHub Actions job summary only** (workflow log + run UI; no commit, no PR comment, no issue comment).
- **(b) Issue comment** on a tracking issue (lifecycle-level no-op log).
- **(c) Workflow log only** (no UI surface beyond raw logs).

**Recommendation: (a).** Operator-stated direction in cnos#495 body: *"workflow/job summary only by default."* Operator-stated in long-arc partition: *"workflow/job summary."* Job summary is structured + UI-visible without polluting durable activation memory.

**When to fall back to (b):** when admin-wake's `class: dispatch-summary` entry lands (Sub 2; HELD). Until then, job summary is sufficient.

**Reject (c):** invisible to operator without diving into raw logs; operator-final-read defense-in-depth (T-486-12 P1) is weakened.

## §6. α prompt (verbatim, ready to copy-paste)

```text
You are α@cdd.cnos, running as the implementer for cycle/496.

Branch: cycle/496 (at base main SHA 3f57210d).
Parent issue: cnos#496.
Umbrella: cnos#495 (Sub 1 — first mechanical enforcement of the long-arc partition).
Scaffold: /home/user/cnos/.cdd/unreleased/496/gamma-scaffold.md (READ THIS FIRST in full).

## What you implement (per γ-scaffold §3 Surface inventory)

Implement R0 against the 7 ACs (with AC7 superseded by the amendment in [comment #4792858087](https://github.com/usurobor/cnos/issues/496#issuecomment-4792858087)). Apply γ's recommended fork choices: D1 (c), D2 (a), D3 (a), D4 (a). Do NOT deviate without recording an FN explaining why.

## Implementation order (suggested; minimizes inter-step risk)

1. **Convention amendment (AC1) — `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`.** Add §0.1 (or equivalent location) "Wake-class writer ownership (same-repo)" naming admin-wake as sole writer of `.cn-{agent}/logs/`; naming package-dispatch wakes as non-writers; carving out the §0 Writer Locality invariant explicitly. Cross-reference cnos#495 umbrella + cnos#496 enforcement implementation + the long-arc partition framing.
2. **Provider field declaration (AC2) — wake-provider manifests.**
   - Add `"activation_log_writer": false` to `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json`.
   - Add `"activation_log_writer": true` to `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` (explicit declaration; not relying on default).
   - Update cds-dispatch's `responsibilities[]` to name the mechanical guard.
3. **cds-dispatch prompt scrub (AC3 surface) — `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md`.** Update the disallowed-surfaces section to state explicitly: "The prohibition on writing `.cn-{agent}/logs/` is mechanically enforced — post-run write fence emits `dispatch_activation_log_write_violation` on breach. This is not advisory; it is structural."
4. **Renderer guard (AC4) — `src/packages/cnos.core/commands/install-wake/cn-install-wake`.**
   - Parse `activation_log_writer` field from the manifest (default `true` if absent for backward compat).
   - Detect package-dispatch wakes (`role: dispatch + admin_only: false`).
   - Refuse render when package-dispatch + `activation_log_writer: true` → exit code 4 + stderr message naming the mis-declaration.
   - Refuse render when `activation_log_writer: false` AND `loaded_skills` includes `attach` (or however the renderer enumerates loaded skills) → exit code 4.
   - When `activation_log_writer: false`, omit `attach` from the rendered skills load list in the workflow output.
   - When `activation_log_writer: false`, append the workflow-level write-fence step to the rendered output (per D2 (a)).
5. **Rendered output (AC3 + AC5) — `.github/workflows/cnos-cds-dispatch.yml` + golden.**
   - Re-render via `cn install-wake cds-dispatch --out src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` (the golden update).
   - Then render the production substrate: `cn install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml`.
   - Verify `sha256sum` matches between rendered and golden.
   - The new rendered workflow:
     - Does NOT load `attach`.
     - Has a final step (after the claim/work/push phases) running the write fence: `git status --porcelain -- '.cn-*/logs/'` returning empty; fail with `dispatch_activation_log_write_violation` if non-empty.
6. **Dispatch-protocol skill amendment (AC7 secondary) — `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md`.** Add `dispatch_activation_log_write_violation` to the existing failure-class taxonomy (alongside `dispatch_protocol_missing`, `dispatch_protocol_mismatch`, `dispatch_label_drift`, etc.). Brief 2-paragraph description: when it fires; what evidence is captured; what the runtime behavior is (fail the workflow; do NOT release the claim back to status:todo since the fence fires AFTER work has happened; surface the violation in workflow log + job summary).
7. **CI fixtures (AC7 oracle + AC4 oracle) — `.github/workflows/install-wake-golden.yml` + `src/packages/cnos.core/commands/install-wake/test-fixtures/log-writer-misdeclaration/`.**
   - Create `src/packages/cnos.core/commands/install-wake/test-fixtures/log-writer-misdeclaration/{wake-provider.json,prompt.md}` — mis-declared fixture (role: dispatch + activation_log_writer: true).
   - Extend `.github/workflows/install-wake-golden.yml`:
     - **Positive test (write-fence green-on-no-op):** render cds-dispatch golden; in a fixture sub-shell, simulate the workflow's claim/work phase doing NOTHING to `.cn-{agent}/logs/`; run the write fence; expect zero violations.
     - **Negative test (write-fence catches violation):** in a fixture sub-shell, write a synthetic line to `.cn-sigma/logs/20260624.md`; stage it; run the write fence; expect `dispatch_activation_log_write_violation` exit code.
     - **False-positive resistance test:** in a fixture sub-shell, simulate a parallel admin-wake fixture committing to `.cn-sigma/logs/YYYYMMDD.md` on a fixture-only branch; the dispatch fence (running on a SEPARATE working tree / branch) inspects only its own local writes; verify fence does NOT false-fail.
     - **Renderer refusal test (AC4):** `cn install-wake test-fixtures/log-writer-misdeclaration` exits with code 4 and stderr names the mis-declaration.
8. **Self-coherence §R0 — `.cdd/unreleased/496/self-coherence.md`.** Per cycle/487 pattern: signal review-ready; cite head SHA at signal; declare which ACs are mechanically green (CI-verified) vs which require operator-final-read (T-486-12 P1).

## Variable consistency table seed (per T-486-1 + T-487-1 expanded surface coverage)

Seed table; α extends as implementation reveals more anchor variables:

| Variable | Manifest (wake-provider.json) | Rendered workflow YAML | CI guard | Prompt prose | Smoke procedure | PR body |
|---|---|---|---|---|---|---|
| `activation_log_writer` | `false` (cds-dispatch); `true` (agent-admin) | (implicit: presence/absence of `attach` load + fence step) | (positive + negative fixtures cite it) | Prose names "mechanically enforced" | N/A (no smoke this cycle) | Cites the field by name |
| `dispatch_activation_log_write_violation` | (responsibilities[] reference) | Final step error string | Negative fixture expects this exit error | Prose names this failure class | N/A | Cites it |
| `.cn-{agent}/logs/` | (responsibilities[] disallowed surface) | Fence step's `git status -- '.cn-*/logs/'` glob | Fixtures touch `.cn-sigma/logs/20260624.md` | Prose names the surface | N/A | Cites the surface |
| Renderer exit code 4 | N/A | N/A | Negative fixture expects exit 4 | N/A | N/A | Cites the new exit code |
| Wake class (`role: dispatch` + `admin_only: false`) | Existing fields; no change to enum | (Implicit) | (Implicit) | (No prose change beyond #3) | N/A | (Cites in passing) |

T-487-1 lesson: the table MUST cover prompt prose (column 4), not just structural fields. β's R0 review will walk every column for every variable.

## Signaling β-readiness

After all 8 steps land + push, write self-coherence.md §R0, commit as `α-496 R0: <one-line summary>`, push. Then add a SECOND commit titled `α-496 R0 signal: ready for β review` (zero file changes; signal commit; matches cycle/487 pattern at `e107b7e4`).

## Re-iteration discipline

If β R0 returns iterate (operator-final-read OR β's verdict), enter R1:
- Read β's review carefully; do NOT re-derive ACs from scratch.
- R[N] section in self-coherence.md per cycle/486 T-486-7 off-by-one prevention (R1 section header at the SAME ordinal as the iterate that triggered it; R[N]'s β review is in beta-review.md's §R[N], not §R[N-1]).
- Apply δ-direct R1 pattern (T-486-7) if the iterate is narrow + mechanical (single file/line fix without architectural reframing): γ-interface may apply the fix without re-spawning α. Operator's call.

## Hard constraints

- **Do NOT modify `.cn-sigma/logs/20260624.md`** — AC6 preserves historical evidence.
- **Do NOT pre-empt cnos#497** — use `.cdd/unreleased/496/` paths; do NOT rename to `.cds/`.
- **Do NOT introduce new `cn` subcommands** — D2 (a) workflow-step v0; operator guardrail "keep cnos#496 narrow."
- **Do NOT modify `attach` skill source** — D1 (c) means dispatch wakes don't load it; the skill itself stays unitary.
- **Do NOT skip any AC** — if an AC seems redundant or out-of-scope, file an FN explaining; do NOT silently drop.
- **Commit message pattern** — `α-496 R0 <step name>: <one-line summary>` per file group; final commit `α-496 R0 signal: ready for β review`. Co-Authored-By trailer pattern matches prior cycles.
```

## §7. β prompt (verbatim, ready to copy-paste)

```text
You are β@cdd.cnos, running as the reviewer for cycle/496 R0.

Branch: cycle/496 (read α's R0 commits up to and including the "α-496 R0 signal" commit).
Parent issue: cnos#496.
Umbrella: cnos#495.
Scaffold: /home/user/cnos/.cdd/unreleased/496/gamma-scaffold.md (read in full).
α self-coherence: /home/user/cnos/.cdd/unreleased/496/self-coherence.md §R0 (read in full).
AC7 amendment (load-bearing): [comment #4792858087](https://github.com/usurobor/cnos/issues/496#issuecomment-4792858087).

## Your output

Write `.cdd/unreleased/496/beta-review.md` with §R0 review. Final verdict: **converge** OR **iterate**.

## Pre-AC checklist (mechanical; walk every item)

For each of AC1 through AC7 (with AC7 = amended form), verify the mechanical oracle named in γ-scaffold §4 holds. Record the oracle command run + the observed result for each AC. Mark green/red/ambiguous.

## Variable consistency table walk (T-486-1 + T-487-1 expanded scope — LOAD-BEARING)

Walk α's variable consistency table (in self-coherence.md). For every variable, verify presence/consistency in:
1. Manifest (wake-provider.json).
2. Rendered workflow YAML.
3. CI guard fixtures.
4. **Prompt prose** (T-487-1: the empirical motivator from cycle/487 R0 miss; do NOT skip).
5. **Smoke procedure** (N/A this cycle since no live smoke; record N/A but note for future cycles).
6. **PR body** (β reads the PR body draft α leaves OR predicts what α's PR body will say).

For each variable, flag drift. Cycle/487's R0 missed prompt-prose drift (the "declaration-only" historical block); cycle/496's review MUST cover prompt prose explicitly.

## Per-CI-step bash-e audit (cnos#478 mechanical-injection)

For every new `run:` step α added to `.github/workflows/install-wake-golden.yml`:
- Check `set -e` is in effect (either explicit `set -e` at the top, or implicit via bash-strict header).
- Check no `grep -c ... | true` or pipefail-equivalent class-traps.
- Check `if:` gating semantics match the test's intent (positive test ≠ negative test in `if:`).
- Audit every fixture branch invocation for trap-cleanup discipline (no orphaned fixture branches; no leaked files outside the runner's working tree).

## Local-write-fence false-positive check (LOAD-BEARING; new for cycle/496)

Inspect the rendered workflow's fence step:
1. Verify the fence checks `git status --porcelain -- '.cn-*/logs/'` (or equivalent local-scope check) NOT `git diff origin/main@before...origin/main@after -- .cn-*/logs/` (remote-state delta).
2. Verify the fence's check executes on the workflow's own working tree, NOT against a fetched remote ref.
3. Verify the CI fixture's false-positive resistance test exists AND passes.

If the fence is scoped to remote-state delta: that's a P0 blocker. Iterate immediately. The whole point of the AC7 amendment is to prevent this exact false-positive class.

## Output structure for `.cdd/unreleased/496/beta-review.md`

- **§R0 frontmatter.** YAML with cycle, base SHA, α R0 signal commit SHA, β review timestamp.
- **§R0 verdict.** `converge` OR `iterate` at the top, before any narrative.
- **§R0 AC table.** AC1–AC7 with per-AC green/red/ambiguous + the oracle command + observation.
- **§R0 variable consistency walk.** Every variable from α's table × every column from γ-scaffold §6's seed.
- **§R0 per-CI-step bash-e audit.** Every modified/added `run:` step + check result.
- **§R0 local-write-fence audit.** The 3 false-positive resistance checks above.
- **§R0 friction notes (if any).** FN-β-N entries naming specific drift β surfaced.
- **§R0 recommendation to γ-closeout (regardless of converge/iterate).** Single-line recommendations β predicts will land in cycle/496 closeouts.

## When to iterate vs converge

- **Iterate** if: any AC's oracle returns red; OR fence is remote-scoped (P0); OR prompt prose drifts from manifest field; OR the renderer's mis-declaration test does NOT exit code 4; OR CI fixtures are absent/broken.
- **Converge** if: every AC's oracle returns green; fence is local-scoped + false-positive-resistant; variable consistency holds across all columns including prompt prose; CI fixtures all pass.

## Hard constraints

- **Do NOT iterate for taste** — β iterates only when an oracle returns red OR a load-bearing invariant is violated.
- **Do NOT pre-empt cnos#497** — α's `.cdd/unreleased/496/` paths are correct for this cycle; do NOT flag them as "should be `.cds/`."
- **Do NOT modify any α-written file** — β reviews; does not implement.
- **Honor T-486-12 (operator-final-read P1)** — your converge is not the cycle's final gate; the operator reads after β.
- **Commit message pattern** — `β-496 R0 review: <converge|iterate>`. Co-Authored-By trailer pattern matches prior cycles.
```

## §8. Friction notes (anticipated cluster)

Cycle/487's scaffold opened with anticipated FNs; cycle/496 follows the pattern. These are the FN classes α + β are most likely to surface:

- **FN-1: write-fence false-positive risk (cited 4 times in this scaffold for emphasis).** The operator's amended AC7 is load-bearing precisely because the naive implementation false-positives. α MUST check this in α R0; β MUST verify in β R0. The CI fixture's parallel-admin-wake scenario is the empirical test.
- **FN-2: agent-admin manifest update scope.** AC2 requires agent-admin to declare `activation_log_writer: true` explicitly. The agent-admin's manifest is in `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json`. The rendered `.github/workflows/cnos-agent-admin.yml` does NOT need changes (admin-wake already loads `attach` legitimately). α SHOULD re-render agent-admin's golden to verify byte-equivalence (defensive check that adding the field doesn't perturb the renderer output for admin wakes); but no behavior change.
- **FN-3: dispatch-protocol skill failure-class taxonomy expansion.** The skill currently lists `dispatch_protocol_missing`, `dispatch_protocol_mismatch`, `dispatch_label_drift` (per `responsibilities[]` in cds-dispatch manifest). Adding `dispatch_activation_log_write_violation` is a doctrine micro-amendment. α should make the addition consistent in shape (1-paragraph description; failure trigger; runtime behavior). β should verify the new class doesn't conflict with existing classes' semantics.
- **FN-4: golden regeneration sha256.** Re-rendering `cnos-cds-dispatch.yml` produces a new sha256. The existing CI guard step in `install-wake-golden.yml` that checks "rendered = golden byte-identical" needs the new sha256. Sequence: α (1) updates the golden in the SAME commit that updates the manifest + prompt, (2) re-renders production substrate `.github/workflows/cnos-cds-dispatch.yml` in the SAME commit, (3) updates the CI guard's expected sha256 in the SAME commit. This avoids a transient "rendered ≠ golden" intermediate state.
- **FN-5: prompt prose drift risk.** T-487-1 lesson: cycle/487's R0 missed prompt-prose drift. β R0's variable consistency walk MUST include prompt prose explicitly. The cds-dispatch prompt's `responsibilities[]` enumeration + disallowed-surfaces section are the load-bearing prompt-prose surfaces here.
- **FN-6: D1 (c) edge case — does `cn install-wake` know to skip `attach`?** The renderer reads the manifest's `loaded_skills` field (or equivalent). The cleanest mechanism: when `activation_log_writer: false`, the renderer mechanically omits `attach` from the rendered skills section REGARDLESS of what the manifest declares. Defense in depth: even if a future provider author forgets and declares `loaded_skills: [activate, attach]` alongside `activation_log_writer: false`, the renderer enforces the partition. AC4's mis-declaration check is the primary enforcement; this is the secondary defense.
- **FN-7: bootstrap-δ self-application paradox surface — cnos#496 was claimed by bootstrap-δ; the live `cds-dispatch` wake's selector requires `status:todo`; cnos#496 is at `status:in-progress`. The wake CAN still fire on the issue label-event (issues:labeled trigger) but the SELECTOR check filters it out. Verify post-merge that cds-dispatch did NOT attempt to claim cnos#496 during the cycle. Operator-final-read defense-in-depth.

## §9. Cross-references + carryforwards

| Carryforward | Source | How applied this cycle |
|---|---|---|
| T-486-1 (variable consistency table) | cycle/486 closeouts | α writes seed table; β walks every variable × every surface |
| T-486-7 (R[N] off-by-one prevention; δ-direct R1 pattern) | cycle/486 closeouts | α's iteration discipline cited in §6 |
| T-486-12 (operator-final-read defense-in-depth, P1) | cycle/486 + cycle/487 promotion | β explicitly defers "final gate" to operator-final-read |
| T-486-15 (predecessor-closeout-reading) | cycle/486 | γ-scaffold reads cycle/487 + cycle/491 closeouts; α + β read this scaffold |
| T-487-1 (variable consistency table extended to ALL surfaces — prompt prose + runtime label set + smoke procedure + PR body) | cycle/487 closeouts | β prompt's checklist enumerates 6 columns; prompt prose is column 4 with explicit drift check |
| T-487-2 (dispatch_label_missing doctrine; sibling failure-class concern) | cycle/487 closeouts | FN-3 — `dispatch_activation_log_write_violation` mirrors the same naming pattern + extends the taxonomy |
| Long-arc partition (mechanize the protocol; model the judgment) | cnos#495 umbrella [#4792969173](https://github.com/usurobor/cnos/issues/495#issuecomment-4792969173) | The CYCLE's framing — first concrete enforcement of mechanical orchestration partition |
| KISS preference (no skill split unless forced) | operator directive on cnos#495 | D1 (c) chosen |
| Narrow cnos#496 (don't move orchestration stack into Go) | operator guardrail on cnos#495 | D2 (a) chosen over (c); cycle does not introduce new `cn` subcommands |
| Self-application paradox (bootstrap-δ-claim) | cycle/487 Sub 5C pattern | cnos#496 already at `status:in-progress`; live wake won't claim |

## §10. Non-goals + out-of-scope

- **No admin dispatch-summary implementation.** Deferred Sub 2; admin wake's `class: dispatch-summary` entry waits for both Sub 1 (this cycle) AND Sub 3 (cnos#497 decision) to land.
- **No `.cdd/` → `.cds/` rename.** cnos#497 is design-held (Model A vs Model B undecided); this cycle uses `.cdd/unreleased/496/` per current Model B substrate behavior. The renderer's substitution work in cnos#497B (if Model A wins) is OUT OF SCOPE here.
- **No CDR/CDW dispatch provider work.** Future packages; this cycle is cds-dispatch + agent-admin only.
- **No broader `attach` refactor.** D1 (c) means dispatch doesn't load attach; the skill itself stays unitary. Do NOT split `attach` into `attach-identity` + `attach-channel` (D1 (a)) unless a future cycle empirically forces it.
- **No retroactive rewrite of 2026-06-24 mixed log.** AC6 preserves the historical evidence.
- **No `cn` subcommand additions.** D2 (a) workflow-step v0; operator guardrail "keep cnos#496 narrow."
- **No Go orchestration migration.** The long-arc partition is the FRAME; this cycle is the FIRST STEP. Subsequent cycles (claim sequence to Go, label transitions to Go, δ routing to Go, return-token parsing to Go) are out of scope here.
- **No live smoke required.** Unlike cycle/487 (wave-goal-achievement cell requiring post-merge smoke), cycle/496's implementation is verifiable in PR review: CI fixtures provide the mechanical oracle for the write fence + renderer refusal. Post-merge observation (no .cn-sigma/logs/ commits from cds-dispatch) is operator-final-read confirmation, not a second-stage cycle.

## §11. Self-application paradox handling

cnos#496 was claimed by bootstrap-δ at [comment #4792872514](https://github.com/usurobor/cnos/issues/496#issuecomment-4792872514) (2026-06-24T19:30Z; head SHA `3f57210d`). Labels transitioned: `+ dispatch:cell + protocol:cds + status:in-progress` (from `status:todo` if it had been there; the issue was filed without an initial lifecycle label, so the transition is "no-lifecycle-label → status:in-progress").

The live `cds-dispatch` wake's selector requires `status:todo`; cnos#496 is at `status:in-progress`; the issue is invisible to the live wake's claim mechanism. Live wake will continue normal scheduled sweeps + label-event firings on other candidates without interfering.

**Post-merge confirmation (operator-final-read):** verify that the live wake did NOT attempt to claim cnos#496 during the cycle. `git log -- .cn-sigma/logs/ | grep cycle.*496` should return zero matches authored by cds-dispatch firings. Mirror cycle/487 Sub 5C pattern.

## §12. Cycle-shape forecast

- **R0 converge target.** Aim for R0-converge (cycle/485 + cycle/486 precedent; cycle/487 was R0 + 2 operator-iterates due to Stage 1 + Stage 2 cross-surface state synchronization complexity). Cycle/496's surface count is comparable to cycle/485 (renderer extension): convention amendment + manifest field + prompt scrub + renderer guard + CI fixtures + dispatch-protocol skill micro-amendment. T-487-1 expanded-scope variable consistency walk should catch prompt-prose drift in β R0 (avoiding cycle/487 R0's miss class).
- **Operator-final-read expectation (T-486-12 P1).** The operator is likely to flag at least one of: (a) fence false-positive resistance — verify the local-vs-remote distinction; (b) dispatch-protocol skill failure-class shape — verify `dispatch_activation_log_write_violation` aligns with existing taxonomy; (c) prompt prose explicitness — verify the disallowed-surfaces section names "mechanically enforced" not just "disallowed." γ-interface should run a final variable consistency walk before declaring β-converge sufficient.
- **Two-stage cycle NOT used.** No post-merge smoke required. CI fixtures provide the mechanical oracle for AC7; post-merge observation (no cds-dispatch writes to `.cn-sigma/logs/`) is operator confirmation, not a Stage-2 cycle.
- **Rollback plan.** If any AC's CI fixture fires red post-merge (unlikely if β R0 + operator-final-read both green), the rollback is: revert PR #(TBD when filed); restore prior cds-dispatch behavior (which is the empirical-drift baseline — acceptable as fallback since the production substrate has been operating with this drift for the entire `cds-dispatch` live deployment window without operational breakage). The rollback restores acceptable-baseline state, not a broken state.
- **PR title forecast.** `cycle/496 Sub 1: prevent package dispatch wakes from writing activation logs (mechanical fence)`.

## §13. The doctrine quote (preserved from cnos#495 umbrella)

> Use intelligence where meaning is unresolved.
> Use mechanics where the rule is known.

> Identity can be Sigma.
> Ownership is package-local.
> Memory is summarized upward, not dumped sideways.

This cycle's deliverable: the first mechanical-enforcement primitive crossing the partition. The activation log gets clean; the long-arc direction gets its first witness.

---

Filed by γ@cdd.cnos (bootstrap-δ via δ-interface session), 2026-06-24 (UTC).
