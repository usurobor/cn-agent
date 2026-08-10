package cellcog

import (
	"context"
	"fmt"
	"time"
)

// CodexCLI rents workspace cognition from the Codex CLI. Stateless: a Work
// call is one fresh bounded `codex exec`; no session is started or kept. The
// provider's own workspace-write sandbox confines writes to the working
// directory — the one adapter here with provider-enforced confinement.
type CodexCLI struct {
	Model   string        // exact model id; required
	Bin     string        // default "codex"
	Timeout time.Duration // default 10m
}

func (CodexCLI) Name() string { return "codex-cli" }

func (c CodexCLI) Work(ctx context.Context, dir, prompt string) error {
	if dir == "" {
		return fmt.Errorf("codex-cli: Work needs a directory")
	}
	bin := c.Bin
	if bin == "" {
		bin = "codex"
	}
	args := []string{
		"exec",
		"--model", c.Model,
		"--sandbox", "workspace-write",
		"--skip-git-repo-check", // the worktree is disposable by construction
		"--cd", dir,
		"-", // prompt on stdin
	}
	return runCLI(ctx, bin, dir, prompt, args, c.Timeout)
}
