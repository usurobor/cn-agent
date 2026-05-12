package dispatch

import (
	"fmt"
	"time"
)

// PrintBackend implements a backend that prints the prompt without execution
type PrintBackend struct{}

func (b *PrintBackend) Name() string {
	return "print"
}

func (b *PrintBackend) IsAvailable() error {
	return nil // print is always available
}

func (b *PrintBackend) Dispatch(args *Args, prompt string) (*BackendResult, error) {
	startTime := time.Now()

	// Print the constructed prompt to stdout
	output := fmt.Sprintf("# Dispatch Prompt for %s (issue %d, branch %s)\n\n%s\n",
		args.Role, args.Issue, args.Branch, prompt)

	return &BackendResult{
		ExitCode: 0,
		Duration: time.Since(startTime),
		Stdout:   output,
	}, nil
}