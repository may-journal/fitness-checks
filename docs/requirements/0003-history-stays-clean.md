# 0003 History Stays Clean

## Why

I read a branch's history and release notes and see what each change was, who helped, and which issue it served. Each issue is linked once, and each change has one dated entry.

## Measurement

commits and entry files that pass every history rule

-

commits and entry files checked

## Requirements

```gherkin
SCENARIO: 0003.1
  GIVEN Developer has commit subjects with and without a type and scope
  WHEN the commit checks read the subject line
  THEN only subjects with a listed type and scope pass
    AND each failure quotes the subject

SCENARIO: 0003.2
  GIVEN Developer has commit messages with correct and misspelled `Plan` references
  WHEN the commit checks read the Plan trailer
  THEN every misspelled trailer fails with the fix
    AND prose and other repos and absent trailers pass

SCENARIO: 0003.4
  GIVEN Developer has a branch and open PR that already link some issues
  WHEN the commit checks look for repeated issue links
  THEN a second link to the same issue fails naming the earlier commit or PR

SCENARIO: 0003.5
  GIVEN Developer has changelog entry files with good and bad names and headings
  WHEN the changelog checks read each entry file
  THEN each file needs one dated heading matching its name
    AND upstream `CHANGELOG.md` history passes as it is

SCENARIO: 0003.6
  GIVEN Developer has changelog entry files with different bullet counts and lengths and types
  WHEN the changelog checks read the bullets
  THEN an entry with fewer than three or more than five bullets fails naming it
    AND so does a long bullet or an untyped bullet

SCENARIO: 0003.7
  GIVEN Developer has staged commits with and without a new or edited entry file
  WHEN the changelog checks read the staged commit
  THEN the branch's one entry must name the change
    AND published entries stay as released

SCENARIO: 0003.8
  GIVEN Developer has release manifests and version files that agree or differ
  WHEN the changelog checks read the release versions
  THEN a mismatched or missing version fails naming it
    AND repos with no manifest pass
```
