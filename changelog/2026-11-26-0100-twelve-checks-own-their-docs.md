### 2026.11.26.0100

> Twelve checks own their requirement docs beside their code.

- Docs: twelve checks, from `checkstyle` to `go-test-coverage`, state their behavior as numbered acceptances in their own requirement 0001.
- Docs: their READMEs link that doc, and root requirements 0003, 0004, 0006 and 0007 drop the parts it owns.
- Build: `.fitnessrc.json` names the twelve check folders as package roots.
- Test: each acceptance has its test in the check's own `requirements` folder, and the root tests lose the moved cases.
