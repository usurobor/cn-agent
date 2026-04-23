# Polyglot Packages and Provider Contracts

**Issue:** #170
**Version:** 0.2.0
**Mode:** MCA
**Active Skills:** cdd/design, eng/evolve, eng/process-economics
**Engineering Level:** L7

## Problem

cnos is converging on a package-driven runtime:

- package artifacts are the unit of distribution
- packages can carry doctrine, mindsets, skills, commands, orchestrators, extensions, and templates
- the runtime contract exposes installed/activatable cognition and executable body surfaces
- the Go kernel is becoming the small trusted bootstrap/runtime substrate

But the current design still leaves one important ambiguity unresolved:

> Does "package-driven" imply "same implementation language everywhere," or can the kernel stay Go while package-provided commands and providers are implemented in the language best suited to their job?

This matters now because the likely next shape of the system is:

- **Go** for the kernel/runtime:
  - package manager
  - daemon
  - web surface
  - command registry/dispatch
  - runtime contract
  - doctor/status/update
- **Rust** for some package-provided tools/providers:
  - e.g. git-cn
  - transport-heavy, performance-sensitive, capability-oriented tools
- **packages** continuing to carry non-binary assets:
  - skills
  - doctrine
  - templates
  - orchestrators

Without an explicit design, the system risks drifting into one of two bad shapes:

1. **Monolingual pressure:** everything must be Go because the kernel is Go, even when another language would be better for a package-provided tool/provider
2. **Polyglot drift without contracts:** commands/providers appear in many languages, but discovery, packaging, execution, and policy boundaries become inconsistent

The gap is not "can cnos run binaries?" The gap is:

> What are the explicit contracts that let the kernel remain small and trusted while packages provide commands and providers in whatever language fits the job?

## Decision

Use a Go kernel with language-agnostic packages. The kernel remains small and trusted. Packages remain the only distribution/install substrate. Commands and providers may be implemented in any language that can satisfy the cnos execution contract.

### Core rule

**The package is the unit of distribution. The implementation language inside the package is an implementation detail.**

### Runtime rule

Commands and providers are different runtime surfaces:

- **commands** are exact-dispatch operator actions
- **providers** are executable capability surfaces the runtime can invoke or host

They may ship inside the same package. They do not share the same execution contract.

### Authority rule

The kernel still owns:

- registry construction
- policy
- precedence
- protocol semantics
- runtime contract emission
- trust boundaries

Packages provide implementations. The kernel defines the rules under which they run.

## Constraints

- Keep the bootstrap kernel in Go.
- Do not create a second package system for non-Go artifacts.
- Do not collapse commands into runtime extensions.
- Do not collapse providers into commands.
- Do not require package authors to use reflection, dynamic loading, or classloader tricks.
- Preserve built-in command precedence over repo-local and package-provided commands.
- Preserve package-system content classes as the common install substrate.
- Keep low-level protocol semantics kernel-owned even if concrete transport drivers become package providers.
- Keep execution contracts simple enough that any agent can build and reason about them.

## Challenged Assumption

The challenged assumption is:

> "If the kernel is Go, then every command/provider worth shipping should also be Go."

That assumption is unnecessary. The kernel should be monolingual only in the places where cnos needs one small trusted implementation substrate. The package boundary is the correct place for polyglot behavior.

## Prior Art

This design borrows the right ideas from several families of systems:

- **Git** — exact-dispatch command extensibility through external executables and metadata
- **Terraform / HashiCorp plugin systems** — versioned subprocess protocols for capability providers
- **NuGet / Helm** — package artifacts and package indexes as the distribution/install authority

The goal is to combine the best of these without:

- classloader complexity
- reflection-heavy registration
- runtime resolution magic
- multiple package ecosystems

## Proposal

### 1. One package substrate, polyglot payloads

Packages remain the one installable/distributable unit. A package may carry:

- doctrine
- mindsets
- skills
- commands
- orchestrators
- extensions
- templates

Nothing in the package substrate assumes the language of:

- a command binary
- a provider binary
- a helper executable

The package metadata and runtime contracts define behavior. The payload language remains free.

#### Example

```text
packages/cnos.transport.git/
  cn.package.json
  commands/
    git-cn/
      git-cn-linux-x64
      git-cn-linux-arm64
      git-cn-macos-x64
      git-cn-macos-arm64
  extensions/
    cnos.transport.git/
      manifest.json
```

This package is first-party, but not part of the trusted kernel.

### 2. The kernel is a package too

