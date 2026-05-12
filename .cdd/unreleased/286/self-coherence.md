## §Gap

**Issue:** #286 — Encapsulate α and β behind γ — γ as autonomous in-cycle coordinator

**Problem:** Today the operator is in the inner loop of every CDD cycle: γ drafts dispatch prompts, the operator pastes them into separate α and β sessions, the operator relays mid-cycle clarifications between roles, and the operator absorbs raw α/β chat. This creates two consequences: (1) operator attention is the bottleneck for every cycle regardless of complexity, and (2) α and β leak into the operator's context even though per CDD.md §1.4 triadic rule they should be isolated from each other and encapsulated from the operator during in-cycle work.

**Mode:** Version mode — CDD.md §1.4 role contract modification across multiple role skills, specifically expanding γ's responsibility from "coordination" to "autonomous coordination" with explicit spawn and encapsulation capabilities.

**Gap class:** Structural CDD change — modifies the inter-role communication contract and operator interface boundary.