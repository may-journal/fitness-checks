# 0001 Java Formatting Passes Spotless

## Why

I get Spotless judging my Java formatting with my build's own formatter, in any repo.

## Measurement

repos whose Java builds pass Spotless

-

repos with Maven or Gradle builds

## Requirements

1. Java builds pass Spotless
   1. GIVEN Developer has a repo with or without a Maven or Gradle build
   2. GIVEN a build may declare Spotless
      1. WHEN Spotless would rewrite a file
         1. THEN it fails naming each as `file is not formatted; run spotless:apply`
      2. WHEN the builds are clean
         1. THEN it passes counting the Java files
         2. THEN it runs once from each build root
      3. WHEN a build ships its own wrapper or uses Gradle
         1. THEN the wrapper runs instead of the tool on PATH
         2. THEN Gradle runs `spotlessCheck` without a daemon
      4. WHEN a build does not declare Spotless
         1. THEN it fails with the plugin snippet to add to that build
      5. WHEN the repo has no Maven or Gradle build
         1. THEN it passes with no files
      6. WHEN a failing build names no file
         1. THEN it reports the build's own error lines
      7. WHEN no JDK runs
         1. THEN it fails with an install hint
