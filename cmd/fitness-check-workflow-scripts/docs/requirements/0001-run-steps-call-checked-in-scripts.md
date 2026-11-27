# 0001 Run Steps Call Checked In Scripts

## Why

I get shell that can be linted and tested. A workflow `run:` step holds only a trivial command and real work lives in a checked-in script.

## Measurement

run steps that hold a trivial command

-

run steps

## Requirements

1. Run steps stay trivial
   1. GIVEN Developer has workflow or action files with run steps
      1. WHEN a step is a one line command or a script call
         1. THEN it passes
         2. THEN it counts the files
      2. WHEN a block step holds more than one line of shell
         1. THEN it fails naming the file
         2. THEN it fails naming the line
         3. THEN it names the script to write
      3. WHEN a step holds a pipe or a chain or control flow
         1. THEN it fails naming the file of the step
         2. THEN it fails naming the line of the step
         3. THEN it fails naming what it holds
      4. WHEN operators sit in quotes or expressions or a trailing comment
         1. THEN it passes as plain text
      5. WHEN a plain scalar spans lines or a block holds `run:` text
         1. THEN the scalar is judged
         2. THEN the block text is left alone
      6. WHEN the repo has no workflow or action files
         1. THEN it passes with no files
