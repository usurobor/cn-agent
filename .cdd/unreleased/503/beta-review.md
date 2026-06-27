# β Review — cycle/503 cnos#503 §R0

## Verdict: ITERATE

## Findings

### AC1 — PASS

Both files have `- id: claude` as a list-item key immediately before `uses: anthropics/claude-code-action@v1`.

```
cnos-agent-admin.golden.yml:42:      - id: claude
cnos-agent-admin.golden.yml:43:        uses: anthropics/claude-code-action@v1

cnos-cds-dispatch.golden.yml:56:      - id: claude
cnos-cds-dispatch.golden.yml:57:        uses: anthropics/claude-code-action@v1
```

`id: claude` line numbers (42, 56) are lower than `uses: anthropics/...` line numbers (43, 57) in both files. AC1 PASS.

### AC2 — PASS

All three required elements present in both files:

```
# upload-artifact@v4
cnos-agent-admin.golden.yml:188:        uses: actions/upload-artifact@v4
cnos-cds-dispatch.golden.yml:248:        uses: actions/upload-artifact@v4

# steps.claude.outputs.execution_file gate
cnos-agent-admin.golden.yml:187:        if: steps.claude.outputs.execution_file != ''
cnos-cds-dispatch.golden.yml:247:        if: steps.claude.outputs.execution_file != ''

# path: (file reference, not content)
cnos-agent-admin.golden.yml:191:          path: ${{ steps.claude.outputs.execution_file }}
cnos-cds-dispatch.golden.yml:251:          path: ${{ steps.claude.outputs.execution_file }}

# retention-days: 14
cnos-agent-admin.golden.yml:192:          retention-days: 14
cnos-cds-dispatch.golden.yml:252:          retention-days: 14
```

AC2 PASS.

### AC3 — PASS

All required common fields present in both goldens:

cds-dispatch (lines 279–287): `run_id`, `wake_class`, `claude_session_id`, `execution_artifact`, `workflow`, `conclusion`, `last_known_phase`.

cds-dispatch-only fields: `issue` (line 281), `protocol` (line 282) — both present.

agent-admin (lines 208–214): all seven common fields present.

Trace locations:
- cds-dispatch writes to `.cdd/unreleased/${issue_num}/wake-trace.md` (line 276, 289 of golden).
- agent-admin writes to `.cn-sigma/state/wake-traces/${{ github.run_id }}.md` (lines 205, 215–216 of golden).

AC3 PASS.

### AC4 — PASS

No `cat.*execution` matches in either golden. The upload step uses `path: ${{ steps.claude.outputs.execution_file }}` (file path only — no content embedding). The trace write step emits structured metadata fields only. No raw JSON committed. Retention-days: 14 enforced via `retention-days: 14`.

AC4 PASS.

### AC5 — PASS

`.cdd/unreleased/503/wake-trace-cnos504-note.md` exists and contains:
- Reference to `cnos#504` ✓
- Trace location `.cdd/unreleased/{issue}/wake-trace.md` ✓
- Fields `claude_session_id` and `execution_artifact` ✓

AC5 PASS.

### AC6 — PASS

cds-dispatch golden no-op guard (lines 259–274):
1. `CN_WAKE_BASE_SHA` referenced at line 261 (`base=${CN_WAKE_BASE_SHA:-...}`).
2. `git log "$base..HEAD" --name-only --pretty=format: -- '.cdd/unreleased/'` at line 266.
3. Early `exit 0` when result is empty (line 269) and when no baseline SHA (line 264) and when issue number cannot be extracted (line 274).
4. No commit/push runs when early exit fires (the git commit/push only runs after the guard block completes).

agent-admin HEAD==baseline guard (lines 200–203):
```
grep -n 'git rev-parse HEAD...' -> line 201: [ "$(git rev-parse HEAD)" = "$base" ]
```
Present. AC6 PASS.

### AC7/AC8 — PASS

Renderer audit:
```
grep -ciE 'admin.only|disallowed_surfaces|defer.path|cell_execution' cn-install-wake → 4
grep -ciE 'protocol:cds|cdr|cdw|dispatch:cell|status:todo' cn-install-wake → 0
```

The 4 `admin.only` hits are pre-existing from cycle/496 error messages (confirmed in self-coherence.md). The diff `git diff origin/main...origin/cycle/503 -- cn-install-wake` introduces none of the disallowed strings. `protocol_val` is derived via `jq -r '.protocol // "unknown"' "$manifest_path"` (renderer line 810) — not hardcoded. AC7/AC8 PASS for cycle/503's diff.

### Byte-identity — PASS

