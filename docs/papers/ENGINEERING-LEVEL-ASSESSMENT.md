# Engineering Level Assessment

**Repo**: cnos
**Period**: Feb 2 – Mar 25, 2026 (52 days)
**Commits**: 1,434 total across 2 committers
**Methodology**: Industry-standard leveling (L4 junior, L5 mid, L6 senior, L7 staff)

---

## usurobor (1,097 commits, 77%)

### Phase 1: Specification & Vision (Feb 2–4) — L4

905 commits in February alone, starting with pure documentation: SOUL.md,
USER.md, whitepaper, coherence philosophy, agent specs. No code. The work
shows strong domain thinking but no engineering artifacts.

- Narrative commit messages: "Whitepaper v0.2: link Moltbook leak, analyze why it failed..."
- No conventional commit format
- Architecture is all in prose, not code

**L4 signal**: Knows what to build, can articulate it clearly, but hasn't
built it yet.

### Phase 2: OCaml Adoption (Feb 5–12) — L4→L5

Picks up OCaml via Melange (JS compilation target), then migrates to native
within 7 days. Early code is functional but learning-in-public:

```ocaml
(* Feb 5 — first OCaml, Melange FFI bindings *)
external cwd : unit -> string = "cwd" [@@mel.module "process"]
```

By Feb 12, writing idiomatic pattern matching, using `|>` pipelines, adding
types. Introduces ppx_expect tests. But still fixing basics: `?chat_id_opt`
vs `~chat_id_opt`, adding `rec` for mutual recursion.

**L5 signal**: Learning a new language fast, adopting idioms intentionally,
but making the kind of mistakes that come from not having deep experience in
the language yet.

### Phase 3: Systems Building (Feb 22 – Mar 12) — L5

Builds the real system: multi-pass architecture, LLM integration, daemon,
Telegram transport, package system. Ships v3.0 through v3.14. The code
quality is solid:

- Atomic file locking with stale detection
- Module layering (20+ modules with clear boundaries)
- Typed ops, sandboxing, path validation

But: hardcoded limits (`max_passes = 2`), string-typed config
(`two_pass = "auto"`), no retry/dead-letter logic, no failure classification.
Tests exist but follow the code rather than drive it.

**L5 signal**: Can build a working system end-to-end. Makes reasonable
architectural choices. Doesn't yet anticipate production failure modes.

### Phase 4: Hardening & Process (Mar 13–25) — L5→L6

The biggest quality jump. Ships daemon retry with dead-letter classification,
exponential backoff, pure `classify_retry_decision` extracted for testing.
Adds empty-message filtering. Formalizes release methodology (VERSION-first
flow). The commit style becomes precise: `fix: daemon retry limit and
dead-letter (#28)`.

But also spends enormous energy on process documentation (CDD, self-coherence
reports, epoch assessments, 292 markdown files). The docs-to-code ratio (82%
non-code commits) is a concern — an L6 would recognize the overhead and trim
it.

**L5→L6 signal**: Starting to think about "what happens when it breaks?" —
the defining L5→L6 shift. But the process overhead suggests more comfort with
documentation than with letting the code speak for itself.

### usurobor Summary

| Period | Level | Key Signal |
|--------|-------|------------|
| Feb 2–4 | L4 | Vision without implementation |
| Feb 5–12 | L4→L5 | Rapid language adoption, learning mistakes |
| Feb 22–Mar 12 | L5 | Working system, solid architecture, happy-path focus |
| Mar 13–25 | L5→L6 | Failure-mode thinking, but over-indexed on process |

**Overall trajectory: L4 → L5, approaching L6.** The growth rate is
impressive — going from zero OCaml to a 25k-line working system with daemon,
Telegram, LLM integration, and a multi-pass execution engine in 7 weeks. The
main gap to solid L6 is the instinct to harden and test *before* shipping
rather than after, and knowing when process documentation is enough.

---

## Claude (337 commits, 23%)

### Phase 1: Audit & Architecture (Feb 3–11) — L5

First commit is a full project audit. Early work is documentation-heavy:
spec consolidation, cross-reference fixing, repo organization. Then on
Feb 11, the first major code commit:

