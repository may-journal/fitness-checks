# 0010 Each THEN Has A Tagged Case

## Why

I tie each THEN to the case that proves it, so a reworded THEN fails until its case is reread.

## Measurement

THENs whose tag names their current text

-

THENs

## Requirements

1. A tag names the current text of its THEN
   1. GIVEN Developer has a repo with requirement docs
   2. GIVEN it has tests whose cases carry tags of an ID plus a hash
      1. WHEN a tag names no THEN or holds a hash other than its THEN's
         1. THEN it fails naming the tag
         2. THEN it prints the new hash of a stale tag
      2. WHEN two tags name one THEN
         1. THEN it fails naming the THEN
         2. THEN it names both places
      3. WHEN every tag names a THEN with its current hash
         1. THEN it passes beside cases named for steps
      4. WHEN a tag sits outside the test of its acceptance
         1. THEN it fails naming the test to move it to
2. A tag mode writes the tags
   1. GIVEN Developer has a Go repo whose cases are named for steps
      1. WHEN `--tag` runs
         1. THEN each case of a step with one THEN gains the tag of that THEN
         2. THEN it lists each THEN still untagged
      2. WHEN `--tag --dry-run` runs
         1. THEN it writes no file