```
e7e71c3054b73bcadcafadfba728a735b4859fdcedd5508c8a1eadd5784b46b7
  .github/workflows/cnos-cds-dispatch.yml
  src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml

5312623ab41f22d6a46a47e8a4c2261b829b3840ccb730b0f68b446810fe84bc
  .github/workflows/cnos-agent-admin.yml
  src/packages/cnos.core/orchestrators/agent-admin/cnos-agent-admin.golden.yml
```

Both pairs byte-identical. PASS.

### Idempotence — PASS

```
cn-install-wake: cds-dispatch → cnos-cds-dispatch.golden.yml (unchanged)
CDS IDEMPOTENT

cn-install-wake: agent-admin → cnos-agent-admin.golden.yml (unchanged)
ADMIN IDEMPOTENT
```

Both wakes idempotent. PASS.

### Correctness: YAML structure — PASS

In both goldens, `- id: claude` is the first key of the step's list mapping (at 6-space indent for the dash + 8-space for subsequent keys), and `uses: anthropics/claude-code-action@v1` is a sibling key at the same indent level immediately after. This is the correct YAML mapping structure — not two separate list items.

Verified by inspecting the raw bytes (`cat -A`) at lines 54–60 (cds-dispatch) and 40–46 (agent-admin). PASS.

### Correctness: Commit message quoting — FAIL

In the cds-dispatch golden (line 291):

```yaml
            commit -m 'wake-trace: ${issue_num} run ${{ github.run_id }}'
```

The commit message is wrapped in **single quotes**. At runtime, bash single-quote inhibits variable expansion. Therefore `${issue_num}` will NOT expand — the commit message will be the literal string `wake-trace: ${issue_num} run <run_id>` rather than the intended `wake-trace: 503 run <run_id>`.

Note: `${{ github.run_id }}` is a GitHub Actions expression, not a bash variable — GH Actions substitutes it before the shell runs. But `${issue_num}` is a pure bash runtime variable and single-quote quoting suppresses it.

The renderer source (line 853) emits:
```bash
echo "            commit -m 'wake-trace: \${issue_num} run \${{ github.run_id }}'"
```
The `\${issue_num}` in the renderer's double-quoted string produces `${issue_num}` in the YAML, and the surrounding single quotes in the YAML prevent bash from expanding it at runtime. FAIL.

The agent-admin golden (line 218) is unaffected: `commit -m 'wake-trace: agent-admin run ${{ github.run_id }}'` — no bash runtime variable; `${{ github.run_id }}` is a GH Actions expression handled pre-shell. That line is correct.

### Correctness: Write-fence ordering — PASS

Step name line numbers in cds-dispatch golden:
- Line 246: `name: Upload Claude execution log`
- Line 254: `name: Write wake trace`
- Line 302: `if: always()` (Write fence step)

Ordering: Upload (246) < Write wake trace (254) < Write fence (302). PASS.

### Correctness: Admin surface isolation — PASS

```
grep -n 'cn-sigma' cnos-cds-dispatch.golden.yml | grep -v 'SIGMA_WORKFLOW_PAT|cn-sigma.cnos'
```

Only hit: line 190 — a prose prohibition note in the Claude prompt body (`- `.cn-sigma/logs/` and other `.cn-sigma/` surfaces...`), not a path write. No `.cn-sigma/` filesystem writes by the cds-dispatch trace step. PASS.

---

## If ITERATE: required fixes

**Fix 1 — cds-dispatch commit message quoting (cnos-cds-dispatch.golden.yml line 291 + renderer source line 853)**

The `${issue_num}` bash runtime variable is wrapped in single quotes and will not expand at runtime.

Fix in the renderer source (`src/packages/cnos.core/commands/install-wake/cn-install-wake`, line 853):

Change:
```bash
echo "            commit -m 'wake-trace: \${issue_num} run \${{ github.run_id }}'"
```

To (using double quotes with escaped inner double-quotes):
```bash
echo "            commit -m \"wake-trace: \${issue_num} run \${{ github.run_id }}\""
```

This produces in the YAML:
```yaml
            commit -m "wake-trace: ${issue_num} run ${{ github.run_id }}"
```

At runtime: bash expands `${issue_num}` (runtime variable); GH Actions has already substituted `${{ github.run_id }}` before the shell runs. Both expand correctly.

After fixing the renderer, re-render the golden and re-deploy:
```bash
./src/packages/cnos.core/commands/install-wake/cn-install-wake cds-dispatch
cp src/packages/cnos.cds/orchestrators/cds-dispatch/cnos-cds-dispatch.golden.yml \
   .github/workflows/cnos-cds-dispatch.yml
```

Then verify byte-identity and idempotence again.
