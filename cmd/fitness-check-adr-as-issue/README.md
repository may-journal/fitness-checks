# adr-as-issue

Keeps architecture decisions in `ADR` issues, never in files. An issue can be discussed, linked and superseded in one place, while a file drifts from the talk that made it. The `ADR` issue template and issue-check judge each issue's shape.

## Behavior

[Requirement 0001](docs/requirements/0001-decisions-live-in-adr-issues.md)

## Notes

Run `fitness explain adr-as-issue` for the allow key that lets an ADR file wait.

## Example

```
fitness adr-as-issue
```

Run `fitness explain adr-as-issue` for what each key does. This sets each key to its default:

```json
{
  "adrAsIssue": {
    "allow": []
  }
}
```

