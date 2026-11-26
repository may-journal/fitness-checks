# 0001 Plans Follow The Template

## Why

I read every Plan in one shape: a pitch and its background and its checklist.

## Measurement

Plan bodies that follow the template

-

Plan bodies

## Requirements

1. A Plan follows the template
   1. GIVEN Developer has a Plan body
      1. WHEN it holds a pitch plus a Background section plus a checklist
         1. THEN `plan-structure` passes it
      2. WHEN it has no blockquote pitch
         1. THEN it fails
         2. THEN it asks for a one-line pitch
      3. WHEN its pitch is still the template placeholder
         1. THEN it fails
         2. THEN it asks for the real pitch
      4. WHEN it has no Background section
         1. THEN it fails naming the missing section
      5. WHEN its only checkbox sits inside code
         1. THEN it fails because the checklist has no item
      6. WHEN it has an extra section
         1. THEN it fails naming the unexpected section
      7. WHEN a section is named for open questions
         1. THEN it fails
         2. THEN it asks to turn unknowns into checklist steps
2. A missing Plan body passes
   1. GIVEN Developer runs a check with no Plan body
      1. WHEN `plan-structure` runs
         1. THEN it passes
         2. THEN it judges zero files
