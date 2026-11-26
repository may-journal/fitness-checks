// cspell:ignore mvnw gradlew

// Package requirements proves the java-test check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo, with a JDK 21 and the build tools answering saved responses.
package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain lets this test binary stand in for the JDK and the build tools.
func TestMain(m *testing.M) {
	acceptance.Main(m, "java", "mvn", "mvnw", "gradle", "gradlew", "docker")
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	mavenRepo := func(t *testing.T) string {
		return acceptance.JavaRepo(t, map[string]string{"pom.xml": acceptance.JavaPom(javaTestMaven)})
	}
	mavenWith := func(t *testing.T, run acceptance.Response) (string, int) {
		_, out, code := acceptance.JavaRun(t, mavenRepo(t), "java-test", map[string][]acceptance.Response{"mvn": {run}})
		return out, code
	}
	t.Run("0001.1.2.1.1 it passes with no files checked", func(t *testing.T) {
		t.Parallel()
		_, out, code := acceptance.JavaRun(t, acceptance.Example(t, "happyRepo", nil), "java-test", nil)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "java-test")
	})
	t.Run("0001.1.2.2.1 it fails naming `Class.method failed: message` read from the JUnit XML", func(t *testing.T) {
		t.Parallel()
		out, code := mavenWith(t, javaTestMavenRun(map[string]string{
			"target/surefire-reports/TEST-example.AppTest.xml": javaTestFailed,
			"target/site/jacoco/jacoco.xml":                    javaTestFull,
		}))
		acceptance.Sees(t, out, code, 1, "example.AppTest.subtracts failed: expected: <1> but was: <2>")
		if strings.Contains(acceptance.Read(out), "example.AppTest.adds") {
			t.Errorf("only the failing test may be named:\n%s", out)
		}
	})
	t.Run("0001.1.2.2.2 a report cut short names itself", func(t *testing.T) {
		t.Parallel()
		out, code := mavenWith(t, javaTestMavenRun(map[string]string{"target/surefire-reports/TEST-example.AppTest.xml": "<testsuite><testcase"}))
		acceptance.Sees(t, out, code, 1, "target/surefire-reports/TEST-example.AppTest.xml is not a JUnit XML report; rerun the build")
	})
	t.Run("0001.1.2.3.1 it fails with `report: lines covered/total covered; the goal is 100%`", func(t *testing.T) {
		t.Parallel()
		out, code := mavenWith(t, javaTestMavenRun(map[string]string{
			"target/surefire-reports/TEST-example.AppTest.xml": javaTestPassed,
			"target/site/jacoco/jacoco.xml":                    javaTestCoverage("4", "1"),
		}))
		acceptance.Sees(t, out, code, 1, "target/site/jacoco/jacoco.xml: lines 4/5 covered; the goal is 100%")
	})
	t.Run("0001.1.2.3.2 a report that counts no lines fails alike", func(t *testing.T) {
		t.Parallel()
		out, code := mavenWith(t, javaTestMavenRun(map[string]string{
			"target/surefire-reports/TEST-example.AppTest.xml": javaTestPassed,
			"target/site/jacoco/jacoco.xml":                    `<report name="app"><counter type="INSTRUCTION" missed="0" covered="9"/></report>`,
		}))
		acceptance.Sees(t, out, code, 1, "target/site/jacoco/jacoco.xml counts no lines; compile with debug information")
	})
	files := map[string]string{
		"pom.xml":                      acceptance.JavaParent(javaTestMaven, "core", "extra", "it", "kt", "coverage"),
		"core/pom.xml":                 javaTestModule("core"),
		"core/src/main/java/Core.java": "class Core {}\n",
		// Maven's aggregate covers the modules it holds a group for, so a
		// source another group lists leaves extra uncovered.
		"extra/pom.xml":                 javaTestModule("extra"),
		"extra/src/main/java/Core.java": "class Core {}\n",
		// A module of tests alone has no code of its own to cover.
		"it/pom.xml":                           javaTestModule("it"),
		"it/src/test/java/example/FlowIT.java": "class FlowIT {}\n",
		// A module of other languages has code the aggregate may miss.
		"kt/pom.xml":                         javaTestModule("kt"),
		"kt/src/main/kotlin/example/Util.kt": "package example\nobject Util\n",
		"coverage/pom.xml":                   javaTestModule("coverage"),
	}
	aggregate := func(t *testing.T, report string) (string, int) {
		reports := map[string]string{
			"core/target/surefire-reports/TEST-example.CoreTest.xml": javaTestPassed,
			"core/target/site/jacoco/jacoco.xml":                     javaTestCoverage("2", "3"),
			"extra/target/site/jacoco/jacoco.xml":                    javaTestCoverage("3", "2"),
			"it/target/site/jacoco/jacoco.xml":                       `<report name="it"><counter type="INSTRUCTION" missed="0" covered="0"/></report>`,
			"kt/target/site/jacoco/jacoco.xml":                       javaTestCoverage("1", "4"),
			"coverage/target/site/jacoco-aggregate/jacoco.xml":       report,
		}
		_, out, code := acceptance.JavaRun(t, acceptance.JavaRepo(t, files), "java-test", map[string][]acceptance.Response{"mvn": {javaTestMavenRun(reports)}})
		return out, code
	}
	covered := javaTestAggregate("4", "1", `<group name="core">`+javaTestListed("", "Core.java")+`</group>`)
	t.Run("0001.1.2.4.1 it replaces the module reports", func(t *testing.T) {
		t.Parallel()
		out, code := aggregate(t, covered)
		acceptance.Sees(t, out, code, 1, "coverage/target/site/jacoco-aggregate/jacoco.xml: lines 4/5 covered")
		verifyNotJudged(t, out, "core/target/site/jacoco/jacoco.xml")
		// An aggregate that does not parse covers what is unknown, so it is
		// named alone.
		out, code = aggregate(t, "<report")
		acceptance.Sees(t, out, code, 1, "coverage/target/site/jacoco-aggregate/jacoco.xml is not a JaCoCo XML report; rerun the build")
		verifyNotJudged(t, out, "extra/target/site/jacoco/jacoco.xml")
	})
	t.Run("0001.1.2.4.2 a module it leaves uncovered still fails", func(t *testing.T) {
		t.Parallel()
		out, code := aggregate(t, covered)
		acceptance.Sees(t, out, code, 1, "extra/target/site/jacoco/jacoco.xml: lines 3/5 covered", "kt/target/site/jacoco/jacoco.xml: lines 1/5 covered")
		verifyGradleAggregate(t)
	})
	t.Run("0001.1.2.4.3 a module of tests alone is never judged", func(t *testing.T) {
		t.Parallel()
		out, code := aggregate(t, covered)
		acceptance.Sees(t, out, code, 1)
		verifyNotJudged(t, out, "it/target/site/jacoco/jacoco.xml")
	})
	t.Run("0001.1.2.5.1 it fails naming the build root", func(t *testing.T) {
		t.Parallel()
		out, code := mavenWith(t, javaTestMavenRun(nil))
		acceptance.Sees(t, out, code, 1, "the Maven build at . has product code such as "+acceptance.JavaSource+" and no tests; add tests under src/test/java")
	})
	t.Run("0001.1.2.5.2 test files that ran nothing get a framework hint", func(t *testing.T) {
		t.Parallel()
		karate := acceptance.JavaRepo(t, map[string]string{
			"pom.xml":                        acceptance.JavaPom(javaTestMaven + "<dependency><groupId>io.karatelabs</groupId><artifactId>karate-junit6</artifactId></dependency>"),
			"src/test/resources/app.feature": "Feature: app\n",
		})
		_, out, code := acceptance.JavaRun(t, karate, "java-test", map[string][]acceptance.Response{"mvn": {javaTestMavenRun(nil)}})
		acceptance.Sees(t, out, code, 1, "the Maven build at . has test files but ran no test", "Karate: Karate runs feature files only from a test class with a @Karate.Test method")
	})
	dockerRepo := func(t *testing.T) string {
		return acceptance.JavaRepo(t, map[string]string{
			"pom.xml":                         acceptance.JavaPom(javaTestMaven),
			"src/test/java/example/DbIT.java": "import org.testcontainers.containers.PostgreSQLContainer;\nclass DbIT {}\n",
		})
	}
	t.Run("0001.1.2.6.1 it fails before the run asking to start Docker", func(t *testing.T) {
		t.Parallel()
		down := acceptance.Response{Match: []string{"info"}, Stderr: "Cannot connect to the Docker daemon. Is the docker daemon running?\n", Exit: 1}
		rp, out, code := acceptance.JavaRun(t, dockerRepo(t), "java-test", map[string][]acceptance.Response{"mvn": {javaTestMavenRun(nil)}, "docker": {down}})
		acceptance.Sees(t, out, code, 1, "Testcontainers needs a running Docker; start Docker or set DOCKER_HOST, then rerun")
		if calls := rp.Calls("mvn"); len(calls) != 0 {
			t.Errorf("the build ran %q, want no run without Docker", calls)
		}
	})
	t.Run("0001.1.2.6.2 a hung Docker fails fast alike", func(t *testing.T) {
		t.Parallel()
		hung := acceptance.Response{Match: []string{"info"}, Wait: 60}
		start := time.Now()
		_, out, code := acceptance.JavaRun(t, dockerRepo(t), "java-test", map[string][]acceptance.Response{"mvn": {javaTestMavenRun(nil)}, "docker": {hung}})
		acceptance.Sees(t, out, code, 1, "Testcontainers needs a running Docker")
		if took := time.Since(start); took > 30*time.Second {
			t.Errorf("a hung Docker held the check for %s, want under 30s", took)
		}
	})
	t.Run("0001.1.2.7.1 it fails with the plugin snippet to add", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.JavaRepo(t, map[string]string{"pom.xml": acceptance.JavaPom("")})
		_, out, code := acceptance.JavaRun(t, repo, "java-test", map[string][]acceptance.Response{"mvn": {javaTestMavenRun(nil)}})
		acceptance.Sees(t, out, code, 1, "JaCoCo is not in the Maven build at .; add it:")
		acceptance.JavaSeesSnippet(t, out, "<artifactId>jacoco-maven-plugin</artifactId>")
	})
	t.Run("0001.1.2.7.2 a build with no XML report gets the snippet that enables it", func(t *testing.T) {
		t.Parallel()
		noXML := acceptance.Response{Match: []string{"jacocoTestReport"}, Stdout: "BUILD SUCCESSFUL in 3s\n", Files: map[string]string{
			"build/test-results/test/TEST-example.AppTest.xml": javaTestPassed,
			"build/reports/jacoco/test/html/index.html":        "<html></html>\n",
		}}
		_, out, code := acceptance.JavaRun(t, acceptance.JavaRepo(t, map[string]string{"build.gradle": javaTestGradle}), "java-test", map[string][]acceptance.Response{"gradle": {noXML}})
		acceptance.Sees(t, out, code, 1, "no JaCoCo XML report after the Gradle build at . ran its tests; add:")
		acceptance.JavaSeesSnippet(t, out, "xml.required = true")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	maven := func(t *testing.T, surefire, failsafe string) (acceptance.Replays, string, int) {
		run := javaTestMavenRun(map[string]string{
			"target/surefire-reports/TEST-example.AppTest.xml": surefire,
			"target/failsafe-reports/TEST-example.AppIT.xml":   failsafe,
			"target/site/jacoco/jacoco.xml":                    javaTestFull,
		})
		return acceptance.JavaRun(t, acceptance.JavaRepo(t, map[string]string{"pom.xml": acceptance.JavaPom(javaTestMaven)}), "java-test", map[string][]acceptance.Response{"mvn": {run}})
	}
	t.Run("0001.2.1.1.1 it runs `verify` once", func(t *testing.T) {
		t.Parallel()
		rp, out, code := maven(t, javaTestPassed, javaTestPassed)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.JavaCalledOnce(t, rp, "mvn")
	})
	t.Run("0001.2.1.1.2 it reads Surefire or Failsafe reports", func(t *testing.T) {
		t.Parallel()
		_, out, code := maven(t, javaTestFailed, strings.ReplaceAll(javaTestFailed, "AppTest", "AppIT"))
		acceptance.Sees(t, out, code, 1, "example.AppTest.subtracts failed", "example.AppIT.subtracts failed")
	})
	gradle := func(t *testing.T) (acceptance.Replays, string, int) {
		return acceptance.JavaRun(t, acceptance.JavaRepo(t, map[string]string{"build.gradle": javaTestGradle}), "java-test", map[string][]acceptance.Response{"gradle": {javaTestGradleRun}}, "GRADLE_USER_HOME="+t.TempDir())
	}
	// Clean rules rerun the tasks on every Gradle version; --rerun needs 7.6.
	// The build cache then restores the reports of unchanged tests.
	t.Run("0001.2.1.2.1 it cleans the build", func(t *testing.T) {
		t.Parallel()
		rp, out, code := gradle(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if !rp.Called("gradle", "cleanTest", "cleanJacocoTestReport") || rp.Called("gradle", "--rerun") {
			t.Errorf("gradle calls = %q, want clean rules and no --rerun", rp.Calls("gradle"))
		}
	})
	t.Run("0001.2.1.2.2 it reruns the tests", func(t *testing.T) {
		t.Parallel()
		rp, out, code := gradle(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if !rp.Called("gradle", "cleanTest test") {
			t.Errorf("gradle calls = %q, want cleanTest before test", rp.Calls("gradle"))
		}
	})
	t.Run("0001.2.1.2.3 it reruns JaCoCo", func(t *testing.T) {
		t.Parallel()
		rp, out, code := gradle(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if !rp.Called("gradle", "cleanJacocoTestReport jacocoTestReport") {
			t.Errorf("gradle calls = %q, want cleanJacocoTestReport before jacocoTestReport", rp.Calls("gradle"))
		}
	})
	t.Run("0001.2.1.2.4 it uses `--build-cache`", func(t *testing.T) {
		t.Parallel()
		rp, out, code := gradle(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.JavaCalledOnce(t, rp, "gradle")
		if !rp.Called("gradle", "--build-cache") {
			t.Errorf("gradle calls = %q, want --build-cache", rp.Calls("gradle"))
		}
	})
	t.Run("0001.2.1.3.1 the restored reports count", func(t *testing.T) {
		t.Parallel()
		// Reports the cache restores are written during the run, so they count.
		rp, out, code := acceptance.JavaRun(t, acceptance.JavaRepo(t, map[string]string{"build.gradle": javaTestGradle}), "java-test", map[string][]acceptance.Response{"gradle": {javaTestCached}}, "GRADLE_USER_HOME="+t.TempDir())
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		if !rp.Called("gradle", "--build-cache") {
			t.Errorf("gradle calls = %q, want --build-cache", rp.Calls("gradle"))
		}
	})
	t.Run("0001.2.1.4.1 Gradle runs the tests without the build cache", func(t *testing.T) {
		t.Parallel()
		// Gradle reads any value but true as off, from the build or the user home.
		verifyUncached(t, "org.gradle.caching = no\n", javaTestGradleRun, "GRADLE_USER_HOME="+t.TempDir())
	})
	t.Run("0001.2.1.4.2 a value in the Gradle user home counts the same", func(t *testing.T) {
		t.Parallel()
		user := t.TempDir()
		acceptance.MustDo(t, os.WriteFile(filepath.Join(user, "gradle.properties"), []byte("org.gradle.caching=False\n"), 0o644))
		verifyUncached(t, "org.gradle.caching=true\n", javaTestGradleRun, "GRADLE_USER_HOME="+user)
	})
}

// javaTestGradleRun is a Gradle run that passes with every report written.
var javaTestGradleRun = acceptance.Response{Match: []string{"test", "jacocoTestReport"}, Stdout: "BUILD SUCCESSFUL in 3s\n", Files: map[string]string{
	"build/test-results/test/TEST-example.AppTest.xml":          javaTestPassed,
	"build/test-results/integrationTest/TEST-example.AppIT.xml": javaTestPassed,
	"build/reports/jacoco/test/jacocoTestReport.xml":            javaTestFull,
}}

// cspell:ignore jacoco testcontainers karatelabs opentest sessioninfo Postgre

// javaTestMaven declares JaCoCo in a Maven build.
const javaTestMaven = "<plugin><groupId>org.jacoco</groupId><artifactId>jacoco-maven-plugin</artifactId></plugin>"

// javaTestGradle is a Gradle build that applies JaCoCo.
const javaTestGradle = "plugins { id 'java'; id 'jacoco' }\n"

// javaTestPassed is a JUnit XML report of one passing test.
const javaTestPassed = `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="AppTest" tests="1" failures="0" errors="0" skipped="0"><testcase name="adds" classname="example.AppTest" time="0.01"/></testsuite>
`

// javaTestFailed is a JUnit XML report of one passing and one failing test.
const javaTestFailed = `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="AppTest" tests="2" failures="1" errors="0" skipped="0">
<testcase name="adds" classname="example.AppTest" time="0.01"/>
<testcase name="subtracts" classname="example.AppTest" time="0.01"><failure message="expected: &lt;1&gt; but was: &lt;2&gt;" type="org.opentest4j.AssertionFailedError">org.opentest4j.AssertionFailedError: expected: &lt;1&gt; but was: &lt;2&gt;
	at example.AppTest.subtracts(AppTest.java:12)</failure></testcase>
</testsuite>
`

// javaTestCoverage is a JaCoCo XML report whose own totals hold covered
// and missed lines and cover every branch.
func javaTestCoverage(covered, missed string) string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><!DOCTYPE report PUBLIC "-//JACOCO//DTD Report 1.1//EN" "report.dtd">` +
		`<report name="app"><sessioninfo id="s" start="1" dump="2"/><package name="example"><counter type="LINE" missed="0" covered="1"/></package>` +
		`<counter type="LINE" missed="` + missed + `" covered="` + covered + `"/><counter type="BRANCH" missed="0" covered="2"/></report>
`
}

// javaTestFull is a JaCoCo report with every line and branch covered.
var javaTestFull = javaTestCoverage("5", "0")

// javaTestMavenRun is Maven's answer to verify: it writes files and passes,
// as it does past failing tests.
func javaTestMavenRun(files map[string]string) acceptance.Response {
	return acceptance.Response{Match: []string{"verify"}, Files: files, Stdout: "[INFO] BUILD SUCCESS\n"}
}

// verifyGradleAggregate proves Gradle's aggregate, which lists packages
// without module groups, covers a project when it lists one of its
// sources by declared package.
func verifyGradleAggregate(t *testing.T) {
	t.Helper()
	files := map[string]string{
		"settings.gradle":              "include 'core', 'extra', 'locked'\n",
		"build.gradle":                 "plugins { id 'java'; id 'jacoco-report-aggregation' }\n",
		"core/build.gradle":            javaTestGradle,
		"core/src/main/java/Core.java": "class Core {}\n",
		// A class the aggregate's filters leave out keeps core covered.
		"core/src/main/java/example/Port.java": "package example;\ninterface Port {}\n",
		// The same file name in another package is another source.
		"extra/build.gradle":                    javaTestGradle,
		"extra/src/main/java/example/Core.java": "/* a; b */\npackage example;\nclass Core {}\n",
		// A source that cannot be read is not listed.
		"locked/build.gradle":                    javaTestGradle,
		"locked/src/main/java/example/Core.java": "package example;\nclass Core {}\n",
	}
	repo := acceptance.JavaRepo(t, files)
	acceptance.MustDo(t, os.Chmod(filepath.Join(repo, "locked/src/main/java/example/Core.java"), 0))
	run := acceptance.Response{Match: []string{"testCodeCoverageReport"}, Stdout: "BUILD SUCCESSFUL in 3s\n", Files: map[string]string{
		"core/build/test-results/test/TEST-example.CoreTest.xml":                 javaTestPassed,
		"core/build/reports/jacoco/test/jacocoTestReport.xml":                    javaTestCoverage("2", "3"),
		"extra/build/reports/jacoco/test/jacocoTestReport.xml":                   javaTestCoverage("3", "2"),
		"locked/build/reports/jacoco/test/jacocoTestReport.xml":                  javaTestCoverage("1", "4"),
		"build/reports/jacoco/testCodeCoverageReport/testCodeCoverageReport.xml": javaTestAggregate("5", "0", javaTestListed("", "Core.java")),
	}}
	_, out, code := acceptance.JavaRun(t, repo, "java-test", map[string][]acceptance.Response{"gradle": {run}})
	acceptance.Sees(t, out, code, 1, "extra/build/reports/jacoco/test/jacocoTestReport.xml: lines 3/5 covered", "locked/build/reports/jacoco/test/jacocoTestReport.xml: lines 1/5 covered")
	verifyNotJudged(t, out, "core/build/reports/jacoco/test/jacocoTestReport.xml")
}

// javaTestModule is a Maven module's build file with its artifact id.
func javaTestModule(id string) string {
	return "<project><modelVersion>4.0.0</modelVersion><groupId>com.example</groupId><artifactId>" + id + "</artifactId><version>1</version></project>\n"
}

// javaTestAggregate is an aggregate JaCoCo report with the given line
// totals that lists what listed holds: groups, or packages.
func javaTestAggregate(covered, missed, listed string) string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><report name="coverage"><sessioninfo id="s" start="1" dump="2"/>` + listed +
		`<counter type="LINE" missed="` + missed + `" covered="` + covered + `"/></report>
`
}

// javaTestListed is a report's package listing one source file.
func javaTestListed(pkg, source string) string {
	return `<package name="` + pkg + `"><sourcefile name="` + source + `"/></package>`
}

// verifyNotJudged fails when the output names any of the reports.
func verifyNotJudged(t *testing.T, out string, reports ...string) {
	t.Helper()
	for _, r := range reports {
		if strings.Contains(acceptance.Read(out), r) {
			t.Errorf("%s may not be judged:\n%s", r, out)
		}
	}
}

// javaTestCached is a Gradle run that restores the tests and coverage of
// an unchanged build from its build cache, rewriting their reports.
var javaTestCached = acceptance.Response{Match: []string{"test", "jacocoTestReport"}, Stdout: "> Task :test FROM-CACHE\n> Task :jacocoTestReport FROM-CACHE\nBUILD SUCCESSFUL in 1s\n", Files: map[string]string{
	"build/test-results/test/TEST-example.AppTest.xml": javaTestPassed,
	"build/reports/jacoco/test/jacocoTestReport.xml":   javaTestFull,
}}

// verifyUncached asserts a Gradle build with props in its
// gradle.properties, run with env, reruns its tests without the build
// cache.
func verifyUncached(t *testing.T, props string, gradle acceptance.Response, env string) {
	t.Helper()
	repo := acceptance.JavaRepo(t, map[string]string{"build.gradle": javaTestGradle, "gradle.properties": props})
	rp, out, code := acceptance.JavaRun(t, repo, "java-test", map[string][]acceptance.Response{"gradle": {gradle}}, env)
	acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	if !rp.Called("gradle", "cleanTest test") || rp.Called("gradle", "--build-cache") {
		t.Errorf("gradle calls = %q, want the tests run without --build-cache", rp.Calls("gradle"))
	}
}
