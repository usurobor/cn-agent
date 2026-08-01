# self-coherence.md — cnos#686 (α+β+γ collapsed onto the cell)

## Matter

**Issue:** [cnos#686](https://github.com/usurobor/cnos/issues/686) — "activation protocol: ideal behavior for every activation case (home | attached | unattached-attachable | not-attachable | repo-less) — quiet, evidence-derived, non-confabulating."

**Class:** skill/docs. Per Sigma persona commitment 5 (doctrine-sanctioned for skill/docs matter), α+β+γ collapse onto this implementation working cell; δ accepts the receipt and handles transport. No git commit/push performed — files written only.

**Mode:** design-and-build. The five-case taxonomy, five invariants, and three named checks are the converged design in the issue body; this cell writes the exact skill prose + binding checks that encode them.

**The gap.** `agent/activate/SKILL.md` (+ `agent/attach/SKILL.md`) specified a single activation shape ("told to activate as `<hub-url>`; load Kernel → CA → Persona → Operator → state → confirm") with no behavior differentiation by context, no quiet-output discipline, and no no-confabulation rule. Observed failure (issue): a body woke on a bare box with **no repo** and recited a nonexistent "pause," a stale HEAD (`19cf9470`) as current state, "release 3.82.0 merged," a dualist core-drive ("reduce incoherence between model and reality"), and an operator-coupled "Intelligent Assistants Team" identity — none derivable without a repo, all recited from memory as if current. First-contact activation, especially by a human who did not write Sigma, is the trust-setting moment; confabulation there asserts unverifiable claims as fact.

**What this cell shipped.** `activate/SKILL.md` now resolves **exactly one** of five context-cases from observable facts before loading (§2.0), states the five cross-cutting invariants as **binding rules** (§3.9–§3.13), and specifies the three named checks with fixtures (§5.8–§5.13) + failure catalogue (F10–F13). `attach/SKILL.md` gains the **consequence-disclosure step** (§2.3 step 2) — the Case-C attach discloses the exact repo-side changes before any mutation and attaches only on confirmation, with `SILENT_MUTATION` as the named violation (§3.9, §4.6, A9). Channel transport (orphan refs, cursor ownership) is cross-referenced to cnos#684, not re-specified. `PERSONA.md` was **not** edited (operator-gated) — this cell makes activation *read* identity correctly.

## Skills / integration approach

Read both target skills in full before editing; matched the house prose style (terse imperative, `## N.` numbered Define/Unfold/Rules/Verify skeleton with unnumbered Core Principle + Algorithm framing zone, ❌/✅ example pairs, `§`-cross-referenced sections, named failure modes). Integration constraints honored as generation constraints, not post-hoc checks:

- **No renumber of existing `activate` sections.** The `cn activate` renderer parses §4.1's machine-readable load-order block, and §2.1/§2.5/§4.1 are cross-referenced from this skill, from `attach/SKILL.md`, and from the Go renderer. The case procedure was therefore inserted as **§2.0** (runs first, before the §2.1 load order) rather than renumbering §2.1→§2.2 etc. The renderer contract (§4) is byte-untouched.
- **One-source-of-truth (KERNEL §2.2).** The five-case procedure + per-case behavior is stated once (§2.0); the invariants are stated once as rules (§3.9–§3.13); Verify and the failure catalogue **reference** those sections rather than restating them. §2.0's "Cross-cutting invariants" paragraph names them and points to §3.9–§3.13 / §5.9–§5.11 / F10–F13.
- **Attach step insertion, not bolt-on.** The disclosure step was inserted as §2.3 step 2 in the existing foreign-activation inaugural numbered procedure (steps 3–6 renumbered to 4–7 to keep the ordered list intact); the Algorithm summary, a new Rule §3.9, a Verify check §4.6, and catalogue entry A9 were added to match the skill's existing five-part rhythm.

## Files changed (exact paths)

- `src/packages/cnos.core/skills/agent/activate/SKILL.md` — Core Principle (case-first framing + observed-failure citation); Algorithm (Step 0 — resolve case); **§2.0** new (decision procedure Q1–Q5 + per-case ideal behavior A–E + cross-cutting-invariants pointer); Rules **§3.9–§3.13** (quiet / no-confabulated-state / no-confabulated-identity / attach-is-not-silent / identity-self-contained); Verify **§5.8–§5.13** (exactly-one-case, `ACTIVATION_STATE_CONFABULATION`, quiet-check, `SILENT_MUTATION`, identity-content, live-state); Failure catalogue **F10–F13**; References — new "Activation-case taxonomy (cnos#686)" subsection; section-manifest updated (`unfold-case-procedure`).
- `src/packages/cnos.core/skills/agent/attach/SKILL.md` — Core Principle (silent-mutation discipline); Algorithm (disclose-before-init); **§2.3 step 2** new (consequence disclosure + confirmation gate; steps 4–7 renumbered); Rule **§3.9** (disclose consequences before inaugural attach); Verify **§4.6** (consequence-disclosure check); catalogue **A9** (`SILENT_MUTATION`); References — Peer-skill entry cross-refs activate §2.0; new "Activation-case taxonomy and channel transport" subsection (#686, #684).
- `.cdd/unreleased/686/self-coherence.md` — this receipt.

`PERSONA.md` — **not touched** (operator-gated, per issue scope).

## AC1–AC7 → evidence map

| AC | Status | Evidence |
|---|---|---|
| **AC1** — exactly-one-case resolution (total + disjoint) | **met** | `activate` §2.0 decision procedure: Q1→Q5 ordered questions, each with a single observable discriminator; "resolves **exactly one** of five cases … total … disjoint; run in order, first that resolves wins." Every leaf terminates in exactly one of {A,B,C,D,E}. Verify §5.8 drives the five fixtures (bare box→E; `cn-sigma`→A; registered-with-surfaces→B; attachable-unregistered→C; ordinary/peer→D) and names the zero-or-≥2 outcome as a binding finding. |
| **AC2** — repo identity read, not assumed | **met** | §2.0 Q2: "Identify the repo **from the repo itself** — `git … remote get-url origin`, README, presence of `.cn-sigma/`/`.cdd/`. Never from assumption or memory. Multiple/ambiguous → **disambiguate with the activator** … do not silently pick. Only after the repo is named from a read step → Q3." §5.8 closes with the AC2 disambiguation check; §2.0 ❌ example covers the silent-pick violation. |
| **AC3** — quiet activation | **met** | Rule §3.9 "Activate quietly" (output = resolved case + one action; no cycle/release/pause/HEAD dump; binding). §2.0 per-case behaviors each say "Quiet" for A/B. Verify §5.10 (grep output; pause/`cycle/NNN`/"release X merged"/HEAD-as-current → fails). F12. |
| **AC4** — no-confabulated-state / identity (Case E = intro+question only) | **met** | §2.0 Case E: "minimal canonical identity only + ask … Do **not** run the six-item load … Do **not** assert doctrine, cycles, releases, HEAD, dispatch posture, or a pause … Output = minimal intro + one open question, nothing else." Rule §3.10 (state read from repo or `UNKNOWN`; repo wins). Verify §5.9 names `ACTIVATION_STATE_CONFABULATION`; F10. |
| **AC5** — attach discloses consequences | **met** | `attach` §2.3 step 2: emit the **exact list** (channel/log surface, registry entry + cursor, wake wiring) **before any mutation**; attach **only on confirmation**; non-confirmation → `deferred to operator`, no surface created; "Mutating … without this prior disclosure is `SILENT_MUTATION`." Rule §3.9, Verify §4.6, catalogue A9. Mirrored in `activate` §2.0 Case C + Rule §3.12 + Verify §5.11 + F11. #684 cross-referenced for transport, not re-specified. |
| **AC6** — identity content (self-contained, non-dualist, operator-free) | **partial — gated on PERSONA.md** | Rules §3.11 (never confabulate identity; no operator identity; no dualist gloss) + §3.13 (self-contained: OSS engineer who shipped cnos and other OSS projects; not an "assistants team," not operator-defined). Verify §5.12 + F13 encode the check. **Dependency flagged in-skill and here:** the recited identity *content* lives in `cn-sigma:.cn-sigma/spec/PERSONA.md`, which is operator-gated and out of this cell's scope; §5.12 explicitly states the check "also functions to force the source to be corrected." This cell makes activation *read* identity correctly and *specifies* the content discipline; it cannot fully satisfy AC6 until the PERSONA.md correction lands. |
| **AC7** — attached/home report live state (repo wins) | **met** | §2.0 Case A ("Read home state **from the repo** … never from memory; report current state + the next step") and Case B ("read **current channel state from the repo** … report current state + the potential next step"). Rule §3.10 ("repo wins"). Verify §5.13 (reported state == repo state; divergence → fail). |

## Known gaps / debt

1. **AC6 gated on operator-approved `PERSONA.md` content correction.** Out of this cell's scope by issue design (constitutive file, operator-gated, tracked separately). The skill prose specifies the identity discipline and the check forces the source; full satisfaction awaits the content edit. Named dependency, not silently assumed met.
2. **Mechanical activation-lint harness deferred** (issue "Deferred"). The three named checks (`ACTIVATION_STATE_CONFABULATION`, `SILENT_MUTATION`, quiet-check) are specified in-skill with positive/negative fixtures (§5.8–§5.13, §4.6) and rely on review for v0; mechanization is a possible follow-on sub, not this cell.
3. **#685 folding/closure** (issue closure condition 5) is a repo/issue-management action for δ/operator — the skills already carry the "supersedes #685" pointer (References). Not a file-write this cell performs.
4. **Renderer §4 deliberately untouched.** The case procedure is upstream of the six-step load block the renderer parses; §4.1/§4.2/§4.3/§4.4 are byte-unchanged, so `cn activate`'s observable output contract is preserved. Wiring the renderer to *emit* case-resolution guidance (if ever desired) is out of scope.
5. **No `.go` / no `transitions.json` / no label changes** — pure skill-prose diff (two `.md` files + this receipt).

## Verdict

**CONVERGED.**

Rationale: AC1, AC2, AC3, AC4, AC5, AC7 are fully encoded in the two skills with binding rules, fixture-driven checks, and named-violation catalogue entries; the five cases are total and disjoint; the attach consequence-disclosure step is in place with `SILENT_MUTATION` as its binding check and #684 cross-referenced for transport. AC6 is **partial by construction** — the identity *content* it grades is operator-gated (`PERSONA.md`), explicitly out of scope; this cell ships the reading discipline and the forcing check, and flags the dependency both in-skill (§5.12) and here. That partial is the designed boundary of the cell, not an unmet obligation of it. Integration preserved the renderer contract, existing section numbering, and house prose style; no foreign-looking section was bolted on. Remaining items are declared debt (mechanical lint, #685 closure, PERSONA content), all carried explicitly.
