---
name: activation
description: Bootstrap an existing repository to operate under CDD — scaffold, version pin, labels, identity, CI, notifications, secrets, and operator registry.
artifact_class: skill
kata_surface: embedded
governing_question: What is the minimal correct sequence to bring an existing repo under CDD without re-deriving the protocol from scratch?
visibility: internal
parent: cdd
triggers:
  - activation
  - bootstrap
  - onboarding
  - tenant
scope: task-local
inputs:
  - existing repository with a primary branch
  - GitHub-hosted project (portability surface noted in §2)
  - operator with repository admin permissions
outputs:
  - populated .cdd/ directory scaffold
  - CDD version pin (.cdd/CDD-VERSION)
  - dispatch declaration (.cdd/DISPATCH)
  - release cadence declaration (.cdd/CADENCE)
  - operator handle registry (.cdd/OPERATORS)
  - MCA registry (.cdd/MCAs/INDEX.md)
  - vendored skill bundle or substitute declaration (.cdd/skills/)
  - GitHub issue labels (minimum set)
  - activation checklist completion evidence (§24 verification)
requires:
  - repository exists with a primary branch
  - operator has GitHub repo admin access
  - host supports labels, secrets, and CI (GitHub assumed; see §2 portability surface)
calls:
  - operator/SKILL.md
  - CDD.md
---

# CDD Activation

## §1 Purpose

This skill prescribes the one-time bootstrap sequence for bringing an existing repository under Coherence-Driven Development (CDD). Its purpose is to replace per-tenant improvisation with a canonical, repeatable activation path that every new repo can execute in order, verify step by step, and trust as authoritative. The skill closes the gap between "cdd protocol exists" and "this specific repo is operating under cdd."

Activation is not the same as running a cycle. Activation creates the infrastructure that cycles operate within: the directory scaffold, the version pin, the identity convention, the dispatch configuration, the CI surface, the notification interface, the operator registry. A repo that has not been activated may run cycles, but those cycles lack the prescribed governance surfaces — audit trails are incomplete, CI gates are absent, and tenants must re-derive what activation would have given them for free.

Activation runs once per repo. After activation, tenants run cycles against the scaffolded infrastructure. The `§24 Verification` command confirms activation succeeded; every new tenant should be able to run it against their repo and get a green result before dispatching their first cycle.

---

## §2 Pre-conditions

Before activation begins, three conditions must hold. First, the repository exists and has a designated primary branch (typically `main`). All merge targets and tag anchors during CDD cycles reference this branch; activation cannot proceed without a stable primary branch name. Second, the repository host supports issue labels, repository secrets (for notification adapters), and CI pipelines. This skill assumes GitHub; the portability surface is named in each section that uses a GitHub-specific mechanism (labels in §5, secrets in §11, CI in §9). Adapting to GitLab, Gitea, or another host requires substituting the equivalent platform primitive — the underlying cdd prescription is host-agnostic; only the mechanism names differ. Third, the activating operator has repository admin permissions sufficient to create labels, configure secrets, and add CI workflows. Activation is an admin-class operation; it cannot be delegated to a contributor with read-only or write-without-admin access.

**Portability surface summary:**

| CDD prescription | GitHub mechanism | Portability note |
|---|---|---|
| Issue labels (§5) | GitHub Labels API | GitLab: label resource; Gitea: issue labels |
| Repository secrets (§11) | GitHub repository secrets | GitLab: CI/CD variables; Gitea: actions secrets |
| CI pipelines (§9) | GitHub Actions workflows | GitLab CI YAML; Gitea Actions |
| Branch protection (§6) | GitHub branch rules | Platform-equivalent protection |

---

## §3 .cdd/ Scaffold

Activation begins by creating the `.cdd/` directory tree at the repository root. This directory is the permanent home for all CDD governance artifacts: version pins, dispatch configuration, cadence declaration, operator registry, MCA tracking, cycle close-outs, and the vendored skill bundle. Every path referenced in this skill and in `CDD.md` is relative to this root.

Create the following tree. Each entry below is the path relative to repo root and its one-line purpose:

```
.cdd/
  CDD-VERSION          # cnos commit SHA (and tag if present) that this repo follows — see §4
  DISPATCH             # dispatch configuration declaration (§5.1 or §5.2) — see §8
  CADENCE              # release cadence enum: versioned | rolling-docs | mixed — see §21
  OPERATORS            # table of authorized δ handles — see §19
  MCAs/
    INDEX.md           # table of in-flight multi-cycle architectures — see §15
  skills/
    README.md          # vendored skill bundle or origin-repo declaration — see §16
  releases/            # cycle close-out directories for version-tagged releases
  unreleased/          # cycle close-out directories for in-progress cycles
    {N}/               # one directory per cycle number (created by γ at dispatch)
  iterations/
    INDEX.md           # index of cdd-iteration.md findings across cycles
    cross-repo/
      README.md        # cross-repo trace bundle naming convention — see §13
  proposals/           # design proposals and architectural decision records
```

The `releases/` and `unreleased/` trees are populated by γ during cycle dispatch (`CDD.md §1.4` γ algorithm Phase 1). The operator does not pre-create cycle directories; γ owns that. The remaining files listed above — `CDD-VERSION`, `DISPATCH`, `CADENCE`, `OPERATORS`, `MCAs/INDEX.md`, `skills/README.md`, `iterations/cross-repo/README.md` — are activation artifacts that the operator creates during this bootstrap sequence. The `proposals/` directory starts empty.

