### 2026.11.26.2000

> Old-style requirements fail wherever they hide, in a README or in a step that joins clauses.

- Feat: `readme-behavior` fails a GIVEN, WHEN or THEN list item in any case outside a code fence, naming the line.
- Docs: `readme-behavior` gains 0001.6 and `requirements` gains 0008. `dependency-currency` gains 0001.4, which holds the behavior its README had hidden in Notes.
- Docs: every requirement doc here that the new rule failed splits its joined steps into numbered GIVEN and THEN items.
- Feat: `requirements` fails a Gherkin block keyword such as `Scenario:` not in caps, 0008.2.
- Feat: `requirements` fails GIVEN, WHEN, THEN, AND or BUT inside a step not in caps outside code and quotes, 0008.1.
