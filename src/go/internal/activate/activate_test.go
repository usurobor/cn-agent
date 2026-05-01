package activate

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- Fixture builders ---

// makeSigmaFixture creates a hub mirroring the cn-sigma pattern:
// spec/PERSONA.md + spec/OPERATOR.md + vendored cnos.core with doctrine/KERNEL.md.
func makeSigmaFixture(t *testing.T) string {
	t.Helper()
	hub := t.TempDir()
	mustMkdir(t, hub, ".cn")
	writeConfigJSON(t, hub, "sigma-hub", "1.0.0")
	mustMkdir(t, hub, "spec")
	mustWrite(t, filepath.Join(hub, "spec", "PERSONA.md"), "# Persona\n")
	mustWrite(t, filepath.Join(hub, "spec", "OPERATOR.md"), "# Operator\n")
	// Vendored kernel with package manifest
	vendorCore := filepath.Join(hub, ".cn", "vendor", "packages", "cnos.core")
	mustMkdir(t, hub, ".cn/vendor/packages/cnos.core/doctrine")
	mustWrite(t, filepath.Join(vendorCore, "doctrine", "KERNEL.md"), "# Kernel\n")
	mustWrite(t, filepath.Join(vendorCore, "cn.package.json"),
		`{"schema":"cn.package.v1","name":"cnos.core","version":"3.71.0"}`)
	// Write deps.json declaring cnos.core
	mustWrite(t, filepath.Join(hub, ".cn", "deps.json"),
		`{"schema":"cn.deps.v1","packages":[{"name":"cnos.core","version":"3.71.0"}]}`)
	return hub
}

// makeInitOnlyFixture creates a hub mirroring cn init output:
// spec/SOUL.md only, no .cn/deps.json, no vendor, no PERSONA, no OPERATOR.
func makeInitOnlyFixture(t *testing.T) string {
	t.Helper()
	hub := t.TempDir()
	mustMkdir(t, hub, ".cn")
	writeConfigJSON(t, hub, "init-hub", "1.0.0")
	mustMkdir(t, hub, "spec")
	mustWrite(t, filepath.Join(hub, "spec", "SOUL.md"), "# Soul\n")
	return hub
}

// makeInitPlusSetupFixture creates a hub mirroring cn init + cn setup:
// spec/SOUL.md + .cn/deps.json declaring cnos.core, no vendor, no PERSONA, no OPERATOR.
func makeInitPlusSetupFixture(t *testing.T) string {
	t.Helper()
	hub := t.TempDir()
	mustMkdir(t, hub, ".cn")
	writeConfigJSON(t, hub, "setup-hub", "1.0.0")
	mustMkdir(t, hub, "spec")
	mustWrite(t, filepath.Join(hub, "spec", "SOUL.md"), "# Soul\n")
	mustWrite(t, filepath.Join(hub, ".cn", "deps.json"),
		`{"schema":"cn.deps.v1","packages":[{"name":"cnos.core","version":"3.71.0"}]}`)
	return hub
}

// makeMinimalHub creates the minimal hub fixture used by basic tests.
func makeMinimalHub(t *testing.T, name string) string {
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
			t.Fatalf("makeMinimalHub: mkdir %s: %v", d, err)
		}
	}
	writeConfigJSON(t, hub, name, "1.0.0")
	mustWrite(t, filepath.Join(hub, "spec", "PERSONA.md"), "# Persona\n")
	return hub
}

func mustMkdir(t *testing.T, hub, rel string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(hub, rel), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", rel, err)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func writeConfigJSON(t *testing.T, hub, name, version string) {
	t.Helper()
	cfg := `{"name":"` + name + `","version":"` + version + `","created":"2026-01-01T00:00:00Z"}`
	mustWrite(t, filepath.Join(hub, ".cn", "config.json"), cfg)
}

func run(t *testing.T, hub string) (string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	err := Run(context.Background(), Options{HubPath: hub, Stdout: &stdout, Stderr: &stderr})
	if err != nil {
		t.Fatalf("Run: %v\nstderr: %s", err, stderr.String())
	}
	return stdout.String(), stderr.String()
}

// --- AC11: Sigma-shape fixture ---

