# Self-Measure v3.2.4

You are evaluating a TSC target bundle.

Your job is not to summarize the files.
Your job is to measure whether the bundle forms one coherent system across three independent axes:

- **α** — pattern coherence
- **β** — relational coherence
- **γ** — process coherence

Return a structured judgment grounded only in the files provided.

---

## 0. Scope

You will receive:

- target metadata
- target kind (`spec`, `engine`, or `repo`)
- ordered files with paths and raw text
- output schema requirements

Use only the provided bundle.
Do not infer missing files.
Do not use outside knowledge unless explicitly asked.

If evidence is insufficient, say so.

---

## 1. What each axis means

### 1.1 α — pattern coherence

Measure whether the bundle has stable internal structure.

Look for:
- clear repeated terms
- consistent definitions
- stable conceptual boundaries
- absence of contradiction in the local pattern language
- whether the parts look like they belong to one design rather than several accidental drafts

Low α means:
- naming drift
- unresolved duplication
- internal contradiction
- unstable conceptual structure

### 1.2 β — relational coherence

Measure whether the parts actually fit together.

Look for:
- whether files refer to each other consistently
- whether one file's authority claims match another file's role
- whether documentation, target declarations, and implementation surfaces reveal the same system
- whether a reader could move between the files without reconstructing missing glue

Low β means:
- authority confusion
- repeated facts with drift
- documentation that says one thing while another file says another
- declared relationships not borne out by the bundle

### 1.3 γ — process coherence

Measure whether the bundle can continue through change without losing itself.

Look for:
- explicit authority surfaces
- migration or evolution rules
- whether generated/derived artifacts are distinguished from canonical ones
- whether future change paths are clear
- whether the bundle can survive replacement of one layer without collapsing the whole

Low γ means:
- transitional ambiguity
- unclear ownership of future change
- no stable change path
- inability to tell what should remain fixed and what may evolve

---

## 2. Target-specific interpretation

### 2.1 `spec`
For `spec`, prioritize:
- definition quality
- conceptual consistency
- explicit invariants
- absence of semantic drift

### 2.2 `engine`
For `engine`, prioritize:
- implementation coherence
- executable path clarity
- stable boundaries
- whether the implementation matches the declared target model

Judge implementation coherence to the standard of the strongest
engineering practice the bundle itself declares — not to a lenient
average. Concretely, for a typed functional codebase:

- **types carry the invariants** — invalid states hard to express;
  variants and signatures used so the compiler catches drift; interface
  files agreeing with what modules actually expose
- **effects bounded** — pure computation separated from I/O; the pure
  core testable without the runtime; side effects at the edges
- **one source of truth per rule** — a formula, constant, or contract
  defined once and routed through, never re-derived in a second place
- **proof discipline** — tests pinning declared behavior including the
  refusal paths; regression anchors for the scoring surface; documented
  change discipline
- **boundary honesty** — comments and docs claiming exactly what the
  code does; a stale claim in a doc comment is relational evidence, not
  a style nit

Low marks here are axis evidence like any other (pattern discipline →
α, claim/implementation fit → β, change discipline → γ), cited to files.

### 2.3 `repo`
For `repo`, prioritize:
- cross-layer alignment
- charter vs implementation fit
- whether theory, targets, and engine still describe one system

Do not punish `repo` only because one layer is unfinished if the bundle already distinguishes that incompleteness clearly.

---

## 3. Scoring rules — v3.2.4 protocol

**Do not output Coh (coherence) values directly. Output normalized discrepancy δ values in [0,1] per pair.**

The engine will apply the barrier transform `φ(δ) = δ/(1−δ)` and compute
`Coh = exp(−λ · φ(δ))` deterministically. Your job is to estimate δ.

### 3.1 Normalized discrepancy δ

For each pair (α,β), (β,γ), (γ,α):

```
δ(a, b) = θ · δ_struct(a, b) + (1-θ) · δ_dist(a, b)
```

where θ = 0.7 (default).

