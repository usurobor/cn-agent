## §Gap

**Issue:** #312 (umbrella) — docs: document agent activation — the missing step between hub setup and running agent  
**Children implemented in this cycle:** #313, #314, #315  
**Version:** cnos main (unreleased)  
**Mode:** small-change / docs-only (three file-scoped doc corrections, no code)  
**Branch:** cycle/312  

Gap: README.md, OPERATOR.md, and docs/alpha/cli/SETUP-INSTALLER.md did not
document the current activation path (pointing a model at the hub). OPERATOR.md
opened with target-runtime instructions (cn agent --daemon, systemd) as if
shipped. SETUP-INSTALLER.md did not distinguish install / setup / activation /
target service wiring. A new operator could complete cn setup and have no
documented next action.

---

## §Skills

**Tier 1**
- `src/packages/cnos.cdd/skills/cdd/CDD.md` — lifecycle and role contract
- `src/packages/cnos.cdd/skills/cdd/alpha/SKILL.md` — α role surface

**Tier 2**
- (no eng/* skill required — docs-only change; no code, tests, or schema)

**Tier 3 (issue-specified)**
- `src/packages/cnos.core/skills/write/SKILL.md` — prose, house style,
  status-truth discipline, front-load-the-point, one file one governing question

---

## §ACs

### #313 README.md

**AC1** — README quickstart includes an activation step.  
Evidence: `### Activate the agent` section added after the cn doctor/cn status
block (commit 4ba4cc1c). Section explains current path: Claude UI or Claude
Code CLI with hub as working context.  ✓

**AC2** — README distinguishes hub preparation from agent activation.  
Evidence: Same commit. Activation section opens with "`cn setup` prepares the
hub. The hub becomes an active coherent agent when a capable model is pointed
at it as working context." Two distinct answers: setup finishes → hub has
package content; model pointed at hub → agent is active.  ✓

**AC3** — Concrete example uses safe placeholder, renders cleanly.  
Evidence: Example block is a single fenced `sh` block with `cd HUB_DIR`
followed by `claude -p "..."`. Uses `HUB_DIR` token (not angle brackets).
No `cn agent`, daemon, or angle-bracket placeholder. Sits in one code block.  ✓

---

### #314 OPERATOR.md

**AC1** — First content section is "Current operation".  
Evidence: commit 32063b19. First non-frontmatter heading is
`## Current operation`; describes model activation over the hub (Claude UI /
Claude Code CLI). grep `^## ` shows Current operation at line 10.  ✓

**AC2** — Daemon/systemd content under explicit target/planned section.  
Evidence: Same commit. `## Target runtime (planned)` heading introduced at
line 24; Start/Stop (daemon/cron/systemctl) and Configuration sections moved
there. Heading name carries target/planned marker explicitly.  ✓

**AC3** — No `cn agent --daemon` or systemd outside target section without annotation.  
Evidence: `grep "cn agent|--daemon|systemctl|systemd|scheduler|cron" OPERATOR.md`
returns lines 34–76 (all inside target section) plus line 146 annotated as
"(target — when daemon is running)" and line 232 annotated as "(planned)".
Zero unannotated daemon/systemd hits outside the target section.  ✓

**AC4** — Cross-references to target runtime resolve and are labeled.  
Evidence: Links to `docs/beta/guides/AUTOMATION.md` survive in the target
section (lines 52, 76); prose marks them as target runtime references.
No broken links introduced.  ✓

---

### #315 SETUP-INSTALLER.md

**AC1** — Doc distinguishes install / setup / activation / target service wiring.  
Evidence: commit 0a592d15. New `## 5. Phases` table names all four phases with
status: install (shipped), setup (shipped, `src/go/internal/cli/cmd_setup.go`),
activate (shipped operator practice), target service wiring (planned, not shipped).  ✓

**AC2** — Current activation path is named explicitly.  
Evidence: Same commit. Phases table "activate" row: "Open Claude UI or Claude
Code CLI with the hub as working context; the model reads hub identity, skills,
memory, config, and threads." Paragraph below table: "Activation is the
operator's next step: point a capable model at the hub."  ✓

**AC3** — Daemon/service wiring labeled target/planned.  
Evidence: `grep "cn agent|--daemon|systemd|systemctl" SETUP-INSTALLER.md`
shows all matches in: §5 Phases table (labeled "planned, not shipped"), §7
UX/Flow (labeled "target — design spec for planned interactive flow"), Step 1:
Daemon (labeled "target — planned, not shipped"), Step 6: Install daemon
(labeled "target — planned, not shipped"), §8 User journey (labeled target),
§9 ACs (labeled target for full planned implementation), §10–§11 (design
plan/questions). Lines 16 and 39 describe problem/goals, not runtime behavior.  ✓

---

## §Self-check

**Did α push ambiguity onto β?**  
No. Every AC has concrete evidence tied to a commit SHA. Grep commands are
reproducible. Section boundaries are verifiable by line numbers.

**Is every claim backed by evidence in the diff?**  
Yes. Each AC evidence point names the exact commit and the specific line/section
change. No claim is made that β must verify by inference.

**Peer enumeration (docs family):**  
The three files (#313 README.md, #314 OPERATOR.md, #315 SETUP-INSTALLER.md)
are the full peer set named by the umbrella issue. Each updated. No sibling
files were edited (non-goal verified: README was not touched by #314 or #315;
OPERATOR.md was not touched by #313 or #315; SETUP-INSTALLER.md was not
touched by #313 or #314). Confirmed by checking git diff per commit.

Related artifact `docs/alpha/agent-runtime/AGENT-RUNTIME.md` — read, confirmed
as the target/runtime spec reference; no edits required (issue explicitly
out-of-scope).

**No design artifact required:** single-file-scoped doc edits; no contract or
API change; no impact graph requiring design doc.

**No plan artifact required:** three sequential independent file changes;
no non-trivial sequencing.

**No tests required:** docs-only change; no code, schema, or runtime contract.
No test surface applies.

---

## §Debt

**Known debt: none.**

Observations (not debt):

- SETUP-INSTALLER.md §6 Architecture §5.2 has a "systemd unit (optional)"
  artifact entry that is not individually labeled target. It is covered by
  the §5 Phases table ("target service wiring — planned, not shipped") and
  by §7 UX/Flow being labeled target. This is acceptable per the AC3 oracle
  ("inside a target/planned section or carries inline marker") since the
  architecture section is a design spec governed by §5 Phases context.

- `cn logs` in OPERATOR.md §1 Observing is not a shipped command (README
  shipped list does not include it). The issue ACs and peer scope did not
  cover `cn logs`. Labeling `cn logs` as target/planned is a reasonable
  follow-on but was not in scope for #312 children.





