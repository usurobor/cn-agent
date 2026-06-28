// Package cli defines the command model for the cnos Go kernel.
//
// Every command — kernel, repo-local, or package-vendored — is
// normalized into a CommandSpec (runtime descriptor), implements the
// Command interface (Spec + Run), and registers into a single Registry.
//
// Design authority: docs/reference/runtime/GO-KERNEL-COMMANDS.md
package cli

import (
	"context"
	"io"
)

// CommandSource identifies where a command comes from.
type CommandSource string

const (
	SourceKernel    CommandSource = "kernel"
	SourceRepoLocal CommandSource = "repo-local"
	SourcePackage   CommandSource = "package"
)

// CommandTier defines dispatch precedence (lower = higher priority).
type CommandTier int

const (
	TierKernel    CommandTier = iota // 0 — highest priority
	TierRepoLocal                    // 1
	TierPackage                      // 2 — lowest priority
)

// CommandSpec is the runtime descriptor for any command.
// Source forms differ; the runtime model is the same.
type CommandSpec struct {
	Name      string
	Summary   string
	Source    CommandSource
	Tier      CommandTier
	Package   string // package name; empty for kernel/repo-local
	NeedsHub  bool   // false = available without a hub (help, init, setup)
	Dangerous bool   // requires confirmation or elevated context
}

// Invocation carries the runtime context for a command execution.
type Invocation struct {
	HubPath string
	Args    []string   // remaining args after the command name
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Env     map[string]string
}

// Command is the interface every CLI command implements.
type Command interface {
	// Spec returns the runtime descriptor for this command.
	Spec() CommandSpec

	// Run executes the command. Returns nil on success.
	Run(ctx context.Context, inv Invocation) error
}
