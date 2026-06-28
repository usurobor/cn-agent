# CDD Kata Framework

**Version:** 0.1.0
**Status:** Draft
**Package:** cnos.cdd.kata

## 0. Purpose

This document defines the kata framework used to prove two things:

1. **Runtime viability** — cnos can perform the full package loop: author → build → distribute → install → discover → dispatch → self-describe.
2. **CDD value** — given the same task, CDD produces better evidence, safer execution, and more durable closure than ad hoc execution.

The framework is designed to be:
- runnable by humans or machines
- judged by a fixed rubric
- repeatable every release
- distributable as a package

## 1. Principles

### 1.1 Executable truth first

The source of truth is executable kata commands/scripts. Human-readable docs explain them. They do not replace them.

### 1.2 Fixed artifact contract

Each kata defines:
- scenario
- inputs
- exact commands
- required artifacts
- expected output
- scoring rubric

### 1.3 Fixed judge contract

The judge scores evidence against the rubric. It does not invent new criteria.

### 1.4 Same task, two modes

Where the goal is to prove CDD value, the same task should be runnable in:
- baseline / ad hoc mode
- CDD / artifact-driven mode

That comparison is how the framework demonstrates value.

## 2. Two kata families

### 2.1 Runtime katas

These prove the package-driven runtime loop works mechanically. They answer:
- can the system build package artifacts?
- can it restore them?
- can it discover commands?
- can it dispatch commands?
- can it self-describe installed packages and discovered commands?

These are mostly binary pass/fail.

| ID | Name | Purpose |
|---|---|---|
| R1 | Boot | `cn init && cn deps restore` produces a working hub |
| R2 | Command | `cn build → cn deps restore → cn daily` dispatches from a package |
| R3 | Round-trip | add a new package command, build/install, then dispatch it |
| R4 | Doctor | break package/runtime state and verify `cn doctor` catches it |

### 2.2 Method katas

These prove CDD adds value beyond ad hoc execution. They answer:
- are the right artifacts produced?
- is the reasoning traceable?
- are review and release safer?
- is post-release closure stronger?

These are scored, not just pass/fail.

| ID | Name | Default level target | Purpose |
|---|---|---:|---|
| M1 | Design | L6 | prove design artifact quality and traceability |
| M2 | Review | L6 | prove evidence-bound review quality |
| M3 | Release | L6 | prove release as coherence delta, not just tagging |
| M4 | Post-release | L6 / L7 | prove closure and leverage measurement |
| M5 | Full cycle | L6 / L7 | prove end-to-end CDD execution on one bounded change |

## 3. Package distribution

The framework ships as `src/packages/cnos.cdd.kata/`. Its package contains:
- commands (`kata-list`, `kata-run`, `kata-judge`)
- one optional skill explaining framework use
- the kata corpus bundled inside command directories
- rubrics
- judge prompts
- fixtures

The package system currently has no generic docs/data class. So the authoritative shipped form is command payloads + bundled kata data.

## 4. Runtime surfaces

### 4.1 Commands

The package exposes three commands:

**kata-list** — Lists available katas and their classes.

**kata-run `<id>` `[--mode baseline|cdd]`** — Runs one kata: sets up fixtures, executes required commands, captures artifacts, writes a run bundle.

**kata-judge `<run-dir>`** — Loads the fixed rubric and judge prompt and produces a structured verdict.

### 4.2 Optional skill

An optional skill may explain how to use the framework, interpret verdicts, and add new katas correctly. But the commands remain the source of truth.

## 5. Kata contract

Each kata must contain:

```
kata.md           # human-readable scenario and required artifacts
rubric.json       # machine-readable scoring contract
judge-prompt.md   # fixed evaluation prompt for LLM judge
fixtures/         # input state
expected/         # expected outputs or structural expectations
run.sh            # executable entry
```

### 5.1 kata.md

Required fields: ID, Name, Class, Purpose, Mode(s) supported, Scenario, Inputs, Exact commands, Required artifacts, Expected output, Pass conditions.

### 5.2 rubric.json

Machine-readable scoring contract.

### 5.3 judge-prompt.md

Fixed evaluation prompt for an LLM judge.

### 5.4 fixtures/

Input state for the kata.

### 5.5 expected/

Expected outputs or structural expectations.

## 6. Run output contract

Every kata run writes a bundle:

```
.kata-runs/<run-id>/
  metadata.json
  stdout.log
  stderr.log
  artifacts/
  verdict.json
```

**metadata.json** includes: kata id, mode, runtime version, package versions, timestamp, platform, command exit codes.

**artifacts/** contains the actual produced files: design docs, review notes, release notes, status output, runtime contract output, any generated thread or verification output.

## 7. Judge model

The judge scores evidence, not intention.

### 7.1 Runtime kata verdicts

Binary: pass / fail. Optional details: missing artifacts, wrong dispatch, wrong status output, doctor false negative.

### 7.2 Method kata verdicts

Structured scoring:

```json
{
  "artifact_completeness": 0.0,
  "alpha_quality": 0.0,
  "beta_quality": 0.0,
  "gamma_quality": 0.0,
  "level_evidence": {
    "L5": false,
    "L6": false,
    "L7": false
  },
  "summary": "..."
}
```

### 7.3 Level derivation

- **L5** — local correctness
- **L6** — system-safe execution
- **L7** — system-shaping leverage

A level is granted only if the required evidence passes. The judge must not infer a level from style or ambition alone.

## 8. Proving CDD adds value

For method katas, the framework supports paired runs: baseline and CDD. Same task, same inputs, same scope. The comparison measures: artifact completeness, review findings, closure quality, level evidence.

This is the core proof mechanism for CDD value.

## 9. Initial kata suite

### R1 — Boot

Prove a fresh hub can be initialized, packages restored, and runtime state surfaced correctly.

### R2 — Command

Prove package command discovery/dispatch works: `cn build → cn deps restore → cn daily`.

### R3 — Round-trip

Prove a newly authored package command flows through source → build → install → dispatch.

### R4 — Doctor

Prove doctor catches broken package/runtime state.

### M1 — Design

Prove design artifacts are complete and traceable. Required: named incoherence, impact graph, acceptance criteria, CDD trace.

### M2 — Review

Prove review is evidence-bound and architecture-aware. Required: verdict, evidence, architecture check, active skill consistency.

### M3 — Post-release

Prove release closure and post-release assessment quality. Required: coherence measurement, lag update, process learning, closeout, production verification.

## 10. Governance

### 10.1 Adding a kata

Must define: scenario, exact commands, required artifacts, rubric, judge prompt, expected outputs.

### 10.2 Changing a rubric

Rubric changes must be versioned. Do not silently rescore past runs under changed criteria.

### 10.3 CI use

Runtime katas may run in CI as mechanical proof. Method katas may be run on release branches, candidate PRs, or benchmark runs.

## 11. Short form

- Runtime katas prove the package-driven agent runtime works.
- Method katas prove CDD adds value.
- Commands are the source of truth.
- The judge scores evidence, not vibes.
- Same task, same inputs, two modes: baseline vs CDD.
