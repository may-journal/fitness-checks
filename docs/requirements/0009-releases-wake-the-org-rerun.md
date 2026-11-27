# 0009 Releases Wake The Org Rerun

## Why

A check release changes what a pull request must pass, so each open pull request runs the fitness suite again. The org's .github repo does that work once it hears of the release.

## Measurement

check releases announced to the .github repo
-
check releases

## Requirements

1. Announces a release
   1. GIVEN Developer has a release that published a check
      1. WHEN the release workflow finishes
         1. THEN it sends the `fitness-check-released` dispatch to the .github repo
2. Stays quiet with no release
   1. GIVEN Developer has a push that released no check
      1. WHEN the release workflow finishes
         1. THEN it sends no dispatch
