### 2026.11.23.1430

> Two journey steps that name each other's requirement docs no longer pass.

- Fix: `architecture-journey` fails two steps whose labels read as each other's docs, naming both IDs at each step's line.
- Feat: one step that reads as a doc it does not name adds a note, and the check still passes.
- Test: requirement 0005.6 crosses two step labels, straightens them, and words one step loosely.
- Docs: the doc comment explains how a label is compared with a doc's title words.
