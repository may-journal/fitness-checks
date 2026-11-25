# 0001 Changelog Entry Files Are Dated

## Why

I keep my changelog as dated entry files, one per change, so the folder lists in release order.

## Measurement

entry files holding one dated heading and named for its date

-

entry files

## Requirements

- 0001.1
  - Given Developer has a repo, with or without `changelog/` entry files and a `CHANGELOG.md`
    1. When each entry file holds one `### yyyy.mm.dd.HHMM` heading, its bullets, and that date in its name
       - Then it passes, counting the entry files, with or without front matter
    2. When an entry file holds two dated headings, or an undated one
       - Then it fails, asking for one `### yyyy.mm.dd.HHMM` heading and its bullets
    3. When an entry file's name lacks its heading's date, or carries another
       - Then it fails naming the file it should be, `changelog/<yyyy-mm-dd-HHMM>-<change>.md`
    4. When the repo has no entry files, with or without `CHANGELOG.md`
       - Then it passes with no files
    5. When `CHANGELOG.md` keeps an upstream fork's undated history, or the package version moves
       - Then neither is judged, and the entry files still pass
