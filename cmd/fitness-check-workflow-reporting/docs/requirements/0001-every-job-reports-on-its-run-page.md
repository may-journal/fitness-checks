# 0001 Every Job Reports On Its Run Page

## Why

I see a job summary on every run page, so a failure never shows only a red badge.

## Measurement

workflow jobs that write a summary after any outcome

-

workflow jobs

## Requirements

1. Jobs report
   1. GIVEN Developer has workflow jobs with or without a summary step
      1. WHEN a job has a last step with `if: always()` that writes `$GITHUB_STEP_SUMMARY`
         1. THEN it passes
      2. WHEN a job with steps has no such step
         1. THEN it fails naming the file
         2. THEN it fails naming the job's line
         3. THEN it names the step to add
      3. WHEN a summary step lacks `if: always()`
         1. THEN it fails naming that step's line
      4. WHEN a job only calls a reusable workflow
         1. THEN it passes with no step needed
      5. WHEN a job is a matrix
         1. THEN it is judged once
         2. THEN it is named once
2. Workflow files are valid YAML
   1. GIVEN Developer has workflow files with values GitHub rejects or accepts
      1. WHEN a file is not valid YAML or holds a value GitHub rejects
         1. THEN it fails naming the line
         2. THEN it fails naming the column
         3. THEN the quoted or reshaped forms pass
      2. WHEN the repo has no workflow files
         1. THEN it passes with no files
      3. WHEN the repo's own workflows are read
         1. THEN they pass
