### 2026.11.14.2300

> An installer resolves a check at a major version through a moving tag and one index.

- Feat: `scripts/release-major.sh` points `fitness-check-<name>-vN` at the newest release of major N, and it also backfills the tags for the 1.0.0 releases.
- Feat: `scripts/release-index.sh` writes `index.json` and `scripts/release-index-publish.sh` keeps it on the release `fitness-check-index`.
- Fix: the release scripts read only `vX.Y.Z` tags as a check's last release, so a major tag is never mistaken for one.
- Docs: requirement 0163 records the index format.
- Test: requirement 0163 proves the major tag moves once in any clone, the plan ignores it, and the index lists published releases.
