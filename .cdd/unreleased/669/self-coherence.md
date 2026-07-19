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

# R6/R2 Fix Round — β D1 and γ C2

### Binding

This α repair starts from exact γ clarification commit
`3144369b9066b6fbc25c54b20fc241e4a64c5046`, whose parent is canonical β
R6/R1 REQUEST CHANGES receipt
`62135222a4de8d078cfb5cb1108532252dd19366`. The γ commit supersedes the
scaffold's false tree/direct-diff oracle; α did not edit `gamma-scaffold.md`,
`gamma-clarification.md`, or `beta-review.md`.

The new R6/R2 product-matter commit is
`db9ce95b1beb025d24f43fba2df58692665f80be`, tree
`2e4f08e1d2632a8a5fd008a08187a1b42d231152`, with exact parent
`3144369b9066b6fbc25c54b20fc241e4a64c5046`. Both its author and committer
are `alpha <alpha@cdd.cnos>`. It changes exactly four files:

- `recursive-cell-runner.py` — fail-closed lexical path validation;
- `test-recursive-cell-runner.sh` — positive and symlink regression family;
- recursive-cell `SKILL.md` — authoritative custody boundary; and
- recursive-cell `README.md` — executable path behavior.

Design and a separate plan artifact were not required: β D1 prescribed one
bounded filesystem membrane and its exact regression family. The implementation
order was tests/peer enumeration, runner membrane, authority/docs projection,
then full revalidation.

### D1 resolution — filesystem evidence custody

The runner now checks the lexical output path before resolving or creating it.
The output root itself and every existing traversed component must be
nonsymlinked. In ingestion, all six response paths and the invariant-assessment
path are preflighted as regular files with no symlink component before the
publication lock or stage exists, then rechecked immediately before each copy.
The pre-existing canonical `emission/`, `publication/`, staging, and lock
symlink/race guards remain intact.

The regression peer set is complete:

- one regular output root, six regular responses, and one regular invariant
  still produce six reports, the aggregate, and the success marker;
- a direct output-root symlink and a symlinked output ancestor each refuse
  before creating canonical or redirected output;
- each of `cc662-system` and `cc662-l0` through `cc662-l4` is independently
  replaced by a symlink and refuses before staging/publication;
- a symlinked response-directory component and a symlinked invariant input
  also refuse; and
- every refusal leaves no publication, stage, lock, redirected emission, or
  mixed canonical output.

The CUE result/publication schemas and tracked JSON witness fixtures were
audited and intentionally remain unchanged: D1 changes filesystem custody,
not the typed witness or publication shape. Symlink cases are created
deterministically in the shell harness because a symlink is filesystem
metadata, not JSON fixture content. The tracked and live PR claim is now backed
by runtime enforcement; the tracked proposal below narrows the sentence to the
exact output-component and external-input behavior.

### C2 resolution — corrected patch-stream oracle

The appended γ clarification at `3144369b...` is binding. No R6 acceptance
claim uses whole-tree equality with `ccaf356...` or a direct diff from that
commit.

The immutable replay prefix remains proven by two byte-identical binary patch
streams, excluding `.cdd/unreleased/669/**`:

- superseded `90287522...ccaf356...`; and
- reconstructed `e8ba9954...b2aafc79...`.

Both retain SHA-256
`4c42826034a446e0953e0934085dbfd8ee6bb71febff39254c7a34c44a65d9c0`.
That digest is provenance for the pre-repair prefix, not a constraint against
β's required additive repair.

The exact additive product patch
`3144369b9066b6fbc25c54b20fc241e4a64c5046..db9ce95b1beb025d24f43fba2df58692665f80be`
has binary patch-stream SHA-256
`23749673b9fa678ad832b5b4f6915f5bed2b2865896555daf3504e2dd4e62796`.
Applying that exact stream to a detached superseded R5 matter checkout and
diffing from `90287522...` produces bytes identical to the canonical corrected
R6 stream `e8ba9954...db9ce95b...`; both full streams hash to
`055c1b52ea5e83e64cca7a3dd8efe7c054b76f9e170c3fa90e9484a05d81c6a1`.
This is patch-stream equivalence over the intended replay plus reviewed
additive repair, not final-tree or direct-diff equivalence.

