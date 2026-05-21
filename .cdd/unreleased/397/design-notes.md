# α design notes — cycle/397

## §1 Inventory of `operator/SKILL.md`

Reading the file top to bottom and classifying each section under one of four headings:

- **δR** — δ-as-role (move to `delta/SKILL.md`)
- **OC** — operator-as-coordinator (stays in `operator/SKILL.md`)
- **HM** — harness mechanics (stays; Phase 4b will claim)
- **RE** — release effector mechanics (stays; Phase 4c will claim)

| Lines | Section | Class | Disposition |
|---|---|---|---|
| 1–28 | Frontmatter | OC+δR (operator skill metadata, points at γ) | KEEP. Adjust `description` minimally to note δ-role split. |
| 30–41 | Core Principle ("δ owns what agents cannot") | OC+δR | KEEP as the operator-as-coordinator framing; note `delta/SKILL.md` for the role-skill surface. |
| 42–62 | Algorithm (dispatch γ→α→β→re-dispatch→Gate→Override) | OC | KEEP. The dispatch loop is the coordinator's algorithm. |
| 65–88 | Git identity for role actors | OC | KEEP. Identity convention is project-wide; multiple skills cite it here. |
| 92–121 | §1 Route (dispatch γ first; dispatch α/β sequentially) | OC | KEEP. |
| 124–177 | §2 Wait (do not poll; subscribe to issue) | OC | KEEP. Coordinator's wait-discipline. |
| 180–197 | §3.1 External actions the operator holds (table of actions) | δR (authority naming) | MOVE the **role-authority statement** (δ holds these external actions per AC7's "authority-naming context") to `delta/SKILL.md`. The table itself is δ-role policy. |
| 198–212 | §3.2/§3.3 Execute on request, not observation; Report completion | δR | MOVE to `delta/SKILL.md` (gate-execution discipline = boundary policy). |
| 214–240 | §3.4 Cut the release — disconnect (algorithm + `scripts/release.sh` mechanics) | δR + RE | SPLIT. The role-policy ("δ is sole tag-author; δ blocks release until CI green; δ owns recovery on red"; "the triad's work is not complete until tagged") is δR — moves to `delta/SKILL.md` (authority-naming). The mechanics (`scripts/release.sh` step-by-step, polling commands, branch cleanup runbook) are RE — stays in `operator/SKILL.md` until Phase 4c claims them. The split is clean because the role-policy is one paragraph; the mechanics are several. |
| 242–246 | §3.5 The tag is the signal | δR | MOVE (this is the boundary-decision-as-git-observable doctrine). |
| 250–307 | §3a δ as inward membrane | δR | MOVE in its entirety. Per AC3 the substance MUST appear only in `delta/SKILL.md`; `operator/SKILL.md §3a` replaced by a one-line redirect. |
| 311–339 | §4 Override (when, protocol, not for taste) | δR | MOVE entirely. Override is the degraded boundary action per cnos#367 freeze; it is the canonical δ-as-role authority surface. |
| 343–412 | §5 Dispatch configurations (§5.1 multi-session; §5.2 δ=γ collapse; §5.2.1 quiescence; §5.3 escalation) | OC+HM | KEEP. Dispatch config + collapse rules + harness quiescence are coordinator + harness concerns. Phase 4b will claim §5.2.1 (quiescence is harness mechanics) and possibly more. |
| 416–425 | §6 What the operator does NOT do | OC | KEEP (these are negative coordinator rules — what δ-as-coordinator must not do). |
| 429–444 | §7 Cycle lifecycle from the operator's view | OC | KEEP (lifecycle table — operator coordination). |
| 447–497 | §8 Timeout recovery | HM | KEEP. Recovery procedure is harness-mechanics-adjacent; Phase 4b claims later. |
| 501–547 | §9 Embedded Kata (Kata A normal cycle; Kata B override) | OC+δR | KEEP both katas. Kata A exercises coordinator routing (OC). Kata B is an override exercise (δR-adjacent); since the override doctrine moves to `delta/SKILL.md`, the kata stays here as the operator-as-coordinator's "when override is requested" exercise, with a cross-reference to `delta/SKILL.md` §Override for the doctrine. Alternative: also lift Kata B into `delta/`. Decision: leave both katas in `operator/SKILL.md` to keep the kata bundle intact; `delta/SKILL.md` provides the override doctrine; the kata references it. This avoids splitting kata-bundle pedagogy. |
| 551–737 | §10 Wave Coordination (§10.1 template; §10.2 manifest; §10.3 status; §10.4 closure; §10.5 reporting; §10.6 iteration lifecycle) | OC+HM | KEEP. Wave coordination is coordinator/harness; not δ-role boundary policy. |

**Sections that move (verbatim, no content drops):**
1. §3.1 (external-actions table + framing as δ-held actions)
2. §3.2 (execute on request not observation)
3. §3.3 (report completion)
4. §3.4 **role-policy paragraphs** (δ-as-sole-tag-author; CI-green/red dispatch authority; "the tag is the disconnection point" framing) — NOT the `scripts/release.sh` step list
5. §3.5 (the tag is the signal)
6. §3a (δ as inward membrane — entire section)
7. §4 (override — when / protocol / not for taste)

**Sections that stay in `operator/SKILL.md`:**
- Frontmatter, Core Principle, Git identity table
- §1 Route, §2 Wait
- §3.4 **mechanics paragraphs** (`scripts/release.sh` invocation; manual-tag prohibition; tag-message-generation runbook) — gets a redirect note pointing to `delta/SKILL.md` for the boundary-policy framing
- §5 Dispatch configurations
- §6 What the operator does NOT do
- §7 Cycle lifecycle
- §8 Timeout recovery
- §9 Embedded Kata (with delta/ cross-reference for the override doctrine in Kata B)
- §10 Wave Coordination

## §2 Plan for `delta/SKILL.md` structure

**Frontmatter** — name: delta; artifact_class: skill; parent: cdd; scope: role-local; triggers covering boundary, membrane, override, receipt, verdict, gate; calls operator/SKILL.md and RECEIPT-VALIDATION.md.

**Section plan:**

1. **Preamble + Core Principle** — δ as the cell's boundary actor per CCNF; two-sided membrane framing (outward/inward); relation to operator-as-coordinator (δ-the-role lives here, dispatch-coordination of δ-the-actor stays in `operator/SKILL.md`).

2. **§1 Outward membrane — receipt + V verdict → boundary decision**
   - Subsections lifted from operator §3.1, §3.2, §3.3, §3.5, plus the §3.4 role-policy paragraphs.
   - The five `BoundaryDecision` enumerants per CCNF §Cell Outcomes (accept, release, override, reject, repair_dispatch).
   - Authority-naming context for external actions (tag, release, deploy, push, branch-delete, force-push, auth-refresh) — names them; mechanics stay in `operator/SKILL.md` per AC7 / Phase 4c.
   - Composition rule with V (RECEIPT-VALIDATION.md §Validation Interface): δ reads `ValidationVerdict`, records `BoundaryDecision`.
   - Cross-references: V (RECEIPT-VALIDATION.md), CCNF (COHERENCE-CELL-NORMAL-FORM.md), operator/SKILL.md (for dispatch-coordination + release mechanics).

3. **§2 Inward membrane — γ contract → α-ready dispatch (implementation-contract enrichment)**
   - Lifted verbatim from operator §3a.
   - The 7 implementation-contract axes.
   - δ's review-before-routing duty; fill-or-escalate paths.
   - The four-surface mesh (γ template / δ enrichment / α constraint / β verification).
   - Phase 4a notation: this section now lives here (no longer a "Phase 4 relocation target" — it has been relocated).
   - Empirical anchors (cnos#389, #391, #392, #393) preserved.

4. **§3 Override — degraded boundary action**
   - Lifted from operator §4.
   - **Critical addition per AC4:** explicit verdict-vs-decision distinction. Override does not rewrite the `ValidationVerdict`; override is a `BoundaryDecision`. Cite `RECEIPT-VALIDATION.md` §Q4 and §"ValidationVerdict vs BoundaryDecision" verbatim where the freeze lives.
   - When to override, how to declare, not-for-taste constraint.
   - Cross-reference to operator/SKILL.md §9 Kata B for the embedded exercise.

5. **§4 Composition with V (RECEIPT-VALIDATION.md)**
   - The composition rule from `RECEIPT-VALIDATION.md` §Validation Interface (PASS → {accept,release} or reject/repair-dispatch; FAIL → {reject, repair-dispatch} or override-with-block; PASS+override = invalid; FAIL+accept = invalid).
   - The four cell outcomes from CCNF §Cell Outcomes (accepted, degraded, blocked, invalid).
   - Cross-reference both surfaces.

6. **§5 What δ-as-role does NOT do**
   - Mirror operator §6 from the role perspective: δ does not produce matter (α), does not discriminate matter (β), does not close cells (γ).
   - δ-the-role lives here; δ-as-actor's coordinator hat lives in `operator/SKILL.md` (γ=δ collapse case).

7. **§6 Cross-references and relationships**
   - `operator/SKILL.md` for dispatch coordination (γ=δ collapse, harness mechanics)
   - `RECEIPT-VALIDATION.md` for V verdict / BoundaryDecision freeze
   - `COHERENCE-CELL-NORMAL-FORM.md` for δ's signature in the kernel
   - `COHERENCE-CELL.md` for predecessor doctrine (δ Boundary Complex)
   - `gamma/SKILL.md` §2.5 for the γ-template side of the mesh
   - `alpha/SKILL.md` §3.6 for α-constraint
   - `beta/SKILL.md` Rule 7 for β-verification
   - `release/SKILL.md` and `post-release/SKILL.md` for release-mechanics (Phase 4c claim)

## §3 Cross-reference update plan

Three skill files cite `operator/SKILL.md §3a` (the inward-membrane content) — these MUST update to `delta/SKILL.md`:

1. `gamma/SKILL.md` line 363, line 370 — two refs to `operator/SKILL.md §3a "δ as inward membrane"`. Both update to `delta/SKILL.md §2 "Inward membrane"` (or the chosen section title).
2. `alpha/SKILL.md` line 355 — one ref. Updates to `delta/SKILL.md`.
3. `beta/SKILL.md` line 175 — one ref. Updates to `delta/SKILL.md`.

References that STAY at `operator/SKILL.md` (per scope rule — these point at non-δ-role content):

- `post-release/SKILL.md` line 39 → `operator/SKILL.md §3.4` for **release mechanics** (Phase 4c will move). Add `delta/SKILL.md` for the role-policy companion ref; but path-only update is permissible.
- `release/SKILL.md` lines 35, 238 → same.
- `activation/SKILL.md` multiple refs → identity, dispatch configs, §3.4 mechanics (all coordinator/harness, stay).
- `review/SKILL.md` line 141 → `operator/SKILL.md` §5.2 (dispatch config; stays).
- `CDD.md` multiple refs → identity, §3.3 (operator-coordinator), §3.4 (release mechanics), timeout-recovery (harness). All stay.
- `COHERENCE-CELL.md` doctrinal refs predicting the split — leave as historical anchor; doctrine surface owns its own narrative.
- `RECEIPT-VALIDATION.md` and `COHERENCE-CELL-NORMAL-FORM.md` Phase 4 target references — these mention the target file `operator/SKILL.md` as the Phase 4 surface; leave as historical record (these documents were authored before this cycle and reflect the state then). Future cycles can update if needed; this cycle does not edit doctrine surfaces.

**Question: should `post-release/SKILL.md` line 39 and `release/SKILL.md` lines 35, 238 add a `delta/SKILL.md` co-reference for the role-policy half (since the policy "δ owns tag/release/deploy" is δ-role)?** 

Decision: YES for `post-release/SKILL.md` line 39 ("**δ owns tag/release/deploy** … per `CDD.md` §1.4 β algorithm and `operator/SKILL.md` §3.4") — this is a role-policy citation. Update to point at delta/ AND operator/ (the role-policy lives in delta/; the mechanics live in operator/). Similarly for `release/SKILL.md` lines 35, 238. This is a path swap, not a content change — discharges AC5 cleanly.

## §4 No-content-loss audit

Every line in the moved sections must end up in `delta/SKILL.md` with semantics preserved. The β-rigor check (cnos#393 Rule 7) demands no δ-content drops. The verification at β-collapsed time: a diff-by-diff comparison of operator/SKILL.md `before` vs `after` shows every removed paragraph appears in `delta/SKILL.md`.

The only intentional content rewrites:
- §3a phrase "Phase 4 (δ split) — relocation target" — REWRITE to "This section has landed here as Phase 4a of cnos#366 (the relocation target named in the precursor cnos#393)." Substance is preserved; the futuristic framing becomes the present tense.
- §3.4 "this is δ's authoritative release-cut algorithm" framing — split into role-policy paragraph in `delta/SKILL.md` + mechanics paragraphs in `operator/SKILL.md` with explicit cross-references. No content lost; layering only.

## §5 Implementation-contract conformance (this section is binding per cnos#393 Rule 7)

| Axis | Pinned value (from cnos#397 issue body) | Conformance plan |
|---|---|---|
| Language | Markdown | All edits are .md only. ✓ |
| CLI integration target | N/A | No CLI added or edited. ✓ |
| Package scoping | `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` (new) + edits to `cdd/operator/SKILL.md` | New file at exact path; operator/SKILL.md at exact path. ✓ |
| Existing-binary disposition | operator/SKILL.md keeps operator-as-coordinator + dispatch role; boundary-policy + δ-inward-membrane sections move to delta/SKILL.md. Harness mechanics → Phase 4b. Release effector → Phase 4c. | Plan implements this exactly: §1+§2+§5+§6+§7+§8+§9+§10 stay; §3 splits (role-policy moves, mechanics stay for Phase 4c); §3a moves entirely; §4 moves entirely. ✓ |
| Runtime dependencies | None | No deps added. ✓ |
| JSON/wire contract | N/A | No JSON edited. ✓ |
| Backward compat | All existing operator/SKILL.md refs from other skills updated to delta/ where appropriate; operator/SKILL.md continues to exist. | Plan updates exactly the 4 refs that cite δ-role content (gamma×2, alpha×1, beta×1) + optional co-ref additions in post-release/release. operator/SKILL.md preserved. ✓ |
