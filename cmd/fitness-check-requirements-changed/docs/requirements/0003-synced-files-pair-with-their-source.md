# 0003 Synced Files Pair With Their Source

## Why

I get a synced file judged by the repo it comes from, so a sync pull request needs no doc or test here.

## Measurement

synced files changed that pass as written or fail naming their source

-

synced files changed

## Requirements

1. Synced files pair with their source
   1. GIVEN Developer has a repo with requirement docs
   2. GIVEN its `.github/synced-files.json` names a script by its blob hash
      1. WHEN the script changes to match that hash
         1. THEN it passes with no changed test or requirement doc
      2. WHEN the script is edited by hand so its hash differs
         1. THEN it fails naming the file
         2. THEN it names the source repo to change