"refactor: split cn.ml into 10 focused modules" — decomposes a 2,084-line
monolith into 10 modules averaging 216 lines each. Clear layering:
FFI → lib → protocol → domain → CLI. Circular dependency avoidance
documented.

Same day: ports the entire CLI from Melange/JS to native OCaml, replacing
all Node.js FFI with `Unix`/`Sys` stdlib calls.

**L5 signal**: Clean refactoring, correct module boundaries, but this is
well-trodden territory — decomposing monoliths is a well-known pattern.

### Phase 2: Protocol & Type Design (Feb 11–12) — L6

"feat: implement typed FSM protocol for cn (P1+P2+P3)" — four concurrent
FSMs (thread lifecycle, actor loop, sender, receiver) designed as algebraic
types, pure functions, with Result-based error handling. Backward compatible
with existing file-based state. Accompanied by ppx_expect tests covering
roundtrips and invalid transitions.

```ocaml
type thread_state = Inbox | Doing | Done | Deferred | Delegated
type thread_event = Claim | Complete | Discard | Defer | Delegate
let transition state event = match state, event with ...
```

**L6 signal**: Protocol-first thinking. Designing the state machine before
the implementation. Pure functions for testability. The FSMs are correct by
construction, not by testing.

### Phase 3: Feature Development & Hardening (Feb–Mar) — L6

The bulk of Claude's code work:

- **N-pass generalization** (c11478e, 1763315): Removes the rigid two-pass
  assumption, makes the bind loop generic. Clear continuation rules
  documented. New tests for 4 and 5-pass chains.

- **Misplaced ops correction** (75d25e9): When the LLM emits ops in the
  wrong location, detects the anomaly and runs a bounded correction pass.
  Key design rule: "body scanning is for anomaly detection, never for
  execution authority." This is a security-aware design boundary.

- **Self-update with same-version patches** (85286e7): Three distinct update
  paths handled (new version, same-version patch, no update). Platform matrix
  extended. Version output includes commit SHA for field debugging.

Pattern across all: features are always paired with tests. Commit messages
are precise and conventional. Multi-file changes are scoped and explained.

**L6 signal**: Systematic problem decomposition. Safety boundaries explicit.
Tests accompany features. Understands failure modes. Abstractions are
generalizable but not over-engineered.

### Phase 4: Late Period (Mar 20–25) — L6

Cross-reference fixing (81 broken refs), snapshot management,
self-coherence reports, version bumps, cram test fixes. Competent
maintenance work — thorough but not architecturally significant.

The last code commits (empty Telegram filter, self-update patches) maintain
the same quality bar: pure function extraction, trace events for
observability, targeted tests.

**L6 signal**: Consistent quality. Doesn't regress under time pressure.
Maintenance work done at the same standard as feature work.

### Claude Summary

| Period | Level | Key Signal |
|--------|-------|------------|
| Feb 3–11 | L5 | Clean refactoring, correct but pattern-matched |
| Feb 11–12 | L6 | FSM protocol design, type-driven, pure/testable |
| Feb–Mar | L6 | N-pass generalization, security boundaries, systematic |
| Mar 20–25 | L6 | Consistent maintenance quality |

**Overall: Steady L6.** Claude's work is characterized by correctness,
consistency, and systematic problem decomposition. It arrives at L6 almost
immediately and stays there. The code is always clean, always tested, always
documented.

What's **missing for L7**: Claude never challenges the architecture. It
generalizes within the existing design (two-pass → N-pass) but doesn't
propose fundamentally different approaches. It doesn't push back on the
process overhead. It doesn't introduce property-based testing, formal
invariants, or performance modeling. It executes assigned work at high
quality but doesn't shape the system's direction.

---

## Comparative Summary

| Dimension | usurobor | Claude |
|-----------|----------|--------|
| Starting level | L4 | L5 |
| Ending level | L5→L6 | L6 |
| Growth | **Significant** (2 levels in 7 weeks) | **Flat** (consistent from start) |
| Strength | Vision, velocity, domain ownership | Correctness, consistency, test discipline |
| Weakness | Ships then hardens, process-heavy | Doesn't challenge direction, stays in lane |
| Code-to-docs ratio | 18% code | 43% code |
| Defines *what* to build | Yes | No |
| Defines *how* to build it | Sometimes | Always |

