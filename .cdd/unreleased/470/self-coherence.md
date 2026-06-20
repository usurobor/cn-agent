# α self-coherence — cycle/470

<!--
section-manifest:
  planned: [Gap, Skills, AC1-impl, AC2-impl, AC3-impl, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
  completed: [Gap]
-->

## Gap

**Issue:** [cnos#470](https://github.com/usurobor/cnos/issues/470) — `agent-admin/wake-provider: cnos.core agent-admin wake (admin-only, no cell execution)`. Sub 2 of cnos#467 (master tracker `agent/wake-orchestration`); builds on Sub 1 (cnos#468 — label doctrine, merged at `c0048bef`).

**Version / mode:** **design-and-build** (per cnos#470 mode declaration). Provider declaration is package content; prompt template is content; contract skill is content. No code, no `cn` binary edits, no `.github/workflows/` edits. A design call (manifest schema fields, allowed/disallowed surfaces structure, defer-path shape) lands within the cycle. The γ scaffold (`.cdd/unreleased/470/gamma-scaffold.md`) pins the implementation contract; α executes within its axes.

**Form-choice acknowledgment:** γ pinned the provider declaration entry at `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` (manifest, schema `cn.wake-provider.v1`) + `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` (prompt template), NOT at `cnos.core/commands/install-wake/`. α accepts the γ pin. No structural override is required: inspection of `cnos.core/orchestrators/daily-review/orchestrator.json` confirms that filename + schema are independent (the existing file is named `orchestrator.json` with `kind: cn.orchestrator.v1`, but `cn` does not constrain sibling files in the same directory or constrain `orchestrators/{name}/` to a single schema). The new `wake-provider.json` carries a new schema `cn.wake-provider.v1` and lives as a sibling shape under `orchestrators/agent-admin/`; the Sub 3 renderer will dispatch by schema. α's authoring proceeds on the γ-pinned path.

**Form choice — secondary file (AC1 contract skill):** AC1 names the surface `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` as the canonical (or design-call alternative) location. α uses the canonical path. This skill is a sibling of `cnos.core/skills/agent/{activate, attach, label-doctrine}` and conforms to `cnos.core/skills/skill/SKILL.md` (the skill-format meta-skill).

**Bootstrap exception context:** This cycle executes through the pre-dispatch δ/channel path (γ-interface session acts as bootstrap-δ; γ/α/β spawned as separate sub-agents via the Agent tool; `.cdd/unreleased/470/` artifacts are the shared memory). β is NOT yet spawned; α does not poll for β verdicts during the authoring run. δ runs β after α signals review-readiness. α exits after signaling review-readiness (sequential bounded dispatch per `cnos.cds/skills/cds/CDS.md` §"Field 6: Actor collapse rule"); δ re-dispatches α after β merges for `alpha-closeout.md`.

**Implementation-contract axes (γ-pinned; α MUST NOT improvise):**

| Axis | Pinned value | α compliance plan |
|---|---|---|
| Language | Markdown (skill + prompt template) + JSON (declaration manifest, schema `cn.wake-provider.v1`) | Only `.md` and `.json` files created; no `.go`, no `.sh`, no `.yml` |
| CLI integration target | None for this sub (`cn wake install` is Sub 3) | No edits to `src/go/`; no edits to `cnos.core/cn.package.json`'s `commands` map |
| Package scoping | `src/packages/cnos.core/` only | Only `cnos.core/skills/agent/wake-provider/` + `cnos.core/orchestrators/agent-admin/` touched |
| Existing-binary disposition | N/A — no binaries | n/a |
| Runtime dependencies | None — static data + markdown | manifest is parseable JSON; prompt is markdown |
| JSON/wire contract | Manifest declares `"schema": "cn.wake-provider.v1"` as the first field | wake-provider.json first key = `schema` |
| Backward compat | `.github/workflows/claude-wake.yml` MUST remain byte-identical on `cycle/470` | verified by `git diff origin/main..HEAD -- .github/workflows/claude-wake.yml \| wc -l` MUST return 0 |
