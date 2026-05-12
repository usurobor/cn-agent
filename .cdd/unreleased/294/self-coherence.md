# Self-Coherence — Cycle 294

## Gap

**Issue:** #294 — cn cdd status N — single-shot cycle TLDR (read-only γ tooling)

**Version/Mode:** Standard substantial change cycle per CDD.md §1.1 — introduces new command tooling, spans design/code/tests/docs, requires CDD artifact set.

**Selected Gap:** Missing γ command for cycle state introspection. γ assembles cycle TLDR manually (read issue, git rev-parse, list artifacts, evaluate closure gate, summarize). Information is derivable from git + filesystem + GitHub API but no `cn` command produces structured TLDR. Operator and γ both ask "TLDR current state" 1-2x per cycle without tooling.

**Gap Type:** Tooling gap — mechanical read-only command to project existing state into structured TLDR format.

**Governing CDD.md Clause:** Tooling that reduces manual repetitive coordinator work without changing process semantics (supports §1.4 γ algorithm efficiency).