func TestSigmaShapeActivatesCorrectly(t *testing.T) {
	hub := makeSigmaFixture(t)
	out, _ := run(t, hub)

	// AC3: three layered sections in correct order
	kernelIdx := strings.Index(out, "## Kernel")
	personaIdx := strings.Index(out, "## Persona")
	operatorIdx := strings.Index(out, "## Operator")
	if kernelIdx < 0 || personaIdx < 0 || operatorIdx < 0 {
		t.Fatalf("missing layered sections in output:\n%s", out)
	}
	if !(kernelIdx < personaIdx && personaIdx < operatorIdx) {
		t.Errorf("sections out of order: Kernel=%d Persona=%d Operator=%d", kernelIdx, personaIdx, operatorIdx)
	}

	// AC11 positive: kernel vendored path + persona + operator populated
	if !strings.Contains(out, "vendored at") {
		t.Errorf("## Kernel section must show vendored path:\n%s", out)
	}
	if !strings.Contains(out, "doctrine/KERNEL.md") {
		t.Errorf("## Kernel must reference doctrine/KERNEL.md:\n%s", out)
	}
	if !strings.Contains(out, "@3.71.0") {
		t.Errorf("## Kernel must show version @3.71.0:\n%s", out)
	}
	if !strings.Contains(out, "spec/PERSONA.md: present") {
		t.Errorf("## Persona must show spec/PERSONA.md present:\n%s", out)
	}
	if !strings.Contains(out, "spec/OPERATOR.md: present") {
		t.Errorf("## Operator must show spec/OPERATOR.md present:\n%s", out)
	}

	// AC11 negative: must not contain old "no identity files found"
	if strings.Contains(out, "no identity files found") {
		t.Error("stdout must not contain 'no identity files found'")
	}
	if strings.Contains(out, "spec/SOUL.md") {
		t.Error("stdout must not reference spec/SOUL.md")
	}

	// AC1 negative: no ## Identity header
	if strings.Contains(out, "## Identity") {
		t.Error("stdout must not contain legacy ## Identity bucket")
	}
}

// --- AC12a: Init-only fixture ---

func TestInitOnlyFixture_NoKernelReference(t *testing.T) {
	hub := makeInitOnlyFixture(t)
	out, _ := run(t, hub)

	// AC12a: kernel = "no kernel reference"
	if !strings.Contains(out, "no kernel reference") {
		t.Errorf("init-only fixture must emit 'no kernel reference':\n%s", out)
	}
	// AC12a: must not suggest cn deps restore (no manifest)
	if strings.Contains(out, "run cn deps restore") {
		t.Errorf("init-only fixture must not suggest cn deps restore (no manifest):\n%s", out)
	}
	// AC12a: persona absent explicit
	if !strings.Contains(out, "no spec/PERSONA.md found") {
		t.Errorf("init-only must name absent PERSONA.md:\n%s", out)
	}
	// AC12a: operator absent explicit
	if !strings.Contains(out, "no spec/OPERATOR.md found") {
		t.Errorf("init-only must name absent OPERATOR.md:\n%s", out)
	}
	// AC12a: legacy spec/SOUL.md not surfaced under any layered section
	if strings.Contains(out, "spec/SOUL.md") {
		t.Errorf("init-only must not surface spec/SOUL.md as kernel/persona/operator:\n%s", out)
	}
	// AC1 negative
	if strings.Contains(out, "## Identity") {
		t.Error("must not contain legacy ## Identity bucket")
	}
}

// --- AC12b: Init+setup fixture ---

func TestInitPlusSetupFixture_RestoreGuidance(t *testing.T) {
	hub := makeInitPlusSetupFixture(t)
	out, _ := run(t, hub)

	// AC12b: kernel = manifest-only state
	if !strings.Contains(out, "dependency manifest declares cnos.core; not restored — run cn deps restore") {
		t.Errorf("init+setup must emit manifest-only kernel guidance:\n%s", out)
	}
	// AC12b: must not claim kernel is vendored
	if strings.Contains(out, "vendored at") {
		t.Errorf("init+setup must not claim kernel is vendored:\n%s", out)
	}
	// AC12b: persona absent explicit
	if !strings.Contains(out, "no spec/PERSONA.md found") {
		t.Errorf("init+setup must name absent PERSONA.md:\n%s", out)
	}
	// AC12b: operator absent explicit
	if !strings.Contains(out, "no spec/OPERATOR.md found") {
		t.Errorf("init+setup must name absent OPERATOR.md:\n%s", out)
	}
	// AC12b: legacy spec/SOUL.md not surfaced
	if strings.Contains(out, "spec/SOUL.md") {
		t.Errorf("init+setup must not surface spec/SOUL.md:\n%s", out)
	}
}

