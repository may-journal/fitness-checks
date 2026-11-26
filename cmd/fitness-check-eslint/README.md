# eslint

Runs ESLint over the repo with the shared flat config, by executing the real [ESLint](https://eslint.org) binary. It resolves the binary from `node_modules/.bin`, walking up from the repo root, then PATH — never npx. When no eslint binary resolves, or ESLint crashes or fails without a finding, the check reports info naming the cause. The shared flat config `eslint.config.mjs` is forced via `--config`.

The config resolves in order: the repo's own file, an installed `@mayjournal/fitness-shared` package under `node_modules`, then the copy embedded in the check binary. That embedded copy is materialized on demand.

The shared config imports its plugins (typescript-eslint, jsdoc, perfectionist, prettier) as bare specifiers. Those packages must be installed alongside eslint in the consumer repo — the check's peer contract. Linting is a syntactic parse only, with no type-aware rules, so the runner-default timeout stands.

## Behavior

[Requirement 0001](docs/requirements/0001-javascript-and-typescript-pass-eslint.md)

## Notes

- It enforces `sort-keys` (natural ascending) and the other rules from `eslint.config.mjs`.
- It lints `.ts`, `.tsx`, `.cjs`, `.js` and `.mjs` files, declarations and tests included.
- A repo that chose [Biome](https://github.com/may-journal/fitness-runner/blob/main/fitness/adr-0006/README.md) resolves its binary like eslint, and a missing binary fails with an install hint.

## Example

```
fitness eslint
```

It reads no `.fitnessrc.json` key.

