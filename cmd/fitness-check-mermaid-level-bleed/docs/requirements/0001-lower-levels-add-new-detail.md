# 0001 Lower Levels Add New Detail

## Why

I read each architecture level for what it adds. A level never restates the one above.

## Measurement

Architecture level files that add detail beyond the level above

-

Architecture level files

## Requirements

1. Levels never repeat the one above
   1. GIVEN Developer has numbered architecture level files
   2. GIVEN each has callout tables
      1. WHEN a level repeats a description of the level above
         1. THEN it fails naming the file
         2. THEN it names the text
         3. THEN it names the level
      2. WHEN the repeat differs only in case or spacing
         1. THEN it is reported once
      3. WHEN every description is distinct or blank
         1. THEN the levels pass
      4. WHEN level one shares text with level three
         1. THEN they pass
      5. WHEN level two shares text with level ten
         1. THEN they compare as adjacent
         2. THEN the repeat fails
2. Only description cells are compared
   1. GIVEN Developer has Markdown files in the numbered level layout
   2. GIVEN some files sit outside it
      1. WHEN Markdown sits outside numbered level files in `architecture`
         1. THEN the check passes with no files
      2. WHEN a table has a `Description` header in any case
         1. THEN only that column is compared
      3. WHEN a table has no `Description` header
         1. THEN only its second column is compared
