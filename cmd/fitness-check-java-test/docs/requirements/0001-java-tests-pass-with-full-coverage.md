# 0001 Java Tests Pass With Full Coverage

## Why

I get my Java tests run once per build root and judged by their own JUnit reports and JaCoCo coverage, in any repo.

## Measurement

build roots whose Java tests pass with full coverage

-

build roots with Java code

## Requirements

1. Java test failures and coverage shortfalls are named
   1. GIVEN Developer has a repo with a Maven or Gradle build
   2. GIVEN the build runs JUnit tests with JaCoCo
      1. WHEN the repo has no Maven or Gradle build
         1. THEN it passes with no files checked
      2. WHEN a test fails
         1. THEN it fails naming `Class.method failed: message` read from the JUnit XML
         2. THEN a report cut short names itself
      3. WHEN a JaCoCo line total is below the goal
         1. THEN it fails with `report: lines covered/total covered; the goal is 100%`
         2. THEN a report that counts no lines fails alike
      4. WHEN a build writes an aggregate report
         1. THEN it replaces the module reports
         2. THEN a module it leaves uncovered still fails
         3. THEN a module of tests alone is never judged
      5. WHEN a build root with product code runs no test
         1. THEN it fails naming the build root
         2. THEN test files that ran nothing get a framework hint
      6. WHEN a framework needs Docker and Docker is down
         1. THEN it fails before the run asking to start Docker
         2. THEN a hung Docker fails fast alike
      7. WHEN a build declares no JaCoCo
         1. THEN it fails with the plugin snippet to add
         2. THEN a build with no XML report gets the snippet that enables it
2. Each build runs once with the Gradle build cache
   1. GIVEN Developer has a Maven or Gradle build with passing tests
      1. WHEN the build is Maven
         1. THEN it runs `verify` once
         2. THEN it reads the Surefire and Failsafe reports
      2. WHEN the build is Gradle
         1. THEN it cleans and reruns the tests and the JaCoCo report
         2. THEN it runs with `--build-cache`
      3. WHEN the build cache restores the reports of unchanged tests
         1. THEN the restored reports count
      4. WHEN `org.gradle.caching` is not `true` in the build
         1. THEN Gradle runs the tests without the build cache
         2. THEN a value in the Gradle user home counts the same
