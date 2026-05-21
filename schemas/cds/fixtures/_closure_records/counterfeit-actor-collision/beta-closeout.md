# β Closeout (counterfeit fixture for #389 AC8 / C1)

Actor: same-actor@example

This stub deliberately shares an actor with the sibling `alpha-closeout.md`
to exercise V's counterfeit.actor_separation predicate. The receipt that
references both files omits `mode: collapsed`, so V emits FAIL.
