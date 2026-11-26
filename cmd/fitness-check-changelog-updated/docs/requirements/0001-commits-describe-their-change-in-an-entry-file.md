# 0001 Commits Describe Their Change In An Entry File

## Why

I get each commit's changelog entry checked for being about its change before I commit it.

## Measurement

staged commits whose entry file shares three words with their diff

-

staged commits

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has staged a change in a repo
  WHEN Git is mid-merge
    AND Git is mid-cherry-pick
    AND Git is mid-revert
  THEN it skips
    BUT a replay re-stages old entries
  WHEN nothing is staged
  THEN it skips
    AND it judges no files
  WHEN no entry file is staged
  THEN it fails asking to stage a `changelog/<yyyy-mm-dd-HHMM>-<change>.md` entry file about the change
  WHEN the entry's added lines share 3 words with the rest of the staged diff
  THEN it passes whatever `CHANGELOG.md` holds
    AND it passes even when the commit creates `changelog/`
  WHEN they share fewer
    AND `CHANGELOG.md` counts as part of the rest
  THEN it fails with the count
    AND it suggests up to 10 words from the diff
  WHEN words are short
    AND words are cased differently
    AND words sit on unchanged lines
  THEN only added words of 3 or more letters count
    AND they count lowercased

SCENARIO: 0001.2
  GIVEN Developer has a branch with release tags AND commits that add entry files
  WHEN a commit adds an entry file whose name lacks its heading's date
  THEN it fails naming the file it should be
  WHEN an earlier commit on the branch added its entry
    AND this commit adds another
  THEN it fails naming the branch's entry
    AND it names each new one
  WHEN the branch stacks on another `origin` branch
  THEN that branch's commits do not count
    AND its first commit adds its entry
  WHEN the commit is on the default branch
    AND the repo has no `origin` refs
  THEN each commit is judged alone
    AND it may add its own entry
  WHEN a commit edits an entry file the newest release tag holds
    AND a commit deletes an entry file the newest release tag holds
  THEN it fails asking for a new entry file instead
  WHEN a commit only drops that entry's front matter
  THEN it passes
```
