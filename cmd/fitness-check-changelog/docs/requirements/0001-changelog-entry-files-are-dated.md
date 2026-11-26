# 0001 Changelog Entry Files Are Dated

## Why

I keep my changelog as dated entry files, one per change, so the folder lists in release order.

## Measurement

entry files holding one dated heading and named for its date

-

entry files

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has a repo with or without `changelog/` entry files AND a `CHANGELOG.md`
  WHEN each entry file holds one `### yyyy.mm.dd.HHMM` heading
    AND each entry file holds its bullets
    AND each entry file has that date in its name
  THEN it passes
    AND it counts the entry files
    AND front matter may be present or absent
  WHEN an entry file holds two dated headings
    AND an entry file holds an undated heading
  THEN it fails asking for one `### yyyy.mm.dd.HHMM` heading and its bullets
  WHEN an entry file's name lacks its heading's date
    AND an entry file's name carries another date
  THEN it fails naming the file it should be
    AND it gives `changelog/<yyyy-mm-dd-HHMM>-<change>.md`
  WHEN the repo has no entry files
    AND `CHANGELOG.md` may be present or absent
  THEN it passes with no files
  WHEN `CHANGELOG.md` keeps an upstream fork's undated history
    AND the package version moves
  THEN neither is judged
    AND the entry files still pass
```
