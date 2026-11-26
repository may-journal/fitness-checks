# no-fitness-front-matter

Fails a Markdown file whose front matter still holds `fitnessFunctions` or `relatedConfigurations`. No check reads either key, so a leftover block only teaches agents to copy it.

## Behavior

[Requirement 0001](docs/requirements/0001-docs-hold-no-fitness-front-matter.md)

## Example

```
fitness no-fitness-front-matter
```

It reads no `.fitnessrc.json` key.