---

## §4 Version Pin

Create the file `.cdd/CDD-VERSION` containing the cnos commit SHA that this repository follows. The version pin declares which revision of the cdd skill bundle the repo is locked to; it is the anchor for skill bundle integrity checks (§16) and the reference point for CDD upgrade decisions. The format is one item per line: the full 40-character commit SHA on the first line, followed by the associated tag on the second line if one exists.

```
6d4bb436159fa9adfc96c32e4868f2d60049bdae
3.74.0
```

If no tag exists for the pinned SHA, omit the second line. If the SHA is for a commit between tags, record the SHA only; do not interpolate a tag. The tag line, when present, is informational — the SHA is the authoritative pin. When the operator bumps the version (to pick up skill improvements), they update both lines atomically, then refresh the vendored skill bundle per §16.

The version pin must be set before any cycles run against the repo. CI checks the SHA against the pinned version per §16 and will fail if the pin is absent or malformed. A correctly formed pin contains exactly one 40-character hex SHA on line 1, with an optional non-empty tag on line 2 and no trailing content.

---

## §5 Labels

Create the minimum label set on the repository host. These labels are the vocabulary that issue triage, CI filtering, and MCA tracking depend on. Without them, cycle issues cannot be labeled consistently, MCA auto-suggest (§15) cannot fire, and the CI artifact validator cannot distinguish cycle-type issues.

**Minimum required labels:**

| Label | Color (suggested) | Purpose |
|---|---|---|
| `cdd` | `#0075ca` | Applied to every issue that runs a CDD cycle |
| `mca` | `#e4e669` | Multi-cycle architecture — issues that span ≥2 cycles |
| `P0` | `#d73a4a` | Critical priority |
| `P1` | `#e99695` | High priority |
| `P2` | `#f9d0c4` | Medium priority |
| `P3` | `#fef2c0` | Low / informational |

On GitHub, create labels via `gh label create <name> --color <hex> --description <text>` or through the repository Labels UI. The project may add additional labels beyond this minimum set — for example, tsc adds `kata` to mark issues that include kata authoring — but the six labels above must always be present. Remove or rename these labels only with an explicit migration plan; CI and triage rules reference them by exact name.

---

## §6 Branch Convention

All CDD cycle branches follow the naming pattern `cycle/{N}` where `{N}` is a monotonically increasing positive integer. This convention makes cycle branches discoverable by tooling (`git branch -r --list 'origin/cycle/*'`), prevents namespace collisions, and makes the association between a branch and its issue number legible at a glance. The number `{N}` is the issue number of the governing issue, not a sequence counter separate from the issue tracker.

γ allocates `{N}` at dispatch time by reading the governing issue number. γ creates the branch as `cycle/{N}` from the current `origin/main` HEAD. α and β check out this branch but never create or rename it. The operator (δ) creates no branches either — branch creation is a γ action per `CDD.md §1.4 γ algorithm Phase 1 step 3a`. If a fix round requires a new branch due to harness push restrictions (per `operator/SKILL.md §5.2` consequence 3), γ manages the chain naming.

Branch protection rules for `main` should be set before activation is complete. At minimum: require pull request review before merging, require status checks to pass (CI from §9), and disallow force-push. These protections enforce that every merge to `main` has passed β review and CI validation — the mechanical guarantees that CDD's coherence claims rest on.

---

## §7 Identity Convention

Every CDD role actor configures a git identity in the form `{role}@{project}.cdd.cnos` before making any commits on the cycle branch. This form encodes the protocol origin (cnos), the protocol namespace (cdd), the tenant project, and the role in a single observable string that appears in every `git log` output and commit trailer. The full specification is in `operator/SKILL.md §"Git identity for role actors"`, which is the canonical source; this section summarizes the activation-time setup procedure.

**Special case — cnos itself.** When the project running the cycle is cnos, the redundant form `{role}@cnos.cdd.cnos` is replaced by the elision `{role}@cdd.cnos`. This is the form used in all cnos commit trailers. Non-cnos projects use the full three-part form.

During activation, record the identity forms in a comment or README so that operators know what to `git config` before dispatching. The configuration command is:

```bash
# Non-cnos projects:
git config user.name "{role}"
git config user.email "{role}@{project}.cdd.cnos"

# cnos project (elision form):
git config user.name "{role}"
git config user.email "{role}@cdd.cnos"
```

This must be run at the start of each dispatch session — the identity is not globally persisted. If an operator dispatches without setting the identity, commits will carry the wrong author email and the role-identity-is-git-observable property (`CDD.md §1.4`) will be violated. The pre-review gate in `alpha/SKILL.md §2.6 row 14` requires α to verify its commit author email before signaling review-readiness; CI may additionally check commit author emails against the operator registry (§19) on subsequent cycles.

---

## §8 Dispatch Declaration

