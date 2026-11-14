### 2026.11.14.2045

> Each check's README is generated here from its `doc.go`, so the two never disagree.

- Docs: `scripts/generate-readmes.sh` writes every check README from its doc comments and the config keys it reads, ported from fitness-runner's generator.
- Ci: the Checks workflow fails a README that differs from what the script writes.
- Docs: each README now names this repo's script as its source, and repo-local links in doc comments become relative.
