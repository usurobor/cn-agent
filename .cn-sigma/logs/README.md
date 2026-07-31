# Sigma activation logs

Foreign-to-home channel for Sigma-at-cnos → Sigma-at-cn-sigma. Single writer per file; files sharded per day as `YYYYMMDD.md`. Append-only within each day; new day creates a new file.

Cursor on the home side (`cn-sigma:.cn-sigma/state/activations.md` → `last_read_foreign_log: <sha>`) is a Git commit SHA of cnos — it spans the directory naturally, so home walks `logs/` from cursor SHA forward.

See `cnos:docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md` for the full convention.
