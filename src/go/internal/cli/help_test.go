package cli

import (
	"bytes"
	"strings"
	"testing"
)

// helpStubCmd is a stubCmd that also implements HelpProvider.
type helpStubCmd struct {
	stubCmd
	help string
}

func (c *helpStubCmd) Help() string { return c.help }

var _ Command = (*helpStubCmd)(nil)
var _ HelpProvider = (*helpStubCmd)(nil)

// cnos#612 AC2: a command implementing HelpProvider gets its own Help()
// text verbatim, not the generic fallback.
func TestPrintCommandHelp_UsesHelpProvider(t *testing.T) {
	reg := NewRegistry()
	cmd := &helpStubCmd{
		stubCmd: stubCmd{spec: CommandSpec{Name: "widget", Summary: "Do widget things"}},
		help:    "cn widget - detailed usage\n\nFLAGS:\n  --color\n",
	}
	reg.Register(cmd)

	var buf bytes.Buffer
	PrintCommandHelp(&buf, reg, cmd)

	if buf.String() != cmd.help {
		t.Errorf("PrintCommandHelp = %q, want the HelpProvider's own text %q", buf.String(), cmd.help)
	}
}

// cnos#612 AC2: a command with no custom Help() still prints usable
// "Usage: cn <invocation> [args...]" text with its summary, exit 0
// (checked by the caller — PrintCommandHelp never errors).
func TestPrintCommandHelp_GenericFallback(t *testing.T) {
	reg := NewRegistry()
	cmd := &stubCmd{spec: CommandSpec{Name: "doctor", Summary: "Validate hub health"}}
	reg.Register(cmd)

	var buf bytes.Buffer
	PrintCommandHelp(&buf, reg, cmd)

	out := buf.String()
	if !strings.Contains(out, "Usage: cn doctor [args...]") {
		t.Errorf("generic fallback missing usage line, got: %q", out)
	}
	if !strings.Contains(out, "Validate hub health") {
		t.Errorf("generic fallback missing summary, got: %q", out)
	}
}

// The generic fallback's invocation form must match InvocationName, so a
// hyphenated noun-group command (e.g. "issues-fsm") shows the space form
// that actually dispatches it, not the internal registry key.
func TestPrintCommandHelp_GenericFallback_UsesInvocationName(t *testing.T) {
	reg := newTestRegistry() // has kata-run / kata-list group
	cmd, ok := reg.Lookup("kata-run")
	if !ok {
		t.Fatal("expected kata-run registered in newTestRegistry")
	}

	var buf bytes.Buffer
	PrintCommandHelp(&buf, reg, cmd)

	out := buf.String()
	if !strings.Contains(out, "Usage: cn kata run [args...]") {
		t.Errorf("generic fallback should use space form for shadowed name, got: %q", out)
	}
	if strings.Contains(out, "kata-run") {
		t.Errorf("generic fallback leaked internal hyphenated key: %q", out)
	}
}
