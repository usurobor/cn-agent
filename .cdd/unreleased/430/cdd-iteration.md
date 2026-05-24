# cdd-iteration — Cycle 430

**Cycle:** [cnos#430](https://github.com/usurobor/cnos/issues/430) — macos-x64 cross-compile fix (cnos#429 F3)
**Date:** 2026-05-24
**Mode:** release-hygiene; γ+α+β collapsed on δ
**Authority:** cnos.handoff/skills/handoff/receipt-stream/SKILL.md; cnos#429 F3 (origin)

## §0 Findings

`protocol_gap_count: 0` new. This cycle **resolves** cnos#429 F3 (no new gaps surfaced).

## §1 What changed

1. **release.yml (durable):** build matrix made arch-explicit (`goos`/`goarch`/`native`). `macos-x64` now cross-compiles `GOOS=darwin GOARCH=amd64` on the `macos-14` arm64 runner; native `go test` + Tier-1 kata skipped for that leg via `if: matrix.native` (amd64 can't exec on arm64). Root cause: `macos-latest` is now Apple Silicon, so the old `os: macos-latest` leg built native arm64 and mislabeled it x64.
2. **One-shot rebackfill workflow:** rebuilt all four 3.82.0 binaries with the corrected matrix, regenerated `checksums.txt`, re-uploaded all assets via the release-asset API (body preserved), self-deleted (`4b4a730b`).

## §2 Verification

- `file cn-macos-x64` → `Mach-O 64-bit x86_64 executable` (was arm64).
- `cn-macos-x64` SHA `b6e94f83…` now distinct from `cn-macos-arm64` SHA `8cac032a…`.
- `checksums.txt` regenerated; all four entries match the downloaded assets (`cn-linux-x64`, `cn-macos-x64`, `cn-macos-arm64` independently re-verified).
- release.yml YAML validated; matrix targets `[linux-x64, linux-arm64, macos-arm64, macos-x64]`.

## §3 Pause status

Re-affirmed. No #405 Track A2/B1. No protocol evolution past v3.82.0. Next phase: field application of CDS/CDR + handoff/memory-return.

## §4 Residual

- **AC3 of cnos#429** (release.yml clean on a real tag-push) remains unconfirmed in this environment — needs an operator PAT tag-push or UI "Run workflow" (F2 limitation). The corrected matrix is the thing that next tag will exercise.

## INDEX update

```
| 430 | #430 | 2026-05-24 | 0 | 2 | 0 | 0 | .cdd/unreleased/430/cdd-iteration.md |
```

Filed by δ (operator-side γ) on 2026-05-24.
