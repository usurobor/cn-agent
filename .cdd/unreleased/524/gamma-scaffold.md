---
cycle: 524
parent_issue: cnos#524
protocol: cds
cycle_branch: cycle/524
base_main_sha: db547ebe5b408e4c74092ad3ed56509e605894ef
head_sha_at_scaffold: db547ebe5b408e4c74092ad3ed56509e605894ef
mode: MCA
role: γ
authored_by: γ@cdd.cnos (wake-invoked δ dispatch, W4 scope — clean re-dispatch)
date: 2026-06-30 (UTC)
output_contract: γ-scaffold + α prompt + β prompt
dispatch_scope: W4-implementation
ac_set: AC4 (revised — header-only diff), AC5, AC6, AC7
w0_locked: true
w0_design_ref: .cdd/unreleased/524/w0-design.md
w1_ref: PR #526 (merged; SKILL.md files authored, #Wake CUE schema added)
w2_ref: PR #527 (merged; dual-source parity gate, --parity-check mode, CI step)
w3_ref: PR #529 (merged at 90b9e812; renderer default flipped to SKILL.md)
prior_w4_attempt: INVALIDATED (run 28464981342 — empty run, no PR/commits/closeout; RCA accepted;
  see issue comment "⛔ W4 empty run INVALIDATED")
rca_repair_ref: PR #531 (merged at 7950ab3d; closeout-integrity guard — scripts/ci/check-dispatch-closeout-integrity.sh)
operator_directive_ref: cnos#524 comment id 4847294646 (2026-06-30T19:38:24Z, "W4 final phase, clean re-dispatch")
run_class: repair_pass (issue carries prior status:review → status:changes → status:todo
  history from the invalidated W4 run; see REPAIR-PLAN.md for the finding-to-action map)
---

# γ-scaffold — cnos#524: wake-as-skill W4 implementation (final phase)

**Scope:** W4 — delete `wake-provider.json` + `prompt.md` for both wakes; make `cn-install-wake`
read SKILL.md unconditionally (no JSON fallback, no `--source`, no `--parity-check`); update the
generated golden/workflow header to attribute to `SKILL.md` instead of the now-deleted JSON+prompt
pair; retire or convert JSON-only tests; keep the live workflow behavior, selector, permissions,
concurrency, and prompt semantics byte-identical except for that header attribution change.

**Preconditions verified before this scaffold was authored (all hold against current `main`):**
- `cycle/524` absent from the remote (A1 — confirmed `git ls-remote --heads`).
- Closeout-integrity guard live (A2 — PR #531 merged at `7950ab3d`, confirmed ancestor of `main`
  via `gh pr view 531`; `scripts/ci/check-dispatch-closeout-integrity.sh` present on `main`).
- #524 open; `status:todo` at claim time (re-read confirmed by the dispatch wake before claim).
- No Demo 0 in scope.

