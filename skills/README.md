# skills/ - Coherent Agent Skills

This directory holds reusable skills for this Coherence Network repo.

Each skill lives under `skills/<skill-name>/` and MUST contain at least:

- `SKILL.md` — the skill's spec, with:
  - TERMS — what it assumes (files, tools, state).
  - INPUTS — what events or parameters it receives.
  - EFFECTS — which files or external services it touches.

Optional:

- Scripts or helpers invoked by the runtime or by katas.

To add a skill:

1. Create `skills/<skill-name>/`.
2. Add `skills/<skill-name>/SKILL.md`.
3. Reference the skill from specs (for example `spec/core/HEARTBEAT.md`).
4. Run `./setup.sh` so the updated specs and skill files are pushed and applied to the runtime.
