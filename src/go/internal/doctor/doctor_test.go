package doctor

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunAllHealthyHub(t *testing.T) {
	hub := makeTestHub(t)
	checks := RunAll(context.Background(), hub, "3.48.0", nil)

	if len(checks) == 0 {
		t.Fatal("expected checks")
	}

	// Hub-structure checks should pass for makeTestHub.
	for _, ch := range checks {
		if ch.Name == ".cn/config" && ch.Status != StatusPass {
			t.Errorf(".cn/config should pass, got status=%d value=%q", ch.Status, ch.Value)
		}
		if ch.Name == "spec/SOUL.md" && ch.Status == StatusFail {
			t.Errorf("spec/SOUL.md should not fail, got status=%d value=%q", ch.Status, ch.Value)
		}
	}
}

func TestRunAllMissingConfig(t *testing.T) {
	hub := t.TempDir()
	os.MkdirAll(filepath.Join(hub, ".cn"), 0755)

	checks := RunAll(context.Background(), hub, "3.48.0", nil)

	found := false
	for _, ch := range checks {
		if ch.Name == ".cn/config" && ch.Status == StatusFail {
			found = true
		}
	}
	if !found {
		t.Error("expected .cn/config check to fail")
	}
}

func TestRunAllPackageMissing(t *testing.T) {
	hub := makeTestHub(t)
	// Install cnos.core but lockfile expects cnos.core + cnos.eng.
	os.MkdirAll(filepath.Join(hub, ".cn", "vendor", "packages", "cnos.core"), 0755)
	lockfile := `{"schema":"cn.lock.v2","packages":[{"name":"cnos.core","version":"1.0.0","sha256":"aaa"},{"name":"cnos.eng","version":"1.0.0","sha256":"bbb"}]}`
	os.WriteFile(filepath.Join(hub, ".cn", "deps.lock.json"), []byte(lockfile), 0644)

	checks := RunAll(context.Background(), hub, "3.48.0", nil)

	found := false
	for _, ch := range checks {
		if ch.Name == "packages" && ch.Status == StatusFail && strings.Contains(ch.Value, "missing from lockfile") {
			found = true
		}
	}
	if !found {
		t.Error("expected missing package in packages check")
	}
}

func TestRunAllRuntimeContract(t *testing.T) {
	hub := makeTestHub(t)
	contract := map[string]any{
		"schema":    "cn.runtime_contract.v2",
		"identity":  map[string]any{"cn_version": "3.48.0"},
		"cognition": map[string]any{},
		"body":      map[string]any{},
		"medium":    map[string]any{},
	}
	contractDir := filepath.Join(hub, "state")
	os.MkdirAll(contractDir, 0755)
	data, _ := json.MarshalIndent(contract, "", "  ")
	os.WriteFile(filepath.Join(contractDir, "runtime-contract.json"), data, 0644)

	checks := RunAll(context.Background(), hub, "3.48.0", nil)

	found := false
	for _, ch := range checks {
		if ch.Name == "runtime contract" && ch.Status == StatusPass && strings.Contains(ch.Value, "identity") {
			found = true
		}
	}
	if !found {
		t.Error("expected valid runtime contract check")
	}
}

func TestRunAllIncompleteContract(t *testing.T) {
	hub := makeTestHub(t)
	contract := map[string]any{
		"schema":    "cn.runtime_contract.v2",
		"identity":  map[string]any{},
		"cognition": map[string]any{},
		"body":      map[string]any{},
	}
	contractDir := filepath.Join(hub, "state")
	os.MkdirAll(contractDir, 0755)
	data, _ := json.MarshalIndent(contract, "", "  ")
	os.WriteFile(filepath.Join(contractDir, "runtime-contract.json"), data, 0644)

	checks := RunAll(context.Background(), hub, "3.48.0", nil)

	found := false
	for _, ch := range checks {
		if ch.Name == "runtime contract" && ch.Status == StatusFail && strings.Contains(ch.Value, "medium") {
			found = true
		}
	}
	if !found {
		t.Error("expected incomplete contract with missing medium")
	}
}