// --- AC7: Three kernel states ---

func TestKernelState_Vendored(t *testing.T) {
	hub := makeSigmaFixture(t)
	out, _ := run(t, hub)
	if !strings.Contains(out, "vendored at") || !strings.Contains(out, "doctrine/KERNEL.md") {
		t.Errorf("vendored fixture must show resolved kernel path:\n%s", out)
	}
}

func TestKernelState_ManifestOnly(t *testing.T) {
	hub := makeInitPlusSetupFixture(t)
	out, _ := run(t, hub)
	if !strings.Contains(out, "dependency manifest declares cnos.core; not restored — run cn deps restore") {
		t.Errorf("manifest-only fixture must show restore guidance:\n%s", out)
	}
	if strings.Contains(out, "vendored at") {
		t.Error("manifest-only must not show vendored path")
	}
}

func TestKernelState_None(t *testing.T) {
	hub := makeInitOnlyFixture(t)
	out, _ := run(t, hub)
	if !strings.Contains(out, "no kernel reference") {
		t.Errorf("no-manifest fixture must show 'no kernel reference':\n%s", out)
	}
}

// --- AC2: Three resolvers, no legacy paths ---

func TestResolversRejectLegacyPaths(t *testing.T) {
	hub := t.TempDir()
	mustMkdir(t, hub, ".cn")
	writeConfigJSON(t, hub, "legacy-hub", "1.0.0")
	mustMkdir(t, hub, "spec")
	// Place only legacy files — none should be surfaced under the new sections
	mustWrite(t, filepath.Join(hub, "spec", "SOUL.md"), "# Soul\n")
	mustWrite(t, filepath.Join(hub, "spec", "identity.md"), "# identity\n")
	mustWrite(t, filepath.Join(hub, "agent", "identity.md"), "# agent identity\n")
	mustWrite(t, filepath.Join(hub, "spec", "USER.md"), "# User\n")

	out, _ := run(t, hub)

	for _, legacyPath := range []string{
		"spec/SOUL.md", "spec/identity.md", "agent/identity.md", "spec/USER.md",
	} {
		if strings.Contains(out, legacyPath) {
			t.Errorf("legacy path %q must not appear under layered sections:\n%s", legacyPath, out)
		}
	}
	// All three sections should report absence
	if !strings.Contains(out, "no kernel reference") {
		t.Errorf("must report no kernel reference:\n%s", out)
	}
	if !strings.Contains(out, "no spec/PERSONA.md found") {
		t.Errorf("must report absent PERSONA.md:\n%s", out)
	}
	if !strings.Contains(out, "no spec/OPERATOR.md found") {
		t.Errorf("must report absent OPERATOR.md:\n%s", out)
	}
}

// --- AC8: Deps detection three states ---

func TestDepsState_Restored(t *testing.T) {
	hub := makeSigmaFixture(t)
	out, _ := run(t, hub)
	if !strings.Contains(out, "restored: cnos.core") {
		t.Errorf("restored fixture must list packages:\n%s", out)
	}
}

func TestDepsState_ManifestOnly(t *testing.T) {
	hub := makeInitPlusSetupFixture(t)
	out, _ := run(t, hub)
	if !strings.Contains(out, "dependency manifest present; packages not restored — run cn deps restore") {
		t.Errorf("manifest-only deps must show restore guidance:\n%s", out)
	}
}

func TestDepsState_None(t *testing.T) {
	hub := makeInitOnlyFixture(t)
	out, _ := run(t, hub)
	if !strings.Contains(out, "no dependency manifest") {
		t.Errorf("no-deps fixture must say 'no dependency manifest':\n%s", out)
	}
}

