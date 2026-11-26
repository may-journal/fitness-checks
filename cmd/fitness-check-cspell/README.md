# cspell

Fails unknown words in the checked files. It spell-checks with a native Go engine and dictionaries embedded in the check binary — no external cspell tool, no npm install. The project `cspell.json` words always win over the embedded base dictionaries; it ignores no path.

## Behavior

[Requirement 0001](docs/requirements/0001-unknown-words-fail.md)

## Example

```
fitness cspell
fitness cspell --body-file body.md
```

It reads no `.fitnessrc.json` key.

