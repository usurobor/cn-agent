# Review Checklist

Pre-submit and review checklist. Both author and reviewer must verify.

## P0 — Governance

- [ ] **No self-merge.** Author never merges own diff to main unless human explicitly instructs.
- [ ] **Branch current.** Rebased on main. No missing commits from main.

## P1 — Process

- [ ] **Author rebases.** Reviewer never rebases. Reviewer's time > author's time.
- [ ] **Only author deletes branch.** Reviewer returns branches, never deletes.

## P2 — Quality

- [ ] **Purpose.** Solves stated problem?
- [ ] **Simplicity.** Simplest solution?
- [ ] **Necessity.** No unnecessary additions?
- [ ] **Types.** Correct and semantic?
- [ ] **Edge cases.** Handled?
- [ ] **Tests.** Tested?
- [ ] **History.** Clean commits?

---

**Author:** Verify P0 + P1 before pushing branch.

**Reviewer:** Verify all. If P0 fails → REQUEST REBASE or reject.
