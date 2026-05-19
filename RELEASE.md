# 3.78.0

## Outcome

Coherence delta: C_Σ **B−** (`α B+`, `β B+`, `γ C`) · **Level:** `L6`

cnos's last procedure-in-Go is now procedure-as-skill. The `agent/activate` skill is the single source of truth for agent self-activation at a hub (Kernel → CA skills → Persona → Operator → hub state → identity), and `cn activate`'s renderer reads it. Non-cn bodies (Claude Code on the web, Codex, Claude.ai with WebFetch) can fetch the same skill and self-bootstrap by the same procedure — the "I wake up incoherent by default" failure documented at cn-sigma is closed at the protocol layer.

## Why it matters

Before this release, activation was the one cnos behavior encoded only in Go (`src/go/internal/activate/activate.go::writePrompt`). Every other cnos procedure shipped as a skill artifact with frontmatter, triggers, governing question, and load semantics — activation was the exception. That exception broke the "everything is a skill" invariant for the procedure most foundational to body coherence: a body that cannot self-activate is a body that wakes up incoherent every session. This release restores the invariant by shipping the activation procedure as `src/packages/cnos.core/skills/agent/activate/SKILL.md` (485 lines, conforming to `skill/SKILL.md`) and evolving the Go renderer to parse the skill's §4.1 load-order block on every invocation. The skill is canonical on drift; an in-Go fallback preserves manifest-only and pre-`cn deps restore` behavior.

The release also closes a pre-merge gate gap (β/SKILL.md §pre-merge gate row 3): the row prescribed running "the cycle's own validator" but did not enumerate the validator set. Two contract validators that the cycle's surface ships (I5 SKILL.md frontmatter validation; the R5-activate kata) were not run on the merge tree, producing post-merge CI red on two wiring-class regressions. γ patched both regressions and tightened β/SKILL.md row 3 to enumerate the contract-validator set explicitly so future cycles don't repeat the application gap.

## Added

- **#379** — `src/packages/cnos.core/skills/agent/activate/SKILL.md`: new body-agnostic activation skill defining the six-item canonical load order (Kernel → CA skills → Persona → Operator → hub state → identity), the three-tier body-capability matrix (shell+git preferred, HTTP fetch only, no-fetch operator-injected), a paste-testable README router template hubs adopt verbatim with only the hub URL substituted, and disambiguation from `cnos.cdd/skills/cdd/activation/SKILL.md` (which activates a repo under CDD, not an agent at a hub).
- **#379** — `src/go/internal/activate/activate.go`: `writePrompt` no longer carries the section ordering as in-Go literals; it iterates a `readFirst []readFirstItem` slice produced by `loadActivateSkillOrdering`, which parses the skill's §4.1 marker-bounded block via `parseReadFirstOrderBlock`. An in-Go fallback (`canonicalReadFirstOrdering`) preserves manifest-only behavior when the skill is not yet vendored; the skill remains canonical on drift.
- **#379** — `src/go/internal/activate/activate_test.go`: `TestSkillIsSourceOfTruthForReadFirstOrder` is a two-phase edit test (write fixture skill with kernel-before-persona; assert; swap; assert swap propagates; `out1 != out2` coherence assertion). Companion tests for fallback behavior, parser happy-path, and parser no-markers. Twenty-nine tests in the package, all green.

## Changed

- **#379** — `## Read first` ordering in `cn activate` output changed from persona/operator/kernel/deps/refl to kernel/CA-skills/persona/operator/hub-state/identity per AC3 of the issue. Section headers (`## Read first`, `## Kernel`, `## Persona`, `## Operator`, etc.) are preserved, so consumers parsing by section header are unaffected; consumers doing string-match within `## Read first` see the new ordering. Pre-existing tests (`TestReadFirstSection_OrderedSigma`, `TestReadFirstSection_InitOnlyOrdered`) were updated to the new canonical order — appropriate, not silent.
- **`src/packages/cnos.cdd/skills/cdd/beta/SKILL.md`** §pre-merge gate row 3 — tightened from "the cycle's own validator (or any CI-equivalent the cycle ships)" to an enumerated list naming `validate-skill-frontmatter.sh`, `cn-cdd-verify`, `check-version-consistency.sh`, and the per-cycle kata under `src/packages/cnos.kata/katas/{N}/` (with R5-activate cited as the activate-surface kata exemplar). Empirical anchor (#379 post-merge CI red on I5 and R5-activate P10) named inline. (γ step 13a; lands in this release.)
- **`src/packages/cnos.kata/katas/R5-activate/{run.sh,kata.md}`** P10 — assertion updated to the new canonical `## Read first` ordering (`kernel < persona < operator < {deps, reflection}`); the deps/reflection relative order is no longer asserted because both share one hub-state line in the new ordering. The kata doc names the canonical source-of-truth (`src/packages/cnos.core/skills/agent/activate/SKILL.md §4.1`). (γ step 13a; lands in this release.)
- **`src/packages/cnos.core/skills/agent/activate/SKILL.md`** `calls:` frontmatter — paths rewritten from `cnos.core/...` prefix form to package-skill-root-relative form (`../doctrine/KERNEL.md`, `agent/cap/SKILL.md`, etc.) to satisfy the I5 validator's resolution rule. (γ step 13a; lands in this release.)

## Validation

- **Skill-as-source-of-truth (renderer-level):** `go test -count=1 -run TestSkillIsSourceOfTruth ./internal/activate/...` PASS on cycle HEAD and merge tree (β R1).
- **I5 SKILL.md frontmatter validator:** `./tools/validate-skill-frontmatter.sh` → "66 SKILL.md validated; no findings" after γ patch.
- **R5-activate kata:** expected green on the next CI run against a freshly-built `cn` binary (γ patch lands the assertion update; binary built from this release embeds the new renderer).
- **End-to-end body activation (manual dry-run):** deferred — no body session available during this release's PRA; commit on next σ session checklist.

## Known Issues

- Hub README router template ships in the skill (§2.3) but no hub has adopted it yet; per-hub README patches are deferred to per-hub cycles (next MCA: cn-sigma adoption).
- Renderer fallback `canonicalReadFirstOrdering` duplicates the skill's §4.1 ordering as in-Go constants. Mitigations: skill is preferred when vendored, fallback exercised only when bundle absent, fallback path tested, drift rule documented (skill canonical). Durable fix is `//go:embed` of the skill into the binary; deferred to a future cycle.
- I4 (Repo link validation) and `notify` (telegram bot token) CI workflows remain red on `main` from pre-cycle baseline; not regressed by #379; tracked separately.

## Source

Issue #379 (https://github.com/usurobor/cnos/issues/379). Cross-repo proposal lineage: `usurobor/cn-sigma:.cdd/iterations/cross-repo/cnos/agent-activate-skill/` (bundle commit `1a4e25f75`; AC4 capability-matrix refinement folded in from cn-sigma branch HEAD `bdda457f5`). Disposition: `accepted` with delta named in the issue body §Source Proposal. The `landed` event addition to the source-side STATUS is pending operator merge of the two mirror branches (`cnos:sigma/cross-repo-mirror-agent-activate-skill@212d5239` and `cn-sigma:sigma/cross-repo-status-lineage-update@89049611c`); δ to coordinate.
