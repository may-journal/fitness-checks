# 0001 A Repeated String Becomes A Shared Constant

## Why

I see a string that my JavaScript or TypeScript repeats three times, so I can share it as a constant.

## Measurement

repos whose strings repeat fewer than three times

-

repos with JavaScript or TypeScript sources

## Requirements

1. A string used three times fails
   1. GIVEN Developer has JavaScript or TypeScript sources
      1. WHEN a value appears three times across files
         1. THEN `repeated-string-literals` fails with one error for the value
         2. THEN it lists each location
         3. THEN it asks for a shared constant
      2. WHEN a value appears only twice
         1. THEN it passes
      3. WHEN a value is shorter than three characters
         1. THEN it is ignored
      4. WHEN a test file repeats a value
         1. THEN it counts like any source file
      5. WHEN several values repeat
         1. THEN the errors list the most repeated first
      6. WHEN a value repeats in more than five places
         1. THEN the error lists five locations
         2. THEN the error counts the rest
      7. WHEN the repo has no source files
         1. THEN it passes with zero files
2. Spellings with no shared constant pass
   1. GIVEN Developer has sources with strings that repeat
      1. WHEN a module specifier repeats
         1. THEN it passes
      2. WHEN a buffer encoding or a `typeof` result repeats
         1. THEN it passes
      3. WHEN a template literal or a regex or a comment repeats text
         1. THEN it passes
3. A repo allows a value
   1. GIVEN Developer sets `repeatedStringLiterals.allow` in `.fitnessrc.json`
      1. WHEN a listed value repeats
         1. THEN it passes
      2. WHEN another value repeats
         1. THEN it still fails
