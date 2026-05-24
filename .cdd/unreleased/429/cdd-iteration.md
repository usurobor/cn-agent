# cdd-iteration — Cycle 429

**Cycle:** [cnos#429](https://github.com/usurobor/cnos/issues/429) — Release-pipeline repair
**Date:** 2026-05-24
**Mode:** release-hygiene; γ+α+β collapsed on δ (operator finalization after sigma terminal death)
**Authority:** cnos.handoff/skills/handoff/receipt-stream/SKILL.md (cdd-iteration shape); cnos#426 cdd-iteration (F1/F2 origin)

## §0 Findings

`protocol_gap_count: 3`.

| ID | Finding | Class | Disposition |
|----|---------|-------|-------------|
| F1 | release.yml build-matrix git-auth failure on the 3.82.0 tag-push run | cdd-tooling-gap | **Refined/refuted.** Not the permissions defect sigma hypothesized (cycle/426 + 3.81.0 counterexamples; backfill build green under per-job read perms). Likely transient. Permissions-split kept as hardening. No further patch. |
| F2 | `on.push.tags` does not re-trigger on tag force-update | substrate-limitation | **No-patch.** GITHUB_TOKEN anti-recursion + proxy 403 on operator tag-push. Accepted; workaround = remote-runner one-shots / operator PAT. |
| F3 | `cn-macos-x64` release asset is an arm64 binary (macos-latest is now Apple Silicon) | cdd-tooling-gap | **next-MCA.** Affects release.yml. Fix = cross-compile `GOOS=darwin GOARCH=amd64` or use macos-13. Out of backfill scope. |

## §1 Trigger assessment

- F1 originated as a next-MCA carry-forward from cnos#426. Reassessed against new evidence this cycle and downgraded from "structural permissions bug" to "transient, unreproduced." The honest receipt corrects the earlier inference rather than ratifying it.
- F3 is newly surfaced by inspecting the backfilled `checksums.txt` (duplicate SHA across the two macOS targets). Mechanical discovery, not review churn.

## §2 What actually fixed the user-visible breakage

`cn update` for 3.81.0→3.82.0 was broken because the 3.82.0 release shipped **notes only** — no binaries. The backfill one-shot built and attached `cn-{linux,macos}-{x64,arm64}` + `checksums.txt`. Verified: `cn-linux-x64` downloads, hash matches, reports `3.82.0`. The release.yml permissions-split (sigma) is orthogonal hardening, not the thing that fixed `cn update`.

## §3 Pause status (AC5)

Re-affirmed. No #405 Track A2/B1. No protocol evolution past v3.82.0. Next phase: field application of CDS/CDR + handoff/memory-return testing.

## INDEX update

```
| 429 | #429 | 2026-05-24 | 3 | 1 | 1 | 1 | .cdd/unreleased/429/cdd-iteration.md |
```

Filed by δ (operator-side γ) on 2026-05-24.
