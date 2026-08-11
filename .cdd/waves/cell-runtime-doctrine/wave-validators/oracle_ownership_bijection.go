// cnos#671 R12 — WAVE-BOUNDARY pre-authorization validator (Go + CUE repo; NO Python).
//
// Materializes the whole-wave oracle-ownership / classification BIJECTION named by
// oracle-registry.yaml `wave_predicates.wave_oracle_ownership_bijection_enforced`
// (deferred_owner: "wave", wave_authorization_gated: true).
//
// A PRE-AUTHORIZATION gate runs BEFORE any Working Cell executes, so it cannot be
// deferred to a WC — it must be a real, runnable validator in the plan matter. This
// program is that validator: credential-free, single-file `package main`, standard
// library ONLY (os/exec + encoding/json; no go.mod required, `go run` works anywhere).
//
// R12 SOUNDNESS REPAIR (external-β R11 BLOCKER): the R11 validator hand-parsed YAML by
// indentation/prefix and silently dropped CUE-valid flow-style values (`predicates:
// ["x"]`, flow maps), producing a demonstrated complete-wave FALSE PASS. R12 removes
// ALL ad-hoc YAML parsing: every input is normalized by `cue export --out json` (the
// same CUE that structurally vets this wave) and decoded via encoding/json, so flow and
// block serializations are indistinguishable to this program. Owners are derived from
// SEMANTIC identity (contract `cell.id`, cross-checked against the wave's node ids),
// never from filename text. Any parse/normalization loss, or an empty semantic owner or
// predicate, FAILS CLOSED (exit 2). Changing the serialization style can no longer hide
// an unregistered predicate from the gate.
//
// It proves, over the six child contracts + the authoritative registry:
//
//	(1) BIJECTION — union(child acceptance.predicates) over (owner, predicate) ⇄ the
//	    registry `assurance:` entries over (owner, predicate), EXACTLY: no missing
//	    (a child predicate with no registry entry), no phantom (a registry entry with
//	    no child predicate), no duplicate on either side. (78 ⇄ 78 for the real wave.)
//	(2) OWNERSHIP — every `mechanically-verifiable` predicate binds EXACTLY ONE owner:
//	    exactly one concrete checker|schema path, and its (owner, predicate) is unique
//	    (a predicate "owned twice" is rejected).
//	(3) IDENTITY — every child owner is a real wave node id (contract cell.id ∈ wave
//	    nodes; the contract filename stem must equal cell.id). A contract mis-filed or a
//	    cell.id that is not an owning wave node FAILS CLOSED.
//
// Prints a clear result; EXIT 0 IFF bijective, 1 if a bijection/ownership defect is
// found, 2 on load / normalization / fail-closed error.
//
// Requires `cue` on PATH (the wave already requires it for structural `cue vet`); the
// binary name can be overridden with the CUE env var.
//
// Usage:
//
//	go run oracle_ownership_bijection.go <path>
//	  <path> = a WAVE DIRECTORY  (contains contracts/*.cn-cell-contract-v1.yaml +
//	           wave.cn-wave-v1.yaml + oracle-registry.yaml) — the real wave; OR
//	  <path> = a self-contained FIXTURE FILE (top-level `child_predicates:` +
//	           `assurance:` blocks; optional `wave_nodes:` to exercise the identity
//	           check) — a minimal registry+predicate set.
//
// Examples:
//
//	go run oracle_ownership_bijection.go ..                                  # real wave (from wave-validators/)
//	go run oracle_ownership_bijection.go fixtures/oracle-ownership.one-checker-each.positive.yaml
//	go run oracle_ownership_bijection.go fixtures/oracle-ownership.double-owned.negative.yaml
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// A (owner, predicate) key — the bijection is over these pairs (the same predicate
// STRING may legitimately appear under different owners, so ownership is per-pair).
type pair struct {
	owner     string
	predicate string
}

func (p pair) String() string { return p.owner + " :: " + p.predicate }

// An assurance/registry entry (only the fields the bijection needs).
type entry struct {
	owner          string
	predicate      string
	classification string
	checker        string
	schema         string
}

