# 0170 Checks Requirements Renumber Map

## Why

I can renumber requirement docs and still trace each moved ID to its new home.

## Measurement

renumber map entries whose new ID a doc defines

-

renumber map entries

## Requirements

- 0170.1
  - Given Developer has folded a doc into another and listed `old -> new` in docs/requirements/renumber.txt
    - When a new doc reuses the old doc's number and acceptances
      - Then the check passes, while the same reuse without the entry fails
- 0170.2
  - Given Developer has listed one moved acceptance, but not its removed doc, in the renumber map
    - When a new doc reuses both IDs
      - Then only the unlisted doc ID fails as deleted before
- 0170.3
  - Given Developer has a renumber map line that is malformed, repeats or maps to itself, or points at an ID no doc defines
    - When fitness runs the requirements check
      - Then it fails naming the map's line
- 0170.4
  - Given Developer has moved a doc to a new number and listed it in the renumber map
    - When I stage the doc without its tests
      - Then it fails until the tests move with it
- 0170.5
  - Given Developer has reused a number the renumber map moved, and a later commit deleted that new doc
    - When another doc brings the number back
      - Then it fails as deleted before
