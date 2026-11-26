# changelog

Judges the `changelog/` entry files, the only changelog source. Each must hold one dated `### yyyy.mm.dd.HHMM` heading and its bullets, and its name must start with that heading's date. `CHANGELOG.md` is not judged, so a repo may keep, freeze, or drop it, and an upstream fork's `### Added` history passes as it is. The package version is not judged either.

## Behavior

[Requirement 0001](docs/requirements/0001-changelog-entry-files-are-dated.md) states what this check does, and a test case named for each step proves it.

## Entry files

A pull request adds its entry as its own `changelog/<yyyy-mm-dd-HHMM>-<change>.md` file, with the change in kebab-case. Two open pull requests never edit the same lines that way. It holds one `### yyyy.mm.dd.HHMM` heading and its bullets. The name's date is that heading's, dots read as hyphens, so the folder lists in release order.
