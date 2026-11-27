# 0001 Dependencies Stay At Their Latest Version

## Why

I see every declared npm dependency that has an installable update, so my project stays at peak fitness.

## Measurement

repos whose declared dependencies are all at their latest version

-

repos with a `package.json`

## Requirements

1. Current dependencies pass
   1. GIVEN Developer has a repo with or without a root `package.json`
   2. GIVEN it may hold workspaces
      1. WHEN every declared dependency is at its latest published version
         1. THEN it passes
      2. WHEN a dependency is an internal `@mayjournal/*` package
         1. THEN it is skipped
      3. WHEN the repo has no `package.json`
         1. THEN it passes with no files
2. Stale dependencies fail
   1. GIVEN Developer has a repo with a root `package.json`
   2. GIVEN it has an install
      1. WHEN a dependency of any kind is behind its latest version
         1. THEN it fails naming each as `name: current → latest`
      2. WHEN a declared dependency is not installed
         1. THEN it fails as `name: missing → latest`
      3. WHEN a workspace declares a stale dependency
         1. THEN it fails naming it
      4. WHEN workspaces share one install
         1. THEN identical lines collapse to one
3. A silent registry only informs
   1. GIVEN Developer has a repo with a stale dependency
      1. WHEN the registry is unreachable or answers an error or garbage
         1. THEN it reports info counting the packages it could not query
         2. THEN it reports no finding
      2. WHEN the registry accepts a connection and never answers
         1. THEN it closes the connection within its request timeout
4. A blocked GitHub reads no registry
   1. GIVEN Developer has a repo with a stale dependency behind a proxy
      1. WHEN the runner finds GitHub's API host unreachable
         1. THEN it reads no registry
         2. THEN it reports info naming the host
         3. THEN it reports info naming the proxy
