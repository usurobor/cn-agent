//go:build !unix

package activate

import "fmt"

// CheckSpawnBinary on non-Unix platforms reports the unsupported-platform
// state immediately. The no-flag default path (Run → stdout) remains
// functional on every platform — only --claude / --codex spawn requires
// the syscall.Exec primitive that lives in spawn.go's unix-tagged twin.
func CheckSpawnBinary(_, flag string) error {
	return fmt.Errorf("%s spawn is not supported on this platform — use the default stdout form (cn activate HUB_DIR | <cli>) instead", flag)
}

// Spawn on non-Unix platforms returns the same unsupported-platform error.
// Reachable only if a caller bypasses CheckSpawnBinary; defensive in case.
func Spawn(binary, _ string) error {
	return fmt.Errorf("%s spawn is not supported on this platform", binary)
}
