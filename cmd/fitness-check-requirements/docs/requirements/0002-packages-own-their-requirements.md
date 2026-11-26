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
