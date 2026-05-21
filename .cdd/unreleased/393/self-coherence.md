# α self-coherence — cycle/393

**Issue:** cnos#393
**Mode:** design-and-build, γ+α+β-collapsed-on-δ
**Branch:** `cycle/393` from `origin/main` base `e531dba0`

## §Gap

cnos#393 names a **protocol gap** empirically anchored in cnos#389 (Python-not-Go) and cnos#391 (wrong package scoping + separate binary): the CDD role-skill surfaces did not name a role that owns implementation-contract decisions (language, CLI integration target, package scoping, runtime dependencies, existing-binary disposition, JSON/wire contract preservation, backward-compat invariants). δ did it implicitly, but the responsibility was not written down; α had room to improvise; β's behavior-only AC oracles APPROVE-d without catching the implementation-contract drift; γ's dispatch prompts under-specified the contract.

cnos#392 was the first cycle where δ pinned the 7-axis implementation contract at dispatch as an ad-hoc operator action. Its `cdd-iteration.md` F1–F4 forecast the four patches cnos#393 ships. This cycle makes the doctrine explicit across the four role surfaces (α, β, γ, operator/δ) and surfaces it as a Phase 4 (δ split) design input on cnos#366.

§Gap is empirical, not negative. `rg "Implementation contract" src/packages/cnos.cdd/skills/cdd/{alpha,beta,gamma,operator}/SKILL.md` returned 0 hits before this cycle's commits (verified at base SHA `e531dba0`). The 4-patch surface is genuinely additive; no rule renumbering of existing α/β/γ/operator content; no rewrite of existing γ §2.5 template rows; no removal of any existing operator section.

## §Skills

**Tier 1 (canonical):** `CDD.md`, `cdd/alpha/SKILL.md`, `cdd/beta/SKILL.md`, `cdd/gamma/SKILL.md`, `cdd/operator/SKILL.md`, `cdd/design/SKILL.md`, `cdd/issue/SKILL.md` proof half.

**Tier 2 (always-applicable):** none — this cycle is markdown-only; no `eng/*` skills load.

**Tier 3 (cycle-specific):** the issue named all four patch-target skills as Tier 3 to load (α, β, γ, operator) — all read in full before authoring. `cnos.cdd/skills/cdd/COHERENCE-CELL-NORMAL-FORM.md` (read for δ-as-boundary framing; the inward face is the additive extension), `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` (read for "Separate persona, protocol, and project" framing), `ROLES.md §4a` (five-layer chain — δ's two-sided membrane fits layer 2 operator-contract surface).

## §ACs

| AC | Status | Evidence |
|---|---|---|
| **AC1** — `alpha/SKILL.md` Rule 8/equivalent lands; `rg "α MUST NOT change the implementation-contract\|Implementation contract is δ"` ≥1 hit; cites #389 + #391 | **PASS** | `alpha/SKILL.md` §3.6 "Implementation contract is δ's, not α's" landed at commit `a49de338` (lines 344–365 of current file). `rg "α MUST NOT change the implementation-contract\|Implementation contract is δ"` returns 2 hits (heading + body). Empirical anchor paragraph cites both cnos#389 and cnos#391 explicitly. Numbering note: α's `## 3. Rules` uses 3.1–3.5 single-period style (not "Rule N" as β does); new rule lands as 3.6 to match α convention. |
| **AC2** — `beta/SKILL.md` Rule 7 lands; `rg "Implementation-contract coherence\|Rule 7"` ≥1 hit; cites #389 R1 + #391 R1 | **PASS** | `beta/SKILL.md` §Role Rules `### 7. Implementation-contract coherence` landed at commit `a49de338`. `rg "Implementation-contract coherence\|Rule 7"` returns 3 hits (heading + 2 body refs). Empirical anchor explicitly cites cnos#389 R1 + cnos#391 R1 behavior-only APPROVE failures. |
| **AC3** — `gamma/SKILL.md` §2.5 dispatch-prompt template includes `## Implementation contract`; enumerates 7 axes; rule "γ MUST NOT dispatch with empty rows" named | **PASS** | `gamma/SKILL.md` §2.5 Step 3b new subsection `##### Implementation contract section (required for α prompt)` landed at commit `a49de338`. The template block enumerates all 7 axes: Language, CLI integration target, Package scoping, Existing-binary disposition, Runtime dependencies, JSON/wire contract preservation, Backward-compat invariant. The rule "**γ MUST NOT dispatch with empty / 'TBD' rows**" is named explicitly in bold. |
| **AC4** — `operator/SKILL.md` names δ-inward membrane; `rg "inward membrane\|implementation-contract enrichment"` ≥1 hit; names δ two-sided; references γ template | **PASS** | `operator/SKILL.md` §3a `## 3a. δ as inward membrane: implementation-contract enrichment at dispatch` landed at commit `a49de338`. `rg "inward membrane\|implementation-contract enrichment"` returns 7 hits in operator/SKILL.md. The section explicitly states "δ is also a **two-sided membrane**" and references `gamma/SKILL.md` §2.5 Step 3b (the Patch C template) by name + path. |
| **AC5** — cnos#366 Phase 4 body updated to name δ as explicitly two-sided membrane; absorbs cnos#393 as Phase 4 input | **PASS** | cnos#366 body updated via `mcp__github__issue_write method=update` (REST API write recorded at id `4452434047`). Phase 4 section now opens with "δ is **two-sided** — **outward** ... **AND inward** ... Phase 4 must carve both surfaces explicitly." Phase 4 `delta/SKILL.md` deliverable bullet now reads "boundary complex with two-sided membrane: **outward** ... **AND inward** ..." A new "Update — 2026-05-21 (cnos#393 close-out)" callout sits at the top of the issue body. cnos#393 appears as the 5th Phase 4 input bullet. |
| **AC6** — Cross-references coherent: α cites γ template; β cites α rule + γ template; γ template cites δ scope; δ scope cites γ template | **PASS** | Mesh bidirectional (see §Cross-reference mesh evidence below). Each of the 4 patches cites the other three role-skill files by name + path + §reference. α §3.6 cites γ §2.5 Step 3b, β Rule 7, δ §3a. β Rule 7 cites α §3.6, γ §2.5 Step 3b, δ §3a. γ §2.5 cites α §3.6, β Rule 7, δ §3a. δ §3a cites α §3.6, β Rule 7, γ §2.5 Step 3b. |
| **AC7** — Each patch cites cnos#389 + cnos#391 as empirical anchors | **PASS** | All 4 patches contain a "**Empirical anchor.**" paragraph naming both #389 and #391 by issue number. Mechanical: `grep -c "#389" src/packages/cnos.cdd/skills/cdd/{alpha,beta,gamma,operator}/SKILL.md` returns 1, 1, 1, 2 (operator has 2 because it appears in both the empirical anchor and Phase 4 cross-reference). Same for #391. Each patch also cites cnos#392 as the proof-of-concept cycle where the 7-axis pinning succeeded. |

All 7 ACs PASS.

## §Self-check

**Did α's work push ambiguity onto β?** No. The 4 patches are surgically additive — each insertion point preserves existing rule numbering and existing template rows. The cross-reference mesh is concrete (each citation names the other patch's section heading + file path), not handwavy. AC oracles are mechanically grep-able as authored.