### R6/R2 validation and self-check

Validation on the repaired bytes passes:

- `git diff --check`;
- the targeted symlink family and the complete deterministic runner suite
  (`96` captured output lines), including every inherited custody, witness,
  atomic publication, replay, gate, math, and refusal case;
- frozen composite instruction SHA-256
  `de0414bac1a086d6ed3ac5c7d84cd6293517eb83c920f6c4a0723f5dfac06a3a`
  and exact #662 target preflight with `30/4/4/8/12/11` files;
- frontmatter positive/negative self-test and all `100 SKILL.md` files;
- `cnos.cdd` package layout and installed-activation refusal;
- Python and shell syntax plus all seven workflow YAML files;
- a depth-one checkout of the new R6/R2 α history after fetching only pinned
  objects `a0d39293...` and `2d6b93cc...`, with all skills and the complete
  runner green; and
- pinned real `coh 0.12.0 (26aab50)` six-target emit/ingest using the regular
  tracked synthetic witnesses, producing six reports, 13 passing invariants,
  `hold`, and standing `none`.

The change does not push ambiguity to β: the producer/consumer set is the CLI
output argument, six response paths, invariant path, snapshot function,
runner stage, deterministic harness, SKILL authority, README projection, and
tracked PR projection; every member is enforced, tested, or explicitly
unchanged because its typed shape is unaffected. Remote CI on the resulting
review-ready SHA is still a β input and is not preclaimed here.

### R6/R2 debt and handoff

State-A remains source-checkout-only and fail-closed. State-B arbitrary
methodology/target loading and installed activation remain unshipped. The
candidate retains standing `none`; real-coh output is synthetic compatibility
evidence, not a semantic sample, consistency/admissibility claim, held-out
anchor, ratification, or standing promotion.

The R6/R1 β verdict remains REQUEST CHANGES historical authority over its exact
reviewed SHA; D1 and C2 are repaired here but only a fresh context-isolated β
review of the final immutable R6/R2 α head can issue a new verdict. No α action
performs β/γ/CC judgment, comment projection, ready mark, merge, label change,
operator acceptance, standing promotion, or cell/wave FSM transition.

# R6/R3 Fix Round — precise custody claims

### Binding

This α projection repair starts from exact γ clarification commit
`0a871512c812e25fa7c0d3557c6229c74ad0e3fa`, whose parent is the independent
β R6/R2 REQUEST CHANGES receipt
`6b2e75a75c599b734d7984631c16fb914770e8f5`. β reviewed exact R6/R2 α
commit `8df927352c2f323e7ac109a847abd547a225caed`, tree
`ed802cb3dc738614fbc822b8c1b6ceb7b44af7e7`, and found only D3: the current
tracked/live PR body's phrase `raced ... inputs and outputs` generalized beyond
the implemented filesystem boundary.

γ selected the precise-claim repair rather than a new concurrency mechanism.
No runner, test, schema, witness, aggregate, CUE, H01-H13, TSC, State-A/B, or
standing behavior changes here. Product matter remains exact commit
`db9ce95b1beb025d24f43fba2df58692665f80be`, tree
`2e4f08e1d2632a8a5fd008a08187a1b42d231152`, with additive patch-stream
SHA-256 `23749673b9fa678ad832b5b4f6915f5bed2b2865896555daf3504e2dd4e62796`
and corrected full patch-stream SHA-256
`055c1b52ea5e83e64cca7a3dd8efe7c054b76f9e170c3fa90e9484a05d81c6a1`.

### Projection repair and grep audit

The current proposal no longer claims universal refusal of raced inputs and
outputs. It now names only these implemented classes:

