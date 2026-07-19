# α Self-Coherence — cnos#669 R6 provenance recovery

## Gap

The technically converged cnos#669 lineage ended at
`ccaf35607520ffb43d9450e35759e30d742355d5`, but its matter commits were
authored and committed as γ and it lacked a pre-dispatch γ scaffold. Those
bytes therefore could not become canonical by retroactive attestation. R6
reconstructs the cycle from current-main base
`e8ba9954764d58e7b808104d633504e25aa615cc`, with γ's scaffold first and the
same six matter patches freshly authored and committed by α.

The final R6 matter boundary is
`b2aafc79e2c5952dc2ecf22038ae9dfcc1aaa544`, tree
`f85a62401329ced89cc2299d09bf12814c9384c8`. The later γ commit
`35bfa8beaf4efafad032819fe2ebce88afd1a5ea` adds historical snapshots only;
it does not move product matter or supply a current review verdict.

## Skills

- `cnos.cdd/skills/cdd/gamma/SKILL.md` §2.5: the γ scaffold must precede
  α dispatch.
- `cnos.cdd/skills/cdd/review/SKILL.md` §3.11b: required cycle artifacts must
  be complete before review can bind.
- `cnos.cdd/skills/cdd/harness/SKILL.md` §3: Git author and committer identity
  are part of the role audit trail.
- `cnos.cdd/skills/cdd/alpha/SKILL.md` §2.6 row 14: α verifies its canonical
  identity before committing matter or self-coherence.
- `cnos.cdd/skills/cdd/beta/SKILL.md` pre-merge rows 1 and 4: fresh β must
  verify its own identity and the pre-existing scaffold.
- `cnos.handoff/skills/handoff/artifact-channel/SKILL.md`: α, β, and γ own
  distinct artifacts in a sequential handoff.

## ACs

1. **Scaffold-before-α topology — PASS.** The exact opening edge is
   `e8ba9954764d58e7b808104d633504e25aa615cc` →
   `6e07b4cc42a851a3b35499212f624264ff49557a` (γ scaffold) →
   `6764c42ee0493816b9446c6e9d3b73cd7ffc6a94` (first α matter). The
   scaffold is the first R6 commit and was never edited by α.
2. **Matter provenance — PASS.** All six replay commits preserve the original
   subject and stable patch ID, and both author and committer are exactly
   `alpha <alpha@cdd.cnos>`.
3. **Cycle-patch equivalence — PASS.** The byte streams from
   `git diff --binary 90287522ebf55b54f875973631506500a5f4f578..ccaf35607520ffb43d9450e35759e30d742355d5 -- . ':(exclude).cdd/unreleased/669/**'`
   and
   `git diff --binary e8ba9954764d58e7b808104d633504e25aa615cc..b2aafc79e2c5952dc2ecf22038ae9dfcc1aaa544 -- . ':(exclude).cdd/unreleased/669/**'`
   are identical. Both have SHA-256
   `4c42826034a446e0953e0934085dbfd8ee6bb71febff39254c7a34c44a65d9c0`.
4. **Historical custody — PASS.** γ preserved the pre-R6 R2 snapshots in
   `4d93a9302ca867e09e6cbd0983b4369da35bc0be` and the superseded R5
   snapshots in `35bfa8beaf4efafad032819fe2ebce88afd1a5ea`. Their bytes remain
   unchanged and are evidence only.
5. **Role audit — PASS.** Every R6 α commit through the matter boundary has
   α author and committer identity; every R6 γ commit through the input
   boundary has γ author and committer identity. No β-authored canonical
   R6 review exists yet.
6. **State truth — PASS.** The candidate remains source-checkout-only with
   standing `none`. State-B arbitrary methodology/target loading and installed
   activation remain unshipped. The tracked real-coh route is synthetic
   compatibility evidence, not a semantic sample or calibration warrant.
7. **Control plane — PASS.** PR #670 remains a draft candidate. R6 authorizes
   no ready mark, merge, label change, acceptance, standing promotion, or FSM
   transition.

## CDD Trace

The reconstructed matter mapping is exact:

| Superseded matter commit | Canonical R6 α commit | Subject |
|---|---|---|
| `9cfbcc45fd4a57f9f58741a6c1f27f3577a038d8` | `6764c42ee0493816b9446c6e9d3b73cd7ffc6a94` | `feat: add recursive cell coherence methodology` |
| `2b9b5e51d90e9379202d25cbdd8f39f7424946bb` | `624c7569ecf515e18cfd17607888e12cb84164e5` | `fix(cdd): make recursive-cell measurement replayable` |
| `9cc9c3bbe97b11a94dbab6e47bdc95a2fb0a418c` | `6be13f40ba40a1e21ef1166e4c67cc0f4d66fd01` | `ci: fetch recursive-cell calibration revision for I5` |
| `bc6f6aa946cbbb555bb409c8bc180f5d49685b6e` | `b5fbdb000cd4b3d7f5e559106c129fee43c10177` | `cdd: execute recursive-cell state-a route` |
| `8baf9ead0ff773a23812db1153d1bcfa27bedf2f` | `84cf339fcc6c8f6ed4ffe3a8fb99970bbcfced2e` | `cdd: harden recursive-cell publication` |
| `ccaf35607520ffb43d9450e35759e30d742355d5` | `b2aafc79e2c5952dc2ecf22038ae9dfcc1aaa544` | `cdd: align recursive-cell witness paths` |

