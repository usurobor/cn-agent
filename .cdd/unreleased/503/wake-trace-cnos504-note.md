# cnos#504 cross-reference note

The wake-trace surface established by cnos#503 is the intended diagnostic
starting point for stale-claim recovery (cnos#504).

Trace location for cell-dispatch wakes: `.cdd/unreleased/{issue}/wake-trace.md`

The `claude_session_id` and `execution_artifact` fields in the trace file
point to the Claude session and the uploaded execution artifact (workflow
artifact, 14-day retention). These are the primary diagnostic inputs for
cnos#504 stale-claim recovery.

This note satisfies AC5 of cnos#503 (documentation-only AC).
