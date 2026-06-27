# γ R0 Scaffold — cnos#503

**Cycle**: cycle/503
**Protocol**: cds
**Mode**: substantial (design + build)
**γ**: sigma@cnos.cn-sigma.cnos
**Scaffold rev**: R0
**Base SHA**: 111c29701f41c65ce87c44d3304d4a9eabf61ac5

---

## 1. Issue — cnos#503

**Title**: wake-trace: persist Claude execution logs as artifacts + cnos run trace

**Gap**: The rendered wake workflows (`cnos-cds-dispatch.yml`, `cnos-agent-admin.yml`) invoke
`anthropics/claude-code-action@v1` without:

- Giving the Claude step an `id:` (so subsequent steps cannot reference
  `steps.<id>.outputs.execution_file` or `steps.<id>.outputs.session_id`)
- Uploading the execution file as a workflow artifact
- Writing a durable run-trace pointer after meaningful work

The action exposes both `execution_file` and `session_id` outputs, but the rendered substrate does
nothing with them. The execution file exists during the run and disappears when the job ends.

**Two-stage proof**:

- Stage 1: implement + PR merge (all 6 ACs verified via CI fixtures + golden checks)
- Stage 2: post-merge live smoke (trigger a smoke cell; verify trace appears; verify no-op sweep
  produces no trace)

Do NOT close cnos#503 until Stage 2 evidence lands.

**Trace-location split**:

| Wake class | Durable trace pointer location |
|---|---|
| Cell dispatch wakes (cds-dispatch) | `.cdd/unreleased/{issue}/wake-trace.md` |
| Agent-admin wakes | Admin-owned surface (NOT in package cell directory) |
| Raw Claude execution logs (all) | Short-retention workflow artifacts only; never committed |

---

## 2. Mode

**Substantial (design + build)** — the implementation requires:

- Extending the renderer (`cn-install-wake`) to emit three new YAML blocks in both rendered
  workflows: (a) `id: claude` on the action step, (b) a conditional `upload-artifact` step,
  (c) a conditional wake-trace write + commit step.
- Updating both golden fixtures to match the new renderer output.
- Re-rendering deployed workflows (`cnos-cds-dispatch.yml`, `cnos-agent-admin.yml`) so they are
  byte-identical with goldens.
- Authoring AC-specific CI fixtures for AC1, AC2, AC3, AC6 in `install-wake-golden.yml`.
- AC5 is documentation-only: updating the stale-claim recovery runbook or leaving a note for
  cnos#504 that the trace surface exists.

The renderer must remain free of role-decision strings (AC7/AC8 audit continues to pass).

---

## 3. Surfaces α will touch

| Surface | Location | Change |
|---|---|---|
| Renderer | `src/packages/cnos.core/commands/install-wake/cn-install-wake` | Add `id:` to Claude step; add upload-artifact step; add wake-trace write step |
| CDS dispatch golden | `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` | Re-rendered (three new sections) |
| Agent-admin golden | `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` | Re-rendered (three new sections, admin-surface trace path) |
| Deployed cds-dispatch | `.github/workflows/cnos-cds-dispatch.yml` | Re-rendered; byte-identical with golden |
| Deployed agent-admin | `.github/workflows/cnos-agent-admin.yml` | Re-rendered; byte-identical with golden |
| CI golden check | `.github/workflows/install-wake-golden.yml` | New steps for AC1, AC2, AC3, AC6 oracles |
| Trace template (optional) | renderer inline or shell here-doc | The fields written to wake-trace.md |
| AC5 doc note | `.cdd/unreleased/503/wake-trace.md` or in scaffold | Note for cnos#504 naming trace surface |

α MUST NOT:
- Hand-edit golden YAML files (use re-render only)
- Modify `anthropics/claude-code-action@v1` itself
- Commit raw JSON execution logs to the repo
- Write to `.cn-sigma/logs/` or other admin surfaces

---

## 4. AC Oracle Approach

**AC1 — Claude step exposes outputs**

The rendered YAML must have `id: claude` on the `uses: anthropics/claude-code-action@v1` step.
Oracle: `grep -nE '^\s+id: claude$'` in each golden file; assert non-empty. The CI golden-check
workflow adds a step asserting this grep succeeds on both golden YAMLs. The renderer emits `id:
claude` unconditionally for both wake classes (admin and dispatch). The `id:` field is additive to
the existing `uses:` block — it inserts above the `with:` block in YAML, ordering: `id:`, `uses:`,
`with:`.

**AC2 — Execution file uploaded**

The rendered YAML must contain an `actions/upload-artifact@v4` step whose `if:` condition gates on
`steps.claude.outputs.execution_file != ''`. Oracle: `grep -nE 'upload-artifact@v4'` in each
golden file; assert non-empty. Additionally: grep for `steps.claude.outputs.execution_file` to
confirm the gate is wired correctly. The step must specify `retention-days: 14`. Existing CI
idempotence + sha256 checks continue to pass because the golden is re-rendered from source.

**AC3 — Durable trace record**

