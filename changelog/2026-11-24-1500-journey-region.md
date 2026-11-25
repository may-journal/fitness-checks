### 2026.11.24.1500

> The user journey names each requirement doc from the doc itself, in generated regions.

- Feat: `architecture-journey` fails a marked journey region that differs from the docs its marker lists, and names `fitness architecture-journey --fix`.
- Feat: a start marker such as `%% journey-nodes:start 0001 0002` lists the docs its region holds; a doc in no region or two fails.
- Feat: `--fix` rewrites every region, and writes nothing when a doc lacks its title.
- Docs: `05-user-journey.md` generates its eight nodes inside the hand-written phases, which keep their edges. `requirement.Document` carries its `Title`.
- Test: requirement 0005.6 proves a stale region, an unlisted doc and a twice-listed doc fail, a fresh region passes, and `--fix` rewrites every region.
