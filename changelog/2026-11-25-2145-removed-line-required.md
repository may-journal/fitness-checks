### 2026.11.25.2145

> A numbered requirement doc can no longer drop an acceptance without a `Removed` line, even its last one.

- Fix: `requirements` fails a numbered doc missing an acceptance ID its git history held, as an acceptance or a `Removed` line.
- Fix: the failure names the ID and the `N. Removed in YYYY.MM.DD.HHMM: <reason>` line to add.
- Docs: the check's requirement 0003.8 states it and its README says so.
