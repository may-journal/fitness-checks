# 0002 GitHub Rejected Values And Files Read

## Why

I see each value GitHub rejects at its place, and only the workflow files a change touches, so a fix is small.

## Measurement

rejected workflow values that fail at their place

-

rejected workflow values

## Requirements

1. Rejected values fail at their place
   1. GIVEN Developer has a workflow job whose `runs-on` holds a plain value
      1. WHEN the value holds a colon plus space or ends in a colon
         1. THEN it fails at its line plus column
         2. THEN the quoted value passes
      2. WHEN the value starts with a character YAML reserves
         1. THEN it fails naming the character
      3. WHEN a list or a key or a value opens on its key's line
         1. THEN it fails telling to move it onto its own line
      4. WHEN text other than a comment follows a closing bracket
         1. THEN it fails saying the text follows the bracket
2. Each run reads only the workflow files it needs
   1. GIVEN Developer has a workflow whose job lacks a report step
      1. WHEN a staged commit or a pull request changes another workflow
         1. THEN the untouched workflow is not judged
      2. WHEN a workflow has a step that writes a file
         1. THEN no step is run
