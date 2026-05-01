# R5 — Activate kata

**Tier:** 2 — runtime/package
**Class:** runtime
**Proves:** `cn activate` generates a bootstrap prompt from local hub state.

## Scenario

Given a hub created by `cn init`, `cn activate` produces a model-ready
bootstrap prompt on stdout without invoking a model, starting a daemon,
or leaking secrets.

## Pass conditions

### P1: cwd hub discovery

Run `cn activate` from inside a hub directory.

- Exit code 0.
- stdout contains the activation header (`You are activating a cnos hub.`).
- stdout contains the hub path.
- stderr may contain diagnostics but stdout is prompt-only.

### P2: explicit HUB_DIR

Run `cn activate HUB_DIR` from outside the hub (e.g. `/tmp`).

- Exit code 0.
- stdout contains the activation header and the explicit hub path.

### P3: stdout/stderr separation

Run `cn activate HUB_DIR > prompt.md` and inspect `prompt.md`.

- File contains only the prompt (no warnings, progress, or diagnostics).

### P4: missing hub fails clearly

Run `cn activate` from a directory that is not a hub, with no argument.

- Exit code non-zero.
- stderr contains a diagnostic mentioning hub discovery failure.
- stdout is empty.

### P5: explicit bad path fails

Run `cn activate /nonexistent/path`.

- Exit code non-zero.
- stderr contains a diagnostic.
- stdout is empty.

### P6: no secrets in prompt

Run `cn activate HUB_DIR` on a hub that contains `.cn/secrets.env` and `.env`.

- stdout does not contain the contents of those files.
- stdout does not contain tokens, keys, or passwords from those files.

### P7: triad split — three sections, no ## Identity bucket

Run `cn activate HUB_DIR` on a sigma-shape hub (spec/PERSONA.md,
spec/OPERATOR.md, vendored cnos.core/doctrine/KERNEL.md).

- stdout contains `## Kernel`, `## Persona`, `## Operator` sections.
- stdout does not contain `## Identity`.
- stdout does not contain `no identity files found`.

### P8: kernel states

**P8a: vendored**
Run on a hub with `.cn/vendor/packages/cnos.core/doctrine/KERNEL.md`
present and a `cn.package.json` declaring a version.

- `## Kernel` section contains `vendored at` with the path and `@<version>`.
- stdout does not suggest `cn deps restore`.

**P8b: manifest-only**
Run on a hub with `.cn/deps.json` declaring `cnos.core` but no vendor directory.

- `## Kernel` section contains `dependency manifest declares cnos.core; not restored — run cn deps restore`.
- stdout does not contain `vendored at`.

**P8c: none**
Run on a hub with no deps.json and no vendor.

- `## Kernel` section contains `no kernel reference`.
- stdout does not contain `cn deps restore`.

### P9: deps states

**P9a: restored**
Run on a hub with `.cn/deps.json` and at least one installed package
under `.cn/vendor/packages/`.

- `## Dependencies` section contains `restored: <package-list>`.

**P9b: manifest-only**
Run on a hub with `.cn/deps.json` but no installed vendor packages.

- `## Dependencies` section contains `dependency manifest present; packages not restored — run cn deps restore`.

**P9c: none**
Run on a hub with no `.cn/deps.json`.

- `## Dependencies` section contains `no dependency manifest`.

### P10: read-first ordering

Run `cn activate HUB_DIR` on a sigma-shape hub with a daily reflection present.

- stdout contains `## Read first` section.
- In that section, persona appears before operator, operator before kernel,
  kernel before deps manifest, deps manifest before latest reflection.
- The order must be: 1) PERSONA.md 2) OPERATOR.md 3) KERNEL.md 4) deps.json
  5) latest reflection (when present).

### P11: latest reflection pointer

**P11a: present**
Run on a hub with two files under `threads/reflections/daily/` (e.g.
`2026-04-29.md` and `2026-05-01.md`).

- `## Read first` section includes the path to `2026-05-01.md`.
- `2026-04-29.md` is not listed as the latest.

**P11b: absent**
Run on a hub with an empty or absent `threads/reflections/daily/`.

- stdout does not contain `not present` for the reflection entry.
- The latest reflection line is omitted rather than shown as absent.

## Inputs

The kata relies on a hub created by `cn init` or a fixture hub with
at minimum `.cn/config.json`.

For P7–P11, use these fixture shapes:

- **sigma-shape**: `spec/PERSONA.md`, `spec/OPERATOR.md`,
  `.cn/vendor/packages/cnos.core/doctrine/KERNEL.md`,
  `.cn/vendor/packages/cnos.core/cn.package.json` (with version),
  `.cn/deps.json` declaring cnos.core
- **init-only**: `spec/SOUL.md` only, no deps.json, no vendor, no PERSONA, no OPERATOR
- **init+setup**: `spec/SOUL.md` + `.cn/deps.json` declaring cnos.core, no vendor

## Notes

`cn activate` was added in 3.70.0 (#320). The Kernel/Persona/Operator
triad was introduced in 3.71.0 (#321). P1–P6 validate structural
correctness; P7–P11 validate the triad design.

This kata validates command behavior, not prompt quality or model
compatibility.
