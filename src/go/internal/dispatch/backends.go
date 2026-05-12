package dispatch

import (
	"fmt"
)

// getBackend returns the appropriate backend implementation
func (d *Dispatcher) getBackend(name string) (Backend, error) {
	switch name {
	case "claude":
		backend := &ClaudeBackend{}
		if err := backend.IsAvailable(); err != nil {
			return nil, err
		}
		return backend, nil
	case "stub":
		return &StubBackend{}, nil
	case "print":
		return &PrintBackend{}, nil
	default:
		return nil, fmt.Errorf("unknown backend %q", name)
	}
}