The complete pre-artifact R6 topology is linear:

`e8ba9954` → `6e07b4cc` (γ scaffold) → `6764c42e` (α) →
`624c7569` (α) → `6be13f40` (α) → `784eef75` (γ clarification) →
`4d93a930` (γ pre-R6 snapshots) → `b5fbdb00` (α) → `84cf339f`
(α) → `b2aafc79` (α final matter) → `35bfa8be` (γ R5 snapshots).

The immutable historical snapshot set is:

| Snapshot | Preserving γ commit | SHA-256 |
|---|---|---|
| `beta-review-R2.md` | `4d93a9302ca867e09e6cbd0983b4369da35bc0be` | `620590e6b5d476c1f47514db97798a20fb9489d70fe608175261e36d73ed4944` |
| `gamma-closeout-pre-R6.md` | `4d93a9302ca867e09e6cbd0983b4369da35bc0be` | `a7badc28d74929f2743839d9e125205b99ef9ac47a18df5300862736ea79b6cd` |
| `alpha-receipt-R5-comment.md` | `35bfa8beaf4efafad032819fe2ebce88afd1a5ea` | `f540689cf4817f6e380d4f205cb2fb9946f843240a8fb0593a8ccc7638eb3f0e` |
| `beta-review-R5-snapshot.md` | `35bfa8beaf4efafad032819fe2ebce88afd1a5ea` | `c7383f6f96973b61e8bf29ba1bf0252dc1964f19a64837455f7693bca06f835f` |
| `gamma-closeout-R5.md` | `35bfa8beaf4efafad032819fe2ebce88afd1a5ea` | `4e31637431856c86d6e49612c749b6da4f523075e96bba120c0ba9c17550479f` |

The R2 and R5 `CONVERGE` verdicts, their γ closeouts, and any associated CC
recommendations bind only their old immutable lineages. They are explicitly
superseded as current authority. A fresh context-isolated β must review the
exact commit containing this R6 α artifact and return a new exact-SHA verdict;
only fresh β convergence may be followed by new γ binding and CC
re-adjudication.

## Self-check

Final validation on the reconstructed bytes passed:

- `git diff --check` for the cycle and the final working bytes;
- frozen composite assembly with digest
  `de0414bac1a086d6ed3ac5c7d84cd6293517eb83c920f6c4a0723f5dfac06a3a`;
- frozen #662 target preflight at receipt `a0d39293...` and matter
  `2d6b93cc...`, resolving `30/4/4/8/12/11` files;
- the complete deterministic recursive-cell runner, including custody,
  digest, typed-witness, atomic publication, replay, gate, math, and refusal
  cases;
- frontmatter positive/negative self-test and every repository `SKILL.md`;
- `cnos.cdd` package layout plus installed-activation fail-closed smoke;
- relevant workflow YAML parsing and changed shell-script syntax checks;
- a depth-one clone of the reconstructed history after fetching only pinned
  objects `a0d39293...` and `2d6b93cc...`, with all-skill validation and the
  complete runner passing there; and
- pinned real `coh 0.12.0 (26aab50)` six-target emit/ingest using tracked
  synthetic witnesses, producing six reports and an all-pass `hold` result
  with standing `none`.

Stable patch-ID equality was rechecked for all six old/new pairs. The Git role
log was rechecked over `e8ba9954..HEAD`: the six matter commits are α/α,
and `6e07b4cc`, `784eef75`, `4d93a930`, and `35bfa8be` are γ/γ. The γ
scaffold, clarification, and historical snapshots remained byte-identical.

## Debt

- Fresh β review over the exact immutable R6 α artifact commit is required;
  there is deliberately no canonical R6 β verdict yet.
- Only after fresh β convergence may γ create a new closeout and CC perform
  fresh process/coherence adjudication.
- State-B arbitrary loading and installed-package activation remain unshipped.
- The candidate retains zero standing. The real-coh run is synthetic
  compatibility evidence; it supplies no independent semantic sample,
  admissibility, consistency, held-out, hosting-independent, or off-house
  warrant.
- PR #670 must remain draft; merge, ready mark, standing promotion, and FSM
  transition remain outside this artifact's authority.