// --- AC9: Ordered ## Read first section ---

func TestReadFirstSection_OrderedSigma(t *testing.T) {
	hub := makeSigmaFixture(t)
	// Add a daily reflection so latest reflection appears.
	mustWrite(t, filepath.Join(hub, "threads", "reflections", "daily", "2026-05-01.md"), "# Reflection\n")

	out, _ := run(t, hub)

	readFirstIdx := strings.Index(out, "## Read first")
	if readFirstIdx < 0 {
		t.Fatalf("missing ## Read first section:\n%s", out)
	}
	// Persona (1.) must come before operator (2.) before kernel (3.)
	idx1 := strings.Index(out[readFirstIdx:], "spec/PERSONA.md")
	idx2 := strings.Index(out[readFirstIdx:], "spec/OPERATOR.md")
	idx3 := strings.Index(out[readFirstIdx:], "doctrine/KERNEL.md")
	idx4 := strings.Index(out[readFirstIdx:], "deps.json")
	idx5 := strings.Index(out[readFirstIdx:], "latest reflection")
	for name, idx := range map[string]int{
		"PERSONA.md":  idx1,
		"OPERATOR.md": idx2,
		"KERNEL.md":   idx3,
		"deps.json":   idx4,
		"reflection":  idx5,
	} {
		if idx < 0 {
			t.Errorf("## Read first missing %s:\n%s", name, out)
		}
	}
	if !(idx1 < idx2 && idx2 < idx3 && idx3 < idx4 && idx4 < idx5) {
		t.Errorf("## Read first order wrong: PERSONA=%d OPERATOR=%d KERNEL=%d deps=%d reflection=%d",
			idx1, idx2, idx3, idx4, idx5)
	}

	// AC9 negative: must not contain old "inspect files directly" as the only guidance
	if strings.Contains(out, "inspect files directly") {
		t.Error("old footer guidance must not appear")
	}
}

func TestReadFirstSection_InitOnlyOrdered(t *testing.T) {
	hub := makeInitOnlyFixture(t)
	out, _ := run(t, hub)

	readFirstIdx := strings.Index(out, "## Read first")
	if readFirstIdx < 0 {
		t.Fatalf("missing ## Read first section:\n%s", out)
	}
	// Absent entries should still appear in order with "not found" message
	idx1 := strings.Index(out[readFirstIdx:], "PERSONA.md")
	idx2 := strings.Index(out[readFirstIdx:], "OPERATOR.md")
	idx3 := strings.Index(out[readFirstIdx:], "kernel")
	if idx1 < 0 || idx2 < 0 || idx3 < 0 {
		t.Errorf("## Read first must list all three layers even when absent:\n%s", out)
	}
	if !(idx1 < idx2 && idx2 < idx3) {
		t.Errorf("## Read first order wrong for init-only: PERSONA=%d OPERATOR=%d kernel=%d",
			idx1, idx2, idx3)
	}
}

// --- AC10: Latest reflection pointer ---

func TestLatestReflection_Present(t *testing.T) {
	hub := makeSigmaFixture(t)
	mustWrite(t, filepath.Join(hub, "threads", "reflections", "daily", "2026-04-29.md"), "older\n")
	mustWrite(t, filepath.Join(hub, "threads", "reflections", "daily", "2026-05-01.md"), "latest\n")

	out, _ := run(t, hub)

	if !strings.Contains(out, "2026-05-01.md") {
		t.Errorf("latest reflection must reference 2026-05-01.md:\n%s", out)
	}
	if strings.Contains(out, "2026-04-29.md") && strings.Index(out, "2026-04-29.md") < strings.Index(out, "2026-05-01.md") {
		// Older file should not be listed as the latest
		t.Error("older reflection file must not precede latest in Read first")
	}
}

func TestLatestReflection_Empty(t *testing.T) {
	hub := makeInitOnlyFixture(t)
	out, _ := run(t, hub)

	// No reflection dir → no "latest reflection" line and no "not present" text
	if strings.Contains(out, "not present") {
		t.Error("empty reflection dir must not produce 'not present' text")
	}
}

// --- AC1: scanIdentity / ## Identity removed ---

