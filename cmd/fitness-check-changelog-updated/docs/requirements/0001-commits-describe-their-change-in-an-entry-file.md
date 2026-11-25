# 0001 Commits Describe Their Change In An Entry File

## Why

I get each commit's changelog entry checked for being about its change before I commit it.

## Measurement

staged commits whose entry file shares three words with their diff

-

staged commits

## Requirements

- 0001.1
  - Given Developer has staged a change in a repo
    1. When Git is mid-merge, mid-cherry-pick or mid-revert
       - Then it skips, since a replay re-stages old entries
    2. When nothing is staged
       - Then it skips, judging no files
    3. When no entry file is staged
       - Then it fails asking to stage a `changelog/<yyyy-mm-dd-HHMM>-<change>.md` entry file about the change
    4. When the entry's added lines share 3 words with the rest of the staged diff
       - Then it passes, whatever `CHANGELOG.md` holds, even when the commit creates `changelog/`
    5. When they share fewer, with `CHANGELOG.md` counted as part of the rest
       - Then it fails with the count, suggesting up to 10 words from the diff
    6. When words are short, cased differently, or on unchanged lines
       - Then only added words of 3 or more letters count, lowercased
- 0001.2
  - Given Developer has a branch, with release tags, whose commits add entry files
    1. When a commit adds an entry file whose name lacks its heading's date
       - Then it fails naming the file it should be
    2. When an earlier commit on the branch added its entry, and this one adds another
       - Then it fails naming the branch's entry and each new one
    3. When the branch stacks on another `origin` branch
       - Then that branch's commits do not count, so its first commit adds its entry
    4. When the commit is on the default branch, or the repo has no `origin` refs
       - Then each commit is judged alone, and may add its own entry
    5. When a commit edits or deletes an entry file the newest release tag holds
       - Then it fails asking for a new entry file instead
    6. When a commit only drops that entry's front matter
       - Then it passes
