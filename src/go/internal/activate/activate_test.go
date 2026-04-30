package activate

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// makeFixtureHub creates a minimal hub fixture at dir/hub-name and returns the hub path.
func makeFixtureHub(t *testing.T, name string) string {
	t.Helper()
	hub := filepath.Join(t.TempDir(), name)
	dirs := []string{
		".cn",
		"spec",
		"state",
		"threads/in",
		"threads/mail",
		"threads/reflections/daily",
		"threads/reflections/weekly",
		"threads/adhoc",
		"threads/archived",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(hub, d), 0755); err != nil {
			t.Fatalf("makeFixtureHub: mkdir %s: %v", d, err)
		}
	}
	config := `{"name":"` + name + `","version":"1.0.0","created":"2026-01-01T00:00:00Z"}`
	if err := os.WriteFile(filepath.Join(hub, ".cn", "config.json"), []byte(config), 0644); err != nil {
		t.Fatalf("makeFixtureHub: write config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(hub, "spec", "SOUL.md"), []byte("# Soul\n"), 0644); err != nil {
		t.Fatalf("makeFixtureHub: write SOUL.md: %v", err)
	}
	return hub
}

func TestRunPositive_CwdHub(t *testing.T) {
	hub := makeFixtureHub(t, "test-agent")

	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), Options{
		HubPath: hub,
		Stdout:  &stdout,
		Stderr:  &stderr,
	})
	if err != nil {
		t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
	}

	out := stdout.String()

	// AC1: prompt is printed to stdout
	if !strings.Contains(out, "You are activating a cnos hub.") {
		t.Error("prompt missing activation header")
	}
	// Hub path appears in prompt
	if !strings.Contains(out, hub) {
		t.Errorf("prompt missing hub path %q", hub)
	}
	// Hub name from config
	if !strings.Contains(out, "Hub name: test-agent") {
		t.Error("prompt missing hub name")
	}
	// Identity file listed
	if !strings.Contains(out, "spec/SOUL.md: present") {
		t.Error("prompt missing identity file")
	}
}

func TestRunPositive_ExplicitHubDir(t *testing.T) {
	hub := makeFixtureHub(t, "explicit-agent")

	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), Options{
		HubPath: hub,
		Stdout:  &stdout,
		Stderr:  &stderr,
	})
	if err != nil {
		t.Fatalf("Run with explicit HubPath: %v", err)
	}
	if !strings.Contains(stdout.String(), "You are activating a cnos hub.") {
		t.Error("prompt missing header")
	}
}

func TestRunPositive_PackagesListed(t *testing.T) {
	hub := makeFixtureHub(t, "pkg-agent")
	// Install two fake packages.
	for _, pkg := range []string{"cnos.core", "cnos.eng"} {
		if err := os.MkdirAll(filepath.Join(hub, ".cn", "vendor", "packages", pkg), 0755); err != nil {
			t.Fatal(err)
		}
	}

	var stdout, stderr bytes.Buffer
	if err := Run(context.Background(), Options{HubPath: hub, Stdout: &stdout, Stderr: &stderr}); err != nil {
		t.Fatalf("Run: %v", err)
	}
	out := stdout.String()
	if !strings.Contains(out, "cnos.core") || !strings.Contains(out, "cnos.eng") {
		t.Errorf("prompt missing package names:\n%s", out)
	}
	if !strings.Contains(out, "2 installed") {
		t.Errorf("prompt missing package count:\n%s", out)
	}
}

