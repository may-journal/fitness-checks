// Package requirements proves the checkstyle check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo, with a JDK and the build tools answering saved responses.
package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain lets this test binary stand in for the JDK and the build tools. Run
// as one of them, it answers from saved responses instead of running tests.
func TestMain(m *testing.M) {
	acceptance.ServeStandIns([]string{"java", "mvn", "gradle"})
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 warning fails with position", func(t *testing.T) {
		t.Parallel()
		repo := checkstyleRepo(t, checkstyleMaven, map[string]string{checkstyleConfig: checkstyleCeiling})
		_, out, code := javaRun(t, repo, "checkstyle", map[string][]acceptance.Response{"mvn": {checkstyleWarning(repo)}})
		acceptance.Sees(t, out, code, 1, javaSource+":3:5 - Name 'Foo' must match pattern '^[a-z][a-zA-Z0-9]*$' (MethodName)")
		// Gradle prints a warning only when its task runs, never when it holds
		// the task up to date, so each task must rerun.
		repo = javaRepo(t, map[string]string{"build.gradle": "plugins { id 'java'; id 'checkstyle' }\n", checkstyleConfig: checkstyleCeiling})
		rp, out, code := javaRun(t, repo, "checkstyle", map[string][]acceptance.Response{"gradle": {checkstyleGradleWarning(repo)}})
		acceptance.Sees(t, out, code, 1, javaSource+":3:5 - Name 'Foo' must match pattern '^[a-z][a-zA-Z0-9]*$' (MethodName)")
		if !rp.Called("gradle", "checkstyleMain --rerun", "checkstyleTest --rerun") {
			t.Errorf("gradle calls = %q, want both Checkstyle tasks forced to rerun", rp.Calls("gradle"))
		}
	})
	t.Run("0001.1.2 clean sources pass checkstyle, counting the Java files", func(t *testing.T) {
		t.Parallel()
		repo := checkstyleRepo(t, checkstyleMaven, map[string]string{
			checkstyleConfig:                    checkstyleCeiling,
			"config/checkstyle/strict.xml":      checkstyleCeiling,
			"src/test/resources/no-ceiling.xml": checkstyleNoCeiling,
		})
		clean := acceptance.Response{Match: []string{"checkstyle:check", "checkstyle.config.location=", "/" + checkstyleConfig}, Stdout: checkstyleClean.Stdout}
		_, out, code := javaRun(t, repo, "checkstyle", map[string][]acceptance.Response{"mvn": {clean}})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "checkstyle", "1")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.1 build without checkstyle plugin gets snippet", func(t *testing.T) {
		t.Parallel()
		repo := checkstyleRepo(t, "", map[string]string{checkstyleConfig: checkstyleCeiling})
		_, out, code := javaRun(t, repo, "checkstyle", map[string][]acceptance.Response{"mvn": {checkstyleClean}})
		acceptance.Sees(t, out, code, 1, "Checkstyle is not in the Maven build at .; add it:")
		javaSeesSnippet(t, out, "<artifactId>maven-checkstyle-plugin</artifactId>")
	})
	t.Run("0001.2.2 build without checkstyle config fails", func(t *testing.T) {
		t.Parallel()
		repo := checkstyleRepo(t, checkstyleMaven, nil)
		_, out, code := javaRun(t, repo, "checkstyle", map[string][]acceptance.Response{"mvn": {checkstyleClean}})
		acceptance.Sees(t, out, code, 1, "no Checkstyle config sets a complexity ceiling; add "+checkstyleConfig+" holding")
	})
	t.Run("0001.2.3 config without complexity ceiling fails naming the config", func(t *testing.T) {
		t.Parallel()
		repo := checkstyleRepo(t, checkstyleMaven, map[string]string{checkstyleConfig: checkstyleNoCeiling})
		_, out, code := javaRun(t, repo, "checkstyle", map[string][]acceptance.Response{"mvn": {checkstyleClean}})
		acceptance.Sees(t, out, code, 1, checkstyleConfig+" sets no complexity ceiling; add inside TreeWalker:")
		javaSeesSnippet(t, out, `<module name="CyclomaticComplexity">`)
	})
	t.Run("0001.2.4 maven build with hinted config passes", func(t *testing.T) {
		t.Parallel()
		_, out, _ := javaRun(t, checkstyleRepo(t, checkstyleMaven, nil), "checkstyle", map[string][]acceptance.Response{"mvn": {checkstyleClean}})
		hinted := checkstyleHinted(t, out)
		if !strings.HasPrefix(hinted, "<!DOCTYPE module PUBLIC \"-//Checkstyle//DTD Checkstyle Configuration 1.3//EN\"") {
			t.Errorf("the hinted config lacks the document type line Maven needs:\n%s", hinted)
		}
		repo := checkstyleRepo(t, checkstyleMaven, map[string]string{checkstyleConfig: hinted})
		_, out, code := javaRun(t, repo, "checkstyle", map[string][]acceptance.Response{"mvn": {checkstyleClean}})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	_, out, code := javaRun(t, acceptance.Example(t, "happyRepo", nil), "checkstyle", nil)
	acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	acceptance.PassedWithNoFiles(t, out, "checkstyle")
}

// javaJDK is what `java -version` prints, on stderr, for a running JDK 21.
var javaJDK = acceptance.Response{
	Match:  []string{"-version"},
	Stderr: "openjdk version \"21.0.2\" 2024-01-16\nOpenJDK Runtime Environment (build 21.0.2+13)\n",
}

// javaSource is a Java source file every Java example repo holds.
const javaSource = "src/main/java/App.java"

// javaPom is a Maven build file whose build declares plugins.
func javaPom(plugins string) string {
	return "<project><modelVersion>4.0.0</modelVersion><groupId>com.example</groupId>" +
		"<artifactId>app</artifactId><version>1</version><build><plugins>" + plugins + "</plugins></build></project>\n"
}

// javaRepo makes happyRepo with a Java source and files written over it.
func javaRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	all := map[string]string{javaSource: "class App {}\n"}
	for name, body := range files {
		all[name] = body
	}
	return acceptance.Example(t, "happyRepo", all)
}

