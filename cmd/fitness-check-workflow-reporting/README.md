# workflow-reporting

Fails a workflow job that can fail without a report on its run page ([#181](https://github.com/may-journal/fitness-runner/issues/181)). [ADR 0003](https://github.com/may-journal/fitness-runner/blob/main/fitness/adr-0003/README.md) asks every job for a job summary, so a failure never shows only a red badge.

## Behavior

[Requirement 0001](docs/requirements/0001-every-job-reports-on-its-run-page.md) states what this check passes and fails, and a test case named for each step proves it. It reads the `.yml` and `.yaml` workflow files directly under `.github/workflows`, keeping each value's line.

A job with `steps` needs one step with `if: always()`, bare or as `${{ always() }}`, whose `run` writes to `$GITHUB_STEP_SUMMARY`.

## Values GitHub rejects

These values fail at their line and column, as GitHub fails them:

- A plain value holding a `:` followed by a space, or ending in `:`. Quote the value.
- A plain value starting with `,`, `]`, `}`, `%`, `@` or a backtick, which YAML reserves. Quote the value.
- A list, key or value opened with `- `, `? ` or `: ` on its key's line. Move it onto its own indented line.
- Text other than a comment after the closing bracket of a `[...]` or `{...}` value.

In a pull request, or a staged commit, it reads only the changed workflow files. It reads the files only: it never runs a step or fetches a called workflow.

## Errors

```
.github/workflows/ci.yml:8: job build has no step that runs with if: always() and writes $GITHUB_STEP_SUMMARY, so a failed run shows only a badge; add a last step with `if: ${{ always() }}` whose run writes the job's status to "$GITHUB_STEP_SUMMARY"
.github/workflows/ci.yml:3: this workflow is not valid YAML: a tab indents this line; YAML indents with spaces; fix the YAML on this line
.github/workflows/ci.yml:5:30: this workflow is not valid YAML: a `: ` inside a plain value starts a mapping, which YAML does not allow here; quote the value; fix the YAML on this line
```

## Adding the step

A last step like this one reports every outcome:

```
steps:
  - name: Report job outcome
    if: ${{ always() }}
    run: printf '### Job %s: %s\n' "$GITHUB_JOB" "${{ job.status }}" >> "$GITHUB_STEP_SUMMARY"
```