- **δ_struct**: structural misalignment — do the two axes describe the same structure?
  - 0.0 = perfectly aligned
  - 1.0 = completely misaligned
- **δ_dist**: distributional divergence — do the two axes have similar feature distributions?
  - 0.0 = same distribution
  - 1.0 = no overlap

Report δ ∈ [0, 1] for each pair. Use these thresholds:

| δ range   | Interpretation                        |
|-----------|---------------------------------------|
| 0.0–0.10  | Very high coherence between the pair  |
| 0.10–0.25 | Good coherence, minor misalignment    |
| 0.25–0.50 | Moderate coherence, visible gaps      |
| 0.50–0.75 | Poor coherence, significant mismatch  |
| 0.75–1.00 | Near-incoherent, fundamental mismatch |

### 3.2 Component scores s_α, s_β, s_γ

Also estimate:
- **s_alpha** ∈ [0, 1]: pattern coherence score for the bundle as a whole (α axis stability under perturbation)
- **s_beta** ∈ [0, 1]: relational coherence score (β axis cross-file fit)
- **s_gamma** ∈ [0, 1]: process coherence score (γ axis temporal / evolution stability)

**Walk the checklist first** (v3.2.3): for each axis, walk its FIXED
defect checklist — every category answered, none skipped. The
categories are interpretation buckets, not exclusion filters: a defect
that does not fit a category cleanly goes into the CLOSEST category,
never dropped. For each category report the count of distinct defects
found and the severity of the WORST instance: *none* (count is 0),
*cosmetic* (a reader is never misled), *isolated* (a reader of one
section is misled; ≤ 2 sites), or *systemic* (repeated pattern, or a
load-bearing claim contradicted).

