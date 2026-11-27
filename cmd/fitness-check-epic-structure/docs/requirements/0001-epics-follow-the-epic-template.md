# 0001 Epics Follow The Epic Template

## Why

I get each epic issue body judged against the Epic template, so its list of issues reads the same everywhere.

## Measurement

epic bodies that follow the template

-

epic bodies

## Requirements

1. Epics follow the template
   1. GIVEN Developer has an epic issue body that reaches the check by `--body-file` or `--body`
      1. WHEN the body holds a one-line pitch plus a Background section plus a list under What needs to happen
         1. THEN it passes
      2. WHEN the body lacks the pitch or one of the sections
         1. THEN it fails naming what is missing
      3. WHEN the body holds another `##` section
         1. THEN it fails naming the section
      4. WHEN the repo has its own Epic template
         1. THEN that template judges the sections
      5. WHEN fenced code holds a `##` heading or a list item
         1. THEN it never counts
2. The list holds plain references in run order
   1. GIVEN Developer has an epic issue body with a list under What needs to happen
      1. WHEN an item is a checkbox or holds an emoji or status word or is no plain reference or precedes its `after` target
         1. THEN it fails naming the item
      2. WHEN a reference is listed twice
         1. THEN it fails naming the reference
      3. WHEN the body holds a table
         1. THEN it fails
3. No body means nothing to check
   1. GIVEN Developer has no epic body
      1. WHEN the check runs with no input
         1. THEN it passes with no files checked