- missing or non-regular response and invariant inputs;
- canonical target and invariant-assessment prompts changed after emission;
- response bytes inconsistent with the invariant assessment's bound digests;
- reuse of canonical emission or publication;
- contention for the exclusive publication lock;
- private publication staging followed by atomic rename, so refusal leaves no
  partial canonical publication; and
- static no-symlink-component checks for the output path and the six response
  plus invariant paths before snapshot.

Repository-wide tracked grep found the superseded broad wording in the current
proposal, the historical R5 proposal snapshot, β's D3 finding, and γ's binding
clarification. The current-proposal hit is removed; R5 and the β/γ hits remain
immutable historical/review evidence. Other `raced` hits concern unrelated
dispatch, issue-scan, label-doctor, CI-comment, and historical log behavior.
The historical R5 proposal snapshot is preserved byte-for-byte.

The recursive-cell `SKILL.md` and `README.md` already state the exact static
boundary: the lexical output path must contain no symlink component, and the
six response plus invariant paths must be regular and contain no symlink
component before staging and immediately before snapshot. Their independent
internal emission/publication/staging/lock guards remain specifically scoped;
the earlier R6/R2 shorthand `symlink/race guards remain intact` denotes only
those named internal guards, not a universal pre-snapshot external-pathname
race guarantee.

At the starting SHA, live draft PR #670 was open at exact head
`0a871512c812e25fa7c0d3557c6229c74ad0e3fa`. Its prior κ projection already
carried the exact R6/R2 α artifact identity and the clarified frozen-main-base
description; those live semantics are back-projected here so they do not
regress in the canonical tracked proposal. After this edit, the full tracked
proposal's only semantic delta from the live body is the narrowed custody
bullet; κ live-body projection is deliberately pending.

### R6/R3 validation and self-check

The repair changes exactly two α-owned artifacts:

- `.cdd/unreleased/669/r6-proposed-pr-body.md`; and
- `.cdd/unreleased/669/self-coherence.md`.

`git diff --check` passes. Diff and tree checks confirm no source, runner,
test, schema, or matter delta from `db9ce95b...`; the matter tree remains
`2e4f08e1...`. Current-authority grep confirms the proposal, `SKILL.md`, and
`README.md` contain no universal pre-snapshot pathname-race promise. The live
body comparison confirms all non-D3 proposal semantics already match and only
the κ-owned live projection remains pending. No real-coh rerun is warranted for
a two-Markdown claim repair over unchanged product bytes.

### R6/R3 debt and handoff

State-A remains source-checkout-only and fail-closed. State-B arbitrary loading
and installed activation remain unshipped. The candidate retains standing
`none`; synthetic real-coh compatibility evidence remains neither semantic
calibration nor ratification.

β's R6/R2 REQUEST CHANGES verdict remains historical authority over its exact
reviewed SHA. A fresh context-isolated β must bind the final immutable R6/R3 α
head and verify both tracked wording and the pending live-body projection. No
State-A or State-B transition, γ closeout, CC judgment, operator final read,
comment/body mutation, ready mark, merge, label change, standing promotion, or
cell/wave FSM transition is authorized by this artifact.

### R6/R3 receipt-truth correction

The proposal no longer presents predecessor `8df92735...` and its tree as the
current review-ready artifact. Those identifiers remain only the historical
R6/R2 α receipt reviewed by β. This supersedes the earlier R6/R3 statement
that back-projected that predecessor identity as a current proposal semantic.

The immutable product matter remains `db9ce95b...`, tree `2e4f08e1...`. The
exact R6/R3 α review head is necessarily the commit containing the final
proposal and self-coherence bytes; the artifact does not attempt to embed its
own not-yet-known SHA. κ's handoff and fresh β intake must bind that head
externally. Current pending-review references now say R6/R3; historical R6/R2
round labels and receipts remain unchanged.

This receipt-only correction changes the same two α-owned Markdown artifacts,
leaves all product matter and other custody evidence byte-identical, and
authorizes no κ, β, γ, CC, PR-state, standing, or FSM action.