The rendered YAML must contain a wake-trace write step whose output location is
wake-class-appropriate:
- cds-dispatch: writes `.cdd/unreleased/$ISSUE_NUMBER/wake-trace.md`
- agent-admin: writes an admin-owned trace path (e.g. `.cn-sigma/logs/wake-trace-$GITHUB_RUN_ID.md`
  or `.cn-sigma/state/wake-traces/$GITHUB_RUN_ID.md`)

Required fields in the written file: `run_id`, `wake_class`, `claude_session_id`,
`execution_artifact`, `workflow`, `conclusion`, `last_known_phase`. For cell-dispatch:
additionally `issue` and `protocol`.

Oracle: grep each golden for the required field names in the trace write step's shell block. The
step has an AC6-gating `if:` condition — see AC6 below. CI asserts the field list is present in
the golden.

**AC4 — Privacy/retention guard**

No raw execution log JSON is committed to the repo. The upload-artifact step targets the path from
`steps.claude.outputs.execution_file` only (the file the action itself produced; no raw content
expansion in commits). The `retention-days: 14` value is asserted in the golden. The wake-trace
write step writes only the structured metadata fields (run_id, session_id, etc.) — not the raw
execution log content.

Oracle: grep golden for `retention-days: 14`; assert no raw JSON expansion in the trace write
block (no `cat $execution_file` redirected into a committed file).

**AC5 — Stale-claim recovery integration**

Documentation-only AC (no code gate). The cnos#504 runbook (or a note here) names the
`.cdd/unreleased/{N}/wake-trace.md` surface as the diagnostic starting point for stale-claim
recovery. This AC is satisfied by the trace surface existing (from AC3) and a reference in the
cycle/503 scaffold or a note committed to `.cdd/unreleased/503/` naming the cnos#504 dependency.

Oracle: at Stage 1, verify `.cdd/unreleased/503/` contains a note naming cnos#504 and the trace
surface. β confirms the note is present and readable.

**AC6 — No-op sweep produces no committed trace**

The wake-trace write + commit step must have an `if:` guard that prevents execution when the
selector scan was empty (no meaningful work was done). For cds-dispatch: check whether
`.cdd/unreleased/` gained new paths in this run — the predicate is:

```sh
git log "$CN_WAKE_BASE_SHA..HEAD" --name-only -- '.cdd/unreleased/' | grep -v '^$'
```

Non-empty = meaningful work done; write the trace. Empty = no-op firing; skip the trace step
entirely. For agent-admin: check whether HEAD advanced from baseline — if HEAD == CN_WAKE_BASE_SHA
(heartbeat no-op), skip the trace; if HEAD != CN_WAKE_BASE_SHA (any commit was made), write the
trace.

Oracle: trigger a no-op cds-dispatch firing (empty queue) in CI fixture; verify
`git log -- .cdd/unreleased/` shows no new wake-trace commits. CI adds a fixture-based smoke step
to `install-wake-golden.yml` that simulates the no-op condition and asserts the trace step
predicate returns "skip".

---

## 5. Source-of-Truth Table

| Artifact | Governing reference |
|---|---|
| Renderer logic | `src/packages/cnos.core/commands/install-wake/cn-install-wake` (source of truth; all downstream artifacts are derived) |
| CDS dispatch golden | `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` (rendered artifact; must match renderer byte-for-byte) |
| Agent-admin golden | `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` (rendered artifact; must match renderer byte-for-byte) |
| Deployed cds-dispatch workflow | `.github/workflows/cnos-cds-dispatch.yml` (deployed artifact; must be byte-identical with cds golden per sha256 check in CI) |
| Deployed agent-admin workflow | `.github/workflows/cnos-agent-admin.yml` (deployed artifact; must be byte-identical with agent-admin golden — CI will check once this AC lands) |
| CI golden-check workflow | `.github/workflows/install-wake-golden.yml` (test harness; adds AC1/AC2/AC3/AC6 oracle steps) |
| Wake-provider schema | `cn.wake-provider.v1` (NOT extended; trace machinery is renderer authority only) |
| AC7/AC8 audit pattern | `install-wake-golden.yml` step "AC8 (cycle/476) + AC7 (cycle/485) renderer-side authority audit" — grep: `admin.only\|disallowed_surfaces\|defer.path\|cell_execution` and `protocol:cds\|cdr\|cdw\|dispatch:cell\|status:todo` must not appear in renderer source |
| Activation-log write fence | Rendered by renderer as final step on dispatch wakes; unchanged by this cycle |

---

## 6. Implementation Contract (pinned by δ — α MUST NOT improvise any axis)

| Axis | Pinned value |
|---|---|
| Language | POSIX sh (renderer changes are shell) + YAML (rendered output) |
| CLI integration target | `cn install-wake` (existing command, extended); re-render via `./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch` (golden) and `./src/packages/cnos.core/commands/install-wake/cn-install-wake agent-admin` (golden), then `--out .github/workflows/cnos-cds-dispatch.yml` and `--out .github/workflows/cnos-agent-admin.yml` (deployed) |
| Package scoping | Renderer: `src/packages/cnos.core/commands/install-wake/cn-install-wake`. CDS golden: `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml`. Admin golden: `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml`. Deployed: `.github/workflows/cnos-cds-dispatch.yml` and `.github/workflows/cnos-agent-admin.yml`. |
| Existing-binary disposition | Coexist (extend the renderer; interface unchanged) |
| Runtime dependencies | `actions/upload-artifact@v4` (GitHub-provided, no new binary deps); `git` + `jq` already present in runner |
| JSON/wire contract | `cn.wake-provider.v1` schema NOT extended; golden YAML gains new steps; no semantic changes to existing fields |
| Backward compat | AC6 invariant: no-op sweeps produce no committed traces; existing write-fence behavior preserved; `id:` on Claude step is additive; upload-artifact step is conditional; re-rendered goldens must match byte-for-byte with fresh render |

