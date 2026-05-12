package dispatch

import (
	"fmt"
	"os/exec"
	"time"
)

// ClaudeBackend implements dispatch via Claude CLI
type ClaudeBackend struct{}

func (b *ClaudeBackend) Name() string {
	return "claude"
}

func (b *ClaudeBackend) IsAvailable() error {
	// Check if claude binary is available
	if _, err := exec.LookPath("claude"); err != nil {
		return fmt.Errorf("claude binary not found in PATH. Install Claude Code CLI first")
	}

	// TODO: Check if authenticated (could run `claude --version` or similar)
	// For now, just check binary exists

	return nil
}

func (b *ClaudeBackend) Dispatch(args *Args, prompt string) (*BackendResult, error) {
	startTime := time.Now()

	// Construct claude command with required flags
	cmdArgs := []string{
		"-p",                              // prompt mode (fresh session)
		"--output-format", "stream-json",  // for observability
		"--verbose",                       // required with stream-json
		"--permission-mode", "acceptEdits", // trust mode for file writes
		"--allowedTools", "Read,Write,Bash", // tool scoping
	}

	// Create command
	cmd := exec.Command("claude", cmdArgs...)

	// Set up stdin pipe to send prompt
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("create stdin pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start claude command: %w", err)
	}

	// Send prompt via stdin
	if _, err := stdin.Write([]byte(prompt)); err != nil {
		stdin.Close()
		cmd.Wait()
		return nil, fmt.Errorf("write prompt to stdin: %w", err)
	}
	stdin.Close()

	// Wait for completion
	err = cmd.Wait()
	duration := time.Since(startTime)

	// Prepare result
	result := &BackendResult{
		ExitCode: 0,
		Duration: duration,
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = 1
		}

		failureKind := "backend_execution_failed"
		result.FailureKind = &failureKind
		diagnostic := err.Error()
		result.Diagnostic = &diagnostic

		return result, fmt.Errorf("claude execution failed: %w", err)
	}

	return result, nil
}