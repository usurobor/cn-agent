# Remote-runner receipt — backfill 3.82.0 binaries

**Doctrine:** docs/gamma/essays/BOX-AND-THE-RUNNER.md
**Skill:** src/packages/cnos.cdd/skills/cdd/delta/SKILL.md §8
**Cycle:** cnos#429 (AC1/AC2)

| Field | Value |
|-------|-------|
| **Who asked** | δ (operator), finishing cnos#429 after sigma terminal death |
| **What runs** | `.github/workflows/backfill-3.82.0-binaries.yml` — build 4 platform binaries at tag 3.82.0, verify (go test + Tier-1 kata), generate checksums.txt, upload as release assets via GitHub asset API |
| **Where** | GitHub Actions runners (ubuntu-latest, ubuntu-24.04-arm, macos-latest, macos-14); target = `usurobor/cnos` release `3.82.0` |
| **Authority** | `GITHUB_TOKEN`, `contents: write` scoped to publish job only; build/checkout use `contents: read` |
| **Evidence** | Self-delete commit `2bd35595` on main (publish completed); release asset count 2→16; `cn-linux-x64` downloaded, SHA `b8872062…bba72` matches checksums.txt; binary reports `Current version: 3.82.0` |
| **Who may accept** | Operator (δ) — accepted: AC1/AC2 satisfied. Caveat F3: `cn-macos-x64` is arm64-mislabeled (next-MCA) |

**Trigger:** push to main on the workflow's own path. **Lifecycle:** one-shot; self-deleted its own file after publish so the latent-execution authority closed.

**Body-preservation note:** assets were uploaded via the release-asset API (not softprops with `body_path`), deliberately, so the cycle/427 tightened release notes were not reverted to the pre-rewrite RELEASE.md present at tag commit fd1d654e.
