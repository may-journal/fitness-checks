// cspell:ignore mvnw diffplug

// Package requirements proves the spotless check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo, with the JDK and build tools answering saved responses.
package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
	"github.com/may-journal/fitness-checks/req"
)

func TestMain(m *testing.M) {
	acceptance.Main(m, "java", "mvn", "mvnw", "gradle")
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 unformatted files fail naming each", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.1.1", "cfda2072", func(_ req.A) {
			repo := acceptance.JavaRepo(t, map[string]string{"pom.xml": acceptance.JavaPom(spotlessMaven)})
			_, out, code := acceptance.JavaRun(t, repo, "spotless", map[string][]acceptance.Response{"mvn": {spotlessViolations}})
			acceptance.Sees(t, out, code, 1, acceptance.JavaSource+" is not formatted; run spotless:apply")
		})
	})
	t.Run("0001.1.2 clean builds pass and run once from each build root", func(t *testing.T) {
		t.Parallel()
		modules := map[string]string{
			"pom.xml":                         acceptance.JavaParent(spotlessMaven, "core", "web"),
			"core/pom.xml":                    acceptance.JavaPom(""),
			"core/src/main/java/Core.java":    "class Core {}\n",
			"core/src/test/resources/pom.xml": "<project/>\n",
			"web/pom.xml":                     acceptance.JavaPom(""),
			"web/src/main/java/web/Page.java": "class Page {}\n",
		}
		rp, out, code := acceptance.JavaRun(t, acceptance.JavaRepo(t, modules), "spotless", map[string][]acceptance.Response{"mvn": {spotlessClean}})
		// A nested build no pom lists as a module is a build root of its own.
		modules["tools/pom.xml"] = acceptance.JavaPom(spotlessMaven)
		rp2, out2, code2 := acceptance.JavaRun(t, acceptance.JavaRepo(t, modules), "spotless", map[string][]acceptance.Response{"mvn": {spotlessClean}})
		req.Then(t, "0001.1.2.2.1", "66148298", func(_ req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWith(t, out, "spotless", "3")
		})
		req.Then(t, "0001.1.2.2.2", "2410823a", func(a req.A) {
			acceptance.JavaCalledOnce(t, rp, "mvn")
			acceptance.Sees(t, out2, code2, 0, "All 1 checks passed")
			a.Equal(2, len(rp2.Calls("mvn")))
		})
	})
	t.Run("0001.1.3 the build wrapper runs and Gradle runs without a daemon", func(t *testing.T) {
		t.Parallel()
		repo := acceptance.JavaRepo(t, map[string]string{"pom.xml": acceptance.JavaPom(spotlessMaven)})
		self, err := os.Executable()
		acceptance.MustDo(t, err)
		acceptance.MustDo(t, os.Symlink(self, filepath.Join(repo, "mvnw")))
		rp, out, code := acceptance.JavaRun(t, repo, "spotless", map[string][]acceptance.Response{"mvn": {spotlessClean}, "mvnw": {spotlessClean}})
		rp2, out2, code2 := acceptance.JavaRun(t, spotlessGradleRepo(t), "spotless", map[string][]acceptance.Response{"gradle": {spotlessGradleClean}})
		req.Then(t, "0001.1.2.3.1", "691904a6", func(a req.A) {
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.JavaCalledOnce(t, rp, "mvnw")
			a.Equal(0, len(rp.Calls("mvn")))
		})
		req.Then(t, "0001.1.2.3.2", "860a5d42", func(a req.A) {
			acceptance.Sees(t, out2, code2, 0, "All 1 checks passed")
			a.True(rp2.Called("gradle", "--no-daemon", "spotlessCheck"), "gradle calls = %q, want spotlessCheck with --no-daemon", rp2.Calls("gradle"))
		})
	})
	t.Run("0001.1.4 a build without Spotless gets the plugin snippet", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.4.1", "99ba78b0", func(_ req.A) {
			repo := acceptance.JavaRepo(t, map[string]string{
				"pom.xml":                         acceptance.JavaPom(""),
				"tools/pom.xml":                   acceptance.JavaPom(spotlessMaven),
				"src/test/resources/fake/pom.xml": acceptance.JavaPom(spotlessMaven),
			})
			_, out, code := acceptance.JavaRun(t, repo, "spotless", map[string][]acceptance.Response{"mvn": {spotlessClean}})
			acceptance.Sees(t, out, code, 1, "Spotless is not in the Maven build at .; add it:")
			acceptance.JavaSeesSnippet(t, out, "<artifactId>spotless-maven-plugin</artifactId>")
			if strings.Contains(acceptance.Read(out), "build at tools") {
				t.Errorf("the nested build declares Spotless, so only the root may fail:\n%s", out)
			}
		})
	})
	t.Run("0001.1.5 a repo with no build passes with no files", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.5.1", "857dd667", func(_ req.A) {
			_, out, code := acceptance.JavaRun(t, acceptance.Example(t, "happyRepo", nil), "spotless", nil)
			acceptance.Sees(t, out, code, 0, "All 1 checks passed")
			acceptance.PassedWithNoFiles(t, out, "spotless")
		})
	})
	t.Run("0001.1.6 a failure naming no file shows the build's error", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.6.1", "610904a2", func(_ req.A) {
			_, out, code := acceptance.JavaRun(t, spotlessGradleRepo(t), "spotless", map[string][]acceptance.Response{"gradle": {spotlessGradleNoJDK}})
			acceptance.Sees(t, out, code, 1, "Could not determine the dependencies of task ':spotlessJava'.")
			acceptance.JavaSeesSnippet(t, out, "> Cannot find a Java installation on your machine")
			acceptance.JavaSeesSnippet(t, out, "matching: {languageVersion=17,")
		})
	})
	t.Run("0001.1.7 no running JDK shows the install hint", func(t *testing.T) {
		t.Parallel()
		req.Then(t, "0001.1.2.7.1", "b2a40761", func(_ req.A) {
			repo := acceptance.JavaRepo(t, map[string]string{"pom.xml": acceptance.JavaPom(spotlessMaven)})
			out, code := acceptance.Fitness(t, repo, []string{acceptance.PathWithout(t, "java"), "JAVA_HOME="}, "spotless")
			acceptance.Sees(t, out, code, 1, "Java not installed: install a JDK")
		})
	})
}