---

## 7. α Prompt

```
You are α (alpha) executing the implementation phase of CDD cycle/503 for cnos#503.

Issue: wake-trace: persist Claude execution logs as artifacts + cnos run trace
Mode: substantial (design + build)
Protocol: cds
Base SHA: 111c29701f41c65ce87c44d3304d4a9eabf61ac5
Cycle branch: cycle/503

Your contract is the γ scaffold at .cdd/unreleased/503/gamma-scaffold.md.
Read it fully before acting. Every pinned value in §6 (implementation contract)
is binding — do not improvise on any axis.

---

## What you must implement

### Overview

Extend the `cn install-wake` renderer to emit three new YAML blocks in both
wake flavors (cds-dispatch and agent-admin):

1. `id: claude` on the `anthropics/claude-code-action@v1` step
2. A conditional `actions/upload-artifact@v4` step
3. A conditional wake-trace write + commit step

Then re-render all downstream artifacts and add CI oracle steps.

---

### Step 1 — Read these files first (do not skip)

1. `src/packages/cnos.core/commands/install-wake/cn-install-wake` — the renderer (914 lines)
2. `src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml` — cds golden
3. `src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml` — admin golden
4. `.github/workflows/install-wake-golden.yml` — CI golden-check (has AC7/AC8 audit greps)
5. `.github/workflows/cnos-cds-dispatch.yml` — deployed cds-dispatch (must match golden)
6. `.github/workflows/cnos-agent-admin.yml` — deployed agent-admin (must match golden)

Pay special attention to:
- Line 758: `echo "      - uses: anthropics/claude-code-action@v1"` — this is where you insert
  `id: claude` BEFORE the `uses:` line. In YAML, `id:` must come before `uses:`.
- Lines 831-898: the activation-log write-fence block — your new steps must be inserted BEFORE
  this existing fence block (since the fence runs `if: always()` and must remain the final step).
- The AC7/AC8 audit greps in install-wake-golden.yml: your new renderer code must NOT contain
  `admin_only`, `disallowed_surfaces`, `defer_path`, `cell_execution`, `protocol:cds`, `cdr`,
  `cdw`, `dispatch:cell`, or `status:todo` as literal strings.

---

### Step 2 — Renderer changes (cn-install-wake)

#### 2a. Add `id: claude` before `uses: anthropics/claude-code-action@v1`

Find the section at approximately line 758 that emits the Claude action step:
```
  echo "      - uses: anthropics/claude-code-action@v1"
```

Insert immediately before it:
```
  echo "      - id: claude"
```

This is unconditional for both wake classes.

#### 2b. After the prompt body injection (after line 778 `sed ... prompt_path`), add the
upload-artifact step

Insert a new block unconditionally after the prompt body is inlined (between the prompt sed
and the write-fence block). For BOTH wake classes (dispatch and admin):

```sh
{
  echo ""
  echo "      - name: Upload Claude execution log"
  echo "        if: steps.claude.outputs.execution_file != ''"
  echo "        uses: actions/upload-artifact@v4"
  echo "        with:"
  echo "          name: claude-execution-\${{ github.run_id }}"
  echo "          path: \${{ steps.claude.outputs.execution_file }}"
  echo "          retention-days: 14"
} >> "$tmp_out"
```

#### 2c. Add the wake-trace write step

This step must be emitted BEFORE the write-fence block (the write-fence has `if: always()` and
is the final step; the trace step runs only on meaningful work and must precede it).

The predicate and trace location differ by wake class. Use the `role` variable (already extracted
from the manifest) to dispatch:

For `role = dispatch` (cds-dispatch):

```sh
if [ "$role" = "dispatch" ]; then
  {
    echo ""
    echo "      - name: Write wake trace"
    echo "        if: |"
    echo "          \${{ steps.claude.outputs.session_id != '' &&"
    echo "               hashFiles('.cdd/unreleased/**') != '' }}"
    echo "        shell: bash"
    echo "        run: |"
    echo "          set -euo pipefail"
    echo "          # Determine if meaningful work was done in this run."
    echo "          # AC6: no-op sweeps (empty selector scan) MUST NOT produce"
    echo "          # committed trace entries. Gate: check if .cdd/unreleased/"
    echo "          # gained new paths since CN_WAKE_BASE_SHA."
    echo "          base=\${CN_WAKE_BASE_SHA:-\${GITHUB_SHA:-}}"
    echo "          if [ -z \"\$base\" ]; then"
    echo "            echo 'wake-trace: no baseline SHA available; skipping trace commit.'"
    echo "            exit 0"
    echo "          fi"
    echo "          new_paths=\$(git log \"\$base..HEAD\" --name-only --pretty=format: -- '.cdd/unreleased/' 2>/dev/null | grep -v '^\$' || true)"
    echo "          if [ -z \"\$new_paths\" ]; then"
    echo "            echo 'wake-trace: no new .cdd/unreleased/ paths in this run — no-op sweep; skipping trace.'"
    echo "            exit 0"
    echo "          fi"
    echo "          # Extract the issue number from the first new path."
    echo "          issue_num=\$(echo \"\$new_paths\" | grep -oE '.cdd/unreleased/[0-9]+/' | head -1 | grep -oE '[0-9]+')"
    echo "          if [ -z \"\$issue_num\" ]; then"
    echo "            echo 'wake-trace: could not extract issue number from new paths; skipping trace.'"
    echo "            exit 0"
    echo "          fi"
    echo "          trace_path=\".cdd/unreleased/\${issue_num}/wake-trace.md\""
    echo "          cat > \"\$trace_path\" <<EOF"
    echo "# wake-trace — cnos#\${issue_num}"
    echo "run_id: \${{ github.run_id }}"
    echo "wake_class: cds-dispatch"
    echo "issue: \${issue_num}"
    echo "protocol: cds"
    echo "claude_session_id: \${{ steps.claude.outputs.session_id }}"
    echo "execution_artifact: claude-execution-\${{ github.run_id }}"
    echo "workflow: \${{ github.workflow }}"
    echo "conclusion: \${{ job.status }}"
    echo "last_known_phase: post-claude"
    echo "EOF"
    echo "          git add \"\$trace_path\""
    echo "          git -c user.name=\"\${bot_name}\" -c user.email=\"\${bot_name}\" \\"
    echo "            commit -m \"wake-trace: cnos#\${issue_num} run \${{ github.run_id }}\""
    echo "          git push"
  } >> "$tmp_out"
