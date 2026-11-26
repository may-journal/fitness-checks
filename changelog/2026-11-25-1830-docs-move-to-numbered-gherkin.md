### 2026.11.25.1830

> This repo's own requirement docs now read as numbered Gherkin lists, and the requirements template shows that form.

- Docs: every requirement doc, root and package, holds numbered acceptances with caps keywords, each with its own short name.
- Docs: each comma and chained "or" is split into the next numbered GIVEN or THEN.
- Refactor: remove migration `0003.2-gherkin`, `--up`, `--down` and the `--fix` conversion; any other flag now fails naming it.
- Docs: the requirements README template shows the numbered form in place of bullets.
