# α Close-out — #308

skill(cdd/review): split phase 2 into review/issue-contract, review/diff-context, review/architecture siblings

---

## Cycle summary

Content-preserving structural relocation: split `review/implementation/SKILL.md` (three bundled cognitive modes — issue contract walk §2.0, diff/context inspection §2.1, architecture/design check §2.2) into three independently-loadable sub-skills each owning one reason to change. Orchestrator (`review/SKILL.md`) updated: calls array enumerates all four sub-skills; Phase 2 body loads them sequentially; deferred-split note removed. Monolith deleted. M2-review kata, contract/SKILL.md footer, and schemas/skill.cue comment updated to reference new paths. Single review round; no findings.

## Findings

None. R1 verdict: APPROVED (β, no findings).

## Friction log

- `review/architecture/SKILL.md` initially declared `calls: src/packages/cnos.core/skills/design/SKILL.md` — a cross-package static call that the I5 validator cannot resolve (package-skill-root-relative, not repo-root-relative). Fixed in commit `63866d92` to `calls: []` with prose instruction retained. Same pattern as original `implementation/SKILL.md`; I5 cross-package ref limitation is an existing constraint (O6 from #301 close-out).

## Observations and patterns

- Small-change P3 issue. Six ACs, all content-preserving; no design artifact required (decision was named in the original orchestrator's deferred-split note). Implementation across eight file operations (3 creates, 4 updates, 1 delete) was deterministic.
- The I5 cross-package `calls` limitation surfaced again (see #301 O6). The pattern: when a skill genuinely loads a cross-package skill, the static `calls` field cannot represent it; prose instruction is the current v0 workaround. The fix was one commit on-branch before the readiness signal.
- Self-coherence carried AC evidence to concrete oracle level (frontmatter parse, find oracle, grep, line counts) on first draft. No β findings against evidence quality.
