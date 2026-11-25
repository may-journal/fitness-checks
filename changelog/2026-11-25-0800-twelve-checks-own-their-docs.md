### 2026.11.25.0800

> Twelve more checks, `adr-as-issue` to `changelog`, own their requirement docs beside their code, as the `swiftlint` pilot does.

- Docs: `adr-as-issue`, `adr-checks`, the six `architecture-*` checks, `build-output-untracked` and the three `changelog` checks state their behavior as numbered steps of their own requirement 0001, which each README links.
- Docs: root requirements 0003, 0004, 0005, 0006 and 0007 keep only what remains once those checks' acceptances moved out.
- Build: `.fitnessrc.json` names each moved check's folder as a package root.
- Test: each step has a named case in its check's own test, with the shared architecture and local check fixtures moved to `internal/acceptance`.
