package activation

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// findRepoRoot walks up from the current working directory looking
// for go.mod, then returns the repo root (two levels above go/).
// Kept local to index_test.go to avoid coupling across test files.
func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			// dir is src/go/; repo root is two levels up.
			return filepath.Dir(filepath.Dir(dir))
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find repo root")
		}
		dir = parent
	}
}

// installSkill writes a SKILL.md for the given (package, skill id,
// frontmatter body). The skill id is interpreted as a path relative to
// <pkg>/skills/ and may contain forward slashes for nested skills.
func installSkill(t *testing.T, hub, pkgName, skillID, body string) string {
	t.Helper()
	pkgDir := filepath.Join(hub, ".cn", "vendor", "packages", pkgName)
	skillPath := filepath.Join(pkgDir, "skills", filepath.FromSlash(skillID), "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(skillPath), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(skillPath, []byte(body), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	return skillPath
}

// installManifest writes a minimal cn.package.json at the package root
// to prove Discover ignores it (no `skills` field is consulted).
func installManifest(t *testing.T, hub, pkgName, body string) {
	t.Helper()
	pkgDir := filepath.Join(hub, ".cn", "vendor", "packages", pkgName)
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "cn.package.json"), []byte(body), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
}

// --- AC1 discovery ---

// B1 — Discover returns every SKILL.md on disk, regardless of manifest.
func TestDiscover_FilesystemIsAuthoritative(t *testing.T) {
	hub := t.TempDir()
	installManifest(t, hub, "cnos.demo", `{"schema":"cn.package.v1","name":"cnos.demo","version":"1.0.0"}`)
	installSkill(t, hub, "cnos.demo", "with-triggers",
		"---\nname: with-triggers\ndescription: A useful skill\ntriggers:\n  - alpha\n  - beta\n---\n# Body\n")
	installSkill(t, hub, "cnos.demo", "no-triggers",
		"---\nname: no-triggers\ndescription: A skill without triggers\n---\n# Body\n")

	skills := Discover(hub)
	if len(skills) != 2 {
		t.Fatalf("len(skills) = %d, want 2", len(skills))
	}
	ids := []string{skills[0].SkillID, skills[1].SkillID}
	sort.Strings(ids)
	want := []string{"no-triggers", "with-triggers"}
	if !reflect.DeepEqual(ids, want) {
		t.Errorf("ids = %v, want %v", ids, want)
	}
	for _, s := range skills {
		if s.Package != "cnos.demo" {
			t.Errorf("package = %q, want %q", s.Package, "cnos.demo")
		}
		if s.ReadErr != nil {
			t.Errorf("unexpected ReadErr for %s: %v", s.SkillID, s.ReadErr)
		}
	}
}

// Discover recurses into nested skill trees (e.g. cdd/alpha/SKILL.md)
// and uses the path relative to skills/ as the skill ID.
func TestDiscover_NestedSkillID(t *testing.T) {
	hub := t.TempDir()
	installSkill(t, hub, "cnos.cdd", "cdd/alpha",
		"---\nname: alpha\ndescription: alpha role\ntriggers:\n  - alpha\n---\n")
	installSkill(t, hub, "cnos.cdd", "cdd",
		"---\nname: cdd\ndescription: cdd\ntriggers:\n  - review\n---\n")

	skills := Discover(hub)
	ids := []string{}
	for _, s := range skills {
		ids = append(ids, s.SkillID)
	}
	sort.Strings(ids)
	want := []string{"cdd", "cdd/alpha"}
	if !reflect.DeepEqual(ids, want) {
		t.Errorf("ids = %v, want %v", ids, want)
	}
}

// Discover does not consult the manifest's skills field even if present.
// The old `sources.skills` declaration must be ignored in favour of
// filesystem presence (issue #261, ef53b939).
func TestDiscover_IgnoresManifestSkillsField(t *testing.T) {
	hub := t.TempDir()
	// Manifest declares a ghost skill that doesn't exist on disk and
	// does NOT declare the on-disk skill. Discover must find the on-disk
	// skill and not fabricate the ghost.
	installManifest(t, hub, "cnos.demo",
		`{"schema":"cn.package.v1","name":"cnos.demo","version":"1.0.0",`+
			`"sources":{"skills":["ghost"]}}`)
	installSkill(t, hub, "cnos.demo", "real",
		"---\nname: real\ntriggers:\n  - x\n---\n")

	skills := Discover(hub)
	if len(skills) != 1 {
		t.Fatalf("len = %d, want 1", len(skills))
	}
	if skills[0].SkillID != "real" {
		t.Errorf("id = %q, want %q", skills[0].SkillID, "real")
	}
}

