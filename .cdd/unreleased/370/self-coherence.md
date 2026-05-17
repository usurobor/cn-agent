<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness] -->
<!-- completed: [Gap, Skills, ACs, Self-check, Debt, CDD-Trace, Review-readiness] -->

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

## Self-check

**Did α's work push ambiguity onto β?** No. The nine ACs are each mapped to specific evidence in `COHERENCE-CELL-NORMAL-FORM.md` (with section + clause locations) and four ACs are backed by directly executable oracles with captured stdout (AC1, AC2, AC7, AC9). β does not need to re-derive the AC mapping or re-execute the shell oracles to verify them — the prose-evidence sections name exact section anchors and the oracle blocks above are self-contained.

**Is every claim backed by evidence in the diff?**
- AC1, AC9: backed by `git diff --stat` evidence (`+ COHERENCE-CELL-NORMAL-FORM.md` plus cycle-evidence files only).
- AC2: backed by literal-token rg counts (`draft doctrine` ×2, `COHERENCE-CELL.md` ×4, `RECEIPT-VALIDATION.md` ×6).
- AC3–AC6: backed by named §Kernel / §Cell Outcomes / §Recursion Modes / §Scope-Lift sections with verbatim quotes of the load-bearing statements (five-step signatures, evidence-binding rule, four-outcome preconditions, both recursion modes, three projections, projection-not-renaming framing, β/γ-no-upward-projection clause).
- AC7: backed by the `awk` extraction + word-bounded `rg` substrate scan returning no matches over a 275-line kernel slice with exactly four section headers.
- AC8: backed by §Two-Layer Separation enumerating four realization peers with explicit bolded headings.

**Peer enumeration.** The peer family at risk of drift is the *predecessor + receptor + canonical algorithm* set (`COHERENCE-CELL.md`, `RECEIPT-VALIDATION.md`, `CDD.md`). The cycle's surface-containment contract leaves all three predecessors unedited (verified by `git diff --stat`). The new file is *additive* — a kernel-layer companion citing the predecessors as load-bearing inputs — and explicitly states "this document does not replace, supersede, or silently override `COHERENCE-CELL.md`." No peer surface drift.

**Sibling / harness audit.** No code, schemas, harness, or CI workflows touched; docs-only. The schema-bearing harness surfaces (`schemas/cdd/`, `cn-cdd-verify`, role skills) are explicitly out-of-scope for this cycle and cited only as labelled realization peers in §Two-Layer Separation — they do not consume kernel statements operationally yet; they will when their respective phases land.

**Substrate-leakage failure mode.** γ-scaffold §"Failure modes" #1 pre-flagged substrate leakage in kernel sections as the primary risk. The kernel-slice scan was re-executed after every section landed; after one false start (a `[#369](https://github.com/...)` URL in §Cell Outcomes that the kernel-slice scan caught immediately on the first run), the URL was replaced with the bare `#369` reference and the scan returned no matches. The mechanism worked as designed.

**Length discipline.** Target was 200–400 lines. Final: 435 lines (kernel slice: 275 lines; preamble + two-layer separation + non-goals + closure: 160 lines). 9% over the upper target. The non-kernel sections carry: predecessor citation framing (AC2), realization-peer enumeration (AC8), non-goals enumeration (AC9), and phase-inheritance table. The kernel slice itself is well within target. Recorded as observation, not debt — see §Debt.

**Role boundary.** α did α's work: produced the kernel-layer doc, mapped ACs to evidence, ran executable oracles, recorded self-coherence. α did not edit predecessor surfaces, did not edit harness, did not anticipate β's triage frame.

## Debt

**Known debt:**

1. **Length: 435 lines vs 200–400 target.** Doc lands 35 lines (9%) over the upper target. The overage lives in non-kernel sections (preamble + two-layer separation + non-goals + closure). The kernel slice alone is 275 lines, well within target. Disposition: observation, not requiring fix — the non-kernel sections carry AC2/AC8/AC9 evidence and the phase-inheritance handoff. β triage may flag for trim if the assessment is that decisive-over-exhaustive discipline drifted; α's reading is that each non-kernel paragraph load-bears.

