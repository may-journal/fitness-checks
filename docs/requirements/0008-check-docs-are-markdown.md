# 0008 Check Docs Are Markdown

## Why

I edit each check's rule as Markdown that GitHub renders and the Markdown checks lint, not as Go comments.

## Measurement

checks whose rule is a README

-

checks

## Requirements

1. The doc move keeps Markdown
   1. GIVEN Developer has a check whose rule is a doc.go comment
   2. GIVEN the comment has a heading and list and link and code and spell checker line
      1. WHEN the doc move script runs on its folder
         1. THEN its README holds them as Markdown
         2. THEN its doc.go is one line that names and embeds the README
         3. THEN the described summary stays
2. A doc with extras fails
   1. GIVEN Developer has a doc.go that holds more than its comment
      1. WHEN the doc move script runs on its folder
         1. THEN it fails naming the folder
         2. THEN it leaves both files as they were
3. Every check is moved
   1. GIVEN Developer has every check moved to Markdown
      1. WHEN the README test reads each folder
         1. THEN each doc.go is one line naming its README and embedding it
         2. THEN the repo README alone links the architecture page
4. Summaries read cleanly
   1. GIVEN Developer has a README with a parenthesized link URL
   2. GIVEN another README has a first block that is not prose
      1. WHEN the check describes itself
         1. THEN the summary has no stray parenthesis
         2. THEN the summary is empty without prose
5. Accented capitals survive
   1. GIVEN Developer has a doc.go comment that opens with a multibyte letter
      1. WHEN the doc move script runs on its folder
         1. THEN the README opens with that letter as a whole capital
6. READMEs match describe
   1. GIVEN Developer has each check's README
      1. WHEN the README test reads it against the check's `--describe`
         1. THEN the README names only config keys its check declares besides check names
         2. THEN it says to run `fitness explain` when it declares keys
