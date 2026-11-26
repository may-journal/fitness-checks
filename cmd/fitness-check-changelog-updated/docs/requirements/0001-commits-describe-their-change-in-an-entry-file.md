# 0001 Commits Describe Their Change In An Entry File

## Why

I get each commit's changelog entry checked for being about its change before I commit it.

## Measurement

staged commits whose entry file shares three words with their diff

-

staged commits

## Requirements

1. Staged changes carry their entry
   1. GIVEN Developer has staged a change in a repo
      1. WHEN Git is mid-merge
         1. THEN it skips
         2. THEN a mid-cherry-pick run skips too
         3. THEN a mid-revert run skips too
      2. WHEN nothing is staged
         1. THEN it skips
      3. WHEN no entry file is staged
         1. THEN it fails asking to stage a `changelog/<yyyy-mm-dd-HHMM>-<change>.md` entry file about the change
      4. WHEN the entry's added lines share 3 words with the rest of the staged diff
         1. THEN it passes whatever `CHANGELOG.md` holds
      5. WHEN they share fewer words
         1. THEN it counts `CHANGELOG.md` as part of the rest
         2. THEN it fails with the count
      6. WHEN words are short or cased differently or on unchanged lines
         1. THEN only added words of 3 or more letters count
2. A branch adds one entry per commit
   1. GIVEN Developer has a branch with release tags
   2. GIVEN its commits add entry files
      1. WHEN a commit adds an entry file whose name lacks its heading's date
         1. THEN it fails naming the file it should be
      2. WHEN an earlier commit on the branch added its entry and this one adds another
         1. THEN it fails naming the branch's entry
         2. THEN it names each new one
      3. WHEN the branch stacks on another `origin` branch
         1. THEN that branch's commits do not count
         2. THEN its first commit adds its entry
      4. WHEN the commit is on the default branch
         1. THEN each commit is judged alone
         2. THEN it may add its own entry
      5. WHEN a commit edits an entry file the newest release tag holds
         1. THEN it fails asking for a new entry file instead
      6. WHEN a commit only drops that entry's front matter
         1. THEN it passes
