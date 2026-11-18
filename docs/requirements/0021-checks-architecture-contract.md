# 0021 Architecture Contract

## Why

I read one model from the diagrams and callout tables I already keep, so checks can share the same facts.

## Measurement

architecture views read or rejected with clear source locations
-
architecture views checked

## Requirements

- 0021.11
    - Given Developer has no product code, only tests in any test folder, and no architecture
        - When I check its contract
            - Then it passes with zero files
- 0021.2
    - Given Developer has complete diagrams in either root
        - When I renumber their callouts
            - Then the same diagram IDs pass
- 0021.10
    - Given Developer has a partial architecture, or Go, Swift, TypeScript, JavaScript or Java product code and no architecture
        - When I check it
            - Then it fails asking for each missing view
- 0021.4
    - Given Developer has duplicate, invalid, or unresolved diagram IDs or actor names
        - When I check them
            - Then errors name their diagram lines
- 0021.5
    - Given Developer has two elements naming one source path
        - When I check ownership
            - Then errors name both claims
- 0021.6
    - Given Developer has a source link the check cannot read
        - When I check evidence
            - Then errors name its row and source
- 0021.7
    - Given Developer has invalid requirement links
        - When I check the contract
            - Then errors name only the broken reference
- 0021.9
    - Given Developer has unchanged broken evidence, container drift, or component drift
        - When I check unrelated changes
            - Then the same findings remain
