# adr-checks

Pairs each architecture decision, an `ADR` issue, with the check in `fitness/adr-NNNN/` that proves it. An ADR is prose, so a decision holds only while a check proves it. The check's README links the issue by its full URL, and adr-as-issue keeps the decision itself out of files. A requirement doc holds the same way, while its tests pass.

## Behavior

[Requirement 0001](docs/requirements/0001-each-decision-has-its-check.md)

## Example

```
fitness adr-checks
```

It reads no `.fitnessrc.json` key.

