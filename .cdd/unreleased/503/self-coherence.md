# Self-Coherence — cycle/503 cnos#503

## §R0

### What was built

Three renderer changes to `src/packages/cnos.core/commands/install-wake/cn-install-wake`:

1. **AC1 — id: claude on the Claude step**: Added `echo "      - id: claude"` immediately before `echo "        uses: anthropics/claude-code-action@v1"`. This is unconditional for both wake classes (dispatch and admin). In the rendered YAML the step now reads `- id: claude` followed by `uses:`.

2. **AC2 — Upload-artifact step**: Added an `actions/upload-artifact@v4` step gated on `steps.claude.outputs.execution_file != ''`, with `retention-days: 14` and artifact name `claude-execution-${{ github.run_id }}`. This step is unconditional for both wake classes and is inserted after the Claude step (after the prompt body injection).

3. **AC3/AC6 — Wake-trace write step**: Added a role-conditional trace write step after the upload-artifact step. For dispatch wakes: checks if `.cdd/unreleased/` gained new paths since `CN_WAKE_BASE_SHA` (AC6 no-op guard), extracts the issue number, and writes `.cdd/unreleased/{issue}/wake-trace.md` with required fields. For admin wakes: checks if HEAD changed from baseline (AC6 heartbeat guard), and writes `.cn-sigma/state/wake-traces/${{ github.run_id }}.md`.

Both goldens re-rendered and deployed workflows updated to be byte-identical with goldens.

CI oracle steps added to `install-wake-golden.yml`:
- AC1 oracle: `id: claude` in both goldens
- AC2 oracle: `upload-artifact@v4` with gate and `retention-days: 14`
- AC3 oracle: all required trace fields in both goldens
- AC6 oracle: `CN_WAKE_BASE_SHA` and `.cdd/unreleased/` in cds-dispatch golden
- Byte-identity check for agent-admin deployed workflow vs golden (mirrors existing cds-dispatch check)
- Path trigger extended to include `.github/workflows/cnos-agent-admin.yml`

AC5 documentation note: `.cdd/unreleased/503/wake-trace-cnos504-note.md` naming cnos#504 as the stale-claim recovery consumer of the trace surface.

### AC verification

- AC1: VERIFIED — `id: claude` present in both goldens at correct position (before `uses:`). CI oracle uses `grep -qF 'id: claude'`. Manual oracle check passes.
- AC2: VERIFIED — `upload-artifact@v4` step with `if: steps.claude.outputs.execution_file != ''` and `retention-days: 14` present in both goldens. Manual oracle check passes.
- AC3: VERIFIED — all required fields (run_id, wake_class, claude_session_id, execution_artifact, workflow, conclusion, last_known_phase) present in both goldens. Cell-dispatch-only fields (issue, protocol) present in cds-dispatch golden. Manual oracle check passes.
- AC4: VERIFIED — no raw execution log content committed; upload step uses `path: ${{ steps.claude.outputs.execution_file }}` (file path only); retention-days: 14 enforced. Trace write step writes only structured metadata fields, not raw log content.
- AC5: VERIFIED — `.cdd/unreleased/503/wake-trace-cnos504-note.md` committed, names cnos#504 and the trace surface.
- AC6: VERIFIED — no-op guard present in cds-dispatch trace step: checks `git log "$base..HEAD" -- '.cdd/unreleased/'` inside `run:` block; exits 0 on empty result. Admin wake guard checks HEAD == baseline. Both guards are inside `run:` blocks (not step-level `if:`). Manual oracle check passes.

### Review-ready signal

READY FOR REVIEW

### Debt

- AC7/AC8 audit: 4 pre-existing `admin_only` matches in the renderer (from cycle/496 error messages). These were present before cycle/503 and are pre-existing failures in CI on main. My changes added zero new matches — verified by `git diff` inspection.
- Stage 2 (live smoke): out of scope per gamma-scaffold — post-merge only. Not closing cnos#503 until Stage 2 evidence lands.
- The `protocol_val` variable uses `jq -r '.protocol // "unknown"' "$manifest_path"` at render time, so the rendered YAML contains the literal `cds` (from the manifest). The renderer source carries no `protocol:cds` literal — it comes from the manifest via jq — satisfying AC7.

### Friction notes

- FN-dash-echo: Dash's `echo` in double-quoted strings interprets `\n` as a literal newline (not the backslash-n sequence). Changed `printf '%s\n'` approach to a `{ echo '...'; echo '...'; } > file` block, which avoids the `\n` quoting problem entirely.
- FN-grep-pattern: `grep -qF '- id: claude'` fails on GNU grep when the pattern starts with `-` (misinterpreted as flag). Changed AC1 oracle to use `grep -qF 'id: claude'` (without leading dash) which is still unambiguous since no other step has that string.

## §R1

### What was fixed

β R0 finding: commit message in dispatch wake-trace step used single quotes (`'...'`), which inhibits bash variable expansion of `${issue_num}` at runtime. Fixed by changing the renderer echo line to use escaped double quotes so `${issue_num}` expands to the actual issue number at job runtime.

### AC verification

- Commit message quoting: FIXED — cds-dispatch golden now uses `"wake-trace: ${issue_num} run ${{ github.run_id }}"` (double-quoted).
- All other ACs unchanged from R0.
- Byte-identity: VERIFIED — deployed workflows match goldens.
- Idempotence: VERIFIED.

### Review-ready signal

READY FOR REVIEW
