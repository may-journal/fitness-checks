### 2026.11.25.1230

> A requirement doc can number its acceptances, each a name over numbered GIVEN, WHEN and THEN items in caps.

- Feat: `requirements` sets each numbered acceptance's ID from its place, reads it with the Cucumber parser and pairs its steps.
- Feat: it fails a dash, a misnumbered or misplaced item, a lowercase keyword, AND or BUT, a comma, or a bad name.
- Feat: a `Removed in <version>: <reason>` line deletes an ID in place, counts toward no limit, and history fails a move or reuse.
- Feat: migration `0003.2-gherkin` numbers bullets through `--fix` or `--up`, and `--down` turns them back.
- Docs: the check's requirement 0003 states the form, its README shows it, and root 0005.6 reads numbered Givens.
