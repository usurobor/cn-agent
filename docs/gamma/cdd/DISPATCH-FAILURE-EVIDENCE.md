# Dispatch Failure Evidence

Evidence log for #295 (`cn dispatch`). Each triad run records failures, recovery actions, and dispatch-mechanism implications.

## Format

```
### Cycle #<issue> — <title>
Date: YYYY-MM-DD
α model: <model>
β model: <model>

#### Failure N: <category>
- Phase: α | β | γ
- Symptom: <what happened>
- Root cause: <why>
- Recovery: <what γ did>
- Implication for #295: <what cn dispatch must handle>
- Time cost: <wasted minutes>
```

## Categories

- **OOM** — process killed by kernel
- **TIMEOUT** — session exceeded time budget
- **TOOL_ERROR** — gh/git/tool command failed
- **LIFECYCLE** — process parenting, backgrounding, signal handling
- **PROMPT** — α/β misunderstood instructions, wrong scope
- **CONTEXT** — missing file content, stale state
- **AUTH** — permission denied, rate limit
- **MODEL** — token limit, refusal, hallucination
- **WORKTREE** — branch conflict, dirty state, merge failure

## Evidence from prior dispatch tests (2026-04-30)

### Cycle #307 — move issue katas to external package

#### Failure 1: OOM (β round 1)
- Phase: β
- Symptom: process disappeared, no output, no error
- Root cause: 2GB VPS, no swap, two Node processes competing
- Recovery: added 2GB swap
- Implication for #295: dispatch must check available memory before spawning; should surface OOM as structured failure, not silent death
- Time cost: ~5min

#### Failure 2: TOOL_ERROR (β round 2)
- Phase: β
- Symptom: session died after 9 events, tool_result error
- Root cause: `gh issue view` without `--json` hits deprecated GraphQL field
- Recovery: switched to `gh issue view N --json`
- Implication for #295: dispatch prompts must use `gh --json`; dispatch should capture tool errors from stream
- Time cost: ~3min

#### Failure 3: LIFECYCLE (β rounds 3-4)
- Phase: β
- Symptom: silent exit, partial output, no error
- Root cause: background `&` process killed when parent shell ends
- Recovery: switched to foreground-only
- Implication for #295: dispatch MUST run foreground, blocking; no backgrounding
- Time cost: ~10min

---

## Evidence from sequential triad runs (2026-05-01)

### Cycle #309 — eng/troubleshoot skill

Date: 2026-05-01
α model: sonnet
β model: sonnet
Result: **SUCCESS** — R1 approved, 0 findings, 3.65.0 tagged

#### Observations (no failures)

- α completed in ~5 minutes: 1 skill file (361 lines) + 7 self-coherence commits
- β completed in ~8 minutes: full review + merge + release + branch cleanup + artifact move
- β did NOT push main after merge — δ had to push. **Implication for #295:** dispatch must verify remote push after merge, or β prompt must explicitly include `git push origin main`.
- β tagged v3.65.0 and did full release (version bump, CHANGELOG, tag push, branch delete, artifact move) — exceeded typical β scope. Not a failure, but #295 should clarify whether β or δ owns release.
- First attempt at nesting γ inside `claude -p` (γ dispatching α/β) produced no visible output and was hard to monitor. **Implication for #295:** γ-as-subprocess-of-subprocess is fragile; δ should orchestrate α/β directly or γ must have a structured output channel.
- Total wall time: ~15 minutes including prompt composition

#### Failure 1: AUTH — β exceeded role boundary (tagged release)
- Phase: β
- Symptom: β tagged 3.65.0, bumped version, updated CHANGELOG, deleted cycle branch, moved .cdd artifacts — all δ responsibilities per CDD §1.4
- Root cause: `--allowedTools "Read,Write,Bash"` gives unrestricted shell. No technical gate prevents `git tag` / `git push --tags`. β prompt said "If Approved, merge" but didn't say "do NOT tag/release." β read CDD, interpreted its scope expansively.
- Recovery: None needed (release was correct), but role boundary was violated
- Implication for #295: dispatch must either (a) scope tools per role (β cannot tag/push tags), (b) wrap bash in a role-specific allowlist, or (c) at minimum include explicit "you do NOT" constraints in prompt templates. CDD says merge is β's authority; tag/release/deploy is δ's.
- Time cost: 0 (no rework), but sets a bad precedent for role isolation

