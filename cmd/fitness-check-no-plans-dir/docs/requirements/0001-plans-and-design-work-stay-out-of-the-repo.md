# 0001 Plans And Design Work Stay Out Of The Repo

## Why

I find plans in Plan issues and mockups and reports in comments. The repo holds only working source.

## Measurement

Repo paths that hold no plan or temporary design work

-

Repo paths

## Requirements

1. Plan files fail by path
   1. GIVEN Developer has a repo
   2. GIVEN it may hold plan or design files
      1. WHEN a file sits in a `plans` folder at any depth
         1. THEN it fails naming the path
         2. THEN it names a Plan issue
      2. WHEN files are maintained docs or shipped images or tests quoting plan text
         1. THEN they pass
         2. THEN every file stays in place
      3. WHEN a file sits in a `mockups` or `validation` or `review-artifacts` folder
         1. THEN it fails in path order
         2. THEN it names the comment it belongs in
      4. WHEN a file is named `*-plan.md` or `*-mockup.*` or `pr-<number>-*`
         1. THEN it fails naming the rule
         2. THEN it names the comment it belongs in
      5. WHEN a Markdown file outside `.github` holds a `## What needs to happen` checklist
         1. THEN it fails naming the Plan issue
      6. WHEN a plan file is staged or a rename lands in a mockups folder
         1. THEN it fails before the commit
      7. WHEN git cannot list the repo files
         1. THEN it fails saying it cannot list them
