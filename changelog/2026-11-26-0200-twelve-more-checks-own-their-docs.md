### 2026.11.26.0200

> Twelve more checks, `go-test` to `markdown-links`, own their requirement docs beside their code.

- Docs: each check states its behavior as numbered acceptances in its own requirement 0001, which its README links.
- Docs: root requirements 0003, 0004, 0006 and 0002 keep only what remains, with each moved acceptance left as a Removed line.
- Build: `.fitnessrc.json` names each moved check's folder as a package root.
- Test: each step has a named case in its check's own test, with shared fixtures in `internal/acceptance`.
