# cnos#671 R11 — WAVE-BOUNDARY pre-authorization gate runner (Go + CUE repo; NO Python).
#
# This is the credential-free invocation that BINDS wave authorization to the PASS of the
# wave-oracle-ownership-bijection validator. A PRE-AUTHORIZATION gate runs BEFORE any Working
# Cell executes, so it cannot be deferred to a WC — it is materialized and run here.
#
# Run from this dir:  make -C .cdd/waves/cell-runtime-doctrine/wave-validators
#
#   make bijection  -> run the validator over the REAL six contracts + registry (../); MUST exit 0
#                      (bijective: true — 78 child acceptance predicates ⇄ 78 registry entries).
#   make fixtures   -> positive fixture MUST exit 0; negative (double-owned) fixture MUST exit non-zero.
#   make all        -> bijection + fixtures. This is the authorization-ready gate: `make all` exits 0
#                      IFF the real-wave bijection holds AND both fixtures behave — i.e. the validator
#                      exists, is uncorrupted, and PASSes. Remove/corrupt the validator and this fails.
#
# Requires `go` on PATH (go 1.21+). No modules, no network, no credentials.
GO      ?= go
VALIDATOR := oracle_ownership_bijection.go
WAVEDIR   := ..
POS := fixtures/oracle-ownership.one-checker-each.positive.yaml
NEG := fixtures/oracle-ownership.double-owned.negative.yaml

.PHONY: all bijection fixtures
all: bijection fixtures
	@echo "wave-validators: PASS — wave-boundary oracle-ownership bijection holds; authorization-ready gate satisfied"

# ---- the real-wave pre-authorization bijection (MUST PASS) ----
bijection:
	$(GO) run $(VALIDATOR) $(WAVEDIR)

# ---- fixtures: positive PASSes (exit 0); negative double-own is REJECTED (non-zero) ----
fixtures:
	$(GO) run $(VALIDATOR) $(POS)
	@if $(GO) run $(VALIDATOR) $(NEG) >/dev/null 2>&1; then \
	  echo "FAIL (unexpected pass): negative fixture $(NEG) was accepted"; exit 1; \
	else echo "ok rejected (non-zero): $(NEG)"; fi
