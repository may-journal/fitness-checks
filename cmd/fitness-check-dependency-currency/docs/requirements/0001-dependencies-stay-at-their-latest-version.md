# 0001 Dependencies Stay At Their Latest Version

## Why

I see every declared npm dependency that has an installable update, so my project stays at peak fitness.

## Measurement

repos whose declared dependencies are all at their latest version

-

repos with a `package.json`

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a repo with or without a root `package.json` and workspaces
  WHEN every declared dependency is at its latest published version
  THEN it passes
  WHEN a dependency of any kind is behind its latest version
  THEN it fails naming each as `name: current → latest`
  WHEN a declared dependency is not installed
  THEN it fails as `name: missing → latest`
  WHEN a dependency is an internal `@mayjournal/*` package
  THEN it is skipped
  WHEN a workspace declares a stale dependency
  THEN it fails naming it
  WHEN workspaces share one install
  THEN identical lines collapse to one
  WHEN the registry is unreachable or answers an error or garbage
  THEN it reports info counting the packages it could not query and no finding
  WHEN the repo has no `package.json`
  THEN it passes with no files
```