2. **Predecessor doctrine has a different `## Recursion Equation` formulation than this kernel's §Kernel.** `COHERENCE-CELL.md` §Recursion Equation states `receiptₙ := closed_cellₙ.receipt` as a separate line and treats `γ.close` as fusing `(contract, matter, review)` without naming evidence explicitly. This kernel's §Kernel folds receipt emission into `γₙ.close(contractₙ, matterₙ, reviewₙ, evidenceₙ)` and pins the evidence-binding rule as load-bearing. The two are consistent (the kernel adds the evidence input to the predecessor's three-input close), but the predecessor's surface does not state the evidence-binding rule. This is *deliberate*: the gap is exactly what this cycle exists to close. Disposition: not debt; this is the kernel's contribution. Recorded for β's awareness.

3. **`#369` AC4 and AC3 alignment is cited at doctrine level only.** This kernel commits to the four outcomes and three projections as the *what*; `#369`'s schema work commits to the *type*. The two cycles run in parallel and have not yet reconciled their concrete shapes. If `#369`'s schema work pins outcome preconditions inconsistent with the kernel's, the kernel's statement is the load-bearing claim and `#369` aligns to it — but the actual schema files are not in this diff and cannot be co-verified by this cycle. Disposition: known coordination point with the parallel cycle; not a blocker for this cycle's closure.

4. **No debt under AC1–AC9.** All nine acceptance criteria pass under the per-AC oracles defined by the issue body.

**No deferred sub-ACs.** Every AC the issue body names is fully exercised in this cycle.

## CDD-Trace

1. **Issue:** `#370` — Phase 1.5 of `#366` (coherence-cell executability roadmap). Docs-only doctrine companion. 9 ACs.

2. **Active skills:** Tier 1 (`CDD.md`, `alpha/SKILL.md`); Tier 2 (`write/SKILL.md`); Tier 3 (`design/SKILL.md`, `issue/SKILL.md` — not exercised, `COHERENCE-CELL.md` + `RECEIPT-VALIDATION.md` as source-of-truth reads).

3. **Design:** Issue body + γ-scaffold (`gamma-scaffold.md`) carry design. The γ-scaffold pre-flagged nine failure modes; α's authoring discipline executed against all nine, including section-header conformance to the AC7 awk pattern.

4. **Plan:** Implicit in the dispatch — author the doc section-by-section with one commit per section per `CDD.md` §1.4 large-file authoring rule. Section order: Preamble → Kernel → Cell Outcomes → Recursion Modes → Scope-Lift → Two-Layer Separation → Non-goals → Closure. Manifest at top of file updated after each section landed.

5. **Tests:** Not applicable (docs-only). Acceptance evidence is the per-AC oracle bank in the issue body — executed and recorded inline in §ACs above.

6. **Code:** None. Single new doctrine file:
   - `src/packages/cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` (435 lines, status `A`) — the kernel-layer companion this cycle produces. Authored across six commits (one per major section + scaffold). The kernel slice (the four `## Kernel`, `## Cell Outcomes`, `## Recursion Modes`, `## Scope-Lift` sections) is 275 lines. Each section's authoring is verifiable from the per-section commit on `cycle/370`.

7. **Docs:** Done as the cycle's primary artifact (the new file *is* doc). Cycle-evidence files under `.cdd/unreleased/370/` (`alpha-codex-prompt.md`, `beta-codex-prompt.md`, `gamma-scaffold.md`, `self-coherence.md`) carry α's and γ's working state. No other docs touched.

8. **Self-coherence:** This file. §Gap, §Skills, §ACs (with executable-oracle results inline), §Self-check, §Debt, §CDD-Trace landed across multiple commits; §Review-readiness lands as the final commit signaling β polling.

## Review-readiness | round 1

| Property | Value |
|---|---|
| Round | 1 |
| Base SHA (observed at signal time) | `704365d23378fcbfcf1e33679025809af6b81100` — matches scaffold's recorded base; `origin/main` not advanced since dispatch (re-fetched at signal time) |
| Implementation SHA | `10b41762` — last implementation commit prior to this readiness-signal commit (per `alpha/SKILL.md` §2.6 SHA convention; this signal advances HEAD but cites the stable prior SHA) |
| Branch | `cycle/370` |
| Branch CI | N/A — docs-only cycle, no CI workflows trigger on `cycle/*` |
| α identity | `alpha@cdd.cnos` (verified via `git log -1 --format='%ae'`) |
| AC verdict | AC1–AC9: PASS (all per-AC oracles executed and recorded inline above) |
| Pre-review gate | All applicable rows pass (§"Pre-review gate" alpha/SKILL.md §2.6); doc rows N/A logged below |
| Ready for β | Yes |

**Pre-review gate row-by-row (alpha/SKILL.md §2.6):**

1. ✓ `origin/cycle/370` is on `origin/main`'s line; base SHA `704365d2` matches; no rebase needed.
2. ✓ CDD Trace through step 7 (§CDD-Trace step 7 names docs as the cycle's primary artifact; step 8 is this self-coherence file).
3. N/A — tests not applicable for a docs-only doctrine cycle; the per-AC oracle bank in §ACs is the analogous executable surface, and all PASS.
4. ✓ Every AC has evidence (§ACs).
5. ✓ Known debt explicit (§Debt).
6. N/A — no schema-bearing contract changed in this cycle; the kernel positions `V` and the four outcomes at a doctrine level only.
7. ✓ Peer enumeration completed (§Self-check); predecessor + receptor + canonical-algorithm peers all unedited (verified by diff).
8. N/A — no schema-bearing contract changed; no harness audit applicable.
9. N/A — no mid-cycle patches; clean section-by-section authoring.
10. N/A — docs-only; no CI gate.
11. ✓ Artifact enumeration matches diff. Every file in `git diff --stat origin/main..HEAD` is named in this self-coherence: `COHERENCE-CELL-NORMAL-FORM.md` (§CDD-Trace step 6); `.cdd/unreleased/370/gamma-scaffold.md`, `.cdd/unreleased/370/alpha-codex-prompt.md`, `.cdd/unreleased/370/beta-codex-prompt.md` (dispatch artifacts produced by γ, named in §CDD-Trace step 3); `.cdd/unreleased/370/self-coherence.md` (this file, §CDD-Trace step 8).
12. N/A — no new modules or functions; docs-only.
13. N/A — no test runner output applicable. Executable oracle bank (AC1, AC2 ×3, AC7, AC9) ran with PASS results captured inline in §ACs.
14. ✓ α commit author email matches canonical pattern: `alpha@cdd.cnos` per `operator/SKILL.md` git identity for the cnos project (using elision form). All α commits on the cycle branch from `5b772567` onward carry this identity.

**β polling begins now.** α will poll `.cdd/unreleased/370/beta-review.md` and the issue every 60s per `alpha/SKILL.md` §2.7 until β returns a verdict.
