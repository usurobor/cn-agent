# β Dispatch — Cycle #344 (Cycle A)

```
You are β. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/beta/SKILL.md and follow its load order.
Issue: gh issue view 344 --json title,body,state,comments
Branch: cycle/344
```

## Context for δ

- Allowed tools: `--allowedTools "Read,Write,Bash"` (β must not tag, release, or run arbitrary commands; Bash is read-only git/gh operations only)
- Git identity for β commits on this cycle: `Beta <beta@cdd.cnos>`
- β dispatched after α signals review-readiness: `status: review-ready` in `.cdd/unreleased/344/self-coherence.md` on `origin/cycle/344`.
- β writes review passes incrementally to `.cdd/unreleased/344/beta-review.md`, committing + pushing after each pass.
- On approval: β merges `cycle/344` into main with `Closes #344` in the merge commit, then writes `.cdd/unreleased/344/beta-closeout.md`.

## What β must review

This is a docs-only cycle (Cycle A of a 3-cycle meta-issue). The diff will contain:

1. **`src/packages/cnos.cdd/skills/cdd/activation/SKILL.md`** — new file, 24 sections, 600–1100 lines
2. **`src/packages/cnos.cdd/skills/cdd/CDD.md`** — ≤3-line pointer to activation skill
3. **`src/packages/cnos.cdd/skills/cdd/operator/SKILL.md`** — ≤3-line pointer to activation skill
4. **`src/packages/cnos.cdd/skills/cdd/post-release/SKILL.md`** — ≤3-line pointer to activation skill
5. **`.cdd/CDD-VERSION`** — cnos commit SHA pin (+ tag if present)
6. **`.cdd/DISPATCH`** — single-line dispatch configuration record
7. **`.cdd/CADENCE`** — `rolling-docs`
8. **`.cdd/MCAs/INDEX.md`** — in-flight MCA table
9. **`.cdd/OPERATORS`** — authorized δ handles table
10. **`.cdd/skills/README.md`** — cnos-as-source declaration (pointer artifact, not a vendored copy)

### A.AC1 — activation/SKILL.md structure

- `rg '^## §' src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` → exactly 24 hits (§1–§24)
- `wc -l src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` → 600–1100
- Every section is substantive (≥3 sentences), no stubs
- No `TODO` or `tbd` markers
- §10 contains a "transport contract" subsection enumerating ≥4 event types; no Telegram-specific assumptions in the contract (A.AC3)
- §11 contains a labeled table of secrets; no example tokens in the file; all secret names in `CDD_*` namespace (A.AC4)
- §22 references the ε role and `epsilon/SKILL.md` for cdd-iteration work-product attribution (per cycle #345 landed context)
- §23 is a numbered checklist with 18–24 steps; each step has a one-line verification; no forward dependencies (A.AC5)

**Honest-claim checks (review/SKILL.md §3.13):**

- **(a) Reproducibility:** Any measurement or example in the skill (e.g., "≤24h SLA," "§23 step count 18–24") must be grounded in the issue body or existing cdd artifacts — not invented.
- **(b) Source-of-truth alignment:** Every term used in `activation/SKILL.md` traces to canonical definitions in existing cdd skill files. Check: "dispatch configuration" traces to `operator/SKILL.md §5`; "identity convention" traces to `operator/SKILL.md §"Git identity"`; "cdd-iteration" traces to `post-release/SKILL.md` Step 5.6b; "ε" traces to `epsilon/SKILL.md §1` and `ROLES.md §1`.
- **(c) Wiring claims:** If activation/SKILL.md says "§8 references `operator/SKILL.md §5`," β grep-checks the cross-reference is actually present in the skill text and points to the correct section.

### A.AC2 — cross-references

- `rg 'cdd/activation/SKILL.md' src/packages/cnos.cdd/skills/cdd/` → ≥3 hits
- Each pointer appears in a contextually appropriate section (not bolted on; integrated)
- No broken relative links (β resolves each path against the repo root)
- One pointer per target file, no duplicates

### A.AC3 — transport-agnostic notification interface

- §10 prose mentions Telegram only as one example
- §10 "transport contract" subsection enumerates ≥4 events
- A hypothetical Slack adapter could be implemented without changing the skill
- No Telegram-specific keys, URLs, or assumptions in the contract section itself

### A.AC4 — concrete secrets prescription

- §11 contains a labeled table (not inline prose) of secrets
- All secret names in `CDD_*` namespace; at minimum `CDD_TELEGRAM_BOT_TOKEN` and `CDD_TELEGRAM_CHAT_ID`
- Explicit prohibition on committing tokens
- No example tokens (not even fake ones)

### A.AC5 — numbered runnable checklist

- §23 step count: 18–24
- Each step has a one-line expected outcome / verification
- Steps are ordered with no forward dependency (step N does not require completing step N+2 first)
- Each step ≤2 sentences

### A.AC6 — cnos self-activation marker files

All six must exist in the diff:

- `.cdd/CDD-VERSION` — contains a valid cnos commit SHA; tag line present if the tip is tagged
- `.cdd/DISPATCH` — single-line or minimal format per the dispatch decision in self-coherence.md; references §5.1 or §5.2 of `operator/SKILL.md`
- `.cdd/CADENCE` — value `rolling-docs`
- `.cdd/MCAs/INDEX.md` — table present (may be empty-row if no in-flight MCAs); structure matches §15 prescription
- `.cdd/OPERATORS` — table present; columns include handle, identity-email, role-scopes-permitted, added-cycle
- `.cdd/skills/README.md` — declares cnos as its own canonical source with the vendored path location

**Self-application honest-claim check:** cnos must pass activation/SKILL.md §24 verification as of the merge commit. β runs the §24 check against the branch state and confirms it passes. If §24 check fails, it is a D-level finding.

### Open-question decisions in self-coherence.md

The issue carries 37 open questions. β verifies:
- `self-coherence.md §Open-Question Decisions` contains a numbered table or numbered list with 37 entries
- Each entry records the decision taken and, if deviating from the recommendation, a rationale
- The skill text implements the decision (e.g., if OQ #1 recommends "SHA + tag one per line," `.cdd/CDD-VERSION` contains two lines; if OQ #2 recommends single-line DISPATCH, `.cdd/DISPATCH` is one line)

If fewer than 37 decision entries exist, request changes citing which question numbers are missing.

## §2.5b disconnect note

This is a docs-only cycle (no version bump, no tag). The disconnect is the merge commit on main. β must NOT run `scripts/release.sh`, bump `VERSION`, or push a tag. Per `release/SKILL.md §2.5b`: the merge commit hash is the disconnect signal. The cycle directory move (`.cdd/unreleased/344/` → `.cdd/releases/docs/{ISO-date}/344/`) is γ's step post-merge.

## Dispatch note

δ: dispatch β with `claude -p` against the cnos repo. β polls `origin/cycle/344` for the review-readiness signal in `.cdd/unreleased/344/self-coherence.md` before beginning the review. Pass the prompt text above. This cycle has a substantial primary artifact (24-section skill file) — β should write review passes incrementally (contract integrity → AC coverage → honest-claim checks → verdict), committing after each pass, to avoid session timeout loss.
