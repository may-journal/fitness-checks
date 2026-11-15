### 2026.11.15.0700

> A workflow or action step that does real work is a script file, since shell inside YAML cannot be linted or tested.

- Feat: `workflow-scripts` fails a `run:` step in a workflow or `action.yml` that holds more than one trivial command.
- Test: requirement 0007.9 runs it on workflows and actions.
- Docs: the README names the threshold and shows how to move a step into a script.
- Test: requirement 0002 lists workflow-scripts and the workflow and action files it reaches.
- Fix: the check reads quoted and multi-line steps, and skips comments, escaped semicolons and block bodies.
