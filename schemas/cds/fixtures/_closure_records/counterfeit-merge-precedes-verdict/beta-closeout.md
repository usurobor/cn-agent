# β Closeout (counterfeit fixture for #389 AC8 / C2)

Actor: beta-actor@example

This β closeout's git commit time is fresh (cycle/389 commit), but the
receipt that references it declares `boundary_decision.decided_at`
deep in the past (2024). V's rule C2 compares the two timestamps and
flags the β-postdates-δ violation.
