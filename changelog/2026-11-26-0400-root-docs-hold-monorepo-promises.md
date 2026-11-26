### 2026.11.26.0400

> The root requirement docs keep only monorepo promises, and the requirements check owns the rules that stayed in them.

- Docs: root requirements 0003 to 0007 are gone, and 0002.5 and 0002.7 left, each moved to its check.
- Docs: `requirements` gains packages and root scope acceptances, and docs 0006 and 0007. `requirements-changed` gains doc 0002.
- Docs: the user journey shows the root docs 0001, 0002 and 0008.
- Test: the cases move out of the root tests to those checks, and root tests a check already proved are deleted.
- Fix: `requirements` fails a monorepo root acceptance that names a single `fitness-check-` package.
