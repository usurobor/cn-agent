<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness] -->
<!-- completed: [Gap, Skills, ACs] -->

# α self-coherence — #370

## Gap

`COHERENCE-CELL.md` (#364, merge `32b126e4`) names the coherence-cell *organism* — receipt rule, four-way structural separation, role-as-cell-function. `RECEIPT-VALIDATION.md` (#367, merge `37ac1c75`) names the parent-facing *receptor* — `V`'s typed interface, verdict/decision distinction, δ-authoritative validation. `#366`'s roadmap §"Recursion the implementation must instantiate" carries the algorithm in informal mathematical notation but is a roadmap, not a doctrine surface.

The *algorithm* — the recursion that takes a closed cell at scope `n` and projects it as α-matter, β-discrimination, and γ-coordination at scope `n+1` — is implicit across these four surfaces but stated nowhere. Without a kernel-layer doc that names the algorithm, Phase 3 (validator) would have to derive V's contract from a recursion no doctrine names; Phase 4 (δ split) would intuit δ's signature per cycle; Phase 6 (ε relocation) would move ε without a kernel statement of its signature; Phase 7 (`CDD.md` rewrite) would have to derive the kernel mid-rewrite — exactly the failure mode two-layer separation exists to prevent.

This cycle (`#370`, Phase 1.5 of `#366`) closes that gap by producing `COHERENCE-CELL-NORMAL-FORM.md` — the substrate-independent kernel-layer companion at the doctrine layer.

**Mode:** docs-only (single new doctrine file under `src/packages/cnos.cdd/skills/cdd/` + cycle evidence under `.cdd/unreleased/370/`).

**Cycle:** 370. **Branch:** `cycle/370`. **Base SHA:** `704365d2`.

## Skills

**Tier 1 (always-on for α CDD):**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical algorithm
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface (esp. §1.4 large-file authoring discipline, §2.5 incremental self-coherence, §2.6 pre-review gate, §2.7 review-request)

**Tier 2 (always-applicable):**
- `src/packages/cnos.eng/skills/eng/*` — markdown / prose authoring (applied implicitly via the write skill)
- `src/packages/cnos.core/skills/write/SKILL.md` — every α output is a written artifact

**Tier 3 (issue-specific):**
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — form authority for any issue-pack reconciliation during cycle (not exercised — the issue and γ-scaffold were stable across the cycle)
- `src/packages/cnos.core/skills/design/SKILL.md` — kernel-articulation discipline: substrate-independence, two-layer separation, single source of truth, defer realization choices. **Load-bearing** for this cycle's authoring — without it the kernel would drift into operational details and fail AC7
- `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL.md` — predecessor doctrine; read as source-of-truth for receipt rule, four-way separation, role-as-cell-function; the organism this kernel describes algorithmically
- `src/packages/cnos.cdd/skills/cdd/RECEIPT-VALIDATION.md` — typed-interface design peer; read as source-of-truth for V's interface, verdict/decision distinction, δ-authoritative validation; the receptor this kernel positions at step 4–5

## ACs

### AC1 — `COHERENCE-CELL-NORMAL-FORM.md` exists at canonical path

**Oracle (executed):**

```
$ test -f src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md
PASS: file exists
```

**Evidence:** File created at the canonical path next to `COHERENCE-CELL.md` and `RECEIPT-VALIDATION.md`. Diff stat shows `+ src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` (435 lines, status `A`).

**Verdict:** PASS.

### AC2 — Draft-doctrine status declared and predecessors cited

**Oracle (executed):**

```
$ rg -c 'draft doctrine' src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md
2

$ rg -c 'COHERENCE-CELL\.md' src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md
4

$ rg -c 'RECEIPT-VALIDATION\.md' src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md
6
```

All three literal-token oracles return ≥1.

**Prose evidence:**
- Status frontmatter: `**Status:** Draft doctrine. Phase 1.5 of #366`.
- §Preamble opens: `This document is **draft doctrine**. It is a companion at the kernel layer to COHERENCE-CELL.md … and to RECEIPT-VALIDATION.md`.
- §Preamble §"Companion, not replacement" states explicitly: `COHERENCE-CELL.md remains the predecessor doctrine. … This document does not replace, supersede, or silently override COHERENCE-CELL.md` and `RECEIPT-VALIDATION.md remains the parent-facing receptor design`.
- §Preamble §"What this document is" declares this doc as the kernel layer (`Coherence-Cell Normal Form (CCNF) is the substrate-independent recursion algorithm of the coherence cell`).

**Verdict:** PASS.

### AC3 — Kernel as five-step recursion at scope `n` with evidence-binding rule

**Prose check (executed against §Kernel):**

- Five steps named in order, each with its typed signature:
  1. `matterₙ := αₙ.produce(contractₙ)` — §Kernel §"Step 1"
  2. `reviewₙ := βₙ.review(contractₙ, matterₙ)` — §Kernel §"Step 2"; explicitly states β consumes matter only and "β's signature explicitly excludes evidence"
  3. `receiptₙ := γₙ.close(contractₙ, matterₙ, reviewₙ, evidenceₙ)` — §Kernel §"Step 3"; γ binds evidence
  4. `verdictₙ := V(contractₙ, receiptₙ)` — §Kernel §"Step 4"; V dereferences evidence from receipt (signature explicitly shows two inputs, not three)
  5. `decisionₙ := δₙ.decide(receiptₙ, verdictₙ)` — §Kernel §"Step 5"; explicitly states "δ's signature explicitly excludes evidence" and "δ does not re-read the evidence graph"
- Evidence-binding rule pinned at kernel level in §Kernel §"Step 3" as a blockquote:
  > Evidence accumulates during α and β work. γ binds it into the receipt at close-out as typed references. V dereferences the references to validate. β never consumes evidence directly. δ never re-reads evidence.
- Four enforcement points named immediately after the rule (α: no receipt yet; β: signature excludes evidence; γ: signature includes evidence as input γ binds; δ: signature excludes evidence).

**Verdict:** PASS.

### AC4 — Four closed-cell outcomes pinned with verdict × decision preconditions

**Prose check (executed against §Cell Outcomes):**

- Section opens: `A closed cell at scope n terminates in exactly one of four outcomes.`
- Four outcomes named with explicit `(verdict, decision)` preconditions (verbatim from §"The four outcomes"):
  - `accepted := (verdictₙ = PASS)  ∧ (decisionₙ ∈ {accept, release})`
  - `degraded := (verdictₙ ≠ PASS)  ∧ (decisionₙ = override)`
  - `blocked  := (decisionₙ ∈ {reject, repair_dispatch})` — any verdict
  - `invalid  := (verdictₙ = PASS ∧ decisionₙ = override) ∨ (verdictₙ ≠ PASS ∧ decisionₙ ∈ {accept, release})`
- `invalid` declared non-terminal: §"`invalid` is non-terminal" states `**invalid cells do not close.** … δ must re-decide … before the cycle can terminate.`
- Alignment with `#369` AC4 explicit: §"Alignment with the receptor design" — `The four-outcome table aligns with #369 AC4 — the verdict × action × transmissibility table that the receptor's schema work is typing.`

**Verdict:** PASS.

### AC5 — Two recursion modes pinned: within-scope vs cross-scope

**Prose check (executed against §Recursion Modes):**

- Section opens: `The kernel's recursion has two distinct modes. They share the same recursion operator … but differ in how the scope index advances.`
- Within-scope mode named in §"Within-scope mode — repair-dispatch (same scope index)" with precondition `decisionₙ = repair_dispatch`, behaviour `cell at scope n stays open / child cell runs at scope n under repair contract / on child accept: parent γₙ re-emits a fresh receiptₙ / V re-fires at step 4 / δ re-decides at step 5`, and scope-index property `Scope index: unchanged`.
- Cross-scope mode named in §"Cross-scope mode — accept / degraded (scope index advances)" with preconditions covering `accept`/`release` for PASS and `override` for non-PASS, behaviour `cell at scope n closes / closed_cellₙ projects as α-matter at scope n+1`, and scope-index property `Scope index: advances`.
- Same-operator framing explicit: §"Same operator, different scope behaviour" — `The kernel's recursion is a single operator … The two modes are how that operator's output relates to the scope index`.
- Non-mode guardrails explicit: §"What the recursion does not do" — `Repair-dispatch is not scope-advancing` and `Accept is not same-scope`.

**Verdict:** PASS.

### AC6 — Three scope-lift projections named with projection-not-renaming framing

**Prose check (executed against §Scope-Lift):**

- Three projections named in §"The three projections":
  - `Projection 1: closed (αₙ, βₙ, γₙ) cell → αₙ₊₁ matter`
  - `Projection 2: δₙ boundary decision → βₙ₊₁-like discrimination`
  - `Projection 3: εₙ receipt-stream observation → γₙ₊₁-like coordination / evolution`
- Projection-not-renaming framing explicit in §"Projection-not-renaming framing" as a blockquote: `This is a projection under scope-lift. It is not a flat role-renaming inside a single cell. δₙ is not literally βₙ₊₁; εₙ is not literally γₙ₊₁.`
- β/γ-no-upward-projection clause explicit in §"β and γ have no upward projection": `**βₙ and γₙ have no projection to scope n+1.** Their work is intra-scope … Neither βₙ nor γₙ produces an object that crosses to scope n+1 as a separately typed projection.`
- Alignment with `#369` AC3 explicit: §"Alignment with the receptor design" — `The three-projection structure aligns with #369 AC3 — the scope-lift projection that the receptor's schema work is typing.`
- Non-projection guardrails explicit: §"What scope-lift does not do" — `β has no upward projection / γ has no upward projection / Scope-lift is not a single-cell loop`.

**Verdict:** PASS.

### AC7 — Substrate-independence enforced in the kernel section (section-bounded oracle)

**Oracle (executed):**

```
$ awk '/^## (Kernel|Cell Outcomes|Recursion Modes|Scope-Lift)/ {keep=1} \
       /^## / && !/^## (Kernel|Cell Outcomes|Recursion Modes|Scope-Lift)/ {keep=0} \
       keep {print}' \
       src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md \
   > /tmp/ccnf-kernel.md

$ wc -l /tmp/ccnf-kernel.md
275

$ rg '^## ' /tmp/ccnf-kernel.md
## Kernel
## Cell Outcomes
## Recursion Modes
## Scope-Lift

$ rg -i '\b(github|cue|cn-cdd-verify|cn dispatch|claude|gh)\b' /tmp/ccnf-kernel.md
(no matches → exit 1)
```

All four section headers extracted exactly. The substrate-term scan against the extracted kernel slice returns no matches; the `! rg -i …` assertion holds (the asserted exit-1 outcome).

**Prose evidence:** Realization-substrate names (`github`, `cue`, `cn-cdd-verify`, `CDD.md`, etc.) appear only in §Preamble (companion-citation context), §Two-Layer Separation (realization-peer citations), and §Non-goals (out-of-scope path enumeration). None appear inside `## Kernel`, `## Cell Outcomes`, `## Recursion Modes`, or `## Scope-Lift`.

**Verdict:** PASS.

### AC8 — Two-layer separation declared with realization peers cited

**Prose check (executed against §Two-Layer Separation):**

- Section opens with the kernel-is-the-*what* / realization-is-the-*how-on-this-substrate* framing as a blockquote: `Kernel is the what. Realization is the how-on-this-substrate.`
- Four realization peers cited under §"Realization peers", each with a bolded heading:
  - `RECEIPT-VALIDATION.md — parent-facing typed-interface design`
  - `schemas/cdd/ (Phase 2, #369, in flight) — typed schemas for receipt, contract, boundary decision`
  - `cn-cdd-verify (Phase 3, deferred) — V's command-wrapper implementation`
  - `CDD.md (Phase 7 rewrite, deferred) — canonical executable algorithm`
- Direction-of-dependency rule stated in §"What lives where": `realization cites kernel; kernel does not cite realization`.

**Oracle (token presence in the file):**

```
$ rg -c 'RECEIPT-VALIDATION\.md' …  → 6
$ rg -c 'schemas/cdd/'         …    → 3
$ rg -c 'cn-cdd-verify'        …    → 5
$ rg -c 'CDD\.md'              …    → 11
```

All four realization peers cited (counts include the §Two-Layer Separation section plus secondary references in §Preamble, §Non-goals, and §Closure).

**Verdict:** PASS.

### AC9 — Non-goals respected; surface containment proved

**Oracle (executed):**

```
$ git diff origin/main..HEAD --stat
 .cdd/unreleased/370/alpha-codex-prompt.md          |  58 +++
 .cdd/unreleased/370/beta-codex-prompt.md           |  45 +++
 .cdd/unreleased/370/gamma-scaffold.md              | 140 +++++++
 .cdd/unreleased/370/self-coherence.md              |  32 ++   (will grow as ACs/Self-check/Debt/CDD-Trace/Review-readiness land)
 .../skills/cdd/COHERENCE-CELL-NORMAL-FORM.md       | 435 +++++++++++++++++++++
 5 files changed
```

Surfaces present in the diff: exactly the new doc + cycle evidence under `.cdd/unreleased/370/`. No edits to `COHERENCE-CELL.md`, `RECEIPT-VALIDATION.md`, `CDD.md`, `schemas/cdd/`, `cn-cdd-verify`, `operator/SKILL.md`, `gamma/SKILL.md`, `epsilon/SKILL.md`, `ROLES.md`, or CI workflows.

**Doc-body restatement of non-goals:** §Non-goals enumerates each prohibited surface with the matching phase reference (Phases 2–7 + predecessor doctrine + receptor design + CI / package boundaries).

**Verdict:** PASS.
