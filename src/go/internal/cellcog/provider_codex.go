package cellcog

import (
	"context"
	"fmt"
	"time"
)

// CodexCLI rents workspace cognition from the Codex CLI. Stateless: a Work
// call is one fresh bounded `codex exec`; no rollout state is kept.
type CodexCLI struct {
	Model   string        // exact model id; required
	Bin     string        // default "codex"
	Timeout time.Duration // default 10m
}

func (CodexCLI) Name() string { return "codex-cli" }

// CodexArgv is the exact invocation, built purely so it can be asserted
// without spawning anything.
//
//   - `--ephemeral` keeps the run stateless; without it Codex persists rollout
//     state between invocations.
//   - `--ignore-user-config` and `--ignore-rules` keep ambient configuration
//     from becoming a second, unreceipted component definition. Authentication
//     stays ambient, as intended — credentials are the operator's, never the
//     cell's, and never enter a receipt.
//   - `--sandbox workspace-write` is the provider's own write confinement; it
//     is the only adapter here with provider-enforced scoping, and even that
//     is not claimed as OS confinement by this package.
func CodexArgv(model, dir string) []string {
	return []string{
		"exec",
		"--model", model,
		"--ephemeral",
		"--ignore-user-config",
		"--ignore-rules",
		"--sandbox", "workspace-write",
		"--skip-git-repo-check", // the worktree is disposable by construction
		"--cd", dir,
		"-", // prompt on stdin
	}
}

func (c CodexCLI) Work(ctx context.Context, dir, prompt string) error {
	if dir == "" {
		return fmt.Errorf("codex-cli: Work needs a directory")
	}
	bin := c.Bin
	if bin == "" {
		bin = "codex"
	}
	return runCLI(ctx, bin, dir, prompt, CodexArgv(c.Model, dir), c.Timeout)
}
