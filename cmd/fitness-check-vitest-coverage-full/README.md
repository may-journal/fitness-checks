# vitest-coverage-full

Requires Vitest coverage thresholds of 100 in the repo and the shared config, then runs `vitest run --coverage`.

## Behavior

[Requirement 0001](docs/requirements/0001-coverage-is-full-and-tests-pass.md) states what this check does, and a test case named for each step proves it.

A pass needs two things. The project and the fitness-runner shared config set Vitest coverage thresholds of 100 for branches, functions, lines, and statements. Then `vitest run --coverage` exits 0 in the project root. Both configs are read statically, through their code views, so a commented-out threshold never counts.

- Threshold and exclude settings come from the project's vitest config, else the shared `vitest.config.mjs` (installed `@mayjournal/fitness-shared` package or the binary's embedded copy).
- When the project has no local vitest config, `vitest run --coverage` runs with `--config` pointing at that resolved file.
- The vitest binary resolves from `node_modules/.bin` (walking up from root), then PATH — never npx. A missing binary fails with an install hint.

## Errors

```
Vitest coverage thresholds must be 100 for branches, functions, lines, and statements.
Fitness-runner package must have Vitest coverage thresholds set to 100 for branches, functions, lines, and statements.
```

A failing coverage run reports the last line of Vitest output when Vitest judged the code. It judged the code when a threshold is missed or a test failed. A run that exits non-zero with neither, could not start, or was killed reports info quoting its last lines. A missing coverage provider is one example.
