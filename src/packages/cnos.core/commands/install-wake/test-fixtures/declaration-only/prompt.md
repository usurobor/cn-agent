# Synthetic test fixture — declaration-only

This prompt body is never inlined into any rendered workflow because the
manifest's `activation_state` is `declaration-only` and the renderer
refuses (exit 3) before reaching the prompt-inlining stage.

Existence required by the renderer's prompt-template existence check
(per wake-provider/SKILL.md §2.1: `prompt_template` is a required field
and the file must exist at the resolved path); content intentionally
trivial.
