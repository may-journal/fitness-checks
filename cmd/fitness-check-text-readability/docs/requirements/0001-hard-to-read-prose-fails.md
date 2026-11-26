# 0001 Hard To Read Prose Fails

## Why

I can read any Markdown doc or code comment, because a check fails prose that scores as hard to read.

## Measurement

Markdown files and comment blocks whose prose passes the readability check

-

Markdown files and comment blocks judged

## Requirements

1. Dense documents fail
   1. GIVEN Developer has Markdown files with plain or dense prose
      1. WHEN a file scores in alarm on at least two of the three formulas
         1. THEN it fails naming the file
         2. THEN it names its prose words plus its scores
         3. THEN it adds the formulas plus the fix
         4. THEN it adds the methodology
      2. WHEN a file has short plain sentences
         1. THEN it passes
      3. WHEN a file has fewer than 100 prose words
         1. THEN it passes without a score
      4. WHEN only one formula is in alarm
         1. THEN it passes
      5. WHEN dense text sits inside fenced code
         1. THEN it passes since code is never scored
2. Judging follows the repo and the run
   1. GIVEN Developer has a repo with Markdown files
      1. WHEN an issue or a pull request body is passed with `--body-file`
         1. THEN a dense body fails naming the description
      2. WHEN the repo sets lower bands in `.fitnessrc.json`
         1. THEN the files are judged with those bands
      3. WHEN `--report` is passed
         1. THEN it prints every file's word count with its scores
         2. THEN the verdict is unchanged
3. Dense comments fail
   1. GIVEN Developer has a code comment block of 100 prose words or more
      1. WHEN the block scores in alarm
         1. THEN it fails naming the block as `path:line`
