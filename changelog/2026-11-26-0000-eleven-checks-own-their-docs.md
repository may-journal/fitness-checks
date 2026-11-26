### 2026.11.26.0000

> Eleven checks, `node-version` to `repeated-string-literals`, own their requirement docs beside their code.

- Docs: each states its behavior as numbered acceptances in its own requirement doc, which its README links.
- Docs: `prose-budget` gains requirement 0002 for comments and extends 0001 with the prose limits.
- Docs: root requirements 0002.4, 0003.2, 0003.8, 0004.7 and part of 0004.1, 0004.3, 0004.6, 0006.5 and 0006.7 leave with their checks.
- Build: `.fitnessrc.json` names the eleven as package roots.
- Test: each step has a named case in its check's own test, moved out of the root tests where one existed.