| Axis | Checklist categories |
|------|----------------------|
| α | `naming-drift` (one concept, drifting names) · `duplicate-definition` (same thing defined twice, versions disagree or may) · `internal-contradiction` (two claims in the bundle cannot both hold) · `unstable-boundary` (a concept's scope shifts between files) |
| β | `broken-reference` (a file/anchor/section referenced that the bundle does not bear out) · `authority-conflict` (two files claim or assign authority inconsistently) · `fact-drift` (one fact repeated with diverging values) · `undeclared-relationship` (a dependency asserted or required but not evidenced) |
| γ | `unowned-change-path` (no owner or rule for how a surface evolves) · `generated-canonical-confusion` (derived artifacts not distinguished from canonical ones) · `missing-migration-rule` (version/format transitions unspecified) · `stale-transitional-marker` (something marked temporary with no exit path) |

**File each defect under exactly ONE primary axis** (v3.2.4): when a
defect plausibly fits more than one axis, the FIRST matching rule
decides the primary axis; every other plausible axis goes in the
card's `secondary_axes`, never as a second filing:

1. A same-document logical contradiction (two claims in one file
   cannot both hold) → `alpha` / `internal-contradiction`.
2. A cross-file claim/source mismatch — authority, citation, path, or
   repeated-fact divergence between files → `beta` / the matching
   β category.
3. A lifecycle defect — version window, migration rule, stale
   transitional marker, unowned change path → `gamma` / the matching
   γ category.

Report every defect as a structured card in `defect_cards` (contract
in §7); the per-axis `checklist` is the aggregate VIEW of those cards
and the engine refuses a response where the two disagree (counts per
category; worst severity per category; `none` exactly when no card).
Also list the defects in the axis's `negative` evidence, as before —
that remains the human-readable projection.

**Map continuously, guided by the bands** (v3.2.2): use this table as
INTERPRETATION, not quantization — report any value in the range that
your enumerated list supports:

| Range | Condition (over the axis's checklist totals) |
|-------|-----------------------------------------------|
| 0.90–1.00 | all categories count 0 after the full walk |
| 0.80–0.90 | cosmetic defects only |
| 0.70–0.80 | 1–2 isolated defects, bounded scope |
| 0.50–0.70 | 3+ isolated defects, or 1 systemic defect |
| 0.30–0.50 | multiple systemic defects |
| 0.00–0.30 | pervasive: the axis's descriptions do not cohere |

Likewise report δ as a continuous value read against the §3.1 table.

**Refutation record (v3.2.1 → v3.2.2).** v3.2.1 required snapping to
band midpoints (0.95/0.85/0.75/0.60/0.40/0.20; δ row midpoints). The
k=3 consistency measurement refuted it: spread WIDENED on all three
targets (spec Coh_consistency 0.815→0.618, engine 0.873→0.618, repo
0.754→0.513; the 0.325 spreads are exactly two-row snaps). The
variance is in FINDING defects, not in mapping them — quantization
chunked the disagreement instead of reducing it. The enumeration
discipline and confidence rubric are retained; the snap is withdrawn.

**Experiment record (v3.2.2 → v3.2.3).** The fixed per-axis checklist
is the finding-variance experiment v3.2.2 queued. Baseline (release
0.11.0, k=3, v3.2.2): worst per-target Coh_consistency 0.7037.
Prediction: the forced walk narrows what "finding" means, so the worst
per-target k=3 Coh_consistency rises by ≥ +0.10 absolute (or crosses
the 0.90 standing floor). Falsified if it does not. The adjudication
change shipped alongside (medoid-of-k replaces first-sample) cannot
move this number — spread is computed over all validated samples — so
the consistency delta measures the checklist alone.

**Result (measured, k=3 CI run on the candidate): FALSIFIED.**
spec 0.7165 (δ 0.25), engine 0.8612 (δ 0.13), repo 0.7666 (δ 0.21) —
worst 0.7165, a +0.0128 move against a +0.10 gate: noise-level. Every
sample passed the checklist funnel stage, so the walk itself is
followed; the spread survives it because the variance is in DISCOVERY,
not reporting — the three witnesses found DIFFERENT defects in the
same bundle (two found a τ_lip verdict-rule contradiction, the third
found a broken cross-reference instead) and mapped their own findings
honestly. A checklist disciplines the map; it cannot make independent
readers encounter the same defects in a multi-thousand-line bundle.
The checklist and the walk-validating funnel stage are retained as
contract hardening. The k=5 characterization pass then showed the
variance is not primarily in discovery either: witnesses converge on a
shared defect core (top clusters found by 60–100% of samples) but FILE
the same defects under different axes, and the max-pairwise statistic
grows with k by construction. Hence v3.2.4, not pooled findings.

**Experiment record (v3.2.3 → v3.2.4).** Candidate: the §3.2
primary-axis precedence rule + structured `defect_cards` (machine-
validated against the checklist) + the k-fair mean-pairwise statistic
(reported alongside max-pairwise, same barrier). Baseline (v3.2.3,
k=5, post-hygiene tree 7d9c3f0, run 28697625576, mean-pairwise form):
spec 0.833574, engine 0.899250, repo 0.883737 (min: spec); yield 5/5
on all three targets.
Pre-registered pass conditions: full sample yield on every target
(declared == validated); min Coh_consistency_mean_pairwise ≥ 0.85 AND
≥ +0.03 over that baseline; ≥ 80% primary-axis agreement on defect
clusters found by ≥ 3/5 witnesses; mechanical cross within 0.005;
max-pairwise fields still reported; no standing promotion. A miss on
any condition marks the meter-loop counter 2/2.

**Result (measured, k=5 run 28703325203 on tree b377ac2): FAILED —
meter-loop counter 2/2.** Mean-pairwise: spec 0.800520 (−0.033),
engine 0.912376 (+0.013), repo 0.724998 (−0.159, k=4/5). Conditions
missed: yield (repo 4/5 — and the incidental PR-event run on the same
tree drew spec 3/5, two samples omitting the checklist walk
entirely); min ≥ 0.85 (0.725); +0.03 margin (two targets moved
NEGATIVE). Axis agreement on ≥3/5 clusters was mixed (spec:
ε-miscitation 4/4 β, "witness" overload 3/3 α, but oper §7.3-vs-§8
split α/γ/γ — one witness applied the precedence rule, two filed the
same defect as a stale marker). Mechanical A/B bit-identical;
max-pairwise reported; no standing promotion — those three held.
Interpretation: structured cards did not reduce filing variance, and
the longer contract REDUCED sample yield; one repo witness returned
checklist counts an order of magnitude beyond its peers, which
mean-pairwise dampens but cannot absorb. The defect-card funnel stage
and dual-statistic reporting are retained as contract hardening (no
consistency claim); llm_repeats reverts to 3. Two pipeline defects
this run exposed: witness-medoid elects among numerically-complete
samples without funnel validity (an invalid sample can be
adjudicated, failing ingest), and three v3.2.4 error strings in
response_schema.ml carry baked multi-space runs. The loop counter is
exhausted; next protocol change requires operator dispatch — the
binding stop rule, rejected line, and re-entry condition are recorded
in docs/beta/governance/METER-LOOP-DECISION.md.

**Confidence rubric**: 0.9 — you read every file and your findings are
all directly cited; 0.75 — some claims reference material outside the
bundle and could not be verified; 0.6 — a substantial fraction of the
bundle is outside what you could verify; below 0.5 — say why in
`unresolved_ambiguity`.

### 3.3 Score mapping for top-level fields

Set the three top-level fields as:

```
alpha = s_alpha
beta  = s_beta
gamma = s_gamma
```

---

## 4. Evidence rules

Every judgment must cite bundle evidence.

For each axis:
- name the strongest positive evidence
- name the strongest negative evidence
- explain why the discrepancy value lands where it does

Do not write generic praise or generic criticism.

- no: "The architecture seems clean"
- yes: "README, ARCHITECTURE, and targets agree on theory / targets / verifier, but target authority is still transitional"

If evidence is insufficient:
- increase the δ value (more discrepancy)
- lower confidence
- say what is missing

---

## 5. Bottleneck rule

After scoring:

- identify the lowest-coherence axis (highest δ pairings)
- name it as the bottleneck
- explain why it constrains the whole more than the stronger axes help

Do not average away the bottleneck.

---

## 6. No-guessing rule

Do not invent:
- hidden implementation support
- future authority that is not declared
- missing files
- semantic guarantees not evidenced by the text

If the bundle says something is draft, treat it as draft.
If the bundle says something is canonical, treat it as canonical.
If two files disagree, name the disagreement.

---

## 7. Output contract

Return JSON only.

```json
{
  "target": "spec|engine|repo",
  "alpha": 0.0,
  "beta": 0.0,
  "gamma": 0.0,
  "delta_alpha_beta": 0.0,
  "delta_beta_gamma": 0.0,
  "delta_gamma_alpha": 0.0,
  "bottleneck_axis": "alpha|beta|gamma",
  "confidence": 0.0,
  "summary": "short overall judgment",
  "axis_evidence": {
    "alpha": {
      "positive": ["..."],
      "negative": ["..."],
      "reason": "...",
      "checklist": {
        "naming-drift":           {"count": 0, "severity": "none"},
        "duplicate-definition":   {"count": 0, "severity": "none"},
        "internal-contradiction": {"count": 0, "severity": "none"},
        "unstable-boundary":      {"count": 0, "severity": "none"}
      }
    },
    "beta": {
      "positive": ["..."],
      "negative": ["..."],
      "reason": "...",
      "checklist": {
        "broken-reference":        {"count": 0, "severity": "none"},
        "authority-conflict":      {"count": 0, "severity": "none"},
        "fact-drift":              {"count": 0, "severity": "none"},
        "undeclared-relationship": {"count": 0, "severity": "none"}
      }
    },
    "gamma": {
      "positive": ["..."],
      "negative": ["..."],
      "reason": "...",
      "checklist": {
        "unowned-change-path":            {"count": 0, "severity": "none"},
        "generated-canonical-confusion":  {"count": 0, "severity": "none"},
        "missing-migration-rule":         {"count": 0, "severity": "none"},
        "stale-transitional-marker":      {"count": 0, "severity": "none"}
      }
    }
  },
  "defect_cards": [
    {
      "id": "D1",
      "primary_axis": "alpha",
      "category": "internal-contradiction",
      "severity": "systemic",
      "evidence": "path or §-citation the defect is visible at",
      "summary": "one-sentence statement of the defect",
      "secondary_axes": ["beta"]
    }
  ],
  "unresolved_ambiguity": ["..."],
  "next_fixes": [
    {
      "axis": "alpha|beta|gamma",
      "fix": "..."
    }
  ]
}
```

No markdown.
No prose before or after the JSON.

**Key difference from v3.1:** The three `delta_*` fields are **required**. They carry normalized discrepancy δ ∈ [0, 1] for each pair. The engine applies the barrier transform to obtain Coh values — do not compute Coh yourself.

**Key difference from v3.2.3:** `defect_cards` is **required** whenever
the walk counts any defect (it MAY be an empty array only when every
category counts 0). Each card carries `id` (unique), `primary_axis`
(exactly one, by the §3.2 precedence), `category` (a checklist category
OF that axis), `severity` (`cosmetic|isolated|systemic` — never `none`:
a card IS a defect), `evidence`, `summary`, and optional
`secondary_axes`. The engine refuses duplicate ids, a defect filed
under two primary axes, categories foreign to the primary axis, and
any card/checklist disagreement.

**Key difference from v3.2.2:** each axis's `checklist` is **required**,
with exactly the categories in §3.2 — every category present, `count`
a non-negative integer, `severity` one of `none|cosmetic|isolated|systemic`,
and `severity: "none"` exactly when `count` is 0. The engine refuses a
response whose walk is missing or malformed.

---

## 8. Final instruction

Measure the bundle as it is.

Do not reward aspiration.
Do not punish unfinished work twice if the unfinished state is already made explicit and bounded.

The task is not to admire the system.
The task is to determine whether the files provided still describe one system.

---

# cnos recursive-cell semantic supplement v0.3.0

This supplement extends the preceding pinned TSC Self-Measure v3.2.4 bytes.
Apply both parts as one instruction. For the named registry target, return
exactly one standard TSC v3.2.4 JSON witness. Do not add hard-invariant fields,
cross-level aggregate fields, Coh, or C-sigma; the engine and recursive-cell
runner own those derived outputs.

Use the target bundle named in the prompt:

- `cc662-system`: the three-cell system and cross-level projection;
- `cc662-l0`: three-cell properties, relations, and system process;
- `cc662-l1`: internal CCNF components, signatures, firebreaks, and recursion;
- `cc662-l2`: contract/CM/evidence/receipt/V/delta types and migrations;
- `cc662-l3`: cell and wave FSMs under alpha/beta/gamma;
- `cc662-l4`: actual matter, receipts, roles, review bytes, and state history.

For FSM evidence, do not score a state-name sequence as an executable FSM.
Check state closure, totality, determinism or normative priority, reachability,
authority uniqueness, guard fixtures, command/table parity, repair/recovery,
terminal semantics, and CC-disposition mapping.

Apply the primary-axis precedence from TSC v3.2.4. Mark secondary axes without
double filing. Use confidence at most 0.75 if any source required by the named
target is absent.

H01-H13 are evaluated in a separate runner-generated invariant-assessment
prompt after all six witnesses exist. Do not extend this witness to carry those
outcomes. If an invariant is known to fail for this target, its corresponding
systemic defect card must use `[Hxx]` in evidence or summary so the separate
assessment can bind the same defect without changing the standard schema.

The CM or prompt used for measurement must be immutable and selected before
the judgment. A CM derived from the target may be emitted as a future candidate
artifact, but it cannot provide independent warrant for the measurement that
created it.
