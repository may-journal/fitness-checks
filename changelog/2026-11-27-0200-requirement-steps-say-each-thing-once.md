### 2026.11.27.0200

> A requirement's steps say each thing once, so the check catches a repeat instead of a reviewer.

- Feat: `requirements` fails sibling WHENs that open with the same three words that carry meaning, naming the words to move into the GIVEN.
- Feat: `requirements` fails a THEN that matches an earlier THEN of its acceptance, naming both lines and asking to group the WHENs.
- Docs: `requirements` gains requirement 0009, with a case for each step. Its own docs and every check's docs group their WHENs to pass.