// javaRun runs `fitness-install -- <check>` on repo with a JDK 21 and build
// tools answering saved responses, JAVA_HOME unset so the JDK on PATH runs.
// env adds to the run's environment.
func javaRun(t *testing.T, repo, check string, saved map[string][]acceptance.Response, env ...string) (acceptance.Replays, string, int) {
	t.Helper()
	all := map[string][]acceptance.Response{"java": {javaJDK}}
	for name, responses := range saved {
		all[name] = responses
	}
	rp := acceptance.StandIn(t, all)
	out, code := acceptance.Fitness(t, repo, append(append(rp.Env, "JAVA_HOME="), env...), check)
	return rp, out, code
}

// javaSeesSnippet asserts the user saw a build file snippet, which the
// results table may break inside a long word, by comparing without spaces.
func javaSeesSnippet(t *testing.T, out, snippet string) {
	t.Helper()
	squeeze := func(s string) string { return strings.Join(strings.Fields(acceptance.Read(s)), "") }
	if !strings.Contains(squeeze(out), squeeze(snippet)) {
		t.Errorf("output is missing the snippet %q:\n%s", snippet, out)
	}
}

// checkstyleMaven declares Checkstyle in a Maven build.
const checkstyleMaven = "<plugin><groupId>org.apache.maven.plugins</groupId><artifactId>maven-checkstyle-plugin</artifactId></plugin>"

// checkstyleConfig is where the Checkstyle config sits.
const checkstyleConfig = "config/checkstyle/checkstyle.xml"

// checkstyleCeiling is a Checkstyle config that caps complexity.
const checkstyleCeiling = `<?xml version="1.0"?>
<module name="Checker">
  <module name="TreeWalker">
    <module name="CyclomaticComplexity"><property name="max" value="5"/></module>
  </module>
</module>
`

// checkstyleNoCeiling is a Checkstyle config with rules but no ceiling.
const checkstyleNoCeiling = `<?xml version="1.0"?>
<module name="Checker">
  <module name="TreeWalker">
    <module name="MethodName"/>
  </module>
</module>
`

// checkstyleClean is Maven's answer for a Checkstyle run with no finding.
var checkstyleClean = acceptance.Response{
	Match:  []string{"checkstyle:check"},
	Stdout: "[INFO] You have 0 Checkstyle violations.\n[INFO] BUILD SUCCESS\n",
}

// checkstyleWarning is Maven's answer for one warning in the Java source of
// repo, printed with its absolute path as the plugin does.
func checkstyleWarning(repo string) acceptance.Response {
	return acceptance.Response{
		Match: []string{"checkstyle:check", "checkstyle.violationSeverity=warning"},
		Stdout: "[INFO] --- checkstyle:3.6.0:check (default-cli) @ app ---\n" +
			"[WARN] " + filepath.Join(repo, javaSource) + ":[3,5] (naming) MethodName: Name 'Foo' must match pattern '^[a-z][a-zA-Z0-9]*$'.\n" +
			"[ERROR] Failed to execute goal org.apache.maven.plugins:maven-checkstyle-plugin:3.6.0:check (default-cli) on project app: You have 1 Checkstyle violation. -> [Help 1]\n",
		Exit: 1,
	}
}

// checkstyleRepo is a Maven repo declaring plugins, with files written over it.
func checkstyleRepo(t *testing.T, plugins string, files map[string]string) string {
	t.Helper()
	all := map[string]string{"pom.xml": javaPom(plugins)}
	for name, body := range files {
		all[name] = body
	}
	return javaRepo(t, all)
}

// checkstyleGradleWarning is Gradle's answer for one Checkstyle warning in
// repo's Java source, passing the build because warnings do not fail it.
func checkstyleGradleWarning(repo string) acceptance.Response {
	return acceptance.Response{
		Match: []string{"checkstyleMain"},
		Stdout: "> Task :checkstyleMain\n" +
			"[ant:checkstyle] [WARN] " + filepath.Join(repo, javaSource) + ":3:5: Name 'Foo' must match pattern '^[a-z][a-zA-Z0-9]*$'. [MethodName]\n" +
			"\nBUILD SUCCESSFUL in 3s\n",
	}
}

// checkstyleHinted reads, from fitness's output, the config the hint asks
// a repo without one to add.
func checkstyleHinted(t *testing.T, out string) string {
	t.Helper()
	_, after, ok := strings.Cut(acceptance.Read(out), checkstyleConfig+" holding ")
	config, _, found := strings.Cut(after, "</module></module></module>")
	if !ok || !found {
		t.Fatalf("output gives no config to add:\n%s", out)
	}
	return config + "</module></module></module>\n"
}
