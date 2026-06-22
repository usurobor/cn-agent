---
cycle: 485
parent_issue: cnos#485
master_tracker: cnos#467 (Sub 5A)
cycle_branch: cycle/485
authored_by: β@cdd.cnos (bootstrap-δ session)
date: 2026-06-22 (UTC)
pr: https://github.com/usurobor/cnos/pull/488
head_sha: 9b9ac165aaf3fb149712ef748e687530d22ba2df
alpha_signal_sha: 550a6907320d64c9032116c70a0526787cab52d1
agent_admin_golden_sha256_invariant: fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72
cds_dispatch_golden_sha256: 75e0406662ea12b3caa5ad01d8dbe6bdfb620d0c047fd910917a8e930855605e
verdict: converge
---

# β review — cnos#485 (cn-install-wake renderer extension for dispatch wake providers)

## §R0 — Verdict

**verdict: converge.** All 9 ACs verified independently with passing oracles; OG-1/OG-2/OG-3 each pass empirical inspection; per-CI-step bash-e audit clean (12 steps audited; α's audit matches β's independent audit row-for-row); doctrine consistency intact (zero `cnos.cdd` literals in renderer source; zero stale `protocol:cdd` / `cdd-dispatch` references); non-goals preserved (production cds-dispatch manifest still `declaration-only`; agent-admin golden byte-identical at `fa6b8c0c…`); install-wake-golden CI on the PR HEAD shows the `Re-render + diff per-package goldens` job green; the three concurrent red checks (I4/I5/I6) are inherited caps per the cnos#467 wave inventory, not cycle/485-introduced. No blockers; no findings.

---

## §Per-AC verification table

