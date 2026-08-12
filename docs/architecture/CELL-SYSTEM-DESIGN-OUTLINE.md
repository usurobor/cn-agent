# Cell System Design — Outline

**Status:** Design outline · pre-ratification · not an implementation contract
**Purpose:** Establish the questions, types, boundaries, and proof obligations the
canonical cell design must settle before Case 3 or later cell work continues.
**Review mode:** Pi/Sigma design convergence. Review this artifact as architecture;
do not implement around an unresolved question.
**Current baseline:** branch `claude/cds-case3-rented-beta`, parent `d94ca9f7`.
That head carries Case-2 as merged **plus** unmerged Case-3 work this document
describes as executed: the `Answerer` port, the `cds.review` fill, the typed
`contract.task` slot with `internal/cdsissue`, and the matter gate. Statements
below marked `executed` mean executed *there*, not on `main`.
**Companions:** CCNF owns the substrate-independent closure kernel; TSC owns
coherence-methodology semantics; this document will own their end-to-end CNOS cell
realization once promoted from outline to design.

> **This is deliberately an outline of the design, not the design itself.**
> Working positions are marked as such. Questions remain questions until Pi, Sigma,
> and the operator converge on them. JSON below is a candidate used to test whether
> the conceptual model is complete; it is not yet a frozen schema.

---

## 0. Why this outline exists

### Named incoherence

CNOS has a kernel equation, several strong conceptual papers, a generic episode
runner, and a CDS Case-2 realization, but it has no single current-state design
that names all cell components, their construction, their inputs and outputs,
their authority, and their failure paths.

Implementation is therefore discovering architecture locally. The same missing
design has appeared as several apparently separate defects:

- alpha was offered file editing without the engineering capability needed to run
  the tests that justify its work;
- beta was asked for a verdict from a diff that did not contain enough information
  to verify its own claims, and invented a concrete false defect;
- the issue was referred to in a one-line goal but was absent from the episode;
- admission, production, review, evidence, telemetry, and boundary policy have
  been added one at a time without a common type model;
- acceptance criteria were treated as if they enumerated every affected surface,
  leaving no rule for newly discovered impact outside those criteria;
- implementation and documentation repeatedly made stronger claims than the
  executed boundary established.

Sigma's current Case-3 r0 records the first two failures directly. The active Case-3
branch then began adding a typed issue and admission logic. Those changes are useful
evidence, but they are not a substitute for deciding the cell design first.

### Executed Case-3 evidence for the design-first pause

Sigma's design-first event 68 adds a stronger empirical claim: eight material defects
were all discovered after implementation even though each was decidable from a design.
The outline accepts that evidence and requires the final design to make each class
falsifiable before code resumes.

| Observed defect | Missing design proof |
|---|---|
| Diff measured from mutable `HEAD`, not the pinned base | complete value flow and consumer catalogue |
| Alpha withheld Bash and called that containment | enforce/declare/claim-nothing authority table |
| Shell was available but ordinary commands were not approved | separate capability availability, approval, and containment semantics |
| Beta emitted a confident false finding from insufficient input | per-seat input-sufficiency argument and an `unverified` outcome |
| New task matter was not traced into every closure/schema surface | cross-boundary value and scope-lift catalogue |
| An admission barrier leaked into unrelated specs and fixtures | barrier location and type-boundary specification |
| CUE and Go implemented different definitions of nonblank | rule owner, mirror, and executable parity witness |
| Local corpus ran a stale binary | oracle artifact provenance and proof-limit statement |

These are not eight independent patch requests. They are evidence that the final design
must expose the system's dataflow and authority before another Case-3 implementation
round.

### Challenged assumptions

This design challenges the following assumptions:

1. An issue or prose goal alone is enough to contract implementation.
2. A reviewer can reconstruct hidden context from candidate matter.
3. Withholding a capability constitutes containment.
4. Skills are merely prompt fragments rather than part of a methodology.
5. Acceptance criteria are a complete impact boundary.
6. Authored configuration, normalized IR, compiled plan, and run request can be one
   undifferentiated JSON object.
7. Another local implementation round can settle these questions more cheaply than
   a design.

### Design discipline

The final design must satisfy both current CNOS design disciplines:

- `cdd/design`: named incoherence, constraints, challenged assumption, impact
  graph, typed proposal, leverage, negative leverage, alternatives, and falsifiable
  acceptance criteria;
- `cnos.core/design`: one reason to change per boundary, policy above detail,
  truthful interfaces, explicit runtime surfaces, no duplicate source of truth,
  and no service-locator or registry drift.

---

## Part I — The seven inventories

These come first deliberately. They are the load-bearing part of the design:
concrete, checkable against the code that exists today, and each one closes a
class of defect that this workstream discovered only after implementation.
Part II's narrative sections explain decisions the inventories cannot express;
where the two disagree, the inventories are what a reviewer checks.

**Status vocabulary used throughout.** `executed` — exists and runs today;
`partial` — exists but not in the position this design gives it; `designed` —
decided here, not built; `absent` — named as a gap with no decision yet.

**Scope.** These inventories describe the CDS implementation cell. A second
cell kind will re-instantiate them; the shapes are meant to survive that, the
CDS-specific rows are not.

---

### Inventory 1 — Components

One row per component: its single reason to change, what it owns, and what it
must not know. "Must not know" is the load-bearing column — most boundary
erosion in this codebase has been a component learning something it had no
business learning.

