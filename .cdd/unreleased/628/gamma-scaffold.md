# γ scaffold — cnos#628 (S1 doctrine: CCNF-kernel + WC/PC/CC deployment classes; CM↔V settlement)

## Shape of this cell (read first)

This is **not** a standard γ-scaffolds → α-implements → β-reviews cell. The issue body's own
"Ratification path (CDD receipt / β-independence)" section pins **Path A**: the α-matter is
**already authored** as reconciled PR #629 (`sigma/cell-runtime-arch-note` → `main`), written by
κ (Sigma-as-author) outside the dispatch flow, precisely because κ cannot author its own
independent β-review without violating the α≠β contagion firebreak
(`COHERENCE-CELL.md` §"β Independence"). The dispatch wake's δ therefore:

- does **not** dispatch a fresh α round (§9.3 step 2 is skipped — there is no new implementation
  to produce; re-deriving the doctrine would duplicate, not verify, κ's matter);
- **does** run an independent β review of PR #629 against #628's AC1–AC8, as the wake's own
  cycle/628 matter;
- lands the receipt set (`self-coherence.md`, `beta-review.md`, closeouts) under
  `.cdd/unreleased/628/` on `cycle/628`, distinct from PR #629's own branch.

`run_class: first_pass` (per `cds-dispatch/SKILL.md` §"Repair re-entry preflight" Step A — no
rejection evidence, no scanner "MECHANICAL reversion" comment, no prior-merged-round continuation
evidence on #628 itself).

## Source of truth

| Surface | Source | Status |
|---|---|---|
| Cell contract | issue #628 body (AC1–AC8, Scope, Non-goals) | authoritative |
| α-matter under review | PR #629 (`sigma/cell-runtime-arch-note` → `main`), 1 commit, `MERGEABLE`, all 10 checks green | open |
| Parent wave | #627 | open |
| Prior sibling scope-continuation precedent | `.cdd/unreleased/626/self-coherence.md` §"run_class note"; `gamma-closeout.md` L309–315, L382–385 | landed on `main` |
| δ contract for this shape | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` §9 (no subsection names this exact "review pre-existing external α-matter" shape explicitly; closest analogues are §9.11 resumed-from-mechanical-reversion's "don't re-scaffold existing matter" principle and §9.13 scope-continuation's "don't re-litigate merged rounds" principle — neither is a perfect fit since #628 has no prior *cnos.cdd-dispatch-cycle* history at all, only externally-authored matter) | this scaffold names the gap explicitly rather than silently misapplying §9.10/§9.11/§9.13 |

## Per-AC oracle list (from issue #628)

| AC | Oracle |
|---|---|
| AC1 | Doctrine states the four orthogonal axes; ε named cross-cell (not internal-closure role); CCNF kernel named canonical unit. Grep/citation check against `docs/architecture/CELL-RUNTIME.md`. |
| AC2 | Doctrine states WC/PC/CC are output-telos deployment classes running the full internal role loop (not α/β/γ re-labeled); "PC = relator"; S₃-invariance caveat; selector is `--class`. |
| AC3 | `CELL-KINDS.md` re-headed in place (banner + body) as domain vocabulary; filename preserved; heterogeneity classification + `cell_kind`→`cell_class`+`matter_domain` compat note present; #614's `planning` domain + mandatory learning section, and #640's process-loop, preserved (re-headed, not removed); physical rename deferred; `CDD.md` pointer updated. |
| AC4 | CM↔V settled without changing the shipped `#ValidationVerdict` enum; CM may carry graded/tsc-style measurement + escalation predicates; V emits only `PASS\|FAIL`; escalation is runner routing; override is a δ `#BoundaryDecision` action; `WARN` is not a verdict (CCNF corrected to match). |
| AC5 | `docs/architecture/CELL-RUNTIME.md` present and cross-linked from `docs/architecture/README.md`, `COHERENCE-CELL-NORMAL-FORM.md`, and `CELL-KINDS.md`. |
| AC6 | Escalation-predicate + no-silent-no-op invariants stated. |
| AC7 | #366 reconciliation recorded in the PR; no new phase issue filed until S1 is reviewed. |
| AC8 | Doctrine receipt (PR body) names every settled decision (1–5), the constitutive-change/migration path, the ratification path, and residual questions, with citations. |

## Non-goals to verify absent from the diff

No schemas, no `cn cell` implementation, no wakes, no CDS behavior change, no `CDD.md`
**rewrite** (one-line pointer fix only), no launching S2–S8, no change to `ValidationVerdict`
enum, no physical `CELL-KINDS.md`→`CELL-DOMAINS.md` rename. #614's `planning`/learning and #640's
process-loop content must not be deleted.

## β prompt (this wake's own task)

Independently walk AC1–AC8 above against PR #629's actual diff (`git diff main...sigma/cell-runtime-arch-note`)
and the landed doctrine files it touches. Do not trust the PR body's own claims — verify by
reading the diff and the resulting file content directly (grep/citation, per the issue's own
"Oracles are grep/citation checks against the landed doctrine" line). Confirm the non-goals list
above holds. Record verdict `converge` or `iterate` with per-AC findings in `beta-review.md`.

## Scope guardrails

This cell does not modify PR #629's branch. This cell's own matter is the receipt set under
`.cdd/unreleased/628/` on `cycle/628`. Merge of PR #629 is gated on this receipt's verdict per
the issue body's explicit Path-A framing ("Merge of #629 stays gated: it lands only on that β
receipt (preferred) or an explicit operator human-gate override").

## Friction notes

1. The δ skill's §9 does not yet name a subsection for "review pre-existing, externally-authored
   α-matter, no fresh α dispatch" as a first-class resume shape (distinct from §9.10 rejection-resume,
   §9.11 mechanical-reversion-resume, §9.13 scope-continuation). This scaffold's own reasoning above
   is a first empirical witness for that gap; a future amendment may want to name it (candidate:
   `external-matter-review` shape) the way §9.13 formalized the `scope_continuation` shape after
   cnos#626's self-disclosed gap.
