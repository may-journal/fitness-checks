# prettier

Runs Prettier `--check` over the repo's files, to ensure they are formatted. It uses a local Prettier config when present, such as `.prettierrc.json`, `prettier.config.cjs`, or a `package.json` `"prettier"` field. Otherwise it falls back to the shared `prettier.config.cjs`. That comes from an installed `@mayjournal/fitness-shared` package when present, else the copy embedded in the check binary, materialized on demand.

[Prettier](https://prettier.io) itself is a peer tool: the check execs the real binary, never npx. It resolves from `node_modules/.bin` (walking up from the repo root) then PATH. A missing binary reports info with a one-line install hint, and a crashed one info quoting its output.

Run `fitness prettier` from your app root to check formatting. Pass through args to Prettier: `fitness prettier --write .` or `fitness prettier --write src/`. The runner forwards args after the check name to the check, which runs at the repo root.

## Behavior

[Requirement 0001](docs/requirements/0001-files-are-formatted-with-the-repos-formatter.md)

## Notes

Errors list file paths that need formatting.

## Example

```
fitness prettier
```

It reads no `.fitnessrc.json` key.

