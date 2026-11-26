# 0001 Epics Follow The Epic Template

## Why

I get each epic issue body judged against the Epic template, so its list of issues reads the same everywhere.

## Measurement

epic bodies that follow the template

-

epic bodies

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has an epic issue body whose text reaches the check by `--body-file` or `--body` or stdin
  WHEN the body holds a one-line pitch and a Background section and a list of plain issue references under What needs to happen
  THEN it passes
  WHEN the body lacks one of those or holds another `##` section
  THEN it fails naming the section
  WHEN a list item is a plain `#N` or `org/repo#N` reference with an optional short note after a comma
  THEN it passes
  WHEN a list item is a checkbox or holds an emoji or a status word or is not a plain reference
  THEN it fails naming the item
  WHEN a reference is listed twice or an item comes before an item its `after` note names
  THEN it fails naming the item

SCENARIO: 0001.2
  GIVEN Developer has an epic issue body with a table or fenced code or its own template or no text
  WHEN the body holds a table
  THEN it fails
  WHEN fenced code holds a `##` heading or a list item
  THEN it never counts
  WHEN the repo has its own Epic template
  THEN that template judges the sections
  WHEN no body reaches the check
  THEN it passes with no files checked
```
