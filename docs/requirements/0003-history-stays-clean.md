# 0003 History Stays Clean

## Why

I read a branch's history and release notes and see what each change was, who helped, and which issue it served. Each issue is linked once, and each change has one dated entry.

## Measurement

commits and entry files that pass every history rule

-

commits and entry files checked

## Requirements

1. Subjects carry type and scope
   1. GIVEN Developer has commit subjects with and without a type
      1. WHEN the commit checks read the subject line
         1. THEN only subjects with a listed type and scope pass
         2. THEN each failure quotes the subject
2. Removed in 2026.11.26.0000: moved to the plan-trailer check's own requirement
3. Removed in 2026.11.26.0100: moved to the commit-attribution doc
4. Issue links are not repeated
   1. GIVEN Developer has a branch and PR linking issues
      1. WHEN the commit checks look for repeated issue links
         1. THEN a second link fails naming the earlier commit or PR
5. Changelog files are dated
   1. GIVEN Developer has changelog entry files with good and bad names and headings
      1. WHEN the changelog checks read each entry file
         1. THEN it needs one dated heading matching its name
         2. THEN `CHANGELOG.md` history passes
6. Changelog bullets are short
   1. GIVEN Developer has entry files with different bullets
      1. WHEN the changelog checks read the bullets
         1. THEN fewer than three or more than five fails
         2. THEN a long or untyped bullet fails
7. Each branch adds one entry
   1. GIVEN Developer has staged commits with and without an entry file
      1. WHEN the changelog checks read the staged commit
         1. THEN the branch's one entry must name the change
         2. THEN published entries stay as released
8. Removed in 2026.11.26.0000: moved to the release-changelog check's own requirement
