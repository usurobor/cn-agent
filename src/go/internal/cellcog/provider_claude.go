package cellcog

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Adapter defaults. The bounds exist because a rented seat is the one part of
// an episode the runtime cannot predict: the running system's wake had no
// timeout, no turn cap and no output cap, and a runaway cell burned its whole
// job budget. A cell bounds what it rents.
const (
	defaultClaudeBin = "claude"
	defaultTimeout   = 10 * time.Minute
	defaultMaxBytes  = 4 << 20 // matches the kernel's aggregate artifact bound
	maxStderrBytes   = 8 << 10 // diagnostics only
	waitDelay        = 2 * time.Second
)

// ClaudeCLI rents cognition from the Claude Code CLI. It is an adapter: the
// subprocess is the only effect in this package, and it produces text — the
// seat, not the provider, decides what that text means.
//
// The seat is text-in/text-out on purpose. It passes no tool, permission, or
// workspace flags, so this provider cannot touch the filesystem; a producing
// alpha that edits a working copy arrives with the matter substrate (G1),
// where the workspace is materialized outside the kernel and passed in.
type ClaudeCLI struct {
	Bin      string        // default "claude"
	Timeout  time.Duration // default 10m
	MaxBytes int           // default 4 MiB
}

func (ClaudeCLI) Name() string { return "claude-cli" }

func (c ClaudeCLI) Complete(ctx context.Context, prompt string) (string, error) {
	return c.run(ctx, "", prompt)
}

// Work implements Coder: the same binary, granted file tools and confined to
// dir. Tools are FILE tools only — read, write, edit, search. No shell: a
// producing seat needs to change files, and a shell is a capability nobody
// asked for. Running tests is a mechanical step or a reviewer's job, not
// something the produced-work seat gets to do to the host.
func (c ClaudeCLI) Work(ctx context.Context, dir, prompt string) error {
	if dir == "" {
		return fmt.Errorf("cellcog: Work needs a directory")
	}
	_, err := c.run(ctx, dir, prompt)
	return err
}

func (c ClaudeCLI) run(ctx context.Context, dir, prompt string) (string, error) {
	bin := c.Bin
	if bin == "" {
		bin = defaultClaudeBin
	}
	if _, err := exec.LookPath(bin); err != nil {
		return "", fmt.Errorf("%q not found in PATH: %w", bin, err)
	}
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	maxBytes := c.MaxBytes
	if maxBytes <= 0 {
		maxBytes = defaultMaxBytes
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	args := []string{"-p", "--output-format", "text"}
	if dir != "" {
		args = append(args, "--permission-mode", "acceptEdits",
			"--allowedTools", "Read,Write,Edit,Glob,Grep")
	}
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = dir // "" keeps the caller's cwd; a text-only Complete writes nothing
	cmd.Stdin = strings.NewReader(prompt)
	stdout := &boundedBuffer{max: maxBytes}
	stderr := &boundedBuffer{max: maxStderrBytes}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	// Killing the provider does not by itself unblock Wait: anything it spawned
	// inherits the output pipe and holds it open, so a timeout would wait for
	// the grandchild instead of the child it just killed. WaitDelay bounds that
	// second wait — after it, the pipes are closed and Wait returns.
	cmd.WaitDelay = waitDelay

	err := cmd.Run()
	if ctxErr := ctx.Err(); ctxErr != nil {
		return "", fmt.Errorf("%s did not answer within %s: %w", bin, timeout, ctxErr)
	}
	if err != nil {
		return "", fmt.Errorf("%s failed: %w (stderr: %s)", bin, err, strings.TrimSpace(stderr.String()))
	}
	if stdout.truncated {
		return "", fmt.Errorf("%s answered with more than %d bytes", bin, maxBytes)
	}
	return stdout.String(), nil
}

// boundedBuffer captures at most max bytes and remembers that it had to stop.
// It never reports a short write, so the bound fails the answer here rather
// than killing the child mid-stream with a broken pipe.
type boundedBuffer struct {
	max       int
	buf       bytes.Buffer
	truncated bool
}

func (b *boundedBuffer) Write(p []byte) (int, error) {
	switch room := b.max - b.buf.Len(); {
	case room >= len(p):
		b.buf.Write(p)
	case room > 0:
		b.buf.Write(p[:room])
		b.truncated = true
	default:
		b.truncated = true
	}
	return len(p), nil
}

func (b *boundedBuffer) String() string { return b.buf.String() }
