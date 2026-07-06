# self-coherence — cycle #608

**manifest:** sections = [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness]
**completed:** [Gap, Skills]

---

## §Gap

**Issue:** [cnos#608](https://github.com/usurobor/cnos/issues/608) — "cds-install Sub 1 (Cn=1 / PR 1): implement `cn repo install` — base installer (dispatch none)"

**Mode:** `design-and-build` (per the issue body's own `**Mode:**` field).
**cell_kind:** `implementation` (per the issue body's own `**cell_kind:**` field).

**Parent:** #607 (wave master tracker — "first-class CDS repo installer"). This
cycle is PR1 of the wave (base installer, dispatch none). PR2 (#609, renderer
generalization), PR3 (#610, dispatch install), PR4 (#611, bootstrap delegation)
are explicitly out of scope for this cycle — only #608 is dispatched now.

**Gap being closed:** cnos ships the hub/package path (`cn init`, `cn setup`,
`cn deps restore`, `cn doctor`, `cn build`, `cn update`) but has no first-class
"install CNOS/CDS into this repo" command. Onboarding an arbitrary Git
repository into the CDS method requires ~7 manual steps (install `cn`, fetch
the release index, hand-write `.cn/deps.json`, run `cn deps lock`, run
`cn deps restore`, edit `.gitignore`, commit) and `cn init` is the wrong tool
for this — it scaffolds a full agent-hub (`spec/SOUL.md`, `agent/`, `threads/`,
`state/`) into a `cn-<name>/` subdirectory, not a package-substrate install
into the current repo (cnos#606 C4).

**What closes it:** a new kernel command, `cn repo install`, that:
- runs in a plain `git init`-only repo before any `.cn/` hub or package state
  exists;
- resolves a cnos release (default: latest) or an explicit `--index`
  path/URL;
- writes deterministic `.cn/deps.json` (schema `cn.deps.v1`) pinning the
  default package set (`cnos.core`, `cnos.cdd`, `cnos.cds`) to exact,
  index-resolved versions;
- reuses the existing `internal/restore` substrate directly
  (`restore.GenerateLockFromIndex`, `restore.Restore`) to write
  `.cn/deps.lock.json` (schema `cn.lock.v2`, SHA-256-pinned) and restore
  packages under `.cn/vendor/packages/<name>/` (name-based, not
  `<name>@<version>/`);
- ensures the target repo's `.gitignore` excludes `.cn/vendor/`;
- supports `--release latest|<tag>`, `--index <path-or-url>`,
  `--packages <csv>`, `--dry-run`;
- gates autonomous dispatch install behind an explicit, currently-failing
  `--dispatch cds` flag (default `--dispatch none`) — base install never
  writes `.github/workflows/` and never requires a workflow/agent secret;
  dispatch install itself is #609/#610 territory, out of scope here.

**Design surface:** `docs/development/design/cn-repo-install-MOCKS.md`
(Mocks A, B, E1 — the base-install half of the wave's design doc; Mocks
C/D/F/G belong to later subs). This file did not exist on `origin/main` at
cycle start; per γ's scaffold it was pulled verbatim (operator-reviewed,
"Enacted" per #607) from `origin/claude/cds-install-guide-2ka54j` rather than
re-authored — landed in commit `f22d5a78` (pre-rebase) / `40a6706c`
(pre-rebase second push), now at `d783ba38`'s ancestor
`f22d5a78` → after the rebase onto `origin/main` documented below, the
current SHA for that commit is unchanged in content (rebase was a clean,
non-conflicting fast-forward-style replay — no file in this cycle's diff
overlapped the intervening `origin/main` commits).

**Base SHA (original, at branch creation):** `3dad64285026582a161549b8fd10108dd67a369e`
**Base SHA (current, after α's pre-review-gate rebase):** `80778d688e04e61c66d38f2bd5962fafb0729e95`
(`origin/main` advanced by two unrelated automation commits — board-map
regeneration + a sigma heartbeat log entry — during this cycle; α rebased
`cycle/608` onto the new tip per alpha/SKILL.md §2.6 row 1. No conflicts;
the two commits touch `docs/development/board/*` and `.cn-sigma/logs/*`,
disjoint from every path this cycle's diff touches.)

---

## §Skills

**Tier 1 (lifecycle):**
- `CDD.md` — canonical lifecycle/role contract.
- `cnos.cdd/skills/cdd/alpha/SKILL.md` — this role's execution surface (§2.1
  dispatch intake, §2.5 self-coherence, §2.6 pre-review gate, §2.7 request
  review, §3.6 implementation-contract constraint).
- `design/SKILL.md` not separately re-loaded as a standalone step: the design
  artifact (`cn-repo-install-MOCKS.md`) was already operator-reviewed/pinned
  per γ's scaffold; this cycle's "design" work was landing it verbatim, not
  authoring new design content requiring the design skill's judgment
  algorithm.
- `plan/SKILL.md` not separately invoked as a standalone artifact: γ's
  scaffold already carried a concrete "Surfaces α is expected to touch" table
  + AC oracle mapping + expected diff scope, which functioned as the plan.
  Sequencing was straightforward (land design doc → domain package → cli
  wrapper → registration → tests → docs → self-coherence) with no
  non-trivial ordering decisions requiring the plan skill's own algorithm.

**Tier 2 (`cnos.eng`):**
- `eng/go/SKILL.md` — load-bearing throughout. §2.18 ("Dispatch boundary:
  cli/ owns dispatch only") directly shaped the repoinstall/cli split;
  §2.17 (parse/read split, parallel-parser prohibition) governed reuse of
  `pkg.ParsePackageIndex`/`pkg.ParseManifest`/`pkg.ParseLockfile` rather than
  inventing new ad hoc JSON decoding; the determinism conventions (sorted
  map-key iteration, no map-order-dependent output) shaped
  `writeManifest`/`rewriteRelativeEntriesFromBase`.

**Tier 3 (issue-specific):** none loaded beyond `eng/go` — this is a
single-language (Go) kernel-command cycle with no CLI-ergonomics-package,
security-review, or other Tier-3 skill named by the issue or scaffold.

**Implementation contract (δ-pinned, per alpha/SKILL.md §3.6 — not
improvised):** all 7 axes were populated in γ's scaffold (Language: Go;
CLI integration target: kernel `cn` subcommand via noun-verb resolution;
Package scoping: `cli/cmd_repo_install.go` + new `internal/repoinstall/`
package, reusing `restore/`, `pkg/`, `binupdate/`; Existing-binary
disposition: coexist, `cn init`/`setup`/`deps` unchanged; Runtime
dependencies: none new; JSON/wire contract preservation: `cn.deps.v1` /
`cn.lock.v2` / `cn.package-index.v1` schemas unchanged, vendor path
name-based; Backward-compat invariant: `cn init`/`setup`/`deps` behavior
unchanged). No row was empty or TBD at dispatch time, so no γ/δ escalation
was needed before coding. Every diff hunk in this cycle maps to one or more
of these pinned rows — see §CDD Trace step 6 for the file-by-file mapping.
