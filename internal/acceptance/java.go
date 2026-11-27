package acceptance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/req"
)

// JavaJDK is what `java -version` prints, on stderr, for a running JDK 21.
var JavaJDK = Response{
	Match:  []string{"-version"},
	Stderr: "openjdk version \"21.0.2\" 2024-01-16\nOpenJDK Runtime Environment (build 21.0.2+13)\n",
}

// JavaSource is a Java source file every Java example repo holds.
const JavaSource = "src/main/java/App.java"

// JavaParent is a Maven build that declares plugins and lists modules.
func JavaParent(plugins string, modules ...string) string {
	list := "<modules>"
	for _, m := range modules {
		list += "<module>" + m + "</module>"
	}
	return strings.Replace(JavaPom(plugins), "<build>", list+"</modules><build>", 1)
}

// JavaRepo makes happyRepo with a Java source and files written over it.
func JavaRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	return Example(t, "happyRepo", With(map[string]string{JavaSource: "class App {}\n"}, files))
}

// JavaRun runs `fitness-install -- <check>` on repo with a JDK 21 and build
// tools answering saved responses, JAVA_HOME unset so the JDK on PATH runs.
// env adds to the run's environment.
func JavaRun(t *testing.T, repo, check string, saved map[string][]Response, env ...string) (Replays, string, int) {
	t.Helper()
	all := map[string][]Response{"java": {JavaJDK}}
	for name, responses := range saved {
		all[name] = responses
	}
	rp := StandIn(t, all)
	out, code := Fitness(t, repo, append(append(rp.Env, "JAVA_HOME="), env...), check)
	return rp, out, code
}

// JavaWrapper puts the stand-in at the root of repo as the build wrapper, to
// answer from its saved responses.
func JavaWrapper(t *testing.T, repo, name string) {
	t.Helper()
	self, err := os.Executable()
	MustDo(t, err)
	MustDo(t, os.Symlink(self, filepath.Join(repo, name)))
}

// JavaSeesSnippet asserts the user saw a build file snippet, which the
// results table may break inside a long word, by comparing without spaces.
func JavaSeesSnippet(t *testing.T, out, snippet string) {
	t.Helper()
	req.Asserted(t)
	squeeze := func(s string) string { return strings.Join(strings.Fields(Read(s)), "") }
	if !strings.Contains(squeeze(out), squeeze(snippet)) {
		t.Errorf("output is missing the snippet %q:\n%s", snippet, out)
	}
}

// JavaCalledOnce asserts the named build tool ran exactly once.
func JavaCalledOnce(t *testing.T, rp Replays, name string) {
	t.Helper()
	req.Asserted(t)
	if calls := rp.Calls(name); len(calls) != 1 {
		t.Errorf("%s calls = %q, want one", name, calls)
	}
}
