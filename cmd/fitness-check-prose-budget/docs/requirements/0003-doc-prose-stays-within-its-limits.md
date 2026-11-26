# 0003 Doc Prose Stays Within Its Limits

## Why

I read docs and bodies whose sentences and paragraphs and sections and lists stay short, with limits my repo can set.

## Measurement

docs and bodies that pass every prose limit

-

docs and bodies the check judges

## Requirements

1. Doc prose stays within its limits
   1. GIVEN Developer has Markdown files and bodies
      1. WHEN a section holds more words than its limit
         1. THEN it fails naming the section and its word count
      2. WHEN `CHANGELOG.md` holds a long list item
         1. THEN it fails like any other doc
      3. WHEN a body read with `--body-file` holds a long sentence
         1. THEN it fails naming the sentence length
      4. WHEN a list holds more items than its level allows
         1. THEN it fails naming the level
      5. WHEN a sentence runs long or a paragraph holds too many sentences
         1. THEN it fails naming the count and the limit
      6. WHEN fenced code or a Measurement dash holds long text
         1. THEN it passes
      7. WHEN a ticked Plan task holds exactly the word limit
         1. THEN it passes
2. A repo sets its own limits
   1. GIVEN Developer sets prose keys in `.fitnessrc.json`
      1. WHEN a repo lowers a prose limit
         1. THEN a sentence past it fails
      2. WHEN a repo sets the comment word numbers
         1. THEN the overlong rule follows them
