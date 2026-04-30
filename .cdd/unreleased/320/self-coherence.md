## Gap — #320 cn activate

**Issue:** #320 — feat(cli): cn activate — bootstrap prompt generation from hub state

**Version/mode:** MCA — adds `cn activate` as a kernel CLI command that generates a bootstrap prompt from local hub state. No runtime is started; no model is invoked.

**Incoherence closed:** `cn init` and `cn setup` exist, but no CLI path helps the operator assemble the activation prompt a capable model needs to orient to the hub. Operators must manually know which hub files matter and assemble the prompt themselves, producing inconsistent results.

---

## Skills

**Tier 1 (always loaded):**
- `CDD.md` (canonical lifecycle)
- `cdd/alpha/SKILL.md` (α role surface)

**Tier 2 (general engineering):**
- `cnos.core/skills/design/SKILL.md`
- `cnos.core/skills/write/SKILL.md`
- `cnos.eng/skills/eng/tool/SKILL.md`
- `cnos.eng/skills/eng/test/SKILL.md`
- `cnos.eng/skills/eng/ux-cli/SKILL.md`
- `cnos.eng/skills/eng/go/SKILL.md`

**Active constraints applied:**
- `design`: prompt generation separated from runtime/daemon/dispatch; `cn activate` is a flat kernel command (not `cn agent activate`); `cli/` owns dispatch only (`eng/go §2.18`)
- `write`: help text and prompt prose say what the command does without overclaiming runtime behavior
- `tool`: stdout = prompt only; stderr = diagnostics; fail-fast on invalid hub
- `test`: invariants named before tests written; positive and negative cases both covered; no model invocation in tests
- `ux-cli`: error messages are causal and actionable; `✗` and `→` glyphs used per semantic table
- `go`: `eng/go §2.17` (Parse/Read split not applicable — no schema parsing introduced); `eng/go §2.18` (domain logic in `internal/activate/`, thin wrapper in `internal/cli/cmd_activate.go`)
