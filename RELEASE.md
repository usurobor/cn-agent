## Outcome

Coherence delta: C_Σ A (`α A`, `β A-`, `γ A`) · **Level:** L6

CDD lifecycle made executable. Every role now knows when it starts, what it owns, what it must write before returning, and what γ must verify before closure. The lifecycle state table (S0–S12) replaces implicit operator memory with an inspectable contract.

## Why it matters

Recent cycles exposed three systematic gaps: RELEASE.md had no owner, `.cdd/unreleased/` directories accumulated across 10 releases, and α close-outs were missing because α exited before β approved with no re-entry path. All three gaps share a root cause: CDD mixed polling-era persistent-session assumptions with the current sequential bounded dispatch model. This release reconciles the two.

## Changed

- **CDD.md:** coordination model declaration (§1.6), lifecycle state table S0–S12 (§4.1a), role/artifact ownership matrix (§5.3b), closure verification checklist with 10 failure-mode rows and positive/negative tests (§8.1), small-change artifact collapse table (§1.2), α re-dispatch prompt formats (§1.6a).
- **alpha/SKILL.md §2.8:** re-dispatch mechanism for close-out after β approval.
- **gamma/SKILL.md:** α close-out re-dispatch protocol (§2.7), δ preflight gate row 13 in closure gate (§2.10), RELEASE.md and artifact move ownership (§2.6).
- **operator/SKILL.md:** algorithm steps 4–6 for re-dispatch paths, lifecycle table updated, tag gated on `gamma-closeout.md`.
- **release/SKILL.md + post-release/SKILL.md:** β/δ authority split corrected — δ owns tag/release/deploy.
- **release.sh:** automatic `.cdd/unreleased/` → `.cdd/releases/{version}/` move.

## Validation

- β approved R3 after 3 rounds (R1: 2D+1C → R2: 1C → R3: clean).
- All 12 ACs met. Lifecycle state table, ownership matrix, and closure checklist cross-verified.
- 10 stale unreleased cycle directories moved to correct release versions.
- CI green on main.
