# 0002 Packages Own Their Requirements

## Why

I keep each package's requirement docs and tests in its own folder, judged as fast as one repo's.

## Measurement

package requirement docs judged on their own history

-

package requirement docs

## Requirements

- 0002.1
  - Given Developer has a monorepo whose `.fitnessrc.json` names package roots that each hold their own requirement docs
    - When `requirements` reads their ID history
      - Then one git log serves them all and each package keeps its own deleted IDs
