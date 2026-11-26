# dependency-currency

Fails a declared npm dependency behind its latest published version — an installable update means the project isn't at peak fitness. It runs in every repo and passes clean without a `package.json`.

## Behavior

[Requirement 0001](docs/requirements/0001-dependencies-stay-at-their-latest-version.md)

## Notes

- `filesChecked` counts the tracked `package.json` manifests.
- Reports against `latest`, not `wanted`: a dependency inside its declared semver range but behind the newest release is still flagged.
- The point is to stay current, not merely in-range.
- Only direct declared dependencies count, since you can only bump what you declare.
- The registry URL comes from the nearest `.npmrc` `registry=` entry, repo root first and then home, and defaults to the public registry.
- One proxy blocks GitHub and the registry alike.

## Example

```
fitness dependency-currency
```

It reads no `.fitnessrc.json` key.
