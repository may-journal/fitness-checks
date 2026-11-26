# 0001 Java Builds Pin A JDK The Machine Has

## Why

I get my Java repo pinning its JDK and my machine's JDKs satisfying the pin, in any repo.

## Measurement

repos whose Java builds pin a JDK the machine satisfies

-

repos with Java builds

## Requirements

1. Pins are judged against the running JDK
   1. GIVEN Developer has a repo with a Maven or Gradle build
   2. GIVEN a JDK may be running
      1. WHEN the repo has no Java build
         1. THEN it passes with no files checked
      2. WHEN no pin exists anywhere
         1. THEN it fails with the snippets that add one
      3. WHEN a pin is at or below the running JDK
         1. THEN it passes since a newer JDK builds for an older release
      4. WHEN a pin is above the running JDK
         1. THEN it fails naming the file
         2. THEN it fails naming the release it needs
      5. WHEN a pin names no release
         1. THEN it fails quoting the pin
      6. WHEN a pom names no release
         1. THEN it takes one from its properties or its parents in the repo
         2. THEN it takes one from a parent in the local Maven repository
      7. WHEN there is no running JDK
         1. THEN it fails with an install hint
2. Gradle toolchains need an installed JDK
   1. GIVEN Developer has a Gradle build with a toolchain
      1. WHEN the toolchain asks for a release no folder holds
         1. THEN it fails naming the JDK to install
         2. THEN a release in Gradle's own download folder passes
         3. THEN the JDK Gradle runs on passes without starting Gradle
      2. WHEN `gradle.properties` lists JDK folders
         1. THEN they count in every properties form
         2. THEN a build in a folder reads its own `gradle.properties`
         3. THEN the Gradle user home overrides the build
      3. WHEN a toolchain sets a vendor
         1. THEN it needs an installed JDK from that vendor
         2. THEN a J9 implementation needs an installed J9 JDK
      4. WHEN a resolver sits in the settings script
         1. THEN the toolchain passes with no installed JDK
         2. THEN a resolver in a user home init script passes alike
         3. THEN an empty or commented-out resolver does not count
         4. THEN `org.gradle.java.installations.auto-download=false` exempts nothing
