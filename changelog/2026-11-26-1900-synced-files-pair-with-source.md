### 2026.11.26.1900

> A file the org sync wrote needs no doc or test where it lands, so sync pull requests pass `requirements-changed`.

- Feat: `requirements-changed` passes a changed file whose staged blob matches `.github/synced-files.json`, the record the sync writes.
- Feat: a synced file edited by hand fails, naming the source repo to change instead.
- Docs: `requirements-changed` gains requirement 0003, with a test case for each step.
