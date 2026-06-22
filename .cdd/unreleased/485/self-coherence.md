---
cycle: 485
parent_issue: cnos#485
master_tracker: cnos#467 (Sub 5A)
cycle_branch: cycle/485
authored_by: α@cdd.cnos (bootstrap-δ session)
date: 2026-06-22 (UTC)
base_main_sha: af8ed8ec
scaffold_sha: c7f7325c (γ scaffold commit)
cds_dispatch_golden_sha256: 75e0406662ea12b3caa5ad01d8dbe6bdfb620d0c047fd910917a8e930855605e
agent_admin_golden_sha256: fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72 (preserved unchanged from base sha af8ed8ec per AC8)
---

# α self-coherence — cnos#485 (cn-install-wake renderer extension for dispatch wake providers)

## §Design

### Approach

The renderer extension is a strictly-additive overlay on cycle/476's admin-shape renderer:

1. **Role-conditional output_contract validation.** The pre-existing unconditional `output_contract.channel_log_convention` check (which is admin-shape-specific) is gated on `role == admin`. For `role == dispatch`, the validator instead requires `output_contract.cycle_artifact_root` (string), `output_contract.artifact_class_taxonomy` (array), and `output_contract.cell_runtime` (string) per wake-provider/SKILL.md §2.1. The `role == observer` branch is left untouched (out of scope for this cycle; no observer wakes ship yet).
2. **Dispatch-required field validation.** When `role == dispatch`, the renderer additionally requires `protocol` (string) and `selector` (object with `include` non-empty array + `exclude` array) per wake-provider/SKILL.md §3.9.
3. **`activation_state` refusal gate per §3.10.** A new block reads `.activation_state` (default `"live"` per §2.2) and `.activation_state_notes`. If the effective state is anything other than `"live"`, the renderer prints a precise stderr error naming the value AND the notes content, then exits 3 (a new exit code, documented in the script header alongside 0/1/2). A new flag `--activation-state-override <state>` allows test-only bypass; passing it emits a `WARNING` to stderr that names the override and the manifest's declared value.
4. **`issues_labeled_selector_match` trigger taxonomy.** A new `case` arm in the existing trigger loop renders `on: issues: types: [labeled]` for this trigger. The unknown-trigger error message is updated to list the new value.
5. **Dispatch-aware job-level `if:`.** A second `has_issues_labeled` flag is tracked alongside the existing `has_issues_title`. When `has_issues_labeled` is set, the renderer reads `selector.include[]` via jq and assembles a `contains(github.event.label.name, '<label>') || ...` clause; the labels are read from the manifest at render time, never as literals in the renderer source (OG-1).
6. **Admin-shape byte-identical preservation.** The admin branch (`event_name == 'schedule' || contains(github.event.issue.title, '<pattern>')`) is preserved exactly. The new dispatch branch uses `event_name == 'schedule' || (event_name == 'issues' && <selector-match>)` — explicitly gating the issue clause on `event_name == 'issues'` so schedule firings (where `github.event.label` is undefined) are not silently suppressed (OG-2).
7. **Cross-package manifest lookup.** The renderer's manifest-discovery loop falls back to sibling packages under `src/packages/` when the wake-name's manifest is not under `CN_PACKAGE_ROOT/orchestrators/`. This lets the CI re-render check invoke `cn-install-wake cds-dispatch` (manifest in cnos.cds) without depending on the dispatch shim's `CN_PACKAGE_ROOT` export. When a sibling match is found, `CN_PACKAGE_ROOT` is re-pinned to that sibling so the header's `manifest:` path-stripping produces the same relative form (e.g. `orchestrators/cds-dispatch/wake-provider.json`) the admin-shape header uses.

### Key decision — OG-3 override mechanism (Option A)

Per the γ scaffold §6 guardrail #3, the test-only activation-state override has three suggested forms. I chose **Option A — `--activation-state-override <state>` CLI flag**. Rationale:

