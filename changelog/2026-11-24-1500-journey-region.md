### 2026.11.24.1500

> The user journey names each requirement doc from the doc itself, in a generated region.

- Feat: `architecture-journey` fails a marked journey region that differs from the docs' IDs and titles, and names `fitness architecture-journey --fix`.
- Feat: `--fix` rewrites the region, and writes nothing when a doc lacks its title.
- Feat: a journey with no markers stays hand-written and passes.
- Docs: `05-user-journey.md` generates its eight nodes. `requirement.Document` carries its `Title`.
- Test: requirement 0005.6 proves a stale region fails, a fresh one passes, and `--fix` rewrites it.
