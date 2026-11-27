### 2026.11.27.0130

> A THEN with no tagged test case fails, and the old step names no longer count.

- Feat: `requirements` fails each THEN that no tag names, and prints the tag to add, 0010.3.
- Fix: a case named for a step no longer stands in for a tag, and 0005.3 is removed.
- Docs: Node, Swift and Java tests tag each THEN too.