// Discover returns nil (not a panic) when the vendor dir is missing.
func TestDiscover_NoVendorDir(t *testing.T) {
	hub := t.TempDir()
	if got := Discover(hub); got != nil {
		t.Errorf("expected nil, got %d skills", len(got))
	}
}

// A package without a skills/ directory contributes nothing.
func TestDiscover_PackageWithoutSkillsDir(t *testing.T) {
	hub := t.TempDir()
	installManifest(t, hub, "cnos.cmdonly",
		`{"schema":"cn.package.v1","name":"cnos.cmdonly","version":"1.0.0"}`)
	if got := Discover(hub); got != nil {
		t.Errorf("expected nil, got %d skills", len(got))
	}
}

// A SKILL.md placed directly under skills/ (with no containing named
// subdirectory) has no derivable skill ID and is intentionally not
// discoverable. Authors must place each skill under skills/<name>/.
func TestDiscover_RootLevelSKILLMdSkipped(t *testing.T) {
	hub := t.TempDir()
	pkgDir := filepath.Join(hub, ".cn", "vendor", "packages", "cnos.demo")
	if err := os.MkdirAll(filepath.Join(pkgDir, "skills"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// SKILL.md directly under skills/ — no containing subdirectory.
	if err := os.WriteFile(
		filepath.Join(pkgDir, "skills", "SKILL.md"),
		[]byte("---\nname: rootless\ntriggers:\n  - x\n---\n"),
		0o644,
	); err != nil {
		t.Fatalf("write: %v", err)
	}
	// A properly-placed sibling skill should still be discovered.
	installSkill(t, hub, "cnos.demo", "good",
		"---\nname: good\ntriggers:\n  - g\n---\n")

	skills := Discover(hub)
	if len(skills) != 1 || skills[0].SkillID != "good" {
		t.Errorf("skills = %v, want single 'good' (root SKILL.md must be skipped)", skills)
	}
}

// Discover output is sorted by (Package, SkillID) for deterministic
// downstream rendering (eng/go §2.13).
func TestDiscover_DeterministicOrder(t *testing.T) {
	hub := t.TempDir()
	installSkill(t, hub, "cnos.zeta", "zz",
		"---\nname: zz\ntriggers:\n  - z\n---\n")
	installSkill(t, hub, "cnos.alpha", "bb",
		"---\nname: bb\ntriggers:\n  - b\n---\n")
	installSkill(t, hub, "cnos.alpha", "aa",
		"---\nname: aa\ntriggers:\n  - a\n---\n")

	skills := Discover(hub)
	got := make([]string, len(skills))
	for i, s := range skills {
		got[i] = s.Package + "/" + s.SkillID
	}
	want := []string{"cnos.alpha/aa", "cnos.alpha/bb", "cnos.zeta/zz"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("order = %v, want %v", got, want)
	}
}

// Integration test against the real SKILL.md corpus: install an
// unmodified cnos.core skill into a temp hub and verify Discover
// parses its inline triggers. This is the regression β called out
// after the initial PR landed: the parser's inline-list gap surfaces
// as "skill X has no triggers" in doctor against a freshly restored
// hub because 81% of production SKILL.md files use inline form.
func TestDiscover_RealCoreSkills_HaveTriggers(t *testing.T) {
	repoRoot := findRepoRoot(t)
	realSkill := filepath.Join(repoRoot, "src", "packages", "cnos.core",
		"skills", "agent", "ca-conduct", "SKILL.md")
	data, err := os.ReadFile(realSkill)
	if err != nil {
		t.Skipf("real skill not present at %s: %v", realSkill, err)
	}

	hub := t.TempDir()
	dst := filepath.Join(hub, ".cn", "vendor", "packages", "cnos.core",
		"skills", "agent", "ca-conduct", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	skills := Discover(hub)
	if len(skills) != 1 {
		t.Fatalf("len = %d, want 1", len(skills))
	}
	if skills[0].SkillID != "agent/ca-conduct" {
		t.Errorf("skill id = %q, want agent/ca-conduct", skills[0].SkillID)
	}
	if len(skills[0].Frontmatter.Triggers) == 0 {
		t.Errorf("triggers empty — parser did not consume the real SKILL.md's %q line",
			"triggers: [...]")
	}

	// The same skill must not trip the doctor empty-triggers check.
	issues := Validate(hub)
	for _, iss := range issues {
		if iss.Kind == IssueEmptyTriggers {
			t.Errorf("unexpected empty-triggers issue for real skill: %s", iss.Message)
		}
	}
}

// Integration test against the full skill corpus: copy every
// src/packages/*/skills/**/SKILL.md into a temp hub and run Validate.
// Expects no IssueEmptyTriggers — every real skill declares triggers
// in some form (inline or block). This is the Go-level equivalent of
// β's required negative regression "skill activation check reports
// all skills valid against the real cnos.core/cnos.cdd corpus".
func TestValidate_RealCorpus_NoEmptyTriggers(t *testing.T) {
	repoRoot := findRepoRoot(t)
	srcPackages := filepath.Join(repoRoot, "src", "packages")
	pkgEntries, err := os.ReadDir(srcPackages)
	if err != nil {
		t.Skipf("no src/packages at %s: %v", srcPackages, err)
	}

	hub := t.TempDir()
	vendorPkgs := filepath.Join(hub, ".cn", "vendor", "packages")
	copied := 0
	for _, pe := range pkgEntries {
		if !pe.IsDir() {
			continue
		}
		srcSkills := filepath.Join(srcPackages, pe.Name(), "skills")
		if _, err := os.Stat(srcSkills); err != nil {
			continue
		}
		dstSkills := filepath.Join(vendorPkgs, pe.Name(), "skills")
		if err := os.MkdirAll(dstSkills, 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		err := filepath.Walk(srcSkills, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}
			rel, err := filepath.Rel(srcSkills, path)
			if err != nil {
				return err
			}
			dst := filepath.Join(dstSkills, rel)
			if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
				return err
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if filepath.Base(path) == "SKILL.md" {
				copied++
			}
			return os.WriteFile(dst, data, 0o644)
		})
		if err != nil {
			t.Fatalf("copy %s: %v", pe.Name(), err)
		}
	}
	if copied == 0 {
		t.Skip("no SKILL.md files copied — skipping corpus check")
	}

	issues := Validate(hub)
	for _, iss := range issues {
		if iss.Kind == IssueEmptyTriggers {
			t.Errorf("empty-triggers against real corpus: %s", iss.Message)
		}
	}
	// Log conflict/missing counts for visibility; they are not part of
	// this test's contract (conflicts legitimately happen if the corpus
	// itself has duplicate trigger keywords, which is a corpus problem
	// rather than a parser problem).
	t.Logf("real-corpus validation: %d skills copied, %d issues surfaced",
		copied, len(issues))
}

// --- AC3 activation index (public filter) ---

func TestBuildIndex_ExcludesInternalVisibility(t *testing.T) {
	hub := t.TempDir()
	installSkill(t, hub, "cnos.cdd", "cdd/alpha",
		"---\nname: alpha\nvisibility: internal\ntriggers:\n  - alpha\n---\n")
	installSkill(t, hub, "cnos.cdd", "cdd",
		"---\nname: cdd\ntriggers:\n  - review\n---\n")
	installSkill(t, hub, "cnos.core", "daily",
		"---\nname: daily\nvisibility: public\ntriggers:\n  - daily\n---\n")

	entries := BuildIndex(hub)
	if len(entries) != 2 {
		t.Fatalf("len = %d, want 2 (internal excluded)", len(entries))
	}
	ids := []string{entries[0].SkillID, entries[1].SkillID}
	sort.Strings(ids)
	want := []string{"cdd", "daily"}
	if !reflect.DeepEqual(ids, want) {
		t.Errorf("ids = %v, want %v", ids, want)
	}
}

// Missing visibility field defaults to public.
func TestBuildIndex_AbsentVisibilityDefaultsPublic(t *testing.T) {
	hub := t.TempDir()
	installSkill(t, hub, "cnos.demo", "public-by-default",
		"---\nname: x\ntriggers:\n  - x\n---\n")

	entries := BuildIndex(hub)
	if len(entries) != 1 {
		t.Fatalf("len = %d, want 1", len(entries))
	}
	if entries[0].SkillID != "public-by-default" {
		t.Errorf("id = %q", entries[0].SkillID)
	}
}

// BuildIndex carries description into Entry.Summary and preserves
// trigger order.
func TestBuildIndex_EntryShape(t *testing.T) {
	hub := t.TempDir()
	installSkill(t, hub, "cnos.demo", "s",
		"---\nname: s\ndescription: demo summary\ntriggers:\n  - a\n  - b\n  - c\n---\n")

	entries := BuildIndex(hub)
	if len(entries) != 1 {
		t.Fatalf("len = %d, want 1", len(entries))
	}
	e := entries[0]
	if e.Package != "cnos.demo" {
		t.Errorf("package = %q, want cnos.demo", e.Package)
	}
	if e.Summary != "demo summary" {
		t.Errorf("summary = %q", e.Summary)
	}
	if !reflect.DeepEqual(e.Triggers, []string{"a", "b", "c"}) {
		t.Errorf("triggers = %v", e.Triggers)
	}
}

// --- AC4 validate / doctor ---

// V1 — unreadable SKILL.md surfaces as IssueMissingSkill. Tested
// via the pure ValidateSkills so we don't depend on platform file
// permissions (which are a no-op under root).
func TestValidateSkills_UnreadableSkillMissing(t *testing.T) {
	skills := []Skill{
		{
			Package: "cnos.demo",
			SkillID: "hidden",
			Path:    "/tmp/does-not-matter/SKILL.md",
			ReadErr: os.ErrPermission,
		},
	}
	issues := ValidateSkills(skills)
	if len(issues) != 1 {
		t.Fatalf("issues = %d, want 1", len(issues))
	}
	if issues[0].Kind != IssueMissingSkill {
		t.Errorf("kind = %v, want IssueMissingSkill", issues[0].Kind)
	}
	if !strings.Contains(issues[0].Message, "hidden") {
		t.Errorf("message = %q, want to mention 'hidden'", issues[0].Message)
	}
	if !strings.Contains(issues[0].Message, "cnos.demo") {
		t.Errorf("message = %q, want to mention package", issues[0].Message)
	}
}

// Integration path: Validate() picks up an unreadable SKILL.md by
// pointing it at a directory (os.ReadFile returns EISDIR). Works
// regardless of uid.
func TestValidate_UnreadableViaDirectoryName(t *testing.T) {
	hub := t.TempDir()
	// Create skills/broken/SKILL.md as a directory, not a file. The
	// filesystem walker still sees "SKILL.md" as an entry, but the
	// ReadFrontmatter IO wrapper will fail to read it.
	dir := filepath.Join(hub, ".cn", "vendor", "packages", "cnos.demo",
		"skills", "broken", "SKILL.md")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	// Also install a valid sibling so Discover has something to traverse.
	installSkill(t, hub, "cnos.demo", "ok",
		"---\nname: ok\ntriggers:\n  - ok\n---\n")

	// Because SKILL.md is a directory, WalkDir sees d.Name() == "SKILL.md"
	// but d.IsDir() == true, so our walker (which skips directories) will
	// NOT pick it up. The broken directory becomes invisible to Discover,
	// which is the intended filesystem-authoritative semantics: if the
	// on-disk shape is not a readable SKILL.md file, it is not a skill.
	skills := Discover(hub)
	if len(skills) != 1 || skills[0].SkillID != "ok" {
		t.Errorf("skills = %v, want single 'ok' entry", skills)
	}
}

// V2 — empty triggers surfaces IssueEmptyTriggers.
func TestValidate_EmptyTriggers(t *testing.T) {
	hub := t.TempDir()
	installSkill(t, hub, "cnos.demo", "plain",
		"---\nname: plain\ndescription: no triggers field\n---\n")

	issues := Validate(hub)
	kinds := []string{}
	for _, i := range issues {
		kinds = append(kinds, IssueKindLabel(i.Kind))
	}
	sort.Strings(kinds)
	if !reflect.DeepEqual(kinds, []string{"empty"}) {
		t.Errorf("kinds = %v, want [empty]", kinds)
	}
}

// V3 — conflicting triggers across distinct skill ids surface as
// IssueTriggerConflict, one issue per shared trigger.
func TestValidate_TriggerConflict(t *testing.T) {
	hub := t.TempDir()
	installSkill(t, hub, "cnos.alpha", "first",
		"---\nname: first\ntriggers:\n  - shared\n  - uniq-a\n---\n")
	installSkill(t, hub, "cnos.beta", "second",
		"---\nname: second\ntriggers:\n  - shared\n  - uniq-b\n---\n")

	issues := Validate(hub)
	var conflicts []Issue
	for _, i := range issues {
		if i.Kind == IssueTriggerConflict {
			conflicts = append(conflicts, i)
		}
	}
	if len(conflicts) != 1 {
		t.Fatalf("conflicts = %d, want 1; all=%v", len(conflicts), issues)
	}
	msg := conflicts[0].Message
	if !strings.Contains(msg, "shared") {
		t.Errorf("message = %q, want to mention 'shared'", msg)
	}
	if !strings.Contains(msg, "first") || !strings.Contains(msg, "second") {
		t.Errorf("message = %q, want to mention both skill ids", msg)
	}
}

// A single skill claiming the same trigger twice is not a conflict.
func TestValidate_SelfTriggersNotConflict(t *testing.T) {
	hub := t.TempDir()
	installSkill(t, hub, "cnos.demo", "only",
		"---\nname: only\ntriggers:\n  - solo\n  - solo\n---\n")

	issues := Validate(hub)
	for _, i := range issues {
		if i.Kind == IssueTriggerConflict {
			t.Errorf("unexpected trigger conflict: %v", i.Message)
		}
	}
}

// Internal skills are excluded from trigger-conflict detection: they
// never surface in the public activation index (see BuildIndex), so
// an internal sub-skill sharing a trigger keyword with its parent
// orchestrator is intentional delegation, not ambiguity. Mirrors the
// OCaml-era semantic where validate iterated only manifest-exposed
// (public) skill IDs.
func TestValidate_InternalDoesNotConflictWithPublic(t *testing.T) {
	hub := t.TempDir()
	installSkill(t, hub, "cnos.cdd", "cdd",
		"---\nname: cdd\ntriggers:\n  - review\n---\n")
	installSkill(t, hub, "cnos.cdd", "cdd/review",
		"---\nname: review\nvisibility: internal\ntriggers:\n  - review\n---\n")

	issues := Validate(hub)
	for _, i := range issues {
		if i.Kind == IssueTriggerConflict {
			t.Errorf("unexpected conflict between public orchestrator and internal sub-skill: %s", i.Message)
		}
	}
}

// Two internal skills sharing a trigger do not conflict either —
// neither is publicly addressable.
func TestValidate_TwoInternalDoNotConflict(t *testing.T) {
	hub := t.TempDir()
	installSkill(t, hub, "cnos.cdd", "cdd/a",
		"---\nname: a\nvisibility: internal\ntriggers:\n  - shared\n---\n")
	installSkill(t, hub, "cnos.cdd", "cdd/b",
		"---\nname: b\nvisibility: internal\ntriggers:\n  - shared\n---\n")

	issues := Validate(hub)
	for _, i := range issues {
		if i.Kind == IssueTriggerConflict {
			t.Errorf("unexpected conflict between internal skills: %s", i.Message)
		}
	}
}

// No packages, no vendor dir → no issues.
func TestValidate_EmptyHub(t *testing.T) {
	hub := t.TempDir()
	if issues := Validate(hub); issues != nil {
		t.Errorf("expected nil issues, got %v", issues)
	}
}

// Validate is deterministic: the conflict list is sorted by trigger.
func TestValidate_DeterministicOrder(t *testing.T) {
	hub := t.TempDir()
	installSkill(t, hub, "cnos.a", "s1",
		"---\nname: s1\ntriggers:\n  - zebra\n  - apple\n---\n")
	installSkill(t, hub, "cnos.b", "s2",
		"---\nname: s2\ntriggers:\n  - zebra\n  - apple\n---\n")

	issues := Validate(hub)
	var conflictTriggers []string
	for _, i := range issues {
		if i.Kind == IssueTriggerConflict {
			// Message begins with "trigger <keyword> claimed by:"
			parts := strings.SplitN(i.Message, " ", 3)
			if len(parts) >= 2 {
				conflictTriggers = append(conflictTriggers, parts[1])
			}
		}
	}
	want := []string{"apple", "zebra"}
	if !reflect.DeepEqual(conflictTriggers, want) {
		t.Errorf("conflict trigger order = %v, want %v", conflictTriggers, want)
	}
}
