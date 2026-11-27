# 0001 Each Acceptance Has One Storybook Page

## Why

I get each requirement acceptance shown by exactly one Storybook page named for it, so no page or acceptance drops out unnoticed.

## Measurement

acceptances shown by exactly one Storybook page

-

acceptances in Storybook apps

## Requirements

1. Storybook pages match acceptances
   1. GIVEN Developer has a Swift repo with requirement docs
      1. WHEN no Swift file imports a Storybook module or declares a `BookPage`
         1. THEN it passes with a note that no Storybook is in the repo
      2. WHEN a Swift file imports a Storybook module or declares a `BookPage`
         1. THEN the repo is judged even with no pages yet
      3. WHEN each acceptance has one page named with its ID
         1. THEN it passes
         2. THEN it reads every Swift file
      4. WHEN a page has no name or a name that is not a plain literal
         1. THEN it fails naming the file
         2. THEN it fails naming the line
         3. THEN it says how to name the page
      5. WHEN a page name does not start with a defined acceptance ID or starts with a chained step
         1. THEN it fails naming the file that holds the page
         2. THEN it fails naming the line that holds the name
         3. THEN it names the page name to use
      6. WHEN two pages repeat a name or an acceptance ID
         1. THEN it fails naming both places
         2. THEN a `maxPagesPerAcceptance` setting allows more pages per acceptance
      7. WHEN an acceptance has no page in a Storybook app
         1. THEN it fails naming the doc line
         2. THEN it names the page name to add