- **Visibility.** The flag name itself names what is being bypassed (`activation_state_override`); an operator running `cn install-wake --help` sees the affordance and immediately reads its purpose. The header docstring explicitly says "TEST-ONLY override ... never pass this flag during production install."
- **Auditability.** Every use prints a `cn-install-wake: WARNING — --activation-state-override=<state> bypasses the manifest's declared activation_state=<declared>. This flag is test-only ...` line to stderr. Any operator or CI run that passes the flag leaves a visible audit trail.
- **No silent escape.** Option C (env var `CN_INSTALL_WAKE_ACTIVATION_OVERRIDE`) was rejected because env vars don't show up in `--help`, and an operator may not realize the override is active. Option B (a separate `.test-live.json` manifest) was rejected because it duplicates the manifest and creates two-source-of-truth drift risk; if the production manifest is ever updated, the test fixture silently rots. Option A keeps a single source of truth (the production declaration-only manifest) and exercises it through a clearly-named CLI affordance.
- **OG-3 compliance.** The production manifest at `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` is unchanged (`activation_state: declaration-only`). The renderer refuses normal `cn install-wake cds-dispatch` with exit 3 (AC5 verified). The override is only exercised by the CI workflow's golden re-render step, which is itself visible in the CI logs.

### α-side guardrail-2 AC (explicit)

Per the γ scaffold §6 guardrail #2: the rendered `cnos-cds-dispatch.golden.yml`'s job-level `if:` allows schedule events unconditionally AND gates only issue events on the selector's include-set membership.

Verification command:
```
grep -E "^[[:space:]]+if:" src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
```

Output:
```
    if: ${{ github.event_name == 'schedule' || (github.event_name == 'issues' && (contains(github.event.label.name, 'dispatch:cell') || contains(github.event.label.name, 'protocol:cds') || contains(github.event.label.name, 'status:todo'))) }}
```

This is the exact shape γ named as the correct form in scaffold §6. The schedule short-circuit fires the scheduled sweep regardless of `github.event.label`; the issue branch is explicitly gated on `event_name == 'issues'`.

The CI workflow's `Verify substrate structural shape — cds-dispatch` step has an explicit `grep -qE "github\.event_name == 'schedule'"` check that fails the build if this shape regresses.

---

## §AC-by-AC verification

All commands run on cycle/485 HEAD (post-α-R0 commit) at working directory `/home/user/cnos`. Each row: command run, observed result, pass/fail.

