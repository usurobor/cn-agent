# α self-coherence — cycle/470

<!--
section-manifest:
  planned: [Gap, Skills, AC1-impl, AC2-impl, AC3-impl, ACs, Self-check, Debt, CDD-Trace, Review-readiness]
  completed: [Gap, Skills, AC1-impl, AC2-impl, AC3-impl, ACs]
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

## ACs

Each AC is verified against the actual diff (β SKILL Rule 6: code is evidence). Oracles match those in the γ scaffold's AC mapping table, with regex evidence re-grepped at implementation SHA `2c7b1437` (the last implementation commit before this self-coherence section).

### AC1 — Wake-provider declaration contract skill exists at cnos.core

**Invariant (verbatim from #470):** A skill at `cnos.core/skills/agent/wake-provider/SKILL.md` (or design-call alternative) defines the wake-provider declaration contract: required fields, optional fields, the substrate-rendering target, and the split between "what the package declares" vs "what the renderer materializes." Concrete enough that cnos.cdd's dispatch wake provider declaration (Sub 4) can be authored against this skill alone.

**Status:** **pass**

**Evidence:**

- File exists: `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` (created at commit `61588ca0`).
- Frontmatter conforms to `cnos.core/skills/skill/SKILL.md` §3.1 (required fields: `name`, `description`, `artifact_class`, `kata_surface`, `governing_question`; this skill declares all five + `visibility`, `parent`, `triggers`, `scope`, `inputs`, `outputs`, `requires`).
- Section headers (33 `#`-prefixed lines per `grep -cE "^#+ " SKILL.md`): includes `## Core Principle`, `## Algorithm`, `## 1. Define`, `## 2. Unfold`, `## 3. Rules`, `## 4. Verify`, `## 5. Cross-references`, `## 6. Failure modes`, `## 7. Authority and stability`.
- Required-fields enumeration: §2.1 table (12 required fields named: `schema`, `name`, `package`, `role`, `responsibilities`, `admin_only`, `input_contract`, `output_contract`, `allowed_surfaces`, `disallowed_surfaces`, `defer_path`, `prompt_template`, `cross_references`).
- Optional-fields enumeration: §2.2 table (6 optional fields named: `description`, `agent_variable`, `permission_intent`, `concurrency_intent`, `superseded_substrate_artifact`, `relationship_to_substrate`).
- Substrate-rendering target: §2.4 ("Today's substrate is GitHub Actions with `anthropics/claude-code-action@v1` as the agent-execution carrier"), naming substrate artifact path, trigger encoding, permission encoding, run-identity binding, secret binding, prompt-receiving field.
- Package-authority-vs-renderer-authority split: §1.2 introduction + §2.5 canonical grep-able table making the split mechanical.
- Cnos#467 / wake-orchestration / label-doctrine / cnos#468 cites: 14 occurrences per `grep -cE "cnos#467|wake-orchestration|label-doctrine|cnos#468" SKILL.md` (≥ 2 per scaffold oracle).
- "Concrete enough for Sub 4 to author cnos.cdd dispatch wake provider against": §2.6 "Authoring a new wake provider against this contract" gives the 6-step procedure; §1.2 table contrasts admin vs dispatch role explicitly; §5 names cnos#467 Sub 4 as a downstream consumer expected to pattern-copy.

### AC2 — Agent-admin wake provider declaration entry exists in cnos.core

**Invariant (verbatim from #470):** A concrete entry at `cnos.core/commands/install-wake/` or `cnos.core/orchestrators/agent-admin/` (design call within cycle; document) carries the agent-admin wake provider declaration — identity, responsibilities (admin-only), input/output contract, allowed/disallowed surfaces, defer-path, prompt template. The declaration is substrate-agnostic.

**Status:** **pass**

**Form choice:** γ-pinned form `cnos.core/orchestrators/agent-admin/` (per §Gap form-choice acknowledgment).

**Evidence:**

- File exists: `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` (created at commit `4c8f30c8`).
- JSON parses: `jq . wake-provider.json > /dev/null` exit 0.
- Schema declaration: `jq -r 'keys_unsorted[0]'` returns `schema`; `jq -e '.schema == "cn.wake-provider.v1"'` returns true. The `schema` field is the first key of the object (γ-scaffold pinned axis).
- Conforms to AC1's contract: all 12 required fields of §2.1 are present (`jq -e '. | has("schema") and has("name") and has("package") and has("role") and has("responsibilities") and has("admin_only") and has("input_contract") and has("output_contract") and has("allowed_surfaces") and has("disallowed_surfaces") and has("defer_path") and has("prompt_template") and has("cross_references")'`).
- Identity: `name: "agent-admin"`, `package: "cnos.core"`, `role: "admin"`, `admin_only: true` (`jq -e '.admin_only == true'` returns true).
- Responsibilities (admin-only enumeration): 8 array entries covering activate, attach, channel sync, status reporting, issue creation/refinement, label application, directive routing, off-role refusal.
- Input contract: `input_contract.triggers = ["schedule", "issues_opened_title_match"]` + `input_contract.inbound` describing home thread + open issues.
- Output contract: `output_contract.channel_log_convention` cites `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`; `class_taxonomy` enumerates 5 values (4 named in scaffold + `cycle-complete` reserved for Sub 6).
- Allowed surfaces: 6 entries enumerating channel-log writes, state, comments, label application, issue creation.
- Disallowed surfaces: 9 entries including `cell_execution` and `.github/workflows/` (both required per AC5 oracle).
- Defer-path: 3 keys (`cell_shaped_directive`, `off_role_directive`, `ambiguous_directive`) with explicit prose.
- Prompt template: `prompt_template: "prompt.md"` (the sibling file at AC3).
- Substrate-agnostic: 5 manifest-side hits on the `github|workflow|yaml|GITHUB_TOKEN|claude-code-action|runs-on` regex; ALL fall within the γ-scaffold carve-out — L53 (AC5 mandatory enumeration of `.github/workflows/` as disallowed surface), L79 (descriptive renderer-mapping prose in `permission_intent_notes` naming "GitHub Actions 'permissions:' YAML block" as the substrate the renderer projects to), L83 (same for `concurrency_intent.notes`), L85 (`superseded_substrate_artifact` field — the AC7 carve-out the contract explicitly provides), L86 (`relationship_to_substrate` — the AC7 prose section). No substrate emission; only descriptive references to the substrate the renderer (Sub 3) will project to.

### AC3 — Prompt template enforces admin-only constraint

**Invariant (verbatim from #470):** The provider's prompt template contains explicit language: enumerates the admin responsibilities (activate, attach, channel sync, status reporting, issue creation/refinement, label application); explicitly forbids cell execution; specifies the defer-path for cell-shaped directives (defer to relevant `protocol:{P}` dispatch wake if installed; surface to operator if not); cites cnos#468 label-doctrine for any label application Sigma may perform.

**Status:** **pass**

**Evidence:**

- File exists: `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` (created at commit `2c7b1437`).
- Activate/attach enumeration: 14 occurrences per `grep -cE "activate|attach" prompt.md` (≥ 2 per scaffold). Specifically: §"Identity and activation" calls out both skills by canonical path; §"Admin responsibilities" steps 1 + 2 name them as the first two responsibilities.
- "Do not execute cells" / "MUST NOT execute": 3 occurrences per `grep -ciE "do not execute cell|never execute cell|no cell execution|MUST NOT execute" prompt.md` (≥ 1 per scaffold). Specifically: §"Admin-only boundary: MUST NOT execute cells" carries the explicit "You MUST NOT execute cells under any circumstance" + "MUST NOT execute the cell inline" + the defer-path naming.
- Defer / dispatch wake: 13 occurrences per `grep -ciE "defer|dispatch wake" prompt.md` (≥ 2 per scaffold). Specifically: §"Defer-path for cell-shaped directives" provides the 4-step deferral procedure; §"Defer-path for off-role and ambiguous directives" provides the parallel off-role and ambiguous procedures.
- Label-doctrine / cnos#468 cite: 6 occurrences per `grep -cE "label-doctrine|cnos#468" prompt.md` (≥ 1 per scaffold). Specifically: §"Admin responsibilities" step 6 cites both `[label doctrine](../../skills/agent/label-doctrine/SKILL.md)` and `(cnos#468)` and grounds the boundary in §4.1 / §4.2 of the doctrine skill.
- Enumerated responsibilities: 19 occurrences per `grep -ciE "status report|channel|issue creat|label appl" prompt.md` (≥ 3 per scaffold). §"Admin responsibilities" enumerates 8 steps including each named responsibility.

### AC4 — Input + output contracts documented

**Invariant (verbatim from #470):** Declaration names triggers (schedule cron slots; optional `issues: opened` for `claude-wake`-titled trigger issues) as inputs, and the terminal-report shape (channel log entry per `AGENT-ACTIVATION-LOG-v0.md`; frontmatter fields; cursor advancement) as outputs. The contract is precise enough that Sub 3 (renderer) can produce a working workflow + Sub 6 (cycle-complete reading) can extend the `class:` taxonomy.

**Status:** **pass**

**Evidence:**

- Input contract section: `input_contract` object in manifest with `triggers: ["schedule", "issues_opened_title_match"]` + per-trigger descriptions in `trigger_descriptions` + `inbound` prose describing home thread + open issues. The taxonomy is logical (renderer maps to substrate per AC1 §2.5 split table).
- Output contract section: `output_contract` object with 4 keys — `channel_log_convention: "docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md"` (cites canonical source), `writer_surface: ".cn-{agent}/logs/YYYYMMDD.md"`, `class_taxonomy: ["heartbeat", "substantive", "inaugural", "directive-out", "cycle-complete"]`, `cursor_advance: true` + `cursor_field` description.
- `AGENT-ACTIVATION-LOG-v0` / activation log cite: 11 total occurrences across manifest + prompt per `grep -cE "AGENT-ACTIVATION-LOG-v0|activation log|conventions"` (≥ 1 per scaffold).
- Four `class:` values from scaffold (heartbeat / substantive / directive-out / cycle-complete): 10 total occurrences across manifest + prompt per `grep -cE "heartbeat|substantive|directive-out|cycle-complete"` (≥ 4 per scaffold).
- Cursor advancement: named explicitly via `output_contract.cursor_advance: true` and `output_contract.cursor_field` description in manifest; named in prompt §"Admin responsibilities" step 7.
- Sub 6 forward-compatibility: `class_taxonomy` enumerates `cycle-complete` with explicit `class_taxonomy_notes` declaring it as "reserved for Sub 6 of cnos#467 ... it is enumerated here so the contract is forward-compatible and the Sub 6 wiring is a prompt-template extension, not a manifest schema change." Design call recorded per the dispatch prompt's behavioral-gap guidance.

### AC5 — Allowed + disallowed admin surfaces enumerated

**Invariant (verbatim from #470):** Declaration enumerates: allowed = `.cn-{agent}/logs/`, `.cn-{agent}/state/`, issue/PR comments, label application (per #468); disallowed = code/test/doc files outside `.cn-{agent}/` and `.cdd/`, `.github/workflows/`, branch protection, repo settings, cell execution.

**Status:** **pass**

**Evidence:**

- `allowed_surfaces` field: 6 entries (`jq -e '.allowed_surfaces | length > 0'` returns true). Each named surface from invariant present: `.cn-{agent}/logs/` (L1), `.cn-{agent}/state/` (L2), issue comments (L3), PR comments (L4), label application (L5; cites cnos#468), issue creation (L6).
- `disallowed_surfaces` field: 9 entries. Each named surface from invariant present: `cell_execution` (L1 — `jq -e '.disallowed_surfaces | any(. == "cell_execution")'` returns true; the literal string is present so the boundary is grep-able from the manifest without prose); `.github/workflows/` (L2); code files outside `.cn-{agent}/` and `.cdd/` (L3); test files (L4); doc files (L5); branch protection (L6); repo settings (L7); label definition (L8 — admin boundary per cnos#468 §4.2); other agents' surfaces (L9 — writer-locality per AGENT-ACTIVATION-LOG-v0 §0).
- `.cn-{agent}` references: 22 total occurrences across manifest + prompt per `grep -cE "\.cn-|cn-\{agent\}"` (≥ 2 per scaffold).
- Cell execution explicit in disallowed: yes, named both as `disallowed_surfaces[0]` (the literal `cell_execution` string) and prose-named in the prompt's §"Disallowed surfaces" first bullet ("Cell execution — the most important disallowed surface").

### AC6 — Cross-references present

**Invariant (verbatim from #470):** Declaration references cnos#468 (label doctrine), cnos#467 (architecture), `cnos.core/skills/agent/activate` + `attach` skills, `cn-sigma:.cn-sigma/spec/OPERATOR.md §"CDD role assignment"`, and `cnos:docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`.

**Status:** **pass**

**Evidence:**

- cnos#468 / label-doctrine: 11 total occurrences (manifest: 5; prompt: 6).
- cnos#467 / wake-orchestration: 13 total occurrences (manifest: 7; prompt: 6).
- activate skill + attach skill: 9 total occurrences (manifest: 5 — both cited in `cross_references.consumed_skills`; prompt: 4 — both cited in §"Cross-references" + invoked in §"Identity and activation").
- `cn-sigma:.cn-sigma/spec/OPERATOR.md`: 2 total occurrences (manifest: 1 — in `cross_references.adjacent_operator_doctrine`; prompt: 1 — in §"Cross-references" with explicit "CDD role assignment" section anchor).
- `AGENT-ACTIVATION-LOG-v0`: 10 total occurrences (manifest: 5; prompt: 5).
- Structured `cross_references` object in manifest carries explicit keys: `architecture`, `predecessors`, `consumed_skills`, `consumed_conventions`, `adjacent_operator_doctrine`, `downstream_consumers` — making the cross-reference graph machine-consumable per AC1 contract skill §3.6.

### AC7 — Relationship to existing claude-wake.yml documented

**Invariant (verbatim from #470):** Declaration includes a short section naming `.github/workflows/claude-wake.yml` as the existing production wake, noting that this provider declaration supersedes it once Sub 3 (renderer) lands, and naming Sub 3 as the cutover point. The existing wake is NOT touched in this sub.

**Status:** **pass**

**Evidence:**

- Manifest `superseded_substrate_artifact: ".github/workflows/claude-wake.yml"` (the AC7 carve-out field defined in AC1 contract skill §2.2).
- Manifest `relationship_to_substrate` carries the AC7 prose section: names the existing wake as substrate-named and owned by no package, notes the collapsed admin/dispatch boundary in the existing wake, declares this declaration as the package-owned replacement, names Sub 3 of cnos#467 (`cn wake install`) as the cutover point, names the rendered artifact (`.github/workflows/cnos-agent-admin.yml`), states "cnos#470 (this cycle) does NOT touch claude-wake.yml — the cutover happens at Sub 3 or later", and references AC7's byte-identical invariant.
- Mechanical proof of the byte-identical invariant: `git diff origin/main..HEAD -- .github/workflows/claude-wake.yml | wc -l` returns 0 at HEAD `2c7b1437` (the implementation SHA; re-validated immediately before this AC evidence block).
- Scope discipline mechanical proof: `git diff --name-only origin/main..HEAD | grep -vE '^(src/packages/cnos\.core/|\.cdd/unreleased/470/)' | wc -l` returns 0. The cycle touches only `src/packages/cnos.core/` (new files) and `.cdd/unreleased/470/` (γ-scaffold + α self-coherence). Nothing under `.github/` is touched.


