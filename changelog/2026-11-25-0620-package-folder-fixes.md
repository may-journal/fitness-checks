### 2026.11.25.0620

> `prose-budget` and `requirements` handle package roots and own the docs for that, so every check can move its docs beside its code.

- Feat: `requirements.maxSteps` sets how many steps an acceptance chains, 7 by default, and `prose-budget` lets package docs chain as many.
- Fix: CI maps a change to a check's own `requirements` folder to that check, rather than naming the folder as a check.
- Perf: `requirements` reads the ID history of every package root in one git log, not one per root, under any git diff config.
- Docs: both checks are package roots with their own requirement docs, and root 0001.2 covers a check's own tests.
- Test: the shared doc and repo fixtures move to `internal/acceptance`, so a check's own tests use them.
