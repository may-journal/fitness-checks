# vitest-coverage-exclude

Fails any Vitest coverage exclude entry, so coverage measures every file.

## Behavior

[Requirement 0001](docs/requirements/0001-coverage-measures-every-file.md)

## Notes

Vitest's own default exclude skips tests and declarations, so set `coverage.exclude: []` to measure them; the shared config does.

It reads `vitest.config.*` or the `package.json` `vitest` field in the project root. If none exist, it uses the shared vitest config. That is an installed `@mayjournal/fitness-shared` package when present, else the copy embedded in the check binary, materialized on demand.

It reads the config as text, so exclude entries built from variables, spreads, or imports are invisible to it.

## Example

```
fitness vitest-coverage-exclude
```

It reads no `.fitnessrc.json` key.