| Component | Realization | One responsibility | Owns | Must not know | Status |
|---|---|---|---|---|---|
| Invocation runner | `cellrun` | turn one invocation into a closure and an exit code | argument parsing, contract reading, closure encoding, exit mapping | any fill's semantics; any provider; git; what a task is | executed |
| Composition root | `cellfills` | bind fill ids to constructors once | the registry passed to `Build` | how any fill works internally | executed |
| Spec normalizer | `cellspec.Parse`, `.Resolve` | source JSON → normalized IR, purely | exact key language, hole grammar, parameter domains, presence-vs-absence | fills, CDS, issue shape, providers, the filesystem | executed |
| Plan compiler | `cellspec.Resolved.Build` + `cellfill` | normalized IR → constructed components | construction order, strict per-fill decoding, the one construction-time effect | seat semantics; what a prompt means | executed (not yet a distinct `CompiledCellPlan` type) |
| Protocol kernel | `cellkernel` | run the fixed closure protocol and emit one receipt | `Contract`/`Matter`/`Review`/`Artifact`/`Receipt`/`Closure`, sealing, γ, V, δ, lift, the single digest | providers, git, CDS, what `Task` bytes mean | executed |
| Cognition adapters | `cellcog` | one bounded, stateless provider invocation | argv recipes, timeouts, output bounds, truthful execution mode | fill semantics, worktrees, what counts as evidence | executed |
| Skill loader | `cellskill` | refs → ordered bodies with content digests | ref grammar, load order, digesting | why a skill was selected | executed |
| **Subject adapter** | `cellwork` | materialize, measure, and reconstruct the subject | `materialize(base)`, `measure(workspace, base)`, `reconstruct(base, matter)`, disposability | contract semantics, prompts, verdicts | **partial** — exists, but constructed privately by the producing fill |
| **Admission** | `cdsissue` + designed `cds.work-contract-admission` | decide whether the run input is executable, before production | the typed issue/design contract, structural rules, the cognitive semantic check, the admission receipt | how the work is done; how it is judged | **partial** — the predicate exists; the component and its receipt do not |
| Producing fill | `cdspatch` | construct and run the producing seat | prompt composition, artifact identity | how the subject is cut (should be the adapter's); how it is judged |  executed |
| Assessing fill | `cdsreview` | construct and run the assessing seat | verdict decoding, refusal to judge unreviewable matter | the producing seat's internals |  executed — but see the CDS predicate in Inventory 3 |
| **Methodology bundle** | — | be the ONE normative source of assessable obligations | the obligations, and both projections of them | which seat is reading it | **absent** — today two independent skill lists |
| **Oracle suite** | — | run declared mechanical checks and emit their receipts | check identity, artifact provenance, proof limits | matter semantics; verdicts | **absent** — `scripts/cell-schema-check.sh` is one, and it lives outside every cell |

**What Inventory 1 already shows.** Three of thirteen components are absent or
mispositioned, and each corresponds to an unresolved question the outline
raised independently: the subject adapter (Pi decision 1), the methodology
bundle (decision 2), and the oracle suite. That is the inventory doing its
job — it makes a missing component a visible empty row rather than a defect
found later.

---

### Inventory 2 — Value catalogue

Every value that crosses a component boundary: who produces it, who consumes
it, which schema surfaces must admit it, and whether it enters the one
scope-lift digest. This closes the class that produced the HEAD-versus-base
measurement bug and the closure-schema omission — both were values whose
consumer list was never written down.

**Rule this inventory exists to enforce:** a value may not gain a consumer or
a schema surface without gaining a row here first.

| Value | Producer | Consumers | Schema surfaces | In digest |
|---|---|---|---|---|
| cell source bytes | author / future compiler | `cellspec.Parse` | `cdd.#CellSpec`, `cds.#CDSCellSpec` | no (its resolved form is) |
| parameters (`--param`) | invoker | `cellspec.Resolve` | `#Param` declaration only; values are Go's authority | no |
| `contract.id`, `contract.goal` | source | both seats, record, V | `#CellSpec`, `#EpisodeRecord` | **yes** |
| `contract.task` (opaque bytes) | source | `cdsissue.Admit` at both seats; record | `#CellSpec.contract.task?` (open), `#CDSIssue` (closed), `#EpisodeRecord.contract.task?` | **yes** |
| `contract.required_evidence` | source | V's producer-authority check, record | `#RequiredRef`, `#CDSCellSpec` order rule | **yes** |
| resolved spec (per-seat declarations) | `Build` | record; `VerifyClosure` re-compares against `RunMeta` | `#CDSPatchAlphaResolved`, `#CDSReviewBetaResolved` | **yes** |
| execution mode | `cellcog.New` (follows the provider) | `RunMeta`, record, `lift` | `#EpisodeClosure` | **yes** |
| `base_sha` (pinned) | subject adapter at materialization | **the measurement** (`measure`), the `base_sha` artifact, and — designed — `reconstruct` | `#CDSPatchAlphaResolved.workspace.base_sha` (40 hex) | **yes**, via artifact and resolved spec |
| workspace directory | subject adapter | producing seat only | none (runtime-local, disposable) | no |
| matter (unified diff) | subject adapter's `measure` | assessing seat; record; the `diff` artifact | `#EpisodeRecord.matter` | **yes** |
| evaluation view | designed `reconstruct(base, matter)` | assessing seat | none yet | no — derived, see Inventory 4 |
| artifacts `base_sha`, `diff` | producing station | V's required-evidence check; record | `#EpisodeRecord.alpha.artifacts` | **yes** |
| skill bodies | `cellskill.LoadAll` | prompts only | none | **no** — deliberately |
| skill refs + sha256 | `cellskill.LoadAll` | resolved declarations, record | `…Resolved.skills` | **yes** |
| review verdict | assessing seat | record; V; δ | `#EpisodeRecord.review` | **yes** |
| verdict / decision / status / repair | γ→V→δ→lift, mechanically | closure; exit code | `#EpisodeClosure` | derived, re-derivable |
| scope-lift digest | `sha256(canonicalBytes(record))` | `VerifyClosure`, any downstream verifier | `#EpisodeClosure.receipt.scope_lift_digest` | is the digest |
| admission receipt | **designed** | closure; downstream custody | **to be defined** | **must be** |
| oracle receipts | **designed** | V; assessment | **to be defined** | **must be** |
| provider credentials | ambient operator configuration | provider child process only | none, ever | **never** |
| telemetry | provider stream | diagnostics on the error path only | none | no |

**Two rows are the design's current debt.** The admission receipt and the
oracle receipts have no schema surface and no digest binding, which is exactly
the state `contract.task` was in before it was traced. They are written here
before they are built.

---

### Inventory 3 — Authority table

Per component: what it **enforces** (something breaks if violated), what it
merely **declares** (a truthful statement carrying no enforcement), and what it
**claims nothing** about. This closes the class that produced the withheld-Bash
mistake, the availability-versus-approval confusion, and the admission barrier
leaking into unrelated fixtures.

**Rule this inventory exists to enforce:** any sentence in code or docs
asserting a property must appear in the *enforces* column of some row, or be
rewritten as a declaration.

| Component | Enforces | Declares | Claims nothing about |
|---|---|---|---|
| `cellspec.Parse`/`Resolve` | exact keys, case-sensitivity, hole grammar, parameter domains, required-parameter presence | that a normalized IR has no unresolved holes | whether the declared components can be constructed |
| `cellfill` strict decoding | closed per-fill key language, including case-insensitive Go aliases | that a fill's arguments are exactly its own | whether the arguments are semantically sensible |
| `cellkernel` | seat isolation by construction (unexported sealed values), record validity bounds, digest recomputation, producer authority for required evidence, re-derivability of verdict/decision/status | the declared protocol id (provenance only; `protocol_validated=false`) | what any fill's declaration means |
| `cellcog` argv | that a cell cannot smuggle argv, env, executable paths, or safety overrides | the offered tool surface and the declared permission mode | **containment of any kind** — a seat with Bash reaches the whole host |
| `--safe-mode` | suppression of USER and PROJECT customization | that this cell's digested skills are the only context *this cell* contributes | vendor-managed substrate policy above the baseline |
| Subject adapter | that matter is computed from the worktree against the pinned base, never taken from a seat's account | disposability of the workspace | anything a seat did outside the workspace |
| Admission (structural) | the typed issue contract: closed keys, non-blank fields, unique criterion ids, a verification route per criterion | that the input is *well-formed* | that the input is *executable* |
| Admission (cognitive, designed) | nothing | an **attested, unverified** judgement that the issue and design are semantically executable | correctness — it is cognition, and it is not itself reviewed |
| Assessing seat | that a verdict decodes to the requested schema or the seat fails; that unreviewable matter is refused without renting cognition | its findings and its citations | anything it could not check from its input — hence `unverified` |
| `scripts/cell-schema-check.sh` | that both authorities agree over one corpus; that the CLI under test is built from source | what green proves (Inventory 6) | any property of a rented provider |

**Barrier locations, stated once.** A barrier belongs at exactly one place and
must not follow the guarded thing around:

| Barrier | Sits at | Explicitly does NOT sit at |
|---|---|---|
| issue admissibility | the door — `Admit`, before cognition is rented | the shape of every CDS document (`#CDSCellSpec` leaves `task` optional for exactly this reason) |
| matter reviewability | the assessing fill, before cognition is rented | the kernel, which never inspects matter |
| capability | the execution substrate | the tool list, which declares and does not contain |
| provider credentials | ambient operator configuration | cell JSON, receipts, telemetry — at any point |

---

### Inventory 4 — Seat specifications

Per cognitive station: the decision it makes, the input it receives, the
**argument that the input suffices for that decision**, its representable
outcomes, and what makes it independent. This closes the class that produced
the confident false finding.

**Rule this inventory exists to enforce:** a station may not be asked for a
decision whose sufficiency argument is missing. If the argument cannot be
written, either the input changes or the decision does.

#### 4.1 Admission

| | |
|---|---|
| **Decision** | is this run input executable? |
| **Input** | the untrusted run-input envelope: contract bytes, design bytes, subject reference |
| **Outcomes** | `admitted` · `rejected` · `incomplete` · `failed` |
| **Sufficiency (structural)** | complete. Every structural rule is decidable from the bytes alone. |
| **Sufficiency (cognitive)** | **partial, and labelled as such.** A model reading Markdown cannot prove total obligation coverage. The outcome is an attestation, not a verification, and must be receipted as one. |
| **Independence** | it judges input, not work; it never sees matter |

#### 4.2 Production

| | |
|---|---|
| **Decision** | none — it produces, it does not judge |
| **Input** | frozen contract (including the admitted issue), the constructive projection of the methodology bundle, a materialized workspace |
| **Product** | the measured diff — **not** its own account of what it did |
| **Capability required** | the full practical engineering surface, including a shell. A seat that cannot run the project's tests cannot check its own work, and produces plausible code it has no way to verify. Withholding that was never containment (Inventory 3). |
| **Outcomes** | matter, possibly empty; empty matter closes the episode unmet |
| **Sufficiency** | the contract must state acceptance criteria and their verification routes; a one-line goal is insufficient and was the root of the false-verdict episode |

#### 4.3 Assessment

| | |
|---|---|
| **Decision** | for each declared obligation: is it satisfied by this matter? |
| **Input** | frozen contract (same bytes as production received), matter, the adversarial projection of the same methodology bundle, and the runtime-derived evaluation view |
| **Outcomes** | `pass` · `finding` (with citation) · `unverified` · `failed` |
| **Sufficiency** | this is the argument the design turns on, below |
| **Independence** | reconstruction, not blindness — below |

**Sufficiency argument.** An obligation is decidable by assessment when its
verification route resolves against the evaluation view. Routes that resolve:
"the diff adds file X", "X states Y", "no call site of Z remains". Routes that
do **not** resolve without an oracle: "the tests pass", "it vets", "it builds".
Those are the oracle suite's obligations, and assessment must be given their
receipts rather than asked to guess — which is why the oracle suite is a
component and not a convenience.

Where a route resolves against neither, the outcome is `unverified`. That arm
is not a courtesy; it is what makes the sufficiency argument honest. Without
it a station under-informed for its question must guess, and the recorded
episode shows what guessing produces: a specific, confident, false claim that a
file lacked an import the file contains at line 126.

**Independence argument.** Assessment's evaluation view is
`reconstruct(base_sha, matter)` — a deterministic runtime function of two
values already inside the CCNF pair. It adds **no information** to
`(contract, matter)`; it changes the **form** from a patch to a tree.

Therefore independence does not rest on the seat being unable to look. It rests
on there being no channel from production to assessment except matter, which
assessment reads anyway. Production cannot influence the view except by
changing the patch, and the patch is the thing under review.

This is checkable rather than promised, and the check is the property to build:
the reconstruction takes exactly `(base_sha, matter)`, and the assessing seat's
constructor receives no other source. It also survives assessment gaining
tools, which the previous property — no tools, therefore cannot look — did not.

---

### Inventory 5 — Rule ownership

Every rule stated in two authorities needs a named owner, a named mirror, and a
**mechanism that checks they agree**. This closes the class that produced eight
whitespace runes admitted by CUE and rejected by Go under a comment asserting
the two were the same predicate.

**Rule this inventory exists to enforce:** a sentence claiming two authorities
agree is not an agreement mechanism. Name the artifact that fails when they
diverge.

| Rule | Normative owner | Mirror | Agreement mechanism | Status |
|---|---|---|---|---|
| parameter/hole name grammar | `cdd.#ParamName` | `cellspec.validParamName` | shared corpus; a name legal in one and not the other fails a fixture | executed |
| concrete-vs-hole value | `cds.#Concrete` | `cellspec` hole detection | `fixtures/invalid/cds-bad-hole-name.json` | executed |
| provider/model pairing | `cds.#Cognition` | `cellcog.New` | `cds-fake-with-model.json`, `cds-modelless-provider.json` | executed |
| issue admissibility | `cds.#CDSIssue` | `cdsissue.Admit` | one corpus read by both: `fixtures/issue/`, 2 positives + 13 single-reason negatives | executed |
| blankness | one pattern string | transcribed into both | Go: predicate compared against `unicode.IsSpace` over the whole rune space. CUE: a fixture carrying the entire whitespace set in one field, so dropping any rune makes it vet clean and fails the gate | executed |
| closure shape | `cdd.#EpisodeClosure` | `cellkernel` record + `VerifyClosure` | live cells' output vetted in the corpus | executed |
| evidence order (diff first) | `#CDSCellSpec` structural list | — | `cds-diff-not-first.json` | executed (CUE-only by choice) |
| producing tool surface | `.github/workflows/cnos-cds-dispatch.yml` | `cellcog.CodingToolSurface` | parity test that **reads the YAML** | executed |
| methodology obligations | **the bundle** (designed) | both projections | **absent** — the projection property is the mechanism to build | designed |
| admission result vocabulary | **absent** | — | — | designed |

**The pattern worth naming.** Every executed row's mechanism is a *fixture or a
test that fails*, never a comment. The two designed rows have no mechanism yet,
and that is the honest reading of their status.

---

### Inventory 6 — Gate specification

Per oracle: how the artifact under test is **obtained**, what a green result
proves, and — the column that matters — what it does not. This closes the class
in which the shared corpus ran a stale binary and reported green after the
guard it was testing had been deleted.

**Rule this inventory exists to enforce:** an oracle must state how it obtains
what it measures. An oracle that does not is measuring something unknown.

| Oracle | Artifact under test | How obtained | Green proves | Green does NOT prove |
|---|---|---|---|---|
| `go test -race -count=1 ./...` | the source tree | compiled per run; `-count=1` because a cached PASS is not a result | the assertions that exist hold | that the assertions can fail — only mutation shows that |
| `go vet ./...` | the source tree | compiled per run | no vet-detectable defect | nothing semantic |
| `cue vet … -d '#X'` | a fixture or a live closure | read from disk | the document satisfies that definition | nothing about the other authority — that is Inventory 5's job |
| `scripts/cell-schema-check.sh` corpus | fixtures **and the CLI** | **the CLI is now built from source per run**; previously `./cn` from the repo root, which meant local runs measured a possibly-stale binary while CI (which builds first) did not | both authorities agree over the corpus, and the live cells close and vet | anything about a rented provider — the corpus rents only the deterministic fake |
| live rented-cognition episode | one real provider run | run by hand from an immutable clean commit; raw closure committed | that the declared authority sufficed on that occasion | reproducibility — the mode is `cognitive` and honestly irreproducible |
| declared cell oracles | the candidate workspace | **designed** — run by the runtime over the evaluation view, receipts to assessment and V | the declared mechanical checks | any obligation no oracle covers |

**The proof-limit statement, once, plainly.** No gate here proves the absence of
defects. Each proves that a specific stated property held for a specific
obtained artifact. Where a gate's artifact provenance is unstated, its green
result is evidence about an unknown object — which is what happened, and what
this inventory exists to prevent recurring.

---

### Inventory 7 — End-to-end flow

Every gate in order, with its typed failure and its exit. `executed` rows are
today's behaviour; `designed` rows are this document's decisions.

| # | Step | Gate | On failure | Exit | Status |
|---|---|---|---|---|---|
| 1 | parse arguments | flags well formed, no duplicates | usage error | 2 | executed |
| 2 | read contract | readable | usage error | 2 | executed |
| 3 | `Parse` | exact keys, case, hole grammar | spec rejected | 2 | executed |
| 4 | `Resolve` | required parameters supplied, values in domain | resolution rejected | 2 | executed |
| 5 | `Build` | fills known; per-fill strict decode; skills load; providers construct | construction rejected | 2 | executed |
| 6 | **admit** | structural, then cognitive | `rejected`/`incomplete` | **1 with an admission closure** | **designed** — today the seats refuse and **no closure is emitted at all** |
| 7 | mint identity | non-empty, pairwise distinct | fail closed before production | 2 | executed |
| 8 | freeze contract + metadata | clone; no seat holds a shared reference | — | — | executed |
| 9 | **materialize subject** | base resolves to a commit | subject failure | 2 | **partial** — inside the producing fill today |
| 10 | produce | provider runs within bounds | seat malfunction | 2 | executed |
| 11 | **measure** | diff against the **pinned base**, not `HEAD` | measurement failure | 2 | executed (since `4e8fe9c8`) |
| 12 | seal α | matter and artifact bounds | invalid record | 2 | executed |
| 13 | **run oracles** | declared mechanical checks | oracle unavailable → assessment `INCOMPLETE` | 1 | **designed** |
| 14 | **reconstruct view** | pure function of `(base, matter)` | reconstruction failure | 2 | **designed** |
| 15 | assess | verdict decodes to the requested schema; unreviewable matter refused before renting | seat malfunction | 2 | executed |
| 16 | seal β | bounds | invalid record | 2 | executed |
| 17 | γ closes | record composed; one digest taken | — | — | executed |
| 18 | V validates | record validity, digest recomputes, contract matches, required evidence present with producer authority, review passes | verdict carries typed failures | — | executed |
| 19 | δ decides | mechanical from receipt + verdict | — | — | executed |
| 20 | lift | status from (verdict, decision, mode) | inconsistent | 2 | executed |
| 21 | self-verify | `VerifyClosure` against **this invocation's** contract and metadata, never the closure's own | closure rejected, empty stdout | 2 | executed |
| 22 | encode + exit | — | — | 0 accepted · 1 non-accepted terminal · 3 simulated | executed |

**The gap this inventory makes visible.** Step 6 is the one refusal path that
currently emits nothing. A taskless CDS cell exits non-zero with no closure, so
the refusal is invisible to every downstream consumer — no receipt, no reason,
no custody. Every other terminal state in this table produces an artifact. That
asymmetry should be an acceptance criterion for promotion, not a footnote.

---

## Part II — Narrative, decisions, and open questions

The sections below explain decisions the inventories cannot express, and hold
the questions still open. Where a narrative sentence and an inventory row
disagree, the row is what a reviewer checks.

---

## 1. Working thesis

### 1.1. Conceptual definition

**Working position:**

> A cell is a bounded, cognition-bearing transformation that applies a declared
> methodology to a contract and a subject, produces candidate matter, assesses it
> independently, and emits a validated receipt before any result crosses its
> boundary.

The cell is not defined by Claude, Codex, or another model. Its identity is its
contract type, subject type, output-matter type, methodology, and authority
boundary. A model is a replaceable execution engine inside a declared component.

### 1.2. Cognitive-only v0

**Working position:** every deployed CNOS cell in v0 contains at least one cognitive
station. A purely mechanical transformation should normally be an ordinary
function, workflow, oracle, or CM property provider; it does not need cell ceremony.

Do not permanently encode this product-scope decision as an ontological prohibition
in the kernel. A future mechanical-only plan may earn cell treatment if uniform
composition and receipts provide real leverage. That is not a v0 requirement.

### 1.3. Four distinct artifacts before a run

The design must distinguish:

```text
Issue             why/what must change; scope and acceptance
Design            how it will change; invariants and impact graph
SubjectSnapshot   exact state being operated on
CellDefinition    reusable methodology and component construction
```

For CDS, implementation must not start without both an issue and a design. A small
change may carry a compact inline design; the logical artifact remains distinct so
"what" cannot silently choose "how".

### 1.4. Three output planes plus the receipt

```text
matter       the proposed semantic result
evidence     runtime-observed witnesses supporting or refuting claims
telemetry    non-authoritative progress/debug/usage events
receipt      immutable closure envelope binding inputs, matter, assessment,
             evidence references, provenance, verdict, and decision
```

Telemetry is never evidence merely because it was logged. A named collector or
oracle must deliberately promote an observation into a digested evidence artifact.

---

## 2. Does the CDS cell now flatten?

### 2.1. Working answer

Yes, substantially—but not into one model call.

The cell-specific, configurable **semantic** spine is three cognitive
stations — this is unchanged, and the subject adapter is deliberately not a
fourth:

```text
admit contract + subject
  → produce candidate matter
  → assess candidate matter
```

The **execution** flow shows where the substrate serves those stations. The
subject adapter is a declared runtime-substrate component, not a station: it
makes no decision and rents no cognition, but two stations depend on it and
therefore it can belong to neither.

```text
admit
  → materialize(base) ──────────→ produce
                                     │ measure(workspace, base) → matter
                                     ▼
                   reconstruct(base, matter) → evaluation view → assess
```

The remaining tail is fixed mechanics:

```text
γ closes → V validates → δ decides
```

**Why the adapter had to be named.** In the executed code, materialization
lives inside the producing fill and the workspace is released when production
returns (`cdspatch.go:193-197`, `defer release()`). Assessment runs afterwards,
so a reconstructed candidate workspace is not merely unbuilt — the subject no
longer exists when assessment needs it. The measurement defect has the same
origin: measuring "what changed relative to what" was a private detail of the
producing seat rather than the whole job of a named component. It measured
against `HEAD` from Case 2 until `4e8fe9c8`, through a merged case and nine
review rounds, because no component's specification made the question its own.

For the bootstrap:

- **admission** is a CDS-owned peer component combining structural validation with
  a receipted, explicitly attested cognitive issue/design review;
- the **subject adapter** is a declared runtime-substrate component — it makes no
  decision and rents no cognition, but production and assessment both depend on
  it, so it belongs to neither;
- **production** is a full engineering cognition component over a disposable,
  pinned subject;
- **assessment** is a fixed cognitive review component using the admitted contract,
  runtime-derived evaluation view, and declared skills;
- **gamma** is mechanical receipt construction;
- **V** is mechanical receipt validation;
- **delta** is boundary policy/authority, not cognition.

In the future, admission and assessment both become executions of compiled Coh
methodologies. Their position stays; their bootstrap implementation changes.

### 2.2. Important correction

Bootstrap beta is **not mechanical** merely because its orchestration and output
shape are fixed. It is rented cognition executing a fixed falsification procedure.
The eventual Coh evaluator is mechanically orchestrated, but individual property
providers may still be cognitive.

### 2.3. Where input checking belongs

The generic runner cannot know what makes a CDS issue, writing brief, research
question, or dataset admissible. Admission must therefore be a first-class
component of the cell definition.

The runner only performs the generic operation:

```text
construct the declared admission component
→ give it the run's contract + subject envelope
→ continue only on an admitted result
```

For CDS v0, a `cds.work-contract-admission` constructor owns:

1. exact structural validation of issue + design;
2. the structural relation between issue, design, and subject;
3. a cognitive semantic check, using issue/design skills, when structure alone
   cannot establish executability;
4. an explicit `admitted | rejected | incomplete | failed` result;
5. a digest binding the admitted issue/design bytes.

Admission is a PEER component of the cell definition, not a constructor nested
under `input`: an untrusted run-input envelope enters it, and success produces an
`AdmittedContract` plus an admission receipt that production then consumes.

**Target state:** malformed or unreviewable input emits an admission receipt and
stops before production. It is never silently discarded.

**Current state, stated exactly:** this does not hold. Both seats call
`cdsissue.Admit` before renting cognition, so an inadmissible issue does stop the
run — but it stops it with a non-zero exit and **no closure at all**
(`cellrun/run.go` returns 2 on a seat error, before any receipt exists). The
refusal is therefore invisible to every downstream consumer: no receipt, no
reason, no custody. Every other terminal state in Inventory 7 produces an
artifact. Closing that asymmetry is an acceptance criterion for promotion.

**Bootstrap cognitive admission is an attestation, not a verification.** It is
rented cognition, it is not itself reviewed, and it must be labelled
`attested / unverified` in its receipt until an independent review or Coh
replaces it. γ, V and δ go unreviewed because they are mechanical and
re-derivable; cognitive admission is neither, and must not borrow their standing.

---

## 3. JSON question — source, IR, and run request

### 3.1. Working answer

We should now be able to express a CDS cell fully in JSON, but we need to stop
calling every JSON stage "the IR."

```text
CellSourceJSON       human/compiler-authored declaration; may contain holes
NormalizedCellIR     holes filled; exact keys/types; no unresolved symbols
CompiledCellPlan     skills loaded/digested; providers validated; capabilities and
                     execution plan bound by the composition root
RunRequest           admitted issue + design + pinned subject for one episode
```

The runner may fill holes while bootstrapping. The executable IR is the resolved,
canonical result, not the hole-bearing source.

### 3.2. Why canonical JSON remains the likely v0 IR

JSON is currently the strongest fit because it is:

- language-neutral between a future Coh/cell compiler and the Go runtime;
- directly validated by CUE as an independent oracle;
- content-addressable under a canonical serialization;
- inspectable and diffable during bootstrap;
- closed enough to reject unknown fields and arbitrary executable configuration;
- already the proven input/output shape of both CNOS cells and TSC CM work.

The final design must still compare:

| Alternative | Main cost |
|---|---|
| Execute CUE directly | conflates constraints/defaulting with operational semantics |
| YAML | duplicate keys, coercion, and weak canonical bytes |
| Protobuf | code-generation and poor review/authoring surface; useful later as transport |
| CBOR | compact canonical transport, poor human/debug surface |
| Compiler-private AST | not portable, independently validatable, or durable |

### 3.3. Candidate authored JSON

This candidate is a completeness probe, not a schema decision:

```json
{
  "format": "cnos.cell.source/0.1",
  "id": "cnos.cds.implement",
  "version": "0.1",
  "protocol": "cnos.cdd.cds.receipt.v1",
  "parameters": {
    "language": {
      "kind": "skill",
      "required": true,
      "domain": [
        "cnos.eng:eng/go",
        "cnos.eng:eng/ocaml",
        "cnos.eng:eng/typescript"
      ]
    },
    "style": {
      "kind": "skill",
      "required": false,
      "default": "cnos.eng:eng/write-functional"
    },
    "provider": {
      "kind": "value",
      "required": true,
      "domain": ["claude-cli", "fake"]
    },
    "model": {
      "kind": "value",
      "required": false,
      "default": ""
    }
  },
  "input": {
    "contract": {
      "kind": "cnos.cds.implementation-contract/0.1"
    },
    "subject": {
      "kind": "git.snapshot/0.1"
    }
  },
  "admit": {
    "fill": "cds.work-contract-admission",
    "cognition": {
      "provider": "$provider",
      "model": "$model"
    },
    "skills": [
      "cnos.cdd:cdd/issue",
      "cnos.cdd:cdd/design"
    ]
  },
  "produce": {
    "fill": "cds.patch",
    "cognition": {
      "provider": "$provider",
      "model": "$model"
    }
  },
  "methodology": {
    "skills": [
      "cnos.eng:eng/code",
      "cnos.eng:eng/test",
      "$language",
      "$style"
    ]
  },
  "assess": {
    "fill": "cds.bootstrap-falsifier",
    "cognition": {
      "provider": "$provider",
      "model": "$model"
    }
  },
  "output": {
    "matter": {
      "kind": "git.patch/0.1"
    },
    "required_evidence": [
      {
        "id": "diff",
        "kind": "diff",
        "producer": "alpha"
      },
      {
        "id": "assessment",
        "kind": "property-assessment",
        "producer": "beta"
      }
    ]
  }
}
```

### 3.4. What the candidate deliberately omits

- **Issue, design, repo, and base SHA are not holes.** They are run-specific
  contract/subject values carried by `RunRequest`.
- **No `producer: "runtime"`.** `cellkernel.Role` is `alpha | beta`
  (`kernel.go:76-77`) and `validateRecord` rejects anything else
  (`kernel.go:617`). Station ownership and observation origin are different
  questions and are not being conflated in this round: runtime-measured
  evidence stays positioned on the producing station's side under the current
  kernel, and the distinction is deferred to a later design round with its own
  evidence.
- **No per-station skills list.** One `methodology` bundle is declared once;
  production receives its constructive projection and assessment its
  adversarial one. A `produce.skills` or `assess.skills` key would be two
  normative sources of obligations wearing one name.
- **Gamma, V, and delta are not configurable components here.** The declared
  protocol selects the fixed mechanical closure/validation/boundary contracts.
- **No arbitrary binary, argv, environment, command, or capability list appears.**
  A fill declares its required capability; runtime policy decides whether and how
  it can be provided.
- **No `bindings.alpha` or component-reference plane exists.** Each component is
  declared once with its constructor properties inline.
- **`subject.kind` names an adapter, not a runner capability.** It selects a
  declared subject adapter at the composition root, exactly as `fill` selects a
  constructor. The generic runner stays unaware of git.
- **`assess.fill` is provisional.** `cds.bootstrap-falsifier` is a placeholder
  for whatever §3.7 settles, and is deliberately not the executable-looking
  `cdd.review`: a generic reviewer must not contain a CDS predicate, and the
  current implementation does — `cdsreview.go:188` gates matter on
  `\ndiff --git `. Either the matter type carries its own admissibility
  predicate or the cell definition supplies one; that is an open decision, not
  a naming choice.
- **No task-specific skills hole appears yet.** The design must decide whether
  admitted issue/design contracts may declare additional canonical skills and how
  the fill merges them without creating an unbounded prompt-extension escape.

### 3.5. Which holes have earned their place?

**Working position:**

| Value | Classification | Why |
|---|---|---|
| language skill | semantic configuration hole | changes applicable engineering properties |
| style skill | project-policy hole with default | not universal to every CDS use |
| provider | execution-configuration hole | replaceable engine, selected before construction |
| model selector | execution-configuration hole | must be explicit in the resolved plan/receipt |
| issue + design | runtime contract, not holes | they define this episode's work |
| repo + base SHA | runtime subject, not holes | they define the exact state operated on |
| timeout/tool policy | runtime/fill policy, not holes yet | making them arbitrary weakens the boundary without demonstrated need |
| issue-specific skills | open decision | should be admitted contract data if allowed, not arbitrary CLI injection |
| methodology bundle | declared once, not a hole | one normative source; `$language` and `$style` are holes INSIDE it |
| subject adapter | selected by `subject.kind`, not a hole | the substrate serving two stations, not per-episode configuration |

The same provider/model holes may populate admission, alpha, and beta while each
station still receives a fresh, independent invocation. We should not add
station-specific provider/model parameters until a real need appears.

### 3.6. Which CDS skills belong in alpha?

The known suitable baseline is currently in `cnos.eng`:

```text
eng/code
eng/test
eng/<language>
eng/write-functional (when project policy selects it)
```

The current root `cnos.cds:cds` skill describes lifecycle/orchestration concerns
that a workspace-edit alpha may not have authority or inputs to perform. Do not
inject it ceremonially. The design must decide whether CDS has a distinct,
capability-matched producer obligation deserving a narrow `cds/implement` skill,
or whether the fill + issue/design contract already supply all CDS-specific
semantics.

### 3.7. Bootstrap beta

The bootstrap intent is sound: beta is a fixed cognitive review procedure so CDS
can implement TSC, after which beta becomes Coh execution.

However, the current `cdd/review` skill cannot simply be assumed suitable. It
contains branch, issue/PR, `.cdd`, CI, and write/merge expectations. A no-workspace
answering seat cannot obey that full contract. Three options were on the table:

1. define a narrow bootstrap review skill over `(issue, design, subject, matter)`;
2. define a mechanically selected projection of the one methodology bundle;
3. give the review component the complete capability/input contract the existing
   skill requires.

**Decided: option 2.** Option 1 creates a second normative source of criteria —
the thing the single-bundle decision forbids and the thing Coh exists to end.
Option 3 grants branch, PR, `.cdd` and CI authority the review does not need,
which makes independence harder to argue rather than easier.

The projection is stated as a PROPERTY rather than a procedure, so it survives
Coh replacing the mechanism: the two views must be derived from one bundle such
that no obligation appears in one and not the other. That is checkable today and
is the missing agreement mechanism in Inventory 5's methodology row.

The fixed review algorithm should be:

```text
for every declared property / acceptance obligation:
  attempt to falsify it against a runtime-derived evaluation view
  emit pass | finding | unverified | failed, with citations
```

It must not guess when the supplied view cannot decide a property.

### 3.8. Resolved IR obligations

After parameter filling and construction, the resolved/compiled representation
must bind at least:

- exact cell source digest and schema version;
- concrete provider adapter id and requested model selector for every cognitive
  station;
- ordered canonical skill refs plus exact content digests;
- fill implementation/version or policy digest;
- contract kind, subject kind, matter kind, and protocol;
- declared admission and assessment result vocabularies;
- capability/resource requirements as runtime-owned policy identifiers;
- no unresolved holes, symbolic component refs, arbitrary commands, or ambient
  configuration dependencies.

Whether this is one `NormalizedCellIR` or a normalized IR followed by a distinct
`CompiledCellPlan` is a design decision. TSC already argues for the two-stage split;
CNOS should not collapse it without evidence.

---

## 4. Architecture sketch the design must refine

```text
AUTHORING / CONFIGURATION
┌──────────┐  ┌──────────┐  ┌────────────────────────┐
│  Issue   │  │  Design  │  │ Cell source            │
│ what/why │  │ how/impact│ │ admission/α/assessment │
└────┬─────┘  └────┬─────┘  └───────────┬────────────┘
     └──────┬──────┘                    │
            │                    compile → CUE vet
            │                    → resolve/link
            └─────────────┬─────────────┘
                          ▼
                   Compiled Cell Plan
                          │
RUNTIME                   ▼
Subject snapshot → RunRequest → admission
                                │
                                ▼
                         α produces matter
                                │
                runtime reconstructs candidate view
                                │
                                ▼
                    β / Coh assesses properties
                                │
                     runtime-owned evidence
                                │
                                ▼
                        γ closes receipt
                                │
                        V validates receipt
                                │
                        δ decides boundary
                                │
                          effect adapter

OBSERVATION
runtime events ──────────→ telemetry sink       non-authoritative
closed receipts ─────────→ ε / learning cells   cross-episode
```

---

## 5. Proposed canonical design-document outline

The final design should answer the following questions in this order.

It must also include seven falsifiable inventories rather than leaving the answers only
in narrative prose:

1. component inventory — responsibility, ownership, and forbidden knowledge;
2. value catalogue — producer, consumers, schema surfaces, and scope-lift binding;
3. authority table — enforces, declares, and claims nothing, including barrier location;
4. seat specifications — decision, inputs, sufficiency argument, outcomes, independence;
5. rule ownership — normative owner, mirror, and executable agreement mechanism;
6. gate specification — artifact provenance, what green proves, and what it does not;
7. end-to-end flow — every gate in order, with typed failure and exit behavior.

The narrative sections below explain and justify those inventories; they do not replace
them.

### 5.1. Document contract and authority

- Is this the canonical current-state realization design?
- Which existing documents remain normative companions?
- Which become rationale, migration notes, or superseded proposals?
- What does CCNF own, what does TSC own, and what does CNOS cell realization own?

### 5.2. Problem, goals, non-goals, and challenged assumptions

- Which repeated failures prove the current boundary is incomplete?
- What future work should become cheaper after this design?
- What does v0 explicitly refuse to solve?

### 5.3. Conceptual definition and invariants

- What distinguishes a cell from a function, workflow, CM, provider, or agent?
- Is cognition required in v0 and/or in the type?
- What determines cell identity and episode identity?
- Which invariants must hold regardless of realization?

### 5.4. Vocabulary and type model

- What are `CellSource`, `NormalizedCellIR`, `CompiledCellPlan`, `RunRequest`,
  `Contract`, `SubjectSnapshot`, `Matter`, `Assessment`, `Evidence`, `Telemetry`,
  `Receipt`, `Verdict`, `BoundaryDecision`, and `ContractFinding`?
- Which are authored, derived, immutable, content-addressed, or runtime-local?

### 5.5. System boundary and complete component diagram

- Which components are inside the episode runtime?
- Are V and delta part of the cell wall or the enclosing parent?
- Where do effectors, storage, transport, worktrees, credentials, and model
  processes live?
- What is each component's one reason to change?

### 5.6. Definition, construction, and dependency injection

- Which properties are configured once?
- Which values are holes, run inputs, or runtime policy?
- How is every component constructed exactly once?
- What lifecycle/scope does each component have?
- Which Spring ideas are adopted and which are rejected?

Working Spring answer: explicit nested component definitions, constructor injection,
one composition root, per-episode scope, and validation before instantiation. Reject
reflection, scanning, autowiring, mutable singletons, service location, and arbitrary
runtime class/command names.

### 5.7. Source language, normalized IR, compiled plan, and CUE

- Is cell authoring a Coh `cell` form or a small composition form containing a CM?
- Which obligations are compiler checks versus CUE checks versus linker checks?
- Why JSON, and which canonicalization contract applies?
- How are schema/IR versions migrated?

### 5.8. Contract, subject, issue, and design

- What belongs in the generic contract?
- What makes an implementation issue executable?
- What must design add beyond the issue?
- May small designs be inline while remaining distinct?
- How are issue and design digests bound to the run?
- How do ACs, scope, non-goals, invariants, and the impact graph differ?

### 5.9. Admission

- What can be rejected structurally before cognition?
- What requires semantic assessment?
- Is semantic admission a preflight child cell or a station in the target cell?
- What exact artifact proves alpha received the admitted bytes?
- What happens on `INCOMPLETE`?

### 5.10. Episode execution semantics

- What exactly does alpha receive and what capabilities must it have?
- How is the subject contained and how is candidate state measured?
- What exactly does beta/Coh receive?
- Which stages may run in parallel?
- When do oracles run?
- What does cancellation mean at each stage?

### 5.11. Methodology, skills, properties, and Coh

- What is the one normative source of assessable properties?
- How are constructive alpha guidance and adversarial assessment projected from it?
- Which instructions are advice rather than assessable obligations?
- How are mechanical and cognitive property providers composed?
- How are refusal and uncertainty represented?

### 5.12. Outputs, evidence, receipt, decision, and custody

- What returns directly and what is stored by reference?
- Does a receipt embed matter or bind its digest?
- Who owns evidence bytes?
- Which artifacts ascend to the parent cell?
- How do retries and superseding attempts preserve history?

### 5.13. Telemetry and learning

- Which progress/tool/usage events exist?
- What is retained or redacted?
- When can telemetry become evidence?
- What may epsilon consume?
- Does losing telemetry ever affect the verdict? Recommended: no.

### 5.14. Failure, discovery, repair, and re-contracting

- How do malformed input, provider failure, matter defect, unverified property,
  contract defect, resource exhaustion, and cancellation differ?
- What happens when alpha discovers a necessary surface outside the design's impact
  graph?
- Can a cell amend its own contract? Recommended: no.
- When is repair valid and when is a new issue/design digest required?

### 5.15. Composition and recursion

- How do child contracts descend and receipts ascend?
- How does alpha propose children without invoking the runtime itself?
- How do `let!` and `and!` lower into a dataflow graph?
- What belongs to the main loop rather than one cell?

### 5.16. Authority, capability, security, and resources

- What is declared, available, approved, and actually contained?
- Which substrate gives engineering alpha full useful capability safely?
- Who chooses provider/model?
- Which credentials may enter execution but never the receipt?
- Which boundary actions require an operator?

### 5.17. Reliability and operational semantics

- What is deterministic and what is not?
- What does reproducibility mean for cognition?
- How are timeouts, retries, IDs, cancellation, concurrency, and canonical event
  order defined?
- Does retry mint a new attempt? Recommended: yes.

### 5.18. Worked realizations

- CDS implementation: admitted issue+design + pinned repo → patch → assessment.
- Writing: admitted brief/style + source text → document.
- Planning: goal + state snapshot → plan/relation matter.
- One mechanical and one cognitive property provider inside a cognitive cell.

### 5.19. Alternatives and decision record

Compare at minimum:

- cognitive-only versus mechanically realizable cells;
- issue-only versus issue+design;
- blind beta versus runtime-reconstructed evaluation view;
- skills-as-prompts versus properties-as-methodology;
- Coh-as-whole-cell versus Coh-as-methodology component;
- JSON versus CUE/YAML/Protobuf/CBOR/private AST;
- nested constructor injection versus component refs/service locator;
- telemetry inside receipts versus separate stream;
- automatic scope widening versus explicit re-contracting.

### 5.20. Migration, impact graph, and conformance

- Which Case-2/Case-3 code survives, changes, or is held?
- What is the smallest bootstrap beta compatible with future Coh?
- When does work switch to TSC?
- Which existing docs are updated or demoted?
- Which fixtures prove every boundary can fail for its own reason?

---

## 6. Contract-discovery rule

Acceptance criteria prove selected outcomes. They do not enumerate every effect a
valid implementation may have.

The design's impact graph, invariants, scope, and non-goals establish the execution
boundary. During production:

| Discovery | Required disposition |
|---|---|
| Necessary change already named by impact graph/invariants | continue and bind evidence |
| Necessary dependency not represented in design | emit `ContractFinding`; stop `needs_recontract` |
| Issue/design contradiction | emit `ContractFinding`; stop `needs_recontract` |
| Optional improvement outside contract | do not implement; emit follow-up candidate |
| Required oracle unavailable | assessment `INCOMPLETE`; no acceptance |
| Forbidden/out-of-scope effect | reject candidate matter |

Partial matter may be retained as a non-transmissible proposal. The cell cannot
amend the issue/design that authorizes itself. An enclosing planning/design cell or
operator produces a new contract digest, then a new episode runs.

---

## 7. Impact graph for the eventual design

### Upstream sources

- CCNF kernel doctrine and receipt-validation boundary;
- `cdd/issue`, `cdd/issue/contract`, `cdd/issue/proof`, and `cdd/design`;
- `cnos.core/design` and `cnos.eng/evolve`;
- TSC `CMSource → NormalizedCMIR → CompiledCM → RunRequest → MeasurementReceipt`;
- Case-2 executable evidence and Sigma's Case-3 r0 findings.

### Downstream consumers

- `schemas/cdd/spec.cue` and `schemas/cds/spec.cue`;
- `cellfill`, `cellspec`, `cellrun`, `cellkernel`, and composition root;
- CDS admission, patch, and bootstrap review fills;
- skill loading and resolved invocation metadata;
- workspace/evidence/telemetry/provider adapters;
- future `.cell`/Coh compiler and TSC runtime;
- CLI and future GitHub invocation adapter;
- planning/cohering/main-loop work;
- receipt custody, epsilon, and memory/learning projections.

### Existing documents requiring an authority disposition

- `COHERENCE-CELL-NORMAL-FORM.md` — retain as kernel authority;
- `CELL-OF-CELLS.md` — retain as recursive system model/rationale;
- `DUMB-MODELS-SMART-CELLS.md` — retain as product/authority rationale;
- `CELL-RUNTIME.md` — reconcile or supersede its proposed realization claims;
- `CDS-CELL-MIGRATION.md` — demote to migration/history after canonical design;
- `CELL-RUNNER-CASES.md` — retain as executable-case record, not architecture source.

---

## 8. Leverage and negative leverage

### Positive leverage

- Stops implementation from choosing architecture one local defect at a time.
- Makes issue and design explicit prerequisites to coding.
- Gives every future cell one construction, admission, evidence, and failure model.
- Lets bootstrap skills be replaced by Coh without changing episode boundaries.
- Makes planning, working, and cohering cells composable through the same artifacts.
- Gives Sigma and Pi one review authority instead of several partially overlapping
  essays and branch-local comments.

### Negative leverage / cost

- Pauses visible Case-3 implementation progress.
- Requires a migration from current CellSpec/closure shapes.
- Adds an explicit design artifact to CDS contracting.
- Requires deciding versioned source/IR/compiled/run artifact boundaries.
- May reveal that some current skills mix incompatible roles and need projection or
  decomposition.

The cost is justified only if the final design deletes repeated local reasoning and
becomes the contract future implementation reviews against.

---

## 9. Non-goals of this outline

- No Case-3 code changes.
- No final JSON/CUE schema.
- No new issue or PR.
- No Coh runtime implementation.
- No generic provider router.
- No repair loop or main-loop implementation.
- No mechanical-only cell implementation.
- No claim that the candidate JSON is accepted by current code.
- No claim that the current `cdd/review` skill is usable unchanged.

---

## 10. Questions for Sigma's design review

Sigma should review the outline and candidate JSON against the live Case-2/Case-3
implementation evidence, then answer each item with `agree`, `change`, or `open`:

1. Is the three-component configurable spine `admit → produce → assess` complete?
2. Is the fixed tail `γ → V → δ` correctly omitted from per-cell configuration?
3. Are issue + design + subject + cell definition the complete pre-run artifact set?
4. Does the Source JSON contain every CDS-specific decision without leaking runner,
   GitHub, arbitrary command, or service-locator semantics?
5. Which holes have actually earned their place? Are any missing?
6. Does `cnos.cds` need a narrow producer skill, or is CDS semantics already supplied
   by the fill + issue/design contract?
7. Can the existing review skill be projected safely, or does bootstrap require a
   new narrow review skill?
8. What must beta's runtime-derived evaluation view contain while preserving the
   CCNF `(contract, matter)` boundary?
9. Should admission be a preflight child cell or an inline station of the run?
10. Should Coh express the whole cell, or remain a measure-only component embedded
    in a cell composition form?
11. Is `NormalizedCellIR` distinct from `CompiledCellPlan`?
12. Which conclusions are contradicted by actual executed evidence?

The requested output is a design response or patch to this outline. It is explicitly
not authorization to continue Case-3 implementation.

---

## 11. Acceptance criteria for promoting outline → design

- [ ] The named incoherence and all observed failure modes are evidenced.
- [ ] Every component appears once in the diagram and responsibility table.
- [ ] Contract, subject, definition, source, normalized IR, compiled plan, run,
      matter, assessment, evidence, telemetry, receipt, verdict, decision, and effect
      are distinct typed artifacts or intentionally unified with a stated proof.
- [ ] The CDS cell is fully expressible as valid candidate JSON without hidden
      semantics in the generic runner.
- [ ] Issue and design admission has structural and semantic paths with typed refusal.
- [ ] Alpha has truthful engineering capability under a named containment boundary.
- [ ] Beta can verify every asked property or return `INCOMPLETE`; no hidden-context
      recovery is required.
- [ ] Gamma is demonstrably mechanical; V and delta authority remain separate.
- [ ] Matter, evidence, telemetry, and receipt channels cannot be confused.
- [ ] Newly discovered impact has a typed re-contract path and cannot self-authorize.
- [ ] Coh/TSC migration preserves the same positions and artifact contracts.
- [ ] Alternatives include at least three structurally different composition/IR moves.
- [ ] The impact graph covers every existing doctrine, schema, runtime, skill, and
      invocation consumer.
- [ ] Every final implementation AC names a specific artifact and an oracle, including
      a negative/mutation witness.
- [ ] Every refusal path emits a closure. A refusal that produces no artifact is
      invisible to custody, repair, and learning.
- [ ] Every oracle states how it obtains the artifact it measures. An oracle with
      unstated provenance reports on an unknown object.
