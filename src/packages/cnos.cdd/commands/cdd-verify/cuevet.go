package cddverify

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
)

// cueVet invokes `cue vet -c -d <def> <schema_files...> <receipt>`.
//
// Returns (ok, diagnostics). On non-zero exit, diagnostics contains the
// stderr/stdout lines verbatim (the JSON Schema requires verbatim cue vet
// diagnostics per schemas/cdd/validation_verdict.schema.json).
//
// Adapter (eng/go §2.6, §3.7): all subprocess args are passed argv-style
// (eng/go §3.10 — never string-concatenate shell commands).
func cueVet(ctx context.Context, receiptPath string, entry DispatchEntry, repoRoot string) (ok bool, diagnostics []string) {
	args := []string{"vet", "-c", "-d", entry.Definition}
	for _, sf := range entry.SchemaFiles {
		args = append(args, filepath.Join(repoRoot, sf))
	}
	args = append(args, receiptPath)

	cmd := exec.CommandContext(ctx, "cue", args...)
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err == nil {
		return true, nil
	}
	for _, line := range strings.Split(strings.TrimRight(string(out), "\n"), "\n") {
		line = strings.TrimRight(line, " \t")
		if line != "" {
			diagnostics = append(diagnostics, line)
		}
	}
	return false, diagnostics
}

// cueAvailable returns true iff `cue version` runs cleanly.
//
// Called from Run before validate() so V's exit-2 (V itself errored) is
// emitted with a clear stderr message, not a vague cue vet failure.
func cueAvailable(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "cue", "version")
	return cmd.Run() == nil
}
