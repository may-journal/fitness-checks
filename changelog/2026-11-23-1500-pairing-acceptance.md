### 2026.11.23.1500

> A changed requirement test passes only beside a changed acceptance, so filler in its doc no longer satisfies pairing.

- Fix: `requirements` pairs a changed test only with a doc whose Requirements section changed.
- Fix: a doc edit outside that section, or a format-only rewrite, no longer lets a changed test through.
- Test: requirement 0006.3 adds a Why line that names the test file, and the changed test still fails.
- Docs: the Pairs section of the doc comment states the rule.
