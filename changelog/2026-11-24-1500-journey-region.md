### 2026.11.24.1500

> The user journey names each requirement doc from the doc itself, in a generated region.

- Feat: `architecture-journey` fails a marked journey region that differs from the docs' IDs and titles, and names `fitness architecture-journey --fix`.
- Feat: `--fix` rewrites the region. A journey with no markers stays hand-written and passes.
- Refactor: `requirement.Document` carries its `Title`.
- Docs: `05-user-journey.md` generates its eight nodes; edges and callouts stay hand-written.
- Test: requirement 0005.6 proves a stale region fails, a fresh one passes, and `--fix` rewrites it.
