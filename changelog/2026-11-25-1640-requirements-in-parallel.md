### 2026.11.25.1640

> The requirement suite runs as twelve parallel jobs, so a pull request gets its answer in under two minutes.

- Build: `tests/requirements.sh` runs one share of the root requirement tests, or the changed checks' own folders, as `REQUIREMENT_SHARD` picks.
- Build: CI runs the root tests as twelve balanced jobs that fail past two minutes.
- Build: only main runs the measure script's slow test.
- Test: cases share the installed bundle's check programs instead of each building every check.
- Docs: requirement 0001.2 chains a step for how CI splits the suite.
- Test: each root test runs in exactly one share, and only the changed checks' folders run.