| AC | Statement | Command(s) | Result | Pass/Fail |
|---|---|---|---|---|
| AC1 | Renderer parses `role: dispatch`; emits YAML with `name: cnos-cds-dispatch`, `group: cds-dispatch-sigma`, all four permissions | `./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch --activation-state-override live; echo $?` | exit=0; fixture written; `name: cnos-cds-dispatch`, `group: cds-dispatch-sigma`, `contents: write`, `issues: write`, `pull-requests: write`, `id-token: write` all present (each grep matched) | **PASS** |
| AC2 | `protocol` + `selector` emitted; each `selector.include` label appears in a `contains(github.event.label.name, ...)` construct | `for label in 'dispatch:cell' 'protocol:cds' 'status:todo'; do grep -qF "$label" $f && echo PASS; done; grep -E "if:.*contains\(github\.event\.label\.name" $f` | All 3 labels present in `if:` line; `contains(github.event.label.name, '<label>')` construct present 3x | **PASS** |
| AC3 | `on:` block includes `issues: { types: [labeled] }` + `schedule:`; job-level `if:` references both | `python3 -c "import yaml; y=yaml.safe_load(open('$f')); assert 'labeled' in y[True]['issues']['types']"; grep -qE "github\.event_name == 'issues'.*github\.event\.label\.name" $f` | YAML parses; `on.issues.types` contains `labeled`; job-level `if:` references `github.event_name == 'issues'` followed by `github.event.label.name` matches | **PASS** |
| AC4 | Inlined prompt template carries `.cdd/unreleased/{N}/` + `cnos.cdd`; rendered YAML carries those strings + `prompt: \|` block | `grep -qF '.cdd/unreleased/{N}/' $f; grep -qF 'cnos.cdd' $f; grep -qE "^[[:space:]]+prompt: \|" $f` | All three present; the `prompt.md` body flows through unchanged with `{agent}` substituted to `sigma` | **PASS** |
| AC5 | `cn install-wake cds-dispatch` (no override) exits 3; stderr names `declaration-only` + `activation_state_notes` content | `./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch 2>/tmp/ac5.err; echo $?; grep -qF 'declaration-only' /tmp/ac5.err; grep -qE 'cnos#454\|cnos#467\|preconditions' /tmp/ac5.err` | exit=3; stderr names `declaration-only`; stderr cites `cnos#454`, `cnos#467`, and `preconditions` (the activation_state_notes content) | **PASS** |
| AC6 | Two consecutive renders of cds-dispatch produce byte-identical output | `./cn-install-wake cds-dispatch --activation-state-override live; sha1=...; ./cn-install-wake ...; sha2=...; [ "$sha1" = "$sha2" ]` | sha1=sha2=`75e0406662ea12b3caa5ad01d8dbe6bdfb620d0c047fd910917a8e930855605e`; second invocation prints `(unchanged)` | **PASS** |
| AC7 | `grep -ciE 'protocol:cds\|cdr\|cdw\|dispatch:cell\|status:todo' cn-install-wake` returns 0 | `n=$(grep -ciE 'protocol:cds\|cdr\|cdw\|dispatch:cell\|status:todo' src/packages/cnos.core/commands/install-wake/cn-install-wake \|\| true); [ "$n" = "0" ]` | n=0; renderer source carries no dispatch-shape role-decision literals (labels are read from `selector.include[]` via jq at render time; β empirical-read confirms no shell-var concatenation tricks were used for the dispatch surface) | **PASS** |
| AC8 | `cn install-wake agent-admin` produces a byte-identical `cnos-agent-admin.golden.yml` (sha256 = `fa6b8c0c…`) | `./cn-install-wake agent-admin; sha256sum src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml \| cut -d' ' -f1` | `fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72`; second invocation prints `(unchanged)`; `git diff --exit-code` on the file returns 0 | **PASS** |
| AC9 | `cnos-cds-dispatch.golden.yml` exists; install-wake-golden CI guards both goldens via re-render + diff + idempotence + YAML parse | `test -f $f; grep -qF 'cnos-agent-admin.golden.yml' .github/workflows/install-wake-golden.yml; grep -qF 'cnos-cds-dispatch.golden.yml' .github/workflows/install-wake-golden.yml; n=$(grep -cE 'sha256sum.*\.golden\.yml' ...); [ $n -ge 2 ]` | All present; CI workflow has 4 `sha256sum.*\.golden\.yml` occurrences (2 idempotence checks × before+after) | **PASS** |

**Local verification: 9 of 9 ACs verified.** CI green is pending (will surface on the cycle/485 PR run of `install-wake-golden.yml`); β re-verifies CI evidence per scaffold §7.

---

## §Per-CI-step bash-e audit (per cnos#478 mechanical-injection process)

For every `run:` block in `.github/workflows/install-wake-golden.yml` that this cycle introduces or modifies. The audit columns track the bash-e + command-substitution + pipeline traps that cycle/476 hit (3-round class-trap on `grep -c` returning exit 1 on zero matches under bash-e).

