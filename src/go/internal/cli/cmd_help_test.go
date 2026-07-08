package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// cnos#612 AC4: `cn help` must list every command by the form that
// actually invokes it. Kernel commands whose name is shadowed by a noun
// group (e.g. "cell-finalize", sibling to "cell-return"/"cell-resume")
// are only dispatchable via the space form ("cell finalize") — the
// hyphenated flat token resolves to the group listing instead (see
// TestResolveCommandFlatHyphenatedRejected in dispatch_test.go). Before
// this fix, help printed the hyphenated internal key verbatim.
func TestHelpCmd_DisplayNamesMatchInvocation(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&stubCmd{spec: CommandSpec{Name: "help", Summary: "Show available commands", Source: SourceKernel}})
	reg.Register(&stubCmd{spec: CommandSpec{Name: "cell-return", Summary: "Deliver operator verdict", Source: SourceKernel}})
	reg.Register(&stubCmd{spec: CommandSpec{Name: "cell-finalize", Summary: "Finalize a cell", Source: SourceKernel}})
	reg.Register(&stubCmd{spec: CommandSpec{Name: "doctor", Summary: "Validate hub health", Source: SourceKernel}})

	h := &HelpCmd{Registry: reg}
	var buf bytes.Buffer
	if err := h.Run(context.Background(), Invocation{Stdout: &buf, Stderr: &buf}); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"cell return", "cell finalize", "doctor"} {
		if !strings.Contains(out, want) {
			t.Errorf("help output missing invocable form %q:\n%s", want, out)
		}
	}
	for _, unwanted := range []string{"cell-return", "cell-finalize"} {
		if strings.Contains(out, unwanted) {
			t.Errorf("help output leaked non-invocable hyphenated key %q:\n%s", unwanted, out)
		}
	}
}
