### 2026.11.25.1230

> A requirement doc can state its acceptances as one all-caps Gherkin block, which the Cucumber parser reads.

- Feat: `requirements` reads a fenced `gherkin` block of SCENARIO acceptances with the Cucumber parser and pairs steps with tests as for bullets.
- Feat: it fails a keyword not in caps, an AND not indented under its step and a comma outside code.
- Feat: migration `0003.6-gherkin` runs up through `--fix` or `--up` and back through `--down`, and pairing reads the move as a reformat.
- Docs: the requirements check's own requirement 0003 states the form, its README shows it, and root 0005.6 reads Gherkin Givens.
