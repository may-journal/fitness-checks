# 0001 Workflow Calls Resolve At Their Refs

## Why

I call a may-journal workflow or action only at a tag or sha that holds its path. No caller then fails every event with a missing-workflow error.

## Measurement

pinned workflow calls whose path exists at their ref

-

pinned workflow calls read

## Requirements

1. Pinned paths exist
   1. GIVEN Developer has workflow files with `uses:` lines
   2. GIVEN GitHub answers for each repo
      1. WHEN each path exists at its tag or sha
         1. THEN it passes
         2. THEN it counts the workflow files
      2. WHEN a workflow or action path is missing at its tag
         1. THEN it fails naming the file and line and path and ref
      3. WHEN a ref is not a tag or a branch or a full sha or a sha names no commit
         1. THEN it fails naming the file and line and ref
      4. WHEN a call is local or another owner's or built from an expression or in a script block
         1. THEN it is not read
         2. THEN it passes
      5. WHEN the repo has no workflow files
         1. THEN it passes with no files
2. GitHub cannot always answer
   1. GIVEN Developer has pinned calls that GitHub cannot always answer
      1. WHEN a repo answers 404 for the token
         1. THEN its calls are skipped with a note
      2. WHEN GitHub refuses a call or fails to answer
         1. THEN the findings so far are kept
         2. THEN the calls left unverified are noted
      3. WHEN the service cannot answer at all
         1. THEN it passes with info
         2. THEN it never fails
      4. WHEN many pinned calls each answer slowly
         1. THEN they are read in parallel
