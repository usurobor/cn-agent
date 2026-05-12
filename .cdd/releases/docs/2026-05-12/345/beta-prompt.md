# β Dispatch — Cycle #345

```
You are β. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 345 --json title,body,state,comments
Branch: cycle/345
```

## Context for δ

- Allowed tools: `--allowedTools "Read,Write"` (β must not tag, release, or run arbitrary commands)
- Git identity for β commits on this cycle: `Beta <beta@cdd.cnos>`
- β dispatched after α signals review-readiness: `status: review-ready` in `.cdd/unreleased/345/self-coherence.md` on `origin/cycle/345`.
- β writes review passes incrementally to `.cdd/unreleased/345/beta-review.md`, committing + pushing after each pass.
- On approval: β merges `cycle/345` into main with `Closes #345` in the merge commit, then writes `.cdd/unreleased/345/beta-closeout.md`.

## What β must review

This is a docs-only cycle. The diff will contain:

1. **`ROLES.md`** at repo root — new file (250–500 lines), §§1–8.
2. **`src/packages/cnos.cdd/skills/cdd/CDD.md`** — top-of-file pointer to `ROLES.md` naming cdd as instantiation.
3. **`src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md`** — new stub file (30–100 lines).
4. **`src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md`** — single-sentence re-attribution at Step 5.6b.

**Key review targets per AC:**

AC1 — `ROLES.md`:
- `rg '^## §' ROLES.md` → exactly 8 hits (§§1–8).
- `wc -l ROLES.md` → 250–500.
- §1 table: exactly 5 rows; verbs exactly: produces / reviews / coordinates / operates / iterates.
- §3: exactly 6 instantiation-contract fields.
- §5: ≤200 words; names cdw as "separate issue" or "future cycle" explicitly; no `cdw/` directory in the diff.
- §§1–3 contain no instantiation-specific language (protocol-agnostic sections).
- §6 names cdd and cdw as the first two letters in the c-d-X scheme.
- §8 glossary: ≥6 defined terms including role, matter, frame, instantiation, scope-escalation, order-of-observation.
- Honest-claim check (rule 3.13b): every term in §§1–3 traces to canonical usage in existing cdd skill files or the cited external references (Bateson, von Foerster). No new normative claims that contradict existing cdd skill language.

AC2 — `CDD.md` pointer:
- `head -20 src/packages/cnos.cdd/skills/cdd/CDD.md` contains `ROLES.md`.
- One pointer; no claim that role structure is cdd-original.

AC3 — `epsilon/SKILL.md`:
- `wc -l src/packages/cnos.cdd/skills/cdd/epsilon/SKILL.md` → 30–100.
- Cross-references `ROLES.md §1` and `post-release/SKILL.md` Step 5.6b.
- §1 names ε's relationship to δ (often same actor; separable when warranted).
- No claim ε requires a separate human/agent.

AC4 — `post-release/SKILL.md` Step 5.6b:
- `rg 'ε' src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md` → hit at Step 5.6b.
- No other prose changes to post-release/SKILL.md beyond the single re-attribution.

AC5 — Self-application constraint:
- α's self-coherence.md should note that AC5 (cdd-iteration.md `(ε)` attribution) is a γ close-out obligation. β confirms the note exists and is accurate. β does not need to verify the cdd-iteration.md itself — it doesn't exist yet at review time.

AC6 — §5 cdw placeholder:
- `rg 'separate issue\|future cycle' ROLES.md` → hit in §5.
- §5 is ≤200 words.
- No `cdw/` directory in the diff (negative invariant).

**Placement open question:** If α placed `ROLES.md` somewhere other than repo root, β verifies the placement decision is documented in self-coherence.md and the choice is coherent. Per the issue recommendation, repo root is preferred.

## §2.5b disconnect note

This is a docs-only cycle (no version bump, no tag). The disconnect is the merge commit on main. β must NOT run `scripts/release.sh`, bump `VERSION`, or push a tag. Per `release/SKILL.md` §2.5b: the merge commit hash is the disconnect signal. The cycle directory move (`.cdd/unreleased/345/` → `.cdd/releases/docs/{ISO-date}/345/`) is γ's step post-merge.

## Dispatch note

δ: dispatch β with `claude -p` against the cnos repo. β polls `origin/cycle/345` for the review-readiness signal in `.cdd/unreleased/345/self-coherence.md` before beginning. Pass the prompt text above.
