# ε cdd-iteration — cycle/425

**Issue:** [cnos#425](https://github.com/usurobor/cnos/issues/425) — Capture remote-runner-delegation primitive + first use (v3.82.0 tag retarget via GH Actions).
**Mode:** docs + skill patch + workflow + receipt; γ+α+β-collapsed-on-δ. **Substantive** entry (not a courtesy stub): records one `cdd-protocol-gap` finding with disposition `patch-landed-in-cycle`.

## Findings

### F1 — cdd-protocol-gap: boundary model incomplete (remote-runner delegation needed naming)

**Class:** `cdd-protocol-gap`
**Severity:** binding (boundary discipline)
**Disposition:** `patch-landed-in-cycle` — this cycle IS the patch.

**The gap.** Prior to this cycle, the cnos boundary model was drawn at the local-shell surface. δ-as-role's outward membrane (per `cnos.cdd/skills/cdd/delta/SKILL.md §1`) governed receipt-level acceptance; the release-effector skill (`release-effector/SKILL.md`) governed locally-invoked release moves; the operator skill (`operator/SKILL.md`) governed gate-action confirmation. **But there was no doctrine for artifacts that, when committed, cause another body to execute** — workflow files, deploy hooks, extension manifests, scheduled jobs, webhook handlers, MCP server registrations. The class is "artifacts that cause another body to execute." The boundary model treated these as configuration text rather than as effect surfaces.

**How the gap surfaced.** During this cycle's authoring session, an agent (under Claude's local harness with shortened hands) discovered that it could not push the `3.82.0` tag locally (harness sandbox restriction) but could write a `.github/workflows/repoint-3.82.0.yml` file that would cause a GitHub Actions runner to push the tag on the agent's behalf. The discovery felt natural rather than exceptional — the workflow file is just text from the harness's perspective, but it is delayed execution from the runner's perspective. The operator then directed the agent to formalize the move into doctrine ("execute as delta or gamma" + "you can create gh action and then run it"). The discovery exposed that the existing boundary model was guarding one limb (the local shell) while leaving other limbs (the CI runner, deploy hooks, third-party services) unnamed.

**Why this is a `cdd-protocol-gap` rather than a `cdd-tooling-gap` or `cdd-skill-gap`.** Tooling-gap means the harness/validator lacked a mechanism; skill-gap means an existing skill's content was incomplete. Protocol-gap means the *protocol shape itself* lacked a naming for a class of moves. This is the protocol case: δ-as-role's authority extended to receipt acceptance and locally-invoked release effectors, but the *class of artifacts that cause another body to execute* was not named as a δ-class effect surface in any skill. Adding this naming is a protocol extension, not a tooling repair or a skill polish.

**The patch (landed in this cycle).** The fix has 4 pieces, all landed together in commit `334f1ca6`:

1. **Essay (`docs/gamma/essays/BOX-AND-THE-RUNNER.md`).** Names the primitive ("remote-runner delegation"), states the rule ("Any artifact that can make another system execute is an effect surface"), pins the 6-field receipt convention, and frames the trick-vs-protocol distinction so future cycles read this as an extension of the boundary model rather than a clever workaround.

