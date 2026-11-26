# checkstyle

Lints Java with the [Checkstyle](https://checkstyle.org) plugin the repo's own build declares, and holds its config to a complexity ceiling. The Java language module owns the build run, the ceiling rule and the judgment. It is on by default.

## Behavior

[Requirement 0001](docs/requirements/0001-java-code-passes-checkstyle.md)

## Notes

- A failure that names no finding reports the build's own error lines.
- A Java repo with no running JDK fails with an install hint.
- Maven runs with the repo's only config, or else with `config/checkstyle/checkstyle.xml`, unless the build names its own.

## Example

```
fitness checkstyle
```

It reads no `.fitnessrc.json` key.

