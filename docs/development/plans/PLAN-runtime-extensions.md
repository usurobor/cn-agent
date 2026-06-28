# Implementation Plan for Runtime Extensions

**Status:** Draft
**Issue:** #73
**Depends on:** Runtime Contract v2, current package restore / doctor / traceability architecture

---

## 0. Coherence Contract

### Gap

cnos has converged on a package-native cognitive substrate, but runtime capabilities still grow by core accretion. Issue #73 identifies the architectural fix:

- packages can already carry doctrine/mindsets/skills/profiles
- runtime capabilities should become equally extensible
- network access (`http_get`, `dns_resolve`) should prove the architecture as the first extension rather than as a permanent core special case

### Mode

MCA — change the runtime architecture so capability families can be installed and loaded as extensions.

### α / β / γ target

- α PATTERN: explicit extension manifests, registry, dispatch, policy, and lifecycle
- β RELATION: package distribution, runtime contract, doctor, traceability, and execution all describe the same extension system
- γ EXIT: future capabilities become additive, and a later registry/app ecosystem can sit on this runtime layer without redesign

### Smallest coherent intervention

Do not rewrite built-in capabilities. Add:

- manifest schema
- extension discovery
- open op registry/dispatch
- subprocess host protocol
- one reference extension (`cnos.net.http`)

---

## 1. Deliverables

### Code

- manifest schema + parser
- extension registry builder
- open op dispatch in CN Shell
- subprocess extension host protocol
- first reference extension package (`cnos.net.http`)
- Runtime Contract integration
- doctor validation
- traceability events

### Docs

- canonical runtime extension spec already exists
- implementation plan (this file)
- updates to Runtime Contract / doctor / traceability docs as needed

### Tests

- manifest parse/validation
- discovery/conflict behavior
- dispatch routing
- policy intersection
- subprocess host protocol
- runtime contract surfacing
- doctor checks
- traceability events

---

## 2. Step Order

### Step 1 — Define the extension manifest schema

**Goal:** Create a machine-readable schema for extensions.

**Work:**
- add schema under docs/alpha/schemas/
- implement parser/validator in code
- support fields from the design:
  - name
  - version
  - interface
  - kind
  - backend
  - ops
  - permissions
  - engines

**Acceptance:**
- manifest can be parsed and validated
- malformed manifests fail clearly
- unknown / missing required fields are rejected

---

### Step 2 — Add extension source/build/install layout support

**Goal:** Make first-party extensions follow current cnos source discipline.

**Work:**
- define canonical source-of-truth layout under `src/agent/extensions/{pkg}/{ext-name}/`
- extend `cn build` so extension content is copied into generated package output
- ensure installed packages preserve:
  - `extensions/<ext-name>/cn.extension.json`
  - `host/`
  - `docs/`
  - `schemas/`

**Acceptance:**
- extension content is source-driven, not authored in generated packages
- generated packages contain extension layout deterministically
- installed layout under `.cn/vendor/packages/<pkg>@<ver>/extensions/` is stable

---

### Step 3 — Build the extension registry at boot

**Goal:** Discover installed extensions from the installed package tree.

**Work:**
- scan `.cn/vendor/packages/<pkg>@<ver>/extensions/**/cn.extension.json`
- build registry entries with:
  - extension identity
  - package source
  - backend kind
  - typed ops provided
  - permissions
  - compatibility status
  - enablement state
  - lifecycle state

**Acceptance:**
- discovery is eager
- activation remains lazy
- duplicate names/op kinds or incompatible engines/interfaces are rejected deterministically

---

### Step 4 — Introduce the open op registry / dispatch model

**Goal:** Replace the closed built-in-only typed-op vocabulary assumption with:
- stable common op envelope
- registry-based dispatch

**Work:**
- keep built-in ops registered by core
- register extension ops from manifests
- route by op kind through the registry
- support two-stage validation:
  1. core envelope validation
  2. provider-specific request validation

**Acceptance:**
- built-in ops still work unchanged
- extension-defined op kinds dispatch without core code edits
- op-kind conflicts are rejected, never shadowed

---

### Step 5 — Implement subprocess extension host protocol

**Goal:** Make subprocess the default extension backend.

**Work:**
Define a minimal host protocol, e.g. stdio JSON:
- describe
- health
- execute
- shutdown

Implement:
- process spawning
- request/response framing
- timeouts
- max-bytes/resource limits
- structured error handling

**Acceptance:**
- host startup is bounded
- host failure is isolated
- runtime receives structured health and execution errors
- no extension host can expand authority on its own

