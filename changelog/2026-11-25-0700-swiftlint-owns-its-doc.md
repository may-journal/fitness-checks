### 2026.11.25.0700

> The `swiftlint` check owns its requirement doc beside its code, as the pilot for every check.

- Docs: `swiftlint` states its behavior as the seven steps of its own requirement 0001.1, which its README links.
- Docs: requirement 0006.6 keeps only the Java tools.
- Build: `.fitnessrc.json` names `cmd/fitness-check-swiftlint` as a package root.
- Fix: `prose-budget` lets a package's own requirement doc chain seven steps, as a root doc may.
- Test: each step has a named case in the check's own test, using stand-ins moved to `internal/acceptance`.