**Is every claim backed by evidence in the diff?** Yes. Each AC row above names the specific section, commit SHA, and mechanical oracle. AC5 is on a GitHub issue body (REST write id `4452434047`), not the branch — confirmed by reading the issue via `mcp__github__issue_read` after the update.

**Implementation contract conformance (self-applying Rule 7).** The implementation contract for this cycle was operator-pinned in the dispatch prompt:

| Axis | Pinned value | Conformance evidence |
|---|---|---|
| Language | Markdown (skill files with YAML frontmatter) | All 4 patches are markdown additions to existing `.md` files; no code changes. |
| CLI integration target | N/A | No commands added. |
| Package scoping | `src/packages/cnos.cdd/skills/cdd/{alpha,beta,gamma,operator}/SKILL.md` | All 4 patches land at exactly those paths. `git diff --stat origin/main..HEAD -- src/` shows exactly the 4 expected files. |
| Existing-binary disposition | N/A | No binaries. |
| Runtime dependencies | None | No new imports / modules / deps. |
| JSON/wire contract | N/A | No schemas / serialization touched. |
| Backward compat | Existing rules (α 3.1–3.5; β Rules 1–6; γ §2.5 template; operator existing sections) preserved; new content additive | Verified by `git diff origin/main..HEAD -- src/packages/cnos.cdd/skills/cdd/`: every hunk is an addition; no deletions of existing rule text; α §3.5 preserved unchanged (the patch inserts §3.6 *after* §3.5); β Rule 6 preserved unchanged (patch inserts Rule 7 after); γ §2.5 existing α/β prompt blocks preserved unchanged (patch inserts the new template subsection *between* them); operator §3.5 preserved (patch inserts §3a *after* §3.5). |

