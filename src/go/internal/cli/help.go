package cli

import (
	"fmt"
	"io"
)

// HelpProvider is implemented by commands whose --help/-h output is richer
// than the generic invocation+summary fallback (flag docs, examples, exit
// codes). Commands that don't implement it get PrintCommandHelp's generic
// fallback for free — see cnos#612 AC2.
type HelpProvider interface {
	Help() string
}

// PrintCommandHelp writes --help/-h output for cmd to w: the command's own
// Help() text when it implements HelpProvider, otherwise a generic
// "Usage: cn <invocation> [args...]" + summary line built from its Spec.
//
// Called from main.go before the hub-requirement gate, so --help works
// even on commands that need a hub to actually run (cnos#612 AC2).
func PrintCommandHelp(w io.Writer, reg *Registry, cmd Command) {
	if hp, ok := cmd.(HelpProvider); ok {
		fmt.Fprint(w, hp.Help())
		return
	}
	spec := cmd.Spec()
	fmt.Fprintf(w, "Usage: cn %s [args...]\n\n%s\n", InvocationName(reg, spec.Name), spec.Summary)
}
