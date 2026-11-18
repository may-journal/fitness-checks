# 0013 Checks Binary Size

## Why

Each check ships as its own program, so a smaller program installs faster. Releases strip debug info, trim paths and compress, and CI fails when a program grows past its baseline.

## Measurement

Check programs within their size baseline
-
Check programs

## Requirements

- 0013.1
    - Given Developer builds a release for one check
        - When the program in its gzip is compared with the same build without `-s -w`
            - Then the program in the gzip is smaller
- 0013.2
    - Given Developer has a size baseline, and a program grown, shrunk or gone from it by over 10%
        - When `scripts/check-binary-sizes.sh` runs
            - Then it exits 1 naming that check, and with the current sizes it exits 0
- 0013.3
    - Given Developer has a repo to judge
        - When `scripts/measure-checks.sh` runs
            - Then `sizes.tsv` holds each check's sizes and start-up time, and `timing.tsv` holds the suite time both ways
- 0013.4
    - Given Developer has releases with a tar only, and with a tar and its gzip
        - When `scripts/release-index.sh` reads them, and a format 1 reader looks up each artifact
            - Then the reader finds the tar's `url` and `sha256` in each, and `gz` appears only where a gzip was published