2. **Skill section (`delta/SKILL.md §8 "Remote-runner delegation — δ-class effect surface").** Operationalizes the doctrine inside δ-as-role: 6 sub-sections covering the required 6-field receipt, the authoring-order rule, the one-shot self-delete discipline, composition with the existing V/δ wall, and the first-use anchor.

3. **First-use artifact (`.github/workflows/repoint-3.82.0.yml`).** A one-shot, self-deleting workflow that force-moves the `3.82.0` tag to `fd1d654e` so `release.yml` publishes the correct body. The artifact exists *as the first instance* of the discipline the essay + skill section name; future cycles authoring similar artifacts inherit the same pattern.

4. **First receipt (`.cdd/unreleased/425/remote-runner-receipt-3.82.0.md`).** A 6-field receipt that instantiates the convention for the first-use artifact. Operator-acceptance criterion explicitly named; failure-modes table with 6 rows + mitigations; acceptance-criteria with 6 numbered criteria + partial/rejection branches.

**Why `patch-landed-in-cycle` is the correct disposition.** The dispatch brief explicitly framed this cycle as "single cycle bundling: capture the remote-runner-delegation primitive as cnos doctrine + execute its first use." Both the capture and the first use are deliverables of this cycle, not of a downstream MCA. The patch is not pending; it is on the branch. ε does not need to propose a follow-on issue or schedule an MCA — the gap is closed by the cycle's matter. `no-patch` would be incorrect (the patch exists); `MCA` would be incorrect (no follow-on action required); only `patch-landed-in-cycle` honestly describes the disposition.

**Forward signal (not a finding; informational).** The doctrine-before-first-use precedence rule (the bundle ships atomically; the workflow's push trigger fires only post-merge by which point the doctrine is on main) is a *structural pattern* this cycle establishes. It is reusable for any future cycle that bundles "name a primitive" + "first exercise of the primitive." ε does not file this as a separate finding because the pattern is documented inside `delta/SKILL.md §8.2`; future cycles can find it there. But ε notes here that the pattern is a candidate for explicit naming in `gamma/SKILL.md` or `dispatch/SKILL.md` as a class of cycle (a "doctrine-and-first-use bundle") if the pattern recurs.

## Protocol-gap signals (across receipt-stream)

This is the single `cdd-protocol-gap` finding for this cycle. Looking across the receipt stream (per `cnos.handoff/skills/handoff/receipt-stream/SKILL.md`), this finding stands alone — no prior `cdd-*-gap` finding in `.cdd/iterations/INDEX.md` names a boundary-model gap; the gap was latent until the harness-vs-runner asymmetry surfaced it.

ε does not currently see a pattern across multiple cycles that this finding reinforces. It is a first-of-its-kind protocol-gap finding. Future ε runs will track whether subsequent cycles authoring remote-runner artifacts (a) cite this doctrine, (b) produce well-formed 6-field receipts, (c) honor the doctrine-before-first-use precedence rule, and (d) one-shot self-delete persistent-vs-ongoing where applicable. If those four properties hold across the next 3–5 remote-runner-using cycles, the doctrine is field-proven; if they drift, ε will file a follow-on `cdd-protocol-gap` finding to strengthen the doctrine (e.g., add a CUE schema for the receipt; add a `cn cdd verify`-side check for workflow-without-receipt).

## Non-findings (worth recording)

- **Verbatim seed preservation held.** The operator's seed text from the dispatch brief is preserved verbatim across the six core prose sections of the essay (Point, Progress note, What I realized, What this means for cnos, What not to celebrate, Memory to keep). Essay-class scaffolding (paragraph transitions, the effect-surfaces code block, cross-references to existing δ/release-effector doctrine) is additive, not substance-displacing. The discipline of preserving the seed where possible is a positive ε observation for future seed-driven essay cycles.

- **Citation closed loop across the 5 deliverables.** Essay → skill section → workflow header → receipt → back to essay. The cross-citations form a closed graph; no artifact stands alone. This is the load-bearing legibility property for the doctrine bundle and is the kind of structural discipline ε aggregation can detect across the receipt stream.

- **β-independence collapse named, not papered over.** This cycle ran as `γ+α+β-collapsed-on-δ`. Per `COHERENCE-CELL.md §β Independence as Contagion Firebreak`, the collapse compromises structural-independence; the receipt is closed-as-degraded at that axis. `beta-closeout.md` names the collapse explicitly with the cycle-414 / cycle-424 precedent inheritance and the dispatch authorization. The discipline of naming the collapse rather than hiding it is good ε hygiene and inherited from prior essay cycles.

- **Doctrine essay tense is descriptive, not prescriptive.** The essay says "cnos must name the move before the move becomes folklore" and "the agent boundary is not the place where the model sits; it is the full path from intention to effect." Both are descriptions of what cnos's coherence model already implies once extended, not prescriptions for new protocol invention. The skill section operationalizes those descriptions inside δ-as-role's existing authority. No new role station, no new schema, no new validator predicate is created — only the existing δ membrane extends inward to govern the authoring of remote-runner-triggering artifacts.

## Verdict

`protocol_gap_count: 1`.

One `cdd-protocol-gap` finding (F1: boundary model incomplete — remote-runner delegation needed naming) with disposition `patch-landed-in-cycle`. The patch is on the branch as commit `334f1ca6` (essay + skill section + workflow + receipt). No follow-on MCA required; no `next-MCA` disposition; no `no-patch` disposition.

INDEX.md row: `findings=1, patches=1, MCAs=0, no-patch=0`.

Filed by ε@cnos (γ+α+β-collapsed-on-δ) on 2026-05-23.
