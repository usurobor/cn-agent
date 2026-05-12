# CDD Activation Templates

Copy-pasteable templates for wiring CDD into a tenant repository.
No edits are needed beyond replacing secret names if you diverged from the
canonical `CDD_` namespace.

---

## telegram-notifier/

Reference implementation of the CDD notification adapter
(`cdd/activation/SKILL.md §10.2`) for Telegram.

| File | Purpose |
|---|---|
| `notify.sh` | Shell script: posts to Telegram Bot API for all 4 CDD events |
| `cdd-notify.yml` | GitHub Actions workflow: triggers notify.sh on the right events |
| `README.md` | Bot registration + secret setup walkthrough |

Copy `notify.sh` to `.github/cdd/notify.sh` and `cdd-notify.yml` to
`.github/workflows/`. See `telegram-notifier/README.md` for the full walkthrough.

## github-actions/

CI workflow templates for CDD artifact governance.

| File | Purpose |
|---|---|
| `cdd-artifact-validate.yml` | Validates `.cdd/` structure on every push to `cycle/**` and `main` |
| `cdd-cycle-on-merge.yml` | Runs project test suite + emits cycle-merge notification on merge to main |

Copy both files to `.github/workflows/` and replace the test command placeholder
in `cdd-cycle-on-merge.yml` with your project's actual test runner.

---

*Governing spec: `cdd/activation/SKILL.md §9–§10`*
