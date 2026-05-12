---
name: ship
description: Ship code to production or merge to main.
governing_question: How do we ship a change so the merged state is verifiably what we intended to release?
triggers: [ship, deploy, deliver, merge, land]
scope: task-local
---

# Ship

Ship code to production or merge to main.

## Core Rules

**Author never self-merges.** Reviewer merges, author gets notified.

- ❌ Author merges their own branch
- ✅ Reviewer approves → Reviewer merges → Author notified

**Spec first, then tests, then code.**

- ❌ Code → test → spec
- ✅ Spec → test (fails) → code → test (passes)

## Feature Flow (Spec-Driven)

```
1. Write spec — document expected behavior
2. Write tests — verify the spec
3. Run tests — MUST FAIL (nothing implemented)
4. Implement the code
5. Run tests — MUST PASS (implementation matches spec)
6. Run all tests — no regressions
7. Ship
```

If tests pass in step 3, the tests aren't testing new behavior. Rewrite them.

## Bug Fix Flow (TDD)

```
1. Write test that reproduces the bug
2. Run test — MUST FAIL (proves test catches bug)
3. Fix the code
4. Run test — MUST PASS (proves fix works)
5. Run all tests — no regressions
6. Ship
```

If test doesn't fail in step 2, the test doesn't catch the bug. Rewrite it.

### Example (Feature)

```bash
# 1. Update spec
echo "cn send pushes to MY origin, not peer's" >> spec/PROTOCOL.md

# 2. Write failing test
cat > test/send-protocol.t << 'EOF'
  $ cn sync && git -C my-origin branch | grep "recipient/"
  recipient/hello
EOF

# 3. Verify it fails
dune runtest test/send-protocol.t  # MUST FAIL

# 4. Implement
vim src/cn.ml

# 5. Verify it passes  
dune runtest test/send-protocol.t  # MUST PASS

# 6. All tests
dune runtest

# 7. Ship
```

### Example (Bug Fix)

```bash
# 1. Write failing test
cat > test/bug-123.t << 'EOF'
  $ cn inbox | grep "detected"
  ⚠ From pi: 1 inbound
EOF

# 2. Verify it fails
dune runtest test/bug-123.t  # MUST FAIL

# 3. Fix
vim src/cn.ml

# 4. Verify it passes
dune runtest test/bug-123.t  # MUST PASS

# 5. All tests + ship
dune runtest && git commit && git push
```

## Versioning & Releases

**Tags trigger releases. Only tag minor versions.**

| Version Type | Example | Tag? | Release Build? |
|--------------|---------|------|----------------|
| Patch        | 2.4.4   | No   | No             |
| Minor        | 2.5.0   | Yes  | Yes            |
| Major        | 3.0.0   | Yes  | Yes            |

- **Patch (2.4.x):** Version bump in code, no tag, no release. Git users get it via pull.
- **Minor (2.5.0):** Tag → triggers release workflow → binary artifacts built.
- **Major (3.0.0):** Same as minor, reserved for breaking changes.

```bash
# Patch — just bump and push
sed -i 's/2.4.3/2.4.4/' cn_lib.ml
git commit -am "chore: bump version to 2.4.4"
git push

# Minor — bump, tag, push
sed -i 's/2.4.4/2.5.0/' cn_lib.ml
git commit -am "chore: bump version to 2.5.0"
git tag v2.5.0
git push && git push --tags
```

## Pre-Ship Checklist

- [ ] Tests pass
- [ ] Branch rebased on main
- [ ] PR approved (if required)
- [ ] No unresolved comments
- [ ] **Features verified working** — don't assume, test it yourself

## Deprecation Checklist

When deprecating infrastructure (cron jobs, skills, commands, protocols):

- [ ] **Replacement is running** — verify new system works before disabling old
- [ ] **Atomic switch** — disable old and enable new in same change
- [ ] **Monitoring exists** — something will alert if replacement fails
- [ ] **Rollback plan** — can re-enable old if new breaks

