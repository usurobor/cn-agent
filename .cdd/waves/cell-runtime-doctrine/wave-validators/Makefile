# cnos#671 R12 — WAVE-BOUNDARY pre-authorization gate runner (Go + CUE repo; NO Python).
#
# This is the credential-free invocation that BINDS wave authorization to the PASS of the
# wave-oracle-ownership-bijection validator. A PRE-AUTHORIZATION gate runs BEFORE any Working
# Cell executes, so it cannot be deferred to a WC — it is materialized and run here.
#
# R12: the validator normalizes every input through `cue export --out json` (no ad-hoc YAML
# parsing), so `cue` is now required on PATH alongside `go` (the wave already requires cue for
# structural `cue vet`). The gate BUILDS the validator once and runs the BINARY, so the EXACT
# program exit code is asserted per fixture (0 = bijective PASS; 1 = bijection/ownership defect;
# 2 = fail-closed load/normalization error) — `go run` collapses every non-zero to 1 and cannot.
#
# Run from this dir:  make -C .cdd/waves/cell-runtime-doctrine/wave-validators
#
#   make bijection  -> run the validator over the REAL six contracts + registry (../); MUST exit 0
#                      (bijective: true — 78 child acceptance predicates ⇄ 78 registry entries).
#   make fixtures   -> every fixture is asserted to exit with its EXACT expected code:
#                      positives 0; bijection-defect negatives 1; fail-closed negatives 2.
#   make all        -> bijection + fixtures. This is the authorization-ready gate: `make all` exits 0
#                      IFF the real-wave bijection holds AND every fixture exits its expected code —
#                      i.e. the validator exists, is uncorrupted, and is SOUND over flow/block YAML.
#                      Remove/corrupt the validator and this fails.
#   make clean      -> remove the build dir.
#
# Requires `go` (1.21+) and `cue` (v0.17+) on PATH. No modules, no network, no credentials.
GO       ?= go
VALIDATOR := oracle_ownership_bijection.go
BUILD     := .build
BIN       := $(BUILD)/oracle_ownership_bijection
WAVEDIR   := ..

.PHONY: all bijection fixtures clean

all: bijection fixtures
	@echo "wave-validators: PASS — wave-boundary oracle-ownership bijection holds; authorization-ready gate satisfied"

$(BIN): $(VALIDATOR)
	@mkdir -p $(BUILD)
	$(GO) build -o $(BIN) ./$(VALIDATOR)

# ---- the real-wave pre-authorization bijection (MUST exit 0) ----
bijection: $(BIN)
	@echo "== real-wave oracle-ownership bijection (must exit 0: bijective true, 78 <-> 78) =="
	$(BIN) $(WAVEDIR)

# ---- fixtures: each asserted to exit its EXACT expected code (0 positive | 1 defect | 2 fail-closed) ----
fixtures: $(BIN)
	@set -e; \
	expect() { \
	  if "$(BIN)" "$$1" >/dev/null 2>&1; then rc=0; else rc=$$?; fi; \
	  if [ "$$rc" != "$$2" ]; then echo "FAIL: $$1 exited $$rc, expected $$2"; exit 1; fi; \
	  echo "ok (exit $$rc == $$2): $$1"; \
	}; \
	echo "== POSITIVES (expect exit 0) =="; \
	expect fixtures/oracle-ownership.one-checker-each.positive.yaml 0; \
	expect fixtures/oracle-ownership.flow-style.positive.yaml 0; \
	expect fixtures/wavedir.positive 0; \
	echo "== BIJECTION-DEFECT NEGATIVES (expect exit 1) =="; \
	expect fixtures/oracle-ownership.double-owned.negative.yaml 1; \
	expect fixtures/oracle-ownership.flow-missing-registry.negative.yaml 1; \
	expect fixtures/oracle-ownership.flow-mismatched-owner.negative.yaml 1; \
	expect fixtures/wavedir-flow-unregistered.negative 1; \
	echo "== FAIL-CLOSED NEGATIVES (expect exit 2) =="; \
	expect fixtures/oracle-ownership.empty-owner.negative.yaml 2; \
	expect fixtures/oracle-ownership.cellid-not-a-node.negative.yaml 2

clean:
	@rm -rf $(BUILD)
