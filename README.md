# fitness-checks

The fitness checks that the [fitness-runner](https://github.com/may-journal/fitness-runner) engine runs. Each check is its own program in `cmd/fitness-check-<name>`, with a short README beside it. It imports the engine's public packages at the version `go.mod` pins, and no other check.

The engine, hooks, workflows and installer stay in fitness-runner. Its bundle keeps its own copy of each check until the engine finds these programs. [Split fitness-runner into an engine repo and a checks repo](https://github.com/may-journal/fitness-runner/issues/918) plans that step.

## License

[MIT](LICENSE)

## Migration

[scripts/migrate-checks-from-engine.sh](scripts/migrate-checks-from-engine.sh) migrates the checks from the fitness-runner engine into this monorepo, run once per migration. It takes a fitness-runner clone and a commit, copies that commit's check folders, helpers, requirement docs and tests, and rewrites their imports. History stays in fitness-runner.

## Releases

Each check releases on its own line, tagged `fitness-check-<name>-vX.Y.Z`, with its programs for linux and darwin attached. On each push to main, [scripts/release-plan.sh](scripts/release-plan.sh) names the checks the push can alter and picks each one's version from its conventional commits. [scripts/release-publish.sh](scripts/release-publish.sh) then builds and publishes each one, stopping when [scripts/describe-patterns](scripts/describe-patterns/main.go) finds a describe pattern that does not compile. `go.mod` and `go.sum` hold every dependency version.

### Release index

Each check's major tag, `fitness-check-<name>-vN`, follows its newest release of that major, which [scripts/release-major.sh](scripts/release-major.sh) moves after each run. The release `fitness-check-index` holds `index.json`, at `https://github.com/may-journal/fitness-checks/releases/download/fitness-check-index/index.json`, which [scripts/release-index.sh](scripts/release-index.sh) writes from the published releases. A check is its folder name without `fitness-check-`. An artifact is a tar of one program, keyed `<os>-<arch>` for linux and darwin on amd64 and arm64.

```json
{
  "format": 1,
  "repository": "may-journal/fitness-checks",
  "checks": {
    "gofmt": {
      "latest": "1.2.0",
      "majors": { "1": "1.2.0" },
      "versions": {
        "1.2.0": {
          "tag": "fitness-check-gofmt-v1.2.0",
          "files": "\\.(go)$",
          "exclude": "(^|/)(testdata|vendor)/",
          "artifacts": {
            "linux-arm64": { "url": "https://github.com/.../fitness-check-gofmt-v1.2.0-linux-arm64.tar", "sha256": "0f3a...", "gz": { "url": "https://github.com/.../fitness-check-gofmt-v1.2.0-linux-arm64.tar.gz", "sha256": "9c1d..." } }
          }
        }
      }
    }
  }
}
```

An installer reads `majors` for a major, takes `versions[version].artifacts[<os>-<arch>]`, and checks the download against its `sha256`. `format` rises only for a change that breaks a reader. An artifact may also hold `gz`, the gzip of its tar with its own `url` and `sha256`. An installer that knows it may fetch that instead.

A version may also hold `files`, `exclude`, `needs` and `requires`, copied from its `<tag>.describe.json` release asset. They say which files reach the check and what else it requires, as fitness-runner's `checkkit.Describe` defines them. An engine can then fetch only the checks a run reaches. A version released before them holds none, which an engine reads as a check every file reaches.

## Tests

`go test ./...` builds every package and runs the tests that need no installed engine. [tests/requirements.sh](tests/requirements.sh) runs the requirement tests in `requirements/` as a user meets each check. It installs the fitness-runner engine `go.mod` pins, with this repo's check programs beside it. On a pull request, CI builds and tests only the checks [scripts/changed-checks.sh](scripts/changed-checks.sh) names, then runs that suite.

Each check's full docs are its `README.md`, which GitHub renders and `fitness-check-<name> --describe` summarizes. A check not yet moved keeps them in its `doc.go`, which [pkg.go.dev](https://pkg.go.dev/github.com/may-journal/fitness-checks) renders, and its README only links them.
