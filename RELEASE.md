# RELEASE.md

## Outcome

Patch release: δ post-cycle disconnect for 3.59.0.

No Go kernel changes. No binary changes. Spec, skill, and doc patches only — all from the #230 cycle's close-out triage and the first live δ (operator) session.

## Why it matters

The 3.59.0 cycle produced 10 commits of post-release work (PRA, close-outs, skill patches, δ operator skill, CTB v0.1). Without a tag, the triad's final output remains entangled with whatever comes next. This release is the disconnection point — the tagged snapshot of everything the #230 cycle and its δ session produced.

## Added

- **δ operator skill** (`operator/SKILL.md`): New CDD role — δ owns external gates, session routing, override authority. First addition to the CDD role model since the triad. Per `COHERENCE-FOR-AGENTS.md`: δ is whole-to-whole composition — a new one-as-two boundary between the triad-as-whole and the platform, not a fourth triad role.
- **CTB Language Spec v0.1** (`docs/alpha/ctb/LANGUAGE-SPEC.md`): Normative reference for skill modules — signatures, scope, dispatch, composition, effect-plan boundary. First concrete `calls_dynamic` migration on `alpha/SKILL.md`.
- **CTB Semantics Notes** (`docs/alpha/ctb/SEMANTICS-NOTES.md`): Non-normative conceptual rationale behind the spec.
- **Reflect §3.6** — decision-basis capture: when a daily records a triage/classification/disposition, record criteria + per-item basis. A conclusion without basis is a claim the next session cannot verify.
- **Daily template** `## Decisions` section with basis prompt.

## Changed

- **CDD.md §Tracking:** reachability preflight added as 3rd mandatory polling part. Branch glob broadened. Synchronous baseline step before transition loop. Baseline rule stated explicitly.
- **CDD.md γ algorithm:** step 1 = git identity (`gamma@cdd.{project}`). All subsequent steps renumbered.
- **CDD.md β step 8:** defer to δ (was "γ/operator"), reference δ §3.5 signal.
- **β SKILL.md Rule 1:** refusal is not terminal — polling continues regardless.
- **γ SKILL.md:** δ added to inputs, calls, and load order (step 5). Waits for δ completion signal before close-out triage. Step map and refs updated for renumbering.
- **δ §3.2/§2.1:** execute on request, not on observation. Heartbeat observation ≠ gate request.
- **δ §3.4:** post-cycle release as mandatory disconnection (not optional assessment).
- **δ §3.5:** signal γ after release-phase gates.
- **post-release/SKILL.md:** step refs updated (12a → 13a).
- **alpha/SKILL.md:** `calls_dynamic` frontmatter (`issue.tier2_bundles`, `issue.tier3_skills`) replaces prose "Tier 2 and Tier 3 skills named by the issue."
- **CTB README.md:** updated document map with authority rules across Vision/Spec/Notes/kernel.

## Validation

- No Go code changes — binary is unchanged from 3.59.0.
- All changes are markdown spec/skill/doc files in `src/packages/` and `docs/`.