**Cross-reference mesh — is it actually bidirectional?** Each of the 4 patches names the *other three* role-skill files explicitly. Self-references are 0 (a patch does not cite itself; that would be circular). Confirmed via `grep -c <path>` for each pair. No circular dependency in the rule logic — each rule is locally self-justifying via the empirical anchors (#389, #391); the mesh is for discoverability, not for logical justification.

**Did α run AC oracles before signaling review-readiness?** Yes — all AC1–AC4 + AC7 oracle greps were run against the diff before commit `a49de338`; AC5 was verified by re-reading the cnos#366 body after the issue_write. AC6 mesh evidence is in §Cross-reference mesh evidence below.

## §Debt

None of substance. The Patch D content lives in `operator/SKILL.md` §3a today; per the issue's design constraint and the cnos#366 Phase 4 absorption, this section is intentionally a design-prerequisite anchor that Phase 4 will relocate to `delta/SKILL.md`. This is **declared next-cycle work**, not undisclosed debt — cnos#366 now names it explicitly as a Phase 4 input.

One minor authorial choice worth flagging: the issue body suggested "Rule 8" for α; α's actual numbering convention is `3.N` (one-period), not "Rule N" as β uses. I added §3.6 to match α's style; cross-references in β/γ/operator use "α §3.6" not "α Rule 8." This is a faithful translation, not a deviation.

## §Cross-reference mesh evidence (AC6 detail)

Each patch cites the other three by file path + section ref:

```
α §3.6 cites:
  - gamma/SKILL.md §2.5 Step 3b `## Implementation contract` template (line 354)
  - operator/SKILL.md §3a "δ as inward membrane" (line 355)
  - beta/SKILL.md §Role Rules Rule 7 "Implementation-contract coherence" (line 356)

β Rule 7 cites:
  - gamma/SKILL.md §2.5 Step 3b template (lines 165, 174)
  - alpha/SKILL.md §3.6 "Implementation contract is δ's, not α's" (line 173)
  - operator/SKILL.md §3a (line 175)

γ §2.5 new subsection cites:
  - operator/SKILL.md §3a (lines 363, 370)
  - alpha/SKILL.md §3.6 (line 372)
  - beta/SKILL.md §Role Rules Rule 7 (line 374)

operator §3a cites:
  - gamma/SKILL.md §2.5 Step 3b (lines 264, 283)
  - alpha/SKILL.md §3.6 (line 288)
  - beta/SKILL.md §Role Rules Rule 7 (line 291)
```

Total mesh edges: 12 (= 4 patches × 3 other-patches-cited). No circular logical dependencies — each rule's empirical anchor (#389 + #391) is the local justification surface; the mesh is for discoverability.

## §CDD Trace

1. **Receive** — dispatched as δ-as-agent for cycle #393 with operator-pinned implementation contract on the 7 axes. Read issue body via `mcp__github__issue_read`. Mode: design-and-build, γ+α+β-collapsed-on-δ.
2. **Produce — design half** — `gamma-scaffold.md` first (γ-side pre-dispatch gate; binding rule per `gamma/SKILL.md` §2.5 Step 3b), then `design-notes.md` (insertion-point + content sketch per patch). Both committed and pushed to `origin/cycle/393` before any source edit.
3. **Produce — build half** — 4 patches authored at the planned insertion points in one editing pass: alpha §3.6, beta Rule 7, gamma §2.5 subsection, operator §3a. All edits use the `Edit` tool with the `old_string` chosen to preserve adjacent unchanged content.
4. **Prove** — AC oracle greps run against the diff before commit (AC1–AC4 + AC7). Mesh greps confirm AC6. AC5 verified by reading cnos#366 body after the update.
5. **Commit** — single commit `a49de338` for the 4 patches with a structured message naming each patch and the closures.
6. **β-collapsed review** — see `.cdd/unreleased/393/beta-review.md` (β round 1, written under collapsed roles; verifies AC1–AC7 against the diff with particular rigor on AC6 mesh and AC7 anchors).
7. **Close-outs** — `.cdd/unreleased/393/{alpha,beta,gamma}-closeout.md` plus `cdd-iteration.md` plus INDEX update.

Artifact list mapped to diff:

```
.cdd/unreleased/393/gamma-scaffold.md   (γ-side pre-dispatch artifact)
.cdd/unreleased/393/design-notes.md      (α design half)
src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md       (Patch A)
src/packages/cnos.cdd/skills/cdd/beta/SKILL.md        (Patch B)
src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md       (Patch C)
src/packages/cnos.cdd/skills/cdd/operator/SKILL.md    (Patch D)
.cdd/unreleased/393/self-coherence.md    (this file)
.cdd/unreleased/393/beta-review.md       (β verdict, collapsed)
.cdd/unreleased/393/alpha-closeout.md    (α close-out)
.cdd/unreleased/393/beta-closeout.md     (β close-out)
.cdd/unreleased/393/gamma-closeout.md    (γ close-out + triage)
.cdd/unreleased/393/cdd-iteration.md     (closure-gate artifact)
.cdd/iterations/INDEX.md                  (row added for #393)
```

Plus cnos#366 body edit (off-branch, GitHub-side).

## §Review-readiness

Round 1 | base SHA `e531dba0` | implementation SHA `a49de338` | branch CI: skipped (markdown-only diff; no CI gates fire on `.md`-only changes) | ready for β-collapsed review.

Path (b) on γ-artifact-of-record per `alpha/SKILL.md` §2.6 row 15: **γ-artifact at canonical §5.1 path** — `gamma-scaffold.md` committed at `182337cb` and present on `origin/cycle/393` at the canonical `.cdd/unreleased/393/gamma-scaffold.md` location.

Identity: commits authored as `alpha@cdd.cnos` (build commit + this self-coherence commit); scaffold authored as `gamma@cdd.cnos`. Mixed-identity history is correct for the collapsed-on-δ mode: each commit carries the identity of the role that wrote it.
