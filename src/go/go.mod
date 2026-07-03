module github.com/usurobor/cnos/src/go

go 1.24

// cnos#392: cdd-verify is co-located with the cnos.cdd package per the
// package-scoping pin in the issue's implementation contract. The
// workspace file (../../go.work) and this replace directive both ensure
// the local module is used; the replace makes `go mod tidy` work without
// network resolution.
require github.com/usurobor/cnos/packages/cnos.cdd/commands/cdd-verify v0.0.0-00010101000000-000000000000

replace github.com/usurobor/cnos/packages/cnos.cdd/commands/cdd-verify => ../packages/cnos.cdd/commands/cdd-verify

// cnos#556: issues-map (the `cn issues map` domain implementation) is
// co-located with the cnos.issues package, mirroring the cnos#392
// cdd-verify precedent above exactly. The workspace file (../../go.work)
// and this replace directive both ensure the local module is used; the
// replace makes `go mod tidy` work without network resolution.
require github.com/usurobor/cnos/packages/cnos.issues/commands/issues-map v0.0.0-00010101000000-000000000000

replace github.com/usurobor/cnos/packages/cnos.issues/commands/issues-map => ../packages/cnos.issues/commands/issues-map

// cnos#568: issues-fsm (the `cn issues fsm evaluate` domain implementation)
// is co-located with the cnos.issues package, mirroring the cnos#556
// issues-map precedent above exactly. The workspace file (../../go.work)
// and this replace directive both ensure the local module is used; the
// replace makes `go mod tidy` work without network resolution.
require github.com/usurobor/cnos/packages/cnos.issues/commands/issues-fsm v0.0.0-00010101000000-000000000000

replace github.com/usurobor/cnos/packages/cnos.issues/commands/issues-fsm => ../packages/cnos.issues/commands/issues-fsm
