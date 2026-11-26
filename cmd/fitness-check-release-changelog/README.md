# release-changelog

Binds Release Please's manifest to the version file a reviewer can inspect. The check applies when `.release-please-manifest.json` exists and contains a root release version.

## Behavior

[Requirement 0001](docs/requirements/0001-release-versions-agree.md) states what this check does, and a test case named for each step proves it.

Release notes come from [entry files](../fitness-check-changelog/README.md#entry-files), not `CHANGELOG.md`, so this check asks nothing of `CHANGELOG.md`. Set `"skip-changelog": true` in `release-please-config.json` so Release Please stops writing it.
