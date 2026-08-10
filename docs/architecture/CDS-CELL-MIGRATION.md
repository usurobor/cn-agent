# CDS Cell Migration — from prose wake to a runnable data cell

Status: plan · Scope: `cnos.cdd` (substrate) + `cnos.cds` (program) · Focus:
**run a CDS cell from the CLI.** GitHub-issue mechanics (selector, claim, FSM
label transitions) are explicitly out of scope for this pass — they stay in the
existing wake and are wired back later.

## Goal

Today `cnos.cds` is expressed as a prose **wake** (`orchestrators/cds-dispatch/
SKILL.md`): a selector claims an issue and hands it to δ, which routes γ/α/β by
reading skill markdown. The cell only exists as text.

Migrate CDS so the cell exists as **data the kernel runs**:

```
main.cell  --compile-->  spec.json  --cue vet #CellSpec-->  cellkernel.Run
 (surface)   (cnos.cdd)    (IR)          (oracle)              (Go kernel)
```

This is the exact mirror of the pipeline we already ship on the *output* side
(γ emits a receipt → `cue vet #CDSReceipt` → V). We are adding the *input* side.

## Why — the reasoning behind every choice

This plan is the settled end of a long design dialogue. The reasoning matters
as much as the steps, because each decision closed a specific fork.

1. **The kernel is pure mechanics; it cannot hold a resolver.** `cellkernel.Run`
   executes the five-step closure at one scope and knows *nothing* about the
   work being done (which language, which skills, which repo). Therefore the
   cell spec must be **constructed before/around the run** — never resolved
   inside the kernel. This is the load-bearing separation the whole design
   protects.

2. **The spec is static data "all the way down" — but only one level at a
   time.** For a mono-lingual repo where the language falls out mechanically,
   the whole tree can be compiled up front (Case A). When a decision needs
   *judgment* (a multilingual repo where the issue demands Rust vs Go), the
   child spec does not exist until the parent's α runs — the tree unfolds one
   static level at a time (Case B). Each cell is static *when it runs*; deeper
   specs are produced at runtime by α cognition emitting data.