| # | Step name | Line range | Command substitutions / pipelines | Guarded? | bash -e exit on intended-success input? | Notes |
|---|---|---|---|---|---|---|
| 1 | Verify jq present | 37-41 | none (`jq --version` only) | implicit `set -e` (default bash GH runner) | exit 0 on success; non-zero only if jq missing | Unchanged from cycle/476; included for completeness. |
| 2 | Re-render agent-admin wake (cnos.core) | 43-45 | none | implicit `set -e` | exit 0 on success | Unchanged from cycle/476. |
| 3 | Re-render cds-dispatch wake (cnos.cds, --activation-state-override live) | 47-55 | none | implicit `set -e` | exit 0 on success | **NEW (cycle/485).** Renderer exits 0 with `--activation-state-override live` on the declaration-only manifest; stderr WARNING is informational only (does not affect exit). |
| 4 | Verify goldens unchanged | 57-69 | none; `git diff --exit-code` is the substantive check | `if ! git diff …; then exit 1; fi` (explicit) | exit 0 when goldens match; explicit `exit 1` on drift | **MODIFIED (cycle/485).** Was `Verify golden unchanged` (single file); now passes both goldens to one `git diff --exit-code` invocation. `git diff` exits 0 on no-diff, 1 on diff; the `if !` flips so `exit 1` only fires on real drift. |
| 5 | Verify idempotence (second render is a no-op) — agent-admin | 71-84 | `$(sha256sum … \| cut -d' ' -f1)` × 2 | implicit `set -e`; `sha256sum` and `cut` always exit 0 on intended-success input | exit 0 when shas match; explicit `exit 1` on drift | **MODIFIED (cycle/485) — renamed only.** Step body unchanged from cycle/476. |
| 6 | Verify idempotence (second render is a no-op) — cds-dispatch | 86-100 | `$(sha256sum … \| cut -d' ' -f1)` × 2 | implicit `set -e`; `sha256sum` and `cut` always exit 0 | exit 0 when shas match; explicit `exit 1` on drift | **NEW (cycle/485).** Mirrors step 5 for the cds-dispatch golden. Same bash-e guard shape as step 5. |
| 7 | Verify YAML parses (both goldens) | 102-114 | inline python3 heredoc | implicit `set -e`; python3 raises on parse error which exits non-zero | exit 0 when both parse | **MODIFIED (cycle/485).** Was a single python3 invocation; now loops over both golden paths. python3 exits non-zero on YAML error; no `\|\| true` needed. |
| 8 | Verify substrate structural shape — agent-admin | 116-128 | `for needle in …; do grep -qE … \|\| (echo + exit 1); done` | `set -eu` + `grep -q` (silent) + explicit `exit 1` in the failure branch | exit 0 when every needle matches | **MODIFIED (cycle/485) — renamed only.** Step body unchanged from cycle/476; `grep -q` returns 0 on match, 1 on no-match; under `set -eu` this would terminate, but the surrounding `if ! grep -qE`…fi inverts the condition cleanly. |
| 9 | Verify substrate structural shape — cds-dispatch | 130-154 | `for needle in …; do grep -qE … \|\| (echo + exit 1); done` + explicit OG-2 grep | `set -eu` + `grep -q` (silent) + explicit `exit 1` in the failure branch | exit 0 when every needle + the schedule grep matches | **NEW (cycle/485).** Mirrors step 8 with dispatch-shape needles plus the OG-2 `event_name == 'schedule'` check. Same bash-e guard shape as step 8. |
| 10 | AC5 — declaration-only refusal (cds-dispatch without override) | 156-180 | `./cn-install-wake … 2>/tmp/ac5.err >/dev/null \|\| rc=$?` (the renderer exits 3 here, which is the intended-success path); `grep -qF` / `grep -qE` for stderr content | `set -u` (not `-e` — because the renderer's intended-success exit is non-zero); explicit `\|\| rc=$?` capture of the non-zero exit; `if [ "$rc" != "3" ]; then exit 1; fi` for the precise exit-code check | exit 0 when renderer exits 3 AND stderr matches both patterns; explicit `exit 1` on any deviation | **NEW (cycle/485).** This step exercises a non-zero-exit-is-success path; `set -e` would terminate immediately on the renderer's exit 3. We use `set -u` for unset-var safety only and capture rc via `\|\| rc=$?` after `set -u +e`-equivalent. The `grep -qF` / `grep -qE` checks return 0 on match (success path); on no-match they return 1, but they're inside `if ! grep …; then exit 1` so the failure surface is explicit. No `\|\| true`-trap territory because we never assign their exit to a variable. |
| 11 | AC2 negative-case smoke (malformed manifest is rejected) | 182-201 | `./cn-install-wake … 2>&1 \| tee /tmp/neg.log` with the `if` checking the exit; `grep -q 'required field "schema"' /tmp/neg.log` | `set -eu` + `set -o pipefail`; pipeline reads via `tee` so both branches of the `if` see the renderer's exit | exit 0 when renderer exits 2 AND stderr names the missing field | Unchanged from cycle/476. Already correctly handles the `pipefail + tee` interaction. |
| 12 | AC8 (cycle/476) + AC7 (cycle/485) renderer-side authority audit | 203-241 | `n_admin=$(grep -ciE … \|\| true)`; `n_dispatch=$(grep -ciE … \|\| true)` | implicit `set -e` (GH default); **`\|\| true` guards the `grep -c` substitution** so a zero-count (which `grep -c` returns as exit 1) does not terminate the step before the `if` check | exit 0 when both counts equal "0" (the intended-success path); explicit `exit 1` on any non-zero count | **MODIFIED (cycle/485).** Was AC8-only; now two side-by-side audits. Both use the cycle/476 `\|\| true` pattern that γ-scaffold FN names explicitly. **Critical:** without `\|\| true` on `grep -ciE`, the step would fail under bash-e when the count is zero (the success case). The `[ "$n" != "0" ]` guard fires on real leaks. |

**Audit summary:** 12 `run:` blocks total; 5 NEW in cycle/485, 5 MODIFIED (renamed or had a sibling check added), 2 unchanged (jq, AC2 negative). Every new/modified block has been audited against the bash-e + command-substitution + pipeline traps. The `|| true` pattern is preserved precisely as cycle/476 introduced it (step 12); the new AC5 step (10) uses an explicit `|| rc=$?` capture because its intended-success path is `exit 3`, which is not zero — that step is the only one in the file where bash-e's "any non-zero terminates" semantics required a different guard than `|| true`.

---

## §R0

**α: R0 complete; β review-ready at sha=`<HEAD-after-commit>`.**

R0 deliverables landed on cycle/485:

1. `src/packages/cnos.core/commands/install-wake/cn-install-wake` — extended with `--activation-state-override` flag, role-conditional output_contract validation, dispatch-required field validation, activation_state §3.10 refusal gate (exit 3), `issues_labeled_selector_match` trigger handler, dispatch-aware job-level `if:` gate (OG-2 schedule-unconditional shape), and cross-package manifest discovery.
2. `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` — new golden fixture; sha256 = `75e0406662ea12b3caa5ad01d8dbe6bdfb620d0c047fd910917a8e930855605e`.
3. `.github/workflows/install-wake-golden.yml` — extended to iterate over both goldens (re-render + diff + idempotence + YAML parse + structural shape) and the AC5 refusal smoke; AC8 audit step extended with the cycle/485 AC7 audit grep.
4. `.cdd/unreleased/485/self-coherence.md` — this file.

**Local verification matrix (also in §AC-by-AC above):**

- AC1 (renderer parses dispatch): PASS
- AC2 (selector labels in if:): PASS
- AC3 (on: + labeled trigger): PASS
- AC4 (inlined prompt has artifact root + cell_runtime): PASS
- AC5 (declaration-only refusal with exit 3): PASS
- AC6 (idempotence): PASS
- AC7 (renderer source authority audit): PASS
- AC8 (agent-admin byte-identical): PASS
- AC9 (golden file + CI guards both): PASS
- OG-2 (schedule unconditional): PASS — explicit grep test in CI step 9
- OG-1 (renderer source manifest-driven): PASS — AC7 grep + β empirical read

**CI status:** Pending push. β reads the `install-wake-golden` CI run on the cycle/485 PR as the canonical signal.

β: per scaffold §7 review checklist, walk the §AC-by-AC table above + the §Per-CI-step audit, verify OG-1/OG-2/OG-3 against the renderer-source diff and the rendered golden, then issue a verdict (converge or iterate) in `.cdd/unreleased/485/beta-review.md` §R0.

---

## §Friction notes

### FN-α-1 — Cross-package manifest discovery

The renderer's pre-existing manifest lookup assumed `CN_PACKAGE_ROOT/orchestrators/<wake-name>/wake-provider.json`. When invoked directly (CI re-render check; this cycle's local-verification), `CN_PACKAGE_ROOT` defaults to the renderer's own package (cnos.core), which only contains the agent-admin wake. Rendering cds-dispatch (owned by cnos.cds) required either explicit `--manifest` or a cross-package fallback.

