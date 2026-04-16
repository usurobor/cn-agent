package cli

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

// newTestRegistry builds a registry populated with a representative mix
// of flat kernel commands and a "kata" noun group.
func newTestRegistry() *Registry {
	reg := NewRegistry()
	reg.Register(&stubCmd{spec: CommandSpec{Name: "help", Summary: "Show available commands", Source: SourceKernel}})
	reg.Register(&stubCmd{spec: CommandSpec{Name: "doctor", Summary: "Run diagnostic checks", Source: SourceKernel}})
	reg.Register(&stubCmd{spec: CommandSpec{Name: "deps", Summary: "Manage package dependencies", Source: SourceKernel}})
	reg.Register(&stubCmd{spec: CommandSpec{Name: "kata-run", Summary: "Run a kata", Source: SourcePackage, Tier: TierPackage, Package: "cnos.kata"}})
	reg.Register(&stubCmd{spec: CommandSpec{Name: "kata-list", Summary: "List katas", Source: SourcePackage, Tier: TierPackage, Package: "cnos.kata"}})
	reg.Register(&stubCmd{spec: CommandSpec{Name: "kata-judge", Summary: "Judge a kata bundle", Source: SourcePackage, Tier: TierPackage, Package: "cnos.kata"}})
	return reg
}

func TestResolveCommandNounVerb(t *testing.T) {
	reg := newTestRegistry()

	res := ResolveCommand(reg, []string{"kata", "run", "--verbose", "M0-gap"})

	if res.Command == nil {
		t.Fatal("expected kata-run to resolve, got nil")
	}
	if got := res.Command.Spec().Name; got != "kata-run" {
		t.Errorf("Command.Name = %q, want %q", got, "kata-run")
	}
	if !reflect.DeepEqual(res.Remaining, []string{"--verbose", "M0-gap"}) {
		t.Errorf("Remaining = %v, want [--verbose M0-gap]", res.Remaining)
	}
	if res.Group != "" {
		t.Errorf("Group = %q, want empty", res.Group)
	}
}

func TestResolveCommandFlatHyphenatedRejected(t *testing.T) {
	// Hyphenated forms like "kata-run" are no longer resolved.
	// They fall through to the group listing for "kata".
	reg := newTestRegistry()

	res := ResolveCommand(reg, []string{"kata-run", "--verbose"})

	if res.Command != nil {
		t.Errorf("Command = %v, want nil (hyphenated form should be rejected)", res.Command)
	}
	if res.Group != "kata" {
		t.Errorf("Group = %q, want %q", res.Group, "kata")
	}
}

func TestResolveCommandFlatKernelCommand(t *testing.T) {
	reg := newTestRegistry()

	res := ResolveCommand(reg, []string{"doctor"})

	if res.Command == nil {
		t.Fatal("expected doctor to resolve")
	}
	if got := res.Command.Spec().Name; got != "doctor" {
		t.Errorf("Command.Name = %q, want %q", got, "doctor")
	}
	if len(res.Remaining) != 0 {
		t.Errorf("Remaining = %v, want empty", res.Remaining)
	}
}

func TestResolveCommandDepsSubcommandPreserved(t *testing.T) {
	// `cn deps lock` should fall back to the flat `deps` command with
	// ["lock"] as args, because `deps-lock` is not a separate registry
	// entry. This preserves DepsCmd's internal subcommand dispatch.
	reg := newTestRegistry()

	res := ResolveCommand(reg, []string{"deps", "lock"})

	if res.Command == nil {
		t.Fatal("expected deps to resolve")
	}
	if got := res.Command.Spec().Name; got != "deps" {
		t.Errorf("Command.Name = %q, want %q", got, "deps")
	}
	if !reflect.DeepEqual(res.Remaining, []string{"lock"}) {
		t.Errorf("Remaining = %v, want [lock]", res.Remaining)
	}
}

func TestResolveCommandGroupNoVerb(t *testing.T) {
	reg := newTestRegistry()

	res := ResolveCommand(reg, []string{"kata"})

	if res.Command != nil {
		t.Errorf("Command = %v, want nil (group case)", res.Command)
	}
	if res.Group != "kata" {
		t.Errorf("Group = %q, want %q", res.Group, "kata")
	}
	if len(res.Remaining) != 0 {
		t.Errorf("Remaining = %v, want empty", res.Remaining)
	}
}

