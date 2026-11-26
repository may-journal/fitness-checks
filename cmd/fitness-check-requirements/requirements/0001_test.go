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

// verifyChainLimit fails unless a chain of limit steps passes the limit. One
// more step must fail with every want, such as the chain and the limit.
func verifyChainLimit(t *testing.T, requirements string, limit int, want ...string) {
	t.Helper()
	if out := chainRun(t, requirements, limit); strings.Contains(out, "0001.1 chains") {
		t.Errorf("a chain of %d steps must pass a limit of %d:\n%s", limit, limit, out)
	}
	out := chainRun(t, requirements, limit+1)
	for _, w := range want {
		if !strings.Contains(out, fmt.Sprintf(w, limit+1, limit)) {
			t.Errorf("output is missing %q:\n%s", fmt.Sprintf(w, limit+1, limit), out)
		}
	}
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	const set = `, "requirements": {"maxSteps": 3}`
	t.Run("0001.1.2.1.1 it fails naming the chain", func(t *testing.T) {
		t.Parallel()
		verifyChainLimit(t, set, 3, "0001.1 chains %[1]d steps")
	})
	t.Run("0001.1.2.1.2 it names the limit", func(t *testing.T) {
		t.Parallel()
		verifyChainLimit(t, set, 3, "at most %[2]d")
	})
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	t.Run("0001.2.2.1.1 the default limit of seven steps applies", func(t *testing.T) {
		t.Parallel()
		verifyChainLimit(t, "", 7, "0001.1 chains %d steps, at most %d")
		// A chain holds two steps at least, so a smaller limit keeps the default.
		verifyChainLimit(t, `, "requirements": {"maxSteps": 1}`, 7, "0001.1 chains %d steps, at most %d")
	})
}