func TestNoIdentityBucket(t *testing.T) {
	for _, name := range []string{"sigma", "init-only", "init+setup"} {
		var hub string
		switch name {
		case "sigma":
			hub = makeSigmaFixture(t)
		case "init-only":
			hub = makeInitOnlyFixture(t)
		case "init+setup":
			hub = makeInitPlusSetupFixture(t)
		}
		out, _ := run(t, hub)
		if strings.Contains(out, "## Identity") {
			t.Errorf("[%s] stdout must not contain legacy ## Identity bucket:\n%s", name, out)
		}
		if strings.Contains(out, "no identity files found") {
			t.Errorf("[%s] stdout must not contain 'no identity files found':\n%s", name, out)
		}
	}
}

// --- Original positive tests (updated) ---

func TestRunPositive_CwdHub(t *testing.T) {
	hub := makeMinimalHub(t, "test-agent")

	out, _ := run(t, hub)

	if !strings.Contains(out, "You are activating a cnos hub.") {
		t.Error("prompt missing activation header")
	}
	if !strings.Contains(out, hub) {
		t.Errorf("prompt missing hub path %q", hub)
	}
	if !strings.Contains(out, "Hub name: test-agent") {
		t.Error("prompt missing hub name")
	}
	// Persona file listed (makeMinimalHub creates spec/PERSONA.md)
	if !strings.Contains(out, "spec/PERSONA.md: present") {
		t.Error("prompt missing persona file in ## Persona section")
	}
	// No legacy Identity bucket
	if strings.Contains(out, "## Identity") {
		t.Error("must not have legacy ## Identity section")
	}
}

func TestRunPositive_ExplicitHubDir(t *testing.T) {
	hub := makeMinimalHub(t, "explicit-agent")
	out, _ := run(t, hub)
	if !strings.Contains(out, "You are activating a cnos hub.") {
		t.Error("prompt missing header")
	}
}

func TestRunPositive_PackagesListed(t *testing.T) {
	hub := makeSigmaFixture(t)
	// sigma fixture already has cnos.core restored
	out, _ := run(t, hub)
	if !strings.Contains(out, "cnos.core") {
		t.Errorf("prompt missing package names:\n%s", out)
	}
	if !strings.Contains(out, "restored:") {
		t.Errorf("prompt missing restored state:\n%s", out)
	}
}

func TestRunPositive_SecretsExcluded(t *testing.T) {
	hub := makeMinimalHub(t, "secret-agent")
	secret := "MY_TOKEN=super-secret-token-xyz-do-not-leak"
	mustWrite(t, filepath.Join(hub, ".cn", "secrets.env"), secret)

	out, _ := run(t, hub)

	if strings.Contains(out, "super-secret-token-xyz-do-not-leak") {
		t.Error("stdout contains secret content — secrets.env must not be included in prompt")
	}
	if !strings.Contains(out, "secrets.env") {
		t.Error("prompt should mention secrets.env exclusion in Notes section")
	}
}

func TestRunPositive_StdoutOnly(t *testing.T) {
	hub := makeMinimalHub(t, "stdout-agent")

	var stdout, stderr bytes.Buffer
	if err := Run(context.Background(), Options{HubPath: hub, Stdout: &stdout, Stderr: &stderr}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if strings.Contains(stdout.String(), "→") {
		t.Error("diagnostic arrow (→) found in stdout — diagnostics must go to stderr")
	}
	if !strings.Contains(stderr.String(), "Generating activation prompt") {
		t.Error("expected diagnostic in stderr")
	}
	if !strings.Contains(stdout.String(), "You are activating a cnos hub.") {
		t.Error("prompt missing from stdout")
	}
}

func TestRunPositive_NoModelInvocation(t *testing.T) {
	hub := makeMinimalHub(t, "nomodel-agent")
	out, _ := run(t, hub)
	if !strings.Contains(out, "No model is invoked.") {
		t.Error("prompt should state that no model is invoked")
	}
}

// --- Negative tests ---

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
	hub := makeMinimalHub(t, "cli-example-agent")
	out, _ := run(t, hub)
	if !strings.Contains(out, `claude -p "Activate this cnos hub`) {
		t.Errorf("prompt missing valid Claude CLI example with query:\n%s", out)
	}
}
