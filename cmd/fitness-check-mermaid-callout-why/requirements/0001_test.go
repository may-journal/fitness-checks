// Package requirements proves the mermaid-callout-why check's
// docs/requirements as a user meets them. Each test runs the installer at
// FITNESS_INSTALL in a copied example repo.
package requirements

import (
	"strings"
	"testing"

	"github.com/may-journal/fitness-checks/internal/acceptance"
)

// TestMain keeps the host's git config out of the example repos.
func TestMain(m *testing.M) { acceptance.Main(m) }

const (
	noWhy = "| # | Description |\n| --- | --- |\n| 1 | an edge |"
	head  = "| # | Description | Why |\n| --- | --- | --- |"
)

// why runs the check on happyRepo with files written over it.
func why(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", files)
	return acceptance.Fitness(t, repo, nil, append([]string{"mermaid-callout-why"}, args...)...)
}

// doc is a doc with a mermaid diagram followed by a table made of rows; the
// table's header sits on line 10.
func doc(rows ...string) string {
	return "# Design\n\n```mermaid\nflowchart LR\n    a -->|1| b\n```\n\nThe callouts:\n\n" +
		strings.Join(rows, "\n") + "\n"
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1.1.1 it fails naming the table line", func(t *testing.T) {
		t.Parallel()
		out, code := why(t, map[string]string{"docs/design.md": doc(noWhy)})
		acceptance.Sees(t, out, code, 1, `docs/design.md: callout table at line 10 is missing a "Why" column`)
		out, code = why(t, map[string]string{"docs/design.md": "# Design\n\n" + noWhy + "\n"})
		acceptance.Sees(t, out, code, 1, `docs/design.md: callout table at line 3 is missing a "Why" column`)
	})
	t.Run("0001.1.1.2.1 it fails naming the row with its line", func(t *testing.T) {
		t.Parallel()
		out, code := why(t, map[string]string{"docs/design.md": doc(head, "| [2] | an edge |  |")})
		acceptance.Sees(t, out, code, 1, `docs/design.md: callout table row [2] has an empty "Why" cell (line 10)`)
	})
	t.Run("0001.1.1.2.2 a row with no Why cell fails the same way", func(t *testing.T) {
		t.Parallel()
		out, code := why(t, map[string]string{"docs/design.md": doc(head, "| 3 | an edge |")})
		acceptance.Sees(t, out, code, 1, `docs/design.md: callout table row 3 has an empty "Why" cell (line 10)`)
	})
	t.Run("0001.1.1.3.1 the table passes", func(t *testing.T) {
		t.Parallel()
		d := doc("| # | Description | WHY |\n| --- | --- | --- |", "| 1 | an edge | it carries the request |")
		out, code := why(t, map[string]string{"docs/design.md": d})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
	})
	t.Run("0001.1.1.4.1 the row passes", func(t *testing.T) {
		t.Parallel()
		out, code := why(t, map[string]string{"docs/design.md": doc(head, "| n/a | an edge |  |")})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed · 1 files scanned")
	})
	t.Run("0001.1.1.5.1 it fails naming `(description)`", func(t *testing.T) {
		t.Parallel()
		out, code := why(t, map[string]string{"body.md": "## Design\n\n" + noWhy + "\n"}, "--body-file", "body.md")
		acceptance.Sees(t, out, code, 1, `(description): callout table at line 3 is missing a "Why" column`)
	})
	plain := func(t *testing.T) (string, int) {
		return why(t, map[string]string{"docs/plain.md": "# Plain\n\n| Name | Value |\n| --- | --- |\n| a | b |\n"})
	}
	t.Run("0001.1.1.6.1 it is skipped", func(t *testing.T) {
		t.Parallel()
		out, code := plain(t)
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.1.6.2 it is not counted", func(t *testing.T) {
		t.Parallel()
		out, _ := plain(t)
		acceptance.PassedWithNoFiles(t, out, "mermaid-callout-why")
	})
}
