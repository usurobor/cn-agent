# β Review — cnos#669 R6

**Verdict:** REQUEST CHANGES

**Round:** R6/R1

**Reviewed α artifact SHA:** `d552214a8bcd076d46c2fcd4c3740d02d5e6b40f`

**Reviewed tree:** `2e64bbc79d644cad4c460bbea998f1e3da4ee156`

**Product-matter boundary:** `b2aafc79e2c5952dc2ecf22038ae9dfcc1aaa544` (`f85a62401329ced89cc2299d09bf12814c9384c8`)

**R6 base:** `e8ba9954764d58e7b808104d633504e25aa615cc`

**Current review base:** `origin/main` at `7b95787e6a319781f573ceb86af183f68b97e918`

**Fixed this round:** none; fresh context-isolated review

**Branch CI state:** green on the reviewed SHA; findings below are outside the current CI oracle

**Merge instruction:** withheld while this REQUEST CHANGES verdict is open

## §2.0.0 Contract Integrity

| Check | Result | Notes |
|---|---|---|
| Status truth preserved | no | D1: the live/tracked PR body claims symlinked inputs and outputs refuse, but both classes were accepted and published in a direct probe. |
| Canonical sources/paths verified | yes | `SKILL.md` is public under the declared cnos.cdd source path; repository-root authority and explicit target-root paths resolve. |
| Scope/non-goals consistent | yes | State-B, installed activation, standing promotion, merge, and FSM effects remain unshipped/unauthorized. |
| Constraint strata consistent | yes | H01–H13 remain gating; all-pass is still `hold` at standing `none`. |
| Exceptions field-specific/reasoned | n/a | No exception mechanism is added by this cycle. |
| Path resolution base explicit | yes | `repository-root`, `explicit-coh-root`, and `source-checkout-only` are declared and exercised. |
| Proof shape adequate | no | D1: the symlink negative covers only `output/publication`, not the output root or response/invariant inputs named by the PR claim. |
| Cross-surface projections updated | yes | Schema, validator, fixtures, package README, cnos.cdd README, CI, and methodology agree on the State-A surface. |
| No witness theater / false closure | no | D1 is a runtime-enforcement overclaim; C2 is a false R6 tree/diff-equivalence claim in the scaffold. |
| PR body matches branch files | no | The tracked proposal and live body carry D1's false fail-closed symlink claim. |
| γ artifacts present (`gamma-scaffold.md`) | yes | Rule 3.11b passes: `6e07b4cc…` is the first R6 commit after the base and predates every α matter commit. |

The repository contains no `AGENTS.md` at the reviewed SHA or on
`origin/main`; that absence was verified rather than silently substituting an
unrelated file.

## §2.0 Issue Contract

### AC Coverage

| # | Acceptance criterion | Status | Evidence |
|---|---|---|---|
| 1 | Public measurement skill in canonical cnos.cdd source | pass | `measure/recursive-cell/SKILL.md` has `artifact_class: measurement`, `visibility: public`; package smoke contains the surface. |
| 2 | CUE essence plus I5 authority/target validation | pass | CUE schema and full frontmatter self-test pass; full corpus result is `100 SKILL.md validated; no findings`. |
| 3 | Six L0–L4/system manifests replay exact R7 counts | pass | Pinned real `coh` resolved `30/4/4/8/12/11` files and emitted six prompts/reports. |
| 4 | Exact calibration preflight passes R7 and rejects moved/current targets | pass | Detached `a0d39293…` passes with matter `2d6b93cc…` and SHA-256 `80e0d8c6…`; the reviewed checkout refuses on head mismatch. |
| 5 | Same-author self-measure with no semantic/admissibility/held-out standing | pass | Self artifacts preserve the declared provenance and supersession; State-A result remains `hold`, standing `none`. |
| 6 | Existing fixtures/full corpus green plus missing-methodology negative | pass | Positive/negative self-test and all 100 skills pass on the exact SHA. |

