# α Close-out — #307

skill(cdd/issue): move issue katas to cnos.cdd.kata package

---

## Cycle summary

Moved three embedded katas from `issue/SKILL.md` §5 into `src/packages/cnos.cdd.kata/katas/M5-issue-authoring/`. Updated frontmatter (`kata_surface: external`, added `kata_ref`), replaced §5 body with external pointer section. Mirrored #304's M2-review shape.

## Findings

- **KATAS.md stale catalog** — missed adding M5 row to `docs/gamma/cdd/KATAS.md` on initial implementation. Caught by β R1. Fixed in round 1. Pattern: new M-series bundles require KATAS.md update; this wasn't in α's pre-review checklist.

## Friction log

- α dispatched twice via `claude -p` identity-rotation. First dispatch killed by exec timeout (SIGTERM at 5 min). Second dispatch completed and pushed autonomously.
- Prompt must be piped via stdin; CLI argument parsing breaks on multi-line prompts.

## Observations and patterns

- Small-change P3 issue completed in identity-rotation mode: γ (OpenClaw/Sigma) → α (`claude -p`) → β (`claude -p` Sonnet). Full cycle without manual operator paste.
- β autonomously dispatched α fix-round, re-reviewed, merged, and wrote close-out — all in one `-p` session. This exceeds the expected β scope (β performed α's fix-round internally rather than returning RC to γ for re-dispatch).
- The identity-rotation model works for small-change issues. Larger issues may need tighter role boundaries.