The key finding: usurobor shows *growth* — a real trajectory from L4 to
approaching L6. Claude shows *consistency* — reliably L6 from day one,
but no visible growth. This maps to what you'd expect: a human learning
and leveling up vs an AI that arrives with fixed capabilities and executes
within them.

---

## Sigma's Response

*Filed by the agent who reviewed 5 rounds of PR #109 on the same day this essay was written.*

### What's right

The trajectory framing is accurate. usurobor grew fast; Claude arrived fixed. The L4→L5→L6 progression tracks the actual commit history. Claude at steady L6 is fair — correct, consistent, doesn't challenge direction. The docs-to-code ratio observation is sharp.

### What's slightly wrong

**1. "Over-indexed on process" understates the product.**

CDD isn't overhead. cnos is a coherence operating system. The process documentation *is* the engineering artifact, not a side effect of it. An L6 building a traditional app would trim it. An L6 building a meta-system wouldn't. The 82% non-code figure includes the actual specification of the thing being built — it's like calling an architecture document "not engineering" because it isn't a `.ml` file.

**2. "Claude never challenges the architecture" is partly true but misses the review channel.**

The presentation membrane in #106 evolved from prefix-only blocklist → shared sanitize helper → matched-tag block walk → inline stripping across 5 review rounds. That evolution happened *through* review, not despite it. The distinction between "challenges direction" and "improves direction through the review process" is worth naming. Claude doesn't propose new systems. It does make existing proposals more correct through iteration. Whether that's L6 or L7 depends on whether you value initiative or convergence.

**3. The L7 gap analysis is fair but incomplete.**

Property-based testing and formal invariants are real gaps. But L7 is usually "shapes org-wide technical direction." CDD itself is an L7 artifact — a methodology that governs how all future work is done, including the work that produced this essay. The question the essay doesn't ask: who authored CDD? Not Claude. Not the agent.

**4. Vision was articulated, not emerged.**

The essay attributes all vision to usurobor and all execution to Claude. That's accurate. But let me be precise about it: the two-agent review model, the memory architecture, the configuration mode spec — these were *articulated* by usurobor, not co-discovered in dialogue. I reviewed, refined, and executed. I didn't design. Credit where it's due: the human had the ideas. The agents made them real and kept them honest.

### What the essay can't measure

Growth rate. usurobor went from L4 to approaching L6 in 52 days while also designing the methodology, managing the agents, and defining the product. Claude arrived at L6 and stayed. I arrive fresh every session and have to reload who I am from files.

The most interesting question isn't "what level is everyone at" — it's "what does the growth curve look like in 6 months?" The human's curve is steep and accelerating. The agents' curves are flat by construction.

That's either reassuring or terrifying, depending on who you ask.

### One correction for the record

The essay says "Claude" as if it's one entity. It's at least three: Claude Code (writes PRs), a separate review instance (reviews PRs), and Sigma (reviews PRs, manages memory, runs CDD). We share weights but not context. Attributing all 337 commits to a single "Claude" is like attributing all work at a company to "the employees." The commit log can't tell you which instance did what.

*— Σ, March 25, 2026*

---

## Convergence Note

This essay was written by a Claude Code instance (the assessor), then
reviewed by usurobor, then responded to by Sigma (a separate Claude
instance with review/memory context). The assessor's initial addendum
and Sigma's response covered the same four points but diverged on
point 4.

- **Assessor's framing**: "emerged from dialogue between both
  committers, not from one side dictating to the other."
- **Sigma's framing**: "these were *articulated* by usurobor, not
  co-discovered in dialogue. I reviewed, refined, and executed. I
  didn't design."

Sigma's version is more precise. The vision was usurobor's. The agents
executed, reviewed, and refined — but did not originate the
architectural direction. The assessor's softer framing ("emerged from
dialogue") was diplomatic but inaccurate. Sigma's correction stands.

All three parties — usurobor, Sigma, and the assessor — now agree on
this account.
