# Working Coding Cell 0.1 — implementation plan

**Status:** Plan · pre-ratification · no code authorized by this document
**Design authority:** `docs/architecture/CELL-SYSTEM-DESIGN.md` @ `8949e01a`,
SHA-256 `28e97a8a52d56aad043b047e9309740e1d13edb4a4a152be6a3efacf4f56557c`
**Scope authority:** `msg-cn-pi-cnos-working-coding-cell-scope-77`
**Donor branch:** `origin/claude/cds-case3-rented-beta` @ `3caff954` —
evidence, never merged wholesale (+5,782 non-doc lines)

Every increment starts from `main`. Nothing here authorizes an issue, a PR, a
merge, or a release.

---

## 0. Three changes to the scope, before the plan

### C1 — HIGH. Skill refs must not be assessment catalogue units.

Scope §"Independent bounded assessment" builds the catalogue from
"admitted-contract obligation/AC IDs, **canonical methodology skill refs**,
`check:matter-nonempty`, and one fixed checker unit", and requires β to return
exactly one `pass | finding | unverified` per unit.

A skill ref is not a decidable obligation. Asked "does this patch satisfy
`cnos.eng:eng/go`?", β has no criterion and no verification route — only a body
of guidance. It must then either guess, which is the exact failure this whole
programme exists to end, or return `unverified` for every skill unit, which
makes coverage noise that V has to ignore. Neither is a review.

The evidence is on this branch: the recorded false finding came from a seat
asked to decide something its input could not decide. Skill units reproduce
that shape deliberately.

**Replacement.** Catalogue units are:

```text
one unit per admitted-contract acceptance-criterion id   (each carries a verification route by construction)
check:matter-nonempty
check:project-verify                                      (the one checker unit)
```

Skills stay in the methodology and shape HOW β judges those units. They are
digested and bound as provenance; they are not things to dispose of. When Coh
arrives, property ids replace criterion ids in this same slot — the shape
survives, which is the point of the position.

### C2 — MEDIUM. Cognitive admission is not needed for the first accepted run.

Scope requires "one bounded bootstrap cognitive admission" in Increment 1. By
the design's own authority table it enforces nothing and declares an
`attested`, unverified judgement. Milestone-closure criterion 2 — "rejects a
deliberately bad issue/design before Alpha" — is satisfied by the structural
gate alone for any deliberately bad issue.

What it costs in 0.1: a third cognitive station, a provider fault class, an
attestation vocabulary, a receipt shape, and one provider round-trip on every
run including the deterministic-fake corpus.

**Proposal: keep the position, defer the arm.** Increment 1 ships the admission
component with a structural arm and a declared, unpopulated cognitive arm, so
0.2 adds an implementation rather than a boundary. The milestone still proves
refusal-before-Alpha with zero cognitive invocations, which is the property
that matters and is cheaper to prove.

**Operator note, because this cuts against a stated interest.** The operator
proposed cognitive issue validation directly. I am arguing about ORDER, not
value: semantic admission is the right long-term answer to "is this issue
executable", and it is the second thing I would build. It is not on the
critical path to one accepted patch.

### C3 — LOW, but it fails closed the wrong way if missed.

The obvious checker recipe includes repo-wide `gofmt -l`. **Measured now:
`gofmt -l src/go` lists 18 files on a clean tree.** A checker running it
repo-wide is red for every candidate regardless of the patch, so every episode
takes `finding` and no run can ever accept. A gate that cannot pass is as
useless as one that cannot fail, and this one would look like the cell working.

The recipe in §3 scopes formatting to paths the matter touches.

---

## 1. Increment map

Each increment starts from `main` and is one PR-sized change.

### Increment 1 — profile, run input, admission

| Kind | Artifact |
|---|---|
| new | `src/go/internal/cdsissue/` — **ported from donor, near-verbatim** |
| new | `src/go/internal/cdsdesign/` — the design facet, same shape as the issue |
| new | `src/go/internal/cellinput/` — `UntrustedRunInput` decode + digest |
| new | `src/go/internal/cdsadmit/` — structural arm; cognitive arm declared, deferred (C2) |
| change | `schemas/cds/spec.cue` — `#CDSIssue` ported; add `#CDSDesign`, `#RunInput` |
| change | `schemas/cdd/spec.cue` — profile subset of `#CellSpec` |
| new | `schemas/cds/fixtures/issue/` — **ported, 15 files** |
| new | `schemas/cds/fixtures/design/`, `.../runinput/` — same single-reason discipline |
| change | `scripts/cell-schema-check.sh` — **port the build-`cn`-from-source fix first** |
| new | `docs/architecture/WCC-0.1.md` — the shipped-subset statement |

**ACs.** (1) A malformed issue, a malformed design, and a subject reference of
the wrong kind each reject with their own reason and **zero provider
invocations**, proven by a recording double that fails if called. (2) The
subject pins once: an authored `HEAD` becomes 40 hex in the bound contract, and
the pinned value is what both stations later receive. (3) A one-byte mutation
of issue, design, or subject changes the bound contract digest. (4) Go and CUE
accept and reject the same corpus, every negative singly-invalid — repair the
one defect and both authorities pass.

