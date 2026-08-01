---
name: measurement-missing-instruction
description: Measurement fixture for TSC-compatible methodology validation.
governing_question: Does the bounded methodology fixture resolve?
artifact_class: measurement
kata_surface: none
triggers: [measure]
scope: task-local
inputs: [target]
outputs: [report]
methodology:
  registry: "schemas/fixtures/skill-frontmatter/measurement-assets/valid-registry.tsc"
  calibration_preflight: "schemas/fixtures/skill-frontmatter/measurement-assets/preflight.sh"
  targets: ["fixture-target"]
  cross_target: false
  instruction: "schemas/fixtures/skill-frontmatter/measurement-assets/absent-instruction.md"
  path_base: repository-root
  activation: source-checkout-only
  target_root_contract: explicit-coh-root
  output_root: ".tsc/fixture"
  default_mode: hybrid
  consistency:
    mechanical: identical
    llm_repeats: 3
    llm_spread: exact fixture
  mechanical:
    backend: pinned fixture
    determinism: exact fixture
    signals:
      alpha: ["alpha.fixture"]
      beta: ["beta.fixture"]
      gamma: ["gamma.fixture"]
  llm:
    estimates: ["target", "alpha", "beta", "gamma", "delta_alpha_beta", "delta_beta_gamma", "delta_gamma_alpha", "bottleneck_axis", "confidence", "summary", "axis_evidence", "defect_cards", "unresolved_ambiguity", "next_fixes"]
    must_not: ["guess"]
    validation: exact fixture
    providers:
      local: fixture
      ci: fixture
    ci_prompt: fixture
---

# Measurement fixture
