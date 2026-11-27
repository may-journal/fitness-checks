# 0001 Ignore Patterns Say Why

## Why

I get every `.gitignore` pattern explained by a comment, so ignores stay easy to audit and onboard.

## Measurement

repos whose `.gitignore` patterns each follow an explanatory comment

-

repos with a `.gitignore`

## Requirements

1. Commented patterns pass
   1. GIVEN Developer has a repo with or without a `.gitignore`
      1. WHEN there is no `.gitignore`
         1. THEN it passes with no files
      2. WHEN every pattern has a `#` comment with text on the line directly above it
         1. THEN it passes
         2. THEN it counts the one file
2. Patterns without their own comment fail
   1. GIVEN Developer has a repo with a `.gitignore`
      1. WHEN a pattern has no comment of its own with text on the line directly above it
         1. THEN it fails naming the line
         2. THEN it fails naming the pattern as `.gitignore:<line>: pattern "<text>" has no explanatory # comment on the line above`
