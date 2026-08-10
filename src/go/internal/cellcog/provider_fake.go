package cellcog

import "context"

// Fake answers deterministically without renting anything. It exists so the
// cognition seam — prompt render, provider call, envelope parse, seat return,
// kernel close — runs in CI on every commit, offline and without a model.
//
// A run behind Fake is therefore `mechanical`, never `cognitive`: nothing was
// rented, and the closure must not imply otherwise. The provider name is
// disclosed in the record's resolved parameters, so a reader can always tell
// which of the two happened.
type Fake struct{}

func (Fake) Name() string { return "fake" }

func (Fake) Complete(_ context.Context, _ string) (string, error) {
	return fakeAnswer, nil
}

const fakeAnswer = `{"matter":"fake provider: deterministic answer, no cognition was rented","artifacts":[]}`
