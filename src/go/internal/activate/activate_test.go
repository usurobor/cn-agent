package activate

import (
	"bytes"
	"context"
	"fmt"
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
		"threads/inbox",
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

// --- AC9 / #379: Ordered ## Read first section ---
// Canonical ordering per cnos.core/skills/agent/activate/SKILL.md §2.1 / §4.1:
// kernel → ca-skills → persona → operator → hub-state (deps + reflection) → identity.

func TestReadFirstSection_OrderedSigma(t *testing.T) {
	hub := makeSigmaFixture(t)
	// Add a daily reflection so the hub-state composite line includes it.
	mustWrite(t, filepath.Join(hub, "threads", "reflections", "daily", "2026-05-01.md"), "# Reflection\n")

	out, _ := run(t, hub)

	readFirstIdx := strings.Index(out, "## Read first")
	if readFirstIdx < 0 {
		t.Fatalf("missing ## Read first section:\n%s", out)
	}
	// New canonical order: KERNEL → ca-skills → PERSONA → OPERATOR → deps → reflection → identity.
	idxKernel := strings.Index(out[readFirstIdx:], "doctrine/KERNEL.md")
	idxCASkills := strings.Index(out[readFirstIdx:], "(CA skill set)")
	idxPersona := strings.Index(out[readFirstIdx:], "spec/PERSONA.md")
	idxOperator := strings.Index(out[readFirstIdx:], "spec/OPERATOR.md")
	idxDeps := strings.Index(out[readFirstIdx:], "deps.json")
	idxReflection := strings.Index(out[readFirstIdx:], "latest reflection")
	idxIdentity := strings.Index(out[readFirstIdx:], "identity confirmation")
	for name, idx := range map[string]int{
		"KERNEL.md":     idxKernel,
		"CA skill set":  idxCASkills,
		"PERSONA.md":    idxPersona,
		"OPERATOR.md":   idxOperator,
		"deps.json":     idxDeps,
		"reflection":    idxReflection,
		"identity conf": idxIdentity,
	} {
		if idx < 0 {
			t.Errorf("## Read first missing %s:\n%s", name, out)
		}
	}
	if !(idxKernel < idxCASkills && idxCASkills < idxPersona && idxPersona < idxOperator &&
		idxOperator < idxDeps && idxDeps < idxReflection && idxReflection < idxIdentity) {
		t.Errorf("## Read first order wrong: KERNEL=%d CA=%d PERSONA=%d OPERATOR=%d deps=%d reflection=%d identity=%d",
			idxKernel, idxCASkills, idxPersona, idxOperator, idxDeps, idxReflection, idxIdentity)
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
	// Even with absent kernel/persona/operator, the canonical six-token order is preserved.
	idxKernel := strings.Index(out[readFirstIdx:], "kernel")
	idxCASkills := strings.Index(out[readFirstIdx:], "(CA skill set)")
	idxPersona := strings.Index(out[readFirstIdx:], "PERSONA.md")
	idxOperator := strings.Index(out[readFirstIdx:], "OPERATOR.md")
	if idxKernel < 0 || idxCASkills < 0 || idxPersona < 0 || idxOperator < 0 {
		t.Errorf("## Read first must list all four layers even when absent:\n%s", out)
	}
	if !(idxKernel < idxCASkills && idxCASkills < idxPersona && idxPersona < idxOperator) {
		t.Errorf("## Read first order wrong for init-only: kernel=%d CA=%d PERSONA=%d OPERATOR=%d",
			idxKernel, idxCASkills, idxPersona, idxOperator)
	}
}

// --- #379 AC7: Skill is the source of truth for ## Read first ordering ---

// vendoredActivateSkillPath returns the path inside a hub where the activate
// skill is expected to live for the renderer to consume it.
func vendoredActivateSkillPath(hub string) string {
	return filepath.Join(hub, ".cn", "vendor", "packages", "cnos.core",
		"skills", "agent", "activate", "SKILL.md")
}

// writeVendoredActivateSkill writes a minimal activate SKILL.md to the
// hub's vendored bundle with a §4.1 machine-readable load-order block
// containing the given (token, label) pairs in order.
func writeVendoredActivateSkill(t *testing.T, hub string, items [][2]string) {
	t.Helper()
	var b strings.Builder
	b.WriteString("---\nname: activate\nartifact_class: skill\n---\n\n")
	b.WriteString("# Activate (fixture)\n\n## 4. Renderer contract\n\n")
	b.WriteString("<!-- read-first-order:begin -->\n")
	for i, pair := range items {
		fmt.Fprintf(&b, "%d. %s — %s\n", i+1, pair[0], pair[1])
	}
	b.WriteString("<!-- read-first-order:end -->\n")
	mustWrite(t, vendoredActivateSkillPath(hub), b.String())
}

// TestSkillIsSourceOfTruthForReadFirstOrder demonstrates AC7: when the
// activate skill is vendored, editing its §4.1 block changes the renderer's
// ## Read first emission order. The skill — not in-Go constants — owns
// the ordering decision.
func TestSkillIsSourceOfTruthForReadFirstOrder(t *testing.T) {
	hub := makeSigmaFixture(t)
	// Add a daily reflection so the hub-state composite line is non-trivial.
	mustWrite(t, filepath.Join(hub, "threads", "reflections", "daily", "2026-05-01.md"), "# Reflection\n")

	// Phase 1: vendor an activate skill that places kernel BEFORE persona.
	writeVendoredActivateSkill(t, hub, [][2]string{
		{"kernel", "Kernel doctrine"},
		{"persona", "Persona"},
		{"operator", "Operator"},
	})

	out1, stderr1 := run(t, hub)
	if strings.Contains(stderr1, "activate skill not vendored") {
		t.Fatalf("renderer fell back to built-in ordering despite vendored skill:\nstderr=%s", stderr1)
	}
	readFirstIdx1 := strings.Index(out1, "## Read first")
	if readFirstIdx1 < 0 {
		t.Fatalf("missing ## Read first section:\n%s", out1)
	}
	kIdx1 := strings.Index(out1[readFirstIdx1:], "doctrine/KERNEL.md")
	pIdx1 := strings.Index(out1[readFirstIdx1:], "spec/PERSONA.md")
	if kIdx1 < 0 || pIdx1 < 0 {
		t.Fatalf("phase 1: missing KERNEL.md or PERSONA.md in ## Read first:\n%s", out1)
	}
	if !(kIdx1 < pIdx1) {
		t.Errorf("phase 1: skill ordered kernel BEFORE persona; renderer emitted otherwise (KERNEL=%d PERSONA=%d):\n%s",
			kIdx1, pIdx1, out1)
	}

	// Phase 2: edit the vendored skill — swap kernel and persona. The
	// renderer must pick up the new ordering without code change.
	writeVendoredActivateSkill(t, hub, [][2]string{
		{"persona", "Persona"},
		{"kernel", "Kernel doctrine"},
		{"operator", "Operator"},
	})

	out2, stderr2 := run(t, hub)
	if strings.Contains(stderr2, "activate skill not vendored") {
		t.Fatalf("phase 2: renderer fell back despite vendored skill:\nstderr=%s", stderr2)
	}
	readFirstIdx2 := strings.Index(out2, "## Read first")
	if readFirstIdx2 < 0 {
		t.Fatalf("phase 2: missing ## Read first section:\n%s", out2)
	}
	kIdx2 := strings.Index(out2[readFirstIdx2:], "doctrine/KERNEL.md")
	pIdx2 := strings.Index(out2[readFirstIdx2:], "spec/PERSONA.md")
	if kIdx2 < 0 || pIdx2 < 0 {
		t.Fatalf("phase 2: missing KERNEL.md or PERSONA.md in ## Read first:\n%s", out2)
	}
	if !(pIdx2 < kIdx2) {
		t.Errorf("phase 2: skill ordered persona BEFORE kernel; renderer emitted otherwise (PERSONA=%d KERNEL=%d):\n%s",
			pIdx2, kIdx2, out2)
	}

	// Coherence check: the two outputs must differ in ## Read first ordering.
	if out1 == out2 {
		t.Error("renderer produced identical output for two different skill orderings — skill not source of truth")
	}
}

// TestSkillFallback_NotVendored confirms that when the activate skill is
// absent from the vendored bundle, the renderer falls back to the built-in
// canonical ordering and announces the fallback on stderr. This preserves
// the manifest-only deps state behavior the prior implementation guaranteed.
func TestSkillFallback_NotVendored(t *testing.T) {
	hub := makeInitPlusSetupFixture(t) // manifest-only deps, no vendored bundle

	out, stderr := run(t, hub)

	if !strings.Contains(stderr, "activate skill not vendored") {
		t.Errorf("expected stderr fallback diagnostic, got: %q", stderr)
	}
	// Canonical built-in order must still apply: kernel → ca-skills → persona → operator.
	readFirstIdx := strings.Index(out, "## Read first")
	if readFirstIdx < 0 {
		t.Fatalf("missing ## Read first section:\n%s", out)
	}
	idxKernel := strings.Index(out[readFirstIdx:], "kernel")
	idxCASkills := strings.Index(out[readFirstIdx:], "(CA skill set)")
	idxPersona := strings.Index(out[readFirstIdx:], "PERSONA.md")
	if !(idxKernel >= 0 && idxKernel < idxCASkills && idxCASkills < idxPersona) {
		t.Errorf("fallback ordering wrong: kernel=%d CA=%d PERSONA=%d:\n%s",
			idxKernel, idxCASkills, idxPersona, out)
	}
	// Manifest-only kernel guidance is preserved.
	if !strings.Contains(out, "run cn deps restore") {
		t.Errorf("manifest-only fallback must preserve 'run cn deps restore' guidance:\n%s", out)
	}
}

// TestParseReadFirstOrderBlock_HappyPath unit-tests the parser directly.
func TestParseReadFirstOrderBlock_HappyPath(t *testing.T) {
	content := "preamble\n" +
		"<!-- read-first-order:begin -->\n" +
		"1. kernel — Kernel doctrine\n" +
		"2. ca-skills — CA skill set\n" +
		"3. persona — Persona\n" +
		"<!-- read-first-order:end -->\n" +
		"trailing\n"
	items, ok := parseReadFirstOrderBlock(content)
	if !ok {
		t.Fatal("parser returned !ok")
	}
	if len(items) != 3 {
		t.Fatalf("want 3 items, got %d: %+v", len(items), items)
	}
	want := []readFirstItem{
		{token: "kernel", label: "Kernel doctrine"},
		{token: "ca-skills", label: "CA skill set"},
		{token: "persona", label: "Persona"},
	}
	for i, w := range want {
		if items[i] != w {
			t.Errorf("item[%d]: want %+v got %+v", i, w, items[i])
		}
	}
}

// TestParseReadFirstOrderBlock_NoMarkers returns !ok on a content lacking the
// begin marker, leaving the caller to choose fallback behavior.
func TestParseReadFirstOrderBlock_NoMarkers(t *testing.T) {
	if items, ok := parseReadFirstOrderBlock("nothing structured here\n"); ok {
		t.Errorf("expected !ok for content without markers, got items=%+v", items)
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
