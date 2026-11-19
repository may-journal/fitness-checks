# 0003 History Stays Clean

## Why

I read a branch's history and release notes and see what each change was, who helped, and which issue it served. Each issue is linked once, and each change has one dated entry.

## Measurement

commits and entry files that pass every history rule

-

commits and entry files checked

## Requirements

- 0003.1
  - Given Developer has commit subjects with and without a type and scope
    - When the commit checks read the subject line
      - Then only subjects with a listed type and scope pass, and each failure quotes the subject
- 0003.2
  - Given Developer has commit messages with correct and misspelled `Plan` references
    - When the commit checks read the Plan trailer
      - Then every misspelled trailer fails with the fix, and prose, other repos and absent trailers pass
- 0003.3
  - Given Developer has commits with and without the AI tool and model trailers
    - When the commit checks read the attribution trailers
      - Then an adopting repo fails naming each missing trailer, and merges and reverts pass
- 0003.4
  - Given Developer has a branch and open PR that already link some issues
    - When the commit checks look for repeated issue links
      - Then a second link to the same issue fails naming the earlier commit or PR
- 0003.5
  - Given Developer has changelog entry files with good and bad names and headings
    - When the changelog checks read each entry file
      - Then each file needs one dated heading matching its name, and upstream `CHANGELOG.md` history passes as it is
- 0003.6
  - Given Developer has changelog entry files with different bullet counts, lengths and types
    - When the changelog checks read the bullets
      - Then an entry with fewer than three or more than five bullets, a long bullet, or an untyped bullet fails naming it
- 0003.7
  - Given Developer has staged commits with and without a new or edited entry file
    - When the changelog checks read the staged commit
      - Then the branch's one entry must name the change, and published entries stay as released
- 0003.8
  - Given Developer has release manifests and version files that agree or differ
    - When the changelog checks read the release versions
      - Then a mismatched or missing version fails naming it, and repos with no manifest pass
