# 0001 Each Issue Is Linked Once

## Why

I read an issue's timeline and see one entry for each change that served it.

## Measurement

branches and pull requests that link each issue once

-

branches and pull requests that link issues

## Requirements

1. A commit message repeats no link its branch holds
   1. GIVEN Developer has a branch whose earlier commit links an issue
   2. GIVEN the repo's origin is on GitHub
      1. WHEN a new message links the same issue by number
         1. THEN it fails naming the issue
         2. THEN it fails naming the earlier commit
         3. THEN a link by URL fails the same way
      2. WHEN the issue a new message links belongs to another repo
         1. THEN it passes
      3. WHEN a comment line or a line below the scissors line holds the link
         1. THEN it passes since git drops those lines
      4. WHEN neither a message nor a pull request exists
         1. THEN it passes with no files
2. A commit message repeats no link the open pull request holds
   1. GIVEN Developer has a branch with an open pull request
   2. GIVEN the pull request body links an issue
      1. WHEN a new message links that issue
         1. THEN it fails naming the issue
         2. THEN it fails naming the pull request
3. A pull request links each issue from one commit at most
   1. GIVEN Developer has a pull request run on CI
   2. GIVEN the branch links an issue in more than one place
      1. WHEN a second commit repeats a link
         1. THEN it fails naming the commit
         2. THEN it fails naming the earlier one
         3. THEN it tells how to rewrite with `--force-with-lease`
      2. WHEN a later commit repeats the pull request's link
         1. THEN it fails naming the pull request
      3. WHEN one early link is judged with main merged in
         1. THEN it passes counting the files scanned
