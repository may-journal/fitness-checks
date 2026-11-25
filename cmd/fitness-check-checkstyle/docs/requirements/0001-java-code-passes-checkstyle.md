# 0001 Java Code Passes Checkstyle

## Why

I get Checkstyle judging my Java code with the plugin my build declares and a config that caps complexity.

## Measurement

repos whose Java code passes Checkstyle

-

repos with a Maven or Gradle build

## Requirements

- 0001.1
  - Given Developer has a repo with or without a Maven or Gradle build that declares Checkstyle
    1. When the build is Maven
       - Then it runs once from each build root with warnings counted as violations over main and test sources
    2. When the build is Gradle
       - Then it runs `checkstyleMain` and `checkstyleTest` with `--rerun`
    3. When Checkstyle finds no violation and no warning
       - Then it passes and counts the Java files
    4. When Checkstyle finds a violation or a warning
       - Then it fails naming each as `file:line[:col] - message (Rule)`
    5. When a run fails without naming a finding
       - Then it reports the build's own error lines
    6. When the repo has no Maven or Gradle build
       - Then it passes with no files and never runs a build
    7. When the repo has Java code and no JDK is running
       - Then it fails with the hint to install a JDK
    8. When the build does not declare Checkstyle
       - Then it fails with the plugin snippet to add
    9. When a tracked Checkstyle config has no `CyclomaticComplexity` module
       - Then it fails naming the config and the module to add
    10. When the repo has no Checkstyle config
        - Then it fails asking for `config/checkstyle/checkstyle.xml` with the document type line Maven needs
    11. When a config sits in a test source set's folder
        - Then it is left out as a test input
    12. When a Maven build names no config
        - Then Maven runs with the repo's only config or else `config/checkstyle/checkstyle.xml`
