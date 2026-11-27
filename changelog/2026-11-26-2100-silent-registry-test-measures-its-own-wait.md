### 2026.11.26.2100

> The silent-registry test measures the check's own wait, so machine load no longer fails it.

- Fix: `dependency-currency` reads `FITNESS_DEPENDENCY_CURRENCY_TIMEOUT` to shorten its request timeout and fetch budget.
- Docs: `dependency-currency` gains 0001.3.2, which says a silent registry is dropped within the request timeout.
- Fix: the 0001.3 test times the held connection at the listener, not the whole run.
