## Gap

**Issue:** #320 — feat(cli): cn activate — bootstrap prompt generation from hub state

**Version/mode:** MCA — adds `cn activate` as a kernel CLI command that generates a bootstrap prompt from local hub state. No runtime is started; no model is invoked.

**Incoherence closed:** `cn init` and `cn setup` exist, but no CLI path helps the operator assemble the activation prompt a capable model needs to orient to the hub. Operators must manually know which hub files matter and assemble the prompt themselves, producing inconsistent results.
