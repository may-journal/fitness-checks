# 0007 Tests Fit Their Language

## Why

I name tests the way each language runs them, and each check that proves a requirement counts the same.

## Measurement

acceptances whose owner the check finds in its language

-

acceptances

## Requirements

1. Each language finds its own tests
   1. GIVEN Developer has a repo written in Swift or JavaScript or Java
      1. WHEN a test is declared the way its language runs it
         1. THEN it owns its acceptance
         2. THEN a name in a comment or a string or a helper owns nothing
      2. WHEN an acceptance has no such test
         1. THEN it fails naming the test form to write
      3. WHEN a Node test runs its cases
         1. THEN it needs only its own test until the Node switch is on
      4. WHEN a test is named for no acceptance
         1. THEN it fails naming the form to rename it to
      5. WHEN a Swift or Java test has no case for a THEN
         1. THEN it fails naming the THEN
2. Code needs requirement docs
   1. GIVEN Developer has a repo with code in a covered language
      1. WHEN docs/requirements holds no doc
         1. THEN it fails naming the folder
      2. WHEN the repo holds no covered code
         1. THEN it passes with no files
3. Swift tests use the UI testing library
   1. GIVEN Developer has a Swift package with requirement tests
      1. WHEN a test imports another test framework or drives the app another way
         1. THEN it fails naming the line
         2. THEN it fails naming the library to use
      2. WHEN a Swift Testing test uses UITestingLibrary
         1. THEN it passes
4. A check can prove a requirement
   1. GIVEN Developer has a `fitness/req-NNNN` check beside a requirement doc
      1. WHEN the check passes
         1. THEN the acceptances need no test
      2. WHEN the check has no doc or changes without its doc
         1. THEN it fails naming the check
      3. WHEN the check judges no files
         1. THEN it fails since its scan is broken
