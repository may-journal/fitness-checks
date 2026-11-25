### 2026.11.25.1640

> The requirement suite runs as eight parallel jobs, so a pull request gets its answer in under two minutes.

- Build: `tests/requirements.sh` runs one share of the root requirement tests, or the changed checks' own folders, as `REQUIREMENT_SHARD` picks.
- Build: the Checks workflow runs the root tests as eight jobs balanced by `tests/requirement-times.tsv`, failing a suite step past two minutes, and the measure script's test only on main.
- Docs: requirement 0001.9 states how CI splits the suite.
- Test: each root test runs in exactly one share, and only the changed checks' folders run.
