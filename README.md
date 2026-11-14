# fitness-checks

The fitness checks that the [fitness-runner](https://github.com/may-journal/fitness-runner) engine runs. Each check is its own program in `cmd/fitness-check-<name>`, with its README beside it. It imports the engine's public packages at the version `go.mod` pins, and no other check.

The engine, hooks, workflows and installer stay in fitness-runner. Its bundle keeps its own copy of each check until the engine finds these programs. [Split fitness-runner into an engine repo and a checks repo](https://github.com/may-journal/fitness-runner/issues/918) plans that step.

## License

[MIT](LICENSE)

## Migration

[scripts/migrate-checks-from-engine.sh](scripts/migrate-checks-from-engine.sh) migrates the checks from the fitness-runner engine into this monorepo, run once per migration. It takes a fitness-runner clone and a commit, copies that commit's check folders, helpers, requirement docs and tests, and rewrites their imports. History stays in fitness-runner.

## Tests

`go test ./...` builds every package and runs the tests that need no installed engine. [tests/requirements.sh](tests/requirements.sh) runs the requirement tests in `requirements/` as a user meets each check. It installs the fitness-runner engine `go.mod` pins, with this repo's check programs beside it. On a pull request, CI builds and tests only the checks [scripts/changed-checks.sh](scripts/changed-checks.sh) names, then runs that suite.
