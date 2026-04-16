# cnos Changelog

## Release Coherence Ledger

Each release is scored on two dimensions: **coherence quality** (TSC grades) and **architectural scope** (engineering level).

**TSC grades** use the [TSC](https://github.com/usurobor/tsc) triadic axes, applied as intuition-level letter grades (A+ to F) per the [CLP self-check](https://github.com/usurobor/tsc-practice/tree/main/clp):

- **C_Σ** — Aggregate coherence: `(s_α · s_β · s_γ)^(1/3)`. Do the three axes describe the same system?
- **α (PATTERN)** — Structural consistency and internal logic. Does repeated sampling yield stable structure?
- **β (RELATION)** — Alignment between pattern, relations, and process. Do the parts fit together?
- **γ (EXIT/PROCESS)** — Evolution stability and procedural explicitness. Does the system change consistently?

These are intuition-level ratings, not outputs from a running TSC engine (formal C_Σ ranges 0–1; ≥0.80 = PASS).

**Engineering level** per [ENGINEERING-LEVELS.md](docs/gamma/ENGINEERING-LEVELS.md), capped by CDD §9.1 cycle execution quality:

- **L4** — Pre-architecture: working prototype, no stable boundaries yet
- **L5** — Local correctness: fix works, follows patterns, no boundary change
- **L6** — System-safe: cross-surface coherence, failure modes handled
- **L7** — System-shaping: architecture boundary moved, class of future work eliminated

See [RELEASE-LEVEL-CLASSIFICATION.md](docs/gamma/essays/RELEASE-LEVEL-CLASSIFICATION.md) for the full historical analysis.

| Version | C_Σ | α | β | γ | Level | Coherence note |
|---------|-----|---|---|---|-------|----------------|
| 3.56.2 | A | A | A | A | L6 (diff: L7; cycle cap: L6) | **`cn deps lock` honors `deps.json`** (#250, PR #259). `restore.GenerateLockFromIndex` was iterating the full package index and writing every `(name, version)` pair into `.cn/deps.lock.json` — pins in `deps.json` had no effect and `cn deps restore` over-vendored every package the hub knew about. Fix: `pkg.ParseManifest` (pure) + `restore.ReadManifest` (IO wrapper, mirroring the §2.17 purity split) read the operator manifest; the generator then iterates `m.Packages` and resolves each pin via `idx.Lookup(name, version)`. Unresolved pins are collected into one explicit error naming every missing `name@version` rather than silently dropped. Missing `deps.json` now errors rather than succeeding with an index dump. Six AC-named tests cover single-pin / duplicate-free / restore-only-pinned / missing-manifest / missing-pin / default-profile — all six fail on the buggy code, pass on the fix. Go/OCaml lockfile-generation semantics now agree (`src/ocaml/cmd/cn_deps.ml` `lockfile_for_manifest` already did this correctly — #259 closes the discrepancy, not introduces one). **CI harness sync:** the kata Tier-1 `06-install.sh` post-condition and Tier-2 hub setup in `ci.yml` both relied on the bug; fix exposed two latent errors (Tier-1 pinned to binary `"dev"` version, Tier-2 used object-syntax `{"name":"version"}` which the parser rejects as `[]ManifestDep`). Both now derive real versions from `src/packages/*/cn.package.json` and use the canonical array schema. `cnos.kata` `lib.sh` `write_deps_json` emits the correct shape; `pkg_version_from_source` helper added and threaded into R3/R4 katas (no more hardcoded `3.54.0`). 1 review round, 0 findings, CI 7/7 green on `44d431a`. α: A — purity boundary preserved, error policy explicit, determinism preserved, Go/OCaml parity restored. β: A — independent build + race-test + manual revert-to-bug verification on the worktree; module-truth and sibling-fallback audits both clean. γ: A — dispatch clean, issue scoped precisely (semver/transitive explicitly out-of-scope), cycle closed in ≤1 round. **Scoring convergence (α↔β):** α close-out scored diff-shape L7 (silent-success fallback removed — malformed manifests and stale pins now surface at lock time rather than being masked into quiet over-installs) and cycle-process L6 (Go code L5-clean on first push; two pre-review CI iterations surfaced three bug-masked test-harness defects — Tier-1 weak post-condition, Tier-2 object-syntax schema mismatch, `cnos.kata` `write_deps_json` same schema mismatch — all fixed in the same PR). β originally scored L5 (CHANGELOG at release time), revised to L6 during assessment to converge with α; L7 diff-shape is credible but cycle-level is capped at L6 per §9.1 (some surface drift surfaced at CI rather than authoring time). Invariants: T-004 (source/artifact/installed explicit) tightened — lockfile now mechanically derives from manifest intent rather than ambient index state; T-002, T-003, INV-001, INV-003 preserved. §9.1 triggers: none fired at review time (review rounds 1, mechanical ratio 0/0, no loaded skill failed to prevent a review finding). **Skill gap surfaced in α close-out:** `eng/go` §2.12 schema-and-compatibility covers Go code handling manifests but does not scope to sibling shell test harnesses producing JSON for the same parsers — that gap allowed `cnos.kata/lib.sh`'s `write_deps_json` schema mismatch to survive to CI discovery rather than authoring-time audit. β concurs with α's proposed one-line addendum and **shipped it as immediate output** (§10.1): §2.12 now includes "sibling harnesses are contract surfaces too" with a mechanical grep form. Known debt: `cn setup`'s binary-version pin (#unfiled follow-up candidate) — in dev/CI environments without `-ldflags`, `cn deps lock` now errors explicitly rather than silently succeeding. Issue #250 scope intentionally excludes semver resolution and transitive dependencies; both remain as listed follow-ups if/when the manifest schema grows a `dependencies` field. |
| 3.56.1 | A | A- | A | A | L6 (cycle: L5) | **ContentClasses convergence** (#253, PR #258). `pkgbuild.ContentClasses` (filesystem discovery, 8 classes) and `pkg.FullPackageManifest.ContentClasses()` (JSON-field heuristic, 5 classes incl. a non-existent `providers`) disagreed on what counts as a content class. Converged to a single authoritative list: `pkg.ContentClasses` (pure-types package, exported, 8 entries in canonical §1.1 order) + `pkgbuild.FindContentClasses(pkgDir)` (shared filesystem predicate used by both `cn build --check` and `cn status`). Removed the JSON-field heuristic entirely: `FullPackageManifest.ContentClasses()` method, `Skills`/`Orchestrators`/`Extensions`/`Providers` fields, and the now-orphaned `SkillsJSON` type all deleted — presence on disk is the single authority (PACKAGE-SYSTEM.md §3). `providers` correctly removed as never-a-content-class (it's a runtime capability surface per POLYGLOT-PACKAGES-AND-PROVIDERS.md). Doctrine synced: PACKAGE-SYSTEM.md §1.1 adds `katas` row (8th class — completes the v3.55.0 drift where pkgbuild gained katas but doctrine didn't), §1.2 "directory tree copy" enumeration now includes katas (F1 narrowing round), §3 names `pkg.ContentClasses` and `pkgbuild.FindContentClasses` as the implementation, §4.1/§4.3 counts corrected 7→8, §6 history gains v3.55.0 (katas-added) and v3.56.1 (convergence) rows. 4 new tests covering empty/full/subset/cn-status-end-to-end paths; real packages still validate. 2 review rounds (target ≤2), 1 finding (F1: B mechanical — §1.2 sibling doc row stale after §1.1 update), fixed on-branch in `507e1e8`. α A- because the CDD §5.3 row 7a schema/shape audit missed the §1.1 ↔ §1.2 peer-enumeration pair (table column ↔ prose bullet restating the same column); α close-out proposes a concrete one-line row 7a addendum as MCA, pending γ triage. β A because the divergence was caught by review §2.2.8 authority-surface audit, narrowed cleanly, and review converged in 2 rounds. γ A because the dispatch chain ran cleanly and the cycle closed within target economics (both close-outs committed directly to main per §1.4 step 11). Invariants: T-005 (finite content classes) preserved-tightened — still 8, now single-sourced. §9.1 trigger: loaded skill failed to prevent a finding (row 7a scope gap); MCA candidate proposed and pending γ ship decision. Mechanical ratio 100% (1/1); total findings < 10 so no process issue filed per post-release §5.5. **Cycle level L5** (L6 cap — doc-internal coherence drift reached review); release diff level L6 (all shipped surfaces agree end-to-end). |
| 3.56.0 | A | A | A- | A- | L6 | **Git-style subcommand dispatch** (#254, PR #257). `cn <noun> <verb>` is now the primary command surface — `ResolveCommand` in `internal/cli/dispatch.go` resolves noun-verb pairs to flat commands, with flat-hyphenated form preserved as backward-compat fallback. `cn help <noun>` lists grouped verbs. Lookup order: noun-verb first, flat fallback second. 14 dispatch tests + 2 help tests, AC1–AC5 met. Design pre-converged in DESIGN-CONSTRAINTS §3.1 (1a7c362) — α cycle was pure MCA. 1 review round, 0 findings, β approved Round 1. α close-out: L6-clean, no §9.1 triggers. **Process note:** α close-out was committed on PR branch and destroyed by squash-merge; had to be recovered to main manually — exposed that `.cdd/` artifacts must be committed to main directly, not on PR branches. β did not produce post-release assessment (completed by σ). α: A (zero-finding cycle, design pre-converged, self-coherence applied at authoring time). β: A- (clean review, but failed to produce post-release assessment — closure gap recurring from 3.55.0). γ: A- (cycle dispatched correctly, but did not enforce closure artifact collection from β). |
| 3.55.0 | A | A- | A | A | L7 | **Kata framework separated from kata content; CDD triadic protocol operationalized** (#251 primary, PR #252; bundles #246/#247, #237/#248, and direct-to-main CDD skill patches). **L7 move:** `cnos.kata` 0.2.0 is now the kata framework — `cn kata-run`, `cn kata-list`, `cn kata-judge` — discovering katas across any installed package via the `<pkg>/katas/<id>/kata.md` convention (class parsed from `**Class:** <runtime\|method>`). `cnos.cdd.kata` 0.3.0 becomes content-only (no commands). The class of work "wire up a kata system for a new domain" is eliminated — a future `cnos.eng.kata` / `cnos.ops.kata` just drops a `katas/` dir. `pkgbuild.ContentClasses` expanded 7 → 8 to recognize `katas/` as a first-class package content directory. **Release hardening (#247):** the Tier 1 kata suite now runs inside the release workflow's build matrix on all 4 platforms between Package and Upload — a per-platform fresh-compile regression on the shipped binary is now structurally blocked from publishing. **Tier 2 runtime kata (#248):** `cnos.kata` ships R1 (command dispatch), R2 (author→build→install→dispatch round-trip in isolated tempdir), R3 (`cn doctor` catches chmod -x on installed entrypoint), R4 (`cn status` surfaces installed package + commands); `kata-tier2` CI job gates. **Process (direct-to-main):** identity format converged on `role@cdd.project` (was `<hub>-<role>@cnos.xyz`); CDD.md named canonical and SKILL.md demoted to loader (single source of truth); §1.4 roles rewritten as dyad + coordinator; §4.4 skill tiers formalized; α/β close-outs into γ cycle iteration; γ triage uses CAP; project-vs-agent MCI distinction (§11.12); `cn cdd-verify --version <v> [--triadic]` command added for artifact completeness (c24c75f3). **.cdd/ Phase 1:** first triadic protocol artifacts landed (`.cdd/releases/3.55.0/gamma/CLOSE-OUT.md`, `.cdd/releases/3.55.0/alpha/251.md`, and this release's beta assessment). **Cycle shape:** 2 reviewed PRs (2 rounds each); PR #248 had 1 mechanical finding (D stale-rebase), 1 structural-test bug (C), 1 filed-as-issue (A, became #250); PR #252 had 2C + 1B, all closed on Round 2 with zero new findings introduced. Mechanical ratio 33% on #251 cycle (F2 dual-source) → §9.1 fired; remediation: §2.5b check 6 sharpened to reject addition-only audits (MCA shipped in 89a1b2f). α: A- because two cycles surfaced mechanical findings a sharper §2.5b would have caught. β: A because both cycles closed in ≤2 rounds with explicit narrowing. γ: A because triage was issued, close-outs wired, and the MCA→shipped loop reduced cycle cap from L5 to L7 for the primary refactor. Invariants: T-005 preserved-tightened (finite set grew from 7 to 8); T-002, T-003, INV-001, INV-003 all preserved. |
| 3.54.0 | A | A | A | A | L6 | **Tier 1 bare-binary kata suite + CI gate** (#236, PR #239). Six executable kata scripts (`scripts/kata/0{1..6}-*.sh`) prove the `cn` binary works end-to-end before any package is installed: help → init → status → doctor → build → install. `run-all.sh` stops on first failure; CI job `kata-tier1` (needs `go`) gates merges. Three-tier kata catalog: Tier 1 bare binary (this release), Tier 2 runtime/package (#237), Tier 3 method/CDD (`cnos.cdd.kata`). Two Go preconditions: `cn help` always lists the 8 kernel commands with `(requires hub)` annotation when no hub; `cn doctor` gains three-way Status (Pass/Info/Fail → ✓/○/✗) so fresh-hub pending artifacts (deps.lock.json, runtime contract, origin remote, vendor before any lock) report informationally instead of failing. Dead `Registry.Available()` removed. 1 review round, 0 D/C findings, all architecture checks pass (§2.2.14.G improved — degraded paths now visible via `○` glyph). T-002 preserved (cmd_doctor/cmd_help remain thin wrappers, CI dispatch-boundary grep green). T-003 preserved (stdlib-only). γ+β collapse per CDD §1.4 (operator-authorized). |
| 3.53.0 | A | A | A | A | L5 | **cn status surfaces installed packages and command registry** (#233, PR #234). Go kernel Phase 4 completion: `cn status` reads cn.package.json manifests from vendor dirs (not dir names), displays name/version/content classes, shows command registry grouped by tier (kernel/repo-local/package) with source attribution. Version drift compares `engines.cnos` from manifest against running binary. `hubstatus` accepts pure `CommandInfo` data type from cli layer — no import cycle. 6 tests. 1 review round, 0 blockers (2 A-level), all architecture checks pass. T-002 preserved (cmd_status.go thin wrapper). T-004 preserved (installed state read from manifests). INV-003 preserved (commands/skills/orchestrators distinct in content class display). MVA Step 4 (#228) complete; #192 AC4 met. |
| 3.52.0 | A | A | A | A- | L6 | **Package and repo-local command discovery + dispatch** (#226, PR #231). Go kernel Phase 4: three source forms (kernel, repo-local, package) normalized into one registry, one dispatch. `internal/discover/` owns scanning + exec. `ExecCommand` dispatches via subprocess with `CN_HUB_PATH`/`CN_PACKAGE_ROOT` env vars. `cn help` shows tiered output with `[pkg-name]` attribution. `cn doctor` validates command integrity (missing/non-exec entrypoints, duplicates). Path confinement on manifest entrypoints. 13 new tests. 1 review round, 4 findings (0D, 1B mechanical, 3B/A judgment), all fixed on-branch. T-002 preserved (cmd_*.go thin wrappers, CI-enforced). INV-003 preserved (commands dispatch, no surface conflation). INV-004 preserved (kernel owns precedence). |
| 3.51.0 | A | A | A | A- | L5 | **Distribution pipeline proved end-to-end** (#227, PR #229). `src/packages/ → cn build → dist/ → cn deps restore → .cn/vendor/packages/` round-trip verified with `TestBuildRestoreRoundTrip`. 3 mismatches fixed: VendorPath aligned to BUILD-AND-DIST.md (version-less), local file restore for dev workflow, `cn deps lock` subcommand. Lockfile types consolidated on `pkg.Lockfile`/`pkg.LockedDep`. T-004 tightened. 2 review rounds, 5 findings all resolved. Known debt: #230 (version upgrade skip). |
| 3.50.0 | A+ | A+ | A+ | A+ | L7 | **Go kernel Phase 3 complete.** Slice D: `update` + `setup` commands (#212, PR #225). `cn update` fetches latest release, SHA-256 verifies, atomic installs. `cn update --check` for dry run. `cn setup` creates .cn/, agent/, default deps.json, .gitignore. All 8 kernel commands implemented: help, init, setup, deps, status, doctor, build, update. 63 Go tests. eng/design-principles skill (L6): Parnas/Martin/Liskov + cnos-specific rules. CDD review §2.2.14 architecture check (7 questions A-G) + §2.3.5 verdict gate — turns design principles into a mechanical review gate. BUILD-AND-DIST migration steps explicit. T-004 premature canonicalization fixed. I1 caught design-principles manifest drift — first real coherence CI catch on main. α A+ (all 8 commands, dispatch boundary CI-enforced). β A+ (design-principles skill + architecture review gate). γ A+ (full CDD invariants chain complete: author→handoff→review→close→CI). |
| 3.49.0 | A+ | A | A+ | A+ | L5 | cli/ extraction (#221, PR #222): dispatch boundary enforced. Domain logic extracted from cli/ to internal/doctor/, internal/hubinit/, internal/hubstatus/. FindIndexPath moved to internal/restore/. CI dispatch boundary check added — INVARIANTS.md T-002 is now mechanical. cmd_doctor 238→44 lines, cmd_init 131→27, cmd_status 76→27. INVARIANTS.md v1.2.0: title normalized, stale references cleaned. α A+ because the purity boundary is now machine-enforced, not just review-enforced. γ A+ because the CDD invariants chain (author→handoff→review→close) produced its first mechanicalization. |
| 3.48.0 | A | A+ | A+ | A+ | L7 | Go kernel Phase 3 Slice C: `build` command (#219, PR #220). Three modes: build (tarball + index + checksums), check (CI drift detection), clean. Pure/IO split in `pkgbuild` package. Architectural invariants doc (INVARIANTS.md): 6 active invariants, 5 transition constraints, 3 process constraints — validated each CDD cycle. CDD lifecycle integration: project invariants loaded at §2.4 (pre-code), §2.5a item 5 (handoff), §2.2.13 (review), §6a (post-release). eng/go §2.17 (Parse/Read purity boundary) + §2.18 (cli/ dispatch boundary) — both derive from INVARIANTS.md. Build-and-Distribution design doc: independent package versioning, src/packages/ → dist/packages/ → .cn/vendor/. β A+ because the architectural constraint surface is now explicit and validated. γ A+ because the CDD self-learning loop now covers invariant preservation end-to-end. |
| 3.47.0 | A | A | A+ | A | L6 | Go kernel Phase 3 Slice B: `doctor` command (#212, PR #217). Hub health validation across 5 categories: prerequisites (git, curl), hub structure (.cn/config, SOUL.md, peers), package system (deps.json, deps.lock, installed packages + version drift), runtime contract (v2 four-layer validation), git remote. 30 Go tests. 2 review findings (F1: optional file severity, F2: test assertion gap), both fixed on-branch per §7.0. Architecture design docs: git-cn package design (#218), package restructuring target (#186), capabilities field added to cn.extension.v1 (no new manifest type). β A+ because the design frontier is now fully explicit: polyglot packages, provider contract, package restructuring target, git-cn package shape. |
| 3.46.0 | A | A | A+ | A | L6 | Go kernel Phase 3 Slice A: `init` + `status` commands (#212, PR #213). Hub lifecycle pair — create hubs with full structure + inspect with version drift detection. Input validation confinement fix (#214, PR #215): bare-name validation rejects path traversal. Two architecture design docs: POLYGLOT-PACKAGES-AND-PROVIDERS.md (Go kernel + language-agnostic packages + command/provider contract split) and PROVIDER-CONTRACT-v1.md (subprocess protocol: spawn/handshake/describe/health/execute/shutdown, NDJSON over stdio, kernel-owned policy). 26 Go tests. eng/go stdlib rule. CI OPAMRETRIES fix. β A+ because design docs make the polyglot boundary explicit — provider contract is the missing concrete layer beneath runtime-extensions. |
| 3.45.0 | A | A+ | A | A | L7 | First runnable Go binary: `go build ./cmd/cn && ./cn deps restore`. CLI entrypoint with modular command dispatch per GO-KERNEL-COMMANDS.md v1.2 (#209, PR #210). `CommandSpec` + `Invocation` + `Command` interface + `Registry` with tier precedence and `Available(hasHub)`. `deps restore` and `help` commands. Hub discovery walks up from cwd. 20 Go tests. CDD §7.0 gate: all review findings must be resolved before merge — validated on first use (PR #211 fixed F2+F3 before merge). |
| 3.44.0 | A | A+ | A | A | L6 | Go purity boundary fixed: `internal/pkg/` is now truly pure (zero `os` imports). `Parse*` (pure, `[]byte`) vs `Read*` (IO, path) mirrors `src/lib/` vs `src/cmd/`. HTTP timeout added (300s, matching OCaml curl flags). Directory traversal check hardened with separator suffix. OCaml skill v2 shipped (L6 parity with Go/TypeScript). Two architecture design docs: GO-KERNEL-COMMANDS (the kernel is a package — `cn.package.v1`-compatible manifest, `CommandSpec` runtime descriptor, 3-tier precedence) and HYBRID-LLM-ROUTING (body-plane routing matrix, local/remote/escalation, receipts). β A+ because the Go module now mirrors the OCaml purity split exactly, and the command architecture converged through 3 review iterations to a clean design. |
| 3.43.0 | A | A | A | A | L7 | First Go code in cnos. `go/` subtree implements `cn deps restore` in Go as an additive kernel slice alongside OCaml (#205, #192 Phase 1). Pure types in `internal/pkg/` mirror `src/lib/cn_package.ml`; IO restore path in `internal/restore/` mirrors `src/cmd/cn_deps.ml`. 13 tests (7 pkg + 6 restore), stdlib-only, zero external deps. Go CI workflow added. 1 review round, 0 blockers, 4 non-blocking findings (F2: HTTP timeout, F3: traversal hardening, F4: purity boundary — all follow-up). First-push CI green. New system boundary: Go module coexists with OCaml. |
| 3.42.0 | A | A | A+ | A | L5 | Doc authority convergence + CDD self-learning. Beta package-system doc retired (#180) — the artifact-first model is now the sole authority, no competing Git-native narrative. CDD post-release skill gains mandatory hub memory writes (Step 7): daily reflection + adhoc thread update required before cycle close, with pre-publish gate checks. Proven by v3.41.0 compaction gap where assessment was written but hub memory was not. Release skill gains §2.6a branch cleanup: delete merged branches after tag+push. 13 stale branches cleaned. #192 Go kernel rewrite scoped to narrow 4-phase extraction with explicit non-goals. |
| 3.41.0 | A | A | A+ | A | L7 | Move 2 complete: pure-model gravity into `src/lib/`. Fourth and final slice ships `cn_frontmatter.ml` (12 pure surface items: frontmatter parser, activation validation types, manifest skill-id reader). `src/lib/` now holds 5 pure modules (cn_json, cn_package, cn_contract, cn_workflow_ir, cn_frontmatter) — every pure type and parser in the codebase. `src/cmd/` holds only IO. The Go kernel rewrite (#192) has its boundary contract. First cycle with zero review findings — §2.5b checks 7+8 (#198) validated. CDD §1.4 Roles formalized: two-agent minimum for substantial cycles, reviewer is default releaser, Role column added to lifecycle table. TypeScript skill converged (schema-backed boundaries, branded primitives, explicit error policy, idempotence/receipts). |
| 3.40.0 | A | A | A | A | L6 | Pure-model gravity: Move 2 slices 2+3 (#182). Runtime contract types extracted into `src/lib/cn_contract.ml` (11 types + `zone_to_string` + transitive `activation_entry`). Workflow IR types, parser, validator, and helpers extracted into `src/lib/cn_workflow_ir.ml` (6 types + 10 pure functions — largest `src/lib/` module). `src/lib/` now holds 4 modules (cn_json, cn_package, cn_contract, cn_workflow_ir) covering the two largest pure surfaces. Type-equality re-export pattern validated at scale (6-variant + 7-variant types). Option-(b) principle established: pure types stay in `src/cmd/` when no `src/lib/` consumer exists. CDD §2.5b gate extended with checks 7 (library-name uniqueness) + 8 (CI-green-before-review). Go skill (`eng/go`) converged through 4 iterations — runtime/kernel craft + cnos-specific systems rules. 2 review rounds each for slices 2+3. Move 2 is 3/4 complete; 1 remaining (activation evaluator). |
| 3.38.0 | A | A | A | A | L5 | Pure-model gravity: first slice of Move 2 (#182). Package manifest, lockfile, and index types extracted from `src/cmd/cn_deps.ml` into `src/lib/cn_package.ml`. 6 types + 7 pure helpers + JSON round-trips. Cn_deps re-exports via OCaml type-equality (`type t = Cn_package.t = { ... }`) — zero caller migration. Purity discipline enforced: only stdlib + Cn_json imports. 11 ppx_expect tests (round-trips, parse, lookup, is_first_party). Dune lib comment documents the no-IO constraint. CORE-REFACTOR.md §7 status block added (3 remaining slices). 1 review round, 0 findings. Vision doc: AGENT-NETWORK.md committed (docs/alpha/vision/). |
| 3.37.0 | A- | A | A | B+ | L6 | Command pipeline symmetry (#184). Commands become the 7th content class in `cn build`: `src/agent/commands/<id>/` → `cn build` → `packages/<name>/commands/<id>/` → `cn deps restore` → vendor → runtime discovery. `daily`, `weekly`, `save` migrated from built-in OCaml to shell-script package commands in `cnos.core`. Built-in command set shrinks toward the bootstrap kernel (help, init, setup, deps, build, doctor, status, update). `cn_lib.ml` command variant, `cn.ml` dispatch, and `cn_gtd.ml` run_daily/run_weekly all cleaned. Dual-field manifest: `sources.commands` (string array for build pipeline) + top-level `commands` (object map for runtime metadata). `cn_command.ml` discovery/validation reads top-level `commands`; `cn_build.ml` reads `sources.commands`. `chmod +x` on `cn-*` entrypoints post-copy. 6 new build tests (AC1: copy+exec, clean, drift detection). PACKAGE-SYSTEM.md §1.1 updated. 3 review rounds (over ≤2 target — test fixture used old manifest schema, same root cause across two files, fix only caught one in R1), 1 finding (100% mechanical). Move 1 of #182 core refactor. |
| 3.36.0 | A | A | A | A- | L7 | Orchestrator IR runtime (#174). New execution surface: `cn.orchestrator.v1` mechanical workflow engine (`Cn_workflow`) — parse, validate, discover, execute. 6 step kinds: `op` (dispatches to typed ops via `Cn_executor.execute_op`), `if`/`match` (deterministic branching on bound scalars), `return`/`fail` (terminal), `llm` (guarded stub — parsed, validated, permission-checked, deferred to next cycle). Permission manifest enforced at both validator and executor (defence-in-depth). Execution trace via `Cn_trace.gemit` at every step boundary. `orchestrators/` added as 8th package content class in `cn_build` (source_decl, parse_sources, build/check/clean). `cn doctor` validates orchestrator manifests (load failures + schema issues, fail-stop). `build_orchestrator_registry` rewritten to consume `Cn_workflow.discover`, eliminating the #173 dueling-schema (`[{name, trigger_kinds}]` inline → `["id"]` string array per ORCHESTRATORS.md §9). `daily-review` shipped as first real orchestrator (schema-clean, not end-to-end runnable until `llm` step is wired). PACKAGE-SYSTEM.md §1.1 updated. Sibling audit: 5 pre-existing bare catches cleaned in `cn_build.ml`. 21 cn_workflow_test + 21 cn_runtime_contract_test (3 replaced). 2 review rounds, 5 findings R1 (2D mechanical: rebase artifact + stale README refs; 1C: `inputs` dead code; 2B: unvalidated field + missing `match` execution test), all closed R2. |
| 3.35.0 | A- | A | A- | B+ | L6 | Runtime contract: activation index + command/orchestrator registries (#173). Cognition layer now exposes `activation_index.skills` with declarative triggers from SKILL.md frontmatter. Body layer adds `commands` and `orchestrators` registries. Markdown/JSON parity enforced — markdown projection mirrors JSON schema shape. Silent error suppression eliminated: all manifest parse errors, malformed entries, and leaf-value coercions logged with package context in both `cn_runtime_contract.ml` and `cn_activation.ml`. Sibling audit applied across both modules. Design docs: ORCHESTRATORS.md v0.1.0 CTB connection (#172), #170 dependency graph cleaned. 21 expect-tests in runtime contract, 11 in activation. 3 review rounds on retro-findings PR. Process gap: PR #176 merged without CDD review — corrected via retro-review + findings PR #177. |
| 3.34.0 | A- | A- | A- | B+ | L7 | Package artifact distribution + commands content class (#167). Replaces git-based package restore with HTTP tarball fetch from package index. Lockfile v2: name+version+sha256 only — no source/rev/subdir. `commands` added as 7th package content class. `cn_command.ml`: discovery + dispatch (built-in > repo-local > package). `cn_help.ml`: lists external commands. `cn_doctor.ml`: validates command integrity. `scripts/build-packages.sh` + `packages/index.json` for release workflow. CDD §2.5a: delegated implementation handoff protocol + self-verification gate. CDD.md §5.3: step 6f (delegated handoff) added to artifact manifest. Architecture-evolution skill §5: five L7 diagnostic patterns. Also includes: #155 vendor offline-first + manual recognition, #161 self-update checksum verification, #146 review findings, review PRE-GATE check, release script. Design: PACKAGE-ARTIFACTS.md v0.5.0, ORCHESTRATORS.md v0.1.0 (#170 follow-up). 2 review rounds, 3 findings R1 (bare catches, missing tests, HTTP restore tests), all addressed R2. |
| 3.33.0 | A- | A | A- | A- | L6 | Harden installer: checksum verification, redirect detection, curated notes (#158). Version detection replaced GitHub API JSON parsing with HTTP redirect (`/releases/latest` Location header) — no jq dependency. Release workflow generates `checksums.txt` (SHA-256 per artifact), installer downloads and verifies (sha256sum/shasum, graceful degradation). RELEASE.md curated for v3.33.0. 1 review round, 1 minor finding (grep -F defensive hardening, non-blocking). |
| 3.32.0 | A- | A | A | A- | L6 | Systematic legacy fallback path audit (#152). 29 findings across touched files: 7 removed (dead deprecated functions + flat asset compat), 15 converted (silent `with _ ->` swallows to logged warnings), 7 kept (fail-closed patterns with justification). `run_inbound`, `feed_next_input`, `wake_agent` deleted (dead code since v3.27). All flat hub-path compat helpers deleted (package namespace only since v3.25). Zero bare `with _ ->` in any touched file. 5 review rounds, 7 findings (4 judgment, 3 mechanical/environmental). Gamma improved from B+ (13 rounds in v3.31.0) to A- (5 rounds). Post-release assessment revised gamma: B+ -> A-. |
| 3.31.0 | A- | A | A | B+ | L7 | Canonical packet transport replaces diff-based inbox materialization (#150). `cn_packet.ml`: envelope schema (`cn.packet.v1`), 9-step validation pipeline, dedup/equivocation index, `refs/cn/msg/{sender}/{msg_id}` namespace. Legacy `materialize_branch` + orphan rejection pipeline fully deleted. Send-side creates root commits under packet refs. 16 ppx_expect tests, 8 cram tests rewritten. 13 review rounds, 22 findings (79% mechanical -- no pre-push test sweep). Silent content swap structurally eliminated. Post-release assessment revised gamma: A -> B+. |
| 3.29.1 | A | A | A | A | L6 | Patch: `resolve_bin_path` uses `Unix.readlink` directly instead of spawning shell child (#148). Fixes false version drift + daemon crash loop on Pi. Release skill updated: §2.5 RELEASE.md format, §3.7 mechanical gate. |
| 3.29.0 | A | A | A | A | L6 | Remove hardcoded paths (#146): `resolve_bin_path` derives binary location at runtime (`$CN_BIN` → `/proc/self/exe` → `readlink Sys.executable_name` → fallback), `resolve_repo` derives GitHub repo from build-time `cn.json` extraction (`$CN_REPO` override). `cn_repo_info.ml` dune-generated module. `cn_deps.ml` default source derived from `Cn_lib.cnos_repo`. 7 ppx_expect tests. 2 review rounds, 5 findings R1 (1B/4C), all addressed R2. Build fails if `cn.json` lacks GitHub URL (no silent fallback). |
| 3.28.0 | A | A | A | A- | L6 | Orphan rejection loop fix (#144): deterministic rejection filenames (`rejection_filename` keyed on peer+branch replaces timestamped `make_thread_filename`), `is_already_rejected` dedup across outbox+sent, `get_branch_author` removed (fetched peer identity authoritative), self-send guard in `send_thread`, `from:` field in rejection envelope. 9 ppx_expect tests. Inbox skill v3: meta-skill alignment (artifact_class, kata_surface, governing_question), Define/Unfold/Rules/Verify structure, two embedded katas, pull-only transport model documented with contrastive examples. 1 review round, 1 D-level finding (stale clone wording), fixed same session. |
| 3.27.1 | A | A | A | A | L5 | Remove cron installation management (#138): `update_cron` and `cn-cron` deleted, docs rewritten (AUTOMATION.md daemon-first, README systemd prerequisite, CLI.md cron label removed). Oneshot mode preserved. CDD traceability normalization: lifecycle step numbers (0–13) replace §2.x in PR template and review skill (#143). 2 PRs, 1 review round on #139 (scope narrowing). |
| 3.27.0 | A- | A | A- | B+ | L6 | CDD traceability + skill-loading bridge: artifact manifest (§5.3), execution trace (§5.4), review enforcement (§2.0.8), release trace update, post-release closeout (§6). Active-skill matrix by work shape × level in eng/README. CDD §2.4 binds level labels to concrete skills. Writing skill rewrite. CAP skill v2 (UIE-based). Meta-skill rewrite (artifact classification, kata surface). 18 commits, 5 skill rewrites, 2 new CDD sections. Process: 3 direct-to-main commits retro-closed via #137; §3.7 added as preventive rule. |
| 3.26.1 | A | A | A | A | L5 | Docs: OPERATOR.md (day-2 operations manual), post-release assessment v3.26.0, CLI/troubleshooting updated for cn logs, hub name sanitization across 19 files (strip deployment-specific agent names from public docs). |
| 3.26.0 | A- | A | A- | B+ | L7 | Unified log stream + cn logs CLI (#74 Phase 1): `cn_ulog.ml` append-only JSONL writer (schema cn.ulog.v1, 5 event kinds), `cn_logs.ml` CLI with --since/--msg/--errors/--json/--kind filters, correlation via msg_id across all entries. Runtime emits at 9 points (message lifecycle + silent drops + poll errors). Chunk-accumulation read path preserves chronological order across day boundaries. 21 tests (13 ulog + 8 logs). 6 review rounds across 2 reviewers — multi-day ordering bug (independent) and silent message drops (Σ) both found and fixed. §2.2.1a extended: new data surfaces require write AND read path verification. |
| 3.25.0 | A | A | A | A- | L7 | Structural self-knowledge interception (#64): `Contract_redirect` as first-class receipt status, `check_self_knowledge_path` interceptor in executor (after sandbox, before I/O), membrane covers all 5 observe surfaces — fs_read (redirect), fs_list child filtering, fs_glob result filtering, git_grep/git_observe exclusions. RC authority declaration references structural enforcement. 10 expect tests. 2 review rounds, 3 D-level membrane holes found R1 (fs_list children, fs_glob, git_grep), all closed R2. Last P0 (#64) closed. |
| 3.24.0 | A | A | A | A- | L7 | Template distribution via package system (#119): templates as 6th content class in package system, `read_template` with Result-typed 3-case resolution (Ok/not-installed/not-found), `run_init` reads from installed cnos.core after `setup_assets`, `run_setup` populates missing spec/ from templates without overwriting operator content. PACKAGE-SYSTEM.md architecture doc (content class taxonomy, pipeline, explicit-vs-generic rationale). 11 new tests (7 build + 4 e2e regression). Unicode hygiene added to review skill/checklist. 1 review round, 1 mechanical finding (stale cross-ref). All 5 ACs met. |
| 3.23.0 | A- | A- | A- | A | L7 | Pre-push gate (#117): `scripts/pre-push.sh` runs dune build, dune runtest, cn build --check, VERSION parity against origin/main. CDD §9.1 cycle iteration with structured template, triggers, cycle level (L5/L6/L7), closure gate, skill patching as immediate output. ENGINEERING-LEVELS.md cross-referenced. #64 P0 closed (all 3 root causes verified). SOUL §2.2 MCA-first on error correction. CDD §4.4 "load means read the file." ORIGIN.md genesis narrative. Squash merge enabled. |
| v3.22.0 | A- | A- | A | A- | L6 | Daemon version-drift detection (#110): `version_drift_check_once` as 8th maintenance primitive, `check_binary_version_drift` detects external binary replacement via `--version` comparison, triggers `re_exec` with shared `is_idle` drain guard. 1 review round, 0% mechanical ratio. All 4 #110 ACs met. |
| v3.20.0 | B+ | A- | A- | B | L6 | Runtime Extensions e2e (#73): host command resolution maps bare manifest names to installed paths, returns Result with existence/permission/traversal validation. Extension dispatch pipeline validated end-to-end with real cnos.net.http host binary. 3 review rounds, 0% mechanical ratio. All 7 #73 ACs met. MCI freeze lifted (8 releases). |
| v3.19.0 | A- | A | A- | B+ | L5 | Package System Substrate AC5-AC7 (#113): integrity hash (md5 tree hash) in lock generation + restore verification, doctor validates full desired→resolved→installed chain (metadata, integrity, stale installs), Runtime Contract package truth expanded with source/rev/integrity. 1 review round, 33% mechanical ratio. #113 AC3-AC7 shipped across v3.18.0+v3.19.0. |
| v3.18.0 | B+ | A- | B+ | B- | L5 | Package System Substrate AC3+AC4 (#113): path-consistent restore (copy_tree replaces 5-category hardcoded list), honest third-party rejection (lockfile_for_manifest returns Result), stale cnos.pm cleanup, package manifests synced with skill reorg. 3 review rounds, 57% mechanical ratio, 1 superseded PR. Partial: AC5-AC8 deferred. |
| v3.17.0 | B+ | A | B+ | C+ | L5 | Runtime Extensions Phase 1 (#73): open op registry replaces closed built-in vocabulary, subprocess host protocol, manifest-driven discovery, extension lifecycle (discovered→compatible→enabled/disabled/rejected), Runtime Contract integration (cognition + body), traceability (extension.* events), doctor health checks, first reference extension cnos.net.http. 5 review rounds, 53% mechanical ratio. Partial: build integration, policy intersection, host binary (Phase 2). TSC corrected by post-release assessment. |
| v3.16.2 | A | A | A | A- | L6 | Two-membrane projection integrity (#106): presentation membrane hardened (shared xml_prefixes, block-level + inline stripping, matched-tag close), self-knowledge membrane (anti-probe in Runtime Contract). 5 review rounds, 12 D-level findings resolved. Configuration mode spec + templates shipped. Review skill converged (§2.1.3 mechanical scan, §2.2.1a input enumeration). |
| v3.16.1 | A | A | A | A | L6 | Daemon retry limit + dead-letter (#28): 4xx fail-fast, exponential backoff (1s/2s/4s), offset advancement after dead-letter. Sustainability surface added. |
| v3.16.0 | A | A | A | A- | L5 | End-to-end self-update (#37): same-version patch detection, ARM release matrix, bare version tag trigger, target_commitish SHA validation. CDD OVERVIEW.md. AGILE-PROCESS.md deleted. 3 review rounds (target ≤2), 63% mechanical ratio — no build step before review. |
| v3.15.2 | A | A | A | A | L6 | Empty Telegram filter (#29), CDD v3.15.0 canonical rewrite (authority split resolved), review skill hardened (4 new checks from PR #103 comparison). |
| v3.15.1 | A | A | A | A- | L6 | Fix #22 review blockers: re-exec after binary update (no stale in-process version), stamp-versions.sh derives manifests from VERSION, cn build --check + cn release gate version consistency, unified truth read in update_runtime. Process: reviewed before merge, CI green. |
| v3.15.0 | A- | A | A- | B+ | L5 | Version coherence (#22): VERSION→dune→cn_version.ml chain is sound (α A). β regressed: I1 CI check failed at merge, claim/reality gap on AC3+AC10. γ regressed: 0 review rounds (self-merged ~26min), merged with red CI, 3 fix commits post-implementation — first release after §11.11 review metrics bypassed review entirely. Manifests not yet stamped from VERSION (validate-only, not derive). Update path uses stale in-process version. Core architecture correct; contract partially closed. |
| v3.14.7 | A | A | A | A | L7 | Reduce review round-trips (#97): branch pre-flight validation (§1.5), scope enumeration (§4.0), review quality metrics (§11.11), process debt integration (§11.12), encoding lag type column (§11.6), cross-ref validation in review skill (§2.2.5), finding taxonomy (§5.1). Post-release template updated with review quality section. |
| v3.14.6 | A | A | A | A | L7 | Retroactive epoch assessments (#85): v3.12.0–v3.12.2 and v3.13.0–v3.14.5. CDD §9.11 release gate: previous release must have assessment before tagging. |
| v3.14.5 | A | A | A | A | L5 | Organize gamma docs (#91): 6 gamma root files moved into 3 bundles (cdd, rules, essays). Bundle READMEs, frozen snapshots, 13+ cross-refs updated. Completes docs reorg trilogy: alpha→beta→gamma. |
| v3.14.4 | A | A | A | A | L5 | Organize beta docs (#89): 7 beta root files moved into 4 thematic bundles (architecture, governance, lineage, schema). Bundle READMEs, frozen snapshots, cross-refs updated. CDD §5.1 freeze semantics note. |
| v3.14.3 | A | A | A | A | L5 | Organize alpha docs (#86): all 18 root-level specs moved into 8 thematic bundles (doctrine, protocol, agent-runtime, runtime-extensions, cognitive-substrate, cli, security, ctb). Bundle READMEs created. Cross-references updated across 20+ files. Freeze policy updated: path-only corrections allowed. |
| v3.14.2 | A | A | A | A | L5 | Alpha bundle migration (#81): legacy design docs and plans moved into CDD bundle structure. Version snapshot dirs with manifests. Navigation surfaces updated. |
| v3.14.1 | A | A | A | A | L6 | CDD post-release tightening (#78): encoding lag table mandatory, concrete next-MCA commitment, MCI freeze triggers. First operator troubleshooting guide. LINEAGE.md taxonomy cleanup. |
| v3.14.0 | A | A | A- | A | L7 | Runtime Contract v2 (#62): vertical self-model (identity, cognition, body, medium). Zone classification for all paths. All architecture docs updated. Doctor structural validation. |
| v3.13.0 | A- | A- | A- | A- | L7 | Docs governance (#75): CDD pipeline with per-step artifacts, self-coherence report format, single cnos version lineage, feature bundles, frozen snapshots, bootstrap-first rule. CDD skill synced. |
| v3.12.1 | A+ | A+ | A+ | A+ | L6 | Daemon boot log declares config sources (#61): version, hub, profile, model, secrets provenance, peers. Type-safe secret_source. |
| v3.12.0 | A+ | A+ | A+ | A+ | L7 | Wake-up self-model contract (#56): Runtime Contract replaces overloaded capabilities block. CDD/review skill hardening (#57): issue AC gates, multi-format parity. |
| v3.11.0 | A+ | A+ | A+ | A | L6 | N-pass merge + misplaced ops correction (#51). Structured output reverted (needs rework). |
| v3.10.0 | A+ | A+ | A+ | A+ | L7 | N-pass bind loop (#50): effect pass no longer terminal, generic pass architecture, processing indicators. |
| v3.9.3 | A+ | A+ | A+ | A+ | L6 | Anti-confabulation (#49): explicit op result signals, denial reason surfacing, doctrine. |
| v3.9.2 | A+ | A+ | A+ | A+ | L6 | Block Pass A projection when Pass B fails (#47). CURIOSITY mindset. MIC/MICA/CDD verbs in glossary. |
| v3.9.1 | A+ | A+ | A+ | A+ | L5 | Fix Cn_shell.execute never called (#46), sync DUR skills to packages, add git_stage to protocol contract. |
| v3.9.0 | A+ | A+ | A+ | A+ | L7 | Two-pass wiring (#41), COGNITIVE-SUBSTRATE spec, DUR skill contract, reflect + adhoc-thread + review skills cohered. |
| v3.8.0 | A+ | A+ | A+ | A+ | L6 | Syscall surface coherence: fs_glob, git_stage, fs_read chunking, observe exclusion symmetry, CLI injection hardening. Review-driven. |
| v3.7.3 | A | A+ | A+ | A | L6 | Agent output discipline: ops-in-body detection, peer awareness, coordination op examples, conditional MCA review. |
| v3.7.2 | A | A+ | A+ | A- | L6 | Trace gap closure + boot drain fix + skill hardening. 5/7 trace gaps closed; peer-only heartbeat and outbox structure remain. |
| v3.7.0 | A | A+ | A+ | A- | L7 | Scheduler unification: one loop, two schedulers (oneshot/daemon), 7-primitive maintenance engine. Boot drain gap found post-merge. |
| v3.6.0 | A+ | A+ | A+ | A+ | L7 | Output Plane Separation: sink-safe rendering, typed op output, two-pass execution, CDD skill. |
| v3.5.1 | A+ | A+ | A+ | A+ | L7 | TRACEABILITY: structured event stream, state projections (ready/runtime/coherence), boot sequence telemetry, CDD design doc. |
| v3.5.0 | A+ | A+ | A+ | A+ | L7 | Unified package model + CAA + FOUNDATIONS. Everything cognitive is a package. Doctrinal capstone. Architecture spec. |
| v3.4.0 | A+ | A+ | A+ | A+ | L7 | CAR: three-layer cognitive asset resolver + package system. Fresh hubs wake with full substrate. Git-native deps, lockfile, hub-local overrides. |
| v3.3.1 | A+ | A+ | A+ | A+ | L6 | Agent instruction alignment: canonical ops examples in capabilities block, stale path fixes, output discipline. Prevents hallucinated tool syntaxes. |
| v3.3.0 | A+ | A+ | A+ | A+ | L7 | CN Shell: typed ops, two-pass execution, path sandbox, crash recovery. Pure-pipe preserved — ops are post-call, governed, receipted. Zero new runtime deps. |
| v3.2.0 | A+ | A+ | A+ | A+ | L7 | Structured LLM schema: system blocks with cache control + real multi-turn messages. Mindsets in context packer. Role-weighted skill scoring. Setup installer design. |
| v3.0.0 | A+ | A+ | A+ | A+ | L7 | Native agent runtime. OpenClaw removed. Pure-pipe: LLM = `string → string`, cn = all effects. Zero runtime deps. |
| v2.4.0 | A+ | A+ | A+ | A+ | L7 | Typed FSM protocol. All 4 state machines (sender, receiver, thread, actor) enforced at compile time. |
| v2.3.x | A+ | A+ | A+ | A | L7 | Native OCaml binary, 10-module refactor. No more Node.js dependency. |
| v2.2.0 | A+ | A+ | A+ | A+ | L7 | First hash consensus. Actor model complete: 5-min cron, input/output protocol, bidirectional messaging, verified sync. |
| v2.1.x | A+ | A+ | A+ | A | L6 | Actor model iterations: cn sync/process/queue, auto-commit, wake mechanism fixes. |
| v2.0.0 | A+ | A+ | A+ | A+ | L7 | Everything through cn. CLI v0.1, UX-CLI spec, SYSTEM.md, cn_actions library. Paradigm shift: agent purity enforced. |
| v1.8.0 | A+ | A+ | A | A+ | L7 | Agent purity (agent=brain, cn=body). CN Protocol, skills/eng/, ship/audit/adhoc-thread skills, AGILE-PROCESS, THREADS-UNIFIED. |
| v1.7.0 | A | A | A | A | L7 | Actor model + inbox tool. GTD triage, RCA process, docs/design/ reorg. Erlang-inspired: your repo = your mailbox. |
| v1.6.0 | A− | A− | A− | A− | L7 | Agent coordination: threads/, peer, peer-sync, HANDSHAKE.md, CA loops. First git-CN handshake. |
| v1.5.0 | B+ | A− | A− | B+ | L6 | Robust CLI: rerunnable setup, safe attach, better preflights. |
| v1.4.0 | B+ | A− | A− | B+ | L6 | Repo-quality hardening (CLI tests, input safety, docs aligned). |
| v1.3.2 | B+ | A− | B+ | B+ | L5 | CLI preflights git+gh; same hub/symlink design. |
| v1.3.1 | B+ | A− | B+ | B+ | L5 | Internal tweaks between tags. |
| v1.3.0 | B+ | A− | B+ | B+ | L6 | CLI creates hub + symlinks; self-cohere wires. |
| v1.2.1 | B+ | A− | B+ | B+ | L5 | CLI cue + onboarding tweaks. |
| v1.2.0 | B+ | A− | B+ | B+ | L5 | Audit + restructure; mindsets as dimensions. |
| v1.1.0 | B | B+ | B | B | L5 | Template layout; git-CN naming; CLI added. |
| v1.0.0 | B− | B− | C+ | B− | L5 | First public template; git-CN hub + self-cohere. |
| v0.1.0 | C− | C | C− | D+ | L4 | Moltbook-coupled prototype with SQLite. |

---

## 3.56.2 (2026-04-16)

### Fixed

- **`cn deps lock` honors `.cn/deps.json`** (#250, PR #259). Before: `restore.GenerateLockFromIndex` iterated every `(name, version)` pair in the package index, ignoring the manifest entirely. A hub pinning `cnos.core@3.54.0` got a lockfile with every package in the index at every version, and `cn deps restore` over-vendored. After: the generator reads `.cn/deps.json` first, then resolves each pin via `idx.Lookup(name, version)`. Unresolved pins fail with an explicit error listing every missing `name@version`; missing `deps.json` fails with a clear message instead of succeeding with an index dump.

### Added

- **`pkg.ParseManifest([]byte) (*Manifest, error)`** (`src/go/internal/pkg/pkg.go`): pure parser for `.cn/deps.json`. Mirrors the existing `ParseLockfile` / `ParsePackageIndex` shape; no IO.
- **`restore.ReadManifest(path) (*pkg.Manifest, error)`** (`src/go/internal/restore/restore.go`): IO wrapper around `ParseManifest`. Preserves the `eng/go §2.17` `Parse*`/`Read*` purity split.
- **Six AC-named tests** in `restore_test.go`:
  - `TestGenerateLockFromIndex_FiltersByManifest` — AC1 + AC4 (reproducer from #250: 3-package × 2-version index, 1 pin → lockfile len 1).
  - `TestGenerateLockFromIndex_NameAppearsAtMostOnce` — AC2 (3 versions in index, 1 pinned → exactly one entry).
  - `TestGenerateLockFromIndex_RestoreInstallsOnlyPinned` — AC3 (integration with real tarballs: 2 pkgs in index, 1 pinned → only pinned installed).
  - `TestGenerateLockFromIndex_MissingManifest` — no silent fallback when `deps.json` is absent.
  - `TestGenerateLockFromIndex_PinNotInIndex` — unresolved pins error explicitly with the full `name@version`.
  - `TestGenerateLockFromIndex_DefaultProfile` — the two-entry `cn setup` default (`cnos.core` + `cnos.eng`) round-trips through lock without leaking a third `cnos.cdd` from the index.

### Changed

- **`restore.GenerateLockFromIndex`** rewritten (`src/go/internal/restore/restore.go`): reads manifest, iterates `m.Packages`, looks up each pin, collects unresolved pins into a single error message, writes the sorted lockfile. Godoc spells out the new contract and every error case. Output remains sorted by `(name, version)` for deterministic content.
- **CI `kata-tier2` setup** (`.github/workflows/ci.yml`): now derives `cnos.core` / `cnos.kata` / `cnos.cdd.kata` versions from `src/packages/*/cn.package.json` at job time and writes the canonical array-schema `deps.json`. Previously pinned `cnos.core@3.54.0` (stale) with object-syntax `packages: {name: version}` (the Go and OCaml parsers both reject that shape; the old bug masked the error). The fix mechanically closes both drift paths.
- **Tier-1 kata `scripts/kata/06-install.sh`**: overwrites the `cn setup` default `deps.json` with an explicit pin to the real `cnos.core` version read from `src/packages/cnos.core/cn.package.json`, and tightens the post-condition from "≥ 1 installed" to "exactly `cnos.core` installed" — the over-vendoring regression class is now blocked at CI time.
- **`cnos.kata` `lib.sh` `write_deps_json`** now emits `packages` as an array of `{name, version}` objects. The comment documents why the pre-#250 object-syntax worked "by accident" (lockfile dumped the entire index regardless of parse result).
- **`cnos.kata` `lib.sh` `pkg_version_from_source <pkg> <repo_dist>`** added as a shared helper: reads `src/packages/<pkg>/cn.package.json` and echoes the version. R3 (`doctor-broken`) and R4 (`self-describe`) katas now use it instead of hardcoding `3.54.0`.

### Not in scope (explicit)

- Semver range resolution (`*`, `^`, `~`) — issue #250 defers this; a pin of `*` now errors as "not in index" because `*` is not a real version key.
- Transitive dependency resolution — package manifests do not yet declare a `dependencies` field; when that lands, `GenerateLockFromIndex` will need to walk the graph.
- Schema evolution — `cn.deps.v1` and `cn.lock.v2` unchanged.

### Known debt

- `cn setup` pins `deps.json` to the **binary's** version. In any environment where the binary version is not a real key in the local package index (dev checkouts, CI without `-ldflags`), `cn deps lock` immediately after `cn setup` will now error explicitly. The Tier-1 kata works around this; operators in similar setups can either build the binary with `-X main.version=` or hand-edit `deps.json`. A follow-up to make `cn setup` resolve against the local index is a candidate issue.

---

## 3.54.0 (2026-04-14)

### Added

- **Tier 1 bare-binary kata suite** (#236): six executable scripts at `scripts/kata/0{1..6}-*.sh` proving the `cn` binary works end-to-end before any package is installed — `01-binary` (help listing), `02-init` (hub creation), `03-status` (hub state read), `04-doctor` (clean-hub validation), `05-build` (dist production), `06-install` (package restore). Each script's header comment names exactly what it proves; no overclaiming.
- **`scripts/kata/run-all.sh` stop-on-first-failure** (#236): iterates `[0-9]*.sh` in filename order, exits non-zero on the first failure so CI surfaces the earliest break.
- **CI `kata-tier1` job** (#236): `.github/workflows/ci.yml` runs the suite on every PR, depends on `go` so it sees the freshly-built binary, and gates merges.
- **Three-tier kata catalog** (`docs/gamma/cdd/KATAS.md`, #236): Tier 1 = bare binary (this release), Tier 2 = runtime/package in `cnos.kata` (#237), Tier 3 = method/CDD in `cnos.cdd.kata`. `scripts/smoke/` remains separate for release-bootstrap compatibility.

### Changed

- **`cn help` always lists the 8 kernel commands** (#236): `cmd_help.go` now emits the Kernel section unconditionally, annotating hub-requiring commands with `(requires hub)` when no hub is discovered. Enables kata 01 to run before `cn init`.
- **`cn doctor` three-way Status** (#236): `internal/doctor/doctor.go` replaces the binary pass/fail with `StatusPass` (✓), `StatusInfo` (○), `StatusFail` (✗). Only `StatusFail` drives the exit code; `HasFailures` checks only `StatusFail`. Fresh-hub pending artifacts (`deps.lock.json`, runtime contract, state/runtime-contract.json, git origin remote, `.cn/vendor/packages/` before any lock) report informationally rather than fatally. `deps.json` after `cn setup` remains fatal when missing. Present-but-invalid artifacts still fail.
- **`cmd_doctor.go` glyph renderer**: three-way glyph selection matches the three-way Status; `cli/` stays a thin wrapper per T-002.

### Removed

- **`Registry.Available()` dead code** (#236): previously gated pre-hub vs post-hub command lists; the always-list help policy makes it unused. Removed to keep the registry accessor surface minimal.
- **`scripts/kata/01-boot.sh`, `02-command.sh`, `03-roundtrip.sh`** (#236): post-package behavior moved to Tier 2 (`cnos.kata`, #237). Old `04-doctor.sh` rewritten as the clean-hub doctor check that Tier 1 actually needs.

---

## 3.52.0 (2026-04-12)

### Added

- **Package command scanning** (#226): `internal/discover/ScanPackageCommands` walks `.cn/vendor/packages/*/cn.package.json`, parses `commands` entries, creates exec-backed `Command` implementations registered into the unified registry.
- **Repo-local command scanning** (#226): `internal/discover/ScanRepoLocalCommands` walks `.cn/commands/cn-*`, creates `Command` from filename convention. Operator overrides at tier 1.
- **External command dispatch** (#226): `ExecCommand.Run()` execs entrypoint scripts with `CN_HUB_PATH`, `CN_PACKAGE_ROOT`, `CN_COMMAND_NAME` environment variables.
- **Tiered help output** (#226): `cn help` groups commands by Kernel / Repo-local / Package with source attribution (`[cnos.core]`).
- **Command integrity validation** (#226): `cn doctor` reports missing entrypoints, non-executable files, and duplicate command names within a tier.
- **Entrypoint path confinement** (#226): manifest entrypoints validated to stay within package directory via `filepath.Rel` — rejects path traversal (`../../`).
- **`FullPackageManifest` type** (#226): pure type in `internal/pkg/` for parsing `cn.package.json` with command entries. Deterministic ordering via `slices.Sort`.
- **13 new tests** (#226): 8 in `discover/`, 5 in `pkg/`.

### Changed

- **`doctor.RunAll` signature** (#226): accepts optional `commandIssues` parameter for command integrity reporting.
- **`DoctorCmd` struct** (#226): receives `Registry` pointer from `main.go` for cross-command validation.

---

## 3.41.0 (2026-04-10)

### Added

- **TypeScript skill** (`eng/typescript`): Converged through 2 iterations. Schema-backed trust boundaries (Zod/Valibot/ArkType), branded primitives for domain scalars (`AbsolutePath`, `Checksum`, `PackageName`), explicit error policy (thrown `Error` vs discriminated results by layer), mutation discipline (`readonly` default, return new values), idempotence/retry/receipts for stateful operations. Two katas (package index client, async command runner). Registered in `cnos.eng`.

### Changed

- **Move 2 slice 4 (final)** (#201, PR #202): Activation frontmatter parser + validation types extracted into `src/lib/cn_frontmatter.ml` — 3 types (`frontmatter`, `issue_kind` 3-variant, `issue`) + 1 constant + 6 parsers + 2 helpers. 21 ppx_expect tests. 1 review round, **zero findings** — first Move 2 cycle with none. §2.5b checks 7+8 validated: the failure classes from slices 2+3 did not recur.
- **Move 2 complete.** `src/lib/` now holds 5 pure modules (cn_json, cn_package, cn_contract, cn_workflow_ir, cn_frontmatter). Every pure type and parser lives in `src/lib/`; every IO function lives in `src/cmd/`. The Go kernel rewrite (#192) has its boundary contract.
- **CDD §1.4 Roles formalized**: Two-agent minimum for substantial cycles. Small-change exception for §1.2-qualifying work. Implementer as optional delegated role. Reviewer is the default releaser (carries independent evaluation into assessment). Merge is part of step 9. Role column added to §5.3 lifecycle table.
- **CDD post-release ownership**: Releasing agent owns steps 11–13. Reviewer-default rationale documented.
- **v3.39.0 + v3.40.0 stacked post-release assessments** (#199, PR #200): Both cycles assessed retroactively. #198 cited as §12a corrective for both.

---

## 3.40.0 (2026-04-09)

### Added

- **Go skill** (`eng/go`): Runtime/kernel Go craft skill converged through 4 iterations. Covers package design, concrete types, consumer-owned interfaces, context flow, error wrapping with `errors.Is`/`errors.As`, resource lifecycles (`defer` discipline), structured observability (`log/slog`), side-effect boundaries (adapter injection), concurrency discipline, testing, build. cnos-specific: schema/compatibility, determinism, idempotence, traceability/receipts, runtime boundaries (commands/orchestrators/extensions/skills), shell/archive safety, override precedence. Two katas (package restore, contract extraction). Registered in `cnos.eng` package.
- **CDD §2.5b checks 7+8** (#198): Check 7 — workspace-global library-name uniqueness (`grep` before new dune stanza). Check 8 — CI green before requesting review (draft-until-green pattern). Rooted in #195 F1 + #197 F1.
- **Telegram CI notifications**: All 3 workflows (ci.yml, coherence.yml, release.yml) notify on success/failure. Graceful skip if secrets not set.

### Changed

- **Move 2 slice 2** (#194, PR #195): Runtime contract types extracted into `src/lib/cn_contract.ml` — 11 types (`package_info`, `override_info`, `zone`, `zone_entry`, `identity`, `extension_contract_info`, `command_entry`, `orchestrator_entry`, `cognition`, `body_contract`, `runtime_contract`) + `zone_to_string` + transitive `activation_entry` via chained re-export. `cn_runtime_contract.ml` re-exports via type-equality. 13 ppx_expect tests. 2 review rounds (F1: library name collision, mechanical).
- **Move 2 slice 3** (#196, PR #197): Workflow IR extracted into `src/lib/cn_workflow_ir.ml` — 6 types (`trigger`, `permissions`, `step` 6-variant, `orchestrator`, `issue_kind` 7-variant, `issue`) + 10 pure functions (parsers, validator, helpers, result combinator). Largest `src/lib/` module (316 lines). Option-(b) principle: `load_outcome`/`installed`/`outcome` stay in `cn_workflow.ml` (no `src/lib/` consumer). `Printf.eprintf` stderr precedent documented in CORE-REFACTOR.md §7. 20 ppx_expect tests. 2 review rounds (F1: expect-test stderr mismatch, mechanical).
- **CORE-REFACTOR.md §7**: Discipline extended — stderr diagnostic logging permitted for discovery-time warnings. Status blocks for v3.39.0 + v3.40.0 appended.
- **`src/lib/`**: Now 4 pure modules (cn_json, cn_package, cn_contract, cn_workflow_ir). Move 2 is 3/4 complete.

### Fixed

- **Go skill package placement**: Moved from `cnos.core` to `cnos.eng` — was causing coherence CI drift check failure.
- **v3.38.0 lag table**: TBD entry replaced with #193 (`llm` step execution).

---

## 3.35.0 (2026-04-07)

### Added

- **Runtime contract activation index** (#173): cognition layer now exposes `activation_index.skills` with declarative triggers parsed from SKILL.md frontmatter. `cn_activation.ml`: frontmatter parser, activation table builder, conflict validator.
- **Command registry in runtime contract** (#173): body layer includes `commands` registry via `Cn_command.discover`.
- **Orchestrator registry in runtime contract** (#173): body layer includes `orchestrators` registry from package manifests.
- **ORCHESTRATORS.md design doc** (#170): four-surface architecture (skills decide, commands dispatch, orchestrators execute, extensions provide capability). CTB as source language, Effect Plan as IR.

### Changed

- **Markdown/JSON parity enforced** (#177, F3): markdown projection of `activation_index` now renders two-level nesting (`activation_index: > skills: > entries`) matching JSON schema shape.
- **Silent error suppression eliminated** (#177, F1/F2/R2): all manifest parse errors, malformed entries, and leaf-value coercions logged with package context in `cn_runtime_contract.ml` and `cn_activation.ml`.
- **#170 dependency graph cleaned**: #172 independent, #173 depends on #170 (not #172), #174 explicitly owns orchestrator package-class.
- **ORCHESTRATORS.md**: ETB replaced with CTB terminology (#172).

### Fixed

- **Retro-review findings** (#177): 5 findings from post-merge retro on PR #176 all addressed — silent orchestrator logging (F1/F2), parity (F3), missing tests (F4), stale bare catch confirmed absent (F5).

### Process

- PR #176 merged without CDD review — process violation corrected via retro-review and findings PR #177. WORKFLOW_AUTO.md created as mechanical gate to prevent recurrence.

## 3.32.0 (2026-04-04)

### Fixed

- **Legacy fallback path audit** (#152): systematic audit of all legacy/fallback/deprecated paths in src/. 29 findings classified in frozen audit table.
  - **Removed (7):** `run_inbound`, `feed_next_input`, `wake_agent` (~160 lines dead code since v3.27), `hub_flat_mindsets_path`, `hub_flat_skills_path` + loader blocks (~25 lines, package namespace only since v3.25).
  - **Converted (15):** all bare `with _ ->` in touched files converted to `Printf.eprintf` warning + fallback. Covers `cn_assets.ml` (9), `cn_deps.ml` (3), `cn_agent.ml` (2), plus 1 narrowed catch.
  - **Kept (7):** `git.ml` main/master fallback, `release_lock` cleanup, `extract_inbound_message` defensive fallback, `resolve_render` traced fallback, `cn_dotenv.ml` optional load, `cn_sandbox.validate_path` (fail-closed). All justified.
  - 5 review rounds, 10 findings (3D in R1, 3D in R2, 2 in R3, 1 in R4, 0 in R5). Net: ~185 lines deleted, ~40 added.

### Added

- **Thread event model design** (#153): `docs/alpha/protocol/THREAD-EVENT-MODEL.md` v1.0.1 — canonical ID (`{id}@{author}`) vs locator split, parent-linked publication, Git-first transport-flexible. Implementation plan: `docs/gamma/cdd/3.32.0/PLAN-thread-event-model.md`.
- **Hub placement models design** (#156): `docs/alpha/HUB-PLACEMENT-MODELS.md` — split `hub_root`/`workspace_root` for sandboxed agents, placement manifest (`.cn/placement.json`), nested clone as default backend, submodule optional. Implementation plan: `docs/gamma/cdd/3.33.0/PLAN-hub-placement-models.md`.

### Changed

- **CDD bundle normalization**: all 7 CDD skill artifacts (`cdd`, `design`, `review`, `release`, `post-release`, `plan`, `issue`) have full frontmatter (`artifact_class`, `kata_surface`, `governing_question`). Descriptions sharpened as "How do we..." process questions. Kata expanded to Scenario/Task/Verification/Common failures. Package copies synced.

---

## 3.29.1 (2026-04-02)

**Fix daemon crash loop from false version drift**

### Fixed

- **`resolve_bin_path` uses `Unix.readlink` directly (#148):** Spawning `readlink -f /proc/self/exe` via `Child_process.exec` ran in a bash child whose `/proc/self/exe` was `/usr/bin/bash`. This caused `check_binary_version_drift` to run `bash --version`, detect false drift ("GNU" ≠ "3.29.0"), and trigger a failing re-exec loop every maintenance cycle. Now uses `Unix.readlink "/proc/self/exe"` in-process.

### Changed

- **Release skill §2.5:** Added RELEASE.md step with format template. CI uses `RELEASE.md` as GitHub release body. §3.7 mechanical gate: no tag without RELEASE.md.

## 3.29.0 (2026-04-02)

**Remove hardcoded paths — runtime bin_path + build-time repo**

### Fixed

- **Hardcoded `bin_path`** replaced with `resolve_bin_path()`: `$CN_BIN` env var → `readlink /proc/self/exe` (Linux) → `readlink Sys.executable_name` (macOS/fallback) → `Sys.executable_name`. Self-update now works for any install prefix where user has write access.
- **Hardcoded `repo`** replaced with `resolve_repo()`: `$CN_REPO` env var → `Cn_lib.cnos_repo` (build-time from `cn.json`).
- **`cn_deps.ml` default source** derived from `Cn_lib.cnos_repo` instead of hardcoded `"usurobor/cnos"`.

### Added

- `cn_repo_info.ml` — dune-generated module extracting repo from `cn.json` at build time. Build fails if no GitHub URL found (no silent fallback).
- `test/cmd/cn_selfpath_test.ml` — 7 ppx_expect tests: env var override precedence, fallback behavior, format validation, derived source, negative space.
- `$CN_BIN` and `$CN_REPO` environment variable overrides for non-standard installations.

### Removed

- All `/usr/local/bin/cn` literals from `src/**/*.ml`.
- Hardcoded `"usurobor/cnos"` from `cn_agent.ml`.

## 3.28.0 (2026-04-02)

**Orphan rejection loop fix + inbox skill v3**

### Fixed

- **Orphan rejection deduplication (#144):** `rejection_filename` produces deterministic names keyed on (peer, branch), replacing `make_thread_filename` which prepended timestamps and caused unbounded rejection spam (~12/hour per orphan). `is_already_rejected` checks both outbox and sent dirs before creating a new rejection.
- **Authoritative sender identity (#144):** Removed `get_branch_author` (git log inference). Rejection messages now use fetched peer identity (`peer_name`) exclusively. Added `from:` field to rejection frontmatter envelope.
- **Self-send guard (#144):** `send_thread` rejects `to == my_name` before peer lookup, emitting `outbox.skip` trace event with `reason_code:"self_send"`.
- **Stale rejection hint wording (#144):** Rejection message updated from clone-based model to current pull-only protocol ("Push `peer/topic` to your own hub repo so the recipient can fetch it").

### Added

- `test/cmd/cn_mail_test.ml` — 9 ppx_expect tests: filename determinism (3), dedup across outbox/sent (4), self-send detection (2), no-timestamp-prefix (1), no-clone-wording regression (1).

### Changed

- **Inbox skill v3:** Rewritten to meta-skill standard — `artifact_class: skill`, `kata_surface: embedded`, `governing_question`. Define/Unfold/Rules/Verify structure. Pull-only transport model with contrastive examples. Two embedded katas (review request triage, rejection loop handling). Visible closing artifact requirement for every message.

### Removed

- `get_branch_author` — dead code after sender identity fix.

## 3.27.1 (2026-04-02)

**Remove cron installation management + CDD trace normalization**

### Removed
- `update_cron` function from `cn_system.ml` — no longer reinstalls OS crontabs on `cn update` (#138)
- `cn-cron` shell script (`src/cli/cn-cron`) (#138)
- Cron installation guidance from AUTOMATION.md, README prerequisites, CLI.md (#138)

### Changed
- AUTOMATION.md rewritten: daemon-first setup, oneshot as secondary manual/scripted mode (#138)
- README prerequisites: `systemd (recommended)` replaces `System cron or systemd` (#138)
- CLI.md: `(cron mode)` label removed from `cn agent` (#138)
- PR template headings: `§2.3`/`§2.4` → `step 4`/`step 5` (#143)
- Review skill §2.0.8: `step §2.4` → `step 5` (#143)

### Known debt
- `Agent.mode.Cron` and `run_cron` naming — could be renamed to `Oneshot`/`run_oneshot`
- `CN_CRON_PERIOD_MIN` env var still read in `cn_agent.ml` (dead code)
- Deprecated `run_inbound`, `wake_agent` functions still exist in `cn_agent.ml`

---

## 3.27.0 (2026-04-01)

**CDD traceability + skill-loading bridge**

### Added
- CDD §5.3 artifact manifest: step-to-evidence binding for all 13 lifecycle steps
- CDD §5.4 CDD Trace: lightweight execution trace format for primary branch artifacts
- CDD §3.7: no direct-to-main without retro-closure (refs canonical §12)
- Review §2.0.8: CDD execution trace verification
- Release §2.9: CDD Trace update row after release
- Post-release §6: CDD Closeout block for steps 11–13
- eng/README: default active-skill matrix by work shape × engineering level
- Design template: optional Engineering Level field + CDD Trace section

### Changed
- CDD §2.4: active skills selected by work shape → level → dominant risk; level labels require concrete skill names
- Release §2.4: CHANGELOG bound to canonical ledger format; writing skill required for release notes
- Writing skill: full rewrite — governing question, stable facts, revision pass, word-level test
- CAP skill: v2 rewrite — UIE-based, frontmatter, failure modes, kata surface
- Meta-skill (skill/): rewrite — artifact classification, domain formula, self-demonstration
- Skills README: normalized to match meta-skill contract

### Fixed
- Process debt: 3 direct-to-main commits retro-closed (#137)

### Assessment
- α: A — five skill surfaces now carry the same trace model
- β: A- — CDD ↔ design ↔ review ↔ release ↔ post-release all cross-reference; some older doc refs (§5 numbering in assessment examples) may drift
- γ: B+ — three commits landed direct-to-main before §3.7 existed; retro-closed via #137 but the gate was skipped

---

## v3.15.2 (2026-03-25)

**Empty Telegram filter (#29), CDD canonical rewrite, review skill hardening**

Three concerns addressed in one release: a runtime bugfix, the CDD authority split, and review skill improvements.

### Fixed

- **Empty Telegram messages filtered before enqueuing (#29)** — photos, stickers, edits have `text=""` and caused Claude API 400 → infinite retry loop (#28). `drain_tg` now checks `String.trim msg.text = ""`, advances offset, emits trace event (`reason_code: empty_content`). Mirrors existing `rejected_user` pattern.

### Changed

- **CDD.md rewritten as canonical algorithm spec (v3.15.0)** — absorbs §0 selection, 14-step lifecycle, cycle-close, operational debt override, and all sections previously carried only by the skill. 12 sections: purpose, scope, inputs, selection function, lifecycle, artifact contract, mechanical/judgment boundary, review, gate, assessment, closure, related docs.
- **CDD SKILL.md rewritten as executable summary** — clean DUR structure, canonical doc governs on disagreement, no authority exception blocks. Delegates to sub-skills.
- **RATIONALE.md created** — companion document: why CDD is closed-loop, artifact-driven, and not fully mechanical. Absorbs the "why" that doesn't belong in spec or skill.
- **README.md rewritten** — self-describing topic bundle with reading order and document map.
- **AGILE-PROCESS.md demoted to reference profile** — explicitly one valid implementation of CDD, not the method itself.
- **Review skill §2.0.5** — CDD artifact contract check: verify required CDD artifacts exist for the change class.
- **Review skill §2.2.8** — authority-surface conflict: when multiple surfaces define the same thing, verify they agree.
- **Review skill §2.3.6** — evidence depth: match evidence depth to claim strength (predicate ≠ integration).
- **Review skill §3.7** — CI/release-gate state: don't issue merge instruction when required checks haven't run.

### Process

- **CDD §0 selection** — #29 selected via §0.6 operational infrastructure debt override.
- **PR #103** — authored by usurobor via Claude Code session, reviewed by Sigma. Reset+cherry-pick after history rewrite (Test→usurobor authorship fix via git filter-repo).
- **Git authorship fixed** — all historical `Test <test@test.local>` commits rewritten to `usurobor <usurobor@gmail.com>` via git filter-repo. Force-pushed.

---

## v3.16.0 (2026-03-25)

**End-to-end self-update (#37), CDD OVERVIEW, AGILE-PROCESS deletion**

### Fixed

- **`cn update` end-to-end self-update (#37)** — three independent failures fixed:
  1. Same-version patch detection via commit hash comparison (`Update_patch` variant, `release_info` record, `get_latest_release()` with `Cn_json` parsing)
  2. ARM binary in CI release matrix (`ubuntu-24.04-arm` → `cn-linux-arm64`)
  3. Release workflow triggers on bare version tags (`[0-9]*.[0-9]*.[0-9]*`) in addition to `v*`
  4. `target_commitish` validated as hex SHA before patch comparison — branch names like "main" no longer cause false `Update_patch`
- **`cn --version` shows commit hash** — `cn 3.16.0 (abc1234)` for operator diagnostics
- **Cram tests use prefix match** — `version.t` and `cli.t` handle dynamic commit suffix

### Added

- **CDD OVERVIEW.md** — plain-language introduction: what CDD is, why it exists when agents can write code, what one cycle looks like, what stays mechanical vs judgment-bearing
- **CDD §12 retro-packaging rule** — substantial direct-to-main changes must be followed by retro-snapshot + self-coherence + version-history entry
- **Post-release §5 closure evidence** — template now requires explicit immediate/deferred output accounting

### Changed

- **AGILE-PROCESS.md deleted** — CDD §4.1 lifecycle table supersedes the workflow states. P0–P3 priority bands implied by §3 selection function. One less surface to maintain.
- **CDD bundle reading order** — OVERVIEW → CDD → RATIONALE → epoch assessments
- **Review skill output format** — now includes CI state, merge instruction, and finding type column (closes self-contradiction where rules required them but template omitted them)
- **Release skill §2.3** — VERSION-first flow: `echo X.Y.Z > VERSION` → `stamp-versions.sh` → `check-version-consistency.sh`. Old manual locations removed.
- **Live AGILE-PROCESS refs updated** — ARCHITECTURE.md diagram, AUDIT.md inventory

### Process

- **PR #104** — 4 review rounds (3 D-level blockers found and fixed: syntax error, FSM exhaustive match, bare tag trigger, `target_commitish` validation, `cli.t` prefix match). Authored by usurobor via Claude Code, reviewed by Sigma.
- **Git authorship fixed** — all historical `Test <test@test.local>` commits rewritten to `usurobor <usurobor@gmail.com>` via git filter-repo.
- **CDD retro-packaged** — frozen CDD snapshot in `3.15.2/`, self-coherence report for the canonical rewrite.

---

## v3.16.1 (2026-03-25)

**Daemon retry limit and dead-letter (#28)**

### Fixed

- **Daemon retry limit and dead-letter (#28)** — `drain_tg` no longer retries failed triggers forever. Three changes:
  1. Per-trigger retry counter with exponential backoff (1s, 2s, 4s... capped at 30s)
  2. Error classification: 4xx → dead-letter immediately, 5xx/network → retry with backoff
  3. Dead-letter path: advance Telegram offset, emit trace event, clean up stale state files, continue processing

### Added

- **Sustainability surface** — `.github/FUNDING.yml`, `docs/beta/SUSTAINABILITY.md`, README Support section

### Process

- **PR #105** — 2 review rounds. R1: `Fs.remove` → `Fs.unlink` build fix. AC4 (backoff) pushed from "met via poll interval" to real exponential backoff after review feedback.

---

## v3.13.0 (2026-03-24)

**Docs Governance: CDD Pipeline, Self-Coherence, Feature Bundles (#75)**

Development method and documentation governance made explicit and self-enforcing.

### Added

- **CDD §5.1 pipeline table** — 9 steps (0–8), each with explicit deliverable artifacts and locations. Bootstrap-first rule: first diff on a branch must create the version directory with stubs.
- **CDD §7.8 self-coherence report** — new artifact format (`SELF-COHERENCE.md`) in version directories. Records pipeline compliance, triadic assessment, checklist pass, known debt, and reviewer notes.
- **DOCUMENTATION-SYSTEM.md §3 single version lineage** — all docs use cnos release versions. No independent per-document version lineages. Frozen legacy snapshots explicitly allowed to retain historical versions.
- **Feature bundles** — `docs/alpha/agent-runtime/`, `docs/alpha/runtime-extensions/`, `docs/gamma/cdd/` with bundle READMEs, version directories, and frozen snapshots.
- **CDD skill §4.0 bootstrap** — executable skill now requires version-directory creation before any artifacts.
- **CDD skill §4.8 self-coherence** — executable skill now requires self-coherence report before review.

### Changed

- **CDD §7.3 process artifacts** — fixed stale references: `RULES.md` → `docs/gamma/RULES.md`, removed nonexistent `RELEASE.md`, `PLAN.md` → `docs/gamma/plans/`.
- **CDD §5.0 branch naming** — canonical format `{agent}/{issue}-{scope}-{version}`, tooling suffix rule.
- **Placeholder syntax** — governance docs now use `{placeholder}` instead of `<placeholder>` to prevent GitHub markdown stripping.

### Removed

- **CURIOSITY mindset** — removed `packages/cnos.core/mindsets/CURIOSITY.md`.

---

## v3.12.1 (2026-03-23)

**Daemon Boot Log: Configuration Sources (#61)**

Operators no longer have to probe files to understand what the daemon is running.

### Added

- **Boot banner** — daemon prints structured config declaration before "Daemon started" line: version, hub name, profile, model, secrets source, peers.
- **`secret_source` type** in `cn_dotenv.ml` — `Env | File | Missing` with no value payload. Makes secret leakage structurally impossible at the type level.
- **`probe_source`** function — checks where a secret key is configured without returning its value.
- 5 ppx_expect tests for probe_source (env precedence, file fallback, missing, type safety).

---

## v3.12.0 (2026-03-23)

**Wake-up Self-Model Contract (#56) + CDD/Review Skill Hardening (#57)**

The agent now receives a first-class Runtime Contract at every wake, replacing the overloaded capabilities-only block. Development skills hardened with issue AC reconciliation gates.

### Added

- **cn_runtime_contract.ml** — new module: gather/render_markdown/to_json/write. Structured Runtime Contract with three sub-blocks: self_model (version, hub, packages, overrides), workspace (directories, writable/protected), capabilities (observe/effect ABI, budgets, config).
- **RUNTIME-CONTRACT-v3.10.0.md** — design doc for wake-up self-model contract.
- **COGNITIVE-SUBSTRATE §6.1** — Runtime Contract as normative required runtime-generated block with invariant: agent answers self-knowledge from packed context alone.
- **cn doctor** validates runtime contract (self_model + workspace + capabilities + agent/ existence).
- **cn setup/init** creates `agent/` override directory.
- **CDD skill §2.4** — reconcile coherence contract with issue ACs.
- **CDD skill §4.5.1-3** — AC self-check at commit, multi-format parity rule, build-sync rule.
- **CDD skill §8.5** — author pre-flight AC walk before requesting review.
- **CDD skill §9.1** — issue AC verification as outermost release gate.
- **Review skill §2.0** — walk issue ACs before reading the diff.
- **Review skill §2.2.3** — multi-format parity check for dual-serialization modules.
- 10 new ppx_expect tests for runtime contract.

### Changed

- **cn_capabilities.ml** — `max_passes` now reads from `shell_config.max_passes` instead of hardcoded 2.
- **cn_context.ml** — emits `## Runtime Contract` with three sub-blocks instead of bare `Cn_capabilities.render`.
- **state/runtime-contract.json** — persisted at pack time for operator inspection; JSON now at full parity with markdown (observe/effect/apply_mode/exec_enabled/max_passes/budgets).
- Protected paths derived from `Cn_sandbox.default_denylist` + `protected_files` (single source of truth).

### Fixed

- Package/source sync for agent-ops (N-pass vocabulary) and release skill.

---

## v3.11.0 (2026-03-22)

**N-pass Merge + Misplaced Ops Correction (#51)**

### Added

- **Misplaced ops detection** — detect ops in body text and retry via correction pass (#51).
- **N-pass bind loop merged** to main with review fixes.

### Changed

- Structured output contract (#52) reverted — needs rework for tool_use integration.

---

## v3.10.0 (2026-03-22)

**N-pass Bind Loop (#50)**

### Added

- **Generic N-pass bind loop** — effect pass no longer terminal; loop continues up to `max_passes` as long as there are typed ops.
- **Processing indicators** — agent sees pass labels in context.
- Renamed `two_pass` → `n_pass` across codebase.
- Renamed `run_observe_pass` → `run_pass`.
- Removed observe-centered vocabulary residue.
- Aligned stop reasons and trace reason codes.

---

## v3.9.3 (2026-03-21)

**Anti-Confabulation: Op Result Integrity (#49)**

(Previously: v3.9.1 entry below was mislabeled)

---

## v3.9.1 (2026-03-21)

**Anti-Confabulation: Op Result Integrity (#49)**

Agent no longer confabulates when ops fail or are denied. Explicit result signals in Pass B repack, anti-fabrication instructions in capabilities block and doctrine.

### Added

- **Receipt result signals** — `receipts_summary` now tags each receipt with `[EMPTY_RESULT]`, `[NOT_EXECUTED]`, or `[FAILED]` to distinguish "op succeeded with empty data" from "op was not executed." Pass B context carries these signals so the agent knows what actually happened.
- **Failure WARNING banner** — When any Pass A receipt has status denied or error, the repack includes a `WARNING` block instructing the agent not to fabricate results.
- **Anti-confabulation CRITICAL** in capabilities block — New runtime-generated instruction: "When an op returns status=denied, status=error, or status=ok with zero artifacts, you MUST report the actual result to the user."
- **Op Result Integrity doctrine** — New section in AGENT-OPS.md (both `packages/cnos.core/doctrine/` and `src/agent/doctrine/`) with non-negotiable anti-confabulation rules.
- **6 new orchestrator tests** — `receipts_summary` signal tags for ok/empty/denied/error/skipped/mixed scenarios.

### Changed

- **`receipts_summary` format** — Denial reasons now prefixed with `reason:` for clarity. Each receipt line includes an explicit result signal tag.
- **Capabilities block** — All configs now emit the anti-confabulation CRITICAL line after the ops-in-frontmatter CRITICAL line.

---

## v3.9.0 (2026-03-20)

**Two-Pass Wiring + Cognitive Substrate + DUR Skill Contract**

Runtime two-pass execution wired end-to-end. COGNITIVE-SUBSTRATE spec published. Three skills cohered to DUR contract. Glossary updated.

### Added

- **Two-pass execution wired** (#41) — `Cn_orchestrator.run_two_pass` coordinates Pass A → repack → LLM call → parse → Pass B → coordination gating. `finalize` accepts `?packed` for Pass B LLM re-call. 7 integration tests with mock LLM.
- **COGNITIVE-SUBSTRATE.md** v1.0.0 — canonical classes (Identity, Doctrine, Mindset, Skill, Reflection), distinction rule, placement algorithm, file contracts, promotion/splitting rules, validation checks.
- **DUR glossary entry** — Define/Unfold/Rules documented as canonical skill contract in GLOSSARY.md.
- **`cn_runtime_integration_test.ml`** — new integration test file for two-pass orchestration.
- **`issue-41-pass-b-wiring.md`** — design doc for two-pass wiring.

### Changed

- **Reflect skill** — cohered to DUR contract. Evidence/Interpretation/Conclusion remapped to Define/Unfold/Rules. Axes (α/β/γ) moved under Unfold. Cadence table under Rules.
- **Adhoc-thread skill** — cohered to DUR contract. Added trigger recognition, type-matching judgment, bias-toward-capture rule.
- **Review skill** — cohered to DUR contract with MCI from PR #32 review process.
- **§7.2 mindset contract** — relaxed to structurally flexible, semantically constrained.

### Fixed

- Integration test pattern match updated for `pass_b_output` field.

---

## v3.8.0 (2026-03-20)

**Syscall Surface Coherence Amendment**

The CN Shell ABI is now honest, orthogonal, and hardened. Four incoherences closed: advertised-but-stub ops, implicit compound behavior, lossy observation, and hidden external dependencies. Multi-round review-driven development (PR #32).

### Added

- **fs_glob** — real implementation replacing stub. Pure OCaml glob matching (*, **, ?), sandbox-validated, symlink cycle detection via realpath, bounded output.
- **git_stage** — explicit staging op. Two modes: named files (sandbox-validated per file, directory rejection, `--literal-pathspecs`) or stage-all (NUL-delimited porcelain parser, per-candidate `validate_path`). Separates staging from commit.
- **fs_read chunking** — `offset` and `limit` fields for bounded reads of large files. Budget-capped.
- **cn doctor** — now checks for `patch(1)` (required by fs_patch).
- **Observe exclusion symmetry** — all four git observe ops (status, diff, log, grep) now use shared `git_observe_exclusions` (`-- . :!.cn :!state :!logs`). Runtime-owned paths hidden from agent.
- **SYSCALL-SURFACE-v3.8.0.md** — design doc specifying all four fixes.
- **PLAN-v3.8.0-syscall-surface.md** — implementation plan.

### Changed

- **git_commit** — no longer does implicit `git add -A`. Commits current index only. Returns `Skipped/nothing_staged` when index is empty. Use `git_stage` before `git_commit`.
- **git_grep** — uses `:(literal)` per-path instead of global `GIT_LITERAL_PATHSPECS`, appends observe exclusion pathspecs. Uses `-e` to force query as pattern.
- **git_diff / git_log** — `validate_rev` rejects leading-dash revisions (prevents CLI injection).
- **Observe vs write exclusions split** — protected files (SOUL.md, USER.md, etc.) are readable by git observe ops but not stageable by git_stage.

### Security

- `--literal-pathspecs` on git_stage prevents glob/pathspec magic interpretation
- `validate_rev` blocks `--output=<file>` style injection via leading-dash revisions
- `fs_glob` validates paths before descending, prevents symlink attacks
- `git_stage` stage-all rejects directories and symlinks to protected files

### Tests

- 32 → 69+ ppx_expect tests in cn_executor_test.ml
- Regression pairs for each hardening: positive case (feature works) + negative case (incoherence blocked)
- Capstone: git_status shows src/ changes while hiding .cn/, state/, logs/; prior observe op artifacts don't pollute subsequent git_status

### Coherence Delta

- **Gap:** Four ABI incoherences — stub op, compound behavior, lossy reads, hidden dependency
- **Mode:** MCA — change the runtime surface
- **α:** Syscall ABI now honest and orthogonal — no advertised-but-unimplemented ops
- **β:** All git observe ops use same exclusion boundary; capabilities block matches executor behavior
- **γ:** git_stage/git_commit separation enables future fine-grained staging without special cases

---

## v3.7.1 (unreleased)

**Trace Gap Closure: Full Behavioral Reconstruction**

The trace system now captures everything needed to reconstruct agent behavior without external assumptions.

### Fixed

- **LLM call latency + model** — `llm.call.ok` and `llm.call.error` now include `model`, `latency_ms`. Success events already had `input_tokens`/`output_tokens`/`stop_reason`; error events now get `model` and timing too.
- **Drain trigger_ids** — `drain.complete`/`drain.stopped` include `trigger_ids` array showing exactly which items were processed, plus `duration_ms`.
- **Sync duration** — `sync.ok` and `sync.error` include `duration_ms`. Error events include structured `step` and `exit_code` fields.
- **Config snapshot at boot** — `config.loaded` expanded with `max_tokens`, `sync_interval_sec`, `review_interval_sec`, `oneshot_drain_limit`, `daemon_drain_limit`, `telegram_configured`, `allowed_users` count.
- **Poll heartbeat** — `daemon.heartbeat` (Debug) emits every 60s with `polls_since_last` and `uptime_sec`. Hours of silence now distinguishable from dead daemon.
- **Cleanup trace** — `cleanup_once` emits `cleanup.complete` with status=Ok_ (removed count), Skipped (nothing_to_clean), or Degraded (cleanup_failed). No more silent success.

---

## v3.6.0 (unreleased)

### Fixed

- **Two-pass execution wired in cn_runtime (Issue #41)** — `finalize` now invokes Pass B when `needs_pass_b=true`. Previously, observe ops executed but the LLM was never re-called with evidence, making the sense→act loop dead. New `Cn_orchestrator.run_two_pass` coordinates: Pass A → repack with artifacts/receipts → LLM re-call → Parse → Pass B → coordination gating. Pass-A-safe coordination ops (ack, surface, reply) execute during Pass A; terminal ops (done, fail, send, delegate) defer to Pass B and are gated on effect success. Projection and conversation history use the final (Pass B) output. 7 integration tests added.

- **`cn update` binary-only (Issue #27)** — `cn update` now exclusively downloads pre-built binaries from GitHub Releases, mirroring `install.sh`. Removed the git-based update path entirely: no more `/usr/local/lib/cnos` assumption, no `opam exec -- dune build`, no `has_git_install()` detection. One install method, one update method. `Update_git` variant removed from `update_info` type; replaced by `Update_available`. Self-update check at CLI startup uses the same binary-download path.

---

## v3.5.1 (2026-03-11)

**TRACEABILITY: Structured Observability + CDD**

Operators can now answer "did it boot? what did it load? why did it transition?" from files alone.

### Added

- **`cn_trace.ml`** — Append-only JSONL event stream (`logs/events/YYYYMMDD.jsonl`). Schema `cn.events.v1` with boot_id, monotonic seq, severity, layer (sensor/body/mind/governance/world), reason_code on every transition.
- **`cn_trace_state.ml`** — State projections: `state/ready.json` (mind/body/sensors), `state/runtime.json` (cycle/lock/pass), `state/coherence.json` (structural checks).
- **Boot telemetry** — Mandatory 9-event boot sequence: `boot.start` → `config.loaded` → `deps.lock.loaded` → `assets.validated` → `doctrine.loaded` → `mindsets.loaded` → `skills.indexed` → `capabilities.rendered` → `boot.ready`. Per-package skill counts and hub-override detection.
- **Cycle telemetry** — `cycle.start`/`cycle.recover` → `pack.start`/`pack.complete` → `llm.call.start`/`llm.call.ok` → `effects.execute.start`/`effects.execute.complete` → `projection.*` → `finalize.complete`.
- **Sensor telemetry** — `daemon.poll.*`, `daemon.offset.advanced`/`daemon.offset.blocked` with reason codes (`rejected_user`, `still_queued`, `still_in_flight`, `processing_failed`).
- **Governance events** — `pass.selected`, `ops.classified`, `policy.denied` from orchestrator with reason codes.
- **`Cn_hub.log_action` shim** — Resurrected as compatibility bridge to `Cn_trace.gemit`.
- **`Cn_agent.queue_depth`** — Real queue count in projections.
- **`docs/alpha/security/TRACEABILITY.md`** — Full observability spec (721 lines).
- **`docs/gamma/plans/TRACEABILITY-implementation-plan.md`** — 12-step implementation plan.
- **`docs/gamma/CDD.md`** — Coherence-Driven Development v1.1.0: development method applying CAP to the development process.
- **Doc graph updates** — TRACEABILITY wired into README.md, AGENT-RUNTIME.md updated to v3.3.7.

### Fixed

- **Trace test isolation** — `open Cn_cmd` invalid with `wrapped false`; tmp dir collisions from unseeded `Random.int`.
- **Lockfile `ls-remote` parse** — Stricter: `None -> ""` on malformed output.

---

## v3.5.0 (2026-03-10)

**Unified Package Model + Coherent Agent Architecture + FOUNDATIONS**

Everything cognitive is now a package. Three packages ship: `cnos.core` (doctrine + mindsets + agent skills), `cnos.eng` (engineering skills), `cnos.pm` (PM skills). Role profiles select which packages an agent loads.

### Added

- **Unified package model** — `packages/` directory with `cn.package.json` manifests; two-layer resolution (installed packages → hub-local overrides); replaces the three-layer core/packages/hub-local split from v3.4.0
- **`cnos.core` package** — doctrine (FOUNDATIONS, CAP, CBP, CA-Conduct, COHERENCE, AGENT-OPS) + all mindsets + agent/ops/documenting/testing/release/skill skills
- **`cnos.eng` package** — engineering skills (coding, design, functional, OCaml, RCA, review, ship, testing, tool-writing, UX-CLI)
- **`cnos.pm` package** — PM skills (follow-up, issue, ship)
- **Role profiles** — `profiles/engineer.json`, `profiles/pm.json` — select packages per role
- **CAA v1.0.0** (`docs/alpha/agent-runtime/CAA.md`) — Coherent Agent Architecture spec: structural definition of a coherent agent, 7 invariants, failure mode table, wake-up strata
- **FOUNDATIONS.md** (`docs/explanation/`, `packages/cnos.core/doctrine/`) — doctrinal capstone: first principle through runtime, four doctrinal layers, cognitive loop, hierarchy of guidance
- **Doc graph cleanup** — version numbers removed from design doc filenames (AGENT-RUNTIME-v3.md → AGENT-RUNTIME.md, CAR-v3.4.md → CAR.md, PLAN-v3.3.md → PLAN.md)

### Changed

- **`cn_assets.ml`** — rewritten for unified two-layer resolution (was three-layer)
- **`cn_deps.ml`** — updated for package-based materialization
- **`cn_context.ml`** — delegates to `Cn_assets` for doctrine/mindset/skill loading
- **`cn_system.ml`** — updated for package-aware `cn setup` and `cn init`
- **`AGENTS.md`** — removed template-repo assumption; updated for package model
- **`README.md`** — corrected env var names, removed template-repo references

### Removed

- **`docs/design/CAR-v3.4.md`** — superseded by `docs/alpha/cognitive-substrate/CAR.md` (unified model)
- Template-repo assumption throughout docs

---

## v3.4.0 (2026-03-09)

**Cognitive Asset Resolver & Package System**

Fresh hubs now wake with the full cognitive substrate — mindsets, skills, and agent-ops guidance — without needing a separate template repo checkout.

### Added

- **`cn_assets.ml`** — Three-layer asset resolver (core → packages → hub-local overrides). Deterministic loading order, keyword-matched skill scoring delegates here
- **`cn_deps.ml`** — Dependency manifest (`.cn/deps.json`), lockfile (`.cn/deps.lock.json`), materialize/restore pipeline. Git-native package distribution
- **`cn deps` CLI commands** — `list`, `restore`, `doctor`, `add`, `remove`, `update`, `vendor`
- **Cognitive Assets block** in capabilities — agent sees asset summary (core count, packages, overrides) on every wake-up
- **`cn doctor`** checks for `.cn/vendor/core/`, deps manifest, and lockfile
- **`cn setup`** and `cn init` now materialize cognitive assets automatically — hub is wake-ready out of the box
- **CAR design doc** (`docs/design/CAR-v3.4.md`) — problem statement, three-layer model, package format, resolution spec
- **Implementation plan** (`docs/gamma/plans/CAR-implementation-plan.md`) — 8 ordered steps

### Changed

- **`cn_context.ml`** — delegates mindset/skill loading to `Cn_assets`; fails fast if core assets missing
- **`cn_capabilities.ml`** — `render` accepts optional `~assets` summary
- **`cn setup`** no longer requires root — focuses on hub asset materialization (system-level setup moves to `--system`)

### Removed

- Inline `load_mindsets` and skill-walking code from `cn_context.ml` (moved to `cn_assets.ml`)

---

## v3.3.1 (2026-03-07)

**Agent Instruction Alignment**

Fixes Pi hallucinating XML-style tool syntax (`<observe><fs_read>...`) by ensuring the agent sees the correct `ops:` JSON format on every call.

### Added

- **Canonical ops examples in capabilities block** — `cn_capabilities.ml` now injects `syntax:`, `example_observe:`, and `example_effect:` lines so the model sees the exact frontmatter JSON format on every wake-up
- **Output Discipline section** in `AGENTS.md` — explicit ban on XML pseudo-tools, markdown tool-call blocks, and ad-hoc shell snippets
- **Typed Capability Ops section** in `agent-ops/SKILL.md` — full JSON examples for observe, effect, and two-pass patterns with "do not" list

### Fixed

- **Stale paths** in `AGENTS.md` — `skills/` → `src/agent/skills/`, `mindsets/` → `src/agent/mindsets/`, `mindsets/OPERATIONS.md` → `src/agent/mindsets/OPERATIONS.md`
- **Session contract** in `AGENTS.md` — split into runtime mode (context pre-packed, don't re-read) vs manual/debugging mode

### Changed

- **Capabilities block** conditionally omits `example_effect:` when `apply_mode=off` (observe-only)
- All capabilities tests updated for new example lines

---

## v3.3.0 (2026-03-06)

**CN Shell — Typed Ops, Two-Pass Execution, Path Sandbox**

The agent can now read files, inspect git state, write patches, and run allowlisted commands — all as governed, post-call typed ops with receipts and artifact hashing. The pure-pipe boundary is preserved: no in-call tools, no tool loop. Ops are proposed in output.md, validated by cn, and executed with full audit trail.

### Added

- **cn_shell.ml** — Typed op vocabulary (7 observe + 5 effect kinds) + manifest parser with phase validation, auto-assigned op_ids, and duplicate detection
- **cn_sandbox.ml** — Path sandbox: normalization, `..` collapse, symlink resolution via `realpath`, denylist on resolved canonical paths, protected file enforcement. Catches symlinked-parent bypass attacks
- **cn_executor.ml** — Op dispatcher for fs_read/list, git_status/diff/log/grep, fs_write/patch, git_branch/commit, exec (allowlisted + env-scrubbed). Produces receipts and SHA-256-hashed artifacts
- **cn_orchestrator.ml** — Two-pass execution engine: Pass A executes observe ops and defers effects; Pass B executes effects and denies new observe ops. Coordination op gating (terminal ops blocked on effect failure). Context repack with bounded receipts + artifact excerpts for Pass B
- **cn_projection.ml** — Idempotent projection markers using `O_CREAT|O_EXCL` for crash-recovery deduplication of outbound messages
- **cn_capabilities.ml** — Runtime capability discovery block injected into packed context so the agent knows what ops are available before proposing them
- **cn_dotenv.ml** — `.cn/secrets.env` loader with env-var-first resolution and `0600` permission enforcement
- **cn_sha256.ml** — Pure OCaml SHA-256 (FIPS 180-4), zero external dependencies
- **Telegram UX** — 🤔 reaction on inbound message + typing indicator; reaction cleared on success. Uses `request_once` (no retry, 3s cap) for cosmetic calls
- **Crash recovery** — `ops_done` checkpoint prevents duplicate side effects on retry; projection markers prevent duplicate Telegram sends; conversation dedup by trigger_id; ordered cleanup (state files before markers) for crash safety
- **175+ new tests** — ppx_expect tests across all new modules: dotenv (18), shell (28), sandbox (30), executor (28), orchestrator (29), projection (16), capabilities (12), SHA-256 (10)

### Changed

- **cn_config.ml** — Loads CN Shell settings from `.cn/config.json` runtime section: `two_pass`, `apply_mode`, `exec_enabled`, `exec_allowlist`, budgets
- **cn_context.ml** — Optional `~shell_config` parameter injects capabilities block into packed context
- **cn_runtime.ml** — Passes `shell_config` through pack and recovery paths; stale `ops_done` GC on idle State 3 entry

### Docs

- **AGENT-RUNTIME.md** — Appendix C with 6 normative worked examples (v3.3.6)
- **PLAN.md** — 7-step implementation plan (all steps complete)
- **README, ARCHITECTURE, CLI, AUTOMATION** — Updated for new modules and hub structure
- **BUILD-RELEASE.md** — Accurate 7-step release process with RELEASE.md support

### CI

- **ci.yml** — `TMPDIR` isolation + `-j 1` for ppx_expect temp file race
- **release.yml** — Same TMPDIR fix + conditional `RELEASE.md` body for GitHub Releases

---

## v3.2.0 (2026-03-04)

**Structured LLM Request Schema**

The LLM now receives a structured prompt instead of a flat string. This unlocks prompt caching, proper system prompts, and real multi-turn conversation.

### The Big Picture

```
Before (v3.0):  { messages: [{ role: "user", content: "<everything>" }] }
After  (v3.2):  { system: [stable_block, dynamic_block], messages: [turn, turn, ...] }
```

### Added

- **Structured system prompt** — Two system blocks: (1) stable identity/user/mindsets with `cache_control` for Anthropic prompt caching (~90% cache hits on repeat calls), (2) dynamic reflections/skills (refreshed each call)
- **Real multi-turn messages** — Conversation history as actual `user`/`assistant` turns instead of markdown rendered inside a single user message
- **`system_block` / `message_turn` types** in `cn_llm.ml` with manual `to_json` helpers (zero-dep, no ppx)
- **`packed` type** in `cn_context.ml` — returns structured data + `audit_text` for backward-compatible `input.md` logging
- **Design doc** — `docs/DESIGN-LLM-SCHEMA.md` covers problem, schema, OCaml serialization decision, and all tradeoffs
- **Mindsets in context packer** (v3.1.x) — `src/agent/mindsets/` loaded in deterministic order between USER and reflections
- **Role-weighted skill scoring** (v3.1.x) — `runtime.role` from `.cn/config.json` applies +2 bonus to matching skill paths
- **Setup installer design doc** — `docs/alpha/cli/SETUP-INSTALLER.md`: guided `cn setup` flow with Telegram auto-detection, persona picker, secrets management
- **UX-CLI terminal conventions** — Color-to-semantics map, `NO_COLOR` support, actionable error messages
- **New skills** — `agent/coherent`, `eng/coding`, `eng/functional`, `eng/testing`

### Changed

- `Cn_llm.call` signature: `~content:string` → `~system:system_block list ~messages:message_turn list`
- `Cn_context.pack` returns `packed` record (system blocks + messages + audit_text) instead of flat `content` string
- `Cn_runtime.ml` passes structured data to LLM; recovery path (State 2) re-packs from extracted inbound message
- Tests updated to use `packed.audit_text` instead of `packed.content`

### Fixed

- `Filename.temp_dir` replaced with OCaml 4.14-compatible helper
- `runtime.role` normalized to lowercase (tolerates "Engineer", "PM", etc.)
- Expect test promotion: tokenize order + JSON error messages

---

## v3.0.0 (2026-02-22)

**Native Agent Runtime**

cnos agents now run natively — no external orchestrator required. The runtime implements a pure-pipe architecture: the LLM is a stateless `string → string` function, `cn` handles all I/O, effects, and coordination. OpenClaw is fully removed.

### The Big Picture

```
Before (v2.x):  Telegram → OpenClaw → cn → agent ops
After  (v3.0):  Telegram → cn agent → Claude API → cn → agent ops
                           ^^^^^^^^^^^^^^^^^^^^^^^^
                           All native OCaml, ~1,500 lines
```

The agent runtime replaces OpenClaw's:
- Telegram bot infrastructure → `cn_telegram.ml`
- LLM orchestration → `cn_llm.ml` (single call, no tool loop)
- Context management → `cn_context.ml` (pack everything upfront)
- Session handling → `cn_runtime.ml` (orchestrator)

### Added

**6 new modules (~1,500 lines total):**

| Module | Lines | Purpose |
|--------|-------|---------|
| `cn_runtime.ml` | 557 | Orchestrator: pack → call LLM → write → archive → execute → project |
| `cn_json.ml` | 310 | Zero-dependency JSON parser/emitter |
| `cn_context.ml` | 177 | Context packer (SOUL, USER, skills, conversation, message) |
| `cn_telegram.ml` | 170 | Telegram Bot API client (long-poll + send) |
| `cn_llm.ml` | 154 | Claude Messages API client (single request/response) |
| `cn_config.ml` | 88 | Environment + config.json loader |

**CLI interface:**
- `cn agent` — Cron entry point (dequeue, pack, call LLM, execute ops)
- `cn agent --daemon` — Long-running Telegram polling loop
- `cn agent --stdio` — Interactive testing mode
- `cn agent --process` — Process one queued item directly

**Documentation:**
- `AGENT-RUNTIME.md` — Full design doc with rationale, architecture, migration plan
- Updated `ARCHITECTURE.md`, `CLI.md`, `AUTOMATION.md`

### Architecture: Pure Pipe

The core insight: LLMs don't need tools if you pack context upfront.

```
cn packs context → LLM plans → cn executes ops
```

**Context stuffing vs. tool loops:**

| Approach | Token cost | API calls | Latency |
|----------|------------|-----------|---------|
| Tool loop (3 retrievals) | ~15K+ | 3-5 | 3-10s |
| Context stuffing | ~6.5K | **1** | **1-3s** |

For cnos's bounded, predictable context (~20-30 hub files), stuffing wins.

**What gets packed into `state/input.md`:**
- `spec/SOUL.md` — Agent identity (~500 tokens)
- `spec/USER.md` — User context (~300 tokens)
- Last 3 daily reflections (~1,500 tokens)
- Matched skills, top 3 (~1,500 tokens)
- Conversation history, last 10 (~2,000 tokens)
- Inbound message (~200 tokens)

### Changed

- **Removed OpenClaw dependency** — No external orchestrator
- **System deps only** — Requires `git` + `curl`, no opam runtime deps
- **Config location** — `.cn/config.json` (reuses existing hub discovery)
- **Secrets via env vars** — `TELEGRAM_TOKEN`, `ANTHROPIC_KEY`, `CN_MODEL`

### Security Model (Preserved)

The CN security invariant is enforced:

> Agent interface: `state/input.md → state/output.md`  
> LLM reality: sees text in, produces text out. cn does all effects.

- LLM never touches files, commands, or network
- Ops parsed from output.md frontmatter, validated before execution
- Raw IO pairs archived before effects (crash-safe audit trail)
- API keys via env vars, never logged

### Breaking Changes

- **OpenClaw no longer required** — Remove OC config after migration
- **`cn agent` is the new entry point** — Replaces OC heartbeat/cron
- **Telegram handled natively** — OC bot infrastructure bypassed

### Migration

1. Set env vars: `TELEGRAM_TOKEN`, `ANTHROPIC_KEY`
2. Create `.cn/config.json` with `allowed_users`
3. Start daemon: `cn agent --daemon` (or systemd unit)
4. Add cron: `* * * * * cn agent` (processes queue)
5. Disable OpenClaw after verification

Rollback: `systemctl stop cn-<name> && systemctl start <previous-service>`

### Technical Highlights

- **Zero opam runtime deps** — stdlib + Unix only
- **curl-backed HTTP** — No OCaml HTTP stack complexity
- **Body consumption rules** — Markdown body = full response, frontmatter = notification
- **FSM integration** — Reuses existing `cn_protocol.ml` actor FSM
- **617 lines of new tests** — `cn_cmd_test.ml` + `cn_json_test.ml`

---

## v2.4.0 (2026-02-11)

**Typed FSM Protocol**

The cn protocol is now modeled as four typed finite state machines with compile-time safety.

### Added
- **4 typed FSMs** — Transport Sender, Transport Receiver, Thread Lifecycle, Actor Loop
- **cn_protocol.ml** — unified protocol implementation replacing scattered logic
- **385 lines of protocol tests** — ppx_expect exhaustive transition testing
- **ARCHITECTURE.md** — system overview documentation
- **PROTOCOL.md** — FSM specifications with ASCII diagrams
- **AUDIT.md** — design audit methodology

### Changed
- Archived 7 superseded design docs to `docs/design/_archive/`
- Unified `cn_io.ml` + `cn_mail.ml` transport logic into single FSM

### Fixed
- Invalid state transitions now caught at compile time, not runtime

---

## v2.3.1 (2026-02-11)

**Branch Cleanup Fix**

### Fixed
- Delete local branches after inbox sync (MCA cleanup)
- Skip self-update check if install dir missing
- CI: add ppx_expect to opam deps

---

## v2.3.0 (2026-02-11)

**Native OCaml Binary**

cn is now a native binary — no Node.js required.

### Added
- **Native OCaml build** — `dune build` produces standalone binary
- **Release workflow** — pre-built binaries for Linux/macOS
- **10-module refactor** — cn.ml split into focused modules:
  - `cn_agent.ml`, `cn_commands.ml`, `cn_ffi.ml`, `cn_fmt.ml`
  - `cn_gtd.ml`, `cn_hub.ml`, `cn_mail.ml`, `cn_mca.ml`
  - `cn_protocol.ml`, `cn_system.ml`

### Changed
- `cn update` now copies `bin/*` to `/usr/local/bin/` (2.2.28 backport)
- `cn-cron` now runs `cn in` after sync to wake agent (2.2.27 backport)

### Fixed
- Filter inbound to `threads/in/` only (2.2.26 backport)
- Native argv handling (drops 1, not 2 like Node.js)

---

## v2.2.0 (2026-02-07)

**First Hash Consensus**

Two AI agents achieved verified hash consensus via git-CN protocol — the actor model is complete.

### Milestone
- **Hash consensus** — Pi and Sigma independently converged on cnos `d1cb82c`
- **Verified via git** — runtime.md pushed to GitHub, human-verified
- **No central coordinator** — pure git-based decentralized coordination

### Added
- **cn update auto-commit** — runtime.md auto-commits and pushes (P1)
- **Output auto-reply** — output.md responses auto-flow to sender (CLP accepted)
- **MCA injection** — coherence check every 5th cycle
- **Sync dedup** — checks `_archived/` before materializing
- **agent-ops skill** — operational procedures for CN agents
- **CREDITS.md** — lineage acknowledgment

### Fixed
- Wake mechanism: `openclaw system event` instead of curl
- Version string sync between package.json and cn_lib.ml
- Stale branch cleanup procedures

### Protocol
- 5-minute cron interval standardized
- input.md/output.md protocol documented in SYSTEM.md
- Queue system (state/queue/) for ordered processing

---

## v2.0.0 (2026-02-05)

**Everything Through cn**

Paradigm shift: agents no longer run git directly. Everything goes through the `cn` CLI.

### Added
- **CN CLI v0.1** — `tools/src/cn/`: modular OCaml implementation
  - `cn init`, `cn status`, `cn inbox`, `cn peer`, `cn doctor`
  - Best-in-class patterns from npm, git, gh, cargo
- **UX-CLI spec** — `skills/ux-cli/SKILL.md`: terminal UX standard
  - Semantic color channels
  - Failure-first design
  - "No blame, no vibes"
- **SYSTEM.md** — `spec/SYSTEM.md`: definitive system specification
- **cn_actions library** — Unix-style atomic primitives
- **Coordination First** principle — unblock others before your own loop
- **Erlang model** — fire-and-forget, sender tracks

### Architecture
- Agent = brain (decisions only)
- cn = body (all execution)
- Agent purity enforced by design

### Breaking Changes
- Agents should use `cn` commands, not raw git
- Direct git execution deprecated (will be blocked in future)

---

## v1.8.0 (2026-02-05)

**Agent Purity + Process Maturity**

Core architectural principle established: Agent = brain (decisions only), cn = body (all execution).

### Added
- **CN Protocol** — `docs/design/CN-PROTOCOL.md`: action vocabulary, protocol enforcement rules
- **Inbox Architecture** — `docs/design/INBOX-ARCHITECTURE.md`: agent purity constraint, thread-as-interface
- **Engineering skills** — `skills/eng/` with coding, ocaml, design, review, rca
- **Ship skill** — `skills/ship/`: branch → review → merge → delete workflow
- **Audit skill** — `skills/audit/`: periodic health checks
- **Adhoc-thread skill** — `skills/adhoc-thread/`: when/how to create threads (β test)
- **THINKING mindset** — evidence-based reasoning, know vs guess
- **AGILE-PROCESS** — `docs/gamma/AGILE-PROCESS.md`: backlog to done workflow
- **THREADS-UNIFIED** — `docs/design/THREADS-UNIFIED.md`: backlog as threads, GTD everywhere

### Changed
- **PM.md** — "Only creator deletes branch" rule added
- **ENGINEERING.md** — "Always rebase before review" principle

### Fixed
- Branch deletion violation documented in RCA

### Process
- No self-merge enforced
- Rebase before review required
- Author deletes own branches after merge

---

## v1.7.0 (2026-02-05)

**Actor Model + GTD Inbox**

Major coordination upgrade. Agents now use Erlang-inspired actor model: your repo is your mailbox.

### Added
- **inbox tool** — replaces peer-sync. GTD triage with Do/Defer/Delegate/Delete
- **Actor model design** — `docs/design/ACTOR-MODEL-DESIGN.md`
- **RCA process** — `docs/beta/evidence/rca/` with template and first incident
- **FUNCTIONAL.md** — mindset for OCaml/FP patterns
- **PM.md** — product management mindset with user stories, no self-merge
- **FOUNDATIONS.md** — the coherence stack explained (C≡ → TSC → CTB → cnos)
- **APHORISMS.md** — collected wisdom ("Tokens for thinking. Electrons for clockwork.")
- **ROADMAP.md** — official project roadmap
- **GitHub Actions CI** — OCaml tests + native build

### Changed
- **docs/ reorganized** — whitepapers/design docs moved to `docs/design/`
- **Governance** — no self-merge rule: engineer writes → PM merges

### Deprecated
- **peer-sync** — use `inbox` instead

### Fixed
- 4-hour coordination failure (RCA documented, protocol fixed)
