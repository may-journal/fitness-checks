### 2026.11.19.0900

> A repo can restart its whole requirement set from 0001, while everyday changes still never reuse a number.

- Feat: requirements lets a commit that removes or renames away every existing doc reuse any number.
- Feat: after such a restructure, the never-reuse rule counts only removals from it onward.
- Feat: a staged restructure passes before its commit, so the pre-commit hook lets it land.
- Test: a staged and committed restructure passes, and dropping then restoring an ID after it still fails.
- Docs: the check's Reuse rule and acceptance 0006.3 describe the restructure exception.
