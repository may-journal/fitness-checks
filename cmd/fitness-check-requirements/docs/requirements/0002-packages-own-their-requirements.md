# 0002 Packages Own Their Requirements

## Why

I keep each package's requirement docs and tests in its own folder, judged as fast as one repo's.

## Measurement

package requirement docs judged on their own history

-

package requirement docs

## Requirements

1. One log serves every package
   1. GIVEN Developer has a monorepo whose `.fitnessrc.json` names package roots
   2. GIVEN each package root holds its own requirement docs
      1. WHEN `requirements` reads their ID history
         1. THEN one git log serves them all
         2. THEN each package keeps its own deleted IDs
2. Each package owns its docs and tests
   1. GIVEN Developer has a monorepo whose `.fitnessrc.json` names its package folders
      1. WHEN `requirements` judges its docs and tests
         1. THEN each package with code owns its docs and tests
         2. THEN a package without code needs none
      2. WHEN two packages number a doc alike
         1. THEN the number repeats only across packages
      3. WHEN `.fitnessrc.json` names no packages
         1. THEN the repo is one package
3. The repo root keeps only monorepo promises
   1. GIVEN Developer has a monorepo whose `.fitnessrc.json` names its package folders
      1. WHEN a root doc acceptance names a single fitness-check package
         1. THEN it fails naming the acceptance and the check
      2. WHEN a package doc or a Removed line names a check
         1. THEN it passes
