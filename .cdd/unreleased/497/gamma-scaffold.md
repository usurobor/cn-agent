# γ scaffold — cnos#497

cycle: 497
protocol: cds
mode: design / decision
base_sha: c82750d24381b878c30cf80f09b0ccf4e50494e5
cycle_branch: cycle/497
scaffold_authored_by: gamma@cdd.cnos (δ-collapsed)

---

## Cell summary

**Governing question:** Does `.cdd/` mean "owned by the cnos.cdd package" (Model A) or "the repo's CDD receipt ledger surface" (Model B)?

**Deliverable:** A written decision artifact answering all 7 open questions from the issue body. No code change. No migration. Implementation issues (497B, 497C) file only if Model A wins — and only after this issue's verdict lands.

**Cell mode:** design / decision — docs-only. Per persona commitment §5, γ+α+β-collapsed-on-δ applies. The AC oracle is mechanical (decision artifact exists, all 7 questions answered, cross-references named). Review-independence risk is structurally low.

---

## Source of truth

| Claim / surface | Canonical source | Status |
|---|---|---|
| Current artifact paths | `.cdd/unreleased/{N}/`, `.cdd/releases/{X.Y.Z}/{N}/`, `.cdd/releases/docs/{ISO-date}/{N}/` | In use on main |
| Release mechanics | `src/packages/cnos.cdd/skills/release/SKILL.md` (or current path) | In use |
| Dispatch-protocol artifact root | `dispatch-protocol/SKILL.md §2.2` referencing `.cdd/unreleased/` | On main |
| Log-write locality (cnos#496) | `dispatch-protocol/SKILL.md §2.7`, `AGENT-ACTIVATION-LOG-v0.md §0.1` | On main |
| Label split doctrine | PR #480 (cnos.cdd framework vs cnos.cds protocol) | MERGED |
| Wake-provider artifact output contract | `cnos.cds/orchestrators/cds-dispatch/wake-provider.json` `output_contract.cycle_artifact_root` = `.cdd/unreleased/{N}/` | On main |

---

## Per-AC oracle

| AC | Oracle | Pass criterion |
|---|---|---|
| AC1: Decision artifact exists | File exists at `.cdd/unreleased/497/self-coherence.md` OR a `docs/gamma/decisions/` ADR | `ls .cdd/unreleased/497/self-coherence.md` exits 0 |
| AC2: All 7 open questions answered | Decision artifact answers Q1–Q7 explicitly | grep for "Q1" through "Q7" (or equivalent headings) in the decision artifact |
| AC3: Model chosen and reasoned | Decision artifact names the chosen model with explicit reasoning, not just a preference | The decision names Model A or Model B; cites at least 2 structural arguments |
| AC4: Affected surfaces named | Cross-references identify all surfaces the decision impacts (or explicitly states "none" for Model B) | Cross-reference list exists in the decision artifact |
| AC5: Follow-up issues identified | If Model A: 497B + 497C are defined with crisp scope. If Model B: explicitly states "497B + 497C never file" | Statement present in decision artifact |
| AC6: No implementation work in this cycle | No rename, no migration, no renderer change, no CI guard change in diff | `git diff cycle/497..main -- ':!.cdd/'` exits 0 (no changes outside `.cdd/`) |

---

## Scope guardrails

**In scope:** Decision artifact (text). Cross-references to affected surfaces. Follow-up issue identification.

**Out of scope (hard):** ANY rename of `.cdd/unreleased/{N}/`. ANY migration. ANY renderer change. ANY CI guard change. ANY provider field change.

---

## α prompt

**Context:** You are α for cnos#497. The claimed cell's mode is `design / decision`. You are to produce the decision artifact that answers the 7 open questions from the issue body and declares which model (A or B) the project adopts for artifact-path ownership. This is a docs-only cycle — no code changes.

**Your deliverable:** Write `.cdd/unreleased/497/self-coherence.md` containing:
- §R0: The decision artifact itself — Model A or Model B, with reasoning, and explicit answers to all 7 open questions
- Source of truth table (citing this scaffold)
- AC verification for AC1–AC6

**Constraints:**
- No code changes. No renames. No migrations.
- Answer ALL 7 open questions from the issue body explicitly.
- Name the chosen model clearly: "Verdict: Model B" or "Verdict: Model A".
- If Model B: state "497B + 497C never file; this issue closes with the decision artifact."
- If Model A: state scope of 497B + 497C per the issue's "Follow-up implementation shape" section.
- Cross-reference the specific surfaces the decision impacts (or explicitly "none changed" for Model B).

---

## β prompt

**Context:** You are β for cnos#497. Review α's decision artifact at `.cdd/unreleased/497/self-coherence.md §R0`.

**Your review oracle:**
1. Is the decision artifact present?
2. Are all 7 open questions from the issue body explicitly answered?
3. Is a model (A or B) clearly named with at least 2 structural arguments?
4. Are affected surfaces named (or "none changed" declared)?
5. Are follow-up issues (or non-filing) declared?
6. Is there any diff outside `.cdd/unreleased/497/` in this cycle? (Must be none.)

**Verdict criteria:** `converge` if all 6 pass. `iterate` with specific findings otherwise.

---

## Friction notes

- The cell's mode (`design / decision`) is not one of the four standard CDD modes. It is a specialization of docs-only — the artifact is a decision text rather than code or skill patches.
- The operator has explicitly stated this is design-first: "implementation premature without ledger-ownership decision." Respect this constraint; produce the decision, not the implementation.
- The collapse (γ+α+β onto δ) is justified: AC oracle is mechanical (file exists, headings present, no diff outside `.cdd/`); review-independence risk is low for a decision-text artifact.
