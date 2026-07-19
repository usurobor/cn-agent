# R4 Working Contract and Fresh Review Request — cnos#669

## Control-plane binding

- Issue: `#669`
- PR: `#670`
- Exact reviewed R3 head / R4 implementation base:
  `bc6f6aa946cbbb555bb409c8bc180f5d49685b6e`
- External β R3 verdict: **ITERATE**
- β comment ID: `5016263226`
- β comment URL:
  `https://github.com/usurobor/cnos/pull/670#issuecomment-5016263226`

R4 is one focused successor commit. It does not amend the reviewed R3 SHA.
The earlier R2 convergence/γ closeout and the R3 review are historical evidence
only; none applies as a verdict over R4 matter.

## Canonical output and scope

This remains a WC whose canonical output is the repository-native recursive-
cell measurement package. R4 is limited to β's two required corrections and
the directly coupled self-coherence repairs:

1. validate `next_fixes` as a list of exact `{axis, fix}` objects, with
   `axis ∈ {alpha,beta,gamma}` and a nonblank string `fix`;
2. require a fresh run root and atomically create an immutable `emission/`
   bundle under an exclusive emission lock;
3. snapshot that emission, the six external responses, and invariant assessment
   into publication staging; validate/hash/use only those snapshots; bind them
   all in the typed success marker with the six reports and aggregate; then
   publish once by an atomic rename under an exclusive publication lock;
4. refuse existing, symlinked, or concurrently locked canonical output; leave
   no canonical or staging output after mid-target or final-CUE failure;
5. bind the invariant assessment to the exact six response SHA-256 values as
   well as the six prompt digests;
6. carry each canonical TSC report's `bottleneck_axis` into its level record and
   use that axis for the selected lowest-`C_sigma` level;
7. include `SKILL.md` and the instruction assembly script in the path-framed CM
   authority bundle; and
8. prove the positive, malformed, reuse, contention, symlink, mid-target, and
   final-validation paths with deterministic fixtures, plus compatibility with
   the real pinned coh executable.

Out of scope remains unchanged: State-B arbitrary methodology/target loading,
installed-package activation, a new semantic sample, independent calibration,
standing promotion, #662 ratification, merge, or any FSM/control-plane
transition.

## Acceptance and review request

R4 is reviewable only as a new immutable commit whose focused diff is based on
the exact R3 SHA above; all local and shallow-clone checks pass; unrelated #429
material remains untouched; the draft PR body names R4's fresh-review state;
and α posts an immutable receipt naming the final R4 SHA and evidence.

After α stops writing, request a fresh context-isolated β review of that exact
R4 SHA. Only a new β convergence may permit a subsequent γ binding and CC
re-adjudication. Shared OpenAI Codex hosting and GitHub identity must remain
disclosed; functional separation is not hosting independence or off-house
standing.
