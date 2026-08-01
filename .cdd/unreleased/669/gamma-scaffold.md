# γ Scaffold — cnos#669 R6 provenance recovery

**Issue:** `cnos#669`
**PR:** `cnos#670`
**Branch:** `cycle/669`
**cell_kind:** `implementation`
**Mode:** substantial-cycle `repair_pass`
**Base:** `e8ba9954764d58e7b808104d633504e25aa615cc` (`origin/main`)

## Gap

The recursive-cell CM implementation reached technical convergence in the
superseded `54ec50a4...` lineage, but that lineage cannot canonically converge
under CDD:

1. `.cdd/unreleased/669/gamma-scaffold.md` was absent and issue #669 carries
   no discoverable protocol exemption or qualifying wave manifest. Review rule
   3.11b therefore made every prior β approval procedurally unavailable.
2. Every commit in the superseded PR lineage used `gamma@cdd.cnos`, including
   commits described as WC α matter. The audit trail therefore did not make
   α≠β structurally observable, and no β-authored canonical review artifact
   existed in Git.
3. The historical `.cdd/unreleased/669/gamma-closeout.md` contained four
   trailing-space findings in the full PR diff.

R6 reconstructs the cycle from the same base rather than attaching a
retroactive attestation. This scaffold is the first cycle commit. α then
replays the final matter under `alpha@cdd.cnos`; β independently reviews and
commits canonical `beta-review.md` under `beta@cdd.cnos`; γ closes only after
that review. The old SHAs and comments remain historical evidence, not current
authority.

## Skills

- `cnos.cdd/skills/cdd/gamma/SKILL.md` §2.5 — pre-dispatch scaffold gate.
- `cnos.cdd/skills/cdd/review/SKILL.md` §3.11b — binding artifact completeness.
- `cnos.cdd/skills/cdd/harness/SKILL.md` §3 — canonical role identities and
  wrong-author history-rewrite recovery.
- `cnos.cdd/skills/cdd/alpha/SKILL.md` §2.6 row 14 — α identity gate.
- `cnos.cdd/skills/cdd/beta/SKILL.md` pre-merge rows 1 and 4 — β identity and
  scaffold checks.
- `cnos.handoff/skills/handoff/artifact-channel/SKILL.md` — one owning role per
  artifact and sequential α→β→γ handoff.

## Source-of-truth table

| Fact | Authority | R6 realization |
|---|---|---|
| γ artifact required | `review/SKILL.md` §3.11b | this file, committed before α |
| Role identity | `harness/SKILL.md` §3 | author and committer emails per role |
| α matter | final technical tree from superseded R5 | replayed by α without semantic drift |
| β judgment | fresh exact-SHA review | canonical `beta-review.md`, committed by β |
| γ closure | post-β receipt binding | canonical `gamma-closeout.md`, committed by γ |
| Candidate standing | recursive-cell `SKILL.md` | source-checkout-only, standing `none` |

## Surfaces and expected diff

- Reconstructed CM/schema/CI matter already present in the superseded R5 tree.
- `.cdd/unreleased/669/` role-local artifacts and R6 provenance map.
- Git commit author/committer identities and commit topology.
- PR body and review comments rebound to the new immutable SHAs.

No CM semantic expansion, State-B implementation, standing promotion, #662
matter edit, ready mark, merge, label change, or FSM transition is in scope.

## Implementation contract

| Axis | Binding |
|---|---|
| Language | Existing Python, shell, CUE, YAML, JSON, and Markdown bytes are replayed; R6 adds only process Markdown unless a validator exposes drift. |
| CLI integration target | Existing recursive-cell runner and repository CI commands remain unchanged. Git is the provenance substrate. |
| Package scoping | `cnos.cdd` recursive-cell CM remains the only product package in scope. |
| Existing-binary disposition | No binary replacement or new executable surface. |
| Runtime dependencies | No new runtime dependency; Git/GitHub are process evidence channels. |
| JSON/wire preservation | TSC v3.2.4 witness, aggregate, and publication contracts remain byte-for-byte equivalent to superseded R5. |
| Backward compatibility | Product diff against superseded R5 must be empty outside `.cdd/unreleased/669/`; zero-standing and State-A/B boundaries remain unchanged. |

## ACs

1. This scaffold is the first commit after `origin/main` and precedes every α
   matter commit in the canonical cycle history.
2. Every replayed matter/repair commit is authored and committed by
   `alpha <alpha@cdd.cnos>`; γ commits only γ-owned coordination/receipt
   artifacts; fresh β commits canonical `beta-review.md` as
   `beta <beta@cdd.cnos>`.
3. The final product tree outside `.cdd/unreleased/669/` is identical to the
   technically converged superseded R5 matter
   `ccaf35607520ffb43d9450e35759e30d742355d5`.
4. `git diff --check origin/main...HEAD` passes, including the historical γ
   closeout projection retained for provenance.
5. The complete recursive-cell runner, frontmatter, package, frozen-target,
   shallow-clone, and pinned real-coh compatibility evidence remains green.
6. Prior β/γ/CC verdicts are explicitly superseded for authority. A fresh
   exact-SHA β→γ→CC chain binds only the reconstructed lineage.
7. PR #670 remains draft throughout; no actor performs merge, acceptance,
   standing promotion, or FSM transition.

## AC oracle approach

- `git log --format='%H %an <%ae> %cn <%ce>' origin/main..HEAD` proves role
  authorship and topology.
- `git diff --quiet ccaf35607520ffb43d9450e35759e30d742355d5 HEAD -- .
  ':(exclude).cdd/unreleased/669/**'` proves product-byte equivalence.
- `git diff --check origin/main...HEAD` proves whitespace cleanliness.
- Repository test commands and the pinned `coh` route prove unchanged runtime
  compatibility.
- β cites this exact path under review rule 3.11b and commits its own artifact.

## Canonical α artifact shape

α writes `.cdd/unreleased/669/self-coherence.md` using the bare canonical
headers `## Gap`, `## Skills`, `## ACs`, `## CDD Trace`, `## Self-check`, and
`## Debt`. Fix-round evidence stays in α's artifact; α never edits β or γ
artifacts.
