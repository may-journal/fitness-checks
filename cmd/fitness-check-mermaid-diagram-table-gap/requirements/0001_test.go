// Package requirements proves the mermaid-diagram-table-gap check's
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

// The numbered diagram, its legend, the caption, and a matching callout table
// that each Given arranges around the diagram-to-table gap.
const (
	diagram = "```mermaid\nContainer(app, \"1\", \"UI\")\nRel(app, api, \"2\")\n```"
	legend  = "```mermaid\nPerson(p, \"Person\")\nSystem(s, \"System\")\n```"
	caption = "Numbers on nodes and arrows match the callout table."
	table   = "| # | Description | Why |\n| --- | --- | --- |\n| 1 | the app | because |\n| 2 | talks to | because |"
	prose   = "Four containers sit inside one app."
)

// gap runs the check on happyRepo with files written over it.
func gap(t *testing.T, files map[string]string, args ...string) (string, int) {
	t.Helper()
	repo := acceptance.Example(t, "happyRepo", files)
	return acceptance.Fitness(t, repo, nil, append([]string{"mermaid-diagram-table-gap"}, args...)...)
}

// doc joins parts into a doc, one blank line between each.
func doc(parts ...string) string {
	return "# Doc\n\n" + strings.Join(parts, "\n\n") + "\n"
}

// finding is the message a reader sees for prose at line in file.
func finding(file, line string) string {
	return file + ": prose between the diagram and its callout table (line " + line + ") — move detail into the callout table"
}

func Test0001_1(t *testing.T) {
	t.Parallel()
	t.Run("0001.1.1 a paragraph between diagram and table fails", func(t *testing.T) {
		t.Parallel()
		out, code := gap(t, map[string]string{"doc.md": doc(diagram, prose, table)})
		acceptance.Sees(t, out, code, 1, finding("doc.md", "8"))
	})
	t.Run("0001.1.2 a paragraph between legend and table fails", func(t *testing.T) {
		t.Parallel()
		out, code := gap(t, map[string]string{"doc.md": doc(diagram, legend, prose, table)})
		acceptance.Sees(t, out, code, 1, finding("doc.md", "13"))
	})
	t.Run("0001.1.3 a matching caption alone passes", func(t *testing.T) {
		t.Parallel()
		out, code := gap(t, map[string]string{
			"a.md": doc(diagram, legend, caption, table),
			"b.md": doc(diagram, "Numbers on classes and relationships match the callout table.", table),
			"c.md": doc(diagram, `Numbers match the callout table. "Screen" is a terminal state.`, table),
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.4 a paragraph after the caption fails but not the caption", func(t *testing.T) {
		t.Parallel()
		out, code := gap(t, map[string]string{"doc.md": doc(diagram, legend, caption, prose, table)})
		acceptance.Sees(t, out, code, 1, finding("doc.md", "15"))
		if strings.Contains(acceptance.Read(out), acceptance.Read(finding("doc.md", "13"))) {
			t.Errorf("the caption at line 13 must not be flagged:\n%s", out)
		}
	})
	t.Run("0001.1.5 prose outside the gap passes", func(t *testing.T) {
		t.Parallel()
		out, code := gap(t, map[string]string{
			"a.md": doc("Intro prose before the diagram.", diagram, caption, table, "Trailing prose after the table."),
			"b.md": doc("Prose above.", "```mermaid\nContainer(app, \"UI\")\n```", "Prose below."),
		})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
	})
	t.Run("0001.1.6 a file with no diagram is skipped", func(t *testing.T) {
		t.Parallel()
		out, code := gap(t, map[string]string{"doc.md": doc("Just prose, nothing to see here.")})
		acceptance.Sees(t, out, code, 0, "All 1 checks passed")
		acceptance.PassedWithNoFiles(t, out, "mermaid-diagram-table-gap")
	})
	t.Run("0001.1.7 a PR body paragraph in the gap fails naming the description", func(t *testing.T) {
		t.Parallel()
		body := strings.Join([]string{diagram, prose, table}, "\n\n") + "\n"
		out, code := gap(t, map[string]string{"body.md": body}, "--body-file", "body.md")
		acceptance.Sees(t, out, code, 1, finding("(description)", "6"))
	})
}