### Named Files and Artifact Contract

| Surface | Status | Notes |
|---|---|---|
| `gamma-scaffold.md` | present, γ-owned | First cycle commit; unchanged afterward; C2 applies to two overstrong equivalence sentences. |
| `self-coherence.md` | present, α-owned | Exact reviewed α artifact at `d552214a…`; prior R2/R5 verdicts are correctly marked historical only. |
| `beta-review.md` | present, β-owned | This R6/R1 round; it binds the immutable α SHA and does not edit matter. |
| `gamma-closeout.md` | not yet due | A fresh R6 closeout is prohibited until a fresh β convergence. |
| Live PR #670 | open + draft | No ready mark, label, merge, standing, acceptance, or FSM transition observed. |

## §2.1 Diff, Topology, and Runtime Evidence

### Scaffold, topology, and identity

The R6 parent chain is linear. Its opening edge is
`e8ba9954…` → `6e07b4cc…` (γ scaffold) → `6764c42e…` (first α matter).
All six replayed matter commits have both author and committer
`alpha <alpha@cdd.cnos>`; the scaffold, clarification, and snapshot commits
have both author and committer `gamma <gamma@cdd.cnos>`; `d552214a…` is
α/α. The scaffold path was touched only by `6e07b4cc…`.

### Product-patch equivalence and whitespace

The old
`90287522ebf55b54f875973631506500a5f4f578..ccaf35607520ffb43d9450e35759e30d742355d5`
and R6
`e8ba9954764d58e7b808104d633504e25aa615cc..b2aafc79e2c5952dc2ecf22038ae9dfcc1aaa544`
binary patch streams, excluding
`.cdd/unreleased/669/**`, compare byte-for-byte equal. Both reproduce SHA-256:

`4c42826034a446e0953e0934085dbfd8ee6bb71febff39254c7a34c44a65d9c0`

`b2aafc79e2c5952dc2ecf22038ae9dfcc1aaa544..d552214a8bcd076d46c2fcd4c3740d02d5e6b40f`
is empty outside the cycle directory, and `git diff --check` passes for
`e8ba9954764d58e7b808104d633504e25aa615cc...d552214a8bcd076d46c2fcd4c3740d02d5e6b40f`.
C2 distinguishes this true
patch-stream statement from the scaffold's false tree/direct-diff statement.

### Independent execution

- Composite assembly passes at SHA-256
  `de0414bac1a086d6ed3ac5c7d84cd6293517eb83c920f6c4a0723f5dfac06a3a`.
- `validate-skill-frontmatter.sh --self-test` passes; full validation reports
  `100 SKILL.md validated; no findings`.
- `test-recursive-cell-runner.sh` passes its deterministic six-target,
  H01–H13, witness-shape, aggregate, custody, replay, lock, mutation, failure,
  and atomic-publication suite.
- `smoke-recursive-cell-package.sh` passes package layout and installed
  activation refusal.
- An independent pinned `coh 0.12.0 (26aab50)` run at TSC revision
  `26aab5023f03dc7d0abf82e5fdba20134fc6adad` emitted and ingested six exact
  targets, retained six reports, validated 13 invariant items and the CUE
  publication, computed aggregate `0.5667621299882342`, and emitted
  `hold` / standing `none`. This is synthetic compatibility evidence only.
- `git merge-tree --write-tree origin/main d552214a…` succeeds without an
  unmerged path; GitHub reports the draft PR mergeable/clean.

### CI status

Main protection declares no required status contexts, so the fallback is every
workflow running for this cycle head. Both Build workflows on exact SHA
`d552214a…` completed successfully:

- push: https://github.com/usurobor/cnos/actions/runs/29700479454
- pull request: https://github.com/usurobor/cnos/actions/runs/29700480533

All exposed jobs are successful, including Go, I1, I2, I4, I5, I6, dispatch
guards, workflow parsing, binary verification, and package verification.

## Architecture Check

