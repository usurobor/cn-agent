# Go Kernel Command Architecture

## The kernel is a package

**Version:** 1.1.0
**Status:** Draft
**Addresses:** #192 (Go kernel rewrite), #209 (CLI dispatch)
**Related:**
- docs/reference/packages/PACKAGE-SYSTEM.md (package content classes)
- src/cmd/cn_command.ml (current OCaml command dispatch)
- src/lib/cn_lib.ml (current OCaml command variants)

---

## 0. Core Principle

**The platform is an instance of its own runtime command model, not a special case outside it.**

Every command in cnos — kernel, repo-local, or package-vendored — is normalized into the same runtime command descriptor, registered into the same registry, and dispatched through the same command execution interface. The source forms differ. The runtime model does not.

---

## 1. Why

### 1.1. The problem

The OCaml CLI has two command models:

- **Built-in commands**: hardcoded variant type in `cn_lib.ml`, match arm in `cn.ml`, implemented as OCaml function calls. No manifest, no entrypoint file.
- **Package commands**: declared in `cn.package.json` under `commands`, entrypoint is a shell script (`cn-<name>`), discovered at runtime from `.cn/vendor/packages/`, dispatched by execing the entrypoint.

This split means:
- Adding a built-in command touches the variant type, the parser, and the dispatch match — three files minimum
- Adding a package command touches one `cn.package.json` and one script — two files, self-contained
- `cn help` must special-case built-ins vs discovered commands
- Testing built-ins requires the full binary; testing package commands requires only the entrypoint
- The two models have different shapes, different lifecycles, different documentation

### 1.2. The direction already taken

v3.37.0 migrated `daily`, `weekly`, `save` from built-in to `cnos.core` package commands. The built-in set was explicitly shrinking toward a bootstrap kernel: `help`, `init`, `setup`, `deps`, `build`, `doctor`, `status`, `update`. The content-class model already treats commands as a package surface.

### 1.3. The logical endpoint

If commands are a content class, and the kernel is shrinking toward a minimal set, then the kernel's commands should follow the same runtime command model. The kernel is a package.

---

## 2. Design

### 2.1. Runtime command descriptor

All command sources — kernel, repo-local, package — normalize into the same runtime descriptor:

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
    Tier      CommandTier
    Package   string // package name, empty for kernel/repo-local
    NeedsHub  bool   // false = available without a hub (help, init, setup)
    Dangerous bool   // requires confirmation or elevated context
}
```

Kernel manifest, package manifest, and repo-local convention all produce `CommandSpec`. The source forms differ; the runtime descriptor is the same. Repo-local commands may remain convention-based in v1. If needed later, they may add optional sidecar metadata.

### 2.2. Command runner interface

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

`NeedsHub` and similar properties live in `Spec()`, not in ad hoc dispatch logic. IO streams, env, and hub path travel via `Invocation`, not as bare function parameters.

### 2.3. One registry

```go
type Registry struct {
    commands map[string]Command
    order    []string // insertion order for help
}

