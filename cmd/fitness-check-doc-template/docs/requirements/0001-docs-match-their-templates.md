# 0001 Docs Match Their Templates

## Why

I keep every markdown doc in the shape of the template that governs its folder.

## Measurement

markdown docs that match their templates

-

markdown docs governed by a template

## Requirements

```gherkin
SCENARIO: 0001.1
  GIVEN Developer has markdown docs beside `template.md` or `*.template.md` files or a `templates` folder
  WHEN a doc has a `##` section its template lacks or lacks one its template has
  THEN it fails naming the `Unexpected` or `Missing` section
  WHEN a template section has no guiding `<!-- comment -->`
  THEN it fails naming the section
  WHEN a GitHub PR or issue template sits under `.github`
  THEN it gets only the comment rule and governs no file
  WHEN docs sit in nested folders with templates at several levels
  THEN the first folder walking up with a governing template decides
  WHEN a folder holds several templates and a doc matches none
  THEN one error names all of them
  WHEN a template in a `templates` folder shares a neighbor's name
  THEN it holds that neighbor alone to its sections and legends and diagram types and callout tables
  WHEN a view sits under `docs/architecture/` or `architecture/`
  THEN it follows the shipped template unless the repo's own copy governs as a named template
  WHEN a C4 or class diagram uses a kind or style the template does not
  THEN it fails naming the element and the palette it needs
  WHEN a repo has no templates
  THEN it passes
```
