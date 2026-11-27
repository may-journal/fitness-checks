# 0002 Comments Stay Short And Plain

## Why

I read code comments that are short and plain, and they say why. They hold no dead code or filler.

## Measurement

comments that pass every comment rule

-

comments the check judges

## Requirements

1. Comments stay short
   1. GIVEN Developer has code comments in source files
      1. WHEN a comment sentence runs long or a paragraph holds too many sentences
         1. THEN `prose-budget` fails naming the line
      2. WHEN a sentence starts with So
         1. THEN it fails asking to join the sentence before it or drop it
      3. WHEN a sentence has three words or fewer
         1. THEN it fails with the `comment-fragment` rule
      4. WHEN a comment holds three sentences
         1. THEN it fails with the `comment-one-statement` rule
      5. WHEN a plain comment is short
         1. THEN it passes
2. A hook judges only changed comments
   1. GIVEN Developer stages a change beside comments that already break a rule
      1. WHEN the staged change adds a comment that breaks a rule
         1. THEN only that comment fails
         2. THEN a full scan fails both comments
3. Dead weight fails
   1. GIVEN Developer has comments that hold code or filler
      1. WHEN a comment holds commented out code
         1. THEN it fails with the `commented-out-code` rule
      2. WHEN a comment restates the next code line
         1. THEN it fails with the `restates-code` rule
      3. WHEN a comment holds an emoji or a banner
         1. THEN it fails with the `decoration` rule
      4. WHEN a comment narrates a change
         1. THEN it fails with the `change-narration` rule
      5. WHEN a comment holds a hedging word
         1. THEN it fails with the `hedging` rule
      6. WHEN a comment opens like a tutorial
         1. THEN it fails with the `tutorial-tone` rule
4. Similar words pass
   1. GIVEN Developer has comments with words that resemble filler
      1. WHEN a comment holds a hedge-like word or a tutorial-like phrase
         1. THEN it passes
5. Function comments stay proportionate
   1. GIVEN Developer has functions with comments
      1. WHEN a function body holds many comment words for few code tokens
         1. THEN it fails with the `overlong` rule
      2. WHEN a function holds one long comment or a fenced example
         1. THEN it passes
