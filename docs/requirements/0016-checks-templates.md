# 0016 Checks Templates

## Why

Anyone can open any requirement and find its reason, its measure, and its proof in the same place. Every rule for one area sits in one doc.

## Measurement

requirement docs with no shape errors
-
requirement docs

## Requirements

- 0016.1
    - Given Developer has a doc whose file name and title do not carry the same ID
        - When fitness runs
            - Then it fails
- 0016.2
    - Given Developer has two docs with the same ID
        - When fitness runs
            - Then it fails
- 0016.3
    - Given Developer has a doc whose sections are not Why, Measurement, and Requirements, each with content
        - When fitness runs
            - Then it fails
- 0016.4
    - Given Developer has a Measurement that is not a numerator, a `-` line, and a denominator
        - When fitness runs
            - Then it fails
- 0016.8
    - Given Developer has an acceptance ID that is repeated or names another requirement
        - When fitness runs
            - Then it fails
- 0016.9
    - Given Developer has an acceptance with a second Given, When, or Then, or a chain out of order
        - When fitness runs
            - Then it fails
- 0016.10
    - Given Developer has two docs with acceptances sharing a When and similar Givens
        - When fitness runs
            - Then it fails naming both docs and linking ADR 0004 in the README of its check
- 0016.11
    - Given Developer has a doc whose Measurement ends in a Source line citing code that moved
        - When I run `fitness requirements --fix` and commit
            - Then the check fails naming the line, and once the fix deletes it, passes with no test change
