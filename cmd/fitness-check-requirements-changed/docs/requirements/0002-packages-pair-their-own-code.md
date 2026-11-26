# 0002 Packages Pair Their Own Code

## Why

I change a package's code with its own test plus doc, since each package owns them.

## Measurement

package code changes that come with their own test and doc

-

package code changes

## Requirements

1. Packages pair their own code
   1. GIVEN Developer has a monorepo whose `.fitnessrc.json` names its package folders
   2. GIVEN each package holds its own requirement docs
   3. GIVEN each package holds its own tests
      1. WHEN a package's code changes beside another package's test plus doc
         1. THEN it fails naming the file
      2. WHEN a package's code changes beside its own test plus doc
         1. THEN it passes
      3. WHEN a package's code is deleted while a doc names it
         1. THEN it fails naming the file with the doc
