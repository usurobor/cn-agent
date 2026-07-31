# Round 04 — F17 micro-closure + convergence declaration

**Branch:** `claude/repo-cleanup-newcomer`
**Predecessor:** Round-03 (β→α→γ) closed the `docs/` + root newcomer surface. Round-03 γ ran the first **repo-wide** sweep (β had scoped only `docs/` + root) and surfaced exactly one residual β missed — **F17** — with the explicit off-ramp that it could ride the code/ops pass if `.cn-sigma/` is ruled operational substrate.

## Scope adjudication (orchestrator)

F17 is **in scope for the docs pass**, off-ramp declined:

- The target is `.cn-sigma/logs/README.md` — **a directory README (documentation), not code, schema, package skill, or `src/` behavioral content.** The operator's deferral was explicitly "code"; a broken link in a README is not code.
- The cite is **current-tense** ("See X for the full convention"), not a frozen dated log. The frozen sigma logs are the dated `YYYYMMDD.md` files; the directory README is standing documentation.
- `.cn-sigma/` is **surfaced in the root `README.md` tree map (line 125)**, so a newcomer reaches this README directly. A dead link here violates operator criterion #3 ("all links must work") and #5 ("friendly/clear to newcomers").

Deferring a one-line link fix to a "code pass" when it has nothing to do with code behavior would be miscategorization, not scope discipline.

## Matter

Single mechanical swap — target verified on disk (`docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md`, present):

| File:line | Before | After |
|---|---|---|
| `.cn-sigma/logs/README.md:7` | `cnos:docs/gamma/conventions/AGENT-ACTIVATION-LOG-v0.md` | `cnos:docs/reference/conventions/AGENT-ACTIVATION-LOG-v0.md` |

No full β→α→γ triad was spawned for Round 04: Round-03 γ already performed the exhaustive independent adversarial repo-wide sweep and specified the fix with the on-disk target verified. Re-spawning three agents to swap one string would be the exact noise-generation the cleaning loop exists to eliminate (write skill: minimum text required; kernel: smallest real fix). γ's Round-03 sweep **is** the adversarial check; this closure applies its one specified fix and re-verifies.

## Verification (re-derived, orchestrator)

- **F17 cleared:** `grep gamma/conventions .cn-sigma/logs/README.md` → empty.
- **Definitive live-authority sweep** (`docs/` + root + `.cn-sigma` READMEs, excluding frozen buckets): **zero** live-authority retired-dir cites remain.
- **Sole repo-wide survivors** outside the closed newcomer surface all bucket cleanly:
  - frozen version-keyed records — `CHANGELOG`, `RELEASE.md`, dated plans (`docs/development/plans/PLAN-v3.22.0-eng-lane-clarity.md` `Status: Complete`, et al.), `.cdd/iterations/**` work-records, dated `.cn-sigma/logs/YYYYMMDD.md`;
  - intentional retirement / forbid statements (`DOCUMENTATION-SYSTEM.md`, `GLOSSARY.md`);
  - F16 authoring-intent prose (`CCNF-AND-TYPED-TRUST.md:383`, `DECREASING-INCOHERENCE.md:528`);
  - subset-(C) `docs/gamma/cdd/` snapshot convention — a live behavioral contract in package skills/schemas, owned by the deferred **code pass**.

## Convergence verdict

**CONVERGED.** The `docs/` + structure + root + newcomer-reachable surface is pristine against all six operator criteria. The docs cleaning loop terminates here.

Remaining `docs/gamma/cdd/` references and all `src/` / `schemas/` / package-skill content are the explicitly-deferred **code pass** — out of scope for this cleanup by operator instruction ("don't touch code yet").

## Cell lineage (docs cleaning loop)

```
R1  85655a5c β → ca450f92 α → 1a1caba7 γ   (F1–F7)
R2  d633d2ec β → 94a6862b α → e0849745 γ   (F8–F11)
R3  893d4978 β → db673386 α → <this> γ+R4  (F12–F16, + F17 repo-wide residual)
R4  F17 micro-closure + convergence         (this file)
```

Docs surface converged. Next: open a PR for operator review of the whole cleanup as one coherent diff; the code pass (subset-C `docs/gamma/cdd/` + `src/`/`schemas/`) is a separate follow-on.
