# architecture-contract

Validates declarations beside all five architecture views before more specific checks resolve their proof. It reads one model from the diagrams and callout tables already in each view, and adds no fields or columns. Later checks share the model through `internal/architecture.Load` and requirement clauses through `internal/requirement.Load`.

## Behavior

[Requirement 0001](docs/requirements/0001-architecture-views-form-one-valid-model.md)

## Example

```
fitness architecture-contract
```

It reads no `.fitnessrc.json` key.

