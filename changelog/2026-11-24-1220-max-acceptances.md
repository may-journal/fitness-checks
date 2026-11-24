### 2026.11.24.1220

> A requirement doc holds at most eight acceptances by default, and the `requirements` check fails the ninth.

- Feat: `requirements` fails a doc holding more than `requirements.maxAcceptances` acceptances, 8 by default, at the first acceptance past the cap.
- Feat: `requirements` declares `requirements.maxAcceptances` in `--describe`, so `fitness config schema` and `fitness explain` list it.
- Test: requirement 0006.1 passes eight acceptances and fails nine; a key of 10 passes nine and fails eleven, and 3 fails four.
- Fix: the README key test reads a code span naming a check, such as `requirements`, as the check, not a key.