type model struct {
	mode      string // "wave-dir" | "fixture-file"
	child     []pair
	assurance []entry
}

// ---- CUE-normalized input (NO ad-hoc YAML parsing) ----

// cueExportJSON runs `cue export <path> --out json` and decodes the normalized JSON into
// v. This is the ONLY input path: CUE collapses flow/block serialization to one JSON
// shape, so a flow-style list/map can no longer be silently dropped. A non-zero cue exit
// or a decode error FAILS CLOSED (the caller exits 2).
func cueExportJSON(path string, v any) error {
	cueBin := os.Getenv("CUE")
	if cueBin == "" {
		cueBin = "cue"
	}
	cmd := exec.Command(cueBin, "export", path, "--out", "json")
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cue export %s failed (fail-closed): %v: %s", path, err, strings.TrimSpace(errb.String()))
	}
	if err := json.Unmarshal(out.Bytes(), v); err != nil {
		return fmt.Errorf("decoding cue-export JSON for %s (fail-closed): %w", path, err)
	}
	return nil
}

// Decode targets — the SEMANTIC shape (a subset of cn.cell.contract.v1 / cn.wave.v1 /
// the registry). Fields absent in the source decode to zero values; we fail closed on
// any required semantic field that is empty.
type contractDoc struct {
	Cell struct {
		ID string `json:"id"`
	} `json:"cell"`
	Acceptance struct {
		Predicates []string `json:"predicates"`
	} `json:"acceptance"`
}

type waveDoc struct {
	Nodes []struct {
		ID string `json:"id"`
	} `json:"nodes"`
}

type entryJSON struct {
	Owner          string `json:"owner"`
	Predicate      string `json:"predicate"`
	Classification string `json:"classification"`
	Checker        string `json:"checker"`
	Schema         string `json:"schema"`
}

type registryDoc struct {
	Assurance []entryJSON `json:"assurance"`
}

type fixtureDoc struct {
	ChildPredicates []struct {
		Owner     string `json:"owner"`
		Predicate string `json:"predicate"`
	} `json:"child_predicates"`
	Assurance []entryJSON `json:"assurance"`
	WaveNodes []string    `json:"wave_nodes"`
}

func toEntries(js []entryJSON) []entry {
	var out []entry
	for _, e := range js {
		out = append(out, entry{
			owner:          e.Owner,
			predicate:      e.Predicate,
			classification: e.Classification,
			checker:        e.Checker,
			schema:         e.Schema,
		})
	}
	return out
}

// loadWaveDir reads the real wave: the six child contracts (owner = cell.id, cross-checked
// against filename stem AND the wave node set), and the registry assurance entries.
func loadWaveDir(dir string) (*model, error) {
	glob := filepath.Join(dir, "contracts", "wc-*.cn-cell-contract-v1.yaml")
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no child contracts matched %s", glob)
	}
	sort.Strings(files)

	// Wave node ids — the authoritative owner namespace (identity, not filename).
	var wave waveDoc
	if err := cueExportJSON(filepath.Join(dir, "wave.cn-wave-v1.yaml"), &wave); err != nil {
		return nil, err
	}
	nodeSet := map[string]bool{}
	for _, n := range wave.Nodes {
		if strings.TrimSpace(n.ID) == "" {
			return nil, fmt.Errorf("wave node with empty id (fail-closed)")
		}
		nodeSet[n.ID] = true
	}
	if len(nodeSet) == 0 {
		return nil, fmt.Errorf("wave has no nodes (fail-closed)")
	}

	m := &model{mode: "wave-dir"}
	for _, f := range files {
		var c contractDoc
		if err := cueExportJSON(f, &c); err != nil {
			return nil, err
		}
		owner := strings.TrimSpace(c.Cell.ID) // SEMANTIC owner — the contract's declared identity
		if owner == "" {
			return nil, fmt.Errorf("%s: empty cell.id (fail-closed)", f)
		}
		// filename stem must equal cell.id (a mis-filed contract is a fail-closed error).
		stem := strings.SplitN(filepath.Base(f), ".", 2)[0]
		if stem != owner {
			return nil, fmt.Errorf("%s: filename stem %q != cell.id %q (fail-closed)", f, stem, owner)
		}
		// cell.id must be a real owning wave node.
		if !nodeSet[owner] {
			return nil, fmt.Errorf("%s: cell.id %q is not a wave node id (fail-closed)", f, owner)
		}
		if len(c.Acceptance.Predicates) == 0 {
			return nil, fmt.Errorf("%s: no acceptance.predicates (fail-closed)", f)
		}
		for _, p := range c.Acceptance.Predicates {
			if strings.TrimSpace(p) == "" {
				return nil, fmt.Errorf("%s: empty acceptance predicate (fail-closed)", f)
			}
			m.child = append(m.child, pair{owner: owner, predicate: p})
		}
	}

	var reg registryDoc
	regPath := filepath.Join(dir, "oracle-registry.yaml")
	if err := cueExportJSON(regPath, &reg); err != nil {
		return nil, err
	}
	m.assurance = toEntries(reg.Assurance)
	if len(m.assurance) == 0 {
		return nil, fmt.Errorf("no assurance entries in %s (fail-closed)", regPath)
	}
	if err := requireNonEmpty(m.assurance, regPath); err != nil {
		return nil, err
	}
	return m, nil
}

