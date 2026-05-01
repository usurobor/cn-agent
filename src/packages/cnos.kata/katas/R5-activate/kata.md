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

## Inputs

The kata relies on a hub created by `cn init` or a fixture hub with
at minimum `.cn/config.json` and one identity file (`spec/SOUL.md`).

## Notes

`cn activate` was added in 3.70.0 (#320). This kata validates the
command exists and behaves correctly. It does not test prompt quality
or model compatibility — only structural correctness.
