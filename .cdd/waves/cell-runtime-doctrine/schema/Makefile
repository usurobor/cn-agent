# cnos#671 R9 — structural pre-authorization validation runner (Go + CUE repo; NO Python).
#
# Run from the schema/ dir:  make -C .cdd/waves/cell-runtime-doctrine/schema
#
#   make clean        -> cue vet the real plan artifacts + clean/canonical bases; all PASS.
#                        The 6 real child contracts validate against #WorkingCellContract (the named
#                        WC refinement: class=working, requested_output.kind=artifact). The 4
#                        canonical variants (nullable scope, optional prior_receipt, relation_graph,
#                        judgment) validate against the faithful canonical #CellContract.
#   make regressions  -> cue vet every bad fixture; each MUST be rejected (non-zero). Fails if any passes.
#   make all          -> clean + regressions.
#
# Requires `cue` on PATH (cue v0.17+). Install: go install cuelang.org/go/cmd/cue@latest
CUE ?= cue
SCHEMA := ./
WAVEDIR := ..

.PHONY: all clean regressions
all: clean regressions

# ---- STRUCTURAL pre-authorization validation of the real plan artifacts (must PASS) ----
clean:
	$(CUE) vet $(SCHEMA) $(WAVEDIR)/wave.cn-wave-v1.yaml -d '#Wave'
	# The 6 real child contracts are Working Cells -> validated against the named WC refinement.
	$(CUE) vet $(SCHEMA) $(WAVEDIR)/contracts/wc-1.cn-cell-contract-v1.yaml -d '#WorkingCellContract'
	$(CUE) vet $(SCHEMA) $(WAVEDIR)/contracts/wc-2.cn-cell-contract-v1.yaml -d '#WorkingCellContract'
	$(CUE) vet $(SCHEMA) $(WAVEDIR)/contracts/wc-3a.cn-cell-contract-v1.yaml -d '#WorkingCellContract'
	$(CUE) vet $(SCHEMA) $(WAVEDIR)/contracts/wc-3b.cn-cell-contract-v1.yaml -d '#WorkingCellContract'
	$(CUE) vet $(SCHEMA) $(WAVEDIR)/contracts/wc-4.cn-cell-contract-v1.yaml -d '#WorkingCellContract'
	$(CUE) vet $(SCHEMA) $(WAVEDIR)/contracts/wc-5.cn-cell-contract-v1.yaml -d '#WorkingCellContract'
	$(CUE) vet $(SCHEMA) $(WAVEDIR)/oracle-registry.yaml -d '#AssuranceRegistry'
	$(CUE) vet $(SCHEMA) $(WAVEDIR)/intent.cn-intent-v1.yaml -d '#Intent'
	$(CUE) vet $(SCHEMA) ./regressions/wave.clean.yaml -d '#Wave'
	$(CUE) vet $(SCHEMA) ./regressions/contract.clean.yaml -d '#CellContract'
	# Canonical §2 variants (the FAITHFUL #CellContract, not the WC refinement) — all MUST PASS.
	$(CUE) vet $(SCHEMA) ./regressions/contract.canonical-nullable-scope.yaml -d '#CellContract'
	$(CUE) vet $(SCHEMA) ./regressions/contract.canonical-optional-prior-receipt.yaml -d '#CellContract'
	$(CUE) vet $(SCHEMA) ./regressions/contract.canonical-relation-graph.yaml -d '#CellContract'
	$(CUE) vet $(SCHEMA) ./regressions/contract.canonical-judgment.yaml -d '#CellContract'
	@echo "clean: all artifacts PASS cue vet"

# ---- REGRESSIONS: each blocker/§2-drift mutation MUST be rejected by cue vet (non-zero). ----
regressions:
	@rc=0; \
	for f in ./regressions/wave.bad-*.yaml; do \
	  if $(CUE) vet $(SCHEMA) "$$f" -d '#Wave' >/dev/null 2>&1; then \
	    echo "FAIL (unexpected pass): $$f"; rc=1; \
	  else echo "ok rejected: $$f"; fi; \
	done; \
	for f in ./regressions/contract.bad-*.yaml; do \
	  if $(CUE) vet $(SCHEMA) "$$f" -d '#CellContract' >/dev/null 2>&1; then \
	    echo "FAIL (unexpected pass): $$f"; rc=1; \
	  else echo "ok rejected: $$f"; fi; \
	done; \
	if [ $$rc -ne 0 ]; then echo "regressions: a bad fixture was NOT rejected"; exit 1; fi; \
	echo "regressions: every bad fixture rejected by cue vet"
