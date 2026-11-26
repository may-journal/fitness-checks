# 0001 Java Code Passes Checkstyle

## Why

I get Checkstyle judging my Java code with the plugin my build declares and a config that caps complexity.

## Measurement

repos whose Java code passes Checkstyle

-

repos with a Maven or Gradle build

## Requirements

1. Findings fail
   1. GIVEN Developer has a repo with Java code and a Maven or Gradle build that declares Checkstyle
      1. WHEN Checkstyle finds a violation or a warning
         1. THEN it fails naming each as `file:line[:col] - message (Rule)`
         2. THEN Gradle runs `checkstyleMain` and `checkstyleTest` with `--rerun`
      2. WHEN Checkstyle finds no violation and no warning
         1. THEN it passes and counts the Java files
         2. THEN a config in a test source set's folder is left out as a test input
2. The build declares Checkstyle and a complexity ceiling
   1. GIVEN Developer has a repo with Java code and a Maven build
      1. WHEN the build does not declare Checkstyle
         1. THEN it fails with the plugin snippet to add
      2. WHEN the repo has no Checkstyle config
         1. THEN it fails asking for `config/checkstyle/checkstyle.xml`
         2. THEN it gives the document type line Maven needs
      3. WHEN a tracked Checkstyle config has no `CyclomaticComplexity` module
         1. THEN it fails naming the config and the module to add
      4. WHEN the config the hint gave is added
         1. THEN Maven runs with it and the check passes
3. A repo with no build passes
   1. GIVEN Developer has a repo with no Maven or Gradle build
      1. WHEN the check runs
         1. THEN it passes with no files
