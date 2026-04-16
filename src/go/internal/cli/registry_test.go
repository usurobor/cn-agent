package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// stubCmd is a minimal Command implementation for testing.
type stubCmd struct {
	spec CommandSpec
	ran  bool
}

func (c *stubCmd) Spec() CommandSpec          { return c.spec }
func (c *stubCmd) Run(_ context.Context, _ Invocation) error {
	c.ran = true
	return nil
}

func TestRegistryLookupHit(t *testing.T) {
	reg := NewRegistry()
	cmd := &stubCmd{spec: CommandSpec{Name: "test", Summary: "a test"}}
	reg.Register(cmd)

	got, ok := reg.Lookup("test")
	if !ok {
		t.Fatal("expected hit")
	}
	if got.Spec().Name != "test" {
		t.Errorf("name = %q, want %q", got.Spec().Name, "test")
	}
}

func TestRegistryLookupMiss(t *testing.T) {
	reg := NewRegistry()
	_, ok := reg.Lookup("nonexistent")
	if ok {
		t.Fatal("expected miss")
	}
}

func TestRegistryAllPreservesOrder(t *testing.T) {
	reg := NewRegistry()
	reg.Register(&stubCmd{spec: CommandSpec{Name: "help"}})
	reg.Register(&stubCmd{spec: CommandSpec{Name: "deps"}})
	reg.Register(&stubCmd{spec: CommandSpec{Name: "doctor"}})

	cmds := reg.All()
	if len(cmds) != 3 {
		t.Fatalf("len = %d, want 3", len(cmds))
	}
	want := []string{"help", "deps", "doctor"}
	for i, cmd := range cmds {
		if cmd.Spec().Name != want[i] {
			t.Errorf("cmds[%d].Name = %q, want %q", i, cmd.Spec().Name, want[i])
		}
	}
}

func TestRegistryTierPrecedence(t *testing.T) {
	reg := NewRegistry()
	// Register a package-tier command first.
	reg.Register(&stubCmd{spec: CommandSpec{Name: "deps", Tier: TierPackage, Summary: "from package"}})
	// Then a kernel-tier command with the same name — should win.
	reg.Register(&stubCmd{spec: CommandSpec{Name: "deps", Tier: TierKernel, Summary: "from kernel"}})

	cmd, ok := reg.Lookup("deps")
	if !ok {
		t.Fatal("expected hit")
	}
	if cmd.Spec().Summary != "from kernel" {
		t.Errorf("summary = %q, want %q", cmd.Spec().Summary, "from kernel")
	}
}

