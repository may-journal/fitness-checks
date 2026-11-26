### 2026.11.25.0800

> Twelve more checks, `adr-as-issue` to `changelog`, own their requirement docs beside their code, as the `swiftlint` pilot does.

- Docs: twelve checks state their behavior as numbered steps of their own requirement 0001, which each README links.
- Docs: root requirements 0003 to 0007 keep only what remains, and the level-repeat rule moves to the unused 0005.9.
- Build: `.fitnessrc.json` names each moved check's folder as a package root.
- Test: each step has a named case in its check's own test, with shared fixtures moved to `internal/acceptance`.
