<!--
section-manifest:
  planned: [gap, skills, acs, self-check, debt, cdd-trace, review-readiness]
  completed: [gap, skills, acs]
-->

# Self-coherence — α #379

## Gap

**Issue:** [#379 — agent/activate skill: single source of truth for agent self-activation](https://github.com/usurobor/cnos/issues/379).
**Mode:** design-and-build. The design surface is the issue body (Source of Truth table, AC3 load order, AC4 capability matrix, AC7 renderer-interpolation contract). α executes design and build in the same cycle — no separate committed `DESIGN.md`.
**Version:** unreleased (cycle directory `.cdd/unreleased/379/`); release boundary follows β's merge and γ's cycle close.
**Base SHA:** `7a7f7152af649ea93cebe5f909fbf1397d809547` (γ scaffold base = origin/main HEAD at dispatch time; cycle branch rebased onto current `origin/main` — confirmed in §Review-readiness).

**Named gap.** cnos exposes activation as procedure-in-Go: `src/go/internal/activate/activate.go` `writePrompt` hardcodes the section ordering and the "Read first" list in in-Go literals. Every other cnos behavior is a skill (`cnos.core/skills/agent/{cap,clp,mca,mci,coherent,agent-ops}/SKILL.md` is the precedent). Two consequences follow: (1) the "everything is a skill" invariant is violated for activation; (2) non-cn bodies (Claude Code on the web, Codex sessions, Claude.ai with WebFetch) cannot reach the procedure because it does not exist as a fetchable artifact — they fall into the "I wake up incoherent by default" failure named in cn-sigma `threads/adhoc/20260325-session2-learnings.md` §1.

**Closure shape.** This cycle creates `src/packages/cnos.core/skills/agent/activate/SKILL.md` as the single source of truth for agent activation, and evolves `src/go/internal/activate/activate.go` to render from it. The skill prescribes a six-item canonical load order (Kernel → CA skills → Persona → Operator → hub state → identity confirmation), a three-tier body-capability matrix (shell+git preferred, HTTP fetch only, no-fetch operator-injected), a README router template hubs adopt verbatim, and a disambiguation paragraph distinguishing this skill from `cnos.cdd/skills/cdd/activation/SKILL.md` (CDD repo activation, a distinct concern). The renderer parses the skill's §4.1 machine-readable load-order block at every invocation and emits in that order; an in-Go fallback preserves manifest-only and no-deps behavior when the skill is not yet vendored.

**Non-goals (held).** No hub README patches (cn-sigma, cn-pi) — deferred to per-hub cycles. No `cnos.xyz` rendering service. No `cn doctor` validation of router presence. No edits to `cnos.cdd/skills/cdd/activation/SKILL.md` beyond citation in §2.4 disambiguation. No changes to the `cn activate` CLI surface (flags, arguments, exit codes — `src/go/internal/cli/cmd_activate.go` unchanged).

## Skills

**Tier 1 — always-on for α/CDD:**

- `src/packages/cnos.cdd/skills/cdd/CDD.md` — canonical algorithm and lifecycle authority. §1.4 large-file authoring rule applied: SKILL.md authored section-by-section with HTML-comment manifest; self-coherence.md authored section-by-section per α/SKILL.md §2.5.
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface. §2.2 artifact order followed (skill artifact → tests → code → docs → self-coherence); §2.5 incremental self-coherence; §2.6 pre-review gate exercised in §Review-readiness; §2.6 row 14 retroactive author-amend executed (worktree config had γ identity from session start; rebase --exec --reset-author applied across α's six commits before signaling).

**Tier 2 — always-applicable engineering:**

- `src/packages/cnos.eng/skills/eng/go/SKILL.md` — Go authoring conventions. Applied to renderer evolution (`activate.go`): small functions, explicit fallback path, error-free reads on optional skill, line-oriented parser instead of regex.
- `src/packages/cnos.core/skills/write/SKILL.md` — prose discipline for the SKILL.md artifact (clarity over decoration; rule of three: principle → unfold → kata).

**Tier 3 — issue-specific (per issue body §Skills to load + γ scaffold §Tier 3 list):**

- `src/packages/cnos.core/skills/skill/SKILL.md` — load-bearing for AC2 conformance. Frontmatter keys (`name`, `description`, `artifact_class: skill`, `kata_surface`, `governing_question`, `triggers`, `scope`, `inputs`, `outputs`, `parent`, `visibility`, `requires`, `calls`) all present; body follows Define → Unfold → Rules → Verify → Final test → Kata progression; embedded kata at §7 satisfies `kata_surface: embedded`.
- `src/packages/cnos.core/skills/write/SKILL.md` — prose structure.
- `src/packages/cnos.core/skills/design/SKILL.md` — design-and-build mode design step (substrate independence between cnos and per-hub layers; single source of truth for the activation procedure; defer realization choices about which body capability is used).
- `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` — cited and disambiguated in SKILL.md §2.4 (AC6). Not modified — non-goal per γ scaffold failure-mode 10.
- `src/packages/cnos.eng/skills/eng/go/SKILL.md` — renderer evolution (AC7).
- `src/packages/cnos.cdd/skills/cdd/issue/SKILL.md` — form authority. No in-cycle issue-pack reconciliation was required; the issue body remained authoritative throughout.

## ACs

Per-AC oracle results recorded inline. AC numbering follows the issue body.

### AC1 — activate skill exists at canonical path

**Invariant.** `src/packages/cnos.core/skills/agent/activate/SKILL.md` exists, is non-empty, and is ≥200 lines of substantive content.

**Oracle.** `test -f src/packages/cnos.core/skills/agent/activate/SKILL.md && [ -s src/packages/cnos.core/skills/agent/activate/SKILL.md ] && wc -l src/packages/cnos.core/skills/agent/activate/SKILL.md`

**Result.**
```
485 src/packages/cnos.core/skills/agent/activate/SKILL.md
```
File present at exact path, non-empty, 485 lines (well above 200). Sibling to `agent/{cap,clp,mca,mci,coherent,agent-ops,...}/SKILL.md`. **PASS.**

### AC2 — Skill conforms to skill/SKILL.md format

**Invariant.** Frontmatter contains `name`, `description`, `artifact_class: skill`, `governing_question`, `triggers`, `scope`, `inputs`, `outputs` (plus the skill skill's `kata_surface` requirement). Body follows the Define → Unfold → Rules → Verify progression with an embedded kata surface.

**Oracle.** Frontmatter parsed via awk; required keys grep'd; body section structure inspected.

**Result.**
```
=== AC2 oracle: frontmatter keys present ===
name: activate
description: Activate an agent identity at a cnos hub. Body-agnostic procedure for loading Kernel + CA skills + Persona + Operator + hub state, in that order, …
artifact_class: skill
kata_surface: embedded
governing_question: How does an AI body, knowing only a hub URL, reach a state where it can name its identity, its operator, and its current orientation without operator improvisation?
triggers: [activate, activation, bootstrap, self-activation, wake up, cn activate]
scope: task-local
inputs: [hub URL or hub path, body capabilities, access to this skill's content]
outputs: [loaded Kernel + CA skills + Persona + Operator + hub state surveyed + identity named]
visibility: public
parent: agent
requires: [hub URL or hub path, fetch capability or operator injection]
calls: [KERNEL.md + cap/clp/mca/mci/coherent/agent-ops SKILL.md]
```

Body sections present in order: `## Core Principle`, `## Algorithm`, `## 1. Define` (§1.1 parts, §1.2 fit, §1.3 failure mode), `## 2. Unfold` (§2.1 load order, §2.2 capability matrix, §2.3 README router template, §2.4 disambiguation), `## 3. Rules` (3.1–3.7 imperative rules with ❌/✅ pairs), `## 4. Renderer contract` (§4.1 machine-readable block, §4.2 interpolation surface, §4.3 fallback, §4.4 observable-output preservation), `## 5. Verify` (5.1–5.6 checks), `## 6. Failure modes catalogue`, `## 7. Embedded kata` (scenario / task / inputs / artifacts / procedure / verification / common failures / reflection — satisfying skill/SKILL.md §2.4 kata surface contract), `## 8. References`. **PASS.**

### AC3 — Skill prescribes the body-agnostic six-item load order

**Invariant.** SKILL.md describes the sequence as (1) Kernel from cnos, (2) CA skills from cnos, (3) Persona from hub, (4) Operator from hub, (5) hub state (deps, latest reflection, memory, threads), (6) identity confirmation. The soul/identity layering rule is cited.

**Oracle.** awk-extract of §2.1 numbered list; awk-extract of §4.1 machine-readable block; grep for layering-rule citation.

**Result — §2.1 ordered list:**
```
1. **Kernel from cnos.** …
2. **CA skills from cnos.** …
3. **Persona from hub.** …
4. **Operator from hub.** …
5. **Hub state.** …
6. **Identity confirmation.** …
```

**Result — §4.1 machine-readable block:**
```
<!-- read-first-order:begin -->
1. kernel — Kernel doctrine (what kind of agent)
2. ca-skills — CA skill set (behavioral instructions)
3. persona — Persona (who you are at this hub)
4. operator — Operator (whom you serve at this hub)
5. hub-state — Hub state (deps, latest reflection, memory, threads)
6. identity — Identity confirmation
<!-- read-first-order:end -->
```

**Result — layering rule citation (§2.1 layering paragraph):**
> _The layering rule is from cn-sigma `threads/adhoc/20260325-session2-learnings.md` §3: "Soul = what kind of agent. Identity = which agent. Don't mix them."_

Both the human-readable list (§2.1) and the machine-readable block (§4.1) are in the canonical order. §5.2 of the skill names the §2.1/§4.1 coherence rule. **PASS.**

### AC4 — Body-agnostic three-tier capability matrix

**Invariant.** SKILL.md enumerates at least three body-capability tiers, each with its load path. Tier (a) shell + git is named preferred when available. The no-fetch degraded path is named honestly.

**Oracle.** grep for tier labels; inspect §2.2 table rows; confirm preferred-when-available statement.

**Result.**
```
§2.2 §157: | (a) shell + git (preferred when available) | Body can execute shell commands and `git clone` | git clone https://github.com/usurobor/cnos.git, then read files from the local checkout | git clone <hub-url> … | One network operation per repo; atomic; later reads are filesystem-fast |
§2.2 §158: | (b) HTTP fetch only | Body can fetch raw URLs but cannot execute shell | Fetch raw GitHub URLs per file, e.g. https://raw.githubusercontent.com/usurobor/cnos/main/src/packages/cnos.core/doctrine/KERNEL.md | …per-file… | One network operation per file; non-atomic; may hit rate limits |
§2.2 §159: | (c) no fetch | Body has no shell, no git, no HTTP fetch | The body cannot self-activate. The operator must paste the bootstrap contents into the body's prompt window directly. | The operator must inject. | Operator labor per session; degraded path |

§2.2 §161: **Tier (a) shell + git is preferred when available.** Reasons: atomic / local / inspectable / honest about reality (most modern bodies have shell+git).
§2.2 §170: **Tier (c) no fetch is the degraded path, named honestly.** A body with no shell, no git, and no HTTP fetch cannot self-activate. Activation in this tier requires the operator to inject the bootstrap content directly…
```

Three tiers explicit, each names its load path, tier (a) marked preferred with reasons, tier (c) named honestly as degraded with operator-injection requirement. Tier-selection rule states bodies MUST observe their own capability and pick the matching tier. **PASS.**

### AC5 — README router template documented in the skill

**Invariant.** SKILL.md includes a labeled "README router template" section with a fenced markdown block hubs adopt verbatim, requiring only the hub URL substitution.

**Oracle.** grep for heading; inspect fenced block; copy block into /tmp and confirm only one placeholder.

**Result.**
```
§2.3 §182: ### 2.3. README router template
```

The block uses `~~~markdown` fences (so the inner ``` blocks pass through cleanly), is self-contained, and contains exactly one placeholder: `<HUB-URL>`. Paste-test executed:

```sh
$ awk '/^~~~markdown$/,/^~~~$/' src/packages/cnos.core/skills/agent/activate/SKILL.md | grep -c '<HUB-URL>'
3   # three occurrences of <HUB-URL>, no other placeholders
$ awk '/^~~~markdown$/,/^~~~$/' src/packages/cnos.core/skills/agent/activate/SKILL.md | grep -E '<[A-Z_-]+>' | sort -u
<HUB-URL>
```

The skill's §2.3 adoption-rule paragraph states explicitly: "The only per-hub edit is the `<HUB-URL>` substitution." The template names tier (a)/(b)/(c) paths in the same order as §2.2. **PASS.**

### AC6 — Naming disambiguation from cdd/activation/SKILL.md

**Invariant.** SKILL.md cites `src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` and `src/packages/cnos.core/skills/agent/activate/SKILL.md`, exact paths, in the same paragraph or section, with a one-sentence distinction.

**Oracle.** `rg` both paths inside the skill; confirm both appear in §2.4 with the one-sentence distinction.

**Result.**
```
$ rg 'cnos\.cdd/skills/cdd/activation/SKILL\.md' src/packages/cnos.core/skills/agent/activate/SKILL.md | wc -l
4
$ rg 'cnos\.core/skills/agent/activate/SKILL\.md' src/packages/cnos.core/skills/agent/activate/SKILL.md | wc -l
7
```

§2.4 paragraph (excerpt):
> _`src/packages/cnos.cdd/skills/cdd/activation/SKILL.md` is **CDD repo activation** — the one-time bootstrap sequence that brings an existing repository under the Coherence-Driven Development protocol … `src/packages/cnos.core/skills/agent/activate/SKILL.md` (this skill) is **agent activation at a hub** — the procedure an AI body follows on every wake-up to load Kernel, CA skills, Persona, Operator, and hub state, and to confirm its identity at the named hub. … They share a word, not a concern: one activates a repo under CDD; the other activates a body at a hub._

Both paths exact (leaf `activate` for core; leaf `activation` for cdd — the common typo guarded against in γ scaffold failure-mode 5). One-sentence distinction present. **PASS.**

### AC7 — `cn activate` reads the skill instead of hardcoding the procedure

**Invariant.** `src/go/internal/activate/activate.go` `writePrompt` no longer sources its section ordering / "Read first" list from in-Go literals. It reads the activate skill from the vendored bundle and parses §4.1; a fallback preserves manifest-only deps behavior when the skill is not vendored. `activate_test.go` contains a test that demonstrates the skill is the source of truth.

**Oracle.**
1. `cd src/go && go test ./internal/activate/... -race` → must pass.
2. Inspect `writePrompt`: confirm no hardcoded ordering literals.
3. Inspect `activate_test.go`: confirm a test asserts that editing the skill's ordering changes the renderer's output.

**Result — go test:**
```
$ cd src/go && go test -race ./internal/activate/...
ok  	github.com/usurobor/cnos/src/go/internal/activate	1.132s

$ cd src/go && go test ./...   # full src/go package suite — no regression elsewhere
ok  	github.com/usurobor/cnos/src/go/internal/activate
ok  	github.com/usurobor/cnos/src/go/internal/activation
ok  	github.com/usurobor/cnos/src/go/internal/binupdate
ok  	github.com/usurobor/cnos/src/go/internal/cli
ok  	github.com/usurobor/cnos/src/go/internal/discover
ok  	github.com/usurobor/cnos/src/go/internal/dispatch
ok  	github.com/usurobor/cnos/src/go/internal/doctor
ok  	github.com/usurobor/cnos/src/go/internal/hubinit
ok  	github.com/usurobor/cnos/src/go/internal/hubsetup
ok  	github.com/usurobor/cnos/src/go/internal/hubstatus
ok  	github.com/usurobor/cnos/src/go/internal/pkg
ok  	github.com/usurobor/cnos/src/go/internal/pkgbuild
ok  	github.com/usurobor/cnos/src/go/internal/restore

$ go vet ./internal/activate/...
(no output — clean)
```

**Result — writePrompt no longer hardcodes ordering.** The new signature is `writePrompt(w, absPath, cfg, kernel, persona, operator, deps, latest, memory, threads, readFirst []readFirstItem)`. The `## Read first` block is emitted by iterating `readFirst` (parsed from the skill or canonical fallback):
```go
fmt.Fprintf(w, "## Read first\n")
for i, item := range readFirst {
    fmt.Fprintf(w, "%d. %s\n", i+1, renderReadFirstLine(item, kernel, persona, operator, deps, latest))
}
```
The renderer's `loadActivateSkillOrdering` reads from `<cnDir>/vendor/packages/cnos.core/skills/agent/activate/SKILL.md` and parses §4.1 via `parseReadFirstOrderBlock`. The fallback `canonicalReadFirstOrdering` mirrors §4.1 and is exercised when the skill is not vendored; the skill remains authoritative on drift.

**Result — skill-as-source-of-truth test.** `TestSkillIsSourceOfTruthForReadFirstOrder` in `activate_test.go`:
- Phase 1: writes a vendored fixture skill with `kernel` before `persona`; runs `Run`; asserts `KERNEL.md < PERSONA.md` in `## Read first`.
- Phase 2: rewrites the same fixture skill to swap `kernel` and `persona`; reruns `Run`; asserts `PERSONA.md < KERNEL.md`.
- Coherence assertion: `out1 != out2` (the two skill orderings produce different renderer outputs).

Companion tests: `TestSkillFallback_NotVendored` confirms the manifest-only fallback path emits the canonical ordering, surfaces the stderr diagnostic ("activate skill not vendored"), and preserves `run cn deps restore` guidance. `TestParseReadFirstOrderBlock_HappyPath` / `_NoMarkers` unit-test the parser.

**Result — observable-output deltas (documented per AC7 surface-containment rule).**
The §4.1 canonical ordering changes the order of items in `cn activate`'s `## Read first` block:
- **Before:** persona → operator → kernel → deps → latest reflection.
- **After:** kernel → ca-skills → persona → operator → hub-state (deps + latest reflection composite) → identity confirmation.

Other observable surfaces preserved:
- `## Read first`, `## Kernel`, `## Persona`, `## Operator`, `## Dependencies`, `## Config summary`, `## Memory and reflection`, `## Inbox and threads`, `## Notes` — all section headers retained in their original order in `writePrompt`.
- Manifest-only kernel guidance ("dependency manifest declares cnos.core; not restored — run cn deps restore") preserved.
- Manifest-only deps guidance ("dependency manifest present; packages not restored — run cn deps restore") preserved.
- `## Notes` claude -p example unchanged.
- Stderr-only diagnostics preserved; stdout still contains only the prompt.

**PASS.**

### Surface containment

```
$ git diff --stat origin/main..HEAD
 .cdd/unreleased/379/gamma-scaffold.md              | 183 ++++++++
 .cdd/unreleased/379/self-coherence.md              | (this file)
 src/go/internal/activate/activate.go               | 211 ++++++--
 src/go/internal/activate/activate_test.go          | 216 ++++++++--
 src/packages/cnos.core/skills/agent/activate/SKILL.md | 485 +++++++
```

Five files changed, all expected per γ scaffold failure-mode 10. No edits to `cnos.cdd/skills/cdd/activation/SKILL.md`. No hub README patches. No `cn activate` CLI surface change (`src/go/internal/cli/cmd_activate.go` untouched).
