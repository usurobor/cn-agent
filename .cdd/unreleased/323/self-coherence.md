# Self-Coherence — Cycle #323

## Gap

**Issue:** #323 — fix(activate): scanner misses threads/inbox/

**Problem:** `cn activate` scans `threads/in`, `threads/mail`, and `threads/archived` but does not scan `threads/inbox/`. The canonical reference hub (cn-sigma) uses `threads/inbox/` as its primary inbound message surface, making activation output blind to the actual inbox.

**Expected:** `cn activate` should scan `threads/inbox/` in addition to the current surfaces.

**Version/Mode:** MCA (Minimal Correct Action) — direct fix to existing scanner function.