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
	"github.com/may-journal/fitness-checks/req"
)

// TestMain lets this test binary stand in for the JDK and the build tools.
func TestMain(m *testing.M) {
	acceptance.Main(m, "java", "mvn", "mvnw", "gradle", "gradlew", "docker")
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a repo with no Maven or Gradle build passes with no files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.1.1", "0a671ccb", func(_ req.A) {
			_, out, code := acceptance.JavaRun(t, acceptance.Example(t, "happyRepo", nil), "java-test", nil)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, out, "java-test")
		})
	})
	t.Run("0001.1.2 a failing test is named from its JUnit report", func(t *testing.T) {
		t.Parallel()
		out, code := javaTestMavenOutcome(t, javaTestFailed, javaTestFull)
		// A report cut short names itself, never a missing test run.
		cutOut, cutCode := javaTestMavenOutcome(t, "<testsuite><testcase", "")
		req.Then(t, "0001.1.2.2.1", "9254cd81", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "example.AppTest.subtracts failed: expected: <1> but was: <2>")
			a.True(!strings.Contains(acceptance.Read(out), "example.AppTest.adds"), "only the failing test may be named:\n%s", out)
		})
		req.Then(t, "0001.1.2.2.2", "f82ec8a2", func(a req.A) {
			acceptance.Sees(t, cutOut, cutCode, 1, "target/surefire-reports/TEST-example.AppTest.xml is not a JUnit XML report; rerun the build")
		})
	})
	t.Run("0001.1.3 coverage below the goal fails naming the report", func(t *testing.T) {
		t.Parallel()
		out, code := javaTestMavenOutcome(t, javaTestPassed, javaTestCoverage("4", "1"))
		// A report without line counts can't show the goal is met.
		bareOut, bareCode := javaTestMavenOutcome(t, javaTestPassed, `<report name="app"><counter type="INSTRUCTION" missed="0" covered="9"/></report>`)
		req.Then(t, "0001.1.2.3.1", "8746ab20", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "target/site/jacoco/jacoco.xml: lines 4/5 covered; the goal is 100%")
		})
		req.Then(t, "0001.1.2.3.2", "ee5fea70", func(a req.A) {
			acceptance.Sees(t, bareOut, bareCode, 1, "target/site/jacoco/jacoco.xml counts no lines; compile with debug information")
		})
	})
	t.Run("0001.1.4 an aggregate report replaces the module reports", func(t *testing.T) {
		t.Parallel()
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
		reports := map[string]string{
			"core/target/surefire-reports/TEST-example.CoreTest.xml": javaTestPassed,
			"core/target/site/jacoco/jacoco.xml":                     javaTestCoverage("2", "3"),
			"extra/target/site/jacoco/jacoco.xml":                    javaTestCoverage("3", "2"),
			"it/target/site/jacoco/jacoco.xml":                       `<report name="it"><counter type="INSTRUCTION" missed="0" covered="0"/></report>`,
			"kt/target/site/jacoco/jacoco.xml":                       javaTestCoverage("1", "4"),
			"coverage/target/site/jacoco-aggregate/jacoco.xml":       javaTestAggregate("4", "1", `<group name="core">`+javaTestListed("", "Core.java")+`</group>`),
		}
		_, out, code := acceptance.JavaRun(t, acceptance.JavaRepo(t, files), "java-test", map[string][]acceptance.Response{"mvn": {javaTestMavenRun(reports)}})
		// An aggregate that does not parse covers what is unknown, so it is
		// named alone.
		reports["coverage/target/site/jacoco-aggregate/jacoco.xml"] = "<report"
		_, badOut, badCode := acceptance.JavaRun(t, acceptance.JavaRepo(t, files), "java-test", map[string][]acceptance.Response{"mvn": {javaTestMavenRun(reports)}})
		req.Then(t, "0001.1.2.4.1", "62f31101", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "coverage/target/site/jacoco-aggregate/jacoco.xml: lines 4/5 covered")
			javaTestNotJudged(t, out, "core/target/site/jacoco/jacoco.xml")
			acceptance.Sees(t, badOut, badCode, 1, "coverage/target/site/jacoco-aggregate/jacoco.xml is not a JaCoCo XML report; rerun the build")
			javaTestNotJudged(t, badOut, "extra/target/site/jacoco/jacoco.xml")
		})
		req.Then(t, "0001.1.2.4.2", "b32905dc", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "extra/target/site/jacoco/jacoco.xml: lines 3/5 covered", "kt/target/site/jacoco/jacoco.xml: lines 1/5 covered")
			javaTestGradleAggregate(t)
		})
		req.Then(t, "0001.1.2.4.3", "f0e54417", func(a req.A) {
			javaTestNotJudged(t, out, "it/target/site/jacoco/jacoco.xml")
		})
	})
	t.Run("0001.1.5 a build root with product code and no test fails", func(t *testing.T) {
		t.Parallel()
		_, out, code := acceptance.JavaRun(t, acceptance.JavaRepo(t, map[string]string{"pom.xml": acceptance.JavaPom(javaTestMaven)}), "java-test", map[string][]acceptance.Response{"mvn": {javaTestMavenRun(nil)}})
		karate := acceptance.JavaRepo(t, map[string]string{
			"pom.xml":                        acceptance.JavaPom(javaTestMaven + "<dependency><groupId>io.karatelabs</groupId><artifactId>karate-junit6</artifactId></dependency>"),
			"src/test/resources/app.feature": "Feature: app\n",
		})
		_, hintOut, hintCode := acceptance.JavaRun(t, karate, "java-test", map[string][]acceptance.Response{"mvn": {javaTestMavenRun(nil)}})
		req.Then(t, "0001.1.2.5.1", "8cd0b757", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "the Maven build at . has product code such as "+acceptance.JavaSource+" and no tests; add tests under src/test/java")
		})
		req.Then(t, "0001.1.2.5.2", "c4212b08", func(a req.A) {
			acceptance.Sees(t, hintOut, hintCode, 1, "the Maven build at . has test files but ran no test", "Karate: Karate runs feature files only from a test class with a @Karate.Test method")
		})
	})
	t.Run("0001.1.6 a framework needing Docker that is down fails fast", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.JavaRepo(t, map[string]string{
			"pom.xml":                         acceptance.JavaPom(javaTestMaven),
			"src/test/java/example/DbIT.java": "import org.testcontainers.containers.PostgreSQLContainer;\nclass DbIT {}\n",
		})
		down := acceptance.Response{Match: []string{"info"}, Stderr: "Cannot connect to the Docker daemon. Is the docker daemon running?\n", Exit: 1}
		rp, out, code := acceptance.JavaRun(t, repo, "java-test", map[string][]acceptance.Response{"mvn": {javaTestMavenRun(nil)}, "docker": {down}})
		// A hung Docker fails fast rather than holding the check to its limit.
		hung := acceptance.Response{Match: []string{"info"}, Wait: 60}
		start := time.Now()
		_, hungOut, hungCode := acceptance.JavaRun(t, repo, "java-test", map[string][]acceptance.Response{"mvn": {javaTestMavenRun(nil)}, "docker": {hung}})
		took := time.Since(start)
		req.Then(t, "0001.1.2.6.1", "5ff9679a", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "Testcontainers needs a running Docker; start Docker or set DOCKER_HOST, then rerun")
			calls := rp.Calls("mvn")
			a.True(len(calls) == 0, "the build ran %q, want no run without Docker", calls)
		})
		req.Then(t, "0001.1.2.6.2", "363de3c1", func(a req.A) {
			acceptance.Sees(t, hungOut, hungCode, 1, "Testcontainers needs a running Docker")
			a.True(took <= 30*time.Second, "a hung Docker held the check for %s, want under 30s", took)
		})
	})
	t.Run("0001.1.7 a build without JaCoCo reports gets the snippet to add", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.JavaRepo(t, map[string]string{"pom.xml": acceptance.JavaPom("")})
		_, out, code := acceptance.JavaRun(t, repo, "java-test", map[string][]acceptance.Response{"mvn": {javaTestMavenRun(nil)}})
		noXML := acceptance.Response{Match: []string{"jacocoTestReport"}, Stdout: "BUILD SUCCESSFUL in 3s\n", Files: map[string]string{
			"build/test-results/test/TEST-example.AppTest.xml": javaTestPassed,
			"build/reports/jacoco/test/html/index.html":        "<html></html>\n",
		}}
		_, xmlOut, xmlCode := acceptance.JavaRun(t, acceptance.JavaRepo(t, map[string]string{"build.gradle": javaTestGradle}), "java-test", map[string][]acceptance.Response{"gradle": {noXML}})
		req.Then(t, "0001.1.2.7.1", "10cdf9c6", func(a req.A) {
			acceptance.Sees(t, out, code, 1, "JaCoCo is not in the Maven build at .; add it:")
			acceptance.JavaSeesSnippet(t, out, "<artifactId>jacoco-maven-plugin</artifactId>")
		})
		req.Then(t, "0001.1.2.7.2", "e9b05980", func(a req.A) {
			acceptance.Sees(t, xmlOut, xmlCode, 1, "no JaCoCo XML report after the Gradle build at . ran its tests; add:")
			acceptance.JavaSeesSnippet(t, xmlOut, "xml.required = true")
		})
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 a Maven build runs verify once then reads the Surefire reports then the Failsafe reports", func(t *testing.T) {
		t.Parallel()
		maven := javaTestMavenRun(map[string]string{
			"target/surefire-reports/TEST-example.AppTest.xml": javaTestPassed,
			"target/failsafe-reports/TEST-example.AppIT.xml":   javaTestPassed,
			"target/site/jacoco/jacoco.xml":                    javaTestFull,
		})
		rp, out, code := acceptance.JavaRun(t, acceptance.JavaRepo(t, map[string]string{"pom.xml": acceptance.JavaPom(javaTestMaven)}), "java-test", map[string][]acceptance.Response{"mvn": {maven}})
		req.Then(t, "0001.2.1.1.1", "74d5578b", func(a req.A) {
			acceptance.JavaCalledOnce(t, rp, "mvn")
			a.True(rp.Called("mvn", "verify"), "mvn calls = %q, want verify", rp.Calls("mvn"))
		})
		req.Then(t, "0001.2.1.1.2", "e035cbb4", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		})
	})
	t.Run("0001.2.2 a Gradle build reruns its tests with the build cache", func(t *testing.T) {
		t.Parallel()
		rp, out, code := acceptance.JavaRun(t, acceptance.JavaRepo(t, map[string]string{"build.gradle": javaTestGradle}), "java-test", map[string][]acceptance.Response{"gradle": {javaTestGradleRun}}, "GRADLE_USER_HOME="+t.TempDir())
		// Clean rules rerun the tasks on every Gradle version; --rerun needs 7.6.
		// The build cache then restores the reports of unchanged tests.
		req.Then(t, "0001.2.1.2.1", "0cd576a9", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.JavaCalledOnce(t, rp, "gradle")
			a.True(rp.Called("gradle", "cleanTest") && rp.Called("gradle", "cleanJacocoTestReport"), "gradle calls = %q, want clean rules before each task", rp.Calls("gradle"))
		})
		req.Then(t, "0001.2.1.2.2", "a109aca9", func(a req.A) {
			a.True(rp.Called("gradle", "cleanTest test") && !rp.Called("gradle", "--rerun"), "gradle calls = %q, want the tests rerun by clean rules and no --rerun", rp.Calls("gradle"))
		})
		req.Then(t, "0001.2.1.2.3", "2fb1c8f3", func(a req.A) {
			a.True(rp.Called("gradle", "cleanJacocoTestReport jacocoTestReport") && !rp.Called("gradle", "--rerun"), "gradle calls = %q, want JaCoCo rerun by clean rules and no --rerun", rp.Calls("gradle"))
		})
		req.Then(t, "0001.2.1.2.4", "bd5dda75", func(a req.A) {
			a.True(rp.Called("gradle", "--build-cache"), "gradle calls = %q, want --build-cache", rp.Calls("gradle"))
		})
	})
	t.Run("0001.2.3 reports the build cache restores count", func(t *testing.T) {
		t.Parallel()
		// Reports the cache restores are written during the run, so they count.
		req.Then(t, "0001.2.1.3.1", "d48501e5", func(_ req.A) {
			rp, out, code := acceptance.JavaRun(t, acceptance.JavaRepo(t, map[string]string{"build.gradle": javaTestGradle}), "java-test", map[string][]acceptance.Response{"gradle": {javaTestCached}}, "GRADLE_USER_HOME="+t.TempDir())
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			if !rp.Called("gradle", "--build-cache") {
				t.Errorf("gradle calls = %q, want --build-cache", rp.Calls("gradle"))
			}
		})
	})
	t.Run("0001.2.4 a build that does not enable the cache runs without it", func(t *testing.T) {
		t.Parallel()
		// Gradle reads any value but true as off, from the build or the user home.
		user := t.TempDir()
		acceptance.MustDo(t, os.WriteFile(filepath.Join(user, "gradle.properties"), []byte("org.gradle.caching=False\n"), 0o644))
		req.Then(t, "0001.2.1.4.1", "a397e6a4", func(a req.A) {
			javaTestUncached(t, "org.gradle.caching = no\n", javaTestGradleRun, "GRADLE_USER_HOME="+t.TempDir())
		})
		req.Then(t, "0001.2.1.4.2", "b8cbe131", func(a req.A) {
			javaTestUncached(t, "org.gradle.caching=true\n", javaTestGradleRun, "GRADLE_USER_HOME="+user)
		})
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

// javaTestMavenOutcome runs a Maven build that writes one JUnit report and one JaCoCo report.
func javaTestMavenOutcome(t *testing.T, unit, jacoco string) (string, int) {
	t.Helper()
	files := map[string]string{"target/surefire-reports/TEST-example.AppTest.xml": unit}
	if jacoco != "" {
		files["target/site/jacoco/jacoco.xml"] = jacoco
	}
	_, out, code := acceptance.JavaRun(t, acceptance.JavaRepo(t, map[string]string{"pom.xml": acceptance.JavaPom(javaTestMaven)}), "java-test", map[string][]acceptance.Response{"mvn": {javaTestMavenRun(files)}})
	return out, code
}

// javaTestGradleAggregate proves Gradle's aggregate, which lists packages
// without module groups, covers a project when it lists one of its
// sources by declared package.
func javaTestGradleAggregate(t *testing.T) {
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
	javaTestNotJudged(t, out, "core/build/reports/jacoco/test/jacocoTestReport.xml")
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

// javaTestNotJudged fails when the output names any of the reports.
func javaTestNotJudged(t *testing.T, out string, reports ...string) {
	t.Helper()
	req.Asserted(t)
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

// javaTestUncached asserts a Gradle build with props in its
// gradle.properties, run with env, reruns its tests without the build
// cache.
func javaTestUncached(t *testing.T, props string, gradle acceptance.Response, env string) {
	t.Helper()
	repo := acceptance.JavaRepo(t, map[string]string{"build.gradle": javaTestGradle, "gradle.properties": props})
	rp, out, code := acceptance.JavaRun(t, repo, "java-test", map[string][]acceptance.Response{"gradle": {gradle}}, env)
	acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	if !rp.Called("gradle", "cleanTest test") || rp.Called("gradle", "--build-cache") {
		t.Errorf("gradle calls = %q, want the tests run without --build-cache", rp.Calls("gradle"))
	}
}
