package dispatch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// StubBackend implements a test backend that records prompts without execution
type StubBackend struct{}

func (b *StubBackend) Name() string {
	return "stub"
}

func (b *StubBackend) IsAvailable() error {
	return nil // stub is always available
}

func (b *StubBackend) Dispatch(args *Args, prompt string) (*BackendResult, error) {
	startTime := time.Now()

	// Record the prompt payload for test verification
	if err := b.recordPrompt(args, prompt); err != nil {
		failureKind := "stub_record_failed"
		diagnostic := err.Error()
		return &BackendResult{
			ExitCode:    1,
			Duration:    time.Since(startTime),
			FailureKind: &failureKind,
			Diagnostic:  &diagnostic,
		}, err
	}

	// Simulate successful execution
	return &BackendResult{
		ExitCode: 0,
		Duration: time.Since(startTime),
		Stdout:   fmt.Sprintf("✓ Stub backend recorded %s dispatch for issue %d\n", args.Role, args.Issue),
	}, nil
}

// recordPrompt saves the prompt to .cdd/dispatch-log/ for test verification
func (b *StubBackend) recordPrompt(args *Args, prompt string) error {
	// Create log directory
	logDir := ".cdd/dispatch-log"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("create log directory: %w", err)
	}

	// Create log entry
	logEntry := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"role":      args.Role,
		"branch":    args.Branch,
		"issue":     args.Issue,
		"project":   args.Project,
		"backend":   "stub",
		"prompt":    prompt,
	}

	// Write to log file
	filename := fmt.Sprintf("stub-%s-%d-%d.json", args.Role, args.Issue, time.Now().Unix())
	logFile := filepath.Join(logDir, filename)

	logJSON, err := json.MarshalIndent(logEntry, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal log entry: %w", err)
	}

	if err := os.WriteFile(logFile, logJSON, 0644); err != nil {
		return fmt.Errorf("write log file: %w", err)
	}

	return nil
}