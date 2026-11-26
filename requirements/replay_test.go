// cspell:ignore mvnw gradlew

package requirements

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// standIns are the external programs tests replace with saved responses: gh,
// the JavaScript and Swift tools, the JDK, Java build tools and Docker.
var standIns = []string{"gh", "eslint", "prettier", "biome", "vitest", "swiftlint", "java", "mvn", "mvnw", "gradle", "gradlew", "docker"}

// The stand-ins answer from saved responses, as acceptance documents.
type response = acceptance.Response

var logCall = acceptance.LogCall

// replays is a test's stand-ins, with the environment that puts them first on
// PATH as env.
type replays struct {
	acceptance.Replays
	env []string
}

// standIn puts a stand-in for each named program first on PATH, answering from
// saved responses.
func standIn(t *testing.T, saved map[string][]response) replays {
	t.Helper()
	r := acceptance.StandIn(t, saved)
	return replays{Replays: r, env: r.Env}
}

// calls returns every call the named stand-in received, in order.
func (r replays) calls(name string) []string { return r.Calls(name) }

// called reports whether some call to name holds every want.
func (r replays) called(name string, want ...string) bool { return r.Called(name, want...) }

// count returns how many calls to name hold every want.
func (r replays) count(name string, want ...string) int { return r.Count(name, want...) }

// TestMain lets this test binary stand in for those programs. Run under one
// of their names, it answers from saved responses instead of running tests.
func TestMain(m *testing.M) {
	acceptance.ServeStandIns(standIns)
	if name := filepath.Base(os.Args[0]); strings.HasPrefix(name, timedPrefix) {
		os.Exit(timedCheck(name, os.Args[1:]))
	}
	acceptance.KeepHostGitOut()
	os.Exit(m.Run())
}