func (r *Registry) Register(cmd Command)
func (r *Registry) Lookup(name string) (Command, bool)
func (r *Registry) All() []Command
func (r *Registry) Available(hasHub bool) []Command
```

The registry stores all commands. `All()` returns everything (for doctor/diagnostics). `Available()` filters by current runtime state (e.g., `NeedsHub` commands excluded when no hub is found). Help and dispatch use `Available()`.

Registration and availability are distinct. Examples:
- `deps` is registered but unavailable without a hub
- `help` is always available
- `init` and `setup` are available without a hub

The runtime makes this visible instead of smearing it into hidden dispatch rules.

### 2.4. Three source forms, one model

| Source | Source form | Runtime normalization |
|--------|-----------|----------------------|
| `cnos.kernel` | Embedded `cn.package.json` + constructor map | `CommandSpec` + Go implementation |
| Repo-local | `.cn/commands/cn-<name>` (+ optional sidecar later) | Synthetic `CommandSpec` + exec-backed implementation |
| Package-vendored | Vendored `cn.package.json` `commands` entries | `CommandSpec` + exec-backed implementation |

The source forms differ. The runtime descriptor is the same.

### 2.5. The kernel package

`cnos.kernel` is a `cn.package.v1`-compatible manifest embedded in the binary. It uses the same schema, the same `commands` object, and the same `sources.commands` array as any vendored package. It is not installable through the package index, but it is manifest-compatible with the package model — the same parser that reads `cnos.core/cn.package.json` can read the kernel manifest.

The embedded manifest is the source of truth for:

- command names
- summaries
- metadata (`needs_hub`, `dangerous`)
- the `sources.commands` list (which commands this package provides)

```json
{
  "schema": "cn.package.v1",
  "name": "cnos.kernel",
  "version": "3.43.0",
  "kind": "kernel",
  "description": "cnos bootstrap kernel — commands that exist before any packages are installed",
  "commands": {
    "help":   { "summary": "Show available commands", "needs_hub": false },
    "init":   { "summary": "Initialize a new hub", "needs_hub": false },
    "setup":  { "summary": "Configure cnos", "needs_hub": false },
    "deps":   { "summary": "Manage package dependencies" },
    "build":  { "summary": "Build packages from source" },
    "doctor": { "summary": "Validate hub health" },
    "status": { "summary": "Show hub status" },
    "update": { "summary": "Update cnos binary", "needs_hub": false }
  },
  "sources": {
    "commands": ["help", "init", "setup", "deps", "build", "doctor", "status", "update"]
  }
}
```

Notes:
- `schema: cn.package.v1` — real package schema, not a custom kernel format
- `kind: kernel` — distinguishes from installable packages without breaking schema compatibility
- `sources.commands` — same split as other packages (ids in sources, metadata in top-level `commands`)
- No `engines` field — the kernel does not depend on itself
- Not in `packages/index.json` — not installable, but inspectable with standard package tooling

Its constructor map is the source of truth for implementation binding:

```go
var kernelConstructors = map[string]func() Command{
    "help":   NewHelp,
    "init":   NewInit,
    "setup":  NewSetup,
    "deps":   NewDeps,
    "build":  NewBuild,
    "doctor": NewDoctor,
    "status": NewStatus,
    "update": NewUpdate,
}
```

This is explicit and boring. No reflection, no `init()` magic, no code generation. Adding a kernel command is: one file + one manifest entry + one constructor map entry.

### 2.6. Startup sequence

```
1. Load kernel manifest (embedded)
   → Parse command metadata into CommandSpec values

2. Bind kernel command implementations via explicit constructor map
   → Register kernel Commands into registry

3. Discover hub (walk up from cwd to find .cn/)
   → If no hub: registry has all commands, but NeedsHub ones are unavailable

4. Discover repo-local commands (.cn/commands/cn-<name>)
   → Produce synthetic CommandSpec + exec-backed Command
   → Register; skip if name conflicts with kernel (kernel wins)

