---
name: test-log-writer-misdeclaration
description: "Synthetic mis-declared wake fixture (cnos#524). The wake SKILL.md the AC4/cycle-496 mis-declaration smoke renders from. Declares role:dispatch + admin_only:false + activation_log_writer:true — the mis-declaration — so the renderer refuses with exit 4 + stderr 'activation_log_writer mis-declaration:'. Never rendered to a substrate workflow."
governing_question: Does the renderer refuse (exit 4) an activation_log_writer mis-declaration when reading from the SKILL.md default source?
artifact_class: wake
scope: global
kata_surface: none
triggers:
  - dispatch-wake
inputs:
  - "n/a (synthetic fixture; never invoked)"
outputs:
  - "n/a (synthetic fixture; refused at exit 4 before any output)"
wake:
  role: dispatch
  package: cnos.cds
  admin_only: false
  activation_log_writer: true
  activation_state: live
  protocol: cds
  selector:
    include:
      - dispatch:cell
      - protocol:cds
      - status:todo
    exclude: []
  input:
    triggers:
      - schedule
  output:
    cycle_artifact_root: ".cdd/unreleased/{N}/"
    artifact_class_taxonomy:
      - test
    cell_runtime: cnos.cdd
  permission_intent:
    - contents.write
    - issues.write
    - pull_requests.write
    - id_token.write
  concurrency:
    serialize: true
    group: "test-log-writer-misdeclaration"
  agent_variable:
    name: agent
    default: sigma
  surfaces:
    allowed:
      - "n/a (synthetic fixture)"
    disallowed:
      - "everything (synthetic fixture; demonstrates refusal, not rendered)"
  defer_path:
    cell_shaped_directive: "n/a (synthetic fixture)"
    off_role_directive: "n/a (synthetic fixture)"
    ambiguous_directive: "n/a (synthetic fixture)"
---

# Synthetic test fixture — activation_log_writer mis-declaration

This body is never inlined into any rendered workflow: the wake declares
`role: dispatch` + `admin_only: false` + `activation_log_writer: true` — the
mis-declaration — so the renderer refuses (exit 4) at the activation_log_writer
gate (per `wake-provider/SKILL.md` §3 + cycle/496) before reaching the
prompt-inlining stage. It exists so the renderer's SKILL.md body-extraction
step has a body to read, and to keep the AC4/cycle-496 negative-case smoke
alive against the W3 default source (`--source skill`). The JSON twin in this
directory is retained for W3 dual-source parity and is deleted in W4.
