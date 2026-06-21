# γ scaffold — cycle/470

**Issue:** [cnos#470](https://github.com/usurobor/cnos/issues/470) — `agent-admin/wake-provider: cnos.core agent-admin wake (admin-only, no cell execution)`. Sub 2 of cnos#467 (master tracker `agent/wake-orchestration`); builds on Sub 1 (cnos#468 — label doctrine, merged at `c0048bef`).

**Mode:** **design-and-build** (per #470 mode declaration). Provider declaration is package content; prompt template is content; contract skill is content. No code, no `cn` binary edits, no `.github/workflows/` edits. A design call (declaration form — `commands/install-wake/` vs `orchestrators/agent-admin/`) lands within the cycle; this scaffold pins the recommended form below (§Form choice). α may override **only** if it discovers a structural reason during authoring; the override must be recorded in `self-coherence.md` §Gap with the reason cited.

**Branch:** `cycle/470` from `origin/main` at `c0048befb89c9b4aa083dd5fdb5c6c5547966ab9` (verified by `git rev-parse HEAD` at γ scaffold time — this is the merge commit of PR #469 / cycle/468 which landed the label doctrine).

**Cycle execution mode:** **bootstrap (pre-dispatch δ/channel path).** The γ-interface session acts as bootstrap-δ; γ/α/β are spawned as separate sub-agents via the Agent tool; `.cdd/unreleased/470/` artifacts are the shared memory across roles (no chat-state continuity). The package-owned dispatch wakes that would normally claim a `dispatch:cell + protocol:cdd + status:todo` cycle do not exist yet — this issue is **part of building them** (Sub 2 builds the cnos.core agent-admin provider; Sub 3 builds the renderer that materializes it; Sub 4 builds the cnos.cdd dispatch wake provider; Sub 5 wires δ's wake-invoked mode; Sub 6 wires end-to-end smoke + cycle-complete artifact reading). A friction log will be appended to `gamma-closeout.md` capturing what context each spawned role needed to act without shared memory; that log feeds Sub 5's eventual dispatch-prompt template spec.

---

## Form choice (pinned by γ; α may override only with a recorded structural reason)

The issue offers two surfaces for the provider declaration entry: `cnos.core/commands/install-wake/` or `cnos.core/orchestrators/agent-admin/`. γ pins **`cnos.core/orchestrators/agent-admin/`** with a new sibling manifest file `wake-provider.json` (carrying a new schema `cn.wake-provider.v1`) — NOT overloading the existing `cn.orchestrator.v1` workflow-engine IR schema. The prompt template lands as `cnos.core/orchestrators/agent-admin/prompt.md` (the renderer in Sub 3 will inline-substitute its body into the rendered workflow's `prompt:` field).

Rationale:

1. **Semantic separation from Sub 3.** `commands/install-wake/` would be a shell command named `cn-install-wake` — but executing the install IS Sub 3's renderer. Placing the declaration under a command name that implies execution conflates the declaration (Sub 2) with the renderer (Sub 3). The declaration is *data + prompt*, not an action verb.
2. **`orchestrators/` is the package-content-class for named long-form work-shape declarations.** The existing `orchestrators/daily-review/orchestrator.json` is *one* shape within `orchestrators/` (the `cn.orchestrator.v1` mechanical-workflow IR). Nothing in the package-content-class rules forbids a sibling shape (a wake-provider declaration). Schemas are file-versioned, not directory-versioned.
3. **No new package-content-class.** Per cnos#467 active design constraint, wake provision lands as `commands/` or `orchestrators/` first; a new `wakes/` class is explicitly deferred. `orchestrators/agent-admin/` honors that constraint.
4. **Forward compatibility with Sub 4.** Sub 4 (cnos.cdd dispatch wake provider) will pattern-copy from this form. `cnos.cdd/orchestrators/cdd-dispatch/wake-provider.json + prompt.md` parallels cleanly. The `commands/install-wake/` shape would force every protocol package to ship its own `cn-install-wake` script — duplicative and substrate-coupled.

**If α discovers a structural reason to use `commands/install-wake/` instead** (e.g., an undocumented `cn` invariant that orchestrators must conform to `cn.orchestrator.v1`), α records the reason in `self-coherence.md` §Gap and proceeds with the alternative form. γ's pin is not absolute; γ has not yet seen the Go-side dispatch code that consumes either directory class.

---

## Expected touched surfaces (α's prediction; α records actual in self-coherence §CDD Trace step 6)

| File | Purpose | New / edit |
|---|---|---|
| `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` | AC1 — wake-provider declaration **contract** skill (substrate-agnostic; required + optional fields; what the package declares vs what the renderer materializes; concrete enough for Sub 4 to author cnos.cdd's dispatch wake provider against it alone) | new |
| `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` | AC2 — agent-admin **provider declaration** entry (identity, responsibilities, input/output contract, allowed/disallowed surfaces, defer-path) conforming to AC1's contract; schema `cn.wake-provider.v1` | new |
| `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` | AC3 — admin-only **prompt template** (substrate-agnostic; admin responsibilities enumerated, cell execution forbidden, defer-path explicit, cnos#468 label-doctrine cited) | new |
| `src/packages/cnos.core/orchestrators/agent-admin/README.md` *(optional)* | Cross-references (AC6) + relationship-to-claude-wake.yml section (AC7); may instead live in `wake-provider.json` description fields and/or `prompt.md` header — α's call | optional |
| `.cdd/unreleased/470/self-coherence.md` | α's review-readiness signal + CDD Trace | new (α writes incrementally per `alpha/SKILL.md` §2.5) |
| `.cdd/unreleased/470/beta-review.md` | β's R1+ verdict + findings | new (β writes per `beta/SKILL.md`) |
| `.cdd/unreleased/470/alpha-closeout.md` | α's post-merge cycle narrative | new (α writes after β merge per re-dispatch) |
| `.cdd/unreleased/470/beta-closeout.md` | β's post-merge release context | new (β writes at merge) |
| `.cdd/unreleased/470/gamma-closeout.md` | γ's cycle closure declaration + friction log | new (γ writes after β passes; out of scope for α/β prompts) |
| `.cdd/unreleased/470/gamma-scaffold.md` | this file | (already written) |

**Explicitly NOT touched:**
- `.github/workflows/claude-wake.yml` — AC7 invariant: the existing hand-written wake remains byte-identical on `cycle/470`. Cutover happens at Sub 3 or later.
- `.github/workflows/*` — no GitHub Actions YAML emission in this sub (substrate is rendered, not hand-edited).
- Any package other than `cnos.core` — Sub 4 builds the cnos.cdd dispatch wake provider; not this sub.
- `cn` Go binary (`src/go/`) — `cn wake install` is Sub 3.
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` — δ wake-invoked mode is Sub 5.
- `cn-sigma:.cn-sigma/spec/OPERATOR.md` — out of scope; cited from this provider, not edited.

---

## Implementation contract (pinned by δ at dispatch; verified at this scaffold)

| Axis | Pinned value |
|---|---|
| Language | Markdown (skill + prompt template) + JSON (declaration manifest with schema `cn.wake-provider.v1`). No Go, no shell, no YAML. |
| CLI integration target | **None for this sub.** `cn wake install` is Sub 3. The declaration must be parseable as data; no `cn` subcommand change. |
| Package scoping | `src/packages/cnos.core/` only. New files under `cnos.core/skills/agent/wake-provider/` (contract skill) and `cnos.core/orchestrators/agent-admin/` (provider declaration entry). |
| Existing-binary disposition | N/A — no binaries touched. `cnos.core/cn.package.json`'s `commands` map is NOT edited; the new entry lives under `orchestrators/`, which is data the future renderer consumes, not a registered command. |
| Runtime dependencies | None. The declaration is static data + a markdown prompt template. |
| JSON/wire contract | The provider declaration manifest format MUST be machine-consumable by the Sub 3 renderer (`cn wake install`). α declares the schema with a `"schema": "cn.wake-provider.v1"` field so the renderer can dispatch by schema. Required vs optional fields are enumerated by AC1's contract skill. The contract skill is the canonical schema definition; the manifest is one instance of it. |
| Backward compat | `.github/workflows/claude-wake.yml` MUST remain byte-identical on `cycle/470`. Verified by `git diff origin/main -- .github/workflows/claude-wake.yml` returning 0 lines. The existing wake is **temporary**; Sub 3 cutover supersedes it, not Sub 2. |

---

## AC mapping (oracles + surfaces; α records actual evidence in `self-coherence.md` §ACs)

| AC | Invariant (per #470) | Oracle (mechanical) | Surface |
|----|----------------------|---------------------|---------|
| AC1 | Wake-provider declaration contract skill exists at cnos.core; defines required + optional fields, substrate-rendering target, split between "what the package declares" vs "what the renderer materializes"; concrete enough for cnos.cdd's dispatch wake provider (Sub 4) to be authored against it alone. Cites cnos#467 architecture + cnos#468 label doctrine. | `test -f src/packages/cnos.core/skills/agent/wake-provider/SKILL.md`; `grep -E "^#+ " SKILL.md` shows contract sections (required fields, optional fields, rendering target, what-package-declares-vs-renderer-materializes split); `grep -c "cnos#467\|wake-orchestration\|label-doctrine\|cnos#468" SKILL.md` ≥ 2; skill frontmatter conforms to `cnos.core/skills/skill/SKILL.md` format (`name`, `description`, `artifact_class`, `visibility`, `parent`, `triggers`, `scope`). | `cnos.core/skills/agent/wake-provider/SKILL.md` |
| AC2 | Agent-admin wake provider declaration entry exists at the form γ pinned (or α's documented override); carries identity, responsibilities (admin-only), input/output contract, allowed/disallowed surfaces, defer-path, prompt template. Substrate-agnostic — no GitHub Actions YAML in the declaration itself. | `test -f src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` (γ-pinned form) OR `test -f src/packages/cnos.core/commands/install-wake/...` (α-override form); JSON parses (`jq . wake-provider.json > /dev/null`); manifest schema declared via `"schema": "cn.wake-provider.v1"`; conforms to AC1's contract skill (required fields present); `grep -c "github\|workflow\|yaml\|GITHUB_TOKEN" wake-provider.json prompt.md` MUST equal 0 (substrate-agnostic — Sub 3 emits GH-Actions YAML, not Sub 2); admin-only constraint visible (the manifest carries an `admin_only: true` field OR equivalent named field per AC1's schema). | `cnos.core/orchestrators/agent-admin/wake-provider.json` + `prompt.md` (γ-pinned form) |
| AC3 | Prompt template enforces admin-only constraint: enumerates admin responsibilities (activate, attach, channel sync, status reporting, issue creation/refinement, label application), explicitly forbids cell execution, specifies defer-path for cell-shaped directives (defer to relevant `protocol:{P}` dispatch wake if installed; surface to operator if not), cites cnos#468 label-doctrine for label-application discipline. | Reading `prompt.md`: `grep -E "activate\|attach" prompt.md` ≥ 2; `grep -iE "do not execute cell\|never execute cell\|no cell execution\|MUST NOT execute" prompt.md` ≥ 1; `grep -iE "defer\|dispatch wake" prompt.md` ≥ 2 covering the defer-path; `grep -c "label-doctrine\|cnos#468" prompt.md` ≥ 1; `grep -iE "status report\|channel\|issue creat\|label appl" prompt.md` ≥ 3 covering the enumerated responsibilities. | `cnos.core/orchestrators/agent-admin/prompt.md` |
| AC4 | Input + output contracts documented: triggers (schedule cron slots; optional `issues: opened` for `claude-wake`-titled trigger issues) as inputs; terminal-report shape (channel log entry per `AGENT-ACTIVATION-LOG-v0.md`; frontmatter fields including `class:` taxonomy `heartbeat / substantive / directive-out / cycle-complete`; cursor advancement) as outputs. Precise enough that Sub 3 (renderer) can produce a working workflow + Sub 6 (cycle-complete reading) can extend the `class:` taxonomy. | Reading the declaration (`wake-provider.json` + `prompt.md`): input contract section names triggers (cron + optional issues-opened-with-title-match); output contract section cites `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` and enumerates the four `class:` values; cursor advancement named explicitly. `grep -c "AGENT-ACTIVATION-LOG-v0\|activation log\|conventions" wake-provider.json prompt.md` ≥ 1; `grep -c "heartbeat\|substantive\|directive-out\|cycle-complete" wake-provider.json prompt.md` ≥ 4. | `cnos.core/orchestrators/agent-admin/wake-provider.json` + `prompt.md` |
| AC5 | Allowed + disallowed admin surfaces enumerated: allowed = `.cn-{agent}/logs/`, `.cn-{agent}/state/` (subject to per-install policy), issue/PR comments, label application per #468; disallowed = code/test/doc files outside `.cn-{agent}/` and `.cdd/`, `.github/workflows/`, branch protection, repo settings, **cell execution**. | Reading the declaration: `allowed_surfaces` and `disallowed_surfaces` fields (or named-prose equivalents) both present; both lists are non-empty; `cell execution` (or equivalent phrase) explicit in disallowed; `.github/workflows/` explicit in disallowed; `.cn-{agent}/logs/` + `.cn-{agent}/state/` explicit in allowed. `grep -c "\.cn-\|cn-{agent}\|cn-\\\\{agent\\\\}" wake-provider.json prompt.md` ≥ 2. | `cnos.core/orchestrators/agent-admin/wake-provider.json` |
| AC6 | Cross-references present: cnos#468 (label doctrine), cnos#467 (architecture), `cnos.core/skills/agent/activate` + `attach` skills, `cn-sigma:.cn-sigma/spec/OPERATOR.md §"CDD role assignment"`, `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md`. | `grep -c "cnos#468\|label-doctrine" {declaration-files}` ≥ 1; `grep -c "cnos#467\|wake-orchestration" {declaration-files}` ≥ 1; `grep -c "activate/SKILL\|attach/SKILL\|cnos.core/skills/agent/activate\|cnos.core/skills/agent/attach" {declaration-files}` ≥ 2; `grep -c "cn-sigma:.cn-sigma/spec/OPERATOR\|cn-sigma.*OPERATOR" {declaration-files}` ≥ 1; `grep -c "AGENT-ACTIVATION-LOG-v0\|activation log convention" {declaration-files}` ≥ 1. | `wake-provider.json` + `prompt.md` (whichever carries the reference; both is fine) |
| AC7 | Relationship to `.github/workflows/claude-wake.yml` documented: short section naming it as existing production wake, noting this provider supersedes it once Sub 3 (renderer) lands, naming Sub 3 as the cutover point. The existing wake is NOT touched in this sub. | Reading the declaration: section heading or named field mentioning `claude-wake.yml` exists; cutover-via-Sub-3 stated; `git diff origin/main -- .github/workflows/claude-wake.yml` returns 0 lines on `cycle/470` HEAD. | declaration body + absence of changes under `.github/` |

**Mechanical gate (run by α at pre-review-gate row 9 polyglot re-audit and by β at row 3 pre-merge non-destructive merge-test):**

```bash
# Mechanical proof of AC7 (no .github/ changes):
git diff --name-only origin/main..HEAD -- .github/ | wc -l  # MUST be 0

# Mechanical proof of substrate-agnostic AC2:
grep -ciE 'github|workflow|yaml|GITHUB_TOKEN|claude-code-action|runs-on' src/packages/cnos.core/orchestrators/agent-admin/*.json src/packages/cnos.core/orchestrators/agent-admin/*.md 2>/dev/null  # MUST be 0 except for the AC7 reference section which may name "claude-wake.yml" and ".github/workflows/" as the surface being-superseded-by-Sub-3

# Mechanical proof of scope discipline:
git diff --name-only origin/main..HEAD | grep -vE '^(src/packages/cnos\.core/|\.cdd/unreleased/470/)' | wc -l  # MUST be 0
```

---

## α dispatch prompt

(Below is the full self-contained prompt δ feeds to the spawned α session. α has zero shared memory with γ or β. The prompt carries everything α needs.)

---

```
You are α (alpha) for CDD cycle cnos#470.

# Your one job
Implement cnos#470 (`agent-admin/wake-provider: cnos.core agent-admin wake (admin-only, no cell execution)`). Produce the review-ready artifact set on `cycle/470`. Write `.cdd/unreleased/470/self-coherence.md` incrementally per the α SKILL §2.5 (one section per commit, pushed). Signal review-readiness by appending the `## Review-readiness` section to `self-coherence.md` and pushing. **Do not spawn anyone. Do not dispatch β.** δ (the operator-routed bootstrap session) will dispatch β after you signal review-readiness.

# Environment
- Working directory: `/home/user/cnos` (the cnos repository working tree)
- Branch: **`cycle/470`** (created by γ from `origin/main` at SHA `c0048befb89c9b4aa083dd5fdb5c6c5547966ab9` — the merge commit of cnos#469 / cycle/468 which landed cnos#468 label doctrine)
- Repository: `usurobor/cnos`
- Issue URL: https://github.com/usurobor/cnos/issues/470
- Master tracker URL: https://github.com/usurobor/cnos/issues/467 (read the "Foundational architecture — package-owned wake providers (authoritative)" section)
- Sub 1 (merged predecessor): https://github.com/usurobor/cnos/issues/468

# Git identity
Before any commit, set:
  git config user.email "alpha@cdd.cnos"
  git config user.name "alpha"
Verify with `git config --get user.email`. (Per `operator/SKILL.md §Git identity for role actors` + `alpha/SKILL.md §2.6 row 14`.)

# Cycle execution mode (read before claiming any rule does not apply)
This cycle is executed through the **pre-dispatch δ/channel path**: γ-interface session acts as bootstrap-δ, spawns α (you) via the Agent tool, then will spawn β after your review-readiness signal. The package-owned dispatch wakes that would normally claim this cycle do not exist yet — Sub 2 (this cycle) is part of building them. β is NOT yet spawned; do not poll for β verdicts during your authoring run. δ will run β after you exit on review-readiness. α exits after signaling review-readiness (sequential bounded dispatch per `cnos.cds/skills/cds/CDS.md` §"Field 6: Actor collapse rule"); δ will re-dispatch you after β merges for `alpha-closeout.md`.

# Skills to load (in this order, before any implementation commit)
Tier 1:
1. `src/packages/cnos.cdd/skills/cdd/CDD.md` (if present; else skip — it may not exist)
2. `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` (your role)

Lifecycle sub-skills as work requires:
3. `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` (to verify ACs are independently testable)
4. `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` (this is design-and-build — load before any wake-provider.json schema decision; mark "not required" only with explicit justification per `alpha/SKILL.md §2.2`)

Tier 3 (issue-specific):
5. `src/packages/cnos.core/skills/agent/activate/SKILL.md` (the prompt template you author will direct the wake to invoke this)
6. `src/packages/cnos.core/skills/agent/attach/SKILL.md` (same)
7. `src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` (cnos#468 merged predecessor; the prompt template cites this for label-application discipline; the manifest's allowed/disallowed surfaces conform to its boundary)
8. `src/packages/cnos.core/skills/skill/SKILL.md` (the skill-format meta-skill; AC1's contract skill conforms to it)
9. `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` (the terminal-report shape AC4 names)
10. `src/packages/cnos.core/cn.package.json` + `ls src/packages/cnos.core/commands/ src/packages/cnos.core/orchestrators/` (existing package structure for shape reference)
11. `.github/workflows/claude-wake.yml` (the existing wake your provider supersedes at Sub 3 cutover; READ ONLY — AC7 invariant says you do NOT modify this file)

Also re-read the γ scaffold for full context including form-choice rationale and AC mapping:
12. `.cdd/unreleased/470/gamma-scaffold.md`

# Implementation contract (pinned by δ; α MUST NOT improvise these axes)

| Axis | Pinned value |
|---|---|
| Language | Markdown (skill + prompt template) + JSON (declaration manifest, schema `cn.wake-provider.v1`). No Go, no shell, no YAML. |
| CLI integration target | None for this sub (`cn wake install` is Sub 3). No `cn` subcommand change. No edits to `src/go/`. No edits to `cnos.core/cn.package.json`'s `commands` map. |
| Package scoping | `src/packages/cnos.core/` only. Files under `cnos.core/skills/agent/wake-provider/` and `cnos.core/orchestrators/agent-admin/`. |
| Existing-binary disposition | N/A — no binaries. |
| Runtime dependencies | None — static data + markdown. |
| JSON/wire contract | The provider declaration manifest format MUST be machine-consumable by the Sub 3 renderer. Declare `"schema": "cn.wake-provider.v1"` as the first field of the manifest. Required + optional fields are enumerated by AC1's contract skill (which you also author this cycle); the manifest is one instance of that contract. |
| Backward compat | `.github/workflows/claude-wake.yml` MUST remain byte-identical on `cycle/470`. Mechanical check before signaling review-readiness: `git diff origin/main..HEAD -- .github/workflows/claude-wake.yml | wc -l` MUST return 0. |

α MUST NOT change any pinned axis. If you find a structural reason to change one, surface to γ via `.cdd/unreleased/470/gamma-clarification.md` (write it, commit it on `cycle/470`, and stop) before making any contradicting commit.

# Form choice (γ-pinned; override only with recorded structural reason)
γ pins the provider declaration entry at **`src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json`** (manifest) + **`src/packages/cnos.core/orchestrators/agent-admin/prompt.md`** (prompt template). NOT at `cnos.core/commands/install-wake/`. Rationale is in `.cdd/unreleased/470/gamma-scaffold.md` §"Form choice"; read it before authoring. You MAY override to `commands/install-wake/` ONLY if you discover a structural reason in `cn`'s dispatch logic (e.g., `orchestrators/` is constrained to a single `orchestrator.json` filename or schema). If you override, record the reason in `self-coherence.md` §Gap citing the file + line that forces the override.

# Acceptance criteria (verbatim from cnos#470 — AC1 through AC7)

[Paste from cnos#470 issue body §"Acceptance criteria" verbatim: AC1 Wake-provider declaration contract skill, AC2 Agent-admin wake provider declaration entry, AC3 Prompt template enforces admin-only constraint, AC4 Input + output contracts documented, AC5 Allowed + disallowed admin surfaces enumerated, AC6 Cross-references present, AC7 Relationship to existing claude-wake.yml documented. Read the issue body via `mcp__github__issue_read` with method=get, owner=usurobor, repo=cnos, issue_number=470 — or via the GitHub web URL if MCP github is unavailable. Each AC's Invariant / Oracle / Surface block must be read in full; the γ scaffold's AC mapping table is a navigation aid, not a substitute.]

# Expected surfaces (γ scaffold §"Expected touched surfaces" — re-read it; this is summary only)

NEW:
- `src/packages/cnos.core/skills/agent/wake-provider/SKILL.md` (AC1 contract skill)
- `src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` (AC2 declaration, schema `cn.wake-provider.v1`)
- `src/packages/cnos.core/orchestrators/agent-admin/prompt.md` (AC3 admin-only prompt template)
- `src/packages/cnos.core/orchestrators/agent-admin/README.md` (optional; for AC6 cross-refs + AC7 relationship section if you prefer them separated from the manifest)
- `.cdd/unreleased/470/self-coherence.md` (your review-readiness signal)

NOT TOUCHED (refusal conditions — if you find yourself editing these, STOP and re-read this prompt):
- `.github/workflows/claude-wake.yml` (AC7 byte-identical invariant)
- Anything else under `.github/`
- Any package other than `cnos.core` (esp. `cnos.cdd/`, `cnos.cdr/`, `cnos.kata/`)
- `src/go/` (no `cn wake install` — Sub 3)
- `src/packages/cnos.cdd/skills/cdd/delta/SKILL.md` (Sub 5)
- `src/packages/cnos.core/cn.package.json`'s `commands` map (the new entry is data under `orchestrators/`, not a registered command)
- `cn-sigma:` anything (different repo)

# Refusal conditions (binding)
α MUST refuse and surface to γ via `.cdd/unreleased/470/gamma-clarification.md` (commit on cycle/470, then stop) if:
- The cycle requires editing `.github/workflows/claude-wake.yml` to satisfy any AC (it does not; AC7 says the opposite).
- The cycle requires implementing `cn wake install` or any renderer (those are Sub 3).
- The cycle requires implementing any dispatch claim protocol (cnos#454; not this sub's concern).
- The cycle requires authoring a CDD dispatch wake (Sub 4) or extending δ (Sub 5).
- A pinned implementation-contract axis appears unsatisfiable from this scope.
- An AC oracle from the issue body appears unsatisfiable from the pinned form / scope.

If you find a *behavioral* gap in the declaration design (e.g., a field the renderer in Sub 3 will need that this scaffold did not anticipate), do NOT defer — design it in (this is design-and-build mode) and record the design call in `self-coherence.md` §Gap with the reasoning. The cycle's job is to ship a declaration the Sub 3 renderer can consume; if it lacks a field the renderer needs, the declaration is incomplete.

# Process

1. Set git identity (above).
2. `git fetch origin && git switch cycle/470` (the branch exists — γ created it).
3. Read everything in §"Skills to load" above.
4. Re-read `.cdd/unreleased/470/gamma-scaffold.md` (this scaffold) in full.
5. Read cnos#470 issue body in full (verbatim ACs).
6. Read cnos#467 master tracker — specifically §"Foundational architecture — package-owned wake providers (authoritative)" and §"Acceptance criteria (wave-level)".
7. Begin writing `.cdd/unreleased/470/self-coherence.md` incrementally per α SKILL §2.5:
   - §Gap (issue, version/mode = `design-and-build`, form-choice acknowledgment or documented override)
   - §Skills (Tier 1/2/3 loaded)
   - Then implement: AC1 contract skill → AC2 manifest → AC3 prompt template → AC4-AC7 fields in manifest/prompt as appropriate
   - §ACs (AC-by-AC evidence pointing at the diff)
   - §Self-check (did α push ambiguity onto β? is every claim backed by diff evidence?)
   - §Debt (known debt; "none" is acceptable if true)
   - §CDD Trace (steps 0-7)
   - §Review-readiness (the signal that exits α to β)
8. Before §Review-readiness: run the mechanical gates from the γ scaffold §"AC mapping" → "Mechanical gate" block. ALL must pass.
9. Run α SKILL §2.6 pre-review-gate rows 1-15. Row 9 (polyglot re-audit) for this cycle = JSON validity + Markdown grep + AC oracle re-run.
10. Append §Review-readiness with the implementation SHA (per α SKILL §2.6 "SHA convention for readiness signal").
11. Push to `origin/cycle/470`. Exit. Do NOT poll for β; δ runs β.

# Output artifacts (you write)
- `.cdd/unreleased/470/self-coherence.md` (incremental, one section per commit)
- Implementation files under `src/packages/cnos.core/` per §"Expected surfaces"

You do NOT write:
- `beta-review.md` (β writes it)
- `alpha-closeout.md` (you write it later, after δ re-dispatches you post-β-merge)
- `gamma-closeout.md` (γ writes it)
- Any `cdd-iteration.md` (γ writes it if findings warrant)

# Commit message guidance
- One commit per self-coherence section: `α-470 self-coherence: §Gap`, `α-470 self-coherence: §Skills`, etc.
- Implementation commits: `α-470: wake-provider contract skill (AC1)`, `α-470: agent-admin wake-provider.json (AC2)`, `α-470: agent-admin prompt.md (AC3)`, etc.
- Push after each commit. The cycle branch is the coordination surface.

# Stop condition
Stop when §Review-readiness is appended to `self-coherence.md` and pushed to `origin/cycle/470`. Do NOT wait for β. Do NOT message anyone. Return a short summary: which commits landed, where review-readiness was signaled, any known debt declared in §Debt, any γ-clarification you needed to file. δ will read your summary and dispatch β.
```

---

## β dispatch prompt

(Below is the full self-contained prompt δ feeds to the spawned β session after α signals review-readiness. β has zero shared memory with γ or α.)

---

```
You are β (beta) for CDD cycle cnos#470.

# Your one job
Review α's implementation of cnos#470 (`agent-admin/wake-provider: cnos.core agent-admin wake (admin-only, no cell execution)`) on `cycle/470`. Write `.cdd/unreleased/470/beta-review.md` with a per-round verdict (R1, R2, ...) and binding findings. If verdict = APPROVE, merge to `main` (per β SKILL §"Phase map" → β merge); if RC, request changes and exit (α will re-dispatch). After merge, write `.cdd/unreleased/470/beta-closeout.md`. **Do not spawn anyone. Do not dispatch γ or α.** δ (the operator-routed bootstrap session) will dispatch α for close-out after merge.

# Environment
- Working directory: `/home/user/cnos`
- Branch under review: **`cycle/470`** (γ created from `origin/main` at SHA `c0048befb89c9b4aa083dd5fdb5c6c5547966ab9`; α implements on it)
- Branch to merge into: `main`
- Repository: `usurobor/cnos`
- Issue URL: https://github.com/usurobor/cnos/issues/470
- Master tracker URL: https://github.com/usurobor/cnos/issues/467
- PR URL: {PR may or may not exist; β operates against `cycle/470` branch directly per β SKILL §"Phase map" — β reads `origin/cycle/{N}` and merges via `git merge`. If δ has opened a PR for tracking, the URL will be in the bootstrap message δ passes to you; otherwise β operates branch-direct.}

# Git identity
Before any commit (including the merge commit), set:
  git config user.email "beta@cdd.cnos"
  git config user.name "beta"
Verify with `git config --get user.email`. Per β SKILL §"Pre-merge gate" row 1 (Identity truth) — if you use a worktree for merge-testing, use `git config --worktree user.email "beta@cdd.cnos"` to avoid leaking to shared config.

# Cycle execution mode (read before claiming any rule does not apply)
This cycle is executed through the **pre-dispatch δ/channel path**: γ-interface session acts as bootstrap-δ, spawned α (already exited at review-readiness), now spawning β (you). The package-owned dispatch wakes that would normally claim this cycle do not exist yet — Sub 2 IS part of building them. β operates in sequential bounded dispatch (per `cnos.cds/skills/cds/CDS.md` §"Field 6: Actor collapse rule"). After your merge + `beta-closeout.md`, you exit; δ re-dispatches α for `alpha-closeout.md` and then γ writes `gamma-closeout.md`.

# Skills to load (in this order, before reviewing)
Tier 1:
1. `src/packages/cnos.cdd/skills/cdd/CDD.md` (if present)
2. `src/packages/cnos.cdd/skills/cdd/beta/SKILL.md` (your role)
3. `src/packages/cnos.cdd/skills/cdd/review/SKILL.md`
4. `src/packages/cnos.cdd/skills/cdd/release/SKILL.md` (for the merge mechanics — but NOT the full release; tag/deploy is δ's release boundary, not yours, and this cycle does not cut a release anyway)

Tier 3 (issue-specific):
5. `src/packages/cnos.core/skills/agent/activate/SKILL.md` (α's prompt template will reference this)
6. `src/packages/cnos.core/skills/agent/attach/SKILL.md` (same)
7. `src/packages/cnos.core/skills/agent/label-doctrine/SKILL.md` (cnos#468 — the prompt template cites this; AC6 cross-reference)
8. `src/packages/cnos.core/skills/skill/SKILL.md` (the skill-format meta-skill α's AC1 contract skill must conform to)
9. `docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` (AC4 terminal-report shape)
10. `.github/workflows/claude-wake.yml` (AC7 byte-identical invariant — confirm unchanged)

Also re-read for context:
11. `.cdd/unreleased/470/gamma-scaffold.md` (γ's scaffold — AC mapping table is your review checklist; form-choice rationale is γ's pin and α may have overridden it; verify which path the diff took)
12. `.cdd/unreleased/470/self-coherence.md` (α's review-readiness signal — CDD Trace, AC evidence, known debt)

# Issue context (for verbatim ACs)
Read cnos#470 issue body in full via `mcp__github__issue_read` (owner=usurobor, repo=cnos, issue_number=470, method=get). ACs AC1-AC7 are the binding gate. β SKILL Rule 6 anchors oracle evidence on code (the diff), not doc — re-grep α's actual artifacts; do not trust α's claims without verifying against the diff.

# Pre-merge gate (β SKILL §"Pre-merge gate"; binding before `git merge`)
Run all four rows. If any row fails, RC with severity D on that row's failure mode.

Row 1 — Identity truth: `git config user.email` returns `beta@cdd.cnos`.

Row 2 — Canonical-skill freshness: `git fetch --verbose origin main && git rev-parse origin/main` returns SHA matching session-start snapshot. If `origin/main` advanced mid-cycle, re-load Tier-1 skills and re-evaluate.

Row 3 — Non-destructive merge-test (this cycle's specific validators):
- `git worktree add /tmp/cnos-merge-470/wt origin/main && cd /tmp/cnos-merge-470/wt && git config --worktree user.name beta && git config --worktree user.email beta@cdd.cnos && git merge --no-ff --no-commit origin/cycle/470`
- Confirm zero unmerged paths.
- Validators specific to this cycle's surface:
  - `./tools/validate-skill-frontmatter.sh` IF the cycle adds/modifies any SKILL.md frontmatter (it does — AC1's `wake-provider/SKILL.md` is new)
  - `jq . src/packages/cnos.core/orchestrators/agent-admin/wake-provider.json` (or α's override path) — JSON parses
  - `cn-cdd-verify` IF present (for `.cdd/` artifacts)
  - Skill-format conformance check: the new `wake-provider/SKILL.md` frontmatter has `name`, `description`, `artifact_class`, `visibility`, `parent`, `triggers`, `scope` per `src/packages/cnos.core/skills/skill/SKILL.md`
- Mechanical AC7 check: `git diff origin/main HEAD -- .github/workflows/claude-wake.yml | wc -l` MUST be 0
- Mechanical scope check: `git diff --name-only origin/main HEAD | grep -vE '^(src/packages/cnos\.core/|\.cdd/unreleased/470/)' | wc -l` MUST be 0
- Mechanical substrate-agnostic check (AC2): `grep -ciE 'github|workflow|yaml|GITHUB_TOKEN|claude-code-action|runs-on' src/packages/cnos.core/orchestrators/agent-admin/*.json src/packages/cnos.core/orchestrators/agent-admin/*.md 2>/dev/null` MUST be 0 except for the AC7 reference section which legitimately names `claude-wake.yml` / `.github/workflows/` as the surface being superseded.
- Tear down worktree.

Row 4 — γ artifact completeness: `git ls-tree -r origin/cycle/470 .cdd/unreleased/470/gamma-scaffold.md` returns a non-empty line. (γ scaffold should exist; it was written before α dispatch.)

# Review approach (β SKILL Rule 6 + Rule 7)

For each AC (AC1-AC7 per cnos#470 body):
1. Read the AC's Invariant / Oracle / Surface verbatim from the issue body.
2. Re-grep the cycle's diff (NOT α's `self-coherence.md` claims) to verify the oracle passes against actual artifacts. β SKILL Rule 6: code is evidence, doc is hypothesis.
3. Per Rule 6a: if a γ-scaffold oracle regex is brittle (misses a real literal the code emits), widen β-side and record both the scaffold regex and the widened regex in `beta-review.md`. Do NOT request α to reshape implementation to fit a brittle scaffold regex.
4. Per Rule 6b: if a name overpromises (e.g., a field named `responsibilities` that's empty, or `defer_path` that's a placeholder), flag as name-overpromise — α decides rename-or-implement.

Rule 7 — Implementation-contract coherence (binding):
For each of the 7 axes in the dispatch contract (γ scaffold §"Implementation contract"), confirm the diff conforms. The diff hunks must map onto the pinned rows:
- Language: only `.md` + `.json` added; no `.go`, no `.sh`, no `.yml`.
- CLI integration target = None: `cnos.core/cn.package.json` `commands` map unchanged; no `src/go/` edits.
- Package scoping: only `src/packages/cnos.core/skills/agent/wake-provider/` and `src/packages/cnos.core/orchestrators/agent-admin/` new content; nothing else.
- Existing-binary disposition: N/A — no binary changes.
- Runtime dependencies: declaration has no runtime; if α added runtime, RC.
- JSON/wire contract: manifest declares `"schema": "cn.wake-provider.v1"`; AC1 contract skill is the canonical schema definition.
- Backward compat: `claude-wake.yml` byte-identical (mechanical AC7).

If any axis diverges without a `.cdd/unreleased/470/gamma-clarification.md` updating it, RC with D-severity `implementation-contract` classification.

# Verdict format

Write `.cdd/unreleased/470/beta-review.md` with this structure (per `review/SKILL.md` shape):

```
# β review — cycle/470

## R1 ({YYYY-MM-DD HH:MM UTC} — base origin/main: {SHA} — head origin/cycle/470: {SHA})

### Verdict: {APPROVE | REQUEST CHANGES | REJECT}

### Contract Integrity
{Identity, branch, base/head SHAs verified}

### Pre-merge gate
| Row | Result | Notes |
| 1 (Identity truth) | pass/fail | |
| 2 (Canonical-skill freshness) | pass/fail | |
| 3 (Non-destructive merge-test) | pass/fail | each validator named |
| 4 (γ artifact completeness) | pass/fail | |

### AC coverage
| AC | Status | Evidence (re-grepped from diff) | Notes |
| AC1 | pass/fail | {grep output, file path, line} | |
| AC2 | pass/fail | | |
| AC3 | pass/fail | | |
| AC4 | pass/fail | | |
| AC5 | pass/fail | | |
| AC6 | pass/fail | | |
| AC7 | pass/fail | (mechanical: 0 diff lines) | |

### Implementation-contract coherence (Rule 7)
| Axis | Status | Evidence |
| Language | pass/fail | |
| CLI integration | pass/fail | |
| ...

### Binding findings (if RC or REJECT)
- F1 (severity D/C/B/A; classification {protocol-compliance|implementation-contract|correctness|...}): {description; surface; oracle; remediation}
- F2 ...

### Approval line (if APPROVE)
"Approved at head {SHA}; merging via `git merge --no-ff origin/cycle/470` to main."
```

# On APPROVE: merge

If R1 (or RN) verdict is APPROVE and pre-merge gate passes:
1. From a clean checkout on `main`: `git fetch origin && git switch main && git pull --ff-only`.
2. `git merge --no-ff origin/cycle/470 -m "Merge pull request #470 from usurobor/cycle/470 — agent-admin wake-provider (admin-only)"` (or use a PR-merge if δ opened one; either is fine — β SKILL Rule 5 is "review → merge → β close-out" same session).
3. Verify the merge SHA. Push: `git push origin main`.
4. Write `.cdd/unreleased/470/beta-closeout.md` (per β SKILL Rule 5 §"Closure discipline" — review summary, implementation assessment, technical review, process observations; NOT release notes — this cycle does not cut a release; NOT a triage table — γ owns that).
5. Commit `beta-closeout.md` on `main`; push.
6. Exit. Do NOT spawn α or γ for close-out; δ re-dispatches.

# On REQUEST CHANGES: exit cleanly

If R1 (or RN) verdict is RC:
1. Ensure `beta-review.md` has the binding findings clearly enumerated with severities.
2. Push `beta-review.md` to `origin/cycle/470` (NOT to main — cycle artifacts live on cycle/470 until merge).
3. Exit. δ will re-dispatch α to fix findings; α appends a fix-round section to `self-coherence.md` and re-signals review-readiness; δ re-dispatches you for R2.

# Refusal conditions (binding)
β MUST refuse and surface to γ via `.cdd/unreleased/470/gamma-clarification.md` (commit on `cycle/470`, then exit without verdict) if:
- The operator (or anyone) attempts to instruct β directly on implementation, scope, or merge — β communicates only via γ-via-artifact-channel during a cycle (β SKILL Rule 1).
- The harness offers β a separate `claude/{slug}-{rand}`-style branch as the review/merge surface — refuse; review verdict lives on `origin/cycle/470` (β SKILL Rule 1).
- The cycle requires β to author the fix it judges — β does not author α's work (β SKILL Rule 1: independence).
- Approve-with-follow-up is requested — β SKILL Rule 4: no "approve with follow-up" except explicitly named design-scope deferrals filed before merge.

# Stop condition
Stop when (a) merge completed + `beta-closeout.md` written and pushed, OR (b) `beta-review.md` with RC verdict pushed to `cycle/470`. In either case, return a short summary: round number, verdict, merge SHA (if APPROVE), findings (if RC), any γ-clarification you filed. δ reads your summary and either re-dispatches α (RC path) or re-dispatches α for close-out (APPROVE path).
```

---

## Non-goals (verbatim from cnos#470 §"Non-goals" + operator constraints)

From cnos#470:
- **Not implementing `cn wake install`** — Sub 3.
- **Not touching `.github/workflows/`** — out of scope; cutover at Sub 3 or later.
- **Not implementing any dispatch claim protocol** — cnos#454.
- **Not creating per-package wake provider templates** — only the agent-admin provider here.
- **Not implementing CDD dispatch wake** — Sub 4.
- **Not migrating existing labels** — separate consideration.

From operator constraints (binding for this scaffold):
- **MUST NOT** implement `cn wake install`, the wake renderer, any package dispatch wake, CDD dispatch runtime, claim protocol, cell runtime, or touch `.github/workflows/claude-wake.yml`.
- **MUST NOT** edit `cnos.core/cn.package.json`'s `commands` map (the new declaration is data under `orchestrators/`, not a registered command).
- **Release notes are NOT automatic.** γ's eventual closeout records *candidate* release notes if the change warrants. γ does NOT manufacture `RELEASE.md` or cut a release for this cycle. δ does not tag.
- If the provider declaration is structurally complete without touching `.github/`, do not touch it.

## Refusal conditions for α (consolidated; mirrored from α dispatch prompt for visibility)

α MUST refuse (via `gamma-clarification.md` and stop) when:
1. The cycle requires editing `.github/workflows/claude-wake.yml`.
2. The cycle requires implementing `cn wake install` or any renderer.
3. The cycle requires authoring a CDD dispatch wake (Sub 4), δ wake-invoked-mode (Sub 5), or end-to-end smoke (Sub 6).
4. A pinned implementation-contract axis (γ scaffold table above) appears unsatisfiable from this scope.
5. An AC oracle from the issue body appears unsatisfiable from the pinned form / scope.
6. The cycle requires editing any path outside the §"Expected touched surfaces" allowed list.

α MAY (without refusal) design within the declaration shape — the manifest fields, the prompt template wording, the contract skill's required/optional split — and override γ's form-pin (with a structural reason recorded). This is `design-and-build` mode; the cycle's job is to design AND build, not to mechanically transcribe.

---

## Cycle scope sizing (per #470 — kept whole)

| Factor | Reading | Splitting signal? |
|---|---|---|
| (a) New code surface | 1 provider declaration entry (manifest + prompt template) + 1 contract skill; both in cnos.core | No |
| (b) Cross-module breadth | cnos.core only; no other package, no `.github/`, no `cn` binary | No |
| (c) Lifecycle span | doctrine (cycle) + small new artifact set | No |
| (d) MCA preconditions | partial; form-choice resolved by γ pin (this scaffold); manifest-field design is α's design call within the cycle | design-and-build |
| (e) Independent shippability | Yes — ships independently; Sub 3 consumes; Sub 4 parallels | n/a |

**Decision: keep whole (design-and-build).** Per cnos#470 §"Cycle scope sizing" decision.

---

Filed by γ@cnos (γ-spawned-as-sub-agent for cycle/470, bootstrap-δ path) on 2026-06-20.