**W0 reference:** `.cdd/unreleased/524/w0-design.md` §F (W4) — locked; governs the design.
**W1 state:** Both wake `SKILL.md` files exist; `#Wake` CUE schema in `schemas/skill.cue` (PR #526).
**W2 state:** `cn-install-wake` can parse SKILL.md via `--source skill`; parity proved byte-for-byte
including headers (PR #527).
**W3 state:** `cn-install-wake` default flipped to `source_type="skill"`; `--parity-check` now
proves render(JSON+prompt) == render(SKILL.md) (PR #529, merged `90b9e812`).

---

## §0. Repair re-entry note (cnos#516 preflight — see REPAIR-PLAN.md)

This claim is a `repair_pass` per the dispatch-protocol §Repair re-entry preflight: the issue's
comment history shows `status:review → status:changes` (the W4 RCA invalidation) followed by
`status:changes → status:todo` (this re-dispatch). **What was rejected was the *process*, not a
line of code** — the invalidated W4 run produced zero commits, zero PR, zero closeout; there is no
prior W4 branch state to repair against (`cycle/524` was deleted per A1). The "rejected finding" is
therefore process-shaped: *"a dispatch run can reach `status:review` with no deliverable."* That
finding was already fully remediated by a separate, already-merged cell (PR #531, the
closeout-integrity guard) — not by this cycle. This cycle's job is to deliver the actual W4 work
under the now-guarded process. `.cdd/unreleased/524/REPAIR-PLAN.md` records this mapping; α/β/δ
honor it; the converge closeout's `repair_evidence` block points back to it.

---

## §1. W4 scope

**What W4 delivers (the only writable changes):**

1. **Delete JSON+prompt pairs:**
   - `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` (delete)
   - `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` (delete)
   - `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` (delete)
   - `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` (delete)
   - `src/packages/cnos.core/commands/install-wake/test-fixtures/declaration-only/wake-provider.json` (delete)
   - `src/packages/cnos.core/commands/install-wake/test-fixtures/declaration-only/prompt.md` (delete)
   - `src/packages/cnos.core/commands/install-wake/test-fixtures/log-writer-misdeclaration/wake-provider.json` (delete)
   - `src/packages/cnos.core/commands/install-wake/test-fixtures/log-writer-misdeclaration/prompt.md` (delete)
   - (Each of these four directories already has a `SKILL.md` twin — confirmed present. The
     `SKILL.md` is what survives.)

2. **`src/packages/cnos.core/commands/install-wake/cn-install-wake` — remove the JSON source path
   entirely:**
   - Remove the `--source json|skill` flag (no flag needed once `skill` is the only source —
     OR keep `--source` accepting only `skill` for forward-compat and `die` on `json`; α's call,
     name the choice in `self-coherence.md`. Either is acceptable; the renderer MUST NOT silently
     accept `--source json` and try to read a now-absent file).
   - Remove the `--parity-check` flag and its execution block (§ ~1109-1135) — there is nothing
     left to prove parity against once JSON is gone.
   - Remove the `source_type` json/skill branching: the `skill_to_json_manifest` + `skill_body`
     extraction (currently gated by `if [ "$source_type" = "skill" ]`) becomes unconditional.
   - Remove the JSON-branch `required_fields` set (the `else` branch with `responsibilities`,
     `prompt_template`, `cross_references` as top-level required fields) — only the reduced
     SKILL.md-shape required-field set remains.
   - Remove the JSON-branch prompt resolution (`prompt_template_relpath` / `jq -r '.prompt_template'
     "$manifest_path"` on the raw JSON) — only the SKILL.md body-as-prompt path remains.
   - **Manifest resolution must now look for `SKILL.md`, not `wake-provider.json`:** the default
     lookup (`${manifest_dir}/wake-provider.json`) and the cross-package sibling-search fallback
     (`candidate="${sibling}/orchestrators/${wake_name}/wake-provider.json"`) both need to resolve
     `SKILL.md` instead. The `--manifest <path>` override (used by the CI fixture smokes) needs to
     accept a `SKILL.md` path directly (or a directory containing one — α's call; document it in
     `--help`).
   - **Header attribution must change from JSON+prompt to SKILL.md:** the header currently emits
     two `#` lines (`manifest: .../wake-provider.json` and `prompt: .../prompt.md`), derived from
     `json_manifest_path` / `json_prompt_path` (which exist purely to keep the header stable across
     `--source` values — that reason is gone). Replace with a single source line naming the
     `SKILL.md` path, e.g.:
     ```
     # DO NOT EDIT. Rendered by `cn install-wake <name>` from:
     #   source: orchestrators/<name>/SKILL.md
     ```
     This is the **only** permitted non-comment-prose change to the rendered golden/workflow
     bytes — see §1.1 below. Exact wording is α's call; keep it accurate and parseable by a human
     diffing the golden.
   - Update the file's top-of-script doc comment (Usage/Arguments/Exit-codes block) to drop
     `--source`/`--parity-check` documentation and the now-dead exit code 5.
   - Renumber or retire exit code 5 (parity-check failure) — it no longer has a caller. Either
     remove it from the exit-code table or repurpose it (α's call; if repurposed, document why;
     if removed, leave a gap-note rather than silently reusing 5 for something unrelated).

3. **`.github/workflows/install-wake-golden.yml`:**
   - Remove the **"W3 parity check — render(JSON+prompt) == render(SKILL.md)"** step entirely
     (nothing left to compare).
   - Retire or convert the **"AC2 negative-case smoke (malformed manifest is rejected)"** step
     (its own comment already says *"Retire/replace in W4 when JSON is removed"*). Convert to a
     SKILL-only negative case: a temp `SKILL.md` fixture missing a required `wake:` field (e.g. no
     `role:`) should still produce exit 2 + a precise `required field "..." missing` message via
     the renderer's default (skill) path. This preserves AC2's oracle (malformed declarations are
     rejected, not silently defaulted) without depending on JSON.
   - Update the **"AC5 — declaration-only refusal"** step's `--manifest` argument: it currently
     points at `test-fixtures/declaration-only/wake-provider.json`; repoint it at
     `test-fixtures/declaration-only/SKILL.md` (or the directory, per however α resolved
     `--manifest` above). The fixture's `SKILL.md` already declares `activation_state:
     declaration-only` (confirmed) — same refusal-gate assertions (exit 3, stderr names
     `declaration-only`, stderr mentions cnos#454/cnos#467/preconditions) still apply unchanged.
   - Update the **"AC4 (cycle/496) — renderer refusal on mis-declaration"** step's `--manifest`
     argument: repoint at `test-fixtures/log-writer-misdeclaration/SKILL.md`. Same assertions
     (exit 4, stderr contains `activation_log_writer mis-declaration`) apply unchanged.
   - Leave the rest of the workflow (re-render steps, golden-diff step, sha256 live-vs-golden
     check, idempotence checks, YAML-parses check, substrate-shape checks, AC8/AC7 renderer-
     authority-leak audit) structurally unchanged — only the header-diff expectation shifts (see
     §1.1).

4. **Two golden files + two live workflow files — re-render, header-only diff:**
   - `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`
   - `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`
   - `.github/workflows/cnos-agent-admin.yml`
   - `.github/workflows/cnos-cds-dispatch.yml`
   Each currently carries the two-line `# manifest: .../wake-provider.json` / `# prompt:
   .../prompt.md` header. After the renderer change, re-running `cn install-wake <name>` (and
   `cn install-wake cds-dispatch --out .github/workflows/cnos-cds-dispatch.yml` /
   `cn install-wake agent-admin --out .github/workflows/cnos-agent-admin.yml`) must change **only**
   those header lines. Every other byte (triggers, permissions, concurrency, job body, the embedded
   prompt) MUST be unchanged. Commit the re-rendered files.

5. **`.cdd/unreleased/524/` records** — this scaffold (already written), `REPAIR-PLAN.md` (already
   written), `self-coherence.md §R3`, `beta-review.md §R3`, and on converge: `alpha-closeout.md`,
   `beta-closeout.md`, `gamma-closeout.md` amendments (append, do not overwrite prior W0/W1/W2/W3
   history) carrying the `repair_evidence` and `deliverable_evidence` blocks.

### §1.1 The header-only-diff invariant (read before touching anything else)

The operator's directive is explicit and absolute on this point — quoting verbatim:

> **Stop if:** any non-header workflow bytes change · wake behavior changes · deleting JSON/prompt
> changes semantics · renderer needs old JSON compatibility to pass · #524 would be closed · Demo 0
> gets invoked.

**This is a hard scope boundary, not a target to optimize toward.** Two `SKILL.md` files — Surf
`src/packages/cnos.core/orchestrators/agent-admin/SKILL.md` and
`src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` — contain **body prose** that mentions
`wake-provider.json` (e.g. `cds-dispatch/SKILL.md` line ~94: a markdown link
`[wake-provider.json](wake-provider.json)` describing it as "the machine-readable source of
truth"; both files' `## Responsibilities (from wake-provider.json; body reference)` and `##
Cross-references (from wake-provider.json)` section headers). **The SKILL.md body IS the rendered
prompt** (per the W0 design §E body-as-prompt rule) — editing this prose changes non-header bytes
of the rendered golden/workflow, which directly violates the stop condition above and the
"no prompt semantic change" scope constraint.

**α MUST NOT edit `agent-admin/SKILL.md` or `cds-dispatch/SKILL.md` body prose in this cycle**, even
though doing so would make the "no active references to wake-provider.json/prompt.md remain"
required-proof bullet literally true of the *entire* repository. Resolve the apparent tension this
way: **"no active references remain" is satisfied by (a) the files being deleted (AC5) and (b) the
renderer's resolution/parsing logic no longer reading or referencing them (§1 item 2) — not by
scrubbing every prose mention inside frozen prompt bodies.** The stale `[wake-provider.json
](wake-provider.json)` link and the two stale section-header parentheticals become a dangling
reference once the files are deleted (the link target 404s). **Name this explicitly as a known gap**
in `gamma-closeout.md`'s process-gap audit and file a tracked follow-up issue for a future cell to
correct the SKILL.md body prose (that correction is itself a legitimate, separately-scoped prompt
content change, not a renderer/substrate change, and belongs in its own cell). Do not fix it here.

**Also out of scope (do not touch, even though they textually match `wake-provider.json`/
`prompt.md`):** `schemas/skill.cue`, `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md`,
`src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md`,
`src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`, `docs/**`, `src/go/internal/dispatch/*_test.go`,
`src/go/internal/cli/cmd_activate.go`, any `.cdd/releases/**` or `.cdd/unreleased/{N}/` for N≠524,
any kata file. These are historical/design-doc surfaces describing the migration narrative or
unrelated subsystems; editing them is scope creep per the non-goals list ("no `.cdd/releases/`
work", "no Demo 0") and risks the "no non-goal noun appears in the diff" closure condition.

---

## §2. AC oracle list (W4-relevant ACs)

### AC4 (revised) — Byte-identical goldens, header-only diff

**Invariant:** Re-rendering both wakes after the renderer change produces a diff confined to the
two header lines (or whatever single replacement line α chooses) on both goldens and both live
workflows. Zero bytes differ anywhere else.

**Oracle (pass conditions):**
- `git diff -- <golden|workflow files>` shows changes ONLY on the header `#` lines.
- `install-wake-golden` CI's "Re-render" + "Verify goldens unchanged" steps are necessarily
  **expected to show a diff this one time** (the header changed) — this CI check as currently
  written treats ANY diff as a failure. α/β must confirm the **committed** goldens already reflect
  the new header (i.e. α commits the re-rendered output, so CI's re-render-and-diff is comparing
  the new renderer against the new committed golden, which match — clean). Do not rely on CI
  tolerating a diff; commit the new render as the new golden.
- sha256 live-vs-golden check (cds-dispatch) still passes (both are the freshly re-rendered byte
  stream).
- Idempotence checks still pass (second render of the new output is a no-op).
- YAML-parses + substrate-structural-shape checks still pass unchanged (header is comments only,
  doesn't affect YAML structure).

**Failure scenario:** Any job body, trigger, permission, concurrency, or prompt byte changes.
STOP — do not commit; surface a `status:blocked` per §9.6 of the δ skill instead of patching
around it.

---

### AC5 — JSON + prompt deleted

**Invariant:** `wake-provider.json` and `prompt.md` are absent from both orchestrator dirs and
both test-fixture dirs; the renderer reads only `SKILL.md`.

**Oracle:** `find src/packages/cnos.core/orchestrators src/packages/cnos.cds/orchestrators
src/packages/cnos.core/commands/install-wake/test-fixtures -name 'wake-provider.json' -o -name
'prompt.md'` returns empty. `cn install-wake agent-admin` and `cn install-wake cds-dispatch` still
render byte-identical (header-only-diff, per AC4) goldens with no `wake-provider.json`/`prompt.md`
present anywhere on the filesystem.

**Failure scenario:** The renderer's manifest-resolution fallback (cross-package sibling search)
still globs for `wake-provider.json` and silently produces "manifest not found" instead of finding
the `SKILL.md` — a regression that breaks `cn install-wake <name>` for any new package without an
explicit `--manifest`. Verify the fallback path explicitly (e.g. by removing `CN_PACKAGE_ROOT`
from the environment and re-running from repo root, the same way "direct invocation" is described
in the script's own header comment).

---

### AC6 — Renderer refusals preserved (defense-in-depth)

**Invariant:** `activation_state` declaration-only → exit 3; `activation_log_writer`
mis-declaration → exit 4. Both fire from the SKILL.md-only path (there is no other path now).

**Oracle:** CI's "AC5 — declaration-only refusal" step (repointed at
`test-fixtures/declaration-only/SKILL.md`) and "AC4 (cycle/496) — renderer refusal on
mis-declaration" step (repointed at `test-fixtures/log-writer-misdeclaration/SKILL.md`) both still
assert the same exit codes + stderr substrings as before W4.

**Failure scenario:** Removing the JSON branch accidentally removes a shared validation helper
(`require_field`, `manifest_error`, `log_writer_mis_declaration`) that the SKILL.md path also
depends on. Run both refusal smokes locally before declaring AC6 satisfied.

---

### AC7 — All gates green

**Invariant:** No inherited-cap, no new red.

**Oracle:** I1 / I2 / I4 / I5 / I6 / Go / Package / Binary / `install-wake golden` /
`dispatch-repair-preflight` (cnos#516) / `dispatch-closeout-integrity` (cnos#524 W4 RCA, #531) all
green on this cell's PR. The AC8/AC7 renderer-authority-leak audit (grep for admin-shape and
dispatch-shape role-decision strings in the renderer source) must remain green — removing the JSON
branch should, if anything, shrink the surface that audit scans, not add new leaks.

---

## §3. Implementation contract (for α)

### §3.1 Required-proof checklist (operator's verbatim list — α self-verifies all of these before
requesting β review)

- [ ] agent-admin renders from `SKILL.md`
- [ ] cds-dispatch renders from `SKILL.md`
- [ ] `wake-provider.json` and `prompt.md` are gone (both orchestrator dirs + both test-fixture dirs)
- [ ] no active references to `wake-provider.json` / `prompt.md` remain **in the renderer's
      resolution/parsing logic or in any file path the renderer reads at runtime** (per §1.1 scope
      boundary — body prose in the two wake SKILL.md files is explicitly NOT in scope)
- [ ] live workflows match regenerated goldens (sha256 check, cds-dispatch)
- [ ] golden diff is header-only / attribution-only (confirm by diffing old → new committed goldens)
- [ ] `install-wake-golden` green
- [ ] `dispatch-repair-preflight` green
- [ ] `dispatch-closeout-integrity` green
- [ ] I1/I2/I4/I5/I6/Go/Binary/Package green

### §3.2 Scope guardrails for α

**MUST change:**
- `src/packages/cnos.core/commands/install-wake/cn-install-wake`
- `.github/workflows/install-wake-golden.yml`
- `src/packages/cnos.core/orchestrators/agent-admin/{wake-provider.json,prompt.md}` (delete)
- `src/packages/cnos.cds/orchestrators/cds-dispatch/{wake-provider.json,prompt.md}` (delete)
- `src/packages/cnos.core/commands/install-wake/test-fixtures/declaration-only/{wake-provider.json,prompt.md}` (delete)
- `src/packages/cnos.core/commands/install-wake/test-fixtures/log-writer-misdeclaration/{wake-provider.json,prompt.md}` (delete)
- `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` (header-only diff)
- `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` (header-only diff)
- `.github/workflows/cnos-agent-admin.yml` (header-only diff)
- `.github/workflows/cnos-cds-dispatch.yml` (header-only diff)
- `.cdd/unreleased/524/` records

**MUST NOT change:**
- `src/packages/cnos.core/orchestrators/agent-admin/SKILL.md` (body prose — see §1.1)
- `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md` (body prose — see §1.1)
- `schemas/skill.cue`
- `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md`
- `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md`
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md`
- `src/go/**`
- Any `docs/**`, any `.cdd/releases/**`, any other `.cdd/unreleased/{N}/`
- Any non-header byte of either golden or either live workflow

**PR linkage invariant:** Every commit message and the PR body MUST use `Refs #524` or
`Part of #524` — NEVER `Closes/Fixes/Resolves/close/fix/resolve #524` (including negated forms).
Issue #524 must remain open after W4 merges (operator does the final close/cycle-complete read).

### §3.3 Friction notes (FN) for α

**FN-W4-1: order of operations matters for a clean diff.** Implement the renderer change FIRST,
verify it renders byte-identical-except-header against the *current* (W3) committed goldens by
temporarily keeping the JSON files in place, THEN delete the JSON/prompt files, THEN re-render and
commit the final goldens. This gives a clean two-step verification: "renderer change alone is
header-only" then "JSON deletion doesn't break anything because the renderer no longer needs it."
Squash into clean commits as you see fit, but verify in that order.

**FN-W4-2: `--manifest` override contract change.** The CI fixture smokes (AC2/AC5/AC4-496) call
`--manifest <path-to-wake-provider.json>` today. Decide a clear new contract for `--manifest`
(accepts a `SKILL.md` path directly is the natural choice, consistent with the default lookup
change) and update all three CI steps' invocations to match. Document the new contract in the
script's `--help` text.

**FN-W4-3: exit code 5 retirement.** `--parity-check` owned exit 5. Once removed, decide: drop exit
5 from the documented exit-code table (clean), or leave a one-line note that 5 is retired/reserved.
Either is fine; just don't leave stale prose claiming exit 5 still means "parity failure" when
nothing can trigger it anymore.

**FN-W4-4: the cross-package sibling-search fallback is easy to miss.** It's a second occurrence of
the `wake-provider.json` filename glob (the first is the default per-package lookup) — search the
whole file for the literal string `wake-provider.json` after your edits to confirm zero renderer-
logic occurrences remain (comments documenting *why* JSON was removed are fine to keep; resolution
logic must not still glob for it).

**FN-W4-5: AC2 conversion detail.** A SKILL.md missing `wake.role` will fail inside
`skill_to_json_manifest` (the python YAML→JSON bridge) or downstream in `require_field` — confirm
which, and make sure the resulting stderr message is precise enough for the CI step's `grep -q
'required field "..." missing'`-style assertion to still mean something. If the python bridge
errors before reaching `require_field`, the message format will differ from the old JSON-path
message — that's fine, just make the CI assertion match whatever the actual message is, and keep
the assertion narrative honest (don't just assert exit 2, assert exit 2 + an informative substring).

---

## §4. β review prompt

β reviews the W4 implementation against this scaffold and the operator's verbatim "Required proof"
+ "Stop if" lists (quoted in full in §1.1 and §3.1 above). β's primary checks, in order:

1. **Header-only diff, mechanically verified.** Diff the old committed goldens (pre-W4, in git
   history) against the new committed goldens — confirm the change set is confined to the header
   lines. Do the same for both live workflow files. This is the single highest-stakes check in this
   cycle; do not take α's word for it, run the diff yourself.
2. **Files actually deleted.** `find ... -name wake-provider.json -o -name prompt.md` under the two
   orchestrator dirs and two test-fixture dirs returns empty.
3. **Renderer has no JSON path left.** Read the diff on `cn-install-wake`; confirm no
   `source_type = "json"` branch, no `--source`/`--parity-check` flag handling that still works,
   no remaining `wake-provider.json` glob in resolution logic (comments OK).
4. **SKILL.md body prose untouched.** Confirm `agent-admin/SKILL.md` and `cds-dispatch/SKILL.md`
   have ZERO diff (per §1.1 — this is a scope violation if touched, not a nice-to-have).
5. **CI workflow updated correctly.** "W3 parity check" step is gone. AC2 negative case converts to
   a SKILL.md fixture (not silently deleted — the oracle must still exist in some form). AC5 and
   AC4(cycle/496) steps' `--manifest` args point at the SKILL.md fixtures and still pass.
6. **All required-proof checklist items (§3.1) are independently confirmed**, not just asserted in
   `self-coherence.md`.
7. **`repair_evidence` present** (this is a `repair_pass` per §0 / `REPAIR-PLAN.md`) — confirm the
   closeout names the prior rejection, what was required, what's completed, and cites evidence the
   branch state now differs from the invalidated empty run (trivially true here: this branch has
   actual commits; the invalidated run had none — but the block must still be written out).
8. **`deliverable_evidence` present** (cnos#524 W4 RCA / #531 guard) — PR number, head/base SHA,
   commits-beyond-base count > 0, the six closeout artifact filenames, before any `status:review`
   transition is even considered.
9. **Scope compliance.** Diff confined to the §3.2 MUST-change list; nothing in the MUST-NOT list
   touched; no non-goal noun (Demo 0, releases work, unrelated skill docs) appears.
10. **PR linkage.** `Refs #524` / `Part of #524` only; #524 stays open.
11. **CI green.** I1/I2/I4/I5/I6/Go/Package/Binary/`install-wake golden`/`dispatch-repair-preflight`/
    `dispatch-closeout-integrity` all green on the PR.

β verdict MUST be one of: `converge` (all checks pass; proceed to closeouts + PR + `status:review`
return token) or `iterate` (specific findings enumerated; α repairs before converge — cell stays
`status:in-progress` during the iteration per dispatch-protocol §3.6, this is an INTERNAL loop, not
a new `status:changes`).

---

_Authored by γ@cdd.cnos (wake-invoked δ dispatch, W4 scope — clean re-dispatch), 2026-06-30 (UTC).
W4 operator directive: cnos#524 comment id 4847294646 (2026-06-30T19:38:24Z). W0 design:
`.cdd/unreleased/524/w0-design.md` §F. W3 baseline: PR #529 merged at `90b9e812`. RCA + guard:
PR #531 merged at `7950ab3d`. Repair-re-entry map: `.cdd/unreleased/524/REPAIR-PLAN.md`._