func TestHelpCmdOutput(t *testing.T) {
	reg := NewRegistry()
	help := &HelpCmd{Registry: reg}
	reg.Register(help)
	reg.Register(&stubCmd{spec: CommandSpec{Name: "deps", Summary: "Manage dependencies", Source: SourceKernel, NeedsHub: true}})

	var stdout bytes.Buffer
	inv := Invocation{
		HubPath: "/some/hub",
		Stdout:  &stdout,
		Stderr:  &bytes.Buffer{},
	}
	if err := help.Run(context.Background(), inv); err != nil {
		t.Fatalf("help.Run: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "cn — cnos kernel") {
		t.Error("expected header in help output")
	}
	if !strings.Contains(out, "Kernel commands:") {
		t.Error("expected 'Kernel commands:' section header")
	}
	if !strings.Contains(out, "deps") {
		t.Error("expected 'deps' in help output")
	}
	if !strings.Contains(out, "Manage dependencies") {
		t.Error("expected deps summary in help output")
	}
}

func TestHelpCmdNoHub(t *testing.T) {
	reg := NewRegistry()
	help := &HelpCmd{Registry: reg}
	reg.Register(help)
	reg.Register(&stubCmd{spec: CommandSpec{Name: "deps", Summary: "Manage dependencies", Source: SourceKernel, NeedsHub: true}})
	reg.Register(&stubCmd{spec: CommandSpec{Name: "deploy", Summary: "Deploy to prod", Source: SourceRepoLocal, Tier: TierRepoLocal, NeedsHub: true}})

	var stdout bytes.Buffer
	inv := Invocation{
		HubPath: "", // no hub
		Stdout:  &stdout,
		Stderr:  &bytes.Buffer{},
	}
	help.Run(context.Background(), inv)

	out := stdout.String()
	// Kernel commands are always listed so users can see the full
	// kernel surface without needing to create a hub first.
	if !strings.Contains(out, "deps") {
		t.Error("kernel 'deps' should be listed even when no hub")
	}
	if !strings.Contains(out, "(requires hub)") {
		t.Error("expected '(requires hub)' annotation on hub-needing kernel commands")
	}
	// Non-kernel commands are still filtered out when no hub exists.
	if strings.Contains(out, "deploy") {
		t.Error("repo-local 'deploy' should be hidden when no hub")
	}
	if !strings.Contains(out, "No hub found") {
		t.Error("expected no-hub warning")
	}
}

func TestHelpCmdGroupListing(t *testing.T) {
	// `cn help kata` — HelpCmd must dispatch to the group listing
	// instead of the full kernel listing.
	reg := NewRegistry()
	help := &HelpCmd{Registry: reg}
	reg.Register(help)
	reg.Register(&stubCmd{spec: CommandSpec{Name: "kata-run", Summary: "Run a kata", Source: SourcePackage, Tier: TierPackage}})
	reg.Register(&stubCmd{spec: CommandSpec{Name: "kata-list", Summary: "List katas", Source: SourcePackage, Tier: TierPackage}})

	var stdout bytes.Buffer
	inv := Invocation{
		HubPath: "/some/hub",
		Args:    []string{"kata"},
		Stdout:  &stdout,
		Stderr:  &bytes.Buffer{},
	}
	if err := help.Run(context.Background(), inv); err != nil {
		t.Fatalf("help.Run: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "kata run") {
		t.Errorf("expected 'kata run' in group listing, got:\n%s", out)
	}
	if !strings.Contains(out, "kata list") {
		t.Errorf("expected 'kata list' in group listing, got:\n%s", out)
	}
	if strings.Contains(out, "Kernel commands:") {
		t.Errorf("group listing should not include the full kernel listing:\n%s", out)
	}
}

func TestHelpCmdUnknownGroupFallsThrough(t *testing.T) {
	// `cn help bogus` — unknown noun falls through to the full listing
	// rather than erroring.
	reg := NewRegistry()
	help := &HelpCmd{Registry: reg}
	reg.Register(help)
	reg.Register(&stubCmd{spec: CommandSpec{Name: "deps", Summary: "Manage package dependencies", Source: SourceKernel}})

	var stdout bytes.Buffer
	inv := Invocation{
		HubPath: "/some/hub",
		Args:    []string{"bogus"},
		Stdout:  &stdout,
		Stderr:  &bytes.Buffer{},
	}
	if err := help.Run(context.Background(), inv); err != nil {
		t.Fatalf("help.Run: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "Kernel commands:") {
		t.Errorf("expected full listing fallback, got:\n%s", out)
	}
	if !strings.Contains(out, "deps") {
		t.Errorf("expected 'deps' in fallback full listing, got:\n%s", out)
	}
}

func TestHelpCmdTieredOutput(t *testing.T) {
	reg := NewRegistry()
	help := &HelpCmd{Registry: reg}
	reg.Register(help)
	reg.Register(&stubCmd{spec: CommandSpec{Name: "deps", Summary: "Manage dependencies", Source: SourceKernel, NeedsHub: true}})
	reg.Register(&stubCmd{spec: CommandSpec{Name: "deploy", Summary: "Deploy to prod", Source: SourceRepoLocal, Tier: TierRepoLocal, NeedsHub: true}})
	reg.Register(&stubCmd{spec: CommandSpec{Name: "daily", Summary: "Daily reflection", Source: SourcePackage, Tier: TierPackage, Package: "cnos.core", NeedsHub: true}})

	var stdout bytes.Buffer
	inv := Invocation{
		HubPath: "/some/hub",
		Stdout:  &stdout,
		Stderr:  &bytes.Buffer{},
	}
	help.Run(context.Background(), inv)

	out := stdout.String()
	if !strings.Contains(out, "Kernel commands:") {
		t.Error("expected 'Kernel commands:' section")
	}
	if !strings.Contains(out, "Repo-local commands:") {
		t.Error("expected 'Repo-local commands:' section")
	}
	if !strings.Contains(out, "Package commands:") {
		t.Error("expected 'Package commands:' section")
	}
	if !strings.Contains(out, "[cnos.core]") {
		t.Error("expected package attribution '[cnos.core]'")
	}
}
