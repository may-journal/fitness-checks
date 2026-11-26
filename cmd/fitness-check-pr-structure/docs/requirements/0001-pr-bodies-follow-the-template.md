# 0001 PR Bodies Follow The Template

## Why

I read every PR in one shape: a summary and its background and its changelog.

## Measurement

PR bodies that follow the template

-

PR bodies

## Requirements

1. A PR body follows the template
   1. GIVEN Developer has a PR body
      1. WHEN it holds a summary plus a Background section plus a Changelog bullet
         1. THEN `pr-structure` passes it
      2. WHEN it has no blockquote summary
         1. THEN it fails
         2. THEN it asks for a one-line summary
      3. WHEN its summary is still the template placeholder
         1. THEN it fails
         2. THEN it asks for the real summary
      4. WHEN it has no Background section
         1. THEN it fails naming the missing section
      5. WHEN its Changelog has no bullet
         1. THEN it fails naming the empty section
      6. WHEN it has an extra section
         1. THEN it fails naming the unexpected section
      7. WHEN an extra heading sits only inside fenced code
         1. THEN it passes
2. A missing PR body passes
   1. GIVEN Developer runs a check with no PR body
      1. WHEN `pr-structure` runs
         1. THEN it passes
         2. THEN it judges zero files