All commands run on `/home/user/cnos`, branch `cycle/485`, HEAD `9b9ac165`. β re-ran every oracle independently (not just read α's claim).

| # | AC statement | Oracle command(s) | Observed result | Verdict |
|---|---|---|---|---|
| AC1 | Renderer parses `role: dispatch`; emits YAML with `name: cnos-cds-dispatch`, `group: cds-dispatch-sigma`, all 4 permissions | `./cn-install-wake cds-dispatch --activation-state-override live; echo $?`; grep `^name:`, `^  group:`, each permission line | rc=0; WARNING printed to stderr; output "(unchanged)" (file already at golden); `name: cnos-cds-dispatch` present; `  group: cds-dispatch-sigma` present; `contents: write`, `issues: write`, `pull-requests: write`, `id-token: write` all matched | **PASS** |
| AC2 | `protocol` + `selector` emitted; each include label appears in a `contains(github.event.label.name, ...)` construct | `for l in dispatch:cell protocol:cds status:todo; do grep -qF "$l" $f; done; grep -oE "contains\(github\.event\.label\.name, '[^']+'\)" $f` | All 3 include labels present; 3 distinct `contains(github.event.label.name, '<label>')` clauses emitted (one per label), wrapped in parentheses inside the issues conjunction | **PASS** |
| AC3 | `on:` block has `issues: { types: [labeled] }` + `schedule:`; job-level `if:` references both | `python3 yaml.safe_load` → `on.keys()` = `['schedule', 'issues']`; `on.issues.types` includes `'labeled'`; schedule has 4 cron entries | YAML parses cleanly; both top-level triggers present; `labeled` in types list; job-level `if:` references `github.event_name == 'schedule'` AND `github.event_name == 'issues'` AND `github.event.label.name` | **PASS** |
| AC4 | Inlined prompt has `.cdd/unreleased/{N}/` + `cnos.cdd`; rendered YAML carries the `prompt: \|` block-scalar | `grep -F '.cdd/unreleased/{N}/' $f; grep -F 'cnos.cdd' $f; grep -E "^[[:space:]]+prompt: \|" $f` | `.cdd/unreleased/{N}/` present (multiple lines from prompt.md); `cnos.cdd` present 7 times in the inlined prompt; `          prompt: \|` present at the start of the prompt block | **PASS** |
| AC5 | `cn install-wake cds-dispatch` (no override) → exit 3; stderr names `declaration-only` + activation_state_notes content | `./cn-install-wake cds-dispatch 2>/tmp/ac5.err; echo $?; grep -qF 'declaration-only' /tmp/ac5.err; grep -qE 'cnos#454\|cnos#467\|preconditions' /tmp/ac5.err` | rc=3; stderr names `activation_state="declaration-only"`; stderr inlines the entire `activation_state_notes` field which contains `cnos#454`, `cnos#467`, AND `preconditions`; stderr names the test-only escape hatch with full `--activation-state-override` flag name | **PASS** |
| AC6 | Two consecutive renders of cds-dispatch produce byte-identical output | First `--override live` render, capture sha; second `--override live` render, capture sha; compare | `sha_a=sha_b=75e0406662ea12b3caa5ad01d8dbe6bdfb620d0c047fd910917a8e930855605e`; second invocation reports "(unchanged)" — α's idempotent-write logic correctly detects no-diff and does not rewrite | **PASS** |
| AC7 | `grep -ciE 'protocol:cds\|cdr\|cdw\|dispatch:cell\|status:todo' cn-install-wake` returns 0 | `n=$(grep -ciE … \|\| true); [ "$n" = "0" ]` | n=0; renderer source carries zero dispatch-shape role-decision literals. β's empirical line-by-line read (see OG-1 below) confirms labels are read via `jq -r ".selector.include[$j]" "$manifest_path"` (line 597), not from any shell variable concatenation | **PASS** |
| AC8 | `cn install-wake agent-admin` produces sha256 = `fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72` (verified value, not issue's stale `47824628…`) | `./cn-install-wake agent-admin; sha256sum src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml \| cut -d' ' -f1` | `fa6b8c0cd64fb626a5e1e991128cbb27fb883b6d1594914543032a2b0d2d3e72` — exact match to the verified-at-branch-base value; second invocation reports "(unchanged)"; `git diff --exit-code` on the file returns 0 | **PASS** |
| AC9 | cds-dispatch golden exists; CI workflow iterates over both goldens; YAML parses for both | `test -f $f; grep -F 'cnos-agent-admin.golden.yml' install-wake-golden.yml; grep -F 'cnos-cds-dispatch.golden.yml' install-wake-golden.yml; grep -cE 'sha256sum.*\.golden\.yml' install-wake-golden.yml` | golden file present; both fixtures named in the CI workflow (in step 4 `git diff` and steps 5/6 idempotence); `sha256sum.*\.golden\.yml` occurs 4 times (2 idempotence checks × 2 goldens); `Verify YAML parses (both goldens)` step loads both via `yaml.safe_load` in a single python3 invocation | **PASS** |

**9/9 ACs verified.**

---

## §Operator-guardrail verification

### OG-1 — Renderer source must remain package-driven, NOT hard-coded

**Verdict: PASS.**

**β empirical read.** Per γ FN-4: the AC7 grep is defeatable by shell-variable concatenation. β walked every line in `cn-install-wake` that touches `protocol`, `dispatch`, `selector`, or `status` content (renderer source @ HEAD `9b9ac165`):

- **Label values flow only through jq.** Line 597: `label_value="$(jq -r ".selector.include[$j]" "$manifest_path")"`. The OR-chain emitted into the rendered `if:` clause is constructed entirely from this value (line 598: `clause="contains(github.event.label.name, '${label_value}')"`). No shell-variable substitute path for the label values exists.
- **`protocol` field used only for type validation.** Lines 340-342 read `.protocol | type` from the manifest to enforce string-type required-field check; the *value* of `.protocol` is never emitted into the rendered YAML by the renderer. (The rendered prompt body contains `protocol:cds` references, but those come from the inlined `prompt.md` content, which is package authority — the renderer flows it through verbatim with only `{agent}` substitution.)
- **`dispatch` substring in the source.** All 26 occurrences of the substring `dispatch` in the renderer source are either (a) comments referencing the dispatch shape conceptually, (b) the role-enum value `dispatch` in the role-validation `case` switch (line 318) which is a *role-enum literal*, not a label-emission literal — this is exactly what wake-provider/SKILL.md §2.1 names as the package-authority value the renderer dispatches on, or (c) variable names like `has_issues_labeled`. None emit `dispatch:cell` into rendered output.
- **`status` substring.** Zero occurrences in the renderer source (`grep -n status …` returned empty).
- **Concat pattern audit.** Three concat-variable assignments at lines 283-285: `_a="admin"; _u="_only"`, `_d="dis"; _as="allowed_surfaces"`, `_p="defer"; _path="_path"`. These predate cycle/485 (cycle/476 §Debt) and assemble admin-shape *manifest-field names* (`admin_only`, `disallowed_surfaces`, `defer_path`) for the `require_field` validation loop — they are NOT role-decision label-strings, and α did NOT add any analogous concat pattern for the dispatch shape (no `_d="dis"; _c="patch:cell"` or similar). The dispatch-shape required-field validation (lines 340-350) names fields directly without concat (`.protocol`, `.selector.include`, `.output_contract.cycle_artifact_root`, etc.) — these are field-path expressions, not label literals AC7 catches.

**Conclusion.** No literal-string trickery; no obfuscation via concat. The renderer's dispatch-shape behavior is genuinely manifest-driven. AC7's grep oracle returns 0, AND β's empirical read confirms genuine manifest-read at line 597.

### OG-2 — `issues_labeled_selector_match` MUST NOT suppress the scheduled sweep

**Verdict: PASS.**

**The `if:` line in the rendered cds-dispatch golden (line 35):**
```
    if: ${{ github.event_name == 'schedule' || (github.event_name == 'issues' && (contains(github.event.label.name, 'dispatch:cell') || contains(github.event.label.name, 'protocol:cds') || contains(github.event.label.name, 'status:todo'))) }}
```

This is the exact shape γ's scaffold §6 named as correct. Mental simulation:
- **Schedule firing:** `github.event_name == 'schedule'` evaluates true; the OR short-circuits; the entire expression is true regardless of `github.event.label` (which is undefined at schedule time). **Job fires.**
- **Labeled-issue event with selector match:** `github.event_name == 'schedule'` false; `github.event_name == 'issues'` true; `contains(github.event.label.name, '<one of dispatch:cell/protocol:cds/status:todo>')` true. **Job fires.**
- **Labeled-issue event with non-selector label:** `github.event_name == 'issues'` true; all three `contains(...)` false. Expression false. **Job does not fire.**
- **Any other event (push, workflow_dispatch, etc.):** Neither schedule nor issues. Expression false. **Job does not fire** — which is correct; the workflow has no other triggers in its `on:` block.

The renderer source at lines 622-629 explicitly distinguishes the admin shape (where the schedule short-circuit + `contains(github.event.issue.title, ...)` is sufficient because `issue.title` is undefined at schedule firing time and evaluates to non-match) from the dispatch shape (where `event_name == 'issues'` explicit-gate is necessary because `github.event.label` is undefined at schedule time). α's comment block lines 558-566 explicitly cites γ's OG-2 directive. CI step 9 has an explicit `grep -qE "github\.event_name == 'schedule'"` regression test that would fail loudly if this shape regresses.

### OG-3 — `activation_state: declaration-only` override mechanism must be visibly test-only

**Verdict: PASS.**

α chose Option A (`--activation-state-override <state>` CLI flag). β audited along five dimensions:

1. **Documented in `--help`?** Yes. The flag and its TEST-ONLY semantics are documented in the `Arguments:` section of the header comment block (lines 39-54), which is what `--help` prints (the help handler at line 174 reads the header docstring via `sed -n '2,/^$/p' "$0"`). The help text says verbatim: "TEST-ONLY override … The flag name itself names what is being bypassed; a stderr warning is emitted on every use. Never pass this flag during production install".
2. **WARNING emitted on every use?** Yes. Line 364 unconditionally emits the WARNING whenever `$activation_state_override` is non-empty, naming both the override value AND the manifest's declared value — visible audit trail in any CI log or operator terminal.
3. **Restricted to bypass-`live`-only?** Yes. Line 370's check `if [ "$activation_state_effective" != "live" ]` means any override value OTHER than `live` still results in refusal (exit 3) + WARNING. β verified empirically: `--activation-state-override renderer-pending` → WARNING + exit 3; `--activation-state-override declaration-only` → WARNING + exit 3; `--activation-state-override live` → WARNING + render. So the flag's *only* effective use is `--activation-state-override live`, which is precisely the test/CI mode α/γ documented.
4. **Empty-string override safe?** Yes. β tested `--activation-state-override ""`. Line 363's `[ -n "$activation_state_override" ]` is false on empty string; the WARNING is suppressed and the manifest's declared value flows through (still `declaration-only` → exit 3). No silent bypass via empty string.
5. **Discoverability for an operator?** Yes. The flag name `--activation-state-override` is self-documenting (names what is being bypassed); the WARNING is impossible to miss in CI logs; the help-text TEST-ONLY warning is explicit. The renderer also tells the operator about the flag in the refusal error itself (line 381: "to render anyway for test/diff-review purposes, pass --activation-state-override live (test-only; see --help).").

**Could a normal operator accidentally use it?** No. There is no environment variable, no implicit gating, no fixture-path indirection. The operator must explicitly pass an explicit flag with an explicit value, and every use leaves a WARNING in the substrate run logs.

---

## §Per-CI-step bash-e audit (independently re-populated)

`.github/workflows/install-wake-golden.yml` has 12 `run:` blocks. β re-audited each independently from α's table. Discrepancies with α's audit: **none.** All α's per-step characterizations match β's independent reading.

| # | Step name | Line range | Command substitutions / pipelines used | Guarded? | bash -e exit on intended-success input | Notes |
|---|---|---|---|---|---|---|
| 1 | Verify jq present | 37-41 | none (`jq --version` only) | implicit `set -e` (GH runner default) | exit 0 on success | Unchanged from cycle/476. No risk surface. |
| 2 | Re-render agent-admin wake (cnos.core) | 43-45 | none | implicit `set -e` | exit 0 on success | Unchanged. |
| 3 | Re-render cds-dispatch wake (cnos.cds, `--activation-state-override live`) | 47-55 | none | implicit `set -e` | exit 0 on success | **NEW.** Renderer's WARNING goes to stderr; renderer exits 0 because override='live' is the only value that bypasses refusal. β verified locally: rc=0. |
| 4 | Verify goldens unchanged | 57-69 | none; `git diff --exit-code` returns 0 on no-diff, 1 on diff | `if ! git diff …; then exit 1; fi` — non-zero exit consumed by `if !` | exit 0 when goldens match; explicit `exit 1` only on actual drift | **MODIFIED.** Single `git diff` invocation lists both goldens. The `if !` inversion is the correct guard pattern for bash-e — without it, the diff-found exit would terminate the step before the error annotation is emitted. |
| 5 | Verify idempotence — agent-admin | 71-84 | `$(sha256sum file \| cut -d' ' -f1)` × 2 | implicit `set -e`; both `sha256sum` and `cut` are exit-0 on the well-formed file inputs they see (file exists from step 2's re-render) | exit 0 when shas match; explicit `exit 1` on drift | **MODIFIED (renamed only).** No pipefail set, but neither side of the pipe can return non-zero on intended-success input. Safe. |
| 6 | Verify idempotence — cds-dispatch | 86-100 | `$(sha256sum file \| cut -d' ' -f1)` × 2 | implicit `set -e`; same shape as step 5 | exit 0 when shas match; explicit `exit 1` on drift | **NEW.** Mirrors step 5; uses `--activation-state-override live` because the manifest is declaration-only. Same safe pipe shape. |
| 7 | Verify YAML parses (both goldens) | 102-114 | inline python3 heredoc; `yaml.safe_load` raises on parse error | implicit `set -e`; python3 exits non-zero on uncaught exception | exit 0 when both parse cleanly | **MODIFIED.** Loops over both golden paths. Single python3 invocation; PyYAML available on ubuntu-latest by default. No substitution-into-shell guards needed. |
| 8 | Verify substrate structural shape — agent-admin | 116-128 | `for needle in …; do if ! grep -qE "$needle" "$f"; then exit 1; fi; done` | `set -eu` + `grep -q` (silent) wrapped in `if !` — non-zero from grep on no-match consumed cleanly | exit 0 when every needle matches | **MODIFIED (renamed only).** Step body unchanged from cycle/476. |
| 9 | Verify substrate structural shape — cds-dispatch | 130-154 | `for needle in …; do if ! grep -qE …; then exit 1; fi; done` plus explicit OG-2 `if ! grep -qE "github\.event_name == 'schedule'"` | `set -eu`; all `grep -qE` calls inside `if !` | exit 0 when every needle + the schedule grep matches | **NEW.** Mirrors step 8 with dispatch-shape needles + the explicit OG-2 schedule-event regression test. Same safe shape. |
| 10 | AC5 — declaration-only refusal | 156-180 | `./cn-install-wake … 2>/tmp/ac5.err >/dev/null \|\| rc=$?` (renderer's intended-success exit is **3**); subsequent `if ! grep -qF …; then exit 1` checks | `set -u` (NOT `-e`, deliberately, because intended-success exit is non-zero); `rc=$?` captures the non-zero exit cleanly; subsequent commands (echo/cat/grep) are either inside `if !` or have no failure path on intended input | exit 0 when renderer exits 3 AND stderr matches both patterns | **NEW.** This is the one step where α deliberately uses `set -u` only — the renderer's intended-success exit is non-zero, and `set -e` would terminate immediately on the renderer's exit 3. The `\|\| rc=$?` pattern is exactly the right idiom for "I expect this to fail; capture the failure code for inspection". β re-verified: `set -u` without `-e` means subsequent unguarded commands won't terminate on failure — but all subsequent commands in step 10 are either `echo`/`cat` (no failure path) or inside `if !` guards, so no hidden trap. The cycle/476 class-trap shape (a substitution exit propagating before the inspection runs) is not present here. |
| 11 | AC2 negative-case smoke | 182-201 | `if ./cn-install-wake … 2>&1 \| tee /tmp/neg.log; then …`; `tee` always exits 0; renderer is expected to exit 2 (manifest error) | `set -eu` + `set -o pipefail` + `if !`-style — pipefail makes the pipeline return the renderer's non-zero exit; `if` consumes it | exit 0 when renderer exits non-zero AND error message names "schema" | Unchanged from cycle/476. The pipefail + tee + `if` interaction is correctly composed. |
| 12 | AC8 (cycle/476) + AC7 (cycle/485) renderer-side authority audit | 203-241 | `n_admin=$(grep -ciE … \|\| true)`; `n_dispatch=$(grep -ciE … \|\| true)` | implicit `set -e`; **`\|\| true` neutralizes `grep -c`'s exit 1 on zero matches** (the intended-success path) | exit 0 when both counts are "0" (zero leaks); explicit `exit 1` on any non-zero count | **MODIFIED.** Was AC8 only; now two side-by-side audits. Both use the cycle/476 `\|\| true` discipline. **This is the cycle/476 3-round class-trap location.** β verified the `\|\| true` is correctly placed on both substitutions; without it, the substitution exit (1 on zero matches) would terminate the step under bash-e before the diagnostic `[ "$n" != "0" ]` check fires. The comment block lines 217-224 documents this discipline explicitly. |

**Audit summary.** 12 `run:` blocks, 5 NEW in cycle/485 (steps 3, 6, 9, 10; step 12 added a new sibling grep), 5 MODIFIED (renamed or had a sibling added), 2 unchanged (steps 1, 11). Every NEW/MODIFIED block has been audited row-by-row by β independently of α's claims; α's audit and β's audit agree on every column. The cycle/476 `|| true` pattern is preserved precisely (step 12); step 10's deliberate non-`-e` posture is the right idiom for a step whose intended-success exit is non-zero. **No bash-e regressions detected.**

---

## §Doctrine consistency check

**Verdict: PASS.**

Per PR #480 doctrine correction + PR #466 dispatch-protocol skill + PR #483 Sub 4: the renderer extension MUST NOT hard-code `cnos.cdd` as the cell runtime in the rendered workflow's prompt. Runtime selection is via the manifest's `cell_runtime` field.

- **Renderer source `cnos.cdd` occurrences:** ZERO. `grep -n 'cnos\.cdd' cn-install-wake` returned empty. The renderer does not hard-code any cell-runtime target.
- **Rendered cds-dispatch golden `cnos.cdd` occurrences:** 7 lines. β read each: all are inside the inlined `prompt: |` block-scalar (i.e., from `prompt.md` content), and all are in the form `cnos.cdd cell-runtime framework's δ role contract` or similar prose that names the framework correctly per cell_runtime authority. None are injected by the renderer source.
- **Stale references audit (`protocol:cdd`, `cdd-dispatch`, `cnos.cdd dispatch wake`):** ZERO new `+` lines in the cycle/485 diff carry any of these stale forms. (Re-ran the grep on the full `+` set across renderer + golden + CI workflow.)
- **`cell_runtime: cnos.cdd` inlined into the rendered prompt:** Yes — α correctly relied on `prompt.md` containing the framework name; the renderer flows it through verbatim. This is exactly the cell_runtime authority shape (package authority sets the runtime name; renderer is substrate-only).

---

## §CI evidence

PR #488 HEAD `9b9ac165` triggered the install-wake-golden CI workflow. β read check_runs via `mcp__github__pull_request_read`:

- **`Re-render + diff per-package goldens`** (the install-wake-golden main job): **success** (run https://github.com/usurobor/cnos/actions/runs/27977271314/job/82798098377). This is the canonical green signal that AC1+AC2+AC3+AC5+AC6+AC7+AC8+AC9 all hold in the substrate environment.
- **`Go build & test`:** success.
- **`Package/source drift (I1)`:** success.
- **`Protocol contract schema sync (I2)`:** success.
- **`Package verification`:** success.
- **`Binary verification`:** success.
- **`Repo link validation (I4)`:** **failure** — inherited cap per the cnos#467 wave inventory (PR #483 comment-4770338676 I4/I5/I6 are pre-existing problems unrelated to this cycle).
- **`SKILL.md frontmatter validation (I5)`:** **failure** — inherited cap (I5).
- **`CDD artifact ledger validation (I6)`:** **failure** — inherited cap (I6).

The three red checks correspond exactly to the three inherited caps named in the wave inventory and accepted as not-this-cycle's-fault. The named checks β was instructed to focus on (per the dispatch prompt) are all green. **CI evidence aligns with local AC verification.**

---

## §Non-goal verification

`git diff main...HEAD --name-only` returns exactly 5 files:

```
.cdd/unreleased/485/gamma-scaffold.md
.cdd/unreleased/485/self-coherence.md
.github/workflows/install-wake-golden.yml
src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
src/packages/cnos.core/commands/install-wake/cn-install-wake
```

None of the forbidden paths (γ scaffold §4 "α MUST NOT touch") appear:

- `src/packages/cnos.cds/orchestrators/cds-dispatch/wake-provider.json` — NOT touched. β confirmed via `jq -r '.activation_state' wake-provider.json` → still `declaration-only`, as required.
- `src/packages/cnos.cds/orchestrators/cds-dispatch/prompt.md` — NOT touched.
- `src/packages/cnos.core/orchestrators/agent-admin/*` — NOT touched (any modification would have broken the AC8 byte-identical invariant; β verified the working-tree golden's sha is unchanged from base).
- `.github/workflows/cnos-agent-admin.yml` — NOT touched (live admin wake preserved).
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` — NOT touched (Sub 5B scope).

α stayed strictly in cycle/485 scope.

---

## §Findings

**No findings.** β did not identify any blocker or non-blocking issue. The work is converge-ready.

---

## §Friction notes

These are β's meta-observations for future δ template / cycle improvements. None are blockers; none require action from α/γ this cycle.

**FN-β-1.** α correctly inherited γ's verified agent-admin sha256 (`fa6b8c0c…`) rather than the issue body's stale `47824628…` value. The "verify cited-sha against filesystem state before relying on it" discipline that γ surfaced in scaffold §11 FN-1 is the right one; α's self-coherence §FN-α-5 cites it explicitly. β verified the value is still correct at HEAD (`9b9ac165`) by reading the file directly: matches `fa6b8c0c…`. Cycle/485 is the third independent verification of this value (γ at base, α post-extension, β at PR HEAD); the discipline is producing reliable handoffs.

**FN-β-2.** The per-CI-step bash-e audit α populated in self-coherence.md §Per-CI-step is high quality. β's independent re-audit agreed on every cell of every row. This is the first cycle since cnos#478 mechanical-injection mandate where the audit table is populated at *both* α implementation time AND β review time, and the two populations agree. The discipline appears to have absorbed the cycle/476 3-round class-trap lesson; consider this empirical evidence that the per-CI-step table is now reliably catching the trap class it was introduced to prevent.

**FN-β-3.** OG-3 Option A (CLI flag with WARNING + exit-3 on non-`live`) is a strong "visibly test-only" implementation. The five-property check β walked (help-documented; WARNING-on-every-use; bypass-`live`-only; empty-string-safe; refusal-quoting-the-flag) is a concrete operationalization of γ's "is this genuinely test-only?" question. Future γ scaffolds in this template family might pin Option A as the preferred mechanism and embed the five-property β-checklist directly rather than leaving the mechanism choice open with three options — the open choice consumed α's reasoning time, β's review time, and the AC has converged on Option A anyway.

**FN-β-4.** No bash-e regressions detected. The `|| true` discipline on `grep -c` substitutions in step 12 + the deliberate `set -u`-without-`-e` posture in step 10 are both correctly applied. α/γ should consider FN-β-2 + this observation as evidence that the per-CI-step table format is doing its job and merits being hoisted into a shared β-skill (per γ FN-6 / FN-8) so future cycles can `@include` rather than restate.

**FN-β-5.** The CI green-vs-red accounting (CI evidence section above) was straightforward: three inherited-cap reds, six greens including the canonical signal. The "inherited cap" framing γ named in the dispatch prompt is the right pattern — β can confidently issue converge despite red checks because the red checks are pre-existing and traceable to the wave inventory. The δ wake-invoked mode skill (cnos#486) should formalize this affordance: β verdicts are not gated on inherited-cap failures.
