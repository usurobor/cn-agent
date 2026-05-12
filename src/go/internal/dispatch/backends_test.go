package dispatch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestStubBackend(t *testing.T) {
	// Setup test directory
	testDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	os.Chdir(testDir)

	backend := &StubBackend{}

	// Test availability
	if err := backend.IsAvailable(); err != nil {
		t.Errorf("StubBackend should always be available, got: %v", err)
	}

	// Test dispatch
	args := &Args{
		Role:    "alpha",
		Branch:  "cycle/295",
		Issue:   295,
		Project: "cnos",
	}
	prompt := "Test prompt for α dispatch"

	result, err := backend.Dispatch(args, prompt)
	if err != nil {
		t.Errorf("StubBackend.Dispatch failed: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	if result.Duration == 0 {
		t.Errorf("Expected non-zero duration")
	}

	if !strings.Contains(result.Stdout, "Stub backend recorded") {
		t.Errorf("Expected success message in stdout")
	}

	// Verify log file was created
	logFiles, err := filepath.Glob(".cdd/dispatch-log/stub-*.json")
	if err != nil || len(logFiles) == 0 {
		t.Errorf("Expected stub log file to be created")
	}
}

func TestPrintBackend(t *testing.T) {
	backend := &PrintBackend{}

	// Test availability
	if err := backend.IsAvailable(); err != nil {
		t.Errorf("PrintBackend should always be available, got: %v", err)
	}

	// Test dispatch
	args := &Args{
		Role:    "beta",
		Branch:  "cycle/300",
		Issue:   300,
		Project: "cnos",
	}
	prompt := "Test prompt for β dispatch"

	result, err := backend.Dispatch(args, prompt)
	if err != nil {
		t.Errorf("PrintBackend.Dispatch failed: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}

	if result.Duration == 0 {
		t.Errorf("Expected non-zero duration")
	}

	expectedContent := []string{
		"# Dispatch Prompt for beta",
		"issue 300",
		"branch cycle/300",
		"Test prompt for β dispatch",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(result.Stdout, expected) {
			t.Errorf("Expected %q in output, got: %s", expected, result.Stdout)
		}
	}
}

func TestClaudeBackendAvailability(t *testing.T) {
	backend := &ClaudeBackend{}

	// This will likely fail since claude binary won't be available in test environment
	err := backend.IsAvailable()
	if err == nil {
		t.Skip("Claude binary is available, skipping negative test")
	}

	if !strings.Contains(err.Error(), "claude binary not found") {
		t.Errorf("Expected specific error message about claude binary, got: %v", err)
	}
}

func TestBackendFactory(t *testing.T) {
	dispatcher := &Dispatcher{}

	// Test valid backends
	validBackends := []string{"claude", "stub", "print"}
	for _, name := range validBackends {
		backend, err := dispatcher.getBackend(name)
		if name == "claude" && err != nil {
			// Claude backend may not be available in test environment
			continue
		}

		if err != nil {
			t.Errorf("Expected %s backend to be available, got: %v", name, err)
			continue
		}

		if backend.Name() != name {
			t.Errorf("Expected backend name %s, got %s", name, backend.Name())
		}
	}

	// Test invalid backend
	_, err := dispatcher.getBackend("invalid")
	if err == nil {
		t.Errorf("Expected error for invalid backend")
	}
}