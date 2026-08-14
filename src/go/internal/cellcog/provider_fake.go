package cellcog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// FakeCoder makes one deterministic, real change so CI exercises the whole
// substrate — worktree, edit, measured diff — offline and without a model. The
// run is `mechanical`: nothing was rented and no closure may imply otherwise.
type FakeCoder struct{}

func (FakeCoder) Name() string { return "fake" }

func (FakeCoder) Work(_ context.Context, dir, prompt string) error {
	note := "fake coder: deterministic change, no cognition was rented\n" +
		fmt.Sprintf("prompt bytes: %d\n", len(prompt))
	return os.WriteFile(filepath.Join(dir, "CELL-FAKE-CHANGE.txt"), []byte(note), 0o600)
}
