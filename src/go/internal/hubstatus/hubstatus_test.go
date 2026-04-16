package hubstatus

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeManifest creates a cn.package.json in the vendor dir for the
// given package name. Helper — keeps test setup concise.
func writeManifest(t *testing.T, hub, name, manifest string) {
	t.Helper()
	dir := filepath.Join(hub, ".cn", "vendor", "packages", name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "cn.package.json"), []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}
}

// writeContentDir creates a non-empty content-class directory inside an
// installed package dir. Used by tests that exercise the filesystem-based
// content-class discovery path (pkgbuild.FindContentClasses).
func writeContentDir(t *testing.T, hub, pkgName, class string) {
	t.Helper()
	dir := filepath.Join(hub, ".cn", "vendor", "packages", pkgName, class)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	// A single marker file satisfies the "non-empty" check.
	if err := os.WriteFile(filepath.Join(dir, ".keep"), []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestRunShowsPackages(t *testing.T) {
	hub := t.TempDir()
	writeManifest(t, hub, "cnos.core", `{
		"schema": "cn.package.v1",
		"name": "cnos.core",
		"version": "3.52.0",
		"kind": "package",
		"engines": {"cnos": "3.52.0"},
		"skills": {"exposed": ["agent/cap"]},
		"commands": {"daily": {"entrypoint": "commands/daily/cn-daily", "summary": "Daily"}}
	}`)
	writeManifest(t, hub, "cnos.eng", `{
		"schema": "cn.package.v1",
		"name": "cnos.eng",
		"version": "3.52.0",
		"kind": "package",
		"engines": {"cnos": "3.52.0"},
		"skills": {"exposed": ["eng/go"]}
	}`)

	var stdout bytes.Buffer
	if err := Run(hub, "3.52.0", nil, &stdout); err != nil {
		t.Fatalf("status: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "cnos.core") {
		t.Error("expected cnos.core in output")
	}
	if !strings.Contains(out, "cnos.eng") {
		t.Error("expected cnos.eng in output")
	}
	if !strings.Contains(out, "3.52.0") {
		t.Error("expected version 3.52.0")
	}
	if !strings.Contains(out, "\u2713") {
		t.Error("expected check mark for matching version")
	}
}

func TestRunContentClasses(t *testing.T) {
	hub := t.TempDir()
	writeManifest(t, hub, "cnos.core", `{
		"schema": "cn.package.v1",
		"name": "cnos.core",
		"version": "3.52.0",
		"kind": "package",
		"engines": {"cnos": "3.52.0"}
	}`)
	// Content classes are discovered from the package directory tree,
	// not from manifest JSON fields. Presence is the single authority
	// (PACKAGE-SYSTEM.md §3).
	writeContentDir(t, hub, "cnos.core", "skills")
	writeContentDir(t, hub, "cnos.core", "commands")

	var stdout bytes.Buffer
	if err := Run(hub, "3.52.0", nil, &stdout); err != nil {
		t.Fatalf("status: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "[skills, commands]") {
		t.Errorf("expected content classes [skills, commands] in output, got:\n%s", out)
	}
}

// TestRunContentClassesAllEight verifies that cn status surfaces every
// content class enumerated in pkg.ContentClasses when the installed
// package contains all of them. This is the cn status side of the
// convergence proved by TestFindContentClassesAll in pkgbuild.
func TestRunContentClassesAllEight(t *testing.T) {
	hub := t.TempDir()
	writeManifest(t, hub, "cnos.full", `{
		"schema": "cn.package.v1",
		"name": "cnos.full",
		"version": "3.52.0",
		"kind": "package",
		"engines": {"cnos": "3.52.0"}
	}`)
	for _, class := range []string{
		"doctrine", "mindsets", "skills", "extensions",
		"templates", "commands", "orchestrators", "katas",
	} {
		writeContentDir(t, hub, "cnos.full", class)
	}

	var stdout bytes.Buffer
	if err := Run(hub, "3.52.0", nil, &stdout); err != nil {
		t.Fatalf("status: %v", err)
	}

	// Expect the canonical order from pkg.ContentClasses.
	out := stdout.String()
	want := "[doctrine, mindsets, skills, extensions, templates, commands, orchestrators, katas]"
	if !strings.Contains(out, want) {
		t.Errorf("expected content classes %s in output, got:\n%s", want, out)
	}
}

func TestRunVersionDrift(t *testing.T) {
	hub := t.TempDir()
	writeManifest(t, hub, "cnos.core", `{
		"schema": "cn.package.v1",
		"name": "cnos.core",
		"version": "3.48.0",
		"kind": "package",
		"engines": {"cnos": "3.48.0"}
	}`)

	var stdout bytes.Buffer
	Run(hub, "3.52.0", nil, &stdout)

	out := stdout.String()
	if !strings.Contains(out, "\u2717") {
		t.Error("expected cross mark for drift")
	}
	if !strings.Contains(out, "engines.cnos 3.48.0") {
		t.Error("expected engines.cnos drift info")
	}
	if !strings.Contains(out, "version_drift") {
		t.Error("expected version_drift warning")
	}
}

func TestRunNoPackages(t *testing.T) {
	hub := t.TempDir()

	var stdout bytes.Buffer
	if err := Run(hub, "3.52.0", nil, &stdout); err != nil {
		t.Fatalf("status: %v", err)
	}

	if !strings.Contains(stdout.String(), "No packages installed") {
		t.Error("expected 'No packages installed'")
	}
}

func TestRunCommandRegistry(t *testing.T) {
	hub := t.TempDir()
	// Create a minimal vendor dir so we get past the package section.
	writeManifest(t, hub, "cnos.core", `{
		"schema": "cn.package.v1",
		"name": "cnos.core",
		"version": "3.52.0",
		"kind": "package",
		"engines": {"cnos": "3.52.0"}
	}`)

	commands := []CommandInfo{
		{Name: "help", Summary: "Show available commands", Tier: "kernel"},
		{Name: "deploy", Summary: "Deploy to prod", Tier: "repo-local"},
		{Name: "daily", Summary: "Daily reflection", Tier: "package", Package: "cnos.core"},
	}

	var stdout bytes.Buffer
	if err := Run(hub, "3.52.0", commands, &stdout); err != nil {
		t.Fatalf("status: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "Commands:") {
		t.Error("expected Commands: header")
	}
	if !strings.Contains(out, "kernel:") {
		t.Error("expected kernel: tier group")
	}
	if !strings.Contains(out, "repo-local:") {
		t.Error("expected repo-local: tier group")
	}
	if !strings.Contains(out, "package:") {
		t.Error("expected package: tier group")
	}
	if !strings.Contains(out, "daily") {
		t.Error("expected daily command listed")
	}
	if !strings.Contains(out, "(cnos.core)") {
		t.Error("expected package attribution for daily command")
	}
}

func TestRunSkipsDirsWithoutManifest(t *testing.T) {
	hub := t.TempDir()
	// Create a dir without cn.package.json — should be silently skipped.
	os.MkdirAll(filepath.Join(hub, ".cn", "vendor", "packages", "junk"), 0755)
	writeManifest(t, hub, "cnos.core", `{
		"schema": "cn.package.v1",
		"name": "cnos.core",
		"version": "3.52.0",
		"kind": "package",
		"engines": {"cnos": "3.52.0"}
	}`)

	var stdout bytes.Buffer
	if err := Run(hub, "3.52.0", nil, &stdout); err != nil {
		t.Fatalf("status: %v", err)
	}

	out := stdout.String()
	if strings.Contains(out, "junk") {
		t.Error("junk dir without manifest should not appear")
	}
	if !strings.Contains(out, "cnos.core") {
		t.Error("expected cnos.core in output")
	}
}