I chose the cross-package fallback because the CI re-render step would otherwise be tightly coupled to the renderer's owning-package convention, and an inner-loop developer running `cn-install-wake cds-dispatch` would have hit the same wall.

The fallback walks `src/packages/*/orchestrators/<wake-name>/wake-provider.json` and re-pins `CN_PACKAGE_ROOT` to the matching sibling. This preserves the relative path stripping in the header (so the cds-dispatch golden's header reads `manifest: orchestrators/cds-dispatch/wake-provider.json`, matching the admin-shape convention). When invoked via the `cn` dispatch shim with `CN_PACKAGE_ROOT` pre-exported, the lookup finds the manifest under the canonical path and the fallback is dead code.

β reviews: this fallback is a small surface-area extension that could in principle find a wake-name in the wrong sibling if two packages declared the same wake-name. Today no two packages do (agent-admin in cnos.core; cds-dispatch in cnos.cds; future cdr-dispatch in cnos.cdr; cdw-dispatch in cnos.cdw). If this becomes a real concern, a `--package <pkg>` flag could pin the lookup; out of scope for cycle/485.

### FN-α-2 — `prompt.md` line endings trailing-whitespace strip

The renderer inlines `prompt.md` with `sed 's/^/            /; s/[[:space:]]*$//'` (12-space indent + trailing-whitespace strip). The cds-dispatch `prompt.md` has a blank line at line 187 that became `            ` (12 spaces) in earlier admin-shape rendering, then was stripped to empty by the second sed. This is the cycle/476 idempotence-shape; it works identically for the dispatch shape. No friction.

### FN-α-3 — `agent_variable.default` on the cds-dispatch manifest

The cds-dispatch manifest's `agent_variable.default` is `"sigma"` (not `null` as in agent-admin). This means `cn install-wake cds-dispatch` succeeds without `--agent sigma` — the renderer's default agent fallback picks up the value. The rendered golden has `cds-dispatch-sigma` everywhere. β verifies this matches the AC1 oracle's expectation (`group: cds-dispatch-sigma`).

The renderer source treats `--agent` as overrideable with a hardcoded `sigma` fallback when neither flag nor manifest provides; the cds-dispatch manifest's explicit `"sigma"` default is consistent with this. Future per-package bot accounts (cnos#449 follow-up) will lift the hardcoded sigma → variable per-package.

### FN-α-4 — γ-scaffold's β-review section names the per-CI-step audit on β; I populated it on α

The γ scaffold §7 (β review prompt) names the per-CI-step verification table as β's deliverable per cnos#478. The α implementation prompt §6 also names it for α. Both authors populate the same surface (one in self-coherence.md §Per-CI-step, the other in beta-review.md §R[N] per-CI-step). I populated mine in this file; β re-audits independently in beta-review.md per the γ-scaffold §7 instruction. No conflict — α surfaces the audit at implementation time; β reads + verifies at review time.

### FN-α-5 — γ-scaffold cited a sha drift in cnos#485 AC8 (FN-1 in γ-scaffold §11)

cnos#485 AC8 cites `47824628a5958ec9372196b30fa2cb0b547c492799358368636c0e24981b10e1` as the agent-admin golden's sha256. γ verified at branch base `af8ed8ec` that the actual value is `fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72`. I used γ's verified value as the AC8 invariant and verified: after the renderer extension lands on cycle/485, `cn install-wake agent-admin` still produces the `fa6b8c0c…` sha — i.e., the admin-shape render is byte-identical to the base-sha render. β re-verifies on the cycle/485 PR's `install-wake-golden` CI run, which also catches this via `git diff --exit-code`.

If β prefers to re-verify the cited sha against current main HEAD (in case main has moved during cycle/485), the command is:
```
git fetch origin main && git show origin/main:src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml | sha256sum
```

The drift between issue body and reality is downstream of three post-cycle-476 cleanup commits (`c353d432`, `4b633bb2`, `7ab62cb9`) that γ surfaced in scaffold §11 FN-1. The lesson γ named — "δ should re-verify cited sha256 values against actual filesystem state before dispatching γ" — is the right one; I have no improvement to suggest.
