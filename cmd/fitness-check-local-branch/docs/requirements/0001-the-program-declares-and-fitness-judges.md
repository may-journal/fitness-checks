# 0001 The Program Declares And Fitness Judges

## Why

I get local branch judgments from fitness itself and never from a program that cannot make them.

## Measurement

runs of the program that declare or refuse as stated

-

runs of the program

## Requirements

1. Local branch declares itself and refuses to judge alone
   1. GIVEN Developer has the `local-branch` program built
      1. WHEN it is asked to describe itself
         1. THEN it declares the `local-branch` check
         2. THEN it declares the scope `github`
         3. THEN it declares protocol 1
      2. WHEN it is run alone on a folder
         1. THEN it exits 3 saying it runs only inside fitness
