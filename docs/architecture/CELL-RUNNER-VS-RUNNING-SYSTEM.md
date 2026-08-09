# Cell runner vs. the running system — honest evaluation and synthesis

**Status:** evaluation note (operator-directed step-back before Phase 3 /
rented cognition). Compares the new GitHub-free kernel
(`src/go/internal/cellkernel` + `cn cell run`, PR #718, FIDO doctrine per
`msg-cn-pi-cnos-cell-runner-fido-functional-44`) against the **currently
running** cell implementation: `cnos.cdd` (CCNF kernel + role skills),
`cnos.cds` (dispatch wake, FSM, lifecycle), `cdd-verify` (V), and the real
closure corpus under `.cdd/`.

Sources read: `cnos.cdd/skills/cdd/**` (CDD.md v4, role SKILLs, close-out),
`cnos.cds/**` (dispatch SKILL, `fsm/transitions.json`, lifecycle/selection),
`commands/cdd-verify/*.go` (V + counterfeit predicates), the live workflow +
golden, and real cycle records (e.g. `.cdd/unreleased/524/`).

---

## 1. The honest headline

The two systems solve **different halves of the same problem**, and each is
weak exactly where the other is strong.

- The **running system** is a *process* implementation: lifecycle, iteration,
  recovery, review discipline, learning capture — encoded as prose gates that
  language models follow, with GitHub as its state machine and git as its
  evidence substrate. It has real operational scar tissue (false-complete
  runs, re-certification-instead-of-repair, stale-ref poisoning) and real
  fixes for each. Its weakness, named in its own docs: *every invariant is a
  prose gate a role may skip*; V checks artifact presence and actor metadata,
  not content derivation; nothing is transactional.
- The **new kernel** is a *trust-boundary* implementation: immutable seat
  scopes, sealed results, positional ownership, one verifiable record, typed
  failure routing — invariants held by types and a verifier, not by prose. Its
  weakness: it runs one in-memory episode of a toy shape. It has **no matter
  substrate** (matter is a string; real matter is a diff at a SHA), **no
  durability** (process death loses everything; the running system reconstructs
  from the branch), **no iteration**, **no findings structure**, and **no
  actor identity** beyond ids minted by the same process.

Neither replaces the other today. The synthesis is: **the kernel becomes the
trust core *inside* the running system's proven process shell**, and before
cognition we port the process lessons the shell learned the hard way.

## 2. What the running system does better (borrow list)

1. **Matter lives on a durable substrate.** A cycle's matter is the
   `cycle/{N}` branch: commits vs a pinned `base_sha`, artifacts on disk,
   reconstructable by any fresh process (*"the branch state IS the iteration
   state"*). The write-fence (`CN_WAKE_BASE_SHA`) + the `if: always()`
   mechanical finalizer checkpoint matter even when the cognitive session
   dies. → Our kernel needs a **matter-substrate seam** (worktree adapter
   outside the kernel; `base_sha`/`head_sha`/diff bound into the record)
   before a rented α can produce anything real.
2. **Review is SHA-bound, multi-pass, and findings-shaped.** Real
   `beta-review.md` frontmatter binds `base_sha` + exact `alpha_commit_shas`;
   β's method is disciplined (scope-compliance against the diff **before**
   content; field-coverage walks; severity classes D/C/B/A; findings →
   dispositions; *"APPROVED is a conjunction"* — auto-RC on contradiction).
   → Our `Review{pass, notes}` is impoverished. **Typed findings** (severity,
   location, disposition) belong in `BetaOutput`; they are what a repair round
   consumes.
3. **Scope contracts.** Cycles declare α's permitted writable paths; β checks
   the diff against them first (cycle 524 did exactly this). → A **writable
   scope** field on the contract is a cheap, powerful containment for rented α.
4. **Iteration and repair are solved problems.** Internal α↔β rounds
   (R0…RN) stay inside the cell; cross-firing repair has the
   **REPAIR-PLAN-first** discipline (finding→action→evidence map written
   before any other file — the #514/#516 scar: three wakes "repaired 0 of 41
   required items" by re-certifying), `repair_evidence` blocks, run-class
   preflight, and **never-blind-requeue-over-matter**. → This is the design
   for our deferred `Drive`, already field-tested.
5. **Actor identity and temporal counterfeits.** V's C1 (α actor ≠ β actor
   via closure metadata + `git log`) and C2 (β's verdict commit precedes δ's
   decision) check things our in-process ids cannot: *who* held a seat and
   *when* relative to the boundary decision. In-process these are vacuous;
   **cross-process (rented cognition) they become necessary again** — bind
   actor/provider identity into `StationRecord` and let V check it at the
   process boundary.
6. **Honest degradation: `mode: collapsed`.** When independence is degraded
   and *admitted*, V downgrades C1/C2 to warnings instead of rejecting. →
   Adopt as a closure concept alongside `simulated`: `simulated` = fake work,
   `collapsed` = real work, degraded independence, honestly declared. (Under
   FIDO, single-model-both-seats without fresh contexts *is* collapsed mode.)
7. **Skill loading tiers.** Tier 1a (role + protocol, unconditional), 1b/1c
   (lifecycle by phase), 2 (eng bundle), 3 (issue-named, e.g. `eng/go`) — with
   *"naming a skill without reading it is not loading it."* → This is the
   worked resolver design for our `alpha_skills`/`beta_skills` lists and the
   Phase-2/3 skill-loading step.
