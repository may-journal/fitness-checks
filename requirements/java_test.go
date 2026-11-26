// cspell:ignore mvnw

package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// The Java fixtures live in internal/acceptance, shared with the Java checks'
// own requirement tests.
var (
	javaJDK         = acceptance.JavaJDK
	javaParent      = acceptance.JavaParent
	javaWrapper     = acceptance.JavaWrapper
	javaSeesSnippet = acceptance.JavaSeesSnippet
)

const javaSource = acceptance.JavaSource

// javaRepo makes happyRepo with a Java source and files written over it.
func javaRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	return acceptance.JavaRepo(t, files)
}

// javaRun runs `fitness-install -- <check>` on repo with a JDK 21 and build
// tools answering saved responses. env adds to the run's environment.
func javaRun(t *testing.T, repo, check string, saved map[string][]response, env ...string) (replays, string, int) {
	t.Helper()
	rp, out, code := acceptance.JavaRun(t, repo, check, saved, env...)
	return replays{Replays: rp, env: rp.Env}, out, code
}

// javaCalledOnce asserts the named build tool ran exactly once.
func javaCalledOnce(t *testing.T, rp replays, name string) {
	t.Helper()
	acceptance.JavaCalledOnce(t, rp.Replays, name)
}
