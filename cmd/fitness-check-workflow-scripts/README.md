# workflow-scripts

Fails a `run:` step in a workflow or action that holds more than a trivial command ([#946](https://github.com/may-journal/fitness-runner/issues/946)). Shell inside YAML cannot be linted or tested, so real work belongs in a checked-in script.

## Behavior

Run `fitness explain workflow-scripts` for the files it reads. [Requirement 0001](docs/requirements/0001-run-steps-call-checked-in-scripts.md) states what passes and fails, and a test case named for each step proves it. A trivial step is one line with no control flow, heredoc, pipe, `;`, `&&` or `||`, such as `run: scripts/release.sh`.

In a pull request, or a staged commit, it reads only the changed files, and it never runs a step.

## Errors

```
.github/workflows/ci.yml:12: this run step holds 4 lines of shell; move it into a checked-in script such as scripts/<name>.sh, and call that script
action.yml:30: this run step holds a pipe; move it into a checked-in script such as scripts/<name>.sh, and call that script
```

## Fixing a step

Put the shell in `scripts/<name>.sh`, pass values through `env:`, and call the script from the step:

```
env:
  TAG: ${{ steps.release.outputs.tag }}
run: scripts/move-major-tag.sh
```
