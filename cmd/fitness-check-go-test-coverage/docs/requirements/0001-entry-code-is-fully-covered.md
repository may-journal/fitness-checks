# 0001 Entry Code Is Fully Covered

## Why

I get every command package and declared library entry reaching full statement coverage and a report of tests that overlap or run unstably.

## Measurement

entry packages at full statement coverage

-

entry packages

## Requirements

1. Entries are declared
   1. GIVEN Developer has a repo with or without Go modules
      1. WHEN the repo has no Go modules
         1. THEN it passes with no files
      2. WHEN Go is not on PATH
         1. THEN it reports info with the hint `see https://go.dev/dl`
      3. WHEN a module has a library with no declared entry
         1. THEN it fails asking to declare entries
      4. WHEN `--entry` names an unknown package
         1. THEN it fails naming that package
2. Entries reach full coverage
   1. GIVEN Developer has a repo with Go modules with declared entries
      1. WHEN an entry runs fewer statements than it holds
         1. THEN it fails naming each entry
         2. THEN it names the entry's covered statements
         3. THEN it names the entry's total statements
      2. WHEN a top-level test passes a named helper to `t.Run`
         1. THEN it fails naming the file with its line
         2. THEN it fails naming the helper
3. Audits write their map and report overlap
   1. GIVEN Developer has a repo with Go modules that have tests
      1. WHEN `--audit-map` names a new path outside the repo
         1. THEN it passes
         2. THEN it writes the source map there
      2. WHEN a test covers different blocks on each run
         1. THEN it reports the test as unstable
         2. THEN it excludes the test from overlap advice
      3. WHEN the audit runs parallel tests plus serial tests
         1. THEN parallel tests run alongside each other
         2. THEN serial packages run side by side
         3. THEN tests within one serial package never run side by side
