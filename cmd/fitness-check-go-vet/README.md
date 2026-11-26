# go-vet

Runs `go vet ./...` in every Go module of the repo, catching suspicious code the compiler accepts. It is on by default.

## Behavior

[Requirement 0001](docs/requirements/0001-go-code-passes-vet.md)

## Example

```
fitness go-vet
```

It reads no `.fitnessrc.json` key.