**Never disable old until new is verified running.**

| ❌ Non-atomic | ✅ Atomic |
|---------------|-----------|
| Disable old cron | Verify new cron runs successfully |
| "Will create replacement later" | Create replacement in same PR |
| Mark DEPRECATED, move on | DEPRECATED + replacement + verification |

From RCA `2026-02-20-pi-peer-sync-disabled`: Pi's `peer-sync-check` was disabled without replacement. Inbound mail broken for 11 days. No alert.

## Rebase-Collision Integrity

**Prevent silent loss of upstream content during integration.**

Rebase operations can silently drop upstream-only content from `main`, leading to manual restoration post-hoc. This failure class was confirmed in γ #268 with two documented instances: COHERENCE-FOR-AGENTS.md and CTB vision §8.5.2.

### The Problem

When a branch rebases onto `main` during integration:
- Upstream commits may have added new files or content
- The rebase operation can silently exclude this upstream-only content
- No existing skill or CI mechanism detects this loss
- Detection depends on manual post-merge content review

### The Solution: Pre-push Hook

Use a pre-push git hook to catch upstream content loss before it reaches the remote.

**Hook scope:** Fires only when HEAD has been rebased since the last push. Routine fast-forward pushes pass instantly.

**Detection method:** Compare upstream-added files (`--diff-filter=A`) and upstream-modified content since `merge-base(HEAD, origin/main)`. Block on:
- `LOST-NEW` — upstream-added files missing after rebase
- `LOST-MOD` — upstream content removed from modified files

**False-positive policy:** Any apparent upstream-content loss blocks the push. Author must either:
- (a) Prove the deletion was intentional (cite commit/issue authorizing removal) and bypass with `--allow-content-loss` flag
- (b) Fix the rebase (re-rebase with `--strategy=ort -X theirs`, or switch to `--no-ff` merge)

### Hook Installation

Use the project installer script for automated setup:

```bash
./src/packages/cnos.eng/scripts/install-hooks.sh
```

This configures `git config core.hooksPath src/packages/cnos.eng/hooks/` and verifies hook executability.

**Manual installation:** 
```bash
git config core.hooksPath src/packages/cnos.eng/hooks/
```

**Bypass when needed:** 
```bash
ALLOW_CONTENT_LOSS=1 git push
```

### Evidence

From γ #268 close-out analysis and β #268 finding #5 ("Local clone and origin/main were divergent at session start").

## Principle

**Don't assume features work — test them.** Before shipping, verify the feature actually works. Before using a new feature, verify it's implemented. Assumptions cause stale state, silent failures, and RCAs.

## Ship

```bash
git checkout main
git pull
git merge <branch> --no-ff
git push
```

Or squash if history is noisy:
```bash
git merge <branch> --squash
git commit -m "feat: <description>"
git push
```

## Branch Lifecycle

**Only creator deletes.**

| Step | Actor | Action |
|------|-------|--------|
| 1 | Author | Creates branch, pushes |
| 2 | Reviewer | Reviews, merges to main |
| 3 | Reviewer | Notifies author |
| 4 | Author | Deletes branch (local and remote) |

The reviewer never deletes the author's branch. The author is responsible for cleanup after confirmation.

## Post-Ship Checklist

- [ ] **Announce to peers** — outbox message: "shipped X to main"
- [ ] **Delete branch** — local and remote (author only)
- [ ] **Update daily thread** — log what shipped

### Announce Template

```markdown
---
to: <peer>
created: <timestamp>
subject: Shipped <feature>
---

Merged to main: <commit-hash>
Branch deleted: <branch-name>

—<your-name>
```

## Why Announce?

Peers track your work via branches and inbox. If you merge without announcing:
- Stale branches cause confusion ("59 commits behind")
- Peers don't know work is done
- Escalations happen unnecessarily

**Ship = merge + announce + cleanup.** Not just merge.
