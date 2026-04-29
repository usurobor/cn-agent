# cnos Documentation

A recurrent coherence system with Git as its lowest durable substrate.

**Version:** 3.61.1
**Date:** 2026-04-29

---

## Start Here

**[THESIS.md](./THESIS.md)** — what cnos is. The whole, above the triad.

Then choose your path:

| You want to... | Start with |
|----------------|------------|
| Understand what cnos is | [THESIS.md](./THESIS.md) → [COHERENCE-SYSTEM.md](./alpha/essays/COHERENCE-SYSTEM.md) → [FOUNDATIONS.md](./alpha/essays/FOUNDATIONS.md) |
| Build or run a cnos agent | [AGENT-RUNTIME.md](./alpha/agent-runtime/AGENT-RUNTIME.md) → [CLI.md](./alpha/cli/CLI.md) → [HANDSHAKE.md](./beta/guides/HANDSHAKE.md) |
| Contribute code | [CDD.md](./gamma/cdd/CDD.md) → [ARCHITECTURE.md](./beta/architecture/ARCHITECTURE.md) |
| Understand agent composition / CTB | [CTB README](./alpha/ctb/README.md) → [v0.1 spec](./alpha/ctb/LANGUAGE-SPEC.md) → [v0.2 draft](./alpha/ctb/LANGUAGE-SPEC-v0.2-draft.md) |
| Understand the formal foundation | [TSC repo](https://github.com/usurobor/tsc) → [C≡](https://github.com/usurobor/tsc/blob/main/spec/c-equiv.md) → [TSC Core](https://github.com/usurobor/tsc/blob/main/spec/tsc-core.md) → [TSC Oper](https://github.com/usurobor/tsc/blob/main/spec/tsc-oper.md) |
| Understand the runtime extensions model | [RUNTIME-EXTENSIONS.md](./alpha/runtime-extensions/RUNTIME-EXTENSIONS.md) → [runtime-extensions bundle](./alpha/runtime-extensions/) |
| Write or modify a skill | [WRITE-A-SKILL.md](./beta/guides/WRITE-A-SKILL.md) → [COGNITIVE-SUBSTRATE.md](./alpha/cognitive-substrate/COGNITIVE-SUBSTRATE.md) |
| Do a release | [release skill](../src/packages/cnos.cdd/skills/cdd/release/SKILL.md) → [BUILD-RELEASE.md](./beta/guides/BUILD-RELEASE.md) |

---

## Reading Model

The docs tree has two dimensions:

### Triad axis — ontological character

Every document has a dominant axis. Read top-down for breadth.

| Axis | Directory | Question |
|------|-----------|----------|
| **Whole** | `docs/` | What is cnos? |
| **α Pattern** | `docs/alpha/` | What has been articulated? |
| **β Relation** | `docs/beta/` | Do the parts cohere? |
| **γ Evolution** | `docs/gamma/` | How does it change? |

### Feature bundles — depth within a feature

Related documents (spec, snapshots, design docs) are grouped into **feature bundles** within `alpha/`. Read into a bundle for depth on a single feature.

| Bundle | Canonical spec | What it covers |
|--------|---------------|----------------|
| [agent-runtime/](./alpha/agent-runtime/) | [AGENT-RUNTIME.md](./alpha/agent-runtime/AGENT-RUNTIME.md) | CN Shell, typed ops, N-pass orchestration, receipts |
| [runtime-extensions/](./alpha/runtime-extensions/) | [RUNTIME-EXTENSIONS.md](./alpha/runtime-extensions/RUNTIME-EXTENSIONS.md) | Capability providers, discovery, isolation |

See [DOCUMENTATION-SYSTEM.md](./beta/governance/DOCUMENTATION-SYSTEM.md) for the full taxonomy and rules.

---

## The Triad

### α — Pattern: What Is Articulated

The substance of the system — doctrine, specs, definitions.

*"What has been made explicit? Is it internally consistent?"*

**Anchor docs** (canonical, evolve in place):

| Document | Scope |
|----------|-------|
| [COHERENCE-SYSTEM.md](./alpha/essays/COHERENCE-SYSTEM.md) | Meta-model: coherence as primary; the instruction set |
| [FOUNDATIONS.md](./alpha/essays/FOUNDATIONS.md) | The coherence stack — doctrinal layers |
| [MANIFESTO.md](./alpha/essays/MANIFESTO.md) | Principles and values |
| [CAA.md](./alpha/agent-runtime/CAA.md) | Coherent agent architecture |
| [AGENT-RUNTIME.md](./alpha/agent-runtime/AGENT-RUNTIME.md) | Runtime spec: CN Shell, typed ops, N-pass orchestration, receipts |
| [RUNTIME-EXTENSIONS.md](./alpha/runtime-extensions/RUNTIME-EXTENSIONS.md) | Capability providers, discovery, and isolation |
| [COGNITIVE-SUBSTRATE.md](./alpha/cognitive-substrate/COGNITIVE-SUBSTRATE.md) | Cognitive asset classes — doctrine, mindsets, skills |
| [CAR.md](./alpha/cognitive-substrate/CAR.md) | Cognitive asset resolver — local, versioned cognition |
| [WHITEPAPER.md](./alpha/protocol/WHITEPAPER.md) | CN protocol specification (v3.0.0) |
| [PROTOCOL.md](./alpha/protocol/PROTOCOL.md) | FSM design — state diagrams, transition tables |
| [TRACEABILITY.md](./alpha/security/TRACEABILITY.md) | Observability — event stream, state projections, readiness |
| [SECURITY-MODEL.md](./alpha/security/SECURITY-MODEL.md) | Security architecture — sandbox, FSM enforcement, audit trail |
| [CLI.md](./alpha/cli/CLI.md) | CLI command reference |
| [DAEMON.md](./alpha/cli/DAEMON.md) | Daemon mode design |
| [SETUP-INSTALLER.md](./alpha/cli/SETUP-INSTALLER.md) | Install script specification |
| [THREAD-API.md](./alpha/protocol/THREAD-API.md) | Agent content API |
| [MESSAGE-PACKET-TRANSPORT.md](./alpha/protocol/MESSAGE-PACKET-TRANSPORT.md) | Fail-closed inbox materialization — packet-based transport (#150) |

**Feature-scoped design docs** (version-stamped, see [migration path](./beta/governance/DOCUMENTATION-SYSTEM.md#6-migration-path-for-legacy-filenames)):

| Document | Feature | Version |
|----------|---------|---------|
| [RUNTIME-CONTRACT-v3.10.0.md](./alpha/agent-runtime/3.10.0/DESIGN.md) | Wake-up self-model contract | 3.10.0 |
| [N-PASS-BIND-v3.8.0.md](./alpha/agent-runtime/3.8.0/N-PASS-BIND.md) | N-pass bind loop and indicators | 3.8.0 |
| [SYSCALL-SURFACE-v3.8.0.md](./alpha/agent-runtime/3.8.0/SYSCALL-SURFACE.md) | Syscall surface redesign | 3.8.0 |
| [SCHEDULER-v3.7.0.md](./alpha/agent-runtime/3.7.0/DESIGN.md) | Scheduler design | 3.7.0 |
| [CTB-v4.0.0-VISION.md](./alpha/ctb/CTB-v4.0.0-VISION.md) | CTB v4.0.0 vision: agent-composition language | 4.0.0 |

**Feature bundles:**

| Bundle | Contents |
|--------|----------|
| [agent-runtime/](./alpha/agent-runtime/) | Runtime spec, CAA, runtime contract, version-scoped design docs |
| [cli/](./alpha/cli/) | CLI reference, daemon mode, setup installer |
| [cognitive-substrate/](./alpha/cognitive-substrate/) | Cognitive asset classes, CAR resolver |
| [ctb/](./alpha/ctb/) | CTB — triadic agent-composition language (draft); [v0.1 baseline](./alpha/ctb/LANGUAGE-SPEC.md), [v0.2 draft](./alpha/ctb/LANGUAGE-SPEC-v0.2-draft.md), [notes](./alpha/ctb/SEMANTICS-NOTES.md) |
| [doctrine/](./alpha/doctrine/) | Doctrine sub-packages (coherence, ethics, judgment, inheritance) |
| [essays/](./alpha/essays/) | System doctrine: coherence system, foundations, manifesto |
| [protocol/](./alpha/protocol/) | Whitepaper, protocol FSMs, thread API |
| [runtime-extensions/](./alpha/runtime-extensions/) | Extensions spec, version snapshots |
| [security/](./alpha/security/) | Security model, traceability |

### β — Relation: How Parts Cohere

The graph of relations — system overview, shared vocabulary, operator connection, model-reality evidence.

*"Do the parts reveal one system? Does model match reality?"*

| Document | Scope |
|----------|-------|
| [architecture/](./beta/architecture/) | System overview — how the parts relate |
| [governance/](./beta/governance/) | Doc system, naming conventions, glossary |
| [lineage/](./beta/lineage/) | Origin narrative and structural lineage |
| [schema/](./beta/schema/) | LLM schema design |
| [guides/](./beta/guides/) | Operator ↔ system: handshake, automation, migration, skills, dojo, build/release |
| [evidence/](./beta/evidence/) | Model ↔ reality: audit, RCAs |

### γ — Evolution: How the System Moves

The process that moves the system — method, plans, gates.

*"How does it change without losing itself?"*

| Bundle | Contents |
|--------|----------|
| [cdd/](./gamma/cdd/) | CDD method, agile process, frozen snapshots |
| [rules/](./gamma/rules/) | Non-negotiable rules, coherence invariants |
| [essays/](./gamma/essays/) | Position papers: stateless agency, executable skills |
| [plans/](./gamma/plans/) | Implementation plans for specific releases |
| [checklists/](./gamma/checklists/) | Review checklists by discipline |

---

## Directory Map

```
docs/
├── THESIS.md                          # The whole
├── README.md                          # This file — navigation
├── alpha/                             # α Pattern
│   ├── agent-runtime/                 # Runtime spec, CAA, runtime contract
│   │   ├── README.md
│   │   ├── AGENT-RUNTIME.md           # Canonical spec (v3.8.0)
│   │   ├── CAA.md                     # Coherent agent architecture
│   │   ├── RUNTIME-CONTRACT-v2.md     # Runtime contract v2
│   │   └── 3.x.0/                    # Version-scoped design docs
│   ├── cli/                           # CLI, daemon, installer
│   │   ├── CLI.md, DAEMON.md, SETUP-INSTALLER.md
│   ├── cognitive-substrate/           # Cognitive assets, CAR
│   │   ├── COGNITIVE-SUBSTRATE.md, CAR.md
│   ├── ctb/                           # CTB agent-composition language (draft)
│   │   ├── README.md                  # Document map + authority
│   │   ├── LANGUAGE-SPEC.md           # v0.1 baseline (normative)
│   │   ├── LANGUAGE-SPEC-v0.2-draft.md # v0.2 agent-module target (draft)
│   │   ├── SEMANTICS-NOTES.md         # Conceptual rationale (non-normative)
│   │   └── CTB-v4.0.0-VISION.md      # Strategy + roadmap (non-normative)
│   ├── doctrine/                      # Doctrine sub-packages
│   │   └── coherence-for-agents/, ethics-for-agents/, ...
│   ├── essays/                        # System doctrine and long-form essays
│   │   ├── COHERENCE-SYSTEM.md        # Meta-model: coherence as primary
│   │   ├── FOUNDATIONS.md             # Core doctrine and coherence loop
│   │   └── MANIFESTO.md              # Human+AI commons, sovereignty
│   ├── protocol/                      # Whitepaper, protocol, thread API
│   │   ├── PROTOCOL.md, WHITEPAPER.md, THREAD-API.md
│   ├── runtime-extensions/            # Extensions spec + snapshots
│   │   ├── README.md
│   │   ├── RUNTIME-EXTENSIONS.md      # Canonical spec (v1.0.6)
│   │   └── 1.0.6/                    # Version directory
│   │       ├── README.md             # Snapshot manifest
│   │       └── SPEC.md               # Frozen spec at v1.0.6
│   ├── security/                      # Security model, traceability
│   │   ├── SECURITY-MODEL.md, TRACEABILITY.md
│   └── schemas/                       # JSON schemas
├── beta/                              # β Relation
│   ├── architecture/                  # System architecture
│   │   ├── README.md
│   │   ├── ARCHITECTURE.md
│   │   └── 3.14.4/
│   ├── governance/                    # Doc system, naming, glossary
│   │   ├── README.md
│   │   ├── DOCUMENTATION-SYSTEM.md
│   │   ├── NAMING.md
│   │   ├── GLOSSARY.md
│   │   └── 3.14.4/
│   ├── lineage/                       # Origin and structural lineage
│   │   ├── README.md
│   │   ├── ORIGIN.md
│   │   ├── LINEAGE.md
│   │   └── 3.14.4/
│   ├── schema/                        # LLM schema design
│   │   ├── README.md
│   │   ├── DESIGN-LLM-SCHEMA.md
│   │   └── 3.14.4/
│   ├── guides/                        # Operator guides
│   └── evidence/                      # Audits, RCAs
└── gamma/                             # γ Evolution
    ├── cdd/                           # CDD method, agile process
    │   ├── README.md
    │   ├── CDD.md
    │   ├── RATIONALE.md
    │   └── 3.x.x/                    # Frozen snapshots
    ├── rules/                         # Project rules, invariants
    │   ├── RULES.md, INVARIANTS.md
    ├── essays/                        # Position papers
    │   ├── STATELESS-AGENCY.md, EXECUTABLE-SKILLS.md
    ├── plans/                         # Implementation plans
    └── checklists/                    # Review checklists
```
