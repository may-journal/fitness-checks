// Package requirements proves the issue-checklist check's docs/requirements as
// a user meets them. Each test runs the installer at FITNESS_INSTALL in a copied
// example repo, with an issue body for the check to read.
package requirements

import (
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

func TestMain(m *testing.M) {
	acceptance.Main(m)
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 unchecked items fail, each quoted", func(t *testing.T) {
		t.Parallel()
		out, code := checklist(t, "- [x] Write the doc.\n- [ ] Write the tests.\n- [ ] Ship it.\n")
		acceptance.Sees(t, out, code, 1, `unchecked item: "Write the tests."`, `unchecked item: "Ship it."`)
	})
	t.Run("0001.1.2 ticked items and no list pass", func(t *testing.T) {
		t.Parallel()
		out, code := checklist(t, "- [x] Write the doc.\n- [x] Write the tests.\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWith(t, out, "issue-checklist", "1")
		out, code = checklist(t, "Just a description.\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.3 items inside fenced code never count", func(t *testing.T) {
		t.Parallel()
		out, code := checklist(t, "- [x] Write the doc.\n\n```\n- [ ] An example item.\n```\n")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.4 no input passes with zero files checked", func(t *testing.T) {
		t.Parallel()
		out, code := acceptance.Fitness(t, acceptance.Example(t, "happyRepo", nil), nil, "issue-checklist")
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "issue-checklist")
	})
}

// checklist runs `fitness-install -- issue-checklist` on an issue body kept in
// the repo as issue.md.
func checklist(t *testing.T, body string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", map[string]string{"issue.md": body})
	return acceptance.Fitness(t, repo, nil, "issue-checklist", "--body-file", "issue.md")
}
