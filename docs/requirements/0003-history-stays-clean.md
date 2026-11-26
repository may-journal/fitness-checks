# 0003 History Stays Clean

## Why

I read a branch's history and release notes and see what each change was, who helped, and which issue it served. Each issue is linked once.

## Measurement

commits and entry files that pass every history rule

-

commits and entry files checked

## Requirements

1. Subjects carry type and scope
   1. GIVEN Developer has hooks installed on a fork branch
      1. WHEN I commit a subject
         1. THEN a valid subject commits without judging inherited history
         2. THEN a subject with no scope fails quoting only that subject
2. Removed in 2026.11.26.0000: moved to the plan-trailer check's own requirement
3. Removed in 2026.11.26.0100: moved to the commit-attribution doc
4. Removed in 2026.11.26.0200: moved to the issue-link-once doc
5. Removed in 2026.11.25.0800: moved to the `changelog` doc
6. Removed in 2026.11.25.0800: moved to the `changelog-bullets` doc
7. Removed in 2026.11.25.0800: moved to the `changelog-updated` doc
8. Removed in 2026.11.26.0000: moved to the release-changelog check's own requirement
