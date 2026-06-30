---
name: test-declaration-only
description: "Synthetic declaration-only wake fixture (cnos#524 W3). The typed SKILL.md twin of this directory's wake-provider.json, added when W3 flipped cn-install-wake's default source to SKILL.md so the AC5 refusal smoke exercises the SKILL.md default path (not the legacy JSON path). Never rendered to a substrate workflow — the renderer refuses (exit 3) on activation_state: declaration-only before prompt inlining."
governing_question: Does the renderer refuse to render a declaration-only wake (exit 3) when reading from the SKILL.md default source?
artifact_class: wake
scope: global
kata_surface: none
triggers:
  - admin-wake
inputs:
  - "n/a (declaration-only fixture; never invoked)"
outputs:
  - "n/a (declaration-only fixture; refused at exit 3 before any output)"
wake:
  role: admin
  package: cnos.core
  admin_only: true
  activation_log_writer: false
  activation_state: declaration-only
  input:
    triggers:
      - schedule
  output:
    channel_log_convention: "n/a (declaration-only fixture)"
    writer_surface: "n/a (declaration-only fixture; never writes)"
    class_taxonomy:
      - n/a
    cursor_advance: false
    cursor_field: "n/a (declaration-only fixture)"
  permission_intent:
    - contents.write
  concurrency:
    serialize: true
    group: "test-declaration-only"
  agent_variable:
    name: agent
    default: null
  surfaces:
    allowed:
      - "n/a (declaration-only fixture)"
    disallowed:
      - "everything (declaration-only fixture)"
  defer_path:
    cell_shaped_directive: "n/a (declaration-only fixture; defer to operator for everything)"
    off_role_directive: "n/a (declaration-only fixture)"
    ambiguous_directive: "n/a (declaration-only fixture)"
---

# Synthetic test fixture — declaration-only

This body is never inlined into any rendered workflow: the `wake.activation_state`
is `declaration-only`, so the renderer refuses (exit 3) at the activation_state
gate (per `wake-provider/SKILL.md` §3.10) before reaching the prompt-inlining
stage. It exists so the renderer's SKILL.md body-extraction step has a body to
read, and to keep the AC5 negative-case smoke alive against the W3 default
(`--source skill`).

## Activation preconditions (would be required to flip to `live`)

This fixture deliberately documents preconditions — cnos#454 and cnos#467 — so
the contract it represents matches the legacy `wake-provider.json` twin in this
directory. The renderer's refusal message names the activation_state value and,
when the synthesized manifest carries no `activation_state_notes`, states that
the manifest should name the preconditions for flipping to live. The JSON twin
is retained for W3 dual-source parity and is deleted in W4.