**Stop.** If admission needs to know how work is produced or judged, or if the
run input needs a field the design does not name.

### Increment 2 — Git custody, full α, reconstruction

| Kind | Artifact |
|---|---|
| change | `src/go/internal/cellcog/` — port `8dad3d6e`, `34a505c5`, `23121305` with tests |
| change | `src/go/internal/cellwork/cellwork.go` — port base-relative `Diff` (`4e8fe9c8`) + `--binary` |
| new | `src/go/internal/cellwork/subject.go` — **ported from `3caff954`** |
| new | `src/go/internal/cellwork/reconstruct.go` — **ported, including the symlink and bound fixes** |
| change | `src/go/internal/cdspatch/` — refactored: no `workspace`, subject from the contract |

**ACs.** (1) α commits its own work and the change is still measured — the
recorded numbers are 0 bytes HEAD-relative against 112 base-relative. (2) Every
change class survives measurement and reconstruction: added, modified, deleted,
renamed, mode-change, CRLF, unicode and embedded-newline paths, binary,
symlink. (3) `reconstruct(base, matter)` equals an independently applied patch,
compared file by file by the test rather than asserted. (4) A tampered patch
fails to apply and says so distinctly from a bad base and a bad repository.
(5) The caller checkout is byte-, status-, and HEAD-unchanged after a run.
(6) A symlink is reported as a link and its target is **not** followed — the
witness plants a secret outside the repo and asserts its bytes are absent.
(7) The view's bound is decided from metadata: a one-line patch to a large file
allocates less than the bound.

**Stop.** If measurement needs anything from α's own account, or if the
reconstruction needs a channel other than `(pinned subject, sealed matter)`.

### Increment 3 — one methodology, one checker, typed β

| Kind | Artifact |
|---|---|
| new | `src/go/internal/cellmethod/` — load one ordered bundle once, digest it, project two role-wrapped views |
| new | `src/go/internal/cellcheck/` — the closed `cnos.project-verify.v0` recipe |
| new | `src/go/internal/cdsassess/` — catalogue construction + typed β (replaces donor `cdsreview`) |
| change | `src/go/internal/cellcog/` — `Answerer` port ported from `b339d78e` |
| discard | donor `#CDSReviewBetaAuthored.skills` and its fixtures — C1 and the single-bundle rule remove them |

**ACs.** (1) One bundle reaches both seats: identical ordered refs and body
digests, asserted from one contract, and no `beta.skills` key survives anywhere
— proven by a decoder negative, not a grep. (2) Catalogue coverage is exact:
missing, duplicate, reordered, unknown, or malformed β entries are each a
distinct fault, not a review result. (3) A checker `fail` forces `finding` and a
cognitive `pass` on that unit is a fault — the witness makes β assert pass over
a failing checker and requires the run to fault. (4) Unavailable tooling forces
`unverified`, distinctly from a checker that ran and failed. (5) β is offered
`--tools ""` and receives no path, handle, session, or transcript — asserted on
the argv and on the constructor's inputs.

