<!-- sections: [gap, skills, acs, self-check, debt, cdd-trace, review-readiness] -->
<!-- completed: [gap, skills, acs, self-check, debt, cdd-trace, review-readiness] -->

# Self-Coherence — Issue #359

## Gap

`operator/SKILL.md §5.2` titled "Single-session δ-as-γ via Agent tool" did not state which role-pair the collapse covers. The intended reading is "δ acts as γ; α and β remain isolated sub-agents." The observed misread (tsc #49 wave-1) is "all roles fused into one sub-agent," which violates `CDD.md §1.4` Triadic rule. §5.2 needed an explicit scope statement and an explicit violation shape.

Mode: docs-only. Branch: `cycle/359`. Base: `23e28e45`.

## Skills

- Tier 1 / role: `CDD.md`, `alpha/SKILL.md`. Lifecycle sub-skills not loaded — single-paragraph skill patch with three ACs needed no design or plan artifact.
- Tier 2: none — change is not engineering code.
- Tier 3: `src/packages/cnos.core/skills/write/SKILL.md` (governing question, brevity-is-earned, state-what-it-is, front-load-the-point).

## ACs

Implementation SHA: `22e9e7eb` (`α-359: clarify §5.2 collapses δ↔γ only — γ↔α↔β stays separate`). Diff: `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md +4 / -0`.

**AC1 — `operator/SKILL.md §5.2` contains explicit statement that §5.2 collapses δ↔γ only and γ↔α↔β remain structurally separate per §1.4.** Evidence: the new paragraph in `operator/SKILL.md §5.2` reads "§5.2 collapses **δ↔γ only**. γ↔α↔β remain structurally separate per `CDD.md §1.4` Triadic rule: γ scaffolds and coordinates in the parent session, α implements in its own sub-agent, β reviews and merges in its own sub-agent." `git show 22e9e7eb -- src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` confirms the line lands at the head of the three-consequences list. Met.

**AC2 — names what violation looks like (single subagent doing γ+α+β work).** Evidence: the second new paragraph reads "A single sub-agent that performs γ-selection plus α-implementation plus β-review is not §5.2 — it is a §1.4 violation. §5.2 requires three execution contexts: the parent session (γ, also δ), a separate α sub-agent, and a separate β sub-agent. Lumping γ+α+β into one sub-agent breaks role-isolation (α gains access to β's reasoning and vice versa) and is rejected." Met.

**AC3 — `operator/SKILL.md §5.2` consistent with the patch.** Evidence: re-read §5.2 mechanism paragraph ("α and β are dispatched as sub-agents") and consequence (1) ("δ=γ collapse"). The new "Scope of the collapse" block does not contradict these — it sharpens the mechanism paragraph (which already implied α/β separation) and explains the *scope* that consequence (1) is internal to. Downstream references to §5.2 in `activation/SKILL.md §8` (dispatch declaration), `release/SKILL.md §3.8` (configuration-floor clause), and `operator/SKILL.md §5.3` (escalation criteria) all describe "γ/δ separation absent" — none asserts α/β collapse, so the patch leaves them coherent. `rg "§5\.2|δ-as-γ|δ=γ|δ↔γ" src/packages/cnos.cdd/skills/cdd/` enumerated and inspected: 4 files match, all consistent. Met.

## Self-check

**Did α push ambiguity onto β?** No. The patch adds two paragraphs of explicit scope language and an explicit violation shape; the diff is bounded (operator/SKILL.md +4 lines), the wording is positive (states what §5.2 *is*, not just what it is not), and the three pre-existing structural consequences below the new block remain unchanged — β does not need to reverse-engineer intent.

**Is every claim backed by evidence in the diff?** Yes. AC1 / AC2 quote the exact paragraphs added to operator/SKILL.md §5.2. AC3 lists the four downstream sites that reference §5.2 (`operator/SKILL.md`, `release/SKILL.md`, `alpha/SKILL.md`, `activation/SKILL.md`) and notes that none assert α/β collapse, so the patch leaves them coherent. The `rg` enumeration in AC3 is reproducible.

**Skill-class peer enumeration.** The diff touches one role skill (`operator/SKILL.md`). Role-skill peers (`alpha/`, `beta/`, `gamma/`) are exempt: this clarification is δ-scoped and does not change the γ/α/β contract — γ↔α↔β separation is *reaffirmed*, not redefined. Lifecycle-skill peers (`review/`, `release/`, `post-release/`, `design/`, `plan/`, `issue/`) are exempt: none encodes §5.2's collapse scope as an operational contract; `release/SKILL.md §3.8` references §5.2 only as a grading-floor trigger ("γ/δ separation absent"), which the patch leaves true.

**Polyglot re-audit.** Diff languages: Markdown only (operator/SKILL.md, self-coherence.md, gamma-scaffold.md). Markdown checks: table shape unchanged, no new cross-references introduced (the new paragraphs reference `CDD.md §1.4` which already exists in §5.2's surrounding prose), prose reads coherently with neighboring "Three structural consequences follow" list.

**Artifact enumeration matches diff.** `git diff --stat origin/main..HEAD` lists three files: `operator/SKILL.md` (the patch — AC1/AC2/AC3 evidence), `gamma-scaffold.md` (γ artifact, not authored by α — referenced in CDD-Trace step 6), `self-coherence.md` (this file). All three accounted for.

## Debt

None. The patch is a contained clarification: two paragraphs at the head of §5.2's "Three structural consequences" list. No deferred work, no partial AC closure, no untouched peers. The pre-existing TSC #49 wave-1 misread was the originating signal and is closed by the new violation-shape paragraph (AC2).

## CDD Trace

1. **Design** — not required. Issue body, ACs, and the recommended patch are concrete enough to author directly; no impact graph beyond §5.2 itself and no new contract surface. Documented as "not required" per `alpha/SKILL.md §2.2`.
2. **Coherence contract** — §Gap above names the incoherence ("§5.2 did not state which role-pair the collapse covers") and the originating misread (tsc #49 wave-1, recorded in `.cdd/releases/0.10.0/49/cdd-iteration.md` F1).
3. **Plan** — not required. Single-file edit, three ACs, ordering trivial: write paragraph 1 (scope statement, AC1), write paragraph 2 (violation shape, AC2), re-grep §5.2 references to verify consistency (AC3). Documented as "not required."
4. **Tests** — not applicable. Skill content is prose-only; no executable surface. The acceptance check is the §ACs evidence in this artifact (text presence + downstream-reference enumeration).
5. **Code** — not applicable (no executable code touched). The implementation is the docs change in step 6.
6. **Docs** — `src/packages/cnos.cdd/skills/cdd/operator/SKILL.md` patched at §5.2 head (commit `22e9e7eb`, +4 lines). γ artifact `.cdd/unreleased/359/gamma-scaffold.md` (commit `88a573de`, authored by γ, not α). Caller-path trace: not applicable — no new modules or functions, only inline prose added to an existing section. Downstream §5.2 references in `release/SKILL.md §3.8`, `alpha/SKILL.md` (review-readiness fix-round protocol), and `activation/SKILL.md §8` (dispatch declaration) inspected and left coherent (AC3 evidence).
7. **Self-coherence** — this file. Sections committed incrementally: `f6672240` §Gap, `c7977e21` §Skills, `9b4c71c0` §ACs, `039dcc7a` §Self-check, `057150c1` §Debt, this commit §CDD-Trace. Review-readiness signal follows in a separate commit per `alpha/SKILL.md §2.6` SHA-convention guidance.

## Review-readiness | round 1

- **Implementation SHA:** `22e9e7eb` (`α-359: clarify §5.2 collapses δ↔γ only — γ↔α↔β stays separate`)
- **Base SHA:** `23e28e45` — equals `origin/main` HEAD observed 2026-05-14 at append time (no rebase needed; cycle branch is rooted on current main)
- **Branch CI:** green — Build `success` on HEAD `f7f79415` per `gh run list --branch cycle/359` at 2026-05-14 20:12:02Z
- **Author email:** `alpha@cdd.cnos` on all 8 α commits per `git log --format='%ae' origin/main..HEAD` (γ scaffold commit `88a573de` is `gamma@cdd.cnos`, expected)
- **Pre-review gate (`alpha/SKILL.md §2.6`):**
  1. cycle branch rebased on main ✓ (base = current main HEAD)
  2. CDD Trace through step 7 ✓
  3. tests N/A — skill content is prose-only ✓ (§ACs documents reason)
  4. every AC has evidence ✓ (§ACs)
  5. known debt explicit ✓ ("None" in §Debt)
  6. schema / shape audit N/A — no contract surface changed ✓
  7. peer enumeration ✓ (§Self-check — role-skill peers exempt, lifecycle-skill peers exempt with reason)
  8. harness audit N/A — no schema-bearing contract ✓
  9. polyglot re-audit ✓ — Markdown-only diff (§Self-check)
  10. CI green on HEAD ✓
  11. artifact enumeration matches diff ✓ (§Self-check: 3 files in `git diff --stat`, all accounted for)
  12. caller-path trace N/A — no new modules ✓
  13. test assertion count N/A — no tests ✓
  14. author email matches `alpha@cdd.cnos` ✓
- **Status:** ready for β.
