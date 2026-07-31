# α Close-out — Code Pass, Round 01 (execution)

**Role:** α (implementer). **Branch:** `claude/repo-cleanup-code-pass` (stacked on docs PR #679). **Base HEAD:** `a66f4cf0`.
**Mandate:** the working tree reflects only current state; git history + dotdirs (`.cdd/`, `.cn-sigma/`, `.github/`) are the archive.
**Method:** executed the operator-approved scope in phase order per β's verified coupling tables; self-verified against the falsification gate; no commit, no PR (σ commits).

**Headline:** All four falsification gates pass. 7 files deleted, 3 relocated, 48 edited. Zero dangling refs to any of the 6 deleted basenames; zero old-root links to the 3 relocated files; skill-frontmatter validator green (99/99, self-test green); CHANGELOG / RELEASE / board / release-gate / ledger.go untouched. The CTB cascade surface was ~5× wider than β's Item 4 estimate (see Phase 3 note); no dependent doc was left missing content — every inbound ref was a pointer/citation repaired to stand alone.

---

## Phase 1 — doctrine + root-cause

### Item 5 — Frozen-history doctrine rewrite (ratified)
`docs/reference/governance/DOCUMENTATION-SYSTEM.md` §5 "Frozen history", paragraph 2 replaced with β's drafted text. Paragraph 1 and the "Supersession" (`Supersedes:` / `Superseded by:`) subsection left unchanged. Applied text:

> A completed, dated, or superseded record — a design decision, a cycle log, a finished plan — is **not kept on the reader surface**. Git history is its archive: to read it, check out the release or commit under which it was written. The working tree on `main` reflects only current state. Records that must persist as intentional history live in the dotdirs (`.cdd/`, `.cn-sigma/`, `.github/`), which are exempt from the current-state rule. A document that remains on the reader surface is, by that fact, current — its paths and contents are kept correct, not frozen.

### Item 6 — Promotion contract (ratified) — BODY EDITS ONLY, frontmatter untouched
`src/packages/cnos.cdd/skills/cdd/design/SKILL.md`:
- **§3.1 "Output Format"** annotated as *cycle-artifact form (working cell) — not the promoted reader-surface form*, with a forward pointer to §3.2.
- **New §3.2 "Promotion"** inserted after §3.1 (existing §3.2 "Keep the CDD Trace…" → §3.3; §3.3 "Pre-Submission Checklist" → §3.4). The Promotion rule: on promotion to a canonical `docs/reference/**` | `docs/architecture/**` doc, **retain** governed header (`Version:`/`Status:`/`Doc-Class:`) + `## Purpose` + spec body + `## Alternatives`/`## Migration`/`## Non-goals`, fold real limits into `## Limitations`; **drop** `Mode`/`Active Skills`/`Engineering Level`/`## Impact Graph`/`## File Changes`/`## Acceptance Criteria`/`## CDD Trace` (they live in the cycle record; git history preserves them). Includes one ❌/✅ pair.
- **Renumber safety:** the only live (non-`.cdd`) cross-ref to this file's §-numbering is `CDS.md:1710 → §3.1` (Output Format), unaffected. The `§3.2`/`§3.3`/`§3.10` refs in `CDS.md`/`RECEIPT-VALIDATION.md` describe rules ("one source of truth", "package cohesion") that map to the separate `cnos.core` design skill's grammar, not this file — verified they do not resolve against this file's headings, so the renumber introduces no new mismatch.
- **Stale example** refreshed at old line 257: `"Issue: #113. Plan: PLAN-package-system.md. Prior art: EXTENSION-REGISTRY.md."` → `"Issue: #404. Plan: design-notes.md. Prior art: RUNTIME-EXTENSIONS.md."` (both new referents are current/real).

**Template scaffolding notes** (one HTML-comment line each, at top):
- `.github/PULL_REQUEST_TEMPLATE.md` and `docs/development/cdd/SELF-COHERENCE-TEMPLATE.md` now carry: *"Cycle scaffolding: this is the working cell. When a design artifact is promoted to a canonical `docs/reference/**` or `docs/architecture/**` doc, strip the cycle apparatus per the Promotion rule in `src/packages/cnos.cdd/skills/cdd/design/SKILL.md` §3.2."* `PLAN-TEMPLATE.md` left alone (clean, per β).

**Frontmatter validator (post-edit):** `scripts/ci/validate-skill-frontmatter.sh` — self-test **PASS** (exit 0); full run **PASS** (99 SKILL.md, no findings). (`cue v0.13.2` installed locally to the scratchpad to match build.yml's pin.)

---

## Phase 2 — delete 4 genuine records + repair citations

**Deleted:** `docs/reference/packages/DESIGN-266-dist-out-of-git.md`; `docs/papers/ENGINEERING-LEVEL-ASSESSMENT.md`; `docs/development/cdd/DISPATCH-FAILURE-EVIDENCE.md`; `docs/development/design/cn-repo-install-MOCKS.md` (self-refs at DESIGN-266:228/241/247 and MOCKS:260 went with the files).

**Citation repairs (all line edits; kept lesson/provenance, dropped filename):**
- `.gitignore:37` — "Never committed — regenerated build output, not source."
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md:176` — kept the #266 F3 / F3-bis intra-doc-repetition lesson; "a design doc carried one count across 4 sentences…" (filename dropped).
- `docs/development/ENGINEERING-LEVELS.md` — dropped the assessment bullet from `Related:`; removed §10 "Relationship to the assessment essay" (entire section was about the deleted essay; no external §-refs) and renumbered §11 Summary → §10.
- `docs/papers/README.md` — removed the now-empty "## Engineering and release" section (its only bullet was the deleted essay).
- `src/packages/cnos.eng/skills/eng/README.md` — dropped the assessment bullet from `Related:`; removed "## Relationship to the level assessment essay" (essay-specific; normativity already stated in Purpose).
- `src/packages/cnos.cdd/skills/cdd/release-effector/SKILL.md:79` — kept "manual tagging is not allowed"; dropped the `DISPATCH-FAILURE-EVIDENCE.md, cycle #84 failure 3` parenthetical.
- CI **comments** (verified comment-only; no workflow reads these files, no golden fixture asserts the strings): `.github/workflows/build.yml:72` (kept `cnos#612 Mock F`), `.github/workflows/install-wake-golden.yml:692` (kept `cnos#609 / Mock C`), `src/packages/cnos.core/commands/install-wake/cn-install-wake:40/1070/1082` (kept Mock C / E2/E4 / E provenance; dropped file path — incl. the echoed generated-comment line, confirmed absent from all golden fixtures).
- `docs/guides/INSTALL-CDS.md:253` — removed the "Related" bullet citing the MOCKS file.

**Directory lifecycle:** `docs/development/design/README.md` was a **stub index** ("No decision surveys are active here now… Past surveys are archived in git history"). After deleting `cn-repo-install-MOCKS.md` the dir held only that stub with no external referrers → **removed the dir + README** per the decision rule.

**Verification:** `git grep` for the 4 basenames (excl `.cdd/`, `CHANGELOG.md`) → **zero**.

---

## Phase 3 — CTB pair deletion (cascade) + root relocation

### CTB deletion (operator override of β's keep — full cascade)
**Deleted:** `docs/reference/ctb/CTB-v4.0.0-VISION.md`, `docs/reference/ctb/LANGUAGE-SPEC-v0.2-draft.md`.

**Cascade repairs — every inbound ref rewritten so the doc stands alone (β's list + prose "the Vision" governance refs β's filename-grep missed):**
- `docs/reference/ctb/README.md` — bundle README rewritten: "four documents" → "two documents"; Document Map now lists only `LANGUAGE-SPEC.md` (named as **the** canonical spec, satisfying the §7 one-canonical-spec CI invariant) + `SEMANTICS-NOTES.md`; Authority section drops the v0.2-promotion clause and the Vision-governs clause.
- `docs/reference/ctb/LANGUAGE-SPEC.md` — lines 15 (Vision-surfaces-seams prose), 19 (strategy-document "It is not" bullet), 23 (Vision-disagree clause), 315 (`(Vision §8.5.1 E1)` provenance tag), 508 (Vision-governs Authority bullet), 513 (Vision-disagree Authority clause) all removed/rewritten to stand alone.
- `docs/reference/ctb/SEMANTICS-NOTES.md` — lines 7 (Companion-to Vision), 13 ("The Vision sets direction"), 610 (Vision-governs clause) rewritten/removed.
- `docs/reference/runtime/ORCHESTRATORS.md` — §325 "See" link repointed to `../ctb/LANGUAGE-SPEC.md`; §339 dropped four `(CTB-v4.0.0-VISION §3.1.x)` anchors (kept the property bullets); §371 "This matches CTB-v4.0.0-VISION.md exactly:" → "This is the CTB model:" (kept bullets, dropped §-anchors); §390 "The CTB vision doc says it explicitly:" → "on the premise that:" (kept the three premise claims inline, dropped the deleted-doc attribution + §1 anchors).
- `src/packages/cnos.core/skills/agent/emoji-language/SKILL.md:70` — "See:" link repointed from the Vision to `docs/reference/ctb/README.md`.
- `docs/papers/ACTIVATION-NOT-DEPLOYMENT.md:178` — illustrative JSON receipt's `artifacts[]` v0.2-draft path repointed to the surviving `LANGUAGE-SPEC.md`.
- `schemas/README.md:184` — removed the v0.2-draft bullet (kept the v0.1 `LANGUAGE-SPEC.md` bullet).
- **Extra ref beyond β's list, caught + fixed:** `src/packages/cnos.eng/skills/eng/ship/SKILL.md:159` cited "CTB vision §8.5.2" as a rebase-drop incident example → generalized to drop the deleted-doc section reference (kept the surviving `COHERENCE-FOR-AGENTS.md` instance).

**Confirmed:** `git grep` for both CTB basenames (excl `.cdd/`, `CHANGELOG.md`) → **zero**; no surviving "the Vision"/"CTB vision" governance prose remains in the CTB bundle, ORCHESTRATORS, or the two skills.

**FLAG — substance removed by the operator's delete (informational, not a broken dependency):**
- **No dependent doc was left missing content it needs.** Every inbound reference to either deleted file was a pointer, cross-reference, citation, or (for ORCHESTRATORS §390) a quotation whose claims I preserved inline. All were repaired to stand alone.
- **However**, the delete does remove genuine current-state substance that now lives **only in git history**: (a) `LANGUAGE-SPEC-v0.2-draft.md` was the **sole** definition of the v0.2 agent-module migration target — agent-as-primitive, triadic carrier, and the composition-operator algebra (`>>`, `>>=`, `|||`, `case`, `fix`, `wait`, `try`). No surviving doc carries that grammar. (b) `CTB-v4.0.0-VISION.md` was the only consolidated CTB strategy/motivation/roadmap. The "why CTB exists" motivation survives in `docs/papers/EXECUTABLE-SKILLS.md`; the §1 premise survives inline in ORCHESTRATORS §7.5; the roadmap does not survive on the reader surface. This is the operator's explicit decision (de-historicize the forward-looking roadmap); flagged so σ is aware the v0.2 grammar is retrievable only via `git show`/checkout.

### Root relocation (operator decision)
Homes chosen by governing question + existing tree, via `git mv`:
- **`ROLES.md` → `docs/concepts/ROLES.md`** — the generic role-scope-ladder doctrine (SSOT for the role structure) is the mental model `concepts/` answers. Added to `docs/concepts/README.md` "Doctrine and origin" nav.
- **`OPERATOR.md` → `docs/guides/OPERATOR.md`** — the Operator Manual is a day-2 task how-to. Added to `docs/guides/README.md` "## Operator" nav.
- **`SUSTAINABILITY.md` → `docs/SUSTAINABILITY.md`** — project-meta funding/sponsorship stance (`Doc-Class: reference`) that fits no single reader-intent bucket cleanly; kept at `docs/` top-level as a THESIS-peer project doc, discovered via the root README "Support" link. (Noted tension with DOCUMENTATION-SYSTEM §1 "root = THESIS + README"; §6 covers the no-clean-bucket case. σ may prefer forcing it into an intent dir — flagged as a judgment call.)

**Link-repair surface (materially larger than β's Item 4 — β listed ~15 sites; the true ROLES.md surface is ~19 files / ~80 link sites, dominated by the CDR/CDS/CDD skill trees + schema receipts):**
- **Inbound relative links** (`../…/ROLES.md`, all engineered to hit repo root) uniformly repointed `/ROLES.md` → `/docs/concepts/ROLES.md` across 19 files (schemas/{cdd,cdr,cds}/README.md; cnos.cdd activation/epsilon/gamma/post-release; cnos.cdr README + docs + CDR.md + alpha/beta/epsilon/gamma/operator; cnos.cds README + docs + CDS.md + SKILL.md), preserving `../` depth and `#anchors`. Verified: no non-root `/ROLES.md` existed, so the transform was safe and uniform.
- **Papers `related:` frontmatter** (`- ROLES.md` → `- docs/concepts/ROLES.md`): BOX-AND-THE-RUNNER, CCNF-AND-TYPED-TRUST, CELL-OF-CELLS, DECREASING-INCOHERENCE.
- **Moved ROLES.md's own outbound links** (`](src/packages/…)` → `](../../src/packages/…)`, 6 links).
- **"at the repo root" location prose** that the move falsifies, corrected in `RECEIPT-VALIDATION.md` (×3), `CDR.md:9`, `CDS.md:9`.
- **`SUSTAINABILITY.md`** `**Canonical-Path:**` stamp → `docs/SUSTAINABILITY.md`; root `README.md:195` link → `docs/SUSTAINABILITY.md`.
- **OPERATOR.md outbound links** (from `docs/guides/`): `](docs/…)` → `](../…)`, `](README.md…)` → `](../../README.md…)` (9 links).
- **Left as-is (correct):** bare-name prose/comment citations `` `ROLES.md §X` `` without a path (basename unchanged → still accurate; e.g. the `.cue`/`.yaml` receipt comments, CDD.md:95/107, many skill-body prose citations). These are name citations, not path links, so they neither break nor point at an old path.
- **Not touched (correct):** every `spec/OPERATOR.md` reference (per-hub identity spec — `activate.go`, `activate_test.go`, golden fixtures, R5-activate kata); `git status` on `src/go/internal/activate/` is clean.

**On-disk link resolution spot-checked:** inbound (schemas/cdd, cds/CDS, epsilon → ROLES), ROLES outbound (→ epsilon), OPERATOR outbound (README, SETUP-INSTALLER, AUTOMATION, CLI, PACKAGE-SYSTEM), emoji-language → ctb/README, and `docs/SUSTAINABILITY.md` — all resolve.

---

## Self-verification (falsification gate)

1. **6 deleted basenames**, `git grep` excl `.cdd/` + `CHANGELOG.md` → **ZERO dangling refs.** (PASS)
2. **Old-root links to ROLES/OPERATOR/SUSTAINABILITY** → **zero.** The only `](ROLES.md)` / `](OPERATOR.md)` hits are the correct **same-directory** nav links I added inside `docs/concepts/README.md` and `docs/guides/README.md` (the files now live in those dirs). No frontmatter `- ROLES.md` root entry remains. `spec/OPERATOR.md` refs confirmed untouched. (PASS)
3. **`scripts/ci/validate-skill-frontmatter.sh`** — self-test exit 0; full run exit 0 (99 SKILL.md, no findings). (PASS)
4. **Hard-constraint paths untouched** — `git status` clean for `CHANGELOG.md`, `RELEASE.md`, `docs/development/board/`, `scripts/release.sh`, `scripts/validate-release-gate.sh`, `ledger.go`. (PASS)

## Hard constraints — confirmed intact
`CHANGELOG.md`, `RELEASE.md`, `docs/development/board/`, `scripts/release*.sh`, `validate-release-gate.sh`, `ledger.go`, the release gate: **untouched.** `design/SKILL.md`: **body only**, YAML frontmatter untouched (validator confirms). `.cdd/**`: untouched except this receipt.

## Change-set totals
- **Deleted (7):** 4 records + `docs/development/design/README.md` + 2 CTB docs.
- **Relocated (3):** ROLES → `docs/concepts/`, OPERATOR → `docs/guides/`, SUSTAINABILITY → `docs/`.
- **Edited (48):** doctrine (1), design skill + 2 templates (3), Phase-2 citation repairs (8), CTB cascade (9), root-relocation link repairs (26 incl. the moved files' own edits, the 2 destination nav READMEs, and the root README).

**No commit, no PR** — left for σ.
