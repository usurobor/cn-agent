//go:build unix

package activate

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// execFunc matches the signature of syscall.Exec. The dry-run / test hook
// substitutes a recording function so spawn_test.go can verify the planned
// argv without actually replacing the current process.
type execFunc func(argv0 string, argv []string, envv []string) error

// lookPathFunc matches the signature of exec.LookPath. Tests substitute a
// fake lookup to drive the missing-binary path without touching $PATH.
type lookPathFunc func(name string) (string, error)

// defaultExecFn is the production exec primitive — syscall.Exec replaces the
// cn process image with the target binary, preserving the operator's TTY so
// the resulting claude / codex session lands in an interactive REPL.
var defaultExecFn execFunc = syscall.Exec

// defaultLookPath is the production binary-lookup primitive.
var defaultLookPath lookPathFunc = exec.LookPath

// CheckSpawnBinary verifies that the requested binary is on $PATH. It is
// called before any rendering so an operator who set --claude / --codex
// without the matching CLI installed sees a clear error early. flag is the
// user-visible flag name (e.g. "--claude") so the diagnostic names both the
// missing binary and the flag that requested it.
func CheckSpawnBinary(binary, flag string) error {
	return checkSpawnBinary(binary, flag, defaultLookPath)
}

func checkSpawnBinary(binary, flag string, lookPath lookPathFunc) error {
	if _, err := lookPath(binary); err != nil {
		return fmt.Errorf(
			"%s (requested by %s) not found in PATH — install it or ensure $PATH includes its directory",
			binary, flag,
		)
	}
	return nil
}

// Spawn replaces the cn process image with `<binary> <prompt>` using
// syscall.Exec, preserving the parent terminal's TTY / stdin / stdout / stderr
// so the resulting session is interactive. On success it does not return;
// any returned error is a setup failure (binary missing from PATH, exec
// syscall failure). Callers MUST have run CheckSpawnBinary before any
// rendering, so that the pre-render error path stays clean.
func Spawn(binary, prompt string) error {
	return spawnWith(binary, prompt, defaultLookPath, defaultExecFn)
}

func spawnWith(binary, prompt string, lookPath lookPathFunc, execFn execFunc) error {
	path, err := lookPath(binary)
	if err != nil {
		return fmt.Errorf("%s not found in PATH: %w", binary, err)
	}
	// argv[0] is the binary name as the program sees it (Unix exec
	// convention); the actual file is identified by the resolved path
	// passed to syscall.Exec as argv0. The single remaining argv element
	// is the rendered prompt — bare positional, no `-p`, no subcommand,
	// which is what drops claude / codex into interactive mode.
	argv := []string{binary, prompt}
	return execFn(path, argv, os.Environ())
}
