# 0001 Docs Match Their Templates

## Why

I keep every markdown doc in the shape of the template that governs its folder.

## Measurement

markdown docs that match their templates

-

markdown docs governed by a template

## Requirements

1. Docs match their folder's template
   1. GIVEN Developer has markdown docs beside `template.md` or `*.template.md` files or a `templates` folder
      1. WHEN a doc has a `##` section its template lacks or lacks one its template has
         1. THEN it fails naming the `Unexpected` or `Missing` section
      2. WHEN docs sit in nested folders with templates at several levels
         1. THEN the first folder walking up with a governing template decides
      3. WHEN a doc matches none of the several templates in its folder
         1. THEN one error names all of them
      4. WHEN a GitHub PR or issue template sits under `.github`
         1. THEN it governs no file
2. Templates guide each section
   1. GIVEN Developer has templates with or without guiding comments
      1. WHEN a template section has no guiding `<!-- comment -->`
         1. THEN it fails naming the template
         2. THEN it fails naming the section
         3. THEN a GitHub PR or issue template gets only this rule
3. Architecture views follow their templates
   1. GIVEN Developer has views under `docs/architecture` beside the shipped templates
      1. WHEN a view lacks a section or legend or callout table its template has
         1. THEN it fails naming the line
         2. THEN it fails naming what the template has
      2. WHEN the repo keeps no copy of the templates
         1. THEN the shipped templates govern the views
      3. WHEN a C4 or class diagram uses a kind or style the template does not
         1. THEN it fails naming the element
         2. THEN it fails naming the palette it needs