fi
```

IMPORTANT: The `bot_name` variable is already set in the renderer; use it for the git commit
identity. The `echo` lines containing the heredoc body are tricky — see Friction Note FN-1 below
for the correct approach.

For `role = admin` (agent-admin):

```sh
if [ "$role" = "admin" ]; then
  {
    echo ""
    echo "      - name: Write wake trace"
    echo "        if: steps.claude.outputs.session_id != ''"
    echo "        shell: bash"
    echo "        run: |"
    echo "          set -euo pipefail"
    echo "          # AC6 gate for admin wakes: skip trace if no commits were made"
    echo "          # (heartbeat no-op — HEAD equals baseline)."
    echo "          base=\${CN_WAKE_BASE_SHA:-\${GITHUB_SHA:-}}"
    echo "          if [ -n \"\$base\" ] && [ \"\$(git rev-parse HEAD)\" = \"\$base\" ]; then"
    echo "            echo 'wake-trace: HEAD unchanged from baseline — heartbeat no-op; skipping trace.'"
    echo "            exit 0"
    echo "          fi"
    echo "          trace_dir=\".cn-sigma/state/wake-traces\""
    echo "          mkdir -p \"\$trace_dir\""
    echo "          trace_path=\"\${trace_dir}/\${{ github.run_id }}.md\""
    echo "          cat > \"\$trace_path\" <<EOF"
    echo "# wake-trace — agent-admin"
    echo "run_id: \${{ github.run_id }}"
    echo "wake_class: agent-admin"
    echo "claude_session_id: \${{ steps.claude.outputs.session_id }}"
    echo "execution_artifact: claude-execution-\${{ github.run_id }}"
    echo "workflow: \${{ github.workflow }}"
    echo "conclusion: \${{ job.status }}"
    echo "last_known_phase: post-claude"
    echo "EOF"
    echo "          git add \"\$trace_path\""
    echo "          git -c user.name=\"\${bot_name}\" -c user.email=\"\${bot_name}\" \\"
    echo "            commit -m \"wake-trace: agent-admin run \${{ github.run_id }}\""
    echo "          git push"
  } >> "$tmp_out"
fi
```

Note: agent-admin already has no write-fence (write-fence is only for dispatch wakes with
`activation_log_writer:false`). The trace step is inserted between the prompt sed and the (absent)
fence — so for agent-admin it is simply the last real step.

**CRITICAL placement**: In the renderer, after the `sed ... prompt_path` line (line ~778), the
code immediately checks `if [ "$activation_log_writer" = "false" ]; then` (line 831) to emit the
write-fence. The new blocks (upload-artifact + wake-trace) must be inserted BETWEEN the prompt sed
and the write-fence conditional. The structure must be:

```
1. [existing] prompt sed injection
2. [NEW] upload-artifact step (unconditional; for all roles)
3. [NEW] wake-trace step (role-conditional: dispatch vs admin)
4. [existing] write-fence step (dispatch only; if activation_log_writer=false)
```

#### 2d. git identity in the trace commit step

The renderer emits `${bot_name}` into the shell block. Since the renderer's `echo` lines are
inside the renderer's sh code, `${bot_name}` must be interpolated at render time (it is a
renderer variable). Do NOT use `\$` for it inside the echo'd heredoc. However, GitHub Actions
expressions like `${{ github.run_id }}` must be emitted literally — escape the `$` in echo as
`\${{ ... }}`.

---

### Step 3 — Re-render goldens

After editing the renderer, re-render both goldens:

```sh
./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch
./src/packages/cnos.core/commands/install-wake/cn-install-wake agent-admin
```

Verify the goldens changed as expected:
```sh
git diff src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
git diff src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
```

The diffs must show exactly the three new blocks added: `id: claude`, the upload-artifact step,
and the wake-trace step. No other changes.

---

### Step 4 — Re-render deployed workflows

```sh
./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch \
  --out .github/workflows/cnos-cds-dispatch.yml
