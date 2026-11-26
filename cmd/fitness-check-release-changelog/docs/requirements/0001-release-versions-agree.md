# 0001 Release Versions Agree

## Why

I can trust that the version Release Please records is the version my repo ships.

## Measurement

repos whose manifest and version file agree

-

repos with a Release Please manifest

## Requirements

1. Release versions agree
   1. GIVEN Developer has a repo that may hold a Release Please manifest
   2. GIVEN it may hold a `version.txt`
   3. GIVEN it may hold a `CHANGELOG.md`
      1. WHEN the repo has no manifest
         1. THEN `release-changelog` passes with zero files
      2. WHEN the manifest is an empty JSON object
         1. THEN it passes as a bootstrap
      3. WHEN the manifest is malformed or lacks a numeric root version
         1. THEN it fails naming the problem
      4. WHEN `version.txt` differs from the manifest
         1. THEN it fails naming both versions
      5. WHEN `version.txt` is missing
         1. THEN it fails naming the version it must match
      6. WHEN the versions match
         1. THEN it passes with or without a `CHANGELOG.md`
