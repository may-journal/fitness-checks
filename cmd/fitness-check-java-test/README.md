# java-test

Runs Java tests once per build root and judges their JUnit XML results and [JaCoCo](https://www.jacoco.org/jacoco/) coverage. The Java language module owns the build run and the judgment. It is on by default.

## Behavior

[Requirement 0001](docs/requirements/0001-java-tests-pass-with-full-coverage.md)

## Notes

The [Java module](https://github.com/may-journal/fitness-runner/blob/main/go/pkg/lang/java/java.go) says what a build root is and which wrapper or tool runs it. A Java repo with no running JDK fails with an install hint.

## Reports

Only reports the run itself wrote count. The JUnit XML comes from `target/surefire-reports`, `target/failsafe-reports` and Gradle's `build/test-results/<task>`.

Coverage reads each report's own totals, never a sum of its classes. An aggregate report, Maven's `jacoco-aggregate` or Gradle's `testCodeCoverageReport`, replaces the module reports.

## Frameworks

The [frameworks](https://github.com/may-journal/fitness-runner/blob/main/go/pkg/lang/java/tests/frameworks) folder holds one commented file per framework, spotted from build files and imports. Each says how its tests reach the JUnit XML. A framework that needs a running service, as Testcontainers needs Docker, fails before the run when the service is down.

## Example

```
fitness java-test
```

It reads no `.fitnessrc.json` key.

<!-- cspell:ignore jacoco -->
