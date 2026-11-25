### 2026.11.25.1640

> The requirement suite runs as twelve parallel jobs, so a pull request gets its answer in under two minutes.

- Build: CI runs the root requirement tests as twelve balanced jobs that fail past two minutes, and the changed checks' folders in one more.
- Build: only main runs the slow release build and measure script cases.
- Test: the 0001 cases share the installed bundle's check programs instead of building every check.
- Docs: requirement 0001.2 chains a step for how CI splits the suite, which its test proves.
