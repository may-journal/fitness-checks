### 2026.11.17.1500

> Requirements sit together by area in the order a user meets them, numbered 0001 to 0088 with no gaps.

- Docs: requirement docs renumber by area, from engine contract and requirements through architecture, commits, changelog, languages, Markdown, Mermaid and prose.
- Test: requirement test files, Test functions and numbered helpers take their doc's new number.
- Docs: the user journey labels and check `doc.go` links name the new numbers.
- Build: `scripts/renumber-requirements.sh` applies the old-to-new mapping in `scripts/renumber-requirements.txt`.