// spotlessMaven declares Spotless in a Maven build.
const spotlessMaven = "<plugin><groupId>com.diffplug.spotless</groupId><artifactId>spotless-maven-plugin</artifactId></plugin>"

// spotlessGradleBuild is a Gradle build that declares Spotless.
const spotlessGradleBuild = "plugins { id 'java'; id 'com.diffplug.spotless' version '8.10.4' }\n"

// spotlessClean is the build's answer when Spotless finds nothing to rewrite.
var spotlessClean = acceptance.Response{Match: []string{"spotless:check"}, Stdout: "[INFO] BUILD SUCCESS\n"}

// spotlessGradleClean is Gradle's answer for a clean spotlessCheck.
var spotlessGradleClean = acceptance.Response{Match: []string{"spotlessCheck"}, Stdout: "> Task :spotlessCheck\n\nBUILD SUCCESSFUL in 2s\n"}

// spotlessViolations is Maven's answer when Spotless would rewrite the
// Java source: the file under the violations header, then its diff.
var spotlessViolations = acceptance.Response{
	Match: []string{"spotless:check"},
	Stdout: "[INFO] Scanning for projects...\n" +
		"[ERROR] Failed to execute goal com.diffplug.spotless:spotless-maven-plugin:3.10.4:check (default-cli) on project app: The following files had format violations:\n" +
		"[ERROR]     " + acceptance.JavaSource + "\n" +
		"[ERROR]         @@ -1 +1 @@\n" +
		"[ERROR]         -class App {  }\n" +
		"[ERROR]         +class App {}\n" +
		"[ERROR] Run 'mvn spotless:apply' to fix these violations.\n",
	Exit: 1,
}

// spotlessGradleNoJDK is Gradle's answer, saved from a real run, when the
// toolchain needs a JDK the machine lacks.
var spotlessGradleNoJDK = acceptance.Response{
	Match: []string{"spotlessCheck"},
	Stdout: "FAILURE: Build failed with an exception.\n\n* What went wrong:\n" +
		"Could not determine the dependencies of task ':spotlessJava'.\n" +
		"> Could not resolve all dependencies for configuration ':compileClasspath'.\n" +
		"   > Failed to calculate the value of task ':compileJava' property 'javaCompiler'.\n" +
		"      > Cannot find a Java installation on your machine (Mac OS X 27.0.1 aarch64) matching: {languageVersion=17, vendor=any vendor, implementation=vendor-specific, nativeImageCapable=false}. Toolchain download repositories have not been configured.\n" +
		"\n* Try:\n> Run with --stacktrace option to get the stack trace.\n\nBUILD FAILED in 1s\n",
	Exit: 1,
}

// spotlessGradleRepo is a Java repo whose Gradle build declares Spotless.
func spotlessGradleRepo(t *testing.T) string {
	t.Helper()
	return acceptance.JavaRepo(t, map[string]string{"build.gradle": spotlessGradleBuild})
}
