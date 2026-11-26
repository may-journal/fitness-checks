# 0001 Changelog Entry Files Are Dated

## Why

I keep my changelog as dated entry files, one per change, so the folder lists in release order.

## Measurement

entry files holding one dated heading and named for its date

-

entry files

## Requirements

1. Entry files hold one dated heading
   1. GIVEN Developer has a repo
   2. GIVEN it may hold `changelog/` entry files
   3. GIVEN it may hold a `CHANGELOG.md`
      1. WHEN each entry file holds one `### yyyy.mm.dd.HHMM` heading and its bullets and that date in its name
         1. THEN it passes counting the entry files
         2. THEN front matter is optional
      2. WHEN an entry file holds two dated headings
         1. THEN it fails asking for one `### yyyy.mm.dd.HHMM` heading and its bullets
         2. THEN an undated heading fails the same way
      3. WHEN an entry file's name lacks its heading's date
         1. THEN it fails naming the file it should be `changelog/<yyyy-mm-dd-HHMM>-<change>.md`
         2. THEN a name that carries another date fails the same way
      4. WHEN the repo has no entry files
         1. THEN it passes with no files
         2. THEN a `CHANGELOG.md` makes no difference
      5. WHEN `CHANGELOG.md` keeps an upstream fork's undated history
         1. THEN it is not judged
         2. THEN a moved package version is not judged either
         3. THEN the entry files still pass
