## Gap

**Issue:** cnos#684 — "Sigma activation channel — rename `.cn-sigma/logs`, move to symmetric append-only orphan refs (independent stream, off main HEAD)."

**Version / mode:** Mode `explore` (issue frontmatter: `Mode: explore`, `Status: design-first (explore)`). This is a design-first cell, not an implementation cell for orphan refs or history rewriting. The deliverable is a design document + mechanism spec + dry-run migration plan, produced entirely inside the `cnos` repo on `cycle/684`.

**What is not-current about the pre-684 state:** `.cn-sigma/logs/` — dated `YYYYMMDD.md` files under `main`'s tree — is the Sigma foreign→home activation channel, but (a) the name "logs" is too generic for what is specifically an activation-stream, and (b) the channel is an independent communication timeline, not product state, yet it materializes in `main`'s HEAD tree (~1024 commits / ~10 MB, the dominant source of `main` HEAD churn per the operator's dispatch-authorization comment). The v0 convention (`AGENT-ACTIVATION-LOG-v0.md`) also assumes a bare-SHA cursor and a commit-to-`main` attach contract, both of which stop working once the stream moves off `main`.

**What this cycle closes:** the go-forward mechanism design (rename, orphan-ref names + writer/reader roles, amended registration schema, attach-contract sequence, orphan-ref invariants + enforcement, promotion boundary) plus a dry-run-only migration plan and a specified (not executed) content-preservation verification procedure for the already-committed `.cn-sigma/logs/**` history on `main`. It explicitly does **not** close: any physical write to `.cn-sigma/**`, any history rewrite, any ref creation/deletion, any force-push. Those remain operator-executed, separately-gated actions per the operator's binding dispatch-authorization comment (cnos#684, @2026-08-01T07:50:17Z).

**Dispatch chain:** operator authorization (2026-08-01T07:50:17Z, "κ dispatch authorization + bounded scope refinement") → γ R0 scaffold (`.cdd/unreleased/684/gamma-scaffold.md`) → δ inward-membrane enrichment confirming γ's two scope-resolution calls (`.cdd/unreleased/684/gamma-clarification.md`) → this α R0 cycle.