// TestFreshHubLifecycleChecksAreInfo is the Tier 1 kata contract
// (issue #236 AC4): a freshly-init'd + set-up hub has pending
// lifecycle artifacts (deps.lock.json, runtime contract, vendor
// dir pre-lock, origin remote). These must report as StatusInfo,
// not StatusFail — otherwise `cn doctor` exits non-zero on a
// healthy fresh hub.
//
// The test is scoped to hub-lifecycle state; environmental
// prerequisites (git/curl/identity) depend on the host and are
// tested separately.
func TestFreshHubLifecycleChecksAreInfo(t *testing.T) {
	hub := makeTestHub(t)
	// `cn setup` writes deps.json. The rest is as hubinit leaves it.
	os.WriteFile(filepath.Join(hub, ".cn", "deps.json"),
		[]byte(`{"schema":"cn.deps.v1","profile":"engineer","packages":[]}`), 0644)

	checks := RunAll(context.Background(), hub, "3.48.0", nil)

	wantInfo := map[string]bool{
		".cn/deps.lock.json": true,
		"packages":           true,
		"runtime contract":   true,
		"origin remote":      true,
	}
	seen := map[string]Status{}
	for _, ch := range checks {
		if wantInfo[ch.Name] {
			seen[ch.Name] = ch.Status
		}
	}
	for name := range wantInfo {
		s, ok := seen[name]
		if !ok {
			t.Errorf("expected lifecycle check %q in results", name)
			continue
		}
		if s != StatusInfo {
			t.Errorf("lifecycle check %q on fresh hub: want StatusInfo, got %d", name, s)
		}
	}
}

// TestSkillActivationCheck_PassEmpty: no installed packages means
// no skill activation problems; the check passes.
func TestSkillActivationCheck_PassEmpty(t *testing.T) {
	hub := makeTestHub(t)
	checks := RunAll(context.Background(), hub, "3.48.0", nil)

	found := false
	for _, ch := range checks {
		if ch.Name == "skill activation" {
			found = true
			if ch.Status != StatusPass {
				t.Errorf("skill activation on empty hub: status=%d value=%q, want StatusPass",
					ch.Status, ch.Value)
			}
		}
	}
	if !found {
		t.Error("expected 'skill activation' check in results")
	}
}

// TestSkillActivationCheck_WarnsOnTriggerConflict: public-public
// trigger overlap is StatusInfo (doctor prints ○, hub stays rc=0).
// Overlapping activation keywords are a legitimate authoring pattern;
// only structural breakage (unreadable SKILL.md) fails the hub.
func TestSkillActivationCheck_WarnsOnTriggerConflict(t *testing.T) {
	hub := makeTestHub(t)
	writeSkill := func(pkg, skillID, body string) {
		path := filepath.Join(hub, ".cn", "vendor", "packages", pkg,
			"skills", skillID, "SKILL.md")
		os.MkdirAll(filepath.Dir(path), 0o755)
		os.WriteFile(path, []byte(body), 0o644)
	}
	writeSkill("cnos.a", "first",
		"---\nname: first\ntriggers:\n  - shared\n---\n")
	writeSkill("cnos.b", "second",
		"---\nname: second\ntriggers:\n  - shared\n---\n")

	checks := RunAll(context.Background(), hub, "3.48.0", nil)

	found := false
	for _, ch := range checks {
		if ch.Name == "skill activation" {
			found = true
			if ch.Status != StatusInfo {
				t.Errorf("status = %d, want StatusInfo", ch.Status)
			}
			if !strings.Contains(ch.Value, "conflict") || !strings.Contains(ch.Value, "shared") {
				t.Errorf("value = %q, want to mention conflict + 'shared' trigger", ch.Value)
			}
		}
	}
	if !found {
		t.Error("expected 'skill activation' check in results")
	}
	// Scoped assertion: no skill-activation check alone should escalate
	// to StatusFail on a hub whose only "problem" is overlapping
	// triggers. (The base test hub has other legitimate failures
	// around env/hub setup that are orthogonal to activation.)
	for _, ch := range checks {
		if ch.Name == "skill activation" && ch.Status == StatusFail {
			t.Errorf("skill activation escalated to StatusFail on a conflict-only hub: %q", ch.Value)
		}
	}
}

// Note: the StatusFail mapping for IssueMissingSkill is one-line
// wiring (see checkSkillActivation). The underlying Issue generation
// is covered by activation.TestValidateSkills_UnreadableSkillMissing.
// Adding a doctor-level test of the mapping would require platform-
// specific tricks to produce an unreadable file on linux+root; the
// coverage cost/benefit doesn't justify it.

func makeTestHub(t *testing.T) string {
	t.Helper()
	hub := t.TempDir()
	for _, d := range []string{".cn", "spec", "state", "agent"} {
		os.MkdirAll(filepath.Join(hub, d), 0755)
	}
	os.WriteFile(filepath.Join(hub, ".cn", "config.json"),
		[]byte(`{"name":"test","version":"1.0.0"}`), 0644)
	os.WriteFile(filepath.Join(hub, "spec", "SOUL.md"),
		[]byte("# SOUL\n"), 0644)
	os.WriteFile(filepath.Join(hub, "state", "peers.md"),
		[]byte("# Peers\n"), 0644)
	return hub
}
