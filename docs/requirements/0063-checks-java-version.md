# 0063 Checks Java Version

## Why

I learn at once when my repo pins no JDK or my JDK is older than its pin, before a build fails strangely.

## Measurement

Java repos whose running JDK satisfies every pin
-
Java repos checked

## Requirements

- 0063.1
    - Given Developer has a build that pins no JDK
        - When fitness runs java-version
            - Then it fails with the snippets that add a pin
- 0063.2
    - Given Developer has a version file at or below the running JDK
        - When fitness runs java-version
            - Then it passes
- 0063.3
    - Given Developer has a build release above the running JDK
        - When fitness runs java-version
            - Then it fails naming the file and the release it needs
- 0063.4
    - Given Developer has a Maven release set through a property, in the pom or a parent it inherits
        - When fitness runs java-version
            - Then it judges the property's release
- 0063.5
    - Given Developer has a pin that names no release
        - When fitness runs java-version
            - Then it fails naming the pin
- 0063.6
    - Given Developer has a build and no running JDK
        - When fitness runs java-version
            - Then it fails telling me to install a JDK
- 0063.7
    - Given Developer has a repo with no Maven or Gradle build
        - When fitness runs java-version
            - Then it passes with zero files
- 0063.8
    - Given Developer has a Gradle toolchain no installed JDK matches, and no resolver
        - When fitness runs java-version
            - Then it fails naming the JDK to install
