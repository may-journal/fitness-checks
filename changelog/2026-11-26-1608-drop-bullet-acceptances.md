### 2026.11.26.1608

> `requirements` reads only numbered acceptances, so a doc still in the bullet form fails until it moves.

- Feat: a bullet acceptance fails at its line and points to the numbered form in the `requirements` README.
- Refactor: the bullet grammar leaves `internal/requirement`, while history still counts bullet-era IDs as held.
- Docs: `requirements` gains 0004.6, rewords 0004.1.3 and 0004.4.2 and removes 0004.3. Root 0002.8 has one check fail a bullet doc.
- Fix: pairing compares two numbered versions of a doc as written but for acceptance names, as `requirements` 0006.1.7 states.
- Test: the shared fixtures write numbered acceptances, and `architecture-actors` 0001.1.7 covers bullet docs.
