### 2026.11.15.0700

> A workflow or action step that does real work is a script file, since shell inside YAML cannot be linted or tested.

- Feat: `workflow-scripts` fails a `run:` step in a workflow or `action.yml` that holds more than one trivial command.
- Fix: the check reads quoted and multi-line steps, and skips comments, escaped semicolons and block bodies.
- Test: requirements 0167.1 to 0167.5 run it on workflows and actions.
- Test: the catalog lists workflow-scripts, and requirements 0160.2 and 0160.3 expect no bundled copy of it.
- Docs: its README is the fixed stub that points at the doc comment.
