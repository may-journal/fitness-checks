# 0001 READMEs Show Their Requirements And API

## Why

I read a check's behavior in one place and its whole API in one example.

## Measurement

READMEs whose Behavior is one link and whose Example shows the API

-

READMEs the check judges

## Requirements

1. Behavior is one link
   1. GIVEN Developer has a README with a Behavior section
      1. WHEN the section holds one link to its folder's requirements
         1. THEN it passes
      2. WHEN the section holds more than that link
         1. THEN it fails naming the line
      3. WHEN the section holds no link
         1. THEN it fails naming the heading
2. The link is the README's own requirements
   1. GIVEN Developer has a README whose Behavior holds one link
      1. WHEN the link names a doc in its folder's `docs/requirements`
         1. THEN it passes
      2. WHEN the link names another folder
         1. THEN it fails naming the line
      3. WHEN the link names nothing that exists
         1. THEN it fails naming the line
3. A check README has an Example
   1. GIVEN Developer has the README of a check program
      1. WHEN it has no Example section
         1. THEN it fails naming the README
      2. WHEN its Example leaves out the command `fitness <name>`
         1. THEN it fails naming the command
4. The Example shows each declared flag
   1. GIVEN Developer has a check README whose program declares a flag
      1. WHEN its Example leaves the flag out
         1. THEN it fails naming the flag
      2. WHEN its Example shows the flag
         1. THEN it passes
5. The Example sets each declared key to its default
   1. GIVEN Developer has a check README whose program declares config keys
      1. WHEN its Example leaves a key out
         1. THEN it fails naming the key
      2. WHEN a `.fitnessrc.json` snippet sets each key to its default
         1. THEN it passes
6. Steps stay out of the README
   1. GIVEN Developer has a README
      1. WHEN a bulleted or numbered item opens with GIVEN or WHEN or THEN in any case
         1. THEN it fails naming the line
      2. WHEN the same words sit in a code fence or open a plain sentence
         1. THEN it passes
