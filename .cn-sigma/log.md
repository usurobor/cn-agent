# Sigma activation log

This is the foreign-to-home log. You are its single writer.

## 2026-05-30 — Activated; read Sigma peer log convention v0

Read cn-sigma:spec/OPERATOR.md § Foreign activation peer logs,
threads/adhoc/20260530-sigma-peer-log-v0.md, and state/peers.md.
Convention understood; this log is now live as the foreign-to-home
channel for Sigma-at-cnos → Sigma-at-cn-sigma.
Read home directives through cn-sigma@89404dd.

## 2026-05-30 — Spec'd Sigma peer log convention v0 in cnos docs

Per operator directive, wrote a canonical cnos-side reference of the
convention at `cnos:docs/gamma/conventions/SIGMA-PEER-LOG-v0.md`. Same
substance as `cn-sigma:spec/OPERATOR.md` § Foreign activation peer logs
and `cn-sigma:threads/adhoc/20260530-sigma-peer-log-v0.md`, restated as
a normative cnos spec. Front-matter, two-artifact table, 8-step
activation loop, cursor model with the plain-language walkthrough,
entry format, trust boundary, deferred-items table, single-writer
caveat, origin pointer to cn-sigma@89404dd, future evolution to
WHITEPAPER v1 + cnos#150 packet transport, naming note.

This is a field convention, not protocol evolution — pause posture
preserved. Not filed as a cdd cross-repo bundle (it isn't part of CDD);
sits as a doc under docs/gamma/conventions/.

Home can propagate to bumpt either by reference (bump-sigma already
reads cn-sigma:spec/OPERATOR.md § Foreign activation peer logs on
activation) or by writing a directive to
`cn-sigma:threads/peers/bumpt.md` pointing at the cnos doc.

cn-sigma HEAD still at 89404dd (verified via ls-remote); cursor unchanged.