func TestRunPositive_SecretsExcluded(t *testing.T) {
	hub := makeFixtureHub(t, "secret-agent")
	// Write a fake secrets file — content must never appear in stdout.
	secret := "MY_TOKEN=super-secret-token-xyz-do-not-leak"
	if err := os.WriteFile(filepath.Join(hub, ".cn", "secrets.env"), []byte(secret), 0600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := Run(context.Background(), Options{HubPath: hub, Stdout: &stdout, Stderr: &stderr}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	// AC6 negative oracle: secret content must not appear in stdout.
	if strings.Contains(stdout.String(), "super-secret-token-xyz-do-not-leak") {
		t.Error("stdout contains secret content — secrets.env must not be included in prompt")
	}
	// AC6 positive: prompt notes that secrets are excluded.
	if !strings.Contains(stdout.String(), "secrets.env") {
		t.Error("prompt should mention secrets.env exclusion in Notes section")
	}
}

func TestRunPositive_StdoutOnly(t *testing.T) {
	hub := makeFixtureHub(t, "stdout-agent")

	var stdout, stderr bytes.Buffer
	if err := Run(context.Background(), Options{HubPath: hub, Stdout: &stdout, Stderr: &stderr}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	// AC8: diagnostics go to stderr, not stdout.
	stderrOut := stderr.String()
	stdoutOut := stdout.String()
	if strings.Contains(stdoutOut, "→") {
		t.Error("diagnostic arrow (→) found in stdout — diagnostics must go to stderr")
	}
	if !strings.Contains(stderrOut, "Generating activation prompt") {
		t.Error("expected diagnostic in stderr")
	}
	// stdout must contain the prompt
	if !strings.Contains(stdoutOut, "You are activating a cnos hub.") {
		t.Error("prompt missing from stdout")
	}
}

func TestRunPositive_NoModelInvocation(t *testing.T) {
	// AC7: verifiable by inspection — Run() has no subprocess calls.
	// This test confirms Run() completes without calling any external process.
	hub := makeFixtureHub(t, "nomodel-agent")

	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), Options{HubPath: hub, Stdout: &stdout, Stderr: &stderr})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	// Prompt says so explicitly.
	if !strings.Contains(stdout.String(), "No model is invoked.") {
		t.Error("prompt should state that no model is invoked")
	}
}

func TestRunNegative_EmptyHubPath(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), Options{HubPath: "", Stdout: &stdout, Stderr: &stderr})
	if err == nil {
		t.Fatal("expected error for empty HubPath")
	}
	if stdout.Len() > 0 {
		t.Errorf("stdout should be empty on failure, got: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "No hub found") {
		t.Errorf("stderr should describe the error, got: %q", stderr.String())
	}
}

func TestRunNegative_MissingPath(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), Options{
		HubPath: "/nonexistent/path/that/does/not/exist",
		Stdout:  &stdout,
		Stderr:  &stderr,
	})
	if err == nil {
		t.Fatal("expected error for nonexistent path")
	}
	if stdout.Len() > 0 {
		t.Errorf("stdout should be empty on failure, got: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "not found") {
		t.Errorf("expected 'not found' in stderr, got: %q", stderr.String())
	}
}

func TestRunNegative_PathWithoutDotCn(t *testing.T) {
	// A directory that exists but has no .cn/ — not a hub.
	dir := t.TempDir()

	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), Options{HubPath: dir, Stdout: &stdout, Stderr: &stderr})
	if err == nil {
		t.Fatal("expected error for directory without .cn/")
	}
	if stdout.Len() > 0 {
		t.Errorf("stdout should be empty on failure, got: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Not a hub") {
		t.Errorf("expected 'Not a hub' in stderr, got: %q", stderr.String())
	}
}

func TestRunPositive_ClaudeCliExampleInPrompt(t *testing.T) {
	hub := makeFixtureHub(t, "cli-example-agent")

	var stdout, stderr bytes.Buffer
	if err := Run(context.Background(), Options{HubPath: hub, Stdout: &stdout, Stderr: &stderr}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	// AC9: the example must include a query argument with -p, not bare "claude -p".
	out := stdout.String()
	if !strings.Contains(out, `claude -p "Activate this cnos hub`) {
		t.Errorf("prompt missing valid Claude CLI example with query:\n%s", out)
	}
}
