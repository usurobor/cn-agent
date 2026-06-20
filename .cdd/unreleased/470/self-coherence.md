# α self-coherence — cycle/470

<!--
section-manifest:
  planned: [Gap, Skills, AC1-impl, AC2-impl, AC3-impl, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
  completed: [Gap, Skills]
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

## Skills

**Tier 1 (lifecycle):**

- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (α role; §2.5 incremental self-coherence; §2.6 pre-review gate rows 1–15; §3.6 implementation-contract is δ's not α's)
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` (to verify ACs are independently testable; AC oracle/surface re-read)
- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` (this is design-and-build mode; loaded before authoring the manifest schema and prompt template; impact graph traced — see CDD Trace step 4)

Note: `src/packages/cnos.cdd/skills/cdd/CDD.md` is referenced by `alpha/SKILL.md` §"Load Order" Tier 1; absent from this checkout (only `alpha/`, `beta/`, `gamma/`, `delta/`, `issue/`, `design/`, `review/`, `release/`, ... role/lifecycle subskills exist). α treats the canonical lifecycle contract as expressed in `alpha/SKILL.md` itself + `cnos.cds/skills/cds/CDS.md` (CDS canonical artifact order / lifecycle step table) per the redirect in `alpha/SKILL.md` §"Load Order".

**Tier 3 (issue-specific):**

- `src/packages/cnos.core/skills/agent/activate/SKILL.md` (the prompt template α authors directs the wake to invoke this; §2.1 six-item load order; identity-confirmation gate)
- `src/packages/cnos.core/skills/agent/attach/SKILL.md` (same; §2.4 follow-up sync procedure; §2.5 ephemeral attach for ephemeral mode; AGENT-ACTIVATION-LOG-v0 cursor mechanics)
- `src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` (cnos#468 merged predecessor; the prompt cites it for label-application discipline; §4.1 Sigma may; §4.2 Sigma may not; the manifest's allowed/disallowed surfaces conform to the boundary at §4)
- `src/packages/cnos.core/skills/skill/SKILL.md` (skill-format meta-skill; AC1's contract skill conforms; §1.0 artifact_class; §3.1 frontmatter requirements)
- `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` (terminal-report shape AC4 names; §5 entry format; `class:` taxonomy `heartbeat / substantive / inaugural / directive-out`; per AC4 prose the cycle-complete class will be wired in Sub 6 but α enumerates all five at the contract level)

**Package-shape reference (read for shape, not modified):**

- `src/packages/cnos.core/cn.package.json` (existing package shape; AC2 manifest does NOT register a new command in the `commands` map per pinned axis)
- `src/packages/cnos.core/orchestrators/daily-review/orchestrator.json` (existing sibling orchestrator shape; confirms `orchestrators/{name}/` is not constrained to a single file; new sibling `wake-provider.json` is valid)
- `src/packages/cnos.core/commands/{daily,weekly,save}/` (existing command shape; α did NOT use this form — γ pinned `orchestrators/agent-admin/`)
- `.github/workflows/claude-wake.yml` (existing wake the provider supersedes at Sub 3 cutover; READ ONLY — AC7 invariant)

**γ scaffold (read for context):**

- `.cdd/unreleased/470/gamma-scaffold.md` (full form-choice rationale, expected touched surfaces, AC mapping table with mechanical gate block, implementation contract pin)

**Skills explicitly NOT loaded (and why):**

- `cnos.cdd/skills/cdd/beta/SKILL.md` — β role; α does not load β skills (alpha/SKILL.md §2.1 "do not load β or γ role skills")
- `cnos.cdd/skills/cdd/gamma/SKILL.md` — γ role; same
- `cnos.cdd/skills/cdd/delta/SKILL.md` — δ role; same; also explicitly out of scope for Sub 2 (Sub 5 territory)
- `cnos.cdd/skills/cdd/plan/SKILL.md` — implementation sequencing is non-trivial only across many files; this cycle creates 3 files with clear dependency order (AC1 contract skill → AC2 manifest → AC3 prompt template); plan declared "not required" per α/SKILL.md §2.2 with this justification
- Tier 2 `eng/*` — no eng-language-specific work; manifest is JSON + skill/prompt are markdown; no Go/Python/shell

