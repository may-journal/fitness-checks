# requirements-changed

Fails a change that adds or changes product code without also changing a test and a requirement doc. The [requirements](../fitness-check-requirements/README.md) check judges the docs and tests that exist; this check judges what a commit or pull request leaves out.

## Behavior

[Requirements](docs/requirements)

## Example

```
fitness requirements-changed
fitness requirements-changed --message "text"
```

It reads no `.fitnessrc.json` key.
