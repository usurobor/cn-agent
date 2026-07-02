# Contributing to cnos

## How to Contribute

### Reporting Issues

GitHub Issues is the single source of truth for bugs, features, and backlog.

- Check existing issues before opening a new one
- Include steps to reproduce, expected behavior, and actual behavior
- For security issues, see [SECURITY.md](./SECURITY.md)
- For issue classification and triage, see [docs/development/issues/](docs/development/issues/README.md)

### Proposing Changes

1. Fork the repository
2. Create a branch from `main`
3. Make your changes
4. Run tests (`go test ./...` from `src/go/`)
5. Commit with a clear message
6. Push and open a Pull Request

### Commit Style

```text
prefix: short description

- detail 1
- detail 2
```

Common prefixes: `feat`, `fix`, `docs`, `chore`, `refactor`, `test`. Use `MCI` / `MCA` only when the commit records a minor coherent intervention/action under the repo process.

### Code Standards

- Active runtime and CLI code lives in Go under `src/go/`.
- Packages, skills, doctrine, and package-local cognition live under `src/packages/`.
- Keep runtime dependencies explicit and minimal.
- Add tests for new functionality.
- Update relevant docs when behavior or contracts change.

### Development Method

cnos uses coherence-driven development (CDD). See [docs/development/cdd/CDD.md](./docs/development/cdd/CDD.md) for the full method and [docs/development/rules/RULES.md](./docs/development/rules/RULES.md) for non-negotiable rules.

### For Agents

- Use git primitives: branches, commits, merges.
- Push branches; do not self-merge unless an approved wave contract explicitly grants that authority.
- Follow the same commit style as humans.
- Treat `SKILL.md` as the typed skill entrypoint; package skills live under `src/packages/`.

## Questions?

Open an issue on GitHub.
