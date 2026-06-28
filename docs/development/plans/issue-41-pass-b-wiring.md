# Design: Wire Two-Pass Execution in cn_runtime (Issue #41)

**Status:** Complete
**Date:** 2026-03-20
**Issue:** #41 — pass B not invoked
**Mode:** MCA — spec correct, code incomplete
**Branch:** claude/cdd-executable-skill-PfdYZ

---

## Coherence Contract

**Gap:** `cn_runtime.ml:finalize` computes `needs_pass_b=true` after Pass A but
never invokes Pass B. Observe ops execute and artifacts are written to disk,
but the LLM is never re-called with evidence. The sense→act loop is broken.

**Scope:** Runtime integration only. All building blocks exist and are tested:
- `Cn_orchestrator.run_pass_b` — fully implemented
- `Cn_orchestrator.repack_for_pass_b` — fully implemented
- `Cn_executor.write_receipts` — appends correctly
- `Cn_orchestrator.gate_coordination` — fully implemented

**Expected triadic effect:**
- β (RELATION): docs/runtime alignment restored — AGENT-RUNTIME spec matches code behavior
- α (PATTERN): internal consistency — observe ops now produce visible evidence
- γ (EXIT): evolution path — two-pass is the foundation for future observe+effect workflows

**Failure if skipped:** Observe ops are dead code. Agent can request evidence but
never receives it, breaking the core sensing→picture→action loop.

---

## Design

### Current Flow (cn_runtime.ml:finalize, lines 446-622)

```
1. archive_raw
2. parse output
3. execute ops:
   a. coordination ops → Cn_agent.execute_op
   b. parser denial receipts → Cn_orchestrator.write_denial_receipts
   c. typed ops → Cn_orchestrator.run_pass_a
      └─ _pass_a_result.needs_pass_b logged but IGNORED
4. mark_ops_done
5. project to Telegram
6. append conversation
7. cleanup
```

### Target Flow

```
1. archive_raw
2. parse output (Pass A)
3. execute ops:
   a. coordination ops (Pass-A-safe only when two-pass)
   b. parser denial receipts
   c. typed ops → run_pass_a
   d. IF needs_pass_b:
      i.   repack context with artifacts/receipts
      ii.  archive Pass A IO pair
      iii. call LLM (second API call)
      iv.  archive Pass B output
      v.   parse Pass B output
      vi.  execute Pass B typed ops → run_pass_b
      vii. execute Pass B coordination ops (gated on effect success)
4. mark_ops_done
5. project to Telegram (uses Pass B output when two-pass)
6. append conversation (uses Pass B output when two-pass)
7. cleanup
```

### Key Design Decisions

1. **Coordination ops in two-pass:** During Pass A, only pass-A-safe ops execute
   (ack, surface, reply with idempotent projection). Terminal ops (done, fail,
   send, delegate, defer, delete) are deferred to after Pass B completes, where
   they are gated on effect success.

2. **Pass B repacking:** Uses `Cn_orchestrator.repack_for_pass_b` which injects
   receipts summary + artifact excerpts. Skills are NOT re-scored (fixed at
   Pack time per spec).

3. **Pass B context structure:** The Pass B LLM call uses the SAME system blocks
   as Pass A (skills fixed), but appends the repack content as an additional
   assistant+user turn pair in the messages array.

4. **Projection uses final output:** When two-pass occurs, Telegram projection
   and conversation history use Pass B's output, not Pass A's.

5. **Recovery:** Pass B output is written to `state/output.md` (overwriting
   Pass A). If crash occurs mid-Pass-B, State 1 recovery will re-finalize with
   Pass B output.

6. **`finalize` receives context for LLM re-call:** `finalize` needs access to
   the packed context (system blocks + messages) to make the Pass B LLM call.
   Currently it only receives `output_content`. Solution: pass the packed
   context as an additional parameter.

---

## Implementation Plan

### Step 1: Write failing integration tests

New file: `test/cmd/cn_runtime_integration_test.ml`

Tests that exercise the full `finalize` pipeline with mocked LLM:
1. Two-pass: observe+effect → Pass B invoked, both passes receipted
2. Single-pass: effect-only → no Pass B, normal flow
3. Coordination gating: Pass B effects fail → terminal ops blocked
4. Pass-A-safe coordination: ack/surface execute in Pass A, done deferred

### Step 2: Modify `finalize` signature

Add `~packed` parameter carrying the context needed for Pass B LLM call.
Update all three call sites in `process_one` (State 1, State 2, State 3).

### Step 3: Wire Pass B in `finalize`

After `run_pass_a`, when `needs_pass_b=true`:
1. Classify and execute Pass-A-safe coordination ops only
2. Call `repack_for_pass_b`
3. Build Pass B messages (append repack as conversation turns)
4. Call `Cn_llm.call` with Pass B context
5. Archive Pass B IO pair
6. Parse Pass B output
7. Call `run_pass_b` with typed ops
8. Gate and execute Pass B coordination ops

### Step 4: Update projection and conversation to use final output

When two-pass occurred, use Pass B parsed output for projection and
conversation history.

---

## File Change Summary

| File | Action | Lines |
|------|--------|-------|
| `test/cmd/cn_runtime_integration_test.ml` | **New** — integration tests | ~200 |
| `test/cmd/dune` | Edit — add test library | ~8 |
| `src/cmd/cn_runtime.ml` | Edit — wire Pass B in finalize | ~80 |
| **Total** | | ~290 |
