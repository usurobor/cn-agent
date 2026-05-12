package dispatch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// IsValidBackend checks if the backend name is supported
func IsValidBackend(backend string) bool {
	switch backend {
	case "claude", "stub", "print":
		return true
	default:
		return false
	}
}

// DetectProject extracts the project name from the hub path.
// Uses the directory name containing .cn/ as the project name.
func DetectProject(hubPath string) string {
	if hubPath == "" {
		return "cnos" // fallback
	}
	return filepath.Base(hubPath)
}

// ResolveBackend selects the backend: explicit flag > env > default "claude".
func ResolveBackend(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if env := os.Getenv("CN_DISPATCH_BACKEND"); env != "" {
		return env
	}
	return "claude"
}

// RecordAttempt writes the attempt descriptor to the cycle directory.
func RecordAttempt(hubPath string, issue int, result *Result) error {
	if issue == 0 {
		return nil // no issue number provided
	}

	attemptDir := filepath.Join(hubPath, ".cdd", "unreleased", fmt.Sprintf("%d", issue), "dispatch")
	if err := os.MkdirAll(attemptDir, 0755); err != nil {
		return fmt.Errorf("create dispatch directory: %w", err)
	}

	attemptFile := filepath.Join(attemptDir, result.AttemptID+".json")
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}

	if err := os.WriteFile(attemptFile, resultJSON, 0644); err != nil {
		return fmt.Errorf("write attempt file: %w", err)
	}

	return nil
}
