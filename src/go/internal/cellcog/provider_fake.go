package cellcog

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// FakeCoder makes one deterministic, real change so CI exercises the whole
// substrate — worktree, edit, measured diff — offline and without a model. A
// run behind it is `mechanical`: nothing was rented, and the closure must not
// imply otherwise.
type FakeCoder struct{}

func (FakeCoder) Name() string { return "fake" }

func (FakeCoder) Work(_ context.Context, dir, prompt string) error {
	note := "fake coder: deterministic change, no cognition was rented\n" +
		fmt.Sprintf("prompt bytes: %d\n", len(prompt))
	return os.WriteFile(filepath.Join(dir, "CELL-FAKE-CHANGE.txt"), []byte(note), 0o600)
}

// FakeAnswerer is the deterministic answering counterpart. It exists so CI can
// exercise construction, prompt rendering, schema handling and decode offline
// — NOT so a review can be simulated.
//
// It therefore never passes. A fake that returned `pass: true` would be a
// false completion wearing a reviewer's clothes, and a fake that guessed from
// the prompt would be inventing judgement it does not have. Refusing is the
// only honest deterministic verdict, and it is the same stance
// `cdd.mechanical-unmet` takes for the same reason.
type FakeAnswerer struct{}

func (FakeAnswerer) Name() string { return "fake" }

func (FakeAnswerer) Answer(_ context.Context, prompt string, schema json.RawMessage) (json.RawMessage, error) {
	if len(schema) == 0 {
		return nil, fmt.Errorf("fake: Answer needs an answer schema")
	}
	return json.Marshal(struct {
		Pass  bool   `json:"pass"`
		Notes string `json:"notes"`
	}{
		Pass: false,
		Notes: "fake answerer: no cognition was rented, so no judgement was formed " +
			fmt.Sprintf("(prompt bytes: %d)", len(prompt)),
	})
}