8. **Learning capture (ε).** `gamma-closeout.md` mandates a `learning:` block
   (observations / process_deltas / reusable_patterns / followups /
   operator_burden). → This is precisely the **distilled-knowledge record** of
   the #698 discussion; put it in the closure. (The running system's admitted
   gap — *"detect-recurrence has no owner"* — is ours to fix on the read side.)
9. **Override discipline.** *Degraded ⇔ override ≠ null*; per-predicate
   override only; the original verdict is never rewritten (C3). → Our
   `Override` decision is an enum with no payload contract; adopt this shape
   when δ grows policy.
10. **The outer shell's good architecture.** `transitions.json` as data + a
    generic guard engine; marker-file transition *requests* with the FSM as
    the **sole label writer** (the wake demoted to requester); pre-claim
    recovery scanner; golden/drift CI on the workflow *and on the prompt
    text*. → Keep the shell as-is; it is the S8 adapter's foundation. Note the
    convergence: demoting the wake to requester is the *same move* as our
    kernel-owned γ/V/δ — untrusted parties propose, mechanical authority
    disposes.

## 3. What the new kernel does better (keep list)

1. **The receipt is computable — actually.** CDD doctrine says *"the receipt
   is computable, not authored"*, but operationally the five closure records
   are authored prose and V checks presence, headers, and actor lines. Our
   `EpisodeRecord` + one scope-lift digest + `VerifyClosure` makes the
   doctrine mechanical: verdict → decision → status literally recompute from
   the record. **The old system's aspiration is the new system's
   implementation.**
2. **Invariants by construction, not prose.** Sealed results (unforgeable
   outside the kernel), positional ownership (no authority fields on the seat
   surface to forge), frozen contracts, β's projection isolation — the
   firebreak the running system enforces by instructing a model (δ §9.3's
   context table) is a **function signature** here. The running docs admit V
   *"cannot prove independence"* and C1 is defeatable by declaring collapsed
   mode; our isolation is structural where it can be, with tests.
3. **Typed failure routing.** `contract_unmet` vs `invalid_record`/
   `invalid_identity`, integrity failing closed to `rejected` — the running V
   emits one FAIL and a model (δ) interprets prose predicates.
4. **Deterministic and testable without GitHub.** One binary; adversarial
   race-clean test suite; CUE + live-CLI corpus in CI. The running system
   needs live GitHub to evaluate (`run_active` reads the substrate it's meant
   to be independent of), and its CI can only grep that the *instruction text
   exists*.
5. **Bounded and fail-closed.** Size bounds, UTF-8 artifact contract,
   cancellation between stations, fail-closed identity, explicit exit codes,
   `simulated` for non-authoritative runs. The running wake has **no**
   mechanical bounds at all — no timeout, no turn cap, no round cap; a
   runaway cell burns the 6-hour job budget and dies mid-cell.
6. **No shared mutable state.** The running control plane is labels +
   comments + branch state mutated by three concurrent writers (wake,
   operator, scanner) with TOCTOU windows that are *detected* (release back
   to queue) rather than prevented.

## 4. Where the new kernel is honestly behind (gaps before cognition)

Ordered; the first two block a real CDS episode entirely.

- **G1 — matter substrate.** `Matter{Data string}` cannot carry a code
  change. Needed: contract gains `base_sha` (+ optional writable scope); a
  **workspace adapter** (worktree) materializes α's working copy *outside*
  the kernel; the sealed α output binds `head_sha` and the diff as an
  artifact. The kernel stays substrate-free; the adapter is the seam.
- **G2 — durability.** The closure must land on disk
  (`.cdd/unreleased/{N}/closure.json` beside — eventually instead of — the
  five prose files), so a dead process leaves a reconstructable trail and the
  finalizer pattern applies. The typed closure replaces the *mechanical* role
  of the prose records; the prose narratives survive as artifacts referenced
  from the record.
- **G3 — findings.** `BetaOutput` grows typed findings
  (severity/location/text/disposition); the repair contract derives from
  them.
- **G4 — rounds.** `Drive` (deferred by Pi's gate, correctly) adopts the
  running system's shape: internal α↔β rounds within one cell; REPAIR-PLAN
  discipline on re-entry; never-blind-requeue-over-matter.
- **G5 — actor identity.** `StationRecord` gains provider/actor identity when
  seats cross processes; C1/C2-equivalents return as V predicates at that
  boundary; `collapsed` joins the status/mode vocabulary for honestly
  degraded independence.
- **G6 — learning.** The closure carries the ε `learning:` block (the #698
  distilled-knowledge record).

## 5. Synthesis — the migration posture

**Kernel inside the shell, not beside it.** The running system keeps the
outer loop it is good at (FSM + marker requests + scanner + finalizer +
labels), unchanged. The kernel replaces, in stages, what the shell is bad at:
the *inside* of a cell — where today ~300 lines of prose instruct a model to
route α/β and hand-author closure records, tomorrow the wake's δ invokes
`cn cell run` and gets back a verifiable closure. That is exactly the frozen
plan's S8 ("reduce the GitHub Action to fact observation, core invocation,
and labels/PR projection") — this evaluation confirms S8's shape against the
real code.

Pre-cognition order (proposal): **G1 + G3** (matter seam + findings — the CDS
profile becomes real), then **G2** (durable closure — the finalizer pattern
applies to us), then cognition behind the provider seam, then G4–G6 with the
`Drive` work. Each lands under the existing Pi β gate.

What eventually gets retired from the running system (only after parity):
prose claim mechanics inside the wake (→ CLI), header-grep ledger checks
(→ closure vet), authored receipts (→ computed closures). What is *never*
retired by this work: the FSM, the marker-request pattern, the scanner, the
finalizer, the golden guards, and the operator's label authority.
