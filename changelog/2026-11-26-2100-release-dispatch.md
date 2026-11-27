### 2026.11.26.2100

> A check release tells the org's .github repo, which reruns the fitness suite on every open pull request.

- Feat: the Release workflow sends a `fitness-check-released` dispatch to the .github repo after it publishes a check.
- Feat: a push that released no check sends none.
- Docs: requirement 0009 states the announcement, with a test for each acceptance.
