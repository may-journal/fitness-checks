# markdown-links

Requires every relative link in every markdown file to resolve to an existing file or directory. Broken relative links survive refactors silently — a docs move or a file deletion leaves pointers that no other check sees. This check makes them commit-blocking.

## Behavior

[Requirement 0001](docs/requirements/0001-relative-links-resolve.md) states what this check does, and a test case named for each step proves it.

Run `fitness explain markdown-links` for the files it reads.

## Errors

```
docs/competition.md:11: broken relative link: ./README.md
```
