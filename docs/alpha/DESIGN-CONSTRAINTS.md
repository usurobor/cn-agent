# Design Constraints

Active constraints governing cnos architecture and interfaces. Referenced by CDD review (§2.2.13) and issue skills.

---

## §1. Source of truth

Each fact lives in one place. Derived copies must not diverge from the authoritative source.

## §2. Package system

### §2.1 One package substrate

All packages use `cn.package.v1` schema. One manifest format, one build pipeline, one install path.

### §2.2 Source / artifact / installed clarity

What is authored (`src/`), what is built (`dist/`), and what is installed (`.cn/vendor/`) are always distinct and explicit.

## §3. Command surface

### §3.1 Git-style subcommands

Commands use `cn <noun> <verb>` (git-style): `cn kata run`, `cn cdd verify`, `cn deps lock`.

Internal registry key is hyphenated (`kata-run`). User-facing surface is space-separated. `cn <noun>` with no verb lists available subcommands.

### §3.2 Dispatch boundary

`cli/` owns dispatch. Domain packages own logic. No `cmd_*.go` file may contain domain logic beyond argument parsing and dispatch.

## §4. Runtime surfaces

### §4.1 Surface separation

Skills, commands, orchestrators, and providers are distinct runtime surfaces with separate contracts. Transport implementation is distinct from protocol semantics.

### §4.2 Registry normalization

Different source forms normalize into one runtime descriptor model. Help, doctor, status, and runtime contract reveal the same registry truth.

## §5. Language and runtime

### §5.0 OCaml is deprecated

The OCaml codebase (`src/ocaml/`, `test/cmd/*.ml`) is deprecated. **Do not modify, extend, or fix OCaml files.** All new work and all fixes to existing behavior must be in Go (`src/go/`). The OCaml code will be removed once the Go port is complete. If a fix appears to require changing OCaml, the correct action is to implement the equivalent in Go and mark the OCaml path as superseded.

## §6. Architecture

### §6.1 Reason to change

Each module/package/boundary has one real reason to change. No convenience buckets.

### §6.2 Policy above detail

Policy remains in the kernel. Packages cannot widen their own authority.

### §6.3 Degraded-path visibility

Fallback and degraded behavior is explicit, testable, and inspectable via doctor/status.

---

_This document is the canonical constraints reference. CDD review §2.2.13 loads it. Issue skill §2.4.2 links to it. When a constraint is revised, name the revision explicitly._