3. **Composition is one reentrant mechanism, not a new construct — but α
   *proposes*, the runtime *executes*.** A parent decomposes by having its α
   emit *child contracts* (data); the **runtime/Drive** runs each child episode;
   the accepted child receipts become the parent's matter; parent-β reviews the
   composition against the parent contract. This is the settled industry answer
   — Railway-Oriented Programming (F# `Result` composition), the GoF Composite +
   Decorator patterns (uniform interface so a child stands in for the parent's
   input), the interpreter-over-data / Free-monad pattern (recursion lives in
   the *data + a recursive walker*, not a grammar), and Unix pipes (compose only
   through the typed payload — for us, the receipt).

   **Correction from Pi β (`msg-cn-pi-cnos-cell-runner-cases-review-31`, C1):**
   an earlier framing made parent-α *call the kernel on the child* directly
   (`α = extract ∘ Run ∘ resolve`, the `.NET Rx` `SelectMany`). That fuses
   *decomposition judgment* with *runtime execution* inside a cognition seat. A
   composite-α adapter may be a later convenience, but it must not be the
   normative boundary: **the runtime executes; α only proposes.** The surface
   `let!` (below) is sugar that *lowers to* "α emits child contract → runtime
   runs it → receipt → matter" — not to α invoking `Run` itself. This keeps
   child execution mechanical and auditable (route receipts, custody) rather
   than buried in α.

4. **The static-vs-dynamic "spec expression" dilemma was never a language
   problem.** We weighed compiling everything to a static resolved spec against
   a dynamic expression language carrying `resolve` at runtime. TSC's
   `cm-language` (the F#-shaped `coh` surface) settled it by example: it is a
   *compiler to a CUE-validated data IR* with **no runtime** — it names
   providers but executes nothing. The dynamism we feared (`α = resolve |>
   coding_cell`) is not a language operator; it is α cognition *emitting data*
   that the kernel then loops over. So: **spec is data, `resolve` is a library
   function, recursion is the kernel's loop.** No runtime expression language.

5. **We adopt `cm-language`'s *shape and toolchain architecture*, not its
   *semantics*.** cm is measure-only by construction (`forbid compile, admit,
   authorize, repair, self_authorize`) — the categorical opposite of a cell
   whose α *produces* and whose δ *repair-dispatches*. So we do not adopt the
   language. We adopt its form: the compact F#-computation-expression surface
   (`NAME (params) -> ReturnType { … }`, `let!`/`and!` binders) compiled to a
   CUE-validated data IR with **CUE as the independent oracle**. It is the most
   compact legible way we have seen to express a bounded, typed, composable
   workflow; the noise of raw YAML (`mode: rented` everywhere) is exactly what
   it removes.

6. **No OCaml.** cm chose OCaml because TSC *is* an OCaml toolchain. Our kernel
   is Go; introducing OCaml means a second toolchain and a subprocess/FFI seam
   between compiler and kernel — pure friction. A ~150-line Go `participle`
   grammar (or Starlark to bootstrap) targets `cellkernel.Spec` natively.

7. **Compile to JSON + `cue vet` — because it mirrors what we already run.** The
   receipt path is already `γ emits → cue vet #CDSReceipt → V`. Adding the
   input-side mirror (`compile → cue vet #CellSpec → Run`) means both ends share
   one compile-vet-execute shape. That symmetry is why the remaining work is
   three small artifacts, not a new architecture.

8. **Parameters are Unix typed holes — because that is the seam that lets us
   start simple and grow.** Filling `language` from the CLI now, from a parent
   cell (mechanically or cognitively) later, changes only *who fills the hole* —
   never the cell. Unix already solved value→implementation resolution (`$PATH`)
   and required/optional inputs (positional vs flag); we borrow it wholesale.

9. **The four protocols already differ in only two things** (verified against
   `schemas/cds/receipt.cue` and `schemas/cdr/receipt.cue`: both are
   `cdd.#Receipt & { protocol_id, evidence_refs }`). A protocol *is*
   `(protocol_id, evidence-key set, skills)`. So one generic cell + a thin
   overlay expresses cdd/cds/cdr/cdw — the surface differs only in the header
   protocol and the two skill lines. This is what makes "CDS is one cell the
   runner runs" concrete rather than aspirational.

## Two settled decisions

1. **Entrypoint file:** `main.cell` (per package). Mirrors `main.go`; the file
   holds cell definitions, the CLI runs one by name (`cn cell run cds`).
2. **Parameters are Unix-shaped typed holes.** See "Parameters → skills" below.

## The surface (F#-inspired, same shape as TSC `coh`/cm — our own vocabulary)

The CDS cell, concretely, grounded in `schemas/cds/receipt.cue`:

```
cell cds (issue, repo, language: skill, style?: skill = functional)
    -> cnos.cdd.cds.receipt.v1 {
  produce with eng, $language, $style     # α: implement
  review  with $language, cds-review       # β: review
}                                          # γ/V/δ: mechanical kernel defaults
```

- The header's `-> cnos.cdd.cds.receipt.v1` selects the receipt schema; V
  dispatches on it and cue-vets against `#CDSReceipt`.
- The five per-cycle artifacts (`self_coherence, alpha_closeout, beta_review,
  beta_closeout, gamma_closeout` + `diff`) are **not authored** — they are the
  protocol's evidence contract; **γ binds them mechanically**.
- cdr / cdw are the *same body shape*, differing only in header protocol and
  the two skill lines (cdr overlay exists; cdw schema not yet shipped).

A cell whose α is a subcell (the parent/dispatch case, later) is one `let!`:

```
cell cds-dispatch (issue, repo, language: skill) -> cnos.cdd.cds.receipt.v1 {
  let! impl = cds (issue, repo, $language)   # α PROPOSES this child contract;
                                             # runtime runs it; receipt → matter
  review with dispatch-review                 # β reviews the composition/choice
}
```

`let!` is surface sugar: it lowers to "parent-α emits the `cds` child contract →
the runtime executes that child episode → the accepted child receipt becomes
the parent's matter." It does **not** compile to α calling the kernel itself
(Pi β C1). Composition and single-episode work therefore share one runtime, not
two cell algorithms.

## Parameters → skills (the resolution model)

A parameter is a **typed hole**, resolved exactly like a Unix command resolves
argv against `$PATH`:

| Concern | Unix analogue | Here |
|---|---|---|
| required vs optional | positional arg vs flag-with-default | `language: skill` (required) vs `style?: skill = functional` |
| value → implementation | `ls` → `/bin/ls` via `$PATH` | `"ocaml"` → the `ocaml` skill via the **skill path** |
| domain / typo check | — | CUE constrains `language: "go"|"ocaml"|"rust"`; bad value fails `vet` |
| splice into the seat | — | `$language` in `produce with` |
| pass down to a child | wrapper `exec cmd "$@"` | `let! impl = cds(issue, repo, $language)` |

Four layers, matching the rest of the architecture:

1. **Surface** declares holes (`language: skill`, `style?: skill = functional`)
   and splices them (`$language`).
2. **CUE `#CellSpec`** carries the holes and (via the cds overlay) constrains
   their domains — a missing *required* hole or an out-of-domain value fails
   vet. Mirror of how `#CDSReceipt` overlays `#Receipt` with `evidence_refs`.
3. **Runner** fills holes from the CLI (`--language go`), resolves each value to
   a concrete skill ref against the skill path, loads it into the seat.
4. **Kernel / seats** receive *loaded skills*. α/β never see the string
   `"ocaml"` — they get the resolved skill. (Seats consume skills, never
   resolve them — same top-down rule as everywhere.)

**Staging of who fills the hole (the cell never changes):**
- *Now (CLI bootstrap):* the runner fills every hole from CLI flags.
- *Later:* a parent cell fills `language` mechanically (`language_of(repo)`) or
  cognitively (`triage(issue)`). Only the *filler* changes; resolution
  (value→skill) and the `cds` cell body are untouched.

## Piece inventory

| Piece | Owner | Status |
|---|---|---|
| Receipt schema (output contract) | cnos.cdd + cds | ✅ shipped |
| CCNF kernel (`RunEpisode` → verifiable `Closure`) | src/go `internal/cellkernel` | ✅ shipped (PR #718; hardened through Pi β #31–#45) |
| Workspace-edit cognition (bounded, stateless provider adapters) | src/go `internal/cellcog` (`Coder` port; `ClaudeCLI`, `FakeCoder` over one process seam) | ✅ shipped — explicit model, typed argv, sealed permission mode, no arbitrary command from cell JSON. `codex-cli` HELD (see below). `Coder` is a workspace-EDIT port: directory in, no value out. It does not yet serve planning or research. |
| `#CellSpec` CUE schema (input contract) | cnos.cdd (`schemas/cdd/spec.cue`) | ✅ shipped |
| cds params-domain overlay (`#CDSCellSpec`) | cnos.cds (`schemas/cds/spec.cue`) | ✅ shipped (canonical diff-first evidence rule) |
| spec loader/binder (strict parse → kernel Spec) | src/go `internal/cellspec` | ✅ shipped |
| `cn cell run` (fill holes, run, emit closure) | src/go `internal/cellrun` (+ thin `cli` wrapper) | ✅ shipped (exits 0/1/2/3; closure self-verifies) |
| skill resolution + loading | src/go `internal/cellskill` (one installed root under the hub) | ✅ shipped — bodies are LOADED and injected; ordered refs + content digests recorded in the closure |
| fill-owned seat construction | src/go `internal/cellfill` (registry) + `internal/cdspatch` (`cds.patch`) + `internal/cellfills` (composition root) | ✅ shipped |
| matter substrate (diff at a base SHA) | src/go `internal/cellwork` (worktree adapter, outside the kernel) | ✅ shipped (G1) — the runtime measures the diff; a seat cannot claim one |
| typed findings on `BetaOutput` | src/go `internal/cellkernel` | ❌ G3 — the findings ARE the repair plan |
| `main.cell` compiler (surface → spec.json) | cnos.cdd | ❌ Phase 4 sugar (deliberately deferred) |
| `cnos.cds/main.cell` | cnos.cds | ❌ Phase 4 (JSON fixtures are the current authored form) |

## Phases (Pi's KISS ladder, `msg-…-review-31` D5; smallest-first)

**Phase K — correct the kernel (shipped).** The kernel corrections Pi's β #31
required (honest closure taxonomy, kernel-owned γ/V/δ with no injectable
certification, the evidence seam, fail-closed seats) landed and were then
superseded by the stronger FIDO refactor below — the shipped kernel is the
reference; the piece inventory above reflects it.

**Phase 0 — the input contract.** Write `schemas/cdd/spec.cue` (`#CellSpec`:
`{version, contract (producer-attributed required_evidence), protocol_id,
params, alpha:{fill, …}, beta:{fill, …}}`) and the cds overlay
(`#CDSCellSpec` pinning `protocol_id` + param domains). Hand-write a
`cds.spec.json` for a real issue; `cue vet` it. *Proves the data shape before
any Go or compiler.* (`budget` was removed as decorative per Pi #32 D5; it
returns only when it is actually enforced.)

**Phase 1 — CLI 0: one local episode, honest receipt (no cognition).**
`internal/cellspec` loader (strict parse → kernel `Spec`), and `cn cell run
--contract <path|-> [--param k=v]` that fills holes and calls `RunEpisode`.
Emits a **generic** `cnos.cellkernel.episode-closure.v0` with
`protocol_validated=false` — the declared `protocol_id` is provenance, never a
validated CDS claim. Each seat names its own `fill`: `cdd.stub` (a
non-authoritative **`simulated`** smoke run, exit 3), `cdd.bool` +
`cdd.bool-check` (a real mechanical episode with a genuine review predicate),
or `cds.patch` + `cdd.mechanical-unmet` (Case 2). **Zero GitHub/network.**

Implementation follows the operator-ratified **FIDO/functional doctrine**
(`msg-cn-pi-cnos-cell-runner-fido-functional-44`): no mutable shared episode
state; immutable seat scopes (`AlphaInput → Result<AlphaOutput>`, sealed before
crossing scope; β gets a fresh projection of the sealed α matter only —
artifact validation stays with V); **positional
ownership** (no seat-authored provenance — a required α artifact must sit under
`record.alpha`); **one** immutable `EpisodeRecord` with **one** scope-lift
digest; `VerifyClosure` as the single verification boundary (digest recomputes;
verdict/decision/status re-derive purely); typed failure routing (integrity →
`rejected`, never the α repair path); fail-closed identity; UTF-8 +
size-bounded artifacts. A CI job (`.github/workflows/cell-schema.yml`) vets the
CUE schemas + actual `cn cell run` output (accepted / needs-repair / simulated)
against a shared positive/negative corpus (Pi #31–#33 + PR-#718 β + #44).

**Phase 2 — skill resolution.** The `$PATH`-like resolver (`value → skill ref`)
+ required/optional/default hole logic; `cn cell run` errors on unfilled
required holes with a Unix-style usage line. Seats load the resolved skills
(cognition still stubbed). *Proves parameters map to real skills.*

**Phase 3 — rented α + CDS patch fill.** ◐ *First half shipped.* Stated
precisely: **`cds.patch` is the fill; cognition is one of its constructor
dependencies**, not a fill itself and not a new architecture.
`internal/cellcog` constructs bounded,
stateless provider adapters (claude-cli and a deterministic fake; codex-cli is
HELD, see below) with an explicit model, typed argv and a sealed permission
mode — a cell cannot smuggle flags into one, nor inherit edit authority from
the host.
`internal/cellskill` resolves canonical installed refs and **loads** the
bodies, recording ordered refs and content digests in the closure, because
naming a skill is not loading it. `internal/cdspatch` composes those with the
worktree substrate into one provider-neutral patch alpha; the generic runner
learns none of it.

**G1 shipped with it.** The runtime **measures** the change: `cellwork` cuts a
disposable worktree at a commit pinned during construction, and computes the
diff itself. A seat that claims a sweeping change and wrote nothing produces
no diff, and an episode with no diff cannot satisfy a contract requiring one —
false completion (the #514/#516 scar) is unrepresentable rather than a review
failure to catch later. Bounds are applied as output streams, not after it is
buffered.

*Still open before CDS is complete:* **Case 3** (an independent rented β
reviewing that diff — the first real judgement in the loop, and a one-field
`beta.fill` change under this boundary) and **G3** typed findings.
`cn cds run --issue N` and the escalation predicate (Pi Q2) follow those.

**Phase 4 — the surface + compiler (sugar, later).** `cn cell compile main.cell
→ spec.json` (Go `participle`, ~150 lines; CUE stays the independent oracle) +
author `cnos.cds/main.cell`. Pure sugar over Phase-1 JSON — **optional until it
earns its keep**, and explicitly *not* on the frozen S4/S6 critical path.

**Horizon (out of this pass):** `Drive` bounded-repair loop (Case 4); the
`cds-dispatch` parent (`let!` = α-proposes-child / runtime-executes, Case 5);
rewiring the GitHub wake to a thin adapter that fills holes + lands artifacts.

## HELD — GitHub Actions provisioning (captured, not implemented)

**Status: held.** Nothing below is built, and nothing in it belongs to
`cellrun`, `cds.patch`, `cellcog`, or cell JSON. It is written down now so the
later invocation adapter inherits decisions rather than rediscovering them.

1. **CLI installation and authentication belong to the runner/workflow image**,
   never to the runner code or a cell declaration.
2. **Provision exact pinned Claude and Codex CLI versions before an episode**,
   and fail before α if the selected executable is absent. Never
   opportunistically download during cognition — a cell that installs its own
   tools mid-episode has an unbounded, unreceipted dependency.
3. **Credentials stay secrets/environment supplied by the workflow** and never
   enter cell JSON or a receipt. Only the selected provider's credential should
   reach its child process.
4. **Model remains the explicit fill property.** The later execution receipt
   should record the resolved executable identity, observed CLI version (and
   artifact digest where available), provider-policy version, and requested
   model — never secrets, never the full environment.
5. **A workflow may install both pinned CLIs** so provider selection stays in
   the one alpha declaration; a prebuilt image is a later latency
   optimization, not a design change.
6. **The child environment must eventually be an explicit provider-specific
   allowlist**, not arbitrary inherited ambient configuration. Outer OS
   sandboxing remains a separate execution-substrate concern — this project
   claims no OS confinement.

*Empirical source notes (inspected commits, preserved so the later work does
not re-derive them):*

- Anthropic's action at `6b082c41935b4c8a3b8b0ef85ba4ba4d9eeb8975` is a
  composite action: it installs a pinned native Claude CLI during the job (or
  accepts a supplied executable path) and injects API/OAuth/WIF authentication
  from the workflow. Borrow the **provisioning boundary**, not the GitHub
  orchestration.
  <https://github.com/anthropics/claude-code-action/blob/6b082c41935b4c8a3b8b0ef85ba4ba4d9eeb8975/action.yml>
- OpenAI's action at `52fe01ec70a42f454c9d2ebd47598f9fd6893d56` is also
  composite: it installs npm CLI/proxy packages, starts a loopback Responses
  API proxy for the API-key path, then runs `codex exec` under declared
  permissions. **Correction (Pi #55 B1):** an earlier reading of this called
  the installed versions exact. They are not — the action's `codex-version`
  input defaults to blank, which tracks npm latest, and is exact only when the
  workflow supplies a version. Our decision is unchanged and now rests on the
  right reason: CNOS must pin the action commit *and* pass an explicit
  version, because the default is a moving target.
  <https://learn.chatgpt.com/docs/github-action> ·
  <https://github.com/openai/codex-action/blob/52fe01ec70a42f454c9d2ebd47598f9fd6893d56/action.yml>

## HELD — codex-cli as a provider (withdrawn from Case 2, not abandoned)

**Status: held.** `codex-cli` is absent from the admitted provider set in both
authorities — `cellcog.New` and `#Cognition` — and `provider_codex.go` is
deleted rather than left unreachable, because code whose comments claim an
isolation it does not deliver is worse than no code.

**Why it was withdrawn (Pi #55 D1).** A seat may carry only the fill's
ordered, digested skills; anything else is a second, unreceipted component
definition. The flags we had do not reach that far:

- `--ignore-user-config` suppresses only `$CODEX_HOME/config.toml`;
- `--ignore-rules` suppresses only user/project execpolicy `.rules`;
- global and project `AGENTS.md` guidance still loads;
- skills are still discovered from repository, user, admin and system
  locations.

The adapter's comments and tests claimed broader isolation than that, and
because `codex-cli` was admitted by both authorities this was a live
constructor boundary rather than a documentation error.

**What returning requires**, in order:

1. Typed suppression knobs in the installed Codex version, sealed into the
   argv — never env, config, or argv supplied by cell JSON.
2. A dedicated clean `CODEX_HOME` provided by the execution substrate.
3. A real run proving that a poisoned project `AGENTS.md`, a poisoned global
   `AGENTS.md`, and an ambient discoverable skill all fail to reach the
   invocation. This is the part that cannot be done in the current
   environment: `codex` is not installed and no credential is available, so
   the claim would be an assertion, not a measurement.

**Preserved argv research**, so re-enabling does not re-derive it: `exec`,
`--model <exact>`, `--ephemeral` (without it Codex persists rollout state
between invocations and the adapter stops being stateless),
`--ignore-user-config`, `--ignore-rules`, `--sandbox workspace-write`,
`--skip-git-repo-check`, `--cd <dir>`, `-` for prompt-on-stdin. Forbidden in
any future revival: `danger-full-access`, `--yolo`, `--full-auto`,
`--dangerously-bypass-approvals-and-sandbox`.

References: <https://learn.chatgpt.com/docs/developer-commands?surface=cli> ·
<https://learn.chatgpt.com/docs/agent-configuration/agents-md> ·
<https://learn.chatgpt.com/docs/build-skills>

## HELD — promote the generic cognition schema to CDD (captured, not implemented)

**Status: held until a second cognitive fill exists.** CDD is the generic
layer, CDS a concrete one; a mechanism every fill can rent belongs to CDD, and
no future fill may depend on CDS to reach it.

The **Go** layer already obeys this, and the dependency graph is what enforces
it rather than a naming convention:

```text
cellcog    -> (no internal deps)     generic cognition
cellskill  -> (no internal deps)     generic skill loading
cellwork   -> (no internal deps)     generic worktree substrate
cellfill   -> cellkernel             generic registry + cdd fills
cdspatch   -> cellcog, cellskill, cellwork, cellfill, cellkernel
```

`cellcog` is a leaf; nothing depends on `cdspatch` but the composition root. A
later `text.write` or `research` fill sits exactly where `cdspatch` sits — it
rents cognition, it does not reimplement it, and it never imports CDS.

The **CUE** layer does not yet obey it. Two generic definitions live in the
concrete overlay `schemas/cds/spec.cue`:

- `#Cognition` — mirrors `cellcog.New`, which is a generic mechanism.
- `#Hole` — a generic cellspec param-resolution concept, not a CDS one.

A second cognitive fill would have to duplicate them or import `schemas/cds`.
Both belong in `schemas/cdd/`.

What correctly stays in CDS is the **composition**, not the pieces:
`#CDSPatchAlphaResolved` binds cognition + workspace + skills into the shape a
diff-producing alpha takes. Likewise `cellwork` is generic substrate that only
diff-producing fills rent — a researcher or text cell would never touch it.

Deliberately deferred: promoting a definition with exactly one consumer is
speculative, and the move is mechanical once a real second consumer exists to
prove the boundary. Do it when the second cognitive fill lands, not before.

## Boundary the kernel never owns (Pi β #31, C3)

The episode kernel owns no GitHub, ref, PR, branch, cursor, writer-locality, or
custody policy. CLI and GitHub are **invokers/projections**. Persistence
defaults to content-like typed-receipt custody; CHAIN (`.cdd` publication)
remains an opt-in mechanism a CLI adapter may drive — never the kernel.

## Alignment with the frozen shipping plan

This migration instantiates, it does not reopen, the frozen plan
(`msg-cn-pi-cnos-shipping-plan-lock-27`): Phase K + Phase 1 = S4 core + S6 CLI
scaffolding; Phase 3 = the #717/F local CDS field proof (S7). The `.cell`
surface (Phase 4) is a CNOS addition layered *after* the mechanical CLI, not a
plan change. Repair, composition recursion, and the GitHub thin-adapter (S8)
follow only when their preceding cases have executable evidence.

## Ownership split

- **cnos.cdd** owns the *substrate*: the `#CellSpec` schema, the compiler, the
  loader, the runner (`cn cell compile` / `cn cell run`), and the skill-path
  resolver. It is the language + kernel.
- **cnos.cds** owns a *program written in it*: `main.cell`, the params-domain
  overlay, and its α/β skills. cds is one fill family; cdw/cdr are siblings of the
  same shape.

### The abstract cell (cdd) vs. the concrete cells (cds/cdr/cdw)

The "abstract cell" is a **type cdd owns**, not a base program to inherit. It is
`schemas/cdd/spec.cue #CellSpec` (+ the kernel + the empty cell). Specializations
are unifications, exactly mirroring the receipt side:

```
#CDSCellSpec: cdd.#CellSpec & { protocol_id: "…cds…", params: {language: …} }
   ⟺  #CDSReceipt: cdd.#Receipt & { protocol_id, evidence_refs }
```

There is deliberately **no language-level inheritance** (no `cdd/main.cell` base
that cds `extends`): a cell has no shared *body* to inherit — its skills and
protocol are precisely what a specialization supplies, and the shared machinery
is the kernel. The abstraction lives in the type + the kernel.

cdd's canonical instance of the abstract cell is the **empty / identity cell**
(`schemas/cdd/fixtures/empty-cell-spec.json`) — the smallest well-formed
`#CellSpec`, the data analogue of the kernel's `EmptySpec`/`cell-0`, and the
runner's reference. It runs today (`cn cell run --contract …/empty-cell-spec.json`
→ `simulated`, exit 3 — its `stub` profile is honest, non-authoritative smoke);
at Phase 4 it becomes `cdd/main.cell`. The generic cdd cell and a
concrete cds cell close through the **same kernel**; only the overlay differs.
