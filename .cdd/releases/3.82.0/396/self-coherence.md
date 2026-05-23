# α self-coherence — cycle/396

**Issue:** [cnos#396](https://github.com/usurobor/cnos/issues/396).
**Surface authored:** `src/packages/cnos.cdr/docs/empirical-anchor-cph.md` (217+ lines, single new file).

## Coherence checks (α-side, pre-β-review)

### Surface ↔ design

- All 13 issue-body-enumerated artifact classes appear in §1 (rows 1–13). ✔
- The two survey-discovered extensions appear in §1 (rows 14–15) and are flagged with `**(extension)**` so a reader can distinguish them from issue-body classes. ✔
- The verdict-vocabulary cross-walk in §2 enumerates every cph verdict observed during the survey. ✔
- §3 transposes the mapping by CDR.md field; all six fields plus the persona/protocol/project layer are exercised. ✔

### AC dry-run (α-side, before β-collapsed self-review)

- **AC1 (file exists, non-empty):** `test -s` PASS; file is 217+ lines.
- **AC2 (12+ classes covered):** all 12 issue-body oracle substrings present (`rg -c` returns ≥1 hit per pattern after the row-label adjustment for `wave-manifest`, `wave-receipt`, `diagnostic-oracle`).
- **AC3 (zero "no surface" rows):** `rg "TBD|no surface|not mapped"` returns 2 hits, both in clearly-prose-context per AC3's allowance: §4 (titled "Non-surface concerns") and §6.1 (a meta-statement "No 'no surface' rows" announcing zero such rows in §1). The §1 mapping table itself has zero TBD/no-surface/not-mapped rows.
- **AC4 (cnos#376 AC3 satisfied + cnos#376 cite):** the doc preamble (§Scope) names "Closes cnos#376 AC3"; §6 title is "Closing — cnos#376 AC3 satisfied"; §6 ¶"**cnos#376 AC3 is therefore satisfied**" is the explicit close.
- **AC5 (citations, no extended quotation):** 23 instances of the cph-SHA pin (`@2bf73ad6` or full SHA). Citations are path + SHA pairs. The only verbatim cph content quoted is short identifying snippets (e.g. "GO with bounded scope (path (a) inferred bilateral)", "REVISE for friend pre-pilot", "Construct survives R4 contact subject to anchor caveats", "close Sub C as NO-GO with rationale per protocol") — each ≤ 1 line, used as a verdict-vocabulary anchor, not a prose embedding. No multi-line cph paragraphs reproduced.

### Citation-discipline notes

- The doc cites cph at commit SHA `2bf73ad64b7aa338cc6c8bc48067fbbba818cffa` (resolved at filing time). Short SHA `2bf73ad6` is used throughout for readability; full SHA documented in §0.1 as the authoritative pin.
- The cph#32 reference is by issue-link, not by content embedding. The doc names the failure pattern in mechanism-language ("numbers cited without `/opt/gait-data/` mount") that is closer to the issue body's terms than a paraphrase but is sufficiently shortened that it is mechanism-naming rather than quotation.
- Row 11's full SHA-256 hash for the OpenCap archive is reproduced as evidence of the hash-pinning practice; this is data (a 64-char hex string), not cph-local prose.

### Cross-reference integrity

- `CDR.md` link (`../skills/cdr/CDR.md`) resolves: file exists at `src/packages/cnos.cdr/skills/cdr/CDR.md` (Sub 1, merged at `c0eb9656`).
- `schemas/cdr/receipt.cue` link (`../../../../schemas/cdr/receipt.cue`) resolves: file exists.
- `ROLES.md` link (`../../../../ROLES.md`) resolves: file exists at repo root.
- `docs/gamma/essays/CCNF-AND-TYPED-TRUST.md` link (`../../../../docs/gamma/essays/CCNF-AND-TYPED-TRUST.md`) resolves.
- Role-overlay file paths (`cnos.cdr/skills/cdr/{alpha,beta,gamma,delta,epsilon}/SKILL.md`) are referenced by name only; per design-notes §"Frozen non-goals", the doc does not depend on their merge state. They will resolve when Sub 3 (cnos#395) merges.
- cph link `https://github.com/usurobor/cph` and cph#32 link `https://github.com/usurobor/cph/issues/32` resolve (cph is public).

### Frozen non-goals respected

- No CDR doctrine authored. The doc cites CDR.md sections; it does not extend them. ✔
- No role-overlay files modified or authored. Sub 3's deliverables are referenced by name. ✔
- No cph prose embedded. Citations carry `<path>@<sha>`; short identifying snippets only. ✔
- No CDR-shape changes proposed. The §1 mapping fits within CDR.md v0.1 surfaces; the §6.1 known-gaps section names what would warrant a Sub 1 reopen (no such finding from this survey). ✔

## CDD Trace

| Step | Artifact | Decision |
|------|----------|----------|
| 0 Observe | mapping table draft | 15 rows, all surfaces mapped |
| 1 Select | AC dry-run results | all 5 ACs pass on α-side |
| 4 Gap | none surfaced | mapping is structurally complete |
| 5 Mode | docs-only | hand off to β-collapsed self-review |
| 6 Artifacts | self-coherence.md (this file) | next: beta-review.md |