5. Discover vendor package commands (.cn/vendor/packages/*/cn.package.json)
   → Parse commands entries into CommandSpec + exec-backed Command
   → Register; skip if name conflicts with kernel or repo-local

6. Dispatch: registry.Lookup(args[0])
   → Check availability (NeedsHub vs current state)
   → If unavailable: clear error naming why
   → If available: cmd.Run(ctx, invocation)
```

### 2.7. Precedence

Flat 3-tier, not a tree:

1. **Kernel** (`cnos.kernel`) — always wins. These are bootstrap commands.
2. **Repo-local** (`.cn/commands/cn-<name>`) — operator overrides.
3. **Package-vendored** — installed package commands.

Same precedence as the current OCaml implementation (`cn_command.ml`). Conflicts at the same tier are doctor errors, not silent shadowing.

### 2.8. Availability

Registration and availability are distinct:

| Command | NeedsHub | Available without hub |
|---------|----------|----------------------|
| `help`  | false    | ✅ always |
| `init`  | false    | ✅ always |
| `setup` | false    | ✅ always |
| `deps`  | true     | ❌ needs hub |
| `build` | true     | ❌ needs hub |
| `doctor`| true     | ❌ needs hub |
| `status`| true     | ❌ needs hub |
| `update`| false    | ✅ always (updates binary, not hub) |

`help` shows all commands with availability state, or only available commands, by policy. Either way, the runtime makes the distinction explicit — no hidden "exists but silently fails" behavior.

### 2.9. Bootstrap constraint

`cnos.kernel` must be self-sufficient before any packages are installed. This means:

- `deps restore` is a kernel command (it installs packages)
- `help` is a kernel command (it works without packages)
- `init` and `setup` are kernel commands (they create the hub)

The test: if you delete `.cn/vendor/`, can you still `cn deps restore` to get it back? If yes, the kernel is correctly self-sufficient.

---

## 3. Prior art

### 3.1. Git external subcommands

`git foo` looks for `git-foo` on `$PATH`. Built-in commands and external subcommands coexist under one invocation grammar. The unification is at dispatch, not source format — built-ins are C functions compiled into the binary, externals are separate executables.

**What cnos takes:** dispatch-level unification. The operator types `cn foo` regardless of whether `foo` is kernel, repo-local, or package.
**What cnos adds:** a runtime command descriptor that normalizes metadata across sources. Git has no equivalent — external subcommands have no structured metadata.

### 3.2. Cargo custom subcommands

`cargo foo` dispatches to `cargo-foo` if no built-in matches. Exact-dispatch extensibility without language/runtime complexity.

**What cnos takes:** naming convention for extensibility (`cn-<name>`).
**What cnos adds:** package-manifest-declared commands (richer than naming convention alone).

### 3.3. kubectl plugins

kubectl uses a `plugin.yaml` descriptor and explicit plugin listing. Precedence and visibility questions are close to cnos's help/registry problem.

**What cnos takes:** explicit listing with availability state. `cn help` is like `kubectl plugin list` — it shows what's available, not just what's installed.
**What cnos adds:** tiered precedence (kernel > repo-local > package).

### 3.4. npm bin

Packages expose commands through the `bin` field in `package.json`. Install puts them on `$PATH`. The closest analogy for vendored package commands.

**What cnos takes:** package metadata as the source of command declarations.
**What cnos adds:** kernel commands follow the same metadata shape (embedded manifest).

### 3.5. The shared principle

All of these systems unify **dispatch** across sources. None of them unify **source format** — built-ins, plugins, and package commands are authored differently but invoked the same way. cnos follows this pattern: one runtime command descriptor, one registry, one dispatch interface, three source forms.

The additional principle cnos adds: **the platform's own commands are described by the same metadata model as user/package commands.** Git and Cargo don't do this — their built-ins have no user-visible manifest. cnos's kernel has an embedded `cn.package.json`, making it inspectable with the same tools.

---

## 4. Repo structure

```
go/
  cmd/cn/
    main.go                 # entrypoint: hub discovery, registry setup, dispatch
    kernel.json             # cn.package.json for cnos.kernel (embedded via //go:embed)
  internal/
    cli/
      command.go            # CommandSpec, Invocation, Command interface
      registry.go           # Registry: Register, Lookup, All, Available
      registry_test.go
      discover.go           # scan vendor + repo-local for commands
      discover_test.go
      exec_command.go       # exec-backed Command for package/repo-local commands
    kernel/
      kernel.go             # loads embedded manifest, constructor map, registers kernel commands
      cmd_deps.go           # deps: restore, list, doctor
      cmd_help.go           # help: list all registered commands
      cmd_doctor.go         # doctor (Phase 3)
      cmd_status.go         # status (Phase 3)
      cmd_update.go         # update (Phase 3)
      cmd_init.go           # init (Phase 3)
      cmd_setup.go          # setup (Phase 3)
      cmd_build.go          # build (Phase 3)
    pkg/                    # unchanged — pure types
    restore/                # unchanged — restore logic
```

### 4.1. Adding a kernel command

One file + one manifest entry + one constructor map entry:

1. Create `internal/kernel/cmd_<name>.go` implementing `Command`
2. Add entry to `cmd/cn/kernel.json` under `commands`
3. Add constructor to `kernelConstructors` map in `kernel.go`

Explicit and boring. No magic.

### 4.2. Adding a package command

One script + one manifest entry (unchanged from current model):

1. Create `src/agent/commands/<name>/cn-<name>` (shell script)
2. Add entry to package's `cn.package.json` under `sources.commands` and `commands`
3. `cn build` + `cn deps restore` makes it available

Same shape, slightly different touchpoints.

---

## 5. Phasing

| Phase | What | Commands |
|-------|------|----------|
| Phase 2 (#209) | CLI entrypoint + registry + kernel scaffold | `deps restore`, `help` |
| Phase 3 | Remaining kernel commands | `doctor`, `status`, `update`, `init`, `setup`, `build` |
| Phase 4 | Package command discovery | repo-local + vendor scanning, exec-backed `Command` |
| Phase 5 | Feature parity | all OCaml built-in behavior reimplemented in Go |

Phase 4 is when the "one model" promise is fully realized — kernel and package commands coexist in the same registry at runtime.

---

## 6. Open questions

1. **Should `cnos.kernel` appear in `packages/index.json`?** It's not installable via `deps restore`, but listing it would make the index a complete catalog of all cnos packages. Leaning no — the index is for installable packages.

2. **Should kernel commands be overridable?** Currently kernel always wins. An operator might want to replace `deps restore` with a custom version. Leaning no for now — kernel commands are bootstrap-critical. Repo-local can extend but not replace.

3. **Version stamping:** The kernel manifest version should match the binary version. Use `//go:embed` for the manifest and stamp version at build time via `-ldflags` or a generated file. Same pattern as OCaml's dune-generated version module.

4. **Repo-local sidecar metadata:** v1 uses convention only (`cn-<name>` existence). Future versions could support an optional `.cn/commands/cn-<name>.json` sidecar with summary, NeedsHub, etc. Not needed yet.