./src/packages/cnos.core/commands/install-wake/cn-install-wake agent-admin \
  --out .github/workflows/cnos-agent-admin.yml
```

Verify byte-identity with goldens:
```sh
sha256sum .github/workflows/cnos-cds-dispatch.yml \
  src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
sha256sum .github/workflows/cnos-agent-admin.yml \
  src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
```

Both pairs must match.

---

### Step 5 — Verify AC7/AC8 audit still passes

```sh
n_admin=$(grep -ciE 'admin.only|disallowed_surfaces|defer.path|cell_execution' \
  src/packages/cnos.core/commands/install-wake/cn-install-wake || true)
n_dispatch=$(grep -ciE 'protocol:cds|cdr|cdw|dispatch:cell|status:todo' \
  src/packages/cnos.core/commands/install-wake/cn-install-wake || true)
echo "admin leaks: $n_admin  dispatch leaks: $n_dispatch"
```

Both must be 0. If your new renderer code introduced any of those strings (even in comments),
the CI will fail.

---

### Step 6 — Extend install-wake-golden.yml with AC1/AC2/AC3/AC6 oracle steps

Add the following steps to `.github/workflows/install-wake-golden.yml` after the existing
"Verify substrate structural shape" steps:

**AC1 oracle** — `id: claude` present in both goldens:
```yaml
      - name: AC1 (cnos#503) — Claude step has id
        run: |
          set -eu
          for f in \
            src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml \
            src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml; do
            if ! grep -qE '^\s+id: claude$' "$f"; then
              echo "::error::AC1 cnos#503: missing 'id: claude' on Claude action step in $f"
              exit 1
            fi
          done
          echo "AC1: 'id: claude' present in both goldens."
```

**AC2 oracle** — upload-artifact step with correct gate and retention:
```yaml
      - name: AC2 (cnos#503) — upload-artifact step present with retention-days 14
        run: |
          set -eu
          for f in \
            src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml \
            src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml; do
            if ! grep -qF 'upload-artifact@v4' "$f"; then
              echo "::error::AC2 cnos#503: missing upload-artifact@v4 step in $f"
              exit 1
            fi
            if ! grep -qF 'steps.claude.outputs.execution_file' "$f"; then
              echo "::error::AC2 cnos#503: upload-artifact step not gated on steps.claude.outputs.execution_file in $f"
              exit 1
            fi
            if ! grep -qE 'retention-days: 14' "$f"; then
              echo "::error::AC2 cnos#503: retention-days: 14 not found in $f"
              exit 1
            fi
          done
          echo "AC2: upload-artifact step present with correct gate and retention in both goldens."
```

**AC3 oracle** — trace write step with required fields:
```yaml
      - name: AC3 (cnos#503) — wake-trace write step with required fields
        run: |
          set -eu
          # Required fields per AC3: run_id, wake_class, claude_session_id,
          # execution_artifact, workflow, conclusion, last_known_phase.
          # Cell-dispatch additionally: issue, protocol.
          cds=src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
          admin=src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
          base_fields="run_id wake_class claude_session_id execution_artifact workflow conclusion last_known_phase"
          for field in $base_fields; do
            if ! grep -qF "$field" "$cds"; then
              echo "::error::AC3 cnos#503: missing required trace field '$field' in cds-dispatch golden"
              exit 1
            fi
            if ! grep -qF "$field" "$admin"; then
              echo "::error::AC3 cnos#503: missing required trace field '$field' in agent-admin golden"
              exit 1
            fi
          done
          # Cell-dispatch-only fields.
          for field in issue protocol; do
            if ! grep -qF "$field" "$cds"; then
              echo "::error::AC3 cnos#503: missing cell-dispatch-only trace field '$field' in cds-dispatch golden"
              exit 1
            fi
          done
          echo "AC3: all required trace fields present in both goldens."
```

**AC6 oracle** — no-op guard present in cds-dispatch golden:
```yaml
      - name: AC6 (cnos#503) — no-op guard present in cds-dispatch golden
        run: |
          set -eu
          f=src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
          # The no-op guard reads git log $CN_WAKE_BASE_SHA..HEAD to check
          # if .cdd/unreleased/ gained new paths. Assert the key predicate
          # is present in the rendered trace step.
          if ! grep -qF '.cdd/unreleased/' "$f"; then
            echo "::error::AC6 cnos#503: cds-dispatch golden does not reference .cdd/unreleased/ in no-op guard"
            exit 1
          fi
          if ! grep -qF 'CN_WAKE_BASE_SHA' "$f"; then
            echo "::error::AC6 cnos#503: cds-dispatch golden does not reference CN_WAKE_BASE_SHA in no-op guard"
            exit 1
          fi
          echo "AC6: no-op guard referencing .cdd/unreleased/ and CN_WAKE_BASE_SHA present in cds-dispatch golden."
```

Also add `.github/workflows/cnos-agent-admin.yml` to the `install-wake-golden.yml` path trigger
(it is not currently listed) so the golden-check CI fires on agent-admin deployed-workflow changes.

---

### Step 7 — AC5 documentation note

Create `.cdd/unreleased/503/wake-trace-cnos504-note.md` with content:

```markdown
# cnos#504 cross-reference note

The wake-trace surface established by cnos#503 is the intended diagnostic
starting point for stale-claim recovery (cnos#504).

Trace location for cell-dispatch wakes: `.cdd/unreleased/{issue}/wake-trace.md`

The `claude_session_id` and `execution_artifact` fields in the trace file
point to the Claude session and the uploaded execution artifact (workflow
artifact, 14-day retention). These are the primary diagnostic inputs for
cnos#504 stale-claim recovery.

This note satisfies AC5 of cnos#503 (documentation-only AC).
```

---

### Step 8 — Author self-coherence.md

Author `.cdd/unreleased/503/self-coherence.md` covering:
- §R0: initial pass — what you built, what you verified, what edge cases exist
- §Debt: any open questions or deferred work
- §Friction: any traps you encountered

---

### Step 9 — Commit and push

```sh
git add \
  src/packages/cnos.core/commands/install-wake/cn-install-wake \
  src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml \
  src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml \
  .github/workflows/cnos-cds-dispatch.yml \
  .github/workflows/cnos-agent-admin.yml \
  .github/workflows/install-wake-golden.yml \
  .cdd/unreleased/503/self-coherence.md \
  .cdd/unreleased/503/wake-trace-cnos504-note.md

git commit -m "cnos#503: wake-trace — id:claude, upload-artifact, trace step

AC1: id: claude on anthropics/claude-code-action@v1 step
AC2: conditional upload-artifact@v4 with retention-days:14
AC3: wake-trace write step (wake-class-appropriate location)
AC4: raw logs not committed; retention enforced
AC5: cnos#504 cross-reference note
AC6: no-op guard via CN_WAKE_BASE_SHA..HEAD check on .cdd/unreleased/"

git push -u origin cycle/503
```

---

### Critical constraints (re-stated)

1. The renderer MUST NOT contain the literal strings `admin_only`, `disallowed_surfaces`,
   `defer_path`, `cell_execution`, `protocol:cds`, `dispatch:cell`, `status:todo` — the AC7/AC8
   CI audit will fail if they appear anywhere in the renderer source, including comments.

2. `id:` must come BEFORE `uses:` in the emitted YAML. Wrong order will cause `actions/checkout`
   to fail or produce unparseable YAML.

3. The `if:` condition on the upload-artifact step must use the expression syntax exactly:
   `steps.claude.outputs.execution_file != ''` (no `${{ }}` wrapper needed in the `if:` field
   when using the expression form). Verify in the golden YAML that this parses correctly.

4. The wake-trace commit uses `git -c user.name=... -c user.email=... commit` — do NOT use
   `git config --global` (it would persist across steps and is an unintended side effect).

5. The `git push` in the trace step must succeed on the cycle branch. The renderer emits the
   branch name statically (it doesn't know the branch at render time). The `SIGMA_WORKFLOW_PAT`
   token (already present) provides the auth. The `actions/checkout@v4` step already checked out
   the branch; `git push` with no args pushes to the tracking remote.

6. All new `if:` conditions on steps must use GitHub Actions expression syntax correctly. For a
   step-level `if:`, the `${{ }}` wrapper is optional when the expression is the only content.
   Use `if: steps.claude.outputs.execution_file != ''` (no outer `${{ }}`).

7. After all changes, the second render (idempotence check) must produce byte-identical output.
   Run the idempotence check manually before committing:
   ```sh
   sha1=$(sha256sum src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml)
   ./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch
   sha2=$(sha256sum src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml)
   [ "$sha1" = "$sha2" ] && echo "idempotent" || echo "NON-IDEMPOTENT — fix renderer"
   ```
```

---

## 8. β Prompt

```
You are β (beta) reviewing the implementation of CDD cycle/503 for cnos#503.

Issue: wake-trace: persist Claude execution logs as artifacts + cnos run trace
Cycle branch: cycle/503
Mode: substantial (design + build)

Read the γ scaffold at .cdd/unreleased/503/gamma-scaffold.md before reviewing.
Your job is to produce .cdd/unreleased/503/beta-review.md with a verdict:
CONVERGE or ITERATE (with specific findings if ITERATE).

---

## What to verify

### AC1 — id: claude on the Claude step

```sh
grep -nE '^\s+id: claude$' \
  src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml \
  src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
```

Expected: at least one match per file, immediately before the `uses:
anthropics/claude-code-action@v1` line.

Also verify the `id: claude` appears BEFORE `uses:` in the YAML (ordering matters for YAML
parsers).

### AC2 — upload-artifact step

```sh
grep -nF 'upload-artifact@v4' \
  src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml \
  src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
grep -nF 'steps.claude.outputs.execution_file' \
  src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml \
  src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
grep -nF 'retention-days: 14' \
  src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml \
  src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
```

All three must return matches in both files.

### AC3 — wake-trace write step with required fields

Required fields for all wakes: `run_id`, `wake_class`, `claude_session_id`, `execution_artifact`,
`workflow`, `conclusion`, `last_known_phase`.

Additionally for cds-dispatch: `issue`, `protocol`.

```sh
for field in run_id wake_class claude_session_id execution_artifact workflow conclusion last_known_phase issue protocol; do
  echo "=== $field ==="
  grep -nF "$field" src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
done
for field in run_id wake_class claude_session_id execution_artifact workflow conclusion last_known_phase; do
  echo "=== $field ==="
  grep -nF "$field" src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
done
```

All fields must be present in the respective golden.

Verify the trace write step writes to the correct location:
- cds-dispatch: `.cdd/unreleased/${issue_num}/wake-trace.md` (uses variable extracted from git log)
- agent-admin: `.cn-sigma/state/wake-traces/${{ github.run_id }}.md` or similar admin-owned path
  (NOT `.cdd/unreleased/`)

### AC4 — No raw JSON committed; 14-day retention

Verify the upload-artifact step uses `path: ${{ steps.claude.outputs.execution_file }}` — it
uploads the file path, not the file content embedded in YAML. No `cat $execution_file` appears in
any committed step.

Verify `retention-days: 14` is present (already checked in AC2).

### AC5 — Documentation note for cnos#504

```sh
ls .cdd/unreleased/503/
cat .cdd/unreleased/503/wake-trace-cnos504-note.md
```

Verify the file exists and mentions cnos#504, the trace location
`.cdd/unreleased/{issue}/wake-trace.md`, and the `claude_session_id` / `execution_artifact`
fields.

### AC6 — No-op guard

In cds-dispatch golden, verify:
1. The wake-trace step references `CN_WAKE_BASE_SHA`
2. The step checks `git log "$base..HEAD" -- '.cdd/unreleased/'` (or equivalent)
3. The step exits 0 (skips) when the result is empty
4. The step does NOT commit when the result is empty

```sh
grep -n 'CN_WAKE_BASE_SHA' \
  src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
grep -n '.cdd/unreleased' \
  src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
```

For agent-admin, verify the no-op guard checks `HEAD == baseline` (heartbeat case) and skips the
trace in that case.

### Golden-check CI consistency

Verify the install-wake-golden.yml was updated with four new steps (AC1, AC2, AC3, AC6 oracles).

```sh
grep -n 'AC1\|AC2\|AC3\|AC6' .github/workflows/install-wake-golden.yml | grep 'cnos#503'
```

Expected: at least four matches (one per AC).

### AC7/AC8 — Renderer does not contain role-decision strings

```sh
grep -ciE 'admin.only|disallowed_surfaces|defer.path|cell_execution' \
  src/packages/cnos.core/commands/install-wake/cn-install-wake
grep -ciE 'protocol:cds|cdr|cdw|dispatch:cell|status:todo' \
  src/packages/cnos.core/commands/install-wake/cn-install-wake
```

Both must return 0.

### Deployed workflow byte-identity with goldens

```sh
sha256sum .github/workflows/cnos-cds-dispatch.yml \
  src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml
sha256sum .github/workflows/cnos-agent-admin.yml \
  src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
```

Each pair must be byte-identical.

### Idempotence of renderer

```sh
sha1=$(sha256sum src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml | cut -d' ' -f1)
./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch >/dev/null
sha2=$(sha256sum src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml | cut -d' ' -f1)
[ "$sha1" = "$sha2" ] && echo "IDEMPOTENT" || echo "NON-IDEMPOTENT — renderer is not deterministic"
```

Must output "IDEMPOTENT".

### Self-coherence

Read `.cdd/unreleased/503/self-coherence.md`. Verify it names what was built, what was verified,
and what debt or open questions remain.

---

## Verdict criteria

CONVERGE if all of the following hold:
- AC1 through AC6 oracle checks above all pass
- AC7/AC8 renderer audit passes (0 role-decision string leaks)
- Deployed workflows byte-identical with goldens
- Renderer is idempotent
- self-coherence.md is present and coherent
- No raw JSON committed; no admin-surface writes by the cds-dispatch trace step
- The cds-dispatch wake-trace step does NOT write to `.cn-sigma/` (that surface is admin-only)

ITERATE if any check fails. Name the specific failure(s) and the exact fix required. Do not
name stylistic issues unless they affect correctness or AC compliance.
```

---

## 9. Scope Guardrails

**Explicit non-goals** (out of scope for cycle/503):

- Modifying `anthropics/claude-code-action@v1` itself — the action is consumed, not modified
- Replacing claude-code-action with a different action
- Adding a UI or console for viewing traces
- Raw log storage beyond the 14-day artifact (no long-term log archival)
- Extending the `cn.wake-provider.v1` schema — trace machinery is renderer authority only
- Modifying the write-fence logic (the write-fence is preserved unchanged)
- Changing the concurrency model or cron schedule
- Adding new labels or changing label semantics

**Boundary rules**:

- The trace write step for cds-dispatch MUST write to `.cdd/unreleased/{N}/wake-trace.md` on the
  cycle branch — NOT to `main`, NOT to `.cn-sigma/`, NOT to any admin-owned surface.
- The trace write step for agent-admin MUST write to an admin-owned surface (e.g.
  `.cn-sigma/state/wake-traces/`) — NOT to `.cdd/unreleased/`.
- The write-fence (dispatch-wake-only) remains the final step. The new trace step is inserted
  before the write-fence.
- The renderer MUST NOT hand-patch YAML. Every downstream artifact must be derived via re-render.
- Stage 2 (live smoke) is post-merge; α does not attempt live smoke in this cycle.

---

## 10. Friction Notes

**FN-1 — Shell heredoc inside echo chains**

The trace write step emits a heredoc (`cat > "$trace_path" <<EOF ... EOF`) inside shell `echo`
lines in the renderer. This creates quoting fragility: the renderer's `echo` lines must:
(a) not expand `${{ ... }}` expressions at render time (GHA expressions must be emitted literally
    into the YAML)
(b) not expand `${issue_num}` and similar shell variables at render time (they must run at job
    time in the runner)
(c) DO expand `${bot_name}` at render time (it's a renderer variable bound to the agent identity)

The safest approach: use a separate shell function or `printf '%s\n'` for the heredoc body lines
rather than embedding them in `echo` chains. Alternatively, write the trace step shell block
into a variable and emit it line-by-line. See how the write-fence block handles this problem
(lines 831-898 of cn-install-wake) — it uses many `echo "          ..."` lines with careful
double/single-quote nesting to control what expands at render time vs run time.

**FN-2 — `if:` syntax for step conditionals**

GitHub Actions step-level `if:` fields accept expressions with or without `${{ }}`. When the
content is a bare expression (no string interpolation), omit the outer braces:

```yaml
if: steps.claude.outputs.execution_file != ''
```

NOT:
```yaml
if: ${{ steps.claude.outputs.execution_file != '' }}
```

Both are valid, but the bare form is what the existing write-fence step uses (see golden line 253:
`if: always()`). The renderer's `echo` lines should emit the form without outer braces for
consistency.

**FN-3 — YAML ordering: `id:` before `uses:`**

In a GitHub Actions step, `id:` must appear before `uses:` (or `run:`). The renderer currently
emits `- uses: anthropics/claude-code-action@v1` without a preceding `id:`. α must insert
`- id: claude` as a separate `echo` line BEFORE the `echo "      - uses: anthropics/..."` line.
Do not combine them on one line.

**FN-4 — Agent-admin golden lacks a write-fence**

The agent-admin wake does NOT have an activation-log write fence (write-fence is only for dispatch
wakes with `activation_log_writer:false`). The agent-admin manifest has `activation_log_writer`
defaulting to `true` (writer). So for agent-admin, the trace step is the last step in the
rendered YAML — there is no write-fence to insert before. This is intentional; the renderer's
conditional (`if [ "$activation_log_writer" = "false" ]; then`) governs the fence, not the trace.

**FN-5 — AC6 predicate must not use `hashFiles` in the `if:` expression**

`hashFiles()` in a step-level `if:` expression computes a hash of file contents at job evaluation
time — it does not check whether files were ADDED in this run's commits. The correct predicate is
a shell check inside the step's `run:` block: `git log "$CN_WAKE_BASE_SHA..HEAD" --name-only --
'.cdd/unreleased/'`. The γ scaffold AC oracle prompt uses `hashFiles` in a draft expression; α
should NOT use `hashFiles` as the gate. The gate must be inside the `run:` block, not in the
step-level `if:`.

**FN-6 — install-wake-golden.yml path trigger for agent-admin deployed workflow**

The current `install-wake-golden.yml` does not trigger on changes to
`.github/workflows/cnos-agent-admin.yml` (only `cnos-cds-dispatch.yml` is listed). α must add
`.github/workflows/cnos-agent-admin.yml` to the `paths:` list in that workflow's `on:` block so
the CI fires when the deployed agent-admin workflow changes.

**FN-7 — git push in the trace step needs the right remote and branch**

The trace step emits `git push` (bare). This works only if the runner's checkout has a tracking
remote set. `actions/checkout@v4` with a PAT sets the remote automatically. However, for the
cds-dispatch case, the Claude action runs on a cycle branch (e.g. `cycle/503`); the runner
checked out `main` initially. The trace commit goes to the cycle branch, which α's work already
pushed. The `git push` in the trace step will push to the branch the runner is currently on.

If the runner is on `main` (no cycle work happened because the scan was a no-op), the AC6 guard
exits early — no commit, no push. If the runner is on a cycle branch, the push goes to that
branch. Verify in the golden that the `git push` is unconditional inside the "meaningful work"
branch and is reachable only after the AC6 gate passes.

**FN-8 — Main SHA drift**

At γ dispatch time, `origin/main` was at SHA `6175a73e6a7db1f748bc8b49302855daa6e51fa4` (δ noted
this; the cycle branch was created from the pinned SHA `111c29701f41c65ce87c44d3304d4a9eabf61ac5`
per dispatch instructions). α works on `cycle/503` which diverges from `111c297`. When α opens
the PR, the merge base will be `111c297` not `6175a73`. This is expected and acceptable; the PR
targets `main` and CI will run against the merge-commit state.