| Check | Result | Notes |
|---|---|---|
| Reason to change preserved | yes | Methodology, runner, schema, validation, and fixtures have distinct roles. |
| Policy above detail preserved | yes | H01–H13 and disposition policy remain in the methodology/runner contract. |
| Interfaces remain truthful | no | D1: the advertised symlink refusal is broader than the implemented interface. |
| Registry model remains unified | yes | Six exact targets resolve through one pinned registry/manifests model. |
| Source/artifact/installed boundary preserved | yes | Source activation works; installed activation explicitly refuses. |
| Runtime surfaces remain distinct | yes | Measurement, CUE validation, CC/FSM judgment, and operator effects remain separate. |
| Degraded paths visible and testable | no | D1's symlinked root/input bypass is accepted, not surfaced as a refusal. |

## Findings

### D1 — Runtime accepts symlinked inputs and output root despite the PR's fail-closed claim

**Severity:** D

**Type:** contract, honest-claim, mechanical, runtime

**Evidence:** The live and tracked R6 PR body says the runner “refuses changed,
missing, reused, symlinked, raced, or partially published inputs and outputs”
(`r6-proposed-pr-body.md` lines 17–18). In
`recursive-cell-runner.py`, `args.output.resolve()` follows a symlink before
the output root is checked (lines 338–340), while `snapshot_file()` uses
`shutil.copyfile()` without an `lstat`/symlink refusal (lines 85–90) for all
response and invariant inputs (lines 489–495). The shipped symlink test covers
only a symlink created at `output/publication`.

An independent negative probe supplied a symlinked output root and six
symlinked response files. `emit` wrote through the symlink and `ingest`
published successfully:

`PROBE: symlinked output root and all six symlinked response inputs were accepted`

This is a demonstrated false enforcement claim and an evidence-custody gap,
not a theoretical objection.

**Required repair:** Refuse the output-root path itself when it or any
disallowed ancestor is a symlink before resolving/creating it; refuse
symlinked response and invariant inputs before snapshotting; keep the existing
emission/publication internal symlink checks; narrow prose only if some
symlink class is intentionally supported.

**Regression pair:**

- Positive: regular fresh output root plus six regular response files and one
  regular invariant file publishes the six reports, aggregate, and success
  marker.
- Negative: a symlinked output root, any symlinked response input, and a
  symlinked invariant input each explicitly refuse before canonical emission
  or publication and leave no redirected/mixed canonical output.

### C2 — The γ scaffold confuses patch equivalence with tree/direct-diff equivalence

**Severity:** C

**Type:** contract, honest-claim, mechanical

**Evidence:** `gamma-scaffold.md` AC3 says the R6 final product tree outside the
cycle directory is identical to R5 (`ccaf356…`), its implementation-contract
row says the direct product diff must be empty, and its oracle uses
`git diff --quiet ccaf356… HEAD`. That command fails. The exact direct diff
contains three base-inherited files:

- `.cn-sigma/logs/20260719.md`
- `docs/development/board/board-data.json`
- `docs/development/board/index.html`

The narrower claim actually proven by self-coherence and the PR body is that
the two cycle patch streams are byte-identical at SHA-256 `4c428260…`.

**Required repair:** γ must supersede the false tree/direct-diff wording in a
γ-owned clarification with the exact two-range binary patch-stream oracle,
and α must ensure the next readiness record cites only that corrected
contract. Do not rewrite the pre-α scaffold commit or product matter.

## Verdict

**REQUEST CHANGES.** Issue AC1–AC6, rule 3.11b topology, role identities,
product-patch digest, whitespace, deterministic runner suite, real-coh route,
CUE publication, zero-standing boundary, source/install partition, and remote
CI all pass at the immutable α SHA. D1 and C2 remain unresolved, so this round
cannot converge. No merge, ready mark, label/state change, γ closeout,
standing promotion, CC judgment, operator acceptance, or FSM transition is
authorized by this review.