// loadFixture reads a self-contained fixture file (child_predicates + assurance + optional
// wave_nodes). Because it too goes through cue export, flow-style fixtures normalize
// identically; empty owner/predicate fails closed. If wave_nodes is present, every child
// owner must be one of them (exercises the identity check without a full wave dir).
func loadFixture(path string) (*model, error) {
	var fx fixtureDoc
	if err := cueExportJSON(path, &fx); err != nil {
		return nil, err
	}
	m := &model{mode: "fixture-file"}

	var nodeSet map[string]bool
	if len(fx.WaveNodes) > 0 {
		nodeSet = map[string]bool{}
		for _, id := range fx.WaveNodes {
			if strings.TrimSpace(id) == "" {
				return nil, fmt.Errorf("%s: empty wave_nodes id (fail-closed)", path)
			}
			nodeSet[id] = true
		}
	}
	for i, cp := range fx.ChildPredicates {
		owner := strings.TrimSpace(cp.Owner)
		predicate := strings.TrimSpace(cp.Predicate)
		if owner == "" || predicate == "" {
			return nil, fmt.Errorf("%s: child_predicates[%d] has empty owner or predicate (fail-closed)", path, i)
		}
		if nodeSet != nil && !nodeSet[owner] {
			return nil, fmt.Errorf("%s: child_predicates[%d] owner %q is not in wave_nodes (fail-closed)", path, i, owner)
		}
		m.child = append(m.child, pair{owner: owner, predicate: predicate})
	}
	m.assurance = toEntries(fx.Assurance)
	if len(m.child) == 0 && len(m.assurance) == 0 {
		return nil, fmt.Errorf("fixture %s has neither child_predicates nor assurance blocks (fail-closed)", path)
	}
	if err := requireNonEmpty(m.assurance, path); err != nil {
		return nil, err
	}
	return m, nil
}

// requireNonEmpty fails closed if any assurance entry has an empty owner or predicate —
// a normalization loss must never be read as an empty-string slot that collapses to a
// spurious match.
func requireNonEmpty(entries []entry, src string) error {
	for i, e := range entries {
		if strings.TrimSpace(e.owner) == "" || strings.TrimSpace(e.predicate) == "" {
			return fmt.Errorf("%s: assurance[%d] has empty owner or predicate (fail-closed)", src, i)
		}
	}
	return nil
}

func load(path string) (*model, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return loadWaveDir(path)
	}
	return loadFixture(path)
}

