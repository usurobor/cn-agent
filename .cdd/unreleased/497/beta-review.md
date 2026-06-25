# beta-review — cnos#497

cycle: 497
role: beta
round: R0

---

## §R0

### Pre-merge gate walk

**Row 1 — AC oracle (per gamma-scaffold.md):**

| AC | Oracle result | Notes |
|---|---|---|
| AC1: Decision artifact exists | PASS | `.cdd/unreleased/497/self-coherence.md` present |
| AC2: All 7 open questions answered | PASS | Q1–Q7 answered under explicit headings in §R0 |
| AC3: Model named with ≥2 structural arguments | PASS | "Verdict: Model B"; 5 structural arguments present |
| AC4: Affected surfaces named | PASS | All major surfaces listed and confirmed unchanged |
| AC5: Follow-up issues identified | PASS | "497B and 497C never file" — explicit declaration |
| AC6: No implementation diff outside `.cdd/` | PASS | Only `.cdd/unreleased/497/` changed in cycle/497 |

**Row 2 — Canonical skill staleness:** γ scaffold authored at SHA `c82750d`; no skill-affecting commits on main since claim. Not stale.

**Row 3 — Source-of-truth alignment:** Decision cites `dispatch-protocol/SKILL.md`, `wake-provider.json output_contract`, and cnos#496 sibling context. All references are present on main. No broken or stale citations.

**Row 4 — Independent AC walk:**

1. **Q1 ("What owns the receipt ledger?")** — α's answer is internally consistent and grounded. The ledger-vs-namespace distinction argument is sound: `.cdd/` is the framework's accounting surface; concrete protocols are writers, not owners. The analogy (`.git/` as git's metadata directory) is apt. No gap.

2. **Q2 ("What does release do with `.cds/unreleased/`?")** — Correctly answered: that path does not exist under Model B; release reads `.cdd/unreleased/` as today. No implementation change required.

3. **Q3 ("Does `.cdd/releases/` still exist?")** — Yes, confirmed unchanged. Concise and correct.

4. **Q4 ("Where do docs-only disconnects move?")** — Correctly stated as unchanged under Model B: `.cdd/releases/docs/{ISO-date}/{N}/`.

5. **Q5 ("What does cn-cdd-verify verify?")** — Correct: inspects `.cdd/` as unified ledger; no rename; no per-package variant. Correctly notes that protocol-specific routing (if ever needed) would be a runtime value inside the receipt, not a directory variable.

6. **Q6 ("How does a package-agnostic parent find child receipts?")** — Single walk over `.cdd/unreleased/` confirmed. Correctly cites wave-master tracking already in use (cnos#467 sub-issue pattern). Solid.

7. **Q7 ("Mechanical-orchestration trajectory?")** — Correct: `cn dispatch claim` writes to `.cdd/unreleased/{N}/` always; path is a framework constant; protocol label drives dispatch routing but not path derivation. The orthogonality claim is well-supported.

**Row 5 — No implementation work in diff:** Confirmed. Only `.cdd/unreleased/497/` artifacts created. No code, no renames, no migrations.

**Row 6 — Reasoning quality check:**
- Five structural arguments presented; each is independent (ledger-vs-namespace, write-locality orthogonality, migration cost, O(1) resolution, mechanical-orchestration simplicity).
- No argument is a preference ("Model B feels better") — all arguments are structural (tooling invariance, information-theoretic sufficiency of in-receipt protocol fields, complexity reduction).
- Model A's case is not strawmanned; the issue body's honest presentation of both models is preserved in the decision.

**Row 7 — Scope discipline:** Decision artifact correctly stays within the declared scope: "The decision is text. The implementation is separate." No creep toward renaming, migration, or CI changes.

---

### Verdict

`verdict: converge`

All 6 ACs pass. No findings requiring iteration. The decision artifact is complete, internally coherent, and consistent with the issue's declared scope. The reasoning is grounded in structural arguments, not preference. All 7 open questions are answered explicitly.

**Merge recommendation:** Merge cycle/497 → main. Close cnos#497 with the decision artifact as the deliverable. 497B and 497C do not file per the verdict.
