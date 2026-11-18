### 2026.11.18.0400

> The 88 requirement docs become seven, 0001 to 0007, one for each thing a person gets from these checks.

- Docs: seven requirement docs replace the old ones, each with at most eight acceptances.
- Test: every old test is now a named case in its new acceptance's table test, one file per requirement.
- Test: shared helpers drop old numbers, and repeated run code moves into `stagedCheck` and `forbiddenMajor`.
- Docs: journey labels and check docs point at the new requirement numbers.
- Refactor: the mermaid checks share their run and JavaScript whitespace helpers through `internal/diagrams`.
