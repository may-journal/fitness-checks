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

// numberedDoc is a requirement doc of limit numbered acceptances, each
// followed by a Removed line, so it holds twice the list limit.
func numberedDoc(limit int) string {
	reqs := ""
	for n := 1; n <= 2*limit; n += 2 {
		reqs += fmt.Sprintf("%d. Adds pair %d\n   1. GIVEN two numbers\n      1. WHEN I add them\n         1. THEN I get their sum\n", n, n)
		reqs += fmt.Sprintf("%d. Removed in 2026.11.25.1230: folded into 0001.%d\n", n+1, n)
	}
	return acceptance.DocWith(reqs)
}

// prosed runs prose-budget on a repo whose root requirement doc is doc.
func prosed(t *testing.T, doc string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", map[string]string{
		".fitnessrc.json": `{"checkVersions": {}}` + "\n",
		docs[0]:           doc,
	})
	return acceptance.Fitness(t, repo, nil, "prose-budget")
}

// items is n numbered items of word at indent, numbered from 1.
func items(indent, word string, n int) string {
	out := ""
	for i := 1; i <= n; i++ {
		out += fmt.Sprintf("%s%d. %s %d\n", indent, i, word, i)
	}
	return out
}

func Test0001_3(t *testing.T) {
	t.Parallel()
	t.Run("0001.3.1 Removed lines count toward no limit", func(t *testing.T) {
		t.Parallel()
		out, code := prosed(t, numberedDoc(8))
		acceptance.Sees(t, out, code, 0)
		// The same count of acceptances alone passes the limit and fails.
		out, code = prosed(t, strings.ReplaceAll(numberedDoc(8), ". Removed in 2026.11.25.1230: folded into", ". Adds"))
		acceptance.Sees(t, out, code, 1, docs[0]+": a list has more than 8 items at level 1")
	})
	t.Run("0001.3.2 GIVENs or a WHEN's THENs past half the list limit fail at that level", func(t *testing.T) {
		t.Parallel()
		within := "1. Adds\n" + items("   ", "GIVEN two numbers", 4) + "      1. WHEN I add them\n" + items("         ", "THEN I get their sum", 4)
		out, code := prosed(t, acceptance.DocWith(within))
		acceptance.Sees(t, out, code, 0)
		out, code = prosed(t, acceptance.DocWith("1. Adds\n"+items("   ", "GIVEN two numbers", 5)+"      1. WHEN I add them\n         1. THEN I get their sum\n"))
		acceptance.Sees(t, out, code, 1, docs[0]+": a list has more than 4 items at level 2")
		out, code = prosed(t, acceptance.DocWith("1. Adds\n   1. GIVEN two numbers\n      1. WHEN I add them\n"+items("         ", "THEN I get their sum", 5)))
		acceptance.Sees(t, out, code, 1, docs[0]+": a list has more than 4 items at level 4")
	})
}