`cnos.kernel` remains an embedded `cn.package.v1` manifest in the Go binary. This keeps the platform inside its own model.

The kernel manifest uses the same schema and parser as normal packages. The only difference is installability: it is embedded, not restored.

#### Kernel responsibilities

The kernel owns:

- package restore/install
- command registry
- provider registry
- orchestrator registry
- runtime contract
- doctor/status/help/update
- daemon/web surface
- trust and precedence policy
- protocol semantics

#### Kernel non-responsibilities

The kernel does not need to implement every command or provider itself.

### 3. Commands vs providers

#### 3.1 Commands

Commands are:

- exact-dispatch operator actions
- invoked as `cn <name>`
- packaged as command entries
- discovered through kernel/repo-local/package precedence
- suitable for one-shot execution

Examples:

- `daily`
- `weekly`
- `save`
- `git-cn` as an operator tool
- `telegram-send` as an operator tool

#### 3.2 Providers

Providers are:

- executable capability surfaces
- invoked by the runtime, not by command-line dispatch
- suitable for long-lived or protocol-based interaction
- represented in the runtime contract as body/capability or provider surfaces

Examples:

- `cnos.net.http`
- `cnos.transport.git`
- `cnos.llm.local`
- `cnos.llm.anthropic`

#### 3.3 Rule

Commands and providers may live in the same package. They do not become the same runtime thing.

### 4. Command execution contract

Every command, regardless of source, is normalized into one runtime descriptor model.

#### Runtime descriptor

```go
type CommandSource string

const (
    SourceKernel   CommandSource = "kernel"
    SourceRepoLocal CommandSource = "repo-local"
    SourcePackage  CommandSource = "package"
)

type CommandTier int

const (
    TierKernel CommandTier = iota
    TierRepoLocal
    TierPackage
)

type CommandSpec struct {
    Name      string
    Summary   string
    Source    CommandSource
    Tier     CommandTier
    Package  string
    NeedsHub bool
    Dangerous bool
}
```

#### Execution interface

```go
type Invocation struct {
    HubPath string
    Args    []string
    Stdin   io.Reader
    Stdout  io.Writer
    Stderr  io.Writer
    Env     map[string]string
}

type Command interface {
    Spec() CommandSpec
    Run(ctx context.Context, inv Invocation) error
}
```

#### Rule

The runtime command model is one. Source forms differ:

- kernel manifest + constructor map
- repo-local convention (optional sidecar later)
- package manifest metadata

### 5. Provider execution contract

Providers need a separate contract.

#### 5.1 Why separate

A provider is not just a command with a different name. It may need:

- capability declaration
- health checks
- version handshake
- long-lived subprocess or socket lifecycle
- receipt/error protocol
- request/response framing

#### 5.2 Provider declaration

There is no separate provider manifest. Providers are declared via `cn.extension.v1` with a `capabilities` field:

```json
{
  "schema": "cn.extension.v1",
  "name": "cnos.transport.git",
  "version": "0.1.0",
  "interface": "cn.ext.v1",
  "kind": "capability-provider",
  "backend": {
    "kind": "subprocess",
    "command": ["bin/git-cn-linux-x64"]
  },
  "capabilities": [
    "transport.git.sync",
    "transport.git.packet"
  ],
  "ops": [
    {
      "kind": "git_sync",
      "class": "effect",
      "request_schema": "schemas/git_sync.json",
      "response_format": "artifact:json"
    }
  ],
  "permissions": {
    "git": true,
    "default_read_only": false,
    "allow_secrets": []
  },
  "engines": {
    "cnos": ">=3.40.0 <4.0.0"
  }
}
```

`capabilities` is optional but strongly recommended. See RUNTIME-EXTENSIONS.md §5.1a.

#### 5.3 Runtime model

The kernel launches providers as subprocesses and speaks a versioned protocol to them. No reflection. No in-process plugin loading. No classloader model.

#### 5.4 First-class example

`git-cn` should be treated as a provider candidate, not just a command:

- **command form** for operator use
- **provider form** for runtime transport use

This is the exact kind of case where Rust is a good implementation choice without forcing the kernel to be Rust.

### 6. A2A / transport split

This design sharpens the earlier A2A split.

**Kernel-owned:**

- protocol semantics
- envelope rules
- receipts
- routing invariants
- dedupe / rejection policy
- trust boundaries

**Package/provider-owned:**

- concrete transport drivers
- Git transport implementation
- HTTP transport implementation
- Telegram transport implementation if that ever exists

That means:

- A2A semantics stay core
- A2A transport implementations can become providers

