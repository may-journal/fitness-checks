### 2026.11.14.2100

> Each check releases on its own line, and a push releases exactly the checks it can alter.

- Feat: `release.yml` runs `scripts/release-plan.sh` on each push to main, which tags each check the push can alter, `fitness-check-<name>-vX.Y.Z`, bumped by its conventional commits.
- Feat: `scripts/release-publish.sh` builds each tagged check for four platforms, unstripped and uncompressed, and publishes it with a checksum and notes.
- Feat: `scripts/changed-checks.sh --each` names the checks each commit can alter, so CI and release share one decision.
- Chore: Release Please and its config, manifest and `version.txt` leave this repo, since it releases one line per check.
- Test: requirement 0162 proves a check change, a helper change and a `go.mod` change each release exactly the checks they alter.
