# 0003 History Stays Clean

## Why

I read a branch's history and release notes and see what each change was, who helped, and which issue it served. Each issue is linked once.

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
5. Removed in 2026.11.25.0800: moved to the `changelog` doc
6. Removed in 2026.11.25.0800: moved to the `changelog-bullets` doc
7. Removed in 2026.11.25.0800: moved to the `changelog-updated` doc
8. Removed in 2026.11.26.0000: moved to the release-changelog check's own requirement
