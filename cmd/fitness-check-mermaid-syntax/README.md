# mermaid-syntax

Fails a mermaid diagram that would not render, so a broken diagram never reaches a reader. It is part of the mermaid diagram and callout table fitness-function set ([#28](https://github.com/may-journal/fitness-runner/issues/28)).

## Behavior

[Requirement 0001](docs/requirements/0001-diagrams-render.md)

## Example

```
fitness mermaid-syntax
fitness mermaid-syntax --body-file body.md
```

It reads no `.fitnessrc.json` key.

