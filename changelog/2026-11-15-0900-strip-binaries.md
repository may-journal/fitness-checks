### 2026.11.15.0900

> Release programs are stripped and path-trimmed, with a gzip beside each tar, and CI fails when one grows.

- Feat: `scripts/release-build.sh` builds with path trimming and `-s -w`, and makes a gzip of each tar beside it.
- Feat: the index gets an optional `gz` field, and `url`, `sha256` and `format` 1 stay as the shipped installer reads them.
- Ci: `binary-size.yml` fails when a program is over 10% off `tests/size-baseline.tsv`, built with the go.mod toolchain.
- Test: requirement 0167 proves the stripping, the size gate and a format 1 reader, and the baseline fits the pinned engine.
- Docs: `scripts/measure-checks.sh` records sizes and suite times, and the `local-branch` README is regenerated.