Create the file `.cdd/DISPATCH` declaring which dispatch configuration this repository uses. The two valid configurations are defined in `operator/SKILL.md §5`: §5.1 (canonical multi-session, one `claude -p` process per role) and §5.2 (single-session δ-as-γ via Agent tool, Claude Code activation). The dispatch declaration exists so that β and γ know which grading floor applies (`release/SKILL.md §3.8` configuration-floor clause) and so that future operators joining the project know how cycles are run without reading prior γ-closeouts.

The format is a single line identifying the operator skill section followed by a one-line description:

```
§5.2 — single-session δ-as-γ via Agent tool (Claude Code)
```

or

```
§5.1 — canonical multi-session (one claude -p per role)
```

Do not use richer YAML for the dispatch file unless dispatch sub-modes multiply beyond the two currently defined configurations; added complexity must earn its weight. Per-cycle overrides are declared in `releases/{version}/{N}/DISPATCH-OVERRIDE` (§18) and do not modify `.cdd/DISPATCH` — the root file records the repo's default configuration.

The dispatch configuration affects the honest-grading floor: §5.2 cycles carry a lower grading ceiling on the γ axis because γ/δ separation is structurally absent (a single parent session holds both roles). This is not a defect — it is a declared trade-off. Recording it in `.cdd/DISPATCH` makes the trade-off visible to every reviewer.

---

## §9 CI Integration

CI is the mechanical enforcement surface for CDD governance. A repo running CDD cycles without CI is relying entirely on role-actor discipline; CI converts that discipline into a verifiable gate. The minimum CI surface for a CDD-activated repo has three layers: artifact validation, test runs, and project-specific progressions.

