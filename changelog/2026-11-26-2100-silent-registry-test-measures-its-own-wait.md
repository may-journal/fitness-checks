### 2026.11.26.2100

> The silent-registry test measures the check's own wait, so machine load no longer fails it.

- Fix: `dependency-currency` reads `FITNESS_DEPENDENCY_CURRENCY_TIMEOUT` to shorten its request timeout and fetch budget.
- Docs: requirement 0001.3's test times the held connection at the listener, not the whole run.
