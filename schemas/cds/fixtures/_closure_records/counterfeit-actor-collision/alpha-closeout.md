# α Closeout (counterfeit fixture for #389 AC8 / C1)

Actor: same-actor@example

This stub exists to exercise V's counterfeit.actor_separation predicate. The
companion `beta-closeout.md` declares the same actor. V's rule C1 should
flag this as a failed_predicate (since the receipt does NOT declare
`mode: collapsed`).
