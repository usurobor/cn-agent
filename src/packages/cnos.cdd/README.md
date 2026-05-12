# CDD — Coherence-Driven Development

**CDD (Coherence-Driven Development)** is the γ-axis evolution machinery for how the cnos codebase moves coherently. It defines the method, captures evidence, and provides the triadic role structure (α/β/γ) that governs substantial changes.

## What CDD Does

CDD is the development method used to evolve cnos coherently. Its purpose is not merely to ship features, but to reduce incoherence across the system as a whole: doctrine, architecture, implementation, runtime behavior, operator understanding, release state, and development process itself.

## Package Structure

This package contains two main surfaces:

### `/skills/` — The Method
- **`cdd/CDD.md`** — Canonical algorithm spec for development cycles
- **Role skills** — `alpha/`, `beta/`, `gamma/`, `operator/` define triadic responsibilities  
- **Lifecycle skills** — `issue/`, `review/`, `release/`, `design/`, `plan/` implement phases
- **Activation** — `activation/SKILL.md` for one-time repository bootstrap

### Repository Root `.cdd/` — The Evidence
- **`.cdd/releases/{version}/{cycle}/`** — Per-release evidence directories ("receipts")
- **`.cdd/unreleased/{N}/`** — Active cycle coordination artifacts
- **`.cdd/waves/`** — Wave coordination for related issue sequences

## For Essay Readers

If you arrived here from the agent-first essay, the `.cdd/releases/` directories contain the per-cycle artifacts that demonstrate measured coherence improvement: self-coherence reports, review records, close-outs, and post-release assessments. These are the "receipts" the essay describes—durable evidence that substantial changes were executed through the CDD method rather than informally.

## Quick Start

**New to CDD?** Start with:
1. `skills/cdd/CDD.md` — Core algorithm and role definitions
2. `activation/SKILL.md` — One-time repository setup
3. Browse `.cdd/releases/` — See CDD in practice across real cycles

**Looking for a specific cycle?** Check `.cdd/releases/{version}/{issue-number}/` for the complete artifact set from that cycle.

## License

Part of the cnos project.