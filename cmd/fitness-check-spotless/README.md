# spotless

Checks Java formatting with the [Spotless](https://github.com/diffplug/spotless) plugin the repo's own build declares. The Java language module owns the build run and its judgment. It is on by default.

## Behavior

[Requirement 0001](docs/requirements/0001-java-formatting-passes-spotless.md)

## Notes

It runs `spotless:check` for Maven and `spotlessCheck` for Gradle. The [Java module](https://github.com/may-journal/fitness-runner/blob/main/go/pkg/lang/java/java.go) says what a build root is and which wrapper or tool runs it.

## Formatter

The build's own Spotless block picks the formatter, so the check never imposes a style. The hint suggests Google Java Format.

## Example

```
fitness spotless
```

It reads no `.fitnessrc.json` key.

