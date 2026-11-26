### 2026.11.26.0200

> Ten more checks, `requirements-changed` to `workflow-scripts`, own their requirement docs beside their code, and the requirements doc rules become acceptances.

- Docs: each of the ten states its behavior as numbered acceptances in its own requirement doc, which its README links.
- Docs: `requirements` gains 0004 for the template rules and 0005 for ownership, which its README links.
- Docs: root requirements 0004.1, 0006.1, 0007.7 and 0007.8 are removed and 0002.3, 0003.1, 0004.3, 0006.2, 0006.3, 0006.6 and 0006.8 keep only what remains.
- Build: `.fitnessrc.json` names the ten as package roots.
- Test: each step has a named case in its check's own test, moved out of the root tests where one existed.
