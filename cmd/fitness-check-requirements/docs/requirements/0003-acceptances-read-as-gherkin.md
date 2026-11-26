# 0003 Acceptances Read As Gherkin

## Why

I number each acceptance and give it a name. Under it I write numbered GIVEN, WHEN and THEN items in caps, one condition or result each. The standard Gherkin parser reads them.

## Measurement

Numbered acceptances the check reads and pairs with their tests

-

Numbered acceptances

## Requirements

- 0003.1
  - Given Developer has numbered acceptances `requirements` reads
    1. When each level numbers from 1 with keywords in caps
       - Then it pairs steps with tests
    2. When a keyword is not in caps
       - Then it names the keyword
    3. When an item is an AND or BUT
       - Then it asks for the next GIVEN or THEN
    4. When an item holds a comma outside code
       - Then it asks to split it
    5. When an item is misnumbered or misnamed or misplaced
       - Then it names what it needs
    6. When the parser rejects an item
       - Then it fails at that line
- 0003.2
  - Given Developer runs migration `0003.2-gherkin`
    1. When `--fix` or `--up` runs on bullets
       - Then each acceptance is named for its first WHEN
    2. When `--down` runs
       - Then extra GIVENs and THENs join with and
    3. When a flag names no migration
       - Then it lists the migrations
- 0003.3
  - Given Developer moves a bullet doc to numbered items
    - When `requirements` judges it
      - Then its tests need no change
- 0003.4
  - Given Developer has committed numbered acceptances
    1. When one is renamed in place
       - Then it passes
    2. When a commit moves one
       - Then it names both IDs
    3. When it moves back
       - Then it passes
- 0003.5
  - Given Developer has committed numbered acceptances
    1. When a commit removes one with its test
       - Then it passes
    2. When a `Removed` line holds items or lacks a version or reason
       - Then it fails
    3. When a test names a `Removed` acceptance
       - Then it fails
    4. When it is restored
       - Then it fails as a reused ID
    5. When `--up` meets a skipped ID
       - Then it writes a versioned `Removed` line
    6. When `--down` meets a `Removed` line
       - Then the bullets skip its ID
    7. When one commit removes one and appends it
       - Then it passes
- 0003.6
  - Given Developer has a doc at the acceptance cap with `Removed` lines
    - When `requirements` and `prose-budget` read it
      - Then both pass
