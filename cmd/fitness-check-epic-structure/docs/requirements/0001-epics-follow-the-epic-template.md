# 0001 Epics Follow The Epic Template

## Why

I get each epic issue body judged against the Epic template, so its list of issues reads the same everywhere.

## Measurement

epic bodies that follow the template

-

epic bodies

## Requirements

- 0001.1
  - Given Developer has an epic issue body whose text reaches the check by `--body-file` or `--body` or stdin
    1. When the body holds a one-line pitch and a Background section and a list of plain issue references under What needs to happen
       - Then it passes
    2. When the body lacks one of those or holds another `##` section
       - Then it fails naming the section
    3. When a list item is a plain `#N` or `org/repo#N` reference with an optional short note after a comma
       - Then it passes
    4. When a list item is a checkbox or holds an emoji or a status word or is not a plain reference
       - Then it fails naming the item
    5. When a reference is listed twice or an item comes before an item its `after` note names
       - Then it fails naming the item
- 0001.2
  - Given Developer has an epic issue body with a table or fenced code or its own template or no text
    1. When the body holds a table
       - Then it fails
    2. When fenced code holds a `##` heading or a list item
       - Then it never counts
    3. When the repo has its own Epic template
       - Then that template judges the sections
    4. When no body reaches the check
       - Then it passes with no files checked
