# β review — cnos#628, independent review of PR #629

**Independence note:** this review is authored by the `cds-dispatch` wake's δ-orchestrated β pass,
which did not author PR #629 and has no prior-session history with it — satisfies the α≠β
contagion firebreak the issue body's ratification-path section names as the reason κ could not
self-review. Verified by reading the actual diff and file content, not by trusting the PR body's
claims (per the issue's own "Oracles are grep/citation checks against the landed doctrine" line).

## §R0

### Headline finding — dispositive, blocks convergence

**PR #629's branch is NOT actually rebased onto current `main` as its own body claims.** Verified:

```
git merge-base --is-ancestor main origin/sigma/cell-runtime-arch-note
# => main is NOT an ancestor of the PR branch
```

The PR branch silently **deletes** already-landed doctrine that exists on current `main` and is
absent from the diff's stated scope:

- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.13 "scope-continuation shape (cnos#639)" —
  entire subsection (48 lines) removed with **no replacement**. This is not an edit; `git diff
  main..origin/sigma/cell-runtime-arch-note -- .../delta/SKILL.md` shows only `-` lines for this
  hunk.
- `.cdd/unreleased/639/*` — 8 files, ~1300 lines (gamma-scaffold.md, self-coherence.md,
  beta-review.md, alpha-closeout.md, beta-closeout.md, gamma-closeout.md, deliverable-evidence.md,
  CLAIM-REQUEST.yml, REVIEW-REQUEST.yml) — all absent from the PR branch.
- `.github/workflows/cnos-cds-dispatch.yml`, `src/packages/cnos.cds/orchestrators/cds-dispatch/SKILL.md`,
  `.../cds-dispatch/cnos-cds-dispatch.golden.yml`, `src/packages/cnos.core/skills/agent/dispatch-protocol/SKILL.md`
  — every hunk in all four is a pure removal of the `scope_continuation`/cnos#639 content; the PR
  branch adds nothing new to any of these four files (independently confirmed — no `+` content
  beyond context lines).

The PR body's own "Changes (6 files)" accounting never mentions any of these 21-total touched
files beyond the 5 doctrine files it intends to change (`CELL-RUNTIME.md`, `CCNF.md`,
`CELL-KINDS.md`, `CDD.md`, `README.md`) — the regression is undisclosed in the PR's own summary,
not just present.

**Why this is dispositive:** merging this PR as-is would silently regress `main`, deleting a
landed, cross-referenced doctrine section (this very cell's own `gamma-scaffold.md` cites §9.13 as
prior art for classification). This is exactly the failure mode CDD's receipt/review discipline
exists to catch before it reaches `main`.

### Per-AC table (substantive content, independent of the headline finding)

| AC | Verdict | Evidence |
|---|---|---|
| AC1 — four orthogonal axes; ε cross-cell; CCNF kernel canonical | **Satisfied** | `docs/architecture/CELL-RUNTIME.md` (new file) "Four orthogonal classifications" table; explicit "ε is not a step in that closure — it observes receipt streams *across* cells." |
| AC2 — WC/PC/CC full role loop (not α/β/γ relabeled); PC=relator; S₃-invariance; `--class` selector | **Satisfied** | "Every WC, PC, and CC runs its own internal, independent α→β→γ→V→δ"; "relator / relational-structurer" framing for PC; explicit "S₃-invariance" subsection; `--class working\|planning\|cohering` selector named (not `--kind`). |
| AC3 — CELL-KINDS.md re-headed in place; filename preserved; heterogeneity classification; compat note; #614/#640 preserved; no physical rename; CDD.md minimal pointer | **Satisfied** | Banner + body re-headed (`# CELL-KINDS — Legacy Cell-Kind Taxonomy, now Domain Vocabulary`); heterogeneity table (matter domains / contract modes / coordination-boundary) present; `cell_kind`→`cell_class`+`matter_domain` compat note present; `## Mandatory terminal learning section` (#614) and `## Process self-improvement loop (cnos#640)` (#640) both still present on the PR branch (re-headed contextually, not deleted); no `CELL-DOMAINS.md` added; `CDD.md` diff is a single-line pointer fix. |
| AC4 — CM↔V settled w/o schema change; V=PASS\|FAIL only; escalation=runner routing; override=δ action; "WARN is not a verdict"; CCNF corrected | **Satisfied** | `schemas/` untouched in the diff; CCNF diff states "WARN is not a verdict value in the shipped schema (schemas/cdd/receipt.cue pins #ValidationVerdict.verdict ∈ {PASS, FAIL})", removing prior "or (optionally) WARN" language; CELL-RUNTIME.md states escalation is runner routing and override is a δ `#Override` boundary decision, not a verdict. |
| AC5 — CELL-RUNTIME.md exists, cross-linked from README + CCNF + CELL-KINDS | **Satisfied** | 111-line substantive new file; README.md adds a "## Runtime" section linking it; CCNF adds a citing paragraph; CELL-KINDS.md Cross-references section adds a citing line. |
| AC6 — escalation-predicate + no-silent-no-op invariants | **Satisfied** | CELL-RUNTIME.md "Mechanical-first, no silent no-op" section: deterministic/logged escalation predicate; every pulse emits a receipt stating why, never an empty green. |
| AC7 — #366 reconciliation recorded; no new phase issue until S1 reviewed | **Satisfied (as far as diff-checkable)** | PR body has a dedicated "#366 reconciliation (bounded)" section; not independently verifiable beyond the PR body's own text. |
| AC8 — doctrine receipt names settlements 1–5, migration path, ratification path, residual questions, with citations | **Satisfied, minor nit** | PR body's "Doctrine receipt (settled decisions)" section lists all 5 settlements with citations (CCNF #370, COHERENCE-CELL.md #364, CELL-KINDS.md #570, schemas/cdd/{receipt,boundary_decision}.cue #369, #366, tsc spec); "Ratification path" section present; residual questions under "Known Debt." Nit: the literal words "constitutive change"/"migration" don't appear in the PR body proper (they appear inside the new CELL-RUNTIME.md doctrine file itself) — substance is present, vocabulary match is inexact. Non-blocking. |

### Non-goals check

- `schemas/` untouched — confirmed.
- No `cn cell` Go implementation changes — confirmed, no `.go` files in the diff.
- No physical `CELL-KINDS.md`→`CELL-DOMAINS.md` rename — confirmed.
- The workflow/SKILL/golden-yml touches are entirely explained by the headline staleness finding
  (pure removals of #639 content, no new content added by the PR's own hand) — not a *separate*
  non-goal violation, but corroborating evidence of the regression's scope.

### Minor non-blocking finding

CELL-RUNTIME.md's runner-routing enum (`continue, hold, escalate, human_gate, no_op_with_reason`)
and its no-op-reason enum (`hold | blocked | waiting_human_gate | already_satisfied |
no_unblocked_cell | evidence_missing`) use overlapping-but-inconsistent terms (`human_gate` vs
`waiting_human_gate`). Worth a follow-up tidy; does not block AC6.

## Verdict

**`verdict: iterate`**

Driven entirely by the headline finding (silent deletion of landed cnos#639 content) — dispositive
on its own. The substantive new doctrine content (CELL-RUNTIME.md + the CCNF/CELL-KINDS/CDD.md/
README.md edits) is independently sound and satisfies AC1–AC8 on the merits. **Required repair:**
rebase `sigma/cell-runtime-arch-note` onto current `main` (restoring §9.13 and the other
#639-derived content) and re-verify no other content was silently dropped, before this cell's
independent β can converge.

**Merge disposition:** per the issue body's explicit framing ("Merge of #629 stays gated: it lands
only on that β receipt"), this receipt does **not** authorize merging PR #629. Findings posted as
a PR review on #629 directly.

## §R1 (repair re-entry — `run_class: repair_pass`)

**Repair context:** κ rebased `sigma/cell-runtime-arch-note` onto current `main` (new head
`68797cf9`, parent `44aa9f84`) and re-applied `status:todo` on 2026-07-10T08:22:17Z, per the
issue comment naming the resolved staleness. This round re-verifies that repair independently
(`.cdd/unreleased/628/REPAIR-PLAN.md`) rather than trusting the claim.

### Independent re-verification of the headline finding

```
$ git merge-base --is-ancestor origin/main origin/sigma/cell-runtime-arch-note
main IS ancestor of the PR branch (exit 0)          # was: NOT an ancestor, in §R0
```

Note: `main`'s tip had advanced one further commit (`b33bc822`, an unrelated automated
`board-map` regeneration of `docs/development/board/{board-data.json,index.html}`) after κ's
rebase point (`44aa9f84`). The two-dot diff (`main..PR-branch`) therefore shows that incidental
board-data delta as noise; the **three-dot** (merge-base) diff is the correct comparison and
shows exactly the 5 originally-intended doctrine files, matching κ's own accounting:

```
$ git diff --stat origin/main...origin/sigma/cell-runtime-arch-note
 docs/architecture/CELL-RUNTIME.md                  | 111 +++++++++++++++++++++
 docs/architecture/README.md                        |   4 +
 src/packages/cnos.cdd/skills/cdd/CDD.md            |   2 +-
 src/packages/cnos.cdd/skills/cdd/CELL-KINDS.md     |  43 +++++---
 .../skills/cdd/COHERENCE-CELL-NORMAL-FORM.md       |   9 +-
 5 files changed, 149 insertions(+), 20 deletions(-)
```

Confirmed present and unchanged on the rebased branch tip (independently, not by trusting the PR
body): `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9.13 "scope-continuation shape
(cnos#639)" (grep hit at its expected line); all 8 `.cdd/unreleased/639/*` files (`git ls-tree`);
all 4 previously-deleted dispatch-substrate files (`cnos-cds-dispatch.yml`, `cds-dispatch/SKILL.md`,
`cds-dispatch.golden.yml`, `dispatch-protocol/SKILL.md`). **Zero false deletions — the headline
finding is resolved.**

### Independent re-verification of AC1–AC8 substance

The 5-file diff scope and net change (`+149/-20`) are identical to what §R0 already reviewed
line-by-line; this round spot-checked (not re-derived) the specific citations §R0's per-AC table
relies on directly against the rebased branch tip, and all are present verbatim:

- AC1/AC2 — `docs/architecture/README.md`'s new "## Runtime" section citing `CELL-RUNTIME.md`;
  `COHERENCE-CELL-NORMAL-FORM.md`'s new "Proposed operating-scale realization" paragraph naming
  WC/PC/CC as classes of the kernel (not a fork) and pointing back to #628.
- AC3 — `CELL-KINDS.md` banner still reads "Legacy Cell-Kind Taxonomy, now Domain Vocabulary" with
  the "Demoted (proposed, #628)" status line; `Mandatory terminal learning section` (#614) and
  `Process self-improvement loop (cnos#640)` (#640) both still present, untouched.
- AC4 — `COHERENCE-CELL-NORMAL-FORM.md`'s `V` diff still reads "`PASS` or `FAIL`" with the
  explicit "`WARN` is **not** a verdict value in the shipped schema" sentence; `schemas/` still
  untouched in the three-dot diff.
- AC5 — `CELL-RUNTIME.md` still present (111 lines) and still cross-linked from README + CCNF +
  CELL-KINDS.
- AC6/AC7/AC8 — unaffected by the rebase (no diff-scope change); §R0's findings carry over
  unmodified, including the non-blocking `human_gate`/`waiting_human_gate` terminology nit (still
  present, still non-blocking).

No new content drift, no additional deletions beyond the 5-file intended scope. §R0's per-AC table
is not re-derived from scratch — REPAIR-PLAN.md's scope-discipline note applies — it is confirmed
still accurate against the rebased tip.

### Verdict

**`verdict: converge`**

The single dispositive blocker from §R0 (silent deletion of landed cnos#639 content via a stale
base) is resolved by κ's rebase, independently re-verified above. AC1–AC8 substance, already
independently sound per §R0, is confirmed unchanged. The non-blocking `human_gate` terminology nit
from §R0 remains open as a non-blocking follow-up, not a convergence condition. This receipt
authorizes merging PR #629 per the issue body's Path-A framing ("Merge of #629 stays gated: it
lands only on that β receipt").
