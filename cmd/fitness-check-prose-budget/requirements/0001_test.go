// Package requirements proves the prose-budget check's docs/requirements as a
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

// docs are requirement 0001 in the repo root and in package pkgs/adder.
var docs = []string{"docs/requirements/0001-adds-numbers.md", "pkgs/adder/docs/requirements/0001-adds-numbers.md"}

// chained runs prose-budget on a repo whose doc chains steps, under the
// .fitnessrc.json requirements key given. A package doc's repo names pkgs/* as
// package roots.
func chained(t *testing.T, requirements, doc string, steps int) (string, int) {
	t.Helper()
	rc := `{"checkVersions": {}` + requirements + "}\n"
	if strings.HasPrefix(doc, "pkgs/") {
		rc = `{"packages": ["pkgs/*"]` + requirements + "}\n"
	}
	repo := acceptance.Example(t, "happyRepo", map[string]string{
		".fitnessrc.json": rc,
		doc:               acceptance.DocWith(acceptance.Chain("0001.1", steps)),
	})
	return acceptance.Fitness(t, repo, nil, "prose-budget")
}

// chainsWithin fails unless each doc's chain of limit steps passes and one
// more step fails naming the limit.
func chainsWithin(t *testing.T, requirements string, limit int) {
	t.Helper()
	for _, doc := range docs {
		out, code := chained(t, requirements, doc, limit)
		acceptance.Sees(t, out, code, 0)
		out, code = chained(t, requirements, doc, limit+1)
		acceptance.Sees(t, out, code, 1, fmt.Sprintf("%s: a list has more than %d items at level 3", doc, limit))
	}
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	chainsWithin(t, `, "requirements": {"maxSteps": 5}`, 5)
}

func Test0001_2(t *testing.T) {
	t.Parallel()
	chainsWithin(t, "", 7)
	// A chain holds two steps at least, so a smaller limit keeps the default.
	chainsWithin(t, `, "requirements": {"maxSteps": 1}`, 7)
}