This is the cleanest way to support `git-cn` as a package/provider without hollowing out the kernel.

### 7. Artifact structure and platform targeting

Because packages may carry native binaries, the package system must allow platform-specific payloads.

#### Rule

The package identity is platform-neutral:

- name
- version
- content classes
- metadata

The package artifact resolution may be platform-specific.

#### Package index implication

The package index should be able to resolve:

- package metadata
- package artifact URL
- platform-specific binary payloads where relevant

This does not require multiple package ecosystems. It only requires platform-aware package artifacts.

### 8. Runtime contract implications

The runtime contract should eventually expose:

- installed package commands
- installed providers
- provider availability / health
- command source/tier/package
- routing/provider choices where relevant

This makes the runtime self-model explicit and auditable.

### 9. Migration path

**Phase 1:**

- keep Go as the kernel language
- allow package-provided commands to be polyglot binaries
- normalize command registration through one runtime descriptor model

**Phase 2:**

- add provider contract and provider registry
- move one concrete capability/provider into package form
- likely first example: `cnos.transport.git`

**Phase 3:**

- teach runtime contract / doctor / status to surface provider state
- move more concrete transports or integrations to package/provider form

## Leverage

This design gives cnos:

- one trusted kernel language
- one package substrate
- one command registry model
- one provider registry model
- freedom to use the best language for the job at package boundaries
- no pressure to rewrite everything in one language
- no second plugin system

It lets:

- Go own the boring kernel
- Rust own the transport-heavy native tools where appropriate
- packages remain the common install/distribution unit

## Negative Leverage

This adds:

- a second execution contract (commands vs providers)
- platform-aware package payload concerns
- more explicit packaging discipline
- more doctor/status/runtime-contract surface area

It also requires resisting the temptation to let every package invent its own execution model.

## Alternatives Considered

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| Make everything Go | One language | Forces bad language choices for package-provided tools/providers | Rejected |
| Make commands and providers the same thing | Simpler surface | Blurs exact-dispatch tools with capability providers | Rejected |
| Separate package systems for kernel and polyglot tools | Clear on paper | Two ecosystems, two indexes, two restore paths | Rejected |
| Go kernel + polyglot packages + separate command/provider contracts | Simple kernel, flexible packages, one substrate | More explicit contracts to maintain | **Chosen** |

## File Changes

**Create:**

- `docs/alpha/agent-runtime/POLYGLOT-PACKAGES-AND-PROVIDERS.md`

**Edit:**

- `docs/alpha/package-system/PACKAGE-SYSTEM.md`
  - clarify language-agnostic payload principle
  - mention platform-aware payloads where relevant
- `docs/alpha/runtime-extensions/RUNTIME-EXTENSIONS.md`
  - clarify provider contract vs package distribution
- `docs/alpha/agent-runtime/RUNTIME-CONTRACT-v2.md`
  - later, expose provider registry / health
- command/provider design docs
  - align with one descriptor model for commands
  - one provider contract for providers

## Acceptance Criteria

- cnos documents one package substrate with language-agnostic payloads
- command execution contract is explicit and distinct from provider contract
- `cnos.kernel` remains package-compatible while embedded
- package-provided commands can be implemented in non-Go languages
- provider packages can declare executable providers without becoming commands
- low-level protocol semantics remain kernel-owned
- transport implementations can move to provider packages without redefining protocol truth
- the package system does not split into separate ecosystems for Go vs non-Go payloads

## Known Debt

- exact provider protocol still needs its own dedicated spec
- platform-aware package artifact resolution needs a concrete index format later
- first provider migration (likely `cnos.transport.git`) is still future work
- doctor/status/runtime-contract surfacing of providers remains follow-up work

## CDD Trace

| Step | Artifact | Skills loaded | Decision |
|------|----------|--------------|----------|
| 0 | Observe | — | Current design now points to a Go kernel, but package/tool/provider needs are becoming polyglot |
| 1 | Select | — | Selected gap: the package/runtime design does not yet explicitly explain how a Go kernel hosts non-Go commands/providers coherently |
| 4 | Gap | this artifact | Named incoherence: package substrate and runtime execution contracts are not yet explicit enough for a polyglot future |
| 5 | Mode | this artifact | L7 MCA; kernel/package/provider architecture clarification |
| 6 | Artifacts | this artifact | Design drafted; implementation sequencing remains future work |

---

If you want, the next useful artifact is the **provider contract spec** for packages like `cnos.transport.git`, since that is the place where this design becomes concrete.
