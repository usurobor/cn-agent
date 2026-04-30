## Outcome

Coherence delta: C_Σ A- (`α A-`, `β A-`, `γ A-`) · **Level:** `L5`

`eng/troubleshoot` fills the live diagnosis gap exposed by the 2026-04-30 identity-rotation dispatch test. When an agent hits a live failure mid-task, it now has a structured triage skill instead of improvised poking. The skill teaches the triage order, hypothesis discipline, one-change rule, and RCA handoff boundary — making the diagnosis path explicit and repeatable.

## Why it matters

Five β dispatch failures across three root-cause classes were diagnosed ad-hoc during the dispatch test. Each failure shared the same anti-pattern: wrong first hypothesis (token limit, rate limit) before checking process state, kernel log, tool output, resource pressure, and parent-process lifecycle. A structured skill forces the cheaper checks first and records negative results so future agents do not re-test ruled-out hypotheses.

## Added

- **`eng/troubleshoot` skill** (#309): live diagnosis skill for active environmental, runtime, process, and tool failures. Includes: six-class triage order (process state → kernel/OOM → tool output → resource pressure → lifecycle/parent → application/model), hypothesis discipline (state before test, oracle before test, cheapest first, one change), three worked examples from the dispatch test (OOM kill, `gh` GraphQL error, background-process lifecycle kill), RCA handoff triggers, and embedded kata.

## Validation

- β merged `cycle/309` with no-ff into `main` (`Closes #309`).
- Non-destructive merge test confirmed zero unmerged paths.
- `scripts/check-version-consistency.sh` passed: VERSION, cn.json, all package manifests agree at 3.65.0.
- Frontmatter manually validated against `schemas/skill.cue` (cue CLI not installed on host — exit 2 = prerequisite missing, not schema failure; all hard-gate and spec-required-exception-backed fields present).
- I5 frontmatter CI job will run on this push; expected exit 2 (cue not available in CI environment) rather than exit 1 (schema violation).
