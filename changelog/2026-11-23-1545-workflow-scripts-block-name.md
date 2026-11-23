### 2026.11.23.1545

> A run step is judged even when its step opens with a block scalar name.

- Fix: `workflow-scripts` counts the list dash in a block scalar's indent, so the step's other keys are no longer its text.
- Test: requirement 0007.8 names a step with `>-` and its piped run step fails.
- Docs: the acceptance says each failure names its own line.
