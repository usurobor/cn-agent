# Design: Package Artifact Distribution and Command Content Class

**Issue:** #167
**Version:** 0.5.0
**Mode:** MCA
**Active Skills:** design, cap, writing
**Engineering Level:** L7

## Problem

Two gaps that share a root cause:

**1. Released packages are not yet the distribution unit.** `cn deps restore` fetches a git repository commit and copies a subtree out of it. That couples package installation to git protocol capabilities. When the protocol fails (#155: Claude Code sandbox, shallow-fetch-by-SHA rejected), there is no simple recovery. The lockfile carries `source`, `rev`, `subdir` — transport metadata consumers don't need. `cn_deps.ml` has 27 references to these fields.

**2. Commands are not yet a package content class.** The CLI command surface is a closed built-in variant (`commit`, `push`, `save`, `daily`, `weekly`, `release` compiled into core). Adding operator-facing commands requires core edits (#162). The package content model already has explicit classes — doctrine, mindsets, skills, extensions, templates — but commands are not one of them.

**Evidence:**
- #155: `cn deps restore` failed in Claude Code sandbox (git 2.43.0, shallow-fetch-by-SHA unsupported). Three patches shipped to work around git protocol limitations.
- `cn_deps.ml`: 27 references to `source_kind`/`rev`/`subdir` — complexity that exists only because distribution is coupled to repository structure.
- #162: modular CLI commands requested, no current path without core modification.

**What fails if skipped:** Every new environment that restricts git protocol requires another fetch workaround. CLI command extension requires either core modification or a second plugin framework.

## Constraints

- **PACKAGE-SYSTEM.md** governs package schema, content classes, and restore semantics — must be updated to add `commands` as a content class and document artifact-first distribution
- **RUNTIME-EXTENSIONS.md** governs runtime capability providers — stays separate, not touched
- **cn.package.json** schema is consumed by `cn deps list`, `cn deps restore`, `cn doctor` — `commands` addition must follow the existing explicit content-class pattern
- **Release workflow** (`.github/workflows/release.yml`) already builds binaries and checksums — must be extended, not replaced

## Challenged Assumption

Two assumptions replaced:

1. **Packages are distributed by fetching repository objects and extracting subtrees.** Replaced by: released packages are distributed as versioned artifacts over HTTP. The repository is the source of truth for development; the artifact is the unit of distribution.

2. **CLI commands are compiled built-ins or nothing.** Replaced by: commands are a package content class, discoverable at runtime, following the same explicit-class pattern as doctrine/mindsets/skills/extensions/templates.

## Design Decisions

Three explicit design decisions, kept separate:

### Decision 1: Artifact-first distribution for released packages

Released first-party packages restore over HTTPS from versioned tarballs. A package index resolves name+version to URL+SHA-256. Git becomes a development/exception path, not the consumer default.

### Decision 2: Commands as a first-class package content class

This issue adds `commands` as a package content class alongside doctrine, mindsets, skills, extensions, and templates. Commands are declared explicitly in `cn.package.json` under a `commands` key, following the existing content-class pattern. This is not smuggling commands into some other class — the package-system model prefers explicit classes over a generic tree, and the docs explain why that explicitness is preferred at the current scale.

### Decision 3: Runtime capability extensions remain separate

Runtime extensions = typed capability providers for the agent/runtime plane (network, browser, device, API). Package commands = operator/control-plane commands. Skills/templates/doctrine = cognitive substrate. Three-way split:

| Surface | What | Example |
|---------|------|---------|
| Runtime extensions | Typed capability providers for agent execution | `cnos.net.http` |
| Package commands | Operator-facing CLI commands | `cn daily`, `cn save` |
| Cognitive substrate | Skills, doctrine, mindsets, templates | `coherent`, `ENGINEERING.md` |

## Impact Graph

### Downstream consumers
- `cn deps restore` — switch to HTTP tarball fetch
- `cn deps vendor` — understand new lockfile format
- `cn deps list` — reads installed packages, minimal change
- `cn doctor` — validate package integrity + command integrity
- `cn help` — discover and list package commands + repo-local commands
- `cn <name>` dispatch — resolve by precedence (built-in → repo-local → package)
- `cn_runtime_contract.ml` — emit skill activation table in cognition layer at wake
- Lockfile schema — simplifies (drops git transport fields)
- `cn.package.json` schema — extends (`commands` content class)
- `PACKAGE-SYSTEM.md` — must document `commands` as content class and artifact-first distribution

### Upstream producers
- Release workflow — build and publish package tarballs + index
- `scripts/stamp-versions.sh` — stamp package artifact versions
- `scripts/check-version-consistency.sh` — validate package artifacts
- Package index — new artifact, generated at release time
- `scripts/build-packages.sh` — new script, builds tarballs from `packages/` dirs

### Copies and authority
- Package index (authority for name+version → URL+hash resolution)
- Lockfile (authority for what a hub has installed — name+version+hash)
- `cn.package.json` inside tarball (authority for package metadata including commands)
- `cn.package.json` in repo source (development source, generates the above)

## Proposal

### 1. Package artifact

A package artifact is a tarball:

```
<name>-<version>.tar.gz
```

#### 1.1 Package directory layout

A package follows a well-known directory structure. Each content class has a conventional directory. The runtime knows how to discover and load each class from these directories — the manifest declares what exists, the runtime owns loading semantics. (NuGet model: dumb package, smart framework.)

```
<name>@<version>/
├── cn.package.json          # manifest (required)
├── doctrine/                # doctrinal specs
│   ├── FOUNDATIONS.md
│   └── ...
├── mindsets/                # agent mindsets
│   ├── ENGINEERING.md
│   └── ...
├── skills/                  # agent skills
│   ├── cdd/
│   │   ├── SKILL.md         # entry point (convention)
│   │   ├── design/
│   │   │   └── SKILL.md
│   │   ├── review/
│   │   │   └── SKILL.md
│   │   └── ...
│   ├── agent/
│   │   ├── cap/
│   │   │   └── SKILL.md
│   │   └── ...
│   └── ...
├── extensions/              # runtime capability providers
│   └── cnos.net.http/
│       └── ...
├── templates/               # hub initialization templates
│   ├── SOUL.md
│   └── USER.md
└── commands/                # CLI command executables
    ├── cn-daily
    └── cn-weekly
```

#### 1.2 Loading conventions

The package is the delivery unit. The runtime and agent own loading:

| Content class | Convention | Loaded by | When |
|--------------|-----------|-----------|------|
| **doctrine** | `doctrine/*.md` | Runtime at wake | Always |
| **mindsets** | `mindsets/*.md` | Runtime at wake | Always |
| **skills** | `skills/<name>/SKILL.md` | Agent skill-loading gate | By trigger (keyword, work shape, explicit) |
| **extensions** | `extensions/<name>/` | Runtime extension host | At capability registration |
| **templates** | `templates/*.md` | `cn setup` / `cn init` | At hub initialization |
| **commands** | `commands/cn-<name>` | `cn` CLI dispatch | By operator invocation |

**Key principle:** loading logic does not live in the manifest. The manifest declares what content classes exist and what items belong to each class. The runtime and agent know the conventions for each class. This keeps packages simple and loading behavior consistent across all packages.

#### 1.3 Skill activation and encapsulation

**Activation:** Skills declare activation keywords in their SKILL.md frontmatter:

```yaml
---
name: cdd
description: Apply Coherence-Driven Development...
triggers: [review, PR, release, issue, design, plan, assess, post-release]
---
```

When a package is installed, the runtime walks `<pkg>/skills/` for every SKILL.md on disk, parses each frontmatter, and builds a public activation table from the skills whose `visibility` is not `internal`. Agent encounters a trigger keyword → runtime loads the matching skill. No per-agent configuration needed and no manifest-declared inventory.

**Encapsulation:** Top-level orchestrator skills stay public; sub-skills are marked internal in their own frontmatter. The manifest does not enumerate skills.

`cdd/design`, `cdd/review`, `cdd/release`, etc. are internal to the CDD skill. They live in the `skills/cdd/` directory and declare `visibility: internal` in their frontmatter, which excludes them from the public activation index. The orchestrator skill (`cdd/SKILL.md`) owns all delegation — it decides which sub-skill to load at which pipeline step.

Sub-skills declare their parent and visibility:

```yaml
---
name: review
description: CLP review protocol...
parent: cdd
visibility: internal
---
```

**Lifecycle:**

1. **Install:** `cn deps restore` extracts package → runtime walks `<pkg>/skills/` for SKILL.md files → parses each frontmatter → trigger keywords from public (non-internal) skills are added to the agent's activation table, referencing the orchestrator skill
2. **Activate:** Agent encounters trigger keyword → runtime looks up activation table → loads the orchestrator skill → orchestrator delegates to internal sub-skills as needed
3. **Uninstall:** Package directory removed → runtime rebuilds activation table by rescanning remaining packages → trigger keywords from the removed package disappear

The activation table is derived, not configured. It is always the union of trigger keywords from every public skill discovered on disk across installed packages. No manual maintenance and no manifest-declared skill inventory.

**Where the activation table lives:** in the runtime contract, emitted at every wake. The runtime contract (`cn_runtime_contract.ml`) already describes the agent's identity, cognition, body, and medium. The activation table belongs in **cognition** — it tells the agent what skills it has and when to use them.

```
## Cognition

Packages:
  cnos.core 3.34.0 (doctrine: 12, mindsets: 9, skills: 4, commands: 2)
  cnos.eng 3.34.0 (skills: 15)

Skills:
  cdd       [triggers: review, PR, release, issue, design, plan, assess, post-release]
  coherent  [triggers: coherence, check, verify]
  cap       [triggers: action, MCA, MCI, change]
  agent-ops [triggers: ops, status, doctor, hub]
```

The agent reads this at wake and knows: what skills are available, what keywords activate them. No BOOTSTRAP.md table, no per-agent configuration. The runtime contract is the agent's complete self-knowledge — identity, cognition (including skill activation), body (capabilities), and medium (environment). Everything derived from installed packages.

**Result:**
- Agent sees one skill per concern (e.g., "cdd"), not N sub-skills
- Activation table has one entry per exposed skill
- Sub-skills are private implementation — loaded by the orchestrator, not by the runtime
- Install/uninstall a package → activation keywords appear/disappear automatically

#### 1.4 Worked example: CDD skill package

End-to-end flow from package install to skill execution.

**Package structure on disk** (`cnos.cdd/`):

```
skills/
└── cdd/
    ├── SKILL.md              ← orchestrator (public by default)
    ├── design/
    │   └── SKILL.md          ← sub-skill (visibility: internal, parent: cdd)
    ├── review/
    │   └── SKILL.md
    ├── release/
    │   └── SKILL.md
    ├── post-release/
    │   └── SKILL.md
    ├── plan/
    │   └── SKILL.md
    └── issue/
        └── SKILL.md
```

**Manifest** (`cn.package.json`):

```json
{
  "schema": "cn.package.v1",
  "name": "cnos.cdd",
  "version": "1.0.0",
  "kind": "package"
}
```

No `skills` field — skills are discovered by walking `skills/` on disk. Sub-skills are hidden from the public activation index by declaring `visibility: internal` in their own frontmatter.

**Orchestrator frontmatter** (`cdd/SKILL.md`):

```yaml
---
name: cdd
description: Coherence-Driven Development — full lifecycle from observed gap to closed cycle.
triggers: [review, PR, release, issue, design, plan, assess, post-release, ship, tag, approve]
---
```

**Sub-skill frontmatter** (`cdd/review/SKILL.md`):

```yaml
---
name: review
description: CLP review protocol for CDD step 8.
parent: cdd
visibility: internal
---
```

**Loading sequence:**

1. `cn deps restore` installs `cnos.cdd` → extracts to `.cn/vendor/packages/cnos.cdd/`
2. Runtime walks `.cn/vendor/packages/cnos.cdd/skills/` → parses each SKILL.md → filters to `visibility != internal` → builds activation table: `{review, PR, release, ...} → cdd`
3. Agent encounters "review this PR"
4. Runtime matches "review" → loads `skills/cdd/SKILL.md`
5. CDD SKILL.md §5 delegation table: review = pipeline step 8 → load `cdd/review/SKILL.md`
6. Agent loads the sub-skill, executes review within CDD's pipeline context

The agent never loads `cdd/review` directly. It always enters through CDD. CDD provides the pipeline context (what step, what came before, what comes next). The sub-skill provides step-level execution detail.

The package artifact is the normal distribution unit for released packages.

### 2. Package index

A small JSON file resolving name+version to URL+SHA-256:

```json
{
  "schema": "cn.package-index.v1",
  "packages": {
    "cnos.core": {
      "3.34.0": {
        "url": "https://github.com/usurobor/cnos/releases/download/3.34.0/cnos.core-3.34.0.tar.gz",
        "sha256": "..."
      }
    },
    "cnos.eng": {
      "3.34.0": {
        "url": "https://github.com/usurobor/cnos/releases/download/3.34.0/cnos.eng-3.34.0.tar.gz",
        "sha256": "..."
      }
    }
  }
}
```

The index is the package-resolution authority. Lives in repo at `packages/index.json`, fetchable via raw GitHub URL. Updated at release time.

### 3. Lockfile

The lockfile stores package identity, not git transport detail:

```json
{
  "schema": "cn.lock.v2",
  "packages": [
    { "name": "cnos.core", "version": "3.34.0", "sha256": "..." },
    { "name": "cnos.eng", "version": "3.34.0", "sha256": "..." }
  ]
}
```

The lockfile says what to install. The package index says where to fetch it.

### 4. Restore flow

1. Read lockfile
2. Look up name+version in the package index
3. Download artifact over HTTPS
4. Verify SHA-256
5. Extract into `.cn/vendor/packages/<name>@<version>/`
6. Validate `cn.package.json`

No git required in the normal path.

### 5. Development override

For local development only:
- `cn deps restore --from-local <path>`
- Per-package local path override in manifest

This is the only development escape hatch needed in v1.

### 6. Commands content class

`commands` joins the existing content classes in `cn.package.json`:

```json
{
  "schema": "cn.package.v1",
  "name": "cnos.core",
  "version": "3.34.0",
  "sources": {
    "doctrine": ["*"],
    "mindsets": ["ENGINEERING.md", "..."],
    "skills": ["agent/cap", "agent/clp", "..."],
    "extensions": ["cnos.net.http"],
    "templates": ["SOUL.md", "USER.md"],
    "commands": {
      "daily": {
        "entrypoint": "commands/cn-daily",
        "summary": "Create or show the daily reflection thread"
      },
      "weekly": {
        "entrypoint": "commands/cn-weekly",
        "summary": "Create or show the weekly review thread"
      }
    }
  }
}
```

Command files live under `commands/` in the package tree, alongside `skills/`, `doctrine/`, etc.

### 7. Command discovery precedence

1. **Built-in commands** — core compiled commands, always authoritative
2. **Repo-local commands** — `.cn/commands/cn-<name>`
3. **Vendored package commands** — from installed package manifests

No PATH fallback in v1. If two external commands claim the same name at the same precedence level, that is a doctor error.

### 8. Help and doctor integration

- `cn help` lists built-ins, repo-local commands, and package commands with source and summary
- `cn doctor` reports: duplicate command names, missing entrypoints, non-executable command files, malformed command metadata

## Leverage

- **Eliminates #155's failure class entirely** — no git protocol dependency for package install
- **Eliminates #162's need for a separate plugin system** — commands are a package content class
- **Lockfile simplification** — drops `source_kind`, `rev`, `subdir` (27 references in cn_deps.ml)
- **Every restricted environment works** — sandboxes, containers, airgapped (with vendored packages)
- **Content-class consistency** — commands follow the same explicit pattern as doctrine/mindsets/skills/extensions/templates
- **Hosting portability** — lockfile stores name+version+hash, not URLs. Move hosting without touching lockfiles.

## Negative Leverage

- Release workflow must build package tarballs (N packages × build step)
- Package index must be maintained (generated, but still an artifact to manage)

## Alternatives Considered

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| Embed cnos.core in binary | Zero network bootstrap | Special-cases one package, doesn't fix distribution generally | Rejected |
| Fetch by tag instead of SHA | Minimal code change | Still uses git as transport, doesn't fix the real problem | Rejected |
| Separate command plugin framework | Clean separation | Second extension surface, more machinery | Rejected — commands as content class is simpler |
| PATH-based command discovery | Familiar (cargo, kubectl) | Harder to audit, no integrity checking, no doctor validation | Deferred to post-v1 if needed |
| Commands as implicit package tree content | No schema change | Violates explicit content-class pattern; smuggles commands in | Rejected |

## Process Cost / Automation Boundary

- Package tarball generation: automated in release workflow
- Index generation: automated in release workflow
- Checksum verification: automated in `cn deps restore`
- Command discovery: automated in `cn` dispatch
- Human judgment remains for: package content decisions, command API design, version policy

## Non-goals

- Arbitrary third-party Git repos as normal package sources
- PATH-wide command discovery in v1
- Third-party package hosting generally
- Redesigning runtime extensions
- Semantic version range solving in v1
- Package signing beyond checksum verification

## File Changes

### Create
- `packages/index.json` — package index
- `scripts/build-packages.sh` — build package tarballs from `packages/` dirs
- `src/cmd/cn_command.ml` — command discovery and dispatch

### Edit
- `src/cmd/cn_deps.ml` — replace git fetch with HTTP tarball fetch, new lockfile format, delete `source_kind`/`rev`/`subdir`
- `src/lib/cn_lib.ml` — lockfile v2 types, package index types, command content class types
- `src/cmd/cn_doctor.ml` — package command validation
- `src/cmd/cn_help.ml` — list package and repo-local commands
- `src/cli/cn.ml` — command dispatch precedence
- `.github/workflows/release.yml` — build and publish package tarballs + index
- `scripts/stamp-versions.sh` — stamp package artifact versions
- `scripts/check-version-consistency.sh` — validate package artifacts
- `docs/alpha/package-system/PACKAGE-SYSTEM.md` — add `commands` content class, document artifact-first distribution
- `packages/cnos.core/cn.package.json` — add `commands` content class entries
- `packages/cnos.eng/cn.package.json` — add `commands` content class entries if applicable
- `src/cmd/cn_runtime_contract.ml` — add skill activation table to cognition layer (scan exposed skills' frontmatter for trigger keywords)

### Delete
- Git fetch/clone logic in `cn_deps.ml`
- `source_kind`, `rev`, `subdir` from lockfile schema

## Acceptance Criteria

- [ ] AC1: Released first-party packages are published as tarball artifacts by the release workflow
- [ ] AC2: `packages/index.json` exists and resolves name+version to URL+SHA-256
- [ ] AC3: `cn deps restore` restores first-party packages via HTTP without git
- [ ] AC4: `cn deps restore` verifies package SHA-256 before install
- [ ] AC5: Extracted packages validate `cn.package.json`
- [ ] AC6: `commands` is an explicit content class in `cn.package.json` schema, documented in PACKAGE-SYSTEM.md
- [ ] AC7: `cn help` lists package-provided and repo-local commands
- [ ] AC8: `cn doctor` validates package command integrity (duplicates, missing entrypoints, non-executable, malformed metadata)
- [ ] AC9: `.cn/commands/` repo-local commands are discoverable
- [ ] AC10: Runtime capability extensions remain a separate system (no regression)

## Migration

One hub exists (sigma). No backward compatibility needed. Cut over:

1. Ship artifact-first restore + package index + command discovery
2. Regenerate lockfile in new format
3. Delete git fetch code

No phased migration. No v1/v2 lockfile coexistence.

## Known Debt

- Package index is a flat file; no versioned API or caching layer
- No third-party package hosting story yet
- No package signing beyond SHA-256 checksums
- Version solving is exact-match only (no ranges)

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|---------------|----------|
| 0 Observe | CHANGELOG TSC, open issues, #155 post-mortem | — | Git-based package fetch fails in restricted environments; #162 needs command extension |
| 1 Select | — | — | Package distribution is the root cause of #155; commands are not a content class (#162) |
| 4 Gap | this artifact | — | Released packages not yet the distribution unit; commands not yet a package content class |
| 5 Mode | this artifact | design, cap, writing | MCA, L7, system-shaping: artifact-first distribution + commands as content class |
| 6 Artifacts | PACKAGE-ARTIFACTS.md v0.2.0 | design | Design doc tightened per operator review |
