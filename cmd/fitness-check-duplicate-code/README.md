# duplicate-code

Fails a repo whose duplicated lines pass 1% of its scanned lines. A native Go clone-detection engine judges them: a generic lexer plus a rolling-hash detector in `internal/clonedetect`. Its thresholds are modeled on [jscpd](https://github.com/kucherenko/jscpd)'s, but no jscpd tool runs, and percentage parity with it is not the goal. It is on by default.

## Behavior

[Requirement 0001](docs/requirements/0001-duplicate-code-stays-under-one-percent.md) states what this check does, and a test case named for each step proves it.

## Escape hatch

For duplication that's a real, justified constraint, such as boilerplate a framework requires per type, wrap it in inline-comment markers. Don't force a refactor:

```
// duplicate-code:ignore-start
...
// duplicate-code:ignore-end
```

The older `jscpd:ignore-start` and `jscpd:ignore-end` markers still work as deprecated aliases, so repos that carry them keep passing.
