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

### Cycle #84 — make reflection a core CA continuity requirement

Date: 2026-04-30
α model: sonnet
β model: sonnet
Result: **SUCCESS** — R1 RC (2 mechanical findings), R2 approved, merged. β respected role boundary (no tag/release).

#### Failure 1: AUTH — β cannot review with Read,Write only
- Phase: β (R1 attempt 1)
- Symptom: β sat idle for ~5min, never committed. Could not run `git diff`, `gh issue view`, or any shell commands.
- Root cause: `--allowedTools "Read,Write"` excludes Bash. β needs shell for read-only git/gh operations.
- Recovery: killed β, re-dispatched with `--allowedTools "Read,Write,Bash"`
- Implication for #295: β needs Bash. Role boundary enforcement must be in the skill/prompt, not tool scoping. Future: restricted bash allowlist per role.
- Time cost: ~5min

#### Failure 2: CONTEXT — α fix round couldn't write files
- Phase: α (fix round)
- Symptom: α described the fix but exited without applying it, saying "please grant permission to edit"
- Root cause: unclear — α had Read,Write,Bash. Possibly claude -p file write permission issue with .cdd/ paths, or α misunderstood its own tools.
- Recovery: δ applied the mechanical fix directly (SHA correction + review-readiness section)
- Implication for #295: dispatch should verify α actually committed, not just exited cleanly. Structured exit codes or artifact checks needed.
- Time cost: ~3min

#### Failure 3: RELEASE — package manifests not stamped before tag
- Phase: δ (release gate)
- Symptom: release CI smoke test failed — binary pins cnos.core@3.66.0 but package index only has cnos.core@3.65.0
- Root cause: δ bumped VERSION and cn.json manually but did not run `stamp-versions.sh` to propagate to package `cn.package.json` files
- Recovery: ran stamp-versions.sh, committed, force-pushed tag 3.66.0
- Implication for #295: release gate must include stamp-versions.sh + check-version-consistency.sh. Mechanical dispatcher should run both before tagging.
- Time cost: ~5min

#### Observation: β respected role boundary after skill fix
- β merged and wrote close-out but did NOT tag, release, push tags, or delete branches
- The beta/SKILL.md fix (commit 0e94fe3d) successfully constrained β's behavior
- Implication for #295: role boundary enforcement via skill text works when the skill is clear. Tool scoping (Read,Write only) is too restrictive for β.

---

#### Failure 1: AUTH — β exceeded role boundary (tagged release)
- Phase: β
- Symptom: β tagged 3.65.0, bumped version, updated CHANGELOG, deleted cycle branch, moved .cdd artifacts — all δ responsibilities per CDD §1.4
- Root cause: `--allowedTools "Read,Write,Bash"` gives unrestricted shell. No technical gate prevents `git tag` / `git push --tags`. β prompt said "If Approved, merge" but didn't say "do NOT tag/release." β read CDD, interpreted its scope expansively.
- Recovery: None needed (release was correct), but role boundary was violated
- Implication for #295: dispatch must either (a) scope tools per role (β cannot tag/push tags), (b) wrap bash in a role-specific allowlist, or (c) at minimum include explicit "you do NOT" constraints in prompt templates. CDD says merge is β's authority; tag/release/deploy is δ's.
- Time cost: 0 (no rework), but sets a bad precedent for role isolation

