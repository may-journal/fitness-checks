### 2026.11.27.0030

> Each THEN ties to its test case by its ID plus a hash of its text, so test names are free.

- Feat: `requirements` fails a tag that names no THEN, a stale hash, a doubled tag or a tag outside its acceptance's test, 0010.1.
- Feat: `requirements --tag` tags each Go case that asserts through a counting helper, and `--dry-run` writes nothing, 0010.2.
- Docs: every requirement test here tags its THENs, and the old step names stay until the next release.
