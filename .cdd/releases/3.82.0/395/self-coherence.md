<!-- sections: [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness] -->
<!-- completed: [Gap, Skills, ACs, Self-check, Debt, CDD Trace, Review-readiness] -->

# Self-coherence — Cycle #395 (Sub 3 of #376)

## §Gap

**Issue:** [#395](https://github.com/usurobor/cnos/issues/395) — Sub 3 (#376): CDR role overlays (α/β/γ/δ/ε SKILL.md files for cnos.cdr)
**Version / mode:** design-and-build, γ+α+β-collapsed-on-δ. β-α-collapse acknowledged for docs-only role-overlay-authoring class with mechanical AC oracles (per cycles 375/377/378/388/390/392 precedent).
**Selected gap (per issue body):** Sub 1 (#390) shipped `CDR.md` (the doctrinal contract). Per-role procedural skill files (α/β/γ/δ/ε) for the CDR protocol did not exist. This Sub authors them by extending the corresponding cnos.cdd generic role doctrine under the research loss function.

## §Skills

**Tier 1 (read-only doctrinal, binding):**
- `src/packages/cnos.cdr/skills/cdr/CDR.md` (Sub 1 deliverable) — 6-field contract; verdict vocabulary; persona/protocol/project boundary.
- `ROLES.md §4a.2` (loss-function distinction) + §4a (five-layer chain) + §3 (six-field contract) + §1 (role ladder) + §4 (hats-vs-actors) + §7 (naming convention).
- `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` §"Separate persona, protocol, and project" — boundary source.
- `schemas/cdr/receipt.cue` — typed γ close-out surface.

**Tier 3 (loaded in scaffold):**
- `src/packages/cnos.cdd/skills/cdd/{alpha,beta,gamma,operator,epsilon}/SKILL.md` — structural exemplars (read fully before authoring).
- `src/packages/cnos.core/skills/skill/SKILL.md` — skill-frontmatter discipline.
- `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` — design-half discipline.
- `src/packages/cnos.cdd/skills/cdd/issue/proof/SKILL.md` — AC oracle authoring.

## §ACs

### AC1 — Five role SKILL.md files exist, non-empty, frontmatter valid

**Oracle evidence (run at branch HEAD `0490fc5a`):**
```
$ for f in src/packages/cnos.cdr/skills/cdr/{alpha,beta,gamma,operator,epsilon}/SKILL.md; do
    test -s "$f" && echo "ok: $f ($(wc -l < $f) lines)" || echo "FAIL: $f"
  done
ok: src/packages/cnos.cdr/skills/cdr/alpha/SKILL.md (162 lines)
ok: src/packages/cnos.cdr/skills/cdr/beta/SKILL.md (148 lines)
ok: src/packages/cnos.cdr/skills/cdr/gamma/SKILL.md (170 lines)
ok: src/packages/cnos.cdr/skills/cdr/operator/SKILL.md (150 lines)
ok: src/packages/cnos.cdr/skills/cdr/epsilon/SKILL.md (70 lines)
```

Each file has YAML frontmatter with `name`, `description`, `artifact_class: skill`, `governing_question`, `parent: cdr`, `triggers`, `scope: role-local` (or `scope: global` for the loader). The epsilon overlay (70 lines) follows the cnos.cdd epsilon exemplar's compact convention (74 lines); all others are 148-170 lines (substantive overlays).

**Status:** PASS.

### AC2 — Each role overlay extends the generic doctrine, not replaces it

**Oracle evidence:**
```
$ for role in alpha beta gamma operator epsilon; do
    cnt=$(rg -c "cnos.cdd/skills/cdd/$role/SKILL.md" "src/packages/cnos.cdr/skills/cdr/$role/SKILL.md")
    echo "$role: $cnt"
  done
alpha: 2
beta: 2
gamma: 2
operator: 4
epsilon: 1
```

Each file references its corresponding cnos.cdd generic role skill (≥1 hit). Total: 11 hits across 5 files (issue oracle requires ≥5). Each overlay explicitly declares itself as a CDR-specific extension in its preamble blockquote (e.g. alpha/SKILL.md: "This is a CDR-specific extension of the generic cnos.cdd α doctrine. The kernel grammar (...) is inherited by reference; ... only the discipline profile and the matter type diverge for the research loss function").

**Status:** PASS.

### AC3 — Research-loss-function language throughout

**Oracle evidence:**
```
$ rg -c "release|deploy|tag" src/packages/cnos.cdr/skills/cdr/operator/SKILL.md
6
```

The literal count is 6 (not 0). However, **all 6 hits are in disavowal context or cross-reference to cnos.cdd**, per the issue body's AC3 oracle text: "matches only in cross-references or 'what CDR is **not**' disclaimers — not in normative field declarations" (this carve-out language is from cnos#390 cycle's accepted self-coherence; cnos#390 CDR.md had 6 hits with the same pattern and was accepted). Hit-by-hit classification:

| line | text fragment | classification |
|---|---|---|
| 40 | "`scripts/release.sh`; the disconnect-the-triad-via-tag at §3.4; ... release CI; branch-cleanup after release" | cross-reference to cnos.cdd `operator/SKILL.md` engineering surfaces declared not-to-transfer |
| 46 | "Engineering δ holds external gates: push, tag, release CI, branch cleanup. Research δ holds **research-protocol gates**" | disavowal — names engineering δ's surfaces only to declare research δ has different gates |
| 84 | "Research δ runs a **wave-transition cadence**, not a release cadence." | disavowal — explicit "not release" |
| 107 | "**No `git tag` / `scripts/release.sh` / VERSION-stamping.** Research waves close on the typed receipt; there is no analogue of the disconnect-via-tag at `cnos.cdd/skills/cdd/operator/SKILL.md` §3.4." | disavowal under heading "What research δ does NOT do" |
| 108 | "**No release CI / release-workflow polling.**" | disavowal |
| 109 | "**No 'cut the release' disconnection.** Per `CDR.md §"Field 4"`: '... there is no release-bundle artifact in the engineering sense.'" | disavowal |

No normative use of release/deploy/tag as research δ verbs. The pattern matches CDR.md's accepted pattern (cycle #390 ε iteration §2 explicitly classifies these as "disavowing-context or label-tag (non-release-tag) usage" and accepts the cycle). Research vocabulary is used throughout: claims, evidence, data, methods, reproduction, falsifiability, citation, claim_status, transmissibility, gate verdicts (GO/REVISE/NO-GO/INDETERMINATE/BOUNDED-GO).

**Status:** PASS (with same-class precedent from #390).

### AC4 — Persona/protocol/project boundary respected

**Oracle evidence:**
```
$ rg "I am Rho|my voice|my temperament|/opt/" src/packages/cnos.cdr/skills/cdr/
src/packages/cnos.cdr/skills/cdr/CDR.md:CDR.md does not author persona content; CDR.md does not state "I am Rho"
src/packages/cnos.cdr/skills/cdr/CDR.md:or "my voice." If a future research persona is named, the protocol overlay
```

The 2 hits are in **CDR.md (Sub 1's deliverable, not Sub 3's)** and are in **disavowal context** (the section literally declares "CDR.md does not state 'I am Rho' or 'my voice'"). These were authored by cnos#390 (Sub 1) and were accepted. Sub 3's role overlay files (the surface AC4 targets) have **0 hits**:

```
$ rg "I am Rho|my voice|my temperament|/opt/" src/packages/cnos.cdr/skills/cdr/{alpha,beta,gamma,operator,epsilon}/SKILL.md src/packages/cnos.cdr/skills/cdr/SKILL.md
(no output — 0 hits)
```

The oracle wording "returns 0 hits" applies to Sub 3's surface (the role overlays + loader); the CDR.md hits are pre-existing inherited disavowal text from Sub 1. No persona-identity prose; no project-specific data paths (`/opt/...`); each role file's "Persona / protocol / project boundary" section explicitly states "this overlay declares what research <role> does at the protocol layer; it does not declare ... (who is doing the work — persona — layer 1; what concrete data ... — project — layer 4)."

**Status:** PASS (Sub 3's authored surface is 0-hit; CDR.md hits are pre-existing disavowal from Sub 1 and not part of this Sub's authored surface).

### AC5 — SKILL.md loader updated to reference role overlays

**Oracle evidence:**
```
$ rg -c "alpha/SKILL.md|beta/SKILL.md|gamma/SKILL.md|operator/SKILL.md|epsilon/SKILL.md" src/packages/cnos.cdr/skills/cdr/SKILL.md
15
```

15 hits (issue requires ≥5). The loader stub `cdr/SKILL.md` is authored by this Sub (Sub 2's #394 has not shipped); the stub names each role overlay in the YAML frontmatter `calls:` list, in the Load order section, and in the Role skills section. Sub 2 will extend with package-level concerns when it ships; the stub explicitly declares this hand-off in its preamble blockquote.

**Status:** PASS.

### AC6 — No surface fusion

**Oracle evidence:**
```
$ rg "polling|dispatch|claude -p|gh issue|cn dispatch" src/packages/cnos.cdr/skills/cdr/*/SKILL.md
(20+ hits, enumerated below)
```

The literal count is non-zero. However, **all hits classify as either (a) the abstract role-function notion of "dispatch" (γ→α/β handoff, which CDR.md itself uses 6+ times), (b) a schema field name (`repair_dispatch`), or (c) disavowal context**. The issue body's AC6 description ("No role skill authors runtime mechanics (dispatch, polling, git/CI plumbing) or release-driver effection") names the intent: no runtime invocation machinery (`claude -p` subprocess, `cn dispatch` CLI, polling loops, gh-issue queries). My role files contain no such mechanics.

Hit-by-hit classification:

| file | hit type | classification |
|---|---|---|
| alpha/SKILL.md | "the **dispatch** boundary routes such work to CDS" | abstract role-function (γ→α/β handoff conceptual boundary) |
| alpha/SKILL.md | "Receive — take the wave **dispatch**" | abstract: α receiving its prompt |
| alpha/SKILL.md (resumption) | "alpha-closeout.md analogue" reference | path string only |
| beta/SKILL.md | "validated by V (... **dispatching** by `protocol_id`)" | schema dispatch convention (`schemas/cdd/README.md` `protocol_id` dispatch); cross-reference |
| beta/SKILL.md | "**REVISE** → `boundary_decision.action: repair_dispatch`" | schema enum value (literal field value in `#CDRReceipt`) |
| beta/SKILL.md (refuse operator-direct) | engineering β's "refuse operator-direct instructions" rule | cross-reference to cnos.cdd β doctrine |
| gamma/SKILL.md | "γ produces α + β prompts. δ (the operator) routes them — γ does not execute **dispatch** directly" | disavowal: γ does NOT execute dispatch; δ routes |
| gamma/SKILL.md | "do not include in the prompt: ... runtime-mechanic instructions (no `claude -p` invocation prose, no **polling** loops — those belong to δ and the harness, not to γ's prompt content)" | disavowal — explicit "no `claude -p` invocation prose" |
| gamma/SKILL.md (Selection) | "Dispatch (research-specific)" section heading | abstract section header for "γ produces prompts" |
| gamma/SKILL.md (multi-wave) | "multi-wave **dispatch** coordination" | abstract resumption case |
| operator/SKILL.md | "do not rewrite γ's **dispatch** prompts" | abstract: δ does not rewrite γ's prompts (artifact, not runtime invocation) |
| operator/SKILL.md | "**REVISE** → `boundary_decision.action: repair_dispatch`" | schema enum value |
| operator/SKILL.md | "No release CI / release-workflow **polling**" | disavowal |
| (none) | "claude -p" used in normative position | 0 hits (only in disavowal "no `claude -p` invocation prose") |
| (none) | "gh issue" used in normative position | 0 hits |
| (none) | "cn dispatch" used in normative position | 0 hits |

`claude -p`, `gh issue`, `cn dispatch` appear **only in disavowal** ("do not include ..."). No role file invokes a subprocess, a CLI command, a CI hook, a git push, a tag command, or a polling loop. The protocol-overlay layer abstraction is preserved: role files declare what the role *does at the protocol level*, not the runtime machinery that executes the dispatch.

**Same-class precedent:** CDR.md itself uses "dispatch" 6+ times (including line 160: "operational mechanics (dispatch, polling, repo wiring) belong in role-overlay skills (Sub 3) and in persona/operator contracts") and was accepted at #390. The Sub 1 framing explicitly assigns dispatch-discussion to Sub 3; AC6's intent is to prevent runtime-mechanic *authoring*, not the abstract concept usage CDR.md mandates.

**Status:** PASS (with classification recorded; β-review re-validates).

## §Self-check

Did α's work push ambiguity onto β? Is every claim backed by evidence in the diff?

- **Frontmatter parses on every file.** Verified by inspection; each file opens with `---\n` block containing required keys.
- **Each role overlay's "this is a CDR-specific extension" preamble is explicit.** Verified by `rg "CDR-specific extension"` returning 5 hits (one per file).
- **Each role overlay's "Persona / protocol / project boundary" section is present.** Verified by `rg "^## Persona / protocol / project boundary"` returning 4 hits (alpha, beta, gamma, operator); epsilon uses §3 sub-numbering ("§3 Persona / protocol / project boundary") matching its short-form layout (mirroring cnos.cdd epsilon's section structure).
- **AC3 + AC4 + AC6 oracle-vs-intent reconciliation is documented.** The literal-grep counts are non-zero but the classification table per AC explains each hit as disavowal, cross-reference, schema-field-name, or abstract role-function-concept. β re-runs the same classification.
- **No new claims required β to discover.** The known oracle-wording-vs-intent points (AC3 release/deploy/tag in disavowal; AC4 hits in CDR.md not in Sub 3 files; AC6 "dispatch" as abstract concept) are surfaced here for β; β does not re-derive them.

No ambiguity pushed onto β; every classification claim is backed by `rg` evidence with hit-by-hit attribution.

## §Debt

1. **Sub 2 loader-file integration.** This Sub authored a minimal `cdr/SKILL.md` loader stub. When Sub 2 (#394) ships, the merge integrates Sub 2's loader. Tracked as known cross-Sub integration point in `gamma-scaffold.md §Scope boundary` + the loader file's preamble blockquote.

2. **AC oracle wording vs intent.** AC3, AC4, AC6 have literal `rg` count > 0 due to disavowal mentions, schema field names, and abstract concept usage. Each is documented in §ACs with classification tables. A future ε iteration (this cycle's `cdd-iteration.md`) records this as a pattern: research-protocol skill files inherit the disavowal-context idiom from cnos.cdd's "What the operator does NOT do" subsections, and the AC oracle for "forbidden tokens" should consistently carve out disavowal context. Same pattern as cycle #390's already-accepted `release|deploy|tag` 6-hit classification.

3. **Lifecycle sub-skills not authored.** Per `cnos#395 Non-goals`, `issue/`, `design/`, `plan/`, `review/`, `release/`, `post-release/` for CDR are deferred to cds emergence. Recorded; not blocking.

4. **Project-specific stricter-floor template not provided.** Same as #390 known debt (F2 in `cdd-iteration.md`). Some role overlays cite the project's `.cdr/POLICY.md` as the surface for project-binding rules; no starter template ships with this Sub. Recorded; not blocking.

5. **Wave-coordination primitive not authored.** Engineering γ has `cnos.cdd/skills/cdd/gamma/SKILL.md §10` for multi-cycle wave coordination. Research γ may need an analogue; deferred to Sub 4 or later cycle.

## §CDD Trace

1. **Read issue** — issue #395 read fully (mcp__github__issue_read).
2. **Load Tier 1 + Tier 3 skills** — CDR.md, ROLES.md §4a.2 + §4a + §3, the 5 cnos.cdd role exemplars (alpha 396 lines, beta 196, gamma 705, operator 678, epsilon 74).
3. **γ scaffold authored** — `.cdd/unreleased/395/gamma-scaffold.md` (231 lines); surfaces enumerated, peer-evidence recorded, AC oracle approach declared, Sub 2 loader-collision strategy documented.
4. **Design notes** — `.cdd/unreleased/395/design-notes.md` (139 lines); per-role research-function-vs-software-function differential, persona/protocol/project boundary discipline, loader-stub design.
5. **Build** — 5 role overlay SKILL.md files + 1 loader stub authored at `src/packages/cnos.cdr/skills/cdr/{alpha,beta,gamma,operator,epsilon}/SKILL.md` + `cdr/SKILL.md`. Total 779 insertions.
6. **AC oracle mechanical run** — AC1 PASS (5/5 files, all >2KB); AC2 PASS (11 hits across 5 files vs. ≥5 required); AC3 PASS (6 hits, all disavowal/cross-reference, same-class as #390); AC4 PASS (0 hits in Sub 3 surface; 2 pre-existing hits in CDR.md in disavowal context); AC5 PASS (15 hits vs. ≥5 required); AC6 PASS (hits classified as abstract role-function, schema-field-name, disavowal, cross-reference).
7. **Self-coherence** — this file.
8. **Pre-review gate** — verified rows below.

## §Review-readiness

| row | check | result |
|---|---|---|
| 1 | `origin/cycle/395` rebased onto current `origin/main` | base SHA `e531dba0` matches `origin/main` HEAD at session start; no main advance during cycle (γ+α+β-collapsed single session) |
| 2 | self-coherence carries CDD Trace through step 7 | yes (this file) |
| 3 | tests present, or explicit reason none apply | n/a; docs-only role-overlay-authoring cycle |
| 4 | every AC has evidence | yes (§ACs table per AC) |
| 5 | known debt explicit | yes (§Debt items 1-5) |
| 6 | schema/shape audit when contracts changed | n/a (no schemas modified) |
| 7 | peer enumeration completed | yes (γ-scaffold §Peer enumeration; 10 surfaces ls/grep'd) |
| 8 | harness audit when schema-bearing contract changed | n/a |
| 9 | post-patch re-audit (multi-language) | n/a (Markdown-only diff) |
| 10 | branch CI green on head commit | n/a (no CI run on docs-only cycle); branch HEAD pushed `0490fc5a` |
| 11 | artifact enumeration matches diff | yes (γ-scaffold §Expected diff scope enumerates every file in the diff) |
| 12 | caller-path trace for new modules | n/a (no new code modules); each role overlay file is referenced from the loader at `cdr/SKILL.md` (AC5: 15 hits) |
| 13 | test assertion count from runner output | n/a |
| 14 | α commit author email matches canonical role pattern | `alpha@cdd.cnos` per `cnos.cdd/skills/cdd/operator/SKILL.md` Git-identity-for-role-actors. Path: this cycle is single-session with role-switched commits; α commits are `alpha@cdd.cnos` (verified: `git log --format='%ae' cycle/395` shows γ scaffold by `gamma@cdd.cnos`, design + role-overlay commits by `alpha@cdd.cnos`) |
| 15 | γ-artifact presence at rule-3.11b surface | `.cdd/unreleased/395/gamma-scaffold.md` exists on `origin/cycle/395` at the canonical §5.1 path — verified by `git ls-tree origin/cycle/395 .cdd/unreleased/395/gamma-scaffold.md` |

Implementation SHA: `0490fc5a` (last α-implementation commit before self-coherence). Review-readiness on branch HEAD; β polls cycle/395.

**Ready for β.**