func TestResolveCommandGroupUnknownVerb(t *testing.T) {
	// `cn kata bogus` — not a known verb; kata-* exists.
	// The resolver reports the group and the caller decides what to do
	// with the unknown verb (we keep it in Remaining for the caller).
	reg := newTestRegistry()

	res := ResolveCommand(reg, []string{"kata", "bogus"})

	if res.Command != nil {
		t.Errorf("Command = %v, want nil", res.Command)
	}
	if res.Group != "kata" {
		t.Errorf("Group = %q, want %q", res.Group, "kata")
	}
	if !reflect.DeepEqual(res.Remaining, []string{"bogus"}) {
		t.Errorf("Remaining = %v, want [bogus]", res.Remaining)
	}
}

func TestResolveCommandGroupHelpFlag(t *testing.T) {
	// `cn kata --help` — --help is not a verb, falls through to group.
	reg := newTestRegistry()

	res := ResolveCommand(reg, []string{"kata", "--help"})

	if res.Group != "kata" {
		t.Errorf("Group = %q, want %q", res.Group, "kata")
	}
	if !reflect.DeepEqual(res.Remaining, []string{"--help"}) {
		t.Errorf("Remaining = %v, want [--help]", res.Remaining)
	}
}

func TestResolveCommandUnknown(t *testing.T) {
	reg := newTestRegistry()

	res := ResolveCommand(reg, []string{"bogus"})

	if res.Command != nil {
		t.Errorf("Command = %v, want nil", res.Command)
	}
	if res.Group != "" {
		t.Errorf("Group = %q, want empty", res.Group)
	}
}

func TestResolveCommandEmpty(t *testing.T) {
	reg := newTestRegistry()

	res := ResolveCommand(reg, nil)

	if res.Command != nil || res.Group != "" {
		t.Errorf("empty args should give empty Resolution, got %+v", res)
	}
}

func TestGroupMembersOrder(t *testing.T) {
	reg := newTestRegistry()

	members := GroupMembers(reg, "kata")

	if len(members) != 3 {
		t.Fatalf("len(members) = %d, want 3", len(members))
	}
	want := []string{"kata-run", "kata-list", "kata-judge"}
	for i, m := range members {
		if m.Spec().Name != want[i] {
			t.Errorf("members[%d].Name = %q, want %q", i, m.Spec().Name, want[i])
		}
	}
}

func TestGroupMembersEmpty(t *testing.T) {
	reg := newTestRegistry()

	if got := GroupMembers(reg, "bogus"); len(got) != 0 {
		t.Errorf("GroupMembers(bogus) = %v, want empty", got)
	}
}

func TestGroupMembersExactNameNotMember(t *testing.T) {
	// A flat command named "deps" should not be returned as a member of
	// the "deps" group — only "deps-<verb>" commands are.
	reg := newTestRegistry()

	if got := GroupMembers(reg, "deps"); len(got) != 0 {
		t.Errorf("GroupMembers(deps) = %v, want empty (deps is flat, not a group)", got)
	}
}

func TestPrintGroup(t *testing.T) {
	reg := newTestRegistry()

	var buf bytes.Buffer
	ok := PrintGroup(&buf, reg, "kata")

	if !ok {
		t.Fatal("PrintGroup returned false for a known group")
	}
	out := buf.String()
	// User-facing form uses "kata run", not "kata-run".
	for _, want := range []string{"kata run", "kata list", "kata judge", "Run a kata", "List katas", "Judge a kata bundle"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\noutput:\n%s", want, out)
		}
	}
	// Internal hyphenated keys must not leak to the user-facing listing.
	if strings.Contains(out, "kata-run") {
		t.Errorf("output leaks internal hyphenated key kata-run:\n%s", out)
	}
}

func TestPrintGroupUnknownPrefix(t *testing.T) {
	reg := newTestRegistry()

	var buf bytes.Buffer
	ok := PrintGroup(&buf, reg, "bogus")

	if ok {
		t.Error("PrintGroup should return false for unknown prefix")
	}
	if buf.Len() != 0 {
		t.Errorf("PrintGroup wrote %q for unknown prefix, want nothing", buf.String())
	}
}
