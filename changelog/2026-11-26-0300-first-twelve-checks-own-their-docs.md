### 2026.11.26.0300

> Twelve checks, `markdown-no-bold-italic` to `no-plans-dir`, own their requirement docs beside their code.

- Docs: each of the twelve states its behavior as numbered acceptances in its own requirement 0001, which its README links.
- Docs: root requirements 0002.5, 0005.1 to 0005.4 and 0007.2 leave with their checks, and 0004.2 and 0006.5 keep what remains.
- Build: `.fitnessrc.json` names the twelve as package roots.
- Test: each step has a named case in its check's own test, moved out of the root tests.
