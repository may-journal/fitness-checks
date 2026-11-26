# mermaid-level-bleed

Warns when a callout description at one C4 level repeats the previous level verbatim. Lower levels should add level-specific rationale, not restate the parent. It is part of the mermaid diagram and callout table fitness-function set ([#28](https://github.com/may-journal/fitness-runner/issues/28)).

## Behavior

[Requirement 0001](docs/requirements/0001-lower-levels-add-new-detail.md)

## Example

```
fitness mermaid-level-bleed
```

It reads no `.fitnessrc.json` key.