**Stop.** If β needs a third outer input, filesystem tools, or any α-private
state; or if any catalogue unit has no verification route (C1's boundary).

### Increment 4 — verified closure, corpus, real witness

| Kind | Artifact |
|---|---|
| change | `src/go/internal/cellkernel/` — bind the new values through the existing spine; no second kernel |
| change | `src/go/internal/cellrun/` — typed outcomes and documented exits |
| change | `scripts/cell-schema-check.sh` — deterministic-fake corpus, positive and negative |
| new | `docs/architecture/WCC-0.1-EVIDENCE.md` + committed raw closure |

**ACs.** (1) Mutation negatives for issue, design, base, profile, methodology,
patch, assessment, checker evidence, and provider policy — each flips the
closure's self-verification, each asserted separately. (2) Admission refusal,
runtime fault, non-accepting closure, and accepted closure are four
distinguishable results with documented exits. (3) One real Claude α + Claude β
episode accepts a bounded CNOS patch on an exact green head; its patch is
applied independently and the declared checker run against it by hand. (4) The
gate builds `cn` from the revision under review — already true on the donor and
ported in Increment 1.

**Stop.** If binding the new values requires domain semantics in `cellrun` or a
second kernel.

---

## 2. Donor classification

Commit-level, against `origin/main...3caff954`.

| Commit | Subject | Disposition |
|---|---|---|
| `8dad3d6e` | timeout hang diagnostics | **reuse** — self-contained, Inc 2 |
| `34a505c5` | streamed provider output | **reuse** — self-contained, Inc 2 |
| `23121305` | producing seat's real tool surface | **reuse** — with its dispatch-parity test that reads the YAML |
| `4e8fe9c8` | base-relative measurement; `--allowedTools Bash`; honest parity test | **reuse** — the highest-value donor commit |
| `d94ca9f7` | admission at the door; **corpus builds `cn`** | **split**: the corpus fix is **reuse** and lands first in Inc 1; the optional-`task` decision is **discard**, since 0.1 carries a separate run input and the question dissolves |
| `ada2a445` | typed issue, `Contract.Task`, one renderer | **refactor** — `cdsissue` ports near-verbatim; the opaque-slot mechanism is superseded by the admitted contract |
| `3caff954` | pinned subject, `Reconstruct`, symlink + bound fixes | **refactor** — `subject.go` and `reconstruct.go` port with their witnesses; `Contract.Subject` as an opaque slot is superseded |
| `b339d78e` | `Answerer` port + `cds.review` fill | **split**: the `Answerer` port is **reuse**; the fill is **discard** and rewritten as `cdsassess` |
| `5bfe2f69` | `cds.review` CUE overlay, corpus, witnesses | **discard** — encodes `beta.skills`, which C1 and the single-bundle rule remove |
| docs commits | outline → design → corrections | **retain** — `CELL-SYSTEM-DESIGN.md` is the authority; the outline is already deleted |

File-level, the pieces that port with their tests: `cdsissue.go` +
`cdsissue_test.go` + `fixtures/issue/` (15); `cellwork/subject.go` +
`subject_test.go` + `fixtures/subject/` (7); `cellwork/reconstruct.go` +
`reconstruct_test.go` + `symlink_test.go`; the `cellcog` deltas above; the
corpus script's build-from-source preamble.

Discarded outright: `cdsreview.go` and its tests, `#CDSReviewBetaAuthored`,
`cds-review-with-workspace.json`, `cdsissue/seats_test.go` (its property moves
to the methodology projection test), `Contract.Task` and `Contract.Subject`
opaque slots including `opaqueSlotIntegrity` and both `*_test.go`.

---

## 3. The checker recipe, exactly

`cnos.project-verify.v0`, run by the runtime against the **reconstructed
candidate**, never against α's workspace, never from cell JSON:

```text
1. go -C src/go build ./...
2. go -C src/go vet ./...
3. go -C src/go test -count=1 ./...
4. gofmt -l <only the .go paths this matter touches>     # see C3
```

Fixed properties: the recipe is a compiled-in constant; `-count=1` because a
cached PASS is not a result; step 4 is scoped to touched paths because
repo-wide it is red on a clean tree today (18 files) and would make acceptance
impossible. `cue vet` is deliberately **absent** from 0.1 — it needs a `cue`
binary the profile does not declare, and its absence must be `unverified`
rather than a silent skip; that is 0.2.

Provenance: the `go` toolchain is the caller's, recorded by version in the
observation. The observation records recipe id, each step's exit, bounded
stdout/stderr tails, and the digest of the candidate it ran against. A step
that cannot start is `unavailable` → `unverified`, distinct from a step that
ran and failed.

**Bounded review-view limits.** Content bound `1 MiB` total, spent in path
order so the same `(subject, matter)` yields the same view. Per file: content
carried for regular UTF-8 text; deleted paths carry status and no content;
binary named and omitted; over-bound named and omitted with the view marked
truncated; symlinks carry their target path as content and are never followed.
Truncation is always reported, never silent, and a truncated view instructs β
to return `unverified` for any unit needing content it was not shown.

---

## 4. What is unnecessary for the first accepted shell run

Beyond C2: the `AdmissionRefusalReceipt` needs only a typed JSON result and a
distinct exit code in 0.1, not a custody-grade receipt type; telemetry needs
nothing beyond the existing timeout diagnostics; and `check:matter-nonempty`
is already implied by "empty matter is a valid non-accepting result" — it can
be one catalogue unit derived mechanically rather than a declared check.

---

## 5. Draft canonical issue body

> **cell: Working Coding Cell 0.1 — admitted CDS contract to verified local closure**
>
> **Problem.** What exists: a generic episode runner, a rented-cognition
> adapter, and a converged design. What is expected: one shell command that
> takes an admitted issue and design plus a pinned repository, and returns a
> verified patch an operator can inspect and apply. Where they diverge: no
> invocation path connects them, and the reviewing seat cannot see the state
> its judgement is about.
>
> **Impact.** Until this exists, every coherence cycle is driven by hand.
>
> **Source of truth.** `docs/architecture/CELL-SYSTEM-DESIGN.md` @ `8949e01a`
> (design authority); `docs/architecture/WCC-0.1-PLAN.md` (route).
>
> **Scope.** The four increments above.
> **Out of scope.** Everything in scope §"What 0.1 explicitly does not ship",
> reproduced verbatim in the issue.
>
> **Acceptance.** The seven milestone-closure conditions, each naming its
> oracle and each with a mutation witness.
>
> **Closure.** One exact green head on which one real α+β episode accepts a
> bounded CNOS patch, that patch applies independently and passes the declared
> checker, and the caller checkout is unchanged.

---

## 6. Honest limits of this plan

The increment boundaries are guesses about review risk, not measurements. The
one I least trust is Increment 3: catalogue construction, methodology
projection, and typed β are three separable things, and if the catalogue design
moves under contact, β should be split out rather than dragged along.

Nothing here has been executed. The donor evidence is real; the route is not
yet witnessed.
