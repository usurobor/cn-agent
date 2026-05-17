---
cycle: 370
date: "2026-05-17"
issue: "https://github.com/usurobor/cnos/issues/370"
merge_sha: "0d9f7498"
findings_count: 2
patches_count: 1
mcas_count: 1
no_patch_count: 0
---

# CDD Iteration — #370

Two cdd-`*`-gap findings surfaced during cycle #370 close-out triage (one cycle, two roles independently catching the same class twice plus one validator-vs-skill drift).

Format per `cdd/post-release/SKILL.md` Step 5.6b.

## Finding 1 — validator-literal vs skill-prose drift (`§CDD-Trace` vs `§CDD Trace`)

**Class:** `cdd-tooling-gap` + `cdd-skill-gap` (one finding, two surfaces drifting from each other).

**Surfaced by:** β R1 (binding finding F1; CI red on review SHA `aa10f902`, Build run `25990518085`); α close-out F1 (root-cause attribution to `alpha/SKILL.md` §2.5 §-shorthand vs validator literal grep).

**Drift surfaces:**

- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md:217` — `**§CDD-Trace** — CDD Trace through step 7` (hyphen in §-shorthand)
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md:352` — resumption enumeration `[Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness]` (hyphen)
- `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify:480` — comment `# Required sections: ..., CDD-Trace` (hyphen) while the same file's grep at L495,L573 is `^## CDD Trace` (space) — internal drift inside the validator
- Convention precedent on `origin/main`: released #367 self-coherence uses `## CDD Trace` (space)

**Class definition:** When a validator's contract is a literal `grep` against a section header (`^## CDD Trace`) and the doctrine surface enumerates section names in a `§Name` shorthand that uses a different separator (`§CDD-Trace`), α may copy the doctrine shorthand verbatim into the artifact and fail the validator. The drift is invisible to every α-side authoring-time check except (a) running the validator or (b) cross-referencing released-cycle precedent on `origin/main`.

**Disposition:** **Patch landed now** — γ step 13a commit `4a0115d2` on main aligns all three drift surfaces to the validator's space form. cn-cdd-verify --unreleased green after patch.

**Patch evidence:**

- `git show 4a0115d2 --stat` shows `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (4 -- / 4 ++) and `src/packages/cnos.cdd/commands/cdd-verify/cn-cdd-verify` (1 -- / 1 ++).
- All three hyphen forms removed; space form now consistent across SKILL.md, validator comment, validator grep, and released-cycle precedent.

## Finding 2 — `extensions.worktreeConfig=true` identity-leak class

**Class:** `cdd-skill-gap` + `cdd-protocol-gap` (role-skill miss; pattern repeating across roles and across cycles).

**Surfaced by:** β R1 row 1 (worktree-config leak during merge-test recipe; β surfaced inside §pre-merge gate row 3 with reference to cycle #301 O8); α close-out F4 (independent surfacing inside the close-out re-dispatch when γ-acting-as-δ overwrote shared `.git/config` between α's session start and first commit attempt).

**Pattern (per α close-out O5):** When the repo has `extensions.worktreeConfig=true` enabled at the shared `.git/config` layer (this repo does) and one or more sibling worktrees carry per-worktree `config.worktree` files, plain `git config user.email X` (no `--worktree` flag) writes to the shared layer. The write succeeds and the immediate `git config --get user.email` returns X. But any subsequent process (a sibling worktree's command, a hook, an unrelated tool) that writes to the shared `.git/config` overwrites the value, and the next read returns the overwriting role's identity. β surfaced the class inside the merge-test recipe (R1 row 1); α surfaced it inside the close-out re-dispatch (this cycle, F4). Two roles, two surfaces, same shared-config blast radius. Convention #301 O8 surfaced the class once before; the gap remained.

**Detection cost:** bounded — one config read at session start (`git config --get extensions.worktreeConfig` returning `true`) determines whether the role must use `--worktree` from the first identity write.

**Disposition:** **Next-MCA committed** — file follow-up issue for an `alpha/SKILL.md` §2.6 + `beta/SKILL.md` §pre-merge gate row 1 + `release/SKILL.md` §2.1 + (future) `delta/SKILL.md` (Phase 4) preventive-write update. The fix is multi-surface and design-shaped (not mechanical); each role's session-start identity step needs a `extensions.worktreeConfig` check that mandates the `--worktree` flag. Bundling under a single follow-on cycle that visits every role's identity-write surface produces one coherent skill patch rather than three drifted edits.

**MCA first AC:** `git config --get extensions.worktreeConfig` returning `true` at session-start makes every subsequent identity write of the form `git config user.email X` (no `--worktree` flag) a binding skill-gate fail; each role-skill's session-start identity step is updated to use `--worktree` from the first write, and `beta/SKILL.md` §pre-merge gate row 1 + the merge-test recipe in `release/SKILL.md` §2.1 carry the same preventive pattern verbatim. Target: a single follow-on cycle producing the multi-role skill patch.

**MCA owner:** γ (cdd-protocol owner); filed as **#373** (https://github.com/usurobor/cnos/issues/373) — "Preventive --worktree identity write across all role skills when extensions.worktreeConfig=true", P2, parent #366.

## Summary

| Finding | Class | Disposition | Patch / MCA path |
|---------|-------|-------------|------------------|
| 1 — validator-literal-vs-skill-prose drift | cdd-tooling-gap + cdd-skill-gap | patch | `4a0115d2` on main (γ step 13a; cycle #370 same-day) |
| 2 — `extensions.worktreeConfig=true` identity-leak class | cdd-skill-gap + cdd-protocol-gap | next-MCA | #373 (filed pre-closure; P2; parent #366) |

## Aggregator update

`.cdd/iterations/INDEX.md` row added at cycle #370 same-day; this artifact's path is `.cdd/releases/docs/2026-05-17/370/cdd-iteration.md` after the cycle dir move.
