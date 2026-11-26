// Package requirements proves the requirements check's docs/requirements as a
// user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo.
package requirements

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.KeepHostGitOut()
	acceptance.TakeTurn()
	os.Exit(m.Run())
}

// rootDoc is the repo root's requirement 0001.
const rootDoc = "docs/requirements/0001-adds-numbers.md"

// chainRun runs requirements on a Go repo whose root doc chains steps, under
// the .fitnessrc.json requirements key given.
func chainRun(t *testing.T, requirements string, steps int) string {
	t.Helper()
	rc := `{"policy": "external", "checks": ["requirements"]` + requirements + "}\n"
	repo := acceptance.Example(t, "happyRepo", acceptance.With(acceptance.GoRepo(), map[string]string{
		".fitnessrc.json": rc,
		rootDoc:           acceptance.DocWith(acceptance.Chain("0001.1", steps)),
	}))
	out, _ := acceptance.Fitness(t, repo, nil, "requirements")
	return acceptance.Read(out)
}

// chainsWithin fails unless a chain of steps passes the limit and one more
// step fails naming both.
func chainsWithin(t *testing.T, requirements string, limit int) {
	t.Helper()
	if out := chainRun(t, requirements, limit); strings.Contains(out, "0001.1 chains") {
		t.Errorf("a chain of %d steps must pass a limit of %d:\n%s", limit, limit, out)
	}
	want := fmt.Sprintf("0001.1 chains %d steps, at most %d", limit+1, limit)
	if out := chainRun(t, requirements, limit+1); !strings.Contains(out, want) {
		t.Errorf("output is missing %q:\n%s", want, out)
	}
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	chainsWithin(t, `, "requirements": {"maxSteps": 3}`, 3)
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	chainsWithin(t, "", 7)
	// A chain holds two steps at least, so a smaller limit keeps the default.
	chainsWithin(t, `, "requirements": {"maxSteps": 1}`, 7)
}
