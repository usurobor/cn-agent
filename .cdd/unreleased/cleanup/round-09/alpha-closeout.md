# α Closeout — Round 09 (terminal docs TRIM: numbered/embedded CDD apparatus in 4 canonical specs)

**Role:** α (implementer). **Branch:** `claude/repo-cleanup-newcomer`. **HEAD at start:** `b84230cc`.
**Mandate:** de-historicize the non-dot reader surface to current-state only; git history is the archive. Dotdirs EXEMPT. TRIM round — the 4 targets are current-state specs; strip embedded CDD design-cycle apparatus, preserve 100% of spec substance. No commit, no PR.
**Contract:** γ Round-08 closeout §2/§5 — the bounded residual class of 4 canonical specs carrying the apparatus in numbered/embedded form.

**Headline:** All 4 trimmed. 18 ins / 240 del across exactly the 4 target docs; zero defer paths touched; zero deletions of files. Every removed block is process apparatus; all rationale/limitation substance preserved (Coherence-Contract motivation and failure-modes rescued into `## Purpose` sections). Numbering renumbered contiguously, all internal `§N` refs resolve, no anchor targets broken. **Convergence grep: the numbered/embedded apparatus class is fully cleared on the in-scope reader surface — no 5th in-scope spec.** Remaining grep hits are all templates, the HARD-DEFER DESIGN-266, or legitimate kata rubric content (classified below).

---

## Per-doc trim

### 1. `docs/reference/protocol/cn/THREAD-EVENT-MODEL.md`
- **Removed:** `## 0. Coherence Contract` (Gap / Named incoherence / Failure modes / `### Mode` MCA / `### α / β / γ target` / Smallest coherent intervention); `## 20. File Changes` (create/edit plan incl. `src/transport/cn_thread_event.ml` etc.); `## 21. Acceptance Criteria` (11 unchecked `[ ]`).
- **Rescued → `## Purpose` (new, unnumbered, replaces §0):** the Gap's what-cnos-has/lacks framing folded into one paragraph; **the 8 Failure modes preserved verbatim** as "defends against these failure modes" (genuine design rationale — same treatment γ verified for MESSAGE-PACKET in R8). The "Named incoherence" (too-many-roles) and "Smallest coherent intervention" content is already covered by `## 1. Core Decision` (canonical stack) and `## 19. Non-goals` — no unique fact lost.
- **`## 22. Known Debt` → `## 20. Limitations`, per-item:** kept 5 genuine current limitations (reflections outside event model v1; subscription UX/policy undefined; feed pagination/index undesigned; chain/nostr/mailbox locators unsupported; event-store/locator retention-GC undesigned), rephrased from "deferred" to current-state "not yet". **Dropped 1** pure future-work/migration-planning item ("Legacy thread-file migration needs its own plan" — migration is owned by the kept body; not a current capability gap).
- **Renumber:** §1–19 unchanged; Known Debt §22 → §20 Limitations. Sequence now …17,18,19,20 contiguous.
- **No substance lost.**

### 2. `docs/reference/packages/EXTENSION-REGISTRY.md`
- **Removed:** `## 0. Coherence Contract` (Gap / `### Mode` MCA / `### α / β / γ target` / Smallest coherent intervention); `## 20. Acceptance Criteria` ("This is done when: 1…6").
- **Rescued → `## Purpose` (new, replaces §0):** the Gap's ecosystem-layer motivation (publish / discover / trust / install-lifecycle / channels / bundles / local-vs-registry) folded into one paragraph. Smallest-coherent-intervention content (registry / trust / lifecycle / bundle layer above Runtime Extensions) is already covered by `## 1. Relationship to Runtime Extensions` — no unique fact lost.
- **Kept:** §18 Alternatives Considered, §19 Migration Strategy, §21 Summary.
- **Renumber:** §21 Summary → §20 Summary. Sequence …17,18,19,20 contiguous.
- **No Known Debt in this doc.** No substance lost.

