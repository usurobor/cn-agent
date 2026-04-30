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


