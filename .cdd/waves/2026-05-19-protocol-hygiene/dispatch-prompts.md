# Dispatch prompts — Wave 2026-05-19 Protocol Hygiene

These are the γ prompts δ dispatches for each cycle. δ uses `cn dispatch --role γ --branch cycle/{N}` (or `claude -p` equivalent per `operator/SKILL.md`).

Each γ then produces α and β prompts after creating the cycle branch and authoring `.cdd/unreleased/{N}/gamma-scaffold.md` (per rule 3.11b).

Per `gamma/SKILL.md` §2.5, the prompt format is intentionally short — point at the issue, name the skills to load, do not paraphrase the issue body.

---

## Cycle 375 — γ-side pre-dispatch gate

```text
You are γ. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md and follow its load order.
Issue: gh issue view 375 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/375 (you create from origin/main per gamma/SKILL.md §2.5 Step 3a)
Wave: .cdd/waves/2026-05-19-protocol-hygiene/manifest.md (this wave's standing permissions and timeout budgets apply)
Tier 3 skills (issue-specific):
- cdd/gamma (the patch's primary surface)
- cdd/CDD (alternative patch surface)
- cdd/review (the dual β-side surface to reference)
```

α/β skill loading (γ adds when dispatching):
- α loads: cnos.core/skills/write/SKILL.md, cdd/alpha/SKILL.md, cdd/gamma/SKILL.md (read-only for context on the patch surface), cdd/review/SKILL.md (read-only — rule 3.11b is the β-side dual)
- β loads: cdd/beta/SKILL.md, cdd/review/SKILL.md (review surface — rule 3.11b lives there)

---

## Cycle 378 — rule 3.11b discoverability under wave-mode

```text
You are γ. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md and follow its load order.
Issue: gh issue view 378 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/378 (you create from origin/main per gamma/SKILL.md §2.5 Step 3a)
Wave: .cdd/waves/2026-05-19-protocol-hygiene/manifest.md (this wave's standing permissions and timeout budgets apply)
Tier 3 skills (issue-specific):
- cdd/review (rule 3.11b primary patch surface)
- cdd/alpha (§2.6 row patch surface)
- cdd/operator (§5.2 + §10 wave-manifest cross-reference)
- cdd/CDD (γ=δ collapse rule context)
Empirical anchor (cite in patched skill text):
- usurobor/cph cdr-refactor wave 2026-05-18 — master cph#11; subs cph#12, #13, #14, #15
- Four sub-uniform configuration; three distinct β substantive reads of the same wave-manifest-as-γ-artifact pattern
```

α/β skill loading:
- α loads: cnos.core/skills/write/SKILL.md, cdd/alpha/SKILL.md, cdd/review/SKILL.md (read-only — primary patch surface), cdd/operator/SKILL.md (read-only — §5.2 wave-mode), cdd/CDD.md (read-only — γ=δ collapse)
- β loads: cdd/beta/SKILL.md, cdd/review/SKILL.md (review's own rule 3.11b — meta-review)

---

## Cycle 377 — cross-repo coordination protocol

```text
You are γ. Project: cnos.
Load src/packages/cnos.cdd/skills/cdd/gamma/SKILL.md and follow its load order.
Issue: gh issue view 377 --repo usurobor/cnos --json title,body,state,comments
Branch: cycle/377 (you create from origin/main per gamma/SKILL.md §2.5 Step 3a)
Wave: .cdd/waves/2026-05-19-protocol-hygiene/manifest.md (this wave's standing permissions and timeout budgets apply)
Mode: design-and-build (per issue body; design phase converges state machine + directional cases; build phase writes the skill and updates referencing skills)
Tier 3 skills (issue-specific):
- cnos.core/skills/skill (skill authoring meta-skill — needed for new SKILL.md frontmatter and shape)
- cdd/design (design phase doctrine)
- cdd/gamma (read-only — extracting §2.1 + §2.7 fragments from here)
- cdd/post-release (read-only — extracting §5.6b from here)
- cdd/issue/labels (already loaded indirectly via gamma; no re-load needed)
Empirical anchors (the new protocol must validate retroactively against these without contradiction):
- .cdd/iterations/cross-repo/tsc/cdd-supercycle/LINEAGE.md (outbound iteration trace; 6 patches)
- .cdd/iterations/cross-repo/gait-support-paths/bootstrap-cdr/{LINEAGE.md,FEEDBACK.patch,STATUS} (inbound proposal; master/sub)
- .cdd/iterations/cross-repo/cn-sigma/agent-activate-skill/{LINEAGE.md,FEEDBACK.patch} (inbound proposal; 1:1)
- .cdd/iterations/cross-repo/cph/coherence-drift-sweep-followup-2026-05-18/{LINEAGE.md,cdr-changelog-rule.patch,cdd-iteration.md} (cross-repo ε iteration with patches landing in both repos)
Additional protocol observations to fold (filed in cnos#377 comments + recent bundle LINEAGEs):
- "drafted" event used by cn-sigma in STATUS (not in documented 5-event vocabulary)
- bundle file set varies: bootstrap-cdr has 4 files (ISSUE/STATUS/README/LINEAGE); agent-activate-skill has 3 (no README)
- source ## Source Proposal block with placeholders for target γ to fill is a cleaner author-side convention than insert-fresh-block
- hat-collapse attribution (γ@A + ε@B same session) has no protocol-defined form
- inbound mirror at cnos:.cdd/iterations/cross-repo/{source}/{slug}/ is dual-purpose (iteration content + patch trace) with no codified shape
```

α/β skill loading:
- α loads: cnos.core/skills/write/SKILL.md, cnos.core/skills/skill/SKILL.md, cdd/alpha/SKILL.md, cdd/design/SKILL.md, cdd/gamma/SKILL.md (read-only — extraction source), cdd/post-release/SKILL.md (read-only — extraction source), cdd/issue/SKILL.md (read-only — Source Proposal block format)
- β loads: cdd/beta/SKILL.md, cdd/review/SKILL.md, cdd/gamma/SKILL.md (the canonical surface being reshaped), cdd/post-release/SKILL.md (same)

---

## Wave-level dispatch sequence (δ executes)

Per the wave manifest, the three cycles are parallel. δ may dispatch all three γ sessions concurrently; each γ then runs its α and β sequentially.

If δ's environment cannot run three concurrent γ sessions, dispatch in the order suggested in `status.md`: 375 (smallest, γ-side-only) → 378 (medium, multi-skill) → 377 (largest, new skill file). The order is suggestive; β handles merge conflicts at merge time.

After all three cycles merge, ε writes `.cdd/iterations/wave-2026-05-19-protocol-hygiene.md` consolidating any cross-cycle findings.
