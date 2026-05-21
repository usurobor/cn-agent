<!-- sections: [intake, role-collapse, contract-summary, dispatch-pinning, dispatch-record] -->
<!-- completed: [intake, role-collapse, contract-summary, dispatch-pinning, dispatch-record] -->

# γ scaffold — cycle/402

## Intake

- **Issue:** [cnos#402](https://github.com/usurobor/cnos/issues/402) — Phase 7 (#366): CDD.md rewrite — compress to CCNF spine; software-specific content named for cds extraction
- **Parent:** [cnos#366](https://github.com/usurobor/cnos/issues/366) — coherence-cell executable-protocol roadmap (terminal phase)
- **Mode:** design-and-build
- **Priority:** P1
- **Labels:** P1, cdd, design
- **Branch:** `cycle/402` from `origin/main` (base `8c3d573b`)

## Role collapse

This cycle runs under δ-as-agent with γ+α+β collapsed on δ (per dispatch). The collapse is acknowledged by:

- γ writes scaffold (this file), dispatch pin, and close-out triage.
- α writes design notes, performs the CDD.md rewrite, updates cross-references, and writes self-coherence.
- β (collapsed-on-δ) runs the AC1–AC7 sweep with the extra rigor δ requested on AC7 (hard rule preconditions explicit and verified).
- δ executes the merge to main and the post-merge closure on cnos#366.

The collapse does not relax authoring discipline. Each role's artifacts land at their canonical filenames in `.cdd/unreleased/402/`. The β review at AC7 is given the extra rigor δ pinned in the dispatch.

## Contract summary (from issue body)

**Target deliverable:** Rewrite `src/packages/cnos.cdd/skills/cdd/CDD.md` (currently 1344 lines) into a compact CCNF-spine document. The operative target structure is pinned in the issue body and reproduced below as the implementation contract.

**ACs (7):**
1. CDD.md `wc -l` ≤ 300 (from 1344, ~22% target).
2. CCNF five-step recursion equation is the structural spine — appears verbatim from COHERENCE-CELL-NORMAL-FORM.md in first ~50 lines.
3. Domain packages named (CDS, CDR, future c-d-X).
4. Pointers section names all canonical surfaces (doctrine, schemas, roles, runtime substrate, realization peers).
5. Software-specific content disposition explicit — every pre-cycle §§3–10 section either moved to cnos.cds or named "pending cds extraction (tracked at cnos#NN)". No silent drops.
6. Cross-references from other skill files resolve in the new CDD.md (`rg "CDD.md §" src/packages/cnos.cdd/skills/cdd/` — every hit resolves or is updated to content-citation).
7. Hard rule satisfied: V is executable (#392); domain evidence has homes (schemas/cds/ + schemas/cdr/ exist per #388; cdr v0.1 per #376; cds extraction tracked).

**Hard-rule precondition verification (β-collapsed will recheck):**
- V executable: `src/go/cn cdd verify --help` includes `--receipt <path>` (the V validator dispatch) — confirmed at cycle intake. Verified buildable from `src/go/cmd/cn/main.go` against `src/packages/cnos.cdd/commands/cdd-verify/`.
- Domain evidence homes: `schemas/cdd/` (generic), `schemas/cds/` (software, fixtures only — package itself pending bootstrap), `schemas/cdr/` (research, receipt.cue + fixtures shipped) all exist on `origin/main`.

## Dispatch pinning

The implementation contract is pinned by the issue body and δ's dispatch; α MUST NOT improvise on the operative target structure. The key constraints:

- **CCNF recursion equation appears verbatim** from COHERENCE-CELL-NORMAL-FORM.md §Kernel (lines 58–64 of that file) — no paraphrase.
- **Line count target ≤ 300** is binding; surface as refusal condition if the cross-reference work makes it unreachable.
- **No CCNF-X formalization** — Phase 7 ships the conservative compact rewrite. Mode enum, sizing predicate, master+sub graph, dispatch-prompt schema, findings state machine are CCNF-X follow-on, NOT this cycle.
- **No cnos.cds package creation** — cnos.cds does not exist at cycle time (verified: `ls src/packages/` returns `cnos.cdd cnos.cdd.kata cnos.cdr cnos.core cnos.eng cnos.kata`). Therefore §§3–10 content stays in CDD.md with "pending cds extraction" notes pointing at a tracker issue that δ files as part of this cycle (or under a "named for future cds bootstrap" pointer if the tracker issue is unavailable).
- **All cross-references resolve.** `rg "CDD.md §" src/packages/cnos.cdd/skills/cdd/` returns hits across post-release, release-effector, release, gamma, harness, operator, alpha, review, beta, activation, plus the CI workflow templates. Each hit is either re-pointed to a new section, re-cited by content ("the kernel" / "the cycle algorithm"), or kept as a software-realization pointer with the "pending cds extraction" disposition.

## Dispatch record

γ → α: design-and-build per pinned contract. α produces the rewritten CDD.md, updates cross-references, writes self-coherence with §"Hard rule satisfied" subsection and §"Cross-reference map" table.

γ → β-collapsed: AC1–AC7 sweep with particular rigor on AC7 — verify both hard-rule preconditions (V executable + domain evidence homes) hold; verify no domain-specific evidence requirements appear in the generic CDD.md surface.

γ → δ: post-α/β-collapsed, execute merge to main, comment on cnos#366 declaring the roadmap COMPLETE with all phase merge SHAs cited, close cnos#366 with state_reason=completed, update INDEX.md row for #402, attempt branch cleanup (expect 403).