**Layer 1 — CDD artifact validation.** On every push to any branch, CI checks that the cycle directory (if present) contains the required file set defined in §12. The canonical example for this check is `scripts/validate-release-gate.sh` (landed in cnos cycle #339), which performs mechanical pre-merge gate validation of the `.cdd/` structure and required artifact presence. Tenant repos may vendor this script or implement an equivalent; what matters is that file-presence validation fires on every push and fails loudly with a named missing artifact.

**Layer 2 — Test and spec runs.** On every push, CI runs the project's test suite and spec validation (if applicable). For docs-only repos these may be link-check or markdown-lint runs; for code repos these are unit and integration test suites. Failing tests block merge. This layer ensures that no CDD cycle can merge a change that breaks the project's own correctness surface.

**Layer 3 — Project-specific progressions.** Additional CI jobs that matter for the specific project: for tsc, kata runs (one job per kata, all must pass); for cnos, skill-bundle integrity checks (§16); for any project, linting or formatting gates that the project has adopted. These are listed in the activation checklist (§23 step 9) and become part of the repo's CI baseline.

CI jobs must be configured to run on `cycle/*` branches as well as `main`. A job that only runs on `main` is invisible during review and provides no signal to α and β during the cycle. The canonical GitHub Actions trigger is `on: push: branches: ['main', 'cycle/**']`; without the `cycle/**` glob, CI silently passes review on branches where it has never run.

**Minimum CI workflow structure:**

```yaml
name: cdd-validate
on:
  push:
    branches: ['main', 'cycle/**']
  pull_request:
    branches: ['main']

jobs:
  artifact-validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Validate CDD artifact presence
        run: bash scripts/validate-release-gate.sh

  tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run project test suite
        run: # project-specific test command
```

The workflow is a template; tenant repos substitute their actual test command and add project-specific jobs. The `artifact-validate` job is mandatory and non-negotiable; the `tests` job shape adapts to the project's language and toolchain.

---

## §10 Notification Integration

Notifications give the operator (δ) and project stakeholders a passive view of the CDD pipeline without polling the repository continuously. This section defines the transport-agnostic notification interface: the event vocabulary, the adapter contract, and the extension point where transport-specific implementations plug in.

### §10.1 Event vocabulary

Every notification adapter must handle at minimum the following four events:

| Event | Trigger | Payload minimum |
|---|---|---|
| `cycle-open` | Push to `cycle/{N}` branch (new branch) | cycle number, issue title, branch name |
| `beta-verdict` | Push of `beta-review.md` to cycle branch | cycle number, verdict (APPROVED / REQUEST CHANGES), round number |
| `cycle-rc` | Push of a fix-round artifact (α re-dispatch after RC) | cycle number, RC round number, fix-round commit SHA |
| `cycle-merge` | Merge of cycle branch to `main` | cycle number, merge commit SHA, version (if tagged) |

Adapters may implement additional events (e.g., `dispatch-start`, `ci-failure`, `sla-breach`) but must implement all four above to be considered conformant. Notification fatigue is a real cost; the four milestone events above represent the coarsest useful granularity — one notification per cycle phase transition. Per-commit or per-push notifications for intermediate α work are not recommended.

### §10.2 Adapter contract

A conformant notification adapter is any executable (script, binary, or GitHub Action workflow step) that:

1. Accepts the event name as the first argument (or as an environment variable `CDD_EVENT`)
2. Accepts event payload fields as environment variables (`CDD_CYCLE_N`, `CDD_VERDICT`, `CDD_ROUND`, `CDD_SHA`, `CDD_TITLE`)
3. Returns exit code 0 on successful delivery, non-zero on failure
4. Logs the delivery attempt (event name, destination, timestamp) to stdout for CI log capture
5. Does not require modification to `cdd/activation/SKILL.md` to add a new transport

The contract is intentionally minimal so that a Telegram adapter, a Slack adapter, a Discord adapter, or an email relay can each be implemented without changing the interface. The reference adapter for Telegram is defined in `cdd/activation/templates/telegram-notifier/` (Cycle B deliverable); it is the canonical example of the contract, not the contract itself. A hypothetical Slack adapter would implement the same five properties above using Slack's Incoming Webhooks API instead of Telegram's Bot API — no changes to the contract or to this skill would be required.

---

## §11 Secrets

Notification adapters and other CDD integrations require credentials that must never appear in repository files or commit history. This section prescribes exact secret names, their storage location, rotation cadence, and the hard prohibition on committing tokens.

**Secrets table:**

| Secret name | Purpose | Storage | Rotation cadence |
|---|---|---|---|
| `CDD_TELEGRAM_BOT_TOKEN` | Telegram Bot API token for notification adapter | GitHub repository secrets | Quarterly default; every-cycle for high-sensitivity repos |
| `CDD_TELEGRAM_CHAT_ID` | Telegram chat or channel ID that receives notifications | GitHub repository secrets | On channel change |

All CDD-managed secrets use the `CDD_` prefix. This namespace prevents collision with project-level secrets and makes CDD-managed credentials identifiable in the secrets UI. New secrets introduced by future adapters (Slack, Discord) follow the same naming pattern: `CDD_{ADAPTER}_{CREDENTIAL_TYPE}`.

**Storage rule:** All secrets live in GitHub repository secrets (Settings → Secrets and variables → Actions). They are never stored in `.env` files, YAML workflow files, shell scripts, or any other version-controlled artifact. They are never logged, echoed, or printed in CI output. GitHub's secret masking provides a second layer of protection, but the first layer is never writing the value anywhere that masking does not cover.

**Prohibition on example tokens:** No example tokens — even fabricated ones — appear in this file or in any activation template. Fabricated tokens are easy to mistake for real ones during copy-paste onboarding, and they create a false sense of "I just need to swap this value out" that leads to leakage.

**Rotation procedure:** To rotate `CDD_TELEGRAM_BOT_TOKEN`, generate a new token via BotFather, update the GitHub repository secret, and verify that the next CI run posts a test notification. The old token is invalidated by BotFather; no code changes are required. Record the rotation in `.cdd/OPERATORS` or a cycle close-out for audit traceability.

---

## §12 Cycle-README Template

Every cycle directory — `releases/{version}/{N}/` for versioned releases or `releases/docs/{ISO-date}/{N}/` for rolling-docs releases — must contain the following five files. These files constitute the cycle close-out record and are what CI artifact validation (§9 Layer 1) checks for:

| File | Author | Purpose |
|---|---|---|
| `self-coherence.md` | α | Gap statement, AC evidence, CDD Trace, review-readiness signal |
| `alpha-closeout.md` | α (re-dispatched) | Post-merge cycle findings, friction log, engineering observations |
| `beta-review.md` | β | Review verdict (APPROVED / REQUEST CHANGES), per-finding record, merge evidence |
| `gamma-closeout.md` | γ | TSC Grades, PRA, dispatch record, deferred-output list |
| `cdd-iteration.md` | ε (often δ=ε) | Protocol-iteration findings for this cycle; see §22 |

γ creates the cycle directory and stubs during dispatch. α, β, and γ fill files in the order shown. The cycle directory is initially created under `unreleased/{N}/`; at release time, γ moves it to the versioned path per `release/SKILL.md §2.5a`. CI artifact validation checks that all five files exist and are non-empty before allowing merge to `main`.

**Note on alpha-closeout.md:** Per `alpha/SKILL.md §2.8`, α writes `alpha-closeout.md` after β merge, not before. At review-readiness time, `alpha-closeout.md` either does not exist (standard path) or exists as `[provisional — pending β outcome]`. CI should not require `alpha-closeout.md` before merge; it should check for it in the post-merge gate only.

---

## §13 Cross-Repo Trace Bundle Init

When a cycle's work produces deliverables in another repository — or when a cycle coordinates with a cycle in a sibling repo — a cross-repo trace bundle tracks the coordination. Pre-creating the bundle root at activation time makes the convention discoverable before any cross-repo work occurs; the README documents the naming pattern so that future cycles can create bundles without re-deriving the structure.

Create `.cdd/iterations/cross-repo/README.md` with the following content during activation:

```markdown
# Cross-Repo Trace Bundles

Each cross-repo coordination unit lives at:

  .cdd/iterations/cross-repo/{target}/{slug}/

Where:
- {target}   — the target repository name (e.g., tsc, acme-api)
- {slug}     — a short, lowercase, hyphen-separated descriptor of the coordination
               (e.g., supercycle-v0.7.0, schema-migration-v2)

## When to create a bundle

Create a bundle when:
- This repo's cycle produces deliverables that land in {target}
- A cycle in this repo coordinates sequentially or in parallel with a cycle in {target}
- A design decision in this repo has binding downstream effect on {target}

## Bundle contents

Each bundle directory contains:
- README.md  — bundle purpose, originating and target cycles, current status
- STATUS     — one of: open | converging | closed (updated as cycles complete)

A bundle is closed when both the originating-repo cycle and the target-repo cycle
have shipped close-outs and the STATUS file reads 'closed'.
```

The nested `{target}/{slug}/` form (OQ #10) matches the existing cnos-tsc supercycle layout and keeps bundles grouped by target repo, making `ls .cdd/iterations/cross-repo/` a readable index of all active cross-repo relationships.

---

## §14 Honest-Claim Manifest Convention

Every cycle produces a `claims.md` at `releases/{version}/{N}/claims.md` (or `releases/docs/{ISO-date}/{N}/claims.md` for rolling-docs releases). The manifest enumerates every honest claim made in the cycle — measurements, source-of-truth alignments, and wiring assertions — with a verification command or artifact reference for each. No empty manifests are allowed; every cycle must produce at least one claim statement even if the claim is "no behavior change — this cycle adds only documentation."

**Minimum claim categories (per `review/SKILL.md §3.13`):**

| Category | What it covers | Example claim |
|---|---|---|
| Reproducibility | Measurements quoted in cycle artifacts are reproducible from committed artifacts | `wc -l activation/SKILL.md → 623` (run from repo root at HEAD) |
| Source-of-truth alignment | Terms used in docs trace to their canonical definitions | `§7 identity form matches operator/SKILL.md §"Git identity" table` |
| Wiring | Integration points (CI hooks, secrets, adapters) connect as described | `CDD_TELEGRAM_BOT_TOKEN secret present → notifier posts on cycle-open event` |

Cycles may add categories beyond the three above with rationale (e.g., a cycle introducing a new protocol surface might add "contract completeness" claims). The floor is not exhaustive; it is the minimum that prevents the most common class of honest-claim failures surfaced in review.

The manifest is authored by α as part of the cycle close-out artifact set and reviewed by β alongside `self-coherence.md`. β's `review/SKILL.md §3.13` honest-claim check uses `claims.md` as its primary audit surface for the cycle. Leaving it empty is a D-level finding.

**Suggested claims.md structure:**

```markdown
# Honest-Claim Manifest — Cycle #N

## Reproducibility claims
- Claim: `wc -l activation/SKILL.md → 623`
  Verification: `wc -l src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` at HEAD

## Source-of-truth alignment claims
- Claim: §7 identity form matches operator/SKILL.md §"Git identity" table
  Verification: diff §7 identity table against operator/SKILL.md identity table — zero delta

## Wiring claims
- Claim: cross-references in CDD.md, operator/SKILL.md, post-release/SKILL.md point to activation/SKILL.md
  Verification: `rg 'activation/SKILL.md' src/packages/cnos.cdd/skills/cdd/` → ≥3 hits
```

---

## §15 MCA Registry

Multi-cycle architectures (MCAs) are design efforts that span more than one cycle. Without a central index, MCAs are tracked only through issue labels and per-cycle close-out references — they are discoverable but not indexed. The MCA registry makes the set of active MCAs visible as a single artifact.

Create `.cdd/MCAs/INDEX.md` with a table of in-flight MCAs:

```markdown
# MCA Registry

| slug | originating-cycle | target-close-cycle | owner | status |
|---|---|---|---|---|
| (none) | — | — | — | — |
```

For each active MCA, also create `.cdd/MCAs/{slug}/README.md` containing: the slug, the originating cycle number, the target cycles (all N where the MCA expects to deliver), the owner (δ handle), the current state description, and close criteria.

**Close criteria for an MCA:** An MCA moves from `converging` to `closed` when (a) all named target cycles have shipped close-outs and (b) the close-out cycle's `cdd-iteration.md` confirms that the axis the MCA tracked has collapsed — i.e., the recurring finding that prompted the MCA no longer appears in the most recent cycle.

**Auto-suggest trigger:** When ≥3 same-axis findings appear in 5 or more consecutive cycle `cdd-iteration.md` files, CI emits an alert suggesting an MCA for that axis. The alert is informational — the operator (δ) decides whether to formalize it. CI does not create MCAs automatically; auto-creation would produce MCAs without operator judgment about whether the pattern is real or coincidental.

---

## §16 Skill Bundle Pull/Sync

Tenant repos need a local copy of the cdd skill bundle so that cycles can load `operator/SKILL.md`, `alpha/SKILL.md`, and the rest without requiring network access to cnos. The canonical approach is a vendored copy under `.cdd/skills/`, frozen at the SHA pinned in `.cdd/CDD-VERSION`.

**Vendored copy procedure:**

1. At activation time, clone or download the cnos repo at the SHA in `.cdd/CDD-VERSION`
2. Copy `src/packages/cnos.cdd/skills/cdd/` into `.cdd/skills/cdd/`
3. Commit the vendored copy as part of the activation commit

**Refresh cadence:** Refresh only when bumping `.cdd/CDD-VERSION`. The operator explicitly bumps the version with a rationale commit message, then re-runs the copy procedure. Automatic refresh on every cycle would silently change the skill surface mid-project; the explicit bump-and-refresh pattern makes skill version changes deliberate and audit-traceable.

**Integrity check in CI:** One CI step compares the vendored skill file SHAs against the source tree at the cnos commit SHA declared in `.cdd/CDD-VERSION`. If the vendored copy has diverged (manual edit, incomplete refresh, partial merge), CI fails with a named diff. This check is the mechanical guarantee that the vendored skills match the declared version.

**Special case — cnos itself:** cnos is the origin repository for the skill bundle. A vendored copy of itself is recursive and unnecessary. For cnos, `.cdd/skills/README.md` declares: "cnos is the origin repository for this skill bundle. The canonical skill location is `src/packages/cnos.cdd/skills/cdd/`. Tenant repos vendor from the cnos commit SHA pinned in `.cdd/CDD-VERSION`."

---

## §17 Pre-Commit / Pre-Push Hooks

Optional git hooks provide fast local feedback before CI sees a push. They are opt-in — developer machines should not have hooks installed by a clone or build step; many developers have strong opinions about hook installation and CI is the authoritative gate anyway. The `.cdd/hooks/` directory contains the hook sources; developers who want local enforcement install them via symlink.

**Installation (opt-in):**

```bash
ln -sf ../../.cdd/hooks/validate-structure .git/hooks/pre-commit
ln -sf ../../.cdd/hooks/validate-cycle-readme .git/hooks/pre-push
```

**Hook scope — structure only locally.** Pre-commit and pre-push hooks check `.cdd/` structure: are the required directories present? Is `.cdd/CDD-VERSION` a 40-character SHA? Is the `DISPATCH` file non-empty? Content checks (claim manifest non-empty, README template fields filled, honest-claim categories present) run in CI, not locally. Keeping local hooks to structure-only means they run in milliseconds and do not require the full CI toolchain to be installed.

**Failure mode:** Block on malformed `.cdd/` structure (missing required files or directories); warn on missing optional fields (e.g., no `MCAs/INDEX.md` entry when the issue has `mca` label). A block stops the commit or push and prints the exact failing check; a warning prints and proceeds. The distinction matters: blocks protect the CI surface from a known class of failures; warnings surface best-practice gaps without imposing hard gates on developer flow.

---

## §18 Per-Cycle Dispatch Override

The default dispatch configuration is declared in `.cdd/DISPATCH` (§8). Individual cycles may deviate with an explicit per-cycle override. Overrides exist for cases where the default configuration is inappropriate for a specific cycle's scale or complexity — for example, escalating from §5.2 to §5.1 when a cycle has ≥7 ACs (see `operator/SKILL.md §5.3` escalation criteria).

**Override file location:** `releases/{version}/{N}/DISPATCH-OVERRIDE` (or `releases/docs/{ISO-date}/{N}/DISPATCH-OVERRIDE` for rolling-docs).

**Override format:**

```
§5.1 — canonical multi-session
Reason: ≥7 ACs in this cycle; escalation criteria per operator/SKILL.md §5.3 apply.
```

The first line is the section reference of the configuration actually used. The second line is a mandatory one-line reason. An override without a reason field is treated as a malformed override and flagged by CI artifact validation.

**Visibility in close-out:** γ-closeout TSC Grades records the dispatch configuration actually used for the cycle, noting whether it was the repo default or an override. This makes the grading floor visible in the close-out record without requiring the reader to cross-reference `.cdd/DISPATCH` and the override file separately.

---

## §19 Operator Handle Registry

The operator handle registry lists every authorized δ actor for the repository. Its purpose is to make the set of authorized operators observable, provide an audit surface for commit identity checks, and give CI a reference to validate commit author emails against.

Create `.cdd/OPERATORS` as a plain-text table:

```
# CDD Operator Registry
# Columns: handle | identity-email | role-scopes-permitted | added-cycle
#
alpha    | alpha@cdd.cnos    | alpha               | 0
beta     | beta@cdd.cnos     | beta                | 0
gamma    | gamma@cdd.cnos    | gamma               | 0
```

**Columns:**

| Column | Meaning |
|---|---|
| `handle` | Short name used to refer to the operator in close-outs and overrides |
| `identity-email` | The canonical git author email for this operator's CDD commits (`{role}@{project}.cdd.cnos` per §7) |
| `role-scopes-permitted` | Comma-separated list of CDD roles this operator may act as (alpha, beta, gamma, delta, or `*`) |
| `added-cycle` | Cycle number at which this operator was added to the registry (0 for activation-time entries) |

**CI enforcement:** For the first cycle after activation, CI warns on unregistered operator commits but does not block merge. This grace period lets the activation-time operator complete onboarding without being immediately blocked by the check they just created. From the second cycle onward, an unregistered commit author email causes CI to fail the artifact validation job.

**Handle revocation:** When an operator is removed from active participation, populate their row's `removed-cycle` column (add the column if not present) with the cycle number at which they were removed. Do not delete the row. Historical commits by that operator remain valid; the registry is forward-only. A removed operator's commits on cycles predating the removal are grandfathered; new commits after the removal cycle would fail the CI check.

---

## §20 Close-Out SLA

CDD cycles must produce their close-out artifacts within a defined window after merge to `main`. Without an SLA, close-out debt accumulates silently — merge happens, cycle-dir files sit in `unreleased/`, and the next cycle starts before the prior one is fully closed. The SLA makes the expectation explicit and gives CI a surface to alert on.

**Default SLA:** 24 hours from merge to `main` to all five cycle close-out files present in the merged commit on `main` (per §12). This is the default for all repos that have not overridden it.

**Override:** The SLA may be overridden per-repo in `.cdd/CADENCE`. Add a line `close-out-sla: 48h` (or the desired duration) to override the default. The override requires a rationale comment in `.cdd/CADENCE`.

**Measurement:** "Close-out artifact visibility" means the close-out files are present in a commit merged to `main` — not merely pushed to the cycle branch. The SLA clock starts at the merge commit timestamp and stops when the last required close-out file appears in a `main` commit.

**Enforcement posture by cycle:** In Cycle B, the notification adapter emits a `sla-breach` notification when the SLA window passes without close-out file detection. No merge-blocking occurs in Cycle B — the alert is informational. Block-on-breach (failing the CI artifact gate for the next cycle until the prior cycle's close-out files are present) is deferred to a follow-on cycle after operator experience with alert-only confirms that blocking is warranted.

---

## §21 Release Cadence Declaration

Create `.cdd/CADENCE` declaring the repository's release cadence. The cadence affects directory layout (§3), tagging discipline, and the SLA override mechanism (§20).

**Valid values:**

| Value | Meaning | Directory layout |
|---|---|---|
| `versioned` | Cycles map to semantic version releases (e.g., `v0.7.0`) | `releases/{semver}/{N}/` |
| `rolling-docs` | Cycles are date-anchored documentation releases with no version tag | `releases/docs/{ISO-date}/{N}/` |
| `mixed` | Some cycles are versioned engine releases; others are rolling-docs cycles | Both layouts in use |

**File format:**

```
versioned
```

or

```
rolling-docs
```

or

```
mixed
# versioned: engine releases (e.g., v0.7.0)
# rolling-docs: cdd skill/docs cycles (date-anchored)
```

**Tagging discipline per cadence:**

- `versioned`: every cycle that ships a release receives a semantic version tag at merge (e.g., `git tag 3.74.0`). Tagging is a δ action executed via `scripts/release.sh` per `operator/SKILL.md §3.4`; manual `git tag` is not allowed.
- `rolling-docs`: cycles do not receive version tags. The date anchor in the directory name is the release identifier. Tagging a rolling-docs cycle would create a misleading "release" artifact with no version semantics.
- `mixed`: tag versioned cycles; do not tag rolling-docs cycles. γ-closeout records the cadence type for each cycle in the TSC Grades section so the close-out record is self-explanatory.

---

## §22 cdd-iteration Cadence

Every cycle produces a `cdd-iteration.md` in its cycle directory. This file is the work product of ε (the protocol-iteration role, per `ROLES.md §1` row 5 and `epsilon/SKILL.md §1`). In most repos δ=ε — the operator doubles as the protocol-iteration actor — and `cdd-iteration.md` is written by the same person who writes `gamma-closeout.md`. Whether or not there is a separate ε actor, the file must be produced.

**Why every cycle, not only cycles with findings:** An empty findings list is itself a signal — it says the cycle ran cleanly with no protocol gaps surfaced. Conditioning `cdd-iteration.md` on the existence of findings produces a survivorship bias: the absence of the file is ambiguous (clean cycle, or close-out skipped?). A file present with content "no findings this cycle" is unambiguous. This prescription extends the `epsilon/SKILL.md §1` "only when findings exist" language; the activation skill is the authoritative prescription here.

**Findings severity scale:** Reuse the severity scale from `review/SKILL.md §3.2` for consistency: D (blocker-class protocol incoherence), C (significant incoherence, non-blocking), B (improvement opportunity), A (polish). Add `info` for non-actionable observations that are worth recording but require no follow-up action. Using the same scale as β review findings keeps the vocabulary consistent across the cycle close-out artifact set.

**Auto-spawn MCA trigger:** When ≥3 same-axis findings appear in 5 or more consecutive cycle `cdd-iteration.md` files, CI emits an operator alert. The operator (δ) decides whether to formalize the axis into an MCA entry in `.cdd/MCAs/INDEX.md`. N=5 is the window for the sliding count; the window slides forward one cycle at a time. Auto-creation is explicitly not performed — CI alerts, operator approves materialization. This prevents MCA proliferation from CI sensitivity while still surfacing persistent patterns.

**Reference cdd-iteration.md per-finding shape** follows `post-release/SKILL.md §5.6b` (authoring procedure). Each finding has: axis, severity, description, evidence (commit or artifact reference), and disposition (MCA available / protocol MCI filed / agent MCI updated / dropped-one-off with rationale).

---

## §23 Activation Checklist

Run these steps in order on a new repo. Each step references the governing section and states its verification. No step depends on a later step; the list is a valid linearization of the activation sequence.

1. **Create `.cdd/` directory tree** (§3). Verify: `ls .cdd/` shows `CDD-VERSION`, `DISPATCH`, `CADENCE`, `OPERATORS`, `MCAs/`, `skills/`, `releases/`, `unreleased/`, `iterations/`, `proposals/`.

2. **Create `.cdd/CDD-VERSION`** (§4). Populate with the cnos HEAD SHA and tag (one per line). Verify: `cat .cdd/CDD-VERSION` shows a 40-char SHA on line 1 and a version tag on line 2 (if a tag exists at that SHA).

3. **Create `.cdd/DISPATCH`** (§8). Choose §5.1 or §5.2 based on your dispatch environment; record one line with the section reference and a brief description. Verify: `cat .cdd/DISPATCH` is non-empty and references `§5.1` or `§5.2`.

4. **Create `.cdd/CADENCE`** (§21). Set `versioned`, `rolling-docs`, or `mixed`. Verify: `cat .cdd/CADENCE` matches one of the three values.

5. **Create `.cdd/OPERATORS`** (§19). Populate with the repo's initial δ handles; use role-actor emails in `{role}@{project}.cdd.cnos` form. Verify: `cat .cdd/OPERATORS` shows at least one row for each of alpha, beta, gamma with correct identity-email form.

6. **Create `.cdd/MCAs/INDEX.md`** (§15). Populate the table (empty `(none)` row if no active MCAs). Verify: `cat .cdd/MCAs/INDEX.md` contains the required table header columns.

7. **Create `.cdd/skills/README.md`** (§16). For non-cnos repos: vendor the skill bundle from the SHA in `CDD-VERSION`. For cnos: write the origin-repo declaration. Verify: `cat .cdd/skills/README.md` is non-empty and references the bundle location.

8. **Create `.cdd/iterations/cross-repo/README.md`** (§13). Verify: file exists and contains the `{target}/{slug}/` naming convention text.

9. **Create GitHub issue labels** (§5). Run `gh label list` and confirm `cdd`, `mca`, `P0`, `P1`, `P2`, `P3` are all present. Verify: `gh label list | grep -E '^(cdd|mca|P[0-3])'` returns 6 lines.

10. **Configure git identity for role actors** (§7). Document the identity forms in an internal README or operator guide. Verify: `git config user.email` for a test commit shows `{role}@{project}.cdd.cnos` form.

11. **Configure CI artifact validation** (§9 Layer 1). Add or reference `scripts/validate-release-gate.sh` (or equivalent) to the CI pipeline. Verify: a test push with a cycle directory missing one required file causes CI to fail.

12. **Configure CI test and spec runs** (§9 Layer 2). Verify: `gh run list` shows a recent run with test jobs for the project's test suite.

13. **Configure project-specific CI progressions** (§9 Layer 3). Add kata runs, lint jobs, or other project-specific gates. Verify: `gh workflow list` shows the expected project-specific jobs.

14. **Create GitHub repository secrets for notification adapter** (§11). Add `CDD_TELEGRAM_BOT_TOKEN` and `CDD_TELEGRAM_CHAT_ID` (or the secrets for the chosen adapter). Verify: `gh secret list` shows both secret names (values are masked).

15. **Configure branch protection on `main`** (§6). Require PR review, require CI status checks, disallow force-push. Verify: attempting a direct push to `main` is rejected with a branch protection error.

16. **Install optional pre-commit/pre-push hooks if desired** (§17). Run `ln -sf ../../.cdd/hooks/validate-structure .git/hooks/pre-commit` if local hook enforcement is wanted. Verify: `ls -la .git/hooks/pre-commit` shows the symlink (if installed), or confirm CI-only enforcement is acceptable.

17. **Configure close-out SLA override if needed** (§20). If the default 24h SLA is unsuitable, add `close-out-sla: {duration}` to `.cdd/CADENCE`. Verify: `cat .cdd/CADENCE` reflects the override, or confirm 24h default is accepted.

18. **Populate `.cdd/iterations/INDEX.md`** (§3) — create if not already present during scaffold. This file indexes cdd-iteration findings across cycles. Verify: `cat .cdd/iterations/INDEX.md` is present (may be empty initially).

19. **Commit all activation artifacts to `main`** in a single activation commit. Verify: `git log --oneline -1` shows the activation commit; `git diff main~1..main --stat` lists all seven activation marker files.

20. **Run the §24 verification command** against the activated repo. Verify: all checks pass with no `MISSING` or `MALFORMED` output.

---

## §24 Verification

Run the following command from the repository root to confirm activation succeeded. The command checks for the presence and non-emptiness of every activation marker file created in §23.

```bash
for f in \
  .cdd/CDD-VERSION \
  .cdd/DISPATCH \
  .cdd/CADENCE \
  .cdd/OPERATORS \
  .cdd/MCAs/INDEX.md \
  .cdd/skills/README.md \
  .cdd/iterations/cross-repo/README.md \
  .cdd/iterations/INDEX.md
do
  if [ ! -f "$f" ]; then
    echo "MISSING: $f"
  elif [ ! -s "$f" ]; then
    echo "EMPTY: $f"
  else
    echo "OK: $f"
  fi
done
sha=$(head -1 .cdd/CDD-VERSION 2>/dev/null)
if echo "$sha" | grep -qE '^[0-9a-f]{40}$'; then
  echo "OK: .cdd/CDD-VERSION SHA format valid"
else
  echo "MALFORMED: .cdd/CDD-VERSION — expected 40-char hex on line 1, got: $sha"
fi
```

A fully activated repository produces only `OK:` lines. Any `MISSING:`, `EMPTY:`, or `MALFORMED:` line indicates an incomplete activation step. Fix the named file and re-run until all lines are `OK:`.

The verification command is idempotent — running it on an already-activated repo is safe. It is also the canonical oracle that CI uses for §9 Layer 1 artifact presence checks; the CI job wraps this script or an equivalent and fails the run if any non-OK line appears.
