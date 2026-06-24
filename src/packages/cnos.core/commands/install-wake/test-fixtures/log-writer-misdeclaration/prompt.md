# test-log-writer-misdeclaration fixture prompt

Synthetic minimal prompt for the cycle/496 AC4 negative-case smoke
(activation_log_writer mis-declaration refusal). Never rendered to a
substrate workflow; exists solely so the renderer can resolve the
prompt_template path during its validation walk.

The renderer must refuse this fixture at exit code 4 BEFORE reaching
the prompt-inlining stage, so the body content is immaterial — only
the file's presence at the manifest-declared path matters.