### 3. `docs/reference/runtime/extensions/RUNTIME-EXTENSIONS.md`
- **Removed:** `## 0. Coherence Contract` (Gap [cites #67] / `### Mode` MCA / `### α / β / γ target` / Smallest coherent intervention); `## 15. Acceptance Criteria` ("This is done when: 1…7"); trailing **`## CLP Summary`** (design-doc self-grade: bare `α: A / β: A / γ: A` + "No material axis weakness remains inside this artifact. The next work lies outside the doc" — pure cycle-completion grading, no spec substance).
- **Rescued → `## Purpose` (new, replaces §0):** the Gap's incoherence ("cognitive substrate is packageable and local, while runtime capabilities are not yet equally extensible"), the #67/network-access driver, and the additive-capability-families intent folded into one paragraph. `## 1. Core Decision` already carries the extension-vs-core rule; no unique fact lost. (Header `**Addresses:** #67` and `**Iteration history:**` left intact — version-lineage metadata, not in the apparatus strip list, and γ scoped this doc to the two blocks only.)
- **Kept:** §13 Alternatives, §14 Migration Strategy, §16 Summary, §17 Marketplace-Readiness.
- **Renumber:** §16 Summary → §15; §17 Marketplace-Readiness → §16 (incl. its subsections 17.1/17.2/17.3 → 16.1/16.2/16.3). Sequence …13,14,15,16 contiguous.
- **No Known Debt section.** No substance lost.

### 4. `docs/architecture/cognitive-substrate/COGNITIVE-SUBSTRATE.md`
- **Removed:** trailing `## Coherence Contract for This Document` (`**Gap:**` / `**Mode:** MCI` / `**Scope:**` / `**Expected effect:**` α-β-γ / `**Failure if skipped:**` / `**CLP:**` TERMS/POINTER/EXIT).
- **No reframe needed:** this doc already opens with a current-state `## 0. Purpose`; the Gap's "no single normative document for classifying substrate" motivation is already stated there. The "Known misclassifications" (COHERENCE-as-mindset) is already captured in `## 10.5`; the alignment/defer obligation is already captured in `## 12` ("Any document that explains doctrine/mindsets/skills MUST defer to this spec"). The block's "Remaining MCA: align CAR/AGENT-RUNTIME/WHITEPAPER/skills-README" is future-work planning — no current-state fact lost.
- **Kept:** §12 Compatibility and Migration (current-state normativity statement) and all of §1–11.
- **No renumbering** (removed block was unnumbered trailing; §0–12 untouched). Internal refs §6.1 and §7.3 still resolve.
- **No substance lost.**

---

## Facts rescued into body sections (what → where)
- THREAD-EVENT-MODEL: the 8 named **Failure modes** (transport/semantic collapse, authority duplication, routing leakage, reply-propagation ambiguity, thread-discovery ambiguity, identity/locator confusion, locator fragility, projection-authority drift) → `## Purpose`. Gap has/lacks framing → `## Purpose`.
- EXTENSION-REGISTRY: ecosystem-layer concern list (publish/discover/trust/lifecycle/channels/bundles/local-vs-registry) → `## Purpose`.
- RUNTIME-EXTENSIONS: the packageable-cognition-vs-non-extensible-capabilities incoherence + #67 network-access driver → `## Purpose`.
- COGNITIVE-SUBSTRATE: none needed — all Coherence-Contract substance already present in §0 Purpose / §10.5 / §12.

## Limitations preserved
- THREAD-EVENT-MODEL `## 20. Limitations`: 5 kept (reflections-outside-model / subscription-UX / feed-pagination / chain-nostr-mailbox-locators / retention-GC), 1 dropped (legacy-migration planning).
- The other 3 docs had no `Known Debt` section.

---

## MANDATORY convergence grep — result

Mandated pattern (numbered + unnumbered heading forms + header stamps) across all `docs/**`:

```
git grep -nE '^#{1,4} +([0-9]+\. +)?(Coherence Contract|Acceptance Criteria|CDD Trace|Impact Graph|File Changes)|^(Issue|Mode|Active Skills|Engineering Level|Cycle|Branch): ' -- 'docs/**'
```

Hits (all classified — none is an in-scope spec needing a trim):

| Hit | Classification |
|-----|----------------|
| `docs/development/cdd/PLAN-TEMPLATE.md:3,5,6` (`Issue: #NN` / `Mode: MCA / MCI` / `Active Skills: ...`) | **TEMPLATE — legitimate KEEP.** Placeholder scaffolding (`#NN`, both-options, ellipsis) showing the CDD plan format; current-state tooling. Upstream template apparatus is the code-pass's concern (γ R8 §5 CDD-templates handoff item). |
| `docs/development/cdd/SELF-COHERENCE-TEMPLATE.md:3,5,6,20` (+`## Acceptance Criteria Check`) | **TEMPLATE — legitimate KEEP.** Explicitly cleared by γ R8. Same class as PLAN-TEMPLATE. |
| `docs/reference/packages/DESIGN-266-dist-out-of-git.md:47,203,239,259` (`## Impact Graph` / `## File Changes` / `## Acceptance Criteria` / `## CDD Trace`) | **HARD-DEFER.** One of the 6 code-coupled do-not-touch docs; deletion + coupling-repair owned by the code pass. |

Supplementary sweep (Known Debt / CLP Summary / α-β-γ target / Smallest coherent / `**Mode:** MC*` / Failure if skipped / "Coherence Contract for This"), deferred-doc set excluded:

| Hit | Classification |
|-----|----------------|
| `docs/development/cdd/SELF-COHERENCE-TEMPLATE.md:33 ## Known Debt` | **TEMPLATE — legitimate KEEP** (same template as above). |
| `docs/development/kata/{A1,B1,C1}-*.md ## CLP Summary — v1.0.1 (Evaluator-side, derived)` and `KATA-EVALUATION.md:474 ## 17. CLP Summary` | **Legitimate KEEP — kata rubric content, NOT design-cycle apparatus.** These `## CLP Summary` sections are the kata's intrinsic TERMS/POINTER/EXIT rubric (what the exercise's invariant is, the fix pointer, the exit criteria) — the deliverable substance of a kata/evaluation artifact. Distinct from the RUNTIME-EXTENSIONS `## CLP Summary` I removed, which was a design-doc self-grade (`α: A / β: A / γ: A` + "next work lies outside the doc"). The α/β/γ in KATA-EVALUATION §17 is the scoring-vector *methodology* ("Score α / β / γ / efficiency as a vector"), not a cycle grade. No trim needed. |

**Verdict: the numbered/embedded CDD-apparatus class is FULLY CLEARED on the in-scope reader surface. No 5th in-scope spec. No Round-10 docs matter surfaced.** The only surviving apparatus-signature hits are (a) CDD doc templates — current-state tooling whose upstream apparatus is the code-pass template-fix item, (b) DESIGN-266 — HARD-DEFER code-coupled, and (c) kata rubric content — legitimate current-state substance.

---

## Coherence / falsification
- **Renumbering contiguous** in all 4: THREAD-EVENT-MODEL …17,18,19,20(Limitations); EXTENSION-REGISTRY …17,18,19,20(Summary); RUNTIME-EXTENSIONS …13,14,15(Summary),16(Marketplace-Readiness, subsections 16.1–16.3); COGNITIVE-SUBSTRATE §0–12 untouched.
- **Internal `§N` cross-refs:** `git grep '§[0-9]|section [0-9]'` in each doc → only COGNITIVE-SUBSTRATE §6.1 (×2) and §7.3, both still-existing targets. Zero dangling refs. Other 3 docs carry no internal `§` refs.
- **Anchor/inbound links:** no repo file links to any removed heading of the 4 docs as a `#`-fragment (`git grep` for `<file>.md#` and for `#coherence-contract|#acceptance-criteria|#file-changes|#clp-summary|#N-...` → clean). The §20/21/22 hits in `src/packages/cnos.cdd/skills/**` are that SKILL.md's own self-numbering, unrelated to the trimmed docs.
- **Re-read** each trimmed doc: coherent current-state spec, opens with Purpose/spec body, no lost architecture/contract/schema/rationale.

## Defer paths — confirmed untouched
`git status --short` = exactly 4 `M` under `docs/` (the 4 targets) + this receipt. No `.cdd/**` (bar this receipt), `.cn-sigma/`, `src/`, `schemas/`, `scripts/`, `tests/`, `.github/`, `install.sh`, `cn.json`, no root `RELEASE/CHANGELOG/ROLES/OPERATOR/SUSTAINABILITY`, no `docs/development/board/`, none of the 6 code-coupled docs, no CDD skill templates. Diffstat: 18 ins / 240 del, 4 files.

## Convergence read
The docs pass converges after this round. Rounds 1–7 cleared the whole-doc historical-record class; R8 cleared the header-stamp CDD apparatus on the 6 primaries + 5 siblings; R9 clears the terminal numbered/embedded-apparatus class on the last 4 canonical specs. The non-dot reader surface is now current-state-only modulo the explicitly HARD-DEFER code-coupled set (DESIGN-266 + 5 others) and the CDD doc/skill templates whose upstream apparatus-prescription is the code pass's to reconcile. No further docs-only apparatus class remains.
