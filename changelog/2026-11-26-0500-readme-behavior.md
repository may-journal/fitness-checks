### 2026.11.26.0500

> A new `readme-behavior` check holds each check README to one Behavior link and a full Example.

- Feat: `readme-behavior` fails a Behavior section that holds more than one link to its own requirements.
- Feat: it fails a check README whose Example misses its command or a flag or config key `--describe` declares.
- Docs: every check README keeps one Behavior link and gains an Example, with first paragraphs unchanged.
- Docs: root requirement 0008.7 moves to the new check, and the doc move script writes the single link.
