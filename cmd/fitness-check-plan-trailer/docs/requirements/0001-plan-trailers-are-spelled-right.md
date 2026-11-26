# 0001 Plan Trailers Are Spelled Right

## Why

I link a commit to its Plan with one trailer that every tool reads the same way.

## Measurement

commit messages whose Plan trailer is spelled right

-

commit messages

## Requirements

1. Plan trailers read one way
   1. GIVEN Developer has a commit message that may name a Plan
      1. WHEN it holds no Plan reference
         1. THEN `plan-trailer` passes it
      2. WHEN it holds a trailer reading `Plan #63`
         1. THEN it passes
      3. WHEN it holds a Plan line in another spelling
         1. THEN it fails once per line
         2. THEN its fix names all three spellings
      4. WHEN body prose starts with the word Plan
         1. THEN it passes because prose is never a trailer
      5. WHEN no message is supplied
         1. THEN it reads the HEAD commit
      6. WHEN it names a Plan of another repo as `Plan <owner>/<repo>#<number>` or as an issue URL
         1. THEN it passes
      7. WHEN a Plan line of another repo is misspelled
         1. THEN it fails with the same fix
