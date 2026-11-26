# go-complexity

Enforces a cyclomatic complexity ceiling on Go functions, the Go-native counterpart of the house eslint rule (`complexity: max 5`).

## Behavior

[Requirement 0001](docs/requirements/0001-go-functions-stay-simple.md) states what this check does, and a test case named for each step proves it.

## Ceiling

A repo can override the ceiling in `.fitnessrc.json`. Run `fitness explain go-complexity` for the key and default:

```
{
  "goComplexity": { "max": 8 }
}
```

A native Go engine (`go/ast` and `go/parser` from the standard library) scores the code; no external tool is required.

## Notes

- A file it cannot read is skipped and named as info, and the other files are still judged.
- `filesChecked` counts every Go file scanned, tests included.