---

### Step 6 — Add policy / sandbox intersection

**Goal:** Ensure extensions declare needs but core runtime decides effective authority.

**Work:**
- intersect:
  - extension-declared permissions
  - runtime config
  - package/profile policy
- model states:
  - discovered
  - disabled
  - rejected
  - unavailable
  - enabled
  - loaded
  - active

**First reference policy set for cnos.net.http:**
- domain allowlist
- methods allowed
- timeout
- byte budgets
- allowed secrets
- explicit enable/disable

**Acceptance:**
- risky extensions can be discovered but disabled
- no ambient credentials
- secret injection is explicit and policy-controlled

---

### Step 7 — Runtime Contract integration

**Goal:** Surface installed and active extensions to the agent at wake.

**Work:**
Extend Runtime Contract v2:
- `cognition.extensions.installed`
- `body.extensions.active`
- capability surface includes extension-provided op kinds
- include enable/disable state and effective permission envelope

**Acceptance:**
At wake, the agent can determine:
- which extensions are installed
- which are enabled
- which ops they add
- what the effective permission envelope is

---

### Step 8 — cn doctor integration

**Goal:** Validate extension compatibility and configuration.

**Work:**
Teach doctor to check:
- installed manifests parse
- engine compatibility
- interface compatibility
- op-kind conflicts
- enable/disable state
- installed roots vs lockfile
- provider health (where safe/cheap)
- config/policy coherence

**Acceptance:**
Doctor can explain:
- what is installed
- what is enabled
- what is rejected
- why

---

### Step 9 — Traceability integration

**Goal:** Make extension lifecycle and execution reconstructable.

**Required events:**
- extension.discovered
- extension.loaded
- extension.rejected
- extension.disabled
- extension.health.ok
- extension.health.error
- extension.op.start
- extension.op.ok
- extension.op.error

**Acceptance:**
An operator can answer:
- was the extension discovered?
- was it rejected?
- was it disabled?
- what happened when its op ran?

---

### Step 10 — First reference extension: cnos.net.http

**Goal:** Prove the model with a useful but bounded extension.

**Package:** `cnos.net.http`

**Initial ops:**
- `http_get` (observe)
- `dns_resolve` (observe)

**Constraints:**
- subprocess backend
- domain allowlist
- read-only by default
- no ambient credentials

**Acceptance:**
- package installs cleanly
- ops appear in Runtime Contract
- ops dispatch through extension registry
- requests are sandboxed and traced
- no core runtime code changes are required for the op kinds themselves beyond the generic extension machinery

---

## 3. Tests

### Unit / schema tests
- manifest parsing
- invalid manifest rejection
- op-kind uniqueness / conflict policy
- two-stage validation

### Registry tests
- discovery from installed layout
- compatible / incompatible engine handling
- duplicate op rejection
- enable/disable states

### Host protocol tests
- describe
- health
- execute
- timeout/error behavior

### Runtime tests
- built-in op still dispatches correctly
- extension op dispatches correctly
- disabled extension op is rejected coherently
- Runtime Contract includes installed/active extension state

### Doctor tests
- malformed manifest
- rejected engine/interface
- disabled extension
- duplicate op conflict

### Traceability tests
- `extension.*` event family emitted correctly

---

## 4. CI / Release Gating

Before shipping:
- schema validation green
- registry/dispatch tests green
- doctor tests green
- traceability tests green
- cnos.net.http reference extension tests green

---

## 5. Non-goals

This plan does not include:
- full registry/marketplace design
- publication workflow
- bundle/app install UX
- native plugin backend in v1
- migration of existing built-in capabilities out of core

Those belong to later phases or companion docs.

---

## 6. Acceptance Criteria

This implementation is complete when:

1. A package can ship an extension without modifying core runtime code.
2. cnos discovers installed extensions automatically at boot.
3. The runtime can dispatch extension-defined op kinds without core code changes.
4. The Runtime Contract tells the agent which extension ops exist and whether they are enabled.
5. Extension execution is sandboxed and traced.
6. cn doctor validates extension compatibility and configuration.
7. `cnos.net.http` proves the model by shipping `http_get` as the first extension.

---

## 7. Summary

This plan keeps the cnos core minimal while making capability growth coherent:

- schema
- discovery
- open dispatch
- subprocess isolation
- policy intersection
- Runtime Contract surfacing
- doctor
- traceability
- first reference extension

That is the smallest implementation that proves the Runtime Extensions architecture is real.
