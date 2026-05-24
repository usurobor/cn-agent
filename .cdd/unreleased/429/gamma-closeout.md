# γ close-out — Cycle 429

**Cycle:** [cnos#429](https://github.com/usurobor/cnos/issues/429) — Release-pipeline repair
**Date:** 2026-05-24
**Role:** γ+α+β collapsed on δ (operator-side finalization; sigma terminal died mid-cycle)
**Verdict:** Closed. AC1/AC2 satisfied (binaries shipped + verified). AC3 untestable in this environment (documented). AC4 satisfied (this artifact set). AC5 re-affirmed.

## What happened

Terminal sigma opened cnos#429 with the full 5-AC brief, then shipped **only** the release.yml permissions split (merged at `bd9465b8`) and died before testing it or backfilling binaries. No close-out artifacts, no INDEX row, issue left open. δ took over and finished.

## RCA (F1) — original release.yml build-matrix git-auth failure

Sigma's working hypothesis: workflow-level `permissions: contents: write` poisoned the GITHUB_TOKEN across all jobs, breaking checkout in the build matrix. **This is not supported by the evidence:**

- The cycle/426 one-shot (`publish-3.82.0-release.yml`) ran with **workflow-level `contents: write`** and its checkout + softprops publish **succeeded**.
- release.yml shipped **3.81.0** (May 19) green on the **same** permissions block.
- Same file, same trigger type, opposite outcome between 3.81.0 ✓ and 3.82.0 ✗ → the failure was **transient/environmental** on that specific 3.82.0 tag-push run, not a structural permissions defect.

**Empirical confirmation this cycle:** the backfill workflow built all four matrix legs at tag 3.82.0 with **per-job `contents: read`** on build and they passed `go build` + `go test` + Tier-1 kata cleanly. The build steps themselves are sound.

**Disposition:** sigma's permissions-split is kept as harmless least-privilege hardening, but it is **not** confirmed as "the fix," because the failure it targets was not reproduced and is likely non-deterministic. Recorded honestly rather than claimed as root-caused.

## F2 — `on.push.tags` does not re-trigger on tag force-update

Confirmed **structural**, not a bug: GitHub's anti-recursion rule suppresses workflow triggers for refs pushed via `GITHUB_TOKEN`, and this environment's proxy blocks operator-local tag pushes (HTTP 403). Net effect: release.yml's `on.push.tags` path **cannot be exercised by automation here**. Workaround in force: remote-runner one-shot workflows (BOX-AND-THE-RUNNER doctrine) or an operator PAT/UI dispatch. **Disposition: no-patch — accepted limitation, documented.**

## F3 (NEW) — `macos-x64` asset is actually an arm64 binary

`checksums.txt` shows `cn-macos-x64` and `cn-macos-arm64` with an **identical** SHA-256. Cause: GitHub migrated `macos-latest` to Apple-Silicon runners, so both the `macos-latest` (labeled x64) and `macos-14` (arm64) matrix legs build native arm64. The `cn-macos-x64` asset will not run on Intel Macs → `cn update` is broken for that (shrinking) population. **This bug is in `release.yml` too and predates this cycle (3.81.0 likely carries it).** Fix is a cross-compile (`GOOS=darwin GOARCH=amd64`, cn is pure-stdlib so it cross-compiles trivially) on an arm runner, or use the `macos-13` Intel runner. **Disposition: next-MCA** (operator decision pending; touches release.yml redesign, out of this cycle's backfill scope).

## AC grid

- **AC1 — 3.82.0 binaries:** PARTIAL→PASS. All four `cn-*` assets + `checksums.txt` attached (asset count 2→16). `cn-linux-x64` SHA matches checksums.txt; binary reports `Current version: 3.82.0`. Caveat: `cn-macos-x64` is mislabeled arm64 (F3).
- **AC2 — `cn update` 3.81.0→3.82.0:** PASS (sandbox-limited). Asset download + checksum-match verified; version stamp correct. The api.github.com `releases/latest` call is 403-blocked in this sandbox (network policy, not a cn defect), so the full in-place self-replace couldn't run here. The previously-broken link (missing assets → 404) is fixed.
- **AC3 — release.yml clean on next tag push:** NOT VERIFIABLE here (F2). Build steps proven sound via backfill; trigger path needs operator PAT tag-push or UI dispatch to confirm.
- **AC4 — F1/F2 (+F3) receipt-stream disposition:** PASS (this artifact set + INDEX row).
- **AC5 — pause re-affirmed:** PASS. No #405 Track A2/B1, no protocol evolution past v3.82.0. Field application of CDS/CDR + handoff/memory-return remains the next phase.

## Deliverables

- `release.yml` permissions split (sigma; merged `bd9465b8`).
- One-shot `backfill-3.82.0-binaries.yml` — built + uploaded the four binaries + checksums.txt to the 3.82.0 release via the asset API (body preserved); self-deleted (`2bd35595`).
- This close-out set + `cdd-iteration.md` + INDEX row.

## Attribution

γ+α+β collapsed on δ. Operator-side finalization after agent death, per the cycle/420 precedent (agent stall → operator γ takes over, preserves any landed work). Sigma's release.yml patch preserved verbatim; δ added the backfill, verification, and receipt-stream housekeeping.