func sortedPairs(m map[pair]bool) []pair {
	var ps []pair
	for p := range m {
		ps = append(ps, p)
	}
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].owner != ps[j].owner {
			return ps[i].owner < ps[j].owner
		}
		return ps[i].predicate < ps[j].predicate
	})
	return ps
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run oracle_ownership_bijection.go <wave-dir|fixture-file>")
		os.Exit(2)
	}
	path := os.Args[1]
	m, err := load(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load error: %v\n", err)
		os.Exit(2)
	}

	// Child side: (owner, predicate) multiset.
	childCount := map[pair]int{}
	for _, p := range m.child {
		childCount[p]++
	}
	// Registry side: (owner, predicate) multiset + classification/checker lookup.
	assrCount := map[pair]int{}
	mechCount := map[pair]int{}    // mechanically-verifiable entries per pair
	mechConcrete := map[pair]int{} // # of concrete checker|schema bound per mech pair
	for _, e := range m.assurance {
		k := pair{owner: e.owner, predicate: e.predicate}
		assrCount[k]++
		if e.classification == "mechanically-verifiable" {
			mechCount[k]++
			n := 0
			if strings.TrimSpace(e.checker) != "" {
				n++
			}
			if strings.TrimSpace(e.schema) != "" {
				n++
			}
			mechConcrete[k] += n
		}
	}

	childSet := map[pair]bool{}
	for p := range childCount {
		childSet[p] = true
	}
	assrSet := map[pair]bool{}
	for p := range assrCount {
		assrSet[p] = true
	}

	// Bijection defects.
	missing := map[pair]bool{}     // child -> no registry entry
	phantom := map[pair]bool{}     // registry -> no child predicate
	childDup := map[pair]bool{}    // duplicate child (owner,predicate)
	doubleOwned := map[pair]bool{} // duplicate registry entry — a predicate "owned twice"
	for p, c := range childCount {
		if c > 1 {
			childDup[p] = true
		}
		if !assrSet[p] {
			missing[p] = true
		}
	}
	for p, c := range assrCount {
		if c > 1 {
			doubleOwned[p] = true
		}
		if !childSet[p] {
			phantom[p] = true
		}
	}
	// Ownership: every mechanically-verifiable predicate binds EXACTLY ONE checker|schema.
	mechBadOwner := map[pair]bool{}
	for p, c := range mechCount {
		// c>1 already flagged by doubleOwned; the concrete count must be exactly one PER entry.
		if mechConcrete[p] != c {
			mechBadOwner[p] = true
		}
	}

	bijective := len(missing) == 0 && len(phantom) == 0 && len(childDup) == 0 &&
		len(doubleOwned) == 0 && len(mechBadOwner) == 0

	// ---- report ----
	fmt.Println("wave-oracle-ownership-bijection validator (cnos#671 R12, wave-boundary pre-authorization; cue-export normalized)")
	fmt.Printf("  input:                        %s  (mode: %s)\n", path, m.mode)
	fmt.Printf("  child acceptance predicates:  %d\n", len(m.child))
	fmt.Printf("  registry assurance entries:   %d\n", len(m.assurance))
	fmt.Printf("  mechanically-verifiable:      %d\n", len(mechCount))
	fmt.Printf("  missing   (child -> registry): %d\n", len(missing))
	fmt.Printf("  phantom   (registry -> child): %d\n", len(phantom))
	fmt.Printf("  child duplicates:              %d\n", len(childDup))
	fmt.Printf("  double-owned (registry dup):   %d\n", len(doubleOwned))
	fmt.Printf("  mech missing/≠1 checker owner: %d\n", len(mechBadOwner))

	report := func(label string, s map[pair]bool) {
		if len(s) == 0 {
			return
		}
		fmt.Printf("  --- %s ---\n", label)
		for _, p := range sortedPairs(s) {
			fmt.Printf("      %s\n", p)
		}
	}
	report("MISSING (child acceptance predicate with no registry entry)", missing)
	report("PHANTOM (registry entry with no child acceptance predicate)", phantom)
	report("DUPLICATE child (owner,predicate)", childDup)
	report("DOUBLE-OWNED registry (owner,predicate) — a predicate owned twice", doubleOwned)
	report("mechanically-verifiable NOT bound to exactly one checker|schema", mechBadOwner)

	fmt.Printf("  bijective: %v\n", bijective)
	if bijective {
		fmt.Println("RESULT: PASS — oracle-ownership / classification bijection holds (exit 0)")
		os.Exit(0)
	}
	fmt.Println("RESULT: FAIL — oracle-ownership / classification bijection broken (exit 1)")
	os.Exit(1)
}
