### 2026.11.14.1800

> The checks move here from fitness-runner commit `eb0488a`, and each one builds as its own program.

- Build: `scripts/migrate-checks-from-engine.sh` copies the checks, their helpers, docs and tests from one fitness-runner commit and rewrites their imports.
- Docs: the README says where each check lives, how the migration runs and how to test; GitHub folds the word lists.
- Chore: `local-branch` stays in fitness-runner for now, since it judges only with the engine's branch metadata; requirement 0150 stays with its own check.
- Build: each check is a `package main` of this module, importing fitness-runner's public packages at the version `go.mod` pins, and `--describe` adds protocol 1.
- Test: `tests/requirements.sh` runs the moved requirement tests through the pinned engine, which starts this repo's check programs in place of its built-in ones.
